You are the **runner** for **T3b** — the backlog's moves. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` first, then
[work/tui/plan.md](../tui/plan.md) § 4 for the reasoning, and
[work/tui/README.md](../tui/README.md) for the state and the protocol.

Read [work/briefs/tui-t3a.md](tui-t3a.md) — the slice before this one — because
**T3a built the ground this slice stands on**: the item's detail, and the domain's
answer to *"may this be promoted?"* (`core.PromotionOf`), which the console asks.

## The chef's ruling — the backlog's moves

> *Q6 — may the console create backlog items?* **yes** — a backlog item is not a
> recipe, so ADR-0003 does not forbid it.

T3a's refusal ends with the remedy: *"shape it into a small task first"*. **This slice
is that remedy, and the two other moves the tab wants.**

Every operation this slice needs **already exists** in the domain:
`CreateBacklogItem`, `UpdateBacklogItem`, `DismissBacklogItem`. This slice is the
console catching up to its own domain, and it needs no CLI work.

## Your slice — T3b — the backlog's moves

### AC-01 — the domain answers what move the item has

**The target:** the console **asks** the domain what the chef can do with the selected
item, and the footer and the key read **that** answer.

T3a's `core.PromotionOf` answers *"may this be promoted, and why not"*. This slice
needs the **move** as well — an open **idea** has a move that is not a promotion. So:

- The domain answers the move, in **one** place, and the console asks it.
- **The console does not decide.** A console-side "is it a small task?" rule is the
  defect F-114 just closed; do not reopen it.
- If the answer is a **new exported function**, the wire applies to it, both ways: the
  **`// view`** marker, a registration in `core.DerivedViews`, and a row in
  `docs/03` § 8 — or a declared helper in `core.NonViews`, named one by one. **Say
  which you chose and why.** The table's name comes from the runtime.

**The three moves**, and what each one means:

| The move | On | What it does |
|---|---|---|
| **shape it** | an open **idea** | the kind changes to a small task, **and a scope is added** — `docs/03` § 1.1: *"An **update**, not a promotion: the kind changes and a scope is added, and the id does not change"* |
| **promote to recipe** | an open **small task** | the promotion T3a kept — one criterion, signable |
| **none** | a promoted, dismissed, or otherwise settled item | no move, and the footer says so |

### AC-02 — shaping an idea

**The target:** the chef can shape an idea into a small task, and the item keeps its
id.

- **The kind and the scope** — a scope is *added*, so the chef supplies one. The
  console's prompt asks for what the domain needs and no more.
- **The id does not change**: `BL01` stays `BL01` — the id is kind-neutral, so the
  trail survives the shaping. The domain's `UpdateBacklogItem` does it; the console
  asks.
- **This is an update, not a promotion.** No recipe is created, and the item is not
  promoted.

### AC-03 — dismissing an item

**The target:** the chef can dismiss an item, **with a reason**.

- The domain appends the reason to the item's body (`DismissBacklogItem`), so a
  dismissal **without** a reason says nothing to the next reader. **Ask for it.**
- The reason is the chef's own words, in the console's existing text prompt.

### AC-04 — creating an item

**The target:** the chef can write an item down from the console.

- **The chef's ruling: yes.** A backlog item is **not** a recipe, so ADR-0003 — *the
  console never authors the recipe* — is not broken. The console creates a **backlog
  item** and nothing else.
- The item starts as an **idea** — the kind a hunch starts as — and the chef shapes
  it after, with AC-01's move. Say in your report what the console asks for.
- **No recipe is created by this slice**, and no recipe is authored by the console.

### AC-05 — the footer tells the truth, again

- The footer names the move the selected item has — **`space shape it`** on an open
  idea, **`space promote to recipe`** on an open small task, and **no `space`** on an
  item with no move.
- T3a closed the *refusal* on `space` for an idea. This slice gives an idea a **move
  instead** — so `space` on an idea now **shapes** it, and there is nothing to refuse.
  The domain's refusal of the *promotion* of an idea stands: `maitre-cli backlog
  promote` on an idea is still refused, with the reason, and the journal still records
  that an agent asked.

### AC-06 — the goldens

**The goldens are expected to be unchanged.** They hold the **Recipes** tab, and this
slice changes the **Backlog** tab. If a golden moved, find out why and say so.

### AC-07 — the state row is updated

`work/tui/README.md`: T3b goes to `closed` with its commits, and the next session
becomes the first row that is not `closed`. Commit it with the slice.

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. **No dependency is added.**
- **`core.DomainVersion` is 7.** No format, no field, no directory.
- The tab count is 7 and the tab names do not change.
- **One fact, one answer** — a fact the console shows is computed in the domain.
- **The domain layer has zero UI imports.**
- **The console never authors a recipe.** It authors a *backlog item* — which is not
  one. If a change would cross that line, **stop and report it**.
- The gate is `make check`. A focused check is `make test PKG=… RUN=…`.

## The scope — what not to touch

| | |
|---|---|
| **The free edit of every field** | **not this slice.** The three moves are what the chef needs and what the refusal points at. Editing an item's title, body and tags as a form is a later slice, if it is wanted at all |
| **The promotion** | T3a's, unchanged: a small task promotes, with its one criterion; an idea does not |
| **The recipe** | the console does not author one, and this slice creates none |
| **`maitre-cli`** | no verb is added and no command changes |
| **The tabs** | no tab is added, removed or renamed |
| **The domain's writes** | `core.Init` still writes what it wrote. The change is in the **reader** and in the **answer the console asks for** |
| **The mouse** | T2's rules are closed. A click on a backlog row opens the item's detail — unchanged |
| **The docs** | `docs/05-human-interface.md` describes the screen and the keys: if this slice makes a line of it false, fix the line and say so. And a **new answer** in the domain means a row in `docs/03` § 8 — see AC-01 |
| **`docs/audit.md`** | the ledger. **Never swept.** This slice files no finding |
| **`work/`** | the orchestrator's, except the state row (AC-07) |

## The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. If a file you did not
write has changed, **stop and report it**, and do not touch it: it is most likely the
reviewer's mutation, in flight.

## The evidence this cycle must produce

For **each** criterion, a command and its output — not a claim:

| Criterion | The evidence |
|---|---|
| **AC-01** | the footer's move on each state of an item — an open idea, an open small task, a promoted item — and that the console **asks** the domain for it. **Mutate it**: make the console decide for itself, and the suite must go **red** |
| **AC-02** | shaping an idea changes its kind, **adds a scope**, and **keeps its id** — `BL01` is still `BL01` — and creates **no** recipe |
| **AC-03** | dismissing an item records the chef's **reason**, and the item's body carries it |
| **AC-04** | creating an item from the console writes a **backlog item** — and **no** recipe, and nothing outside the plan's backlog |
| **AC-05** | `space` on an idea **shapes** it — nothing is refused — and the CLI's promotion of an idea is **still** refused, with the reason |
| **AC-06** | `make goldens` leaves `git diff` empty |
| **AC-07** | the state row names the slice's commits, and the next row is the next session |

Plus: `make check` green, `make check-race` green, and the tree clean.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t2b.md <<EOF
Runner: T3b done. Commits: <n>. Per-AC evidence: AC-01..AC-07 stamped. Pre-notify gate: pass. The moves: <what the domain answers, and where the answer lives>. Shape it: <what changes, and that the id does not>. Dismiss: <that the reason is recorded>. Create: <what the console asks for, and that no recipe is made>. The footer: <what it says on an idea, on a small task, and on a settled item>. The domain's answer: <the function you added, and whether it is a view with a row in docs/03 or a declared helper>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <what the next session should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t3b"
herdr agent prompt "w1H:p1" "Runner: T3b done. Commits: <n>. Per-AC evidence: AC-01..AC-07 stamped. Pre-notify gate: pass. The moves: <what the domain answers, and where the answer lives>. Shape it: <what changes, and that the id does not>. Dismiss: <that the reason is recorded>. Create: <what the console asks for, and that no recipe is made>. The footer: <what it says on an idea, on a small task, and on a settled item>. The domain's answer: <the function you added, and whether it is a view with a row in docs/03 or a declared helper>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <what the next session should know>."
```

**Write the receipt, commit it, and *then* nudge** — the nudge only reaches an idle
receiver; the receipt is the half that cannot be lost.

**One write per file, and commit as soon as the file is coherent.** Every edit to a
large file re-renders that file's whole diff into the terminal, and a turn that renders
dozens of diffs runs out before the work is done.

## Before you report

```bash
make check && git status --short
```

Green, and the tree clean. Report the **real** numbers: the commits, the gate's
result, and what each move does.
