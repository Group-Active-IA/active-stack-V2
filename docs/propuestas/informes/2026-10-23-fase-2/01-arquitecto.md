Informe ARQUITECTO — PROPOSE ALTO. Revisión de diseño; no constituye aprobación humana.

**JUICIO general:** la separación entre aprobación, captura, adaptadores y entrega de roles es viable. Recomiendo corregir un bloqueo de preservación antes de presentar el alcance para aprobación.

**ARQ-1 — Bloqueante: configuración OpenCode inválida puede perder contenido ajeno.**

EVIDENCIA: `permissions/install.go`, `mergeJSONFile`, entrega directamente los bytes existentes a `filemerge.MergeJSONObjects`. En `internal/filemerge/json_merge.go:9`, un error al decodificar la base se convierte deliberadamente en objeto vacío; el merge continúa y produce sólo el overlay. La propuesta §10 exige preservar hermanos personales y settings ajenos, pero sus pruebas previstas no incluyen base inválida.

JUICIO: un snapshot permite recuperar, pero no satisface la preservación prometida ni justifica reportar éxito tras descartar contenido. Agregar al contrato: configuración preexistente no interpretable → error visible y cero escritura; reparación sólo mediante alcance aprobado aparte. Puede implementarse una validación previa en el installer sin modificar `filemerge`. Exigir prueba discriminante con bytes inválidos y contenido personal, verificando conservación exacta.

**Notas de arquitectura para los changes posteriores:**

- **ARQ-2 — Dispatch de harness config.** EVIDENCIA: `internal/harness/config/install.go` compone e inyecta el orquestador para cualquier `HarnessConfig`; no distingue `role-definitions` por ID. JUICIO: explicitar dispatch por harness y dependencia/orden para que no se publique referencia a roles antes de instalarlos. La tabla general de integración todavía no define ese punto.
- **ARQ-3 — Recuperación del paso parcialmente fallido.** EVIDENCIA: `internal/pipeline/rollback.go` documenta rollback de pasos apply exitosos. JUICIO: roles e instalación de hooks deben recuperar también sus propias escrituras parciales antes de retornar error, o registrar recuperación incompleta. Reutilizar pipeline no demuestra esa garantía.
- **ARQ-4 — Aprobación sin autorreferencia documental.** JUICIO: concretar dónde vive la tabla humana de A y qué diferencia documental autorizada introduce respecto de P. De lo contrario, completar la propia tabla puede parecer una alteración posterior del documento aprobado. La separación A/P propuesta permite resolverlo sin rediseño.

**Fuentes leídas:** `AGENTS.md`; exploración y propuesta fase 2 completas; `internal/harness/config/{compose,install}.go`; `permissions/install.go`; `internal/filemerge/json_merge.go`; `internal/agents/interface.go`; referencias localizadas mediante `rg` en overlays, tests de merge/tier, adapters, pipeline y dispatch de roles. No leí informes de otros revisores.

**Límites:** EVIDENCIA exclusivamente por lectura; no ejecuté tests, build, Compose ni clientes. No escribí archivos ni accedí a configuración global. Las afirmaciones externas de la exploración no fueron verificadas independientemente. Filesystem compartido: independencia procedimental por fuentes elegidas; aislamiento técnico no demostrado. SUPOSICIÓN: futuros changes resolverán los `(TBD)` antes de habilitar funcionalidades afectadas.
