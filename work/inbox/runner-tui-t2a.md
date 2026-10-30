Runner: T2a done. Commits: 8 (a72ea98, 115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 2bdc563). Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass.

**The wheel step: three lines** (`wheelStep`, internal/ui/mouse.go). One line is too slow to be useful and a page is too fast to aim with; three is what a terminal's own wheel reports per notch, and it is small enough to stop on. The wheel over the list drives `m.move` — the same scroll the keyboard drives, so the selection stays visible by construction — and the wheel over the detail drives `m.scrollDetail`.

**The scripted mouse spelling: `<verb>[:<arg>][@<x>,<y>]`** — `wheel:up`, `wheel:down`, `wheel:down@12,5`. The coordinates are the terminal's own: zero-based, (0,0) is the upper-left corner. A step with no position goes over the list (the console's default target) and over the detail when the list is not on screen. `ui.Drive` draws the frame before each step, so a mouse event is answered against the frame on screen. `parseMouseStep` already reads the whole spelling; `MouseFor` returns false for a verb it does not know, so **"click@x,y" is free** — T2b only adds the verb.

**The goldens: unchanged.** The mouse changes what the console answers, not what it draws, and the view's `MouseMode` is not in the content. `make goldens-check` regenerates them and finds no diff.

Notes for T2b:

| The thing | Where it is |
|---|---|
| the map | `internal/ui/hit.go` — `hitMap{tabs, rows, list, detail, banner}`, on `Model.hits`, rebuilt by every render |
| the lookups | `hitMap.tabAt(x, y)` → the tab index; `hitMap.rowAt(x, y)` → (the item, whether there is one) |
| what a row holds | the index the tab's **own selection** uses — the Recipes tab's visible-row index, the backlog index, the path index, the board's recipe index, the log event index, the palette index |
| the detail's region | the split layout: `listColumn()+1 .. width-1`; full-screen detail: the whole body under the header |
| the geometry | `tabBarRow` = `headerRows - 2`; the list starts at `headerRows + bannerRows` |
| the mouse's answer | `internal/ui/mouse.go` — `handleMouse`; a click is **not** implemented, on purpose |

**The prompt swallows the mouse.** `handleMouse` returns early while a prompt is open: every input belongs to the prompt, so the selection cannot move out from under a half-typed answer. The click must follow the same rule.

**One thing I did not do:** no click handler exists — a click that does nothing yet would look like a thing that works. The error and opening views clear the map, so no event can land on a tab that is not on screen.
