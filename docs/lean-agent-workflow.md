# Flujo compacto de agentes

`lean-workflow` es un toggle predeterminado de `sdd-orchestrator`. Cambia las instrucciones generadas; no ejecuta por sí mismo agentes, plugins ni pruebas. Las definiciones instaladas son contratos de trabajo, no una garantía de obediencia del modelo.

## Fases y roles

Los comandos OPSX conservan sus nombres; el orquestador asigna los roles dentro de cada fase. Tiene todos los recursos (roles, fases, plugins configurados, memoria del proyecto) y decide cómo usarlos dentro de estas reglas.

| Fase | Roles |
|---|---|
| `propose` | planificador, después revisor del plan |
| `apply` | constructor |
| `verify` | rompedor, después auditor |
| `archive` | cierre con APROBADO del auditor y actualización del conocimiento |

## Secuencia y responsabilidad

| Rol | Recibe | Entrega |
|---|---|---|
| Planificador (`role-planner`) | Objetivo, límites y decisiones pertinentes | Proposal, design, tasks y specs aplicables |
| Revisor del plan (`role-plan-reviewer`) | Documentos y fuentes necesarias | Contradicciones, escenarios faltantes, viabilidad y dictamen; correcciones al planificador |
| Constructor (`role-constructor`) | Plan revisado, archivos y consumidores nombrados | Código con baseline y RED/GREEN ejecutados, pruebas y limitaciones |
| Rompedor (`role-breaker`) | Revisión/diff quieto y consumidores | Casos adversos de lógica, wiring y documentación; fallos reproducidos vuelven al constructor |
| Auditor (`role-auditor`) | Revisión actual y originales completos de los informes requeridos | Veredicto sobre conformidad, evidencia, cobertura y pendientes; deriva cada defecto al responsable |

El agente principal conserva la visión del proyecto y coordina los pasos. No añade automáticamente siete especialistas ni un senior. Puede incluir criterios especializados en el encargo del rol correspondiente. Las correcciones repiten las verificaciones afectadas; no se aprovecha un informe verde de otra revisión para cerrar. Los bloqueantes o informes faltantes impiden cierre. Tras tres intentos fallidos se pide dirección humana.

## Clasificación del change

La propuesta registra LOW, MEDIUM, HIGH o CRITICAL con riesgo, alcance y justificación. LOW suele ser un cambio local de bajo impacto; MEDIUM cruza lógica o módulos; HIGH afecta ampliamente datos/configuración; CRITICAL alcanza autenticación, pagos, seguridad o límites de datos sensibles. El tamaño pequeño no rebaja un riesgo crítico. Los niveles ajustan profundidad y rigor, nunca el TDD de código. Desde MEDIUM se despachan cinco roles separados; en LOW el orquestador puede planificar y revisar el plan él mismo, dejándolo registrado en `tasks.md`, y constructor, rompedor y auditor siguen siendo agentes separados. Con governance activo, HIGH/CRITICAL requieren aprobación humana explícita antes de construcción: el informe del revisor no es esa aprobación.

## Profundidad por nivel

Los mínimos nunca se saltean; el orquestador gasta esfuerzo donde está el riesgo.

| Nivel | Plan | TDD | Verificación |
|---|---|---|---|
| LOW | el orquestador puede fusionar plan y revisión | RED y GREEN por comportamiento | una ronda del rompedor, después auditor |
| MEDIUM | cinco roles | más bordes y errores | rondas hasta quedar limpio, después auditor |
| HIGH | cinco roles y aprobación humana | más regresiones de consumidores y mutante en aserciones clave | rondas hasta quedar limpio; el auditor ejecuta chequeos propios |
| CRITICAL | cinco roles y aprobación humana | mutante en cada aserción de seguridad o datos | como HIGH; con proveedor de identidad externo, recorrido humano real antes de archivar |

Un fallo reproducido siempre abre otra ronda. Al tercer intento fallido el orquestador ofrece tres opciones: alcance nuevo aprobado, deuda justificada o detener el change; sin respuesta escrita nada avanza. Retomar la sesión no reinicia ni incrementa el contador.

## Contexto y estado

Cada encargo contiene objetivo, archivos/diff concretos, decisiones relevantes, consumidores externos al área, restricciones, criterios verificables y ruta de salida. El orquestador guarda un resumen operativo y consulta los originales cuando hacen falta; no reenvía toda la conversación o todo el repositorio a cada rol. El auditor recibe los originales completos necesarios para evaluar la evidencia. Constructor, rompedor y auditor leen las lecciones numeradas (L-NNN) si existen; con un registro de skills, el encargo pega sólo sus reglas pertinentes.

Cada informe original se guarda en `openspec/changes/<change>/informes/<fase>-<intento>-<rol>.md` y se indexa en `LEEME.md` con autor, ejecutor/proveedor declarado, familia compartida, revisión e informes faltantes. Una identidad sin registro de despacho es SUPOSICIÓN. Una ronda de corrección arregla sólo la lista del auditor, agrupada por archivo; lo nuevo que no bloquea pasa a DEUDA.

El estado durable del change registra rol activo, revisión o diff, intento, rutas de informes, bloqueantes y próxima acción. Se consulta OpenSpec para el estado de artefactos. Las notas y el resumen no sustituyen esa consulta. Una interrupción no cuenta como un nuevo intento ni convierte ausencia de evidencia en aprobación.

## TDD para código

Para cada comportamiento nuevo o modificado: ejecutar baseline, escribir prueba, observar fallo de aserción por el motivo esperado, implementar el mínimo, observar GREEN, agregar casos diferentes y refactorizar con pruebas verdes. Una corrección empieza por una reproducción que falla. Un test nacido verde no acredita RED; un error de compilación tampoco demuestra la aserción. Registrar comandos, salida y hora bajo la tarea.

Documentación sin comportamiento ejecutable usa validación de formato/contrato y declara la excepción; no fabrica un ciclo RED/GREEN. Ningún toggle histórico desactiva el contrato TDD de código en lean.

## Proveedores y recursos

Los roles no fijan proveedor o modelo. Cada proyecto declara quién ejecuta cada rol en una tabla `## Role executors` de su CLAUDE.md o AGENTS.md:

| Role | Executor |
|---|---|
| breaker | codex plugin |
| constructor | <plugin o proveedor/modelo> |

Un rol sin fila corre como subagente nativo (`role-*` en Claude) o inline donde no hay subagentes. Si el ejecutor declarado no está disponible no se sustituye en silencio: se informa y se pregunta o se usa el predeterminado registrando la sustitución. La tabla no lleva credenciales, sólo dónde viven. Active Stack no agrega conectores API ni promete trasladar permisos entre procesos.

Como el ejecutor puede ser Codex u otro proveedor que no hereda instrucciones, cada encargo es autosuficiente: contrato del rol, formato de informe tipado con su ruta y el contexto acotado, nunca la conversación entera. En Claude el subagente nativo carga su propio procedimiento y el orquestador no lo repite. Los siete especialistas son opcionales: sus descripciones dicen que el orquestador los consulta si el cambio toca su área, sin convocatoria automática, y carecen de modelo fijo.

## Instalación sobre instalaciones previas

- Un archivo `role-*.md` idéntico a una versión ya publicada por el stack se trata como propio: backup y actualización.
- Un archivo desconocido se preserva; la versión del stack queda en `<archivo>.active-stack.md` y el conflicto se informa. Si es uno de los cinco roles de flujo, se marca BLOQUEANTE y el chequeo de salud lo reporta como `foreign`.
- Con el plugin `engram@engram` habilitado en Claude no se registra otra vez el MCP `engram`.
- El instalador nunca borra archivos que no pueda probar propios.

## Precedencia de toggles

| Selección | Resultado con `lean-workflow` activo |
|---|---|
| `delegation`, `roles`, `review`, `auditor`, `tdd` presentes o ausentes | La secuencia de cinco roles, contexto acotado y TDD de código siguen siendo obligatorios |
| `model-routing` | Ignorado: la ejecución sale de `## Role executors` de cada proyecto |
| `senior` | Contexto histórico a cargo del auditor; no agrega una etapa |
| `governance` | Agrega contrato compacto de clasificación/aprobación humana |
| `evidence` | Agrega formato compacto de reportes tipados |
| `knowledge-base` | Agrega continuidad y actualización del destino documental del proyecto |
| `engram` | Agrega memoria persistente si está disponible |
| `ship` | Agrega controles de publicación y CI |
| Quitar `lean-workflow` | Selecciona la composición legacy completa con sus reglas históricas; no mezcla ambos ciclos |

Los nombres desconocidos siguen siendo errores; el orden o repetición de toggles no cambia el resultado. Lite y Full seleccionan lean por el catálogo; una configuración explícita antigua sin el toggle conserva legacy. Las pruebas de composición cubren ambos modos.

## Conocimiento y publicación

La fundación sigue generando KB y roadmap. Cada proyecto declara su destino, incluido `obsidian.md` si lo usa. Leer el contexto relevante al iniciar y actualizar hechos, decisiones y próximos pasos al cerrar, conservando historia válida y excluyendo secretos. No crear una nota duplicada ni copiarla completa a cada encargo. Engram no sustituye OpenSpec ni la documentación canónica.

El cierre documental no es una orden de publicar. Con ship activo, exigir revisión de estado, aprobación correspondiente y CI completo sobre el SHA publicado; no merge automático. Este cambio del repositorio no reinstala dotfiles, no borra especialistas personales y no promete ahorro monetario. Las comparaciones de bytes/líneas describen tamaño generado, no tokens efectivos, caché ni facturación.
