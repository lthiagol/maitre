You are a **fresh runner**. The lane that wrote T4 is gone; its commits are in the log.
The orchestrator's pane id is `w1H:pJ`; use it verbatim in the notify command below.

This is **T4's fix cycle 1**: **three findings the orchestrator ruled in**, and **one
note**.

## The review's verdict

`ok`, **0 blockers**. The slice's shape held: the console lands on Home, the seven
names are right, the Sous-Chef's content is inside Home's automation block, every
number comes from the domain, the no-items invariant is held by a test, the art yields
to the content in both dimensions, and the goldens moved and are claimed.

## F-122 (medium) — a long plan name pushes the footer off the bottom

**The target:** a row that does not fit is **truncated**, and the footer stays on
screen.

**The hole.** The plan row truncates the **path** to 56 columns, but the **name** beside
it is not bounded — so a long plan name **wraps** the row, and the wrap **pushes the
footer off the bottom**. The review measured it: 100 characters at 80 columns.

That is a class the house has filed twice as medium — **F-48** (*a refusal that does
not fit takes rows of its own*) and **F-94** (*the art pushes the last event off the
bottom*). **The footer is the console's promise**, and the content does not push it
off the screen.

**And the Settings tab's `root` row has the same shape** — the review noted it
*tolerates* the whole path today. **One rule, both places**: a row that does not fit is
truncated to the width.

**What lands:** the plan name is bounded like the path is, and the Settings root row
gets the same treatment. **The evidence:** a plan with a long name renders **one row**
and the footer is still the last row.

## F-123 (medium) — the automation's honesty dissolved with the tab

**The target:** the console does not imply that an unbuilt thing is running.

**The hole.** The Sous-Chef tab said the automation was **not built** — that was its
honesty, and it was in the tab the plan told you to dissolve. Home's automation block
carries the vital signs and the last event, and the console **nowhere says the
automation is not built**. So the block reads like a running automation.

**What lands:** Home's automation block carries the honesty the tab carried. The
orchestrator's `AGENTS.md` already says it (*"the automation is not built: the Home
tab's automation block holds its vital signs"*) — the **screen** must say it too.

**The evidence:** the Home tab's automation block says the automation is not built, in
the console's own words.

## F-124 (medium) — the rows' parts must add up

**The target:** a row whose parts do not add up to the number beside them is a row
that **misleads**.

**The hole.** Two rows, two shapes of the same thing:

| The row | It shows | It leaves out |
|---|---|---|
| the recipes | *total 7 · need you 2 · on the line 2 · served 1* | **being shaped** (2) — the triage has four sections |
| the evidence | *criteria 12 · proven 6 · failed 2 · rejected 0 · unproven 0* | **pending** (4) — the view carries it |

The parts sum to 5 and to 8. **The orchestrator's plan is corrected** (the rows now
name the states), and the console must match it.

**What lands:** the rows name every state, so the parts **add up**. **The row must stay
on one line** — the width is a rule, and a row that wraps pushes the footer off the
bottom, which is F-122. **If a row cannot carry every state on one line, say so** —
and say what you would do instead; that is a design question the orchestrator answers.

## F-125 (medium) — `docs/05` makes a claim that is not true

**The target:** every line the console's doc says is true of the console.

**The hole.** `docs/05`'s Home paragraph says Home shows *"the counts the board does
not carry"* — and the board **does** carry the recipes' counts: the Recipes tab's
section headers are the triage's counts. So the doc claims a distinction the console
does not make.

**What lands:** the paragraph says what is true. The distinction that **does** hold is
**numbers against items** — the board shows the recipes; Home shows their counts. The
recipes' counts appearing on both is the **chef's ruling** (Q2: the full set), not an
accident, and the doc must not deny it.

## The note the review left, and the orchestrator's ruling

**The `ADR-0009` question — the ADR stays.** `docs/decisions/0009` still names the
tab *Sous-Chef*: it is a **decision of record**, in the vocabulary of its day, and the
console's tab names moved on after it — the same ruling as `ADR-0003`'s job table, and
for the same reason. `docs/07`, the vocabulary of record, **was** updated. **Nothing to
do.**

## The scope — what not to touch

| | |
|---|---|
| **The other tabs** | their content does not move |
| **The numbers' origins** | unchanged — every number still comes from the domain's view |
| **The domain's view** | `core.PlanStatsOf` keeps its registration and its row in `docs/03` § 8, and the wire still fails both ways |
| **The CLI** | no verb added, no command changed |
| **`core.DomainVersion`** | stays **7** |
| **The budget** | **0 / 0 / 28** against **0 / 0 / 30** |
| **The art** | slice G's rule stands: the decoration yields to the content, in both dimensions |
| **`work/tui/plan.md`** | the **orchestrator's** file — the orchestrator corrected the two rows. **Do not touch it** |
| **`docs/audit.md`** | the ledger is the **reviewer's**. This cycle files no finding and touches no row |
| **`work/`** | the orchestrator's, except the state row |

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. **No dependency is added.**
- **The domain layer has zero UI imports.**
- **The tab count is 7**, and the names are `Home · Recipes · Backlog · Path · Board ·
  Log · Settings`.
- **One fact, one answer.**
- The gate is `make check`. A focused check is `make test PKG=… RUN=…`.

## The evidence this cycle must produce

| Finding | The evidence |
|---|---|
| **F-122** | a plan with a **long name** renders the plan row on **one line**, and the **footer is still the last row**. Say the name you used and the width |
| **F-123** | the Home tab's automation block **says the automation is not built**, in the console's own words. Show the line |
| **F-124** | each of the two rows' parts **add up** — show the arithmetic — and each row is **one line** |
| **F-125** | `docs/05`'s Home paragraph says what is **true**, and the numbers-versus-items distinction is the one it makes |

Plus: `make check` green, `make check-race` green, the goldens regenerated and claimed, and the tree clean.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t4-fix1.md <<EOF
Runner: T4 fix 1 done. Commits: <n>. Fixed: F-122, F-123, F-124, F-125. The long name: <the name you used, the width, and that the footer is still last>. The automation's honesty: <the line the block carries now>. The rows that add up: <the arithmetic for both rows>. The doc: <what docs/05 says now>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t4 fix 1"
herdr agent prompt "w1H:pJ" "Runner: T4 fix 1 done. Commits: <n>. Fixed: F-122, F-123, F-124, F-125. The long name: <the name you used, the width, and that the footer is still last>. The automation's honesty: <the line the block carries now>. The rows that add up: <the arithmetic for both rows>. The doc: <what docs/05 says now>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
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
result, and the arithmetic.
