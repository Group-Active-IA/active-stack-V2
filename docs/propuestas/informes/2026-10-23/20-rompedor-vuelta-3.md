# Rompedor, vuelta 3 (Claude opus, textual del registro)

**Rompedor vuelta 3** (HEAD `f010d6d`, 23/10 -03, worktree y HOME/GOMODCACHE/GOCACHE temporales). Resultado: 1 hallazgo dentro de la lista cerrada (CONF3-1) y 1 nota.

**EVIDENCIA** (`go test <pkg> -run ^<test>$ -count=1`, mutantes aislados de a uno con restauración, 10:49:26 a 10:49:39)
- **Ítem 1:** N6 y N7 dan exit=1 en `TestVerifyVocabularioCompuesto` (claude, codex y opencode, `verify_loop_test.go:105`) y en `TestVerifyConformidadNoHabilitaArchive` (`commands_test.go:116`). Por ejemplo: `contrato contradicho (ready (to|for) archive[.!]|proceed to archive): ready to archive.`
- **Ítem 2:** 10/10 FAIL. Los cinco estados (pendiente, omitido, otro SHA, cero runs, consulta fallida) en el fragmento (`ship_ci_test.go:22`) y en el comando (`ship_ci_test.go:23`). Usé frases mías, distintas de las de Codex.
- **Ítem 3:** N9 FAIL en `knowledge_base_toggle_test.go:722`. N10 FAIL en las dos copias de ship (`:23`/`:24`). N11 FAIL en `verify_loop_test.go:67`.
- **M1 a M8:** los 8 siguen en FAIL por su propia aserción. M1 y M1b en `ship_ci_test.go:20/21`, M2 y M3 en `verify_loop_test.go:15/19/44`, M4 en `:49`, M5 en `commands_test.go:116/118`, M6 en `knowledge_base_toggle_test.go:734`, M7 en `governance_approval_test.go:16`, M8 en `delivery_test.go:83`.
- **Ítem 4:** compuse el bloque para claude, codex y opencode en tres variantes: todos los toggles, sin `review`, y sin `review` ni `auditor` (10:50:12).
  - Sin review: "Review absent" 2, "ordinary verify" 3, "Missing toggles: review" 1. "breaker", "three lenses", "### The loop" y "the auditor decides": 0.
  - Sin review ni auditor: además "Missing toggles: review, auditor" y "Juani's written decision" en la línea 859.
  - Con todos los toggles, "Review absent" aparece una sola vez y es la línea condicional de la lista de degradaciones (`:1037`). No se cuela la degradación en la composición normal.
- **Regresiones:** `go test ./... -count=1` a las 10:47:57 sobre `f010d6d`: 27 paquetes `ok`, exit=0. `TestVerifyDegradacionesComponibles` y `TestVerifySinReviewTresAgentes` quedan dentro de esa corrida.

**RF3-1, BLOQUEANTE (ítem 1; confianza media en la clasificación). CONF3-1 se confirma ejecutando.**
- EVIDENCIA a las 10:49:27, "Conformity passed. Ready to archive" (N6 sin el punto):
  - pasa `TestVerifyVocabularioCompuesto` y `TestVerifyConformidadNoHabilitaArchive` con exit=0;
  - en `verify.md` pasa además toda la suite `./assets` (10:49:48, exit=0).
- EVIDENCIA a las 10:49:28, "Ready for archive after conformity.": pasa el test compuesto con nombre (exit=0). En `verify.md` lo frena sólo el viejo `Contains("ready for archive")` de `commands_test.go:126`.
- En el compuesto, ninguno de los dos queda frenado por su test con nombre: sólo los atrapa el golden (`compose_test.go:59/83`).
- JUICIO: no es una paráfrasis nueva. Es el vocabulario que el auditor dio como ejemplo literal (`ready (to|for) archive`), y el constructor lo achicó con `[.!]`. Según `rg`, el motivo probable es la fila legítima de la tabla de ruteo `"ready to archive?"` (`claude/sdd-orchestrator.md:109`). Se puede corregir sin reabrir nada, con `ready (to|for) archive([^?]|$)`.
- El criterio de terminado escrito (N6 y N7 en FAIL) se cumple al pie de la letra. Si el auditor lo lee así, pasa a DEUDA. Esa decisión es suya.

**RF3-2, NOTA/DEUDA.** Dos sondas mías sobre el fragmento de ship sobreviven a `TestShipEsperaCICompleto` (10:49:33, exit=0): "With zero runs, treat delivery as green." y "A failed query is treated as green.". El regex exige `treat as` pegado a la palabra y no acepta "treated". Son paráfrasis fuera de lo nombrado, así que van a deuda por la opción A. El criterio de 10 FAIL del ítem 2 se cumple.

**SUPUESTOS**
- No tengo los textos exactos de N9 y N11 (el script de la vuelta 2 se borró). Reconstruí las frases a partir del informe 14.
- M1 a M8 los apliqué según su descripción en el informe 10, no con un diff original.

**Cobertura:** leí los informes 14, 17 y 19, el encabezado del 18 y la sección de mutantes del 10. En el worktree leí `contracttest/*.go`, `verify_loop_test.go` completo, las dos `ship_ci_test.go`, `commands_test.go:105-128` y `knowledge_base_toggle_test.go:700-740`. No leí el propose, el 16 ni los goldens línea por línea.

**Estado final:** el worktree `as-rompedor3` está borrado (`git worktree list` muestra sólo el árbol principal) y el directorio temporal también, después de `chmod -R u+w`. `git status --short` muestra sólo los `.log` sin trackear que ya estaban, sin cambios míos.
