# Active Stack: cómo era antes y cómo es ahora

Fecha del informe: 2026-10-23. Rango cubierto: `e9694b6` (ANTES, main antes de la fase 1, el stack que venía de Juan Cruz) hasta `de0ea76` (AHORA, HEAD de main). Fase 1: `e9694b6..0204022` (44 commits). Fase 2: `0204022..de0ea76` (51 commits). Total del rango completo: 95 commits.

## 1. Resumen en 10 líneas

1. Antes, `active-stack` ya era un instalador methodology-first con catálogo de harnesses, modos lite/full/custom y toggles del orquestador (`delegation`, `governance`, `evidence`, `auditor`, `roles`, `review`, `knowledge-base`, `senior`, `ship`, `engram`, `tdd`), pero varios de esos toggles eran más declaración que mecanismo.
2. `review` hacía "tres pases independientes" que en la práctica no eran independientes (compartían change directory y memoria Engram sin aislar claves), y la propia instrucción lo admitía como límite conocido.
3. `senior` corría en todas las fases y no tenía autoridad de cierre definida frente al auditor; `ship` no esperaba el CI completo del SHA pusheado; `archive` no tenía mecanismo real de reescritura de la base de conocimiento.
4. La fase 1 (9 commits `feat`/`fix` de código más 35 de proceso documental) convirtió esos toggles en contratos verificables: bucle constructor/rompedor/auditor de hasta tres vueltas en verify, aprobación humana escrita obligatoria para HIGH/CRITICAL, TDD con evidencia roja/verde pegada por tarea, archive que reescribe la KB, senior limitado a verify, ship que espera el CI.
5. La fase 1 se cerró con 8 changes archivados y una vuelta 3 de verify que el auditor declaró NO APROBADA; Juani decidió corregir afuera del bucle (opción B) y dio el OK escrito para archivar con deuda explícita.
6. La fase 2 construyó piezas que antes no existían como binario: `active-stack approval` (lectura de aprobación escrita, sin firma criptográfica, por decisión explícita de Juani), `active-stack reports capture` (captura de informes con escaneo de secretos y estado RETENIDO) y la instalación real de los 7 `role-*.md` de Juani como recursos del harness `role-definitions`.
7. También corrigió un bug heredado de permisos de OpenCode (reemplazaba todo el objeto `permission`, perdía reglas personales) por un reemplazo acotado a `permission.bash`, preservando claves hermanas.
8. El uso real (no sólo `go test`) encontró bugs que la suite en verde no había atrapado: el comando `approval` fallaba al llamarlo como lo llama `main`, un evento parcial con secreto podía escaparse del escaneo, `"task": ""` aparecía en el `opencode.json` real, y Windows no simula `chmod` para el test de permiso denegado.
9. Un ensayo aislado de reinstalación completa (sandbox Bubblewrap, HOME falso, sin tocar la máquina real de Juani) acreditó instalación headless, reinstalación idempotente, falla inducida con rollback parcial reportado como tal, y restore byte a byte verificado por SHA256; quedó pendiente la reinstalación en la máquina real de Juani, que es decisión suya.
10. Queda deuda explícita y escrita: el paso 5 (lector semántico de verify) y el paso 6 (documentación de la fase 2) no se construyeron todavía, N1/N1b de reintentos de CI quedó resuelto sólo como regla pero sin test de rerun real, y hay ítems D-11, D-20, D-21 abiertos sobre OpenCode y la captura automática por hooks.

## 2. Cómo era antes (e9694b6)

### Qué hacía

Active Stack ya era, en `e9694b6`, "un binario, un catálogo de harnesses, tu agente de IA listo para trabajar con método" (`README.md` en `e9694b6`). El catálogo (`internal/catalog/harnesses.yaml`) agrupaba harnesses en tres tipos: `skill` (SKILL.md clonado), `config` (texto bundleado en el binario) y `external` (binario/servicio de terceros, vía CLI o MCP).

### Cómo se instalaba

```
active-stack install                 # TUI interactiva (Bubbletea)
active-stack install --dry-run       # plan, sin ejecutar
active-stack install --mode full     # headless
active-stack install --agent claude  # agente concreto
active-stack install --help          # flags
```

Tres modos: **Lite** (sustrato mínimo: OpenSpec, Engram, Context7, `sdd-orchestrator`, `permissions`, `starter-add-command`), **Full** (Lite + fundación guiada: `active-orchestrator` y las skills de `kb-creator`, `project-scaffold`, `roadmap-generator`, `agent-instruction`, `skill-registry`, `web-scraper`, `discovery-research`, `find-skill`, `skill-creator`) y **Custom** (elección harness por harness, con `permissions` forzado siempre: "la seguridad no es opcional"). Agentes soportados: Claude Code, Codex, OpenCode (el README de `e9694b6` ya listaba estos tres en el badge de portada).

El harness clave, `sdd-orchestrator`, es de tipo `config` y se compone por **toggles** (`internal/harness/config/compose.go`) a partir de fragmentos en `internal/harness/config/assets/*.md`: `delegation`, `model-routing`, `governance`, `evidence`, `auditor`, `roles`, `review`, `knowledge-base`, `senior`, `ship`, `engram`, `tdd`. El resultado se inyecta en `CLAUDE.md`/`AGENTS.md` del usuario entre markers idempotentes (`internal/harness/config/inject.go`), con backup previo (`internal/backup`) y merge por markers (`internal/filemerge`).

### Cómo era el ciclo OPSX

El ciclo `explore → propose → apply → verify → archive` (más `ship` al publicar) ya existía con estos comportamientos (archivos en `internal/harness/config/assets/`, estado en `e9694b6`):

- **`review.md`**: proponía "tres pases independientes", cada uno con un lente (correctitud, cableado, comentarios/documentación). El propio texto admitía que la independencia era "pedida y parcialmente lograda": los pases podían leer el directorio del change y compartían la misma `topic_key` de Engram (`opsx/{change}/{phase}`), lo que permitía contaminación entre pases a pesar de la instrucción de no mirarse.
- **`senior.md`**: el senior se convocaba en **las cinco fases** (explore, propose, apply, verify, archive), con una tabla de "qué pregunta en cada fase". Etiquetaba hallazgos como BLOCKING o NOTE, y el texto decía "the senior labels, you hold": la decisión de frenar el ciclo quedaba en el orquestador humano/agente, sin que el auditor mediara.
- **`ship.md`**: ya tenía el diagnóstico de los seis estados del repositorio (commit directo a main, branch sin remoto, PR abierto, main adelantado, cambios sin commitear, sin remoto) y la regla de nunca mergear el PR, pero el cierre era "el ciclo termina en el PR abierto", sin exigencia de esperar su resultado de CI.
- **`governance.md`**: definía los cuatro niveles (CRITICAL/HIGH/MEDIUM/LOW) y "a HIGH o CRITICAL: describir el cambio y esperar confirmación", pero no existía el artefacto escrito obligatorio (`## Aprobación` en `proposal.md`) ni ningún comando que lo verificara.
- **`roles.md`**: ya tenía la tabla de los 7 roles con su "terminado" observable y la tabla de convocatoria por fase, pero no existía la obligación explícita de convocar 3+ especialistas antes de la firma en propose HIGH/CRITICAL.
- **`auditor.md`**: podía "quedarse callado" si todo estaba bien, sin un contrato de salida fijo (APROBADO/BLOQUEANTES) ni relación con el senior.
- **`knowledge-base.md`**: distinguía REWRITE de APPEND, pero no tenía mecanismo de "unión entre impacto declarado y diff real" ni etiqueta de sospecha para un archive que sólo hizo APPEND.
- **`strict-tdd.md`**: pedía una tabla resumen de evidencia (checkmarks) al final del reporte, sin exigir pegar comando/salida/hora bajo cada tarea en `tasks.md`.

El binario ya tenía `active-stack roles` (añadido por Juan Cruz en el commit `a697153`, anterior a `e9694b6`): lee el diff y devuelve qué revisores corresponden, sin convocar a nadie. **No existían** `active-stack approval` ni `active-stack reports capture`: no hay archivos `internal/approval/`, `internal/reports/`, `cmd/active-stack/approval_dispatch.go` en el árbol de `e9694b6` (confirmado con `git ls-tree -r e9694b6`).

### Qué no garantizaba

- Ninguna verificación de que la tabla de aprobación humana existiera o estuviera completa antes de construir en HIGH/CRITICAL.
- Ninguna captura determinística y versionada de los informes de los roles: dependía de que el orquestador transcribiera o de copiar a mano.
- Ningún escaneo de secretos antes de versionar un informe.
- `permission.bash` de OpenCode podía arrastrar residuo de un tier anterior al cambiar de tier (bug conocido, documentado más abajo).
- Los 7 `role-*.md` existían como texto en `roles.md` (la convocatoria), pero no como archivos instalados en el agente: convocar un rol no garantizaba que el sub-agente tuviera definición real.
- `senior` podía pedir frenar el ciclo directamente, sin pasar por un auditor que decidiera.
- `ship` podía declarar el trabajo publicado sin haber esperado el resultado del CI.

## 3. Qué cambió, por área

### 3.1 Apply: constructor con evidencia TDD real

**Antes:** `strict-tdd.md` pedía, al final del apply, una tabla resumen con checkmarks (RED ✅ Written, GREEN ✅ Passed, etc.) sin exigir el comando y la salida real pegados en `tasks.md`.

**Ahora:** cada tarea exige RED y GREEN pegados bajo la tarea en `tasks.md` (comando, `date +%T`, últimas líneas de salida); si un test nació verde, hay que escribir un mutante de 2 a 4 líneas que lo haga fallar en su aserción, pegar ese diff y revertirlo; el cierre corre las suites completas sobre el commit entregado y registra `git rev-parse HEAD`; el constructor declara consumidores fuera del área asignada y su familia de modelo (`internal/harness/config/assets/strict-tdd.md`, diff `e9694b6..de0ea76`; commit `f3fc02a` "feat(apply-constructor): exigir evidencia TDD y límites del constructor").

**Dónde se nota:** en `tasks.md` de cualquier change nuevo (hay bloques RED/GREEN reales, no una tabla de visto); en el propio repo, los informes de apply de fase 2 (`docs/propuestas/informes/2026-10-23-fase-2/pasos-1-4/apply/*.md`) siguen exactamente este formato con hora, SHA y salida pegada.

### 3.2 Verify: constructor/rompedor/auditor con bucle de hasta tres vueltas

**Antes:** `review.md` describía "tres pases independientes", cada uno con un lente distinto, sin auditor que decidiera entre ellos ni contador de vueltas persistido.

**Ahora:** un único rompedor cubre las tres lentes obligatorias (correctitud, cableado, verdad de comentarios/documentación) en un solo informe; el auditor lee ese informe (y el del senior, si está activo) y decide APROBADO o BLOQUEANTES, con DEUDA y LECCIONES; hay hasta tres vueltas, contadas en `## Bucle de verify` dentro de `tasks.md` del change (no en Engram); en la vuelta 3, con bloqueantes, se exige decisión humana escrita de Juani (corregir con nueva aprobación, aceptar deuda justificada, o frenar); nunca hay una cuarta vuelta automática (`internal/harness/config/assets/review.md`, diff `e9694b6..de0ea76`; commit `db3de54` "feat(verify-rompedor-auditor): definir bucle y decisión única del auditor").

**Dónde se nota:** la propia fase 1 vivió este bucle: vuelta 1 no aprobada (`8d05c4e`), vuelta 2 no aprobada (`248173f`), vuelta 3 cerrada por el auditor con bloqueantes (`18bb293`, informe `docs/propuestas/informes/2026-10-23/21-auditor-vuelta-3.md`) y decisión de Juani de corregir afuera del bucle (opción B, `e46294a`), con OK final para archivar (`5ab0f4f`).

### 3.3 Aprobación humana escrita y `active-stack approval`

**Antes:** `governance.md` pedía "a HIGH o CRITICAL, describir el plan y esperar confirmación", sin artefacto escrito obligatorio ni comando que lo verificara.

**Ahora:** `governance.md` exige la tabla `## Aprobación` en `proposal.md` (elección, persona, fecha, alcance) escrita por el humano, nunca completada por un agente (`internal/harness/config/assets/governance.md`, diff citado arriba; commit `9f2d03b` "feat(governance-aprobacion): exigir aprobación humana con alcance escrito"). El comando nuevo `active-stack approval --change <id> --commit <C>` (fase 2, commit `50de000`) verifica esa tabla sin firma criptográfica: por decisión explícita de Juani del 23/10 ("Juani decide no configurar firma criptográfica de commits", `docs/propuestas/informes/2026-10-23-fase-2/segunda-revision/06-decisiones-juani.md`), el estado máximo es `WRITTEN_UNVERIFIED` (código de salida 2), con el mensaje literal `"estructura válida; identidad humana no acreditada"` (`internal/approval/approval.go:73`, constante `writtenUnverifiedMessage`; usado en `cmd/active-stack/approval_dispatch.go:106` y probado en `cmd/active-stack/approval_dispatch_test.go:75`). Nunca se emite `VERIFIED`.

**Dónde se nota:** corriendo `active-stack approval --change <proposal.md> --commit <sha>` sobre cualquier propuesta ALTO/CRÍTICO; el código de salida (0 no existe salvo verificado criptográficamente, que no está implementado; 1 para MISSING/INVALID/STALE; 2 para WRITTEN_UNVERIFIED).

### 3.4 `active-stack roles` y la instalación real de los 7 roles

**Antes:** `active-stack roles` (de Juan Cruz, `a697153`) leía el diff y devolvía qué revisores correspondían, sin verificar que el `role-*.md` convocado existiera instalado en el agente.

**Ahora:** el mismo comando agrega un código de salida nuevo, `3`, cuando un rol **convocado** no tiene su `role-*.md` instalado o está vacío (`cmd/active-stack/roles_dispatch.go`, constante `exitDefinitionMissing`, diff `e9694b6..de0ea76`), consultando `internal/harness/roledefs.Health`. Los 7 roles (`role-analyst`, `role-architect`, `role-backend`, `role-dba`, `role-devops`, `role-frontend`, `role-qa`) ahora existen como recursos reales embebidos en `internal/harness/roledefs/assets/role-*.md` e instalados por el harness `role-definitions` (commits `2526eee` "instalar los siete role-*.md como sub-agentes reales" y `5ce2c0f` "los siete roles con el texto real de Juani (J4)"). Son copia literal de los agentes personales de Juani en `~/.claude/agents/`, bajo licencia MIT del repo, revisados antes de copiar para que no quedaran nombres, rutas, correos ni claves (decisión J4, `06-decisiones-juani.md`); quedan afuera `role-auditor`, `role-senior` y `role-implementer` porque no son de los siete de la lista canónica. La instalación respeta ownership: ausente crea, idéntico es no-op, propio intacto se actualiza con backup, personal/desconocido/modificado se preserva con conflicto y ruta alternativa, nunca se borra por glob ni sólo por nombre (`docs/propuestas/2026-10-23-fase-2-pasos-1-4.md`, sección "4, instalar-roles").

**Dónde se nota:** `ls internal/harness/roledefs/assets/` muestra los 7 archivos; `active-stack roles --files <diff>` devuelve código 3 si falta la definición de un rol convocado.

### 3.5 Captura de informes (`reports capture`): adapters y RETENIDO

**Antes:** no existía mecanismo de captura determinística; cualquier informe de rol dependía de que el orquestador lo transcribiera o de copiarlo a mano. No hay `internal/reports/` en `e9694b6`.

**Ahora:** `active-stack reports capture` (núcleo en commit `ac7b3dd` "feat(reports): núcleo de captura de informes antes del despacho") registra repo/worktree, change, fase, vuelta, rol, dispatch ID y SHA antes de aceptar un informe; antes de versionar cualquier original corre un escaneo de secretos (patrones propios: tokens, claves, correos, teléfonos), y si detecta algo el original queda local, ignorado por Git, con estado **RETENIDO** y causa declarada, sin volcar el valor detectado (`docs/propuestas/2026-10-23-fase-2-pasos-1-4.md`, "Contrato común y modo rápido"). Hay adapters confirmados para **Claude** (hooks `PostToolUse`/`SubagentHandback` + `SubagentStop`, versión 2.1.271) y **Codex** (`codex exec --json` 0.132.0, JSONL); **OpenCode no tiene contrato habilitado** y devuelve siempre `UNSUPPORTED` ("Claude usa su contrato de hooks. OpenCode no tiene contrato habilitado.", `internal/reports/provider.go:19`; commit `16e9a2b` "feat(reports): capturar eventos versionados de Claude y Codex" cubre sólo Claude/Codex).

**Dónde se nota:** `active-stack reports capture --change ... --report-text "..."` devuelve `COMPLETE` (código 0), `PENDING` (código 2), `RETAINED` (código 4) o `UNSUPPORTED` (código 3) según el caso; `README.md` lo documenta en la sección "Captura de informes por proveedor" agregada en esta fase.

### 3.6 Permisos de OpenCode: sólo `permission.bash`

**Antes:** un fix previo (`ddeb94f`, mencionado como alternativa rechazada en el propose) reemplazaba el objeto `permission` **entero** al cambiar de tier, lo que borraba cualquier clave hermana personal (por ejemplo `mcp_tool_notion_search`).

**Ahora:** `composeOpencodeOverlay` envuelve en el sentinel `__replace__` sólo `permission.bash` (la única clave cuyo valor es un objeto anidado que `MergeJSONObjects` fusionaba en lugar de reemplazar); las demás claves siguen siendo overwrites simples porque cada tier las declara todas (`internal/harness/config/permissions/overlays.go`, diff `e9694b6..de0ea76`; commit `aaadd10` "fix(opencode): reemplazar solo permission.bash, no el objeto permission entero"). El comentario en el propio código documenta la pérdida aceptada: una regla personal agregada a mano **dentro** de `permission.bash` no sobrevive un reinstall (confirmado también en el uso real, ver sección 5); eso es una pérdida deliberada aceptada por el propose, pendiente el gate de diff/aprobación explícita (deuda **D-11**, ver sección 6).

**Dónde se nota:** reinstalando sobre un `opencode.json` real con reglas personales: las claves hermanas de `permission` (`theme`, `mcp.notion`, `mcp_tool_notion_search`) sobreviven; una regla agregada a mano dentro de `bash` no sobrevive.

### 3.7 Archive que reescribe la base de conocimiento

**Antes:** `knowledge-base.md` distinguía REWRITE de APPEND como concepto, pero sin mecanismo de cobertura (comparar impacto declarado contra el diff real) ni señal de sospecha cuando el archive sólo apendeaba.

**Ahora:** archive cubre la unión entre `## Knowledge Base Impact` de `proposal.md` y el diff real; si sólo hubo APPEND sin REWRITE, el reporte se etiqueta `SUSPICIOUS: append-only`; las lecciones (`L-NNN`) y la deuda se toman de las listas LECCIONES/DEUDA del auditor, nunca inventadas; si el auditor está ausente se declara "auditor source absent" y no se inventan lecciones (`internal/harness/config/assets/knowledge-base.md`, diff citado; commit `84c570d` "feat(archive-base-completa): reescribir la base y conservar lecciones con evidencia").

**Dónde se nota:** `docs/metodologia-opsx.md`, creado en el archive de los 8 changes de fase 1 (commit `8d0b3b2` "docs(archive): cerrar ocho changes y actualizar la base metodológica"), es exactamente ese artefacto: tiene secciones de Decisiones (D-001 a D-003), Lecciones (L-001 a L-004) y una tabla de Deuda consolidada con 8 ítems, cada uno con su evidencia tipada.

### 3.8 Senior sólo en verify

**Antes:** `senior.md` se convocaba en las cinco fases del ciclo, con una tabla de qué preguntar en cada una; el senior podía marcar un hallazgo BLOCKING y frenar el ciclo directamente ("the senior labels, you hold").

**Ahora:** el senior corre **únicamente en verify**, después del rompedor y antes del auditor; ya no etiqueta BLOCKING/NOTE de forma final, sino que entrega contexto (decisiones fechadas, contradicciones) y el **auditor** decide, citando evidencia o una decisión fechada; el orquestador ejecuta el veredicto del auditor, no el del senior; el costo pasa de "cinco por change" a "como máximo tres por change, uno por vuelta de verify" (`internal/harness/config/assets/senior.md`, diff `e9694b6..de0ea76`; commit `66ad39e` "feat(senior-solo-en-verify): limitar contexto histórico a las vueltas de verify").

**Dónde se nota:** no hay informe de senior en explore/propose/apply/archive de ningún change nuevo; sólo aparece en verify (por ejemplo `docs/propuestas/informes/2026-10-23/07-senior-vuelta-2.md` o equivalentes de fase 2, `15-senior-vuelta-2.md`).

### 3.9 Ship que espera el CI completo

**Antes:** `ship.md` terminaba el ciclo en el PR abierto, sin exigir esperar su resultado de CI; "green checks" ni siquiera se mencionaba como condición de cierre explícita.

**Ahora:** "el ciclo termina en el PR abierto **y su resultado de CI**"; el orquestador debe descubrir las corridas (`gh run list --commit <sha>`), esperarlas (`gh run watch <run-id>`) e inspeccionar cada job (`gh run view <run-id> --json headSha,event,workflowName,jobs`); define explícitamente qué NO es verde: un job exitoso con otro pendiente, `queued`/`in_progress`/`cancelled`/`skipped`/`neutral`/`timed_out`, un resultado de otro SHA, cero corridas o una consulta fallida; para active-stack exige `ci.yml` en push con test en Ubuntu y Windows (`internal/harness/config/assets/ship.md`, diff citado; commit `60b5e83` "feat(ship-espera-ci): exigir CI completo sobre el SHA pusheado"). Además incorpora la política N1/N1b de Juani: un job cancelado o fallido por la plataforma, reintentado sobre **el mismo SHA** y en éxito, cuenta como verde; sobre otro SHA no cuenta (decisión del 23/10, `docs/propuestas/informes/2026-10-23-fase-2/segunda-revision/06-decisiones-juani.md`, motivada por un caso real: PropTech, 22/10, job `web` cancelado por instalación de dependencias y reintentado en verde sobre el mismo commit).

**Dónde se nota:** la sección "Complete CI on the pushed SHA" de `internal/harness/config/assets/ship.md`, visible en el `CLAUDE.md`/`AGENTS.md` instalado con el toggle `ship` activo.

### 3.10 Instalación sana: verify, idempotencia, rollback, restore byte a byte

**Antes:** no había verificación explícita de que lo instalado coincidiera con `Compose`, ni reporte de versión/SHA del binario.

**Ahora:** `active-stack version` imprime versión y SHA (vía metadatos Go); el verify de instalación valida identidad del binario `cmd/active-stack`, corre `roles --help`, compara el bloque instalado contra `Compose` para la variante y toggles efectivos, y rechaza markers duplicados (commit `84fcf0e` "feat(instalacion-sana): verificar binario y contenido instalado e informar versión"; fuentes `cmd/active-stack/version.go`, `internal/verify/{binary,delivery}.go`). Esto se ensayó de punta a punta en un sandbox aislado (ver sección 5): instalación headless Full/Claude, reinstalación idempotente (mismo SHA256 antes y después), falla inducida a mitad de instalación con rollback reportado como `partial` (nunca como éxito ficticio), y restore byte a byte verificado por SHA256 y `cmp` contra el `backup-primera-instalacion/`. **CONF-1 sigue abierto**: ese ensayo acredita el comportamiento aislado, no la instalación real en la máquina de Juani, que queda pendiente de su decisión (`docs/metodologia-opsx.md`, tabla de deuda consolidada, ítem CONF-1).

**Dónde se nota:** `active-stack version`; el verify-hook al final de cualquier `active-stack install`; el informe completo del ensayo en `docs/informes/2026-10-23-ensayo-reinstalacion.md`.

## 4. Tabla comparativa

| Área | ANTES (`e9694b6`) | AHORA (`de0ea76`) |
|---|---|---|
| Apply (constructor) | Tabla resumen de checkmarks al final del reporte | RED/GREEN pegados por tarea en `tasks.md`, mutante si nació verde, cierre sobre `git rev-parse HEAD` |
| Verify | "Tres pases independientes", sin auditor que decida ni contador persistido | Un rompedor con tres lentes obligatorias, auditor decide, bucle de hasta 3 vueltas contado en `tasks.md` |
| Aprobación HIGH/CRITICAL | "Esperar confirmación", sin artefacto verificable | Tabla `## Aprobación` escrita obligatoria + `active-stack approval` (WRITTEN_UNVERIFIED, sin firma criptográfica por decisión de Juani) |
| `active-stack roles` | Lee diff, dice qué revisores corresponden (ya existía, de Juan Cruz) | Igual, más código de salida 3 si falta la definición del rol convocado |
| Informes de roles | Sin captura determinística; dependía del orquestador | `active-stack reports capture`: escaneo de secretos, estado RETENIDO, adapters Claude/Codex confirmados, OpenCode UNSUPPORTED |
| Permisos OpenCode | Bug conocido: cambiar de tier podía arrastrar residuo de `bash`; alternativa rechazada reemplazaba todo `permission` y perdía reglas personales | `__replace__` acotado a `permission.bash`; claves hermanas personales sobreviven; regla personal dentro de `bash` no sobrevive (pérdida aceptada, D-11) |
| Roles instalados | Convocatoria sin archivo instalado garantizado | 7 `role-*.md` reales (texto de Juani) instalados por el harness `role-definitions`, con ownership y conflicto declarado |
| Archive / KB | REWRITE/APPEND como concepto, sin cobertura ni señal de sospecha | Cubre unión impacto/diff; `SUSPICIOUS: append-only` si falta REWRITE; lecciones/deuda tomadas del auditor |
| Senior | Convocado en las 5 fases, podía frenar el ciclo directamente | Sólo en verify, entre rompedor y auditor; el auditor decide, no el senior |
| Ship | Terminaba en el PR abierto, sin exigir CI verde | Espera CI completo del SHA pusheado, con N1/N1b para reintentos sobre el mismo SHA |
| Instalación sana | Sin verify de identidad/versión/Compose | `version`, verify de binario y Compose; ensayado con reinstalación, falla inducida y restore byte a byte en sandbox aislado |

## 5. Cómo se probó: hallazgos reales y la lección

El patrón se repite en casi todos los pasos de fase 2: **la suite de Go en verde no acreditaba que el binario funcionara en un entorno tan permisivo como el real**, porque los tests usaban directorios/temporales/`.gitignore` más favorables que el mundo real. Por instrucción explícita (surgida después del primer hallazgo), cada paso siguiente corrió el binario compilado contra un HOME falso o un repo git temporal nuevo, no sólo `go test`.

**1. El comando `approval` fallaba con la suite en verde.** Con los tests de `go test ./internal/approval ./cmd/active-stack -run Approval` en verde sobre el commit `6642a62`, la primera llamada real falló:
```
$ active-stack approval --change docs/propuestas/2026-10-23-fase-2-pasos-1-4.md --commit d26131e
error: approval.Verify: RepoDir, ProposalPath y Commit son obligatorios   (exit 1)
```
Causa: `main` llamaba a `runApprovalDispatch` con el directorio vacío y ningún test cubría ese camino. Arreglo en `53d4d4e`. Un segundo hallazgo real (deuda D-3 anticipada): `isDirty` usaba `git status -uall`, así que cualquier archivo `.log` sin trackear invalidaba una firma válida; arreglo en `1edd7f7` con `--untracked-files=no`. Fuente: `docs/propuestas/informes/2026-10-23-fase-2/pasos-1-4/apply/01-constructor-paso-1.md:192-215`. Lección del propio informe: "tests verdes no acreditan que el comando ande; correrlo una vez sí."

**2. Un evento parcial con secreto podía escapar al escaneo.** El test `TestProviderPartialSecretRetained` (RED, 2026-10-23T14:32:29) encontró que un evento parcial sensible quedaba en `PENDING` en vez de `RETAINED`:
```
provider_test.go:226: evento parcial sensible: status=PENDING err=<nil>
```
Causa: `Source.Raw` no se inicializaba antes de un retorno temprano por evento parcial, así que el escaneo de secretos nunca llegaba a correr sobre ese texto. Arreglo en `482ca98` ("fix(reports): retener secretos de eventos parciales de proveedores"), verde confirmado a las 14:32:51. Fuente: `docs/propuestas/informes/2026-10-23-fase-2/pasos-1-4/apply/05-constructor-pasos-2bcd.md:80-115`.

**3. "Filtración de retenidos": el test del constructor usaba un entorno más permisivo que el real.** El repo de prueba (`newRepo`) copiaba el `.gitignore` de active-stack, que ya ignora `informes/**/.retenidos/`; en un repo de usuario SIN ese `.gitignore`, un original con secreto (clave AWS + correo) quedaba RETENIDO pero el archivo original igual aparecía en `git add -A`:
```
$ git add -A && git diff --cached --name-only
openspec/changes/x/informes/verify/round-1/qa/.retenidos/97474ccae81b.txt   <- el original con la clave
```
Arreglo en `26f7ddc` ("fix(reports): la carpeta de retenidos se ignora sola en cualquier repo"): se escribe un `.retenidos/.gitignore` con `*` antes del original, para cualquier repo lo tenga o no. Fuente: `docs/propuestas/informes/2026-10-23-fase-2/pasos-1-4/apply/02-constructor-paso-2a.md:382-398`. Es la segunda vez en la misma tanda que el verde del constructor venía de un entorno más permisivo que el real (JUICIO del propio informe, confianza alta).

**4. `"task": ""` apareció en un `opencode.json` real.** La primera corrida real del binario (sin el fix) escribió una clave vacía en el archivo de permisos real del usuario de prueba. Causa: al reescribir `composeOpencodeOverlay` con `map[string]any` se perdió el `omitempty` que el struct original tenía para `Task` (ningún tier lo fija). Ningún test lo atrapó porque comparaban claves puntuales, no recorrían todas las claves buscando un string vacío. Arreglo: test nuevo `TestOpencodeNeverWritesEmptyPermissionVerb` y omitir la clave cuando está vacía. Fuente: `docs/propuestas/informes/2026-10-23-fase-2/pasos-1-4/apply/03-constructor-paso-3.md:262-276`; el comentario que documenta el hallazgo quedó también en el código, `internal/harness/config/permissions/overlays.go` (comentario sobre `TestOpencodeNeverWritesEmptyPermissionVerb`, "found the hard way").

**5. Windows no simula `chmod` para el test de permiso denegado.** El último commit de la fase, `de0ea76` ("test(reports): saltear permiso denegado en Windows, donde chmod no lo simula"), saltea en CI de Windows el caso de prueba que depende de `chmod` para simular "permiso denegado" al escribir un informe, porque el modelo de permisos de Windows no lo reproduce igual. Fuente: `.github/workflows/ci.yml` (matriz `ubuntu-latest, windows-latest` ya existía desde antes de la fase 1, confirmada sin diff entre `e9694b6` y `de0ea76`) y el commit mencionado.

**La lección repetida, en una frase:** un test que arma su propio entorno (repo temporal con `.gitignore` copiado, struct con `omitempty` implícito, directorio siempre pasado explícito) puede estar probando una versión del mundo más cómoda que la real. La fase 2 adoptó, paso a paso, la costumbre de compilar el binario y correrlo contra un HOME o repo real antes de declarar un paso terminado (ver `docs/propuestas/informes/2026-10-23-fase-2/pasos-1-4/apply/02-constructor-paso-2a.md:398`: "para los pasos que quedan, la corrida real va en un repo que NO sea active-stack").

El ensayo de reinstalación (sección 3.10) usó la misma lógica a escala de todo el instalador: sandbox Bubblewrap con HOME falso, `/opt` privado, `/dev` privado, sin sudo, sin tokens heredados, PATH acotado al bin falso. Resultado headless Full/Claude: instalación en verde, reinstalación idempotente (SHA256 idéntico), falla inducida (`read-only file system` en `commands/active/init.md`) con código de salida 1 y rollback reportado como **partial** por una limitación declarada (brew), restore byte a byte verificado por SHA256 (`25d2ecaf5a3f369e44c6e4aecb6d1f5fd5f646c5d97c378910d9d91955e6dfb7` antes y después) y verify final de 17 aprobados, 0 fallidos, 1 omitido (Context7, check no implementado). Fuente completa: `docs/informes/2026-10-23-ensayo-reinstalacion.md`.

## 6. Lo que todavía falta o es deuda

- **Paso 5, lector semántico de verify.** Juani aprobó la opción A (lector separado, bajo `review`, antes del senior y del auditor, dentro del mismo tope de tres vueltas) en la decisión J5 del 23/10 (`docs/propuestas/informes/2026-10-23-fase-2/segunda-revision/06-decisiones-juani.md`), pero el `propose` de pasos 1-4 lo deja explícitamente fuera de alcance ("quedan fuera... paso 5", `docs/propuestas/2026-10-23-fase-2-pasos-1-4.md:5`) y la tabla de orden de construcción lo marca como change candidato todavía no construido (`docs/propuestas/2026-10-23-metodologia-fase-2.md`, fila "5 | `verify-lectura-semantica`").
- **N1/N1b, reintentos de CI.** La regla está escrita y aprobada por Juani (un reintento exitoso sobre el mismo SHA cuenta como verde), pero el corpus de tests discriminantes para reruns reales (N1/N1b junto con N2b, N3, N4, N5, N8) sigue como deuda abierta heredada de fase 1 (`docs/metodologia-opsx.md`, tabla de deuda, fila "N1, N1b, N2b, N3, N4, N5, N8").
- **Paso 6, documentar fase 2.** Según la tabla de orden de construcción, requiere "verify completo según método vigente, deuda explícita" antes de tocar la documentación; no se construyó en este rango (`docs/propuestas/2026-10-23-metodologia-fase-2.md`, fila "6 | `documentar-fase2`").
- **D-11.** Reemplazar `permission.bash` descarta sin aviso cualquier regla personal agregada a mano dentro de ese objeto; no hay mecanismo de mostrar el diff ni pedir aprobación explícita antes de la pérdida (gate J6, TBD en el propose). Fuente: `docs/propuestas/informes/2026-10-23-fase-2/pasos-1-4/apply/03-constructor-paso-3.md:374-379`.
- **D-20.** El desacople entre el `verify-hook` del pipeline de instalación y `internal/verify.ChecksForHarness` no se investigó: puede significar que el health check de `role-definitions` nunca corre en una instalación real vía `active-stack install`, sólo cuando algo externo llama a `ChecksForHarness` directamente. Fuente: `docs/propuestas/informes/2026-10-23-fase-2/pasos-1-4/apply/04-constructor-paso-4.md:309-315`.
- **D-21.** La captura automática por hooks de Claude Code (paso 2b) no quedó cableada: la captura funciona invocando `reports capture` a mano, no automáticamente desde un hook instalado. Se mergeó así por decisión explícita de Juani ("terminate los adaptadores y mergeamos, demasiado por ahora"). Fuente: `docs/propuestas/informes/2026-10-23-fase-2/pasos-1-4/apply/05-constructor-pasos-2bcd.md:323-325`.
- **CONF-1, reinstalación real en la máquina de Juani.** Sigue abierta: el ensayo aislado (sección 5) acredita comportamiento de instalación/rollback/restore, pero no constituye instalación real; la decisión de fase 2 fue explícita en dejarla afuera de este alcance ("Reinstalar la fase 1 en la máquina de Juani ya" quedó como intención del modo rápido, pero el propose de pasos 1-4 la deja fuera: "quedan fuera... reinstalación en la máquina de Juani y cierre de CONF-1", `docs/propuestas/2026-10-23-fase-2-pasos-1-4.md:5`). Es decisión de Juani cuándo ejecutarla.
- **Deuda heredada de fase 1 sin pagar aún** (de `docs/metodologia-opsx.md`): CONF-3/RF2-3 (el senior exige auditor aun apagado, puede esperar un veredicto imposible), RF2-4/RF1-4 (restos de autoridad vieja del reviewer mezclados con la del auditor), RF1-6 (claves de Engram por fase en vez de por vuelta/rol en algunos fragmentos), RF3-2 (frases tipo "treated as green" todavía pasan el test de contratos de ship).
- Deuda de fase 2 sin resolver, además de D-11/D-20/D-21: **AR-1** (dirección de dependencia entre `approval` y `reports` sin fijar), **AR-2** (el enlace posterior a `Compose` no tiene propietario de módulo asignado), **AR-3** (costo de revertir persistencia/adapters no precisado), **BE-4** (falta fuente/esquema concreto por proveedor antes de habilitar cada adapter), **QA-1/QA-2** (regresión JSONC vs JSON estricto, y fallo entre adapters del mismo paso), todos listados en `docs/propuestas/2026-10-23-fase-2-pasos-1-4.md`, sección "Deuda aceptada (modo rápido)".

## 7. Apéndice: commits clave

**Fase 1** (`e9694b6..0204022`):
- `f3fc02a` feat(apply-constructor): exigir evidencia TDD y límites del constructor
- `db3de54` feat(verify-rompedor-auditor): definir bucle y decisión única del auditor
- `66ad39e` feat(senior-solo-en-verify): limitar contexto histórico a las vueltas de verify
- `84c570d` feat(archive-base-completa): reescribir la base y conservar lecciones con evidencia
- `9f2d03b` feat(governance-aprobacion): exigir aprobación humana con alcance escrito
- `531e1b6` feat(propose-multirol): exigir especialistas y auditor antes de la firma
- `60b5e83` feat(ship-espera-ci): exigir CI completo sobre el SHA pusheado
- `84fcf0e` feat(instalacion-sana): verificar binario y contenido instalado e informar versión
- `8c6a06b` fix(opencode-mcp-command-array): usar command como array para MCP locales
- `18bb293` decisión del auditor, vuelta 3 de verify: el bucle se cierra (no aprobado)
- `e46294a` decisión de Juani sobre RF3-1, opción B
- `ce7899e` fix(contracttest): rechazar el anuncio de archive sin puntuación (RF3-1)
- `5ab0f4f` OK de Juani para archivar la fase 1
- `8d0b3b2` docs(archive): cerrar ocho changes y actualizar la base metodológica

**Fase 2** (`0204022..de0ea76`):
- `50de000` feat(approval): lectura de aprobación escrita sin firma criptográfica
- `53d4d4e` / `1edd7f7` fix(approval): bugs reales encontrados en el primer uso
- `ac7b3dd` feat(reports): núcleo de captura de informes antes del despacho
- `26f7ddc` fix(reports): la carpeta de retenidos se ignora sola en cualquier repo
- `482ca98` fix(reports): retener secretos de eventos parciales de proveedores
- `16e9a2b` feat(reports): capturar eventos versionados de Claude y Codex
- `aaadd10` fix(opencode): reemplazar solo permission.bash, no el objeto permission entero
- `2526eee` feat(roles): instalar los siete role-*.md como sub-agentes reales
- `5ce2c0f` feat(roles): los siete roles con el texto real de Juani (J4)
- `a697153` (anterior a `e9694b6`, de Juan Cruz) feat(roles): subcomando que lee el diff
- `de0ea76` test(reports): saltear permiso denegado en Windows, donde chmod no lo simula
- `39701cb` docs(propuestas): deuda D-21, hooks automáticos de captura

## Nota de verificación

Todas las afirmaciones de este informe salen de lectura directa de archivos y `git show`/`git diff`/`git log` sobre el repo en `/home/juanisarmiento/Proyectos/active-stack-V2` (sin builds ni `go test`), más la lectura completa del ensayo de reinstalación en `docs/informes/2026-10-23-ensayo-reinstalacion.md`. No se verificó: el contenido de `.github/workflows/release.yml` ni del instalador Windows (`installer/windows/`), que no fueron pedidos; el estado actual de D-17/D-18/D-19 (capability de `AgentsDir` por proveedor, archivo de conflicto `*.active-stack.md`) más allá de lo citado en los informes de apply, que no se leyó con detalle adicional; y no se corrió ningún comando del binario en este informe (sólo lectura), así que los resultados de comandos citados en la sección 5 son los que los propios informes de apply registraron, no una repetición propia.
