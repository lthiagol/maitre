Reviewer: slice C cycle 2 done. Verdict: ok. F-74..F-77 closed: yes. Filed 0 new
blockers (none). Mutation checks: (F-74) removed `m.legacy = msg.legacy` from
`Model.applyReload` — both `TestAnOldPlanStillRenders` and
`TestTheLiveConsoleCarriesTheLegacyGroup` fail, and the live console's Backlog tab
loses the FROM THE OLD PLAN group; reverted. (F-75) removed the legacy scan from
`reservedRecipeIDs` — `TestALegacyDanglingPointerResolvesToNothing` fails ("the next
id is RCP10, but that is the one the old plan points at"); reverted. (F-76) made
`listArg` return an empty slice for a flag that was not given — the CLI test fails
("the update dropped the tags"); and made `setList` treat an empty list as "not
given" — the core test fails ("the empty lists did not clear the fields"); both
reverted. (F-77) the code did not move — `promotedCriteria` still leaves the
verification empty and `TestThePromotionWeakensNoGate` still passes.

What I checked beyond the four mutations:
- F-74: `applyReload` is now the only place that assigns the model's plan fields
  (grep for `m.recipes =`/`m.items =`/`m.legacy =`/`m.journal =` returns only
  applyReload plus the unrelated `m.sections =` in the order handler). The test
  helper no longer copies fields by hand — `loadModel` and `theme_test.go` both send
  the reload message through the model's own `Update`.
- F-75: the format on disk is the two-word `"recipe RCP10"`, and that is reserved.
  The older `"recipe R10"` spelling reserves its number too (parseNumericID strips the
  letters). A bare `"RCP10"` is one word and is skipped by `referencedIDs` —
  deliberate, since the writer only ever wrote the two-word form; a legacy
  special's own pointer (`"prep-note N01"`) is correctly not a recipe id.
- F-76: both list fields obey one rule. `backlog update --tag x` leaves the scope
  alone; `backlog update --scope x` leaves the tags alone; an update naming neither
  leaves both. The other verbs are byte-identical to the base: `--in-scope`,
  `--out-of-scope`, `--step --covers` and `--ac add` all behave exactly as before.
  The JSON edge: no document carries `null` for a list — the backlog item's `scope`
  and `tags` are omitted when empty, and `jq` reads them as `null`, which is the
  same shape the base's `special`/`prep note` had; `emptyIfNil` keeps the in-memory
  value a list. The parser still refuses `--tag` with no value (`--tag needs a
  value`) — the chef's open question was not touched.
- F-77: the row now says the criterion makes the work **checkable** and becomes
  falsifiable when the recipe's own life writes its verification; the paragraph says
  the same in the other direction. The code did not move.

What the fixes introduced: nothing outside the four findings' own code. The gate
table is byte-identical, the plan format facts did not move, the command count is
still 42, `core.DomainVersion` is still 5 and the recorded fingerprint is unchanged.
The only addition beyond the findings is three lines in `internal/ui/detail.go`
showing "Promoted to …" for a legacy artifact — beyond the finding, but honest (the
pointer is read now) and it puts the legacy pane on the same footing as the current
one, so I judge it belongs.

Scope and invariants: `internal/core` and `internal/cli` stay stdlib-only,
`internal/ui` keeps 28 external modules, the domain layer has zero UI imports.
`make lint` clean, `gofmt -l .` empty, `make test` green, `go test -race` green on
all three packages, `go test -count=2 ./...` green, and the demo and chaos canaries
both write a plan that renders.

LOW notes (not filed, per the brief):
- A legacy plan's file that is not valid JSON now fails recipe creation as well as
  the render — the reservation scan reads the legacy directories. The base behaved
  the same way for a corrupt prep note, so this is the old shape, not a new one.
- The runner's own note 2 (the brief's "the second out-of-scope item" reading) is
  still open. I read the brief the way the runner did — the promotion writes no
  out-of-scope item, and the recipe's own life writes them — because the brief's
  sentence lists the out-of-scope items among the things that are "the recipe's own
  life, exactly as for any recipe". It is the chef's call, and I am not filing it.
- The runner's note 1 (a list cannot be emptied from the command line) is also
  still open and is the chef's; the parser's rule was deliberately left alone.
