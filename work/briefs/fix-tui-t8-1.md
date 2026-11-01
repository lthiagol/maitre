You are a **fresh runner** on a **fix cycle**. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/briefs/tui-t8.md](tui-t8.md) — the contract the slice was built to — and then
this brief. This brief wins where they differ. The orchestrator's pane id is `w1H:p2`.

This is **T8, fix cycle 1**. The slice landed and the review filed three blockers
and three lows. **All six are one shape: a rule that nothing holds.** You are
closing that shape, not redesigning the menu.

The findings are in `docs/audit.md` as **F-131** … **F-136** — read those rows, they
carry the proof. The reviewer's receipt is
[work/inbox/reviewer-tui-t8.md](../inbox/reviewer-tui-t8.md).

## The one idea in this cycle

**A rule whose only witness is another rule is a rule that can go when the other one
moves.** Every finding here is a claim that is true today and held by nothing: the
reviewer's mutations put the false forms back and **no test turned red**.

The fix is not to write the sentences better. It is to **make each claim say itself
where it lives, and to give every claim a test that fails when the claim is false.**

## F-131 and F-133 — the text commands

Three commands print **text**: `help`, `skill help` and `skill load <name>`. The
claim lives in three places, and they do not agree:

| | says |
|---|---|
| `AGENTS.md:133` | *"Every command prints one JSON document"* — **no exception at all** |
| `docs/04-operations.md:182` | *"Two commands print text instead: `help`, and `skill load <name>`"* — **undercounts**; three print text |
| `internal/cli/cli.go:225-227` | the **shipped binary's own help** carries the same wrong count |

**The target behaviour:** the three copies say the **same thing**, and the thing is
**true**. The tool must not miscount its own surface in the place an agent reads
first — `maitre-cli help` is the first thing anyone reads.

**And it must be held.** A test asserts that the commands which print text are
**exactly** the set the docs and the help name — one source of truth, three copies
checked against it. Put the wrong count back and the test must go red.

## F-132 — one verb per operation

`AGENTS.md:131` and `docs/04-operations.md:159` both say *"one verb per operation,
no aliases"*. Before the slice the 36 verbs mapped one-for-one onto the 36
operations. The slice added **three verbs with no operation behind them** — they
read an `embed.FS` and never touch the plan.

**The target behaviour:** the sentence states its exception, the way the
text-printing exception is stated: the `skill` group is a **read-only surface over
the tool's own embedded files**, and it is not an operation. Both copies say it.

**And it must be held.** A test asserts that **every verb outside the `skill` group
maps to an operation**, so the sentence cannot rot silently as more verbs are added.

## F-134 — the nameless skill

The refusal for a skill whose frontmatter has **no name** is held by nothing. Delete
that refusal and no test turns red: the skill is still refused, but **by accident** —
the directory-agreement rule catches the same file — and the diagnostic degrades to
*"its directory is probe-noname and its frontmatter says its name is "* — an empty
name in the middle of a sentence.

**The target behaviour:** the nameless skill is refused **for its own reason**, and
that reason says what is missing. A test holds it: deleting the refusal turns the
test red.

## F-135 — the skill with no description

A skill with a `name:` and **no `description:`** is served, and the list carries
`{"name":"probe-nodesc","description":""}`. Nothing warns.

**The list is what an agent reads to route.** An entry it cannot describe is an
entry it cannot route on — which is the whole job of the description.

**The target behaviour:** a skill with no description is **refused**, with the reason
naming the missing description, and it does not appear in the list. A test holds it.

**The shipped skill is unaffected** — `orchestrator-role` has both a name and a
description.

## F-136 — the list's names

`cmdSkillList` taking each name from a **written map** instead of from the skill
turns **no** test red, because the tree holds one skill and the map happens to name
it. The same mutation with a second skill in the tree turns the test red — so the
hole closes itself at the first second skill.

**The target behaviour:** a test puts **two** skills in the tree and asserts the
list carries **both**, with each name and description taken from that skill. The
written-map mutation must turn it red.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the commands that print **text** are **exactly** `help`, `skill help` and `skill load <name>` — and `AGENTS.md`, `docs/04-operations.md` and the CLI's own help say **the same thing**. A test holds it: put the wrong count back and the test goes red |
| **AC-02** | the *"one verb per operation"* sentence states its exception — the `skill` group is a read-only surface over the tool's embedded files, not an operation — in **both** files. A test holds it: every verb outside the `skill` group maps to an operation |
| **AC-03** | a skill whose frontmatter has **no name** is refused **for its own reason**, and the reason says what is missing. A test holds it |
| **AC-04** | a skill with a **name and no description** is refused, the reason names the missing description, and it does not appear in the list. A test holds it |
| **AC-05** | the list's names and descriptions come from the skills themselves, proved with **two** skills in the tree. A written map turns the test red |
| **AC-06** | `work/tui/README.md`'s T8 row names this fix cycle's commits. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7** |
| **`internal/core`** | does not change |
| **`internal/ui`** | does not change. The console is **T9** |
| **`internal/cli`** | stays **stdlib-only** |
| **The skill's content** | the shipped `orchestrator-role` does not change. Its description already routes |
| **The menu's shape** | `skill help` · `skill list` · `skill load <name>` stay. This cycle does not add or remove a verb |
| **The dependency budget** | unchanged: 0 / 0 / ≤ 30 |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **`docs/audit.md`** | the reviewer's. The findings are closed by the cycle-2 reviewer, not by you |
| **`docs/03`, `docs/05`, `docs/06`** | not this cycle |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`work/tui/plan-automation.md`** | the orchestrator's |
| **`work/`** | the orchestrator's, except the T8 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t8-fix1.md <<'EOF'
Runner: T8 fix 1 done. Commits: <n>. The text commands: <the three, and that the three copies say the same thing>. Its witness: <the test, and what it turns red>. One verb per operation: <the exception, in both files, and the test that holds it>. The nameless skill: <the reason it now gives, and the test>. The skill with no description: <the refusal, and the test>. The list's names: <the two skills, and the test the written map turns red>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t8-fix1.md && git commit -q -m "inbox: runner reports tui-t8 fix 1"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
