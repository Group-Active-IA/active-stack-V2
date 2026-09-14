package roles

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command(gitPath(t), args...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t",
		"GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func gitPath(t *testing.T) string {
	t.Helper()
	p, err := exec.LookPath("git")
	if err != nil {
		t.Skipf("sin git en el PATH: %v", err)
	}
	return p
}

func repo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	git(t, dir, "init", "-q", "-b", "main")
	writeFile(t, dir, "README.md", "hola\n")
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-qm", "inicial")
	return dir
}

func writeFile(t *testing.T, dir, rel, body string) {
	t.Helper()
	full := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestArchivosExplicitosNoTocanGit: --files es la tercera preferencia y no
// necesita repositorio.
func TestArchivosExplicitosNoTocanGit(t *testing.T) {
	dir := t.TempDir() // NO es un repo

	got, err := Collect(Options{Dir: dir, Files: []string{"a.go", "b.sql"}})
	if err != nil {
		t.Fatalf("Collect con --files falló fuera de un repo: %v", err)
	}
	if strings.Join(got, ",") != "a.go,b.sql" {
		t.Errorf("Collect devolvió %v y esperaba la lista explícita tal cual", got)
	}
}

// TestSinRepoYSinArchivosEsUnError: no se inventa un diff vacío. Un diff vacío
// significa "nada que revisar" y eso es una respuesta, no un fallo silencioso.
func TestSinRepoYSinArchivosEsUnError(t *testing.T) {
	_, err := Collect(Options{Dir: t.TempDir()})
	if err == nil {
		t.Fatal("Collect fuera de un repo y sin --files devolvió nil: un diff vacío " +
			"inventado se lee igual que 'nada que revisar', y no es lo mismo")
	}
	if !strings.Contains(err.Error(), "--files") {
		t.Errorf("el error no dice cómo salir de la situación: %v", err)
	}
}

// TestSinArgumentosSonLosCambiosNoCommiteados.
func TestSinArgumentosSonLosCambiosNoCommiteados(t *testing.T) {
	dir := repo(t)
	writeFile(t, dir, "migrations/003_tenant.sql", "create table t();\n") // sin trackear
	writeFile(t, dir, "README.md", "hola\nmundo\n")                       // modificado

	got, err := Collect(Options{Dir: dir})
	if err != nil {
		t.Fatalf("Collect error = %v", err)
	}
	quiero := map[string]bool{"migrations/003_tenant.sql": true, "README.md": true}
	for f := range quiero {
		if !contiene(got, f) {
			t.Errorf("Collect no trajo %q; devolvió %v", f, got)
		}
	}
}

// TestBaseComparaContraLaReferencia.
func TestBaseComparaContraLaReferencia(t *testing.T) {
	dir := repo(t)
	base := strings.TrimSpace(git(t, dir, "rev-parse", "HEAD"))
	writeFile(t, dir, "api/auth.py", "def login(): ...\n")
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-qm", "auth")

	got, err := Collect(Options{Dir: dir, Base: base})
	if err != nil {
		t.Fatalf("Collect error = %v", err)
	}
	if !contiene(got, "api/auth.py") {
		t.Errorf("--base no trajo el archivo del commit posterior: %v", got)
	}
	if contiene(got, "README.md") {
		t.Errorf("--base trajo un archivo que no cambió desde la referencia: %v", got)
	}
}

// TestLasRutasVienenConBarraNormal: git ya las devuelve así, y el clasificador
// depende de eso. El CI cruza Linux y GOOS=windows porque ya cazó bugs de rutas.
func TestLasRutasVienenConBarraNormal(t *testing.T) {
	dir := repo(t)
	writeFile(t, dir, "web/components/Boton.tsx", "export default 1\n")

	got, err := Collect(Options{Dir: dir})
	if err != nil {
		t.Fatalf("Collect error = %v", err)
	}
	if !contiene(got, "web/components/Boton.tsx") {
		t.Fatalf("Collect no trajo el archivo: %v (sin esto el test de barras pasa en vacío)", got)
	}
	for _, f := range got {
		if strings.Contains(f, "\\") {
			t.Errorf("Collect devolvió una ruta con barra invertida: %q", f)
		}
	}
}
