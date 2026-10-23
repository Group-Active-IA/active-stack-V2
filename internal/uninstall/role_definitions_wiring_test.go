package uninstall_test

import (
	"os"
	"path/filepath"
	"testing"

	roledefsinstaller "github.com/Group-Active-IA/active-stack/internal/harness/roledefs"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/pipeline"
	"github.com/Group-Active-IA/active-stack/internal/uninstall"
)

// agentsDirFakeAdapter extends fakeAdapter with AgentsDir, mirroring
// internal/install's own role_definitions_wiring_test.go.
type agentsDirFakeAdapter struct {
	fakeAdapter
}

func (a agentsDirFakeAdapter) AgentsDir(homeDir string) string {
	return filepath.Join(homeDir, ".claude", "agents")
}

func roleDefinitionsHarness() model.Harness {
	return model.Harness{
		ID:           "role-definitions",
		Name:         "Role Definitions",
		Type:         model.HarnessConfig,
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
		Agents:       []model.Agent{model.AgentClaude},
	}
}

// TestBuildPlan_RoleDefinitions_RemovesOwnedFilesForReal is RED first via
// mutant (documented in the apply report): removing the `h.ID ==
// "role-definitions"` branch in buildUninstallStep makes this fail, because
// the files roledefs.Install wrote would otherwise never be removed by the
// generic markerRemovalStep (which edits InstructionsPath, not AgentsDir).
//
// Setup goes straight through roledefs.Install (already covered end-to-end
// by internal/install's own role_definitions_wiring_test.go) so this test's
// only job is the uninstall side of the wiring: real uninstall.BuildPlan +
// pipeline.Orchestrator.Execute reaching the real filesystem, preserving a
// personal file with a colliding name.
func TestBuildPlan_RoleDefinitions_RemovesOwnedFilesForReal(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")

	adapter := agentsDirFakeAdapter{fakeAdapter{agent: model.AgentClaude}}
	res, err := roledefsinstaller.DefaultResources()
	if err != nil {
		t.Fatalf("DefaultResources: %v", err)
	}
	if _, err := roledefsinstaller.Install([]roledefsinstaller.AgentAdapter{adapter}, home, backupDir, res); err != nil {
		t.Fatalf("roledefs.Install (setup): %v", err)
	}

	// Personal file, same directory, never touched by active-stack.
	personalPath := filepath.Join(agentsDir, "my-own-agent.md")
	if err := os.WriteFile(personalPath, []byte("mío"), 0o644); err != nil {
		t.Fatal(err)
	}

	cat := &fakeCatalog{harnesses: []model.Harness{roleDefinitionsHarness()}}
	reg := &fakeRegistry{adapters: map[model.Agent]uninstall.AgentAdapter{
		model.AgentClaude: adapter,
	}}
	uIntent := uninstall.Intent{
		Agents:   []model.Agent{model.AgentClaude},
		Mode:     model.ModeCustom,
		Custom:   []string{"role-definitions"},
		Strategy: uninstall.StrategyTargeted,
	}
	uPlan, err := uninstall.BuildPlan(cat, uIntent, buildUninstallOptions(home, reg))
	if err != nil {
		t.Fatalf("uninstall.BuildPlan: %v", err)
	}
	if res := pipeline.NewOrchestrator(pipeline.DefaultRollbackPolicy()).Execute(uPlan.StagePlan); res.Err != nil {
		t.Fatalf("uninstall.Execute: %v", res.Err)
	}

	if _, statErr := os.Stat(filepath.Join(agentsDir, "role-dba.md")); !os.IsNotExist(statErr) {
		t.Errorf("role-dba.md still present after uninstall (stat err: %v)", statErr)
	}
	got, err := os.ReadFile(personalPath)
	if err != nil || string(got) != "mío" {
		t.Errorf("personal file was not preserved by uninstall: data=%q err=%v", got, err)
	}
}
