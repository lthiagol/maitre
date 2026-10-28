You are the **reviewer** — the same lane that reviewed slice E. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 3 — the last fixes to slice E

The runner fixed F-86 and the two notes you left. **Cycle 3 is blockers only**: file
**HIGH or MEDIUM** blockers. No nits, no cosmetics, no style.

### F-86 — the instrument was blind, not the tree

**Check the instrument, because the tree did not change.**

- The sweep's command now carries **all seven** retired names — *maître d'* among
  them — and it **excludes the glossary's `_Avoid_` lines**, so a warning about a
  retired word is not counted as drift.
- **Run the command yourself** and check the output is the same the commit body
  shows. The command is in the commit body and in `work/inbox/runner-slice-e-fix-2.md`.
- **Check the count reconciles.** You counted **77 lines in 19 files**. The runner
  reports **72 in 18**, and explains the movement: the old command counted **five of
  the glossary's six `_Avoid_` lines** as drift, and the complete command sees all six
  and excludes all six. **Do the arithmetic**: 77 − 5 = 72, and one file drops out
  (the glossary's own). Confirm that the line the new command finds where the old
  could not — *maître d'* — is in a place that is legitimately excluded.
- **Check the drift that remains is all legitimate**, file by file, as in cycle 2.

### The two notes

- **The glossary's heading is now "### The surfaces"**, and the section says *"The
  console is a **surface**, not a job"* — ADR-0003's word. Check the taxonomy reads:
  `grep -n "^### " GLOSSARY.md` → **The plan**, **The crew**, **The work**, **The
  surfaces**. The console must not be filed as a job, and must not be filed as a
  string.
- **`docs/03` § 6's pointer now lands** on a new `### 6.2 The version guard`, and
  the implementation is named in the same paragraph (`core.DomainVersion`, and
  `Store.checkDomainVersion` in `internal/core/lock.go`). **Follow the link** — it
  must resolve, and the section it lands on must say what the code says.

### The mascot and the tree

**The mascot did not land** — the art is in `logo/` (untracked), and no placement
carries a drawing. **AC-01 moves to the drawing session.** Check that nothing of the
other session's work was committed, and that `kitchen/activity.json` did not move.

### Scope

F-86's commit **changed no file** — the fix was the instrument, in the commit body
and the brief. Check that is true, and that the two notes' commits touched only the
glossary's heading and `docs/03`.

### The ledger — and this is the last round

Move **F-86** from open to closed, with the evidence that now holds it. **Add the
row for slice E itself**: the rename's last slice — the glossary rewritten, ADR 0009
written, the tree swept, and **the mascot's acceptance criterion carried to the
session that owns the drawing**. Commit the ledger.

The ledger is not swept: keep the findings' wording as they were filed.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.**

### Hard rules

- Do **not** edit implementation code. Read-only for everything except
  `docs/audit.md`.
- Revert every mutation you make. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
- **Do not touch `logo/`** — it is another session's work.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-e-cycle-3.md <<EOF
Reviewer: slice E cycle 3 done. Verdict: <ok|changes-needed>. F-86 closed: <yes|no>. Filed <N> new blockers (<F-87.. or none>). Sweep count: <the number, and whether it reconciles>. Mascot: <did it land? yes/no>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-e cycle 3"
herdr agent prompt "w1Z:p1" "Reviewer: slice E cycle 3 done. Verdict: <ok|changes-needed>. F-86 closed: <yes|no>. Filed <N> new blockers (<F-87.. or none>). Sweep count: <the number, and whether it reconciles>. Mascot: <did it land? yes/no>. LOW notes: <any>."
```
