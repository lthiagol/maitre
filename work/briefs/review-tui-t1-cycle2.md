You are the **reviewer** for **T1** — cycle 2. The same lane that reviewed cycle 1.
The orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command below.

Read [work/tui/README.md](../tui/README.md) for the state and the protocol, and
[work/briefs/tui-t1.md](tui-t1.md) for the slice, plus
[work/briefs/fix-tui-t1-1.md](fix-tui-t1-1.md) for the fix cycle you are reviewing.

## Cycle 2 is **blockers only**

File **HIGH or MEDIUM** blockers. No nits, no cosmetics, no style, and **do not
re-file what cycle 1 closed**. Two things are *not* blockers: a LOW note, and a
question that has an answer somewhere you have not looked yet — go and look.

## What cycle 1 decided, and what happened after

Cycle 1's verdict was **`ok`, 0 blockers**, and the slice's six acceptance criteria
held. The orchestrator then took **four** of cycle 1's five LOW notes and ruled them
into a fix cycle — two of them **upgraded to medium**, because checking them showed
they were not later-round material:

| # | What | Why it was upgraded |
|--:|---|---|
| **F-99** | the name the console renders is sanitised | the plan's own `project.json` name reaches the header and the window title, and `maitre-cli init` writes that name from the directory's basename — so a directory with a newline in it breaks the header into two rows |
| **F-100** | the plan's name has one answer | three states answered "what is this plan called?" three ways, and the third answered an **absolute path**, which is what lands in the window title |
| **F-101** | the row labelled `root` shows the root | it showed the plan directory's **parent**, so `kitchen/` was nowhere in the console — and the row's label did not name what it showed |
| **F-102** | `docs/05` says what the window title carries | it said the console *owns* the window title, and nothing said what it carries |

**The numbers F-99..F-102 are reserved for those four**, and cycle 2 writes their rows
in the ledger. The runner's fix cycle landed in five commits:
`5a79638` · `67b455d` · `b4a3052` · `55bdc81` · `e218ec3`.

## What to verify — each by command, not by reading

### 1. F-99 — the name is sanitised

- The header and the window title show a **one-row** name for a plan whose directory
  basename carries a newline — and **the file on disk keeps what was written**
  (`project.json` still holds the newline). The sanitising is a rendering concern.
- **Mutate it and prove the test bites:** take the `sanitize` call out of the name's
  path in `internal/ui/identity.go`, and the suite must fail. Revert it, and confirm
  the tree is exactly as the runner left it.

### 2. F-100 — the name has one answer

- A plan whose `project.json` names it **empty** and a plan whose `project.json` is
  **missing** show the **same** answer — and that answer is a **name**, never a path.
- **The answer lives in the domain.** `core.Store.ProjectName` is the one place that
  answers; `core.LoadProject` answers the same for an empty name as for a missing
  file; and the console **asks** rather than answering a second time. Check that the
  console really does ask — `grep` it, and read `internal/ui/identity.go`.
- **The writer is untouched.** `core.Init` still writes what it wrote before, and the
  plan format on disk does not move. `core.DomainVersion` is **7**, and
  `TestDomainShapeIsVersioned` passes.
- **The CLI did not change.** `maitre-cli plan` prints the name the domain now
  answers — check that is the *whole* of the change, and that no verb was added.

### 3. F-101 — the row labelled `root` shows the root

- The Settings tab's `root` row shows the **plan directory**, humanised —
  `~/Code/Pessoal/maitre/kitchen` — and **the themes are still listed**.
- The row is labelled `root`, and the domain's root is the plan directory:
  `core.FindRoot` returns it, `store.Root` holds it, and the glossary names
  `kitchen/` *the plan directory*.

### 4. F-102 — `docs/05` says what the window title carries

- The paragraph is in `docs/05-human-interface.md`, it says what the window title
  carries, and it says the window title is **not** the screen's row 0.
- T1's brief was corrected in the same commit, and the correction says so.

### 5. The scope held

| | |
|---|---|
| **The goldens** | unchanged — `make goldens` leaves `git diff` empty. A golden that moved for no stated reason is the slice saying it did more than it said |
| **The dependency budget** | **0 / 0 / 28** against **0 / 0 / 30** |
| **The tabs** | seven, the same names |
| **`internal/core` and `internal/cli`** | stdlib-only; the domain layer has zero UI imports |
| **The art** | the mark's drawings live in `internal/brand` and nowhere else |
| **`docs/audit.md`** | the ledger is **yours**. The runner touched no row |

## Two things the orchestrator observed, and did not rule on

**These are observations, not findings.** You judge them. Each one is cheap to check
by command — and if you think either is a blocker, file it; if you think it is a note,
say so in your report and it goes in the ledger for a later round.

### A. The `root` row renders the plan's **path** unsanitised

The fix cycle sanitised the **name**. The row that shows the **path** does not go
through the same door — and a plan directory whose name carries a newline splits the
row:

```
  root        /tmp/vfy/we
ird/kitchen
```

The path is not agent-written, so the question is a judgement: **does the console
survive a plan whose directory name carries a newline?** The header now does; the
Settings row may not.

### B. The state row's count and its list disagree

`work/tui/README.md`'s T1 row says **12 commits** and names **nine**. The three it
does not name are the three commits that wrote the row — and a commit cannot carry
its own hash, but the *next* row-writing commit can name the previous one.

Cycle 1 filed this as a LOW note: *"the state row's commit list is short by its own
commits... The count is right; the list is not complete."* The fix cycle repeated it.

The row is the **resume record**: the thing a cold session reads to know where the
work stopped. The question is what its shape should be so that it cannot drift again:
**the count counts what the list names**, or the list names every commit it can, or
the row says the count excludes the commits that write it.

## One question with an answer to look for

**Does the plan's name belong in `docs/03` § 8?**

The house rule is that a view is registered in `core.DerivedViews` **and** written
down in `docs/03` § 8, and the wire checks both directions. `core.Store.ProjectName`
is **not** marked `// view`, so the coverage wire accepts it under the `Store.*`
family — `make check` is green, and that is the instrument.

But `docs/03` § 8 is the table of **derived views**, and the house rule about *one
fact, one answer* says a fact the console shows is computed in the domain. The plan's
name is a fact the console shows. Look at how § 8 treats the other store readers —
`LoadProject`, `ListRecipes` — and say whether the plan's name belongs in that table
or beside it. **If you file it, file it as a MEDIUM**: the doc and the code would
disagree about what the domain computes, and the doc drift is the class this project
keeps paying for.

## The ledger — `docs/audit.md` is yours

- **F-99..F-102** get their rows: closed, with the commit that closed each and the
  test that holds it.
- The **slice T1** row and the **cycle 1** row are there. Add the **cycle 2** row.
- Anything you file is **F-103** onward.
- **Never sweep the ledger.** Its rows keep the names as they were when each finding
  was filed.

## The verdict

| | |
|---|---|
| No HIGH or MEDIUM findings | `verdict: ok` |
| Any HIGH or MEDIUM finding | `verdict: changes-needed` |

**Do not invent reasons to fail.** A slice that does what it said is `ok`.

## Hard rules

- Do **not** edit implementation code — **except** `docs/audit.md`, which is yours.
- Revert every mutation you make, **including** the ones the check asks for. The tree
  must be exactly as the runner left it, except your ledger commit — confirm it with
  `git status --short` and `git diff` at the end.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
- **One commit** for the ledger. One commit for the receipt.

## Notify the orchestrator

**The receipt comes FIRST** — the file cannot be lost, the nudge is only the bell.

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-tui-t1-cycle2.md <<EOF
Reviewer: T1 cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-103.. or none>). F-99..F-102 verified: <how, and the mutation that proved each test bites>. The path row: <your judgement>. The state row: <your judgement>. The name in docs/03: <your judgement>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t1 cycle 2"
herdr agent prompt "w1H:p1" "Reviewer: T1 cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-103.. or none>). F-99..F-102 verified: <how, and the mutation that proved each test bites>. The path row: <your judgement>. The state row: <your judgement>. The name in docs/03: <your judgement>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.

**Nothing else waits on you.** After the nudge, end your turn.
