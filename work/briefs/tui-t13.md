You are a **fresh runner**. The repository is `/Users/thiago/Code/Pessoal/maitre`.
Read `AGENTS.md` first, then this brief. This brief wins where the design and the
docs differ. The orchestrator's pane id is `w1H:p2`.

This slice is **T13 — the surfaces**. It changes **no behaviour**. It writes down
the boundary between the three surfaces the machine already has, and holds it.

## The boundary

| | who writes it | with what |
|---|---|---|
| **the plan's records** — a recipe's lifecycle, its evidence, its review | the **crew**, through the **CLI** | `maitre-cli` verbs, each one an operation through `core.Operations`, so the gates and the audit apply |
| **the automation** — start and stop the crew, the round, the label | the **TUI** | `herdr tab create` · `herdr pane run` · `herdr tab rename` |
| **the automation's own bookkeeping** — the briefs, the receipts, the ledger | the orchestrator | the repo's `work/` and `docs/`, **not** the plan |

**One sentence, which is the slice:** *the plan's records move through the CLI; the
crew moves through herdr; the CLI never controls the automation.*

## The hole this slice closes

**The shipped skill mentions `herdr` 32 times and `maitre-cli` ZERO times.**

```
$ grep -c "maitre-cli" internal/cli/skills/orchestrator-role/SKILL.md
0
$ grep -c "herdr"        internal/cli/skills/orchestrator-role/SKILL.md
32
```

The skill's sections are: the role · the crew · the roster's labels · what you must
know · the protocol · the herdr mechanics · what you must not do · when something
goes wrong. **Not one of them names the CLI.**

And the consequence is not academic. The protocol ends *"close the slice → the
ledger"*, and the agent's only stated machinery is herdr. **How does a recipe's
lifecycle move?** The skill never says. An agent that guesses will **hand-edit the
plan's JSON** — the write the domain refuses, bypassing every gate and the audit
trail. The house's own rule names that failure: *the tool never writes a record it
will not read*.

**A boundary nobody can read is a boundary the agent guesses at.**

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the skill **names the CLI** and says **which surface moves what**: the recipe's lifecycle, its evidence and its review move through `maitre-cli`; the lanes, the labels and the round move through herdr; the automation is started from the **console**, not from the CLI |
| **AC-02** | the **CLI has no verb that controls the automation** — and a test holds it. The CLI's verbs are the plan's records and nothing else: no verb starts a crew, stops one, renames a lane's tab, or reads the crew's state. The wire reads the CLI's own dispatch table, the way the other wires in this repo do |
| **AC-03** | the skill tells the orchestrator to move a recipe's state **through the CLI** — never by editing the plan's files — and says so where the protocol says *"close the slice"* |
| **AC-04** | the **docs** say the same: `docs/03`'s architecture says which surface is which, and `AGENTS.md`'s layer table names the CLI as the artifact's surface and the console as the automation's |
| **AC-05** | `work/tui/README.md`'s T13 row names this slice's commits. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **What the CLI does** | **unchanged.** Its verbs, its output, its count — all as they are. This slice adds no verb and removes none |
| **What the console does** | **unchanged.** The tab still starts and stops the crew the way it does |
| **`core.DomainVersion`** | stays **7**. No format change |
| **`internal/core`** | does not change |
| **The goldens** | **not one moves.** The console's screens do not change |
| **The budget** | unchanged: 0 / 0 / 0 / 28 |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **The CLI's verbs** | do not add one, do not remove one. AC-02 **holds** the set; it does not change it |
| **The console** | does not change. The Sous-Chef tab is closed; its behaviour stays |
| **Courses** | not this slice. The ladder's rung 3 is the next round, and it lands on the boundary this slice writes down |
| **The runtime's protocol** | not this slice's. The console's fence against an unpinned herdr is its own small slice |
| **`docs/audit.md`** | the reviewer's |
| **`work/`** | the orchestrator's, except the T13 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t13.md <<'EOF'
Runner: T13 done. Commits: <n>. The skill: <what it now says about the CLI, the herdr and the console, and where in it>. The wire: <how the test holds that the CLI has no verb that controls the automation>. The recipe's state: <what the skill says moves it, and through what>. The docs: <what docs/03 and AGENTS.md now say>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t13.md && git commit -q -m "inbox: runner reports tui-t13"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
