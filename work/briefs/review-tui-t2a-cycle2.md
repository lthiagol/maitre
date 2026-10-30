You are the **reviewer** for **T2a** — cycle 2. The same lane that reviewed cycle 1.
The orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command below.

Read [work/tui/README.md](../tui/README.md) for the state and the protocol, and the
slice's briefs: [tui-t2a.md](tui-t2a.md) · [fix-tui-t2a-1.md](fix-tui-t2a-1.md).

## Cycle 2 is **blockers only**

File **HIGH or MEDIUM** blockers. No nits, no cosmetics, no style, and **do not
re-file what cycle 1 closed**. A LOW note goes in the ledger for a later round.

## What cycle 1 decided

Cycle 1's verdict was **`changes-needed`**, with **F-106** and **F-107** filed and
**F-108** noted. The slice's shape held: the map **is** built by the render, the wheel
**is** the keyboard's scroll, and all six of cycle 1's mutations failed the suite.

The fix cycle landed in four commits: `619b342` · `47be828` · `f1cc1ca` · `90e4e96`.

## What to verify — each by command, not by reading

### 1. F-106 — the detail scrolls in both layouts

- The scroll's offset is applied **once**, where the detail is drawn, and **both** the
  full-screen detail and the split layout's detail show it.
- **The bound comes from the width the pane draws with.** The split pane draws at its
  own width and the full-screen detail at its own, and the scroll's limit is computed
  from the same one — the detail's row count must not have two answers.
- **Prove it end to end.** The split layout is reached by opening the detail and
  backing out of it:

  ```bash
  maitre-cli demo --dir /tmp/rev && cd /tmp/rev
  maitre --render --width 200 --keys "enter,esc" --plain > /tmp/p0.txt
  maitre --render --width 200 --keys "enter,esc,wheel:down@(150,6)" --plain > /tmp/p1.txt
  diff /tmp/p0.txt /tmp/p1.txt     # the detail moved, and the list beside it did not
  ```

  Say what moved and what did not. **And check the other direction**: a wheel over
  the **list** in the split layout must move the list and **not** the detail.

### 2. F-107 — the position form survives the command line

- `maitre --render --keys "wheel:down@(12,5)"` **works**, and the steps are still
  comma-separated — `--keys "tab,down,space"` must keep working.
- The **old** spelling is refused with a message a reader can act on, not a bare
  "unknown step".
- **`docs/05` prints a command that works.** That is the whole of the finding: the doc
  printed a command that failed.
- **One grammar, one place.** The splitter that cuts the steps apart and the parser
  that reads a step must not be able to disagree — check where each lives and whether
  one could drift from the other.

### 3. F-108 — a region that is not drawn answers for nothing

- The list layout has **no** detail region, and the cleared map — the one the error
  and the not-ready views install — answers **nothing** at the origin.
- **A wheel at the origin over the list scrolls the list.**
- **Mutate it and prove the test bites**: make the zero box answer again, and the
  suite must go red.

## The judgement this cycle must make — a **reversed** test

**The runner changed a behaviour in this fix cycle, and inverted a test that pinned
it.** Cycle 1's brief said *"the wheel over the header does nothing"*, and the slice
pinned it in `TestTheWheelOverTheHeaderDoesNothing`. That test is now
`TestTheWheelOverTheHeaderScrollsTheList`, and it asserts the **opposite**.

**Where it came from: the orchestrator's own brief.** Cycle 1's fix brief demanded the
evidence *"a wheel at the origin over the list scrolls the list"* — and the origin is
the header's first cell. With the origin no longer a phantom detail, the runner
generalised: **the list is the console's one scrollable thing, and the wheel scrolls
it everywhere except over the detail.**

**The orchestrator accepts the behaviour, and says why:** the wheel is a scroll
gesture, the list is the console's one scrollable thing, and a wheel that dies over
the header is a papercut. `docs/05` now states the fallback in its wheel table and in
prose.

**But the reversal is yours to judge.** A fix cycle that inverts a pinned assertion is
the thing a reviewer must look at hardest:

- Is the **new** behaviour right, and is `docs/05`'s statement of it **true**?
- Does the **full-screen** detail still ignore the wheel over the header? There is no
  list on screen there, so the fallback has nowhere to go.
- Is the old behaviour — *the wheel answers what is under the pointer* — still the
  console's model anywhere, or has the console now got **two** rules about what the
  wheel does? If it has two, one of them is a second answer to one question.

**Judge it and say so.** If you disagree with the orchestrator's ruling, file it —
with your reasoning, not as a nit.

## The scope — what held

| | |
|---|---|
| **No click** | the click is **T2b**. There must be **no** click handler — a partial click that looks finished is worse than no click. `grep` for it |
| **`maitre-cli`** | no verb is added and no command changes |
| **The tabs** | seven, the same names |
| **The domain** | `internal/core` untouched — the hit map is a fact about the *screen*, not about the plan |
| **`core.DomainVersion`** | still **7** |
| **The budget** | **0 / 0 / 28** against **0 / 0 / 30** |
| **The goldens** | `make goldens` leaves `git diff` empty |
| **`docs/03`** | untouched |
| **`docs/05`** | the lines this slice added and changed are **true** — check them, including the wheel's fallback and the bracketed position |
| **`docs/audit.md`** | the ledger is **yours** |

## The mutation checks

**Mutate, run, revert**, and report what you mutated and what the suite said:

| # | Mutation | What must fail |
|--:|---|---|
| 1 | take the offset out of the **split** layout's detail again | the split-detail test |
| 2 | compute the scroll's bound at the **full-screen** width again | the bound test |
| 3 | make the splitter cut a **bracketed** position at its comma | the position test |
| 4 | make the zero box answer at the origin again | the absent-region test |
| 5 | make the wheel over the **detail** scroll the **list** instead | the detail-targeting test |

**Revert every mutation**, and **confirm the revert** — `git status --short` and
`git diff` at the end. The tree must be exactly as the runner left it, except your
ledger commit.

## The ledger — `docs/audit.md` is yours

- File **HIGH or MEDIUM** findings. Number them from **F-106** — T2a's cycle 1 filed
  F-106 and F-107, and the fix cycle closed them, so the next free number is
  **F-108**.
- **LOW findings: note them in your message, do not file them.**
- The **slice T2a** row and the **cycle 1** row are there. Add the **cycle 2** row.
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
cat > work/inbox/reviewer-tui-t2a-cycle2.md <<EOF
Reviewer: T2a cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-108.. or none>). F-106: <what you proved, and the command you used>. The bound: <the width it uses now>. F-107: <the command that works, and the old spelling's message>. One grammar: <where the splitter and the parser live>. F-108: <what you proved it with>. The reversed test: <your judgement on the wheel's fallback>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t2a cycle 2"
herdr agent prompt "w1H:p1" "Reviewer: T2a cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-108.. or none>). F-106: <what you proved, and the command you used>. The bound: <the width it uses now>. F-107: <the command that works, and the old spelling's message>. One grammar: <where the splitter and the parser live>. F-108: <what you proved it with>. The reversed test: <your judgement on the wheel's fallback>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
