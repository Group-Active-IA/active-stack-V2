package approval

import (
	"fmt"
	"os/exec"
	"strings"
)

// runGit corre git en dir y devuelve stdout. Mismo patrón que
// internal/roles/source.go: nada de libgit2, el binario real.
func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(ee.Stderr)))
		}
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return string(out), nil
}

// resolveCommit devuelve el SHA completo de ref, o error si no resuelve a un
// commit (SHA ambiguo, inexistente, o que no es un commit).
func resolveCommit(dir, ref string) (string, error) {
	out, err := runGit(dir, "rev-parse", "--verify", ref+"^{commit}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// isDirty dice si el árbol de trabajo tiene cambios (trackeados o no) frente a
// HEAD. Un árbol documental sucio invalida la comparación: lo que se está
// leyendo en disco puede no ser lo que está commiteado.
func isDirty(dir string) (bool, error) {
	out, err := runGit(dir, "status", "--porcelain=v1", "-uall")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

// changedPaths devuelve los archivos que el commit c modifica respecto de su
// primer padre. Un commit sin padre (raíz) se compara contra el árbol vacío.
func changedPaths(dir, c string) ([]string, error) {
	hasParent, err := commitHasParent(dir, c)
	if err != nil {
		return nil, err
	}
	var out string
	if hasParent {
		out, err = runGit(dir, "diff-tree", "--no-commit-id", "--name-only", "-r", c)
	} else {
		out, err = runGit(dir, "diff-tree", "--no-commit-id", "--name-only", "-r", "--root", c)
	}
	if err != nil {
		return nil, err
	}
	return splitLines(out), nil
}

func commitHasParent(dir, c string) (bool, error) {
	_, err := runGit(dir, "rev-parse", "--verify", c+"^1")
	if err != nil {
		// rev-parse sobre ^1 de un commit raíz sale con error: no hay padre.
		return false, nil
	}
	return true, nil
}

// parentOf devuelve el SHA completo del primer padre de c, o "" si c es raíz.
func parentOf(dir, c string) (string, error) {
	has, err := commitHasParent(dir, c)
	if err != nil || !has {
		return "", err
	}
	out, err := runGit(dir, "rev-parse", "--verify", c+"^1")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// showFile devuelve el contenido de path en la revisión rev. ok=false cuando
// el archivo no existe en esa revisión (no es un error de git: es un dato).
func showFile(dir, rev, path string) (content string, ok bool, err error) {
	out, err := runGit(dir, "show", rev+":"+path)
	if err != nil {
		// git show sale con error tanto si el path no existe como por otros
		// fallos; distinguimos con cat-file -e, que es exacto para "existe".
		exists, errCheck := pathExistsAt(dir, rev, path)
		if errCheck != nil {
			return "", false, errCheck
		}
		if !exists {
			return "", false, nil
		}
		return "", false, err
	}
	return out, true, nil
}

func pathExistsAt(dir, rev, path string) (bool, error) {
	cmd := exec.Command("git", "cat-file", "-e", rev+":"+path)
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// blobHash devuelve el hash de git para path en rev, u ok=false si no existe.
func blobHash(dir, rev, path string) (hash string, ok bool, err error) {
	exists, err := pathExistsAt(dir, rev, path)
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}
	out, err := runGit(dir, "rev-parse", "--verify", rev+":"+path)
	if err != nil {
		return "", false, err
	}
	return strings.TrimSpace(out), true, nil
}

// isAncestor dice si p es ancestro de c (o el mismo commit).
func isAncestor(dir, p, c string) (bool, error) {
	cmd := exec.Command("git", "merge-base", "--is-ancestor", p, c)
	cmd.Dir = dir
	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 1 {
		return false, nil
	}
	return false, err
}

// commitMessage devuelve el mensaje completo (asunto + cuerpo) del commit c.
func commitMessage(dir, c string) (string, error) {
	out, err := runGit(dir, "log", "-1", "--format=%B", c)
	if err != nil {
		return "", err
	}
	return out, nil
}

func splitLines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		l = strings.TrimSpace(l)
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}
