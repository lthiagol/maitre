You are the **reviewer** for **T8 fix cycle 1**. You are a fresh lane. You did not
write the fix. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[fix-tui-t8-1.md](fix-tui-t8-1.md) — the contract for this cycle — and then
[tui-t8.md](tui-t8.md), the contract the slice was built to. The orchestrator's pane
id is `w1H:p2`. Use it verbatim.

The fix commits: `706560c` (the text commands and the list that holds them),
`dc4f4e4` (one verb per operation, and the exception), `6d48c1f` (the nameless
skill), `1801999` (the skill with no description), `27483df` (the list's names, with
two skills), and the row `9b8f034`. The receipt is
[work/inbox/runner-tui-t8-fix1.md](../inbox/runner-tui-t8-fix1.md).

## This cycle is six findings

File **HIGH or MEDIUM** blockers. Do not re-review T8 from scratch. New findings
start at **F-137**. A LOW note goes in the message, not in the ledger as a finding.

| Finding | The target | The commit |
|---|---|---|
| **F-131** | `AGENTS.md` said every command prints one JSON document, with no exception | `706560c` |
| **F-133** | "Two commands print text" undercounted — three print text — and the **shipped binary's own help** carried the same wrong count | `706560c` |
| **F-132** | "one verb per operation" was false in both files, and neither stated the exception | `dc4f4e4` |
| **F-134** | the refusal for a skill with no name was held by nothing — the directory rule caught it by accident | `6d48c1f` |
| **F-135** | a skill with a name and no description was served, listed with an empty description | `1801999` |
| **F-136** | the list's names were held only while there was one skill | `27483df` |

Close all six in `docs/audit.md` against those commits, and write the slice row for
this cycle.

**The cycle's one idea:** a rule whose only witness is another rule is a rule that can
go when the other one moves. Each finding must now have **its own** witness, and
putting the false form back must turn **that** test red.

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | `AGENTS.md`, `docs/04-operations.md` and the **binary's own help** say the **same thing** about the text commands, and that thing is **true**: exactly `help`, `skill help` and `skill load <name>` print text |
| 2 | The list lives in **one place** — `internal/cli/textcommands.go` — and the help is **built from** it, while the two documents are **checked against** it. Put "Three" → "Two" in `AGENTS.md`, drop `skill help` from the list, and change the binary's help: **each** must turn a test red |
| 3 | The *"one verb per operation"* sentence states its exception in **both** files, and the test **asks the code** — it reads the dispatch table, follows each verb's calls, and finds the groups whose verbs never open the plan. The answer is `skill`, and only `skill`. Take the exception out and the test goes red |
| 4 | A skill whose frontmatter has **no name** is refused **for its own reason** — the reason names what is missing and is not the directory rule's. Deleting the refusal turns the tests red |
| 5 | A skill with a **name and no description** is refused, the reason names the missing description, and it does not appear in the list. The **shipped** `orchestrator-role` is unaffected — it has both |
| 6 | The list's names come from the skills themselves, proved with **two** skills in the tree through the **CLI surface**. A written map, and a filter on a written name, each turn the test red |
| 7 | `internal/core` and `internal/ui` did not change. `core.DomainVersion` is **7**. The budget is 0 / 0 / 28. The command count did **not** move: still 45 commands, 39 verbs in 7 groups, plus 6 workspace commands |
| 8 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The three judgement calls the fix made

**The test fixes the sentence's wording.** The documents are found by the marker
*"commands print text instead:"* and compared with the list. So the wording is a
contract between the test and the documents. Is a marker-based test the right
shape, or is it brittle — a reworded sentence turns a test red for a reason that
has nothing to do with the truth? Say which.

**The witness for the verbs is source-level.** It reads the CLI's dispatch table and
the `Runner`'s methods and follows the calls — so `step done` → `stepStatus` →
`openStore` counts as reaching the plan. What happens when a verb starts reaching
the plan, or when the dispatch table is refactored? Is this witness one that holds,
or one that has to be rewritten whenever the code moves?

**The workspace commands `init`, `demo` and `gates` also run without a plan** — but
they are not verbs, and the documents say *one verb per operation*. Is the claim's
scope right as it stands, or does the sentence need to name them?

## Mutate, run, revert

Put each finding back, and report the tests that went red:

| | The mutation |
|--:|---|
| 1 | take the exception out of `AGENTS.md` |
| 2 | change the count word in `docs/04-operations.md` |
| 3 | change what the binary's own help says |
| 4 | delete the nameless-skill refusal |
| 5 | delete the descriptionless-skill refusal |
| 6 | make `cmdSkillList` take the name from a written map |
| 7 | add a verb to the dispatch table that never opens the plan — does the exception's witness notice? |

Revert every mutation. `git status --short` at the end is empty except your ledger
commit.

## The questions

Read [work/questions.md](../questions.md). Q11–Q13 are **closed** and Q16–Q28 are
**ruled** — do not reopen them and do not answer them. Say whether anything was
newly raised.

## The scope

| | |
|---|---|
| **Implementation code** | do not edit it, except to mutate and revert |
| **`docs/audit.md`** | yours. One commit for the ledger, and the slice row for this cycle |
| **`work/`** | the orchestrator's, except your receipt |
| **The runner** | do not prompt it. Do not spawn agents |
| **`docs/03`, `docs/05`, `docs/06`** | not yours. Report a doc that contradicts the code; do not fix it |

## The verdict

| | |
|---|---|
| No HIGH or MEDIUM findings | `verdict: ok` |
| Any HIGH or MEDIUM finding | `verdict: changes-needed` |

Do not invent reasons to fail.

## Notify the orchestrator

The receipt comes first.

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-tui-t8-cycle2.md <<'EOF'
Reviewer: T8 cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-137.. or none>). F-131 and F-133: <the three copies, the one list, and the witness>. F-132: <the exception in both files, and the source-level witness>. F-134: <the nameless skill's own reason, and its witness>. F-135: <the descriptionless skill, and its witness>. F-136: <the two skills, and the witness>. The wording: <whether a marker-based test is the right shape>. The source-level witness: <whether it holds, or has to be rewritten when the code moves>. The scope: <the workspace commands, and whether the claim is right as it stands>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t8-cycle2.md && git commit -q -m "inbox: reviewer reports tui-t8 cycle 2"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
