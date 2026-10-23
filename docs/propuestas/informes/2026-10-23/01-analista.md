JUICIO (analista funcional; confianza alta, contra el “terminado” del rol): faltan escenarios para firmar. BLOQUEANTE significa resolver antes de firmar; la decisión sigue siendo de Juani.
EVIDENCIA (2026-10-23): `git branch --show-current` → `feat/metodologia-proptech`; `git diff --stat` → vacío. Leí el rol entero; no consulté otros informes ni ejecuté builds/tests.
Referencias: P = `docs/propuestas/2026-10-23-metodologia-proptech-en-opsx.md`; E = `docs/exploraciones/2026-10-23-metodologia-proptech.md`; C = `internal/harness/config/`.

- **AN-1 BLOQUEANTE | P:85-87. Autoridades contradictorias.**
  EVIDENCIA (23/10): `nl -ba C/assets/senior.md` → 132-134 ordena presentar todo BLOCKING a la persona y esperar; P:87 conserva esa regla, mientras P:85-86 hace decidir al auditor.
  JUICIO (analista; alta): ante un bloqueante del senior y APROBADO del auditor, el orquestador recibe dos salidas incompatibles. Escribir el escenario: actor orquestador, disparador desacuerdo, resultado observable y precedencia. Hoy falta esa resolución.

- **AN-2 BLOQUEANTE | P:33-35,85,187-188. El contador y el alcance no tienen contrato único.**
  EVIDENCIA (23/10): `nl -ba P` → 33 permite volver desde vueltas 1 y 2; 84 habla de corregir en vueltas 2 y 3; 188 limita la revisión a lo corregido.
  JUICIO (analista; alta): definir qué incrementa la vuelta, cómo se recupera tras compactación y qué decisiones puede tomar la persona al tercer cierre. Caso actual sin respuesta definida: aparece una regresión fuera de las líneas corregidas en vuelta 2. Debe reportarse y clasificarse, sin perderse por el recorte.

- **AN-3 BLOQUEANTE | P:85,88,90,146. Combinaciones de toggles incompletas.**
  EVIDENCIA (23/10): `nl -ba C/compose.go` → 119,142 activa auditor y review por separado; `nl -ba C/review_cost_rules_test.go` → 47-49 contempla sin review, auditor o evidence.
  JUICIO (analista; alta): para el orquestador que ejecuta verify con review activo y auditor apagado, falta resultado observable: ¿detiene, degrada o termina sin aprobación? También falta quién produce DEUDA/LECCIONES al archivar sin auditor. Fijar escenarios por combinación relevante.

- **AN-4 BLOQUEANTE | P:90; assets/commands/claude/opsx/verify.md:137-140. Conformidad y aprobación no están reconciliadas.**
  EVIDENCIA (23/10): `nl -ba assets/commands/claude/opsx/verify.md` → CRITICAL exige corregir antes de archive; sin críticos declara “Ready for archive”. P:90 sólo agrega un paso 9.
  JUICIO (analista; alta): definir qué hace el auditor con CRITICAL de conformidad y dónde queda la única decisión final. Caso que hoy falla: conformidad limpia anuncia listo sin haber corrido rompedor/auditor. El escenario debe impedir esa habilitación anticipada.

- **AN-5 BLOQUEANTE | P:84-94; E:71-74,104-116. Reglas retenidas sin destino verificable.**
  EVIDENCIA (23/10): `rg -n 'identidad|teclado|L-032|leen antes|lectura previa|R17|R13' P` → sin coincidencias; E mantiene verificación humana con proveedor externo y lectura previa de lecciones.
  JUICIO (analista; alta): falta decidir inclusión o exclusión explícita. Escenarios: constructor/rompedor leen lecciones antes de actuar; orquestador, ante identidad externa no recorrida, obtiene verificación humana o registra deuda. Hoy P no especifica esas salidas.

- **AN-6 BLOQUEANTE | P:89,103. “CI completo” carece de criterio de completitud.**
  EVIDENCIA (23/10): `nl -ba P` → 89 propone listar runs por SHA y observarlos, sin definir workflows esperados ni tratamiento de cancelados, omitidos o pendientes.
  JUICIO (analista; alta): actor orquestador, disparador push, resultado listo sólo al completar el conjunto esperado. Caso sin resolver: un workflow verde y otro pendiente sobre el mismo SHA. Agregar escenario negativo; listar un run verde no demuestra completitud.

- **AN-7 BLOQUEANTE | P:94,109. Independencia de propose y cantidad de informes ambiguas.**
  EVIDENCIA (23/10): `nl -ba P` → 94 pide al menos tres especialistas, adicionales según capas y luego auditor “sobre los cuatro”; 109 sólo fija presencia de roles y anterioridad.
  JUICIO (analista; alta): precisar qué informes recibe el auditor y cómo evita el siguiente especialista leer los ya guardados o sus memorias. Caso actual: segundo revisor abre el informe del primero en el change. Debe declararse contaminación, no independencia.

- **AN-8 NOTA | P:99-112,122-135. Falta trazabilidad uniforme al criterio de terminado.**
  EVIDENCIA (23/10): `nl -ba C/compose_test.go` → 59,83 compara salida con goldens; P:109 identifica un rojo concreto para multirol, pero las demás filas no presentan escenarios completos.
  JUICIO (analista; alta): completar por regla actor, disparador, salida y caso que hoy falla, incluidos pasos 0 y 8-12. Los tests textuales verifican el bloque emitido; no prueban que un agente cumpla el proceso.

JUICIO (analista; confianza alta salvo donde indico), sobre las recomendaciones provisorias:
1. **Un rompedor con tres lentes:** aceptable como cobertura, sin llamarlo tres revisiones independientes. Exigir resultado por lente aun sin hallazgos; en vueltas posteriores incluir consumidores y regresiones.
2. **Senior antes del auditor:** favorable, condicionado a resolver AN-1. Una decisión vieja aporta contexto; falta definir cuándo justifica deuda y cuándo sigue existiendo un defecto.
3. **Firma en fase 2, configurada por Juani:** razonable por etapas. Antes de firmar, definir qué canal acepta el constructor mientras falta la verificación y qué hace si no puede validarlo. P:131 sólo menciona `git verify-commit`, insuficiente como escenario de autorización del alcance.
4. **Modelo declarado en informe:** favorable si queda como SUPUESTO cuando no hay registro verificable; indicar modelo/proveedor declarado y si comparte familia con los revisores. No presentarlo como independencia comprobada.
5. **Archive a sonnet:** favorable como preferencia (confianza media), sin garantía de calidad. EVIDENCIA (23/10): `nl -ba C/assets/antigravity/sdd-orchestrator.md` → 434 todavía ordena dedicar menos esfuerzo a archive; cambiar también esa instrucción contradictoria.

SUPUESTO: evalué contratos documentados y contenido estático; no observé ejecuciones reales de agentes ni comprobé los incidentes históricos citados por la exploración.