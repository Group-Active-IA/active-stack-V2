package config_test

// ──────────────────────────────────────────────────────────────────────────────
// El auditor: alguien que revisa a los que revisan
//
// Hoy nadie revisa al revisor. Si un sub-agente reporta "revisé y está todo
// bien", no hay forma de saber si revisó. Llega el veredicto, no lo que hizo
// para llegar a él.
//
// Es el mismo problema de evidencia contra juicio, una capa más arriba, y por
// eso el auditor va DESPUÉS del reporte tipado y no antes: sin declaraciones
// tipadas no tiene qué auditar, sólo prosa sobre la que opinar.
//
// El 2026-09-08 dos revisiones independientes sobre overlays.go encontraron el
// mismo bug. Eso fue señal de que era real. Si hubieran coincidido palabra por
// palabra habría sido lo contrario: una sola opinión repetida. Nadie estaba
// mirando eso, y es exactamente lo que este agente mira.
// ──────────────────────────────────────────────────────────────────────────────

import (
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

func auditorFragment(t *testing.T) string {
	t.Helper()
	body, err := config.Compose([]string{"auditor"}, "claude")
	if err != nil {
		t.Fatalf("Compose([auditor]): %v", err)
	}
	return strings.ToLower(body)
}

func TestCompose_AuditorToggle_IsKnown(t *testing.T) {
	for _, v := range evidenceVariants {
		if _, err := config.Compose([]string{"auditor"}, v); err != nil {
			t.Errorf("Compose([auditor], %q): %v", v, err)
		}
	}
}

func TestCompose_AuditorToggle_Subtractive(t *testing.T) {
	off, err := config.Compose([]string{"delegation"}, "claude")
	if err != nil {
		t.Fatalf("Compose without auditor: %v", err)
	}
	if strings.Contains(off, "## The Auditor") {
		t.Error("the auditor fragment leaks into a block that did not ask for it")
	}
}

// TestCompose_AuditorToggle_ChecksTheEvidenceClaims asserts the auditor verifies
// that what was filed as evidence actually carries evidence.
//
// This is its cheapest and highest-yield check: an agent claiming it ran
// something, with no command and no output attached, is the failure the typed
// report was built to expose.
func TestCompose_AuditorToggle_ChecksTheEvidenceClaims(t *testing.T) {
	body := auditorFragment(t)
	required := map[string]string{
		"without the command": "evidence with no command behind it is a claim",
		"stale":               "output from an earlier run describes a tree that no longer exists",
		"promoted":            "an assumption filed as evidence is the whole failure mode",
	}
	for phrase, why := range required {
		if !strings.Contains(body, phrase) {
			t.Errorf("the auditor fragment is missing %q: %s", phrase, why)
		}
	}
}

// TestCompose_AuditorToggle_DetectsSuspiciousAgreement asserts it looks at
// whether independent reviewers really were independent.
//
// Three reviewers saying the same thing in the same words are not three
// opinions: they are one that got repeated. Only somebody reading all three can
// see it.
func TestCompose_AuditorToggle_DetectsSuspiciousAgreement(t *testing.T) {
	body := auditorFragment(t)
	for _, phrase := range []string{"same words", "one opinion"} {
		if !strings.Contains(body, phrase) {
			t.Errorf("the auditor fragment never checks for suspicious agreement (%q): "+
				"convergence is only evidence when the reviewers were actually independent", phrase)
		}
	}
}

// TestCompose_AuditorToggle_ReportsToTheUser asserts the auditor answers to the
// person, not to the orchestrator.
//
// If the orchestrator could filter what the auditor says, the auditor is the
// boss picking which complaints reach you. This is the closest thing to a
// separation of powers available without separate OS identities.
func TestCompose_AuditorToggle_ReportsToTheUser(t *testing.T) {
	body := auditorFragment(t)
	for _, phrase := range []string{"reports to the user", "must not summarise"} {
		if !strings.Contains(body, phrase) {
			t.Errorf("the auditor fragment is missing %q: an auditor the orchestrator "+
				"can edit is not an auditor", phrase)
		}
	}
}

// TestCompose_AuditorToggle_HasASilenceThreshold asserts it is allowed to say
// nothing.
//
// An auditor that always finds something gets ignored by the third change, and
// then it is not auditing anything.
func TestCompose_AuditorToggle_HasASilenceThreshold(t *testing.T) {
	body := auditorFragment(t)
	for _, phrase := range []string{"aprobado", "does not know"} {
		if !strings.Contains(body, phrase) {
			t.Errorf("the auditor fragment is missing %q: without a silence threshold "+
				"it becomes noise and gets switched off", phrase)
		}
	}
}
