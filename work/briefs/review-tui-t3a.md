You are the **reviewer** for **T3a** — the backlog is readable. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` first, then
[work/tui/plan.md](../tui/plan.md) § 4 for the reasoning, and
[work/tui/README.md](../tui/README.md) for the state and the protocol.

Read [work/briefs/tui-t3a.md](tui-t3a.md) — the slice's contract.

**You are a fresh lane.** You did not write this code and you did not read it while
it was written. Re-derive the facts from the code — not from the runner's report, and
not from the brief's claims.

## The chef's ruling — the mouse, and what it opened

> *some tabs we can use the scroll, but the click does not work. also, we cant select
> the tabs using the mouse.*

The mouse works now. What it exposed is the **Backlog**: the cursor moves over an item
you cannot read, and one keystroke takes it somewhere it cannot come back from.

## The slice's six acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | a backlog item has a **detail pane** — title, body, scope, kind, status, tags, dates, and `promoted_to` when it has one. And **a click** on a backlog row opens it |
| **AC-02** | the console offers **no move the domain refuses**: `space` on an idea is refused **with the domain's reason**, and the footer does not offer a move the domain refuses |
| **AC-03** | **D-3** — the domain refuses to promote an **idea**, with the reason; a **small task** still promotes, with its one criterion |
| **AC-04** | **D-5** — every writer writes the same shape: a list is `[]`, never `null` |
| **AC-05** | the goldens are expected to be **unchanged** |
| **AC-06** | `work/tui/README.md`'s state row is updated and committed with the slice |

## What to check — and how

**Each check below is a command, not a reading.** A check you cannot run, you do not
report as passed.

### 1. The item's detail, and the click that opens it

- `enter` on a backlog item shows **its** detail — the id and title, the kind and
  status, the dates, the body, the scope, the tags, and `promoted_to` when it has one.
- It is the console's **one** detail pane: the same geometry, the same scroll, the
  same `esc`.
- **A click on a backlog row opens it**, exactly as a click on a recipe row opens its
  recipe. T2b made the backlog's click select-and-stop, because the Backlog had no
  detail. Now it has one.

### 2. The refusal is the domain's

- `space` on an **idea** is refused **with the domain's reason** — not with a
  console-side wording. The console **asks**; it does not decide by itself.
- The **footer** does not offer `space promote to recipe` on an idea, and still does
  on a small task. The console's own rule is in its code: *a footer that promises a
  triage the space refuses is a lie.*
- `maitre-cli backlog promote BL01` on an idea is refused with the same reason, and
  **exits 1** like every other refusal. No verb changed.

### 3. D-3 — the rule is in one place

- The promotion and the console read **one** answer, and the reason names the remedy:
  *shape it into a small task first*.
- A **small task** still promotes, and still gets **one** criterion.

### 4. D-5 — every writer writes the same shape

```bash
maitre-cli recipe create --title "A plain recipe"     # "acceptance": []   "scope": {"in_scope": [], ...}
maitre-cli backlog promote BL02                       # a small task       "scope": {"in_scope": [], ...}
```

- The rule lives in the **writer** — the store normalises before it writes — so a
  recipe written **first** time says what a recipe read and written back says.
- **The digest does not move.** The approval's payload canonicalises its lists, so
  `null` and `[]` hash the same. Check the repo's own plan still verifies, and that a
  fresh approve → write → read → recompute still verifies.

### 5. The wire — the new view

`core.PromotionOf` is a **new derived view**: the `// view` marker, registered in
`core.DerivedViews` as `promotion`, and given a row in `docs/03` § 8. The wire checks
**both** directions, and the table's name comes from the runtime.

- Check the registration, the row, and that the wire can **fail**.
- Check that the console really asks the domain, and does not carry its own answer.

### 6. The scope held

| | |
|---|---|
| **The backlog's other moves** | **T3b** — dismiss, create, edit. This slice made the item readable and removed the trap |
| **The domain's writes** | `core.Init` still writes what it wrote. The change is in the **reader** and in the **shape the writer writes** |
| **`maitre-cli`** | no verb added, no command changed. The **behaviour** of `backlog promote` on an idea changed — that is AC-03 |
| **The tabs** | seven, the same names |
| **`core.DomainVersion`** | still **7** — the facts the format already had do not become new facts, and the version's instrument says so |
| **The budget** | **0 / 0 / 28** against **0 / 0 / 30** |
| **The goldens** | `make goldens` leaves `git diff` empty |
| **`docs/03` and `docs/07`** | the slice corrected two places that said an idea could be promoted straight to a recipe. **Check those corrections are right** |
| **`docs/audit.md`** | the ledger is **yours** |

## Two questions the orchestrator did not rule on

**These are questions, not findings.** Answer them in your report; file them if you
think they are findings.

### A. Does a refused `space` write to the plan?

The domain journals a refusal (`refuseLocked`), so the console's `space` on an idea
**leaves a refusal event in the plan's journal** — and the footer does not even offer
the move. A misclick writes to the record.

Is that right? The house journals refusals on purpose (a refused command is evidence
it was tried), and the journal is append-only. But the console *knows* there is no move
— it read the answer before the chef pressed anything. Say what you think, and why.

### B. The backlog item's own lists are written **absent**

The recipe's lists are `[]`. The **backlog item's** `scope` and `tags` are written
**absent** — the field is `omitempty` — so a reader that asks for a missing key gets
`null`, which is the hazard the domain's own rule names.

The runner ruled this out of scope: *making them present would be a format change and
a version bump*. **Check that claim** — the fingerprint hashes the field **names**, and
`BacklogItem.scope` is already one of them — and say whether the item's absent list is
the same defect as the recipe's `null`, or a different thing.

## The mutation checks — this is what makes a review worth something

A green suite proves the code passes the tests it was written with. It does not prove
the tests *can* fail. **Mutate, run, revert** — and report what you mutated and what
the suite said:

| # | Mutation | What must fail |
|--:|---|---|
| 1 | make the promotion of an idea **succeed** | the refusal test |
| 2 | make the console **decide** promotability itself, instead of asking | the footer's test |
| 3 | take the writer's normalisation out | the shape test |
| 4 | make `PromotionOf` answer for **every** kind | the refusal test |
| 5 | remove the view's registration, or its row in `docs/03` § 8 | the wire |

**Revert every mutation**, and **confirm the revert** — `git status --short` and
`git diff` at the end. The tree must be exactly as the runner left it, except your
ledger commit.

## The ledger — `docs/audit.md` is yours

- File **HIGH or MEDIUM** findings. Number them from **F-110** — T2b's last was
  **F-109**.
- **LOW findings: note them in your message, do not file them.**
- Add the **slice T3a** row to the ledger's slice table.
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
cat > work/inbox/reviewer-tui-t3a.md <<EOF
Reviewer: T3a review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-110.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The detail: <what it shows, and the click that opens it>. The refusal: <the reason an idea is refused with, and where the reason comes from>. The footer: <what it says on an idea and on a small task>. The shape: <what the promotion writes, and what recipe create writes>. The digest: <that the repo's plan and a fresh approve still verify>. The view: <its registration, its row, and that the wire can fail>. A: <your judgement on the refusal event>. B: <your judgement on the item's absent list>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t3a"
herdr agent prompt "w1H:p1" "Reviewer: T3a review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-110.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The detail: <what it shows, and the click that opens it>. The refusal: <the reason an idea is refused with, and where the reason comes from>. The footer: <what it says on an idea and on a small task>. The shape: <what the promotion writes, and what recipe create writes>. The digest: <that the repo's plan and a fresh approve still verify>. The view: <its registration, its row, and that the wire can fail>. A: <your judgement on the refusal event>. B: <your judgement on the item's absent list>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
