# Informe del arquitecto (role-architect, opus) · 2026-10-23

Extraído textual del transcript del sub-agente con jq (último mensaje de texto), sin transcripción del orquestador.

Revisé las diez piezas, cada una con dónde iría, la dirección de sus dependencias, la alternativa que descarté y lo que costaría revertirla. Antes de recomendar cualquier cosa hay tres hallazgos:

- **El binario instalado es un binario de tests.** `go version -m` sobre `~/.local/bin/active-stack` dice `cmd/active-stack/headless.test`, fechado hoy 09:05. Por eso `roles` falla: hoy el binario instalado no trae el comando.
- **El bloque depende de agentes que Active Stack no instala.** El texto del orquestador y el clasificador nombran `role-dba`, `role-auditor` y el resto, pero esos archivos no vienen de este repo.
- **La firma por autor de commit no prueba nada en esta máquina.** Los agentes commitean como Juani, así que el autor no distingue a la persona del agente.

Todo lo que sigue es de sólo lectura, de hoy (2026-10-23) entre las 09:09 y las 09:12.

## EVIDENCIA

**Estructura del bloque**
- E1 `compose.go:34-38`: el orden actual es delegation, model-routing, governance, evidence, auditor, roles, review, knowledge-base, senior, ship, engram, tdd. El mensaje de error de `compose.go:66` repite la lista a mano.
- E2 `wc -w` sobre la base de Claude más los fragmentos da **17311** palabras. Sólo la base son 4297.
- E3 La base ya tiene una convención para nombrar otros toggles sin depender de ellos: lo hace de forma condicional ("Where the `roles` fragment is on the page", `sdd-orchestrator.md:369,385,465`).
- E4 `ship_toggle_test.go:200` exige que ship quede después de review y senior y antes de engram. `ship_toggle_test.go:554-627` tiene frases prohibidas como `"waits for confirmation"` y `"asks and waits"`.

**Agentes `role-*` (pieza 3)**
- E5 `git log --all -S'role-implementer'` no devuelve nada. Los `HarnessType` válidos son skill, config, external y command (`model/harness.go:10-19`), y el catálogo no tiene ningún tipo "agent".
- E6 `eza -l ~/.claude/agents/` muestra los diez `role-*.md` fechados el 9 de septiembre, sin dueño en este repo.
- E7 Igual los nombran `sdd-orchestrator.md:356-359,454-457` y `internal/roles/classify.go:37-40`.
- E8 `role-implementer.md:178-193` dice: "No escribís código de autenticación... **no existe el canal**". Para ese archivo, nada que un agente pueda escribir cuenta como aprobación.

**Ship (pieza 6)**
- E9 La doctrina de ship está dos veces: en el fragmento `assets/ship.md` (2135 palabras) y en el comando `assets/commands/claude/opsx/ship.md` (3786 palabras, lo instala `internal/harness/command/installer.go:41`).

**Binario y bloque instalado (pieza 9)**
```
$ go version -m ~/.local/bin/active-stack
  path github.com/Group-Active-IA/active-stack/cmd/active-stack/headless.test
$ ~/.local/bin/active-stack foo     → flag provided but not defined: -project / Usage of uninstall
$ eza -l ~/.local/bin/active-stack  → 9,4M  23 oct 09:05
```
- E10 `rg 'active-stack roles' ~/.claude/CLAUDE.md` no encuentra nada. `internal/verify` importa model y permissions, pero no `harness/config`.

**Transcripts de sub-agentes (pieza 7)**
- E11 Hay 857 archivos `agent-*.jsonl` en `~/.claude/projects`, y **ninguno** bajo `-home-juanisarmiento-Proyectos-PropTech`. La última línea de una muestra tiene `type=assistant`, `isSidechain=True` y `agentId`.

**Commits de juancruz (pieza 10)**
- E12 `git show fae8bbf -- internal/harness/external/mcp.go | git apply --check` aplica limpio. Lo que falla es sólo `mcp_test.go:484`. En V2, `mcp.go` no cambió desde `a691ec0`.
- E13 `ddeb94f` falla en `overlays.go:191`. `4b31cc3` no está en `a691ec0`.
- E14 `tier_switch_test.go:233` (`TestOpencodeBashResidueSurvivesTierSwitch`) **exige** que el residuo de balanceado siga presente. Con `ddeb94f` ese test se pone rojo, y su propio mensaje pide borrarlo en ese caso.
- E15 Los patch-id de `e649ef5` y `9f8d976` son distintos (39add73… contra dfea770…). El `diff` entre los dos parches sólo cambia índices y contexto de `harnesses.yaml`.
- E16 `fae8bbf` y `ddeb94f` traen `Co-Authored-By: Claude Sonnet 5`.

## JUICIO (arquitecto; comparé contra `compose.go`, la convención de E3 y los tests citados)

**1. Ciclo constructor, rompedor, auditor → (b) toggle nuevo `build-break-audit`, entre `review` y `knowledge-base`.**
- Depende de review, auditor y roles, todos compuestos antes, y los nombra sólo en forma condicional (E3).
- Regla para no crear un ciclo: ningún fragmento anterior lo menciona a él. Ni review.md ni auditor.md se editan para nombrar el loop.
- El "veredicto textual a la persona" ya está en `auditor.md` ("Who it answers to"), así que no se duplica: se cita.
- Rechazo convertir `review.md` en el loop secuencial. Borraría las tres lentes que `sdd-orchestrator.md:472` declara que no son overhead, y revertirlo costaría reescribir `review_toggle_test.go` entero. El toggle nuevo se revierte sacando el toggle.
- Tests que toca: `build_break_audit_toggle_test.go` nuevo, los dos goldens, `catalog_toggles_test.go`, `docs_toggles_test.go` (ARCHITECTURE §2.1, CLAUDE §2, `harnesses.yaml:54-56`) y `compose.go:66`. Confianza alta.

**2. Constructor de otra familia de modelo → dentro del fragmento de la pieza 1, como invariante ("quien rompe y quien audita son de otra familia que quien construye"). No va en `model-routing`.**
- `model-routing` es la tabla de alias que se pasan al parámetro `model` del Agent tool. Codex se invoca con `codex exec` por Bash, no es un alias. Meterlo ahí haría depender un marcador sustractivo de un binario externo.
- Lo específico de PropTech (sin sandbox, `docker ps`, cuotas) → (e): no va en Active Stack.
- Dependencia nueva: fragmento → CLI `codex`, que hoy no instala nadie. Revertir es borrar un párrafo. Confianza media.

**3. Aprobación humana verificable → (a) en `governance.md`, que defina el artefacto (tabla `## Aprobación` en proposal.md), más (d) un subcomando `active-stack approval --change X` que lo verifique.**
- El autor del commit **no discrimina**, porque los agentes commitean con el `user.name` de Juani. Para que el subcomando decida algo necesita `git verify-commit`, o sea una firma criptográfica. Sin eso, verifica forma y no autoría.
- Hoy Active Stack no instala los `role-*` (E5-E7) y aun así el bloque y el clasificador cuentan con ellos. Es una dependencia que nadie nombró ni versionó.
- Para que `role-implementer` acepte la firma hay que elegir entre (c) instalar los `role-*` como harness, con un tipo nuevo `agent` que reuse el patrón de `harness/command`, o (e) editarlos a mano fuera del repo. Recomiendo (c): si no, la regla de E8 queda sin dueño.
- Revertir un tipo de harness nuevo es caro: toca el modelo, el catálogo, el uninstall y verify. Gobernanza CRITICAL, decide la persona. Confianza media.

**4. TDD con rojo y verde pegados en tasks.md, y mutante → (a) modificar `strict-tdd.md`.**
- Sin dependencias nuevas.
- El mutante incluye lo de L-013: copia previa (`stash`) y decirlo en el informe.
- Rechazo un toggle `mutation-check` aparte, porque partiría la definición de TDD en dos.
- Tests: `compose_test.go:128` (TDDToggleOff) y los goldens. No hay `tdd_toggle_test.go`: conviene crearlo. Confianza alta.

**5. Lecciones L-NNN → (a) en `knowledge-base.md`, como tercera operación (sólo agregar), con todas las reglas de lecciones ahí.**
- Las que nombran al auditor y al loop van en forma condicional: "si el auditor está, su veredicto cierra con LECCIONES"; "si build-break-audit está, los briefs de constructor y rompedor llevan el documento". Así todas las flechas van de knowledge-base hacia fragmentos anteriores y no hay ciclo.
- Rechazo poner "proponé LECCIONES" en `auditor.md` (posición 4c): se referiría hacia adelante, a un concepto que define un toggle posterior.
- El nombre del archivo no se fija. `knowledge-base.md` ya prohíbe asumir los nombres de kb-creator. Que kb-creator cree el documento es (e), en su propio repo.
- Tests: `knowledge_base_toggle_test.go` y goldens. Confianza media-alta.

**6. Ship espera el CI → (a) modificar `ship.md` y también `assets/commands/claude/opsx/ship.md` (E9).**
- "Listo" pasa a significar CI verde sobre el SHA pusheado (`gh run list --commit`). "No hay CI configurado" se reporta como tal, no como verde.
- Hay que esquivar las frases prohibidas de E4.
- La dependencia de `gh` ya está declarada como suposición en ship.md, no es nueva.
- Riesgo: que las dos copias diverjan. Revertir es barato. Confianza alta.

**7. Informes textuales fuera del repo → (d) recomiendo un hook `SubagentStop` instalado por Active Stack que copie el último mensaje a `.git/informes-<cambio>/`.**
- Así no pasa por ningún LLM y nadie necesita permiso de Write.
- Fallback: un subcomando que lea el jsonl.
- Rechazo darle Write acotado por carpeta a los revisores. El launch pattern no lleva lista de herramientas (`roles.md:174-190`), los permisos por ruta son de sesión y no de sub-agente, y reabre justo el agujero de "el revisor edita".
- El subcomando solo es frágil. El formato jsonl es interno y sólo existe en Claude, y en PropTech no hay ninguno (E11).
- Dependencia: Active Stack → el contrato de hooks de Claude Code, una frontera externa que versiona Anthropic. Confianza media-baja.

**8. Bloque demasiado grande → combinar dos cosas.**
- En CLAUDE.md queda la base y, por toggle, un resumen con su "terminado" y las reglas que deciden antes de actuar: push a main pregunta, el gate humano, auditor textual.
- El detalle va a canales que **Active Stack ya posee**: `/opsx:verify` (review, auditor, roles) y `/opsx:ship`. El de ship hoy está duplicado (E9).
- Rechazo una skill por toggle cargada bajo demanda: agrega una arista config → skill en el instalador, y si la skill no se carga la regla no se aplica, sin ninguna señal.
- Costo de revertir: alto. Casi todos los `*_toggle_test.go` afirman frases sobre la salida de `Compose`.
- Antes de decidir conviene medir qué secciones se citan de verdad. Confianza media.

**9. Para que `roles` se use → (d) en `verify` y en el pipeline.**
- (i) Un check en `verify` que compare el bloque instalado contra `Compose(toggles)` y avise si quedó viejo. Crea la arista verify → harness/config, que no existe hoy (E10) y no forma ciclo.
- (ii) Un check de binario: que `go version -m` diga el path `cmd/active-stack` y no `.test`, y que `active-stack roles --help` salga con código 0.
- (iii) Un subcomando `version`, que no existe.
- Lo primero de todo es encontrar qué test de `cmd/active-stack/headless` hace self-install sobre el HOME real. Eso le corresponde a QA y backend, no a mí. Confianza alta en el síntoma, baja en la causa.

**10. Commits de juancruz**
- `fae8bbf` → **cherry-pick con resolución**, sólo de `mcp_test.go`, y reescribir el mensaje para sacar el `Co-Authored-By` (E16, va contra la regla de CLAUDE.md). Al resolver, actualizar también los comentarios de `mcp.go:246`, que siguen describiendo command+args para opencode.
- `ddeb94f` → **reimplementar**. Choca con `4b31cc3` en texto (E13) y en semántica (E14). Además envuelve todo `permission` en `__replace__`, lo que borra claves que el usuario haya agregado a mano. `4b31cc3` proponía `__replace__` sólo sobre `bash`. Es dominio ALTO: el alcance lo elige la persona. Hay que portar sus tests.
- `merge` → rechazado. Arrastra los trailers de E16 y suma `e649ef5`, que repite `9f8d976` con otro contexto (E15).

## SUPUESTOS

- S1 Que `ddeb94f` tiene el comportamiento que digo (borrar claves `permission` del usuario) lo leí del diff. No ejecuté el merge.
- S2 Que el payload de `SubagentStop` trae la ruta del transcript del sub-agente: no lo verifiqué.
- S3 Que la ausencia de transcripts en PropTech se debe a la retención: no lo verifiqué.
- S4 Que el script de firma de PropTech no usa GPG: no lo abrí.
- S5 Que la causa del binario de tests es un test con self-install: es una hipótesis. Lo único verificado es el síntoma.

## Orden recomendado

1. **Pieza 9.** Sin binario ni bloque al día, nada de lo demás llega a la máquina de quien lo usa.
2. **Pieza 10, la parte de `fae8bbf`.** Bug real y barato (Engram no queda activo en OpenCode).
3. **Pieza 6.** Cambio aislado y de bajo riesgo.
4. **Pieza 4.** Aislado, sin dependencias.
5. **Pieza 1 junto con la 2.** El toggle nuevo, que fija dónde va el loop.
6. **Pieza 5.** Depende de que la 1 ya exista para poder nombrarla.
7. **Pieza 7.** Primero un spike sobre el hook.
8. **Pieza 3.** Gate CRITICAL y decide la persona.
9. **Pieza 8.** El que más tests remueve. Va último, cuando el contenido deje de moverse.
10. **El tier de `ddeb94f`.** En paralelo cuando haya decisión sobre el alcance del `__replace__`.

Archivos que abrí, todos bajo `/home/juanisarmiento/Proyectos/active-stack-V2/` salvo donde se indica:
- `ARCHITECTURE.md`, `CLAUDE.md`
- `internal/harness/config/compose.go`, `assets/claude/sdd-orchestrator.md:334-500`, `ship_toggle_test.go:200-230,554-600`, `permissions/tier_switch_test.go:222-275`
- `cmd/active-stack/main.go:85-120`, `internal/catalog/harnesses.yaml:139-155`, `internal/harness/command/installer.go`
- `/home/juanisarmiento/Proyectos/PropTech/knowledge-base/14_forma_de_trabajo.md` y las líneas 1-60 de `15_lecciones.md`
- `/home/juanisarmiento/.claude/agents/role-implementer.md`, líneas 1-12 y 175-195
- Además, búsquedas con `rg` en `review.md`, `auditor.md`, `knowledge-base.md`, `verify/` y `roles/`, y los `git show` de `fae8bbf`, `ddeb94f` y `4b31cc3`.

No abrí `docs/exploraciones/` (incluido el informe del analista) para no contaminar esta lectura. No guardé nada en engram.
