Runner: T2b done. Commits: 6 (4dfd5a4, 187d706, f64e6a2, ac8df77, 588cf5f, 9c716ce) — the state row is 909ec28. Per-AC evidence: AC-01..AC-07 stamped. Pre-notify gate: pass (`make check`, `make check-race`; the tree clean).

**A click on a tab** switches to that tab, and it is the keyboard's `tab` made through one function: `Model.switchTab` — the same `m.tab = …`, the same `clampSelections()`, the same status reset. The click lands on the **label**; the separator between two labels belongs to no tab and changes nothing (`TestAClickOnTheGapBetweenTabsDoesNothing`). `TestTheClickOnATabIsTheKeyboardsSwitch` drives the click and the key into two consoles and compares every cursor and the status.

**A click on a row** selects it. The index comes from the **map** (`m.hits.rowAt`), the index the tab's own selection uses — the one the render drew the row from. A header row holds no item and selects nothing.

| Layout | What the click does | Why |
|---|---|---|
| **split** (≥160 columns, the detail beside the list) | selects; the detail follows the selection | the detail is already on screen — nothing else to do |
| **narrow** (the list alone) | selects **and** opens the detail | the detail is not on screen at all; a click that only selected would show the chef nothing — so it also makes the keyboard's `enter` |
| **Backlog** | selects and stops | the Backlog has no detail yet — that is T3a |
| **Settings** | selects the theme, and does **not** apply it | `enter` applies the theme, and applying it is a write |

**A click on the detail pane** is the console's `enter` on that recipe: the detail takes the screen, `esc` brings it back. In the full-screen detail it does nothing — because `enter` does nothing there. The brief's premise was false and the fix is the reason: the wheel already follows the pointer, so the click had nothing left to mean *for the wheel*.

**The click's spelling** is `click@(x,y)`, and `click` with no position (over the list, the console's default target). It is the **same grammar** as the wheel: the same parser reads it (`parseMouseStep`), and the position is bracketed, so `ui.SplitSteps` does not cut the comma inside it — `SplitSteps("click@(12,5),tab")` is two steps (`TestTheClickSurvivesTheSplitter`). The command line works: `maitre --render --keys "click@(12,2)" --plain` switches to the Backlog. An unknown step is refused with a message a reader can act on (`click:sideways`, `click@(1`, `click@1,5`).

**The goldens: unchanged.** A click changes what the console *answers*, not what it draws; `make goldens-check` regenerated them and found no diff, and no golden file is in the slice's commits. The Settings tab's own **KEYS** list was left alone for the same reason — adding the mouse there would have moved the draw.

Notes for the next session:

- **The map was not extended.** The tab labels, the rows and the detail pane were all already there; the click needed no new region. T2a's `there` flag and the row's item index did the work as they stand.
- **The console answers the press, not the release.** `handleMouse` now takes the `tea.MouseMsg` and switches on the concrete type (`tea.MouseWheelMsg`, `tea.MouseClickMsg`); a release would otherwise make every move happen twice.
- **One answer for the tab and for the detail.** `Model.switchTab` and `Model.openDetail` are the keyboard's moves, and the click calls them — a click that set the tab and skipped the rest would be a second answer to one question. The wheel's rule was not touched.
- `docs/05-human-interface.md`: the mouse table gained the three click rows, and the key table said `tab` is "next recipe" — it is the next tab; that line was false before this slice and is fixed.
- No domain change, no CLI change, no dependency: `core.DomainVersion` is 7, the tab count is 7, `internal/ui` still holds 28 external modules.
- Tripwires proved by mutation: dropping the narrow layout's open fails `TestAClickOnARowInTheNarrowLayoutOpensTheDetail`; dropping the row's selection fails the two row tests; emptying `switchTab` fails `TestTheClickOnATabClampsTheCursors`.
