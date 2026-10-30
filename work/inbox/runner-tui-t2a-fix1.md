Runner: T2a fix 1 done. Commits: 4 (619b342, 47be828, f1cc1ca, plus the state row). Fixed: F-106, F-107, F-108.

**The detail's scroll.** The offset is applied **once**, where the detail is drawn: `Model.detailScrollRows` (internal/ui/detail.go) is called by `renderDetailFull`, by `renderSplit` and by `viewDetail`. `renderSplit` looked at the offset never — that was the hole: the wheel moved `m.detailScroll` from 0 to 3 and the frame came out byte-for-byte identical, in a pane the map advertises.

How I proved the split layout moves — two instruments:

* in the package: `go test ./internal/ui -run 'TestTheWheelScrollsTheSplitDetail' -v` → PASS. The test draws the split frame, wheels over the detail, and asserts the offset moved **and** that the frame differs.
* from the command line: `--render` alone never draws the split layout, but `--keys` drives the real `View`, so
  `diff <(maitre --render --width 200 --keys "wheel:up@(150,6)" --plain) <(maitre --render --width 200 --keys "wheel:down@(150,6)" --plain)`
  shows the detail column moving by three rows.

**The row count's one answer.** The bound is counted at the width the pane **draws with**: `Model.detailPaneWidth()` (internal/ui/detail.go) — the full screen's `width-4` in the full-screen detail, and `width - listColumn() - 3` in the split layout, which is the width `renderSplit` draws with (it now asks the same function). `detailLineCount()` uses it, so the bound and the pane cannot disagree: `go test ./internal/ui -run 'TestTheScrollsBoundIsThePanesWidth' -v` → PASS (the bound is the pane's own row count, and the scroll reaches the pane's last row).

**The position form: `wheel:down@(12,5)` — the position is bracketed.** The steps are separated by commas, and so were the coordinates, so the command line split the step in half before the driver saw it. The bracket says "this is one thing", and the comma inside it is not a separator. The alternative — a different coordinate separator — would have made the pair unreadable (`12;5`, `12:5`); the alternative — changing the step separator — would have broken `maitre --keys "tab,down,space"`, which the plan says stays. The splitter that knows the grammar (`ui.SplitSteps`) lives beside the parser that reads it (`parseMouseStep`), so the command line and the driver cannot disagree about what a step is.

The command that works:

```
maitre --render --keys "wheel:down@(12,5)" --plain   # exit 0
```

The old spelling now fails loudly instead of half-working: `maitre --render --keys "wheel:down@12,5" --plain` → `maitre: malformed mouse step "wheel:down@12"`.

**The doc.** `docs/05-human-interface.md` prints `maitre --keys "wheel:down@(12,5)"` — the form that works — and says why the position is bracketed. It also now says what the wheel does where the pointer is not over the detail: the list answers, because the list is the console's one scrollable thing.

**The absent region.** `hitBox` carries a `there` flag: the zero box is **not there**, and `contains`/`empty` say so. The regions that are drawn are built with `newHitBox`. The wheel falls back to the list where the pointer is not over the detail — so the wheel at the origin scrolls the **list** (the phantom detail cell is gone), and the cleared map of the error and opening views answers nothing at all.

Proved with a mutation, not a claim: taking the `there` guard out of `contains`/`empty` — the pre-fix behaviour — and running `go test ./internal/ui -run 'TestTheOriginIsNotAPhantomDetail|TestAClearedMapAnswersNothing'` fails with the finding's own symptoms (*the wheel at the origin scrolled the detail to 3*). With the fix, both PASS. The CLI agrees: `--keys "wheel:down@(0,0)"` moves the list, and the footer loses *space sign off* because the selection moved off RCP02.

DomainVersion: 7. The goldens: **unchanged** — `make goldens-check` regenerates them and finds no diff; the mouse and this cycle change what the console answers, not what it draws. The gate: `make check` green, `make check-race` green. `internal/ui` still 28 external modules, no dependency added, `internal/core` and `internal/cli` untouched.

**One behaviour change worth the reviewer's eye:** the wheel over the **header** now scrolls the list. Before, it did nothing, and the reviewer's own evidence for F-108 asks for a wheel at the origin to scroll the list. The full-screen detail has no list beside it, so there the wheel over the header still does nothing — the list answers only when the list is on screen.

Notes for T2b:

| The thing | Where it is |
|---|---|
| the step grammar | `wheel[:up\|down][@(x,y)]` — the position is **bracketed**; `ui.SplitSteps` splits a command line into steps, `parseMouseStep` reads one |
| the click's spelling is free | `mouseVerbs` holds the gestures the console answers; `MouseFor` refuses a verb it does not know, and `isMouseStep` tells a malformed gesture from an unknown key |
| the regions | `hitMap.tabs`, `.rows`, `.list`, `.detail`, `.banner` on `Model.hits`; a region that is not drawn is the **zero** `hitBox` and answers nothing |
| the lookups | `tabAt(x, y)` → the tab index; `rowAt(x, y)` → (the item, whether there is one); the item is the index the tab's own selection uses |
| the detail's width | `detailPaneWidth()` — the pane's own width, and the bound's width. The click does not change it |
| the wheel's fallback | the detail answers over the detail; the list answers everywhere else, when the list is on screen |
