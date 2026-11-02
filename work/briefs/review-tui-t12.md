You are the **reviewer** for **T12**. You are a fresh lane. You did not write this
code. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[tui-t12.md](tui-t12.md) — the contract. The orchestrator's pane id is `w1H:p2`. Use it
verbatim.

The runner's commits: `d197c41` (AC-01, the render path forks no process), `b0040e0`
(AC-02, the crew read costs one process), `ade00a2` (AC-03, the console still asks),
`3bfdb85` (the read says what it holds), and the row `5b77cc7`. The receipt is
[work/inbox/runner-tui-t12.md](../inbox/runner-tui-t12.md).

## This slice

**T12 — the render path stops waiting.** It changes **no behaviour the chef can
see**. The Sous-Chef tab still shows the crew, and the live console still asks the
runtime. It changes **when** the console asks, and **how long** a frame takes.

The chain the slice was written against:

```
render.go:41   snapshot()          →  waited
crew.go:54     readCrew()           →  waited
herdr.go:76    exec.Command(herdr)  →  two spawns
```

## What to verify — each by command, not by reading

| # | Check |
|:--:|---|
| 1 | **No process spawn in a snapshot.** Put a **fake `herdr`** on the PATH that counts its invocations, render a frame, and the count is **zero**. The frame forks nothing |
| 2 | The crew read is **one process**: `herdr api snapshot` — **one** document, **both** answers (the crew's members and each one's job). A fake herdr on the PATH is called **once** |
| 3 | The live console **still asks**: the poll refreshes the crew while the tab is visible, and the tab shows **what it last heard**. F-140's witness still holds **at the process level** — the mouse's switch and the filter line's switch reach the runtime, and a click never leaves the tab saying *the automation is not running* while it is |
| 4 | The **goldens** are unchanged — `git diff` over `internal/ui/testdata/` is empty — **and** they render with **herdr off the PATH** |
| 5 | The **poll's frequency** is unchanged. It is deliberate: F-140's witness requires the poll to ask while the tab is visible. What changed is **what** the poll asks, not **how often** |
| 6 | **What the tab shows** is unchanged: the crew's rows, their states, their jobs, the blocked agent's question |
| 7 | `internal/core` did not change. `core.DomainVersion` is **7**. The budget is unchanged: 0 / 0 / 0 / 28 |
| 8 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The judgement call the runner made

**The tab has a third state.** A console that **has not asked** says *"the crew has not
been asked"* and offers **no move** — rather than guess *"the automation is not
running"*.

That is the same rule the round has been teaching: **the console never says a fact it
does not know.** A console that has not asked the runtime does not know whether the
crew is running; saying "the automation is not running" would be a guess, and the
footer offering a move on a guess is the move that starts a second crew (F-140).

Judge it: is the third state right, is the wording right, and does the tab still
answer the chef's question — *can I start the crew?* — when the console has not asked
yet?

## Mutate, run, revert

| | The mutation |
|:--:|---|
| 1 | make the snapshot spawn a process again |
| 2 | make the crew read cost **two** processes again |
| 3 | make the poll stop asking while the tab is visible |
| 4 | make the tab say *the automation is not running* when it has simply not asked |

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
| **`docs/audit.md`** | yours. One commit for the ledger, and the slice row for this cycle — and, if T12's findings are all closed, **the row that closes the slice** and the round |
| **`work/`** | the orchestrator's, except your receipt |
| **The runner** | do not prompt it. Do not spawn agents |
| **`docs/03`, `docs/04`, `docs/06`** | not yours. Report a doc that contradicts the code; do not fix it |

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
cat > work/inbox/reviewer-tui-t12.md <<'EOF'
Reviewer: T12 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-147.. or none>). The snapshot: <that it forks nothing, and the fake-herdr witness that counts zero calls>. The read: <that it is one process, one document, both answers>. The console: <that it still asks, what the tab shows, and that F-140's witness holds at the process level>. The goldens: <unchanged, and that they render with herdr off the PATH>. The poll: <that its frequency is unchanged>. The third state: <the wording, and whether it answers the chef's question>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t12.md && git commit -q -m "inbox: reviewer reports tui-t12"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
