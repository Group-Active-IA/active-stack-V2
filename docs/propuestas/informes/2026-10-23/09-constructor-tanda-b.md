# Constructor: tanda B

Alcance: pasos 4, 5, 6, 7 y 7b de la propuesta firmada en 4748680.
Constructor: Codex, proveedor OpenAI; identidad y familia SUPUESTAS sin registro de despacho. No constituye revisión independiente ni aprobación de auditor.
Fuentes: CLAUDE.md, propuesta, fragmentos, tests, documentación y tanda A. KB numerada ausente; lecciones leídas en las referencias L-NNN de la propuesta. Casos dictados por la propuesta auditada; no se afirma convergencia independiente.
Consumidores: Compose de ocho variantes, comandos embebidos, catálogo y documentación. Pruebas textuales, sin acreditar conducta real de agentes. Sin instalación, build de binario ni push. Los informes previos del usuario quedan fuera de estos commits.

## verify-rompedor-auditor: rojo

### rojo

```text
Hora: 2026-10-23T10:07:01.725020-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./internal/harness/config ./assets -run "TestVerify|TestAuditorDecide|TestCompose_ReviewToggle_(ThreeIndependent|IndependenceIs|NamesTheShared)|TestCompose_SeniorToggle_Splits|TestCompose_AuditorToggle_HasASilence" -count=1
--- FAIL: TestCompose_AuditorToggle_HasASilenceThreshold (0.00s)
    auditor_toggle_test.go:115: the auditor fragment is missing "aprobado": without a silence threshold it becomes noise and gets switched off
--- FAIL: TestCompose_ReviewToggle_ThreeIndependentPasses (0.00s)
    review_toggle_test.go:145: the review fragment has no section starting with "### one breaker"
--- FAIL: TestCompose_ReviewToggle_IndependenceIsByConstructionNotByMechanism (0.00s)
    review_toggle_test.go:159: the review fragment has no section starting with "### one breaker"
--- FAIL: TestCompose_ReviewToggle_NamesTheSharedMemoryChannel (0.00s)
    review_toggle_test.go:300: the review fragment has no section starting with "### one breaker"
--- FAIL: TestCompose_SeniorToggle_SplitsLabellingFromHolding (0.00s)
    senior_toggle_test.go:413: the threshold section is missing "the senior supplies context; the auditor classifies and decides" — the split, in one line
    senior_toggle_test.go:413: the threshold section is missing "the orchestrator executes the auditor's verdict" — where the stop actually happens
--- FAIL: TestVerifyContratoBucle (0.00s)
    --- FAIL: TestVerifyContratoBucle/conformity_->_breaker_->_senior_(if_active)_->_auditor (0.00s)
        verify_loop_test.go:22: falta contrato: conformity -> breaker -> senior (if active) -> auditor
    --- FAIL: TestVerifyContratoBucle/only_aprobado_from_the_auditor,_with_conformity_satisfied_and_complete_reports,_enables_archive (0.00s)
        verify_loop_test.go:22: falta contrato: only aprobado from the auditor, with conformity satisfied and complete reports, enables archive
    --- FAIL: TestVerifyContratoBucle/unresolved_critical_prevents_aprobado (0.00s)
        verify_loop_test.go:22: falta contrato: unresolved critical prevents aprobado
    --- FAIL: TestVerifyContratoBucle/##_bucle_de_verify (0.00s)
        verify_loop_test.go:22: falta contrato: ## bucle de verify
    --- FAIL: TestVerifyContratoBucle/before_convening_reviewers (0.00s)
        verify_loop_test.go:22: falta contrato: before convening reviewers
    --- FAIL: TestVerifyContratoBucle/round_1,_2_or_3 (0.00s)
        verify_loop_test.go:22: falta contrato: round 1, 2 or 3
    --- FAIL: TestVerifyContratoBucle/in_progress_or_closed (0.00s)
        verify_loop_test.go:22: falta contrato: in progress or closed
    --- FAIL: TestVerifyContratoBucle/sha_or_reviewed_diff (0.00s)
        verify_loop_test.go:22: falta contrato: sha or reviewed diff
    --- FAIL: TestVerifyContratoBucle/report_paths,_blockers_and_decision (0.00s)
        verify_loop_test.go:22: falta contrato: report paths, blockers and decision
    --- FAIL: TestVerifyContratoBucle/resume_reads_tasks.md,_not_engram_as_authority (0.00s)
        verify_loop_test.go:22: falta contrato: resume reads tasks.md, not engram as authority
    --- FAIL: TestVerifyContratoBucle/failed_capture_does_not_increment_or_reset_the_counter (0.00s)
        verify_loop_test.go:22: falta contrato: failed capture does not increment or reset the counter
    --- FAIL: TestVerifyContratoBucle/blockers_in_round_1_or_2_return_only_the_auditor's_list_to_apply (0.00s)
        verify_loop_test.go:22: falta contrato: blockers in round 1 or 2 return only the auditor's list to apply
    --- FAIL: TestVerifyContratoBucle/round_3_closes_the_loop (0.00s)
        verify_loop_test.go:22: falta contrato: round 3 closes the loop
    --- FAIL: TestVerifyContratoBucle/no_automatic_fourth_round (0.00s)
        verify_loop_test.go:22: falta contrato: no automatic fourth round
    --- FAIL: TestVerifyContratoBucle/without_a_written_answer_nothing_advances (0.00s)
        verify_loop_test.go:22: falta contrato: without a written answer nothing advances
    --- FAIL: TestVerifyContratoBucle/a_new_rule_requires_human_signature_before_construction_(l-029) (0.00s)
        verify_loop_test.go:22: falta contrato: a new rule requires human signature before construction (l-029)
    --- FAIL: TestVerifyContratoBucle/correctness:_record_boundary_cases,_command,_output_and_result (0.00s)
        verify_loop_test.go:22: falta contrato: correctness: record boundary cases, command, output and result
    --- FAIL: TestVerifyContratoBucle/wiring:_trace_consumers_and_connections,_command,_output_and_result (0.00s)
        verify_loop_test.go:22: falta contrato: wiring: trace consumers and connections, command, output and result
    --- FAIL: TestVerifyContratoBucle/comments/documentation:_compare_contracts_with_behavior,_sources_and_result (0.00s)
        verify_loop_test.go:22: falta contrato: comments/documentation: compare contracts with behavior, sources and result
    --- FAIL: TestVerifyContratoBucle/a_failing_test_or_a_surviving_mutant (0.00s)
        verify_loop_test.go:22: falta contrato: a failing test or a surviving mutant
    --- FAIL: TestVerifyContratoBucle/outside_the_diff (0.00s)
        verify_loop_test.go:22: falta contrato: outside the diff
    --- FAIL: TestVerifyContratoBucle/first_change_that_uses_the_piece (0.00s)
        verify_loop_test.go:22: falta contrato: first change that uses the piece
    --- FAIL: TestVerifyContratoBucle/count_tests_themselves (0.00s)
        verify_loop_test.go:22: falta contrato: count tests themselves
    --- FAIL: TestVerifyContratoBucle/read_numbered_lessons (0.00s)
        verify_loop_test.go:22: falta contrato: read numbered lessons
    --- FAIL: TestVerifyContratoBucle/at_the_keyboard (0.00s)
        verify_loop_test.go:22: falta contrato: at the keyboard
    --- FAIL: TestVerifyContratoBucle/explicit_debt_and_a_limitation_for_the_auditor (0.00s)
        verify_loop_test.go:22: falta contrato: explicit debt and a limitation for the auditor
    --- FAIL: TestVerifyContratoBucle/announce_missing_toggles_at_verify_start_and_record_degradation_in_tasks.md (0.00s)
        verify_loop_test.go:22: falta contrato: announce missing toggles at verify start and record degradation in tasks.md
    --- FAIL: TestVerifyContratoBucle/auditor_absent_with_review_active (0.00s)
        verify_loop_test.go:22: falta contrato: auditor absent with review active
    --- FAIL: TestVerifyContratoBucle/no_automatic_approval_or_return_to_apply (0.00s)
        verify_loop_test.go:22: falta contrato: no automatic approval or return to apply
    --- FAIL: TestVerifyContratoBucle/tdd_absent (0.00s)
        verify_loop_test.go:22: falta contrato: tdd absent
    --- FAIL: TestVerifyContratoBucle/do_not_claim_strict_tdd (0.00s)
        verify_loop_test.go:22: falta contrato: do not claim strict tdd
    --- FAIL: TestVerifyContratoBucle/senior_absent (0.00s)
        verify_loop_test.go:22: falta contrato: senior absent
    --- FAIL: TestVerifyContratoBucle/do_not_invent_historical_context (0.00s)
        verify_loop_test.go:22: falta contrato: do not invent historical context
    --- FAIL: TestVerifyContratoBucle/knowledge-base_absent (0.00s)
        verify_loop_test.go:22: falta contrato: knowledge-base absent
    --- FAIL: TestVerifyContratoBucle/destination_pending (0.00s)
        verify_loop_test.go:22: falta contrato: destination pending
    --- FAIL: TestVerifyContratoBucle/review_absent (0.00s)
        verify_loop_test.go:22: falta contrato: review absent
    --- FAIL: TestVerifyContratoBucle/no_breaker_or_three-round_loop (0.00s)
        verify_loop_test.go:22: falta contrato: no breaker or three-round loop
    --- FAIL: TestVerifyContratoBucle/knowledge-base_active_with_auditor_absent (0.00s)
        verify_loop_test.go:22: falta contrato: knowledge-base active with auditor absent
    --- FAIL: TestVerifyContratoBucle/do_not_invent_auditor_lessons_or_debt (0.00s)
        verify_loop_test.go:22: falta contrato: do not invent auditor lessons or debt
    --- FAIL: TestVerifyContratoBucle/missing_required_report_stops_closure (0.00s)
        verify_loop_test.go:22: falta contrato: missing required report stops closure
    --- FAIL: TestVerifyContratoBucle/never_replace_it_with_a_summary (0.00s)
        verify_loop_test.go:22: falta contrato: never replace it with a summary
    --- FAIL: TestVerifyContratoBucle/informes/ (0.00s)
        verify_loop_test.go:22: falta contrato: informes/
    --- FAIL: TestVerifyContratoBucle/leeme.md (0.00s)
        verify_loop_test.go:22: falta contrato: leeme.md
    --- FAIL: TestVerifyContratoBucle/model/provider (0.00s)
        verify_loop_test.go:22: falta contrato: model/provider
    --- FAIL: TestVerifyContratoBucle/provenance (0.00s)
        verify_loop_test.go:22: falta contrato: provenance
    --- FAIL: TestVerifyContratoBucle/verify-round-{n}.md_references_the_originals (0.00s)
        verify_loop_test.go:22: falta contrato: verify-round-{n}.md references the originals
    --- FAIL: TestVerifyContratoBucle/reviewers_receive_no_write (0.00s)
        verify_loop_test.go:22: falta contrato: reviewers receive no write
    --- FAIL: TestVerifyContratoBucle/conformity,_breaker,_senior_and_auditor (0.00s)
        verify_loop_test.go:22: falta contrato: conformity, breaker, senior and auditor
    --- FAIL: TestVerifyContratoBucle/full_return_or_textual_transcript (0.00s)
        verify_loop_test.go:22: falta contrato: full return or textual transcript
--- FAIL: TestAuditorDecideSinAutoridadParalela (0.00s)
    verify_loop_test.go:28: falta decisión: aprobado
    verify_loop_test.go:28: falta decisión: bloqueantes
    verify_loop_test.go:28: falta decisión: deuda
    verify_loop_test.go:28: falta decisión: lecciones
    verify_loop_test.go:28: falta decisión: unresolved critical prevents aprobado
    verify_loop_test.go:28: falta decisión: classify disagreements with the senior
    verify_loop_test.go:28: falta decisión: dated decision
    verify_loop_test.go:28: falta decisión: count tests independently
    verify_loop_test.go:31: falta secuencia: after the breaker returns and before the auditor
    verify_loop_test.go:31: falta secuencia: the auditor classifies and decides
    verify_loop_test.go:31: falta secuencia: no parallel human wait
    verify_loop_test.go:32: autoridad vieja: the two run together
    verify_loop_test.go:32: autoridad vieja: when the senior returns a finding labelled blocking, do not start
    verify_loop_test.go:32: autoridad vieja: the senior labels, you hold
--- FAIL: TestVerifyDegradacionesComponibles (0.00s)
    --- FAIL: TestVerifyDegradacionesComponibles/#00 (0.00s)
        verify_loop_test.go:40: falta degradación visible
    --- FAIL: TestVerifyDegradacionesComponibles/auditor (0.00s)
        verify_loop_test.go:40: falta degradación visible
    --- FAIL: TestVerifyDegradacionesComponibles/tdd (0.00s)
        verify_loop_test.go:40: falta degradación visible
    --- FAIL: TestVerifyDegradacionesComponibles/senior (0.00s)
        verify_loop_test.go:40: falta degradación visible
    --- FAIL: TestVerifyDegradacionesComponibles/knowledge-base (0.00s)
        verify_loop_test.go:40: falta degradación visible
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.012s
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:111: falta encadenamiento: 9. **
    commands_test.go:111: falta encadenamiento: pending final decision
    commands_test.go:111: falta encadenamiento: only the auditor
    commands_test.go:111: falta encadenamiento: review fragment
    commands_test.go:111: falta encadenamiento: missing toggles
    commands_test.go:111: falta encadenamiento: written decision
    commands_test.go:111: falta encadenamiento: missing required report
    commands_test.go:112: conformidad anuncia archive antes del auditor
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.002s
FAIL
Salida: 1
```

## verify-rompedor-auditor: rojo-aserciones

### rojo-aserciones

```text
Hora: 2026-10-23T10:08:07.184803-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./internal/harness/config -run "TestCompose_ReviewToggle_(ThreeIndependent|IndependenceIs|NamesTheShared)" -count=1
--- FAIL: TestCompose_ReviewToggle_ThreeIndependentPasses (0.00s)
    review_toggle_test.go:146: falta contrato: one breaker covers all three mandatory lenses
    review_toggle_test.go:147: falta contrato: not independent between lenses
--- FAIL: TestCompose_ReviewToggle_IndependenceIsByConstructionNotByMechanism (0.00s)
    review_toggle_test.go:160: falta contrato: declare sources read
    review_toggle_test.go:161: falta contrato: contamination
    review_toggle_test.go:162: falta contrato: not proof of isolation
--- FAIL: TestCompose_ReviewToggle_NamesTheSharedMemoryChannel (0.00s)
    review_toggle_test.go:301: falta contrato: opsx/{change}/verify/round-{n}/{role}
    review_toggle_test.go:302: falta contrato: do not search sibling roles
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
Salida: 1
```

## verify-rompedor-auditor: verde-dirigido

### verde-dirigido

```text
Hora: 2026-10-23T10:08:07.627191-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./internal/harness/config ./assets -run "TestVerify|TestAuditor|TestCompose_Review|TestCompose_Senior|TestCompose_Auditor|TestNinguna" -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.071s
ok  	github.com/Group-Active-IA/active-stack/assets	0.003s
Salida: 0
```

Mutante temporal sin-wiring, copia previa en memoria desde archivo versionado: `Wiring: trace consumers and connections, command, output and result.` pasa a `Wiring: omitted.`.

## verify-rompedor-auditor: mutante-sin-wiring

### mutante-sin-wiring

```text
Hora: 2026-10-23T10:08:33.306398-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./internal/harness/config -run ^TestVerifyContratoBucle$ -count=1
--- FAIL: TestVerifyContratoBucle (0.00s)
    --- FAIL: TestVerifyContratoBucle/wiring:_trace_consumers_and_connections,_command,_output_and_result (0.00s)
        verify_loop_test.go:25: falta contrato: wiring: trace consumers and connections, command, output and result
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.005s
FAIL
Salida: 1
```

Mutante temporal sin-comentarios, copia previa en memoria desde archivo versionado: `Comments/documentation: compare contracts with behavior, sources and result.` pasa a `Comments/documentation: omitted.`.

## verify-rompedor-auditor: mutante-sin-comentarios

### mutante-sin-comentarios

```text
Hora: 2026-10-23T10:08:33.770934-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./internal/harness/config -run ^TestVerifyContratoBucle$ -count=1
--- FAIL: TestVerifyContratoBucle (0.00s)
    --- FAIL: TestVerifyContratoBucle/comments/documentation:_compare_contracts_with_behavior,_sources_and_result (0.00s)
        verify_loop_test.go:25: falta contrato: comments/documentation: compare contracts with behavior, sources and result
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
Salida: 1
```

Mutante temporal critical-aprueba, copia previa en memoria desde archivo versionado: `Unresolved CRITICAL prevents APROBADO.` pasa a `Unresolved CRITICAL allows APROBADO.`.

## verify-rompedor-auditor: mutante-critical-aprueba

### mutante-critical-aprueba

```text
Hora: 2026-10-23T10:08:34.218075-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./internal/harness/config -run ^TestVerifyContratoBucle$ -count=1
--- FAIL: TestVerifyContratoBucle (0.00s)
    --- FAIL: TestVerifyContratoBucle/unresolved_critical_prevents_aprobado (0.00s)
        verify_loop_test.go:25: falta contrato: unresolved critical prevents aprobado
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.005s
FAIL
Salida: 1
```

Mutante temporal archive-prematuro, copia previa en memoria desde archivo versionado: `Only APROBADO from the auditor, with conformity satisfied and complete reports, enables archive.` pasa a `Conformity enables archive before the auditor.`.

## verify-rompedor-auditor: mutante-archive-prematuro

### mutante-archive-prematuro

```text
Hora: 2026-10-23T10:08:34.732194-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./internal/harness/config -run ^TestVerifyContratoBucle$ -count=1
--- FAIL: TestVerifyContratoBucle (0.00s)
    --- FAIL: TestVerifyContratoBucle/only_aprobado_from_the_auditor,_with_conformity_satisfied_and_complete_reports,_enables_archive (0.00s)
        verify_loop_test.go:25: falta contrato: only aprobado from the auditor, with conformity satisfied and complete reports, enables archive
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
Salida: 1
```

Mutante temporal contador-reinicia, copia previa en memoria desde archivo versionado: `Resume reads tasks.md, not Engram as authority.` pasa a `Resume resets the counter from Engram.`.

## verify-rompedor-auditor: mutante-contador-reinicia

### mutante-contador-reinicia

```text
Hora: 2026-10-23T10:08:35.246360-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./internal/harness/config -run ^TestVerifyContratoBucle$ -count=1
--- FAIL: TestVerifyContratoBucle (0.00s)
    --- FAIL: TestVerifyContratoBucle/resume_reads_tasks.md,_not_engram_as_authority (0.00s)
        verify_loop_test.go:25: falta contrato: resume reads tasks.md, not engram as authority
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
Salida: 1
```

Mutante temporal cuarta-vuelta, copia previa en memoria desde archivo versionado: `No automatic fourth round.` pasa a `Allow an automatic fourth round.`.

## verify-rompedor-auditor: mutante-cuarta-vuelta

### mutante-cuarta-vuelta

```text
Hora: 2026-10-23T10:08:35.730292-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./internal/harness/config -run ^TestVerifyContratoBucle$ -count=1
--- FAIL: TestVerifyContratoBucle (0.00s)
    --- FAIL: TestVerifyContratoBucle/no_automatic_fourth_round (0.00s)
        verify_loop_test.go:25: falta contrato: no automatic fourth round
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
Salida: 1
```

## verify-rompedor-auditor: mutante-auditor-antes-senior

### mutante-auditor-antes-senior

```text
Hora: 2026-10-23T10:08:36.200119-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./internal/harness/config -run ^TestAuditorDecideSinAutoridadParalela$ -count=1
--- FAIL: TestAuditorDecideSinAutoridadParalela (0.00s)
    verify_loop_test.go:41: falta secuencia: after the breaker returns and before the auditor
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
Salida: 1
```

## verify-rompedor-auditor: goldens

### goldens

```text
Hora: 2026-10-23T10:08:36.689839-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./internal/harness/config -run "TestCompose_AllToggles" -update -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.006s
Salida: 0
```

## verify-rompedor-auditor: verde

### verde

```text
Hora: 2026-10-23T10:08:36.994921-03:00
SHA: f3fc02aae836d9d458a41e60a5b8e4d20e944f90
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	0.013s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.627s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.534s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.644s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	0.011s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.657s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	0.041s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.225s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	0.088s
ok  	github.com/Group-Active-IA/active-stack/internal/system	0.022s
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.035s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	0.029s
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.044s
Salida: 0
```

Revisión del diff de goldens de verify-rompedor-auditor: El diff reemplaza las tres pasadas por un rompedor con tres lentes, agrega contador y degradaciones, exige veredicto del auditor y ordena senior antes del auditor; conserva los umbrales del senior para el paso 5. Sólo se regeneran Claude y OpenCode por cambios de texto compuesto. No se modifican markers ni lógica de instalación.

## verify-rompedor-auditor: cierre-post-commit

### cierre-post-commit

```text
Hora: 2026-10-23T10:09:13.503862-03:00
SHA: db3de543858475a1cd8e574ce77b6f60933679f1
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
Salida: 0
```

## senior-solo-en-verify: rojo

### rojo

```text
Hora: 2026-10-23T10:09:54.390356-03:00
SHA: db3de543858475a1cd8e574ce77b6f60933679f1
Comando: go test ./internal/harness/config -run "TestCompose_SeniorToggle_(Preamble|Asks|Runs|IsNotAnEighth|Declares)|TestDocsSenior|TestSeniorCatalogo" -count=1
--- FAIL: TestCompose_SeniorToggle_PreambleSaysItIsEveryPhase (0.00s)
    senior_toggle_test.go:197: falta alcance de senior: only in **verify**
--- FAIL: TestCompose_SeniorToggle_AsksADifferentQuestionInEachPhase (0.00s)
    senior_toggle_test.go:341: falta alcance de senior: dated decisions
--- FAIL: TestCompose_SeniorToggle_RunsAtTheCloseOfEachPhase (0.00s)
    senior_toggle_test.go:356: falta alcance de senior: after the breaker returns and before the auditor
    senior_toggle_test.go:357: falta alcance de senior: quiet diff
--- FAIL: TestCompose_SeniorToggle_IsNotAnEighthRole (0.00s)
    senior_toggle_test.go:374: falta alcance de senior: once per verify round when active
--- FAIL: TestCompose_SeniorToggle_DeclaresItsCost (0.00s)
    senior_toggle_test.go:495: falta alcance de senior: one more sub-agent per verify round
    senior_toggle_test.go:496: falta alcance de senior: at most three per change
--- FAIL: TestSeniorCatalogoDescribeVerify (0.00s)
    senior_toggle_test.go:581: catálogo no describe senior sólo en verify
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.014s
FAIL
Salida: 1
```

## senior-solo-en-verify: rojo-docs

### rojo-docs

```text
Hora: 2026-10-23T10:10:25.765304-03:00
SHA: db3de543858475a1cd8e574ce77b6f60933679f1
Comando: go test ./internal/harness/config -run ^TestDocsSeniorSoloVerify$ -count=1
--- FAIL: TestDocsSeniorSoloVerify (0.00s)
    docs_toggles_test.go:268: README.md:197: falta senior sólo en verify
    docs_toggles_test.go:268: ARCHITECTURE.md:84: falta senior sólo en verify
    docs_toggles_test.go:268: AGENTS.md:86: falta senior sólo en verify
    docs_toggles_test.go:268: CLAUDE.md:92: falta senior sólo en verify
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
Salida: 1
```

## senior-solo-en-verify: verde-dirigido

### verde-dirigido

```text
Hora: 2026-10-23T10:10:26.224105-03:00
SHA: db3de543858475a1cd8e574ce77b6f60933679f1
Comando: go test ./internal/harness/config ./internal/catalog -run "TestCompose_Senior|TestDocsSenior|TestSeniorCatalogo|TestLoad" -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.015s
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	0.008s
Salida: 0
```

## senior-solo-en-verify: goldens

### goldens

```text
Hora: 2026-10-23T10:10:38.590568-03:00
SHA: db3de543858475a1cd8e574ce77b6f60933679f1
Comando: go test ./internal/harness/config -run TestCompose_AllToggles -update -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.005s
Salida: 0
```

## senior-solo-en-verify: verde

### verde

```text
Hora: 2026-10-23T10:10:38.850914-03:00
SHA: db3de543858475a1cd8e574ce77b6f60933679f1
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.533s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.422s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.630s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	0.090s
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.659s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.225s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.034s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.044s
Salida: 0
```

Revisión del diff de goldens de senior-solo-en-verify: El diff elimina la convocatoria en cinco fases, deja verify después del rompedor y antes del auditor, mantiene memoria y umbrales, y cambia el costo a uno por vuelta con tope de tres. Se alinean las descripciones senior y review en los cuatro documentos y el catálogo, validado mediante catalog.Load(). Sólo se regeneran Claude y OpenCode por cambios de texto compuesto. No se modifican markers ni lógica de instalación.

## senior-solo-en-verify: cierre-post-commit

### cierre-post-commit

```text
Hora: 2026-10-23T10:10:43.032451-03:00
SHA: 66ad39ec9b61364a69e518b749c14c023f9b7f14
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
Salida: 0
```

## archive-base-completa: rojo

### rojo

```text
Hora: 2026-10-23T10:11:13.207316-03:00
SHA: 66ad39ec9b61364a69e518b749c14c023f9b7f14
Comando: go test ./internal/harness/config -run "TestArchive|TestCompose_KnowledgeBaseToggle_TheDetector" -count=1
--- FAIL: TestCompose_KnowledgeBaseToggle_TheDetectorIsTheRewriteList (0.00s)
    knowledge_base_toggle_test.go:396: the detector section is missing "suspicious: append-only" — the detector itself
--- FAIL: TestArchiveBaseCompletaConEvidencia (0.00s)
    knowledge_base_toggle_test.go:720: falta contrato de archive: union of knowledge base impact and the actual diff
    knowledge_base_toggle_test.go:720: falta contrato de archive: even when the impact section exists
    knowledge_base_toggle_test.go:720: falta contrato de archive: lessons l-nnn
    knowledge_base_toggle_test.go:720: falta contrato de archive: only append lessons
    knowledge_base_toggle_test.go:720: falta contrato de archive: lecciones
    knowledge_base_toggle_test.go:720: falta contrato de archive: deuda
    knowledge_base_toggle_test.go:720: falta contrato de archive: each entry
    knowledge_base_toggle_test.go:720: falta contrato de archive: what changed
    knowledge_base_toggle_test.go:720: falta contrato de archive: sha, test or decision
    knowledge_base_toggle_test.go:720: falta contrato de archive: typed evidence
    knowledge_base_toggle_test.go:720: falta contrato de archive: auditor source absent
    knowledge_base_toggle_test.go:720: falta contrato de archive: written decision
    knowledge_base_toggle_test.go:720: falta contrato de archive: do not invent
    knowledge_base_toggle_test.go:736: prompt no entrega: actual diff
    knowledge_base_toggle_test.go:736: prompt no entrega: lessons l-nnn
    knowledge_base_toggle_test.go:736: prompt no entrega: lecciones
    knowledge_base_toggle_test.go:736: prompt no entrega: deuda
    knowledge_base_toggle_test.go:736: prompt no entrega: sha, test or decision
    knowledge_base_toggle_test.go:736: prompt no entrega: typed evidence
    knowledge_base_toggle_test.go:736: prompt no entrega: suspicious: append-only
--- FAIL: TestArchiveSonnetEnOchoVariantes (0.00s)
    --- FAIL: TestArchiveSonnetEnOchoVariantes/claude (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/opencode (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/codex (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/cursor (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/gemini (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/windsurf (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: less on archive
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/antigravity (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: less on archive
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/generic (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.007s
FAIL
Salida: 1
```

## archive-base-completa: verde-dirigido

### verde-dirigido

```text
Hora: 2026-10-23T10:11:47.574315-03:00
SHA: 66ad39ec9b61364a69e518b749c14c023f9b7f14
Comando: go test ./internal/harness/config -run "TestArchive|TestCompose_KnowledgeBase|TestOrchestrator" -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.023s
Salida: 0
```

## archive-base-completa: goldens

### goldens

```text
Hora: 2026-10-23T10:11:48.104631-03:00
SHA: 66ad39ec9b61364a69e518b749c14c023f9b7f14
Comando: go test ./internal/harness/config -run TestCompose_AllToggles -update -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.006s
Salida: 0
```

## archive-base-completa: verde

### verde

```text
Hora: 2026-10-23T10:11:48.387198-03:00
SHA: 66ad39ec9b61364a69e518b749c14c023f9b7f14
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.564s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.595s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.654s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.862s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.226s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.033s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.034s
Salida: 0
```

Revisión del diff de goldens de archive-base-completa: El diff cambia archive de haiku a sonnet con sustitución declarada, exige la unión de impacto y diff, agrega decisiones, lecciones L-NNN y deuda con evidencia, y etiqueta APPEND sin REWRITE como SUSPICIOUS: append-only. Conserva reporte por documento y arrastre PENDING. Sólo se regeneran Claude y OpenCode por cambios de texto compuesto. No se modifican markers ni lógica de instalación.

## archive-base-completa: cierre-post-commit

### cierre-post-commit

```text
Hora: 2026-10-23T10:11:52.445183-03:00
SHA: 84c570d86cb5622c659a975a4e1c02b10d81db8c
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
Salida: 0
```

## governance-aprobacion: rojo

### rojo

```text
Hora: 2026-10-23T10:12:17.443422-03:00
SHA: 84c570d86cb5622c659a975a4e1c02b10d81db8c
Comando: go test ./internal/harness/config -run ^TestGovernanceAprobacionHumanaEnProposal$ -count=1
--- FAIL: TestGovernanceAprobacionHumanaEnProposal (0.00s)
    governance_approval_test.go:17: falta contrato de firma: ## aprobación
    governance_approval_test.go:17: falta contrato de firma: proposal.md
    governance_approval_test.go:17: falta contrato de firma: | elección | persona | fecha | alcance |
    governance_approval_test.go:17: falta contrato de firma: no agent fills
    governance_approval_test.go:17: falta contrato de firma: written, dated and explicit scope
    governance_approval_test.go:17: falta contrato de firma: missing or ambiguous approval stops construction
    governance_approval_test.go:17: falta contrato de firma: new rule after signature
    governance_approval_test.go:17: falta contrato de firma: own row
    governance_approval_test.go:17: falta contrato de firma: before construction
    governance_approval_test.go:17: falta contrato de firma: l-029
    governance_approval_test.go:17: falta contrato de firma: cryptographic verification
    governance_approval_test.go:17: falta contrato de firma: not implemented
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.003s
FAIL
Salida: 1
```

## governance-aprobacion: verde-dirigido

### verde-dirigido

```text
Hora: 2026-10-23T10:12:17.842384-03:00
SHA: 84c570d86cb5622c659a975a4e1c02b10d81db8c
Comando: go test ./internal/harness/config -run ^TestGovernanceAprobacionHumanaEnProposal$ -count=1
--- FAIL: TestGovernanceAprobacionHumanaEnProposal (0.00s)
    governance_approval_test.go:17: falta contrato de firma: | elección | persona | fecha | alcance |
    governance_approval_test.go:17: falta contrato de firma: no agent fills
    governance_approval_test.go:17: falta contrato de firma: missing or ambiguous approval stops construction
    governance_approval_test.go:17: falta contrato de firma: new rule after signature
    governance_approval_test.go:17: falta contrato de firma: own row
    governance_approval_test.go:17: falta contrato de firma: before construction
    governance_approval_test.go:17: falta contrato de firma: l-029
    governance_approval_test.go:17: falta contrato de firma: cryptographic verification
    governance_approval_test.go:17: falta contrato de firma: not implemented
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.003s
FAIL
Salida: 1
```

## governance-aprobacion: goldens

### goldens

```text
Hora: 2026-10-23T10:12:18.326285-03:00
SHA: 84c570d86cb5622c659a975a4e1c02b10d81db8c
Comando: go test ./internal/harness/config -run TestCompose_AllToggles -update -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
Salida: 0
```

## governance-aprobacion: verde

### verde

```text
Hora: 2026-10-23T10:12:18.582404-03:00
SHA: 84c570d86cb5622c659a975a4e1c02b10d81db8c
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.517s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.333s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.599s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
--- FAIL: TestGovernanceAprobacionHumanaEnProposal (0.00s)
    governance_approval_test.go:17: falta contrato de firma: | elección | persona | fecha | alcance |
    governance_approval_test.go:17: falta contrato de firma: no agent fills
    governance_approval_test.go:17: falta contrato de firma: missing or ambiguous approval stops construction
    governance_approval_test.go:17: falta contrato de firma: new rule after signature
    governance_approval_test.go:17: falta contrato de firma: own row
    governance_approval_test.go:17: falta contrato de firma: before construction
    governance_approval_test.go:17: falta contrato de firma: l-029
    governance_approval_test.go:17: falta contrato de firma: cryptographic verification
    governance_approval_test.go:17: falta contrato de firma: not implemented
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	2.612s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.210s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.032s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.034s
FAIL
Salida: 1
```

Revisión del diff de goldens de governance-aprobacion: El diff agrega únicamente el artefacto de aprobación humana en proposal.md, tabla vacía con elección, persona, fecha y alcance, freno ante firma ausente o ambigua y fila propia previa para reglas nuevas L-029. No implementa verificación criptográfica ni completa la firma de nadie. Sólo se regeneran Claude y OpenCode por cambios de texto compuesto. No se modifican markers ni lógica de instalación.

## governance-aprobacion: cierre-post-commit

### cierre-post-commit

```text
Hora: 2026-10-23T10:12:22.226524-03:00
SHA: 7f467d534485318778bc17db4efcf98e182ded8f
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
--- FAIL: TestGovernanceAprobacionHumanaEnProposal (0.00s)
    governance_approval_test.go:17: falta contrato de firma: | elección | persona | fecha | alcance |
    governance_approval_test.go:17: falta contrato de firma: no agent fills
    governance_approval_test.go:17: falta contrato de firma: missing or ambiguous approval stops construction
    governance_approval_test.go:17: falta contrato de firma: new rule after signature
    governance_approval_test.go:17: falta contrato de firma: own row
    governance_approval_test.go:17: falta contrato de firma: before construction
    governance_approval_test.go:17: falta contrato de firma: l-029
    governance_approval_test.go:17: falta contrato de firma: cryptographic verification
    governance_approval_test.go:17: falta contrato de firma: not implemented
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	2.350s
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
FAIL
Salida: 1
```

## governance-aprobacion: rojo-aserciones-helper-corregido

### rojo-aserciones-helper-corregido

```text
Hora: 2026-10-23T10:12:43.792366-03:00
SHA: 7f467d534485318778bc17db4efcf98e182ded8f
Comando: go test ./internal/harness/config -run ^TestGovernanceAprobacionHumanaEnProposal$ -count=1
--- FAIL: TestGovernanceAprobacionHumanaEnProposal (0.00s)
    governance_approval_test.go:22: falta contrato de firma: ## aprobación
    governance_approval_test.go:22: falta contrato de firma: proposal.md
    governance_approval_test.go:22: falta contrato de firma: | elección | persona | fecha | alcance |
    governance_approval_test.go:22: falta contrato de firma: no agent fills
    governance_approval_test.go:22: falta contrato de firma: written, dated and explicit scope
    governance_approval_test.go:22: falta contrato de firma: missing or ambiguous approval stops construction
    governance_approval_test.go:22: falta contrato de firma: new rule after signature
    governance_approval_test.go:22: falta contrato de firma: own row
    governance_approval_test.go:22: falta contrato de firma: before construction
    governance_approval_test.go:22: falta contrato de firma: l-029
    governance_approval_test.go:22: falta contrato de firma: cryptographic verification
    governance_approval_test.go:22: falta contrato de firma: not implemented
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.003s
FAIL
Salida: 1
```

## governance-aprobacion: verde-corregido

### verde-corregido

```text
Hora: 2026-10-23T10:12:44.233810-03:00
SHA: 7f467d534485318778bc17db4efcf98e182ded8f
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
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.370s
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
Salida: 0
```

Corrección del constructor en paso 7: el helper de secciones cortaba en el encabezado del ejemplo Markdown. Se acotó por el inicio del fragmento completo, se repitió rojo contra governance del padre 84c570d y verde completo. Se enmienda el commit local para conservar uno por paso; se conservan los intentos fallidos como evidencia, no se los presenta como verde.
