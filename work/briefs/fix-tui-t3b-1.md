You are the **runner**, again — the same lane that wrote T3b. The orchestrator's pane id
is `w1H:p1`; use it verbatim in the notify command below.

This is **T3b's fix cycle 1**: **one blocker in the orchestrator's own file** (already
fixed — see below), **one blocker of yours**, and **one note the orchestrator ruled
in**.

## The review's verdict

`changes-needed`, **one blocker: F-117**. The slice's shape held — the domain answers
the move and the console asks it, shaping changes the kind and adds a scope without
touching the id, dismissal records the reason, create writes a backlog item and
nothing else, and the footer names the move. All four of the review's mutations
failed the suite.

## F-117 (medium) — the plan promised a slice the brief ruled out

**This one is the orchestrator's, and it is already fixed.** `work/tui/plan.md` — the
plan is the **orchestrator's** file, and the orchestrator had corrected two of its
three places and missed the third: the opening table still said *"dismiss · create ·
edit reach the console"*. The orchestrator has corrected it, and the table's
*Promote* row, which said the kind change is `idea → feature` — narrower than what
the console does, where the chef picks the kind.

**Nothing for you to do on that file** — the orchestrator's commit says which. **But
check the correction** when you review, and check that **no other doc** promises the
free edit: the house's rule is that a doc which promises what is not built is a bug.

## F-118 (medium) — `docs/05` still promises the edit

**The target:** the console's own doc says what the console does.

`docs/05-human-interface.md` opens with the console's **jobs**, and one of them is:

```
| **Manage the backlog** | Create, edit, shape, promote, dismiss backlog items |
```

**The console does not edit an item's fields.** It shapes, promotes, dismisses and
creates — and the free edit is **not built** (a later round, if it is wanted at all).
So the doc promises a job the console does not have — the same class as F-117, in the
file **you** own.

**What lands:** the job list says what the console does. **Say in your report** what
you wrote.

## F-119 (low, ruled in by the orchestrator) — the prompt hardcodes the domain's kinds

**The target:** the kinds the shape prompt offers come from the **domain**, not from a
list the console carries.

The shape prompt offers *f) feature · b) bug* — hardcoded in the console. The domain
owns the kinds (`core.BacklogKinds`) and owns the rule that says which of them are
**small tasks**. So the console carries a copy of a domain fact, and if the domain
gained a small-task kind, the console's prompt would be stale — the same shape as
**F-114**, which the review closed two cycles ago: *the console decides for itself,
and nothing holds it.*

**What lands:** the kinds the prompt offers come from the domain. If the domain's move
answer carries them, that is the answer. If it needs a new domain answer, **say so
rather than reaching for one** — and the wire's rules apply to anything new and
exported (the **`// view`** marker, a registration in `core.DerivedViews`, and a row
in `docs/03` § 8 — or a declared helper in `core.NonViews`, named one by one).

**The evidence:** the mutation must fail the suite — a console-side kind list that
disagrees with the domain's.

## The notes the review left, and the orchestrator's ruling

**Two stay for a later round, and here is why:**

- **After a write the cursor keeps its index, and the re-grouping moves the item.**
  This is **pre-existing and console-wide**: a **sign-off** on the Recipes tab moves
  the recipe between the triage sections the same way, and the console has always
  behaved this way. It is a console-wide question about the selection, not T3b's.
  It goes in the ledger for a later round.
- **The shaping test's id assertion is near-tautological.** The id does not change
  because the domain never changes an item's id, so the assertion is a guard rather
  than a proof. It stays as it is — the guard is worth having.

**One is a note, and it is right as it stands:** the moves act while the detail pane
is open. That is the console's model — the moves act on the selection, and the detail
shows the selection — and the recipe path does the same.

## The scope — what not to touch

| | |
|---|---|
| **The moves** | **closed.** Shaping, dismissing and creating work, and the domain answers the move |
| **The promotion** | unchanged: a small task promotes with its one criterion; an idea does not, and the CLI's refusal stands |
| **The recipe** | the console does not author one, and this cycle creates none |
| **`internal/core`** | only if F-119's answer needs it — and then **say so** rather than reaching for it |
| **`internal/cli`** | unchanged |
| **The tabs** | seven, the same names |
| **`core.DomainVersion`** | stays **7** |
| **The goldens** | expected unchanged — they hold the **Recipes** tab |
| **`work/tui/plan.md`** | the **orchestrator's** file — the orchestrator corrected it. **Do not touch it** |
| **`docs/audit.md`** | the ledger is the **reviewer's**. This cycle files no finding and touches no row |
| **`work/`** | the orchestrator's, except the state row |

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. **No dependency is added.**
- **The domain layer has zero UI imports.**
- **The tab count is 7** and the tab names do not change.
- **One fact, one answer** — a fact the console shows is computed in the domain.
- The gate is `make check`. A focused check is `make test PKG=… RUN=…`.

## The evidence this cycle must produce

| Finding | The evidence |
|---|---|
| **F-118** | `docs/05`'s job list says what the console does — and **no doc in the tree** promises the free edit. Say which files you checked |
| **F-119** | the kinds the prompt offers come from the domain — and the mutation (a console-side kind list that disagrees with the domain's) turns the suite **red**. Say where the kinds come from now |

Plus: `make check` green, `make check-race` green, the goldens unchanged, and the tree
clean.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t3b-fix1.md <<EOF
Runner: T3b fix 1 done. Commits: <n>. Fixed: F-118, F-119. The doc: <what you wrote, and every file you checked for the free edit>. The kinds: <where they come from now, and the mutation that now fails it>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t3b fix 1"
herdr agent prompt "w1H:p1" "Runner: T3b fix 1 done. Commits: <n>. Fixed: F-118, F-119. The doc: <what you wrote, and every file you checked for the free edit>. The kinds: <where they come from now, and the mutation that now fails it>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
```

**Write the receipt, commit it, and *then* nudge** — the nudge only reaches an idle
receiver; the receipt is the half that cannot be lost.

**One write per file, and commit as soon as the file is coherent.**

**And if you find a file you did not write has changed: stop and report it**, and do
not touch it.

## Before you report

```bash
make check && git status --short
```

Green, and the tree clean. Report the **real** numbers: the commits, the gate's
result, and where the kinds come from.
