You are the **runner**, again — the same lane that wrote T2b. The orchestrator's pane id
is `w1H:p1`; use it verbatim in the notify command below.

This is **T2b's fix cycle 1**: **two findings, and two notes the orchestrator ruled
in**.

## The review's verdict

`ok`, **0 blockers**. The clicks are the keyboard's moves, the two layouts differ as
the brief ruled, the Backlog and the theme row do not over-reach, the grammar is one
grammar, and all six of the review's mutations failed the suite. **The behaviour is
right.**

## F-110 (medium) — the header-click test cannot fail for what it names

**The target:** the test that says *"a click on a header row selects nothing"* is
able to fail.

**The hole.** `TestAClickOnAHeaderRowSelectsNothing` clicks a header row and asserts
the cursor did not move. **But the cursor starts at 0** — so a header click that
selects item **0** leaves it exactly where it was, and the test cannot see it. The
reviewer's mutation — make a header click select item 0 — **left this test green**.

**What lands:** the cursor is somewhere else before the click, so a header click that
moves it **fails the test**. The test must be able to fail for the thing it names.

**The evidence:** move the cursor, click a header row, and the mutation must turn the
suite **red**.

## F-111 (medium) — the row-select test's second assertion means nothing

**The target:** every assertion in a test means something.

**The hole.** `TestAClickOnARowSelectsIt` ends with:

```go
rows := m.visibleRows()
if got := rows[m.selRecipes].ID; got != rows[1].ID {
	t.Errorf("the cursor is on %s, want %s", got, rows[1].ID)
}
```

The first assertion already established `m.selRecipes == 1`, so this one compares
`rows[1].ID` with `rows[1].ID` — it holds **whatever the map says**, and it would
hold even if the click had selected the wrong recipe. It reads like a second
guarantee and is not one.

**What lands:** the assertion means something, **or it goes**. If it goes, the test
still holds the map's contract — say so in your report. If it stays, it must be able
to fail.

**The evidence:** break the thing it claims to check, and the suite must go **red**.

## F-112 (low, ruled in) — the duplicated doc comment

**The target:** the file says what it means, once.

`internal/ui/mouse.go`'s doc comment prints the click's two example lines **twice** —
a copy-paste slip. The orchestrator found it and put it in the review brief as a
note; the review filed it as a note too. **This cycle touches that file, so this
cycle tidies it.** Four lines, one of them deleted.

## F-113 (low, ruled in) — the message names the wrong cause

**The target:** a message names **its own** cause, so the reader can act on it.

**The hole, proved by command:**

```bash
maitre --render --keys "click@(1,2" --plain
# maitre: unknown step "click@(1,2"
```

The step **is** known — it is a mouse step — and it is **malformed**: the bracket is
never closed. The message says *unknown step*, which sends the reader looking for a
spelling that does not exist, when the thing to fix is the **bracket**.

The house's rule, from the rename: **each cause names its own remedy.** The wheel's
malformed position already says *malformed mouse step*; the click's does not.

**What lands:** a malformed position is refused with a message that names **that**
cause — and the reader can tell a malformed step from an unknown one.

**The evidence:** the command above must name the malformed position, and the
unknown-step message must still be the unknown-step message.

## The notes the review left, and did not file

**Do not touch these in this cycle:**

- The **Log tab's** detail is unreachable in the narrow layout — for the mouse and
  for the keyboard alike. It is **pre-existing**, and it is not this slice's.
- The **press or the release** is settled: the console answers the **press**, and
  that is right — it agrees with the wheel, which follows the pointer, and a
  release-based rule would let a drag that ends over the detail open it by accident.

## The scope — what not to touch

| | |
|---|---|
| **The behaviour** | the clicks are right and the review proved it. This cycle is about **the evidence and the messages**, not about what a click does |
| **The wheel's rule** | **closed.** Do not reopen it |
| **The map** | T2a's. Do not extend it for this cycle |
| **The clicks** | the tab click, the row click and the detail click — **unchanged** |
| **The CLI** | `maitre-cli` does not change |
| **The tabs** | seven, the same names |
| **The domain** | no domain change. `internal/core` is untouched |
| **`core.DomainVersion`** | stays **7** |
| **The goldens** | expected unchanged — a golden that moved for no stated reason is the slice saying it did more than it said |
| **`docs/05`** | unchanged unless a line of it became false — and nothing here changes the console's behaviour |
| **`docs/audit.md`** | the ledger is the **reviewer's**. This cycle files no finding and touches no row |
| **`work/`** | the orchestrator's, except the state row |

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
| **F-110** | the mutation — a header click that selects item 0 — turns the suite **red**, and the test that failed is the **header-click** test |
| **F-111** | the mutation — break what the second assertion claims to check — turns the suite **red**. Say what you broke, and which assertion caught it |
| **F-112** | the doc comment prints the click's examples **once** |
| **F-113** | `maitre --render --keys "click@(1,2" --plain` names the **malformed position**; an unknown step still says unknown step |

Plus: `make check` green, `make check-race` green, the goldens unchanged, and the tree
clean.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t2b-fix1.md <<EOF
Runner: T2b fix 1 done. Commits: <n>. Fixed: F-110, F-111, F-112, F-113. The header test: <what you changed, and the mutation that now fails it>. The row test: <what the second assertion does now, or that it went, and the mutation that fails it>. The doc comment: <that it prints the examples once>. The message: <what each cause says now>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t2b fix 1"
herdr agent prompt "w1H:p1" "Runner: T2b fix 1 done. Commits: <n>. Fixed: F-110, F-111, F-112, F-113. The header test: <what you changed, and the mutation that now fails it>. The row test: <what the second assertion does now, or that it went, and the mutation that fails it>. The doc comment: <that it prints the examples once>. The message: <what each cause says now>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
```

**Write the receipt, commit it, and *then* nudge** — the nudge only reaches an idle
receiver; the receipt is the half that cannot be lost.

**One write per file, and commit as soon as the file is coherent.**

## Before you report

```bash
make check && git status --short
```

Green, and the tree clean. Report the **real** numbers: the commits, the gate's
result, and which test caught each mutation.
