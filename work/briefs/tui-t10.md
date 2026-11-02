You are a **fresh runner**. The repository is `/Users/thiago/Code/Pessoal/maitre`.
Read `AGENTS.md`, then [work/tui/plan-automation.md](../tui/plan-automation.md) § 3–§ 6,
§ 9 and § 14–§ 15 — the design — and then this brief. This brief wins where they
differ. The orchestrator's pane id is `w1H:p2`.

This slice is **T10 — the Sous-Chef tab**. It is the **last** slice of the
automation's work: the tab that starts the automation, and shows it working.

**The design is drawn.** The mockup is
[work/tui/sous-chef-mockups.html](work/tui/sous-chef-mockups.html) — it carries a
simulator, so the tab can be driven and watched. The rejected designs are in
[sous-chef-mockups-archived.html](work/tui/sous-chef-mockups-archived.html). The
mockup is the design; this brief is the contract.

**If the session cannot hold the whole slice, stop at a criterion boundary and
report.** Every criterion is one commit, so a session that ends leaves the work
committed and the rest is a brief away.

## The tab

| | |
|---|---|
| **the name** | **Sous-Chef** — the tab's name, not a job's. The rows carry the glossary's job names: **orchestrator**, **runner**, **reviewer** |
| **the place** | **between Board and Log**. **8 tabs**: Home · Recipes · Backlog · Ready · Board · **Sous-Chef** · Log · Settings |
| **the move** | one key, two meanings: `space` **starts** the crew when it is not running, and **stops** it when it is. The footer says which |

The tab count goes from **7** to **8**, so `AGENTS.md`, `docs/05` and the console's
own tab list all move with the code.

## The ticket

The tab is a **ticket**: an order at the top, what is queued under it, and a stamp at
the foot saying where the work stands.

```
 THE ORDER
  harness  ▸ opencode        pi         custom
            command: opencode
  crew     ▸ 3 · orchestrator, runner, reviewer
           2 · runner, and the orchestrator reviews
           1 · one agent, and the chef reviews
  sources  ▸ ready (4)       recipes (6)     backlog (4)
  budget   3 rounds · 12 interactions
 THE QUEUE
  BL05   Wire the review queue                         next
  ...
 THE STAMP · round 2 of 3 · 6 of 12 interactions · plan 7 · build 7 · they agree
  working  orchestrator  routing RCP04
  working  runner        BL05 · 2 of 3 criteria proven
  idle     reviewer      waiting on the runner
 space stop the crew · ←→ choose · enter apply · q quit
```

**The crew row states the consequence, not a number.** "2 · runner, and the
orchestrator reviews" · "1 · one agent, and **the chef reviews**". A lone agent
cannot review itself — the close gate refuses a self-reviewed recipe
(`core.SelfReviewed`) — so the one-agent crew is the one where the **chef** holds
the reviewer's chair.

**Policy above, progress below.** The order says the budget as a **policy** ("3
rounds · 12 interactions"); the stamp says the **progress** ("round 2 of 3 · 6 of
12"). Two different facts, two different words, neither repeated.

**The vital signs** are on the stamp: the plan's domain version, this build's, and
whether they agree. When they do not, the stamp says so.

## The settings' `crew` section

The record gains the tab's configuration. **The schema's own version moves 1 → 2**,
and the `crew` section joins `view`:

```json
{
  "version": 2,
  "view": { "theme": "Kitchen", "order": "triage", "showServed": true },
  "crew": { "harness": "opencode", "count": 3, "sources": "ready", "rounds": 3, "interactions": 12 }
}
```

**The console writes it** when the chef changes a pick, and reads it when it
starts. **The plan's format does not move** — `core.DomainVersion` stays **7**;
this is the *file's* version, not the plan's. **A T9 build refuses a T10 file by
version** — that is the fence, and it is the answer to "what does forward
compatible mean" (see `plan-automation.md` § 15).

## The herdr integration

**herdr is a runtime dependency, not a module.** The `herdr` binary on the PATH.
`internal/ui`'s budget does not move.

| fact | owner | where it comes from |
|---|---|---|
| **the crew exists** | herdr | `herdr agent list` — the agents whose `cwd` is **the plan's root**. A workspace's own listing carries no path, so the **agents'** cwd is the anchor |
| **each one's state** | herdr | the agent's `agent_status`: `working` · `idle` · `blocked` · `done` |
| **each one's job** | the tab's **label** | `herdr tab list`, matched by `tab_id`. The label is the roster line: `maitre: runner · RCP03` |

**The console writes the first label**, when it spawns the orchestrator. After that
the **orchestrator** keeps the labels current — it renames the runner's tab when
the runner picks up a recipe.

**Why the label and not the plan:** herdr knows an agent is an *mcode*, not which
job it plays. The bridge registers the pane's **state**; it has no vocabulary for
*runner*. The label is the one place the job can live without the plan's format
moving.

## The move

**`space` starts the crew.** The console opens a herdr tab on the plan's root —
in the **chef's own workspace** when the console is inside a herdr session, or in a
workspace made **for the plan** when it is not — and runs the orchestrator in it:

```bash
herdr pane run <the pane> mcode-herdr "<the startup task>"
```

**The startup task is short**, because the role is not in it:

> You are the orchestrator of the plan at <the plan's root>. Run
> `maitre-cli skill load orchestrator-role` and follow it. The recipes to execute
> are: <the work>.

**The skill carries the role; the startup task carries the work.**

**`space` stops the crew** when it is running. The tab then says the automation is
not running, and the move comes back.

## The blocked agent

**A blocked agent raises a message to the orchestrator, asking for directions.**
The console gains **no move** for it: the stamp shows that the reviewer asked and
that the orchestrator has it. The chef watches; the orchestrator decides. The
escalation to the chef is the plan's own path — the recipe goes under ***Needs
you*** on the Recipes tab, and the chef signs off there.

## The failures — each says what the tab says

| | the tab |
|---|---|
| **`herdr` is not on the PATH** | *"herdr is not here"* and offers **nothing**. A tab that cannot start the automation must not look like one that can |
| **the spawn fails** | says so, and **the move stays**. A failed start is not a running automation |
| **the crew is gone** | every agent exited: the tab says the automation is not running, and the move comes back |
| **the orchestrator is gone but the crew is not** | says so. A crew with no orchestrator has nobody to ask for directions, and the tab must not pretend otherwise |
| **the plan's root moved** | the console finds the crew by cwd. A plan that moved is a plan with no automation |

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the **8th tab** — **Sous-Chef**, between Board and Log. The tab count moves, and `AGENTS.md`, `docs/05` and the console's tab list move with the code |
| **AC-02** | the **order**: the harness (with the command it will run), the crew (**with who reviews**), the sources and the budget. The console reads them from the record and writes them when the chef changes one |
| **AC-03** | the **queue** and the **stamp**: the queued items and which is next; the round, the interactions, and the **vital signs** — the plan's version, the build's, and whether they agree |
| **AC-04** | the **crew's rows**: the console finds the crew by the agents' cwd, reads each one's state, and shows the job from the tab's **label** — `maitre: runner · RCP03` |
| **AC-05** | the record's **schema moves 1 → 2** with the `crew` section; the console reads and writes it; a file at version 1 is refused by name, and the console says so rather than guessing |
| **AC-06** | the **moves**: `space` starts the crew — the tab, the spawn, and the **startup task** that loads the role — and `space` stops it. The footer names which |
| **AC-07** | the **failures**: the five cases above, each with what the tab says and what it offers. A key that cannot do its job says so |
| **AC-08** | `work/tui/README.md`'s T10 row names this slice's commits. `docs/05` says what the tab is. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. The plan's format does not move — the *settings file's* version moves, which is a different thing |
| **`internal/core`** | does not change |
| **`internal/cli`** | does not change. The skills menu is **T8**; this slice uses it |
| **`internal/ui`**'s budget | unchanged: 28 of 30. **herdr is a runtime dependency, not a module** |
| **No new view** | nothing is registered in `core.DerivedViews`, so `docs/03` § 8 does not move |
| **The console never authors a recipe** | ADR-0003. The console starts a process; the crew does the work |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **`internal/core`** | does not change. If you think it must, **report it** — do not change it |
| **`internal/cli`** | does not change. `maitre-cli skill load orchestrator-role` is **T8** and it works |
| **The skills** | do not change. The role is T8's and it landed |
| **`docs/03`, `docs/04`, `docs/06`** | not this slice, except `docs/05` and the CLI surface if the count moves — it does not |
| **`kitchen/`** | do not rewrite the repo's plan. The settings file is the tool's and it is gitignored |
| **`docs/audit.md`** | the reviewer's |
| **`work/tui/plan-automation.md`**, **`work/`** | the orchestrator's, except the T10 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t10.md <<'EOF'
Runner: T10 done. Commits: <n>. The tab: <the name, the place, the count>. The order: <the four controls, and the crew's consequence>. The queue and the stamp: <what the queue shows, what the stamp says, and the vital signs>. The crew: <how the console finds it, the states, and the labels>. The record: <the schema's version, and the refusal>. The moves: <the start, the startup task, the stop, and what the footer says>. The failures: <each one, and what the tab says>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t10.md && git commit -q -m "inbox: runner reports tui-t10"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
