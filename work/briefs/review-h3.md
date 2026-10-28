You are the **reviewer** — the same lane that reviewed slice H. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 3 — the last fix of slice H

The runner fixed the goldens' invocation and the three one-liners. **Cycle 3 is
blockers only**: file **HIGH or MEDIUM** blockers.

**Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's), or your own LOW notes
from cycle 2 — except the ones the orchestrator ruled on, which are listed below.

### 1. The goldens are a test invocation, and the gate is one script

> The ruling was *one place decides how a test run is invoked*. The goldens are a
> third invocation, and it bypassed the zero-match refusal — so a `goldens-check` that
> regenerated nothing passed with an empty diff.

**Verify, by command:**

- **the golden tests renamed away** (temporarily) must make `make goldens-check`
  **fail** — it passed before;
- `goldens` and `goldens-check` are **two targets of one recipe** — read it and judge:
  two spellings of the invocation are two answers to one question, and they will drift;
- and the goldens go through **`scripts/go-test.sh`**, the same script as `test` and
  `test-race`, so the refusal applies to them too.

### 2. The three one-liners

| Note | What to check |
|---|---|
| `AGENTS.md`'s `make check` row omits `versions-check` | the doc now names it. Check it. |
| the version failure said *"run `make tools`"* but not how to fix a PATH order | it must now say **both** — and the reason the second half matters is that an older tool *earlier* on the `PATH` still wins after `make tools`. **Read the message the user sees** and judge whether someone who follows it would be unblocked. |
| the doc check looked only at the **first** matching row | it must read **every** mention. **Prove it:** a second mention of the same tool, with another version, must fail the check. |

### 3. The note that is accepted as it stands — do not "fix" it

> `check-tool-versions.sh` calls `require-tool.sh` by a relative path and works only
> from the repo root.

**The orchestrator accepted this, and the reason is the chef's ruling: the interface
is `make`.** The scripts are `make`'s implementation, and `make` runs its recipes from
the repo root. Every script in `scripts/` behaves this way. **Confirm it is consistent
across the scripts** — and do not propose a change.

### 4. The gate, whole

- **Every target** runs alone: `fmt` · `fmt-check` · `vet` · `staticcheck` · `lint` ·
  `test` · `test-race` · `goldens` · `goldens-check` · `budget` · `tidy-check` ·
  `versions-check` · `check` · `check-race` · `tools` · `vuln` · `ci` · `install-hooks`.
- `make check`, `make -j4 check` and `make ci` are green.
- The gate **does not skip a check**: with a tool's directory off the `PATH`, it fails
  with the instruction.
- The budget is **0 / 0 / 28** against **0 / 0 / 30**, and `AGENTS.md`'s table is
  checked **both ways** — doctor the doc's cell, and doctor the script's ceiling.
- The tool versions are checked **both ways** — the doc's mention, and the version of
  the tool the gate runs.
- `make vuln` reports **0 reachable**; the **goldens did not move**; `make tidy-check`
  puts the tree back when tidy fails.

### 5. Scope, and the ledger

- **No behaviour changes** in the code. `core.DomainVersion` is **7**; the gate table,
  the plan format, the command count (**42**) and the names did not move.
- **The ledger** — `docs/audit.md` is yours. **F-98** must be **closed** (it was
  filed in cycle 1 and the runner fixed it), and the row for **slice H** must be there:
  the gate, the tools, the reasons, and the two tripwires the ledger's table gained —
  the dependency budget and the tool pins.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.**

### Hard rules

- Do **not** edit implementation code — except `docs/audit.md`, which is yours.
- Revert every mutation you make — **including** the renamed golden tests, the
  doctored doc and the doctored ceiling. The tree must be exactly as the runner left
  it, except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-h-cycle-3.md <<EOF
Reviewer: slice H cycle 3 done. Verdict: <ok|changes-needed>. F-98 closed: <yes|no>. Filed <N> new blockers (<F-99.. or none>). Mutation checks: <what you mutated and reverted>. The goldens with their tests renamed away: <what it did>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-h cycle 3"
herdr agent prompt "w1Z:p1" "Reviewer: slice H cycle 3 done. Verdict: <ok|changes-needed>. F-98 closed: <yes|no>. Filed <N> new blockers (<F-99.. or none>). Mutation checks: <what you mutated and reverted>. The goldens with their tests renamed away: <what it did>. LOW notes: <any>."
```
