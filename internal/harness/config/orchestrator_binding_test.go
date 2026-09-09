package config

import (
	"regexp"
	"strings"
	"testing"
)

// variantes is every agent flavour of the base block. Each one is a separate
// file, so a fix applied to one says nothing about the other seven.
var variantes = []string{
	"antigravity", "claude", "codex", "cursor",
	"gemini", "generic", "opencode", "windsurf",
}

// normalizarMarkdown strips emphasis, backticks and repeated whitespace so an
// assert compares CONTENT and not formatting.
//
// Without it these checks are evadable by accident: wrapping a phrase in
// **bold** or `code` breaks a raw substring match, so a positive assert goes
// red on a harmless edit and — worse — an anti-assert goes green while the
// forbidden sentence is still sitting there in italics.
func normalizarMarkdown(s string) string {
	s = strings.NewReplacer("*", "", "`", "", "_", "").Replace(s)
	return strings.ToLower(regexp.MustCompile(`\s+`).ReplaceAllString(s, " "))
}

// TestBloqueSeDirigeALaSesionPrincipal pins WHO the orchestrator block is
// addressed to, in both directions.
//
// The block used to open with "Bind this to the dedicated `sdd-orchestrator`
// agent only. Do NOT apply it to executor agents." That sentence had no
// reader: no installer ever creates an `sdd-orchestrator` agent file (verified
// on 2026-09-09 against a full install — the eight sdd-* phase agents exist,
// the orchestrator does not), and the session that DOES orchestrate read line
// 3 and correctly concluded the block was not addressed to it.
//
// The result was the worst of both: ~1500 lines of instructions that every
// agent read and none applied. Observed live — a full foundation run on a real
// project produced good artifacts and convened zero reviewers, zero
// specialists and zero auditors, because nothing considered itself the
// addressee.
//
// The block is written for the main conversation and cannot be otherwise: it
// asks the orchestrator to put a finding in front of the person and wait for
// the answer, and a sub-agent returns its report and terminates. This block
// says so itself, further down.
//
// The datum that discriminates: the composed text must NAME the main
// conversation as the orchestrator (positive), and must NOT bind itself to a
// dedicated agent instead (anti-assert). One direction alone is not enough —
// a block could add the positive sentence and keep the old exclusion right
// under it, and each half on its own would still pass.
func TestBloqueSeDirigeALaSesionPrincipal(t *testing.T) {
	// Frases que DEBEN estar: el bloque nombra a la conversacion principal.
	requeridas := []string{
		"the session that talks to the user",
		"the main conversation is the orchestrator",
	}
	// Frases PROHIBIDAS: atarse a un agente dedicado excluye a quien orquesta.
	prohibidas := []string{
		"bind this to the dedicated",
		"bind this to the sdd-orchestrator",
	}

	for _, variante := range variantes {
		t.Run(variante, func(t *testing.T) {
			bloque, err := Compose(KnownToggles(), variante)
			if err != nil {
				t.Fatalf("Compose(%q) error = %v", variante, err)
			}
			texto := normalizarMarkdown(bloque)

			for _, frase := range requeridas {
				if !strings.Contains(texto, normalizarMarkdown(frase)) {
					t.Errorf("el bloque de %q no dice a quien se dirige: falta %q", variante, frase)
				}
			}
			for _, frase := range prohibidas {
				if strings.Contains(texto, normalizarMarkdown(frase)) {
					t.Errorf("el bloque de %q vuelve a atarse a un agente dedicado: contiene %q", variante, frase)
				}
			}
		})
	}
}
