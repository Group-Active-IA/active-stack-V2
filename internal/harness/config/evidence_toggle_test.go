package config_test

// ──────────────────────────────────────────────────────────────────────────────
// Typed reporting: separating what the agent knows from what it assumes
//
// Today an agent hands everything back in the same tone: what it executed, what
// it thinks, and what it took for granted all arrive wearing the same face and
// cannot be told apart by reading.
//
// The case that motivates this: an agent reports "the billing integration works
// correctly" after testing against a mock it wrote itself, shaped the way it
// imagined the response looks. It never called the service. That sentence is
// indistinguishable from "the tests pass", which is true.
//
// The cost does not show up there. It shows up in production, when the service
// answers XML. By then the project record states the integration was verified,
// and nobody knows that claim never was.
//
// Separating evidence from judgment comes from Alberto Cortez's architecture
// report (v3.1, section 5.2): technical judgments remain evaluations by agents
// or people, recorded "without presenting them as facts demonstrated by the
// validator". The third category, the declared assumption, is this stack's
// addition: a badly founded judgment can be argued with, a silent assumption
// cannot even be argued with, because nobody knows it is there.
// ──────────────────────────────────────────────────────────────────────────────

import (
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

var evidenceVariants = []string{
	"claude", "opencode", "codex", "cursor",
	"gemini", "windsurf", "antigravity", "generic",
}

func evidenceFragment(t *testing.T) string {
	t.Helper()
	body, err := config.Compose([]string{"evidence"}, "claude")
	if err != nil {
		t.Fatalf("Compose([evidence]): %v", err)
	}
	return strings.ToLower(body)
}

// TestCompose_EvidenceToggle_IsKnown asserts the toggle is recognised.
func TestCompose_EvidenceToggle_IsKnown(t *testing.T) {
	for _, v := range evidenceVariants {
		if _, err := config.Compose([]string{"evidence"}, v); err != nil {
			t.Errorf("Compose([evidence], %q): %v", v, err)
		}
	}
}

// TestCompose_EvidenceToggle_Subtractive asserts a block that did not ask for it
// does not carry it.
//
// The marker is the fragment's own heading, and that is deliberate: a bare word
// like "evidence" can be written by any other fragment, and a marker somebody
// else can write does not distinguish anything.
func TestCompose_EvidenceToggle_Subtractive(t *testing.T) {
	off, err := config.Compose([]string{"delegation"}, "claude")
	if err != nil {
		t.Fatalf("Compose without evidence: %v", err)
	}
	if strings.Contains(off, "## Typed Reporting") {
		t.Error("the evidence fragment leaks into a block that did not ask for it")
	}
}

// TestCompose_EvidenceToggle_HasTheThreeCategories asserts the three categories
// are named. Two of them are not enough: the split that matters is between what
// was verified and what was taken for granted.
func TestCompose_EvidenceToggle_HasTheThreeCategories(t *testing.T) {
	body := evidenceFragment(t)
	for _, cat := range []string{"evidence", "judgment", "assumption"} {
		if !strings.Contains(body, cat) {
			t.Errorf("the fragment never names the %q category", cat)
		}
	}
}

// TestCompose_EvidenceToggle_DefinesWhatEachRequires asserts each category has an
// admission rule. Without them the categories are decorative labels and
// everything ends up filed as evidence.
func TestCompose_EvidenceToggle_DefinesWhatEachRequires(t *testing.T) {
	body := evidenceFragment(t)
	required := map[string]string{
		"the command":         "without the command that produced it, evidence is a memory",
		"timestamp":           "output from two hours ago describes a tree that no longer exists",
		"confidence":          "a judgment with no confidence level reads as a fact",
		"even if nobody asks": "what gets lost is not what the agent ignores, it is what it knew it was assuming and did not say",
	}
	for phrase, why := range required {
		if !strings.Contains(body, phrase) {
			t.Errorf("the fragment is missing %q: %s", phrase, why)
		}
	}
}

// TestCompose_EvidenceToggle_HasTheDowngradeRule asserts nothing may be promoted
// to a stronger category by convenience.
//
// This is the rule that keeps the format from becoming ceremony: without it,
// everything gets filed as evidence and the three sections say nothing.
func TestCompose_EvidenceToggle_HasTheDowngradeRule(t *testing.T) {
	body := evidenceFragment(t)
	for _, phrase := range []string{"never moves up", "when in doubt"} {
		if !strings.Contains(body, phrase) {
			t.Errorf("missing the downgrade rule (%q): without it everything ends up "+
				"filed as evidence and the format becomes decoration", phrase)
		}
	}
}

// TestCompose_EvidenceToggle_SaysDeclaringIsGood asserts the fragment states that
// declaring an assumption is a merit, not an admission of failure.
//
// This is the mitigation for the most likely failure mode: if declaring feels
// like confessing, the cheapest way out is to stop declaring, and the metric
// goes to zero while the problem stays.
func TestCompose_EvidenceToggle_SaysDeclaringIsGood(t *testing.T) {
	if !strings.Contains(evidenceFragment(t), "is a merit") {
		t.Error("the fragment never says declaring an assumption is a merit: " +
			"if declaring feels like confessing, agents stop declaring")
	}
}

// TestCompose_EvidenceToggle_TellsTheOrchestratorToDemandIt asserts the fragment
// instructs the orchestrator to REQUEST the format when launching a sub-agent.
//
// Without this the cable stops halfway: the orchestrator reports typed to the
// person, but everything it receives from sub-agents is prose, so it has to
// classify second-hand what it did not observe. The one who ran the command is
// the only one who can say whether it is evidence.
func TestCompose_EvidenceToggle_TellsTheOrchestratorToDemandIt(t *testing.T) {
	body := evidenceFragment(t)
	for _, phrase := range []string{"when you launch a sub-agent", "ask for"} {
		if !strings.Contains(body, phrase) {
			t.Errorf("the fragment never tells the orchestrator to demand the format "+
				"(%q): a summary that arrives as prose cannot be re-typed by whoever "+
				"did not run it", phrase)
		}
	}
}
