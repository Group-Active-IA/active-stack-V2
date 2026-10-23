# Informe del analista funcional (role-analyst, opus) · 2026-10-23

Extraído textual del transcript del sub-agente con jq (último mensaje de texto), sin transcripción del orquestador.

Análisis de las 14 piezas como analista funcional, sin editar ni guardar nada (2026-10-23).

**Resumen:** 4 piezas quedan como requisito verificable (D, G, K, N). Otras 3 se pueden verificar con una condición previa (A, C, L). Las otras 7 no son verificables tal como están escritas (B, E, F, H, I, J, M). Hay dos contradicciones fuertes: C contra `role-implementer`, y la de A y H con la regla del bloque de lanzar en paralelo lo que no depende de nada.

**Archivos que abrí:** PropTech `knowledge-base/14_forma_de_trabajo.md` y `15_lecciones.md` (enteros), `.git/informes-2b/reglas-2b.md`, `.git/informes-2b/auditor/informe.md`, `proposal.md:263-293` del 2b; `~/.claude/agents/role-implementer.md:107-205`. En active-stack-V2: `assets/strict-tdd.md` (entero), `assets/claude/sdd-orchestrator.md:270-429` y `compose.go:40-79`. En `review.md`, `ship.md`, `roles.md`, `auditor.md`, `senior.md` y `governance.md` sólo hice búsquedas (rg). **No abrí** `classify.go`, `roles_dispatch.go` ni los fragmentos completos.

## EVIDENCIA (lecturas del 2026-10-23; no corrí ningún comando de shell)
- E1. `strict-tdd.md:64-73`: la evidencia de TDD es una tabla de tildes ("✅ Written", "✅ Passed") en el informe. No pide hora, salida, `tasks.md` ni mutante. rg de `tasks.md|mutant|timestamp` en ese archivo: 0 coincidencias.
- E2. `ship.md:170`: "the cycle ends at the open PR". `:172-173`: los checks en verde no son aprobación. No hay ningún paso que espere el run de CI de ese SHA.
- E3. `claude/sdd-orchestrator.md:316-317`: los agentes independientes se lanzan "in the same message". `:336`: "Done when `active-stack roles` was run".
- E4. rg de `active-stack roles` en `~/.claude/CLAUDE.md`: 0. rg en `/home/juanisarmiento/Proyectos/PropTech`: 0 archivos.
- E5. `role-implementer.md:178-182`: "No escribís código de autenticación... **no existe el canal**". `:186-192`: ningún archivo que pueda escribir un agente sirve como autorización.
- E6. `proposal.md:265`: "La completa, a mano, uno de los dos socios. Ningún agente la escribe". `:290-293`: tabla Socio/Fecha/Alcance.
- E7. rg de `lesson|lecci|model family|codex|browser|round|.git/informes` en los assets: ninguna regla. Sólo menciones de paso (`roles.md:80`, `senior.md:202`).
- E8. `compose.go:45-59`: los fragmentos se agregan enteros según el toggle. No hay selección por fase.
- E9. `claude/sdd-orchestrator.md:305-306`: el propio repo mide un `CLAUDE.md` de 1.702 líneas, 89% del bloque, heredado por 39 agentes.
- E10. Auditor del 2b, fila 7: de las revisiones del propose "no quedó el informe original, sólo el resumen del orquestador".
- E11. `review.md:42`: "Launch all three". Son tres pasadas, no un rompedor.

## JUICIO (analista funcional), pieza por pieza
Formato: actor / disparador / resultado observable / caso que falla hoy.

**A. Ciclo de a uno**
- Actor: orquestador. Disparador: `apply` termina en un cambio de riesgo completo.
- Resultado: queda un registro con, por vuelta, constructor → rompedor → auditor en mensajes distintos. A la tercera vuelta sin OK, el rompedor va directo al auditor. El veredicto aparece byte a byte en la respuesta a la persona.
- Falla hoy: no hay rol de constructor ni de rompedor en los assets (E7). `review.md` lanza tres pasadas por lentes y no tiene tope de vueltas (E11).
- Verificable si "textual" se define como "diff vacío entre el informe del auditor y lo que se le mostró a la persona". Confianza: alta.

**B. Constructor de otra familia**
- Actor: orquestador. Disparador: lanza `apply`.
- Resultado: el registro de despacho anota proveedor y modelo de cada agente, y el del constructor no coincide con el del rompedor ni con el del auditor.
- Falla hoy: la tabla de modelos sólo tiene alias de Claude (opus, sonnet, haiku) y no hay campo de proveedor.
- Así no es verificable: "otra familia" no tiene dueño de dato. Reformulación: "el registro declara `provider` por agente y el instalador rechaza que coincidan". Confianza: media. No leí `model-routing`.

**C. Firma humana**
- Actor: la persona. Disparador: una tarea toca un dominio CRITICAL.
- Resultado: el constructor sigue sólo si existe un commit que agrega la fila de la tabla `## Aprobación`, firmado por la persona.
- Falla hoy: `role-implementer` se niega igual (E5).
- **No es verificable como está escrita.** "Commiteada por ella" no se distingue de un commit que hace un agente: en esta máquina `git user` es Juani, y Codex corre sin sandbox y commitea (14_forma_de_trabajo:121-126).
- Reformulación: firma criptográfica (`git verify-commit` con una clave que tenga passphrase o esté en un token, que el agente no pueda usar) sobre el commit que toca sólo la sección `## Aprobación`. Sin eso, la pieza reabre exactamente el canal que E5 cierra, y E5 tiene razón. Confianza: alta.

**D. TDD con evidencia pegada**
- Actor: constructor. Disparador: marca `[x]` una tarea.
- Resultado: debajo de la tarea en `tasks.md` hay rojo y verde con comando, `date +%T` y salida. El rojo es anterior al verde. Si el test nació en verde, hay un diff de mutante cuya falla cae en la aserción que el test nombra (L-034).
- Falla hoy: E1.
- "Cae en la aserción" es verificable si se pide pegar la línea de la falla (archivo:línea de la aserción). Si no, es juicio. Confianza: alta.

**E. Lecciones**
- Actor: auditor (propone) y orquestador (agrega al archivar).
- Resultado: el veredicto tiene una sección `LECCIONES` (o "ninguna"), y el archivo de lecciones sólo crece: `git log -p` sin líneas borradas.
- Falla hoy: E7. No existe el archivo, la sección ni la lectura previa.
- "Leen antes de empezar" no es observable. Reformulación: el brief nombra las L-NNN que aplican y el informe las cita. Confianza: alta.

**F. Verificación real con proveedor externo**
- Actor: la persona al teclado. Disparador: `archive` de un cambio que integra un proveedor externo.
- Resultado: existe una tabla de pasos con hora y efecto observado contra la instancia real, o deuda declarada con su alcance.
- Falla hoy: nada en los assets lo pide. `review.md:14` sólo habla de que el proveedor "devuelve basura".
- Así no es verificable: "proveedor externo" no se detecta. Reformulación: el proposal declara `proveedores-externos:` y `archive` rechaza si hay alguno sin tabla. Confianza: media.

**G. Ship espera el CI**
- Actor: orquestador. Disparador: después del push.
- Resultado: el informe trae el estado final del run de ese SHA (`gh run list --commit <sha>`), y no hay un segundo push encima mientras el primero no terminó.
- Falla hoy: E2. Confianza: alta.

**H. Auditoría de cimientos**
- Así no es verificable: "núcleo crítico" no tiene definición.
- Reformulación: el proposal marca `cimientos: si`, y entonces corre `active-stack roles` sobre el árbol completo (no sobre un diff) y quedan informes de cada rol en la carpeta de L.
- Falla hoy: el clasificador lee un diff (`claude/sdd-orchestrator.md:362`). No confirmé que acepte "todo el árbol"; no abrí `classify.go`. Confianza: baja.

**I. UI mirada en un navegador**
- Actor: orquestador. Disparador: se cierra un bloque con archivos de UI.
- Resultado: capturas por pantalla con SHA y viewport, y un shape aprobado con fecha anterior al primer commit de UI.
- Falla hoy: E7. Ningún asset menciona navegador ni capturas.
- "Mirar" no se puede observar; las capturas sí. Confianza: media.

**J. Conflictos del brief**
- Resultado: el brief tiene dos listas, "resuelvo y lo declaro" y "freno". Cada resolución propia aparece en `tasks.md` atribuida al constructor.
- Falla hoy: no hay plantilla de brief para el constructor en los assets.
- Así no es verificable: la frontera "orden o redacción" vs "permisos" es juicio, y el 2b R2B-3 la cruzó sin que nadie lo notara (auditor fila 5). Reformulación: "todo lo que toca una regla firmada frena". Confianza: media.

**K. Suites de cierre**
- Resultado: cada corrida de cierre anota el SHA, y ese SHA es igual al HEAD al marcar el cierre.
- Falla hoy: ship pide SHA sólo al publicar, no al cerrar el bloque. El 2b corrió la web sobre 9eb0976 con HEAD en 2b570b3 (auditor fila 6). Confianza: alta.

**L. Informes fuera del repo**
- Resultado: existe `.git/informes-<cambio>/<rol>/informe.md` por cada agente del registro de despacho.
- Falla hoy: E10.
- Condición previa: si el orquestador escribe el archivo, vuelve a ser transcripción. Hay que pedir que lo escriba el agente y comparar contra su respuesta. Confianza: media.

**M. Bloque liviano**
- Así no es verificable: "~17.300 palabras" no lo medí, y "modo compacto" no tiene criterio de terminado.
- Reformulación: "con todos los toggles, lo que se inyecta en la fase X pesa ≤ N tokens y contiene las secciones Y", medido con un test sobre `Compose`.
- Falla hoy: E8 y E9. Confianza: media.

**N. Clasificador sin usar**
- Actor: orquestador. Disparador: convoca revisores.
- Resultado: el registro de despacho trae la salida pegada de `active-stack roles`, o "no pudo correr".
- Falla hoy, y no por un defecto de criterio: el bloque del repo ya lo exige (E3), pero el bloque instalado no lo tiene (E4), y PropTech no lo muestra nunca. Primero es un problema de reinstalación, después de cumplimiento. Confianza: alta.

## Contradicciones
1. **C contra `role-implementer`** (E5 contra E6). No se resuelve agregando texto: hace falta firma criptográfica, o dejar que la persona escriba ese código. Confianza: alta.
2. **A y H contra el punto 2 del despacho** (E3). En A no hay conflicto real: el rompedor necesita el resultado del constructor, así que pasa la prueba de dependencia de `:319-321`. En H sí: los revisores de cimientos son independientes, y el bloque exige mandarlos juntos. La memoria del usuario dice "de a uno, pocos agentes". Confianza: alta.
3. **A contra `review.md`**: un rompedor con tope de 3 vueltas contra tres pasadas paralelas por lentes, sin vueltas (E11). Hay que elegir uno o definir cómo conviven.
4. **A contra `roles.md:58`**: en la fila de apply no se convoca a nadie. El ciclo del método mete al rompedor *entre* vueltas de apply, cuando el diff todavía se mueve.
5. **14_forma_de_trabajo:156 (push a main sin preguntar, D-049) contra `ship.md`** ("Ask every time", "including the tenth time"). Si se importa el método, esa autorización de PropTech choca de frente con ship.
6. **L-024 contra `senior.md`**: L-024 cierra cada vuelta con revisores de rol; el senior se convoca una vez por fase. Un revisor de rol por vuelta multiplica agentes, contra el "pocos agentes" de la memoria.

## SUPUESTO
- No sé de dónde sale `role-implementer.md`: el texto de E5 no está en active-stack-V2 (rg: 0). Supuse que es de otro instalador o que se editó a mano.
- No corrí `active-stack roles` ni ningún test, ni conté palabras del bloque (M queda sin verificar).
- Leí `review.md`, `ship.md`, `roles.md` y `senior.md` por búsqueda, no enteros. Puede haber reglas que no capturé con los patrones.
- Las conclusiones sobre B (tabla de modelos) salen del bloque cargado en mi contexto, no del fragmento `model-routing` del repo.

## JUICIO (orquestador del informe): roles no convocados
- No aplica: este informe es sólo del analista. Confianza: alta.
