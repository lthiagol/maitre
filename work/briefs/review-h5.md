You are the **reviewer** — the same lane that reviewed slice H. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 5 — the last cycle of slice H

The runner fixed the last two notes. **Cycle 5 is blockers only**: file **HIGH or
MEDIUM** blockers. No nits, no cosmetics, no style.

**This is the last cycle for slice H.** Anything still open after it is a note for the
ledger — **not** another round.

### 1. The gate sees a golden nobody committed

> `git diff --exit-code` does not see untracked files, so a new golden that a test
> writes but nobody commits passed the gate.

**Verify, by command:**

- **a new golden no test claims** → `make goldens-check` **fails**, names the file,
  and the message names the **remedy** for that cause (name the test so the filter
  matches it, or extend the filter) — and says which cause it is;
- **a new golden a test claims but nobody committed** → **fails**, naming the file;
- **a golden the commit would drop** (staged for deletion) → **fails** the same way;
- and the check must **leave the tree exactly as it was** — no intent-to-add left
  behind, and a staged deletion still staged (the runner says the index is put back
  with `git rm --cached --force` rather than `git reset`, so a staged deletion stays
  staged — **check that**).

### 2. The messages say what to do

The rule the orchestrator has applied to this gate twice now: **every message tells
the reader what to do.** Check the orphan message and the commit message both name a
**remedy**, and that each names the **right** remedy for the cause it reports.

**And check they cannot cry wolf:** the reviewer should be able to follow each message
and be unblocked, without reading the script.

### 3. The gate, whole

- **Every target** runs alone: `fmt` · `fmt-check` · `vet` · `staticcheck` · `lint` ·
  `test` · `test-race` · `goldens` · `goldens-check` · `budget` · `tidy-check` ·
  `versions-check` · `check` · `check-race` · `tools` · `vuln` · `ci` · `install-hooks`.
- `make check`, `make -j4 check` and `make ci` are green.
- The gate **does not skip a check**: with a tool's directory off the `PATH`, it fails
  with the instruction.
- The budget is **0 / 0 / 28** against **0 / 0 / 30**, checked **both ways**.
- The tool versions are checked **both ways** — the doc's mention and the tool the
  gate runs.
- The goldens are checked three ways: **unchanged**, **claimed by a test**, and **in
  the commit**.
- `make vuln` reports **0 reachable**; `tidy-check` puts the tree back byte-for-byte
  when tidy fails.

### 4. Scope, and the ledger

- **No behaviour changes** in the code. `core.DomainVersion` is **7**; the gate table,
  the plan format, the command count (**42**) and the names did not move.
- **The ledger** — `docs/audit.md` is yours. **F-98 closed**, the **slice H** row, the
  cycle rows, and the discipline table's tripwires: the dependency budget, the tool
  pins, the goldens' claims — and now **the goldens' commit check**.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.** They go in the
  ledger for a later round — **not** another cycle.

### Hard rules

- Do **not** edit implementation code — except `docs/audit.md`, which is yours.
- Revert every mutation you make — **including** the new goldens, the staged
  deletion and the intent-to-add. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-h-cycle-5.md <<EOF
Reviewer: slice H cycle 5 done. Verdict: <ok|changes-needed>. F-98 closed: <yes|no>. Filed <N> new blockers (<F-99.. or none>). Mutation checks: <what you mutated and reverted>. The goldens' three checks: <what each did>. The messages: <what they tell the reader to do>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-h cycle 5"
herdr agent prompt "w1Z:p1" "Reviewer: slice H cycle 5 done. Verdict: <ok|changes-needed>. F-98 closed: <yes|no>. Filed <N> new blockers (<F-99.. or none>). Mutation checks: <what you mutated and reverted>. The goldens' three checks: <what each did>. The messages: <what they tell the reader to do>. LOW notes: <any>."
```
