You are the **runner** for **T3a** — the backlog is readable. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` first, then
[work/tui/plan.md](../tui/plan.md) § 4 for the reasoning, and
[work/tui/README.md](../tui/README.md) for the state and the protocol.

**This is the first slice of the second half of the round.** T1 gave the console its
identity; T2 gave it the mouse. **T3a makes the backlog readable** — and removes the
trap the chef can walk into with one keystroke.

## The chef's ruling — the mouse, and what it opened

> *some tabs we can use the scroll, but the click does not work. also, we cant select
> the tabs using the mouse.*

The mouse now works. What it exposed is the **Backlog** tab: you can move the cursor
over an item you cannot read, and one keystroke takes it somewhere it cannot come back
from.

## What the plan asked for, and what the code says today

| | Today | This slice |
|---|---|---|
| **Read** | `enter` on a backlog item does **nothing** — the detail pane is empty. The item's body, its scope, its tags and its status are unreachable | a detail pane for the item |
| **Promote** | `space` on **every** row — including an **idea** — and the promotion is not gated by kind | `space` on a **small task** promotes. On an **idea** it is **refused, with the reason** |

**The trap, proved by command:**

```bash
maitre-cli backlog create --title "The board needs a way to show evidence" --kind idea
maitre-cli backlog promote BL01
# "acceptance": []          <- no criteria
```

A recipe with no acceptance criteria can never pass `ac.falsifiable`, so it can never
be signed off. **The console offers, in one keystroke, a move that leads to a dead
end.** The domain already knows the difference — `isSmallTask` says *"the work is
known, and it is ready to run"* — and it gives a small task one criterion when it
promotes. An idea gets none.

## Your slice — T3a — the backlog is readable

### AC-01 — a backlog item has a detail pane

`enter` on a backlog item opens its detail: the **title**, the **body**, the **scope**,
the **kind**, the **status**, the **tags**, the **dates**, and **`promoted_to`** when it
has one.

- The item is **not** a recipe, and the console's detail machinery is recipe-shaped
  today — so the item needs its own rendering. **You find the seam.**
- It is the console's **one** detail pane: the same geometry, the same scroll, the
  same `esc`. Not a second pane.
- The fields are **stored facts** the domain reads; the console renders them.
- **The click opens it too.** T2b made a click on a row select it and stop, because
  the Backlog had no detail. Now it has one: **a click on a backlog row opens the
  item's detail**, exactly as a click on a recipe row opens its recipe's.
- **The goldens do not hold the Backlog tab** — so the change is invisible to them.

### AC-02 — the console offers no move the domain refuses

**The target:** `space` on an **idea** is refused **with the reason**, and the footer
does not promise a move the domain refuses.

**The refusal is the domain's.** The console asks the domain to promote the item and
shows **the domain's reason**. It does not decide by itself that an idea is not
promotable — that would be a second answer to one question, and the domain already
answers it (AC-03).

**The footer tells the truth.** Today it says `space promote to recipe` on **every**
backlog row. On an idea there is no such move, so the console must not say it. The
console's own rule is in its code: *a footer that promises a triage the space refuses
is a lie*.

**The domain answers "may this be promoted?" in one place**, and both the console's
footer and the promotion's refusal read that one answer.

> **My ruling, and it is a change to the plan.** The plan proposed that `space` on an
> idea *shape it* — the kind change `idea → feature`, in one keystroke. **This slice
> refuses instead.** A kind change is an **edit**, and the plan's own table puts
> *Edit* under the backlog's moves — **T3b** owns it. And one keystroke that turns an
> idea into a *small task* files the work as **known** — which is the domain's own
> definition of a small task — and an idea by definition does not say that. So the
> refusal **names the remedy** — *shape it into a small task* — and the way to do that
> arrives with T3b. **If you think the refusal is the wrong answer, say so in your
> report, with your reasons** — this is a judgement and the orchestrator has ruled it,
> not discovered it.

### AC-03 — D-3: an idea cannot be promoted into a dead end

**The target:** the domain refuses to promote an **idea**, and the refusal says why.

- The reason names what to do: an idea has no acceptance criteria, and a recipe
  without them can never be signed off — **shape it into a small task first**.
- **A small task still promotes**, and still gets its one criterion. That rule does not
  move.
- The rule lives in **one place**, and both the refusal and the console's answer read
  it. `isSmallTask` already says it is *"the rule lives here, once, because two
  surfaces need it"* — two surfaces now need it: the promotion and the console.

### AC-04 — D-5: every writer writes the same shape

**The target:** a list is a list — an empty list is `[]`, never `null` — **in every
file the tool writes**.

**The hole, proved by command.** Two writers, two shapes:

| The writer | The file it writes |
|---|---|
| `recipe create` | `"scope": {"in_scope": [], "out_of_scope": []}` |
| `backlog promote` | `"scope": {"in_scope": null, "out_of_scope": null}` |

The domain's own rule says: *"a list is a list: an empty list is `[]`, never `null`.
The plan is read by agents, and the standard readers — `jq`, a JSON parser in any
language — break on `null`."* And there is a third shape, from the same command:

```
recipe create --title "A plain recipe"   →  "acceptance": null
```

**Same format, two answers** — and one of the answers is `null`.

**Where the rule belongs.** The store normalises the lists when it **reads** a recipe,
so a recipe that is read and written back is normalised. A recipe written **first**
time — the promotion's, and `recipe create`'s — never goes through it.

> **My ruling: the rule belongs in the writer.** The store normalises every list
> **before it writes**, so every writer agrees — the promotion's scope becomes `[]`,
> and `acceptance: null` becomes `[]` too. One place decides, and the domain's own
> comment already says the contract holds *"everywhere the plan is read"* — this makes
> it hold everywhere the plan is **written** too.

**The digest does not move.** The approval's payload canonicalises its lists —
`canonicalList`, `canonicalCriteria` — so `null` and `[]` hash the same and **no stored
digest changes**. Say so in your report if you disagree, but check it.

### AC-05 — the goldens

**The goldens are expected to be unchanged.** They hold the **Recipes** tab, and this
slice changes the **Backlog** tab and the shape a recipe is *written* in — not what a
recipe *reads* as. If a golden moves, find out why and say so.

### AC-06 — the state row is updated

`work/tui/README.md`: T3a goes to `closed` with its commits, and **T3b** — the
backlog's moves: dismiss, create, edit — becomes the next session. Commit it with the
slice.

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. **No dependency is added.**
- **`core.DomainVersion` is 7.** The facts the format already had do not become new
  facts: *"a list is `[]`, never `null`"* is the format's own rule, and the writer now
  keeps it — so **the version does not move**. `TestDomainShapeIsVersioned` is the
  instrument; if it fails, **stop and report it**.
- The tab count is 7 and the tab names do not change.
- **One fact, one answer** — a fact the console shows is computed in the domain.
- **The domain layer has zero UI imports.**
- The gate is `make check`. A focused check is `make test PKG=… RUN=…`.

## The scope — what not to touch

| | |
|---|---|
| **The backlog's moves** | **T3b** — dismiss, create, edit. This slice makes the item readable and removes the trap; it does not add the console's surfaces for the other moves |
| **The domain's writes** | `core.Init` still writes what it wrote. The change is in the **reader** and in the **shape the writer writes** |
| **`maitre-cli`** | no verb is added and no command changes. The **behaviour** of `backlog promote` on an idea changes — that is AC-03, and it is a refusal with a reason, not a removed verb |
| **The tabs** | no tab is added, removed or renamed |
| **The mouse** | T2's. A click on a backlog row now opens the item's detail (AC-01) — that is the mouse catching up with the backlog, not a change to the mouse's rules |
| **The docs** | `docs/05-human-interface.md` describes the screen and the keys; if this slice makes a line of it false, fix the line and say so. And a **new view** means a row in `docs/03` § 8 — see below |
| **`docs/audit.md`** | the ledger. **Never swept.** This slice files no finding |
| **`work/`** | the orchestrator's, except the state row (AC-06) |

## A new view, and the wire that holds it

If the domain's answer to *"may this item be promoted?"* is a **new exported
function**, the house's wire applies to it, and it applies **both ways**:

- The function carries the **`// view`** marker, it is **registered** in
  `core.DerivedViews`, and it is **written down in `docs/03` § 8** — the wire checks
  both directions, and a view the table does not name is a view a reader cannot find.
- The table's row names the function, and the **name comes from the runtime** — a
  renamed function leaves the table naming something that is not there.

**Or** the answer is a plain method on the item, declared by name in `core.NonViews`
— the list of the exported functions that are not views, which the same wire checks
(the list cannot rot: a name there that no longer names a function fails the check).

**Say which you chose and why.** The console asks the domain; it does not decide.

## The evidence this cycle must produce

For **each** criterion, a command and its output — not a claim:

| Criterion | The evidence |
|---|---|
| **AC-01** | `enter` on a backlog item shows its detail — title, body, scope, kind, status, tags, dates, and `promoted_to` when it has one. And **a click** on a backlog row opens it |
| **AC-02** | `space` on an **idea** is refused **with the domain's reason**, and the footer does **not** offer `space promote to recipe` on an idea — while it still does on a small task |
| **AC-03** | `maitre-cli backlog promote BL01` on an **idea** is refused with the reason. And a **small task** still promotes, with its one criterion |
| **AC-04** | the promotion writes `"scope": {"in_scope": [], "out_of_scope": []}` and `"acceptance": []` — the **same shape** `recipe create` writes. And an **existing plan's approval digest** still verifies |
| **AC-05** | `make goldens` leaves `git diff` empty |
| **AC-06** | the state row names the slice's commits, and T3b is the next row |

Plus: `make check` green, `make check-race` green, and the tree clean.

## The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. If a file you did not
write has changed, **stop and report it**.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t3a.md <<EOF
Runner: T3a done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. The item's detail: <what it shows>. The refusal: <the reason an idea is refused with, and where the reason comes from>. The footer: <what it says on an idea and on a small task>. The shape: <what the promotion writes now, and what recipe create writes>. The digest: <that an existing plan's approval still verifies>. The view: <the function you added, and whether it is a view with a row in docs/03 or a declared helper>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <what T3b needs to know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t3a"
herdr agent prompt "w1H:p1" "Runner: T3a done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. The item's detail: <what it shows>. The refusal: <the reason an idea is refused with, and where the reason comes from>. The footer: <what it says on an idea and on a small task>. The shape: <what the promotion writes now, and what recipe create writes>. The digest: <that an existing plan's approval still verifies>. The view: <the function you added, and whether it is a view with a row in docs/03 or a declared helper>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <what T3b needs to know>."
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
result, and what an idea's `space` says now.
