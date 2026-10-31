You are the **reviewer** for **T3b** — cycle 3, the **last** one. The same lane that
reviewed cycles 1 and 2. The orchestrator's pane id is `w1H:p1`; use it verbatim in the
notify command below.

Read [work/tui/README.md](../tui/README.md) for the state and the protocol, and the
slice's briefs: [tui-t3b.md](../briefs/tui-t3b.md) ·
[fix-tui-t3b-1.md](../briefs/fix-tui-t3b-1.md) ·
[fix-tui-t3b-2.md](../briefs/fix-tui-t3b-2.md).

## Cycle 3 is the **last** cycle

**Anything still open after it is a note for the ledger**, not another round: the
orchestrator escalates to the chef. File **HIGH** blockers only — a MEDIUM finding at
this point is a **note for the ledger**, because the cycle cap says the slice ends and
a note is how a finding ends without a fourth round.

## What cycle 2 decided

Cycle 2's verdict was **`changes-needed`**, with **F-120** filed — and it caught an
**orchestrator's ruling** in the act: cycle 1's review had noted, as a LOW note, that
`docs/05`'s job list promised the console's free edit, the orchestrator ruled it into
a fix cycle, and the fix made the table speak in **two voices** — one row in build
state, three in remit.

**The orchestrator's ruling, corrected: the table is the REMIT.** It says what the
console's job *is*; what is *built* is the keys table below it. The build state leaves
`docs/05` and stays in the plan. The same table, the same voice, in
`docs/decisions/0003-surface-split.md` — which **stays**, in the vocabulary of its day.

The fix cycle landed: `c0f1624` · `89eaecc`.

## What to verify — each by command, not by reading

### 1. F-120 — the job list speaks in one voice

- **One line above the table** says what the table is — the console's remit, whether
  or not it is built — and that **what is built is the keys table below**.
- **Every row is in that voice.** The backlog's row carries **no build state**: the
  clause *"the free edit of an item's fields is not built"* is **out** of `docs/05`,
  and it stays in the plan.
- The table may name what the console is **for** and is not built yet. That is what a
  remit is.

### 2. F-121 — every kind the prompt offers is reachable

- The prompt's choices carry **distinct** keys: the kind's initial, then the shortest
  prefix no other kind has taken, with a fallback so no choice is ever unreachable.
- **Mutate it and prove the test bites:** the first byte and nothing else must turn the
  suite **red**.
- The **letter** is the console's own affordance — the domain owns the kind, the
  console owns the key. The rule belongs where the key is made.

### 3. The **Run** row — the orchestrator's answer to the runner's question

**The runner asked, and left the row untouched. Here is the orchestrator's answer:
the row stays.**

Every verb in it names something in the **chef's** voice:

| The verb | Whose |
|---|---|
| **start** | the **chef's decision** that the work starts — the sign-off that puts the recipe on the line. The **claim** is the agent's: `EventStart` routes to `ClaimRecipe`, and the console's own code says *"Claiming the work is an agent's move, and the console is not an agent — it must not record itself as one."* The console never records a claim, and the row does not ask it to |
| **pause**, **watch**, **steer**, **see what the agents are doing now** | the chef's controls over work in flight — built or not, and a remit may name what is not built |

So the row crosses no line: it names no move the console is **forbidden** to make, and
no **other agent's** work as the console's own. **The console's forbidden move — the
claim — is not in the row.**

**Check it against the code and say so if you disagree.** The orchestrator was wrong
once already in this slice, and the reviewer's correction was better than the ruling.

### 4. The slice's shape still holds

- The moves, the shaping, the dismissal, the create, the footer, the goldens, the
  scope, and `core.DomainVersion` (**7**).

## The note that waits for a later round

**Do not treat it as a finding, and do not expect it fixed:** `BacklogMoveOf` builds
the kinds list on **every call**, and the footer asks it every frame — a per-frame
allocation on the console's render path, which redraws the whole screen anyway. It is
a performance question about the render path, and the orchestrator has ruled it for a
later round.

## The mutation checks

**Mutate, run, revert**, and report what you mutated and what the suite said:

| # | Mutation | What must fail |
|--:|---|---|
| 1 | the console carrying its **own** kind list | the prompt's-choices tests |
| 2 | the key taken as the kind's **first byte** and nothing else | the distinct-keys test |
| 3 | the console deciding the move by itself | the footer's tests |

**Revert every mutation**, and **confirm the revert**.

## The ledger — `docs/audit.md` is yours

- **F-120 and F-121** get their rows: closed, with the commit that closed each and
  the test that holds it.
- The **slice T3b** row and the **cycle 1** and **cycle 2** rows are there. Add the
  **cycle 3** row.
- Anything you file is **F-121** onward — F-117, F-118, F-119 and F-120 are closed.
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
cat > work/inbox/reviewer-tui-t3b-cycle3.md <<EOF
Reviewer: T3b cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> HIGH blockers (<F-121.. or none>). F-120: <that the job list says once what it is, and that every row is in that voice>. F-121: <that the keys are distinct, and the mutation that fails it>. The Run row: <your judgement on the orchestrator's answer - the row stays>. The shape: <that the slice's criteria still hold>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t3b cycle 3"
herdr agent prompt "w1H:p1" "Reviewer: T3b cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> HIGH blockers (<F-121.. or none>). F-120: <that the job list says once what it is, and that every row is in that voice>. F-121: <that the keys are distinct, and the mutation that fails it>. The Run row: <your judgement on the orchestrator's answer - the row stays>. The shape: <that the slice's criteria still hold>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
