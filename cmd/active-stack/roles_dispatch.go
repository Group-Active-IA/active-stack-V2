package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Group-Active-IA/active-stack/internal/agents/claude"
	roledefsinstaller "github.com/Group-Active-IA/active-stack/internal/harness/roledefs"
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
// exitDefinitionMissing es el código de salida cuando un rol CONVOCADO no
// tiene su role-*.md instalado (ausente o vacío). Distinto de 1 (error de
// uso/git) a propósito: "clasificación vacía" y "definición faltante" son
// afirmaciones distintas (propose fase 2 paso 4, "## 4 — instalar-roles"),
// y un script que sólo mira "código != 0" no tiene por qué distinguirlas,
// pero uno que sí mira el código puede.
const exitDefinitionMissing = 3

// Devuelve el código de salida del proceso (0 = ok, 1 = error de uso o de
// git, 3 = un rol convocado no tiene su definición instalada o está vacía).
func runRolesDispatch(args []string, dir string, w io.Writer) int {
	opts, asJSON, homeDir, err := parseRolesFlags(args)
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
	} else {
		fmt.Fprint(w, roles.RenderText(rep))
	}

	missing, err := missingDefinitionsForConvened(rep.Convened(), homeDir)
	if err != nil {
		// No se pudo ni mirar el disco (ej. DefaultResources roto): esto no
		// es "no hay rol convocado", es "no pude averiguarlo" — mismo
		// principio que roles.ErrNoRepo arriba.
		fmt.Fprintf(w, "error: no se pudo verificar la instalación de los roles: %v\n", err)
		return 1
	}
	if len(missing) > 0 {
		fmt.Fprintf(w, "\nbloqueado: definición de rol faltante o vacía — %s\n", strings.Join(missing, ", "))
		fmt.Fprintln(w, "convocar o cerrar con un rol convocado sin definición instalada "+
			"invoca un sub-agente que no existe. Instalar el harness role-definitions "+
			"(active-stack install, harness role-definitions) antes de continuar.")
		return exitDefinitionMissing
	}

	return 0
}

// missingDefinitionsForConvened filtra los recursos canónicos a sólo los
// roles CONVOCADOS por el diff y pregunta a roledefs.Health si alguno falta
// o está vacío. Hoy sólo Claude Code tiene AgentsDir soportado (J3/J4, TBD),
// así que esto es lo único que se verifica — ver
// internal/harness/roledefs/doc.go.
func missingDefinitionsForConvened(convened []string, homeDir string) ([]string, error) {
	if len(convened) == 0 {
		return nil, nil
	}
	all, err := roledefsinstaller.DefaultResources()
	if err != nil {
		return nil, err
	}
	want := make(map[string]bool, len(convened))
	for _, r := range convened {
		want[r+".md"] = true
	}
	var filtered []roledefsinstaller.Resource
	for _, res := range all {
		if want[res.Name] {
			filtered = append(filtered, res)
		}
	}
	if len(filtered) == 0 {
		return nil, nil
	}

	issues, err := roledefsinstaller.Health([]roledefsinstaller.AgentAdapter{claude.NewAdapter()}, homeDir, filtered)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(issues))
	for _, iss := range issues {
		out = append(out, fmt.Sprintf("%s (%s)", iss.Resource, iss.Problem))
	}
	return out, nil
}

// parseRolesFlags parsea a mano porque --files es variádico y el paquete flag
// no lo soporta: `--files a.go b.go` toma un solo valor y deja b.go suelto.
func parseRolesFlags(args []string) (roles.Options, bool, string, error) {
	var opts roles.Options
	var asJSON bool
	var homeDir string

	for i := 0; i < len(args); i++ {
		switch arg := args[i]; arg {
		case "--json":
			asJSON = true
		case "--base":
			if i+1 >= len(args) {
				return opts, false, "", errors.New("--base necesita una referencia (ej: --base main)")
			}
			i++
			opts.Base = args[i]
		case "--home":
			// Sólo para pruebas y para el check de definiciones instaladas:
			// roles.Collect nunca lo usa. Vacío (el default) resuelve al
			// home real en tiempo de chequeo, no acá — así una corrida sin
			// --home sigue funcionando aunque os.UserHomeDir() falle en un
			// entorno sin HOME, porque ese fallo recién importaría si hay
			// algún rol convocado que necesite el chequeo.
			if i+1 >= len(args) {
				return opts, false, "", errors.New("--home necesita una ruta")
			}
			i++
			homeDir = args[i]
		case "--files":
			// Se come todo lo que sigue hasta la próxima flag. Una lista vacía
			// acá NO es "sin archivos": es un uso mal escrito, y dejarla pasar
			// la convertiría en un diff vacío inventado.
			for i+1 < len(args) && !isFlag(args[i+1]) {
				i++
				opts.Files = append(opts.Files, args[i])
			}
			if len(opts.Files) == 0 {
				return opts, false, "", errors.New("--files necesita al menos un archivo")
			}
		case "--help", "-h":
			return opts, false, "", errRolesHelp
		default:
			return opts, false, "", fmt.Errorf("flag desconocida: %q", arg)
		}
	}

	if opts.Base != "" && len(opts.Files) > 0 {
		return opts, false, "", errors.New("--base y --files son excluyentes: " +
			"--files no toca git y --base sólo tiene sentido contra el repo")
	}
	if homeDir == "" {
		if h, err := os.UserHomeDir(); err == nil {
			homeDir = h
		}
	}
	return opts, asJSON, homeDir, nil
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
  --home <ruta>        home a usar para el chequeo de definiciones instaladas
                        (default: el home real; sólo para pruebas o setups
                        no estándar)

Los roles que NO corresponden también salen, con el motivo: una ausencia sin
motivo es indistinguible de un olvido.

Un rol CONVOCADO sin su role-*.md instalado (harness role-definitions)
bloquea con código 3, distinto del 1 de uso/git: convocar a un sub-agente
que no existe no es lo mismo que un diff vacío.`
