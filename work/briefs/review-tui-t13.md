You are the **reviewer** for **T13**. You are a fresh lane. You did not write this
code. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[tui-t13.md](tui-t13.md) — the contract. The orchestrator's pane id is `w1H:p2`. Use it
verbatim.

The runner's commits: `75c4b96` (the skill's *The surfaces* section), `0e3b25b` (the
wire), `32fd028` (the recipe's state moves through the CLI), `4ac2158` (the docs), and
the row `17c8d99`. The receipt is [work/inbox/runner-tui-t13.md](../inbox/runner-tui-t13.md).

## This slice

**T13 — the surfaces.** It changes **no behaviour**. It writes down the boundary
between the three surfaces the machine already has, and holds it.

> *the plan's records move through the CLI; the crew moves through herdr; the CLI
> never controls the automation.*

## What to verify — each by command, not by reading

| # | Check |
|:--:|---|
| 1 | The skill has a section that **names `maitre-cli`**, herdr and the console, and says **which moves what**: the plan's records through the CLI, the automation from the console, the bookkeeping in `work/` and `docs/`. It went from **0** mentions of `maitre-cli` to some |
| 2 | The wire: `internal/cli/surfaces_test.go` reads the CLI's **own dispatch table** out of `cli.go`, holds that **every verb is classified**, and holds that **no verb is named after the automation**. The CLI's verb set did **not** move — 46 commands, 40 verbs in 8 groups |
| 3 | The wire **bites**: the runner's own bite test proves the check can fail. Reproduce it, and try your own: a verb named after the automation must turn it red |
| 4 | The skill says a recipe's state moves **through `maitre-cli`** — *recipe advance*, *ac verify*, *finding resolve* — **never by editing the plan's files** — and says it where the protocol says *"close the slice"*, and again in *What you must know* and *What you must not do* |
| 5 | The docs: `docs/03` has a section naming the three surfaces and the two rules; `AGENTS.md`'s layer table names `internal/cli` the artifact's surface and `internal/ui` the automation's |
| 6 | `internal/core` did not change. `core.DomainVersion` is **7**. The goldens did not move. The budget is 0 / 0 / 0 / 28 |
| 7 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The three words left out, and why it matters

The runner left **three words deliberately out** of the check — `start`, `run`, `kill` —
and put the reasoning in the receipt. `step start` is a legitimate verb on a plan's
record, so the word alone cannot be forbidden.

**This is the judgement to make.** The wire has two arms — every verb **classified**,
and no verb **named after** the automation — and the second arm's vocabulary is a
**hand-written list**. That is the shape this repo has been burned by twice: the
slice-number check's hand-written list (F-139) and the fixture's hand-written list of
field names (F-146). So:

| | |
|:--:|---|
| **does the classification arm catch a new verb?** | a verb that is **not** classified must fail the check — so a new `automation start` verb is caught by the *classification*, even if `automation` is not in the word list |
| **is the word list the wrong shape?** | or is a forbidden-words list the right tool here, given that the words are the **runtime's** vocabulary and the verbs are the **plan's** — two different vocabularies that meet in one dispatch table? |
| **what happens to the three excluded words** | they are the seam: `step start` is a plan record, and `start` is also the automation's. Say whether the seam is where it should be |

## Mutate, run, revert

| | The mutation |
|:--:|---|
| 1 | add a verb that starts a crew — the wire must go red |
| 2 | add a verb that is **not classified** — the wire must go red |
| 3 | take the CLI's name out of the skill's *The surfaces* section — a wire must catch it |
| 4 | make the skill say a recipe's state moves by editing the plan's files |

Revert every mutation. `git status --short` at the end is empty except your ledger
commit.

## The questions

Read [work/questions.md](../questions.md). Every row in it is struck — Q11–Q13 and
Q25 closed, Q16–Q30 ruled. Do not reopen them and do not answer them. Say whether
anything was newly raised.

## The scope

| | |
|:--:|---|
| **Implementation code** | do not edit it, except to mutate and revert |
| **`docs/audit.md`** | yours. One commit for the ledger, and the slice row for this cycle |
| **`work/`** | the orchestrator's, except your receipt |
| **The runner** | do not prompt it. Do not spawn agents |
| **`docs/05`, `docs/04`, `docs/06`** | not yours. Report a doc that contradicts the code; do not fix it |
| **The CLI's verbs** | do not add one and do not remove one. The wire **holds** the set; it does not change it |

## The verdict

| | |
|:--:|---|
| No HIGH or MEDIUM findings | `verdict: ok` |
| Any HIGH or MEDIUM finding | `verdict: changes-needed` |

Do not invent reasons to fail.

## Notify the orchestrator

The receipt comes first.

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-tui-t13.md <<'EOF'
Reviewer: T13 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-147.. or none>). The skill: <what the surfaces section says, and the count of the CLI's mentions>. The wire: <that it reads the dispatch table, what it classifies, what it forbids, and that it bites>. The recipe's state: <what the skill says moves it, and through what>. The three excluded words: <why they are out, and whether the seam is where it should be>. The docs: <what docs/03 and AGENTS.md say>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t13.md && git commit -q -m "inbox: reviewer reports tui-t13"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
