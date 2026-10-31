You are the **reviewer** for **T4 fix cycle 1**. You are a fresh lane. You did not write this code. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then [fix-tui-t4-1.md](fix-tui-t4-1.md) — the contract — and [work/tui/README.md](../tui/README.md) for the state. The orchestrator's pane id is `w1H:pJ`. Use it verbatim in the notify command.

The runner's commits: `446fd27`, `a29935c`, `7488e02`, `a75d867`, `29f6396`. The receipt is [work/inbox/runner-tui-t4-fix1.md](../inbox/runner-tui-t4-fix1.md).

## This cycle is the four findings, and nothing else

File **HIGH or MEDIUM** blockers. Do not re-review T4 from scratch. Do not re-file what this cycle closed. A LOW note goes in the message, not in the ledger as a finding.

| Finding | The target | The commits |
|---|---|---|
| **F-122** | a row that does not fit is **one row**, and the **footer is the last row**. The plan name and the Settings `root` row. The runner also bounded the header's name — check that was required for the footer to stay | `446fd27` |
| **F-123** | Home's automation block **says the automation is not built** | `a29935c` |
| **F-124** | the recipes row and the evidence row name every state, the parts **add up**, and each row is **one line** on the seeded plan | `7488e02` |
| **F-125** | `docs/05` says what is true: the board shows the recipes, Home shows their counts. The counts on both is Q2, the full set | `a75d867` |

**F-122..F-125 were numbered by the orchestrator** from cycle 1's notes. They are not rows in `docs/audit.md` yet. **Write them as closed** against those commits, and write the slice row for this cycle. New findings start at **F-126**.

## The orchestrator's ruling — the two-digit row

The runner reported that a plan with a two-digit count in every section is **84 columns** and would wrap at 80. The seeded plan fits at **78**.

| | |
|---|---|
| **Do not truncate that row** | the tail is a state. A truncated state is a row that lies |
| **Do not drop `total`** | the row shows the parts and the number they add up to |
| **This cycle's bar** | the seeded plan: one line, the parts add up, the footer is last |
| **The all-99 wrap** | a **LOW note**, not a blocker. Shorter labels are the remedy if a later cycle must hold the footer at 80 for two-digit counts in every section. **Not this cycle** |

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | A plan whose **name is 100 characters**, at **80 columns**: the plan row is **one line**, and the **footer is the last row**. Say the name |
| 2 | The Settings **root** row, same width: one line, footer last |
| 3 | The automation block contains the words **the automation is not built** |
| 4 | The recipes' parts add up to the total, and the evidence's parts add up to the criteria count. Show the arithmetic. Each row is one line on the seeded plan |
| 5 | `docs/05` no longer says Home shows counts the board does not carry |
| 6 | `core.DomainVersion` is still **7**. No CLI verb. The other tabs' content did not move |
| 7 | `make check` and `make check-race` are green. `make goldens` leaves `git diff` empty. Say what moved |

**Mutate, run, revert.** One mutation per finding, the smallest edit that breaks that finding's test. Report the test that went red. Revert every mutation. `git status --short` at the end is empty except your ledger commit.

## The questions

Read [work/questions.md](../questions.md). Say whether a question was **answered, missed, or newly raised**.

| # | What this cycle should say |
|--:|---|
| Q8, Q9, Q10 | **answered** — ruled 2026-10-31. This cycle does not implement Q9 or Q10 |
| Q11, Q12, Q13 | still open. They are the orchestrator's. Do not answer them |
| The all-99 row | **not a new question** — the ruling above closed it for this cycle |

## The scope

| | |
|---|---|
| **Implementation code** | do not edit it, except to mutate and revert |
| **`docs/audit.md`** | yours. One commit for the ledger |
| **`work/`** | the orchestrator's, except your receipt |
| **The runner's pane** | do not prompt it. Do not spawn agents |

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
cat > work/inbox/reviewer-tui-t4-cycle2.md <<EOF
Reviewer: T4 cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-126.. or none>). F-122: <one line, footer last, the name, the width>. F-123: <the line on the screen>. F-124: <the arithmetic, and that each row is one line>. F-125: <what docs/05 says now>. Mutations: <what you mutated, what went red, and that you reverted it>. The all-99 row: <noted, not filed>. The questions: <answered, missed, or newly raised>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t4-cycle2.md && git commit -q -m "inbox: reviewer reports tui-t4 cycle 2"
herdr agent prompt "w1H:pJ" "Reviewer: T4 cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-126.. or none>). F-122: <one line, footer last, the name, the width>. F-123: <the line on the screen>. F-124: <the arithmetic, and that each row is one line>. F-125: <what docs/05 says now>. Mutations: <what you mutated, what went red, and that you reverted it>. The all-99 row: <noted, not filed>. The questions: <answered, missed, or newly raised>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>."
```

After the nudge, end your turn.
