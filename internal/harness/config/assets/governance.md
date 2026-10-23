## Agent Governance — Autonomy by Domain

This governance model (from the Active Stack Methodology, Stage 4) defines how much autonomy the agent has based on the criticality of the domain being modified. It applies to **all agents**, regardless of provider.

### Governance Levels

| Level | Typical Domains | Agent Behavior |
|-------|-----------------|----------------|
| **CRITICAL** | Auth, Billing, Security, Audit logs | Analysis only; no code written without explicit human approval. |
| **HIGH** | Config files that affect user data, file injection into user dotfiles, backup/restore | Propose and wait for review before writing. |
| **MEDIUM** | Business logic, domain adapters, pipelines | Implement with checkpoints; surface decisions to the user. |
| **LOW** | Simple CRUDs, catalogs, type definitions, read-only utilities | Full autonomy if tests pass. |

### How to Apply

1. **Before any non-trivial action**, identify the governance level of the domain you are about to modify.
2. At **HIGH** or **CRITICAL** level: describe the planned change and wait for user confirmation before writing.
3. At **MEDIUM** level: implement in steps; surface non-obvious decisions for review.
4. At **LOW** level: proceed autonomously; report what was done in the summary.

### Written approval artifact

At HIGH or CRITICAL, the human records approval in proposal.md under `## Aprobación`. It must be written, dated and explicit scope, naming the choice, person, date and approved scope:

```markdown
## Aprobación

| Elección | Persona | Fecha | Alcance |
|---|---|---|---|
```

No agent fills this table, signs for the person or infers approval from a brief or commit author. Missing or ambiguous approval stops construction. A new rule after signature requires its own row completed by the human before construction (L-029); the original approval does not authorize it. Wording conflicts are resolved and declared, not presented as a new human decision.

Phase 1 accepts written human approval. Cryptographic verification is not implemented here; do not claim that a commit author proves a signature. Its design and implementation remain (TBD) for a separately approved phase.
