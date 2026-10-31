Reviewer: T7 done. Verdict: changes-needed. Filed 2 blockers (F-128, F-129) and 1 low (F-130).

**The line.** One row under the summary, between the summary and the list. The list gives up exactly one row while the line is open (`listAreaHeight` 19 → 18 at 80×24), the frame stays 24 rows, and the footer is the last row. The split layout at 160×30 holds too. The line reads ` filter: <text>▏`.

**The keys.** Live: the list narrows as the chef types, no second keypress — `menu` filters to one recipe with no `enter`. `enter` keeps the filter and closes the line. `esc` clears it and closes it, and the list is whole again.

**The match.** One function, `matchFilter` — case-insensitive, a substring — over a recipe's id · title · intent (outcome and problem) · scope (in and out), and over a backlog item's id · title · body · scope · tags. A recipe and a backlog item share it, and one test holds it for both kinds.

**The summary.** `filter: <text> · N of M` — `menu · 1 of 7`, `board · 3 of 6`, `evidence · 1 of 2`. The section headers count the items shown (*Needs you 1*), and a section with no match is not drawn at all.

**The selection.** It follows by id: RCP03 stays on RCP03 under a filter that keeps it; filter `sharpen` takes RCP05 from row 3 to row 0 and it stays on RCP05. It goes to the first match when the item is filtered away, and the list scrolls to it. The filter survives a tab switch: Recipes → Backlog keeps the text `the` and the new tab's list is filtered (5 of 6).

**The other tabs.** Board, Log, Home and Settings say *the filter does not apply here*, and nothing else changes — the list, the selection and the journal are the same (`activity.json` byte-for-byte). **But at 80 columns the Board says nothing at all** — F-129. Its summary is 47 columns wide, so the 30-column message does not fit, `summaryLine` truncates the left instead, and the message is never drawn. `maitre --render --width 80 --keys "tab,tab,tab,tab,/"` shows *the tree · 6 recipes · 4 can go now · 2 waiting* and no message; at 100 the same command shows the message. Home, Log and Settings are unaffected. A key that does nothing must say so; at 80 the key does nothing **and says nothing**.

**F-128 — the filter changed a row's content.** `shownSections` rebuilds the sections from the *filtered* recipe set, and `core.SummaryOf` reads readiness from the plan it is given — so a recipe that waits on a closed dependency loses the `●` (ready) mark as soon as the filter hides that dependency. On a plan where RCP04 depends on the closed RCP06: unfiltered the row reads `RCP04 approved  Add the backlog view  1/1 ●`; filtered to `backlog` the same recipe reads `1/1` with **no `●`**, while the Ready tab still lists it as adoptable; filtered to `v1` (matching both) the `●` is back. The console contradicts itself between two tabs of the same screen. A filter is a view of the plan; it must not change what the plan says about a recipe.

**The scripts.** `text:<word>` drives the filter in one step, and the single-character steps still work one character each. Outside the filter line it does nothing: `text:quit` did not quit the console (the tab and the filter were untouched). A `text:` step is one word, not a run of key presses.

**Mutations** — all reverted, `git status` clean:

| | The mutation | What went red |
|--:|---|---|
| 1 | `listAreaHeight` not dropping by one while the line is open | **4 tests**: `TestTheFilterLineTakesOneRow`, `TestTheFilterTextCannotBreakTheScreen`, `TestALongFilterDoesNotBreakTheScreen`, `TestTheFilterMatchesTheGolden` |
| 2 | the match case-sensitive | **1 test**: `TestTheMatchRuleHoldsForBothKinds` |
| 3 | the summary's `N` from a second computation, not the list's | **none** — and it *does* diverge: with the served section hidden and the filter `glossary`, the list draws 0 rows and the count says `0 of 6`; the second computation says `1 of 6`. **F-130** |
| 4 | the selection by index, not by id | **1 test**: `TestTheSelectionGoesToTheFirstMatch` |

**The commits.** The nine are coherent. Seven carry one criterion each (AC-01 `7a76e66`, AC-02 `0f264ac`, AC-03 `0dce0c1`, AC-04 `c542916`, AC-05 `5ab1660`, AC-06 `1983b1a`, AC-07 `74070ef`), and the two extras are honest: `af9deb1` carries AC-01's missing Ready-tab test *and* the new golden — the same criterion's evidence, and the golden is the evidence for it — and `65fdbba` closes two holes in one invariant (the line's own: the footer is the last row). Nothing carries two criteria.

**The footer.** The runner's reason is **true**, not a dodge: the reviewer's own mutation — `/ filter` added to the Recipes footer — turns `TestConsoleMatchesTheGolden` and `TestTheStressPlanMatchesTheGolden` red, and the brief says the four existing goldens do not move. The key is named in the Settings tab's key list and in `docs/05`'s interaction table. Whether a key that opens a *mode* belongs in the footer (which names what the **next** key does) or in the Settings list is a question about the console's language, not a defect. Recorded as a note, not a finding — but the chef has to look in Settings to find it.

**The questions.** Q11, Q12 and Q13 stay open and unanswered, as the brief said. **Nothing newly raised.** Two low notes recorded in the ledger instead: the *stay-on-the-item* test filters with `RCP0`, which matches all seven ids, so it cannot fail under any implementation — the behaviour is right and `TestTheSelectionGoesToTheFirstMatch` holds the rule; and the count's provenance is held by nothing.

**The goldens.** The four existing goldens did **not** move (the filter is off by default). `console-filter.golden` is new, and `make goldens-check` says *every golden file is claimed by a test*. `make goldens` leaves `git diff` empty.

**DomainVersion: 7.** `internal/core` and `internal/cli` did not change at all; `core.DerivedViews` did not change, so `docs/03` § 8 does not need to.

**The gate.** `make check` green. `make check-race` green — confirmed by this review, not taken from the receipt.

**LOW notes.** F-130 (the count's rule held by nothing); the *stay-on-the-item* test is a witness that saw nothing; the filter key is in the Settings list and not in the footer.
