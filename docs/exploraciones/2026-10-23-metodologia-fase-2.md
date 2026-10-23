# Exploración: metodología PropTech en OPSX, fase 2

Fecha: 2026-10-23. Rama observada: `feat/metodologia-fase-2`. Base observada: `0204022`, con fase 1 integrada. Estado: EXPLORE documental; ninguna decisión de implementación aprobada.

## 1. Pregunta y fuentes internas

¿Cómo convertir aprobación, originales de revisores, permisos y entrega de roles en contratos comprobables, y detectar contradicciones de sentido que sobreviven a los tests de texto?

Se leyeron la [propuesta de fase 1](../propuestas/2026-10-23-metodologia-proptech-en-opsx.md), especialmente §§5 y 8; la [exploración original](2026-10-23-metodologia-proptech.md), §§7–8; los informes [21](../propuestas/informes/2026-10-23/21-auditor-vuelta-3.md), [24](../propuestas/informes/2026-10-23/24-archive.md) y [decisión 16](../propuestas/informes/2026-10-23/16-decision-juani-rf2-1.md); la [base vigente](../metodologia-opsx.md), y las fuentes enumeradas abajo. Los informes históricos no se reescriben.

EVIDENCIA de inspección en esta sesión: `git status --short --branch` mostró la rama indicada y sólo logs preexistentes sin seguimiento; `git log -4 --oneline` mostró `0204022`; `openspec list` mostró únicamente `instalacion-sana`, 4/8 tareas; `openspec status --change instalacion-sana` mostró proposal/tasks completos y design/specs ausentes. No había change fase 2 registrado al momento de esa consulta. Después se creó con `openspec new change metodologia-fase-2`: alcance documental, proposal referencial y `skip_specs: true` limitado a esta entrega sin cambios de comportamiento. `openspec list --specs` devolvió `No specs found` (11:14 -03). Los futuros changes de implementación no heredan esa exención. Este documento y el propose son artefactos documentales OPSX; no se acredita apply ni archive por crearlos. La captura del horario exacto de esas consultas no está disponible: no se inventa.

## 2. Lo que compone hoy el repo

EVIDENCIA por lectura, sin ejecutar Compose: `internal/harness/config/compose.go` reconoce doce toggles y agrega governance, evidence, auditor, roles, review, knowledge-base, senior, ship, engram y tdd después de la base, con delegation/model-routing sustractivos. Usa variantes de ocho agentes, con goldens completos Claude/OpenCode. La inspección de `compose_all_toggles_claude.golden` confirmó en el resultado compuesto la clave compartida en líneas 169/1797 frente a round/role en 936, la sección antigua sin veto en 959–986 y la espera incondicional del auditor en 1474; no es sólo una inferencia de fragmentos aislados. `assets/commands/claude/opsx/verify.md` y `ship.md` son consumidores adicionales, no equivalentes al bloque de todos los agentes.

`review.md` ya establece conformidad → rompedor → senior activo → auditor, tres vueltas, originales sin Write para revisores e índice de procedencia. Conserva una sección que permite cerrar hallazgos por decisión personal sin ubicarla claramente frente al auditor (RF1-4). `senior.md` exige un auditor incluso apagado (CONF-3/RF2-3). `engram-protocol.md` y las bases aún ofrecen claves por fase, frente a las claves por vuelta/rol de review (RF1-6). `internal/contracttest/text.go` limita vocabulario conocido: no interpreta reglas arbitrarias.

`cmd/active-stack/roles_dispatch.go` clasifica un diff mediante `internal/roles`; no instala ni convoca roles. La enumeración de assets no muestra `role-*.md` bundleados. El adapter expone instrucciones, skills y comandos; un destino de definiciones de rol necesita contrato propio o ubicación explícita dentro de skills. No se debe deducir que todo agente admite agentes Markdown nativos.

## 3. Paso 8: firma no es identidad ni alcance por sí sola

Fuentes oficiales consultadas el 23/10: [git verify-commit](https://git-scm.com/docs/git-verify-commit) valida la firma de un commit; [git-config](https://git-scm.com/docs/git-config), `gpg.format`, `gpg.ssh.allowedSignersFile`, `gpg.ssh.revocationFile`, describe configuración de GPG/SSH y confianza SSH. Es necesario fijar identidad autorizada, no confiar en nombre/email de autor ni en configuración controlable por el agente.

JUICIO: proponer `active-stack approval --change`, inicialmente verificador de sólo lectura. Un commit de aprobación A firmado por identidad autorizada debe contener la decisión y un manifiesto con change, elección, alcance, SHA completo P del propose revisado y hashes de sus documentos/informes. A no puede referenciar su propio SHA. Se verifica A y se resuelve P; cualquier cambio de alcance o material normativo exige otra aprobación. Rebases que cambien P requieren revalidación. Ni el agente completa la tabla ni una firma de documentación incidental autoriza apply.

Actualización por [decisión escrita de Juani del 23/10/2026](../propuestas/informes/2026-10-23-fase-2/segunda-revision/06-decisiones-juani.md): aprobación escrita para este propose y SSH con lista de firmantes confiables como objetivo del primer change. Raíz de confianza/configuración final (TBD). Si se adopta modo firmado y falta llave/verificador/confianza, apply ALTO se detiene. El modo escrito conserva la regla vigente de fase 1 y su limitación probatoria; tabla completa sin A se desarrolla en §8 del propose. Los originales se versionarán bajo openspec/changes/*/informes/ mediante excepción acotada, todavía no instalada en este trabajo documental. Con senior encendido y auditor apagado, el senior informa y decide Juani por escrito, con y sin review.

## 4. Paso 9: contratos reales de captura

Las siguientes son observaciones documentales, no ensayos de los clientes instalados. Versiones mínimas ensayadas: (TBD); en apply se fijarán versiones y fixtures por proveedor.

| Proveedor y fuente oficial abierta el 23/10 | Contrato observado | Consecuencia de diseño |
|---|---|---|
| [Claude Code hooks](https://code.claude.com/docs/en/hooks#subagentstop) | `SubagentStop` incluye `agent_id`, `agent_type`, `agent_transcript_path`, `last_assistant_message` y `stop_hook_active`; `transcript_path` pertenece al padre. Desde v2.1.271, con `SubagentHandback`, el informe está en `tool_input.message`, recibido por hook de esa herramienta; el último texto puede ser sólo despedida. También hay eventos de agentes internos. | Filtrar despachos registrados; capturar Handback cuando aplique y usar Stop como cierre correlacionado. No tratar todo Stop ni su último texto como informe íntegro. Transcript puede ir atrasado; recuperar con validación, sin resumir. |
| [Codex hooks](https://learn.chatgpt.com/docs/hooks#subagentstop) | `SubagentStop` documenta identificadores, `turn_id`, texto y transcript opcionales; formato del transcript no estable. En salida exitosa espera JSON; `decision: block` continúa el subagente, `continue: false` prevalece. | Adapter propio, no copiar salida Claude a ciegas. Campos nulos implican captura pendiente. Preferir evento contractual, conservar bruto y versionar cualquier lector de registro. No confundir continuación con aprobación. |
| [OpenCode plugins](https://opencode.ai/docs/plugins/) y [SDK](https://opencode.ai/docs/sdk/#sessions) | Plugins reciben eventos de mensajes y `session.idle`/`session.status`; el SDK permite `session.messages` y `session.message`, con `info` y `parts`. Estas páginas no establecen un equivalente portable a SubagentStop ni que idle pruebe informe completo. | Adaptar evento como disparador y leer mensaje de la sesión hija identificada en el despacho. Correlación padre/hija y finalización exacta requieren fixture/versionado (TBD). No raspar una base privada como API estable. |

JUICIO: núcleo común `reports capture` (nombre candidato), con adaptadores por proveedor. El hook entrega eventos al capturador; el subcomando recupera desde un registro original identificado, nunca desde una paráfrasis del orquestador. Evento/disco/API originales quedan separados de Markdown derivado. Un hash prueba integridad posterior, no autenticidad frente al mismo usuario con acceso al filesystem. Sin original recuperable, informe FALTANTE y cierre detenido.

## 5. Paso 10: reimplementación acotada y compatibilidad

EVIDENCIA: `git show ddeb94f -- internal/harness/config/permissions/overlays.go` muestra reemplazo de `permission` entero. La fuente actual usa `opencodePermissionShape`, `orderedBashBlock` y claves escalares explícitas en todos los tiers. `TestOpencodeBashResidueSurvivesTierSwitch` en `tier_switch_test.go` exige que sobrevivan reglas anteriores y ordena sustituirlo cuando exista el invariante fuerte.

JUICIO: reimplementar sin cherry-pick: sentinel interno `__replace__` sólo en `permission.bash`. Conservar claves hermanas personales, reemplazar deliberadamente reglas dentro de bash sólo con advertencia concreta, snapshot y aprobación de alcance; no prometer conservarlas a la vez que se reemplaza ese objeto. El sentinel es del merge del instalador, nunca debe quedar en settings finales. El orden final en disco debe conservar wildcard primero y deny al final, también después de MergeJSONObjects.

La [documentación oficial de permisos de OpenCode](https://opencode.ai/docs/permissions/) consultada usa hoy `shell`, y documenta última coincidencia ganadora. El repo usa `bash`: incompatibilidad potencial que exige fixture de versión real antes de apply. No se renombra a shell silenciosamente ni se acredita compatibilidad con un cliente actual por probar sólo JSON. Si la versión objetivo no acepta bash, proponer ajuste de alcance y obtener aprobación ALTO; no expandir este paso automáticamente.

## 6. Paso 11 y S

Roles: proponer harness `config` con assets bundleados, sin tipo nuevo `agent`. Ubicación portable recomendada: definiciones Markdown en directorio propio de Active Stack bajo skills resuelto por adapter, leídas por ruta explícita; traducción a agentes nativos fuera del primer incremento. Fuente canónica/licencia de los siete `role-*`: (TBD), no copiar configuración personal. Un rol personal homónimo se informa y permanece intacto; nunca se asume que pertenece al instalador por su nombre. Manifiesto de propiedad/hashes y backup preceden actualización. Si se inyectan referencias, usar markers idempotentes.

S es regla nueva según decisión 16, todavía no adoptada. JUICIO: un lector separado recibe bloque compuesto efectivo y comandos aplicables, toggle set, decisiones fechadas y manifiesto de revisión. Busca permisos semánticos opuestos incluso junto a frases correctas; informa pares de cláusulas, caso que diverge y alcance. No reemplaza lentes, senior ni auditor. Su informe se incorpora al auditor y al máximo existente de tres vueltas. Disponibilidad de modelo, contaminación, costo y falibilidad se declaran. No es garantía automática ni prueba de conducta real.

## 7. Deuda que conviene pagar aquí

| Deuda de 21/24 | Propuesta de destino | Criterio de cierre futuro |
|---|---|---|
| CONF-3/RF2-3 | Antes de S: senior y degradaciones | Sin auditor no espera veredicto inexistente; decisión humana escrita, sin APROBADO fabricado |
| RF1-4 y RF2-4 | Texto review/KB y contratos antes de S | Autoridad consistente en todo el bloque; corregir `Mark` sin fijar capitalización incidental |
| RF1-6 | Con captura y claves de reportes | Ningún save de verify comparte clave entre roles/vueltas; tampoco instrucción base contradictoria |
| RF3-2 | Corpus S y pruebas ship | Cero runs y `treated as green` rechazan entrega; casos positivos válidos no se rechazan por citar el ejemplo |
| N1/N1b/N2b/N3/N4/N5/N8 | Corpus semántico, bloque y comandos | Contradicciones aditivas y paráfrasis revisadas, no sólo borrado de frases |
| CONF-1 | Sigue en instalacion-sana | No pagada por docs ni por esta fase; receta §3.6 y OK escrito para reinstalar |

RF3-1 ya fue cerrado excepcionalmente con evidencia y autorización posterior; no se reabre una cuarta vuelta. Compactación (12), UI y reinstalación real quedan fuera.

## 8. Resultado y límites

La viabilidad exige separación de verificación, confianza y alcance; captura por proveedor y versión; propiedad explícita de archivos; y revisión semántica falible con auditor como decisor. El [propose de fase 2](../propuestas/2026-10-23-metodologia-fase-2.md) convierte estos hallazgos en pasos revisables. Sólo J1 y las otras cuatro decisiones de segunda revisión constan respondidas en el propose; las restantes preguntas siguen (TBD). No se corrieron tests de producto ni build: sólo lectura, investigación y documentos.
