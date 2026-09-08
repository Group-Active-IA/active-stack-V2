## Role Specialists — Convene A Role, Not A Sub-Agent You Invented

By default every delegation produces a generic sub-agent briefed from scratch,
different each time. Two changes that touch the same layer get reviewed against
two different standards, and neither standard was ever written down.

A role is a **standing standard**: the same criterion applied every time that
layer is touched. What makes it a role and not a job title is that its
definition of done is **observable**. "Writes good UI" cannot be checked.
"Controls operable by keyboard with visible focus" can.

### The seven roles, and what done means for each

**Functional analyst.** Done when every requirement is written as a scenario with
a named actor, a trigger, and an outcome someone can observe, and each one has at
least one case that **fails today**. A requirement nobody can show to be false is
not a requirement, it is a wish.

**Architect.** Done when every boundary the change crosses is named with its
dependency direction, and each decision is recorded next to
**the alternative it rejected** and what reversing it would cost. No new cycle
in the dependency graph.

**Backend.** Done when every endpoint touched declares its **error contract** —
status code and body shape for the failure paths, not only the happy one — input
is validated at the boundary rather than at the caller, no secret appears in code
or logs, and the suite covering the touched packages runs green.

**Frontend.** Done when controls are **operable by keyboard with visible focus**,
every error is associated to its field programmatically and not only by position,
loading and empty states exist and can be reached on demand, and the layout
survives a 320px viewport with no horizontal scroll.

**DBA.** Done when every migration has a
**down that restores the previous shape**, applied and rolled back against a
copy of real data, and every query the change adds carries its plan with **no
full scan on a table that grows**. Each new index names the query it serves.

**DevOps.** Done when the deploy runs from a **clean clone** with the documented
commands, the **rollback was exercised** and not merely described, every secret
comes from the environment and never from the repository, and this change's
failure mode is visible in a log or a metric somebody actually reads.

**QA.** Done when there is at least one test that
**fails without the change and passes with it**, boundary and error paths are
covered and not only the happy path, and every bug reported carries reproduction
steps another person can run without asking a question.

### Who gets convened in each phase

Not everybody in every phase. Choose by **what the change touches**: a database
migration convenes no frontend, and a stylesheet change convenes no DBA.

| Phase | Always | Depending on what the change touches |
|---|---|---|
| explore | functional analyst, architect | nobody else — nothing is built yet, so there is nothing to review |
| propose | architect | backend, frontend, DBA, DevOps for the layers the change names; QA, to say how the result would be proven false |
| apply | QA | the specialists of the layers the diff actually touches, convened **as reviewers** once the implementation sub-agent has returned |
| verify | — | **this phase does not exist in this block yet.** Nothing is assigned here |
| archive | functional analyst | architect, and only if the change moved a boundary |

The review pass sits in `apply` because `apply` is the last phase this block can
actually convene. The cycle documented here is explore → propose → apply →
archive: there is no `/opsx:verify` command, no verify skill, and no verify row
in the model table. A floor parked in a phase nobody can reach looks like a floor
and runs never.

Convening all seven every time costs the same as convening nobody: it stops being
a signal, and by the third change it gets skimmed.

### The roles you did NOT convene go in the report

Every role left out gets **one line**, naming it and **the reason** it was left
out. Two rules about where that line goes:

It belongs **inside the typed report**, not beside it. File it as **JUDGMENT**,
because that is what it is — your call about what this change did not need — and
give it the **confidence** the format asks of any judgment. A ledger in flat
prose is not just a fourth thing to hand back next to the base summary and the
typed report: it is the untyped claim the reporting format already refuses.

```
JUDGMENT  (orchestrator)  Roles left out of this change.
  DBA — the diff touches no migration and no query.
  Frontend — no file under the UI tree.
  DevOps — deploy and configuration unchanged.
  Confidence: high — I read the whole diff before choosing.
```

A role that was not called **is a decision**, and decisions get written down.
Right now it is a silence, and a silence cannot be argued with. The day somebody
reads "did not convene the backend reviewer" on a change that touches login,
that change stops before it ships.

### The floor that is not delegated

When the change touches **authentication, payments or user data**, these
reviewers are convened **regardless of what the orchestrator judges necessary**:

- **backend** — where the check either exists or does not
- **QA** — the failure paths, which is where these bugs live
- **DBA** — whenever that data comes to rest, whatever the diff looks like
- **architect** — for authentication and payments, because who is allowed to do
  what is a boundary, not an implementation detail

This choice is **not the orchestrator's to make**. A floor that may be skipped
whenever the orchestrator feels confident is not a floor: it gets skipped
precisely on the change where confidence came cheapest.

**This floor does not replace the human gate.** If governance is active, Auth,
Billing, Security and Audit logs are CRITICAL there: no code gets written without
explicit human approval, and this floor is a weaker and differently worded
trigger sitting on top of the same domains. Convening four reviewers is not that
approval — **a reviewer sub-agent is not a human**. The floor decides who reads
the change; governance decides who may let it proceed. Where both apply, the
human gate is the one that blocks.

### The reviewer should not get edit tools — and today that is only a request

Ask the reviewer not to edit: a reviewer that can fix what it finds fixes it
quietly, and nobody ever learns the code was wrong. The report is the entire
reason it was convened.

Be honest about what backs that up. Right now this
**is an instruction, not a guarantee.** The launch pattern documented in this
block is `Agent({ description, model, prompt })` — three keys, and
**none of them is a tool list**. The restriction can only travel inside the
prompt, which is precisely the instruction it would like not to be, and a
reviewer that ignores it leaves no trace. So do not report a review as
tool-constrained, and do not treat a reviewer's silence about edits as proof it
made none: check the diff.

**What is missing** to make it real: the launch pattern has to accept a list of
allowed or denied tools, so the orchestrator can withhold Edit and Write instead
of asking for them to go unused. Until that key exists, this section is a request
with a reason attached, and it says so.
