You are the **runner** — the same lane that wrote slice H. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## The last fix for slice H — the ruling, applied twice

The review came back **`ok`**, with six LOW notes. The orchestrator ruled on them, and
**two of them are the same finding as F-98** — the fix was applied to `test` and
`test-race` and not to the third invocation. So this cycle is short.

### The goldens are a third invocation, and it bypasses the refusal

> `make goldens`/`make goldens-check` run a raw `go test -run Golden -update` instead of
> `scripts/go-test.sh`, so the zero-match refusal does not apply there — only the
> stale-golden check would go blind.

**The ruling that was already made: *one place decides how a test run is invoked*.**
The goldens are a test invocation. So:

1. **The goldens go through the same script**, so the **zero-match refusal applies
   there too**: a `goldens-check` that regenerated nothing because the golden tests had
   vanished must **fail**, not pass with an empty diff.
2. **The golden invocation is written once.** `goldens` and `goldens-check` spell it
   out twice today — two spellings, two answers. They become two targets of **one**
   recipe, as `test` and `test-race` are now.

**Prove it:** a `goldens-check` whose golden test has been renamed away (temporarily)
must **fail**. Revert it.

### Three one-liners that ride along

They are in the files this cycle already touches, and one of them names a fix that does
not fix anything.

| Note | The ruling |
|---|---|
| **`AGENTS.md`'s `make check` row omits `versions-check`** | the doc is behind the code — in the safe direction, but behind. **Add it.** |
| **the version failure says "run `make tools`" but not "and put `GOPATH/bin` first on the `PATH`"** | the message names a fix that does not fix the problem: an older tool earlier on the `PATH` still fails after `make tools`. **The message must say both** — and it should point at what `make tools` already says about where the binaries landed. |
| **the doc check looks only at the *first* matching row** | so a second mention of the same tool in the doc could drift unchecked. **Check every mention**, not the first. |

### The note that is accepted as it stands

> `check-tool-versions.sh` calls `require-tool.sh` by a relative path, so it only works
> from the repo root.

**Accepted, and the reason is the chef's ruling: the interface is `make`.** The
scripts are `make`'s implementation, and `make` runs its recipes from the repo root.
Every script in `scripts/` behaves this way. **No change** — and say so in your report,
so the next reader does not "fix" it.

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. The tools do not touch `go.mod`.
- One fact, one answer: **one place decides how a test run is invoked.**
- `core.DomainVersion` is **7**. Do not move it.
- The gate table, the plan format and the command count (**42**) must not move.

## The scope — what not to touch

| | |
|---|---|
| **The code** | **no behaviour changes.** This cycle is the gate's own plumbing. |
| **`docs/audit.md`** | the reviewer's. The reviewer already updated its discipline table in cycle 2. |
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
- the goldens go through the same invocation as every other test run, and a goldens-check that regenerates nothing fails; make check && make check-race exit 0
```

**One write per file, and commit as soon as the file is coherent.**

## Finish with the receipt — update the one you wrote when you started

```bash
cat > work/inbox/runner-slice-h.md <<EOF
Runner: slice H fix 2 done. Commits: <n>. The goldens go through the same invocation; the version failure says how to fix it; the doc check reads every mention. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-h fix 2"
herdr agent prompt "w1Z:p1" "Runner: slice H fix 2 done. Commits: <n>. The goldens go through the same invocation; the version failure says how to fix it; the doc check reads every mention. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.**
