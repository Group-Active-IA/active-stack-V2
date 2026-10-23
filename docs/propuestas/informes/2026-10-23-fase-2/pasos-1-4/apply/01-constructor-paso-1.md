# Constructor — paso 1 (aprobación verificada, sin firma criptográfica)

Fecha: 2026-10-23 · Rama: `feat/metodologia-fase-2` · Modo rápido.

Implementa **únicamente** la sección `## 1 — Aprobación verificada, sin firma
criptográfica (J1 revisada)` de
`docs/propuestas/2026-10-23-fase-2-pasos-1-4.md`, aprobada por Juani en el
commit `d26131e`.

## Commits

| SHA | Mensaje |
|---|---|
| `50de0008e38ec0173899254e278f45a38eb092e4` | `feat(approval): lectura de aprobación escrita sin firma criptográfica` |
| (este informe) | `docs(propuestas): informe del constructor del paso 1` |

Archivos: `internal/approval/{approval,git,table,manifest,render,verify}.go` +
sus `_test.go`, `cmd/active-stack/approval_dispatch.go` (+ test), y el wiring
del subcomando `approval` en `cmd/active-stack/main.go`.

## Qué construí

`active-stack approval --change <ruta> --commit <C>`, sólo lectura. Resuelve
C, exige que el árbol documental esté limpio, que C modifique **sólo**
`--change` (ni "sin tocarlo" ni "mezclado con código"), parsea la última fila
de la tabla `## Aprobación` en C y en el padre de C (para exigir que la fila
no exista YA completa antes — "commit documental separado"), lee en el
mensaje de C un trailer `Propuesta-revisada: <SHA>` como referencia a P,
verifica que P resuelva y sea ancestro del padre de C, construye el
manifiesto determinista (propuesta + informes que la propuesta referencia
vía `(informes/....md)`) hasheado en P, y por último compara el blob de la
propuesta en C contra `--current-ref` (default HEAD) para declarar STALE si
cambió después de la aprobación. Devuelve `WRITTEN_UNVERIFIED` (código 2,
mensaje literal `"estructura válida; identidad humana no acreditada"`) o
`MISSING`/`INVALID`/`STALE` (código 1). No existe código 0 ni un estado
`VERIFIED`.

## Decisión de diseño no explícita en el propose (JUICIO)

El propose dice "la decisión referencia el SHA completo P" pero la tabla de
Aprobación (columnas Elección/Persona/Fecha/Alcance) ya está firmada y no
tiene lugar para P ni para un manifiesto. Decidí que esa referencia vive en
el **mensaje del commit C** como trailer `Propuesta-revisada: <SHA>`, no en
la tabla. Alternativa descartada: agregar una columna a la tabla — hubiera
reabierto el formato que Juani ya firmó, que el propose prohíbe ("no
reabrir sus funcionalidades"). Confianza: media. Esto es una interpretación
mía de una frase ambigua, no algo que el propose deletree; si Juani prefiere
otra convención, es un cambio de `extractProposalRef` en
`internal/approval/verify.go`, no un rediseño.

Consecuencia verificada: la aprobación real `d26131e` **no** cumple este
contrato (su mensaje no tiene el trailer), porque es anterior a que este
contrato existiera. No se intentó ni se debía hacer pasar retroactivamente.

## Rojo / verde por caso

| Caso | Comando | Hora | Resultado |
|---|---|---|---|
| Seam compila (sin lógica) | `go build ./internal/approval/...` | 12:46:3x | exit 0 |
| **RED** caso feliz | `go test ./internal/approval/... -run TestVerify_ValidWrittenApproval_ReturnsWrittenUnverified -v` | 12:47:03 | `Status = MISSING, quería WRITTEN_UNVERIFIED... Cause="sin implementar"` — FAIL por aserción, no por compilación |
| GREEN intermedio (bug real) | mismo comando | 12:47:39 | `Status = INVALID ... "el commit de aprobación no tiene padre"` — bug en `commitHasParent` (`^@` no es sintaxis válida para `rev-parse --verify`); corregido a `^1` |
| **GREEN** caso feliz | mismo comando | 12:47:52 | PASS |
| Triangulación (13 casos nuevos: tabla vacía/ausente, sin C, mezclado con código, persona/fecha/alcance faltantes, tabla ya completa antes, P sin referenciar, P ajeno, informe faltante, STALE, SHA ambiguo, cero escrituras, cero-valor no aprobado, sin VERIFIED) | `go test ./internal/approval/... -v` | 12:49:16 | 1 FAIL: `TestVerify_ProposalNotAncestor_ReturnsInvalid` — fixture reutilizaba una propuesta ya aprobada y disparaba la regla de "tabla ya completa" antes de llegar a la de ancestro |
| Fix del fixture (no de producción) + re-run | mismo comando | 12:49:37 | PASS, 17 tests + 3 subtests |
| Parser de tabla (3 casos: última fila gana, sin heading, corta en el siguiente heading) | `go test ./internal/approval/... -run TestParseApprovalTable -v` | 12:49:53 | PASS directo (caracterización contra el parser ya implementado) |
| CLI end-to-end (5 casos: válido texto, válido JSON, flags faltantes, --help, tabla ausente) | `go test ./cmd/active-stack/... -run Approval -v` | 12:51:15 | PASS directo |
| **Gate exacto del propose** (post-rename de tests para que `-run Approval` los alcance a todos — ver Aprendizaje) | `go test ./internal/approval ./cmd/active-stack -run Approval -count=1 -v` | 12:51:58 | PASS — 20 tests (+3 subtests) en `internal/approval`, 5 en `cmd/active-stack` |

## Verde final

```
$ git rev-parse HEAD
50de0008e38ec0173899254e278f45a38eb092e4
$ go vet ./...
(sin salida, exit 0)
$ go test ./... -count=1
ok  .../assets                                  0.039s
ok  .../cmd/active-stack                        0.856s
ok  .../cmd/active-stack/headless                2.582s
ok  .../e2e/starter                              0.806s
ok  .../internal/agents                          0.016s
ok  .../internal/agents/claude                   0.009s
ok  .../internal/agents/codex                    0.009s
ok  .../internal/agents/opencode                 0.008s
ok  .../internal/approval                        1.052s
ok  .../internal/backup                          0.155s
ok  .../internal/catalog                         0.204s
?   .../internal/contracttest           [no test files]
ok  .../internal/filemerge                       0.021s
ok  .../internal/harness/command                 0.021s
ok  .../internal/harness/config                 27.862s
ok  .../internal/harness/config/permissions      0.125s
ok  .../internal/harness/external                0.088s
ok  .../internal/harness/skill                   0.026s
ok  .../internal/i18n                            0.018s
ok  .../internal/install                         0.261s
ok  .../internal/model                           0.011s
ok  .../internal/pipeline                        0.007s
ok  .../internal/planner                         0.008s
ok  .../internal/roles                           0.100s
ok  .../internal/system                          0.030s
ok  .../internal/tui                             0.033s
ok  .../internal/uninstall                       0.023s
ok  .../internal/verify                          0.041s
```
2026-10-23T12:53:09-03:00 a 12:53:38-03:00. 27 paquetes, 0 fallas, 0 paquetes
nuevos rotos. Baseline previo al cambio (mismo comando, 2026-10-23T12:44:04 a
12:44:xx): 26 paquetes ok, idéntica lista salvo que `internal/approval` no
existía — ningún test preexistente se rompió.

## Aprendizaje (para quien reuse el gate)

El comando de rojo-antes del propose, `-run Approval`, es un filtro de
**nombre de test por regex**, no de paquete. Mis primeros nombres
(`TestVerify_ProposalNotAncestor_...`, `TestZeroValueResult_IsNotApproved`,
etc.) no contienen el string "Approval" y ese comando los saltea en
silencio — corrí la batería completa primero y vi que con `-run Approval`
sólo entraban 7 de 20 tests. Renombré todos con el prefijo `TestApproval*`
para que el gate exacto del propose ejecute la batería completa. Si a futuro
se agregan tests a este paquete, tienen que incluir "Approval" en el nombre
o este mismo gate los va a ignorar sin avisar.

## EVIDENCIA

- Build del seam: `go build ./internal/approval/...` exit 0, antes de escribir
  ningún test (2026-10-23, ~12:46).
- Rojo por aserción (no por compilación) documentado arriba, 12:47:03.
- Verde final documentado arriba, con SHA y hora: `go test ./... -count=1` a
  las 2026-10-23T12:53:09/38-03:00 sobre el commit `50de000`.
- `go vet ./...` sin salida, exit 0, mismo momento.
- `TestApprovalVerify_MakesNoWrites` compara una huella (`git status
  --porcelain=v1 -uall` + `rev-parse HEAD` + listado de archivos no-`.git`)
  antes y después de `Verify` sobre un repo con aprobación válida: idénticas.
  Es la evidencia operacional de "approval nunca escribe", no sólo una
  revisión de que el código no llama a `os.WriteFile`/`git commit`.

## JUICIO

- La convención del trailer `Propuesta-revisada:` (arriba) — autor: yo
  (constructor). Confianza: media, es una interpretación de una frase
  ambigua del propose, no algo verificado contra un tercer documento.
- El manifiesto sólo cubre la propuesta y los informes que ella misma
  referencia con un link `(informes/....md)` relativo a su propio
  directorio. Comparado contra cómo la propuesta real
  (`2026-10-23-fase-2-pasos-1-4.md`) cita sus informes (exactamente ese
  patrón, p. ej. `informes/2026-10-23-fase-2/pasos-1-4/01-arquitecto.md`).
  Confianza: alta para esta convención real; baja para cualquier otra forma
  de citar un informe (ruta absoluta, sin paréntesis, etc.) — esos casos
  simplemente no entran al manifiesto y no se detecta su ausencia.
- `isDirty` bloquea si CUALQUIER cosa en el repo está sin commitear, no sólo
  la propuesta. Es más estricto de lo que el propose pide literalmente
  ("árbol documental sucio") pero evitar acotar el scope a un subárbol es lo
  mínimo defendible sin un segundo flag de alcance. Confianza: media.

## SUPUESTO (no verificado)

- Asumí que `--change <id>` es la **ruta al archivo de la propuesta**
  relativa al repo (p. ej. `docs/propuestas/2026-10-23-fase-2-pasos-1-4.md`),
  no un slug que se resuelve a esa ruta por convención. El propose no fija
  ningún mecanismo de resolución de `<id>` y el repo no tiene ninguno hoy
  (verificado: `rg -n "\-\-change" --type go` antes de este cambio no
  encontró nada). Si en el futuro se decide un esquema de slugs, este CLI
  necesita un paso de resolución antes de llamar a `approval.Verify`.
- Asumí que el commit C puede no tener padre (commit raíz) y en ese caso
  bloqueo con INVALID ("P no puede ser su ancestro de la nada"). No hay
  fixture de un repo real con ese caso límite fuera del test unitario; es
  coherente con la lógica pero no lo vi fallar en un escenario de uso real.
- No probé el comportamiento con un `--current-ref` explícito distinto de
  HEAD (el flag existe y se parsea, pero ningún test lo ejercita con un
  valor no-default). Queda sin triangular.
- No hice nada con los demás pasos del propose (2a–2d, 3, 4): fuera de
  alcance de este informe, tal como se me pidió.

## Deuda nueva (modo rápido)

- **D-1** — `--change <id>` no resuelve slugs, sólo rutas literales (ver
  SUPUESTO). Si OPSX necesita convocar `approval` por nombre de change en vez
  de por ruta de archivo, falta ese mapeo.
- **D-2** — El trailer `Propuesta-revisada:` es una convención nueva que
  ningún commit real usa todavía (incluido `d26131e`, anterior a este
  contrato). Falta decidir si se re-aprueba con el trailer o si se acepta un
  mecanismo alternativo para aprobaciones ya firmadas antes de este cambio.
- **D-3** — `isDirty` mira el repo completo, no sólo la ruta de la propuesta;
  un cambio sin commitear en un archivo no relacionado bloquea la
  validación. Acotar el chequeo a la ruta (y a los informes referenciados)
  es una mejora de precisión, no de corrección.
- **D-4** — Ningún test cubre `--current-ref` con un valor explícito ni un
  historial con `rebase` real (el propose menciona "rebase" como causa de
  bloqueo; hoy sólo se prueba vía SHA ambiguo/STALE, no un rebase real
  simulado con `git rebase`).

## Primer uso real (orquestador, 23/10/2026 12:55 a 12:58)

EVIDENCIA. Con la suite en verde sobre `6642a62`, el comando real fallaba siempre:

```
$ active-stack approval --change docs/propuestas/2026-10-23-fase-2-pasos-1-4.md --commit d26131e
error: approval.Verify: RepoDir, ProposalPath y Commit son obligatorios   (exit 1)
```

`main` llamaba `runApprovalDispatch(..., "", ...)` y ningún test cubría el directorio vacío. Rojo por aserción en `TestRunApprovalDispatch_EmptyDirUsesWorkingDirectory` (12:55:52, código 1 en vez de 2); arreglo en `53d4d4e`.

Segundo hallazgo de uso real (era la deuda D-3): `isDirty` usaba `-uall` y cualquier `.log` sin trackear volvía INVALID una firma válida. Rojo en `TestApprovalVerify_UntrackedUnrelatedFile_StillWrittenUnverified` (12:56:31); arreglo en `1edd7f7` (`--untracked-files=no`); `TestApprovalVerify_TrackedProposalModified_ReturnsInvalid` cubre que un archivo trackeado modificado sigue bloqueando. `gofmt` además reformateó `render.go`.

Estado final sobre la firma real:

```
$ active-stack approval --change docs/propuestas/2026-10-23-fase-2-pasos-1-4.md --commit d26131e
estado: INVALID (código 1)
causa: el commit de aprobación no referencia el SHA completo (P) de la propuesta revisada
```

JUICIO (orquestador, confianza alta): diagnóstico correcto. `d26131e` es anterior al contrato (D-2). No se reescribe el commit de Juani; los scripts de firma siguientes agregan `Propuesta-revisada: <sha completo de P>`.

Lección para el método: el constructor probó la función y el dispatch con un dir explícito, nunca el binario como lo llama `main`. Tests verdes no acreditan que el comando ande; correrlo una vez sí.
