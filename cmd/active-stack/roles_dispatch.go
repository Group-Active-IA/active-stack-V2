package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/Group-Active-IA/active-stack/internal/roles"
)

// runRolesDispatch maneja el subcomando "roles": lee un diff y devuelve qué
// revisores corresponden.
//
// POR QUÉ EXISTE. El bloque del orquestador dice "convocá al revisor que la capa
// toca". Está bien escrito y no funciona: el que decide si una capa "se toca" es
// el mismo que se ahorra el trabajo de convocar, y nadie revisa esa decisión.
// Medido sobre tres corridas reales, de diez roles disponibles el orquestador
// convocó siempre los mismos cuatro y el revisor senior no corrió ni una vez.
// No convocar no deja rastro, así que nadie se entera.
//
// Este comando convierte esa decisión en un dato: el orquestador deja de juzgar
// y pasa a ejecutar y leer. NO CONVOCA A NADIE — devuelve la lista, y la
// convocatoria sigue siendo suya, que es el único que puede hablar con la
// persona.
//
// dir es el directorio desde donde se lee git (vacío = el actual). Es un
// parámetro y no una llamada a os.Getwd() adentro para que el test pueda
// apuntarlo a un temporal que NO es un repo, que es el caso que hay que probar.
//
// Devuelve el código de salida del proceso (0 = ok, 1 = error de uso o de git).
func runRolesDispatch(args []string, dir string, w io.Writer) int {
	opts, asJSON, err := parseRolesFlags(args)
	if errors.Is(err, errRolesHelp) {
		fmt.Fprintln(w, rolesUsage)
		return 0
	}
	if err != nil {
		fmt.Fprintf(w, "error: %v\n", err)
		fmt.Fprintln(w, rolesUsage)
		return 1
	}
	opts.Dir = dir

	files, err := roles.Collect(opts)
	if err != nil {
		// Fuera de un repo y sin --files no se inventa un diff vacío: "no pude
		// averiguarlo" y "no hay nada que revisar" se leen igual en la salida y
		// no son lo mismo.
		fmt.Fprintf(w, "error: %v\n", err)
		if errors.Is(err, roles.ErrNoRepo) {
			fmt.Fprintln(w, rolesUsage)
		}
		return 1
	}

	rep := roles.Classify(files)

	if asJSON {
		out, err := roles.RenderJSON(rep)
		if err != nil {
			fmt.Fprintf(w, "error: %v\n", err)
			return 1
		}
		if _, err := w.Write(out); err != nil {
			fmt.Fprintf(w, "error: escribir la salida: %v\n", err)
			return 1
		}
		return 0
	}

	fmt.Fprint(w, roles.RenderText(rep))
	return 0
}

// parseRolesFlags parsea a mano porque --files es variádico y el paquete flag
// no lo soporta: `--files a.go b.go` toma un solo valor y deja b.go suelto.
func parseRolesFlags(args []string) (roles.Options, bool, error) {
	var opts roles.Options
	var asJSON bool

	for i := 0; i < len(args); i++ {
		switch arg := args[i]; arg {
		case "--json":
			asJSON = true
		case "--base":
			if i+1 >= len(args) {
				return opts, false, errors.New("--base necesita una referencia (ej: --base main)")
			}
			i++
			opts.Base = args[i]
		case "--files":
			// Se come todo lo que sigue hasta la próxima flag. Una lista vacía
			// acá NO es "sin archivos": es un uso mal escrito, y dejarla pasar
			// la convertiría en un diff vacío inventado.
			for i+1 < len(args) && !isFlag(args[i+1]) {
				i++
				opts.Files = append(opts.Files, args[i])
			}
			if len(opts.Files) == 0 {
				return opts, false, errors.New("--files necesita al menos un archivo")
			}
		case "--help", "-h":
			return opts, false, errRolesHelp
		default:
			return opts, false, fmt.Errorf("flag desconocida: %q", arg)
		}
	}

	if opts.Base != "" && len(opts.Files) > 0 {
		return opts, false, errors.New("--base y --files son excluyentes: " +
			"--files no toca git y --base sólo tiene sentido contra el repo")
	}
	return opts, asJSON, nil
}

func isFlag(s string) bool { return len(s) > 1 && s[0] == '-' }

// errRolesHelp separa "pediste ayuda" de "te equivocaste". Sin la distinción,
// --help sale por stderr con código 1 y cualquier script que lo llame lo lee
// como una falla.
var errRolesHelp = errors.New("mostrar el usage")

const rolesUsage = `Usage: active-stack roles [flags]

Lee un diff y devuelve qué revisores corresponden. No convoca a nadie: la
convocatoria sigue siendo del orquestador.

Sin flags se usan los archivos cambiados y no commiteados.

Flags:
  --base <ref>         comparar contra esa referencia (suma los sin trackear)
  --files a.go b.go    lista explícita, sin tocar git
  --json               salida para consumir sin parsear texto

Los roles que NO corresponden también salen, con el motivo: una ausencia sin
motivo es indistinguible de un olvido.`
