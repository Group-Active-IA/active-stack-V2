# Auditor del paso 0, vuelta 1 (Claude opus, textual del registro)

**Auditor de verify, paso 0 `contratos-fase2-deuda`, vuelta 1. Fecha: 2026-10-23 11:54 -03.**

EVIDENCIA: con `git diff --stat 7f1ef88..a6cfc2d -- . ':!docs/propuestas/informes'` salen 11 archivos: 9 de código, test o golden y 2 de `informes/` del change. Corrí `fdfind -I` y `git check-ignore`: `tasks.md` existe pero está ignorado (`.gitignore:40`), así que el rojo/verde versionado vive sólo en `apply-r1-constructor.md`.

**01-constructor (Codex)**
- C1 (evidencia): limpio. Cada rojo y verde trae comando, hora y SHA. La suite final sobre `46daf3b` está en `tasks.md:156`.
- C2 (suposiciones): hallazgo. "No acepta que el fallback… o una frase correcta escondan la contradicción" y "Cubre N2b, N3…" van como hechos, sin categoría, y el M1 del rompedor (N3) sobrevive. Además su mutante (`hits = nil`) prueba que se llama a la función, no que el test reconozca el vocabulario. La sección SUPUESTO sí tiene contenido real.
- C3 (cobertura): limpio. Lista lo que leyó y cubre los 9 archivos que tocó.
- C4 (acuerdo): no aplica a un autor. Sí cuenta que constructor y conformidad son la misma familia (Codex), algo que él mismo declara.
- C5 (lo que nadie miró): nada propio.

**02-conformidad (Codex)**
- C1: la suite cumple (`go test ./...`, 11:46:14 a 11:47:12, exit 0), aunque no escribe el SHA. En cambio "revisé `7f1ef88..46daf3b`" está archivado como EVIDENCIA sin comando ni salida.
- C2: no tiene sección SUPUESTO, su JUICIO no trae confianza y CONF-1 a 3 van como "ADVERTENCIA", que no es ninguna de las tres categorías.
- C3: hallazgo. No dice qué archivos abrió. Cita líneas de 2 archivos y dictamina "implementados" sobre el paso entero.
- C4: es el mismo proveedor que el constructor, así que no es un par independiente.
- C5: no dice nada de `compose.go`, `knowledge_base_toggle_test.go` ni los goldens.

**03-rompedor (Claude)**
- C1: los mutantes traen comando, hora y salida. Dos afirmaciones no: "compuse claude… (11:48)" en RP0-1 y "`debtViolations` da 0" en RP0-3 van sin comando. Las líneas que cita, `review.md:104/106` y `senior-without-auditor.md:3/7`, las verifiqué y existen.
- C2: limpio. Declara que no leyó al constructor y que sólo corrió `internal/harness/config`.
- C3: dice "No abrí ship", pero RP0-3 afirma "0 en `ship.md`". Le corrió el detector sin leerlo. Es menor.
- C4: hallazgo. Escribe "CONF-1 confirmado", "CONF-3 confirmado" y "consistente con CONF-2": leyó 02 antes de empezar. La convergencia 02/03 es una sola opinión confirmada. Lo que sí suma es su camino propio: 3 agentes y 6 mutantes.
- C5: nadie leyó `ship.md` y no hay golden de codex. El rompedor declara las dos cosas; la conformidad no.

**Choque de nombres.** El CONF-1 del propose (instalacion-sana, "fuera de pago") no es el CONF-1 de 02, y el constructor usa el del propose. Abajo, CONF-n significa el de 02.

**Clasificación**
- CONF-1 = RP0-1: BLOQUEANTE. Es la degradación sin auditor que el paso 0 tenía que pagar.
- CONF-2 + RP0-2: BLOQUEANTE. M1 y M2 sobreviven en verde, y M3, M5 y M6 sólo los atrapa un golden, que se regenera y acepta la inversión.
- CONF-3 = RP0-3: BLOQUEANTE (subo la NOTA del rompedor). Es la lección de la fase 1, cuesta una aserción en el mismo archivo, y hoy el corpus es atribuible sólo por un hecho que ningún test exige.
- RP0-4: DEUDA. El sentido es compatible; va al paso S.
- N1: fuera de la vuelta 2. Es regla nueva y su corpus es del paso S/ship (propose L104 y L116). La decisión de Juani del 23/10 sirve como su aprobación escrita, y ahí tiene que ir el control positivo "reintento en `success` sobre el mismo SHA no se rechaza".

**Lista cerrada para la vuelta 2.** Nada fuera de estos tres puntos.
1. Con review y senior, sin auditor, en claude, codex y opencode: ninguna cláusula de la composición puede exigir veredicto o APROBADO del auditor para cerrar o archivar. Con auditor encendido se sigue exigiendo. Rojo visto sobre `8469d94`.
2. Aserciones que no sean golden, por clase de vocabulario: (a) el senior o su informe cierra, aprueba o archiva; (b) esperar o mantener abierto por un auditor ausente; (c) el rompedor autoriza el cierre; (d) completar o inventar las listas del auditor. Cada mutación parte del original limpio y el mensaje nombra la clase. Terminado: M1 a M6 de 03 más al menos 4 redacciones nuevas del rompedor de la vuelta 2, que el constructor no ve, mueren todas por aserción.
3. `TestDebtCorpusAdditiveCopies` exige `debtViolations(original)==0` y que el hit nuevo venga del agregado. Terminado: si se agrega "The review does not block." a `review.md`, el test falla por "original no limpio" en vez de pasar 28/28.

DECISIÓN: NO APROBADO (vuelta 1 de 3)
