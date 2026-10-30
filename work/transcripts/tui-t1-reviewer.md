# The reviewer's lane — T1, captured 2026-10-30

Dumped from pane `w1G:p4` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that judged T1: three cycles, and the findings they closed —
F-99..F-105.
The lane that wrote it is in [tui-t1-runner.md](tui-t1-runner.md).

    ┃  (5 earlier lines)
    ┃   5 files changed, 133 insertions(+), 26 deletions(-)
    ┃  === docs/03 untouched?
    ┃  (empty = untouched)
    ┃  === cli untouched?
    ┃  (empty = untouched)
    ┃  go build -o bin/maitre ./cmd/maitre
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && W=/private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/t1c3 && rm -rf "$W" && mkdir -p "$W/demo2" && cd "$W/demo2" && /Users/thiago
    ┃  /Code/Pessoal/maitre/bin/maitre-cli init >/dev/null…
    ┃
    ┃  (3 earlier lines)
    ┃  === state 2: the file names it empty
    ┃  CLI:      'demo2'
    ┃  console:   ⋈ maitre · demo2
    ┃  === state 3: the file is missing
    ┃  CLI:      'demo2'
    ┃  console:   ⋈ maitre · demo2
    ┃

       + Thought · 1.0s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'EOF'
    ┃  p='internal/core/store.go'…
    ┃
    ┃  (5 earlier lines)
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/core  22.563s
    ┃  --- FAIL: TestTheNameHasOneAnswer (0.80s)
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  78.867s
    ┃  FAIL
    ┃

       + Thought · 1.3s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git checkout -- internal/core/store.go && git status --short && echo REVERTED && echo "=== humanPath call sites:" && grep -rn "humanPath" intern
    ┃  al/ui/*.go | grep -v _test && echo "=== any other path rendered? store.Root / Root in ui:" && grep -rn "store.Root\|m.Root\|filepath\." internal/ui/*.go | grep -v _test
    ┃
    ┃  (2 earlier lines)
    ┃  internal/ui/identity.go:50:// humanPath shortens a path for a human: the home directory becomes `~`.
    ┃  internal/ui/identity.go:58:func humanPath(path string) string {
    ┃  internal/ui/views.go:1030:  add("  " + m.styDim.Render("root        ") + humanPath(m.store.Root))
    ┃  === any other path rendered? store.Root / Root in ui:
    ┃  internal/ui/identity.go:67:  if strings.HasPrefix(path, home+string(filepath.Separator)) {
    ┃  internal/ui/views.go:1030:  add("  " + m.styDim.Render("root        ") + humanPath(m.store.Root))
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make build 2>&1|tail -1 && W=/private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/t1c3 && cd "$W" && mkdir -p "we
    ┃  ird" && cd "we…
    ┃
    ┃  (3 earlier lines)
    ┃    root        /private/var/folders/5l/jx0d52nd65v87ltmtrwwt\
    ┃  f5h0000gn/T/opencode/t1c3/we ird/kitchen$
    ┃    recipes     0$
    ┃    backlog     0$
    ┃  === the header (also one row):
    ┃   ⋈ maitre · we ird$
    ┃

       + Thought · 2.0s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && row=$(grep -n "^| \*\*T1\*\*" work/tui/README.md) && echo "$row" | sed 's/.*\(✅ closed[^|]*\)|\?.*//' >/dev/null; python3 - <<'EOF'
    ┃  import re,subprocess…
    ┃
    ┃  (9 earlier lines)
    ┃  b4a3052 commit fix(ui): the row labelled root shows the root
    ┃  55bdc81 commit docs: the window title carries the project, and the brief's
    ┃  8a6866b commit fix(core): the name's rule lives in one place
    ┃  cf8c05e commit fix(ui): the path the console renders is sanitised
    ┃  === the protocol line:
    ┃  35:did not finish its report. **The count is the number of ids the row names**, and the
    ┃

       + Thought · 5.2s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "GetProject()\|LoadProject()" --include="*.go" . | grep -v "_test.go" ; echo "=== the rule in the spec?"; grep -rn "project.json" docs/
    ┃  *.md | head -20
    ┃
    ┃  (10 earlier lines)
    ┃   still carries a control character*, *the header does not carry the sanitised name*, *the header is 41 rows, want 40*), and the proof by command: a plan under a directory named `we\nird`
    ┃  renders `⋈ maitre · we ird` on **one** row while `project.json` still holds `"we\nird"`. | `internal/ui/identity.go` — `projectName`; `internal/ui/identity_test.go` —
    ┃  `TestTheNameIsSanitised` |
    ┃  docs/audit.md:204:| F-100 | medium | **The plan's name has one answer — and it had three.** A plan that names itself answered with its name, a plan without a `project.json` answered with
    ┃  the project's directory name, and a plan whose `project.json` named it **empty** answered with the **absolute path** of the plan directory — which is what landed in the window title.
    ┃  *Closed in* `internal/core/store.go` — `LoadProject` and the new `Store.ProjectName` answer the same for an empty name as for a missing file, and **a name is never a path**; `internal/
    ┃  ui/identity.go` — the console **asks** the domain (`m.store.ProjectName()`) instead of answering the question a second time. The writer is untouched: `Init` writes the name it wrote
    ┃  before, and the plan format on disk does not move. The proof by command: the empty-name plan and the plan without a `project.json` both show `⋈ maitre · demo`, and `maitre-cli plan`
    ┃  prints the same name. **The rule is written in two places and only one copy is pinned — that is F-103.** | `internal/core/store.go` — `LoadProject`, `ProjectName`, `defaultProjectName`;
    ┃  `internal/core/store_test.go` — `TestProjectNameHasOneAnswer`; `internal/ui/identity_test.go` — `TestTheNameHasOneAnswer` |
    ┃

       + Thought · 11.7s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'EOF'
    ┃  p='internal/ui/identity.go'…
    ┃
    ┃  59:  // MUTATION: the path is not sanitised.
    ┃  --- FAIL: TestTheRootRowIsOneRow (0.66s)
    ┃      identity_test.go:288: the root row does not carry the sanitised path
    ┃      identity_test.go:291: the root row carries a newline
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  0.791s
    ┃  FAIL
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git checkout -- internal/ui/identity.go && git status --short && git diff --stat && echo REVERTED && make check 2>&1 | tail -16
    ┃
    ┃  (10 earlier lines)
    ┃  GO=go scripts/dep-budget.sh
    ┃    internal/core: 0 external modules (ceiling 0)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check-race 2>&1 | tail -18
    ┃
    ┃   Background
    ┃

       + Thought · 2.9s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "defaultProjectName\|projectName(" internal/core/*.go | grep -v _test && echo "=== the core rule, read once more:" && sed -n '105,165p'
    ┃   internal/core/store.go
    ┃
    ┃  (64 earlier lines)
    ┃    }
    ┃    return doc, nil
    ┃  }
    ┃
    ┃  // defaultProjectName is the name a plan that names nothing answers with: the
    ┃  // project's directory.
    ┃

       + Thought · 19.8s

    ┃
    ┃  $ sleep 240; tail -16 /Users/thiago/.local/share/opencode/shell/d375263be99b00fa738c01ef792930354b838569/sh_18cf2040f001Qgfj7M3w236TtI.out 2>/dev/null; pgrep -fl "go test" | head -2
    ┃
    ┃  (9 earlier lines)
    ┃  GO=go scripts/dep-budget.sh
    ┃    internal/core: 0 external modules (ceiling 0)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃

       ↳ Shell finished · cd /Users/thiago/Code/Pessoal/maitre && make check-race 2>&1 | tail -18

       + Thought · 505ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make goldens 2>&1|tail -1 && git status --short && git diff --stat && echo "GOLDENS UNCHANGED / TREE CLEAN" && grep -n "DomainVersion = " intern
    ┃  al/core/lifecycle.go && grep -n "TabNames = " internal/ui/model.go && grep -n "window title carries the project" docs/05-human-interface.md && grep -rn "⋈" --include="*.go" . && e…
    ┃
    ┃  (2 earlier lines)
    ┃  225:const DomainVersion = 7
    ┃  32:var TabNames = []string{"Recipes", "Backlog", "Path", "Board", "Log", "Sous-Chef", "Settings"}
    ┃  234:**The window title carries the project.** It is the row the operating system puts on
    ┃  ./internal/brand/mark.go:14:const Compact = "⋈"
    ┃  === docs/03 in this cycle:
    ┃  (empty)
    ┃

       Everything holds. Now the ledger — the cycle-3 row, the three closed findings, and the notes:

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    60   | Slice H · cycle 5 | The last cycle of slice H: the gate sees a golden nobody          60   | Slice H · cycle 5 | The last cycle of slice H: the gate sees a golden nobody
    ┃         committed, and every message says what to do | The goldens are checked **three ways** —      committed, and every message says what to do | The goldens are checked **three ways**
    ┃         *unchanged*, *claimed by a test*, *in the commit*. The commit check is the new one:          — *unchanged*, *claimed by a test*, *in the commit*. The commit check is the new one:
    ┃         `git diff` does not see an untracked file, so the files are **intent-to-add** before         `git diff` does not see an untracked file, so the files are **intent-to-add** before
    ┃         the diff and a new golden shows up in it as a change; the index is put back as it was (      the diff and a new golden shows up in it as a change; the index is put back as it was
    ┃         `git rm --cached --force`, not `git reset`, so **a staged deletion stays staged**). The      (`git rm --cached --force`, not `git reset`, so **a staged deletion stays staged**).
    ┃         reviewer's probes: a new golden **no test claims** → the check fails, names the file         The reviewer's probes: a new golden **no test claims** → the check fails, names the
    ┃         and says which cause (*"no test source mentions it, so no test claims it"*) and its          file and says which cause (*"no test source mentions it, so no test claims it"*) and
    ┃         remedy (*"write a test that pins it, or delete the file"*); a new golden **a test            its remedy (*"write a test that pins it, or delete the file"*); a new golden **a test
    ┃         claims but nobody committed** → the check fails and names it (*"these golden files are       claims but nobody committed** → the check fails and names it (*"these golden files are
    ┃         not in the commit"*); a golden **staged for deletion** → the same failure, and the           not in the commit"*); a golden **staged for deletion** → the same failure, and the
    ┃         index afterwards is the index before (*"D "* plus the file on disk, untracked). The          index afterwards is the index before (*"D "* plus the file on disk, untracked). The
    ┃         messages were brought to the rule: each cause names **its own** remedy — *"name the          messages were brought to the rule: each cause names **its own** remedy — *"name the
    ┃         test so the filter matches it — its name must carry 'Golden' — or extend the filter"*,       test so the filter matches it — its name must carry 'Golden' — or extend the filter"*,
    ┃         and *"write a test that pins it, or delete the file"* — and the reviewer could follow        and *"write a test that pins it, or delete the file"* — and the reviewer could follow
    ┃         each one and be unblocked without reading the script. The doc check is scoped to the         each one and be unblocked without reading the script. The doc check is scoped to the
    ┃         **pin's shape**: the release version (`staticcheck 2026.2.2`) passes and a mention           **pin's shape**: the release version (`staticcheck 2026.2.2`) passes and a mention
    ┃         carrying the pin's shape with the wrong value (`v0.8.1`) fails. The cycle touched **no       carrying the pin's shape with the wrong value (`v0.8.1`) fails. The cycle touched **no
    ┃         Go source** — **no blockers filed** |                                                        Go source** — **no blockers filed** |
    ┃    61   | Slice T1 | The identity: the window title carries the project; the header's name      61   | Slice T1 | The identity: the window title carries the project; the header's name
    ┃         stops being a footnote; `--project` overrides the name for the session; the Settings         stops being a footnote; `--project` overrides the name for the session; the Settings
    ┃         tab shows the plan's path | The name has **one** source — `Model.projectName` in             tab shows the plan's path | The name has **one** source — `Model.projectName` in
    ┃         `internal/ui/identity.go` — and the header and the window title both read it; the title      `internal/ui/identity.go` — and the header and the window title both read it; the
    ┃         reads `maitre · <project>`, and a plan whose `project.json` has no name falls back to        title reads `maitre · <project>`, and a plan whose `project.json` has no name falls
    ┃         its root. `--project` reaches all three constructors (`New`, `NewLive`, `NewScripted`)       back to its root. `--project` reaches all three constructors (`New`, `NewLive`,
    ┃         and `--render` and `--keys` honour it, because they build the same model — and it            `NewScripted`) and `--render` and `--keys` honour it, because they build the same
    ┃         **writes nothing**: after a run against a copy of a plan, `project.json` and `activity.      model — and it **writes nothing**: after a run against a copy of a plan, `project.
    ┃         json` are **byte-for-byte unchanged**. The header's name wears the accent (`39`) where       json` and `activity.json` are **byte-for-byte unchanged**. The header's name wears the
    ┃         the faint (`240`) read as a footnote, and the mark, the product's name and the               accent (`39`) where the faint (`240`) read as a footnote, and the mark, the product's
    ┃         separator did not move. The Settings tab shows `~/Code/Pessoal/maitre` — the                 name and the separator did not move. The Settings tab shows `~/Code/Pessoal/maitre` —
    ┃         **project's** directory, not the plan directory inside it — and the themes are still         the **project's** directory, not the plan directory inside it — and the themes are
    ┃         listed. **The goldens did not move**: `make goldens` leaves `git diff` empty, because a      still listed. **The goldens did not move**: `make goldens` leaves `git diff` empty,
    ┃         style is not text and the Settings tab is not in the golden. The scope held: `internal/      because a style is not text and the Settings tab is not in the golden. The scope held:
    ┃         core` and `internal/cli` are untouched, `core.DomainVersion` is still **7**, the budget      `internal/core` and `internal/cli` are untouched, `core.DomainVersion` is still **7**,
    ┃         is still **0 / 0 / 28**, the seven tabs kept their names, the mark's glyph is still          the budget is still **0 / 0 / 28**, the seven tabs kept their names, the mark's glyph
    ┃         only in `internal/brand` — **no blockers filed** |                                           is still only in `internal/brand` — **no blockers filed** |
    ┃    62   | Slice T1 · cycle 1 | The review of the identity: four mutations, the write-nothing    62   | Slice T1 · cycle 1 | The review of the identity: four mutations, the write-nothing
    ┃         proof, and the gate | Each mutation fails a test and was reverted, and the tree was          proof, and the gate | Each mutation fails a test and was reverted, and the tree was
    ┃         left exactly as the runner left it: the hardcoded `"maitre"` fails **four** window-          left exactly as the runner left it: the hardcoded `"maitre"` fails **four** window-
    ┃         title tests; a `--project` that writes the name into `project.json` fails *"the              title tests; a `--project` that writes the name into `project.json` fails *"the
    ┃         override writes nothing"* and *"the plan keeps its own name"*; the name back in              override writes nothing"* and *"the plan keeps its own name"*; the name back in
    ┃         `styFaint` fails *"the header's name is not a footnote"*; a title that computes its own      `styFaint` fails *"the header's name is not a footnote"*; a title that computes its
    ┃         name fails *"the override reaches the console"*. The write-nothing claim was proved          own name fails *"the override reaches the console"*. The write-nothing claim was
    ┃         **by command** and not by the diff: a copy of a plan, `maitre --render --project …`          proved **by command** and not by the diff: a copy of a plan, `maitre --render --
    ┃         over it, and `cmp` on both `project.json` and `activity.json`. `make check` and `make        project …` over it, and `cmp` on both `project.json` and `activity.json`. `make check`
    ┃         check-race` are green, `make goldens` leaves the tree clean, and the state row is            and `make check-race` are green, `make goldens` leaves the tree clean, and the state
    ┃         **in** the slice's commits — **no blockers filed** |                                         row is **in** the slice's commits — **no blockers filed** |
    ┃    63   | Slice T1 · cycle 2 | The review of the fix cycle: the four findings, the two          63   | Slice T1 · cycle 2 | The review of the fix cycle: the four findings, the two
    ┃         observations, and the `docs/03` question | F-99..F-102 are closed, and each was proved       observations, and the `docs/03` question | F-99..F-102 are closed, and each was proved
    ┃         by a mutation: the `sanitize` call taken out of the name's path fails                        by a mutation: the `sanitize` call taken out of the name's path fails
    ┃         `TestTheNameIsSanitised` three times over; the `root` row put back on the parent fails       `TestTheNameIsSanitised` three times over; the `root` row put back on the parent fails
    ┃         `TestTheSettingsTabShowsThePlanRoot`; the empty-name branch taken out of `LoadProject`       `TestTheSettingsTabShowsThePlanRoot`; the empty-name branch taken out of `LoadProject`
    ┃         leaves **the whole suite green** — and that is **F-103**, the one blocker filed. The         leaves **the whole suite green** — and that is **F-103**, the one blocker filed. The
    ┃         name has one answer in the code, but the rule is written **twice** and only one copy is      name has one answer in the code, but the rule is written **twice** and only one copy
    ┃         pinned: `maitre-cli plan` then prints `"Name": ""` while the console shows the               is pinned: `maitre-cli plan` then prints `"Name": ""` while the console shows the
    ┃         project's directory name. The **path row** is judged a note and not a blocker: the path      project's directory name. The **path row** is judged a note and not a blocker: the
    ┃         is not the plan's content, the unsanitised path **predates** the slice, and the              path is not the plan's content, the unsanitised path **predates** the slice, and the
    ┃         consequence is one extra visual row in a settings list — the header and the window           consequence is one extra visual row in a settings list — the header and the window
    ┃         title both survive. The **state row** is judged a note as well: the count should count       title both survive. The **state row** is judged a note as well: the count should count
    ┃         what the list names, or the row should say the count excludes the commits that wrote         what the list names, or the row should say the count excludes the commits that wrote
    ┃         it. The plan's name does **not** belong in `docs/03` § 8: that table is the derived          it. The plan's name does **not** belong in `docs/03` § 8: that table is the derived
    ┃         views — *computed, never stored* — and the name is a **stored** fact the domain reads,       views — *computed, never stored* — and the name is a **stored** fact the domain reads,
    ┃         like `LoadProject` and `ListRecipes`, which the table does not carry either. `make           like `LoadProject` and `ListRecipes`, which the table does not carry either. `make
    ┃         check` and `make check-race` are green, the goldens did not move, and the tree was left      check` and `make check-race` are green, the goldens did not move, and the tree was
    ┃         exactly as the runner left it — **one blocker filed** (F-103) |                              left exactly as the runner left it — **one blocker filed** (F-103) |
    ┃                                                                                                 64 + | Slice T1 · cycle 3 | The **last** cycle of T1: the blocker, the two notes, and the
    ┃                                                                                                      question that is the chef's | F-103..F-105 are closed, each proved by a mutation and
    ┃                                                                                                      each reverted: the rule taken out of the one copy turns the suite **red** in both
    ┃                                                                                                      `internal/core` and `internal/ui` — **three** tests, as the runner claimed; the
    ┃                                                                                                      `sanitize` call taken out of `humanPath` fails the root-row test twice over; the state
    ┃                                                                                                      row's count is the number of ids the row names, and the **protocol section** carries
    ┃                                                                                                      the shape in one line. The two surfaces agree in **all three states** — the plan names
    ┃                                                                                                      itself, the file names it empty, the file is missing — and the CLI's behaviour
    ┃                                                                                                      **changed** to match the console's, which is the point of F-103. The rule is written
    ┃                                                                                                      down in the **code** (`core.Store.projectName`'s doc comment) and **nowhere in the
    ┃                                                                                                      spec**: `docs/03` § 6 names `project.json` as *project meta · charter · settings* and
    ┃                                                                                                      says nothing about the name's fallback, and `docs/05` does not either — where it
    ┃                                                                                                      belongs is the **chef's** call, and **nothing was written into `docs/03`**. The
    ┃                                                                                                      goldens did not move, the budget is still **0 / 0 / 28**, `core.DomainVersion` is
    ┃                                                                                                      still **7**, the seven tabs kept their names, the CLI's surface did not change, and
    ┃                                                                                                      the tree was left exactly as the runner left it — **no blockers filed** |
    ┃    64   | The rounds inside a cycle | Fix → verify → back to the executor → repeat, until the   65   | The rounds inside a cycle | Fix → verify → back to the executor → repeat, until the
    ┃         verifier is satisfied | Three rounds in cycle 8: the first fix introduced three new          verifier is satisfied | Three rounds in cycle 8: the first fix introduced three new
    ┃         mediums, the second closed them |                                                            mediums, the second closed them |
    ┃    65                                                                                           66
    ┃    66   ### The discipline every round runs                                                     67   ### The discipline every round runs
    ┃    67                                                                                           68
    ┃

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    112   |---|---|---|---|                                                                      112   |---|---|---|---|
    ┃    113   | **F-59** | medium | **`core.Audit.Decisions` has a shape and no writer** — it is     113   | **F-59** | medium | **`core.Audit.Decisions` has a shape and no writer** — it is
    ┃          where a decision *about the plan* would live, and nothing writes it. `kitchen/               where a decision *about the plan* would live, and nothing writes it. `kitchen/
    ┃          decisions.json` is named by the plan format and opened by nothing. The remedy is on          decisions.json` is named by the plan format and opened by nothing. The remedy is on
    ┃          the table in [06-proposals.md](06-proposals.md) § 2. | `internal/core/model.go:270`;         the table in [06-proposals.md](06-proposals.md) § 2. | `internal/core/model.go:270`;
    ┃          `internal/core/store.go:24` |                                                                `internal/core/store.go:24` |
    ┃    114   | **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.   114   | **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.
    ┃          resolved` gate can never fire** — nothing ever sets a question, so the list is always        resolved` gate can never fire** — nothing ever sets a question, so the list is always
    ┃          empty and the gate always passes. A gate nobody can fail is not a guard. | `internal/        empty and the gate always passes. A gate nobody can fail is not a guard. | `internal/
    ┃          core/model.go:31`; `internal/core/lifecycle.go:197` |                                        core/model.go:31`; `internal/core/lifecycle.go:197` |
    ┃    115   | **F-61** | medium | **The `courses/` directory is named by the plan format, created  115   | **F-61** | medium | **The `courses/` directory is named by the plan format, created
    ┃          by `Init`, and never written** — a directory no writer touches. Per [06-proposals.md](       by `Init`, and never written** — a directory no writer touches. Per [06-proposals.md]
    ┃          06-proposals.md) § 1 the honest answer is to stop naming it: a course is the relation,       (06-proposals.md) § 1 the honest answer is to stop naming it: a course is the
    ┃          not a file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.md) |             relation, not a file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.
    ┃                                                                                                       md) |
    ┃    116 - | **F-103** | medium | **The domain's one-answer rule for the plan's name is written
    ┃          twice, and only one copy is pinned.** `core.Store.LoadProject` answers an empty name
    ┃          with the project's directory name, and `core.Store.ProjectName` answers the same
    ┃          question again for itself — so the plan's name has two answers, and the test only
    ┃          holds the second. The reviewer's mutation — taking the empty-name branch out of
    ┃          `LoadProject` — leaves **the whole suite green** (every package), and then the two
    ┃          shipped surfaces disagree: `maitre-cli plan` prints `"Name": ""` while the console
    ┃          shows the project's directory name. `LoadProject` has **no test of its own**: the only
    ┃          test that reads it (`TestThePlanKeepsItsOwnName`) reads a plan that names itself. The
    ┃          remedy: pin `LoadProject`'s answer in `internal/core/store_test.go` for the same two
    ┃          cases `ProjectName` is pinned for, or drop `ProjectName`'s redundant fallback so the
    ┃          regression shows through it. | `internal/core/store.go` — `LoadProject`,
    ┃          `ProjectName`; `internal/core/store_test.go` — `TestProjectNameHasOneAnswer` |
    ┃    117   | F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and   116   | F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and
    ┃          the size.** It is path + mtime + size per file, so two different contents can share          the size.** It is path + mtime + size per file, so two different contents can share
    ┃          one fingerprint. The tool's own writes always change the mtime, so in practice only a        one fingerprint. The tool's own writes always change the mtime, so in practice only a
    ┃          hand edit that restores the mtime is missed — and the chef's `r` reloads it. |               hand edit that restores the mtime is missed — and the chef's `r` reloads it. |
    ┃          `internal/core/store.go` — `DiskFingerprint`. |                                              `internal/core/store.go` — `DiskFingerprint`. |
    ┃    118   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-      117   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-
    ┃          scope=`. | Recorded so nobody calls it a regression later. Not a defect. |                   scope=`. | Recorded so nobody calls it a regression later. Not a defect. |
    ┃    119                                                                                          118
    ┃    120   ### Notes left for a later round                                                       119   ### Notes left for a later round
    ┃
    ┃
    ┃

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    199   | F-95 | medium | **The wire that says the art is declared only in `internal/brand`    199   | F-95 | medium | **The wire that says the art is declared only in `internal/brand`
    ┃          reads only `.go` files, and skips the brand package itself — so a second copy of the         reads only `.go` files, and skips the brand package itself — so a second copy of the
    ┃          drawing passes.** *Closed in* `internal/brand/mark_test.go` — the wire is now **two          drawing passes.** *Closed in* `internal/brand/mark_test.go` — the wire is now **two
    ┃          rules, one per thing**, and the exemption is **per declaration, found by parsing**:          rules, one per thing**, and the exemption is **per declaration, found by parsing**:
    ┃          `checkArtScope` finds the file that declares `Full` and `Half` (`declarationFile`,           `checkArtScope` finds the file that declares `Full` and `Half` (`declarationFile`,
    ┃          `declares`) and exempts **that file**, not a package; the **drawings** are matched by        `declares`) and exempts **that file**, not a package; the **drawings** are matched by
    ┃          their **lines** in **every** file of any type; the **compact glyph** is refused in any       their **lines** in **every** file of any type; the **compact glyph** is refused in
    ┃          **`.go`** file but not the declaration. Every directory that is not read is named and        any **`.go`** file but not the declaration. Every directory that is not read is named
    ┃          is a **record**: `.git` (the repository's own objects), `bin` (build output — the            and is a **record**: `.git` (the repository's own objects), `bin` (build output — the
    ┃          binaries carry the art as string literals), `skills` (vendored third-party), `work` (        binaries carry the art as string literals), `skills` (vendored third-party), `work` (
    ┃          the night's scratch: the briefs, the transcripts and the receipts, which quote what          the night's scratch: the briefs, the transcripts and the receipts, which quote what
    ┃          the console shows). The reviewer's two mutations now **fail** — a copy in `internal/         the console shows). The reviewer's two mutations now **fail** — a copy in `internal/
    ┃          brand/zz_probe.go` and a copy in `internal/ui/testdata/zz_probe.txt` — and a third,          brand/zz_probe.go` and a copy in `internal/ui/testdata/zz_probe.txt` — and a third,
    ┃          the glyph pasted into another code file, fails too. The **withdrawn claim** is               the glyph pasted into another code file, fails too. The **withdrawn claim** is
    ┃          withdrawn: what is true is that the *drawings* appear nowhere else, and the glyph is         withdrawn: what is true is that the *drawings* appear nowhere else, and the glyph is
    ┃          in the brand package, in the rendered output and in the card's art. | `internal/brand/       in the brand package, in the rendered output and in the card's art. | `internal/
    ┃          mark_test.go` — `TestTheArtIsDeclaredOnlyHere`, `TestTheArtScopeWireCanFail`,                brand/mark_test.go` — `TestTheArtIsDeclaredOnlyHere`, `TestTheArtScopeWireCanFail`,
    ┃          `checkArtScope`, `declarationFile`, `drawingLines` |                                         `checkArtScope`, `declarationFile`, `drawingLines` |
    ┃    200   | F-96 | medium | **The wire that keeps the domain free of the mark is defeated by any 200   | F-96 | medium | **The wire that keeps the domain free of the mark is defeated by
    ┃          import spelling other than a bare grouped one — and a single-line import is the one          any import spelling other than a bare grouped one — and a single-line import is the
    ┃          `gofmt` writes.** *Closed in* `internal/brand/mark_test.go` — the imports are                one `gofmt` writes.** *Closed in* `internal/brand/mark_test.go` — the imports are
    ┃          **parsed, not read**: `importsIn` walks the files with `go/parser` (`parser.                 **parsed, not read**: `importsIn` walks the files with `go/parser` (`parser.
    ┃          ImportsOnly`) and reads `file.Imports`, so **every** spelling is caught — single-line,       ImportsOnly`) and reads `file.Imports`, so **every** spelling is caught — single-
    ┃          grouped, aliased and blank — because a wire that catches one spelling waits for the          line, grouped, aliased and blank — because a wire that catches one spelling waits for
    ┃          other. The scope is in the comment: the rule is about the code that **ships**, so the        the other. The scope is in the comment: the rule is about the code that **ships**, so
    ┃          wire reads the **non-test** files of `internal/core`. The reviewer's two mutations now       the wire reads the **non-test** files of `internal/core`. The reviewer's two
    ┃          **fail** the wire. | `internal/brand/mark_test.go` — `TestTheMarkIsNotTheDomain`,            mutations now **fail** the wire. | `internal/brand/mark_test.go` —
    ┃          `TestTheBoundaryWireCanFail`, `importsIn` |                                                  `TestTheMarkIsNotTheDomain`, `TestTheBoundaryWireCanFail`, `importsIn` |
    ┃    201   | F-97 | medium | **The scope wire exempts the whole declaration *file*, so a second   201   | F-97 | medium | **The scope wire exempts the whole declaration *file*, so a second
    ┃          copy of the drawing pasted into the very file that declares it is invisible — and the        copy of the drawing pasted into the very file that declares it is invisible — and the
    ┃          exemption is per file, not per declaration.** *Closed in* `internal/brand/mark_test.         exemption is per file, not per declaration.** *Closed in* `internal/brand/mark_test.
    ┃          go` — the wire's question is now **is every occurrence of the drawing the canonical          go` — the wire's question is now **is every occurrence of the drawing the canonical
    ┃          declaration?**: `findCanonicalDeclaration` parses the package and returns the **byte         declaration?**: `findCanonicalDeclaration` parses the package and returns the **byte
    ┃          range of the two `var` literals** that hold the drawings, and `maskSpan` takes those         range of the two `var` literals** that hold the drawings, and `maskSpan` takes those
    ┃          ranges out of the text before the search (`drawingDrift`). So the drawing may be in          ranges out of the text before the search (`drawingDrift`). So the drawing may be in
    ┃          the literal that holds it, and **nowhere else** — not in another file and not in             the literal that holds it, and **nowhere else** — not in another file and not in
    ┃          another corner of the declaring file. The reviewer's probe — a copy appended to              another corner of the declaring file. The reviewer's probe — a copy appended to
    ┃          `internal/brand/mark.go` **beside** the declaration — now **fails the wire**, and the        `internal/brand/mark.go` **beside** the declaration — now **fails the wire**, and the
    ┃          whole suite fails with it. A copy of the drawing also now fails in a golden, a doc,          whole suite fails with it. A copy of the drawing also now fails in a golden, a doc,
    ┃          the plan (`kitchen/`), `cmd/`, `images/` and a non-`.go` file, while the legitimate          the plan (`kitchen/`), `cmd/`, `images/` and a non-`.go` file, while the legitimate
    ┃          copies still pass: the `.txt` sources, the declaration itself, and for the glyph the         copies still pass: the `.txt` sources, the declaration itself, and for the glyph the
    ┃          rendered output, the art and the prose. The **glyph's scope is a category, not an            rendered output, the art and the prose. The **glyph's scope is a category, not an
    ┃          extension**: `categoryOf` tells a file what it **is** — `code` (`.go`, `.py`, `.sh`,         extension**: `categoryOf` tells a file what it **is** — `code` (`.go`, `.py`, `.sh`,
    ┃          `.Makefile`…), `output` (`.golden`), `art` (under `images/`), `prose` (the rest) — so        `.Makefile`…), `output` (`.golden`), `art` (under `images/`), `prose` (the rest) — so
    ┃          the glyph pasted into `scripts/console-shot.py` now fails, while the goldens, the            the glyph pasted into `scripts/console-shot.py` now fails, while the goldens, the
    ┃          wallpaper and the brand sheet pass. The four exclusions are still excluded, each for         wallpaper and the brand sheet pass. The four exclusions are still excluded, each for
    ┃          the reason the wire's comment gives: `.git` is the repository's own objects, `bin` is        the reason the wire's comment gives: `.git` is the repository's own objects, `bin` is
    ┃          build output, `skills` is vendored third-party, and `work` is the night's scratch,           build output, `skills` is vendored third-party, and `work` is the night's scratch,
    ┃          whose briefs, transcripts and receipts quote what the console shows, art and all. |          whose briefs, transcripts and receipts quote what the console shows, art and all. |
    ┃          `internal/brand/mark_test.go` — `findCanonicalDeclaration`, `maskSpan`,                      `internal/brand/mark_test.go` — `findCanonicalDeclaration`, `maskSpan`,
    ┃          `drawingDrift`, `glyphDrift`, `categoryOf`, `TestTheDrawingBesideTheDeclarationFails`,       `drawingDrift`, `glyphDrift`, `categoryOf`,
    ┃          `TestTheArtScopeWireCanFail` |                                                               `TestTheDrawingBesideTheDeclarationFails`, `TestTheArtScopeWireCanFail` |
    ┃    202   | F-98 | medium | **The slow gate could pass having tested nothing — `make test-race`  202   | F-98 | medium | **The slow gate could pass having tested nothing — `make test-race`
    ┃          dropped the filter it was given.** *Closed in* `Makefile` — **one place decides how a        dropped the filter it was given.** *Closed in* `Makefile` — **one place decides how a
    ┃          test run is invoked**: `test` and `test-race` are two targets of **one recipe**, and         test run is invoked**: `test` and `test-race` are two targets of **one recipe**, and
    ┃          they differ in **the flags they add** (`TEST_FLAGS`) and in nothing else — the filter        they differ in **the flags they add** (`TEST_FLAGS`) and in nothing else — the filter
    ┃          and the package always reach the script, and the zero-match check always applies. The        and the package always reach the script, and the zero-match check always applies. The
    ┃          reviewer's two probes: `make test-race PKG=./internal/brand RUN='-run NoSuchTest'` now       reviewer's two probes: `make test-race PKG=./internal/brand RUN='-run NoSuchTest'`
    ┃          **fails** with *"the filter matched no test — a green that proves nothing"*, and             now **fails** with *"the filter matched no test — a green that proves nothing"*, and
    ┃          `RUN='-run TestTheDrawingsMatchTheSource'` **runs that test** (`=== RUN                      `RUN='-run TestTheDrawingsMatchTheSource'` **runs that test** (`=== RUN
    ┃          TestTheDrawingsMatchTheSource`, twice — `-count=2`) and passes. The seven notes went         TestTheDrawingsMatchTheSource`, twice — `-count=2`) and passes. The seven notes went
    ┃          with it: the gate now checks the **version it runs** against the pin (an old                 with it: the gate now checks the **version it runs** against the pin (an old
    ┃          `staticcheck` on the `PATH` fails it with *"run `make tools`"*), `docs/testing.md`'s         `staticcheck` on the `PATH` fails it with *"run `make tools`"*), `docs/testing.md`'s
    ┃          versions are the **checked copy** of the `Makefile`'s pins, `dep-budget.sh` honours          versions are the **checked copy** of the `Makefile`'s pins, `dep-budget.sh` honours
    ┃          `GO`, `tidy-check` puts `go.mod` and `go.sum` back when tidy fails, `install-hooks`          `GO`, `tidy-check` puts `go.mod` and `go.sum` back when tidy fails, `install-hooks`
    ┃          says what it replaced, `.NOTPARALLEL:` keeps `make -j` safe, and the ledger's                says what it replaced, `.NOTPARALLEL:` keeps `make -j` safe, and the ledger's
    ┃          discipline table names the make interface. | `Makefile` — `TEST_FLAGS`, `test test-          discipline table names the make interface. | `Makefile` — `TEST_FLAGS`, `test test-
    ┃          race`; `scripts/go-test.sh` — the zero-match refusal, now on both targets |                  race`; `scripts/go-test.sh` — the zero-match refusal, now on both targets |
    ┃                                                                                                 203 + | F-103 | medium | **The domain's one-answer rule for the plan's name is written
    ┃                                                                                                       twice, and only one copy is pinned.** `core.Store.LoadProject` answered an empty name
    ┃                                                                                                       with the project's directory name, and `core.Store.ProjectName` answered the same
    ┃                                                                                                       question again for itself — so the plan's name had two answers, and the test only
    ┃                                                                                                       held the second; `LoadProject` had **no test of its own**, and the reviewer's
    ┃                                                                                                       mutation left **the whole suite green** while `maitre-cli plan` printed `"Name": ""`
    ┃                                                                                                       and the console printed the project's directory name. *Closed in* `internal/core/
    ┃                                                                                                       store.go` — the rule now lives in **one** private function, `projectName`, and
    ┃                                                                                                       `LoadProject` and `Store.ProjectName` both **ask** it; the copy that answers is
    ┃                                                                                                       pinned by a test of its own. The reviewer's mutation — taking the rule out of the one
    ┃                                                                                                       copy — turns the suite **red** in both packages, in **three** tests as the runner
    ┃                                                                                                       claimed: `internal/core` (`TestProjectNameHasOneAnswer`,
    ┃                                                                                                       `TestLoadProjectAnswersTheSameName`) and `internal/ui` (`TestTheNameHasOneAnswer`).
    ┃                                                                                                       The two surfaces now agree in all three states — the plan names itself, the file
    ┃                                                                                                       names it **empty**, the file is **missing**: `maitre-cli plan` and `maitre --render`
    ┃                                                                                                       both print `demo2`, and the CLI's behaviour **changed** to match the console's, which
    ┃                                                                                                       is the point of the finding. The writer is untouched and `core.DomainVersion` is
    ┃                                                                                                       still **7**. | `internal/core/store.go` — `projectName`, `LoadProject`,
    ┃                                                                                                       `ProjectName`; `internal/core/store_test.go` — `TestProjectNameHasOneAnswer`,
    ┃                                                                                                       `TestLoadProjectAnswersTheSameName`; `internal/ui/identity_test.go` —
    ┃                                                                                                       `TestTheNameHasOneAnswer` |
    ┃                                                                                                 204 + | F-104 | low | **The Settings tab's `root` row renders the plan's path unsanitised**
    ┃                                                                                                       — the name went through `sanitize` and the path did not, so a plan whose directory
    ┃                                                                                                       name carried a newline split the row. *Closed in* `internal/ui/identity.go` — the
    ┃                                                                                                       path is sanitised **inside `humanPath`**, so every path the console draws goes
    ┃                                                                                                       through it **by construction** and not by a call site's care (`humanPath` is the only
    ┃                                                                                                       path the console renders). The proof by command: a plan under a directory named `we\
    ┃                                                                                                       nird` shows `root  …/we ird/kitchen` on **one** row; the reviewer's mutation — the
    ┃                                                                                                       `sanitize` call taken out of `humanPath` — fails `TestTheRootRowIsOneRow` twice over.
    ┃                                                                                                       | `internal/ui/identity.go` — `humanPath`; `internal/ui/identity_test.go` —
    ┃                                                                                                       `TestTheRootRowIsOneRow` |
    ┃                                                                                                 205 + | F-105 | low | **The state row's count and its list disagree** — the row said *12
    ┃                                                                                                       commits* and named nine, and the three it did not name were the three commits that
    ┃                                                                                                       wrote the row. *Closed in* `work/tui/README.md` — the shape is one that cannot drift:
    ┃                                                                                                       **the count is the number of ids the row names**, and the commits that wrote the row
    ┃                                                                                                       are not among them (a commit cannot carry its own hash). The row now says *11
    ┃                                                                                                       commits* and names **eleven real ids** — every one of them in the log — and the
    ┃                                                                                                       **protocol section** carries the shape in one line, so the next writer knows it. |
    ┃                                                                                                       `work/tui/README.md` — the T1 row and the protocol section |
    ┃    203   | F-99 | medium | **The name the console renders is not sanitised — the plan's own     206   | F-99 | medium | **The name the console renders is not sanitised — the plan's own
    ┃          `project.json` name reaches the header and the window title, and `maitre-cli init`           `project.json` name reaches the header and the window title, and `maitre-cli init`
    ┃          writes that name from the directory's basename.** A directory with a newline in it           writes that name from the directory's basename.** A directory with a newline in it
    ┃          broke the header into two rows. *Closed in* `internal/ui/identity.go` — the name goes        broke the header into two rows. *Closed in* `internal/ui/identity.go` — the name goes
    ┃          through `sanitize`, the same door as every other rendered string, and the **file on          through `sanitize`, the same door as every other rendered string, and the **file on
    ┃          disk keeps what was written**: the sanitising is a rendering concern. The reviewer's         disk keeps what was written**: the sanitising is a rendering concern. The reviewer's
    ┃          mutation — taking the `sanitize` call out of the name's path — fails the test three          mutation — taking the `sanitize` call out of the name's path — fails the test three
    ┃          times over (*the name still carries a control character*, *the header does not carry         times over (*the name still carries a control character*, *the header does not carry
    ┃          the sanitised name*, *the header is 41 rows, want 40*), and the proof by command: a          the sanitised name*, *the header is 41 rows, want 40*), and the proof by command: a
    ┃          plan under a directory named `we\nird` renders `⋈ maitre · we ird` on **one** row            plan under a directory named `we\nird` renders `⋈ maitre · we ird` on **one** row
    ┃          while `project.json` still holds `"we\nird"`. | `internal/ui/identity.go` —                  while `project.json` still holds `"we\nird"`. | `internal/ui/identity.go` —
    ┃          `projectName`; `internal/ui/identity_test.go` — `TestTheNameIsSanitised` |                   `projectName`; `internal/ui/identity_test.go` — `TestTheNameIsSanitised` |
    ┃    204   | F-100 | medium | **The plan's name has one answer — and it had three.** A plan that  207   | F-100 | medium | **The plan's name has one answer — and it had three.** A plan that
    ┃          names itself answered with its name, a plan without a `project.json` answered with the       names itself answered with its name, a plan without a `project.json` answered with
    ┃          project's directory name, and a plan whose `project.json` named it **empty** answered        the project's directory name, and a plan whose `project.json` named it **empty**
    ┃          with the **absolute path** of the plan directory — which is what landed in the window        answered with the **absolute path** of the plan directory — which is what landed in
    ┃          title. *Closed in* `internal/core/store.go` — `LoadProject` and the new `Store.              the window title. *Closed in* `internal/core/store.go` — `LoadProject` and the new
    ┃          ProjectName` answer the same for an empty name as for a missing file, and **a name is        `Store.ProjectName` answer the same for an empty name as for a missing file, and **a
    ┃          never a path**; `internal/ui/identity.go` — the console **asks** the domain (`m.store.       name is never a path**; `internal/ui/identity.go` — the console **asks** the domain (
    ┃          ProjectName()`) instead of answering the question a second time. The writer is               `m.store.ProjectName()`) instead of answering the question a second time. The writer
    ┃          untouched: `Init` writes the name it wrote before, and the plan format on disk does          is untouched: `Init` writes the name it wrote before, and the plan format on disk
    ┃          not move. The proof by command: the empty-name plan and the plan without a `project.         does not move. The proof by command: the empty-name plan and the plan without a
    ┃          json` both show `⋈ maitre · demo`, and `maitre-cli plan` prints the same name. **The         `project.json` both show `⋈ maitre · demo`, and `maitre-cli plan` prints the same
    ┃          rule is written in two places and only one copy is pinned — that is F-103.** |               name. **The rule is written in two places and only one copy is pinned — that is F-
    ┃          `internal/core/store.go` — `LoadProject`, `ProjectName`, `defaultProjectName`;               103.** | `internal/core/store.go` — `LoadProject`, `ProjectName`,
    ┃          `internal/core/store_test.go` — `TestProjectNameHasOneAnswer`; `internal/ui/                 `defaultProjectName`; `internal/core/store_test.go` — `TestProjectNameHasOneAnswer`;
    ┃          identity_test.go` — `TestTheNameHasOneAnswer` |                                              `internal/ui/identity_test.go` — `TestTheNameHasOneAnswer` |
    ┃    205   | F-101 | medium | **The row labelled `root` showed the plan directory's *parent*, so  208   | F-101 | medium | **The row labelled `root` showed the plan directory's *parent*, so
    ┃          `kitchen/` was nowhere in the console — and the label did not name what it showed.**         `kitchen/` was nowhere in the console — and the label did not name what it showed.**
    ┃          *Closed in* `internal/ui/views.go` — the row shows `humanPath(m.store.Root)`, the plan       *Closed in* `internal/ui/views.go` — the row shows `humanPath(m.store.Root)`, the
    ┃          directory (`~/Code/Pessoal/maitre/kitchen`), and the themes are still listed. This one       plan directory (`~/Code/Pessoal/maitre/kitchen`), and the themes are still listed.
    ┃          was **the brief's fault**: T1's brief gave the project's directory as the example, and       This one was **the brief's fault**: T1's brief gave the project's directory as the
    ┃          the runner and the first review both followed the example; the words win, and `work/         example, and the runner and the first review both followed the example; the words
    ┃          briefs/tui-t1.md` was corrected in the same commit. The reviewer's mutation — the row        win, and `work/briefs/tui-t1.md` was corrected in the same commit. The reviewer's
    ┃          put back on the parent — fails `TestTheSettingsTabShowsThePlanRoot` twice over. |            mutation — the row put back on the parent — fails
    ┃          `internal/ui/views.go` — `settingsLines`; `internal/ui/identity_test.go` —                   `TestTheSettingsTabShowsThePlanRoot` twice over. | `internal/ui/views.go` —
    ┃          `TestTheSettingsTabShowsThePlanRoot` |                                                       `settingsLines`; `internal/ui/identity_test.go` —
    ┃                                                                                                       `TestTheSettingsTabShowsThePlanRoot` |
    ┃    206   | F-102 | low | **`docs/05-human-interface.md` said the console *owns* the window      209   | F-102 | low | **`docs/05-human-interface.md` said the console *owns* the window
    ┃          title and nothing said what it carries.** *Closed in* `docs/05-human-interface.md` —         title and nothing said what it carries.** *Closed in* `docs/05-human-interface.md` —
    ┃          one paragraph: the title is the row the operating system puts on the terminal's tab,         one paragraph: the title is the row the operating system puts on the terminal's tab,
    ┃          it reads `maitre · <project>`, and it is **not** the screen's row 0 (that row is the         it reads `maitre · <project>`, and it is **not** the screen's row 0 (that row is the
    ┃          header). The screen sketch stays. **Nothing holds the paragraph** — the reviewer's           header). The screen sketch stays. **Nothing holds the paragraph** — the reviewer's
    ┃          probe, deleting it, leaves the doc check green. | — (a doc fix; nothing holds it) |          probe, deleting it, leaves the doc check green. | — (a doc fix; nothing holds it) |
    ┃
    ┃

       Now the notes table — the two closed notes come out, the new one goes in:

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    128   | The goldens' checks read **one directory** | `internal/ui/testdata` — a golden file  128   | The goldens' checks read **one directory** | `internal/ui/testdata` — a golden file
    ┃          elsewhere in the tree is neither claimed nor checked for the commit. |                       elsewhere in the tree is neither claimed nor checked for the commit. |
    ┃    129   | `go-test.sh` hands an **unknown flag to every binary** | The package goes first so   129   | `go-test.sh` hands an **unknown flag to every binary** | The package goes first so
    ┃          that the goldens' `-update` reaches the test binary — but a flag `go test` does not          that the goldens' `-update` reaches the test binary — but a flag `go test` does not
    ┃          know is now passed to every test binary in a run, not only to the goldens'. |                know is now passed to every test binary in a run, not only to the goldens'. |
    ┃    130   | A golden whose name carries a **space** | The scripts walk the goldens by word, so a 130   | A golden whose name carries a **space** | The scripts walk the goldens by word, so
    ┃          name with a space in it would be mishandled. There is none today. |                          a name with a space in it would be mishandled. There is none today. |
    ┃    131   | The state row's **commit list is short by its own commits** | T1's row says *7       131   | The state row's **commit list is short by its own commits** | T1's row says *7
    ┃          commits* and names five. The two it cannot name are the two commits that write the row       commits* and names five. The two it cannot name are the two commits that write the
    ┃          — a commit cannot carry its own hash — but the second of them could have named the           row — a commit cannot carry its own hash — but the second of them could have named
    ┃          first. The count is right; the list is not complete. The fix cycle repeated it: the          the first. The count is right; the list is not complete. The fix cycle repeated it:
    ┃          row says *12 commits* and names nine, and the three it does not name are the three           the row says *12 commits* and names nine, and the three it does not name are the
    ┃          commits that wrote the row. The shape that cannot drift is a **count that counts what        three commits that wrote the row. The shape that cannot drift is a **count that
    ┃          the list names** — or a row that says the count excludes the commits that wrote it. |        counts what the list names** — or a row that says the count excludes the commits that
    ┃                                                                                                       wrote it. |
    ┃    132 - | The Settings tab's **path row is not sanitised** | The name goes through `sanitize`
    ┃          (F-99); the path does not. A plan directory whose name carries a newline splits the
    ┃          row — `root  /tmp/vfy/we` / `ird/kitchen` — and pushes the list one visual row, so the
    ┃          pinned footer can slip a row. The header and the window title both survive: the name
    ┃          is the plan's content, and the path is the filesystem's. The row's unsanitised path
    ┃          **predates** the slice, so this is a note and not a blocker — but `sanitize(humanPath(
    ┃          ...))` would close it. |
    ┃    133   | `docs/05`'s window-title paragraph is **held by nothing** | The paragraph that says  132   | `docs/05`'s window-title paragraph is **held by nothing** | The paragraph that says
    ┃          what the window title carries (F-102) has no wire: deleting it leaves the doc check          what the window title carries (F-102) has no wire: deleting it leaves the doc check
    ┃          green. A doc fix with nothing holding it is a doc fix that can be undone by the next         green. A doc fix with nothing holding it is a doc fix that can be undone by the next
    ┃          reader. |                                                                                    reader. |
    ┃    134   | The plan's name rule is **nowhere in the spec** | `docs/03` names `project.json` as  133   | The plan's name rule is **nowhere in the spec** | `docs/03` names `project.json` as
    ┃          *project meta · charter · settings* and never says what the reader answers when the          *project meta · charter · settings* and never says what the reader answers when the
    ┃          name is empty — now the project's directory name. The rule is in the code and (from F-       name is empty — now the project's directory name. The rule is in the code and (from
    ┃          103) only half-pinned; a sentence beside the plan format would put it where a reader         F-103) only half-pinned; a sentence beside the plan format would put it where a
    ┃          looks. |                                                                                     reader looks. |
    ┃    135                                                                                          134
    ┃    136   ## Closed findings                                                                     135   ## Closed findings
    ┃
    ┃

       + Thought · 1.1s

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    127   | A failed claim leaves a **timestamp in the past** | An orphan golden keeps its 2000- 127   | A failed claim leaves a **timestamp in the past** | An orphan golden keeps its
    ┃          01-01 mtime on disk after a failed `goldens-check`. Git does not track it; the next          2000-01-01 mtime on disk after a failed `goldens-check`. Git does not track it; the
    ┃          `make goldens` rewrites the file. |                                                          next `make goldens` rewrites the file. |
    ┃    128   | The goldens' checks read **one directory** | `internal/ui/testdata` — a golden file  128   | The goldens' checks read **one directory** | `internal/ui/testdata` — a golden file
    ┃          elsewhere in the tree is neither claimed nor checked for the commit. |                       elsewhere in the tree is neither claimed nor checked for the commit. |
    ┃    129   | `go-test.sh` hands an **unknown flag to every binary** | The package goes first so   129   | `go-test.sh` hands an **unknown flag to every binary** | The package goes first so
    ┃          that the goldens' `-update` reaches the test binary — but a flag `go test` does not          that the goldens' `-update` reaches the test binary — but a flag `go test` does not
    ┃          know is now passed to every test binary in a run, not only to the goldens'. |                know is now passed to every test binary in a run, not only to the goldens'. |
    ┃    130   | A golden whose name carries a **space** | The scripts walk the goldens by word, so a 130   | A golden whose name carries a **space** | The scripts walk the goldens by word, so
    ┃          name with a space in it would be mishandled. There is none today. |                          a name with a space in it would be mishandled. There is none today. |
    ┃    131 - | The state row's **commit list is short by its own commits** | T1's row says *7
    ┃          commits* and names five. The two it cannot name are the two commits that write the row
    ┃          — a commit cannot carry its own hash — but the second of them could have named the
    ┃          first. The count is right; the list is not complete. The fix cycle repeated it: the
    ┃          row says *12 commits* and names nine, and the three it does not name are the three
    ┃          commits that wrote the row. The shape that cannot drift is a **count that counts what
    ┃          the list names** — or a row that says the count excludes the commits that wrote it. |
    ┃    132   | `docs/05`'s window-title paragraph is **held by nothing** | The paragraph that says  131   | `docs/05`'s window-title paragraph is **held by nothing** | The paragraph that says
    ┃          what the window title carries (F-102) has no wire: deleting it leaves the doc check          what the window title carries (F-102) has no wire: deleting it leaves the doc check
    ┃          green. A doc fix with nothing holding it is a doc fix that can be undone by the next         green. A doc fix with nothing holding it is a doc fix that can be undone by the next
    ┃          reader. |                                                                                    reader. |
    ┃    133   | The plan's name rule is **nowhere in the spec** | `docs/03` names `project.json` as  132   | The plan's name rule is **nowhere in the spec** | `docs/03` names `project.json` as
    ┃          *project meta · charter · settings* and never says what the reader answers when the          *project meta · charter · settings* and never says what the reader answers when the
    ┃          name is empty — now the project's directory name. The rule is in the code and (from F-       name is empty — now the project's directory name. The rule is in the code and (from
    ┃          103) only half-pinned; a sentence beside the plan format would put it where a reader         F-103) only half-pinned; a sentence beside the plan format would put it where a
    ┃          looks. |                                                                                     reader looks. |
    ┃                                                                                                 133 + | A test assertion that **cannot fail** | `TestTheRootRowIsOneRow` splits the rows on
    ┃                                                                                                       `\n` and then checks each line for a `\n` — a line that came out of a split cannot
    ┃                                                                                                       carry one. The test still bites through its other two assertions (*the root row does
    ┃                                                                                                       not carry the sanitised path*, *the root row carries a newline*), so the rule is
    ┃                                                                                                       held; the loop is decoration. |
    ┃    134                                                                                          134
    ┃    135   ## Closed findings                                                                     135   ## Closed findings
    ┃    136                                                                                          136
    ┃    137   | # | Severity | Closed in | Held by |                                                 137   | # | Severity | Closed in | Held by |
    ┃
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && awk '/^\|/ {n=gsub(/\|/,"|"); if ((NR>=63&&NR<=64)||(NR>=113&&NR<=116)||(NR>=127&&NR<=132)||(NR>=200&&NR<=210)) print NR": pipes="n" :: "substr(
    ┃  $0,1,55)}' docs/audit.md && echo "=== open findings table:" && sed -n '/^## Open findings/,/^### Notes/p' docs/audit.md | cut -c1-90
    ┃
    ┃  (27 earlier lines)
    ┃  | **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.resolv
    ┃  | **F-61** | medium | **The `courses/` directory is named by the plan format, created by `
    ┃  | F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and the s
    ┃  | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-scope=`.
    ┃
    ┃  ### Notes left for a later round
    ┃

       + Thought · 1.3s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'EOF'
    ┃  p='docs/audit.md'…
    ┃
    ┃  (2 earlier lines)
    ┃  203:| F-99 | medium | **The name the console renders is not
    ┃  204:| F-101 | medium | **The row labelled `root` showed the
    ┃  205:| F-102 | low | **`docs/05-human-interface.md` said the
    ┃  206:| F-103 | medium | **The domain's one-answer rule for th
    ┃  207:| F-104 | low | **The Settings tab's `root` row renders
    ┃  208:| F-105 | low | **The state row's count and its list dis
    ┃

       + Thought · 489ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && awk 'NR>=201 && NR<=209 {print NR": "substr($0,1,40)}' docs/audit.md
    ┃
    ┃  (2 earlier lines)
    ┃  203: | F-99 | medium | **The name the console
    ┃  204: | F-101 | medium | **The row labelled `r
    ┃  205: | F-102 | low | **`docs/05-human-interfa
    ┃  206: | F-103 | medium | **The domain's one-an
    ┃  207: | F-104 | low | **The Settings tab's `ro
    ┃  208: | F-105 | low | **The state row's count
    ┃  209:
    ┃

       + Thought · 2.2s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git checkout -- docs/audit.md && grep -n "^| F-10[0-5]\|^| F-99\|^\*\*Cycles 1–7" docs/audit.md | cut -c1-50
    ┃
    ┃  203:| F-99 | medium | **The name the console rende
    ┃  204:| F-100 | medium | **The plan's name has one a
    ┃  205:| F-101 | medium | **The row labelled `root` s
    ┃  206:| F-102 | low | **`docs/05-human-interface.md`
    ┃  208:**Cycles 1–7's own findings are not recoverabl
    ┃

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    60   | Slice H · cycle 5 | The last cycle of slice H: the gate sees a golden nobody          60   | Slice H · cycle 5 | The last cycle of slice H: the gate sees a golden nobody
    ┃         committed, and every message says what to do | The goldens are checked **three ways** —      committed, and every message says what to do | The goldens are checked **three ways**
    ┃         *unchanged*, *claimed by a test*, *in the commit*. The commit check is the new one:          — *unchanged*, *claimed by a test*, *in the commit*. The commit check is the new one:
    ┃         `git diff` does not see an untracked file, so the files are **intent-to-add** before         `git diff` does not see an untracked file, so the files are **intent-to-add** before
    ┃         the diff and a new golden shows up in it as a change; the index is put back as it was (      the diff and a new golden shows up in it as a change; the index is put back as it was
    ┃         `git rm --cached --force`, not `git reset`, so **a staged deletion stays staged**). The      (`git rm --cached --force`, not `git reset`, so **a staged deletion stays staged**).
    ┃         reviewer's probes: a new golden **no test claims** → the check fails, names the file         The reviewer's probes: a new golden **no test claims** → the check fails, names the
    ┃         and says which cause (*"no test source mentions it, so no test claims it"*) and its          file and says which cause (*"no test source mentions it, so no test claims it"*) and
    ┃         remedy (*"write a test that pins it, or delete the file"*); a new golden **a test            its remedy (*"write a test that pins it, or delete the file"*); a new golden **a test
    ┃         claims but nobody committed** → the check fails and names it (*"these golden files are       claims but nobody committed** → the check fails and names it (*"these golden files are
    ┃         not in the commit"*); a golden **staged for deletion** → the same failure, and the           not in the commit"*); a golden **staged for deletion** → the same failure, and the
    ┃         index afterwards is the index before (*"D "* plus the file on disk, untracked). The          index afterwards is the index before (*"D "* plus the file on disk, untracked). The
    ┃         messages were brought to the rule: each cause names **its own** remedy — *"name the          messages were brought to the rule: each cause names **its own** remedy — *"name the
    ┃         test so the filter matches it — its name must carry 'Golden' — or extend the filter"*,       test so the filter matches it — its name must carry 'Golden' — or extend the filter"*,
    ┃         and *"write a test that pins it, or delete the file"* — and the reviewer could follow        and *"write a test that pins it, or delete the file"* — and the reviewer could follow
    ┃         each one and be unblocked without reading the script. The doc check is scoped to the         each one and be unblocked without reading the script. The doc check is scoped to the
    ┃         **pin's shape**: the release version (`staticcheck 2026.2.2`) passes and a mention           **pin's shape**: the release version (`staticcheck 2026.2.2`) passes and a mention
    ┃         carrying the pin's shape with the wrong value (`v0.8.1`) fails. The cycle touched **no       carrying the pin's shape with the wrong value (`v0.8.1`) fails. The cycle touched **no
    ┃         Go source** — **no blockers filed** |                                                        Go source** — **no blockers filed** |
    ┃    61   | Slice T1 | The identity: the window title carries the project; the header's name      61   | Slice T1 | The identity: the window title carries the project; the header's name
    ┃         stops being a footnote; `--project` overrides the name for the session; the Settings         stops being a footnote; `--project` overrides the name for the session; the Settings
    ┃         tab shows the plan's path | The name has **one** source — `Model.projectName` in             tab shows the plan's path | The name has **one** source — `Model.projectName` in
    ┃         `internal/ui/identity.go` — and the header and the window title both read it; the title      `internal/ui/identity.go` — and the header and the window title both read it; the
    ┃         reads `maitre · <project>`, and a plan whose `project.json` has no name falls back to        title reads `maitre · <project>`, and a plan whose `project.json` has no name falls
    ┃         its root. `--project` reaches all three constructors (`New`, `NewLive`, `NewScripted`)       back to its root. `--project` reaches all three constructors (`New`, `NewLive`,
    ┃         and `--render` and `--keys` honour it, because they build the same model — and it            `NewScripted`) and `--render` and `--keys` honour it, because they build the same
    ┃         **writes nothing**: after a run against a copy of a plan, `project.json` and `activity.      model — and it **writes nothing**: after a run against a copy of a plan, `project.
    ┃         json` are **byte-for-byte unchanged**. The header's name wears the accent (`39`) where       json` and `activity.json` are **byte-for-byte unchanged**. The header's name wears the
    ┃         the faint (`240`) read as a footnote, and the mark, the product's name and the               accent (`39`) where the faint (`240`) read as a footnote, and the mark, the product's
    ┃         separator did not move. The Settings tab shows `~/Code/Pessoal/maitre` — the                 name and the separator did not move. The Settings tab shows `~/Code/Pessoal/maitre` —
    ┃         **project's** directory, not the plan directory inside it — and the themes are still         the **project's** directory, not the plan directory inside it — and the themes are
    ┃         listed. **The goldens did not move**: `make goldens` leaves `git diff` empty, because a      still listed. **The goldens did not move**: `make goldens` leaves `git diff` empty,
    ┃         style is not text and the Settings tab is not in the golden. The scope held: `internal/      because a style is not text and the Settings tab is not in the golden. The scope held:
    ┃         core` and `internal/cli` are untouched, `core.DomainVersion` is still **7**, the budget      `internal/core` and `internal/cli` are untouched, `core.DomainVersion` is still **7**,
    ┃         is still **0 / 0 / 28**, the seven tabs kept their names, the mark's glyph is still          the budget is still **0 / 0 / 28**, the seven tabs kept their names, the mark's glyph
    ┃         only in `internal/brand` — **no blockers filed** |                                           is still only in `internal/brand` — **no blockers filed** |
    ┃    62   | Slice T1 · cycle 1 | The review of the identity: four mutations, the write-nothing    62   | Slice T1 · cycle 1 | The review of the identity: four mutations, the write-nothing
    ┃         proof, and the gate | Each mutation fails a test and was reverted, and the tree was          proof, and the gate | Each mutation fails a test and was reverted, and the tree was
    ┃         left exactly as the runner left it: the hardcoded `"maitre"` fails **four** window-          left exactly as the runner left it: the hardcoded `"maitre"` fails **four** window-
    ┃         title tests; a `--project` that writes the name into `project.json` fails *"the              title tests; a `--project` that writes the name into `project.json` fails *"the
    ┃         override writes nothing"* and *"the plan keeps its own name"*; the name back in              override writes nothing"* and *"the plan keeps its own name"*; the name back in
    ┃         `styFaint` fails *"the header's name is not a footnote"*; a title that computes its own      `styFaint` fails *"the header's name is not a footnote"*; a title that computes its
    ┃         name fails *"the override reaches the console"*. The write-nothing claim was proved          own name fails *"the override reaches the console"*. The write-nothing claim was
    ┃         **by command** and not by the diff: a copy of a plan, `maitre --render --project …`          proved **by command** and not by the diff: a copy of a plan, `maitre --render --
    ┃         over it, and `cmp` on both `project.json` and `activity.json`. `make check` and `make        project …` over it, and `cmp` on both `project.json` and `activity.json`. `make check`
    ┃         check-race` are green, `make goldens` leaves the tree clean, and the state row is            and `make check-race` are green, `make goldens` leaves the tree clean, and the state
    ┃         **in** the slice's commits — **no blockers filed** |                                         row is **in** the slice's commits — **no blockers filed** |
    ┃    63   | Slice T1 · cycle 2 | The review of the fix cycle: the four findings, the two          63   | Slice T1 · cycle 2 | The review of the fix cycle: the four findings, the two
    ┃         observations, and the `docs/03` question | F-99..F-102 are closed, and each was proved       observations, and the `docs/03` question | F-99..F-102 are closed, and each was proved
    ┃         by a mutation: the `sanitize` call taken out of the name's path fails                        by a mutation: the `sanitize` call taken out of the name's path fails
    ┃         `TestTheNameIsSanitised` three times over; the `root` row put back on the parent fails       `TestTheNameIsSanitised` three times over; the `root` row put back on the parent fails
    ┃         `TestTheSettingsTabShowsThePlanRoot`; the empty-name branch taken out of `LoadProject`       `TestTheSettingsTabShowsThePlanRoot`; the empty-name branch taken out of `LoadProject`
    ┃         leaves **the whole suite green** — and that is **F-103**, the one blocker filed. The         leaves **the whole suite green** — and that is **F-103**, the one blocker filed. The
    ┃         name has one answer in the code, but the rule is written **twice** and only one copy is      name has one answer in the code, but the rule is written **twice** and only one copy
    ┃         pinned: `maitre-cli plan` then prints `"Name": ""` while the console shows the               is pinned: `maitre-cli plan` then prints `"Name": ""` while the console shows the
    ┃         project's directory name. The **path row** is judged a note and not a blocker: the path      project's directory name. The **path row** is judged a note and not a blocker: the
    ┃         is not the plan's content, the unsanitised path **predates** the slice, and the              path is not the plan's content, the unsanitised path **predates** the slice, and the
    ┃         consequence is one extra visual row in a settings list — the header and the window           consequence is one extra visual row in a settings list — the header and the window
    ┃         title both survive. The **state row** is judged a note as well: the count should count       title both survive. The **state row** is judged a note as well: the count should count
    ┃         what the list names, or the row should say the count excludes the commits that wrote         what the list names, or the row should say the count excludes the commits that wrote
    ┃         it. The plan's name does **not** belong in `docs/03` § 8: that table is the derived          it. The plan's name does **not** belong in `docs/03` § 8: that table is the derived
    ┃         views — *computed, never stored* — and the name is a **stored** fact the domain reads,       views — *computed, never stored* — and the name is a **stored** fact the domain reads,
    ┃         like `LoadProject` and `ListRecipes`, which the table does not carry either. `make           like `LoadProject` and `ListRecipes`, which the table does not carry either. `make
    ┃         check` and `make check-race` are green, the goldens did not move, and the tree was left      check` and `make check-race` are green, the goldens did not move, and the tree was
    ┃         exactly as the runner left it — **one blocker filed** (F-103) |                              left exactly as the runner left it — **one blocker filed** (F-103) |
    ┃                                                                                                 64 + | Slice T1 · cycle 3 | The **last** cycle of T1: the blocker, the two notes, and the
    ┃                                                                                                      question that is the chef's | F-103..F-105 are closed, each proved by a mutation and
    ┃                                                                                                      each reverted: the rule taken out of the one copy turns the suite **red** in both
    ┃                                                                                                      `internal/core` and `internal/ui` — **three** tests, as the runner claimed; the
    ┃                                                                                                      `sanitize` call taken out of `humanPath` fails the root-row test twice over; the state
    ┃                                                                                                      row's count is the number of ids the row names, and the **protocol section** carries
    ┃                                                                                                      the shape in one line. The two surfaces agree in **all three states** — the plan names
    ┃                                                                                                      itself, the file names it empty, the file is missing — and the CLI's behaviour
    ┃                                                                                                      **changed** to match the console's, which is the point of F-103. The rule is written
    ┃                                                                                                      down in the **code** (`core.Store.projectName`'s doc comment) and **nowhere in the
    ┃                                                                                                      spec**: `docs/03` § 6 names `project.json` as *project meta · charter · settings* and
    ┃                                                                                                      says nothing about the name's fallback, and `docs/05` does not either — where it
    ┃                                                                                                      belongs is the **chef's** call, and **nothing was written into `docs/03`**. The
    ┃                                                                                                      goldens did not move, the budget is still **0 / 0 / 28**, `core.DomainVersion` is
    ┃                                                                                                      still **7**, the seven tabs kept their names, the CLI's surface did not change, and
    ┃                                                                                                      the tree was left exactly as the runner left it — **no blockers filed** |
    ┃    64   | The rounds inside a cycle | Fix → verify → back to the executor → repeat, until the   65   | The rounds inside a cycle | Fix → verify → back to the executor → repeat, until the
    ┃         verifier is satisfied | Three rounds in cycle 8: the first fix introduced three new          verifier is satisfied | Three rounds in cycle 8: the first fix introduced three new
    ┃         mediums, the second closed them |                                                            mediums, the second closed them |
    ┃    65                                                                                           66
    ┃    66   ### The discipline every round runs                                                     67   ### The discipline every round runs
    ┃    67                                                                                           68
    ┃

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    112   |---|---|---|---|                                                                      112   |---|---|---|---|
    ┃    113   | **F-59** | medium | **`core.Audit.Decisions` has a shape and no writer** — it is     113   | **F-59** | medium | **`core.Audit.Decisions` has a shape and no writer** — it is
    ┃          where a decision *about the plan* would live, and nothing writes it. `kitchen/               where a decision *about the plan* would live, and nothing writes it. `kitchen/
    ┃          decisions.json` is named by the plan format and opened by nothing. The remedy is on          decisions.json` is named by the plan format and opened by nothing. The remedy is on
    ┃          the table in [06-proposals.md](06-proposals.md) § 2. | `internal/core/model.go:270`;         the table in [06-proposals.md](06-proposals.md) § 2. | `internal/core/model.go:270`;
    ┃          `internal/core/store.go:24` |                                                                `internal/core/store.go:24` |
    ┃    114   | **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.   114   | **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.
    ┃          resolved` gate can never fire** — nothing ever sets a question, so the list is always        resolved` gate can never fire** — nothing ever sets a question, so the list is always
    ┃          empty and the gate always passes. A gate nobody can fail is not a guard. | `internal/        empty and the gate always passes. A gate nobody can fail is not a guard. | `internal/
    ┃          core/model.go:31`; `internal/core/lifecycle.go:197` |                                        core/model.go:31`; `internal/core/lifecycle.go:197` |
    ┃    115   | **F-61** | medium | **The `courses/` directory is named by the plan format, created  115   | **F-61** | medium | **The `courses/` directory is named by the plan format, created
    ┃          by `Init`, and never written** — a directory no writer touches. Per [06-proposals.md](       by `Init`, and never written** — a directory no writer touches. Per [06-proposals.md]
    ┃          06-proposals.md) § 1 the honest answer is to stop naming it: a course is the relation,       (06-proposals.md) § 1 the honest answer is to stop naming it: a course is the
    ┃          not a file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.md) |             relation, not a file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.
    ┃                                                                                                       md) |
    ┃    116 - | **F-103** | medium | **The domain's one-answer rule for the plan's name is written
    ┃          twice, and only one copy is pinned.** `core.Store.LoadProject` answers an empty name
    ┃          with the project's directory name, and `core.Store.ProjectName` answers the same
    ┃          question again for itself — so the plan's name has two answers, and the test only
    ┃          holds the second. The reviewer's mutation — taking the empty-name branch out of
    ┃          `LoadProject` — leaves **the whole suite green** (every package), and then the two
    ┃          shipped surfaces disagree: `maitre-cli plan` prints `"Name": ""` while the console
    ┃          shows the project's directory name. `LoadProject` has **no test of its own**: the only
    ┃          test that reads it (`TestThePlanKeepsItsOwnName`) reads a plan that names itself. The
    ┃          remedy: pin `LoadProject`'s answer in `internal/core/store_test.go` for the same two
    ┃          cases `ProjectName` is pinned for, or drop `ProjectName`'s redundant fallback so the
    ┃          regression shows through it. | `internal/core/store.go` — `LoadProject`,
    ┃          `ProjectName`; `internal/core/store_test.go` — `TestProjectNameHasOneAnswer` |
    ┃    117   | F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and   116   | F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and
    ┃          the size.** It is path + mtime + size per file, so two different contents can share          the size.** It is path + mtime + size per file, so two different contents can share
    ┃          one fingerprint. The tool's own writes always change the mtime, so in practice only a        one fingerprint. The tool's own writes always change the mtime, so in practice only a
    ┃          hand edit that restores the mtime is missed — and the chef's `r` reloads it. |               hand edit that restores the mtime is missed — and the chef's `r` reloads it. |
    ┃          `internal/core/store.go` — `DiskFingerprint`. |                                              `internal/core/store.go` — `DiskFingerprint`. |
    ┃    118   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-      117   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-
    ┃          scope=`. | Recorded so nobody calls it a regression later. Not a defect. |                   scope=`. | Recorded so nobody calls it a regression later. Not a defect. |
    ┃    119                                                                                          118
    ┃    120   ### Notes left for a later round                                                       119   ### Notes left for a later round
    ┃
    ┃
    ┃

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    127   | A failed claim leaves a **timestamp in the past** | An orphan golden keeps its 2000- 127   | A failed claim leaves a **timestamp in the past** | An orphan golden keeps its
    ┃          01-01 mtime on disk after a failed `goldens-check`. Git does not track it; the next          2000-01-01 mtime on disk after a failed `goldens-check`. Git does not track it; the
    ┃          `make goldens` rewrites the file. |                                                          next `make goldens` rewrites the file. |
    ┃    128   | The goldens' checks read **one directory** | `internal/ui/testdata` — a golden file  128   | The goldens' checks read **one directory** | `internal/ui/testdata` — a golden file
    ┃          elsewhere in the tree is neither claimed nor checked for the commit. |                       elsewhere in the tree is neither claimed nor checked for the commit. |
    ┃    129   | `go-test.sh` hands an **unknown flag to every binary** | The package goes first so   129   | `go-test.sh` hands an **unknown flag to every binary** | The package goes first so
    ┃          that the goldens' `-update` reaches the test binary — but a flag `go test` does not          that the goldens' `-update` reaches the test binary — but a flag `go test` does not
    ┃          know is now passed to every test binary in a run, not only to the goldens'. |                know is now passed to every test binary in a run, not only to the goldens'. |
    ┃    130   | A golden whose name carries a **space** | The scripts walk the goldens by word, so a 130   | A golden whose name carries a **space** | The scripts walk the goldens by word, so
    ┃          name with a space in it would be mishandled. There is none today. |                          a name with a space in it would be mishandled. There is none today. |
    ┃    131 - | The state row's **commit list is short by its own commits** | T1's row says *7
    ┃          commits* and names five. The two it cannot name are the two commits that write the row
    ┃          — a commit cannot carry its own hash — but the second of them could have named the
    ┃          first. The count is right; the list is not complete. The fix cycle repeated it: the
    ┃          row says *12 commits* and names nine, and the three it does not name are the three
    ┃          commits that wrote the row. The shape that cannot drift is a **count that counts what
    ┃          the list names** — or a row that says the count excludes the commits that wrote it. |
    ┃    132 - | The Settings tab's **path row is not sanitised** | The name goes through `sanitize`
    ┃          (F-99); the path does not. A plan directory whose name carries a newline splits the
    ┃          row — `root  /tmp/vfy/we` / `ird/kitchen` — and pushes the list one visual row, so the
    ┃          pinned footer can slip a row. The header and the window title both survive: the name
    ┃          is the plan's content, and the path is the filesystem's. The row's unsanitised path
    ┃          **predates** the slice, so this is a note and not a blocker — but `sanitize(humanPath(
    ┃          ...))` would close it. |
    ┃    133   | `docs/05`'s window-title paragraph is **held by nothing** | The paragraph that says  131   | `docs/05`'s window-title paragraph is **held by nothing** | The paragraph that says
    ┃          what the window title carries (F-102) has no wire: deleting it leaves the doc check          what the window title carries (F-102) has no wire: deleting it leaves the doc check
    ┃          green. A doc fix with nothing holding it is a doc fix that can be undone by the next         green. A doc fix with nothing holding it is a doc fix that can be undone by the next
    ┃          reader. |                                                                                    reader. |
    ┃    134   | The plan's name rule is **nowhere in the spec** | `docs/03` names `project.json` as  132   | The plan's name rule is **nowhere in the spec** | `docs/03` names `project.json` as
    ┃          *project meta · charter · settings* and never says what the reader answers when the          *project meta · charter · settings* and never says what the reader answers when the
    ┃          name is empty — now the project's directory name. The rule is in the code and (from F-       name is empty — now the project's directory name. The rule is in the code and (from
    ┃          103) only half-pinned; a sentence beside the plan format would put it where a reader         F-103) only half-pinned; a sentence beside the plan format would put it where a
    ┃          looks. |                                                                                     reader looks. |
    ┃                                                                                                 133 + | A test assertion that **cannot fail** | `TestTheRootRowIsOneRow` splits the rows on
    ┃                                                                                                       `\n` and then checks each line for a `\n` — a line that came out of a split cannot
    ┃                                                                                                       carry one. The test still bites through its other two assertions (*the root row does
    ┃                                                                                                       not carry the sanitised path*, *the root row carries a newline*), so the rule is
    ┃                                                                                                       held; the loop is decoration. |
    ┃    135                                                                                          134
    ┃    136   ## Closed findings                                                                     135   ## Closed findings
    ┃    137                                                                                          136
    ┃    138   | # | Severity | Closed in | Held by |                                                 137   | # | Severity | Closed in | Held by |
    ┃
    ┃

       Now the three closed findings, appended after F-102:

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    202   | F-99 | medium | **The name the console renders is not sanitised — the plan's own     202   | F-99 | medium | **The name the console renders is not sanitised — the plan's own
    ┃          `project.json` name reaches the header and the window title, and `maitre-cli init`           `project.json` name reaches the header and the window title, and `maitre-cli init`
    ┃          writes that name from the directory's basename.** A directory with a newline in it           writes that name from the directory's basename.** A directory with a newline in it
    ┃          broke the header into two rows. *Closed in* `internal/ui/identity.go` — the name goes        broke the header into two rows. *Closed in* `internal/ui/identity.go` — the name goes
    ┃          through `sanitize`, the same door as every other rendered string, and the **file on          through `sanitize`, the same door as every other rendered string, and the **file on
    ┃          disk keeps what was written**: the sanitising is a rendering concern. The reviewer's         disk keeps what was written**: the sanitising is a rendering concern. The reviewer's
    ┃          mutation — taking the `sanitize` call out of the name's path — fails the test three          mutation — taking the `sanitize` call out of the name's path — fails the test three
    ┃          times over (*the name still carries a control character*, *the header does not carry         times over (*the name still carries a control character*, *the header does not carry
    ┃          the sanitised name*, *the header is 41 rows, want 40*), and the proof by command: a          the sanitised name*, *the header is 41 rows, want 40*), and the proof by command: a
    ┃          plan under a directory named `we\nird` renders `⋈ maitre · we ird` on **one** row            plan under a directory named `we\nird` renders `⋈ maitre · we ird` on **one** row
    ┃          while `project.json` still holds `"we\nird"`. | `internal/ui/identity.go` —                  while `project.json` still holds `"we\nird"`. | `internal/ui/identity.go` —
    ┃          `projectName`; `internal/ui/identity_test.go` — `TestTheNameIsSanitised` |                   `projectName`; `internal/ui/identity_test.go` — `TestTheNameIsSanitised` |
    ┃    203   | F-100 | medium | **The plan's name has one answer — and it had three.** A plan that  203   | F-100 | medium | **The plan's name has one answer — and it had three.** A plan that
    ┃          names itself answered with its name, a plan without a `project.json` answered with the       names itself answered with its name, a plan without a `project.json` answered with
    ┃          project's directory name, and a plan whose `project.json` named it **empty** answered        the project's directory name, and a plan whose `project.json` named it **empty**
    ┃          with the **absolute path** of the plan directory — which is what landed in the window        answered with the **absolute path** of the plan directory — which is what landed in
    ┃          title. *Closed in* `internal/core/store.go` — `LoadProject` and the new `Store.              the window title. *Closed in* `internal/core/store.go` — `LoadProject` and the new
    ┃          ProjectName` answer the same for an empty name as for a missing file, and **a name is        `Store.ProjectName` answer the same for an empty name as for a missing file, and **a
    ┃          never a path**; `internal/ui/identity.go` — the console **asks** the domain (`m.store.       name is never a path**; `internal/ui/identity.go` — the console **asks** the domain (
    ┃          ProjectName()`) instead of answering the question a second time. The writer is               `m.store.ProjectName()`) instead of answering the question a second time. The writer
    ┃          untouched: `Init` writes the name it wrote before, and the plan format on disk does          is untouched: `Init` writes the name it wrote before, and the plan format on disk
    ┃          not move. The proof by command: the empty-name plan and the plan without a `project.         does not move. The proof by command: the empty-name plan and the plan without a
    ┃          json` both show `⋈ maitre · demo`, and `maitre-cli plan` prints the same name. **The         `project.json` both show `⋈ maitre · demo`, and `maitre-cli plan` prints the same
    ┃          rule is written in two places and only one copy is pinned — that is F-103.** |               name. **The rule is written in two places and only one copy is pinned — that is F-
    ┃          `internal/core/store.go` — `LoadProject`, `ProjectName`, `defaultProjectName`;               103.** | `internal/core/store.go` — `LoadProject`, `ProjectName`,
    ┃          `internal/core/store_test.go` — `TestProjectNameHasOneAnswer`; `internal/ui/                 `defaultProjectName`; `internal/core/store_test.go` — `TestProjectNameHasOneAnswer`;
    ┃          identity_test.go` — `TestTheNameHasOneAnswer` |                                              `internal/ui/identity_test.go` — `TestTheNameHasOneAnswer` |
    ┃    204   | F-101 | medium | **The row labelled `root` showed the plan directory's *parent*, so  204   | F-101 | medium | **The row labelled `root` showed the plan directory's *parent*, so
    ┃          `kitchen/` was nowhere in the console — and the label did not name what it showed.**         `kitchen/` was nowhere in the console — and the label did not name what it showed.**
    ┃          *Closed in* `internal/ui/views.go` — the row shows `humanPath(m.store.Root)`, the plan       *Closed in* `internal/ui/views.go` — the row shows `humanPath(m.store.Root)`, the
    ┃          directory (`~/Code/Pessoal/maitre/kitchen`), and the themes are still listed. This one       plan directory (`~/Code/Pessoal/maitre/kitchen`), and the themes are still listed.
    ┃          was **the brief's fault**: T1's brief gave the project's directory as the example, and       This one was **the brief's fault**: T1's brief gave the project's directory as the
    ┃          the runner and the first review both followed the example; the words win, and `work/         example, and the runner and the first review both followed the example; the words
    ┃          briefs/tui-t1.md` was corrected in the same commit. The reviewer's mutation — the row        win, and `work/briefs/tui-t1.md` was corrected in the same commit. The reviewer's
    ┃          put back on the parent — fails `TestTheSettingsTabShowsThePlanRoot` twice over. |            mutation — the row put back on the parent — fails
    ┃          `internal/ui/views.go` — `settingsLines`; `internal/ui/identity_test.go` —                   `TestTheSettingsTabShowsThePlanRoot` twice over. | `internal/ui/views.go` —
    ┃          `TestTheSettingsTabShowsThePlanRoot` |                                                       `settingsLines`; `internal/ui/identity_test.go` —
    ┃                                                                                                       `TestTheSettingsTabShowsThePlanRoot` |
    ┃    205   | F-102 | low | **`docs/05-human-interface.md` said the console *owns* the window      205   | F-102 | low | **`docs/05-human-interface.md` said the console *owns* the window
    ┃          title and nothing said what it carries.** *Closed in* `docs/05-human-interface.md` —         title and nothing said what it carries.** *Closed in* `docs/05-human-interface.md` —
    ┃          one paragraph: the title is the row the operating system puts on the terminal's tab,         one paragraph: the title is the row the operating system puts on the terminal's tab,
    ┃          it reads `maitre · <project>`, and it is **not** the screen's row 0 (that row is the         it reads `maitre · <project>`, and it is **not** the screen's row 0 (that row is the
    ┃          header). The screen sketch stays. **Nothing holds the paragraph** — the reviewer's           header). The screen sketch stays. **Nothing holds the paragraph** — the reviewer's
    ┃          probe, deleting it, leaves the doc check green. | — (a doc fix; nothing holds it) |          probe, deleting it, leaves the doc check green. | — (a doc fix; nothing holds it) |
    ┃                                                                                                 206 + | F-103 | medium | **The domain's one-answer rule for the plan's name is written
    ┃                                                                                                       twice, and only one copy is pinned.** `core.Store.LoadProject` answered an empty name
    ┃                                                                                                       with the project's directory name, and `core.Store.ProjectName` answered the same
    ┃                                                                                                       question again for itself — so the plan's name had two answers, and the test only
    ┃                                                                                                       held the second; `LoadProject` had **no test of its own**, and the reviewer's
    ┃                                                                                                       mutation left **the whole suite green** while `maitre-cli plan` printed `"Name": ""`
    ┃                                                                                                       and the console printed the project's directory name. *Closed in* `internal/core/
    ┃                                                                                                       store.go` — the rule now lives in **one** private function, `projectName`, and
    ┃                                                                                                       `LoadProject` and `Store.ProjectName` both **ask** it; the copy that answers is
    ┃                                                                                                       pinned by a test of its own. The reviewer's mutation — taking the rule out of the one
    ┃                                                                                                       copy — turns the suite **red** in both packages, in **three** tests as the runner
    ┃                                                                                                       claimed: `internal/core` (`TestProjectNameHasOneAnswer`,
    ┃                                                                                                       `TestLoadProjectAnswersTheSameName`) and `internal/ui` (`TestTheNameHasOneAnswer`).
    ┃                                                                                                       The two surfaces now agree in all three states — the plan names itself, the file
    ┃                                                                                                       names it **empty**, the file is **missing**: `maitre-cli plan` and `maitre --render`
    ┃                                                                                                       both print `demo2`, and the CLI's behaviour **changed** to match the console's, which
    ┃                                                                                                       is the point of the finding. The writer is untouched and `core.DomainVersion` is
    ┃                                                                                                       still **7**. | `internal/core/store.go` — `projectName`, `LoadProject`,
    ┃                                                                                                       `ProjectName`; `internal/core/store_test.go` — `TestProjectNameHasOneAnswer`,
    ┃                                                                                                       `TestLoadProjectAnswersTheSameName`; `internal/ui/identity_test.go` —
    ┃                                                                                                       `TestTheNameHasOneAnswer` |
    ┃                                                                                                 207 + | F-104 | low | **The Settings tab's `root` row renders the plan's path unsanitised**
    ┃                                                                                                       — the name went through `sanitize` and the path did not, so a plan whose directory
    ┃                                                                                                       name carried a newline split the row. *Closed in* `internal/ui/identity.go` — the
    ┃                                                                                                       path is sanitised **inside `humanPath`**, so every path the console draws goes
    ┃                                                                                                       through it **by construction** and not by a call site's care (`humanPath` is the only
    ┃                                                                                                       path the console renders). The proof by command: a plan under a directory named `we\
    ┃                                                                                                       nird` shows `root  …/we ird/kitchen` on **one** row; the reviewer's mutation — the
    ┃                                                                                                       `sanitize` call taken out of `humanPath` — fails `TestTheRootRowIsOneRow` twice over.
    ┃                                                                                                       | `internal/ui/identity.go` — `humanPath`; `internal/ui/identity_test.go` —
    ┃                                                                                                       `TestTheRootRowIsOneRow` |
    ┃                                                                                                 208 + | F-105 | low | **The state row's count and its list disagree** — the row said *12
    ┃                                                                                                       commits* and named nine, and the three it did not name were the three commits that
    ┃                                                                                                       wrote the row. *Closed in* `work/tui/README.md` — the shape is one that cannot drift:
    ┃                                                                                                       **the count is the number of ids the row names**, and the commits that wrote the row
    ┃                                                                                                       are not among them (a commit cannot carry its own hash). The row now says *11
    ┃                                                                                                       commits* and names **eleven real ids** — every one of them in the log — and the
    ┃                                                                                                       **protocol section** carries the shape in one line, so the next writer knows it. |
    ┃                                                                                                       `work/tui/README.md` — the T1 row and the protocol section |
    ┃    206                                                                                          209
    ┃    207   **Cycles 1–7's own findings are not recoverable.** Their per-finding reports lived in  210   **Cycles 1–7's own findings are not recoverable.** Their per-finding reports lived in
    ┃          a                                                                                            a
    ┃    208   scratch directory that was wiped mid-cycle. The *fixes* are all in the repository, and 211   scratch directory that was wiped mid-cycle. The *fixes* are all in the repository,
    ┃                                                                                                       and
    ┃    209   cycle 8's report re-verified fifteen of them with evidence — those are the rows above. 212   cycle 8's report re-verified fifteen of them with evidence — those are the rows
    ┃                                                                                                       above.
    ┃
    ┃
    ┃
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "^| F-" docs/audit.md | tail -12 | cut -c1-45 && echo "=== open findings:" && sed -n '/^## Open findings/,/^### Notes/p' docs/audit.md |
    ┃   grep -c "^| \*\*F-\|^| F-" && sed -n '/^## Open findings/,/^### Notes/p' docs/audit.md | cut -c1-45 && echo "=== notes:" && sed -n '/^### Notes left/,/^## Closed/p' docs/audit.md…
    ┃
    ┃  (30 earlier lines)
    ┃  | The goldens' checks read **one directory** | `internal/ui/
    ┃  | `go-test.sh` hands an **unknown flag to every binary** | T
    ┃  | A golden whose name carries a **space** | The scripts walk
    ┃  | `docs/05`'s window-title paragraph is **held by nothing**
    ┃  | The plan's name rule is **nowhere in the spec** | `docs/03
    ┃  | A test assertion that **cannot fail** | `TestTheRootRowIsO
    ┃

       + Thought · 1.3s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff --stat && git diff docs/audit.md | grep "^+" | grep -c "" && git add docs/audit.md && git commit -q -m "docs: T1's last cycle in the le
    ┃  dger - F-103..F-105 closed, and the rule's place in the spec is the chef's" && git log --oneline -1 && git status --short
    ┃
    ┃  docs/audit.md | 8 +++++---
    ┃   1 file changed, 5 insertions(+), 3 deletions(-)
    ┃  6
    ┃  ecc5368 docs: T1's last cycle in the ledger - F-103..F-105 closed, and the rule's place in the spec is the chef's
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && mkdir -p work/inbox && cat > work/inbox/reviewer-tui-t1-cycle3.md <<EOF
    ┃  Reviewer: T1 cycle 3 done. Verdict: ok. Filed 0 blockers (none).…
    ┃
    ┃  05c6aff inbox: reviewer reports tui-t1 cycle 3
    ┃  ecc5368 docs: T1's last cycle in the ledger - F-103..F-105 closed, and the rule's place in the spec is the chef's
    ┃  0e9ed5c work: T1's cycle 3 review brief - the last cycle
    ┃  TREE CLEAN
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1H:p1" "Reviewer: T1 cycle 3 done. Verdict: ok. Filed 0 blockers (none). F-103 verified: the rule now lives in one private
    ┃  function - core.Store.projectName - and both LoadProject and ProjectName ask it (defaultProjectName has one caller). Mutation - I took the rule out of that one copy and ran the FU…
    ┃
    ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e737d6341ffeTdoEYTgykU25M8"},
    ┃  "agent_status":"done","completion_seq":84,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1H:p1","revision":15,
    ┃  "screen_detection_skipped":true,"state_change_seq":84,"tab_id":"w1H:t1","terminal_id":"term_65f0cd2aa89d26","terminal_title":"OC | Project work proposal","terminal_title_stripped":"OC |
    ┃  Project work proposal","workspace_id":"w1H"},"type":"agent_prompted"}}
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && git diff && git log --oneline -3 && echo "TREE CLEAN — only the review's two commits"
    ┃
    ┃  05c6aff inbox: reviewer reports tui-t1 cycle 3
    ┃  ecc5368 docs: T1's last cycle in the ledger - F-103..F-105 closed, and the rule's place in the spec is the chef's
    ┃  0e9ed5c work: T1's cycle 3 review brief - the last cycle
    ┃  TREE CLEAN — only the review's two commits
    ┃

       verdict: ok — 0 blockers filed.

       F-103..F-105 — all three closed, each proved by a mutation I reproduced

       ┌───────┬───────────────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────┐
       │ #     │ How it was proved                                                                 │ The mutation                                                                               │
       ├───────┼───────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-103 │ The rule lives in one private function, core.Store.projectName; LoadProject and   │ Taking the rule out of that one copy turns the full suite red in both packages, in three   │
       │       │ ProjectName both ask it, and defaultProjectName has a single caller               │ tests — internal/core: TestProjectNameHasOneAnswer, TestLoadProjectAnswersTheSameName;     │
       │       │                                                                                   │ internal/ui: TestTheNameHasOneAnswer — exactly as the runner claimed                       │
       ├───────┼───────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-104 │ humanPath sanitises, and it is the only path the console renders (one call site)  │ sanitize taken out of humanPath → TestTheRootRowIsOneRow fails twice                       │
       │       │ — by construction, not by a call site's care. A plan under we\nird shows root  …/ │                                                                                            │
       │       │ we ird/kitchen on one row                                                         │                                                                                            │
       ├───────┼───────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-105 │ The row says 11 commits and names 11 ids — they agree; every id is a real commit  │ —                                                                                          │
       │       │ of the slice, and the four row-writing commits are not among them. The protocol   │                                                                                            │
       │       │ section carries the shape in one line                                             │                                                                                            │
       └───────┴───────────────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────┘

       The two surfaces agree in all three states — the plan names itself · the file names it empty · the file is missing: maitre-cli plan and maitre --render both show demo2. The CLI's
       behaviour changed (it used to print ""), which was the point of F-103, and the agreement is not a fixture coincidence — the mutation turns both packages red.

       The rule in the spec

       Not written down anywhere in the spec. It is written down in the code — the doc comment on core.Store.projectName. docs/03 § 6 names project.json as project meta · charter · settings and
       says nothing about the name's fallback; docs/05 does not either. I wrote nothing into docs/03 — that placement is the chef's call, and the orchestrator has it.

       Scope and gate

       The goldens are unchanged (make goldens → empty diff); docs/03 untouched this cycle; docs/05's paragraph stays; internal/cli/cmd/ untouched; budget 0 / 0 / 28; seven tabs, same names;
       core.DomainVersion still 7; the glyph only in internal/brand; docs/audit.md untouched by the runner. make check and make check-race green. Every mutation reverted — git status and git
       diff are empty apart from my two commits (ecc5368 ledger, 05c6aff receipt).

       LOW notes (ledger's notes table): TestTheRootRowIsOneRow has one assertion that cannot fail (it checks each line of a \n-split for a \n — the test still bites through its other two
       assertions, so the rule is held); still open from earlier rounds: docs/05's window-title paragraph is held by nothing, and the plan's name rule is nowhere in the spec.

       The nudge was sent to w1H:p1; the receipt is in work/inbox/reviewer-tui-t1-cycle3.md.

       Build · MiniMax-M3.1-Flash · 10m 12s · 103.4 tok/s

    ┃
    ┃
    ┃
    ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
    ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
    ~/Code/Pessoal/maitre:master                                                                                                                                    205.3K (21%)  ctrl+p commands
```
