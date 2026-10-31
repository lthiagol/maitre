You are the **reviewer** for **T6 fix cycle 1**. You are a fresh lane. You did not write this code. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then [fix-tui-t6-1.md](fix-tui-t6-1.md) — the contract. The orchestrator's pane id is `w1H:pJ`. Use it verbatim.

The runner's commits: `a5df3ce`, `df6d387`, and the state row `45af5ab`. The receipt is [work/inbox/runner-tui-t6-fix1.md](../inbox/runner-tui-t6-fix1.md).

## This cycle is two findings

File **HIGH or MEDIUM** blockers. Do not re-review T6 from scratch. New findings start at **F-128**. A LOW note goes in the message, not in the ledger as a finding.

| Finding | The target | The commit |
|---|---|---|
| **F-126** | `space` on the Board writes nothing, including on RCP02, which is ready to sign off. The footer names no `space`. The Recipes tab still signs that recipe off | `a5df3ce` |
| **F-127** | at 80 columns every Board row is one line, and the footer is the last row. The seeded branch rows were 82. The cycle rows were 114. A truncated cycle note still says **dependency cycle**, and the titles stay | `df6d387` |

Close both in `docs/audit.md` against those commits, and write the slice row for this cycle.

The two drawings stay. Below 160 the connectors are `├─` and `└─`. At 160 the trunk is left of the fork. A click still opens the detail. The wide drawing and the hit map were not part of this fix — confirm they still hold, and do not file them again.

The 100-column golden was reported unchanged, because the longest seeded row there is 85. If `make goldens` moves it, say why.

**`DomainVersion` stays 7.**

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | On the Board, `down` to RCP02, `space`. `activity.json` is byte for byte the same, and the recipe is still `ready` |
| 2 | On the Recipes tab, `space` on that same recipe signs it off, `ready` → `approved` |
| 3 | Seeded plan at 80: every Board row is one line, the widest is 80, the footer is the last row of a 40-line screen |
| 4 | The cycle fixture at 80: the widest row is 80, the titles `Cycle left` and `Cycle right` are whole, the note still says `dependency cycle`, the footer is last |
| 5 | `make check` and `make check-race` are green. `make goldens` leaves `git diff` empty |

**Mutate, run, revert.** Put `actOnRecipe` back on the Board, and stop truncating the narrow row. Report the tests that went red. Revert both. `git status --short` at the end is empty except your ledger commit.

## The questions

Read [work/questions.md](../questions.md). Q8, Q9, Q10 and Q3 are answered and not this cycle. Q11, Q12 and Q13 stay open. Do not answer them. Say whether anything was newly raised.

## The scope

| | |
|---|---|
| **Implementation code** | do not edit it, except to mutate and revert |
| **`docs/audit.md`** | yours. One commit for the ledger |
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
cat > work/inbox/reviewer-tui-t6-cycle2.md <<EOF
Reviewer: T6 cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-128.. or none>). F-126: <the journal unchanged, and Recipes still signs off>. F-127: <the widths, one line, the cycle note, the footer>. The drawings: <connectors, columns, the click>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. The goldens: <unchanged / what moved>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t6-cycle2.md && git commit -q -m "inbox: reviewer reports tui-t6 cycle 2"
herdr agent prompt "w1H:pJ" "Reviewer: T6 cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-128.. or none>). F-126: <the journal unchanged, and Recipes still signs off>. F-127: <the widths, one line, the cycle note, the footer>. The drawings: <connectors, columns, the click>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. The goldens: <unchanged / what moved>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>."
```

After the nudge, end your turn.
