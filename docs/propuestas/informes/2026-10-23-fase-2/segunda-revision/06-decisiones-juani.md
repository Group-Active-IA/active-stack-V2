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

## N1/N1b, política de reintentos de CI (23/10/2026)

Juani eligió **A**: un job cancelado (o fallido por la plataforma) que se reintenta sobre **el mismo SHA** y termina en `success` cuenta como verde para ese job, siempre que el último intento de cada job requerido de ese SHA esté en `success`. Un reintento sobre otro SHA no cuenta, y un cancelado sin reintento exitoso sigue sin ser verde. Caso real que lo motiva: PropTech, 22/10/2026, job `web` cancelado por la instalación de dependencias y reintentado en verde sobre el mismo commit.
