package roledefs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/roles"
)

// TestUninstall_RemovesOwnIntactFiles is RED first: Uninstall doesn't exist
// yet. Lo que instalamos y nadie tocó se borra.
func TestUninstall_RemovesOwnIntactFiles(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}

	res := mustResources(t)
	if _, err := Install(adapters, home, backupDir, res); err != nil {
		t.Fatalf("Install: %v", err)
	}

	if _, err := Uninstall(adapters, home, res); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	for _, id := range canonicalRoleIDs {
		path := filepath.Join(agentsDir, id+".md")
		if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
			t.Errorf("%s: still present after Uninstall (stat err: %v)", id, statErr)
		}
	}
}

// TestUninstall_NeverRemovesTheAgentsDirItself triangula: el directorio en sí
// (compartido con agentes personales) nunca se borra, aunque quede "vacío"
// desde nuestra perspectiva.
func TestUninstall_NeverRemovesTheAgentsDirItself(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}

	res := mustResources(t)
	if _, err := Install(adapters, home, backupDir, res); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if _, err := Uninstall(adapters, home, res); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	info, err := os.Stat(agentsDir)
	if err != nil || !info.IsDir() {
		t.Fatalf("agents dir itself was removed or is gone: err=%v", err)
	}
}

// TestUninstall_PreservesPersonalAndUserModifiedFiles es el escenario
// central: un archivo personal preexistente, y uno que instalamos pero el
// usuario editó después, sobreviven a Uninstall sin avisos silenciosos de
// pérdida — "nunca borrar por glob ni sólo por nombre/manifiesto".
func TestUninstall_PreservesPersonalAndUserModifiedFiles(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}

	res := mustResources(t)
	if _, err := Install(adapters, home, backupDir, res); err != nil {
		t.Fatalf("Install: %v", err)
	}

	// role-backend: el usuario lo editó después de instalarlo.
	backendPath := filepath.Join(agentsDir, roles.RoleBackend+".md")
	edited := "editado a mano después de instalar\n"
	if err := os.WriteFile(backendPath, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}

	// role-frontend: nunca lo tocamos nosotros — es un archivo personal con
	// ese mismo nombre, puesto ahí DESPUÉS del install (simula el caso real:
	// manifiesto conoce role-frontend con un hash viejo que ya no corresponde
	// porque Install lo escribió, pero el usuario lo reemplazó entero).
	frontendPath := filepath.Join(agentsDir, roles.RoleFrontend+".md")
	personal := "esto es mío, no de active-stack\n"
	if err := os.WriteFile(frontendPath, []byte(personal), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := Uninstall(adapters, home, res)
	if err != nil {
		t.Fatalf("Uninstall: %v", err)
	}

	gotBackend, err := os.ReadFile(backendPath)
	if err != nil {
		t.Fatalf("role-backend.md disappeared: %v", err)
	}
	if string(gotBackend) != edited {
		t.Fatalf("role-backend.md was not preserved: got %q", gotBackend)
	}
	gotFrontend, err := os.ReadFile(frontendPath)
	if err != nil {
		t.Fatalf("role-frontend.md disappeared: %v", err)
	}
	if string(gotFrontend) != personal {
		t.Fatalf("role-frontend.md was not preserved: got %q", gotFrontend)
	}

	// Las demás (dba, devops, qa, architect, analyst) sí se borraron: no es
	// que Uninstall no borre nada, es que no borra lo que no es suyo.
	stillInstalled := 0
	for _, id := range canonicalRoleIDs {
		if id == roles.RoleBackend || id == roles.RoleFrontend {
			continue
		}
		if _, statErr := os.Stat(filepath.Join(agentsDir, id+".md")); statErr == nil {
			stillInstalled++
		}
	}
	if stillInstalled != 0 {
		t.Errorf("%d of the untouched roles were not removed by Uninstall", stillInstalled)
	}

	preserved := result.Resources
	foundBackend, foundFrontend := false, false
	for _, rr := range preserved {
		if rr.Resource == roles.RoleBackend+".md" {
			foundBackend = true
			if rr.Outcome != OutcomePreserved {
				t.Errorf("role-backend.md outcome = %q, want %q", rr.Outcome, OutcomePreserved)
			}
		}
		if rr.Resource == roles.RoleFrontend+".md" {
			foundFrontend = true
			if rr.Outcome != OutcomePreserved {
				t.Errorf("role-frontend.md outcome = %q, want %q", rr.Outcome, OutcomePreserved)
			}
		}
	}
	if !foundBackend || !foundFrontend {
		t.Fatalf("Uninstall result missing the preserved entries: %+v", preserved)
	}
}

// TestUninstall_IsIdempotent corre dos veces: la segunda no debe fallar ni
// reportar nada raro sobre archivos que ya no están.
func TestUninstall_IsIdempotent(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}
	res := mustResources(t)

	if _, err := Install(adapters, home, backupDir, res); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if _, err := Uninstall(adapters, home, res); err != nil {
		t.Fatalf("first Uninstall: %v", err)
	}
	if _, err := Uninstall(adapters, home, res); err != nil {
		t.Fatalf("second Uninstall (idempotency): %v", err)
	}
}
