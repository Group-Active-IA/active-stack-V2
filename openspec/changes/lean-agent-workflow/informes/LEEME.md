# Informes completos y procedencia

## Apply / roledefs

- Original completo: [apply-roledefs.md](apply-roledefs.md), copia textual íntegra de `/tmp/roles-impl-report.md`; no resumido ni suavizado.
- Autor: subagente `roles_impl`, scope de dispatch `internal/harness/roledefs/`.
- Proveedor disponible: OpenAI. Modelo heredado del orquestador; identidad exacta no consultada ni acreditada independientemente. Familia compartida con el orquestador; no se usó proveedor alternativo y no se declara independencia entre proveedores.
- Revisión: HEAD `61887a442ecf8abed1f5106915cf2d56e6bbb278` más diff sin commit de roledefs descrito en el informe. Las pruebas no acreditan un commit entregado nuevo.
- Mediación: el orquestador escribió el encargo y solicitó fortalecer las comprobaciones del auditor; el informe completo se conserva sin editar. Este índice y las entradas de tasks son consolidación de `roles_impl`, no un nuevo veredicto del auditor.
- Fuentes y contaminación: declaradas en el informe original. Instrucciones/contexto heredados y filesystem compartidos; no aislamiento garantizado. No leyó conclusiones de otros revisores.
- Faltantes: informes de apply de composición y catálogo/comandos/docs, y review de plan/conformidad/rompedor/auditor aún no recibidos en este índice. No se fabrican aprobaciones ni cierre.

## Dispatch

EVIDENCIA (transcripción de asignaciones del orquestador; horas de lanzamiento no capturadas):

- `roles_impl`: exclusivamente `internal/harness/roledefs/`, cinco roles nuevos + siete especialistas sin modelo fijo, pruebas y documentación local. Después consolidó su evidencia y este índice, sin tocar implementación de otras zonas.
- `map_workflow`: composición/configuración del perfil lean; scope `internal/harness/config/` según dispatch del orquestador.
- `propose_workflow`: propuesta inicial y posteriormente catálogo/comandos/documentación, según dispatch del orquestador.
- El orquestador indicó trabajo independiente en paralelo después del plan. No se declaran horarios o simultaneidad medidos más allá de las asignaciones conservadas en la conversación.

EVIDENCIA remitida por el orquestador (fecha 2026-10-29; hora no registrada): `$ command -v active-stack` → exit 1, sin salida. `active-stack roles` no estaba disponible; no se ejecutó vía `go run` o build por regla del repositorio. El ledger NO procede de una lista generada por ese comando.

JUICIO (orquestador, conforme al nuevo flujo explícitamente acordado con Juani): los alcances de QA/conformidad, rompedor y auditor se asignarán por el flujo de cinco roles. Senior se omite porque no es etapa del flujo acordado. Los siete especialistas quedan como criterios opcionales; su presencia instalada no significa asistencia obligatoria. Confianza alta sobre el acuerdo conversacional; los informes de revisión todavía faltan.

## Originales agregados al corte

- apply-compose.md — retorno íntegro constructor configuración, mediación: copia sin resumir desde /tmp/compose-impl-report.md.
- apply-integration.md — retorno íntegro integración, mediación: copia desde /tmp/integration-impl-report.md.
- integration-checks.md — comprobaciones integradas, copia desde /tmp/lean-integrated-checks.md; logs originales en logs/.
- Revisiones de conformidad, breaker y auditor: NO REALIZADAS al pedido de corte. Mismo proveedor/familia compartidos en la implementación, sin independencia entre proveedores acreditada.

Informe general en la raíz del repositorio: INFORME-ESTADO-2026-10-29.md.

## Grupo 5 — installer meets previous installs

| Informe | Autor | Ejecutor | Revisión | Mediación |
|---|---|---|---|---|
| [apply-1-constructor.md](apply-1-constructor.md) | constructor | subagente Claude general-purpose | diff de trabajo sobre `61887a4` (sin commit) | escrito por el propio constructor; sin resumir. Tareas 5.1–5.4. Revisión de rompedor/auditor: NO REALIZADA aún |

## 29/10 tarde — trabajo directo por decisión de Juani

| Trabajo | Autor | Ejecutor | Evidencia |
|---|---|---|---|
| Tareas 4.1–4.7 (fases, LOW, encargos autosuficientes, ejecutores por proyecto, formato de informe, profundidad por nivel, 11 mecanismos legacy) | orquestador | sesión principal Claude | bajo cada tarea en tasks.md |
| Comparación legacy vs lean | agente Explore de sólo lectura | subagente Claude | resumida en tasks.md 4.7; transcripción no copiada a este directorio; hallazgos 1, 2 y 4 reverificados por el orquestador |
| Tareas 6.1–6.3 y docs 3.1 | orquestador | sesión principal Claude | bajo cada tarea en tasks.md |

DESVÍO DECLARADO: por pedido explícito de Juani no se convocaron revisor del plan, rompedor ni auditor para los grupos 4, 5 y 6. El orquestador implementó código (4.x, 6.x), cosa que el flujo reserva al constructor. Familia/proveedor compartido en todo el trabajo; sin independencia entre proveedores.
