# AGENTS.md

> Project root. Read this first. It is the agent entry point and the house style.

## What this project is

**maitre** — a ground-up rewrite of
[`master-plan`](https://github.com/lthiagol/master-plan) v1 (v1.0.0-rc6), built on
[CPD](https://github.com/lthiagol/crafted-prompt-development) as its foundation.

Every line is new. v1 contributes lessons and vocabulary, not code.

| | |
|---|---|
| Agent CLI | `maitre-cli` |
| Human console | `maitre` |
| Plan directory | `kitchen/` |
| Module | `github.com/lthiagol/maitre` |
| The artifact | a **recipe** (CPD: a crafted prompt) |

## The tool today

| | |
|---|---|
| `maitre-cli` | the agent CLI — 42 commands: 36 verbs in 6 groups, plus 6 workspace commands |
| `maitre` | the human console — 5 tabs, fullscreen, alt screen |

The plan lives in `kitchen/`: `recipes/`, `courses/`, `backlog/`, plus
`activity.json` — the only event log — and `project.json`.

The console's tabs: **Recipes** (grouped by what needs the chef) · **Backlog** (ideas
and small tasks) · **Path** (the execution order) · **Board** (the plan as a dependency
tree) · **Log** (the journal) · **Sous-Chef** (the automation's vital signs) ·
**Settings**.

```bash
maitre-cli demo --dir /tmp/kb && cd /tmp/kb && maitre           # see it
maitre-cli demo --chaos --dir /tmp/kb                           # the pathological plan
scripts/console-shot.py --bin ./bin/maitre --cwd /tmp/kb   # capture the screen
```

## House style

**Lead with the answer.**

| Audience | Rule |
|---|---|
| **Humans** | Short. Tables over prose. Say it once. Cut anything a reader could reconstruct. |
| **Agents** | Enough to execute, not enough to read comfortably. Exact paths, exact commands, exact names, the "why" when it changes what to do. |

**Also:**
- State the target behaviour rather than the forbidden one. A prohibition drags the
  forbidden thing into context and half-reads as an instruction to do it.
- Prefer a table to a bulleted list of prose.
- Skip `Overview`, `Summary`, `Background` and `Next steps` sections.
- Quote only when the exact wording matters.
- Code over configuration over description.

Applies to docs, code comments, commit messages, and chat replies.

Full guidance: [`skills/writing-for-agents/SKILL.md`](skills/writing-for-agents/SKILL.md).

## Premises

| # | Premise |
|---|---|
| 1 | **CPD is the foundation.** The methodology is the spec; the tool implements it. CPD gets updated as we learn. Extensions stay additive and general. |
| 2 | **Build by using.** Intent is front-loaded; implementation is discovered by building. Tests protect domain invariants. |
| 3 | **The human interface is the main part.** The agent CLI is secondary. Both matter; focus on the human one. |
| 4 | **Agents author the artifact through the CLI**, in an interactive session with the human. The console *manages*; it does not author. |
| 5 | **Every line is new.** |
| 6 | **The domain layer has zero UI imports.** The interface is replaceable. |
| 7 | **CPD keeps generic vocabulary; the tool uses specific names.** [GLOSSARY.md](GLOSSARY.md) is the single source for the mapping. |

## Current phase

**Phase 3 — build.** The domain model is designed explicitly; everything else is
discovered by using it.

| Order | Artifact | Status |
|--:|---|---|
| 1 | `docs/03-domain-model.md` | 🔄 draft |
| 2 | `docs/04-operations.md` | 🔄 draft |
| 3 | `docs/05-human-interface.md` | 🔄 draft |
| 4 | `docs/06-cli.md` | ⏳ |
| — | `internal/core` — model, lifecycle, gates, store, views, verifier | ✅ |
| — | `internal/core/operations.go` — the public verbs | ✅ slice 3a |
| — | `internal/core/seed.go` — the seeder | ✅ slice 3b |
| — | `internal/cli` — the `maitre-cli` command surface | ✅ slice 1 |
| — | `internal/ui` — the `maitre` console | ✅ slice 4 |

**Slice 1** — `maitre-cli init`, `maitre-cli recipe create|list|show|advance`, `maitre-cli gates`.
The gate table, the transition table and the atomic store.

**Slice 2** — the `maitre` board: lanes, recipe rows, detail pane, `maitre --render` for a
static snapshot. Derived views (`Lane`, `ReviewState`, `Readiness`, `EvidenceSummary`) live
in `internal/core/views.go` — computed, never stored.

**Slice 3a — the operation layer.** `core.Operations` holds the store *and* the audit
context, so every write is stamped and no caller can forget to identify itself. One
journal (`activity.json`) is the only event log; a recipe's audit holds only its authored
deviations and decisions. See [ADR-0007](docs/decisions/0007-operation-layer.md).

**Slice 3b — the seeder.** `core.Seed` builds a realistic plan *through* the operations,
so the fixture cannot drift from the model. `maitre-cli demo` seeds one; the console's tests
use the same code.

**Slice 4 — the console.** Five tabs (Recipes · Backlog · Path · Log · Settings), the
triage grouping, the served checkbox, the detail pane with `glamour` markdown, the
primary action and the refusal banner. The lower rungs — the backlog's ideas and small
tasks, with one-way promotion — came with it, so the Backlog tab has data.

**Slice 5a — the fullscreen console.** The console owns the terminal: alt screen, its
own background, and a window title. The list fills the height and scrolls to keep the
selection visible; the footer is pinned to the bottom and says what `space` will do for
the current selection. A resize re-lays out. The detail pane appears beside the list at
>=160 columns.

**Slice 6a — the console survives real data.** Two things:

- **Sanitising.** Every title, note and command the console renders was written by an
  agent. Control characters, ANSI escapes and Unicode direction overrides are stripped on
  the way in (`internal/ui/sanitize.go`). A newline in a title would corrupt the screen.
- **Adversarial fixtures.** `core.SeedStress` adds the cases that break layouts: a title
  with a newline, a title longer than the column, thirty criteria, fifty steps, twelve
  findings, an empty recipe, and text that is not ASCII. `maitre-cli demo --chaos` writes it.

**Slice 6b — the CLI.** `maitre-cli` is now a thin adapter over the operations: 42 commands
(36 verbs in 6 groups, plus 6 workspace commands), one verb per operation, no aliases. Every write is stamped with an actor, a
session and a role — the session comes from the caller (`MAITRE_CLI_SESSION`), because CPD's
L5 law is about sessions. Every command prints one JSON document; a gate refusal is
exit 2. The command surface is in
[`docs/04-operations.md`](docs/04-operations.md#the-cli-surface).

**Slice 7 — the console stops lying.** The first dogfooded recipe (RCP01) found three
places where the console's story and the domain's story diverged:

- **The triage hid the wait.** A recipe the executor reviewed itself cannot close —
  the gate refuses it — but the console filed it under *On the line* and counted
  *0 need you*. The chef was never told. Now `core.SelfReviewed` is one fact read by
  both the gate and the triage (docs/03-domain-model.md § 8.1).
- **The detail pane could not scroll.** The offset existed and no key moved it, so
  the bottom of a recipe — its deviations, its review thread — was unreachable.
- **The detail pane overflowed its height.** It counted entries, not rows, so a
  markdown paragraph that wrapped pushed the footer off the screen.

**Slice 8 — the chef's moves.** `space` is the chef's move, derived by
`core.ChefMoveOf` so the footer cannot offer a move the gates would refuse: *sign off*
· *review* · *serve*. The review prompt asks for a verdict and a reason; it never
invents one. A script driving the console writes under the `console-script` session,
so the journal shows a driven verdict apart from a given one.

**Slice 9 — the invariant.** *"A recipe is filed under **Needs you** if and only if the
console has a chef move for it."* Written down, it catches the whole class of drift
instead of one instance at a time. It found three more: a recipe on hold still offered
its move, an approved and idle recipe sat under *Needs you* with nothing to press, and
an open tasting note did too.

**Slice 10 — the console triages a tasting note.** The last exception came out. The
detail pane has a cursor over the recipe's open notes (`n`/`N`), and the note under it
wears the console's own selected-row style. `space` opens the triage prompt: **dismiss**
the note, or **comment** in its thread. **Resolve is not offered** — that is the
fixer's move, and the console must not do an agent's job ([ADR-0003](docs/decisions/0003-surface-split.md)
§4). The prompt is the review prompt grown to carry two choices and a text, not a
second framework. An open note outranks the review: it is the concrete blocker — the
close gate refuses a recipe with an open note — and a verdict can wait.

**The round trip ran.** RCP01 was authored through `maitre-cli`, executed with real evidence,
self-reviewed, and then reviewed independently in `maitre` — the first time the loop
ran on real content. See `kitchen/` for the artifact and the journal.

**Slice 11 — the console live-updates.** A change on disk shows up without the chef
pressing `r`: the console stats the plan's files once a second and reloads when one of
them has changed. No dependency, no watcher — stat'ing twenty files is microseconds,
and the plan is read only when something actually changed. The chef keeps their place
across a reload: the recipe, the scroll, the note cursor, the theme, the order, the
served checkbox and the tab. `--render` and `--keys` do not poll: they are not
interactive programs.

**Slice 12 — the Board.** The Path tab is the flat order; the Board is the same
plan as a **dependency tree** — the shape v1's `raul` had, with an execution trunk and
labelled side branches. The trunk is `core.Path`'s ready list, in the same order (one
fact, one answer); a recipe that waits on something forks at its first unmet
dependency, and its row names the dependencies it waits on. A dependency that is not in
the plan is a broken plan, and the tree says so — and so does a dependency cycle: the
tree draws what the walk cannot reach rather than dropping it, because a tree that
forgets a recipe tells the chef a story the plan does not. Derived, never stored:
`core.TreeOf`.

**Slice 13 — the Sous-Chef tab.** The orchestrator agent gets a home in the console. It is a
**placeholder**: the autopilot is not built, and the tab says so plainly — a tab that
is waiting must not look like a tab that is failing. What it shows for now is the
automation's vital signs: the domain version the plan was last written under, this
build's, whether they agree, and the last event in the journal. That is why the tab
exists before the autopilot does — the chef can already see whether the automation is
running on the same rules as the plan.

**Open gaps:** Courses — ladder rung 3 — are not built. The Sous-Chef tab is a placeholder:
the autopilot itself is not built.

**Slice 14 — the specification's table of derived views is checked.** `docs/03` § 8
lists the views the domain computes, and it had drifted in both directions: the chef's
moves gained a fourth one, the tree got no row, and two rows claimed views the domain
does not compute. Three things now hold it:

- **The code owns the set.** `core.DerivedViews` registers each view with the
  function that computes it, and the table is checked against it in both directions.
- **The names come from the runtime.** A view's function name is derived from the
  function itself, so a function that is renamed leaves the table naming something
  that is not there, and the wire says so.
- **The wire fires when a view is added, not only when the doc is edited.** Every
  exported function of `internal/core` is either a view — it carries the `// view`
  marker, which says so where the function is defined — or a declared helper. The
  write path and the boilerplate are exempt as a family, not as seventy names; the
  read-side methods are named one by one, because a model's own methods are exactly
  where the domain put a view. The judgement is **symmetric**: a function is a view if
  and only if it is marked *and* registered, so a function cannot be reclassified by
  deleting a word, and one that is judged twice fails.

**The version guard is built.** The plan records the version of the rules it was last
written under, in its journal, and **a tool writes a plan only if the plan's version
is the tool's version** — the console included. There is no upgrade path: a newer
tool does not write an older plan, because it would apply rules the plan was not
written under. A stale `maitre-cli` on the `PATH` can read the plan but cannot write
to it. `core.DomainVersion` says which rules the tool speaks;
`TestDomainShapeIsVersioned` fails when the rules change without a bump.

## Decisions

| # | Decision |
|---|---|
| [0001](docs/decisions/0001-stack.md) | Stack: **Go + Charm** (bubbletea v2, lipgloss v2, bubbles v2, glamour v2) |
| [0002](docs/decisions/0002-human-interface-first.md) | Build order: human interface first, CLI derived from it |
| [0003](docs/decisions/0003-surface-split.md) | Who writes what: artifact vs. management |
| [0004](docs/decisions/0004-naming.md) | Naming — superseded by 0009 |
| [0005](docs/decisions/0005-vocabulary.md) | Vocabulary: generic methodology, specific product — superseded in part by 0009 |
| [0006](docs/decisions/0006-verification.md) | Verification — who runs it, and what evidence is |
| [0007](docs/decisions/0007-operation-layer.md) | The operation layer |
| [0008](docs/decisions/0008-artifact-is-the-channel.md) | The artifact is the channel — no agent-to-agent messaging |
| [0009](docs/decisions/0009-vocabulary.md) | Vocabulary — the kitchen, the crew, and the machine's records; the decision of record is [docs/07-vocabulary.md](docs/07-vocabulary.md) |

## Layout

```
cmd/maitre-cli/  the agent CLI
cmd/maitre/     the human console
internal/core/  domain model · lifecycle · gates · store · views · verifier   ← stdlib only
internal/cli/   the maitre-cli command surface                                   ← stdlib only
internal/ui/    the maitre console (bubbletea v2)
skills/         vendored third-party skills, pinned to an upstream commit
templates/      our own skills (to come)
scripts/        console-shot.py — capture what the console actually shows
                stamp-plan/ — record a plan's id migration in its journal
prototypes/     throwaway design sketches
work/           the work in progress — briefs, transcripts, the state of the rename
  briefs/       the dispatch contracts, one per slice
  transcripts/  the lanes' scrollback, captured at each slice boundary
docs/           design docs
  decisions/    ADRs
  audits/       the per-cycle review reports behind the ledger
research/       reference material on CPD and v1
```

## Build commands

```bash
make build         # both binaries into ./bin
make test          # go test ./...
make lint          # gofmt check + go vet
make run           # go run ./cmd/maitre-cli
make regen-goldens # update the console's golden snapshot
```

**Try the console:**
```bash
maitre-cli demo --dir /tmp/kb && cd /tmp/kb && maitre
maitre-cli demo --chaos --dir /tmp/kb     # the pathological plan, to break the console on
maitre --render --plain              # a static snapshot
maitre --keys "down,space" --plain  # drive it from a script
```

**Charm v2 module paths are `charm.land/<name>/v2`**, not `github.com/charmbracelet/`.

```bash
go get charm.land/bubbletea/v2 charm.land/lipgloss/v2 charm.land/bubbles/v2 charm.land/glamour/v2
```

## Dependency policy

| Package | External deps |
|---|---|
| `internal/core` | **0** — stdlib only |
| `internal/cli` | **0** — stdlib only |
| `internal/ui` | ≤ 30 (currently 28) |

A dependency needs a written reason. The core keeps its zero.

Count them yourself, the way the table counts them — the external modules a package
pulls in:

```bash
go list -deps -f '{{if not .Standard}}{{.Module.Path}}{{end}}' ./internal/ui | sort -u | grep -v lthiagol/maitre | wc -l
```

## Working rules

- **Show the commit message, wait for the go-ahead, then commit.** Push when asked.
- **Keep the docs current.** A doc that contradicts the code is a bug (v1's
  `stages.toml` vs `SKILL.md` drift: 7 of 12 cells, with the lint still green).
- **Cite the lesson.** Design claims trace to a v1 lesson or a CPD principle. A claim
  with neither is a preference — label it as one.
- **Prefer a table.** Every list of three or more similar things wants to be a table.
- **Propose before restructuring** the doc tree, the decision numbering, or the phase
  order.
- **This repo is under git** since the vocabulary rename (2026-10-27). The first
  commit is the night's work as it stood; the rename lands on top of it.
- **A brief is written into `work/briefs/` before it is dispatched**, never into
  `/tmp`. `/tmp` took half the night's record once. The lanes' scrollback is captured
  into `work/transcripts/` at every slice boundary — it is bounded and rotates.
- **A lane commits per acceptance criterion, as it goes.** A lane's turn can end
  without reporting — it happened on 2026-10-27, and 38 files sat uncommitted in the
  tree. The per-criterion commit is what makes that recoverable: without it, the work
  exists only in a working tree and a pane's scrollback.
- **A lane reports twice: the file first, then the nudge.** The receipt goes in
  `work/inbox/` and is committed; the nudge is `herdr agent prompt "<ORCH_PANE>"`.
  **The receiver waits, it does not poll** — a nudge only lands while the receiver is
  idle, so a turn blocked on a shell `sleep` cancels out the one channel that carries
  the report. End the turn and let the nudge arrive; read `work/inbox/` on the way
  back in. Never re-dispatch into a busy lane to check on it: a queued prompt looks
  exactly like a stall.
- **One write per file, and commit as soon as it is coherent.** Every edit to a large
  file re-renders that file's whole diff into the terminal, and a turn that renders
  dozens of diffs runs out before the work is done. Apply a file's changes in **one**
  scripted pass, then commit it. A commit is a landing place: a turn that ends leaves
  only what was committed. (Slice D's runner lost a turn this way on 2026-10-27, with
  two documents edited and neither committed.)

## Reference material

| Doc | What it holds |
|---|---|
| [GLOSSARY.md](GLOSSARY.md) | The domain vocabulary and the CPD↔tool mapping — single source |
| [research/00-legacy-analysis.md](research/00-legacy-analysis.md) | v1: what it got right, what it got wrong, the root causes |
| [research/01-cpd-foundation.md](research/01-cpd-foundation.md) | CPD: 14 principles, 13 stages, the artifact ladder, the L5 law, gaps |
| [docs/01-vision.md](docs/01-vision.md) | Thesis, goals, non-goals, open questions |
| [docs/02-principles.md](docs/02-principles.md) | 13 design principles, each traced to a lesson |
| [docs/03-domain-model.md](docs/03-domain-model.md) | Anatomy, lifecycle, gates, plan format, operations |
| [docs/04-operations.md](docs/04-operations.md) | The operation set — what the tool can do, and who does it |
| [docs/05-human-interface.md](docs/05-human-interface.md) | The console — four jobs, four layout candidates, the lane graph |
| [docs/06-proposals.md](docs/06-proposals.md) | The design proposals — courses, `design_decisions`, brownfield |
| [docs/07-vocabulary.md](docs/07-vocabulary.md) | The vocabulary decision of record — the scope rule and the rename table |
| [docs/testing.md](docs/testing.md) | How to test the console |
| [docs/audit.md](docs/audit.md) | The audit ledger — every finding, open or closed |
| [docs/audits/](docs/audits/) | The per-cycle review reports behind the ledger's rows |
| [work/README.md](work/README.md) | The work in progress — the state of the rename, and how to resume it |
| [prototypes/README.md](prototypes/README.md) | The board layout prototypes |
| [skills/README.md](skills/README.md) | The vendored skills, why we took them, how to re-sync |

## v1 lessons that keep biting

Carry these into every design decision:

| Lesson | The trap |
|---|---|
| Zero-match verification | A green gate that proves nothing |
| Gate guard is dead code | Transition and gate were separately invocable; the guard never ran |
| `approve` was an alias | `draft → approved` in one call, with zero steps |
| Derived state, stored twice | 5 implementations of one projection; 4 were wrong |
| Index drift | A second source of truth, reconciled only by accident |
| Evidence clobbering | Close-out ceremony overwrote per-AC evidence |
| Prompt ≠ enforcement | Role separation failed because it was prose |
| Doc drift | The same claim in 2–4 files, nothing checking the copies |
| Silent enum fallback | A safe unknown value disabled a whole feature, with green tests |
| `len()+1` ID allocation | Produced `[AC-01, AC-03, AC-03]` |
| Lexical mismatch | `spec` vs `prompt` vs `milestone` — three words, no canonical choice |
| Stale vendored copies | Pinned at `main`, never re-synced, drifted by 31 and 73 lines |
