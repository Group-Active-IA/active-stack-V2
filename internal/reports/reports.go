// Package reports captura, nunca interpreta, informes de revisores antes del
// despacho de la fase siguiente del ciclo OPSX.
//
// POR QUÉ EXISTE (propose 2026-10-23-fase-2-pasos-1-4.md §2a, "Captura de
// informes: núcleo"). El ciclo adversarial (verify) lanza varios pases
// independientes; si nadie registra el informe de cada uno ANTES de que el
// orquestador lo resuma, lo que sobrevive es el resumen del orquestador, no
// el original. Este paquete persiste el original — crudo, texto derivado,
// manifiesto e índice — y nunca lo sustituye por un resumen.
//
// ResolveProvider resuelve contratos versionados de Claude y Codex; OpenCode
// permanece UNSUPPORTED. El núcleo Capture valida y persiste el Source resuelto.
// Identidad (agente/sesión/turno) nunca sale del TEXTO del informe: si el
// adapter no la leyó del origen (hook, evento, respuesta de API), el campo
// queda vacío y la asociación es ambigua (ver evaluateStatus).
package reports

import (
	"encoding/json"
	"fmt"
)

// Status es el resultado de capturar un informe. El cero-valor es
// StatusPending a propósito: ningún consumidor que ignore el error puede leer
// un cero-valor como COMPLETO. Nunca existe un estado que signifique "éxito
// genérico" fuera de StatusComplete, y StatusComplete exige artefactos
// persistidos y contrastados — nunca es el cero-valor.
type Status int

const (
	// StatusPending: contrato conocido pero contenido vacío/en blanco, o
	// identidad ambigua (falta DispatchID, SHA o la vuelta no es positiva).
	// Se persiste lo que hay; no es COMPLETO hasta que se resuelva.
	StatusPending Status = iota
	// StatusUnsupported: el contrato no está registrado (adapter sin
	// implementar o versión no fijada). No se persiste nada: no hay nada que
	// el núcleo pueda interpretar todavía.
	StatusUnsupported
	// StatusRetained (RETENIDO): el escaneo de secretos encontró algo. El
	// original NUNCA entra al árbol versionado; queda local, en una ruta
	// ignorada por git, y el índice lo muestra con causa y tipos detectados,
	// nunca con el valor.
	StatusRetained
	// StatusConflict: misma identidad (mismo Key), bytes distintos a lo ya
	// persistido. El original existente NUNCA se pisa; este intento queda en
	// un directorio de conflicto separado para decisión humana.
	StatusConflict
	// StatusComplete: artefactos persistidos, identidad completa, sin
	// secretos, sin conflicto. El único estado que approval/auditor pueden
	// tratar como informe disponible para contraste.
	StatusComplete
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "PENDING"
	case StatusUnsupported:
		return "UNSUPPORTED"
	case StatusRetained:
		return "RETAINED"
	case StatusConflict:
		return "CONFLICT"
	case StatusComplete:
		return "COMPLETE"
	default:
		return "PENDING"
	}
}

// MarshalJSON serializa Status como su nombre ("COMPLETE", "PENDING", ...),
// nunca como el entero crudo — la salida --json de `reports capture` tiene
// que ser tan explícita como la de texto.
func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// ExitCode es el código de proceso que usa `active-stack reports capture`.
// 0 únicamente para COMPLETE: es el único estado que certifica un informe
// disponible para contraste. Nunca hay dos estados distintos con el mismo
// código distinto de éxito/no-éxito genérico — cada código identifica un solo
// Status, para que un script no confunda PENDING con RETAINED.
func (s Status) ExitCode() int {
	switch s {
	case StatusComplete:
		return 0
	case StatusPending:
		return 2
	case StatusUnsupported:
		return 3
	case StatusRetained:
		return 4
	case StatusConflict:
		return 5
	default:
		return 2
	}
}

// Source es lo que un adapter de proveedor (2b Claude, 2c Codex, 2d OpenCode
// — ver ResolveProvider) resuelve antes de llamar a
// Capture. El núcleo no interpreta hooks ni respuestas de API: el adapter ya
// hizo ese trabajo y entrega texto e identidad resueltos.
//
// Contrato de extensión para 2b/2c/2d: cada adapter define su propio
// Contract (p. ej. "claude-hook-v1") y lo registra con RegisterContract antes
// de que Capture lo acepte; un Contract no registrado siempre es
// StatusUnsupported, nunca un intento de interpretarlo a medias.
type Source struct {
	// Raw conserva el evento original separado del texto derivado.
	Raw string
	// Incomplete impide certificar eventos parciales o asociaciones dudosas.
	Incomplete bool
	// Contract identifica la versión de hook/API que resolvió este Source.
	// "" o no registrado => StatusUnsupported.
	Contract string
	// Text es el informe íntegro ya resuelto — el informe, no el evento ni
	// el transcript completo que lo contenía ("bruto limitado al
	// evento/informe necesario").
	Text string
	// Provider y Model son EVIDENCIA sólo si el adapter los leyó del origen
	// (header del hook, campo de la respuesta de API). Nunca se infieren del
	// texto del modelo: un informe que dice "soy GPT-4" en su prosa no hace
	// que Model valga "GPT-4".
	Provider string
	Model    string
	// Agent, Session y Turn son los identificadores que el adapter pudo
	// resolver de su origen. Vacíos = no disponibles; eso no bloquea la
	// captura (el propose los pide "disponibles", no obligatorios) pero sí
	// cuenta para identidad ambigua junto con DispatchID/SHA.
	Agent   string
	Session string
	Turn    string
}

// CaptureInput son los parámetros de Capture. RepoDir, Change, Phase, Role,
// DispatchID y SHA identifican DÓNDE y QUÉ se está capturando; Source es lo
// que se captura.
type CaptureInput struct {
	// RepoDir es el worktree canónico. Tiene que ser la raíz del worktree
	// (lo que devuelve `git rev-parse --show-toplevel`), no un subdirectorio
	// ni un worktree distinto ("rechazar... worktree ajeno").
	RepoDir string
	// Change es el slug del change OPSX (el directorio bajo
	// openspec/changes/<Change>/), p. ej. "metodologia-fase-2".
	Change string
	// Phase es la fase OPSX que despacha este informe (p. ej. "verify",
	// "apply").
	Phase string
	// Round es la vuelta OPSX. Reintentos del mismo intento NO la
	// incrementan — eso lo decide quien llama a Capture, nunca este paquete.
	// Round <= 0 cuenta como identidad ambigua.
	Round int
	// Role es el rol del revisor/ejecutor (p. ej. "backend", "qa"). Dos
	// roles nunca comparten Key: Role entra en la clave de identidad.
	Role string
	// DispatchID identifica el despacho puntual que generó este informe.
	DispatchID string
	// SHA es el commit revisado. DiffRef es, opcionalmente, el rango de diff
	// revisado (p. ej. "base..head").
	SHA     string
	DiffRef string

	Source Source
}

// CaptureResult es la salida de Capture.
type CaptureResult struct {
	Status Status `json:"status"`
	// Key es la identidad determinista: ver BuildKey.
	Key string `json:"key"`
	// Dir es la ruta (relativa a RepoDir) donde quedaron — o hubieran
	// quedado, para Unsupported — los artefactos.
	Dir string `json:"dir"`
	// ContentHash es sha256 de Source.Text en hex.
	ContentHash string `json:"content_hash,omitempty"`
	// Cause explica Pending/Unsupported/Retained/Conflict. Nunca incluye el
	// valor de un secreto ni el contenido del informe.
	Cause string `json:"cause,omitempty"`
	// SecretKinds son los TIPOS de secreto que el escaneo encontró (p. ej.
	// "github-token"), nunca el valor ni la posición exacta en texto libre.
	SecretKinds []string `json:"secret_kinds,omitempty"`
	// RetainedPath es la ruta local (ignorada por git) donde quedó el
	// original cuando Status es Retained. Vacío en cualquier otro estado.
	RetainedPath string `json:"retained_path,omitempty"`
}

// Capture valida RepoDir/Change/Phase/Role y persiste (o rechaza) el Source
// según el contrato: ver evaluateStatus para la matriz de estados y
// capture.go para la escritura atómica real.
func Capture(in CaptureInput) (CaptureResult, error) {
	if in.RepoDir == "" || in.Change == "" || in.Phase == "" || in.Role == "" {
		return CaptureResult{}, fmt.Errorf("reports.Capture: RepoDir, Change, Phase y Role son obligatorios")
	}
	return capture(in)
}
