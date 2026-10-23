# Constructor, tanda A

Alcance: pasos 0, 1, 2 y 3 del propose firmado en 4748680. Paso 0 limitado a código y diagnóstico por pedido de Juani. Sin reinstalación, instalaciones, push ni lectura de .env. No se tocaron los informes preexistentes sin versionar.

Procedencia: informe escrito por el constructor en esta sesión, con salida textual de herramientas. Proveedor OpenAI; identidad Codex. Modelo/familia exactos sin registro de despacho: SUPUESTO. No se convocaron revisores ni se afirma independencia. Se leyeron CLAUDE.md, propuesta y exploración; no existe knowledge-base/ ni archivo de lecciones local. No se recibieron casos dictados por un auditor en esta entrega.

Las pruebas corren con bwrap: raíz en solo lectura, escritura habilitada en el checkout, /tmp y caché de Go. Los destinos reales ~/.claude, ~/.local/bin y /usr/local/bin quedan en solo lectura. No se hizo go build. Son pruebas locales, no observación de CI remoto.

## Paso 0: instalacion-sana

Diagnóstico estático: internal/install/self_install.go usa os.Executable cuando sourcePath está vacío y copia ese archivo como active-stack. En un test eso puede ser install.test. TestExternalStepCallsExternalInstaller, TestVerifyHookStepRunsAfterApply y TestVerifyHookNilIsNoOp, entre otros, ejecutan planes con el helper buildOptions, que deja NoSelfInstall=false y SelfInstallBinaryPath vacío. No aíslan el resolver del directorio binario. En el código actual pasan HomeDir=t.TempDir(): si /usr/local/bin no es escribible, el destino es ese home temporal, no el home real. Si es escribible, eligen /usr/local/bin. TestBuildPlan_SelfInstallGoosMatchesSnapshotGoos también copia el test, pero con resolver inyectado bajo t.TempDir().

La vía histórica que ignoraba --home está documentada y corregida en self_install.go y TestSelfInstall_WritesBinaryUnderGivenHome. Este último usa un archivo fuente ficticio y dos homes temporales. JUICIO: el auto-install no aislado explica cómo una suite puede copiar un .test; no hay evidencia suficiente para atribuir el archivo real a un test único ni confirmar cuándo ocurrió. No se reprodujo una escritura en destinos reales.

Cambios: version informa versión y revisión disponibles mediante build info; verify valida identidad cmd/active-stack antes de ejecutar roles --help, con timeout. El hook de instalación incorpora el chequeo salvo --no-self-install. Los adapters reales verifican además contenido contra Compose(variante, toggles), incluyendo el prompt primario de OpenCode en settings. Los checks de markers y mode:primary siguen vigentes. No se modificaron backup, filemerge, pipeline ni configuración del usuario. Sin cambios de goldens.

### Rojo observado

```text
Hora: 2026-10-23 09:54:02 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./cmd/active-stack -run ^TestVersionSubcommand$ -count=1
--- FAIL: TestVersionSubcommand (0.00s)
    version_test.go:20: version debe salir 0 e identificar el ejecutable: exit status 1, unknown command: "version"
FAIL
FAIL	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.007s
FAIL
```

```text
Hora: 2026-10-23 09:54:19 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./internal/verify -run ^TestEntregaComparaCompose$ -count=1
--- FAIL: TestEntregaComparaCompose (0.00s)
    --- FAIL: TestEntregaComparaCompose/claude (0.00s)
        delivery_test.go:38: falla=false, esperada=true: [{ID:config:sdd-orchestrator:claude Description:idempotent marker present exactly once in claude instructions (sdd-orchestrator) Status:passed Error:}]
    --- FAIL: TestEntregaComparaCompose/codex (0.00s)
        delivery_test.go:38: falla=false, esperada=true: [{ID:config:sdd-orchestrator:codex Description:idempotent marker present exactly once in codex instructions (sdd-orchestrator) Status:passed Error:}]
    --- FAIL: TestEntregaComparaCompose/opencode (0.00s)
        delivery_test.go:38: falla=false, esperada=true: [{ID:config:sdd-orchestrator:opencode Description:primary agent "sdd-orchestrator" registered (mode:primary) in opencode settings Status:passed Error:}]
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/verify	0.007s
FAIL
```

```text
Hora: 2026-10-23 09:54:35 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./internal/verify -run ^TestBinarioInstalado$ -count=1
FAIL	github.com/Group-Active-IA/active-stack/internal/verify [build failed]
FAIL
```

```text
Hora: 2026-10-23 09:55:46 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./cmd/active-stack -run ^TestVerifyInstalacionIncluyeBinario$ -count=1
FAIL	github.com/Group-Active-IA/active-stack/cmd/active-stack [build failed]
FAIL
```

El rojo de binario y wiring es de compilación por API todavía inexistente; version y Compose fallaron en sus aserciones contra el comportamiento viejo.

### Verde previo al commit

```text
Hora: 2026-10-23 09:55:10 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./internal/verify ./cmd/active-stack -run TestBinarioInstalado|TestEntregaComparaCompose|TestVersionSubcommand -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.011s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.008s
Hora: 2026-10-23 09:56:00 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	0.030s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.835s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.813s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.904s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	0.012s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	0.009s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	0.013s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	0.007s
ok  	github.com/Group-Active-IA/active-stack/internal/backup	0.114s
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	0.188s
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	0.018s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	0.025s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.812s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	0.123s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	0.082s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	0.037s
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	0.013s
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.256s
ok  	github.com/Group-Active-IA/active-stack/internal/model	0.007s
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	0.004s
ok  	github.com/Group-Active-IA/active-stack/internal/planner	0.005s
ok  	github.com/Group-Active-IA/active-stack/internal/roles	0.100s
ok  	github.com/Group-Active-IA/active-stack/internal/system	0.020s
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.031s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	0.021s
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.034s
```

### Cierre del paso 0, después del commit

```text
Hora: 2026-10-23 09:57:03 -0300
SHA: 84fcf0ebc0a6f647c04b5b0eb4dfbbc9939c866d
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	(cached)
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	(cached)
```

## Paso 1: opencode-mcp-command-array

Cherry-pick sin commit de fae8bbf, disponible desde juancruz (remote local ya agregado). Conflicto de mcp_test.go resuelto conservando el comentario del schema command array con formato gofmt. Para observar el rojo se guardó la versión entrante de mcp.go, se restituyó temporalmente la de HEAD, se corrieron los tests entrantes y recién después se aplicó producción. Se agregó aserción de environment porque el cherry-pick también cambia ese campo. No se tomó ddeb94f ni se tocaron permisos. Se eliminó el trailer de coautoría y se usa el nombre del change en el conventional commit en español. Sin goldens.

### Rojo observado

```text
Hora: 2026-10-23 09:57:29 -0300
SHA: 84fcf0ebc0a6f647c04b5b0eb4dfbbc9939c866d
Comando: go test ./internal/harness/external -run ^TestBuildStdioOverlay_OpenCode$ -count=1
--- FAIL: TestBuildStdioOverlay_OpenCode (0.00s)
    mcp_test.go:521: command = engram, want [engram mcp] (OpenCode schema needs binary+args combined)
    mcp_test.go:524: OpenCode local overlay must not have a separate 'args' key
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/external	0.004s
FAIL
```

```text
Hora: 2026-10-23 09:57:45 -0300
SHA: 84fcf0ebc0a6f647c04b5b0eb4dfbbc9939c866d
Comando: go test ./internal/harness/external -run ^TestBuildStdioOverlay_OpenCode$ -count=1
--- FAIL: TestBuildStdioOverlay_OpenCode (0.00s)
    mcp_test.go:521: command = engram, want [engram mcp] (OpenCode schema needs binary+args combined)
    mcp_test.go:524: OpenCode local overlay must not have a separate 'args' key
    mcp_test.go:527: environment = <nil>, se esperaba LOG=1
    mcp_test.go:529: OpenCode no admite env separado de environment
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/external	0.004s
FAIL
```

### Verde previo al commit

```text
Hora: 2026-10-23 09:57:46 -0300
SHA: 84fcf0ebc0a6f647c04b5b0eb4dfbbc9939c866d
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.705s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.869s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.825s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	0.080s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.299s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.037s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	(cached)
```

### Cierre del paso 1, después del commit

```text
Hora: 2026-10-23 09:58:10 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	(cached)
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	(cached)
```

## Paso 2: ship-espera-ci

Se actualizan ambas copias: el cierre incluye resultado de CI sobre el SHA pusheado, con todos los jobs requeridos en success; se registra conjunto esperado, eventos y matriz antes de evaluar, se inspecciona el SHA origen del PR y se invalida el resultado tras otro push. Sin runs, sin CI, error de consulta o estados incompletos dejan la entrega sin confirmar. La pregunta por main y la prohibición de merge siguen vigentes. No se hizo push: no hay evidencia de CI remoto.

Tests propios por copia, primero rojos contra los textos viejos. Mutantes aplicados uno por vez sobre el asset real con copia previa en /tmp, restaurados en finally, sin commit de mutaciones. Los siete por copia aceptaban parcial, cancelado, omitido, otro SHA, cero runs, consulta fallida o un solo job suficiente. Todos fallaron en la aserción del contrato de CI.

Goldens regenerados únicamente después del rojo específico: 35 líneas agregadas y 3 reemplazadas en cada variante Claude/OpenCode. Diff revisado: cierre PR + resultado CI, sección de CI completo y costo de consultas. No cambia otro fragmento ni los modelos.

### Rojo observado

```text
Hora: 2026-10-23 09:58:49 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./internal/harness/config ./assets -run TestShipEsperaCICompleto|TestComandoShipEsperaCICompleto -count=1
--- FAIL: TestShipEsperaCICompleto (0.00s)
    --- FAIL: TestShipEsperaCICompleto/Before_evaluating_CI,_record_the_pushed_SHA,_expected_workflows/events_and_expanded_matrix_jobs. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: Before evaluating CI, record the pushed SHA, expected workflows/events and expanded matrix jobs.
    --- FAIL: TestShipEsperaCICompleto/Every_required_job_must_finish_in_success_on_the_pushed_SHA. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: Every required job must finish in success on the pushed SHA.
    --- FAIL: TestShipEsperaCICompleto/gh_run_list_--commit_<sha> (0.00s)
        ship_ci_test.go:29: falta contrato de CI: gh run list --commit <sha>
    --- FAIL: TestShipEsperaCICompleto/gh_run_watch (0.00s)
        ship_ci_test.go:29: falta contrato de CI: gh run watch
    --- FAIL: TestShipEsperaCICompleto/Inspect_the_jobs_of_every_expected_run;_never_choose_only_the_favorable_run. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: Inspect the jobs of every expected run; never choose only the favorable run.
    --- FAIL: TestShipEsperaCICompleto/For_this_repository,_require_.github/workflows/ci.yml_on_push,_with_test_on_Ubuntu_and_Windows. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: For this repository, require .github/workflows/ci.yml on push, with test on Ubuntu and Windows.
    --- FAIL: TestShipEsperaCICompleto/For_a_PR,_inspect_its_required_checks_and_verify_the_source_SHA,_not_the_synthetic_merge_SHA. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: For a PR, inspect its required checks and verify the source SHA, not the synthetic merge SHA.
    --- FAIL: TestShipEsperaCICompleto/In_other_repositories,_derive_the_required_set_from_their_workflows_and_required_checks;_a_tag_release_is_not_required_for_a_branch_push. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: In other repositories, derive the required set from their workflows and required checks; a tag release is not required for a branch push.
    --- FAIL: TestShipEsperaCICompleto/One_successful_job_with_another_pending_is_not_green. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: One successful job with another pending is not green.
    --- FAIL: TestShipEsperaCICompleto/queued,_in_progress,_cancelled,_skipped,_neutral_and_timed_out_are_not_green. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: queued, in_progress, cancelled, skipped, neutral and timed_out are not green.
    --- FAIL: TestShipEsperaCICompleto/A_result_for_another_SHA_is_not_green. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: A result for another SHA is not green.
    --- FAIL: TestShipEsperaCICompleto/Zero_runs,_no_CI_or_a_failed_query_are_not_green;_report_the_state_and_leave_delivery_unconfirmed. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: Zero runs, no CI or a failed query are not green; report the state and leave delivery unconfirmed.
    --- FAIL: TestShipEsperaCICompleto/On_failure,_report_it_without_stacking_another_fix_as_if_the_push_were_confirmed. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: On failure, report it without stacking another fix as if the push were confirmed.
    --- FAIL: TestShipEsperaCICompleto/A_new_push_invalidates_the_previous_result_and_requires_evaluating_the_new_SHA. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: A new push invalidates the previous result and requires evaluating the new SHA.
    --- FAIL: TestShipEsperaCICompleto/The_orchestrator_waits_for_the_complete_CI_result_before_declaring_delivery_ready. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: The orchestrator waits for the complete CI result before declaring delivery ready.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.005s
--- FAIL: TestComandoShipEsperaCICompleto (0.00s)
    --- FAIL: TestComandoShipEsperaCICompleto/Before_evaluating_CI,_record_the_pushed_SHA,_expected_workflows/events_and_expanded_matrix_jobs. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: Before evaluating CI, record the pushed SHA, expected workflows/events and expanded matrix jobs.
    --- FAIL: TestComandoShipEsperaCICompleto/Every_required_job_must_finish_in_success_on_the_pushed_SHA. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: Every required job must finish in success on the pushed SHA.
    --- FAIL: TestComandoShipEsperaCICompleto/gh_run_list_--commit_<sha> (0.00s)
        ship_ci_test.go:35: falta contrato de CI: gh run list --commit <sha>
    --- FAIL: TestComandoShipEsperaCICompleto/gh_run_watch (0.00s)
        ship_ci_test.go:35: falta contrato de CI: gh run watch
    --- FAIL: TestComandoShipEsperaCICompleto/Inspect_the_jobs_of_every_expected_run;_never_choose_only_the_favorable_run. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: Inspect the jobs of every expected run; never choose only the favorable run.
    --- FAIL: TestComandoShipEsperaCICompleto/For_this_repository,_require_.github/workflows/ci.yml_on_push,_with_test_on_Ubuntu_and_Windows. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: For this repository, require .github/workflows/ci.yml on push, with test on Ubuntu and Windows.
    --- FAIL: TestComandoShipEsperaCICompleto/For_a_PR,_inspect_its_required_checks_and_verify_the_source_SHA,_not_the_synthetic_merge_SHA. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: For a PR, inspect its required checks and verify the source SHA, not the synthetic merge SHA.
    --- FAIL: TestComandoShipEsperaCICompleto/In_other_repositories,_derive_the_required_set_from_their_workflows_and_required_checks;_a_tag_release_is_not_required_for_a_branch_push. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: In other repositories, derive the required set from their workflows and required checks; a tag release is not required for a branch push.
    --- FAIL: TestComandoShipEsperaCICompleto/One_successful_job_with_another_pending_is_not_green. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: One successful job with another pending is not green.
    --- FAIL: TestComandoShipEsperaCICompleto/queued,_in_progress,_cancelled,_skipped,_neutral_and_timed_out_are_not_green. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: queued, in_progress, cancelled, skipped, neutral and timed_out are not green.
    --- FAIL: TestComandoShipEsperaCICompleto/A_result_for_another_SHA_is_not_green. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: A result for another SHA is not green.
    --- FAIL: TestComandoShipEsperaCICompleto/Zero_runs,_no_CI_or_a_failed_query_are_not_green;_report_the_state_and_leave_delivery_unconfirmed. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: Zero runs, no CI or a failed query are not green; report the state and leave delivery unconfirmed.
    --- FAIL: TestComandoShipEsperaCICompleto/On_failure,_report_it_without_stacking_another_fix_as_if_the_push_were_confirmed. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: On failure, report it without stacking another fix as if the push were confirmed.
    --- FAIL: TestComandoShipEsperaCICompleto/A_new_push_invalidates_the_previous_result_and_requires_evaluating_the_new_SHA. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: A new push invalidates the previous result and requires evaluating the new SHA.
    --- FAIL: TestComandoShipEsperaCICompleto/The_orchestrator_waits_for_the_complete_CI_result_before_declaring_delivery_ready. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: The orchestrator waits for the complete CI result before declaring delivery ready.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.003s
FAIL
```

```text
Hora: 2026-10-23 09:59:16 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./internal/harness/config -run ^TestCompose_AllToggles -count=1
--- FAIL: TestCompose_AllToggles (0.00s)
    compose_test.go:59: Compose output does not match golden.
        got len=104023 want len=102616
--- FAIL: TestCompose_AllToggles_opencode (0.00s)
    compose_test.go:83: Compose output does not match golden.
        got len=104365 want len=102958
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.006s
FAIL
```

### Mutantes, salida textual

2-mutante-comando-cancelado

```text
Hora: 2026-10-23 09:59:43 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1
--- FAIL: TestComandoShipEsperaCICompleto (0.00s)
    --- FAIL: TestComandoShipEsperaCICompleto/queued,_in_progress,_cancelled,_skipped,_neutral_and_timed_out_are_not_green. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: queued, in_progress, cancelled, skipped, neutral and timed_out are not green.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.003s
FAIL
```

2-mutante-comando-cero-runs

```text
Hora: 2026-10-23 09:59:44 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1
--- FAIL: TestComandoShipEsperaCICompleto (0.00s)
    --- FAIL: TestComandoShipEsperaCICompleto/Zero_runs,_no_CI_or_a_failed_query_are_not_green;_report_the_state_and_leave_delivery_unconfirmed. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: Zero runs, no CI or a failed query are not green; report the state and leave delivery unconfirmed.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.002s
FAIL
```

2-mutante-comando-error-consulta

```text
Hora: 2026-10-23 09:59:44 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1
--- FAIL: TestComandoShipEsperaCICompleto (0.00s)
    --- FAIL: TestComandoShipEsperaCICompleto/Zero_runs,_no_CI_or_a_failed_query_are_not_green;_report_the_state_and_leave_delivery_unconfirmed. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: Zero runs, no CI or a failed query are not green; report the state and leave delivery unconfirmed.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.003s
FAIL
```

2-mutante-comando-job-incompleto

```text
Hora: 2026-10-23 09:59:45 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1
--- FAIL: TestComandoShipEsperaCICompleto (0.00s)
    --- FAIL: TestComandoShipEsperaCICompleto/Every_required_job_must_finish_in_success_on_the_pushed_SHA. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: Every required job must finish in success on the pushed SHA.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.003s
FAIL
```

2-mutante-comando-omitido

```text
Hora: 2026-10-23 09:59:44 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1
--- FAIL: TestComandoShipEsperaCICompleto (0.00s)
    --- FAIL: TestComandoShipEsperaCICompleto/queued,_in_progress,_cancelled,_skipped,_neutral_and_timed_out_are_not_green. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: queued, in_progress, cancelled, skipped, neutral and timed_out are not green.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.003s
FAIL
```

2-mutante-comando-otro-sha

```text
Hora: 2026-10-23 09:59:44 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1
--- FAIL: TestComandoShipEsperaCICompleto (0.00s)
    --- FAIL: TestComandoShipEsperaCICompleto/A_result_for_another_SHA_is_not_green. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: A result for another SHA is not green.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.003s
FAIL
```

2-mutante-comando-parcial

```text
Hora: 2026-10-23 09:59:43 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1
--- FAIL: TestComandoShipEsperaCICompleto (0.00s)
    --- FAIL: TestComandoShipEsperaCICompleto/One_successful_job_with_another_pending_is_not_green. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: One successful job with another pending is not green.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.003s
FAIL
```

2-mutante-fragmento-cancelado

```text
Hora: 2026-10-23 09:59:40 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1
--- FAIL: TestShipEsperaCICompleto (0.00s)
    --- FAIL: TestShipEsperaCICompleto/queued,_in_progress,_cancelled,_skipped,_neutral_and_timed_out_are_not_green. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: queued, in_progress, cancelled, skipped, neutral and timed_out are not green.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
```

2-mutante-fragmento-cero-runs

```text
Hora: 2026-10-23 09:59:42 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1
--- FAIL: TestShipEsperaCICompleto (0.00s)
    --- FAIL: TestShipEsperaCICompleto/Zero_runs,_no_CI_or_a_failed_query_are_not_green;_report_the_state_and_leave_delivery_unconfirmed. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: Zero runs, no CI or a failed query are not green; report the state and leave delivery unconfirmed.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
```

2-mutante-fragmento-error-consulta

```text
Hora: 2026-10-23 09:59:42 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1
--- FAIL: TestShipEsperaCICompleto (0.00s)
    --- FAIL: TestShipEsperaCICompleto/Zero_runs,_no_CI_or_a_failed_query_are_not_green;_report_the_state_and_leave_delivery_unconfirmed. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: Zero runs, no CI or a failed query are not green; report the state and leave delivery unconfirmed.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
```

2-mutante-fragmento-job-incompleto

```text
Hora: 2026-10-23 09:59:43 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1
--- FAIL: TestShipEsperaCICompleto (0.00s)
    --- FAIL: TestShipEsperaCICompleto/Every_required_job_must_finish_in_success_on_the_pushed_SHA. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: Every required job must finish in success on the pushed SHA.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
```

2-mutante-fragmento-omitido

```text
Hora: 2026-10-23 09:59:41 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1
--- FAIL: TestShipEsperaCICompleto (0.00s)
    --- FAIL: TestShipEsperaCICompleto/queued,_in_progress,_cancelled,_skipped,_neutral_and_timed_out_are_not_green. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: queued, in_progress, cancelled, skipped, neutral and timed_out are not green.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
```

2-mutante-fragmento-otro-sha

```text
Hora: 2026-10-23 09:59:41 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1
--- FAIL: TestShipEsperaCICompleto (0.00s)
    --- FAIL: TestShipEsperaCICompleto/A_result_for_another_SHA_is_not_green. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: A result for another SHA is not green.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
```

2-mutante-fragmento-parcial

```text
Hora: 2026-10-23 09:59:40 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1
--- FAIL: TestShipEsperaCICompleto (0.00s)
    --- FAIL: TestShipEsperaCICompleto/One_successful_job_with_another_pending_is_not_green. (0.00s)
        ship_ci_test.go:29: falta contrato de CI: One successful job with another pending is not green.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
```

### Verde previo al commit

```text
Hora: 2026-10-23 09:59:15 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./internal/harness/config ./assets -run TestShipEsperaCICompleto|TestComandoShipEsperaCICompleto|TestCompose_ShipToggle -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.015s
ok  	github.com/Group-Active-IA/active-stack/assets	0.003s
```

```text
Hora: 2026-10-23 09:59:34 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./internal/harness/config -run ^TestCompose_AllToggles -update -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.007s
```

```text
Hora: 2026-10-23 09:59:34 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	0.004s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.553s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.400s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.639s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	0.011s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.489s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.217s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.032s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	0.024s
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.039s
```
