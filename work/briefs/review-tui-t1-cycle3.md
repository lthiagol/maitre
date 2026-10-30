You are the **reviewer** for **T1** — cycle 3. The same lane that reviewed cycles 1
and 2. The orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command
below.

Read [work/tui/README.md](../tui/README.md) for the state and the protocol, and the
three briefs the slice has: [tui-t1.md](tui-t1.md) ·
[fix-tui-t1-1.md](fix-tui-t1-1.md) · [fix-tui-t1-2.md](fix-tui-t1-2.md).

## Cycle 3 is **the last cycle** — and the cap says what happens next

**This is the last cycle for T1.** Anything still open after it is a **note for the
ledger**, not another round: the orchestrator escalates to the chef.

File **HIGH or MEDIUM** blockers. No nits, no cosmetics, no style, and **do not
re-file what cycles 1 and 2 closed**. A LOW note goes in the ledger for a later
round.

## What cycles 1 and 2 decided

| Cycle | Verdict | What it closed |
|--:|---|---|
| 1 | `ok`, 0 blockers | the six acceptance criteria held; five LOW notes filed |
| 2 | `changes-needed`, 1 blocker | **F-103** — the name's rule was written twice and only one copy was pinned. **F-104** and **F-105** were filed as LOW notes |

F-99..F-102 were cycle 1's four fix-cycle findings, and they are closed. **F-103..F-105
are what this cycle judges.** The fix cycle landed in three commits: `8a6866b` ·
`cf8c05e` · `4e67b31`.

## What to verify — each by command, not by reading

### 1. F-103 — the rule is in one place, and it is pinned

- The rule that answers *"what is this plan called?"* lives in **one** place, and
  `core.LoadProject` and `core.Store.ProjectName` both **ask** it.
- **The mutation fails.** Take the rule out of the one copy — **the suite must go
  red**, in both `internal/core` and `internal/ui`. The runner's claim is that it
  turned red in three tests; **reproduce it yourself** and confirm it.
- **The two surfaces agree** in all three states — the plan names itself · the file
  names it **empty** · the file is **missing**:

  ```bash
  maitre-cli plan          # the CLI's answer
  maitre --render --plain  # the console's answer
  ```

  The CLI's behaviour **changed** in this cycle — it used to print `""` where the
  console printed the project's directory name, and it now prints the domain's answer.
  That was allowed by the brief, and it is the point of the finding. **Check that the
  two now agree, and that the agreement is not a coincidence of the fixture.**
- **The writer is untouched.** `core.Init` still writes what it wrote, and
  `core.DomainVersion` is still **7** — `TestDomainShapeIsVersioned` passes.

### 2. F-104 — the path the console renders is sanitised

- A plan whose **directory** name carries a newline shows the Settings tab's `root`
  row as **one row**.
- The sanitising is **inside `humanPath`**, so every path the console draws goes
  through it by construction — not by a call site's care. Check that it is not one
  call site that happens to be sanitised.

### 3. F-105 — the state row's shape

- The row's **count is the number of ids it names**, and the commits that wrote the
  row are not among them — a commit cannot carry its own hash. Check the row against
  the log: **the count and the list agree**.
- The **protocol section** says so in one line, so the next writer knows the shape.
  That line is in `work/tui/README.md`; check it is there.

### 4. The scope held

| | |
|---|---|
| **The goldens** | unchanged — `make goldens` leaves `git diff` empty |
| **The dependency budget** | **0 / 0 / 28** against **0 / 0 / 30** |
| **The tabs** | seven, the same names |
| **`internal/core` and `internal/cli`** | stdlib-only; the domain layer has zero UI imports |
| **`docs/03`** | **untouched in this cycle.** The plan's name is not a derived view, so it does not belong in § 8 — and where the rule itself belongs in the spec is **the chef's call**, which the orchestrator will put to them |
| **`docs/05`** | the paragraph cycle 1 added stays |
| **`docs/audit.md`** | the ledger is **yours** |

## The one question this cycle must **not** answer

> *where the rule that answers "what is this plan called?" should live in the spec.*

Cycle 2 answered the neighbouring question — the name is not a derived view, so it
does not belong in `docs/03` § 8 — and noted that **the rule itself is nowhere in the
spec**. The runner's reading is that it belongs with the plan's anatomy. That is a
**chef's** call, and the orchestrator is putting it to them.

**Do not write it into `docs/03`.** Say in your report whether the rule is written
down anywhere, so the orchestrator knows what to put to the chef.

## The ledger — `docs/audit.md` is yours

- **F-103..F-105** get their rows: closed, with the commit that closed each and the
  test that holds it.
- The **slice T1** row and the **cycle 1** and **cycle 2** rows are there. Add the
  **cycle 3** row.
- Anything you file is **F-106** onward.
- **Never sweep the ledger.**

## The verdict

| | |
|---|---|
| No HIGH or MEDIUM findings | `verdict: ok` |
| Any HIGH or MEDIUM finding | `verdict: changes-needed` |

**Do not invent reasons to fail.** A slice that does what it said is `ok`.

## Hard rules

- Do **not** edit implementation code — **except** `docs/audit.md`, which is yours.
- Revert every mutation you make. The tree must be exactly as the runner left it,
  except your ledger commit — confirm it with `git status --short` and `git diff`.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
- **One commit** for the ledger. One commit for the receipt.

## Notify the orchestrator

**The receipt comes FIRST** — the file cannot be lost, the nudge is only the bell.

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-tui-t1-cycle3.md <<EOF
Reviewer: T1 cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-106.. or none>). F-103 verified: <the mutation you took, and the suite's answer>. The two surfaces: <what they showed, in all three states>. F-104: <what you proved it with>. F-105: <the row's count, its list, and whether they agree>. The rule in the spec: <written down anywhere, or not>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t1 cycle 3"
herdr agent prompt "w1H:p1" "Reviewer: T1 cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-106.. or none>). F-103 verified: <the mutation you took, and the suite's answer>. The two surfaces: <what they showed, in all three states>. F-104: <what you proved it with>. F-105: <the row's count, its list, and whether they agree>. The rule in the spec: <written down anywhere, or not>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
