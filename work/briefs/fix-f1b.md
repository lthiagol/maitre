You are the **runner** for a **small, focused fix**. The orchestrator's pane id is
`w1Z:p1` — it is already substituted in this file; use it verbatim in the notify
command below. The repository is `/Users/thiago/Code/Pessoal/maitre`.

## Step 0 — write your receipt FIRST, before any code

```bash
cat > work/inbox/runner-slice-f-fix-1b.md <<EOF
Runner: starting F-90 and the notes. Commits so far: 0. Notes: nothing done yet.
EOF
git add work/inbox && git commit -q -m "inbox: runner starts slice-f fix 1b"
```

**Then do the work, and update that receipt when you finish.** Two lanes in this
project have ended a turn without reporting; a receipt written *before* the work is
how the next one knows where the lane stopped. It costs one commit.

## What you are fixing

The reviewer filed four findings against slice F. **Three are already fixed** — F-87
(a retired name cannot buy a gate), F-88 (the operations doc), F-89 (the digest's
encoding is a format fact and the version moves to 7). The commits are in the log:

```bash
git log --oneline -6
```

**You are doing F-90 and the notes that ride with it.** The other lane's turn ended
before it reached them; the tree is clean, so there is nothing to recover.

### F-90 — the "one format" wire has three doors open

`internal/core/retired_test.go`'s `TestTheRetiredNamesAreOnlyWhereTheyBelong` is
the wire for the ruling *"one format"*. It currently:

1. scans **`internal/core` only** — a retired name in `internal/ui` passes;
2. **exempts `model.go` wholesale** — a retired name anywhere in that file passes;
3. has a **token list that omits the retired id literals** (`R01`, `S01`, `N01`) —
   so `R01` passes, and so does an invented token like `LegacyThing`.

**The fix:**

1. **The wire scans the whole tree** — every package, not one.
2. **The exemption is per place, not per file.** The retired names may live in
   exactly two places: the **detection that refuses them** (`retired.go`) and the
   **deny-list that refuses them as names** (`model.go`). Narrow the exemption to
   those places, so a retired name cannot hide in `model.go` behind the exemption.
3. **The token list is complete** — the retired id literals join it, and so does every
   other retired token the code could grow back.
4. **The token choice is unambiguous.** The bare word `evidence` cannot be a token:
   it is a **current** word (the record's name), and the retired stamp field was
   called that too. Pick a token that names the retired thing and nothing else.
5. **The wire can fail** — prove it, on a token in each package it now scans.

### The three notes that ride with it

| Note | The ruling |
|---|---|
| **the file name's id form is not checked** | a file named `R01-….json` in the current format is a retired id form on disk, and the lookup normalises it into the current record. **One format means the file's name is in the current format too.** Decide what the tool does with a retired form in a file name — refuse it, or say what it is — and make it a rule with a test. |
| **a plan with a retired artifact can still be written to** | the plan is not supported: *not supported* must mean **no writes**, not only a read-time message. The detection must gate the write path. |
| `GLOSSARY.md:97` names `coordinator` in prose | the glossary's rule is that a retired word is named **once**, in the `_Avoid_` line. |
| `activity.go`'s version-zero comment is stale | a comment that no longer matches the code is a bug. |

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps **28**
  external modules. The domain has **zero UI imports**.
- `core.DomainVersion` is **7** (F-89 moved it). Do not move it again.
- The **journal is append-only** — its history keeps its old names.
- The gate table, the plan format and the command count (**42**) must not move.

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
| The gate table, the plan format, the command count (42) | a fix may touch the **version** and the **facts list**. Nothing else. If a fix would touch anything else, **stop and report it**. |

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

**One write per file, and commit as soon as the file is coherent.**

## Finish with the receipt — update the one you wrote in step 0

```bash
cat > work/inbox/runner-slice-f-fix-1b.md <<EOF
Runner: slice F fix 1b done. Commits: <n>. Fixed: F-90 and the notes. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-f fix 1b"
herdr agent prompt "w1Z:p1" "Runner: slice F fix 1b done. Commits: <n>. Fixed: F-90 and the notes. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.** The nudge only reaches an idle orchestrator.
