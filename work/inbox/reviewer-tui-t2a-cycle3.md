Reviewer: T2a cycle 3 done. Verdict: ok. Filed 0 HIGH blockers (none).

**F-109: the tripwire bites.** The mutation: `internal/ui/model.go` — `detailLineCount()` counted at the full-screen width again (`m.width - 4` instead of `m.detailPaneWidth()`), which is the pre-fix defect. The whole `internal/ui` suite goes **red**, and it fails **for the bound's width** — one test, three messages:

```
--- FAIL: TestTheScrollsBoundIsThePanesWidth (0.30s)
    bound_test.go:77: the pane has 29 rows at 107 columns; the full screen has 27 at 196
    bound_test.go:81: the bound is 27 rows; the pane has 29 rows at 107 columns
    bound_test.go:90: the scroll stopped at 26, want the pane's last row 28
```

Nothing else in the suite failed, so the tripwire fails for the bound's width and not for an unrelated number. The mutation was reverted with `git checkout --`.

**The test does not ask the code under test for the answer.** The pane's width is computed **in the test** — `m.width - m.listColumn() - 3` — and the pane's rows are counted **at that width** (`detailRows(m.detailLines(paneWidth))`), not through `detailPaneWidth()`. The old test compared two values that both went through `detailPaneWidth()`; the new one states the width independently.

**The fixture.** `longDetailModel` (`internal/ui/bound_test.go`) builds a plan through the core operations — one recipe whose outcome is a long paragraph, so the narrower pane wraps it into more rows. The row counts, printed by the test itself: **29 rows at the pane's 107 columns, 27 at the full screen's 196** — two rows apart. The test fails out loud if the two widths ever agree again (*the fixture must make the two widths differ*), so it cannot go vacuous a second time. The scroll reaches the pane's **last row** (28 of 29), asserted by wheeling 500 times.

**The document: the table reads whole.** It now says:

| Gesture | Where | What |
|---|---|---|
| wheel | over the detail pane | scroll the detail |
| wheel | anywhere else, **while the list is on screen** | scroll the list, keeping the selection visible |

and the prose names the case the table carves out: *"The full-screen detail has no list beside it: there the wheel over the header does nothing at all, and the wheel over the detail scrolls it."* Every row is true, proved by command:
- the split layout, the wheel over the header at (150,0): the **list column changes** (the list scrolls — the list is on screen);
- the full-screen detail, the wheel over the header at (150,0): `--keys "enter"` and `--keys "enter,wheel:down@(150,0)"` are **byte-identical** — nothing happens, and the table's carve-out says so;
- the full-screen detail, the wheel over the detail at (150,6): the frame **changes** — the detail scrolls.

**The scope.** No click handler anywhere (`grep` for `click` finds only two comments; the click is T2b). The wheel's fallback is **untouched** since cycle 2 — the fix cycle changed only `internal/ui/bound_test.go`, `internal/ui/mouse_test.go` and `docs/05`. `maitre-cli` gained no verb and no command changed. The seven tabs kept their names. `internal/core`, `internal/cli` and `docs/03` are untouched since cycle 1. `core.DomainVersion` is still **7**. The budget is still **0 / 0 / 28** against **0 / 0 / 30**.

**The goldens: unchanged.** `make goldens` leaves `git diff` empty.

**The gate:** `make check` green and `make check-race` green. The tree is exactly as the runner left it: `git status --short` is empty and `git diff` names only `docs/audit.md`, my ledger commit.

**Notes for the ledger** (added to the notes table, not filed as findings — I do not think either is a finding):
- `detailScrollRows` writes `m.detailScroll` during a render (the clamping). It moved there from `renderDetailFull`; it is not new.
- The wheel over the list in the split layout moves the selection, so the pane shows another recipe from its top and the detail does not scroll. That is the console's model.

**The slice is done.** Nothing is open: F-106, F-107, F-108 and F-109 are all closed, and the ledger's open table holds no T2a finding.
