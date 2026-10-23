# Revisiones de EXPLORE/PROPOSE, fase 2

Fecha: 2026-10-23. Base de código observada: `0204022`. Árbol documental en elaboración, sin SHA propio al revisar. No son informes de verify ni evidencia de producto construido.

## Procedencia

**Primera revisión, archivos 01 a 05: custodia no verificable, no cuenta como revisión.** Se conservan como historia. La segunda auditoría no pudo contrastar los textos archivados con despachos y respuestas originales; las afirmaciones anteriores de originales completos e independencia no se acreditan. El veredicto histórico PRESENTABLE no cubre el estado actual del propose.

| Archivo histórico | Estado |
|---|---|
| 01-arquitecto.md | custodia no verificable, no cuenta como revisión |
| 02-funcional.md | custodia no verificable, no cuenta como revisión |
| 03-qa.md | custodia no verificable, no cuenta como revisión |
| 04-devops-seguridad.md | custodia no verificable, no cuenta como revisión |
| 05-auditor.md | custodia no verificable, no cuenta como revisión |

La [segunda revisión](segunda-revision/05-auditor.md) registra AN-1, AR-1, AN-2, AR-2, OPS-1, QA-1 y AUD-1. Las [decisiones de Juani](segunda-revision/06-decisiones-juani.md) disponen corrección del constructor y confirmación individual del emisor sobre su delta, sin otra ronda completa. Correcciones presentadas en el propose; confirmaciones pendientes. No se acredita aislamiento técnico ni diversidad de modelos. Los logs ocultos preexistentes se conservan sin agregarlos al commit.

## Respuesta histórica del autor a los hallazgos (sin acreditar revisión)

- ARQ-1: antes del merge, JSON inválido aborta sin escritura; prueba de bytes exactos. ARQ-2: dispatch por ID y assets antes de referencias. ARQ-3: recuperación interna del paso fallido. ARQ-4: delta permitido P→A definido, tabla humana sin autorreferencia de hashes.
- Funcional 1: aprobación escrita por alcance; crypto sólo tras elección J1. Funcional 2: degradación prospectiva no convierte faltante en válido ni permite cierre ordinario.
- QA B-QA-1: ausencia de firma sólo produce WRITTEN_UNVERIFIED si decisión completa; MISSING/INVALID/STALE discriminados y probados. Gates de S, cortes de persistencia y concurrencia explícitos.
- DevOps: eventos como datos, quoting por plataforma, hooks personales conservan orden/semántica; roles validan ownership actual/concurrencia/symlinks; confianza Git maliciosa incluida en pruebas futuras.

Backend: cubierto por arquitecto en fronteras de paquetes y contratos, sin implementación; especialista backend adicional no convocado. DBA/frontend: no hay persistencia SQL ni interfaz; no convocados. Senior: no corresponde a propose. El rol DevOps incluye revisión de seguridad de confianza y filesystem; no se afirma pentest ni revisión de cliente ejecutado.

## Validaciones documentales

`openspec validate metodologia-fase-2` → válido, con aviso informativo de skip_specs para entrega documental sin cambios de comportamiento. `git diff --check` → exit 0 sobre rastreados; se verificará además índice para incluir archivos nuevos. No se ejecutó Go ni build. OpenSpec local ignorado por Git; sólo proposal documental referencial, sin apply, archive ni autorización humana inventada.
