package approval

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Mismo patrón que internal/roles/source_test.go: repos git temporales de
// verdad, sin mockear git.

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

func newRepo(t *testing.T) string {
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

func commitAll(t *testing.T, dir, msg string) string {
	t.Helper()
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-qm", msg)
	return strings.TrimSpace(git(t, dir, "rev-parse", "HEAD"))
}

// proposalBody arma el cuerpo de una propuesta con la sección Aprobación.
// Sin fila de datos: la propuesta "revisada" todavía no está aprobada.
func proposalBody(extra string) string {
	return "# Propose\n\n" + extra + "\n## Aprobación\n\n" +
		"| Elección | Persona | Fecha | Alcance |\n" +
		"|---|---|---|---|\n"
}

// approvalRowLine es la fila que el commit documental de aprobación agrega.
func approvalRowLine(eleccion, persona, fecha, alcance string) string {
	return "| " + eleccion + " | " + persona + " | " + fecha + " | " + alcance + " |\n"
}
