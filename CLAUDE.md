# CLAUDE.md — Active Stack (instalador del harness metodológico)

> Constitución operativa del proyecto. Versión canónica (Claude).
> Espejo de `AGENTS.md`: **si modificás uno, actualizás el otro.**
> Leé este archivo antes de cualquier acción no trivial.

> **Nota sobre el orquestador SDD.** El orquestador SDD global (coordinación,
> delegación, ruteo de modelos, memoria y el resto de los toggles del
> `sdd-orchestrator`) ya vive en tu `~/.claude/CLAUDE.md`. **Este archivo NO lo duplica.** Si necesitás
> el detalle del orquestador, referenciá ese archivo global. Acá viven solo las
> reglas **específicas de este proyecto**.

---

## 1. Stack y topología

- **Lenguaje**: Go 1.26.
- **TUI**: Bubbletea + Lipgloss (sin el theme cosmético del repo viejo).
- **Distribución**: binario único, cross-platform — Windows, macOS, Linux, WSL y Termux.
- **Entrypoint**: `cmd/active-stack/`.
- **Catálogo**: embebido en el binario vía `//go:embed` (`internal/catalog/harnesses.yaml`).

Qué es esto: un **instalador methodology-first**. Materializa el `MANUAL-METODOLOGICO.md`.
Un comando (`active-stack install`) instala/configura el sustrato (harnesses); un único
orquestador de fundación deja el proyecto listo para el ciclo OPSX
(`explore → propose → apply → verify → archive`).

**NO es** (fuera de scope, decisión firme — ver ARCHITECTURE.md §1):
GGA (code review en commit), themes/statusline/keybindings, persona como producto
de marketing, framing "supercharge any agent" / relación con Gentleman.Dots.

### Estructura de paquetes (ARCHITECTURE.md §5)

```
cmd/active-stack/            entrypoint CLI
internal/
  system/                detección OS/arch/WSL/Termux, deps, guards
  catalog/               parseo del harnesses.yaml embebido
  model/                 tipos de dominio (harness, agente, modo)
  planner/               grafo de dependencias, orden, review payload
  agents/                adapters por agente (claude/opencode/...)
  harness/               install/inject por tipo de harness
    skill/  config/  external/
  filemerge/             merge por markers (inyectar sin pisar)
  backup/                snapshot + restore de configs
  pipeline/              ejecución por etapas + rollback
  verify/                health checks post-install
  tui/                   Bubbletea (sin theme cosmético del viejo)
assets/                  catálogo + configs bundleadas (sdd-orchestrator, etc.)
```

Procedencia histórica: port del repo viejo (`E:\ESCRITORIO\programar\2026\jr-stack`), limpiado de branding excluido.
Los paquetes listados ya existen. Se descarta: `components/gga`, `components/theme`,
`components/persona` (marketing). `components/permissions` **se mantiene** (la
seguridad no es opcional).

---

## 2. Modelo de dominio (resumen)

Detalle real en `internal/model/harness.go`. Tipos clave:

- **`Harness`** — módulo instalable/configurable. Campos: `ID`, `Name`, `Type`,
  `Source` (skill), `External` (tool), `Toggles` (config), `InstallModes`,
  `DependsOn`, `Agents`. Métodos: `InMode(mode)`, `SupportsAgent(agent)`.
- **`HarnessType`** — cómo se materializa un harness:
  - `skill` → `SKILL.md` + assets, clonado de un repo y copiado al dir de skills del agente.
  - `config` → texto/archivos bundleados que configuran el agente (ej. `sdd-orchestrator`, `permissions`).
  - `external` → binario/servicio de terceros que instalamos/configuramos pero no son nuestros (Engram, OpenSpec, Context7).
- **`InstallMode`** — `lite` | `full` | `custom`. Convención: un harness Lite lista
  `[lite, full]` (Full incluye a Lite); uno exclusivo de Full lista `[full]`; Custom matchea todos.
- **`Agent`** — `claude`, `opencode`, `gemini`, `codex`, `cursor`, `vscode`, `windsurf`, `antigravity`.
- **Catálogo** (`internal/catalog`) — `Load()` parsea y valida el YAML embebido;
  un catálogo malformado es error de build/release (falla ruidoso). Métodos:
  `ByID`, `ForMode`, `ForAgent`.

**El harness `sdd-orchestrator` es clave**: es de tipo `config` y se compone a
partir de **toggles modulares**. El resultado es el bloque de instrucciones del
orquestador que se inyecta en `CLAUDE.md`/`AGENTS.md` del proyecto destino.

Los toggles, en orden de composición, y qué agrega cada uno (la lista canónica es
`config.KnownToggles()`; el detalle largo, `ARCHITECTURE.md` §2.1):

- `delegation` — qué se resuelve inline y qué se delega, para no inflar contexto.
- `model-routing` — modelo por fase, con sustitución declarada si no hay acceso.
- `governance` — autonomía según criticidad del dominio (CRITICAL → LOW).
- `evidence` — reporte tipado: EVIDENCIA (comando, salida, hora), JUICIO, SUPOSICIÓN.
- `auditor` — lee los reportes de los otros y marca lo que no se sostiene.
- `roles` — convoca un rol con estándar fijo y "listo" observable.
- `review`: un rompedor con tres lentes obligatorias y bucle de hasta tres vueltas; el auditor decide y el orquestador ejecuta.
- `knowledge-base`: en archive reescribe la base que impacto y diff tocaron; agrega decisiones, lecciones L-NNN y deuda del auditor con evidencia por entrada. Sólo APPEND sin REWRITE se reporta SUSPICIOUS: append-only.
- `senior`: sólo en verify, después del rompedor y antes del auditor; aporta contexto de decisiones fechadas, uno por vuelta y hasta tres por change.
- `ship` — diagnostica la publicación y espera el CI completo del SHA pusheado; `main` pregunta, PR o directo se propone, nunca mergea.
- `engram` — memoria persistente entre sesiones y compactaciones.
- `tdd` — TDD estricto en apply: sin test que falle primero no hay código.

---

## 3. Reglas críticas NO negociables

Formuladas como "NUNCA X → hacer Y". Estas reglas son duras: violarlas invalida el trabajo.

- **NUNCA pisar config del usuario sin backup** → SIEMPRE snapshot vía `internal/backup` antes de escribir.
- **SIEMPRE inyectar con markers idempotentes** → usar `internal/filemerge` (merge por markers); reinstalar no debe duplicar bloques.
- **NUNCA hardcodear paths de agente** → resolver SIEMPRE vía el adapter del agente (`internal/agents`).
- **NUNCA commitear sin pedido explícito** del responsable; *conventional commits* exclusivamente; **NUNCA** atribución de coautoría a la IA.
- **NUNCA meter en el repo lo que se saca** → GGA, theme, statusline, keybindings, persona-marketing y el framing "supercharge any agent" NO van. La persona puede sobrevivir solo como harness `config` opcional, sin envoltorio de marketing.
- **SIEMPRE limpiar leftovers `gentle-ai` / `gentle-stack` / `Gentleman.Dots`** al portar código del repo viejo (paths, strings, nombres, branding).
- **NUNCA instalar "repos"** → el instalador instala **harnesses**, y cada harness sabe cómo se materializa (skill/config/external).
- **NUNCA editar `harnesses.yaml` sin pasar por `catalog.Load()` validando** → un catálogo inválido rompe el release.
- **NUNCA build después de cambios** salvo pedido explícito (regla del operador).
- **SIEMPRE marcar `(TBD)`** cuando una decisión no está tomada; nunca inventar.

---

## 4. Governance por dominio

El nivel de autonomía es proporcional a la criticidad del paquete tocado.

| Nivel | Dominios (paquetes) | Comportamiento del agente |
|---|---|---|
| **ALTO** | `backup` (snapshot/restore), `filemerge` (merge por markers), `pipeline` (rollback) | Propone y espera review. Pueden **destruir config del usuario**. |
| **MEDIO** | `agents` (adapters por agente), `harness/*` (installers) | Implementa con checkpoints. |
| **BAJO** | `catalog`, `model`, `tui` | Autonomía completa si los tests pasan. |

> No hay dominios CRÍTICO en este proyecto (no manejamos auth/billing/secrets de
> producción). El máximo es ALTO porque destruir la config del usuario es el peor
> daño posible del instalador.

---

## 5. Mapa de navegación ("Necesito X → Leer Y")

| Necesito… | Leer / mirar |
|---|---|
| El blueprint del diseño (fuente de verdad) | `ARCHITECTURE.md` |
| La metodología (etapas, governance, flujo) | `../MANUAL-METODOLOGICO.md` |
| Roadmap de changes, deps, camino crítico | `CHANGES.md` |
| Qué harnesses existen y de qué tipo | `internal/catalog/harnesses.yaml` |
| Los tipos de dominio | `internal/model/harness.go` |
| Cargar/validar el catálogo | `internal/catalog/catalog.go` |
| Infra implementada (backup, filemerge, planner, agents, pipeline, verify, tui) | `internal/`; procedencia histórica: repo viejo `jr-stack/internal/` |
| Configs bundleadas (sdd-orchestrator por agente) | `internal/harness/config/assets/` |
| Estado real de los changes | `openspec` CLI (`openspec list`, `openspec status`) — fuente de verdad |
| El orquestador SDD global | `~/.claude/CLAUDE.md` (NO duplicar acá) |

---

## 6. Ciclo vigente y evidencia

El ciclo vigente es `explore → propose → apply → verify → archive`, seguido de `ship` cuando se publica:

- **Propose:** con governance activo, ALTO o CRÍTICO requiere informes de arquitecto, analista funcional y QA, más especialistas por capa, convocados de a uno sin leer los otros informes; el auditor lee los N originales antes de la aprobación humana escrita. La tabla `## Aprobación` registra elección, persona, fecha y alcance; ningún agente la completa. Una regla nueva necesita aprobación propia.
- **Apply constructor:** rojo y verde por tarea en `tasks.md`, con comando, salida, hora y SHA; si el test nació verde, mutante que falle en su aserción y se revierta. Declara consumidores, conflictos y familia del modelo. En vueltas 2 y 3 corrige sólo la lista del auditor.
- **Verify:** conformidad → rompedor con tres lentes (correctitud, cableado, documentación) → senior si está activo → auditor. La conformidad sola no habilita archive. El auditor decide y el orquestador ejecuta; sólo APROBADO con informes completos y conformidad satisfecha habilita el cierre ordinario. El contador y los informes viven en el change: máximo tres vueltas, sin cuarta automática. En la tercera con bloqueantes se exige decisión humana escrita sobre alcance nuevo, deuda justificada o freno; no se fabrica APROBADO.
- **Archive:** lee la unión del impacto declarado y el diff; reescribe las descripciones actuales y agrega decisiones, lecciones numeradas y deuda del auditor con fecha, costo y evidencia tipada. Deja `knowledge-base-report.md` por change con REWRITTEN, APPENDED y PENDING por documento; sólo APPEND se justifica y se señala como `SUSPICIOUS: append-only` cuando correspondía reescribir.
- **Ship:** espera todo el CI requerido sobre el SHA pusheado, inspeccionando workflows, eventos y jobs. En active-stack exige `ci.yml` en push con tests Ubuntu y Windows; en otros repos deriva su conjunto. Cero runs, consulta fallida, SHA ajeno o estados incompletos no confirman entrega. `main` exige permiso y nunca se mergea el PR automáticamente.

Si falta un toggle se declara la degradación; sin review no hay rompedor ni bucle. Sin auditor se necesita decisión escrita de Juani para cerrar con alcance degradado. Un informe obligatorio faltante detiene el cierre.

`active-stack version` informa versión y SHA disponibles. El chequeo instalado valida `cmd/active-stack`, `roles --help` y contenido contra `Compose`; no acredita reinstalación real. CONF-1 mantiene `instalacion-sana` abierto hasta §3.6 puntos 1–4 y OK escrito para reinstalar.

EVIDENCIA: `84fcf0e`, `c325955` y [base metodológica](docs/metodologia-opsx.md) (commits por fase, decisiones, lecciones y deuda; informe 22 autoriza sólo ocho archives).

## 7. Notas

- `openspec/` está versionado-ignorado por git (dogfooding interno).
- Cada incremento del roadmap es un change OPSX (ver `CHANGES.md`).
- Skills de dominio relevantes para este repo: `go-testing` (tests Go + Bubbletea TUI).
