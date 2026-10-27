# 0006. Verification — who runs it, and what a proof is

- **Status:** ✅ accepted
- **Date:** 2026-10-25
- **Deciders:** TBD

## Context

CPD principle 9: *"Verification is auditable, not aspirational… 'The tests pass' is a
claim. `cargo test -p foo` → exit 0 is evidence. Bypass is recorded as explicit debt."*

v1 implemented the principle's *words* and missed its *point*. The agent supplied a
verification string; the tool ran it; the tool recorded the exit code; and a filter
matching **zero** tests was recorded as a pass. The dogfood log called it *"worse than no
gate, because it manufactures confidence."*

Two failure modes, one cause: the tool never asked whether the proof had matched
anything.

The decision has a trust dimension — the tool executing commands read from a plan that
an agent authored.

## Decision

### 1. The tool runs the verification

`scale verify <recipe> <AC>` executes the recipe's declared command, captures the exit
code and an output digest, derives a match count where it can, and records the proof.

The trust boundary is the **sign-off gate**: the human approves the recipe — *which
contains the verification commands* — before anything executes. By the time a proof is
recorded, the human has read the command. This is CPD's structure doing the work, and it
is why "the tool executes the plan" is safe here in a way it was not in v1 (whose
approval gate was an alias and ran nothing).

### 2. A proof has one of three match states

| State | Meaning | Remedy |
|---|---|---|
| **proven** | Match count known, and > 0 | — |
| **rejected** | Match count known, and == 0 — **the filter matched nothing** | Fix the spec. This is not a test failure. |
| **unproven** | Match count unknown — the tool cannot tell | Make the verification machine-checkable, or get a human sign-off |

Only **proven** satisfies the `ac.all_proven` gate. A rejected or unproven proof is
recorded, visible, and blocks — each with a different remedy, which is why they are
distinct states.

> v1's runner skill had this right and the tool did not: *"the runner does NOT silently
> fix bad verifications — that is a spec defect, not a runner defect."*

### 3. The match count is derived, never asserted

The tool counts matches from the output it captured. It never takes a match count on
trust.

- `go test …` → the count is known. `--- PASS:` lines, or `ok <pkg>` lines without
  `-v`. `[no tests to run]` → zero matched.
- Any other command → the count is unknown, so the proof is **unproven**.

A command the tool cannot interpret is honestly labelled, rather than silently passed.

### 4. Shell composition is allowed

The verification runs through a shell, so a verification may be a pipeline. v1 rejected
`&&` and `;`, which forced one wrapper script per acceptance criterion — a tax v1's own
log called out.

### 5. The proof is append-only

Re-running appends. It never overwrites. The criterion's status derives from the latest
proof. This is the direct fix for v1's evidence-clobbering bug, where a close-out
ceremony overwrote per-criterion evidence on three separate milestones.

## Consequences

### Positive
- A criterion cannot be `passed` on a proof that matched nothing.
- The human reads every command before it runs.
- Evidence accumulates rather than being replaced — the audit trail is trustworthy.
- Pipelines are allowed, so a verification is expressible in one line.

### Negative
- **The tool executes shell commands from the plan.** The sign-off gate is the
  mitigation, which means it must actually work — if approval is ever bypassed, this
  becomes arbitrary execution. (Tracked: the permission deny-list, below.)
- Only `go test` gives a known match count today. Other verification shapes are
  *unproven* and need a human sign-off until a resolver is written for them.
- Running a verification takes as long as the command. The default timeout must
  accommodate a realistic full suite, and a timeout must be **distinguishable from a
  failure** — v1's timeout message was identical to a real failure, which cost three
  milestones of confusion.

### Follow-up
- **The permission deny-list.** v1's runner could call `reviews pass` on its own work.
  The mechanism is a permission boundary that makes the runner physically unable to
  perform a human action. Unbuilt; it is the enforcement half of CPD's L5 law.
- **More resolvers.** Each command shape we want a machine-checkable proof for needs a
  counter. Until then it is unproven, which is honest.
- **Human sign-off for unproven criteria.** Lives in the console
  ([`docs/05-human-interface.md`](../05-human-interface.md)), not the CLI — a human
  signing off is a human action.
