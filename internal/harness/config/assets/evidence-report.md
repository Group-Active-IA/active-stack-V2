## Typed Reporting — Separate What You Know From What You Assume

When you report back, the person reading cannot tell your categories apart unless
you name them. "The tests pass", "the error handling is consistent" and "the
billing integration works" arrive in the same tone, and only the first one is
backed by anything.

So label every claim you make. Three categories, and nothing else:

### EVIDENCE — you ran it

Requires **the command and its output**, with a **timestamp**.

Without the command that produced it, evidence is a memory. And the timestamp is
not decoration: output from two hours ago describes a tree that no longer exists,
and "the tests pass" about a previous state of the code is not a true statement
about this one.

```
EVIDENCE
  $ go test ./ventas/...
  ok  ventas  12 passed, 0 failed        (2026-09-08 14:32)
```

### JUDGMENT — you think so

Requires **who says it** and **what they compared it against**, plus a
**confidence** level.

A judgment with no confidence reads as a fact. Saying "high, I compared it
against three handlers in this repo" and saying "low, I do not know the traffic
this endpoint takes" are both useful; saying neither is what turns an opinion
into an apparent finding.

```
JUDGMENT  (apply agent)
  Error handling follows the repo's pattern.
  Compared against pedidos/handler.go and usuarios/handler.go.
  Confidence: high.
```

### ASSUMPTION — you took it for granted

Declare it **even if nobody asks**, and say what you did instead of verifying.

This is the category that does the work. What gets lost in a report is not what
the agent ignores: it is what the agent knew it was assuming and did not say.

```
ASSUMPTION  (not verified)
  Assumed the billing service answers JSON.
  Never called it: the 12 tests use a mock I wrote with the shape I guessed.
```

### The downgrade rule

A claim **never moves up** a category for convenience. Partially verified is an
assumption until it is fully verified. **When in doubt, the lower category.**

Without this rule everything ends up filed as evidence and the three sections
say nothing.

### Declaring an assumption is a merit

State it plainly, because it decides whether this works: **declaring an
assumption is a merit, not an admission of failure.** An assumption you name is
a bug that did not reach production. One you keep quiet is a bug with a delay
timer.

If declaring feels like confessing, the cheapest way out is to stop declaring,
and then the report looks cleaner while the project gets worse.

### A finding you did not reproduce is not a blocker

The downgrade rule above is about claims. It applies to **findings** in exactly
the same way, and findings are where it gets ignored: a reviewer suspects
something and files it as blocking, because blocking is the label that gets read.

So before reporting anything as blocking, **go and reproduce it**. If you cannot —
no tool for it, the environment will not come up, the datum is not there —
**report it anyway, as a suspicion, naming what you were missing to confirm it**.

The measurement, from a full cycle on a real project on 2026-09-09: of 27
findings reported as blocking, 12 were real. The other 15 were **paid for three
times** — the reviewer to find it, the auditor to review it, and the agent fixing
things to discard it.

**And the risk, out loud, because it is how this rule fails:** a reviewer
required to reproduce everything **goes quiet about what it could not prove**,
and some of those things are real. So the rule is **downgrade, never go silent**.
An unreproduced finding filed as a suspicion costs one line; the same finding
kept quiet costs a bug in production. Between the two,
**silence is the worse failure**, and it is the one that leaves no trace.

```
JUDGMENT  (correctness pass)  SUSPICION — not reproduced.
  The migration looks like it breaks on a password containing `%`.
  Missing to confirm it: I could not raise a database to run it against.
  Confidence: low.
```

### Practical notes

- Skip empty sections. A one-line change does not need three headings.
- Use the reader's language for the labels. The categories are the point, not
  the words.
- **When you launch a sub-agent, ask for its summary in this format.** The one
  who ran the command is the only one who can say whether something is evidence:
  if the summary arrives as prose, you would be classifying second-hand what you
  did not observe, and the assumption the sub-agent made silently stays silent.
