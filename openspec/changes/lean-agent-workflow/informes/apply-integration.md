# Integration implementation report

## EVIDENCE

Author: integration sub-agent, Codex runtime. Declared family: shared with orchestrator; no independent-provider claim. Sources read: AGENTS.md, proposal/design/specs, catalogue, command installer tests and actual verify asset, README/ARCHITECTURE and lean composition test matrix. No other reviewer conclusions used; parent supplied compatibility boundaries. No auditor dictated exact cases. This is constructor/integration work, not an independent review.

2026-10-29, local -03, completed 10:58:27:

- Baseline command `go test ./internal/catalog ./internal/harness/command` first failed to write sandboxed Go cache; rerun with requested escalation passed: catalog 0.071s, command 0.008s. This was environment setup, not a RED assertion.
- RED after adding tests, same command: `TestDefaultOrchestratorSelectsLeanWorkflow` failed at lean_workflow_test.go:13 (`default orchestrator does not select lean-workflow`); `TestInstalledVerifySelectsWorkflowBeforeLegacyInstructions` failed at lean_verify_test.go:25 (`installed verify lacks explicit lean-first and legacy-only dispatch`). Both were assertion failures before changing catalogue/command asset.
- GREEN after changes, same command: `ok .../internal/catalog 0.081s`, `ok .../internal/harness/command 0.011s`.
- `gofmt -w internal/catalog/lean_workflow_test.go internal/harness/command/lean_verify_test.go` completed.
- `git diff --check` returned no output; `cmp AGENTS.md CLAUDE.md` returned no output (identical). Checked at 10:58:27 -03.

Changes: catalogue enables lean-workflow by default with bilingual description; role-definitions describes five workflow roles plus seven optional specialists without fixed models. Installed verify asset dispatches lean before legacy and explicitly prevents fallthrough; lean uses breaker then auditor, mandatory code TDD and current originals. Legacy branch remains under explicit opt-out. README/ARCHITECTURE describe the new default and historical compatibility. AGENTS/CLAUDE mirrors reflect the current source contract and preserve no-build/no-commit rules. New docs/lean-agent-workflow.md covers roles, TDD, risk, context/state, toggle matrix, plugins, KB/obsidian and evidence limits.

Tests exercise catalogue Load default data and the actual command installer output, including roles bound and legacy fallthrough prevention. No code build, user dotfile installation, commit or push performed. Report and tests are integration evidence; parent consolidates tasks.md and full-suite evidence.

## JUDGMENT

Follow-up integration (same run, 11:00 local approximately): existing config documentation guards exposed catalogue descriptions omitting historical names and root toggle enumerations lacking lean-workflow. Updated the actual descriptions/enumerations, preserving legacy senior-only-verify wording explicitly scoped to legacy. No config tests edited. `go test ./internal/harness/config -run 'TestRootDocs_ToggleEnumerationsAreNotStale|TestDocsSeniorSoloVerify|TestSeniorCatalogoDescribeVerify|TestCatalog_OrchestratorDescriptionIsNotStaleAboutToggles' -count=1` now passes (`ok .../internal/harness/config 0.016s`). Final catalogue/command run passes (`catalog 0.068s`, command cached). Matrix matches config constructor: optional governance/evidence formatting/KB/Engram/ship; required five-role sequence/TDD with senior folded into auditor. No measured savings claimed in docs.

Compared generated/default catalogue and command bindings against agreed lean design and config agent's reported matrix. Confidence high for these static/installation contracts. Mandatory agent behavior and economic savings require runtime evaluation; these tests do not prove agent obedience.

## ASSUMPTION / LIMITS

Report timestamps summarize this run; exact per-command wall timestamps were not separately emitted. No API/plugin execution or cache/billing measurement. Full-suite checks and final review remain parent responsibilities. The change is uncommitted, so no delivered-commit test claim is made. No mutation needed: new tests were observed red at their assertions before asset changes.
