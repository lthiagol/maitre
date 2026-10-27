Runner: slice C done. Commits: 7 (one per AC, plus one docs commit). Per-AC
evidence: AC-01..AC-06 stamped in each commit body. Pre-notify gate: pass
(make lint clean, gofmt -l . empty, make test green, git status clean,
make regen-goldens clean — the Recipes tab's golden did not move).

What the slice did: the backlog is **one artifact with a kind**. `Special` and
`PrepNote` are gone; `kitchen/backlog/` holds `BL01` items with `kind` ∈ idea |
feature | bug. `maitre-cli special` and `maitre-cli prepnote` are one group,
`maitre-cli backlog` (list · create · update · dismiss · promote). The count
went 46 → 42, and the doc, AGENTS.md and the counting tripwire say 42
together. `core.DomainVersion` 4 → 5, the fingerprint is re-recorded, and
docs/03 § 6 says 5.

Notes the reviewer should have:

1. **The repo's own plan had no backlog item.** The brief said "it has a
   promoted item and a served recipe — the item's file moves to
   kitchen/backlog/". It does not: `kitchen/specials/` and
   `kitchen/prep-notes/` were **empty** (untracked, and never in git — the
   repo's plan has only ever held the recipe and the journal). So there was
   nothing to move, and I did not fabricate one. What I did migrate: the plan's
   **shape** (the format change is recorded in the journal as an append —
   `git diff` on kitchen/activity.json shows an addition, never a rewrite),
   the directories (`kitchen/backlog/` replaces the two), and
   **F-73**: `kitchen/project.json` no longer says `"name": "new-project"`,
   so the console's header reads `maitre · maitre`. The old ids still
   resolve: `TestThePreMigrationBacklogIdStillResolves` (the pre-migration id
   `S01` names the item `BL01`) and
   `TestTheRepoPlansJournalIdsStillResolve` (every id the repo's own journal
   names resolves).

2. **An old plan does not crash the console, and the console is honest.** The
   core reads `kitchen/specials/` and `kitchen/prep-notes/` through
   `core.Store.ListLegacyItems`, and the Backlog tab shows those artifacts in
   a third group, *From the old plan*, labelled with the rung they came from
   (`special`, `prep note`). The console does **not** write the old format, so
   it offers no move for them — the footer says nothing and `space` says why.
   Test: `TestAnOldPlanStillRenders`.

3. **The promotion weakened nothing.** No gate rule changed, no gate was
   weakened, and no plan-format field was added for the promotion. The
   promoted recipe is **not** pre-groomed: no out-of-scope items, no steps, no
   work packages, and the criterion's *verification* is left empty — a
   criterion with no verification cannot pass `ac.falsifiable`, so the groom
   gate refuses a promoted recipe exactly as it refuses any other half-written
   recipe. `TestThePromotionWeakensNoGate` proves it by name
   (`scope.out_of_scope_declared`, `ac.falsifiable`).
   **A judgement call to check:** the item's *title* becomes the recipe's
   `intent.outcome` and the item's *body* becomes `intent.problem`; the
   scope becomes `scope.in_scope`. The old prep note's intent became the
   outcome; this maps the two text fields to the two fields they actually
   mean. An idea promotes with **no** criterion (the criteria are what the
   recipe is groomed into); a small task promotes with **exactly one**.

4. **The fence still stands, in its designed direction.** A plan stamped
   *newer* than the tool is refused (6 → refused); a plan stamped 4 is
   *older*, so the tool proceeds and re-stamps it to 5 on its next write —
   that is the upgrade path the guard was built for
   (`TestDomainVersionRefusesANewerPlan`,
   `TestDomainVersionStampsAnOlderPlan`). If you meant "a plan stamped 4 must
   be refused", that is the opposite of the guard's direction and it would be
   a change to the guard's rule, not a slice-C detail — tell me and I will do
   it.

5. **A gap I did not close, on purpose.** There is no verb that migrates a
   legacy special into the backlog: the CLI speaks the new format only, and
   the console shows the old artifacts read-only. A migration verb would move
   the command count this slice pins at 42, so I left it to you.

6. **docs/07-vocabulary.md still says "Today the backlog has two file
   types".** It is the decision of record and it is yours; I did not rewrite
   it. The current state lives in docs/03, docs/04, docs/05 and GLOSSARY.md.

7. **Two tripwires were missing and one was half-blind.** The kinds are a
   closed enum and the doc must state the set (`TestTheDocListsTheBacklogKinds`,
   both ways); an unknown `--kind` must be refused with the whole set in the
   message (`TestAnUnknownKindNamesTheKinds`); and the command-count wire
   checked the total but not the verbs-and-groups the same claim is told in —
   it now checks both halves.

Mutation checks I ran and reverted: the id normaliser made to accept only
`BL` (ten tests fail, including the recipe's); the command count moved to 41
(the wire fires); the verbs-and-groups moved to "37 verbs in 7 groups" (the
wire fires).

LOW notes: none.
