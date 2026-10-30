Runner: T3b fix 1 done. Commits: 2 (4e14f19, 58069e7) + the state row (af20145). Fixed: F-118, F-119.

The doc: docs/05-human-interface.md line 14 - the console's job list. It said "Create, edit, shape, promote, dismiss backlog items"; it now says "Create, shape, promote, dismiss backlog items - the **free edit** of an item's fields is **not built**". The console shapes, promotes, dismisses and creates; it does not edit an item's fields, and the doc said otherwise.

Every file I checked for the same promise: docs/05-human-interface.md (fixed - the job list was the only place in it), docs/03-domain-model.md (the domain's shaping semantics - the shaping row now also says the domain names the kinds a shaping may pick; no console promise anywhere), docs/04-operations.md (the CLI's own `backlog update --title --body --scope --kind --tag` verb - the CLI is unchanged and does have the free edit), docs/06-proposals.md, GLOSSARY.md, AGENTS.md, the root README, and work/tui/plan.md (the orchestrator's own correction, 05058d4 - I checked it and did not touch it). The one place left that says "edit" is docs/decisions/0003-surface-split.md line 50, and it stays: it is a decision of record about who owns the backlog's management (the console, not the agent), written in the pre-rename vocabulary ("ideas and briefs"), and it is not a statement about what is built today. Rewriting a decision to match a build state would erase the decision. Say so if the reviewer disagrees.

The kinds: they come from the DOMAIN, through the move answer the console already asks. core.BacklogMoveOf now carries the kinds a shaping may pick - the domain's small tasks, read from the rule (isSmallTask) and not copied from it - so a kind the rule gains is a kind the console's prompt offers the same day. The console builds the prompt's choices from that answer and carries no list of its own; the letter that stands for a kind stays the console's own affordance (the domain owns the kind, not the key). No new exported function: the answer is the move answer, so the wire has nothing new to check - the // view marker, the registration in core.DerivedViews and the row in docs/03 § 8 all still hold, and the row now says the answer carries the kinds.

The mutation that now fails it: the console's shapeChoices was mutated to carry its own list - feature and a "chore" that the domain does not know - and the suite went RED: TestTheShapePromptOffersTheDomainsKinds ("the prompt offers \"chore\", want the domain's \"bug\"") and TestTheShapePromptAsksForTheKindAndTheScope ("the prompt should offer \"b) bug\""). Reverted.

DomainVersion: 7. No format, no field, no directory. The tab count is 7 and the names did not change. No dependency added: internal/core 0, internal/cli 0, internal/ui 28 (ceiling 30).

The goldens: unchanged. `make goldens` leaves git diff empty - the goldens hold the Recipes tab, and this cycle changed the Backlog tab's prompt and one doc line.

The gate: make check green, make check-race green, the tree clean.

Notes for the reviewer:
- F-117 (the orchestrator's file): I checked 05058d4's correction. The opening table's row 4 now says "**shape it** . dismiss . create reach the console" and says the free edit is not in this round, and the Promote row now says the kind changes to a small task "the chef picks which" - which matches what the console does. The correction is in the three places the finding named.
- The prompt's choices are now built at runtime from the domain's answer, so the console's shapeChoices is a function, not a package-level list. If a reviewer wants the kinds to be a separate domain answer (a registered view of their own) rather than a field on the move, say so - the brief allowed either, and I took the one that adds no exported name.
- Two console-wide observations I did not touch: the cursor keeps its index after a write (pre-existing, console-wide), and the moves act while the detail pane is open (the console's model).
- work/tui/plan.md's § 4 table still has a "Today" column that says "not in the console" for the shape - a snapshot written before the slice landed. The orchestrator's file; I did not touch it.
