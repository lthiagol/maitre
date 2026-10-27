Runner: slice C fix 1 done. Commits: 4 (one per finding). Fixed: F-74, F-75,
F-76, F-77. Pre-notify gate: pass (make lint clean, gofmt -l . empty, make test
green, git status clean). The gate table, the plan format and the command
count (42) did not move — the tripwires for them still pass.

**F-74 — the live console shows an old plan's specials again.** The live
reload path dropped `m.legacy`; the render path assigned it. Now
`Model.applyReload` puts a reloaded plan into the model and **both** paths go
through it — one place knows the model's shape, so a field the render carries
and the live path forgets cannot happen. The test drives the real path now:
`loadModel` sends the reload message through the model's own `Update` instead
of copying the fields by hand (the copy was the bug), and
`TestTheLiveConsoleCarriesTheLegacyGroup` checks the message *and* the model.
Proof by mutation: removing `m.legacy = msg.legacy` from `applyReload` fails
both tests; reverted.

**F-75 — a dangling pointer reserves its id again.** `core.LegacyItem` carries
`PromotedTo` and `reservedRecipeIDs` scans the old plan's pointers as well as
the backlog's. The test creates nine recipes so the numbering reaches `RCP10`
and then checks that `RCP10` is not handed out. Proof by mutation: dropping the
legacy scan fails the test; reverted. (Note: the reservation moves the whole
sequence — with a dangling `RCP10` the first new recipe is `RCP11`. That is the
pre-existing rule, unchanged: the reserved id is skipped, not filled.)

**F-76 — a list flag that is not given leaves its field alone.** The parser
now knows which list flags were given, and `flagSet.listArg` returns **nil**
for a flag that was not given. The operation reads nil as "not given" and a
list as "replace the list it names", through one helper (`setList`) for the
scope and the tags alike. An empty list that *was* given clears the field. The
other verbs are untouched: `StringList` still returns an empty slice, so
`--in-scope`, `--out-of-scope`, `--ac`, `--ac-verify`, `--depends-on` and the
rest behave exactly as before. No null reaches a JSON document: the domain
never stores a nil list (`emptyIfNil`), and the item's list fields are omitted
when empty, exactly as they were before this fix. Proof by mutation: making
`listArg` return an empty slice for a flag that was not given reproduces the
bug exactly ("the update dropped the tags"); reverted. A second mutation — a
`setList` that treats an empty list as "not given" — fails the core test.

**F-77 — the spec says what is true.** § 1.1's row now says the criterion
makes the work **checkable** and becomes **falsifiable** when the recipe's own
life writes its verification; the paragraph says the same in the other
direction ("checkable, not yet falsifiable"). The code did not move.

Notes for the reviewer:

1. **The list-clearing question is open, and I did not decide it.** The parser
   refuses a list flag with no values (`--tag needs a value`), so a list can be
   *replaced* but never *emptied* from the command line. The domain can clear a
   list — an empty list is a list that was given — but the CLI has no way to
   say it. Changing the parser's rule would touch every list flag, so I left it
   alone and stated the current behaviour in docs/04. The question is the
   chef's.

2. **One reading I did not change, and you may want to put to the chef.** The
   slice C brief said the promotion does not pre-groom "the second out-of-scope
   item" — which can be read as *the promotion writes one and the recipe's own
   life writes the second*. I read it as *the recipe's own life writes them*,
   and the code writes **no** out-of-scope items; docs/03 § 1.1 says so
   plainly, and the reviewer filed no finding about it. If the first reading is
   right, the promotion must write one out-of-scope item — a behaviour change,
   not a documentation fix, so I did not make it.

3. The `docs/audit.md` ledger and `docs/07-vocabulary.md` are untouched — they
   are yours and the orchestrator's. The slice C brief's falsifiable phrase was
   already corrected in work/briefs/slice-c.md before this cycle.

LOW notes: none.
