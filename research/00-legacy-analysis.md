# Legacy analysis — `master-plan` v1 (v1.0.0-rc6)

> Reference. Facts about the original, no prescriptions.
> Source: `/Users/thiago/Github/master-plan` @ `wip`, commit `c7fdfc7`, 2026-10-25.

## What it is

A project-management tool for coding agents, with a TUI for humans.

| | |
|---|---|
| `mp` | The *writer*. CLI for coding agents. JSON in/out. |
| `raul` | The *reader*. Read-only TUI for humans. |
| Methodology | An implementation of [CPD](https://github.com/lthiagol/crafted-prompt-development) |
| Storage | A directory of JSON files owned by `mp`. Never hand-edited. Validated on every write. |
| Distribution | Homebrew (stable + dev), `cargo install`, `mp install --harness opencode,cursor,pi` |

### Scale

| Metric | Value |
|---|---|
| Rust files / LOC | 767 / ~234,000 |
| — `mp` src / tests | 75,274 / 76,962 |
| — `raul` src / tests | 34,149 / 42,280 |
| — `mp-model` | 4,111 |
| `#[test]` functions | ~4,201 |
| CLI surface | **47 top-level commands, 265 subcommands** (312 paths) |
| Crates | 4 (`mp`, `mp-model`, `mp-oracle`, `raul`) |
| Skills | 8 (3 core, 5 catalog) — 27 files, **2,633 LOC** |
| Dogfood milestones | 248 |
| Release | v1.0.0-rc6, branches `stable` + `wip` |

> **The teaching layer is 2.6k LOC. The implementation is 234k — 89× larger.** The
> product's value proposition lives in the skills; the code carries it.

### Architecture

```
main.rs → cli/ (clap) → app/dispatch.rs → commands/*.rs (handlers)
                                        → domain modules (milestone/, step.rs, wp.rs, …)
                                        → store.rs (JSON I/O, atomic writes)
                                        → schema.rs → mini_schema.rs
                                        → assets.rs (embedded templates + schemas)
```

- **Persistence: JSON.** Migrated from TOML at milestone M92. `ARCHITECTURE.md` is stale
  and still describes the TOML layout.
- **Validation:** `mini_schema.rs` — a hand-rolled JSON Schema subset (10 keywords), so
  the binary avoids `jsonschema`'s 30-crate tree. `mp-oracle` holds the dev-only parity
  tests.
- **Two binaries.** `raul` shells out to `mp --format json`.
- **Skills:** YAML frontmatter + `SKILL.md` + `manifest.json`, deployed to harness dirs.

---

## Assets worth carrying forward

Knowledge, not code. These are the hardest-won parts of v1.

| # | Asset | Why |
|--:|---|---|
| A1 | **The lesson catalog** — `crates/mp/tests/fixtures/code-review-lessons.md`, 2,489 lines, **L1–L63**. 5 lessons encoded as machine-checked `Pattern:` blocks driven by `crates/mp/tests/code_review_patterns.rs` | A lesson enforced by a test, not remembered. **The single most valuable asset in the repo.** |
| A2 | **Single-source manifest + lint** — `stages.toml` is canonical; `scripts/mp_flow_lint.py` asserts `SKILL.md` matches (12 stages, per-stage commands, every role-binding cell) | The mechanism that prevents doc↔manifest drift |
| A3 | **Consumer-surface hygiene** — `scripts/check-consumer-surface.sh` blocks internal IDs (`M\d+`), lesson codes (`L\d+`) and dead doc paths from shipping | Enforces "skills ship to people with no access to your repo" automatically |
| A4 | **Autopilot verifier's typed detectors** — 7 role-boundary violations + 3-source cross-check (milestone JSON + `reviews.json` + `activity.json`) | Turns "prompt text is not enforcement" into enforcement |
| A5 | **Per-AC evidence contract** — evidence must contain (1) the runnable command, (2) `exit <code>`, (3) `(<passed>/<total> pass)`, in order | Mechanically rejects generic summaries |
| A6 | **L5 hand-off audit** — `mp reviews l5-check` detects `same_session_across_role_boundary`, `missing_session_identity`, `role_inversion` | Makes "the reviewer must be a different session" checkable |
| A7 | **Two-zone rule** — plan zone (`master-plan/`, all I/O via `mp`) vs code zone | Prevents hand-edits of plan files |
| A8 | **Fragment-first CRUD** — `mp milestone ac\|step\|wp show\|add\|update\|remove <id> [<fragment-id>]`; `--replace-arrays` rejected by default | Agents edit one thing at a time |
| A9 | **Atomic writes + advisory lock** — `flock` on the plan dir; temp-file + fsync + rename | Prevents corrupt plans from concurrent `mp` |
| A10 | **The pure transition table** — `MilestonePhase` + `MilestoneEvent` + `transition()`, deterministic and I/O-free | v1's strongest domain artifact |
| A11 | **Evidence-first AC doctrine** — "evidence is test output, not prose"; `--force` recorded as debt | Correct and rare; v1 just couldn't enforce it |
| A12 | **Session-start read discipline** — `status` → `next` → `path` → `inbox`, with `--fields`/`--summary` instead of piping to `jq` | Good agent ergonomics intent |

---

## Defects

### 5.1 CLI surface — 47 commands for one workflow

| Class | Duplication |
|---|---|
| Read a milestone | `show` · `show --summary` · `log` · `trace` · `spec review` · `spec diff` · `stage list` · `impact` (8) |
| List / discover | `list ×6` · `path` · `next` · `status` · `overview` · `search` · `inbox` · `graph` · `reviews pending` · `activity` · `digest` (11) |
| "Pending review?" | `reviews status` · `reviews pending` · `reviews lifecycle` · `list-pending-review` · `inbox --filter …` (5) |
| Record a problem | `challenge add` · `finding add` · `annotation create` · `question add` · `decision add` · `design-decision add` (6) |
| Record a note | `note add` · `autopilot note add` · `brief add` · `annotation create` (4) |
| Intake lifecycles | `idea`/`backlog`/`track`/`session`/`brief` all promote into the same 3 destinations (5) |
| Archive lifecycles | `milestone`/`track`/`idea` each with archive\|restore\|purge + `list archived` + `show archived` (4) |
| Status setters | `set-lifecycle` + `set-spec-status` (alias) + `set-status` (alias) (3) |
| Migration paths | `edit migrate-lifecycle` + `migrate --kinds` + `bulk set-lifecycle` (3) |
| Activity readers | `activity` · `digest` · `overview` · `autopilot tail` (4) |

Other surface defects:

- `milestone ac` is a **full duplicate tree** of `milestone criterion` (8 extra paths) —
  a separate enum variant, not a clap alias.
- `--fields` is global but honored by only **16 of 47** commands; the rest ignore it
  silently. Meanwhile `--include` silently ignores unknown values and `--type` returns
  empty results — but an unknown `--fields` path is a hard error. **Inconsistent error
  semantics.**
- `--quiet` / `--verbose` are declared global flags and **read nowhere** in the crate.
- 11 commands are entirely undocumented in `docs/mp/commands.md`.
- `docs/mp/commands.md` documents `mp milestone handoff` — **which does not exist**.
- `mp reviews sweep`'s doc comment claims `(alias: triage)` — **no alias is declared**.
- `milestone decompose` and `milestone plan` are near-duplicate scaffolders.

### 5.2 Two sources of truth, five divergent derivations

`plan.json` holds a milestone index; each `milestones/*.json` holds the authoritative
copy. Both are written. They drift.

The `lifecycle → (spec_status, execution_status)` projection is implemented **five
separate times**:

| # | Location | Bug |
|--:|---|---|
| 1–2 | `mp-model/milestone.rs` | canonical |
| 3 | `MilestoneFile::effective_execution_status` | different overlay precedence |
| 4–5 | `validate/plan.rs` `effective_spec_status` / `effective_execution_status` | **missing `"executed"`** |
| 6 | `sync.rs` `derive_index_status` | **missing `"executed"`** |

**Live gate hole:** a migrated milestone at `lifecycle=executed` with cleared legacy
fields escapes **G3, G4 and G7** entirely. Reproduced on M235 and M231 — *"systematic to
the `complete` path, not a one-off."*

No `reindex` command exists. `mp validate` warns about drift (W01/W03) but never errors
and never heals.

### 5.3 The gates don't gate

**Ordering is inverted.** `approve_milestone` is literally
`apply_spec_status(ctx, id, "ready")`. A bare draft with zero work packages and zero
steps goes `draft → approved` in one call, while decomposition is blocked behind
approval. v1's own dogfood log: *"ordering deadlock: decomposition requires approval,
approval should require decomposition."*

**Post-approval edits are invisible.** An approved milestone can be rewritten (title,
problem, scope, ACs, steps) and `needs_regrooming` stays `false`, with no warning. The
field, the transition event, and the raul renderer all exist — **no CLI path ever sets
it**. The docs already describe the feature as working.

**G8 over-gates.** The dependency gate guards *every* `execution_status` transition, not
just terminal ones. 5/6 steps done + 5/6 ACs passed still reads `planned`. There is no way
to express "blocked on a dependency but partly executed."

Other dead or misfiring gates:

- `require_min_acceptance_criteria` is a documented config knob — **never read**.
- Schema validation **strips `steps` and `work_packages`** before validating. The
  `$defs/step` and `$defs/work_package` schemas are dead code.
- Step status: **three sources, three answers** — Rust
  `pending|in-progress|done|skipped|failed`, schema `pending|in-progress|done|skipped`,
  docs `pending|in-progress|done|failed`.
- `risk`: Rust `"medium"`, schema `"med"` — reconciled by rewriting the doc before
  validation, so disk and schema disagree.
- AC status: code writes `passed`/`failed`; `docs/milestones/anatomy.md` documents
  `pass`/`fail`.
- `LIFECYCLE_STATES` (9 strings) vs `MilestonePhase` (7 variants) — `self-reviewed` and
  `reviewed` exist as strings no transition ever produces.
- `TransitionContext::default()` is `gates_satisfied: true` and every production caller
  uses the default → **the guard at the top of `transition()` is dead code.**

### 5.4 Verification is the weakest link

| Failure | Consequence |
|---|---|
| A `nextest` filter matching **zero** tests | A green gate that proves nothing — *"worse than no gate, because it manufactures confidence"* |
| `&&`, `;`, newlines rejected in `verification` | Every AC must be one command → v1 shipped **one wrapper shell script per AC** (M232/M233/M241) |
| Default **300 s** deadline | Full-suite ACs time out; the message (*"verification command did not run"*) is indistinguishable from a real failure |
| Verification pointing at `/tmp/verify_*.sh` | Unrunnable at close-out (exit 127) |
| Same command in an AC's `verification` and a step's `tests` | The suite runs **twice** in one completion; fixing the AC alone doesn't unblock the gate |
| Verification greps `--help` output | M116 AC-02 "verified" a flag existed without exercising the parser |
| No lint for "does this filter match ≥1 test?" | Proposed in the dogfood log; **never built** |

### 5.5 Prompt text is not a control plane

| Failure | Evidence |
|---|---|
| Runner called `mp reviews pass` on its own work | M203 violation |
| "Notify the orchestrator" executed as *printing* the message | M200 cycle 1 |
| Lifecycle claim **fabricated** with no backing | M201 cycle 1 |
| Findings shipped with placeholder strings | M202 cycle 1 |
| Typed primitives implemented but never wired to the production hot path | M215, M222, M225, M226, M229 — *"the runner's discipline is 'implement the spec' but the spec never says 'wire to production'"* |
| Evidence clobbered by a close-out ceremony | M200 F-12, M201 F-05, M202 F-13 — three milestones, three occurrences |

v1's own conclusion (`decisions.json` D-40/L7): *"Prompt-text enforcement was already
tried and failed … the runner will complete it by any available means, including violating
role boundaries."* The proposed fix — a **permission DENY list** — never shipped for the
non-autopilot path.

### 5.6 Docs and CLI are one system with no lock

| Drift | Detail |
|---|---|
| `stages.toml` vs `SKILL.md` | **7 of 12** Name cells drifted *while the lint stayed green* — though the manifest's header claims the lint fails on divergence |
| The scratch recipe | 3 files ship a broken version (`mp scratch new` prints JSON; the recipe expects a bare path) |
| `AGENTS.md` | cross-references `MP-COMMANDS.md#mp-scratch` — a **section that doesn't exist** — and claims `list`/`drop` subcommands that don't ship |
| `AGENTS.md` → `docs/code-review-lessons.md` | path **doesn't exist** (it's under `crates/mp/tests/fixtures/`) |
| `schemas/milestone.schema.json` → `docs/IDS.md` | path **doesn't exist** |
| Tap docs | advertise `mp-dev`/`raul-dev` binaries **no formula installs** |
| `ARCHITECTURE.md` | describes the **pre-M92 TOML** layout; claims gates `B1–B3` (only B1 exists) |

> **Root cause:** the same claim lives in 2–4 files and nothing checks the copies. Lesson
> L60: *"Process-prompting docs live in two places, only one checked."*

### 5.7 The agent keeps tripping on the same things

| Failure | Occurrences |
|---|---|
| Close-out ceremony overwrites per-AC evidence | 3+ milestones |
| Placeholder AC verifications shipped | **~17** across 3 sessions |
| `test(/…/)` regex matching 0 tests | silent pass, masked 14+ failures |
| Milestone record lies about its own state | L58 |
| Bulk paths bypassing single-path gates | L6 |
| Dry-run skipping the validation it should preview | L7 — *"Dry-run was a liar"* |
| `len()+1` ID allocation colliding after a removal | L56 → `[AC-01, AC-03, AC-03]` |
| `max()` derivation masking a wrong input | L39 — reported a milestone **terminal `complete`, skipping the entire review phase** |
| Collation fixed at one site, 8 sibling sort sites untouched | L34 |
| Serde-default sentinel colliding with real "empty" | L35 |
| Docstring lists fields, implementation handles a subset | L26 |
| Wire-format key drop silently zeroing consumers | L48 |
| Safe unknown-value fallback silently disabling a whole feature | L63 — auto-remediation 100% record-only, **with green unit tests** |

### 5.8 Release engineering is a separate product with its own bugs

- The Homebrew **dev formula took the tag tarball's sha256** — a *standing* defect: every
  tag bump writes a guaranteed-wrong checksum.
- Tap CI watched `stable` while the default branch is `main` → `ruby -c` and `brew style`
  had **not gated a single tap commit** since the rename.
- The stable formula's `version` was stale on **every RC bump**.
- Parallel `cargo test` races on the workspace lock → **~13/30 intermittent failures**;
  the canonical workaround is `--test-threads=1`.
- A new clippy lint sat red across **4 milestones** before being closed.

### 5.9 Smaller but structural

- `plan.json` is **239 KB** (248 index entries + a full handoff snapshot); `reviews.json`
  is 226 KB; `backlog.json` 89 KB.
- One flat file per artifact type, unbounded. 58 throwaway `.mp-scratch/` dirs at the repo
  root.
- `.gitignore` has no `__pycache__`/`*.pyc` — an orchestrator lane dropped one into the
  worktree.
- `raul` shells out to `mp` and can pick up a **stale installed binary**.

---

## Root causes

Four structural causes generate most of the defects.

### RC1 — Redundant state, redundantly derived

`lifecycle`, `spec_status`, `execution_status`, the `plan.json` index and `flow_stages`
(a second 12-stage state machine) all encode overlapping state. Every derivation is a
chance to diverge.

**Fix direction:** derive, never store.

### RC2 — The teach/enforce boundary is in the wrong place

`SKILL.md` *teaches*; the CLI *enforces*. But hard invariants (role separation, "wire to
production", "evidence is not prose") lived in the teaching layer, where they are
advisory. The autopilot path partially fixed this with typed detectors — but only there.

**Fix direction:** every invariant an agent must obey needs a machine check.

### RC3 — N×N surfaces over one small domain

8 read paths, 11 list paths, 6 problem-recording paths, 5 note paths, 4 archive schemes,
3 status setters. Each is a place to add a flag, and a place for docs to drift.

**Fix direction:** one verb per intent; a small, closed command set.

### RC4 — The model and the CLI are two contracts

`CREATE_MILESTONE_KEYS` omits fields the struct carries. `step add` lacks
`--depends-on-steps`. `bulk` can't set effort/risk — *"the only remedy is a forbidden
hand-edit."* `list milestones` omits `priority`, so the TUI's sort keybinding **cannot
work at all**.

**Fix direction:** the CLI surface is generated from, or checked against, the model.

---

## The question this raises

v1 is at v1.0.0-rc6: 4,201 tests, a rigorous skills contract, an embedded schema
validator, atomic writes, a real TUI, three-agent orchestration, and a 63-lesson
institutional memory.

**And its central promise still leaks:** an AC can be marked `passed` with evidence that
proves nothing, and a milestone can reach `complete` having skipped the review entirely.

> **What is the smallest system that makes the core promise impossible to violate?**

Everything else — language, binary count, TUI, orchestration, distribution — follows from
that answer.

---

## Files worth reading in the original

| File | What it teaches |
|---|---|
| `crates/mp/tests/fixtures/code-review-lessons.md` | L1–L63 — **read first** |
| `mp-dogfood-log.md` | 20 real findings with verdicts (`wontfix\|backlog\|spec-gap\|bug`) |
| `master-plan/decisions.json` | D-19…D-44 — the autopilot design decisions |
| `feedback/lazybrew/mp-adoption-plan-review.md` | An external adopter's migration review |
| `templates/skills/mp-flow/SKILL.md` | The 12-stage contract and 4 hand-off points |
| `templates/skills/mp-flow/stages.toml` | The canonical manifest |
| `crates/mp/src/autopilot/verifier.rs` | The 7 typed detectors — the enforcement model that worked |
| `crates/mp-model/src/milestone.rs` | The transition table |
| `ARCHITECTURE.md` | The map — **stale on persistence format** |
