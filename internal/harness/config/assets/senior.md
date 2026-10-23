## The Senior — A Companion, Not A Phase

Every other role in this block is attached to a point in the cycle: the explorer
explores, the reviewer reads a diff, the archivist archives. This one is
attached to none of them. It is convened in **every phase**, and what it brings
is not a skill the others lack — it is **everything that happened before**.

### It is not a fifth reviewer

Run it beside the adversarial passes and the objection is fair: one more launch
on a phase that already convenes several, buying what?

The difference is not depth. It is what each one carries.

**The adversarial passes have focus and no memory.** A pass sees a diff, asks one
question of it, and terminates. Each one **starts clean**, and that is precisely
what makes it independent — a reviewer arriving with the last three changes in
its head arrives with the last three opinions too.

**The senior has memory and no focus.** It saw the earlier phases of this change,
and the changes that closed before it. Its value is not finding the bug in the
diff: three sharper readers are already doing that, through three lenses, and it
would be the fourth pass through a lens already covered. Its value is saying
**this contradicts a decision taken two weeks ago** — a sentence that is
unavailable to all three of them, because **none of them was there**.

So in `verify` the senior **does not hunt for bugs. It judges the findings** the
others returned. Whether a reported defect is a real problem in this project, or
a rule that was argued over and deliberately set aside, is a question the pass
that raised it **has no context to answer**.

**That is not the auditor's job.** Run the senior after the breaker returns and before the auditor. There is no parallel human wait. The auditor then reads the senior report and decides. The auditor reads
the same reports and asks whether they **hold up**: whether what is filed as
evidence carries its command, whether a verdict covers what it claims to cover,
whether reviewers who agreed got there by different routes. That question is
about **the report**. The senior's question is about **the project** — a finding
can be impeccably evidenced and still be a rule this project argued over and
decided against. The auditor is not looking for that, and could not see it if it
were.

### Where its memory actually comes from

The senior is a sub-agent like the rest, and the base block above states that
**each sub-agent starts with no context**. So the memory is not a property it
has. It arrives through three channels, and **not one of them is a guarantee**:
two are prompts, and the third is a directory somebody has to open.

**The prompt you write.** The launch template above offers five slots and
**the history of this change is none of them** — its context slot is explicit
that what travels is artifact **paths, not content**. Nothing in the pattern
carries what came before, so if you do not write it, it does not arrive.
**When you launch the senior, add this to its prompt:**

```markdown
## What came before
Read the earlier phases' artifacts of this change under
`openspec/changes/{change-name}/`, and list `openspec/changes/archive/` for the
changes that closed before it. If engram is available, run `mem_search "opsx"`
for this project and `mem_context`. Say which of the three you actually used,
and if you used none, say that instead of reporting no prior context.
```

That is **an instruction, not a guarantee**, on the same footing as every other
prompt in this block: nothing makes you paste it, and nothing makes the
sub-agent obey it once pasted.

**The change directory and the archive, which it reads for itself.** This
change's artifacts are under `openspec/changes/{change-name}/`, and closed
changes under `openspec/changes/archive/YYYY-MM-DD-{name}/`. That is where "the
decision taken two weeks ago" is actually stored, and it is
**readable without any memory system at all**. **ASSUMPTION**
(**read out of the base block, not exercised here**): that archiving really names
the directory that way, dated by the day it ran. If it does not, every citation
below is dated off a bad clock, and a wrong date is more persuasive than none.

**Engram, when its toggle is on** — it is a separate toggle, so
**write the prompt as if it may be off**.

What it offers is more than a keyword guess, and the protocol composed below
spells it out. Every completed OPSX action is saved under the **fixed word
`opsx`**, with a title naming the action and the change, so
**`mem_search "opsx"` enumerates the changes that came before** without knowing
anything about them in advance. `mem_context` recovers recent sessions, and
`mem_get_observation` returns an entry whole. And a decision that outlives one
change **already has somewhere to live**: the protocol's stable `topic_key`,
whose own example is `architecture/auth-model`, **is not scoped to a change at
all**.

So the gap is not in the mechanism. It is that
**all of it is asked for and none of it is enforced** — `mem_search` returns
what somebody actually saved, and every save in that protocol is a prompt
exactly like this one. An empty result therefore means
**"nothing was saved", never "nothing was decided"**, and the senior reports
which of the two it is looking at rather than reading silence as an all-clear.

### What it looks at, phase by phase

One prompt reused five times returns the same generic observation five times, and
a generic observation is exactly the noise that gets the whole thing turned off.
The useful question is different at every point:

| Phase | What the senior asks |
|---|---|
| explore | Is this **already solved in the repo**? Is something about to be rewritten that already exists? |
| propose | Is this **the real problem or the symptom**? Which alternative was **discarded without saying so**? |
| apply | Does this **look like the rest of the code**? Was this test **watched failing**? |
| verify | Are the findings **real or noise**? **What was not tested**? |
| archive | What was learned here that is **lost if nobody writes it down**? |

**Launch it at the close of each phase**, on what the phase returned — never
alongside a phase still running. In `apply` that is not a preference: the role
specialists convene nobody while the implementation sub-agent is still writing,
because **a diff that is still moving gets read three times and judged on none of
its versions**, and that reason applies to the senior unchanged. Wait for the
implementation sub-agent to return, then read the diff that stopped moving —
which is also the only point where declining to open the next phase costs
nothing.

**The senior is not one of the seven role specialists.** It is not chosen by what
the change touches, and **the ledger of roles you did not convene does not
account for it**: it is convened in every phase, or its toggle is off. Filing it
as an eighth role puts it in a table that decides attendance by layer, which is
the one criterion that never applies to it.

### It stops nothing. You do.

**A sub-agent cannot stop anything.** It returns its report and terminates:
nothing of it is still running by the time the report is read. The only thing in
this loop that can decline to launch the next phase is **you, the orchestrator**,
so every rule about holding the cycle is addressed to you.

Do not start the next phase until the auditor has read this context and issued its verdict. The senior never triggers a separate human wait.

Written the other way round, as a power the senior holds over the cycle, this
**describes a mechanism that does not exist** and hands the reader a guarantee
nobody implements. Fragments in this repo have already shipped that error and had
to be corrected.

### The threshold, or you hold the cycle every time

An agent asked to analyse and question **will always find something**. If every
observation it returns holds the cycle, the person is interrupted six times per
change and inside a week the toggle is **switched off** — and a senior that is
off advises nobody, which is worse than one that spoke only when it mattered.

So the work is split: the senior supplies context; the auditor classifies and decides. The orchestrator executes the auditor's verdict. The senior proposes BLOCKING or NOTE against the criteria below; the auditor resolves disagreements with evidence or dated decisions.

**BLOCKING**, and only these four:

- it **contradicts a decision already taken**, and the report **cites it with its date**;
- **something that already exists is being rewritten**;
- **the test just written proves nothing**, and the report shows why;
- the change reached something critical that **nobody asked to touch**.

**NOTE** — it stays in the report and **the phase continues** only according to the auditor verdict: style, naming,
preferences, improvements that are not problems, and everything that opens with
"**you could consider**".

The line, in one sentence: **hold the cycle when the report carries something the
person does not know.** Do not hold it for an opinion about something
**already in plain sight**.

That line is **the auditor's threshold word for word**, and deliberately so: it
is the only cut that separates a report worth an interruption from one that
teaches the reader to skim. What differs between the two is not when to speak.
It is what each one looked at before speaking.

### What a BLOCKING finding has to carry

Three conditions. Missing any one of them, the interruption is worse than what it
interrupted, and two of the three are the senior's to satisfy before you ever see
the report.

**The datum, not the suspicion.** Not "this looks off" but "this contradicts X,
look here", with the file, the line and the date. The finding names in advance
the single observation that **separates the claim from its alternative**, and
goes and gets that one. Where there is none, it is a hunch, and
**a hunch is a NOTE**.

**Options, not an open question.** "What do you want to do?"
**hands the whole problem back**, and the person pays for framing the choice as
well as for making it. "It contradicts X: change it, update the decision, or
carry on and it stays on the record" **hands back a choice**.

**A safe default, and this one is yours.** Apply it when the auditor verdict requires a human decision, not as a parallel wait on the senior report. The senior has already terminated; the
absence of an answer is a state only you can hold. If no answer comes,
**nothing advances**. Never an "I'll assume yes" about the very thing that was
worth stopping for: an assumption that outlives an unanswered question is
**the interruption converted into permission**.

Report it in the typed categories like everything else — the citation with its
file and its date is EVIDENCE, "this contradicts it" is a JUDGMENT with its
confidence, and anything taken for granted is an ASSUMPTION.

### What it costs

**One more sub-agent per phase: five per change** when the cycle runs once end to
end, and six or more when apply or verify go round twice. That is paid on every
change, including the ones where the senior has **nothing to say**.

It is a toggle. If the reports stop being read, **turn it off** rather than
leaving it running and ignored — an ignored senior costs exactly the same as one
that works.
