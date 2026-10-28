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
| Slice A | The vocabulary rename: the module, the binaries, the recipe ids, the version | The ids and the guard hold; two mediums and four lows filed |
| Cycle 2 | The fixes to slice A: the usage lines, the two-way wire, the plan's stamp, the format facts | Six notes closed, and the fence checked by mutation |
| Slice B | The vocabulary rename: the roles and the tabs | The roles resolve as written *and* as a role, so a plan written before the rename still cannot close on its own executor's review; the tabs are Board and Sous-Chef, and the honesty of the waiting tab held; one medium filed in the spec |
| Slice B · cycle 2 | The fixes to slice B: the spec's vocabulary rule, and the table under it | F-72 closed — the sentence, the table and the decision of record say the same thing; the ledger's two pointers follow the rename |
| Slice C | The vocabulary rename: the backlog becomes one artifact with a kind | The id, the kind and the fence hold; the live console hides what the snapshot shows, and the old plan's dangling pointer no longer reserves its id — one high and three mediums filed |
| Slice C · cycle 2 | The fixes to slice C: the live reload path, the dangling pointer, the list rule and the spec's row | F-74..F-77 closed — each verified by a mutation, and the four fixes moved nothing but the four findings' own code |
| Slice D | The vocabulary rename: the per-criterion record is evidence, the recipe-level stamp is completion | The two words hold, the reader tolerates the old names and the tripwires hold both ways; the approval digest did not travel with the rename, so a pre-rename plan is refused as tampered — one high and two mediums filed |
| Slice D · cycle 2 | The fixes to slice D: the digest travels with the rename, the glossary covers the record it names, the version paragraph stops contradicting the reader | F-78..F-80 closed — each verified by a mutation, and the fixes moved nothing but the three findings' own code |
| Slice E | The vocabulary rename: the paper — the glossary rewritten, ADR 0009 recorded, the tree swept | The ids and the guard hold and the sweep's 78 lines are all records of something that still exists; the retired-word pattern holds for seven renamed terms and misses the console's, and the spec and the `DomainVersion` comment still describe the fingerprint the way the slice found it — five mediums filed |
| Slice E · cycle 2 | The fixes to slice E: the glossary's retired names, the console's place, the spec's two lists and the version/fingerprint split — **the last slice of the rename** | F-81..F-85 closed, each verified by a command and not by the diff; the sweep's count reconciles (78 → 77: F-83 removed one prose line and nothing was added, the console's `_Avoid_` line already existed — but the published pattern cannot see the name F-81 put on it) — one medium filed, against the sweep's command, which does not look for the console's retired name |
| Slice E · cycle 3 | The last fixes to slice E: the sweep's instrument, and the two notes the review left | F-86 closed — the command carries all seven retired names and tells a warning from the drift (72 lines in 18 files; the count reconciles, 77 − 5), the glossary's heading is **The surfaces** and `docs/03`'s pointer lands on § 6.2 — **the rename's last slice, closed**: the glossary rewritten, ADR 0009 written, the tree swept, and **the mascot's acceptance criterion carried to the session that owns the drawing** (AC-01) |
| Slice E · closed | The rename's last slice — the glossary rewritten, ADR 0009 written, the tree swept | F-81..F-86 closed across three rounds, each verified by a command; the sweep's command carries all seven retired names and tells a warning from the drift, and its **72 lines in 18 files** are all records of something that still exists — **AC-01, the mascot's drawing, is carried to the session that owns it**, and nothing of that session's work was committed here |
| Slice F | The tolerance machinery leaves: the reader stops reading the retired formats, the deny-list keeps the retired names, the resolution goes, the digest fences on content, the version fence is one rule | The retired formats are gone and a plan carrying one is refused plainly; a rename does not move the digest; the fence refuses in both directions — and **the independence gate can now be bought with a retired name**, the operations doc still describes the resolution that is gone, the digest's encoding moved without either fence moving, and the new "one format" wire has three holes — four mediums filed |
| The rounds inside a cycle | Fix → verify → back to the executor → repeat, until the verifier is satisfied | Three rounds in cycle 8: the first fix introduced three new mediums, the second closed them |

### The discipline every round runs

| | |
|---|---|
| Cold build | `make clean && make build`, then `make test`, `make lint`, `gofmt -l .` |
| Concurrency | `go test -race` on every package, and `go test -count=2 ./...` |
| The canaries | `maitre-cli demo --dir /tmp/kb && cd /tmp/kb && maitre --render --plain`, and `maitre-cli demo --chaos` the same way |
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
| The plan format | Every directory and file the code creates is named in `docs/03` § 6, and every directory and file § 6 names is the code's — both ways — and the section's version is the code's |
| The usage lines | Every usage line the CLI prints names a binary that exists, and the console's errors name it too |
| The deny-list | The spec names every role the code refuses, and every name the rename retired — a plan written before the rename still meets those |
| The one format | The retired names may live in exactly two places — the detection that refuses them (`retired.go`) and the deny-list that refuses them as names (`model.go`) — and the scan reads the **code**, not the prose |

## Open findings

| # | Severity | What | Where |
|---|---|---|---|
| **F-59** | medium | **`core.Audit.Decisions` has a shape and no writer** — it is where a decision *about the plan* would live, and nothing writes it. `kitchen/decisions.json` is named by the plan format and opened by nothing. The remedy is on the table in [06-proposals.md](06-proposals.md) § 2. | `internal/core/model.go:270`; `internal/core/store.go:24` |
| **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.resolved` gate can never fire** — nothing ever sets a question, so the list is always empty and the gate always passes. A gate nobody can fail is not a guard. | `internal/core/model.go:31`; `internal/core/lifecycle.go:197` |
| **F-61** | medium | **The `courses/` directory is named by the plan format, created by `Init`, and never written** — a directory no writer touches. Per [06-proposals.md](06-proposals.md) § 1 the honest answer is to stop naming it: a course is the relation, not a file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.md) |
| F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and the size.** It is path + mtime + size per file, so two different contents can share one fingerprint. The tool's own writes always change the mtime, so in practice only a hand edit that restores the mtime is missed — and the chef's `r` reloads it. | `internal/core/store.go` — `DiskFingerprint`. |
| P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-scope=`. | Recorded so nobody calls it a regression later. Not a defect. |
| **F-87** | medium | **The independence gate can be bought with a retired name.** `sameAgent` no longer resolves the retired names, so `review.reviewer_independent` passes for a recipe claimed by `line-cook` and reviewed by `runner` — one agent, two spellings, and the gate calls it independent. The retired names are still carried by the code (the deny-list refuses them as approvers), and the repository's own live record carries one: `kitchen/recipes/RCP01-*.json` says `"claimed_by": "line-cook"`. The comment above `sameAgent` says the executor and the reviewer "are compared as **agents**, not as strings" — the code compares strings, and a retired name is now a second spelling of the runner. The remedy that keeps the ruling: refuse the retired names where an actor name is written into a live record (the claim and the review round), which is what "one format" means for a record — and then the repository's own plan's `claimed_by` moves with it. | `internal/core/views.go:422,447`; `internal/core/model.go:530`; `kitchen/recipes/RCP01-the-console-shows-a-finding-s-conversation.json:120` |
| **F-88** | medium | **`docs/04-operations.md` still describes the resolution the slice removed.** The deny-list paragraph says the retired names "**still name the role they stood for**, so a plan written before the rename is refused exactly like one written after it". The resolution is gone: a retired name is a name the list refuses and nothing more. The refusal still happens — the deny-list carries the names — so the paragraph's conclusion holds and its reason does not. `docs/03` was corrected and `docs/07` was corrected; `docs/04` was not, and the deny-list wire reads `docs/03` only, so nothing catches it. | `docs/04-operations.md:242` |
| **F-89** | medium | **The approval digest's encoding changed and neither fence moved.** The canonical payload dropped the frozen criterion's `evidence` key, so the digest of the same recipe is a different document: the tool now writes a different `approval.fingerprint`, and `planFormatFacts` — the field names *on disk* — cannot see a payload that is not on disk. A plan approved by the previous build therefore passes the version guard (both are version 6) and is refused at the first gated transition with `immutability.respected (frozen sections changed since approval)`: the tool calls a plan nobody touched tampered, and with one encoding and no candidate set there is nothing to re-stamp it. The repository's own plan was re-stamped by hand and the journal carries a hand-written migration event — the migration the ruling forbids, and the proof that the encoding was a rule the fences should have seen. | `internal/core/model.go:453`; `internal/core/format.go:20`; `internal/core/lock.go:79` |
| **F-90** | medium | **The new "one format" wire has three holes.** It scans `internal/core` only — the console, the package that carried the retired group until this slice, is not scanned, and a mutation that puts `"specials"` in `internal/ui/views.go` passes. Its allow-list is a **file**, so `model.go` is wholly exempt and a reader can hide in the file that holds the deny-list. And its token list names four spellings of "legacy" (`ListLegacyItems`, `LegacyItem`, `legacyFingerprint`, `restampLegacy`) while omitting the retired id literals the slice retired (`R01`, `S01`, `N01`) — a fifth spelling of "legacy" and a mapping from the retired id form both pass. The retired recipe-level stamp's name (`evidence`) cannot be a token at all: the criterion's current field shares the word, so that one retired name needs another instrument. | `internal/core/retired_test.go:210` |

## Closed findings

| # | Severity | Closed in | Held by |
|---|---|---|---|
| F-14 | medium | `internal/cli/recipe.go` — `maitre-cli recipe update --depends-on` | `internal/cli/depends_test.go` |
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
| F-66 | medium | `internal/cli/recipe.go`, `fragments.go`, `backlog.go` — every usage line names `maitre-cli`; `cmd/maitre/main.go` — the console's errors name `maitre-cli init` | `internal/cli/usage_test.go` — the usage-lines and console-errors checks |
| F-67 | medium | `internal/core/planformat_test.go` — the plan-format wire checks both ways: the code's directories and files are in the doc, and the doc's are in the code | `internal/core/planformat_test.go` — `TestTheDocNamesThePlanFormatTheCodeHas` |
| F-68 | low | `internal/core/migration.go` — `RecordMigration` appends a `control` event carrying the version, through the domain's own journal path; `scripts/stamp-plan/` — the script that records it | `internal/core/migration_test.go` |
| F-69 | low | `internal/cli/cli.go` — `MAITRE_CLI_ACTOR`, `MAITRE_CLI_SESSION` and `MAITRE_CLI_ROLE` are read; the old `SCALE_*` names still work as a fallback and warn | `internal/cli/env_test.go` |
| F-70 | low | `docs/decisions/0004-naming.md` — the table of what 0009 supersedes; `docs/decisions/0001-stack.md`; the six `prototypes/*.txt` sketches; the code comments | — (a doc fix; nothing to hold it) |
| F-71 | low | `internal/core/format.go` — the facts of the plan's format; `internal/core/version_test.go` — the fingerprint hashes them, so a format change moves the version; `docs/03-domain-model.md` §6 — the version paragraph | `internal/core/version_test.go` — the shape test; `internal/core/planformat_test.go` — the version check |
| P-1 | medium | `internal/core/lifecycle.go`, `internal/core/version_test.go` — the version guard sees the rules T3 and F-50 changed | `internal/core/version_test.go` |
| P-2 | low | `internal/core/ladder.go` — the overlay refusals are journalled too | `internal/core/refusal_test.go` |
| P-3 | low | `internal/ui/model.go` — `refusalText` reads a wrapped error with `errors.As` | `internal/ui/refusal_test.go` |
| P-4 | low | `internal/core/operations.go` — a comment on a note that does not exist is refused | `internal/cli/comment_test.go` |
| P-5 | low | `internal/cli/recipe.go` — a comment's author defaults to the actor | `internal/cli/comment_test.go` |
| P-6 | low | `AGENTS.md` — the dependency count reads 28, with the counting command written down | — |
| P-7 | trivial | `internal/ui/model.go` — the dead `m.detail` field is gone | — |
| T5 | feature | `internal/core/store.go` — `DiskFingerprint`; `internal/ui/model.go` — the poll; `cmd/maitre/main.go` — the live console | `internal/ui/live_test.go` |
| T4 | feature | `internal/core/views.go` — `TreeOf`; `internal/ui/model.go`, `internal/ui/views.go` — the Chopboard tab | `internal/core/tree_test.go`, `internal/ui/board_test.go` |
| F-55 | medium | `internal/core/views.go` — `TreeOf` draws what the walk cannot reach, as `core.Tree.Stuck` | `internal/core/tree_test.go` |
| T6 | feature | `internal/ui/model.go`, `internal/ui/views.go` — the Expo tab and the automation's vital signs | `internal/ui/souschef_test.go` |
| F-57 | medium | `docs/03-domain-model.md` §8's table and `docs/05-human-interface.md`'s footer paragraph — the moves and the tree are named, the two views the domain does not compute are marked, and `core.DerivedViews` owns the sets | `internal/core/docs_test.go` — the derived-views check |
| F-72 | medium | `docs/03-domain-model.md` § "The CPD ↔ tool mapping" — the rule sentence and the table say the same thing now, and both say what the decision of record says: the kitchen names the artifact and the place, the agents' jobs stay plain, and the **Kept plain** cell carries `chef` · `sous-chef` · `runner` · `reviewer` · `orchestrator`. `maître d'` goes: it was the human-facing console's themed name, the console is `maitre` now, and it was never a job. | — (a doc fix; the deny-list wire holds the role list, nothing holds the sentence) |
| F-73 | low | **`kitchen/project.json` no longer carries `"name": "new-project"`** — the decision of record names the project **maitre**, and the console's header renders the plan's name, so the repo's own console reads `maitre · maitre`. It is the *plan's* data, not the code's: slice A renamed the module and the binaries and left this file, and slice C was the slice that migrated the repo's own plan. | `kitchen/project.json:8` |
| F-74 | high | **`internal/ui/model.go` — `Model.applyReload` puts a reloaded plan into the model, and the live path and the static render both go through it, so a field the render carries and the live path forgets cannot happen; the test drives the real path** — `reloadMsg` carried the legacy artifacts and `Update`'s `reloadMsg` branch dropped them, so an old plan's backlog read as empty in the console the chef actually uses, while `--render` showed them. One place now knows the model's shape, and `loadModel` sends the reload message through the model's own `Update` instead of copying the fields by hand — the copy was the bug. | `internal/ui/model.go` — `applyReload`; `internal/ui/render.go` — `snapshot`; `internal/ui/console_test.go` — `loadModel`, `TestAnOldPlanStillRenders`, `TestTheLiveConsoleCarriesTheLegacyGroup` |
| F-75 | medium | **`internal/core/store.go` — `reservedRecipeIDs` scans the old plan's pointers as well as the backlog's, and `core.LegacyItem` carries `PromotedTo`** — the reader dropped the pointers in `kitchen/prep-notes/*.json`, so an id a dangling pointer names could be handed to a stranger — the exact failure `NextRecipeID`'s own comment says must never happen. The pointer is read in the format that is actually on disk (`"recipe RCP10"`), and a pointer in the older `"recipe R10"` spelling reserves its number too; a bare `"RCP10"` is skipped, which is deliberate — the writer only ever wrote the two-word form. | `internal/core/store.go` — `reservedRecipeIDs`; `internal/core/legacy.go` — `LegacyItem.PromotedTo`; `internal/core/backlog_concurrency_test.go` — `TestALegacyDanglingPointerResolvesToNothing` |
| F-76 | medium | **`internal/cli/flags.go` — the parser knows which list flags were given, and `listArg` returns nil for a flag that was not given; `internal/core/backlog.go` — one rule (`setList`) for both list fields** — the tag list was guarded with `!= nil` while the scope list, one line away, used `len(...) > 0`, so an update that never mentioned the tags dropped them. nil now means "not given" and a list means "replace the list it names" — an empty list that *was* given clears the field. `emptyIfNil` keeps the JSON edge a list: an empty one is `[]`, never `null`. The other verbs are untouched: `--in-scope`, `--out-of-scope`, `--ac`, `--ac-verify` and `--depends-on` behave exactly as before. | `internal/cli/flags.go` — `given`, `listArg`; `internal/core/backlog.go` — `setList`, `emptyIfNil`; `internal/cli/backlog_test.go` and `internal/core/backlog_test.go` — `TestTheBacklogUpdateLeavesTheListsAlone` |
| F-77 | medium | **`docs/03-domain-model.md` § 1.1 — the row and the paragraph now agree, and the code did not move** — the row claimed a promoted small task's recipe is "falsifiable from the start" while the paragraph below said its criterion carries no verification and "cannot pass `ac.falsifiable`". The criterion makes the work **checkable** and becomes falsifiable when the recipe's own life writes its verification; a promoted recipe is **checkable, not yet falsifiable**. `promotedCriteria` still leaves the verification empty, and `TestThePromotionWeakensNoGate` still pins it. | `docs/03-domain-model.md` § 1.1; `internal/core/backlog.go` — `promotedCriteria`; `internal/core/promotion_test.go` — `TestThePromotionWeakensNoGate` |
| F-78 | high | **The approval digest did not travel with the rename — a plan written before it is now refused as tampered.** *Closed in* `internal/core/model.go` — `legacyFingerprints` enumerates the **field name** as well as the list forms, so the old name's digest is a recognised encoding and `restampLegacyApproval` re-stamps it canonically; the candidates are the product of the two axes, and a tampered recipe still matches none of them. The repository's own plan is migrated: `kitchen/recipes/RCP01-the-console-shows-a-finding-s-conversation.json` carries the canonical digest, and the journal records it as an **append** — event 37. | `internal/core/approval_test.go` — `TestARenamedFieldIsNotATamper`, `TestTheRepoPlansApprovalDigestMatchesItsContent`; `internal/core/model.go` — `legacyFingerprints`, `legacyFingerprint` |
| F-79 | medium | **The glossary's entry for _Evidence_ is narrower than the record it names.** *Closed in* `GLOSSARY.md` — the **Evidence** entry carries the exception now (*"One verification result — a machine's, or the chef's … A result that matched nothing does not count; the chef's sign-off is the one exception, and it claims no match"*), and `docs/03` § 5 says the same thing, so the single source and the spec no longer contradict each other. | `internal/core/approval_test.go` — `TestASignedOffCriterionIsPassedButNotProven`; `internal/core/lifecycle_test.go` — `TestFinishRequiresProvenCriteria`; `docs/03-domain-model.md` § 5 — (the two docs' agreement is not wire-checked) |
| F-80 | medium | **`docs/03` § 6 says the reader does not see the old records, two paragraphs after saying it reads them.** *Closed in* `docs/03-domain-model.md` § 6 and `internal/core/lifecycle.go`'s version notes — the version-5 and version-6 sentences now say what the reader does: it reads the old field name, it shows the old specials read-only in their own group, and it never looks for new ones. The version bump is unchanged, because the **writer** changed, and that is what the version fences. | `internal/core/legacy.go` — `readLegacyEvidence`, `ListLegacyItems`; `internal/core/legacy_test.go`; `internal/ui/console_test.go` — `TestTheLiveConsoleCarriesTheLegacyGroup` |
| F-81 | medium | **The console's retired name is carried nowhere in the glossary, and the proof offered for the pattern does not reproduce.** *Closed in* `GLOSSARY.md` — the console's entry carries *maître d'* once, on its `_Avoid_` line, like every other retired name. The proof now reproduces: the grep over the seven retired names returns seven hits, one per name, and **0** of them outside an `_Avoid_` line. | `GLOSSARY.md:161`; `grep -n -i -E` over the seven retired names in `GLOSSARY.md` — 7 hits, 0 outside `_Avoid_` |
| F-82 | medium | **The glossary files the console among the crew, where its own new tables say it does not belong.** *Closed in* `GLOSSARY.md` — "### The crew" holds five entries (Chef · Sous-chef · Runner · Reviewer · Orchestrator); the console moved to "### The strings", headed `maitre` the way the names table spells it (`:32`), and the section says the console is a **surface**, not a job, and names what it carries — the window title, the tab bar and the CLI's help. The entry itself still describes the console, not a string. | `GLOSSARY.md:89`; `GLOSSARY.md:151-161`; `GLOSSARY.md:32` |
| F-83 | medium | **Two retired words are named in prose as well as in the `_Avoid_` line.** *Closed in* `GLOSSARY.md` — the backlog item's prose points at [ADR 0009](decisions/0009-vocabulary.md) § 3, which says what the two file types were, so the reader still learns what the backlog item replaced; the `_Avoid_` line keeps the retired words. | `GLOSSARY.md:59-63`; `docs/decisions/0009-vocabulary.md` § 3 |
| F-84 | medium | **The spec lists the fingerprint's facts twice and the two lists disagree.** *Closed in* `docs/03-domain-model.md` and `internal/core/lifecycle.go` — all three lists (two in the spec, one in the `DomainVersion` comment) now say the same four things: the id prefixes, the directories, the root files and the field names on disk, and that is what the code hashes. | `docs/03-domain-model.md:451`; `docs/03-domain-model.md:552`; `internal/core/lifecycle.go:200`; `internal/core/format.go:20` — `planFormatFacts` |
| F-85 | medium | **The spec still states the coupling the slice split.** *Closed in* `docs/03-domain-model.md` § 6, § 6.1 and the version-4 history, `internal/core/lifecycle.go`'s `DomainVersion` comment, and `docs/07-vocabulary.md` — the **version fences the writer** and the **fingerprint is the tripwire**: a change to what the tool writes moves both, a change to the facts list moves the fingerprint alone. The code did not move: `core.DomainVersion` is still **6**, and the guard still refuses a plan stamped 7 and accepts one stamped 5. | `internal/core/version_test.go` — the shape test (`recordedDomainShape = "b84b230161f4ab03"`), `TestDomainVersionRefusesANewerPlan`, `TestDomainVersionStampsAnOlderPlan`; `internal/core/lifecycle.go:216` |
| F-86 | medium | **The rename's sweep cannot see the console's retired name — the name is missing from the instrument's own list of old names.** *Closed in* the **instrument**, not the tree: AC-04's list of names (`work/briefs/slice-e.md`) now carries *maître d'* and says the command must carry every name and must tell a retired word from a **warning about** one; the command (commit `cea2380`, and `work/inbox/runner-slice-e-fix-2.md`) matches all seven retired names and excludes the glossary's `_Avoid_` lines. No file of the tree changed — the tree was clean, and the proof was blind. | the command, run at HEAD — **72 lines in 18 files**, and the count reconciles: 77 − 5 (the glossary's `_Avoid_` lines the old command counted as drift) = 72; the complete pattern without the glossary exclusion returns 78, the extra line being the console's own `_Avoid_` line, and every other hit of *maître d'* is in the excluded dated records (ADR 0004, 0005, 0009) or in the ledger |

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
  the tree at the end of the cycle 2 review. If a fix lands and the row is not moved,
  the ledger is wrong — which is the failure this document exists to prevent.
