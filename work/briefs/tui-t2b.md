You are the **runner** for **T2b** — the mouse, part 2. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` first, then
[work/tui/plan.md](../tui/plan.md) § 3 for the reasoning, and
[work/tui/README.md](../tui/README.md) for the state and the protocol.

**T2 was two sessions.** **T2a** landed the mouse mode, the **hit map** and the
wheel. This session lands the **clicks** — on a tab, on a row, on the detail pane.

**Read T2a's briefs and its commits before you start:** [tui-t2a.md](tui-t2a.md) ·
[fix-tui-t2a-1.md](fix-tui-t2a-1.md) · [fix-tui-t2a-2.md](fix-tui-t2a-2.md). The map
T2a built is the instrument this slice plays, and **T2a's fix cycles changed it**.

## What T2a left behind — the shape you are playing against

| | |
|---|---|
| the map | `internal/ui/hit.go` — `Model.hits` holds `tabs`, `rows`, `list`, `detail` and `banner`; `tabAt(x, y)` and `rowAt(x, y)` look up |
| a row | holds the **item index the tab's own selection uses** — not one recomputed |
| a region | `hitBox` carries a **`there` flag**: a region that was never drawn is *not there* and answers for nothing. The zero box is not a one-cell box at the origin any more |
| the mouse | `internal/ui/mouse.go` — `handleMouse` answers the wheel and **ignores the mouse while a prompt is open**; there is **no click handler yet**, on purpose |
| the grammar | one file, one grammar: `ui.SplitSteps` cuts the command line's steps apart and `parseMouseStep` reads a step. The spelling is `wheel[:up\|down][@(x,y)]` — the position is **bracketed**, because the steps are comma-separated and the comma inside a position must not be one |
| the wheel | the **list** is the console's one scrollable thing: the wheel scrolls the detail **where the pointer is over it**, and the list everywhere else — **except** the full-screen detail, where there is no list and the wheel does nothing. That rule is **closed**; do not reopen it |

## The chef's ruling — the mouse

> *some tabs we can use the scroll, but the click does not work. also, we cant select
> the tabs using the mouse.*

## The console's vocabulary — the clicks are keyboard moves

A click is a move, and the console's moves are the keyboard's. **Name the move you
are making**, and make the same one the keyboard makes:

| The keyboard's | What it does |
|---|---|
| `tab` / `right` | the **next tab**: `m.tab = …`, `clampSelections()`, the status reset |
| `shift+tab` / `left` | the **previous tab** |
| `↑` `↓` | in the full-screen detail, the detail's scroll; otherwise the **selection** moves |
| `enter` | the **detail takes the screen** (`ModeDetail`) |
| `esc` | back to the split view |

There are only **two** modes — the split view (the default: the detail beside the list
when the terminal has room for it, the list alone when it does not) and the full-screen
detail. The pointer is a **third** thing, and it is the mouse's alone.

## Your slice — T2b — the clicks

### AC-01 — a click on a tab label switches tab

A click inside a tab label's range switches to that tab.

- It is the **same switch** the keyboard's `tab` makes: the same `m.tab = …`, the
  same `clampSelections()`, the same status reset. A click that sets the tab and
  skips the rest is a second answer to one question.
- The click lands on the **label**, not on the gap between labels. The tab bar joins
  its parts with a separator, and the separator belongs to no tab.

### AC-02 — a click on a row selects it

A click on a list row puts the cursor on that row's item.

- The row's item index comes from **the map** — the index the render used. Not one
  recomputed from the list: two indexes are two answers.
- A click on a **header** row (the triage sections, the counts) selects **nothing**:
  the map marks those rows as holding no item, and a header is not a thing to select.
- Selecting is the keyboard's `↑`/`↓` landing on it — the same selection state, the
  same clamping, the same scroll-into-view.

### AC-03 — a click on a row shows that recipe's detail

A click on a row selects it, and the **detail shows that recipe**.

- In the **split view with room beside the list**, the detail is already on screen:
  the click selects, and the detail follows the selection. **Say so in your report.**
- In the **narrow layout** — where the list is alone and the detail is not on screen
  at all — a click that only selected would show the chef **nothing**. There the click
  also **opens** the detail, the keyboard's `enter`. **Say which move you make in each
  layout, and why.**
- **The Backlog tab has no detail yet** — that is **T3a**. So on the Backlog tab a
  click selects and stops there. Do not build a backlog detail in this slice.
- The **Settings** tab has no detail either. A click on a theme row **selects** it;
  it does **not** apply the theme — that is `enter`'s move, and it is a write.

### AC-04 — a click on the detail pane is `enter` on that recipe

**The premise this brief started with is false, and the fix is the reason.** The first
draft of this brief said a click on the detail pane gives it the focus *"so the wheel
scrolls the detail rather than the list"* — and added *"today the wheel goes to the
list always"*.

**The wheel does not go to the list always.** T2a's fix cycles made the wheel follow
the **pointer**: it scrolls the detail where the pointer is over it, and the list
everywhere else. So the click on the detail has nothing left to mean *for the wheel* —
and a click that changes nothing is not a click.

**What lands:** a click inside the detail pane is the console's **`enter` on that
recipe** — the detail takes the screen, and `esc` brings it back. That is the
console's existing move, and it adds **no new state**: the console keeps two notions
of where the action goes (the **mode** and the **pointer**) and this adds no third
one.

**If you think a different meaning is better, say so in your report** — and say why.
This is a judgement, and the orchestrator has ruled it, not discovered it.

### AC-05 — the mouse is drivable from a script

T2a gave `ui.Drive` a mouse step. This slice drives the **clicks** with it.

- **The spelling.** The grammar is `wheel[:up|down][@(x,y)]`, and the position is
  **bracketed**. The click's spelling is yours to choose, and it must be **the same
  grammar**: `click@(x,y)` is the obvious one. **Whatever you choose, it must survive
  `ui.SplitSteps`** — the splitter that cuts the command line's steps apart — and it
  must be read by the same parser that reads the wheel's. **One grammar, one place.**
- A test clicks a **tab label** and asserts the tab changed; a test clicks a **row**
  and asserts the selection moved; and the click is answered against the **last drawn
  frame** — T2a's rule: the driver draws before it clicks.
- An **unknown** step is refused with a message a reader can act on.

### AC-06 — the goldens

**The goldens are expected to be unchanged.** A click changes what the console
*answers*, not what it draws. If a golden moved, something else changed the render:
find it and say so.

### AC-07 — the state row is updated

`work/tui/README.md`: T2b goes to `closed` with its commits, and the next session
becomes the first row that is not `closed`. Commit it with the slice.

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. **No dependency is added.**
- **`core.DomainVersion` is 7.** No format, no domain code.
- The tab count is 7 and the tab names do not change.
- **The console offers no move the gates would refuse.** A click is a move: clicking
  a row must not do what `space` would refuse to do.
- The gate is `make check`. A focused check is `make test PKG=… RUN=…`.

## The scope — what not to touch

| | |
|---|---|
| **The wheel** | landed in T2a and its rule is **closed**. If the wheel and a click disagree about where the pointer is, **fix it** — and say so |
| **The map** | T2a's. Extend it if a click needs a region it does not have — and say so |
| **The backlog's detail** | **T3a.** The Backlog tab has no detail pane yet |
| **The CLI** | `maitre-cli` does not change. `maitre`'s own `--keys` help text may need a word about the click — that is not a verb |
| **The tabs** | no tab is added, removed or renamed |
| **The domain** | no domain change. If you believe one is needed, **stop and report it** |
| **The docs** | `docs/05-human-interface.md` § *Interaction* describes the keys and the mouse; if this slice makes a line of it false, fix the line and say so |
| **`docs/audit.md`** | the ledger. **Never swept.** This slice files no finding |
| **`work/`** | the orchestrator's, except the state row (AC-07) |

## The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. If a file you did not
write has changed, **stop and report it**.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t2b.md <<EOF
Runner: T2b done. Commits: <n>. Per-AC evidence: AC-01..AC-07 stamped. Pre-notify gate: pass. A click on a tab: <what it does>. A click on a row: <what it does in each layout, and why>. A click on the detail: <what it does>. The click's spelling: <the syntax you chose, and that it survives the splitter>. The goldens: <unchanged / changed, and why>. Notes: <what the next session should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t2b"
herdr agent prompt "w1H:p1" "Runner: T2b done. Commits: <n>. Per-AC evidence: AC-01..AC-07 stamped. Pre-notify gate: pass. A click on a tab: <what it does>. A click on a row: <what it does in each layout, and why>. A click on the detail: <what it does>. The click's spelling: <the syntax you chose, and that it survives the splitter>. The goldens: <unchanged / changed, and why>. Notes: <what the next session should know>."
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
result, and what each click does.
