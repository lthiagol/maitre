You are the **runner**, again — the same lane that wrote T2a. The orchestrator's pane id
is `w1H:p1`; use it verbatim in the notify command below.

This is **T2a's fix cycle 1**: **two blockers and one note the orchestrator ruled in**.

## The review's verdict

`changes-needed`, **two blockers: F-106 and F-107**. The slice's shape held — the map
**is** built by the render, the wheel **is** the keyboard's scroll, and all six of the
review's mutations failed the suite. Two things did not.

## F-106 (medium, blocker) — the split layout's detail never scrolls

**The target:** the detail pane scrolls the same way in **both** layouts — the
full-screen detail and the detail beside the list.

**The hole, proved by command.** The offset is applied in two places:

| Where | What it does |
|---|---|
| `internal/ui/views.go` — `renderDetailFull` | applies `m.detailScroll` to the rows it draws |
| `internal/ui/detail.go` — `viewDetail` | applies it too |
| **`internal/ui/views.go` — `renderSplit`** | **does not** — it draws `detailRows(m.detailLines(detailWidth))` and never looks at the offset |

So in the **split** layout the wheel over the detail pane moves `m.detailScroll` from 0
to 3 and the frame comes out **byte-for-byte identical**. The scroll state moves and
the screen does not. And the slice's map **advertises** that region — `renderSplit`
records the detail's box — so the pane looks like it should answer and does not.

**There is a second half, and it is the same finding.** The scroll's bound is
computed against a **different width** than the pane draws with:

| | |
|---|---|
| `renderSplit` draws the detail at | `m.width - listColumn() - 3` |
| `detailLineCount()` computes the bound at | `m.width - 4` — the **full-screen** width |

A narrower pane wraps the same text into **more** rows, so the bound is computed
against a row count the pane never has. The detail's row count has **two answers**,
and the scroll stops early because of it.

**What lands:** one answer for the detail's scroll, in both layouts. The offset is
applied where the detail is drawn — **once**, not twice — and the bound is computed
from the width the pane actually draws with.

**The evidence this finding demands:** the wheel over the detail pane moves the
**split** layout's detail, and the frame **differs**. Say in your report how you
proved it — the review notes that `--render` never draws the split layout, so a
command from the command line cannot show it and **a probe in the package can**.

## F-107 (medium, blocker) — the position form is unreachable, and the doc prints it

**The target:** the mouse's position form is reachable **through `maitre --keys`** —
the way every other step is — and **whatever `docs/05` prints works**.

**The hole, proved by command.**

```bash
maitre --render --keys "wheel:down@12,5" --plain
# maitre: unknown step "wheel:down@12"
```

The steps are **comma-separated**, and the coordinates are **comma-separated** too —
so the command line splits the step in half before the driver ever sees it. The
position form exists in the driver's API and cannot be used from the console's own
flag.

And `docs/05-human-interface.md` prints exactly that command:

```
maitre --keys "wheel:down"        # the wheel, over the list
maitre --keys "wheel:down@12,5"   # the wheel, at column 12, row 5
```

**The doc prints a command that fails.** A doc that contradicts the code is a bug.

**What lands:** a spelling that survives the comma. **The choice is yours** — the
coordinate separator may change, the step separator may change, the CLI may split
steps more carefully — but two things must hold:

1. **`maitre --keys` can reach the position form.** The command must work.
2. **`docs/05` prints a command that works.** Whichever spelling you choose, the doc
   shows **that** spelling.

**Say in your report which spelling you chose and why** — the shape of the existing
key names (`"tab"`, `"space"`, `"up"`) is the guide, and `ui.Drive` takes a
`[]string`, so a test can use any spelling you like.

## F-108 (low — the orchestrator ruled it in) — the zero box is not empty

**The target:** a region that is **absent** is absent.

The review filed this as a note, and the orchestrator ruled it into this cycle — it is
in the slice's own **targeting** path, and it is cheap.

**The hole.** `hitBox`'s zero value is **not** empty:

```go
func (b hitBox) contains(x, y int) bool {
	return x >= b.x0 && x <= b.x1 && y >= b.y0 && y <= b.y1
}
func (b hitBox) empty() bool { return b.x1 < b.x0 || b.y1 < b.y0 }
```

With every field zero, `contains(0, 0)` is **true** — the zero value is a one-cell
box at the origin. So:

- the **list** layout's detail is a phantom one-cell region at `(0, 0)`, and a wheel
  at exactly the origin is treated as *over the detail*;
- the **cleared** map — the one the error view and the not-ready view install — still
  answers `overList(0, 0)` **true**.

**What lands:** an absent region is absent. The list layout has no detail region, and
a cleared map answers nothing. **The test must be able to fail**: a wheel at the
origin over the list scrolls the list.

## What lands

| | |
|---|---|
| **`internal/ui`** | the detail's scroll, in both layouts, from one answer (F-106); a spelling the command line can carry (F-107); an absent region is absent (F-108) |
| **`docs/05-human-interface.md`** | the command it prints must work (F-107) |
| **`work/tui/README.md`** | the state row, which says `closed` and is not — this cycle is open (AC) |

## The scope — what not to touch

| | |
|---|---|
| **The clicks** | **T2b.** No click handler. A partial click that looks finished is worse than no click |
| **`maitre-cli`** | no verb is added and no command changes |
| **The tabs** | seven, the same names |
| **The domain** | no domain change. `internal/core` is untouched — the hit map is a fact about the *screen*, not about the plan |
| **`core.DomainVersion`** | stays **7**. This cycle changes no format |
| **The goldens** | `make goldens` regenerates; `make goldens-check` holds them. **Expected unchanged** — a golden that moved for no stated reason is the slice saying it did more than it said |
| **`docs/03`** | untouched |
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

For **each** finding, a command and its output — not a claim:

| Finding | The evidence |
|---|---|
| **F-106** | the wheel over the detail pane moves the **split** layout's detail and the frame **differs**. Say how you proved it — `--render` never draws the split layout, so a probe **in the package** is the instrument. And the scroll's **bound** comes from the width the pane draws with |
| **F-107** | `maitre --render --keys "<the position form>" --plain` **works**, and `docs/05` prints **that** form |
| **F-108** | a wheel at the origin over the **list** scrolls the **list** — the phantom detail cell is gone |

Plus: `make check` green, `make check-race` green, the goldens unchanged, and the tree
clean.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t2a-fix1.md <<EOF
Runner: T2a fix 1 done. Commits: <n>. Fixed: F-106, F-107, F-108. The detail's scroll: <what you changed, and how you proved the split layout moves>. The row count's one answer: <the width the bound uses now>. The position form: <the spelling you chose, why, and the command that works>. The doc: <the command docs/05 prints now>. The absent region: <what you proved it with>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <what T2b needs to know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t2a fix 1"
herdr agent prompt "w1H:p1" "Runner: T2a fix 1 done. Commits: <n>. Fixed: F-106, F-107, F-108. The detail's scroll: <what you changed, and how you proved the split layout moves>. The row count's one answer: <the width the bound uses now>. The position form: <the spelling you chose, why, and the command that works>. The doc: <the command docs/05 prints now>. The absent region: <what you proved it with>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <what T2b needs to know>."
```

**Write the receipt, commit it, and *then* nudge** — the nudge only reaches an idle
receiver; the receipt is the half that cannot be lost.

**One write per file, and commit as soon as the file is coherent.** Every edit to a
large file re-renders that file's whole diff into the terminal, and a turn that renders
dozens of diffs runs out before the work is done.

## Before you report

```bash
make check && git status --short
```

Green, and the tree clean. Report the **real** numbers: the commits, the gate's
result, and the command that works.
