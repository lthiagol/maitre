You are a **fresh runner**. The repository is `/Users/thiago/Code/Pessoal/maitre`.
Read `AGENTS.md`, then [work/tui/plan-automation.md](../tui/plan-automation.md) — the
design — and then this brief. This brief wins where they differ. The orchestrator's
pane id is `w1H:p2`.

This slice is **T8 — the skills**. `maitre-cli` gains a skills menu, and the roles
ship **inside the binary**.

## The verbs

| | |
|---|---|
| **`maitre-cli skill help`** | what the menu is, what the three verbs do, and that the skills ship inside the binary. **Text**, like the top-level `help` |
| **`maitre-cli skill list`** | every embedded skill: its **name** and its **description**. A JSON list, like every other list the CLI prints |
| **`maitre-cli skill load <name>`** | the skill's content, on stdout, **byte for byte as embedded**. The one command that does not print a JSON document — its output is meant to be **read**, not parsed |

## The skills are embedded

**The skill and the tool's version cannot drift.** That is the whole reason: the
orchestrator is spawned by the tool, and it loads the role the tool carries.

- the skills live in the **source tree**, under `internal/cli/skills/`, and are
  embedded at build time with Go's `embed`
- `maitre-cli` run from **anywhere** lists them — no directory, no `PATH`, no
  project, no file on disk
- `embed` is stdlib, so **`internal/cli` stays stdlib-only** and
  `internal/core` does not change at all

**One shape.** A skill is a **directory** with a `SKILL.md` inside, and the
directory's name is the skill's name:

```
internal/cli/skills/orchestrator-role/SKILL.md
```

**The frontmatter carries the name and the description**, and the directory name
must agree with it. A skill whose directory and frontmatter disagree is a
**mistake**, not a variant: a test holds it, and the CLI refuses to serve it.

**The list is computed, not written.** `skill list` reads the embedded skills and
takes each one's name and description **from the skill itself**. There is no second
list to keep in step. A test holds it: a new embedded skill appears in the list
with no other edit.

**The description is a trigger, not a title.** It says what the skill is *for* and
when an agent should load it — What, When, and Near misses. An agent reads the
description before deciding, so a description that says only "the orchestrator
skill" is a description that cannot route.

## The first skill: `orchestrator-role`

This is the skill the spawned orchestrator loads. The chef's words:

> *"it would describe what an orchestrator does and what it should know to act as
> an orchestrator, aside that the information necessary to work, like what
> recipe(s) should be executed."*

**The skill is the role. The startup task is the work.** The skill must **not**
carry which recipes to run — that arrives in the task the console hands the
orchestrator.

**What it carries:**

| | |
|---|---|
| **the role** | what an orchestrator is — the glossary's definition: *calls the recipes, routes them, and keeps the line moving* |
| **what it must know** | the repo's own rules: the brief is the contract, one commit per acceptance criterion, stage by name, `make check` is the gate, the receipt is the record, the open questions belong to the chef |
| **the protocol** | write the brief → open the lanes → dispatch the runner → wait → receive the receipt → brief the reviewer → close the slice → the ledger |
| **the herdr mechanics** | the commands that work, and **the two that do not**: `herdr agent start` has no mcode kind, and `herdr agent prompt` answers `agent_not_ready` because the bridge registers a pane's state and not a named session |
| **what it does not do** | author a recipe — the crew does that. Answer the open questions — those are the chef's |

**A first draft exists**, written in this session:
`/Users/thiago/Code/Pessoal/mcode-herdr-flow.txt`. It is a **draft, not the
skill**: it is the herdr mechanics and the lane protocol, and it says nothing
about the role itself. The skill is the role **and** the protocol. Use the draft
for the protocol; write the role from the repo's own words — `GLOSSARY.md`, the
ADR on the surface split, and the state file.

**The skill is self-contained.** An agent that loads it may know nothing about
this repo. Where a fact belongs to the plan rather than to the role — the state
file's path, the questions file — the skill **names where to look**, and says so
plainly.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | `maitre-cli skill help` prints text: what the menu is, the three verbs, and that the skills ship inside the binary |
| **AC-02** | `maitre-cli skill list` prints a JSON list of every embedded skill's **name and description**, read from the skill itself. The list is `[]`, never `null`. A test holds it: a new embedded skill appears in the list with no other edit |
| **AC-03** | `maitre-cli skill load <name>` prints the skill's content on stdout, **byte for byte as embedded**. An unknown name is refused with a message that names the skill and lists the ones that exist. A test holds the bytes |
| **AC-04** | `orchestrator-role` exists and carries the role, what the orchestrator must know, and the protocol — including the herdr commands that work and the two that do not. It does **not** carry the work: no recipe ids, no slice names |
| **AC-05** | the skills are **embedded in the binary**: `maitre-cli` run from any directory lists them, with no file, no `PATH` and no project. `internal/cli` stays **stdlib-only**. `internal/core` does not change. `core.DomainVersion` stays **7** |
| **AC-06** | `docs/04-operations.md` names the group, the three verbs, and that `skill load` prints raw. `AGENTS.md`'s command counts move to **45 commands · 39 verbs in 7 groups**, and the top-level `help` lists the new group. `work/tui/README.md`'s T8 row names this slice's commits |

One commit per criterion. The state row's own commit is not one of the ids it names.

**The command count is a wire.** `internal/cli/doc_test.go` counts the dispatch
table and checks `docs/04-operations.md` **and** `AGENTS.md` against it — every
occurrence of the count, in both files. The counts move with the code, and the
test says so if you forget.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. No format change — the skills are in the binary, not in the plan |
| **`internal/core`** | does not change at all |
| **`internal/cli`** | stdlib only. `embed` is stdlib |
| **The plan format** | does not move. The repo's own plan is not stamped forward |
| **The vendored `skills/`** | at the repo's root — the third-party skills this repo *uses*. **Do not touch it.** The tool's own roles live in `internal/cli/skills/` |
| **The console** | `internal/ui` does not change. The tab that spawns the orchestrator is **T9** |
| **`work/tui/plan.md`, `work/questions.md`** | the orchestrator's |
| **`kitchen/`** | do not rewrite the repo's plan |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **`internal/core`** | not this slice |
| **`internal/ui`** | not this slice — the tab is T9 |
| **`docs/03`, `docs/05`, `docs/06`, `docs/07`** | not this slice |
| **`docs/audit.md`** | the reviewer's |
| **`work/tui/plan.md`** | the orchestrator's |
| **`work/`** | the orchestrator's, except the T8 state row and your receipt |
| **A `maitre-cli skill run`** | the menu is **help · list · load**. Loading is the whole of it for now |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t8.md <<'EOF'
Runner: T8 done. Commits: <n>. The verbs: <help, list, load, and what each prints>. The embedding: <where the skills live, and that maitre-cli lists them from any directory>. The list: <the shape of the JSON, and where the name and the description come from>. The load: <the bytes, and what an unknown name says>. The skill: <what orchestrator-role carries, and that it carries no work>. The counts: <the command count, the groups, and what the wire said>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t8.md && git commit -q -m "inbox: runner reports tui-t8"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses
a pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
