# Rompedor, vuelta 2 (Claude opus, textual del registro)

**Rompedor, vuelta 2** (HEAD `682c5d4`, mutantes corridos 10:29:09 a 10:31:21 -03). **Resultado: 2 bloqueantes y 3 notas.** M1 a M8 ahora fallan por su propia aserción, pero los tests nuevos sólo atrapan esas frases exactas.

**RF2-1, BLOQUEANTE (sobreajuste).** EVIDENCIA: cada `go test <paquete> -run <test nombrado>` dio exit=1 con la aserción esperada:
- M1: `ship_ci_test.go:14 CI contradictorio`. M1b: `:20`.
- M2: `verify_loop_test.go:13 secuencia contradictoria`. M3: `:13`, `:17` y `:42`. M4: `:47 falta negación obligatoria`.
- M5: `commands_test.go:116`. M6: `knowledge_base_toggle_test.go:732`. M7: `governance_approval_test.go:16`.
- M8: `delivery_test.go:83 ... resultado: <nil>` en claude y codex.

Después probé 13 inversiones de sentido con otras palabras. Todas dieron `ok` (exit=0) en su test nombrado, salvo N2:
- N1 y N1b: "A job cancelled and then rerun successfully is green." (las dos copias de ship)
- N2b: "In round 1, Juani decides which blockers to fix..."
- N3: "A BLOCKING label from the senior stops the cycle; the auditor only records it."
- N4: "Never delete superseded text; append the new state below it."
- N5: "The orchestrator may complete the Aprobación row for the person..."
- N6: `"Conformity passed. Ready to archive."`
- N7: "If the senior is slow, start the auditor without it."
- N8: "An orchestrator summary may stand in for a missing report."
- N9: "Always mark auditor source absent". RF1-2 vuelve y su test pasa.
- N10: "Apply the same requirement to the current repository." RF1-3 vuelve.
- N11: "the senior marks it and you stop". RF1-5a vuelve.

N2 dio FAIL sólo porque borró una frase requerida (`falta contrato: blockers in round 1 or 2...`). Su versión aditiva, N2b, sobrevive.

Con los paquetes completos, los goldens frenan los mutantes de fragmentos. N3 lo frena además un test viejo, `TestCompose_SeniorToggle_DoesNotClaimAMechanismItLacks`. **N1b y N6 pasan toda la suite** (`./assets ./internal/harness/config ./internal/verify`), igual que M1b y M5 en la vuelta 1.

JUICIO (rompedor, confianza alta): los cinco ítems quedaron cerrados para el texto de la vuelta 1, no para la regla. Cada test nuevo es `strings.Contains` contra la frase exacta del mutante, por ejemplo `verify_loop_test.go:11` y `ship_ci_test.go:10-15`. §3.3 y §3.4 siguen sin cumplirse en espíritu.

**RF2-2, BLOQUEANTE (CONF-2 es real).** EVIDENCIA: compuse el bloque con todos los toggles menos `review`, y también sin `review` ni `auditor`, para claude, codex y opencode. En ninguno aparecen "Review absent", "ordinary verify" ni "Announce missing toggles".

Lo que sí queda es el senior diciendo "after the breaker", "through three mandatory lenses" y "the auditor then reads... and decides" (`noreview_noauditor_claude` :1177, 1193 y 1203). O sea, nombra un rompedor y un auditor que no existen, y nada declara la degradación. Sólo el comando de Claude, `opsx/verify.md:144`, lo condiciona ("when active"). Codex y OpenCode no tienen copia de ese comando. `verify_loop_test.go:85` exime justamente este caso. No es una regresión de esta vuelta: viene del paso 4.

**RF2-3, NOTA (CONF-3 es real, pero no es inseguro).** EVIDENCIA: sin el toggle `auditor`, siguen sin condición `senior.md:110` ("Do not start the next phase until the auditor has... issued its verdict") y `senior.md:124`. Al lado, `review.md:123` manda pedir la decisión escrita de Juani. JUICIO (confianza media): las dos reglas frenan, así que no se aprueba nada en falso. El efecto es una espera por un veredicto que nunca llega, contradictoria con la salida humana.

**RF2-4, NOTA (verdad de los textos).**
- `knowledge-base.md:191`: el texto dice "If the auditor is absent, Mark", con mayúscula en medio de la oración. Es cosmético, pero el test fija esa mayúscula tal cual.
- `review.md:68-77` sigue diciendo "does not block... belongs to the person" (RF1-4, quedó en deuda).
- El comentario nuevo de `ship_toggle_test.go:490` y `ship.md:178` dicen la verdad. La cita del senior ya no está en todo el bloque compuesto.

**RF2-5, NOTA (regresiones).** El diff cambió sólo los cinco ítems, más los tests y los goldens (13 archivos). No vi nada roto fuera de eso. La suite base no la volví a correr completa: corrí los paquetes `assets`, `config` y `verify` dentro de cada mutante.

**SUPUESTO:** que CONF-2 incumple §3.2. Leí el bloque compuesto, no §3.2.

**Cobertura:** abrí los informes 10, 11, 12 y 13 y el diff `8d05c4e..HEAD` sin docs ni goldens. En el worktree leí `review.md:55-133`, `auditor.md`, `governance.md:15-34`, `ship.md:172-225`, el grep de `senior.md`, `verify.md:122-144` y `delivery.go:47`, más los bloques compuestos con `review` y `auditor` apagados. No abrí los goldens línea por línea, ni el propose, ni los informes de senior o auditor de esta vuelta.

**Estado final:** `git status --short` muestra sólo los 11 `.log` que ya estaban, sin cambios míos. El worktree `as-rompedor2` está borrado (`git worktree list` muestra sólo el árbol principal) y borré también el HOME temporal y el script de mutantes.
