You are the **runner**, again — the same lane that wrote T1 and its fix cycle 1. The
orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command below.

This is **T1's fix cycle 2**. It is small: **one blocker, two notes, one question**.

## Cycle 2's verdict

`changes-needed`, **one blocker: F-103**. The cycle-1 fixes (F-99..F-102) were all
verified and hold. F-103 is a real one, and the reviewer found it by mutation — so
read the proof before you fix it.

## F-103 (medium, blocker) — the name's rule is written twice, and only one copy is pinned

**The target:** the rule that answers *"what is this plan called?"* is written **once**,
and **every** surface that shows the name reads that one answer.

**The hole, proved by command.** The rule lives in two places:

| Where | What it says |
|---|---|
| `core.LoadProject` | a plan whose `project.json` names it **empty** answers with the project's directory name |
| `core.Store.ProjectName` | `p.Name == ""` answers with the project's directory name — **again** |

Only the second copy is pinned by a test. `core.LoadProject` has **no test of its own**,
so the first copy can be deleted and nothing notices. The reviewer's mutation, which
the orchestrator reproduced:

```bash
# take the empty-name branch out of LoadProject
go test ./...          # GREEN — the whole suite
maitre-cli plan        # "Name": ""     <- the CLI's answer
maitre --render        # ⋈ maitre · demo2  <- the console's answer
```

**Two surfaces, two answers, and a green gate.** That is the house's standing
failure — *two answers to one question* — and the ledger has three rows of it already.

**What lands:** one copy of the rule, and it is pinned. Choose which copy survives —
that is a judgement, and you own it:

- **Drop `ProjectName`'s redundant branch** and let `LoadProject` be the one place
  that answers. Then **`LoadProject`'s answer is pinned by a test of its own**, and
  `ProjectName` is the console's way of asking.
- **Or keep the rule in `ProjectName` alone** and have `LoadProject` answer raw.

**Either is fine. What is not fine is both.** Say in your report which you chose and
why.

**The evidence this finding demands:**

| | |
|---|---|
| the **mutation fails** | take the rule out of whichever copy you did not keep — **the suite must go red**. A green suite with the rule broken is the defect, not the fix |
| the **two surfaces agree** | `maitre-cli plan` and `maitre --render` show the **same** name for the same plan, in all three states: the plan names itself · the file names it empty · the file is missing |
| **the answer is pinned where the rule lives** | the copy that answers has a test that fails when the rule is removed from **that** copy |

**The CLI's behaviour may change, and that is allowed.** Today `maitre-cli plan` prints
`"Name": ""` for a plan that names nothing, while the console shows the project's
directory name. After the fix, both show the same thing. If your choice makes the CLI
print the same name as the console, say so in your report — that is the point.

## F-104 (low) — the `root` row renders the plan's path unsanitised

**The target:** a newline anywhere in the path the Settings tab shows does not split
the row.

The fix cycle sanitised the **name**. The **path** goes through a different door, and
a plan whose directory name carries a newline splits the row:

```
  root        /tmp/vfy/we
ird/kitchen
```

**What lands:** the path the console renders goes through the same sanitising as
everything else the console renders. The path is not the plan's content — the review
agreed — but the console still draws it on one line.

## F-105 (low) — the state row's count and its list must agree

**The target:** the row cannot be read two ways.

Cycle 1 filed this as a note, and the fix cycle repeated it: the row says **12
commits** and names **nine**. The three it does not name are the three commits that
wrote the row — and a commit cannot carry its own hash, but the *next* row-writing
commit can name the previous one.

The row is the **resume record**: the thing a cold session reads to know where the
work stopped. **What lands:** the row's shape is one that cannot drift again — the
count counts what the list names, **or** the list names every commit it can and the
row says the count excludes the commits that wrote it. **Say which you choose.**

This is the row **you** write, in **this** cycle: `work/tui/README.md`'s T1 row. The
orchestrator has already set it to `🔄 fix cycle 2`; your commit moves it to closed
with the shape you chose.

## The question the reviewer answered — nothing to do, and why

> *the name in `docs/03`: it does NOT belong in § 8 — that table is the derived
> views (computed, never stored) and the name is a stored fact the domain reads,
> like `LoadProject` and `ListRecipes`, which the table does not carry either; no
> MEDIUM filed on that question, but the rule itself is nowhere in the spec.*

**The reviewer is right, and nothing lands here.** The plan's name is not a derived
view: it is stored, in `project.json`, and read. `docs/03` § 8 carries the views the
domain *computes*, and it does not carry `LoadProject` or `ListRecipes` either. So the
name does not belong there.

**But the reviewer also said the rule is nowhere in the spec** — and that is worth
saying where it can be found. **Where the rule should live** is the chef's call, not
yours and not the orchestrator's. **Do not write it into `docs/03`.** If you think
the rule belongs somewhere specific, say so in your report and the orchestrator will
put it to the chef.

## The scope — what not to touch

| | |
|---|---|
| **The CLI** | no verb is added, and `maitre-cli`'s command surface does not change. `maitre-cli plan` prints the name the domain answers, and that is the whole of the change |
| **The writer** | `core.Init` still writes what it wrote. The change is in the **reader** |
| **`core.DomainVersion`** | stays **7**. The plan format on disk does not move: the same fields, the same values. `TestDomainShapeIsVersioned` is the instrument — if it fails, **stop and report** |
| **The tabs** | seven, the same names |
| **The goldens** | `make goldens` regenerates; `make goldens-check` holds them. **Expected unchanged** — the header's words do not move. If a golden moves, say so and say why |
| **`docs/03`** | **this cycle touches no `docs/03`.** The name is not a view, and the question of where the rule lives is the chef's |
| **`docs/05`** | the paragraph cycle 1 added stays |
| **`docs/audit.md`** | the ledger is the **reviewer's**. This cycle files no finding and touches no row |
| **`work/`** | the orchestrator's, except the state row (F-105) |

## The scope — what lands

| | |
|---|---|
| **`internal/core`** | one copy of the name's rule, and it is pinned (F-103) |
| **`internal/ui`** | the path the console renders is sanitised (F-104) |
| **`work/tui/README.md`** | the state row's shape, so it cannot drift again (F-105) |

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget (**≤ 30**). **No dependency is added by this cycle.**
- **The domain layer has zero UI imports.**
- **`core.DomainVersion` is 7.**
- The tab count is 7 and the tab names do not change.
- **One fact, one answer.**
- The gate is `make check`. A focused check is `make test PKG=… RUN='-run …'`.

## The evidence this cycle must produce

For **each** finding, a command and its output — not a claim:

| Finding | The evidence |
|---|---|
| **F-103** | the **mutation fails**: removing the rule from the copy you did not keep turns the suite **red**. And the two surfaces — `maitre-cli plan` and `maitre --render` — show the **same** name in all three states |
| **F-104** | a plan whose directory name carries a newline shows a **one-row** `root` row |
| **F-105** | the row's count and its list agree, and you say in your report which shape you chose and why |

Plus: `make check` green, `make check-race` green, the goldens unchanged, and the tree
clean.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t1-fix2.md <<EOF
Runner: T1 fix 2 done. Commits: <n>. Fixed: F-103, F-104, F-105. The rule lives in: <which copy you kept, and why>. The mutation: <what you removed, and the suite's answer>. The two surfaces: <what maitre-cli plan and maitre --render show, in all three states>. The path is sanitised: <what you proved it with>. The state row: <the shape you chose, and why>. DomainVersion: <7, and the version test's result>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t1 fix 2"
herdr agent prompt "w1H:p1" "Runner: T1 fix 2 done. Commits: <n>. Fixed: F-103, F-104, F-105. The rule lives in: <which copy you kept, and why>. The mutation: <what you removed, and the suite's answer>. The two surfaces: <what maitre-cli plan and maitre --render show, in all three states>. The path is sanitised: <what you proved it with>. The state row: <the shape you chose, and why>. DomainVersion: <7, and the version test's result>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
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
result, and what the two surfaces show.
