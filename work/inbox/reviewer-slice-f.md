Reviewer: slice F review done. Verdict: changes-needed. Closed: none — the slice filed no
finding of its own, so the review had nothing to close. Filed 4 new blockers (F-87..F-90).
Mutation checks (every mutation reverted; the tree is as the runner left it, plus this
ledger): renaming the criterion's record field (`evidence` → `proofs` in `model.go`) left the
digest where it was — `TestARenamedFieldIsNotATamper` and
`TestTheRepoPlansApprovalDigestMatchesItsContent` still pass; dropping the intent from the
canonical payload made `TestTamperedFrozenSectionsAreCaught` fail (the content fence is not
vacuous); making `checkDomainVersion` accept an older plan made
`TestDomainVersionRefusesAnOlderPlan` fail (the fence's new rule is tested); putting
`"specials"` in `internal/core/store.go` made the new wire fail (it can fail), while the same
token in `internal/ui/views.go`, the same token in `internal/core/model.go`, and the literals
`"R01"` and `LegacyThing` in `internal/core/store.go` all passed it (F-90). Two probes in a
temporary test file, since removed: an approval by `line-cook`/`taster`/`expo`/`coordinator`/
`Line Cook` is refused ("is an agent role") — AC-02's proof; a recipe claimed by `line-cook` and
reviewed by `runner` is **not** self-reviewed and passes `review.reviewer_independent` (F-87);
and a plan whose approval digest was written by the previous build passes the version guard
(6 == 6) and is then refused with `immutability.respected (frozen sections changed since
approval)` — the tool calls a plan nobody touched tampered (F-89). The eight deleted tests:
`legacy_test.go`'s five (the retired field names) — the live counterparts hold, the writer cannot
write a retired field name (`internal/core/fields_test.go`) and the reader refuses one
(`TestTheRetiredFieldNamesAreRefused`); `TestANormalisedLegacyApprovalIsRestamped` and
`TestALegacyApprovalIsRestampedNotRefused` — the content fence holds
(`TestTamperedFrozenSectionsAreCaught`), the form change holds (`TestAFormChangeIsNotATamper`),
a rename holds (`TestARenamedFieldIsNotATamper`), but "an approval written under an older
encoding is not tampering" went with them, and that is F-89;
`TestALegacyDanglingPointerResolvesToNothing` — held by
`TestADanglingPointerResolvesToNothing`; `TestASelfReviewSurvivesTheRename` — the same-name
property is held by `internal/core/selfreview_test.go` (`TestSelfReviewedIsOneAnswer`,
`TestTheGateAndTheTriageAgree`), the retired-spelling property went with it, and that is F-87;
`TestAnOldPlanStillRenders` and `TestTheLiveConsoleCarriesTheLegacyGroup` — held by
`TestAPlanWithRetiredArtifactsIsNotSupported` (the console refuses the plan and says so, on the
live path too); `TestDomainVersionStampsAnOlderPlan` — replaced by
`TestDomainVersionRefusesAnOlderPlan`, and `TestTheRepoPlanStaysWritableUnderANewerVersion` by
`TestTheRepoPlanIsWrittenUnderThisToolsVersion`. The journal: 37 events unchanged, the 38th is
the migration (`git diff` on `kitchen/activity.json` shows one appended event), and the
journal's old ids still resolve (`TestTheJournalsOldIdStillResolves`, kept). The judgement the
brief asked for: a retired id form resolving to a live record is **the same answer written
differently**, not a second answer — the number is the identity, the directory carries the kind,
and there is one record per number, refused when it is in a retired form; so the lookup's
tolerance stands. The wire's token list should nevertheless name the retired id *literals*
(`R01`, `S01`, `N01`): the code never needs to write them, and a mapping that names them is the
affordance the ruling removed. LOW notes: `internal/core/activity.go:81` still says a
version-zero plan is "an older plan, which a newer tool may write to" — the rule is now the
opposite, and only a plan with **no** journal is writable; the retired check reads the id inside
the document, not the file's name, so a recipe file named `R01-*.json` whose content says
`RCP01` is still read; the reader refuses a plan carrying a retired artifact but the writer does
not, so `maitre-cli backlog create` still writes to such a plan; `GLOSSARY.md:97` names
"coordinator" in prose outside an `_Avoid_` line (pre-existing, F-83's shape); the retired
recipe-level stamp's name (`evidence`) cannot be a token in the new wire — the criterion's
current field shares the word — so that one retired name needs another instrument. The mascot:
`logo/` is untracked and nothing of it is in the slice's commits; `docs/audit.md` was untouched
by the runner.
