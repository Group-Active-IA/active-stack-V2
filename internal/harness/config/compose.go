package config

import (
	"fmt"
	"strings"

	"github.com/Group-Active-IA/active-stack/internal/filemerge"
)

// knownToggles is the catalogue of recognised toggle names.
// Any toggle not in this set causes Compose to return an error.
var knownToggles = map[string]bool{
	"delegation":    true,
	"model-routing": true,
	"engram":        true,
	"tdd":           true,
	"governance":    true,
	"evidence":      true,
	"auditor":       true,
	"roles":         true,
	"review":        true,
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
		"delegation", "model-routing", "governance", "evidence",
		"auditor", "roles", "review", "engram", "tdd",
	}
}

// Compose assembles the sdd-orchestrator block for the given variant and
// active toggles. It is a PURE function: same inputs always produce the same
// output. No I/O is performed; all assets are read from the embedded FS.
//
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
//  5. Additive: if "engram" active → append assets/engram-protocol.md
//  6. Additive: if "tdd" active → append inline TDD flag + assets/strict-tdd.md
//
// Returns an error if any toggle is not in the known catalogue.
func Compose(toggles []string, variantKey string) (string, error) {
	// Validate all toggles first.
	for _, t := range toggles {
		if !knownToggles[t] {
			return "", fmt.Errorf("config compose: unknown toggle %q (known: delegation, model-routing, engram, tdd, governance, evidence, auditor, roles, review)", t)
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
