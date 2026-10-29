package install_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/external"
	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/pipeline"
	"github.com/Group-Active-IA/active-stack/internal/system"
)

// TestExternalStep_ReportsInstallerNotes (lean-agent-workflow 5.3): cuando
// el installer externo decide no registrar el MCP engram porque el plugin
// engram@engram ya lo provee, el salto llega al canal de progreso del
// pipeline como evento "skipped" con el mensaje — no se pierde en el Result.
func TestExternalStep_ReportsInstallerNotes(t *testing.T) {
	note := `MCP "engram" no registrado en /h/.claude.json: el plugin engram@engram ya está habilitado y lo provee`
	restore := install.SetExternalInstallFn(func(context.Context, model.Harness, system.PlatformProfile, []external.AgentAdapter, string) (external.Result, error) {
		return external.Result{BinaryPath: "/bin/engram", AlreadyInstalled: true, Notes: []string{note}}, nil
	})
	defer restore()

	h := model.Harness{
		ID: "engram", Type: model.HarnessExternal,
		External:     &model.External{Method: "homebrew", MCP: &model.MCP{Name: "engram", Command: "engram"}},
		InstallModes: []model.InstallMode{model.ModeLite},
	}
	cat := &fakeCatalog{harnesses: []model.Harness{h}}
	reg := &fakeRegistry{adapters: map[model.Agent]install.AgentAdapter{model.AgentClaude: fakeAdapter{agent: model.AgentClaude}}}
	var events []pipeline.ProgressEvent
	opts := buildOptionsWithProgress(t.TempDir(), reg, nil, func(ev pipeline.ProgressEvent) { events = append(events, ev) })
	opts.NoSelfInstall = true
	plan, err := install.BuildPlan(cat, install.Intent{Agents: []model.Agent{model.AgentClaude}, Mode: model.ModeLite}, opts)
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	result := pipeline.NewOrchestrator(pipeline.DefaultRollbackPolicy(), pipeline.WithProgressFunc(plan.OnProgress)).Execute(plan.StagePlan)
	if result.Err != nil {
		t.Fatalf("Execute: %v", result.Err)
	}

	var got []string
	for _, ev := range events {
		if ev.StepID == "external:engram" && ev.Status == pipeline.StepStatusSkipped && ev.Err != nil {
			got = append(got, ev.Err.Error())
		}
	}
	if len(got) != 1 || !strings.Contains(got[0], "engram@engram") {
		t.Fatalf("want one skipped event with the plugin note, got %q (events %+v)", got, events)
	}
}
