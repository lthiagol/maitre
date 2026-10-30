You are the **reviewer** for **T2a** — cycle 3. The same lane that reviewed cycles 1
and 2. The orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command
below.

Read [work/tui/README.md](../tui/README.md) for the state and the protocol, and the
slice's briefs: [tui-t2a.md](tui-t2a.md) · [fix-tui-t2a-1.md](fix-tui-t2a-1.md) ·
[fix-tui-t2a-2.md](fix-tui-t2a-2.md).

## Cycle 3 is the **last** cycle

**This is the last cycle for T2a.** Anything still open after it is a **note for the
ledger**, not another round: the orchestrator escalates to the chef.

File **HIGH** blockers only. A MEDIUM finding at this point is a **note for the
ledger** — the cycle cap says the slice ends, and a note is how a finding ends without
a fourth round.

## What cycles 1 and 2 decided

| Cycle | Verdict | What it closed |
|--:|---|---|
| 1 | `changes-needed` | **F-106** (the split layout's detail never scrolled) · **F-107** (the position form could not survive the command line) · **F-108** (a region that is not drawn answered for one cell) |
| 2 | `changes-needed` | **F-109** — the bound's width was held by nothing, because the test compared two values that **both** went through the same function, and the fixture's detail was too short for the two widths to differ |

**The wheel's fallback is closed.** Cycle 2 judged the reversed test **sound**: the
list is the console's one scrollable thing, the wheel scrolls it everywhere except
over the detail, and the rule is stated once. **Do not reopen it.**

The fix cycle landed in three commits: `a332158` · `d8d7b95` · `22b9de6`.

## What to verify — each by command, not by reading

### 1. F-109 — the bound's width is held by a test that can fail

- The bound is counted at the width the pane **draws with**, and the scroll reaches
  the pane's **last row**.
- **The tripwire must bite.** Count the bound at the full-screen width again — the
  suite must go **red**, and it must fail **for the bound's width** and not for some
  unrelated number.
- **The test must not ask the code under test for the answer.** The old test compared
  two values that both went through `detailPaneWidth()`. The new one computes the
  pane's width **in the test**. Check that.
- **The fixture must make the two widths differ**, and the test must say so if they
  ever stop differing — a fixture that silently agreed again would make the tripwire
  vacuous a second time.

### 2. The document's wheel table reads whole

- Every row of the table is **true**, including the row that says the wheel scrolls
  the list *"anywhere else, while the list is on screen"*.
- The prose and the table agree — check that the **full-screen detail** is the case
  the table carves out, and that the console really does what the table says there.

### 3. The scope held

| | |
|---|---|
| **No click** | the click is **T2b**. There must be **no** click handler. `grep` for it |
| **The wheel's fallback** | unchanged since cycle 2 |
| **`maitre-cli`** | no verb is added and no command changes |
| **The tabs** | seven, the same names |
| **The domain** | `internal/core` untouched |
| **`core.DomainVersion`** | still **7** |
| **The budget** | **0 / 0 / 28** against **0 / 0 / 30** |
| **The goldens** | `make goldens` leaves `git diff` empty |
| **`docs/03`** | untouched |
| **`docs/audit.md`** | the ledger is **yours** |

## The notes the fix cycle left for a later round

**Do not treat these as findings.** The fix cycle was told to leave them, and a cycle
that widens drifts:

- `detailScrollRows` writes `m.detailScroll` during a **render** — the clamping. It
  moved there from `renderDetailFull`; it is not new.
- The wheel over the list in the split layout moves the **selection**, so the pane
  shows another recipe from its top and the detail does not scroll. That is the
  console's model.

**If you think either is a finding**, say so in your message — the ledger's notes table
is where a note goes, and it is the orchestrator who decides that.

## The ledger — `docs/audit.md` is yours

- The **slice T2a** row and the **cycle 1** and **cycle 2** rows are there. Add the
  **cycle 3** row.
- Anything you file is **F-108** onward — F-106 and F-107 were closed in cycle 1's
  fix, and F-109 was closed in cycle 2's fix.
- **Never sweep the ledger.**

## The verdict

| | |
|---|---|
| No HIGH finding | `verdict: ok` |
| Any HIGH finding | `verdict: changes-needed` |

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
cat > work/inbox/reviewer-tui-t2a-cycle3.md <<EOF
Reviewer: T2a cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> HIGH blockers (<F-108.. or none>). F-109: <the mutation you took, what failed, and that it failed for the bound's width>. The fixture: <how the two widths differ, and the row counts>. The document: <what the table says now, and whether every row is true>. The scope: <what you checked>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. Notes for the ledger: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t2a cycle 3"
herdr agent prompt "w1H:p1" "Reviewer: T2a cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> HIGH blockers (<F-108.. or none>). F-109: <the mutation you took, what failed, and that it failed for the bound's width>. The fixture: <how the two widths differ, and the row counts>. The document: <what the table says now, and whether every row is true>. The scope: <what you checked>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. Notes for the ledger: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
