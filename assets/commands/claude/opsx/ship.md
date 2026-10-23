---
name: "OPSX: Ship"
description: Publish a closed change — diagnose the repository, check the four preconditions, push, and stop at the open PR
category: Workflow
tags: [workflow, ship, publish, experimental]
---

Publish the work of a closed change. This is the last phase of the OPSX cycle and the
only one that leaves the machine.

Everything before this is undone by editing a file. **A push is not.** By the time it is
wrong somebody has already fetched it, and force-pushing over it is not an undo — it is a
second event on top of the first that everybody who pulled in between has to reconcile.
The cost of getting this one wrong is not paid by you.

**This is not a sequence of commands. It is a diagnosis.** The right answer depends
entirely on where the repository is standing when you arrive: the same change, on a
branch with a PR open and on the trunk with three commits nobody has seen, needs two
different answers and shares none of them. **Look first, decide second.** Four of the six
states below end in something that is not a push.

**Input**: Optionally specify a change name after `/opsx:ship` (e.g. `/opsx:ship add-auth`).
If omitted, check if it can be inferred from conversation context. If vague or ambiguous
you MUST prompt.

---

**Steps**

1. **Resolve the change name**

   ```bash
   openspec list --json
   ```

   Returns `{"changes":[{"name","completedTasks","totalTasks","lastModified","status"}]}` —
   **active changes only**. An archived change is absent from this list by design. But so
   is a name that never existed: absence alone proves nothing, and P1 below is what tells
   the two apart.

   If no name was given and it cannot be inferred, use the **AskUserQuestion tool** to let
   the user choose. Offer both:
   - changes still listed as active (these will fail precondition 1),
   - directories under `openspec/changes/archive/` (these are the shippable ones).

   ```bash
   fd -t d . openspec/changes/archive/ -d 1 2>/dev/null
   ```

   **IMPORTANT**: Do NOT guess or auto-select a change.

2. **THE GATE — the four preconditions, all four before anything leaves**

   Run all four, **in order**. **If one is missing, name which one and stop.** Not a
   warning and then the push: a warning followed by the push *is* the push, and the
   warning is the part nobody reads. The report is the precondition that failed, what you
   ran to find out, and nothing else.

   **The only thing that reopens a closed gate is the person saying so — never your own
   judgment.** Two of the four have a legitimate path forward: P3, which archiving breaks
   by design, and P2 on a repository with no detectable suite. Both of those paths run
   through **AskUserQuestion**, and neither is a call the command makes for itself. There
   is no third case, and "it seemed fine" is not one.

   **And a precondition skipped with permission gets written down.** It goes into the
   report's precondition table as `✗ skipped — authorized by the person`, naming what was
   missing, and it stays there even though the push happened. Skipped with permission is a
   decision; skipped in silence is an accident somebody discovers a week later.

   **P1 — The change exists, and it is closed.**

   ```bash
   openspec list --json
   fd -t d "^\d{4}-\d{2}-\d{2}-<name>$" openspec/changes/archive/ -d 1 2>/dev/null
   ```

   No output from that `fd` means the same thing whether the directory is missing or holds
   nothing matching: **there is no archived copy of this change**.

   Absence from `list --json` has **two causes, and they are not the same one**: the change
   was archived (which is what this precondition wants), or it **never existed** — a typo
   in the name you were handed. The datum that discriminates them is the archive directory,
   so that directory is the **confirmation**, never an alternative:

   - Still in `list --json` → **stop**. The answer is `/opsx:archive`, not a push.
   - Absent from the list **and** an `openspec/changes/archive/YYYY-MM-DD-<name>/`
     directory exists → P1 passes.
   - Absent from the list **and no archive directory** → **that change does not exist**.
     Report the name you were given, offer the names that do exist (step 1's two lists),
     and **stop here** — before a single test runs.

   That last branch is the whole point of the split. Read as an `or`, a typo like
   `/opsx:ship add-atuh` *passes* P1, pays for the entire suite in P2, sails through P3,
   and only meets the truth at P4, where `openspec status --change "add-atuh" --json`
   exits 1 with `Error: Change 'add-atuh' not found`.

   **P2 — The suite is green on THIS commit.** Run, not remembered.

   ```bash
   git rev-parse HEAD
   ```

   Detect the suite from the repository (`package.json` scripts.test, `go test ./...`,
   `pytest`, `Makefile` target, whatever this project actually uses) and **run it**.
   Delegate the run to a sub-agent if it is long; what must come back is
   **the command, its output, and the SHA it ran at**.

   - Green → record `command + output + SHA` as EVIDENCE.
   - Red → **stop**, and report which tests failed.
   - No suite could be detected, or the runner is unavailable → you did **not** verify P2.
     Say exactly that, and use **AskUserQuestion** to let the person choose: *"this repo
     has no suite, proceed"* / *"the command is: …"* / *"cancel"*. **Never report green for
     a suite you did not run**, and never let "I could not check" arrive in the shape of
     "it passed".

   "The tests passed" with no commit attached is true of some commit and says nothing
   about this one. Without the command and the SHA in it, green is a JUDGMENT, not
   EVIDENCE.

   **P3 — The tree is clean.**

   ```bash
   git status --porcelain
   ```

   Passes only if it prints nothing.

   If it prints something, look at *what*: archiving moves the whole change directory, so
   **P1 breaks P3 by design** — dirty paths that are exactly `openspec/changes/…` are the
   archive move waiting to be committed. Say so, and go to state 5 of the diagnosis
   (uncommitted changes). **Never decide alone** whether uncommitted work belongs in this
   change.

   **P4 — The `/opsx:verify` verdict is on the record, whatever it says.**

   ```bash
   openspec status --change "<name>" --json
   ```

   Returns `{"changeName","schemaName","isComplete","applyRequires":[…],"artifacts":[{"id","outputPath","status"}]}`.

   A failing verdict does not stop the publish on its own — that call belongs to the
   person. **A missing one does**, because then nothing has judged this work against its
   own artifacts at all.

   Be honest about where the verdict lives. **No artifact id in the schema is the verify
   verdict**, and nothing guarantees a file. Look, in order:
   - a `review-verdict.md` (or similarly named verdict artifact) inside the change
     directory or its archived copy,
   - an artifact in `openspec status --json` that names verification,
   - failing both: the verdict exists **only in this conversation**.

   If it is only in the conversation, **report it that way**: *"verify ran this session and
   said X"* is a weaker claim than *"here is the artifact"*, and reporting the weak one in
   the shape of the strong one is how a precondition becomes a formality.

   **If verify never ran at all: stop.** The answer is `/opsx:verify` — *if that command is
   installed here*. Check before you point at it, because it is not always there:

   ```bash
   fd -t f '^verify\.md$' ~/.claude/commands/opsx/ -d 1 2>/dev/null    # global install
   fd -t f '^verify\.md$' .claude/commands/opsx/  -d 1 2>/dev/null    # this project
   ```

   Two separate calls on purpose: a project with no `.claude/commands/opsx/` makes `fd`
   print an error and still **exit 0**, so a single combined call mixes a missing path in
   with a real result.

   - **Found** → **stop**, and run `/opsx:verify` first.
   - **Not found** → do **not** send the person to a command that does not exist. A gate
     whose only exit is a missing command has no exit, and the only way anybody gets past
     it is by skipping P4 — which is exactly the formality this precondition exists to
     prevent. The honest report is the two facts: **nothing has judged this work against
     its own artifacts**, and `/opsx:verify` is not available here. That is a finding, not
     a verdict, and the decision is the person's: use **AskUserQuestion** — *"ship anyway,
     unverified and on the record"* / *"cancel, and install the verify command first"*.
     **Never** treat the missing command as P4 passing. If the person authorizes
     proceeding, P4 goes into the report as **skipped with permission**, per the rule at
     the top of this gate.

3. **THE DIAGNOSIS — read the repository before deciding anything**

   ```bash
   git branch --show-current
   git remote -v
   git symbolic-ref --short refs/remotes/origin/HEAD
   git rev-parse --abbrev-ref @{upstream}
   git status --porcelain
   ```

   How to read each answer:

   | What you want to know | What answers it | How it fails |
   |---|---|---|
   | Which branch you are on | `git branch --show-current` | empty on a detached HEAD |
   | **What the trunk is called** | `git symbolic-ref --short refs/remotes/origin/HEAD` → e.g. `origin/main`; strip the `origin/` prefix | **exits non-zero when `origin/HEAD` is not set.** That is "I could not tell", NOT "the trunk is `main`". Recover it with `git remote show origin` (network) or ask. **Never hardcode `main`** — it diagnoses "we are not on the trunk" on every repo whose trunk is `master` |
   | Whether there is a remote | `git remote -v` | prints nothing on a repo that has none → state 6 |
   | Whether this branch tracks anything | `git rev-parse --abbrev-ref @{upstream}` | **exits non-zero when it does not track. That failure is the answer, not an error to report** → state 2 |
   | Whether the tree is clean | `git status --porcelain` | anything printed → state 5 |

   **Then fetch, and only then compare.**

   ```bash
   git fetch --quiet origin && echo FETCH_OK || echo FETCH_FAILED
   git log --oneline origin/<trunk>..HEAD    # what you have that it does not
   git log --oneline HEAD..origin/<trunk>    # what it has that you do not
   ```

   `origin/<trunk>` is a **ref on your disk**, a local cache. Without a fetch first it
   answers with whatever was true the last time somebody synced: on a machine that has not
   fetched in a week it reports **no divergence forever**. Both directions matter — one
   says what you have that the trunk does not, the other what the trunk has that you do
   not.

   **If the fetch failed, report that the comparison did not happen.** Never "there was no
   divergence": an answer that arrives because the check never ran looks identical to a
   clean one, and that is the whole trap.

   **Is there already a PR?**

   ```bash
   gh pr list --head "<branch>" --state open --json number,url,title
   ```

   `gh` may be absent, or present but unauthenticated — **the two halves fail separately**.
   On **any** `gh` failure, non-zero exit, or auth error, the report is
   **"I could not tell whether a PR exists"**, and you hand over the branch name so the
   person can look. It is **never** "there is no PR": that confusion is precisely what
   opens the second PR on a change that already had one.

4. **MATCH THE STATE — and four of the six do not end in a push**

   | What you found | What corresponds |
   |---|---|
   | 1. The work was committed straight onto the trunk | **Stop.** Propose a retroactive branch and moving the commits — see below |
   | 2. A new branch with no remote yet (`@{upstream}` exited non-zero) | `git push -u origin <branch>`, then propose a PR (step 5) |
   | 3. A branch whose PR is already open | Push to it and let that PR update. **Never open a second one**, and **skip step 5** — the destination is settled |
   | 4. The trunk has moved on since the branch left it (`HEAD..origin/<trunk>` is non-empty) | Reconcile **first**: **merge** if the branch is already published, rebase only if it is not. Say in the report **which of the two you did** |
   | 5. The tree still has uncommitted changes | **Ask** whether they belong in this change. **Never decide that alone** |
   | 6. No remote is configured | Say so and **stop**. **Do not invent one** |

   **The table reads in singular and a repository can match three rows at once** — standing
   on the trunk, with a dirty tree, with the trunk moved. So the order is fixed, and it is
   not the order of the rows:

   1. **State 1 first, always.** It is the only row that ends in *stop*, and a stop is not
      something the other rows get to override.
   2. **State 6** — no remote means there is nothing to push to; everything below is moot.
   3. **State 5** — resolve the tree before anything gets committed, by asking.
   4. **State 4** — the trunk moved; it has to be reconciled before any push can succeed.
   5. **States 2 and 3** last. They are the two that actually push, and they are mutually
      exclusive: either the branch tracks something or it does not.

   **States 3 and 4 together — the ordinary case, and the one that traps.** A branch with
   an open PR whose trunk has moved on is what every repository with more than one person
   looks like, and the two instructions collide: state 4 says rebase or merge, state 3 says
   `git push` without force, and the guardrails forbid force over published history.

   Rebase is the half that breaks. It rewrites commits that are **already on the remote**,
   git then refuses the push as non-fast-forward, and the `--force` that would get past it
   is both forbidden here and destructive to the tree of everyone who pulled that branch.

   **At that crossing the answer is merge, not rebase:**

   ```bash
   git merge origin/<trunk>          # NOT git rebase: the branch is published
   git push origin <branch>          # fast-forward, no force, the PR updates itself
   ```

   Rebase stays available only where nothing was published yet — state 4 crossing **state
   2**, a branch with no upstream. Say in the report which of the two you did, and why.

   **State 1 — on the trunk, stop before you tidy.** The repair (a branch at HEAD, then the
   trunk moved back to where it forked) **rewrites the person's local trunk**, and that is
   not something to do to somebody's repository on your own initiative. Use
   **AskUserQuestion** before touching anything.

   And check first whether the repair is even available:

   ```bash
   git log --oneline origin/<trunk>..HEAD
   ```

   If local trunk is **not ahead** of `origin/<trunk>`, those commits are already
   published and **no local surgery unpublishes them**. Then the honest report is
   *"this already went to the trunk, here is what went"* — not a plan to move commits that
   are gone.

   **State 3 — an open PR is the destination, not an obstacle.** Pushing to its branch
   updates it. Opening a second one splits the review across two places, and the person
   now has to read both to know the state of one change.

5. **PR or direct — recommend with grounds, do not decide alone**

   **Skip this step entirely if step 3 found an open PR for this branch** (state 3). The
   repository already answered: the destination is that PR, and asking anyway is a question
   whose answer is on the screen. The ordinary case — a change on its own branch with the
   PR already open — is **zero questions and the reads**, and it only stays that way if
   this step knows when not to run. A question that answers itself teaches the person to
   click without reading, and that is the road by which the whole thing ends up switched
   off.

   *"I could not tell whether a PR exists"* is **not** a found PR. Ask.

   This is the one call the diagnosis genuinely cannot make. It turns on three facts
   **that are not in the repository**: whether this repo is the person's own or shared,
   whether the trunk is protected, and whether somebody else is working in these files
   right now.

   What you can do is **argue**. Measure the change first:

   ```bash
   git diff --stat origin/<trunk>..HEAD
   ```

   Optionally read the protection, and read the failure honestly:

   ```bash
   gh api "repos/{owner}/{repo}/branches/<trunk>/protection" --jq '.required_pull_request_reviews != null'
   ```

   If that call succeeds, the answer is EVIDENCE and worth having. If it **errors or comes
   back empty, that is "I could not check"** — not "the branch is unprotected". Those two
   must never arrive in the same sentence.

   Then use **AskUserQuestion** with your recommendation **and what it rests on** — size,
   blast radius, reversibility:

   - *"this touches authentication and it is 400 lines across six files — **PR**"*
   - *"this is one paragraph in the README — **direct**"*

   A bare "PR or direct?" hands the person the framing as well as the choice. A
   recommendation with its grounds is a decision they make in five seconds.

6. **The push — asking depends on who pays if it goes wrong**

   **Push to a branch you created: do NOT ask.** It affects nobody. Nothing builds from
   it, nobody has it checked out, and if it is wrong you push again. An agent that asks for
   everything gets switched off, after which it asks for nothing — including the push that
   mattered.

   ```bash
   git push -u origin <branch>      # state 2
   git push origin <branch>          # state 3, updating the open PR
   ```

   **Push to the trunk: ASK EVERY TIME, via AskUserQuestion — including the tenth time.**
   It is what CI builds, what the next person pulls, and on many repositories what
   deploys. The repository where the answer was yes nine times is exactly where the tenth
   is the one somebody else is halfway through the same files.

   **If no answer comes, nothing gets pushed and the change waits.** There is no safe
   default here: everywhere else in the cycle an unanswered question costs a day; here the
   wrong default costs a revert that somebody else performs.

   Never `--force` or `--force-with-lease` over published history without the person
   explicitly saying so for that specific push.

7. **Open the PR and evaluate its CI result**

   ```bash
   gh pr create --base "<trunk>" --head "<branch>" --title "<title>" --body "<body>"
   ```

   Only if step 5 actually ran and chose PR. If step 3 found an open PR, step 5 was
   skipped and **there is nothing to create** — the push already updated it. And if step 3
   said *"I could not tell whether a PR exists"*, do **not** create one either: hand over
   the branch name and let the person look.

   **It does not merge. Ever.**
   - No `gh pr merge`.
   - No `--auto` / auto-merge: a merge on a timer is still a merge, taken by the same
     process, with the deciding moved out of sight.
   - Green checks are **not** the missing approval. "Nobody found a problem automatically"
     is a different claim from "somebody read it".

   Pushing a branch and opening a PR is publishing something for somebody to look at.
   Merging is deciding on everybody's behalf that it is good, and that decision belongs to
   whoever owns the repository — which is not the process that wrote the code.

   **The cycle ends at the open PR and its CI result.**

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

8. **Report — in typed categories**

   ```
   ## Ship: <change-name>

   ### Preconditions
   | # | Precondition | Result |
   |---|---|---|
   | 1 | Change exists and is closed | ✓ archived 2026-09-09 / ✗ still active / ✗ NO SUCH CHANGE |
   | 2 | Suite green on this commit | ✓ <cmd> — N passed @ <sha> / ✗ red / ✗ skipped — authorized by the person (no suite) |
   | 3 | Tree clean | ✓ / ✗ <paths> / ✗ skipped — authorized by the person (archive move committed) |
   | 4 | Verify verdict on record | ✓ <artifact> / conversation only / ✗ skipped — authorized by the person (verify never ran / command not installed) |

   Any row reading `skipped — authorized` **stays in this table even though the push
   happened**. That is the record.

   ### Diagnosis
   Branch: <branch>   Trunk: <trunk> (from symbolic-ref / could not read)
   Upstream: <upstream or "none — untracked">
   Fetch: OK / FAILED (comparison NOT made)
   Divergence: <n> ahead, <m> behind / not compared
   PR: #<n> <url> / none found / COULD NOT TELL
   States matched: <all that applied>   Resolved in order: <the one acted on, and why>
   Trunk reconciled by: merge (branch published) / rebase (nothing published) / n/a
   Step 5 (PR or direct): asked — <recommendation> / skipped — PR #<n> already open

   ### What happened
   <the push, the PR, or the reason nothing left the machine>

   EVIDENCE   <command + output + timestamp>
   JUDGMENT   <recommendation, what it was compared against, confidence>
   ASSUMPTION <what was taken for granted and what you did instead of verifying>
   ```

---

**Guardrails**

- **The gate is a gate.** A failed precondition ends the command. Never warn-then-push. The
  **only** thing that reopens it is the person authorizing it through AskUserQuestion —
  never the command's own judgment — and the skip is named in the report's precondition
  table even after the push succeeds.
- **A change that is in neither the active list nor the archive does not exist.** Say so
  and stop before running a single test; never read that absence as P1 passing.
- **P4 pointing at `/opsx:verify` requires `/opsx:verify` to be installed.** Where it is
  not, report that nothing judged this work against its artifacts and hand the decision
  over. A gate whose exit does not exist is not a gate, it is a detour.
- **Look before you act.** No push before the diagnosis of step 3 is complete.
- **Published branch + moved trunk → merge, never rebase.** Rebasing published history
  forces a force-push, and force over published history is forbidden here.
- **States are evaluated in a fixed order** (1, 6, 5, 4, then 2/3), because a repository
  routinely matches several at once and only state 1 ends the command.
- **Do not ask what the repository already answered.** An open PR found in step 3 skips the
  PR-or-direct question entirely.
- **Fetch before comparing** against `origin/<trunk>`. A failed fetch means the comparison
  did not happen, never that nothing diverged.
- **The trunk is not always `main`.** Read it; if you cannot, say you could not.
- **`gh` failing means "I could not tell", never "there is none".**
- **Ask for the trunk push, every single time. Never ask for a push to your own branch.**
- **Never decide alone**: the trunk push, PR-or-direct, and what to do with uncommitted
  changes. A sub-agent cannot ask these — it returns the question and terminates. The
  asking happens in this thread or it does not happen.
- **No merge, no auto-merge, no force over published history.**
- If the reads cannot be run at all, the diagnosis is unavailable — **and that is the
  report**.
