You are the **reviewer** — the same lane that reviewed slice H. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 2 — the fixes to slice H

The runner fixed F-98 and the seven notes. **Cycle 2 is blockers only**: file **HIGH or
MEDIUM** blockers. No nits, no cosmetics, no style.

**Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's).

### 1. F-98 — the slow gate cannot pass having tested nothing

> `make test-race PKG=./internal/brand RUN='-run NoSuchTest'` exited 0, while
> `make test` with the same filter exits 2.

**Verify, by command:**

- `make test-race PKG=./internal/brand RUN='-run NoSuchTest'` must **fail** — with the
  zero-match message, not a bare error;
- `make test-race PKG=./internal/brand RUN='-run TestTheDrawingsMatchTheSource'` must
  pass **and actually run that test** — check the command that reached `go test`;
- **and the two targets must not be able to drift again.** The brief said *one place
  decides how a test run is invoked*. **Read the diff and judge**: is the invocation
  now written once, or are there two recipes that still spell it out separately? Two
  spellings are two answers to one question, and they will drift again.

### 2. The seven notes

| Note | What to check |
|---|---|
| the **installed** tool's version is never checked | the gate now has `make versions-check`: it checks the version of the tool **the gate runs** against the pin. **Prove it** — put an older `staticcheck` on the `PATH` and confirm the gate **fails** with the instruction, not silently runs it. |
| the tool versions are **written twice** and only one is used | the `Makefile` holds the pins; `docs/testing.md` is the **checked copy**. **Prove the doc cannot drift**: change the version in the doc and confirm the check fails. And check the *other* direction — that the check compares the **installed** tool too. |
| `dep-budget.sh` ignores `GO` | it honours `GO` now, like every other target. Check it. |
| `tidy-check` leaves the tree dirty if tidy fails | the tree must be **as it was**: snapshot `go.mod`/`go.sum`, restore them if tidy fails. **Prove it** — make tidy fail (no network, or a broken `go.mod`), and confirm the tree is clean afterwards. |
| `install-hooks` overwrites silently | it says what it replaced. Check it. |
| `make -j` is unsafe | `.NOTPARALLEL:` is set, with a comment saying why. **Prove it**: `make -j4 check` must not interleave two targets that touch the same files. |
| the ledger's discipline table still says to run `go test` by hand | **this one is yours.** `docs/audit.md` is the reviewer's file, and the runner was told not to touch it. Update it: the discipline table should name the make interface. |

### 3. The gate, whole

- **Every target** still exists and runs alone: `fmt` · `fmt-check` · `vet` ·
  `staticcheck` · `lint` · `test` · `test-race` · `goldens` · `goldens-check` ·
  `budget` · `tidy-check` · `versions-check` · `check` · `check-race` · `tools` ·
  `vuln` · `ci` · `install-hooks`.
- `make check` and `make ci` are green.
- **The gate does not skip a check:** with a tool's directory off the `PATH`, it
  **fails** with the instruction.
- The budget is **0 / 0 / 28** against the ceilings **0 / 0 / 30**, and `AGENTS.md`'s
  table is checked against it. **Prove both ways**: doctor the doc's cell, and doctor
  the script's ceiling — each must fail.
- `make vuln` reports **0 reachable** vulnerabilities.
- The **goldens did not move**.

### 4. Scope, and the paper

- **No behaviour changes** in the code — the cycle was the gate's own plumbing.
  `core.DomainVersion` is **7**; the gate table, the plan format, the command count
  (**42**) and the names did not move.
- The **paper**: `AGENTS.md`'s "The quality gate" section (every target, a written
  reason for each tool), the working rule, the checked dependency table, and
  `docs/testing.md`'s gate section — including that the doc's tool versions are now
  **checked** and that the gate verifies **the version it runs**.
- `docs/audit.md` is yours — the ledger, with the row for slice H.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.**

### Hard rules

- Do **not** edit implementation code — except `docs/audit.md`, which is yours.
- Revert every mutation you make — **including** the older tool on the `PATH`, the
  doctored doc and the failed tidy. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-h-cycle-2.md <<EOF
Reviewer: slice H cycle 2 done. Verdict: <ok|changes-needed>. F-98 closed: <yes|no>. Filed <N> new blockers (<F-99.. or none>). Mutation checks: <what you mutated and reverted>. The gate with an old tool: <what it did>. The ledger's table: <what you changed>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-h cycle 2"
herdr agent prompt "w1Z:p1" "Reviewer: slice H cycle 2 done. Verdict: <ok|changes-needed>. F-98 closed: <yes|no>. Filed <N> new blockers (<F-99.. or none>). Mutation checks: <what you mutated and reverted>. The gate with an old tool: <what it did>. The ledger's table: <what you changed>. LOW notes: <any>."
```
