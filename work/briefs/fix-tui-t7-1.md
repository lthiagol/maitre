You are a **fresh runner** on a **fix cycle**. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/tui/plan.md](../tui/plan.md) § 6, then [tui-t7.md](tui-t7.md) — the contract
the slice was built to — and then this brief. This brief wins where they differ.
The orchestrator's pane id is `w1H:p2`.

This is **T7, fix cycle 1**. The slice landed and the review filed two blockers and
one low. You are fixing three things and strengthening one test. You are **not**
re-implementing the filter.

The reviewer's receipt is [work/inbox/reviewer-tui-t7.md](../inbox/reviewer-tui-t7.md).
The findings are in `docs/audit.md` as **F-128**, **F-129** and **F-130**. Read
those rows — they carry the proof.

## F-128 — the filter changed a row's *content*

**The defect.** `shownSections` rebuilds the sections from the **filtered** recipe
set, and `core.SummaryOf` reads a recipe's readiness from the plan it is handed.
So a recipe that waits on a closed dependency loses its `●` (ready) mark the
moment the filter hides that dependency — while the Ready tab still lists it as
adoptable. The console contradicts itself between two tabs of the same screen.

**The proof.** On a plan where RCP04 depends on the closed RCP06:

| | the RCP04 row |
|---|---|
| unfiltered | `RCP04 approved  Add the backlog view  1/1 ●` |
| filtered to `backlog` | `1/1` — **no `●`** |
| filtered to `v1` (matches both) | the `●` is back |

**The target behaviour.** **A filter is a view of the plan; it must not change what
the plan says about a recipe.** The row's content must be **byte-for-byte the same**
filtered and unfiltered. The readiness mark is derived from the **whole** plan, so
the filter must filter the *rows*, not the plan the rows are computed against.

The fix belongs in `internal/ui` — the rows keep their content, and the filter
drops rows afterwards. If you conclude the domain must change, **report it**; do
not change `internal/core`.

## F-129 — at 80 columns the Board said nothing

**The defect.** The filter's *"the filter does not apply here"* message goes on the
summary row's right, and the Board's own summary is 47 columns wide. At 80 there is
no room for the 30-column message, so `summaryLine` truncates the left and **the
message is never drawn**.

**The proof.**

```bash
maitre --render --width 80  --keys "tab,tab,tab,tab,/"   # no message
maitre --render --width 100 --keys "tab,tab,tab,tab,/"   # the message
```

**The target behaviour.** A key that does nothing must say so, **at every width**.
At 80 columns the chef presses `/` on the Board and must be told the filter does
not apply here — the chef cannot tell the filter from a dead key.

The console already has a precedent for this: **the refusal banner takes rows of
its own** when the summary row cannot hold it (`bannerLines`, `writeBanner`). Use
the shape that fits — the message may take a row of its own, or the summary row may
be rearranged so the message fits — but the message is never dropped.

**The frame stays 24 rows and the footer stays the last row.** A message that
pushes the footer off the screen is the same defect as the one being fixed. The
list gives up a row for the message exactly as it gives up one for the filter
line, and no more.

Home, Log and Settings are unaffected — their summaries are short enough. The fix
must not disturb them.

## F-130 — the count's own rule was held by nothing

**The defect.** The brief and the code both say the count *"comes from the same
computation the list comes from"*. Nothing held it: the reviewer's mutation —
computing `N` by a **second** pass over the recipes instead of the list's rows —
turned **no** test red. The two computations agree on every seeded case, so the
rule is a claim, not a guard.

**The proof that they can drift.** With the served section hidden (`c`) and the
filter `glossary`: the list draws **0** rows and the count says `0 of 6`; the
second computation says `1 of 6`, because it counts a recipe the list does not
draw.

**The target behaviour.** `N` is the number of rows the list actually draws, taken
from the same computation that draws them — one fact, one answer. And a **test
holds it**: a test that goes red when `N` comes from a second computation. The
reviewer's divergence case is the one to pin.

## The test that saw nothing

The *stay-on-the-item* test filters with `RCP0`, which matches **all seven ids** —
so it cannot fail under any implementation. The behaviour is right and
`TestTheSelectionGoesToTheFirstMatch` holds the rule; this one witness is blind.

Make the filter text a **strict subset**, so the test can fail if the selection
does not stay on its item. Keep the test's name and its intent: the selection
stays on the same item when it still matches.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the row's content is identical filtered and unfiltered — the `●` survives — and the Recipes tab and the Ready tab agree about the same recipe. A test holds it |
| **AC-02** | at 80 columns, `/` on the Board draws *"the filter does not apply here"*. The frame stays 24 rows, the footer stays the last row, and Home, Log and Settings are unchanged |
| **AC-03** | the count's `N` is the number of rows drawn, from the list's own computation, and a test holds it — a test that goes red when `N` comes from a second computation |
| **AC-04** | the *stay-on-the-item* test filters a strict subset, so it can fail. The selection still stays on its item |
| **AC-05** | `work/tui/README.md`'s T7 row names this fix cycle's commits. `docs/05` still says what the filter is. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. No format change |
| **No new view** | nothing is registered in `core.DerivedViews`, so `docs/03` § 8 does not move |
| **`internal/core`, `internal/cli`** | do not change. The fix is the console's |
| **The tab count** | stays **7** |
| **The goldens** | the **four existing** goldens do **not** move. `console-filter.golden` **may** move — regenerate it with `make goldens` and commit it; it must stay claimed by a test |
| **The filter's contract** | the line, the keys, the match rule, the summary's shape and the tabs it does not reach all stay as they are. This cycle fixes what is broken; it does not redesign the filter |
| **Q11, Q12, Q13** | still open. Do not answer them |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **`internal/core`** | the domain does not change. If you think it must, report it |
| **`docs/audit.md`** | the reviewer's. The findings are closed by the cycle-2 reviewer, not by you |
| **`docs/03`, `docs/06`, `docs/07`** | not this cycle |
| **The footer** | do not name the filter key in the footer. The chef has the question open |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`work/tui/plan.md`** | the orchestrator's |
| **`work/`** | the orchestrator's, except the T7 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t7-fix1.md <<'EOF'
Runner: T7 fix 1 done. Commits: <n>. F-128: <the row's content under a filter, and that the Recipes tab and the Ready tab agree>. F-129: <the Board at 80, the message, the frame, the footer, and that the other tabs did not move>. F-130: <where N comes from, and the test that holds it>. The test: <the filter text now, and why it can fail>. The goldens: <the four, and console-filter.golden if it moved>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t7-fix1.md && git commit -q -m "inbox: runner reports tui-t7 fix 1"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses
a pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
