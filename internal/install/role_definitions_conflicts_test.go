package install_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/pipeline"
)

// runRoleDefinitionsPlan corre el plan REAL de role-definitions (BuildPlan +
// Orchestrator + roledefs.Install, sin fakes) sobre home y devuelve el
// resultado y los eventos de progreso que vio el canal del pipeline.
func runRoleDefinitionsPlan(t *testing.T, home string) (pipeline.ExecutionResult, []pipeline.ProgressEvent) {
	t.Helper()
	cat := &fakeCatalog{harnesses: []model.Harness{roleDefinitionsHarness()}}
	reg := &fakeRegistry{adapters: map[model.Agent]install.AgentAdapter{
		model.AgentClaude: agentsDirFakeAdapter{fakeAdapter{agent: model.AgentClaude}},
	}}
	intent := install.Intent{Agents: []model.Agent{model.AgentClaude}, Mode: model.ModeCustom, Custom: []string{"role-definitions"}}

	var events []pipeline.ProgressEvent
	opts := buildOptionsWithProgress(home, reg, nil, func(ev pipeline.ProgressEvent) { events = append(events, ev) })
	plan, err := install.BuildPlan(cat, intent, opts)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}
	orch := pipeline.NewOrchestrator(pipeline.DefaultRollbackPolicy(), pipeline.WithProgressFunc(plan.OnProgress))
	return orch.Execute(plan.StagePlan), events
}

func degradedFor(events []pipeline.ProgressEvent, stepID string) []pipeline.ProgressEvent {
	var out []pipeline.ProgressEvent
	for _, ev := range events {
		if ev.StepID == stepID && ev.Status == pipeline.StepStatusDegraded {
			out = append(out, ev)
		}
	}
	return out
}

// TestRoleDefinitions_ConflictsAreReported_WorkflowRoleIsBlocking es 5.1: el
// Result de roledefs.Install ya no se descarta. Un archivo personal en un rol
// de FLUJO (role-auditor.md) se informa como BLOQUEANTE del flujo, con la
// ruta del archivo preservado y la del lateral donde quedó la versión del
// stack; uno en un ESPECIALISTA se informa también, sin la marca de
// bloqueante. La instalación del resto NO se rompe: el paso devuelve nil y
// el pipeline no revierte nada.
func TestRoleDefinitions_ConflictsAreReported_WorkflowRoleIsBlocking(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	personalAuditor := "---\nname: role-auditor\n---\n\nmi auditor de antes\n"
	personalDBA := "---\nname: role-dba\n---\n\nmi DBA\n"
	auditorPath := filepath.Join(agentsDir, "role-auditor.md")
	dbaPath := filepath.Join(agentsDir, "role-dba.md")
	if err := os.WriteFile(auditorPath, []byte(personalAuditor), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dbaPath, []byte(personalDBA), 0o644); err != nil {
		t.Fatal(err)
	}

	result, events := runRoleDefinitionsPlan(t, home)
	if result.Err != nil {
		t.Fatalf("a role conflict must not break the rest of the install: %v", result.Err)
	}

	deg := degradedFor(events, "role-definitions:role-definitions")
	if len(deg) != 1 || deg[0].Err == nil {
		t.Fatalf("want exactly one degraded event carrying the conflicts, got %+v (all events: %+v)", deg, events)
	}
	msg := deg[0].Err.Error()
	for _, want := range []string{
		"BLOQUEANTE",
		auditorPath,
		auditorPath + ".active-stack.md",
		dbaPath + ".active-stack.md",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("conflict report missing %q:\n%s", want, msg)
		}
	}
	// La marca de bloqueante es para el rol de flujo, no para el especialista.
	for _, line := range strings.Split(msg, "\n") {
		if strings.Contains(line, "role-dba.md") && strings.Contains(line, "BLOQUEANTE") {
			t.Errorf("a specialist conflict must not be marked blocking: %q", line)
		}
	}

	for path, want := range map[string]string{auditorPath: personalAuditor, dbaPath: personalDBA} {
		got, err := os.ReadFile(path)
		if err != nil || string(got) != want {
			t.Errorf("personal file %s was touched: %q err=%v", path, got, err)
		}
	}
}

// TestRoleDefinitions_NoConflicts_NoDegradedEvent triangula: una instalación
// limpia no emite ningún aviso (el canal no se llena de ruido).
func TestRoleDefinitions_NoConflicts_NoDegradedEvent(t *testing.T) {
	home := t.TempDir()
	result, events := runRoleDefinitionsPlan(t, home)
	if result.Err != nil {
		t.Fatalf("Execute: %v", result.Err)
	}
	if deg := degradedFor(events, "role-definitions:role-definitions"); len(deg) != 0 {
		t.Fatalf("clean install emitted degraded events: %+v", deg)
	}
	// Reinstalar encima (todo noop) tampoco.
	result, events = runRoleDefinitionsPlan(t, home)
	if result.Err != nil {
		t.Fatalf("Execute (2nd): %v", result.Err)
	}
	if deg := degradedFor(events, "role-definitions:role-definitions"); len(deg) != 0 {
		t.Fatalf("idempotent reinstall emitted degraded events: %+v", deg)
	}
}

// TestRoleDefinitions_OnlySpecialistConflict_IsReportedNotBlocking: sin
// conflicto en roles de flujo, el aviso existe pero no dice BLOQUEANTE.
func TestRoleDefinitions_OnlySpecialistConflict_IsReportedNotBlocking(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(agentsDir, "role-qa.md"), []byte("mi QA\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	result, events := runRoleDefinitionsPlan(t, home)
	if result.Err != nil {
		t.Fatalf("Execute: %v", result.Err)
	}
	deg := degradedFor(events, "role-definitions:role-definitions")
	if len(deg) != 1 || deg[0].Err == nil {
		t.Fatalf("want one degraded event for the specialist conflict, got %+v", deg)
	}
	if msg := deg[0].Err.Error(); strings.Contains(msg, "BLOQUEANTE") || !strings.Contains(msg, "role-qa.md") {
		t.Fatalf("specialist-only conflict report wrong:\n%s", msg)
	}
}
