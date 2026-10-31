You are the **runner**, again — the same lane that wrote T3b. The orchestrator's pane id
is `w1H:p1`; use it verbatim in the notify command below.

This is **T3b's fix cycle 2**: **one blocker**, and **one note the orchestrator ruled
in**.

## The review's verdict

`changes-needed`, **one blocker: F-120**. The slice's shape held, and F-117, F-118 and
F-119 are closed — the reviewer confirmed each by mutation.

## F-120 (medium) — `docs/05`'s job list speaks in two voices

**The target:** the job list says **one** thing.

**The hole.** Cycle 1's review noted, as a LOW note, that `docs/05`'s job list
promised the console's free edit — and the orchestrator ruled it into a fix cycle.
**That ruling was wrong**, and the reviewer has caught it: the fix made the backlog's
row speak in **build state** (*"the free edit of an item's fields is **not built**"*),
while the other three rows speak in **remit** — what the console's job *is*:

| Row | What it names, and whether the console does it |
|---|---|
| **View** | *"findings anchored to code"* — the domain's `Finding.Anchor` is **never read** by `internal/ui` |
| **Decide** | *"accept or reject a scope"* and *"choose what is next"* — the console does neither |
| **Run** | *"start"* and *"pause"* — the console does neither; the agent runs the work |
| **Manage the backlog** | now says what is **built** — while the three above say what is **remit** |

So the table promises three things the console does not do, and says one thing it does
— in two voices.

**My ruling, and it is the orchestrator's — the table is the REMIT.** The console's
job list says **what the console is for**, and the console's **surface** is described
by the keys table below it. That is how ADR-0003 puts it — *the human interface is a
**management surface*** — and the orchestrator's ruling on F-117 (that the ADR's own job
table **stays**, in the vocabulary of its day) says the same thing about the same
table: it is a statement of the **division of labour**, not of the build.

**What lands:**

1. **One line above the table** says what the table is — the console's remit — and
   that **what is built is the keys table below**.
2. **The backlog's row returns to the remit's wording.** The clause *"the free edit of
   an item's fields is not built"* comes **out** of `docs/05` — it is a **build state**,
   and this table is not where build states live. **It stays in the plan**, which says
   the free edit is *not in this round*.

**Say in your report** what you wrote, and that the three other rows are now in the
same voice.

**One thing to be careful about.** The table is the console's **remit**, so it may name
things the console does not do **yet** — that is the point of a remit. But it must not
name things **no console will ever do**, and it must not name **another agent's** work
as the console's own: *"start, pause"* is the **agent's** job, and the chef's console
drives it. If a row reads that way, **say so in your report** — the orchestrator
ruled the table's voice, and a row that crosses the line is a question the orchestrator
answers, not one you settle.

## F-121 (low, ruled in by the orchestrator) — the prompt's keys must not collide

**The target:** every kind the prompt offers is reachable by its own key.

**The hole.** `kindKey` takes the kind's **first byte** — so two kinds sharing an
initial share the prompt's key, and one of them becomes **unreachable**. Today the
kinds are `feature` and `bug`, so `f` and `b` — but the domain's kinds are a `var`,
and a future kind whose name starts with the same letter would silently disappear
from the prompt.

The **letter** that stands for a kind is the console's own affordance — the runner
said so, and it is right: the domain owns the kind, the console owns the key. So the
console's own rule must be that the letters it hands out do not collide.

**What lands:** the prompt's choices carry **distinct** keys, and the test says so.

**The evidence:** the mutation must fail the suite — two kinds sharing an initial must
not produce two choices with the same key.

## The notes the review left, and the orchestrator's ruling

**One stays for a later round, and here is why:** `BacklogMoveOf` builds the kinds
list on **every call**, and the footer asks it on every frame. That is a per-frame
allocation — a performance question about the console's render path, and the console's
render path has its own shape (it redraws the whole screen). It goes in the ledger for
a later round; it is not this slice's.

## The scope — what not to touch

| | |
|---|---|
| **The moves** | **closed.** Shaping, dismissing, creating, and the domain's answer for the move |
| **`internal/core`** | untouched by this cycle, unless F-121's distinct letters need it — and then **say so** rather than reaching for it |
| **`internal/cli`** | unchanged |
| **The tabs** | seven, the same names |
| **`core.DomainVersion`** | stays **7** |
| **The goldens** | expected unchanged |
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
| **F-120** | `docs/05`'s job list says once what the table is, and **every row is in that voice** — the backlog's row no longer carries a build state. Say what you wrote and how each row now reads |
| **F-121** | the prompt's choices carry **distinct** keys, and the mutation — two kinds sharing an initial — turns the suite **red** |

Plus: `make check` green, `make check-race` green, the goldens unchanged, and the tree
clean.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t3b-fix2.md <<EOF
Runner: T3b fix 2 done. Commits: <n>. Fixed: F-120, F-121. The job list: <the line you added, and how each row now reads>. The keys: <how the prompt's choices keep distinct keys, and the mutation that now fails it>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t3b fix 2"
herdr agent prompt "w1H:p1" "Runner: T3b fix 2 done. Commits: <n>. Fixed: F-120, F-121. The job list: <the line you added, and how each row now reads>. The keys: <how the prompt's choices keep distinct keys, and the mutation that now fails it>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
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
result, and how the job list reads now.
