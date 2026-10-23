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

### Cierre del paso 2, después del commit

```text
Hora: 2026-10-23 10:00:14 -0300
SHA: 60b5e83704d10cbbe6e14e34ed831aa726e273ff
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

## Paso 3: apply-constructor

El módulo strict-tdd ahora define el brief, lectura previa de lecciones si existe KB, consumidores externos, distinción entre redacción y regla nueva con firma propia, evidencia ejecutada por tarea, mutantes con copia previa y ubicación de aserción, cierre posterior al commit, workflow completo y límites de las vueltas 2 y 3. Declara proveedor/familia sin simular independencia. Retiene informes íntegros y fija conteo propio de tests por rompedor/auditor, contaminación declarada e identidad no recorrida como deuda explícita. Se mantiene la entrega íntegra de apply/archive por el prompt en fase 1; no se agregó hook ni control de proveedor.

Se corrigió la contradicción del gate viejo que pedía sólo test escrito: ahora exige ejecutarlo y verlo fallar contra la conducta anterior. La tabla de tildes se reemplaza por evidencia en tasks.md con comando, hora y últimas líneas. También se quita la referencia a esa tabla ya retirada. Los tests existentes TDDToggleOff y ElRepartoNoAutorizaRecortarLaRevision siguen pasando sin modificarlos.

Goldens: 85 líneas agregadas y 19 reemplazadas por variante Claude/OpenCode, sólo dentro de strict-tdd: nuevos contratos, retiro de la tabla de tildes y gate ejecutado. Se reemplaza la puntuación con em dash del fragmento por dos puntos según el pedido. No cambia ship ni ninguna otra fase o variante base.

### Rojo observado

```text
Hora: 2026-10-23 10:01:00 -0300
SHA: 60b5e83704d10cbbe6e14e34ed831aa726e273ff
Comando: go test ./internal/harness/config -run ^TestConstructorContratoApply$ -count=1
--- FAIL: TestConstructorContratoApply (0.00s)
    --- FAIL: TestConstructorContratoApply/The_constructor_receives_the_full_apply_and_archive_instructions_in_its_prompt;_phase_1_does_not_move_them_out_of_this_block. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: The constructor receives the full apply and archive instructions in its prompt; phase 1 does not move them out of this block.
    --- FAIL: TestConstructorContratoApply/Before_acting,_read_the_numbered_lessons_when_a_knowledge_base_is_available_and_record_the_path_read_or_its_absence_in_the_apply_report. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Before acting, read the numbered lessons when a knowledge base is available and record the path read or its absence in the apply report.
    --- FAIL: TestConstructorContratoApply/The_brief_names_consumers_outside_the_assigned_area,_including_other_modules,_infrastructure_and_scripts_(L-023). (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: The brief names consumers outside the assigned area, including other modules, infrastructure and scripts (L-023).
    --- FAIL: TestConstructorContratoApply/Resolve_wording_or_ordering_conflicts_and_declare_the_resolution_as_your_own_(L-030). (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Resolve wording or ordering conflicts and declare the resolution as your own (L-030).
    --- FAIL: TestConstructorContratoApply/A_new_rule_after_signature_requires_its_own_human_approval_before_construction;_stop_rather_than_infer_permission_from_the_brief_(L-029). (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: A new rule after signature requires its own human approval before construction; stop rather than infer permission from the brief (L-029).
    --- FAIL: TestConstructorContratoApply/For_each_task,_paste_RED_and_GREEN_under_that_task_in_tasks.md:_command,_date_+%T_and_the_last_output_lines. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: For each task, paste RED and GREEN under that task in tasks.md: command, date +%T and the last output lines.
    --- FAIL: TestConstructorContratoApply/Mark_[x]_only_after_both_executed_outputs_are_pasted;_a_written_test_alone_is_not_RED. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Mark [x] only after both executed outputs are pasted; a written test alone is not RED.
    --- FAIL: TestConstructorContratoApply/If_a_test_was_born_green,_make_a_temporary_2-4_line_mutant_that_fails_at_the_assertion_the_test_claims_to_prove_(L-019,_L-034). (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: If a test was born green, make a temporary 2-4 line mutant that fails at the assertion the test claims to prove (L-019, L-034).
    --- FAIL: TestConstructorContratoApply/Paste_the_mutant_diff,_command,_output_and_failing_file:line;_a_setup_or_compilation_failure_does_not_prove_that_assertion. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Paste the mutant diff, command, output and failing file:line; a setup or compilation failure does not prove that assertion.
    --- FAIL: TestConstructorContratoApply/Revert_the_mutant_and_rerun_green;_never_commit_the_mutant. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Revert the mutant and rerun green; never commit the mutant.
    --- FAIL: TestConstructorContratoApply/Without_a_prior_commit,_copy_the_original_file_before_mutating_and_declare_that_backup_(L-013). (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Without a prior commit, copy the original file before mutating and declare that backup (L-013).
    --- FAIL: TestConstructorContratoApply/Run_closing_suites_after_the_last_commit_and_record_git_rev-parse_HEAD_beside_command,_time_and_output_(L-031). (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Run closing suites after the last commit and record git rev-parse HEAD beside command, time and output (L-031).
    --- FAIL: TestConstructorContratoApply/Any_later_code_change_invalidates_that_closing_run,_even_if_an_earlier_run_was_green. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Any later code change invalidates that closing run, even if an earlier run was green.
    --- FAIL: TestConstructorContratoApply/Green_means_all_required_workflow_steps,_including_formatting,_lint,_types_and_tests;_a_test-only_run_is_not_complete_CI_(L-005). (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Green means all required workflow steps, including formatting, lint, types and tests; a test-only run is not complete CI (L-005).
    --- FAIL: TestConstructorContratoApply/The_person_adding_a_file_runs_the_applicable_workflow_checks;_report_unavailable_or_unauthorized_checks_explicitly. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: The person adding a file runs the applicable workflow checks; report unavailable or unauthorized checks explicitly.
    --- FAIL: TestConstructorContratoApply/In_verify_rounds_2_and_3,_the_constructor_fixes_only_the_auditor's_list;_new_scope_goes_to_debt. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: In verify rounds 2 and 3, the constructor fixes only the auditor's list; new scope goes to debt.
    --- FAIL: TestConstructorContratoApply/When_another_provider_is_available,_the_constructor_may_run_in_another_model_family,_for_example_codex_exec. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: When another provider is available, the constructor may run in another model family, for example codex exec.
    --- FAIL: TestConstructorContratoApply/If_no_other_provider_is_available,_say_so_in_the_report_and_do_not_claim_independence. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: If no other provider is available, say so in the report and do not claim independence.
    --- FAIL: TestConstructorContratoApply/Declare_the_model/provider_and_shared_family_in_each_report;_without_a_dispatch_record,_label_identity_as_ASSUMPTION. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Declare the model/provider and shared family in each report; without a dispatch record, label identity as ASSUMPTION.
    --- FAIL: TestConstructorContratoApply/The_breaker_and_auditor_count_the_tests_themselves;_do_not_accept_the_constructor's_count_as_their_evidence. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: The breaker and auditor count the tests themselves; do not accept the constructor's count as their evidence.
    --- FAIL: TestConstructorContratoApply/Declare_which_sources_were_read,_whether_the_auditor_dictated_exact_cases_and_any_contamination;_dictated-case_agreement_is_not_independent_convergence. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Declare which sources were read, whether the auditor dictated exact cases and any contamination; dictated-case agreement is not independent convergence.
    --- FAIL: TestConstructorContratoApply/With_an_external_identity_provider,_request_a_real_walkthrough_with_the_person_at_the_keyboard_before_archive;_record_unvisited_paths_as_explicit_debt_and_a_limitation_for_the_auditor,_never_as_validated. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: With an external identity provider, request a real walkthrough with the person at the keyboard before archive; record unvisited paths as explicit debt and a limitation for the auditor, never as validated.
    --- FAIL: TestConstructorContratoApply/Preserve_the_complete_apply_report_in_openspec/changes/<change>/informes/_with_phase_and_delivery_in_its_name,_separate_from_verify_reports. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Preserve the complete apply report in openspec/changes/<change>/informes/ with phase and delivery in its name, separate from verify reports.
    --- FAIL: TestConstructorContratoApply/Index_author,_declared_model/provider,_SHA,_provenance_and_missing_reports_in_LEEME.md;_preserve_the_full_return_or_textual_transcript_and_declare_orchestration_mediation. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: Index author, declared model/provider, SHA, provenance and missing reports in LEEME.md; preserve the full return or textual transcript and declare orchestration mediation.
    --- FAIL: TestConstructorContratoApply/If_a_complete_report_cannot_be_recovered,_report_it_missing_and_stop_closure. (0.00s)
        tdd_toggle_test.go:48: falta contrato del constructor: If a complete report cannot be recovered, report it missing and stop closure.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.005s
FAIL
```

```text
Hora: 2026-10-23 10:01:53 -0300
SHA: 60b5e83704d10cbbe6e14e34ed831aa726e273ff
Comando: go test ./internal/harness/config -run ^TestCompose_AllToggles -count=1
--- FAIL: TestCompose_AllToggles (0.00s)
    compose_test.go:59: Compose output does not match golden.
        got len=107396 want len=104023
--- FAIL: TestCompose_AllToggles_opencode (0.00s)
    compose_test.go:83: Compose output does not match golden.
        got len=107738 want len=104365
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.005s
FAIL
```

### Verde previo al commit

```text
Hora: 2026-10-23 10:01:52 -0300
SHA: 60b5e83704d10cbbe6e14e34ed831aa726e273ff
Comando: go test ./internal/harness/config -run TestConstructorContratoApply|TestCompose_TDDToggleOff|TestElRepartoNoAutorizaRecortarLaRevision -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.055s
```

```text
Hora: 2026-10-23 10:03:02 -0300
SHA: 60b5e83704d10cbbe6e14e34ed831aa726e273ff
Comando: go test ./internal/harness/config -run ^TestCompose_AllToggles -update -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.006s
```

```text
Hora: 2026-10-23 10:03:02 -0300
SHA: 60b5e83704d10cbbe6e14e34ed831aa726e273ff
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.565s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.436s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.665s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.621s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.232s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.034s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.036s
```

### Límites de la entrega

Esto es apply del constructor, sin veredicto de verify ni autorización de archive. Se dejan los cuatro changes abiertos. No se corrió go build por la regla del operador; no se afirma CI completo. La aceptación de reinstalación y restore de §3.6 queda pendiente del OK separado de Juani, fuera de esta tanda. El cierre posterior al commit de este último paso queda pegado en openspec/changes/apply-constructor/tasks.md e informes/apply-entrega-1-constructor.md, para no crear un segundo commit del paso sólo para referenciar su propio SHA.
