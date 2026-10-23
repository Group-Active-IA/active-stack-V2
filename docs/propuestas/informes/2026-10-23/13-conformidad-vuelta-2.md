EVIDENCIA: revisión estática de `4748680..a79c333`, rama `feat/metodologia-proptech`, contra §3 a §5 y los nueve `tasks.md`. Excluí docs e informes como implementación. No ejecuté builds ni tests; las salidas citadas son registros existentes.

| Paso | Completitud y evidencia | JUICIO sobre correctitud y coherencia |
|---|---|---|
| 0 | Cuatro tareas tildadas, rojo/verde y cierre registrados; faltan ensayos de §3.6. | Identidad del binario, `version` y comparación con Compose implementadas; CONF-1. |
| 1 | Dos tareas tildadas con rojo/verde y cierre. | `command` array y `environment` implementados en `internal/harness/external/mcp.go:271`; sin desvío observado. |
| 2 | Contratos en ambas copias, rojo/verde y goldens registrados. | SHA, conjunto requerido y estados incompletos contemplados; cobertura de mutantes parcial, CONF-4. |
| 3 | Contrato del constructor y goldens con rojo/verde. | Evidencia por tarea, límites, mutantes y familia declarada presentes en `internal/harness/config/assets/strict-tdd.md:67`; sin desvío observado. |
| 4 | Bucle, autoridad, tres lentes y mutantes registrados. | Secuencia principal conforme; degradación incompleta al componer, CONF-2 y CONF-3. |
| 5 | Contratos y goldens registrados. | Senior limitado a verify y máximo tres; contradicción con auditor apagado, CONF-3. |
| 6 | Contratos y goldens registrados. | Unión impacto/diff, lecciones, deuda y evidencia presentes; ocho variantes con archive a sonnet. |
| 7 | Contratos, corrección de assertions y cierres registrados. | Tabla humana, alcance y L-029 presentes en `internal/harness/config/assets/governance.md:23`; sin desvío observado. |
| 7b | Contratos y goldens registrados. | N especialistas más auditor, secuencia e informes íntegros presentes en `internal/harness/config/assets/roles.md:192`; sin desvío observado. |

- **CONF-1 · ADVERTENCIA · JUICIO, completitud:** `openspec/changes/instalacion-sana/tasks.md:239` declara ausencia de reinstalación; no acredita snapshot/manifiesto, reinstalación repetida, falla inducida, restore ni aceptación por destino exigidos por §3.6. Los cuatro checks tildados no cubren todo el paso firmado.
  **SUPUESTO:** la exclusión atribuida a Juani en `tasks.md:241` modifica el alcance; este recorte no aporta la decisión primaria que lo respalde.

- **CONF-2 · CRÍTICO · JUICIO, correctitud/coherencia:** `internal/harness/config/assets/review.md:128` contiene la regla para operar sin `review`, pero ese fragmento desaparece al apagarlo. El comando Claude conserva la alternativa; los bloques de los demás agentes no reciben esa regla. Se incumple §3.2, especialmente sin review ni auditor.
  **EVIDENCIA:** `internal/harness/config/verify_loop_test.go:85` exceptúa precisamente esa ausencia del requisito de degradación visible; sólo comprueba que desaparezca el bucle.

- **CONF-3 · ADVERTENCIA · JUICIO, coherencia:** `internal/harness/config/assets/senior.md:110` exige esperar siempre el veredicto del auditor, y `senior.md:124` le asigna siempre la decisión. Con senior activo y auditor apagado contradice la salida por decisión humana escrita de `review.md:123`. Condicionar ambas instrucciones a auditor activo.

- **CONF-4 · ADVERTENCIA · JUICIO, completitud de pruebas:** `internal/harness/config/ship_ci_test.go:15` y `assets/ship_ci_test.go:21` fijan frases para los estados negativos, pero sólo rechazan explícitamente una contradicción sobre cancelación. No se acredita la batería completa de mutantes exigida por §3.4 para pendiente, omitido, SHA ajeno, cero runs y consulta fallida en ambas copias.

EVIDENCIA: `git diff --check 4748680..HEAD` no informó errores. JUICIO: no observé implementación de los pasos 8 a 12 fuera del alcance firmado.