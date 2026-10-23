# Exploración: sumar a Active Stack V2 la metodología de PropTech

Fecha: 2026-10-23 · Rama: `feat/metodologia-proptech` · Fase: explore (no se implementa nada).

Informes de los roles convocados, guardados textuales (extraídos del transcript del sub-agente con
`jq`, sin pasar por el orquestador):

- `docs/exploraciones/informes/2026-10-23-analista-funcional.md` (role-analyst, opus)
- `docs/exploraciones/informes/2026-10-23-arquitecto.md` (role-architect, opus)

Los dos se convocaron de a uno, con prompts independientes; el arquitecto declara que no abrió el
informe del analista.

> **Actualizado el 23/10 con las decisiones de Juani (§8).** Lo que esas decisiones reemplazan quedó
> en su lugar, marcado con *Reemplazado (§8)*, para que se vea qué se pensó y por qué cambió. La
> propuesta que sale de acá está en `docs/propuestas/2026-10-23-metodologia-proptech-en-opsx.md`.

---

## 0. Tres hallazgos que van antes que todo

**EVIDENCIA** (orquestador, 2026-10-23 09:05 a 09:13)

```
$ go version -m ~/.local/bin/active-stack
  path  github.com/Group-Active-IA/active-stack/cmd/active-stack/headless.test
$ active-stack roles --files internal/auth/login.go
  flag provided but not defined: -project
  Usage of uninstall: ...
$ rg -c "active-stack roles" ~/.claude/CLAUDE.md        → 0
$ eza -l ~/.claude/CLAUDE.md                           → 9 sep 16:09
$ git log -1 --format='%h %aI' e513aa7                 → 2026-09-14 (el bloque nombra el comando desde acá)
$ rg -n "role-implementer" -g '!*.golden' .            → nada en el repo; los role-*.md viven sólo en ~/.claude/agents (9 sep)
```

1. **El binario instalado en `~/.local/bin` es un binario de tests** (`headless.test`), no
   `cmd/active-stack`. Por eso `active-stack roles` no funciona en esta máquina.
2. **El bloque instalado es del 9/09**, anterior al subcomando `roles` (14/09). El orquestador de
   PropTech nunca leyó la instrucción de correrlo.
3. **Active Stack no instala los agentes `role-*`** que su propio bloque y su clasificador
   (`internal/roles/classify.go:37-43`) nombran. Esos archivos no tienen dueño ni versión.

**JUICIO** (orquestador, confianza alta): el problema 5 (clasificador sin usar) es primero de
instalación y recién después de cumplimiento. Y el problema 1 (role-implementer) es un archivo que
Active Stack ni siquiera controla.

**SUPUESTO**: no sé qué escribió el binario de tests en `~/.local/bin` (el arquitecto sospecha un test
de `cmd/active-stack/headless` con self-install sobre el HOME real; no verificado).

---

## 1. La metodología de PropTech, como reglas

Fuentes: `PropTech/knowledge-base/14_forma_de_trabajo.md`, `15_lecciones.md`, `CLAUDE.md`, el change
archivado `2026-10-22-login-y-usuarios`, `.git/informes-2b/` y `.git/informes-cimientos/`, y las
memorias `opsx-secuencial-con-auditor`, `codex-cli-en-paralelo` y `proptech-3b-en-curso`.

| # | Regla | Lo que la respalda |
|---|---|---|
| R1 | El orquestador no escribe código: arma el pedido, lanza **de a un agente**, lee, decide y le cuenta a la persona. Ningún agente lanza otros. | 14 "Quiénes" y "Reglas del proceso"; memoria `opsx-secuencial-con-auditor` ("cuidemos los tokens", 09/10) |
| R2 | Ciclo **constructor → rompedor → auditor**. El rompedor deja cada falla como test que falla; **el cambio está aprobado cuando el auditor da OK**. | 14 "El ciclo de cada cambio" |
| R3 | **Tope de 3 vueltas**: a la tercera el rompedor no vuelve al constructor; su informe con diagnóstico va al auditor, que separa bloqueante de deuda, y el veredicto llega **textual** a la persona, que decide. | D-092; ajuste del 14/10 en la memoria (dos vueltas del 1c en loop sin auditor) |
| R4 | Escalado por riesgo: plata o peso legal → ciclo completo; chico o sólo lectura → constructor + una pasada que rompe y audita. | 14 "Escalado por riesgo" |
| R5 | **Constructor de otra familia de modelo** que el rompedor y el auditor (Codex construye, Claude rompe y audita). | L-025 (18 hallazgos en la primera lectura externa del código de Codex); memoria `codex-cli-en-paralelo` |
| R6 | El brief del constructor **nombra los consumidores fuera de su área** (otros módulos, `infra/`, scripts). | L-023, L-025 |
| R7 | CRITICAL/HIGH: la propuesta la **firma la persona por escrito y con fecha** en la tabla `## Aprobación` del proposal, que ningún agente completa; queda en un commit propio. | `proposal.md:263-293` del 2b; commits `0e6cf9c`, `2212c9e`, `30c2f9b`, `2eb8271` |
| R8 | **Una regla nueva después de la firma lleva firma propia**; el constructor no infiere la aprobación del brief. | L-029 (RR-1 a RR-4, `610eb4a`); auditor del 2b, fila 1 (E17 sin firma, ALTA) |
| R9 | TDD con **rojo y verde pegados** en `tasks.md` bajo la tarea: comando, `date +%T`, últimas líneas. `[x]` sólo con los dos pegados. | `.git/informes-2b/reglas-2b.md`; `tasks.md` del 2b, tarea 2.1 |
| R10 | Todo test tiene que poder dar rojo: si nació verde, **mutante** (diff de 2-4 líneas, corrida, revertido) que falle **en la aserción que el test dice probar**, no en un paso anterior. Sin commit, copia previa y decirlo. | L-013, L-019, L-034, L-035 |
| R11 | "Verde" son **todos los pasos del workflow** (lint, formato, tipos, tests), y quien agrega un archivo los corre. Rompedor y auditor **cuentan** los tests ellos mismos. | L-003, L-005 |
| R12 | Las **suites de cierre van después del último commit**, con el SHA anotado; si el código cambia, la corrida vence aunque siga verde. | L-010, L-031; auditor del 2b, fila 6 (web sobre 9eb0976 con HEAD 2b570b3) |
| R13 | **Lecciones numeradas** L-NNN, sólo se agregan; constructor y rompedor las leen antes; el auditor cierra con una sección LECCIONES; el orquestador las agrega al archivar. | 14 "Lecciones"; `15_lecciones.md` (37 entradas) |
| R14 | El auditor lee los **informes originales**, nunca un resumen del orquestador. | D-092; auditor del 2b, fila 7 (de las revisiones del propose quedó sólo el resumen) |
| R15 | Cada vuelta de arreglos se cierra con **revisores de rol de sólo lectura sobre lo que tocó esa vuelta**, no con otra vuelta entera de rompedor. | L-024; auditor del 2b, fila 3 |
| R16 | Cuando el auditor dicta el caso exacto, la coincidencia con el rompedor no es convergencia: la verificación independiente sale del encuadre o lo declara. | L-004, L-021 |
| R17 | **Verificación real con la persona al teclado** antes de archivar cuando hay un proveedor externo de identidad; lo no recorrido queda como deuda declarada, nunca asumido. | 14 "Verificación real"; L-032 (`__clerk_ticket`, tras 2259 + 385 tests verdes) |
| R18 | **"Pusheado no es listo"**: después de un push a `main` se espera el run completo de CI sobre ese commit; no se apila un arreglo sobre otro sin confirmar. | L-037 (dos rojos de `main` el 22/10) |
| R19 | **Auditoría de cimientos** con revisores de rol de a uno + auditor antes de un núcleo crítico. | `.git/informes-cimientos/01..09`; L-035, L-036 salieron de ahí |
| R20 | Toda UI pasa por **impeccable/shape aprobado** antes, y se **mira en un navegador real** después de cada bloque de UI. | L-026 (rediseño D-105a sobre un bloque "terminado") |
| R21 | El brief separa **conflictos de orden o redacción** (los resuelve el constructor y lo declara como suyo) de **permisos** (frena). | L-030 |
| R22 | Los **informes viven fuera del repo**, en `.git/informes-<cambio>/`, y un `LEEME.md` indexa qué corrió quién. | `.git/informes-2b/LEEME.md` |
| R23 | El orquestador verifica en cada bloque lo que el constructor externo no controla: `docker ps` antes y después, atribución de cada commit, suite completa. | 14 "Constructor con Codex" |
| R24 | La deuda se asigna al **primer cambio que de verdad usa la pieza**. | L-012 |

**JUICIO** (orquestador, confianza alta): R1-R3, R7-R10, R12-R13, R17-R18 son el núcleo transferible.
R4, R6, R11, R16, R21 y R24 son criterio de brief y de auditor, no estructura. R23 es específico de
correr Codex sin sandbox y no va al instalador.

---

## 2. Pieza por pieza: qué cubre Active Stack, qué falta, qué cambiar

| Pieza | Hoy en Active Stack V2 (EVIDENCIA) | Falta | Cambio propuesto |
|---|---|---|---|
| Orquestador que no codea y delega | base `## Role`, `## Delegation Rules` | nada | ninguno |
| Secuencial, de a uno | base, Dispatch punto 2 (`sdd-orchestrator.md:314`): lo independiente sale **en un mismo mensaje** | el método pide secuencia | permitir modo secuencial declarado (ver §4, `build-break-audit`) |
| Constructor → rompedor → auditor con tope de vueltas | `review.md` (tres pasadas por lentes, sin vueltas); `auditor.md` (lee informes) | el loop, el rompedor que deja tests, el tope, la decisión humana a la tercera | toggle nuevo `build-break-audit` |
| Veredicto textual a la persona | `auditor.md` "Who it answers to" | nada | se cita, no se duplica |
| Constructor de otra familia | `model-routing` sólo tiene alias de Claude | proveedor por agente | invariante dentro de `build-break-audit`; nada en `model-routing` |
| Firma humana verificable | `governance.md`: "explicit human approval", sin artefacto | dónde vive la aprobación y cómo se verifica; `role-implementer` sin canal (`role-implementer.md:178-193`) | artefacto en `governance.md` + subcomando de verificación + firma criptográfica (§3, problema 1) |
| Firma de correcciones post-firma | nada | todo | regla en `governance.md` |
| Rojo/verde pegados, mutante | `strict-tdd.md:64-73`: tabla de tildes, sin hora, salida ni mutante | todo | modificar `strict-tdd.md` |
| Suites después del último commit con SHA | `ship.md` pide SHA sólo al publicar | SHA por corrida de cierre | regla en `strict-tdd.md` o en `build-break-audit` |
| Lecciones L-NNN | `knowledge-base.md` tiene reescribir/agregar, no lecciones | archivo, sección LECCIONES, lectura previa | tercera operación en `knowledge-base.md` |
| Informes originales fuera del repo | `roles.md:174-190` admite que no hay lista de herramientas; revisores sin Write | canal que no pase por el orquestador | hook `SubagentStop` + subcomando de respaldo (§3, problema 3) |
| Verificación real con persona al teclado | `review.md:14` sólo habla de proveedor que "devuelve basura" | todo | regla en `build-break-audit` (cierre) o en `review.md` |
| Esperar CI tras el push | `ship.md:170` "the cycle ends at the open PR"; ningún `gh run` | todo | modificar `ship.md` y `assets/commands/claude/opsx/ship.md` |
| Auditoría de cimientos | nada; el clasificador lee un diff | modo "árbol completo" | flag en `active-stack roles` (p. ej. `--all`), más una regla corta |
| UI: shape antes, mirar después | nada en los assets | todo | regla en `roles.md` (fila frontend) o skill aparte; ver pregunta 6 |
| Brief con criterio de conflictos | no hay plantilla de brief del constructor | todo | parte de `build-break-audit` |
| Clasificador de revisores | base, Dispatch punto 3; `active-stack roles` | que llegue instalado y funcione | chequeos en `verify` (§3, problema 5) |

> *Reemplazado (§8)*: las filas "Secuencial", "Constructor → rompedor → auditor", "Constructor de otra
> familia", "Brief con criterio de conflictos" y "Verificación real" ya no van a un toggle
> `build-break-audit`. Van dentro de `apply` (`strict-tdd.md`) y de `verify` (`review.md`,
> `auditor.md`, `senior.md`).

---

## 3. Los cinco problemas reales del 22/10 y su propuesta

### Problema 1: `role-implementer` se negó a escribir auth ya firmado

**EVIDENCIA**: `role-implementer.md:178-193` ("no existe el canal"; ningún archivo que un agente pueda
escribir cuenta como aprobación). El 2b se firmó en `0e6cf9c` y el constructor terminó siendo
`general-purpose` con el OK de Juani en el prompt (memoria `opsx-secuencial-con-auditor`).

**JUICIO** (analista y arquitecto por caminos distintos, confianza alta): **verificar "tabla de
Aprobación con commit de Juani" no alcanza.** En esta máquina los agentes commitean con el
`user.name` de Juani y Codex corre sin sandbox y commitea; el autor del commit no distingue a la
persona del agente. `role-implementer` tiene razón en negarse ante esa prueba.

**Propuesta**:
1. `governance.md` define el artefacto: sección `## Aprobación` en `proposal.md` con elección,
   persona, fecha y alcance; correcciones posteriores en una fila nueva (L-029).
2. Subcomando `active-stack approval --change <nombre>` que verifica que el commit que tocó esa
   sección **está firmado** (`git verify-commit`) con una clave que el agente no puede usar
   (passphrase o token físico) y que no tocó nada más.
3. `role-implementer` acepta como canal **sólo la salida de ese comando**, no el texto de la tabla.
4. Para que (3) tenga dueño, Active Stack tiene que **instalar los `role-*.md`** (tipo de harness nuevo
   `agent`, reusando el patrón de `internal/harness/command`). Hoy no los instala.

**SUPUESTO**: que el "script de firma" que corre Juani en PropTech no usa GPG (no se encontró el script;
el `LEEME` del 2b sólo dice "Firma de Juani E1-E15 por script (0e6cf9c)").

### Problema 2: `ship` termina en el push

**EVIDENCIA**: `ship.md:170`; ningún `gh run` en el fragmento; L-037.

**Propuesta**: "listo" pasa a ser **CI verde sobre el SHA pusheado** (`gh run list --commit <sha>` y
`gh run watch`), con tres salidas distintas: verde, rojo (se reporta y no se apila otro arreglo), y
"no hay CI configurado" o "no pude consultar" (nunca leído como verde). Se cambian las **dos copias**
de la doctrina: el fragmento y el comando `/opsx:ship` (3786 palabras), que hoy están duplicadas.
Cuidado con las frases prohibidas de `ship_toggle_test.go:554-627` ("waits for confirmation").

### Problema 3: los revisores no pueden escribir su informe

**EVIDENCIA**: `.git/informes-cimientos/01-arquitecto.md`: "Copiado textual por el orquestador: el rol no
tiene permiso de escritura" (9 informes así). Auditor del 2b, fila 7.

**JUICIO** (arquitecto, confianza media-baja; orquestador de acuerdo con la dirección): no dar Write
acotado por carpeta. Claude Code no restringe Write por ruta para un sub-agente, y reabre el agujero de
"el revisor edita". Mejor un canal que **no pase por ningún LLM**:

- **Primario**: hook `SubagentStop` que instala Active Stack y copia el último mensaje del sub-agente a
  `.git/informes-<cambio>/<rol>.md`.
- **Respaldo**: subcomando que extrae el último mensaje del transcript
  `~/.claude/projects/<proyecto>/<sesión>/subagents/agent-*.jsonl`.

Esta misma exploración usó el respaldo a mano (`jq` sobre el transcript) para guardar los dos informes
de §0: funciona, y el orquestador no tocó el texto.

**SUPUESTO**: que el payload de `SubagentStop` trae la ruta del transcript; que el formato jsonl es
estable (es interno de Claude Code y sólo existe en Claude). En PropTech no hay transcripts de
sub-agentes guardados (arquitecto, E11), así que no se puede reconstruir lo de octubre.

### Problema 4: el bloque es enorme y se carga siempre

**EVIDENCIA**: `wc -w` de la base de Claude + fragmentos: 17.311 palabras (la base sola, 4.297); el
golden `compose_all_toggles_claude.golden` 17.315; `~/.claude/CLAUDE.md` instalado 14.315.

**Propuesta** (arquitecto, confianza media): en `CLAUDE.md` queda la base más, **por toggle, un resumen
corto** con su "terminado" y las reglas que deciden antes de actuar (push a main pregunta, gate
humano, auditor textual). El detalle se mueve a canales que Active Stack **ya posee y se cargan al
invocarlos**: `/opsx:verify` (review, auditor, roles), `/opsx:ship`, y los `role-*.md` si pasan a
instalarse. Rechazada: una skill por toggle cargada bajo demanda (si no se carga, la regla no corre y
no queda señal). **Medir antes** qué secciones se usan de verdad. Es el cambio que más tests mueve:
va último.

### Problema 5: el clasificador existía y no se usó

**EVIDENCIA**: §0, puntos 1 y 2.

**Propuesta**: en `internal/verify`, (i) comparar el bloque instalado contra `Compose(toggles)` y avisar
si quedó viejo; (ii) chequear que el binario sea `cmd/active-stack` y no un `.test`, y que
`active-stack roles --help` salga con 0; (iii) subcomando `version`. Primero, encontrar qué test
escribe el binario de tests sobre el HOME real.

---

## 4. Propuesta de diseño

### Toggles nuevos

> *Reemplazado (§8)*: no hay toggle nuevo. Juani decidió que el método va dentro de las fases OPSX
> que ya existen. La tabla queda como registro de la alternativa rechazada.

| Toggle | Qué agrega | Posición | Por qué ahí |
|---|---|---|---|
| `build-break-audit` | Ciclo constructor → rompedor → auditor **de a uno**, tope de 3 vueltas con diagnóstico al auditor y decisión de la persona; escalado por riesgo; constructor de otra familia que rompedor y auditor (invariante, Codex como ejemplo); brief con consumidores fuera del área y criterio de conflictos (L-030); cierre de vuelta con revisores de rol sobre lo último (L-024); verificación real con la persona al teclado si hay proveedor externo. Declara que en este modo el punto 2 del Dispatch (paralelo) no aplica entre constructor, rompedor y auditor. | entre `review` y `knowledge-base` | depende de `auditor`, `roles` y `review` (anteriores) y los nombra sólo en forma condicional ("where the X fragment is on the page"); ningún fragmento anterior lo nombra a él, así no hay ciclo |

Nada más nuevo como toggle. El resto son cambios a lo existente.

### Cambios a lo existente

| Pieza | Cambio |
|---|---|
| `strict-tdd.md` | rojo y verde pegados en `tasks.md` con comando, hora y salida; mutante que falla en la aserción (L-019, L-034); copia previa si no hay commit (L-013); suites de cierre después del último commit con SHA (L-031); "verde" son todos los pasos del CI (L-005) |
| `ship.md` + `commands/claude/opsx/ship.md` | esperar el CI del SHA pusheado; tres salidas |
| `knowledge-base.md` | tercera operación: lecciones que sólo se agregan, con número, origen, regla y a quién aplica; "si el auditor está, su veredicto cierra con LECCIONES" (condicional) |
| `governance.md` | artefacto `## Aprobación`, correcciones con firma propia, y que el gate se verifica con `active-stack approval` |
| `roles.md` | corregir la sección "no tool list": con `subagent_type` los `role-*` sí restringen herramientas; nombrar el canal de informes; fila de UI con mirar en navegador |
| base, Dispatch punto 2 | excepción explícita cuando `build-break-audit` está activo |
| `internal/roles` + `roles_dispatch.go` | flag para clasificar el árbol completo (auditoría de cimientos) |
| subcomando nuevo `approval` | verifica la firma de la sección de aprobación |
| subcomando o hook de informes | `SubagentStop` + extracción del transcript |
| `internal/verify` | bloque instalado al día; binario correcto |
| catálogo | tipo de harness `agent` para instalar los `role-*.md` |

### Orden (el del arquitecto, con el que coincido)

> *Reemplazado en parte (§8)*: el paso 5 ya no es un toggle sino cambios a `strict-tdd.md`,
> `review.md`, `auditor.md` y `senior.md`; el orden nuevo está en la propuesta.

1. **Instalación sana** (problema 5): binario y bloque al día en `verify`. Sin esto, nada de lo demás llega.
2. **`fae8bbf`** de juancruz (ver §7): bug real y barato.
3. **`ship` espera el CI** (problema 2): aislado.
4. **`strict-tdd.md`**: aislado.
5. **`build-break-audit`**: fija dónde vive el loop.
6. **Lecciones en `knowledge-base.md`**: puede nombrar el loop recién cuando existe.
7. **Canal de informes** (problema 3): primero un spike del hook.
8. **Aprobación verificable + instalar `role-*`** (problema 1): CRITICAL, decide Juani.
9. **Bloque compacto** (problema 4): último, cuando el contenido deje de moverse.

### Tests y goldens que toca cada uno

| Change | Tests |
|---|---|
| `build-break-audit` | `build_break_audit_toggle_test.go` nuevo; `compose_all_toggles_{claude,opencode}.golden`; `catalog_toggles_test.go`; `docs_toggles_test.go` (ARCHITECTURE §2.1, CLAUDE §2, `harnesses.yaml:54-56`); `compose_test.go` (`TestKnownToggles_MatchesWhatComposeAccepts`) y el mensaje de error repetido a mano en `compose.go:66`; `dispatch_rules_test.go` si se toca el punto 2 |
| `strict-tdd.md` | `compose_test.go:128` (`TDDToggleOff`), goldens; crear `tdd_toggle_test.go`, que hoy no existe |
| `ship.md` | `ship_toggle_test.go` (orden y frases prohibidas), goldens, tests del comando en `internal/harness/command` |
| `knowledge-base.md` | `knowledge_base_toggle_test.go`, goldens |
| `governance.md` + `approval` | goldens; tests del subcomando en `cmd/active-stack` |
| informes | tests del hook/subcomando nuevos |
| `verify` | `internal/verify/*_test.go` |
| bloque compacto | casi todos los `*_toggle_test.go` (afirman frases sobre la salida de `Compose`) y los dos goldens |

---

## 5. Riesgos

**JUICIO** (orquestador salvo donde se indica):

- **Codex como constructor exige que el usuario tenga Codex** (y cuota: en el 3a se agotó la ventana
  de Plus en 2,5 h). Por eso va como invariante ("otra familia") y no como dependencia instalada. Si
  no hay otra familia disponible, el bloque tiene que decirlo en el informe, no fingir independencia.
  Confianza alta.
- **Más bloque.** `build-break-audit` suma palabras a un bloque que ya es el problema 4. Mitigación:
  escribirlo corto desde el principio y que el detalle viva en el comando de apply. Confianza media.
- **Choque con `review`**: tres pasadas paralelas por lentes contra un rompedor con vueltas. Si se
  activan los dos sin regla de convivencia, el costo se duplica (memoria: Juani no quiere tres pases
  adversariales). Hay que decidir si `build-break-audit` reemplaza a `review` cuando están los dos.
- ~~**Choque con `ship`**: PropTech autoriza push a `main` sin preguntar al archivar (D-049); `ship`
  pregunta siempre. Ver pregunta 4.~~ **Corregido (§8): es falso.** En PropTech siempre se preguntó
  antes de pushear a `main`. No hay choque: `ship` sigue preguntando siempre.
- **Choque con `roles.md:58`**: en apply no se convoca a nadie; el rompedor actúa entre vueltas de
  apply. Hay que redefinir "diff que se mueve" como "vuelta cerrada". (analista, contradicción 4)
  *Resuelto por §8*: el rompedor vive en `verify`, que siempre corre sobre un diff que dejó de
  moverse; la fila de apply de `roles.md` sigue siendo cierta.
- *Reemplazado (§8)*: "Más bloque" y "Choque con `review`" eran riesgos del toggle nuevo. Con el
  método dentro de las fases, el texto se agrega a fragmentos que ya existen y el choque con `review`
  desaparece porque el rompedor es `review`.
- **Firma criptográfica**: si Juani no tiene GPG o SSH signing configurado, el subcomando `approval`
  no puede decir nada, y el riesgo es que se relaje a "verificar el autor", que no prueba nada.
- **Hooks de Claude Code**: frontera externa que versiona Anthropic; el formato del transcript es
  interno. Sólo cubre Claude, no Codex ni OpenCode.
- **Bloque compacto**: el de más costo de reversión; mueve casi toda la batería de tests.
- **Instalar los `role-*`**: tipo de harness nuevo toca modelo, catálogo, uninstall y verify; y los
  archivos actuales de `~/.claude/agents` se pisarían (backup obligatorio, regla dura del repo).

---

## 6. Preguntas para Juani

1. **Firma**: ¿tenés o querés configurar firma de commits (GPG o SSH) con passphrase? Sin eso, "firmado
   por Juani" no se puede distinguir de "commiteado por un agente con tu nombre".
2. ¿El constructor de otra familia (Codex) lo querés como **regla** del bloque o como **opción** que se
   activa si está instalado? **Respondida (§8)**: opción, si el usuario tiene otro proveedor.
3. Cuando `build-break-audit` y `review` estén prendidos juntos, ¿el rompedor **reemplaza** a las tres
   pasadas o van las dos cosas? **Respondida en parte (§8)**: no hay `build-break-audit`; el rompedor
   *es* `review.md`. Queda abierto si es un agente con las tres lentes o tres pasadas (ver la propuesta).
4. ~~En PropTech pusheás a `main` sin preguntar al archivar (D-049). ¿Eso lo querés en Active Stack como
   excepción configurable, o `ship` sigue preguntando siempre?~~ **Corregido (§8)**: la premisa era
   falsa, en PropTech siempre se preguntó. `ship` sigue preguntando siempre por `main`.
5. ¿Active Stack tiene que **instalar los agentes `role-*`**, aunque eso pise los que tenés hoy en
   `~/.claude/agents` (con backup)?
6. Lo de mirar las pantallas en un navegador después de cada bloque de UI: ¿va en el bloque para
   cualquier proyecto, o sólo donde está impeccable?
7. ¿Arrancamos por el binario roto y el bloque viejo (puntos 1 y 2 de §0) antes que cualquier cosa de
   la metodología?

---

## 7. Lo nuevo de active-stack (Juan Cruz)

Repo `Group-Active-IA/active-stack`, clon en `/home/juanisarmiento/Proyectos/active-stack` (sólo
lectura). **Se agregó en V2 el remote `juancruz`** apuntando a ese clon, con las refs
`juancruz/main` (su `main` local) y `juancruz/origin-main` (su `origin/main`).

**EVIDENCIA** (orquestador, 2026-10-23 09:05-09:06)

```
$ git merge-base HEAD juancruz/origin-main          → a691ec0 (2026-08-17, "detectar dependencias también en el instalador gráfico")
$ git merge-base --is-ancestor a691ec0 HEAD         → sí
$ git rev-list --count juancruz/origin-main..HEAD   → 32   (V2 tiene 32 que ellos no)
$ git log juancruz/main..juancruz/origin-main       → e649ef5 (10/09), fae8bbf (25/09), ddeb94f (25/09)
$ git cherry -v HEAD juancruz/origin-main           → los tres con "+"
$ diff de las líneas +/- de harnesses.yaml en e649ef5 y 9f8d976 → iguales
$ git format-patch -1 --stdout fae8bbf | git apply --check → falla en mcp_test.go:484
$ git format-patch -1 --stdout ddeb94f | git apply --check → falla en overlays.go:191
$ diff de ids del catálogo                           → V2 tiene project-scaffold; ellos no tienen nada que V2 no tenga
$ toggles en su catálogo                             → [tdd, engram, model-routing, delegation, governance]
```

**Comparten historia**: V2 es `a691ec0` de ellos más 32 commits. Todo lo que mencionaba el encargo como
posible novedad de su lado (uninstall, backups en Windows, detección de dependencias en el instalador
gráfico, markers de jr-stack: `3252036`, `e0d791b`, `1f127d1`, `7ab8325`, `da4b0d6`, `a691ec0`) **ya está
en V2** por ascendencia.

| Qué | Lo tiene active-stack | Lo tiene V2 | Cómo traerlo | Riesgo |
|---|---|---|---|---|
| `e649ef5` harnesses discovery-research y web-scraper | sí | sí (`9f8d976`, mismo cambio, otro contexto) | nada | ninguno |
| `fae8bbf` OpenCode: MCP local con `command` como array y `environment` | sí | **no** (`mcp.go:268` sigue con `command` + `args`): Engram no queda activo en OpenCode | **cherry-pick con resolución** sólo en `mcp_test.go`; sacar el `Co-Authored-By` del mensaje (regla del repo); actualizar el comentario de `mcp.go:246` | bajo; `mcp.go` aplica limpio (arquitecto, E12) |
| `ddeb94f` OpenCode: `permission` entero en `__replace__` al cambiar de tier | sí | **parcial**: `4b31cc3` arregla las claves escalares; el residuo de overrides de `bash` sigue y `tier_switch_test.go:233` **exige** que siga | **reimplementar**: `__replace__` sólo sobre `bash`, no sobre todo `permission` (que borraría claves del usuario); portar sus tests y borrar el test que fija el residuo | ALTO (dominio de permisos, regla de governance del repo): alcance lo decide Juani |
| toggles evidence, auditor, roles, review, knowledge-base, senior, ship; subcomando `roles`; harness `project-scaffold`; CI ubuntu+windows; fixes de `--home`, rollback, self-install, filemerge, permisos por tier | no | sí | (esto es lo que V2 tiene y ellos no) | si Juan Cruz quiere converger, el camino es que su repo traiga V2, no al revés |

**Rechazado: merge** de `juancruz/origin-main`. Arrastra los trailers `Co-Authored-By` de `fae8bbf` y
`ddeb94f` (contra la regla del repo) y un duplicado de `9f8d976`. (arquitecto, confianza alta)

**SUPUESTO**: que `ddeb94f` borra claves de `permission` agregadas por el usuario sale de leer el diff;
no se ejecutó.

---

## Roles convocados y no convocados

**JUICIO** (orquestador, confianza alta, leí el encargo y no hay diff de código)

- Convocados: analista funcional y arquitecto (los de la fase explore), de a uno.
- Backend, frontend, DBA, DevOps, QA: no corresponden en explore; nada se construye.
- Senior: no convocado; la memoria `opsx-secuencial-con-auditor` dice que Juani no lo quiere en cada fase.
- Auditor: no convocado; este documento no tiene veredicto que auditar. Si va a propose, conviene que un
  auditor lea los dos informes contra este documento.
- `active-stack roles`: **no se corrió con éxito**; el binario instalado es de tests (§0).

---

## 8. Decisiones de Juani del 23/10

Dichas por Juani al orquestador el 2026-10-23, después de leer este documento. Reemplazan lo marcado
*Reemplazado (§8)* más arriba; el resto de la exploración se mantiene.

**El método no va como toggle separado.** Va dentro de las fases OPSX que ya existen:

1. **`/opsx:apply` es el CONSTRUCTOR.** Escribe con TDD estricto y pega el rojo (y el verde) en
   `tasks.md`. Si el usuario tiene otro proveedor disponible (Codex/GPT), el constructor puede ser de
   otra familia de modelo, para que la revisión sea independiente. Es opción, no requisito.
2. **`/opsx:verify` es conformidad (lo que ya hace) → ROMPEDOR → AUDITOR.** El rompedor es la
   revisión adversarial que ya existe (`review.md`) con el estilo de PropTech: reproducciones
   concretas, mutantes que fallan por la aserción (L-019, L-034), contra un diff que dejó de moverse.
   El auditor (`auditor.md`) lee los informes y decide.
3. **El bucle.** Si hay bloqueantes, vuelven a `/opsx:apply`, el constructor corrige y se repite
   `/opsx:verify`. Se sale cuando el auditor aprueba o a la tercera vuelta; en la tercera, lo que
   queda abierto se le presenta a la persona escrito como decisión suya, nunca en silencio. Una
   corrección que agrega una regla nueva vuelve a pedir firma humana antes de construirse (L-029).
4. **Senior sólo en `verify`**, para separar hallazgos reales de los que contradicen una decisión
   vieja. Deja de estar en todas las fases.
5. **`/opsx:archive` actualiza toda la base de conocimiento que el change tocó**: REWRITE de lo que
   dice cómo es ahora; APPEND de decisiones, lecciones numeradas y deuda. Cada cambio a la base sale
   con la evidencia de su razonamiento (qué cambió y por qué: el commit, el test o la decisión que lo
   respalda), tipado EVIDENCIA/JUICIO/SUPUESTO. Un archive que sólo agrega y no reescribe se reporta
   como sospechoso.
6. **Corrección del push.** La pregunta 4 y el riesgo "Choque con `ship`" decían que en PropTech se
   pusheaba a `main` sin preguntar al archivar. **Es falso: siempre se preguntó.** `ship` sigue
   preguntando siempre por `main`. Además, `ship` no da nada por listo hasta que el CI completo esté
   en verde sobre el commit pusheado (L-037).

La contradicción 5 del informe del analista se apoya en la misma lectura y queda sin efecto. El
informe no se edita: es textual.

**Se mantiene de la exploración**: arreglar la instalación rota de esta máquina (binario de tests en
`~/.local/bin/active-stack`, bloque del 9/09) y que `verify` detecte un bloque instalado
desactualizado; traer los dos arreglos de OpenCode de Juan Cruz (`fae8bbf` por cherry-pick,
`ddeb94f` reimplementado); informes guardados por hook `SubagentStop` o subcomando, sin dar Write a
los revisores; aprobación verificable con firma de commits GPG/SSH (sigue siendo pregunta para
Juani); instalar los `role-*`; bloque compacto al final.
