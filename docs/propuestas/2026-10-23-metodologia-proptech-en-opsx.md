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

**JUICIO** (recomendación provisoria, no decisión atribuida a Juani): el senior va **entre** el
rompedor y el auditor, para que el auditor decida leyendo también qué hallazgos contradicen una
decisión vieja. Se adopta como recomendación provisoria para fase 1 (§8).

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
| `review.md` (toggle `review`, verify) | tres pasadas por lentes, reportan, no bloquean, `review-verdict.md` | **Rompedor + bucle.** Conformidad primero, como requisito y sin habilitar archive (§3.1). El rompedor rompe con **reproducciones concretas** (comando y salida) y deja cada falla como test que falla o mutante que sobrevive; un agente cubre las tres lentes obligatorias (§3.3). Sección nueva "The loop": vueltas numeradas, tope 3, qué vuelve a apply, la salida por aprobación del auditor o por la vuelta 3 con decisión escrita de la persona, y L-029. "It reports, it does not block" se mantiene: **el rompedor reporta, decide el auditor, y quien frena es el orquestador**. Clave de Engram por vuelta y rol (`opsx/{change}/verify/round-{n}/{role}`), que el propio fragmento ya pide. Veredictos por vuelta en el change (`verify-round-{n}.md`). |
| `auditor.md` (toggle `auditor`) | lee informes, habla textual a la persona, puede callar | **Decide.** Su veredicto es APROBADO o BLOQUEANTES (lista) + DEUDA (lista) + LECCIONES (o "ninguna"). Sigue reportando a la persona, textual. "Puede callar" pasa a "si todo se sostiene, una línea: APROBADO". No nombra el bucle (va antes de `review` en la composición). |
| `senior.md` (toggle `senior`) | "convened in **every phase**", tabla de cinco fases, "five per change" | **Sólo verify.** Preámbulo, tabla y costo pasan a una fase: separa hallazgos reales de los que contradicen una decisión fechada. Uno por vuelta de verify (máx. 3 por change). Se mantienen el umbral BLOQUEANTE/NOTA y los canales de memoria. Se elimina "the senior labels, you hold": el senior aporta contexto, el auditor clasifica y decide (§3.1). |
| `knowledge-base.md` (toggle `knowledge-base`, archive) | REWRITE + APPEND + PENDING, reporte como condición de cierre | **Toda la base que el change tocó**: lo que nombra `## Knowledge Base Impact` **más** lo que muestra el diff. APPEND tiene tres clases: decisiones, **lecciones L-NNN** (sólo se agregan; salen de la sección LECCIONES del auditor) y **deuda** (sale de la lista DEUDA del auditor). **Cada entrada lleva su evidencia**: qué cambió y qué lo respalda (SHA, test o decisión), tipada. El detector cambia de "if you only appended, you did not update" a una etiqueta en el reporte: **SUSPICIOUS: append-only**. |
| `ship.md` + `assets/commands/claude/opsx/ship.md` | termina en el PR abierto; `main` pregunta siempre | `main` **sigue preguntando siempre** (sin cambio). Nuevo: nada está listo hasta que el CI completo esté verde sobre el SHA pusheado (L-037): `gh run list --commit <sha>` y `gh run watch`; estados y conjunto requerido según §3.4, sin interpretar incompletos como verde. "The cycle ends at the open PR" pasa a "...at the open PR **and its CI result**". Se cambian las dos copias. |
| `assets/commands/claude/opsx/verify.md` | conformidad sola | Reemplazar también "Ready for archive" por "conformidad evaluada, pendiente de decisión final". Paso 9: seguir §3.1 y §3.2; sólo el auditor habilita archive en el ciclo completo. La regla vive en `review.md`; el comando sólo la encadena. |
| `governance.md` | "explicit human approval", sin artefacto | Tabla `## Aprobación` en `proposal.md` (elección, persona, fecha, alcance), que ningún agente completa; una regla nueva después de la firma lleva fila propia (L-029). La verificación criptográfica queda para el paso 8. |
| Base, tabla de modelos (8 variantes) | archive `haiku` | archive `sonnet`: reescribir la base con evidencia no es "file operations". Pregunta 5. |
| `roles.md`, base Dispatch punto 2 | apply no convoca; lo independiente sale junto | **Sin cambio.** El rompedor vive en verify, que siempre lee un diff quieto; constructor → rompedor → auditor dependen uno del resultado del otro, así que la prueba de dependencia del punto 2 ya los manda en secuencia. |
| `roles.md` (toggle `roles`), fila `propose` y sección nueva tras "The floor that is not delegated" | propose convoca siempre al arquitecto y al resto "según lo que el change nombra"; nada liga el nivel de governance con cuántos roles leen el plan | **Regla de Juani del 23/10: en ALTO o CRÍTICO el propose lo analizan varios roles antes de la firma.** Varios = al menos tres informes independientes (arquitecto, analista funcional y QA, más los de las capas que el plan nombra), convocados de a uno, sin ver los informes de los otros, y después el auditor sobre todos los N informes (§3.5). Los informes se guardan textuales en el change y la tabla `## Aprobación` no se completa hasta que existan; si falta uno, el orquestador dice cuál y por qué. Vive en `roles.md` y no en `governance.md` porque es una regla de a quién se convoca; nombra a governance en forma condicional ("where governance is active"), igual que hoy lo hace el piso. **No es una aprobación**: los revisores no reemplazan la firma humana, la preceden. |
| `compose.go`, catálogo | 12 toggles | **Sin toggle nuevo.** Sí cambian las descripciones de `senior` y `review` en `harnesses.yaml`, `CLAUDE.md`/`AGENTS.md` §2, `ARCHITECTURE.md` §2.1 y `README.md`. |

### 3.1. Una sola autoridad de cierre

El orquestador ejecuta conformidad → rompedor → senior (si está activo) → auditor. Conformidad es requisito: un CRITICAL sin resolver impide APROBADO. Su resultado nunca dice "listo para archivar" ni "Ready for archive". El senior aporta contexto antes del auditor; no dispara una espera humana paralela. Si discrepan, el auditor debe clasificar el hallazgo y justificar por evidencia o decisión fechada por qué sigue como bloqueante o pasa a deuda. El orquestador ejecuta ese veredicto, no lo sustituye. Sólo APROBADO del auditor, con conformidad satisfecha e informes completos, habilita archive en el ciclo completo.

El paso 4 cambia juntos el cierre de conformidad, la decisión del auditor y las frases del senior que hoy ordenan ejecución conjunta o espera inmediata. El paso 5 termina de limitar al senior a verify y actualiza su costo y documentación. Así no queda una transición con dos autoridades.

### 3.2. Toggles opcionales y degradación visible

No se agregan dependencias obligatorias al catálogo. El orquestador anuncia al iniciar verify qué toggles faltan y registra la degradación en `tasks.md`:

| Falta | Regla |
|---|---|
| `auditor` con `review` activo | El rompedor informa, pero no hay aprobación automática ni retorno automático a apply. Se detiene el cierre y se pide a Juani decisión escrita en `tasks.md` para habilitar auditor o aceptar explícitamente el alcance degradado. Nadie se presenta como auditor sustituto. |
| `tdd` | Apply conserva su comportamiento habitual; no se afirma TDD estricto ni rojo/verde acreditado. Verify mantiene su bucle y evaluación de evidencia; la ausencia se declara al auditor. |
| `senior` | Se omite ese informe, se avisa y el auditor decide con conformidad y rompedor. No se inventa contexto histórico. |
| `knowledge-base` | No se promete actualización automática de KB. El veredicto conserva LECCIONES y DEUDA en el change y señala su destino pendiente a Juani. |
| `auditor` con `knowledge-base` activo | El archivista puede actualizar hechos con evidencia, pero no inventa LECCIONES/DEUDA del auditor: marca esa fuente ausente. El cierre requiere la decisión escrita de Juani sobre el alcance degradado. |
| `review` | No hay rompedor ni bucle de tres vueltas. Se declara el verify habitual; con auditor activo éste decide sobre los informes disponibles. Sin auditor, el orquestador informa conformidad y pide decisión escrita de Juani, sin anunciar aprobación del ciclo completo. |

Informe requerido que no volvió o no se puede recuperar íntegro no equivale a toggle apagado: frena el cierre, se identifica el faltante y se recupera o escala a Juani. No se rellena con un resumen ni con APROBADO supuesto.

### 3.3. Vueltas, cobertura y pruebas del contrato

La fuente de verdad es `tasks.md` del change: sección `## Bucle de verify`, con vuelta 1, 2 o 3, estado (en curso/cerrada), SHA o diff revisado, rutas de informes, bloqueantes y decisión. Se escribe antes de convocar revisores. Una vuelta comienza con verify sobre un diff quieto y termina con el veredicto del auditor. Sólo después de corregir los bloqueantes de vuelta 1 o 2 se incrementa al siguiente verify; recuperar sesión o repetir una captura fallida no incrementa ni reinicia el contador. Al retomar se lee esa sección, no Engram como autoridad.

La tercera vuelta cierra el bucle incluso si quedan problemas. No hay cuarta automática: lo abierto queda como decisión escrita para Juani en `tasks.md`, con opciones corregir en un nuevo alcance aprobado, aceptar deuda justificada o detener el change. Sin respuesta queda detenido. La aceptación humana no convierte un CRITICAL pendiente en conformidad satisfecha ni fabrica APROBADO del auditor.

Un rompedor entrega tres secciones con trabajo y resultado, aun si dice "sin hallazgos": corrección (casos límite y reproducciones), wiring (consumidores y conexiones) y comentarios/documentación (contratos frente a conducta). En vueltas 2 y 3 el constructor corrige la lista del auditor; el rompedor revisa además consumidores y regresiones de esas correcciones fuera del diff. Lo nuevo se reporta al auditor: si bloquea, entra al veredicto; lo demás va a deuda asignada al primer change que use la pieza. No se oculta una regresión por estar fuera de las líneas cambiadas.

En apply se agregan pruebas textuales positivas de instrucciones y resultado por lente, no sólo títulos ni antiasserts. Mutantes: quitar wiring o comentarios, aprobar con CRITICAL, anunciar archive antes del auditor, ejecutar auditor antes del retorno del senior, reiniciar contador y permitir cuarta vuelta. Matriz mínima: todos activos; cada ausencia de §3.2; KB sin auditor; informe faltante; APROBADO; bloqueantes en vueltas 1/2; tercera con y sin respuesta. Deben fallar antes del cambio correspondiente. Los goldens se revisan por diff después de esas pruebas; no prueban conducta real de agentes.

### 3.4. CI completo, igual en ambas copias

"CI en verde" significa **todos los jobs requeridos del workflow en `success` sobre el SHA pusheado**. El orquestador registra el SHA, workflows/eventos esperados y jobs expandidos de matriz antes de evaluar resultados. En este repo el mínimo es `.github/workflows/ci.yml`, evento push, jobs `test` en Ubuntu y Windows; si hay PR se inspeccionan también sus checks requeridos, verificando el SHA de origen, sin confundirlo con el SHA sintético de merge. Release por tag no se exige en un push de rama. En otros proyectos se deriva el conjunto de sus workflows y checks requeridos, sin elegir sólo el run favorable.

`gh run list --commit <sha>` descubre runs y `gh run watch` espera su final; se inspeccionan además los jobs del conjunto esperado. Si todavía no apareció un run, sigue pendiente. En curso, queued, cancelado, skipped, neutral, timeout, sin runs, sin CI o consulta fallida nunca son verde: se informa estado y se deja sin confirmar. Fallo se reporta sin apilar otro arreglo como si el push estuviera confirmado. Un nuevo push invalida el resultado anterior y exige evaluar el nuevo SHA. `main` sigue preguntando siempre; nunca se mergea.

Se exige un test propio para `ship.md` compuesto (`ship_toggle_test.go`) y otro para `assets/commands/claude/opsx/ship.md` (`internal/harness/command` o `assets/commands_test.go`), aunque la ruta no cambie. Cada uno fija conjunto completo, SHA y `success`; casos negativos: un job verde y otro pendiente, cancelado, omitido, otro SHA, cero runs y error de consulta. Mutantes que acepten cualquiera de esos estados deben dar rojo en ambas copias. Son pruebas del contrato textual; la observación de CI real queda registrada en apply cuando haya un push autorizado.

### 3.5. Informes y reglas retenidas

Para propose ALTO/CRÍTICO hay **N informes especialistas + 1 del auditor**: N ≥ 3, siempre arquitecto, analista y QA, más un informe por cada especialidad adicional que el alcance requiera. El auditor recibe todos los N originales, no "los cuatro" por regla fija. Para este propose N = 4 por DevOps: cinco informes en total.

Para cada vuelta completa de verify hay cuatro informes: conformidad, rompedor (tres lentes en un informe), senior y auditor. Sin senior son tres; las otras degradaciones se registran como §3.2. El constructor deja un informe de apply por entrega, separado de esos cuatro. Se guardan en `openspec/changes/<change>/informes/`, con fase, vuelta y rol en el nombre; `LEEME.md` indexa quién, modelo/proveedor declarado, SHA, procedencia y faltantes. `verify-round-{n}.md` referencia los originales y veredicto, no los reemplaza. Esta revisión documental usa excepcionalmente las rutas versionadas de §Roles para preservar los cinco originales pedidos.

Hasta tener hook, el orquestador conserva el retorno íntegro o copia textual del transcript y declara esa mediación. Los revisores no reciben Write. Los especialistas de propose se convocan de a uno con propuesta y fuentes, excluyendo informes y memorias de los otros. Cada informe declara qué vio, si recibió casos dictados por el auditor y si hubo contaminación; si leyó otro informe, se conserva pero no se cuenta como independiente y se obtiene reemplazo si falta el mínimo. Rompedor que verifica un caso dictado lo declara, sin vender coincidencia como convergencia independiente (R14, R16). En este lote sólo hay declaraciones de aislamiento, no prueba de ausencia de contaminación.

Trazabilidad retenida de la exploración: R7/R8 van a governance (aprobación humana y regla nueva); R9/R10/R12 a strict-tdd; R6/R21 al brief. R11 exige que rompedor y auditor cuenten tests por sí mismos, sin aceptar el conteo del constructor. R13 exige que constructor y rompedor lean las lecciones antes de actuar, registrando ruta o ausencia en el informe; vive en strict-tdd y review, condicional a KB disponible. R17/L-032 vive en review: si hay proveedor externo de identidad, se pide recorrido real con Juani al teclado antes de archive; lo no recorrido se registra como deuda explícita y limitación para el auditor, nunca como validado. R18 va a ambas copias de ship. R15 se conserva con revisión de correcciones y sus regresiones (§3.3); R22 usa almacenamiento manual en el change durante fase 1, sin comprometer aún el destino del hook. R23 queda excluida por específica de Codex sin sandbox, como decidió la exploración. Los tests de apply fijan lectura previa, conteo propio, identidad no recorrida y contaminación declarada.

### 3.6. Entrega y recuperación en la máquina de Juani

El paso 0 ejecuta en apply, no durante este propose, esta receta registrada en `tasks.md`:

1. Desde checkout limpio del SHA aprobado, registrar versión, ejecutable efectivo en PATH, destinos resueltos por adapter, variante y toggles efectivos. Inspeccionar también `/usr/local/bin` para aislar ensayos; `--home` solo no garantiza aislamiento. No atribuir el `.test` a un test concreto sin reproducción.
2. Antes de escribir, snapshot vía `internal/backup` de cada archivo afectado: binario previo, instrucciones, comandos y settings que se modifiquen. Registrar ID y ruta absoluta del snapshot, fecha, SHA y manifiesto con destino, existencia previa, hash de contenido y permisos. Si el backup no se puede leer o está incompleto, no instalar.
3. Ensayar primero en destinos aislados con binario y PATH aislados: instalar, reinstalar sin bloques duplicados, inducir falla y restaurar. El restore repone contenido y permisos del manifiesto y elimina sólo archivos nuevos de esa operación; compara hashes y permisos con el estado previo. Registrar errores y archivos no recuperados; no declarar restore exitoso si difieren.
4. Recién con ese ensayo satisfactorio, reinstalar en la máquina de Juani con snapshot propio identificable. Ante error de escritura, fallo de health check o diferencia de entrega, restaurar ese snapshot y detener el paso. La aceptación exige binario de `cmd/active-stack` (no `.test`), `roles --help` con salida 0 y contenido instalado igual a `Compose` para variante y toggles registrados, preservando texto del usuario fuera de markers.

Aceptación por destino: Claude, bloque y comandos verify/ship; OpenCode, agente primario y referencia en settings; Codex, instrucciones mediante su adapter sin inventar directorio de comandos. OpenCode debe recuperar instrucciones y settings conjuntamente y conservar ambos respaldos, sin colisión. Los demás adapters se comprueban según su entrega real, con escenarios aislados. Ningún agente depende del hook en fase 1: usa informes manuales íntegros de §3.5. Los goldens Claude/OpenCode no sustituyen estos checks de entrega. No se modifica ahora backup, filemerge ni rollback: si apply necesita hacerlo, requiere propuesta y review ALTO propios.


---

## 4. Tests y goldens que se mueven

| Paso | Tests |
|---|---|
| ship espera el CI | `ship_toggle_test.go`: `NeverMerges` (la frase "the cycle ends at the open pr"), lista nueva de lecturas; cuidar las frases prohibidas de `DoesNotClaimAMechanismItLacks` ("waits for confirmation", "asks and waits"): el texto dice que **el orquestador** espera el run, no que algo lo garantice. Goldens `compose_all_toggles_{claude,opencode}.golden`. Test separado del comando en `internal/harness/command`, aunque no cambie la ruta. Ambas copias deben satisfacer §3.4. |
| constructor (`strict-tdd.md`) | `tdd_toggle_test.go` **nuevo** (hoy no existe); `compose_test.go` `TestCompose_TDDToggleOff`; `dispatch_rules_test.go` `TestElRepartoNoAutorizaRecortarLaRevision` (no tocar sus frases); goldens. |
| rompedor + bucle (`review.md`, `auditor.md`, `verify.md`) | `review_toggle_test.go`: `ThreeIndependentPasses` y `IndependenceIsByConstruction...` cambian para un agente con tres lentes; nuevos para el bucle (tope 3, salida por aprobación, vuelta 3 como decisión escrita, L-029, clave por vuelta); `ReportsDoesNotBlock` se mantiene. `auditor_toggle_test.go`: nuevos (APROBADO/BLOQUEANTES/DEUDA/LECCIONES). `review_cost_rules_test.go` `TestNingunaReglaNuevaRecortaAlAuditorNiLasLentes` **no se toca**: guarda que las lentes no se recorten. `verify_phase_test.go` y `assets/commands_test.go` por el paso 9. Goldens. |
| senior sólo verify | `senior_toggle_test.go`: `PreambleSaysItIsEveryPhase`, `AsksADifferentQuestionInEachPhase`, `RunsAtTheCloseOfEachPhase`, `IsNotAnEighthRole`, `DeclaresItsCost` se reescriben; `InVerifyItJudgesTheFindings` se mantiene. `docs_toggles_test.go` (descripciones en docs y catálogo). Goldens. |
| archive con evidencia | `knowledge_base_toggle_test.go`: `TheDetectorIsTheRewriteList` (etiqueta SUSPICIOUS), `ArchivePromptCarriesTheReport` (el prompt pide lecciones, deuda y evidencia), nuevos para L-NNN y "base que el diff tocó". `verify_phase_test.go`: verificar archive a sonnet en las ocho variantes y quitar instrucciones contradictorias de menor esfuerzo (incluida Antigravity). Goldens. |
| governance | goldens; un test nuevo para la tabla `## Aprobación` y L-029. |
| propose multi-rol (`roles.md`) | `roles_toggle_test.go` nuevo `TestCompose_RolesToggle_HighProposeIsReadBySeveralRolesBeforeSignature`: con `roles`+`governance`, la fila `propose` y la sección nombran HIGH/CRITICAL, "before" la aprobación, al menos tres roles y que no sustituye la firma humana; se ve fallar contra el `roles.md` actual (su fila `propose` sólo dice `architect`). Otro test fija que el texto sea condicional (`roles` sin `governance` no afirma que exista un nivel HIGH). `TestCompose_RolesToggle_MapsTheOPSXPhases` y `DoesNotConveneEverybodyEveryTime` no se tocan: la regla es por nivel, no "todos siempre". Goldens. |

Regla para todos: **cada test nuevo se ve fallar antes** contra el fragmento viejo (el repo ya tropezó
siete veces con asserts que pasaban sin probar nada, según `review_toggle_test.go:213-217`).

---

## 5. Orden de construcción, en pasos chicos

Cada paso es un change OPSX de ≤ 12 tareas, con su propia firma si es ALTO.

| # | Change | Qué | Governance | Depende de |
|---|---|---|---|---|
| 0 | `instalacion-sana` | Encontrar qué test escribe un `.test` en `~/.local/bin`; `internal/verify` chequea que el binario sea `cmd/active-stack`, que `active-stack roles --help` salga 0 y que el bloque instalado sea igual a `Compose(toggles)`; subcomando `version`. Después, reinstalar en la máquina de Juani con la receta y aceptación de §3.6. Diagnosticar sin dar por probada la causa del binario `.test`. | MEDIO | ; |
| 1 | `opencode-mcp-command-array` | Cherry-pick de `fae8bbf` con resolución en `mcp_test.go`, sin el trailer de coautoría. | MEDIO | ; |
| 2 | `ship-espera-ci` | §3 fila ship, en las dos copias. | BAJO | ; |
| 3 | `apply-constructor` | §3 fila `strict-tdd.md`. | BAJO | ; |
| 4 | `verify-rompedor-auditor` | `review.md` (rompedor + bucle), `auditor.md` (decide), cierre y paso 9 de `verify.md`; corregir ya en `senior.md` la secuencia y la autoridad (§3.1), sin esperar al paso 5. | BAJO | 3 (el bucle nombra al constructor) |
| 5 | `senior-solo-en-verify` | `senior.md` + docs y catálogo. | BAJO | 4 (su lugar en la secuencia) |
| 6 | `archive-base-completa` | `knowledge-base.md` (lecciones, deuda, evidencia, SUSPICIOUS); archive a `sonnet` como recomendación provisoria (§8). | BAJO | 4 (LECCIONES y DEUDA salen del auditor) |
| 7 | `governance-aprobacion` | Tabla `## Aprobación` + L-029 en `governance.md`. | BAJO | 4 |
| 7b | `propose-multirol` | §3 fila `roles.md`: en ALTO/CRÍTICO, varios roles + auditor antes de la firma. | BAJO | 7 (la firma que precede es la tabla `## Aprobación`) |
| 8 | `aprobacion-firmada` | Subcomando `active-stack approval --change` con `git verify-commit`. Candidato, fuera de fase 1: TBD de §8. | ALTO | 7 |
| 9 | `informes-por-hook` | Investigar captura de informes. Fuera de fase 1: TBD de §8; hook, subcomando y destino aún no aprobados. | MEDIO | 4 |
| 10 | `opencode-permission-replace` | Reimplementar `ddeb94f`: `__replace__` sólo sobre `bash`. | ALTO | ; |
| 11 | `instalar-roles` | Investigar instalación de `role-*`. Fuera de fase 1: TBD de §8; tipo de harness y reemplazo aún no aprobados. | ALTO | 9 |
| 12 | `bloque-compacto` | Investigar compactación y destinos completos de todas las fases. Fuera de fase 1: TBD de §8, sin mover texto ahora. | MEDIO | todo lo anterior |

0, 1, 2 y 3 son independientes entre sí. 2 a 7, incluido 7b, son texto de fragmentos y tests de Go; 8 a 12 tocan
código, instalación o permisos.

---

## 6. Riesgos

**JUICIO** (sub-agente de propose; confianza entre paréntesis)

- **Apagar `review` apaga el bucle.** Con el método dentro de `review.md`, quien quiera el ciclo sin
  rompedor se queda también sin vueltas. Es el costo de no tener toggle propio (alta).
- **Un rompedor con tres lentes por vuelta**, hasta tres rompedores, tres auditores y tres seniors por change. Se pierde independencia entre lentes; se declara (alta).
- **La vuelta 3 sin respuesta.** Si la persona no contesta, nada avanza; el fragmento tiene que
  decirlo como default seguro del orquestador, igual que hoy dice `senior.md` (alta).
- **Los informes siguen pasando por el orquestador** hasta el paso 9: el auditor puede leer una
  copia, no el original (R14). Mientras tanto, se guarda copia textual del retorno o transcript disponible y se declara su procedencia; si no se puede recuperar íntegro, se informa faltante y se detiene el cierre (media).
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
`openspec init` y quedaríamos manteniendo una copia de algo ajeno. Las reglas viajan por el bloque. Revertir esta decisión requeriría mantener dos comandos propios, sus tests y su compatibilidad con OpenSpec; no se hace en fase 1.

---

## 8. Recomendaciones provisorias y límites de firma

Se mantienen las cinco recomendaciones del orquestador para fase 1, sujetas a la firma de Juani:

1. Un rompedor con tres lentes obligatorias. En vueltas 2 y 3 revisa correcciones, consumidores y regresiones, incluso fuera del diff. AN-2 y QA-4 justifican ampliar el recorte original; no se presenta como independencia entre lentes.
2. Senior antes del auditor, sin autoridad paralela de cierre (§3.1).
3. Firma criptográfica sólo en fase 2, si Juani la configura. En fase 1 vale la aprobación escrita, fechada y de alcance explícito de Juani en `## Aprobación`, en commit propio; si falta o es ambigua, el constructor frena. No acredita verificación criptográfica.
4. Modelo y proveedor declarados en cada informe, incluida familia compartida; sin registro del despacho se tipan SUPUESTO. No se agrega control de proveedor a la herramienta.
5. Archive a sonnet con sustitución declarada si no está disponible. No garantiza calidad. Haiku queda como alternativa; revertir requiere cambiar ocho bases, tests y dos goldens.

**El alcance propuesto para firma abarca fase 1, pasos 0 a 7, incluido 7b.** Esta delimitación precisa el roadmap general de `CHANGES.md`, que no autoriza los diseños posteriores. Los pasos 8 a 12 son intenciones de roadmap, no diseños aprobados ni autorización para implementarlos:

- Firma, paso 8: **(TBD: se decide en el propose de la fase 2 o 3)**. Requiere configurar identidad autorizada y vínculo entre firma, alcance y commit; `approval --change` y `git verify-commit` son candidatos, no contratos cerrados.
- Hook e instalación de roles, pasos 9 y 11: **(TBD: se decide en el propose de la fase 2 o 3)**. Falta validar captura íntegra, persistencia sin sobrescritura, worktrees, fallas visibles y convivencia con hooks; tampoco se aprobó un tipo `agent` ni pisar roles personales.
- Compactar, paso 12: **(TBD: se decide en el propose de la fase 2 o 3)**. Antes de mover texto hay que definir destinos y entrega completa al constructor y archivista; fase 1 conserva apply/archive en el bloque y los transmite íntegros en sus prompts.
- Permisos OpenCode, paso 10: **(TBD: se decide en el propose de la fase 2 o 3)**. Necesita su revisión ALTO antes de escribir configuración.
- Shape y navegador para toda UI o sólo impeccable: **(TBD: se decide en el propose de la fase 2 o 3)**. La fase 1 no cambia UI.

Cada propose posterior debe fijar ubicación, consumidor y proveedor, alternativa descartada y costo de reversión antes de firmar. No se toma ahora una decisión de frontera que fase 1 no necesita.

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

La primera versión la escribió un constructor sin revisores. Después se recibieron cuatro informes de especialistas y uno del auditor, preservados textuales en `docs/propuestas/informes/2026-10-23/`: analista funcional (01), arquitecto (02), DevOps (03), QA (04) y auditor (05). DevOps sí aplica por instalación y recuperación. Esta corrección la hace el constructor leyendo los cinco, por lo que no es una nueva revisión independiente.

Los cuatro especialistas declararon no consultar los otros informes. El auditor leyó los cuatro. Eso no prueba aislamiento: independencia SUPUESTA, contaminación no acreditada ni descartada. No se afirma una nueva aprobación del auditor sobre esta corrección. Los ensayos propuestos se ejecutan en apply, no se dan por aprobados por esta revisión documental.

## Aprobación

La completa Juani, a mano. Ningún agente la escribe.

| Elección | Persona | Fecha | Alcance |
|---|---|---|---|
| Fase 1 del propose (pasos 0 a 7, incluido 7b) con las 5 recomendaciones del orquestador | Juani | 23/10/2026 | Tal como está en el propose en el commit 409d04a, revisado por analista, arquitecto, DevOps y QA y auditado en dos vueltas (docs/propuestas/informes/2026-10-23/). Los pasos 8 a 12 quedan TBD y NO aprobados. Cada paso es un change con su propia firma si es ALTO |


## Respuesta a la revisión

| Bloqueante deduplicado | IDs de origen | Cómo quedó para fase 1 |
|---|---|---|
| Autoridad y cierre único | AN-1, AN-4, AR-2, QA-3 | §3.1: conformidad requerida, senior aporta antes, auditor decide; paso 4 corrige juntas las contradicciones y elimina el anuncio prematuro de archive. |
| Dependencias opcionales | AN-3, AR-1, QA-1 | §3.2: degradación anunciada por toggle, sin aprobación inventada; informe faltante frena y ausencia de auditor requiere decisión escrita de Juani. Matriz de pruebas en §3.3. |
| Vueltas y cobertura | AN-2, QA-4 | §3.3: contador recuperable en tasks.md, tercera cierra y escala lo abierto a Juani; tres lentes con resultados y regresiones fuera del diff, con mutantes positivos. |
| CI completo y dos copias | AN-6, OPS-2, QA-2 | §3.4: todos los jobs requeridos en success sobre SHA pusheado; cancelado o en curso no es verde. Un test por copia, incluidos estados incompletos y errores. |
| Trazabilidad e independencia | AN-5, AN-7 | §3.5: N especialistas + auditor, rutas y destinatarios, contaminación declarada; lectura previa de lecciones, conteo propio e identidad al teclado con deuda explícita. AN-5 se resuelve como omisión de contrato, sin convertir su búsqueda negativa en prueba de ausencia semántica. |
| Entrega y recuperación | OPS-1, OPS-4 | §3.6: snapshot identificable, manifiesto y restore por contenido/permisos, aceptación por adapter y reinstalación repetida. Se ejecutan en apply; no se afirma recuperación ya probada. |
| Fronteras del alcance | AR-3, AR-4, OPS-3 | §8: firma, hook, roles y compactación excluidos de la firma de fase 1 con TBD y motivo. Apply/archive siguen en el bloque; los contratos, alternativas y reversión futuros requieren propose propio. |

Se conserva el criterio de las cinco recomendaciones provisorias. Las razones fuertes para ajustar su redacción son AN-1/AR-2/QA-3 (autoridad y transición) y AN-2/QA-4 (regresiones fuera del recorte); no justifican sumar infraestructura en fase 1. Las notas AR-5/QA-5 se cubren con revisión de goldens y test de ocho variantes; OPS-5 con aislamiento de PATH y binario. La tabla de aprobación sigue vacía.
