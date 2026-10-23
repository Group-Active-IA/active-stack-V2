# Revisiones de EXPLORE/PROPOSE, fase 2

Fecha: 2026-10-23. Base de código observada: `0204022`. Árbol documental en elaboración, sin SHA propio al revisar. No son informes de verify ni evidencia de producto construido.

## Procedencia

01–04 preservan íntegro el retorno final de cada subagente convocado secuencialmente, copiado por el orquestador a archivo separado. Revisores recibieron sólo lectura y declararon no leer otros informes; el auditor recibió los cuatro originales. Mediación: SÍ, copia textual del retorno por el orquestador; no captura automática ni prueba criptográfica. Logs ocultos preexistentes, incluido `.00.log`, no son fuente de estos informes ni se agregan al commit.

Familia/modelo/proveedor de los revisores: configuración heredada del agente principal, sin override; nombre exacto y proveedor verificables en registro de despacho no disponibles aquí: SUPOSICIÓN, no atribución de diversidad. Filesystem compartido, sin aislamiento técnico demostrado. Especialistas revisaron revisiones sucesivas del borrador, no un SHA documental congelado común; auditor lee el resultado final y los originales, sin presentar coincidencia como convergencia independiente sobre idéntico árbol.

| Archivo | Rol / despacho | Estado |
|---|---|---|
| 01-arquitecto.md | arquitecto_fase2 | Original completo; ARQ-1 y notas incorporados |
| 02-funcional.md | funcional_fase2 | Original completo; ambigüedades de crypto y faltantes corregidas |
| 03-qa.md | qa_fase2 | Original completo; estados y gates incorporados |
| 04-devops-seguridad.md | devops_fase2 | Original completo; contratos de hooks, paths y confianza agregados |
| 05-auditor.md | auditor_propose_fase2 | Original completo; PRESENTABLE para decisión humana, sin autorización de apply |

## Respuesta del autor a los hallazgos

- ARQ-1: antes del merge, JSON inválido aborta sin escritura; prueba de bytes exactos. ARQ-2: dispatch por ID y assets antes de referencias. ARQ-3: recuperación interna del paso fallido. ARQ-4: delta permitido P→A definido, tabla humana sin autorreferencia de hashes.
- Funcional 1: aprobación escrita por alcance; crypto sólo tras elección J1. Funcional 2: degradación prospectiva no convierte faltante en válido ni permite cierre ordinario.
- QA B-QA-1: ausencia de firma sólo produce WRITTEN_UNVERIFIED si decisión completa; MISSING/INVALID/STALE discriminados y probados. Gates de S, cortes de persistencia y concurrencia explícitos.
- DevOps: eventos como datos, quoting por plataforma, hooks personales conservan orden/semántica; roles validan ownership actual/concurrencia/symlinks; confianza Git maliciosa incluida en pruebas futuras.

Backend: cubierto por arquitecto en fronteras de paquetes y contratos, sin implementación; especialista backend adicional no convocado. DBA/frontend: no hay persistencia SQL ni interfaz; no convocados. Senior: no corresponde a propose. El rol DevOps incluye revisión de seguridad de confianza y filesystem; no se afirma pentest ni revisión de cliente ejecutado.

## Validaciones documentales

`openspec validate metodologia-fase-2` → válido, con aviso informativo de skip_specs para entrega documental sin cambios de comportamiento. `git diff --check` → exit 0 sobre rastreados; se verificará además índice para incluir archivos nuevos. No se ejecutó Go ni build. OpenSpec local ignorado por Git; sólo proposal documental referencial, sin apply, archive ni autorización humana inventada.
