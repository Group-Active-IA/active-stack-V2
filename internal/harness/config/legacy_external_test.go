package config_test

import "github.com/Group-Active-IA/active-stack/internal/harness/config"

// legacyToggles selects the explicit legacy mode for existing contract tests.
func legacyToggles() []string {
	var out []string
	for _, toggle := range config.KnownToggles() {
		if toggle != "lean-workflow" {
			out = append(out, toggle)
		}
	}
	return out
}
