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
//   - La fase `verify` NO EXISTE en este bloque. El ciclo documentado es
//     explore → propose → apply → archive; "verify" aparece una sola vez en el
//     base y es la fila "Read 1-3 files to decide/verify" de la tabla de
//     delegación. Asignarle el piso de revisión a una fase que el orquestador no
//     puede convocar es un piso que no corre nunca.
//   - El único lanzamiento que el bloque documenta es
//     Agent({description, model, prompt}). Tres claves, ninguna es una lista de
//     herramientas. "No es una instrucción, es qué herramientas recibe" era
//     falso: la frase sólo puede viajar dentro del prompt, o sea, como la
//     instrucción que decía no ser.
//
// Por eso varios de estos tests asertan el MECANISMO y no el vocabulario, y uno
// asserta que la afirmación de más NO vuelva.
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
func rolesPhaseRows(t *testing.T) map[string]string {
	t.Helper()
	rows := map[string]string{}
	for _, line := range strings.Split(rolesFragment(t), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		if len(cells) < 2 {
			continue
		}
		phase := strings.TrimSpace(cells[0])
		rows[phase] = strings.Join(cells[1:], " ")
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
	// The cycle this block documents: explore → propose → apply → archive.
	for _, phase := range []string{"explore", "propose", "apply", "archive"} {
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

// TestCompose_RolesToggle_DoesNotStaffAPhaseThatDoesNotExist asserts the fragment
// does not hand the mandatory review to `verify`.
//
// `verify` is not a phase of this block: there is no /opsx:verify command, no
// verify skill, and no verify row in the model table. Its single appearance in
// the base is "Read 1-3 files to decide/verify" inside the delegation table.
// Assigning the floor there is a floor that never runs, so the row must say so
// in those words and the reviewer pass must live in a phase that exists.
func TestCompose_RolesToggle_DoesNotStaffAPhaseThatDoesNotExist(t *testing.T) {
	rows := rolesPhaseRows(t)
	if row, ok := rows["verify"]; ok {
		if !strings.Contains(row, "does not exist in this block") {
			t.Error("the phase table keeps a `verify` row without saying that the phase " +
				"does not exist in this block: work assigned there is never convened")
		}
		for _, role := range theSevenRoles {
			if strings.Contains(row, role) {
				t.Errorf("the `verify` row assigns %q to a phase the orchestrator cannot "+
					"convene: that reviewer never runs", role)
			}
		}
	}
	apply, ok := rows["apply"]
	if !ok {
		t.Fatal("the phase table has no `apply` row")
	}
	if !strings.Contains(apply, "qa") {
		t.Error("the mandatory QA is not in `apply`, the last phase this block can " +
			"actually convene: it was parked in a phase that does not exist")
	}
	if !strings.Contains(apply, "reviewer") {
		t.Error("the reviewer pass is not in `apply`: the whole no-edit-tools regime " +
			"hangs off a phase that is never reached")
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
