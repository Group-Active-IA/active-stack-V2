package install_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/pipeline"
)

// agentsDirFakeAdapter extends fakeAdapter with AgentsDir, exercising the
// real capability the "role-definitions" harness (propose fase 2 paso 4)
// depends on — fakeAdapter itself deliberately does NOT implement it (see
// install_test.go), covering the "capability not supported, skip" branch
// elsewhere; this one covers the opposite branch end-to-end.
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

// TestBuildPlan_RoleDefinitions_DispatchesToRealInstaller is RED first via
// mutant (see the apply report): disabling the `h.ID == "role-definitions"`
// branch in buildHarnessStep makes this fail, because without it the
// generic configStep would try to inject role-*.md into CLAUDE.md instead
// of writing them under AgentsDir.
//
// End-to-end: real catalog entry, real BuildPlan, real pipeline.Orchestrator,
// real roledefs.Install (no fakes injected) — confirms the wiring through
// install.go/steps.go actually reaches the filesystem, not just that it
// compiles.
func TestBuildPlan_RoleDefinitions_DispatchesToRealInstaller(t *testing.T) {
	cat := &fakeCatalog{harnesses: []model.Harness{roleDefinitionsHarness()}}
	home := t.TempDir()
	reg := &fakeRegistry{adapters: map[model.Agent]install.AgentAdapter{
		model.AgentClaude: agentsDirFakeAdapter{fakeAdapter{agent: model.AgentClaude}},
	}}

	intent := install.Intent{
		Agents: []model.Agent{model.AgentClaude},
		Mode:   model.ModeCustom,
		Custom: []string{"role-definitions"},
	}

	plan, err := install.BuildPlan(cat, intent, buildOptions(home, reg, nil))
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}

	orch := pipeline.NewOrchestrator(pipeline.DefaultRollbackPolicy())
	result := orch.Execute(plan.StagePlan)
	if result.Err != nil {
		t.Fatalf("Execute() error = %v", result.Err)
	}

	for _, id := range []string{"role-dba", "role-backend", "role-frontend", "role-devops", "role-qa", "role-architect", "role-analyst"} {
		path := filepath.Join(home, ".claude", "agents", id+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: not written by the real pipeline: %v", id, err)
			continue
		}
		if len(data) == 0 {
			t.Errorf("%s: written empty", id)
		}
	}

	// The step ID the dedicated dispatch produces is "role-definitions:
	// role-definitions" — confirms it did NOT fall into the generic
	// configStep (which would produce "config:role-definitions" instead).
	// applyStepIDs only strips a fixed list of known prefixes, and
	// "role-definitions:" is not one of them, so the full ID is asserted
	// directly rather than widening that shared helper for one caller.
	var rawIDs []string
	for _, s := range plan.StagePlan.Apply {
		rawIDs = append(rawIDs, s.ID())
	}
	if !containsID(rawIDs, "role-definitions:role-definitions") {
		t.Errorf("expected step id %q in the plan, got %v", "role-definitions:role-definitions", rawIDs)
	}
}
