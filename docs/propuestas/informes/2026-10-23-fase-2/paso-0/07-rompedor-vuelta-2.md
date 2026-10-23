# Rompedor del paso 0, vuelta 2 (Claude opus, textual del registro)

**RP0b, rompedor de verify del paso 0, vuelta 2 de 3.** Corrí sobre HEAD `1e67e89`, en un worktree propio con HOME, GOCACHE y GOMODCACHE temporales, ya borrados, con `git worktree prune` hecho. No toqué el árbol principal. Los `.log` sin trackear ya estaban antes. Hay 1 BLOQUEANTE y 1 NOTA.

**RP0b-1 BLOQUEANTE (punto 2 de la lista cerrada, no cumplido)**
- EVIDENCIA. Script `mut.py`: cada mutación parte de `git checkout -- assets` y después corre `go test ./internal/harness/config -count=1`.
  - M1, M2, M3, M5 y M6 mueren por `TestDebtContractsComposition` y `TestAuthorityFragmentMutations`, con el mensaje de su clase. M4 muere por `TestCompose_ReviewToggle_ReportAuthority`. Horas: M1 12:07:12, M2 12:07:41, M3 12:08:09, M4 12:08:37, M5 12:09:04, M6 12:09:32.
  - **Mis 8 mutantes nuevos, por sustitución (12:10 a 12:13).** `TestDebtContractsComposition` pasa en los 8. El único fallo es `TestAuthorityFragmentMutations` con "punto de mutación ambiguo": mi reemplazo borró el ancla de su tabla. Es un fallo por ancla, no una aserción semántica.
  - **Los mismos 8 por inserción, sin borrar el original (12:13:48 a 12:17:15).** En `senior-without-auditor.md`, N1a a N4b dan `ok`. En review.md y knowledge-base.md, N5c a N8d fallan sólo por golden (`TestCompose_AllToggles`, `TestCompose_AllToggles_opencode`), con `asercion=[]`.
  - **Insertados en el bloque compuesto** sin auditor de claude, codex y opencode, a las 12:17:54: `authorityViolations(block+m, true)` da `hits=0` en los 24 casos.
  - Los mutantes:
    - Clase a: N1a «El senior cierra la ronda cuando no hay auditor.» y N2a «The senior's word settles the round and the change ships.»
    - Clase b: N3b «Hold the change until the auditor weighs in, even if it is switched off.» y N4b «Mantené la ronda abierta hasta que hable el auditor.»
    - Clase c: N5c «A clean breaker report lets the orchestrator ship the change.» y N6c «The breaker, who cannot be overruled, decides that the change is done.»
    - Clase d: N7d «Write its lists yourself.» y N8d «Make up LECCIONES and DEUDA from the diff.»
- JUICIO (rompedor, confianza alta). La aserción es una lista cerrada de verbos en inglés:
  - No reconoce "cierra", "settles", "ships", "decides", "done", "write" ni "make up".
  - La clase b exige "verdict|decision|returns", así que "weighs in" y "hable" pasan.
  - La negación cubre la cláusula entera: "cannot" en una subordinada ("who cannot be overruled") apaga la detección de la principal.
  - Mataron 0 de 8, y el mínimo era 4. Los 6 originales mueren porque la tabla del test calca su redacción, lo mismo que pasaba con CONF-2.

**Punto 1: cumplido**
- EVIDENCIA. Compuse claude, codex y opencode con kb, review, senior y engram, sin auditor, a las 12:06:12:
  - `rg -i '(^|\. )only APROBADO'` da 0 en los tres bloques.
  - Toda exigencia de APROBADO o veredicto del auditor lleva "With auditor active". Con auditor ausente, el cierre va a la decisión escrita de Juani (claude L567, L921).
  - Con el auditor encendido, el requisito sigue: "with auditor active, only APROBADO…enables archive" y "closes with the auditor verdict" aparecen 1 vez cada uno en los tres bloques.

**Punto 3: cumplido**
- EVIDENCIA, 12:17:59. `debt_contract_test.go:222: original no limpio: [review sin autoridad: The review does not block]` y después `--- FAIL: TestDebtCorpusAdditiveCopies`.

**RP0b-2 NOTA**
- EVIDENCIA. Sin auditor, el bloque compuesto dice en claude L592: "A full round has conformity, breaker, senior and auditor reports; omit only the senior…", y en L588: "A missing required report stops closure".
- JUICIO (confianza media). Remite a "declare other degradation as above", así que no exige veredicto. Igual, una lectura literal toma al auditor ausente como un informe requerido que falta.

**Regresiones**
- EVIDENCIA, 12:18:00 y 12:19:13. Con el árbol limpio, `go test ./... -count=1` da 26 `ok` y 0 líneas que no sean ok.

**SUPUESTO.** Mi clasificación "por aserción" excluye los goldens y la falla del ancla. No leí 05 ni 06.

**Cobertura.** Leí 04 (la lista cerrada), 03, el diff `c9dd12b..HEAD` (review.md y los tests), `senior-without-auditor.md` y los 6 bloques compuestos. No abrí `ship.md` ni las copias de comandos de codex y opencode.
