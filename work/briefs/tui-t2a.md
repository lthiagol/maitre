You are the **runner** for **T2a** — the mouse, part 1. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` first, then
[work/tui/plan.md](../tui/plan.md) § 3 for the reasoning, and
[work/tui/README.md](../tui/README.md) for the state and the protocol.

**T2 is two sessions.** This one lands the **hit map and the wheel** — the risky
half, and the half everything else about the mouse needs. The clicks are **T2b**, a
separate session.

## The chef's ruling — the mouse

> *some tabs we can use the scroll, but the click does not work. also, we cant select
> the tabs using the mouse.*

**The console receives mouse events today and drops them.** `internal/ui` has no
`tea.MouseMsg` case anywhere. In bubbletea v2 the mouse mode is **per view**, not a
program option (that was v1):

```go
v := tea.NewView(m.renderLayout())
v.MouseMode = tea.MouseModeCellMotion   // ← the whole of the enabling
```

## Your slice — T2a — the hit map and the wheel

### AC-01 — the mouse is on

`Model.View` sets `v.MouseMode`. It must be `tea.MouseModeCellMotion`, so the
terminal reports motion and wheel events.

- The mode must be set on **every** view the console draws — the normal view and the
  error view (`views.go`, `Model.View` has an early return for `m.err` and for
  `!m.ready`).
- `tea.Mouse{X, Y, Button, Mod}` carries the coordinates; `MouseWheelUp` and
  `MouseWheelDown` are the wheel buttons.

### AC-02 — the render builds the hit map

**This is the slice.** The console must know *where it drew things*: the render is a
string, and the mouse sends coordinates.

The map holds:

| What | The box |
|---|---|
| Each tab label | its column range, on the tab bar's row |
| Each list row | its row range, and the item index that row holds |
| The detail pane | its region |
| The banner rows | their rows — a refusal banner sits between the summary and the list, so the list starts lower when it is there |

**The map is built by the render, not beside it.** A second function that re-derives
the tab bar's columns would be a second answer to one question — the house's own
standing failure. The renderer emits the boxes as it draws. The row the tab bar is
on, the column the list ends at and the row the list starts on are all facts the
renderer already has.

**One map, rebuilt on every render.** It is stored on the model, because the mouse
event arrives *after* the frame was drawn — and a click is answered against the
frame the chef is looking at.

**The geometry is the console's, so read it where it is written:**

| The fact | Where it lives today |
|---|---|
| the header's rows | `headerRows = 4` (`views.go:20`) |
| the footer's row | `footerRows = 1` |
| the list's column | `Model.listColumn()` |
| the split's threshold | `splitThreshold = 160` |
| the list's rows | `Model.listLines()` and `Model.offset` |

### AC-03 — the wheel scrolls

| The pointer is over | The wheel does |
|---|---|
| the list | scroll the list — and the selection stays visible, exactly as the keyboard's scroll does |
| the detail pane | scroll the detail pane |

- The keyboard's `pgup`/`pgdown` move by `Model.listAreaHeight()`. The wheel's step
  is its own — **say in the report what step you chose and why**. A wheel that
  scrolls one line is too slow; a wheel that scrolls a page is too fast.
- The scroll is the **same** scroll the keyboard drives — the same offset, the same
  clamping. A second scroll variable would be a second answer to one question.

### AC-04 — the mouse is drivable from a script

CI has no terminal, and the console is driven from a script by `ui.Drive`
(`internal/ui/keys.go`). **The mouse has to be drivable the same way.**

- `ui.Drive` grows a mouse step. Choose the spelling and say it in the report; the
  shape of the existing key names (`"tab"`, `"space"`, `"up"`) is the guide.
- A test drives a wheel over the list and asserts that the list scrolled.
- **The click is answered against the last drawn frame.** `Drive` renders at the end,
  today — so a scripted mouse step must **draw before it clicks**, or it clicks
  against a map from a frame that is not on screen. Make the driver draw between
  steps, and say how.

### AC-05 — the goldens

**The goldens are expected to be unchanged.** The mouse changes what the console
*answers*, not what it draws, and the golden is `StripANSI(m.View().Content)` — the
view's `MouseMode` is not in the content. If a golden moved, something else changed
the render: find it and say so.

### AC-06 — the state row is updated

`work/tui/README.md`: T2a goes to `closed` with its commits, and **T2b becomes the
next session**. The row is the resume record — commit it with the slice.

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. **No dependency is added** — bubbletea already carries the mouse.
- **`core.DomainVersion` is 7.** This slice changes no format and no domain code.
- The tab count is 7 and the tab names do not change.
- **The domain layer has zero UI imports.** The hit map is a fact about *the screen*,
  not about the plan — it belongs to `internal/ui`.
- The gate is `make check`. A focused check is `make test PKG=… RUN=…`.

## The scope — what not to touch

| | |
|---|---|
| **The clicks** | **T2b.** This session lands the map and the wheel. A click that is half-implemented is worse than no click: do not add a partial one and leave it looking finished |
| **The CLI** | `maitre-cli` does not change |
| **The tabs** | no tab is added, removed or renamed |
| **The domain** | no domain change. If you believe one is needed, **stop and report it** |
| **The docs** | `docs/05-human-interface.md` § *Interaction* describes the keys; if this slice makes a line of it false, fix the line and say so |
| **`docs/audit.md`** | the ledger. **Never swept.** This slice files no finding |
| **`work/`** | the orchestrator's, except the state row (AC-06) |

## The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. If a file you did not
write has changed, **stop and report it**.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t2a.md <<EOF
Runner: T2a done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. The wheel step: <what you chose, and why>. The scripted mouse spelling: <the syntax you chose>. The goldens: <unchanged / changed, and why>. Notes: <what T2b needs to know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t2a"
```

**Write the receipt, commit it, and *then* nudge** — if a lane is waiting for it.

**One write per file, and commit as soon as the file is coherent.**

## Before you report

```bash
make check && git status --short
```

Green, and the tree clean. Report the **real** numbers: the commits, the gate's
result, the wheel's step and the scripted spelling.
