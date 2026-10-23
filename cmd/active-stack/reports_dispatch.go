package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Group-Active-IA/active-stack/internal/reports"
)

// runReportsDispatch maneja el subcomando "reports". Hoy sólo existe
// "reports capture" (núcleo, propose 2026-10-23-fase-2-pasos-1-4.md §2a):
// persiste un informe ya resuelto (texto, identidad, SHA/diff) antes del
// despacho de la fase siguiente. Los adapters de proveedor (2b Claude, 2c
// Codex, 2d OpenCode) se resuelven según su contrato versionado.
//
// dir es el repositorio a inspeccionar (vacío = el actual), mismo patrón que
// approval_dispatch.go y roles_dispatch.go: parámetro, no os.Getwd() adentro,
// para que el test pueda apuntarlo a un repo temporal.
func runReportsDispatch(args []string, dir string, w io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(w, reportsUsage)
		return 1
	}
	switch args[0] {
	case "capture":
		return runReportsCaptureDispatch(args[1:], dir, w)
	case "--help", "-help", "help":
		fmt.Fprintln(w, reportsUsage)
		return 0
	default:
		fmt.Fprintf(w, "error: subcomando desconocido %q\n", args[0])
		fmt.Fprintln(w, reportsUsage)
		return 1
	}
}

func runReportsCaptureDispatch(args []string, dir string, w io.Writer) int {
	in, asJSON, err := parseReportsCaptureFlags(args)
	if errors.Is(err, errReportsCaptureHelp) {
		fmt.Fprintln(w, reportsCaptureUsage)
		return 0
	}
	if err != nil {
		fmt.Fprintf(w, "error: %v\n", err)
		fmt.Fprintln(w, reportsCaptureUsage)
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
	in.RepoDir = dir

	res, err := reports.Capture(in)
	if err != nil {
		fmt.Fprintf(w, "error: %v\n", err)
		return 1
	}

	if asJSON {
		out, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Fprintf(w, "error: %v\n", err)
			return 1
		}
		if _, err := w.Write(append(out, '\n')); err != nil {
			fmt.Fprintf(w, "error: escribir la salida: %v\n", err)
			return 1
		}
		return res.Status.ExitCode()
	}

	fmt.Fprintf(w, "estado: %s (código %d)\n", res.Status, res.Status.ExitCode())
	fmt.Fprintf(w, "clave: %s\n", res.Key)
	fmt.Fprintf(w, "directorio: %s\n", res.Dir)
	if res.Cause != "" {
		fmt.Fprintf(w, "causa: %s\n", res.Cause)
	}
	if res.RetainedPath != "" {
		fmt.Fprintf(w, "original retenido (local, ignorado por git): %s\n", res.RetainedPath)
	}
	return res.Status.ExitCode()
}

func parseReportsCaptureFlags(args []string) (reports.CaptureInput, bool, error) {
	fs := flag.NewFlagSet("reports capture", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	change := fs.String("change", "", "slug del change OPSX (openspec/changes/<change>/)")
	phase := fs.String("phase", "", "fase OPSX que despacha el informe (ej. verify, apply)")
	round := fs.Int("round", 0, "vuelta OPSX; los reintentos del mismo intento NO la incrementan")
	role := fs.String("role", "", "rol del revisor/ejecutor (ej. backend, qa)")
	dispatchID := fs.String("dispatch-id", "", "identificador del despacho puntual")
	agent := fs.String("agent", "", "agente que generó el informe, si está disponible")
	session := fs.String("session", "", "sesión del agente, si está disponible")
	turn := fs.String("turn", "", "turno dentro de la sesión, si está disponible")
	sha := fs.String("sha", "", "commit revisado")
	diff := fs.String("diff", "", "rango de diff revisado (opcional)")
	sourceProvider := fs.String("source-provider", "", "adapter de eventos: claude, codex u opencode")
	sourceVersion := fs.String("source-version", "", "versión exacta admitida")
	sourceSession := fs.String("source-session", "", "sesión esperada del registro del despacho")
	sourceAgent := fs.String("source-agent", "", "hija Claude esperada del registro del despacho")
	contract := fs.String("contract", reports.ContractRawText, "contrato manual del Source")
	provider := fs.String("provider", "", "proveedor del informe, SÓLO si consta en el origen (nunca inferido del texto)")
	model := fs.String("model", "", "modelo exacto, SÓLO si consta en el origen")
	reportPath := fs.String("report", "", "archivo con el texto íntegro del informe")
	reportText := fs.String("report-text", "", "texto íntegro del informe, en vez de --report")
	asJSON := fs.Bool("json", false, "salida para consumir sin parsear texto")
	help := fs.Bool("help", false, "")

	if err := fs.Parse(args); err != nil {
		return reports.CaptureInput{}, false, err
	}
	if *help {
		return reports.CaptureInput{}, false, errReportsCaptureHelp
	}
	if *change == "" {
		return reports.CaptureInput{}, false, errors.New("--change es obligatorio")
	}
	if *phase == "" {
		return reports.CaptureInput{}, false, errors.New("--phase es obligatorio")
	}
	if *role == "" {
		return reports.CaptureInput{}, false, errors.New("--role es obligatorio")
	}
	if *reportPath != "" && *reportText != "" {
		return reports.CaptureInput{}, false, errors.New("--report y --report-text son mutuamente excluyentes")
	}

	if *sourceProvider != "" && (*reportPath == "" || *reportText != "" || *contract != reports.ContractRawText || *agent != "" || *session != "" || *turn != "" || *provider != "" || *model != "") {
		return reports.CaptureInput{}, false, errors.New("adapter requiere --report y metadatos del origen sin overrides")
	}
	if *sourceProvider == "" && (*contract == reports.ContractCodexExec132 || *contract == reports.ContractClaude271) {
		return reports.CaptureInput{}, false, errors.New("contrato manual desconocido; usar --source-provider para eventos")
	}
	text := *reportText
	if *reportPath != "" {
		f, err := os.Open(*reportPath)
		if err != nil {
			return reports.CaptureInput{}, false, errors.New("no se pudo abrir --report")
		}
		b, err := io.ReadAll(io.LimitReader(f, reports.MaxContentBytes+1))
		f.Close()
		if len(b) > reports.MaxContentBytes {
			return reports.CaptureInput{}, false, errors.New("--report supera MaxContentBytes")
		}
		if err != nil {
			return reports.CaptureInput{}, false, fmt.Errorf("leer --report %s: %w", *reportPath, err)
		}
		text = string(b)
	}

	source := reports.Source{Contract: *contract, Text: text, Provider: *provider, Model: *model, Agent: *agent, Session: *session, Turn: *turn}
	if *sourceProvider != "" {
		if !strings.Contains("|claude|codex|opencode|", "|"+*sourceProvider+"|") {
			return reports.CaptureInput{}, false, errors.New("proveedor desconocido")
		}
		if *sourceProvider == "claude" {
			source = reports.ResolveProvider(*sourceProvider, *sourceVersion, *sourceSession, []byte(text), *sourceAgent)
		} else {
			source = reports.ResolveProvider(*sourceProvider, *sourceVersion, *sourceSession, []byte(text))
		}
	}
	return reports.CaptureInput{
		Change:     *change,
		Phase:      *phase,
		Round:      *round,
		Role:       *role,
		DispatchID: *dispatchID,
		SHA:        *sha,
		DiffRef:    *diff,
		Source:     source,
	}, *asJSON, nil
}

var errReportsCaptureHelp = errors.New("mostrar el usage")

const reportsUsage = `Usage: active-stack reports <subcomando>

Subcomandos:
  capture   Captura un informe ya resuelto antes del despacho de la fase siguiente.`

const reportsCaptureUsage = `Usage: active-stack reports capture --change <slug> --phase <fase> --round <n> --role <rol> [flags]

Persiste texto manual o eventos de un único despacho.
Claude 2.1.271 (Handback + SubagentStop); Codex exec 0.132.0 JSONL.
OpenCode y otras versiones: UNSUPPORTED. No instala hooks ni invoca proveedores.

Código de salida: 0 COMPLETE, 2 PENDING, 3 UNSUPPORTED, 4 RETAINED, 5 CONFLICT.
1 es un error de infraestructura (flags inválidos, repo ausente, I/O).

Flags:
  --change <slug>       slug del change OPSX
  --phase <fase>        fase OPSX (ej. verify, apply)
  --round <n>           vuelta OPSX (los reintentos no la incrementan)
  --role <rol>          rol del revisor/ejecutor
  --dispatch-id <id>    identificador del despacho puntual
  --sha <sha>           commit revisado
  --diff <rango>        rango de diff revisado (opcional)
  --agent/--session/--turn <v>   identidad disponible del origen (opcional)
  --provider/--model <v>         EVIDENCIA sólo si consta en el origen
  --contract <id>       contrato del Source (default "raw-text")
  --report <ruta>       archivo con el informe íntegro
  --report-text <texto> informe íntegro inline, en vez de --report
  --source-provider <p> claude | codex | opencode (--report contiene eventos)
  --source-version <v> versión exacta del formato
  --source-session <s> sesión esperada registrada antes del despacho
  --source-agent <a>   hija esperada (Claude)
  --json                salida para consumir sin parsear texto`
