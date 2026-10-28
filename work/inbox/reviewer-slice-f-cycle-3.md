Reviewer: slice F cycle 3 done. Verdict: changes-needed. F-91 closed: yes — the probe now
fails and writes nothing, and the rule lives where a name enters a record. Filed 2 new
blockers (F-92, F-93).

F-91 closed, proved by command: `ops.AdvanceWith(id, core.EventApprove,
core.TransitionContext{Approver: "line-cook", Actor: "chef"})` is refused ("\"line-cook\" is
not a name this tool accepts") and **writes nothing** — the plan still reads, the approval is
not there, the lifecycle did not move. The claim, the review, the approval, the sign-off,
the comment and the request-changes are each refused for `line-cook`, `line_cook`,
`Line Cook`, `LINE-COOK` and ` line-cook `. The **audit role** is refused at the write's
entry: nothing landed in the recipes directory and nothing reached the journal. A **new
journal event** may not carry a retired name, and the **history** is not what stopped it:
on a plan whose journal carries a hand-written `line-cook` event, a good event still lands
and the old name is still read. The funnel: `Store.SaveRecipe` runs `checkRecordNames` on
every record, `Operations.advance` checks the transition's approver (the one place
`ApproveRecipe`, `RequestChanges` and `AdvanceWith` pass through), `Operations.lock` checks
the actor and the role, and `appendJournalEventLocked` checks a new event — so the rule
does not depend on an operation remembering to check.

Mutation checks (every mutation reverted; the tree is as the runner left it, plus this
ledger): `specials` in `internal/ui/views.go`, `specials` in `internal/core/model.go`,
`"R01"` in `internal/core/store.go` and `LegacyThing` in `internal/core/store.go` — all four
still fail the wire; and a name added to the code's own `RetiredNames` ("wizard") is caught
in `internal/ui` **without touching the wire** — the wire reads the code's lists. One more
mutation proved F-93: a name-holding field the reader's list knows and the writer's does
not — the writer accepts `"line-cook"`, and the reader then refuses the plan the writer
just wrote ("this plan carries artifacts of a retired format (a signed-by (\"line-cook\"))").

F-92 (medium): `docs/03` § 1 still says *"a file named `R01-<slug>.json` is read as
readily as one named `RCP01-<slug>.json`"* and that "a reader that cannot resolve them
cannot read the past" — the code **refuses** that file (the file's name is part of the
format), and the same document says so twice (§ 6 and § 6.2). The spec contradicts itself,
and the passage a reader meets first is the wrong one. The lookup is still tolerant and the
journal keeps its old ids; the plan's records do not.

F-93 (medium): the record's name fields are enumerated twice — `actorNames` (the writer's)
and `retiredActorNames` (the reader's) — with the same doc comment, and `actorNames`'s own
comment says "the **reader** reads them when it reads", which it does not. Two answers to
one question; the drift's consequence is the hole F-91 just closed, and no wire compares
the two lists. The remedy is one list and two views of it.

The two renamed tests: `TestARetiredNameCannotBuyTheIndependenceGate` writes the bad
record **by hand** (the writer refuses it now) and still asserts the property it was
written for — the claim is refused at the door, and a record carrying a retired name is
refused by the reader — plus the writer's refusal; `TestAnOldRoleNameStillReads` and
`TestTheJournalKeepsTheOldNames` write the old event **by hand** and still assert the
history's property — the old names are kept as they were written and the reader does not
fail — and now also that the history is not what stops a new event. Neither is a lost
property wearing a new name.

LOW notes (not filed): `docs/03` § 6's compressed clause says the facts-list rule is
excepted "unless the fact is one a plan carries", which read literally would have made
slice E's change (the root files — facts a plan carries) a version bump; § 6.2 and
`docs/07` state the rule correctly ("a fact added to the list that the format already had
is not a format change; a fact the format *changed* is"). The record's name fields are a
hand-written list, so a **new** actor field added to the model must be added to it by hand
— every name-holding field in the model is in it today (ClaimedBy, Author, Reviewer, the
three `by`s and the Evidence's), but nothing checks that. `normalizeAgentName` now maps an
underscore to a dash, so `line_cook` is refused — but the same normalisation makes
`sous_chef` an agent role, which the deny-list wire does not name (the *roles* the code
refuses are unchanged: agent, sous-chef, runner, reviewer, orchestrator). The sweep's
command was rebuilt (literals only) and now reports the code's own detection and deny-list
as drift — the rule is satisfied and the one-format wire is the instrument that holds it.
