package reports

import (
	"fmt"
	"strings"
)

// contractsRegistered son los Contract que Capture acepta. "raw-text" es el
// único contrato que este núcleo resuelve por sí mismo — captura directa de
// texto ya resuelto, sin adapter de proveedor (uso manual, CLI, tests). Los
// adapters 2b/2c/2d (sin implementar) llaman a RegisterContract con su propio
// identificador antes de que Capture los acepte.
var contractsRegistered = map[string]bool{
	ContractRawText: true,
}

// ContractRawText es el contrato base: el llamador ya tiene el texto del
// informe resuelto y no hay adapter de proveedor de por medio.
const ContractRawText = "raw-text"

// RegisterContract habilita un Contract adicional para Capture. Pensado para
// que internal/reports/claude, internal/reports/codex e
// internal/reports/opencode (2b/2c/2d, sin implementar en este paquete) lo
// llamen en su init() con su propio identificador versionado, p. ej.
// "claude-hook-v1". Un Contract no registrado siempre es StatusUnsupported.
func RegisterContract(id string) {
	if id == "" {
		return
	}
	contractsRegistered[id] = true
}

// knownContract dice si Contract está registrado.
func knownContract(id string) bool {
	return contractsRegistered[id]
}

// safeSegment rechaza vacío, separadores de ruta (/ y \), ".." en cualquier
// posición y un segmento que sea exactamente "." — lo mínimo para que Change/
// Phase/Role/round-N nunca puedan escapar del directorio que Capture
// construye con ellos ("rechazar traversal"). No es una lista blanca de
// caracteres: eso quedaría como deuda si en el futuro un change necesita
// unicode o espacios.
func safeSegment(s string) error {
	if s == "" {
		return fmt.Errorf("segmento vacío")
	}
	if s == "." || s == ".." {
		return fmt.Errorf("segmento %q no puede ser \".\" ni \"..\"", s)
	}
	if strings.ContainsAny(s, "/\\") {
		return fmt.Errorf("segmento %q no puede contener separadores de ruta", s)
	}
	if strings.Contains(s, "..") {
		return fmt.Errorf("segmento %q no puede contener \"..\"", s)
	}
	return nil
}

// BuildKey construye la identidad determinista de un informe. El propose
// fija literalmente "opsx/{change}/verify/round-{n}/{role}" para la fase
// verify (RF1-6); generalizamos esa misma forma a cualquier fase porque el
// propose no da una alternativa — DEUDA: sólo la forma de verify está en el
// texto del propose, el resto es extrapolación de este constructor
// (informe de apply, sección JUICIO).
func BuildKey(change, phase string, round int, role string) (string, error) {
	for _, seg := range []string{change, phase, role} {
		if err := safeSegment(seg); err != nil {
			return "", err
		}
	}
	if round <= 0 {
		return "", fmt.Errorf("round tiene que ser positivo, recibí %d", round)
	}
	return formatKeyUnchecked(change, phase, round, role), nil
}

// formatKeyUnchecked da formato a la clave sin validar nada — lo llaman
// BuildKey (que sí valida y exige round positivo) y capture() (que necesita
// una clave para informar incluso en el caso ambiguo, round <= 0 incluido,
// porque ahí el Status es Pending, no un error).
func formatKeyUnchecked(change, phase string, round int, role string) string {
	return fmt.Sprintf("opsx/%s/%s/round-%d/%s", change, phase, round, role)
}

// dirSegments devuelve los segmentos de directorio (relativos a
// openspec/changes/<change>/informes/) para la identidad dada, validando lo
// mismo que BuildKey. round puede ser <= 0 acá (identidad ambigua, no
// traversal) — ese caso se resuelve a round-0 y evaluateStatus lo marca
// Pending, nunca Unsupported ni error.
func dirSegments(change, phase string, round int, role string) ([]string, error) {
	for _, seg := range []string{change, phase, role} {
		if err := safeSegment(seg); err != nil {
			return nil, err
		}
	}
	r := round
	if r < 0 {
		r = 0
	}
	return []string{"openspec", "changes", change, "informes", phase, fmt.Sprintf("round-%d", r), role}, nil
}
