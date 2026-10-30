You are the **runner**, again — the same lane that wrote T3a. The orchestrator's pane id
is `w1H:p1`; use it verbatim in the notify command below.

This is **T3a's fix cycle 1**: **two blockers**. Both are about the console's side of
the rule the domain now owns.

**The runner's tree report is in the record** (`work/inbox/runner-tui-t3a-tree.md`):
the reviewer's mutation was in the tree while the gate ran, and the runner stopped,
reported it, and did not touch the file. That was right, and it is now in the round's
protocol.

## The review's verdict

`changes-needed`, **two blockers: F-114 and F-115**. The slice's shape held: the
item's detail, the domain's refusal, the footer's truth, the shape every writer
writes, the digest, the view's registration and its row — all proved by command, and
four of the review's five mutations failed the suite.

## F-114 (medium) — the footer's guard must bite

**The target:** the console **asks** the domain whether the item may be promoted, and
a test **holds** that it asks.

**The hole.** `internal/ui/views.go`'s `primaryActionLabel` asks `core.PromotionOf` —
and **nothing holds it**. The reviewer's mutation — the footer deciding by the item's
**kind** instead of asking the domain — leaves the **whole `internal/ui` suite
green**.

Why: the two cases the test holds are an **idea** and an **small task**, and those
are exactly the two where the console's own rule and the domain's answer **agree**.
The case where they **differ** is untested: a **promoted** small task — the domain
refuses it, because it is no longer open, while the console's own rule ("is it a
small task?") would say yes. The mutated console then offers `space promote to
recipe` on it. The reviewer proved that by command.

**What lands:** a test on the case where the two answers **differ** — the promoted
small task — so the console's own rule cannot stand in for the domain's. And the
mutation must turn the suite **red**.

**Say in your report** which case your test uses and why it is the one that bites.

## F-115 (medium) — the console must not journal a move it never offered

**The target:** pressing `space` where the footer offers **no** move writes **nothing**
to the plan.

**The hole.** The console asks the domain by **calling the operation**
(`PromoteBacklogItem`), and the operation journals every refusal (`refuseLocked`) —
so pressing `space` on an idea appends a **`refusal` event** to `activity.json`, the
plan's only event log, **append-only, with no way to take it back**. And the footer,
**in the same frame**, offers no move at all: nobody attempted anything, and the
journal now says the chef tried to promote an idea.

The journal's own rule: a refusal is *"an attempt that was turned away"*.

**What lands:** the console asks **`core.PromotionOf`** — a **read** — and shows the
domain's reason, **writing nothing**. The console's own precedent is the recipe path,
which asks `core.ChefMoveOf` and does **not** call the operation when there is no
move. **The CLI's refusal stays**: an agent asked, and the plan said no.

**The evidence:** `space` on an idea leaves `activity.json` **byte-for-byte
unchanged**, and the reason still reaches the chef — in the banner, or in the item's
detail, or wherever the console puts it. **Say where the reason lands now.**

## What lands

| | |
|---|---|
| **`internal/ui`** | the footer asks the domain and a test holds it (F-114); the refusal writes **nothing** (F-115) |
| **`internal/core`** | nothing, unless the evidence above forces it — and if it does, **say so** rather than reaching for it |
| **`internal/cli`** | unchanged — the CLI's refusal is the CLI's |
| **The domain's writes** | unchanged: `core.Init` writes what it wrote, and every recipe is written in the shape the reader reads |
| **`core.DomainVersion`** | stays **7**. The facts the format already had do not become new facts, and the version's instrument says so |
| **The tabs** | seven, the same names |
| **The goldens** | expected unchanged |
| **`docs/05`, `docs/03`, `docs/07`** | unchanged unless a line became false — and nothing here changes the console's *behaviour* except what the refusal writes |
| **`docs/audit.md`** | the ledger is the **reviewer's**. This cycle files no finding and touches no row |
| **`work/`** | the orchestrator's, except the state row |

## The notes the review left, and the orchestrator's ruling

**Do not fix the backlog item's own lists in this cycle.** The review ruled: the
item's absent `scope`/`tags` are the **same defect in kind** as the recipe's `null`
— a reader that asks for the list gets nothing — but the **runner's reason** for
leaving it was wrong: the fingerprint hashes the field **names**, so writing `[]`
would not move the version. The note is in the ledger with the corrected reason, and
it waits for a later round. **This cycle does not widen.**

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
| **F-114** | the mutation — the console deciding promotability for itself — turns the suite **red**, and the test that fails is the one on the case where the two answers **differ**. Say which case it is, and why it bites |
| **F-115** | `space` on an idea leaves `activity.json` **byte-for-byte unchanged** — show the checksum before and after — and the domain's reason still reaches the chef. Say where it lands |

Plus: `make check` green, `make check-race` green, the goldens unchanged, and the tree
clean.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t3a-fix1.md <<EOF
Runner: T3a fix 1 done. Commits: <n>. Fixed: F-114, F-115. The footer's guard: <the case your test uses, and the mutation that now fails it>. The refusal writes nothing: <the checksums before and after, and where the reason lands now>. The CLI's refusal: <that it still refuses, with the reason>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t3a fix 1"
herdr agent prompt "w1H:p1" "Runner: T3a fix 1 done. Commits: <n>. Fixed: F-114, F-115. The footer's guard: <the case your test uses, and the mutation that now fails it>. The refusal writes nothing: <the checksums before and after, and where the reason lands now>. The CLI's refusal: <that it still refuses, with the reason>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
```

**Write the receipt, commit it, and *then* nudge** — the nudge only reaches an idle
receiver; the receipt is the half that cannot be lost.

**One write per file, and commit as soon as the file is coherent.** Every edit to a
large file re-renders that file's whole diff into the terminal, and a turn that renders
dozens of diffs runs out before the work is done.

**And if you find a file you did not write has changed: stop and report it**, and do
not touch it. It is most likely the reviewer's mutation, in flight — and a second
revert would lose its experiment. Say whose you think it is.

## Before you report

```bash
make check && git status --short
```

Green, and the tree clean. Report the **real** numbers: the commits, the gate's
result, and the checksums that show the refusal wrote nothing.
