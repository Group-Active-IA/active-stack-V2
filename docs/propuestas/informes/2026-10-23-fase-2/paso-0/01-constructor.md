# Constructor — paso 0 contratos-fase2-deuda

Fecha: 2026-10-23. Rol: CONSTRUCTOR; vuelta 1 de 3. Sin veredicto de verify.

## Alcance y autorización

Rama feat/metodologia-fase-2; autorización escrita de Juani en 7f1ef88ecd64f83e78c51ed00162df3ff4dd84b6. Leí AGENTS.md, propuesta fase 2 (incluida aprobación), Compose y bloque compuesto materializado en golden Claude (igualdad verificada por TestCompose_AllToggles), fragmentos de apply/strict-tdd, review, senior, verify ordinario, KB, base Codex/Claude, goldens y metodología con lecciones L-001…L-004. Referencias PropTech L-019/L-023/L-029/L-030/L-031/L-034 leídas en strict-tdd; no confundo su numeración con la local.

## Implementación y consumidores

- RF1-4: sección vieja de review deja explícito que rompedor informa, auditor clasifica/decide y orquestador ejecuta; deuda justificada y decisión excepcional siguen el bucle aprobado. Eliminadas frases incompatibles que entregaban autoridad paralela a la persona en cualquier vuelta.
- CONF-3/RF2-3: Compose elige senior-without-auditor.md con review activo/auditor ausente. Senior informa después del rompedor y se conserva el informe para decisión escrita de Juani. Sin review usa el fragmento ordinario ya existente, que conserva senior y decisión escrita. No hay espera de auditor inexistente.
- RF2-4: mark en minúscula en la oración de KB; el test previo normaliza mayúsculas y mantiene condición, en vez de fijar Mark.
- Goldens Claude/OpenCode regenerados secuencialmente y diff revisado: sólo review y mark; no cambios incidentales de variantes ni catálogo.

Consumidores externos al área: instaladores que llaman Compose, health check que reconstruye Compose, agentes que obedecen el bloque y comandos verify/ship. Sin cambios de contratos de API, adapters, persistencia o paths. No se modifica backup/filemerge/pipeline, catálogo, .gitignore, homes personales, binarios, .env ni permisos.

## Tests de regla y cobertura

TestDebtContractsComposition recorre 8 variantes × 16 combinaciones review/auditor/senior/engram, con KB activo. Comprueba contradicciones por cláusula y vocabulario, y contexto propio del senior cuando falta auditor. No acepta que el fallback de otro fragmento o una frase correcta escondan la contradicción.

TestDebtVocabularyOwnVariants usa sujetos, redacción, idioma y mayúsculas distintos, controles positivos de negación e historia y Mark inicial válido. TestDebtCorpusAdditiveCopies conserva el bloque correcto y agrega contradicciones a Compose y a los comandos verify/ship. Cubre N2b, N3, N4, N5, N8 y RF3-2 con variantes propias; cancelado considerado verde se rechaza, sin decidir qué rerun satisface una política concreta.

JUICIO: son guardas deterministas acotadas, no evaluación semántica general ni corpus ciego. No afirmo que deuda del auditor esté pagada. Falsos positivos observados durante desarrollo: senior mencionado en una cláusula cuyo sujeto decisor era auditor; corregí limitando el patrón al sujeto senior y variantes adverbiales. Limitaciones: pronombres, negaciones complejas, orden distante de sujeto/verbo, sinónimos fuera del vocabulario y exenciones de ejemplo histórico pueden producir falsos negativos; no se midieron tasas sobre un conjunto independiente. N1/N1b exitoso tras cancelación no se califica automáticamente sin política aplicable.

## Conflictos resueltos y desviaciones

Resolución propia: el test antiguo exigía “does not block” y la narrativa vieja mientras el bucle ya asignaba autoridad al auditor. Actualicé ese test a relaciones de autoridad, en lugar de conservar la contradicción para hacerlo pasar. El bloque strict-tdd ordena no build; CI incluye build, por lo que se omite por instrucción explícita del operador y no se afirma CI completo.

Desviación de orden documental: creé el change con CLI al inicio, pero completé proposal/design/tasks después de ejecutar el primer rojo y cambios de texto. No presento esos artefactos como previos a implementación; el propose firmado fue la fuente de alcance y el rojo sí precedió a producción. `openspec instructions apply` detectó los artefactos faltantes; se completaron después, sin fabricar aprobación. La trazabilidad roja/verde se recupera de salidas reales guardadas en /tmp/contratos-fase2-evidencia.

## SUPUESTO

Modelo/proveedor: familia OpenAI/Codex según contexto de sesión; identidad como SUPUESTO porque no hay registro externo de despacho. No se demuestra independencia frente a familia del orquestador. Sin subagentes ni revisores convocados. Los casos propios son del constructor: no son convergencia independiente ni evaluación ciega; los IDs de deuda/casos conocidos vienen de informes históricos 14 y 20, no de un auditor de esta vuelta. No se consultó memoria Engram; no se afirma disponibilidad ni recuperación. No se instala ni se reclama instalación sana.

## Pregunta para Juani — fuera de construcción de regla nueva

N1/N1b (TBD): ¿qué política permite considerar satisfecho un job cancelado cuando una corrida posterior exitosa del mismo SHA lo reemplaza, y cómo debe acreditarse el conjunto requerido? No implementé rechazo universal de reruns ni adopté S; requieren sus decisiones propias. Esta pregunta no sustituye una aprobación ni cierra esa deuda.

## EVIDENCIA ejecutada

Baseline antes de modificar producción:

```text
SHA base: 7f1ef88ecd64f83e78c51ed00162df3ff4dd84b6 (working tree modificado)
Hora: 2026-10-23T11:36:59-03:00
$ go test ./internal/harness/config ./assets -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.797s
ok  	github.com/Group-Active-IA/active-stack/assets	0.020s
```

A continuación, copia de tasks.md con rojo/verde y mutante revertido:

# Tasks — contratos-fase2-deuda

Autorización: Juani, 7f1ef88; exclusivamente paso 0.

## 1. Deuda textual

- [x] 1.1 Corregir autoridad, senior sin auditor y Mark; verificar TestDebtContractsComposition.

RED anterior a cambios de producción; incluye las tres deudas en aserciones.

```text
SHA base: 7f1ef88ecd64f83e78c51ed00162df3ff4dd84b6 (working tree modificado)
Hora: 2026-10-23T11:37:44-03:00
$ go test ./internal/harness/config -run TestDebt -count=1
    --- FAIL: TestDebtContractsComposition/generic/10 (0.06s)
        debt_contract_test.go:91: Mark incidental:  If the auditor is absent, Mark auditor source absent and do not invent its lists
    --- FAIL: TestDebtContractsComposition/generic/11 (0.07s)
        debt_contract_test.go:91: review sin autoridad: The review **does not block**
        debt_contract_test.go:91: autoridad personal: it costs to fix now, and that call belongs to the person, not to the reviewer
        debt_contract_test.go:91: Mark incidental:  If the auditor is absent, Mark auditor source absent and do not invent its lists
    --- FAIL: TestDebtContractsComposition/generic/12 (0.05s)
        debt_contract_test.go:91: Mark incidental:  If the auditor is absent, Mark auditor source absent and do not invent its lists
    --- FAIL: TestDebtContractsComposition/generic/13 (0.08s)
        debt_contract_test.go:91: review sin autoridad: The review **does not block**
        debt_contract_test.go:91: autoridad personal: it costs to fix now, and that call belongs to the person, not to the reviewer
        debt_contract_test.go:91: Mark incidental:  If the auditor is absent, Mark auditor source absent and do not invent its lists
        debt_contract_test.go:109: senior sin auditor: falta vocabulario (?i)auditor.*absent
        debt_contract_test.go:109: senior sin auditor: falta vocabulario (?i)Juani.*written decision
        debt_contract_test.go:113: espera imposible del auditor ausente
    --- FAIL: TestDebtContractsComposition/generic/14 (0.06s)
        debt_contract_test.go:91: Mark incidental:  If the auditor is absent, Mark auditor source absent and do not invent its lists
    --- FAIL: TestDebtContractsComposition/generic/15 (0.08s)
        debt_contract_test.go:91: review sin autoridad: The review **does not block**
        debt_contract_test.go:91: autoridad personal: it costs to fix now, and that call belongs to the person, not to the reviewer
        debt_contract_test.go:91: Mark incidental:  If the auditor is absent, Mark auditor source absent and do not invent its lists
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	7.959s
FAIL
```

GREEN

```text
SHA base: 7f1ef88ecd64f83e78c51ed00162df3ff4dd84b6 (working tree modificado)
Hora: 2026-10-23T11:38:37-03:00
$ go test ./internal/harness/config -run TestDebt -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	7.547s
```

## 2. Corpus y goldens

- [x] 2.1 Probar vocabulario con variantes propias y mutación aditiva, revisar diff de goldens; verificar tests TestDebt y TestCompose_AllToggles.

Corpus nació verde: backup declarado en /tmp/contratos-fase2-evidencia/debt_contract_test.go.original. Mutante temporal revertido antes de verde.

```diff
diff --git a/tmp/contratos-fase2-evidencia/debt_contract_test.go.original b/internal/harness/config/debt_contract_test.go
index 7af31c3..ccfd7de 100644
--- a/tmp/contratos-fase2-evidencia/debt_contract_test.go.original
+++ b/internal/harness/config/debt_contract_test.go
@@ -47,7 +47,9 @@ func debtViolations(body string) []string {
 			}
 		}
 	}
-	return hits
+	// Mutante temporal: silenciar todas las contradicciones.
+ hits = nil
+ return hits
 }
 
 func TestDebtVocabularyOwnVariants(t *testing.T) {

```
RED del mutante (falla de aserciones, no compilación)

```text
SHA base: 7f1ef88ecd64f83e78c51ed00162df3ff4dd84b6 (working tree modificado)
Hora: 2026-10-23T11:40:36-03:00
$ go test ./internal/harness/config -run 'TestDebt(Vocabulary|Corpus)' -count=1
        debt_contract_test.go:164: contradicción aditiva sobrevivió: Cancelled checks are considered green
    --- FAIL: TestDebtCorpusAdditiveCopies/ship/command/In_round_1_the_person_chooses_blockers_to_fix (0.07s)
        debt_contract_test.go:164: contradicción aditiva sobrevivió: In round 1 the person chooses blockers to fix
    --- FAIL: TestDebtCorpusAdditiveCopies/ship/command/The_senior_alone_approves_archive (0.08s)
        debt_contract_test.go:164: contradicción aditiva sobrevivió: The senior alone approves archive
    --- FAIL: TestDebtCorpusAdditiveCopies/ship/command/Never_rewrite_old_text;_append_below_it (0.08s)
        debt_contract_test.go:164: contradicción aditiva sobrevivió: Never rewrite old text; append below it
    --- FAIL: TestDebtCorpusAdditiveCopies/ship/command/The_agent_may_fill_the_approval_table (0.07s)
        debt_contract_test.go:164: contradicción aditiva sobrevivió: The agent may fill the approval table
    --- FAIL: TestDebtCorpusAdditiveCopies/ship/command/A_summary_can_substitute_for_the_report (0.07s)
        debt_contract_test.go:164: contradicción aditiva sobrevivió: A summary can substitute for the report
    --- FAIL: TestDebtCorpusAdditiveCopies/ship/command/With_zero_runs,_treat_delivery_as_green (0.08s)
        debt_contract_test.go:164: contradicción aditiva sobrevivió: With zero runs, treat delivery as green
    --- FAIL: TestDebtCorpusAdditiveCopies/ship/command/Cancelled_checks_are_considered_green (0.07s)
        debt_contract_test.go:164: contradicción aditiva sobrevivió: Cancelled checks are considered green
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	5.510s
FAIL
```

GREEN tras restaurar backup

```text
SHA base: 7f1ef88ecd64f83e78c51ed00162df3ff4dd84b6 (working tree modificado)
Hora: 2026-10-23T11:43:10.383387-03:00
$ go test ./internal/harness/config -run TestDebt -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	25.926s
```

RED de compatibilidad del test viejo y goldens (expectativa histórica contradictoria)

```text
SHA base: 7f1ef88ecd64f83e78c51ed00162df3ff4dd84b6 (working tree modificado)
Hora: 2026-10-23T11:43:10.383667-03:00
$ go test ./internal/harness/config ./assets -count=1
--- FAIL: TestCompose_AllToggles (0.00s)
    compose_test.go:59: Compose output does not match golden.
        got len=114455 want len=114471
--- FAIL: TestCompose_AllToggles_opencode (0.00s)
    compose_test.go:83: Compose output does not match golden.
        got len=114797 want len=114813
--- FAIL: TestCompose_ReviewToggle_ReportsDoesNotBlock (0.00s)
    review_toggle_test.go:246: the reports-not-blocks section is missing "does not block" — the rule itself
    review_toggle_test.go:246: the reports-not-blocks section is missing "a defect from a taste" — the distinction whose absence turns a reviewer into an obstacle
    review_toggle_test.go:246: the reports-not-blocks section is missing "switched off in a week" — the cost of blocking on taste, named so it is not abstract
    review_toggle_test.go:246: the reports-not-blocks section is missing "closing over an open finding" — the legitimate case, stated so nobody treats a finding as a veto
    review_toggle_test.go:246: the reports-not-blocks section is missing "invisible" — the illegitimate case: the finding that nobody can see was never made
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	11.428s
ok  	github.com/Group-Active-IA/active-stack/assets	0.023s
FAIL
```

GREEN después de corregir expectativa y regenerar goldens secuencialmente

```text
SHA base: 7f1ef88ecd64f83e78c51ed00162df3ff4dd84b6 (working tree modificado)
Hora: 2026-10-23T11:43:10.383787-03:00
$ go test ./internal/harness/config ./assets -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	10.910s
ok  	github.com/Group-Active-IA/active-stack/assets	0.020s
```

## Bucle de verify

Vuelta: 1 de 3. Estado: pendiente de inicio de verify, apply del constructor entregado; no se incrementa contador.
SHA: se registra al entregar el commit. Informes: docs/propuestas/informes/2026-10-23-fase-2/paso-0/01-constructor.md; copia íntegra local informes/apply-r1-constructor.md.
Conformidad, rompedor, senior y auditor: FALTANTES, no convocados por el constructor.
Bloqueantes: por determinar en verify. Decisión: PENDIENTE, sin APROBADO ni archive.

## Workflow follow-up

- Verify por orquestador y roles correspondientes, con informes completos y auditor; no hay cierre acreditado.
- N1/N1b: política rerun (TBD), pregunta a Juani; no se inventa prohibición universal.
- CONF-1 sigue abierto; RF1-6 y S no se implementan en este paso.

## Entrega y comprobación sobre commit de implementación

Commit: 8469d942ff2ad0c29bd93bdf79bcea1513541085 (`fix(metodologia): corregir contratos y deuda del paso cero`).

```text
Hora: 2026-10-23T11:43:35-03:00
HOME temporal: /tmp/contratos-fase2-home.QSgmFU
$ HOME="$task_home" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" go test ./... -count=1
(En la ejecución, caches resueltas antes de cambiar HOME.)
ok  	github.com/Group-Active-IA/active-stack/assets	0.045s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.806s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.915s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.971s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	0.007s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	0.011s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	0.016s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	0.008s
ok  	github.com/Group-Active-IA/active-stack/internal/backup	0.138s
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	0.219s
?   	github.com/Group-Active-IA/active-stack/internal/contracttest	[no test files]
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	0.022s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	0.040s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	28.719s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	0.118s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	0.111s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	0.052s
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	0.026s
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.291s
ok  	github.com/Group-Active-IA/active-stack/internal/model	0.007s
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	0.004s
ok  	github.com/Group-Active-IA/active-stack/internal/planner	0.010s
ok  	github.com/Group-Active-IA/active-stack/internal/roles	0.148s
ok  	github.com/Group-Active-IA/active-stack/internal/system	0.026s
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.040s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	0.050s
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.053s

Exit: 0
$ HOME="$task_home" GOMODCACHE="$task_modcache" GOCACHE="$task_gocache" go vet ./...
Exit: 0
$ git diff --check
Sin salida, exit 0
$ gofmt -l internal/harness/config
Sin salida, exit 0
$ openspec validate contratos-fase2-deuda
Change contratos-fase2-deuda is valid (skip_specs documental explícito)
```

La entrega documental tiene commit separado. Se repite suite tras ese último commit; SHA/hora/salida quedan bajo Cierre final en tasks.md y en el informe íntegro local del change. No se ejecuta build (prohibido), push, instalación, archive ni CI remoto. No se acredita Windows desde este entorno Linux.

Estado: constructor entregado; verify vuelta 1 de 3 pendiente. No hay aprobación del auditor. La deuda sigue abierta hasta evaluación independiente.
