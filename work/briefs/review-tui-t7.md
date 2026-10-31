You are the **reviewer** for **T7**. You are a fresh lane. You did not write this
code. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[tui-t7.md](tui-t7.md) — the contract. The orchestrator's pane id is `w1H:p2`. Use it
verbatim.

The runner's commits: `7a76e66`, `0f264ac`, `0dce0c1`, `c542916`, `5ab1660`,
`1983b1a`, `af9deb1`, `65fdbba`, `74070ef`. The receipt is
[work/inbox/runner-tui-t7.md](../inbox/runner-tui-t7.md).

## The slice

**T7 — the filter.** `/` filters the list of the current tab, on **Recipes**,
**Backlog** and **Ready**. Board, Log, Home and Settings say the filter does not
apply there, and change nothing.

| | |
|---|---|
| **The line** | one row under the summary, while the chef types in it |
| **The keys** | `/` opens it · typing filters **live** · `enter` keeps the filter and closes the line · `esc` clears the filter and closes it |
| **The match** | one function, case-insensitive, over a recipe's id · title · intent (outcome and problem) · scope (in and out), and over a backlog item's id · title · body · scope · tags |
| **The count** | the summary says `filter: <text> · N of M` |
| **The selection** | stays on the same item (by **id**) when it still matches, goes to the first match when it does not, and the list scrolls to it |
| **The tabs** | the filter survives a tab switch and applies to the new tab |

**`DomainVersion` is 7.** Nothing outside `internal/ui` changed.

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | On Recipes, `/` opens the line. The list gives up **exactly one row**, and the footer is the last row of a 24-line screen at 80 columns |
| 2 | Typing filters **live** — no second keypress applies it. `esc` clears it and the list is whole again |
| 3 | The match is **case-insensitive**, and it reaches a recipe's **intent** and **scope**, and a backlog item's **tags** — not only the title. Test a match that only the intent could satisfy |
| 4 | The summary row reads `filter: <text> · N of M`, and **N is the number of rows shown** |
| 5 | The section headers count the items **shown**. A section with no match is **not drawn** |
| 6 | The selection **follows by id**: with the selection on RCP03, filter to a text that keeps RCP03 — it stays on RCP03, not on whatever is now at the same index. Filter away from it — it goes to the first match, and the list scrolls to it |
| 7 | The filter **survives a tab switch** and applies to the new tab's list |
| 8 | On Board, Log, Home and Settings, `/` says the filter does not apply here, and **nothing else changes** — the list, the selection and the journal are the same |
| 9 | `text:<word>` drives the filter in one step. **Outside the filter line it does nothing** — a script that meant to name a filter cannot quit the console |
| 10 | The four existing goldens did **not** move; the new `console-filter.golden` is **claimed by a test** |
| 11 | `internal/core` and `internal/cli` did not change. `core.DomainVersion` is **7**. `core.DerivedViews` did not change, so `docs/03` § 8 does not need to |
| 12 | `make check` and `make check-race` are green. `make goldens` leaves `git diff` empty |

**The receipt's claim about the race gate.** The receipt was written while
`make check-race` was still running. Confirm the gate yourself — do not take the
receipt's word for it.

**Nine commits for seven criteria.** The runner explains the two extra ones: one
carries AC-01's missing Ready-tab test and the new golden, one closes two holes in
the line's own invariant. Say whether that is coherent or whether a commit is
carrying two things.

**The footer does not name the filter key.** The runner's reason: the goldens pin the
footer and the brief says they do not move, so the key is named in the Settings tab's
key list instead. Say whether that is right, or whether the console is now hiding a
key the chef has to know to find.

## Mutate, run, revert

Break the slice's own invariants and report the tests that went red:

| | The mutation |
|--:|---|
| 1 | make `listAreaHeight` **not** drop by one while the line is open — the footer must fall off the screen |
| 2 | make the match **case-sensitive** |
| 3 | make the summary's **N** come from a second computation, not the list's |
| 4 | make the selection follow by **index** instead of by id |

Revert every mutation. `git status --short` at the end is empty except your ledger
commit.

## The questions

Read [work/questions.md](../questions.md). Q1–Q10 and Q14 are answered and not this
cycle. **Q11, Q12 and Q13 stay open — do not answer them.** Say whether anything was
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
cat > work/inbox/reviewer-tui-t7.md <<'EOF'
Reviewer: T7 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-128.. or none>). The line: <the row it takes, the row the list lost, and the footer>. The keys: <live, enter keeps, esc clears>. The match: <the rule, the fields, and that a recipe and an item share it>. The summary: <the filter row, N of M, the section counts>. The selection: <it follows by id, it moves, it scrolls, and the tab switch>. The other tabs: <what they say, and that nothing else changed>. The scripts: <text:<word>, and that it does nothing outside the line>. Mutations: <what you mutated, what went red, and that you reverted it>. The commits: <whether the nine are coherent>. The footer: <whether it should name the filter key>. The questions: <answered, missed, or newly raised>. The goldens: <unchanged / what moved>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t7.md && git commit -q -m "inbox: reviewer reports tui-t7"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
lane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
