package install

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Group-Active-IA/active-stack/internal/backup"
	cmdinstaller "github.com/Group-Active-IA/active-stack/internal/harness/command"
	cfginstaller "github.com/Group-Active-IA/active-stack/internal/harness/config"
	perminstaller "github.com/Group-Active-IA/active-stack/internal/harness/config/permissions"
	extinstaller "github.com/Group-Active-IA/active-stack/internal/harness/external"
	roledefsinstaller "github.com/Group-Active-IA/active-stack/internal/harness/roledefs"
	skillinstaller "github.com/Group-Active-IA/active-stack/internal/harness/skill"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/pipeline"
	"github.com/Group-Active-IA/active-stack/internal/system"
)

// ─────────────────────────────────────────────────────────────────
// snapshotStep — Prepare stage; creates the backup snapshot.
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
		return fmt.Errorf("snapshot: %w", err)
	}
	if s.manifestOut != nil {
		*s.manifestOut = manifest
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────
// externalStep — wraps external.Install.
// ─────────────────────────────────────────────────────────────────

// externalInstallFn is the backing function for external installation.
// It is a package-level variable so tests can inject a fake.
var externalInstallFn = func(
	ctx context.Context,
	h model.Harness,
	profile system.PlatformProfile,
	adapters []extinstaller.AgentAdapter,
	homeDir string,
	downloadFn extinstaller.DownloadEventFunc,
) (extinstaller.Result, error) {
	return extinstaller.Install(ctx, h, profile, adapters, homeDir, downloadFn)
}

type externalStep struct {
	h          model.Harness
	adapters   []AgentAdapter
	homeDir    string
	profile    system.PlatformProfile
	manifest   *backup.Manifest
	onDownload DownloadEventFunc

	// result is what the installer reported, captured by Run so Rollback knows
	// what this run actually created. It used to be discarded
	// (`_, err := externalInstallFn(...)`), which is why rollback had no idea a
	// binary had been written and reported a clean success while leaving it.
	result extinstaller.Result
	ran    bool
}

func (s *externalStep) ID() string                     { return "external:" + s.h.ID }
func (s *externalStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *externalStep) Run() error {
	ext := toExternalAdapters(s.adapters)
	var downloadFn extinstaller.DownloadEventFunc
	if s.onDownload != nil {
		stepID := s.ID()
		downloadFn = func(ev extinstaller.DownloadEvent) {
			s.onDownload(DownloadEvent{
				StepID:  stepID,
				Type:    string(ev.Type),
				URL:     ev.URL,
				Message: ev.Message,
			})
		}
	}
	res, err := externalInstallFn(context.Background(), s.h, s.profile, ext, s.homeDir, downloadFn)
	if err != nil {
		return err
	}
	s.result = res
	s.ran = true
	return nil
}

func (s *externalStep) Rollback() error {
	// A binary this run downloaded into active-stack's own install dir is ours
	// to delete: nothing else on the machine put it there. Anything installed
	// by npm or brew, and anything that was already present, is not — see
	// Unreverted.
	if s.ran && s.result.Downloaded && s.result.BinaryPath != "" {
		if err := os.Remove(s.result.BinaryPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove downloaded binary %q: %w", s.result.BinaryPath, err)
		}
	}
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// Unreverted implements pipeline.UnrevertedReporter. It names the side effects
// this step performed that rollback has no way to undo, so the installer can
// report a partial rollback instead of claiming the machine is clean.
func (s *externalStep) Unreverted() []string {
	if !s.ran || s.result.AlreadyInstalled || s.h.External == nil {
		return nil
	}
	pkg := s.h.External.Pkg
	if pkg == "" {
		pkg = s.h.ID
	}

	switch s.h.External.Method {
	case "npm":
		// Deleting the shim would leave a half-removed package behind, so we
		// deliberately do not touch it and say so instead.
		return []string{fmt.Sprintf("npm package %q stays installed globally (undo with: npm uninstall -g %s)", pkg, pkg)}
	case "homebrew":
		if s.result.Downloaded || s.result.BinaryPath == "" {
			// Downloaded by us and already removed in Rollback, or never installed.
			return nil
		}
		return []string{fmt.Sprintf("brew formula %q stays installed (undo with: brew uninstall %s)", pkg, pkg)}
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────
// skillStep — wraps skill.NewInstaller(...).Install.
// ─────────────────────────────────────────────────────────────────

// skillInstallFn is the backing function for skill installation.
// It is a package-level variable so tests can inject a fake.
var skillInstallFn = func(
	runner skillRunner,
	embeddedFS fs.FS,
	ctx context.Context,
	h model.Harness,
	adapters []skillinstaller.AgentAdapter,
	homeDir, backupDir string,
) ([]skillinstaller.Result, error) {
	ins := skillinstaller.NewInstaller(runner, embeddedFS)
	return ins.Install(ctx, h, adapters, homeDir, backupDir)
}

type skillStep struct {
	h          model.Harness
	adapters   []AgentAdapter
	homeDir    string
	backupDir  string
	embeddedFS fs.FS
	runner     skillRunner
	manifest   *backup.Manifest
	// bestEffort mirrors h.BestEffort: when true a Run() failure is soft
	// (warning emitted, nil returned) so the pipeline does not abort/rollback.
	bestEffort bool
	// onProgress is the optional progress callback forwarded from Options.
	// When non-nil it receives a warning ProgressEvent on best-effort failure.
	onProgress pipeline.ProgressFunc
}

func (s *skillStep) ID() string                     { return "skill:" + s.h.ID }
func (s *skillStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *skillStep) Run() error {
	skill := toSkillAdapters(s.adapters)
	_, err := skillInstallFn(s.runner, s.embeddedFS, context.Background(), s.h, skill, s.homeDir, s.backupDir)
	if err == nil {
		return nil
	}
	if !s.bestEffort {
		return err
	}
	// Best-effort: emit a degraded event (C-32: distinct from hard failure) and
	// return nil so the pipeline continues. The degraded status is NOT the same
	// as StepStatusFailed — it must not be conflated with an abort-causing failure.
	warningMsg := fmt.Sprintf("[best-effort] skill %q install failed (continuing): %v", s.h.ID, err)
	if s.onProgress != nil {
		s.onProgress(pipeline.ProgressEvent{
			StepID: s.ID(),
			Stage:  pipeline.StageApply,
			Status: pipeline.StepStatusDegraded,
			Err:    fmt.Errorf("%s", warningMsg),
		})
	} else {
		fmt.Fprintln(os.Stderr, warningMsg)
	}
	return nil
}

func (s *skillStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// ─────────────────────────────────────────────────────────────────
// configStep — wraps config.Install.
// ─────────────────────────────────────────────────────────────────

// configInstallFn is the backing function for config installation.
// It is a package-level variable so tests can inject a fake.
var configInstallFn = func(
	h model.Harness,
	adapters []cfginstaller.AgentAdapter,
	homeDir string,
) (cfginstaller.Result, error) {
	return cfginstaller.Install(h, adapters, homeDir)
}

type configStep struct {
	h        model.Harness
	adapters []AgentAdapter
	homeDir  string
	manifest *backup.Manifest
}

func (s *configStep) ID() string                     { return "config:" + s.h.ID }
func (s *configStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *configStep) Run() error {
	cfg := toConfigAdapters(s.adapters)
	_, err := configInstallFn(s.h, cfg, s.homeDir)
	return err
}

func (s *configStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// ─────────────────────────────────────────────────────────────────
// permissionsStep — wraps permissions.Install.
// ─────────────────────────────────────────────────────────────────

// permissionsInstallFn is the backing function for permissions installation.
// It is a package-level variable so tests can inject a fake.
var permissionsInstallFn = func(
	homeDir string,
	adapters []perminstaller.PermissionsAdapter,
	tier model.PermissionTier,
) (perminstaller.Result, error) {
	return perminstaller.Install(homeDir, adapters, tier)
}

type permissionsStep struct {
	h        model.Harness
	adapters []AgentAdapter
	homeDir  string
	tier     model.PermissionTier
	manifest *backup.Manifest
}

func (s *permissionsStep) ID() string                     { return "permissions:" + s.h.ID }
func (s *permissionsStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *permissionsStep) Run() error {
	perm := toPermissionsAdapters(s.adapters)
	_, err := permissionsInstallFn(s.homeDir, perm, s.tier)
	return err
}

func (s *permissionsStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// ─────────────────────────────────────────────────────────────────
// roleDefinitionsStep — wraps roledefs.Install (harness `role-definitions`,
// propose fase 2 paso 4). Dispatched by HarnessID, same pattern as
// "permissions" — see buildHarnessStep in plan.go for why: a HarnessConfig
// harness with no special case falls into the generic toggle-injection
// configStep, which would try to inject the role-*.md resources as if they
// were CLAUDE.md prose instead of materializing them on disk.
//
// roledefs.Install already recovers its own partial writes across adapters
// when one of them fails mid-call (QA-2, see roledefs/doc.go) — this step's
// Rollback() is the OTHER half: the pipeline's cross-step case, where this
// step fully succeeded and a LATER step in the same plan fails. That one
// still needs the pre-Run snapshot, same as permissionsStep.
// ─────────────────────────────────────────────────────────────────

// roleDefinitionsInstallFn is the backing function for role-definitions
// installation. Package-level var so tests can inject a fake.
var roleDefinitionsInstallFn = func(
	adapters []AgentAdapter,
	homeDir, backupDir string,
) (roledefsinstaller.Result, error) {
	res, err := roledefsinstaller.DefaultResources()
	if err != nil {
		return roledefsinstaller.Result{}, err
	}
	return roledefsinstaller.Install(toRoleDefinitionsAdapters(adapters), homeDir, backupDir, res)
}

type roleDefinitionsStep struct {
	h         model.Harness
	adapters  []AgentAdapter
	homeDir   string
	backupDir string
	manifest  *backup.Manifest
}

func (s *roleDefinitionsStep) ID() string                     { return "role-definitions:" + s.h.ID }
func (s *roleDefinitionsStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *roleDefinitionsStep) Run() error {
	_, err := roleDefinitionsInstallFn(s.adapters, s.homeDir, s.backupDir)
	return err
}

func (s *roleDefinitionsStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// roleDefinitionsAdapterShim narrows the full AgentAdapter to
// roledefs.AgentAdapter (Agent + AgentsDir). AgentsDir is not part of the
// install.AgentAdapter interface (low blast radius — see
// internal/harness/roledefs/doc.go), so it is recovered via a runtime type
// assertion against the concrete adapter instead of widening the shared
// interface. An adapter that does not implement it resolves to "" —
// capability not supported, same as a real Codex/OpenCode adapter today.
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
// mcpWriteStep — writes a single model.MCP to the agent's project
// config path and merges it idempotently (C-28 D5).
// Governance ALTO: backup before write + Rollback() restores snapshot.
// ─────────────────────────────────────────────────────────────────

// writeMCPEntry writes a project-scoped MCP entry to the given config path
// using the resolved strategy. For MCPStrategySingleFileMerge (Claude project),
// it delegates to external.WriteMCPProjectEntry which handles backup + merge.
// The snapshotDir is derived from the configPath's parent directory.
func writeMCPEntry(mcp model.MCP, configPath string, strategy model.MCPStrategy) error {
	switch strategy {
	case model.MCPStrategySingleFileMerge, model.MCPStrategyMergeIntoSettings:
		// Both project merge strategies use the same write path:
		// backup + MergeJSONObjects + WriteFileAtomic.
		snapshotDir := filepath.Join(filepath.Dir(configPath), ".active-stack", "backups", "mcp", mcp.Name)
		_, err := extinstaller.WriteMCPProjectEntry(mcp, configPath, snapshotDir)
		return err
	case model.MCPStrategyMergeIntoTOML:
		snapshotDir := filepath.Join(filepath.Dir(configPath), ".active-stack", "backups", "mcp", mcp.Name)
		_, err := extinstaller.WriteCodexMCPProjectEntry(mcp, configPath, snapshotDir)
		return err
	default:
		// MCPStrategySeparateFile (machine target) is handled by the existing
		// harness-based installMCP flow, not by starter MCP wiring.
		return fmt.Errorf("writeMCPEntry: strategy %d not supported for starter MCPs", strategy)
	}
}

// mcpWriteFn is the backing function for writing a project MCP entry.
// It is a package-level variable so tests can inject a fake.
var mcpWriteFn = func(
	mcp model.MCP,
	configPath string,
	strategy model.MCPStrategy,
) error {
	return writeMCPEntry(mcp, configPath, strategy)
}

type mcpWriteStep struct {
	id         string
	mcp        model.MCP
	configPath string
	strategy   model.MCPStrategy
	manifest   *backup.Manifest
}

func (s *mcpWriteStep) ID() string                     { return s.id }
func (s *mcpWriteStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *mcpWriteStep) Run() error {
	return mcpWriteFn(s.mcp, s.configPath, s.strategy)
}

func (s *mcpWriteStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}

// ─────────────────────────────────────────────────────────────────
// Adapter coercion helpers
// ─────────────────────────────────────────────────────────────────

// toExternalAdapters narrows the full AgentAdapter to external.AgentAdapter.
func toExternalAdapters(adapters []AgentAdapter) []extinstaller.AgentAdapter {
	out := make([]extinstaller.AgentAdapter, len(adapters))
	for i, a := range adapters {
		out[i] = a
	}
	return out
}

// toSkillAdapters narrows the full AgentAdapter to skill.AgentAdapter.
func toSkillAdapters(adapters []AgentAdapter) []skillinstaller.AgentAdapter {
	out := make([]skillinstaller.AgentAdapter, len(adapters))
	for i, a := range adapters {
		out[i] = a
	}
	return out
}

// toConfigAdapters narrows the full AgentAdapter to config.AgentAdapter.
func toConfigAdapters(adapters []AgentAdapter) []cfginstaller.AgentAdapter {
	out := make([]cfginstaller.AgentAdapter, len(adapters))
	for i, a := range adapters {
		out[i] = a
	}
	return out
}

// toPermissionsAdapters narrows the full AgentAdapter to permissions.PermissionsAdapter.
func toPermissionsAdapters(adapters []AgentAdapter) []perminstaller.PermissionsAdapter {
	out := make([]perminstaller.PermissionsAdapter, len(adapters))
	for i, a := range adapters {
		out[i] = a
	}
	return out
}

// toCommandAdapters narrows the full AgentAdapter to command.AgentAdapter.
func toCommandAdapters(adapters []AgentAdapter) []cmdinstaller.AgentAdapter {
	out := make([]cmdinstaller.AgentAdapter, len(adapters))
	for i, a := range adapters {
		out[i] = a
	}
	return out
}

// ─────────────────────────────────────────────────────────────────
// commandStep — materializes a slash-command file per focused agent.
// Added in C-31 (HarnessCommand / TBD-2 option a).
// Governance ALTO: backup before write + Rollback() via snapshot manifest.
// ─────────────────────────────────────────────────────────────────

// commandInstallFn is the backing function for command installation.
// It is a package-level variable so tests can inject a fake without real writes.
var commandInstallFn = func(
	adapters []AgentAdapter,
	homeDir, backupDir string,
) error {
	ins := cmdinstaller.NewInstaller(embeddedCommandsFS)
	cmdAdapters := toCommandAdapters(adapters)
	_, err := ins.Install(cmdAdapters, homeDir, backupDir)
	return err
}

// embeddedCommandsFS is the fs.FS for command assets. It is set via
// WithEmbeddedCommandsFS so the production binary passes assets.CommandsFS.
// When nil, commandInstallFn falls back to an empty FS (nothing installed).
var embeddedCommandsFS fs.FS

type commandStep struct {
	h         model.Harness
	adapters  []AgentAdapter
	homeDir   string
	backupDir string
	manifest  *backup.Manifest
}

func (s *commandStep) ID() string                     { return "command:" + s.h.ID }
func (s *commandStep) setManifest(m *backup.Manifest) { s.manifest = m }

func (s *commandStep) Run() error {
	return commandInstallFn(s.adapters, s.homeDir, s.backupDir)
}

func (s *commandStep) Rollback() error {
	if s.manifest == nil {
		return nil
	}
	return restoreFn(*s.manifest)
}
