Runner: T3b fix 2 done. Commits: 2 (c0f1624, 89eaecc) + the state row (1932a0a). Fixed: F-120, F-121.

The job list: one line above the table now says what the table is -
"**This is the console's remit** - what the console is *for*, whether or not it is built yet. What is **built** is the [keys](#interaction) below." - and the backlog's row returns to the remit's wording: "Create, edit, shape, promote, dismiss backlog items". The clause "the free edit of an item's fields is not built" is OUT of docs/05: it was a build state, and this table is not where build states live. It stays in the plan, which says the free edit is not in this round.

How each row now reads, all in the remit's voice - what the console is FOR, whether or not it is built yet:
- View: "Read the plan well - the recipe as prose, the acceptance criteria with their evidence, the audit trail, findings anchored to code". Unchanged. The console does not read the anchor yet (Finding.Anchor is never read by internal/ui), but reading a finding is the chef's own reading, not another agent's work - so the row is in the right voice and the gap is a build gap.
- Decide: "Sign off, triage tasting notes, accept or reject a scope, choose what is next". Unchanged. The console does not accept a scope or choose what is next yet; both are the CHEF's decisions, which is the console's remit, not another agent's work.
- Manage the backlog: "Create, edit, shape, promote, dismiss backlog items" - the remit wording, as it was before cycle 1's fix.
- Run: "Start, pause, watch, steer execution". Unchanged, and this is the row that CROSSES the line by the brief's own test: start and pause are the AGENT's work, and the chef's console drives it. I did not settle it - it is a question for the orchestrator - and I left the row exactly as it was. It is the one thing in the table I would not answer myself.

The keys: the console's own rule now is that the letters it hands out do not collide. The kind's initial, and then the SHORTEST PREFIX of the kind that no other kind has taken - so "feature" and "funnel" get "f" and "fu", and both are reachable. The domain owns the kind; the console owns the key, and the console's rule is about the key. The degenerate case (every prefix taken, which means two kinds are the same word - the domain does not allow it) falls back to the first free letter, so no choice is ever unreachable. internal/core was NOT touched by this cycle.

The mutation that now fails it: kindKey was mutated back to "the first byte and nothing else" and the suite went RED - TestTheShapePromptGivesEveryKindItsOwnKey ("two kinds share the key \"f\" - one of them is unreachable"). Reverted. The domain's own kinds do not collide today (feature, bug -> f, b), so the guard is the unit test over two kinds sharing an initial - the prompt-level test (TestTheShapePromptsKeysAreDistinct) holds the keys the chef sees.

DomainVersion: 7. No format, no field, no directory. The tab count is 7 and the names did not change. No dependency added: internal/core 0, internal/cli 0, internal/ui 28 (ceiling 30).

The goldens: unchanged. `make goldens` leaves git diff empty - the goldens hold the Recipes tab, and this cycle changed the Backlog tab's prompt keys and one doc section.

The gate: make check green, make check-race green, the tree clean.

Notes for the reviewer:
- The keys rule is the console's, and the shortest-unique-prefix rule is the console's too - the domain's kinds are a var, but the KEY is an affordance, so the collision rule belongs where the key is made. Say so if you would rather the domain answered it.
- The Run row is the open question: "Start, pause" is the agent's work and the chef's console drives it. The table is the remit, so naming it is allowed; naming another agent's work as the console's own is the line the brief drew. I left it, and the report says so rather than answering it.
- The cycle-1 note about BacklogMoveOf building the kinds list on every call stays in the ledger for a later round; I did not touch it.
- The console still works end to end: the shape prompt offers "f) feature . b) bug", and shaping BL01 from the console gives it the kind feature, the scope "evidence", and the id BL01.
