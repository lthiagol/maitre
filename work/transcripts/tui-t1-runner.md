# The runner's lane — T1, captured 2026-10-30

Dumped from pane `w1G:p3` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that wrote T1: the window title carries the project, the header's name stops
being a footnote, `--project`, and the Settings tab's path.
The review that judged it is in [tui-t1-reviewer.md](tui-t1-reviewer.md).

    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/05-human-interface.md work/briefs/tui-t1.md && git commit -q -m "docs: the window title carries the project, and the brief's example
    ┃   was wrong…
    ┃
    ┃  55bdc81 docs: the window title carries the project, and the brief's example was wrong
    ┃

       + Thought · 2.4s

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃    11   ## Where it stands                                                                      11   ## Where it stands
    ┃    12                                                                                           12
    ┃    13   | Session | Feature | What it lands | State | Brief |                                   13   | Session | Feature | What it lands | State | Brief |
    ┃    14   |--:|---|---|---|---|                                                                   14   |--:|---|---|---|---|
    ┃    15 - | **T1** | the identity | the window title carries the project; the header's name stops 15 + | **T1** | the identity | the window title carries the project; the header's name stops
    ┃         being a footnote; `--project`; the Settings tab shows the plan's path | 🔄 **review ok       being a footnote; `--project`; the Settings tab shows the plan's path | ✅ closed · 12
    ┃         · fix cycle 1** — 7 commits · 94a53ff, 8fc7438, f0b2aa5, 82289c6, 0bbcb35 | [tui-t1.md]      commits · 94a53ff, 8fc7438, f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052,
    ┃         (../briefs/tui-t1.md) · [fix 1](../briefs/fix-tui-t1-1.md) |                                 55bdc81 | [tui-t1.md](../briefs/tui-t1.md) · [fix 1](../briefs/fix-tui-t1-1.md) |
    ┃    16   | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the     16   | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the
    ┃         wheel scrolls | 🔄 **briefed** | [tui-t2a.md](../briefs/tui-t2a.md) |                        wheel scrolls | 🔄 **briefed** | [tui-t2a.md](../briefs/tui-t2a.md) |
    ┃    17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable  17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable
    ┃         from a script | 🔄 **briefed** | [tui-t2b.md](../briefs/tui-t2b.md) |                        from a script | 🔄 **briefed** | [tui-t2b.md](../briefs/tui-t2b.md) |
    ┃    18   | **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an 18   | **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an
    ┃         idea; D-3, D-5 | ⏳ planned | — |                                                            idea; D-3, D-5 | ⏳ planned | — |
    ┃    19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |              19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/tui/README.md && git commit -q -m "work: T1's fix cycle 1 closes the slice" && git status --short && make check 2>&1 | tail -25
    ┃
    ┃  (18 earlier lines)
    ┃  GO=go scripts/dep-budget.sh
    ┃    internal/core: 0 external modules (ceiling 0)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check-race 2>&1 | tail -14 && echo "=== STATUS ===" && git status --short && echo "(clean if empty)" && git log --oneline -8
    ┃
    ┃  (17 earlier lines)
    ┃  55bdc81 docs: the window title carries the project, and the brief's example was wrong
    ┃  b4a3052 fix(ui): the row labelled root shows the root
    ┃  67b455d fix(ui): the name the console renders is sanitised
    ┃  5a79638 feat(core): the plan's name has one answer
    ┃  e031fe8 work: T1's fix cycle 1 — four things the review surfaced
    ┃  164c883 inbox: reviewer reports tui-t1
    ┃  42dafdb docs: slice T1 in the ledger - the identity holds, and five low notes
    ┃

       + Thought · 430ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make test PKG=./internal/core RUN='-run DomainShapeIsVersioned -count=1 -v' 2>&1 | tail -8 && make build 2>&1 | tail -2
    ┃
    ┃  ==> test
    ┃  GO=go scripts/go-test.sh " -run DomainShapeIsVersioned -count=1 -v" "./internal/core"
    ┃  === RUN   TestDomainShapeIsVersioned
    ┃  --- PASS: TestDomainShapeIsVersioned (0.00s)
    ┃  PASS
    ┃  ok    github.com/lthiagol/maitre/internal/core  0.236s
    ┃  go build -o bin/maitre-cli ./cmd/maitre-cli
    ┃  go build -o bin/maitre ./cmd/maitre
    ┃

    ┃
    ┃  $ cd /tmp && rm -rf f99x && mkdir -p "f99x/we
    ┃  ird" && cd "f99x/we…
    ┃
    ┃  (4 earlier lines)
    ┃  === F-99 file on disk ===
    ┃    "name": "we\nird",
    ┃  === F-100 empty name ===
    ┃   ⋈ maitre · f100x
    ┃  === F-100 missing file ===
    ┃   ⋈ maitre · f100x
    ┃

    ┃
    ┃  $ cd /tmp/f100x && echo "=== F-101 root row ===" && /Users/thiago/Code/Pessoal/maitre/bin/maitre --render --all --plain | grep -A 2 "PLAN" && echo "=== F-101 themes ===" && /Users/thiago/
    ┃  Code/Pessoal/maitre/bin/maitre --render --all --plain | grep -A 8 "THEME" | head -9 && echo "=== F-102 doc ===" && grep -n "window title carries the project" -A 8 /Users/thiago/Cod…
    ┃
    ┃  (18 earlier lines)
    ┃  237-my four terminals is this?*, and it is **not** the screen's row 0: that row is the
    ┃  238-console's header, drawn inside the terminal, and it carries the same name.
    ┃  239-
    ┃  240-- The **list scrolls** to keep the selection visible. Section headers scroll with their
    ┃  241-  items, so a header can disappear while its items stay selected.
    ┃  242-- The **footer is pinned** to the last row and is **contextual**: `space` means different
    ┃

       + Thought · 4.1s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "GetProject\|LoadProject\|ProjectName" internal/ui/ internal/cli/ internal/core/*.go | grep -v "_test.go"
    ┃
    ┃  (8 earlier lines)
    ┃  internal/core/store.go:124:    name = s.defaultProjectName()
    ┃  internal/core/store.go:129:// ProjectName returns the plan's name.
    ┃  internal/core/store.go:137:func (s *Store) ProjectName() string {
    ┃  internal/core/store.go:138:  p, err := s.LoadProject()
    ┃  internal/core/store.go:140:    return s.defaultProjectName()
    ┃  internal/core/store.go:145:// defaultProjectName is the name a plan that names nothing answers with: the
    ┃  internal/core/store.go:150:func (s *Store) defaultProjectName() string {
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && mkdir -p work/inbox && cat > work/inbox/runner-tui-t1-fix1.md <<EOF
    ┃  Runner: T1 fix 1 done. Commits: 5 (5a79638, 67b455d, b4a3052, 55bdc81, e218ec3). Fixed: F-99, F-100, F-101, F-102.…
    ┃
    ┃  caa545e inbox: runner reports tui-t1 fix 1
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1H:p1" "Runner: T1 fix 1 done. Commits: 5 (5a79638, 67b455d, b4a3052, 55bdc81, e218ec3). Fixed: F-99, F-100, F-101, F-102. T
    ┃  he name is sanitised: a plan whose directory basename carries a newline (maitre-cli init writes 'we\nird') renders the header as ONE row — ' ⋈ maitre · we ird' — and kitchen/projec…
    ┃
    ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e737d6341ffeTdoEYTgykU25M8"},
    ┃  "agent_status":"done","completion_seq":63,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1H:p1","revision":15,
    ┃  "screen_detection_skipped":true,"state_change_seq":63,"tab_id":"w1H:t1","terminal_id":"term_65f0cd2aa89d26","terminal_title":"OC | Project work proposal","terminal_title_stripped":"OC |
    ┃  Project work proposal","workspace_id":"w1H"},"type":"agent_prompted"}}
    ┃

       T1 fix cycle 1 — closed. 6 commits (5 slice + the receipt), both gates green, tree clean.

       ┌─────────┬────────────────────────────────────────────────────────────────────────────────────┬──────────────────────────────────────────────────────────────────────────────────────────┐
       │ Finding │ What landed                                                                        │ Evidence                                                                                 │
       ├─────────┼────────────────────────────────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-99    │ the name goes through sanitize, like every other rendered string                   │ a plan whose basename carries a newline (maitre-cli init writes "we\nird") renders  ⋈    │
       │         │                                                                                    │ maitre · we ird as one row; kitchen/project.json still holds "we\nird"                   │
       ├─────────┼────────────────────────────────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-100   │ the domain answers the name once: core.Store.ProjectName(); core.LoadProject now   │ empty →  ⋈ maitre · f100x; missing file →  ⋈ maitre · f100x — the same answer, a name (  │
       │         │ answers the same for an empty name as for a missing file                           │ before: /tmp/f100x/kitchen vs f100x)                                                     │
       ├─────────┼────────────────────────────────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-101   │ the root row shows the plan directory                                              │ /tmp/f100x/kitchen, and ~/maitre-t1-evidence/kitchen under the home directory; the seven │
       │         │                                                                                    │ themes are still listed                                                                  │
       ├─────────┼────────────────────────────────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-102   │ one paragraph in docs/05-human-interface.md                                        │ it says the window title carries the project and that it is not the screen's row 0; the  │
       │         │                                                                                    │ sketch is unchanged                                                                      │
       └─────────┴────────────────────────────────────────────────────────────────────────────────────┴──────────────────────────────────────────────────────────────────────────────────────────┘

       Commits — 5a79638 the domain's one answer · 67b455d the name is sanitised · b4a3052 the root row shows the root · 55bdc81 the doc paragraph and the brief's correction · e218ec3 the state
       row.

       Invariants — core.DomainVersion is 7 and TestDomainShapeIsVersioned passes: the change is in the reader only (core.Init still writes the name it writes today, the format on disk does not
       move). Budget 0 / 0 / 28 against 0 / 0 / 30. The tabs are the same seven. docs/audit.md untouched. The goldens are unchanged — the fixture's project is demo and the header's words did not
       move.

       Notes for the reviewer

       - The new core.Store.ProjectName sits inside the declared Store.* family in core.NonViews, so it needs no row in docs/03 § 8 — the specification does not move.
       - The console no longer computes a name of its own: model.projectName() asks the domain and sanitises what comes back. grep GetProject internal/ui/ finds no call site.
       - The CLI did not change: maitre-cli plan prints the name the domain now answers (f100b) — that is the whole of the change.
       - The Settings row's test change landed in 67b455d with the other UI test changes; the row's own change is b4a3052.
       - T1's brief was corrected in 55bdc81 where it was wrong — the root's example and the name's source — and each correction says so in the brief.

       The receipt is committed at work/inbox/runner-tui-t1-fix1.md, and the nudge to w1H:p1 landed (agent_prompted).

       Build · MiniMax-M3.1-Flash · 9m 53s · 131.5 tok/s

    ┃
    ┃  You are the **runner**, again — the same lane that wrote T1 and its fix cycle 1. The
    ┃  orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command below.
    ┃
    ┃  This is **T1's fix cycle 2**. It is small: **one blocker, two notes, one question**.
    ┃
    ┃  ## Cycle 2's verdict
    ┃
    ┃  `changes-needed`, **one blocker: F-103**. The cycle-1 fixes (F-99..F-102) were all
    ┃  verified and hold. F-103 is a real one, and the reviewer found it by mutation — so
    ┃  read the proof before you fix it.
    ┃
    ┃  ## F-103 (medium, blocker) — the name's rule is written twice, and only one copy is pinned
    ┃
    ┃  **The target:** the rule that answers *"what is this plan called?"* is written **once**,
    ┃  and **every** surface that shows the name reads that one answer.
    ┃
    ┃  **The hole, proved by command.** The rule lives in two places:
    ┃
    ┃  | Where | What it says |
    ┃  |---|---|
    ┃  | `core.LoadProject` | a plan whose `project.json` names it **empty** answers with the project's directory name |
    ┃  | `core.Store.ProjectName` | `p.Name == ""` answers with the project's directory name — **again** |
    ┃
    ┃  Only the second copy is pinned by a test. `core.LoadProject` has **no test of its own**,
    ┃  so the first copy can be deleted and nothing notices. The reviewer's mutation, which
    ┃  the orchestrator reproduced:
    ┃
    ┃  ```bash
    ┃  # take the empty-name branch out of LoadProject
    ┃  go test ./...          # GREEN — the whole suite
    ┃  maitre-cli plan        # "Name": ""     <- the CLI's answer
    ┃  maitre --render        # ⋈ maitre · demo2  <- the console's answer
    ┃  ```
    ┃
    ┃  **Two surfaces, two answers, and a green gate.** That is the house's standing
    ┃  failure — *two answers to one question* — and the ledger has three rows of it already.
    ┃
    ┃  **What lands:** one copy of the rule, and it is pinned. Choose which copy survives —
    ┃  that is a judgement, and you own it:
    ┃
    ┃  - **Drop `ProjectName`'s redundant branch** and let `LoadProject` be the one place
    ┃    that answers. Then **`LoadProject`'s answer is pinned by a test of its own**, and
    ┃    `ProjectName` is the console's way of asking.
    ┃  - **Or keep the rule in `ProjectName` alone** and have `LoadProject` answer raw.
    ┃
    ┃  **Either is fine. What is not fine is both.** Say in your report which you chose and
    ┃  why.
    ┃
    ┃  **The evidence this finding demands:**
    ┃
    ┃  | | |
    ┃  |---|---|
    ┃  | the **mutation fails** | take the rule out of whichever copy you did not keep — **the suite must go red**. A green suite with the rule broken is the defect, not the fix |
    ┃  | the **two surfaces agree** | `maitre-cli plan` and `maitre --render` show the **same** name for the same plan, in all three states: the plan names itself · the file names it empty · the
    ┃  file is missing |
    ┃  | **the answer is pinned where the rule lives** | the copy that answers has a test that fails when the rule is removed from **that** copy |
    ┃
    ┃  **The CLI's behaviour may change, and that is allowed.** Today `maitre-cli plan` prints
    ┃  `"Name": ""` for a plan that names nothing, while the console shows the project's
    ┃  directory name. After the fix, both show the same thing. If your choice makes the CLI
    ┃  print the same name as the console, say so in your report — that is the point.
    ┃
    ┃  ## F-104 (low) — the `root` row renders the plan's path unsanitised
    ┃
    ┃  **The target:** a newline anywhere in the path the Settings tab shows does not split
    ┃  the row.
    ┃
    ┃  The fix cycle sanitised the **name**. The **path** goes through a different door, and
    ┃  a plan whose directory name carries a newline splits the row:
    ┃
    ┃  ```
    ┃    root        /tmp/vfy/we
    ┃  ird/kitchen
    ┃  ```
    ┃
    ┃  **What lands:** the path the console renders goes through the same sanitising as
    ┃  everything else the console renders. The path is not the plan's content — the review
    ┃  agreed — but the console still draws it on one line.
    ┃
    ┃  ## F-105 (low) — the state row's count and its list must agree
    ┃
    ┃  **The target:** the row cannot be read two ways.
    ┃
    ┃  Cycle 1 filed this as a note, and the fix cycle repeated it: the row says **12
    ┃  commits** and names **nine**. The three it does not name are the three commits that
    ┃  wrote the row — and a commit cannot carry its own hash, but the *next* row-writing
    ┃  commit can name the previous one.
    ┃
    ┃  The row is the **resume record**: the thing a cold session reads to know where the
    ┃  work stopped. **What lands:** the row's shape is one that cannot drift again — the
    ┃  count counts what the list names, **or** the list names every commit it can and the
    ┃  row says the count excludes the commits that wrote it. **Say which you choose.**
    ┃
    ┃  This is the row **you** write, in **this** cycle: `work/tui/README.md`'s T1 row. The
    ┃  orchestrator has already set it to `🔄 fix cycle 2`; your commit moves it to closed
    ┃  with the shape you chose.
    ┃
    ┃  ## The question the reviewer answered — nothing to do, and why
    ┃
    ┃  > *the name in `docs/03`: it does NOT belong in § 8 — that table is the derived
    ┃  > views (computed, never stored) and the name is a stored fact the domain reads,
    ┃  > like `LoadProject` and `ListRecipes`, which the table does not carry either; no
    ┃  > MEDIUM filed on that question, but the rule itself is nowhere in the spec.*
    ┃
    ┃  **The reviewer is right, and nothing lands here.** The plan's name is not a derived
    ┃  view: it is stored, in `project.json`, and read. `docs/03` § 8 carries the views the
    ┃  domain *computes*, and it does not carry `LoadProject` or `ListRecipes` either. So the
    ┃  name does not belong there.
    ┃
    ┃  **But the reviewer also said the rule is nowhere in the spec** — and that is worth
    ┃  saying where it can be found. **Where the rule should live** is the chef's call, not
    ┃  yours and not the orchestrator's. **Do not write it into `docs/03`.** If you think
    ┃  the rule belongs somewhere specific, say so in your report and the orchestrator will
    ┃  put it to the chef.
    ┃
    ┃  ## The scope — what not to touch
    ┃
    ┃  | | |
    ┃  |---|---|
    ┃  | **The CLI** | no verb is added, and `maitre-cli`'s command surface does not change. `maitre-cli plan` prints the name the domain answers, and that is the whole of the change |
    ┃  | **The writer** | `core.Init` still writes what it wrote. The change is in the **reader** |
    ┃  | **`core.DomainVersion`** | stays **7**. The plan format on disk does not move: the same fields, the same values. `TestDomainShapeIsVersioned` is the instrument — if it fails, **stop and
    ┃  report** |
    ┃  | **The tabs** | seven, the same names |
    ┃  | **The goldens** | `make goldens` regenerates; `make goldens-check` holds them. **Expected unchanged** — the header's words do not move. If a golden moves, say so and say why |
    ┃  | **`docs/03`** | **this cycle touches no `docs/03`.** The name is not a view, and the question of where the rule lives is the chef's |
    ┃  | **`docs/05`** | the paragraph cycle 1 added stays |
    ┃  | **`docs/audit.md`** | the ledger is the **reviewer's**. This cycle files no finding and touches no row |
    ┃  | **`work/`** | the orchestrator's, except the state row (F-105) |
    ┃
    ┃  ## The scope — what lands
    ┃
    ┃  | | |
    ┃  |---|---|
    ┃  | **`internal/core`** | one copy of the name's rule, and it is pinned (F-103) |
    ┃  | **`internal/ui`** | the path the console renders is sanitised (F-104) |
    ┃  | **`work/tui/README.md`** | the state row's shape, so it cannot drift again (F-105) |
    ┃
    ┃  ## The invariants
    ┃
    ┃  - `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
    ┃    budget (**≤ 30**). **No dependency is added by this cycle.**
    ┃  - **The domain layer has zero UI imports.**
    ┃  - **`core.DomainVersion` is 7.**
    ┃  - The tab count is 7 and the tab names do not change.
    ┃  - **One fact, one answer.**
    ┃  - The gate is `make check`. A focused check is `make test PKG=… RUN='-run …'`.
    ┃
    ┃  ## The evidence this cycle must produce
    ┃
    ┃  For **each** finding, a command and its output — not a claim:
    ┃
    ┃  | Finding | The evidence |
    ┃  |---|---|
    ┃  | **F-103** | the **mutation fails**: removing the rule from the copy you did not keep turns the suite **red**. And the two surfaces — `maitre-cli plan` and `maitre --render` — show the
    ┃  **same** name in all three states |
    ┃  | **F-104** | a plan whose directory name carries a newline shows a **one-row** `root` row |
    ┃  | **F-105** | the row's count and its list agree, and you say in your report which shape you chose and why |
    ┃
    ┃  Plus: `make check` green, `make check-race` green, the goldens unchanged, and the tree
    ┃  clean.
    ┃
    ┃  ## The report — the receipt comes FIRST
    ┃
    ┃  ```bash
    ┃  mkdir -p work/inbox
    ┃  cat > work/inbox/runner-tui-t1-fix2.md <<EOF
    ┃  Runner: T1 fix 2 done. Commits: <n>. Fixed: F-103, F-104, F-105. The rule lives in: <which copy you kept, and why>. The mutation: <what you removed, and the suite's answer>. The two
    ┃  surfaces: <what maitre-cli plan and maitre --render show, in all three states>. The path is sanitised: <what you proved it with>. The state row: <the shape you chose, and why>.
    ┃  DomainVersion: <7, and the version test's result>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
    ┃  EOF
    ┃  git add work/inbox && git commit -q -m "inbox: runner reports tui-t1 fix 2"
    ┃  herdr agent prompt "w1H:p1" "Runner: T1 fix 2 done. Commits: <n>. Fixed: F-103, F-104, F-105. The rule lives in: <which copy you kept, and why>. The mutation: <what you removed, and the
    ┃  suite's answer>. The two surfaces: <what maitre-cli plan and maitre --render show, in all three states>. The path is sanitised: <what you proved it with>. The state row: <the shape you
    ┃  chose, and why>. DomainVersion: <7, and the version test's result>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer
    ┃  should know>."
    ┃  ```
    ┃
    ┃  **Write the receipt, commit it, and *then* nudge** — the nudge only reaches an idle
    ┃  receiver; the receipt is the half that cannot be lost.
    ┃
    ┃  **One write per file, and commit as soon as the file is coherent.** Every edit to a
    ┃  large file re-renders that file's whole diff into the terminal, and a turn that renders
    ┃  dozens of diffs runs out before the work is done.
    ┃
    ┃  ## Before you report
    ┃
    ┃  ```bash
    ┃  make check && git status --short
    ┃  ```
    ┃
    ┃  Green, and the tree clean. Report the **real** numbers: the commits, the gate's
    ┃  result, and what the two surfaces show.
    ┃
    ┃
    ┃  ---
    ┃
    ┃  ## The notify command — verbatim, use it as it is
    ┃
    ┃  The orchestrator's pane id is **w1H:p1**. Your receipt is the record; the nudge is
    ┃  the bell. Send **both**, in this order, and the nudge **last**:
    ┃
    ┃  ```bash
    ┃  mkdir -p work/inbox
    ┃  cat > work/inbox/runner-tui-t1-fix2.md <<EOF
    ┃  Runner: T1 fix 2 done. Commits: <n>. Fixed: F-103, F-104, F-105. The rule lives in: <which copy you kept, and why>. The mutation: <what you removed, and the suite's answer>. The two
    ┃  surfaces: <what maitre-cli plan and maitre --render show, in all three states>. The path is sanitised: <what you proved it with>. The state row: <the shape you chose, and why>.
    ┃  DomainVersion: <7, and the version test's result>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
    ┃  EOF
    ┃  git add work/inbox && git commit -q -m "inbox: runner reports tui-t1 fix 2"
    ┃  herdr agent prompt "w1H:p1" "Runner: T1 fix 2 done. Commits: <n>. Fixed: F-103, F-104, F-105. The rule lives in: <which copy you kept, and why>. The mutation: <what you removed, and the
    ┃  suite's answer>. The two surfaces: <what maitre-cli plan and maitre --render show, in all three states>. The path is sanitised: <what you proved it with>. The state row: <the shape you
    ┃  chose, and why>. DomainVersion: <7, and the version test's result>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer
    ┃  should know>."
    ┃  ```
    ┃
    ┃  **If the nudge fails** — it only reaches an idle receiver — the receipt is the half
    ┃  that cannot be lost. Write the receipt anyway, and say so in it.
    ┃
    ┃  **Nothing else waits on you.** After the nudge, end your turn. Do not poll for a
    ┃  reply, do not re-dispatch yourself, do not start another slice. The orchestrator
    ┃  decides what happens next.
    ┃
    ┃  ## Two things this brief cannot know
    ┃
    ┃  - **The blocker is real, and the reviewer proved it by mutation.** The orchestrator
    ┃    reproduced the mutation independently before ruling: the suite stays green with the
    ┃    rule broken. Read the proof in the brief before you fix it - the fix is a judgement
    ┃    about *which copy survives*, and you own that judgement.
    ┃  - **The reviewer answered the docs/03 question: the name is not a view.** So this
    ┃    cycle touches no docs/03, and the question of where the *rule* should live is the
    ┃    chef's - the orchestrator will put it to the chef, not to you.
    ┃

       + Thought · 46.0s

       → Explored: 2 reads

       + Thought · 15.6s

       → Explored: 1 search



    ┃
    ┃  ← Edit internal/core/store.go
    ┃
    ┃    104     Settings map[string]any `json:"settings"`                                            104     Settings map[string]any `json:"settings"`
    ┃    105   }                                                                                      105   }
    ┃    106                                                                                          106
    ┃    107   // LoadProject reads the plan-level metadata.                                          107   // LoadProject reads the plan-level metadata.
    ┃                                                                                                 108 + //
    ┃                                                                                                 109 + // The name is not answered here: the domain answers it once — see ProjectName — and
    ┃                                                                                                 110 + // this asks. Two copies of the rule would be two answers to one question.
    ┃    108   func (s *Store) LoadProject() (Project, error) {                                       111   func (s *Store) LoadProject() (Project, error) {
    ┃    109 -   raw, err := os.ReadFile(filepath.Join(s.Root, FileProject))                          112 +   doc, err := s.readProjectFile()
    ┃    110     if err != nil {                                                                      113     if err != nil {
    ┃    111       if os.IsNotExist(err) {                                                            114       if os.IsNotExist(err) {
    ┃    112 -       return Project{Name: s.defaultProjectName()}, nil                                115 +       return Project{Name: s.projectName(doc)}, nil
    ┃    113       }                                                                                  116       }
    ┃    114 -     return Project{}, fmt.Errorf("read project: %w", err)                              117 +     // The name is still the answer even when the file cannot be read: the
    ┃                                                                                                 118 +     // caller that ignores the error still gets the plan's name.
    ┃                                                                                                 119 +     return Project{Name: s.projectName(doc)}, fmt.Errorf("read project: %w", err)
    ┃    115     }                                                                                    120     }
    ┃    116 -   var doc projectFile                                                                  121 +   return Project{Name: s.projectName(doc), Created: doc.Created}, nil
    ┃    117 -   if err := json.Unmarshal(raw, &doc); err != nil {
    ┃    118 -     return Project{}, fmt.Errorf("parse project: %w", err)
    ┃    119 -   }
    ┃    120 -   name := doc.Name
    ┃    121 -   if name == "" {
    ┃    122 -     // A plan that names nothing answers with the same name a plan without a
    ┃    123 -     // project file answers with. One question, one answer.
    ┃    124 -     name = s.defaultProjectName()
    ┃    125 -   }
    ┃    126 -   return Project{Name: name, Created: doc.Created}, nil
    ┃    127   }                                                                                      122   }
    ┃    128                                                                                          123
    ┃    129 - // ProjectName returns the plan's name.                                                124 + // ProjectName returns the plan's name — the answer to *what is this plan called?*.
    ┃    130   //                                                                                     125   //
    ┃    131 - // One question, one answer: the name the plan's own file gives it, and the project's  126 + // It asks the rule and gives the answer back, so every reader of the name — the
    ┃    132 - // directory name when the plan names nothing. A missing file and an empty name are    127 + // console, the CLI — gets the same one.
    ┃    133 - // the same case, and they answer the same.
    ┃    134   //                                                                                     128   //
    ┃    135 - // The answer is **never a path**. A path in the header would say *which directory*    129 + // The answer is **never a path**. A path would say *which directory* where the
    ┃    136 - // where the chef asked *which project*, and the two are different questions.          130 + // chef asked *which project*, and the two are different questions.
    ┃    137   func (s *Store) ProjectName() string {                                                 131   func (s *Store) ProjectName() string {
    ┃    138 -   p, err := s.LoadProject()                                                            132 +   doc, _ := s.readProjectFile()
    ┃    139 -   if err != nil || p.Name == "" {                                                      133 +   return s.projectName(doc)
    ┃                                                                                                 134 + }
    ┃                                                                                                 135 +
    ┃                                                                                                 136 + // projectName is the rule, in one place.
    ┃                                                                                                 137 + //
    ┃                                                                                                 138 + // The plan's own name, and the project's directory name when the plan names
    ┃                                                                                                 139 + // nothing: a missing file, an unreadable file and an empty name are the same case.
    ┃                                                                                                 140 + func (s *Store) projectName(doc projectFile) string {
    ┃                                                                                                 141 +   if doc.Name == "" {
    ┃    140       return s.defaultProjectName()                                                      142       return s.defaultProjectName()
    ┃    141     }                                                                                    143     }
    ┃    142 -   return p.Name                                                                        144 +   return doc.Name
    ┃    143   }                                                                                      145   }
    ┃    144                                                                                          146
    ┃                                                                                                 147 + // readProjectFile reads and parses the plan's project file.
    ┃                                                                                                 148 + //
    ┃                                                                                                 149 + // A missing file is an error here, and the caller decides what a plan without one
    ┃                                                                                                 150 + // answers — the file carries no name to read.
    ┃                                                                                                 151 + func (s *Store) readProjectFile() (projectFile, error) {
    ┃                                                                                                 152 +   raw, err := os.ReadFile(filepath.Join(s.Root, FileProject))
    ┃                                                                                                 153 +   if err != nil {
    ┃                                                                                                 154 +     return projectFile{}, err
    ┃                                                                                                 155 +   }
    ┃                                                                                                 156 +   var doc projectFile
    ┃                                                                                                 157 +   if err := json.Unmarshal(raw, &doc); err != nil {
    ┃                                                                                                 158 +     return projectFile{}, err
    ┃                                                                                                 159 +   }
    ┃                                                                                                 160 +   return doc, nil
    ┃                                                                                                 161 + }
    ┃                                                                                                 162 +
    ┃    145   // defaultProjectName is the name a plan that names nothing answers with: the          163   // defaultProjectName is the name a plan that names nothing answers with: the
    ┃    146   // project's directory.                                                                164   // project's directory.
    ┃    147   //                                                                                     165   //
    ┃    148   // The plan directory is `<project>/kitchen`, so the project's directory is the plan   166   // The plan directory is `<project>/kitchen`, so the project's directory is the plan
    ┃

       + Thought · 2.2s

       → Explored: 1 search

       + Thought · 3.3s

    ┃
    ┃  ← Edit internal/core/store.go
    ┃
    ┃    115         return Project{Name: s.projectName(doc)}, nil                                    115         return Project{Name: s.projectName(doc)}, nil
    ┃    116       }                                                                                  116       }
    ┃    117       // The name is still the answer even when the file cannot be read: the             117       // The name is still the answer even when the file cannot be read: the
    ┃    118       // caller that ignores the error still gets the plan's name.                       118       // caller that ignores the error still gets the plan's name.
    ┃    119 -     return Project{Name: s.projectName(doc)}, fmt.Errorf("read project: %w", err)      119 +     return Project{Name: s.projectName(doc)}, err
    ┃    120     }                                                                                    120     }
    ┃    121     return Project{Name: s.projectName(doc), Created: doc.Created}, nil                  121     return Project{Name: s.projectName(doc), Created: doc.Created}, nil
    ┃    122   }                                                                                      122   }
    ┃    123                                                                                          123
    ┃

    ┃
    ┃  ← Edit internal/core/store.go
    ┃
    ┃    150   // answers — the file carries no name to read.                                         150   // answers — the file carries no name to read.
    ┃    151   func (s *Store) readProjectFile() (projectFile, error) {                               151   func (s *Store) readProjectFile() (projectFile, error) {
    ┃    152     raw, err := os.ReadFile(filepath.Join(s.Root, FileProject))                          152     raw, err := os.ReadFile(filepath.Join(s.Root, FileProject))
    ┃    153     if err != nil {                                                                      153     if err != nil {
    ┃    154 -     return projectFile{}, err                                                          154 +     return projectFile{}, fmt.Errorf("read project: %w", err)
    ┃    155     }                                                                                    155     }
    ┃    156     var doc projectFile                                                                  156     var doc projectFile
    ┃    157     if err := json.Unmarshal(raw, &doc); err != nil {                                    157     if err := json.Unmarshal(raw, &doc); err != nil {
    ┃    158 -     return projectFile{}, err                                                          158 +     return projectFile{}, fmt.Errorf("parse project: %w", err)
    ┃    159     }                                                                                    159     }
    ┃    160     return doc, nil                                                                      160     return doc, nil
    ┃    161   }                                                                                      161   }
    ┃    162                                                                                          162
    ┃

       Now the pinned test for the rule's home, and the path sanitising (F-104):

    ┃
    ┃  ← Edit internal/core/store_test.go
    ┃
    ┃     60       t.Errorf("the name is a path: %q", s.ProjectName())                                 60       t.Errorf("the name is a path: %q", s.ProjectName())
    ┃     61     }                                                                                     61     }
    ┃     62   }                                                                                       62   }
    ┃     63                                                                                           63
    ┃                                                                                                  64 + // TestLoadProjectAnswersTheSameName proves the two surfaces agree.
    ┃                                                                                                  65 + //
    ┃                                                                                                  66 + // `maitre-cli plan` reads the project through LoadProject, and the console reads it
    ┃                                                                                                  67 + // through ProjectName. Two answers to one question is the house's standing failure,
    ┃                                                                                                  68 + // so the rule lives in one place and both surfaces ask it.
    ┃                                                                                                  69 + func TestLoadProjectAnswersTheSameName(t *testing.T) {
    ┃                                                                                                  70 +   root := filepath.Join(t.TempDir(), "demo", PlanDirName)
    ┃                                                                                                  71 +   s := Open(root)
    ┃                                                                                                  72 +   if err := s.Init(); err != nil {
    ┃                                                                                                  73 +     t.Fatalf("init: %v", err)
    ┃                                                                                                  74 +   }
    ┃                                                                                                  75 +   project := filepath.Join(s.Root, FileProject)
    ┃                                                                                                  76 +
    ┃                                                                                                  77 +   cases := []struct {
    ┃                                                                                                  78 +     what string
    ┃                                                                                                  79 +     plan func()
    ┃                                                                                                  80 +   }{
    ┃                                                                                                  81 +     {"the plan names itself", func() {}},
    ┃                                                                                                  82 +     {"the plan names nothing", func() {
    ┃                                                                                                  83 +       raw, err := os.ReadFile(project)
    ┃                                                                                                  84 +       if err != nil {
    ┃                                                                                                  85 +         t.Fatalf("read project: %v", err)
    ┃                                                                                                  86 +       }
    ┃                                                                                                  87 +       emptied := strings.Replace(string(raw), `"name": "demo"`, `"name": ""`, 1)
    ┃                                                                                                  88 +       if err := os.WriteFile(project, []byte(emptied), 0o644); err != nil {
    ┃                                                                                                  89 +         t.Fatalf("write project: %v", err)
    ┃                                                                                                  90 +       }
    ┃                                                                                                  91 +     }},
    ┃                                                                                                  92 +     {"the plan has no project file", func() {
    ┃                                                                                                  93 +       if err := os.Remove(project); err != nil {
    ┃                                                                                                  94 +         t.Fatalf("remove project: %v", err)
    ┃                                                                                                  95 +       }
    ┃                                                                                                  96 +     }},
    ┃                                                                                                  97 +   }
    ┃                                                                                                  98 +
    ┃                                                                                                  99 +   for _, c := range cases {
    ┃                                                                                                 100 +     c.plan()
    ┃                                                                                                 101 +     p, err := s.LoadProject()
    ┃                                                                                                 102 +     if err != nil {
    ┃                                                                                                 103 +       t.Fatalf("%s: load project: %v", c.what, err)
    ┃                                                                                                 104 +     }
    ┃                                                                                                 105 +     if p.Name != s.ProjectName() {
    ┃                                                                                                 106 +       t.Errorf("%s: LoadProject answers %q and ProjectName answers %q", c.what, p.
    ┃                                                                                                       Name, s.ProjectName())
    ┃                                                                                                 107 +     }
    ┃                                                                                                 108 +     if p.Name != "demo" {
    ┃                                                                                                 109 +       t.Errorf("%s: the name is %q, want %q", c.what, p.Name, "demo")
    ┃                                                                                                 110 +     }
    ┃                                                                                                 111 +   }
    ┃                                                                                                 112 + }
    ┃                                                                                                 113 +
    ┃     64   func TestInitCreatesTheLayout(t *testing.T) {                                          114   func TestInitCreatesTheLayout(t *testing.T) {
    ┃     65     s := newStore(t)                                                                     115     s := newStore(t)
    ┃     66     for _, d := range []string{DirRecipes, DirCourses, DirBacklog} {                     116     for _, d := range []string{DirRecipes, DirCourses, DirBacklog} {
    ┃     67       info, err := os.Stat(filepath.Join(s.Root, d))                                     117       info, err := os.Stat(filepath.Join(s.Root, d))
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && gofmt -w internal/core/store.go internal/core/store_test.go && make test PKG=./internal/core RUN='-run ProjectNameHasOneAnswer|LoadProjectAnswers
    ┃  TheSameName -count=1 -v' 2>&1 | tail -12
    ┃
    ┃  (5 earlier lines)
    ┃      store_test.go:103: the plan has no project file: load project: read project: open /var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/TestLoadProjectAnswersTheSameName2449662189/001/
    ┃  demo/kitchen/project.json: no such file or directory
    ┃  --- FAIL: TestLoadProjectAnswersTheSameName (0.01s)
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/core  0.361s
    ┃  FAIL
    ┃  make: *** [test] Error 1
    ┃

       + Thought · 2.2s

    ┃
    ┃  ← Edit internal/core/store.go
    ┃
    ┃    110   // this asks. Two copies of the rule would be two answers to one question.             110   // this asks. Two copies of the rule would be two answers to one question.
    ┃    111   func (s *Store) LoadProject() (Project, error) {                                       111   func (s *Store) LoadProject() (Project, error) {
    ┃    112     doc, err := s.readProjectFile()                                                      112     doc, err := s.readProjectFile()
    ┃    113     if err != nil {                                                                      113     if err != nil {
    ┃    114 -     if os.IsNotExist(err) {                                                            114 +     // The error is wrapped, so the test for a missing file has to unwrap it.
    ┃                                                                                                 115 +     if errors.Is(err, os.ErrNotExist) {
    ┃    115         return Project{Name: s.projectName(doc)}, nil                                    116         return Project{Name: s.projectName(doc)}, nil
    ┃    116       }                                                                                  117       }
    ┃    117       // The name is still the answer even when the file cannot be read: the             118       // The name is still the answer even when the file cannot be read: the
    ┃    118       // caller that ignores the error still gets the plan's name.                       119       // caller that ignores the error still gets the plan's name.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && gofmt -w internal/core/store.go && make test PKG=./internal/core RUN='-run ProjectNameHasOneAnswer|LoadProjectAnswersTheSameName -count=1 -v' 2>&
    ┃  1 | tail -10
    ┃
    ┃  (2 earlier lines)
    ┃  === RUN   TestProjectNameHasOneAnswer
    ┃  --- PASS: TestProjectNameHasOneAnswer (0.01s)
    ┃  === RUN   TestLoadProjectAnswersTheSameName
    ┃  --- PASS: TestLoadProjectAnswersTheSameName (0.01s)
    ┃  PASS
    ┃  ok    github.com/lthiagol/maitre/internal/core  0.239s
    ┃

       Now the mutation check — the rule removed from its one home must turn the suite red:

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/core/store.go /tmp/store.go.bak && python3 - <<'EOF'
    ┃  p='internal/core/store.go'…
    ┃
    ┃  (16 earlier lines)
    ┃  --- FAIL: TestTheNameHasOneAnswer (0.71s)
    ┃      identity_test.go:85: the name is "", want the project's directory name "demo"
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  65.508s
    ┃  FAIL
    ┃  make: *** [test] Error 1
    ┃

       + Thought · 507ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp /tmp/store.go.bak internal/core/store.go && git diff --stat internal/core/store.go && make test PKG=./internal/core RUN='-run ProjectName|Load
    ┃  Project -count=1' 2>&1 | tail -4
    ┃
    ┃  internal/core/store.go | 67 ++++++++++++++++++++++++++++++++------------------
    ┃   1 file changed, 43 insertions(+), 24 deletions(-)
    ┃
    ┃  ==> test
    ┃  GO=go scripts/go-test.sh " -run ProjectName|LoadProject -count=1" "./internal/core"
    ┃  ok    github.com/lthiagol/maitre/internal/core  0.283s
    ┃

       Now F-104 — the path the console renders is sanitised:

    ┃
    ┃  ← Edit internal/ui/views.go
    ┃
    ┃    1026     var out []listLine                                                                  1026     var out []listLine
    ┃    1027     add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }   1027     add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }
    ┃    1028     add("")                                                                             1028     add("")
    ┃    1029     add(" " + m.stySec.Render("PLAN"))                                                  1029     add(" " + m.stySec.Render("PLAN"))
    ┃    1030 -   add("  " + m.styDim.Render("root        ") + humanPath(m.store.Root))               1030 +   add("  " + m.styDim.Render("root        ") + humanPath(sanitize(m.store.Root)))
    ┃    1031     add("  " + m.styDim.Render("recipes     ") + fmt.Sprint(len(m.recipes)))            1031     add("  " + m.styDim.Render("recipes     ") + fmt.Sprint(len(m.recipes)))
    ┃    1032     add("  " + m.styDim.Render("backlog     ") + fmt.Sprint(len(m.items)))              1032     add("  " + m.styDim.Render("backlog     ") + fmt.Sprint(len(m.items)))
    ┃    1033     add("  " + m.styDim.Render("events      ") + fmt.Sprint(len(m.journal.Events)))     1033     add("  " + m.styDim.Render("events      ") + fmt.Sprint(len(m.journal.Events)))
    ┃    1034     add("")                                                                             1034     add("")
    ┃

    ┃
    ┃  ← Edit internal/ui/identity.go
    ┃
    ┃    51   //                                                                                      51   //
    ┃    52   // A path with the home directory spelled out is longer than it needs to be, and the    52   // A path with the home directory spelled out is longer than it needs to be, and the
    ┃    53   // part that is the same in every one of the chef's paths is the part that says         53   // part that is the same in every one of the chef's paths is the part that says
    ┃    54   // nothing.                                                                             54   // nothing.
    ┃                                                                                                 55 + //
    ┃                                                                                                 56 + // The path is **sanitised** first: it is drawn on one row, and a directory whose
    ┃                                                                                                 57 + // name carries a newline would split it.
    ┃    55   func humanPath(path string) string {                                                    58   func humanPath(path string) string {
    ┃                                                                                                 59 +   path = sanitize(path)
    ┃    56     home, err := os.UserHomeDir()                                                         60     home, err := os.UserHomeDir()
    ┃    57     if err != nil || home == "" {                                                         61     if err != nil || home == "" {
    ┃    58       return path                                                                         62       return path
    ┃    59     }                                                                                     63     }
    ┃

       + Thought · 4.5s

    ┃
    ┃  ← Edit internal/ui/views.go
    ┃
    ┃    1026     var out []listLine                                                                  1026     var out []listLine
    ┃    1027     add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }   1027     add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }
    ┃    1028     add("")                                                                             1028     add("")
    ┃    1029     add(" " + m.stySec.Render("PLAN"))                                                  1029     add(" " + m.stySec.Render("PLAN"))
    ┃    1030 -   add("  " + m.styDim.Render("root        ") + humanPath(sanitize(m.store.Root)))     1030 +   add("  " + m.styDim.Render("root        ") + humanPath(m.store.Root))
    ┃    1031     add("  " + m.styDim.Render("recipes     ") + fmt.Sprint(len(m.recipes)))            1031     add("  " + m.styDim.Render("recipes     ") + fmt.Sprint(len(m.recipes)))
    ┃    1032     add("  " + m.styDim.Render("backlog     ") + fmt.Sprint(len(m.items)))              1032     add("  " + m.styDim.Render("backlog     ") + fmt.Sprint(len(m.items)))
    ┃    1033     add("  " + m.styDim.Render("events      ") + fmt.Sprint(len(m.journal.Events)))     1033     add("  " + m.styDim.Render("events      ") + fmt.Sprint(len(m.journal.Events)))
    ┃    1034     add("")                                                                             1034     add("")
    ┃

       Now a test for the path's row (F-104):

    ┃
    ┃  ← Edit internal/ui/identity_test.go
    ┃
    ┃    258       }                                                                                  258       }
    ┃    259     }                                                                                    259     }
    ┃    260   }                                                                                      260   }
    ┃    261                                                                                          261
    ┃                                                                                                 262 + // TestTheRootRowIsOneRow proves a path with a newline in it does not split the row.
    ┃                                                                                                 263 + //
    ┃                                                                                                 264 + // The path is not the plan's content — but the console draws it on one line, and a
    ┃                                                                                                 265 + // directory whose name carries a newline would take two.
    ┃                                                                                                 266 + func TestTheRootRowIsOneRow(t *testing.T) {
    ┃                                                                                                 267 +   // A plan whose directory name carries a newline.
    ┃                                                                                                 268 +   root := filepath.Join(t.TempDir(), "we\nird", core.PlanDirName)
    ┃                                                                                                 269 +   store := core.Open(root)
    ┃                                                                                                 270 +   if err := core.Seed(store); err != nil {
    ┃                                                                                                 271 +     t.Fatalf("seed: %v", err)
    ┃                                                                                                 272 +   }
    ┃                                                                                                 273 +   m := loadInto(t, New(store))
    ┃                                                                                                 274 +   m.tab = TabSettings
    ┃                                                                                                 275 +
    ┃                                                                                                 276 +   var text strings.Builder
    ┃                                                                                                 277 +   for _, line := range m.settingsLines() {
    ┃                                                                                                 278 +     text.WriteString(line.text + "\n")
    ┃                                                                                                 279 +   }
    ┃                                                                                                 280 +   rows := text.String()
    ┃                                                                                                 281 +
    ┃                                                                                                 282 +   for _, line := range strings.Split(strings.TrimRight(rows, "\n"), "\n") {
    ┃                                                                                                 283 +     if strings.Contains(line, "\n") || strings.Contains(line, "\r") {
    ┃                                                                                                 284 +       t.Errorf("the root row is split: %q", line)
    ┃                                                                                                 285 +     }
    ┃                                                                                                 286 +   }
    ┃                                                                                                 287 +   if !strings.Contains(rows, "we ird") {
    ┃                                                                                                 288 +     t.Errorf("the root row does not carry the sanitised path\n%s", rows)
    ┃                                                                                                 289 +   }
    ┃                                                                                                 290 +   if strings.Contains(rows, "we\nird") {
    ┃                                                                                                 291 +     t.Errorf("the root row carries a newline\n%s", rows)
    ┃                                                                                                 292 +   }
    ┃                                                                                                 293 + }
    ┃                                                                                                 294 +
    ┃    262   // TestHumanPathShortensTheHomeDirectory proves the path is formatted for a human.     295   // TestHumanPathShortensTheHomeDirectory proves the path is formatted for a human.
    ┃    263   //                                                                                     296   //
    ┃    264   // The home directory is the same in every one of the chef's paths, and the part       297   // The home directory is the same in every one of the chef's paths, and the part
    ┃    265   // that is the same is the part that says nothing.                                     298   // that is the same is the part that says nothing.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && gofmt -w internal/ui/identity.go internal/ui/identity_test.go && make test PKG=./internal/ui RUN='-run RootRow|HumanPath|SettingsTab -count=1 -v'
    ┃   2>&1 | tail -12
    ┃
    ┃  (6 earlier lines)
    ┃  === RUN   TestHumanPathShortensTheHomeDirectory
    ┃  --- PASS: TestHumanPathShortensTheHomeDirectory (0.00s)
    ┃  === RUN   TestTheSettingsTabPicksTheTheme
    ┃  --- PASS: TestTheSettingsTabPicksTheTheme (0.79s)
    ┃  PASS
    ┃  ok    github.com/lthiagol/maitre/internal/ui  3.180s
    ┃

       + Thought · 535ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/ui/identity.go /tmp/identity.go.bak && python3 - <<'EOF'
    ┃  p='internal/ui/identity.go'…
    ┃
    ┃  (7 earlier lines)
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  0.957s
    ┃  FAIL
    ┃  make: *** [test] Error 1
    ┃  ==> test
    ┃  GO=go scripts/go-test.sh " -run RootRow -count=1" "./internal/ui"
    ┃  ok    github.com/lthiagol/maitre/internal/ui  1.128s
    ┃

       + Thought · 4.2s

       → Explored: 1 read

       + Thought · 1.1s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make test PKG=./internal/core 2>&1 | tail -3 && make test PKG=./internal/ui 2>&1 | tail -3 && make goldens-check 2>&1 | tail -3 && git status --s
    ┃  hort
    ┃
    ┃  (7 earlier lines)
    ┃  ok    github.com/lthiagol/maitre/internal/ui  5.003s
    ┃    every golden file is claimed by a test
    ┃   M internal/core/store.go
    ┃   M internal/core/store_test.go
    ┃   M internal/ui/identity.go
    ┃   M internal/ui/identity_test.go
    ┃

       + Thought · 367ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/core/store.go internal/core/store_test.go && git commit -q -m "fix(core): the name's rule lives in one place
    ┃  …
    ┃
    ┃  8a6866b fix(core): the name's rule lives in one place
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/ui/identity.go internal/ui/identity_test.go && git commit -q -m "fix(ui): the path the console renders is sanitised
    ┃  …
    ┃
    ┃  cf8c05e fix(ui): the path the console renders is sanitised
    ┃

       Now the state row (F-105):

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃    11   ## Where it stands                                                                      11   ## Where it stands
    ┃    12                                                                                           12
    ┃    13   | Session | Feature | What it lands | State | Brief |                                   13   | Session | Feature | What it lands | State | Brief |
    ┃    14   |--:|---|---|---|---|                                                                   14   |--:|---|---|---|---|
    ┃    15 - | **T1** | the identity | the window title carries the project; the header's name stops 15 + | **T1** | the identity | the window title carries the project; the header's name stops
    ┃         being a footnote; `--project`; the Settings tab shows the plan's path | 🔄 **fix cycle       being a footnote; `--project`; the Settings tab shows the plan's path | ✅ closed · 11
    ┃         2** — 9 commits · 94a53ff, 8fc7438, f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d,             commits · 94a53ff, 8fc7438, f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052,
    ┃         b4a3052, 55bdc81 | [tui-t1.md](../briefs/tui-t1.md) · [fix 1](../briefs/fix-tui-t1-1.        55bdc81, 8a6866b, cf8c05e | [tui-t1.md](../briefs/tui-t1.md) · [fix 1](../briefs/fix-
    ┃         md) · [fix 2](../briefs/fix-tui-t1-2.md) |                                                   tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |
    ┃    16   | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the     16   | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the
    ┃         wheel scrolls | 🔄 **briefed** | [tui-t2a.md](../briefs/tui-t2a.md) |                        wheel scrolls | 🔄 **briefed** | [tui-t2a.md](../briefs/tui-t2a.md) |
    ┃    17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable  17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable
    ┃         from a script | 🔄 **briefed** | [tui-t2b.md](../briefs/tui-t2b.md) |                        from a script | 🔄 **briefed** | [tui-t2b.md](../briefs/tui-t2b.md) |
    ┃    18   | **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an 18   | **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an
    ┃         idea; D-3, D-5 | ⏳ planned | — |                                                            idea; D-3, D-5 | ⏳ planned | — |
    ┃    19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |              19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |
    ┃

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃    31   their session — a brief written today is stale by the time the slice runs.              31   their session — a brief written today is stale by the time the slice runs.
    ┃    32                                                                                           32
    ┃    33   **The commits** land in the table's own row as they happen: `T1 · 2 commits ·           33   **The commits** land in the table's own row as they happen: `T1 · 2 commits ·
    ┃    34   3f9c1a2, 8b0e441`. A row with commits under it and no `closed` is a session that        34   3f9c1a2, 8b0e441`. A row with commits under it and no `closed` is a session that
    ┃    35 - did not finish its report.                                                              35 + did not finish its report. **The count is the number of ids the row names**, and the
    ┃                                                                                                 36 + commits that wrote the row are not among them — a commit cannot carry its own hash.
    ┃    36                                                                                           37
    ┃    37   ## The session protocol                                                                 38   ## The session protocol
    ┃    38                                                                                           39
    ┃    39   **One session, one slice.** The slices are independent enough that they can be done     40   **One session, one slice.** The slices are independent enough that they can be done
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/tui/README.md && git commit -q -m "work: T1's fix cycle 2 closes the slice, and the row's shape is written down
    ┃  …
    ┃
    ┃  4e67b31 work: T1's fix cycle 2 closes the slice, and the row's shape is written down
    ┃  cf8c05e fix(ui): the path the console renders is sanitised
    ┃  8a6866b fix(core): the name's rule lives in one place
    ┃  264f0c3 work: T1's fix cycle 2 — one blocker, two notes, one question
    ┃

       Now the evidence, by command:

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make build 2>&1 | tail -2 && cd /tmp && rm -rf f2 && mkdir -p f2 && cd f2 && B=/Users/thiago/Code/Pessoal/maitre/bin && $B/maitre-cli init >/dev/
    ┃  null && echo "=== state 1: the plan names itself ===" && $B/maitre-cli plan | python3 -c "import json,sys;print('cli   :',json.load(sys.stdin)['project']['Name'])" && $B/maitre --r…
    ┃
    ┃  (5 earlier lines)
    ┃  === state 2: the file names it empty ===
    ┃  cli   : f2
    ┃   ⋈ maitre · f2
    ┃  === state 3: the file is missing ===
    ┃  cli   : f2
    ┃   ⋈ maitre · f2
    ┃

       + Thought · 880ms

    ┃
    ┃  $ cd /tmp && rm -rf f2x && mkdir -p "f2x/we
    ┃  ird" && cd "f2x/we…
    ┃
    ┃  (2 earlier lines)
    ┃    root        /tmp/f2x/we ird/kitchen$
    ┃    recipes     0$
    ┃  === the header, one row ===
    ┃   ⋈ maitre · we ird$
    ┃  ───────────────────────────────────────────────────────────\
    ┃  ─────────────────────────────────────────$
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check 2>&1 | tail -16 && echo "=== STATUS ===" && git status --short && echo "(clean if empty)"
    ┃
    ┃  (11 earlier lines)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃  === STATUS ===
    ┃  (clean if empty)
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check-race 2>&1 | tail -10 && make test PKG=./internal/core RUN='-run DomainShapeIsVersioned -count=1 -v' 2>&1 | tail -5
    ┃
    ┃  (8 earlier lines)
    ┃
    ┃  ==> tidy-check
    ┃  GO=go scripts/go-test.sh " -run DomainShapeIsVersioned -count=1 -v" "./internal/core"
    ┃  === RUN   TestDomainShapeIsVersioned
    ┃  --- PASS: TestDomainShapeIsVersioned (0.00s)
    ┃  PASS
    ┃  ok    github.com/lthiagol/maitre/internal/core  0.206s
    ┃

       + Thought · 6.2s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && mkdir -p work/inbox && cat > work/inbox/runner-tui-t1-fix2.md <<EOF
    ┃  Runner: T1 fix 2 done. Commits: 3 (8a6866b, cf8c05e, 4e67b31). Fixed: F-103, F-104, F-105.…
    ┃
    ┃  6165daf inbox: runner reports tui-t1 fix 2
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1H:p1" "Runner: T1 fix 2 done. Commits: 3 (8a6866b, cf8c05e, 4e67b31). Fixed: F-103, F-104, F-105. The rule lives in: one pr
    ┃  ivate function, core.Store.projectName(doc), and everything asks it — LoadProject, ProjectName, and through them the CLI and the console. I chose the helper over letting LoadProjec…
    ┃
    ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e737d6341ffeTdoEYTgykU25M8"},
    ┃  "agent_status":"done","completion_seq":80,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1H:p1","revision":15,
    ┃  "screen_detection_skipped":true,"state_change_seq":80,"tab_id":"w1H:t1","terminal_id":"term_65f0cd2aa89d26","terminal_title":"OC | Project work proposal","terminal_title_stripped":"OC |
    ┃  Project work proposal","workspace_id":"w1H"},"type":"agent_prompted"}}
    ┃

       T1 fix cycle 2 — closed. 4 commits (3 slice + the receipt), both gates green, tree clean.

       ┌─────────┬────────────────────────────────────────────────────────────────────────────────────┬──────────────────────────────────────────────────────────────────────────────────────────┐
       │ Finding │ What landed                                                                        │ Evidence                                                                                 │
       ├─────────┼────────────────────────────────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-103   │ the rule lives in one private function — core.Store.projectName(doc) — and         │ the mutation (removing the empty-name branch) turned the suite red: internal/core failed │
       │         │ LoadProject, ProjectName, the CLI and the console all ask it                       │ TestProjectNameHasOneAnswer + TestLoadProjectAnswersTheSameName, internal/ui failed      │
       │         │                                                                                    │ TestTheNameHasOneAnswer. Reverted.                                                       │
       ├─────────┼────────────────────────────────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-104   │ the path is sanitised inside humanPath, so every path the console draws is         │ a plan under we\nird shows   root        /tmp/f2x/we ird/kitchen — one row (sed -n l     │
       │         │ sanitised by construction                                                          │ prints a single line)                                                                    │
       ├─────────┼────────────────────────────────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-105   │ the count counts what the list names: 11 commits, 11 ids; the row-writing commits  │ the row                                                                                  │
       │         │ are not named (a commit cannot carry its own hash), and the protocol section now   │                                                                                          │
       │         │ says so in one line                                                                │                                                                                          │
       └─────────┴────────────────────────────────────────────────────────────────────────────────────┴──────────────────────────────────────────────────────────────────────────────────────────┘

       The two surfaces agree in all three states — the plan names itself · the file names it empty · the file is missing:

       cli : f2      ⋈ maitre · f2
       cli   : f2      ⋈ maitre · f2
       cli   : f2      ⋈ maitre · f2

       Before the fix the CLI printed "" for the empty-name state while the console printed the directory name.

       Why the rule moved into a private helper — LoadProject and ProjectName are each asked by a different surface (the CLI and the console), so a rule either of them carries can be deleted
       without the other noticing. That was exactly F-103. One copy, three askers, and the pin sits on the copy.

       Notes for the reviewer

       - The CLI's behaviour changed, as the brief allowed: maitre-cli plan prints the domain's answer (f2) instead of "".
       - The error strings did not change (read project: %w, parse project: %w); a missing file is now found with errors.Is(err, os.ErrNotExist) because the error is wrapped. LoadProject carries
         the name even when the file cannot be read, so a caller that ignores the error still gets it.
       - docs/03 was not touched — the name is not a view. Where the rule belongs in the spec is the chef's call; my reading is that it belongs with the plan's anatomy (the format section), not
         § 8. I wrote it nowhere — it is in the code, where it is enforced.
       - core.DomainVersion is7 and TestDomainShapeIsVersioned passes; the writer and the format on disk are untouched. Budget 0 / 0 / 28.

       The receipt is committed at work/inbox/runner-tui-t1-fix2.md, and the nudge to w1H:p1 landed (agent_prompted).

       Build · MiniMax-M3.1-Flash · 11m 34s · 105.2 tok/s

    ┃
    ┃
    ┃
    ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
    ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
    ~/Code/Pessoal/maitre:master                                                                                                                                     182.1K (18%)  ctrl+p commands
```
