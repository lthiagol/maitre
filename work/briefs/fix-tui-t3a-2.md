You are the **runner**, again — the same lane that wrote T3a. The orchestrator's pane id
is `w1H:p1`; use it verbatim in the notify command below.

This is **T3a's fix cycle 2**, and it is **small**: two things in the seam T3a
introduced.

## Why this cycle exists

Your fix cycle 1 asked a question, and the orchestrator has answered it: **the Log tab's
detail rides in this cycle.** You said it was one line in `hasDetail` and one test.
It is a little more than that, and the reason is worth writing down.

## F-116 (medium) — the key and the pane disagree on the Log tab, and the footer promises a key that does nothing

**The target:** the console's **one** detail pane and the key that opens it read
**one** answer, on every tab — and the footer never promises a key that does nothing.

**The hole.** T3a introduced `Model.hasDetail()` so that *"the pane and the key read
the same answer, so the key can never open something the pane does not show"*. On the
**Log** tab they do not:

| | The Log tab, today |
|---|---|
| the **pane** | shows the journal entry — `detailLines` has a `case TabLog:` and renders `m.logDetail(width)` |
| the **key** | `enter` opens **nothing** — `hasDetail()` has no `case TabLog:` |
| the **footer** | says **`enter detail`** — the default case adds that key on every tab but Settings |

So the footer promises a key that does nothing, and the pane shows something the key
cannot open. The same promise is wrong on the **Sous-Chef**: the pane says *"Nothing to
open here — the Sous-Chef is not a list of recipes"*, and the footer still says
`enter detail`.

**What lands — two changes, both in the seam:**

1. **`hasDetail()` answers for the Log tab.** The pane already renders a journal
   entry, so the key must be able to open it. One case.
2. **The footer offers `enter detail` only where there is a detail to open** — the
   same answer the key reads. So the console never promises a key that does nothing:
   not on the Log before (1), and not on the Sous-Chef, and not on a tab with nothing
   selected.

**The evidence:**

| | |
|---|---|
| the **Log** | `enter` on the Log tab opens the event's detail — in the **narrow** layout, where the detail is not beside the list. Show it by command: the frame changes, and it shows the event |
| the **footer** | on the **Log** tab the footer offers `enter detail`; on the **Sous-Chef** it does **not**; and on the **Recipes** tab with **nothing selected** it does not |
| **one answer** | the key and the pane agree on every tab — check them **together**, tab by tab |

## What does *not* ride in this cycle

**The detail's arrow on a promoted item** — the one that reads like an instruction
rather than a fact — **stays for a later round.** It wants a different marker for the
domain's answer than the recipe detail's *"why it needs the chef"* arrow, and that is a
question about the console's two arrow shapes, not a line.

**The backlog item's own absent lists** stay, as ruled — the note is in the ledger
with the corrected reason.

## The scope — what not to touch

| | |
|---|---|
| **The refusals** | **closed.** F-114 and F-115 are fixed and verified: the footer asks the domain and a test holds it, and a refused `space` writes nothing to the plan |
| **The domain** | `internal/core` is untouched by this cycle. If the seam needs a domain change, **stop and report it** — do not reach for it |
| **The Sous-Chef's pane** | its message — *"Nothing to open here"* — stays. This cycle changes the **footer's promise**, not that message |
| **`internal/cli`** | unchanged |
| **The tabs** | seven, the same names |
| **`core.DomainVersion`** | stays **7** |
| **The goldens** | expected unchanged — but **say so with the evidence**: the footer's key list is rendered text, and the goldens hold the **Recipes** tab, where the footer's keys do not change |
| **`docs/05`** | unchanged unless a line became false — and if the footer's promise changes, check that the doc does not promise it either |
| **`docs/audit.md`** | the ledger is the **reviewer's**. This cycle files no finding and touches no row |
| **`work/`** | the orchestrator's, except the state row |

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. **No dependency is added.**
- **The domain layer has zero UI imports.**
- **The tab count is 7** and the tab names do not change.
- **One fact, one answer** — a fact the console shows is computed in the domain.
- The gate is `make check`. A focused check is `make test PKG=… RUN=…`.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t3a-fix2.md <<EOF
Runner: T3a fix 2 done. Commits: <n>. Fixed: F-116. The Log's detail: <what enter opens there now, and the command that shows it>. The footer's promise: <what it offers on the Log, on the Sous-Chef, and on Recipes with nothing selected>. One answer: <the tabs you checked, and that the key and the pane agree>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t3a fix 2"
herdr agent prompt "w1H:p1" "Runner: T3a fix 2 done. Commits: <n>. Fixed: F-116. The Log's detail: <what enter opens there now, and the command that shows it>. The footer's promise: <what it offers on the Log, on the Sous-Chef, and on Recipes with nothing selected>. One answer: <the tabs you checked, and that the key and the pane agree>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
```

**Write the receipt, commit it, and *then* nudge** — the nudge only reaches an idle
receiver; the receipt is the half that cannot be lost.

**One write per file, and commit as soon as the file is coherent.**

**And if you find a file you did not write has changed: stop and report it**, and do
not touch it. It is most likely the reviewer's mutation, in flight.

## Before you report

```bash
make check && git status --short
```

Green, and the tree clean. Report the **real** numbers: the commits, the gate's
result, and what the footer promises on each tab.
