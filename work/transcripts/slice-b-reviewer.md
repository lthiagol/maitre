# The reviewer's lane — slice B, captured 2026-10-27

Dumped from pane `w1Z:p4` with `herdr pane read --source recent-unwrapped --lines 5000`.
A fresh lane: it re-derived the facts from the code, filed F-72 (medium) and F-73
(low), and verified the two fix cycles by mutation — `agentRole` returning `""`,
the `line-cook` mapping dropped, `sameAgent` ignoring the role, the spec losing the
word `runner`.

```
  ┃  ### IDs
  ┃
  ┃  | Entity | Form | Scope |
  ┃  |---|---|---|
  ┃  | recipe | `RCP01`, children `RCP01.1` | workspace |
  ┃

     + Thought · 5.3s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log -L 71,80:docs/03-domain-model.md --oneline 2>/dev/null | head -60
  ┃
  ┃  (44 earlier lines)
  ┃  +[ADR-0005](decisions/0005-vocabulary.md).
  ┃  +
  ┃  +| Themed | Kept plain |
  ┃  +|---|---|
  ┃  +| recipe · special · prep note · course · menu | `acceptance_criteria` · `scope` · `evidence` · `findings` · `work_packages` |
  ┃  +| chef · maître d' · sous chef · line cook · taster · expo | lifecycle states · gate names · `status` · `review_state` |
  ┃  +| `scale` · `maitre` · the skills | |
  ┃

     + Thought · 7.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "themed\|Kept plain\|kept plain" --include="*.md" docs/ AGENTS.md GLOSSARY.md | grep -v "docs/audit.md"
  ┃
  ┃  docs/05-human-interface.md:417:**`glamour`, themed to our palette.** A recipe's intent, problem and audit trail are
  ┃  docs/05-human-interface.md:492:| Prose | `glamour`, themed to our palette |
  ┃  docs/03-domain-model.md:72:The rule behind the split: the **artifact and the roles** get themed names; everything
  ┃  docs/03-domain-model.md:76:| Themed | Kept plain |
  ┃  docs/06-proposals.md:33:| The glossary already names the rung: *"One intent across several recipes — a themed set of dishes served together."* | [GLOSSARY.md](../GLOSSARY.md) |
  ┃  docs/decisions/0005-vocabulary.md:43:| Themed | Kept plain |
  ┃  GLOSSARY.md:33:One intent across several recipes — a themed set of dishes served together.
  ┃

     + Thought · 9.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git show af6a723 -- docs/03-domain-model.md docs/05-human-interface.md | grep -E "^[+-]" | grep -vE "^(\+\+\+|---)"
  ┃
  ┃  (19 earlier lines)
  ┃  +   into a machine check instead of a prose reminder (P2). The names the rename retired
  ┃  +   — `line-cook`, `taster`, `expo`, `coordinator` — still name the role they stood for,
  ┃  +   so a plan written before the rename is refused exactly like one written after it.
  ┃  -      recipe: "RCP07", actor: "human", session: "s-3", role: "coordinator",
  ┃  +      recipe: "RCP07", actor: "human", session: "s-3", role: "sous-chef",
  ┃  -console claiming work as `line-cook`.
  ┃  +console claiming work as `runner`.
  ┃

     + Thought · 2.6s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cat internal/ui/testdata/console-recipes.golden
  ┃
  ┃  (33 earlier lines)
  ┃
  ┃
  ┃
  ┃
  ┃
  ┃
  ┃   ↑↓ move · tab switch · enter detail · space sign off · 1 2 3 order · c served · r reload · q quit
  ┃

     + Thought · 973ms

  ┃
  ┃  $ cd /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/kb && /Users/thiago/Code/Pessoal/maitre/bin/maitre --keys "enter" --plain 2>&1 | head -24
  ┃
  ┃  (17 earlier lines)
  ┃    in   the sign-off action
  ┃    out  web shell
  ┃    out  orchestration
  ┃
  ┃  Acceptance criteria
  ┃    AC-01   passed   the board renders 100 recipes
  ┃            verify: go test ./internal/core -run Board
  ┃

     + Thought · 4.5s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "func ReadinessOf" -A40 internal/core/views.go | head -60
  ┃
  ┃  (34 earlier lines)
  ┃  521-// ProofSummary counts the evidence attached to a recipe.
  ┃  522-//
  ┃  523-// Rejected and unproven are separate: a rejected proof matched nothing (a spec defect),
  ┃  524-// an unproven one could not be attributed (not a machine proof). Different remedies.
  ┃  525-type ProofSummary struct {
  ┃  526-  Criteria int
  ┃  527-  Proven   int
  ┃

     + Thought · 10.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '300,345p' internal/cli/cli.go
  ┃
  ┃  (39 earlier lines)
  ┃      return v
  ┃    }
  ┃    v := os.Getenv(deprecated)
  ┃    if v == "" {
  ┃      return ""
  ┃    }
  ┃    r.warn(fmt.Sprintf("$%s is deprecated and will be removed; use $%s instead", deprecated, current))
  ┃

     + Thought · 4.8s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "The rule behind the split" docs/03-domain-model.md && grep -n "runner · reviewer · orchestrator" docs/03-domain-model.md && grep -n "new-project" kitchen/project.json
  ┃
  ┃  72:The rule behind the split: the **artifact and the roles** get themed names; everything
  ┃  79:| chef · maître d' · sous chef | runner · reviewer · orchestrator · lifecycle states · gate names · `status` · `review_state` |
  ┃  8:  "name": "new-project",
  ┃

     + Thought · 350ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '10,22p' docs/07-vocabulary.md
  ┃
  ┃  (6 earlier lines)
  ┃  | The jobs | plain | chef, sous-chef, **runner**, **reviewer**, **orchestrator** |
  ┃  | The plan's loose end | plain | **backlog item**, **idea** / **feature** / **bug** |
  ┃
  ┃  Two consequences that used to be accidents and are now rules:
  ┃
  ┃  - The console says *tasting note*; the code says `Finding`. That is a **translation**,
  ┃    not a mismatch — the console is the surface where the kitchen is free to be charming.
  ┃

     + Thought · 4.1s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && git log --oneline -3 && echo "=== TREE CLEAN ==="
  ┃
  ┃  3436be2 work: the review brief names the orchestrator's pane
  ┃  c3d2110 work: the lanes wait, they do not poll
  ┃  a3e1b9a inbox: runner reports slice-b
  ┃  === TREE CLEAN ===
  ┃

     + Thought · 153ms

     → Explored: 2 reads



  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    75   |---|---|---|---|                                                                                                                    75   |---|---|---|---|
  ┃    76   | **F-59** | medium | **`core.Audit.Decisions` has a shape and no writer** — it is where a decision *about the plan* would live, and 76   | **F-59** | medium | **`core.Audit.Decisions` has a shape and no writer** — it is where a decision *about the plan* would live, and
  ┃         nothing writes it. `kitchen/decisions.json` is named by the plan format and opened by nothing. The remedy is on the table in [06-         nothing writes it. `kitchen/decisions.json` is named by the plan format and opened by nothing. The remedy is on the table in [06-
  ┃         proposals.md](06-proposals.md) § 2. | `internal/core/model.go:270`; `internal/core/store.go:24` |                                         proposals.md](06-proposals.md) § 2. | `internal/core/model.go:270`; `internal/core/store.go:24` |
  ┃    77   | **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.resolved` gate can never fire** — nothing ever   77   | **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.resolved` gate can never fire** — nothing ever
  ┃         sets a question, so the list is always empty and the gate always passes. A gate nobody can fail is not a guard. | `internal/core/         sets a question, so the list is always empty and the gate always passes. A gate nobody can fail is not a guard. | `internal/core/
  ┃         model.go:31`; `internal/core/lifecycle.go:197` |                                                                                          model.go:31`; `internal/core/lifecycle.go:197` |
  ┃    78   | **F-61** | medium | **The `courses/` directory is named by the plan format, created by `Init`, and never written** — a directory   78   | **F-61** | medium | **The `courses/` directory is named by the plan format, created by `Init`, and never written** — a directory
  ┃         no writer touches. Per [06-proposals.md](06-proposals.md) § 1 the honest answer is to stop naming it: a course is the relation, not       no writer touches. Per [06-proposals.md](06-proposals.md) § 1 the honest answer is to stop naming it: a course is the relation, not
  ┃         a file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.md) |                                                              a file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.md) |
  ┃                                                                                                                                              79 + | **F-72** | medium | **The spec's vocabulary rule contradicts its own table after the role rename** — `docs/03-domain-model.md` §
  ┃                                                                                                                                                   "The CPD ↔ tool mapping" says *"the **artifact and the roles** get themed names"*, but the table now keeps `runner · reviewer ·
  ┃                                                                                                                                                   orchestrator` in the **Kept plain** column. The table is the right half: the decision of record ([docs/07-vocabulary.md](07-
  ┃                                                                                                                                                   vocabulary.md) § "The rule") puts the agents' jobs in the plain vocabulary, and the code says the same. The sentence was not touched
  ┃                                                                                                                                                   when the roles moved, so the section states the opposite of what it implements, seven lines apart. Nothing catches it: the deny-list
  ┃                                                                                                                                                   tripwire only asks whether the doc *names* a role, not which column it sits in. The rule sentence and the table must say the same
  ┃                                                                                                                                                   thing. | `docs/03-domain-model.md:72`; the table at `docs/03-domain-model.md:79` |
  ┃    79   | F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and the size.** It is path + mtime + size per file, 80   | F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and the size.** It is path + mtime + size per file,
  ┃         so two different contents can share one fingerprint. The tool's own writes always change the mtime, so in practice only a hand edit       so two different contents can share one fingerprint. The tool's own writes always change the mtime, so in practice only a hand edit
  ┃         that restores the mtime is missed — and the chef's `r` reloads it. | `internal/core/store.go` — `DiskFingerprint`. |                      that restores the mtime is missed — and the chef's `r` reloads it. | `internal/core/store.go` — `DiskFingerprint`. |
  ┃                                                                                                                                              81 + | F-73 | low | **`kitchen/project.json` still carries `"name": "new-project"`** — the decision of record names the project
  ┃                                                                                                                                                   **maitre**, and the console's header renders the plan's name, so the repo's own console reads `maitre · new-project`. It is the
  ┃                                                                                                                                                   *plan's* data, not the code's: slice A renamed the module and the binaries and left this file, and slice C is the slice that
  ┃                                                                                                                                                   migrates the repo's own plan. | `kitchen/project.json:8` |
  ┃    80   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-scope=`. | Recorded so nobody calls it a regression 82   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-scope=`. | Recorded so nobody calls it a regression
  ┃         later. Not a defect. |                                                                                                                    later. Not a defect. |
  ┃    81                                                                                                                                        83
  ┃    82   ## Closed findings                                                                                                                   84   ## Closed findings
  ┃    83                                                                                                                                        85
  ┃

     + Thought · 956ms

     → Explored: 1 read

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    35   | Cycle 14 | The last sweep: the check reads the package, and the marker travels with the function | One note closed, and the        35   | Cycle 14 | The last sweep: the check reads the package, and the marker travels with the function | One note closed, and the
  ┃         families that admit a view without a classification |                                                                                     families that admit a view without a classification |
  ┃    36   | Cycle 15 | The closing review of the night: no read-side family, and a symmetric judgement | One note closed, and the shape judged 36   | Cycle 15 | The closing review of the night: no read-side family, and a symmetric judgement | One note closed, and the shape judged
  ┃         |                                                                                                                                         |
  ┃    37   | Slice A | The vocabulary rename: the module, the binaries, the recipe ids, the version | The ids and the guard hold; two mediums   37   | Slice A | The vocabulary rename: the module, the binaries, the recipe ids, the version | The ids and the guard hold; two mediums
  ┃         and four lows filed |                                                                                                                     and four lows filed |
  ┃    38   | Cycle 2 | The fixes to slice A: the usage lines, the two-way wire, the plan's stamp, the format facts | Six notes closed, and the  38   | Cycle 2 | The fixes to slice A: the usage lines, the two-way wire, the plan's stamp, the format facts | Six notes closed, and the
  ┃         fence checked by mutation |                                                                                                               fence checked by mutation |
  ┃                                                                                                                                              39 + | Slice B | The vocabulary rename: the roles and the tabs | The roles resolve as written *and* as a role, so a plan written before
  ┃                                                                                                                                                   the rename still cannot close on its own executor's review; the tabs are Board and Sous-Chef, and the honesty of the waiting tab
  ┃                                                                                                                                                   held; one medium filed in the spec |
  ┃    39   | The rounds inside a cycle | Fix → verify → back to the executor → repeat, until the verifier is satisfied | Three rounds in cycle  40   | The rounds inside a cycle | Fix → verify → back to the executor → repeat, until the verifier is satisfied | Three rounds in cycle
  ┃         8: the first fix introduced three new mediums, the second closed them |                                                                   8: the first fix introduced three new mediums, the second closed them |
  ┃    40                                                                                                                                        41
  ┃    41   ### The discipline every round runs                                                                                                  42   ### The discipline every round runs
  ┃    42                                                                                                                                        43
  ┃

     + Thought · 923ms

     → Explored: 1 read

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    68   | The lifecycles | The state table and the code's lifecycles, both ways |                                                            68   | The lifecycles | The state table and the code's lifecycles, both ways |
  ┃    69   | The interface | The footer paragraph names the triage |                                                                            69   | The interface | The footer paragraph names the triage |
  ┃    70   | The plan format | Every directory and file the code creates is named in `docs/03` § 6, and every directory and file § 6 names is   70   | The plan format | Every directory and file the code creates is named in `docs/03` § 6, and every directory and file § 6 names is
  ┃         the code's — both ways — and the section's version is the code's |                                                                        the code's — both ways — and the section's version is the code's |
  ┃    71   | The usage lines | Every usage line the CLI prints names a binary that exists, and the console's errors name it too |               71   | The usage lines | Every usage line the CLI prints names a binary that exists, and the console's errors name it too |
  ┃                                                                                                                                              72 + | The deny-list | The spec names every role the code refuses, and every name the rename retired — a plan written before the rename
  ┃                                                                                                                                                   still meets those |
  ┃    72                                                                                                                                        73
  ┃    73   ## Open findings                                                                                                                     74   ## Open findings
  ┃    74                                                                                                                                        75
  ┃    75   | # | Severity | What | Where |                                                                                                      76   | # | Severity | What | Where |
  ┃

     + Thought · 619ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff docs/audit.md
  ┃
  ┃  (23 earlier lines)
  ┃   | **F-61** | medium | **The `courses/` directory is named by the plan format, created by `Init`, and never written** — a directory no writer touches. Per [06-proposals.md](06-proposals.md) § 1 the honest answer is to stop naming it: a course is the relation, not a file. |
  ┃  `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.md) |
  ┃  +| **F-72** | medium | **The spec's vocabulary rule contradicts its own table after the role rename** — `docs/03-domain-model.md` § "The CPD ↔ tool mapping" says *"the **artifact and the roles** get themed names"*, but the table now keeps `runner · reviewer · orchestrator` in
  ┃  the **Kept plain** column. The table is the right half: the decision of record ([docs/07-vocabulary.md](07-vocabulary.md) § "The rule") puts the agents' jobs in the plain vocabulary, and the code says the same. The sentence was not touched when the roles moved, so the section
  ┃  states the opposite of what it implements, seven lines apart. Nothing catches it: the deny-list tripwire only asks whether the doc *names* a role, not which column it sits in. The rule sentence and the table must say the same thing. | `docs/03-domain-model.md:72`; the table at
  ┃  `docs/03-domain-model.md:79` |
  ┃   | F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and the size.** It is path + mtime + size per file, so two different contents can share one fingerprint. The tool's own writes always change the mtime, so in practice only a hand edit that
  ┃  restores the mtime is missed — and the chef's `r` reloads it. | `internal/core/store.go` — `DiskFingerprint`. |
  ┃  +| F-73 | low | **`kitchen/project.json` still carries `"name": "new-project"`** — the decision of record names the project **maitre**, and the console's header renders the plan's name, so the repo's own console reads `maitre · new-project`. It is the *plan's* data, not the
  ┃  code's: slice A renamed the module and the binaries and left this file, and slice C is the slice that migrates the repo's own plan. | `kitchen/project.json:8` |
  ┃   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-scope=`. | Recorded so nobody calls it a regression later. Not a defect. |
  ┃
  ┃   ## Closed findings
  ┃

     + Thought · 206ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' 2>&1 | tail -5 && go test ./... 2>&1 | tail -6
  ┃
  ┃  (2 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  0.279s [no tests to run]
  ┃  ?     github.com/lthiagol/maitre/cmd/maitre  [no test files]
  ┃  ?     github.com/lthiagol/maitre/cmd/maitre-cli  [no test files]
  ┃  ok    github.com/lthiagol/maitre/internal/cli  (cached)
  ┃  ok    github.com/lthiagol/maitre/internal/core  19.529s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  60.664s
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make regen-goldens >/dev/null 2>&1 && git status --short && echo "--- only the ledger should be modified ---"
  ┃
  ┃  M docs/audit.md
  ┃  --- only the ledger should be modified ---
  ┃

     Now the inbox receipt and the commits.

  ┃
  ┃  # Wrote work/inbox/reviewer-slice-b.md
  ┃
  ┃    1 Reviewer: slice B review done. Verdict: changes-needed. Closed: none. Filed 2 new blockers (F-72, F-73 — F-72 is the only MEDIUM; F-73 is the low `kitchen/project.json` drift).
  ┃    2
  ┃    3 ## What held
  ┃    4
  ┃    5 - **The sharp edge holds.** `core.agentRole()` resolves a name as written *and* as a role,
  ┃    6   and `core.sameAgent()` compares the executor with the reviewer through it. A plan
  ┃    7   written before the rename cannot close on its own executor's review — proved end to
  ┃    8   end on a hand-edited plan (`claimed_by: line-cook`): the executor reviewing under
  ┃    9   `line-cook` is refused, the same executor reviewing under the **new** name `runner`
  ┃   10   is refused, and an independent `reviewer` is accepted.
  ┃   11 - **The deny-list got tighter, not looser.** Every current and every retired role, in
  ┃   12   every spelling (`Line Cook`, `LINE-COOK`, ` line-cook `, `Runner`), is refused the
  ┃   13   approval; `chef` and an arbitrary human name are not. An unknown role still fails to
  ┃   14   match the list, exactly as before — no silent enum fallback.
  ┃   15 - **The tabs.** `TabBoard` / `TabSousChef` throughout; `Chopboard` and `Expo` survive
  ┃   16   nowhere outside the three record exclusions. The Sous-Chef tab still says the
  ┃   17   automation is not built yet, still shows the vital signs, and grew no crew view. The
  ┃   18   header is one line; the goldens moved by the tab bar and nothing else.
  ┃   19 - **One word per fact.** The row shows the plain state word, the lane is a section
  ┃   20   header, and the detail pane's meta line no longer carries the lane's label beside
  ┃   21   the state. The blocked reason kept its own line, so nothing was lost.
  ┃   22 - **Invariants.** `internal/core` and `internal/cli` are stdlib-only; `internal/ui` is
  ┃   23   at 28 external modules; the domain layer imports no UI; `core.DomainVersion` is 4;
  ┃   24   `kitchen/activity.json` is untouched — 33 events, the 32 old ones still `R01`, the
  ┃   25   migration the 33rd, a `control` event stamped 4.
  ┃   26 - **Scope.** The roles and the tabs only. The backlog, the ids, the module path, the
  ┃   27   binaries, `Proof`→evidence, the glossary and the mascot are untouched. The WIP diff
  ┃   28   (38 files) is entirely the rename — no other change rode along.
  ┃   29
  ┃   30 ## The blocker
  ┃   31
  ┃   32 **F-72 (medium)** — `docs/03-domain-model.md` § "The CPD ↔ tool mapping" now contradicts
  ┃   33 itself. The table keeps `runner · reviewer · orchestrator` in the **Kept plain** column
  ┃   34 (this slice moved them), but the rule two lines above still says *"the **artifact and the
  ┃   35 roles** get themed names"*. The table is the right half — the decision of record puts the
  ┃   36 agents' jobs in the plain vocabulary — so the sentence is what is stale. Nothing checks
  ┃   37 it: the new deny-list tripwire asks whether the doc *names* a role, not which column it
  ┃   38 sits in.
  ┃   39
  ┃   40 ## Mutation checks (all reverted; the tree is as the runner left it)
  ┃   41
  ┃   42 | Mutation | Result |
  ┃   43 |---|---|
  ┃   44 | `core.agentRole` returns `""` | `TestASelfReviewSurvivesTheRename`, `TestTheRenamedRolesAreTheAgentRoles` and 13 more fail |
  ┃   45 | drop the `line-cook` → `runner` mapping | `TestASelfReviewSurvivesTheRename`, `TestTheRenamedRolesAreTheAgentRoles` fail |
  ┃   46 | `core.sameAgent` ignores the role | `TestASelfReviewSurvivesTheRename` fails |
  ┃   47 | `docs/03` no longer names `runner` | the new deny-list tripwire fails on real drift |
  ┃   48
  ┃   49 ## LOW notes (not filed)
  ┃   50
  ┃   51 - `kitchen/project.json` still says `"name": "new-project"`; the console's header
  ┃   52   renders it, so the repo's own console reads `maitre · new-project`. **Filed as F-73
  ┃   53   (low)** — it is plan data, it belongs to slice C's plan migration, and it is not
  ┃   54   slice B's to change.
  ┃   55 - `docs/audit.md`'s closed rows T4 and T6 still name `internal/ui/chop_test.go` and
  ┃   56   `internal/ui/expo_test.go`, which this slice renamed. A closed row claims a test
  ┃   57   holds the fix, and those two point at files that no longer exist. Not swept here.
  ┃   58 - `docs/04-operations.md` names the roles in prose as well as `docs/03` does; only
  ┃   59   `docs/03` is wired to the code.
  ┃   60 - The CLI's documented example (`--actor runner`) has no tripwire; a future rename
  ┃   61   would silently put a retired name back in the help text.
  ┃   62 - The detail pane's readiness reason restates the state in kitchen words ("ready" and
  ┃   63   "awaiting sign-off"). The lane's *label* is gone from the pane, so the invariant as
  ┃   64   stated holds; the sentence remains.
  ┃   65 - `GLOSSARY.md` still maps `runner`→`line cook`, `reviewer`→`taster`,
  ┃   66   `orchestrator`→`expo` — slice E's rewrite, and it is the file the whole rename
  ┃   67   points at.
  ┃   68
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/audit.md && git commit -q -m "docs: the ledger files slice B's spec drift
  ┃  …
  ┃
  ┃  26902c1 docs: the ledger files slice B's spec drift
  ┃  3436be2 work: the review brief names the orchestrator's pane
  ┃  ?? work/inbox/reviewer-slice-b.md
  ┃

     + Thought · 229ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log -3 --format='%H%n%an <%ae>%n%b---'
  ┃
  ┃  (29 earlier lines)
  ┃  Also: the repo moved to /Users/thiago/Code/Pessoal/maitre, so the briefs carry
  ┃  the new path, and the slice B brief names the orchestrator's pane outright
  ┃  rather than leaving $ORCH_PANE to be resolved by whoever reads it.
  ┃
  ┃  The review brief for slice B now proves the retired tab names are gone, and it
  ┃  does not false-positive on the decision of record's own rename table.
  ┃  ---
  ┃

     + Thought · 272ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git commit -q --amend -m "docs: the ledger files slice B's spec drift
  ┃  …
  ┃
  ┃  (9 earlier lines)
  ┃
  ┃  The project's own plan still calls itself new-project; the console renders
  ┃  that name in its header. It is the plan's data and it belongs to the slice
  ┃  that migrates the plan.
  ┃
  ┃  ?? work/inbox/reviewer-slice-b.md
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/inbox && git commit -q -m "inbox: reviewer reports slice-b" && git log --oneline -3 && git status --short && echo "--- CLEAN ---"
  ┃
  ┃  6410b65 inbox: reviewer reports slice-b
  ┃  41075bb docs: the ledger files slice B's spec drift
  ┃  3436be2 work: the review brief names the orchestrator's pane
  ┃  --- CLEAN ---
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1Z:p1" "Reviewer: slice B review done. Verdict: changes-needed. Closed: none. Filed 2 new blockers (F-72 medium, F-73 low). Mutation checks: agentRole returns \"\" (roles+self-review tests fail), dropped the line-c
  ┃  ook mapping (roles+self-review tests fail), sameAgent ignores the role (self-review test fails), docs/03 no longer names runner (the new deny-list tripwire fires). All reverted; the tree is as the runner left it. LOW notes: the ledger's T4/T6 rows still name the renamed…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e8138fe5effeUJvY7OF96j92fE"},"agent_status":"idle","completion_seq":348,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":
  ┃  false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":16,"screen_detection_skipped":true,"state_change_seq":348,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722","terminal_title":"OC | Maitre vocabulary rename handover",
  ┃  "terminal_title_stripped":"OC | Maitre vocabulary rename handover","workspace_id":"w1Z"},"type":"agent_prompted"}}
  ┃  NOTIFIED
  ┃

     + Thought · 300ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff 3436be2..HEAD --stat && echo "=== tree ===" && git status --short && echo CLEAN
  ┃
  ┃  docs/audit.md                  |  4 +++
  ┃   work/inbox/reviewer-slice-b.md | 67 ++++++++++++++++++++++++++++++++++++++++++
  ┃   2 files changed, 71 insertions(+)
  ┃  === tree ===
  ┃  CLEAN
  ┃

     Review of slice B is done. Verdict: changes-needed — one MEDIUM blocker filed.

     What held

     ┌─────────────────────────┬─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Check                   │ Result                                                                                                                                                                                                                                                  │
     ├─────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ The sharp edge          │ Holds, end to end. On a plan written before the rename (claimed_by: line-cook): the executor reviewing under line-cook → refused; the same executor under the new name runner → refused; an independent reviewer → accepted.                            │
     ├─────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ No silent enum fallback │ The deny-list got tighter. Every current and retired role, in every spelling (Line Cook, LINE-COOK, line-cook, Runner), is refused the approval; chef and an arbitrary human name are not.                                                              │
     ├─────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ The tabs                │ TabBoard / TabSousChef throughout. Chopboard/Expo survive nowhere outside the three record exclusions. The Sous-Chef tab still says not built yet, still shows the vital signs, grew no crew view. Header is one line; goldens moved by the tab bar     │
     │                         │ only.                                                                                                                                                                                                                                                   │
     ├─────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ One word per fact       │ The row shows the plain state word, the lane is a section header, and the detail's meta line no longer carries the lane's label. The blocked reason kept its own line — nothing lost.                                                                   │
     ├─────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ Invariants              │ core/cli stdlib-only, ui at 28 modules, zero UI imports, DomainVersion 4, activity.json untouched (33 events, the 32 old ones still R01, the migration the 33rd, control, stamped 4).                                                                   │
     ├─────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ Scope                   │ Roles and tabs only. The 38-file WIP diff is entirely the rename — nothing else rode along.                                                                                                                                                             │
     └─────────────────────────┴─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

     The blocker

     F-72 (medium) — docs/03-domain-model.md § "The CPD ↔ tool mapping" now contradicts itself. The table keeps runner · reviewer · orchestrator in the Kept plain column (this slice moved them), but the rule two lines above still says "the artifact and the roles get themed names".
     The table is the right half — the decision of record puts the agents' jobs in the plain vocabulary — so the sentence is the stale half. Nothing catches it: the new deny-list tripwire asks whether the doc names a role, not which column it sits in.

     F-73 (low) — kitchen/project.json still says "name": "new-project"; the console renders it, so the repo's own console reads maitre · new-project. It is drift, but it is plan data and belongs to slice C's plan migration — filed, not changed.

     Mutation checks (all reverted, tree clean)

     ┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Mutation                                                                                                            │ Result                                                                                                                                                      │
     ├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ agentRole returns ""                                                                                                │ 15 tests fail, incl. TestASelfReviewSurvivesTheRename                                                                                                       │
     ├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ drop the line-cook→runner mapping                                                                                   │ TestASelfReviewSurvivesTheRename, TestTheRenamedRolesAreTheAgentRoles fail                                                                                  │
     ├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ sameAgent ignores the role                                                                                          │ TestASelfReviewSurvivesTheRename fails                                                                                                                      │
     ├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ docs/03 stops naming runner                                                                                         │ the new deny-list tripwire fires on real drift                                                                                                              │
     └─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘



     Build · MiniMax-M3.1-Flash · 13m 23s · 164.7 tok/s

  ┃
  ┃  You are the **reviewer** — the same lane that reviewed slice B. The orchestrator's pane
  ┃  id is `w1Z:p1`; use it verbatim in the notify command below.
  ┃
  ┃  ## Cycle 2 — the fix to slice B
  ┃
  ┃  The runner fixed the one finding you filed. **Cycle 2 is blockers only**: file **HIGH
  ┃  or MEDIUM** blockers. No nits, no cosmetics, no style.
  ┃
  ┃  **Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's) or **F-73** (the plan's
  ┃  own name — the orchestrator assigned it to slice C, and it stays open).
  ┃
  ┃  ### 1. Is F-72 actually closed?
  ┃
  ┃  **F-72** — `docs/03-domain-model.md` § *The CPD ↔ tool mapping* said *"the artifact
  ┃  **and the roles** get themed names"* while the table right below kept `runner ·
  ┃  reviewer · orchestrator` in the **Kept plain** column.
  ┃
  ┃  The fix rewrites the rule sentence so it agrees with the table and with the decision
  ┃  of record. Verify **the sentence and the table say the same thing** — read them
  ┃  together, and read `docs/07-vocabulary.md` § *The rule* beside them, since the
  ┃  decision of record is what both must match.
  ┃
  ┃  A sentence that no longer contradicts the table but now says something the decision of
  ┃  record does not say is **not closed**. Judge the sentence against the decision, not
  ┃  against the table alone.
  ┃
  ┃  Then confirm the runner did what the brief said and nothing more:
  ┃
  ┃  - **No code changed.** The finding was a sentence in a document. `git diff` over the
  ┃    fix's commit must touch `docs/03-domain-model.md` and nothing in `internal/`.
  ┃  - **No new tripwire was added.** A wire that parses prose to ask which column a role
  ┃    sits in is brittle; the deny-list wire already covers the rename. A new wire here is
  ┃    scope creep, not a fix.
  ┃  - **`docs/decisions/0005-vocabulary.md` was not edited.** An ADR is a dated record;
  ┃    slice E writes ADR 0009 to supersede it.
  ┃
  ┃  ```bash
  ┃  go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v   # the wires still hold
  ┃  make lint && gofmt -l . && make test
  ┃  ```
  ┃
  ┃  ### 2. Two ledger rows point at files that no longer exist
  ┃
  ┃  The finding is yours, so the pointer is yours to fix.
  ┃
  ┃  | Row | Says | Now |
  ┃  |---|---|---|
  ┃  | **T4** | `internal/ui/chop_test.go` — `core.TreeOf`; the console's Board tab | `internal/ui/board_test.go` |
  ┃  | **T6** | `internal/ui/expo_test.go` — the Sous-Chef tab and the automation's vital signs | `internal/ui/souschef_test.go` |
  ┃
  ┃  A closed row claims *a test holds the fix*. A pointer at a file that was renamed
  ┃  makes that claim false — the test still holds, but the row points at a ghost.
  ┃
  ┃  **Update the file pointers only. The finding's wording stays as it was filed** — the
  ┃  ledger is the record of what was found, and it is not swept. This is a pointer that
  ┃  follows a rename, not a rewrite.
  ┃
  ┃  ### 3. Close F-72
  ┃
  ┃  When it is closed, move its row from open to closed, with the file and the test that
  ┃  now holds it. Commit the ledger.
  ┃
  ┃  ### Verdict
  ┃
  ┃  - No HIGH or MEDIUM findings → `verdict: ok`.
  ┃  - Any HIGH or MEDIUM → `verdict: changes-needed`.
  ┃  - LOW findings: **note them in your message, do not file them.**
  ┃
  ┃  ### Hard rules
  ┃
  ┃  - Do **not** edit implementation code, and do **not** edit `docs/03` — the runner
  ┃    owns the fix. Read-only for everything except `docs/audit.md`.
  ┃  - Revert every mutation you make. The tree must be exactly as the runner left it,
  ┃    except your ledger commit.
  ┃  - No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
  ┃
  ┃  ### Notify the orchestrator
  ┃
  ┃  ```bash
  ┃  mkdir -p work/inbox
  ┃  cat > work/inbox/reviewer-slice-b-cycle-2.md <<EOF
  ┃  Reviewer: slice B cycle 2 done. Verdict: <ok|changes-needed>. F-72 closed: <yes|no>. Rows T4/T6 pointers updated. Filed <N> new blockers (<F-74.. or none>). LOW notes: <any>.
  ┃  EOF
  ┃  git add work/inbox && git commit -q -m "inbox: reviewer reports slice-b cycle 2"
  ┃  herdr agent prompt "w1Z:p1" "Reviewer: slice B cycle 2 done. Verdict: <ok|changes-needed>. F-72 closed: <yes|no>. Rows T4/T6 pointers updated. Filed <N> new blockers (<F-74.. or none>). LOW notes: <any>."
  ┃  ```
  ┃

     + Thought · 60ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -8 && echo "---STATUS---" && git status --short
  ┃
  ┃  (2 earlier lines)
  ┃  3bc90f9 work: fix cycle 2 — the table is the stale half, not the sentence
  ┃  1c386c5 inbox: runner reports slice-b fix 1
  ┃  8c394cd F-72: the spec's rule sentence agrees with its table
  ┃  f4751d1 work: slice B's state, the lanes' panes, and where each finding lands
  ┃  18cc2e8 work: cycle 2's review brief
  ┃  77814f1 work: the fix cycle's brief, and the orchestrator's rulings
  ┃  ---STATUS---
  ┃

     + Thought · 335ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git show 8c394cd && echo "=====================" && git show edc08b6
  ┃
  ┃  (84 earlier lines)
  ┃   |---|---|
  ┃   | recipe · special · prep note · course · menu | `acceptance_criteria` · `scope` · `evidence` · `findings` · `work_packages` |
  ┃  -| chef · maître d' · sous chef | runner · reviewer · orchestrator · lifecycle states · gate names · `status` · `review_state` |
  ┃  +| | `chef` · `sous-chef` · runner · reviewer · orchestrator · lifecycle states · gate names · `status` · `review_state` |
  ┃   | `maitre-cli` · `maitre` · the skills | |
  ┃
  ┃   ### IDs
  ┃

     + Thought · 3.8s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "maître d'\|maitre d'\|maître" --include="*.md" --include="*.go" . | grep -v '^\./research/' | grep -v '^\./work/' | grep -v '^\./docs/audit.md'
  ┃
  ┃  ./GLOSSARY.md:124:| the human-facing console | **maître d'** |
  ┃  ./docs/decisions/0004-naming.md:59:| `maitre` | Human console | the maître d' serves the Chef |
  ┃  ./docs/decisions/0004-naming.md:63:`pass`/`board`: the maître d' is the one who greets the chef and takes the order, which
  ┃  ./docs/decisions/0004-naming.md:76:Chef · maître d' · sous chef · line cook · taster · expo — defined in
  ┃  ./docs/decisions/0005-vocabulary.md:46:| Roles: chef · maître d' · sous chef · line cook · taster · expo | lifecycle states · lanes · gate names · `status` · `review_state` |
  ┃

     + Thought · 10.4s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && echo "=== diff since my last commit (6410b65) ===" && git diff 6410b65..HEAD --stat && echo "=== any change under internal/ or cmd/ ? ===" && git diff 6410b65..HEAD --stat -- internal/ cmd/ scripts/ kitchen/ go.mod && echo "(empty abov
  ┃  e = no code touched)"
  ┃
  ┃  (5 earlier lines)
  ┃   work/briefs/review-b2.md           | 86 ++++++++++++++++++++++++++++++++++
  ┃   work/inbox/runner-slice-b-fix-1.md | 68 +++++++++++++++++++++++++++
  ┃   work/inbox/runner-slice-b-fix-2.md | 73 +++++++++++++++++++++++++++++
  ┃   7 files changed, 411 insertions(+), 8 deletions(-)
  ┃  === any change under internal/ or cmd/ ? ===
  ┃  (empty above = no code touched)
  ┃

     + Thought · 335ms

     → Explored: 1 read

     + Thought · 9.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cat work/briefs/fix-b1.md && echo "=================" && cat work/briefs/fix-b2.md
  ┃
  ┃  (162 earlier lines)
  ┃  mkdir -p work/inbox
  ┃  cat > work/inbox/runner-slice-b-fix-2.md <<EOF
  ┃  Runner: slice B fix 2 done. Commits: <n>. F-72 closed in docs/03's table. Notes: <anything the reviewer should know>.
  ┃  EOF
  ┃  git add work/inbox && git commit -q -m "inbox: runner reports slice-b fix 2"
  ┃  herdr agent prompt "w1Z:p1" "Runner: slice B fix 2 done. Commits: <n>. F-72 closed in docs/03's table. Notes: <anything the reviewer should know>."
  ┃  ```
  ┃

     + Thought · 8.1s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v 2>&1 | grep -E "^(=== RUN|--- (PASS|FAIL)|ok|FAIL)"
  ┃
  ┃  (24 earlier lines)
  ┃  --- PASS: TestEveryDocumentedStepStatusIsReachable (0.07s)
  ┃  === RUN   TestAFailurePrintsItsDocumentOnStderr
  ┃  --- PASS: TestAFailurePrintsItsDocumentOnStderr (0.01s)
  ┃  === RUN   TestTheDocsCommandCountMatchesTheCode
  ┃  --- PASS: TestTheDocsCommandCountMatchesTheCode (0.00s)
  ┃  ok    github.com/lthiagol/maitre/internal/cli  (cached)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  (cached) [no tests to run]
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff 6410b65..HEAD -- internal/core/docs_test.go internal/core/planformat_test.go internal/cli/ && echo "(empty = no new wire, no code touched)" && git diff 6410b65..HEAD -- docs/decisions/ && echo "(empty = no ADR edited)"
  ┃
  ┃  (empty = no new wire, no code touched)
  ┃  (empty = no ADR edited)
  ┃

     + Thought · 1.2s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "runner\|reviewer\|orchestrator\|sous.chef\|sous chef\|\bagent\b\|line-cook\|taster\|expo\|coordinator" docs/03-domain-model.md | grep -v "^79:" | head -20
  ┃
  ┃  (13 earlier lines)
  ┃  275:    reviewed --> closed : close<br/>review.verdict_ok · review.reviewer_independent · review.clean · audit.documented · immutability.respected
  ┃  311:   one of the agent roles. The chef is the human; the sous chef, the runner, the
  ┃  312:   reviewer and the orchestrator are agents, and so is the CLI's default actor `agent`.
  ┃  314:   matches the actor's **name**, so an agent that calls itself `chef` passes it: the
  ┃  316:   principle 7 (*"Would I trust a capable agent with this prompt as it stands?"*) turned
  ┃  318:   — `line-cook`, `taster`, `expo`, `coordinator` — still name the role they stood for,
  ┃  347:| `approval.by_human` | the actor signing off is not one of the agent roles — only the chef approves | `approve` |
  ┃

     + Thought · 2.2s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp docs/03-domain-model.md /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/03c2.bak && python3 - <<'EOF'
  ┃  p='docs/03-domain-model.md'…
  ┃
  ┃  --- MUTATION: the doc stops naming 'orchestrator' (table and prose) ---
  ┃      docs_test.go:1576: the doc does not name "orchestrator" — the code refuses it as an agent
  ┃  --- FAIL: TestTheDocNamesTheRolesTheCodeRefuses (0.00s)
  ┃  REVERTED
  ┃

     + Thought · 630ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'EOF'
  ┃  p='docs/03-domain-model.md'…
  ┃
  ┃  --- MUTATION: the whole mapping table is gone; does the wire still hold on the prose alone? ---
  ┃  --- PASS: TestTheDocNamesTheRolesTheCodeRefuses (0.00s)
  ┃  REVERTED
  ┃

     + Thought · 292ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make lint && echo "LINT OK" && gofmt -l . && echo "GOFMT CLEAN" && make test 2>&1 | tail -6
  ┃
  ┃  (2 earlier lines)
  ┃  GOFMT CLEAN
  ┃  ?     github.com/lthiagol/maitre/cmd/maitre  [no test files]
  ┃  ?     github.com/lthiagol/maitre/cmd/maitre-cli  [no test files]
  ┃  ok    github.com/lthiagol/maitre/internal/cli  (cached)
  ┃  ok    github.com/lthiagol/maitre/internal/core  12.692s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  (cached)
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃

     + Thought · 238ms

     → Explored: 1 read

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "| T4 \|| T6 \|| F-72 " docs/audit.md
  ┃
  ┃  122:| T4 | feature | `internal/core/views.go` — `TreeOf`; `internal/ui/model.go`, `internal/ui/views.go` — the Chopboard tab | `internal/core/tree_test.go`, `internal/ui/chop_test.go` |
  ┃  124:| T6 | feature | `internal/ui/model.go`, `internal/ui/views.go` — the Expo tab and the automation's vital signs | `internal/ui/expo_test.go` |
  ┃

     + Thought · 690ms

     → Explored: 2 reads

     Now the ledger edits: close F-72 and fix the two stale pointers.

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    77   |---|---|---|---|                                                                                                                    77   |---|---|---|---|
  ┃    78   | **F-59** | medium | **`core.Audit.Decisions` has a shape and no writer** — it is where a decision *about the plan* would live, and 78   | **F-59** | medium | **`core.Audit.Decisions` has a shape and no writer** — it is where a decision *about the plan* would live, and
  ┃         nothing writes it. `kitchen/decisions.json` is named by the plan format and opened by nothing. The remedy is on the table in [06-         nothing writes it. `kitchen/decisions.json` is named by the plan format and opened by nothing. The remedy is on the table in [06-
  ┃         proposals.md](06-proposals.md) § 2. | `internal/core/model.go:270`; `internal/core/store.go:24` |                                         proposals.md](06-proposals.md) § 2. | `internal/core/model.go:270`; `internal/core/store.go:24` |
  ┃    79   | **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.resolved` gate can never fire** — nothing ever   79   | **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.resolved` gate can never fire** — nothing ever
  ┃         sets a question, so the list is always empty and the gate always passes. A gate nobody can fail is not a guard. | `internal/core/         sets a question, so the list is always empty and the gate always passes. A gate nobody can fail is not a guard. | `internal/core/
  ┃         model.go:31`; `internal/core/lifecycle.go:197` |                                                                                          model.go:31`; `internal/core/lifecycle.go:197` |
  ┃    80   | **F-61** | medium | **The `courses/` directory is named by the plan format, created by `Init`, and never written** — a directory   80   | **F-61** | medium | **The `courses/` directory is named by the plan format, created by `Init`, and never written** — a directory
  ┃         no writer touches. Per [06-proposals.md](06-proposals.md) § 1 the honest answer is to stop naming it: a course is the relation, not       no writer touches. Per [06-proposals.md](06-proposals.md) § 1 the honest answer is to stop naming it: a course is the relation, not
  ┃         a file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.md) |                                                              a file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.md) |
  ┃    81 - | **F-72** | medium | **The spec's vocabulary rule contradicts its own table after the role rename** — `docs/03-domain-model.md` §
  ┃         "The CPD ↔ tool mapping" says *"the **artifact and the roles** get themed names"*, but the table now keeps `runner · reviewer ·
  ┃         orchestrator` in the **Kept plain** column. The table is the right half: the decision of record ([docs/07-vocabulary.md](07-
  ┃         vocabulary.md) § "The rule") puts the agents' jobs in the plain vocabulary, and the code says the same. The sentence was not touched
  ┃         when the roles moved, so the section states the opposite of what it implements, seven lines apart. Nothing catches it: the deny-list
  ┃         tripwire only asks whether the doc *names* a role, not which column it sits in. The rule sentence and the table must say the same
  ┃         thing. | `docs/03-domain-model.md:72`; the table at `docs/03-domain-model.md:79` |
  ┃    82   | F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and the size.** It is path + mtime + size per file, 81   | F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and the size.** It is path + mtime + size per file,
  ┃         so two different contents can share one fingerprint. The tool's own writes always change the mtime, so in practice only a hand edit       so two different contents can share one fingerprint. The tool's own writes always change the mtime, so in practice only a hand edit
  ┃         that restores the mtime is missed — and the chef's `r` reloads it. | `internal/core/store.go` — `DiskFingerprint`. |                      that restores the mtime is missed — and the chef's `r` reloads it. | `internal/core/store.go` — `DiskFingerprint`. |
  ┃    83   | F-73 | low | **`kitchen/project.json` still carries `"name": "new-project"`** — the decision of record names the project           82   | F-73 | low | **`kitchen/project.json` still carries `"name": "new-project"`** — the decision of record names the project
  ┃         **maitre**, and the console's header renders the plan's name, so the repo's own console reads `maitre · new-project`. It is the           **maitre**, and the console's header renders the plan's name, so the repo's own console reads `maitre · new-project`. It is the
  ┃         *plan's* data, not the code's: slice A renamed the module and the binaries and left this file, and slice C is the slice that              *plan's* data, not the code's: slice A renamed the module and the binaries and left this file, and slice C is the slice that
  ┃         migrates the repo's own plan. | `kitchen/project.json:8` |                                                                                migrates the repo's own plan. | `kitchen/project.json:8` |
  ┃    84   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-scope=`. | Recorded so nobody calls it a regression 83   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-scope=`. | Recorded so nobody calls it a regression
  ┃         later. Not a defect. |                                                                                                                    later. Not a defect. |
  ┃    85                                                                                                                                        84
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    121   | T4 | feature | `internal/core/views.go` — `TreeOf`; `internal/ui/model.go`, `internal/ui/views.go` — the Chopboard tab |          121   | T4 | feature | `internal/core/views.go` — `TreeOf`; `internal/ui/model.go`, `internal/ui/views.go` — the Chopboard tab |
  ┃          `internal/core/tree_test.go`, `internal/ui/chop_test.go` |                                                                                `internal/core/tree_test.go`, `internal/ui/chop_test.go` |
  ┃    122   | F-55 | medium | `internal/core/views.go` — `TreeOf` draws what the walk cannot reach, as `core.Tree.Stuck` | `internal/core/      122   | F-55 | medium | `internal/core/views.go` — `TreeOf` draws what the walk cannot reach, as `core.Tree.Stuck` | `internal/core/
  ┃          tree_test.go` |                                                                                                                           tree_test.go` |
  ┃    123   | T6 | feature | `internal/ui/model.go`, `internal/ui/views.go` — the Expo tab and the automation's vital signs | `internal/ui/     123   | T6 | feature | `internal/ui/model.go`, `internal/ui/views.go` — the Expo tab and the automation's vital signs | `internal/ui/
  ┃          expo_test.go` |                                                                                                                           expo_test.go` |
  ┃    124   | F-57 | medium | `docs/03-domain-model.md` §8's table and `docs/05-human-interface.md`'s footer paragraph — the moves and the tree 124   | F-57 | medium | `docs/03-domain-model.md` §8's table and `docs/05-human-interface.md`'s footer paragraph — the moves and the tree
  ┃          are named, the two views the domain does not compute are marked, and `core.DerivedViews` owns the sets | `internal/core/docs_test.        are named, the two views the domain does not compute are marked, and `core.DerivedViews` owns the sets | `internal/core/docs_test.
  ┃          go` — the derived-views check |                                                                                                           go` — the derived-views check |
  ┃                                                                                                                                              125 + | F-72 | medium | `docs/03-domain-model.md` § "The CPD ↔ tool mapping" — the rule sentence and the table say the same thing now,
  ┃                                                                                                                                                    and both say what the decision of record says: the kitchen names the artifact and the place, the agents' jobs stay plain, and the
  ┃                                                                                                                                                    **Kept plain** cell carries `chef` · `sous-chef` · `runner` · `reviewer` · `orchestrator`. `maître d'` goes: it was the human-
  ┃                                                                                                                                                    facing console's themed name, the console is `maitre` now, and it was never a job. | — (a doc fix; the deny-list wire holds the
  ┃                                                                                                                                                    role list, nothing holds the sentence) |
  ┃    125                                                                                                                                       126
  ┃    126   ## What this ledger cannot tell you                                                                                                 127   ## What this ledger cannot tell you
  ┃    127                                                                                                                                       128
  ┃    128   **Cycles 1–7's own findings are not recoverable.** Their per-finding reports lived in a                                             129   **Cycles 1–7's own findings are not recoverable.** Their per-finding reports lived in a
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    117   | P-5 | low | `internal/cli/recipe.go` — a comment's author defaults to the actor | `internal/cli/comment_test.go` |                117   | P-5 | low | `internal/cli/recipe.go` — a comment's author defaults to the actor | `internal/cli/comment_test.go` |
  ┃    118   | P-6 | low | `AGENTS.md` — the dependency count reads 28, with the counting command written down | — |                             118   | P-6 | low | `AGENTS.md` — the dependency count reads 28, with the counting command written down | — |
  ┃    119   | P-7 | trivial | `internal/ui/model.go` — the dead `m.detail` field is gone | — |                                                  119   | P-7 | trivial | `internal/ui/model.go` — the dead `m.detail` field is gone | — |
  ┃    120   | T5 | feature | `internal/core/store.go` — `DiskFingerprint`; `internal/ui/model.go` — the poll; `cmd/maitre/main.go` — the live   120   | T5 | feature | `internal/core/store.go` — `DiskFingerprint`; `internal/ui/model.go` — the poll; `cmd/maitre/main.go` — the live
  ┃          console | `internal/ui/live_test.go` |                                                                                                    console | `internal/ui/live_test.go` |
  ┃    121 - | T4 | feature | `internal/core/views.go` — `TreeOf`; `internal/ui/model.go`, `internal/ui/views.go` — the Chopboard tab |          121 + | T4 | feature | `internal/core/views.go` — `TreeOf`; `internal/ui/model.go`, `internal/ui/views.go` — the Chopboard tab |
  ┃          `internal/core/tree_test.go`, `internal/ui/chop_test.go` |                                                                                `internal/core/tree_test.go`, `internal/ui/board_test.go` |
  ┃    122   | F-55 | medium | `internal/core/views.go` — `TreeOf` draws what the walk cannot reach, as `core.Tree.Stuck` | `internal/core/      122   | F-55 | medium | `internal/core/views.go` — `TreeOf` draws what the walk cannot reach, as `core.Tree.Stuck` | `internal/core/
  ┃          tree_test.go` |                                                                                                                           tree_test.go` |
  ┃    123   | T6 | feature | `internal/ui/model.go`, `internal/ui/views.go` — the Expo tab and the automation's vital signs | `internal/ui/     123   | T6 | feature | `internal/ui/model.go`, `internal/ui/views.go` — the Expo tab and the automation's vital signs | `internal/ui/
  ┃          expo_test.go` |                                                                                                                           expo_test.go` |
  ┃    124   | F-57 | medium | `docs/03-domain-model.md` §8's table and `docs/05-human-interface.md`'s footer paragraph — the moves and the tree 124   | F-57 | medium | `docs/03-domain-model.md` §8's table and `docs/05-human-interface.md`'s footer paragraph — the moves and the tree
  ┃          are named, the two views the domain does not compute are marked, and `core.DerivedViews` owns the sets | `internal/core/docs_test.        are named, the two views the domain does not compute are marked, and `core.DerivedViews` owns the sets | `internal/core/docs_test.
  ┃          go` — the derived-views check |                                                                                                           go` — the derived-views check |
  ┃    125   | F-72 | medium | `docs/03-domain-model.md` § "The CPD ↔ tool mapping" — the rule sentence and the table say the same thing now,    125   | F-72 | medium | `docs/03-domain-model.md` § "The CPD ↔ tool mapping" — the rule sentence and the table say the same thing now,
  ┃          and both say what the decision of record says: the kitchen names the artifact and the place, the agents' jobs stay plain, and the         and both say what the decision of record says: the kitchen names the artifact and the place, the agents' jobs stay plain, and the
  ┃          **Kept plain** cell carries `chef` · `sous-chef` · `runner` · `reviewer` · `orchestrator`. `maître d'` goes: it was the human-            **Kept plain** cell carries `chef` · `sous-chef` · `runner` · `reviewer` · `orchestrator`. `maître d'` goes: it was the human-
  ┃          facing console's themed name, the console is `maitre` now, and it was never a job. | — (a doc fix; the deny-list wire holds the           facing console's themed name, the console is `maitre` now, and it was never a job. | — (a doc fix; the deny-list wire holds the
  ┃          role list, nothing holds the sentence) |                                                                                                  role list, nothing holds the sentence) |
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    119   | P-7 | trivial | `internal/ui/model.go` — the dead `m.detail` field is gone | — |                                                  119   | P-7 | trivial | `internal/ui/model.go` — the dead `m.detail` field is gone | — |
  ┃    120   | T5 | feature | `internal/core/store.go` — `DiskFingerprint`; `internal/ui/model.go` — the poll; `cmd/maitre/main.go` — the live   120   | T5 | feature | `internal/core/store.go` — `DiskFingerprint`; `internal/ui/model.go` — the poll; `cmd/maitre/main.go` — the live
  ┃          console | `internal/ui/live_test.go` |                                                                                                    console | `internal/ui/live_test.go` |
  ┃    121   | T4 | feature | `internal/core/views.go` — `TreeOf`; `internal/ui/model.go`, `internal/ui/views.go` — the Chopboard tab |          121   | T4 | feature | `internal/core/views.go` — `TreeOf`; `internal/ui/model.go`, `internal/ui/views.go` — the Chopboard tab |
  ┃          `internal/core/tree_test.go`, `internal/ui/board_test.go` |                                                                               `internal/core/tree_test.go`, `internal/ui/board_test.go` |
  ┃    122   | F-55 | medium | `internal/core/views.go` — `TreeOf` draws what the walk cannot reach, as `core.Tree.Stuck` | `internal/core/      122   | F-55 | medium | `internal/core/views.go` — `TreeOf` draws what the walk cannot reach, as `core.Tree.Stuck` | `internal/core/
  ┃          tree_test.go` |                                                                                                                           tree_test.go` |
  ┃    123 - | T6 | feature | `internal/ui/model.go`, `internal/ui/views.go` — the Expo tab and the automation's vital signs | `internal/ui/     123 + | T6 | feature | `internal/ui/model.go`, `internal/ui/views.go` — the Expo tab and the automation's vital signs | `internal/ui/
  ┃          expo_test.go` |                                                                                                                           souschef_test.go` |
  ┃    124   | F-57 | medium | `docs/03-domain-model.md` §8's table and `docs/05-human-interface.md`'s footer paragraph — the moves and the tree 124   | F-57 | medium | `docs/03-domain-model.md` §8's table and `docs/05-human-interface.md`'s footer paragraph — the moves and the tree
  ┃          are named, the two views the domain does not compute are marked, and `core.DerivedViews` owns the sets | `internal/core/docs_test.        are named, the two views the domain does not compute are marked, and `core.DerivedViews` owns the sets | `internal/core/docs_test.
  ┃          go` — the derived-views check |                                                                                                           go` — the derived-views check |
  ┃    125   | F-72 | medium | `docs/03-domain-model.md` § "The CPD ↔ tool mapping" — the rule sentence and the table say the same thing now,    125   | F-72 | medium | `docs/03-domain-model.md` § "The CPD ↔ tool mapping" — the rule sentence and the table say the same thing now,
  ┃          and both say what the decision of record says: the kitchen names the artifact and the place, the agents' jobs stay plain, and the         and both say what the decision of record says: the kitchen names the artifact and the place, the agents' jobs stay plain, and the
  ┃          **Kept plain** cell carries `chef` · `sous-chef` · `runner` · `reviewer` · `orchestrator`. `maître d'` goes: it was the human-            **Kept plain** cell carries `chef` · `sous-chef` · `runner` · `reviewer` · `orchestrator`. `maître d'` goes: it was the human-
  ┃          facing console's themed name, the console is `maitre` now, and it was never a job. | — (a doc fix; the deny-list wire holds the           facing console's themed name, the console is `maitre` now, and it was never a job. | — (a doc fix; the deny-list wire holds the
  ┃          role list, nothing holds the sentence) |                                                                                                  role list, nothing holds the sentence) |
  ┃    126                                                                                                                                       126
  ┃    127   ## What this ledger cannot tell you                                                                                                 127   ## What this ledger cannot tell you
  ┃

     + Thought · 4.2s

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    36   | Cycle 15 | The closing review of the night: no read-side family, and a symmetric judgement | One note closed, and the shape judged 36   | Cycle 15 | The closing review of the night: no read-side family, and a symmetric judgement | One note closed, and the shape judged
  ┃         |                                                                                                                                         |
  ┃    37   | Slice A | The vocabulary rename: the module, the binaries, the recipe ids, the version | The ids and the guard hold; two mediums   37   | Slice A | The vocabulary rename: the module, the binaries, the recipe ids, the version | The ids and the guard hold; two mediums
  ┃         and four lows filed |                                                                                                                     and four lows filed |
  ┃    38   | Cycle 2 | The fixes to slice A: the usage lines, the two-way wire, the plan's stamp, the format facts | Six notes closed, and the  38   | Cycle 2 | The fixes to slice A: the usage lines, the two-way wire, the plan's stamp, the format facts | Six notes closed, and the
  ┃         fence checked by mutation |                                                                                                               fence checked by mutation |
  ┃    39   | Slice B | The vocabulary rename: the roles and the tabs | The roles resolve as written *and* as a role, so a plan written before   39   | Slice B | The vocabulary rename: the roles and the tabs | The roles resolve as written *and* as a role, so a plan written before
  ┃         the rename still cannot close on its own executor's review; the tabs are Board and Sous-Chef, and the honesty of the waiting tab          the rename still cannot close on its own executor's review; the tabs are Board and Sous-Chef, and the honesty of the waiting tab
  ┃         held; one medium filed in the spec |                                                                                                      held; one medium filed in the spec |
  ┃                                                                                                                                              40 + | Slice B · cycle 2 | The fixes to slice B: the spec's vocabulary rule, and the table under it | F-72 closed — the sentence, the
  ┃                                                                                                                                                   table and the decision of record say the same thing; the ledger's two pointers follow the rename |
  ┃    40   | The rounds inside a cycle | Fix → verify → back to the executor → repeat, until the verifier is satisfied | Three rounds in cycle  41   | The rounds inside a cycle | Fix → verify → back to the executor → repeat, until the verifier is satisfied | Three rounds in cycle
  ┃         8: the first fix introduced three new mediums, the second closed them |                                                                   8: the first fix introduced three new mediums, the second closed them |
  ┃    41                                                                                                                                        42
  ┃    42   ### The discipline every round runs                                                                                                  43   ### The discipline every round runs
  ┃    43                                                                                                                                        44
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff docs/audit.md
  ┃
  ┃  (26 earlier lines)
  ┃   | F-55 | medium | `internal/core/views.go` — `TreeOf` draws what the walk cannot reach, as `core.Tree.Stuck` | `internal/core/tree_test.go` |
  ┃  -| T6 | feature | `internal/ui/model.go`, `internal/ui/views.go` — the Expo tab and the automation's vital signs | `internal/ui/expo_test.go` |
  ┃  +| T6 | feature | `internal/ui/model.go`, `internal/ui/views.go` — the Expo tab and the automation's vital signs | `internal/ui/souschef_test.go` |
  ┃   | F-57 | medium | `docs/03-domain-model.md` §8's table and `docs/05-human-interface.md`'s footer paragraph — the moves and the tree are named, the two views the domain does not compute are marked, and `core.DerivedViews` owns the sets | `internal/core/docs_test.go` — the
  ┃  derived-views check |
  ┃  +| F-72 | medium | `docs/03-domain-model.md` § "The CPD ↔ tool mapping" — the rule sentence and the table say the same thing now, and both say what the decision of record says: the kitchen names the artifact and the place, the agents' jobs stay plain, and the **Kept plain**
  ┃  cell carries `chef` · `sous-chef` · `runner` · `reviewer` · `orchestrator`. `maître d'` goes: it was the human-facing console's themed name, the console is `maitre` now, and it was never a job. | — (a doc fix; the deny-list wire holds the role list, nothing holds the sentence)
  ┃  |
  ┃
  ┃   ## What this ledger cannot tell you
  ┃

     + Thought · 618ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && ls internal/ui/board_test.go internal/ui/souschef_test.go && git status --short && go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' 2>&1 | tail -4
  ┃
  ┃  internal/ui/board_test.go
  ┃  internal/ui/souschef_test.go
  ┃   M docs/audit.md
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.086s
  ┃  ok    github.com/lthiagol/maitre/internal/cli  (cached)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  (cached) [no tests to run]
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/audit.md && git commit -q -m "docs: F-72 is closed, and two pointers follow the rename
  ┃  …
  ┃
  ┃  df64692 docs: F-72 is closed, and two pointers follow the rename
  ┃  c4b6217 work: the remaining briefs carry the pane, the path, and what is already done
  ┃  723d991 inbox: runner reports slice-b fix 2
  ┃

  ┃
  ┃  # Wrote work/inbox/reviewer-slice-b-cycle-2.md
  ┃
  ┃   1 Reviewer: slice B cycle 2 done. Verdict: ok. F-72 closed: yes. Rows T4/T6 pointers updated. Filed 0 new blockers (none). LOW notes: the spec's mapping table still keeps the tools (`maitre-cli` · `maitre` · the skills) in the Themed column with no rule sentence covering
  ┃     them — pre-existing since the first commit, and the fix's sentence neither names nor excludes them; `maître d'` now survives only in `GLOSSARY.md` and the two old ADRs, which are slice E's; in the fixed table `chef` and `sous-chef` are backticked and `runner` · `reviewer`
  ┃     · `orchestrator` are not — cosmetic, and the deny-list wire reads the whole document either way.
  ┃   2
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/inbox && git commit -q -m "inbox: reviewer reports slice-b cycle 2" && herdr agent prompt "w1Z:p1" "Reviewer: slice B cycle 2 done. Verdict: ok. F-72 closed: yes. Rows T4/T6 pointers updated. Filed 0 new blockers (none). L
  ┃  OW notes: the mapping table still keeps the tools in the Themed column with no rule sentence covering them (pre-existing, untouched by the fix); maitre d' now survives only in GLOSSARY.md and the two old ADRs, which are slice E's." && git log --oneline -3 && git status …
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e8138fe5effeUJvY7OF96j92fE"},"agent_status":"done","completion_seq":360,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":
  ┃  false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":16,"screen_detection_skipped":true,"state_change_seq":360,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722","terminal_title":"OC | Maitre vocabulary rename handover",
  ┃  "terminal_title_stripped":"OC | Maitre vocabulary rename handover","workspace_id":"w1Z"},"type":"agent_prompted"}}
  ┃  b9e1b22 inbox: reviewer reports slice-b cycle 2
  ┃  df64692 docs: F-72 is closed, and two pointers follow the rename
  ┃  c4b6217 work: the remaining briefs carry the pane, the path, and what is already done
  ┃  --- TREE CLEAN ---
  ┃

     + Thought · 435ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff c4b6217..HEAD --stat && echo "=== the runner's tree, untouched ===" && git diff c4b6217..HEAD --stat -- internal/ cmd/ scripts/ kitchen/ go.mod docs/03-domain-model.md docs/07-vocabulary.md && echo "(empty = nothing but the le
  ┃  dger and the receipt)"
  ┃
  ┃  docs/audit.md                          | 7 ++++---
  ┃   work/inbox/reviewer-slice-b-cycle-2.md | 1 +
  ┃   2 files changed, 5 insertions(+), 3 deletions(-)
  ┃  === the runner's tree, untouched ===
  ┃  (empty = nothing but the ledger and the receipt)
  ┃

     Cycle 2 is done. Verdict: ok — F-72 is closed, no new blockers.

     Is F-72 actually closed?

     The runner fixed it in two commits (8c394cd, then edc08b6 for the residual the orchestrator ruled on). The final state of docs/03-domain-model.md § The CPD ↔ tool mapping:

     │ The rule behind the split: the kitchen names the artifact and the place; the agents' jobs and everything an agent reasons about mechanically stay plain. See docs/07-vocabulary.md § The rule.

     ┌──────────────────────────────────────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Themed                                                                                                   │ Kept plain                                                                                                                                                             │
     ├──────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ recipe · special · prep note · course · menu                                                             │ acceptance_criteria · scope · evidence · findings · work_packages                                                                                                      │
     ├──────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │                                                                                                          │ chef · sous-chef · runner · reviewer · orchestrator · lifecycle states · gate names · status · review_state                                                            │
     ├──────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ maitre-cli · maitre · the skills                                                                         │                                                                                                                                                                        │
     └──────────────────────────────────────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

     Against the decision of record (docs/07 § The rule), not the table alone:

     ┌──────────────────────────────────────────────────┬──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────────────────────┐
     │                                                  │ docs/07                                                                                                                                              │ docs/03 now                                                             │
     ├──────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────┤
     │ The rule                                         │ "The kitchen names the artifact and the place. Plain words name the machine's records and the agents' jobs."                                         │ same, in its own words                                                  │
     ├──────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────┤
     │ The jobs                                         │ chef, sous-chef, runner, reviewer, orchestrator — all plain                                                                                          │ the same five, all in Kept plain                                        │
     └──────────────────────────────────────────────────┴──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────────────────────────┘

     The sentence, the table and the decision say the same thing. maître d' goes — it was the old glossary's themed name for the human-facing console, the console is maitre (already in the table's third row), and it was never a job.

     The brief's three checks

     ┌──────────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Check                                                                        │ Result                                                                                                                                                                                             │
     ├──────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ No code changed                                                              │ git diff 6410b65..HEAD over internal/, cmd/, scripts/, kitchen/, go.mod is empty. Only docs/03-domain-model.md (8 lines) and the orchestrator's work/ files moved.                                 │
     ├──────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ No new tripwire                                                              │ internal/core/docs_test.go untouched — no wire was added.                                                                                                                                          │
     ├──────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ docs/decisions/0005-vocabulary.md not edited                                 │ untouched; maître d' and the old roles survive there as the dated record they are.                                                                                                                 │
     └──────────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

     The gate: go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v — 13/13 pass; make lint, gofmt -l . clean, make test green.

     The deny-list wire is still live and still reads the whole document. Two mutations, both reverted: renaming every orchestrator in docs/03 makes the wire fail (the doc does not name "orchestrator"); deleting the whole mapping table leaves it passing — the prose at § the roles
     names every role and every retired name on its own, so the table's cell surgery could not have dropped the last mention of anything.

     The two ghost pointers

     Updated the pointers only, the wording as filed:

     ┌────────────────────────────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Row                                                                        │ Was                                                                                             │ Now                                                                                                │
     ├────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ T4                                                                         │ internal/ui/chop_test.go                                                                        │ internal/ui/board_test.go                                                                          │
     ├────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ T6                                                                         │ internal/ui/expo_test.go                                                                        │ internal/ui/souschef_test.go                                                                       │
     └────────────────────────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────┘

     Both files exist. F-72 moved from open to closed; F-73 stays open, untouched.

     LOW notes (not filed)

     - The mapping table still keeps the tools (maitre-cli · maitre · the skills) in the Themed column with no rule sentence covering them — pre-existing since the first commit; the fix's sentence neither names nor excludes them.
     - maître d' now survives only in GLOSSARY.md and the two old ADRs — slice E's sweep.
     - In the fixed table chef and sous-chef are backticked and runner · reviewer · orchestrator are not — cosmetic, and the wire reads the whole document either way.

     Ledger committed (df64692), receipt committed (b9e1b22), orchestrator notified. The tree is exactly as the runner left it plus those two commits.

     Build · MiniMax-M3.1-Flash · 1m 44s · 170.9 tok/s

  ┃
  ┃
  ┃
  ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
  ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
  ~/Code/Pessoal/maitre:master                                                                                                                                                                                                                               178.3K (18%)  ctrl+p commands

```
