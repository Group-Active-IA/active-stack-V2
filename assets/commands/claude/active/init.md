---
name: "Active: Init"
description: Bootstrap the Active Stack foundation in the current project and hand the conversation to the orchestrator
category: Workflow
tags:
  - active-stack
  - foundation
  - orchestrator
  - workflow
---

Initialize the Active Stack foundation in the project you are standing in, then let the
conversation take over.

This is the **first rung of the ladder**: the user creates a folder by hand, opens Claude,
and runs one command. Everything after that is the orchestrator asking what you are going
to build — the knowledge base, the roadmap, the skills and the project `CLAUDE.md` come
out of that conversation, not out of this command.

**This command reimplements nothing.** `active-orchestrator` owns the whole foundation
flow: entry-point detection, `openspec init`, the phase dispatch
(kb-creator → roadmap-generator → find-skill → skill-registry → agent-instruction), the
shared state file, the checkpoints and the graceful degradation when a sub-skill is
missing. This command's only job is to invoke that skill correctly and get out of the way.

**Input**: none. **Do not ask for a project path** — the user is already inside the
directory they want initialized. If arguments are passed anyway, forward them verbatim as
extra context to the skill; do not interpret them yourself.

---

**Steps**

1. **Check the skill is available — and stop honestly if it is not**

   `active-orchestrator` must be reachable through the **Skill tool**. Confirm it appears
   in the available-skills list for this session.

   If it is **not** available, say so explicitly and stop. Report:
   - which skill is missing (`active-orchestrator`),
   - that the foundation flow needs the **full** install mode of the stack:
     ```bash
     active-stack install --mode full
     ```
     or, for the skill alone:
     ```bash
     npx skills add Group-Active-IA/active-orchestrator -g
     ```
   - that its sub-skills (`kb-creator`, `roadmap-generator`, `find-skill`,
     `skill-registry`, `agent-instruction`) ship with that same mode.

   **Do NOT improvise the phases by hand.** Running `openspec init` and writing a
   knowledge base yourself produces a foundation nothing downstream recognizes: no
   `.active-orchestrator-state.json`, so the next `/active:init` cannot resume and the
   OPSX orchestrator cannot tell an initialized project from a virgin one. Missing skill →
   report and stop. That is the whole answer.

   > A missing **sub**-skill is a different case and is **not** yours to handle: the
   > orchestrator degrades gracefully per phase (offers to install, or skips and records
   > it in state). Only the orchestrator itself being absent stops this command.

2. **Do NOT diagnose the project yourself**

   Resist the urge to look around first. Do **not** run `openspec status`, do **not** read
   `.active-orchestrator-state.json`, and do **not** decide which phase "should" run.

   The skill does all of that in its own Step 0 + resume logic, and it is the only thing
   that reads the state file correctly. A diagnosis made here is a **second opinion the
   skill never sees** — at best duplicated work, at worst you announcing a state and then
   the skill acting on a different one. Hand over the control, do not pre-chew it.

   **One carve-out, and only one**: the knowledge-base check in Step 3. It is not a
   diagnosis and it does not route anything — it is a single defensive read that protects
   something the skill does not protect. Everything else stays off-limits.

3. **The one read you DO make: is there already a knowledge base?**

   Before invoking the orchestrator, check whether `knowledge-base/` exists **and has
   content**. Run exactly this, from the project root:

   ```bash
   if [ -d knowledge-base ]; then
     set -- knowledge-base/*
     if [ -e "$1" ]; then echo KB_CON_CONTENIDO; else echo KB_VACIA; fi
   else
     echo KB_AUSENTE
   fi
   ```

   **Nothing here is installed by anyone.** `[`, `set` and `echo` are shell builtins, so
   this runs on any machine, in `bash`, `sh` or `zsh`, with no `ls`, no `find`, no `fd`.
   That matters because this file **travels as an asset to other people's projects**:
   `fd`/`fdfind`, `rg`, `bat` and `eza` are a preference of one machine, and a gate that
   depends on them is not a gate anywhere else. The trick in the middle: a glob that
   matches nothing is left **literal** in `$1`, and `[ -e ]` on that literal is false.

   Read the answer as **three states, not two**:

   - **`KB_VACIA` or `KB_AUSENTE`** → nothing to protect. Go straight to Step 4 without
     saying anything about it.
   - **`KB_CON_CONTENIDO`** → stop and ask (below).
   - **Neither token came back** — the command errored, you had no way to run it, or the
     output is something you did not expect → **you did not check.** Treat it like
     `KB_CON_CONTENIDO`: ask the same question, and say plainly *"no pude verificar si ya
     hay una base de conocimiento"* instead of naming files you did not see.

   That third state is the whole point. **"I could not check" is not "there is nothing
   there."** If a failed check fell through to silence, the error would become a green
   light: the command would find no knowledge base, conclude there was nothing to protect,
   and walk straight into overwriting one. A check that never ran and a check that came
   back clean must never produce the same behaviour.

   **When you ask** — via `AskUserQuestion`, before invoking anything — name what is
   there: how many files and, if you can get them portably, when they were last written.
   If you cannot date them, **list the names and say the dates were unavailable**; do not
   drop the detail silently. Offer three options:

   - **"Conservarla — saltear la fase de KB"** → invoke the orchestrator as in Step 4, and
     state in the same message that the knowledge base already exists and must NOT be
     regenerated; if the flow reaches the kb phase, the user wants it skipped.
   - **"Regenerarla desde cero"** → invoke normally; the KB will be rewritten.
   - **"Cancelar"** → stop here. Do not invoke the skill. Say nothing was touched.

   STOP and wait for the answer. Do not pick a default.

   **Why this gate exists — and why it is here and not in the skill.** The knowledge base
   is the most expensive thing the whole flow produces and the only artifact that cannot
   be rebuilt without having the entire discovery conversation again. `kb-creator` guards
   it in **Mode B** only: its interactive trigger requires `knowledge-base/` to be absent.
   **Mode A does not guard it at all** — it fires whenever `docs/` holds source files, and
   its own workflow says it generates the ten canonical files *"sin hacer preguntas.
   Fire-and-forget."* Its `Don't use when` section mentions the case, but that is usage
   guidance for a reader, not a gate the flow executes. So on a project with a `docs/`
   folder, a second `/active:init` would overwrite a three-week-old knowledge base without
   asking anyone. This is one question standing in front of an irreversible overwrite.

4. **Invoke the orchestrator on the full flow**

   ```
   Skill("active-orchestrator")
   ```

   Per the skill's Step 0 entry-point table, a **direct `Skill` call with no command**
   defaults to the **full flow** (Step 1 — `openspec init`, then every foundation phase in
   order). That is exactly what `/active:init` means, so no extra routing is needed.

   State the entry point plainly when you invoke it, so nothing is left to inference:

   > Entrypoint: full foundation flow (equivalent to `/active-orchestrator:init`).
   > The user ran `/active:init` from the project root. Apply your resume logic before
   > anything else.

   If Step 3 ended in "Conservarla", add that instruction to the same message.

   > **Naming note (for you, not for the user)**: the skill's frontmatter declares its
   > trigger as `/active-orchestrator:init`. **That command does not exist.** The real
   > command is `/active:init`, in the shorter `active` namespace that also holds
   > `/active:starter-add`. The skill was deliberately left untouched, so it will keep
   > emitting the long name at the user in its own messages — see Step 7 for the part the
   > user needs to hear.

5. **What actually happens if this is run twice**

   Idempotence is **not** implemented here, and it is **not complete** in the skill
   either. What follows is what the skill really does, branch by branch. Do not describe
   it as a safety property it does not have.

   The mechanism is `.active-orchestrator-state.json` at the project root, which the
   orchestrator loads at the start of **every** invocation:

   - **State file absent** → virgin project. Full flow from Step 1.
   - **State file present with `step != "done"`** → half-initialized. The skill asks, via
     AskUserQuestion, *"Hay una fundación en progreso (paso: `{step}`). ¿Continuamos desde
     ahí o empezamos de cero?"* — on "Continuar" it jumps to the recorded step and reuses
     the persisted sections (`kb`, `roadmap`, `skills`, `registry`, `agents`) as context;
     on "Empezar de cero" it deletes the state file first. Both branches are behind that
     answer.
   - **State file present with `step == "done"`** → **this is the hole.** The skill's
     resume logic only branches on `step != "done"`; there is no `else`. Control falls
     through to Step 0, whose entry-point table sends a direct `Skill` call to the **full
     flow (Step 1)**. So an already-founded project re-runs the whole flow. The only line
     in the skill that mentions `step == "done"` lives in *Step 4 — Summary*, which runs
     at the **end** of a flow, not at its start.
   - **`openspec/` already exists** → the skill asks *"¿Lo reuso tal cual o lo regenero?"*.
     That question is about the **OpenSpec scaffolding only** — it says nothing about the
     knowledge base, and **both** of its answers keep going: "Reusarlo" skips to Step 2 and
     dispatches `kb-creator` with no further check, and "Regenerar" additionally rewrites
     the state file.

   So a second `/active:init` **can** reach `kb-creator` on a project that already has a
   knowledge base. Downstream, `kb-creator`'s Mode B refuses to fire when `knowledge-base/`
   exists, but **Mode A does not check at all** and overwrites it silently. That is what
   Step 3's gate is standing in front of.

   **If the user asks whether running it twice is safe, tell them this — and check.**
   Reading `knowledge-base/` and `.active-orchestrator-state.json` to answer a direct
   question is not pre-diagnosis; it is answering with evidence instead of with a promise.

   **Re-running a single phase has no command today.** The skill's frontmatter and its own
   messages cite `/active-orchestrator:kb`, `:rules`, `:find-skill`, `:registry` and
   `:openspec`; **none of those exist as slash commands.** The `active` namespace ships
   exactly two: `/active:init` and `/active:starter-add`. The only door is the full flow —
   with Step 3's gate in front of it.

6. **Let the checkpoints reach the user — do not answer them**

   The full flow holds a checkpoint at **every** phase boundary: a post-phase advance gate
   after each phase (*Continuar / Ajustar / Parar*), plus a pre-phase context gate before
   the roadmap. Those questions are for the person.

   **Never answer them on the user's behalf**, never assume "Continuar", and never
   summarize a phase's output in a way that reads as an answer. If a checkpoint appears,
   stop and wait. The whole point of the gates is that the person sees each artifact
   before the next one is built on top of it.

7. **End inside the conversation, not in a summary**

   This command is done at the **first question the orchestrator puts to the user**. On a
   normal run that is the **advance gate right after `openspec init`** — *"Fase `openspec`
   lista. ¿Cómo seguimos?"* with *Continuar / Ajustar / Parar*. It is **not** the
   `kb-creator` discovery: that only opens after the user answers this gate.

   Other legitimate stopping points, all of them questions: the resume prompt on a
   half-initialized project, the `openspec/`-already-exists prompt, or — if the `openspec`
   CLI is missing — the skill's install message and a stop.

   **That first gate is a checkpoint, not an unfinished step.** Do not read it as "I am
   not done yet" and answer it yourself; Step 6 forbids exactly that.

   **Do not close with "listo" or a recap of what you set up.** The command is a hand-off,
   not a task with a report: the first thing the user should see after the foundation
   starts is a question. If you find yourself writing a final summary, you stopped one
   step short.

   **The one thing you do say to the user**, in a single line before handing over:

   > Para volver a esto, el comando es `/active:init`. Si en algún mensaje ves citado
   > `/active-orchestrator:init`, es este mismo comando — el nombre largo lo emite la
   > skill, pero no existe como comando.

   Say it to **the user**, in their language, not as a note to yourself. The skill prints
   the long name in its own resume hints, its skipped-phase hints and its final summary,
   in sessions where this file is not loaded; if the user copies it from the scrollback
   the next day, it will not exist.

---

**Guardrails**

- **Invoke, do not reimplement.** No `openspec init`, no knowledge base, no `CHANGES.md`,
  no `CLAUDE.md` written from this command. Every one of those belongs to a sub-skill.
- **Do not pre-diagnose** — with the single exception of Step 3's knowledge-base check,
  which decides nothing about routing and only asks a question. The skill's resume logic
  stays the only reader of `.active-orchestrator-state.json`.
- **Never overwrite a knowledge base without an explicit answer.** It is the one artifact
  the flow cannot rebuild without redoing the whole discovery.
- **Assume no tooling is installed.** This file ships to other people's machines. Step 3's
  check uses shell builtins only, on purpose: `fd`/`fdfind`, `rg`, `bat` and `eza` are one
  machine's convention, and `fdfind` does not even exist outside Debian and Ubuntu. If you
  ever add a check here, it survives on a bare box or it does not go in.
- **A check that failed is not a check that passed.** Anywhere you look before acting,
  "I could not tell" is its own outcome and gets said out loud — never folded into the
  reassuring one.
- **Do not ask for the folder.** The current directory is the project.
- **Missing orchestrator → report what is missing and stop.** Never fail silently, never
  improvise the flow by hand.
- **Do not install anything on your own initiative.** `find-skill` recommends and the user
  picks at the skill's own install gate; writing into the global skills is a HIGH-governance
  side effect.
- **Do not commit or push.** Founding a project is not a reason to touch git.
- **Do not claim a safety property you have not checked.** If asked whether a second run
  is safe, look — do not recite.
- **The last word is the orchestrator's question, not yours.**
