Reviewer: T2a review done. Verdict: changes-needed. Filed 2 blockers (F-106, F-107).

**The map is built by the render: yes.** `tabBar` emits the tab boxes as it draws them; `hitMapOf` is called from the three render functions (`renderSingle`, `renderSplit`, `renderDetailFull`) and from nowhere else; the banner pushes the list's rows down in the map exactly as it does on the screen (`hitMapOf` takes the banner's row count the writer returned); `headerRows`, `footerRows` and `splitThreshold` are read where they are written (no recomputed numbers). The map is rebuilt on every render and kept on `Model.hits`.

**The wheel is the keyboard's scroll: yes, and it is the same variable.** The wheel over the list drives `m.move` — the keyboard's own move — so there is one selection and one offset. Compared by command: on a plan of 47 recipes, `maitre --render --keys "wheel:down"` and `maitre --render --keys "down,down,down"` are **byte-identical**; the unit test compares the selection and the offset and both land on the same place.

**The scripted mouse:** the spelling is `<verb>[:<arg>][@<x>,<y>]` — `wheel:up`, `wheel:down`, `wheel:down@12,5`. `ui.Drive` draws the frame **before** each step (`m.View()` at the top of the loop), so a mouse step is answered against the frame on screen; an unknown step is refused with `unknown step "..."`. **But the position form cannot be reached through `maitre --keys`** — the steps are split on `,` and so are the coordinates: `maitre --render --keys "wheel:down@12,5"` prints `maitre: unknown step "wheel:down@12"` and exits 1. That is **F-107**, and `docs/05-human-interface.md` prints exactly that command as the way to drive the wheel at a place.

**The end-to-end check:** `maitre --render --keys "enter,wheel:down"` moved the **full-screen** detail (the recipe's first three rows scrolled off, the frame diffed). The **list** scroll needed a plan that overflows — I made one (47 recipes): `wheel:down` and `down,down,down` are byte-identical, and 40 `down`s differ. What I could **not** prove this way: the wheel over the detail in the **split** layout — `--render` draws the tabs one after another and never the split layout, so that one is only observable in a probe: at 200x30 the wheel over the detail's own box takes `detailScroll` from 0 to 3 and the frame is **byte-identical** — `renderSplit` never reads the offset — while the same wheel in the full-screen detail does move it. That is **F-106**.

**Mutation checks (all six failed a test, all reverted):**
| # | Mutation | What failed |
|--:|---|---|
| 1 | the `MouseMode` off the error view | `TestTheMouseIsOn` — "the error view's mouse mode = 0, want cell motion" |
| 2 | the render stops emitting the map | `TestTheHitMapKnowsWhereTheTabsAre`, `…WhereTheRowsAre`, `TestTheBannerPushesTheListDown`, `TestTheWheelScrollsTheList`, `TestTheWheelIsTheKeyboardsScroll` — five tests |
| 3 | the wheel scrolls the offset instead of the keyboard's scroll | `TestTheWheelScrollsTheList`, `TestTheWheelIsTheKeyboardsScroll` |
| 4 | the wheel over the header scrolls the list | `TestTheWheelOverTheHeaderDoesNothing` |
| 5 | the wheel acts while a prompt is open | `TestTheWheelLeavesAnOpenPromptAlone` (both the list and the detail) |
| 6 | the scripted step acts without drawing first | `TestTheMouseIsDrivableFromAScript` |

Every mutation was reverted with `git checkout --`; the tree is exactly as the runner left it (`git diff e4b7061..HEAD` touches only `docs/audit.md`, the orchestrator's brief, its inbox and the state row — no implementation file).

**The goldens: unchanged.** `make goldens` leaves `git diff` empty.

**The gate:** `make check` green and `make check-race` green. The scope held: no click handler anywhere (the click is T2b), `internal/core` and `internal/cli` untouched, `core.DomainVersion` still 7, the budget still 0 / 0 / 28, seven tabs with the same names, `docs/audit.md` not swept by the slice, the state row updated and committed.

**LOW notes (not filed):**
- The "no region" box is the zero value, and the zero value is **not** empty: `hitBox{}.empty()` is false and it contains (0,0). So in the list layout the map's detail region is a phantom cell at (0,0) — a wheel on the very first cell scrolls an invisible detail — and the error and not-ready views, which clear the map, still answer `overList(0,0)` and `overDetail(0,0)` as true. The tabs are cleared, so AC-01 holds; the list and the detail are not. T2b should give the empty box a real shape.
- The wheel over the list moves the **selection** by three, so it can change the recipe the split layout's detail pane shows. It is the keyboard's own move, so it is one answer, but the plan's table says "scroll the list, keeping the selection visible" and the implementation moves the cursor.
- The detail's row count has two answers: the pane's width (the split render) and `m.width - 4` (`detailLineCount`). Folded into F-106 — the fix has to answer it once.
