# Constructor — paso 0, vuelta 2 de 3

Fecha: 2026-10-23T12:02:50.142033-03:00. Change: `contratos-fase2-deuda`.
Familia declarada: Codex / GPT-6. Rol: CONSTRUCTOR. Base: `c9dd12b29bc39a216d3b80adfd7fe2ab973199b6`. Implementación: `32612f96a17e9db4da388e01819b3f7721caf36e`.

## Alcance cerrado y fuentes

Leí AGENTS.md, el auditor `04-auditor-vuelta-1.md` completo (incluida su lista cerrada), el rompedor `03-rompedor-vuelta-1.md`, `tasks.md`, `compose.go` (selección de review y senior), `review_toggle_test.go`, `debt_contract_test.go`, los fragmentos review, senior-without-auditor y el contrato de archive de knowledge-base. Consulté `openspec status --change contratos-fase2-deuda`: planning completo, specs omitidas por declaración del change. No convoqué revisores ni actué como auditor.

Se corrigen exactamente los tres puntos. N1, N1b, RP0-4/paso S y la reinstalación CONF-1 del propose quedan fuera. No hay build, push, instalación ni escrituras en ~/.claude, ~/.local/bin o .env. Los logs ajenos sin trackear se conservaron.

## Punto 1 — composición sin espera imposible

EVIDENCIA: `review.md` acota a **With auditor active** la autoridad, APROBADO para archive, deuda justificada, cierre de ronda y retorno de la lista a apply. Cuando falta auditor, el cierre se detiene ante la decisión escrita de Juani: habilitar auditor, aceptar alcance degradado justificado o frenar. No se atribuye esa decisión al senior ni al rompedor. Se preserva el orden conformidad → rompedor → senior → auditor cuando está activo.

Consumidores: `Compose` carga el mismo fragmento para claude, codex y opencode; no se cambió el selector. `TestDebtContractsComposition` inspecciona el cuerpo completo, con 16 combinaciones por variante, e incluye ambos estados del auditor. Con review y auditor activos exige expresamente el requisito de APROBADO para archive. Sin auditor, la clase b busca esperas/requisitos sobre su veredicto en todas las cláusulas; otro fragmento con el fallback correcto no borra un hit. Se actualizaron sólo los dos goldens afectados por el texto.

## Punto 2 — clases, no igualdad con las frases mutantes

EVIDENCIA: `authorityViolations` combina familias de vocabulario por cláusula, normaliza mayúsculas y formato, e identifica también sujetos informe/pronombre y construcciones pasivas. No busca igualdad con los casos de la tabla:

| Clase | Familias verificadas | Mensaje |
|---|---|---|
| a | senior/historian/its report + cierre, aprobación, archive, veredicto, autorización y equivalentes | clase a: senior con autoridad |
| b | auditor + esperar, mantener abierto, until, pending o requisito de cierre + veredicto/aprobación/decisión | clase b: espera de auditor ausente |
| c | breaker/rompedor + vocabulario de cierre/aprobación | clase c: rompedor con autoridad |
| d | auditor/its + listas/lecciones/deuda + completar, inventar, derivar, reconstruir y equivalentes | clase d: listas del auditor inventadas |

Controles positivos cubren reportar contexto, prohibir autoridad al rompedor, prohibir inventar listas, prohibir esperar al auditor ausente y exigir APROBADO con auditor activo. La ausencia expresada como “auditor is not available” no niega la aprobación atribuida al senior. Las referencias históricas a changes cerrados y rutas de archive no se interpretan como autorización actual.

`TestAuthorityFragmentMutations` conserva originales inmutables: cada subtest lee el fragmento limpio, exige cero infracciones, exige una única ubicación de reemplazo, muta una copia y exige el hit de **su clase**. Incluye M1, M2, M3, M5, M6 y cuatro variantes propias V1–V4. M4 sigue cubierto por `TestCompose_ReviewToggle_ReportAuthority`: se acotó la relación sujeto/verbo para impedir que el prefijo “With auditor active” haga pasar la inversión de roles.

Además ejecuté las diez mutaciones físicamente sobre assets, una por vez, con restauración byte a byte mediante `finally`. Cada ejecución arranca del original limpio; los comandos excluyen todos los goldens. M4 prueba inversión de roles, distinta de las cuatro clases de vocabulario. Los otros nueve fallan con el nombre de clase correspondiente. Las redacciones propias no se presentan como las cuatro futuras del rompedor.

## Punto 3 — atribución del corpus aditivo

EVIDENCIA: `TestDebtCorpusAdditiveCopies` exige `debtViolations(original)==0` antes de probar agregados, y exige que la cláusula del hit nuevo pertenezca al agregado. Separar por punto o punto y coma no impide atribuir el hit a la cláusula agregada. El original limpio impide que una infracción previa enmascare las 28 mutaciones. Agregar “The review does not block.” al asset falla explícitamente por `original no limpio`, antes de aceptar los agregados.

## Rojos y verdes

Cada bloque declara su SHA de ejecución. Los rojos iniciales parten de la base con working tree modificado; las mutaciones finales parten del primer commit de implementación con el ajuste de secuencia aún sin commitear. No se atribuyen al SHA los cambios del working tree.

### Rojo inicial de puntos 1 y 2

Comando: `go test ./internal/harness/config -run 'TestAuthorityVocabularyClasses|TestDebtContractsComposition|TestDebtCorpusAdditiveCopies' -count=1`.
SHA base: `c9dd12b29bc39a216d3b80adfd7fe2ab973199b6`. Salida parcial literal (las demás variantes repetían los mismos diagnósticos):

```text
2026-10-23T11:56:47-03:00
c9dd12b29bc39a216d3b80adfd7fe2ab973199b6
--- FAIL: TestAuthorityVocabularyClasses (0.00s)
    --- FAIL: TestAuthorityVocabularyClasses/a/Its_report_is_the_closing_verdict_whenever_no_auditor_is_configured. (0.00s)
        debt_contract_test.go:99: clase a sobrevivió: Its report is the closing verdict whenever no auditor is configured. ([])
--- FAIL: TestDebtContractsComposition (20.68s)
    --- FAIL: TestDebtContractsComposition/claude/01 (0.17s)
        debt_contract_test.go:165: clase b: espera de auditor ausente:  an open finding may remain only as debt explicitly justified by the auditor's verdict
        debt_contract_test.go:165: clase b: espera de auditor ausente:  only aprobado from the auditor, with conformity satisfied and complete reports, enables archive
        debt_contract_test.go:165: clase b: espera de auditor ausente:  a round begins with verify on a quiet diff and closes with the auditor verdict
        debt_contract_test.go:165: clase b: espera de auditor ausente:   blockers in round 1 or 2 return only the auditor's list to apply
    --- FAIL: TestDebtContractsComposition/claude/04 (0.14s)
        debt_contract_test.go:165: clase a: senior con autoridad:  the senior runs only in verify, once, to compare the ordinary verification findings against dated decisions in this change's artifacts and closed changes
    --- FAIL: TestDebtContractsComposition/claude/05 (0.17s)
        debt_contract_test.go:165: clase b: espera de auditor ausente:  an open finding may remain only as debt explicitly justified by the auditor's verdict
        debt_contract_test.go:165: clase b: espera de auditor ausente:  only aprobado from the auditor, with conformity satisfied and complete reports, enables archive
        debt_contract_test.go:165: clase b: espera de auditor ausente:  a round begins with verify on a quiet diff and closes with the auditor verdict
        debt_contract_test.go:165: clase b: espera de auditor ausente:   blockers in round 1 or 2 return only the auditor's list to apply
        debt_contract_test.go:165: clase a: senior con autoridad:  the senior reports historical context against dated decisions from this change and closed changes
    --- FAIL: TestDebtContractsComposition/claude/06 (0.14s)
        debt_contract_test.go:165: clase a: senior con autoridad:  the senior runs only in verify, once, to compare the ordinary verification findings against dated decisions in this change's artifacts and closed changes
    --- FAIL: TestDebtContractsComposition/claude/07 (0.19s)
        debt_contract_test.go:165: clase a: senior con autoridad:  when you launch the senior, add this to its prompt:  markdown ## what came before read the earlier phases' artifacts of this change under openspec/changes/{change-name}/, and list openspec/changes/archive/ for the changes that closed before it
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	26.079s
FAIL
2026-10-23T11:57:14-03:00
```

La primera implementación de la atribución del corpus rechazó indebidamente el caso con punto y coma al comparar con toda la frase; se corrigió la atribución a la cláusula. La primera ampliación del detector señaló contexto histórico y negaciones válidas; se acotaron esas excepciones y se agregó control de “not available”. Estas ejecuciones intermedias no se declaran verdes.

### Verde de controles, composición y corpus

Comando: `go test ./internal/harness/config -run 'TestAuthority|TestDebtContractsComposition|TestDebtCorpusAdditiveCopies|TestCompose_ReviewToggle_ReportAuthority' -count=1`.

```text
2026-10-23T11:59:44-03:00
c9dd12b29bc39a216d3b80adfd7fe2ab973199b6
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	23.597s
2026-10-23T12:00:08-03:00
```

### Rojos de M1–M6, V1–V4 y original contaminado

SHA al ejecutar: `b255ee44b32c27f1f01c1e8078108bc4f816c91a`, working tree con el ajuste final de secuencia.  Cada entrada registra clase, comando, hora y salida completa. Exit 1 es el resultado esperado del mutante; no son fallas del original. El script aborta si un mutante sobrevive y restaura el asset aun ante error.

```text

M1 clase a; original limpio restaurado antes de mutar
2026-10-23T12:02:14.493172-03:00
$ go test ./internal/harness/config -run TestDebtContractsComposition/claude/05|TestCompose_ReviewToggle_ReportAuthority -count=1
--- FAIL: TestDebtContractsComposition (0.15s)
    --- FAIL: TestDebtContractsComposition/claude/05 (0.15s)
        debt_contract_test.go:175: clase a: senior con autoridad:  its report is the closing verdict whenever no auditor is configured
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.164s
FAIL
exit 1; fin 2026-10-23T12:02:15.108078-03:00

M2 clase b; original limpio restaurado antes de mutar
2026-10-23T12:02:15.108671-03:00
$ go test ./internal/harness/config -run TestDebtContractsComposition/claude/05|TestCompose_ReviewToggle_ReportAuthority -count=1
--- FAIL: TestDebtContractsComposition (0.15s)
    --- FAIL: TestDebtContractsComposition/claude/05 (0.15s)
        debt_contract_test.go:175: clase b: espera de auditor ausente:  keep the round open until the auditor returns its verdict
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.160s
FAIL
exit 1; fin 2026-10-23T12:02:15.727262-03:00

M3 clase c; original limpio restaurado antes de mutar
2026-10-23T12:02:15.727733-03:00
$ go test ./internal/harness/config -run TestDebtContractsComposition/claude/05|TestCompose_ReviewToggle_ReportAuthority -count=1
--- FAIL: TestDebtContractsComposition (0.15s)
    --- FAIL: TestDebtContractsComposition/claude/05 (0.15s)
        debt_contract_test.go:175: clase c: rompedor con autoridad:  the breaker's report is enough to close the change
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.160s
FAIL
exit 1; fin 2026-10-23T12:02:16.338675-03:00

M4 clase autoridad invertida; original limpio restaurado antes de mutar
2026-10-23T12:02:16.340342-03:00
$ go test ./internal/harness/config -run TestDebtContractsComposition/claude/05|TestCompose_ReviewToggle_ReportAuthority -count=1
--- FAIL: TestCompose_ReviewToggle_ReportAuthority (0.00s)
    review_toggle_test.go:248: falta relación de autoridad: auditor\s+classifies[^;.]*decides
    review_toggle_test.go:248: falta relación de autoridad: orchestrator\s+executes
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.159s
FAIL
exit 1; fin 2026-10-23T12:02:16.956952-03:00

M5 clase d; original limpio restaurado antes de mutar
2026-10-23T12:02:16.957436-03:00
$ go test ./internal/harness/config -run TestDebtContractsComposition/claude/05|TestCompose_ReviewToggle_ReportAuthority -count=1
--- FAIL: TestDebtContractsComposition (0.15s)
    --- FAIL: TestDebtContractsComposition/claude/05 (0.15s)
        debt_contract_test.go:175: clase d: listas del auditor inventadas:  if the auditor is absent, mark auditor source absent and fill its lists from your own reading of the diff
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.159s
FAIL
exit 1; fin 2026-10-23T12:02:17.573853-03:00

M6 clase a; original limpio restaurado antes de mutar
2026-10-23T12:02:17.574268-03:00
$ go test ./internal/harness/config -run TestDebtContractsComposition/claude/05|TestCompose_ReviewToggle_ReportAuthority -count=1
--- FAIL: TestDebtContractsComposition (0.15s)
    --- FAIL: TestDebtContractsComposition/claude/05 (0.15s)
        debt_contract_test.go:175: clase a: senior con autoridad:  with auditor absent, closure proceeds on the senior's report alone
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.158s
FAIL
exit 1; fin 2026-10-23T12:02:18.194245-03:00

V1 clase a; original limpio restaurado antes de mutar
2026-10-23T12:02:18.194643-03:00
$ go test ./internal/harness/config -run TestDebtContractsComposition/claude/05|TestCompose_ReviewToggle_ReportAuthority -count=1
--- FAIL: TestDebtContractsComposition (0.15s)
    --- FAIL: TestDebtContractsComposition/claude/05 (0.15s)
        debt_contract_test.go:175: clase a: senior con autoridad:  the historian grants approval for archiving
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.159s
FAIL
exit 1; fin 2026-10-23T12:02:18.810392-03:00

V2 clase b; original limpio restaurado antes de mutar
2026-10-23T12:02:18.811991-03:00
$ go test ./internal/harness/config -run TestDebtContractsComposition/claude/05|TestCompose_ReviewToggle_ReportAuthority -count=1
--- FAIL: TestDebtContractsComposition (0.15s)
    --- FAIL: TestDebtContractsComposition/claude/05 (0.15s)
        debt_contract_test.go:175: clase b: espera de auditor ausente:  archiving awaits the absent auditor's approval
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.164s
FAIL
exit 1; fin 2026-10-23T12:02:19.437765-03:00

V3 clase c; original limpio restaurado antes de mutar
2026-10-23T12:02:19.438155-03:00
$ go test ./internal/harness/config -run TestDebtContractsComposition/claude/05|TestCompose_ReviewToggle_ReportAuthority -count=1
--- FAIL: TestDebtContractsComposition (0.15s)
    --- FAIL: TestDebtContractsComposition/claude/05 (0.15s)
        debt_contract_test.go:175: clase c: rompedor con autoridad:  approval for archive rests with the breaker
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.160s
FAIL
exit 1; fin 2026-10-23T12:02:20.047680-03:00

V4 clase d; original limpio restaurado antes de mutar
2026-10-23T12:02:20.048130-03:00
$ go test ./internal/harness/config -run TestDebtContractsComposition/claude/05|TestCompose_ReviewToggle_ReportAuthority -count=1
--- FAIL: TestDebtContractsComposition (0.15s)
    --- FAIL: TestDebtContractsComposition/claude/05 (0.14s)
        debt_contract_test.go:175: clase d: listas del auditor inventadas:  if the auditor is absent, mark auditor source absent and reconstruct auditor lists from available evidence
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.156s
FAIL
exit 1; fin 2026-10-23T12:02:20.659919-03:00

ORIGINAL clase original no limpio; original limpio restaurado antes de mutar
2026-10-23T12:02:20.660401-03:00
$ go test ./internal/harness/config -run TestDebtCorpusAdditiveCopies -count=1
--- FAIL: TestDebtCorpusAdditiveCopies (0.47s)
    debt_contract_test.go:222: original no limpio: [review sin autoridad: The review does not block]
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.471s
FAIL
exit 1; fin 2026-10-23T12:02:21.599738-03:00
```

### Actualización de expectativas completas

Sólo después de los rojos por aserción se regeneraron goldens:

```text
$ go test ./internal/harness/config -run '^TestCompose_AllToggles(_opencode)?$' -update -count=1
2026-10-23T12:00:18-03:00
c9dd12b29bc39a216d3b80adfd7fe2ab973199b6
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
2026-10-23T12:00:19-03:00
```

## JUICIO

Autor constructor, confianza media: los controles ejercitados sostienen los tres cambios solicitados y detectan las inversiones conocidas sin depender de goldens. No acreditan independencia ni aprobación de verify. El detector sigue siendo un guard léxico acotado, no un intérprete semántico de instrucciones arbitrarias.

## SUPUESTO

- Las cuatro redacciones futuras del rompedor son desconocidas. SUPUESTO: compartirán familias léxicas de las clases a–d; **no verifico ni afirmo** que ya hayan muerto. El criterio externo completo del punto 2 queda pendiente del rompedor y auditor de vuelta 2.
- SUPUESTO: preservar las instrucciones condicionales en el texto compuesto expresa correctamente la degradación al consumidor; no se ejercitó una sesión real de un agente instalado.
- Sin conflictos con cambios de terceros trackeados al inicio; sólo había logs ajenos no trackeados. Constructor de la misma familia Codex que la conformidad anterior; no constituye revisión independiente.

## Entrega y límites

Vuelta 2 de 3, apply corregido; verify pendiente. La decisión vigente sigue siendo NO APROBADO de vuelta 1 hasta nuevo dictamen. No se cierra ni archiva el change. `tasks.md` permanece local e ignorado por el repositorio; su actualización no se fuerza a git. El reporte versionado deja el contador y los resultados disponibles al siguiente rol.

## Regresión detectada por la suite completa y corrección

EVIDENCIA: la primera suite completa falló en dos tests existentes de secuencia, por cambiar las expresiones de orden completo y espera obligatoria. Se restauraron ambas expresiones bajo la condición `With auditor active`; se mantuvo separada la secuencia sin auditor. Se repitieron M1–M6, V1–V4 y el original contaminado (salidas anteriores de este informe corresponden a esa última repetición), se regeneraron los dos goldens y se commiteó la corrección. No se relajaron esos tests existentes.

```text
2026-10-23T12:00:44.967354-03:00
b255ee44b32c27f1f01c1e8078108bc4f816c91a
HOME temporal: /tmp/contratos-v2-home.gkk7ngal
$ HOME="$task_home" GOMODCACHE="$task_modcache" GOCACHE="$task_gocache" go test ./... -count=1
ok  	github.com/Group-Active-IA/active-stack/assets	0.050s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.567s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.301s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.683s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	0.005s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	0.006s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	0.012s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	0.006s
ok  	github.com/Group-Active-IA/active-stack/internal/backup	0.093s
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	0.140s
?   	github.com/Group-Active-IA/active-stack/internal/contracttest	[no test files]
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	0.016s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	0.020s
--- FAIL: TestVerifyContratoBucle (0.00s)
    verify_loop_test.go:19: falta espera obligatoria; no se permite: Launch the auditor in parallel with the senior.
    --- FAIL: TestVerifyContratoBucle/conformity_->_breaker_->_senior_(if_active)_->_auditor (0.00s)
        verify_loop_test.go:35: falta contrato: conformity -> breaker -> senior (if active) -> auditor
--- FAIL: TestAuditorDecideSinAutoridadParalela (0.02s)
    verify_loop_test.go:44: secuencia contradictoria: Launch the auditor in parallel with the senior.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	26.836s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	0.100s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	0.077s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	0.025s
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	0.011s
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.226s
ok  	github.com/Group-Active-IA/active-stack/internal/model	0.007s
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	0.012s
ok  	github.com/Group-Active-IA/active-stack/internal/planner	0.010s
ok  	github.com/Group-Active-IA/active-stack/internal/roles	0.106s
ok  	github.com/Group-Active-IA/active-stack/internal/system	0.026s
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.034s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	0.037s
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.034s
FAIL
Exit: 1
2026-10-23T12:01:13.052397-03:00
```

## Verde final — suite completa con HOME temporal

EVIDENCIA: caches de Go resueltas antes de cambiar HOME; HOME creado y eliminado mediante TemporaryDirectory. Sin build separado. Código y tests sin cambios después de este resultado; sólo documentación de entrega.

```text
2026-10-23T12:02:22.115798-03:00
32612f96a17e9db4da388e01819b3f7721caf36e
HOME temporal: /tmp/contratos-v2-home.cu14jj0l
$ HOME="$task_home" GOMODCACHE="$task_modcache" GOCACHE="$task_gocache" go test ./... -count=1
ok  	github.com/Group-Active-IA/active-stack/assets	0.051s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.570s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.245s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.709s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	0.006s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	0.010s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	0.008s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	0.018s
ok  	github.com/Group-Active-IA/active-stack/internal/backup	0.106s
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	0.176s
?   	github.com/Group-Active-IA/active-stack/internal/contracttest	[no test files]
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	0.012s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	0.017s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	26.719s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	0.105s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	0.081s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	0.038s
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	0.015s
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.275s
ok  	github.com/Group-Active-IA/active-stack/internal/model	0.013s
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	0.006s
ok  	github.com/Group-Active-IA/active-stack/internal/planner	0.006s
ok  	github.com/Group-Active-IA/active-stack/internal/roles	0.113s
ok  	github.com/Group-Active-IA/active-stack/internal/system	0.031s
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.042s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	0.037s
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.036s
Exit: 0
2026-10-23T12:02:50.087814-03:00
```
