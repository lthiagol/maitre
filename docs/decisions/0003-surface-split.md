# 0003. Surface split — who writes what

- **Status:** 🔄 proposed
- **Date:** 2026-10-25
- **Deciders:** TBD

## Context

### The two surfaces

| | Human interface | Agent CLI |
|---|---|---|
| Who drives it | the human | an agent |
| What it does with the plan | **view · decide · manage the backlog · run execution** | **author the artifact** · read state · record progress and evidence |
| What it writes | decisions, triage, backlog changes | the artifact (intent, problem, scope, ACs, guidance, audit trail) |

### Where the artifact comes from

CPD stage 0 is *"Hold a chat session: explore the problem, challenge assumptions, organize
thoughts, understand goals. **No prompt yet.**"* Stages 1–4 (Draft, Groom, Specify,
Approve) are all coordinator work, driven by conversation with the human.

So the artifact is **authored by an agent, in an interactive session with the human**. The
human contributes through that conversation — challenging scope, rejecting an AC, naming
the out-of-scope — and the agent materializes it.

### The mediation rule survives

> *The plan should not be directly edited, except for edge cases — but edited/created
> using the CLI.*

v1's rule ("the plan zone is mediated; never hand-edit") is **kept**. Its purpose was
never to exclude humans; it was that the plan's invariants hold and the audit trail is
trustworthy. The agent writing through the CLI preserves both. A human editing a
recipe JSON by hand destroys both.

## Decision

### 1. The artifact is authored by agents, through the CLI

The human never hand-writes intent, problem, scope, ACs or guidance. They *direct* that
authoring in conversation; the agent materializes it via the CLI.

### 2. The human interface is a **management surface**, not an authoring surface

| Job | What it looks like |
|---|---|
| **View** | Read the plan well — the artifact as prose, the AC table, evidence, the audit trail, findings anchored to code |
| **Decide** | Approve, triage findings, accept or reject a scope, choose the next work, accept or reject a remediation |
| **Manage the backlog** | Create, edit, promote, dismiss ideas and briefs |
| **Run & execute** | Start, pause, watch, steer execution; see what agents are doing now |

**No artifact editor.** No rich structured-text format, no parse/serialize round-trip, no
merge UX.

### 3. Both surfaces are first-class writers — through the same rules

The interface is not a read-only view. When the human approves a recipe or resolves a
finding, the interface writes to the plan — through the domain layer, so every rule, gate
and invariant applies exactly as it does for the agent.

```
                    ┌──────────────────────────────┐
   human ──────────▶│                              │
                    │   core/  (domain + rules)     │────▶ plan files
   agent ──────────▶│   no UI imports              │
                    │                              │
                    └──────────────────────────────┘
                       ▲            ▲
                       │            │
              ui/ (Charm)      cli/
```

The mediation rule is now *stronger*: **there is no path to the plan that doesn't go
through the core.**

### 4. The console offers the chef's moves only

The console's primary action is the *chef's* action for the selected recipe. It signs
off, reviews, serves. It does **not** claim the work — that is an agent's move, and the
console is not an agent.

This was found the hard way: an early build let the console claim a recipe as
`runner`, which put a false claim in the artifact — the journal said `chef`, the
recipe said `runner`. Offering to do an agent's job is how a surface ends up lying
about who did what.

The chef's moves are derived (`core.ChefMoveOf`), so the footer cannot offer a move the
gates would refuse. Where a move needs a judgment — a review verdict — the console
**asks** and records the answer with the chef's name on it. A verdict the console
invented would be a verdict nobody gave.

### 5. Two operation sets, one domain

```
domain operations
├── artifact lifecycle   →  CLI only        (agent authors via conversation + CLI)
├── decision & triage    →  UI primary, CLI available
├── backlog management   →  UI primary, CLI available
├── execution control    →  UI primary, CLI available
├── state read           →  both
└── evidence & findings  →  CLI primary (the agent has the output), UI views it
```

The bloat guard becomes: **every operation must be justified by a real workflow — human
or agent — and neither surface may invent an operation the domain doesn't support.**

## Consequences

### Positive
- **The human interface gets dramatically simpler.** No artifact editor. Its weight goes
  into *reading* — where a plan manager actually needs craft — and into *deciding*.
- **The CLI's write path is the agent's critical path**, and it is narrow and
  well-defined — exactly what an agent handles well.
- **The mediation rule holds**, and gets stronger.
- **CPD's stage 0 gets a real home**: the human↔agent conversation.
- The "no structured-text format" outcome is a large, quiet saving.

### Negative
- **A human cannot fix a bad artifact directly.** They must go back to the agent. A
  deliberate cost — the honest reading of "the artifact is the agent's, mediated through
  the CLI."
- The interface's value depends heavily on **reading** quality — markdown, AC tables,
  evidence, audit trail. That is where the design effort goes.
- Two clients means two implementations of the same operations. Shared code keeps them
  consistent; the *presentation* of a decision differs by surface.

### Follow-up

- **The domain layer must expose every operation as a callable unit, not as a CLI
  handler.** Otherwise the UI will be forced to shell out to the CLI and we rebuild v1's
  version-skew bug. **This is the single most important architectural consequence of this
  decision.**
- The plan format needs a first-class notion of a **decision** (approve, triage,
  accept-scope) as distinct from a **mutation** of the artifact.

## Open questions

- **What is the artifact called?** → [ADR-0004](0004-naming.md): *recipe*.
- **Escape hatch for direct edits.** "Except for edge cases" needs a documented, narrow,
  auditable path — not a general "edit the file" mode. Candidates: a `repair` command for
  specific corrupt states; a UI "fix" action that routes through the agent.
- **Does the human need to reorder or re-prioritise the backlog?** Almost certainly —
  and that is a decision, not an artifact edit.
