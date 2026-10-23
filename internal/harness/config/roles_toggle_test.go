package config_test

// ──────────────────────────────────────────────────────────────────────────────
// Especialistas por rol: convocar a alguien, no fabricar un genérico
//
// Hoy cada delegación produce un sub-agente genérico briefeado desde cero,
// distinto cada vez. Dos cambios que tocan la misma capa los revisan dos
// criterios distintos, y ninguno de los dos quedó escrito en ningún lado.
//
// Un rol es un criterio estable: el mismo estándar cada vez que se toca esa
// capa. Lo que lo convierte en rol y no en título es que su terminado sea
// OBSERVABLE. "Hace buena UI" no se puede verificar; "controles operables por
// teclado con foco visible" sí.
//
// Los siete roles salen del informe de arquitectura de Alberto Cortez. Las tres
// reglas del final son aporte nuestro: registrar a quién NO se convocó, un piso
// de revisores que no elige el orquestador, y el revisor sin edición.
//
// DOS COSAS QUE ESTE ARCHIVO CUIDA, porque el fragmento ya las afirmó de más:
//
//   - **CAMBIÓ**: la fase `verify` sí existe. Cuando se escribió este archivo el
//     bloque documentaba explore → propose → apply → archive y la única
//     aparición de "verify" era la fila "Read 1-3 files to decide/verify" de la
//     tabla de delegación; asignarle el piso de revisión a una fase que el
//     orquestador no podía convocar habría sido un piso que no corre nunca. Pero
//     la fase existía afuera: `openspec init` instala `/opsx:verify` y la skill
//     `openspec-verify-change`. Lo que faltaba era que el bloque la nombrara, no
//     que la fase apareciera. Nombrada la fase, el piso vuelve a `verify`, que
//     es su lugar, y el parche que lo tenía estacionado en `apply` se retira.
//   - El único lanzamiento que el bloque documenta sigue siendo
//     Agent({description, model, prompt}). Tres claves, ninguna es una lista de
//     herramientas. "No es una instrucción, es qué herramientas recibe" era
//     falso: la frase sólo puede viajar dentro del prompt, o sea, como la
//     instrucción que decía no ser. Ese límite NO cambió.
//
// Por eso varios de estos tests asertan el MECANISMO y no el vocabulario, y dos
// assertan que afirmaciones de más NO vuelvan.
// ──────────────────────────────────────────────────────────────────────────────

import (
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

// rolesHeading is the fragment's own heading, and it is the marker on purpose:
// a bare word like "roles" or "architect" can be written by any other fragment
// or by the base block, and a marker somebody else can write distinguishes
// nothing.
const rolesHeading = "## role specialists"

// theSevenRoles is the roster, used both to check the roles are defined and to
// check the phase table actually assigns somebody.
var theSevenRoles = []string{
	"functional analyst", "architect", "backend", "frontend", "dba", "devops", "qa",
}

// rolesFragment returns ONLY the roles fragment, not the whole composed block.
//
// This matters more here than in the other toggles: the base orchestrator
// already contains the words "explore", "propose", "apply" and "archive". A test
// asserting them against the full output would pass with the fragment empty,
// which is the exact shape of a test that verifies nothing.
func rolesFragment(t *testing.T) string {
	t.Helper()
	body, err := config.Compose([]string{"roles"}, "claude")
	if err != nil {
		t.Fatalf("Compose([roles]): %v", err)
	}
	lower := strings.ToLower(body)
	i := strings.Index(lower, rolesHeading)
	if i < 0 {
		t.Fatalf("the composed block does not carry the roles fragment (heading %q missing)", rolesHeading)
	}
	return lower[i:]
}

// rolesSection returns the BODY of one "### " section of the fragment, with its
// own heading line dropped.
//
// Dropping the heading is the point. A previous version of the ledger test
// searched for "not convened", a phrase the section title itself writes: the
// instruction and the example could both be deleted and the test stayed green.
func rolesSection(t *testing.T, headingPrefix string) string {
	t.Helper()
	body := rolesFragment(t)
	i := strings.Index(body, headingPrefix)
	if i < 0 {
		t.Fatalf("the roles fragment has no section starting with %q", headingPrefix)
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

// rolesPhaseRows parses the phase table into phase → the rest of the row.
//
// Parsing instead of substring-matching is what makes the mapping test a test:
// the previous version asserted the five phase words were spelled somewhere, so
// emptying every "who is convened" cell left it green.
// It delegates to the shared markdown table parser in verify_phase_test.go and
// flattens each row, which is safe HERE because every consumer below asks only
// "does this row mention X anywhere". It was NOT safe in the copy this function
// used to carry: that one flattened inside the parser, so a caller that needed
// one specific column could not get it and silently compared the whole row
// instead. One parser now, and the flattening happens at the call site that
// wants it.
func rolesPhaseRows(t *testing.T) map[string]string {
	t.Helper()
	parsed, _ := tableRows(rolesFragment(t))
	rows := make(map[string]string, len(parsed))
	for phase, cells := range parsed {
		rows[phase] = rowText(cells)
	}
	if len(rows) == 0 {
		t.Fatal("the roles fragment has no phase table at all")
	}
	return rows
}

// TestCompose_RolesToggle_IsKnown asserts the toggle is recognised by every
// variant: the fragment is shared, so no variant may reject it.
func TestCompose_RolesToggle_IsKnown(t *testing.T) {
	for _, v := range evidenceVariants {
		if _, err := config.Compose([]string{"roles"}, v); err != nil {
			t.Errorf("Compose([roles], %q): %v", v, err)
		}
	}
}

// TestCompose_RolesToggle_Subtractive asserts a block that did not ask for it
// does not carry it.
func TestCompose_RolesToggle_Subtractive(t *testing.T) {
	off, err := config.Compose([]string{"delegation", "model-routing", "engram", "tdd", "governance", "evidence", "auditor"}, "claude")
	if err != nil {
		t.Fatalf("Compose without roles: %v", err)
	}
	if strings.Contains(strings.ToLower(off), rolesHeading) {
		t.Error("the roles fragment leaks into a block that did not ask for it")
	}
}

// TestCompose_RolesToggle_NamesTheSevenRoles asserts every role is defined in the
// roster section — not merely mentioned somewhere in the fragment.
//
// Scoping this is not pedantry: with the whole fragment as the haystack,
// deleting a role's entire definition still left this green, because the phase
// table happens to write its name too.
func TestCompose_RolesToggle_NamesTheSevenRoles(t *testing.T) {
	body := rolesSection(t, "### the seven roles")
	for _, role := range theSevenRoles {
		if !strings.Contains(body, role) {
			t.Errorf("the roster never defines the %q role", role)
		}
	}
}

// TestCompose_RolesToggle_EachRoleHasAnObservableDone asserts each role carries a
// criterion somebody can check, not a description of what it is good at.
//
// This is the whole difference between a role and a job title. A criterion that
// cannot be observed cannot be failed, and a standard that cannot be failed is
// decoration.
func TestCompose_RolesToggle_EachRoleHasAnObservableDone(t *testing.T) {
	body := rolesSection(t, "### the seven roles")
	required := map[string]string{
		"fails today":                                 "functional analyst: a requirement nobody can show false is a wish",
		"the alternative it rejected":                 "architect: a decision with no rejected alternative cannot be reviewed",
		"error contract":                              "backend: the happy path alone is not a contract",
		"operable by keyboard with visible focus":     "frontend: this is the check that replaces 'good UI'",
		"down that restores the previous shape":       "dba: a migration with no exercised down is a one-way door",
		"rollback was exercised":                      "devops: a rollback that was only described has never run",
		"fails without the change and passes with it": "qa: a test that passes before the change tests nothing",
	}
	for phrase, why := range required {
		if !strings.Contains(body, phrase) {
			t.Errorf("the roster is missing %q — %s", phrase, why)
		}
	}
}

// TestCompose_RolesToggle_MapsTheOPSXPhases asserts each phase the block can
// actually convene has somebody assigned to it.
//
// It parses the table and looks for role names in the row, so emptying the
// columns turns it red. The previous version only checked the five phase words
// were spelled, which is orthography, not mapping.
func TestCompose_RolesToggle_MapsTheOPSXPhases(t *testing.T) {
	rows := rolesPhaseRows(t)
	// The phases that convene somebody. `apply` is deliberately not among them:
	// see TestCompose_RolesToggle_ApplyConvenesNobodyOnPurpose.
	for _, phase := range []string{"explore", "propose", "verify", "archive"} {
		row, ok := rows[phase]
		if !ok {
			t.Errorf("the phase table has no row for %q", phase)
			continue
		}
		named := false
		for _, role := range theSevenRoles {
			if strings.Contains(row, role) {
				named = true
				break
			}
		}
		if !named {
			t.Errorf("the %q row names no role at all: a phase with an empty row assigns nobody", phase)
		}
	}
}

// TestCompose_RolesToggle_StaffsVerifyBecauseThePhaseNowExists is the successor
// of TestCompose_RolesToggle_DoesNotStaffAPhaseThatDoesNotExist, and it asserts
// the OPPOSITE of what that test asserted.
//
// POR QUÉ CAMBIÓ LA PROPIEDAD. The old test existed to stop somebody staffing a
// phase the orchestrator could not convene: a floor parked where nothing reaches
// it looks like a floor and runs never. That reasoning was correct and the guard
// was worth having. What changed is the world, not the test's taste: `openspec
// init --tools claude` installs `/opsx:verify` and the `openspec-verify-change`
// skill, and this block now names the phase in all five places that enumerate
// the cycle. The phase always existed in the tool; what was missing was the
// block telling the agent so.
//
// So the guard inverts. `verify` must now carry the mandatory QA and the
// reviewer pass — the two things the previous commit had to park in `apply` —
// and `apply` must have given them back.
func TestCompose_RolesToggle_StaffsVerifyBecauseThePhaseNowExists(t *testing.T) {
	rows := rolesPhaseRows(t)
	verify, ok := rows["verify"]
	if !ok {
		t.Fatal("the phase table has no `verify` row: the phase exists in the cycle " +
			"and the roster assigns nobody to it")
	}
	if !strings.Contains(verify, "qa") {
		t.Error("the mandatory QA is not in `verify`: it is the phase whose whole job " +
			"is checking the implementation, and it is where the floor belongs")
	}
	if !strings.Contains(verify, "reviewer") {
		t.Error("the reviewer pass is not in `verify`: the pass was only ever in `apply` " +
			"because `verify` was not reachable from this block, and now it is")
	}
	// And the patch is gone: `apply` no longer holds the pass it was lent.
	apply, ok := rows["apply"]
	if !ok {
		t.Fatal("the phase table has no `apply` row")
	}
	if strings.Contains(apply, "reviewer") {
		t.Error("`apply` still runs the reviewer pass: with `verify` reachable the pass " +
			"belongs there, and running it twice is the redundancy the roster refuses")
	}
}

// TestCompose_RolesToggle_ApplyConvenesNobodyOnPurpose asserts the emptied
// `apply` row says why it is empty.
//
// This fragment's own rule is that a role not convened is a decision and gets
// written down. A row that simply lost its content reads as an oversight, and
// the next person to read it puts the reviewers back.
func TestCompose_RolesToggle_ApplyConvenesNobodyOnPurpose(t *testing.T) {
	rows := rolesPhaseRows(t)
	apply, ok := rows["apply"]
	if !ok {
		t.Fatal("the phase table has no `apply` row")
	}
	if !strings.Contains(apply, "still moving") {
		t.Error("the `apply` row does not say why it convenes nobody: reviewing a diff " +
			"that is still being written is the reason, and an unexplained empty row " +
			"reads as an omission somebody will helpfully fix")
	}
}

// TestCompose_RolesToggle_DoesNotClaimVerifyIsUnreachable is the guard on the
// retired claim.
//
// The sentence below was true when it was written and is false now. A false
// statement about what the agent can reach is worse than a missing one: it tells
// the orchestrator not to try.
func TestCompose_RolesToggle_DoesNotClaimVerifyIsUnreachable(t *testing.T) {
	body := rolesFragment(t)
	for _, banned := range []string{
		"does not exist in this block",
		"there is no `/opsx:verify` command",
	} {
		if strings.Contains(body, banned) {
			t.Errorf("the roles fragment still claims %q. `openspec init` installs "+
				"/opsx:verify and openspec-verify-change, and this block now lists the "+
				"phase in all five places that enumerate the cycle", banned)
		}
	}
}

// TestCompose_RolesToggle_DoesNotConveneEverybodyEveryTime asserts the choice is
// driven by what the change touches.
//
// Convening all seven every time costs the same as convening nobody: it stops
// being a signal and the orchestrator learns to ignore it.
func TestCompose_RolesToggle_DoesNotConveneEverybodyEveryTime(t *testing.T) {
	body := rolesFragment(t)
	for _, phrase := range []string{"migration convenes no frontend", "convenes no dba"} {
		if !strings.Contains(body, phrase) {
			t.Errorf("the roles fragment is missing %q: without a selection rule the "+
				"orchestrator convenes everybody, which is the same as convening nobody", phrase)
		}
	}
}

// TestCompose_RolesToggle_RecordsWhoWasNotConvened asserts the ledger has a
// mechanism: one line per role left out, its reason, and a place inside the
// typed report to live.
//
// Asserted against the section BODY, without its heading. The previous version
// searched for "not convened" — a phrase the heading itself writes — so deleting
// the instruction and the example left it green. That is the same bug shape this
// file already found once, one section higher.
func TestCompose_RolesToggle_RecordsWhoWasNotConvened(t *testing.T) {
	body := rolesSection(t, "### the roles you did not convene")
	required := map[string]string{
		"one line":      "without a per-role line the ledger collapses into one vague sentence",
		"the reason":    "a name with no reason is a list, and a list cannot be argued with",
		"is a decision": "the point of writing it down: a role not called is a decision, not a silence",
	}
	for phrase, why := range required {
		if !strings.Contains(body, phrase) {
			t.Errorf("the ledger section is missing %q: %s", phrase, why)
		}
	}
	// It must name at least two roles in its worked example: an instruction with
	// no example is what the orchestrator skips first.
	named := 0
	for _, role := range theSevenRoles {
		if strings.Contains(body, role) {
			named++
		}
	}
	if named < 2 {
		t.Errorf("the ledger section shows %d roles in its example: without a worked "+
			"example naming real roles, the format is a suggestion", named)
	}
}

// TestCompose_RolesToggle_LedgerLivesInsideTheTypedReport asserts the ledger is
// filed under a category of the typed report instead of inventing a parallel
// format.
//
// Three instructions now compete over what a sub-agent hands back: the base asks
// for five bullets, evidence-report.md asks for the typed categories, and this
// fragment adds the ledger. A flat-prose ledger is not merely a fourth format —
// evidence-report.md forbids it, since anything that is not evidence must arrive
// labelled, with its author and its confidence.
func TestCompose_RolesToggle_LedgerLivesInsideTheTypedReport(t *testing.T) {
	body := rolesSection(t, "### the roles you did not convene")
	for _, phrase := range []string{"judgment", "confidence", "typed report"} {
		if !strings.Contains(body, phrase) {
			t.Errorf("the ledger section is missing %q: an untyped ledger next to the "+
				"report is the parallel format the evidence fragment already forbids", phrase)
		}
	}
}

// TestCompose_RolesToggle_HasAFloorTheOrchestratorCannotSkip asserts that some
// reviewers go in regardless of the orchestrator's judgment.
//
// A floor the orchestrator may skip when it feels confident is not a floor: it
// gets skipped precisely on the change where confidence was cheapest.
func TestCompose_RolesToggle_HasAFloorTheOrchestratorCannotSkip(t *testing.T) {
	body := rolesSection(t, "### the floor that is not delegated")
	for _, phrase := range []string{
		"authentication, payments or user data",
		"regardless of what the orchestrator",
		"not the orchestrator's to make",
	} {
		if !strings.Contains(body, phrase) {
			t.Errorf("the floor section is missing %q: without it, the decision of who "+
				"reviews auth and payments belongs to the same agent that wrote the code", phrase)
		}
	}
}

// TestCompose_RolesToggle_FloorDoesNotStandInForTheHumanGate asserts the floor
// says it is not a substitute for governance's CRITICAL approval.
//
// governance.md marks Auth, Billing, Security and Audit logs CRITICAL: no code
// written without explicit human approval. This floor convenes four reviewer
// sub-agents on an overlapping trigger, and four sub-agents can easily read as
// "the review requirement was met". A sub-agent is not a human.
func TestCompose_RolesToggle_FloorDoesNotStandInForTheHumanGate(t *testing.T) {
	body := rolesSection(t, "### the floor that is not delegated")
	for _, phrase := range []string{"governance", "not a human", "does not replace"} {
		if !strings.Contains(body, phrase) {
			t.Errorf("the floor section is missing %q: convening reviewers must not read "+
				"as having satisfied the human approval gate", phrase)
		}
	}
}

// TestCompose_RolesToggle_ReviewerNoEditIsAskedNotGuaranteed asserts the fragment
// tells the truth about what it can enforce.
//
// The only launch pattern this block documents is
// Agent({description, model, prompt}): three keys, none of them a tool list. So
// "you have no edit tools" can only travel inside the prompt — which is the
// instruction the earlier wording claimed it was not. The fragment must ask for
// the restraint, admit it is an instruction today, and name the missing lever.
func TestCompose_RolesToggle_ReviewerNoEditIsAskedNotGuaranteed(t *testing.T) {
	body := rolesSection(t, "### the reviewer should not get edit tools")
	for _, phrase := range []string{
		"is an instruction, not a guarantee",
		"none of them is a tool list",
		"what is missing",
	} {
		if !strings.Contains(body, phrase) {
			t.Errorf("the reviewer section is missing %q: a guarantee the launch pattern "+
				"cannot express is a claim, and claiming it hides that the lever is absent", phrase)
		}
	}
}

// TestCompose_RolesToggle_DoesNotClaimAToolLeverItLacks asserts the overclaim
// does not come back.
//
// This is the assert the first version of this file did not have: it checked the
// sentence was written, never that the mechanism behind it existed. The wording
// below is exactly what was there before, and it was false.
func TestCompose_RolesToggle_DoesNotClaimAToolLeverItLacks(t *testing.T) {
	body := rolesFragment(t)
	for _, banned := range []string{"not an instruction", "which tools it receives"} {
		if strings.Contains(body, banned) {
			t.Errorf("the roles fragment claims %q, but the documented launch is "+
				"Agent({description, model, prompt}) — three keys, none a tool list. "+
				"The sentence degrades into the prompt, i.e. into the instruction it "+
				"denies being", banned)
		}
	}
}

// rolesRowMatching returns the flattened row whose FIRST CELL contains want.
//
// Indexing by a fragment of the first cell is what lets a test assert the
// mapping instead of the vocabulary: the row keys here are descriptions of what
// a diff touches, and they are meant to be re-worded. What may not change is
// which reviewer the row lands on.
func rolesRowMatching(t *testing.T, want string) string {
	t.Helper()
	rows := rolesPhaseRows(t)
	for key, row := range rows {
		if strings.Contains(key, want) {
			return row
		}
	}
	t.Fatalf("no table row in the roles fragment describes a diff that touches %q", want)
	return ""
}

// TestCompose_RolesToggle_SizeIsReadOffTheDiff asserts the review size is
// DERIVED from what the change touches, not picked as a level by the agent that
// coordinated writing it.
//
// WHY THIS PROPERTY AND NOT THE OBVIOUS ONE. The first version of this rule hung
// the size off the CRITICAL / HIGH / MEDIUM / LOW ladder and asked the
// orchestrator to classify. Two holes, and the run of 2026-09-09 shows both at
// once:
//
//   - The classification was WRONG on the very run that motivated the rule. That
//     change was filed as harmless scaffolding, and it touched secrets in
//     environment variables, database migrations, the CI and the `.env`. Its
//     twelve real problems were real BECAUSE it touched those. A rule keyed on
//     the chosen level would have convened the small review on the one change we
//     know had twelve things broken in it.
//   - The classifier is the interested party. The agent that coordinated writing
//     the code grading its own change is the incentive pointing at "low".
//
// Reading it off the diff closes both with one change: nobody can file as small
// a change with a migration inside it, because the migration is right there.
//
// The datum that discriminates: the rows must MAP a diff fact to a reviewer. A
// section that merely says "it depends on what it touches" and stops defers to
// nothing, so this parses the table and checks where each row lands.
func TestCompose_RolesToggle_SizeIsReadOffTheDiff(t *testing.T) {
	body := rolesSection(t, "### how many reviewers")
	required := map[string]string{
		"seventeen agents":              "without the measured number the rule is a preference, and preferences get argued with instead of applied",
		"is not a level somebody picks": "the whole correction: a level the orchestrator chooses is a level the orchestrator gets wrong in its own favour",
		"it is read off the diff":       "the replacement mechanism has to be named, not implied",
		"migration is right there":      "the argument for why reading beats judging, and it is what a reader has to be able to repeat",
		"no matter how small the diff":  "the floor is what this rule may not touch, and it must say so where it is read",
		"when in doubt, you convene":    "a size rule with no default errs towards convening nobody, which is the failure it was written to avoid",
		"read this rule as a licence":   "it must name what it does NOT license: the auditor and the lenses were the parts that paid off",
	}
	for phrase, why := range required {
		if !strings.Contains(body, phrase) {
			t.Errorf("the size section is missing %q: %s", phrase, why)
		}
	}

	// El mecanismo: cada hecho del diff cae en un revisor concreto.
	mapeo := map[string]string{
		"migration": "dba",
		"secrets":   "floor",
		"ci":        "devops",
		"ui tree":   "frontend",
		"endpoint":  "backend",
	}
	for hecho, revisor := range mapeo {
		row := rolesRowMatching(t, hecho)
		if !strings.Contains(row, revisor) {
			t.Errorf("a diff touching %q does not land on %q (row: %q): a table that "+
				"names the fact and not the reviewer classifies nothing", hecho, revisor, row)
		}
	}

	// Y la clasificación NO puede volver a ser un nivel que alguien elige.
	for _, banned := range []string{
		"which level you used",
		"classifying by blast radius",
		"the change's governance level",
	} {
		if strings.Contains(body, banned) {
			t.Errorf("the size section still asks somebody to CHOOSE a level (%q). The "+
				"run of 2026-09-09 filed a change touching secrets, migrations and CI as "+
				"scaffolding: the classifier is the interested party", banned)
		}
	}
}

func TestCompose_RolesToggle_HighProposeIsReadBySeveralRolesBeforeSignature(t *testing.T) {
	body, err := config.Compose([]string{"roles", "governance"}, "claude")
	if err != nil {
		t.Fatal(err)
	}
	body = strings.ToLower(body)
	i := strings.Index(body, rolesHeading)
	if i < 0 {
		t.Fatal("falta roles")
	}
	body = body[i:]
	rows, _ := tableRows(body)
	row := rowText(rows["propose"])
	for _, frase := range []string{"where governance is active", "high", "critical", "at least three", "functional analyst", "qa", "auditor", "before human approval"} {
		if !strings.Contains(row, frase) {
			t.Errorf("fila propose no exige: %s", frase)
		}
	}
	for _, frase := range []string{
		"at least three independent specialist reports", "architect, functional analyst and qa", "one additional report for each specialty the scope requires", "one at a time", "exclude other specialists' reports and memories", "auditor reads all n original reports", "n >= 3", "n + 1", "before human approval", "do not replace the human signature", "## aprobación", "until all required reports exist", "identify which report is missing and why", "informes/", "leeme.md", "model/provider", "sha", "provenance", "full return or textual transcript", "contamination", "does not count as independent", "obtain a replacement", "reviewers receive no write", "not proof of isolation",
	} {
		if !strings.Contains(seniorFlat(body), frase) {
			t.Errorf("falta propose multirol: %s", frase)
		}
	}
	// La regla nueva debe vivir después del piso y antes del pedido de herramientas.
	piso := strings.Index(body, "### the floor that is not delegated")
	regla := strings.Index(body, "### high/critical propose before signature")
	tools := strings.Index(body, "### the reviewer should not get edit tools")
	if !(piso >= 0 && piso < regla && regla < tools) {
		t.Error("regla multirol fuera del lugar aprobado")
	}
}

func TestCompose_RolesToggle_HighProposeIsConditionalOnGovernance(t *testing.T) {
	body := seniorFlat(rolesFragment(t))
	for _, frase := range []string{"where governance is active and marks the change high or critical", "if governance is absent, do not invent a high/critical classification"} {
		if !strings.Contains(body, frase) {
			t.Errorf("falta condición: %s", frase)
		}
	}
}
