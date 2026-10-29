package catalog

import "testing"

// Selecting the catalogue's default orchestrator must activate the lean path.
func TestDefaultOrchestratorSelectsLeanWorkflow(t *testing.T) {
	cat, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	h, ok := cat.ByID("sdd-orchestrator")
	if !ok {
		t.Fatal("orchestrator missing")
	}
	found := false
	for _, toggle := range h.Toggles {
		if toggle == "lean-workflow" {
			found = true
		}
	}
	if !found {
		t.Fatal("default orchestrator does not select lean-workflow")
	}
}
