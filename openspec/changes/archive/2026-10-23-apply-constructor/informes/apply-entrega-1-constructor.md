# Apply, entrega 1, constructor: apply-constructor

Autor: Codex. Proveedor declarado: OpenAI. Modelo/familia exactos sin registro de despacho: SUPUESTO. SHA de entrega: f3fc02aae836d9d458a41e60a5b8e4d20e944f90.
Procedencia: informe del constructor y salidas textuales de herramientas de esta sesión; copia manual mediada por el constructor, sin captura por hook. No se afirma independencia ni revisión de terceros. Se leyeron las reglas, propuesta y exploración; no hay KB local ni archivo de lecciones. No se recibieron casos dictados por el auditor en esta entrega.

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

## Suite posterior al commit de esta entrega

```text
Hora: 2026-10-23 10:03:44 -0300
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
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
