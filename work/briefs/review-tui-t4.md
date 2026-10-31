You are the **reviewer** for **T4** — the Home tab. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` first, then
[work/tui/plan.md](../tui/plan.md) § 2 for the reasoning, and
[work/tui/README.md](../tui/README.md) for the state, the protocol, and the ruling.

Read [work/briefs/tui-t4.md](tui-t4.md) — the slice's contract.

**You are a fresh lane.** You did not write this code and you did not read it while
it was written. Re-derive the facts from the code — not from the runner's report, and
not from the brief's claims.

## The chef's ruling, and the precedent this slice had to answer

| # | Question | Answer |
|--:|---|---|
| **Q1** | the Home tab's name | **Home** |
| **Q2** | the Home tab's numbers | **the full set** — plan · recipes · evidence · tasting · backlog · automation |

**The precedent:** v1 had an **Overview** lane, and v2 dropped it — *"the triage-ordered
board **is** the overview; a second tab would restate it."* The difference is written
down as an invariant:

> **The Home tab shows numbers; it shows no items.**
> No row on Home is selectable, `space` does nothing there, and the tab's item count
> is zero. A test holds it.

## The slice's six acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the Home tab exists and is the **landing** tab; the **Sous-Chef dissolves into it**; the tab count stays **7** |
| **AC-02** | the six blocks of numbers, **every one computed in the domain** — and the domain's answer is a **view**, registered, with a row in `docs/03` § 8 |
| **AC-03** | Home shows **no items**, and a test holds it |
| **AC-04** | the art moves to Home, sits **beside** the numbers, and **the content still wins** — at 80×24 the numbers are whole and the art may be absent |
| **AC-05** | the goldens are expected to **change** — and the change is **claimed by a test** |
| **AC-06** | `work/tui/README.md`'s state row is updated and committed with the slice |

## What to check — and how

**Each check below is a command, not a reading.** A check you cannot run, you do not
report as passed.

### 1. The tab

- The console **lands on Home**, the tab bar reads `▸ Home    Recipes    Backlog
  Path    Board    Log    Settings`, and the count is **7**. The **Sous-Chef** is gone,
  and its content is Home's **automation** block — not lost.
- **No other tab moved.** The Path, the Board, the Log and the Settings keep their
  content. (The Path's **rename** to *Ready* is **T5**.)

### 2. The numbers, and the domain's answer

- Every row of the six blocks is on the screen, with the number it promises.
- **The console computes none of them.** Find each number and say where it comes
  from — the plan's name and path, the recipes' counts, the evidence's five, the
  tasting's two, the backlog's two, and the automation's three.
- The domain's answer is a **view**: the **`// view`** marker, a registration in
  `core.DerivedViews`, and a **row in `docs/03` § 8** whose name comes from the
  runtime. **The wire must fail when either side goes** — remove the registration,
  and remove the row, and see.
- **The stored facts are not in the view.** The plan's name and its path are *stored*,
  not computed — check the view does not pretend otherwise.

### 3. No items — and the test that holds it

- No row on Home is selectable, `space` does nothing there, and the tab's item count
  is zero.
- **Mutate it and prove the test bites:** make a row on Home selectable, and the suite
  must go **red**. A Home row that could be pressed is the tab v1's Overview was
  dropped for.
- The **footer promises nothing** on Home — a footer that promises a move the console
  refuses is a lie (that is T3a's F-116).

### 4. The art, and the content that wins

- The art sits **beside** the numbers, not above them.
- **At 80×24 the numbers are whole and the art is absent** — the tab is what it would
  be without the art. Slice G's rule, unchanged: *the decoration yields to the
  content, in both dimensions*.
- **Check both dimensions separately** — the width and the height, each on its own.
  A mark that fits the width but not the height is still a mark that pushes content
  off the screen.
- The mark's drawings live in `internal/brand` and **nowhere else** — the brand's own
  wire holds that, and it will fail if a copy is pasted.

### 5. The goldens

- `make goldens-check` is green, the changed goldens are **claimed by a test**, and
  `git diff` after `make goldens` is empty.
- **Say what moved and why** — a golden that moved for no stated reason is the slice
  saying it did more than it said.

### 6. The scope

| | |
|---|---|
| **The other tabs** | their content did not move |
| **The CLI** | no verb added, no command changed |
| **`core.DomainVersion`** | still **7** — a new *view* is not a format change |
| **The budget** | **0 / 0 / 28** against **0 / 0 / 30** |
| **The docs** | `docs/05` and `AGENTS.md` say what the console is now — **check every line they changed is true**, and that no line was left stale |
| **`docs/03` § 8** | the new row is there, and the wire holds it |
| **`docs/audit.md`** | the ledger is **yours** |

## The mutation checks

**Mutate, run, revert**, and report what you mutated and what the suite said:

| # | Mutation | What must fail |
|--:|---|---|
| 1 | a row on Home **selectable** | the no-items test |
| 2 | the view's **registration** removed | the wire |
| 3 | the view's **row in `docs/03` § 8** removed | the wire |
| 4 | the art drawn at **80×24** | the art-needs-both-dimensions test |

**Revert every mutation**, and **confirm the revert** — `git status --short` and
`git diff` at the end. The tree must be exactly as the runner left it, except your
ledger commit.

## The ledger — `docs/audit.md` is yours

- The **slice T4** row goes in the ledger's slice table.
- File **HIGH or MEDIUM** findings. Number them from **F-122** — T3b's last was
  **F-121**.
- **LOW findings: note them in your message, do not file them.**
- **Never sweep the ledger.**

## The verdict

| | |
|---|---|
| No HIGH or MEDIUM findings | `verdict: ok` |
| Any HIGH or MEDIUM finding | `verdict: changes-needed` |

**Do not invent reasons to fail.** A slice that does what it said is `ok`.

## Hard rules

- Do **not** edit implementation code — **except** `docs/audit.md`, which is yours.
- Revert every mutation you make, and **confirm the revert**.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
- **One commit** for the ledger. One commit for the receipt.

## Notify the orchestrator

**The receipt comes FIRST** — the file cannot be lost, the nudge is only the bell.

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-tui-t4.md <<EOF
Reviewer: T4 review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-122.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The tab: <that the console lands on Home, the seven names, and that the Sous-Chef's content is there>. The numbers: <each row, and where its number comes from>. The domain's answer: <the view's registration, its row, and that the wire fails both ways>. No items: <the mutation that now fails it>. The art: <what it looks like at 80x24 and at a taller terminal, and which dimension governs>. The goldens: <what moved, and why>. The docs: <that the lines they changed are true>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t4"
herdr agent prompt "w1H:p1" "Reviewer: T4 review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-122.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The tab: <that the console lands on Home, the seven names, and that the Sous-Chef's content is there>. The numbers: <each row, and where its number comes from>. The domain's answer: <the view's registration, its row, and that the wire fails both ways>. No items: <the mutation that now fails it>. The art: <what it looks like at 80x24 and at a taller terminal, and which dimension governs>. The goldens: <what moved, and why>. The docs: <that the lines they changed are true>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
