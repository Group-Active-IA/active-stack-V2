package config_test

import (
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

func TestLeanWorkflow_SequentialContractsAndTDD(t *testing.T) {
	for _, variant := range []string{"claude", "codex", "opencode", "generic", "unknown"} {
		t.Run(variant, func(t *testing.T) {
			body, err := config.Compose([]string{"lean-workflow"}, variant)
			if err != nil {
				t.Fatal(err)
			}
			previous := -1
			for _, role := range []string{"### Planner", "### Plan reviewer", "### Constructor", "### Breaker", "### Auditor"} {
				at := strings.Index(body, role)
				if at <= previous {
					t.Fatalf("missing or out-of-order contract %q", role)
				}
				previous = at
			}
			for _, required := range []string{"LOW", "MEDIUM", "HIGH", "CRITICAL", "risk and scope", "TDD is mandatory for every code change", "execute RED before production code", "reverification", "tasks.md", "bounded context", "No fixed model or provider", "documentation-only"} {
				if !strings.Contains(body, required) {
					t.Errorf("missing lean contract %q", required)
				}
			}
		})
	}
}

func TestLeanWorkflow_OptionalContractsAndBudget(t *testing.T) {
	minimal, err := config.Compose([]string{"lean-workflow"}, "claude")
	if err != nil {
		t.Fatal(err)
	}
	for _, option := range []struct{ toggle, heading string }{
		{"governance", "## Human approval"}, {"knowledge-base", "## Project knowledge"},
		{"engram", "## Persistent memory"}, {"ship", "## Publication"}, {"evidence", "## Typed reports"},
	} {
		if strings.Contains(minimal, option.heading) {
			t.Errorf("optional contract leaked: %s", option.toggle)
		}
		body, err := config.Compose([]string{"lean-workflow", option.toggle}, "claude")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(body, option.heading) {
			t.Errorf("missing optional contract %s", option.toggle)
		}
	}
	all := config.KnownToggles()
	body, err := config.Compose(all, "claude")
	if err != nil {
		t.Fatal(err)
	}
	if len(body) > 18000 {
		t.Fatalf("lean output is %d bytes; budget 18000", len(body))
	}
	if !strings.Contains(body, "lean-workflow takes precedence") || !strings.Contains(body, "senior context is handled by the auditor") {
		t.Fatal("missing precedence contract")
	}
	for _, fixed := range []string{"model: opus", "model: sonnet", "| opus |", "| sonnet |"} {
		if strings.Contains(body, fixed) {
			t.Fatalf("fixed model leaked: %s", fixed)
		}
	}
	for i, j := 0, len(all)-1; i < j; i, j = i+1, j-1 {
		all[i], all[j] = all[j], all[i]
	}
	reordered, err := config.Compose(append(all, "lean-workflow"), "claude")
	if err != nil {
		t.Fatal(err)
	}
	if body != reordered {
		t.Fatal("toggle order/duplication changed output")
	}
	if _, err := config.Compose([]string{"lean-workflow", "not-a-toggle"}, "claude"); err == nil {
		t.Fatal("lean accepted unknown toggle")
	}
}

func TestLeanWorkflow_PhaseMapLowMergeAndSelfContainedBriefs(t *testing.T) {
	for _, variant := range []string{"claude", "codex", "opencode", "generic"} {
		t.Run(variant, func(t *testing.T) {
			body, err := config.Compose([]string{"lean-workflow"}, variant)
			if err != nil {
				t.Fatal(err)
			}
			for _, required := range []string{
				"| propose | planner, then plan-reviewer |",
				"| apply | constructor |",
				"| verify | breaker, then auditor |",
				"| archive |",
				"LOW may merge planning and plan review into the orchestrator",
				"MEDIUM and above dispatch five separate roles",
				"self-contained",
				"another provider's model",
				"typed report format (EVIDENCIA, JUICIO, SUPOSICIÓN, sources read, contamination by other reports)",
			} {
				if !strings.Contains(body, required) {
					t.Errorf("missing orchestrator resource %q", required)
				}
			}
		})
	}
}

func TestLeanWorkflow_DepthByLevelAndRecoveredLegacyMechanisms(t *testing.T) {
	core := []string{
		// rigor scales with the change level; TDD and rounds never disappear
		"## Depth by level",
		"| LOW | orchestrator may merge plan and review | RED then GREEN per behavior | one breaker round, then auditor |",
		"| CRITICAL |",
		"mutant",
		// 1 reports and provenance index
		"openspec/changes/<change>/informes/<phase>-<attempt>-<role>.md",
		"LEEME.md",
		"Identity without a dispatch record is SUPOSICIÓN",
		// 2 lessons are read before acting
		"read the numbered lessons (L-NNN)",
		// 3 auditor memory channels
		"openspec/changes/archive/",
		"declares which sources it used",
		// 4 lessons and debt survive without a knowledge base
		"destination pending",
		// 5 bounded correction round
		"fixes only the auditor's list",
		// 6 escalation with options
		"new approved scope, accepted justified debt or stopping the change",
		// 7 stable attempt counter
		"Session recovery neither increments nor resets the attempt counter",
		// 8 external identity provider walkthrough
		"external identity provider",
		// 9 skill registry rules in briefs
		"skill registry",
		// 10 group corrections by file
		"grouped by file or area",
	}
	for _, variant := range []string{"claude", "generic"} {
		body, err := config.Compose([]string{"lean-workflow"}, variant)
		if err != nil {
			t.Fatal(err)
		}
		for _, required := range core {
			if !strings.Contains(body, required) {
				t.Errorf("%s: missing %q", variant, required)
			}
		}
	}
	// 11 publication diagnosis gaps live in the ship contract
	ship, err := config.Compose([]string{"lean-workflow", "ship"}, "claude")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"No remote: stop", "uncommitted changes: ask", "recommend PR or direct push with reasons", "source SHA, not the synthetic merge SHA"} {
		if !strings.Contains(ship, required) {
			t.Errorf("ship: missing %q", required)
		}
	}
}

func TestLeanWorkflow_PerProjectRoleExecutors(t *testing.T) {
	for _, variant := range []string{"claude", "generic"} {
		body, err := config.Compose([]string{"lean-workflow"}, variant)
		if err != nil {
			t.Fatal(err)
		}
		for _, required := range []string{
			"## Role executors",
			"declared per project",
			"| breaker | codex plugin |",
			"A role without a row runs as a native subagent",
			"never substitute silently",
			"No credentials",
		} {
			if !strings.Contains(body, required) {
				t.Errorf("%s: missing executor contract %q", variant, required)
			}
		}
	}
	claude, err := config.Compose([]string{"lean-workflow"}, "claude")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(claude, "unless the project's Role executors assign it elsewhere") {
		t.Fatal("claude dispatch ignores per-project executors")
	}
}

func TestLeanWorkflow_NativeDispatchOnlyForClaude(t *testing.T) {
	claude, err := config.Compose([]string{"lean-workflow"}, "claude")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(claude, "subagent_type") {
		t.Fatal("claude variant lacks native role dispatch")
	}
	generic, err := config.Compose([]string{"lean-workflow"}, "generic")
	if err != nil {
		t.Fatal(err)
	}
	for _, variant := range []string{"codex", "opencode"} {
		body, err := config.Compose([]string{"lean-workflow"}, variant)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(body, "subagent_type") || body != generic {
			t.Fatalf("%s variant must equal generic without Claude dispatch", variant)
		}
	}
}

func TestLeanWorkflow_GenericFallbackAndLegacy(t *testing.T) {
	a, err := config.Compose([]string{"lean-workflow"}, "generic")
	if err != nil {
		t.Fatal(err)
	}
	b, err := config.Compose([]string{"lean-workflow"}, "unknown")
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Fatal("unknown variant did not use generic")
	}
	legacy, err := config.Compose([]string{"tdd", "review"}, "claude")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(legacy, "### Planner") {
		t.Fatal("lean workflow leaked into legacy")
	}
}

func TestLeanWorkflow_Goldens(t *testing.T) {
	for _, variant := range []string{"claude", "codex", "opencode", "generic"} {
		t.Run(variant, func(t *testing.T) {
			body, err := config.Compose(config.KnownToggles(), variant)
			if err != nil {
				t.Fatal(err)
			}
			name := "compose_lean_" + variant
			if *update {
				writeGolden(t, name, body)
				return
			}
			if body != loadGolden(t, name) {
				t.Fatalf("lean %s does not match golden", variant)
			}
		})
	}
}

func TestLeanWorkflow_HistoricalSwitchesCannotDisableContracts(t *testing.T) {
	for _, omitted := range []string{"tdd", "review", "auditor", "roles", "delegation", "model-routing", "senior"} {
		t.Run(omitted, func(t *testing.T) {
			var toggles []string
			for _, toggle := range config.KnownToggles() {
				if toggle != omitted {
					toggles = append(toggles, toggle)
				}
			}
			body, err := config.Compose(toggles, "claude")
			if err != nil {
				t.Fatal(err)
			}
			for _, required := range []string{"### Planner", "### Plan reviewer", "### Constructor", "### Breaker", "### Auditor", "TDD is mandatory", "required original reports", "explicit, dated human approval"} {
				if !strings.Contains(body, required) {
					t.Errorf("omitting %s removed %q", omitted, required)
				}
			}
			if strings.Contains(body, "## The Senior") || strings.Contains(body, "## Role Specialists") {
				t.Fatal("legacy attendance leaked")
			}
		})
	}
}

// The five roles report with EVIDENCIA/JUICIO/SUPOSICIÓN and verdicts are
// APROBADO/BLOQUEANTES; the core asked for EVIDENCE/JUDGMENT/ASSUMPTION, so
// the auditor received labels it was not told to look for (review of
// 2026-10-29, finding 10.6). One label set across core and roles.
func TestLeanWorkflow_OneReportLabelSet(t *testing.T) {
	body, err := config.Compose(config.KnownToggles(), "claude")
	if err != nil {
		t.Fatal(err)
	}
	for _, label := range []string{"EVIDENCIA", "JUICIO", "SUPOSICIÓN"} {
		if !strings.Contains(body, label) {
			t.Errorf("lean core lacks report label %q used by the roles", label)
		}
	}
	for _, label := range []string{"EVIDENCE", "JUDGMENT", "ASSUMPTION"} {
		if strings.Contains(body, label) {
			t.Errorf("lean core still uses English report label %q", label)
		}
	}
}
