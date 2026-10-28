You are the **runner** — the same lane that wrote slice H. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## The fix cycle for slice H — one blocker and seven notes

The reviewer filed **F-98**. The orchestrator ruled: **accepted**, and **the seven LOW
notes are fixed in the same cycle** — they are all the same class of thing: a place
where the gate could be wrong and nobody would know.

### F-98 (medium) — `make test-race` silently drops the filter

> `make test-race PKG=./internal/brand RUN='-run NoSuchTest'` **exits 0**, while
> `make test` with the same filter exits 2.

The cause is visible in the `Makefile`:

```make
test:
	GO=$(GO) scripts/go-test.sh "$(RUN)" "$(PKG)"

test-race:
	GO=$(GO) scripts/go-test.sh "-race -shuffle=on -count=2" "$(PKG)"
```

**`test-race` never passes `$(RUN)`.** The flags land in the script's *filter* slot and
the package in its *package* slot — so the filter is empty, **and the zero-match check
never fires**, and a run that tested **nothing** exits 0. The slow gate is the one CI
runs: a gate that can pass having tested nothing is the worst kind of broken gate.

**The fix — one place decides how a test run is invoked.** The two targets may differ
in **the flags they add** and in nothing else. The **filter** and the **package** are
always passed, and the **zero-match check applies to both**.

**Prove it:** `make test-race PKG=./internal/brand RUN='-run NoSuchTest'` must **fail**,
and `make test-race PKG=./internal/brand RUN='-run TestTheDrawingsMatchTheSource'` must
pass **and run that test**. And the two targets must not be able to drift again: if
they are two recipes that spell out the invocation separately, they will.

### The seven notes

| Note | The ruling |
|---|---|
| **the tool's version is never checked** | the gate runs whatever `staticcheck` is on the `PATH`, which may be an older one. **The gate checks the version it runs** against the pin, and fails with *"run `make tools`"* if it does not match. A gate that runs a tool it did not pin is a gate nobody knows what it ran. |
| **the tool versions are written twice and only one is used** | the `Makefile` pins them and `docs/testing.md` names them — and the doc's copy is decorative. **The doc's version is the *checked* copy**, the same shape as the dependency budget: the gate holds the pin and checks the doc against it. Say which direction you chose and why. |
| **`dep-budget.sh` ignores `GO`** | the `Makefile` honours `GO ?= go`; the script calls `go` outright. **It honours `GO`**, like every other target. |
| **`tidy-check` leaves the tree dirty if tidy itself fails** | if `go mod tidy` fails halfway, `go.mod` is half-changed and the gate stops with a dirty tree. **The tree must be as it was**: snapshot `go.mod`/`go.sum` before, restore them if tidy fails. |
| **`install-hooks` overwrites silently** | it says nothing when it replaces an existing hook. **It says so** — and it says what it replaced. |
| **`make -j` is unsafe** | two targets touch the same files. **`.NOTPARALLEL:`** in the `Makefile`, and a comment saying why. |
| **the ledger's discipline table still says to run `go test` by hand** | **`docs/audit.md` is the reviewer's** — the reviewer said so itself. It is theirs to update, and the cycle-2 brief will ask them to. **You do not touch it.** |

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. The tools do not touch `go.mod`.
- The domain layer has **zero UI imports**.
- One fact, one answer: the pins, the numbers, the versions — each in one place, with
  the other copies checked.
- `core.DomainVersion` is **7**. Do not move it.
- The gate table, the plan format and the command count (**42**) must not move.

## The scope — what not to touch

| | |
|---|---|
| **The code** | **no behaviour changes.** This cycle is the gate's own plumbing. If a fix would change what the tool does, **stop and report it**. |
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

## One commit for the blocker, one for the notes

```
Per-AC evidence:
- F-98: make test-race honours the filter and fails on a zero-match run; make test-race PKG=./internal/brand RUN='-run NoSuchTest' exits 2
```

**One write per file, and commit as soon as the file is coherent.**

## Finish with the receipt — update the one you wrote when you started

```bash
cat > work/inbox/runner-slice-h.md <<EOF
Runner: slice H fix 1 done. Commits: <n>. Fixed: F-98 and the seven notes. The tool versions live in: <where, and why>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-h fix 1"
herdr agent prompt "w1Z:p1" "Runner: slice H fix 1 done. Commits: <n>. Fixed: F-98 and the seven notes. The tool versions live in: <where, and why>. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.**
