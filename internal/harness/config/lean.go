package config

import "strings"

// composeLean keeps universal role contracts in the permanent context because
// not every adapter has a native agents directory. Optional integrations are
// concise and deterministic; they never append a historical lifecycle. Only
// Claude gets an extra native-dispatch note; every other variant is generic.
func composeLean(toggles []string, variantKey string) (string, error) {
	core, err := loadFragment("lean-workflow.md")
	if err != nil {
		return "", err
	}
	active := make(map[string]bool, len(toggles))
	for _, toggle := range toggles {
		active[toggle] = true
	}
	var out strings.Builder
	out.WriteString(core)
	if variantKey == "claude" {
		out.WriteString("\n\n")
		out.WriteString(leanClaudeDispatch)
	}
	for _, option := range []struct{ toggle, text string }{
		{"governance", leanGovernance}, {"evidence", leanEvidence},
		{"knowledge-base", leanKnowledge}, {"engram", leanMemory}, {"ship", leanShip},
	} {
		if active[option.toggle] {
			out.WriteString("\n\n")
			out.WriteString(option.text)
		}
	}
	return out.String(), nil
}

const leanClaudeDispatch = `## Claude native roles

Dispatch each role inside Claude, unless the project's Role executors assign it elsewhere, with the Agent tool and subagent_type: role-planner, role-plan-reviewer, role-constructor, role-breaker or role-auditor; the installed resource carries its procedure, so the brief adds only bounded context. Roles assigned to a plugin or another provider receive the self-contained brief instead.`

const leanGovernance = `## Human approval

Classify by risk and scope, not diff length. HIGH or CRITICAL construction requires explicit, dated human approval naming the choice and scope in proposal.md under ## Aprobación. Only the human fills that approval; no agent signs for them or treats reviewer approval as human authorization. New scope beyond approval needs a new human decision before construction. Authentication, billing, security, audit integrity and sensitive-data boundaries retain CRITICAL classification even for tiny diffs. Missing or ambiguous approval stops construction. Resolve wording conflicts explicitly, without inventing another human choice.`

const leanEvidence = `## Typed reports

EVIDENCE carries command, timestamp, reviewed revision and relevant output; preserve full output in an artifact when long. JUDGMENT names its author, comparison and confidence. ASSUMPTION states what was not verified and why. Unreproduced defects remain visible suspicions, not demonstrated blockers. Never promote a partial check into evidence. The original report is authoritative; summaries link to it and do not replace it. Deliver the auditor's verdict verbatim to the human.`

const leanKnowledge = `## Project knowledge

List the project's actual knowledge directory; do not assume a fixed set of filenames. Check per-document git author dates against newer archived changes, report named drift or unavailable dates and do not silently rewrite the map at session start. Planner records ## Knowledge Base Impact in proposal.md, naming each affected document and REWRITE or APPEND, or a justified none. At closure inspect the union of impact and actual diff: rewrite current facts/status, append dated decisions, numbered lessons and debt from the auditor, preserving history. Update the existing obsidian.md when project instructions designate it, including symlink targets, without creating duplicate notes or copying the whole note into every brief. Never include secrets. Keep technical canonical documentation distinct from personal notes. Write knowledge-base-report.md with REWRITTEN, APPENDED and PENDING by document, reason for no rewrites, and carry pending items from the previous report with their originating change/date. If that report or the KB is missing, declare the gap and preserve debt in the change; do not invent auditor lessons. Flag unjustified append-only work as SUSPICIOUS.`

const leanMemory = `## Persistent memory

When configured, recover relevant Engram context at session start and before resuming; do not search sibling reviewers' conclusions to manufacture independent agreement. Save significant decisions and completed actions with per-change, per-attempt, per-role keys such as opsx/{change}/{phase}/attempt-{n}/{role}. Keep summaries short and artifact paths explicit; OpenSpec and on-disk reports remain workflow authority. An empty search means nothing was saved, not that no decision existed. Update the designated project note for human continuity without duplicating full reports in every memory store. At session end record accomplishments, remaining work and relevant files. If unavailable, declare it and recover from artifacts. Never export cross-project memory; public memory export requires explicit human authorization.`

const leanShip = `## Publication

Publish only when explicitly requested. No remote: stop and report, never invent one. With uncommitted changes: ask, never decide alone. Before any push require archived/closed change, recorded auditor approval, a clean tree and the required checks green on the delivered revision. Diagnose branch, default remote branch, upstream and existing PR; fetch before comparing divergence. Failed lookups are unknown, not absence. Work directly on trunk requires a human decision before moving commits; pushing to trunk always asks. When publishing, recommend PR or direct push with reasons, and name rebase or merge when integrating. Update an existing PR rather than opening a duplicate; never merge or enable auto-merge. For a PR, check the source SHA, not the synthetic merge SHA. After push record its SHA and wait for every expected workflow/event and expanded matrix job on that SHA; inspect all jobs, including failures. Zero runs, failed queries, another SHA, pending, skipped, cancelled, neutral or timed-out jobs are not success. For Active Stack require ci.yml push tests on Ubuntu and Windows; derive required checks from workflows for other projects. A new push invalidates earlier CI evidence. Report incomplete delivery without claiming it ready.`
