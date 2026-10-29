package catalog

import (
	"regexp"
	"testing"
)

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

// The npm installer only compares versions when the catalog declares one;
// without it any openspec on PATH passes (review 2026-10-29, 8.2).
func TestOpenSpecDeclaresMinVersion(t *testing.T) {
	cat, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	h, ok := cat.ByID("openspec")
	if !ok || h.External == nil {
		t.Fatal("openspec harness missing")
	}
	if !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(h.External.MinVersion) {
		t.Fatalf("openspec min_version = %q, want x.y.z", h.External.MinVersion)
	}
}
