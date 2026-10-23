## The Knowledge Base — The Map Nothing Writes Back To

`kb-creator` runs once, during the project's foundation, and produces the
knowledge base: the vision, the stack, the actors, the data model, the business
rules, the stories, the flows, the architecture, the decisions and the open
questions. Then it is marked done.

**Nothing in the cycle ever writes back to it.** The startup routine above says
so plainly: if the foundation is already in place, go straight to `openspec
list`. From that point on the knowledge base is never opened again.

Three months later the data-model document describes a schema that changed two
changes ago, and the agent believes it, because it is the canonical document. It
hallucinated nothing: **it believed the map, and the map was old.**

A wrong map with authority costs more than no map at all. Whoever has no map
asks. Whoever has a wrong one **builds on top of it**.

### Read the directory before you believe any list of documents

The filenames belong to `kb-creator`, not to this block. List the knowledge base
directory and work from what is actually in it.

**ASSUMPTION** (not verified from this repo): that the directory is called
`knowledge-base/` and holds those ten documents under those names. `kb-creator`
is cloned from a separate repository and nothing here pins what it writes. If it
renames or renumbers, a hardcoded list matches nothing and you report
**"this project has no knowledge base"** about a project that has one — the same
wrong answer this section exists to prevent, reached from the other side.

### At session start: report the drift, do not repair it

Before working on a project that has a knowledge base, compare **when its
documents were last written** against **the changes archived after that**. If
closed changes are newer than the map's last edit, **say it out loud before
doing anything else**, and **name them**:

```
JUDGMENT  (orchestrator)  The knowledge base may be behind.
  Changes archived since the map was last written:
    add-billing-webhooks (2026-07-02), split-user-table (2026-08-19).
  Per document, last written:
    data model      2026-06-14  — behind both.
    business rules  2026-06-14  — behind both.
    architecture    2026-08-22  — behind neither.
    (the other seven: 2026-06-14, behind both.)
  Date used: git log --format=%aI per document, see below.
  Confidence: medium — the dates say when bytes changed, not when anybody
  checked the document was still true.
```

**One date for the whole knowledge base is the same mistake as no names.** It
tells the reader that something went stale without telling them
**which of the ten documents to open**, so the report costs a minute of reading
and hands back a shrug. Date each document separately: the archive section
already reports per document, and a drift report that does not
**cannot be acted on without redoing the work it just did**.

Then stop. **Do not rewrite it on your own.** It is the most authoritative
document in the project, and a document that carries that weight
**is not edited by a process nobody is watching**. The next agent reads your
silent edit with exactly the authority the original had, and nothing records
that it came from a guess. You report the drift; **the person decides**.

**A drift report with no names is noise.** "The knowledge base may be out of
date" is true of every knowledge base on every day, so it teaches the reader to
skip the line. The names are what make it checkable in a minute.

And say it **before the work, not after**. Afterwards the map has already been
believed, which is the entire failure.

### What the two dates actually mean

The comparison rests on two dates and **neither one means what you want it to**.
Say which one you used, and never turn "I found no drift" into a statement that
the documents are still true.

**The archive date** comes from the directory name. Archived changes land in
`openspec/changes/archive/YYYY-MM-DD-<name>/`, so listing that directory yields
the names and the days without opening a single file. **ASSUMPTION** (read out
of the block above, not exercised here): that `openspec archive` really names the
directory that way, dated by the day it ran. If it does not, **one side of the
comparison is empty** and the check reports no drift every time, forever.

**The document date is the harder half**, and it has three failure modes rather
than one.

A file's modification time is not the date somebody edited it: `git clone` stamps
**every file with the checkout time**, so on a fresh clone the whole knowledge
base looks written today and the check finds nothing — on the machine that just
arrived and has the least context of all.

Git is better, but **use the author date and not the committer date**:
`git log -1 --format=%aI -- <path>`. The committer date `%cI` is
**rewritten by rebase, by amend, and above all by a squash-merge** — GitHub's
squash stamps every document touched anywhere in the branch with the time of the
merge, **without a byte changing in any of them**. That is the same failure as
the fresh clone, arriving by the road offered as the remedy. `%aI` survives all
three; it is still only **when the bytes last changed**, so a bulk reformat or a
mass rename **refreshes documents nobody reviewed**.

And **this side goes empty too.** `git log -1 -- <path>` on a path git is not
tracking prints **nothing and exits 0** — not an error, so a script that only
checks the exit code reads the empty string as a date and compares against it.
If `knowledge-base/` is in `.gitignore`, or if the knowledge base was generated
and **not committed yet** — the ordinary state immediately after `kb-creator`
runs, which is the exact moment this fragment is about — you get no date and
report no drift, permanently. It is the same mechanism and the same consequence
as the archive side above, so treat an empty date as
**"I could not check", never as "nothing changed"**, and say which one it was.

So what this check answers is
**"were bytes written to this document after that change closed"**, not "is this
document still true". It catches the common case — months of changes against an
untouched map — and it is **blind to the document edited last week and wrong
anyway**. That case needs a reader, and this section does not provide one.

### In propose: name the documents the change will touch

The plan declares **which knowledge-base documents this change will affect**, and
for each one **which of the two operations** it needs: a rewrite or an append.

Naming them in advance is what gives `archive` something concrete to act on.
**Asked afterwards the question gets a shrug**, because by then the change
touched a bit of everything and "which documents did this affect" has no better
answer than "several, probably".

It goes in `proposal.md` under a heading of its own — `## Knowledge Base Impact`
— so archive can open it by name instead of hunting for it in the prose.

**"Nothing to update" is an answer**, and it belongs in the section. A change
that genuinely touches no document says so in one line, because an absent
section and a considered "none" **look identical afterwards**, and the second is
much rarer than the first.

**When you launch the propose sub-agent, add this to its prompt.** Nothing else
carries it. `openspec instructions` generates the artifact templates and
**this heading is not in them**; the sub-agent prompt template in the block above
asks for five things and **the knowledge base is none of them**; and that same
block opens by saying **each sub-agent starts with no context**. A rule written
in this section and nowhere else is addressed to somebody who never reads it —
**this section's reader is you, and the sub-agent's reader is the prompt you
write**.

```markdown
## Knowledge Base
List the knowledge base directory. Add a `## Knowledge Base Impact` section to
proposal.md naming each document this change will affect and, for each one,
whether it needs a REWRITE (status, next steps, frontmatter) or an APPEND
(decisions, gotchas). If it affects none, write "none" and one line of why.
```

Two links, and **both are requests**: nothing makes you paste this, and nothing
makes the sub-agent obey it once pasted. Finding no such section at archive time
is the ordinary case, not an error — **fall back to reading the diff**, and say
in the report that that is what you did.

### In archive: rewriting and appending are two operations, never one

Cover the union of Knowledge Base Impact and the actual diff, even when the impact section exists. Inspect all KB documents touched by either source; reconcile omissions in the report.

If auditor is absent, mark "auditor source absent": update facts with evidence but do not invent LECCIONES or DEUDA. Closure needs Juani's written decision accepting degraded scope or enabling auditor. If knowledge-base is absent, preserve the verdict lessons and debt in the change and tell Juani the destination is pending.

Archiving is already the moment the system declares that something stopped being
a proposal and became true. That moment existed and **nobody was using it**.

**Rewrite** — status, next steps and frontmatter. If the data model changed, the
document **says the new schema and the old sentences are gone**. Not appended
underneath, with the superseded paragraph and its replacement sitting there
looking equally current: a reader who meets two contradicting descriptions of the
schema, with no way to date either, **believes whichever comes first**.

**Append**: decisions, lessons L-NNN and debt, each with **its date and what it cost**. Only append lessons; never rewrite numbered lesson history. Take lessons from the auditor LECCIONES and debt from its DEUDA list. Each entry records what changed and typed evidence: SHA, test or decision. Preserve distinctions between EVIDENCE, JUDGMENT and ASSUMPTION. That
is history, and **history accumulates, it is not overwritten**. Deleting a
superseded decision destroys the only record of why the current one won, and the
next person re-litigates it from scratch.

The line between them: a section that answers **"how is it now"** gets rewritten,
a section that answers **"how did we get here"** gets appended to.

Nothing in the tooling does this. `openspec archive` syncs the delta specs into
the main specs and **knows nothing about the knowledge base**: no CLI command
performs this step, and no part of the block above mentions it either.

**So when you launch the archive sub-agent, add this to its prompt:**

```markdown
## Knowledge Base
Read `## Knowledge Base Impact` in proposal.md and the actual diff, even when the section exists. Cover their union; if the section is missing, say so. Then, for each document either source names: REWRITE status, next
steps and frontmatter so they describe the new state and the old text is gone;
APPEND decisions, lessons L-NNN from auditor LECCIONES and debt from auditor DEUDA, dated, with what each one cost. Only append lessons. Each entry states what changed and typed evidence: SHA, test or decision. If the auditor is absent, mark auditor source absent and do not invent its lists; closure requires Juani's written decision on degraded scope. Return the two
lists separately and by document, headed REWRITTEN and APPENDED. If REWRITTEN
is empty, say why in one line and label the report SUSPICIOUS: append-only when only APPEND was done. Then open the previous closed change's
`knowledge-base-report.md`, carry whatever it left undone into a third list
headed PENDING with the change and the date that incurred it, and write all
three lists into this change as `knowledge-base-report.md`. If there is no
previous report, say so instead of writing an empty PENDING.
```

That is one more request on top of another: nothing makes you paste it, and
nothing makes the sub-agent do it. **Asked for, not enforced** — the same footing
as everything else in this block, stated where the asking actually happens
instead of one level too late.

### The detector: appending alone is the half that comes free

If only APPEND was done without REWRITE, label the report **SUSPICIOUS: append-only** and explain why; this flags suspicion, not an invented rewrite.

Appending comes out easy and it feels like work: one new entry at the bottom,
dated, with a reason. Nothing had to be re-read and nothing had to be deleted.
**Rewriting is the half that forces you to re-read what was written three months
ago** and delete what stopped being true — and that deletion is the whole point,
because the stale line is the one the next agent acts on.

So report the two lists **separately and by document**:

```
JUDGMENT  (archive agent)  Knowledge base after "split-user-table".
  REWRITTEN  data model — users is now two tables; removed the three
             paragraphs describing the single-table shape.
  APPENDED   decisions — 2026-09-08, the split cost a migration whose down
             had to be exercised against a copy of production data.
  Confidence: high — I opened both documents and read them whole.
```

**An empty REWRITTEN list is a claim, not a blank.** A change can legitimately
add only history, and saying so in one line with the reason is a complete answer.
Leaving it blank is not, because **"I updated the knowledge base" is the same
sentence** whether the hard half was done or skipped, and without the split
**nobody can tell those two apart** afterwards.

### The report is the change's closing condition

Everything above is how to do the work. This is the part that says the work
happened — and archive needs one, because **the phase has no verifiable output at
all**. Across two measured projects, archive returned
**not a single finding in ten changes**. That is not ten clean archives: the
three omissions that did surface were found by the auditor reading the reports
afterwards, never by the phase whose job it was.

The other run is the sharper one, because it is dated. A Java project carries the
commit `448d9ab docs: ejecutar los tres REWRITE de archive que quedaron
pendientes de C-01 y C-02` — **three rewrites missed across two whole changes**,
and one of them was the document the next change was reading as its acceptance
criterion. Both changes closed on time. Nothing objected, because
**nothing had anything to check**.

So: **the change is not closed until the two lists are on the record.** Not a
reminder and then the archive — a reminder followed by the archive is the
archive, and the reminder is the half nobody reads. If the report is not there,
name which list is missing and stop before declaring the change closed.

**What holds it is you, not the CLI.** `openspec archive` runs whether or not the
report exists, for the same reason the section above already gives. The hold is
the orchestrator declining to call the change closed, which is the only lever
anybody here actually has — and it is the same one a missing precondition uses
before a publish.

**The two lists stay two, and stay per document.** REWRITTEN names each document
whose state text was replaced and what was deleted out of it; APPENDED names each
document that gained a dated entry and what it cost. Collapsed into
"I updated the knowledge base" the report is worth nothing, because that sentence
is **true of the easy half on its own**, and the easy half is the one that
survives a hurried archive.

**And an empty REWRITTEN list closes the change only with its reason beside it.**
A change that genuinely moved no state — a decision recorded, a gotcha written
down, nothing that answers "how is it now" — is a complete answer when it says
so: `REWRITTEN none — added a retry policy, changed no documented behaviour`.
What does not close the change is the blank. Take that obligation out and the
whole section falls back to the sentence it was written to replace: the report
**reads the same whether the hard half was done or skipped**, and three rewrites
go missing for two changes before anybody notices.

### Where the report lives, and the one assumption it borrows

Write it into the change as `knowledge-base-report.md`, beside the other
artifacts, in the typed categories — the JUDGMENT shape the detector above
already shows, not a second format next to it. A report that exists only in the
archive sub-agent's return **terminates with the sub-agent**, and the section
below has to open it again next month.

**This report shares one assumption with the review verdict, and it is one
assumption of the system, not two**: that archiving moves the whole change
directory, so what was filed inside it is still there afterwards. Where the
`review` fragment is on the page that assumption is already declared there in
full, and it is not restated here. If archiving instead copies a known artifact
list, the verdict and this report are dropped together and for the same reason —
and the carried debt below goes with them.

### The debt travels forward, one hop at a time

A report that covers only the change being closed leaves the previous change's
misses invisible, and that is exactly how three rewrites survived two archives in
the Java project: by the time anybody looked, each one was somebody else's
omission. Whoever wrote `448d9ab` had to **go and find them**, months later, in
an archive nobody reads for pleasure.

So the report carries a third list, and it is short on purpose. **PENDING** holds
what the *previous* change's report left undone: every document its
`## Knowledge Base Impact` named for a rewrite that never turned up in its
REWRITTEN list, plus whatever that report's own PENDING still carried. Each entry
keeps **the change that incurred it and the date**, so a debt on its third hop
says so without anybody counting.

**One hop back, never the whole archive.** The look-back is one file — the last
closed change's report — and not a sweep of `openspec/changes/archive/`. A report
that re-audits every change since the beginning is an inventory, an inventory
gets skimmed, and a skimmed list is worth what no list is worth while costing the
whole archive to produce. Debt nobody pays is handed by each report to the next
one: it stays visible by **travelling**, not by being searched for.

**And a missing previous report is not an empty PENDING.** If the last change
closed before this rule existed, or left no `knowledge-base-report.md` at all,
write **"could not carry debt forward: <change> left no report"** and nothing
else. Those two lines look identical on the page and only one of them means there
is no debt — the same failure as an empty date in the drift check above, arriving
at the other end of the cycle.
