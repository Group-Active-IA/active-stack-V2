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
