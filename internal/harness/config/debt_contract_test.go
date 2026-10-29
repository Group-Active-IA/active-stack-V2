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
		{"bloqueantes vuelta 1", `\b(round 1|vuelta 1)\b.*\b(juani|person|persona|human owner)\b.*\b(decides|chooses|decide|elige)\b.*\b(blockers|bloqueantes)\b`},
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

// authorityViolations checks vocabulary families within individual clauses.
// Report subjects and passive constructions count as authority too.
func authorityViolations(body string, absent bool) []string {
	var hits []string
	cache := map[string]*regexp.Regexp{}
	for _, clause := range strings.FieldsFunc(strings.ToLower(strings.NewReplacer("*", "", "`", "", "\n", " ").Replace(body)), func(r rune) bool { return r == '.' || r == ';' }) {
		clause = strings.NewReplacer("closed changes", "past changes", "changes that closed before", "past changes", "openspec/changes/archive/", "past/").Replace(clause)
		match := func(p string) bool {
			re := cache[p]
			if re == nil {
				re = regexp.MustCompile(p)
				cache[p] = re
			}
			return re.MatchString(clause)
		}
		negated := match(`\b(no closing authority|not (?:a |as a )?parallel (?:closing authority|wait)|(?:do|does|must|may|shall|should) not (?:invent|fill|complete|create|authorize|approve|close|archive)|never (?:invent|fill|complete|authorize|approve|close|archive)|cannot|can't)\b`)
		if !negated {
			authority := `\b(close[sd]?|closing|closure|approv[a-z]*|aprobado|archiv[a-z]*|verdict|sign.off|green.light|final decision|final say|go.ahead|finish[a-z]*|authoriz[a-z]*|cierre|aprueb[a-z]*)\b`
			if match(`\b(senior|historian|its report)\b`) && match(authority) {
				hits = append(hits, "clase a: senior con autoridad: "+clause)
			}
			if match(`\b(breaker|rompedor)\b`) && match(authority) {
				hits = append(hits, "clase c: rompedor con autoridad: "+clause)
			}
			if match(`\b(auditor|its)\b`) && match(`\b(list[a-z]*|lessons?|debts?|lecciones|deuda)\b`) && match(`\b(fill[a-z]*|complet[a-z]*|invent[a-z]*|creat[a-z]*|generat[a-z]*|reconstruct[a-z]*|draft[a-z]*|synthesi[a-z]*|populat[a-z]*|derive[a-z]*|supply|supplies|fabricat[a-z]*|produc[a-z]*|assembl[a-z]*|author[a-z]*)\b`) {
				hits = append(hits, "clase d: listas del auditor inventadas: "+clause)
			}
		}
		if absent && !match(`\b(with auditor (active|enabled|present)|if (the )?auditor is (active|enabled|present)|when (the )?auditor is (active|enabled|present))\b`) {
			// Negation applies to the wait, not merely to the auditor's presence.
			waitNegated := match(`\b(do not|never|must not|cannot|no need to) (wait|keep|maintain|require|hold|remain)\b`)
			if !waitNegated && match(`\bauditor\b`) && match(`\b(wait[a-z]*|await[a-z]*|until|pending|requires?|only|closes? with|remain[a-z]* open|keep[a-z]*.*open|hold[a-z]*.*open)\b`) && match(`\b(verdict|decision|aprobado|approval|returns?|respond[a-z]*|sign.off)\b`) {
				hits = append(hits, "clase b: espera de auditor ausente: "+clause)
			}
		}
	}
	return hits
}

func TestAuthorityVocabularyClasses(t *testing.T) {
	cases := []struct{ class, text string }{
		{"a", "Its report is the closing verdict whenever no auditor is configured."},
		{"a", "Closure proceeds on the senior's report alone."},
		{"a", "Archive is approved by the historian."},
		{"a", "If the auditor is not available, the senior approves closure."},
		{"b", "Keep the round open until the auditor returns its verdict."},
		{"b", "Archiving awaits the absent auditor's approval."},
		{"c", "The breaker's report is enough to close the change."},
		{"c", "Approval for archive rests with the breaker."},
		{"d", "Fill its lists from your own reading of the diff."},
		{"d", "Reconstruct auditor lists from available evidence."},
	}
	for _, c := range cases {
		t.Run(c.class+"/"+c.text, func(t *testing.T) {
			hits := authorityViolations(c.text, true)
			if len(hits) == 0 || !strings.Contains(strings.Join(hits, " "), "clase "+c.class+":") {
				t.Fatalf("clase %s sobrevivió: %s (%v)", c.class, c.text, hits)
			}
		})
	}
	for _, good := range []string{"The senior reports context", "The breaker cannot authorize closure", "Do not invent auditor lists", "Do not wait for a verdict from an absent auditor", "With auditor active, only APROBADO from the auditor enables archive"} {
		if hits := authorityViolations(good, true); len(hits) > 0 {
			t.Errorf("control positivo: %v", hits)
		}
	}
}

func TestDebtVocabularyOwnVariants(t *testing.T) {
	for _, bad := range []string{
		"In round 1, the human owner decides which blockers to fix",
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
				for _, hit := range append(debtViolations(body), authorityViolations(body, mask&2 == 0)...) {
					t.Error(hit)
				}
				if mask&3 == 3 && !regexp.MustCompile(`(?i)with auditor active, only APROBADO from the auditor.*enables archive`).MatchString(body) {
					t.Error("auditor activo: falta requisito de APROBADO para archive")
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
					required := []string{`(?i)senior.*(report|context)`, `(?i)auditor.*absent`, `(?i)the human owner.*written decision`}
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
	block, err := config.Compose(legacyToggles(), "claude")
	if err != nil {
		t.Fatal(err)
	}
	for _, phase := range []string{"verify", "ship"} {
		command, err := os.ReadFile("../../../assets/commands/claude/opsx/" + phase + ".md")
		if err != nil {
			t.Fatal(err)
		}
		for copyName, body := range map[string]string{"compose": block, "command": string(command)} {
			if hits := debtViolations(body); len(hits) != 0 {
				t.Fatalf("original no limpio: %v", hits)
			}
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
					hits := debtViolations(body + "\n" + bad)
					attributable := false
					for _, hit := range hits {
						if strings.Contains(bad, strings.TrimSpace(strings.SplitN(hit, ": ", 2)[1])) {
							attributable = true
						}
					}
					if !attributable {
						t.Errorf("contradicción aditiva sobrevivió: %s", bad)
					}
				})
			}
		}
	}
}

// Every replacement starts from an immutable clean fragment, never a previous mutant.
func TestAuthorityFragmentMutations(t *testing.T) {
	cases := []struct{ name, class, file, old, replacement string }{
		{"M1", "a", "senior-without-auditor.md", "It has no closing authority.", "Its report is the closing verdict whenever no auditor is configured."},
		{"M2", "b", "senior-without-auditor.md", "Do not wait for a verdict from an absent auditor, impersonate it or authorize an automatic return to apply.", "Keep the round open until the auditor returns its verdict; meanwhile send blockers back to apply."},
		{"M3", "c", "review.md", "the breaker cannot authorize closure.", "the breaker's report is enough to close the change."},
		{"M5", "d", "knowledge-base.md", "do not invent its lists", "fill its lists from your own reading of the diff"},
		{"M6", "a", "review.md", "closure stops for the human owner's written decision on degraded scope", "closure proceeds on the senior's report alone"},
		{"V1", "a", "senior-without-auditor.md", "It has no closing authority.", "The historian grants approval for archiving."},
		{"V2", "b", "senior-without-auditor.md", "Do not wait for a verdict from an absent auditor, impersonate it or authorize an automatic return to apply.", "Archiving awaits the absent auditor's approval."},
		{"V3", "c", "review.md", "the breaker cannot authorize closure.", "Approval for archive rests with the breaker."},
		{"V4", "d", "knowledge-base.md", "do not invent its lists", "reconstruct auditor lists from available evidence"},
	}
	for _, c := range cases {
		t.Run(c.name+"/clase "+c.class, func(t *testing.T) {
			original, err := os.ReadFile("assets/" + c.file)
			if err != nil {
				t.Fatal(err)
			}
			clean := string(original)
			if hits := authorityViolations(clean, true); len(hits) > 0 {
				t.Fatalf("original no limpio clase %s: %v", c.class, hits)
			}
			if strings.Count(clean, c.old) != 1 {
				t.Fatalf("clase %s: punto de mutación ambiguo", c.class)
			}
			mutant := strings.Replace(clean, c.old, c.replacement, 1)
			hits := authorityViolations(mutant, true)
			if !strings.Contains(strings.Join(hits, "\n"), "clase "+c.class+":") {
				t.Fatalf("clase %s: %s sobrevivió: %v", c.class, c.name, hits)
			}
		})
	}
}
