You are the **reviewer** for **T6**. You are a fresh lane. You did not write this code. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then [tui-t6.md](tui-t6.md) — the contract. The orchestrator's pane id is `w1H:pJ`. Use it verbatim.

The runner's commits: `5713472`, `92413be`, `cdfdc16`, and the state row `b2fe3a1`. The receipt is [work/inbox/runner-tui-t6.md](../inbox/runner-tui-t6.md).

## This slice draws the tree

File **HIGH or MEDIUM** blockers. New findings start at **F-126**. A LOW note goes in the message, not in the ledger as a finding.

| What | The target |
|---|---|
| **Below 160** | today's vertical tree. The connectors are `├─` and `└─`. `console-board-narrow.golden` pins it, and the other goldens did not move for this |
| **At 160 and above** | left to right. The trunk's column is left of the fork's column, on the wide capture the trunk is column 2 and the fork is column 81. The node shows the title. The dependency is still readable |
| **Nothing dropped** | an orphan is drawn. A cycle is drawn, and it says so. `core.TreeOf` is still the tree |
| **The footer** | the last row, at both widths. A node that does not fit is truncated, not wrapped |
| **The click** | wide and narrow. The hit map comes from the render. A click selects that recipe and opens its detail. `ui.Drive` expresses it. `enter` still opens the detail |
| **`space`** | offers no claim |

**`DomainVersion` stays 7.** No new view. The Board is still the Board. Ready is untouched.

## The orchestrator's ruling — the wide board takes the width

At 160 the other tabs put the detail beside the list. This slice does not.

| | |
|---|---|
| **The Board at ≥ 160** | uses the full width, so a fork has a column. The detail opens on `enter` and on a click. It does not sit beside the tree |
| **The other tabs** | unchanged. Their detail still sits beside the list at ≥ 160 |
| **Not a finding** | the full width is how the fork gets its column |

The pyte capture smudging the footer's middle dot is **not a finding**. Read the view's footer. The runner says it is `· r reload`.

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | Below 160 the connectors are `├─` and `└─`, and the narrow golden matches |
| 2 | At 160 the trunk's column is left of the fork's. The node shows the title. Say the columns |
| 3 | An orphan is on the board. A cycle is on the board and the row says so |
| 4 | The footer is the last row at 80 and at 160 |
| 5 | A driven click, wide and narrow, selects the recipe the box was drawn for and opens its detail. The box covers the title. Taking the box out of the render turns the click test red |
| 6 | `space` writes nothing on the Board. The other tabs' beside-detail at ≥ 160 still holds |
| 7 | `make check` and `make check-race` are green. `make goldens` leaves `git diff` empty |

**Mutate, run, revert.** One mutation that draws the wide tree below 160, and one that gives the trunk and the fork the same column. Report the tests that went red. Revert both. `git status --short` at the end is empty except your ledger commit.

## The questions

Read [work/questions.md](../questions.md).

| # | What to say |
|--:|---|
| Q8, Q9, Q10, Q3 | answered, not this slice |
| Q11, Q12, Q13 | still open. Do not answer them |
| The full-width board | not a new question. The ruling above closed it |

## The scope

| | |
|---|---|
| **Implementation code** | do not edit it, except to mutate and revert |
| **`docs/audit.md`** | yours. One commit for the ledger. Write the slice row |
| **`work/`** | the orchestrator's, except your receipt |
| **The runner** | do not prompt it. Do not spawn agents |

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
cat > work/inbox/reviewer-tui-t6.md <<EOF
Reviewer: T6 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-126.. or none>). The narrow tree: <connectors, and the golden>. The wide tree: <the columns, the title, the orphan, the cycle, the footer>. The click: <wide and narrow, the recipe, the detail, and that the map comes from the render>. The full width: <the Board takes it, the other tabs do not>. space: <writes nothing>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t6.md && git commit -q -m "inbox: reviewer reports tui-t6"
herdr agent prompt "w1H:pJ" "Reviewer: T6 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-126.. or none>). The narrow tree: <connectors, and the golden>. The wide tree: <the columns, the title, the orphan, the cycle, the footer>. The click: <wide and narrow, the recipe, the detail, and that the map comes from the render>. The full width: <the Board takes it, the other tabs do not>. space: <writes nothing>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>."
```

After the nudge, end your turn.
