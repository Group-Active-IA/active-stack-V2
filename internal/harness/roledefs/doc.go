// Package roledefs implements the `role-definitions` harness (propose
// docs/propuestas/2026-10-23-fase-2-pasos-1-4.md, sección 4 — "instalar-roles").
//
// QUÉ INSTALA. Los siete role-*.md que active-orchestrator necesita para
// convocar sub-agentes reales por `subagent_type` (internal/roles define los
// siete IDs canónicos: role-dba, role-backend, role-frontend, role-devops,
// role-qa, role-architect, role-analyst). Sin estos archivos en el directorio
// de agentes del usuario, "convocar al rol" es una instrucción sin lector: el
// Agent tool no tiene de dónde sacar la definición.
//
// DÓNDE SE INSTALAN, Y POR QUÉ NO ES LO QUE EL PROPOSE DICE LITERAL. El
// propose pide "directorio propio bajo SkillsDir resuelto por capability del
// adapter" — el mecanismo de internal/harness/skill, que coloca un
// SKILL.md + recursos en <SkillsDir>/<harnessID>/. Ese mecanismo es correcto
// para una skill que un agente DESCUBRE y LEE, pero los siete roles no son
// una skill: son subagentes que Claude Code resuelve desde archivos PLANOS en
// ~/.claude/agents/<nombre>.md, uno por subagente, sin subdirectorio (así
// es como el propio orquestador de esta sesión invoca role-dba, role-qa,
// etc. — ver el listado de subagentes disponibles). Meter los siete adentro
// de ~/.claude/skills/role-definitions/ los deja legibles pero NO los
// registra como subagente invocable: el mecanismo que el propose nombra no
// produce el efecto que el propose busca.
//
// La instrucción operativa de esta tarea de apply (no el propose) fija el
// escenario real a probar explícitamente sobre ~/.claude/agents/, con
// agentes personales preexistentes ahí — eso es la autoridad que se sigue
// acá. Se documenta como SUPUESTO/deuda en el informe de apply, no se
// resuelve en silencio.
//
// POR QUÉ NO HAY UN ÍNDICE SKILL.md. Consecuencia directa de lo anterior:
// ~/.claude/agents/ es un directorio que Claude Code recorre buscando
// definiciones de subagente; soltar un SKILL.md ahí (que no es una) es
// terreno no probado y no es necesario para el mecanismo. El "índice" que
// el propose pide se resuelve en cambio como el manifiesto de ownership de
// este paquete (manifest.go) — JSON, no Markdown, y fuera de
// ~/.claude/agents/ — que de todas formas es lo único que necesita un
// consumidor programático (health check, CLI) para saber qué instaló este
// harness y con qué hash. Deuda nueva, no silenciada.
//
// AR-2 (propose, "Deuda aceptada"): "el enlace posterior a Compose no tiene
// módulo propietario ni dependencias declaradas; asignar propietario y
// declarar que config no depende de install/verify". Este paquete ES ese
// propietario: es el "enlace puro posterior a Compose" — el paso que, con
// contexto explícito (homeDir, adapter, backupDir), materializa en disco lo
// que la sección "## role specialists" de internal/harness/config/assets/
// roles.md declara en prosa (los siete roles por nombre). config.Compose
// sigue siendo puro (texto, sin I/O); roledefs es quien enlaza ese texto con
// el filesystem, DESPUÉS de que Compose devolvió el string, nunca adentro de
// Compose. La dirección de dependencia es roledefs → config (lee los IDs
// canónicos de internal/roles, no de config); config NUNCA importa
// internal/install, internal/verify ni este paquete — ver
// internal/harness/config/dependency_direction_test.go, que lo verifica
// escaneando los imports reales, no sólo lo declara en un comentario.
//
// QA-2 (propose, "Deuda aceptada"): "falta el caso de fallo entre adapters
// del mismo paso: el primero escribe, el segundo falla; exigir recuperación
// del primero y conservación de ediciones posteriores". Install, en este
// paquete, recupera por sí mismo las escrituras parciales de TODOS los
// adapters ya procesados en la misma llamada antes de devolver el error —
// no depende del rollback genérico del pipeline, que sólo repite Rollback()
// sobre pasos que ya terminaron en éxito (internal/pipeline/rollback.go:
// ExecuteRollback itera result.Steps con Status == StepStatusSucceeded; un
// paso que falló nunca entra a esa lista). Ver TestInstall_SecondAdapterFails_
// FirstAdapterWriteIsUndone en install_test.go.
package roledefs
