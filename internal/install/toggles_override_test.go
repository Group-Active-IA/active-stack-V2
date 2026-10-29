package install_test

import (
	"reflect"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/backup"
	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

func orchestratorToggles(t *testing.T, override []string) []string {
	t.Helper()
	restoreSnap := install.SetSnapshotCreate(func(string, []string) (backup.Manifest, error) { return backup.Manifest{}, nil })
	t.Cleanup(restoreSnap)
	var seen []string
	restoreConfig := install.SetConfigInstallFn(func(h model.Harness, _ interface{}, _ string) error {
		if h.ID == "sdd-orchestrator" {
			seen = h.Toggles
		}
		return nil
	})
	t.Cleanup(restoreConfig)
	h := model.Harness{
		ID: "sdd-orchestrator", Type: model.HarnessConfig,
		Toggles:      []string{"lean-workflow", "governance", "ship"},
		InstallModes: []model.InstallMode{model.ModeFull}, Agents: []model.Agent{model.AgentClaude},
	}
	cat := &fakeCatalog{harnesses: []model.Harness{h}}
	reg := &fakeRegistry{adapters: map[model.Agent]install.AgentAdapter{model.AgentClaude: fakeAdapter{agent: model.AgentClaude}}}
	plan, err := install.BuildPlan(cat, install.Intent{Agents: []model.Agent{model.AgentClaude}, Mode: model.ModeFull, Toggles: override},
		install.Options{HomeDir: t.TempDir(), Registry: reg, NoSelfInstall: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range plan.Apply {
		if err := s.Run(); err != nil {
			t.Fatalf("%s: %v", s.ID(), err)
		}
	}
	return seen
}

func TestBuildPlan_TogglesOverrideTheOrchestrator(t *testing.T) {
	got := orchestratorToggles(t, []string{"lean-workflow", "evidence"})
	if !reflect.DeepEqual(got, []string{"lean-workflow", "evidence"}) {
		t.Fatalf("orchestrator composed with %v, want the --toggles list", got)
	}
}

func TestBuildPlan_NoOverrideKeepsCatalogToggles(t *testing.T) {
	got := orchestratorToggles(t, nil)
	if !reflect.DeepEqual(got, []string{"lean-workflow", "governance", "ship"}) {
		t.Fatalf("orchestrator composed with %v, want the catalog list", got)
	}
}

// The post-install check compares the installed block with Compose(h.Toggles):
// SelectHarnesses feeds it, so the override must reach it too, or a --toggles
// install fails verification and rolls back (found installing with the real
// binary, 2026-10-29).
func TestSelectHarnesses_CarriesTogglesOverride(t *testing.T) {
	h := model.Harness{
		ID: "sdd-orchestrator", Type: model.HarnessConfig,
		Toggles:      []string{"lean-workflow", "ship"},
		InstallModes: []model.InstallMode{model.ModeFull}, Agents: []model.Agent{model.AgentClaude},
	}
	got, err := install.SelectHarnesses(&fakeCatalog{harnesses: []model.Harness{h}},
		install.Intent{Agents: []model.Agent{model.AgentClaude}, Mode: model.ModeFull, Toggles: []string{"lean-workflow"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range got {
		if s.ID == "sdd-orchestrator" && !reflect.DeepEqual(s.Toggles, []string{"lean-workflow"}) {
			t.Fatalf("verification would compare against %v, want the --toggles list", s.Toggles)
		}
	}
}
