package roles

import (
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// ErrNoRepo es lo que devuelve Collect cuando no hay repositorio git y tampoco
// se pasaron archivos.
//
// POR QUÉ ES UN ERROR Y NO UNA LISTA VACÍA: una lista vacía significa "nada que
// revisar", y eso es una respuesta legítima. Devolverla cuando en realidad no se
// pudo mirar convierte "no pude averiguarlo" en "no hay nada", que son la misma
// salida y no son lo mismo. Es el mismo modo de falla que el bloque del
// orquestador describe para `git log` sobre una ruta no trackeada: imprime nada,
// sale 0, y quien sólo mira el código de salida lee el vacío como un dato.
var ErrNoRepo = errors.New(
	"no hay repositorio git acá y no se pasaron archivos: usá --files <a> <b> " +
		"o corré el comando dentro de un repo")

// Options describe de dónde sale la lista de archivos, en orden de preferencia:
// Files explícito, si no Base, si no los cambios no commiteados.
type Options struct {
	// Dir es el directorio desde donde se corre git. Vacío = el actual.
	Dir string
	// Base es la referencia contra la cual comparar (--base).
	Base string
	// Files es la lista explícita (--files). Si está, git no se toca.
	Files []string
}

// Collect devuelve los archivos del diff.
//
// Con --base se compara el árbol de trabajo contra la referencia y ADEMÁS se
// suman los archivos sin trackear. Un archivo nuevo que todavía no entró al
// índice es parte del cambio aunque `git diff` no lo vea, y dejarlo afuera es
// exactamente el error que este comando existe para no cometer: ante la duda,
// se convoca.
func Collect(o Options) ([]string, error) {
	if len(o.Files) > 0 {
		return o.Files, nil
	}

	if _, err := runGit(o.Dir, "rev-parse", "--show-toplevel"); err != nil {
		return nil, ErrNoRepo
	}

	visto := map[string]bool{}
	add := func(p string) {
		p = strings.ReplaceAll(strings.TrimSpace(p), "\\", "/")
		if p != "" && !visto[p] {
			visto[p] = true
		}
	}

	if o.Base != "" {
		out, err := runGit(o.Dir, "diff", "--name-only", "-z", o.Base, "--")
		if err != nil {
			return nil, fmt.Errorf("git diff contra %q: %w", o.Base, err)
		}
		for _, p := range splitNUL(out) {
			add(p)
		}
	} else {
		// -uall NO es opcional. Sin él `git status` colapsa un directorio sin
		// trackear en una sola entrada terminada en "/": un `svc/` nuevo con
		// `auth.py` adentro llega como "svc/", ningún token del piso matchea, y
		// el piso NO se activa. Ese es exactamente el error caro —un revisor de
		// MENOS, que cuesta un bug en producción, contra uno de más que cuesta
		// tokens— y es invisible en la salida: se lee igual que un diff que
		// legítimamente no tocaba autenticación. Cazado por
		// TestSinArgumentosSonLosCambiosNoCommiteados (2026-09-13), que con la
		// versión sin -uall recibía "migrations/" en lugar del .sql.
		out, err := runGit(o.Dir, "status", "--porcelain=v1", "-z", "-uall")
		if err != nil {
			return nil, fmt.Errorf("git status: %w", err)
		}
		for _, p := range rutasDeStatus(out) {
			add(p)
		}
	}

	// Los sin trackear en los dos modos. `git status -z` ya los trae; con
	// --base hace falta pedirlos aparte.
	if o.Base != "" {
		out, err := runGit(o.Dir, "ls-files", "--others", "--exclude-standard", "-z")
		if err == nil {
			for _, p := range splitNUL(out) {
				add(p)
			}
		}
	}

	files := make([]string, 0, len(visto))
	for p := range visto {
		files = append(files, p)
	}
	// Orden estable: el informe se pega en el registro del reparto y dos
	// corridas sobre el mismo árbol tienen que dar el mismo texto. El orden de
	// un map en Go está aleatorizado a propósito.
	sort.Strings(files)
	return files, nil
}

// rutasDeStatus parsea `git status --porcelain=v1 -z`.
//
// El formato NO es una ruta por registro: un renombrado ocupa DOS registros
// consecutivos, `R  destino\0origen\0`. Leído como uno solo, el origen del
// renombrado entra como si fuera un archivo cambiado más — y como acá eso sólo
// puede AGREGAR revisores, el error sería invisible en la salida y sólo se
// notaría como un revisor de más. Se consume explícitamente.
func rutasDeStatus(out string) []string {
	campos := splitNUL(out)
	var rutas []string
	for i := 0; i < len(campos); i++ {
		entrada := campos[i]
		if len(entrada) < 4 {
			continue
		}
		x, y := entrada[0], entrada[1]
		ruta := entrada[3:]
		rutas = append(rutas, ruta)
		if x == 'R' || x == 'C' || y == 'R' || y == 'C' {
			i++ // el registro siguiente es la ruta original: se descarta
		}
	}
	return rutas
}

func splitNUL(s string) []string {
	var out []string
	for _, p := range strings.Split(s, "\x00") {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
