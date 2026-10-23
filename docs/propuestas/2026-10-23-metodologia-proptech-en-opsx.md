# Propuesta: la metodología de PropTech dentro de las fases OPSX

Fecha: 2026-10-23 · Rama: `feat/metodologia-proptech` · Fase: propose (no se implementa nada).
Sale de: `docs/exploraciones/2026-10-23-metodologia-proptech.md`, §8 (decisiones de Juani del 23/10).

**Dónde vive esta propuesta y por qué.** La convención del repo es que cada incremento sea un change
OPSX y que el plan esté en `CHANGES.md`. `openspec/` está ignorado por git y en esta máquina no está
inicializado (`openspec list` → "No OpenSpec root found"), así que un `proposal.md` ahí no quedaría en
ningún commit. Por eso el plan va acá y en `CHANGES.md` (epic nuevo, al final). Cuando Juani firme,
cada paso del §5 se abre con `openspec new change "<id>"` y su `proposal.md` cita esta página.

---

## 1. Qué cambia y por qué

Hoy el bloque instalado arma el ciclo con piezas sueltas: `review.md` lanza tres pasadas que
reportan y no bloquean, `auditor.md` lee informes y le habla a la persona, `senior.md` acompaña las
cinco fases, `knowledge-base.md` reescribe y agrega en archive, `ship.md` termina en el PR abierto.
Ninguna dice qué pasa con un hallazgo bloqueante, ni cuántas veces se vuelve a construir.

PropTech probó otro ciclo y Juani decidió traerlo **dentro de las fases que ya existen**, sin toggle
nuevo:

```
/opsx:apply   CONSTRUCTOR   TDD estricto; rojo y verde pegados en tasks.md;
                            de otra familia de modelo si el usuario tiene otro proveedor
      │
      ▼
/opsx:verify  conformidad → ROMPEDOR → SENIOR → AUDITOR
                            (lo de hoy)  (review.md)  (senior.md)  (auditor.md: decide)
      │
      ├── auditor aprueba ─────────────────────────────► /opsx:archive
      ├── bloqueantes, vuelta 1 o 2 ──► /opsx:apply (corrige sólo eso) ──► /opsx:verify
      └── vuelta 3 con algo abierto ──► la persona decide, por escrito; nada avanza sin respuesta
      (una corrección que agrega una regla nueva vuelve a pedir firma antes de construirse, L-029)

/opsx:archive REWRITE de "cómo es ahora" + APPEND de decisiones, lecciones L-NNN y deuda,
              cada entrada con su evidencia (commit, test o decisión), tipada.
              Sólo-APPEND sin REWRITE = se reporta SOSPECHOSO.
/opsx:ship    pregunta siempre por main; "listo" = CI completo en verde sobre el SHA pusheado
```

**SUPUESTO** (no está en las decisiones de Juani, lo propongo yo): el senior va **entre** el
rompedor y el auditor, para que el auditor decida leyendo también qué hallazgos contradicen una
decisión vieja. Es la pregunta abierta 2.

---

## 2. Evidencia de partida

**EVIDENCIA** (sub-agente de propose, 2026-10-23 09:27)

```
$ eza assets/commands/claude/opsx              → ship.md  verify.md
$ rg -c "method: embed" internal/catalog/harnesses.yaml   → 0
$ rg -n "archive | haiku" internal/harness/config/assets/*/sdd-orchestrator.md | wc -l   → 8
$ rg -n "in \*\*every phase\*\*" internal/harness/config/assets/senior.md   → 5: ...convened in **every phase**...
$ wc -w internal/harness/config/assets/{review,auditor,senior,knowledge-base,strict-tdd,ship}.md
    1211 / 475 / 1927 / 2935 / 606 / 2135
```

Lo que eso dice:

1. **Active Stack sólo instala dos comandos OPSX**: `/opsx:verify` y `/opsx:ship`. `/opsx:apply` y
   `/opsx:archive` los genera `openspec init`. Las skills `assets/skills/openspec-*` son placeholders
   de 24 palabras que ningún harness instala. Entonces **las reglas de apply y archive sólo pueden
   viajar por el bloque** (`strict-tdd.md` y `knowledge-base.md`); las de verify pueden ir al bloque
   y al comando.
2. **El orden de composición manda dónde va cada texto** (`compose.go:45-59`): governance → evidence
   → auditor → roles → review → knowledge-base → senior → ship → engram → tdd. Un fragmento puede
   nombrar a uno anterior; a uno posterior sólo en forma condicional ("where the X fragment is on the
   page"). Por eso **el bucle va en `review.md`**: es el primer fragmento posterior al auditor y
   anterior al senior y a la base de conocimiento.
3. El modelo de archive es `haiku` en las 8 variantes del bloque base.

---

## 3. Qué cambia, fragmento por fragmento

Los fragmentos están en inglés; acá va qué dirán, no el texto final.

| Pieza | Hoy | Cambio |
|---|---|---|
| `strict-tdd.md` (toggle `tdd`, fase apply) | tabla de tildes en el informe, sin hora ni salida (E1 del analista) | **Constructor.** Rojo y verde pegados en `tasks.md` bajo la tarea: comando, `date +%T`, últimas líneas; `[x]` sólo con los dos. Test que nació verde: mutante de 2-4 líneas que falla **en la aserción que el test dice probar**, pegando `archivo:línea` de la falla, y revertido (L-019, L-034); sin commit, copia previa (L-013). Suites de cierre después del último commit, con el SHA (L-031). "Verde" son todos los pasos del CI (L-005). Brief: nombra consumidores fuera del área (L-023) y separa conflicto de redacción (resuelve y declara) de regla nueva (frena, L-029/L-030). En vuelta 2 y 3 corrige **sólo** la lista del auditor; lo nuevo va a deuda. **Otra familia**: si hay otro proveedor (p. ej. `codex exec`), el constructor puede correr ahí; si no, el informe lo dice y no finge independencia. |
| `review.md` (toggle `review`, verify) | tres pasadas por lentes, reportan, no bloquean, `review-verdict.md` | **Rompedor + bucle.** Conformidad primero (sin cambios). El rompedor rompe con **reproducciones concretas** (comando y salida) y deja cada falla como test que falla o mutante que sobrevive; las tres lentes siguen siendo obligatorias (ver pregunta 1). Sección nueva "The loop": vueltas numeradas, tope 3, qué vuelve a apply, la salida por aprobación del auditor o por la vuelta 3 con decisión escrita de la persona, y L-029. "It reports, it does not block" se mantiene: **el rompedor reporta, decide el auditor, y quien frena es el orquestador**. Clave de Engram por vuelta y rol (`opsx/{change}/verify/round-{n}/{role}`), que el propio fragmento ya pide. Veredictos por vuelta en el change (`verify-round-{n}.md`). |
| `auditor.md` (toggle `auditor`) | lee informes, habla textual a la persona, puede callar | **Decide.** Su veredicto es APROBADO o BLOQUEANTES (lista) + DEUDA (lista) + LECCIONES (o "ninguna"). Sigue reportando a la persona, textual. "Puede callar" pasa a "si todo se sostiene, una línea: APROBADO". No nombra el bucle (va antes de `review` en la composición). |
| `senior.md` (toggle `senior`) | "convened in **every phase**", tabla de cinco fases, "five per change" | **Sólo verify.** Preámbulo, tabla y costo pasan a una fase: separa hallazgos reales de los que contradicen una decisión fechada. Uno por vuelta de verify (máx. 3 por change). Se mantienen el umbral BLOQUEANTE/NOTA, "the senior labels, you hold" y los canales de memoria. |
| `knowledge-base.md` (toggle `knowledge-base`, archive) | REWRITE + APPEND + PENDING, reporte como condición de cierre | **Toda la base que el change tocó**: lo que nombra `## Knowledge Base Impact` **más** lo que muestra el diff. APPEND tiene tres clases: decisiones, **lecciones L-NNN** (sólo se agregan; salen de la sección LECCIONES del auditor) y **deuda** (sale de la lista DEUDA del auditor). **Cada entrada lleva su evidencia**: qué cambió y qué lo respalda (SHA, test o decisión), tipada. El detector cambia de "if you only appended, you did not update" a una etiqueta en el reporte: **SUSPICIOUS: append-only**. |
| `ship.md` + `assets/commands/claude/opsx/ship.md` | termina en el PR abierto; `main` pregunta siempre | `main` **sigue preguntando siempre** (sin cambio). Nuevo: nada está listo hasta que el CI completo esté verde sobre el SHA pusheado (L-037): `gh run list --commit <sha>` y `gh run watch`; tres salidas: verde, rojo (se reporta, no se apila otro arreglo), "no hay CI / no pude consultar" (nunca leído como verde). "The cycle ends at the open PR" pasa a "...at the open PR **and its CI result**". Se cambian las dos copias. |
| `assets/commands/claude/opsx/verify.md` | conformidad sola | Paso 9 corto: "si el bloque trae rompedor y auditor, seguí con ellos; el veredicto del auditor decide si se archiva". La regla vive en `review.md`; el comando sólo la encadena. |
| `governance.md` | "explicit human approval", sin artefacto | Tabla `## Aprobación` en `proposal.md` (elección, persona, fecha, alcance), que ningún agente completa; una regla nueva después de la firma lleva fila propia (L-029). La verificación criptográfica queda para el paso 8. |
| Base, tabla de modelos (8 variantes) | archive `haiku` | archive `sonnet`: reescribir la base con evidencia no es "file operations". Pregunta 5. |
| `roles.md`, base Dispatch punto 2 | apply no convoca; lo independiente sale junto | **Sin cambio.** El rompedor vive en verify, que siempre lee un diff quieto; constructor → rompedor → auditor dependen uno del resultado del otro, así que la prueba de dependencia del punto 2 ya los manda en secuencia. |
| `compose.go`, catálogo | 12 toggles | **Sin toggle nuevo.** Sí cambian las descripciones de `senior` y `review` en `harnesses.yaml`, `CLAUDE.md`/`AGENTS.md` §2, `ARCHITECTURE.md` §2.1 y `README.md`. |

---

## 4. Tests y goldens que se mueven

| Paso | Tests |
|---|---|
| ship espera el CI | `ship_toggle_test.go`: `NeverMerges` (la frase "the cycle ends at the open pr"), lista nueva de lecturas; cuidar las frases prohibidas de `DoesNotClaimAMechanismItLacks` ("waits for confirmation", "asks and waits"): el texto dice que **el orquestador** espera el run, no que algo lo garantice. Goldens `compose_all_toggles_{claude,opencode}.golden`. Tests del comando en `internal/harness/command` sólo si cambia la ruta. |
| constructor (`strict-tdd.md`) | `tdd_toggle_test.go` **nuevo** (hoy no existe); `compose_test.go` `TestCompose_TDDToggleOff`; `dispatch_rules_test.go` `TestElRepartoNoAutorizaRecortarLaRevision` (no tocar sus frases); goldens. |
| rompedor + bucle (`review.md`, `auditor.md`, `verify.md`) | `review_toggle_test.go`: `ThreeIndependentPasses` y `IndependenceIsByConstruction...` cambian si la pregunta 1 sale "un agente"; nuevos para el bucle (tope 3, salida por aprobación, vuelta 3 como decisión escrita, L-029, clave por vuelta); `ReportsDoesNotBlock` se mantiene. `auditor_toggle_test.go`: nuevos (APROBADO/BLOQUEANTES/DEUDA/LECCIONES). `review_cost_rules_test.go` `TestNingunaReglaNuevaRecortaAlAuditorNiLasLentes` **no se toca**: guarda que las lentes no se recorten. `verify_phase_test.go` y `assets/commands_test.go` por el paso 9. Goldens. |
| senior sólo verify | `senior_toggle_test.go`: `PreambleSaysItIsEveryPhase`, `AsksADifferentQuestionInEachPhase`, `RunsAtTheCloseOfEachPhase`, `IsNotAnEighthRole`, `DeclaresItsCost` se reescriben; `InVerifyItJudgesTheFindings` se mantiene. `docs_toggles_test.go` (descripciones en docs y catálogo). Goldens. |
| archive con evidencia | `knowledge_base_toggle_test.go`: `TheDetectorIsTheRewriteList` (etiqueta SUSPICIOUS), `ArchivePromptCarriesTheReport` (el prompt pide lecciones, deuda y evidencia), nuevos para L-NNN y "base que el diff tocó". `verify_phase_test.go` si cambia la fila de archive. Goldens. |
| governance | goldens; un test nuevo para la tabla `## Aprobación` y L-029. |

Regla para todos: **cada test nuevo se ve fallar antes** contra el fragmento viejo (el repo ya tropezó
siete veces con asserts que pasaban sin probar nada, según `review_toggle_test.go:213-217`).

---

## 5. Orden de construcción, en pasos chicos

Cada paso es un change OPSX de ≤ 12 tareas, con su propia firma si es ALTO.

| # | Change | Qué | Governance | Depende de |
|---|---|---|---|---|
| 0 | `instalacion-sana` | Encontrar qué test escribe un `.test` en `~/.local/bin`; `internal/verify` chequea que el binario sea `cmd/active-stack`, que `active-stack roles --help` salga 0 y que el bloque instalado sea igual a `Compose(toggles)`; subcomando `version`. Después, reinstalar en la máquina de Juani (con backup). | MEDIO | — |
| 1 | `opencode-mcp-command-array` | Cherry-pick de `fae8bbf` con resolución en `mcp_test.go`, sin el trailer de coautoría. | MEDIO | — |
| 2 | `ship-espera-ci` | §3 fila ship, en las dos copias. | BAJO | — |
| 3 | `apply-constructor` | §3 fila `strict-tdd.md`. | BAJO | — |
| 4 | `verify-rompedor-auditor` | `review.md` (rompedor + bucle), `auditor.md` (decide), paso 9 de `verify.md`. | BAJO | 3 (el bucle nombra al constructor) |
| 5 | `senior-solo-en-verify` | `senior.md` + docs y catálogo. | BAJO | 4 (su lugar en la secuencia) |
| 6 | `archive-base-completa` | `knowledge-base.md` (lecciones, deuda, evidencia, SUSPICIOUS); archive a `sonnet` si Juani dice que sí. | BAJO | 4 (LECCIONES y DEUDA salen del auditor) |
| 7 | `governance-aprobacion` | Tabla `## Aprobación` + L-029 en `governance.md`. | BAJO | 4 |
| 8 | `aprobacion-firmada` | Subcomando `active-stack approval --change` con `git verify-commit`. **Sólo si Juani configura firma** (pregunta 3). | ALTO | 7 |
| 9 | `informes-por-hook` | Spike de `SubagentStop`; después hook + subcomando de respaldo que copia el último mensaje del sub-agente a `.git/informes-<change>/`. | MEDIO | 4 |
| 10 | `opencode-permission-replace` | Reimplementar `ddeb94f`: `__replace__` sólo sobre `bash`. | ALTO | — |
| 11 | `instalar-roles` | Tipo de harness `agent` para los `role-*.md` (backup de `~/.claude/agents`). | ALTO | 9 |
| 12 | `bloque-compacto` | Base + resumen corto por toggle; el detalle a `/opsx:verify`, `/opsx:ship` y los `role-*`. | MEDIO | todo lo anterior |

0, 1, 2 y 3 son independientes entre sí. 2 a 7 son texto de fragmentos y tests de Go; 8 a 12 tocan
código, instalación o permisos.

---

## 6. Riesgos

**JUICIO** (sub-agente de propose; confianza entre paréntesis)

- **Apagar `review` apaga el bucle.** Con el método dentro de `review.md`, quien quiera el ciclo sin
  rompedor se queda también sin vueltas. Es el costo de no tener toggle propio (alta).
- **Tres pasadas por vuelta × tres vueltas** = hasta nueve lanzamientos de rompedor más tres
  auditores y tres seniors, contra el "pocos agentes" de Juani. Depende de la pregunta 1 (alta).
- **La vuelta 3 sin respuesta.** Si la persona no contesta, nada avanza; el fragmento tiene que
  decirlo como default seguro del orquestador, igual que hoy dice `senior.md` (alta).
- **Los informes siguen pasando por el orquestador** hasta el paso 9: el auditor puede leer una
  copia, no el original (R14). Mientras tanto, la copia se saca del transcript con `jq` y se declara
  como tal (media).
- **"Otra familia" no se verifica.** Sin campo de proveedor en el registro del despacho, que el
  constructor fue Codex es palabra del orquestador. Se declara como SUPUESTO en el informe (media).
- **Firma**: sin GPG/SSH configurado, L-029 queda como pedido sin prueba, y el riesgo es relajarlo a
  "el autor del commit", que no prueba nada en esta máquina (alta).
- **Goldens**: cada paso de 2 a 7 regenera los dos goldens; si dos pasos corren en paralelo chocan
  en esos archivos. Van en secuencia (alta).
- **`openspec` cambia sus comandos.** Apply y archive no son nuestros; si `openspec` cambia
  `/opsx:archive`, nuestras reglas siguen en el bloque pero el comando puede contradecirlas (baja).

---

## 7. Alternativas rechazadas

**Toggle separado `build-break-audit`** (la propuesta de la exploración, §4). Rechazada por Juani el
23/10: el método no es un modo aparte, es cómo corren apply, verify y archive. Además duplicaba al
rompedor con `review` y necesitaba una regla de convivencia.

*Costo de revertir* (JUICIO, confianza media): pasar el bucle a un toggle propio más adelante es un
change de unas 10 tareas: sacar la sección "The loop" de `review.md` y la decisión de `auditor.md` a
un fragmento nuevo, sumar el toggle en `compose.go` (`knownToggles`, el orden y el mensaje de error
escrito a mano en `compose.go:66`), catálogo, cuatro documentos y un archivo de tests nuevo, y volver
a mover los tests de los pasos 4 y 5. No toca código de instalación ni datos del usuario.

**Comandos propios `/opsx:apply` y `/opsx:archive` embebidos.** Rechazada: pisarían los que genera
`openspec init` y quedaríamos manteniendo una copia de algo ajeno. Las reglas viajan por el bloque.

---

## 8. Preguntas abiertas para Juani

1. **¿El rompedor es un agente o tres?** Un agente que cubre las tres lentes ahorra tokens pero
   pierde la independencia entre lentes; tres pasadas de a una la mantienen y cuestan el triple por
   vuelta. Mi recomendación: un agente con las tres lentes como secciones obligatorias en la vuelta 1
   y sólo sobre lo corregido en las vueltas 2 y 3.
2. **¿El senior va antes o después del auditor?** Propongo antes, para que el auditor decida sabiendo
   qué hallazgos chocan con decisiones viejas.
3. **¿Configurás firma de commits (GPG o SSH, con passphrase)?** Sin eso, el paso 8 no se hace y
   L-029 queda como pedido sin prueba.
4. **Constructor de otra familia**: ¿alcanza con que el informe diga qué modelo construyó, o querés
   que Active Stack lo chequee (campo de proveedor en el registro del despacho)?
5. **¿Archive pasa de `haiku` a `sonnet`?** Reescribir la base con evidencia es juicio, no mover
   archivos.
6. Siguen abiertas de la exploración: instalar los `role-*` pisando los tuyos (con backup), y si lo
   de mirar la UI en un navegador va para todo proyecto o sólo con impeccable.

---

## Knowledge Base Impact

Este repo no tiene `knowledge-base/`; sus documentos de referencia son estos:

- `ARCHITECTURE.md` §2.1: **REWRITE** de las descripciones de `review`, `auditor`, `senior`,
  `knowledge-base`, `ship` y `tdd`.
- `CLAUDE.md` y `AGENTS.md` §2: **REWRITE** de las mismas líneas (son espejo).
- `README.md`: **REWRITE** de la lista de toggles.
- `CHANGES.md`: **APPEND** del epic nuevo (ya hecho en este commit).

---

## Roles convocados y no convocados

**JUICIO** (sub-agente de propose, confianza alta: leí los fragmentos y los tests que nombro)

- Este documento lo escribió un solo sub-agente. **No convoqué al arquitecto**, que el bloque pide
  siempre en propose, ni a QA para decir cómo se probaría falso el resultado: eso lo decide el
  orquestador antes de que Juani firme. Lo digo para que no se lea como revisado.
- Backend, frontend, DBA, DevOps: no aplican a una propuesta de texto; los pasos 0, 8, 9, 10 y 11
  tocan código e instalación y los convocan en su propio propose.

**SUPUESTO**: que los números de línea citados de los tests y de `compose.go` siguen iguales al
construir; son del árbol en `4213a3f`.

---

## Aprobación

La completa Juani, a mano. Ningún agente la escribe.

| Elección | Persona | Fecha | Alcance |
|---|---|---|---|
| | | | |
