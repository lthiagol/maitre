You are the **reviewer** — the independent reviewer lane, fresh for slice H. The
orchestrator's pane id is `w1Z:p1` — it is already substituted in this file; wherever
this brief says `$ORCH_PANE`, it means `w1Z:p1`. Use it verbatim in the notify command
below. The repository is `/Users/thiago/Code/Pessoal/maitre`.

## Your slice — the review of slice H

Slice H built **the quality gate**: one interface (`make`), the tools the gate needs,
and CI. The chef's ruling: *we use just the make tool instead of running direct
testing code.*

```bash
cd /Users/thiago/Code/Pessoal/maitre
git log --oneline -12
git diff <base>..HEAD --stat      # the base is the commit before slice H's first commit
```

**Cycle 1: blockers only.** File **HIGH or MEDIUM** blockers. No nits, no cosmetics,
no style, no LOW notes filed as rows. **Do not re-file F-56, F-59, F-60, F-61, P-8**
(the chef's).

### 1. The gate — by command, not by reading

**Every target must exist and work on its own.** Run each one and see what it does:

`fmt` · `fmt-check` · `vet` · `staticcheck` · `lint` · `test` · `test-race` ·
`goldens` · `goldens-check` · `budget` · `tidy-check` · `check` · `check-race` ·
`tools` · `vuln` · `ci` · `install-hooks`

**And two properties that matter more than the list:**

- **The gate must not skip a check.** If `staticcheck` is not on the `PATH`, the gate
  must **fail** with a message that says how to get it — not skip it. **Prove it:**
  run the gate with the tool's directory removed from the `PATH` and confirm it fails
  loudly. *A gate that quietly skips a check is a gate that has stopped seeing.*
- **A focused run that matches nothing must fail.** `make test PKG=./internal/core
  RUN='-run NoSuchTest'` must **fail**, not pass with "no tests to run". **Prove it.**
  And prove `PKG` and `RUN` really reach `go test` — a green that proves nothing is
  the failure this repo knows best.

### 2. The budget — one place, and the doc cannot drift

`scripts/dep-budget.sh` holds the expected counts and the ceilings; the gate checks
`AGENTS.md`'s table against them. **Verify:**

- the live counts are computed **the way `AGENTS.md` says to compute them** — the
  external modules a package pulls in;
- the ceilings hold: `internal/core` **0**, `internal/cli` **0**, `internal/ui` **≤ 30**;
- **`AGENTS.md`'s table is checked**, so the doc cannot drift from the gate — including
  the *current* count it states;
- **there is ONE place for the numbers.** The runner says the script holds them and
  the doc is the checked copy. **Check it** — if the Makefile also carries a copy,
  that is a second answer;
- the failure message names **the live count, the expected count and the file to
  change**;
- and **a dependency that moves the count fails the gate.** Prove it: add a
  dependency to `internal/ui` temporarily, run `make budget`, and confirm it fails and
  says what to change. Revert it.

### 3. The tools, pinned, and CI

- **`make tools`** installs `staticcheck@v0.8.2` and `govulncheck@v1.9.0` — **pinned,
  not `@latest`** — and says where they landed and whether that directory is on the
  `PATH`.
- **CI** (`.github/workflows/ci.yml`): on **push** and **pull_request**;
  `permissions: contents: read`; the Go version is the **declared floor**
  (`go.mod` says `go 1.26.0`); the module cache; the tools installed and **on the
  `PATH`**; and the gate run **through `make`**, so the list of checks lives in one
  place and not in the workflow.
- **The pre-push hook** is committed, `make install-hooks` copies it, and it runs the
  **fast** gate — the slow one stays in CI, because it takes minutes (the runner
  measured `make check-race` at **3m06s**). Check the hook runs the fast gate and
  that it is executable.

### 4. What the slice found, and fixed

Three things happened here that are **not** tooling, and each needs its own check:

1. **staticcheck reported 7 issues** on today's tree, and the runner fixed them. The
   brief said *no behaviour changes*. **Read the diff and judge each fix**: one of them
   **removed a function** (`Operations.journal`, apparently dead). Confirm it really
   was dead — nothing calls it, and nothing outside a test called it either. A removed
   function that something used is a behaviour change.
2. **`go.mod` was untidy** and is now tidy. `make tidy-check` must produce no diff.
3. **govulncheck found 2 reachable vulnerabilities** (`golang.org/x/text`,
   `github.com/yuin/goldmark`), fixed by a version bump. **Check**: `make vuln` reports
   **0 reachable** vulnerabilities; the **module count did not move** (28 — a version
   bump is not a new module); and **the goldens did not move**.

### 5. The paper

**A tool earns its place with a written reason**, the way a dependency does. Check
that each of these is *there*, and that the reason is the real one:

| | |
|---|---|
| `AGENTS.md`'s **"The quality gate"** section | every target, and **why each tool** — `staticcheck`, `govulncheck`, `-race`, `-shuffle=on`, `-count=2`, `go mod tidy`, CI |
| `AGENTS.md`'s **working rules** | *the gate is `make check`; a focused check is `make test PKG=… RUN=…`; the tools are not run directly* |
| `AGENTS.md`'s dependency table | the **checked** copy, with the ceiling and the current count |
| `docs/testing.md` | a section on the gate, the tools and CI — it already said *"No pty test in CI"*, so that is where the CI story belongs |
| The paper says plainly | these are **dev tools**, not module dependencies: `go.mod` does not move, and the budget is untouched by them |

### 6. Scope

**No behaviour changes** — except the staticcheck fixes, which are checked above.
`core.DomainVersion` is **7** and did not move. The gate table, the plan format, the
command count (**42**) and the names did not move. `docs/audit.md` is the reviewer's.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.**

### The ledger

`docs/audit.md` is **yours** to write. The runner did not touch it — confirm that with
`git diff`. Add any NEW blocker as a new row with an id continuing from **F-97**. Do
not re-file what is closed; the ledger is not swept. **Add the row for slice H** — the
gate, the tools, and the reasons.

### Hard rules

- Do **not** edit implementation code. Read-only for everything except
  `docs/audit.md`.
- Revert every mutation you make — **including the temporary dependency** you add to
  prove the budget fails. The tree must be exactly as the runner left it, except your
  ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-h.md <<EOF
Reviewer: slice H review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-98.. or none>). Mutation checks: <what you mutated and reverted>. The gate without a tool: <what it did>. The budget with a dependency added: <what it said>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-h"
herdr agent prompt "w1Z:p1" "Reviewer: slice H review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-98.. or none>). Mutation checks: <what you mutated and reverted>. The gate without a tool: <what it did>. The budget with a dependency added: <what it said>. LOW notes: <any>."
```
