# The runner's lane — slice B, captured 2026-10-27

Dumped from pane `w1Z:p2` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that wrote slice B: AC-01 to AC-04, then the two fix cycles (F-72).
The review that judged it is in [slice-b-reviewer.md](slice-b-reviewer.md).

```

     + Thought · 1.3s

  ┃
  ┃  $ grep -n "func readDoc" -A 20 internal/core/docs_test.go; grep -n "TestTheDocListsTheDerivedViews" -A 45 internal/core/docs_test.go | head -60
  ┃
  ┃  (64 earlier lines)
  ┃  819-
  ┃  820-  cases := []struct {
  ┃  821-    name  string
  ┃  822-    doc   string
  ┃  823-    check func(string) error
  ┃  824-  }{
  ┃  825-    {"a chef move removed", strings.Replace(doc, " · `triage`", "", 1), checkDerivedViews},
  ┃

     + Thought · 1.7s

  ┃
  ┃  $ tail -40 internal/core/docs_test.go; wc -l internal/core/docs_test.go
  ┃
  ┃  (34 earlier lines)
  ┃      }
  ┃    }
  ┃    if strings.Contains(doc, "cancel.reason_present") {
  ┃      t.Error("the doc names cancel.reason_present as a gate — the reason is checked inline")
  ┃    }
  ┃  }
  ┃      1545 internal/core/docs_test.go
  ┃

  ┃
  ┃  $ cat >> internal/core/docs_test.go <<'EOF'
  ┃  …
  ┃
  ┃  (9 earlier lines)
  ┃  === RUN   TestTheRenamedRolesAreTheAgentRoles
  ┃  --- PASS: TestTheRenamedRolesAreTheAgentRoles (0.00s)
  ┃  === RUN   TestThePlanJournalsTheNewRoles
  ┃  --- PASS: TestThePlanJournalsTheNewRoles (0.02s)
  ┃  PASS
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.219s
  ┃

     + Thought · 2.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
  ┃  import re…
  ┃
  ┃  # github.com/lthiagol/maitre/internal/core [github.com/lthiagol/maitre/internal/core.test]
  ┃  internal/core/docs_test.go:1556:11: undefined: fmt
  ┃  internal/core/docs_test.go:1561:11: undefined: fmt
  ┃  FAIL  github.com/lthiagol/maitre/internal/core [build failed]
  ┃  FAIL
  ┃

  ┃
  ┃  $ sed -n '1,12p' internal/core/docs_test.go
  ┃
  ┃  (5 earlier lines)
  ┃    "reflect"
  ┃    "regexp"
  ┃    "strings"
  ┃    "testing"
  ┃  )
  ┃
  ┃  // readDoc reads one of the project's design docs.
  ┃

  ┃
  ┃  ← Edit internal/core/docs_test.go
  ┃
  ┃    1   package core                                                                                                                          1   package core
  ┃    2                                                                                                                                         2
  ┃    3   import (                                                                                                                              3   import (
  ┃                                                                                                                                              4 +   "fmt"
  ┃    4     "os"                                                                                                                                5     "os"
  ┃    5     "path/filepath"                                                                                                                     6     "path/filepath"
  ┃    6     "reflect"                                                                                                                           7     "reflect"
  ┃    7     "regexp"                                                                                                                            8     "regexp"
  ┃

  ┃
  ┃  $ gofmt -w internal/core/docs_test.go && go test ./internal/core -run 'DenyList' -v 2>&1 | head -20
  ┃
  ┃  === RUN   TestTheDocDoesNotOverstateTheDenyList
  ┃  --- PASS: TestTheDocDoesNotOverstateTheDenyList (0.00s)
  ┃  === RUN   TestTheDenyListTripwireFailsOnRealDrift
  ┃  --- PASS: TestTheDenyListTripwireFailsOnRealDrift (0.00s)
  ┃  PASS
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.204s
  ┃

     + Thought · 505ms

  ┃
  ┃  $ go test ./internal/core -run 'Doc|DenyList|Roles' -v 2>&1 | grep -E "^(--- |ok |FAIL)"
  ┃
  ┃  (7 earlier lines)
  ┃  --- PASS: TestTheDenyListTripwireFailsOnRealDrift (0.00s)
  ┃  --- PASS: TestTheDocDescribesWhatTheCodeDoes (0.00s)
  ┃  --- PASS: TestTheDocStatesThePlanFormatVersion (0.00s)
  ┃  --- PASS: TestTheDocNamesThePlanFormatTheCodeHas (0.00s)
  ┃  --- PASS: TestTheRenamedRolesAreTheAgentRoles (0.00s)
  ┃  --- PASS: TestThePlanJournalsTheNewRoles (0.03s)
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.206s
  ┃

     + Thought · 537ms

  ┃
  ┃  $ grep -rn "DomainVersion = \|DomainVersion=" internal/core/*.go | grep -v _test | head; go run ./cmd/maitre-cli gates 2>/dev/null | head -3
  ┃
  ┃  internal/core/activity.go:139:  e.DomainVersion = DomainVersion
  ┃  internal/core/activity.go:141:    e.DomainVersion = PlanDomainVersion(j)
  ┃  internal/core/lifecycle.go:179:const DomainVersion = 4
  ┃  {
  ┃    "gates": [
  ┃      "abort.recorded",
  ┃

  ┃
  ┃  $ make lint && gofmt -l . ; go test ./... 2>&1 | tail -6 && make regen-goldens && git status --short
  ┃
  ┃  (7 earlier lines)
  ┃  go test ./internal/ui -run Golden -update
  ┃  ok    github.com/lthiagol/maitre/internal/ui  4.730s
  ┃   M internal/core/docs_test.go
  ┃   M work/README.md
  ┃   M work/briefs/slice-b-resume.md
  ┃   M work/handover.md
  ┃  ?? work/briefs/review-b.md
  ┃

  ┃
  ┃  $ git add internal/core/docs_test.go && git commit -q -F - <<'EOF'
  ┃  AC-04: the tripwires still hold — and one was missing…
  ┃
  ┃  773f7f4 AC-04: the tripwires still hold — and one was missing
  ┃  32df9d4 AC-03: the console shows one word per fact
  ┃  d2d862d AC-02: the tabs are Board and Sous-Chef
  ┃  876cb5f work: the paths follow the repo
  ┃  af6a723 AC-01: the roles are runner, reviewer, orchestrator
  ┃

     + Thought · 1.0s

  ┃
  ┃  $ grep -rn "Chopboard\|Expo\|chopboard" internal/ cmd/ scripts/ docs/05-human-interface.md docs/06-proposals.md docs/03-domain-model.md docs/04-operations.md AGENTS.md 2>/dev/null | grep -v "exported\|exposes" | head -20; echo "--- roles ---"; grep -rn "line-cook\|line cook\|t
  ┃  aster\|\"expo\"\|coordinator" internal/ cmd/ scripts/ 2>/dev/null | head -20
  ┃
  ┃  (20 earlier lines)
  ┃  internal/core/model.go:570:  "line-cook":   "runner",
  ┃  internal/core/model.go:571:  "taster":      "reviewer",
  ┃  internal/core/model.go:572:  "expo":        "orchestrator",
  ┃  internal/core/model.go:573:  "coordinator": "sous-chef",
  ┃  internal/core/model.go:580:// and "line-cook" are the same name — and "line-cook" is the runner.
  ┃  internal/core/views.go:443:// The names are compared as written and as roles: `line-cook` and `runner` are one
  ┃

     + Thought · 1.4s

  ┃
  ┃  $ git status --short kitchen/; git diff --stat HEAD~4..HEAD; echo "--- deps ---"; go list -deps -f '{{if not .Standard}}{{.Module.Path}}{{end}}' ./internal/ui | sort -u | grep -v lthiagol/maitre | wc -l; echo "--- core ui imports ---"; grep -rn "internal/ui" internal/core/ | w
  ┃  c -l
  ┃
  ┃  (17 earlier lines)
  ┃   work/handover.md                               |   2 +-
  ┃   18 files changed, 327 insertions(+), 139 deletions(-)
  ┃  --- deps ---
  ┃        28
  ┃  --- core ui imports ---
  ┃         9
  ┃

     + Thought · 266ms

  ┃
  ┃  $ grep -rn "internal/ui" internal/core/
  ┃
  ┃  (2 earlier lines)
  ┃  internal/core/chaos.go:62:      {Description: "the columns still line up", Verification: "go test ./internal/ui -run Width"},
  ┃  internal/core/chaos.go:79:      {Description: "the row still lines up", Verification: "go test ./internal/ui -run Short"},
  ┃  internal/core/chaos.go:171:      {Description: "the notes are all listed", Verification: "go test ./internal/ui -run Findings"},
  ┃  internal/core/chaos.go:222:      {Description: "the note reads", Verification: "go test ./internal/ui -run Deviation"},
  ┃  internal/core/seed.go:71:      {Description: "the console offers sign off", Verification: "go test ./internal/ui -run SignOff"},
  ┃  internal/core/seed.go:102:      {Description: "the filter narrows the queue", Verification: "go test ./internal/ui -run Filter"},
  ┃  internal/core/seed.go:103:      {Description: "findings close without a restart", Verification: "go test ./internal/ui -run Close"},
  ┃

     + Thought · 779ms

  ┃
  ┃  $ go list -deps ./internal/core | grep "maitre/internal" ; echo "(no ui import above = good)"; make lint && gofmt -l . ; make test 2>&1 | tail -6 && make regen-goldens && git status --short
  ┃
  ┃  (8 earlier lines)
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃  go test ./internal/ui -run Golden -update
  ┃  ok    github.com/lthiagol/maitre/internal/ui  7.422s
  ┃   M work/README.md
  ┃   M work/briefs/slice-b-resume.md
  ┃   M work/handover.md
  ┃  ?? work/briefs/review-b.md
  ┃

     + Thought · 1.3s

  ┃
  ┃  $ which herdr; echo "ORCH_PANE=${ORCH_PANE:-unset}"; echo "HERDR_PANE_ID=${HERDR_PANE_ID:-unset}"; cat work/inbox/README.md
  ┃
  ┃  (10 earlier lines)
  ┃
  ┃  | File | Written by |
  ┃  |---|---|
  ┃  | `runner-<slice>.md` | the runner, when a slice is done |
  ┃  | `reviewer-<slice>-<cycle>.md` | the reviewer, when a cycle is done |
  ┃
  ┃  The orchestrator checks this directory when a lane goes quiet.
  ┃

     + Thought · 748ms

  ┃
  ┃  $ herdr --help 2>&1 | head -40
  ┃
  ┃  (33 earlier lines)
  ┃    herdr server reload-config       Reload config.toml in the running server
  ┃    herdr config reset-keys          Back up config.toml and remove custom keybindings
  ┃    herdr channel <subcommand>       Manage the stable or preview update channel
  ┃    herdr machine <subcommand>       Manage saved SSH machines
  ┃    herdr api <subcommand>           Inspect socket API metadata and live runtime state
  ┃    herdr workspace <subcommand>     Workspace helpers over the socket API
  ┃    herdr worktree <subcommand>      Git worktree helpers over the socket API
  ┃

  ┃
  ┃  $ herdr pane list 2>&1 | head -30
  ┃
  ┃  (1 earlier line)
  ┃  status":"unknown","cwd":"/Users/thiago/Code/Moises/moises-training-troubleshoot/training-prod","focused":false,"foreground_cwd":"/Users/thiago/Code/Moises/moises-training-troubleshoot/training-prod","pane_id":"w1W:p3","revision":2,"scroll":{"max_offset_from_bottom":0,
  ┃  "offset_from_bottom":0,"viewport_rows":86},"tab_id":"w1W:t3","terminal_id":"term_65ec60d16b89d9","terminal_title":"../scripts/change-gate.py tui","terminal_title_stripped":"../scripts/change-gate.py tui","workspace_id":"w1W"},{"agent":"cursor","agent_session":{"agent":
  ┃  "cursor","kind":"id","source":"herdr:cursor","value":"0f6926a0-e472-4adb-8845-4e6126832dfe"},"agent_status":"idle","cwd":"/Users/thiago/Code/Pessoal/mcode-herder","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/mcode-herder","pane_id":"w1X:p1","revision":14,
  ┃  "scroll":{"max_offset_from_bottom":6490,"offset_from_bottom":0,"viewport_rows":37},"tab_id":"w1X:t1","terminal_id":"term_65ed43b3b35141a","terminal_title":"Code Integration Analysis","terminal_title_stripped":"Code Integration Analysis","workspace_id":"w1X"},{"agent":
  ┃  "opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e8138fe5effeUJvY7OF96j92fE"},"agent_status":"idle","cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:
  ┃  p1","revision":16,"scroll":{"max_offset_from_bottom":0,"offset_from_bottom":0,"viewport_rows":86},"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722","terminal_title":"OC | Maitre vocabulary rename handover","terminal_title_stripped":"OC | Maitre vocabulary rename
  ┃  handover","workspace_id":"w1Z"},{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e81361ce3ffeHKJiGcrf284UWh"},"agent_status":"working","cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/
  ┃  thiago/Code/Pessoal/maitre","pane_id":"w1Z:p2","revision":4,"scroll":{"max_offset_from_bottom":0,"offset_from_bottom":0,"viewport_rows":86},"tab_id":"w1Z:t2","terminal_id":"term_65ed74394378f23","terminal_title":"OC | Slice B roles and tabs","terminal_title_stripped":"OC |
  ┃  Slice B roles and tabs","workspace_id":"w1Z"},{"agent_status":"unknown","cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p3","revision":7,"scroll":{"max_offset_from_bottom":38,"offset_from_bottom":0,
  ┃  "viewport_rows":86},"tab_id":"w1Z:t3","terminal_id":"term_65ed74a08674624","terminal_title":"thiago@donatello:~/Code/Pessoal/maitre","terminal_title_stripped":"thiago@donatello:~/Code/Pessoal/maitre","workspace_id":"w1Z"}],"type":"pane_list"}}
  ┃

     + Thought · 5.2s

  ┃
  ┃  # Wrote work/inbox/runner-slice-b.md
  ┃
  ┃    1 # Runner — slice B report
  ┃    2
  ┃    3 **Commits:** 4 acceptance criteria + 1 work commit.
  ┃    4
  ┃    5 | Commit | What |
  ┃    6 |---|---|
  ┃    7 | `af6a723` | **AC-01** — the roles are runner, reviewer, orchestrator |
  ┃    8 | `d2d862d` | **AC-02** — the tabs are Board and Sous-Chef |
  ┃    9 | `32df9d4` | **AC-03** — the console shows one word per fact |
  ┃   10 | `773f7f4` | **AC-04** — the tripwires still hold, and one was missing |
  ┃   11 | `876cb5f` | work: the paths follow the repo (the repo moved; the briefs carried the old path) |
  ┃   12
  ┃   13 **Per-AC evidence** — stamped in each commit body:
  ┃   14
  ┃   15 | AC | Evidence |
  ┃   16 |---|---|
  ┃   17 | AC-01 | `go test ./internal/core -run 'Role\|SelfReview\|Renamed\|AreTwoAgents\|OldPlan' -v` exit 0 (8/8) |
  ┃   18 | AC-02 | `go test ./internal/ui` exit 0; goldens regenerated; `maitre --render --plain` shows `Recipes Backlog Path Board Log Sous-Chef Settings` |
  ┃   19 | AC-03 | `go test ./internal/ui -run 'Word\|Lane\|State\|Blocked' -v` exit 0 (3/3); the meta-line test was checked to **fail** when the lane word is put back |
  ┃   20 | AC-04 | `go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v` exit 0 (13/13); command count still 46; `core.DomainVersion` still 4 |
  ┃   21
  ┃   22 **Discipline gate:** `make lint` clean · `gofmt -l .` empty · `make test` green ·
  ┃   23 `make regen-goldens` clean · `git status` clean apart from the orchestrator's own
  ┃   24 `work/` edits.
  ┃   25
  ┃   26 ## The invariants held
  ┃   27
  ┃   28 | | |
  ┃   29 |---|---|
  ┃   30 | `internal/core` and `internal/cli` | **0** external modules |
  ┃   31 | `internal/ui` | **28** external modules |
  ┃   32 | The domain layer | **zero** UI imports |
  ┃   33 | `kitchen/activity.json` | **33** events, untouched — the 32 historical keep `R01`, the migration is event 33, `control`, stamped 4 |
  ┃   34 | `core.DomainVersion` | **4** — the roles are vocabulary, not a format change |
  ┃   35
  ┃   36 ## Notes for the reviewer
  ┃   37
  ┃   38 1. **The WIP was kept, not reverted.** The dead lane's work was coherent and green, so
  ┃   39    it became AC-01's commit, with the brief's leftovers fixed on top: the CLI's
  ┃   40    documented examples (`internal/cli/cli.go`, `internal/cli/flags.go`), `docs/05:292`,
  ┃   41    `docs/04`'s role examples, `docs/03`'s role rows, and `scripts/stamp-plan`'s actor.
  ┃   42    The sharp edge — a review round recorded under `line-cook` cannot close its own
  ┃   43    recipe — **was already tested** (`internal/core/roles_test.go`,
  ┃   44    `TestASelfReviewSurvivesTheRename`); the brief asked for it and it is there.
  ┃   45
  ┃   46 2. **AC-01 renamed the internal names too, not only the labels.** `retiredAgentRoles`,
  ┃   47    `agentRole()`, `sameAgent()` — a retired word left in the code is drift waiting for a
  ┃   48    reader. `internal/core/model.go`'s `IsAgentRole` keeps its name and its signature.
  ┃   49
  ┃   50 3. **AC-02 renamed the internal names too.** `TabChopboard` → `TabBoard`,
  ┃   51    `TabExpo` → `TabSousChef`, and with them `chopTree` → `boardTree`,
  ┃   52    `chopBuilder` → `boardBuilder`, `expoLines` → `sousChefLines`, and the two test
  ┃   53    files (`chop_test.go` → `board_test.go`, `expo_test.go` → `souschef_test.go`).
  ┃   54    The **goldens changed in one line each** — the tab bar — and nothing else.
  ┃   55
  ┃   56 4. **AC-03 found one real drift, and it was not in the row.** The row and the tree rows
  ┃   57    already showed the plain state word. The **detail pane** showed the state word *and*
  ┃   58    the lane's word side by side: `closed · Served`. That is two names for one fact, so
  ┃   59    the lane's word is gone from the meta line.
  ┃   60    **Nothing was lost with it:** the lane was the only place the pane said a recipe was
  ┃   61    blocked, so a blocked recipe now says so on its own line, with the reason — and
  ┃   62    `TestABlockedRecipeStillSaysItIsOnHold` guards that.
  ┃   63
  ┃   64 5. **AC-04 needed a commit, so a tripwire was added.** The four tripwires already held
  ┃   65    and needed no change. But the rename exposed a gap: the **deny-list** had no
  ┃   66    tripwire tying `core.AgentRoles` to the spec's prose. Now the spec must name every
  ┃   67    role the code refuses and every name the rename retired — and the tripwire is proved
  ┃   68    able to fail, the way the derived-views tripwire is.
  ┃   69
  ┃   70 6. **`kitchen/project.json` still says `"name": "new-project"`.** Per the brief, it was
  ┃   71    **not** changed — it is the plan's data, not the code's, and outside this slice. The
  ┃   72    reviewer judges it.
  ┃   73
  ┃   74 7. **`GLOSSARY.md` was not touched.** It still maps `runner` → *line cook*, `reviewer`
  ┃   75    → *taster*, `orchestrator` → *expo*, and calls the tab *Expo*. The brief puts the
  ┃   76    glossary in slice E, so it is left for that slice — but it now contradicts the code,
  ┃   77    and the slice that owns it should know.
  ┃   78
  ┃   79 8. **Two ADRs still carry the old names** (`docs/decisions/0004-naming.md`,
  ┃   80    `docs/decisions/0005-vocabulary.md`) and `docs/audit.md` still calls the reviewer
  ┃   81    *the taster*. Those are dated records — what was decided, and what the ledger says —
  ┃   82    and rewriting them would rewrite history. `docs/07-vocabulary.md` is the decision of
  ┃   83    record and already carries the rename table.
  ┃   84
  ┃   85 9. **`work/briefs/review-b.md` is untracked** and was left alone: the reviewer wrote it
  ┃   86    while this lane worked.
  ┃   87
  ┃

  ┃
  ┃  $ git add work/inbox && git commit -q -m "inbox: runner reports slice-b" && git log --oneline -3 && git status --short
  ┃
  ┃  a3e1b9a inbox: runner reports slice-b
  ┃  773f7f4 AC-04: the tripwires still hold — and one was missing
  ┃  32df9d4 AC-03: the console shows one word per fact
  ┃   M work/README.md
  ┃   M work/briefs/slice-b-resume.md
  ┃   M work/handover.md
  ┃  ?? work/briefs/review-b.md
  ┃

  ┃
  ┃  $ herdr agent prompt "w1Z:p1" "Runner: slice B done. Commits: 4 AC + 1 work + 1 inbox. Per-AC evidence: AC-01..AC-04 stamped. Pre-notify gate: pass. Notes: the WIP was kept and became AC-01 (the retired-name self-review test was already there). AC-02 renamed the internal names
  ┃   too (TabBoard/TabSousChef, chopTree->boardTree); the goldens changed one line each. AC-03 found the real drift in the detail pane, not the row: it showed 'closed . Served' — two names for one fact; the lane's word is gone and a blocked recipe now says 'On hold: <reason…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e8138fe5effeUJvY7OF96j92fE"},"agent_status":"idle","completion_seq":341,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":
  ┃  false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":16,"screen_detection_skipped":true,"state_change_seq":341,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722","terminal_title":"OC | Maitre vocabulary rename handover",
  ┃  "terminal_title_stripped":"OC | Maitre vocabulary rename handover","workspace_id":"w1Z"},"type":"agent_prompted"}}
  ┃

     + Thought · 287ms

     Slice B is complete. Summary:

     What landed

     ┌─────────┬───────┬─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Commit  │ AC    │ What                                                                                                                                                                                                                                                            │
     ├─────────┼───────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ af6a723 │ AC-01 │ The roles are runner, reviewer, orchestrator. The dead lane's WIP was kept (it was coherent and green) and became this commit, with the brief's leftovers fixed on top: CLI examples, docs/03/04/05, scripts/stamp-plan. The brief's sharp edge — a review      │
     │         │       │ under line-cook cannot close its own recipe — was already tested.                                                                                                                                                                                               │
     ├─────────┼───────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ d2d862d │ AC-02 │ TabChopboard → TabBoard, TabExpo → TabSousChef, plus the internal names (chopTree → boardTree, expoLines → sousChefLines, both test files renamed). The Sous-Chef keeps its honest content.                                                                     │
     ├─────────┼───────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 32df9d4 │ AC-03 │ One word per fact.                                                                                                                                                                                                                                              │
     ├─────────┼───────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 773f7f4 │ AC-04 │ The tripwires hold; a missing one was added.                                                                                                                                                                                                                    │
     └─────────┴───────┴─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

     The two findings worth your attention

     AC-03's drift was not in the row — the row and the tree rows already showed the plain state word. The detail pane showed the state word and the lane's word side by side: closed · Served. Two names for one fact. The lane's word is gone. Nothing was lost: the lane was the only
     place the pane said a recipe was blocked, so a blocked recipe now says On hold: <reason> on its own line, and a test guards it.

     AC-04 needed a commit, so a tripwire was added. The four tripwires already held. But the rename exposed a gap: the deny-list had no tripwire tying core.AgentRoles to the spec's prose — a role renamed in the code would leave the spec talking about a role the tool no longer
     refuses. Now the spec must name every role the code refuses and every retired name, and the tripwire is proved able to fail.

     Invariants

     28 external modules in internal/ui, 0 in core/cli, zero UI imports in the domain, DomainVersion still 4, kitchen/activity.json untouched at 33 events (migration at 33, stamped 4). Goldens changed one line each — the tab bar.

     Left for the reviewer

     - kitchen/project.json still says "name": "new-project" — not changed, per the brief.
     - GLOSSARY.md still carries the old names and calls the tab Expo — slice E's job, but it now contradicts the code.
     - Two ADRs and docs/audit.md still say the taster — dated records, left alone.

     Reported to the orchestrator: work/inbox/runner-slice-b.md (committed) and a nudge to pane w1Z:p1.

     Build · MiniMax-M3.1-Flash · 12m 47s · 115.7 tok/s

     ◈ Instructions updated: core/instructions

  ┃
  ┃  You are the **runner** — the same lane that wrote slice B. The orchestrator's pane id is
  ┃  `w1Z:p1`; use it verbatim in the notify command below.
  ┃
  ┃  ## The fix cycle for slice B — one finding
  ┃
  ┃  The reviewer filed **F-72 (medium)** and the orchestrator ruled on it: **accepted**.
  ┃  One finding, one fix. Do not touch anything else.
  ┃
  ┃  ### F-72 — the spec's rule sentence contradicts its own table
  ┃
  ┃  `docs/03-domain-model.md` § *The CPD ↔ tool mapping* says:
  ┃
  ┃  > The rule behind the split: the **artifact and the roles** get themed names;
  ┃  > everything an agent reasons about mechanically stays descriptive.
  ┃
  ┃  but the table right below it keeps `runner · reviewer · orchestrator` in the
  ┃  **Kept plain** column. The **table is the right half** — the decision of record
  ┃  ([docs/07-vocabulary.md](../../docs/07-vocabulary.md)) puts the agents' jobs in the
  ┃  plain vocabulary, and the code says the same. The **sentence is what is stale**, and it
  ┃  states the opposite of what the section implements, seven lines apart.
  ┃
  ┃  **The fix:** rewrite the rule sentence so it agrees with the table and with
  ┃  `docs/07-vocabulary.md`. Say it in your own words; it must be true of the table as it
  ┃  now stands — the *artifact and the place* are themed, the *agents' jobs* are plain.
  ┃
  ┃  Two things the fix is **not**:
  ┃
  ┃  - **No new tripwire.** A wire that parses prose to ask which *column* a role sits in is
  ┃    brittle, and the deny-list wire already fires when a role is renamed in the code and
  ┃    the spec is not renamed with it. Do not add one.
  ┃  - **No ADR edit.** `docs/decisions/0005-vocabulary.md` still carries the old wording,
  ┃    but an ADR is a dated record of a choice, and slice E writes ADR 0009 to supersede
  ┃    it. Say in your report that 0005 is now stale and points at 0009.
  ┃
  ┃  **Read `docs/07-vocabulary.md` § *The rule* before you write the sentence.** If the
  ┃  decision of record and `docs/03`'s table disagree, the decision of record wins and
  ┃  say so in your report.
  ┃
  ┃  ### Not yours — do not touch
  ┃
  ┃  | | |
  ┃  |---|---|
  ┃  | `docs/audit.md` | the reviewer's to write. It files F-72 and closes it. |
  ┃  | `kitchen/project.json` | F-73 — the plan's own name — is **assigned to slice C**, not here. Leave it. |
  ┃  | `docs/04-operations.md`, `GLOSSARY.md`, the old ADRs | slice E's sweep. |
  ┃  | `internal/core`, `internal/cli`, `internal/ui` | **no code changes in this cycle.** The finding is a sentence in a document. |
  ┃
  ┃  ### The gate
  ┃
  ┃  ```bash
  ┃  make lint && gofmt -l . && make test && git status --short
  ┃  ```
  ┃
  ┃  Green, and the tree clean apart from the orchestrator's own `work/` files.
  ┃
  ┃  ### One commit
  ┃
  ┃  One commit for the fix, with the evidence stamped in the body:
  ┃
  ┃  ```
  ┃  Per-AC evidence:
  ┃  - F-72: docs/03's rule sentence now matches its table; go test ./internal/core -run 'Doc' -v exit 0
  ┃  ```
  ┃
  ┃  ### Notify the orchestrator
  ┃
  ┃  ```bash
  ┃  mkdir -p work/inbox
  ┃  cat > work/inbox/runner-slice-b-fix-1.md <<EOF
  ┃  Runner: slice B fix 1 done. Commits: <n>. F-72 fixed in docs/03. Notes: <anything the reviewer should know>.
  ┃  EOF
  ┃  git add work/inbox && git commit -q -m "inbox: runner reports slice-b fix 1"
  ┃  herdr agent prompt "w1Z:p1" "Runner: slice B fix 1 done. Commits: <n>. F-72 fixed in docs/03. Notes: <anything the reviewer should know>."
  ┃  ```
  ┃

     + Thought · 242ms

     → Explored: 2 reads

     + Thought · 14.6s

  ┃
  ┃  $ grep -rn "0005\|decisions/" internal/core/docs_test.go internal/cli/*_test.go | head -10; echo "---"; grep -rn "decisions/" docs/03-domain-model.md | head
  ┃
  ┃  ---
  ┃  docs/03-domain-model.md:3:> 🔄 Draft. The one artifact designed explicitly up front ([ADR-0002](decisions/0002-human-interface-first.md));
  ┃  docs/03-domain-model.md:29:CPD's five rungs, named for a board ([ADR-0004](decisions/0004-naming.md)).
  ┃  docs/03-domain-model.md:74:[ADR-0005](decisions/0005-vocabulary.md).
  ┃  docs/03-domain-model.md:527:[ADR-0003](decisions/0003-surface-split.md) for the split.
  ┃

     + Thought · 474ms

  ┃
  ┃  ← Edit docs/03-domain-model.md
  ┃
  ┃    68   **[GLOSSARY.md](../GLOSSARY.md) owns it.** Every term there names its CPD                                                            68   **[GLOSSARY.md](../GLOSSARY.md) owns it.** Every term there names its CPD
  ┃    69   counterpart, so the two vocabularies never drift. This doc links to it rather than                                                   69   counterpart, so the two vocabularies never drift. This doc links to it rather than
  ┃    70   repeating it.                                                                                                                        70   repeating it.
  ┃    71                                                                                                                                        71
  ┃    72 - The rule behind the split: the **artifact and the roles** get themed names; everything                                               72 + The rule behind the split: the kitchen names the **artifact** and the **place**; the
  ┃    73 - an agent reasons about mechanically stays descriptive. See                                                                           73 + **agents' jobs** and everything an agent reasons about mechanically stay plain. See
  ┃    74 - [ADR-0005](decisions/0005-vocabulary.md).                                                                                            74 + [docs/07-vocabulary.md § The rule](../07-vocabulary.md#the-rule).
  ┃    75                                                                                                                                        75
  ┃    76   | Themed | Kept plain |                                                                                                              76   | Themed | Kept plain |
  ┃    77   |---|---|                                                                                                                            77   |---|---|
  ┃    78   | recipe · special · prep note · course · menu | `acceptance_criteria` · `scope` · `evidence` · `findings` · `work_packages` |       78   | recipe · special · prep note · course · menu | `acceptance_criteria` · `scope` · `evidence` · `findings` · `work_packages` |
  ┃

     → Explored: 1 read

     + Thought · 8.4s

  ┃
  ┃  $ go test ./internal/core -run 'Doc' -v 2>&1 | grep -E "^(--- |ok |FAIL)"; echo "--- gate ---"; make lint && gofmt -l . ; make test 2>&1 | tail -6; git status --short
  ┃
  ┃  (13 earlier lines)
  ┃  ?     github.com/lthiagol/maitre/cmd/maitre  [no test files]
  ┃  ?     github.com/lthiagol/maitre/cmd/maitre-cli  [no test files]
  ┃  ok    github.com/lthiagol/maitre/internal/cli  (cached)
  ┃  ok    github.com/lthiagol/maitre/internal/core  19.083s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  61.309s
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃   M docs/03-domain-model.md
  ┃

     + Thought · 1.3s

  ┃
  ┃  $ git add docs/03-domain-model.md && git commit -q -F - <<'EOF'
  ┃  F-72: the spec's rule sentence agrees with its table…
  ┃
  ┃  8c394cd F-72: the spec's rule sentence agrees with its table
  ┃  f4751d1 work: slice B's state, the lanes' panes, and where each finding lands
  ┃  18cc2e8 work: cycle 2's review brief
  ┃

     + Thought · 2.9s

  ┃
  ┃  # Wrote work/inbox/runner-slice-b-fix-1.md
  ┃
  ┃    1 # Runner — slice B, fix cycle 1
  ┃    2
  ┃    3 **Commits:** 1.
  ┃    4
  ┃    5 | Commit | What |
  ┃    6 |---|---|
  ┃    7 | `8c394cd` | **F-72** — the spec's rule sentence agrees with its table |
  ┃    8
  ┃    9 ## What changed
  ┃   10
  ┃   11 `docs/03-domain-model.md` § *The CPD ↔ tool mapping*. One sentence, and the
  ┃   12 citation that carried the stale wording.
  ┃   13
  ┃   14 **Before**
  ┃   15
  ┃   16 > The rule behind the split: the **artifact and the roles** get themed names;
  ┃   17 > everything an agent reasons about mechanically stays descriptive. See
  ┃   18 > [ADR-0005](decisions/0005-vocabulary.md).
  ┃   19
  ┃   20 **After**
  ┃   21
  ┃   22 > The rule behind the split: the kitchen names the **artifact** and the
  ┃   23 > **place**; the **agents' jobs** and everything an agent reasons about
  ┃   24 > mechanically stay plain. See
  ┃   25 > [docs/07-vocabulary.md § The rule](../07-vocabulary.md#the-rule).
  ┃   26
  ┃   27 The sentence is now true of the table as it stands: the Themed column holds the
  ┃   28 artifact (recipe · special · prep note · course · menu) and the tool's own
  ┃   29 names; the Kept-plain column holds the machine's records, the **agents' jobs**
  ┃   30 (runner · reviewer · orchestrator) and the machine's facts.
  ┃   31
  ┃   32 ## Notes for the reviewer
  ┃   33
  ┃   34 1. **The citation moved with the sentence.** The rule now points at the
  ┃   35    **decision of record** — `docs/07-vocabulary.md` § *The rule* — and not at
  ┃   36    ADR-0005. That is deliberate: ADR-0005 is the source of the stale wording
  ┃   37    (it keeps the roles in the themed column), and pointing a live sentence at a
  ┃   38    dated record that contradicts it is how the contradiction came back. The ADR
  ┃   39    itself was **not** touched — it is a dated record, and slice E writes ADR
  ┃   40    0009 to supersede it.
  ┃   41
  ┃   42 2. **ADR-0005 is now stale.** It still says the roles are themed vocabulary and
  ┃   43    still names `line cook`, `taster` and `expo`. Slice E's ADR-0009 supersedes
  ┃   44    it; this cycle left it alone, as instructed.
  ┃   45
  ┃   46 3. **One residual difference between the decision of record and `docs/03`'s
  ┃   47    table — reported, not fixed.** `docs/07`'s layer table puts **chef** and
  ┃   48    **sous-chef** in *The jobs | plain*, while `docs/03`'s table keeps
  ┃   49    `chef · maître d' · sous chef` in the **Themed** column. The two docs still
  ┃   50    differ on the cast's placement. The ruling was that the table stands, so the
  ┃   51    table was not touched, and the new sentence says nothing about the cast
  ┃   52    either way — it is silent rather than wrong. If the chef wants the cast
  ┃   53    settled, it belongs with slice E's ADR-0009 and the glossary sweep.
  ┃   54
  ┃   55 4. **No new tripwire**, as instructed: the deny-list wire already fires when a
  ┃   56    role is renamed in the code and the spec is not renamed with it.
  ┃   57
  ┃   58 5. **No code changes.** `internal/core`, `internal/cli` and `internal/ui` are
  ┃   59    untouched.
  ┃   60
  ┃   61 ## The gate
  ┃   62
  ┃   63 ```
  ┃   64 make lint && gofmt -l . && make test && git status --short
  ┃   65 ```
  ┃   66
  ┃   67 Green: `go vet` clean, `gofmt -l .` empty, core · cli · ui all pass, and the
  ┃   68 tree clean apart from the orchestrator's own `work/` files.
  ┃   69
  ┃

  ┃
  ┃  $ git add work/inbox && git commit -q -m "inbox: runner reports slice-b fix 1" && git log --oneline -2 && git status --short
  ┃
  ┃  1c386c5 inbox: runner reports slice-b fix 1
  ┃  8c394cd F-72: the spec's rule sentence agrees with its table
  ┃

  ┃
  ┃  $ herdr agent prompt "w1Z:p1" "Runner: slice B fix 1 done. Commits: 1 (8c394cd) + this inbox receipt. F-72 fixed in docs/03. Notes: the rule sentence now says 'the kitchen names the artifact and the place; the agents' jobs and everything an agent reasons about mechanically sta
  ┃  y plain', and its citation moved from the stale ADR-0005 to docs/07-vocabulary.md § The rule - the decision of record. ADR-0005 itself untouched and now stale; slice E's ADR-0009 supersedes it. One residual difference reported, not fixed: docs/07's layer table puts chef…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e8138fe5effeUJvY7OF96j92fE"},"agent_status":"done","completion_seq":352,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":
  ┃  false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":16,"screen_detection_skipped":true,"state_change_seq":352,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722","terminal_title":"OC | Maitre vocabulary rename handover",
  ┃  "terminal_title_stripped":"OC | Maitre vocabulary rename handover","workspace_id":"w1Z"},"type":"agent_prompted"}}
  ┃

     F-72 is fixed.

     The change

     docs/03-domain-model.md § The CPD ↔ tool mapping — one sentence and its citation:

     │ The rule behind the split: the kitchen names the artifact and the place; the agents' jobs and everything an agent reasons about mechanically stay plain. See docs/07-vocabulary.md § The rule.

     The old sentence said "the artifact and the roles get themed names" while the table below it keeps runner · reviewer · orchestrator in the Kept plain column. The new one is true of the table as it stands.

     The citation moved with it. The rule now points at the decision of record (docs/07 § The rule) and not at ADR-0005 — that ADR is the source of the stale wording, so pointing a live sentence at it is how the contradiction came back. The ADR itself was not touched; slice E
     writes ADR-0009 to supersede it.

     Two things for the reviewer

     - ADR-0005 is now stale — it still says the roles are themed and still names line cook, taster, expo.
     - One residual difference between the decision of record and docs/03's table, reported and not fixed: docs/07's layer table puts chef and sous-chef under The jobs | plain, while docs/03's table keeps chef · maître d' · sous chef in the Themed column. The ruling was that the
       table stands, so I did not touch it, and the new sentence is silent on the cast rather than wrong about it. Settling the cast belongs with slice E's ADR-0009 and the glossary sweep.

     No new tripwire, no code changes.

     Gate: go vet clean · gofmt -l . empty · core · cli · ui all pass · tree clean apart from the orchestrator's work/ files.

     Reported: work/inbox/runner-slice-b-fix-1.md (committed, 1c386c5) and a nudge to pane w1Z:p1.

     Build · MiniMax-M3.1-Flash · 1m 54s · 210.9 tok/s

  ┃
  ┃  You are the **runner** — the same lane that wrote slice B and its first fix. The
  ┃  orchestrator's pane id is `w1Z:p1`; use it verbatim in the notify command below.
  ┃
  ┃  ## Fix cycle 2 — F-72 is not closed yet
  ┃
  ┃  You fixed the rule sentence. **Thank you for reporting the residual difference instead
  ┃  of picking a side** — that was the right call, and it is why this round is short.
  ┃
  ┃  The orchestrator ruled: **F-72 stays open.** The section still contradicts itself — the
  ┃  sentence now says *the agents' jobs stay plain*, and the table right below it keeps two
  ┃  jobs in the **Themed** column. The fix and the residual are the same defect, one layer
  ┃  down.
  ┃
  ┃  ### The three sources, and which one wins
  ┃
  ┃  | Source | What it says about the cast |
  ┃  |---|---|
  ┃  | `docs/07-vocabulary.md` § *The rule* — **the decision of record** | The jobs \| plain \| `chef`, `sous-chef`, `runner`, `reviewer`, `orchestrator` — **all plain** |
  ┃  | `docs/03-domain-model.md`'s table — the spec | Themed \| `chef` · `maître d'` · `sous chef` \| Kept plain \| `runner` · `reviewer` · `orchestrator` |
  ┃  | `GLOSSARY.md` — the old mapping (slice E rewrites it) | `coordinator → sous chef` · `the human validator → chef` · `the human-facing console → maître d'` |
  ┃
  ┃  **The decision of record wins.** `docs/03`'s table is wrong twice over:
  ┃
  ┃  1. **It themes two jobs that are plain.** `chef` and `sous chef` are the agents' jobs,
  ┃     and `docs/07` puts every job in the plain column. Move them.
  ┃  2. **It keeps `maître d'`, which was never a job.** Per the old glossary it was *the
  ┃     human-facing console*'s themed name. The console is `maitre` now — and it is already
  ┃     in the themed column's third row. So `maître d'` does not move: **it goes**.
  ┃
  ┃  ### The fix — one table, `docs/03-domain-model.md`
  ┃
  ┃  | Column | Becomes |
  ┃  |---|---|
  ┃  | **Themed** | `recipe` · `special` · `prep note` · `course` · `menu` — **unchanged** |
  ┃  | **Kept plain** | `chef` · `sous-chef` · `runner` · `reviewer` · `orchestrator` · lifecycle states · gate names · `status` · `review_state` |
  ┃
  ┃  Write the two jobs as the code spells them — `chef`, `sous-chef` — not `sous chef`.
  ┃  The code's role is `sous-chef`; the prose in the spec may say *sous chef* when it is
  ┃  spoken, but the table is a list of names.
  ┃
  ┃  `special` and `prep note` stay: the backlog is slice C's business, and the table
  ┃  describes the format as it stands today.
  ┃
  ┃  ### The wire must still hold
  ┃
  ┃  `docs/03`'s prose must keep naming every role the code refuses — `agent`,
  ┃  `sous-chef`, `runner`, `reviewer`, `orchestrator` — and every retired name. The
  ┃  prose at § the roles already does; the tripwire reads the **whole document**, not the
  ┃  table, so moving the table's cells must not remove the last mention of anything.
  ┃
  ┃  ```bash
  ┃  go test ./internal/core -run 'Doc|DenyList' -v    # the wires still hold
  ┃  ```
  ┃
  ┃  If a name would vanish from the document entirely, **that is a signal you are about to
  ┃  break a wire** — restore the name in the prose, and say so in your report.
  ┃
  ┃  ### Not yours — do not touch
  ┃
  ┃  | | |
  ┃  |---|---|
  ┃  | `docs/07-vocabulary.md` | the decision of record. It wins; it does not change. |
  ┃  | `GLOSSARY.md` | slice E's rewrite. It still carries the old mapping — leave it, and note it. |
  ┃  | `docs/audit.md` | the reviewer's to write. |
  ┃  | `docs/decisions/0005-vocabulary.md` | a dated record; slice E's ADR 0009 supersedes it. |
  ┃  | `internal/core`, `internal/cli`, `internal/ui` | **no code changes in this cycle.** |
  ┃
  ┃  ### The gate
  ┃
  ┃  ```bash
  ┃  make lint && gofmt -l . && make test && git status --short
  ┃  ```
  ┃
  ┃  Green, and the tree clean apart from the orchestrator's own `work/` files.
  ┃
  ┃  ### One commit
  ┃
  ┃  One commit for the fix, with the evidence stamped in the body:
  ┃
  ┃  ```
  ┃  Per-AC evidence:
  ┃  - F-72: docs/03's table and its rule sentence now agree with docs/07; go test ./internal/core -run 'Doc|DenyList' -v exit 0
  ┃  ```
  ┃
  ┃  ### Notify the orchestrator
  ┃
  ┃  ```bash
  ┃  mkdir -p work/inbox
  ┃  cat > work/inbox/runner-slice-b-fix-2.md <<EOF
  ┃  Runner: slice B fix 2 done. Commits: <n>. F-72 closed in docs/03's table. Notes: <anything the reviewer should know>.
  ┃  EOF
  ┃  git add work/inbox && git commit -q -m "inbox: runner reports slice-b fix 2"
  ┃  herdr agent prompt "w1Z:p1" "Runner: slice B fix 2 done. Commits: <n>. F-72 closed in docs/03's table. Notes: <anything the reviewer should know>."
  ┃  ```
  ┃

     + Thought · 10.1s

  ┃
  ┃  $ grep -n "maître d'\|maitre d'\|themed\|Themed" docs/03-domain-model.md; echo "--- roles prose ---"; sed -n '300,325p' docs/03-domain-model.md
  ┃
  ┃  (22 earlier lines)
  ┃     so a plan written before the rename is refused exactly like one written after it.
  ┃
  ┃  ### 3.3 Overlays
  ┃
  ┃  | Overlay | Fields | Notes |
  ┃  |---|---|---|
  ┃  | `blocked` | `reason`, `by`, `at` | A recipe can be blocked *and partly executed* — v1's G8 made that inexpressible, so fragment state carried the truth while recipe state lied |
  ┃

     + Thought · 5.3s

  ┃
  ┃  ← Edit docs/03-domain-model.md
  ┃
  ┃    75                                                                                                                                        75
  ┃    76   | Themed | Kept plain |                                                                                                              76   | Themed | Kept plain |
  ┃    77   |---|---|                                                                                                                            77   |---|---|
  ┃    78   | recipe · special · prep note · course · menu | `acceptance_criteria` · `scope` · `evidence` · `findings` · `work_packages` |       78   | recipe · special · prep note · course · menu | `acceptance_criteria` · `scope` · `evidence` · `findings` · `work_packages` |
  ┃    79 - | chef · maître d' · sous chef | runner · reviewer · orchestrator · lifecycle states · gate names · `status` · `review_state` |      79 + | | `chef` · `sous-chef` · runner · reviewer · orchestrator · lifecycle states · gate names · `status` · `review_state` |
  ┃    80   | `maitre-cli` · `maitre` · the skills | |                                                                                           80   | `maitre-cli` · `maitre` · the skills | |
  ┃    81                                                                                                                                        81
  ┃    82   ### IDs                                                                                                                              82   ### IDs
  ┃    83                                                                                                                                        83
  ┃

     → Explored: 1 read

  ┃
  ┃  $ go test ./internal/core -run 'Doc|DenyList' -v 2>&1 | grep -E "^(--- |ok |FAIL)"
  ┃
  ┃  (5 earlier lines)
  ┃  --- PASS: TestTheDocDoesNotOverstateTheDenyList (0.00s)
  ┃  --- PASS: TestTheDocNamesTheRolesTheCodeRefuses (0.00s)
  ┃  --- PASS: TestTheDenyListTripwireFailsOnRealDrift (0.00s)
  ┃  --- PASS: TestTheDocDescribesWhatTheCodeDoes (0.00s)
  ┃  --- PASS: TestTheDocStatesThePlanFormatVersion (0.00s)
  ┃  --- PASS: TestTheDocNamesThePlanFormatTheCodeHas (0.00s)
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.322s
  ┃

     + Thought · 730ms

  ┃
  ┃  $ for n in agent sous-chef runner reviewer orchestrator line-cook taster expo coordinator; do printf "%-14s %s\n" "$n" "$(grep -c -- "$n" docs/03-domain-model.md)"; done; echo "--- maître d' ---"; grep -c "maître d'" docs/03-domain-model.md || true
  ┃
  ┃  (4 earlier lines)
  ┃  orchestrator   2
  ┃  line-cook      1
  ┃  taster         1
  ┃  expo           1
  ┃  coordinator    1
  ┃  --- maître d' ---
  ┃  0
  ┃

     + Thought · 1.2s

  ┃
  ┃  $ grep -rn "maître d'\|maitre d'" docs/ GLOSSARY.md README.md 2>/dev/null
  ┃
  ┃  docs/decisions/0004-naming.md:59:| `maitre` | Human console | the maître d' serves the Chef |
  ┃  docs/decisions/0004-naming.md:63:`pass`/`board`: the maître d' is the one who greets the chef and takes the order, which
  ┃  docs/decisions/0004-naming.md:76:Chef · maître d' · sous chef · line cook · taster · expo — defined in
  ┃  docs/decisions/0005-vocabulary.md:46:| Roles: chef · maître d' · sous chef · line cook · taster · expo | lifecycle states · lanes · gate names · `status` · `review_state` |
  ┃  GLOSSARY.md:124:| the human-facing console | **maître d'** |
  ┃

     + Thought · 472ms

  ┃
  ┃  $ make lint && gofmt -l . ; make test 2>&1 | tail -6; git status --short
  ┃
  ┃  go vet ./...
  ┃  ?     github.com/lthiagol/maitre/cmd/maitre  [no test files]
  ┃  ?     github.com/lthiagol/maitre/cmd/maitre-cli  [no test files]
  ┃  ok    github.com/lthiagol/maitre/internal/cli  (cached)
  ┃  ok    github.com/lthiagol/maitre/internal/core  14.308s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  (cached)
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃   M docs/03-domain-model.md
  ┃

  ┃
  ┃  $ git add docs/03-domain-model.md && git commit -q -F - <<'EOF'
  ┃  F-72: the table's roles move to the plain column, and `maître d'` goes…
  ┃
  ┃  edc08b6 F-72: the table's roles move to the plain column, and `maître d'` goes
  ┃  3bc90f9 work: fix cycle 2 — the table is the stale half, not the sentence
  ┃

  ┃
  ┃  # Wrote work/inbox/runner-slice-b-fix-2.md
  ┃
  ┃    1 # Runner — slice B, fix cycle 2
  ┃    2
  ┃    3 **Commits:** 1.
  ┃    4
  ┃    5 | Commit | What |
  ┃    6 |---|---|
  ┃    7 | `edc08b6` | **F-72** — the table's roles move to the plain column, and `maître d'` goes |
  ┃    8
  ┃    9 ## What changed
  ┃   10
  ┃   11 `docs/03-domain-model.md` § *The CPD ↔ tool mapping* — the table only.
  ┃   12
  ┃   13 **Before**
  ┃   14
  ┃   15 | Themed | Kept plain |
  ┃   16 |---|---|
  ┃   17 | recipe · special · prep note · course · menu | `acceptance_criteria` · `scope` · `evidence` · `findings` · `work_packages` |
  ┃   18 | chef · maître d' · sous chef | runner · reviewer · orchestrator · lifecycle states · gate names · `status` · `review_state` |
  ┃   19 | `maitre-cli` · `maitre` · the skills | |
  ┃   20
  ┃   21 **After**
  ┃   22
  ┃   23 | Themed | Kept plain |
  ┃   24 |---|---|
  ┃   25 | recipe · special · prep note · course · menu | `acceptance_criteria` · `scope` · `evidence` · `findings` · `work_packages` |
  ┃   26 | | `chef` · `sous-chef` · runner · reviewer · orchestrator · lifecycle states · gate names · `status` · `review_state` |
  ┃   27 | `maitre-cli` · `maitre` · the skills | |
  ┃   28
  ┃   29 - `chef` and `sous-chef` **moved** to the plain column, spelled as the code
  ┃   30   spells them.
  ┃   31 - `maître d'` **went**. It was never a job — per the old glossary it was the
  ┃   32   *human-facing console*'s themed name, and the console is `maitre` now, already
  ┃   33   in the themed column's third row.
  ┃   34 - The artifact row and the records are untouched, and `special` and `prep note`
  ┃   35   stay: the backlog is slice C's business.
  ┃   36
  ┃   37 The rule sentence above the table now agrees with it, and both agree with
  ┃   38 `docs/07-vocabulary.md` § *The rule*.
  ┃   39
  ┃   40 ## Notes for the reviewer
  ┃   41
  ┃   42 1. **No name vanished from the document.** The deny-list wire reads the whole
  ┃   43    document, so moving the cells could have taken the last mention of a role
  ┃   44    with it. It did not: the prose at § *the roles* (3.2) still names every role
  ┃   45    the code refuses and every retired name. Counts in `docs/03` after the move —
  ┃   46    `agent` 16, `sous-chef` 4, `runner` 6, `reviewer` 10, `orchestrator` 2,
  ┃   47    `line-cook` 1, `taster` 1, `expo` 1, `coordinator` 1. Nothing was restored
  ┃   48    into the prose, because nothing was lost from it.
  ┃   49
  ┃   50 2. **`maître d'` is gone from `docs/03` entirely** — the table was its only
  ┃   51    mention there. It survives in `docs/decisions/0004-naming.md` and
  ┃   52    `0005-vocabulary.md` (dated records, untouched) and in **`GLOSSARY.md`**, which
  ┃   53    still maps *the human-facing console → maître d'*. That mapping is now wrong
  ┃   54    on two counts — the console is `maitre`, and the cast is plain — and it is
  ┃   55    slice E's rewrite. Noted, not touched.
  ┃   56
  ┃   57 3. **The table's second row now has an empty Themed cell.** That is the honest
  ┃   58    shape: every job is plain, so nothing in that row is themed, and the
  ┃   59    row's plain cell still carries the machine's facts (lifecycle states, gate
  ┃   60    names, `status`, `review_state`). Nothing was restructured beyond what the
  ┃   61    ruling asked for.
  ┃   62
  ┃   63 4. **No code changes.** `internal/core`, `internal/cli` and `internal/ui` are
  ┃   64    untouched.
  ┃   65
  ┃   66 ## The gate
  ┃   67
  ┃   68 ```
  ┃   69 make lint && gofmt -l . && make test && git status --short
  ┃   70 ```
  ┃   71
  ┃   72 Green: `go vet` clean, `gofmt -l .` empty, core · cli · ui all pass, and the
  ┃   73 tree clean apart from the orchestrator's own `work/` files.
  ┃   74
  ┃

  ┃
  ┃  $ git add work/inbox && git commit -q -m "inbox: runner reports slice-b fix 2" && git log --oneline -2 && git status --short
  ┃
  ┃  723d991 inbox: runner reports slice-b fix 2
  ┃  edc08b6 F-72: the table's roles move to the plain column, and `maître d'` goes
  ┃

  ┃
  ┃  $ herdr agent prompt "w1Z:p1" "Runner: slice B fix 2 done. Commits: 1 (edc08b6) + this inbox receipt. F-72 closed in docs/03's table. Notes: chef and sous-chef moved from the Themed column into the Kept-plain one, spelled as the code spells them (chef, sous-chef); 'maître d' i
  ┃  s gone - it was never a job, only the old glossary's themed name for the human-facing console, and the console is maitre, already in the themed column's third row. The table's second row now has an empty Themed cell, which is the honest shape: every job is plain, and th…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e8138fe5effeUJvY7OF96j92fE"},"agent_status":"done","completion_seq":356,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":
  ┃  false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":16,"screen_detection_skipped":true,"state_change_seq":356,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722","terminal_title":"OC | Maitre vocabulary rename handover",
  ┃  "terminal_title_stripped":"OC | Maitre vocabulary rename handover","workspace_id":"w1Z"},"type":"agent_prompted"}}
  ┃

     F-72 is closed at both layers.

     The change

     docs/03-domain-model.md § The CPD ↔ tool mapping — the table only:

     ┌──────────────────────────────────────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Themed                                                                                                   │ Kept plain                                                                                                                                                             │
     ├──────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ recipe · special · prep note · course · menu                                                             │ acceptance_criteria · scope · evidence · findings · work_packages                                                                                                      │
     ├──────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │                                                                                                          │ chef · sous-chef · runner · reviewer · orchestrator · lifecycle states · gate names · status · review_state                                                            │
     ├──────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ maitre-cli · maitre · the skills                                                                         │                                                                                                                                                                        │
     └──────────────────────────────────────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

     - chef and sous-chef moved to the plain column, spelled as the code spells them.
     - maître d' went — never a job, only the old glossary's themed name for the human-facing console, and the console is maitre now (already in the themed column's third row).
     - Artifact row, records, and special/prep note untouched.

     The sentence above the table and the table now agree with each other and with docs/07 § The rule.

     Three notes for the reviewer

     1. No name vanished from the document. The deny-list wire reads the whole doc, so moving cells could have taken a role's last mention with it. It did not — the prose at § 3.2 still names every role the code refuses and every retired name (agent 16, sous-chef 4, runner 6,
        reviewer 10, orchestrator 2, line-cook 1, taster 1, expo 1, coordinator 1). Nothing had to be restored.

     2. maître d' is gone from docs/03 entirely — the table was its only mention there. It survives in ADRs 0004/0005 (dated records) and in GLOSSARY.md, which still maps the human-facing console → maître d' — now wrong on two counts, and slice E's rewrite. Noted, not touched.

     3. The table's second row now has an empty Themed cell. That is the honest shape: every job is plain, so nothing in that row is themed, and its plain cell still carries the machine's facts.

     Gate: go vet clean · gofmt -l . empty · core · cli · ui all pass · tree clean apart from the orchestrator's work/ files.

     Reported: work/inbox/runner-slice-b-fix-2.md (committed, 723d991) and a nudge to pane w1Z:p1.

     Build · MiniMax-M3.1-Flash · 1m 3s · 185.6 tok/s

  ┃
  ┃
  ┃
  ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
  ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
  ~/Code/Pessoal/maitre:master                                                                                                                                                                                                                               147.4K (15%)  ctrl+p commands

```
