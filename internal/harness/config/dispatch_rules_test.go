package config

import (
	"strings"
	"testing"
)

// ──────────────────────────────────────────────────────────────────────────────
// El terminado observable del reparto
//
// El bloque define siete roles especialistas, y lo que los hace roles y no
// cargos es que su definición de terminado SE PUEDE OBSERVAR. El de QA dice
// "hay al menos un test que falla sin el cambio y pasa con él". Eso se comprueba
// mirando algo.
//
// El orquestador era el único componente del sistema sin un criterio así, y es
// el que más decide: a quién convocar, en qué orden, con qué contexto. Las tres
// reglas que le hablan —qué delegar, qué modelo usar, cómo armar un encargo— son
// orientativas: ninguna se puede verificar después.
//
// Los cinco defectos del reparto, medidos sobre tres corridas reales:
//
//  1. No paralelizó. 39 agentes sumaron 385 minutos contra 360 de reloj: 18
//     minutos de trabajo simultáneo, el 5%. Esperó a un agente el 94% del tiempo.
//  2. Brifeó con rutas en lugar de contexto. Ocho revisores costaron 1.351.186
//     tokens, casi exactamente lo que costó escribir todo el código (1.325.302),
//     porque cada uno recorrió el proyecto entero antes de mirar su rincón.
//  3. Convocó sin criterio. Diecisiete agentes sobre un cambio de andamiaje, y
//     cuatro reportaron el mismo hallazgo por el mismo camino de dos comandos.
//  4. Mandó lo que no aplica. El CLAUDE.md son 1.702 líneas y el bloque es el
//     89%; heredado por 39 agentes son ~780.000 tokens, el 16% de la jornada.
//  5. No corrió él lo que podía correr él: los chequeos mecánicos se ejecutaron
//     una vez por revisor.
//
// Y LO QUE HACE FALTA NO OLVIDAR: las cinco salieron a la luz sólo porque
// alguien preguntó al final. Ninguna dejó rastro mientras ocurría. Por eso el
// quinto punto —el registro— es el que sostiene a los otros cuatro, y por eso
// tiene su propio test más abajo.
//
// POR QUÉ VA EN `## Rules` (inmediatamente después, porque no entra en un
// bullet) Y NO EN UN FRAGMENTO. Es la misma razón que las dos reglas de costo:
// un fragmento desaparece con su toggle, y desaparecería justo en la instalación
// donde nadie eligió esos fragmentos. `## Rules` es la única sección del base que
// ningún toggle borra — el toggle `delegation` borra la tabla de delegación
// entera, y esta sección no puede irse con ella. De ahí que la matriz incluya la
// fila "sin delegation".
// ──────────────────────────────────────────────────────────────────────────────

// criterioDelPunto1 es la frase que dice DÓNDE se verifica el primer punto, y
// dónde no.
//
// La primera versión de esta sección lo mandaba a verificar en el informe del
// agente, y ahí era ciego: nada del formato tipado le pide a un agente que
// declare qué leyó primero, y es justo lo que no tiene ningún motivo para
// ofrecer por su cuenta. Agregarle esa línea al formato tampoco servía —el
// formato lo hereda cada sub-agente, así que cada línea se multiplica por la
// cantidad de agentes de la jornada, que es exactamente lo que esta sección
// viene a achicar.
//
// El encargo sí se puede leer, y no tiene cómo esconder lo que no dice.
const criterioDelPunto1 = "you check it in the prompt you wrote, not in the agent's report"

// alcancePorAgente es lo que ata el punto 1 al punto 5, y sin eso el punto 1
// vuelve a ser incomprobable pasada la sesión: el encargo se lee abriéndolo, y
// el único lugar donde el encargo queda después es el registro.
func alcancePorAgente() []string {
	return []string{
		"the scope per agent is the part that makes point 1 checkable",
		"what it was actually pointed at",
	}
}

// puntosDelReparto son las frases portantes de los cuatro primeros puntos, una
// por punto, con el número que impide que la regla se borre después como obvia.
//
// Se listan aparte del quinto a propósito: el anti-assert de más abajo necesita
// poder preguntar "¿están estos cuatro y falta el registro?", que es exactamente
// la forma en que esta sección se recorta sin que nadie lo note.
func puntosDelReparto() map[string][]string {
	return map[string][]string{
		"1 alcance": {
			"names what the agent is to work on",
			criterioDelPunto1,
			"a prompt cannot hide what it does not say",
			"do not re-paste into the prompt what the agent already inherits",
		},
		"2 paralelo": {
			"launched in the same message",
			"does the second need the first's result",
			"eighteen minutes of simultaneous work, 5%",
			// Lo que el número ES y lo que NO es. Ver el test de abajo.
			"it is the incident that justifies the rule, not a check you can repeat",
			"the observable part is the message",
		},
		"3 convocatoria": {
			"the line of the diff that called for it",
			"this point connects the two instead of restating them",
		},
		"4 chequeos": {
			"the deterministic output every reviewer needed was produced one time",
			"it is part of the dispatch's done",
		},
	}
}

// registroDelReparto son las frases portantes del quinto punto. Tres, y cada una
// cubre un modo distinto de vaciarlo sin borrarlo:
//
//   - la enumeración, sin la cual "dejá registro" no dice de qué;
//   - la frase que dice qué ES el registro, sin la cual se lee como documentación
//     y se saltea como documentación;
//   - el dato que lo sostiene, sin el cual la regla queda sin su porqué y es la
//     primera que alguien borra en la próxima reescritura.
func registroDelReparto() []string {
	return []string{
		"what was launched and with what scope, in what order, what went out in parallel, who was not convened and why",
		"that record is not documentation. it is the evidence that the dispatch happened at all",
		"nobody can tell whether the orchestrator did this job or improvised it",
		"the five failures above came to light only because somebody asked at the end",
	}
}

// TestElRepartoTieneTerminadoObservable exige los cinco puntos en las ocho
// variantes y contra toda la matriz de toggles, incluida la fila que apaga
// `delegation`.
//
// El dato que discrimina: componer con el toggle de delegación APAGADO. Si la
// sección sobrevive a eso, no vive en la sección removible.
func TestElRepartoTieneTerminadoObservable(t *testing.T) {
	for _, variante := range variantes {
		for nombreCombo, toggles := range combosDeCosto() {
			t.Run(variante+"/"+nombreCombo, func(t *testing.T) {
				bloque, err := Compose(toggles, variante)
				if err != nil {
					t.Fatalf("Compose(%q) error = %v", variante, err)
				}
				texto := normalizarMarkdown(bloque)

				for punto, frases := range puntosDelReparto() {
					for _, frase := range frases {
						if !strings.Contains(texto, normalizarMarkdown(frase)) {
							t.Errorf("al reparto le falta el punto %q: no está %q. "+
								"Un criterio que no se puede mirar es una intención, "+
								"y el orquestador ya tenía tres de ésas", punto, frase)
						}
					}
				}
				for _, frase := range registroDelReparto() {
					if !strings.Contains(texto, normalizarMarkdown(frase)) {
						t.Errorf("al reparto le falta el punto 5 (el registro): no está %q. "+
							"Sin registro los otros cuatro vuelven a ser pedidos que nadie "+
							"puede comprobar", frase)
					}
				}
			})
		}
	}
}

// TestElRegistroNoSePuedeRecortarDelReparto es el anti-assert, y es el que hace
// verificables a los otros cuatro puntos.
//
// La forma en que esta sección se degrada no es borrándola: es dejando los
// cuatro puntos que suenan a buena práctica y sacando el quinto, que es el único
// que cuesta trabajo. Los cuatro primeros describen cómo repartir; el registro es
// lo único que después permite saber si se repartió así o se improvisó. Un ciclo
// sin ese registro y un ciclo bien repartido dejan exactamente la misma
// evidencia: ninguna.
//
// Por eso el test no pregunta "¿está el registro?" —eso ya lo hace el de
// arriba—, sino "¿están los otros cuatro Y falta el registro?". Esa combinación
// es el recorte, y tiene que ponerse en rojo nombrándolo.
func TestElRegistroNoSePuedeRecortarDelReparto(t *testing.T) {
	for _, variante := range variantes {
		for nombreCombo, toggles := range combosDeCosto() {
			t.Run(variante+"/"+nombreCombo, func(t *testing.T) {
				bloque, err := Compose(toggles, variante)
				if err != nil {
					t.Fatalf("Compose(%q) error = %v", variante, err)
				}
				texto := normalizarMarkdown(bloque)

				var presentes []string
				for punto, frases := range puntosDelReparto() {
					if strings.Contains(texto, normalizarMarkdown(frases[0])) {
						presentes = append(presentes, punto)
					}
				}
				if len(presentes) == 0 {
					// La sección entera no está: ése es el fallo del test de
					// arriba, y decirlo dos veces sólo agrega ruido.
					return
				}

				for _, frase := range registroDelReparto() {
					if !strings.Contains(texto, normalizarMarkdown(frase)) {
						t.Errorf("el reparto de %q conserva los puntos %v pero le falta el "+
							"registro (%q). Recortado así, la sección queda diciendo cómo "+
							"repartir y sin nada que permita comprobar que se repartió: "+
							"las cinco fallas que la motivaron salieron a la luz sólo "+
							"porque alguien preguntó al final, y ninguna dejó rastro "+
							"mientras ocurría", variante, presentes, frase)
					}
				}
			})
		}
	}
}

// TestElRegistroLlevaElAlcancePorAgente cuida la pieza que hace verificable al
// punto 1, y falla nombrando ESA consecuencia y no la falta de una frase.
//
// La cadena es de tres eslabones y se corta en el del medio sin ruido:
//
//  1. el punto 1 se verifica abriendo el encargo que escribió el orquestador;
//  2. el encargo sólo queda en algún lado si el registro lo guarda;
//  3. el registro lo guarda si lleva, por agente, el alcance que se le dio.
//
// Alguien que recorte el eslabón 3 deja el 1 diciendo "abrí el encargo" sobre un
// encargo que ya no existe pasada la sesión, y el punto vuelve a donde estaba:
// observable en principio, ciego en la práctica. Por eso el test es condicional
// —si está el criterio del punto 1, tiene que estar el alcance en el registro— y
// no un simple "¿está la frase?": el fallo interesante no es la ausencia, es la
// combinación.
//
// Cuida además que la referencia cruzada a `review.md` no quede colgada. La
// sección dice que comparte UNA suposición con el veredicto —que archivar mueve
// el directorio entero del cambio— y que no la repite porque allá ya está
// declarada. Si allá se borra, acá queda un puntero a nada, que es peor que
// haberla repetido.
func TestElRegistroLlevaElAlcancePorAgente(t *testing.T) {
	for _, variante := range variantes {
		for nombreCombo, toggles := range combosDeCosto() {
			t.Run(variante+"/"+nombreCombo, func(t *testing.T) {
				bloque, err := Compose(toggles, variante)
				if err != nil {
					t.Fatalf("Compose(%q) error = %v", variante, err)
				}
				texto := normalizarMarkdown(bloque)

				if !strings.Contains(texto, normalizarMarkdown(criterioDelPunto1)) {
					// Sin el punto 1 no hay cadena que cortar: eso lo reporta
					// `TestElRepartoTieneTerminadoObservable`.
					return
				}
				for _, frase := range alcancePorAgente() {
					if !strings.Contains(texto, normalizarMarkdown(frase)) {
						t.Errorf("el reparto de %q verifica el punto 1 abriendo el encargo "+
							"(%q) pero su registro no guarda el alcance por agente: falta "+
							"%q. El encargo no queda en ningún lado si el registro no lo "+
							"lleva, así que el criterio del punto 1 es incomprobable "+
							"pasada la sesión, que es exactamente de donde se lo sacó",
							variante, criterioDelPunto1, frase)
					}
				}
			})
		}
	}
}

// TestLaSuposicionDelArchivadoEsUnaSola comprueba que la referencia cruzada de
// la sección apunta a algo que existe.
//
// Corre sólo con el toggle `review` activo, que es la única condición bajo la
// cual la sección afirma que la suposición "ya está declarada allá".
func TestLaSuposicionDelArchivadoEsUnaSola(t *testing.T) {
	for _, variante := range variantes {
		t.Run(variante, func(t *testing.T) {
			bloque, err := Compose(KnownToggles(), variante)
			if err != nil {
				t.Fatalf("Compose(%q) error = %v", variante, err)
			}
			texto := normalizarMarkdown(bloque)

			if !strings.Contains(texto, normalizarMarkdown(
				"this record shares one assumption with the review verdict")) {
				t.Errorf("el reparto de %q no dice que comparte la suposición del archivado "+
					"con el veredicto. Declarada dos veces por separado, se lee como dos "+
					"suposiciones del sistema cuando es una", variante)
			}
			if !strings.Contains(texto, normalizarMarkdown(
				"openspec archive` moves the whole change directory")) {
				t.Errorf("el reparto de %q referencia una suposición que ya no está "+
					"declarada en ningún lado: `review.md` la perdió y acá quedó un "+
					"puntero a nada, que es peor que haberla repetido", variante)
			}
		})
	}
}

// TestElRepartoNoAutorizaRecortarLaRevision es el segundo anti-assert, y corre
// en las dos direcciones.
//
// Los cinco puntos hablan de gastar mejor el reparto, y una regla de ahorro sin
// su límite se lee como permiso para ahorrar en lo único que rendía. El auditor
// y las tres lentes ya están cubiertos por
// `TestNingunaReglaNuevaRecortaAlAuditorNiLasLentes`, que prohíbe esas frases en
// todo el bloque compuesto. Lo que ese test NO cubre es el TDD, que esta sección
// sí nombra: la mitad cara del ciclo es el test rojo que viste fallar, y es
// justo la que se recorta cuando alguien lee "gastar menos".
//
// Una sola dirección no alcanza: el bloque podría agregar la cláusula protectora
// y dejar la frase permisiva tres párrafos más abajo, y cada mitad por separado
// pasaría igual.
func TestElRepartoNoAutorizaRecortarLaRevision(t *testing.T) {
	prohibidas := []string{
		"the tdd cycle is optional",
		"skip the tdd cycle",
		"drop the tdd cycle",
		"tdd is dispatch overhead",
		"the red test is dispatch overhead",
	}

	for _, variante := range variantes {
		t.Run(variante, func(t *testing.T) {
			bloque, err := Compose(KnownToggles(), variante)
			if err != nil {
				t.Fatalf("Compose(%q) error = %v", variante, err)
			}
			texto := normalizarMarkdown(bloque)

			for _, frase := range prohibidas {
				if strings.Contains(texto, normalizarMarkdown(frase)) {
					t.Errorf("el bloque de %q autoriza recortar el ciclo que sí valía lo "+
						"que costaba: contiene %q. El reparto se gasta mejor convocando "+
						"lo mismo con mejor encargo, no convocando menos", variante, frase)
				}
			}
			if !strings.Contains(texto, normalizarMarkdown(
				"the auditor, the three adversarial lenses and the tdd cycle are not dispatch overhead")) {
				t.Errorf("el reparto de %q no dice qué NO alcanza. Sin esa cláusula, cinco "+
					"puntos sobre gastar mejor se leen como permiso para revisar menos, "+
					"que es la lectura inversa a la que tienen", variante)
			}
			if !strings.Contains(texto, normalizarMarkdown(
				"read this section as a reason to convene fewer reviewers and you have read it backwards")) {
				t.Errorf("el reparto de %q no nombra su lectura inversa. La regla falla "+
					"por ahí: quien lee cinco puntos de eficiencia concluye que sobran "+
					"revisores, no que sobran recorridas del proyecto", variante)
			}
		})
	}
}
