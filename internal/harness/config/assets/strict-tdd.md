# Strict TDD Module: Apply Phase

> **This module is loaded ONLY when Strict TDD Mode is enabled AND a test runner is available.**
> If you are reading this, the orchestrator already verified both conditions. Follow every instruction.

## TDD Philosophy

TDD is not testing. TDD is **software design driven by tests**. You write a test that describes what the code SHOULD do, then write the minimum code to make it real. The tests design the API, the contracts, the behavior. Code is a side effect of tests.

### The Three Laws

1. **Do NOT write production code** until you have a failing test
2. **Do NOT write more test** than is necessary to fail
3. **Do NOT write more code** than is necessary to pass the test

## TDD Implementation Cycle

For EVERY task assigned to you, follow this cycle strictly:

```
FOR EACH TASK:
├── 0. SAFETY NET (only if modifying existing files)
│   ├── Run existing tests for files being modified
│   ├── Capture baseline: "{N} tests passing"
│   ├── If any FAIL → STOP, report as "pre-existing failure"
│   │   (do NOT fix pre-existing failures: report to orchestrator)
│   └── This baseline proves you did not break what already worked
│
├── 1. UNDERSTAND
│   ├── Read the task description
│   ├── Read relevant spec scenarios (these ARE your acceptance criteria)
│   ├── Read the design decisions (these CONSTRAIN your approach)
│   ├── Read existing code and test patterns (match the style)
│   └── Determine test layer (see "Choosing Test Layer" below)
│
├── 2. RED: Write a failing test FIRST
│   ├── Write test(s) that describe the expected behavior from the spec
│   ├── Prefer pure functions where possible (no side effects = easy to test)
│   ├── The test describes required behavior through new or existing production code
│   └── GATE: Execute the test against the old behavior and observe its failure before GREEN
│
├── 3. GREEN: Write the MINIMUM code to pass
│   ├── Implement ONLY what the failing test needs
│   ├── Fake It is VALID here (hardcoded return values are OK)
│   ├── EXECUTE tests → must PASS
│   └── GATE: Do NOT proceed until GREEN is confirmed by execution
│
├── 4. TRIANGULATE (MANDATORY for most tasks)
│   ├── Add a second test case with DIFFERENT inputs/expected outputs
│   ├── EXECUTE tests → if Fake It breaks, generalize to real logic
│   ├── Repeat until ALL spec scenarios for this task are covered
│   ├── MINIMUM: at least 2 test cases per behavior (happy path + one edge case)
│   └── GATE: All spec scenarios for this task must have tests before REFACTOR
│
├── 5. REFACTOR: Improve without changing behavior
│   ├── Extract constants, functions; improve naming; remove duplication
│   ├── EXECUTE tests after EACH refactoring step → must STILL PASS
│   └── GATE: Tests green after EVERY refactoring change
│
├── 6. Mark task complete [x] only with executed RED and GREEN pasted in tasks.md
└── 7. Note any deviations or issues discovered
```

## Rules

The cycle above is binding (its GATEs, the Three Laws, the Safety Net, the evidence under each task). One rule it does NOT make obvious:

- NEVER write trivial assertions (tautologies, type-only checks, ghost loops): they pass without exercising behavior.

## Constructor brief and boundaries

The constructor receives the full apply and archive instructions in its prompt; phase 1 does not move them out of this block.

Before acting, read the numbered lessons when a knowledge base is available and record the path read or its absence in the apply report.

The brief names consumers outside the assigned area, including other modules, infrastructure and scripts (L-023).

Resolve wording or ordering conflicts and declare the resolution as your own (L-030).

A new rule after signature requires its own human approval before construction; stop rather than infer permission from the brief (L-029).

## Evidence under each task

For each task, paste RED and GREEN under that task in tasks.md: command, date +%T and the last output lines.

Mark [x] only after both executed outputs are pasted; a written test alone is not RED.

If a test was born green, make a temporary 2-4 line mutant that fails at the assertion the test claims to prove (L-019, L-034).

Paste the mutant diff, command, output and failing file:line; a setup or compilation failure does not prove that assertion.

Revert the mutant and rerun green; never commit the mutant.

Without a prior commit, copy the original file before mutating and declare that backup (L-013).

## Closing on the delivered commit

Run closing suites after the last commit and record git rev-parse HEAD beside command, time and output (L-031).

Any later code change invalidates that closing run, even if an earlier run was green.

Green means all required workflow steps, including formatting, lint, types and tests; a test-only run is not complete CI (L-005).

The person adding a file runs the applicable workflow checks; report unavailable or unauthorized checks explicitly.

## Corrections and model family

In verify rounds 2 and 3, the constructor fixes only the auditor's list; new scope goes to debt.

When another provider is available, the constructor may run in another model family, for example codex exec.

If no other provider is available, say so in the report and do not claim independence.

Declare the model/provider and shared family in each report; without a dispatch record, label identity as ASSUMPTION.

## Independent evidence and complete reports

The breaker and auditor count the tests themselves; do not accept the constructor's count as their evidence.

Declare which sources were read, whether the auditor dictated exact cases and any contamination; dictated-case agreement is not independent convergence.

With an external identity provider, request a real walkthrough with the person at the keyboard before archive; record unvisited paths as explicit debt and a limitation for the auditor, never as validated.

Preserve the complete apply report in openspec/changes/<change>/informes/ with phase and delivery in its name, separate from verify reports.

Index author, declared model/provider, SHA, provenance and missing reports in LEEME.md; preserve the full return or textual transcript and declare orchestration mediation.

If a complete report cannot be recovered, report it missing and stop closure.

Under each task, retain the executed evidence, for example:

```text
RED
command: <focused test command against the old behavior>
time: <date +%T>
output: <last lines, including assertion failure and file:line>
GREEN
command: <same focused test command after the change>
time: <date +%T>
output: <last lines, including passed result>
closing suite: <command, time, git rev-parse HEAD and last lines after commit>
```

The return summary points to these task entries and lists any missing check.
The old checkmark-only summary cannot replace executed evidence.
