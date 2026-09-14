package config

import (
	"strings"
	"testing"
)

// ──────────────────────────────────────────────────────────────────────────────
// El punto 3 del reparto nombra al comando que lo ejecuta
//
// Un mecanismo que nadie invoca es lo mismo que no tenerlo, y en este repo ya
// costó caro dos veces en una semana: un bloque dirigido a un agente que ningún
// instalador creaba (`TestBloqueSeDirigeALaSesionPrincipal`) y una skill cuyos
// disparadores declarados no existían como comandos. Las dos veces el texto
// estaba bien escrito y no tenía lector.
//
// `active-stack roles` lee el diff y devuelve qué revisores corresponden, con el
// motivo de cada uno, incluidos los que NO van. Le saca al orquestador la
// decisión de a quién convocar. Si el bloque no lo nombra, el comando existe y
// nadie lo corre, y el punto 3 vuelve a ser lo que era: un criterio.
//
// POR QUÉ ESO IMPORTA, medido y no opinado: escrito como criterio, en tres
// corridas el orquestador convocó siempre los mismos cuatro roles de diez y el
// revisor senior no corrió ni una vez. El defecto es estructural —el que juzga
// si una capa "se toca" es el mismo que se ahorra el trabajo de convocar, y no
// convocar no deja rastro—, así que no se arregla redactando mejor el criterio.
// Se arregla reemplazándolo por un procedimiento.
//
// POR QUÉ VA CONTRA TODA LA MATRIZ DE TOGGLES, incluida la fila que apaga
// `delegation`: el punto 3 vive en `## The Dispatch`, que está en la única parte
// del base que ningún toggle borra. Componer con `delegation` apagado es el dato
// que discrimina "está en la sección que sobrevive" de "está en la que no".
// ──────────────────────────────────────────────────────────────────────────────

// elComando es el nombre exacto que hay que poder tipear. Se chequea el
// invocable completo y no sólo "roles": la palabra suelta ya aparece media docena
// de veces en el bloque (el fragmento `roles`, los siete roles, el registro de
// roles no convocados), así que un assert sobre "roles" pasaría en verde con el
// comando sin nombrar en ninguna parte, que es exactamente el estado que este
// test existe para poner en rojo.
const elComando = "active-stack roles"

// procedimiento son las frases que convierten el punto 3 de criterio en
// procedimiento. Una por pieza, y cada una cubre un modo distinto de vaciarlo:
//
//   - que el comando se corra Y que su salida sea la convocatoria, sin lo cual
//     nombrarlo es una mención y no una instrucción;
//   - que el orquestador LEA en lugar de decidir, que es el punto entero;
//   - el dato medido que lo justifica, sin el cual la regla queda sin porqué y
//     es la primera que alguien borra en la próxima reescritura como obvia;
//   - el motivo estructural, sin el cual la lectura natural es "convocá con más
//     ganas", que es la que ya se midió fallando.
func procedimientoDelPunto3() []string {
	return []string{
		"the roles it returned are the roles that went out",
		"the orchestrator does not decide here, it reads",
		"the same four roles out of ten, every time",
		"whoever judges whether a layer \"is touched\" is the same one who saves the work of convening",
	}
}

// losQueNoVan son las frases de la otra mitad del comando: el rol no convocado
// sale con su motivo, y esa lista es la que llena el registro del punto 5.
//
// Se chequea aparte del procedimiento porque es la mitad que se pierde primero:
// quien recorta deja "corré el comando y convocá lo que dice" y se lleva puesto
// el motivo de los que no fueron, que es justo lo que no se puede reconstruir
// después. Una ausencia sin motivo es indistinguible de un olvido.
func losQueNoVanDelPunto3() []string {
	return []string{
		"the roles that do not go out come back too, each with its reason",
		"indistinguishable from having forgotten it",
		"the same list that fills the \"who was not convened and why\" of the record in point 5",
	}
}

// losDosQueNoDecideElDiff cuida la distinción que el comando hace y que el texto
// tiene que respetar: hay roles que el diff no decide.
//
// El arquitecto y el analista funcional se convocan por fase, no por un archivo,
// y el comando devuelve `not-decided-by-diff` en lugar de afirmar que no
// corresponden. Borrar esta distinción tiene una consecuencia concreta y en una
// sola dirección: el texto pasaría a decir que la lista cubre los siete, el
// orquestador dejaría de convocar a esos dos por fase, y `roles.md` —que es la
// especificación— los pide siempre en propose, explore y archive.
//
// La única excepción es el arquitecto bajo el piso, y el texto la nombra: ahí
// `Classify` sí lo devuelve convocado. Describirlo como "el diff nunca decide al
// arquitecto" sería más prolijo y sería falso, y una descripción falsa del
// comando dentro del bloque es exactamente la falla que este archivo cuida.
func losDosQueNoDecideElDiff() []string {
	return []string{
		"five of the seven come off the list outright",
		"not-decided-by-diff",
		"declining to assert something it never looked at",
	}
}

// elComandoNoDisponible es la respuesta honesta a "¿y si el binario no está en
// el PATH, o el proyecto no es un repo git?".
//
// No es "entonces decidí vos como antes": eso devuelve el punto 3 al criterio
// que se midió fallando, y por la puerta de atrás, en la corrida donde menos se
// nota. Es que el reparto se hace igual y el registro dice que se hizo sin la
// lista, porque las dos versiones dejan las mismas líneas escritas y no son la
// misma afirmación.
func elComandoNoDisponible() []string {
	return []string{
		"when the command is not available you do not fall back to deciding: you declare it",
		"that it happened without the list",
		"are the same four lines on the page and are not the same claim",
	}
}

// TestElPunto3NombraAlComandoQueLoEjecuta exige el comando y su procedimiento en
// las ocho variantes y contra toda la matriz de toggles.
func TestElPunto3NombraAlComandoQueLoEjecuta(t *testing.T) {
	for _, variante := range variantes {
		for nombreCombo, toggles := range combosDeCosto() {
			t.Run(variante+"/"+nombreCombo, func(t *testing.T) {
				texto := componerNormalizado(t, toggles, variante)

				if !strings.Contains(texto, normalizarMarkdown(elComando)) {
					t.Fatalf("el reparto de %q no nombra %q. El comando existe, le saca al "+
						"orquestador la decisión de a quién convocar, y nada le dice que lo "+
						"corra: un mecanismo que nadie invoca es lo mismo que no tenerlo",
						variante, elComando)
				}
				for _, frase := range procedimientoDelPunto3() {
					if !strings.Contains(texto, normalizarMarkdown(frase)) {
						t.Errorf("el punto 3 de %q nombra el comando pero no es un "+
							"procedimiento: falta %q. Mencionarlo sin decir que su salida ES "+
							"la convocatoria lo deja como un criterio con una herramienta al "+
							"lado, y como criterio ya se midió convocando los mismos cuatro "+
							"roles de diez en tres corridas", variante, frase)
					}
				}
				for _, frase := range losQueNoVanDelPunto3() {
					if !strings.Contains(texto, normalizarMarkdown(frase)) {
						t.Errorf("el punto 3 de %q pierde la mitad del valor del comando: "+
							"falta %q. El rol que no se convoca no aparece hoy en ningún "+
							"lado, y esa ausencia es indistinguible de un olvido", variante, frase)
					}
				}
				for _, frase := range losDosQueNoDecideElDiff() {
					if !strings.Contains(texto, normalizarMarkdown(frase)) {
						t.Errorf("el punto 3 de %q borra la distinción que el comando hace: "+
							"falta %q. El arquitecto y el analista funcional se convocan por "+
							"fase y no por un archivo; leer la lista como si cubriera los "+
							"siete es dejar de convocarlos, y `roles.md` los pide siempre",
							variante, frase)
					}
				}
				for _, frase := range elComandoNoDisponible() {
					if !strings.Contains(texto, normalizarMarkdown(frase)) {
						t.Errorf("el punto 3 de %q no dice qué hacer sin el comando: falta "+
							"%q. El binario puede no estar en el PATH y el proyecto puede no "+
							"ser un repo git, y sin esa respuesta el hueco se llena solo con "+
							"la que el punto vino a sacar", variante, frase)
					}
				}
			})
		}
	}
}

// TestElPunto3NoAutorizaConvocarSinCorrerElComando es el anti-assert, y corre en
// las dos direcciones.
//
// La forma en que este punto se degrada no es borrando el comando: es dejándolo
// nombrado y agregando tres palabras que lo vuelven opcional. "Corré el comando,
// o usá tu criterio si te parece" nombra el comando, pasaría el test de arriba
// entero, y devuelve el punto exactamente a donde estaba —porque el criterio es
// gratis y correr el comando no—, sin que el diff se vea sospechoso.
//
// Una sola dirección no alcanza. Prohibir las frases permisivas deja pasar un
// texto que simplemente no dice que la lista sea vinculante, y exigir la
// cláusula vinculante deja pasar un texto que la tiene y tres párrafos más abajo
// se desdice. Las dos mitades por separado pasan en verde; juntas, no.
func TestElPunto3NoAutorizaConvocarSinCorrerElComando(t *testing.T) {
	// Redacciones que convierten el procedimiento en una sugerencia. No
	// pretenden ser exhaustivas —ninguna lista de frases lo es—, sino cubrir las
	// formas en que un procedimiento se ablanda: declararlo opcional, permitir
	// saltearlo, o devolverle la decisión al criterio de quien convoca.
	prohibidas := []string{
		"running the command is optional",
		"the command is optional",
		"you may skip the command",
		"you can skip the command",
		"skip the command",
		"use your judgment if you prefer",
		"use your own judgment instead",
		"decide who to convene yourself",
		"if you prefer, decide",
		"or decide it yourself",
		"at your discretion",
	}

	for _, variante := range variantes {
		for nombreCombo, toggles := range combosDeCosto() {
			t.Run(variante+"/"+nombreCombo, func(t *testing.T) {
				texto := componerNormalizado(t, toggles, variante)

				for _, frase := range prohibidas {
					if strings.Contains(texto, normalizarMarkdown(frase)) {
						t.Errorf("el bloque de %q autoriza convocar sin correr el comando: "+
							"contiene %q. Con esa puerta abierta el punto 3 vuelve a ser un "+
							"criterio, porque el criterio es gratis y correr el comando no, "+
							"y el criterio ya se midió convocando los mismos cuatro roles de "+
							"diez en tres corridas", variante, frase)
					}
				}

				if !strings.Contains(texto, normalizarMarkdown(elComando)) {
					// El comando no está en ninguna parte: ése es el fallo de
					// `TestElPunto3NombraAlComandoQueLoEjecuta`, y decirlo dos
					// veces sólo agrega ruido.
					return
				}
				vinculante := []string{
					"the roles it returned are the roles that went out",
					"the orchestrator does not decide here, it reads",
				}
				for _, frase := range vinculante {
					if !strings.Contains(texto, normalizarMarkdown(frase)) {
						t.Errorf("el bloque de %q nombra %q pero no dice que su salida sea "+
							"vinculante: falta %q. Un comando nombrado sin esa cláusula es "+
							"una herramienta sugerida, y el punto 3 ya fue eso", variante,
							elComando, frase)
					}
				}
			})
		}
	}
}

// componerNormalizado compone una variante y devuelve el texto listo para
// comparar CONTENIDO y no formato. Ver `normalizarMarkdown`: sin eso, envolver
// una frase en **negrita** pondría en rojo un assert positivo por un cambio
// inofensivo y —peor— dejaría en verde un anti-assert con la frase prohibida
// todavía ahí, en itálica.
func componerNormalizado(t *testing.T, toggles []string, variante string) string {
	t.Helper()
	bloque, err := Compose(toggles, variante)
	if err != nil {
		t.Fatalf("Compose(%q) error = %v", variante, err)
	}
	return normalizarMarkdown(bloque)
}
