You are the **reviewer** for **T2b** — cycle 2. The same lane that reviewed cycle 1.
The orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command below.

Read [work/tui/README.md](../tui/README.md) for the state and the protocol, and the
slice's briefs: [tui-t2b.md](tui-t2b.md) · [fix-tui-t2b-1.md](fix-tui-t2b-1.md).

## Cycle 2 is **blockers only**

File **HIGH or MEDIUM** blockers. No nits, no cosmetics, no style, and **do not
re-file what cycle 1 closed**. A LOW note goes in the ledger for a later round.

## What cycle 1 decided

Cycle 1's verdict was **`ok`, 0 blockers** — the clicks are the keyboard's moves, the
two layouts differ as the brief ruled, the Backlog and the theme row do not over-reach,
the grammar is one grammar, and the press-not-the-release rule is right.

Cycle 1 filed five **notes**, and the orchestrator took **two** of them into a fix
cycle, because they were about the slice's own **evidence**:

| # | What |
|--:|---|
| **F-110** | the header-click test could not fail for what it names — the cursor started at **0**, so a header click that selected item 0 left it exactly where it was, and the reviewer's mutation left the test **green** |
| **F-111** | the row-select test's second assertion compared `rows[1].ID` with `rows[1].ID` — it held whatever the map said, and it read like a second guarantee without being one |
| **F-112** | *(ruled in)* the duplicated doc comment in `internal/ui/mouse.go` — this cycle touched that file, so this cycle tidied it |
| **F-113** | *(ruled in)* a malformed position was refused as *"unknown step"*, which sends the reader looking for a spelling that does not exist, when the thing to fix is the **bracket** |

The fix cycle landed in three commits: `0144216` · `6bb7545` · `f8969ef`.

## What to verify — each by command, not by reading

### 1. F-110 — the header-click test can fail

- The cursor is **somewhere else** before the click, so a header click that moves it
  is visible.
- **Mutate it and prove the test bites:** make a header row hold item 0 — the map's
  contract broken — and the suite must go **red**, and the **header-click test** must
  be the one that fails.

### 2. F-111 — the row-select test's second assertion means something

- The two sides must be read from **different places** — the frame and the model's
  rows — so the comparison can fail.
- **Mutate it and prove the test bites:** break the render's item counter, so the
  map's item index drifts from the drawn recipe, and the suite must go **red**.

### 3. F-112 — the doc comment

- The click's two example lines are printed **once**.

### 4. F-113 — the message names its own cause

```bash
maitre --render --keys "click@(1,2" --plain   # malformed mouse step "click@(1,2": the position must be (x,y)
maitre --render --keys "wibble" --plain       # unknown step "wibble"
```

- Each cause names **its own** remedy, and an unknown step still says unknown step.

## One question for you — the **unbracketed** position

The fix cycle gives the **bracketed** malformed position the remedy clause:

```
maitre --render --keys "click@(1,2"   →  malformed mouse step "click@(1,2": the position must be (x,y)
maitre --render --keys "wheel:down@12" →  malformed mouse step "wheel:down@12"
```

Both are the same cause with **the same remedy** — the position must be bracketed —
but only the first names it. The second names the cause and leaves the reader to work
out the bracket.

**Is that right?** Say so in your report. If it is a finding, file it — it is one
line. If it is a note, it goes in the ledger's notes table with the next round.

## The scope — what held

| | |
|---|---|
| **The behaviour** | the clicks did not change in this cycle. Check it: the tab click, the row click in both layouts, the detail click, the Backlog, the theme row |
| **The wheel's rule** | unchanged and **closed** |
| **The map** | T2a's, unextended |
| **`maitre-cli`** | unchanged |
| **The tabs** | seven, the same names |
| **The domain** | `internal/core` untouched |
| **`core.DomainVersion`** | still **7** |
| **The budget** | **0 / 0 / 28** against **0 / 0 / 30** |
| **The goldens** | `make goldens` leaves `git diff` empty |
| **`docs/05`** | unchanged, and no line of it became false |
| **`docs/audit.md`** | the ledger is **yours** |

## The notes the review left, and the orchestrator's ruling

- The **Log tab's** detail is unreachable in the narrow layout — for the mouse and for
  the keyboard alike. **Pre-existing**, and not this slice's.
- The **press or the release** is settled: the press, and it is right.

## The mutation checks

**Mutate, run, revert**, and report what you mutated and what the suite said:

| # | Mutation | What must fail |
|--:|---|---|
| 1 | a header row holds item **0** | the header-click test |
| 2 | break the render's **item counter** | the row-select test's second assertion |
| 3 | make the malformed-position message say **"unknown step"** | the message test |
| 4 | print the click's example lines **twice** in the doc comment | whatever holds the doc comment |

**Revert every mutation**, and **confirm the revert**.

## The ledger — `docs/audit.md` is yours

- The **slice T2b** row and the **cycle 1** row are there. Add the **cycle 2** row.
- Anything you file is **F-110** onward — F-110 and F-111 were closed by cycle 1's
  fix, so the next free number is **F-112**.
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
cat > work/inbox/reviewer-tui-t2b-cycle2.md <<EOF
Reviewer: T2b cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-112.. or none>). F-110: <the mutation you took, and the test that failed>. F-111: <what the second assertion reads now, and the mutation that fails it>. F-112: <that the examples print once>. F-113: <what each cause says now>. The unbracketed position: <your judgement>. The behaviour: <that the clicks did not change, and what you checked>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t2b cycle 2"
herdr agent prompt "w1H:p1" "Reviewer: T2b cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-112.. or none>). F-110: <the mutation you took, and the test that failed>. F-111: <what the second assertion reads now, and the mutation that fails it>. F-112: <that the examples print once>. F-113: <what each cause says now>. The unbracketed position: <your judgement>. The behaviour: <that the clicks did not change, and what you checked>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
