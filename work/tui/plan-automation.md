# The automation — the skills, and the Sous-Chef tab

> **This is the design.** Nothing is built. Each slice's brief is written *after*
> the chef rules on the questions in § 10, and the brief wins over this file where
> they differ.
>
> The design the T1–T7 round followed is [plan.md](plan.md); the round's state is
> [README.md](README.md). This is the first work **out** of that round — the
> orchestrator's tab was explicitly *not in this plan*, and it comes back now
> because the orchestrator does.
>
> **It is two slices, not one.** The house rule is one slice per session: two
> slices in one session means two sets of context, and the second is the one that
> breaks. The split is in § 11.

## The answer

A tab that **starts the automation** and then **shows it working**.

| | |
|---|---|
| **The name** | **Sous-Chef** — ruled 2026-10-31. The tab's old name: the automation's tab, when it was a placeholder that said so plainly |
| **The place** | **between Board and Log** — ruled 2026-10-31. **8 tabs**: Home · Recipes · Backlog · Ready · Board · **Sous-Chef** · Log · Settings |
| **The move** | one: **start the automation**. `space` on the tab spawns the orchestrator |
| **The role** | the spawned orchestrator **loads its own role**: `maitre-cli skill load orchestrator-role`. The skill is the role's knowledge; the startup task adds the work |
| **What it shows** | the crew — the **orchestrator**, the **runner**, the **reviewer** — each with its state and what it is on |
| **Where the state comes from** | **herdr**, the runtime. The console asks it; the plan does not carry it |
| **The blocked agent** | raises a message to the **orchestrator**, asking for directions. The orchestrator decides |
| **The domain** | **no new view, no format change.** The console's action spawns a process and writes nothing to the plan |
| **The CLI** | gains a **skills menu**: `maitre-cli skill help` · `list` · `load <name>` |
| **The design** | **the ticket** — an order, the queue under it, and a stamp at the foot. Ruled 2026-11-01, § 14 |
| **The settings** | **in the plan's folder** — `kitchen/settings.json`, owned by the **tool**. The plan's format does not move, § 15 |
| **`DomainVersion`** | **7** — the settings are the tool's file, not a plan record. Q28, Q29 |

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
| **not running** | the **order**, unsigned: the harness, the crew, the sources, the budget, and the queue | **start the automation** |
| **running** | the **stamp**: the round, the interactions, the vital signs, and each agent's state and what it is on | **stop the automation** — § 14 |

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
nothing to say about it.

**The move is offered twice**, and the footer says which: `space` **starts** the
crew when it is not running, and **stops** it when it is. That revises Q22 — see
§ 10.

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

The startup task is **short**, because the role is not in it:

> You are the orchestrator of the plan at &lt;the plan's root&gt;.
> Run `maitre-cli skill load orchestrator-role` and follow it.
> The recipes to execute are: &lt;the work&gt;.

**The skill carries the role; the startup task carries the work.** The chef's
ruling: *"it would describe what an orchestrator does and what it should know to
act as an orchestrator, aside that the information necessary to work, like what
recipe(s) should be executed."*

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
| **no format change** | `core.DomainVersion` stays **7**. The settings are the **tool's** file, not a plan record, so the plan's format does not move and the repo's own plan is not stamped. Q28, Q29 |
| **the console gains a dependency** | the `herdr` binary, on the PATH. It is a runtime dependency, not a module — `internal/ui`'s budget is unchanged |
| **the console's new action** | spawns a process. It writes **nothing** to the plan — Q20. The *settings* it writes are the tool's own file, not a plan record |
| **`internal/core`** | does not change. The skills are **files**, and the CLI reads them; the domain computes nothing new |
| **`internal/cli`** | gains three verbs — `skill help`, `skill list`, `skill load <name>` — and the file reading behind them. Still **stdlib-only** |

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
| **Q20** | does starting the automation leave a trace in the plan's journal? | **ruled: no.** The journal records what happened *to the plan*; the crew's state is runtime state. The orchestrator's own writes are recorded through the CLI as usual |
| **Q21** | what does the startup task say? | **ruled: it loads the role.** `maitre-cli skill load orchestrator-role`, plus the work. The skill carries the role's knowledge; the startup task carries the work |
| **Q22** | can the chef stop the automation from the tab? | **superseded by Q27.** It was *no*; the approved ticket offers a stop, and the footer names which key does what |
| **Q23** | where do the skills live? | **ruled: embedded in the binary.** The chef's reason: *"a way to guarantee that each version has the correct version of the skill."* The skill and the tool cannot drift, because the tool **carries** the skill |
| **Q24** | one slice or two? | **ruled: two** — **T8** the CLI's skills menu and the `orchestrator-role` skill; **T9** the Sous-Chef tab. The tab's startup task loads the skill, so the skill must exist first |
| **Q26** | the tab's design | **ruled: the ticket** — an order, the queue under it, and a stamp at the foot. § 14 |
| **Q27** | can the chef stop the automation from the tab? | **ruled: yes.** `space` stops the crew, and the footer says which. It revises Q22 |
| **Q28** | where do the console's settings live — the view's and the crew's? | **ruled: in the plan's folder, in a new file** — `kitchen/settings.json`. § 15 |
| **Q29** | who owns that file — the domain, or the tool? | **ruled: the tool.** The settings are not a plan *record* — the glossary's records are the machine's records *about the work* — and they carry no rules the plan was written under. So the plan's format does not move: `DomainVersion` stays **7**, no stamp, no migration. § 15 |

## 11 — The two slices

| | |
|---|---|
| **T8 — the skills** | `maitre-cli skill help` · `list` · `load <name>`, the skills' home, and the **`orchestrator-role`** skill. Nothing here touches the console |
| **T9 — the Sous-Chef tab** | the tab, the move that starts the automation, the crew's rows, and the startup task that loads the role |

**The order is forced.** The tab's startup task runs `maitre-cli skill load
orchestrator-role`, so the skill has to exist before the tab can spawn anything
that works. T8 first.

## 12 — The skills

A skill is **a file the CLI spits out**. The chef's words: *"similar to a cat on a
file and it would spit the skill content."*

| | |
|---|---|
| **`maitre-cli skill help`** | what the menu is, and what the verbs do |
| **`maitre-cli skill list`** | every skill's name, and a line about each |
| **`maitre-cli skill load <name>`** | the skill's content, on stdout |

**The skill is the role.** `orchestrator-role` is what the spawned agent runs to
become an orchestrator: what the job is, what it must know, and the protocol it
follows — brief the crew, open the lanes, wait, close the slice. The work itself
is never in the skill; it arrives in the startup task.

**The skills are embedded in the binary** — ruled 2026-10-31. They live in the
source tree and are carried by the tool at build time, so **the skill and the
tool's version cannot drift**: the orchestrator the tool spawns loads the role the
tool carries. `maitre-cli` run from anywhere lists them — no directory, no `PATH`,
no project, no file on disk. `embed` is stdlib, so `internal/cli` stays
stdlib-only and `internal/core` does not change.

**A skill is not a recipe.** A recipe is the plan's unit of work, signed off and
tasted. A skill is the machine's own instruction sheet for a job — it is not
signed off, it has no acceptance criteria, and it is not shared between plans
unless the project carries it. That is the same line the plan draws for **tags**:
*a label is not an artifact*.

**This is the file written this session.** The `mcode-herdr-flow.txt` in the
parent directory is a first draft of exactly this skill — the herdr mechanics, the
lane protocol, the receipt and nudge. It becomes `orchestrator-role` when the
skills' home is ruled (Q23).

## 13 — Not in this slice

| | |
|---|---|
| **Starting or stopping a crew member** | the console starts the **orchestrator**, and the orchestrator dispatches the crew. The console is not the one that picks the runner's next recipe |
| **The chef's own row** | the chef is not a row. The tab is about the machine |
| **The recipe rows' contents** | unchanged. The Sous-Chef tab does not re-file the plan's items |
| **The blocked agent's escalation to the chef** | the plan's own path: the recipe goes under *Needs you* and the chef signs off on the Recipes tab |
| **`tags` on the recipe** | the plan's option A, later, with the CLI |
| **The CLI's verbs** | the automation is started from the console. A CLI verb for it is later, and the CLI is the chef's own surface |

## 14 — The tab's design: the ticket

**Ruled 2026-11-01.** The mockup is
[sous-chef-mockups.html](sous-chef-mockups.html) — it carries a simulator, so the
chef can drive the configuration and watch the tab redraw. The rejected designs are
[sous-chef-mockups-archived.html](sous-chef-mockups-archived.html).

The tab is a **ticket**: an order at the top, what is queued under it, and a stamp at
the foot saying where the work stands.

```
 ⋈ maitre · kb
──────────────────────────────────────────────────────────────────────────────
  Home    Recipes    Backlog    Ready    Board  ▸ Sous-Chef    Log    Settings
  THE ORDER
 ─────────────────────────────────────────────────────────────────────────────
  harness  ▸ opencode        pi         custom
            command: opencode
  crew     ▸ 3 · orchestrator, runner, reviewer
           2 · runner, and the orchestrator reviews
           1 · one agent, and the chef reviews
  sources  ▸ ready (4)       recipes (6)     backlog (4)
  budget   3 rounds · 12 interactions
 ─────────────────────────────────────────────────────────────────────────────
  THE QUEUE
  BL05   Wire the review queue                         next
  BL06   Sharpen the evidence contract
  BL11   Wire the review queue
  BL12   Sharpen the evidence contract
 ─────────────────────────────────────────────────────────────────────────────
  THE STAMP · round 2 of 3 · 6 of 12 interactions · plan 7 · build 7 · they agree
  working  orchestrator  routing RCP04
  working  runner        BL05 · 2 of 3 criteria proven
  idle     reviewer      waiting on the runner
 space stop the crew · ←→ choose · enter apply · q quit
```

### The four controls, and what each one says

| control | says |
|---|---|
| **harness** | `opencode` · `pi` · `custom` — and the **command** it will run, so `custom` is never a mystery |
| **crew** | `3` · `2` · `1`, each with **who reviews** |
| **sources** | `ready (4)` · `recipes (6)` · `backlog (4)` — where the work comes from |
| **budget** | `3 rounds · 12 interactions` — the **cap** the chef sets |

**The count is a composition, not a number.** Three agents, one per role · two, and
the orchestrator reviews the runner · one, and **the chef reviews**. A lone agent
cannot review itself — the close gate refuses a self-reviewed recipe
(`core.SelfReviewed`) — so the one-agent crew is the one where the chef holds the
reviewer's chair.

### Policy above, progress below

The order says the budget as a **policy** — "3 rounds · 12 interactions". The stamp
says the **progress** — "round 2 of 3 · 6 of 12". Two different facts, two different
words, neither one repeated.

### The vital signs

The stamp carries the plan's domain version, this build's, and whether they agree.
When they do not, the stamp says so: the tab that starts the automation has to be able
to report that it is not speaking the plan's rules.

### The blocked agent

The stamp shows `reviewer asked the orchestrator · the orchestrator has it`. The
console gains **no move** for it (Q18): the chef watches, the orchestrator decides.

### One task at a time

Three sous-chefs is not three tasks in parallel. The roles are sequential — one works,
the others wait — and the stamp shows exactly one of them working.

## 15 — Where the settings live

**The question the chef asked**: how the configurations on this tab, and the
preferences on the Settings tab, are persisted.

The console has two sets of choices that are not the plan's work: the **view**
(theme, order, show served) and the **crew** (harness, count, sources, budget).
Today they live in memory and die with the console — `docs/05` records the theme's
choice as *for the session*.

**Ruled 2026-11-01 — Q28 and Q29.**

**Where:** in the plan's folder, in a new file.

```
kitchen/settings.json
```

**Who owns it:** the **tool**. Not the domain.

**Why the tool, and not the domain.** A **record** is the machine's record *about
the work* — the glossary's own list is criterion · evidence · completion ·
lifecycle · `depends_on`. The settings are not one: a theme, an order, a harness, a
crew size and a budget are the chef's preferences and the automation's setup, and
none of them is a fact about the work. The version fence exists for exactly one
reason — a newer tool must not write an older plan, because it would apply rules
the plan was not written under — and the settings carry **no rules**. So:

| | |
|---|---|
| **the plan's format** | does not move. `planFormatFacts()` does not change, `core.DomainVersion` stays **7**, and the repo's own plan is **not** stamped |
| **the file's own version** | the settings carry `"version": 1` — their schema has a version of its own, the way the plan's records do and the way the skill does |
| **an unknown version** | the tool **refuses** it, by name. It does not guess — the house's *silent enum fallback* lesson |
| **the domain** | does not read the file, does not write it, does not know it |
| **the two writers** | the console and the CLI are both tools, so they share one small stdlib package — `internal/settings` — holding the schema, the read, the atomic write and the version check |
| **the plan's readers** | ignore it. A file in the plan's directory that is not a plan record does not change what a plan is |

```json
{
  "version": 1,
  "view": { "theme": "Kitchen", "order": "triage", "showServed": true }
}
```

| | |
|---|---|
| **the shape** | two parts — the version, and the sections. **T9** carries the `view` section, the one with a consumer today: the console's Settings tab. **T10** adds the `crew` section — the ticket's harness, count, sources and budget — and moves the schema's own version to **2** |
| **never runtime state** | the settings hold what the chef **chose**. The crew's *state* is what **herdr** knows (§ 5). A pane id is never written |
| **it travels with the plan** | the folder is the plan's folder, so the settings travel with it — the team works the plan the same way. A checkout that wants its own keeps it out of git |
| **forward compatible, precisely** | the **version fence wins**: a file at a version this build does not know is **refused by name**, before its keys are read. Unknown keys are ignored only **within** the version this build knows. T10 moves the schema to version 2, so a T9 build refuses a T10 file |
| **atomic** | written to a temporary file and renamed, so a half-written file cannot lose the chef's choices |
| **the poll already watches it** | `core.DiskFingerprint` hashes every `.json` in the plan's folder, so a settings write is seen by the console's reload scan |

**What the tool gains**, in the slice that lands it:

| | |
|---|---|
| **a package** | `internal/settings` — stdlib-only, 0 external modules. It owns the schema, the read, the atomic write and the version check |
| **a CLI verb** | `maitre-cli settings show` — **read-only**. The console is the only writer, because the console owns the vocabulary. There is no `settings set`: a second writer is a second place that could write a value the console does not know |
| **the console** | the Settings tab's VIEW section reads and writes the record, so the chef's choices survive a restart |
