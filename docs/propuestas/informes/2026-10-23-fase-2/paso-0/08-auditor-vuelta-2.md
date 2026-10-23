# Auditor del paso 0, vuelta 2, cierre en modo rápido (Claude opus, textual del registro)

**Auditor de verify, paso 0 `contratos-fase2-deuda`, vuelta 2 (cierre en modo rápido). Fecha: 2026-10-23.**

EVIDENCIA: `git diff --stat 32612f9..HEAD -- internal` no tiene salida. El código no cambió después del constructor, así que la suite de 06 (sobre `b42ea48`) y los mutantes de 07 (sobre `1e67e89`) corrieron sobre el árbol actual y no están vencidos. Por fuera de los informes, `git diff --stat c9dd12b..ebffa80` toca 5 archivos: `review.md`, `debt_contract_test.go`, `review_toggle_test.go` y dos goldens.

**05-constructor (Codex)**
- C1/C2: hallazgo. En `05-constructor-v2.md:20` va como EVIDENCIA, sin comando, que `authorityViolations` "no busca igualdad con los casos de la tabla", y 07 lo refuta: 0 de 8 redacciones nuevas mueren. En cambio su SUPUESTO (L236) es honesto: dice que no verifica las redacciones ciegas. El que exagera es su resumen: "Corregidos los tres puntos" (`05-constructor-v2-respuesta.md:1`).
- C3: limpio. L8 lista lo que leyó y cubre los 5 archivos. C4 no aplica porque es un autor. C5: nada propio.

**06-conformidad (Codex)**
- C1: limpio. La suite trae comando, SHA y hora. El punto 3 queda como JUICIO y declara "No ejecuté esa mutación", así que no lo sube de categoría.
- C2: queda una "ADVERTENCIA" (L6) que no es ninguna de las tres categorías. Es menor. Ahora sí tiene SUPUESTO.
- C3/C5: limpio. Abrió los 5 archivos, incluidos los dos goldens (L17-18).

**07-rompedor (Claude)**
- C1/C2: limpio. Cada mutante trae texto, hora y salida, y el SUPUESTO declara que excluye golden y ancla (L38).
- C3: declara lo que leyó (L40). No abrió `ship.md` ni las copias de comandos de codex y opencode, que están fuera de la lista cerrada.
- C4: el acuerdo con 06 sobre el punto 2 es real. Llegaron por caminos distintos: 06 por procedimiento (V1-V4 son del constructor) y 07 empíricamente (8 mutantes ciegos). Además 07 dice no haber leído 05 ni 06.
- C5: nadie leyó `ship.md`, que es deuda heredada de la vuelta 1.

**Clasificación.** Punto 1: cumplido, con evidencia en 07 L23-26. Punto 3: cumplido, con evidencia en 07 L29. Punto 2: NO cumplido. RP0b-2: DEUDA.

**DEUDA del paso 0, a cobrar en el paso 5 (lectura semántica):**
- D1 (grave): el guard de autoridad es una lista cerrada de verbos en inglés. Cualquier paráfrasis que le dé el cierre al senior o al rompedor, o que haga esperar a un auditor apagado, pasa la suite en verde. Ejemplos: "cierra", "settles", "weighs in", o "cannot" en una subordinada, que apaga la detección de toda la cláusula. Sólo la atrapa un golden, que se regenera y la acepta. Riesgo: una edición futura de `review.md` o `senior-without-auditor.md` reintroduce la inversión de autoridad sin que nada falle. El paso 5 tiene que usar N1a a N8d (07 L12-15) como corpus del lector semántico, y su terminado es que mueran los 8.
- D2: el bloque sin auditor dice "A full round has … auditor reports" y "A missing required report stops closure" (claude L588/L592). Una lectura literal frena la ronda por un informe del auditor que no puede existir.

DECISIÓN: CERRADO CON DEUDA ACEPTADA (modo rápido)
