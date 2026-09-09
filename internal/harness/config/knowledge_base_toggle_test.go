package config_test

// ──────────────────────────────────────────────────────────────────────────────
// La knowledge-base: el mapa al que nada le vuelve a escribir
//
// `kb-creator` corre UNA vez, en la fundación, y deja diez documentos. Después
// se marca terminado y nada del ciclo vuelve a abrirlos. El arranque del bloque
// base lo dice explícito: si ya hay fundación, andá directo a `openspec list`.
//
// A los tres meses el documento del modelo de datos describe un esquema que
// cambió hace dos changes, y el agente le cree porque es el documento canónico.
// No alucinó: le creyó al mapa, y el mapa estaba viejo.
//
// LO QUE ESTE ARCHIVO CUIDA, además de que las reglas estén escritas:
//
//   - Cada test acota su sección y le saca el encabezado, por el bug que este
//     repo ya cazó diez veces: un test que busca una frase que el título de su
//     propia sección también escribe pasa con la sección vaciada.
//   - Los DOS MECANISMOS que no existen y el fragmento tiene que declarar:
//       · la fecha del documento. mtime NO es fecha de edición: `git clone`
//         estampa todos los archivos con la hora del checkout, así que en un
//         clon fresco la KB entera parece escrita hoy y el chequeo no reporta
//         nada — justo en la máquina que menos contexto tiene. `git log` es
//         mejor y sigue siendo "cuándo cambiaron los bytes", no "sigue siendo
//         cierto".
//       · el update en archive. `openspec archive` sincroniza specs y no sabe
//         que la KB existe: ningún comando lo hace. Viaja en el prompt del
//         sub-agente, que es un pedido, igual que en roles.md y review.md.
//     Un anti-assert cuida que no vuelvan a aparecer garantías falsas.
//   - El detector: si sólo apilaste, no actualizaste. Apilar sale gratis y se
//     siente trabajo; reescribir obliga a releer y BORRAR. "Actualicé la KB" es
//     la misma frase hayas hecho la mitad difícil o no, así que las dos listas
//     van separadas y una lista REWRITTEN vacía se afirma, no se deja en blanco.
// ──────────────────────────────────────────────────────────────────────────────

import (
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/catalog"
	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

// kbHeading is the fragment's own heading, and it is the marker on purpose: the
// base block writes "kb-creator" and "foundation" on its own, so a bare word
// would match a block composed without this fragment at all.
const kbHeading = "## the knowledge base"

// kbFragment returns ONLY the knowledge-base fragment, never the whole block.
func kbFragment(t *testing.T) string {
	t.Helper()
	body, err := config.Compose([]string{"knowledge-base"}, "claude")
	if err != nil {
		t.Fatalf("Compose([knowledge-base]): %v", err)
	}
	lower := strings.ToLower(body)
	i := strings.Index(lower, kbHeading)
	if i < 0 {
		t.Fatalf("the composed block does not carry the knowledge-base fragment (heading %q missing)", kbHeading)
	}
	return lower[i:]
}

// kbFlat collapses every run of whitespace to a single space.
//
// The fragment is hard-wrapped at 80 columns, so a phrase asserted below can sit
// astride a newline and be absent as a literal substring while being present as
// prose. Without this, re-wrapping a paragraph — which changes nothing anybody
// reads — turns these tests red, and the fix somebody reaches for is deleting
// the assert.
func kbFlat(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// kbSection returns the BODY of one "### " section with its own heading line
// dropped. Dropping the heading is the whole point: a test that matches a phrase
// its section title writes stays green with the section emptied.
func kbSection(t *testing.T, headingPrefix string) string {
	t.Helper()
	body := kbFragment(t)
	i := strings.Index(body, headingPrefix)
	if i < 0 {
		t.Fatalf("the knowledge-base fragment has no section starting with %q", headingPrefix)
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
	return kbFlat(rest)
}

// kbPreamble returns the body of the "## " heading: everything from the
// fragment's own heading line to the first "### ".
//
// It exists because the section-by-section degradation ran over the six "###"
// sections, and the fragment has six sections AND a preamble. Deleting all four
// opening paragraphs — the ten documents, "nothing in the cycle ever writes back
// to it", "it believed the map, and the map was old" — turned nothing red. The
// fragment would have kept every rule and lost the whole argument for why the
// rules exist, which is the half a reader needs to apply them to a case nobody
// wrote down.
func kbPreamble(t *testing.T) string {
	t.Helper()
	body := kbFragment(t)
	nl := strings.Index(body, "\n")
	if nl < 0 {
		t.Fatal("the knowledge-base fragment is a single line")
	}
	rest := body[nl+1:]
	if j := strings.Index(rest, "\n### "); j >= 0 {
		rest = rest[:j]
	}
	return kbFlat(rest)
}

// TestCompose_KnowledgeBaseToggle_PreambleCarriesTheArgument asserts the opening
// keeps the reason the fragment exists.
//
// Rules without their argument get applied literally and stop at the first case
// nobody anticipated. This preamble carries the diagnosis: the map is written
// once and never revisited, the agent that trusts it is not hallucinating, and a
// wrong map with authority is worse than no map — the sentence that tells a
// reader why "report, do not repair" is not excessive caution.
func TestCompose_KnowledgeBaseToggle_PreambleCarriesTheArgument(t *testing.T) {
	body := kbPreamble(t)
	requirePhrases(t, body, map[string]string{
		"runs once":                                "kb-creator writes the map one time, which is the root of the whole problem",
		"nothing in the cycle ever writes back":    "the gap, stated as the defect it is",
		"it believed the map, and the map was old": "the failure named, so it is not filed as hallucination",
		"costs more than no map at all":            "why a stale map is worse than an absent one",
		"builds on top of it":                      "the asymmetry: no map makes you ask, a wrong map makes you build",
	}, "the fragment preamble")
}

// TestCompose_KnowledgeBaseToggle_IsKnown asserts every variant accepts it: the
// fragment is shared, so no variant may reject it.
func TestCompose_KnowledgeBaseToggle_IsKnown(t *testing.T) {
	for _, v := range evidenceVariants {
		if _, err := config.Compose([]string{"knowledge-base"}, v); err != nil {
			t.Errorf("Compose([knowledge-base], %q): %v", v, err)
		}
	}
}

// TestCompose_KnowledgeBaseToggle_Subtractive asserts a block that did not ask
// for it does not carry it.
func TestCompose_KnowledgeBaseToggle_Subtractive(t *testing.T) {
	off, err := config.Compose([]string{"delegation", "model-routing", "engram", "tdd", "governance", "evidence", "auditor", "roles", "review"}, "claude")
	if err != nil {
		t.Fatalf("Compose without knowledge-base: %v", err)
	}
	if strings.Contains(strings.ToLower(off), kbHeading) {
		t.Error("the knowledge-base fragment leaks into a block that did not ask for it")
	}
}

// TestCompose_KnowledgeBaseToggle_ComposesInsideTheReviewFamily pins the
// position argued for in compose.go.
//
// Order is argued in a comment and fixed by nothing else. The golden covers it
// too, but a golden fails on any asset edit and says only "the bytes differ":
// this one says which neighbour moved.
func TestCompose_KnowledgeBaseToggle_ComposesInsideTheReviewFamily(t *testing.T) {
	all, err := config.Compose(config.KnownToggles(), "claude")
	if err != nil {
		t.Fatalf("Compose(all): %v", err)
	}
	low := strings.ToLower(all)
	kb := strings.Index(low, kbHeading)
	rv := strings.Index(low, reviewHeading)
	// engram-protocol.md's own heading, in full. NOT "## engram": the base block
	// writes that exact markdown heading inside the sub-agent prompt template, so
	// the short marker matched the BASE at offset 6801 and the order assert
	// compared the fragment against a heading the fragment does not own. Caught by
	// running it, which is the only way this class of bug ever surfaces.
	eng := strings.Index(low, "## engram persistent memory")
	if kb < 0 || rv < 0 || eng < 0 {
		t.Fatalf("missing heading in the all-toggles block (kb=%d review=%d engram=%d)", kb, rv, eng)
	}
	if !(rv < kb && kb < eng) {
		t.Errorf("the knowledge-base fragment composes out of place (review=%d kb=%d engram=%d). "+
			"It belongs after review — everything above it is scoped to the change in flight and "+
			"this is the only section whose subject outlives it — and before engram, so engram's "+
			"protocol does not read as if it already covered a map on disk", rv, kb, eng)
	}
}

// TestCompose_KnowledgeBaseToggle_ReportsDriftDoesNotRepairIt asserts the
// session-start rule keeps the two halves apart.
//
// A fragment that says "check the knowledge base at startup" and stops there
// gets read as permission to fix what it finds. The most authoritative document
// in the project is the last one that should be edited by a process nobody is
// watching.
func TestCompose_KnowledgeBaseToggle_ReportsDriftDoesNotRepairIt(t *testing.T) {
	body := kbSection(t, "### at session start")
	requirePhrases(t, body, map[string]string{
		// NOT the bare words "report" or "drift": the section title writes both,
		// and the fragment's opening paragraphs write "report" again.
		"do not rewrite it on your own":                 "the prohibition, which is the half that gets dropped",
		"is not edited by a process nobody is watching": "why: authority is exactly what makes an unwatched edit expensive",
		"the person decides":                            "who owns the call once the drift is reported",
		"before the work, not after":                    "when: afterwards the map has already been believed",
	}, "the session-start section")
}

// TestCompose_KnowledgeBaseToggle_NamesTheChangesInTheDriftReport asserts the
// report carries the change names.
//
// "The knowledge base may be out of date" is true of every knowledge base on
// every day. A warning that is always true trains the reader to skip it, and
// then the drift report costs attention and buys nothing.
func TestCompose_KnowledgeBaseToggle_NamesTheChangesInTheDriftReport(t *testing.T) {
	body := kbSection(t, "### at session start")
	requirePhrases(t, body, map[string]string{
		"name them":              "the instruction: the names are the payload",
		"with no names is noise": "what an unnamed report degrades into",
		"checkable in a minute":  "what the names buy the reader",
		"archived since":         "the worked shape, so the report is executable and not aspirational",
		// The names of the CHANGES were there from the start; the names of the
		// DOCUMENTS were not. One date for the whole knowledge base tells the reader
		// something is stale without saying which of the ten to open — the
		// fragment's own argument about names, turned against its own example.
		"per document, last written":         "the report is dated per document, not once for the whole map",
		"which of the ten documents to open": "what a single collapsed date fails to answer",
	}, "the session-start section")
}

// TestCompose_KnowledgeBaseToggle_SaysWhatTheTwoDatesActuallyMean is the guard
// on the mechanism this check rests on.
//
// EL MECANISMO. The comparison is between two dates and neither means "is this
// document still true":
//
//   - the archive date is read off `openspec/changes/archive/YYYY-MM-DD-<name>/`,
//     which is a path convention documented in the base block and never
//     exercised here. If archive names the directory differently, one side of
//     the comparison is empty and the check reports no drift forever;
//   - the document date is worse. mtime is the checkout time after a `git
//     clone`, so on a fresh clone the whole knowledge base looks written today.
//     `git log` is the better date and is still only when the bytes changed.
//
// A fragment that skips this ships a check whose failure mode is silence.
func TestCompose_KnowledgeBaseToggle_SaysWhatTheTwoDatesActuallyMean(t *testing.T) {
	body := kbSection(t, "### what the two dates")
	requirePhrases(t, body, map[string]string{
		"neither one means what you want it to": "the headline: both dates are proxies",
		"every file with the checkout time":     "the concrete gotcha — a fresh clone makes the whole map look current",
		"git log -1 --format=%ai --":            "the better date, in a form somebody can run — the AUTHOR date",
		"when the bytes last changed":           "the limit of even the better date",
		"refreshes documents nobody reviewed":   "the false negative a bulk reformat creates",
		"one side of the comparison is empty":   "what it costs if the archive path convention is wrong",
		// The committer date is the trap: it is the obvious remedy for the mtime
		// problem and fails the same way.
		"rewritten by rebase, by amend, and above all by a squash-merge": "why %cI is not the fix: a squash restamps every document in the branch",
		"without a byte changing in any of them":                         "the consequence, in the same terms as the fresh-clone failure it mirrors",
		// The document side goes empty too. This was declared for the archive side
		// only: same mechanism, same consequence, so it is declared on both.
		"nothing and exits 0": "an untracked path yields no date and no error, so a caller reads the empty string as one",
		"not committed yet":   "the ordinary state right after kb-creator runs — the exact moment this fragment is about",
		"i could not check":   "how an empty date must be read, so it never becomes a silent 'no drift'",
	}, "the two-dates section")
	// The failure mode of the whole check, stated so it is not sold as a
	// freshness guarantee. NOT the bare word "blind": nothing else writes it,
	// but the half that matters is WHICH case it misses.
	if !strings.Contains(body, "blind to the document edited last week and wrong") {
		t.Error("the two-dates section never names the case this check cannot see: " +
			"a document rewritten recently and wrong anyway passes every date comparison, " +
			"and a check sold without its blind spot gets trusted past it")
	}
}

// TestCompose_KnowledgeBaseToggle_ProposeNamesTheDocumentsUpFront asserts the
// declaration happens in propose and not later.
//
// Asked at archive time the question has no good answer: by then the change
// touched a bit of everything. Naming the documents while the plan is still a
// plan is what gives archive something concrete to act on.
func TestCompose_KnowledgeBaseToggle_ProposeNamesTheDocumentsUpFront(t *testing.T) {
	body := kbSection(t, "### in propose")
	requirePhrases(t, body, map[string]string{
		"which of the two operations":                "not just which documents: which kind of edit each one needs",
		"asked afterwards the question gets a shrug": "why up front, stated as the failure it prevents",
		// NOT the bare "knowledge base impact" and NOT the bare "proposal.md": the
		// prompt block added below writes both again, so each had two occurrences in
		// this one section and deleting either left the assert green. Caught by
		// degrading phrase by phrase, the pass that sees what emptying a whole
		// section cannot. Pin the two occurrences that carry different weight: the
		// prose naming where the section lives, and the prompt line that is what
		// actually reaches the sub-agent.
		"goes in `proposal.md` under a heading of its own":        "where the section lives, named in the prose the orchestrator reads",
		"add a `## knowledge base impact` section to proposal.md": "the same instruction inside the prompt, which is the copy the sub-agent gets",
		"nothing to update":         "the considered 'none', which an absent section cannot express",
		"look identical afterwards": "why the explicit none is required at all",
	}, "the propose section")
}

// TestCompose_KnowledgeBaseToggle_ProposeSectionIsARequestNotAContract asserts
// the fragment admits what carries the new heading.
//
// `openspec instructions` generates the artifact templates and this heading is
// not one of them. THE BLOCKING DEFECT: the fragment used to say, in the
// indicative, that the orchestrator's prompt "asks for it". It does not. The
// sub-agent prompt template in the base block asks for five things and the
// knowledge base is none of them; `rg -ci 'knowledge.base'` over all eight base
// variants returns zero; and the base opens by stating that each sub-agent
// starts with no context. So this section was addressed to a reader that never
// received it. It now ORDERS the orchestrator to paste the lines, which is the
// lever it actually holds — the same one evidence-report.md uses for the report
// format. The honesty did not disappear, it moved to the right level: two
// links, and neither is enforced.
func TestCompose_KnowledgeBaseToggle_ProposeSectionIsARequestNotAContract(t *testing.T) {
	body := kbSection(t, "### in propose")
	requirePhrases(t, body, map[string]string{
		"this heading is not in them":                                   "the templates do not carry it",
		"when you launch the propose sub-agent, add this to its prompt": "the imperative, without which this section has no addressee",
		"the knowledge base is none of them":                            "why: the base template's five bullets do not include it",
		"each sub-agent starts with no context":                         "the base block's own sentence, which is what makes an unpasted rule unread",
		"both are requests":                                             "the honesty, now at the level where the asking happens",
		"fall back to reading the diff":                                 "the fallback, so a missing section is not a dead end",
	}, "the propose section")
}

// TestCompose_KnowledgeBaseToggle_ArchiveSplitsRewriteFromAppend asserts the two
// operations stay two.
//
// Collapsed into one they both become "append", because appending is what comes
// out when nobody drew the line. The result is a document where the old schema
// and the new one sit next to each other looking equally current.
func TestCompose_KnowledgeBaseToggle_ArchiveSplitsRewriteFromAppend(t *testing.T) {
	body := kbSection(t, "### in archive")
	requirePhrases(t, body, map[string]string{
		// NOT "the new thing": the rule and its counter-example used to write that
		// phrase in the SAME LINE, so deleting the rule left the counter-example
		// holding the assert up. Verified by degrading exactly that and watching it
		// stay green — the twelfth time this repo has hit this shape, and the first
		// the section-level degradation could not see, because emptying the body
		// kills both occurrences at once. The prose was reworded so the rule owns
		// its phrase alone.
		"says the new schema and the old sentences are gone": "what a rewritten section leaves behind: the current state, with the superseded text deleted",
		"believes whichever comes first":                     "the cost of two undated contradicting descriptions",
		"its date and what it cost":                          "what an appended entry carries",
		"history accumulates":                                "why appending is right for the other half",
		"how is it now":                                      "the test that sorts a section into rewrite",
		"how did we get here":                                "the test that sorts a section into append",
	}, "the archive section")
	// The two operation names as emphasised labels, not as bare words: "rewrite"
	// and "append" are written by the propose section, by the detector section and
	// by this section's own title. Bare words left this assert green with the
	// whole split deleted — verified by degrading it.
	for _, op := range []string{"**rewrite**", "**append**"} {
		if !strings.Contains(body, op) {
			t.Errorf("the archive section has no %q operation labelled: a split nobody "+
				"labelled collapses into appending, which is the one that comes free", op)
		}
	}
}

// TestCompose_KnowledgeBaseToggle_ArchiveDoesNotClaimTheCLIDoesIt asserts the
// fragment names who actually performs the update.
//
// `openspec archive` syncs delta specs into the main specs. It knows nothing
// about the knowledge base, so no CLI command performs this step: it travels in
// the archive sub-agent's prompt. Saying otherwise would make an unperformed
// step look like an automated one, which is the failure mode this whole fragment
// was written about.
func TestCompose_KnowledgeBaseToggle_ArchiveDoesNotClaimTheCLIDoesIt(t *testing.T) {
	body := kbSection(t, "### in archive")
	requirePhrases(t, body, map[string]string{
		"knows nothing about the knowledge base": "what the CLI actually does not do",
		"no cli command performs this step":      "stated plainly, so nobody waits for one",
		"asked for, not enforced":                "the footing, in roles.md's and review.md's own terms",
		// Same blocking defect as propose: these rules were written for the archive
		// sub-agent and nothing put them in front of it.
		"when you launch the archive sub-agent, add this to its prompt": "the imperative that gives these rules a reader",
		"headed rewritten and appended":                                 "the two lists demanded in the prompt itself, not only in this section",
		"if rewritten is empty, say why":                                "the detector, carried into the prompt where the sub-agent can act on it",
	}, "the archive section")
}

// TestCompose_KnowledgeBaseToggle_TheDetectorIsTheRewriteList asserts the check
// that tells a real update from half of one.
//
// EL DETECTOR. Appending comes out easy: one dated entry at the bottom, nothing
// re-read, nothing deleted. Rewriting is the half that forces a reader back
// through what was written three months ago to delete what stopped being true —
// and that deletion is the point, because the stale line is the one the next
// agent acts on.
func TestCompose_KnowledgeBaseToggle_TheDetectorIsTheRewriteList(t *testing.T) {
	body := kbSection(t, "### the detector")
	requirePhrases(t, body, map[string]string{
		"if you only appended, you did not update": "the detector itself",
		"re-read what was written three months":    "what rewriting costs, which is why it gets skipped",
		// NOT the bare "separately and by document": the archive prompt block writes
		// "return the two lists separately and by document" too, so deleting the
		// detector's own rule left the archive copy holding this up. Pin the verb.
		"so report the two lists **separately and by document**": "how the two lists are reported: one merged list hides the split again",
		// NOT the bare "rewritten": the section writes it again in "an empty
		// REWRITTEN list is a claim", so deleting the whole worked report left the
		// assert green. Pin both labels as they appear in the report itself.
		"rewritten data model": "the hard half, labelled in the reported shape",
		"appended decisions":   "the easy half, labelled beside it — one list without the other hides the split",
		"is the same sentence": "why reporting the split is required: the claim reads identically either way",
	}, "the detector section")
}

// TestCompose_KnowledgeBaseToggle_EmptyRewriteListIsAClaim asserts silence is
// not an allowed output for the hard half.
//
// Mirrors review.md's empty verdict. A change can legitimately add only history,
// but "nothing to rewrite" left blank and "nothing to rewrite" never considered
// leave identical evidence behind, and the second is far more common.
func TestCompose_KnowledgeBaseToggle_EmptyRewriteListIsAClaim(t *testing.T) {
	body := kbSection(t, "### the detector")
	requirePhrases(t, body, map[string]string{
		"a claim, not a blank":            "the rule",
		"nobody can tell those two apart": "the reason: a blank and a considered none are the same evidence",
	}, "the detector section")
}

// TestCompose_KnowledgeBaseToggle_LabelsTheDocumentNamesAsAssumption asserts the
// fragment does not hardcode filenames it cannot verify.
//
// `kb-creator` is cloned from a separate repository and nothing in this repo
// pins what it writes. A hardcoded list that stops matching reports "this
// project has no knowledge base" about a project that has one — the same wrong
// answer the fragment exists to prevent, reached from the other side.
func TestCompose_KnowledgeBaseToggle_LabelsTheDocumentNamesAsAssumption(t *testing.T) {
	body := kbSection(t, "### read the directory")
	// NOT the bare word "assumption": the drift and archive sections both use the
	// typed categories, and the word travels with them. What is asserted is the
	// label ON THIS CLAIM plus the consequence of it being false.
	requirePhrases(t, body, map[string]string{
		"not verified from this repo":        "the label attached to this specific claim",
		"nothing here pins what it writes":   "why it cannot be verified from here",
		"this project has no knowledge base": "what it costs if the names drift: a false negative that reads like a fact",
		"list the knowledge base directory":  "the instruction that survives a rename",
	}, "the read-the-directory section")
}

// TestCompose_KnowledgeBaseToggle_DoesNotClaimAMechanismItLacks is the anti-
// assert.
//
// Two fragments in this repo have already shipped guarantees that did not exist
// and had to be corrected. Everything here is a prompt: an agent that ignores it
// ignores it, and there is no detector, no sync and no automatic anything behind
// any of these words.
func TestCompose_KnowledgeBaseToggle_DoesNotClaimAMechanismItLacks(t *testing.T) {
	body := kbFlat(kbFragment(t))
	for _, banned := range []string{
		"automatically updated",
		"automatically detects",
		"kept in sync",
		"keeps the knowledge base in sync",
		"guaranteed current",
		"always up to date",
		"cannot go stale",
		"openspec archive updates",
		"proves the map is current",
		"the modification time is the edit",
	} {
		if strings.Contains(body, banned) {
			t.Errorf("the knowledge-base fragment claims %q. Nothing here is a mechanism: "+
				"the drift check is two proxy dates the agent reads by hand, and the archive "+
				"update travels in a sub-agent prompt that nothing enforces. Stating it as a "+
				"guarantee is how a stale map gets trusted, which is the failure this fragment "+
				"was written about", banned)
		}
	}
}

// TestCatalog_OffersTheKnowledgeBaseToggle asserts the fragment is reachable.
//
// A fragment that composes but that the catalog never advertises is dead code
// with tests: nothing in the product can turn it on.
func TestCatalog_OffersTheKnowledgeBaseToggle(t *testing.T) {
	toggles := orchestratorToggles(t)
	for _, tg := range toggles {
		if tg == "knowledge-base" {
			return
		}
	}
	t.Errorf("the catalog does not offer the `knowledge-base` toggle (offers: %s): the "+
		"fragment composes and no user can ask for it", strings.Join(toggles, ", "))
}

// TestCatalog_OrchestratorDescriptionIsNotStaleAboutToggles is the fragment's
// own lesson aimed at this repo.
//
// The sdd-orchestrator long_description enumerates the toggles for the person
// choosing in the TUI. It listed five while the catalog offered nine, so the
// document the user reads described a product two toggles out of date. That is
// the same defect the fragment describes, in the same repo that ships it.
func TestCatalog_OrchestratorDescriptionIsNotStaleAboutToggles(t *testing.T) {
	c, err := catalog.Load()
	if err != nil {
		t.Fatalf("catalog.Load: %v", err)
	}
	for _, h := range c.Harnesses {
		if h.ID != "sdd-orchestrator" {
			continue
		}
		for _, lang := range []string{"es", "en"} {
			low := strings.ToLower(h.LongDescription[lang])
			if low == "" {
				t.Errorf("the sdd-orchestrator harness has no %s long_description", lang)
				continue
			}
			// Iterate the catalog's own toggle list instead of naming one. Hardcoding
			// "knowledge-base" goes stale the moment toggle eleven lands: the
			// description falls behind again and this test stays green, which is
			// exactly the state it was written to end.
			for _, tg := range orchestratorToggles(t) {
				if !strings.Contains(low, tg) {
					t.Errorf("the sdd-orchestrator %s long_description enumerates the toggles "+
						"and omits %q: the user picking in the TUI reads a list that is already "+
						"behind the catalog", lang, tg)
				}
			}
		}
		return
	}
	t.Fatal("the catalog has no sdd-orchestrator harness")
}
