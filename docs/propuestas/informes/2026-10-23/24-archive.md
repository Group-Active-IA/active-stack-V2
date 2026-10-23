# Archive de fase 1 — 2026-10-23

**JUICIO del archivista:** se cubre la unión del alcance firmado, los ocho proposals y el diff `4748680..5ab0f4f`. Ningún proposal tiene `## Knowledge Base Impact`; se aplicó el fallback al diff, sin interpretar ausencia como “ningún impacto”. Se leyeron el fragmento `internal/harness/config/assets/knowledge-base.md` (archive, reporte y arrastre), la propuesta firmada, fuentes de implementación, informes de las tres vueltas y OK escrito. No se lanzó otro agente ni se introdujo regla metodológica nueva.

## Autorización y límite

**EVIDENCIA:** `5ab0f4f`, informe 22: OK escrito de Juani tras corrección `ce7899e` y verificación independiente; permite archivar ocho changes, excluye `instalacion-sana`. Informe 21 (`18bb293`): vuelta 3 de 3 NO APROBADA, RF3-1 abierto entonces. **JUICIO:** el cierre usa la decisión humana posterior, no un APROBADO retroactivo ni una cuarta vuelta. El informe 23 registra la evidencia de corrección (suite 26 ok, 0 fallas, 10:56, según decisión 22); no se atribuye esa ejecución al archivista.

La aprobación se conserva como fuente humana; no se modifica su tabla ni se inventa firma criptográfica. No build, push, reinstalación, escritura en ~/.claude, ~/.local/bin ni .env. Los logs ocultos preexistentes quedan sin agregar al commit. `openspec/` está ignorado por git: los archives persisten localmente y el registro versionado de esta operación es este informe.

## Alcance por change y evidencia

| Change | Implementación / motivo documental |
|---|---|
| opencode-mcp-command-array | `8c6a06b`: MCP local command array y environment; base metodológica registra el contrato actual |
| ship-espera-ci | `60b5e83`, `140f5d5`, `61b23ea`, `a04a6ba`: CI completo/SHA, alcance por repo y autoridad; quedan paráfrasis RF3-2 |
| apply-constructor | `f3fc02a`: rojo/verde por tarea, consumidores, mutantes y límites del constructor |
| verify-rompedor-auditor | `db3de54`, `0b5d0d4`, `ce7899e`: bucle tres vueltas, degradación y corrección RF3-1 |
| senior-solo-en-verify | `66ad39e`: verify tras rompedor y antes de auditor, uno por vuelta; conserva CONF-3/RF2-3 |
| archive-base-completa | `84c570d`, `13ce7f2`: unión impacto/diff, REWRITE/APPEND, evidencia y ausencia de auditor condicional |
| governance-aprobacion | `9f2d03b`: aprobación humana escrita y alcance, regla nueva con aprobación propia |
| propose-multirol | `531e1b6`: arquitecto, analista, QA y especialistas, auditor sobre originales antes de firma |

**EVIDENCIA:** `git diff --stat 4748680..HEAD` inicial incluyó los cuatro documentos raíz ya tocados, fragmentos, comandos, catálogo, tests y docs de informes. Se cubrieron README, ARCHITECTURE, CHANGES, ambas constituciones y la propuesta; los informes fechados y la exploración son historia, no contratos actuales que deban sustituirse. `docs/archive/` y `docs/superpowers/` no describen este nuevo ciclo y no recibieron cambios. No se editó catálogo ni código Go; la validación previa del catálogo consta en los informes del constructor, no se simula una validación nueva.

## REWRITTEN — por documento

| Documento | Qué se sustituyó / qué describe ahora | Evidencia y razonamiento tipados |
|---|---|---|
| README.md | Ficha ship ya exige CI completo; estado deja de afirmar CI remoto verde sin consulta; ciclo operativo y version/checks con límite CONF-1 | EVIDENCIA `60b5e83`, `84fcf0e`, informe 22. JUICIO: el usuario necesita distinguir instrucciones implementadas de entrega real |
| ARCHITECTURE.md | Se elimina “planificación”, “esqueleto en curso” y paquetes por construir/portar; ciclo genérico pasa a secuencia, autoridad, tres vueltas, archive y CI; roadmap actual con CONF-1 | EVIDENCIA paquetes existentes en internal/, commits de tabla anterior y CHANGES C-01…C-25. JUICIO: el blueprint viejo orientaba como futuro lo ya implementado |
| AGENTS.md | Topología ya implementada y rutas actuales de assets; ficha ship y contrato del ciclo/instalación, con evidencia y enlace a historia | EVIDENCIA `84fcf0e`, `c325955`, `db3de54`, `531e1b6`. JUICIO: evitar reglas locales incompletas sobre cierre y entrega |
| CLAUDE.md | Mismos cambios operativos que AGENTS; preservadas diferencias introductorias preexistentes | EVIDENCIA diff de ambos documentos. JUICIO: mantener espejo sin editar configuración global |
| CHANGES.md | “PROPUESTO, sin firma” sustituido por estado autorizado, paso 7b omitido incorporado y siguiente trabajo metodológico fijado en CONF-1/fases futuras | EVIDENCIA `4748680`, `5ab0f4f`, informe 22. JUICIO: las fichas históricas de otros epics no se reinterpretan como trabajo de esta fase |
| docs/propuestas/2026-10-23-metodologia-proptech-en-opsx.md | Encabezado de propose y “OpenSpec no inicializado” sustituidos por estado firmado/implementado con archive parcial; diagnóstico original del cuerpo identificado como histórico | EVIDENCIA openspec list inicial, `4748680`, `5ab0f4f`. JUICIO: preservar alcance y aprobación originales sin presentar el diagnóstico antiguo como presente |
| docs/metodologia-opsx.md | Documento nuevo, no REWRITE de archivo previo: crea referencia actual completa para ciclo, degradación, instalación/version y límites de cobertura | EVIDENCIA fuentes y SHAs por sección. JUICIO: la base estaba distribuida; la referencia evita repetir deuda e historia en cinco documentos |
| openspec/changes/instalacion-sana/tasks.md | Estado obsoleto “vuelta 2 en curso” sustituido por vuelta 3 cerrada y exclusión humana; se corrige la lectura de “excluida” como exención | EVIDENCIA informes 17, 19, 21 y 22. JUICIO: Complete por checks de código no equivale a archive habilitado |

## APPENDED — por documento

| Documento | Historia agregada / costo |
|---|---|
| README.md | Ninguna historia nueva; REWRITE del estado y guía del ciclo. Decisiones, lecciones y deuda enlazadas a docs/metodologia-opsx.md |
| ARCHITECTURE.md | Ninguna entrada histórica nueva; REWRITE del diseño/estado; conserva decisiones anteriores y enlaza la historia de fase 1 |
| AGENTS.md / CLAUDE.md | Ninguna entrada histórica nueva; REWRITE operativo pareado, historia enlazada a la base |
| CHANGES.md | Registro fechado del cierre excepcional, condición CONF-1 y próximos pasos; costo tres vueltas y corrección extra. EVIDENCIA informes 21–23 |
| docs/propuestas/2026-10-23-metodologia-proptech-en-opsx.md | Registro fechado del archive parcial; costo corrección fuera del bucle y deuda conservada. EVIDENCIA informes 21–23 |
| docs/metodologia-opsx.md | D-001…D-003; lecciones locales L-001…L-004; deuda consolidada íntegra del informe 21 y RF3-1 cerrado. Cada entrada con fecha, costo y EVIDENCIA/JUICIO/SUPOSICIÓN; no se renumera historia externa PropTech |
| openspec/changes/instalacion-sana/tasks.md | Checklist pendiente de cuatro puntos CONF-1, sin tildar trabajo no ejecutado. EVIDENCIA §3.6 e informe 22 |
| knowledge-base-report.md en cada uno de los ocho changes | Sólo APPEND/creación: son registros de cierre nuevos, no contratos actuales que reescribir. Incluyen ambas listas y PENDING; costo documental de conservar trazabilidad por change |
| docs/propuestas/informes/2026-10-23/24-archive.md | Sólo APPEND/creación: informe nuevo de esta operación; no existe estado anterior a sustituir. No es SUSPICIOUS: append-only de la base, que sí recibió REWRITE en seis documentos existentes |

## Lecciones y deuda: procedencia y límites

**EVIDENCIA:** informe 21 no tiene sección separada LECCIONES, pero sí el juicio explícito sobre criterio literal/objetivo y los límites de los tests; informes 11, 14, 17 y 20 sostienen las lecciones de texto. L-003 incorpora además la observación directa solicitada por Juani: informe 11 prueba pérdida de miles de líneas y autorreferencias, pero no demuestra por sí solo la causa `codex exec -o`. Se tipa como instrucción/observación humana y JUICIO, no como reproducción del archivista. L-004 usa el comando y error EROFS del informe 19; no afirma que todo sandbox impida Go.

**PENDING:** no había archive previo ni knowledge-base-report.md anterior al iniciar (`openspec/changes/archive/` no existía). No se pudo arrastrar deuda desde un change previo; no se equipara a deuda vacía. La tanda encadena sus reportes, sin revisar todo el archivo histórico. Conserva CONF-1; CONF-3/RF2-3; RF2-4; RF1-4; RF1-6; N1, N1b, N2b, N3, N4, N5, N8; RF3-2. RF3-1 queda separado como cerrado. Asignaciones y costo/riesgo por entrada en docs/metodologia-opsx.md; identificadores futuros no decididos = (TBD). No se reabre el bucle ni se corrige producto en esta fase.

## Ejecución y comprobación

### EVIDENCIA de ejecución — 2026-10-23, controles terminados 11:03:37 -03

Se ejecutó, uno por uno y sin flags de omisión de validación, `openspec archive <change> --yes` para los ocho nombres de la tabla de alcance, en su orden. Ocho exit 0; cada salida confirmó `Change '<change>' archived as '2026-10-23-<change>'.` Los ocho directorios quedaron en `openspec/changes/archive/2026-10-23-<change>/`, conservando `knowledge-base-report.md` y copias íntegras de informes 21/22. El CLI está disponible; no se movieron directorios manualmente.

Los ocho comandos emitieron la advertencia no bloqueante de propuestas sin `## Why` / `## What Changes` y recomendación de deltas specs. No se ocultó ni se desactivó validación. No había delta specs en estos changes: el CLI archivó artefactos sin anunciar sincronización de specs. JUICIO: se conserva el formato histórico de propuestas referenciadas a la propuesta firmada, sin inventar deltas ni aprobación nueva.

`openspec list` → sólo `instalacion-sana 4/8 tasks`. Las cuatro tareas pendientes hacen visible CONF-1, sin destildar las cuatro pruebas de código ya realizadas. `openspec status --change instalacion-sana` → raíz activa, 2/4 artefactos (proposal/tasks completos; specs/design ausentes). JUICIO: ese contador de artefactos no elimina la condición humana de cierre.

Control documental, fuente de código `5ab0f4f` con diff sólo documental de esta entrega:

```text
go test ./internal/harness/config -run '^(TestRootDocs_ToggleEnumerationsAreNotStale|TestDocsSeniorSoloVerify)$' -count=1 -v
=== RUN   TestRootDocs_ToggleEnumerationsAreNotStale
--- PASS: TestRootDocs_ToggleEnumerationsAreNotStale (0.00s)
=== RUN   TestDocsSeniorSoloVerify
--- PASS: TestDocsSeniorSoloVerify (0.00s)
PASS
ok github.com/Group-Active-IA/active-stack/internal/harness/config 0.004s
```

Exit 0. `git diff --check` → sin salida, exit 0. Se había ejecutado además un filtro inicial de tests de descripción, exit 0; el comando explícito de arriba verifica que corrieron los dos contratos documentales pertinentes. No se ejecutó suite completa nueva: sólo cambió documentación; la suite de producto citada sigue siendo la registrada en informes 22/23. No se ejecutó `go build` ni CI remoto. La nueva redacción “archivado” y este registro de resultados se revisaron después del CLI; no cambian descripciones de toggles probadas.

**EVIDENCIA de alcance final:** ocho reportes preservados en archives, una única raíz activa con CONF-1. Los archivos ignorados de OpenSpec no forman parte del commit; este informe versionado conserva nombres, comandos, salidas y limitaciones. Las únicas modificaciones versionadas previstas son README, ARCHITECTURE, CHANGES, CLAUDE, AGENTS, propuesta, base metodológica y este informe. Se relee el reporte inmediatamente anterior en cada salto y se conserva toda la lista PENDING; no se afirma deuda vacía.

Control del índice: `git diff --cached --check` detectó una línea vacía final en este informe, no incluida por el control inicial de archivos sin seguimiento. Se retiró antes de enmendar el commit; el control del índice final quedó sin salida. No se cambió producto.
