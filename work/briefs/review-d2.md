You are the **reviewer** — the same lane that reviewed slice D. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 2 — the fixes to slice D

The runner fixed the three findings you filed. **Cycle 2 is blockers only**: file
**HIGH or MEDIUM** blockers. No nits, no cosmetics, no style.

**Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's), or your own LOW notes
from cycle 1 (the grammar of *"an evidence"* — folded into F-79; `docs/07`'s
alternatives table — the orchestrator ruled it stays; `research/` — history; the
glossary's brigade section — slice E's; `scripts/stamp-plan`'s usage line — fixed
here; `AGENTS.md` and `work/README.md` — the orchestrator's, and correct; and § 6's
"does not look" sentence — folded into F-80).

### 1. Are the three findings actually closed?

Verify each with a **command or a mutation**, not by reading the diff and believing it.

**F-78 (high) — the approval digest travels with the rename.**

- `legacyFingerprints` now enumerates the **field name** as well as the list forms.
  **Check the set is still COMPLETE.** The comment above the function used to argue
  why the set was complete; that argument had to be rewritten, not extended in
  silence. Read the new argument and check it: the candidates must be the **cross
  product** of every encoding an older build could have written — the old field name
  *and* the `null`/`[]` list forms — and nothing else.
- **The safety property is the one that matters now.** The candidate set grew, and
  every candidate is another way a digest could accidentally match. **Prove that a
  recipe whose content changed after approval still matches NO candidate**, and that
  the `immutability.respected` gate still refuses it. A gate that can be fooled by a
  wider candidate set is worse than the bug you filed.
- **The repository's own plan is migrated** and the standing discipline is a test
  now — *the repository's own plan's stored approval digest must still match the
  canonical recomputation*. Check that the test reads the digest **off the disk** and
  compares it with the recomputation, and that it would fail if the file were wrong.
- **The journal** took the migration as an append (event 37). Check with `git diff`
  on `kitchen/activity.json`: events 1–36 untouched.

**F-79 (medium) — the glossary's entry covers the record it names.**

- The **Evidence** entry names the sign-off exception, and `docs/03` § 5 says the same
  thing. **Check the two say the same thing** — the single source and the spec must
  not contradict each other, which is what the finding was.
- **Check the sentence is true**: a result that matched nothing does not count, *and
  the chef's sign-off is the exception, and it claims no match*. The code stores a
  sign-off as `Criterion.Evidence{Match: "signed-off", Matched: 0, SignedOff: true}`
  and `ac.all_proven` takes it as a pass.
- *"An evidence"* is gone from the tree outside `docs/audit.md`. Check that the word
  was not simply **banished** where it was right: `ac.all_proven`, the match states
  (`proven` · `rejected` · `unknown`), the old field names, and prose that means the
  *idea* of a proof may keep it.

**F-80 (medium) — the version paragraph stops contradicting the reader.**

- `docs/03` § 6 and `internal/core/lifecycle.go`'s version notes now say what the
  reader **does**: it reads the old field name, it shows the old specials read-only
  in their own group, and it never looks for new ones. **Check each claim against the
  code**, and that the version bump is unchanged (the *writer* changed — that is what
  the version fences).
- `scripts/stamp-plan`'s usage line names the field-name migrations.

### 2. The ADR line — judge it against the rule, not against taste

The runner's F-79 commit renamed the vocabulary **incidentally** in ADR-0002,
ADR-0006, ADR-0008 and the ADR index, and left **ADR-0004 (naming)** and **ADR-0005
(vocabulary)** untouched.

**The orchestrator's ruling, which you judge against:**

| | |
|---|---|
| **An ADR whose subject is the mechanism** (0002, 0006, 0008) | its *incidental* vocabulary may move — it decides about verification or about the interface, not about the word, and a doc that contradicts the code is a bug |
| **An ADR whose subject is the vocabulary** (0004, 0005) | it **keeps its tables as the record of what it said**, and ADR 0009 supersedes it in slice E. Those are the two the runner left alone. |

**Judge two things:** that the runner drew this line correctly, and that in the ADRs
it did touch, only the *vocabulary* moved — **not the substance**. An ADR that
decides something different now, or loses a constraint, is a rewrite of the record
and a blocker. Read those four ADRs' diffs closely.

### 3. What did the fixes introduce?

1. **Scope.** The gate table, the plan format, the command count (**42**) and
   `core.DomainVersion` (**6**) must not have moved. Check the whole diff of the
   three commits — the F-79 commit touched 21 files, and each one needs a reason.
2. **The shared tree.** Another session is drawing the mascot in this repository and
   has files in the tree (`logo/`). **Confirm the runner committed nothing of it** —
   no file of the other session's, and no file the runner did not write. The
   runner's report says it staged every commit by name.
3. **The plan's writes.** The only write to the plan should be the migration event.
   Check that the journal took **exactly one** new event and that nothing else in
   `kitchen/` moved except the recipe's own migrated field.

### 4. The ledger

Move **F-78, F-79, F-80** from open to closed, with the file and the test that now
holds each. Commit the ledger. The ledger is not swept: keep the findings' wording
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
- **Do not touch `logo/`** — it is another session's work.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-d-cycle-2.md <<EOF
Reviewer: slice D cycle 2 done. Verdict: <ok|changes-needed>. F-78..F-80 closed: <yes|no>. Filed <N> new blockers (<F-81.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-d cycle 2"
herdr agent prompt "w1Z:p1" "Reviewer: slice D cycle 2 done. Verdict: <ok|changes-needed>. F-78..F-80 closed: <yes|no>. Filed <N> new blockers (<F-81.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>."
```
