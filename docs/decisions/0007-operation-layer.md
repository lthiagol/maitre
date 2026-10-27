# 0007. The operation layer

- **Status:** ✅ accepted
- **Date:** 2026-10-25
- **Deciders:** TBD

## Context

[ADR-0003](0003-surface-split.md) requires the domain layer to expose every operation as
a **callable unit, not as a CLI handler** — otherwise the console would have to shell
out to the CLI and we would rebuild v1's version-skew bug.

That has not been built. `internal/cli` currently loads, mutates, gates and saves by
hand. That is fine while there is one writer. It stops being fine the moment `maitre`
writes, because then two surfaces implement the same orchestration and drift.

## Decision

### One type: `core.Operations`

```go
type Operations struct {
    store   *Store
    context AuditContext   // actor · session · role
}

func NewOperations(store *Store, ctx AuditContext) *Operations
```

Methods return `(result, error)`. Both surfaces construct one and call it.

### The audit context lives in the struct, not in the call

`Operations` carries the actor, session and role for the lifetime of the operation set.
**Every write is stamped automatically, so a caller cannot forget to identify itself.**

This is not tidiness. CPD's L5 law says the session that produces code is never the
session that fixes findings, and that the session that plans is never the session that
reviews. Auditing that requires every write to carry its session — so the stamping has
to be *unskippable*. A per-call parameter is a parameter every call site can forget.

The alternative — methods on `Store` — conflates persistence with domain orchestration
and produces a god object. Free functions leave no home for shared context.

### `Store` stays the persistence seam

`Store` reads and writes files. It knows nothing about gates, lanes or lifecycle. The
domain rules live in `Operations` and the pure functions they call (`Transition`,
`LaneOf`, `ReviewStateOf`).

### One operation, one name

Every operation is justified by a real workflow — a human one or an agent one — and
neither surface may invent an operation the domain does not support. v1's 47 commands
were mostly *syntax* variations of a handful of operations: an alias, a bulk variant, a
`--json` path.

## Consequences

### Positive
- The CLI and the console cannot drift: they call the same code.
- Every write is attributable — actor, session and role are always present.
- The console can write without shelling out, so the version-skew bug class is closed.
- New surfaces are cheap: a new adapter, no new orchestration.

### Negative
- `Operations` is a big type. It grows with the operation set — the discipline against
  bloat is the "one operation, one name" rule, not a smaller type.
- One `Operations` per surface per session, which means the audit context is fixed for
  the life of the surface. If a session's role changes mid-flight, the context must be
  rebuilt. (Correct behaviour, mildly awkward.)

### Follow-up
- The canonical operation list lives in
  [`docs/03-domain-model.md` § 7](../03-domain-model.md#7-the-operation-set) and is
  **grown as operations are added**, not written up front.
