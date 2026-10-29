package config

import (
	"fmt"
	"strings"

	"github.com/Group-Active-IA/active-stack/internal/filemerge"
)

// knownToggles is the catalogue of recognised toggle names.
// Any toggle not in this set causes Compose to return an error.
var knownToggles = map[string]bool{
	"lean-workflow":  true,
	"delegation":     true,
	"model-routing":  true,
	"engram":         true,
	"tdd":            true,
	"governance":     true,
	"evidence":       true,
	"auditor":        true,
	"roles":          true,
	"review":         true,
	"knowledge-base": true,
	"senior":         true,
	"ship":           true,
}

// KnownToggles returns every recognised toggle name, in composition order.
//
// It exists so a test can say "all toggles" and mean all of them. The golden
// tests used to hardcode five names, so the four toggles added after them were
// never composed together: their assembly ORDER — argued for in the comments
// below — was fixed by nothing, and no golden ever showed the fragments
// side by side. A hardcoded list drifts silently every time the catalogue grows.
func KnownToggles() []string {
	return []string{
		"lean-workflow",
		"delegation", "model-routing", "governance", "evidence",
		"auditor", "roles", "review", "knowledge-base", "senior", "ship", "engram", "tdd",
	}
}

// Compose assembles the sdd-orchestrator block for the given variant and
// active toggles. It is a PURE function: same inputs always produce the same
// output. No I/O is performed; all assets are read from the embedded FS.
//
// With lean-workflow, composition returns compact universal role contracts and
// selected integration contracts, independent of the adapter's native roles.
// Without it, the legacy assembly order is retained below.
// Assembly order (deterministic, independent of toggle slice order):
//
//  1. base (always): assets/<variant>/sdd-orchestrator.md
//  2. Subtractive: if "delegation" NOT active → remove <!-- active-stack:sdd-delegation --> section
//  3. Subtractive: if "model-routing" NOT active → remove <!-- active-stack:sdd-model-assignments --> section
//  4. Additive: if "governance" active → append assets/governance.md
//     4b. Additive: if "evidence" active → append assets/evidence-report.md
//     4c. Additive: if "auditor" active → append assets/auditor.md
//     4d. Additive: if "roles" active → append assets/roles.md
//     4e. Additive: if "review" active → append assets/review.md
//     4f. Additive: if "knowledge-base" active → append assets/knowledge-base.md
//     4g. Additive: if "senior" active → append assets/senior.md
//     4h. Additive: if "ship" active → append assets/ship.md
//  5. Additive: if "engram" active → append assets/engram-protocol.md
//  6. Additive: if "tdd" active → append inline TDD flag + assets/strict-tdd.md
//
// Returns an error if any toggle is not in the known catalogue.
func Compose(toggles []string, variantKey string) (string, error) {
	// Validate all toggles first.
	for _, t := range toggles {
		if !knownToggles[t] {
			return "", fmt.Errorf("config compose: unknown toggle %q (known: %s)", t, strings.Join(KnownToggles(), ", "))
		}
	}

	// The lean lifecycle is a complete contract. Validate before selecting it so
	// unknown options cannot be hidden by the early return. Without this toggle,
	// the historical composition below remains unchanged.
	for _, t := range toggles {
		if t == "lean-workflow" {
			return composeLean(toggles, variantKey)
		}
	}

	// Build an active toggle set for O(1) lookup.
	active := make(map[string]bool, len(toggles))
	for _, t := range toggles {
		active[t] = true
	}

	// Step 1: load base variant.
	base, err := loadVariant(variantKey)
	if err != nil {
		return "", fmt.Errorf("config compose: %w", err)
	}

	result := base

	// Step 2: subtractive — delegation.
	if !active["delegation"] {
		result = filemerge.InjectMarkdownSection(result, "sdd-delegation", "")
	}

	// Step 3: subtractive — model-routing.
	if !active["model-routing"] {
		result = filemerge.InjectMarkdownSection(result, "sdd-model-assignments", "")
	}

	// Steps 4-6: additive fragments (always appended in fixed order).
	var additives []string

	// Step 4: governance.
	if active["governance"] {
		gov, err := loadFragment("governance.md")
		if err != nil {
			return "", fmt.Errorf("config compose: %w", err)
		}
		additives = append(additives, gov)
	}

	// Step 4b: typed reporting. Before engram and TDD because it frames how every
	// other section's output should come back: a summary that reaches the reader
	// untyped is a summary the reader cannot weigh.
	if active["evidence"] {
		ev, err := loadFragment("evidence-report.md")
		if err != nil {
			return "", fmt.Errorf("config compose: %w", err)
		}
		additives = append(additives, ev)
	}

	// Step 4c: the auditor. After the typed report on purpose: with no declared
	// categories it has nothing to audit, only prose to opine about.
	if active["auditor"] {
		au, err := loadFragment("auditor.md")
		if err != nil {
			return "", fmt.Errorf("config compose: %w", err)
		}
		additives = append(additives, au)
	}

	// Step 4d: role specialists. After the auditor and before engram: it decides
	// WHO gets convened, so it must sit downstream of the rules that say how any
	// of them must report back.
	if active["roles"] {
		rl, err := loadFragment("roles.md")
		if err != nil {
			return "", fmt.Errorf("config compose: %w", err)
		}
		additives = append(additives, rl)
	}

	// Step 4e: adversarial review. Last of the review family, and after roles on
	// purpose: roles decides WHO is convened, this decides HOW a review runs and
	// what it may not do with what it finds. Reversing them would define the
	// procedure before there is anybody to follow it.
	if active["review"] {
		rv, err := loadFragment("review.md")
		if err != nil {
			return "", fmt.Errorf("config compose: %w", err)
		}
		additives = append(additives, rv)
	}

	if !active["review"] {
		degradation, err := loadFragment("verify-without-review.md")
		if err != nil {
			return "", fmt.Errorf("config compose: %w", err)
		}
		// List review first, then other absences in canonical order.
		missing := []string{"review"}
		for _, toggle := range KnownToggles() {
			if toggle != "review" && toggle != "lean-workflow" && !active[toggle] {
				missing = append(missing, toggle)
			}
		}
		additives = append(additives, degradation+"\nMissing toggles: "+strings.Join(missing, ", ")+".\n")
	}

	// Step 4f: the knowledge base. After review and before the senior on purpose:
	// everything above it is scoped to the change in flight, and this is the first
	// section whose subject OUTLIVES the change. It also reports in the typed
	// categories and its drift report is a JUDGMENT the auditor can read, so it
	// must sit downstream of both. Still upstream of engram, for the reason it was
	// placed there: engram is memory for the agent and this is a map on disk for
	// the person, so engram composing later must not read as if its protocol
	// already covered this.
	if active["knowledge-base"] {
		kb, err := loadFragment("knowledge-base.md")
		if err != nil {
			return "", fmt.Errorf("config compose: %w", err)
		}
		additives = append(additives, kb)
	}

	// Step 4g: the senior. Last of the additive review family, and after both
	// review and knowledge-base on purpose: it is the only section defined by
	// CONTRAST with the adversarial passes — memory against focus — so a reader
	// who meets it before meeting them cannot see what it adds. Before engram
	// because engram is where the half of its memory that is not on disk comes
	// from, and the fragment says so: it must not read as if engram were already
	// assumed.
	if active["senior"] {
		seniorAsset := "senior.md"
		if !active["review"] {
			seniorAsset = "senior-without-review.md"
		} else if !active["auditor"] {
			seniorAsset = "senior-without-auditor.md"
		}
		sr, err := loadFragment(seniorAsset)
		if err != nil {
			return "", fmt.Errorf("config compose: %w", err)
		}
		additives = append(additives, sr)
	}

	// Step 4h: ship. Last of the additive family, and after the senior on purpose:
	// its division of labour IS the senior's — the diagnosis reports, the
	// orchestrator answers — and a reader who meets the borrowed rule before the
	// fragment that names it meets it without its argument. It also cites the
	// verify verdict the review family produces, so both must already be on the
	// page. Still before engram, which is a separate toggle nothing here may read
	// as already assumed.
	if active["ship"] {
		sp, err := loadFragment("ship.md")
		if err != nil {
			return "", fmt.Errorf("config compose: %w", err)
		}
		additives = append(additives, sp)
	}

	// Step 5: engram.
	if active["engram"] {
		eng, err := loadFragment("engram-protocol.md")
		if err != nil {
			return "", fmt.Errorf("config compose: %w", err)
		}
		additives = append(additives, eng)
	}

	// Step 6: tdd flag + module.
	if active["tdd"] {
		tddModule, err := loadFragment("strict-tdd.md")
		if err != nil {
			return "", fmt.Errorf("config compose: %w", err)
		}
		if !active["review"] {
			tddModule = strings.ReplaceAll(tddModule, "The breaker and auditor count the tests themselves; do not accept the constructor's count as their evidence.", "Ordinary verify counts tests independently; do not accept the constructor's count as evidence.")
		}
		tddBlock := "Strict TDD Mode: enabled\n\n" + tddModule
		additives = append(additives, tddBlock)
	}

	if len(additives) > 0 {
		combined := strings.Join(additives, "\n\n")
		if !strings.HasSuffix(result, "\n") {
			result += "\n"
		}
		result += "\n" + combined
	}

	return result, nil
}
