Reviewer: slice F cycle 2 done. Verdict: changes-needed. F-87..F-90 closed: yes — all four,
each verified by a command and not by the diff. Filed 1 new blocker (F-91).

F-87 closed: the retired names are refused wherever an actor names itself — the claim, the
review, the approval, the sign-off, the comment and the write's actor — and a live record
that carries one is refused by the reader. The resolution is still gone
(`sameAgent("line-cook","runner")` is false). The repository's own plan moved consistently:
`claimed_by` is `runner`, the first review round's reviewer is `runner` (they match), the
last round is the chef's, so the recipe still reads the same (not self-reviewed); the
journal keeps the old names — 18 of its 40 events still say `line-cook`, `taster`, `expo` or
`coordinator`, and the 37 historical events are untouched. F-88 closed: `docs/04` no longer
says a retired name names the role it stood for. F-89 closed: `DomainVersion` is 7, the
digest's encoding is named (`canonical-payload-v1`) and the payload's field names are in
`planFormatFacts`; a plan stamped 6 is refused by the **version** guard ("this plan was
written under domain version 6 and this tool speaks version 7") and never reaches the
digest, whose tamper message means content and only content (a changed criterion's
description is still refused). The repo's plan is stamped 7 and the declaration is an
append (event 40). F-90 closed: the scan reads the whole tree, the exemptions are per
declaration, and the token list carries the id literals, the families and the retired
stamp's access.

Mutation checks (every mutation reverted; the tree is as the runner left it, plus this
ledger): `specials` in `internal/ui/views.go`, `specials` in `internal/core/model.go`,
`"R01"` in `internal/core/store.go` and `LegacyThing` in `internal/core/store.go` — all four
now fail the wire. Probes: the claim, the review, the approval, the sign-off, the comment
and the actor are each refused for `line-cook`; a plan stamped 6 is refused by the version
guard and a plan stamped 7 is written; a changed criterion's description is refused by the
digest; the lookup still takes `RCP01`, `rcp01`, `RCP1`, `1`, `R01` and `r1` for the same
recipe, and a file named `B01-<slug>.json` is read; a write to a plan holding a retired
artifact fails and nothing lands; the repo's plan is stamped 7, is writable, and reads the
same.

F-91 (medium): the rule is enforced per operation, not where a name enters a record, and
the second approve path is not guarded — `ops.AdvanceWith(id, core.EventApprove,
core.TransitionContext{Approver: "line-cook", Actor: "chef"})` succeeds, the approval
carries `"by": "line-cook"`, and from then on the plan is refused by the reader *and* the
writer ("this plan carries artifacts of a retired format (the approver ("line-cook"))") —
the tool bricks the plan it just wrote. The shipped binaries are safe (`recipe advance`
refuses the events that have their own verb; the console approves through
`ApproveRecipe`), so the reach is the domain's own API — where the rule has to hold. The
remedy is one place: refuse the name where the transition records it.

The two renamed tests: `TestAFileNamedWithTheOldIdIsRefused` asserts the new premise (the
file's name is part of the format, and the refusal names the file and the current form) and
keeps the surviving half — the lookup's tolerance, `R01` still names the recipe once the
file is named in the current form. `TestTheRemovalLeavesAnotherKindsFileAlone` still
asserts what it was written for — the removal reads the *whole id*, so the C01-named file
survives the re-save of RCP01 — with the premise moved (the stray record's id is now a
recipe's, because a record whose id is not a recipe's id is an unsupported plan, so the old
scenario could not exist in a supported plan). Neither is a lost property wearing a new
name.

The question cycle 1 left open, now closed: the rule is **one coherent rule** — a *typed
name* is a lookup (the tool strips the leading letters and answers with the number, because
the number is the identity and the directory carries the kind), while a *record* and a
*file's name* on disk are the format and carry one form. The position decides, not the
string: the same text is a lookup in one position and a format violation in the other. The
sweep's list is satisfied in the code — the retired names live only in the detection that
refuses them and the deny-list that refuses them as names, per declaration, over the whole
tree.

LOW notes (not filed): `docs/03` § 6.2 says "Now the shape moves the shape" (a typo) and
still states "a change to the facts list … moves the fingerprint and leaves the version
alone" while this cycle grew the facts list *and* moved the version — the real reason is
that version 6 covered two digest encodings, so the stamp had to be re-cut, not that "the
writer changed"; the same sentence is in `docs/07`. The audit context's **role** is not
checked, so `--role line-cook` writes a retired role into the journal (the journal is
history, so no live record carries it). `RequestChanges(id, by, notes)` does not check `by`,
and the transition drops the name anyway. The retired id *form* is caught by its literals
(`R01`, `S01`, `N01`), so a generic prefix check (`strings.HasPrefix(id, "R")`) would not
trip the wire. The retired names are refused under the tool's own spelling rules (the case
does not matter, a space stands for a dash), so a variant spelling such as `line_cook` is
another name — consistent with the deny-list being a speed bump against the accident, not a
boundary against the lie. Slice E's sweep command now reports the code's own detection and
deny-list as drift (25 lines in `retired_test.go`, 11 in `model.go`, 7 in `retired.go`): the
rule is satisfied, but the instrument needs the new wire's exemptions. A file whose *name*
and *record* disagree (`C01-rogue.json` holding `RCP02`) is read as a recipe — only the
retired name forms are refused in a file's name.
