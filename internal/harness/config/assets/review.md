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

### Three independent passes

Ask the rupture question as **three independent passes**, not as one.

Independent means exactly one thing:
**none of them sees the conclusions of the others** — not their findings, not
their verdicts, not a summary of either. Reviewers who read each other
**converge on the first opinion**, because the second finds the first one's
framing already convincing and spends its effort confirming it. Three reviewers
that converged that way **stop being three**: what comes back looks like
agreement and is repetition.

One part of that the orchestrator genuinely controls, because it writes each
prompt itself: withholding the other passes' output is a matter of what it
chooses to put in there. Two channels stay open behind it, and neither is closed
by a careful prompt.

**The change directory.** A pass **can read the change directory**, so if the
first verdict is already sitting in it, keeping that verdict out of the prompt
bought nothing. Launch all three and collect every result
**before writing any verdict** anywhere a later pass could open it.

**Shared memory, which is the larger one.** The sub-agent prompt this block
constructs carries `topic_key: "opsx/{change-name}/{phase}"`. Three passes on one
change in one phase land on **the same key** — and the Engram protocol, when it
is active, is mandatory rather than optional: it asks for a search before
starting work that may already have been done, and a save
**after every completed action**. So pass two can find pass one's conclusions
with both of them following the rules, and the ordering rule above does not help:
those writes do not happen when the verdict is written, they happen mid-pass.

**What is missing** to close it: **one key per pass** — a topic like
`opsx/{change-name}/verify/pass-{n}` — so the three do not share a topic, plus an
instruction not to search the sibling passes' keys. Until the phase key is split
that way, treat the independence of the three passes as **asked for and partly
achieved**, and say so in the verdict rather than reporting three passes as if
nothing connected them.

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

Give each pass one lens and record which pass got which. A pass told only to
"review it" reviews correctness, because correctness is what everybody reviews.

### It reports, it does not block

The verdict is written into the change as `review-verdict.md`, next to the other
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

If the three passes **found nothing, say so** — one line, naming the three lenses
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
