package tui

// C-32: TUI honest reporting — tasks 4.1, 4.2, 4.3

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/pipeline"
)

// Task 4.1 RED — statusIcon for StepStatusDegraded returns a distinct glyph
// (warning style, not the hard-fail ✗ and not the success ✓).
func TestStatusIcon_Degraded_IsDistinct(t *testing.T) {
	degradedIcon := statusIcon(pipeline.StepStatusDegraded)

	// Must not be empty.
	if degradedIcon == "" {
		t.Fatal("statusIcon(StepStatusDegraded) must return a non-empty glyph")
	}
	// Must differ from the hard-failure glyph.
	failedIcon := statusIcon(pipeline.StepStatusFailed)
	if degradedIcon == failedIcon {
		t.Errorf("StepStatusDegraded icon must differ from StepStatusFailed icon (both = %q)", degradedIcon)
	}
	// Must differ from the success glyph.
	successIcon := statusIcon(pipeline.StepStatusSucceeded)
	if degradedIcon == successIcon {
		t.Errorf("StepStatusDegraded icon must differ from StepStatusSucceeded icon (both = %q)", degradedIcon)
	}
	// Must contain ⚠ (the design-specified warning glyph).
	if !strings.Contains(degradedIcon, "⚠") {
		t.Errorf("StepStatusDegraded icon must contain ⚠; got %q", degradedIcon)
	}
}

// Task 4.3 RED+GREEN — viewComplete includes degraded harnesses in the summary
// while reporting the run as completed/success (not failed).
func TestViewComplete_DegradedHarnessesInSummary(t *testing.T) {
	m := newModel(ModelDeps{})
	m.Screen = ScreenComplete
	// Simulate a successful run (no error) with one degraded step row.
	m.ExecutionResult = pipeline.ExecutionResult{Err: nil}
	m.stepRows = []stepRow{
		{stepID: "skill:code-review-excellence", status: pipeline.StepStatusDegraded, err: fmt.Errorf("upstream path missing")},
		{stepID: "skill:agile-product-owner", status: pipeline.StepStatusSucceeded},
	}

	view := m.viewComplete()

	// The view must report installation as completed (success path).
	if strings.Contains(view, "Installation failed") {
		t.Errorf("viewComplete must not show 'Installation failed' for a degraded run; got:\n%s", view)
	}
	if !strings.Contains(view, "Installation complete") {
		t.Errorf("viewComplete must show 'Installation complete' for a degraded run; got:\n%s", view)
	}

	// The view must list the degraded harness.
	if !strings.Contains(view, "code-review-excellence") {
		t.Errorf("viewComplete must list degraded harness 'code-review-excellence'; got:\n%s", view)
	}

	// The clean harness must NOT appear in a dedicated "degraded" section.
	// (It may appear in the step list, but not in the degraded summary.)
	// Verify the summary has some "degraded" indicator.
	if !strings.Contains(strings.ToLower(view), "degraded") {
		t.Errorf("viewComplete must mention 'degraded' for runs with degraded harnesses; got:\n%s", view)
	}
}

// Triangulate: clean run (no degraded rows) → viewComplete does NOT show a
// degraded section.
func TestViewComplete_CleanRun_NoDegradedSection(t *testing.T) {
	m := newModel(ModelDeps{})
	m.Screen = ScreenComplete
	m.ExecutionResult = pipeline.ExecutionResult{Err: nil}
	m.stepRows = []stepRow{
		{stepID: "skill:engram", status: pipeline.StepStatusSucceeded},
		{stepID: "skill:openspec", status: pipeline.StepStatusSucceeded},
	}

	view := m.viewComplete()

	if strings.Contains(strings.ToLower(view), "degraded") {
		t.Errorf("viewComplete for a clean run must NOT show 'degraded'; got:\n%s", view)
	}
}

// TestApplyProgressEvent_SucceededDoesNotEraseDegraded (lean-agent-workflow
// 5.1): un paso que avisa por el canal degraded y después devuelve nil
// recibe del runner un evento Succeeded inmediatamente después. Si ese
// evento pisa la fila, la pantalla final dice "All steps succeeded." y el
// aviso — p. ej. un conflicto BLOQUEANTE en role-auditor.md — desaparece.
func TestApplyProgressEvent_SucceededDoesNotEraseDegraded(t *testing.T) {
	m := newModel(ModelDeps{})
	m.stepRows = []stepRow{{stepID: "role-definitions:role-definitions", status: pipeline.StepStatusRunning}}

	warning := fmt.Errorf("BLOQUEANTE del flujo lean: role-auditor.md")
	m.applyProgressEvent(pipeline.ProgressEvent{StepID: "role-definitions:role-definitions", Status: pipeline.StepStatusDegraded, Err: warning})
	m.applyProgressEvent(pipeline.ProgressEvent{StepID: "role-definitions:role-definitions", Status: pipeline.StepStatusSucceeded})

	row := m.stepRows[0]
	if row.status != pipeline.StepStatusDegraded || row.err == nil || !strings.Contains(row.err.Error(), "BLOQUEANTE") {
		t.Fatalf("row after degraded+succeeded = %+v, want degraded with the warning kept", row)
	}

	m.Screen = ScreenComplete
	m.ExecutionResult = pipeline.ExecutionResult{}
	view := m.viewComplete()
	if !strings.Contains(view, "BLOQUEANTE") || strings.Contains(view, "All steps succeeded.") {
		t.Fatalf("completion screen hides the warning:\n%s", view)
	}
}

// TestApplyProgressEvent_SkippedNoteIsKeptAndShown (lean-agent-workflow 5.3):
// un aviso informativo (evento skipped con mensaje, p. ej. "MCP engram no
// registrado: el plugin engram@engram lo provee") no cambia el estado del
// paso — sigue siendo éxito — pero aparece en la pantalla final.
func TestApplyProgressEvent_SkippedNoteIsKeptAndShown(t *testing.T) {
	m := newModel(ModelDeps{})
	m.stepRows = []stepRow{{stepID: "external:engram", status: pipeline.StepStatusRunning}}
	m.applyProgressEvent(pipeline.ProgressEvent{StepID: "external:engram", Status: pipeline.StepStatusSkipped, Err: fmt.Errorf("plugin engram@engram lo provee")})
	m.applyProgressEvent(pipeline.ProgressEvent{StepID: "external:engram", Status: pipeline.StepStatusSucceeded})

	if m.stepRows[0].status != pipeline.StepStatusSucceeded {
		t.Fatalf("a note must not change the step status, got %q", m.stepRows[0].status)
	}
	m.Screen = ScreenComplete
	m.ExecutionResult = pipeline.ExecutionResult{}
	view := m.viewComplete()
	if !strings.Contains(view, "engram@engram") || !strings.Contains(view, "All steps succeeded.") {
		t.Fatalf("completion screen must keep success and show the note:\n%s", view)
	}
}
