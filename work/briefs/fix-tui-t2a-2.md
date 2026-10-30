You are the **runner**, again — the same lane that wrote T2a and its fix cycle 1. The
orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command below.

This is **T2a's fix cycle 2**: **one blocker and one line of a document**.

## The review's verdict

`changes-needed`, **one blocker: F-109**. Cycle 2 verified F-106, F-107 and F-108
by command and by mutation, and it **judged the reversed wheel test sound** — the
fallback is one rule, in one place, stated once in `docs/05`. That question is
closed; do not reopen it.

## F-109 (medium, blocker) — the bound's width is not held by anything

**The target:** the scroll stops at the pane's last row, and the width the bound is
counted at is **held by a test that fails when the width moves**.

**The hole, proved by command.** The fix cycle put the bound and the pane on one
answer — `Model.detailPaneWidth()` — and added a test, `TestTheScrollsBoundIsThePanesWidth`,
that checks the bound equals the pane's rows. **The reviewer's mutation leaves the whole
`internal/ui` suite green:**

```bash
# the BOUND counted at the full-screen width again
# internal/core/model.go:  detailLineCount() → m.width - 4
go test ./internal/ui/      # GREEN — every test passes
```

**Why the test does not catch it.** Its first assertion compares two values that
**both** go through the same function:

```go
pane := len(detailRows(m.detailLines(m.detailPaneWidth())))
if m.detailLineCount() != pane { ... }
```

Both sides ask `detailPaneWidth()` — so when the **bound** is computed somewhere else,
the comparison only notices **if the two widths give different row counts**. And the
fixture's detail is too short for that: at the split pane's width and at the full
screen's width it wraps to the **same** number of rows. The test passes with the bound
on the wrong width.

**What lands:** a test that **pins the width**, with a fixture where the two widths
**differ** — a detail long enough that the narrower pane wraps it into more rows.
Then the bound computed at the full-screen width **fails**.

**The evidence this finding demands:**

| | |
|---|---|
| the **mutation fails** | count the bound at the full-screen width, and the suite must go **red**. Say which test failed and what it said |
| the **bound is the pane's** | the scroll reaches the pane's **last row** in the **split** layout — not a row short of it |
| the **widths differ** | your fixture makes the split pane's width and the full screen's width give **different** row counts — say how you know, and how many |

**Watch the trap:** the test must fail because the **bound's width** moved, not
because some unrelated number moved. A test that fails for the wrong reason is not a
tripwire.

## The one line of the document

**The target:** the wheel's rule is stated **whole**.

`docs/05-human-interface.md`'s wheel table says:

| Gesture | Where | What |
|---|---|---|
| wheel | over the detail pane | scroll the detail |
| wheel | anywhere else | scroll the list, keeping the selection visible |

**"Anywhere else" is not true of the full-screen detail**: there is no list on screen
there, so the wheel over the header does **nothing**. The prose below the table says
so — the table does not.

**What lands:** the table says what is true, and the rule reads whole. **One line**,
or one row. Say in your report what you wrote.

## The scope — what not to touch

| | |
|---|---|
| **The clicks** | **T2b.** No click handler |
| **The wheel's fallback** | **closed.** The list is the console's one scrollable thing, and the rule is stated once. Do not change it, and do not reopen it |
| **`maitre-cli`** | no verb is added and no command changes |
| **The tabs** | seven, the same names |
| **The domain** | no domain change. `internal/core` is untouched |
| **`core.DomainVersion`** | stays **7** |
| **The goldens** | expected unchanged — a golden that moved for no stated reason is the slice saying it did more than it said |
| **`docs/03`** | untouched |
| **`docs/audit.md`** | the ledger is the **reviewer's**. This cycle files no finding and touches no row |
| **`work/`** | the orchestrator's, except the state row |

## The notes the review filed, and left for a later round

**Do not fix these in this cycle.** They are notes, and a cycle that widens is a cycle
that drifts:

- `detailScrollRows` writes `m.detailScroll` during a **render** — the clamping. It
  moved there from `renderDetailFull`; it is not new.
- The wheel over the list in the split layout moves the **selection**, so the pane
  shows another recipe from its top and the detail does not scroll. That is the
  console's model.

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. **No dependency is added.**
- **The domain layer has zero UI imports.**
- **The tab count is 7** and the tab names do not change.
- **One fact, one answer.**
- The gate is `make check`. A focused check is `make test PKG=… RUN=…`.

## The evidence this cycle must produce

| Finding | The evidence |
|---|---|
| **F-109** | the mutation fails: counting the bound at the full-screen width turns the suite **red**, and the test that failed is about the **bound's width**. And your fixture makes the two widths give different row counts — say how many |
| **the document** | the wheel's rule reads whole, and every row of the table is **true** |

Plus: `make check` green, `make check-race` green, the goldens unchanged, and the tree
clean.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t2a-fix2.md <<EOF
Runner: T2a fix 2 done. Commits: <n>. Fixed: F-109 and the document's wheel table. The bound's width: <the mutation you took, the test that failed, and what it said>. The fixture: <how the two widths come to differ, and the row counts>. The document: <what you wrote, and the table it reads as now>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t2a fix 2"
herdr agent prompt "w1H:p1" "Runner: T2a fix 2 done. Commits: <n>. Fixed: F-109 and the document's wheel table. The bound's width: <the mutation you took, the test that failed, and what it said>. The fixture: <how the two widths come to differ, and the row counts>. The document: <what you wrote, and the table it reads as now>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
```

**Write the receipt, commit it, and *then* nudge** — the nudge only reaches an idle
receiver; the receipt is the half that cannot be lost.

**One write per file, and commit as soon as the file is coherent.**

## Before you report

```bash
make check && git status --short
```

Green, and the tree clean. Report the **real** numbers: the commits, the gate's
result, and which test caught the mutation.
