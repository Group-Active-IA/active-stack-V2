package reports

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// runGit corre git en dir y devuelve stdout. Mismo patrón que
// internal/approval/git.go: el binario real, nada de libgit2.
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

// canonicalRoot resuelve la raíz real del worktree que contiene dir,
// siguiendo symlinks. "Registrar antes del despacho repo/worktree canónico"
// (propose §2a): si dir no ES esa raíz, es un subdirectorio o un worktree
// ajeno, y Capture lo rechaza en vez de adivinar hacia dónde escribir.
func canonicalRoot(dir string) (string, error) {
	resolvedDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return "", fmt.Errorf("resolver symlinks de %s: %w", dir, err)
	}
	out, err := runGit(resolvedDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("%s no es un worktree git: %w", dir, err)
	}
	top := strings.TrimSpace(out)
	resolvedTop, err := filepath.EvalSymlinks(top)
	if err != nil {
		return "", fmt.Errorf("resolver symlinks de %s: %w", top, err)
	}
	return resolvedTop, nil
}

// requireCanonicalRoot rechaza RepoDir si no es exactamente la raíz de su
// worktree ("worktree ajeno" cubre tanto un subdirectorio como la raíz de
// OTRO worktree vinculado al mismo repo).
func requireCanonicalRoot(dir string) (string, error) {
	resolvedDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return "", fmt.Errorf("resolver symlinks de %s: %w", dir, err)
	}
	top, err := canonicalRoot(dir)
	if err != nil {
		return "", err
	}
	if filepath.Clean(resolvedDir) != filepath.Clean(top) {
		return "", fmt.Errorf("RepoDir %q no es la raíz del worktree canónico (raíz real: %q); pasá la raíz, no un subdirectorio ni un worktree ajeno", dir, top)
	}
	return top, nil
}

// resolveWithinRoot resuelve rel (una ruta limpia, ya validada por
// safeSegment en cada segmento) dentro de root y rechaza si el resultado —
// siguiendo symlinks en cualquier ancestro ya existente — queda afuera de
// root. Cubre "symlink fuera de raíz": un ancestro simlinkeado hacia afuera
// no permite escribir ahí aunque el nombre del segmento sea válido.
func resolveWithinRoot(root string, segments []string) (string, error) {
	rootClean := filepath.Clean(root)
	full := filepath.Join(append([]string{rootClean}, segments...)...)
	if full != rootClean && !strings.HasPrefix(full, rootClean+string(filepath.Separator)) {
		return "", fmt.Errorf("la ruta resuelta %q queda afuera de la raíz %q", full, root)
	}

	// Caminar cada ancestro desde root hacia abajo; si alguno YA existe y es
	// un symlink que resuelve afuera de root, rechazar antes de que Capture
	// cree nada debajo ("symlink fuera de raíz"). Un segmento que todavía no
	// existe no puede ser un symlink: ahí termina la validación.
	cur := rootClean
	for _, seg := range segments {
		cur = filepath.Join(cur, seg)
		info, err := os.Lstat(cur)
		if err != nil {
			break
		}
		if info.Mode()&os.ModeSymlink != 0 {
			resolved, err := filepath.EvalSymlinks(cur)
			if err != nil {
				return "", fmt.Errorf("resolver symlink %q: %w", cur, err)
			}
			if resolved != rootClean && !strings.HasPrefix(resolved, rootClean+string(filepath.Separator)) {
				return "", fmt.Errorf("un ancestro de %q es un symlink que escapa de la raíz %q", full, root)
			}
			cur = resolved
		}
	}
	return full, nil
}
