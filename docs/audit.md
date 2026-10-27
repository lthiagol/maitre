# The audit ledger

Every finding the audit produced, whether it is closed, and where it was closed. It
exists because a finding that lives only in someone's head gets dropped: two notes went
missing from a plan once, and nothing noticed.

**This table is live.** A finding that is closed leaves the open table and enters the
closed one, with the file it was fixed in and the test that holds it.

## The method

Two lanes, always separate.

| Lane | Who | What it may do |
|---|---|---|
| The executor | the line cook | Fix findings. Write the code and the tests. File its own notes in a pending list. |
| The taster | an independent session | Find findings. Read the code, attack the fixes, drive the shipped binaries. **Never fixes.** |

The taster cannot see the executor's reasoning, and the executor cannot see the taster's
findings until the taster reports them. A finding is worth what an independent session
can substantiate with a command.

### The cycle matrix

| Round | What | What came out |
|---|---|---|
| Audit | The whole codebase, no diff — the repository is not under git | 39 findings, filed in one report |
| Cycles 2–7 | One fix cycle per round: fix, then review the fix | Each round closed the notes it named and left the rest |
| Cycle 8 | Two fix cycles at once, plus a feature (the tasting-note triage) | 15 notes verified; the fixes then needed two further rounds |
| Cycle 9 | Two features: the live console and the Chopboard | 15 acceptance criteria verified; one medium found in the tree's shape |
| Cycle 10 | The final sweep: the tree's cycles and the Expo tab | The tree's shape re-checked over thirteen shapes; one new doc drift found |
| Cycle 11 | The closing sweep: the derived-views table | The drift closed, and a tripwire added for the table that had drifted |
| Cycle 12 | The last sweep: the declaration became a mapping, and the triage's row | Two notes closed, and the one hole the mapping could not close |
| Cycle 13 | The closing sweep: the coverage check that fires when a view is added | One note closed, and the scope the check does not reach |
| Cycle 14 | The last sweep: the check reads the package, and the marker travels with the function | One note closed, and the families that admit a view without a classification |
| Cycle 15 | The closing review of the night: no read-side family, and a symmetric judgement | One note closed, and the shape judged |
| The rounds inside a cycle | Fix → verify → back to the executor → repeat, until the verifier is satisfied | Three rounds in cycle 8: the first fix introduced three new mediums, the second closed them |

### The discipline every round runs

| | |
|---|---|
| Cold build | `make clean && make build`, then `make test`, `make lint`, `gofmt -l .` |
| Concurrency | `go test -race` on every package, and `go test -count=2 ./...` |
| The canaries | `scale demo --dir /tmp/kb && cd /tmp/kb && maitre --render --plain`, and `scale demo --chaos` the same way |
| The real plan | The repository's own `kitchen/` must still parse, render, and its stored approval digest must still match the canonical recomputation |
| The goldens | Byte-identical — except where a feature says otherwise. The Chopboard added a tab, so the tab bar moved and nothing else did. |

### The tripwires

The docs are checked against the code, so a doc cannot drift silently. Each one is
**proved non-vacuous** by corrupting the real document in a scratch copy and watching it
fail.

| Tripwire | Checks |
|---|---|
| The gate table | Every gate the code has is listed, and every gate listed exists — both ways |
| The journal kinds | The closed enum, both ways |
| The operation marks | A ✅ mark must be a real operation; a ⏳ mark must not be |
| The command count | **Every** occurrence of "N commands" in every file that states it, not the presence of one |
| The lifecycle diagram | Every event has an edge with the same from, to and gates, and every edge is an event |
| The derived views | Every view the declaration has is in the table, and every row is a declared view; the *Computed by* cell must name the function the code uses, and the enumerations inside the rows are checked the same way |
| The view coverage | Every exported function of the package is a marked view, a declared helper, or a failure — and the judgement is symmetric: a function is a view if and only if it is marked *and* registered, so a function cannot be reclassified by deleting a word |
| The open questions | § 9's links to the proposals, the anchors they point at, the new question, and the seven that were there before |
| The lifecycles | The state table and the code's lifecycles, both ways |
| The interface | The footer paragraph names the triage |

## Open findings

| # | Severity | What | Where |
|---|---|---|---|
| **F-59** | medium | **`core.Audit.Decisions` has a shape and no writer** — it is where a decision *about the plan* would live, and nothing writes it. `kitchen/decisions.json` is named by the plan format and opened by nothing. The remedy is on the table in [06-proposals.md](06-proposals.md) § 2. | `internal/core/model.go:270`; `internal/core/store.go:24` |
| **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.resolved` gate can never fire** — nothing ever sets a question, so the list is always empty and the gate always passes. A gate nobody can fail is not a guard. | `internal/core/model.go:31`; `internal/core/lifecycle.go:197` |
| **F-61** | medium | **The `courses/` directory is named by the plan format, created by `Init`, and never written** — a directory no writer touches. Per [06-proposals.md](06-proposals.md) § 1 the honest answer is to stop naming it: a course is the relation, not a file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.md) |
| F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and the size.** It is path + mtime + size per file, so two different contents can share one fingerprint. The tool's own writes always change the mtime, so in practice only a hand edit that restores the mtime is missed — and the chef's `r` reloads it. | `internal/core/store.go` — `DiskFingerprint`. |
| P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-scope=`. | Recorded so nobody calls it a regression later. Not a defect. |

## Closed findings

| # | Severity | Closed in | Held by |
|---|---|---|---|
| F-14 | medium | `internal/cli/recipe.go` — `scale recipe update --depends-on` | `internal/cli/depends_test.go` |
| F-19 | medium | `internal/core/operations.go` — a refusal is journalled | `internal/core/refusal_test.go`, `internal/cli/refusal_test.go` |
| F-45 | medium | `docs/03-domain-model.md` — the diagram's `close` edge lists every gate | `internal/core/docs_test.go` — the diagram check |
| F-46 | medium | `docs/03-domain-model.md` — the `request_changes` edge | `internal/core/docs_test.go` — the diagram check |
| F-47 | medium | `internal/cli/doc_test.go` — every occurrence of the count is checked | `internal/cli/doc_test.go` — the half-edit test |
| F-48 | medium | `internal/ui/views.go` — a refusal that does not fit takes rows of its own | `internal/ui/refusal_test.go` |
| T3 | feature | `internal/ui/note.go`, `internal/core/views.go` — the chef triages a tasting note; resolve is not offered | `internal/ui/note_test.go`, `internal/ui/note_cursor_test.go`, `internal/core/invariant_test.go` |
| F-49 | medium | `internal/core/store.go` — a rename leaves one file | `internal/core/index_test.go` |
| F-50 | medium | `internal/core/operations.go` — a note's outcome is written once | `internal/core/noteoutcome_test.go` |
| F-51 | medium | `internal/core/lifecycle.go` — `DomainVersion` is 2 and the fingerprint covers the chef moves and the note outcomes | `internal/core/version_test.go` — the shape test |
| F-52 | medium | `internal/core/store.go` — the copy is matched on the whole id | `internal/core/index_test.go` |
| F-53 | medium | `internal/core/operations.go` — the reopen carries the outcome it erased | `internal/core/noteoutcome_test.go` |
| F-54 | medium | `internal/core/store.go` — a failed removal is not a silent partial write | `internal/core/index_test.go` |
| F-58 | medium | `internal/core/views.go` — `core.DerivedViews` holds each view's function as a value and the name the doc must match is derived from the runtime, so a renamed function breaks the wire; `docs/03` §8's table gained a *Computed by* column and marks the two planned views | `internal/core/docs_test.go` — the derived-views check |
| F-65 | medium | `internal/core/views.go` — a view added under an exempt family is registered nowhere; the read-side methods are named one by one, and the judgement is symmetric: a function is a view if and only if it is marked and registered | `internal/core/docs_test.go` — the coverage check |
| F-64 | medium | `internal/core/views.go` — the coverage check reads the whole package; a view is a function that carries the `// view` marker, the write path and the boilerplate are exempt as a family, and a function that is both fails; the read-side families are F-65 | `internal/core/docs_test.go` — the coverage check |
| F-63 | medium | `internal/core/views.go` — a view added to the code fails the wire even when the doc is untouched | `internal/core/docs_test.go` — the coverage check |
| F-62 | medium | `docs/03-domain-model.md` §8 — the triage's row is there; `core.AllTriageSections` holds the board's display order so `TriageSections` and the doc cannot disagree; § 9 links the two proposals and carries the new question | `internal/core/docs_test.go` — the derived-views and open-questions checks |
| P-1 | medium | `internal/core/lifecycle.go`, `internal/core/version_test.go` — the version guard sees the rules T3 and F-50 changed | `internal/core/version_test.go` |
| P-2 | low | `internal/core/ladder.go` — the overlay refusals are journalled too | `internal/core/refusal_test.go` |
| P-3 | low | `internal/ui/model.go` — `refusalText` reads a wrapped error with `errors.As` | `internal/ui/refusal_test.go` |
| P-4 | low | `internal/core/operations.go` — a comment on a note that does not exist is refused | `internal/cli/comment_test.go` |
| P-5 | low | `internal/cli/recipe.go` — a comment's author defaults to the actor | `internal/cli/comment_test.go` |
| P-6 | low | `AGENTS.md` — the dependency count reads 28, with the counting command written down | — |
| P-7 | trivial | `internal/ui/model.go` — the dead `m.detail` field is gone | — |
| T5 | feature | `internal/core/store.go` — `DiskFingerprint`; `internal/ui/model.go` — the poll; `cmd/maitre/main.go` — the live console | `internal/ui/live_test.go` |
| T4 | feature | `internal/core/views.go` — `TreeOf`; `internal/ui/model.go`, `internal/ui/views.go` — the Chopboard tab | `internal/core/tree_test.go`, `internal/ui/chop_test.go` |
| F-55 | medium | `internal/core/views.go` — `TreeOf` draws what the walk cannot reach, as `core.Tree.Stuck` | `internal/core/tree_test.go` |
| T6 | feature | `internal/ui/model.go`, `internal/ui/views.go` — the Expo tab and the automation's vital signs | `internal/ui/expo_test.go` |
| F-57 | medium | `docs/03-domain-model.md` §8's table and `docs/05-human-interface.md`'s footer paragraph — the moves and the tree are named, the two views the domain does not compute are marked, and `core.DerivedViews` owns the sets | `internal/core/docs_test.go` — the derived-views check |

## What this ledger cannot tell you

**Cycles 1–7's own findings are not recoverable.** Their per-finding reports lived in a
scratch directory that was wiped mid-cycle. The *fixes* are all in the repository, and
cycle 8's report re-verified fifteen of them with evidence — those are the rows above.
What is gone is the wording, the severity rationale and the file:line of each of the
original notes, and the notes that were never fixed at all. A number in the F-series
below the ones listed is **not** a gap in the ledger: it is a gap in the record, and this
table cannot invent it back.

| Lost | Recoverable from |
|---|---|
| The exact wording and file:line of cycles 1–7's findings | The fixes in the tree, and the tests that hold them |
| Which notes were closed by which round | Not recoverable |
| The findings of cycles 1–7 that were never fixed | Not recoverable — the tree is the only evidence that they are gone |

Two other limits, stated plainly:

- **A closed row is a claim about the tree, not about the history.** It says the defect
  is gone and a test holds it. It does not say which round closed it, and it cannot.
- **The open table is only as live as the last person to look.** It was written against
  the tree at the end of cycle 15's review. If a fix lands and the row is not moved,
  the ledger is wrong — which is the failure this document exists to prevent.
