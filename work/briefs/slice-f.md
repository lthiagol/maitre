You are the **runner** for slice F of this project. The orchestrator's pane id is
`w1Z:p1` — it is already substituted in this file; wherever this brief says
`$ORCH_PANE`, it means `w1Z:p1`. Use it verbatim in the notify command below. The
repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` (house style and
the working rules) and `docs/07-vocabulary.md` (the decision of record) first.

## The chef's ruling — this is not brownfield

> *"This project at this moment has no worry with brownfield — we're still making it,
> even the tests files we create can be updated to our new env, no worry with
> migrating from old to new. The idea of rewriting this project is exactly this: take
> new decisions and architecture without worrying with migrations or previous
> versions, taking into consideration the lessons learned previously."*

**What that means for the code.** The retired formats are not *supported*; they are
**gone**. No upgrade path, no compatibility layer, no second writer for a retired
format, no reader that accepts an old name. The **tolerance machinery** the rename
built — the legacy reader, the retired-name *resolution*, the old approval digests —
comes out.

**What stays — the lessons, which is the point of the sentence.** Every tripwire, the
append-only journal, **one fact one answer**, the docs-current rule, and the
fingerprint's job.

**One line separates what goes from what stays:** *the journal is history; the plan's
live records are current-format only.* The journal's events keep their old names and
ids **forever** — it is an append-only log — and the reader still resolves what the
journal names. That is not a migration path, that is a log. The **live** records
(recipes, backlog items) are read in the current format and nothing else.

## Your slice — F — one format, no tolerance

### AC-01 — the reader stops reading the retired formats

`internal/core/legacy.go` exists to read what the rename retired:

| It reads | Retired |
|---|---|
| `kitchen/specials/` and `kitchen/prep-notes/` | the two old backlog directories |
| `proofs` (for the record), `evidence` (for the stamp) | the two old field names |
| the old ids (`S01`, `P01`, `R01`) in live records | the old id forms |
| the old directories' files | the legacy console group, `FROM THE OLD PLAN` |

**All of it goes.** A plan that carries retired artifacts is **not supported** — and
the tool says so plainly rather than pretending the plan is empty. The console's
third group goes with it.

**Do not** replace the tolerance with a *different* tolerance. There is one format.

### AC-02 — the deny-list keeps the retired names; the resolution goes

`retiredAgentRoles` currently maps `line-cook` → `runner`, `taster` → `reviewer`,
`expo` → `orchestrator`, `coordinator` → `sous-chef`, and `core.agentRole()` resolves
a name through it. That resolution is a **migration affordance**: it says "the old
name *is* the new role".

**The refusal must stay.** A name like `line-cook` must never be able to sign a
recipe off — that is a safety property of the deny-list, not a compatibility feature.
Today the mapping produces that refusal *by* resolving the name to a role. Without
the resolution, the retired names must be **refused as names**.

So: the retired names become **part of the deny-list**, and the mapping's *resolver*
role goes. `line-cook` is not "the runner" any more; it is a name the list refuses.

**The tripwire follows.** `docs/03` names every role the code refuses **and** every
name the rename retired — check whether that wire still says what it should now that
the two sets are joined, and say what you did.

### AC-03 — the approval digest fences on content, and a rename never moves it

`legacyFingerprints` enumerates the encodings an older build could have written — the
`null`/`[]` list forms and the old field name. All of it goes.

But the **root cause** goes with it, or the same bug returns the next time a field is
renamed: the digest hashes the criterion's **JSON**, and the JSON carries the **field
name** — so a rename moves the digest, and a plan approved before the rename is
refused as *tampered*.

**The fix:** the digest hashes a **canonical payload with fixed field names** — a
shape the code owns, not the struct's JSON tags. Then:

- a rename of a record field does not move the digest at all;
- the digest fences on **content** and only on content;
- the candidate set is unnecessary, because there is only ever one encoding.

**The canonical payload must be a deliberate, named thing** — say what it is and why
it is written that way, in a comment a reader can check. And it must still cover
**exactly** the frozen sections: intent, scope, and each criterion's declared shape.
A criterion's *evidence* is excluded (it grows after approval) and must stay excluded.

### AC-04 — the version fence is one rule, with no upgrade path

Today the fence refuses a write when the plan is stamped **newer** than the tool —
so a new tool may write an old plan. That is the upgrade path, and the ruling
removes it.

**The new rule:** a tool writes a plan only if the plan's version **is** the tool's
version. A plan stamped with an older version is **not supported** — and the refusal
says so, in words a chef can act on, without promising a migration.

**The fingerprint's job stays, and it is not this job.** The shape is the **code-side
tripwire**: it says *the format's facts changed*, so the version must move with
them. The version is the **plan-side** fact. Two facts, two jobs — the prose in
`docs/03` and in `internal/core/lifecycle.go` must say exactly that.

### AC-05 — the tree moves to the current format

- **The repository's own plan** (`kitchen/`) is in the current format already — its
  live records are. **The journal's history stays as it is**: append-only, old names
  and old ids and all. Do not rewrite it.
- **The seeds, the demo, the fixtures and the tests** move to the current format.
- **The tests that assert the tolerance are deleted, not rewritten.** The chef said it:
  *"even the test files we create can be updated to our new env, no worry with
  migrating from old to new."* A test that asserts a retired format is read is a test
  for a feature that no longer exists. **Say in your report which tests you deleted
  and which you changed**, so the reviewer can check that nothing real went with them.
- **The journal's old ids must still resolve** — that is the history, not a migration.
  Keep that test; adapt it if the ids it names changed.

### AC-06 — the tripwires still hold

```bash
go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v
```

The derived-views table, the gate table, the plan-format fields, the command count
(**42**) and the roles the code refuses all agree with the code.

**Watch the deny-list wire carefully** — AC-02 joins two sets, and the wire reads
the spec. If the two sets are now one, the wire says so.

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps **28**
  external modules.
- The domain layer has **zero UI imports**.
- One fact, one answer: **the format is one format**, and the retired names are not
  a second answer to any question.
- **The journal is append-only.** Its history keeps its old names — that is what
  makes it a log.
- `core.DomainVersion` stays **6** unless a format fact really changed — and if it
  did, the version moves **and** the prose that explains the fence moves with it.

## The scope — what not to touch

| | |
|---|---|
| **The roles, the tabs, the backlog, the ids, the binaries** | settled by the rename. Do not undo a name. |
| **`Proof` → evidence** | settled. The record is **evidence**; the stamp is **completion**. |
| **The glossary, ADR 0009, `docs/07`** | the paper. The **glossary's** `_Avoid_` lines name retired words *on purpose* — they are warnings, and the sweep excludes them. Do not "clean" them. |
| **`docs/audit.md`** | the ledger of the night. **Never swept.** Its rows keep the names as they were filed. |
| **The mascot** | another session's work, in `logo/`. Leave it alone. |

## The discipline gate — run before you report

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

All green, and the tree clean apart from the orchestrator's own `work/` files.

**One commit per acceptance criterion**, with the evidence stamped in the body:

```
Per-AC evidence:
- AC-01: go test ./internal/core -run 'Legacy|Backlog' -v exit 0 (N/N pass)
```

**One write per file, and commit as soon as the file is coherent** — a commit is a
landing place, and the diff render is a turn's budget.

## The tree is shared

**Another session is drawing the mascot in this repository** and has files in the tree
(`logo/`). Stage **by name** — never `git add -A`, never `git add .`. Leave anything
you did not write alone. If a file you did not write has changed, **stop and report
it**.

## The report — the receipt comes first

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-f.md <<EOF
Runner: slice F done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. Tests deleted: <which>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-f"
herdr agent prompt "w1Z:p1" "Runner: slice F done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. Tests deleted: <which>. Notes: <anything the reviewer should know>."
```
