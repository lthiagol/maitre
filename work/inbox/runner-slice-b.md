# Runner — slice B report

**Commits:** 4 acceptance criteria + 1 work commit.

| Commit | What |
|---|---|
| `af6a723` | **AC-01** — the roles are runner, reviewer, orchestrator |
| `d2d862d` | **AC-02** — the tabs are Board and Sous-Chef |
| `32df9d4` | **AC-03** — the console shows one word per fact |
| `773f7f4` | **AC-04** — the tripwires still hold, and one was missing |
| `876cb5f` | work: the paths follow the repo (the repo moved; the briefs carried the old path) |

**Per-AC evidence** — stamped in each commit body:

| AC | Evidence |
|---|---|
| AC-01 | `go test ./internal/core -run 'Role\|SelfReview\|Renamed\|AreTwoAgents\|OldPlan' -v` exit 0 (8/8) |
| AC-02 | `go test ./internal/ui` exit 0; goldens regenerated; `maitre --render --plain` shows `Recipes Backlog Path Board Log Sous-Chef Settings` |
| AC-03 | `go test ./internal/ui -run 'Word\|Lane\|State\|Blocked' -v` exit 0 (3/3); the meta-line test was checked to **fail** when the lane word is put back |
| AC-04 | `go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v` exit 0 (13/13); command count still 46; `core.DomainVersion` still 4 |

**Discipline gate:** `make lint` clean · `gofmt -l .` empty · `make test` green ·
`make regen-goldens` clean · `git status` clean apart from the orchestrator's own
`work/` edits.

## The invariants held

| | |
|---|---|
| `internal/core` and `internal/cli` | **0** external modules |
| `internal/ui` | **28** external modules |
| The domain layer | **zero** UI imports |
| `kitchen/activity.json` | **33** events, untouched — the 32 historical keep `R01`, the migration is event 33, `control`, stamped 4 |
| `core.DomainVersion` | **4** — the roles are vocabulary, not a format change |

## Notes for the reviewer

1. **The WIP was kept, not reverted.** The dead lane's work was coherent and green, so
   it became AC-01's commit, with the brief's leftovers fixed on top: the CLI's
   documented examples (`internal/cli/cli.go`, `internal/cli/flags.go`), `docs/05:292`,
   `docs/04`'s role examples, `docs/03`'s role rows, and `scripts/stamp-plan`'s actor.
   The sharp edge — a review round recorded under `line-cook` cannot close its own
   recipe — **was already tested** (`internal/core/roles_test.go`,
   `TestASelfReviewSurvivesTheRename`); the brief asked for it and it is there.

2. **AC-01 renamed the internal names too, not only the labels.** `retiredAgentRoles`,
   `agentRole()`, `sameAgent()` — a retired word left in the code is drift waiting for a
   reader. `internal/core/model.go`'s `IsAgentRole` keeps its name and its signature.

3. **AC-02 renamed the internal names too.** `TabChopboard` → `TabBoard`,
   `TabExpo` → `TabSousChef`, and with them `chopTree` → `boardTree`,
   `chopBuilder` → `boardBuilder`, `expoLines` → `sousChefLines`, and the two test
   files (`chop_test.go` → `board_test.go`, `expo_test.go` → `souschef_test.go`).
   The **goldens changed in one line each** — the tab bar — and nothing else.

4. **AC-03 found one real drift, and it was not in the row.** The row and the tree rows
   already showed the plain state word. The **detail pane** showed the state word *and*
   the lane's word side by side: `closed · Served`. That is two names for one fact, so
   the lane's word is gone from the meta line.
   **Nothing was lost with it:** the lane was the only place the pane said a recipe was
   blocked, so a blocked recipe now says so on its own line, with the reason — and
   `TestABlockedRecipeStillSaysItIsOnHold` guards that.

5. **AC-04 needed a commit, so a tripwire was added.** The four tripwires already held
   and needed no change. But the rename exposed a gap: the **deny-list** had no
   tripwire tying `core.AgentRoles` to the spec's prose. Now the spec must name every
   role the code refuses and every name the rename retired — and the tripwire is proved
   able to fail, the way the derived-views tripwire is.

6. **`kitchen/project.json` still says `"name": "new-project"`.** Per the brief, it was
   **not** changed — it is the plan's data, not the code's, and outside this slice. The
   reviewer judges it.

7. **`GLOSSARY.md` was not touched.** It still maps `runner` → *line cook*, `reviewer`
   → *taster*, `orchestrator` → *expo*, and calls the tab *Expo*. The brief puts the
   glossary in slice E, so it is left for that slice — but it now contradicts the code,
   and the slice that owns it should know.

8. **Two ADRs still carry the old names** (`docs/decisions/0004-naming.md`,
   `docs/decisions/0005-vocabulary.md`) and `docs/audit.md` still calls the reviewer
   *the taster*. Those are dated records — what was decided, and what the ledger says —
   and rewriting them would rewrite history. `docs/07-vocabulary.md` is the decision of
   record and already carries the rename table.

9. **`work/briefs/review-b.md` is untracked** and was left alone: the reviewer wrote it
   while this lane worked.
