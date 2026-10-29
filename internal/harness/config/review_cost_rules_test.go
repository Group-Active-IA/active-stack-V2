package config

import (
	"strings"
	"testing"
)

// ──────────────────────────────────────────────────────────────────────────────
// Las dos reglas de costo que viven en el base
//
// De dónde salen: el 2026-09-09 se corrió el ciclo OPSX completo sobre el primer
// cambio de un proyecto real (`foundation-setup`: scaffolding, gobernanza BAJA,
// sin una línea de dominio). Lo medido, no lo opinado:
//
//   - 22 sub-agentes, 4.900.000 tokens.
//   - La revisión sola: 1.939.129, el 39% del total.
//   - Ocho revisores en verify: 1.351.186 — casi exactamente lo que costó
//     escribir todo el código (1.325.302).
//   - 27 bloqueantes reportados, 12 problemas reales.
//   - El agente que corrigió los doce: 816.762 él solo, el 17%.
//
// Y el auditor dejó escrito el diagnóstico de una de las fugas: de cuatro
// convergencias, una era barata — `workbox-routing` lo reportaron cuatro
// revisores, pero es un chequeo de dos comandos hecho cuatro veces.
//
// POR QUÉ VAN EN `## Rules` Y NO EN UN FRAGMENTO. Las dos hablan de qué hace el
// orquestador ANTES y DESPUÉS de convocar. Estacionadas en `roles.md` o en
// `review.md` desaparecen con su toggle, y desaparecen justo en la instalación
// donde nadie eligió esos fragmentos: la que más barato paga convocar de más.
// `## Rules` es la única sección del base que ningún toggle borra — el toggle
// `delegation` borra la tabla de delegación entera, y estas reglas no pueden
// irse con ella.
// ──────────────────────────────────────────────────────────────────────────────

// combosDeCosto is the toggle matrix these two rules must survive.
//
// The important row is "sin delegation": that toggle deletes the
// <!-- active-stack:sdd-delegation --> section outright, and a rule about how to
// spend delegation parked inside it would be absent exactly where somebody
// turned delegation off on purpose.
func combosDeCosto() map[string][]string {
	return map[string][]string{
		"todos los toggles": legacyToggles(),
		"ningun toggle":     {},
		"sin delegation":    sinToggle("delegation"),
		"sin roles":         sinToggle("roles"),
		"sin review":        sinToggle("review"),
		"sin auditor":       sinToggle("auditor"),
		"sin evidence":      sinToggle("evidence"),
	}
}

// exigirEnTodasLasVariantes composes every variant against every toggle combo
// and asserts each phrase is present, comparing normalised content so an edit
// that only adds emphasis does not turn it red — and, more importantly, so a
// phrase re-worded into italics cannot evade the check.
func exigirEnTodasLasVariantes(t *testing.T, requeridas []string, porQue string) {
	t.Helper()
	for _, variante := range variantes {
		for nombre, toggles := range combosDeCosto() {
			t.Run(variante+"/"+nombre, func(t *testing.T) {
				bloque, err := Compose(toggles, variante)
				if err != nil {
					t.Fatalf("Compose(%q) error = %v", variante, err)
				}
				texto := normalizarMarkdown(bloque)
				for _, frase := range requeridas {
					if !strings.Contains(texto, normalizarMarkdown(frase)) {
						t.Errorf("falta %q — %s", frase, porQue)
					}
				}
			})
		}
	}
}

// TestChequeosMecanicosLosCorreElOrquestador pins that deterministic checks are
// run once, by the orchestrator, and their output is handed to the reviewers.
//
// The defect it guards: eight reviewers in verify cost 1.351.186 tokens, almost
// exactly what writing the entire code cost, and four of them arrived at the
// same finding down the same two-command path. A command returns the same thing
// whoever runs it; what is worth a whole context is READING the output, never
// producing it.
//
// The datum that discriminates: the rule must also carry its LIMIT and its
// BOUND, and it fails differently without each.
//
//   - Without "mechanical and deterministic", "run it yourself" reads as licence
//     to pre-chew the judgment too, which is the one thing delegation exists for,
//     and the rule that meant to save tokens deletes the review instead.
//   - Without the bound on OUTPUT SIZE, the rule reverses outright. Running a
//     verbose suite inline puts the whole dump in the orchestrator's own context
//     — precisely what the delegation table's `Bash for execution → delegate` row
//     exists to prevent — and then pastes it into N prompts on top. It only wins
//     while the output is an exit code, a count, a `git diff --stat`.
//
// And it must RECONCILE with that row instead of contradicting it in silence:
// two rules in one block pointing opposite ways get resolved by whichever the
// agent read last. The row is about who executes the work; this rule is about
// not paying N contexts for a datum already in hand. Past the bound, the row
// wins, and the text has to say so.
func TestChequeosMecanicosLosCorreElOrquestador(t *testing.T) {
	exigirEnTodasLasVariantes(t, []string{
		"the mechanical checks you run yourself",
		"the criterion is not in running the command: it is in reading the output",
		"mechanical and deterministic",
		// La cota. Sin ella la regla se da vuelta.
		"the bound, and without it the rule reverses",
		"if the output is a dump, you pass the command and its verdict in one line",
		// La reconciliación con la fila de la tabla que esta regla contradiría.
		"who executes the work",
		"not paying n contexts for a datum you already have",
		"if the output does not fit the bound, the table wins",
	}, "un revisor convocado para ejecutar un comando paga un contexto entero por un exit code")
}

// TestElQueCorrigeSeRepartePorZona pins that findings are grouped by area of the
// code before they are handed out.
//
// Two halves that counterweight each other, and the rule is worthless with
// either one missing:
//
//   - One agent for everything: by finding twelve it drags the context of the
//     eleven before it. Measured at 816.762 tokens, 17% of the whole run.
//   - One agent per finding: each one pays the project's initial read again, and
//     that read costs the same for one finding as for twelve.
//
// The datum that discriminates: the FILE, not the problem. Two findings in the
// same file handed to two agents overwrite each other, which is a defect the
// token count never shows.
func TestElQueCorrigeSeRepartePorZona(t *testing.T) {
	exigirEnTodasLasVariantes(t, []string{
		"group them by zone of the code",
		"the file decides the cut",
		"pays the entrance fee twelve times",
	}, "uno solo arrastra el contexto de los once anteriores y uno por problema paga doce veces la entrada")
}

// TestNingunaReglaNuevaRecortaAlAuditorNiLasLentes is the anti-assert, and it is
// the reason the two rules above can be read at all without doing damage.
//
// Both new rules are about spending less on review. The same measurement that
// produced them found the two things that were WORTH what they cost:
//
//   - The auditor: 267.082 tokens, 5% of the run, and it found three things that
//     none of the eleven reviewers saw. Best ratio of anything measured.
//   - The three adversarial lenses: they are not three passes of the same
//     review. Correctness found a `%` in a password that kills the migrations,
//     wiring found the CI was going to start red, and comments found a text that
//     ordered deleting the file holding a hard rule. Dropping two lenses does not
//     remove a third of the noise: it removes two kinds of bug.
//
// So the guard runs in both directions. The banned phrases must not appear
// anywhere in the composed block, and the new rules must carry the sentence that
// says out loud they license nothing of the sort. One direction alone is not
// enough: a block can add the protective clause and keep the permissive sentence
// three paragraphs down, and each half on its own would still pass.
func TestNingunaReglaNuevaRecortaAlAuditorNiLasLentes(t *testing.T) {
	prohibidas := []string{
		"skip the auditor",
		"drop the auditor",
		"without the auditor",
		"the auditor is optional",
		"fewer lenses",
		"two lenses",
		"drop a lens",
		"skip a lens",
		"one lens is enough",
	}

	for _, variante := range variantes {
		t.Run(variante, func(t *testing.T) {
			bloque, err := Compose(legacyToggles(), variante)
			if err != nil {
				t.Fatalf("Compose(%q) error = %v", variante, err)
			}
			texto := normalizarMarkdown(bloque)
			for _, frase := range prohibidas {
				if strings.Contains(texto, normalizarMarkdown(frase)) {
					t.Errorf("el bloque de %q autoriza recortar la revisión que sí valía "+
						"lo que costaba: contiene %q. El auditor fue el 5%% del gasto y "+
						"encontró tres cosas que once revisores no vieron; las tres lentes "+
						"no son tres pases del mismo review", variante, frase)
				}
			}
			// La cláusula protectora de la regla del base, que ningún toggle borra.
			if !strings.Contains(texto, normalizarMarkdown(
				"the auditor and the three adversarial lenses are not commands")) {
				t.Errorf("la regla de chequeos mecánicos de %q no dice qué NO alcanza: "+
					"una regla de ahorro sin su límite se lee como permiso para ahorrar "+
					"en lo único que rendía", variante)
			}
		})
	}
}
