# Metodología OPSX vigente — fase 1

Estado al 2026-10-23: ocho changes archivados con el CLI; `instalacion-sana` permanece abierto por CONF-1.

El ciclo vigente es `explore → propose → apply → verify → archive`, seguido de `ship` cuando se publica:

- **Propose:** con governance activo, ALTO o CRÍTICO requiere informes de arquitecto, analista funcional y QA, más especialistas por capa, convocados de a uno sin leer los otros informes; el auditor lee los N originales antes de la aprobación humana escrita. La tabla `## Aprobación` registra elección, persona, fecha y alcance; ningún agente la completa. Una regla nueva necesita aprobación propia.
- **Apply constructor:** rojo y verde por tarea en `tasks.md`, con comando, salida, hora y SHA; si el test nació verde, mutante que falle en su aserción y se revierta. Declara consumidores, conflictos y familia del modelo. En vueltas 2 y 3 corrige sólo la lista del auditor.
- **Verify:** conformidad → rompedor con tres lentes (correctitud, cableado, documentación) → senior si está activo → auditor. La conformidad sola no habilita archive. El auditor decide y el orquestador ejecuta; sólo APROBADO con informes completos y conformidad satisfecha habilita el cierre ordinario. El contador y los informes viven en el change: máximo tres vueltas, sin cuarta automática. En la tercera con bloqueantes se exige decisión humana escrita sobre alcance nuevo, deuda justificada o freno; no se fabrica APROBADO.
- **Archive:** lee la unión del impacto declarado y el diff; reescribe las descripciones actuales y agrega decisiones, lecciones numeradas y deuda del auditor con fecha, costo y evidencia tipada. Deja `knowledge-base-report.md` por change con REWRITTEN, APPENDED y PENDING por documento; sólo APPEND se justifica y se señala como `SUSPICIOUS: append-only` cuando correspondía reescribir.
- **Ship:** espera todo el CI requerido sobre el SHA pusheado, inspeccionando workflows, eventos y jobs. En active-stack exige `ci.yml` en push con tests Ubuntu y Windows; en otros repos deriva su conjunto. Cero runs, consulta fallida, SHA ajeno o estados incompletos no confirman entrega. `main` exige permiso y nunca se mergea el PR automáticamente.

Si falta un toggle se declara la degradación; sin review no hay rompedor ni bucle. Sin auditor se necesita decisión escrita de Juani para cerrar con alcance degradado. Un informe obligatorio faltante detiene el cierre.

**EVIDENCIA — 2026-10-23:** contratos implementados en `f3fc02a`, `db3de54`, `66ad39e`, `84c570d`, `9f2d03b`, `531e1b6`, `60b5e83`; degradación sin review corregida en `0b5d0d4`. El [OK escrito](propuestas/informes/2026-10-23/22-decision-juani-rf3-1.md) acepta el archive tras `ce7899e`. **JUICIO:** estos textos describen instrucciones instalables; no prueban obediencia futura ni instalación real. La [deuda del auditor](propuestas/informes/2026-10-23/21-auditor-vuelta-3.md) limita su cobertura.

## Instalación sana y versión

`active-stack version` imprime `active-stack <versión> (SHA <revisión>)` usando metadatos Go; sin ellos muestra `devel` y `desconocido`. El chequeo del binario valida la identidad `cmd/active-stack` antes de ejecutar `roles --help` (salida 0). La entrega compara el bloque con `Compose` para variante y toggles efectivos y rechaza markers duplicados en Claude/Codex. OpenCode usa `command` como array y `environment` para variables MCP locales.

**EVIDENCIA — 2026-10-23:** `84fcf0e`, `c325955`, `8c6a06b`; fuentes `cmd/active-stack/version.go`, `internal/verify/{binary,delivery}.go`, `internal/harness/external/mcp.go`. **JUICIO:** comprobaciones implementadas no equivalen a reinstalación acreditada.

CONF-1 exige registrar §3.6 de la [propuesta](propuestas/2026-10-23-metodologia-proptech-en-opsx.md): (1) SHA, versión, PATH, destinos por adapter, variante y toggles; (2) snapshot legible con ID, ruta, hashes y permisos de binario, instrucciones, comandos y settings; (3) ensayo aislado con PATH/binario aislados, reinstalación, falla inducida y restore comparado por contenido y permisos; (4) OK escrito separado de Juani y reinstalación real con snapshot propio y aceptación por destino, preservando texto fuera de markers. OpenCode recupera instrucciones y settings conjuntamente. No se atribuye el `.test` a un test concreto sin reproducción.

## Decisiones — historial acumulativo

- **D-001 · 2026-10-23 · EVIDENCIA:** se integra la metodología en toggles existentes, sin toggle nuevo; aprobación inicial `4748680`, propuesta §5. Costo: cambios secuenciales sobre fragmentos y dos goldens compartidos.
- **D-002 · 2026-10-23 · EVIDENCIA:** Juani elige la opción A de RF2-1 (`5e8798f`, informe 16): fortalecer contratos dentro del alcance; lectura semántica obligatoria nueva queda candidata a fase 2, con firma propia. Costo: tercera vuelta y cobertura aún limitada de paráfrasis.
- **D-003 · 2026-10-23 · EVIDENCIA:** vuelta 3 NO APROBADA (`18bb293`, informe 21); decisión B y OK posterior (`e46294a`, `5ab0f4f`, informe 22) permiten corrección externa al bucle `ce7899e` y archive de ocho changes. Costo: corrección y verificación independiente adicionales; no hay vuelta 4 ni APROBADO retroactivo. `instalacion-sana` sigue abierto.

## Lecciones — historial local, no renumera PropTech

La serie L-001… de este documento pertenece a active-stack-V2; las L-029/L-037 de la propuesta son referencias externas de PropTech. Las entradas son JUICIO derivado de las fuentes indicadas, no nuevas reglas aprobadas.

- **L-001 · 2026-10-23 · JUICIO:** `strings.Contains` puede atrapar la frase sin proteger la regla: borrar la frase da rojo, agregar su contradicción o parafrasearla puede dar verde. EVIDENCIA: informes 10, 14 y 20; auditorías 11, 17 y 21; `ce7899e` cierra el caso sin puntuación. Costo observado: tres vueltas y corrección extra; RF3-2 y otras paráfrasis siguen abiertas. Un golden regenerado no sustituye una prueba de sentido.
- **L-002 · 2026-10-23 · JUICIO:** un criterio de terminado literal puede cumplirse sin alcanzar el objetivo. EVIDENCIA: el informe 21 recoge la admisión del auditor (“es la segunda vez”); el criterio previo está en el informe 17; en 21, N6 sin punto evade el criterio ya satisfecho para N6/N7. Costo: cierre rechazado dos veces y decisión humana adicional.
- **L-003 · 2026-10-23 · JUICIO:** `codex exec -o` sobre el archivo donde el agente escribió su informe lo pisa con el retorno final; separar destino del retorno e informe conserva ambos. EVIDENCIA de procedencia: instrucción directa de Juani al archivista en esta sesión (23/10) y pérdida observada por el auditor en informe 11: el 08 pasó de 751 a 13 líneas, el 09 perdió unas 1145, ambos enlazaban a sí mismos; los originales sobrevivían en los changes. El informe 11 acredita la pérdida, no atribuye por sí solo la causa a `-o`; no se reproduce aquí. Costo: pérdida/sustitución del informe y necesidad de recuperarlo; magnitud no medida. SUPOSICIÓN: qué ejecuciones concretas del ciclo sufrieron esa pérdida, sin prueba primaria aquí.
- **L-004 · 2026-10-23 · JUICIO:** revisores con sandbox de sólo lectura no pueden ejecutar Go cuando su directorio temporal es de sólo lectura; la revisión estática debe declarar ese límite. EVIDENCIA: informe 19, 10:45:46 -03, `go test ./assets ./internal/harness/config ./internal/verify ./internal/catalog ./cmd/active-stack` → `mkdir /tmp/go-build159821652: read-only file system`; informe 21 lo confirma. Costo: conformidad sin ejecución propia, necesita evidencia de otro entorno. No se generaliza a todo sandbox que permita temporales escribibles.

## Deuda consolidada — estado al 2026-10-23

Origen de todas las entradas abiertas: DEUDA del informe 21 (`18bb293`); detalles en informes 10, 13 y 14. Fecha de arrastre: 2026-10-23. Asignación fuera de la explícita del auditor: JUICIO del archivista, al primer change que use la pieza; identificador de change futuro (TBD). No se da por pagada deuda por este archive.

| Ítem | Estado y costo/riesgo que se conserva | Destino / evidencia tipada |
|---|---|---|
| CONF-1 | ABIERTA, condiciona archive de `instalacion-sana`; faltan receta y evidencia §3.6 puntos 1–4 y OK del 4; instalación/restore reales no acreditados | `instalacion-sana`; EVIDENCIA informes 19, 21, 22 |
| CONF-3 / RF2-3 | ABIERTA; senior exige auditor aun apagado, puede esperar un veredicto imposible | `senior-solo-en-verify`, archivado con deuda; próximo consumidor de senior sin auditor (TBD); EVIDENCIA informes 13, 14, 21 |
| RF2-4 | ABIERTA; mayúscula `Mark` en medio de oración fijada por test; incluye restos de RF1-4, sin duplicar deuda | próximo cambio de KB/review (TBD); EVIDENCIA informes 14, 21 |
| RF1-4 | ABIERTA; restos “does not block / belongs to the person” junto al bucle con autoridad del auditor | próximo cambio de review (TBD); EVIDENCIA informes 10, 11, 21 |
| RF1-6 | ABIERTA; claves Engram por fase frente a claves por vuelta/rol, riesgo de contaminación | próximo consumidor del protocolo Engram (TBD); EVIDENCIA informes 10, 21 |
| N1, N1b, N2b, N3, N4, N5, N8 | ABIERTAS; paráfrasis de opción A sin cobertura semántica general acreditada, falsa confianza de tests textuales | próximo cambio de contratos (TBD); EVIDENCIA informes 14, 16, 21 |
| RF3-2 | ABIERTA; “treated as green” y “With zero runs, treat delivery as green” pasan `TestShipEsperaCICompleto` | próximo cambio de ship/contratos (TBD); EVIDENCIA informes 20, 21 |
| RF3-1 | CERRADA fuera del bucle; costo: corrección de una línea más mutantes en bloque y comando, verificación independiente y OK escrito | EVIDENCIA `ce7899e`, informes 22 y 23; decisión `5ab0f4f` |

**PENDING documental:** no se puede arrastrar deuda desde un archive previo: no existía `openspec/changes/archive/` al iniciar. En esta tanda los reportes se encadenan y conservan esta deuda. No equivale a afirmar ausencia de deuda histórica del proyecto.
