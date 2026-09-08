package config_test

// ──────────────────────────────────────────────────────────────────────────────
// El catálogo y Compose tienen que hablar del mismo set de toggles
//
// El catálogo (`internal/catalog/harnesses.yaml`) es lo que la TUI le OFRECE al
// usuario. `Compose` es lo que sabe armar. Nada los ataba: un toggle ofrecido en
// el YAML y ausente de `knownToggles` se elige sin fricción y explota recién en
// el install, con "unknown toggle", después de que el usuario ya dijo que sí.
//
// Es el mismo agujero para los cuatro toggles que ya existían. El guardia entra
// ahora porque `review` es el quinto y ya había pasado cuatro veces sin red.
// ──────────────────────────────────────────────────────────────────────────────

import (
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/catalog"
	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

// orchestratorToggles returns the toggle list the catalog advertises for the
// sdd-orchestrator harness.
func orchestratorToggles(t *testing.T) []string {
	t.Helper()
	c, err := catalog.Load()
	if err != nil {
		t.Fatalf("catalog.Load: %v", err)
	}
	for _, h := range c.Harnesses {
		if h.ID == "sdd-orchestrator" {
			if len(h.Toggles) == 0 {
				t.Fatal("the sdd-orchestrator harness advertises no toggles at all")
			}
			return h.Toggles
		}
	}
	t.Fatal("the catalog has no sdd-orchestrator harness")
	return nil
}

// TestCatalog_EveryOfferedToggleComposes asserts Compose accepts everything the
// catalog offers, one by one and all together.
//
// One by one AND all together on purpose: composing the full list only proves
// the set is valid, and the error names just the first offender.
func TestCatalog_EveryOfferedToggleComposes(t *testing.T) {
	toggles := orchestratorToggles(t)
	for _, tg := range toggles {
		if _, err := config.Compose([]string{tg}, "claude"); err != nil {
			t.Errorf("the catalog offers toggle %q and Compose rejects it: %v — the user "+
				"picks it in the TUI and the install fails afterwards", tg, err)
		}
	}
	if _, err := config.Compose(toggles, "claude"); err != nil {
		t.Errorf("Compose rejects the catalog's full toggle list: %v", err)
	}
}

// TestCatalog_OffersTheReviewToggle asserts the adversarial review is reachable.
//
// A fragment that composes but that the catalog never advertises is dead code
// with tests: nothing in the product can turn it on.
func TestCatalog_OffersTheReviewToggle(t *testing.T) {
	toggles := orchestratorToggles(t)
	for _, tg := range toggles {
		if tg == "review" {
			return
		}
	}
	t.Errorf("the catalog does not offer the `review` toggle (offers: %s): the "+
		"fragment composes and no user can ask for it", strings.Join(toggles, ", "))
}

// TestCatalog_OpenspecDescriptionNamesVerify asserts the sixth place that
// enumerates the cycle also names the phase.
//
// The other five are inside the orchestrator block, which only the agent reads.
// This one is the `openspec` harness's long_description: it is what the USER
// reads in the TUI while deciding whether to install the CLI. Teaching the agent
// a five-phase cycle while telling the person it has four is the same defect
// pointed at a different reader.
func TestCatalog_OpenspecDescriptionNamesVerify(t *testing.T) {
	c, err := catalog.Load()
	if err != nil {
		t.Fatalf("catalog.Load: %v", err)
	}
	for _, h := range c.Harnesses {
		if h.ID != "openspec" {
			continue
		}
		for _, lang := range []string{"es", "en"} {
			text := h.LongDescription[lang]
			low := strings.ToLower(text)
			if low == "" {
				continue
			}
			// Both wordings are the localised cycle, so check the localised verb.
			want := "verificar"
			if lang == "en" {
				want = "verify"
			}
			if !strings.Contains(low, want) {
				t.Errorf("the openspec harness's %s long_description describes the cycle "+
					"without %q: the user is told there are four phases and the CLI "+
					"installs five", lang, want)
			}
		}
		return
	}
	t.Fatal("the catalog has no openspec harness")
}
