# Constructor — paso 4 (`instalar-roles`)

Rama `feat/metodologia-fase-2`. Propose firmado `d26131e`
(`docs/propuestas/2026-10-23-fase-2-pasos-1-4.md`, sección "## 4 — `instalar-roles`").
MODO RÁPIDO: mínimo que cumple, el resto a deuda.

## Qué se construyó

Harness `role-definitions`: instala los siete `role-*.md` canónicos
(`role-dba`, `role-backend`, `role-frontend`, `role-devops`, `role-qa`,
`role-architect`, `role-analyst`, vía `internal/roles` para no duplicar la
lista) como archivos reales en `~/.claude/agents/`, con ciclo de vida
completo: ausente crea, idéntico no-op, propio intacto actualiza con backup,
personal/desconocido/modificado se preserva intacto y nuestra versión se
escribe en `<nombre>.active-stack.md`. Nunca borra el directorio
`~/.claude/agents/` (compartido con agentes personales). Rechaza symlinks
externos, hace un chequeo de concurrencia (re-lee antes de escribir), y
revierte sus propias escrituras parciales si un adapter posterior falla en
la misma llamada (QA-2).

## Desviación deliberada del texto literal del propose

El propose pide "directorio propio bajo SkillsDir resuelto por capability
del adapter" — el mecanismo de `internal/harness/skill` (SKILL.md +
recursos en `<SkillsDir>/<harnessID>/`). La instrucción operativa de esta
tarea de apply fija el escenario real a probar explícitamente sobre
`~/.claude/agents/`, con agentes personales preexistentes ahí, y editado a
mano con el mismo nombre que un rol.

JUICIO (constructor, confianza alta): seguí la instrucción operativa, no el
texto literal del propose, por una razón funcional y no sólo de obediencia:
los siete roles existen para ser convocados como `subagent_type` por el
Agent tool, y esa resolución en Claude Code lee archivos PLANOS de
`~/.claude/agents/`, no un `SKILL.md` dentro de `~/.claude/skills/`. Meter
los siete ahí los deja legibles pero no los registra como subagente
invocable — el mecanismo que el propose nombra no produce el efecto que
busca. Documentado en `internal/harness/roledefs/doc.go`.

Consecuencia: no hay índice `SKILL.md` (soltar uno dentro de
`~/.claude/agents/` es terreno no probado: Claude Code podría intentar
parsearlo como definición de subagente). El "índice" se resolvió como un
manifiesto de ownership JSON (`~/.active-stack/state/role-definitions/
<agent>.json`), fuera de `~/.claude/agents/` — ver deuda D-16.

## AR-2 (Deuda aceptada del propose)

"el enlace posterior a Compose no tiene módulo propietario ni dependencias
declaradas; asignar propietario y declarar que config no depende de
install/verify."

- Propietario asignado: `internal/harness/roledefs` (doc.go, comentario de
  paquete). Es el paso que, después de que `config.Compose` devuelve texto
  puro sobre los siete roles, materializa los recursos en disco.
- Dirección de dependencia verificada, no sólo declarada:
  `internal/harness/config/dependency_direction_test.go` escanea con
  `go/parser` los imports reales de los `.go` de producción de
  `internal/harness/config` y falla si alguno importa `install`,
  `uninstall`, `verify` o `roledefs`. Mutante real: agregar
  `import _ ".../internal/install"` a un archivo de ese paquete no sólo
  pone el test en rojo — directamente no compila (ciclo de import, porque
  `install` ya importa `config`). Para `verify`, que SÍ importa `config`
  hoy (`internal/verify/delivery.go`) sin que exista el ciclo inverso, el
  test es la única guarda real: ver "Rojo/verde" más abajo.

## QA-2 (Deuda aceptada del propose)

"falta el caso de fallo entre adapters del mismo paso: el primero escribe,
el segundo falla; exigir recuperación del primero y conservación de
ediciones posteriores."

Cubierto con test, en dos variantes (`create` y `update`) y a nivel de
paquete (`internal/harness/roledefs/qa2_test.go`) y a nivel de wiring real
del pipeline (mutantes sobre `plan.go`, ver tabla). La razón de por qué el
rollback genérico del pipeline NO alcanza: `internal/pipeline/
rollback.go:ExecuteRollback` sólo repite `Rollback()` sobre pasos cuyo
`StepResult.Status == StepStatusSucceeded` — un paso que falló en su propio
`Run()` nunca llega a ese estado, así que su propia recuperación tiene que
vivir adentro de `Install`, no en el genérico. `Install` lleva una lista de
"undo actions" (crear → borrar; actualizar → restaurar desde el snapshot
tomado antes de escribir) y la reproduce en reversa ante cualquier error,
across TODOS los adapters ya procesados en esa llamada — no sólo el que
falló.

## Rojo/verde (tabla, con hora)

| Comportamiento | Comando | Hora | Resultado |
|---|---|---|---|
| `DefaultResources`/`ValidateResources` no existen | `go vet ./internal/harness/roledefs/...` | 2026-10-23 13:44:43 | RED (compilación) |
| idem, implementados | `go test ./internal/harness/roledefs/... -v` | 2026-10-23 13:45:02 | GREEN (5 tests) |
| `loadManifest`/`saveManifest` no existen | `go vet ./internal/harness/roledefs/...` | 13:45:26 | RED (compilación) |
| idem, implementados | `go test ./internal/harness/roledefs/... -v` | 13:45:42 | GREEN (8 tests) |
| `Install` no existe | `go vet ./internal/harness/roledefs/...` | 13:47:45 | RED (compilación) |
| idem, implementado (7 escenarios: crea, no-op, update+backup, conflicto personal, conflicto propio-editado, skip capability, symlink externo) | `go test ./internal/harness/roledefs/... -v` | 13:48:23 | GREEN (15 tests) |
| QA-2 mutante: `rollbackActions` anulado | `go test ./internal/harness/roledefs/... -run SecondAdapter -v` | 13:49:36 | RED por aserción (2 tests, archivos de claude no revertidos) |
| QA-2 restaurado | idem | 13:49:40 | GREEN |
| Concurrencia mutante: re-lectura anulada | `go test ./internal/harness/roledefs/... -run ConcurrentEdit -v` | 13:49:49 | RED por aserción |
| Concurrencia restaurada | idem | 13:49:54 | GREEN (suite completa del paquete, 19 tests) |
| `Uninstall` no existe | `go vet ./internal/harness/roledefs/...` | 13:50:41 | RED (compilación) |
| idem, implementado | `go test ./internal/harness/roledefs/... -v` | 13:51:08 | GREEN (23 tests) |
| `Health` no existe | `go vet ./internal/harness/roledefs/...` | 13:51:32 | RED (compilación) |
| idem, implementado | `go test ./internal/harness/roledefs/... -v` + gofmt + vet | 13:51:43 | GREEN (26 tests), gofmt/vet limpios |
| AR-2: import de `install` agregado a mano a `internal/harness/config` | `go test ./internal/harness/config/ -run TestConfigPackage_DoesNotImportInstallOrVerify -v` | 13:52:50 | no compila (ciclo real) — evidencia de que el test es redundante para `install`, no de que esté roto |
| AR-2 real (sin el import agregado) | idem | 13:53:04 | GREEN |
| Catálogo: conteos duros rotos por la entrada nueva | `go test ./internal/catalog/... -v` | 13:53:35 | RED (2 tests, `wantFoundationGlobal`/`wantCount` desactualizados) |
| Conteos bumped 16→17 / 6,16,16→7,17,17 | idem | 14:04 (re-verificado en la corrida final) | GREEN |
| Wiring install: dispatch `role-definitions` deshabilitado a mano en `plan.go` | `go test ./internal/install/... -run RoleDefinitions -v` | 13:59:28 | RED por aserción (7 archivos no escritos + step id cae en `config:role-definitions`) |
| Wiring install restaurado | idem | 13:59:34 | GREEN |
| Wiring uninstall: dispatch deshabilitado a mano en `plan.go` | `go test ./internal/uninstall/... -run RoleDefinitions -v` | 14:00:31 | RED por aserción (`role-dba.md` seguía presente) |
| Wiring uninstall restaurado | idem | 14:00:35 | GREEN |
| Wiring verify: dispatch deshabilitado a mano en `harness_checks.go` | `go test ./internal/verify/... -run RoleDefinitions -v` | 14:02:21 | RED por aserción (2 de 3 tests: cae en `checkConfig` genérico, busca `InstructionsPath`) |
| Wiring verify restaurado | idem | 14:02:26 | GREEN |
| CLI roles: flag `--home` no existe | `go test ./cmd/active-stack/... -run 'TestRoles_' -v` | 14:03:36 | RED (flag desconocida) |
| CLI roles implementado | idem | 14:04:35 (tras corregir dos tests preexistentes acoplados al home real, ver Aprendizaje) | GREEN |

## Aprendizaje (para quien reuse esto)

Agregar el chequeo de definiciones a `active-stack roles` con
`os.UserHomeDir()` como default rompió DOS tests preexistentes
(`TestRolesConArchivosExplicitosNombraLosDosLados`,
`TestRolesJSONEsConsumibleSinParsearTexto`): dependían implícitamente de que
el HOME real de quien corre la suite no tuviera (o sí tuviera, daba igual
antes) roles instalados, y de golpe empezaron a fallar según el estado de
`~/.claude/agents/` de la máquina. Es el mismo patrón de acoplamiento
ambiental que ya documentaron los pasos 1 y 2a con el entorno de test más
permisivo que el real, pero en la dirección inversa: acá el entorno de test
se volvió sensible a una máquina real que el test nunca debería mirar.
Arreglado haciendo que ambos tests instalen un home controlado con
`roledefs.Install` y lo pasen por `--home`, en vez de depender del default.

## Primer uso real (orquestador/constructor, 2026-10-23 14:06 a 14:08)

Binario compilado a
`/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/apply-p4/bin/active-stack`.
`$FAKEHOME1` con `~/.claude/agents/` ya poblado por Juani (simulado): un
agente personal sin relación (`my-deploy-helper.md`) y uno con el MISMO
NOMBRE que un rol (`role-dba.md`), editado a mano, con su propio criterio.

**Antes:**
```
$ ls .claude/agents/
my-deploy-helper.md  role-dba.md
$ md5sum .claude/agents/*.md
e0bdf1b2...  my-deploy-helper.md
3792f069...  role-dba.md
```

**Instalación real:**
```
$ active-stack install --headless --mode custom --custom role-definitions \
    --agent claude --home "$FAKEHOME" --yes --no-self-install
  → snapshot running...        ✓
  → permissions:permissions    ✓   (forzado por C-21 en modo custom, no es este harness)
  → role-definitions:role-definitions   ✓
  → verify-hook                ✓
exit=0
```

**Después:**
```
$ ls .claude/agents/
my-deploy-helper.md  role-analyst.md  role-architect.md  role-backend.md
role-dba.md  role-dba.md.active-stack.md  role-devops.md  role-frontend.md
role-qa.md
```

- `my-deploy-helper.md` y `role-dba.md`: md5 **idéntico** al de antes
  (`e0bdf1b2...` y `3792f069...`) — nada personal se perdió.
- `role-dba.md.active-stack.md`: nuestra versión, aparte, con el frontmatter
  placeholder (`description: PLACEHOLDER — ...`).
- Los otros 6 roles: creados de cero.
- Manifiesto de ownership (`.active-stack/state/role-definitions/
  claude.json`): contiene los 6 hashes de lo que SÍ es nuestro;
  `role-dba.md` no aparece (correcto: no es nuestro).

**Reinstalación (idempotencia real):**
```
$ md5sum .claude/agents/*.md > antes.txt
$ active-stack install ... (mismo comando)   exit=0
$ md5sum .claude/agents/*.md > despues.txt
$ diff antes.txt despues.txt
(sin diferencias — IDÉNTICO byte a byte, incluidos los dos personales)
```

**Desinstalación real:**
```
$ active-stack uninstall --mode custom --custom role-definitions \
    --agent claude --home "$FAKEHOME" --yes
  → uninstall-snapshot                        ✓
  → role-definitions-removal:role-definitions ✓
exit=0
$ ls .claude/agents/
my-deploy-helper.md  role-dba.md
$ md5sum .claude/agents/*.md
e0bdf1b2...  my-deploy-helper.md     (idéntico)
3792f069...  role-dba.md             (idéntico, sigue siendo el mío)
```
Los 6 roles propios y el `.active-stack.md` de conflicto se borraron; los
dos personales (incluido el que colisiona de nombre) sobrevivieron
byte a byte. Segunda desinstalación (idempotencia): mismo resultado, `exit=0`,
sin error sobre archivos que ya no están.

**CLI `roles` bloqueando de verdad (home sin roles instalados, diff real):**
```
$ active-stack roles --home "$FAKEHOME" --files migrations/001_x.sql
role-dba   migrations/001_x.sql
role-qa    siempre
...
bloqueado: definición de rol faltante o vacía — role-qa.md (missing)
exit=3
```
Tras reinstalar `role-definitions` en ese mismo home, la misma invocación
sale `exit=0`.

**Segundo home, tres agentes (`claude,codex,opencode`), con config previa en
los tres:**
```
$ active-stack install --headless --mode custom --custom role-definitions \
    --agent claude,codex,opencode --home "$FAKEHOME2" --yes --no-self-install
exit=0
```
- `role-backend.md` (personal, con el mismo nombre que un rol): preservado
  byte a byte ("Mío. No tocar."); nuestra versión en
  `role-backend.md.active-stack.md`.
- `.codex/config.toml` y `.config/opencode/opencode.json` SÍ cambiaron —
  pero por el harness `permissions` (forzado en modo custom, C-21), no por
  `role-definitions`: confirmado leyendo el contenido (agrega
  `approval_policy`/`sandbox_mode` a Codex y un `permission.bash` a
  OpenCode, preservando la clave `"mcp": {}` que ya tenían). `role-
  definitions` en sí mismo no escribió nada en ninguno de los dos, porque
  `AgentsDir()` devuelve `""` para ambos hoy (capability no soportada,
  J3/J4) — confirmado porque ningún archivo nuevo apareció bajo `.codex/`
  ni `.config/opencode/` más allá del cambio de `permissions`.

JUICIO (constructor, confianza alta): el mecanismo hace exactamente lo que
el paso pedía probar — personal preservado, colisión de nombre resuelta sin
pérdida, idempotencia real, capability no soportada correctamente
silenciosa — en dos homes distintos y dos combinaciones de agentes.

## JUICIO

- Instalar en `~/.claude/agents/` en vez de bajo `SkillsDir` — autor: yo.
  Confianza alta en que es lo correcto FUNCIONALMENTE (ver desviación
  arriba); baja en que sea lo que el propose "quiso decir" literalmente, si
  alguien lo escribió pensando en el mecanismo de skills.
- Exit code 3 para "definición faltante/vacía de un rol convocado" en
  `active-stack roles` — autor: yo, valor arbitrario. El propose pide que
  bloquee, no fija un código. Confianza media: no colisiona con 0/1
  existentes, pero no está documentado en ningún lado central de códigos de
  salida del CLI (si existe tal lista, no la encontré).
- El manifiesto de ownership vive en `~/.active-stack/state/role-
  definitions/<agent>.json`, no en `~/.active-stack/backups/...` — autor:
  yo. Es un dato persistente, no un backup puntual; sigue el precedente de
  `internal/i18n` (`~/.active-stack/config.json`). Confianza alta.
- El archivo de conflicto se detecta "nuestro" por comparación directa de
  contenido contra el candidato actual, sin pasar por el manifiesto de
  ownership (a diferencia del archivo principal) — autor: yo. Es más simple
  y en la práctica correcto (nadie más escribe un archivo llamado
  `role-x.md.active-stack.md` por coincidencia), pero es una asimetría con
  el resto del diseño. Confianza media.

## SUPUESTO (no verificado)

- Asumí que "capability del adapter" para `AgentsDir` debía agregarse como
  un método NUEVO en cada adapter concreto (`claude`, `codex`, `opencode`),
  no en la interfaz compartida `agents.Adapter` — para no ensanchar esa
  interfaz y mantener el radio de impacto bajo, igual que
  `skill.AgentAdapter` ya hace. No verifiqué si hay una convención escrita
  en otro lado que prefiera lo contrario.
- No probé el flujo con `--mode lite` ni `--mode full` reales (sólo
  `--mode custom --custom role-definitions`): el catálogo declara
  `install_modes: [lite, full]`, y los tests de catálogo confirman el
  conteo, pero no corrí el binario real en esos dos modos.
- No verifiqué que `active-stack verify` (el subcomando, no el hook interno
  de instalación) efectivamente invoque `checkRoleDefinitions` end-to-end:
  el "Verification checks: 0 passed, 0 failed" que mostró la corrida real
  sugiere que el `verify-hook` del pipeline de instalación no recorre el
  catálogo igual que `internal/verify.ChecksForHarness` — no investigué esa
  discrepancia porque está fuera del alcance de este paso (es una pregunta
  sobre el wiring general de verify, no sobre `role-definitions`
  específicamente), pero la dejo nombrada en vez de callada.
- No probé el caso de un `~/.claude/agents/` que no existe en absoluto
  (padre tampoco existe) en la corrida real — sólo en tests unitarios
  (`TestInstall_AbsentCreates` lo cubre con `t.TempDir()`, que tampoco tiene
  el padre).

## Deuda nueva (modo rápido)

- **D-15** — Contenido placeholder: fuente y licencia definitivas de los
  siete roles (J3/J4) siguen sin decidir. El mecanismo de instalación es
  real; el texto que viaja dentro de cada `role-*.md` no lo es.
- **D-16** — No hay índice `SKILL.md` como pide el propose; se sustituyó por
  el manifiesto de ownership JSON. Si en algún momento se decide que el
  índice tiene que ser legible por humanos en Markdown (no sólo
  programático), falta ese artefacto.
- **D-17** — `AgentsDir` sólo está confirmado para Claude Code. Codex y
  OpenCode devuelven `""` (capability no soportada) hasta que J3/J4 decidan
  si tienen un equivalente real y cuál es.
- **D-18** — El gate previo del propose ("Fuente/licencia, siete
  definiciones canónicas, nombre/directorio y modos J3/J4 (TBD) son gates
  antes de bundlear") no se resolvió; se documentó como bloqueado y se
  avanzó sólo con el mecanismo, consistente con MODO RÁPIDO.
- **D-19** — El archivo de conflicto (`*.active-stack.md`) se identifica
  como "nuestro" por igualdad de contenido directa, no por el mismo
  mecanismo de manifiesto de ownership que el archivo principal (ver
  JUICIO). Si alguna vez el contenido candidato cambia de versión a
  versión, un `*.active-stack.md` de una versión VIEJA nuestra puede quedar
  en el disco sin que `Uninstall` lo reconozca como propio (porque ya no
  matchea el candidato ACTUAL) — no hay test para ese caso específico.
- **D-20** — El desacople entre el `verify-hook` del pipeline de instalación
  y `internal/verify.ChecksForHarness` (ver SUPUESTO) no se investigó; puede
  significar que el health check de `role-definitions` escrito en este paso
  nunca corre en una instalación real vía `active-stack install`, sólo
  cuando algo externo llama a `ChecksForHarness` directamente (p. ej. un
  futuro comando `active-stack verify`).
- Hereda, sin resolver: **AR-1**, **AR-3**, **BE-4**, **QA-1**, **D-11** a
  **D-14** (pasos anteriores), ninguno tocado por este paso.

## EVIDENCIA

- Build del seam: `go build ./internal/harness/roledefs/...` antes de
  escribir ningún test, implícito en cada ciclo RED documentado arriba (el
  primer RED de cada archivo nuevo es de compilación, nunca de build roto
  por otra razón).
- RED por aserción (mutante, no sólo RED-antes-de-escribir) para QA-2,
  concurrencia, AR-2, y los tres puntos de wiring (install/uninstall/
  verify) — tabla arriba, con hora y comando, cada uno restaurado y
  reverificado en verde inmediatamente después.
- Suite completa verde repetidas veces durante el desarrollo (28 paquetes
  con sólo `roledefs` agregado, antes de cablear catalog/install/uninstall/
  verify/CLI) y una vez más, completa, sobre el commit final (29 paquetes,
  ver "Verde final").

## Verde final

```
$ git rev-parse HEAD
2526eee1fd8de6ee95ec3bcba7cc1b0fed889509
$ gofmt -l .
(sin salida)
$ go vet ./...
(sin salida, exit 0)
$ go test ./... -count=1
ok  	.../assets                              0.041s
ok  	.../cmd/active-stack                    0.968s
ok  	.../cmd/active-stack/headless           2.585s
ok  	.../e2e/starter                         0.807s
ok  	.../internal/agents                     0.010s
ok  	.../internal/agents/claude              0.017s
ok  	.../internal/agents/codex               0.011s
ok  	.../internal/agents/opencode            0.012s
ok  	.../internal/approval                   1.165s
ok  	.../internal/backup                     0.124s
ok  	.../internal/catalog                    0.216s
?   	.../internal/contracttest               [no test files]
ok  	.../internal/filemerge                  0.010s
ok  	.../internal/harness/command            0.019s
ok  	.../internal/harness/config            27.805s
ok  	.../internal/harness/config/permissions 0.156s
ok  	.../internal/harness/external           0.100s
ok  	.../internal/harness/roledefs           0.094s
ok  	.../internal/harness/skill              0.027s
ok  	.../internal/i18n                       0.025s
ok  	.../internal/install                    0.304s
ok  	.../internal/model                      0.007s
ok  	.../internal/pipeline                   0.008s
ok  	.../internal/planner                    0.008s
ok  	.../internal/reports                    0.519s
ok  	.../internal/roles                      0.128s
ok  	.../internal/system                     0.036s
ok  	.../internal/tui                        0.034s
ok  	.../internal/uninstall                  0.027s
ok  	.../internal/verify                     0.044s
```
2026-10-23T14:10:01-03:00 a 14:10:30-03:00, sobre `2526eee`. 29 paquetes, 0
fallas. Baseline del paso 3 (mismo comando, 28 paquetes ok): la única
diferencia es `internal/harness/roledefs`, que no existía — ningún paquete
preexistente se rompió.

## Roles no convocados en este paso (JUICIO, confianza alta)

Informe de **apply**, no de **verify**: por el contrato del propose ("nadie
convoca arquitecto/QA/auditor en este trabajo; eso es la revisión prevista en
verify, por paso"), no corresponde convocar revisores acá. Arquitecto/
backend/DBA/frontend/DevOps/QA quedan para la vuelta de `verify` de este
paso; este informe no la sustituye.

## Cierre de D-15 (orquestador, 23/10/2026)

J4 decidida por Juani (ver `segunda-revision/06-decisiones-juani.md`): los placeholders se reemplazaron por sus `role-*.md`. `go test ./internal/harness/roledefs` verde después del reemplazo.
