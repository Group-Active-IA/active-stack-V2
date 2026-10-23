# Auditor, vuelta 1 de verify (Claude opus, textual del registro)

No apruebo la vuelta 1: hay cuatro bloqueantes, y además falta el informe de conformidad, así que el auditor no tenía con qué aprobar.

**Chequeos sobre el informe 10, el del rompedor**
1. Evidencia: tiene horas (10:17:48 a 10:18:08) pero no pega ninguna salida, y el comando dice `go test -run <test del contrato>`. JUICIO mío, confianza alta: igual se sostiene. Leí `verify_loop_test.go:9-60` y `ship_ci_test.go:16` y sólo usan `strings.Contains`, así que un texto contradictorio que se agrega no los rompe. M1 a M7 los verifiqué leyendo los tests, no los corrí. M8 no lo verifiqué.
2. Supuestos: declara dos, y son honestos.
3. Cobertura: dictamina sobre `4748680..299889d` entero, pero no abrió `strict-tdd.md` (+104 líneas, paso 3), `internal/harness/external/mcp.go` (paso 1), `version.go`, `hook.go` ni `harnesses.yaml`.
4. Acuerdo: no aplica. Es el único revisor de la vuelta.
5. Lo que nadie miró: `strict-tdd.md` y `mcp.go`. Que nadie los mencione no quiere decir que estén bien.

**Chequeos sobre los informes 08 y 09, los del constructor**
1. Evidencia: en el commit `299889d`, el 08 pasó de 751 líneas a 13 y el 09 perdió unas 1145. Los dos ahora son resúmenes y su link "Informe con salidas" apunta al mismo archivo. Las salidas reales sobreviven en `openspec/changes/*/informes/apply-entrega-1-constructor.md`, con comando, SHA y hora. Los tests se corrieron sobre los commits de código; después de `531e1b6` sólo se tocaron docs (`git diff --stat 531e1b6..HEAD`, 5 archivos de informes), así que esas salidas siguen describiendo el código actual.
2. Supuesto ascendido: lo que el constructor llama "14 + 7 mutantes detectados" son frases borradas, no inversiones de sentido. Todas las salidas dicen `falta contrato: <frase>`, por ejemplo `ship-espera-ci/informes/apply-entrega-1-constructor.md:212`, y el mutante de wiring cambia la frase por `Wiring: omitted.` (`constructor-tanda-b-completo.md:206`).
3. Cobertura: cada informe dice qué tocó en cada paso, y se puede contrastar.
4. No aplica: es un solo autor.
5. No se reinstaló nada en tu máquina (§3.6 punto 4). Queda declarado como "excluida por Juani" en `instalacion-sana/tasks.md:241`; no lo verifiqué.

**Cómo clasifico cada hallazgo**
- RF1-1, **BLOQUEANTE**. §3.3 y §3.4 exigen que esos mutantes den rojo, y hoy pasan los 9.
- RF1-2, **BLOQUEANTE**. Leí `knowledge-base.md:191`: la frase vale siempre, sin condición, y contradice la fila "auditor con KB activo" de §3.2.
- RF1-3, **BLOQUEANTE**. "For this repository" está en `ship.md:207` y en `opsx/ship.md:395`. Como el bloque se instala global, en cualquier otro proyecto exige un CI que no existe.
- RF1-5a, **BLOQUEANTE**. `ship.md:179` sigue citando "the senior labels, you hold", la frase que el propose firmado en §3 elimina. Vuelve a meter la autoridad doble que §3.1 prohíbe. RF1-5b (`senior.md:17`) queda como NOTA.
- RF1-4 es NOTA: `review.md:106` sí dice "the auditor decides", sólo quedan restos en las líneas 68-73. RF1-6 también es NOTA, porque §3 no pidió tocar `engram-protocol`. RF1-7 se suma a RF1-1 como M8. Todas las notas van a deuda.

**Proceso: lo que impide aprobar aunque el código estuviera bien**
No hay informe de conformidad ni de senior (§3.5 pide cuatro por vuelta), y §3.1 exige conformidad satisfecha para dar APROBADO. Además, de los 9 changes, en 4 `## Bucle de verify` dice "No iniciado" y los otros 5 ni tienen la sección (§3.3). Esto lo arregla el orquestador, no el constructor.

**Lo que corrige el constructor en la vuelta 2**
1. M1, M1b, M2, M3, M4, M5, M6 y M7 aplicados de a uno: cada test nombrado da FAIL por su propia aserción, con la frase contradictoria en el mensaje. Se registra con comando, SHA y hora en el informe del change.
2. M8 (`!= 1` cambiado por `< 1` en `delivery.go`): un test de `./internal/verify` falla con un bloque duplicado.
3. RF1-2: la línea del archivista dice que sólo aplica si el auditor está ausente. Un test compone el bloque con auditor y KB activos y falla si la frase aparece sin esa condición.
4. RF1-3: en las dos copias se nombra active-stack explícitamente. Un test de cada copia falla si vuelve "For this repository".
5. RF1-5a: se borra la cita de `ship.md:179`. La prohibición de esa frase pasa a chequear el bloque compuesto entero, no sólo `senior.md`.

DECISIÓN: NO APROBADO (vuelta 1 de 3)
