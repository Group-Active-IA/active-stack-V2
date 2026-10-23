# Constructor — paso 2a (captura de informes: núcleo)

Fecha: 2026-10-23 · Rama: `feat/metodologia-fase-2` · Modo rápido.

Implementa **únicamente** la sección `## 2a — Captura de informes: núcleo` de
`docs/propuestas/2026-10-23-fase-2-pasos-1-4.md`, respetando "Contrato común
y modo rápido" y "Evidencia, entrega e impacto", aprobada por Juani en el
commit `d26131e`. NO implementa los adapters 2b (Claude), 2c (Codex) ni 2d
(OpenCode): deja sólo la interfaz (`reports.Source`, `reports.RegisterContract`)
que ellos van a usar.

## Commits

| SHA | Mensaje |
|---|---|
| `ac7b3dda6824a64ced9c6cccd7a4b03da40d6210` | `feat(reports): núcleo de captura de informes antes del despacho` |
| (este informe) | `docs(propuestas): informe del constructor del paso 2a` |

Archivos: `internal/reports/{reports,key,worktree,secrets,manifest,fsops,capture,leeme}.go`
+ sus `_test.go` (`capture_test.go`, `key_test.go`, `secrets_test.go`,
`testutil_test.go`), `cmd/active-stack/reports_dispatch.go` (+ test), wiring
de `reports` en `cmd/active-stack/main.go`, y una entrada nueva en
`.gitignore` para `**/informes/**/.retenidos/`.

## Qué construí

`active-stack reports capture --change <slug> --phase <fase> --round <n>
--role <rol> --dispatch-id <id> --sha <sha> [--diff <rango>]
[--agent/--session/--turn/--provider/--model <v>] [--contract <id>]
(--report <ruta> | --report-text <texto>) [--json]`.

Persiste en `openspec/changes/<change>/informes/<fase>/round-<n>/<rol>/`:
`raw.txt` + `informe.md` (contenido íntegro, dos escrituras separadas
atómicas — crash entre una y otra es recuperable), `manifest.json`
(estructurado: versión, identidad, SHA/diff, contrato, procedencia sólo si
el `Source` la traía, hash, tamaño, estado, causa) y un `LEEME.md` por
change, **regenerado completo** desde todos los `manifest.json` del árbol en
cada captura (nunca un append incremental, para no duplicar ni desordenar
entre escritores de distintas identidades).

Matriz de estado (`internal/reports/capture.go`): `UNSUPPORTED` si
`Source.Contract` no está registrado (nada se escribe); si no, `PENDING` si
el texto está en blanco; si no, `RETAINED` si el escaneo de secretos
encuentra algo (prioridad sobre identidad ambigua — ver `TestReportsCapture_
SecretTakesPriorityOverAmbiguousIdentity`); si no, `PENDING` si falta
`DispatchID`/`SHA` o `Round <= 0`; si no, `COMPLETE`. Un reintento con el
mismo hash de contenido en la misma identidad es no-op (no reescribe
archivos, no toca `Round`); un reintento con bytes distintos en la misma
identidad cae en `conflict-<hash12>/` sin pisar el original.

Seguridad de rutas (`worktree.go`, `key.go`): `RepoDir` tiene que ser
exactamente la raíz del worktree (`git rev-parse --show-toplevel`), no un
subdirectorio ni un worktree ajeno; `Change`/`Phase`/`Role` se validan contra
separadores de ruta y `".."` antes de construir ningún path; cualquier
ancestro ya existente del destino que sea un symlink hacia afuera de la raíz
se rechaza antes de escribir.

Escaneo de secretos (`secrets.go`): **gitleaks no está en el PATH de este
entorno** (`which gitleaks` → nada), así que el escaneo es enteramente
patrones propios — token de GitHub, Slack, AWS access key, bloque de clave
privada, asignación genérica `key/secret/token/password=...`, email,
teléfono. Nunca se vuelca el valor encontrado, sólo el tipo de patrón. Un
informe `RETAINED` escribe el original sólo bajo `.retenidos/` (ignorado por
`.gitignore`, patrón nuevo `**/informes/**/.retenidos/`) y el
`manifest.json`/`LEEME.md` versionados nunca contienen el valor.

## Rojo / verde por caso

| Caso | Comando | Hora | Resultado |
|---|---|---|---|
| Seam compila (sólo tipos, sin lógica de `capture`) | `go build ./internal/reports/...` | 2026-10-23T13:06:24-03:00 | exit 0 |
| **RED** por aserción (no por compilación) — fixture de test usaba `testGit` para `git check-ignore -v`, que trata exit 1 ("no ignorado") como fallo del helper, no como el resultado a afirmar | `go test ./internal/reports/... -run 'Capture\|Worktree\|Atomic\|Recovery' -count=1 -v` | 2026-10-23T13:07:40-03:00 | `TestReportsCapture_OriginalIsVersionable_GitCheckIgnoreReturnsNotIgnored` FAIL: `git [check-ignore ...]: exit status 1` — las otras 20 pasaron contra la implementación ya escrita (ver Aprendizaje) |
| Fix del fixture (helper `gitCheckIgnore` que distingue "ignorado"/"no ignorado"/error real, no de producción) + re-run | mismo comando | 2026-10-23T13:08:13-03:00 | PASS — 21/21 |
| Triangulación (6 tests nuevos: forma literal de `BuildKey`, rechazo de round no positivo, traversal en cualquier segmento, `RegisterContract`, detección de cada patrón de secreto sin volcar el valor, prioridad RETAINED sobre identidad ambigua) | `go test ./internal/reports/... -count=1 -v` | 2026-10-23T13:09:32-03:00 | PASS — 27/27 |
| Gate exacto del propose | `go test ./internal/reports -run 'Capture\|Worktree\|Atomic\|Recovery' -count=1 -v` | 2026-10-23T13:11:46-03:00 | PASS — 21/21 (los 6 de triangulación no matchean el regex, igual que pasó con `approval` en el paso 1; corren en el run completo de arriba) |
| CLI end-to-end (7 casos: feliz, JSON con `status` como string, secreto → código 4 sin fuga, `dir=""` usa `os.Getwd()`, flags faltantes, subcomando desconocido, `--help`) | `go test ./cmd/active-stack/... -run Reports -count=1 -v` | 2026-10-23T13:11:13-03:00 (primera vez) → 2026-10-23T13:11:47-03:00 (tras el mutante, ver abajo) | PASS — 7/7 |
| **Mutante** (`Status.ExitCode()` de `StatusRetained` devuelve `0` en vez de `4`, para probar que el "nace verde" de la capa CLI de verdad se ejercita) | `go test ./internal/reports/... ./cmd/active-stack/... -run Retained -count=1 -v` | 2026-10-23T13:11:31-03:00 | `TestRunReportsCaptureDispatch_SecretInReport_ReturnsRetainedExitCodeAndNeverLeaksIt` FAIL por aserción (`code = 0, quería 4`); el test de `internal/reports` no lo detectó porque no mira `ExitCode()` — ver Aprendizaje. Revertido antes de continuar |

## Primer uso real (binario real, repo git temporal, 2026-10-23 13:11–13:13)

Instrucción explícita de este paso, por la lección del paso 1 ("tests verdes
no acreditan que el comando ande"): compilé el binario a un path temporal y
corrí `reports capture` de punta a punta, no sólo `go test`.

```
$ go build -o .../active-stack ./cmd/active-stack && echo "build OK"
build OK
```

En un repo git temporal **nuevo** (no el de los tests), con la copia real
del `.gitignore` del proyecto:

```
$ active-stack reports capture --change demo-change --phase verify --round 1 --role backend \
    --dispatch-id d-1 --sha abc123 --report-text "Informe de revisión: todo OK."
estado: COMPLETE (código 0)
clave: opsx/demo-change/verify/round-1/backend
directorio: openspec/changes/demo-change/informes/verify/round-1/backend
exit=0

$ active-stack reports capture --change demo-change --phase verify --round 1 --role backend \
    --dispatch-id d-1 --sha abc123 --report-text "Informe de revisión: todo OK."   # mismo evento
estado: COMPLETE (código 0)        # no-op: ni mtime de manifest.json cambió
exit=0

$ active-stack reports capture --change demo-change --phase verify --round 1 --role qa \
    --dispatch-id d-2 --sha abc123 \
    --report-text "Informe QA. Credencial de prueba: ghp_FAKEFAKEFAKEFAKEFAKEFAKEFAKEFAKEFAKE"
estado: RETAINED (código 4)
clave: opsx/demo-change/verify/round-1/qa
directorio: openspec/changes/demo-change/informes/verify/round-1/qa
causa: escaneo de secretos detectó: github-token
original retenido (local, ignorado por git): openspec/changes/demo-change/informes/verify/round-1/qa/.retenidos/0ff45a9d646a.txt
exit=4
```

Verificación del caso fijado por el propose (token falso → RETENIDO, sin
fuga), a nivel **git**, no sólo de la salida del comando — un `grep -r`
ingenuo con `ugrep` no es prueba suficiente porque esa herramienta respeta
`.gitignore` por default y no busca en archivos ignorados, lo cual hubiera
hecho pasar el chequeo por la razón equivocada:

```
$ git add -A && git status --porcelain
A  openspec/changes/demo-change/informes/LEEME.md
A  openspec/changes/demo-change/informes/verify/round-1/backend/informe.md
A  openspec/changes/demo-change/informes/verify/round-1/backend/manifest.json
A  openspec/changes/demo-change/informes/verify/round-1/backend/raw.txt
A  openspec/changes/demo-change/informes/verify/round-1/qa/manifest.json
# el .retenidos/ de qa NO aparece: git add -A ya lo excluyó

$ git diff --cached | grep -c "ghp_FAKEFAKE"
0

$ git ls-files -z | xargs -0 grep -l "ghp_FAKEFAKE"
# 0 coincidencias en todo lo que git versionaría

$ git check-ignore -v openspec/changes/demo-change/informes/verify/round-1/qa/.retenidos/0ff45a9d646a.txt
.gitignore:53:**/informes/**/.retenidos/	openspec/changes/.../qa/.retenidos/0ff45a9d646a.txt

$ git commit -qm "demo: captura de informes" && echo committed
committed
```

El token sólo vive en el archivo local `.retenidos/0ff45a9d646a.txt`
(confirmado con `cat`), que `git add -A` nunca propuso agregar. `manifest.json`
y `LEEME.md` — los dos archivos versionados que describen esta captura — no
contienen el token (confirmado leyéndolos enteros).

Otros dos casos reales, mismo binario:

```
$ active-stack reports capture ... --report /tmp/informe-demo.txt   # --report (archivo) en vez de --report-text
estado: COMPLETE (código 0)

$ active-stack reports capture ... --contract claude-hook-v1 --report-text "evento de hook sin adapter"
estado: UNSUPPORTED (código 3)
causa: contrato "claude-hook-v1" no está registrado: ningún adapter lo resolvió todavía
```

**Hallazgo de uso real** (no cubierto por los tests unitarios, que siempre
pasan `RepoDir` explícito a la raíz): correr el binario desde un
**subdirectorio** del repo falla, visible y sin escribir nada — es la
consecuencia correcta de "RepoDir tiene que ser la raíz canónica", pero el
CLI no la resuelve sola subiendo directorios:

```
$ cd demo-repo/sub/dir && active-stack reports capture ...
error: RepoDir ".../demo-repo/sub/dir" no es la raíz del worktree canónico
(raíz real: ".../demo-repo"); pasá la raíz, no un subdirectorio ni un worktree ajeno
exit=1
```

Es un comportamiento defendible (rechazo visible, no un error silencioso) y
queda igual que `approval`/`roles` (ninguno resuelve el root tampoco), pero
registro la deuda abajo porque el orquestador OPSX real puede no estar
siempre parado en la raíz.

## Verde final

```
$ git rev-parse HEAD
ac7b3dda6824a64ced9c6cccd7a4b03da40d6210
$ gofmt -l .
(sin salida, exit 0)
$ go vet ./...
(sin salida, exit 0)
$ go test ./... -count=1
ok  .../assets                                  0.048s
ok  .../cmd/active-stack                        1.431s
ok  .../cmd/active-stack/headless                3.126s
ok  .../e2e/starter                              1.158s
ok  .../internal/agents                          0.035s
ok  .../internal/agents/claude                   0.018s
ok  .../internal/agents/codex                    0.025s
ok  .../internal/agents/opencode                 0.019s
ok  .../internal/approval                        1.523s
ok  .../internal/backup                          0.205s
ok  .../internal/catalog                         0.391s
?   .../internal/contracttest           [no test files]
ok  .../internal/filemerge                       0.032s
ok  .../internal/harness/command                 0.048s
ok  .../internal/harness/config                 28.522s
ok  .../internal/harness/config/permissions      0.169s
ok  .../internal/harness/external                0.053s
ok  .../internal/harness/skill                   0.021s
ok  .../internal/i18n                            0.006s
ok  .../internal/install                         0.261s
ok  .../internal/model                           0.006s
ok  .../internal/pipeline                        0.006s
ok  .../internal/planner                         0.005s
ok  .../internal/reports                         0.457s
ok  .../internal/roles                           0.095s
ok  .../internal/system                          0.024s
ok  .../internal/tui                             0.032s
ok  .../internal/uninstall                       0.025s
ok  .../internal/verify                          0.038s
```
2026-10-23T13:13:00-03:00 a 13:13:32-03:00, sobre `ac7b3dd`. 28 paquetes, 0
fallas. Baseline anterior (informe del paso 1, mismo comando,
2026-10-23T12:53:09 a 12:53:38): 27 paquetes ok. La única diferencia es
`internal/reports`, que no existía — ningún paquete preexistente se rompió.

## Aprendizaje (para quien reuse el gate o los tests)

1. Mismo defecto que ya documentó el paso 1 con `-run Approval`: `-run
   'Capture|Worktree|Atomic|Recovery'` es un filtro de **nombre de test**,
   no de paquete ni de "qué prueba qué". Mis tests de triangulación de
   `BuildKey`/`scanSecrets`/`RegisterContract` no contienen ninguna de esas
   cuatro palabras y el gate exacto del propose los saltea en silencio;
   corren igual en `go test ./internal/reports/... -count=1` (sin `-run`),
   que es el comando que de verdad hay que correr para ver todo.
2. `git check-ignore -v` devuelve **exit 1** cuando el path NO está
   ignorado — eso es información, no un fallo del comando. Un helper de
   test que trata cualquier exit no-cero como `t.Fatalf` (el patrón que ya
   existía en `approval_dispatch_test.go`, reusado tal cual en
   `testutil_test.go`) revienta exactamente el caso que había que afirmar.
   Hubo que escribir un segundo helper (`gitCheckIgnore`) que distingue
   "ignorado" (exit 0) de "no ignorado" (exit 1) de un error real
   (cualquier otro código).
3. Herramientas de grep "inteligentes" (`ugrep`, y `rg` tiene el mismo
   comportamiento) **respetan `.gitignore` por default** y no entran a
   archivos ignorados. Un `grep -r TOKEN .` que no encuentra nada puede
   significar "no está" o puede significar "está, pero en un archivo que
   esta herramienta decidió no mirar". La prueba que de verdad importa es a
   nivel git (`git add -A && git status`, `git ls-files | xargs grep`,
   `git diff --cached`), no un grep de shell sobre el árbol de trabajo.

## JUICIO

- La forma de la clave para fases que no son `verify` (`opsx/<change>/<fase>/
  round-<n>/<rol>`) — autor: yo (constructor). Confianza: media. El propose
  fija literalmente `opsx/{change}/verify/round-{n}/{role}` (RF1-6) sólo para
  verify; generalicé esa misma forma a cualquier fase porque no hay una
  alternativa en el texto. `TestBuildKey_MatchesProposeLiteralFormForVerify`
  cubre la forma literal exacta; el resto es extrapolación, documentada en
  el propio código (`key.go`, comentario de `BuildKey`).
- Nombrar siempre `raw.txt` + `informe.md` con esos nombres fijos, sin
  negociar tipo de contenido (markdown vs texto plano según el origen) —
  autor: yo. Confianza: media-alta para el modo rápido (el propose no pide
  negociación de formato), baja si en el futuro un adapter necesita
  preservar, por ejemplo, un `.json` nativo del informe.
- `MaxContentBytes = 8 MiB` — autor: yo, valor arbitrario. El propose exige
  explícitamente que el límite esté "fijado antes de habilitar captura" sin
  decir cuál; 8 MiB evita que esto quede sin ningún límite (que sí sería
  truncamiento-por-omisión-de-diseño) pero no es una decisión de producto.
  Confianza: baja en el valor exacto, alta en que "algún límite explícito,
  con error visible al excederlo" es mejor que ninguno.
- Decidí que `Round <= 0` cuenta como "identidad ambigua" (`PENDING`), no
  como error de validación de ruta — a diferencia de `Change`/`Phase`/`Role`,
  que si tienen forma inválida SÍ son error. Razón: un `Round` ausente es un
  dato de negocio incompleto (el llamador todavía no sabe en qué vuelta
  está), no un intento de escapar del directorio. Confianza: media.
- El índice `LEEME.md` se **regenera completo** en cada captura (recorre
  todos los `manifest.json` del change) en vez de aplicar un append
  incremental. Decisión deliberada para no duplicar/desordenar entre
  distintas identidades concurrentes; NO resuelve dos escritores pisándose
  el `LEEME.md` al mismo tiempo entre sí (sin lock de archivo) — ver deuda
  D-9. Confianza: alta en que regenerar es más simple y menos propenso a
  corrupción que un append; media en que sea suficiente bajo concurrencia
  real.

## SUPUESTO (no verificado)

- Asumí que `Change` es el **slug del directorio** bajo
  `openspec/changes/<change>/` (ej. `"metodologia-fase-2"`), no una ruta a
  un archivo de propuesta como en `approval` (ahí `--change` es una ruta de
  archivo). Son convenciones distintas a propósito: `approval` lee un
  archivo de propuesta en `docs/propuestas/`, `reports` escribe en la
  estructura `openspec/changes/` que ya usa este mismo repo
  (`openspec/changes/metodologia-fase-2/`, `instalacion-sana/`, etc., con su
  propio `informes/LEEME.md` ya existente, hecho a mano). No verifiqué esto
  contra un documento que lo fije explícitamente para `reports`; lo inferí
  de que el propose dice "se versionan en `openspec/changes/<change>/
  informes/`" y de la estructura real de este repo.
- No probé `reports capture` con dos procesos escribiendo el MISMO
  `manifest.json` al mismo tiempo (concurrencia real, no secuencial). Los
  tests de "dos escritores" que corrí son secuenciales (llamada A, después
  llamada B) — cubren no-op/conflicto correctamente, pero no una carrera de
  verdad entre dos `os.Rename` simultáneos sobre el mismo path.
- No implementé timeout ni retención (el propose los menciona junto con el
  límite de tamaño como "TBD, deben quedar fijados antes de habilitar
  captura"). Sólo fijé el límite de tamaño; timeout y retención no existen
  en este código en ninguna forma, ni siquiera como placeholder.
- No verifiqué el comportamiento de `resolveWithinRoot` con symlinks en
  sistemas de archivos que no los soportan (Windows sin privilegios) más
  allá del `t.Skipf` que el test ya tiene para symlinks no disponibles en
  este entorno.
- Asumí que "SHA/diff revisado" (identidad) sólo necesita guardarse como
  dato (`manifest.json`), sin validar contra el repo que ese SHA exista o
  que el diff resuelva — a diferencia de `approval`, que sí resuelve commits
  contra git. `reports` no corre ningún `git rev-parse` sobre `--sha`/
  `--diff`: los toma literales. Si el propose esperaba esa validación
  cruzada, no está hecha.

## Deuda nueva (modo rápido)

- **D-5** — El CLI no resuelve la raíz del worktree solo cuando se invoca
  desde un subdirectorio: falla visible con "no es la raíz canónica" en vez
  de subir directorios hasta encontrarla (ver "Primer uso real"). Mismo
  patrón que `approval`/`roles`, que tampoco lo hacen.
- **D-6** — `MaxContentBytes` (8 MiB) es un valor arbitrario de este
  constructor, no una decisión de producto; igual falta fijar timeout y
  retención, que el propose pide explícitamente junto con el tamaño y que
  no existen en ningún lado de este código.
- **D-7** — `LEEME.md` se regenera completo por change en cada captura, sin
  lock de archivo: dos capturas concurrentes de identidades DISTINTAS del
  mismo change pueden pisarse la regeneración del índice entre sí (los
  `manifest.json` individuales sí quedan correctos, el índice puede quedar
  con una carrera perdida hasta la próxima captura que lo vuelva a generar
  completo).
- **D-8** — `reports capture` no valida `--sha`/`--diff` contra el repo
  (no corre `git rev-parse`); los toma como texto literal. `approval` sí
  resuelve commits. Si el contraste real con git es parte del contrato de
  "SHA/diff revisado", falta.
- **D-9** — La forma de clave para fases distintas de `verify` es
  extrapolación de este constructor (ver JUICIO), no texto del propose.
- **D-10** — Los patrones de secretos son propios y más angostos que una
  herramienta dedicada (ver EVIDENCIA): no cubren secretos de alta entropía
  sin prefijo reconocible (un token pegado sin ningún `key=`/`token=`/
  prefijo tipo `ghp_` alrededor no se detecta). gitleaks no está instalado
  en este entorno y no se instaló (restricción explícita de esta tarea).
- Hereda, sin resolver, **BE-4** (del propose): los adapters 2b/2c/2d no
  están implementados; sólo existe la interfaz (`Source`, `RegisterContract`).

## EVIDENCIA

- Build del seam: `go build ./internal/reports/...` exit 0, 2026-10-23
  13:06:24, antes de escribir ningún test.
- Rojo por aserción (fixture, no producción) y verde documentados arriba con
  hora exacta.
- Mutante aplicado y revertido sobre `Status.ExitCode()` (`StatusRetained`
  `4`→`0`): detectado por `TestRunReportsCaptureDispatch_
  SecretInReport_ReturnsRetainedExitCodeAndNeverLeaksIt` (CLI), NO detectado
  por los tests de `internal/reports` (que no afirman `ExitCode()`
  directamente) — confirma que esa capa de tests sí ejercita el camino real
  y no es un verde de fachada.
- `which gitleaks` → no encontrado; el escaneo de secretos es 100% patrones
  propios, con constancia en el código (`secrets.go`) y en este informe.
- `git check-ignore -v` sobre el retenido real → `.gitignore:53:**/informes/
  **/.retenidos/` (la regla nueva matchea); sobre un original activo → exit
  1 (no ignorado); tras `git add -A`, el original activo aparece en
  `git ls-files` y el retenido NO aparece en `git status --porcelain`.
- `git diff --cached | grep -c TOKEN` → `0` y `git ls-files -z | xargs -0
  grep -l TOKEN` → sin coincidencias, sobre un repo real donde el token SÍ
  existe en disco (confirmado con `cat` del archivo `.retenidos/`).
- Suite completa (`go test ./... -count=1`) verde sobre `ac7b3dd`, con hora
  y comparación contra el baseline del paso 1 (27→28 paquetes, mismo
  resto).

## Roles no convocados en este paso (JUICIO, confianza alta)

Este informe es de **apply**, no de **verify**: por el propio contrato del
propose ("nadie convoca arquitecto/QA/auditor en este trabajo; eso es la
revisión prevista en verify, por paso"), no corresponde convocar revisores
acá. Dejo constancia igual, como pide el bloque de gobierno de roles:
backend/DBA/frontend/DevOps no se convocaron porque la fase actual es
`apply`, que convoca únicamente al implementador; QA/arquitecto/backend
como **revisores** quedan para la vuelta de `verify` de este paso, que este
informe no sustituye.
