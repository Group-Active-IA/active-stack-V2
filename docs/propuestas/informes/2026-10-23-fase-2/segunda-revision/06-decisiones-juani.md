# Decisiones de Juani sobre el propose de la fase 2 (23/10/2026)

Juani aceptó las cinco recomendaciones del auditor (`05-auditor.md`):
1. **J1, firma de commits:** para firmar este propose vale la aprobación escrita (como en la fase 1). La firma SSH con lista de firmantes confiables es el objetivo del primer change de la fase 2.
2. **Originales (AR-1):** viven en el repo, con una excepción acotada en `.gitignore` para `openspec/changes/*/informes/`.
3. **Senior sin auditor (AN-2):** el senior informa y decide Juani (lo que ya dice `senior-without-review.md`).
4. **Cierre de los bloqueantes:** el constructor corrige y cada rol confirma sólo el suyo, sobre su delta; no hay otra ronda completa.
5. **Primera revisión (01 a 05 de la carpeta principal):** se conserva como historia, rotulada "custodia no verificable"; no cuenta como revisión.

## J5 (23/10/2026)

Juani eligió **A** (recomendada): la lectura semántica entra como regla de verify, con un lector separado bajo `review`, antes del senior y del auditor, dentro del mismo tope de tres vueltas. Es una regla nueva: se construye en el paso 5 con esta aprobación escrita como respaldo.

## J1, revisada (23/10/2026)

Juani decide **no configurar firma criptográfica de commits**. El paso 1 (`aprobacion-firmada`) se rediseña sin GPG ni SSH: la aprobación sigue siendo la escrita (tabla `## Aprobación` en commit propio, con alcance explícito), y lo que el paso puede aportar es verificar esa estructura (tabla presente, persona, fecha, alcance, commit separado del código) sin afirmar autenticidad criptográfica. Queda declarado como límite: un agente con acceso al repo puede imitar la firma escrita; el control es de proceso, no criptográfico.
