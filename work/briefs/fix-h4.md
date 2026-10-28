You are the **runner** — the same lane that wrote slice H. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## The last round for slice H — two notes

The review came back **`ok`** for the fourth time. The orchestrator ruled: **these two
are fixed here, and this is the last round for slice H.** The other two notes are
carried to the ledger.

### 1. The gate cannot see an untracked golden

> `git diff --exit-code HEAD -- internal/ui/testdata` does **not** see **untracked**
> files, so a new golden that a test writes but nobody commits **passes the gate**.

**The rule: a gate must not say green when it checked nothing.** This is F-98's shape
again — the check that cannot see the thing — and it is the one the repo has paid for
most often.

**The fix:** the check must see an untracked golden. The remedy is the git one — tell
git to **intent-to-add** the golden files before the diff, so a new file appears in it
as a change rather than as nothing. (`git add -N` is the command; the runner's call
how to wire it.)

**Prove it:** write a **new** golden file that no test claims (temporarily), and
confirm `make goldens-check` **fails** — and that the message names the file. Then
write a new golden file that **is** claimed but is not committed, and confirm the
check **also** fails, because a file nobody committed is not in the commit. Revert
both.

**And the tree must stay clean** — the check must not leave the intent-to-add staged
behind it.

### 2. The orphan message does not say what to do

> the orphan message names the file and the two causes but **not the remedy** — so a
> new golden whose test is not named `...Golden` will be reported, a true positive by
> design that reads like a false alarm to whoever added it.

**Every message in this gate tells the reader what to do.** The budget's message names
the file to change; the version message says both halves; the tool's message says where
the binaries land. This one names a cause and stops.

**The fix:** the orphan message names the **remedy** — the two ways out, which are
the two ways a golden stops being claimed:

1. **name the test so the filter matches it** (the filter is `-run Golden`, so the
   test's name must carry `Golden`), or
2. **extend the filter** so it matches the test.

And it must say which of the two causes it is, so the reader knows which remedy
applies. A message that leaves the reader guessing is a message that will be ignored
the next time.

### The two notes that are carried as they stand

| Note | Ruling |
|---|---|
| after a **failed** `goldens-check`, an orphan golden keeps a `2000-01-01` mtime on disk | **carry it.** A filesystem artifact after a *failed* run is not a git change, and the next `make goldens` rewrites it. The gate already failed — the chef is looking at the failure, not at a green. |
| the orphan check's scope is `internal/ui/testdata` (the `GOLDENS` variable) | **carry it.** The scope is a variable; it grows with the need. Say in the script's comment that the scope is the variable, so the next person knows where to widen it. |

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget.
- One fact, one answer: **one place decides how a test run is invoked.**
- `core.DomainVersion` is **7**. Do not move it.
- The gate table, the plan format and the command count (**42**) must not move.

## The scope — what not to touch

| | |
|---|---|
| **The code** | **no behaviour changes.** This round is the gate's own plumbing. |
| **`docs/audit.md`** | the reviewer's. |
| **`docs/07-vocabulary.md`, `work/`** | the orchestrator's. |
| The names, the glossary, the spec | settled. Do not touch. |

## The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. Leave anything you did not
write alone. If a file you did not write has changed, **stop and report it**.

## The gate

```bash
make check && make check-race && git status --short
```

Both green, and the tree clean.

## One commit

```
Per-AC evidence:
- goldens-check sees an untracked golden, and the orphan message names the remedy; make check && make check-race exit 0
```

**One write per file, and commit as soon as the file is coherent.**

## Finish with the receipt — update the one you wrote when you started

```bash
cat > work/inbox/runner-slice-h.md <<EOF
Runner: slice H fix 4 done. Commits: <n>. The gate sees an untracked golden; the orphan message names the remedy. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-h fix 4"
herdr agent prompt "w1Z:p1" "Runner: slice H fix 4 done. Commits: <n>. The gate sees an untracked golden; the orphan message names the remedy. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.**
