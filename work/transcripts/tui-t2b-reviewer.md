# The reviewer's lane — T2b, captured 2026-10-30

Dumped from pane `w1H:p7` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that judged T2b: one cycle, no blockers, and five notes — two of them about
the slice's own tests reading stronger than they assert.
The lane that wrote it is in [tui-t2b-runner.md](tui-t2b-runner.md).

    ┃    rule under "Needs you" and the frame did not change.
    ┃
    ┃  **The duplicated doc comment** in `internal/ui/mouse.go` is in the brief as a note.
    ┃



       → Explored: 2 reads

       + Thought · 196ms

       → Explored: 1 read

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -15 && git status --short && ls internal/ui/
    ┃
    ┃  (45 earlier lines)
    ┃  term.go
    ┃  testdata
    ┃  theme_test.go
    ┃  theme.go
    ┃  triage_test.go
    ┃  views.go
    ┃  words_test.go
    ┃

       + Thought · 171ms

       → Explored: 2 reads

       + Thought · 519ms

       → Explored: 1 read

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "func (m \*Model) switchTab\|func (m \*Model) moveTo\|func (m \*Model) openDetail\|func (m \*Mode
    ┃  l) clampSelections\|func (m \*Model) move\b\|func (m \*Model) selectByIndex" internal/ui/*.go
    ┃
    ┃  internal/ui/model.go:485:func (m *Model) clampSelections() {
    ┃  internal/ui/model.go:631:func (m *Model) switchTab(t Tab) {
    ┃  internal/ui/model.go:644:func (m *Model) openDetail() {
    ┃  internal/ui/model.go:697:func (m *Model) move(delta int) {
    ┃  internal/ui/model.go:718:func (m *Model) moveTo(index int) {
    ┃

       → Explored: 1 read

       + Thought · 3.9s

       → Explored: 1 read

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "MouseClickMsg\|MouseReleaseMsg\|MousePressMsg" $(go env GOMODCACHE)/charm.land/bubbletea/v2*/mo
    ┃  use.go 2>/dev/null | head -50; ls $(go env GOMODCACHE)/charm.land/ 2>/dev/null
    ┃
    ┃  (11 earlier lines)
    ┃  bubbletea
    ┃  bubbletea@v1.3.10
    ┃  glamour
    ┃  glamour@v1.0.0
    ┃  lipgloss
    ┃  lipgloss@v1.1.0
    ┃

       + Thought · 420ms

    ┃
    ┃  $ cd /Users/thiago/go/pkg/mod/charm.land/bubbletea/v2@v2.1.0 && grep -rn "MouseClickMsg\|MouseReleaseMsg\|MouseMotionMsg" *.go | grep -v "_test" |
    ┃  head -40
    ┃
    ┃  (11 earlier lines)
    ┃  mouse.go:108:func (e MouseReleaseMsg) Mouse() Mouse {
    ┃  mouse.go:127:// MouseMotionMsg represents a mouse motion message.
    ┃  mouse.go:128:type MouseMotionMsg Mouse
    ┃  mouse.go:131:func (e MouseMotionMsg) String() string {
    ┃  mouse.go:142:func (e MouseMotionMsg) Mouse() Mouse {
    ┃  tea.go:826:        case MouseClickMsg, MouseReleaseMsg, MouseWheelMsg, MouseMotionMsg:
    ┃

       → Explored: 1 read

    ┃
    ┃  $ cd /Users/thiago/go/pkg/mod/charm.land/bubbletea/v2@v2.1.0 && grep -rn "MousePress\|mousePress\|X10\|SGR" key.go mouse.go | head -40
    ┃
    ┃  (no output)
    ┃

       + Thought · 2.7s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "hits\|hitMapOf\|listColumn\|detailPaneWidth\|MouseAllMotion\|EnableMouse" internal/ui/*.go cmd/m
    ┃  aitre/*.go | grep -v _test
    ┃
    ┃  (38 earlier lines)
    ┃  internal/ui/views.go:739:    out = append(out, listLine{text: m.styRule.Render(strings.Repeat("─", m.listColumn())), itemIdx: -1})
    ┃  internal/ui/views.go:752:  width := m.listColumn()
    ┃  internal/ui/views.go:806:  width := m.listColumn()
    ┃  internal/ui/views.go:827:  width := m.listColumn()
    ┃  internal/ui/views.go:899:  width := m.listColumn()
    ┃  internal/ui/views.go:998:  art := markFor(m.listColumn())
    ┃

       → Explored: 1 read

       + Thought · 1.1s

       → Explored: 1 read

       + Thought · 976ms

       → Explored: 1 read

       + Thought · 748ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "func (m \*Model) listLines\|func (m \*Model) logLines\|func (m \*Model) settingsLines" -A 30 int
    ┃  ernal/ui/views.go | sed -n '1,120p'
    ┃
    ┃  (89 earlier lines)
    ┃  1091-  }
    ┃  1092-  add("")
    ┃  1093-  add(" " + m.stySec.Render("KEYS"))
    ┃  1094-  enterWhat := "open the recipe full screen"
    ┃  1095-  if m.tab == TabSettings {
    ┃  1096-    enterWhat = "apply the theme"
    ┃

       + Thought · 783ms

       → Explored: 1 read

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "func (m \*Model) tabBar" -A 40 internal/ui/views.go
    ┃
    ┃  (34 earlier lines)
    ┃  308-  style := m.styAcc
    ┃  309-  if m.refusal != "" {
    ┃  310-    style = m.styWarn
    ┃  311-  }
    ┃  312-  return " " + left + strings.Repeat(" ", gap) + style.Render(right)
    ┃  313-}
    ┃  314-
    ┃

       + Thought · 388ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "func (m \*Model) Update" -A 60 internal/ui/model.go
    ┃
    ┃  (54 earlier lines)
    ┃  321-  }
    ┃  322-}
    ┃  323-
    ┃  324-// sectionsByLane groups by the lifecycle lane.
    ┃  325-func (m *Model) sectionsByLane() []core.TriageSection {
    ┃  326-  byLane := map[core.Lane][]core.Recipe{}
    ┃  327-  for _, r := range m.recipes {
    ┃

       + Thought · 561ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && ls cmd/maitre/ && grep -rn "SplitSteps\|keys" cmd/maitre/*.go | head -40
    ┃
    ┃  main.go
    ┃  cmd/maitre/main.go:5:// terminal and works in CI. `--keys` drives the console from a script, so its behaviour
    ┃  cmd/maitre/main.go:35:  keys := fs.String("keys", "", "drive the console with these steps before rendering, comma separated — keys and mouse
    ┃  steps")
    ┃  cmd/maitre/main.go:59:  if *keys != "" || *all || *show != "" || *render {
    ┃  cmd/maitre/main.go:74:    if *keys != "" {
    ┃  cmd/maitre/main.go:75:      out, err := ui.Drive(m, ui.SplitSteps(*keys))
    ┃

       → Explored: 1 read

       + Thought · 390ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make build 2>&1 | tail -5 && git log --oneline -6 --stat | head -60
    ┃
    ┃  (16 earlier lines)
    ┃   internal/ui/click_test.go | 108 ++++++++++++++++++++++++++++++++++++++++++++++
    ┃   internal/ui/mouse.go      |   7 +++
    ┃   2 files changed, 115 insertions(+)
    ┃  ac8df77 feat(ui): a click on the detail pane is enter
    ┃   internal/ui/click_test.go | 36 ++++++++++++++++++++++++++++++++++++
    ┃   internal/ui/mouse.go      |  3 +++
    ┃   2 files changed, 39 insertions(+)
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "Interaction" -A 60 docs/05-human-interface.md | head -100
    ┃
    ┃  (54 earlier lines)
    ┃  305-
    ┃  306-**The theme is the chef's to pick.** The Settings tab lists the palettes under a
    ┃  307-**THEME** heading; the cursor moves over them and `enter` applies the one under it. The
    ┃  308-active theme is marked, so the chef can see which one is in force. The choice is for the
    ┃  309-session: nothing is written to the plan, so a theme is a setting of the console, not of
    ┃  310-the work.
    ┃  311-
    ┃

       + Thought · 226ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git show 9c716ce -- docs/05-human-interface.md
    ┃
    ┃  (58 earlier lines)
    ┃   ```
    ┃   maitre --keys "wheel:down"         # the wheel, over the list
    ┃   maitre --keys "wheel:down@(12,5)"   # the wheel, at column 12, row 5
    ┃  +maitre --keys "click@(12,2)"        # the click, at column 12, row 2
    ┃   ```
    ┃
    ┃   The coordinates are the terminal's own — zero-based, with (0, 0) at the upper-left
    ┃

       + Thought · 1.0s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && rm -rf /tmp/kb2 && ./bin/maitre-cli demo --dir /tmp/kb2 && cd /tmp/kb2 && /Users/thiago/Code/Pessoal/mait
    ┃  re/bin/maitre --render --plain --width 120 | head -8; echo "=== exit $?"
    ┃
    ┃  (8 earlier lines)
    ┃   7 recipes · 2 need you   ☑ show served
    ┃   Needs you                                                                                                          2
    ┃  ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
    ┃    RCP02 ready       Sign off the menu                             1/2
    ┃    RCP03 remediation Wire the review queue                         2/3 1 note
    ┃  === exit 0
    ┃

       + Thought · 857ms

    ┃
    ┃  $ cd /tmp/kb2 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre; $M --render --plain --width 120 --keys "click@(14,2)" | head -8; echo "=== exit ${
    ┃  PIPESTATUS[0]}"; echo "--- split width 200, click a row:"; $M --render --plain --width 200 --keys "click@(5,7)" | head -12
    ┃
    ┃  (16 earlier lines)
    ┃    RCP02 ready       Sign off the menu                             1/2                       → awaiting tasting
    ┃    RCP03 remediation Wire the review queue                         2/3 1 note
    ┃   On the line                                                                            2  Outcome
    ┃  ──────────────────────────────────────────────────────────────────────────────────────────
    ┃    RCP04 approved    Add the backlog view                          1/1 ●                        Findings reach the chef without leaving the board
    ┃    RCP05 executing   Sharpen the evidence contract                 1/2
    ┃

       + Thought · 467ms

    ┃
    ┃  $ cd /tmp/kb2 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre; echo "--- narrow click row 7:"; $M --render --plain --width 120 --keys "click@(5,7
    ┃  )" | head -14; echo "--- exit ${PIPESTATUS[0]}"
    ┃
    ┃  (10 earlier lines)
    ┃
    ┃      Findings reach the chef without leaving the board
    ┃
    ┃  Problem
    ┃
    ┃  --- exit
    ┃

       + Thought · 360ms

    ┃
    ┃  $ cd /tmp/kb2 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre; echo "--- backlog click (narrow) row 6:"; $M --render --plain --width 120 --keys "
    ┃  click@(14,2),click@(5,6)" | head -10; echo "--- settings click (narrow):"; $M --render --plain --width 120 --keys "click@(60,2)" | head -20
    ┃
    ┃  (26 earlier lines)
    ┃
    ┃   THEME
    ┃   ● Kitchen
    ┃     Catppuccin Latte
    ┃     Catppuccin Frappé
    ┃     Catppuccin Macchiato
    ┃

       + Thought · 2.4s

    ┃
    ┃  $ cd /tmp/kb2 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre; $M --render --width 120 --keys "click@(60,2)" > /tmp/base.ansi; $M --render --widt
    ┃  h 120 --keys "click@(60,2),click@(5,17)" > /tmp/click.ansi; diff <(cat -v /tmp/base.ansi) <(cat -v /tmp/click.ansi) | head -20; echo "diff e…
    ┃
    ┃  17,18c17,18
    ┃  <  ^[[38;5;252;48;5;237m�M-^WM-^O Kitchen               ^[[m
    ┃  <    Catppuccin Latte
    ┃  ---
    ┃  >  �M-^WM-^O Kitchen
    ┃  >  ^[[38;5;252;48;5;237m  Catppuccin Latte      ^[[m
    ┃  diff exit 0
    ┃

       + Thought · 5.8s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check 2>&1 | tail -30
    ┃
    ┃  (23 earlier lines)
    ┃  GO=go scripts/dep-budget.sh
    ┃    internal/core: 0 external modules (ceiling 0)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃

       + Thought · 517ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && echo "--- diff:" && git diff --stat && echo "--- DomainVersion:" && grep -rn "Domai
    ┃  nVersion =" internal/core/ && echo "--- tabs:" && grep -rn "TabNames = \|TabNames=" internal/ui/*.go | head -3 && echo "--- commits touching…
    ┃
    ┃  (3 earlier lines)
    ┃  internal/core/activity.go:154:    e.DomainVersion = PlanDomainVersion(j)
    ┃  internal/core/version_test.go:356:  j.Events[len(j.Events)-1].DomainVersion = version
    ┃  internal/core/lifecycle.go:225:const DomainVersion = 7
    ┃  --- tabs:
    ┃  internal/ui/model.go:32:var TabNames = []string{"Recipes", "Backlog", "Path", "Board", "Log", "Sous-Chef", "Settings"}
    ┃  --- commits touching core since T2a:
    ┃

       + Thought · 616ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff --stat eeedff0..HEAD && echo "=== mouse.go diff:" && git diff eeedff0..HEAD -- internal/ui/mouse
    ┃  .go
    ┃
    ┃  (135 earlier lines)
    ┃  +      return nil, false
    ┃  +    }
    ┃  +    return tea.MouseClickMsg{X: x, Y: y, Button: tea.MouseLeft}, true
    ┃     }
    ┃     return nil, false
    ┃   }
    ┃

       + Thought · 1.9s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff eeedff0..HEAD -- internal/ui/model.go
    ┃
    ┃  (76 earlier lines)
    ┃  +  m.mode = ModeDetail
    ┃  +  m.detailScroll = 0
    ┃  +}
    ┃  +
    ┃   // scrollDetail moves the detail pane's window.
    ┃   //
    ┃   // A recipe is longer than the pane — an intent, evidence per criterion, a review
    ┃

       + Thought · 5.5s



    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'EOF'
    ┃  import re…
    ┃
    ┃  (3 earlier lines)
    ┃      click_test.go:412: the click on a tab left the status "signed off RCP02", want it cleared
    ┃      click_test.go:415: the click on a tab left the cursor at 99, want it inside the 6 rows
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  12.586s
    ┃  FAIL
    ┃  make: *** [test] Error 1
    ┃

       + Thought · 599ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git checkout -- internal/ui/mouse.go && git status --short && python3 - <<'EOF'
    ┃  p='internal/ui/hit.go'…
    ┃
    ┃  ==> test
    ┃  GO=go scripts/go-test.sh " -run Click" "./internal/ui"
    ┃  ok    github.com/lthiagol/maitre/internal/ui  12.197s
    ┃

       + Thought · 1.8s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff -- internal/ui/hit.go && make test PKG=./internal/ui RUN='-run TestAClickOnARowSelectsIt' 2>&1 |
    ┃   tail -20
    ┃
    ┃  (11 earlier lines)
    ┃       hits.rows = append(hits.rows, hitRow{
    ┃         hitBox: newHitBox(0, width-1, start+i, start+i),
    ┃
    ┃  ==> test
    ┃  GO=go scripts/go-test.sh " -run TestAClickOnARowSelectsIt" "./internal/ui"
    ┃  ok    github.com/lthiagol/maitre/internal/ui  0.904s
    ┃

       + Thought · 1.8s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cat > internal/ui/zz_tmp_debug_test.go <<'EOF'
    ┃  package ui…
    ┃
    ┃  ==> test
    ┃  GO=go scripts/go-test.sh " -run TestZZDebugRows" "./internal/ui"
    ┃  ok    github.com/lthiagol/maitre/internal/ui  0.995s
    ┃

       + Thought · 458ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui -run TestZZDebugRows -v 2>&1 | head -50
    ┃
    ┃  (23 earlier lines)
    ┃      zz_tmp_debug_test.go:11: row 22 y=26 item=-1
    ┃      zz_tmp_debug_test.go:11: row 23 y=27 item=-1
    ┃      zz_tmp_debug_test.go:11: row 24 y=28 item=-1
    ┃      zz_tmp_debug_test.go:13: visible rows: 7
    ┃  --- PASS: TestZZDebugRows (0.64s)
    ┃  PASS
    ┃  ok    github.com/lthiagol/maitre/internal/ui  0.862s
    ┃

       + Thought · 9.2s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make test PKG=./internal/ui 2>&1 | tail -25
    ┃
    ┃  (18 earlier lines)
    ┃      mouse_test.go:155: the settings row 22 holds item 22, want -1
    ┃      mouse_test.go:155: the settings row 23 holds item 23, want -1
    ┃      mouse_test.go:155: the settings row 24 holds item 24, want -1
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  112.924s
    ┃  FAIL
    ┃  make: *** [test] Error 1
    ┃

       + Thought · 12.7s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && rm internal/ui/zz_tmp_debug_test.go && git checkout -- internal/ui/hit.go && git status --short && python
    ┃  3 - <<'EOF'…
    ┃
    ┃  Traceback (most recent call last):
    ┃    File "<stdin>", line 4, in <module>
    ┃  AssertionError
    ┃
    ┃  ==> test
    ┃  GO=go scripts/go-test.sh " -run Click" "./internal/ui"
    ┃  ok    github.com/lthiagol/maitre/internal/ui  15.676s
    ┃

       + Thought · 1.2s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'EOF'
    ┃  p='internal/ui/hit.go'…
    ┃
    ┃  (8 earlier lines)
    ┃         return r.item, true
    ┃       }
    ┃
    ┃  ==> test
    ┃  GO=go scripts/go-test.sh " -run Click" "./internal/ui"
    ┃  ok    github.com/lthiagol/maitre/internal/ui  14.751s
    ┃

       + Thought · 317ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make test PKG=./internal/ui 2>&1 | tail -15
    ┃
    ┃  (8 earlier lines)
    ┃      mouse_test.go:102: the point on row 27 holds item 0, want -1
    ┃      mouse_test.go:99: the point on row 28 holds item 0 (true), want an item: false
    ┃      mouse_test.go:102: the point on row 28 holds item 0, want -1
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  107.460s
    ┃  FAIL
    ┃  make: *** [test] Error 1
    ┃

       + Thought · 3.0s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git checkout -- internal/ui/hit.go && git status --short && python3 - <<'EOF'
    ┃  p='internal/ui/mouse.go'…
    ┃
    ┃  (10 earlier lines)
    ┃  GO=go scripts/go-test.sh " -run Click" "./internal/ui"
    ┃  # github.com/lthiagol/maitre/internal/ui [github.com/lthiagol/maitre/internal/ui.test]
    ┃  internal/ui/mouse.go:60:5: declared and not used: item
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui [build failed]
    ┃  FAIL
    ┃  make: *** [test] Error 1
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'EOF'
    ┃  p='internal/ui/mouse.go'…
    ┃
    ┃  (9 earlier lines)
    ┃  --- FAIL: TestTheScriptedClickTakesAPosition (0.82s)
    ┃      click_test.go:299: the scripted click left the cursor on 0, want 1
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  15.494s
    ┃  FAIL
    ┃  make: *** [test] Error 1
    ┃

       + Thought · 851ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git checkout -- internal/ui/mouse.go && python3 - <<'EOF'
    ┃  p='internal/ui/mouse.go'…
    ┃
    ┃  (13 earlier lines)
    ┃  --- FAIL: TestAClickOnARowInTheNarrowLayoutOpensTheDetail (0.89s)
    ┃      click_test.go:126: the click in the narrow layout left the console in 0, want the detail on screen
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  14.065s
    ┃  FAIL
    ┃  make: *** [test] Error 1
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git checkout -- internal/ui/mouse.go && python3 - <<'EOF'
    ┃  p='internal/ui/mouse.go'…
    ┃
    ┃  (11 earlier lines)
    ┃  --- FAIL: TestAClickOnTheDetailPaneOpensIt (0.78s)
    ┃      click_test.go:236: the click on the detail left the console in 0, want the detail on screen
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  13.859s
    ┃  FAIL
    ┃  make: *** [test] Error 1
    ┃

       + Thought · 555ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git checkout -- internal/ui/mouse.go && python3 - <<'EOF'
    ┃  p='internal/ui/mouse.go'…
    ┃
    ┃  (17 earlier lines)
    ┃  --- FAIL: TestTheClickSurvivesTheSplitter (0.00s)
    ┃      click_test.go:311: the steps are ["click@(12" "5)" "tab"], want ["click@(12,5)" "tab"]
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  13.728s
    ┃  FAIL
    ┃  make: *** [test] Error 1
    ┃

       + Thought · 1.2s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git checkout -- internal/ui/mouse.go internal/ui/hit.go && git status --short && git diff --stat && echo
    ┃  "CLEAN" && git log --oneline -1
    ┃
    ┃  CLEAN
    ┃  8fa6cbd work: T2b's review brief - the clicks are the keyboard's moves
    ┃

    ┃
    ┃  $ cd /tmp/kb2 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre; echo "--- unknown step:"; $M --render --keys "click:sideways"; echo "exit $?"; ech
    ┃  o "--- malformed:"; $M --render --keys "click@(1,2"; echo "exit $?"; echo "--- unknown verb:"; $M --render --keys "wibble"; echo "exit $?"; …
    ┃
    ┃  (6 earlier lines)
    ┃  --- unknown verb:
    ┃  maitre: unknown step "wibble"
    ┃  exit 1
    ┃  --- unbracketed:
    ┃  maitre: unknown step "click@1"
    ┃  exit 1
    ┃

       + Thought · 4.7s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '966,1010p' internal/ui/views.go
    ┃
    ┃  (38 earlier lines)
    ┃    }
    ┃    return art
    ┃  }
    ┃
    ┃  // sousChefContent is the tab's content, without the mark.
    ┃  //
    ┃  // The Sous-Chef is the orchestrator agent: it calls the tickets, routes them, and
    ┃

       + Thought · 3.6s

       → Explored: 1 read

       + Thought · 1.1s

       → Explored: 1 read

       + Thought · 5.3s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check-race 2>&1 | tail -20
    ┃
    ┃  (13 earlier lines)
    ┃  GO=go scripts/dep-budget.sh
    ┃    internal/core: 0 external modules (ceiling 0)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && echo "TREE-CLEAN" && rm -rf /tmp/kb3 && ./bin/maitre-cli demo --chaos --dir /tmp/kb
    ┃  3 >/dev/null && cd /tmp/kb3 && /Users/thiago/Code/Pessoal/maitre/bin/maitre --render --plain --width 120 --keys "click@(14,2)" | head -8; ec…
    ┃
    ┃  (4 earlier lines)
    ┃   9 items
    ┃   IDEAS                                                                                                              5
    ┃  ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
    ┃    BL01  …S: a backlog title with a newline and a tab idea    open
    ┃    BL04  The board needs a way to show evidence       idea    open
    ┃  exit 0
    ┃

       + Thought · 937ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && ./bin/maitre --render --plain --width 120 --keys "click@(14,2)" | head -6; echo "exit $?"; git status --s
    ┃  hort; echo CLEAN
    ┃
    ┃  (2 earlier lines)
    ┃    Recipes  ▸ Backlog    Path    Board    Log    Sous-Chef    Settings
    ┃   0 items
    ┃
    ┃
    ┃  exit 0
    ┃  CLEAN
    ┃

       + Thought · 4.5s

    ┃
    ┃  $ cd /tmp/kb2 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre; $M --render --plain --width 120 --keys "click@(5,7),click@(14,2)" | head -12; echo
    ┃   "=== keyboard equivalent:"; $M --render --plain --width 120 --keys "enter,tab" | head -12
    ┃
    ┃  (19 earlier lines)
    ┃
    ┃  Why
    ┃
    ┃      Evidence the chef cannot read is a claim, not evidence.
    ┃
    ┃  Tags
    ┃

       + Thought · 1.4s

       → Explored: 1 read

       + Thought · 3.3s

    ┃
    ┃  $ cd /tmp/kb2 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre; $M --render --plain --width 200 --keys "click@(14,2)" | head -12
    ┃
    ┃  (5 earlier lines)
    ┃  ──────────────────────────────────────────────────────────────────────────────────────────  idea · open
    ┃    BL01  The board needs a way to show evidence       idea    open
    ┃    BL02  A recipe with no ACs is probably an idea     idea    open                          Why
    ┃    BL03  The chef wants a weekly digest               idea    open
    ┃    BL04  Trying the board on a real project           idea    open                              Evidence the chef cannot read is a claim, not
    ┃  evidence.
    ┃   SMALL TASKS                                                                            2
    ┃  ────────────────────────────────────────────────────────────────────────────────────────── Tags
    ┃

       + Thought · 6.4s

       Now the ledger. First the cycle-matrix row:

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃     65   | Slice T2a | The mouse, part 1: the mouse is on; the render        65   | Slice T2a | The mouse, part 1: the mouse is on; the render
    ┃          builds the hit map; the wheel scrolls the list and the detail |          builds the hit map; the wheel scrolls the list and the detail |
    ┃          The mouse is on **every** view the console draws — the console's,        The mouse is on **every** view the console draws — the console's,
    ┃          the error's and the opening one's — and the error and opening            the error's and the opening one's — and the error and opening
    ┃          views **clear the map**, so no event can land on a tab that is not       views **clear the map**, so no event can land on a tab that is not
    ┃          on screen. The map is built **by the render**: `tabBar` emits the        on screen. The map is built **by the render**: `tabBar` emits the
    ┃          tab boxes as it draws them, `hitMapOf` is called from the three          tab boxes as it draws them, `hitMapOf` is called from the three
    ┃          render functions and from nowhere else, the banner pushes the            render functions and from nowhere else, the banner pushes the
    ┃          list's rows down in the map as it does on the screen, and the            list's rows down in the map as it does on the screen, and the
    ┃          geometry constants (`headerRows`, `footerRows`, `splitThreshold`)        geometry constants (`headerRows`, `footerRows`, `splitThreshold`)
    ┃          are read where they are written. The wheel over the list drives          are read where they are written. The wheel over the list drives
    ┃          `m.move` — the keyboard's own move — and a scripted `wheel:down`         `m.move` — the keyboard's own move — and a scripted `wheel:down`
    ┃          and three `down`s land on **byte-identical** frames. The wheel           and three `down`s land on **byte-identical** frames. The wheel
    ┃          over the detail drives `m.scrollDetail`, but **the split layout          over the detail drives `m.scrollDetail`, but **the split layout
    ┃          never renders that offset**: at ≥160 columns — where the console         never renders that offset**: at ≥160 columns — where the console
    ┃          draws the list and the detail side by side, and `ModeSplit` is the       draws the list and the detail side by side, and `ModeSplit` is the
    ┃          default — the pane does not move, and the clamp measures the             default — the pane does not move, and the clamp measures the
    ┃          detail at a different width than the pane draws it. The scripted         detail at a different width than the pane draws it. The scripted
    ┃          mouse's **position form** cannot be reached through `maitre --           mouse's **position form** cannot be reached through `maitre --
    ┃          keys`, and the doc prints the command line that fails. The scope         keys`, and the doc prints the command line that fails. The scope
    ┃          held: **no click exists**, on purpose; `internal/core` and               held: **no click exists**, on purpose; `internal/core` and
    ┃          `internal/cli` are untouched; `core.DomainVersion` is still **7**;       `internal/cli` are untouched; `core.DomainVersion` is still **7**;
    ┃          the budget is still **0 / 0 / 28**; the seven tabs kept their            the budget is still **0 / 0 / 28**; the seven tabs kept their
    ┃          names; the goldens did not move — **two blockers filed** (F-106,         names; the goldens did not move — **two blockers filed** (F-106,
    ┃          F-107) |                                                                 F-107) |
    ┃     66   | Slice T2a · cycle 1 | The review of the slice: the map, the       66   | Slice T2a · cycle 1 | The review of the slice: the map, the
    ┃          wheel, and six mutations | The slice's shape held: the map **is**        wheel, and six mutations | The slice's shape held: the map **is**
    ┃          built by the render (`tabBar` emits the tab boxes, `hitMapOf` is         built by the render (`tabBar` emits the tab boxes, `hitMapOf` is
    ┃          called from the three render functions and nowhere else, the             called from the three render functions and nowhere else, the
    ┃          banner pushes the list's rows down in the map, and the geometry          banner pushes the list's rows down in the map, and the geometry
    ┃          constants are read where they are written), the wheel **is** the         constants are read where they are written), the wheel **is** the
    ┃          keyboard's scroll (a scripted `wheel:down` and three `down`s land        keyboard's scroll (a scripted `wheel:down` and three `down`s land
    ┃          on byte-identical frames), and all **six** mutations failed the          on byte-identical frames), and all **six** mutations failed the
    ┃          suite — the mouse mode off the error view, the render not emitting       suite — the mouse mode off the error view, the render not emitting
    ┃          the map, a second scroll variable, the wheel over the header, the        the map, a second scroll variable, the wheel over the header, the
    ┃          wheel under an open prompt, and the scripted step acting without         wheel under an open prompt, and the scripted step acting without
    ┃          drawing first. Two blockers filed: **F-106** (the wheel over the         drawing first. Two blockers filed: **F-106** (the wheel over the
    ┃          detail does nothing in the split layout) and **F-107** (the              detail does nothing in the split layout) and **F-107** (the
    ┃          scripted position form is unreachable through `maitre --keys`, and       scripted position form is unreachable through `maitre --keys`, and
    ┃          the doc prints the command that fails). The zero box's one cell at       the doc prints the command that fails). The zero box's one cell at
    ┃          the origin was noted as **F-108** and ruled into the fix cycle.          the origin was noted as **F-108** and ruled into the fix cycle.
    ┃          The goldens did not move; `internal/core` was untouched; `core.          The goldens did not move; `internal/core` was untouched; `core.
    ┃          DomainVersion` was still **7**; the budget was still **0 / 0 /           DomainVersion` was still **7**; the budget was still **0 / 0 /
    ┃          28** — **verdict: changes-needed** |                                     28** — **verdict: changes-needed** |
    ┃     67   | Slice T2a · cycle 2 | The review of the fix cycle: the split      67   | Slice T2a · cycle 2 | The review of the fix cycle: the split
    ┃          detail, the position form, the absent region, and the **reversed**       detail, the position form, the absent region, and the **reversed**
    ┃          test | F-106, F-107 and F-108 are closed, each proved by a               test | F-106, F-107 and F-108 are closed, each proved by a
    ┃          mutation: taking the offset out of the split layout fails                mutation: taking the offset out of the split layout fails
    ┃          `TestTheWheelScrollsTheSplitDetail`; a splitter that cuts a              `TestTheWheelScrollsTheSplitDetail`; a splitter that cuts a
    ┃          bracketed position at its comma fails the two position tests;            bracketed position at its comma fails the two position tests;
    ┃          taking the `there` guard out of `hitBox` fails three tests. The          taking the `there` guard out of `hitBox` fails three tests. The
    ┃          offset is applied in **one** place and the pane's width has              offset is applied in **one** place and the pane's width has
    ┃          **one** answer (`detailPaneWidth`), and the wheel at the origin          **one** answer (`detailPaneWidth`), and the wheel at the origin
    ┃          over the list scrolls the list. **But the bound's guard does not         over the list scrolls the list. **But the bound's guard does not
    ┃          bite**: the bound put back on the full-screen width leaves the           bite**: the bound put back on the full-screen width leaves the
    ┃          **whole `internal/ui` suite green** — **one blocker filed** (F-          **whole `internal/ui` suite green** — **one blocker filed** (F-
    ┃          109). The **reversed** test is judged **sound**: the rule is one         109). The **reversed** test is judged **sound**: the rule is one
    ┃          rule — the detail answers over the detail, the list answers the          rule — the detail answers over the detail, the list answers the
    ┃          rest — stated in one place in the code and in the doc; the full-         rest — stated in one place in the code and in the doc; the full-
    ┃          screen detail keeps the *no list, nothing happens* answer, pinned        screen detail keeps the *no list, nothing happens* answer, pinned
    ┃          by `TestTheWheelOverTheFullScreenHeaderDoesNothing`; and the old         by `TestTheWheelOverTheFullScreenHeaderDoesNothing`; and the old
    ┃          rule and the fix cycle's own F-108 evidence (*a wheel at the             rule and the fix cycle's own F-108 evidence (*a wheel at the
    ┃          origin over the list scrolls the list*) cannot both hold, because        origin over the list scrolls the list*) cannot both hold, because
    ┃          the origin **is** the header's first cell. The scope held: **no          the origin **is** the header's first cell. The scope held: **no
    ┃          click exists**, on purpose; `internal/core`, `internal/cli` and          click exists**, on purpose; `internal/core`, `internal/cli` and
    ┃          `docs/03` are untouched; `core.DomainVersion` is still **7**; the        `docs/03` are untouched; `core.DomainVersion` is still **7**; the
    ┃          budget is still **0 / 0 / 28**; the seven tabs kept their names;         budget is still **0 / 0 / 28**; the seven tabs kept their names;
    ┃          the goldens did not move; `make check` and `make check-race` are         the goldens did not move; `make check` and `make check-race` are
    ┃          green. |                                                                 green. |
    ┃     68   | Slice T2a · cycle 3 | The **last** cycle of T2a: the bound's      68   | Slice T2a · cycle 3 | The **last** cycle of T2a: the bound's
    ┃          tripwire, the wheel table, and the scope | F-109 is closed: the          tripwire, the wheel table, and the scope | F-109 is closed: the
    ┃          tripwire **bites** — the reviewer's mutation, the bound counted at       tripwire **bites** — the reviewer's mutation, the bound counted at
    ┃          the full-screen width again, turns `internal/ui` **red** in that         the full-screen width again, turns `internal/ui` **red** in that
    ┃          one test and **for the bound's width** (*the bound is 27 rows; the       one test and **for the bound's width** (*the bound is 27 rows; the
    ┃          pane has 29 rows at 107 columns* · *the scroll stopped at 26, want       pane has 29 rows at 107 columns* · *the scroll stopped at 26, want
    ┃          the pane's last row 28*). The new test does not ask the code under       the pane's last row 28*). The new test does not ask the code under
    ┃          test for the answer — the pane's width is computed **in the test**       test for the answer — the pane's width is computed **in the test**
    ┃          — and its fixture **makes the two widths differ** (29 rows at 107        — and its fixture **makes the two widths differ** (29 rows at 107
    ┃          columns, 27 at 196) and says so out loud if they ever agree again.       columns, 27 at 196) and says so out loud if they ever agree again.
    ┃          The document's wheel table reads whole: the second row carries           The document's wheel table reads whole: the second row carries
    ┃          *"while the list is on screen"*, and the prose names the case the        *"while the list is on screen"*, and the prose names the case the
    ┃          table carves out; every row is true, proved by command — over the        table carves out; every row is true, proved by command — over the
    ┃          header in the split layout the list scrolls, and in the full-            header in the split layout the list scrolls, and in the full-
    ┃          screen detail the wheel over the header does **nothing** (`enter`        screen detail the wheel over the header does **nothing** (`enter`
    ┃          and `enter,wheel:down@(150,0)` are byte-identical) while the wheel       and `enter,wheel:down@(150,0)` are byte-identical) while the wheel
    ┃          over the detail scrolls it. The wheel's **fallback is untouched**        over the detail scrolls it. The wheel's **fallback is untouched**
    ┃          since cycle 2. The scope held: **no click**; `internal/core`,            since cycle 2. The scope held: **no click**; `internal/core`,
    ┃          `internal/cli` and `docs/03` untouched; `core.DomainVersion` still       `internal/cli` and `docs/03` untouched; `core.DomainVersion` still
    ┃          **7**; the budget still **0 / 0 / 28**; the seven tabs kept their        **7**; the budget still **0 / 0 / 28**; the seven tabs kept their
    ┃          names; the goldens did not move; `make check` and `make check-           names; the goldens did not move; `make check` and `make check-
    ┃          race` green — **no blockers filed** |                                    race` green — **no blockers filed** |
    ┃                                                                              69 + | Slice T2b | The review of the clicks: the tab, the row, the
    ┃                                                                                   detail, and the script | The clicks **are** the keyboard's moves:
    ┃                                                                                   a click on a tab label goes through `Model.switchTab` — and the
    ┃                                                                                   keyboard's `tab` was moved into that same function, so there is
    ┃                                                                                   **one** switch; a click on a row takes the item **the map holds**
    ┃                                                                                   (`hits.rowAt` → `selectItem` → `moveTo`) and lands the cursor
    ┃                                                                                   where the keyboard's `↑` `↓` would; a click on the detail pane is
    ┃                                                                                   `Model.openDetail` — and the keyboard's `enter` was moved into
    ┃                                                                                   that same function, so there is **one** open. The **two layouts**
    ┃                                                                                   were proved by command, not by reading: in the split layout (`--
    ┃                                                                                   width 200`) a click on row 7 puts the detail column on RCP03 and
    ┃                                                                                   leaves the mode split; in the narrow layout (`--width 120`) the
    ┃                                                                                   same click opens RCP03's detail on its own. The **Backlog**
    ┃                                                                                   selects and stops — `openDetail` refuses it, because the Backlog
    ┃                                                                                   holds no recipe — and a click on a **theme** row moves the
    ┃                                                                                   highlight without applying the theme (the ANSI diff shows the
    ┃                                                                                   highlight moved and the palette did not). The **grammar** lives in
    ┃                                                                                   one file: `SplitSteps` and `parseMouseStep` are both in `internal/
    ┃                                                                                   ui/mouse.go`, `click@(12,5)` survives the splitter, and an unknown
    ┃                                                                                   step is refused with a message a reader can act on. The **six**
    ┃                                                                                   mutations each failed the suite, and each was reverted: the tab
    ┃                                                                                   switch without the rest, the click ignoring the map's item, the
    ┃                                                                                   header click selecting, the narrow layout not opening the detail,
    ┃                                                                                   the detail click doing nothing, and the splitter cutting the
    ┃                                                                                   bracketed position at its comma — **two of the six are caught by
    ┃                                                                                   T2a's `mouse_test.go` and not by the click's own tests**, which is
    ┃                                                                                   a note and not a blocker. The scope held: `internal/core` and
    ┃                                                                                   `internal/cli` are untouched, `core.DomainVersion` is still **7**,
    ┃                                                                                   the budget is still **0 / 0 / 28**, the seven tabs kept their
    ┃                                                                                   names, `maitre-cli` did not change, the goldens did not move, and
    ┃                                                                                   `docs/05`'s click rows are true — the key table's *next recipe* →
    ┃                                                                                   *next tab* correction is right. The console answers the **press**
    ┃                                                                                   and not the release, so a drag acts on the row the chef pressed
    ┃                                                                                   on: the standard answer, and the safe one. **No blockers filed** |
    ┃     69   | The rounds inside a cycle | Fix → verify → back to the executor   70   | The rounds inside a cycle | Fix → verify → back to the executor
    ┃          → repeat, until the verifier is satisfied | Three rounds in cycle        → repeat, until the verifier is satisfied | Three rounds in cycle
    ┃          8: the first fix introduced three new mediums, the second closed         8: the first fix introduced three new mediums, the second closed
    ┃          them |                                                                   them |
    ┃     70                                                                       71
    ┃     71   ### The discipline every round runs                                 72   ### The discipline every round runs
    ┃     72                                                                       73
    ┃

       + Thought · 3.3s

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    137   | The plan's name rule is **nowhere in the spec** | **Closed by    137   | The plan's name rule is **nowhere in the spec** | **Closed by
    ┃          the chef's ruling (2026-10-30, Q7):** the rule lives in `docs/03`        the chef's ruling (2026-10-30, Q7):** the rule lives in `docs/03`
    ┃          § 6, beside the plan format — it is about what the reader answers        § 6, beside the plan format — it is about what the reader answers
    ┃          from `project.json`, and § 6 is where the reader looks. It is            from `project.json`, and § 6 is where the reader looks. It is
    ┃          **not** in § 8: the name is a *stored* fact, not a derived view. A       **not** in § 8: the name is a *stored* fact, not a derived view. A
    ┃          wire (`TestThePlanFormatStatesTheNameRule`) holds the sentence,          wire (`TestThePlanFormatStatesTheNameRule`) holds the sentence,
    ┃          and `TestTheNameRuleTripwireFailsOnRealDrift` proves the wire can        and `TestTheNameRuleTripwireFailsOnRealDrift` proves the wire can
    ┃          fail. |                                                                  fail. |
    ┃    138   | A test assertion that **cannot fail** | `TestTheRootRowIsOneRow` 138   | A test assertion that **cannot fail** | `TestTheRootRowIsOneRow`
    ┃          splits the rows on `\n` and then checks each line for a `\n` — a         splits the rows on `\n` and then checks each line for a `\n` — a
    ┃          line that came out of a split cannot carry one. The test still           line that came out of a split cannot carry one. The test still
    ┃          bites through its other two assertions (*the root row does not           bites through its other two assertions (*the root row does not
    ┃          carry the sanitised path*, *the root row carries a newline*), so         carry the sanitised path*, *the root row carries a newline*), so
    ┃          the rule is held; the loop is decoration. |                              the rule is held; the loop is decoration. |
    ┃    139   | The render **writes** the scroll's offset | `detailScrollRows`   139   | The render **writes** the scroll's offset | `detailScrollRows`
    ┃          clamps `m.detailScroll` — and writes it — while the frame is being       clamps `m.detailScroll` — and writes it — while the frame is being
    ┃          drawn. It moved there from `renderDetailFull`, which clamped it          drawn. It moved there from `renderDetailFull`, which clamped it
    ┃          inline; it is **not new**, and the render already writes `m.hits`.       inline; it is **not new**, and the render already writes `m.hits`.
    ┃          Worth knowing before T2b adds more to the render path. |                 Worth knowing before T2b adds more to the render path. |
    ┃    140   | The wheel over the list moves the **selection** | In the split   140   | The wheel over the list moves the **selection** | In the split
    ┃          layout the wheel over the list is the keyboard's `down`, so the          layout the wheel over the list is the keyboard's `down`, so the
    ┃          selection moves and the detail pane shows **another recipe from          selection moves and the detail pane shows **another recipe from
    ┃          its top** — the detail does not scroll. That is the console's            its top** — the detail does not scroll. That is the console's
    ┃          model (one answer, not two), but a reader of the pane may expect         model (one answer, not two), but a reader of the pane may expect
    ┃          the recipe to stay put. |                                                the recipe to stay put. |
    ┃                                                                             141 + | `internal/ui/mouse.go`'s doc comment prints the click's two
    ┃                                                                                   examples **twice** | A copy-paste slip in `MouseFor`'s comment:
    ┃                                                                                   the `"click"` and `"click@(12,5)"` lines are there twice. Cosmetic
    ┃                                                                                   — the next lane that touches the file tidies it. |
    ┃                                                                             142 + | The click's row test is **self-referential** |
    ┃                                                                                   `TestAClickOnARowSelectsIt` finds the row to click **from the
    ┃                                                                                   map** (`rowFor`) and then compares `visibleRows[sel]` with
    ┃                                                                                   `visibleRows[1]` — a comparison that holds whenever `sel == 1`,
    ┃                                                                                   which the line above already asserted. The map's contract (*a row
    ┃                                                                                   holds the **item** index; a header holds none*) is held by T2a's
    ┃                                                                                   `internal/ui/mouse_test.go`, so a wrong map still turns the suite
    ┃                                                                                   red — the reviewer's mutation (the map holding the list's line
    ┃                                                                                   index) failed there, not here. |
    ┃                                                                             143 + | A header click that selects **item 0** is invisible |
    ┃                                                                                   `TestAClickOnAHeaderRowSelectsNothing` records the cursor before
    ┃                                                                                   the click and compares it after; the cursor starts at **0**, so a
    ┃                                                                                   header click that selects item 0 changes nothing and the test is
    ┃                                                                                   green. The map's *"a header holds no item"* is held by
    ┃                                                                                   `mouse_test.go`. |
    ┃                                                                             144 + | A malformed **position** is refused as an unknown step |
    ┃                                                                                   `isMouseStep` cuts the name on `:`, so `click@(1,2` has the verb
    ┃                                                                                   `click@(1,2` and is not a mouse step at all — the driver says
    ┃                                                                                   *unknown step "click@(1,2"* where *malformed mouse step* would
    ┃                                                                                   name the mistake. `click:sideways` is refused correctly, and
    ┃                                                                                   `click@1,2` is cut by the splitter before it is read. |
    ┃                                                                             145 + | The **Log** tab's detail is unreachable in the narrow layout |
    ┃                                                                                   `openDetail` opens a **recipe**, so neither `enter` nor a click
    ┃                                                                                   opens the journal entry's detail where the list is alone — the
    ┃                                                                                   pane is only reachable in the split layout. This **predates** the
    ┃                                                                                   slice (the keyboard's `enter` never opened it), and the click is
    ┃                                                                                   faithful to the keyboard. T3a's territory. |
    ┃    141                                                                      146
    ┃    142   ## Closed findings                                                 147   ## Closed findings
    ┃    143                                                                      148
    ┃    144   | # | Severity | Closed in | Held by |                             149   | # | Severity | Closed in | Held by |
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff --stat && make fmt-check 2>&1 | tail -3 && git status --short
    ┃
    ┃  docs/audit.md | 6 ++++++
    ┃   1 file changed, 6 insertions(+)
    ┃
    ┃  ==> fmt-check
    ┃   M docs/audit.md
    ┃

       + Thought · 1.8s

    ┃
    ┃  $ cd /tmp && rm -rf kb4 && cp -R kb2 kb4 && cd kb4 && cp kitchen/activity.json /tmp/act.before && /Users/thiago/Code/Pessoal/maitre/bin/maitre --re
    ┃  nder --width 120 --keys "click@(14,2),click@(5,6)" >/dev/null && cmp /tmp/act.before kitchen/activity.json && echo "activity.json UNCHANGED"…
    ┃
    ┃  activity.json UNCHANGED
    ┃         7
    ┃

       + Thought · 893ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/audit.md && git commit -q -m "docs: T2b reviewed — the clicks are the keyboard's moves, and
    ┃  the slice is done" && git log --oneline -2 && git status --short
    ┃
    ┃  ee6a935 docs: T2b reviewed — the clicks are the keyboard's moves, and the slice is done
    ┃  8fa6cbd work: T2b's review brief - the clicks are the keyboard's moves
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && mkdir -p work/inbox && cat > work/inbox/reviewer-tui-t2b.md <<EOF
    ┃  Reviewer: T2b review done. Verdict: ok. Filed 0 blockers (none). Mutation checks: (1) the click set the tab without the switch's other work — \…
    ┃
    ┃  9518071 inbox: reviewer reports tui-t2b
    ┃  ee6a935 docs: T2b reviewed — the clicks are the keyboard's moves, and the slice is done
    ┃  TREE CLEAN
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1H:p1" "Reviewer: T2b review done. Verdict: ok. Filed 0 blockers (none). Mutation ch
    ┃  ecks: 6 mutations, each reverted, tree clean; the tab switch without the rest, the click ignoring the map's item, the header click selecting…
    ┃
    ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":
    ┃  "ses_e737d6341ffeTdoEYTgykU25M8"},"agent_status":"done","completion_seq":185,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,
    ┃  "foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1H:p1","revision":15,"screen_detection_skipped":true,"state_change_seq":185,
    ┃  "tab_id":"w1H:t1","terminal_id":"term_65f0cd2aa89d26","terminal_title":"OC | Project work proposal","terminal_title_stripped":"OC | Project work
    ┃  proposal","workspace_id":"w1H"},"type":"agent_prompted"}}
    ┃

       + Thought · 377ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && git log --oneline -3 && git diff HEAD~2 --stat
    ┃
    ┃  9518071 inbox: reviewer reports tui-t2b
    ┃  ee6a935 docs: T2b reviewed — the clicks are the keyboard's moves, and the slice is done
    ┃  8fa6cbd work: T2b's review brief - the clicks are the keyboard's moves
    ┃   docs/audit.md                  | 6 ++++++
    ┃   work/inbox/reviewer-tui-t2b.md | 1 +
    ┃   2 files changed, 7 insertions(+)
    ┃

       verdict: ok — 0 blockers filed. Receipt: work/inbox/reviewer-tui-t2b.md (committed 9518071); ledger: docs/audit.md (ee6a935). The nudge landed on
       w1H:p1.

       The clicks are the keyboard's moves

       ┌─────────────┬────────────────────────────────────────────┬──────────────────────────────────────────────────────────────────────────────────────┐
       │ Click       │ The function                               │ Note                                                                                 │
       ├─────────────┼────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────┤
       │ tab label   │ Model.switchTab (internal/ui/model.go:631) │ the keyboard's tab/right/l was moved into it — one switch                            │
       ├─────────────┼────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────┤
       │ list row    │ hits.rowAt → selectItem → Model.moveTo     │ the item is the map's, and moveTo is the keyboard's selection state, clamping and    │
       │             │                                            │ scroll-into-view                                                                     │
       ├─────────────┼────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────┤
       │ detail pane │ Model.openDetail (model.go:644)            │ the keyboard's enter was moved into it — one open                                    │
       └─────────────┴────────────────────────────────────────────┴──────────────────────────────────────────────────────────────────────────────────────┘

       Proved by command (not by reading)

       - split --width 200: a click on row 7 puts the detail column on RCP03, mode stays split. narrow --width 120: the same click opens RCP03's detail
         alone.
       - Backlog: selects and stops (openDetail refuses — no recipe). Settings: the theme row is selected, not applied (ANSI diff: the highlight moved,
         the palette did not).
       - maitre --render --keys "click@(14,2)" → exit 0, marker moves to Backlog. A click writes nothing (cmp on activity.json).
       - Grammar: SplitSteps and parseMouseStep are both in internal/ui/mouse.go — one file, one grammar.

       Mutations — all reverted, tree clean

       6 mutations, each failed the suite. Two of them (the map holding the list's line index; a header click selecting) were caught by T2a's mouse_test.
       go, not by the click's own tests — that gap is in the notes table.

       Press or release

       The console answers the press (bubbletea v2's MouseClickMsg is button-down; the release is ignored). A drag acts on the pressed row. Right: it is
       the standard, it agrees with the wheel (which follows the pointer), and a release rule would let a drag ending over the detail open it by accident.

       Scope, gate, goldens

       internal/core/internal/cli untouched, DomainVersion 7, budget 0/0/28, seven tabs, maitre-cli unchanged, goldens unchanged, docs/05's click rows
       true (the next recipe → next tab correction is right). make check and make check-race both pass.

       LOW notes (in the ledger): the duplicated doc comment in internal/ui/mouse.go; the self-referential row-select test; the invisible header click at
       item 0; a malformed position refused as unknown step; the Log tab's narrow-layout detail (pre-existing).

       Build · MiniMax-M3.1-Flash · 16m 11s · 148.5 tok/s

    ┃
    ┃
    ┃
    ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
    ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
    ~/Code/Pessoal/maitre:master                                                                                             125.9K (13%)  ctrl+p commands
```
