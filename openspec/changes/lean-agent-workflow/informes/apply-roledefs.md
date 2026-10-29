# Apply: workflow role definitions

## Scope and provenance

Author: roles_impl subagent, same model family/provider as orchestrator (inherited dispatch; model identity not independently queried). Sources: AGENTS.md, roledefs resources/install/health/docs and tests/assets, lean-agent-workflow design and specs. No other reviewer conclusions read. Orchestrator supplied role requirements and later requested auditor executable checks, not exact test cases. Shared filesystem/history is potential contamination; no claim of provider independence or tool-enforced read-only reviewers.

## EVIDENCE

2026-10-29, local timezone -03.

Baseline: `GOCACHE=/tmp/active-stack-roledefs-go-cache go test ./internal/harness/roledefs` passed, output `ok github.com/Group-Active-IA/active-stack/internal/harness/roledefs 0.024s` before edits. Default cache failed read-only at ~/.cache/go-build; used task-specific /tmp cache. No network fallback/install.

RED 10:54:53: `GOCACHE=/tmp/active-stack-roledefs-go-cache go test ./internal/harness/roledefs -run 'TestDefaultResources_NoFixedModel|TestWorkflowResources_Contracts'` failed as expected:
- resources_test.go:121: seven specialists each fix model instead of caller selection.
- resources_test.go:142: five workflow resources absent.
Complete short output: /tmp/roles-red.log. Earlier full suite also failed TestInstall_AbsentCreates (7 vs 12), demonstrating actual installer resource delivery, not only prose matching.

GREEN 10:56:18: same package suite passed after resources and validation change: `ok .../internal/harness/roledefs 0.030s`.
Independent constructor count 10:56:48: `GOCACHE=/tmp/active-stack-roledefs-go-cache go test -count=1 -json ./internal/harness/roledefs > /tmp/roles-green.json`; parsed pass events: 29 tests, package pass 0.031s. Count is constructor evidence, not independent reviewer count.

Additional RED 10:56:58: `GOCACHE=/tmp/active-stack-roledefs-go-cache go test ./internal/harness/roledefs -run TestWorkflowResources_Contracts` failed at resources_test.go:151 because auditor lacked `Reproducí` and `Ejecutá comprobaciones` contract. Added requirements for independent checks and reproducing blocker claims, without redundant full-suite reruns.

Final GREEN 10:57:06: `GOCACHE=/tmp/active-stack-roledefs-go-cache go test -count=1 ./internal/harness/roledefs` => `ok .../internal/harness/roledefs 0.031s`. `gofmt` applied to changed Go tests/resources. No production changes after this final package run. No commit, build or user installation executed. No mutant necessary: new tests observed genuine assertion RED against prior behavior.

## Changes

Five native resources: role-planner.md, role-plan-reviewer.md, role-constructor.md, role-breaker.md, role-auditor.md. Seven specialists preserved, model frontmatter removed. Exact validation now requires twelve resources; senior remains rejected and duplicate auditor rejected. Installer safety algorithms unchanged; count test now verifies twelve created outcomes. Package comments updated to actual inventory.

Planner creates proposal/design/tasks/scenarios and classifies risk; reviewer questions scenarios/viability, does not claim tests of absent code; constructor requires baseline/assertion RED/minimal GREEN/triangulation/refactor and records task evidence; breaker covers correctness/wiring/documentation with reproducible cases; auditor reads full originals, checks artifacts/code/tests and independently executes pertinent checks before verdict. All use bounded brief, source declarations, timestamped evidence, judgment/confidence, assumptions and contamination reporting. No role hardcodes model or provider.

## JUDGMENT

Confidence high against design/specs and targeted installer/resource tests. Existing adapters resolve delivery unchanged; universal contracts for agents without native role files are outside my scope (compose owner). Specialist presence is compatibility inventory, not mandatory attendance. Frontmatter tools omit Write/Edit in reviewers but Bash can write; restriction is instruction, not sandbox enforcement.

## ASSUMPTIONS / limits

No live agent execution, cost/latency measurement or plugin transport test. Content tests prove shipped text, not agent compliance. No full repository checks or closing delivered-commit test because scope was roledefs and no commit requested. Parent consolidates task evidence and final verification; tasks.md not edited to avoid races. No user dotfiles changed.
