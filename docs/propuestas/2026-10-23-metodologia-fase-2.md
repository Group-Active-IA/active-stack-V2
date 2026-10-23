# Propuesta: metodología PropTech en OPSX, fase 2

Fecha: 2026-10-23 · Rama: `feat/metodologia-fase-2` · Estado: PROPUESTO, sin aprobación humana. Sólo EXPLORE/PROPOSE; no autoriza apply, instalación ni nuevas reglas operativas.

## Why

Fase 1 dejó aprobación escrita sin identidad criptográfica, originales mediados por el orquestador, residuos de bash al cambiar tiers y roles referidos sin instalación. Sus tests de texto dejaron pasar reglas opuestas escritas con otras palabras. La [exploración de fase 2](../exploraciones/2026-10-23-metodologia-fase-2.md) registra fuentes oficiales, contratos y límites; la decisión 16 exige aprobación propia para S. Este diseño cubre pasos 8–11 y S, más deuda ligada a esas piezas. No incluye compactación, UI, código ahora ni resolución de CONF-1.

## What Changes

### 8. Aprobación verificable y modo sin firma — ALTO

Proponer `active-stack approval --change <id> --commit <A>` como lectura: no genera firmas ni completa aprobación. `internal/approval` nuevo separaría resolución de Git, política de identidad y validación del manifiesto; `cmd/active-stack` expondría resultado verificable. Contrato de salida propuesto: VERIFIED, WRITTEN_UNVERIFIED, MISSING, INVALID, STALE, con códigos 0 sólo para VERIFIED, 2 para modo escrito identificado y 1 para faltante/error/invalidez; los consumidores nunca interpretan todo retorno no fatal como aprobación.

En modo firmado, A es un commit de aprobación humana; incluye elección, persona, fecha, change, alcance, SHA completo P del propose revisado y lista de archivos/hashes de propuesta e informes. P es ancestro de A y debe contener los documentos revisados. La tabla humana vive en el documento de aprobación contenido en A; si Juani completa la tabla de este propose, el delta permitido P→A es exclusivamente esa sección más el manifiesto y documento de decisión, sin editar alcance ni informes. El manifiesto hashea los documentos de P antes de completar esa tabla: se compara el delta permitido por estructura y no se exige igualdad imposible del hash completo tras la aprobación. Cualquier otro cambio normativo exige un nuevo P revisado. En ese modo se verifica la firma de A mediante Git y se compara la identidad criptográfica autorizada, no author/email. La política de llaves/principales y revocación viene de una raíz humana confiable fuera de archivos modificables por el change; su ubicación la resuelve un contrato de configuración explícito (TBD), sin escribirla en este trabajo. No ejecutar programas de verificación arbitrarios definidos por el repo: usar entorno/configuración de confianza controlados y Git invocado con argumentos, timeout y errores visibles. Ensayar repo con config Git maliciosa y demostrar que no ejecuta programas ni cambia política de confianza; configuración exacta queda como gate del change 8 (TBD).

El verificador exige tabla humana no vacía y coincidente con el manifiesto; comprueba contenido de P y ausencia de alteraciones normativas posteriores hasta el árbol evaluado. Cambio de alcance, regla o informes usados invalida aprobación; cambios posteriores de implementación dentro de alcance no cambian P, pero se revisan en verify. Si no puede distinguirlos, STALE y consulta humana. Se rechazan SHA abreviado ambiguo, commit ajeno al repo/change, firma desconocida/revocada, firma válida sin alcance, árbol documental sucio y rebase que invalide referencias. Revalidar en apply y antes de cierre; no inferir consentimiento del commit que agrega estos documentos.

Decisión J1 de Juani del 23/10/2026: para este propose vale aprobación escrita; SSH con lista de firmantes confiables es el objetivo del primer change de fase 2. La ubicación de la raíz de confianza y revocación sigue (TBD). Para este propose se mantiene el modo escrito elegido: sólo una decisión escrita completa y coincidente recibe WRITTEN_UNVERIFIED, nunca VERIFIED. Sin tabla/decisión es MISSING; inconsistente es INVALID; alcance obsoleto es STALE. Si Juani elige requerir firma, la ausencia de verificador, llave o confianza detiene apply ALTO. En modo escrito, cada alcance ALTO requiere su aprobación escrita explícita y se declara degradación; no se instala ni configura firma automáticamente. La tabla final queda vacía hasta su intervención.

**Escenario «tabla completa sin commit A» (AN-1):** entrada: tabla/decisión escrita por Juani completa, coincidente con change, alcance y P, árbol documental limpio y originales versionados en P, pero sin commit A. En modo escrito devuelve WRITTEN_UNVERIFIED, código 2; el consumidor habilita apply sólo para ese alcance explícitamente aprobado y registra la degradación. En modo firmado devuelve MISSING, código 1, y bloquea apply: no hay firma de A que verificar. La tabla vacía, ajena, inconsistente o de alcance obsoleto conserva MISSING/INVALID/STALE, código 1. No se exige A para el modo escrito. El constructor prepara un manifiesto determinista de P mediante tooling de sólo lectura (change, archivos, hashes y procedencia de originales); Juani escribe la decisión y el verificador contrasta ambos sin completar la tabla. La decisión escrita puede registrarse después de P en un commit documental sin firma, sin convertirlo en A firmado: debe referenciar el SHA P ya existente; no puede hacerse que P contenga su propio SHA. El delta documental permitido se controla igual que en modo firmado.

### 9. Originales sin transcripción del orquestador — MEDIO, persistencia ALTO

Proponer `internal/reports` y `reports capture`: escritor determinista separado del modelo. Adaptadores Claude/Codex/OpenCode leen eventos o registros originales según la matriz de contratos de la exploración. Para Claude, cubrir también Handback; para Codex, JSON de salida propio y campos opcionales; para OpenCode, resolver sesión hija y mensaje completo mediante API, sin considerar idle sinónimo de informe. Fijar versión objetivo y fixture de cada cliente antes de habilitar soporte. Contrato no confirmado → UNSUPPORTED visible, nunca captura supuesta.

Antes del despacho registrar repo/worktree canónico, change, fase, vuelta, rol, ID de despacho, agente/sesión/turno cuando estén disponibles y SHA/diff revisado. No derivar identidad del nombre escrito por el modelo. Una asociación ambigua queda PENDING. Proveedor/modelo sólo EVIDENCIA si constan en el registro; de otro modo SUPOSICIÓN. Capturar bruto limitado al informe/evento necesario y texto íntegro como derivados separados; no volcar todo transcript personal al repo. Mantener procedencia, versión, hash, tamaño, estado final y mediación automática en índice `LEEME.md` legible y manifiesto estructurado versionado.

Destino versionado elegido por Juani (AR-1): `openspec/changes/<change>/informes/<fase>-r<n>-<rol>-<dispatch-id>/`, relativo a raíz validada y adapter de registros, sin rutas globales hardcodeadas. Informes preexistentes nunca se sobrescriben: creación exclusiva/atómica; mismo evento/mismo hash es idempotente; contenido diferente con misma identidad es conflicto visible y se conserva separado para auditor. Rechazar traversal, symlink fuera de raíz, worktree equivocado, reporte truncado, mezcla de roles y retorno parcial. Ensayar cortes después del original y antes/después de manifiesto, índice y referencias; al reiniciar nunca COMPLETO sin todos los artefactos, ni pérdida de archivos personales. Escritura fallida deja estado recuperable y cierre detenido; índice sólo marca completo después de persistir originales. Reintentos no incrementan vuelta OPSX; límites de tamaño/timeout y retención exactos (TBD) se muestran como error, jamás truncamiento silencioso.

Para incorporar ese destino, el change de captura debe sustituir la exclusión `openspec/` por esta excepción acotada (diseño; este trabajo no edita `.gitignore`):

```gitignore
/openspec/*
!/openspec/changes/
/openspec/changes/*
!/openspec/changes/*/
/openspec/changes/*/*
!/openspec/changes/*/informes/
```

El capturador valida igualdad byte a byte/hash y procedencia contra el registro original identificado antes de marcar COMPLETO; `internal/approval` vuelve a contrastar hashes, identidad de despacho e inclusión de todos los originales en P antes de aceptar el manifiesto. Ambigüedad/conflicto deja PENDING y bloquea aprobación/cierre. El auditor inspecciona esos originales y sus límites de procedencia, sin inferir autenticidad del hash. Gate de aceptación: `git check-ignore openspec/changes/metodologia-fase-2/informes/example.md` debe devolver exit 1; verificar también que propuesta/tareas de OpenSpec siguen ignoradas. Hoy el repo aún ignora el destino: no se acredita que la excepción esté instalada. Se descartan exportar copias a docs y aprobar sólo hashes externos: duplican fuentes o dejan P sin originales. Revertir la elección exige migrar originales/referencias/manifiestos y reaprobar alcances cuyo material cambie.

Revisores sin Write; hook/capturador tiene permiso acotado de persistencia. El comando instalado usa ejecutable resuelto por contrato del instalador y quoting específico por plataforma probado; entradas de evento/change/rol viajan como datos JSON/argumentos, nunca interpolación shell. Capturador ausente o fallido deja informe pendiente y cierre detenido. Hooks personales se conservan, se preservan su orden y semántica al agregar entradas propias identificables con backup previo y merge por markers donde corresponda; si el formato no permite markers usar merge estructurado y manifiesto, sin insertar claves inválidas. Recuperación manual candidata sólo lee original verificable, no acepta resumen como sustituto. Si falta informe obligatorio se detiene el cierre. Este mecanismo reduce mediación textual; no demuestra independencia ni inviolabilidad frente al usuario del sistema.

**Ciclo de vida de hooks propios (OPS-1):** el instalador registra ID propio, ubicación resuelta por adapter, versión y hash de cada entrada gestionada, con snapshot de la configuración. Update compara contenido actual con ownership y hash: propio intacto → backup y reemplazo de esa entrada; idéntico → no-op; ajeno, desconocido o modificado → conservar y reportar conflicto, sin duplicar ni apropiárselo por nombre. Uninstall elimina sólo entradas propias intactas y sus referencias, conserva hooks personales y su orden, y retira el contenedor sólo si fue creado por Active Stack y quedó vacío. Entradas modificadas se preservan con conflicto visible. Se valida concurrencia antes de mutar; fallo parcial recupera localmente sólo escrituras propias, preserva ediciones concurrentes y declara recuperación incompleta si corresponde. El manifiesto solo no autoriza borrar.

### 10. OpenCode: reemplazo de bash — ALTO

Reimplementar intención de `ddeb94f`, sin su reemplazo completo de permission ni su trailer. Cambio previsto en `internal/harness/config/permissions/overlays.go`: `permission.bash` envuelto con sentinel existente `__replace__`, preservando orden de `orderedBashBlock`. `permissions/install.go` debe seguir usando snapshot previo; probar integración con merge final, no sólo overlay. No modificar internals de backup/filemerge/pipeline salvo propuesta y review ALTO adicionales. Validar configuración existente antes de merge: base inválida/no interpretable implica error visible y cero escritura; no aceptar el fallback actual de MergeJSONObjects a objeto vacío. La validación previa puede vivir en installer; reparar configuración exige alcance aparte.

Alcance destructivo deliberado: se reemplaza el objeto bash, incluidas reglas personales dentro de él, únicamente tras mostrar diff y aprobar esa frontera. Conservar todas las claves hermanas personales fuera de las claves escalares que el tier ya administra, además de settings ajenos a permission. Mantener los escalares explícitos actuales por tier, sin reabrir el bug de external_directory. Snapshot identificable de contenido/permisos antes de escritura; fallo de backup impide mutación. Fallo posterior restaura ese snapshot o informa recuperación incompleta, nunca éxito ficticio. Reinstalación idempotente. Comprobar hash/identidad del archivo entre lectura, snapshot y escritura; cambio concurrente aborta sin pisar. Para restaurar, comparar con bytes escritos por esta operación: si alguien editó después, preservar ambas versiones y declarar conflicto de recuperación, sin restore ciego.

La documentación actual habla de shell: compatibilidad bash/versión objetivo es gate previo, no detalle cosmético. Si no se demuestra soporte de bash en versión fijada, este paso queda bloqueado para apply y requiere propuesta ALTO de migración; no escribir ambas claves por las dudas. `__replace__` se consume en el instalador y no llega a OpenCode. Verificar orden y efecto last-match-wins sobre bytes finales; port de tests de Juan Cruz ajustado al alcance reducido. Sustituir `TestOpencodeBashResidueSurvivesTierSwitch` por ausencia de residuo y preservación de hermanos, sin borrar evidencia del bug histórico.

### 11. Instalar definiciones de roles — ALTO

Harness nuevo candidato `role-definitions`, tipo `config`, separado del toggle `roles`: éste convoca, el harness entrega archivos. Primera entrega portable: `role-*.md` bundleados como recursos dentro de un directorio propio de Active Stack bajo `SkillsDir` del adapter. Proponer pequeño SKILL.md índice de esos recursos, sin afirmar que cada archivo es agente nativo. Nombre exacto/directorio, modos lite/full/custom y lista canónica quedan sujetos a J3/J4 (TBD). No extender HarnessType ni inventar destinos de comandos para Codex. Un capability del adapter resuelve ubicación y soporte; vacío significa no soportado, no fallback a un path inventado.

**Enlace de composición (AR-2):** el installer resuelve project/global mediante el capability del adapter y calcula hash/versión desde los assets instalados y validados contra su manifiesto. Dirección: installer → adapter/resolución y validación de recursos → render puro con contexto explícito de ruta/versión/hash → filemerge/inyección. `Compose(toggles, variantKey)` conserva su contrato puro actual; se propone una etapa pura de enlace posterior que recibe ese contexto, sin leer disco ni consultar adapters desde Compose. El chequeo instalado reconstruye Compose más el mismo enlace y compara contenido y recursos; no compara contra Compose sin resolver. El plan exige la precondición «orquestador con toggle roles → role-definitions instalado y validado» antes de publicar referencias, también en custom; no depende del orden incidental. Se descartan referencias relativas estables porque project/global y worktrees no comparten una base garantizada. Revertir requiere cambiar contrato de enlace, goldens, chequeo instalado y migrar referencias ya inyectadas con backup/markers; no basta cambiar el nombre del directorio.

El bloque `roles.md` referenciaría ruta explícita efectiva e índice con versión/hash; `active-stack roles` seguiría clasificando y expondría faltantes como estado distinto de lista vacía. Si falta definición necesaria, detener convocatoria/cierre y pedir resolución; un rol de nombre correcto sin contenido no cuenta. Fuente/licencia y contenido de los siete especialistas de fase 1 deben revisarse antes de bundlearlos (TBD); senior/auditor no se convierten por accidente en esos siete.

En install/update/uninstall aplicar también rechazo de traversal y symlinks fuera de destino, comparación de hashes concurrentes y ownership contra contenido actual; el manifiesto solo no autoriza sobrescribir. Ownership en manifiesto: ausente → crear; existente idéntico → no-op; propio sin cambios personales → snapshot y actualización; personal/desconocido/modificado → preservar y reportar conflicto con alternativa de ruta propia. Nunca borrar por glob `role-*`, ni reemplazar por nombre. Referencias inyectadas mediante filemerge con markers idempotentes y backup. Probar uninstall/rollback: sólo retirar archivos propios sin modificaciones, conservar personales y reportar conflictos. El dispatch actual de HarnessConfig instala el orquestador para todo config: separar explícitamente por ID y probar que role-definitions no vuelve a inyectar el orquestador. Instalar assets antes de publicar referencias; dependencia/orden deben quedar en plan. Hooks y roles recuperan sus escrituras parciales dentro del paso fallido antes de retornar error, porque rollback de pipeline sólo cubre pasos exitosos; registrar recuperación incompleta si no pueden. Integrar planner/pipeline usando contratos existentes, con revisión ALTO de cualquier cambio de rollback. Opción de exportación nativa por agente se difiere.

### S. Lectura semántica por verify — regla nueva, aprobación propia

Proponer contrato dentro de `review`, sin toggle nuevo: con review activo, conformidad → rompedor (tres lentes) → lector semántico → senior activo → auditor (si auditor está activo). Con senior encendido y auditor apagado, el senior informa y decide Juani por escrito sobre el alcance degradado, tanto con review como sin review; no se omite al senior ni se espera un auditor inexistente. El lector es despacho separado del constructor y rompedor; recibe bloque completo efectivamente compuesto para variante/toggles y comandos aplicables, alcance/diff y decisiones fechadas. No recibe conclusiones de otros revisores antes de emitir el suyo; comparte filesystem, por lo que aislamiento es declarado y SUPUESTO salvo control demostrable. Auditor y senior sí reciben ese original según su función. No se duplica el contador: un informe semántico por vuelta, hasta tres.

Busca contradicciones entre cláusulas y con decisiones: quién aprueba, qué frena, cuándo vuelve apply, qué pasa sin toggle y qué evidencia permite archive/ship. Entrega fuentes leídas con hashes, pares de citas/localizaciones, escenario que cambia por obedecer cada una, JUICIO/confianza, límites y SUPOSICIONES. Cero hallazgos también exige cobertura y fuentes; no equivale a prueba de corrección. No reescribe el bloque ni completa aprobación; propone hallazgos al auditor, que decide. Falta de informe obligatorio detiene cierre, sin presumir APROBADO.

Con review apagado se anuncia verify ordinario sin lector ni bucle; si auditor apagado y senior encendido, el senior informa y decide Juani por escrito sobre alcance degradado, sin esperar un veredicto inexistente; con senior apagado decide Juani sin ese informe. Sin modelo disponible o captura completa, se informa faltante; Juani puede habilitar recurso o decidir prospectivamente un alcance sin S mediante decisión escrita que enumere obligaciones exceptuadas. Esa degradación no valida captura incompleta ni sustituye un informe obligatorio perdido: sin recuperación íntegra se detiene el cierre ordinario. No se inventa resultado ni cuarta vuelta. El semantic reader no sustituye las tres lentes del rompedor ni elimina tests, goldens o evidencia conductual. Modelo/proveedor recomendado y presupuesto son (TBD), no se impone acceso a otro proveedor. No se adopta S hasta respuesta propia J5.

## Cambios por fragmento o paquete

| Pieza | Cambio futuro y consumidor |
|---|---|
| `cmd/active-stack`, nuevos `internal/approval` / `internal/reports` | CLI lectura de aprobación, captura y estados estructurados; orquestador/health checks consumen |
| `internal/agents` | Capacidades/rutas de roles y registros por proveedor, project/global y worktrees; sin hardcodear homes |
| `permissions/overlays.go`, `install.go` y tests | Bash replace, preservación, snapshot y orden en settings reales |
| `internal/harness/config`, `internal/catalog/harnesses.yaml` | Harness config de roles, embed/manifiesto e instalación; `catalog.Load()` obligatorio si se edita YAML |
| `governance.md` | A/P/alcance, modo firmado o escrito y nueva regla con aprobación propia |
| `roles.md`, variantes `*/sdd-orchestrator.md` | Referencias reales a definiciones; claves Engram aisladas y entrega completa |
| `review.md`, `auditor.md`, `senior.md`, `senior-without-review.md`, `verify-without-review.md` | S, autoridad consistente y degradaciones; ausencia de auditor no genera espera imposible |
| `engram-protocol.md` | En verify clave `opsx/{change}/verify/round-{n}/{role}`, sin búsquedas de conclusiones hermanas; fuera de verify conservar protocolo explícito |
| `knowledge-base.md`, `ship.md`, comandos Claude verify/ship | RF2-4, originales, S y paráfrasis CI en consumidores reales; no afirmar comandos inexistentes en otros agentes |
| `internal/contracttest`, tests de config/assets/agents/roles y goldens | Corpus aditivo, integración y revisión de composiciones; goldens actualizados sólo tras rojo explicado |
| `backup`, `filemerge`, `pipeline` | Reutilización como dependencia ALTO; no autoriza rediseñar sus garantías |

## Tests que deben verse fallar primero

Son diseño de pruebas, NO resultados ejecutados. Por tarea futura registrar comando exacto, salida, hora y SHA en tasks.md, rojo antes de código y verde después. Si nace verde, mutante que falle en su aserción y se revierta. Error de compilación o falta de entorno no acredita rojo funcional. Tests con repos/home temporales, jamás destinos instalados del operador; no build sin pedido.

| Paso / comando candidato | Fallo discriminante requerido antes de implementar |
|---|---|
| 8: `go test ./internal/approval ./cmd/active-stack -run Approval -count=1 -v` | Firma válida de identidad no autorizada, firma sin scope, P obsoleto o manifiesto alterado rechazados; firma GPG/SSH confiable del P correcto aceptada. Tabla completa sin commit A: modo escrito WRITTEN_UNVERIFIED/código 2 permite sólo apply del alcance escrito; modo firmado MISSING/código 1 bloquea apply. Constructor prepara manifiesto determinista de P y verificador lo contrasta con decisión humana y originales versionados. Sin firma y con decisión escrita completa/coincidente reporta WRITTEN_UNVERIFIED; vacía/ausente MISSING, inconsistente INVALID, obsoleta STALE. Probar códigos y consumo en ambos modos; nunca VERIFIED sin firma confiable. Fixtures efímeros, no llaves personales. Iniciar con seam probado antes de crear paquete completo. |
| 9: `go test ./internal/reports -run Capture -count=1 -v` | Stop con despedida y Handback real guarda informe, no despedida; Codex nulos y transcript incompleto quedan pendientes; OpenCode idle padre no cierra hija. Duplicado idempotente, conflicto preservado, eventos fuera de orden, concurrencia y reinicio; igualdad de bytes/hash del original y derivado documentado. |
| 9 hooks: `go test ./internal/reports ./internal/install ./internal/uninstall -run Hook -count=1 -v` | Instalar → actualizar → desinstalar mediante planes reales: entrada propia intacta se actualiza y elimina con backup; personal/homónima o propia modificada se preserva con conflicto; orden personal intacto, reinstalación sin duplicados, fallo parcial recuperado y edición concurrente conservada. |
| 9 seguridad/persistencia: `go test ./internal/reports -run 'Worktree|Atomic|Recovery' -count=1 -v` | Traversal/symlink/worktree ajeno rechazados; disco lleno/permiso denegado no dejan completo; informe faltante bloquea cierre y no reinicia contador. Orquestador no puede sustituir origen con resumen. |
| 10: `go test ./internal/harness/config/permissions -run Opencode -count=1 -v` | Matriz ordenada 3×3 de tiers (incluye reinstalar mismo), sembrar overrides bash ajenos y claves hermanas personales; exigir sólo bash nuevo, hermanos intactos y orden final deny-last. Viejo test de residuo debe dejar de ser criterio aceptable. Base JSON inválida con contenido personal conserva bytes exactos y devuelve error, sin fallback a vacío. Snapshot fallido deja bytes intactos y fallo inducido restaura hash/permisos; parser cliente fijado reconoce clave. |
| 11: `go test ./internal/harness/config ./internal/agents/... ./internal/roles ./internal/catalog ./internal/install ./internal/uninstall ./internal/verify ./cmd/active-stack -count=1 -v` | Plan real en home temporal: instalar role-definitions, modificar un rol personal, inducir fallo posterior y desinstalar. Assets incluidos en snapshot/recuperación; probar recuperación local del paso fallido y rollback de pasos exitosos, preservando personales. Vaciar o quitar definición: health checks y CLI `roles` informan faltante distinto de clasificación vacía, sin éxito ficticio. Rol personal homónimo/modificado intacto, propio actualizado con backup, rollback no borra personal, project/global separados, adapter sin soporte falla visible, índice referencia archivos existentes. `catalog.Load()` valida harness real. |
| Deuda/S: `go test ./internal/harness/config ./assets -count=1 -v` | Composiciones de las ocho variantes × combinaciones review/auditor/senior/engram: senior encendido/auditor apagado exige informe del senior y decisión escrita de Juani con/sin review, sin espera de auditor; fase base no ordena save compartido; lector sólo cuando review. Las dos copias verify/ship reciben corpus. |

Corpus mínimo aditivo: conservar frase correcta y agregar N1/N1b (CI cancelado y rerun), N2b (persona desplaza auditor en vuelta 1), N3 (senior decide), N4 (append sin rewrite), N5 (agente llena aprobación), N8 (resumen sustituye informe) y RF3-2 (`treated as green`, `With zero runs, treat delivery as green`). Variar idioma, negación, puntuación y sujeto. Incluir positivos: citas históricas de errores, ejemplos negados, cancelación reemplazada por corrida exitosa sólo si la política CI la satisface realmente. No fijar regla nueva sobre reruns sin resolver política aplicable con Juani.

No basta `strings.Contains` ni ampliar regex hasta calcar corpus: tests deterministas delimitan regresiones conocidas, goldens muestran diff, y una evaluación del lector con casos ciegos exige localizar la contradicción y explicar sus dos conductas. Registrar falsos positivos/negativos y casos no cubiertos. Criterio candidato: todos los casos bloqueantes sembrados identificados con su escenario y controles positivos sin bloqueo injustificado; ejecuciones repetidas y variación de modelo declaradas. Separar corpus de desarrollo y evaluación ciega; congelar denominador, repeticiones y tolerancias en el change de S antes del ensayo. Tamaño/margen estadístico (TBD), no afirmar cobertura semántica universal. Ninguna deuda se cierra sólo porque S fue añadido.

## Deuda de fase 1 y aceptación

| Ítem | Incremento y evidencia futura para pagarla |
|---|---|
| CONF-3/RF2-3 | Primer paso: senior encendido y auditor apagado → el senior informa y decide Juani por escrito; comprobar el mismo resultado con/sin review |
| RF1-4 + RF2-4 | Autoridad reviewer/auditor/persona coherente en sección vieja y bucle; corregir mayúscula incidental KB con test de regla, no frase |
| RF1-6 | Paso 9: todas las instrucciones base/protocolo/fragmentos entregan clave por rol y vuelta en verify; ensayo de dos roles no comparte observación |
| RF3-2 | Paso S/ship: corpus en bloque y comando, más interpretación de cero runs y estados incompletos |
| Paráfrasis N1/N1b/N2b/N3/N4/N5/N8 | Paso S: corpus original y casos nuevos del revisor; auditor registra alcance residual |
| CONF-1 | Fuera de pago aquí; sigue abierto instalacion-sana y no autoriza reinstalación |

## Orden de construcción y governance

Cada fila es change candidato OPSX, ≤12 tareas; los de nivel ALTO requieren aprobación propia; ninguna empieza por este documento sin aprobación humana escrita del alcance, con verificación criptográfica sólo si Juani elige ese modo. BAJO conserva su autonomía dentro del alcance autorizado; las reglas nuevas requieren aprobación propia sin convertir todo su dominio en ALTO. Goldens se editan en secuencia.

| Orden | Change candidato | Pasos pequeños / checkpoint | Nivel y dependencia |
|---|---|---|---|
| 0 | `contratos-fase2-deuda` | Reproducir deuda → corregir autoridad/degradación/Mark → corpus conocido y revisar goldens | BAJO en texto; aprobación de cualquier regla nueva; base para S |
| 1 | `aprobacion-firmada` | Implementar objetivo J1 SSH/lista confiable → fixtures de identidad → verificar A/P/alcance → CLI sin escrituras → consumo y fallos | ALTO; arquitecto, funcional, QA, especialista seguridad/DevOps y auditor antes de aprobación humana |
| 2a | `captura-informes-nucleo` | Contrato de despacho → escritor/índice exclusivo → recuperación y aislamiento → RF1-6 | ALTO por persistencia/recuperación; 0 y contrato de aprobación resuelto |
| 2b | `captura-claude` | Fixtures Stop/Handback → adapter → ensayo aislado de hooks personales | MEDIO con checkpoint; escritura/backup ALTO; 2a |
| 2c | `captura-codex` | Fixtures opcionales → protocolo JSON → captura/reintento aislado | MEDIO, integración ALTO; 2a |
| 2d | `captura-opencode` | Fijar versión/sesión hija → API mensaje completo → fallos e idle falso | MEDIO, integración ALTO; 2a; no anunciar paridad antes de prueba |
| 3 | `opencode-permission-replace` | Compatibilidad bash → rojo de residuo → sentinel acotado → orden final → rollback aislado | ALTO; independiente de captura, aprobación específica de bash personal |
| 4 | `instalar-roles` | Resolver fuente/licencia → harness config/adapter → ownership → update/conflicto → uninstall/rollback aislado | ALTO; 2a y contratos de captura para aceptación; sin depender de todos los proveedores |
| 5 | `verify-lectura-semantica` | Aprobar J5 → corpus ciego → despacho/informe → secuencia y degradación → evaluación y deuda residual | Texto BAJO; regla requiere aprobación propia; integración MEDIO; 0 y captura del proveedor usado |
| 6 | `documentar-fase2` | Verify completo según método vigente, deuda explícita → REWRITE/APPEND en archive sólo tras autorización aplicable | Según paquetes tocados; no ship en este trabajo |

La revisión de este propose ALTO convoca arquitecto, analista funcional, QA y DevOps/seguridad de uno en uno sin leer los otros informes; auditor lee los originales completos. Registro de esas revisiones se agrega como documentos separados antes de presentar aprobación. No atribuirles aislamiento demostrado por promesa; no convertir su opinión en firma de Juani.

## Riesgos y alternativas descartadas

| Riesgo / alternativa | Decisión propuesta y costo de revertir (JUICIO) |
|---|---|
| Firma válida pero aprobación vacía o llave del agente | Validar identidad externa y alcance. Rechazar author/email como prueba; volver a sólo texto requiere decisión humana y migrar estados/consumidores, no borrar registros previos |
| No configurar firma impide todo trabajo | Mantener modo escrito explícito elegido por Juani; no fingir crypto. Exigir firma universal para este propose se descarta por la decisión J1 escrita de Juani. Revertir luego: migrar política y reaprobar alcances, sin firmar retroactivamente |
| Hook universal Stop, copia manual o parser privado estable | Descartados por Handback, campos nulos y sesiones. Núcleo+adapters cuesta más fixtures y mantenimiento de versiones; revertir a manual recuperando originales conserva evidencia, pero pierde automatización y debe declarar mediación |
| Captura “inmutable” desde el mismo usuario | No prometer autenticidad ni aislamiento fuertes; hash y registro detectan cambios accidentales. Endurecer con servicio externo requeriría otro diseño, operación y migración de confianza |
| Reemplazar permission entero como ddeb94f | Descartado por pérdida de claves hermanas. Revertir bash al snapshot puede quitar permisos legítimos posteriores: comparar diff y pedir revisión, no restaurar a ciegas |
| Agregar tipo agent o pisar roles personales | Descartado en primer incremento; config portable limita integración nativa. Migrar luego exige adapters/exportadores, manifest y compatibilidad por agente; no hay borrado automático de anteriores |
| Otra regex como solución semántica completa | Descartada; false confidence y tres vueltas ya observadas. Lector cuesta un despacho por verify y puede equivocarse; revertir S requiere aprobación de regla, retirar obligación/prompt/tests y mantener sus informes históricos |
| S duplica autoridad o genera cuarta vuelta | Sólo reporta al auditor y usa contador existente. Casos irresolubles quedan para decisión humana al límite; no reiniciar contador por captura o cambio de modelo |

## Preguntas para Juani con opciones

J1 resuelta por la decisión escrita del 23/10/2026; J2–J6 siguen (TBD). Las cinco decisiones de segunda revisión no aprueban apply ni responden las otras preguntas.

| ID | Pregunta | Opciones |
|---|---|---|
| J1 | Resuelta: aprobación escrita para este propose | Objetivo del primer change: SSH con lista de firmantes confiables; raíz/configuración final (TBD). Fuente: [decisiones de Juani](informes/2026-10-23-fase-2/segunda-revision/06-decisiones-juani.md) |
| J2 | ¿Aceptás núcleo de captura más adaptadores por versión? | A (recomendada): automatizar proveedor probado y fallback lector de originales; B: sólo subcomando lector al inicio. Ninguna opción permite resumen sustituto |
| J3 | ¿Dónde entregar role-*? | A (recomendada): harness config portable con recursos en skills y referencias explícitas; B: exportación nativa por agente con diseño adicional. No pisar personales en ninguna |
| J4 | ¿Qué fuente/licencia y modos aprobás para esos roles? | A: aportar/revisar siete definiciones canónicas y decidir modos en change propio; B: diferir instalación hasta tener fuente. No tomar archivos personales como bundle implícito |
| J5 | ¿Adoptás S como regla nueva en verify? | A (recomendada): lector separado bajo review, antes de senior/auditor, hasta tres vueltas; B: ensayo optativo sin exigirlo para cerrar hasta evaluar costo/calidad. Requiere elección propia, no hereda decisión 16 |
| J6 | ¿Autorizás reemplazar todas las reglas dentro de bash en el paso 10? | A: sí, con diff/snapshot y preservando hermanos; B: diferir ante reglas personales. Compatibilidad shell/bash demostrada primero; sin respuesta no hay escritura |

## Knowledge Base Impact

En apply/archive futuros: REWRITE de `ARCHITECTURE.md`, `README.md`, `AGENTS.md` y `CLAUDE.md` en espejo, `docs/metodologia-opsx.md` y `CHANGES.md` según comportamiento realmente entregado; APPEND de decisiones fechadas, lecciones L-NNN y deuda con costo/evidencia. Esta exploración y propuesta conservan su carácter histórico. Unión de impacto y diff determina alcance final. `knowledge-base-report.md` por change separa REWRITTEN/APPENDED/PENDING; no cerrar deuda sólo por asignarle destino ni presentar append-only como actualización suficiente. El presente trabajo crea únicamente documentos de explore/propose e informes de revisión; no modifica estas fuentes vigentes.

## Revisiones de propose

La primera revisión (01–05 de la carpeta principal) queda como historia: **custodia no verificable, no cuenta como revisión**. El [auditor de segunda revisión](informes/2026-10-23-fase-2/segunda-revision/05-auditor.md) registró seis bloqueantes de los roles más AUD-1; no hay veredicto vigente «sin bloqueantes». Este delta responde los siete; su cierre queda pendiente de confirmación del rol que levantó cada uno, sólo sobre su delta, y del auditor para AUD-1. Juani eligió ese procedimiento, sin otra ronda completa. El constructor no convoca revisores ni se atribuye confirmaciones. No es verify de código ni autorización de apply.

## Aprobación

| Elección | Persona | Fecha | Alcance |
|---|---|---|---|

## Respuesta a la segunda revisión

Fuente de criterios literales: [05-auditor.md](informes/2026-10-23-fase-2/segunda-revision/05-auditor.md). Estado de todos: corrección documental presentada, confirmación del emisor pendiente; no se fabrica cierre ni APROBADO.

### AN-1

Criterio literal: Cierra cuando §8 agrega el escenario "tabla completa sin commit A", con entrada, estado, código de salida, efecto sobre apply y quién arma el manifiesto. Ese caso tiene que aparecer en la fila 8 de pruebas.

Respuesta: §8 y fila 8 distinguen modo escrito (WRITTEN_UNVERIFIED, 2, alcance autorizado) y firmado (MISSING, 1, bloqueo), con manifiesto preparado por constructor y validado por verificador.

### AR-1

Criterio literal: Cierra cuando el propose nombra un destino versionado para los originales y quién valida igualdad y procedencia. Se comprueba con `git check-ignore <destino>`: tiene que salir con exit 1.

Respuesta: §9 fija openspec/changes/<change>/informes/, excepción acotada elegida por Juani y validación por capturador/approval, inspeccionada por auditor. La prueba del patrón se hace en repo temporal; el gate en el repo real sigue pendiente de implementar .gitignore, fuera de este alcance documental.

### AN-2

Criterio literal: Cierra cuando hay un único resultado para "senior encendido, auditor apagado". Las líneas 49/53 y 92 tienen que decir lo mismo, y la fila 82 tiene que nombrarlo.

Respuesta: §S, deuda CONF-3/RF2-3 y fila Deuda/S dicen: el senior informa y decide Juani por escrito, con y sin review.

### AR-2

Criterio literal: Cierra cuando §11 nombra quién resuelve la ruta y el hash, la dirección de la dependencia, la alternativa descartada y el costo de revertir.

Respuesta: §11 asigna resolución/hash al installer con adapter, enlace puro posterior a Compose, precondición del plan, alternativa relativa descartada y migración de contratos/goldens/chequeos/referencias.

### OPS-1

Criterio literal: Cierra cuando §9 define cómo desinstalar y actualizar los hooks propios según ownership, con una fila de prueba instalar → actualizar → desinstalar.

Respuesta: §9 define ownership, backup, conflictos, update y uninstall sin tocar personales; fila 9 hooks recorre el ciclo y fallos.

### QA-1

Criterio literal: Cierra cuando la fila 11 suma `./internal/install ./internal/uninstall ./internal/verify ./cmd/active-stack` con los casos discriminantes de QA-1.

Respuesta: Fila 11 amplía comando y exige plan real/home temporal, assets en snapshot/recuperación, fallo posterior y desinstalación preservando personal, definición vacía/faltante detectada por health checks y CLI distinta de clasificación vacía.

### AUD-1

Criterio literal: Cierra cuando esa línea refleja la segunda revisión y su estado.

Respuesta: Revisiones de propose sustituye el veredicto anterior por los siete bloqueantes de segunda revisión y confirmaciones pendientes; LEEME rotula 01–05 como historia sin custodia verificable.

Las cinco decisiones de Juani se aplican así: aprobación escrita ahora y objetivo SSH en change 1 (J1); originales en repo con excepción acotada (§9); senior informa/Juani decide (§S y deuda); confirmación individual sobre delta, sin convocatoria del constructor (Revisiones); primera revisión histórica sin cómputo (LEEME y Revisiones). La tabla de aprobación permanece sin completar por el constructor.
