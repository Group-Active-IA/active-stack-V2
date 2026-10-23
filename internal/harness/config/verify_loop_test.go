package config_test

import (
	"github.com/Group-Active-IA/active-stack/internal/harness/config"
	"strings"
	"testing"
)

func TestVerifyContratoBucle(t *testing.T) {
	body := seniorFlat(reviewFragment(t))
	for _, forbidden := range []string{"if blockers remain, start round 4 automatically.", "launch the auditor in parallel with the senior."} {
		if strings.Contains(body, forbidden) {
			t.Errorf("secuencia contradictoria: %s", forbidden)
		}
	}
	if !strings.Contains(body, "wait for the senior report before launching the auditor.") {
		t.Error("falta espera obligatoria; no se permite: Launch the auditor in parallel with the senior.")
	}
	for _, frase := range []string{
		"conformity -> breaker -> senior (if active) -> auditor",
		"only aprobado from the auditor, with conformity satisfied and complete reports, enables archive",
		"unresolved critical prevents aprobado",
		"## bucle de verify", "before convening reviewers", "round 1, 2 or 3", "in progress or closed", "sha or reviewed diff", "report paths, blockers and decision",
		"resume reads tasks.md, not engram as authority", "failed capture does not increment or reset the counter",
		"blockers in round 1 or 2 return only the auditor's list to apply", "round 3 closes the loop", "no automatic fourth round", "without a written answer nothing advances", "a new rule requires human signature before construction (l-029)",
		"correctness: record boundary cases, command, output and result", "wiring: trace consumers and connections, command, output and result", "comments/documentation: compare contracts with behavior, sources and result",
		"a failing test or a surviving mutant", "outside the diff", "first change that uses the piece", "count tests themselves", "read numbered lessons", "at the keyboard", "explicit debt and a limitation for the auditor",
		"announce missing toggles at verify start and record degradation in tasks.md", "auditor absent with review active", "no automatic approval or return to apply", "tdd absent", "do not claim strict tdd", "senior absent", "do not invent historical context", "knowledge-base absent", "destination pending", "review absent", "no breaker or three-round loop", "knowledge-base active with auditor absent", "do not invent auditor lessons or debt", "missing required report stops closure", "never replace it with a summary",
		"informes/", "leeme.md", "model/provider", "provenance", "verify-round-{n}.md references the originals", "reviewers receive no write", "conformity, breaker, senior and auditor", "full return or textual transcript",
	} {
		t.Run(frase, func(t *testing.T) {
			if !strings.Contains(body, frase) {
				t.Errorf("falta contrato: %s", frase)
			}
		})
	}
}

func TestAuditorDecideSinAutoridadParalela(t *testing.T) {
	sequence := seniorFlat(reviewFragment(t))
	if strings.Contains(sequence, "launch the auditor in parallel with the senior.") || !strings.Contains(sequence, "wait for the senior report before launching the auditor.") {
		t.Error("secuencia contradictoria: Launch the auditor in parallel with the senior.")
	}

	b := seniorFlat(auditorFragment(t))
	if !strings.Contains(b, "missing required reports prevent approval.") {
		t.Error("falta negación obligatoria: Missing required reports prevent approval.; permitir aprobación sin informes contradice el contrato")
	}
	for _, frase := range []string{"aprobado", "bloqueantes", "deuda", "lecciones", "none", "unresolved critical prevents aprobado", "classify disagreements with the senior", "dated decision", "count tests independently"} {
		if !strings.Contains(b, frase) {
			t.Errorf("falta decisión: %s", frase)
		}
	}
	s := seniorFlat(seniorFragment(t))
	for _, frase := range []string{"after the breaker returns and before the auditor", "the auditor classifies and decides", "no parallel human wait"} {
		if !strings.Contains(s, frase) {
			t.Errorf("falta secuencia: %s", frase)
		}
	}
	all, err := config.Compose(config.KnownToggles(), "claude")
	if err != nil {
		t.Fatal(err)
	}
	s = seniorFlat(all)
	for _, frase := range []string{"the two run together", "when the senior returns a finding labelled blocking, do not start", "the senior labels, you hold"} {
		if strings.Contains(s, frase) {
			t.Errorf("autoridad vieja: %s", frase)
		}
	}
}

func TestVerifyDegradacionesComponibles(t *testing.T) {
	for _, ausente := range []string{"", "review", "auditor", "tdd", "senior", "knowledge-base"} {
		t.Run(ausente, func(t *testing.T) {
			toggles := []string{}
			for _, x := range []string{"review", "auditor", "tdd", "senior", "knowledge-base"} {
				if x != ausente {
					toggles = append(toggles, x)
				}
			}
			b, e := config.Compose(toggles, "claude")
			if e != nil {
				t.Fatal(e)
			}
			if ausente != "review" && !strings.Contains(b, "Announce missing toggles at verify start") {
				t.Fatal("falta degradación visible")
			}
			if ausente == "review" && strings.Contains(b, "### The loop") {
				t.Fatal("bucle filtrado")
			}
		})
	}
}
