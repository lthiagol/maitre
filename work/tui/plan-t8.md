# T8 — the orchestrator's tab

> **This is the design.** Nothing is built. The slice's brief is written *after* the
> chef rules on the questions in § 8, and the brief wins over this file where they
> differ.
>
> The design the T1–T7 round followed is [plan.md](plan.md); the round's state is
> [README.md](README.md). This is the first slice **out** of that round — the
> orchestrator's tab was explicitly *not in this plan*, and it comes back now
> because the orchestrator does.

## The answer

A tab that answers one question:

> **What is the automation doing, and what does it need from the chef?**

| | |
|---|---|
| **The name** | **Q16 — open.** The candidates and the vocabulary reason are in § 5 |
| **The rows** | one row per agent — the **job**, the **state**, and **what it is on** |
| **The move** | **Q18 — open.** The default is read-only: the tab *names* the move and does not offer it |
| **The facts** | four, each with one owner — § 2 |
| **The domain** | a view, registered in `core.DerivedViews` and written in `docs/03` § 8. The wire checks both directions |
| **The source** | **deferred by the chef.** The three options and what each costs are in § 7 |
| **The count** | **Q17 — open.** The tab count is 7 today, and the invariant says it stays |

## 1 — Why the tab exists

Three facts in the record, and they point at the same gap.

| | |
|---|---|
| **The docs promise it** | `AGENTS.md`: *"the orchestrator's tab comes back when the orchestrator does"*. `docs/05`: *"The orchestrator comes back as a tab when the orchestrator does"*. `docs/07` records the automation's tab as **Home**, *"and the tab itself comes back when the orchestrator does"* |
| **The orchestrator does** | it is the session that writes the briefs, opens the lanes, waits on them and closes the slice. It is not a placeholder anymore |
| **The chef cannot see it** | the Home tab's automation block says **"the automation is not built"**, in the warning style, beside the version agreement. It is honest — and it is the console telling the chef there is nothing to look at |

The block is honest because the automation was not built. **It is built now.** A tab
that says so would be the console's way of saying *the chef can see the machine*.

## 2 — The facts, and who owns each

One fact, one answer. A number the console shows is computed in the domain, and a
fact the console reads has one owner.

| fact | owner | where it is today |
|---|---|---|
| **the agents' states** — working · idle · blocked · done | the runtime, or a record of it | nowhere. **The source is Q19 — deferred** |
| **what each agent is on** — a recipe, a slice | the plan | the journal: every event carries a `role` and an actor, and `core.AgentRoles` is `agent · sous-chef · runner · reviewer · orchestrator` |
| **the version agreement** — the plan's version, this build's, whether they agree | the plan and the build | `core.AutomationStats`, on the Home tab |
| **the last event** | the journal | `core.AutomationStats.LastEvent`, on the Home tab |

The last two are **computed** — they are a view, and they are already built. The
first two are the new part.

## 3 — The rows

| row | shows |
|---|---|
| the **job** | the agent's name in the glossary's vocabulary: `orchestrator`, `runner`, `reviewer` |
| the **state** | `working` · `idle` · `blocked` · `done` — the same words the runtime reports |
| **what it is on** | the recipe's id, or the slice's name |

A row is **selectable**: the chef moves down the crew the way they move down a
list. `enter` opens the detail pane — the agent's last report, when there is one.

**A blocked agent says what it needs.** *"waiting on the chef — RCP03 needs your
sign-off"* — and the move stays where it belongs. The tab **names** the move; it
does not offer it. The console offers no move the gates would refuse, and a move
the console cannot perform is a promise it must not make. The Recipes tab signs
RCP03 off. The crew tab says so.

## 4 — The split with Home

Home's invariant is **numbers, never items** — no row on it is selectable, `space`
does nothing, the item count is zero. A tab holds items, so the live state
belongs there and Home keeps the vital signs.

| | Home | the crew tab |
|---|---|---|
| the version agreement | ✅ | — |
| the last event | ✅ | — |
| the agents and their states | — | ✅ |
| what each agent is on | — | ✅ |
| items | none | the agents |

Home's automation block drops the *"the automation is not built"* line **when it is
true** — the block said it because the automation was not built. A line that says
"not built" over a built automation is the console lying.

## 5 — The vocabulary, and the two collisions

**The glossary is the single source** ([GLOSSARY.md](../../GLOSSARY.md),
[docs/07](../../docs/07-vocabulary.md)), and it settles two things that a first
draft gets wrong.

**"Sous-chef" and "orchestrator" are different jobs.** The sous-chef *composes
recipe cards, preps the line and runs the tasting*. The orchestrator *calls the
recipes, routes them and keeps the line moving*. The old tab was **called**
Sous-Chef and carried the automation's vital signs — the tab's name and the job
in it drifted apart, and `docs/07` records the automation's tab as **Home**. So
naming the returning tab *Sous-Chef* names the wrong job.

**"Lane" is taken.** The glossary defines a **lane** as *where a recipe sits in the
chef's workflow* — grooming, sign-off, on the line, tasting, on hold, served — and
lists *lane* among the names to avoid for the **line**, the place where the work
happens. The runner's and the reviewer's work items are **not** lanes. (The
orchestration protocol has been calling them lanes all evening; the glossary is the
record, and the tab is where the word would become a product fact.)

**The candidates** — Q16:

| name | for | against |
|---|---|---|
| **Crew** | the glossary's own word for the people and agents working the kitchen. The tab shows the orchestrator, the runner and the reviewer — the crew minus the chef | the chef is part of the crew, and the chef is not a row |
| **Orchestrator** | the job that comes back, and the name both docs already use | the tab shows three agents, and only one of them is the orchestrator |
| **Sous-Chef** | the tab's old name, and `docs/07`'s UI string | it names the wrong job — see above |

## 6 — The domain work

A view, like every other number the console shows:

- the facts are computed in `internal/core` and **registered in `core.DerivedViews`**,
  or the wire fails — the wire checks `docs/03` § 8 against the registry in **both**
  directions, and a view that is not written down is a drift waiting to happen
- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget
- `core.DomainVersion` moves **only** if the plan's format moves — which depends
  on Q19 (§ 7)
- a test holds the split: the Home tab shows numbers and no items; the crew tab
  shows items

## 7 — Where the state comes from — **Q19, deferred**

The chef deferred this: *"for now, let's focus on the design of the tab. not really
building it, i'll give instructions after these questions."* The design holds
either way — the tab's **shape** is fixed, the **transport** is not.

| option | the console reads | costs |
|---|---|---|
| **the plan** — the orchestrator writes the crew's state into the plan; the console reads only the plan | the plan directory, as it does now | the plan's format moves: a new file → **`DomainVersion` 7 → 8**, and the repo's own plan is stamped forward. The state is as fresh as the last write |
| **the runtime** — the console asks for the live state | the runtime, outside the plan | the console gains a dependency on something that is not the plan. A plan read from disk says nothing about the crew, and the console's one rule — it reads the plan — is broken |
| **both** | the plan's record and the runtime's live state | two sources for one fact, and the console must reconcile them |

**The recommendation, for whenever the chef rules**: the **plan**. The console
keeps its one dependency, the crew becomes part of the artifact, and the plan is
the channel (ADR-0008). The cost is a format change and a version bump — which
the plan's own no-migration ruling makes cheap.

## 8 — The questions for the chef

| # | Question | Recommendation |
|--:|---|---|
| **Q16** | the tab's name | **Crew** — the glossary's word for the people and agents working the kitchen, and the tab shows the crew minus the chef |
| **Q17** | the tab count and where it sits | **8 tabs**: Home · Recipes · Backlog · Ready · Board · **Crew** · Log · Settings. The crew sits with the plan's tabs, before the records and the settings |
| **Q18** | does the tab offer the chef a move? | **no — it names the move.** A blocked agent says what it needs and where the move lives. A move the console cannot perform is a promise it must not make |
| **Q19** | where the state comes from | deferred — § 7. The recommendation is **the plan**, when the chef rules |

## 9 — Not in this slice

| | |
|---|---|
| **Starting or stopping a crew member** | the console manages the plan; it does not run the automation. A key that starts a lane is a new domain operation, and it belongs to the CLI |
| **The chef's own row** | the chef is not a row on the crew tab. The tab is about the machine |
| **The recipe rows' contents** | unchanged. The crew tab does not re-file the plan's items |
| **`tags` on the recipe** | the plan's option A, later, with the CLI. `DomainVersion` 7 → 8 on its own account |
