# T8 — the Sous-Chef tab

> **This is the design.** Nothing is built. The slice's brief is written *after* the
> chef rules on the questions in § 10, and the brief wins over this file where they
> differ.
>
> The design the T1–T7 round followed is [plan.md](plan.md); the round's state is
> [README.md](README.md). This is the first slice **out** of that round — the
> orchestrator's tab was explicitly *not in this plan*, and it comes back now
> because the orchestrator does.

## The answer

A tab that **starts the automation** and then **shows it working**.

| | |
|---|---|
| **The name** | **Sous-Chef** — ruled 2026-10-31. The tab's old name: the automation's tab, when it was a placeholder that said so plainly |
| **The place** | **between Board and Log** — ruled 2026-10-31. **8 tabs**: Home · Recipes · Backlog · Ready · Board · **Sous-Chef** · Log · Settings |
| **The move** | one: **start the automation**. `space` on the tab spawns the orchestrator |
| **What it shows** | the crew — the **orchestrator**, the **runner**, the **reviewer** — each with its state and what it is on |
| **Where the state comes from** | **herdr**, the runtime. The console asks it; the plan does not carry it |
| **The blocked agent** | raises a message to the **orchestrator**, asking for directions. The orchestrator decides |
| **The domain** | **no new view, no format change.** The console's action spawns a process and writes nothing to the plan — unless **Q20** says otherwise |
| **`DomainVersion`** | **7** — unless Q20 says otherwise |

## 1 — Why the tab exists

| | |
|---|---|
| **The docs promise it** | `AGENTS.md`: *"the orchestrator's tab comes back when the orchestrator does"*. `docs/05`: *"The orchestrator comes back as a tab when the orchestrator does"*. `docs/07` records the automation's tab as **Home**, *"and the tab itself comes back when the orchestrator does"* |
| **The orchestrator does** | it writes the briefs, opens the lanes, waits on them and closes the slice. It is not a placeholder anymore |
| **The chef cannot see it** | the Home tab's automation block says **"the automation is not built"**, in the warning style. It was honest. The automation is built now |

And the chef's own design (2026-10-31): *"i expect that the maitre tui spin up a new
tab on herdr, it should detect if its running or not into a herdr session and if
yes, open a new tab on the same space, if not, create a new space on herdr on the
path where the tui is running and spawn a new tab with the orchestrator, it will
take care if something is needed, in case of the runner or agent being blocked, they
should raise back a message to the orchestrator asking for directions."*

## 2 — Two states

| | the tab shows | the move |
|---|---|---|
| **not running** | that the automation is not running | **start the automation** |
| **running** | the crew, each with its state and what it is on | none — the chef watches, and the crew asks the **orchestrator** |

The footer names the move, as it does everywhere: *"space starts the automation"*.

## 3 — The move

`space` on the Sous-Chef tab **starts the automation**: it spawns the orchestrator
in a herdr tab, on the plan's own directory.

This is the console's **first action that is not a plan write**. Every other key
either moves a cursor or asks the domain to do something to the plan. This one
starts a process.

That is management, not authorship (ADR-0003): the console still never writes a
recipe, and the artifact is still the crew's work. But it is a new kind of action,
and the house's rule — *the console offers no move the gates would refuse* — has
nothing to say about it. The move is offered **once**: when the automation is not
running. Once it is running, the tab has no move, and the footer says so.

## 4 — How the console starts it

**1 — is the console inside a herdr session?** herdr exports `HERDR_ENV=1` and
`HERDR_SOCKET_PATH` into every pane it owns.

**2 — open the tab.**

| | |
|---|---|
| **inside a herdr session** | the tab joins the **chef's own workspace**, beside the chef:<br>`herdr tab create --workspace "$HERDR_WORKSPACE_ID" --label "maitre: orchestrator" --cwd <the plan's root> --no-focus` |
| **not in one** | the console makes a **workspace for the plan**:<br>`herdr workspace create --cwd <the plan's root> --label "<the project>" --no-focus`<br>then a tab in it, as above |

The cwd is **the plan's root** — the directory holding `kitchen/`. The automation
works the plan, and a tab on the wrong directory is an automation on nobody's plan.

**3 — spawn the orchestrator in that tab.**

```bash
herdr pane run <the pane> mcode-herdr "<the startup task>"
```

The startup task names the job and the path. What it says beyond that is **Q21**.

**4 — the console remembers nothing.** It does not write a file, and it does not
keep the pane id in memory. It finds the crew again by asking herdr (§ 5), so a
console that was closed and reopened still shows the running automation.

## 5 — How the console finds the crew

Every fact, one owner, one answer.

| fact | owner | where it comes from |
|---|---|---|
| **the crew exists** | herdr | `herdr agent list` — the agents whose `cwd` is **the plan's root**. A workspace's own listing carries no path, so the **agents'** cwd is the anchor |
| **each one's state** | herdr | the agent's `agent_status`: `working` · `idle` · `blocked` · `done` |
| **each one's job** | the tab's **label** | `herdr tab list`, matched by `tab_id`. The label is the crew's roster line |
| **what it is on** | the tab's label | the same line: `maitre: runner · RCP03` |

**The label is the roster, and the orchestrator writes it.** The console writes the
first one, when it spawns the orchestrator. After that the orchestrator is the one
dispatching, so the orchestrator keeps the labels current — it renames the runner's
tab when the runner picks up RCP03. `herdr tab rename <TAB_ID> <LABEL>`.

**Why the label and not the plan:** herdr knows an agent is an *mcode*, not which
job it plays. The bridge registers the pane's **state**; it has no vocabulary for
*runner*. The label is the one place the job can live without the plan's format
moving.

## 6 — When a crew member is blocked

**A blocked agent raises a message to the orchestrator, asking for directions** —
the chef's ruling, and the protocol this session ran:

| | |
|---|---|
| **the record** | the receipt, written and committed — `work/inbox/<lane>-<slice>.md`. It cannot be lost |
| **the bell** | `herdr pane send-text <the orchestrator's pane> "<the receipt>"` + `herdr pane send-keys <the orchestrator's pane> enter` |

`herdr agent prompt` does not work here — the `mcode-herdr` bridge registers a
pane's **state**, not a named session, so the agent-level path answers
`agent_not_ready`. The pane-level path is the one that works, in both directions.

**Then the orchestrator decides**: dispatch again, or escalate. And an escalation
to the chef is the plan's own path — the recipe is filed under ***Needs you*** on
the Recipes tab, and the chef signs off there. **The console gains no move for
it.** The Sous-Chef tab shows the blocked agent's state, so the chef can see that
the automation is waiting on *the orchestrator* rather than on the chef.

## 7 — The vocabulary, and the two collisions

**The glossary is the single source** ([GLOSSARY.md](../../GLOSSARY.md),
[docs/07](../../docs/07-vocabulary.md)).

**"Sous-chef" and "orchestrator" are different jobs.** The sous-chef *composes
recipe cards, preps the line and runs the tasting*. The orchestrator *calls the
recipes, routes them and keeps the line moving*. The old tab was **called**
Sous-Chef and carried the automation's vital signs — the tab's name and the job
in it had drifted apart, and `docs/07` records the automation's tab as **Home**.

**Ruled: the tab is called Sous-Chef.** The name is the **tab's** name, not a
job's. The rows carry the glossary's **job** names — `orchestrator`, `runner`,
`reviewer` — so the tab and the glossary's sous-chef stay two different things,
and the glossary's sous-chef is not something this slice builds.

**"Lane" is taken.** The glossary defines a **lane** as *where a recipe sits in the
chef's workflow*, and lists *lane* among the names to avoid for the **line**, the
place where the work happens. The crew's members are not lanes.

## 8 — The domain work

**Almost none, and that is the point.**

| | |
|---|---|
| **no new view** | the tab's facts are **runtime** facts, not plan facts. Nothing is registered in `core.DerivedViews`, so `docs/03` § 8 does not move |
| **no format change** | `core.DomainVersion` stays **7**. The plan's format does not move, and the repo's own plan is not stamped forward |
| **the console gains a dependency** | the `herdr` binary, on the PATH. It is a runtime dependency, not a module — `internal/ui`'s budget is unchanged |
| **the console's new action** | spawns a process. It writes **nothing** to the plan — unless **Q20** says otherwise |
| **`internal/core`, `internal/cli`** | do not change |

## 9 — Failure handling

| | |
|---|---|
| **`herdr` is not on the PATH** | the tab says *"herdr is not here"* and offers **nothing**. The tab that cannot start the automation must not look like one that can |
| **the spawn fails** | the tab says so, and the move stays. A failed start is not a running automation |
| **the crew is gone** | every agent exited. The tab says the automation is not running, and the move comes back |
| **the orchestrator is gone but the crew is not** | the tab says so. A crew with no orchestrator has nobody to ask for directions — the design's blocked-message path is broken, and the tab must not pretend otherwise |
| **the plan's root moved** | the console finds the crew by cwd. A plan that moved is a plan with no automation |

## 10 — The questions for the chef

| # | Question | Recommendation |
|--:|---|---|
| **Q16** | the tab's name | **ruled: Sous-Chef** |
| **Q17** | the count and the place | **ruled: 8 tabs, between Board and Log** |
| **Q18** | does the tab offer the chef a move? | **ruled: it starts the automation.** The blocked agent asks the **orchestrator**, not the chef |
| **Q19** | where the state comes from | **ruled by the design: herdr**, found by the agents' cwd. The plan carries no crew state |
| **Q20** | does starting the automation leave a trace in the plan's journal? | **no** — the journal records what happened *to the plan*, and the crew's state is runtime state. The orchestrator's own writes are recorded through the CLI as usual. A trace would mean a new operation and a format change |
| **Q21** | what does the startup task say? | **the job and the path** — *"you are the orchestrator of the plan at &lt;path&gt;"* — and nothing else. The orchestrator reads the plan and decides what is needed; the console does not tell it which slice to work on |

## 11 — Not in this slice

| | |
|---|---|
| **Starting or stopping a crew member** | the console starts the **orchestrator**, and the orchestrator dispatches the crew. The console is not the one that picks the runner's next recipe |
| **The chef's own row** | the chef is not a row. The tab is about the machine |
| **The recipe rows' contents** | unchanged. The Sous-Chef tab does not re-file the plan's items |
| **The blocked agent's escalation to the chef** | the plan's own path: the recipe goes under *Needs you* and the chef signs off on the Recipes tab |
| **`tags` on the recipe** | the plan's option A, later, with the CLI |
| **The CLI's verbs** | the automation is started from the console. A CLI verb for it is later, and the CLI is the chef's own surface |
