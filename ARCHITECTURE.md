# Active Stack — Arquitectura vigente

> Instalador **methodology-first** del harness de desarrollo asistido por IA.
> Materializa la metodología OPSX ([docs/metodologia-opsx.md](docs/metodologia-opsx.md), [docs/lean-agent-workflow.md](docs/lean-agent-workflow.md)): un comando instala y configura el
> sustrato (harnesses), y un único orquestador de fundación deja el proyecto
> listo para el ciclo OPSX (`explore → propose → apply → verify → archive`).

Estado actual del proyecto: [docs/ESTADO.md](docs/ESTADO.md).

---

## 1. Qué es (y qué NO es)

**Es:** un configurador del harness metodológico. Toma el/los agente(s) de IA
del usuario y les inyecta, de forma modular y actualizable, los harnesses que
exige la metodología.

**NO es** (decisión: sacar del scope, no van en la metodología):

- ❌ GGA (code review en commit) — harness aparte, no parte del ciclo SDD.
- ❌ Themes / statusline / keybindings — cosmético.
- ❌ Persona como producto de marketing ("Your own Gentleman!").
- ❌ Framing "supercharge any agent" / relación con Gentleman.Dots.

La persona puede sobrevivir como harness **opcional** (configuración), sin el
envoltorio de marketing.

### 1.1 Mantener / Modificar / Sacar

El repo es nuevo por el volumen de cambios, pero NO se tira todo. Buckets:

| Bucket | Qué | Detalle |
|--------|-----|---------|
| **MANTENER** | Núcleo del instalador e interfaz | Detección OS/deps, adapters por agente, backup/rollback, merge por markers, pipeline, verify, y la **interfaz de install/uninstall** de la TUI. |
| **MODIFICAR** | Flujo de instalación + branding | El **flujo de install se rediseña**: antes era selección de componentes por agente; ahora hay **muchos más módulos (harnesses)** → selección/agrupación por harness, árbol de deps entre harnesses, modos. Limpiar todo leftover `gentle-ai` / `gentle-stack` / `Gentleman.Dots`. |
| **SACAR** | Lo no-metodológico | GGA, theme, persona-marketing, statusline, keybindings, framing "supercharge any agent". |

> Nota: el flujo de install **no es un port directo** — es un rediseño. La
> mecánica de bajo nivel (merge, backup, adapters) se reusa; la **orquestación
> de qué se instala y cómo se presenta** se rehace alrededor del catálogo de
> harnesses (§3). El uninstall se mantiene pero debe entender harnesses, no
> componentes viejos.

---

## 2. Definición de Harness

Un **harness** es cualquier módulo que prepara o guía el entorno de la IA.
**No es sinónimo de "repo".** Se materializa de cuatro formas:

| Forma | Qué es | Origen | Ejemplos |
|-------|--------|--------|----------|
| **Skill** | `SKILL.md` + assets, cargada bajo demanda | repo propio o de terceros, se baja al instalar | kb-creator, roadmap-generator, agent-instruction, active-orchestrator, skill-registry (propias); find-skill, skill-creator (terceros) |
| **Configuración** | Texto/archivos que configuran el agente | bundleado en el instalador | `sdd-orchestrator` (orquestador con toggles: TDD, Engram, …), config de `AGENTS.md`/`CLAUDE.md`, permisos, MCP |
| **Herramienta externa** | Binario/servicio de terceros | instalado/configurado, no es nuestro repo | Engram, OpenSpec CLI, Context7 (MCP) |
| **Comando** | Slash-commands nativos del agente | bundleado en el instalador | `/active:init`, `/active:starter-add`, `/opsx:ship`, `/opsx:verify` |

Principio: el instalador no instala "repos", instala **harnesses**, y cada
harness sabe cómo se materializa.

### 2.1 El harness `sdd-orchestrator` (clave)

El modo predeterminado es `lean-workflow`: core compacto con contratos universales y cinco roles secuenciales. En Claude cada rol corre como subagente nativo `role-*`; para otros ejecutores el encargo lleva el contrato universal. Los modelos/proveedores son configurables y los plugins de ejecución quedan fuera del instalador. Ver [flujo compacto de cinco roles](docs/lean-agent-workflow.md) para matriz de precedencia, TDD y estado durable.

Toggles del `sdd-orchestrator` (lista canónica: `config.KnownToggles()`; el orden no cambia el resultado). Con `lean-workflow` activo sólo `governance`, `evidence`, `knowledge-base`, `engram` y `ship` agregan su contrato compacto; los marcados **sólo legacy** se aceptan por compatibilidad pero no pueden apagar nada del flujo lean (TDD y auditor son siempre obligatorios). Se eligen con `active-stack install --toggles`:

- `lean-workflow` (selector) — selector predeterminado; quitarlo activa la composición legacy.
- `delegation` (sólo legacy) — reglas de delegación: qué se resuelve inline y qué se manda a
  un sub-agente, para no inflar el contexto.
- `model-routing` (sólo legacy) — tabla de modelo por fase, con sustitución declarada cuando
  el modelo asignado no está disponible.
- `governance` (lean y legacy) — niveles de autonomía por criticidad del dominio (CRITICAL,
  HIGH, MEDIUM, LOW): a más crítico, menos se escribe sin aprobación humana.
- `evidence` (lean y legacy) — reporte tipado: cada afirmación etiquetada EVIDENCIA, JUICIO o
  SUPOSICIÓN, y la evidencia exige comando, salida y hora.
- `auditor` (sólo legacy) — alguien lee a los revisores: marca evidencia sin comando, salida
  vieja y suposiciones ascendidas a evidencia.
- `roles` (sólo legacy) — convoca un rol con estándar fijo y un "listo" observable, en vez de
  fabricar un sub-agente genérico distinto cada vez.
- `review` (sólo legacy): un rompedor con tres lentes obligatorias y bucle de hasta tres vueltas; el auditor decide y el orquestador ejecuta.
- `knowledge-base` (lean y legacy): en archive reescribe la base que impacto y diff tocaron; agrega decisiones, lecciones L-NNN y deuda del auditor con evidencia por entrada. Sólo APPEND sin REWRITE se reporta SUSPICIOUS: append-only.
- `senior` (sólo legacy): sólo en verify, después del rompedor y antes del auditor; aporta contexto de decisiones fechadas, uno por vuelta y hasta tres por change.
- `ship` (lean y legacy) — publicar en GitHub el trabajo de un change cerrado. No es una
  secuencia de pasos sino un diagnóstico: mira en qué estado está el
  repositorio (trunk, rama sin remoto, PR abierto, trunk que avanzó, cambios
  sin commitear, sin remoto) y decide. Push a `main` siempre pregunta, a una
  rama propia no; PR o directo se propone con fundamento; espera el CI completo sobre el SHA pusheado y nunca mergea.
- `engram` (lean y legacy) — memoria persistente entre sesiones y compactaciones, cableada a
  las fases OPSX.
- `tdd` (sólo legacy) — en lean TDD es siempre obligatorio y este toggle no cambia nada; en legacy: TDD estricto en la fase apply: no hay código de producción sin un
  test que falle primero.

El resultado es el bloque de instrucciones del orquestador que se inyecta en
`CLAUDE.md` / `AGENTS.md`, compuesto a partir de los toggles elegidos.

---

## 3. Catálogo de harnesses (embebido)

El instalador lleva embebido un **catálogo maestro** (p. ej. `harnesses.yaml`)
que es la fuente de verdad de qué harnesses existen. Cada entrada declara:

```yaml
- id: kb-creator
  type: skill            # skill | config | external | command
  source:
    repo: Group-Active-IA/kb-creator
    ref: latest          # tag/branch/commit
  install_modes: [full]  # lite | full | ...
  depends_on: [openspec] # otros harness id
  agents: [claude, opencode, ...]   # dónde aplica
- id: sdd-orchestrator
  type: config
  toggles: [lean-workflow, delegation, model-routing, governance, evidence, auditor, roles, review, knowledge-base, senior, ship, engram, tdd]
  install_modes: [lite, full]
  ...
- id: engram
  type: external
  install: { method: homebrew | download, ... }
```

La TUI lee este catálogo para construir la lista de selección. Para
agregar/actualizar un harness se toca el catálogo y se saca release del
instalador (tradeoff aceptado: simple y versionado con el instalador).

---

## 4. Flujo de uso

### 4.1 Instalación (una vez por máquina)

```
active-stack install  →  TUI (Bubbletea)
  1. Detectar OS/arch/agentes/deps
  2. Elegir agente(s)
  3. Elegir modo (Lite / Full / Custom)
  4. Resolver catálogo → árbol de dependencias de harnesses
  5. Backup de configs existentes
  6. Instalar deps externas (Engram, OpenSpec) + harnesses
  7. Inyectar configs por agente (merge por markers)
  8. Verificar (health checks)
```

### 4.2 Fundación del proyecto (una vez por proyecto)

Un **único comando** que orquesta la fase de fundación con lazy-loading de
skills, en este orden (la skill `active-orchestrator`, antes llamada `jr-orchestrator`).
`discovery-research` corre antes de `kb-creator` si se elige investigar primero. `project-scaffold`
no forma parte de este flujo: es un generador de documentación independiente; dentro de la
fundación, el `AGENTS.md`/`CLAUDE.md` canónico lo genera `agent-instruction`:

```
1. openspec init        (solo la carpeta openspec/)
2. kb-creator           (armar knowledge-base/)
3. roadmap-generator    (armar CHANGES.md)
4. find-skill           (instalar skills faltantes según stack)
5. agent-instruction    (armar AGENTS.md/CLAUDE.md con todas las referencias)
```

### 4.3 Ciclo iterativo (por change)

Con lean activo, el ciclo OpenSpec conserva sus artefactos y el agente principal coordina planificador → revisor del plan → constructor → rompedor → auditor. Los cambios se clasifican por riesgo y alcance; el constructor aplica TDD para código siempre. Los informes originales y el estado durable permiten retomar y repetir verificaciones afectadas. Se conserva la fundación KB/roadmap y la memoria del proyecto. Ver [flujo compacto de cinco roles](docs/lean-agent-workflow.md).

La siguiente descripción corresponde **solo al modo legacy**, sin `lean-workflow`, y se conserva por compatibilidad:

- **Propose:** con governance activo, ALTO o CRÍTICO requiere informes de arquitecto, analista funcional y QA, más especialistas por capa, convocados de a uno sin leer los otros informes; el auditor lee los N originales antes de la aprobación humana escrita. La tabla `## Aprobación` registra elección, persona, fecha y alcance; ningún agente la completa. Una regla nueva necesita aprobación propia.
- **Apply constructor:** rojo y verde por tarea en `tasks.md`, con comando, salida, hora y SHA; si el test nació verde, mutante que falle en su aserción y se revierta. Declara consumidores, conflictos y familia del modelo. En vueltas 2 y 3 corrige sólo la lista del auditor.
- **Verify:** conformidad → rompedor con tres lentes (correctitud, cableado, documentación) → senior si está activo → auditor. La conformidad sola no habilita archive. El auditor decide y el orquestador ejecuta; sólo APROBADO con informes completos y conformidad satisfecha habilita el cierre ordinario. El contador y los informes viven en el change: máximo tres vueltas, sin cuarta automática. En la tercera con bloqueantes se exige decisión humana escrita sobre alcance nuevo, deuda justificada o freno; no se fabrica APROBADO.
- **Archive:** lee la unión del impacto declarado y el diff; reescribe las descripciones actuales y agrega decisiones, lecciones numeradas y deuda del auditor con fecha, costo y evidencia tipada. Deja `knowledge-base-report.md` por change con REWRITTEN, APPENDED y PENDING por documento; sólo APPEND se justifica y se señala como `SUSPICIOUS: append-only` cuando correspondía reescribir.
- **Ship:** espera todo el CI requerido sobre el SHA pusheado, inspeccionando workflows, eventos y jobs. En active-stack exige `ci.yml` en push con tests Ubuntu y Windows; en otros repos deriva su conjunto. Cero runs, consulta fallida, SHA ajeno o estados incompletos no confirman entrega. `main` exige permiso y nunca se mergea el PR automáticamente.

Si falta un toggle se declara la degradación; sin review no hay rompedor ni bucle. Sin auditor se necesita decisión escrita de la persona responsable para cerrar con alcance degradado. Un informe obligatorio faltante detiene el cierre.

EVIDENCIA: [contratos, decisiones y deuda de fase 1](docs/metodologia-opsx.md). OpenSpec CLI conserva el estado; no actualiza por sí solo esta base.

---

## 5. Arquitectura del instalador (Go + Bubbletea)

La infraestructura implementada se organiza así:

```
cmd/active-stack/            entrypoint CLI (install, uninstall, roles, reports, starter, version)
  headless/              ejecución no interactiva de install/uninstall/starter
internal/
  system/                detección OS/arch/WSL/Termux, deps, guards
  catalog/               parseo del harnesses.yaml embebido
  model/                 tipos de dominio (harness, agente, modo)
  planner/               grafo de dependencias, orden, review payload
  agents/                adapters por agente (claude/opencode/codex/...)
  harness/               install/inject por tipo de harness
    skill/  config/  external/  command/  roledefs/
  install/               plan de instalación, snapshot y self-install
  uninstall/             desinstalación targeted y restore desde backup
  filemerge/             merge por markers (inyectar sin pisar)
  backup/                snapshot + restore de configs
  pipeline/              ejecución por etapas + rollback
  verify/                health checks post-install
  reports/               captura de informes por proveedor
  roles/  approval/      clasificador de roles legacy y tabla de aprobación
  tui/                   Bubbletea (sin theme cosmético del viejo)
assets/                  slash-commands y skills embebidas (los fragmentos del
                         orquestador viven en internal/harness/config/assets/)
```

Procedencia histórica: port del repo anterior. Los paquetes listados ya existen.
Del repo viejo se descartó GGA, theme y persona (marketing); permissions se mantiene
como harness `config` (la seguridad no es opcional).

---

## 6. Roadmap de construcción (incrementos)

El roadmap de construcción C-01…C-32 está en `CHANGES.md`. Estado actual del proyecto: [docs/ESTADO.md](docs/ESTADO.md).

EVIDENCIA: propuesta firmada `4748680`, decisión `5ab0f4f`; estado detallado en `CHANGES.md`.

---

## 7. Decisiones tomadas y abiertas

### Resueltas
- **Fetch de skills (C-08) — mixto por tipo**: skills propias (active-orchestrator, etc.)
  → `git clone` del repo; terceros (find-skill, skill-creator) → también `git clone`, con `path` al subdirectorio;
  core openspec (atadas a la versión del orquestador) → embebidas. El modelo de
  skill necesita un campo de método de fetch (clone/npx/embed), análogo a
  `External.Method`.
- **active-orchestrator — activador modular**: orquesta la fase de fundación y activa/llama
  los módulos (kb-creator, roadmap-generator, etc.) según cuáles estén instalados/
  activos (lazy-loading, §4.2). No es una skill monolítica.
- **Toggles del `sdd-orchestrator`**: `lean-workflow` selecciona los contratos compactos por defecto. Al desactivarlo, el texto legacy sale de los assets del stack viejo.
  Base (siempre, con `delegation` y `model-routing` como secciones marcadas que
  se quitan si su toggle falta) + los fragmentos aditivos, que se anexan en este orden:
  `lean-workflow` (selector de modo), `delegation`, `model-routing`, `governance`, `evidence`, `auditor`, `roles`, `review`, `knowledge-base`,
  `senior`, `ship`, `engram` y `tdd`. El orden está argumentado en los
  comentarios de `internal/harness/config/compose.go`, y `config.KnownToggles()`
  es la lista canónica.

### Follow-ups / abiertas
- ~~**Fix engram download**~~ → **RESUELTO**: `model.External` ya separa `Repo`
  (`Gentleman-Programming/engram`) del `Pkg` (ver `internal/model/harness.go` y el catálogo).
- ~~**Nombre del repo remoto en GitHub**~~ → **RESUELTO** (2026-05-30): el viejo
  `JuanCruzRobledo/jr-stack` (repo fuente original) queda **archivado como legacy**; este conserva `active-stack`.
- ~~**Mapeo harness↔modo**~~ → **RESUELTO** (CHANGES.md C-20): Lite = sustrato,
  Full = sustrato + fundación guiada, Custom = todos. `active-orchestrator` movido a Full-only.
- **Empaquetado de skills de terceros** → se clonan como las propias, con `path` al subdirectorio del repo (ej. `find-skill` instala `vercel-labs/skills`, subdir `skills/find-skills`). No se usa `npx`.
- **Custom + `permissions`** → decisión: NO desactivable (security-first); falta implementar (CHANGES.md **C-21**).
