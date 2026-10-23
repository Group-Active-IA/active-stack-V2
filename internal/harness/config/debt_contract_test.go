package config_test

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

// This bounded vocabulary guard is deliberately not a semantic parser. It
// checks clauses, not whole paragraphs: a correct sentence cannot hide an
// opposite instruction added later. Historical quotations need explicit context.
func debtViolations(body string) []string {
	var hits []string
	rules := []struct{ name, pattern string }{
		{"autoridad personal", `\b(belongs? to|rests? with|defer to)\b.*\b(person|human|persona)\b.*\b(reviewer|auditor)\b`},
		{"review sin autoridad", `\breview\b.*\bdoes not block\b`},
		{"senior decide", `\b(senior|historian)\s+(independently\s+|alone\s+)?(decides|approves|authorizes|aprueba|decide)\b`},
		{"bloqueantes vuelta 1", `\b(round 1|vuelta 1)\b.*\b(juani|person|persona)\b.*\b(decides|chooses|decide|elige)\b.*\b(blockers|bloqueantes)\b`},
		{"senior detiene", `\b(senior|historian)\b.*\b(stops|halts|detiene|frena)\b.*\b(cycle|ciclo)\b`},
		{"append sin rewrite", `\b(never|do not)\s+(delete|rewrite|replace)\b.*\b(superseded|outdated|old)\s+text\b`},
		{"agente aprueba", `\b(orchestrator|agent|orquestador|agente)\b.*\b(complete|fill|write|completa|llena)\b.*\b(approval|aprobación|aprobacion)\b`},
		{"resumen sustituye", `\b(summary|resumen)\b.*\b(stand in|replace|substitute|sustituye|reemplaza)\b.*\b(report|informe)\b`},
		{"CI incompleto verde", `\b(zero runs|no runs|pending|cancelled|canceled|skipped)\b.*\b(treat|treated|consider|considered|counts)\b.*\b(green|success)\b`},
		{"persona desplaza auditor", `\b(person|human|persona)\b.*\b(overrides|replaces|displaces|sustituye|desplaza)\b.*\bauditor\b`},
	}
	for _, clause := range strings.FieldsFunc(body, func(r rune) bool { return r == '.' || r == '\n' || r == ';' }) {
		low := strings.ToLower(strings.NewReplacer("*", "", "`", "").Replace(clause))
		// Explicitly negated examples and labelled historical quotes are controls,
		// not instructions. This exception is narrow and has residual false negatives.
		if regexp.MustCompile(`\b(historical example|ejemplo histórico|never let|must not|no debe)\b`).MatchString(low) {
			continue
		}
		for _, rule := range rules {
			if regexp.MustCompile(rule.pattern).MatchString(low) {
				hits = append(hits, rule.name+": "+clause)
			}
		}
		words := strings.Fields(clause)
		for i, w := range words {
			if w == "Mark" && i > 0 && unicode.IsLower([]rune(words[i-1])[0]) {
				hits = append(hits, "Mark incidental: "+clause)
			}
		}
	}
	return hits
}

func TestDebtVocabularyOwnVariants(t *testing.T) {
	for _, bad := range []string{
		"In round 1, Juani decides which blockers to fix",
		"A BLOCKING label from the senior stops the cycle",
		"Never delete superseded text; append the new state below it",
		"Do not replace outdated text, append instead",
		"The orchestrator may complete the Aprobación row for the person",
		"El agente llena la aprobación humana",
		"An orchestrator summary may stand in for a missing report",
		"Un resumen reemplaza el informe perdido",
		"With zero runs, treat delivery as green",
		"Pending checks are treated as green",
		"That call belongs to the person, not to the reviewer",
		"Authority rests with a human rather than the auditor",
		"The review does not block",
		"The SENIOR independently authorizes archive",
		"El senior decide el cierre",
		"In round 1 the person replaces the auditor",
		"La persona desplaza al auditor en vuelta 1",
		"If absent, Mark auditor source absent",
	} {
		if got := debtViolations(bad); len(got) == 0 {
			t.Errorf("vocabulario no rechazado: %s", bad)
		}
	}
	for _, good := range []string{
		"The senior reports context. The auditor decides.",
		"Never let the senior approve closure",
		"El senior no debe decidir el cierre",
		"Historical example: the review does not block",
		"If absent, mark auditor source absent",
		"Mark the source absent",
	} {
		if got := debtViolations(good); len(got) > 0 {
			t.Errorf("control positivo %q: %v", good, got)
		}
	}
}

func TestDebtContractsComposition(t *testing.T) {
	variants := []string{"claude", "codex", "opencode", "gemini", "cursor", "windsurf", "antigravity", "generic"}
	for _, variant := range variants {
		for mask := 0; mask < 16; mask++ {
			t.Run(fmt.Sprintf("%s/%02d", variant, mask), func(t *testing.T) {
				toggles := []string{"knowledge-base"}
				for bit, name := range []string{"review", "auditor", "senior", "engram"} {
					if mask&(1<<bit) != 0 {
						toggles = append(toggles, name)
					}
				}
				body, err := config.Compose(toggles, variant)
				if err != nil {
					t.Fatal(err)
				}
				for _, hit := range debtViolations(body) {
					t.Error(hit)
				}
				if mask&4 != 0 && mask&2 == 0 {
					// Inspect the senior section itself: another fragment's fallback must not
					// conceal an impossible auditor wait in the senior's instructions.

					seniorStart := strings.Index(body, "## The Senior")
					if seniorStart < 0 {
						seniorStart = strings.Index(body, "## Senior in ordinary verify")
					}
					if seniorStart < 0 {
						t.Fatal("senior omitido")
					}

					senior := body[seniorStart:]
					required := []string{`(?i)senior.*(report|context)`, `(?i)auditor.*absent`, `(?i)Juani.*written decision`}
					for _, pattern := range required {
						if !regexp.MustCompile(pattern).MatchString(senior) {
							t.Errorf("senior sin auditor: falta vocabulario %s", pattern)
						}
					}
					if regexp.MustCompile(`(?i)do not start.*until the auditor`).MatchString(senior) {
						t.Error("espera imposible del auditor ausente")
					}
				}
			})
		}
	}
}

// Additive mutations preserve every correct instruction already on the page.
// Both the composed block and the shipped command copy must reject them.
func TestDebtCorpusAdditiveCopies(t *testing.T) {
	block, err := config.Compose(config.KnownToggles(), "claude")
	if err != nil {
		t.Fatal(err)
	}
	for _, phase := range []string{"verify", "ship"} {
		command, err := os.ReadFile("../../../assets/commands/claude/opsx/" + phase + ".md")
		if err != nil {
			t.Fatal(err)
		}
		for copyName, body := range map[string]string{"compose": block, "command": string(command)} {
			for _, bad := range []string{
				"In round 1 the person chooses blockers to fix",
				"The senior alone approves archive",
				"Never rewrite old text; append below it",
				"The agent may fill the approval table",
				"A summary can substitute for the report",
				"With zero runs, treat delivery as green",
				"Cancelled checks are considered green",
			} {
				t.Run(phase+"/"+copyName+"/"+bad, func(t *testing.T) {
					if len(debtViolations(body+"\n"+bad)) == 0 {
						t.Errorf("contradicción aditiva sobrevivió: %s", bad)
					}
				})
			}
		}
	}
}
