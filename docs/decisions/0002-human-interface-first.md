# 0002. Human interface first — the build order

- **Status:** 🔄 proposed
- **Date:** 2026-10-25
- **Deciders:** TBD

## Context

### What v1 did, and why it produced what it did

v1 built the **agent CLI first**, then bolted a TUI onto it.

| v1 artifact | Size | Fate |
|---|---|---|
| `mp` CLI | 47 top-level commands, 265 subcommands | Bloated, redundant, hard to navigate |
| `raul` TUI | 34k LOC, 60 modules | A *view* over a CLI shaped by agents, not humans |
| Docs | 40+ files across 9 trees | Drifted from the CLI (7/12 cells in one table) |
| Gates | G1–G14 + T/R/B/W series | Many gates that gated the wrong thing |

The CLI was designed around **what an agent can express**. The TUI was designed around
**what the CLI happened to expose**. So the human surface inherited the agent surface's
shape — 8 ways to read a milestone, 11 to list — and became a menu of the agent's verbs
rather than a way of *managing a plan*.

That is backwards for a product whose human interface is the main part.

## Decision

### 1. The domain model is designed around how a human manages a plan

Not around what an agent can express in a flag. The primary nouns are the human's nouns:
*a board, a lane, a recipe, a review queue, a finding, a piece of evidence.* The
agent's verbs are a *projection* of the human's operations.

### 2. The human interface is built first, end to end

Before the CLI exists. The first slice: **a human can look at a plan, read a
recipe, and see its state.** Not "the CLI can create a recipe."

### 3. The CLI is derived, never invented

> Every operation must be justified by a real workflow — a human one or an agent one — and
> neither surface may invent an operation the domain doesn't support.

This is the mechanical version of P3 (*one verb per intent*) and P4 (*the model is the
contract*). v1's 47 commands were mostly *syntax* variations — an alias, a bulk variant, a
`--json` path — rather than distinct operations.

### 4. The domain layer has zero UI imports

```
core/     domain model, plan format, gates, evidence rules   ← no UI imports, ever
ui/       human interface (bubbletea)                          ← depends on core
cli/      agent-facing commands                                ← depends on core
```

The insurance policy. It makes the UI replaceable — another TUI framework, a web UI, a
native app — without touching the domain.

### 5. We learn the domain by using the interface, not by specifying it

Premise 2: **no TDD, no spec-driven development.**

| From spec-driven / TDD | From CPD |
|---|---|
| ~~A spec per feature, approved before coding~~ | Intent front-loaded; implementation discovered by building |
| ~~Tests written first, against the spec~~ | Verification as executable evidence, not as a discipline |
| ~~The plan as a second source of truth~~ | The artifact *is* the artifact; it carries the audit trail |
| ~~Gates as blockers before work starts~~ | Gates as *claims about state*, checked on write |

What we keep: **executable verification** (a proof is a command a machine ran), the **audit
trail**, **independent review** (the reviewer is never the author), **size-aware routing**
(ceremony scales with risk), and **intent front-loaded, implementation not**.

What we drop: the ceremony of specifying before building, and test-first authoring. **We
still write tests** — for the domain invariants that must never be violated, not as a TDD
ritual. UI work is validated by using it.

## Consequences

### Positive
- The human interface is a *way of working*, not a menu of CLI verbs.
- The CLI cannot outgrow the domain — v1's 47-command failure mode is structurally
  prevented.
- The domain stays UI-free, so the interface stays replaceable.
- No spec-approval queue.

### Negative
- **We will build the wrong interface first.** Churn is inevitable; the point is that it
  is cheap when the domain is isolated.
- The agent CLI is not useful until later. An adopter wanting agent-driven work in v2.0
  waits. (Mitigated: the CLI is a thin projection once the operations exist.)
- "No spec-driven development" needs a guardrail or it slides into "no design." The
  guardrail: **the domain model and the gate semantics are designed explicitly**;
  *features* are discovered by use.

### Follow-up
- The domain package's public API is the artifact to review most carefully. Everything
  else can be rewritten; that cannot.
- The canonical operation list should live in one place, and both surfaces checked
  against it — see [ADR-0003](0003-surface-split.md).

## Open questions

- ~~**Does the human interface also need to be the thing that writes?**~~ **Resolved →
  [ADR-0003](0003-surface-split.md).** The interface *does* write — decisions, triage,
  backlog — but the **artifact is authored by agents through the CLI**. The mediation
  rule survives and gets stronger.
- **How much of the CLI does an agent actually need?** Beyond artifact authoring: read
  state · claim work · report progress · record evidence · file a finding · hand off.
  A radical reduction from 47 commands — worth confirming as the target.
