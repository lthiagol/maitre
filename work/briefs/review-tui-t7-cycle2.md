You are the **reviewer** for **T7 fix cycle 1**. You are a fresh lane. You did not
write the fix. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[fix-tui-t7-1.md](fix-tui-t7-1.md) — the contract for this cycle — and then
[tui-t7.md](tui-t7.md), the contract the slice was built to. The orchestrator's pane
id is `w1H:p2`. Use it verbatim.

The fix commits: `675b4cc`, `def86f3`, `00c234e`, `8c505d7`, and the docs commit
`ad453df`. The receipt is [work/inbox/runner-tui-t7-fix1.md](../inbox/runner-tui-t7-fix1.md).

## This cycle is three findings

File **HIGH or MEDIUM** blockers. Do not re-review T7 from scratch. New findings
start at **F-131**. A LOW note goes in the message, not in the ledger as a finding.

| Finding | The target | The commit |
|---|---|---|
| **F-128** | the filter changed a row's *content*: `shownSections` rebuilt the sections from the **filtered** set, so `core.SummaryOf` read readiness from a plan the filter had already narrowed, and a recipe waiting on a closed dependency lost its `●` when the filter hid that dependency — while the Ready tab still called it adoptable | `675b4cc` |
| **F-129** | at 80 columns the Board said **nothing** when the chef pressed `/`: the message goes on the summary row's right, the Board's summary is 47 columns, and `summaryLine` truncated the left instead | `def86f3` |
| **F-130** | the count's own rule was held by nothing: `N` computed by a second pass over the recipes turned no test red, and the two computations diverge (served hidden + filter `glossary` → the list draws 0 rows and says `0 of 6`; the second pass says `1 of 6`) | `00c234e` |

Close all three in `docs/audit.md` against those commits, and write the slice row
for this cycle.

**The filter's contract stays.** The line under the summary, the live match, the
one rule (`matchFilter`), `filter: <text> · N of M`, the selection following by id,
the filter surviving a tab switch, `text:<word>`, and the tabs the filter does not
reach — all as they were. Confirm they still hold; do not file them again.

**The footer's question is answered.** The chef ruled (Q15): the filter key lives in
the **Settings tab's key list**, not the footer. Do not file that note again.

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | On a plan where RCP04 depends on the closed RCP06: the RCP04 row under a filter is **byte-for-byte** the row without one — the `●` survives — and the Ready tab still lists it as adoptable. The two tabs agree |
| 2 | `maitre --render --width 80 --keys "tab,tab,tab,tab,/"` — the Board says *the filter does not apply here*. The frame is **24 rows**, the footer is the last row, and the list gives up exactly one row |
| 3 | Home, Log and Settings are **unchanged** by the fix — their summaries still hold the message beside them, and nothing moved |
| 4 | The count: with the served section hidden (`c`) and the filter `glossary`, the list draws **0** rows and the count says `0 of 6`. A count computed by a second pass over the recipes says `1 of 6` — the divergence the finding names |
| 5 | The *stay-on-the-item* test filters a **strict subset** — it filters `review` — so a cursor left on its old index names RCP07 and the test goes red. The selection still stays on its item |
| 6 | **All five goldens** did not move, including `console-filter.golden`. `make goldens` leaves `git diff` empty, and every golden is still claimed by a test |
| 7 | `internal/core` and `internal/cli` did not change. `core.DomainVersion` is **7**. `core.DerivedViews` did not change, so `docs/03` § 8 does not need to |
| 8 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

**The shared message shape.** The runner says a message that does not fit beside
the summary now takes a row of its own, *whether it is a refusal or a status*. That
is a wider change than the finding asked for: check that the **refusal banner**
still behaves — a refusal that fits beside the summary still goes beside it, and a
refusal that does not still takes rows of its own, and the reason is never
truncated away.

## Mutate, run, revert

Put each finding back, and report the tests that went red:

| | The mutation |
|--:|---|
| 1 | feed the **filtered** plan into the rows' readiness again — the `●` must vanish and the test must go red |
| 2 | drop the message's own row — the 80-column Board must go silent and the test must go red |
| 3 | compute `N` by a second pass over the recipes — the count must diverge and the test must go red |
| 4 | leave the selection on its old index — the *stay-on-the-item* test must go red |

Revert every mutation. `git status --short` at the end is empty except your ledger
commit.

## The questions

Read [work/questions.md](../questions.md). Q1–Q10, Q14 and Q15 are answered and not
this cycle. **Q11, Q12 and Q13 stay open — do not answer them.** Say whether
anything was newly raised.

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
cat > work/inbox/reviewer-tui-t7-cycle2.md <<'EOF'
Reviewer: T7 cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-131.. or none>). F-128: <the row's content under a filter, and that the two tabs agree>. F-129: <the Board at 80, the message, the frame, the footer, and the other tabs unmoved>. F-130: <where N comes from, and the divergence case>. The refusal banner: <whether the shared message shape kept it>. The test: <the filter text, and that it can fail>. Mutations: <what you mutated, what went red, and that you reverted it>. The goldens: <all five, unmoved>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t7-cycle2.md && git commit -q -m "inbox: reviewer reports tui-t7 cycle 2"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
