package config

// legacyToggles keeps the historical content/ordering regression suite explicit.
// Lean lifecycle contracts have their own tests; adding a mode must not silently
// repurpose the older suite to exercise a different lifecycle.
func legacyToggles() []string {
	var out []string
	for _, toggle := range KnownToggles() {
		if toggle != "lean-workflow" {
			out = append(out, toggle)
		}
	}
	return out
}
