You are the **reviewer** for **T3b** — the backlog's moves. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` first, then
[work/tui/plan.md](../tui/plan.md) § 4 for the reasoning, and
[work/tui/README.md](../tui/README.md) for the state and the protocol.

Read [work/briefs/tui-t3b.md](tui-t3b.md) — the slice's contract — and
[tui-t3a.md](tui-t3a.md), because T3a built the ground this slice stands on.

**You are a fresh lane.** You did not write this code and you did not read it while
it was written. Re-derive the facts from the code — not from the runner's report, and
not from the brief's claims.

## The slice's seven acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the **domain** answers what move the item has — one place, and the console asks it |
| **AC-02** | **shape it**: the kind changes to a small task, **a scope is added**, and the **id does not change** |
| **AC-03** | **dismiss**: with a **reason** — the domain appends it to the item's body |
| **AC-04** | **create**: the chef may write an item down. It starts as an **idea** |
| **AC-05** | the footer tells the truth — `space shape it` on an open idea, `space promote to recipe` on an open small task, and no `space` on a settled item |
| **AC-06** | the goldens are expected to be **unchanged** |
| **AC-07** | `work/tui/README.md`'s state row is updated and committed with the slice |

**The house rule that carries the slice:** the console **asks** the domain. A
console-side *"is it a small task?"* rule is the defect **F-114** just closed — do not
let it back in.

**And the line the console must not cross:** it never authors a **recipe**. A backlog
item is not a recipe, so creating one is allowed (the chef's ruling Q6); authoring a
recipe is not.

## What to check — and how

**Each check below is a command, not a reading.** A check you cannot run, you do not
report as passed.

### 1. The domain answers, and the console asks

- The move is computed in **one** place, in the domain, and it **reads**
  `core.PromotionOf` — so the console can never offer a promotion the domain refuses.
- The console has **no** rule of its own about kinds. Grep it and say what you find.
- If the answer is a **new exported function**, the wire applies: the **`// view`**
  marker, a registration in `core.DerivedViews`, and a row in `docs/03` § 8 — and the
  wire must **fail** when either is removed.
- **Mutate it and prove the test bites:** make the console decide for itself, and the
  suite must go **red**.

### 2. Shaping

- The kind changes to a small task, a **scope is added**, and the **id does not
  change** — `BL01` stays `BL01`. It is an **update**, not a promotion: **no recipe**
  is created, and the item stays open.
- The chef picks the kind; the domain's rules govern the rest.

### 3. Dismissing

- The chef's **reason** is recorded — the domain appends it to the item's body, so a
  dismissal **without** a reason would say nothing. **Ask for it.**

### 4. Creating

- The console writes a **backlog item** — and **nothing else**: no recipe, and nothing
  outside the plan's backlog.
- The new item starts as an **idea**.

### 5. The footer, and the scope

- The footer names the move the selected item has — **`space shape it`** on an open
  idea, **`space promote to recipe`** on an open small task, and **no `space`** on a
  settled item.
- `space` on an idea now **shapes** it — there is nothing to refuse. And the **CLI's**
  promotion of an idea is **still** refused, with the reason, and the journal still
  records that an agent asked.
- The domain's writes, the tabs, `core.DomainVersion` (**7**), the budget
  (**0 / 0 / 28**), the goldens, and `docs/05` — unchanged, or justified.

## The mutation checks

**Mutate, run, revert**, and report what you mutated and what the suite said:

| # | Mutation | What must fail |
|--:|---|---|
| 1 | the console deciding the move by itself | the footer's test |
| 2 | the domain's move answering **promote** for an open idea | the refusal test |
| 3 | the view's registration, or its row in `docs/03` § 8, removed | the wire |
| 4 | the shaping changing the **id** | the shaping test |

**Revert every mutation**, and **confirm the revert** — `git status --short` and
`git diff` at the end. The tree must be exactly as the runner left it, except your
ledger commit.

## Two things the runner flagged, and one of them is the orchestrator's

- **The plan promised a slice the brief ruled out.** `work/tui/plan.md` said T3b lands
  *"dismiss · create · edit"*, and its § 4 table promised the console editing title,
  body, scope, kind and tags. The brief ruled the **free edit** out: the three moves
  are what the chef needs and what T3a's refusal points at, and a five-field surface
  is a later slice, if it is wanted at all. **The orchestrator has corrected the plan
  and the round's rulings row** — check both corrections, and check that no doc still
  promises the free edit.
- **The scripted driver types one character per step**, so `--keys` can only type
  single letters. That is its old rule, unchanged. Say whether it holds the console's
  prompts to account — the *tests* drive the model directly, so the prompts' text steps
  are reachable there.

## The ledger — `docs/audit.md` is yours

- The **slice T3b** row, and the **cycle 1** row, go in the ledger's table.
- File **HIGH or MEDIUM** findings. Number them from **F-117** — T3a's last was
  **F-116**.
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
cat > work/inbox/reviewer-tui-t3b.md <<EOF
Reviewer: T3b review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-117.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The domain's answer: <where the move lives, and that the console asks it>. Shaping: <what changes, and that the id does not>. Dismiss: <that the reason is recorded>. Create: <that the console writes a backlog item and nothing else>. The footer: <what it says on an idea, on a small task, and on a settled item>. The plan's correction: <that no doc still promises the free edit>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t3b"
herdr agent prompt "w1H:p1" "Reviewer: T3b review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-117.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The domain's answer: <where the move lives, and that the console asks it>. Shaping: <what changes, and that the id does not>. Dismiss: <that the reason is recorded>. Create: <that the console writes a backlog item and nothing else>. The footer: <what it says on an idea, on a small task, and on a settled item>. The plan's correction: <that no doc still promises the free edit>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
