package roledefs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/roles"
)

type fakeAdapter struct {
	agent model.Agent
	dir   string // "" = capability unsupported
}

func (f fakeAdapter) Agent() model.Agent      { return f.agent }
func (f fakeAdapter) AgentsDir(string) string { return f.dir }

func mustResources(t *testing.T) []Resource {
	t.Helper()
	res, err := DefaultResources()
	if err != nil {
		t.Fatalf("DefaultResources: %v", err)
	}
	return res
}

func resourceByRole(t *testing.T, id string) Resource {
	t.Helper()
	for _, r := range mustResources(t) {
		if r.Name == id+".md" {
			return r
		}
	}
	t.Fatalf("no resource named %q", id+".md")
	return Resource{}
}

// TestInstall_AbsentCreates is RED first: Install doesn't exist yet.
// Un home nuevo (sin ~/.claude/agents/) tiene que terminar con los siete
// archivos creados, y el directorio se crea solo.
func TestInstall_AbsentCreates(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")

	res, err := Install([]AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}, home, backupDir, mustResources(t))
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	for _, id := range canonicalRoleIDs {
		path := filepath.Join(agentsDir, id+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: not written: %v", id, err)
			continue
		}
		if len(data) == 0 {
			t.Errorf("%s: written empty", id)
		}
	}
	created := 0
	for _, rr := range res.Resources {
		if rr.Outcome == OutcomeCreated {
			created++
		}
	}
	if created != 12 {
		t.Errorf("got %d OutcomeCreated, want 12 (full result: %+v)", created, res.Resources)
	}
}

// TestInstall_IdenticalIsNoop triangula con una segunda corrida: reinstalar
// sobre contenido idéntico no debe escribir nada ni generar un conflicto.
func TestInstall_IdenticalIsNoop(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}

	if _, err := Install(adapters, home, backupDir, mustResources(t)); err != nil {
		t.Fatalf("first install: %v", err)
	}
	before, err := os.ReadFile(filepath.Join(agentsDir, roles.RoleDBA+".md"))
	if err != nil {
		t.Fatalf("read after first install: %v", err)
	}

	res, err := Install(adapters, home, backupDir, mustResources(t))
	if err != nil {
		t.Fatalf("second install: %v", err)
	}
	after, err := os.ReadFile(filepath.Join(agentsDir, roles.RoleDBA+".md"))
	if err != nil {
		t.Fatalf("read after second install: %v", err)
	}
	if string(before) != string(after) {
		t.Fatalf("content changed on a no-op reinstall")
	}
	for _, rr := range res.Resources {
		if rr.Outcome != OutcomeNoop {
			t.Errorf("resource %q: got outcome %q on an identical reinstall, want noop", rr.Resource, rr.Outcome)
		}
	}
}

// TestInstall_OwnIntactUpdates_WithBackup cubre "propio intacto actualiza con
// backup": si lo que hay en disco es exactamente lo que ESTA harness escribió
// la vez anterior (hash grabado en el manifiesto), una nueva versión del
// contenido candidato se escribe encima, con backup.
func TestInstall_OwnIntactUpdates_WithBackup(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}

	v1 := mustResources(t)
	if _, err := Install(adapters, home, backupDir, v1); err != nil {
		t.Fatalf("install v1: %v", err)
	}

	// Simula una nueva versión del contenido canónico (v2) para role-dba.
	v2 := append([]Resource{}, v1...)
	for i, r := range v2 {
		if r.Name == roles.RoleDBA+".md" {
			v2[i].Content = append(append([]byte{}, r.Content...), []byte("\n\nv2\n")...)
		}
	}

	res, err := Install(adapters, home, backupDir, v2)
	if err != nil {
		t.Fatalf("install v2: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(agentsDir, roles.RoleDBA+".md"))
	if err != nil {
		t.Fatalf("read role-dba after v2: %v", err)
	}
	wantContent := resourceByRole(t, roles.RoleDBA)
	for _, r := range v2 {
		if r.Name == roles.RoleDBA+".md" {
			wantContent = r
		}
	}
	if string(got) != string(wantContent.Content) {
		t.Fatalf("role-dba.md was not updated to the new candidate content")
	}
	var dbaResult ResourceResult
	found := false
	for _, rr := range res.Resources {
		if rr.Resource == roles.RoleDBA+".md" {
			dbaResult = rr
			found = true
		}
	}
	if !found {
		t.Fatalf("no result for role-dba.md")
	}
	if dbaResult.Outcome != OutcomeUpdated {
		t.Fatalf("role-dba.md outcome = %q, want %q", dbaResult.Outcome, OutcomeUpdated)
	}

	// Backup real en disco, no sólo en memoria.
	entries, err := os.ReadDir(backupDir)
	if err != nil || len(entries) == 0 {
		t.Fatalf("expected at least one snapshot under %q, got entries=%v err=%v", backupDir, entries, err)
	}
}

// TestInstall_PersonalFileIsPreserved_WithConflictSideFile es el escenario
// central del paso 4: un archivo PERSONAL preexistente (nunca escrito por
// este harness, no registrado en el manifiesto de ownership) con el mismo
// nombre que un rol tiene que sobrevivir intacto, y la versión candidata
// nuestra aparece aparte, no lo pisa.
func TestInstall_PersonalFileIsPreserved_WithConflictSideFile(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	personal := "---\nname: role-dba\ndescription: mi propio rol de DBA, escrito a mano.\n---\n\nNo me toques.\n"
	dbaPath := filepath.Join(agentsDir, roles.RoleDBA+".md")
	if err := os.WriteFile(dbaPath, []byte(personal), 0o644); err != nil {
		t.Fatal(err)
	}

	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}
	res, err := Install(adapters, home, backupDir, mustResources(t))
	if err != nil {
		t.Fatalf("Install: %v", err)
	}

	stillThere, err := os.ReadFile(dbaPath)
	if err != nil {
		t.Fatalf("read personal file after install: %v", err)
	}
	if string(stillThere) != personal {
		t.Fatalf("personal role-dba.md was modified by Install:\ngot:  %q\nwant: %q", stillThere, personal)
	}

	var dbaResult ResourceResult
	found := false
	for _, rr := range res.Resources {
		if rr.Resource == roles.RoleDBA+".md" {
			dbaResult = rr
			found = true
		}
	}
	if !found {
		t.Fatalf("no result for role-dba.md")
	}
	if dbaResult.Outcome != OutcomeConflict {
		t.Fatalf("role-dba.md outcome = %q, want %q", dbaResult.Outcome, OutcomeConflict)
	}
	if dbaResult.ConflictPath == "" {
		t.Fatalf("conflict with no alternate path recorded")
	}
	ours, err := os.ReadFile(dbaResult.ConflictPath)
	if err != nil {
		t.Fatalf("read conflict side file %q: %v", dbaResult.ConflictPath, err)
	}
	if string(ours) != string(resourceByRole(t, roles.RoleDBA).Content) {
		t.Fatalf("conflict side file does not carry our candidate content")
	}
}

// TestInstall_UserModifiedOwnFile_IsPreservedNotOverwritten triangula el
// otro camino hacia OutcomeConflict: un archivo que SÍ instalamos nosotros
// antes, pero que el usuario editó a mano después, tiene que tratarse igual
// que uno personal — nunca "propio intacto".
func TestInstall_UserModifiedOwnFile_IsPreservedNotOverwritten(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}

	if _, err := Install(adapters, home, backupDir, mustResources(t)); err != nil {
		t.Fatalf("install v1: %v", err)
	}
	dbaPath := filepath.Join(agentsDir, roles.RoleDBA+".md")
	edited := "contenido editado a mano por el usuario, ya no es lo que instalamos\n"
	if err := os.WriteFile(dbaPath, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := Install(adapters, home, backupDir, mustResources(t))
	if err != nil {
		t.Fatalf("install v2: %v", err)
	}
	got, err := os.ReadFile(dbaPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != edited {
		t.Fatalf("user-edited role-dba.md was overwritten: got %q", got)
	}
	for _, rr := range res.Resources {
		if rr.Resource == roles.RoleDBA+".md" && rr.Outcome != OutcomeConflict {
			t.Fatalf("role-dba.md outcome = %q, want %q", rr.Outcome, OutcomeConflict)
		}
	}
}

// TestInstall_SkipsAdapterWithEmptyAgentsDir cubre "capability vacío
// significa no soportado" — el adapter se saltea en silencio, sin error.
func TestInstall_SkipsAdapterWithEmptyAgentsDir(t *testing.T) {
	home := t.TempDir()
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	res, err := Install([]AgentAdapter{fakeAdapter{agent: model.AgentCodex, dir: ""}}, home, backupDir, mustResources(t))
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != model.AgentCodex {
		t.Fatalf("got Skipped=%v, want [codex]", res.Skipped)
	}
	if len(res.Resources) != 0 {
		t.Fatalf("got %d resource results for a skipped adapter, want 0", len(res.Resources))
	}
}

// TestInstall_RejectsExternalSymlink cubre "rechazar traversal/symlinks
// externos": un symlink en el lugar de un rol, apuntando AFUERA del
// directorio de agentes, tiene que abortar — nunca seguirse y sobrescribirse.
func TestInstall_RejectsExternalSymlink(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(home, "outside-secret.md")
	if err := os.WriteFile(outside, []byte("secreto"), 0o644); err != nil {
		t.Fatal(err)
	}
	dbaPath := filepath.Join(agentsDir, roles.RoleDBA+".md")
	if err := os.Symlink(outside, dbaPath); err != nil {
		t.Skipf("symlinks not supported in this environment: %v", err)
	}

	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}
	_, err := Install(adapters, home, backupDir, mustResources(t))
	if err == nil {
		t.Fatal("Install accepted a symlink pointing outside the agents dir, and should not have")
	}

	// El archivo externo no se tocó.
	data, rerr := os.ReadFile(outside)
	if rerr != nil || string(data) != "secreto" {
		t.Fatalf("the external target was touched: data=%q err=%v", data, rerr)
	}
}
