# Compose implementation report

## EVIDENCE

Sources read: AGENTS.md; lean-agent-workflow design/tasks/specs; compose.go, assets.go, install.go, types.go; historical config tests. No knowledge-base directory was supplied for this change; numbered lessons were not read.

- Baseline, 2026-10-29 10:54–10:55 -03: `GOCACHE=/tmp/active-stack-go-cache go test ./internal/harness/config` → `ok github.com/Group-Active-IA/active-stack/internal/harness/config 28.495s`. Initial default-cache attempt failed with read-only ~/.cache/go-build; cache was relocated to /tmp, not treated as a baseline code failure.
- RED, 2026-10-29 10:55 -03: `GOCACHE=/tmp/active-stack-go-cache go test ./internal/harness/config -run '^TestLeanWorkflow' -count=1` → FAIL at lean_workflow_test.go:14, `config compose: unknown toggle "lean-workflow"`; all five variant subtests reached this behavior error; no setup/compiler failure.
- GREEN after implementation, 2026-10-29 10:56 -03: same focused command → `ok .../internal/harness/config 0.004s`.
- Triangulation added optional toggle contracts, historical-toggle omissions, deterministic order/duplicates, unknown toggle rejection, generic fallback and legacy golden checks. Snapshot goldens were recorded after behavior tests, not claimed as independent RED.
- Refactor/check, 2026-10-29 10:59:24 -03: `GOCACHE=/tmp/active-stack-go-cache go test ./internal/harness/config -run 'TestLeanWorkflow|TestCompose_Legacy|TestVerifySinReviewTresAgentes|TestShipEsperaCICompleto|TestBloqueSeDirigeALaSesionPrincipal' -count=1` → `ok .../internal/harness/config 0.076s`.
- `wc -lc internal/harness/config/testdata/compose_lean_claude.golden internal/harness/config/testdata/compose_all_toggles_claude.golden` → lean 149 lines / 12944 bytes; legacy 2050 lines / 114880 bytes. This measures generated size, not actual tokens, cache or money.
- Full config rerun before integration settled: `go test ./internal/harness/config -count=1` with temporary GOCACHE → remaining failures confined to docs/catalog toggle enumerations and senior catalogue wording, delivered to integration owner. No failing lean or legacy behavioral checks in that run.

## Changes

`Compose` recognizes lean-workflow, validates all options before selecting it, then delegates to compact deterministic composition. Legacy path remains unchanged except excluding the new mode name from legacy missing-toggle diagnostics so historical output stays stable. Unknown-toggle diagnostics now list canonical names dynamically.

Universal core defines planner -> plan-reviewer -> constructor -> breaker -> auditor, provider-neutral dispatch, risk+scope classification, original reports, bounded briefs, durable tasks.md state, correction ownership and downstream reverification, three-failure escalation, unconditional code TDD and honest docs-only validation.

Optional contracts: governance, evidence format, project KB/obsidian, Engram and publication. Historical attendance toggles cannot disable the lean sequence or TDD; senior context is absorbed by auditor. Old extensive fragments remain available when lean-workflow is absent.

Historical tests now explicitly use legacyToggles(), with existing goldens byte-for-byte preserved. Four lean goldens and dedicated behavior tests cover compact output. No install, backup, permissions, filemerge, pipeline, catalogue or roledefs code edited by this agent.

## JUDGMENT

Implementation agent: confidence high in pure composition/compatibility covered by executed checks; no claim of runtime agent obedience or savings. Parent integration must rerun complete suites after docs/catalog settle. Context is universal across variants deliberately, because native role resources do not exist in every adapter.

## ASSUMPTIONS / limitations

No user configuration was reinstalled. No plugin or provider execution tested. No build or commit. Golden and text-contract tests verify generated instructions; they cannot establish agents follow them. Model/provider identity: inherited Codex execution, not an independent model family. Parent provided exact requested scenarios; no auditor dictated cases, and prior mapping exploration was used, so this is construction evidence rather than independent review.

## Final integration GREEN

EVIDENCE — 2026-10-29 10:59:58 -03 start, completed about 11:00:35: `GOCACHE=/tmp/active-stack-go-cache go test ./internal/harness/config -count=1` → `ok github.com/Group-Active-IA/active-stack/internal/harness/config 36.904s`. Integration docs/catalog fixes had landed before this run. `git diff --check -- internal/harness/config` and `gofmt -l` for new/edited Go implementation and tests produced no output. Existing two legacy golden files have no diff. No later code changes made by this agent.
