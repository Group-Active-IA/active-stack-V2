## Adversarial Review — Conformity First, Then What Breaks

`/opsx:verify` answers one question, and answers it well. This section is the
other question, the one nothing in the cycle asks today.

### Two questions, not one

**Conformity** — does what was built match the artifacts: the specs, the tasks,
the design? That is what `/opsx:verify` checks, across completeness, correctness
and coherence. Run it first and change nothing about it. A change that does not
do what it said it would do needs no cleverer review than that one.

**Rupture** — what breaks that nobody wrote down? Nobody tests what happens when
the external provider **returns garbage**: a 200 with an empty body, a field that
is suddenly a string, a response that never arrives. That is not in the specs,
and it is missing for a reason — if it were there, somebody would
**already have thought of it**, and the code would already handle it.

The two questions fail differently, so they are asked separately. Conformity
alone certifies that the change does what was imagined, and says nothing at all
about what was not.

### One breaker, three mandatory lenses

One breaker covers all three mandatory lenses in one report. This is not independent between lenses. Declare sources read, whether the auditor dictated exact cases, and contamination; dictated agreement is not independent convergence. These declarations are not proof of isolation: a reviewer can read the change directory and shared memory.

When Engram is active, use topic_key `opsx/{change}/verify/round-{n}/{role}` for every save, including after every completed action; do not search sibling roles or their conclusions. The same key for every role would leak conclusions. Round and role keys reduce contamination but do not enforce isolation.

### Different lenses, not three of the same

Three passes through one lens is redundancy: it finds
**the same kind of bug three times** and hands back three findings that are one.
Three passes through three lenses finds **three kinds**. The count was never the
point — the coverage is.

- **Correctness** — does the logic do what it claims on the inputs it will
  actually see: the boundary ones, the empty one, the one that arrives twice?
- **Wiring** — for the extracted function, the new branch, the handler that
  moved: **who calls it**, **and is that caller tested**? Logic that is correct
  and unreachable is a refactor that quietly deleted a feature, and it passes
  every test of the code that remains.
- **Truth of the comments** — the comments, docstrings and names around the
  changed lines: do they **still tell the truth**? A comment that stopped being
  true is worse than no comment, because it is read with authority and nobody
  goes back to check it.

Give the breaker all three lenses and record work and result for each. A breaker told only to
"review it" reviews correctness, because correctness is what everybody reviews.

### It reports, it does not block

The verdict is written into the change as `verify-round-{n}.md`, next to the other
artifacts, and is **archived with it**. Give it a name, not a place: a verdict
filed as "somewhere in the change" is one nothing can open later.

Write it in the typed categories, not in prose. What a pass ran goes under
**EVIDENCE** with its command and output; every finding is a **JUDGMENT** with
its author and its **confidence**; anything taken for granted is an ASSUMPTION.
This is the fourth thing in this block competing over what comes back, and an
untyped fourth format is the one the reporting rules already refuse.

**ASSUMPTION** (not verified from this repo): that the file survives because
`openspec archive` moves the whole change directory. That is read out of the
installed archive command, not exercised here. If archiving instead copies a
known artifact list, the verdict is dropped on archive and this section is
wrong — check it before relying on the verdict being there next month.

The review **does not block**. It holds no veto and does not need one. A reviewer
that cannot tell **a defect from a taste**, and halts the work for both, gets
**switched off in a week** — and then nothing gets reviewed at all, which is
strictly worse than a review somebody argued with.

**Closing over an open finding** is legitimate. Not every finding is worth what
it costs to fix now, and that call belongs to the person, not to the reviewer.
What is not legitimate is the finding going **invisible**: written down and
knowingly accepted is a decision, while never written down is an accident that
gets discovered later and called a surprise.

### An empty verdict is a finding

If all three lenses **found nothing, say so** — one line, naming the three lenses
and what each one actually covered.

That line is not a formality. Without it, a clean review and a review that never
ran leave identical evidence behind: nothing. Silence cannot tell
**a clean review from one that never ran** apart, and the second is by far the
more common of the two.

### The datum that discriminates

Before accepting or rejecting any finding, name it **in advance of looking** — the
single observation that
**separates the hypothesis from its alternative**, then go get that one.

Not "show me the logs": logs always contain something that fits whatever was
already believed. The form is: "**look for this line**: `retry budget exhausted`.
If it is present the failure is A; **if it is absent** it is B." A datum with only
one branch discriminates nothing, it confirms.

Then check that the evidence is **from this run**. Output from an earlier run
describes a tree that no longer exists, and it is the most persuasive wrong
answer on offer, because it looks exactly like the right one.

### The loop

Run conformity -> breaker -> senior (if active) -> auditor, in that order on a quiet diff. Wait for the senior report before launching the auditor. Only APROBADO from the auditor, with conformity satisfied and complete reports, enables archive. Unresolved CRITICAL prevents APROBADO. The breaker reports; the auditor decides; the orchestrator executes that verdict, never substitutes its own approval.

Before convening reviewers, write `## Bucle de verify` in the change's tasks.md: round 1, 2 or 3, in progress or closed, SHA or reviewed diff, report paths, blockers and decision. Resume reads tasks.md, not Engram as authority. Session recovery or failed capture does not increment or reset the counter. A round begins with verify on a quiet diff and closes with the auditor verdict.

Blockers in round 1 or 2 return only the auditor's list to apply. Increment only after correcting those blockers, at the next verify. In rounds 2 and 3 the breaker checks corrections, consumers and regressions outside the diff. New blocking findings go to the auditor; other findings become debt assigned to the first change that uses the piece. A new rule requires human signature before construction (L-029).

Round 3 closes the loop, including open blockers. No automatic fourth round. Present the open list to Juani for a written decision in tasks.md: correct in a newly approved scope, accept justified debt, or stop the change. Without a written answer nothing advances. Human acceptance neither resolves a CRITICAL nor fabricates APROBADO.

The breaker delivers these three sections even with no findings:
- Correctness: record boundary cases, command, output and result.
- Wiring: trace consumers and connections, command, output and result.
- Comments/documentation: compare contracts with behavior, sources and result.
Each failure has a concrete reproduction with command and output, and a failing test or a surviving mutant. The breaker and auditor count tests themselves, independently of the constructor's count. Before acting, read numbered lessons when KB is available; record the path read or its absence. With an external identity provider, request a real walkthrough with Juani at the keyboard before archive. Unvisited paths are explicit debt and a limitation for the auditor, never validated.

### Degraded operation

Announce missing toggles at verify start and record degradation in tasks.md. These are conditional rules, not mandatory catalog dependencies:
- Auditor absent with review active: no automatic approval or return to apply. Stop closure and ask Juani for a written decision in tasks.md to enable auditor or accept degraded scope. Do not impersonate an auditor.
- TDD absent: apply keeps its usual behavior; do not claim strict TDD or accredited red/green. Disclose absence to the auditor; verify still evaluates evidence and its loop.
- Senior absent: omit its report, announce this and let the auditor decide on conformity and breaker; do not invent historical context.
- Knowledge-base absent: do not promise automatic KB updates. Preserve LECCIONES and DEUDA in the change, with destination pending for Juani.
- Knowledge-base active with auditor absent: archive may update facts with evidence; do not invent auditor lessons or debt. Mark the source absent; closure needs Juani's written decision on degraded scope.
- Review absent: no breaker or three-round loop. Declare ordinary verify. If auditor is active, it decides on available reports; otherwise report conformity and request Juani's written decision, never approval of the complete cycle.
A missing required report stops closure, including an incomplete return: identify it, recover it or escalate to Juani; never replace it with a summary or assumed APROBADO. This is not an absent toggle.

### Complete reports

A full round has conformity, breaker, senior and auditor reports; omit only the senior when its toggle is off and declare other degradation as above. Keep the constructor's apply report separate. Store full return or textual transcript under openspec/changes/<change>/informes/ with phase, round and role in the filename. LEEME.md indexes author, declared model/provider, shared family, SHA, provenance, orchestration mediation and missing reports. Identity without dispatch records is ASSUMPTION. Reviewers receive no Write. verify-round-{n}.md references the originals and verdict, not replacements. A return that cannot be recovered in full stops closure. Reports declare sources, dictated cases and contamination.
