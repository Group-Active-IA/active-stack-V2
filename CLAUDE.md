# AGENTS.md / CLAUDE.md — Active Stack (instalador del harness metodológico)

> Constitución operativa del proyecto. Versión **model-agnostic**.
> `AGENTS.md` y `CLAUDE.md` son idénticos byte a byte (un test lo verifica): **si modificás uno, copiá el otro.**
> Todo agente lee este archivo antes de cualquier acción no trivial.

---

## 1. Stack y topología

- **Lenguaje**: Go 1.26.
- **TUI**: Bubbletea + Lipgloss (sin el theme cosmético del repo viejo).
- **Distribución**: binario único para Windows, Linux y WSL (amd64).
- **Entrypoint**: `cmd/active-stack/`.
- **Catálogo**: embebido en el binario vía `//go:embed` (`internal/catalog/harnesses.yaml`).

Qué es esto: un **instalador methodology-first**. Materializa la metodología OPSX descrita en [docs/metodologia-opsx.md](docs/metodologia-opsx.md) y [docs/lean-agent-workflow.md](docs/lean-agent-workflow.md).
Un comando (`active-stack install`) instala/configura el sustrato (harnesses); un único
orquestador de fundación deja el proyecto listo para el ciclo OPSX
(`explore → propose → apply → verify → archive`).

**NO es** (fuera de scope, decisión firme — ver ARCHITECTURE.md §1):
GGA (code review en commit), themes/statusline/keybindings, persona como producto
de marketing, framing "supercharge any agent" / relación con Gentleman.Dots.

### Estructura de paquetes (ARCHITECTURE.md §5)

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

Procedencia histórica: port del repo anterior (`jr-stack`), limpiado de branding excluido.
Los paquetes listados ya existen. Del repo viejo se descartó GGA, theme y persona
(marketing); permissions **se mantiene** como harness `config` (la seguridad no es opcional).

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
  - `command` → slash-commands nativos del agente, embebidos en el binario (`/active:init`, `/opsx:ship`, `/opsx:verify`).
- **`InstallMode`** — `lite` | `full` | `custom`. Convención: un harness Lite lista
  `[lite, full]` (Full incluye a Lite); uno exclusivo de Full lista `[full]`; Custom matchea todos.
- **`Agent`** — `claude`, `opencode`, `gemini`, `codex`, `cursor`, `vscode`, `windsurf`, `antigravity`.
- **Catálogo** (`internal/catalog`) — `Load()` parsea y valida el YAML embebido;
  un catálogo malformado es error de build/release (falla ruidoso). Métodos:
  `ByID`, `ForMode`, `ForAgent`.

**El harness `sdd-orchestrator` es clave**: es de tipo `config` y se compone a
partir de **toggles modulares**. El resultado es el bloque de instrucciones del
orquestador que se inyecta en `CLAUDE.md`/`AGENTS.md` del proyecto destino.

`lean-workflow` es la opción predeterminada: cinco roles secuenciales, contratos compactos, contexto acotado y TDD obligatorio para código. Los roles van dentro de las fases OPSX, la profundidad escala por nivel y cada proyecto declara quién ejecuta cada rol en `## Role executors`; este repositorio no integra proveedores. Los especialistas son opcionales y senior no agrega una etapa. Ver [flujo compacto de cinco roles](docs/lean-agent-workflow.md).

Toggles del `sdd-orchestrator` (lista canónica: `config.KnownToggles()`; el orden no cambia el resultado). Con `lean-workflow` activo sólo `governance`, `evidence`, `knowledge-base`, `engram` y `ship` agregan su contrato compacto; los marcados **sólo legacy** se aceptan por compatibilidad pero no pueden apagar nada del flujo lean (TDD y auditor son siempre obligatorios). Se eligen con `active-stack install --toggles`:

- `lean-workflow` (selector) — selector predeterminado; quitarlo activa la composición legacy.
- `delegation` (sólo legacy) — qué se resuelve inline y qué se delega, para no inflar contexto.
- `model-routing` (sólo legacy) — modelo por fase, con sustitución declarada si no hay acceso.
- `governance` (lean y legacy) — autonomía según criticidad del dominio (CRITICAL → LOW).
- `evidence` (lean y legacy) — reporte tipado: EVIDENCIA (comando, salida, hora), JUICIO, SUPOSICIÓN.
- `auditor` (sólo legacy) — lee los reportes de los otros y marca lo que no se sostiene.
- `roles` (sólo legacy) — convoca un rol con estándar fijo y "listo" observable.
- `review` (sólo legacy): un rompedor con tres lentes obligatorias y bucle de hasta tres vueltas; el auditor decide y el orquestador ejecuta.
- `knowledge-base` (lean y legacy): en archive reescribe la base que impacto y diff tocaron; agrega decisiones, lecciones L-NNN y deuda del auditor con evidencia por entrada. Sólo APPEND sin REWRITE se reporta SUSPICIOUS: append-only.
- `senior` (sólo legacy): sólo en verify, después del rompedor y antes del auditor; aporta contexto de decisiones fechadas, uno por vuelta y hasta tres por change.
- `ship` (lean y legacy) — diagnostica la publicación y espera el CI completo del SHA pusheado; `main` pregunta, PR o directo se propone, nunca mergea.
- `engram` (lean y legacy) — memoria persistente entre sesiones y compactaciones.
- `tdd` (sólo legacy) — en lean TDD es siempre obligatorio y este toggle no cambia nada; en legacy: TDD estricto en apply: sin test que falle primero no hay código.

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

> Esta tabla es el riesgo de los paquetes de este repo. Los changes usan la escala del
> flujo (LOW/MEDIUM/HIGH/CRITICAL): ALTO equivale a HIGH, MEDIO a MEDIUM y BAJO a LOW.
> No hay dominios CRÍTICO en este proyecto (no manejamos auth/billing/secrets de
> producción). El máximo es ALTO porque destruir la config del usuario es el peor
> daño posible del instalador.

---

## 5. Mapa de navegación ("Necesito X → Leer Y")

| Necesito… | Leer / mirar |
|---|---|
| El blueprint del diseño (fuente de verdad) | `ARCHITECTURE.md` |
| La metodología (etapas, governance, flujo) | `docs/lean-agent-workflow.md` (predeterminada) y `docs/metodologia-opsx.md` (historia y legacy) |
| Roadmap de changes, deps, camino crítico | `CHANGES.md` |
| Qué harnesses existen y de qué tipo | `internal/catalog/harnesses.yaml` |
| Los tipos de dominio | `internal/model/harness.go` |
| Cargar/validar el catálogo | `internal/catalog/catalog.go` |
| Infra implementada (backup, filemerge, planner, agents, pipeline, verify, tui) | `internal/`; procedencia histórica: repo viejo `jr-stack/internal/` |
| Configs bundleadas (sdd-orchestrator por agente) | `internal/harness/config/assets/` |
| Estado real de los changes | `openspec` CLI (`openspec list`, `openspec status`) — fuente de verdad |

---

## 6. Ciclo vigente y evidencia

El ciclo OpenSpec conserva `explore → propose → apply → verify → archive` y publicación separada. En el modo predeterminado lean, el agente principal coordina:

1. Planificador: propuesta, diseño, tareas y escenarios.
2. Revisor del plan: inconsistencias y viabilidad; devuelve correcciones al planificador.
3. Constructor: implementación con baseline, RED de aserción observado, GREEN y comprobaciones de regresión/refactor.
4. Rompedor: casos adversos de lógica, consumidores y documentación; defectos reproducidos vuelven al constructor.
5. Auditor: originales completos, conformidad, evidencia y cobertura; deriva defectos y exige repetir verificaciones afectadas.

Cada change registra LOW/MEDIUM/HIGH/CRITICAL por riesgo y alcance; en LOW el orquestador puede fusionar plan y revisión, y desde MEDIUM trabajan cinco agentes separados. Con governance activo, HIGH/CRITICAL requieren aprobación humana explícita antes de construcción. TDD de código no depende del toggle histórico `tdd`; documentación sin comportamiento ejecutable declara validación apropiada. Nunca se fabrica RED, informe o aprobación.

El orquestador prepara encargos con objetivo, archivos/diff, decisiones, consumidores, límites, criterios de salida y ruta de informe. Persiste rol, revisión, intento, informes, bloqueantes y próximo paso; OpenSpec es fuente de estado de artefactos. Informes faltantes, evidencia desactualizada o bloqueantes impiden cierre. Tres intentos fallidos requieren dirección humana. No se agrega senior ni fan-out obligatorio de especialistas.

El conocimiento se actualiza en el destino declarado por el proyecto (incluido obsidian.md), sin duplicar notas ni copiar toda la memoria a cada encargo. La publicación conserva sus permisos y CI sobre el SHA entregado; nunca merge automático. No reinstalar dotfiles ni ejecutar build por este trabajo.

Al desactivar explícitamente lean-workflow se usan los fragmentos legacy y su procedimiento anterior. La guía [lean-agent-workflow](docs/lean-agent-workflow.md) define precedencia y límites; el historial metodológico sigue en [docs/metodologia-opsx.md](docs/metodologia-opsx.md).

`active-stack version` informa versión y SHA disponibles. El chequeo instalado valida `cmd/active-stack`, `roles --help` y contenido contra `Compose`; no acredita reinstalación real. La instalación real se hizo el 29/10/2026 con OK escrito del responsable, después de una prueba en HOME aislado (instalar, reinstalar, restaurar); falta el cierre formal de `instalacion-sana` contra la receta §3.6 (CONF-1).

EVIDENCIA: `84fcf0e`, `c325955` y [base metodológica](docs/metodologia-opsx.md) (commits por fase, decisiones, lecciones y deuda; informe 22 autoriza sólo ocho archives).

## 7. Notas

- Este archivo gobierna el desarrollo del instalador. Los contratos que instala se generan desde los assets versionados; no se presupone un orquestador global ya instalado ni se reinstala al modificar este repositorio.
- `openspec/` está versionado-ignorado por git (dogfooding interno).
- Cada incremento del roadmap es un change OPSX (ver `CHANGES.md`).
