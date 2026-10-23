# Senior, vuelta 2 (Claude opus, textual del registro)

Mi lectura es que RF2-1 mezcla dos cosas distintas. Una es un defecto del constructor y se corrige en la vuelta 3. La otra es un límite que el propose firmado ya había aceptado. Tratarlo entero como "cualquier frase con el sentido invertido tiene que dar rojo" va contra §3.3, que está firmado.

**Fuentes que usé:** el propose (§3.1 a §3.4, §4, §5 y la Aprobación), la exploración §8, los informes 11, 13 y 14 (de 12-respuesta sólo un grep), y L-019, L-034 y L-035 de PropTech. No abrí `14_forma_de_trabajo.md`. En Engram, `mem_search` con todos los proyectos no devolvió nada, y no existe el proyecto `active-stack-v2`: eso quiere decir que nadie guardó nada, no que nadie lo decidió. No llamé a `mem_context`.

**EVIDENCIA** (comandos corridos hoy, 23/10)
- `git log -1 --format='%h %aI %s' 4748680` → `2026-10-23T09:52:20-03:00`, la aprobación de Juani.
- §3.3 nombra una lista cerrada de seis mutantes: quitar wiring o comentarios, aprobar con CRITICAL, anunciar archive antes del auditor, correr el auditor antes que vuelva el senior, reiniciar el contador y permitir una cuarta vuelta. También dice: "Los goldens se revisan por diff… no prueban conducta real de agentes". §3.4 agrega: "Son pruebas del contrato textual".
- `rg 'ready.{0,10}archive'` en los tests sólo encuentra `assets/commands_test.go:123`, que hace `strings.Contains(body, "ready for archive")`. Por eso N6 ("Ready to archive.") pasa toda la suite.
- `eza openspec/changes/archive` → vacío. El paso 4 (`verify-rompedor-auditor`) sigue abierto.

**RF2-1, BLOQUEANTE** (choca con §3.3 y §3.4, firmados el 23/10, `4748680`)
- **La parte que es defecto (confianza alta).** N6, N7, N1b y N9 a N11 son mutantes que §3.3/§3.4 nombran o fallas RF1 ya cerradas. Siguen vivos porque cada test busca la frase exacta del mutante y no el vocabulario de la regla. Es L-035: un test que existe no es un test que vigila. Además, el auditor de la vuelta 1 ya había avisado en el informe 11, línea 6, que un texto contradictorio que se agrega no rompe los tests. Esto se corrige en la vuelta 3, por ejemplo prohibiendo `ready (to|for) archive` en el bloque compuesto.
- **La parte que es límite.** Cerrar todas las paráfrasis posibles (N3, N4, N5, N8, N2b) no se puede con tests de texto. Para eso §3.3 ya puso los goldens revisados por diff y la lectura del rompedor. Si la vuelta 3 sale a perseguir las 13 frases, la que venga número 14 vuelve a pasar, y la vuelta 3 es la última.
- **Ojo con N1 (confianza media, no lo verifiqué contra GitHub):** un job cancelado y relanzado con éxito sobre el mismo SHA probablemente cuenta como verde según la definición de §3.4. Si es así, N1 no invierte el sentido y es ruido.
- **Opciones para Juani:**
  - (a) El auditor sólo bloquea la parte que es defecto (lista cerrada más el vocabulario de cada regla). Las paráfrasis abiertas van a deuda y quedan cubiertas por goldens y lectura humana.
  - (b) Se amplía §3.3 para exigir además una lectura semántica por agente o persona en verify. Es una regla nueva, así que según L-029 vuelve a firma.
  - (c) Se mantiene RF2-1 entero como bloqueante y se acepta que la vuelta 3 puede cerrar con eso abierto.

**RF2-2, NOTA** (real, no es ruido). La regla firmada es la fila `review` de §3.2. El hallazgo ya bloquea por su cuenta: conformidad lo marcó CRÍTICO (CONF-2), y §3.1 dice que un CRITICAL sin resolver impide APROBADO. No tengo historia que sumarle. Como el paso 4 no está archivado, corregirlo en este verify no reabre nada.

**RF2-3, NOTA.** Sin auditor, las dos reglas frenan el cierre, así que no se aprueba nada en falso. Va a deuda del paso 5 (`senior-solo-en-verify`).

**SUPUESTOS**
- Que "Mutantes de sentido" no aparece en otra parte del propose: busqué con rg y no lo leí completo.
- Que el rompedor tiene razón en qué tests corrió: no reproduje ninguno de sus mutantes.

Archivos relevantes:
- `/home/juanisarmiento/Proyectos/active-stack-V2/docs/propuestas/2026-10-23-metodologia-proptech-en-opsx.md`
- `/home/juanisarmiento/Proyectos/active-stack-V2/assets/commands_test.go`
- `/home/juanisarmiento/Proyectos/active-stack-V2/docs/propuestas/informes/2026-10-23/14-rompedor-vuelta-2.md`
