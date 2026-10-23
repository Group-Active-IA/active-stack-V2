package config_test

// ──────────────────────────────────────────────────────────────────────────────
// La fase `verify` existe, y el bloque no la nombraba
//
// `openspec init --tools claude` (openspec 1.3.1) instala
// `.claude/commands/opsx/verify.md` — 164 líneas — y la skill
// `openspec-verify-change`. O sea: el comando está disponible en el proyecto y
// el bloque de instrucciones del orquestador no le decía al agente que existe.
//
// El ciclo real es explore → propose → apply → verify → archive, y así está
// escrito en el AGENTS.md de este repo desde antes que en los assets.
//
// Cinco lugares del bloque enumeran el ciclo, y ninguno la incluía:
//   1. el diagrama del workflow
//   2. la lista de slash-commands
//   3. la lista de skills cargadas por contexto
//   4. la tabla de "intención del usuario → skill"
//   5. la tabla de asignación de modelos
//
// CÓMO ESTÁN ESCRITOS ESTOS TESTS, y por qué así:
//
// Buscar la palabra "verify" contra el bloque entero no probaría nada: el base
// ya la escribe en la fila "Read 1-3 files to decide/verify" de la tabla de
// delegación. Cada test acota su lugar (la sección, la lista de bullets, la
// tabla parseada) y además asserta la POSICIÓN — entre apply y archive. Una fase
// listada en el lugar equivocado del ciclo desinforma más que una ausente.
// ──────────────────────────────────────────────────────────────────────────────

import (
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

// verifySkillName is the skill openspec actually installs. Not "openspec-verify":
// the installed directory is `openspec-verify-change`, matching the
// `openspec-apply-change` / `openspec-archive-change` the block already lists.
const verifySkillName = "openspec-verify-change"

// baseBlock composes with delegation and model-routing ON, because two of the
// five places live inside those markers: the antigravity intent table sits
// inside <!-- active-stack:sdd-delegation -->, and every model table sits inside
// <!-- active-stack:sdd-model-assignments -->. Composing without them would
// delete the very sections under test and the assertions would look absent
// rather than wrong.
func baseBlock(t *testing.T, variant string) string {
	t.Helper()
	body, err := config.Compose([]string{"delegation", "model-routing"}, variant)
	if err != nil {
		t.Fatalf("Compose([delegation model-routing], %q): %v", variant, err)
	}
	return strings.ToLower(body)
}

// h2Section returns the body of one "## " section, heading line dropped.
func h2Section(t *testing.T, body, heading string) string {
	t.Helper()
	i := strings.Index(body, heading)
	if i < 0 {
		t.Fatalf("the block has no %q section", heading)
	}
	rest := body[i+len(heading):]
	if j := strings.Index(rest, "\n## "); j >= 0 {
		rest = rest[:j]
	}
	return rest
}

// bulletsAfter returns the run of "- " lines that follows a label line.
//
// Scoping to the run is the point: "Commands Available" holds two lists whose
// entries are near-identical, so an assertion against the whole section passes
// when the entry landed in the wrong list.
func bulletsAfter(t *testing.T, section, label string) string {
	t.Helper()
	i := strings.Index(section, label)
	if i < 0 {
		t.Fatalf("the section has no %q list", label)
	}
	var out []string
	started := false
	for _, line := range strings.Split(section[i:], "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			started = true
			out = append(out, line)
			continue
		}
		if started {
			break
		}
	}
	if len(out) == 0 {
		t.Fatalf("the %q list has no entries", label)
	}
	return strings.Join(out, "\n")
}

// tableRows parses every markdown table row in text into firstCell → the
// REMAINING CELLS, kept apart. Separator and header rows are harmless: nothing
// asserts on them.
//
// Keeping the cells separate is the whole contract, and getting it wrong is how
// this helper shipped broken once. The first version joined the tail with a
// space, which deleted every "|" in the row; a caller that then split the result
// back on "|" got a one-element slice holding the entire row, compared THAT
// against a model name, and could never fail. Anything that needs one specific
// column must be able to ask for that column.
func tableRows(text string) (map[string][]string, []string) {
	rows := map[string][]string{}
	var order []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		if len(cells) < 2 {
			continue
		}
		key := strings.TrimSpace(cells[0])
		tail := make([]string, 0, len(cells)-1)
		for _, c := range cells[1:] {
			tail = append(tail, strings.TrimSpace(c))
		}
		if _, seen := rows[key]; !seen {
			order = append(order, key)
		}
		rows[key] = tail
	}
	return rows, order
}

// rowText flattens a parsed row for the callers that only need "does this row
// mention X anywhere". Callers that need a specific column index the slice.
func rowText(cells []string) string { return strings.Join(cells, " ") }

// orderOf reports the index of want in order, or -1.
func orderOf(order []string, want string) int {
	for i, k := range order {
		if k == want {
			return i
		}
	}
	return -1
}

// betweenApplyAndArchive asserts verify sits in its place in the cycle.
//
// Presence alone is a weak assertion: a `verify` bullet appended after `archive`
// documents a cycle that does not exist, and reads as authoritative while doing
// it.
func betweenApplyAndArchive(t *testing.T, variant, where, text string) {
	t.Helper()
	apply := strings.Index(text, "/opsx:apply")
	verify := strings.Index(text, "/opsx:verify")
	archive := strings.Index(text, "/opsx:archive")
	if verify < 0 {
		t.Errorf("[%s] %s never names /opsx:verify: the command is installed by "+
			"`openspec init` and the block does not tell the agent it exists", variant, where)
		return
	}
	if apply < 0 || archive < 0 {
		t.Fatalf("[%s] %s is missing apply or archive", variant, where)
	}
	if !(apply < verify && verify < archive) {
		t.Errorf("[%s] %s places verify outside apply → verify → archive "+
			"(apply=%d verify=%d archive=%d)", variant, where, apply, verify, archive)
	}
}

// ── 1. the workflow diagram ──────────────────────────────────────────────────

func TestOrchestrator_WorkflowDiagram_HasVerifyBetweenApplyAndArchive(t *testing.T) {
	for _, v := range evidenceVariants {
		section := h2Section(t, baseBlock(t, v), "## opsx workflow")
		betweenApplyAndArchive(t, v, "the workflow diagram", section)
	}
}

// ── 2. the slash-command list ────────────────────────────────────────────────

func TestOrchestrator_SlashCommands_ListVerify(t *testing.T) {
	for _, v := range evidenceVariants {
		section := h2Section(t, baseBlock(t, v), "## commands available")
		list := bulletsAfter(t, section, "slash commands (type directly):")
		betweenApplyAndArchive(t, v, "the slash-command list", list)
	}
}

// ── 3. the skills list ───────────────────────────────────────────────────────

func TestOrchestrator_Skills_ListVerifyChange(t *testing.T) {
	for _, v := range evidenceVariants {
		section := h2Section(t, baseBlock(t, v), "## commands available")
		// Prefix, not the full label: `generic` writes
		// "Skills (loaded by context, appear in autocomplete):".
		list := bulletsAfter(t, section, "skills (loaded by context")
		if !strings.Contains(list, verifySkillName) {
			t.Errorf("[%s] the skills list never names %q: the agent is told about the "+
				"slash command but not about the skill behind it", v, verifySkillName)
			continue
		}
		apply := strings.Index(list, "openspec-apply-change")
		verify := strings.Index(list, verifySkillName)
		archive := strings.Index(list, "openspec-archive-change")
		if !(apply < verify && verify < archive) {
			t.Errorf("[%s] the skills list places %q outside apply → verify → archive",
				v, verifySkillName)
		}
	}
}

// ── 4. the "user intent → skill" table ───────────────────────────────────────

// TestOrchestrator_IntentTable_RoutesVerify asserts a TABLE ROW maps a verify
// intent to the verify skill.
//
// Asserted against parsed rows, not against the block: the skills list one
// section up writes the same skill name as a bullet, so a substring search would
// go green with the routing table untouched.
func TestOrchestrator_IntentTable_RoutesVerify(t *testing.T) {
	for _, v := range evidenceVariants {
		rows, _ := tableRows(baseBlock(t, v))
		found := ""
		for intent, cells := range rows {
			if strings.Contains(rowText(cells), verifySkillName) {
				found = intent
				break
			}
		}
		if found == "" {
			t.Errorf("[%s] no routing row sends any user intent to %q: the agent has the "+
				"skill listed and no rule for when to reach for it", v, verifySkillName)
			continue
		}
		if !strings.Contains(found, "verify") {
			t.Errorf("[%s] the row routing to %q is keyed on %q, which never says verify",
				v, verifySkillName, found)
		}
	}
}

// ── 5. the model assignment table ────────────────────────────────────────────

// TestOrchestrator_ModelTable_AssignsVerify asserts verify has its own row.
//
// Without one it falls to `default` — sonnet — by omission rather than by
// choice, and the table stops recording a decision.
func TestOrchestrator_ModelTable_AssignsVerify(t *testing.T) {
	for _, v := range evidenceVariants {
		section := h2Section(t, baseBlock(t, v), "## model assignments")
		rows, order := tableRows(section)
		cells, ok := rows["verify"]
		if !ok {
			t.Errorf("[%s] the model table has no `verify` row: the phase falls to the "+
				"`default` row by omission, which is not an assignment", v)
			continue
		}
		// Column 0 of the tail is the "Default Model" column. Read it by index:
		// this is the assertion the whole model choice rests on, so it has to
		// look at the model cell and nothing else.
		model := cells[0]
		if model == "" || model == "—" || model == "default" || model == "sonnet" {
			t.Errorf("[%s] the `verify` row assigns %q. Falling back to the default "+
				"(sonnet) is the outcome this row exists to override: verify judges "+
				"conformity AND hypothesises failures nobody wrote down, and a cheaper "+
				"reviewer reads faster and agrees", v, model)
		}
		// And the reason column has to say something: a model with no reason is a
		// preference, and the next person to read the table cannot argue with it.
		if len(cells) < 2 || cells[1] == "" {
			t.Errorf("[%s] the `verify` row assigns %q with no reason next to it", v, model)
		}
		ia, iv, iar := orderOf(order, "apply"), orderOf(order, "verify"), orderOf(order, "archive")
		if !(ia < iv && iv < iar) {
			t.Errorf("[%s] the model table lists verify outside apply → verify → archive "+
				"(apply=%d verify=%d archive=%d)", v, ia, iv, iar)
		}
	}
}

// TestOrchestrator_EngramWiring_ListsVerify asserts the SEVENTH place that
// enumerates the cycle also names the phase.
//
// engram-protocol.md's OPSX wiring says to save "after every completed OPSX
// action (explore/propose/apply/archive)" under
// `topic_key: "opsx/{change-name}/{phase}"`. That list governs the very channel
// review.md has to reason about, so a verify pass falling outside it is a phase
// whose progress is never persisted — and the one enumeration whose staleness
// changes behaviour rather than only documentation.
func TestOrchestrator_EngramWiring_ListsVerify(t *testing.T) {
	body, err := config.Compose([]string{"engram"}, "claude")
	if err != nil {
		t.Fatalf("Compose([engram]): %v", err)
	}
	low := strings.ToLower(body)
	i := strings.Index(low, "after every completed opsx action")
	if i < 0 {
		t.Fatal("the engram protocol no longer has the OPSX save trigger")
	}
	// Just the parenthesised phase list on that line, so the assertion cannot be
	// satisfied by the word "verify" appearing anywhere else in the protocol.
	line := low[i:]
	if nl := strings.Index(line, "\n"); nl >= 0 {
		line = line[:nl]
	}
	open, close := strings.Index(line, "("), strings.Index(line, ")")
	if open < 0 || close < open {
		t.Fatalf("the OPSX save trigger no longer lists its phases: %q", line)
	}
	phases := line[open+1 : close]
	if !strings.Contains(phases, "verify") {
		t.Errorf("the engram OPSX wiring saves after (%s): a verify pass is not a "+
			"completed action, so its progress is never persisted and the next "+
			"session cannot tell the review ran", phases)
	}
}

// TestOrchestrator_ModelTable_VerifyRowNamesItsDegradation asserts the verify
// row says what happens when opus is unavailable.
//
// The model-assignments preamble says that without access to the assigned model
// you substitute sonnet and continue — silently. For every other phase that is a
// slower answer; for verify it is a different answer, because the reason this
// row is opus at all is that a cheaper reviewer reads faster and agrees. The
// substitution has to be visible where it happens.
func TestOrchestrator_ModelTable_VerifyRowNamesItsDegradation(t *testing.T) {
	for _, v := range evidenceVariants {
		section := h2Section(t, baseBlock(t, v), "## model assignments")
		rows, _ := tableRows(section)
		cells, ok := rows["verify"]
		if !ok || len(cells) < 2 {
			t.Errorf("[%s] the model table has no `verify` row with a reason", v)
			continue
		}
		if !strings.Contains(cells[1], "sonnet") {
			t.Errorf("[%s] the `verify` reason is %q and never mentions the sonnet "+
				"substitution: the block performs it silently, and for this phase a "+
				"cheaper model does not answer slower, it answers differently", v, cells[1])
		}
	}
}

func TestArchiveSonnetEnOchoVariantes(t *testing.T) {
	for _, v := range evidenceVariants {
		t.Run(v, func(t *testing.T) {
			body := h2Section(t, baseBlock(t, v), "## model assignments")
			rows, _ := tableRows(body)
			row, ok := rows["archive"]
			if !ok || len(row) < 2 {
				t.Fatal("falta fila archive")
			}
			if row[0] != "sonnet" {
				t.Errorf("archive usa %s, se esperaba sonnet", row[0])
			}
			for _, frase := range []string{"evidence", "substitution", "unavailable"} {
				if !strings.Contains(row[1], frase) {
					t.Errorf("motivo archive no declara %s", frase)
				}
			}
			for _, frase := range []string{"less on archive", "file operations"} {
				if strings.Contains(body, frase) {
					t.Errorf("instrucción contradictoria: %s", frase)
				}
			}
		})
	}
}
