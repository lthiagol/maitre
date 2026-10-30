# The runner's lane — T2a, captured 2026-10-30

Dumped from pane `w1H:p4` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that wrote T2a: the mouse is on, the render builds the hit map, the wheel
scrolls, and the mouse is drivable from a script.
The review that judged it is in [tui-t2a-reviewer.md](tui-t2a-reviewer.md).

       + Thought · 3.8s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && (make check-race > /tmp/t2a-race.log 2>&1; echo "EXIT=$?" >> /tmp/t2a-race.log) &
    ┃  sleep 5; echo started
    ┃
    ┃   Background
    ┃

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃    12                                                                        12
    ┃    13   | Session | Feature | What it lands | State | Brief |                13   | Session | Feature | What it lands | State | Brief |
    ┃    14   |--:|---|---|---|---|                                                14   |--:|---|---|---|---|
    ┃    15   | **T1** | the identity | the window title carries the project; the  15   | **T1** | the identity | the window title carries the project; the
    ┃         header's name stops being a footnote; `--project`; the Settings tab       header's name stops being a footnote; `--project`; the Settings tab
    ┃         shows the plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438,        shows the plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438,
    ┃         f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81,            f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81,
    ┃         8a6866b, cf8c05e | [tui-t1.md](../briefs/tui-t1.md) · [fix 1](../         8a6866b, cf8c05e | [tui-t1.md](../briefs/tui-t1.md) · [fix 1](../
    ┃         briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |            briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |
    ┃    16 - | **T2a** | the mouse, part 1 | the mouse is on; the render builds   16 + | **T2a** | the mouse, part 1 | the mouse is on; the render builds
    ┃         the hit map; the wheel scrolls | 🔄 **review ok · fix cycle 1** — 7       the hit map; the wheel scrolls | 🔄 **fix cycle 1 done** — 10
    ┃         commits · a72ea98, 115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7,           commits · a72ea98, 115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7,
    ┃         e4b7061 | [tui-t2a.md](../briefs/tui-t2a.md) · [fix 1](../briefs/         e4b7061, 619b342, 47be828, f1cc1ca | [tui-t2a.md](../briefs/tui-t2a.
    ┃         fix-tui-t2a-1.md) |                                                       md) · [fix 1](../briefs/fix-tui-t2a-1.md) |
    ┃    17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail;  17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail;
    ┃         the mouse drivable from a script | 🔄 **briefed** · next | [tui-t2b.      the mouse drivable from a script | 🔄 **briefed** · next | [tui-t2b.
    ┃         md](../briefs/tui-t2b.md) |                                               md](../briefs/tui-t2b.md) |
    ┃    18   | **T3a** | the backlog is readable | the detail pane; `space` stops 18   | **T3a** | the backlog is readable | the detail pane; `space` stops
    ┃         being a trap on an idea; D-3, D-5 | ⏳ planned | — |                      being a trap on an idea; D-3, D-5 | ⏳ planned | — |
    ┃    19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned 19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned
    ┃         | — |                                                                     | — |
    ┃    20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef    20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef
    ┃         dissolves; D-2 | ⏳ planned | — |                                         dissolves; D-2 | ⏳ planned | — |
    ┃

       + Thought · 4.0s

    ┃
    ┃  $ sleep 120; tail -5 /tmp/t2a-race.log
    ┃
    ┃  ?     github.com/lthiagol/maitre/cmd/maitre  [no test files]
    ┃  ?     github.com/lthiagol/maitre/cmd/maitre-cli  [no test files]
    ┃  ok    github.com/lthiagol/maitre/internal/brand  1.865s
    ┃  ok    github.com/lthiagol/maitre/internal/cli  15.548s
    ┃  ok    github.com/lthiagol/maitre/internal/core  111.028s
    ┃

       ↳ Shell finished · cd /Users/thiago/Code/Pessoal/maitre && (make check-race > /tmp/t2a-race.log 2>&1; echo "EXIT=$?" >> /tmp/t2a-race.log) & sleep 5…

       + Thought · 393ms

    ┃
    ┃  $ sleep 180; tail -12 /tmp/t2a-race.log
    ┃
    ┃  (5 earlier lines)
    ┃    internal/core: 0 external modules (ceiling 0)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃  EXIT=0
    ┃

       + Thought · 959ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/ui/hit.go /tmp/t2a-final/hit.go.keep && cp internal/ui/mouse.go /tmp/t2a-final/mouse.go.keep &&
    ┃   git show 61e64c8:internal/ui/hit.go > internal/ui/hit.go && git show 61e64c8:internal/ui/mouse.go > internal/ui/mouse.go && go test ./interna…
    ┃
    ┃  (4 earlier lines)
    ┃  internal/ui/views.go:287:25: undefined: newHitBox
    ┃  internal/ui/mouse_test.go:437:11: undefined: SplitSteps
    ┃  internal/ui/mouse_test.go:457:24: undefined: SplitSteps
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui [build failed]
    ┃  FAIL
    ┃   M work/tui/README.md
    ┃

       + Thought · 1.5s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
    ┃  p="internal/ui/hit.go"…
    ┃
    ┃  (7 earlier lines)
    ┃      mouse_test.go:509: a cleared map should have no detail
    ┃      mouse_test.go:512: a cleared map should hold nothing
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  4.712s
    ┃  FAIL
    ┃   M work/tui/README.md
    ┃

       + Thought · 2.8s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make build >/dev/null 2>&1 && echo "=== F-106: the split layout's detail, in the package ===" && go test ./
    ┃  internal/ui -run 'TestTheWheelScrollsTheSplitDetail|TestTheScrollsBoundIsThePanesWidth' -v 2>&1 | tail -8 && echo && echo "=== F-106: the spli…
    ┃
    ┃  (15 earlier lines)
    ┃  ---
    ┃  >  Needs you                                                                              2
    ┃  > ────────────────────────────────────────────────────────────────────────────────────────── Outcome
    ┃  >   RCP02 ready       Sign off the menu                             1/2
    ┃  >   RCP03 remediation Wire the review queue                         2/3 1 note                   The chef can approve a recipe from the console
    ┃  >  On the line                                                                            2
    ┃

    ┃
    ┃  $ cd /tmp/kb-t2a && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre && echo "=== F-107: the position form on the command line ===" && $M --render --ke
    ┃  ys "wheel:down@(12,5)" --plain | head -6; echo "exit=$?" && echo && echo "=== F-107: the old spelling ===" && $M --render --keys "wheel:down@1…
    ┃
    ┃  (21 earlier lines)
    ┃  === RUN   TestTheOriginIsNotAPhantomDetail
    ┃  --- PASS: TestTheOriginIsNotAPhantomDetail (0.71s)
    ┃  === RUN   TestAClearedMapAnswersNothing
    ┃  --- PASS: TestAClearedMapAnswersNothing (0.76s)
    ┃  PASS
    ┃  ok    github.com/lthiagol/maitre/internal/ui  1.608s
    ┃

       + Thought · 247ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check 2>&1 | tail -12 && echo "--- status ---" && git status --short && echo "--- goldens ---" && git
    ┃  diff HEAD --stat -- internal/ui/testdata && echo "(no golden diff above = unchanged)"
    ┃
    ┃  (10 earlier lines)
    ┃
    ┃  ==> tidy-check
    ┃  --- status ---
    ┃   M work/tui/README.md
    ┃  --- goldens ---
    ┃  (no golden diff above = unchanged)
    ┃



    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && mkdir -p work/inbox && cat > work/inbox/runner-tui-t2a-fix1.md <<'EOF'
    ┃  Runner: T2a fix 1 done. Commits: 4 (619b342, 47be828, f1cc1ca, plus the state row). Fixed: F-106, F-107, F-108.…
    ┃
    ┃  90e4e96 inbox: runner reports tui-t2a fix 1
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1H:p1" "Runner: T2a fix 1 done. Commits: 4 (619b342, 47be828, f1cc1ca, 90e4e96). Fixed
    ┃  : F-106, F-107, F-108. The detail's scroll: the offset is applied once, where the detail is drawn (Model.detailScrollRows), and renderSplit no…
    ┃
    ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":
    ┃  "ses_e737d6341ffeTdoEYTgykU25M8"},"agent_status":"done","completion_seq":132,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,
    ┃  "foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1H:p1","revision":15,"screen_detection_skipped":true,"state_change_seq":132,
    ┃  "tab_id":"w1H:t1","terminal_id":"term_65f0cd2aa89d26","terminal_title":"OC | Project work proposal","terminal_title_stripped":"OC | Project work
    ┃  proposal","workspace_id":"w1H"},"type":"agent_prompted"}}
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && git log --oneline -6
    ┃
    ┃  90e4e96 inbox: runner reports tui-t2a fix 1
    ┃  f1cc1ca fix(ui): the position form survives the command line
    ┃  47be828 fix(ui): an absent region is absent
    ┃  619b342 fix(ui): the detail pane scrolls in both layouts
    ┃  61e64c8 work: T2a's fix cycle 1 - the split detail, the position form, the absent region
    ┃  4ada8f0 inbox: reviewer reports tui-t2a
    ┃

       T2a's fix cycle 1 is landed. The receipt is committed and the nudge was delivered to w1H:p1.

       The three findings

       ┌─────────┬───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┬─────────┐
       │ Finding │ What landed                                                                                                                   │ Commit  │
       ├─────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼─────────┤
       │ F-106   │ the offset is applied once, where the detail is drawn — Model.detailScrollRows, called by renderDetailFull, renderSplit and   │ 619b342 │
       │         │ viewDetail. The bound is counted at the width the pane draws with: Model.detailPaneWidth() (the full screen's width-4, or     │         │
       │         │ width-listColumn()-3 in the split), which renderSplit now asks too                                                            │         │
       ├─────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼─────────┤
       │ F-107   │ the position is bracketed — wheel:down@(12,5). The steps stay comma-separated, so --keys "tab,down,space" still works; the    │ f1cc1ca │
       │         │ splitter that knows the grammar (ui.SplitSteps) lives beside the parser that reads it                                         │         │
       ├─────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼─────────┤
       │ F-108   │ hitBox carries a there flag: the zero box is not there and answers nothing. The wheel falls back to the list where the        │ 47be828 │
       │         │ pointer is not over the detail                                                                                                │         │
       └─────────┴───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┴─────────┘

       The evidence, as commands:

       - F-106 — go test ./internal/ui -run TestTheWheelScrollsTheSplitDetail -v → PASS (the offset moves and the frame differs). And from the command line:
         --render alone never draws the split layout, but --keys drives the real View, so diff <(maitre --render --width 200 --keys "wheel:up@(150,6)" --
         plain) <(… "wheel:down@(150,6)" …) shows the detail column moving three rows.
       - F-107 — maitre --render --keys "wheel:down@(12,5)" --plain → exit 0. The old spelling now fails loudly: maitre: malformed mouse step "wheel:
         down@12". docs/05 prints the bracketed form.
       - F-108 — proved by mutation, not claim: removing the there guard (the pre-fix behaviour) makes TestTheOriginIsNotAPhantomDetail fail with the
         finding's own symptom — "the wheel at the origin scrolled the detail to 3" — and with the fix the wheel at the origin scrolls the list and the
         cleared map answers nothing.

       make check green · make check-race green · goldens unchanged · tree clean · DomainVersion 7 · 28 external modules, none added.

       One behaviour change to look at: the wheel over the header now scrolls the list (before: nothing). That is what the reviewer's F-108 evidence asks
       for — a wheel at the origin over the list scrolls the list — and the fallback stops where the list is not on screen, so the full-screen detail still
       ignores the wheel over its header. It is called out in the receipt for the re-review.

       Build · MiniMax-M3.1-Flash · 25m 35s · 136.4 tok/s

    ┃
    ┃  You are the **runner**, again — the same lane that wrote T2a and its fix cycle 1. The
    ┃  orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command below.
    ┃
    ┃  This is **T2a's fix cycle 2**: **one blocker and one line of a document**.
    ┃
    ┃  ## The review's verdict
    ┃
    ┃  `changes-needed`, **one blocker: F-109**. Cycle 2 verified F-106, F-107 and F-108
    ┃  by command and by mutation, and it **judged the reversed wheel test sound** — the
    ┃  fallback is one rule, in one place, stated once in `docs/05`. That question is
    ┃  closed; do not reopen it.
    ┃
    ┃  ## F-109 (medium, blocker) — the bound's width is not held by anything
    ┃
    ┃  **The target:** the scroll stops at the pane's last row, and the width the bound is
    ┃  counted at is **held by a test that fails when the width moves**.
    ┃
    ┃  **The hole, proved by command.** The fix cycle put the bound and the pane on one
    ┃  answer — `Model.detailPaneWidth()` — and added a test, `TestTheScrollsBoundIsThePanesWidth`,
    ┃  that checks the bound equals the pane's rows. **The reviewer's mutation leaves the whole
    ┃  `internal/ui` suite green:**
    ┃
    ┃  ```bash
    ┃  # the BOUND counted at the full-screen width again
    ┃  # internal/core/model.go:  detailLineCount() → m.width - 4
    ┃  go test ./internal/ui/      # GREEN — every test passes
    ┃  ```
    ┃
    ┃  **Why the test does not catch it.** Its first assertion compares two values that
    ┃  **both** go through the same function:
    ┃
    ┃  ```go
    ┃  pane := len(detailRows(m.detailLines(m.detailPaneWidth())))
    ┃  if m.detailLineCount() != pane { ... }
    ┃  ```
    ┃
    ┃  Both sides ask `detailPaneWidth()` — so when the **bound** is computed somewhere else,
    ┃  the comparison only notices **if the two widths give different row counts**. And the
    ┃  fixture's detail is too short for that: at the split pane's width and at the full
    ┃  screen's width it wraps to the **same** number of rows. The test passes with the bound
    ┃  on the wrong width.
    ┃
    ┃  **What lands:** a test that **pins the width**, with a fixture where the two widths
    ┃  **differ** — a detail long enough that the narrower pane wraps it into more rows.
    ┃  Then the bound computed at the full-screen width **fails**.
    ┃
    ┃  **The evidence this finding demands:**
    ┃
    ┃  | | |
    ┃  |---|---|
    ┃  | the **mutation fails** | count the bound at the full-screen width, and the suite must go **red**. Say which test failed and what it said |
    ┃  | the **bound is the pane's** | the scroll reaches the pane's **last row** in the **split** layout — not a row short of it |
    ┃  | the **widths differ** | your fixture makes the split pane's width and the full screen's width give **different** row counts — say how you know, and
    ┃  how many |
    ┃
    ┃  **Watch the trap:** the test must fail because the **bound's width** moved, not
    ┃  because some unrelated number moved. A test that fails for the wrong reason is not a
    ┃  tripwire.
    ┃
    ┃  ## The one line of the document
    ┃
    ┃  **The target:** the wheel's rule is stated **whole**.
    ┃
    ┃  `docs/05-human-interface.md`'s wheel table says:
    ┃
    ┃  | Gesture | Where | What |
    ┃  |---|---|---|
    ┃  | wheel | over the detail pane | scroll the detail |
    ┃  | wheel | anywhere else | scroll the list, keeping the selection visible |
    ┃
    ┃  **"Anywhere else" is not true of the full-screen detail**: there is no list on screen
    ┃  there, so the wheel over the header does **nothing**. The prose below the table says
    ┃  so — the table does not.
    ┃
    ┃  **What lands:** the table says what is true, and the rule reads whole. **One line**,
    ┃  or one row. Say in your report what you wrote.
    ┃
    ┃  ## The scope — what not to touch
    ┃
    ┃  | | |
    ┃  |---|---|
    ┃  | **The clicks** | **T2b.** No click handler |
    ┃  | **The wheel's fallback** | **closed.** The list is the console's one scrollable thing, and the rule is stated once. Do not change it, and do not
    ┃  reopen it |
    ┃  | **`maitre-cli`** | no verb is added and no command changes |
    ┃  | **The tabs** | seven, the same names |
    ┃  | **The domain** | no domain change. `internal/core` is untouched |
    ┃  | **`core.DomainVersion`** | stays **7** |
    ┃  | **The goldens** | expected unchanged — a golden that moved for no stated reason is the slice saying it did more than it said |
    ┃  | **`docs/03`** | untouched |
    ┃  | **`docs/audit.md`** | the ledger is the **reviewer's**. This cycle files no finding and touches no row |
    ┃  | **`work/`** | the orchestrator's, except the state row |
    ┃
    ┃  ## The notes the review filed, and left for a later round
    ┃
    ┃  **Do not fix these in this cycle.** They are notes, and a cycle that widens is a cycle
    ┃  that drifts:
    ┃
    ┃  - `detailScrollRows` writes `m.detailScroll` during a **render** — the clamping. It
    ┃    moved there from `renderDetailFull`; it is not new.
    ┃  - The wheel over the list in the split layout moves the **selection**, so the pane
    ┃    shows another recipe from its top and the detail does not scroll. That is the
    ┃    console's model.
    ┃
    ┃  ## The invariants
    ┃
    ┃  - `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
    ┃    budget. **No dependency is added.**
    ┃  - **The domain layer has zero UI imports.**
    ┃  - **The tab count is 7** and the tab names do not change.
    ┃  - **One fact, one answer.**
    ┃  - The gate is `make check`. A focused check is `make test PKG=… RUN=…`.
    ┃
    ┃  ## The evidence this cycle must produce
    ┃
    ┃  | Finding | The evidence |
    ┃  |---|---|
    ┃  | **F-109** | the mutation fails: counting the bound at the full-screen width turns the suite **red**, and the test that failed is about the
    ┃  **bound's width**. And your fixture makes the two widths give different row counts — say how many |
    ┃  | **the document** | the wheel's rule reads whole, and every row of the table is **true** |
    ┃
    ┃  Plus: `make check` green, `make check-race` green, the goldens unchanged, and the tree
    ┃  clean.
    ┃
    ┃  ## The report — the receipt comes FIRST
    ┃
    ┃  ```bash
    ┃  mkdir -p work/inbox
    ┃  cat > work/inbox/runner-tui-t2a-fix2.md <<EOF
    ┃  Runner: T2a fix 2 done. Commits: <n>. Fixed: F-109 and the document's wheel table. The bound's width: <the mutation you took, the test that failed,
    ┃  and what it said>. The fixture: <how the two widths come to differ, and the row counts>. The document: <what you wrote, and the table it reads as
    ┃  now>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should
    ┃  know>.
    ┃  EOF
    ┃  git add work/inbox && git commit -q -m "inbox: runner reports tui-t2a fix 2"
    ┃  herdr agent prompt "w1H:p1" "Runner: T2a fix 2 done. Commits: <n>. Fixed: F-109 and the document's wheel table. The bound's width: <the mutation you
    ┃  took, the test that failed, and what it said>. The fixture: <how the two widths come to differ, and the row counts>. The document: <what you wrote,
    ┃  and the table it reads as now>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes:
    ┃  <anything the reviewer should know>."
    ┃  ```
    ┃
    ┃  **Write the receipt, commit it, and *then* nudge** — the nudge only reaches an idle
    ┃  receiver; the receipt is the half that cannot be lost.
    ┃
    ┃  **One write per file, and commit as soon as the file is coherent.**
    ┃
    ┃  ## Before you report
    ┃
    ┃  ```bash
    ┃  make check && git status --short
    ┃  ```
    ┃
    ┃  Green, and the tree clean. Report the **real** numbers: the commits, the gate's
    ┃  result, and which test caught the mutation.
    ┃
    ┃
    ┃  ---
    ┃
    ┃  ## The notify command — verbatim, use it as it is
    ┃
    ┃  The orchestrator's pane id is **w1H:p1**. Your receipt is the record; the nudge is the
    ┃  bell. Send **both**, in this order, and the nudge **last**:
    ┃
    ┃  ```bash
    ┃  mkdir -p work/inbox
    ┃  cat > work/inbox/runner-tui-t2a-fix2.md <<EOF
    ┃  Runner: T2a fix 2 done. Commits: <n>. Fixed: F-109 and the document's wheel table. The bound's width: <the mutation you took, the test that failed,
    ┃  and what it said>. The fixture: <how the two widths come to differ, and the row counts>. The document: <what you wrote, and the table it reads as
    ┃  now>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should
    ┃  know>.
    ┃  EOF
    ┃  git add work/inbox && git commit -q -m "inbox: runner reports tui-t2a fix 2"
    ┃  herdr agent prompt "w1H:p1" "Runner: T2a fix 2 done. Commits: <n>. Fixed: F-109 and the document's wheel table. The bound's width: <the mutation you
    ┃  took, the test that failed, and what it said>. The fixture: <how the two widths come to differ, and the row counts>. The document: <what you wrote,
    ┃  and the table it reads as now>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes:
    ┃  <anything the reviewer should know>."
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
    ┃  - **The orchestrator reproduced the blocker before ruling.** The reviewer's mutation
    ┃    — counting the bound at the full-screen width — leaves the whole `internal/ui`
    ┃    suite green, exactly as the review said. Read the trap in the brief before you
    ┃    write the test: a test that fails for the wrong reason is not a tripwire.
    ┃  - **The wheel's fallback is closed.** The review judged the reversed test sound, and
    ┃    the rule stands: the list is the console's one scrollable thing, and the wheel
    ┃    scrolls it everywhere except over the detail. Do not change it, and do not reopen
    ┃    it — the document's table is what is being completed, not the behaviour.
    ┃

       + Thought · 23.4s



    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "type RecipeSpec" -A 30 internal/core/*.go | head -45
    ┃
    ┃  (24 earlier lines)
    ┃  internal/core/operations.go-76-}
    ┃  internal/core/operations.go-77-
    ┃  internal/core/operations.go-78-// Fragment is a partial update. A nil field is "leave alone"; a non-nil field replaces.
    ┃  internal/core/operations.go-79-// That is what makes the update fragment-first: one section at a time.
    ┃  internal/core/operations.go-80-type Fragment struct {
    ┃  internal/core/operations.go-81-  Title      *string
    ┃  internal/core/operations.go-82-  Outcome    *string
    ┃

       + Thought · 893ms

    ┃
    ┃  # Wrote internal/ui/bound_test.go
    ┃
    ┃    1 package ui
    ┃    2
    ┃    3 import (
    ┃    4   "path/filepath"
    ┃    5   "strings"
    ┃    6   "testing"
    ┃    7
    ┃    8   "github.com/lthiagol/maitre/internal/core"
    ┃    9 )
    ┃   10
    ┃   11 // longDetailModel builds a console over a recipe whose detail **wraps**.
    ┃   12 //
    ┃   13 // The pane is narrower in the split layout — the list takes a column of it — and a
    ┃   14 // narrower pane wraps the same text into more rows. A tripwire for the bound's
    ┃   15 // width needs a detail where the two widths disagree: with a short one they agree,
    ┃   16 // and the tripwire cannot fail.
    ┃   17 func longDetailModel(t *testing.T) *Model {
    ┃   18   t.Helper()
    ┃   19   root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
    ┃   20   store := core.Open(root)
    ┃   21   if err := store.Init(); err != nil {
    ┃   22     t.Fatal(err)
    ┃   23   }
    ┃   24   ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})
    ┃   25
    ┃   26   // A paragraph long enough to wrap several times at the pane's width, and fewer
    ┃   27   // times at the full screen's.
    ┃   28   outcome := strings.Repeat(
    ┃   29     "The chef reads the whole recipe in the console, so the detail pane has to be "+
    ┃   30       "as tall as the work it shows and never stop one row short of the end. ",
    ┃   31     4)
    ┃   32
    ┃   33   r, err := ops.CreateRecipe(core.RecipeSpec{
    ┃   34     Title:      "The detail pane scrolls to its last row",
    ┃   35     Outcome:    outcome,
    ┃   36     Problem:    "The scroll's bound was counted at the wrong width.",
    ┃   37     OutOfScope: []string{"web shell", "orchestration"},
    ┃   38     Criteria: []core.CriterionSpec{
    ┃   39       {Description: "the scroll reaches the pane's last row", Verification: "go test ./internal/ui -run TestTheScrollsBoundIsThePanesWidth"},
    ┃   40     },
    ┃   41     WorkPackages: []string{"ui"},
    ┃   42     Steps:        []core.StepSpec{{Action: "count the bound at the pane's width", DoneWhen: "the scroll reaches the end", Covers: []string{"AC-
    ┃      01"}}},
    ┃   43   })
    ┃   44   if err != nil {
    ┃   45     t.Fatal(err)
    ┃   46   }
    ┃   47
    ┃   48   m := loadModel(t, store)
    ┃   49   m.width = 200
    ┃   50   m.height = 20
    ┃   51   return m
    ┃   52 }
    ┃   53
    ┃   54 // TestTheScrollsBoundIsThePanesWidth covers the bound the scroll stops at.
    ┃   55 //
    ┃   56 // The bound is counted at the width the pane **draws with**, and this test pins that
    ┃   57 // width: the pane's rows are counted here, at the pane's own width, and not
    ┃   58 // through the console's answer — a tripwire that asks the code under test for the
    ┃   59 // answer is not a tripwire.
    ┃   60 //
    ┃   61 // The fixture makes the two widths disagree, or the comparison would be vacuous: the
    ┃   62 // bound could be counted at the full screen's width and the rows would be the same.
    ┃   63 func TestTheScrollsBoundIsThePanesWidth(t *testing.T) {
    ┃   64   m := longDetailModel(t)
    ┃   65   m.View()
    ┃   66
    ┃   67   // The split pane's width, computed here: the list's column, one column of air,
    ┃   68   // and the rest of the terminal.
    ┃   69   paneWidth := m.width - m.listColumn() - 3
    ┃   70   paneRows := len(detailRows(m.detailLines(paneWidth)))
    ┃   71   fullRows := len(detailRows(m.detailLines(m.width - 4)))
    ┃   72   if paneRows == fullRows {
    ┃   73     t.Fatalf("the fixture must make the two widths differ: %d rows at %d columns and %d at %d",
    ┃   74       paneRows, paneWidth, fullRows, m.width-4)
    ┃   75   }
    ┃   76
    ┃   77   // The bound is the pane's rows.
    ┃   78   if m.detailLineCount() != paneRows {
    ┃   79     t.Errorf("the bound is %d rows; the pane has %d rows at %d columns", m.detailLineCount(), paneRows, paneWidth)
    ┃   80   }
    ┃   81
    ┃   82   // And the scroll reaches the pane's last row.
    ┃   83   x, y := m.hits.detail.centre()
    ┃   84   for i := 0; i < 500; i++ {
    ┃   85     m.Update(tea.MouseWheelMsg{X: x, Y: y, Button: tea.MouseWheelDown})
    ┃   86   }
    ┃   87   if want := paneRows - 1; m.detailScroll != want {
    ┃   88     t.Errorf("the scroll stopped at %d, want the pane's last row %d", m.detailScroll, want)
    ┃   89   }
    ┃   90 }
    ┃   91
    ┃

       + Thought · 184ms

    ┃
    ┃  ← Edit internal/ui/mouse_test.go
    ┃
    ┃    398       t.Errorf("the frame did not move:\n%s", after)                  398       t.Errorf("the frame did not move:\n%s", after)
    ┃    399     }                                                                 399     }
    ┃    400   }                                                                   400   }
    ┃    401                                                                       401
    ┃    402 - // TestTheScrollsBoundIsThePanesWidth covers the bound the scroll
    ┃          stops at.
    ┃    403 - //
    ┃    404 - // The pane is narrower in the split layout, and a narrower pane
    ┃          wraps the same text
    ┃    405 - // into more rows. The bound is counted at the width the pane draws
    ┃          with, so the
    ┃    406 - // scroll can reach the pane's last row.
    ┃    407 - func TestTheScrollsBoundIsThePanesWidth(t *testing.T) {
    ┃    408 -   m := newModel(t)
    ┃    409 -   m.width = 200
    ┃    410 -   m.height = 20
    ┃    411 -   m.View()
    ┃    412 -
    ┃    413 -   pane := len(detailRows(m.detailLines(m.detailPaneWidth())))
    ┃    414 -   if m.detailLineCount() != pane {
    ┃    415 -     t.Errorf("the bound is %d rows, want the pane's %d", m.
    ┃          detailLineCount(), pane)
    ┃    416 -   }
    ┃    417 -   if m.detailPaneWidth() != m.width-m.listColumn()-3 {
    ┃    418 -     t.Errorf("the pane draws at %d columns, want %d", m.
    ┃          detailPaneWidth(), m.width-m.listColumn()-3)
    ┃    419 -   }
    ┃    420 -
    ┃    421 -   // The scroll reaches the pane's last row.
    ┃    422 -   x, y := m.hits.detail.centre()
    ┃    423 -   for i := 0; i < 200; i++ {
    ┃    424 -     m.Update(tea.MouseWheelMsg{X: x, Y: y, Button: tea.
    ┃          MouseWheelDown})
    ┃    425 -   }
    ┃    426 -   if want := m.detailLineCount() - 1; m.detailScroll != want {
    ┃    427 -     t.Errorf("the scroll stopped at %d, want %d", m.detailScroll,
    ┃          want)
    ┃    428 -   }
    ┃    429 - }
    ┃    430 -
    ┃    431   // TestThePositionFormSurvivesTheCommandLine covers the spelling on 402   // TestThePositionFormSurvivesTheCommandLine covers the spelling on
    ┃          the command                                                               the command
    ┃    432   // line.                                                            403   // line.
    ┃    433   //                                                                  404   //
    ┃    434   // The steps are separated by commas, so a position has to be       405   // The steps are separated by commas, so a position has to be
    ┃          bracketed: the comma                                                      bracketed: the comma
    ┃

    ┃
    ┃  ← Edit internal/ui/bound_test.go
    ┃
    ┃     4     "path/filepath"                                                     4     "path/filepath"
    ┃     5     "strings"                                                           5     "strings"
    ┃     6     "testing"                                                           6     "testing"
    ┃     7                                                                         7
    ┃                                                                               8 +   "charm.land/bubbletea/v2"
    ┃                                                                               9 +
    ┃     8     "github.com/lthiagol/maitre/internal/core"                         10     "github.com/lthiagol/maitre/internal/core"
    ┃     9   )                                                                    11   )
    ┃    10                                                                        12
    ┃    11   // longDetailModel builds a console over a recipe whose detail       13   // longDetailModel builds a console over a recipe whose detail
    ┃         **wraps**.                                                                **wraps**.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui -run 'TestTheScrollsBoundIsThePanesWidth' -v 2>&1 | tail -20
    ┃
    ┃  # github.com/lthiagol/maitre/internal/ui [github.com/lthiagol/maitre/internal/ui.test]
    ┃  internal/ui/bound_test.go:35:2: declared and not used: r
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui [build failed]
    ┃  FAIL
    ┃

    ┃
    ┃  ← Edit internal/ui/bound_test.go
    ┃
    ┃    31       "The chef reads the whole recipe in the console, so the detail   31       "The chef reads the whole recipe in the console, so the detail
    ┃         pane has to be "+                                                         pane has to be "+
    ┃    32         "as tall as the work it shows and never stop one row short of  32         "as tall as the work it shows and never stop one row short of
    ┃         the end. ",                                                               the end. ",
    ┃    33       4)                                                               33       4)
    ┃    34                                                                        34
    ┃    35 -   r, err := ops.CreateRecipe(core.RecipeSpec{                        35 +   if _, err := ops.CreateRecipe(core.RecipeSpec{
    ┃    36       Title:      "The detail pane scrolls to its last row",           36       Title:      "The detail pane scrolls to its last row",
    ┃    37       Outcome:    outcome,                                             37       Outcome:    outcome,
    ┃    38       Problem:    "The scroll's bound was counted at the wrong width.  38       Problem:    "The scroll's bound was counted at the wrong width.
    ┃         ",                                                                        ",
    ┃    39       OutOfScope: []string{"web shell", "orchestration"},              39       OutOfScope: []string{"web shell", "orchestration"},
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui -run 'TestTheScrollsBoundIsThePanesWidth' -v 2>&1 | tail -20
    ┃
    ┃  # github.com/lthiagol/maitre/internal/ui
    ┃  internal/ui/bound_test.go:46:2: expected operand, found 'if'
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui [setup failed]
    ┃  FAIL
    ┃

       → Explored: 1 read

    ┃
    ┃  ← Edit internal/ui/bound_test.go
    ┃
    ┃    41         {Description: "the scroll reaches the pane's last row",        41         {Description: "the scroll reaches the pane's last row",
    ┃         Verification: "go test ./internal/ui -run                                 Verification: "go test ./internal/ui -run
    ┃         TestTheScrollsBoundIsThePanesWidth"},                                     TestTheScrollsBoundIsThePanesWidth"},
    ┃    42       },                                                               42       },
    ┃    43       WorkPackages: []string{"ui"},                                    43       WorkPackages: []string{"ui"},
    ┃    44       Steps:        []core.StepSpec{{Action: "count the bound at the   44       Steps:        []core.StepSpec{{Action: "count the bound at the
    ┃         pane's width", DoneWhen: "the scroll reaches the end", Covers: []         pane's width", DoneWhen: "the scroll reaches the end", Covers: []
    ┃         string{"AC-01"}}},                                                        string{"AC-01"}}},
    ┃    45 -   })                                                                 45 +   }); err != nil {
    ┃    46 -   if err != nil {
    ┃    47       t.Fatal(err)                                                     46       t.Fatal(err)
    ┃    48     }                                                                  47     }
    ┃    49                                                                        48
    ┃    50     m := loadModel(t, store)                                           49     m := loadModel(t, store)
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui -run 'TestTheScrollsBoundIsThePanesWidth' -v 2>&1 | tail -20
    ┃
    ┃  === RUN   TestTheScrollsBoundIsThePanesWidth
    ┃  --- PASS: TestTheScrollsBoundIsThePanesWidth (0.18s)
    ┃  PASS
    ┃  ok    github.com/lthiagol/maitre/internal/ui  0.629s
    ┃

       + Thought · 711ms

    ┃
    ┃  ← Edit internal/ui/bound_test.go
    ┃
    ┃    73     if paneRows == fullRows {                                          73     if paneRows == fullRows {
    ┃    74       t.Fatalf("the fixture must make the two widths differ: %d rows   74       t.Fatalf("the fixture must make the two widths differ: %d rows
    ┃         at %d columns and %d at %d",                                              at %d columns and %d at %d",
    ┃    75         paneRows, paneWidth, fullRows, m.width-4)                      75         paneRows, paneWidth, fullRows, m.width-4)
    ┃    76     }                                                                  76     }
    ┃                                                                              77 +   t.Logf("the pane has %d rows at %d columns; the full screen has %d
    ┃                                                                                   at %d", paneRows, paneWidth, fullRows, m.width-4)
    ┃    77                                                                        78
    ┃    78     // The bound is the pane's rows.                                   79     // The bound is the pane's rows.
    ┃    79     if m.detailLineCount() != paneRows {                               80     if m.detailLineCount() != paneRows {
    ┃    80       t.Errorf("the bound is %d rows; the pane has %d rows at %d       81       t.Errorf("the bound is %d rows; the pane has %d rows at %d
    ┃         columns", m.detailLineCount(), paneRows, paneWidth)                       columns", m.detailLineCount(), paneRows, paneWidth)
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui -run 'TestTheScrollsBoundIsThePanesWidth' -v 2>&1 | tail -10
    ┃
    ┃  === RUN   TestTheScrollsBoundIsThePanesWidth
    ┃      bound_test.go:77: the pane has 29 rows at 107 columns; the full screen has 27 at 196
    ┃  --- PASS: TestTheScrollsBoundIsThePanesWidth (0.17s)
    ┃  PASS
    ┃  ok    github.com/lthiagol/maitre/internal/ui  0.427s
    ┃

       + Thought · 394ms

       Now the mutation — the bound counted at the full-screen width:

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/ui/model.go /tmp/t2a-final/model.go.keep && python3 - <<'PY'
    ┃  p="internal/ui/model.go"…
    ┃
    ┃  (2 earlier lines)
    ┃      bound_test.go:77: the pane has 29 rows at 107 columns; the full screen has 27 at 196
    ┃      bound_test.go:81: the bound is 27 rows; the pane has 29 rows at 107 columns
    ┃      bound_test.go:90: the scroll stopped at 26, want the pane's last row 28
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  86.733s
    ┃  FAIL
    ┃

       + Thought · 635ms

       → Explored: 1 read

    ┃
    ┃  ← Edit docs/05-human-interface.md
    ┃
    ┃    268                                                                       268
    ┃    269   | Gesture | Where | What |                                          269   | Gesture | Where | What |
    ┃    270   |---|---|---|                                                       270   |---|---|---|
    ┃    271   | wheel | over the detail pane | scroll the detail |                271   | wheel | over the detail pane | scroll the detail |
    ┃    272 - | wheel | anywhere else | scroll the list, keeping the selection    272 + | wheel | anywhere else, **while the list is on screen** | scroll
    ┃          visible |                                                                 the list, keeping the selection visible |
    ┃    273                                                                       273
    ┃    274   The list is the console's one scrollable thing, so a wheel over the 274   The list is the console's one scrollable thing, so a wheel over the
    ┃          header scrolls the                                                        header scrolls the
    ┃    275 - list rather than nothing. The **full-screen detail** has no list    275 + list rather than nothing. The **full-screen detail** has no list
    ┃          beside it, and there                                                      beside it: there the
    ┃    276 - the wheel scrolls the detail alone.                                 276 + wheel over the header does nothing at all, and the wheel over the
    ┃                                                                                    detail scrolls it.
    ┃    277                                                                       277
    ┃    278   The wheel's step is **three lines**: one is too slow to be useful,  278   The wheel's step is **three lines**: one is too slow to be useful,
    ┃          and a page is too                                                         and a page is too
    ┃    279   fast to aim with.                                                   279   fast to aim with.
    ┃    280                                                                       280
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && go build ./... && go test ./internal/ui 2>&1 | tail -5
    ┃
    ┃  M docs/05-human-interface.md
    ┃   M internal/ui/mouse_test.go
    ┃  ?? internal/ui/bound_test.go
    ┃  ok    github.com/lthiagol/maitre/internal/ui  86.873s
    ┃

       + Thought · 355ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/ui/bound_test.go internal/ui/mouse_test.go && git commit -q -m "test(ui): a tripwire for t
    ┃  he bound's width…
    ┃
    ┃  a332158 test(ui): a tripwire for the bound's width
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/05-human-interface.md && git commit -q -m "docs: the wheel's rule reads whole
    ┃  …
    ┃
    ┃  d8d7b95 docs: the wheel's rule reads whole
    ┃  a332158 test(ui): a tripwire for the bound's width
    ┃

       + Thought · 929ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && (make check-race > /tmp/t2a-race2.log 2>&1; echo "EXIT=$?" >> /tmp/t2a-race2.log) & sleep 3; echo started
    ┃
    ┃   Background
    ┃

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃  Could not find oldString in /Users/thiago/Code/Pessoal/maitre/work/tui/README.md. It must match exactly, including whitespace and indentation.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "T2a" work/tui/README.md | head -5
    ┃
    ┃  16:| **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | 🔄 **fix cycle 2** — 10 commits · a72ea98,
    ┃  115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 619b342, 47be828, f1cc1ca | [tui-t2a.md](../briefs/tui-t2a.md) · [fix 1](../briefs/fix-tui-t2a-
    ┃  1.md) · [fix 2](../briefs/fix-tui-t2a-2.md) |
    ┃  30:**The briefs that exist** are T1, T2a and T2b. The rest are written **just before**
    ┃

       + Thought · 829ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && git log --oneline -4
    ┃
    ┃  d8d7b95 docs: the wheel's rule reads whole
    ┃  a332158 test(ui): a tripwire for the bound's width
    ┃  fddaa0d work: T2a's fix cycle 2 - the bound's width, and one line of the document
    ┃  d939ab1 inbox: reviewer reports tui-t2a cycle 2
    ┃

       ↳ Shell finished · cd /Users/thiago/Code/Pessoal/maitre && (make check-race > /tmp/t2a-race2.log 2>&1; echo "EXIT=$?" >> /tmp/t2a-race2.log) & sleep…

       + Thought · 283ms

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃    12                                                                        12
    ┃    13   | Session | Feature | What it lands | State | Brief |                13   | Session | Feature | What it lands | State | Brief |
    ┃    14   |--:|---|---|---|---|                                                14   |--:|---|---|---|---|
    ┃    15   | **T1** | the identity | the window title carries the project; the  15   | **T1** | the identity | the window title carries the project; the
    ┃         header's name stops being a footnote; `--project`; the Settings tab       header's name stops being a footnote; `--project`; the Settings tab
    ┃         shows the plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438,        shows the plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438,
    ┃         f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81,            f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81,
    ┃         8a6866b, cf8c05e | [tui-t1.md](../briefs/tui-t1.md) · [fix 1](../         8a6866b, cf8c05e | [tui-t1.md](../briefs/tui-t1.md) · [fix 1](../
    ┃         briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |            briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |
    ┃    16 - | **T2a** | the mouse, part 1 | the mouse is on; the render builds   16 + | **T2a** | the mouse, part 1 | the mouse is on; the render builds
    ┃         the hit map; the wheel scrolls | 🔄 **fix cycle 2** — 10 commits ·        the hit map; the wheel scrolls | 🔄 **fix cycle 2 done** — 12
    ┃         a72ea98, 115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061,            commits · a72ea98, 115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7,
    ┃         619b342, 47be828, f1cc1ca | [tui-t2a.md](../briefs/tui-t2a.md) · [        e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](
    ┃         fix 1](../briefs/fix-tui-t2a-1.md) · [fix 2](../briefs/fix-tui-t2a-       ../briefs/tui-t2a.md) · [fix 1](../briefs/fix-tui-t2a-1.md) · [fix
    ┃         2.md) |                                                                   2](../briefs/fix-tui-t2a-1.md) |
    ┃    17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail;  17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail;
    ┃         the mouse drivable from a script | 🔄 **briefed** · next | [tui-t2b.      the mouse drivable from a script | 🔄 **briefed** · next | [tui-t2b.
    ┃         md](../briefs/tui-t2b.md) |                                               md](../briefs/tui-t2b.md) |
    ┃    18   | **T3a** | the backlog is readable | the detail pane; `space` stops 18   | **T3a** | the backlog is readable | the detail pane; `space` stops
    ┃         being a trap on an idea; D-3, D-5 | ⏳ planned | — |                      being a trap on an idea; D-3, D-5 | ⏳ planned | — |
    ┃    19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned 19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned
    ┃         | — |                                                                     | — |
    ┃    20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef    20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef
    ┃         dissolves; D-2 | ⏳ planned | — |                                         dissolves; D-2 | ⏳ planned | — |
    ┃

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃  Could not find oldString in /Users/thiago/Code/Pessoal/maitre/work/tui/README.md. It must match exactly, including whitespace and indentation.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "T2a\*\*" work/tui/README.md | head -3 && ls work/briefs/ | grep t2a
    ┃
    ┃  16:| **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | 🔄 **fix cycle 2 done** — 12 commits ·
    ┃  a72ea98, 115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](../briefs/tui-t2a.md) · [
    ┃  fix 1](../briefs/fix-tui-t2a-1.md) · [fix 2](../briefs/fix-tui-t2a-1.md) |
    ┃  fix-tui-t2a-1.md
    ┃  fix-tui-t2a-2.md
    ┃  review-tui-t2a-cycle2.md
    ┃  review-tui-t2a.md
    ┃  tui-t2a.md
    ┃

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃    12                                                                        12
    ┃    13   | Session | Feature | What it lands | State | Brief |                13   | Session | Feature | What it lands | State | Brief |
    ┃    14   |--:|---|---|---|---|                                                14   |--:|---|---|---|---|
    ┃    15   | **T1** | the identity | the window title carries the project; the  15   | **T1** | the identity | the window title carries the project; the
    ┃         header's name stops being a footnote; `--project`; the Settings tab       header's name stops being a footnote; `--project`; the Settings tab
    ┃         shows the plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438,        shows the plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438,
    ┃         f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81,            f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81,
    ┃         8a6866b, cf8c05e | [tui-t1.md](../briefs/tui-t1.md) · [fix 1](../         8a6866b, cf8c05e | [tui-t1.md](../briefs/tui-t1.md) · [fix 1](../
    ┃         briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |            briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |
    ┃    16 - | **T2a** | the mouse, part 1 | the mouse is on; the render builds   16 + | **T2a** | the mouse, part 1 | the mouse is on; the render builds
    ┃         the hit map; the wheel scrolls | 🔄 **fix cycle 2 done** — 12             the hit map; the wheel scrolls | 🔄 **fix cycle 2 done** — 12
    ┃         commits · a72ea98, 115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7,           commits · a72ea98, 115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7,
    ┃         e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](      e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](
    ┃         ../briefs/tui-t2a.md) · [fix 1](../briefs/fix-tui-t2a-1.md) · [fix        ../briefs/tui-t2a.md) · [fix 1](../briefs/fix-tui-t2a-1.md) · [fix
    ┃         2](../briefs/fix-tui-t2a-1.md) |                                          2](../briefs/fix-tui-t2a-2.md) |
    ┃    17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail;  17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail;
    ┃         the mouse drivable from a script | 🔄 **briefed** · next | [tui-t2b.      the mouse drivable from a script | 🔄 **briefed** · next | [tui-t2b.
    ┃         md](../briefs/tui-t2b.md) |                                               md](../briefs/tui-t2b.md) |
    ┃    18   | **T3a** | the backlog is readable | the detail pane; `space` stops 18   | **T3a** | the backlog is readable | the detail pane; `space` stops
    ┃         being a trap on an idea; D-3, D-5 | ⏳ planned | — |                      being a trap on an idea; D-3, D-5 | ⏳ planned | — |
    ┃    19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned 19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned
    ┃         | — |                                                                     | — |
    ┃    20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef    20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef
    ┃         dissolves; D-2 | ⏳ planned | — |                                         dissolves; D-2 | ⏳ planned | — |
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '16p' work/tui/README.md && git diff --stat
    ┃
    ┃  | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | 🔄 **fix cycle 2 done** — 12 commits · a72ea98,
    ┃  115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](../briefs/tui-t2a.md) · [fix 1](../
    ┃  briefs/fix-tui-t2a-1.md) · [fix 2](../briefs/fix-tui-t2a-2.md) |
    ┃   work/tui/README.md | 2 +-
    ┃   1 file changed, 1 insertion(+), 1 deletion(-)
    ┃

       + Thought · 956ms

    ┃
    ┃  $ sleep 240; tail -4 /tmp/t2a-race2.log
    ┃
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃  EXIT=0
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check 2>&1 | tail -8 && echo "--- status ---" && git status --short && echo "--- goldens ---" && git d
    ┃  iff HEAD --stat -- internal/ui/testdata && echo "(no golden diff = unchanged)"
    ┃
    ┃  (6 earlier lines)
    ┃
    ┃  ==> tidy-check
    ┃  --- status ---
    ┃   M work/tui/README.md
    ┃  --- goldens ---
    ┃  (no golden diff = unchanged)
    ┃

       + Thought · 1.7s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/ui/model.go /tmp/t2a-final/model.go.keep && python3 - <<'PY'
    ┃  p="internal/ui/model.go"…
    ┃
    ┃  (3 earlier lines)
    ┃      bound_test.go:81: the bound is 27 rows; the pane has 29 rows at 107 columns
    ┃      bound_test.go:90: the scroll stopped at 26, want the pane's last row 28
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  84.013s
    ┃  FAIL
    ┃   M work/tui/README.md
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && mkdir -p work/inbox && cat > work/inbox/runner-tui-t2a-fix2.md <<'EOF'
    ┃  Runner: T2a fix 2 done. Commits: 4 (a332158, d8d7b95, plus the state row and this receipt). Fixed: F-109 and the document's wheel table.…
    ┃
    ┃  22b9de6 inbox: runner reports tui-t2a fix 2
    ┃  d8d7b95 docs: the wheel's rule reads whole
    ┃  a332158 test(ui): a tripwire for the bound's width
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1H:p1" "Runner: T2a fix 2 done. Commits: 4 (a332158, d8d7b95, 22b9de6, plus the state
    ┃  row). Fixed: F-109 and the document's wheel table. The bound's width: the mutation — detailLineCount() counted at the full-screen width (m.wid…
    ┃
    ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":
    ┃  "ses_e737d6341ffeTdoEYTgykU25M8"},"agent_status":"done","completion_seq":151,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,
    ┃  "foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1H:p1","revision":15,"screen_detection_skipped":true,"state_change_seq":151,
    ┃  "tab_id":"w1H:t1","terminal_id":"term_65f0cd2aa89d26","terminal_title":"OC | Project work proposal","terminal_title_stripped":"OC | Project work
    ┃  proposal","workspace_id":"w1H"},"type":"agent_prompted"}}
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && git log --oneline -5
    ┃
    ┃  22b9de6 inbox: runner reports tui-t2a fix 2
    ┃  d8d7b95 docs: the wheel's rule reads whole
    ┃  a332158 test(ui): a tripwire for the bound's width
    ┃  fddaa0d work: T2a's fix cycle 2 - the bound's width, and one line of the document
    ┃  d939ab1 inbox: reviewer reports tui-t2a cycle 2
    ┃

       T2a's fix cycle 2 is landed. The receipt is committed and the nudge was delivered to w1H:p1.

       F-109 — the bound's width

       The old test was vacuous twice over: both sides of its comparison went through detailPaneWidth(), and the fixture's detail wrapped to the same row
       count at both widths — so counting the bound at the full-screen width changed nothing the test could see.

       The new tripwire (internal/ui/bound_test.go):

       - the pane's width is computed in the test — m.width - m.listColumn() - 3 — and the pane's rows are counted there, not through the console's answer;
       - the fixture is a new plan whose recipe has a long outcome, and the test fails outright if the two widths ever agree again.

       The mutation fails, and it fails for the right reason. Counting the bound at the full-screen width again turns internal/ui red:

       --- FAIL: TestTheScrollsBoundIsThePanesWidth (0.32s)
           bound_test.go:81: the bound is 27 rows; the pane has 29 rows at 107 columns
           bound_test.go:90: the scroll stopped at 26, want the pane's last row 28

       Nothing else in the suite failed — the failure names the bound's row count against the pane's, and the scroll stopping two rows short of the pane's
       last row.

       The fixture's margin, printed by the test itself: the split pane wraps the outcome into 29 rows at 107 columns, the full screen into 27 at 196 — two
       rows apart. The scroll reaches the pane's last row (28 of 29).

       The document

       The table's second row was "anywhere else", which was not true of the full-screen detail. It now reads:

       ┌────────────────────────┬───────────────────────────────────────────────────────────┬──────────────────────────────────────────────────────────────┐
       │ Gesture                │ Where                                                     │ What                                                         │
       ├────────────────────────┼───────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────┤
       │ wheel                  │ over the detail pane                                      │ scroll the detail                                            │
       ├────────────────────────┼───────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────┤
       │ wheel                  │ anywhere else, while the list is on screen                │ scroll the list, keeping the selection visible               │
       └────────────────────────┴───────────────────────────────────────────────────────────┴──────────────────────────────────────────────────────────────┘

       and the prose says what the wheel does with the full-screen detail open: the wheel over the header does nothing at all, and the wheel over the detail
       scrolls it.

       The wheel's fallback is untouched — one rule in one place — and the review's two notes are left for a later round, as it ruled them.

       make check green · make check-race green · goldens unchanged · tree clean · DomainVersion 7 · no dependency added.

       Build · MiniMax-M3.1-Flash · 13m 11s · 80.3 tok/s

    ┃
    ┃
    ┃
    ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
    ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
    ~/Code/Pessoal/maitre:master                                                                                               232.7K (23%)  ctrl+p commands
```
