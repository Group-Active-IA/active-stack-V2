package config_test

// ──────────────────────────────────────────────────────────────────────────────
// Los documentos de la raíz también envejecen
//
// `TestCatalog_OrchestratorDescriptionIsNotStaleAboutToggles` cuida la
// `long_description` del catálogo: lo que el usuario lee en la TUI mientras
// elige. Nada cuidaba README.md, ARCHITECTURE.md, AGENTS.md ni CLAUDE.md, y los
// cuatro enumeraban CINCO toggles mientras el catálogo ofrecía ONCE. governance,
// evidence, auditor, roles, review, knowledge-base y senior entraron uno por uno
// y ninguno pasó por acá.
//
// Es exactamente el defecto que describe `assets/knowledge-base.md` —un mapa con
// autoridad al que nada del ciclo le vuelve a escribir— en el repo que publica
// ese fragmento. Quien no tiene mapa pregunta; quien tiene uno viejo construye
// arriba.
//
// La lista esperada sale de `config.KnownToggles()`, NUNCA de una constante acá.
// Una lista hardcodeada envejece igual que el documento que vigila, y este repo
// ya tuvo ese caso exacto: el guardia de la descripción del catálogo nació
// verificando UN solo nombre hardcodeado, así que el toggle siguiente lo pasaba
// por al lado.
// ──────────────────────────────────────────────────────────────────────────────

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/config"
)

// docsThatEnumerateToggles son los documentos de la raíz que listan los toggles
// del `sdd-orchestrator`. Hardcodear ARCHIVOS es distinto de hardcodear TOGGLES:
// el set de documentos constitucionales cambia cuando alguien crea o borra uno
// —un acto deliberado y ruidoso— mientras que el catálogo de toggles crece solo,
// que es la deriva que este test existe para cazar.
//
// CHANGES.md queda deliberadamente afuera: es un changelog fechado, con fichas de
// changes ya cerrados. Sus enumeraciones describen el estado del repo en la fecha
// del change; corregirlas sería reescribir historia, no arreglar un mapa.
var docsThatEnumerateToggles = []string{
	"README.md",
	"ARCHITECTURE.md",
	"AGENTS.md",
	"CLAUDE.md",
}

// minToggleNamesForEnumeration es el piso de nombres que convierte una mención en
// una enumeración. Ver el comentario de `isToggleEnumeration`.
const minToggleNamesForEnumeration = 3

// TestRootDocs_ToggleEnumerationsAreNotStale falla si un documento de la raíz
// enumera los toggles y le falta alguno del catálogo real.
//
// Falla también si un documento de la lista dejó de tener enumeración alguna.
// Ese segundo modo no es celo: sin él, el guardia se apaga en silencio el día que
// alguien reescribe una sección y el detector deja de reconocerla, que es el
// mismo final del documento que vigila.
func TestRootDocs_ToggleEnumerationsAreNotStale(t *testing.T) {
	known := config.KnownToggles()
	if len(known) == 0 {
		t.Fatal("config.KnownToggles() is empty: there is nothing to check the docs against")
	}

	for _, name := range docsThatEnumerateToggles {
		path := filepath.Join("..", "..", "..", name)
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("cannot read the root document %s: %v", name, err)
			continue
		}

		found := 0
		for _, st := range markdownStatements(string(raw)) {
			if !isToggleEnumeration(st.text, known) {
				continue
			}
			found++
			for _, tg := range known {
				if namesToggle(st.text, tg) {
					continue
				}
				t.Errorf("%s:%d enumerates the sdd-orchestrator toggles and omits %q: "+
					"the document describes a product the catalogue left behind\n  %s",
					name, st.line, tg, excerpt(st.text))
			}
		}

		if found == 0 {
			t.Errorf("%s has no toggle enumeration at all: either the document stopped "+
				"listing the toggles, or it still lists them in a shape this test no "+
				"longer recognises. A guard that quietly stops guarding is the defect "+
				"it was written against", name)
		}
	}
}

// isToggleEnumeration decide si un fragmento de markdown ENUMERA los toggles o
// apenas menciona alguno al pasar. Dos condiciones, y las dos tienen un
// contraejemplo vivo en este repo:
//
//  1. El fragmento dice "toggle" (o "toggles", o "toggleables"). Sin esto,
//     `CLAUDE.md` §nota —"el orquestador SDD global (coordinación, delegación a
//     sub-agentes, model-routing por fase, protocolo Engram, governance) ya vive
//     en tu ~/.claude/CLAUDE.md"— cae adentro: nombra tres toggles y no está
//     listando el catálogo, está diciendo dónde vive el bloque. Exigirle once
//     nombres sería exigirle once nombres a una frase que habla de otra cosa.
//
//  2. Nombra al menos `minToggleNamesForEnumeration` toggles distintos. Sin esto,
//     `ARCHITECTURE.md` §2 —"`sdd-orchestrator` (orquestador con toggles: TDD,
//     Engram, …)"— cae adentro, y esa celda ELIDE a propósito: el "…" es la
//     abreviatura, no un olvido. Un ejemplo que se declara parcial no miente.
//
// Tres es el piso porque dos nombres todavía se leen como ejemplo ("toggles: TDD,
// Engram, …") y tres ya se lee como lista.
func isToggleEnumeration(text string, known []string) bool {
	if !strings.Contains(text, "toggle") {
		return false
	}
	named := 0
	for _, tg := range known {
		if namesToggle(text, tg) {
			named++
		}
	}
	return named >= minToggleNamesForEnumeration
}

// namesToggle reporta si el texto nombra el toggle como identificador entero.
//
// El chequeo de borde no es cosmético: sin él "tdd" matchea dentro de
// `strict-tdd.md` y "review" dentro de "reviewer", así que un documento pasaría
// el guardia nombrando un ARCHIVO en vez de un toggle. El guion cuenta como
// parte del identificador justamente para eso.
func namesToggle(text, toggle string) bool {
	for i := 0; i+len(toggle) <= len(text); {
		j := strings.Index(text[i:], toggle)
		if j < 0 {
			return false
		}
		start := i + j
		end := start + len(toggle)
		if !isIdentByte(text, start-1) && !isIdentByte(text, end) {
			return true
		}
		i = start + 1
	}
	return false
}

func isIdentByte(text string, i int) bool {
	if i < 0 || i >= len(text) {
		return false
	}
	c := text[i]
	return c == '-' || c == '_' ||
		(c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}

// mdStatement es una afirmación del documento: el texto que un lector lee como
// una unidad, ya desenvuelto del hard-wrap, más la línea donde empieza.
type mdStatement struct {
	line  int
	first string
	text  string
}

// markdownStatements parte el documento en afirmaciones.
//
// No alcanza con ir línea por línea: estos documentos wrappean a ~80 columnas, y
// la lista de toggles de `CLAUDE.md` arranca en una línea y termina en la
// siguiente. Un chequeo por línea le exigiría los once nombres a cada mitad.
//
// Las reglas, y por qué cada una:
//
//   - Una línea en blanco cierra la afirmación. Es el corte de párrafo de
//     markdown.
//   - Un encabezado (`#…`) es una afirmación sola. Sin esta regla "### Resueltas"
//     se pega a la lista que le sigue y la falla se reporta en la línea del
//     título, no en la de la lista.
//   - Una fila de tabla (`|…`) es una afirmación sola. Sin esta regla la tabla
//     entera del README es una unidad, y la celda del `sdd-orchestrator` pasaría
//     el guardia porque OTRA fila —la de kb-creator— escribe "knowledge-base".
//     El lector lee la fila, no la tabla.
//   - Una lista pegada a un pie de párrafo terminado en ":" se une a él. El ":"
//     es lo que hace que la lista sea la continuación de esa frase; separarlas
//     dejaría a la lista de §2.1 de ARCHITECTURE.md sin la palabra "toggle" que
//     la introduce, y por lo tanto invisible para el detector.
//
// El texto se devuelve en minúsculas: los documentos escriben "TDD" y "Engram"
// tanto como `tdd` y `engram`.
func markdownStatements(doc string) []mdStatement {
	var out []mdStatement
	open := -1

	for n, raw := range strings.Split(doc, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			open = -1
			continue
		}
		low := strings.ToLower(line)
		if strings.HasPrefix(line, "#") {
			out = append(out, mdStatement{line: n + 1, first: low, text: low})
			open = -1
			continue
		}
		if strings.HasPrefix(line, "|") {
			out = append(out, mdStatement{line: n + 1, first: low, text: low})
			open = -1
			continue
		}
		if open < 0 {
			out = append(out, mdStatement{line: n + 1, first: low, text: low})
			open = len(out) - 1
			continue
		}
		out[open].text += " " + low
	}

	return mergeListLeadIns(out)
}

// mergeListLeadIns pega cada lista a la frase que la introduce con ":".
func mergeListLeadIns(in []mdStatement) []mdStatement {
	var out []mdStatement
	for _, st := range in {
		if n := len(out); n > 0 && strings.HasSuffix(out[n-1].text, ":") && isListItem(st.first) {
			out[n-1].text += " " + st.text
			continue
		}
		out = append(out, st)
	}
	return out
}

func isListItem(line string) bool {
	for _, m := range []string{"- ", "* ", "+ "} {
		if strings.HasPrefix(line, m) {
			return true
		}
	}
	for i := 0; i < len(line); i++ {
		c := line[i]
		if c >= '0' && c <= '9' {
			continue
		}
		return i > 0 && (c == '.' || c == ')')
	}
	return false
}

func excerpt(text string) string {
	const max = 160
	if len(text) <= max {
		return text
	}
	return text[:max] + "…"
}

func TestDocsSeniorSoloVerify(t *testing.T) {
	for _, nombre := range docsThatEnumerateToggles {
		data, err := os.ReadFile(filepath.Join("..", "..", "..", nombre))
		if err != nil {
			t.Fatal(err)
		}
		encontrada := false
		for n, linea := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(linea, "- `senior`") {
				encontrada = true
				if !strings.Contains(linea, "sólo en verify") {
					t.Errorf("%s:%d: falta senior sólo en verify", nombre, n+1)
				}
			}
		}
		if !encontrada {
			t.Errorf("%s: falta descripción de senior", nombre)
		}
	}
}
