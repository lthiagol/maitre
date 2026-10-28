Runner: slice D fix 1 done. Commits: 3 (ac156ff F-78, 8c5899d F-79, f09a682 F-80).
Fixed: F-78, F-79, F-80. Pre-notify gate: pass (make lint clean, gofmt -l . empty, make
test green, git status shows only the other session's `logo/`). The gate table, the
plan format, the command count (42) and `core.DomainVersion` (6) did not move.

**F-78 — the approval digest travels with the rename** (`ac156ff`). The digest hashes
the criterion's JSON and the JSON carries the field **name**, so the version-5 build
hashed `proofs` and the version-6 build hashes `evidence`. The repo's own RCP01 stored
the old name's digest (8813a266…), and `immutability.respected` read that as
tampering.

| | |
|---|---|
| The tolerance | `legacyFingerprints` enumerates the **field name** now, not only the list forms: `legacyFingerprint` hashes the frozen criterion with the record under its older name. The candidates are the product of the two axes |
| The completeness argument | rewritten, not extended in silence: the payload's only list-shaped parts are the three lists, the only thing an older build could write differently about them is `null` vs `[]`, and the only thing the rename changed is the criterion's field name |
| The plan | migrated: RCP01's stored digest is the canonical one (4984098d…). The journal records it as an **append** — event 37, `approval-fingerprint:proofs → approval-fingerprint:evidence`. Events 1–36 keep their old ids and their old field names |
| The standing discipline | `TestTheRepoPlansApprovalDigestMatchesItsContent` reads the plan's stored digest off the disk and compares it with what the tool recomputes. Nothing held it before |

Proof by mutation: dropping the old-name candidate from `legacyFingerprints` fails
`TestARenamedFieldIsNotATamper` ("the approval should be re-stamped canonically"); the
other tamper tests still pass. Reverted.

**F-79 — the glossary's entry covers the record it names** (`8c5899d`). The entry said
*"A verification result a machine produced"* and *"An evidence that matched nothing is
not evidence"* — while a human sign-off lives in the same field and `ac.all_proven`
takes it as a pass. The entry now reads:

> One verification result — a machine's, or the chef's: the command, its exit code,
> what it matched, and when. A result that matched nothing does not count; the
> chef's sign-off is the one exception, and it claims no match.

`docs/03` § 5 says the same thing (evidence is attached, never prose — a machine's, or
the chef's; zero-match does not count; **a sign-off is evidence too**), so the two no
longer contradict each other.

**F-80 — the version paragraph stops contradicting the reader** (`f09a682`). The
version-5 and version-6 sentences in `docs/03` § 6 now say what the reader does: it
reads the old field name, it shows the old specials read-only in their own group, and
it never looks for new ones. The version bump is unchanged — the *writer* changed, and
that is what the version fences. The same sentences moved with it in
`internal/core/lifecycle.go`'s version notes.

Notes for the reviewer:

1. **"An evidence" is gone from the tree** — outside `docs/audit.md`, which is yours.
   The grammar note was folded in as F-79 says, so the wording moved everywhere the
   word appears: the core's comments, the seeded demo plan, P5 in the principles, C7 in
   the constraints, the operations doc, the brownfield proposal, ADR-0006 and 0002,
   `AGENTS.md` and the decisions index. ADR-0006's **title** moved with it ("what
   evidence is"), and the index and `AGENTS.md`'s decisions row follow.

2. **The stamp-plan usage line names the field-name migrations** now (`proofs` →
   `evidence`), as the ruling asked — one line, in F-78's commit, next to the script
   that recorded them.

3. **The LOW rulings I did not act on, and why.** `docs/07`'s *"work that closes with
   no proof"* stays (it names the idea of a proof, and `docs/07` is yours);
   `research/01-cpd-foundation.md` stays (history); `GLOSSARY.md`'s brigade section
   stays (slice E rewrites the glossary); `AGENTS.md` and `work/README.md` stay (yours,
   and they are right).

4. **The tree is shared, and I staged by name.** Every commit names its files
   explicitly — never `git add -A`. `logo/` is untouched and uncommitted, and I did not
   write to the plan except for the migration event (37), which the stamp-plan script
   appended through the domain's operation.

5. **One thing worth your eye.** The finding said the *stored* digest was the old
   name's digest; it was, and the plan is migrated to the canonical one. The reader
   also recognises the old digest, so a plan written before the fix is still readable
   — the tolerance and the migration are both in place, and neither depends on the
   other.

LOW notes: none.
