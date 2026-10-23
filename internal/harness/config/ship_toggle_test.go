package config_test

// ──────────────────────────────────────────────────────────────────────────────
// Ship: el único verbo del ciclo cuyo resultado lo ve otra gente
//
// Todo lo demás en este bloque pasa adentro de la máquina: explore lee, propose
// escribe un plan, apply escribe código, verify lo juzga, archive mueve un
// directorio. Todos se deshacen editando un archivo. Un push no: cuando está mal
// ya lo bajó alguien más, y el force-push encima no es un undo, es un segundo
// evento arriba del primero.
//
// Y NO ES UNA SECUENCIA DE PASOS: ES UN DIAGNÓSTICO. No hay lista fija de
// comandos porque la respuesta correcta depende de dónde está parado el
// repositorio. Seis estados, seis respuestas, y cuatro de las seis NO terminan
// en un push.
//
// LO QUE ESTE ARCHIVO CUIDA, además de que las reglas estén escritas:
//
//   - LOS SEIS ESTADOS SE ASSERTAN COMO FILA ENTERA. Es el caso trece de este
//     repo, cazado ayer en el test del senior: si el dato que discrimina (el
//     estado) vive en el string de motivo y no en el assert, poner la misma
//     respuesta en las seis filas queda verde. Acá la fila entera —estado Y
//     respuesta— es la frase buscada.
//   - LOS DOS ANTI-ASSERTS, en las dos direcciones. Este repo lleva cuatro
//     fragmentos corregidos por afirmar mecanismos inexistentes, y uno por negar
//     uno documentado. Dirección uno: nada acá frena, garantiza ni espera —un
//     sub-agente devuelve su reporte y termina. Dirección dos: el bloque base SÍ
//     nombra `Bash for state (git, openspec status)`, SÍ tiene comandos que
//     reportan el estado del change, y SÍ tiene fase verify (commit e02c36f).
//     Negar cualquiera de las tres es el mismo defecto al revés.
//   - `gh` NO APARECE EN NINGÚN ASSET. Verificado con rg sobre
//     internal/harness/config/assets/: ni el base ni ningún fragmento lo nombra.
//     Va como ASSUMPTION declarada CON su consecuencia, y la consecuencia es la
//     que importa: un `gh pr view` que falla se parece exactamente a "no hay PR",
//     que es justo la respuesta que abre el segundo PR.
//   - LA CONFIRMACIÓN VA DIRIGIDA AL ORQUESTADOR. Un sub-agente que pregunta no
//     recibe respuesta: devuelve la pregunta y termina. Mismo reparto que
//     `senior.md` ("the senior labels, you hold"), citado por nombre para que el
//     lector vea que es el mismo mecanismo y no uno nuevo.
//   - EL COSTO, y que es apagable.
//
// Cada test acota su sección y le saca el encabezado, por el bug que este repo
// ya cazó doce veces: un test que busca una frase que el título de su propia
// sección también escribe pasa con la sección vaciada.
// ──────────────────────────────────────────────────────────────────────────────

import (
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

// shipHeading is the fragment's own heading, and it is the marker on purpose.
// The bare word "ship" IS written elsewhere — roles.md ("before it ships"),
// senior.md ("have already shipped"), engram-protocol.md ("the CLI already
// ships") — so the marker has to be the heading and not the word.
const shipHeading = "## ship"

// shipFragment returns ONLY the ship fragment, never the whole block.
func shipFragment(t *testing.T) string {
	t.Helper()
	body, err := config.Compose([]string{"ship"}, "claude")
	if err != nil {
		t.Fatalf("Compose([ship]): %v", err)
	}
	lower := strings.ToLower(body)
	i := strings.Index(lower, shipHeading)
	if i < 0 {
		t.Fatalf("the composed block does not carry the ship fragment (heading %q missing)", shipHeading)
	}
	return lower[i:]
}

// shipFlat collapses every run of whitespace to a single space, so a phrase that
// sits astride the 80-column wrap is still found as prose.
func shipFlat(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// shipNoise strips the markdown a writer adds without changing the claim:
// backticks, emphasis markers and quotes. It exists for the anti-assert, and the
// reason is EL CASO CATORCE, found by the coordinator on the first version of
// this file.
//
// The banned list said `"gh is guaranteed"` — the tool named bare. But this
// fragment writes `gh` inside backticks four of the five times it names it,
// because it is a tool and tools go as code. So "`gh` is guaranteed to be
// installed" — the sentence somebody editing this fragment IN ITS OWN STYLE
// would actually write — walked straight past the guard, while the bare form
// nobody would write was the only one caught.
//
// The lesson generalises past this file: an anti-assert that matches raw
// markdown is defeated by the document's own formatting conventions. Compare
// against the claim, not against its typography.
var shipNoise = strings.NewReplacer("`", "", "*", "", "_", "", "\"", "", "'", "", "\u201c", "", "\u201d", "", "\u2019", "")

// shipNormalize flattens the whitespace AND strips the markdown noise, so a
// phrase is found however it was dressed.
func shipNormalize(s string) string {
	return shipFlat(shipNoise.Replace(s))
}

// requireShipPhrases is requirePhrases with BOTH SIDES normalised, and it exists
// for the positive half of the same lesson as the anti-assert.
//
// Measured on the first version of this file: of the 46 asserted phrases written
// without markdown, 42 broke when a single word inside them was bolded in the
// fragment. That is a false RED — the rule is intact, the sentence is on the
// page, and the failure says "is missing X" about a sentence that is present.
// Both directions of the same defect: an assert that matches typography reports
// on the typography.
//
// Normalising both sides removes the coupling without weakening anything, because
// in this fragment no assertion IS its emphasis: the table rows assert the
// state-to-answer pairing and keep their `|`, the addressee is asserted by the
// word "you" and not by the bold around it, and the ASSUMPTION label is asserted
// with its parenthetical so prose cannot satisfy it by accident. Phrases below
// may therefore still be written with their markdown — it is stripped from both
// sides — which keeps them readable next to the fragment they quote.
func requireShipPhrases(t *testing.T, body string, required map[string]string, what string) {
	t.Helper()
	norm := shipNormalize(body)
	for phrase, why := range required {
		if !strings.Contains(norm, shipNormalize(phrase)) {
			t.Errorf("%s is missing %q — %s", what, phrase, why)
		}
	}
}

// shipSection returns the BODY of one "### " section with its own heading line
// dropped. Dropping the heading is the whole point: a test that matches a phrase
// its section title also writes stays green with the section emptied.
func shipSection(t *testing.T, headingPrefix string) string {
	t.Helper()
	body := shipFragment(t)
	i := strings.Index(body, headingPrefix)
	if i < 0 {
		t.Fatalf("the ship fragment has no section starting with %q", headingPrefix)
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
	return shipFlat(rest)
}

// shipPreamble returns the body of the "## " heading: everything from the
// fragment's own heading line to the first "### ".
func shipPreamble(t *testing.T) string {
	t.Helper()
	body := shipFragment(t)
	nl := strings.Index(body, "\n")
	if nl < 0 {
		t.Fatal("the ship fragment is a single line")
	}
	rest := body[nl+1:]
	if j := strings.Index(rest, "\n### "); j >= 0 {
		rest = rest[:j]
	}
	return shipFlat(rest)
}

// TestCompose_ShipToggle_IsKnown asserts every variant accepts it: the fragment
// is shared, so no variant may reject it.
func TestCompose_ShipToggle_IsKnown(t *testing.T) {
	for _, v := range evidenceVariants {
		if _, err := config.Compose([]string{"ship"}, v); err != nil {
			t.Errorf("Compose([ship], %q): %v", v, err)
		}
	}
}

// TestCompose_ShipToggle_Subtractive asserts a block that did not ask for it
// does not carry it.
func TestCompose_ShipToggle_Subtractive(t *testing.T) {
	off, err := config.Compose([]string{"delegation", "model-routing", "engram", "tdd", "governance", "evidence", "auditor", "roles", "review", "knowledge-base", "senior"}, "claude")
	if err != nil {
		t.Fatalf("Compose without ship: %v", err)
	}
	if strings.Contains(strings.ToLower(off), shipHeading) {
		t.Error("the ship fragment leaks into a block that did not ask for it")
	}
}

// TestCompose_ShipToggle_ComposesLastOfTheAdditiveFamily pins the position
// argued for in compose.go.
//
// Order is argued in a comment and fixed by nothing else. Ship composes after
// review and after the senior because it borrows both: its preconditions name
// the verdict the review family produces, and its division of labour is the
// senior's "the senior labels, you hold" under another name — a reader who meets
// the borrowed rule before the original meets it without its argument. And
// before engram, which is a separate toggle nothing here may assume is on.
func TestCompose_ShipToggle_ComposesLastOfTheAdditiveFamily(t *testing.T) {
	all, err := config.Compose(config.KnownToggles(), "claude")
	if err != nil {
		t.Fatalf("Compose(all): %v", err)
	}
	low := strings.ToLower(all)
	sp := strings.Index(low, shipHeading)
	rv := strings.Index(low, reviewHeading)
	sr := strings.Index(low, seniorHeading)
	// engram-protocol.md's own heading in full, NOT "## engram": the base block
	// writes that exact markdown heading inside the sub-agent prompt template, so
	// the short marker matches the BASE and the order assert compares against a
	// heading the fragment does not own.
	eng := strings.Index(low, "## engram persistent memory")
	if sp < 0 || rv < 0 || sr < 0 || eng < 0 {
		t.Fatalf("missing heading in the all-toggles block (ship=%d review=%d senior=%d engram=%d)", sp, rv, sr, eng)
	}
	if !(rv < sp && sr < sp && sp < eng) {
		t.Errorf("the ship fragment composes out of place (review=%d senior=%d ship=%d engram=%d). "+
			"It belongs after review and after the senior — it cites the verdict one produces and "+
			"reuses the split the other names — and before engram, which is a separate toggle it "+
			"may not read as already assumed", rv, sr, sp, eng)
	}
}

// TestCompose_ShipToggle_PreambleSaysItIsADiagnosis asserts the opening keeps the
// two things that make this fragment different from every other one in the block.
//
// A reader who takes ship for a checklist runs the checklist: push, open PR,
// done. The whole value is that there is no fixed list — the answer depends on
// where the repository is standing. And a reader who takes a push for one more
// reversible step treats a mistake here as costing what a bad edit costs.
func TestCompose_ShipToggle_PreambleSaysItIsADiagnosis(t *testing.T) {
	body := shipPreamble(t)
	requireShipPhrases(t, body, map[string]string{
		"not a sequence of steps. it is a diagnosis": "the defining property, stated before any rule",
		"a push leaves the machine":                  "what makes this verb different from the other five",
		"is not an undo":                             "why a bad push is not repaired the way a bad edit is",
		"look first, decide second":                  "the instruction the diagnosis framing exists to produce",
	}, "the fragment preamble")
}

// TestCompose_ShipToggle_ListsTheFourPreconditions asserts each one survives on
// its own, and asserts what each one is CHECKED WITH.
//
// A precondition with no check beside it is a wish. The four are specific on
// purpose: generalised into "make sure everything is in order" they are back to
// the agent's judgement, which is what they were written to replace.
func TestCompose_ShipToggle_ListsTheFourPreconditions(t *testing.T) {
	body := shipSection(t, "### the four preconditions")
	requireShipPhrases(t, body, map[string]string{
		"the change is closed":              "the first precondition",
		"`openspec list --json`":            "and what answers it, which is the CLI this block already calls its eyes",
		"the suite is green on this commit": "the second precondition, with the part that makes it checkable",
		"run, not remembered":               "the distinction the second precondition exists for",
		"`git rev-parse head`":              "the SHA that turns a remembered green into an attributable one",
		"the tree is clean":                 "the third precondition",
		"`git status --porcelain`":          "and what answers it",
		"verdict is on the record":          "the fourth precondition",
		"but a missing one does":            "the asymmetry: a failing verdict is the person's call, an absent one is not",
	}, "the preconditions section")
}

// TestCompose_ShipToggle_StopsOnAMissingPrecondition asserts the failure mode is
// a stop and not a note.
//
// "Warn and carry on" is the shape every one of these degrades into, and it is
// indistinguishable from no precondition at all: the push happened either way and
// the warning is the line nobody read. The report has to name WHICH one failed,
// because "preconditions not met" sends the person to re-derive what you already
// know.
func TestCompose_ShipToggle_StopsOnAMissingPrecondition(t *testing.T) {
	body := shipSection(t, "### the four preconditions")
	requireShipPhrases(t, body, map[string]string{
		"name which one and stop":                    "the instruction, with the naming that makes it actionable",
		"a warning followed by the push is the push": "why 'warn and continue' is not a weaker version of this rule but the absence of it",
	}, "the preconditions section")
}

// TestCompose_ShipToggle_IsHonestAboutWhereTheVerdictLives asserts the fourth
// precondition does not invent a filename.
//
// review.md names `review-verdict.md` for the ADVERSARIAL passes, and that is a
// different artifact from `/opsx:verify`'s verdict; nothing in the block above
// names a file that verify writes. Hardcoding one would send the reader to look
// for a path that may not exist and read its absence as "verify never ran".
func TestCompose_ShipToggle_IsHonestAboutWhereTheVerdictLives(t *testing.T) {
	body := shipSection(t, "### the four preconditions")
	requireShipPhrases(t, body, map[string]string{
		"nothing in the block above names a file that `/opsx:verify` writes": "the honest gap, instead of a path invented to fill it",
		"`openspec status --change": "what does report artifact state, which is the CLI this block already trusts for exactly this",
		"only in this conversation": "the fallback case, which is the ordinary one",
		"is a weaker claim":         "and how it must be reported, so a formality never passes for an artifact",
	}, "the preconditions section")
}

// TestCompose_ShipToggle_DiagnosesSixStates asserts every state gets its own
// answer — and asserts the PAIRING, not the answers.
//
// EL CASO TRECE, otra vez. Asserting the six answers as loose phrases lets you
// point every row at the same state and stay green, because the state names
// would live only in the reason strings, which `requirePhrases` never reads. So
// the whole table row is the assert. An answer attached to the wrong state is
// not a smaller mistake than a missing answer: "push -u and propose a PR" is
// correct on a fresh branch and is a published mistake on the trunk.
func TestCompose_ShipToggle_DiagnosesSixStates(t *testing.T) {
	body := shipSection(t, "### the six states")
	requireShipPhrases(t, body, map[string]string{
		"| the work was committed straight onto `main` | **stop.** propose a retroactive branch and move the commits |":            "worked on the trunk: the only row whose first word is stop",
		"| a new branch with no remote yet | `git push -u`, then propose a pr |":                                                   "new branch, no remote: the row where a push is the right answer",
		"| a branch whose pr is already open | push to it and let that pr update. **never open a second one** |":                   "PR already open: push, and the prohibition that makes it checkable",
		"| `main` has moved on since the branch left it | rebase or merge first, and say in the report which of the two you did |": "the trunk advanced: reconcile before anything, and declare which way",
		"| the tree still has uncommitted changes | ask whether they belong in this change. **never decide that alone** |":         "uncommitted work: a question, never a judgement call",
		"| no remote is configured | say so and stop. **do not invent one** |":                                                     "no remote: the row that ends in nothing at all",
	}, "the six-states section")
}

// TestCompose_ShipToggle_TheMainStateChecksWhetherItAlreadyShipped asserts the
// hardest row carries the two things that make its answer safe.
//
// The retroactive branch rewrites the person's local `main`, which is not
// something to do on your own initiative — and it may not even be available:
// if local `main` is not ahead of `origin/main` those commits are already
// published, and no local surgery unpublishes them. A fragment that skips that
// check proposes moving commits that are already gone.
func TestCompose_ShipToggle_TheMainStateChecksWhetherItAlreadyShipped(t *testing.T) {
	body := shipSection(t, "### the six states")
	requireShipPhrases(t, body, map[string]string{
		"rewrites the person's local `main`":  "why the fix is proposed and not run",
		"not ahead of `origin/main`":          "the check that says whether the fix is even available",
		"no local surgery unpublishes them":   "what is true once they are out, stated so the report does not promise a rollback",
		"splits the review across two places": "why a second PR is worse than an awkward first one",
	}, "the six-states section")
}

// TestCompose_ShipToggle_NamesTheReadBehindEachState asserts the diagnosis is
// executable and not a mood.
//
// Six states with no commands behind them is a taxonomy, and the agent invents
// its own reads — including hardcoding `main` as the trunk name, which
// misdiagnoses every repository whose default branch is called something else.
func TestCompose_ShipToggle_NamesTheReadBehindEachState(t *testing.T) {
	body := shipSection(t, "### the reads")
	requireShipPhrases(t, body, map[string]string{
		"`git branch --show-current`":                         "which branch you are on",
		"`git symbolic-ref --short refs/remotes/origin/head`": "what the trunk is actually called on this repository",
		"the default branch is not always called `main`":      "the reason that read exists rather than a hardcoded name",
		"`git remote -v`":                                     "whether there is a remote at all",
		"`git rev-parse --abbrev-ref @{upstream}`":            "whether this branch tracks anything",
		"exits non-zero":                                      "how that read answers, so its failure is not reported as an error",
		"both directions":                                     "the divergence read, which is wrong if only one side is asked",
	}, "the reads section")
}

// TestCompose_ShipToggle_TheRemoteRefIsALocalCache is the gotcha that makes the
// divergence check real.
//
// `origin/main` is a ref on your disk. Without a fetch it answers with whatever
// was true the last time somebody synced, so "the trunk has not moved" is a
// sentence about the local repository and not about the remote — and on a
// machine that has not fetched in a week it reports no divergence forever. It is
// the same silent-failure shape knowledge-base.md documents for an empty date:
// the check returns a clean answer because it never ran.
func TestCompose_ShipToggle_TheRemoteRefIsALocalCache(t *testing.T) {
	body := shipSection(t, "### the reads")
	requireShipPhrases(t, body, map[string]string{
		"is a local cache":                  "what the remote-tracking ref actually is",
		"without a `git fetch` first":       "the missing step that makes the comparison a lie",
		"reports no divergence forever":     "the consequence, which looks exactly like success",
		"say the comparison did not happen": "how a failed fetch must be reported, so silence never becomes an all-clear",
	}, "the reads section")
}

// TestCompose_ShipToggle_DoesNotAssumeItsTooling is the mechanism guard.
//
// EL MECANISMO. Nothing in this block grants a tool. The base block's delegation
// table says `Bash for state (git, openspec status)` runs inline — that is a rule
// about WHERE to run a command, not a promise the command exists, and the table
// lives inside a section the `delegation` toggle removes. Two variants of the
// block (windsurf, antigravity) do not carry that table at all.
//
// And the table does not classify `git push`: its git row is about reading
// state, its execution row names tests, build and install. The one command here
// that changes something outside the machine is in neither.
func TestCompose_ShipToggle_DoesNotAssumeItsTooling(t *testing.T) {
	body := shipSection(t, "### what is actually guaranteed")
	requireShipPhrases(t, body, map[string]string{
		"nothing in this block grants a tool":       "the honest footing the whole section rests on",
		"`bash for state (git, openspec status)`":   "what the block DOES say, quoted rather than paraphrased or denied",
		"a section the `delegation` toggle removes": "why even that sentence cannot be leaned on",
		"`git push` is not a read":                  "the classification gap, named instead of resolved by assumption",
		"neither row names it":                      "and the consequence: it is its own case, not 'git, therefore inline'",
	}, "the guarantees section")
}

// TestCompose_ShipToggle_GhIsAnAssumptionWithItsConsequence asserts the label,
// what it costs if it is false, and that the tool is written about as ABSENT-ABLE
// throughout.
//
// THIS IS THE LOAD-BEARING HALF of the `gh` guard, and the anti-assert is the
// cheap one. Banning phrasings is a race nobody wins — there is always one more
// way to say "gh will be there" — while requiring the declaration and its
// consequence does not depend on anticipating anybody's wording: a fragment that
// keeps these five phrases and then asserts `gh` is guaranteed contradicts itself
// on the same page, and the contradiction is what the banned list has to catch,
// not the universe of ways to write it.
//
// The consequence is what makes the label worth anything. `gh` is named nowhere
// in any asset in this repo, so the assumption is real — and its failure mode is
// the dangerous kind: a `gh pr view` that errors looks exactly like "there is no
// PR", which is precisely the answer that opens the second PR the table above
// forbids. A bare ASSUMPTION label with no consequence attached is decoration.
func TestCompose_ShipToggle_GhIsAnAssumptionWithItsConsequence(t *testing.T) {
	body := shipSection(t, "### what is actually guaranteed")
	requireShipPhrases(t, body, map[string]string{
		"**assumption** (not verified from this repo)": "the label WITH its parenthetical, in the typed vocabulary the " +
			"reporting rules already define. Asserted whole because the bare word survives normalisation and " +
			"prose could then satisfy it by accident",
		"installed and authenticated":           "the two halves of the assumption, which fail separately",
		"if it is unavailable":                  "the absent branch, written into the prose: the tool is discussed as something that may not be there",
		"looks exactly like \"there is no pr\"": "why the failure is dangerous rather than merely inconvenient",
		"on any `gh` failure":                   "the trigger for the honest report, so it covers every way the call can fail and not just one",
		"i could not tell whether a pr exists":  "the report that is honest about it, quoted so it can be copied",
		"never \"there is none\"":               "the reading that is forbidden, which is the one that opens the second PR",
	}, "the guarantees section")
}

// TestCompose_ShipToggle_MainAsksAndYourOwnBranchDoesNot asserts the asymmetry
// AND its reason.
//
// Without the reason the rule reads as timidity and gets generalised in either
// direction: ask for everything, and the toggle is off within a week; ask for
// nothing, and `main` gets a push nobody chose. The reason is one sentence and it
// is the whole rule — whose day breaks if this is wrong.
func TestCompose_ShipToggle_MainAsksAndYourOwnBranchDoesNot(t *testing.T) {
	body := shipSection(t, "### pushing to `main` asks")
	requireShipPhrases(t, body, map[string]string{
		"whose day breaks if it goes wrong": "the reason, in one sentence",
		"affects nobody":                    "why a push to your own branch is not worth an interruption",
		"gets its toggle switched off":      "what asking for everything costs, which is the rule going away entirely",
		"ask every time":                    "the instruction for the other side",
		"including the tenth time":          "the habituation case, which is where this rule is actually lost",
	}, "the main-asks section")
}

// TestCompose_ShipToggle_ProposesPrOrDirectWithGrounds asserts the call is
// proposed with an argument and not taken alone — and asserts WHAT the agent
// cannot see, which is why it is not the agent's call.
//
// Both examples are asserted whole, with their verdicts attached. Asserting
// "propose a recommendation" alone would stay green with one example, and one
// example is not a contrast: the point is that the SAME machinery answers PR on
// one change and direct on another.
func TestCompose_ShipToggle_ProposesPrOrDirectWithGrounds(t *testing.T) {
	body := shipSection(t, "### pr or direct")
	requireShipPhrases(t, body, map[string]string{
		"whether this repo is the person's own or shared":                               "the first fact that is not in the repository",
		"whether `main` is protected":                                                   "the second",
		"working in these files right now":                                              "the third",
		"\"this touches authentication and it is 400 lines across six files — **pr**\"": "the argued recommendation on one side, with its verdict",
		"\"this is one paragraph in the readme — **direct**\"":                          "and on the other, so the contrast is what is on the page",
		"hands them the framing as well as the choice":                                  "what a bare question costs, which is why the grounds are not optional",
		"is **not** \"the branch is unprotected\"":                                      "the failed protection read, which must not arrive as an answer",
	}, "the pr-or-direct section")
}

// TestCompose_ShipToggle_NeverMerges asserts the line where the agent's cycle
// ends, and the two ways it gets crossed.
//
// Publishing something for somebody to look at and deciding on everybody's
// behalf that it is good are different acts. The two crossings are auto-merge
// and green checks: the first is a merge with a delay, the second is a claim
// that nobody found a problem automatically being read as a claim that somebody
// approved it.
func TestCompose_ShipToggle_NeverMerges(t *testing.T) {
	body := shipSection(t, "### it does not merge")
	requireShipPhrases(t, body, map[string]string{
		"deciding on everybody's behalf":                  "what a merge is, as against what a push is",
		"the cycle ends at the open pr and its ci result": "where the agent stops, stated as a place",
		"`gh pr merge`":            "the command named, so the prohibition is checkable",
		"do not enable auto-merge": "the first way the line gets crossed: a merge with a delay is still a merge",
		"is a different claim from somebody having read it": "the second: green checks are not an approval",
	}, "the no-merge section")
}

// TestCompose_ShipToggle_TheAnswerBelongsToTheOrchestrator is the anti-trap.
//
// LA TRAMPA TÉCNICA. A sub-agent that asks a question does not receive an answer:
// it returns the question in its report and terminates, and nothing of it is
// still running when somebody reads it. So the reads can be delegated and the
// three questions cannot. `senior.md` already names this split — "the senior
// labels, you hold" — and citing it is the point: a reader who meets the same
// mechanism twice under two names counts two mechanisms.
//
// And the safe default is the other half. Everywhere else in this block an
// unanswered question costs a delay. Here the wrong default costs a revert
// somebody else performs, so it has to be written down which way silence falls.
func TestCompose_ShipToggle_TheAnswerBelongsToTheOrchestrator(t *testing.T) {
	body := shipSection(t, "### the diagnosis can be delegated")
	requireShipPhrases(t, body, map[string]string{
		"the senior labels, you hold":                                 "the split quoted from the fragment that names it, rather than reinvented here",
		"a sub-agent that asks a question does not receive an answer": "the mechanical fact the whole section rests on",
		"returns the question in its report and terminates":           "why: there is nothing left running to hear the answer",
		"the person is in a conversation with **you**":                "the addressee, named — without it the asking has no venue",
		"nothing gets pushed":                                         "what silence resolves to, which is the only safe direction here",
	}, "the delegation-split section")
}

// TestCompose_ShipToggle_DeclaresItsCost asserts the price is on the page.
//
// A fragment that does not say what it costs gets evaluated against a price of
// zero, and then the bill arrives as interruptions nobody budgeted for. It is a
// diagnosis before every publish and at most two questions per change — and zero
// on the common case, which is worth saying, because a reader who thinks every
// change costs two interruptions turns it off before trying it.
func TestCompose_ShipToggle_DeclaresItsCost(t *testing.T) {
	body := shipSection(t, "### what it costs")
	requireShipPhrases(t, body, map[string]string{
		"one diagnosis before every publish": "the unit of the cost",
		"at most two interruptions":          "the ceiling, so it can be weighed",
		"zero questions":                     "the common case, without which the ceiling reads as the norm",
		"turn it off":                        "that it is a toggle and switching it off is the supported answer",
	}, "the cost section")
}

// TestCompose_ShipToggle_DoesNotClaimAMechanismItLacks is the anti-assert, and it
// runs in BOTH directions.
//
// Fragments in this repo have already shipped guarantees that did not exist, and
// one shipped the mirror image — it denied a mechanism the block documents. Both
// are banned by name, because both send the reader to act on a map of the system
// that is wrong.
//
// IT MATCHES THE CLAIM, NOT THE TYPOGRAPHY. The comparison runs against
// `shipNormalize`, which strips backticks, emphasis and quotes. The first version
// of this test compared raw, and EL CASO CATORCE followed: it banned
// `"gh is guaranteed"` while the fragment itself writes the tool as `gh`, so the
// backticked form — the one anybody editing this file would write — passed clean.
// A guard defeated by the document's own formatting convention is a guard that
// only catches the phrasing nobody uses.
//
// AND THIS LIST IS A LIST, NOT A RULE. That limit is written here on purpose,
// because a guard believed to be complete is worse than one that declares its
// reach: the first makes the reader stop looking. Every entry below is one WAY OF
// SAYING something false, and there is always another wording that is not in the
// list — "you can rely on gh being there" says exactly what the first block
// forbids using none of its phrases. Normalising kills a whole class of evasions
// (the same claim redressed); it does nothing about a claim rephrased. The guard
// that does not depend on anticipating wordings is the POSITIVE one —
// TestCompose_ShipToggle_GhIsAnAssumptionWithItsConsequence — which requires the
// assumption to be declared with its consequence and its failure report. A
// fragment that keeps those and then asserts `gh` is guaranteed contradicts
// itself on the same page, and this list only has to catch the contradiction, not
// enumerate the universe of ways to state it.
func TestCompose_ShipToggle_DoesNotClaimAMechanismItLacks(t *testing.T) {
	body := shipNormalize(shipFragment(t))

	// Direction one: powers nothing here has. Nothing in this block enforces
	// anything, nothing waits, and no tool is guaranteed to exist. Written bare —
	// the normaliser has already removed the backticks the fragment writes them
	// with, so a bare entry catches both dresses and a backticked entry catches
	// neither.
	for _, banned := range []string{
		"ship blocks",
		"ship prevents",
		"blocks the push",
		"prevents the push",
		"makes it impossible",
		"cannot be pushed",
		"automatically pushes",
		"automatically opens a pull request",
		"waits for your answer",
		"waits for confirmation",
		"the sub-agent waits",
		"asks and waits",
		"gh is guaranteed",
		"gh is always available",
		"gh is always there",
		"gh will be there",
		"gh is installed and",
		"you can rely on gh",
		"rely on gh being",
		"assume gh is",
		"always installed",
		"guarantees the push",
		"cannot be overridden",
		"holds a veto",
		"has veto",
	} {
		if strings.Contains(body, banned) {
			t.Errorf("the ship fragment claims %q. Nothing here is a mechanism: no toggle enforces "+
				"anything, a sub-agent returns its report and terminates rather than waiting for an "+
				"answer, and no part of this block installs or guarantees git or `gh`. Stating any of "+
				"it as a capability is the error this repo has already had to correct four times", banned)
		}
	}

	// Direction two: mechanisms the block DOES document. Denying one is the same
	// defect inverted — it sends the reader to build a workaround for a solved
	// problem.
	for _, denied := range []string{
		"nothing in this block mentions git",
		"the block says nothing about git",
		"no command reports whether the change is closed",
		"there is no way to know whether the change closed",
		"there is no verify phase",
		"the cycle has no verify",
	} {
		if strings.Contains(body, denied) {
			t.Errorf("the ship fragment denies a documented mechanism with %q. The base block's "+
				"delegation table names `Bash for state (git, openspec status)`, `openspec list` and "+
				"`openspec status --change` are what it calls its eyes on change state, and `verify` "+
				"is a phase of the cycle the block names in five places. Denying a mechanism that "+
				"exists is the same error as inventing one that does not", denied)
		}
	}
}

// TestCompose_ShipToggle_AntiAssertSeesThroughMarkdown pins the normaliser
// itself, because a guard whose own defect was invisible needs a test that would
// have caught it.
//
// The three strings below are the coordinator's three variants of one false
// claim. Two of them are the same sentence in two dresses, and the raw
// comparison caught only the undressed one. If somebody replaces
// `shipNormalize` with a raw match again, this fails and says why — otherwise
// the regression is silent, exactly as it was the first time.
func TestCompose_ShipToggle_AntiAssertSeesThroughMarkdown(t *testing.T) {
	for _, claim := range []string{
		"gh is guaranteed to be installed.",
		"`gh` is guaranteed to be installed.",
		"**gh** is *guaranteed* to be installed.",
	} {
		if !strings.Contains(shipNormalize(claim), "gh is guaranteed") {
			t.Errorf("the anti-assert does not see through the markdown in %q: normalising gives %q. "+
				"The fragment writes the tool as `gh` because it is a tool, so a guard that only "+
				"matches the bare name is defeated by the document's own formatting convention", claim, shipNormalize(claim))
		}
	}
}

// TestCatalog_OffersTheShipToggle asserts the fragment is reachable.
//
// A fragment that composes but that the catalog never advertises is dead code
// with tests: nothing in the product can turn it on.
func TestCatalog_OffersTheShipToggle(t *testing.T) {
	toggles := orchestratorToggles(t)
	for _, tg := range toggles {
		if tg == "ship" {
			return
		}
	}
	t.Errorf("the catalog does not offer the `ship` toggle (offers: %s): the "+
		"fragment composes and no user can ask for it", strings.Join(toggles, ", "))
}
