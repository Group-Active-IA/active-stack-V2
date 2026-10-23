package config_test

// ──────────────────────────────────────────────────────────────────────────────
// Revisión adversarial: primero conformidad, después ruptura
//
// `/opsx:verify` verifica CONFORMIDAD: que lo implementado coincida con los
// artefactos (specs, tasks, design). Sus tres dimensiones son completeness,
// correctness y coherence, todas contra lo que ya está escrito. Eso está bien y
// este fragmento no lo toca.
//
// Falta la otra mitad: buscar activamente lo que rompe. Nadie prueba qué pasa
// cuando el proveedor externo devuelve basura, porque eso no está en las specs;
// si estuviera, ya se habría pensado.
//
// LO QUE ESTE ARCHIVO CUIDA, además de que las seis reglas estén escritas:
//
//   - Cada test acota su sección y le saca el encabezado. Es el bug que
//     roles_toggle_test.go ya encontró dos veces: un test que busca una frase
//     que el título de su propia sección también escribe pasa con la sección
//     vaciada.
//   - La independencia de las tres pasadas es PARCIAL, y el fragmento tiene que
//     decirlo. El orquestador controla una cosa: escribe cada prompt y no pega
//     ahí la salida de las otras. Quedan DOS canales abiertos que un prompt
//     prolijo no cierra:
//       · el directorio del change, que una pasada puede leer;
//       · engram, que es el grande. El template de prompt del base le pasa a
//         cada sub-agente `topic_key: "opsx/{change-name}/{phase}"`: tres
//         pasadas del mismo change en la misma fase comparten clave. Y
//         engram-protocol.md compone DESPUÉS de este fragmento — última palabra
//         — declarándose MANDATORY y ALWAYS ACTIVE, con búsqueda proactiva al
//         empezar y `mem_save` después de cada acción completada. La regla de
//         orden no lo cubre: esas escrituras no pasan al escribir el veredicto,
//         pasan en medio de cada pasada.
//     Es el pecado de roles.md invertido: aquel inventaba un mecanismo que no
//     existía, éste negaba un canal que sí. Dos anti-asserts lo cuidan: uno
//     impide afirmar aislamiento del lanzamiento, el otro impide volver a
//     declarar que el prompt es el único canal.
// ──────────────────────────────────────────────────────────────────────────────

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

// reviewHeading is the fragment's own heading, and it is the marker on purpose:
// bare words like "review" or "verify" are written by the base block and by
// several other fragments, and a marker somebody else can write distinguishes
// nothing.
const reviewHeading = "## adversarial review"

// reviewFragment returns ONLY the review fragment, never the whole block.
//
// The base block already writes "verify", "apply" and "archive"; asserting
// against the composed output would pass with the fragment empty.
func reviewFragment(t *testing.T) string {
	t.Helper()
	body, err := config.Compose([]string{"review"}, "claude")
	if err != nil {
		t.Fatalf("Compose([review]): %v", err)
	}
	lower := strings.ToLower(body)
	i := strings.Index(lower, reviewHeading)
	if i < 0 {
		t.Fatalf("the composed block does not carry the review fragment (heading %q missing)", reviewHeading)
	}
	return lower[i:]
}

// reviewSection returns the BODY of one "### " section, its own heading line
// dropped. Dropping the heading is the whole point: a test that matches a phrase
// its section title writes is green with the section emptied.
func reviewSection(t *testing.T, headingPrefix string) string {
	t.Helper()
	body := reviewFragment(t)
	i := strings.Index(body, headingPrefix)
	if i < 0 {
		t.Fatalf("the review fragment has no section starting with %q", headingPrefix)
	}
	rest := body[i:]
	nl := strings.Index(rest, "\n")
	if nl < 0 {
		t.Fatalf("section %q has no body", headingPrefix)
	}
	rest = rest[nl+1:]
	if j := strings.Index(rest, "\n### "); j >= 0 {
		rest = rest[:j]
	}
	return rest
}

func requirePhrases(t *testing.T, body string, required map[string]string, what string) {
	t.Helper()
	for phrase, why := range required {
		if !strings.Contains(body, phrase) {
			t.Errorf("%s is missing %q — %s", what, phrase, why)
		}
	}
}

// TestCompose_ReviewToggle_IsKnown asserts every variant accepts it: the
// fragment is shared, so no variant may reject it.
func TestCompose_ReviewToggle_IsKnown(t *testing.T) {
	for _, v := range evidenceVariants {
		if _, err := config.Compose([]string{"review"}, v); err != nil {
			t.Errorf("Compose([review], %q): %v", v, err)
		}
	}
}

// TestCompose_ReviewToggle_Subtractive asserts a block that did not ask for it
// does not carry it.
func TestCompose_ReviewToggle_Subtractive(t *testing.T) {
	off, err := config.Compose([]string{"delegation", "model-routing", "engram", "tdd", "governance", "evidence", "auditor", "roles"}, "claude")
	if err != nil {
		t.Fatalf("Compose without review: %v", err)
	}
	if strings.Contains(strings.ToLower(off), reviewHeading) {
		t.Error("the review fragment leaks into a block that did not ask for it")
	}
}

// TestCompose_ReviewToggle_ConformityFirstThenRupture asserts the fragment keeps
// both questions and says they are different.
//
// Collapsing them is the failure: conformity alone certifies that a change does
// what somebody already thought of, which is the half nobody needed help with.
func TestCompose_ReviewToggle_ConformityFirstThenRupture(t *testing.T) {
	body := reviewSection(t, "### two questions")
	requirePhrases(t, body, map[string]string{
		"conformity":                 "the first question is the one /opsx:verify already answers, and it stays",
		"already have thought of it": "the reason the rupture question exists: what is in the specs was already imagined",
		"returns garbage":            "the worked example of a failure no spec names",
	}, "the two-questions section")
}

// TestCompose_ReviewToggle_ThreeIndependentPasses asserts it demands three, and
// defines what independent means.
//
// "Three passes" without the definition degrades into one pass run three times:
// the second reviewer reads the first verdict, agrees, and the count is
// decoration.
func TestCompose_ReviewToggle_ThreeIndependentPasses(t *testing.T) {
	body := reviewFragment(t)
	if !strings.Contains(body, "one breaker covers all three mandatory lenses") {
		t.Errorf("falta contrato: %s", "one breaker covers all three mandatory lenses")
	}
	if !strings.Contains(body, "not independent between lenses") {
		t.Errorf("falta contrato: %s", "not independent between lenses")
	}

}

// TestCompose_ReviewToggle_IndependenceIsByConstructionNotByMechanism asserts
// the fragment states the limit of what it can enforce.
//
// The orchestrator writes each prompt, so withholding the other passes' output
// is genuinely in its hands — unlike the tool restriction in roles.md. What it
// cannot stop is a pass reading a verdict already written into the change
// directory, which is why the ordering rule exists.
func TestCompose_ReviewToggle_IndependenceIsByConstructionNotByMechanism(t *testing.T) {
	body := reviewFragment(t)
	if !strings.Contains(body, "declare sources read") {
		t.Errorf("falta contrato: %s", "declare sources read")
	}
	if !strings.Contains(body, "contamination") {
		t.Errorf("falta contrato: %s", "contamination")
	}
	if !strings.Contains(body, "not proof of isolation") {
		t.Errorf("falta contrato: %s", "not proof of isolation")
	}

}

// TestCompose_ReviewToggle_DoesNotClaimTheLaunchEnforcesIsolation is the guard
// on the limit above.
//
// The documented launch is Agent({description, model, prompt}) — three keys,
// none of them an isolation setting. Independence is achieved by what the
// orchestrator chooses to put in a prompt; claiming the launch guarantees it
// would hide that the ordering rule is the only thing holding it up.
func TestCompose_ReviewToggle_DoesNotClaimTheLaunchEnforcesIsolation(t *testing.T) {
	body := reviewFragment(t)
	for _, banned := range []string{
		"guaranteed independent",
		"cannot read",
		"no edit tools",
		"sandboxed",
		// The inverse overclaim, and the one this fragment actually shipped:
		// saying the prompt is the whole story while the composed block hands
		// all three passes one engram topic_key.
		"the only channel",
		"share nothing",
		"nothing connects them",
	} {
		if strings.Contains(body, banned) {
			t.Errorf("the review fragment claims %q. The documented launch is "+
				"Agent({description, model, prompt}): it carries no isolation key and no "+
				"tool list. Independence here is a construction choice, not a guarantee, "+
				"and saying otherwise hides that the ordering rule is what holds it up", banned)
		}
	}
}

// TestCompose_ReviewToggle_LensesAreDifferentNotRedundant asserts the three
// passes look through different lenses, and names them.
//
// Three passes with the same lens are redundancy: they find the same kind of bug
// three times and report it as three findings. The point of three is coverage of
// three kinds.
func TestCompose_ReviewToggle_LensesAreDifferentNotRedundant(t *testing.T) {
	body := reviewSection(t, "### different lenses")
	requirePhrases(t, body, map[string]string{
		"the same kind of bug three times": "what redundancy buys, stated so it can be refused",
		"three kinds":                      "what diversity buys instead",
		"who calls it":                     "lens two: extracted logic with no caller is dead on arrival",
		"still tell the truth":             "lens three: a comment that stopped being true is a lie with authority",
	}, "the lenses section")
	// The three lens LABELS, asserted as labelled bullets rather than as bare
	// words. "correctness" alone would not have worked: the section's closing
	// sentence ("a pass told only to review it reviews correctness") writes the
	// word too, so deleting the whole correctness lens left that assert green.
	// Same bug shape this repo has now hit seven times.
	for _, lens := range []string{"- **correctness**", "- **wiring**", "- **truth of the comments**"} {
		if !strings.Contains(body, lens) {
			t.Errorf("the lenses section has no %q lens: three passes without three "+
				"named lenses is one lens run three times", lens)
		}
	}
	// Lens two is the one that gets dropped first, because a reader who sees
	// "correctness" assumes it is covered. It has two halves and both matter.
	if !strings.Contains(body, "and is that caller tested") {
		t.Error("the wiring lens asks who calls the extracted logic but never whether " +
			"that call is tested: an untested caller is the bug this lens exists to find")
	}
}

// TestCompose_ReviewToggle_ReportAuthority separates reporting from closure.
func TestCompose_ReviewToggle_ReportAuthority(t *testing.T) {
	body := reviewSection(t, "### it reports")
	for _, hit := range debtViolations(body) {
		t.Error(hit)
	}
	for _, pattern := range []string{`auditor.*classifies.*decides`, `orchestrator.*executes`, `debt.*justified.*auditor`, `findings.*visible`} {
		if !regexp.MustCompile(pattern).MatchString(strings.Join(strings.Fields(body), " ")) {
			t.Errorf("falta relación de autoridad: %s", pattern)
		}
	}
}

// TestCompose_ReviewToggle_EmptyVerdictIsAFinding asserts silence is not an
// allowed output.
//
// Without the explicit "three passes, nothing found" line, a clean review and a
// review that never ran produce identical evidence: nothing.
func TestCompose_ReviewToggle_EmptyVerdictIsAFinding(t *testing.T) {
	body := reviewSection(t, "### an empty verdict")
	requirePhrases(t, body, map[string]string{
		"found nothing, say so":                  "the instruction",
		"a clean review from one that never ran": "the reason: the two are indistinguishable otherwise",
	}, "the empty-verdict section")
}

// TestCompose_ReviewToggle_NamesTheDiscriminatingDatum asserts the fragment
// demands the observation that separates a hypothesis from its alternative,
// before the finding is accepted or rejected.
//
// "Show me the logs" collects material that fits whatever was already believed.
// Naming the datum in advance is what makes the check a test rather than a
// search for agreement.
func TestCompose_ReviewToggle_NamesTheDiscriminatingDatum(t *testing.T) {
	body := reviewSection(t, "### the datum that discriminates")
	requirePhrases(t, body, map[string]string{
		"separates the hypothesis from its alternative": "the definition of the datum",
		// NOT the bare word "before": it also opens the sentence "Before accepting
		// or rejecting any finding", so deleting the naming rule left the assert
		// green. Verified by degrading it.
		"in advance of looking": "the datum is named before the evidence is seen, which is what stops it being chosen to fit",
		"look for this line":    "the worked example of the naming, in the form an agent can execute",
		"if it is absent":       "the other branch — a datum with one branch discriminates nothing",
		"from this run":         "evidence from an earlier run describes a tree that no longer exists",
	}, "the discriminating-datum section")
}

// TestCompose_ReviewToggle_NamesTheSharedMemoryChannel asserts the fragment does
// not deny a channel the composed block actually opens.
//
// EL CANAL. The sub-agent prompt template in the base block hands every
// sub-agent `topic_key: "opsx/{change-name}/{phase}"`. Three verify passes on
// one change share change-name AND phase, so they share the key. And
// engram-protocol.md composes AFTER this fragment — last word — declaring itself
// MANDATORY and ALWAYS ACTIVE, ordering a proactive `mem_search` before starting
// work that might have been done before, and a `mem_save` after every completed
// action.
//
// So pass two can read pass one's conclusions without either of them doing
// anything wrong. The ordering rule about the change directory does not cover
// it: those writes do not happen when the verdict is written, they happen
// mid-pass.
//
// This is roles.md's sin inverted. That fragment claimed a mechanism it did not
// have; claiming independence while a shared memory key is open is denying a
// mechanism that is there.
func TestCompose_ReviewToggle_NamesTheSharedMemoryChannel(t *testing.T) {
	body := reviewFragment(t)
	if !strings.Contains(body, "opsx/{change}/verify/round-{n}/{role}") {
		t.Errorf("falta contrato: %s", "opsx/{change}/verify/round-{n}/{role}")
	}
	if !strings.Contains(body, "do not search sibling roles") {
		t.Errorf("falta contrato: %s", "do not search sibling roles")
	}
	if !strings.Contains(body, "after every completed action") {
		t.Errorf("falta contrato: %s", "after every completed action")
	}

}

// TestCompose_ReviewToggle_VerdictIsTypedAndNamed asserts the verdict lives
// inside the typed report and has a findable name.
//
// This mirrors TestCompose_RolesToggle_LedgerLivesInsideTheTypedReport, and for
// the same reason: with governance, evidence, auditor, roles and review all on,
// this verdict is the fourth thing competing over what comes back. Flat prose is
// not merely a fourth format — evidence-report.md forbids it, since anything
// that is not evidence must arrive labelled with its author and its confidence.
// "Into the change" with no filename is also unfindable: nothing can go back and
// read a file whose name was never agreed.
func TestCompose_ReviewToggle_VerdictIsTypedAndNamed(t *testing.T) {
	body := reviewSection(t, "### it reports")
	requirePhrases(t, body, map[string]string{
		"verify-round-{n}.md": "a concrete filename: 'into the change' names nothing anybody can open",
		"evidence":            "the typed categories the verdict is filed under",
		"judgment":            "a finding is a judgment and carries its confidence",
		"confidence":          "the field evidence-report.md demands of any judgment",
	}, "the reports-not-blocks section")
}

// TestCompose_ReviewToggle_LabelsTheArchiveClaimAsAssumption asserts the
// fragment does not state as fact something it cannot check from this repo.
//
// That `openspec archive` moves the whole change directory was read out of the
// installed archive command, not verified by running it here. evidence-report.md
// has a category for exactly this, and using it is the fragment's own rule
// applied to itself.
func TestCompose_ReviewToggle_LabelsTheArchiveClaimAsAssumption(t *testing.T) {
	body := reviewSection(t, "### it reports")
	// NOT the bare word "assumption": the paragraph above already writes
	// "anything taken for granted is an ASSUMPTION" while listing the typed
	// categories, so the bare word stayed green with the label deleted and the
	// claim restated as fact. Verified by degrading it. What is asserted is the
	// label ON THIS CLAIM, plus the consequence of it being wrong — an
	// unfalsifiable assumption is a disclaimer, not a category.
	requirePhrases(t, body, map[string]string{
		"not verified from this repo":       "the label attached to this specific claim",
		"the verdict is dropped on archive": "what it costs if the assumption is false, so it can be checked",
	}, "the archive claim")
}
