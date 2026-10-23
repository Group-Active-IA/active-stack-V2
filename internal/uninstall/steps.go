package uninstall

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Group-Active-IA/active-stack/internal/backup"
	"github.com/Group-Active-IA/active-stack/internal/filemerge"
	"github.com/Group-Active-IA/active-stack/internal/harness/command"
	"github.com/Group-Active-IA/active-stack/internal/harness/config"
	roledefsinstaller "github.com/Group-Active-IA/active-stack/internal/harness/roledefs"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

// ─────────────────────────────────────────────────────────────────
// snapshotStep — Prepare stage; creates the uninstall-time backup snapshot.
// ─────────────────────────────────────────────────────────────────

type snapshotStep struct {
	id          string
	paths       []string
	snapDir     string
	snapCreate  func(dir string, paths []string) (backup.Manifest, error)
	manifestOut *backup.Manifest // written after Run so Apply steps can rollback
}

func (s *snapshotStep) ID() string { return s.id }

func (s *snapshotStep) Run() error {
	manifest, err := s.snapCreate(s.snapDir, s.paths)
	if err != nil {
		return fmt.Errorf("uninstall snapshot: %w", err)
	}
	if s.manifestOut != nil {
		*s.manifestOut = manifest
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────
// markerRemovalStep — config harnesses: removes the marker section.
// ─────────────────────────────────────────────────────────────────

// markerRemovalFn is the backing function for reading a file and calling
// InjectMarkdownSection with empty content (idempotent removal).
// It is a package-level variable so tests can inject a fake.
var markerRemovalFn = func(path, sectionID string) error {
	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist — no marker to remove; treat as no-op.
			return nil
		}
		return fmt.Errorf("read instructions file %q: %w", path, err)
	}
	updated := filemerge.InjectMarkdownSection(string(existing), sectionID, "")
	if updated == string(existing) {
		// No change — section was absent; no-op.
		return nil
	}
	_, err = filemerge.WriteFileAtomic(path, []byte(updated), 0o644)
	return err
}

// stalePurgeFn removes every active-stack-marked section the current installer no
// longer owns (legacy/renamed sections from older layouts: persona,
// engram-protocol, strict-tdd-mode, …). It mirrors the install-time cleanup so
// uninstall leaves no orphaned blocks behind, and reuses config.PurgeStaleSections
// as the single source of truth for the owned/stale policy.
//
// It is a package-level variable so tests can swap it out.
var stalePurgeFn = func(path string) error {
	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist — nothing to purge; no-op.
			return nil
		}
		return fmt.Errorf("read instructions file %q: %w", path, err)
	}
	updated := config.PurgeStaleSections(string(existing))
	if updated == string(existing) {
		// No stale sections present; no-op.
		return nil
	}
	_, err = filemerge.WriteFileAtomic(path, []byte(updated), 0o644)
	return err
}

// primaryAgentRemovalFn removes the agent.<id> entry from a settings JSON file
// (e.g. opencode.json), preserving all other user config. It is the mirror of
// the install-time primary-agent registration. No-op when the file is missing,
// malformed, or already has no such entry. It is a package-level variable so
// tests can inject a fake.
var primaryAgentRemovalFn = func(settingsPath, agentID string) error {
	if settingsPath == "" {
		return nil
	}
	raw, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read settings file %q: %w", settingsPath, err)
	}

	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		// Malformed settings — leave it untouched rather than risk clobbering
		// the user's file. Restore-from-snapshot remains the safety net.
		return nil
	}

	agentSection, ok := root["agent"].(map[string]any)
	if !ok {
		return nil // nothing to remove
	}
	if _, present := agentSection[agentID]; !present {
		return nil // no-op
	}
	delete(agentSection, agentID)
	// Drop the agent object entirely if it is now empty, leaving no footprint.
	if len(agentSection) == 0 {
		delete(root, "agent")
	}

	encoded, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal settings %q: %w", settingsPath, err)
	}
	encoded = append(encoded, '\n')
	if _, err := filemerge.WriteFileAtomic(settingsPath, encoded, 0o644); err != nil {
		return fmt.Errorf("write settings %q: %w", settingsPath, err)
	}
	return nil
}

type markerRemovalStep struct {
	h        model.Harness
	adapters []AgentAdapter
	homeDir  string
	manifest *backup.Manifest
}

func (s *markerRemovalStep) ID() string                     { return "marker:" + s.h.ID }
func (s *markerRemovalStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *markerRemovalStep) Run() error {
	for _, a := range s.adapters {
		path := a.InstructionsPath(s.homeDir)
		// Remove this harness's own section first…
		if err := markerRemovalFn(path, s.h.ID); err != nil {
			return fmt.Errorf("marker removal for harness %q on agent %q: %w", s.h.ID, a.Agent(), err)
		}
		// …then purge any legacy active-stack sections from older layouts so the
		// uninstall leaves no orphaned blocks behind.
		if err := stalePurgeFn(path); err != nil {
			return fmt.Errorf("stale section purge for harness %q on agent %q: %w", s.h.ID, a.Agent(), err)
		}
		// Primary-agent delivery lives in the settings JSON, not the markdown
		// instructions file — remove it there too.
		if a.ConfigDelivery() == model.ConfigDeliveryPrimaryAgent {
			if err := primaryAgentRemovalFn(a.SettingsPath(s.homeDir), s.h.ID); err != nil {
				return fmt.Errorf("primary-agent removal for harness %q on agent %q: %w", s.h.ID, a.Agent(), err)
			}
		}
	}
	return nil
}

func (s *markerRemovalStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// ─────────────────────────────────────────────────────────────────
// skillRemovalStep — skill harnesses: removes the skill directory.
// ─────────────────────────────────────────────────────────────────

// skillRemovalFn is the backing function for removing a skill directory.
// It is a package-level variable so tests can inject a fake.
var skillRemovalFn = func(path string) error {
	return os.RemoveAll(path)
}

type skillRemovalStep struct {
	h        model.Harness
	adapters []AgentAdapter
	homeDir  string
	manifest *backup.Manifest
}

func (s *skillRemovalStep) ID() string                     { return "skill-removal:" + s.h.ID }
func (s *skillRemovalStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *skillRemovalStep) Run() error {
	for _, a := range s.adapters {
		skillPath := filepath.Join(a.SkillsDir(s.homeDir), s.h.ID)
		if err := skillRemovalFn(skillPath); err != nil {
			return fmt.Errorf("skill removal for harness %q on agent %q: %w", s.h.ID, a.Agent(), err)
		}
	}
	return nil
}

func (s *skillRemovalStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// ─────────────────────────────────────────────────────────────────
// commandRemovalStep — command harnesses: removes the per-agent command file.
// ─────────────────────────────────────────────────────────────────

// embeddedCommandsFS is the fs.FS with the command assets, set by the binary
// entry point via WithEmbeddedCommandsFS. The uninstall engine needs it to know
// WHICH files the installer wrote: esa lista sale de los assets, nunca de lo que
// hay en disco (ver commandRelPathsForVariant).
var embeddedCommandsFS fs.FS

// WithEmbeddedCommandsFS sets the embedded commands FS used to compute which
// command files a targeted uninstall removes. Pass assets.CommandsFS from the
// binary entry point, next to install.WithEmbeddedCommandsFS.
//
// Mirrors the install-side global on purpose: C-32 fue el bug de que el binario
// no llamaba a la version de install y el instalador reventaba con
// "commandsFS is nil". Misma forma, mismo lugar de wiring, mismo test.
func WithEmbeddedCommandsFS(f fs.FS) {
	embeddedCommandsFS = f
}

// commandRelPathsForVariant returns the relative paths of the command files a
// targeted uninstall may remove for this variant. Es el mismo walk que corrio el
// instalador, sobre los mismos assets embebidos.
//
// 2026-09-09 — por que sale de los assets y NO de listar el disco:
// listar <commandsDir> y borrar todo lo que haya adentro tiene una ventaja real
// (limpiaria comandos de una version anterior que ya no estan embebidos) y un
// costo que la anula: ese directorio es COMPARTIDO con los comandos que el
// usuario escribe a mano. En ~/.claude/commands/ conviven hoy /opsx:explore,
// /opsx:propose, /opsx:apply, /opsx:archive y los /sdd-*, que active-stack nunca
// instalo; un uninstall que camina el disco se los lleva. Un huerfano molesta,
// un archivo del usuario borrado no se perdona, asi que la lista sale de lo que
// este binario sabe que escribio.
//
// Lo que esta decision NO cubre, dicho para que se sepa: un comando instalado por
// una version vieja, cuyo asset ya no existe en este binario, queda huerfano y
// nadie lo limpia. Distinguirlo de un archivo del usuario pediria un registro de
// lo instalado, y hoy no hay ninguno (solo manifiestos de backup).
//
// Con commandsFS nil cae al path unico de RelPathForVariant, que es el
// comportamiento previo a este cambio: si el entry point se olvida de llamar a
// WithEmbeddedCommandsFS (C-32 fue exactamente ese bug), el uninstall degrada a
// borrar starter-add en vez de convertirse en un no-op silencioso.
func commandRelPathsForVariant(variantKey string) []string {
	if embeddedCommandsFS == nil {
		relPath := command.RelPathForVariant(variantKey)
		if relPath == "" {
			return nil
		}
		return []string{relPath}
	}
	relPaths, err := command.RelPathsForVariant(embeddedCommandsFS, variantKey)
	if err != nil {
		return nil
	}
	return relPaths
}

// commandRemovalFn is the backing function for removing a single command file.
// It is a package-level variable so tests can inject a fake.
var commandRemovalFn = func(path string) error {
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil // already absent — idempotent no-op
	}
	return err
}

type commandRemovalStep struct {
	h        model.Harness
	adapters []AgentAdapter
	homeDir  string
	manifest *backup.Manifest
}

func (s *commandRemovalStep) ID() string                     { return "command-removal:" + s.h.ID }
func (s *commandRemovalStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *commandRemovalStep) Run() error {
	for _, a := range s.adapters {
		commandsDir := a.CommandsDir(s.homeDir)
		if commandsDir == "" {
			// Adapter does not support commands; skip silently.
			continue
		}
		// Unknown variant (o variante sin assets) → lista vacia → skip silencioso.
		for _, relPath := range commandRelPathsForVariant(a.VariantKey()) {
			commandPath := filepath.Join(commandsDir, relPath)
			if err := commandRemovalFn(commandPath); err != nil {
				return fmt.Errorf("command removal for harness %q on agent %q: %w", s.h.ID, a.Agent(), err)
			}
		}
	}
	return nil
}

func (s *commandRemovalStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// ─────────────────────────────────────────────────────────────────
// permissionsRemovalStep — permissions harness: restore from snapshot.
// ─────────────────────────────────────────────────────────────────

type permissionsRemovalStep struct {
	h        model.Harness
	adapters []AgentAdapter
	homeDir  string
	manifest *backup.Manifest
}

func (s *permissionsRemovalStep) ID() string                     { return "permissions-removal:" + s.h.ID }
func (s *permissionsRemovalStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *permissionsRemovalStep) Run() error {
	// (TBD) There is no dedicated mechanism yet to selectively strip just the
	// permissions overlay back out of settings.json — the reverse of
	// internal/harness/config/permissions/install.go's mergeJSONFile/agentOverlay.
	//
	// This used to call restoreFn(*s.manifest) here (a copy of Rollback's own
	// body — see below), which was a real bug: s.manifest is the ENTIRE
	// uninstall-time snapshot shared across every harness being removed in
	// this run, not something scoped to permissions/settings.json alone.
	// Calling it from Run silently reverted every OTHER harness this uninstall
	// had already removed (e.g. markerRemovalStep's CLAUDE.md edit) back to
	// its pre-uninstall state, making a targeted uninstall of anything appear
	// to do nothing whenever "permissions" (the mandatory security floor) was
	// part of the selection — which it always is.
	//
	// Run is a safe no-op until a properly scoped removal exists: the
	// settings.json permissions overlay is left in place rather than risk
	// clobbering sibling harnesses' removals. Rollback (below) still
	// legitimately restores the full manifest — that IS its job when the
	// overall uninstall pipeline fails and everything applied so far must be
	// undone.
	return nil
}

func (s *permissionsRemovalStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// ─────────────────────────────────────────────────────────────────
// roleDefinitionsRemovalStep — unlike permissionsRemovalStep, this one is
// NOT a no-op: roledefs.Uninstall discriminates ownership for real (hash
// recorded in its own manifest vs. current bytes), so Run() safely removes
// only what this harness owns and preserves everything personal or
// user-modified — see internal/harness/roledefs/uninstall.go.
// ─────────────────────────────────────────────────────────────────

var roleDefinitionsUninstallFn = func(adapters []AgentAdapter, homeDir string) error {
	res, err := roledefsinstaller.DefaultResources()
	if err != nil {
		return err
	}
	_, err = roledefsinstaller.Uninstall(toRoleDefinitionsAdapters(adapters), homeDir, res)
	return err
}

type roleDefinitionsRemovalStep struct {
	h        model.Harness
	adapters []AgentAdapter
	homeDir  string
	manifest *backup.Manifest
}

func (s *roleDefinitionsRemovalStep) ID() string                     { return "role-definitions-removal:" + s.h.ID }
func (s *roleDefinitionsRemovalStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *roleDefinitionsRemovalStep) Run() error {
	return roleDefinitionsUninstallFn(s.adapters, s.homeDir)
}

func (s *roleDefinitionsRemovalStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// roleDefinitionsAdapterShim and toRoleDefinitionsAdapters mirror the
// install-package versions (internal/install/steps.go) — AgentsDir is
// recovered via a runtime type assertion against the concrete adapter
// rather than widening the shared AgentAdapter interface.
type roleDefinitionsAdapterShim struct {
	AgentAdapter
	underlying AgentAdapter
}

func (s roleDefinitionsAdapterShim) AgentsDir(homeDir string) string {
	if rd, ok := s.underlying.(interface{ AgentsDir(string) string }); ok {
		return rd.AgentsDir(homeDir)
	}
	return ""
}

func toRoleDefinitionsAdapters(adapters []AgentAdapter) []roledefsinstaller.AgentAdapter {
	out := make([]roledefsinstaller.AgentAdapter, len(adapters))
	for i, a := range adapters {
		out[i] = roleDefinitionsAdapterShim{AgentAdapter: a, underlying: a}
	}
	return out
}

// ─────────────────────────────────────────────────────────────────
// restoreStep — StrategyRestore: replays the install-time backup.
// ─────────────────────────────────────────────────────────────────

type restoreStep struct {
	installManifest backup.Manifest
}

func (s *restoreStep) ID() string { return "restore-from-backup" }

func (s *restoreStep) Run() error {
	return restoreFn(s.installManifest)
}

// restoreStep intentionally does not implement RollbackStep: StrategyRestore
// has no Prepare snapshot to roll back to (the install-time manifest IS the
// recovery point), so there is nothing to revert.

// ─────────────────────────────────────────────────────────────────
// externalSkipStep — external harnesses: recorded no-op.
// ─────────────────────────────────────────────────────────────────

type externalSkipStep struct {
	h model.Harness
}

func (s *externalSkipStep) ID() string { return "external-skip:" + s.h.ID }
func (s *externalSkipStep) Run() error { return nil }

type codexMCPRemovalStep struct {
	h        model.Harness
	adapters []AgentAdapter
	homeDir  string
	manifest *backup.Manifest
}

func (s *codexMCPRemovalStep) ID() string                     { return "external-codex-mcp:" + s.h.ID }
func (s *codexMCPRemovalStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *codexMCPRemovalStep) Run() error {
	name := s.h.ID
	if s.h.External != nil && s.h.External.MCP != nil && s.h.External.MCP.Name != "" {
		name = s.h.External.MCP.Name
	}
	for _, adapter := range s.adapters {
		if adapter.Agent() != model.AgentCodex {
			continue
		}
		path := adapter.SettingsPath(s.homeDir)
		raw, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("read Codex MCP config %q: %w", path, err)
		}
		merged, err := filemerge.RemoveTOMLSection(string(raw), "mcp_servers."+name)
		if err != nil {
			return fmt.Errorf("remove Codex MCP %q: %w", name, err)
		}
		if _, err := filemerge.WriteFileAtomic(path, []byte(merged), 0o644); err != nil {
			return fmt.Errorf("write Codex MCP config %q: %w", path, err)
		}
	}
	return nil
}

func (s *codexMCPRemovalStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// ─────────────────────────────────────────────────────────────────
// Package-level function vars (testseams; real implementations below)
// ─────────────────────────────────────────────────────────────────

// snapshotCreate is the backing function for creating snapshots. It is a
// package-level variable so tests can inject a fake without reopening backup.
var snapshotCreate = func(snapshotDir string, paths []string) (backup.Manifest, error) {
	return backup.NewSnapshotter().Create(snapshotDir, paths)
}

// restoreFn is the backing function for restoring from a snapshot.
// It is a package-level variable so tests can inject a fake.
var restoreFn = func(m backup.Manifest) error {
	return backup.RestoreService{}.Restore(m)
}

// manifestReceiver is implemented by steps that need the snapshot manifest
// for rollback. The manifest is set after Prepare runs.
type manifestReceiver interface {
	setManifest(m *backup.Manifest)
}
