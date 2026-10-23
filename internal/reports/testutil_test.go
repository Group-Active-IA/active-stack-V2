package reports

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// gitCheckIgnore corre `git check-ignore -v` para path y devuelve si está
// ignorado. A diferencia de testGit, NO trata exit 1 como fallo: exit 1 es
// la forma normal en que `check-ignore` dice "esto no está ignorado" — el
// resultado que varios tests necesitan afirmar explícitamente.
func gitCheckIgnore(t *testing.T, dir, path string) (ignored bool, verbose string) {
	t.Helper()
	cmd := exec.Command("git", "check-ignore", "-v", path)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err == nil {
		return true, string(out)
	}
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 1 {
		return false, ""
	}
	t.Fatalf("git check-ignore -v %s: %v\n%s", path, err, out)
	return false, ""
}

func testGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t",
		"GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if _, lookErr := exec.LookPath("git"); lookErr != nil {
			t.Skipf("sin git en el PATH: %v", lookErr)
		}
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

// newRepo arma un repo temporal con la copia real del .gitignore del
// repositorio (no un .gitignore inventado para el test): así el
// StatusRetained se prueba contra las reglas que de verdad corren en este
// proyecto, no contra una aproximación.
func newRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	testGit(t, dir, "init", "-q", "-b", "main")

	realGitignore, err := os.ReadFile(findRepoRoot(t) + "/.gitignore")
	if err != nil {
		t.Fatalf("leer .gitignore real: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), realGitignore, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("repo de prueba\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	testGit(t, dir, "add", ".")
	testGit(t, dir, "commit", "-qm", "inicial")
	return dir
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	out := testGit(t, wd, "rev-parse", "--show-toplevel")
	return trimNL(out)
}

func trimNL(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}

func baseInput(repo string) CaptureInput {
	return CaptureInput{
		RepoDir:    repo,
		Change:     "fase-2-pasos-1-4",
		Phase:      "verify",
		Round:      1,
		Role:       "backend",
		DispatchID: "d-1",
		SHA:        "abc123",
		Source: Source{
			Contract: ContractRawText,
			Text:     "informe de prueba sin nada sensible.",
		},
	}
}
