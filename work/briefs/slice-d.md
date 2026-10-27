You are the **runner** for the vocabulary rename of this project. The orchestrator's pane id is `$ORCH_PANE` — the orchestrator sets it before dispatching. Read `docs/07-vocabulary.md` (the decision of record), `AGENTS.md` (house style) and `docs/03-domain-model.md` (the spec) before you touch anything.

## Your slice — D of 5 — evidence and completion

**Scope discipline — do only this.** Do NOT touch the roles, the tabs, the backlog, the ids, the module path or the binaries. Slices A–C are done; leave their files alone unless this slice forces a change (say so if it does).

Today the two words are **crossed**: `Proof` is the per-criterion record (the command, its exit code, what it matched, who signed it off), and `Evidence` is the recipe-level completion stamp (`completed_at`). Meanwhile the **journal's event kind is already `evidence`** — and it is the kind used when a proof is submitted. So the code says *proof* and the journal says *evidence* for the same thing.

The chef's decision: **Evidence wins.** The code ends up agreeing with its own journal.

### AC-01 — the per-criterion record is evidence

- `core.Proof` becomes `core.Evidence`. `Criterion.Proofs []Proof` (JSON `proofs`) becomes `Criterion.Evidence []Evidence` (JSON `evidence`).
- The verifier, the gates, the operations that record a proof, the CLI's JSON and the views follow.
- The derived view `ProofSummary` becomes `EvidenceSummary` — **and `docs/03` § 8's table of derived views is updated in the same commit.** That table is tripwire-checked in both directions: a view that is registered but not named, or named but not registered, fails the suite.
- The journal's event kind stays `evidence`. It already names the right thing.

### AC-02 — the recipe-level stamp is completion

- The recipe-level `Evidence` struct (one field: `completed_at`) becomes `core.Completion` (JSON `completion`), on the recipe as `r.Completion.CompletedAt`.
- It is set at the same moment it is set today: the recipe-level completion stamp, on the `finish` transition.
- Nothing else about the recipe's life changes.

### AC-03 — the vocabulary in the docs and the glossary

- `docs/03-domain-model.md`: the anatomy, the plan format and the view table say evidence.
- `docs/04-operations.md`: the verbs that take a proof say evidence.
- `GLOSSARY.md`: the **Proof** entry becomes **Evidence** — "a verification result a machine produced — the command, its exit code, what it matched, and when. An evidence that matched nothing is not evidence" — and its `_Avoid_` line drops the word "Evidence" and gains "Proof". The CPD mapping row follows. The recipe-level stamp is described there too, so the glossary does not leave a word unexplained.
- Any prose that says "proof" where it means the per-criterion record says "evidence". **Prose that means the *idea* of a proof — "the close gate checks that every criterion is proven" — may keep it** if it reads better; say in your report which ones you kept and why.

### AC-04 — the tripwires still hold

- The derived-views table, the gate table, the plan-format fields and the command count (42, after slice C) all agree with the code.
- Verification: `go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v` exit 0, and `make test` green.

## Invariants you must not break
- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps 28 external modules.
- One fact, one answer. The domain layer has **zero UI imports**.
- **A criterion's record is append-only** — evidence is added, never rewritten. Do not change that while renaming.
- Do not use `git commit --no-verify`. Do not push. Do not touch other panes. Do not spawn agents.

## Workflow
1. One commit per AC, with the per-AC evidence stamped in the commit body.
2. Before notifying: `make lint` clean, `gofmt -l .` empty, `make test` green, `git status` clean, `make regen-goldens` if the console's output changed.
3. Notify the orchestrator with a **shell command**:
```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-d.md <<EOF
Runner: slice D done. Commits: <n>. Per-AC evidence: AC-01..AC-04 stamped. Pre-notify gate: pass. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-d"
herdr agent prompt "$ORCH_PANE" "Runner: slice D done. Commits: <n>. Per-AC evidence: AC-01..AC-04 stamped. Pre-notify gate: pass. Notes: <anything the reviewer should know>."
```

### The version moves with the format

The domain's fingerprint now covers the **format facts** (the id prefix, the plan's directories, the plan's record field names), so a format change fails `TestDomainShapeIsVersioned` until the version moves. This slice renames a plan record's field (`proofs` → `evidence`), so:

- `core.DomainVersion` goes **5 → 6**, `docs/03` § 6's version paragraph says 6, and the recorded fingerprint is updated.
