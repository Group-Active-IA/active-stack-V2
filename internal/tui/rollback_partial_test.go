package tui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/pipeline"
)

// TestViewComplete_PartialRollbackIsNotShownAsSucceeded covers the TUI half of
// the honest-rollback fix: a rollback whose every step returned nil but that
// left things on the machine must not render as "succeeded".
func TestViewComplete_PartialRollbackIsNotShownAsSucceeded(t *testing.T) {
	m := newModel(ModelDeps{})
	m.Screen = ScreenComplete
	m.ExecutionResult = pipeline.ExecutionResult{
		Err: fmt.Errorf("sudo: npm: command not found"),
		Rollback: pipeline.StageResult{
			Stage:      pipeline.StageRollback,
			Success:    true,
			Steps:      []pipeline.StepResult{{StepID: "external:engram", Status: pipeline.StepStatusRolledBack}},
			Unreverted: []string{`npm package "openspec" stays installed globally`},
		},
	}

	view := m.viewComplete()

	if strings.Contains(view, "succeeded") {
		t.Errorf("a rollback that left things behind must not render as 'succeeded'; got:\n%s", view)
	}
	if !strings.Contains(view, "partial") {
		t.Errorf("view must report the rollback as partial; got:\n%s", view)
	}
	if !strings.Contains(view, "openspec") {
		t.Errorf("view must name what was left in place; got:\n%s", view)
	}
}
