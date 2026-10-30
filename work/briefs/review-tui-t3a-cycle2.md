You are the **reviewer** for **T3a** — cycle 2. The same lane that reviewed cycle 1.
The orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command below.

Read [work/tui/README.md](../tui/README.md) for the state and the protocol, and the
slice's briefs: [tui-t3a.md](tui-t3a.md) · [fix-tui-t3a-1.md](fix-tui-t3a-1.md) ·
[fix-tui-t3a-2.md](fix-tui-t3a-2.md).

## Cycle 2 is **blockers only**

File **HIGH or MEDIUM** blockers. No nits, no cosmetics, no style, and **do not
re-file what cycle 1 closed**. A LOW note goes in the ledger for a later round.

## What cycle 1 decided

Cycle 1's verdict was **`changes-needed`**, with **F-114** and **F-115** filed and one
question answered against the orchestrator: a refused `space` must not write to the
plan's journal.

**F-115 (medium)** — the console asked the domain by *calling the operation*, and
the operation journals every refusal, so a misclick appended a `refusal` event to the
append-only journal while the footer, in the same frame, offered no move at all.

**F-114 (medium)** — the footer asked `core.PromotionOf`, but **nothing held it**:
the mutation that made the footer decide by the item's kind left the whole
`internal/ui` suite green, because the two cases the test held were the two where the
console's own rule and the domain's answer **agree**.

**The orchestrator then ruled one more item in** — **F-116** — from the runner's own
question: the Log tab's detail is unreachable by the key while the pane renders it,
and the footer promises `enter detail` on tabs that have none.

The two fix cycles landed: `aa4ba5f` · `da36146` · `99a7576` (cycle 1) and
`63964ce` (cycle 2).

## What to verify — each by command, not by reading

### 1. F-114 — the footer asks the domain, and a test holds it

- The test must use the case where the console's own rule and the domain's answer
  **differ** — the **promoted** small task, which the domain refuses because it is no
  longer open.
- **Mutate it and prove the test bites:** the footer deciding by the item's kind must
  turn the suite **red**, in that test.

### 2. F-115 — a refused `space` writes nothing

- `space` on an idea leaves `activity.json` **byte-for-byte unchanged** — show the
  checksums before and after — and the domain's reason still reaches the chef.
- The **CLI's** refusal still stands: an agent asked, the plan said no, and the
  journal records that attempt.

### 3. F-116 — the key and the pane agree on every tab

- `enter` on the **Log** tab opens the event's detail — in the **narrow** layout,
  where the detail is not beside the list.
- The footer offers `enter detail` **only** where there is a detail to open: not on
  the **Sous-Chef**, not on the **Recipes** tab with **nothing selected**, and yes on
  the **Log** and the tabs that hold a recipe.
- Check the tabs **together** — the key and the footer must agree, tab by tab.

### 4. The slice's shape still holds

- The item's detail, the domain's refusal, the footer's truth, the shape every writer
  writes, the digest, the view's registration and its row in `docs/03` § 8.
- **The goldens** are unchanged.

### 5. The scope held

| | |
|---|---|
| **The domain** | `internal/core` and `internal/cli` untouched by the two fix cycles |
| **`core.DomainVersion`** | still **7** |
| **The budget** | **0 / 0 / 28** against **0 / 0 / 30** |
| **The tabs** | seven, the same names |
| **`docs/05`** | no line of it became false — check that the interaction table still says what the console does |
| **`docs/03`, `docs/07`, `GLOSSARY.md`** | the corrections that landed are **right** — the docs must not say an idea can be promoted, and the glossary must not read as though it could |
| **`docs/audit.md`** | the ledger is **yours** |

## The notes that wait for a later round

**Do not treat these as findings, and do not expect them fixed:**

- The **detail's arrow** on a promoted item reads like an instruction rather than a
  fact; it wants a different marker from the recipe detail's *"why it needs the chef"*
  arrow. That is a question about two arrow shapes.
- The **backlog item's** own `scope`/`tags` are written **absent** — the same defect
  in kind as the recipe's `null`, and the note is in the ledger with the corrected
  reason (the fingerprint hashes the field **names**, so presence would not move the
  version).
- The **unbracketed** mouse position's message does not name its remedy.

## The mutation checks

**Mutate, run, revert**, and report what you mutated and what the suite said:

| # | Mutation | What must fail |
|--:|---|---|
| 1 | the footer deciding by the item's **kind** | the footer-asks test |
| 2 | the console **calling the operation** instead of asking the domain | the refusal-writes-nothing test |
| 3 | `hasDetail()` without the **Log** case | the Log's detail test |
| 4 | the footer offering `enter detail` **unconditionally** | the footer's promise test |

**Revert every mutation**, and **confirm the revert** — `git status --short` and
`git diff` at the end. The tree must be exactly as the runner left it, except your
ledger commit.

## The ledger — `docs/audit.md` is yours

- **F-114, F-115 and F-116** get their rows: closed, with the commit that closed
  each and the test that holds it.
- The **slice T3a** row and the **cycle 1** row are there. Add the **cycle 2** row.
- Anything you file is **F-117** onward.
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
cat > work/inbox/reviewer-tui-t3a-cycle2.md <<EOF
Reviewer: T3a cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-117.. or none>). F-114: <the mutation you took, and the test that failed>. F-115: <the checksums before and after, and where the reason lands>. F-116: <what the Log's enter opens, and the footer's promise tab by tab>. The shape: <that the slice's six criteria still hold>. The scope: <what you checked>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t3a cycle 2"
herdr agent prompt "w1H:p1" "Reviewer: T3a cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-117.. or none>). F-114: <the mutation you took, and the test that failed>. F-115: <the checksums before and after, and where the reason lands>. F-116: <what the Log's enter opens, and the footer's promise tab by tab>. The shape: <that the slice's six criteria still hold>. The scope: <what you checked>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
