You are the **reviewer** — the same lane that reviewed slice C. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 2 — the fixes to slice C

The runner fixed the four findings you filed. **Cycle 2 is blockers only**: file **HIGH
or MEDIUM** blockers. No nits, no cosmetics, no style.

**Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's), or your own LOW notes
from cycle 1 (the refusal that does not carry the reason; `LegacyItem.Text` returning
a prep note's intent under the `body` field; the glossary's CPD rows). **F-73** is
closed.

### 1. Are the four findings actually closed?

Verify each with a **command or a mutation**, not by reading the diff and believing it.

**F-74 (high) — the live console shows an old plan's specials.**

- The live path and the `--render` path now go through one `Model.applyReload`.
  **Check there is no second place that assigns the model's fields** — a field the
  render carries and the live path forgets is the same bug in a smaller coat.
- The test drives the **real** path: the helper sends the reload message through the
  model's own `Update`. **Confirm the helper no longer copies the fields by hand** —
  a helper that sets up the state it is testing is not a test.
- **Mutate:** remove `m.legacy = msg.legacy` from `applyReload` and confirm the tests
  fail. A test that cannot fail is not a test.

**F-75 (medium) — a dangling pointer reserves its id again.**

- `core.LegacyItem` carries `promoted_to` and the reader reads it; `reservedRecipeIDs`
  scans the legacy items' pointers.
- **Check the format is the real one.** The old prep note wrote
  `"promoted_to": "recipe RCP10"` — two words — and `referencedIDs` only parses that
  shape. A bare `"RCP10"` is one word and is **skipped**. Confirm the fix reserves the
  id in the format that is actually on disk, and that a pointer in any other spelling
  is either handled or deliberately not (and say which).
- **Mutate:** remove the legacy scan and confirm the test fails.

**F-76 (medium) — a list flag that is not given leaves its field alone.**

- Both list fields obey `docs/04-operations.md`: *a list flag replaces the list it
  names; a flag that is not given leaves its field alone*. **One rule, not two
  guards.**
- **Check the other verbs.** `listArg` returns nil for a flag that was not given.
  Confirm the verbs that were not supposed to change behave exactly as before —
  `--in-scope`, `--out-of-scope`, `--ac`, `--ac-verify` and the rest.
- **Check the JSON edge.** `docs/04-operations.md` also says *an empty one is `[]`,
  never `null`*. Confirm **no document** carries `null` for a list, and that a list
  field is not silently **omitted** where the rule says `[]`. Look at a written
  backlog item, not at the code alone.
- **Mutate:** restore the old guard and confirm the test fails.

**F-77 (medium) — the spec says what is true.**

- The row and the paragraph now agree, and the code did not move. Confirm the doc
  says a promoted recipe is **checkable, not yet falsifiable** — and that nothing in
  the code now claims otherwise.

### 2. What did the fixes introduce?

1. **`listArg` and the `given` map.** The parser now knows which flags were given.
   A flag that is not given and a flag given with no value must be different things
   — check that the parser still refuses `--tag` with no value (the chef's question is
   open; the runner did not change it, and must not have).
2. **`emptyIfNil` at the JSON edge.** It is a second place that knows the rule "a list
   is a list". Check that the rule is applied on **every** path that writes a list —
   create, update, and the store's own write — and that a list loaded from a
   hand-edited file (a missing field) cannot reach a document as `null`.
3. **`setList`.** One rule for both fields. Check there is no *other* place in the
   domain where a list field is guarded by `!= nil` in one spot and `len(...) > 0` in
   another — that divergence is what F-76 was.
4. **The legacy detail pane now shows "Promoted to …"** (the runner added three
   lines to `internal/ui/detail.go` under F-75). Judge it: the pointer is read now,
   so showing it is honest — but it is **beyond the finding**. Say whether it belongs
   in this slice or should come out.
5. **Scope.** Read the whole diff of the four commits. Did a fix touch anything the
   findings did not name? The gate table, the plan format and the command count (42)
   must not have moved.

### 3. The ledger

Move **F-74, F-75, F-76, F-77** from open to closed, with the file and the test that
now holds each. Commit the ledger. The ledger is not swept: keep the findings' wording
as it was filed.

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

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-c-cycle-2.md <<EOF
Reviewer: slice C cycle 2 done. Verdict: <ok|changes-needed>. F-74..F-77 closed: <yes|no>. Filed <N> new blockers (<F-78.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-c cycle 2"
herdr agent prompt "w1Z:p1" "Reviewer: slice C cycle 2 done. Verdict: <ok|changes-needed>. F-74..F-77 closed: <yes|no>. Filed <N> new blockers (<F-78.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>."
```
