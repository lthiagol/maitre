You are the **reviewer** — the same lane that reviewed slice H. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 4 — the last round of slice H

The runner built the orphan check and scoped the doc check. **Cycle 4 is blockers
only**: file **HIGH or MEDIUM** blockers. No nits, no cosmetics, no style.

**This is the last cycle for slice H.** Anything still open after it is a note for the
ledger, not another round — the repo's discipline is that a LOW is *noted*, not
chased, and this slice has had three rounds already.

### 1. The orphan check — derived, not declared

> `check-golden-claims.sh mark` moves every `.golden` timestamp into the past, the
> golden run rewrites the ones a test claims, `report` names the ones nobody wrote.
> No manifest, and no scan for a path in a test file, because the paths live in
> constants and a scan would be a guess.

**The mechanism, and why it is sound:** the golden test writes the file
**unconditionally** when `-update` is set — `internal/ui/golden_test.go` calls
`os.WriteFile` before it compares — so a **claimed** golden is always rewritten, and
its timestamp always moves. A file nobody claims keeps its marked timestamp and is
named.

**Verify, by command:**

- **a golden test renamed away** must make `make goldens-check` **fail**, naming the
  file (`internal/ui/testdata/console-recipes.golden`) — and the message must say
  **both** possible causes, because the second is a real hole: *the filter misses the
  test that claims it*, in which case `make goldens` would not regenerate it either;
- **a claimed golden whose content already matches** must **not** be named — check
  that the golden test really does write unconditionally, so the mechanism cannot cry
  wolf on a golden that is simply up to date;
- and the check must leave the **tree clean** — the marked timestamps must not
  survive as a modification.

**And the rule the orchestrator set:** *a check that cries wolf is worse than no
check*. Confirm this one does not.

### 2. The doc check compares the pinned version only

> a line with the release version's shape (`staticcheck 2026.2.2`) does not fail it,
> and a second mention with the pinned shape and the wrong value (`v0.8.1`) still
> does.

**Prove both:** the release-version shape must pass; the pinned shape with the wrong
value must fail; and the **installed** tool's version is still checked against the pin.

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
- `make vuln` reports **0 reachable**; the **goldens did not move**; `tidy-check` puts
  the tree back byte-for-byte when tidy fails.

### 4. Scope, and the ledger

- **No behaviour changes** in the code — the cycle touched no Go source.
  `core.DomainVersion` is **7**; the gate table, the plan format, the command count
  (**42**) and the names did not move.
- **The ledger** — `docs/audit.md` is yours. **F-98 closed**, the **slice H** row
  there, and the cycle rows. The three new tripwires the ledger's table gained — the
  dependency budget, the tool pins, and now the goldens' claims — should be in the
  discipline table too, so a reader of the ledger sees what the gate holds.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.** They go in the
  ledger for a later round — **not** another cycle.

### Hard rules

- Do **not** edit implementation code — except `docs/audit.md`, which is yours.
- Revert every mutation you make — **including** the renamed golden test and the
  marked timestamps. The tree must be exactly as the runner left it, except your
  ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-h-cycle-4.md <<EOF
Reviewer: slice H cycle 4 done. Verdict: <ok|changes-needed>. F-98 closed: <yes|no>. Filed <N> new blockers (<F-99.. or none>). Mutation checks: <what you mutated and reverted>. The orphan check: <what it did>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-h cycle 4"
herdr agent prompt "w1Z:p1" "Reviewer: slice H cycle 4 done. Verdict: <ok|changes-needed>. F-98 closed: <yes|no>. Filed <N> new blockers (<F-99.. or none>). Mutation checks: <what you mutated and reverted>. The orphan check: <what it did>. LOW notes: <any>."
```
