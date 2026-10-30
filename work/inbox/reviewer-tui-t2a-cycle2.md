Reviewer: T2a cycle 2 done. Verdict: changes-needed. Filed 1 blocker (F-109).

**F-106: the detail scrolls in both layouts — proved.** The command:
`maitre --render --width 200 --keys "enter,esc"` against `--keys "enter,esc,wheel:down@(150,6)"`. The split layout is reached by opening the detail and backing out of it. The **detail column** (columns 91-200) changed by 6 lines — the recipe's first three rows gone, three blank rows at the bottom — and the **list column** (columns 1-90, trailing blanks trimmed) is **byte-identical**: the detail moved and the list beside it did not. The other direction: `wheel:down@(10,6)` moves the list and leaves the detail's **scroll** alone (the pane follows the selection and shows another recipe from its top — the console's model, the same as the keyboard's `down`).

**The bound:** it is counted at the width the pane draws with — `Model.detailPaneWidth()` (`internal/ui/detail.go`): `width-4` in the full-screen detail, `width - listColumn() - 3` in the split layout. `renderSplit`, `renderDetailFull` and `detailLineCount` all ask that one function, and the offset is applied in one place (`detailScrollRows`). **But nothing holds it** — see F-109.

**F-107: the position form survives the command line.** `maitre --render --keys "wheel:down@(12,5)"` exits 0; `--keys "tab,wheel:down@(12,5),space"` exits 0; `--keys "tab,down,space"` still exits 0. The **old** spelling is refused with a message a reader can act on: `maitre: malformed mouse step "wheel:down@12"` (exit 1). `docs/05` prints `maitre --keys "wheel:down@(12,5)"` — the form that works.

**One grammar, one place:** the splitter is `ui.SplitSteps` (`internal/ui/mouse.go`) and the parser is `parseMouseStep` — **the same file, side by side**; `cmd/maitre/main.go` calls `ui.SplitSteps` and nothing else, and no other code splits steps on a comma. The splitter counts the brackets, so the comma inside a position is not a separator; an unbalanced bracket keeps the step whole and the parser refuses it with the *malformed mouse step* message.

**F-108: a region that is not drawn answers for nothing — proved.** `hitBox` carries a `there` flag; the list layout has no detail region and the cleared map answers nothing at the origin. By command: `maitre --render --width 120 --keys "wheel:down@(0,0)"` is **byte-identical** to `--keys "down,down,down"` — the wheel at the origin scrolls the **list**. By mutation: taking the `there` guard out of `contains`/`empty` turns the suite **red** in three tests (`TestTheWheelOverTheFullScreenHeaderDoesNothing`, `TestTheOriginIsNotAPhantomDetail`, `TestAClearedMapAnswersNothing`).

**The reversed test: judged sound — I do not file it.** The rule is **one** rule — the detail answers over the detail, the list answers the rest, when the list is on screen — and it lives in one place (`scrollUnderPointer`), stated once in `docs/05`. The console does **not** now have two rules: the old rule ("the wheel answers what is under the pointer") and the fix cycle's own F-108 evidence ("a wheel at the origin over the list scrolls the list") **cannot both hold**, because the origin *is* the header's first cell — so the reversal was forced by the ruling, not chosen. The **full-screen detail** still ignores the wheel over the header (no list on screen, nothing to fall back to), pinned by `TestTheWheelOverTheFullScreenHeaderDoesNothing` and proved by command: `--keys "enter"` and `--keys "enter,wheel:down@(150,0)"` are byte-identical. The new behaviour is the one the orchestrator accepted, and it is the better one: a scroll gesture that dies over the header is a papercut, and the list is the console's one scrollable thing.

**F-109 (medium, filed) — the bound's guard does not bite.** `TestTheScrollsBoundIsThePanesWidth` asks whether the bound equals the pane's row count, but the **seeded** recipe's detail has the same row count at both widths, so the assertion holds whatever the bound is. The reviewer's mutation — `detailLineCount()` back on `m.width - 4`, the full-screen width, which is the pre-fix defect — leaves the **whole `internal/ui` suite green**. The consequence is real: with the bound at the full-screen width the split pane's last rows are unreachable (a probe at 200x30 with wrapping prose: the pane has 41 rows, the bound 36 — five rows the scroll can never reach; with the fix, 41 and 41).

**Mutation checks (four of the five bite; the fifth is F-109):**
| # | Mutation | What failed |
|--:|---|---|
| 1 | the offset out of the split layout's detail | `TestTheWheelScrollsTheSplitDetail` |
| 2 | the bound at the full-screen width | **nothing — the whole `internal/ui` suite is green** (F-109) |
| 3 | the splitter cuts a bracketed position at its comma | `TestThePositionFormSurvivesTheCommandLine`, `TestThePositionFormIsDrivable` |
| 4 | the zero box answers at the origin again | three tests (above) |
| 5 | the wheel over the detail scrolls the list | `TestTheWheelScrollsTheDetail`, `TestTheScriptedMouseTakesAPosition`, `TestTheWheelScrollsTheSplitDetail` |

Every mutation was reverted with `git checkout --`; the tree is exactly as the runner left it — `git status --short` is empty and `git diff` names only `docs/audit.md`, my ledger commit.

**The goldens: unchanged.** `make goldens` leaves `git diff` empty.

**The gate:** `make check` green and `make check-race` green. The scope held: **no click handler** anywhere (the click is T2b); `internal/core`, `internal/cli` and `docs/03` untouched; `core.DomainVersion` still **7**; the budget still **0 / 0 / 28**; the seven tabs kept their names; the state row updated in the fix cycle's commit.

**LOW notes (not filed):**
- `docs/05`'s wheel table says *"wheel | anywhere else | scroll the list"*, and in the **full-screen** detail the header *is* "anywhere else" with no list on screen — the wheel over the header does nothing there, and the prose (*"there the wheel scrolls the detail alone"*) does not say so either. One more line would make the rule whole; the behaviour itself is right and the prose is close.
- `detailScrollRows` clamps **and writes** `m.detailScroll` during a render — a render with a side effect on the model. It was there before this cycle (`renderDetailFull` clamped it inline); the helper only moved it. Worth knowing before T2b adds more to the render path.
- The wheel over the list in the split layout moves the **selection**, so the detail pane shows another recipe from its top — the detail does not scroll, but the pane changes. That is the console's model (the keyboard's `down` does the same), so it is one answer, not two.
