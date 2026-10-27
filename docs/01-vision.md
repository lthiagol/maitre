# 01 — Vision

> 🔄 Draft. Premises and working rules: [AGENTS.md](../AGENTS.md).

## Thesis

> **What is the smallest system that makes CPD's core promise impossible to violate?**

| CPD clause | v1's leak |
|---|---|
| *Verification is auditable, not aspirational* | A `verification` string matching **zero** tests was a green gate — *"worse than no gate, because it manufactures confidence"* |
| *The lifecycle only advances when its claims are proven* | `approve` was an alias for `set-spec-status ready`; `transition()`'s gate guard was **dead code** |
| *The reviewer is never the author* | Enforced by prose; the runner called `reviews pass` on its own work |

### The division of labour

CPD principle 13: *"Code review by agents is mandatory. Code review by humans is an
extra."* Principle 14: *"The human stays in the loop at the seam between executions — not
as a bystander, and not as a micromanager, but as the validator who decides whether the
work is done or what to ask next."*

**The human is not the reviewer. The human is the manager and validator.** The interface
is where the human decides *what to build next*, *whether it was built correctly*, and
*what to ask for next* — the agents do the reviewing.

v1's TUI was read-mostly *by construction*: every write went through the CLI, so the human
surface was a viewer and the CLI was the tool. That is the central thing we fix.

## Goals

| G | Goal | How we know |
|---|---|---|
| G1 | An AC cannot be `passed` without evidence a machine produced | Every `passed` AC carries machine-attached evidence; evidence that matched nothing is rejected |
| G2 | The lifecycle cannot skip a stage | Transitions are a pure function of (state, event, evidence); guards run inside the transition |
| G3 | The human interface is the product | A human can create, groom, execute, review and close a recipe without the CLI |
| G4 | Neither surface outgrows the domain | The CLI's operation set is a subset of the domain's; neither invents an operation |
| G5 | One source of truth per fact | No derived value is stored |
| G6 | Docs can't drift from the code | The command surface and stage manifest are generated or linted |
| G7 | Failure modes are caught by the tool | Every invariant CPD teaches in prose has a machine check behind it |
| G8 | The domain layer is UI-free | `core/` has zero UI imports; the interface is replaceable |

## Non-goals

| # | Non-goal | Why |
|--:|---|---|
| N1 | Compatibility with v1's CLI or plan format | Carries the old bloat |
| N2 | Supporting every artifact v1 had | v1 had 12+ types; each is a lifecycle, a command set and a doc set |
| N3 | Multi-agent orchestration in v2.0 | It consumes the plan model; build the model first |
| N4 | A general-purpose PM tool | This is for coding agents driving code |
| N5 | A spec-driven or TDD workflow for building this | Premise 2 |

## Shape

```
core/    domain model · plan format · gates · evidence rules    ← zero UI imports
ui/      human interface (Go + Charm)
cli/     agent-facing projection of the domain's operations
skills/  what the agent is taught
docs/    CPD updates + user-facing docs
```

- One plan format. JSON, one file per artifact, **no index file**.
- Derived, never stored.
- One domain, two clients. Neither shells out to the other.
- The artifact is authored by agents through the CLI.
- TUI now, web later — the domain being UI-free makes a second shell cheap.

### Who writes what

| Operation | Primary surface |
|---|---|
| Author the artifact (intent, scope, ACs, guidance) | **CLI** — the agent, from a conversation with the human |
| View (prompt, ACs, evidence, audit trail, findings) | **UI** |
| Decide (approve, triage, accept scope, pick next) | **UI** |
| Manage the backlog (ideas, briefs, promotion) | **UI** |
| Run & execute (start, pause, watch, steer) | **UI** |
| Read state, claim work, report progress, record evidence, file a finding | **CLI** — the agent |

## Open questions

### Q1 — How much of CPD in v2.0?

| Option | v2.0 contains | Cost |
|---|---|---|
| A. Lifecycle core | Stages 1–12 + gates + evidence + audit trail + the human interface for them | Medium |
| B. Core + discovery | A + a real instrument for stage 0 (Discover), currently "free-form chat" | Medium |
| C. Core + full ladder | A + all 5 rungs + promotion/demotion | Higher |
| D. Everything + orchestration | C + multi-agent orchestration | High — v1's mistake |

**Lean: B.** Stage 0 is the most valuable stage and the least supported. CPD principle 1
says skipping discovery is *"an explicit, declared choice."* Supporting it properly is
the biggest single contribution we can make to CPD.

### Q2 — What does the human interface look like?

**Scope (answered):** view · decide · manage the backlog · run & execute. It does not
author the artifact. See [ADR-0003](decisions/0003-surface-split.md).

**Visual shape (open):** a board, a detail pane, a review queue, a backlog view, and a
discovery surface.

| Sub-question | Options |
|---|---|
| Board | kanban (columns = lanes) vs. table (rows = recipes). A table scales; a board reads |
| Detail pane | rendered markdown with inline actions vs. tabs (spec / ACs / guidance / evidence / audit). Tabs scale; inline is faster |
| Decisions | inline from the detail pane vs. a queue |
| Live updates | watch a runner working vs. on refresh only |

### Q3 — Does the human write the plan?

**Answered** → [ADR-0003](decisions/0003-surface-split.md). No. The artifact is authored
by agents through the CLI. The mediation rule survives; the UI writes decisions, triage
and backlog changes through the same domain layer.

### Q4 — How small can the agent CLI be?

Beyond artifact authoring: *read state · claim work · report progress · record evidence ·
file a finding · hand off*. v1 had 47 commands. **The most valuable question in this
list** — a CLI that small is far easier for an agent to use correctly.

### Q5 — Does the tool extend CPD, or implement it?

CPD is public and community-facing; changing it carelessly damages a public artifact.
Candidates in
[research/01-cpd-foundation.md § 7](../research/01-cpd-foundation.md#7-extension-candidates).
**Lean: extend, and early** — the tool's discoveries should land in the methodology.

### Q6 — Naming

**Answered** → [ADR-0004](decisions/0004-naming.md): backlog item → recipe → course →
menu. Project **maitre**; binaries `maitre-cli` (agent CLI) and `maitre` (human TUI). CPD's vocabulary updates to match.

### Q7 — Distribution

One brew formula. Possibly two binaries. Possibly a local web UI later.

### Q8 — Does the ladder need a decision-shaped rung?

**Deferred — build and observe.** Watch the recipe path; see whether the gap bites.

The observation that raised it: wayfinder's top rung is a map of open *questions*, each
resolved into a recorded decision. CPD's milestone is implementation-shaped — intent,
scope, ACs, steps. So "a decision we need before we can write the recipe" has no home in
our ladder; right now it would live as prose inside a backlog item.

v1 had the same gap and bolted `decision add` on as a separate plan-level list —
disconnected from the ladder, which is exactly how a thing outside the ladder always
looks. If the gap bites, the fix is a rung (a question-shaped artifact between the
backlog item and the recipe), not another side-list.

Two things worth keeping regardless of the answer:

- **Map-as-index, not map-as-store** — a decision lives in exactly one place; the map
  only gists and links. That is principle P1 (*derive, never store*) applied to a
  planning artifact, arrived at independently.
- **Refer by name, never by bare id** — a wall of `RCP07, RCP11, RCP23` is illegible. Our
  board leads with ids because v1 did. Worth revisiting when we design the board for
  real (slice 3).

## Agreement needed before architecture

1. The thesis and the human-as-validator division.
2. Q1 — how much of CPD, and whether discovery gets an instrument.
3. Q2 — the shape of the human interface.
4. Q4 — the target size of the agent CLI.
5. Q5 — how we extend CPD.
