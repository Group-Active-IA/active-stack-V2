# Rompedor del paso 0, vuelta 1 (Claude opus, textual del registro)

RP0, rompedor de verify del paso 0 `contratos-fase2-deuda`, vuelta 1. Diff `7f1ef88..HEAD` (a6cfc2d). Encontré **2 BLOQUEANTES y 2 NOTAS**. Bloque compuesto para claude, codex y opencode con review+senior, sin auditor, y con senior sin review. Árbol principal sin tocar. Worktree, HOME, GOCACHE y GOMODCACHE temporales ya borrados y `git worktree prune` hecho. Los `.log` sin trackear en `docs/propuestas/informes` ya estaban y no son míos.

**RP0-1 BLOQUEANTE (correctitud, CONF-1 confirmado)**
- EVIDENCIA: compuse claude con review+senior+knowledge-base y sin auditor (11:48). Salen juntas tres reglas que chocan:
  - L565 (`review.md:104`): "Only APROBADO from the auditor … enables archive", sin condición.
  - L567 (`review.md:106`): "A round … closes with the auditor verdict".
  - L921 (`senior-without-auditor.md:7`): "Do not wait for a verdict from an absent auditor"; avanza la decisión escrita de Juani.
- EVIDENCIA: lo mismo pasa en codex (L527/529 contra L883) y en opencode (L567/569 contra L923).
- JUICIO (rompedor, confianza alta): son dos autoridades de cierre para el mismo caso, justo la deuda RF1-4/RF2-4 que el paso 0 tenía que pagar. Una ronda sin auditor nunca cierra y el contador de vueltas no avanza. El comando `assets/commands/claude/opsx/verify.md:144` sí acota con "in the complete cycle"; el fragmento no.

**RP0-2 BLOQUEANTE (verdad de los tests)**
- EVIDENCIA, mutantes de inversión con redacción mía, `go test ./internal/harness/config -count=1`:
  - **M1** 11:49:39, `senior-without-auditor.md:3`: «It has no closing authority.» → «Its report is the closing verdict whenever no auditor is configured.» Resultado: `ok 26.781s`, **sobrevive**.
  - **M2** 11:50:07, `senior-without-auditor.md:7`: «Do not wait … return to apply.» → «Keep the round open until the auditor returns its verdict; meanwhile send blockers back to apply.» Resultado: `ok 26.612s`, **sobrevive**.
  - **M3** 11:50:34, `review.md`: «the breaker cannot authorize closure.» → «the breaker's report is enough to close the change.» Lo atrapan sólo los goldens `TestCompose_AllToggles` y `TestCompose_AllToggles_opencode`.
  - **M4** 11:51:01, `review.md`: «The auditor classifies … the orchestrator executes» → orchestrator y auditor intercambiados. **Muerto por aserción** en `TestCompose_ReviewToggle_ReportAuthority` (más los goldens).
  - **M5** 11:51:31, `knowledge-base.md:191`: «do not invent its lists» → «fill its lists from your own reading of the diff». Sólo los goldens.
  - **M6** 11:52:02, `review.md`: «closure stops for Juani's written decision on degraded scope» → «closure proceeds on the senior's report alone». Sólo los goldens.
- JUICIO (confianza alta): 1 de 6 muere por una aserción atribuible al agregado. Los goldens sólo muestran un diff, y regenerarlos acepta la inversión. `senior-without-auditor.md` no entra en ningún golden porque el golden tiene el auditor encendido, así que en el fragmento nuevo la regla "el senior decide" (N3) y la espera imposible pasan en verde. El chequeo `do not start.*until the auditor` calca una sola frase (consistente con CONF-2).

**RP0-3 NOTA (CONF-3 confirmado, latente)**
- EVIDENCIA 11:52:45: con el original limpio, `debtViolations` da 0 en el bloque claude completo, 0 en `verify.md` y 0 en `ship.md`. Hoy el corpus aditivo es atribuible.
- EVIDENCIA 11:52:59: agregué "The review does not block." a `review.md`. Las 28/28 subpruebas de `TestDebtCorpusAdditiveCopies` pasan igual, así que la infracción previa enmascara todas las mutaciones. Sólo falla `TestDebtContractsComposition/claude/15`. Falta exigir original limpio dentro del mismo test.

**RP0-4 NOTA**
- `review.md:50` conserva el título "### It reports, it does not block", y `reviewSection(t, "### it reports")` depende de él.
- JUICIO (confianza media): el patrón no lo atrapa porque falta "review" en la cláusula. El sentido es compatible ("el rompedor no tiene veto"), pero es el mismo vocabulario que la deuda.

**Cableado**
- EVIDENCIA: `compose.go:192` elige `senior-without-auditor.md` en claude, codex y opencode, y aparece en los tres bloques (L915, L877, L917).
- JUICIO (confianza alta): el cableado es correcto.
- Sin golden de codex.
- El test aditivo sólo lee las copias de los comandos de claude; opencode no tiene `opsx/verify` ni `opsx/ship`.

**Supuestos**
- SUPUESTO: tomé el alcance de deuda de las líneas 84, 102, 112-113 y 125 del propose; no leí el informe del constructor.
- SUPUESTO: no corrí la suite entera, sólo `internal/harness/config`.

**Cobertura**
- Leí: el diff completo de los 9 archivos, la tabla de deuda del propose, CONF-1 a CONF-3, `verify.md`, y los bloques compuestos de claude, codex y opencode (review+senior y senior solo).
- No abrí ship ni N1/N1b (pendiente de Juani).
