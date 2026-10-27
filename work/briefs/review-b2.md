You are the **reviewer** — the same lane that reviewed slice B. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 2 — the fix to slice B

The runner fixed the one finding you filed. **Cycle 2 is blockers only**: file **HIGH
or MEDIUM** blockers. No nits, no cosmetics, no style.

**Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's) or **F-73** (the plan's
own name — the orchestrator assigned it to slice C, and it stays open).

### 1. Is F-72 actually closed?

**F-72** — `docs/03-domain-model.md` § *The CPD ↔ tool mapping* said *"the artifact
**and the roles** get themed names"* while the table right below kept `runner ·
reviewer · orchestrator` in the **Kept plain** column.

The fix rewrites the rule sentence so it agrees with the table and with the decision
of record. Verify **the sentence and the table say the same thing** — read them
together, and read `docs/07-vocabulary.md` § *The rule* beside them, since the
decision of record is what both must match.

A sentence that no longer contradicts the table but now says something the decision of
record does not say is **not closed**. Judge the sentence against the decision, not
against the table alone.

Then confirm the runner did what the brief said and nothing more:

- **No code changed.** The finding was a sentence in a document. `git diff` over the
  fix's commit must touch `docs/03-domain-model.md` and nothing in `internal/`.
- **No new tripwire was added.** A wire that parses prose to ask which column a role
  sits in is brittle; the deny-list wire already covers the rename. A new wire here is
  scope creep, not a fix.
- **`docs/decisions/0005-vocabulary.md` was not edited.** An ADR is a dated record;
  slice E writes ADR 0009 to supersede it.

```bash
go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v   # the wires still hold
make lint && gofmt -l . && make test
```

### 2. Two ledger rows point at files that no longer exist

The finding is yours, so the pointer is yours to fix.

| Row | Says | Now |
|---|---|---|
| **T4** | `internal/ui/chop_test.go` — `core.TreeOf`; the console's Board tab | `internal/ui/board_test.go` |
| **T6** | `internal/ui/expo_test.go` — the Sous-Chef tab and the automation's vital signs | `internal/ui/souschef_test.go` |

A closed row claims *a test holds the fix*. A pointer at a file that was renamed
makes that claim false — the test still holds, but the row points at a ghost.

**Update the file pointers only. The finding's wording stays as it was filed** — the
ledger is the record of what was found, and it is not swept. This is a pointer that
follows a rename, not a rewrite.

### 3. Close F-72

When it is closed, move its row from open to closed, with the file and the test that
now holds it. Commit the ledger.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.**

### Hard rules

- Do **not** edit implementation code, and do **not** edit `docs/03` — the runner
  owns the fix. Read-only for everything except `docs/audit.md`.
- Revert every mutation you make. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-b-cycle-2.md <<EOF
Reviewer: slice B cycle 2 done. Verdict: <ok|changes-needed>. F-72 closed: <yes|no>. Rows T4/T6 pointers updated. Filed <N> new blockers (<F-74.. or none>). LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-b cycle 2"
herdr agent prompt "w1Z:p1" "Reviewer: slice B cycle 2 done. Verdict: <ok|changes-needed>. F-72 closed: <yes|no>. Rows T4/T6 pointers updated. Filed <N> new blockers (<F-74.. or none>). LOW notes: <any>."
```
