You are the **reviewer** for **T3b** — cycle 2. The same lane that reviewed cycle 1.
The orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command below.

Read [work/tui/README.md](../tui/README.md) for the state and the protocol, and the
slice's briefs: [tui-t3b.md](../briefs/tui-t3b.md) ·
[fix-tui-t3b-1.md](../briefs/fix-tui-t3b-1.md).

## Cycle 2 is **blockers only**

File **HIGH or MEDIUM** blockers. No nits, no cosmetics, no style, and **do not
re-file what cycle 1 closed**. A LOW note goes in the ledger for a later round.

## What cycle 1 decided

Cycle 1's verdict was **`changes-needed`**, with **F-117** filed — and F-117 was the
**orchestrator's** file: the plan promised a slice the brief ruled out, in three
places, and the orchestrator had corrected two and missed the third.

The orchestrator has since corrected the third place (`work/tui/plan.md`, the opening
table) and the plan's *Promote* row, which said the kind change is `idea → feature` —
narrower than what the console does, where the chef picks the kind.

The fix cycle landed: `4e14f19` · `58069e7`.

## What to verify — each by command, not by reading

### 1. F-117 — no doc promises the free edit

- Check the orchestrator's correction in `work/tui/plan.md` — the opening table and
  the *Promote* row.
- **Check that no doc in the tree promises the console edits an item's fields.** The
  runner checked nine files and found the promise in `docs/05` (fixed) and in
  `docs/decisions/0003-surface-split.md` — which the runner ruled **stays**.
- **The orchestrator's ruling on the ADR: it stays.** It is a **decision of record**
  about the division of labour — the console manages the backlog, the agent authors
  the recipe — not a statement about what is built today; it is written in the
  vocabulary of its day, and ADR-0009 is the decision of record that replaced those
  words. And the console **does** change backlog items — the kind, the scope, the body
  — so *"edit"* there is not a promise of a five-field form. **If you disagree with
  that ruling, say so in your report** — it is a judgement, and the orchestrator has
  ruled it, not discovered it.

### 2. F-118 — the console's doc says what the console does

- `docs/05-human-interface.md`'s job list says what the console does. **Every row of
  it must be true of the console today.**

### 3. F-119 — the kinds come from the domain

- The kinds the shape prompt offers come from the **domain** — through the move
  answer the console already asks — and the console carries **no list of its own**.
- **Mutate it and prove the test bites:** a console-side kind list that disagrees
  with the domain's must turn the suite **red**.
- The wire has nothing new to check, because no new exported function was added.
  **Confirm that** — that the answer carries the kinds, and that `docs/03` § 8's row
  says so.

### 4. The slice's shape still holds

- The moves, the shaping, the dismissal, the create, the footer, the goldens, the
  scope, and `core.DomainVersion` (**7**).

## The notes that wait for a later round

**Do not treat these as findings, and do not expect them fixed** — the orchestrator
ruled each, with reasons:

- **After a write the cursor keeps its index, and the re-grouping moves the item.**
  Pre-existing and **console-wide**: a sign-off on the Recipes tab moves the recipe
  between the triage sections the same way. It is a console-wide question about the
  selection, not this slice's.
- **The shaping test's id assertion is near-tautological** — the id does not change
  because the domain never changes an item's id. It is a guard, and the guard is worth
  having.
- **The moves act while the detail pane is open.** That is the console's model: the
  moves act on the selection, and the detail shows the selection — and the recipe
  path does the same.

## The mutation checks

**Mutate, run, revert**, and report what you mutated and what the suite said:

| # | Mutation | What must fail |
|--:|---|---|
| 1 | the console carrying its **own** kind list for the shape prompt | the prompt's-choices tests |
| 2 | the console deciding the move by itself | the footer's tests |
| 3 | the view's registration, or its row in `docs/03` § 8, removed | the wire |

**Revert every mutation**, and **confirm the revert**.

## The ledger — `docs/audit.md` is yours

- **F-117, F-118 and F-119** get their rows: closed, with the commit that closed each
  and the test that holds it. **F-117 was closed by the orchestrator's correction** —
  say so in its row.
- The **slice T3b** row and the **cycle 1** row are there. Add the **cycle 2** row.
- Anything you file is **F-120** onward.
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
cat > work/inbox/reviewer-tui-t3b-cycle2.md <<EOF
Reviewer: T3b cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-120.. or none>). F-117: <that no doc promises the free edit, and whether the ADR's staying is right>. F-118: <that the job list is true of the console today>. F-119: <that the kinds come from the domain, and the mutation that fails it>. The shape: <that the slice's criteria still hold>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t3b cycle 2"
herdr agent prompt "w1H:p1" "Reviewer: T3b cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-120.. or none>). F-117: <that no doc promises the free edit, and whether the ADR's staying is right>. F-118: <that the job list is true of the console today>. F-119: <that the kinds come from the domain, and the mutation that fails it>. The shape: <that the slice's criteria still hold>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
