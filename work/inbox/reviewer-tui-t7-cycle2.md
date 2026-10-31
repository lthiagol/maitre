Reviewer: T7 cycle 2 done. Verdict: ok. Filed 0 blockers (none). F-128, F-129 and F-130 are closed.

**F-128.** The row's content under a filter, and the two tabs agree. On a plan where RCP04 depends on the closed RCP06, the row reads `RCP04 approved  Add the backlog view  1/1 ●` **unfiltered and filtered to `backlog` — byte-for-byte the same** — and the Ready tab still lists it as adoptable, so the two tabs agree. `shownSections` now keeps the sections the **whole** plan produced and filters the *items* inside them, so the `core.Summary` a row carries is the one the whole plan computed and the readiness mark cannot move under a filter. The reviewer's probe went further than the fix's own test: **every** recipe's row is byte-for-byte the same filtered and unfiltered, and the section order is preserved under every filter (*Needs you · On the line · Being shaped · Served*).

**F-129.** The Board at 80, the message, the frame, the footer, and the other tabs unmoved. `maitre --render --width 80 --keys "tab,tab,tab,tab,/"` now shows *the filter does not apply here* **on a row of its own** (row 4); the frame is **24 rows**, the footer is row **23**, and the list gives up exactly **one** row (19 → 18). **Home, Log and Settings are untouched** — their summaries are short enough, so the message lands on the summary row (row 3) and the list keeps every row it had. The shape is the one the brief pointed at: a message that does not fit takes a row of its own, the way a refusal already did.

**The refusal banner held.** The runner widened the shape to statuses as well as refusals — wider than the finding asked for, and the brief said to check it. It held: a **short** refusal still sits on the summary row, right-aligned, in the **warning** style (`⛔ refused — review.clean` beside the summary text); a **long** refusal still takes rows of its own with the `⛔` marker and the **whole reason** on the screen; and with the **filter line open** as well, the frame is still 24 rows with the footer last and the banner between them.

**F-130.** Where `N` comes from, and the divergence case. `N` is now the number of rows the list draws, counted from the list's **own lines** — the same function that draws the screen — so the count and the screen cannot disagree by construction. The divergence case: with the served section hidden (`c`) and the filter `glossary`, the list draws **0** rows and the count says **0 of 6**; the second computation would say **1 of 6**. The count equals the rows the list draws on every tab and every filter — *Recipes* 7 of 7 and 1 of 7, *Backlog* 5 of 6 and 3 of 6, *Ready* 2 of 2 and 1 of 2.

**The test.** The *stay-on-the-item* test filters `review` now — a **strict subset** that matches RCP03 and RCP07 and takes the row above the cursor (RCP02) with it. The filtered list is two rows long, so a cursor left on its old index names **RCP07**, and the one that follows its item names **RCP03**. It can now fail.

**Mutations** — all reverted, `git status` clean:

| | The mutation | What went red |
|--:|---|---|
| 1 | the filtered plan fed into the rows' readiness again | **1 test**: `TestTheFilterDoesNotChangeTheRow` |
| 2 | only a refusal takes rows of its own (the message's row dropped) | **1 test**: `TestTheBoardSaysTheFilterDoesNotApplyAtEighty` |
| 3 | `N` by a second pass over the recipes | **1 test**: `TestTheCountIsTheRowsTheListDraws` |
| 4 | the selection left on its old index | **2 tests**: `TestTheSelectionStaysOnTheSameItem`, `TestTheSelectionGoesToTheFirstMatch` |

**The filter's contract stands** (confirmed, not re-filed): the line is one row under the summary (19 → 18 at 80×24, the frame 24 rows, the footer last); the match is live; `enter` keeps and `esc` clears; `matchFilter` is one rule for a recipe and an item; the count is `filter: <text> · N of M`; the selection follows by id (RCP03 → RCP03); the filter survives a tab switch (Recipes → Backlog, 5 of 6); `text:quit` outside the line does not quit the console.

**The goldens.** All **five** unmoved, `console-filter.golden` included; `make goldens` leaves `git diff` empty and every golden file is still claimed by a test.

**DomainVersion: 7.** `internal/core` and `internal/cli` did not change at all; `core.DerivedViews` did not change, so `docs/03` § 8 does not need to.

**The gate.** `make check` green and `make check-race` green — confirmed by this review, not taken from the receipt.

**LOW notes.** The count is now computed from the list's own lines, so a filtered frame costs about a third more than it did: on a **1000-recipe** plan **18.1 ms** against **13.0 ms** before the fix (**7.0 ms** with no filter on); the seeded plan is unchanged (285 µs). It is the same shape as the kinds-list note already in the ledger — a per-frame allocation on the path that redraws the screen — and the count's provenance is the rule F-130 closed, so it is recorded as a note and not a finding. Worth knowing before a plan reaches the thousands. Nothing else raised; Q11, Q12 and Q13 stay open.

Ledger: `docs/audit.md` (`61a910b`).
