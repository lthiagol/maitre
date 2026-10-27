You are the **reviewer** — the independent reviewer lane, fresh for slice D. The
orchestrator's pane id is `w1Z:p1` — it is already substituted in this file; wherever
this brief says `$ORCH_PANE`, it means `w1Z:p1`. Use it verbatim in the notify command
below. The repository is `/Users/thiago/Code/Pessoal/maitre`.

## Your slice — the review of slice D

Slice D crossed the two words: the **per-criterion record** was `Proof` and is now
**evidence**; the **recipe-level stamp** was `Evidence` and is now **completion**. The
chef decided it, against advice: the journal's event kind was already `evidence`, so
the code now agrees with its own journal.

```bash
cd /Users/thiago/Code/Pessoal/maitre
git log --oneline -12
git diff <base>..HEAD --stat      # the base is the commit before slice D's first commit
```

**Cycle 1: blockers only.** File **HIGH or MEDIUM** blockers. No nits, no cosmetics,
no style, no LOW notes filed as rows. **Do not re-file F-56, F-59, F-60, F-61, P-8**
(the chef's).

### Verify each acceptance criterion with a command or a mutation

Not by reading the diff and believing it. Mutate, watch it fail, revert, and prove the
tree is clean afterwards.

**AC-01 — the per-criterion record is evidence.**

- `core.Proof` is `core.Evidence`; `Criterion.Evidence []Evidence` is the JSON
  `evidence`; the derived view is `evidence_summary`, computed by `EvidenceSummaryOf`.
- **`docs/03` § 8's table is tripwire-checked in both directions** — a view that is
  registered but not named fails the suite, and so does a view that is named but not
  registered. Confirm the table and the code agree, **both ways**, and that the wire
  can fail.
- The verifier, the gates, the operations that record a verification, the CLI's JSON
  and the views all follow.
- **The records stay append-only**: evidence is added, never rewritten, and a
  criterion's status still derives from its latest record.

**AC-02 — the recipe-level stamp is completion.**

- `core.Completion`, JSON `completion`, `r.Completion.CompletedAt`, set on `finish`
  exactly where it was. Nothing else about the recipe's life changed.

**AC-03 — the vocabulary in the docs and the glossary.**

- `docs/03`: the anatomy, the plan format and § 8's table say evidence. `docs/04`:
  the verbs that take a record say evidence.
- `GLOSSARY.md`: the **Proof** entry is **Evidence** — *"a verification result a
  machine produced — the command, its exit code, what it matched, and when. An
  evidence that matched nothing is not evidence"* — its `_Avoid_` line drops
  **Evidence** and gains **Proof**, and the CPD mapping row follows. The
  recipe-level stamp is described there too, so the glossary leaves no word
  unexplained.

**AC-04 — the tripwires still hold.**

```bash
go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v
```

The derived-views table, the gate table, the plan-format fields and the command count
(**42**) all agree with the code.

### The prose that kept the word "proof" — judge it

The runner kept the word in six places, and said why. **Judge each one.** The
distinction that matters: the word *proof* may name the **idea** of a verification
(what a gate checks) and not the **record**.

| Kept | Why the runner gave |
|---|---|
| `ac.all_proven` | a gate's name is a *rule*; renaming it is a rule change, not a rename |
| `proven` · `rejected` · `unknown` | they name the **state a criterion is in**, which the record derives — not the record |
| *"a criterion is **proved** by an existing suite"* (`docs/06`) | the verb is the act, not the record |
| **Proof** in the glossary's `_Avoid_` line | the line warns the reader off the word the code dropped |
| `proofs` / `evidence` as **old field names** | they name the format *before* the rename, which is the point |

**And judge the one place that could not keep the old wording:** the sign-off was
*"a sign-off is not a proof"*, and under the new words that would have been **false** —
a sign-off is stored as an evidence record (`signed-off`, claiming no match). It now
says *"a sign-off is not **machine** evidence"*. **Check that the new sentence is
true**, and that the exception it protects is still protected: the sign-off still
satisfies `ac.all_proven` explicitly, and it is still the one thing the zero-match
rule does not apply to.

**Also:** is there any place where the word *evidence* now means **two** things — the
per-criterion record and the recipe-level stamp? The stamp is **completion** now, so
the two must not collide.

### Then: what did the slice introduce?

1. **The reader tolerates the old field names.** A plan written before the rename
   keeps `proofs` (for the record) and `evidence` (for the stamp), and the reader
   reads them. A file carrying **both** names is a broken plan, and the **new** name
   wins — the reader must not merge two answers to one question. Prove both: the old
   name is read, and the new name wins.
2. **The version moved in AC-02 and the fingerprint moved again in AC-01.** The slice
   changes the format twice and `core.DomainVersion` is one number, so it went 5 → 6
   in the commit that changed the format first, and the recorded fingerprint moved
   again in AC-01. **Check the fence still stands in its designed direction**: a plan
   stamped **7** must be refused by this version-6 tool; a plan stamped **4**, **5** or
   **6** must be accepted.
3. **The journal.** `kitchen/activity.json` — the two migration events (35 and 36)
   are **appends**; events 1–34 are untouched. Check with `git diff` on that file, not
   by reading the prose. And the repo's own plan was migrated by hand — the recipe
   file carries the new names, and the journal says so.
4. **The sweep's reach.** The runner's diff touches `docs/01`, `docs/02`, `docs/05`,
   `docs/06`, ADR-0002, ADR-0006, ADR-0008 and `AGENTS.md`. Check each: a doc that
   named a view that no longer exists (`ProofSummary`) had to move, and `docs/02`'s
   principle label had to move because the word would otherwise mean two things.
   **Read the whole diff** — did a change touch something it did not name?
5. **Scope.** The slice is the two words. Confirm it did **not** touch the roles, the
   tabs, the backlog, the ids, the module path or the binaries.

### The invariants — check them, do not assume them

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps **28**
  external modules.
- The domain layer has **zero UI imports**.
- **A criterion's record is append-only.**
- **The journal is append-only** — the historical events keep their old ids and their
  old field names.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.** They go in the
  ledger in a later round.

### The ledger

`docs/audit.md` is **yours** to write. The runner did not touch it — confirm that with
`git diff`. Add any NEW blocker as a new row with an id continuing from **F-77**. Do
not re-file what is closed; the ledger is not swept. Commit the ledger.

### Hard rules

- Do **not** edit implementation code. Read-only for everything except
  `docs/audit.md`.
- Revert every mutation you make. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-d.md <<EOF
Reviewer: slice D review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-78.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-d"
herdr agent prompt "w1Z:p1" "Reviewer: slice D review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-78.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>."
```
