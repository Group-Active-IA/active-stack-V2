# Rompedor, vuelta 1 (Claude opus, sacado textual del registro del sub-agente)

**Informe RF1: rompedor, fase 1** (`4748680..299889d`, 23/10/2026 -03, tres lentes en un pase)

**Veredicto: 3 bloqueantes y 4 notas.** Los tests nuevos sólo buscan frases, así que no detectan cuando el texto se invierte. Los 9 mutantes que probé sobrevivieron a los tests de contrato.

**RF1-1, BLOQUEANTE (verdad de los tests).** EVIDENCIA, 10:17:48 a 10:18:08: apliqué 9 mutantes de una línea en un worktree con HOME temporal, y `go test -run <test del contrato>` dio `ok` en los 9:
- M1, `ship.md:205` del fragmento: `+... timed_out are not green. A cancelled job counts as green when its rerun passes.` Pasa `TestShipEsperaCICompleto`.
- M1b, el mismo cambio en `assets/commands/claude/opsx/ship.md:399`. Pasa el paquete `./assets` completo.
- M2, `review.md` en The loop: `+... If blockers remain, start round 4 automatically.` Pasa `TestVerifyContratoBucle`.
- M3: `Wait for the senior report before launching the auditor.` cambiado por `Launch the auditor in parallel with the senior.` Pasan `TestVerifyContratoBucle` y `TestAuditorDecide*`.
- M4, `auditor.md:51`: borré `Missing required reports prevent approval.` Pasan los tests del auditor.
- M5, `verify.md:140` del comando: `If all clear: "All checks passed. Proceed to archive."` Pasa `TestVerifyConformidadNoHabilitaArchive`.
- M6, prompt de `knowledge-base.md:191`: `Only append lessons.` cambiado por `Rewrite old lessons freely.` Pasa `TestArchiveBaseCompletaConEvidencia`.
- M7, `governance.md`: `+An agent may fill it after chat approval.` Pasa `TestGovernance*`.
- M8, `internal/verify/delivery.go`: `strings.Count(...) != 1` cambiado por `< 1`, que acepta bloques duplicados. Pasan `./internal/verify` y `./internal/install` completos.

Con el paquete completo, los goldens frenan M2. JUICIO: por la misma vía frenarían M1, M3, M4, M6 y M7 (sólo corrí M2). Pero un golden es una foto del texto, y el constructor los regeneró en cada paso. M1b, M5 y M8 no tienen golden y pasan la suite entera.

§3.3 y §3.4 del propose exigían que fallaran los tests propios ante "permitir cuarta vuelta", "auditor antes del senior", "anunciar archive" y "aceptar cancelado en ambas copias". Ninguno de esos casos falla por su aserción. JUICIO, confianza alta: los "14 + 7 mutantes detectados" del constructor eran borrados de frases, no inversiones de sentido.

**RF1-2, BLOQUEANTE (correctitud).** `knowledge-base.md:191`, dentro del prompt que se pega al archivista: `Mark auditor source absent and do not invent its lists; closure requires Juani's written decision on degraded scope.` La frase es incondicional y la misma línea pide sacar las LECCIONES y la DEUDA del auditor. El archivista arranca sin contexto, así que con el auditor activo marcaría la fuente como ausente y pediría firma de degradación en cada archive. El test sólo comprueba que la frase exista.

**RF1-3, BLOQUEANTE (correctitud).** `ship.md:207` del fragmento y `assets/commands/claude/opsx/ship.md:395`: `For this repository, require .github/workflows/ci.yml on push, with test on Ubuntu and Windows.` El bloque se instala global en `~/.claude/CLAUDE.md`, así que en cualquier otro proyecto "this repository" es el repo de turno. Un agente exigiría jobs inexistentes y nunca daría verde, o chocaría con la línea siguiente ("In other repositories..."). §3.4 se refería a active-stack, y el test fija la frase ambigua.

**RF1-4, NOTA (correctitud).** `review.md:69-77` mantiene "The review **does not block**. It holds no veto" y "Closing over an open finding is legitimate... belongs to the person". The loop dice que sólo APROBADO habilita archive y que la persona decide recién en la vuelta 3. Esa sección nunca dice "decide el auditor, frena el orquestador". Confianza media.

**RF1-5, NOTA (comentarios).** `ship.md:179` todavía cita "The senior fragment... names this same split — **the senior labels, you hold**", una frase que el paso 5 borró del senior. `senior.md:17` dice "It saw the earlier phases of this change", y ahora sólo corre en verify.

**RF1-6, NOTA (cableado).** `engram-protocol.md:11` y `claude/sdd-orchestrator.md:169` siguen entregando a todo sub-agente `topic_key: "opsx/{change-name}/{phase}"`. Choca con la clave por vuelta y rol de `review.md`, que avisa que una clave común filtra conclusiones.

**RF1-7, NOTA (cableado).** El cableado se sostiene.
- `BinaryChecks` y `composedDeliveryCheck` se llaman desde el flujo real: `main.go:209/267` y `windows_dispatch.go:123` llaman a `buildInstallVerifyHook`, y éste a `verify.BuildHook`. A su vez, `harness_checks.go:49` agrega el chequeo de entrega.
- Claude, Codex y OpenCode implementan `VariantKey`, así que el chequeo de entrega no se saltea en ninguno.
- `TestEntregaComparaCompose` cubre los tres adapters.
- El hueco es M8: nada prueba la regla "sin bloques duplicados" de §3.6.
- En todo el bloque compuesto no queda ninguna mención de "Ready for archive", "three independent passes", "review-verdict" ni `haiku` como modelo de archive.

**SUPUESTOS**
- Que regenerar goldens en cada paso es lo que anula su poder de detección. Lo leí en el informe 09, no lo reproduje.
- `cmd/active-stack/headless` falló en mi corrida base por no encontrar node en mi PATH restringido. No lo atribuyo al diff.

**Archivos que abrí:** el propose (§3 a §5); los informes 08 y 09; los diffs de `review.md`, `auditor.md`, `governance.md`, `roles.md`, `ship.md`, `knowledge-base.md` y de los comandos `opsx/verify.md` y `opsx/ship.md`; el golden de Claude (secciones auditor, review, senior, ship, TDD y modelos); `internal/verify/{binary,delivery,harness_checks}.go` y sus tests; `cmd/active-stack/verify_binary.go`; `config/install.go`; `system/binary_install_dir.go`; los tests `verify_loop`, `ship_ci` (los dos), `governance_approval`, `tdd_toggle`, `roles_toggle`, `knowledge_base` y `verify_phase`, más `assets/commands_test`. No leí el golden de OpenCode línea por línea: sólo corrí su test en el pase de M2.

**Estado final:**
- `git -C /home/juanisarmiento/Proyectos/active-stack-V2 status --short`: sólo los 9 `?? docs/propuestas/informes/2026-10-23/.0*.log`, que ya estaban antes.
- Sin cambios míos.
- `git worktree list` muestra sólo el árbol principal; el worktree `as-rompedor` está borrado.
