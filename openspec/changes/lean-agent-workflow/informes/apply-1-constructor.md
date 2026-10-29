# Apply 1, constructor: grupo 5 (installer meets previous installs)

- Autor: constructor. Ejecutor: subagente Claude general-purpose (Claude Code). Nivel: HIGH, aprobado por Juani (tasks.md §5).
- Revisión: diff de trabajo sobre `61887a4` (rama `feat/context-efficiency`), sin commit. El árbol ya tenía cambios sin commit de otros roles; no se revirtió nada de eso.
- Restricciones cumplidas: sin commit, push, build ni instalación real. `~/.claude`, `~/.claude.json` y `~/.active-stack` sólo se **leyeron** (hashes y estructura JSON) para fundamentar 5.2 y 5.3; los tests usan `t.TempDir()`.

## Fuentes leídas

AGENTS.md §3/§4; tasks.md §5; informes/LEEME.md; `internal/harness/roledefs/{install,manifest,health,resources,types,uninstall,doc}.go` y tests; `internal/install/{steps,plan,types,self_install}.go`; `internal/pipeline/{result,runner,stages}.go`; `cmd/active-stack/headless/executor.go`, `cmd/active-stack/headless/windows.go` (mapeo de eventos), `cmd/active-stack/main.go`, `roles_dispatch.go`, `verify_binary.go`; `internal/tui/model.go`; `internal/verify/{harness_checks,hook,report}.go`; `internal/harness/external/{install,mcp,types}.go`; `internal/agents/claude/adapter.go`; `internal/catalog/harnesses.yaml` (entrada engram); `internal/uninstall/{plan,steps}.go`; `internal/backup/restore.go`; `internal/harness/skill/clone.go`.

Casos dictados por este encargo: los cinco roles de flujo, `engram@engram`, `settings.json → enabledPlugins` y `installed_plugins.json`, commit 5ce2c0f, preferencia "no romper la instalación del resto pero imposible no verlo".

Contaminación: el encargo trae la hipótesis de que los siete especialistas coinciden con 5ce2c0f; la verifiqué con hashes en vez de darla por cierta. No leí informes de otros revisores salvo LEEME.md (para el formato del índice).

## EVIDENCIA

Todas las horas son -0300 del 2026-10-29. Comando base: `GOCACHE=/tmp/claude-1000/gocache go test -count=1`.

| Hora | Comando | Salida relevante |
|---|---|---|
| 12:06:55 | `go test ./...` (baseline) | exit 0, 29 paquetes `ok`, sin fallos previos |
| (lectura) | `git log --all --format=%h -- internal/harness/roledefs/assets` | `5ce2c0f`, `2526eee` (7 assets cada uno) |
| (lectura) | `sha256sum ~/.claude/agents/role-*.md` | los 7 especialistas = hashes de 5ce2c0f; `role-auditor.md` 70517bf… (desconocido); también hay `role-implementer.md` y `role-senior.md` |
| (lectura) | `ls ~/.active-stack/state/role-definitions/` | no existe: sin manifiesto, el install de hoy da 7 conflictos falsos |
| (lectura) | python sobre `~/.claude/settings.json`, `installed_plugins.json`, `~/.claude.json` | `enabledPlugins: {engram@engram: true, codex@openai-codex: true}`; `plugins.engram@engram` lista con 1 entrada; `mcpServers` ya tiene `engram` (el duplicado existe hoy) |
| 12:08:38 | `./internal/harness/roledefs/ -run 'Shipped\|NonShipped\|SymlinkToShipped'` | RED: `role-dba.md outcome = "conflict", want "updated"`; `5ce2c0f hash … is not in the embedded shipped list`; `shippedRoleHashes is out of date with git history` |
| 12:09:01 | `./internal/harness/roledefs/` | GREEN `ok` |
| 12:09:11 | mutante: `isShippedVersion` ignora el nombre / devuelve siempre true | `de-otro-rol: outcome = "updated", want "conflict"`; y los 5 subtests de borde fallan; revertido |
| 12:09:16 | mutante: `os.Lstat`→`os.Stat` en `installOneResource` | `SymlinkToShipped…: outcome = "updated", want "conflict"`; revertido |
| 12:09:34 | `-run Health` | RED: `got 0 issues, want exactly 1`; `IsWorkflowRole("role-planner.md") = false` |
| 12:09:45 | `./internal/harness/roledefs/` | GREEN `ok` |
| 12:10:09 | `./internal/verify/ -run RoleDefinitions` | RED: `a foreign workflow role must not fail verify (it would roll back the whole install): … Status:failed … role-auditor.md (foreign)` |
| 12:10:24 | `./internal/verify/` | 2 regresiones de consumidores: `got 2 checks, want 1`, `got 2 results, want 1` |
| 12:10:31 | `./internal/verify/ ./cmd/active-stack/...` | `ok` después de actualizar esos 2 tests (ahora también exigen que el check duro nombre `role-dba.md (missing)`) |
| 12:11:07 | `./internal/install/ -run RoleDefinitions_` | RED: `want exactly one degraded event carrying the conflicts, got []` (los eventos muestran `role-definitions:role-definitions … succeeded`) |
| 12:11:27 | `./internal/install/` | GREEN `ok` |
| 12:11:36 | mutantes en `roleConflictReport`: `IsWorkflowRole` → `true` / → `false` | falla `OnlySpecialistConflict…` (`1 BLOQUEA(N)`) y `ConflictsAreReported…` (sin BLOQUEANTE); revertidos, `ok` |
| 12:11:50 | `./internal/tui/ -run SucceededDoesNotErase` | RED: `row after degraded+succeeded = {… status:succeeded err:<nil>}` |
| 12:11:56 | `./internal/tui/ ./cmd/active-stack/...` | GREEN `ok` |
| 12:12:39 | `./internal/agents/claude/ -run PluginMCP` | RED: `InstalledPluginsPath = ""`, `MCPPluginID(engram) = "", want engram@engram` |
| 12:13:11 | `./internal/harness/external/ -run 'Plugin\|PropagatesNotes'` | RED: `~/.claude.json must not be written when the plugin provides engram`; `~/.claude.json was modified`; `settings-corrupto: Notes = [], wantNote=true`; `PropagatesNotes: … Notes:[]` |
| 12:13:38 | `./internal/harness/external/ ./internal/agents/...` | GREEN `ok` |
| 12:13:44 | mutante: no exigir `installed_plugins.json` con la entrada | fallan `habilitado-pero-no-instalado` y `habilitado-instalado-vacio`; revertido |
| 12:14:03 | `./internal/install/ -run ReportsInstallerNotes` | RED: `want one skipped event with the plugin note, got []` |
| 12:14:11 | `./internal/install/` | GREEN `ok` |
| 12:14:30 | `./cmd/active-stack/headless/ -run PrintsSkippedNotes` | RED: `note must appear inline and in the end summary (seen 0 times)` |
| 12:14:38 | `./cmd/active-stack/...` | GREEN `ok` |
| 12:14:51 | `./internal/tui/ -run SkippedNote` | RED: `completion screen must keep success and show the note` |
| 12:15:00 | `./internal/tui/` | GREEN `ok` |
| 12:15:07–12:15:38 | `go test ./...`; `gofmt -l .`; `go vet ./...`; `git diff --check` | exit 0 (29 `ok`, `contracttest` sin tests); gofmt vacío; `vet-ok`; `diffcheck-ok`. Archivos nuevos sin espacios finales (`rg "[ \t]+$"` rc=1) |
| 12:15:53 | `./internal/harness/roledefs/` tras editar doc.go; `-v -run 'MatchGitHistory\|Anchor'` | `ok`; ambos PASS |

## Qué se hizo, por tarea

### 5.1 Conflictos visibles y health de roles de flujo

- `roleDefinitionsStep.Run()` ya no descarta el `Result`. Con conflictos emite **un** `ProgressEvent` `StepStatusDegraded` (el canal que ya usaban `skillStep` best-effort y `selfInstallStep`; sin `onProgress`, cae a stderr como ellos). El mensaje lista cada conflicto con agente, ruta preservada y ruta del `<archivo>.active-stack.md`; los roles de flujo van primero y marcados `BLOQUEANTE del flujo lean`, los especialistas como `conflicto (especialista, no bloquea)`. Sin conflictos no se emite nada (instalación limpia y reinstalación idempotente, testeadas).
- `roledefs.IsWorkflowRole` + `workflowRoleIDs` (la lista canónica se arma desde ahí). `Health` agrega el problema `foreign`: rol de flujo presente, no vacío y con bytes distintos de los del stack. Los especialistas no se comparan (uno personal es legítimo).
- `verify.checkRoleDefinitions` pasa a dos checks: presencia/vacío (duro, igual que antes) y contenido de roles de flujo (`Soft`, mensaje con `BLOQUEANTE del flujo lean` y dónde está la versión del stack).
- TUI: `applyProgressEvent` ya no deja que el `Succeeded` que el runner manda después pise una fila `Degraded`. Era un bug previo: el aviso desaparecía y la pantalla final decía "All steps succeeded.". También afectaba a los skills best-effort de C-32.

### 5.2 Versiones publicadas = del stack

- `internal/harness/roledefs/shipped.go`: `shippedRoleHashes` con los 14 sha256 de los 7 especialistas en 2526eee y 5ce2c0f, embebidos en el binario. `isShippedVersion` sólo vale para el mismo nombre de recurso.
- `installOneResource`: tras el chequeo por manifiesto, un archivo regular cuyo hash está en la lista va por `installOwnIntactUpdate` (backup + relectura anti-concurrencia + escritura + undo), y el manifiesto queda con el hash nuevo.
- `TestShippedRoleHashes_MatchGitHistory` recalcula desde git y exige igualdad exacta. Se saltea sin git o fuera de un checkout. `TestShippedRoleHashes_Anchor5ce2c0f` fija los 7 hashes reales sin depender de git.

### 5.3 Engram por plugin

- Interfaz opcional `external.PluginMCPAdapter` (`SettingsPath`, `InstalledPluginsPath`, `MCPPluginID`), implementada por `claude.Adapter`. `MCPPluginID("engram") = "engram@engram"`, cualquier otro `""`. Las rutas salen del adapter.
- `registerStdioMCP` (el que registra el MCP de los harness external homebrew/npm, incluido engram) consulta `pluginProvidesMCP` antes de hacer backup o escribir. Saltea sólo si se dan las dos condiciones: `enabledPlugins["engram@engram"] == true` y una lista no vacía en `installed_plugins.json`. En ese caso no toca el archivo y deja una `Result.Notes`. Si ya había un `mcpServers.engram` (el caso real de esta máquina), la nota lo nombra para limpieza manual y **no lo borra**.
- Fallback (registra igual que hoy): archivos ausentes, sin `enabledPlugins`, plugin deshabilitado, valor no booleano, habilitado pero no instalado o con lista vacía. Si un archivo está corrupto o no se puede leer, registra **y lo declara** en una nota.
- Canal: `maybeRegisterStdioMCP` propaga las notas. `externalStep` las emite como `StepStatusSkipped` con el mensaje, o a stderr. Headless las imprime inline (`↷`) y en una sección final `Notes:`. La TUI las guarda en la fila sin cambiar el estado y las lista en la pantalla final.

### 5.4 Nadie borra lo que no puede probar propio (sin código)

Revisión estática de todo `os.Remove`/`os.RemoveAll` fuera de tests:

- **Install**: `roledefs.rollbackActions` sólo borra archivos creados en la misma llamada o restaura backups; los conflictos nunca borran. `externalStep.Rollback` sólo borra un binario que esta corrida descargó (`Downloaded`). `self_install` sólo borra su temporal y el binario viejo apartado. `homebrew.go` borra su `aside`. `filemerge` borra su temporal. `backup.Restore` borra entradas con `Existed=false` en el snapshot, que son cosas creadas por la instalación. El skill install sobrescribe con snapshot previo y `copyDir` no borra extras.
- `role-implementer.md` y `role-senior.md` no están en la lista canónica (senior está prohibido en `ValidateResources`), así que ni install ni uninstall los abren. `role-auditor.md` personal termina en conflicto: queda preservado, la versión del stack va al lateral y se avisa BLOQUEANTE. No hay ningún harness skill `sdd-*`: el único `sdd-*` del catálogo es `sdd-orchestrator` (config, por markers).
- **Hallazgos fuera del install (comando `uninstall`)**: `skillRemovalStep` hace `RemoveAll(SkillsDir/<id de harness>)` y `commandRemovalStep` borra las rutas de los assets embebidos, ninguno con prueba de propiedad por hash. Un skill o comando personal con el mismo nombre que uno del catálogo se borraría. Lo mitiga el snapshot de `uninstall-snapshot` en Prepare (`internal/uninstall/plan.go:73-80`). `roledefs.Uninstall` sí exige el hash del manifiesto, y el lateral sólo se borra si es byte-idéntico al nuestro. No lo cambié porque está fuera del alcance de 5.4 ("sin código").

## JUICIO

- **Conflicto en un rol de flujo: advertencia destacada, no error de install** (constructor, confianza alta). Evidencia: un `Run()` que devuelve error dispara `DefaultRollbackPolicy`. Eso revierte todos los harnesses ya aplicados y el propio lateral `.active-stack.md`, que es lo que el usuario necesita para decidir. Lo mismo pasaba en verify: el RED de 12:10:09 muestra que, con el health nuevo y un check duro, la instalación entera habría fallado. "Imposible no verlo": en headless aparece inline y en el resumen final (C-32 ya lo testea para cualquier degraded). En la TUI quedaba borrado por el bug de sobrescritura, que ahora está corregido y testeado. Verify lo muestra como `[??]` con el texto BLOQUEANTE. Costo: el exit code de headless sigue en 0.
- **Exigir enabled + installed para saltear engram** (confianza media-alta). Cualquiera de los dos solo no prueba que el plugin cargue. Ante duda se registra: un duplicado es menos grave que quedarse sin memoria.
- **Lista de hashes sólo con historial git** (confianza media). Es fiel a "versión ya publicada": la versión actual sin commit no lo está, y si es idéntica igual da noop.
- Comparación con lo anterior: antes de esto, en la máquina de Juani (sin manifiesto) un `install` habría dejado los 7 especialistas como conflictos silenciosos y el `role-auditor` personal sin aviso. Además habría vuelto a escribir `mcpServers.engram` junto al plugin.

## SUPOSICIÓN (no verificado)

- Que Claude Code cargue un plugin sólo con `enabledPlugins[id] == true`. Me baso en la forma real de los archivos de esta máquina, no en documentación.
- Windows GUI (`windows.go`): el degraded se mapea a `step_degraded` y el skipped cae al `default` (`step_output` con detalles). No verifiqué cómo los muestra la GUI ni si un `step_succeeded` posterior los tapa.
- No ejecuté la TUI real ni headless contra el home real (prohibido). La cobertura es por tests con homes temporales.

## Archivos tocados

- Nuevos: `internal/harness/roledefs/shipped.go`, `internal/harness/roledefs/shipped_test.go`, `internal/harness/external/plugin.go`, `internal/harness/external/plugin_test.go`, `internal/agents/claude/adapter_plugin_test.go`, `internal/install/role_definitions_conflicts_test.go`, `internal/install/external_notes_test.go`, `cmd/active-stack/headless/notes_executor_test.go`.
- Modificados (sólo los hunks de este trabajo): `internal/harness/roledefs/{install,health,resources,doc}.go`, `internal/harness/roledefs/health_test.go`, `internal/install/{steps,plan}.go`, `internal/verify/harness_checks.go`, `internal/verify/role_definitions_check_test.go`, `internal/harness/external/{install,mcp,types}.go`, `internal/agents/claude/adapter.go`, `cmd/active-stack/headless/executor.go`, `internal/tui/model.go`, `internal/tui/degraded_tui_test.go`, `openspec/changes/lean-agent-workflow/tasks.md` (§5), este informe y `informes/LEEME.md`.

## Limitaciones y decisiones abiertas

1. **Commit futuro de los assets**: cuando se commiteen los 7 especialistas modificados y los 5 roles nuevos, `TestShippedRoleHashes_MatchGitHistory` va a fallar hasta regenerar `shipped.go` en ese mismo cambio (el test imprime la lista exacta). La alternativa es aceptar también los hashes del árbol de trabajo actual. Eso deja el commit verde, pero cada edición de un asset rompería el test mientras se desarrolla. Decide el orquestador o Juani.
2. Exit code de headless con un conflicto BLOQUEANTE: hoy 0, igual que los degraded de C-32. Si se quiere un código distinto (por ejemplo para CI), es un cambio aparte.
3. El encabezado del resumen sigue diciendo "Degraded (best-effort) harnesses", también para role-definitions. El mensaje propio aclara el caso, pero la etiqueta no es exacta.
4. `roledefs.Uninstall` no reconoce hashes publicados: sólo borra con el manifiesto. Tras un install que actualizó por hash publicado, el manifiesto ya queda bien.
5. El salto de engram sólo aplica al target máquina con el adapter Claude concreto. El wrapper de target proyecto no implementa `PluginMCPAdapter` y registra como antes. `installMCP` (método `mcp`, context7) no consulta plugins porque ningún plugin conocido lo provee.
6. La limpieza del duplicado actual en `~/.claude.json` (`mcpServers.engram`) y de `role-implementer.md`, `role-senior.md` y el `role-auditor.md` personal queda para Juani a mano, como pide 5.4.

## Chequeos no ejecutados

`go build` y una instalación real (prohibidos), la GUI de Windows, `openspec validate`, y una revisión independiente (rompedor/auditor) de este diff.
