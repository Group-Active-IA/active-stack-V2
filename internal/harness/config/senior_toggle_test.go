package config_test

// ──────────────────────────────────────────────────────────────────────────────
// El senior: el único acompañante permanente, no una etapa del ciclo
//
// Todos los demás roles de este bloque están atados a un punto: el explorador
// explora, el revisor lee un diff, el archivista archiva. El senior se convoca
// en LAS CINCO fases, y lo que aporta no es una habilidad que a los otros les
// falte: es lo que pasó antes.
//
// LA DIFERENCIA CON LOS ADVERSARIALES, que es lo que justifica el lanzamiento
// extra: ellos tienen foco y no tienen memoria — ven un diff, buscan un bug,
// terminan, y arrancar limpios es lo que los hace independientes. El senior
// tiene memoria y no tiene foco. Por eso en verify NO busca bugs: juzga los
// hallazgos de los otros. Y no es el auditor: el auditor pregunta si el REPORTE
// se sostiene, el senior si el hallazgo IMPORTA en este proyecto.
//
// LO QUE ESTE ARCHIVO CUIDA, además de que las reglas estén escritas:
//
//   - LA TRAMPA TÉCNICA. Un sub-agente no puede frenar nada: devuelve su reporte
//     y termina. El único que puede parar el ciclo es el orquestador, al leer
//     ese reporte. El trabajo está partido — el senior ETIQUETA, el orquestador
//     FRENA — y el anti-assert prohíbe las dos formas de romperlo: atribuirle al
//     senior el freno o la espera.
//   - LOS DOS ANTI-ASSERTS, en las dos direcciones. Los fragmentos anteriores de
//     este repo afirmaban mecanismos inexistentes. La primera versión de ÉSTE
//     cometió el pecado invertido: NEGÓ un mecanismo que el bloque documenta.
//     Decía que engram no permite enumerar changes anteriores "sin saber de
//     antemano una palabra", y `engram-protocol.md` prescribe `mem_search
//     "opsx"` (línea 10) sobre guardados cuyo título lleva la constante "OPSX:"
//     (línea 11), y define `topic_key` estables NO scopeados a change (línea 35,
//     ejemplo `architecture/auth-model`). Pedía como faltante una convención que
//     el protocolo ya define. Las frases falsas están prohibidas por nombre.
//   - EL CANAL 1 NO ES UNA GARANTÍA. El prompt es el canal que MÁS depende de que
//     alguien se acuerde: la plantilla del bloque base tiene cinco slots, ninguno
//     es "qué pasó antes", y su slot de contexto dice explícitamente rutas y NO
//     contenido. Se entrega el bloque literal a pegar y se declara pedido, como
//     hizo knowledge-base.md.
//   - EL MAPEO FASE→PREGUNTA, y por qué se asserta la FILA ENTERA. La primera
//     versión asertaba las ocho preguntas sueltas: puse `verify` en las cinco
//     filas, destruyendo el mapeo completo, y no cayó ni un test. Los nombres de
//     las fases vivían sólo en el string de motivo de `requirePhrases`, que el
//     helper nunca mira. Es una tercera vía de assert flojo, distinta de la frase
//     duplicada y de la sección sin cobertura: EL DATO QUE DISCRIMINA VIVIENDO EN
//     EL MOTIVO Y NO EN EL ASSERT. Mismo molde, dos gemelos ya corregidos: el
//     test de verify nunca asertaba "verify" (cambiar la frase a `apply` quedaba
//     verde) y el de la ASSUMPTION nunca asertaba la ruta que la etiqueta
//     califica.
//   - LOS CHOQUES CON roles.md, que son de razones y no de listas. En apply
//     roles.md no convoca a nadie porque un diff que todavía se mueve se lee tres
//     veces y se juzga en ninguna de sus versiones; las preguntas de apply del
//     senior exigen leer ese diff, así que el fragmento tiene que decir que se
//     lanza al CIERRE de la fase. Y el senior no es el octavo rol: no lo elige lo
//     que el change toca, y el ledger de no convocados no lo cubre.
//   - EL COSTO. Cinco llamadas más por change, y es apagable.
//
// Cada test acota su sección y le saca el encabezado, por el bug que este repo
// ya cazó doce veces: un test que busca una frase que el título de su propia
// sección también escribe pasa con la sección vaciada.
// ──────────────────────────────────────────────────────────────────────────────

import (
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/catalog"
	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

// seniorHeading is the fragment's own heading, and it is the marker on purpose.
// The word "senior" appears nowhere else in any asset — base blocks included —
// so nothing but this fragment can satisfy it.
const seniorHeading = "## the senior"

// seniorFragment returns ONLY the senior fragment, never the whole block.
func seniorFragment(t *testing.T) string {
	t.Helper()
	// This helper inspects the full-cycle fragment; ordinary verify is covered separately.
	body, err := config.Compose([]string{"senior", "review", "auditor"}, "claude")
	if err != nil {
		t.Fatalf("Compose([senior]): %v", err)
	}
	lower := strings.ToLower(body)
	i := strings.Index(lower, seniorHeading)
	if i < 0 {
		t.Fatalf("the composed block does not carry the senior fragment (heading %q missing)", seniorHeading)
	}
	return lower[i:]
}

// seniorFlat collapses every run of whitespace to a single space, so a phrase
// that sits astride the 80-column wrap is still found as prose.
func seniorFlat(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// seniorSection returns the BODY of one "### " section with its own heading line
// dropped. Dropping the heading is the whole point: a test that matches a phrase
// its section title also writes stays green with the section emptied.
func seniorSection(t *testing.T, headingPrefix string) string {
	t.Helper()
	body := seniorFragment(t)
	i := strings.Index(body, headingPrefix)
	if i < 0 {
		t.Fatalf("the senior fragment has no section starting with %q", headingPrefix)
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
	return seniorFlat(rest)
}

// seniorPreamble returns the body of the "## " heading: everything from the
// fragment's own heading line to the first "### ".
func seniorPreamble(t *testing.T) string {
	t.Helper()
	body := seniorFragment(t)
	nl := strings.Index(body, "\n")
	if nl < 0 {
		t.Fatal("the senior fragment is a single line")
	}
	rest := body[nl+1:]
	if j := strings.Index(rest, "\n### "); j >= 0 {
		rest = rest[:j]
	}
	return seniorFlat(rest)
}

// TestCompose_SeniorToggle_IsKnown asserts every variant accepts it: the
// fragment is shared, so no variant may reject it.
func TestCompose_SeniorToggle_IsKnown(t *testing.T) {
	for _, v := range evidenceVariants {
		if _, err := config.Compose([]string{"senior"}, v); err != nil {
			t.Errorf("Compose([senior], %q): %v", v, err)
		}
	}
}

// TestCompose_SeniorToggle_Subtractive asserts a block that did not ask for it
// does not carry it.
func TestCompose_SeniorToggle_Subtractive(t *testing.T) {
	off, err := config.Compose([]string{"delegation", "model-routing", "engram", "tdd", "governance", "evidence", "auditor", "roles", "review", "knowledge-base"}, "claude")
	if err != nil {
		t.Fatalf("Compose without senior: %v", err)
	}
	if strings.Contains(strings.ToLower(off), seniorHeading) {
		t.Error("the senior fragment leaks into a block that did not ask for it")
	}
}

// TestCompose_SeniorToggle_ComposesAfterTheReviewFamily pins the position
// argued for in compose.go.
//
// Order is argued in a comment and fixed by nothing else. It composes after
// review and after knowledge-base because it is defined by contrast with both —
// a reader who meets the senior before meeting the adversarial passes has no
// way to see what it adds — and before engram, whose protocol it cites for the
// half of its memory that is not on disk.
func TestCompose_SeniorToggle_ComposesAfterTheReviewFamily(t *testing.T) {
	all, err := config.Compose(config.KnownToggles(), "claude")
	if err != nil {
		t.Fatalf("Compose(all): %v", err)
	}
	low := strings.ToLower(all)
	sr := strings.Index(low, seniorHeading)
	rv := strings.Index(low, reviewHeading)
	kb := strings.Index(low, kbHeading)
	// engram-protocol.md's own heading in full, NOT "## engram": the base block
	// writes that exact markdown heading inside the sub-agent prompt template, so
	// the short marker matches the BASE and the order assert compares against a
	// heading the fragment does not own.
	eng := strings.Index(low, "## engram persistent memory")
	if sr < 0 || rv < 0 || kb < 0 || eng < 0 {
		t.Fatalf("missing heading in the all-toggles block (senior=%d review=%d kb=%d engram=%d)", sr, rv, kb, eng)
	}
	if !(rv < sr && kb < sr && sr < eng) {
		t.Errorf("the senior fragment composes out of place (review=%d kb=%d senior=%d engram=%d). "+
			"It belongs after review and knowledge-base — it is defined by contrast with the "+
			"adversarial passes, so they must already be on the page — and before engram, "+
			"whose protocol it cites for the half of its memory that is not on disk", rv, kb, sr, eng)
	}
}

// TestCompose_SeniorToggle_PreambleSaysItIsEveryPhase asserts the opening keeps
// what makes this piece different from every other role in the block.
//
// Every other role is attached to one point in the cycle. A reader who takes the
// senior for one more phase-scoped reviewer schedules it once, in verify, and
// the four phases where it had something nobody else could say go unasked.
func TestCompose_SeniorToggle_PreambleSaysItIsEveryPhase(t *testing.T) {
	body := seniorPreamble(t)
	if !strings.Contains(body, "only in **verify**") {
		t.Errorf("falta alcance de senior: %s", "only in **verify**")
	}
	if !strings.Contains(body, "everything that happened before") {
		t.Errorf("falta alcance de senior: %s", "everything that happened before")
	}
}

// TestCompose_SeniorToggle_IsNotAFifthReviewer asserts the asymmetry that
// justifies paying for it.
//
// If it does the same as the adversarial passes, the extra launch buys nothing.
// The difference is not depth: the reviewers have focus and no memory, the
// senior has memory and no focus. Deleting that sentence leaves a section that
// reads like "another reviewer, but wiser", which is exactly the framing that
// gets it merged into the review passes at the first cost review.
func TestCompose_SeniorToggle_IsNotAFifthReviewer(t *testing.T) {
	body := seniorSection(t, "### it is not a fifth reviewer")
	requirePhrases(t, body, map[string]string{
		"focus and no memory": "what the adversarial passes have",
		"memory and no focus": "what the senior has, stated as the mirror image",
		"starts clean":        "why the reviewers have no memory: it is the thing that makes them independent",
		"none of them was there": "the reason the contradiction is unsayable by any of them, " +
			"which is what makes it worth an extra launch",
	}, "the not-a-fifth-reviewer section")
}

// TestCompose_SeniorToggle_InVerifyItJudgesTheFindings asserts the job, and
// asserts WHICH PHASE it belongs to.
//
// The phase name is in the assert and not only in this comment: the first
// version of this test named `verify` in its reason string alone, so changing
// the fragment to "so in `apply` the senior judges the findings" — moving the
// rule to the phase where roles.md convenes nobody — left it green. The datum
// that discriminates has to be inside the assert.
func TestCompose_SeniorToggle_InVerifyItJudgesTheFindings(t *testing.T) {
	body := seniorSection(t, "### it is not a fifth reviewer")
	requirePhrases(t, body, map[string]string{
		"so in `verify` the senior **does not hunt for bugs. it judges the findings**": "the rule WITH its phase: " +
			"the same sentence pointed at any other phase would be a different rule",
		"has no context to answer": "why the pass that raised the finding cannot judge it, which is the whole argument",
	}, "the not-a-fifth-reviewer section")
}

// TestCompose_SeniorToggle_DistinguishesItselfFromTheAuditor asserts the line
// against the OTHER fragment that reads the reviewers' reports.
//
// The auditor already reads the reports and already has a threshold that is this
// one's twin. Without the distinction written down, the honest reading is that
// one of the two is redundant — and the distinction is real: the auditor asks
// whether the report HOLDS UP, the senior whether the finding MATTERS HERE. A
// finding can be impeccably evidenced and still be a rule this project decided
// against, which is invisible to an auditor and is not what it is looking for.
func TestCompose_SeniorToggle_DistinguishesItselfFromTheAuditor(t *testing.T) {
	body := seniorSection(t, "### it is not a fifth reviewer")
	requirePhrases(t, body, map[string]string{
		"that is not the auditor's job": "the collision, named instead of left for the reader to notice",
		"whether they **hold up**":      "the auditor's question",
		"about **the report**":          "what the auditor's question is about",
		"about **the project**":         "what the senior's is about, in the same shape so the two can be compared",
		"impeccably evidenced":          "the case that separates them: it passes every audit and is still wrong here",
	}, "the not-a-fifth-reviewer section")
}

// TestCompose_SeniorToggle_NamesWhereTheMemoryComesFrom is the guard on the
// mechanism the whole fragment rests on.
//
// EL MECANISMO. The base block states that each sub-agent starts with no
// context, so "the senior has memory" is false as written unless the fragment
// says how the memory arrives — through three channels, NONE of them a
// guarantee.
//
// The first version called the prompt "the only guaranteed one". It is the
// opposite: the launch template's context slot is explicit that what travels is
// paths and NOT content, no slot carries the history, and nothing makes the
// orchestrator paste anything. It is the channel that depends most on somebody
// remembering. `roles.md` uses the true word for a channel of the same kind —
// an instruction, not a guarantee — and knowledge-base.md handled it by
// handing over the literal block to paste. So does this.
func TestCompose_SeniorToggle_NamesWhereTheMemoryComesFrom(t *testing.T) {
	body := seniorSection(t, "### where its memory actually comes from")
	requirePhrases(t, body, map[string]string{
		"each sub-agent starts with no context":              "the base block's own sentence, which is what makes the memory a problem to solve",
		"not one of them is a guarantee":                     "the ranking of the three channels, which the first version got backwards",
		"the history of this change is none of them":         "what the launch template's five slots do NOT include",
		"paths, not content":                                 "the template's own words about what its context slot carries",
		"when you launch the senior, add this to its prompt": "the imperative, without which this channel has no addressee",
		"an instruction, not a guarantee":                    "the honest footing, in roles.md's own words",
		"readable without any memory system at all":          "the archive directory, the channel that survives engram being off",
	}, "the memory section")
}

// TestCompose_SeniorToggle_LabelsTheArchivePathAsAssumption asserts the label
// AND the path it qualifies.
//
// The path is in the assert and not only in the reason: the first version
// asserted the ASSUMPTION label alone, so rewriting the directory convention to
// anything at all left the test green while every dated citation in the fragment
// went to a convention that does not exist. The label is worthless without the
// claim it labels.
func TestCompose_SeniorToggle_LabelsTheArchivePathAsAssumption(t *testing.T) {
	body := seniorSection(t, "### where its memory actually comes from")
	requirePhrases(t, body, map[string]string{
		"`openspec/changes/archive/yyyy-mm-dd-{name}/`":  "the path convention the whole dating story rests on",
		"read out of the base block, not exercised here": "the label attached to that specific claim",
		"dated off a bad clock":                          "what it costs if the convention is wrong: a confident citation with a wrong date",
	}, "the memory section")
}

// TestCompose_SeniorToggle_SaysWhatEngramActuallyOffers is the anti-assert's
// positive half.
//
// EL PECADO INVERTIDO. The first version declared that engram cannot enumerate
// prior changes — "you have to already know a word from the decision to find
// it" — and asked, as a missing feature, for a key not scoped to a change. All
// three claims are contradicted by `engram-protocol.md`, which composes two
// fragments later: it prescribes `mem_search "opsx"`, it saves every completed
// action under a title carrying that same constant, and its stable `topic_key`
// example is not scoped to a change at all. Denying a documented mechanism is
// the same defect as inventing one, pointed the other way: it sends the reader
// to build a workaround for a problem that was already solved.
func TestCompose_SeniorToggle_SaysWhatEngramActuallyOffers(t *testing.T) {
	body := seniorSection(t, "### where its memory actually comes from")
	requirePhrases(t, body, map[string]string{
		"write the prompt as if it may be off": "engram is a separate toggle, so nothing may depend on it",
		"fixed word":                           "why no prior knowledge is needed: the search term is a constant",
		"**`mem_search \"opsx\"` enumerates the changes that came before**": "the capability, in a form somebody can run",
		"is not scoped to a change at all":                                  "the stable topic_key, which is where a decision that outlives one change lives",
		"all of it is asked for and none of it is enforced":                 "the REAL limit, which is the footing and not the mechanism",
		"\"nothing was saved\", never \"nothing was decided\"":              "how an empty result must be read, so silence never becomes an all-clear",
	}, "the memory section")
}

// TestCompose_SeniorToggle_AsksADifferentQuestionInEachPhase asserts the five
// phases each get their own question — and asserts the PAIRING, not the
// questions.
//
// EL CASO TRECE. The first version asserted the eight questions as loose
// phrases. Setting every row of the table to `verify`, destroying the entire
// phase-to-question mapping, turned nothing red: the phase names lived only in
// the reason strings, which `requirePhrases` never reads. So the whole table row
// is the assert. A question attached to the wrong phase is not a smaller
// mistake than a missing question — in `apply` it is the difference between
// reading a diff that stopped moving and one that has not.
func TestCompose_SeniorToggle_AsksADifferentQuestionInEachPhase(t *testing.T) {
	body := seniorSection(t, "### what it looks at")
	if !strings.Contains(body, "| verify | are the findings **real or noise**? **what was not tested**? |") {
		t.Errorf("falta alcance de senior: %s", "| verify | are the findings **real or noise**? **what was not tested**? |")
	}
	if !strings.Contains(body, "dated decisions") {
		t.Errorf("falta alcance de senior: %s", "dated decisions")
	}
}

// TestCompose_SeniorToggle_RunsAtTheCloseOfEachPhase asserts the clash with
// roles.md is resolved instead of ignored.
//
// It is a clash of REASONS, not of lists. `roles.md` convenes nobody during
// apply because a diff that is still moving gets read three times and judged on
// none of its versions — and the senior's two apply questions require reading
// exactly that diff. The reason applies unchanged, so the launch goes at the
// close of the phase, on what the phase returned. That is also the only point
// where declining to open the next phase is free.
func TestCompose_SeniorToggle_RunsAtTheCloseOfEachPhase(t *testing.T) {
	body := seniorSection(t, "### what it looks at")
	if !strings.Contains(body, "after the breaker returns and before the auditor") {
		t.Errorf("falta alcance de senior: %s", "after the breaker returns and before the auditor")
	}
	if !strings.Contains(body, "quiet diff") {
		t.Errorf("falta alcance de senior: %s", "quiet diff")
	}
	if !strings.Contains(body, "never alongside a phase still running") {
		t.Errorf("falta alcance de senior: %s", "never alongside a phase still running")
	}
}

// TestCompose_SeniorToggle_IsNotAnEighthRole asserts the second clash with
// roles.md is named.
//
// `roles.md` has seven roles, a convening table driven by what the change
// touches, and a ledger of every role NOT convened. The senior fits none of
// that: nothing about the change decides whether it is called. Filed as an
// eighth role it lands in a table whose only criterion never applies to it, and
// the ledger silently under-reports forever.
func TestCompose_SeniorToggle_IsNotAnEighthRole(t *testing.T) {
	body := seniorSection(t, "### what it looks at")
	if !strings.Contains(body, "not one of the seven role specialists") {
		t.Errorf("falta alcance de senior: %s", "not one of the seven role specialists")
	}
	if !strings.Contains(body, "the ledger of roles you did not convene does not account for it") {
		t.Errorf("falta alcance de senior: %s", "the ledger of roles you did not convene does not account for it")
	}
	if !strings.Contains(body, "once per verify round when active") {
		t.Errorf("falta alcance de senior: %s", "once per verify round when active")
	}
}

// TestCompose_SeniorToggle_TheStopIsAddressedToTheOrchestrator is the anti-trap.
//
// LA TRAMPA TÉCNICA. A sub-agent cannot stop anything: it returns its report and
// terminates, and nothing about it is still running when the report is read. The
// only thing that can decline to launch the next phase is the orchestrator. A
// fragment that writes the stop as a capability of the senior describes a
// mechanism that does not exist, and this repo has already had to correct that
// exact error in other fragments.
func TestCompose_SeniorToggle_TheStopIsAddressedToTheOrchestrator(t *testing.T) {
	body := seniorSection(t, "### it stops nothing")
	requirePhrases(t, body, map[string]string{
		"a sub-agent cannot stop anything":          "the mechanical fact the rest of the section rests on",
		"you, the orchestrator":                     "the addressee, named — without it the instruction has no executor",
		"do not start the next phase":               "the instruction itself, in the imperative and aimed at who can obey it",
		"describes a mechanism that does not exist": "why it is not written the other way round",
	}, "the stop section")
}

// TestCompose_SeniorToggle_SplitsLabellingFromHolding asserts the division of
// labour that keeps the threshold executable.
//
// The senior can label a finding; only the orchestrator can hold the cycle. Left
// as one verb, the threshold reads as something the senior does, which is the
// mechanism that does not exist — the trap re-entering through the section that
// comes after the one that closed it.
func TestCompose_SeniorToggle_SplitsLabellingFromHolding(t *testing.T) {
	body := seniorSection(t, "### the threshold")
	requirePhrases(t, body, map[string]string{
		"the senior supplies context; the auditor classifies and decides": "the split, in one line",
		"the orchestrator executes the auditor's verdict":                 "where the stop actually happens",
		"blocking or note": "the two labels the senior may attach, which is the whole of its half",
	}, "the threshold section")
}

// TestCompose_SeniorToggle_HasAThreshold asserts the cut exists.
//
// EL UMBRAL. An agent asked to analyse and question will always find something.
// With no threshold the cycle is held six times per change and inside a week the
// toggle is switched off — and a senior that is off advises nobody, which is
// strictly worse than one that only ever spoke when it mattered.
func TestCompose_SeniorToggle_HasAThreshold(t *testing.T) {
	body := seniorSection(t, "### the threshold")
	requirePhrases(t, body, map[string]string{
		"will always find something": "why a threshold is required at all",
		"switched off":               "what happens without one",
		"hold the cycle when the report carries something the person does not know": "the line, in one sentence",
		"already in plain sight": "the other half of the line: an opinion about the visible is not a stop",
		"the auditor's threshold word for word": "the collision with the other fragment that has this rule, " +
			"acknowledged rather than presented as new",
	}, "the threshold section")
}

// TestCompose_SeniorToggle_ListsTheFourStopConditions asserts each blocking
// condition survives on its own.
//
// They are four and they are specific on purpose. Generalised into "stop when
// something important comes up" the threshold is back to the agent's judgement,
// which is what it was written to replace.
func TestCompose_SeniorToggle_ListsTheFourStopConditions(t *testing.T) {
	body := seniorSection(t, "### the threshold")
	requirePhrases(t, body, map[string]string{
		"contradicts a decision already taken":             "the first condition",
		"cites it with its date":                           "and the proof it must carry, without which it is a suspicion",
		"something that already exists is being rewritten": "the second condition",
		"the test just written proves nothing":             "the third condition",
		"nobody asked to touch":                            "the fourth condition: the critical thing the change reached on its way past",
	}, "the threshold section")
}

// TestCompose_SeniorToggle_NoteAndCarryOnIsTheOtherHalf asserts the
// non-blocking label is written out.
//
// Without it every observation is a candidate for holding the cycle, because "is
// this worth interrupting for" has no default answer. Style, naming and "you
// could consider" are the bulk of what any reviewer produces.
func TestCompose_SeniorToggle_NoteAndCarryOnIsTheOtherHalf(t *testing.T) {
	body := seniorSection(t, "### the threshold")
	requirePhrases(t, body, map[string]string{
		"**note** — it stays in the report": "the second label and what happens to it",
		"the phase continues":               "what happens to the work while the note is recorded",
		"you could consider":                "the tell that sorts an observation into the non-blocking half",
	}, "the threshold section")
}

// TestCompose_SeniorToggle_StopsWithDatumOptionsAndDefault asserts the three
// conditions that make an interruption cheaper than the problem.
//
// A stop that arrives as "this looks off — what do you want to do?" hands the
// entire problem back and costs more than the defect. The datum is what
// separates a finding from a hunch; the options are what turn the interruption
// into a decision; the safe default is what stops an unanswered question from
// becoming permission — and that last one is the ORCHESTRATOR's, because by then
// the senior has terminated.
func TestCompose_SeniorToggle_StopsWithDatumOptionsAndDefault(t *testing.T) {
	body := seniorSection(t, "### what a blocking finding has to carry")
	requirePhrases(t, body, map[string]string{
		"the datum, not the suspicion":               "the first condition",
		"separates the claim from its alternative":   "what makes a datum a datum and not a confirmation",
		"a hunch is a note":                          "where an observation without that datum belongs",
		"hands the whole problem back":               "what an open question costs the person",
		"hands back a choice":                        "what the options buy instead",
		"a safe default, and this one is yours":      "the third condition, with its owner attached",
		"the senior has already terminated":          "why that one cannot be the senior's",
		"nothing advances":                           "what the default IS: silence never becomes progress",
		"the interruption converted into permission": "the failure the safe default prevents",
	}, "the blocking-finding section")
}

// TestCompose_SeniorToggle_DeclaresItsCost asserts the price is on the page.
//
// A fragment that does not say what it costs gets evaluated against a price of
// zero, and then the bill arrives as latency nobody budgeted for. It is one more
// sub-agent per phase, on every change, including the ones where the senior has
// nothing to say — and it is a toggle, so the answer to a report nobody reads is
// to turn it off, not to leave it running and ignored.
func TestCompose_SeniorToggle_DeclaresItsCost(t *testing.T) {
	body := seniorSection(t, "### what it costs")
	if !strings.Contains(body, "one more sub-agent per verify round") {
		t.Errorf("falta alcance de senior: %s", "one more sub-agent per verify round")
	}
	if !strings.Contains(body, "at most three per change") {
		t.Errorf("falta alcance de senior: %s", "at most three per change")
	}
	if !strings.Contains(body, "nothing to say") {
		t.Errorf("falta alcance de senior: %s", "nothing to say")
	}
	if !strings.Contains(body, "turn it off") {
		t.Errorf("falta alcance de senior: %s", "turn it off")
	}
}

// TestCompose_SeniorToggle_DoesNotClaimAMechanismItLacks is the anti-assert, and
// it runs in BOTH directions.
//
// Fragments in this repo have already shipped guarantees that did not exist. The
// first version of THIS one shipped the mirror image: it denied a mechanism the
// block documents. Both are banned by name, because both send the reader to act
// on a map of the system that is wrong.
func TestCompose_SeniorToggle_DoesNotClaimAMechanismItLacks(t *testing.T) {
	body := seniorFlat(seniorFragment(t))

	// Direction one: powers the senior does not have. A sub-agent returns its
	// report and terminates — it cannot block, wait, or hold anything.
	for _, banned := range []string{
		"the senior blocks",
		"the senior halts",
		"the senior stops",
		"the senior waits",
		"the senior pauses",
		"stop and wait",
		"blocks the phase",
		"halts the cycle",
		"pauses the cycle",
		"automatically pauses",
		"automatically stops",
		"holds a veto",
		"has veto",
		"cannot be overridden",
		"the only guaranteed",
		"remembers every previous change",
		"has access to every previous change",
		"full history of the project",
	} {
		if strings.Contains(body, banned) {
			t.Errorf("the senior fragment claims %q. Nothing here is a mechanism: a sub-agent "+
				"returns its report and terminates, so the only thing that can decline to launch "+
				"the next phase is the orchestrator, and the memory is whatever the prompt carries "+
				"plus what the agent reads off disk. Stating it as a capability of the senior is "+
				"the error this repo has already had to correct in other fragments", banned)
		}
	}

	// Direction two: mechanisms the block DOES document. Denying one is the same
	// defect inverted — it sends the reader to build a workaround for a solved
	// problem, and to file as missing a convention that is already prescribed.
	for _, denied := range []string{
		"already know a word from the decision",
		"nothing in the protocol lists keys by prefix",
		"something like `opsx/decisions`",
		"no way to enumerate",
		"cannot enumerate",
	} {
		if strings.Contains(body, denied) {
			t.Errorf("the senior fragment denies a documented mechanism with %q. engram-protocol.md, "+
				"which composes two fragments later, prescribes `mem_search \"opsx\"`, saves every "+
				"completed OPSX action under a title carrying that same constant, and gives stable "+
				"`topic_key`s that are not scoped to a change (its example is `architecture/auth-model`). "+
				"Denying a mechanism that exists is the same error as inventing one that does not", denied)
		}
	}
}

// TestCatalog_OffersTheSeniorToggle asserts the fragment is reachable.
//
// A fragment that composes but that the catalog never advertises is dead code
// with tests: nothing in the product can turn it on.
func TestCatalog_OffersTheSeniorToggle(t *testing.T) {
	toggles := orchestratorToggles(t)
	for _, tg := range toggles {
		if tg == "senior" {
			return
		}
	}
	t.Errorf("the catalog does not offer the `senior` toggle (offers: %s): the "+
		"fragment composes and no user can ask for it", strings.Join(toggles, ", "))
}

func TestSeniorCatalogoDescribeVerify(t *testing.T) {
	cat, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	h, ok := cat.ByID("sdd-orchestrator")
	if !ok {
		t.Fatal("falta orquestador")
	}
	if !strings.Contains(h.LongDescription["es"], "senior sólo en verify") {
		t.Error("catálogo no describe senior sólo en verify")
	}
}
