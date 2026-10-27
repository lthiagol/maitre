You are the **runner** — the same lane that is writing slice D. Your turn ended without
reporting. **This brief is the continuation.** The orchestrator's pane id is `w1Z:p1`;
use it verbatim in the notify command below.

## Where you are

| | |
|---|---|
| **Landed** | `1d378b2` **AC-01** — the per-criterion record is evidence · `e9710af` **AC-02** — the recipe-level stamp is completion |
| **In flight, uncommitted** | **AC-03** — `docs/03-domain-model.md` and `docs/04-operations.md` are edited and **not committed** |
| **Not started** | **AC-03**'s third file, `GLOSSARY.md` · **AC-04** — the tripwires |
| **The tree** | `go build ./...` OK · `go test ./internal/core ./internal/cli` **green** right now |
| **A copy of the work in flight** | [../state/slice-d-wip.patch](../state/slice-d-wip.patch) (137 lines). If the tree is ever lost: `git apply work/state/slice-d-wip.patch` restores it. |

**Do not start over. Do not revert.** The two commits are good, the journal is a clean
append (event 36 records `proofs -> evidence`, stamped 6, and the 35 historical events
are untouched), and the edits in flight are yours.

## Why your turn ended — and the trap in front of you

Your turn stopped after 4m50s, with no error, in the middle of a large diff of
`docs/03`. The diagnosis: **the diff is the trap.**

Every edit to a 700-line document re-renders that document's whole diff into the
terminal, and your pane was carrying several hundred lines of it. A dozen small
one-line edits to `docs/03` is a dozen full-file diffs — and that is what ran the turn
out, not the thinking.

**The fix is a way of working, not a workaround:**

| | |
|---|---|
| **Batch the edits** | make **one write per file** — a scripted pass (`python3` over the file, or `sed`) that applies every change to `docs/03` at once — rather than a dozen small edits |
| **Commit as soon as the file is coherent** | a commit is a landing place. A turn that ends leaves only what was committed |
| **One file, one commit, one AC** | the per-AC commit is what makes an interrupted turn recoverable — that is the whole reason for the rule |

**Finish `docs/03` and `docs/04` first, commit AC-03, then `GLOSSARY.md`.** If you
change your mind about the order, commit whichever file is coherent when the turn
threatens to end.

## What is left of AC-03 — the vocabulary in the docs and the glossary

**`docs/03-domain-model.md`** — the anatomy, the plan format and § 8's table of
derived views say **evidence**. § 8 is tripwire-checked **in both directions**: a view
that is registered but not named fails the suite, and so does a view that is named but
not registered. The view is now `evidence_summary`, computed by `EvidenceSummaryOf`.

**`docs/04-operations.md`** — the verbs that take a record say **evidence**.

**`GLOSSARY.md`** — this file is not touched yet, and it is the **single source** for
the mapping. Its **Proof** entry becomes **Evidence**:

> **Evidence**:
> a verification result a machine produced — the command, its exit code, what it
> matched, and when. An evidence that matched nothing is not evidence.

Its `_Avoid_` line drops the word **Evidence** and gains **Proof**. The **CPD mapping
row** follows — the row that currently reads `| evidence | **proof** |` becomes
`| evidence | **evidence** |`… and the recipe-level stamp is described there too, so
the glossary does not leave a word unexplained: the recipe-level stamp is
**completion**, not evidence.

> One term, one meaning, one place. `GLOSSARY.md` is where the two vocabularies are
> reconciled; if a reader reads the glossary and the code together, they must not find
> a word that means two things.

**Prose that means the *idea* of a proof may keep the word.** *"the close gate checks
that every criterion is proven"* reads better than *"…is evidenced"*. **Say in your
report which ones you kept and why** — the reviewer will look for exactly that
distinction.

## AC-04 — the tripwires still hold

```bash
go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v
```

The derived-views table, the gate table, the plan-format fields and the command count
(**42**, after slice C) all agree with the code.

## The invariants — do not break them

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps **28**
  external modules.
- The domain layer has **zero UI imports**.
- **A criterion's record is append-only** — evidence is added, never rewritten.
- The **journal is append-only**: the historical events keep their old ids and their
  old field names; the migration is an event, never a rewrite of what came before.
- `core.DomainVersion` stays **6** — it moved with AC-01's format change.
- Do not use `git commit --no-verify`. Do not push. Do not touch other panes. Do not
  spawn agents.

## The gate

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

Green, and the tree clean apart from the orchestrator's own `work/` files.

## The report — the receipt comes first

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-d.md <<EOF
Runner: slice D done. Commits: <n>. Per-AC evidence: AC-01..AC-04 stamped. Pre-notify gate: pass. Notes: <anything the reviewer should know — including the prose you kept the word "proof" in>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-d"
herdr agent prompt "w1Z:p1" "Runner: slice D done. Commits: <n>. Per-AC evidence: AC-01..AC-04 stamped. Pre-notify gate: pass. Notes: <anything the reviewer should know — including the prose you kept the word "proof" in>."
```

**Commit the receipt before you nudge.** The nudge only lands while the orchestrator is
idle, and the receipt is the half that cannot be lost.
