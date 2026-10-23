package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Group-Active-IA/active-stack/internal/approval"
)

// runApprovalDispatch maneja el subcomando "approval": lee una tabla de
// Aprobación ya firmada y dice si la estructura es completa, vigente y está
// en un commit documental separado del código. NUNCA ESCRIBE — no completa la
// tabla, no firma, no toca el árbol de trabajo (J1 revisada, propose
// 2026-10-23-fase-2-pasos-1-4.md §1).
//
// El código más alto que devuelve es 2 (WRITTEN_UNVERIFIED): estructura
// válida, identidad humana NO acreditada. No existe un código 0 de éxito
// genérico a propósito — ver internal/approval, Status.ExitCode.
//
// dir es el repositorio a inspeccionar (vacío = el actual). Parámetro, no
// os.Getwd() adentro, por lo mismo que roles_dispatch.go: el test tiene que
// poder apuntarlo a un repo temporal.
func runApprovalDispatch(args []string, dir string, w io.Writer) int {
	opts, asJSON, err := parseApprovalFlags(args)
	if errors.Is(err, errApprovalHelp) {
		fmt.Fprintln(w, approvalUsage)
		return 0
	}
	if err != nil {
		fmt.Fprintf(w, "error: %v\n", err)
		fmt.Fprintln(w, approvalUsage)
		return 1
	}
	if dir == "" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(w, "error: no pude leer el directorio actual: %v\n", err)
			return 1
		}
		dir = wd
	}
	opts.RepoDir = dir

	res, err := approval.Verify(opts)
	if err != nil {
		fmt.Fprintf(w, "error: %v\n", err)
		return 1
	}

	if asJSON {
		out, err := approval.RenderJSON(res)
		if err != nil {
			fmt.Fprintf(w, "error: %v\n", err)
			return 1
		}
		if _, err := w.Write(out); err != nil {
			fmt.Fprintf(w, "error: escribir la salida: %v\n", err)
			return 1
		}
		return res.Status.ExitCode()
	}

	fmt.Fprint(w, approval.RenderText(res))
	return res.Status.ExitCode()
}

func parseApprovalFlags(args []string) (approval.Options, bool, error) {
	fs := flag.NewFlagSet("approval", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // el usage propio sale por approvalUsage, no el de flag
	change := fs.String("change", "", "ruta de la propuesta a verificar, relativa al repo")
	commit := fs.String("commit", "", "SHA o referencia del commit documental de aprobación (C)")
	currentRef := fs.String("current-ref", "", "revisión contra la que se compara para detectar STALE (default HEAD)")
	asJSON := fs.Bool("json", false, "salida para consumir sin parsear texto")
	help := fs.Bool("help", false, "")

	if err := fs.Parse(args); err != nil {
		return approval.Options{}, false, err
	}
	if *help {
		return approval.Options{}, false, errApprovalHelp
	}
	if *change == "" {
		return approval.Options{}, false, errors.New("--change es obligatorio: ruta de la propuesta a verificar")
	}
	if *commit == "" {
		return approval.Options{}, false, errors.New("--commit es obligatorio: SHA o referencia del commit de aprobación (C)")
	}

	return approval.Options{
		ProposalPath: *change,
		Commit:       *commit,
		CurrentRef:   *currentRef,
	}, *asJSON, nil
}

var errApprovalHelp = errors.New("mostrar el usage")

const approvalUsage = `Usage: active-stack approval --change <ruta> --commit <C> [flags]

Lee una tabla de Aprobación ya firmada. NUNCA ESCRIBE: no completa la tabla,
no firma, no toca el árbol de trabajo.

El código de salida más alto es 2 (estructura válida; identidad humana no
acreditada). 1 bloquea apply/cierre (tabla ausente, inválida u obsoleta).

Flags:
  --change <ruta>       ruta de la propuesta, relativa al repo
  --commit <C>          SHA o referencia del commit documental de aprobación
  --current-ref <ref>   revisión contra la que se detecta STALE (default HEAD)
  --json                salida para consumir sin parsear texto`
