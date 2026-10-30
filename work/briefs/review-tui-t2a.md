You are the **reviewer** for **T2a** — the mouse, part 1. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` first, then
[work/tui/plan.md](../tui/plan.md) § 3 for the reasoning behind the slice, and
[work/tui/README.md](../tui/README.md) for the state and the protocol.

Read [work/briefs/tui-t2a.md](tui-t2a.md) — the slice's contract.

**You are a fresh lane.** You did not write this code and you did not read it while it
was written. Re-derive the facts from the code — not from the runner's report, and not
from the brief's claims. A brief says what was *asked*; the code says what was *done*;
the difference is your job.

## The chef's ruling — the mouse

> *some tabs we can use the scroll, but the click does not work. also, we cant select
> the tabs using the mouse.*

**T2 is two sessions.** This one lands the **hit map and the wheel**; the **clicks are
T2b**. A click that is half-implemented is worse than no click.

## The slice's six acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the mouse is on — `tea.MouseModeCellMotion`, set on **every** view the console draws (the normal one, the error one, the not-ready one) |
| **AC-02** | the render builds the hit map — the tab labels, the list rows, the detail region, the banner rows — and the map is built **by the render**, not beside it |
| **AC-03** | the wheel scrolls: over the list, the list; over the detail, the detail — and it is the **same** scroll the keyboard drives |
| **AC-04** | the mouse is drivable from a script — `ui.Drive` grows a mouse step, and a scripted step draws the frame before it acts |
| **AC-05** | the goldens are expected to be **unchanged** |
| **AC-06** | `work/tui/README.md`'s state row is updated and committed with the slice |

## What to check — and how

**Each check below is a command, not a reading.** A check you cannot run, you do not
report as passed.

### 1. The mouse is on, on every view

- `Model.View` sets `v.MouseMode` on **every** path it returns from — the error view,
  the not-ready view, and the normal one. A view that does not ask for the mouse gets
  no mouse, and the terminal drops the events.
- The error and not-ready views must also **clear the map**: nothing is on screen, so
  no mouse event may land on a tab that is not there.

### 2. The map is built **by the render**

This is the slice. The renderer is the only thing that knows where it put the tabs, the
rows and the detail.

- **The map is a fact of the drawn frame, and it is rebuilt on every render.** It is
  kept on the model because a mouse event arrives *after* the frame was drawn.
- **A second function that re-derives the geometry would be a second answer to one
  question** — the house's standing failure. Read the render path and say whether the
  boxes come from the drawing or from a re-derivation. `grep` for the geometry
  constants (`headerRows`, `footerRows`, `splitThreshold`) and check that each is read
  where it is written, not recomputed.
- **The banner moves the list down.** The refusal banner sits between the summary and
  the list, so the list starts lower when there is one. Check the map's rows move with
  it — a map that kept the old rows would put a click on the wrong line.
- **The detail pane's region** is the pane, not the whole body — and in the
  full-screen detail the **list is not on screen at all**.

### 3. The wheel scrolls, and it is the keyboard's scroll

- The wheel over the **list** scrolls the list; over the **detail**, the detail.
- **The same scroll the keyboard drives** — the same offset, the same clamping. A
  second scroll variable would be a second answer to one question. The runner's claim
  is that the wheel and N keyboard steps land on the same place: **check that claim
  yourself**, and check what the wheel's step is and whether the report's reason for
  it holds.
- **The wheel over the header does nothing.** The header is not the list.
- **An open prompt takes every input** — the wheel must not move the selection out
  from under a half-typed answer.
- **Prove it end to end.** Drive the console from the command line and watch the
  detail move:
  ```bash
  maitre-cli demo --dir /tmp/rev && cd /tmp/rev
  maitre --render --keys "enter" --plain > /tmp/d0.txt
  maitre --render --keys "enter,wheel:down" --plain > /tmp/d1.txt
  diff /tmp/d0.txt /tmp/d1.txt     # the detail moved
  ```
  Note that `--render` alone renders **every** tab, while `--render --keys` renders
  **one** frame — and that the demo plan's list is short enough to fit, so the
  *list* scroll needs a plan that overflows. Say which of your checks proved what.

### 4. The mouse is drivable from a script

- `ui.Drive` grows a mouse step. Check the **spelling** the runner chose, and whether
  it is consistent with the shape of the existing key names.
- **A scripted mouse step draws the frame before it acts.** A mouse event is answered
  against the frame the chef is looking at, so a step that acts against a stale map
  acts on the wrong thing.
- An **unknown** step is refused with a clear message — a name that is neither a key
  nor a mouse step.

### 5. The scope held

| | |
|---|---|
| **No click** | the click is **T2b**. There must be **no** click handler — not a partial one that looks finished. `grep` for it |
| **`maitre-cli`** | no verb is added and no command changes. The console binary's own `--keys` help text may have moved — that is not a verb |
| **The tabs** | seven, the same names |
| **The domain** | no domain change. The hit map is a fact about the *screen*, not about the plan, so it belongs to `internal/ui`. Check `internal/core` is untouched |
| **`core.DomainVersion`** | still **7** |
| **The budget** | **0 / 0 / 28** against **0 / 0 / 30** |
| **The goldens** | `make goldens` leaves `git diff` empty. A golden that moved for no stated reason is the slice saying it did more than it said |
| **`docs/05-human-interface.md`** | the slice added to it; check the lines it added are **true** — and that nothing else in it became false |
| **`docs/audit.md`** | the ledger. **Never swept.** This slice files no finding |
| **`work/`** | the orchestrator's, except the state row (AC-06) |

## The mutation checks — this is what makes a review worth something

A green suite proves the code passes the tests it was written with. It does not prove
the tests *can* fail. **Mutate, run, revert** — and report what you mutated and what
the suite said:

| # | Mutation | What must fail |
|--:|---|---|
| 1 | take the `MouseMode` off the **error** view | the test that says every view asks for the mouse |
| 2 | make the render stop emitting the map — or emit it from somewhere other than the drawing | the map's tests |
| 3 | make the wheel scroll a **second** variable instead of the keyboard's scroll | the wheel-is-the-keyboard's-scroll test |
| 4 | make the wheel over the **header** scroll the list | the header-does-nothing test |
| 5 | make the wheel act while a **prompt** is open | the open-prompt test |
| 6 | make a scripted mouse step act **without drawing first** | the drivable-from-a-script test |

**Revert every mutation.** The tree must be exactly as the runner left it, except
your ledger commit. **Confirm the revert** — `git status --short` and `git diff` at the
end. A reviewer that leaves a mutation behind has damaged the slice it reviewed.

## The ledger — `docs/audit.md` is yours

- File **HIGH or MEDIUM** findings. Number them from **F-106** — T1's last was
  **F-105**.
- **LOW findings: note them in your message, do not file them.** They go in the
  ledger for a later round — **not** another cycle.
- Add the **slice T2a** row to the ledger's slice table.
- **Never sweep the ledger.** Its rows keep the names as they were when each finding
  was filed.

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
cat > work/inbox/reviewer-tui-t2a.md <<EOF
Reviewer: T2a review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-106.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The map is built by the render: <yes/no, and how you know>. The wheel is the keyboard's scroll: <what you compared, and what it showed>. The scripted mouse: <the spelling, and whether it draws first>. The end-to-end check: <what moved, and what you could not prove this way>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t2a"
herdr agent prompt "w1H:p1" "Reviewer: T2a review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-106.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The map is built by the render: <yes/no, and how you know>. The wheel is the keyboard's scroll: <what you compared, and what it showed>. The scripted mouse: <the spelling, and whether it draws first>. The end-to-end check: <what moved, and what you could not prove this way>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
