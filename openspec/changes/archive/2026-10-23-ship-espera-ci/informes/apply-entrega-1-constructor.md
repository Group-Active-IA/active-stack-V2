# Apply, entrega 1, constructor: ship-espera-ci

Autor: Codex. Proveedor declarado: OpenAI. Modelo/familia exactos sin registro de despacho: SUPUESTO. SHA de entrega: 60b5e83704d10cbbe6e14e34ed831aa726e273ff.
Procedencia: informe del constructor y salidas textuales de herramientas de esta sesión; copia manual mediada por el constructor, sin captura por hook. No se afirma independencia ni revisión de terceros. Se leyeron las reglas, propuesta y exploración; no hay KB local ni archivo de lecciones. No se recibieron casos dictados por el auditor en esta entrega.

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


## Suite posterior al commit de esta entrega

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
