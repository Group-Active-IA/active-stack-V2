Revisión arquitectónica sobre `feat/metodologia-fase-2`. Sólo lectura; no leí otros informes de fase 2. “BLOQUEANTE” señala algo a resolver antes de firmar, sin atribuirme veto.

**AR-1 · BLOQUEANTE · Originales fuera de Git frente a aprobación anclada en Git**
EVIDENCIA, 23/10/2026 11:20 -03: `git check-ignore openspec/changes/metodologia-fase-2/informes/example.md` devolvió esa ruta.
La propuesta exige que P contenga documentos e informes revisados en `docs/propuestas/2026-10-23-metodologia-fase-2.md:15`, pero captura los originales bajo `openspec/` en `:27`.
JUICIO del arquitecto, confianza alta: falta definir la frontera `reports → material documental versionado → approval`. Con el destino actual, capturar correctamente no alcanza para construir P.
Antes de firmar: elegir exportación determinista de originales a un destino versionado o incorporación explícita y acotada al commit; definir quién valida igualdad, procedencia y conflictos.
Registrar alternativa descartada y costo de revertir: migrar referencias/manifiestos y reaprobar los alcances cuyo material firmado cambie.

**AR-2 · BLOQUEANTE · Referencias dinámicas sin contrato de composición**
EVIDENCIA, misma fecha/hora: `nl -ba internal/harness/config/compose.go` mostró que `Compose(toggles, variantKey)` no recibe rutas ni manifiestos y declara composición pura en `internal/harness/config/compose.go:41`.
La propuesta promete ruta efectiva e índice con versión/hash en `docs/propuestas/2026-10-23-metodologia-fase-2.md:43`; hoy `internal/harness/config/install.go:63` compone y luego inyecta.
JUICIO del arquitecto, confianza alta: falta decidir dónde se enlazan los recursos instalados con el texto generado, incluyendo project/global y verificación contra Compose.
Antes de firmar: nombrar el responsable y dirección de dependencia, por ejemplo installer → resolución validada → render puro → inyección.
Contrastar esa opción con referencias relativas estables y registrar el costo de revertir: contratos de composición, goldens, chequeo instalado y referencias ya entregadas.

**AR-3 · NOTA · Orden declarado, dependencia todavía sin resolver**
EVIDENCIA, misma fecha/hora: `nl -ba internal/planner/graph.go` mostró que el grafo toma exclusivamente `Harness.DependsOn` en `internal/planner/graph.go:11`.
La propuesta exige instalar recursos antes de publicar referencias en `docs/propuestas/2026-10-23-metodologia-fase-2.md:45`, pero no fija la arista ni cómo depende del toggle `roles`.
JUICIO del arquitecto, confianza media: explicitar `orquestador con roles → role-definitions`, o una precondición equivalente del plan, permitiría revisar selección custom y ausencia de ciclos sin depender del orden incidental.
La separación por ID y la recuperación local del paso fallido sí están contempladas; coinciden con `internal/harness/config/install.go:18` e `internal/pipeline/rollback.go:32`.

No ejecuté tests ni build. No verifiqué contratos externos de proveedores; estos hallazgos se limitan a las fronteras documentadas y al código leído.