**JUICIO (arquitecto, confianza alta, contra el “terminado” del rol): faltan contratos de dependencia y autoridad antes de firmar.** BLOQUEANTE indica qué resolver antes de la firma, no un veto mío. Referencias `P` = `docs/propuestas/2026-10-23-metodologia-proptech-en-opsx.md`.

**EVIDENCIA (2026-10-23, 09:38 ART):** `git branch --show-current` → `feat/metodologia-proptech`; `git status --short` → sólo `?? docs/propuestas/informes/`. No consulté esos informes ni escribí archivos.

**AR-1 · BLOQUEANTE · Dependencias opcionales sin contrato.**
**EVIDENCIA:** `nl -ba internal/harness/config/compose.go` → líneas 64-74 validan nombres; 119, 142, 158 y 173 agregan auditor, review, KB y senior independientemente. `nl -ba P` → 85 y 88 hacen depender bucle, lecciones y deuda del auditor.
**JUICIO (arquitecto, alta):** el orden textual no garantiza presencia. Definí `review → auditor`, `KB → veredicto` y las referencias a senior/TDD como dependencias obligatorias o condicionales, con salida explícita cuando falta el productor. Probá esas combinaciones, además de todos los toggles juntos.

**AR-2 · BLOQUEANTE · Dos autoridades para detener verify.**
**EVIDENCIA:** `nl -ba internal/harness/config/assets/senior.md` → 132-134 exige presentar cualquier BLOCKING a la persona y esperar; `nl -ba P` → 85 decide el auditor, pero 87 conserva “the senior labels, you hold”.
**JUICIO (arquitecto, alta):** senior antes del auditor es coherente como flujo de informes, pero falta resolver quién prevalece ante desacuerdo. Indicá si el auditor clasifica también los hallazgos del senior o si éstos exigen intervención humana inmediata. Los pasos 4 y 5 deben contemplar el estado intermedio.

**AR-3 · BLOQUEANTE · La compactación pierde destinos de apply/archive.**
**EVIDENCIA:** `nl -ba P` → 64-68 conserva esas reglas en el bloque porque sus comandos pertenecen a OpenSpec; 135 propone mover detalle a verify, ship y roles, sin nombrar un destino para apply/archive.
**JUICIO (arquitecto, alta):** falta dibujar `orquestador → prompt del constructor/archivista → reglas completas`. Conservá allí lo necesario o definí cómo se carga y transmite. La compactación no puede depender de que esos agentes lean comandos de otra fase.

**AR-4 · BLOQUEANTE · Alternativas y reversión incompletas.**
**EVIDENCIA:** `nl -ba P` → 168-176 registra alternativa y costo para el toggle separado; 178-179 rechaza comandos propios sin costo de reversión; 131-135 introduce approval, hook, harness agent y compactación sin fichas equivalentes.
**JUICIO (arquitecto, alta):** completá para cada decisión de frontera: ubicación, consumidor → proveedor, alternativa descartada y costo de volver atrás. Puede quedar `(TBD)` para su propose posterior, pero entonces la firma actual debe excluir ese diseño pendiente.

**AR-5 · NOTA · Los goldens pueden bendecir una contradicción.**
**EVIDENCIA:** `sed -n '1,87p' internal/harness/config/compose_test.go` → `-update` escribe y retorna; sólo compara Claude/OpenCode con todos los toggles. Reconstrucción estática mediante `python3` y `Path.read_text()` → ambos goldens coinciden, 1959/1961 líneas; existen ocho bases.
**JUICIO (arquitecto, alta):** serializar evita conflictos de archivos, no errores semánticos. Revisá el diff antes de aceptar la regeneración y fijá con tests separados autoridad, combinaciones parciales y modelo de archive en las ocho variantes.

**JUICIO (arquitecto, confianza indicada, sobre las cinco recomendaciones):**

1. **Un rompedor, tres lentes: favorable, media.** Declarar pérdida de independencia y mantener tres secciones verificables. En vueltas posteriores, incluir consumidores afectados por la corrección, aunque estén fuera del diff.
2. **Senior antes del auditor: favorable, alta**, condicionado a AR-2. Flujo: rompedor → senior → auditor → decisión del orquestador; orden de lectura y orden de ejecución son contratos distintos.
3. **Firma en fase 2, configurada por Juani: favorable, media.** Mientras tanto, aprobación humana declarada sin garantía criptográfica. El diseño posterior debe vincular clave autorizada, alcance aprobado y commit verificado.
4. **Modelo declarado, sin control de la herramienta: favorable, alta.** Evita sumar infraestructura para una opción; la independencia queda como SUPUESTO salvo evidencia del despacho.
5. **Archive a sonnet: favorable, media.** La tarea propuesta requiere juicio; no tengo medición que pruebe mejora. Registrar haiku como alternativa y reversión en ocho bases, tests y dos goldens.

**SUPUESTO / límite:** revisión estática del diseño y sus artefactos; no ejecuté tests, build, hooks, firma ni OpenSpec. No afirmo ausencia de ciclos en componentes futuros cuyo contrato todavía falta.