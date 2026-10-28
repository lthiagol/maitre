You are the **reviewer** — the same lane that reviewed slice F. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 3 — the last fixes to slice F

The runner fixed F-91 and the notes. **Cycle 3 is blockers only**: file **HIGH or
MEDIUM** blockers. No nits, no cosmetics, no style.

**Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's), or your own LOW notes
from cycle 2 (the cross-kind file name — logged for the ledger's next round, with the
orchestrator's agreement).

### 1. Is F-91 closed?

> The rule moved to where the name enters the record: `Store.SaveRecipe` runs
> `checkRecordNames` on every record, the transition's context is checked in
> `Operations.advance`, the audit role is checked at the write's entry, a new journal
> event may not carry a retired name, and `normalizeAgentName` maps a space and an
> underscore to a dash.

**Prove each of these, by command:**

- **The domain's API cannot write a plan it will not read.** `ops.AdvanceWith(id,
  core.EventApprove, core.TransitionContext{Approver: "line-cook", Actor: "chef"})`
  must **fail** and **write nothing** — the probe you filed.
- **The claim, the review, the approval, the sign-off, the comment and the actor** are
  each refused for `line-cook`.
- **The audit role** is refused: `--role line-cook` must not reach the journal.
- **A new journal event may not carry a retired name** — and the **history** keeps
  its old names. Check both, and check the *history* is not what stopped it.
- **`line_cook` is refused** — the same word with an underscore. And `Line Cook`,
  `LINE-COOK`, ` line-cook ` — the same name however it is spelled.

**And the one place, not many.** The rule must live where a name **enters a record**,
not once per operation. Read the diff and judge: is there **one** funnel, or several
checks that happen to cover today's operations? **A check per operation is a rule
waiting for the operation nobody wrote yet.**

### 2. The notes that rode with it

- **The wire now reads the code's canonical lists** (`RetiredNames`,
  `retiredDirectories`, the two field constants, `retiredIDForms`,
  `retiredFamilies`) instead of carrying its own copy. **Check your four probes still
  fail it** — `specials` in `internal/ui`, `specials` in `model.go`, `R01` in
  `store.go`, and `LegacyThing` — and that a retired name added to the code's list is
  caught **without** touching the wire.
- **`docs/03` § 6.2 and § 2.1 say why the version moved**: the facts list grew **and**
  version 6 covered two digest encodings. **Check the sentence is now true in both
  directions** — a fact *added* to the list that the format already had moves the
  fingerprint and not the version (slice E's case), and a fact the format *changed* is
  a version bump.
- `docs/07-vocabulary.md` was the orchestrator's and is corrected in the same spirit
  (`c429bd4`). **Check it agrees with the code and with `docs/03`.**

### 3. The two tests that changed their premise — again

The runner flagged them in the receipt. **Judge the premises, not the words:** does
each test still assert the property it was written for? A renamed test that quietly
asserts something *else* is a lost property wearing a new name.

### 4. Scope, the tree and the ledger

- `core.DomainVersion` is **7** and did not move again. The gate table, the plan
  format and the command count (**42**) must not have moved. The **role names** must
  not have moved.
- **The mascot** — another session's work, in `logo/`. Check the runner committed
  nothing of it.
- **The ledger** — `docs/audit.md` is yours. Move **F-91** to closed, with the file
  and the evidence that holds it, and **add the row for slice F itself**: the ruling
  that this is not brownfield, the tolerance machinery that came out, the one-format
  wire, and the one place where a name enters a record. The ledger is not swept: keep
  the findings' wording as they were filed.

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
cat > work/inbox/reviewer-slice-f-cycle-3.md <<EOF
Reviewer: slice F cycle 3 done. Verdict: <ok|changes-needed>. F-91 closed: <yes|no>. Filed <N> new blockers (<F-92.. or none>). Mutation checks: <what you mutated and reverted>. The two renamed tests: <what each asserts now>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-f cycle 3"
herdr agent prompt "w1Z:p1" "Reviewer: slice F cycle 3 done. Verdict: <ok|changes-needed>. F-91 closed: <yes|no>. Filed <N> new blockers (<F-92.. or none>). Mutation checks: <what you mutated and reverted>. The two renamed tests: <what each asserts now>. LOW notes: <any>."
```
