# 02 — Principles

> 🔄 Draft. Each principle cites the v1 lesson that produced it. A principle with no
> citation is a preference — label it as one.

| # | Principle | Source |
|--:|---|---|
| **P1** | Derive, never store | v1 L39, index drift |
| **P2** | A mechanism, not a reminder | v1 D-40/L7, M203 |
| **P3** | One verb per intent | v1 47 commands / 265 subcommands |
| **P4** | The model is the contract | v1 `CREATE_MILESTONE_KEYS`, B-85, B-87 |
| **P5** | Evidence, not assertion | v1 zero-match gate, D-25/R7 |
| **P6** | One source for docs; lint the copies | v1 `stages.toml` drift 7/12 |
| **P7** | Closed enums, uniform errors | v1 L63, L26 |
| **P8** | The plan zone is a boundary | v1 lazybrew I-02, I-10 |
| **P9** | Fail loud, fail early, fail once | v1 L7, L6, L56 |
| **P10** | Design for the reader with no context | v1 consumer-surface lint |
| **P11** | One domain, two clients | v1 `raul` → stale `mp` |
| **P12** | The domain is UI-free | v1 version skew |
| **P13** | Build by using, not by specifying | CPD's own thesis; v1's breadth |

---

## P1 — Derive, never store

**The failure.** v1 implemented `lifecycle → (spec_status, execution_status)` in **five
places**; four didn't handle `executed`, so a migrated milestone escaped gates G3, G4 and
G7. Separately, `plan.json`'s index was a second copy of state that drifted (reproduced
twice), and `flow_stages` was a *second* 12-stage state machine over the 7-phase
lifecycle, needing four parallel derivation functions.

**Principle.** Any value computable from other state is computed at the read boundary,
never persisted. One source of truth per fact.

**Consequence.** No index file. No legacy status fields. No `flow_stages` — the stage
list *is* the lifecycle.

## P2 — A mechanism, not a reminder

**The failure.** The most-repeated failure class in v1. Role separation was enforced by
prose and failed anyway — the runner called `mp reviews pass` on its own work (M203).
v1's own conclusion (`decisions.json` D-40/L7): *"Prompt-text enforcement was already
tried and failed … the runner will complete it by any available means, including
violating role boundaries."* Also: typed primitives implemented but never wired to the
production hot path (5 milestones); "notify the orchestrator" executed as *printing* the
message (M200).

**Principle.** Every invariant an agent must obey is checked by the tool. If we can't
check it, we don't require it.

**Consequence.** v1's autopilot verifier — 7 typed role-boundary detectors, a 3-source
cross-check, a per-AC evidence contract — is the model, and it belongs in the **core**,
not only in optional orchestration.

## P3 — One verb per intent

**The failure.** v1: 47 top-level commands and 265 subcommands for one linear workflow.
8 ways to read a milestone, 11 to list, 6 to record a problem, 5 to record a note, 4
archive schemes, 3 status setters, 3 migration paths, 3 activity readers, and a full
duplicate tree (`milestone ac` vs `milestone criterion`).

**Principle.** A user intent maps to exactly one command. If two commands answer the same
question, one is a projection of the other.

**Consequence.** The command surface is a designed artifact, reviewed as a set. It fits on
a couple of screens.

## P4 — The model is the contract; the CLI is generated from it

**The failure.** `CREATE_MILESTONE_KEYS` omits fields the struct carries — a decomposed
draft can't be authored through the JSON path. `step add` lacks `--depends-on-steps` that
`step update` has. `bulk` covers priority/spec-status/lifecycle but not effort/risk, so
*"the only remedy is a forbidden hand-edit."* `list milestones` omits `priority` and
`updated`, so the TUI's sort keybinding **cannot work at all**.

**Principle.** Every field in the model is reachable from a surface. Every argument is
derivable from the model. A field that can't be set or read is a bug in one of them.

**Consequence.** Schema/code generation, or an exhaustive model↔surface coverage test.

## P5 — Evidence, not assertion

**The failure.** An AC whose `verification` filter matches **zero** tests is marked
`passed` — *"worse than no gate, because it manufactures confidence."* A
`/tmp/verify_*.sh` script is unrunnable at close-out (exit 127). A `nextest` regex
matching 0 tests is a silent pass that masked 14+ failures. A verification that greps
`--help` output "verified" a flag existed without exercising the parser.

**Principle.** An evidence must be attributable: what it ran, what it matched, what
the result was. An evidence that matched nothing is not evidence.

**Consequence.** The verifier reports match counts and rejects zero-match evidence.
Evidence payloads are structured, not free-text prose.

## P6 — One source for docs; lint the copies

**The failure.** v1's `stages.toml` vs `SKILL.md` role table drifted in **7 of 12** cells
*while the lint stayed green* — even though the manifest's header claims the lint fails
on divergence. The scratch recipe shipped broken in 3 files. `AGENTS.md` linked a
nonexistent section and a nonexistent file. `ARCHITECTURE.md` still described the
pre-M92 TOML layout. Tap docs advertised binaries no formula installed.

**Principle.** A fact written in two places is a bug waiting to happen. Either one place
generates the other, or a check fails when they disagree.

**Consequence.** v1's `mp_flow_lint.py` + `check-consumer-surface.sh` pattern carries
forward, extended to the command reference and the docs tree.

## P7 — Closed enums, uniform error semantics

**The failure.** In v1, an unknown `--fields` path is a hard error, an unknown `--include`
value is **silently ignored**, and an unknown `--type` returns an empty result with exit
0. A safe unknown-value fallback silently disabled auto-remediation **with green unit
tests** (L63). Step status had three sources with three answers; `risk` was `"medium"` in
Rust and `"med"` in the schema.

**Principle.** Every enum-typed argument validates against one closed set. An
unrecognized value is always an error — never a silent empty result.

## P8 — The plan zone is a boundary, not a convention

**The failure.** v1's two-zone rule was good but only social. It shipped a `mp config set
--plan-dir` that *wrote outside the plan directory*, and an `mp init` that emitted empty
placeholder rows adopters had to hand-edit away (the one defect class an external adopter
flagged as purely an `mp` bug).

**Principle.** The plan directory is written by exactly one writer, and the tool's own
commands cannot create unowned artifacts. If a file in the plan zone can be produced by
anything but the tool, that's a bug.

## P9 — Fail loud, fail early, fail once

**The failure.** `effective_lifecycle` used `max()` to derive state, masking a wrong input
and reporting a milestone as **terminal `complete`, skipping the entire review phase**
(L39). A dry-run path skipped the validation it was supposed to preview — *"Dry-run was a
liar"* (L7). Bulk paths bypassed single-path gates (L6). `len()+1` ID generation
collided after a removal, producing `[AC-01, AC-03, AC-03]` (L56).

**Principle.** Validate before you mutate, validate the same way on every path, and
derive IDs from a scan of what exists — never from a count.

## P10 — Design for the reader who has no context

**The failure.** v1's consumer-surface lint exists because skills ship to adopters with no
access to the repo's plan or history. The rules — no internal milestone IDs, no lesson
codes, no pointers to repo-internal files — are good, and v1 had to *enforce* them with
a ripgrep guard, which means the failure mode was "someone shipped internal IDs."

**Principle.** Shipped artifacts name the **capability**, never the history. Enforce it
automatically rather than by review.

## P11 — One domain, two clients

**The failure.** v1's `raul` shelled out to `mp --format json` — a version-skew failure mode
(a stale installed `mp` → *"unexpected argument --all"*) and a release failure mode (two
brew **formulas**, each needing its own checksum; the dev formula took the tag tarball's
sha256 on every bump).

**Principle.** Every operation the human interface performs goes through the domain
layer, so every rule, gate and invariant applies identically to both surfaces. The CLI
may expose operations the interface lacks — the agent authors the artifact through it —
but neither surface may invent an operation the domain does not support.

**Consequence.** The domain model is designed around the human's nouns (board, lane,
recipe, review queue, finding, evidence). The CLI is derived, not invented.

## P12 — The domain is UI-free

**The failure.** Same as P11. v1 had no such separation, so the UI and the CLI could
disagree about what a plan is.

**Principle.** The domain layer has zero UI imports. Every surface is a client of it.

**Consequence.** The human interface is replaceable — a different TUI framework, a local
web UI, a native app — without touching the domain. It is what makes the "extend to
web/native" option real rather than theoretical.

## P13 — Build by using, not by specifying

**The failure.** v1's difficulties came from breadth: 47 commands for one workflow, eleven
artifact types, a 239 KB index file, and a doc tree that drifted because the same claim
lived in 2–4 files. CPD itself argues that front-loading *implementation* is
counterproductive — *"implementation design is discovered by building."*

**Principle.** No spec per feature, no approval queue, no test-first authoring. Build
the smallest end-to-end slice, use it, let the next slice fall out of what we learned.

**What this does not mean.** The domain model and the gate semantics are still designed
explicitly — they are the artifact that cannot be cheaply rewritten. And invariants still
get tests; a test that protects an invariant is not a TDD ritual.

---

## Open questions

- **P1 vs. the stage list.** v1 stored `flow_stages` per milestone; P1 says derive it.
  But the 12-stage flow and the 7-phase lifecycle are different axes — is `hand-off` a
  lifecycle state, or a stage of the flow? (Resolved in
  [03-domain-model.md § 3](03-domain-model.md#3-the-lifecycle): 9 lifecycle states mapped
  to the 13 CPD stages; stage 8 is an activity, not a state.)
- **P4's implementation.** Code generation vs. a coverage test. Generation is stronger but
  adds a build step; a test is cheaper but weaker.
- **P2 vs. agent autonomy.** If the tool enforces everything, how much freedom does the
  agent have to choose its own approach? Is a rigid gate system compatible with "the agent
  follows the spec, deviates with a note"?
