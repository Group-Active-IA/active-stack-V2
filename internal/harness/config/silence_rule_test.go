package config

import (
	"strings"
	"testing"
)

// sinToggle returns legacyToggles() minus one, so a test can compose the block
// with a specific fragment switched off.
func sinToggle(excluido string) []string {
	var out []string
	for _, t := range legacyToggles() {
		if t != excluido {
			out = append(out, t)
		}
	}
	return out
}

// TestReglaDelSilencioLlegaSiempre pins the rule that an un-convened reviewer
// must be reported, and pins it to a part of the block that no toggle removes.
//
// Why it has to survive every combination: the four fragments that convene
// anyone (roles, auditor, review, senior) are each optional, and so is the
// base's own "Delegation Rules" section — the `delegation` toggle deletes it
// outright. A rule about failing to convene, parked inside any of those, is
// absent in exactly the configuration where the failure is most likely: the one
// where delegation was switched off on purpose.
//
// The defect it guards: on 2026-09-09 a full foundation run on a real project
// convened zero reviewers and the report said nothing about it, because nothing
// asked it to. A clean review and a review that never happened leave identical
// evidence — none — so silence gets read as approval, which is the one reading
// that is never true.
//
// The datum that discriminates: compose with the delegation toggle OFF. If the
// rule survives that, it does not live in the removable section.
func TestReglaDelSilencioLlegaSiempre(t *testing.T) {
	requeridas := []string{
		"could not convene",
		"silence about a reviewer reads as approval",
	}

	combos := map[string][]string{
		"todos los toggles": legacyToggles(),
		"ningun toggle":     {},
		"sin delegation":    sinToggle("delegation"),
		"sin roles":         sinToggle("roles"),
		"sin review":        sinToggle("review"),
		"sin auditor":       sinToggle("auditor"),
	}

	for _, variante := range variantes {
		for nombre, toggles := range combos {
			t.Run(variante+"/"+nombre, func(t *testing.T) {
				bloque, err := Compose(toggles, variante)
				if err != nil {
					t.Fatalf("Compose(%q) error = %v", variante, err)
				}
				texto := normalizarMarkdown(bloque)
				for _, frase := range requeridas {
					if !strings.Contains(texto, normalizarMarkdown(frase)) {
						t.Errorf("falta la regla del silencio: %q", frase)
					}
				}
			})
		}
	}
}
