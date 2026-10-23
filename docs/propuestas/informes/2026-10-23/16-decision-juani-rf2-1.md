# Decisión de Juani sobre RF2-1 (23/10/2026)

Pregunta del senior (informe 15), presentada por el orquestador con tres opciones. Juani eligió **A**:

- La vuelta 3 corrige la parte de RF2-1 que es defecto del constructor: los mutantes que §3.3 y §3.4 nombran (lista cerrada) y las fallas RF1 ya cerradas tienen que dar rojo aunque cambie la redacción, con tests que busquen el VOCABULARIO de cada regla y no la frase exacta del mutante (por ejemplo, prohibir `ready (to|for) archive` en el bloque compuesto). También RF2-2 (degradación sin `review`), que conformidad marcó CRÍTICO.
- Las paráfrasis abiertas que no nombra el propose (N2b, N3, N4, N5, N8 del informe 14) pasan a DEUDA: las cubren los goldens revisados por diff y la lectura del rompedor, como ya decía §3.3 ("los goldens no prueban conducta real de agentes") y §3.4 ("pruebas del contrato textual").
- La regla nueva de lectura semántica en cada verify (opción B) NO se adopta ahora; candidata para el propose de la fase 2, con su propia firma (L-029).
