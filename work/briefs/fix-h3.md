You are the **runner** — the same lane that wrote slice H. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## The last round for slice H — two notes, and then it closes

The review came back **`ok`** with six LOW notes. The orchestrator ruled: **two of them
are fixed here, and this is the last round for slice H.** The rest are carried to the
ledger as notes — the repo's own discipline is that a LOW is noted, not chased, and
this slice has had three rounds already.

### 1. The goldens' refusal is package-level

> the zero-match refusal is **package-level** — the filter `-run Golden` matches three
> tests in two files — so renaming the two in `golden_test.go` leaves
> `make goldens-check` green while the **recipes golden is never regenerated**.

**The hole is real**: a golden file nothing checks is a layout nothing pins, and the
check that says "the goldens are current" passes without touching it.

**The orchestrator's proposal — an *orphan* check, not a manifest.** The reviewer said
catching the partial case would take a manifest of the expected test names. It does
not: it takes a **count of orphans**, and the count is **derived**:

> **Every `.golden` file must be claimed by a test.**

A golden file that no test references is a golden that pins nothing. **Prove it:**
rename a golden test away (temporarily) and its file becomes an orphan — the check
must **fail**. Revert it.

**But the rule about the mechanism is the one that matters:**

| | |
|---|---|
| **A check that cries wolf is worse than no check** | the repo has paid for this twice: the old sweep command counted the glossary's `_Avoid_` lines as drift, and the version check once counted warnings as drift. **A false positive must fail *loudly* and name the file** — never silently pass, never falsely fail. |
| **If the orphan check is fragile, say so** | if the golden paths are not literal in the test files — if a test builds the path from a constant — then a reference-scan is a guess, and **a guess in a gate is worse than the hole**. **In that case: do not build it. Say so in your report, and the note is carried with your reasoning.** The orchestrator would rather have the hole recorded than a gate that fails on a Tuesday for no reason. |

**Either way, report which you did and why.**

### 2. The doc check can cry wolf about a version

> the doc check now scans any line mentioning the tool that carries a version-like
> token, so a future `staticcheck 2026.2.2` (the release version) would fail it
> spuriously — loudly, naming the offending version.

**The fix: the check reads the lines it means to read.** The tool reports **two**
versions — the release (`2026.2.2`) and the module (`v0.8.2`) — and the *pin* is the
module one. So the check must compare **the pinned version**, and it must not trip on
a line that carries the *other* version.

**Prove it:** put a line in `docs/testing.md` that mentions the tool with the **other**
version's shape, and confirm the check does **not** fail. And confirm it **still**
fails on a second mention carrying the **pinned** version's shape — that is the case
the fix for the note was about.

### 3. The note that is accepted as it stands

> `go-test.sh` now passes the package first, so a flag unknown to `go test` is handed
> to every test binary in a run, not only the goldens'.

**Accepted.** The ordering is deliberate — `-update` belongs to the test binary — and a
flag nobody knows fails loudly, in every binary, which is fine. **No change.**

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
- the goldens' check catches an orphaned golden file, and the doc check compares the pinned version only; make check && make check-race exit 0
```

**One write per file, and commit as soon as the file is coherent.**

## Finish with the receipt — update the one you wrote when you started

```bash
cat > work/inbox/runner-slice-h.md <<EOF
Runner: slice H fix 3 done. Commits: <n>. The orphan check: <built, or not built and why>. The doc check compares the pinned version only. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-h fix 3"
herdr agent prompt "w1Z:p1" "Runner: slice H fix 3 done. Commits: <n>. The orphan check: <built, or not built and why>. The doc check compares the pinned version only. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.**
