package roledefs

import (
	"path/filepath"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/model"
)

// TestLoadManifest_MissingFileReturnsEmpty is RED first: manifestPath,
// loadManifest and saveManifest don't exist yet. A never-installed agent
// must read as "nothing owned", not as an error — otherwise the very first
// Install run on a fresh machine would fail before writing anything.
func TestLoadManifest_MissingFileReturnsEmpty(t *testing.T) {
	home := t.TempDir()
	m, err := loadManifest(home, model.AgentClaude)
	if err != nil {
		t.Fatalf("loadManifest on a fresh home: %v", err)
	}
	if len(m.Resources) != 0 {
		t.Fatalf("got %d owned resources on a fresh home, want 0: %+v", len(m.Resources), m.Resources)
	}
}

// TestSaveManifest_RoundTrips triangula con un segundo agente y un segundo
// hash para confirmar que no es un acierto de un solo valor.
func TestSaveManifest_RoundTrips(t *testing.T) {
	home := t.TempDir()
	m := ownershipManifest{Resources: map[string]string{
		"role-dba.md":     "aaaa",
		"role-backend.md": "bbbb",
	}}
	if err := saveManifest(home, model.AgentClaude, m); err != nil {
		t.Fatalf("saveManifest: %v", err)
	}
	got, err := loadManifest(home, model.AgentClaude)
	if err != nil {
		t.Fatalf("loadManifest: %v", err)
	}
	if got.Resources["role-dba.md"] != "aaaa" || got.Resources["role-backend.md"] != "bbbb" {
		t.Fatalf("round-trip mismatch: %+v", got.Resources)
	}
}

// TestManifestPath_IsPerAgent asserts two agents don't share a manifest file
// — otherwise installing for codex could read/overwrite claude's ownership
// record.
func TestManifestPath_IsPerAgent(t *testing.T) {
	home := "/fake/home"
	claudePath := manifestPath(home, model.AgentClaude)
	codexPath := manifestPath(home, model.AgentCodex)
	if claudePath == codexPath {
		t.Fatalf("claude and codex share a manifest path: %q", claudePath)
	}
	if filepath.Dir(claudePath) != filepath.Join(home, ".active-stack", "state", "role-definitions") {
		t.Fatalf("manifest path %q not under the expected state directory", claudePath)
	}
}
