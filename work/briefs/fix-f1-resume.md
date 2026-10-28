You are the **runner** — the same lane that wrote slice F and its first fix. **Your turn
ended without reporting.** This brief is the continuation. The orchestrator's pane id
is `w1Z:p1`; use it verbatim in the notify command below.

## Where you are

| | |
|---|---|
| **Landed** | `c9e311f` **F-87** · `84d06f3` **F-88** · `062192d` **F-89** — three of the four findings |
| **Outstanding** | **F-90** and the LOW notes |
| **The tree** | **clean** — only `?? logo/`. Nothing is lost; every commit landed |
| **The receipt** | not written. That is the half of the report that cannot be lost, and it is why the protocol asks for it **before** the nudge |

**Do not start over. Do not revert.** Three findings are done and correct.

## What is left — F-90 and the notes

The full text of F-90 is in [fix-f1.md](fix-f1.md) § *F-90*. In short, the new
"one format" wire has **three doors open**: it scans `internal/core` only, it exempts
`model.go` **wholesale**, and its token list omits the retired id literals
(`R01`/`S01`/`N01`). The reviewer's mutations show it — `'specials'` in
`internal/core` fails the wire; in `internal/ui`, in `model.go`, or as `R01`, it
**passes**.

**The fix:**

1. **The wire scans the whole tree**, not one package.
2. **The exemption is narrowed** to where the retired names legitimately live — the
   **detection that refuses them** and the **deny-list that refuses them as names** —
   and it is narrowed **per place**, so a retired name cannot hide in `model.go`
   behind the deny-list's exemption.
3. **The token list is complete**: the retired *id* literals join it, and so does
   every other retired token the code could grow back.
4. **The token choice is unambiguous.** The bare word `evidence` cannot be a token —
   it is a **current** word (the record's name) and the retired stamp field was called
   that too. Pick a token that names the retired thing and nothing else.
5. **The wire can fail** — prove it, on a token in each package it now scans.

## The LOW notes that ride with it

| Note | The ruling |
|---|---|
| **the file name's id form is not checked** | a file named `R01-….json` in the current format is a retired id form on disk, and the lookup normalises it into the current record. **One format means the file's name is in the current format too.** Decide what the tool does with a retired form in a file name — refuse it, or say what it is — and make it a rule with a test. |
| **a plan with a retired artifact can still be written to** | the plan is not supported: *not supported* must mean **no writes**, not only a read-time message. The detection must gate the write path. |
| `GLOSSARY.md:97` names `coordinator` in prose | the glossary's rule is that a retired word is named **once**, in the `_Avoid_` line. |
| `activity.go`'s version-zero comment is stale | a comment that no longer matches the code is a bug. |

## The tree is shared

**Another session is drawing the mascot in this repository** and has files in the tree
(`logo/`). Stage **by name** — never `git add -A`, never `git add .`. Leave anything
you did not write alone. If a file you did not write has changed, **stop and report
it**.

## Not yours — do not touch

| | |
|---|---|
| `docs/audit.md` | the reviewer's to write. |
| `docs/07-vocabulary.md`, `work/` | the orchestrator's. |
| The gate table, the plan format, the command count (42) | F-89 already moved the version and the facts list. Nothing else. If a fix would touch anything else, **stop and report it**. |

## The gate

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

Green.

## Two commits — F-90, then the notes

```
Per-AC evidence:
- F-90: the wire scans the whole tree and its exemptions are per place; go test ./internal/core -run 'OneFormat|Retired' -v exit 0
```

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-f-fix-1.md <<EOF
Runner: slice F fix 1 done. Commits: <n>. Fixed: F-87..F-90 and the notes. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-f fix 1"
herdr agent prompt "w1Z:p1" "Runner: slice F fix 1 done. Commits: <n>. Fixed: F-87..F-90 and the notes. Notes: <anything the reviewer should know>."
```

**Write the receipt, commit it, and *then* nudge.** The nudge only reaches an idle
orchestrator; the receipt is the half that cannot be lost.

**One write per file, and commit as soon as the file is coherent** — the rule that
would have shortened this turn.
