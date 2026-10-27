# Decision records

One file per decision, `NNNN-short-slug.md`, numbered from `0001`.

## Should this be an ADR?

Four questions. All four must be yes.

| # | Question |
|--:|---|
| 1 | **Hard to reverse?** Would changing it later cost meaningfully? |
| 2 | **Surprising without context?** Will a reader look at the code and wonder *why on earth*? |
| 3 | **A real trade-off?** Were there genuine alternatives, picked for specific reasons? |
| 4 | **Does a reader need the reasoning, or just the decision?** |

Question 4 is the one that separates an ADR from a reference doc. If the decision is
obvious once stated, it belongs in a reference document — a glossary entry, a README
section — not here.

**Fails 1?** It is cheap to reverse: skip the ADR and just reverse it.
**Fails 2?** Nobody will wonder, so there is nothing to record.
**Fails 3?** There was no real alternative, so there is nothing to say beyond "we did the
obvious thing".

## Format

```markdown
# NNNN. <Title>

- **Status:** proposed | accepted | superseded by NNNN
- **Date:** YYYY-MM-DD
- **Deciders:** <who>

## Context
<the forces at play — constraints, evidence>

## Decision
<what we are doing>

## Consequences
### Positive
### Negative
### Follow-up
```

Length is not a virtue — a decision that needs three sections gets three sections. But an
ADR is a *record*, not a breadcrumb: it carries the reasoning and the evidence, because
the reasoning is what survives. The format comes from
[`skills/domain-modeling/ADR-FORMAT.md`](../../skills/domain-modeling/ADR-FORMAT.md).

**Say it once.** If a fact belongs to a glossary, a reference doc or the code, the ADR
points at it rather than repeating it.

## Index

| # | Title | Status |
|---|---|---|
| [0001](0001-stack.md) | Implementation stack — Go + Charm | 🔄 proposed |
| [0002](0002-human-interface-first.md) | Human interface first — the build order | 🔄 proposed |
| [0003](0003-surface-split.md) | Who writes what — artifact vs. management | 🔄 proposed |
| [0004](0004-naming.md) | Naming — project, binaries, plan directory, ladder, roles | ✅ accepted |
| [0005](0005-vocabulary.md) | Vocabulary — generic methodology, specific product | ✅ accepted |
| [0006](0006-verification.md) | Verification — who runs it, and what a proof is | ✅ accepted |
| [0007](0007-operation-layer.md) | The operation layer | ✅ accepted |
| [0008](0008-artifact-is-the-channel.md) | The artifact is the channel — no agent-to-agent messaging | ✅ accepted |
