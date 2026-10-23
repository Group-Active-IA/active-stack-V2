# Decisión de Juani sobre RF3-1 (23/10/2026)

El bucle de verify de la fase 1 se cerró en la vuelta 3 sin aprobación por RF3-1 (informe 21). Juani eligió la opción **B** del auditor:
corrección de una línea fuera del bucle en `internal/contracttest/text.go` (`ready (to|for) archive([^?]|$)`), con el FAIL de N6 sin punto y de "Ready for archive after conformity." en el bloque compuesto y en `opsx/verify.md`, con comando, salida y hora; y su OK escrito antes de archivar `verify-rompedor-auditor`.

**OK escrito de Juani (23/10/2026):** vista la corrección `ce7899e` y su verificación independiente (los tres mutantes de "ready to/for archive" dan FAIL en el bloque compuesto y en `opsx/verify.md`, sin falso positivo para "ready to archive?"; suite completa 26 ok, 0 fallas, 10:56), Juani aprueba archivar la fase 1. `instalacion-sana` NO se archiva: queda abierta con la condición CONF-1 (§3.6 puntos 1 a 4 y su OK para reinstalar).
