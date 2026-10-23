package config_test

import (
	"strings"
	"testing"
)

func TestShipEsperaCICompleto(t *testing.T) {
	body := shipFragment(t)
	if strings.Contains(body, "a cancelled job counts as green when its rerun passes.") {
		t.Error("CI contradictorio: A cancelled job counts as green when its rerun passes.")
	}
	for _, want := range []string{
		"Before evaluating CI, record the pushed SHA, expected workflows/events and expanded matrix jobs.",
		"Every required job must finish in success on the pushed SHA.",
		"gh run list --commit <sha>",
		"gh run watch",
		"Inspect the jobs of every expected run; never choose only the favorable run.",
		"For this repository, require .github/workflows/ci.yml on push, with test on Ubuntu and Windows.",
		"For a PR, inspect its required checks and verify the source SHA, not the synthetic merge SHA.",
		"In other repositories, derive the required set from their workflows and required checks; a tag release is not required for a branch push.",
		"One successful job with another pending is not green.",
		"queued, in_progress, cancelled, skipped, neutral and timed_out are not green.",
		"A result for another SHA is not green.",
		"Zero runs, no CI or a failed query are not green; report the state and leave delivery unconfirmed.",
		"On failure, report it without stacking another fix as if the push were confirmed.",
		"A new push invalidates the previous result and requires evaluating the new SHA.",
		"The orchestrator waits for the complete CI result before declaring delivery ready.",
	} {
		t.Run(want, func(t *testing.T) {
			if !strings.Contains(body, strings.ToLower(want)) {
				t.Errorf("falta contrato de CI: %s", want)
			}
		})
	}
}
