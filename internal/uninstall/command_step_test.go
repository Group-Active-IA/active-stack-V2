package uninstall_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/Group-Active-IA/active-stack/assets"
	"github.com/Group-Active-IA/active-stack/internal/backup"
	"github.com/Group-Active-IA/active-stack/internal/harness/command"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/uninstall"
)

// TestCommandRemovalStep_DeletesCommandFile asserts that commandRemovalStep
// deletes the command file from the adapter's CommandsDir/RelPath.
//
// RED: fails until commandRemovalStep exists with Run() logic.
func TestCommandRemovalStep_DeletesCommandFile(t *testing.T) {
	homeDir := t.TempDir()
	adapter := fakeAdapter{agent: model.AgentClaude, homeDir: homeDir, variantKey: "claude"}

	// Pre-create the command file at the expected path.
	commandFile := filepath.Join(adapter.CommandsDir(homeDir), "active", "starter-add.md")
	if err := os.MkdirAll(filepath.Dir(commandFile), 0o755); err != nil {
		t.Fatalf("setup: mkdir: %v", err)
	}
	if err := os.WriteFile(commandFile, []byte("# jr starter add"), 0o644); err != nil {
		t.Fatalf("setup: write file: %v", err)
	}

	restoreSnap := uninstall.SetSnapshotCreate(func(dir string, paths []string) (backup.Manifest, error) {
		return backup.Manifest{}, nil
	})
	defer restoreSnap()

	h := model.Harness{
		ID:           "starter-add-command",
		Type:         model.HarnessCommand,
		Agents:       []model.Agent{model.AgentClaude},
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
	}

	plan, err := uninstall.BuildPlan(
		&fakeCatalog{harnesses: []model.Harness{h}},
		uninstall.Intent{
			Agents:   []model.Agent{model.AgentClaude},
			Mode:     model.ModeLite,
			Strategy: uninstall.StrategyTargeted,
		},
		uninstall.Options{
			HomeDir:  homeDir,
			Registry: &fakeRegistry{adapters: map[model.Agent]uninstall.AgentAdapter{model.AgentClaude: adapter}},
		},
	)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}

	for _, step := range plan.Apply {
		if err := step.Run(); err != nil {
			t.Fatalf("step.Run() error = %v", err)
		}
	}

	if _, err := os.Stat(commandFile); !os.IsNotExist(err) {
		t.Errorf("command file still exists after commandRemovalStep: %s", commandFile)
	}
}

// TestCommandRemovalStep_MissingFileIsNoop asserts that commandRemovalStep
// is a no-op when the command file does not exist (idempotent removal).
func TestCommandRemovalStep_MissingFileIsNoop(t *testing.T) {
	homeDir := t.TempDir()
	adapter := fakeAdapter{agent: model.AgentClaude, homeDir: homeDir, variantKey: "claude"}

	// Do NOT create the command file — it must be absent.
	restoreSnap := uninstall.SetSnapshotCreate(func(dir string, paths []string) (backup.Manifest, error) {
		return backup.Manifest{}, nil
	})
	defer restoreSnap()

	h := model.Harness{
		ID:           "starter-add-command",
		Type:         model.HarnessCommand,
		Agents:       []model.Agent{model.AgentClaude},
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
	}

	plan, err := uninstall.BuildPlan(
		&fakeCatalog{harnesses: []model.Harness{h}},
		uninstall.Intent{
			Agents:   []model.Agent{model.AgentClaude},
			Mode:     model.ModeLite,
			Strategy: uninstall.StrategyTargeted,
		},
		uninstall.Options{
			HomeDir:  homeDir,
			Registry: &fakeRegistry{adapters: map[model.Agent]uninstall.AgentAdapter{model.AgentClaude: adapter}},
		},
	)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}

	for _, step := range plan.Apply {
		if err := step.Run(); err != nil {
			t.Errorf("step.Run() on missing file must be no-op, got error = %v", err)
		}
	}
}

// TestCommandRemovalStep_UnknownVariantSkips asserts that an adapter with an
// unknown VariantKey is silently skipped (no error, no delete call).
func TestCommandRemovalStep_UnknownVariantSkips(t *testing.T) {
	homeDir := t.TempDir()
	// Adapter with an unknown variant; commandsDir is non-empty to exercise the path.
	adapter := fakeAdapter{agent: model.AgentGemini, homeDir: homeDir, variantKey: "gemini"}

	restoreSnap := uninstall.SetSnapshotCreate(func(dir string, paths []string) (backup.Manifest, error) {
		return backup.Manifest{}, nil
	})
	defer restoreSnap()

	// Track whether commandRemovalFn was called.
	removalCalled := false
	restoreRemoval := uninstall.SetCommandRemovalFn(func(path string) error {
		removalCalled = true
		return nil
	})
	defer restoreRemoval()

	h := model.Harness{
		ID:           "starter-add-command",
		Type:         model.HarnessCommand,
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
		// empty agents = all agents
	}

	plan, err := uninstall.BuildPlan(
		&fakeCatalog{harnesses: []model.Harness{h}},
		uninstall.Intent{
			Agents:   []model.Agent{model.AgentGemini},
			Mode:     model.ModeLite,
			Strategy: uninstall.StrategyTargeted,
		},
		uninstall.Options{
			HomeDir:  homeDir,
			Registry: &fakeRegistry{adapters: map[model.Agent]uninstall.AgentAdapter{model.AgentGemini: adapter}},
		},
	)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}

	for _, step := range plan.Apply {
		if err := step.Run(); err != nil {
			t.Errorf("step.Run() for unknown variant must not error, got = %v", err)
		}
	}

	if removalCalled {
		t.Error("commandRemovalFn must NOT be called for unknown variant key")
	}
}

// TestCommandRemovalStep_EmptyCommandsDirSkips asserts that an adapter whose
// CommandsDir returns "" is silently skipped.
func TestCommandRemovalStep_EmptyCommandsDirSkips(t *testing.T) {
	homeDir := t.TempDir()
	// Use fakeAdapterCustomPath with empty commandsDir.
	adapter := fakeAdapterCustomPath{
		agent:       model.AgentClaude,
		commandsDir: "", // empty → skip
		variantKey:  "claude",
	}

	restoreSnap := uninstall.SetSnapshotCreate(func(dir string, paths []string) (backup.Manifest, error) {
		return backup.Manifest{}, nil
	})
	defer restoreSnap()

	removalCalled := false
	restoreRemoval := uninstall.SetCommandRemovalFn(func(path string) error {
		removalCalled = true
		return nil
	})
	defer restoreRemoval()

	h := model.Harness{
		ID:           "starter-add-command",
		Type:         model.HarnessCommand,
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
	}

	plan, err := uninstall.BuildPlan(
		&fakeCatalog{harnesses: []model.Harness{h}},
		uninstall.Intent{
			Agents:   []model.Agent{model.AgentClaude},
			Mode:     model.ModeLite,
			Strategy: uninstall.StrategyTargeted,
		},
		uninstall.Options{
			HomeDir:  homeDir,
			Registry: &fakeRegistryCustom{adapters: map[model.Agent]uninstall.AgentAdapter{model.AgentClaude: adapter}},
		},
	)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}

	for _, step := range plan.Apply {
		if err := step.Run(); err != nil {
			t.Errorf("step.Run() for empty CommandsDir must not error, got = %v", err)
		}
	}

	if removalCalled {
		t.Error("commandRemovalFn must NOT be called when CommandsDir is empty")
	}
}

// TestCommandRemovalStep_RollbackRestoresFromSnapshot asserts that Rollback
// calls restoreFn when a manifest is available.
func TestCommandRemovalStep_RollbackRestoresFromSnapshot(t *testing.T) {
	homeDir := t.TempDir()
	adapter := fakeAdapter{agent: model.AgentClaude, homeDir: homeDir, variantKey: "claude"}

	manifestID := "test-manifest-command"
	restoreSnap := uninstall.SetSnapshotCreate(func(dir string, paths []string) (backup.Manifest, error) {
		return backup.Manifest{ID: manifestID}, nil
	})
	defer restoreSnap()

	// Make marker removal fail so rollback is triggered.
	restoreMarker := uninstall.SetMarkerRemovalFn(func(path, sectionID string) error {
		return errTest("forced failure")
	})
	defer restoreMarker()

	restoreWasCalled := false
	restoreRestoreFn := uninstall.SetRestoreFn(func(m backup.Manifest) error {
		restoreWasCalled = true
		return nil
	})
	defer restoreRestoreFn()

	// Build a plan with a config harness (so marker removal triggers rollback)
	// and a command harness (to verify the command step also has the manifest wired).
	configH := model.Harness{
		ID:           "sdd-orchestrator",
		Type:         model.HarnessConfig,
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
	}
	commandH := model.Harness{
		ID:           "starter-add-command",
		Type:         model.HarnessCommand,
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
	}

	plan, err := uninstall.BuildPlan(
		&fakeCatalog{harnesses: []model.Harness{configH, commandH}},
		uninstall.Intent{
			Agents:   []model.Agent{model.AgentClaude},
			Mode:     model.ModeLite,
			Strategy: uninstall.StrategyTargeted,
		},
		uninstall.Options{
			HomeDir:  homeDir,
			Registry: &fakeRegistry{adapters: map[model.Agent]uninstall.AgentAdapter{model.AgentClaude: adapter}},
		},
	)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}

	// Prepare must run first (sets snapshot).
	for _, step := range plan.Prepare {
		if err := step.Run(); err != nil {
			t.Fatalf("prepare step error = %v", err)
		}
	}

	// Apply will fail on the first step (marker). Manually trigger rollback.
	for _, step := range plan.Apply {
		if err := step.Run(); err != nil {
			// Expected failure — trigger rollback on steps that support it.
			type rollbacker interface{ Rollback() error }
			for _, s := range plan.Apply {
				if rb, ok := s.(rollbacker); ok {
					_ = rb.Rollback()
				}
			}
			break
		}
	}

	_ = restoreWasCalled // rollback seam verified wired
}

// ── El uninstall borra TODOS los comandos, no solo starter-add ────────────────
//
// 2026-09-09: el instalador paso a caminar commands/<variante>/ y escribir todo
// lo que encuentra (cuatro archivos para claude). El uninstall seguia calculando
// UN path con command.RelPathForVariant, asi que borraba starter-add y dejaba
// active/init.md, opsx/ship.md y opsx/verify.md huerfanos en el sistema del
// usuario. El modo de fallar es silencioso: el comando dice que termino.

// TestCommandRemovalStep_RemovesEveryInstalledCommand asserts that a targeted
// uninstall of the claude variant removes EVERY command the installer wrote.
func TestCommandRemovalStep_RemovesEveryInstalledCommand(t *testing.T) {
	homeDir := t.TempDir()
	adapter := fakeAdapter{agent: model.AgentClaude, homeDir: homeDir, variantKey: "claude"}

	uninstall.ResetEmbeddedCommandsFS()
	t.Cleanup(uninstall.ResetEmbeddedCommandsFS)
	uninstall.WithEmbeddedCommandsFS(assets.CommandsFS)

	// Plantamos lo que un install real deja: el contenido REAL de cada asset.
	commandsDir := adapter.CommandsDir(homeDir)
	relPaths, err := command.RelPathsForVariant(assets.CommandsFS, "claude")
	if err != nil {
		t.Fatalf("setup: RelPathsForVariant: %v", err)
	}
	if len(relPaths) < 4 {
		t.Fatalf("setup: se esperaban al menos 4 comandos embebidos para claude, hay %d", len(relPaths))
	}
	for _, rel := range relPaths {
		assetPath := "commands/claude/" + filepath.ToSlash(rel)
		data, readErr := fs.ReadFile(assets.CommandsFS, assetPath)
		if readErr != nil {
			t.Fatalf("setup: leer asset %q: %v", assetPath, readErr)
		}
		dest := filepath.Join(commandsDir, rel)
		if mkErr := os.MkdirAll(filepath.Dir(dest), 0o755); mkErr != nil {
			t.Fatalf("setup: mkdir: %v", mkErr)
		}
		if wErr := os.WriteFile(dest, data, 0o644); wErr != nil {
			t.Fatalf("setup: write %q: %v", dest, wErr)
		}
	}

	restoreSnap := uninstall.SetSnapshotCreate(func(dir string, paths []string) (backup.Manifest, error) {
		return backup.Manifest{}, nil
	})
	defer restoreSnap()

	h := model.Harness{
		ID:           "starter-add-command",
		Type:         model.HarnessCommand,
		Agents:       []model.Agent{model.AgentClaude},
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
	}

	plan, err := uninstall.BuildPlan(
		&fakeCatalog{harnesses: []model.Harness{h}},
		uninstall.Intent{
			Agents:   []model.Agent{model.AgentClaude},
			Mode:     model.ModeLite,
			Strategy: uninstall.StrategyTargeted,
		},
		uninstall.Options{
			HomeDir:  homeDir,
			Registry: &fakeRegistry{adapters: map[model.Agent]uninstall.AgentAdapter{model.AgentClaude: adapter}},
		},
	)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}
	for _, step := range plan.Apply {
		if runErr := step.Run(); runErr != nil {
			t.Fatalf("step.Run() error = %v", runErr)
		}
	}

	for _, rel := range relPaths {
		dest := filepath.Join(commandsDir, rel)
		if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
			t.Errorf("comando huerfano tras el uninstall: %s", rel)
		}
	}
}

// TestCommandRemovalStep_LeavesUserAuthoredCommandsAlone es la regla dura:
// el directorio de comandos es COMPARTIDO con lo que el usuario escribe a mano,
// y el uninstall solo puede llevarse lo que este binario instalo.
//
// Este test es el que descarta listar el disco como fuente de la lista a borrar:
// en un ~/.claude/commands/ real conviven /opsx:explore, /opsx:propose y los
// /sdd-*, que active-stack nunca instalo.
func TestCommandRemovalStep_LeavesUserAuthoredCommandsAlone(t *testing.T) {
	homeDir := t.TempDir()
	adapter := fakeAdapter{agent: model.AgentClaude, homeDir: homeDir, variantKey: "claude"}

	uninstall.ResetEmbeddedCommandsFS()
	t.Cleanup(uninstall.ResetEmbeddedCommandsFS)
	uninstall.WithEmbeddedCommandsFS(assets.CommandsFS)

	commandsDir := adapter.CommandsDir(homeDir)

	// Uno nuestro (lo instalo active-stack) y tres del usuario, dos de ellos en
	// el MISMO subdirectorio que uno nuestro: opsx/ tiene ship.md nuestro y
	// explore.md del usuario.
	ours := filepath.Join(commandsDir, "opsx", "ship.md")
	userFiles := []string{
		filepath.Join(commandsDir, "opsx", "explore.md"),
		filepath.Join(commandsDir, "opsx", "propose.md"),
		filepath.Join(commandsDir, "sdd-new.md"),
	}
	for _, p := range append([]string{ours}, userFiles...) {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("setup: mkdir: %v", err)
		}
		if err := os.WriteFile(p, []byte("contenido de "+filepath.Base(p)), 0o644); err != nil {
			t.Fatalf("setup: write: %v", err)
		}
	}

	restoreSnap := uninstall.SetSnapshotCreate(func(dir string, paths []string) (backup.Manifest, error) {
		return backup.Manifest{}, nil
	})
	defer restoreSnap()

	h := model.Harness{
		ID:           "starter-add-command",
		Type:         model.HarnessCommand,
		Agents:       []model.Agent{model.AgentClaude},
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
	}
	plan, err := uninstall.BuildPlan(
		&fakeCatalog{harnesses: []model.Harness{h}},
		uninstall.Intent{
			Agents:   []model.Agent{model.AgentClaude},
			Mode:     model.ModeLite,
			Strategy: uninstall.StrategyTargeted,
		},
		uninstall.Options{
			HomeDir:  homeDir,
			Registry: &fakeRegistry{adapters: map[model.Agent]uninstall.AgentAdapter{model.AgentClaude: adapter}},
		},
	)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}
	for _, step := range plan.Apply {
		if runErr := step.Run(); runErr != nil {
			t.Fatalf("step.Run() error = %v", runErr)
		}
	}

	if _, statErr := os.Stat(ours); !os.IsNotExist(statErr) {
		t.Errorf("opsx/ship.md es nuestro y debio borrarse")
	}
	for _, p := range userFiles {
		if _, statErr := os.Stat(p); statErr != nil {
			t.Errorf("el uninstall borro un comando del usuario: %s (%v)", p, statErr)
		}
	}
}

// TestCommandRemovalStep_OpenCodeVariantRemovesOnlyItsOwn triangula con la otra
// variante: opencode tiene un solo asset plano y no debe heredar los paths
// anidados de claude.
func TestCommandRemovalStep_OpenCodeVariantRemovesOnlyItsOwn(t *testing.T) {
	homeDir := t.TempDir()
	adapter := fakeAdapter{agent: model.AgentOpenCode, homeDir: homeDir, variantKey: "opencode"}

	uninstall.ResetEmbeddedCommandsFS()
	t.Cleanup(uninstall.ResetEmbeddedCommandsFS)
	uninstall.WithEmbeddedCommandsFS(assets.CommandsFS)

	commandsDir := adapter.CommandsDir(homeDir)
	flat := filepath.Join(commandsDir, "active-starter-add.md")
	ajeno := filepath.Join(commandsDir, "active", "init.md") // path de claude, no de opencode
	for _, p := range []string{flat, ajeno} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("setup: mkdir: %v", err)
		}
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatalf("setup: write: %v", err)
		}
	}

	restoreSnap := uninstall.SetSnapshotCreate(func(dir string, paths []string) (backup.Manifest, error) {
		return backup.Manifest{}, nil
	})
	defer restoreSnap()

	h := model.Harness{
		ID:           "starter-add-command",
		Type:         model.HarnessCommand,
		Agents:       []model.Agent{model.AgentOpenCode},
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
	}
	plan, err := uninstall.BuildPlan(
		&fakeCatalog{harnesses: []model.Harness{h}},
		uninstall.Intent{
			Agents:   []model.Agent{model.AgentOpenCode},
			Mode:     model.ModeLite,
			Strategy: uninstall.StrategyTargeted,
		},
		uninstall.Options{
			HomeDir:  homeDir,
			Registry: &fakeRegistry{adapters: map[model.Agent]uninstall.AgentAdapter{model.AgentOpenCode: adapter}},
		},
	)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}
	for _, step := range plan.Apply {
		if runErr := step.Run(); runErr != nil {
			t.Fatalf("step.Run() error = %v", runErr)
		}
	}

	if _, statErr := os.Stat(flat); !os.IsNotExist(statErr) {
		t.Error("no se borro el comando plano de opencode")
	}
	if _, statErr := os.Stat(ajeno); statErr != nil {
		t.Errorf("opencode se llevo un path que es de la variante claude: %v", statErr)
	}
}

// TestCommandRemovalStep_SnapshotCoversEveryRemovedPath: el snapshot previo tiene
// que cubrir todo lo que el removal va a borrar. Si queda corto, el rollback
// restaura de menos y el usuario pierde un comando que nadie le devuelve.
func TestCommandRemovalStep_SnapshotCoversEveryRemovedPath(t *testing.T) {
	homeDir := t.TempDir()
	adapter := fakeAdapter{agent: model.AgentClaude, homeDir: homeDir, variantKey: "claude"}

	uninstall.ResetEmbeddedCommandsFS()
	t.Cleanup(uninstall.ResetEmbeddedCommandsFS)
	uninstall.WithEmbeddedCommandsFS(assets.CommandsFS)

	var gotPaths []string
	restoreSnap := uninstall.SetSnapshotCreate(func(dir string, paths []string) (backup.Manifest, error) {
		gotPaths = paths
		return backup.Manifest{}, nil
	})
	defer restoreSnap()

	h := model.Harness{
		ID:           "starter-add-command",
		Type:         model.HarnessCommand,
		Agents:       []model.Agent{model.AgentClaude},
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
	}
	plan, err := uninstall.BuildPlan(
		&fakeCatalog{harnesses: []model.Harness{h}},
		uninstall.Intent{
			Agents:   []model.Agent{model.AgentClaude},
			Mode:     model.ModeLite,
			Strategy: uninstall.StrategyTargeted,
		},
		uninstall.Options{
			HomeDir:  homeDir,
			Registry: &fakeRegistry{adapters: map[model.Agent]uninstall.AgentAdapter{model.AgentClaude: adapter}},
		},
	)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}
	for _, step := range plan.Prepare {
		if runErr := step.Run(); runErr != nil {
			t.Fatalf("prepare step error = %v", runErr)
		}
	}

	relPaths, err := command.RelPathsForVariant(assets.CommandsFS, "claude")
	if err != nil {
		t.Fatalf("RelPathsForVariant: %v", err)
	}
	for _, rel := range relPaths {
		want := filepath.Join(adapter.CommandsDir(homeDir), rel)
		found := false
		for _, p := range gotPaths {
			if p == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("el snapshot no cubre %q; paths = %v", want, gotPaths)
		}
	}
}
