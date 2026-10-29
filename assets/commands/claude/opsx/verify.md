---
name: "OPSX: Verify"
description: Verify implementation matches change artifacts before archiving
category: Workflow
tags: [workflow, verify, experimental]
---

Verify that an implementation matches the change artifacts (specs, tasks, design).

**Input**: Optionally specify a change name after `/opsx:verify` (e.g., `/opsx:verify add-auth`). If omitted, check if it can be inferred from conversation context. If vague or ambiguous you MUST prompt for available changes.

## Lean workflow

First read the installed orchestrator instructions to determine whether `lean-workflow` is active. Do not infer the mode from available files alone. If active, use the compact contract and the steps below; **Do not continue into the legacy steps**.

1. Resolve the change using `openspec list --json`, `openspec status --change "<name>" --json`, and `openspec instructions apply --change "<name>" --json`. Recover its recorded revision, attempt, reports, blockers and next action. Ask only if the change is ambiguous.
2. On a quiet diff, check completeness, correctness and coherence against proposal, design, tasks and scenarios. Preserve conformity evidence in the change; this does not approve closure.
3. Read the change level and the project's Role executors table: the breaker and auditor run on the declared executor, or as native `role-*` subagents when no row exists. Apply the depth table: LOW runs one breaker round unless a failure is reproduced; MEDIUM and above repeat rounds until clean, within three attempts. Dispatch `role-breaker` with the named diff, affected consumers, constraints and existing evidence. Load its native role resource when available, otherwise use the universal contract. It tests correctness boundaries, wiring and documentation truth, reproduces findings, records commands/output/time and distinguishes suspicions from demonstrated failures. Reproduced code defects return to constructor through executed RED then GREEN; repeat invalidated checks and review on the corrected revision.
4. After breaker returns, dispatch `role-auditor` with all original required reports and the resulting revision. It checks conformity, test counts, evidence, assumptions, coverage and unresolved findings. No automatic senior stage or specialist fan-out is added. Audit findings route to planner or constructor according to ownership; rerun the affected downstream verification.
5. Persist full reports and state before transitions: originals at `openspec/changes/<change>/informes/verify-<attempt>-breaker.md` and `verify-<attempt>-auditor.md`, indexed in `LEEME.md` with executor/provider and revision. Missing execution capability, missing/stale required reports or open blockers stop closure; never fabricate approval. Require human direction after three unsuccessful attempts. HIGH/CRITICAL construction retains explicit human approval. TDD is mandatory for every code change even if the historical `tdd` toggle is absent; documentation-only validation must not claim code RED/GREEN. An absent `senior`, `review` or `auditor` legacy toggle cannot weaken the lean role sequence.
6. Report evidence and limitations. Only an auditor verdict with complete current evidence enables closure; update the project-declared KB/`obsidian.md` destination without creating duplicate notes. Do not install configuration, commit, push or archive merely because conformity passed.

## Legacy workflow

Run the following steps only when `lean-workflow` is explicitly disabled. They preserve the earlier toggle-driven procedure and its degradation rules.

**Steps**

1. **If no change name provided, prompt for selection**

   Run `openspec list --json` to get available changes. Use the **AskUserQuestion tool** to let the user select.

   Show changes that have implementation tasks (tasks artifact exists).
   Include the schema used for each change if available.
   Mark changes with incomplete tasks as "(In Progress)".

   **IMPORTANT**: Do NOT guess or auto-select a change. Always let the user choose.

2. **Check status to understand the schema**
   ```bash
   openspec status --change "<name>" --json
   ```
   Parse the JSON to understand:
   - `schemaName`: The workflow being used (e.g., "spec-driven")
   - Which artifacts exist for this change

3. **Get the change directory and load artifacts**

   ```bash
   openspec instructions apply --change "<name>" --json
   ```

   This returns the change directory and `contextFiles` (artifact ID -> array of concrete file paths). Read all available artifacts from `contextFiles`.

4. **Initialize verification report structure**

   Create a report structure with three dimensions:
   - **Completeness**: Track tasks and spec coverage
   - **Correctness**: Track requirement implementation and scenario coverage
   - **Coherence**: Track design adherence and pattern consistency

   Each dimension can have CRITICAL, WARNING, or SUGGESTION issues.

5. **Verify Completeness**

   **Task Completion**:
   - If `contextFiles.tasks` exists, read every file path in it
   - Parse checkboxes: `- [ ]` (incomplete) vs `- [x]` (complete)
   - Count complete vs total tasks
   - If incomplete tasks exist:
     - Add CRITICAL issue for each incomplete task
     - Recommendation: "Complete task: <description>" or "Mark as done if already implemented"

   **Spec Coverage**:
   - If delta specs exist in `openspec/changes/<name>/specs/`:
     - Extract all requirements (marked with "### Requirement:")
     - For each requirement:
       - Search codebase for keywords related to the requirement
       - Assess if implementation likely exists
     - If requirements appear unimplemented:
       - Add CRITICAL issue: "Requirement not found: <requirement name>"
       - Recommendation: "Implement requirement X: <description>"

6. **Verify Correctness**

   **Requirement Implementation Mapping**:
   - For each requirement from delta specs:
     - Search codebase for implementation evidence
     - If found, note file paths and line ranges
     - Assess if implementation matches requirement intent
     - If divergence detected:
       - Add WARNING: "Implementation may diverge from spec: <details>"
       - Recommendation: "Review <file>:<lines> against requirement X"

   **Scenario Coverage**:
   - For each scenario in delta specs (marked with "#### Scenario:"):
     - Check if conditions are handled in code
     - Check if tests exist covering the scenario
     - If scenario appears uncovered:
       - Add WARNING: "Scenario not covered: <scenario name>"
       - Recommendation: "Add test or implementation for scenario: <description>"

7. **Verify Coherence**

   **Design Adherence**:
   - If `contextFiles.design` exists:
     - Extract key decisions (look for sections like "Decision:", "Approach:", "Architecture:")
     - Verify implementation follows those decisions
     - If contradiction detected:
       - Add WARNING: "Design decision not followed: <decision>"
       - Recommendation: "Update implementation or revise design.md to match reality"
   - If no design.md: Skip design adherence check, note "No design.md to verify against"

   **Code Pattern Consistency**:
   - Review new code for consistency with project patterns
   - Check file naming, directory structure, coding style
   - If significant deviations found:
     - Add SUGGESTION: "Code pattern deviation: <details>"
     - Recommendation: "Consider following project pattern: <example>"

8. **Generate Verification Report**

   **Summary Scorecard**:
   ```
   ## Verification Report: <change-name>

   ### Summary
   | Dimension    | Status           |
   |--------------|------------------|
   | Completeness | X/Y tasks, N reqs|
   | Correctness  | M/N reqs covered |
   | Coherence    | Followed/Issues  |
   ```

   **Issues by Priority**:

   1. **CRITICAL** (Must fix before archive):
      - Incomplete tasks
      - Missing requirement implementations
      - Each with specific, actionable recommendation

   2. **WARNING** (Should fix):
      - Spec/design divergences
      - Missing scenario coverage
      - Each with specific recommendation

   3. **SUGGESTION** (Nice to fix):
      - Pattern inconsistencies
      - Minor improvements
      - Each with specific recommendation

   **Final Assessment**:
   - If CRITICAL issues: "X critical issue(s) found. Fix before archiving."
   - If only warnings: "No critical issues. Y warning(s) to consider. Conformity evaluated, pending final decision (with noted improvements)."
   - If all clear: "Conformity evaluated, pending final decision."

9. **Chain the final decision**

   This report evaluates conformity, not approval of the full cycle. Follow the review fragment when active: conformity, breaker, senior (if active), then auditor after the senior returns. Only the auditor enables archive in the complete cycle, with conformity satisfied and complete reports; unresolved CRITICAL prevents APROBADO. Follow the loop and persisted counter in that fragment, not a separate loop here.

   Announce missing toggles and record degradation in tasks.md. If review is absent there is no breaker or three-round loop; an active auditor decides on available reports. If auditor is absent, request the human owner's written decision to enable it or explicitly accept degraded scope, with no automatic approval or return to apply. TDD absent means no strict TDD claim; senior absent means no invented historical context; knowledge-base absent leaves LECCIONES/DEUDA in the change with destination pending. KB active with auditor absent may update facts but must mark auditor sources absent and request that written decision. A missing required report stops closure until recovered or escalated, never replaced by a summary.

**Verification Heuristics**

- **Completeness**: Focus on objective checklist items (checkboxes, requirements list)
- **Correctness**: Use keyword search, file path analysis, reasonable inference - don't require perfect certainty
- **Coherence**: Look for glaring inconsistencies, don't nitpick style
- **False Positives**: When uncertain, prefer SUGGESTION over WARNING, WARNING over CRITICAL
- **Actionability**: Every issue must have a specific recommendation with file/line references where applicable

**Graceful Degradation**

- If only tasks.md exists: verify task completion only, skip spec/design checks
- If tasks + specs exist: verify completeness and correctness, skip design
- If full artifacts: verify all three dimensions
- Always note which checks were skipped and why

**Output Format**

Use clear markdown with:
- Table for summary scorecard
- Grouped lists for issues (CRITICAL/WARNING/SUGGESTION)
- Code references in format: `file.ts:123`
- Specific, actionable recommendations
- No vague suggestions like "consider reviewing"
