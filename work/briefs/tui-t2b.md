You are the **runner** for **T2b** — the mouse, part 2. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` first, then
[work/tui/plan.md](../tui/plan.md) § 3 for the reasoning, and
[work/tui/README.md](../tui/README.md) for the state and the protocol.

**T2 was two sessions.** **T2a** landed the mouse mode, the **hit map** and the
wheel. This session lands the **clicks** — on a tab, on a row, on the detail pane.

**Read T2a's brief and T2a's commits before you start:** [tui-t2a.md](tui-t2a.md). The
map it built is the instrument this slice plays.

## The chef's ruling — the mouse

> *some tabs we can use the scroll, but the click does not work. also, we cant select
> the tabs using the mouse.*

## Your slice — T2b — the clicks

### AC-01 — a click on a tab label switches tab

The tab bar's labels are in the hit map with their column ranges. A click inside a
label's range switches to that tab.

- The switch is the **same** switch the keyboard's `tab` makes: the same
  `m.tab = …`, the same `clampSelections()`, the same status reset. A click that
  sets the tab without the rest of the switch is a second answer to one question.
- The click lands on the label, not on the gap between labels. The tab bar joins its
  parts with a separator, and the separator belongs to no tab.

### AC-02 — a click on a row selects it

A click on a list row puts the cursor on that row's item.

- The row's item index comes from the hit map, which carries the index the render
  used — not one recomputed from the list. Two indexes are two answers.
- A click on a **header** row (the triage sections, the counts) selects nothing: the
  map marks those rows as holding no item, and a header is not a thing to select.
- Selecting a row is the keyboard's `↓`/`↑` landing on it — the same selection
  state, the same clamping, the same scroll-into-view. A click that selects and
  scrolls by a different rule is a second answer to one question.

### AC-03 — a click on a row opens its detail — where a detail exists

On the tabs that have a detail (Recipes, Path, Board), a click selects **and** opens
the detail pane, the way `enter` does.

- **The Backlog tab has no detail yet** — that is **T3a**. So on the Backlog tab a
  click selects and stops there. Say so in the report; do not build a detail for the
  backlog in this slice.
- The Settings tab has no detail either. A click on a theme row selects it; it does
  not apply the theme — that is `enter`'s move, and it is a write.

### AC-04 — a click on the detail pane focuses it

The detail pane is a region in the map. A click inside it gives it the focus, so the
wheel scrolls **the detail** rather than the list.

- The console has one pointer. Saying "the detail is focused" is saying where the
  wheel goes. Today the wheel goes to the list always (T2a) — this slice gives the
  console a notion of **which pane has the pointer**.
- The state is one boolean, and it must not contradict the mode: the full-screen
  detail (`esc` from it) has no list beside it, so the pointer is in the detail and
  the wheel scrolls it.

### AC-05 — the mouse is drivable from a script

T2a gave `ui.Drive` a mouse step. This slice drives the **clicks** with it: a test
clicks a tab label and asserts the tab changed; a test clicks a row and asserts the
selection moved.

- The click is answered against the **last drawn frame** (T2a's rule). The driver
  draws before it clicks.

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
| **The wheel** | landed in T2a. If the wheel and a click disagree about where the pointer is, fix it — and say so |
| **The backlog's detail** | **T3a.** The Backlog tab has no detail pane yet |
| **The CLI** | `maitre-cli` does not change |
| **The tabs** | no tab is added, removed or renamed |
| **The domain** | no domain change. If you believe one is needed, **stop and report it** |
| **The docs** | `docs/05-human-interface.md` § *Interaction* describes the keys; if this slice makes a line of it false, fix the line and say so |
| **`docs/audit.md`** | the ledger. **Never swept.** This slice files no finding |
| **`work/`** | the orchestrator's, except the state row (AC-07) |

## The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. If a file you did not
write has changed, **stop and report it**.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t2b.md <<EOF
Runner: T2b done. Commits: <n>. Per-AC evidence: AC-01..AC-07 stamped. Pre-notify gate: pass. A click on a tab: <what it does>. A click on a row: <what it does>. The goldens: <unchanged / changed, and why>. Notes: <what the next session should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t2b"
```

**Write the receipt, commit it, and *then* nudge** — if a lane is waiting for it.

**One write per file, and commit as soon as the file is coherent.**

## Before you report

```bash
make check && git status --short
```

Green, and the tree clean. Report the **real** numbers: the commits, the gate's
result, and what each click does.
