package roledefs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/roles"
)

// TestInstall_SecondAdapterFails_FirstAdapterWriteIsUndone cubre QA-2
// ("Deuda aceptada" del propose): "falta el caso de fallo entre adapters del
// mismo paso: el primero escribe, el segundo falla; exigir recuperación del
// primero y conservación de ediciones posteriores".
//
// adapter 1 (claude) escribe los archivos de roles con éxito. adapter 2 (codex)
// falla porque su AgentsDir resuelve a una ruta donde YA existe un archivo
// regular en el lugar de un componente del directorio (ENOTDIR real, no un
// fallo simulado con un booleano) — falla determinística sin depender de
// permisos de archivo, que bajo root no fallarían.
//
// Es RED first: hoy Install no revierte nada cuando un adapter posterior
// falla — este test tiene que fallar por su aserción (los archivos de
// claude siguen ahí) antes de implementar el rollback cruzado.
func TestInstall_SecondAdapterFails_FirstAdapterWriteIsUndone(t *testing.T) {
	home := t.TempDir()
	claudeDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")

	// Deja una edición personal, no relacionada con los roles, en el mismo
	// directorio de claude ANTES de instalar — tiene que sobrevivir intacta.
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	personalPath := filepath.Join(claudeDir, "my-own-agent.md")
	personalContent := "---\nname: my-own-agent\n---\nmío, no de active-stack\n"
	if err := os.WriteFile(personalPath, []byte(personalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// La "carpeta" de codex es en realidad un archivo regular: cualquier
	// intento de os.MkdirAll(codexBogusDir, ...) falla con ENOTDIR de verdad.
	codexParent := filepath.Join(home, "codex-blocker")
	if err := os.WriteFile(codexParent, []byte("soy un archivo, no un directorio"), 0o644); err != nil {
		t.Fatal(err)
	}
	codexBogusDir := filepath.Join(codexParent, "agents")

	adapters := []AgentAdapter{
		fakeAdapter{agent: model.AgentClaude, dir: claudeDir},
		fakeAdapter{agent: model.AgentCodex, dir: codexBogusDir},
	}

	_, err := Install(adapters, home, backupDir, mustResources(t))
	if err == nil {
		t.Fatal("Install succeeded despite the second adapter's write failing — the fixture is wrong")
	}

	// El primero escribió: su escritura tiene que haberse revertido.
	for _, id := range canonicalRoleIDs {
		path := filepath.Join(claudeDir, id+".md")
		if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
			t.Errorf("%s: still present after the second adapter failed — not rolled back (stat err: %v)", id, statErr)
		}
	}

	// La edición personal, no tocada por este Install, se conserva.
	got, rerr := os.ReadFile(personalPath)
	if rerr != nil {
		t.Fatalf("personal file disappeared during rollback: %v", rerr)
	}
	if string(got) != personalContent {
		t.Fatalf("personal file was modified during rollback: got %q", got)
	}

	// Tampoco quedó un manifiesto de ownership a medio escribir para claude.
	m, merr := loadManifest(home, model.AgentClaude)
	if merr != nil {
		t.Fatalf("loadManifest: %v", merr)
	}
	if len(m.Resources) != 0 {
		t.Fatalf("ownership manifest for claude was persisted despite the failed install: %+v", m.Resources)
	}
}

// TestInstall_SecondAdapterFails_UndoesAnUpdateToo triangula el mismo
// escenario con un "update" (no sólo un "create"): si el primer adapter
// ACTUALIZÓ un archivo propio intacto y el segundo falla, el archivo tiene
// que volver a su contenido anterior, no quedar en el nuevo.
func TestInstall_SecondAdapterFails_UndoesAnUpdateToo(t *testing.T) {
	home := t.TempDir()
	claudeDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	claudeOnly := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: claudeDir}}

	v1 := mustResources(t)
	if _, err := Install(claudeOnly, home, backupDir, v1); err != nil {
		t.Fatalf("install v1: %v", err)
	}
	before, err := os.ReadFile(filepath.Join(claudeDir, roles.RoleDBA+".md"))
	if err != nil {
		t.Fatal(err)
	}

	v2 := append([]Resource{}, v1...)
	for i, r := range v2 {
		if r.Name == roles.RoleDBA+".md" {
			v2[i].Content = append(append([]byte{}, r.Content...), []byte("\nv2\n")...)
		}
	}

	codexParent := filepath.Join(home, "codex-blocker")
	if err := os.WriteFile(codexParent, []byte("archivo, no directorio"), 0o644); err != nil {
		t.Fatal(err)
	}
	codexBogusDir := filepath.Join(codexParent, "agents")

	adapters := []AgentAdapter{
		fakeAdapter{agent: model.AgentClaude, dir: claudeDir},
		fakeAdapter{agent: model.AgentCodex, dir: codexBogusDir},
	}
	_, err = Install(adapters, home, backupDir, v2)
	if err == nil {
		t.Fatal("Install succeeded despite the second adapter's write failing — the fixture is wrong")
	}

	after, rerr := os.ReadFile(filepath.Join(claudeDir, roles.RoleDBA+".md"))
	if rerr != nil {
		t.Fatalf("role-dba.md disappeared during rollback of an update: %v", rerr)
	}
	if string(after) != string(before) {
		t.Fatalf("role-dba.md was left on the NEW content after a failed install; rollback of an update did not restore it.\nbefore: %q\nafter:  %q", before, after)
	}
}

// TestInstall_ConcurrentEditAbortsWrite cubre "cambio concurrente aborta la
// escritura": si el archivo cambia entre la decisión (propio intacto) y la
// escritura final, Install no pisa esa edición — aborta.
func TestInstall_ConcurrentEditAbortsWrite(t *testing.T) {
	home := t.TempDir()
	claudeDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: claudeDir}}

	v1 := mustResources(t)
	if _, err := Install(adapters, home, backupDir, v1); err != nil {
		t.Fatalf("install v1: %v", err)
	}

	v2 := append([]Resource{}, v1...)
	for i, r := range v2 {
		if r.Name == roles.RoleDBA+".md" {
			v2[i].Content = append(append([]byte{}, r.Content...), []byte("\nv2\n")...)
		}
	}

	dbaPath := filepath.Join(claudeDir, roles.RoleDBA+".md")
	racing := "edición concurrente que llegó justo antes de escribir\n"
	testHookBeforeWrite = func(path string) {
		if path != dbaPath {
			return
		}
		if err := os.WriteFile(path, []byte(racing), 0o644); err != nil {
			t.Fatalf("inject racing edit: %v", err)
		}
	}
	defer func() { testHookBeforeWrite = nil }()

	_, err := Install(adapters, home, backupDir, v2)
	if err == nil {
		t.Fatal("Install reported success despite a concurrent edit landing mid-write")
	}

	got, rerr := os.ReadFile(dbaPath)
	if rerr != nil {
		t.Fatalf("read role-dba.md after aborted concurrent write: %v", rerr)
	}
	if string(got) != racing {
		t.Fatalf("the racing edit was overwritten anyway: got %q, want %q", got, racing)
	}
}
