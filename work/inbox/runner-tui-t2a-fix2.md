Runner: T2a fix 2 done. Commits: 4 (a332158, d8d7b95, plus the state row and this receipt). Fixed: F-109 and the document's wheel table.

**The bound's width.** The old test compared two values that **both** went through `detailPaneWidth()`, so it could not notice the bound being counted elsewhere — and the fixture's detail was too short for the two widths to disagree, so the comparison was vacuous anyway. The new test pins the width:

* the pane's width is computed **in the test** — `m.width - m.listColumn() - 3`, the list's column, one column of air, and the rest of the terminal — and the pane's rows are counted **at that width**, not through the console's answer. A tripwire that asks the code under test for the answer is not a tripwire.
* the fixture makes the two widths **disagree**, and the test says so out loud: if they did not, it fails before it can pass vacuously.

The mutation I took, and what it said — `detailLineCount()` counted at the full-screen width again (`m.width - 4` instead of `m.detailPaneWidth()`), and the whole `internal/ui` suite goes **red**:

```
--- FAIL: TestTheScrollsBoundIsThePanesWidth (0.32s)
    bound_test.go:77: the pane has 29 rows at 107 columns; the full screen has 27 at 196
    bound_test.go:81: the bound is 27 rows; the pane has 29 rows at 107 columns
    bound_test.go:90: the scroll stopped at 26, want the pane's last row 28
```

The failing test is about the **bound's width** — the first line names the bound's row count against the pane's, the second says the scroll stopped two rows short of the pane's last row. Nothing else in the suite failed, so the tripwire fails for the right reason and not for an unrelated number.

**The fixture.** A recipe whose outcome is a long paragraph: the pane wraps it into more rows than the full screen does. The numbers, printed by the test itself:

| | width | rows |
|---|--:|--:|
| the split pane | 107 columns | **29** |
| the full screen | 196 columns | **27** |

Two rows apart — enough for the bound to be caught. The scroll reaches the pane's **last row** (row 28 of 29) in the split layout, which the test asserts by wheeling 500 times and reading the offset.

**The document.** The table's second row was "anywhere else", which was not true of the full-screen detail: there is no list on screen, so the wheel over the header does nothing. The row now reads:

| Gesture | Where | What |
|---|---|---|
| wheel | over the detail pane | scroll the detail |
| wheel | anywhere else, **while the list is on screen** | scroll the list, keeping the selection visible |

And the prose's last sentence says what the wheel does with the full-screen detail open: *there the wheel over the header does nothing at all, and the wheel over the detail scrolls it.* The table and the prose now agree, and every row of the table is true.

The wheel's **fallback** is untouched — one rule, in one place (`Model.scrollUnderPointer`), and the notes the review filed are left for a later round as it ruled them.

DomainVersion: 7. The goldens: **unchanged** — `make goldens-check` regenerates them and finds no diff; nothing in this cycle changes what the console draws. The gate: `make check` green, `make check-race` green. `internal/ui` still 28 external modules, no dependency added, `internal/core` and `internal/cli` untouched.

Notes for the reviewer:

* The tripwire's fixture is a **new plan**, built in the test (`longDetailModel`, internal/ui/bound_test.go) — one recipe with a long outcome. The seeded plan's detail is too short to wrap differently at the two widths, which is why the first test could not fail.
* The pane's width is now written **twice** — once in `Model.detailPaneWidth()` and once in the test, as `m.width - m.listColumn() - 3`. That duplication is the point of a tripwire: the test states the expected width independently of the code under test. If the geometry ever moves, the test is where the house learns about it.
* `go test ./internal/ui -run TestTheScrollsBoundIsThePanesWidth -v` prints the two row counts, so the fixture's margin is visible without reading the test.
