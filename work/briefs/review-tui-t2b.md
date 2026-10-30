You are the **reviewer** for **T2b** — the mouse, part 2. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` first, then
[work/tui/plan.md](../tui/plan.md) § 3 for the reasoning, and
[work/tui/README.md](../tui/README.md) for the state and the protocol.

Read [work/briefs/tui-t2b.md](tui-t2b.md) — the slice's contract — and
[tui-t2a.md](tui-t2a.md) with its two fix briefs, because **T2b plays against what
T2a left**.

**You are a fresh lane.** You did not write this code and you did not read it while
it was written. Re-derive the facts from the code — not from the runner's report, and
not from the brief's claims.

## The chef's ruling — the mouse

> *some tabs we can use the scroll, but the click does not work. also, we cant select
> the tabs using the mouse.*

**T2 was two sessions.** T2a landed the mouse mode, the **hit map** and the wheel.
This one lands the **clicks**.

## The slice's seven acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | a click on a tab label switches tab — through the **same** switch the keyboard's `tab` makes |
| **AC-02** | a click on a row selects it — the item index from the **map**, and a header row selects nothing |
| **AC-03** | a click on a row shows that recipe's detail: in the split layout the detail follows the selection; in the **narrow** layout the click also **opens** it, because there is no room beside the list and a click that only selected would show nothing |
| **AC-04** | a click on the detail pane is the console's **`enter`** on that recipe — the detail takes the screen, `esc` brings it back. The wheel already follows the pointer, so the click has nothing left to mean for the wheel, and **no new state** is added: the console keeps two notions of where the action goes, the mode and the pointer |
| **AC-05** | the click is drivable from a script, in the **same grammar** as the wheel: `click@(x,y)`, and it survives `ui.SplitSteps` |
| **AC-06** | the goldens are expected to be **unchanged** |
| **AC-07** | `work/tui/README.md`'s state row is updated and committed with the slice |

## What to check — and how

**Each check below is a command, not a reading.** A check you cannot run, you do not
report as passed.

### 1. The clicks are the keyboard's moves

- A click on a tab label goes through the **same function** the keyboard's `tab` uses —
  the same `m.tab = …`, the same `clampSelections()`, the same status reset. Read the
  code and say which function it is.
- A click on a row selects the item **the map holds** — the index the render drew the
  row from — and not one recomputed from the list.
- A click on a **header** row (the triage sections, the counts, the rule under them)
  selects **nothing**.
- A click on the **detail pane** is `enter` on that recipe.

### 2. The two layouts, and the Backlog

- In the **split layout** (a wide terminal) a click on a row selects it and the
  detail **follows the selection**.
- In the **narrow** layout the detail is not on screen at all, so the click also
  **opens** it. **Check both layouts by command** — the width decides, and a check
  that only runs one layout proves half.
- On the **Backlog** tab a click selects and **stops**: the Backlog has no detail
  until T3a. A click must not build a backlog detail.
- On the **Settings** tab a click on a theme row **selects** it and does **not**
  apply the theme — that is `enter`'s move, and it is a write.

### 3. The click is drivable from a script

- `click@(x,y)` survives the command line: the steps are comma-separated and the
  position is **bracketed**, so the comma inside a position is not a separator.
  **Prove it**: `maitre --render --keys "click@(14,2)"` must work.
- **One grammar, one place.** The splitter that cuts the steps apart and the parser
  that reads a step must not be able to disagree. Check where each lives.
- The click is answered against the **last drawn frame** — the driver draws before it
  clicks.
- An **unknown** step is refused with a message a reader can act on.

### 4. The scope held

| | |
|---|---|
| **The wheel** | T2a's rule is **closed** — the detail where the pointer is over it, the list everywhere else, and the full-screen detail does nothing. The click and the wheel must **not** disagree about where the pointer is |
| **The map** | T2a's. If a click needed a region the map does not have, the map was **extended** — say so, and check the extension |
| **The CLI** | `maitre-cli` does not change. `maitre`'s own `--keys` help text may have gained a word |
| **The tabs** | seven, the same names |
| **The domain** | `internal/core` untouched — the hit map and the click are facts about the *screen*, not about the plan |
| **`core.DomainVersion`** | still **7** |
| **The budget** | **0 / 0 / 28** against **0 / 0 / 30** |
| **The goldens** | `make goldens` leaves `git diff` empty |
| **`docs/05`** | the click's moves are in the interaction table, and every line of it is **true**. The slice also corrected a **pre-existing** error in the key table — check that correction is right |
| **`docs/audit.md`** | the ledger. **Never swept.** This slice files no finding |
| **`work/`** | the orchestrator's, except the state row (AC-07) |

## The mutation checks — this is what makes a review worth something

A green suite proves the code passes the tests it was written with. It does not prove
the tests *can* fail. **Mutate, run, revert** — and report what you mutated and what
the suite said:

| # | Mutation | What must fail |
|--:|---|---|
| 1 | make the click set the tab **without** the switch's other work — no clamping, no status reset | the tab-switch test |
| 2 | make the click recompute the row's item **from the list** instead of taking the map's | the row-select test |
| 3 | make the click on a **header** row select something | the header test |
| 4 | make the narrow layout's click **not** open the detail | the narrow-layout test |
| 5 | make the click on the **detail** do nothing | the detail-click test |
| 6 | make the splitter cut a **bracketed** click at its comma | the position test |

**Revert every mutation**, and **confirm the revert** — `git status --short` and
`git diff` at the end. The tree must be exactly as the runner left it, except your
ledger commit.

## One question the runner raised

> *the console answers the press and not the release*

Check it and say whether it matters: a **drag** — press, move, release — is a click
whose position is the **release** position. If the console answers the press, a drag
would act on the row the chef pressed; if it answers the release, it acts on the row
the chef let go of. Say which the code does, and whether the answer is right.

## One defect the orchestrator found

**`internal/ui/mouse.go`'s doc comment carries the click's two example lines
**twice** — a copy-paste slip:

```
//	"click"             the click, over the list
//	"click@(12,5)"      the click, at column 12, row 5
//	"click"             the click, over the list      ← the slip
//	"click@(12,5)"      the click, at column 12, row 5 ← the slip
```

It is cosmetic. **Note it** — it goes in the ledger's notes table, and the next lane
that touches that file tidies it. It is **not** a blocker and it does not need a
cycle.

## The ledger — `docs/audit.md` is yours

- File **HIGH or MEDIUM** findings. Number them from **F-110** — T2a's last was
  **F-109**.
- **LOW findings: note them in your message, do not file them.**
- Add the **slice T2b** row to the ledger's slice table.
- **Never sweep the ledger.**

## The verdict

| | |
|---|---|
| No HIGH or MEDIUM findings | `verdict: ok` |
| Any HIGH or MEDIUM finding | `verdict: changes-needed` |

**Do not invent reasons to fail.** A slice that does what it said is `ok`.

## Hard rules

- Do **not** edit implementation code — **except** `docs/audit.md`, which is yours.
- Revert every mutation you make.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
- **One commit** for the ledger. One commit for the receipt.

## Notify the orchestrator

**The receipt comes FIRST** — the file cannot be lost, the nudge is only the bell.

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-tui-t2b.md <<EOF
Reviewer: T2b review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-110.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The clicks are the keyboard's moves: <which function each goes through>. The two layouts: <what each does, and how you proved it>. The Backlog and the theme: <what they do>. The grammar: <where the splitter and the parser live, and that click@(x,y) survives>. The press or the release: <which, and whether it is right>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t2b"
herdr agent prompt "w1H:p1" "Reviewer: T2b review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-110.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The clicks are the keyboard's moves: <which function each goes through>. The two layouts: <what each does, and how you proved it>. The Backlog and the theme: <what they do>. The grammar: <where the splitter and the parser live, and that click@(x,y) survives>. The press or the release: <which, and whether it is right>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
