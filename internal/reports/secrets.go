package reports

import "regexp"

// secretPattern es un patrón propio de escaneo. gitleaks no está en el PATH
// de este entorno (verificado: `which gitleaks` no lo encuentra) — el
// propose permite "reusar gitleaks si el repo destino lo tiene, más patrones
// propios"; sin gitleaks disponible, estos patrones son el escaneo completo,
// no un complemento. DEUDA: son más angostos que una herramienta dedicada —
// no detectan secretos de alta entropía sin prefijo reconocible.
type secretPattern struct {
	kind string
	re   *regexp.Regexp
}

var secretPatterns = []secretPattern{
	{"github-token", regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{20,}`)},
	{"slack-token", regexp.MustCompile(`xox[baprs]-[A-Za-z0-9-]{10,}`)},
	{"aws-access-key-id", regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	{"private-key-block", regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`)},
	{"generic-secret-assignment", regexp.MustCompile(`(?i)(api[_-]?key|secret|token|password|passwd)\s*[:=]\s*['"]?[A-Za-z0-9/_\-\.]{12,}`)},
	{"email", regexp.MustCompile(`[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`)},
	{"phone", regexp.MustCompile(`\+?\d[\d\-\s]{8,}\d`)},
}

// scanSecrets devuelve los TIPOS de secreto que aparecen en text, en orden
// determinista, sin duplicados y sin el valor encontrado — el contrato
// prohíbe "volcar el valor detectado" en cualquier salida, y eso incluye el
// resultado de este escaneo.
func scanSecrets(text string) []string {
	var kinds []string
	for _, p := range secretPatterns {
		if p.re.MatchString(text) {
			kinds = append(kinds, p.kind)
		}
	}
	return kinds
}
