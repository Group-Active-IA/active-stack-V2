## Ship — The One Verb Somebody Else Sees

Everything else in this cycle happens inside the machine. Explore reads, propose
writes a plan, apply writes code, verify judges it, archive moves a directory.
Every one of them is undone by editing a file.

This one is not. **A push leaves the machine**, and by the time it is wrong
somebody has already fetched it. Force-pushing over it **is not an undo**: it is
a second event on top of the first, which everybody who pulled in between now has
to reconcile. The cost of getting this one wrong is not paid by you.

And it is **not a sequence of steps. It is a diagnosis.** There is no fixed list
of commands, because the right answer depends entirely on where the repository is
standing when you arrive: the same change, on a branch with a pull request open,
and on the trunk with three commits nobody has seen, needs two different answers
and shares none of them. **Look first, decide second.**

### The four preconditions, and all four before anything leaves

1. **The change is closed.** `openspec list --json` no longer carries it as
   active, or `openspec/changes/archive/YYYY-MM-DD-{name}/` exists. Archiving
   moves the whole change directory, so those moves are themselves uncommitted
   work — which is the third precondition, broken by the first.
2. **The suite is green on THIS commit** — run, not remembered. Record the SHA it
   ran at (`git rev-parse HEAD`) beside the result. "The tests passed" with no
   commit attached is true of some commit and says nothing about this one, and
   the block above sends test runs to a sub-agent, so what reaches you is a
   report: without the command and the SHA in it, green is a JUDGMENT.
3. **The tree is clean.** `git status --porcelain` prints nothing.
4. **The `/opsx:verify` verdict is on the record, whatever it says.** A failing
   verdict does not stop the publish on its own — that call belongs to the
   person — **but a missing one does**, because then nothing has judged this work
   against its own artifacts at all.

**If one is missing, name which one and stop.** Not a warning and then the push:
**a warning followed by the push is the push**, and the warning is the part
nobody reads. The report is the precondition that failed, what you ran to find
out, and nothing else.

Where that verdict lives is the honest part. **Nothing in the block above names a
file that `/opsx:verify` writes.** What does report artifact state is
`openspec status --change "<name>" --json`, which this block already calls its
eyes — so ask it and read the answer. If the verdict exists **only in this
conversation**, say that instead of pointing at a path: "verify ran this session
and said X" **is a weaker claim** than "here is the artifact", and reporting the
weak one in the shape of the strong one is how a precondition becomes a
formality.

### The six states, and the answer each one has

Read the repository, then match. The answers are not interchangeable, and four of
the six end in something other than a push.

| What you find | What corresponds |
|---|---|
| The work was committed straight onto `main` | **Stop.** Propose a retroactive branch and move the commits |
| A new branch with no remote yet | `git push -u`, then propose a PR |
| A branch whose PR is already open | Push to it and let that PR update. **Never open a second one** |
| `main` has moved on since the branch left it | Rebase or merge first, and say in the report which of the two you did |
| The tree still has uncommitted changes | Ask whether they belong in this change. **Never decide that alone** |
| No remote is configured | Say so and stop. **Do not invent one** |

**On the trunk, stop before you tidy.** The repair — a branch at HEAD, then the
trunk moved back to where it forked — **rewrites the person's local `main`**, and
that is not something to do to somebody's repository on your own initiative. And
check first whether it is even available: if local `main` is **not ahead of
`origin/main`**, those commits are already published and **no local surgery
unpublishes them**. Then the honest report is "this already went to the trunk,
here is what went", not a plan to move commits that are gone.

**An open PR is the destination, not an obstacle.** Pushing to its branch updates
it. Opening a second one **splits the review across two places**, and the person
now has to read both to know the state of one change.

### The reads, and the one that goes stale

The diagnosis is a handful of reads, and several of them come back empty for a
reason that is not the one you will assume.

| What you want to know | What answers it |
|---|---|
| Which branch you are on | `git branch --show-current` |
| Whether that is the trunk | `git symbolic-ref --short refs/remotes/origin/HEAD` — **the default branch is not always called `main`**, and a hardcoded name diagnoses "we are not on the trunk" on every repository whose trunk is `master` |
| Whether there is a remote | `git remote -v`, which prints nothing on a repository that has none |
| Whether this branch tracks anything | `git rev-parse --abbrev-ref @{upstream}`, which **exits non-zero** when it does not — that failure is the answer, not an error to report |
| Whether the tree is clean | `git status --porcelain` |
| Whether the trunk moved | `git log --oneline origin/main..HEAD` and `git log --oneline HEAD..origin/main`, **both directions**: one of them says what you have that it does not, the other what it has that you do not |

**And `origin/main` is a local cache.** It is a ref on your disk, so
**without a `git fetch` first** it answers with whatever was true the last time
somebody synced: "the trunk has not moved" becomes a sentence about this machine
and not about the remote, and on a machine that has not fetched in a week it
**reports no divergence forever**. Fetch, then compare. If the fetch failed,
**say the comparison did not happen** rather than reporting no divergence — an
answer that arrives because the check never ran looks identical to a clean one.

### What is actually guaranteed here: Bash, git and `gh`

**Nothing in this block grants a tool.** The delegation table above says
`Bash for state (git, openspec status)` runs inline and that Bash for execution
gets delegated — that is a rule about **where** to run a command, not a promise
that the command exists, and it lives inside **a section the `delegation` toggle
removes**. Two variants of this block do not carry that table at all. So the
reads above are what you attempt; if there is no way to run them, the diagnosis
is unavailable and that is the report.

**And notice what the table does not classify.** Its git row is about reading
state; its execution row names tests and builds. **`git push` is not a
read**, and it is the one command here that changes something outside this
machine — **neither row names it**. Treat it as its own case rather than filing
it under "it is git, therefore inline".

**ASSUMPTION** (not verified from this repo): that the GitHub CLI is
**installed and authenticated** for this repository. Nothing in this block
installs it, mentions it, or checks it, and the two halves fail separately — an
unauthenticated `gh` is present and answers nothing. If it is unavailable, the
state "a branch whose PR is already open" becomes unreadable, and the failure
**looks exactly like "there is no PR"** — which is precisely the answer that
opens the second PR the table above forbids. So on any `gh` failure report
**"I could not tell whether a PR exists"**, never "there is none", and hand over
the branch name so the person can look.

### Pushing to `main` asks. Pushing to your own branch does not.

The rule is one line and the reason is the whole of it: the difference is
**whose day breaks if it goes wrong**.

A push to a branch you created **affects nobody**. Nothing builds from it, nobody
has it checked out, and if it is wrong you push again. Asking permission for that
spends an interruption and buys nothing — and an agent that asks for everything
**gets its toggle switched off**, after which it asks for nothing, including the
push that mattered.

A push to `main` is the other thing entirely: it is what CI builds, what the next
person pulls, and on many repositories what deploys. **Ask every time** —
**including the tenth time** on a repository where the answer was yes nine times,
because the tenth is the one where somebody else is halfway through the same
files.

### PR or direct: propose with grounds, do not decide alone

This is the one call the diagnosis genuinely cannot make, and pretending
otherwise is worse than asking. It turns on three facts **that are not in the
repository**: **whether this repo is the person's own or shared**,
**whether `main` is protected**, and whether somebody else is
**working in these files right now**.

What you can do is argue. A recommendation with its grounds is a decision the
person makes in five seconds; a bare "PR or direct?" **hands them the framing as
well as the choice**:

- "this touches authentication and it is 400 lines across six files — **PR**"
- "this is one paragraph in the README — **direct**"

Size, blast radius, and whether it is reversible. Say the recommendation, say
what it rests on, and let the person answer.

If the branch protection can actually be read, that answer is EVIDENCE and worth
having. If the call errors or comes back empty, that is **not** "the branch is
unprotected" — it is "I could not check", and those two must never arrive in the
same sentence.

### It does not merge

Pushing a branch and opening a PR is publishing something for somebody to look
at. Merging it is **deciding on everybody's behalf** that it is good, and that
decision belongs to whoever owns the repository — which is not the process that
wrote the code.

So **the cycle ends at the open PR and its CI result**. Do not run `gh pr merge`,
**do not enable auto-merge** — a merge on a timer is still a merge, taken by the
same process, with the deciding moved out of sight — and do not read "the checks
are green" as the missing approval: green checks mean nobody found a problem
automatically, which **is a different claim from somebody having read it**.

### The diagnosis can be delegated. The answer cannot.

The senior fragment, when its toggle is on, names this same split —
**the senior labels, you hold**. Here it is: the diagnosis reports, you answer.
The rule does not depend on that fragment being present.

Reading a repository is exactly the work this block sends to a sub-agent, and the
reads can go there. What cannot is any of the three questions above — the push to
`main`, the PR-or-direct call, the uncommitted changes.
**A sub-agent that asks a question does not receive an answer.** It
**returns the question in its report and terminates**, and nothing of it is still
running by the time somebody reads it. The person is in a conversation with
**you**, so the asking happens in your thread or it does not happen at all.

Which leaves the split: delegate the reads, collect the state, ask, and hold.
**There is no safe default here.** Everywhere else in this block an unanswered
question costs a day; here the wrong default costs a revert that somebody else
performs. So if no answer comes, **nothing gets pushed** and the change waits.

### Complete CI on the pushed SHA

Before evaluating CI, record the pushed SHA, expected workflows/events and expanded matrix jobs.

Every required job must finish in success on the pushed SHA.

Discover runs with `gh run list --commit <sha>`, then the orchestrator uses
`gh run watch <run-id>` to wait for each expected run. Inspect their jobs with
`gh run view <run-id> --json headSha,event,workflowName,jobs` and retain the result.

Inspect the jobs of every expected run; never choose only the favorable run.

For active-stack, require .github/workflows/ci.yml on push, with test on Ubuntu and Windows.

For a PR, inspect its required checks and verify the source SHA, not the synthetic merge SHA.

In other repositories, derive the required set from their workflows and required checks; a tag release is not required for a branch push.

One successful job with another pending is not green.

queued, in_progress, cancelled, skipped, neutral and timed_out are not green.

A result for another SHA is not green.

Zero runs, no CI or a failed query are not green; report the state and leave delivery unconfirmed.

On failure, report it without stacking another fix as if the push were confirmed.

A new push invalidates the previous result and requires evaluating the new SHA.

The orchestrator waits for the complete CI result before declaring delivery ready.

### What it costs

**One diagnosis before every publish** — a handful of git reads, plus the PR lookup and CI run/job reads
when `gh` is there — and **at most two interruptions** per change: the
PR-or-direct recommendation, and the push itself when the target is `main`. On
the ordinary case, a change on its own branch whose PR is already open, it is
**zero questions** and the reads.

That is paid on the one-line changes too, where the diagnosis is nearly always
the same three sentences.

**It is a toggle.** If the questions are being answered without being read,
**turn it off** rather than leaving it running: confirmations that have become
reflex are not a safeguard, they are a delay with a prompt attached.
