You are the **runner** for **T4** — the Home tab. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` first, then
[work/tui/plan.md](../tui/plan.md) § 2 for the reasoning, and
[work/tui/README.md](../tui/README.md) for the state, the protocol, and the **ruling**
this slice inherits.

**This is the first slice of the round's second half.** T1 gave the console its
identity, T2 the mouse, T3 the backlog. **T4 gives the console a landing place.**

## The chef's ruling

| # | Question | Answer |
|--:|---|---|
| **Q1** | the Home tab's name | **Home** |
| **Q2** | the Home tab's numbers | **the full set** — plan · recipes · evidence · tasting · backlog · automation |

## The precedent that must be answered first

v1 had an **Overview** lane, and v2 dropped it: *"the triage-ordered board **is** the
overview; a second tab would restate it."*

The chef's tab is not that tab — but the risk is the same, and the difference has to be
written down:

> **The invariant: the Home tab shows numbers; it shows no items.**
>
> No row on Home is selectable, `space` does nothing there, and the tab's item count
> is zero. **A test holds it**, the way slice 9's invariant held *"a recipe is filed
> under **Needs you** if and only if the console has a chef move for it"*.

**That test is the slice.** It is what stops Home from becoming a second board.

## Your slice — T4 — the Home tab

### AC-01 — the Home tab exists, and it is the landing tab

- A new **Home** tab, and the **Sous-Chef tab dissolves into it**: the tab count
  stays **7** — `Home · Recipes · Backlog · Path · Board · Log · Settings`. The
  Sous-Chef's content does not go away; it becomes Home's **automation** block.
- The tab opens **first** — Home is where the console lands.
- The tab's **name** is `Home` (Q1). The tab count does not change, and no other tab
  is renamed.

### AC-02 — the numbers, and only the numbers

The rows, and where each number comes from:

| Block | The number | Where it comes from |
|---|---|---|
| **The plan** | the project's name and the path | `core.Project` and the plan's root |
| **The recipes** | total · need you · on the line · served | the triage sections, counted |
| **The evidence** | criteria · proven · failed · rejected · unproven | `core.EvidenceSummaryOf`, summed |
| **The tasting** | open notes · open findings | `core.OpenFindingsOf` |
| **The backlog** | ideas · small tasks | `core.BacklogGroups` |
| **The automation** | the plan's domain version · this build's · the last journal event | the journal, and `core.DomainVersion` |

**Every number is computed in the domain.** A number the console computes itself is a
number the CLI would compute differently.

**The domain gains one answer, and it is a view.** The plan calls it **`core.PlanStats`**
— the numbers the Home tab shows. If you add it: the **`// view`** marker, a
registration in `core.DerivedViews`, and a **row in `docs/03` § 8**, whose name comes
from the runtime. The wire checks both directions, and it must **fail** when either is
removed.

**Say what you added and where it lives.** If the numbers fit the **existing** views
without a new one, that is a fine answer — say so, and say why.

### AC-03 — Home shows no items, and a test holds it

- **No row on Home is selectable**, `space` does nothing there, and the tab's item
  count is zero.
- **The test must be able to fail**: make a row selectable on Home, and the suite must
  go **red**.
- The chef's footers are promises — a footer that promises a move the console refuses
  is a lie (that is T3a's F-116). **Home's footer promises nothing.**

### AC-04 — the art moves to Home, and the content still wins

The mark moves from the Sous-Chef to Home — and slice G's rule stands unchanged:

> **The decoration yields to the content, in both dimensions.** The mark is drawn only
> when it fits, in width *and* in height, and at 80×24 the tab is what it would be
> without it.

- The art sits **beside** the numbers, not above them — the point of moving it is that
> it can share the space.
- **At 80×24 the numbers are whole and the art may be absent.** A slice that pushes
  the last number off the screen has failed, whatever the art looks like.
- The mark's drawings live in `internal/brand` and **nowhere else** — the brand's own
  wire holds that, and it will fail if you paste one.

### AC-05 — the goldens

**The goldens are expected to change** — this slice adds a tab and moves the art, and
the goldens hold the console's screen. `make goldens` regenerates, `make goldens-check`
holds, and the change is **claimed by a test**. Say what moved and why.

### AC-06 — the state row is updated

`work/tui/README.md`: T4 goes to `closed` with its commits, and **T5** — the Ready
tab — becomes the next session. Commit it with the slice.

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. **No dependency is added** — bubbletea already carries everything.
- **`core.DomainVersion` is 7.** The plan format does not move: the same fields, the
  same directories, the same file names. A new *view* is not a format change.
- The tab count is **7**, and the names are `Home · Recipes · Backlog · Path · Board ·
  Log · Settings`.
- **One fact, one answer** — a number the console shows is computed in the domain.
- **The domain layer has zero UI imports.**
- The gate is `make check`. A focused check is `make test PKG=… RUN=…`.

## The scope — what not to touch

| | |
|---|---|
| **The other tabs** | their content does not move. The Ready tab (the Path's new name) is **T5** |
| **The filter** | **T7** |
| **The Board** | **T6** |
| **The recipe's detail** | unchanged |
| **The backlog** | unchanged |
| **`maitre-cli`** | no verb is added and no command changes. The CLI may want the same numbers — **it does not get a verb for it** |
| **The docs** | `docs/05-human-interface.md` describes the console's screen: if this slice makes a line of it false, fix the line and say so. And the new view needs a row in `docs/03` § 8 |
| **`docs/audit.md`** | the ledger. **Never swept.** This slice files no finding |
| **`work/`** | the orchestrator's, except the state row (AC-06) |

## The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. If a file you did not
write has changed, **stop and report it**, and do not touch it.

## The evidence this cycle must produce

For **each** criterion, a command and its output — not a claim:

| Criterion | The evidence |
|---|---|
| **AC-01** | the console opens on **Home**, the tab count is **7**, the names are `Home · Recipes · Backlog · Path · Board · Log · Settings`, and the Sous-Chef's content is Home's **automation** block |
| **AC-02** | every row of the table above is on the screen with the number it promises — and **the console computes none of them**: show where each comes from |
| **AC-03** | the mutation — make a row on Home selectable — turns the suite **red** |
| **AC-04** | at **80×24** the numbers are **whole** and the art is **absent**; at a taller terminal the art is drawn **beside** them. Show both |
| **AC-05** | `make goldens` leaves a diff, and `make goldens-check` is green with the change **claimed** |
| **AC-06** | the state row names the slice's commits, and T5 is the next row |

Plus: `make check` green, `make check-race` green, and the tree clean.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t4.md <<EOF
Runner: T4 done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. The tab: <what Home is, and the seven names>. The numbers: <each row, and where its number comes from>. The domain's answer: <what you added, and whether it is a view with a row in docs/03 or the existing views>. No items: <the mutation that now fails it>. The art: <what it looks like at 80x24 and at a taller terminal>. The goldens: <what moved, and why>. The gate: <make check, make check-race>. Notes: <what T5 needs to know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t4"
herdr agent prompt "w1H:p1" "Runner: T4 done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. The tab: <what Home is, and the seven names>. The numbers: <each row, and where its number comes from>. The domain's answer: <what you added, and whether it is a view with a row in docs/03 or the existing views>. No items: <the mutation that now fails it>. The art: <what it looks like at 80x24 and at a taller terminal>. The goldens: <what moved, and why>. The gate: <make check, make check-race>. Notes: <what T5 needs to know>."
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
result, and what Home shows.
