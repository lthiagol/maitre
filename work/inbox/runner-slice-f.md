Runner: slice F done. Commits: 6 (`a0c4524` F-1, `88d3e23` F-2, `19a2284` F-3,
`2d5c546` F-4, `52787a4` F-5, `761bd55` F-6). Per-AC evidence: AC-01..AC-06 stamped in
each commit's body. Pre-notify gate: pass (`make lint && gofmt -l . && make test &&
make regen-goldens && git status --short` — green; the tree clean apart from the other
session's `logo/`).

## What the slice did

| AC | What |
|---|---|
| AC-01 | `internal/core/legacy.go` is gone; `internal/core/retired.go` is new. A plan that carries a retired artifact — `kitchen/specials/`, `kitchen/prep-notes/`, a criterion's `proofs`, a recipe's old `evidence` stamp, a live record under a retired id form — is **not supported**, and the reader says so plainly. The console's third group ("FROM THE OLD PLAN") is gone. |
| AC-02 | `retiredAgentRoles` and `agentRole` are gone. The retired names are **on the deny-list**, so the refusal survives; the resolution does not. The wire now reads the **one** list (`core.AgentRoles`), and the drift test proves it fires on a spec that drops a retired name. |
| AC-03 | The digest hashes a **canonical payload** the code owns (`approvalPayload`, `approvalIntent`, `approvalScope`, `approvalCriterion`), with fixed field names: a rename of a record's field does not move the digest. The candidate set (`legacyFingerprints` and friends) and the re-stamping are gone. A golden digest pins the payload's shape. |
| AC-04 | The fence is **one rule**: a tool writes a plan only if the plan's version **is** the tool's version. No upgrade path. The refusal says what the chef can do and promises no migration. The prose in `docs/03` § 6.2 and in `core.DomainVersion` says the two jobs: the shape is the code-side tripwire, the version is the plan-side fact. |
| AC-05 | The seeds, the demo, the fixtures and the tests are in the current format. The repository's own plan's live records are in the current format; **the journal's history is untouched** (append-only). |
| AC-06 | The tripwires hold, and a new wire holds "one format" (below). |

## Tests deleted, and tests changed

**Deleted** — each one asserted a tolerance that no longer exists:

- `internal/core/legacy_test.go` (the old field names) and `internal/core/legacy.go`.
- `TestALegacyDanglingPointerResolvesToNothing` (the old backlog directories).
- `TestAnOldPlanStillRenders` and `TestTheLiveConsoleCarriesTheLegacyGroup` (the console's third group).
- `TestANormalisedLegacyApprovalIsRestamped` and `TestALegacyApprovalIsRestampedNotRefused` (the candidate set).
- `TestASelfReviewSurvivesTheRename` (the retired-name resolution).
- `TestDomainVersionStampsAnOlderPlan` (the upgrade path).

**Changed** — the subject kept, the tolerance gone:

- `TestTheRetiredFieldNamesAreRefused` (new) and `TestAPlanWithRetiredArtifactsIsNotSupported` (new) assert the **refusal**.
- `TestLoadRecipeTolerantOfIDForm` and `TestListRecipesSortsNumerically`: the records use the current id forms; the *lookup* is still tolerant.
- `TestTheRunnerAndTheReviewerAreTwoAgents`: the current roles.
- `TestARenamedFieldIsNotATamper` and `TestAFormChangeIsNotATamper`: the canonical payload, no candidates.
- `TestDomainVersionRefusesAnOlderPlan` (new) replaces the upgrade test; `TestTheRepoPlanIsWrittenUnderThisToolsVersion` replaces `TestTheRepoPlanStaysWritableUnderANewerVersion`.

## Notes for the reviewer

1. **`core.DomainVersion` stays 6.** No format fact changed — the ids, the directories, the root files and the field names on disk are the same — so the tripwire's shape did not move and the version did not either. The brief's invariant holds. The *rules* around the fence changed (the fence's direction, the reader, the digest's payload), and the prose says so.
2. **The repository's own plan was touched twice, both in F-3:** its approval's digest was re-stamped to the canonical digest (the content did not change — the encoding did; two `null` lists also became `[]`, the plan's own rule), and the journal got **one appended** `control` event recording the digest encoding's move. The journal's history (37 events, the old ids, the old names) was not rewritten.
3. **A lookup is tolerant; a record is not.** `RCP01`, `rcp01`, `RCP1` and `1` name the same recipe — a name resolves to a record, and a name is not a record — and the journal's old ids still resolve. The *record* on disk must be current-format, and a retired one is refused.
4. **A hole the ruling accepts:** with the resolution gone, the domain cannot tell that `line-cook` and `runner` are one agent, so a recipe claimed under the old name and reviewed under the new one no longer reads as self-reviewed. The deny-list refusal is unaffected.
5. **The retired detection's edges**, so the reviewer can check them: an *empty* `specials/` is untidy, not an artifact (the plan still reads); a backlog item with **no `kind`** is of the retired shape and is refused; the recipe's own top-level `evidence` is the old stamp, and a criterion's `evidence` is not.
6. **`docs/07` still describes the third group and the one-sided guard** ("a plan written under the newer version refuses a stale tool"). The ruling replaced both. The brief's scope table puts `docs/07` out of touch, so I left it — say the word and I will add the line. `docs/03`, `AGENTS.md` and the code's prose are current.
7. **A new wire** for "one format": `TestTheRetiredNamesAreOnlyWhereTheyBelong` scans the package's **code** (not the prose) for the retired names; they may live only in `retired.go` (the detection) and `model.go` (the deny-list). It was proved on real drift: a file that reads `proofs` fails it.
8. `internal/core` and `internal/cli` are still stdlib-only; `internal/ui` still pulls **28** external modules.
