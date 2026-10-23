package roledefs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/roles"
)

// TestHealth_AllInstalled_NoIssues is RED first: Health doesn't exist yet.
func TestHealth_AllInstalled_NoIssues(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}
	res := mustResources(t)

	if _, err := Install(adapters, home, backupDir, res); err != nil {
		t.Fatalf("Install: %v", err)
	}

	issues, err := Health(adapters, home, res)
	if err != nil {
		t.Fatalf("Health: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("got %d issues on a fully installed set, want 0: %+v", len(issues), issues)
	}
}

// TestHealth_ReportsMissingAndEmpty_NotJustOneOfThem triangula los dos modos
// de falla que el propose distingue del CLI: definición FALTANTE (el archivo
// no existe) y definición VACÍA (existe, size 0) — no son lo mismo y el
// reporte tiene que distinguirlos.
func TestHealth_ReportsMissingAndEmpty_NotJustOneOfThem(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}
	res := mustResources(t)

	if _, err := Install(adapters, home, backupDir, res); err != nil {
		t.Fatalf("Install: %v", err)
	}

	// role-dba: falta (la borramos a mano, por fuera de Uninstall).
	if err := os.Remove(filepath.Join(agentsDir, roles.RoleDBA+".md")); err != nil {
		t.Fatal(err)
	}
	// role-qa: vacía.
	if err := os.WriteFile(filepath.Join(agentsDir, roles.RoleQA+".md"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	issues, err := Health(adapters, home, res)
	if err != nil {
		t.Fatalf("Health: %v", err)
	}
	byResource := map[string]string{}
	for _, iss := range issues {
		byResource[iss.Resource] = iss.Problem
	}
	if byResource[roles.RoleDBA+".md"] != "missing" {
		t.Errorf("role-dba.md problem = %q, want %q (issues: %+v)", byResource[roles.RoleDBA+".md"], "missing", issues)
	}
	if byResource[roles.RoleQA+".md"] != "empty" {
		t.Errorf("role-qa.md problem = %q, want %q (issues: %+v)", byResource[roles.RoleQA+".md"], "empty", issues)
	}
	if len(issues) != 2 {
		t.Errorf("got %d issues, want exactly 2 (the other 5 roles are healthy): %+v", len(issues), issues)
	}
}

// TestHealth_SkipsUnsupportedAdapters triangula que un adapter sin capability
// no genera issues (no es "falta instalar", es "no aplica").
func TestHealth_SkipsUnsupportedAdapters(t *testing.T) {
	home := t.TempDir()
	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentCodex, dir: ""}}
	issues, err := Health(adapters, home, mustResources(t))
	if err != nil {
		t.Fatalf("Health: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("got %d issues for an unsupported adapter, want 0: %+v", len(issues), issues)
	}
}
