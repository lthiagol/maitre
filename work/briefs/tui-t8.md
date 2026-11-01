You are a **fresh runner**. The repository is `/Users/thiago/Code/Pessoal/maitre`.
Read `AGENTS.md`, then [work/tui/plan-automation.md](../tui/plan-automation.md) § 12
and § 8 — the design — and then this brief. This brief wins where they differ. The
orchestrator's pane id is `w1H:p2`.

This slice is **T8 — the skills**. `maitre-cli` gains a skills menu, and the roles
ship **inside the binary**. Nothing here touches the console.

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
takes each one's name and description **from the skill itself**. There is no
second list to keep in step. A test holds it: a new embedded skill appears in the
list with no other edit.

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

**A draft of the protocol half exists**, written in this session:
`/Users/thiago/Code/Pessoal/mcode-herdr-flow.txt`. It is a **draft, not the
skill**: it is the herdr mechanics and the lane protocol, and it says nothing about
the role itself. **Use the draft for the protocol** — and *rewrite* it, because
that draft was written to be read by a human orchestrator, not by an agent that
knows nothing about this repo.

**The skill is self-contained.** An agent that loads it may know nothing about
this repo. Where a fact belongs to the plan rather than to the role — the state
file, the plan's own path, which recipes to run — the skill says where to find it,
and the startup task carries it.

**Its description must route the reader it is for.** The draft's description is
about *running a two-lane loop*; the skill's description is about *being the
orchestrator*. An agent that is about to act as the orchestrator must find this
skill, and an agent that merely wants to spawn lanes must not.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | `maitre-cli skill help` prints what the menu is, what the three verbs do, and that the skills ship inside the binary |
| **AC-02** | `maitre-cli skill list` prints a JSON list of every embedded skill — its **name** and its **description**, taken **from the skill itself**. A new embedded skill appears in the list with no other edit |
| **AC-03** | `maitre-cli skill load <name>` prints the skill's content on stdout, **byte for byte as embedded**. It is the one command that does not print a JSON document |
| **AC-04** | the `orchestrator-role` skill exists, embedded, and carries **the role and the protocol** — what an orchestrator is, what it must know, the protocol, the herdr mechanics including the two commands that do not work, and what it must not do. It does **not** carry which recipes to run |
| **AC-05** | a skill whose directory and frontmatter disagree is **refused** — the CLI will not serve it — and a test holds the rule |
| **AC-06** | the skill's description routes: it says what the skill is *for* and when an agent should load it, with **What, When and Near misses** |
| **AC-07** | `work/tui/README.md`'s T8 row names this slice's commits. `docs/05` and `docs/03` do not move. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. This slice does not touch the plan's format — the settings record is **T9** |
| **`internal/core`** | does not change at all |
| **`internal/cli`** | stays **stdlib-only**. `embed` is stdlib |
| **`internal/ui`** | does not change. This slice does not touch the console |
| **No new view** | nothing is registered in `core.DerivedViews`, so `docs/03` § 8 does not move |
| **The dependency budget** | unchanged: 0 / 0 / ≤ 30 |
| **The tab count** | stays **7** in the built console. The Sous-Chef tab is **T9** |
| **Q11, Q12, Q13** | closed. Do not reopen them |
| **Q16–Q28** | ruled. Do not re-ask them |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **`internal/ui`** | the console. **T9** owns the Sous-Chef tab |
| **`internal/core`** | the domain. The settings record is **T9** |
| **The plan's format** | `kitchen/settings.json` and the version bump are **T9** |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`docs/03`, `docs/05`, `docs/06`** | not this slice |
| **`docs/audit.md`** | the reviewer's |
| **`work/tui/plan-automation.md`** | the orchestrator's |
| **`work/`** | the orchestrator's, except the T8 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t8.md <<'EOF'
Runner: T8 done. Commits: <n>. help: <what it prints>. list: <the JSON, and that a new skill appears with no other edit>. load: <the content, byte for byte>. The skill: <the role, the protocol, the herdr mechanics, and what it does not do>. Its description: <what it says, and how it routes>. The refusal: <the skill that disagrees with its directory, and that the CLI will not serve it>. The draft: <what you took from mcode-herdr-flow.txt, and what you rewrote>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t8.md && git commit -q -m "inbox: runner reports tui-t8"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
