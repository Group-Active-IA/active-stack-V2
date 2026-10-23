# Constructor — paso 3 (`opencode-permission-replace`)

Fecha: 2026-10-23 · Rama: `feat/metodologia-fase-2` · Modo rápido.

Implementa **únicamente** la sección `## 3 — opencode-permission-replace` de
`docs/propuestas/2026-10-23-fase-2-pasos-1-4.md`, respetando "Contrato común y
modo rápido" y "Evidencia, entrega e impacto", aprobada por Juani en el commit
`d26131e`. Cubre QA-1 ("JSONC válido no se puede romper") con test, por
instrucción explícita de este paso. QA-2 (fallo entre adapters) queda deuda
aceptada, sin tocar.

## Commits

| SHA | Mensaje |
|---|---|
| `aaadd10` | `fix(opencode): reemplazar solo permission.bash, no el objeto permission entero` |
| (este informe) | `docs(propuestas): informe del constructor del paso 3` |

Archivos: `internal/filemerge/{json_merge.go,json_merge_test.go}`,
`internal/harness/config/permissions/{overlays.go,install.go,
permissions_test.go,tier_switch_test.go}`.

## Qué construí

En `overlays.go`, `composeOpencodeOverlay` ya no marshalea el struct
`opencodeSettingsShape` directo: construye el overlay como `map[string]any`
con las claves escalares de `permission` en plano (igual que antes: cada tier
las declara todas, así que un merge normal ya las pisa) y envuelve **solo**
`permission.bash` en `{"__replace__": shape.Bash}` — mismo centinela que ya
usa `buildStdioOverlay` en `internal/harness/external/mcp.go` para las
entradas MCP. `filemerge.mergeOrdered` ya sabía consumir ese centinela antes
de escribir (no hubo que tocar `filemerge` para esto); lo que no existía
todavía era usarlo sobre `permission.bash` en vez de sobre `permission`
entero.

La diferencia deliberada contra la alternativa rechazada por el propose (el
fix de otra rama, commit `ddeb94f`, que envuelve `permission` COMPLETO):
reemplazar solo `bash` deja vivas las claves hermanas de `bash` dentro de
`permission` que ningún tier administra — una regla personal, o un campo que
una versión futura de opencode agregue. `ddeb94f` las destruye.

Pérdida aceptada: una regla agregada a mano **dentro** de `permission.bash`
(no hermana de `bash`, sino entrada del objeto que `bash` reemplaza entero) no
sobrevive a un reinstall. El propose la acepta explícitamente ("pérdida
deliberada de bash personal exige diff/aprobación y snapshot") y deja el gate
de mostrar ese diff + aprobación como TBD (J6) — no lo implementé; ver deuda.

En `filemerge/json_merge.go`, agregué `ValidateJSONObject(raw []byte) error`:
reusa `decodeOrderedOrEmpty` (la misma normalización JSONC — comentarios,
comas finales — que `MergeJSONObjects` ya aplica a la base) pero devuelve el
error de parseo en vez de tragarlo como `{}`. `MergeJSONObjects` no cambió: su
contrato de "base rota se trata como vacía, el overlay prevalece" sigue igual
para todos los demás consumidores (claude, gemini, vscode).

En `install.go`, `installOne` llama a `mergeOpencodeJSONFile` (nueva función,
específica de opencode, con el mismo shape que `mergeCodexTOMLFile` ya tiene
para TOML) en vez de `mergeJSONFile` genérico: lee la base, la valida con
`ValidateJSONObject` **antes** de mergear, y si es inválida devuelve error sin
escribir nada — sin el fallback-a-vacío que `MergeJSONObjects` aplicaría por
default. El backup (`backupIfExists`) sigue corriendo antes, igual que para
todos los agentes; no toqué ese contrato compartido.

## Rojo / verde por caso

| Caso | Comando | Hora | Resultado |
|---|---|---|---|
| RED — `filemerge.ValidateJSONObject` no existe (seam) | `go test ./internal/filemerge/... -run ValidateJSONObject -count=1 -v` | 2026-10-23T13:23:29-03:00 | `undefined: ValidateJSONObject` (build failed) — 4 referencias en el test |
| GREEN tras implementarla | mismo comando | 2026-10-23T13:23:57-03:00 | PASS — 3/3 (rechaza JSON roto, acepta JSONC con comentarios/comas finales, acepta vacío) |
| Suite completa de `filemerge` sin romper nada preexistente | `go test ./internal/filemerge/... -count=1 -v` | 2026-10-23T13:24:02-03:00 | PASS |
| RED por aserción — matriz 3×3 de tiers exigiendo ausencia de residuo en `permission.bash`, y pérdida de una regla bash personal (`TestOpencodeTierSwitchFullyReplacesBash`, `TestOpencodeTierSwitchReplacesPersonalBashRule`) contra el código SIN el fix | `go test ./internal/harness/config/permissions/... -run 'TestOpencodeTierSwitch' -count=1 -v` | 2026-10-23T13:24:55-03:00 | FAIL — 6/9 casos de la matriz con residuo real (ej. `balanceado → estricto` conserva `go test *`/`go build *`/`git status`/`git diff *`); la regla personal sobrevivía (se esperaba que no); la preservación de la clave hermana ya pasaba (no es RED de este fix, es red de no-regresión contra `ddeb94f`) |
| GREEN tras envolver `permission.bash` en `__replace__` | `go test ./internal/harness/config/permissions/... -count=1 -v` | 2026-10-23T13:25:40-03:00 | PASS — 35/35 casos top-level, incluidos todos los preexistentes (idempotencia, `external_directory`, deny-last, shape de claude no filtra a opencode) |
| RED por aserción — base inválida no abortaba (`TestInstallOpenCode_InvalidBaseJSONAbortsWithoutWriting`) contra el `install.go` sin el wiring | `go test ./internal/harness/config/permissions/... -run 'TestInstallOpenCode_' -count=1 -v` | 2026-10-23T13:26:37-03:00 | FAIL en ese caso (`Install() error = nil`); los otros dos (JSONC válido, backup fallido) ya pasaban porque no dependían del fix pendiente |
| GREEN tras wirear `mergeOpencodeJSONFile` en `installOne` | `go build ./...` + mismo comando | 2026-10-23T13:26:55-03:00 / 13:27:06-03:00 | build OK; PASS 3/3, paquete completo 35/35 |
| gofmt / vet sobre el paquete tocado | `gofmt -l .` + `go vet ./...` | 2026-10-23T13:27:23-03:00 | sin salida (limpio) ambos |
| Suite completa antes de la corrida real | `go test ./... -count=1` | 2026-10-23T13:27:30 a 13:28:01-03:00 | 28 paquetes ok, 0 fallas |

## Triangulación

- Matriz 3×3 de tiers **incluida la diagonal** (mismo tier → mismo tier): 9
  casos, todos afirmando el contenido exacto y el orden exacto de
  `permission.bash`, no solo "no hay residuo".
- Caso de preservación: una clave hermana de `bash` dentro de `permission`
  (`mcp_tool_x`) sobrevive a cualquier tier — es la prueba que distingue este
  fix de `ddeb94f`.
- Caso de pérdida aceptada: una regla DENTRO de `permission.bash` no
  sobrevive — documenta el comportamiento para que no se lo lea como bug
  nuevo.
- `ValidateJSONObject`: tres casos (rechaza JSON roto, acepta JSONC con
  comentarios + comas finales, acepta vacío/ausente).
- A nivel `Install()`: base inválida aborta sin escritura (bytes y permisos
  del archivo intactos), JSONC válido con claves personales sigue mergeando,
  backup fallido aborta la escritura.

## Mutante (vía RED real, no revertido)

No aplico un mutante manual separado porque el propio ciclo RED→GREEN de
`TestOpencodeTierSwitchFullyReplacesBash` contra el código sin el fix **es**
el mutante: la matriz falló en 6 de 9 casos con el residuo real documentado
por el bug conocido (mismo síntoma que `TestOpencodeBashResidueSurvivesTierSwitch`,
el test que este trabajo reemplaza), y pasó limpio tras el fix. Es la misma
aserción, antes y después del cambio de producción — equivalente a invertir
el fix y confirmar que el test lo detecta.

## Primer uso real (binario real, HOME falso con `opencode.json` realista)

Instrucción explícita de este paso, por la lección de los pasos 1 y 2a ("tests
verdes no acreditan que el comando ande en un entorno tan permisivo como el
real"). Compilé el binario, armé un `$HOME` falso aparte (no el real, no
`~/.claude`) y un segundo `$HOME` falso (fake home del propio `active-stack`,
para que el binario no toque `~/.local/bin` al resolver su propio estado) y
corrí la instalación end-to-end contra un `opencode.json` escrito a mano con
comentarios, comas finales, claves personales (`theme`, `autoupdate`, `model`,
un bloque `mcp.notion`) y un bloque `permission` previo (tier estricto) con
**dos reglas agregadas a mano**: una clave hermana de `bash`
(`mcp_tool_notion_search`) y una regla dentro de `bash`
(`my-deploy-script.sh *`).

```
$ go build -o .../active-stack ./cmd/active-stack && echo "build OK"
build OK
```

Comando de instalación real (sin TUI, sin self-install, sin red): descubrí que
`--agent opencode` solo no alcanza — `install.BuildPlan` no genera ningún
step sin un `--mode`/`--custom` resuelto — así que usé
`--mode custom --custom permissions` para acotar el plan exactamente al
harness de permisos (confirmado con `--dry-run` antes de ejecutar: 3 steps,
`snapshot` / `permissions:permissions` / `verify-hook`, nada de red ni
self-install).

**ANTES** (`$FAKEHOME/.config/opencode/opencode.json`, `md5=4475ea79...`,
`mode=644`):

```jsonc
{
  // config personal de Juani, escrita a mano en varias sesiones
  "theme": "dracula",
  "autoupdate": false,
  "model": "anthropic/claude-sonnet-5",

  "mcp": {
    "notion": {
      "type": "local",
      "command": ["notion-mcp-server"],
      "enabled": true,
    },
  },

  // instalacion anterior en tier estricto + dos reglas agregadas a mano
  "permission": {
    "read": "ask", "edit": "ask", "glob": "ask", "grep": "ask", "list": "ask",
    "external_directory": "deny", "webfetch": "ask", "websearch": "ask", "lsp": "ask",
    "mcp_tool_notion_search": "allow", // clave personal hermana de bash
    "bash": {
      "*": "ask",
      "my-deploy-script.sh *": "allow", // regla bash agregada a mano
      "rm -rf *": "deny",
      "rm -rf /*": "deny",
    },
  },
}
```

```
$ active-stack install --headless --mode custom --custom permissions \
    --agent opencode --tier balanceado --home "$FAKEHOME" --yes --no-self-install
  → snapshot running...
  ✓ snapshot
  → permissions:permissions running...
  ✓ permissions:permissions
  → verify-hook running...
  ✓ verify-hook
Verification checks: 0 passed, 0 failed, 0 warnings, 0 skipped
Installation verified. You're ready to start building.
exit=0
```

**DESPUÉS** (mismo archivo, release ya construido con el fix):

```json
{
  "theme": "dracula",
  "autoupdate": false,
  "model": "anthropic/claude-sonnet-5",
  "mcp": { "notion": { "type": "local", "command": ["notion-mcp-server"], "enabled": true } },
  "permission": {
    "mcp_tool_notion_search": "allow",
    "bash": {
      "*": "ask", "go test *": "allow", "go build *": "allow",
      "git status": "allow", "git diff *": "allow",
      "rm -rf *": "deny", "rm -rf /*": "deny"
    },
    "edit": "ask", "external_directory": "ask", "glob": "allow", "grep": "allow",
    "list": "allow", "lsp": "allow", "read": "allow", "webfetch": "ask", "websearch": "ask"
  }
}
```

Verificado a ojo y por diff automático (`mode=644` sin cambios):

- **Claves personales de nivel superior sobreviven**: `theme`, `autoupdate`,
  `model`, el bloque `mcp.notion` completo — ninguna tocada.
- **La clave hermana personal sobrevive**: `permission.mcp_tool_notion_search
  = "allow"` — es exactamente la diferencia contra `ddeb94f` (reemplazo del
  objeto `permission` entero la hubiera destruido).
- **La regla bash personal NO sobrevive**: `my-deploy-script.sh *` desapareció
  — pérdida aceptada por el propose, confirmada en el mundo real, no solo en
  el test.
- **Los escalares de `permission` pasaron de estricto a balanceado
  correctamente** (`read: ask→allow`, `edit` sigue `ask`, etc.) y `bash` quedó
  exactamente con los patrones de balanceado, sin ningún residuo de estricto.

Reinstalación del mismo tier (idempotencia real, no solo de test):

```
$ md5sum opencode.json   # 39fc3989...
$ active-stack install ... --tier balanceado ...   # mismo comando
$ md5sum opencode.json   # 39fc3989... — IDÉNTICO, byte a byte
```

Cambio de tier real (`balanceado → bypass`) sobre el mismo archivo ya
instalado, para confirmar la matriz fuera del test harness:

```json
"permission": {
  "mcp_tool_notion_search": "allow",
  "bash": { "*": "allow", "rm -rf /*": "deny", "rm -rf ~/": "deny" },
  "edit": "allow", "external_directory": "allow", "glob": "allow", "grep": "allow",
  "list": "allow", "lsp": "allow", "read": "allow", "webfetch": "allow", "websearch": "allow"
}
```

Sin residuo de balanceado (ni `go test *` ni ninguno de sus allow), clave
personal todavía viva.

Backup real confirmado (no solo por el código, por el disco):

```
$ find "$FAKEHOME/.active-stack/backups" -maxdepth 3
.../backups/install/{snapshot.tar.gz,manifest.json}
.../backups/permissions-opencode/{snapshot.tar.gz,manifest.json}
```

Base inválida real, en un segundo `$FAKEHOME` separado:

```
$ printf '{"permission": {"bash": {' > opencode.json   # roto, sin cerrar
$ md5sum opencode.json   # 86dbceb1...
$ active-stack install --headless --mode custom --custom permissions \
    --agent opencode --tier bypass --home "$FAKEHOME2" --yes --no-self-install
  → snapshot running...
  ✓ snapshot
  → permissions:permissions running...
  ✗ permissions:permissions: permissions: merge settings for agent "opencode":
    opencode settings file ".../opencode.json" is not valid JSON/JSONC: unexpected end of json
Installation failed: ...
Rollback: succeeded
exit=1
$ md5sum opencode.json   # 86dbceb1... — IDÉNTICO, nada se escribió
```

### Hallazgo real que el `go test` no atrapó (y por qué)

Antes de llegar a este resultado limpio, la **primera** corrida real del
binario (sin el fix de `task`) escribió `"task": ""` en el `opencode.json` del
`$FAKEHOME`. Ningún test de `go test ./...` lo había atrapado porque la
reescritura de `composeOpencodeOverlay` con un `map[string]any` perdió el
`omitempty` que el struct original tenía en el campo `Task` (que ningún tier
fija nunca) — los tests existentes comparaban claves/valores puntuales
(`read`, `bash["*"]`, etc.), ninguno recorría **todas** las claves de
`permission` buscando un string vacío. Es exactamente el patrón de los dos
pasos anteriores: el entorno de test no reproducía algo que el uso real sí
mostró. Agregué `TestOpencodeNeverWritesEmptyPermissionVerb` (RED confirmado
contra el binario con el bug, ver tabla de rojo/verde arriba), arreglé
omitiendo la clave cuando está vacía, recompilé y repetí toda la corrida real
de arriba — ya no aparece.

## Verde final

```
$ git rev-parse HEAD
aaadd10...
$ gofmt -l .
(sin salida, exit 0)
$ go vet ./...
(sin salida, exit 0)
$ go test ./... -count=1
ok  .../assets                            0.057s
ok  .../cmd/active-stack                  1.108s
ok  .../cmd/active-stack/headless         2.652s
ok  .../e2e/starter                       0.913s
ok  .../internal/agents                   0.031s
ok  .../internal/agents/claude            0.024s
ok  .../internal/agents/codex             0.021s
ok  .../internal/agents/opencode          0.029s
ok  .../internal/approval                 1.282s
ok  .../internal/backup                   0.174s
ok  .../internal/catalog                  0.219s
?   .../internal/contracttest             [no test files]
ok  .../internal/filemerge                0.020s
ok  .../internal/harness/command          0.023s
ok  .../internal/harness/config          28.346s
ok  .../internal/harness/config/permissions 0.194s
ok  .../internal/harness/external         0.065s
ok  .../internal/harness/skill            0.022s
ok  .../internal/i18n                     0.017s
ok  .../internal/install                  0.267s
ok  .../internal/model                    0.011s
ok  .../internal/pipeline                 0.013s
ok  .../internal/planner                  0.012s
ok  .../internal/reports                  0.481s
ok  .../internal/roles                    0.122s
ok  .../internal/system                   0.024s
ok  .../internal/tui                      0.041s
ok  .../internal/uninstall                0.031s
ok  .../internal/verify                   0.048s
```
2026-10-23T13:31:13-03:00 a 13:31:42-03:00, sobre `8bc706a` (HEAD antes del
commit de este paso). 28 paquetes, 0 fallas — mismo baseline que el paso 2a
(28 paquetes ok), ningún paquete preexistente roto. Suite completa repetida
después de cada uno de los dos fixes de este paso (el de `bash` y el de
`task`), siempre verde.

## JUICIO

- Alcanzar el gate de "base inválida aborta sin escritura" **solo para
  opencode** (no para claude/gemini/vscode vía `mergeJSONFile` genérico) —
  autor: yo (constructor). Confianza: alta. El texto del propose está dentro
  de la sección 3, que es específicamente `opencode-permission-replace`;
  generalizar el cambio a `filemerge.MergeJSONObjects` tocaría el contrato
  compartido que "Contrato común" pide reusar, no rediseñar.
- `ValidateJSONObject` vive en `filemerge` (no en `permissions`) — autor: yo.
  Confianza: alta. Reusa `decodeOrderedOrEmpty`, que es interno al paquete; la
  alternativa (reimplementar la misma normalización JSONC en `permissions`)
  hubiera duplicado lógica que ya existe y puede desincronizarse.
- No implementé el gate J6 (mostrar diff + aprobación antes de reemplazar
  `bash`, perdiendo reglas personales) — autor: yo, siguiendo el propose al
  pie de la letra: lo marca TBD explícitamente, y MODO RÁPIDO pide mínimo que
  cumpla. Confianza: alta en que es lo que el texto pide; baja en que el
  mínimo mecánico (reemplazar sin mostrar nada) sea aceptable para un usuario
  real sin el gate — por eso queda en deuda, no en "hecho".
- No toqué "demostrar bash frente a shell en versión fijada de OpenCode"
  (gate previo del propose, también TBD) — mantuve la clave `"bash"`, que es
  la que ya existía; no validé contra ninguna versión real de opencode si esa
  clave sigue siendo la correcta. Confianza: media — es el mismo estado que
  tenía el código antes de este paso, no una regresión introducida acá.

## SUPUESTO (no verificado)

- Asumí que `--mode custom --custom permissions` es la forma correcta de
  acotar el plan al harness de permisos para la corrida real, inferido de
  `SecurityFirstHarnessID = "permissions"` en `internal/install/plan.go` y
  confirmado con `--dry-run` (3 steps, sin red). No hay una forma documentada
  de "instalar solo permisos" en el `main.go`/README; si existe una oficial
  distinta, no la usé.
- No corrí el binario real contra opencode **de verdad** (el proceso real de
  OpenCode no está instalado en este entorno): la corrida real
  confirma que el archivo que `active-stack` escribe es el que yo esperaba
  (bytes, claves, orden), no que OpenCode lo interpreta como yo creo. Eso
  sigue dependiendo del gate TBD "demostrar bash frente a shell en versión
  fijada" que el propose ya marca pendiente.
- Asumí que `os.FileMode` sin cambios (`644` → `644`) es la única garantía de
  "permisos preservados" que corresponde verificar en este paso; no probé con
  un archivo con permisos no estándar (ej. `600`) en la corrida real, aunque
  sí en el test unitario (`TestInstallOpenCode_InvalidBaseJSONAbortsWithoutWriting`
  usa el modo que `os.WriteFile` deja, `0o644`).
- No verifiqué si `opencode` admite hoy una clave de permiso que este shape
  todavía no declara (por ejemplo algo nuevo en una versión reciente); el
  test de preservación de hermanas (`mcp_tool_x`) cubre el mecanismo, no una
  clave real confirmada contra documentación de opencode.

## Deuda nueva (modo rápido)

- **D-11** — Gate J6 (TBD en el propose): reemplazar `permission.bash`
  descarta sin aviso cualquier regla personal agregada a mano dentro de ese
  objeto. No hay mecanismo de mostrar el diff ni de pedir aprobación explícita
  antes de la pérdida — el propose lo deja pendiente de diseño, no lo
  implementé.
- **D-12** — "Cambio concurrente aborta la escritura" (Contrato común,
  sección 3 del `Rojo antes`) no está implementado para opencode ni existe en
  ningún punto de `internal/backup`/`internal/filemerge` hoy: no hay
  detección de que el archivo cambió entre la lectura y la escritura en
  ningún agente, no solo opencode. Lo que SÍ está cubierto (por la escritura
  atómica de `WriteFileAtomic`, preexistente) es que una falla ANTES del
  `os.Rename` deja el archivo original intacto — eso no es lo mismo que
  detectar una edición concurrente y declarar conflicto.
- **D-13** (hereda, no resuelve) — QA-2 del propose (fallo entre adapters del
  mismo paso: el primero escribe, el segundo falla) sigue sin cubrir, tal
  como la instrucción de este paso pedía dejarlo.
- **D-14** — El gate previo del propose ("demostrar bash frente a shell en
  versión fijada de OpenCode") sigue sin resolver: la clave sigue siendo
  `"bash"` (el estado previo a este paso), sin confirmar contra una versión
  real de opencode si ese nombre sigue siendo el correcto o si versiones
  nuevas esperan `"shell"` u otra forma.
- Hereda, sin resolver: **AR-1**, **AR-2**, **AR-3** (arquitectura,
  propose), no tocadas por este paso — no cruzan `permissions`/`filemerge`.

## EVIDENCIA

- RED/GREEN con hora exacta y comando para los tres cambios de producción
  (`ValidateJSONObject`, `__replace__` sobre `bash`, wiring en `installOne`),
  tabla arriba.
- RED/GREEN real (no de test) del bug de `"task": ""`, encontrado en la
  primera corrida del binario contra el `$FAKEHOME` realista — ver sección
  "Hallazgo real".
- Corrida real completa: antes/después del archivo con claves personales
  (top-level y hermana de `bash`) preservadas, regla bash personal perdida
  (aceptada), reinstalación idéntica byte a byte, cambio de tier real sin
  residuo, backup real en disco, base inválida real con exit 1 y archivo
  intacto (md5 idéntico antes/después).
- Suite completa (`go test ./... -count=1`) verde dos veces (antes y después
  del fix de `task`), gofmt y vet limpios, sobre el mismo baseline de 28
  paquetes que el paso 2a.

## Roles no convocados en este paso (JUICIO, confianza alta)

Informe de **apply**, no de **verify**: por el contrato del propose ("nadie
convoca arquitecto/QA/auditor en este trabajo; eso es la revisión prevista en
verify, por paso"), no corresponde convocar revisores acá. Arquitecto/
backend/DBA/frontend/DevOps/QA quedan para la vuelta de `verify` de este
paso; este informe no la sustituye.
