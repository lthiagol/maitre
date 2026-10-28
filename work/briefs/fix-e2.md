You are the **runner** — the same lane that wrote slice E. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## The last fix cycle for slice E — one finding and two notes

The reviewer came back with **F-86 (medium)** and two LOW notes. The orchestrator
ruled: **all three are fixed here.** F-86 is a beautiful finding — read it, because
it is the kind of thing that is worth remembering.

### F-86 (medium) — the instrument was blind, not the tree

> *the sweep's command and AC-04's list of old names omit* maître d'*, so one of the
> seven renamed terms is outside the instrument's reach; the tree is clean (only dated
> records), and the remedy is the command's: carry the name, exclude the glossary's
> `_Avoid_` lines.*

**The tree is clean.** Nothing needs changing in the docs. **The command could not
see one of the seven renamed terms** — *maître d'*, the console's old themed name.

> An instrument that cannot see a name is not a check. The finding is not "the tree
> is dirty"; it is "the sweep would have passed whatever the tree said about
> *maître d'*, because it never asked."

**The fix:**

1. **The sweep's command carries every one of the seven names** — *maître d'* among
   them — and it **knows the difference between a retired word and a warning about a
   retired word**: the glossary's `_Avoid_` lines name them on purpose, so the sweep
   excludes those lines, and a retired name outside one is drift.
2. **Re-run the sweep with the complete command and show its output** in the commit
   message — the actual output, not a claim about it. The same rule that F-81 taught:
   a command offered as evidence must reproduce.
3. **Say what the count is now, and reconcile it.** The reviewer counted **77 lines
   in 19 files** and found the arithmetic holds (78 → 77 is F-83's one prose line
   removed, nothing added). A command that now also carries *maître d'* will find the
   glossary's `_Avoid_` line — so the count will move, and the commit must explain
   the movement.

**The brief's list is the orchestrator's and is corrected** — *maître d'* is in it
now. Yours is the command.

**No script, deliberately.** The orchestrator considered making the sweep a permanent
tripwire and decided against it, and wants the reason on the record: the sweep is a
**one-time act** over a tree that now has *legitimate* mentions of the old names —
the legacy reader, the retired-role aliases, the dated records. A permanent
tripwire would need a permanent exclusion list, and **a tripwire with a long
exclusion list is a tripwire that has stopped seeing**. The repo already has better
instruments for the live facts (the deny-list wire, the derived-views wire, the
plan-format wire). Do not add a script.

### The two LOW notes

| Note | Ruling |
|---|---|
| the console's entry sits under **"### The strings"** — the section says what the console *carries*, but the heading alone can be misread: the console **is not a string** | **fix it.** A glossary is a **taxonomy**, and a taxonomy that files the console under "the strings" is wrong in the same way F-82 was wrong when it filed the console among the crew. The heading must not claim the console is a string. You choose the heading — say which, and why. |
| `docs/03` § 6 links to `#the-version-guard`, and **there is no such heading** — a broken pointer, pre-existing | **fix it.** The pointer must point at something that exists. The guard's own description is two lines above the link, and the guard lives in `internal/core/lock.go` (`checkDomainVersion`). Point it at whichever is right, and say which. |

## The tree is shared — read this before you commit

**Another session is drawing the mascot in this repository** and has files in the tree
(`logo/`).

| | |
|---|---|
| **Stage by name** | name your files explicitly. **Never `git add -A`**, never `git add .` |
| **`logo/`, and anything you did not write** | leave it alone |
| **`kitchen/`** | re-read before you commit; the other session may write to the plan |

**If a file changed that you did not change, stop and report it.**

## Not yours — do not touch

| | |
|---|---|
| `docs/audit.md` | the reviewer's to write. It files F-86 and closes it. |
| `docs/07-vocabulary.md`, `work/briefs/slice-e.md` | the orchestrator's. Both corrected. |
| `README.md`, `work/README.md` | the orchestrator's. Already corrected. |
| The gate table, the plan format, the command count (42), `core.DomainVersion` (6), `kitchen/activity.json` | none of the three touches them. If a fix would, **stop and report it**. |

## The gate

```bash
make lint && gofmt -l . && make test && git status --short
```

Green.

## Two commits — one per change

```
Per-AC evidence:
- F-86: the sweep's command carries all seven retired names and excludes the glossary's _Avoid_ lines; the output is in the commit body
```

## The report — the receipt comes first

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-e-fix-2.md <<EOF
Runner: slice E fix 2 done. Commits: <n>. Fixed: F-86 and the two notes. Sweep count now: <n>, and it moves by <why>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-e fix 2"
herdr agent prompt "w1Z:p1" "Runner: slice E fix 2 done. Commits: <n>. Fixed: F-86 and the two notes. Sweep count now: <n>, and it moves by <why>. Notes: <anything the reviewer should know>."
```

**One write per file, and commit as soon as the file is coherent.**
