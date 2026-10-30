# The runner's lane — T2b, captured 2026-10-30

Dumped from pane `w1H:p6` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that wrote T2b: the clicks — a tab, a row, a detail — and the click in the
same grammar as the wheel.
The review that judged it is in [tui-t2b-reviewer.md](tui-t2b-reviewer.md).

    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/ui/mouse.go /tmp/mouse.go.bak && cp internal/ui/model.go /tmp/model.go.bak && python3 - <<'EOF'
    ┃  p='internal/ui/model.go'…
    ┃
    ┃  --- FAIL: TestTheClickOnATabClampsTheCursors (0.79s)
    ┃      click_test.go:114: the click on a tab left the status "signed off RCP02", want it cleared
    ┃      click_test.go:117: the click on a tab left the cursor at 99, want it inside the 6 rows
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  1.069s
    ┃  FAIL
    ┃

       + Thought · 13.7s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && mkdir -p /tmp/t2b-final && cp internal/ui/mouse.go internal/ui/model.go internal/ui/click_test.go /tmp/t2b-final/ && ls /tmp/t2b-final
    ┃
    ┃  click_test.go
    ┃  model.go
    ┃  mouse.go
    ┃

       + Thought · 1.7s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
    ┃  import re…
    ┃
    ┃  BUILD_OK
    ┃

    ┃
    ┃  # Wrote internal/ui/click_test.go
    ┃
    ┃     1 package ui
    ┃     2
    ┃     3 import (
    ┃     4   "fmt"
    ┃     5   "testing"
    ┃     6
    ┃     7   "charm.land/bubbletea/v2"
    ┃     8 )
    ┃     9
    ┃    10 // click sends a click to the model at a point.
    ┃    11 //
    ┃    12 // The click is a press, the way the terminal sends it: the release is the same
    ┃    13 // gesture coming back, and the console answers the press only.
    ┃    14 func click(t *testing.T, m *Model, x, y int) {
    ┃    15   t.Helper()
    ┃    16   m.Update(tea.MouseClickMsg{X: x, Y: y, Button: tea.MouseLeft})
    ┃    17 }
    ┃    18
    ┃    19 // selections is every cursor of the console, in one string.
    ┃    20 //
    ┃    21 // Two consoles in the same place have the same cursors: the click and the
    ┃    22 // keyboard must leave the same state behind.
    ┃    23 func selections(m *Model) string {
    ┃    24   return fmt.Sprintf("recipes=%d backlog=%d path=%d board=%d log=%d theme=%d",
    ┃    25     m.selRecipes, m.selBacklog, m.selPath, m.selBoard, m.selLog, m.selTheme)
    ┃    26 }
    ┃    27
    ┃    28 // TestAClickOnATabLabelSwitchesTheTab covers the tab bar.
    ┃    29 func TestAClickOnATabLabelSwitchesTheTab(t *testing.T) {
    ┃    30   m := newModel(t)
    ┃    31   m.width = 120
    ┃    32   m.height = 30
    ┃    33   m.View()
    ┃    34
    ┃    35   x, y := m.hits.tabs[int(TabBacklog)].centre()
    ┃    36   click(t, m, x, y)
    ┃    37   if m.tab != TabBacklog {
    ┃    38     t.Errorf("the click on the %q label left the console on %q", TabBacklog.Label(), m.tab.Label())
    ┃    39   }
    ┃    40 }
    ┃    41
    ┃    42 // TestTheClickOnATabIsTheKeyboardsSwitch checks the click makes the keyboard's move.
    ┃    43 //
    ┃    44 // The keyboard's `tab` and the mouse's click on the next tab's label must land in
    ┃    45 // the same place: one switch, one answer.
    ┃    46 func TestTheClickOnATabIsTheKeyboardsSwitch(t *testing.T) {
    ┃    47   clicked := newModel(t)
    ┃    48   clicked.width = 120
    ┃    49   clicked.height = 30
    ┃    50   clicked.View()
    ┃    51   x, y := clicked.hits.tabs[int(clicked.tab)+1].centre()
    ┃    52   click(t, clicked, x, y)
    ┃    53
    ┃    54   keyed := newModel(t)
    ┃    55   keyed.width = 120
    ┃    56   keyed.height = 30
    ┃    57   keyed.Update(tea.KeyPressMsg{Code: tea.KeyTab})
    ┃    58
    ┃    59   if clicked.tab != keyed.tab {
    ┃    60     t.Errorf("the click landed on %q, the keyboard on %q", clicked.tab.Label(), keyed.tab.Label())
    ┃    61   }
    ┃    62   if selections(clicked) != selections(keyed) {
    ┃    63     t.Errorf("the click and the keyboard left different cursors: %s vs %s",
    ┃    64       selections(clicked), selections(keyed))
    ┃    65   }
    ┃    66   if clicked.status != keyed.status {
    ┃    67     t.Errorf("the click left the status %q, the keyboard %q", clicked.status, keyed.status)
    ┃    68   }
    ┃    69 }
    ┃    70
    ┃    71 // TestTheClickOnATabClampsTheCursors covers the rest of the switch.
    ┃    72 //
    ┃    73 // The switch is the keyboard's: it resets the status and keeps every cursor
    ┃    74 // inside its list. A click that set the tab and skipped the rest would be a
    ┃    75 // second answer to one question.
    ┃    76 func TestTheClickOnATabClampsTheCursors(t *testing.T) {
    ┃    77   m := newModel(t)
    ┃    78   m.width = 120
    ┃    79   m.height = 30
    ┃    80   m.status = "signed off RCP02"
    ┃    81   m.selBacklog = 99
    ┃    82   m.View()
    ┃    83
    ┃    84   x, y := m.hits.tabs[int(TabBacklog)].centre()
    ┃    85   click(t, m, x, y)
    ┃    86   if m.status != "" {
    ┃    87     t.Errorf("the click on a tab left the status %q, want it cleared", m.status)
    ┃    88   }
    ┃    89   if m.selBacklog >= m.backlogCount() {
    ┃    90     t.Errorf("the click on a tab left the cursor at %d, want it inside the %d rows",
    ┃    91       m.selBacklog, m.backlogCount())
    ┃    92   }
    ┃    93 }
    ┃    94
    ┃    95 // TestAClickOnTheGapBetweenTabsDoesNothing covers the separator.
    ┃    96 //
    ┃    97 // The tab bar joins its parts with a separator, and the separator belongs to no
    ┃    98 // tab: a click on it changes nothing.
    ┃    99 func TestAClickOnTheGapBetweenTabsDoesNothing(t *testing.T) {
    ┃   100   m := newModel(t)
    ┃   101   m.width = 120
    ┃   102   m.height = 30
    ┃   103   m.View()
    ┃   104
    ┃   105   before := m.tab
    ┃   106   gap := m.hits.tabs[0].x1 + 1
    ┃   107   click(t, m, gap, tabBarRow)
    ┃   108   if m.tab != before {
    ┃   109     t.Errorf("the click on the gap between the labels went to %q", m.tab.Label())
    ┃   110   }
    ┃   111 }
    ┃   112
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui/ -run 'Click' 2>&1 | tail -5 && staticcheck ./internal/ui/ && echo STATIC_OK
    ┃
    ┃  ok    github.com/lthiagol/maitre/internal/ui  4.414s
    ┃  zsh:1: command not found: staticcheck
    ┃
    ┃  Exited with code 127
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make lint 2>&1 | tail -10
    ┃
    ┃  (3 earlier lines)
    ┃  go vet ./...
    ┃  /Library/Developer/CommandLineTools/usr/bin/make staticcheck
    ┃
    ┃  ==> staticcheck
    ┃  GO=go scripts/check-tool-versions.sh tool staticcheck v0.8.2
    ┃    staticcheck 0.8.2 matches the pin v0.8.2
    ┃  GO=go scripts/require-tool.sh staticcheck ./...
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/ui/model.go internal/ui/mouse.go internal/ui/click_test.go && git commit -q -m "feat(ui): a click on a tab label switches the tab
    ┃  …
    ┃
    ┃  4dfd5a4 feat(ui): a click on a tab label switches the tab
    ┃

       + Thought · 349ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
    ┃  p='internal/ui/mouse.go'…
    ┃
    ┃  BUILD_OK
    ┃

    ┃
    ┃  ← Edit internal/ui/click_test.go
    ┃
    ┃     38       t.Errorf("the click on the %q label left the console on %q", TabBacklog.Label(), m.tab.Label())                                                 38       t.Errorf("the click on the %q label left the console on %q", TabBacklog.Label(), m.tab.Label())
    ┃     39     }                                                                                                                                                 39     }
    ┃     40   }                                                                                                                                                   40   }
    ┃     41                                                                                                                                                       41
    ┃                                                                                                                                                              42 + // rowFor is the screen row that holds an item, and whether there is one.
    ┃                                                                                                                                                              43 + func rowFor(m *Model, item int) (int, bool) {
    ┃                                                                                                                                                              44 +   for _, r := range m.hits.rows {
    ┃                                                                                                                                                              45 +     if r.item == item {
    ┃                                                                                                                                                              46 +       return r.y0, true
    ┃                                                                                                                                                              47 +     }
    ┃                                                                                                                                                              48 +   }
    ┃                                                                                                                                                              49 +   return 0, false
    ┃                                                                                                                                                              50 + }
    ┃                                                                                                                                                              51 +
    ┃                                                                                                                                                              52 + // headerRow is the screen row of a header — a row that holds no item.
    ┃                                                                                                                                                              53 + func headerRow(m *Model) (int, bool) {
    ┃                                                                                                                                                              54 +   for _, r := range m.hits.rows {
    ┃                                                                                                                                                              55 +     if r.item < 0 {
    ┃                                                                                                                                                              56 +       return r.y0, true
    ┃                                                                                                                                                              57 +     }
    ┃                                                                                                                                                              58 +   }
    ┃                                                                                                                                                              59 +   return 0, false
    ┃                                                                                                                                                              60 + }
    ┃                                                                                                                                                              61 +
    ┃                                                                                                                                                              62 + // TestAClickOnARowSelectsIt covers the list.
    ┃                                                                                                                                                              63 + func TestAClickOnARowSelectsIt(t *testing.T) {
    ┃                                                                                                                                                              64 +   m := newModel(t)
    ┃                                                                                                                                                              65 +   m.width = 120
    ┃                                                                                                                                                              66 +   m.height = 30
    ┃                                                                                                                                                              67 +   m.View()
    ┃                                                                                                                                                              68 +
    ┃                                                                                                                                                              69 +   row, ok := rowFor(m, 1)
    ┃                                                                                                                                                              70 +   if !ok {
    ┃                                                                                                                                                              71 +     t.Fatal("the map holds no row with an item")
    ┃                                                                                                                                                              72 +   }
    ┃                                                                                                                                                              73 +   click(t, m, 5, row)
    ┃                                                                                                                                                              74 +   if m.selRecipes != 1 {
    ┃                                                                                                                                                              75 +     t.Errorf("the click on the row left the cursor on %d, want 1", m.selRecipes)
    ┃                                                                                                                                                              76 +   }
    ┃                                                                                                                                                              77 +   // The selection is the tab's own selection: the recipe under it is the recipe
    ┃                                                                                                                                                              78 +   // the row drew.
    ┃                                                                                                                                                              79 +   rows := m.visibleRows()
    ┃                                                                                                                                                              80 +   if got := rows[m.selRecipes].ID; got != rows[1].ID {
    ┃                                                                                                                                                              81 +     t.Errorf("the cursor is on %s, want %s", got, rows[1].ID)
    ┃                                                                                                                                                              82 +   }
    ┃                                                                                                                                                              83 + }
    ┃                                                                                                                                                              84 +
    ┃                                                                                                                                                              85 + // TestAClickOnAHeaderRowSelectsNothing covers the headers.
    ┃                                                                                                                                                              86 + //
    ┃                                                                                                                                                              87 + // A header is not a thing to select: the map says the row holds no item, and the
    ┃                                                                                                                                                              88 + // click leaves the cursor where it was.
    ┃                                                                                                                                                              89 + func TestAClickOnAHeaderRowSelectsNothing(t *testing.T) {
    ┃                                                                                                                                                              90 +   m := newModel(t)
    ┃                                                                                                                                                              91 +   m.width = 120
    ┃                                                                                                                                                              92 +   m.height = 30
    ┃                                                                                                                                                              93 +   m.View()
    ┃                                                                                                                                                              94 +
    ┃                                                                                                                                                              95 +   before := m.selRecipes
    ┃                                                                                                                                                              96 +   row, ok := headerRow(m)
    ┃                                                                                                                                                              97 +   if !ok {
    ┃                                                                                                                                                              98 +     t.Fatal("the map holds no header row")
    ┃                                                                                                                                                              99 +   }
    ┃                                                                                                                                                             100 +   click(t, m, 5, row)
    ┃                                                                                                                                                             101 +   if m.selRecipes != before {
    ┃                                                                                                                                                             102 +     t.Errorf("the click on a header moved the cursor to %d, want %d", m.selRecipes, before)
    ┃                                                                                                                                                             103 +   }
    ┃                                                                                                                                                             104 + }
    ┃                                                                                                                                                             105 +
    ┃     42   // TestTheClickOnATabIsTheKeyboardsSwitch checks the click makes the keyboard's move.                                                              106   // TestTheClickOnATabIsTheKeyboardsSwitch checks the click makes the keyboard's move.
    ┃     43   //                                                                                                                                                 107   //
    ┃     44   // The keyboard's `tab` and the mouse's click on the next tab's label must land in                                                                 108   // The keyboard's `tab` and the mouse's click on the next tab's label must land in
    ┃     45   // the same place: one switch, one answer.                                                                                                         109   // the same place: one switch, one answer.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui/ -run 'Click' 2>&1 | tail -5
    ┃
    ┃  ok    github.com/lthiagol/maitre/internal/ui  6.614s
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/ui/mouse.go internal/ui/click_test.go && git commit -q -m "feat(ui): a click on a row selects it
    ┃  …
    ┃
    ┃  187d706 feat(ui): a click on a row selects it
    ┃

       + Thought · 869ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
    ┃  p='internal/ui/mouse.go'…
    ┃
    ┃  BUILD_OK
    ┃

    ┃
    ┃  ← Edit internal/ui/click_test.go
    ┃
    ┃    102       t.Errorf("the click on a header moved the cursor to %d, want %d", m.selRecipes, before)                                                        102       t.Errorf("the click on a header moved the cursor to %d, want %d", m.selRecipes, before)
    ┃    103     }                                                                                                                                                103     }
    ┃    104   }                                                                                                                                                  104   }
    ┃    105                                                                                                                                                      105
    ┃                                                                                                                                                             106 + // TestAClickOnARowInTheNarrowLayoutOpensTheDetail covers the narrow layout.
    ┃                                                                                                                                                             107 + //
    ┃                                                                                                                                                             108 + // The list is alone there — the detail is not on screen — so a click that only
    ┃                                                                                                                                                             109 + // selected would show the chef nothing. The click also opens the detail: the
    ┃                                                                                                                                                             110 + // keyboard's `enter`.
    ┃                                                                                                                                                             111 + func TestAClickOnARowInTheNarrowLayoutOpensTheDetail(t *testing.T) {
    ┃                                                                                                                                                             112 +   m := newModel(t)
    ┃                                                                                                                                                             113 +   m.width = 120
    ┃                                                                                                                                                             114 +   m.height = 30
    ┃                                                                                                                                                             115 +   m.View()
    ┃                                                                                                                                                             116 +
    ┃                                                                                                                                                             117 +   if !m.hits.detail.empty() {
    ┃                                                                                                                                                             118 +     t.Fatalf("the narrow layout should have no detail on screen: %+v", m.hits.detail)
    ┃                                                                                                                                                             119 +   }
    ┃                                                                                                                                                             120 +   row, ok := rowFor(m, 1)
    ┃                                                                                                                                                             121 +   if !ok {
    ┃                                                                                                                                                             122 +     t.Fatal("the map holds no row with an item")
    ┃                                                                                                                                                             123 +   }
    ┃                                                                                                                                                             124 +   click(t, m, 5, row)
    ┃                                                                                                                                                             125 +   if m.mode != ModeDetail {
    ┃                                                                                                                                                             126 +     t.Errorf("the click in the narrow layout left the console in %v, want the detail on screen", m.mode)
    ┃                                                                                                                                                             127 +   }
    ┃                                                                                                                                                             128 +   rec, ok := m.currentRecipe()
    ┃                                                                                                                                                             129 +   if !ok {
    ┃                                                                                                                                                             130 +     t.Fatal("the detail has no recipe")
    ┃                                                                                                                                                             131 +   }
    ┃                                                                                                                                                             132 +   want := m.visibleRows()[1].ID
    ┃                                                                                                                                                             133 +   if rec.ID != want {
    ┃                                                                                                                                                             134 +     t.Errorf("the detail shows %s, want %s", rec.ID, want)
    ┃                                                                                                                                                             135 +   }
    ┃                                                                                                                                                             136 + }
    ┃                                                                                                                                                             137 +
    ┃                                                                                                                                                             138 + // TestAClickOnARowInTheSplitLayoutOnlySelects covers the split layout.
    ┃                                                                                                                                                             139 + //
    ┃                                                                                                                                                             140 + // The detail is on screen there, beside the list, and it follows the selection —
    ┃                                                                                                                                                             141 + // so the click selects, and nothing else.
    ┃                                                                                                                                                             142 + func TestAClickOnARowInTheSplitLayoutOnlySelects(t *testing.T) {
    ┃                                                                                                                                                             143 +   m := newModel(t)
    ┃                                                                                                                                                             144 +   m.width = 200
    ┃                                                                                                                                                             145 +   m.height = 30
    ┃                                                                                                                                                             146 +   m.View()
    ┃                                                                                                                                                             147 +
    ┃                                                                                                                                                             148 +   if m.hits.detail.empty() {
    ┃                                                                                                                                                             149 +     t.Fatal("the split layout should have the detail on screen")
    ┃                                                                                                                                                             150 +   }
    ┃                                                                                                                                                             151 +   row, ok := rowFor(m, 1)
    ┃                                                                                                                                                             152 +   if !ok {
    ┃                                                                                                                                                             153 +     t.Fatal("the map holds no row with an item")
    ┃                                                                                                                                                             154 +   }
    ┃                                                                                                                                                             155 +   click(t, m, 5, row)
    ┃                                                                                                                                                             156 +   if m.mode != ModeSplit {
    ┃                                                                                                                                                             157 +     t.Errorf("the click in the split layout opened the detail: mode is %v", m.mode)
    ┃                                                                                                                                                             158 +   }
    ┃                                                                                                                                                             159 +   if m.selRecipes != 1 {
    ┃                                                                                                                                                             160 +     t.Errorf("the click left the cursor on %d, want 1", m.selRecipes)
    ┃                                                                                                                                                             161 +   }
    ┃                                                                                                                                                             162 +   rec, ok := m.currentRecipe()
    ┃                                                                                                                                                             163 +   if !ok {
    ┃                                                                                                                                                             164 +     t.Fatal("the detail has no recipe")
    ┃                                                                                                                                                             165 +   }
    ┃                                                                                                                                                             166 +   if want := m.visibleRows()[1].ID; rec.ID != want {
    ┃                                                                                                                                                             167 +     t.Errorf("the detail shows %s, want %s", rec.ID, want)
    ┃                                                                                                                                                             168 +   }
    ┃                                                                                                                                                             169 + }
    ┃                                                                                                                                                             170 +
    ┃                                                                                                                                                             171 + // TestAClickOnTheBacklogSelectsAndStops covers the tab with no detail.
    ┃                                                                                                                                                             172 + //
    ┃                                                                                                                                                             173 + // The Backlog has no detail yet — that is T3a — so a click there selects and
    ┃                                                                                                                                                             174 + // stops.
    ┃                                                                                                                                                             175 + func TestAClickOnTheBacklogSelectsAndStops(t *testing.T) {
    ┃                                                                                                                                                             176 +   m := newModel(t)
    ┃                                                                                                                                                             177 +   m.width = 120
    ┃                                                                                                                                                             178 +   m.height = 30
    ┃                                                                                                                                                             179 +   m.tab = TabBacklog
    ┃                                                                                                                                                             180 +   m.View()
    ┃                                                                                                                                                             181 +
    ┃                                                                                                                                                             182 +   row, ok := rowFor(m, 1)
    ┃                                                                                                                                                             183 +   if !ok {
    ┃                                                                                                                                                             184 +     t.Fatal("the map holds no backlog row with an item")
    ┃                                                                                                                                                             185 +   }
    ┃                                                                                                                                                             186 +   click(t, m, 5, row)
    ┃                                                                                                                                                             187 +   if m.selBacklog != 1 {
    ┃                                                                                                                                                             188 +     t.Errorf("the click left the cursor on %d, want 1", m.selBacklog)
    ┃                                                                                                                                                             189 +   }
    ┃                                                                                                                                                             190 +   if m.mode != ModeSplit {
    ┃                                                                                                                                                             191 +     t.Errorf("the click on the backlog opened a detail: mode is %v", m.mode)
    ┃                                                                                                                                                             192 +   }
    ┃                                                                                                                                                             193 + }
    ┃                                                                                                                                                             194 +
    ┃                                                                                                                                                             195 + // TestAClickOnAThemeRowSelectsItAndDoesNotApplyIt covers the Settings tab.
    ┃                                                                                                                                                             196 + //
    ┃                                                                                                                                                             197 + // The theme is applied by `enter`, and applying it is a write. A click selects
    ┃                                                                                                                                                             198 + // the theme under the pointer and stops there.
    ┃                                                                                                                                                             199 + func TestAClickOnAThemeRowSelectsItAndDoesNotApplyIt(t *testing.T) {
    ┃                                                                                                                                                             200 +   m := newModel(t)
    ┃                                                                                                                                                             201 +   m.width = 120
    ┃                                                                                                                                                             202 +   m.height = 30
    ┃                                                                                                                                                             203 +   m.tab = TabSettings
    ┃                                                                                                                                                             204 +   m.View()
    ┃                                                                                                                                                             205 +
    ┃                                                                                                                                                             206 +   row, ok := rowFor(m, 1)
    ┃                                                                                                                                                             207 +   if !ok {
    ┃                                                                                                                                                             208 +     t.Fatal("the map holds no theme row")
    ┃                                                                                                                                                             209 +   }
    ┃                                                                                                                                                             210 +   before := m.theme.name
    ┃                                                                                                                                                             211 +   click(t, m, 5, row)
    ┃                                                                                                                                                             212 +   if m.selTheme != 1 {
    ┃                                                                                                                                                             213 +     t.Errorf("the click left the cursor on theme %d, want 1", m.selTheme)
    ┃                                                                                                                                                             214 +   }
    ┃                                                                                                                                                             215 +   if m.theme.name != before {
    ┃                                                                                                                                                             216 +     t.Errorf("the click applied the theme %q", m.theme.name)
    ┃                                                                                                                                                             217 +   }
    ┃                                                                                                                                                             218 +   if m.mode != ModeSplit {
    ┃                                                                                                                                                             219 +     t.Errorf("the click on a theme opened a detail: mode is %v", m.mode)
    ┃                                                                                                                                                             220 +   }
    ┃                                                                                                                                                             221 + }
    ┃                                                                                                                                                             222 +
    ┃    106   // TestTheClickOnATabIsTheKeyboardsSwitch checks the click makes the keyboard's move.                                                              223   // TestTheClickOnATabIsTheKeyboardsSwitch checks the click makes the keyboard's move.
    ┃    107   //                                                                                                                                                 224   //
    ┃    108   // The keyboard's `tab` and the mouse's click on the next tab's label must land in                                                                 225   // The keyboard's `tab` and the mouse's click on the next tab's label must land in
    ┃    109   // the same place: one switch, one answer.                                                                                                         226   // the same place: one switch, one answer.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui/ -run 'Click' 2>&1 | tail -8
    ┃
    ┃  ok    github.com/lthiagol/maitre/internal/ui  9.877s
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/ui/mouse.go internal/ui/model.go internal/ui/click_test.go && git commit -q -m "feat(ui): a click on a row shows that recipe
    ┃  …
    ┃
    ┃  f64e6a2 feat(ui): a click on a row shows that recipe
    ┃

       + Thought · 232ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
    ┃  p='internal/ui/mouse.go'…
    ┃
    ┃  BUILD_OK
    ┃

    ┃
    ┃  ← Edit internal/ui/click_test.go
    ┃
    ┃    219       t.Errorf("the click on a theme opened a detail: mode is %v", m.mode)                                                                           219       t.Errorf("the click on a theme opened a detail: mode is %v", m.mode)
    ┃    220     }                                                                                                                                                220     }
    ┃    221   }                                                                                                                                                  221   }
    ┃    222                                                                                                                                                      222
    ┃                                                                                                                                                             223 + // TestAClickOnTheDetailPaneOpensIt covers the detail pane.
    ┃                                                                                                                                                             224 + //
    ┃                                                                                                                                                             225 + // A click inside the detail pane is the console's `enter` on that recipe: the
    ┃                                                                                                                                                             226 + // detail takes the screen, and `esc` brings it back.
    ┃                                                                                                                                                             227 + func TestAClickOnTheDetailPaneOpensIt(t *testing.T) {
    ┃                                                                                                                                                             228 +   m := newModel(t)
    ┃                                                                                                                                                             229 +   m.width = 200
    ┃                                                                                                                                                             230 +   m.height = 30
    ┃                                                                                                                                                             231 +   m.View()
    ┃                                                                                                                                                             232 +
    ┃                                                                                                                                                             233 +   x, y := m.hits.detail.centre()
    ┃                                                                                                                                                             234 +   click(t, m, x, y)
    ┃                                                                                                                                                             235 +   if m.mode != ModeDetail {
    ┃                                                                                                                                                             236 +     t.Errorf("the click on the detail left the console in %v, want the detail on screen", m.mode)
    ┃                                                                                                                                                             237 +   }
    ┃                                                                                                                                                             238 + }
    ┃                                                                                                                                                             239 +
    ┃                                                                                                                                                             240 + // TestAClickOnTheFullScreenDetailDoesNothing covers the full-screen detail.
    ┃                                                                                                                                                             241 + //
    ┃                                                                                                                                                             242 + // The detail already has the screen, and `enter` does nothing there — so the
    ┃                                                                                                                                                             243 + // click does nothing there too.
    ┃                                                                                                                                                             244 + func TestAClickOnTheFullScreenDetailDoesNothing(t *testing.T) {
    ┃                                                                                                                                                             245 +   m := newModel(t)
    ┃                                                                                                                                                             246 +   m.width = 200
    ┃                                                                                                                                                             247 +   m.height = 30
    ┃                                                                                                                                                             248 +   m.View()
    ┃                                                                                                                                                             249 +   m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
    ┃                                                                                                                                                             250 +   m.View()
    ┃                                                                                                                                                             251 +
    ┃                                                                                                                                                             252 +   x, y := m.hits.detail.centre()
    ┃                                                                                                                                                             253 +   click(t, m, x, y)
    ┃                                                                                                                                                             254 +   if m.mode != ModeDetail {
    ┃                                                                                                                                                             255 +     t.Errorf("the click on the open detail left the console in %v", m.mode)
    ┃                                                                                                                                                             256 +   }
    ┃                                                                                                                                                             257 + }
    ┃                                                                                                                                                             258 +
    ┃    223   // TestTheClickOnATabIsTheKeyboardsSwitch checks the click makes the keyboard's move.                                                              259   // TestTheClickOnATabIsTheKeyboardsSwitch checks the click makes the keyboard's move.
    ┃    224   //                                                                                                                                                 260   //
    ┃    225   // The keyboard's `tab` and the mouse's click on the next tab's label must land in                                                                 261   // The keyboard's `tab` and the mouse's click on the next tab's label must land in
    ┃    226   // the same place: one switch, one answer.                                                                                                         262   // the same place: one switch, one answer.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui/ -run 'Click' 2>&1 | tail -8
    ┃
    ┃  ok    github.com/lthiagol/maitre/internal/ui  12.177s
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/ui/mouse.go internal/ui/click_test.go && git commit -q -m "feat(ui): a click on the detail pane is enter
    ┃  …
    ┃
    ┃  ac8df77 feat(ui): a click on the detail pane is enter
    ┃

       + Thought · 908ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
    ┃  p='internal/ui/mouse.go'…
    ┃
    ┃  (3 earlier lines)
    ┃  146-    }
    ┃  147:  case "click":
    ┃  148-    if step.arg != "" {
    ┃  149-      return nil, false
    ┃  150-    }
    ┃  151-    return tea.MouseClickMsg{X: x, Y: y, Button: tea.MouseLeft}, true
    ┃

    ┃
    ┃  ← Edit internal/ui/click_test.go
    ┃
    ┃    255       t.Errorf("the click on the open detail left the console in %v", m.mode)                                                                        255       t.Errorf("the click on the open detail left the console in %v", m.mode)
    ┃    256     }                                                                                                                                                256     }
    ┃    257   }                                                                                                                                                  257   }
    ┃    258                                                                                                                                                      258
    ┃                                                                                                                                                             259 + // stepAt is the scripted spelling of a click at a point.
    ┃                                                                                                                                                             260 + func stepAt(verb string, x, y int) string {
    ┃                                                                                                                                                             261 +   return fmt.Sprintf("%s@(%d,%d)", verb, x, y)
    ┃                                                                                                                                                             262 + }
    ┃                                                                                                                                                             263 +
    ┃                                                                                                                                                             264 + // TestTheClickIsDrivableFromAScript covers the scripted click.
    ┃                                                                                                                                                             265 + //
    ┃                                                                                                                                                             266 + // The console is driven from a script by `maitre --keys`, and CI has no terminal,
    ┃                                                                                                                                                             267 + // so the click has to be drivable the way the keys are.
    ┃                                                                                                                                                             268 + func TestTheClickIsDrivableFromAScript(t *testing.T) {
    ┃                                                                                                                                                             269 +   m := newModel(t)
    ┃                                                                                                                                                             270 +   m.width = 120
    ┃                                                                                                                                                             271 +   m.height = 12
    ┃                                                                                                                                                             272 +   m.View()
    ┃                                                                                                                                                             273 +
    ┃                                                                                                                                                             274 +   // The Backlog label, at the position the last drawn frame says it is at.
    ┃                                                                                                                                                             275 +   x, y := m.hits.tabs[int(TabBacklog)].centre()
    ┃                                                                                                                                                             276 +   if _, err := Drive(m, []string{stepAt("click", x, y)}); err != nil {
    ┃                                                                                                                                                             277 +     t.Fatal(err)
    ┃                                                                                                                                                             278 +   }
    ┃                                                                                                                                                             279 +   if m.tab != TabBacklog {
    ┃                                                                                                                                                             280 +     t.Errorf("the scripted click left the console on %q, want %q", m.tab.Label(), TabBacklog.Label())
    ┃                                                                                                                                                             281 +   }
    ┃                                                                                                                                                             282 + }
    ┃                                                                                                                                                             283 +
    ┃                                                                                                                                                             284 + // TestTheScriptedClickTakesAPosition covers a click on a row.
    ┃                                                                                                                                                             285 + func TestTheScriptedClickTakesAPosition(t *testing.T) {
    ┃                                                                                                                                                             286 +   m := newModel(t)
    ┃                                                                                                                                                             287 +   m.width = 120
    ┃                                                                                                                                                             288 +   m.height = 12
    ┃                                                                                                                                                             289 +   m.View()
    ┃                                                                                                                                                             290 +
    ┃                                                                                                                                                             291 +   row, ok := rowFor(m, 1)
    ┃                                                                                                                                                             292 +   if !ok {
    ┃                                                                                                                                                             293 +     t.Fatal("the map holds no row with an item")
    ┃                                                                                                                                                             294 +   }
    ┃                                                                                                                                                             295 +   if _, err := Drive(m, []string{stepAt("click", 5, row)}); err != nil {
    ┃                                                                                                                                                             296 +     t.Fatal(err)
    ┃                                                                                                                                                             297 +   }
    ┃                                                                                                                                                             298 +   if m.selRecipes != 1 {
    ┃                                                                                                                                                             299 +     t.Errorf("the scripted click left the cursor on %d, want 1", m.selRecipes)
    ┃                                                                                                                                                             300 +   }
    ┃                                                                                                                                                             301 + }
    ┃                                                                                                                                                             302 +
    ┃                                                                                                                                                             303 + // TestTheClickSurvivesTheSplitter covers the command line.
    ┃                                                                                                                                                             304 + //
    ┃                                                                                                                                                             305 + // The steps are separated by commas, and the position is bracketed — so the comma
    ┃                                                                                                                                                             306 + // inside it is not a separator.
    ┃                                                                                                                                                             307 + func TestTheClickSurvivesTheSplitter(t *testing.T) {
    ┃                                                                                                                                                             308 +   steps := SplitSteps("click@(12,5),tab")
    ┃                                                                                                                                                             309 +   want := []string{"click@(12,5)", "tab"}
    ┃                                                                                                                                                             310 +   if len(steps) != len(want) {
    ┃                                                                                                                                                             311 +     t.Fatalf("the steps are %q, want %q", steps, want)
    ┃                                                                                                                                                             312 +   }
    ┃                                                                                                                                                             313 +   for i := range want {
    ┃                                                                                                                                                             314 +     if steps[i] != want[i] {
    ┃                                                                                                                                                             315 +       t.Errorf("step %d is %q, want %q", i, steps[i], want[i])
    ┃                                                                                                                                                             316 +     }
    ┃                                                                                                                                                             317 +   }
    ┃                                                                                                                                                             318 + }
    ┃                                                                                                                                                             319 +
    ┃                                                                                                                                                             320 + // TestAnUnknownClickStepIsRefused covers the driver's refusal.
    ┃                                                                                                                                                             321 + func TestAnUnknownClickStepIsRefused(t *testing.T) {
    ┃                                                                                                                                                             322 +   m := newModel(t)
    ┃                                                                                                                                                             323 +   if _, err := Drive(m, []string{"click:sideways"}); err == nil {
    ┃                                                                                                                                                             324 +     t.Fatal("expected an unknown click to be refused")
    ┃                                                                                                                                                             325 +   }
    ┃                                                                                                                                                             326 +   if _, err := Drive(m, []string{"click@(1"}); err == nil {
    ┃                                                                                                                                                             327 +     t.Fatal("expected a malformed position to be refused")
    ┃                                                                                                                                                             328 +   }
    ┃                                                                                                                                                             329 +   if _, err := Drive(m, []string{"click@1,5"}); err == nil {
    ┃                                                                                                                                                             330 +     t.Fatal("expected an unbracketed position to be refused")
    ┃                                                                                                                                                             331 +   }
    ┃                                                                                                                                                             332 + }
    ┃                                                                                                                                                             333 +
    ┃                                                                                                                                                             334 + // TestTheClickLeavesAnOpenPromptAlone checks the prompt keeps the selection.
    ┃                                                                                                                                                             335 + //
    ┃                                                                                                                                                             336 + // While the prompt is open every input belongs to it — the console must not move
    ┃                                                                                                                                                             337 + // the selection out from under a half-typed answer.
    ┃                                                                                                                                                             338 + func TestTheClickLeavesAnOpenPromptAlone(t *testing.T) {
    ┃                                                                                                                                                             339 +   m := newModel(t)
    ┃                                                                                                                                                             340 +   m.width = 200
    ┃                                                                                                                                                             341 +   m.height = 12
    ┃                                                                                                                                                             342 +   m.View()
    ┃                                                                                                                                                             343 +
    ┃                                                                                                                                                             344 +   m.openReview("RCP02")
    ┃                                                                                                                                                             345 +   if m.prompt == nil {
    ┃                                                                                                                                                             346 +     t.Fatal("the prompt should be open")
    ┃                                                                                                                                                             347 +   }
    ┃                                                                                                                                                             348 +   before := m.selRecipes
    ┃                                                                                                                                                             349 +   beforeTab := m.tab
    ┃                                                                                                                                                             350 +
    ┃                                                                                                                                                             351 +   row, ok := rowFor(m, 1)
    ┃                                                                                                                                                             352 +   if !ok {
    ┃                                                                                                                                                             353 +     t.Fatal("the map holds no row with an item")
    ┃                                                                                                                                                             354 +   }
    ┃                                                                                                                                                             355 +   click(t, m, 5, row)
    ┃                                                                                                                                                             356 +   if m.selRecipes != before {
    ┃                                                                                                                                                             357 +     t.Errorf("the click moved the selection to %d while the prompt was open", m.selRecipes)
    ┃                                                                                                                                                             358 +   }
    ┃                                                                                                                                                             359 +
    ┃                                                                                                                                                             360 +   x, y := m.hits.tabs[int(TabBoard)].centre()
    ┃                                                                                                                                                             361 +   click(t, m, x, y)
    ┃                                                                                                                                                             362 +   if m.tab != beforeTab {
    ┃                                                                                                                                                             363 +     t.Errorf("the click went to %q while the prompt was open", m.tab.Label())
    ┃                                                                                                                                                             364 +   }
    ┃                                                                                                                                                             365 + }
    ┃                                                                                                                                                             366 +
    ┃    259   // TestTheClickOnATabIsTheKeyboardsSwitch checks the click makes the keyboard's move.                                                              367   // TestTheClickOnATabIsTheKeyboardsSwitch checks the click makes the keyboard's move.
    ┃    260   //                                                                                                                                                 368   //
    ┃    261   // The keyboard's `tab` and the mouse's click on the next tab's label must land in                                                                 369   // The keyboard's `tab` and the mouse's click on the next tab's label must land in
    ┃    262   // the same place: one switch, one answer.                                                                                                         370   // the same place: one switch, one answer.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui/ -run 'Click|Unknown|Mouse' 2>&1 | tail -8
    ┃
    ┃  ok    github.com/lthiagol/maitre/internal/ui  21.375s
    ┃

       + Thought · 631ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && rm -rf /tmp/kb-click && go run ./cmd/maitre-cli demo --dir /tmp/kb-click >/dev/null && cd /tmp/kb-click && go run github.com/lthiagol/maitre/cmd/maitre --render --keys "click@(2,2)" --plain 2>&1 | head -5
    ┃
    ┃  no required module provides package github.com/lthiagol/maitre/cmd/maitre: go.mod file not found in current directory or any parent directory; see 'go help modules'
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go build -o /tmp/maitre ./cmd/maitre && cd /tmp/kb-click && /tmp/maitre --render --keys "click@(2,2)" --plain 2>&1 | head -6; echo "---- click on the Backlog label (column 12)"; cd /tmp/kb-click && /tmp/maitre --render --keys "click@(12,2)" --plain 2>&1 | head -6
    ┃
    ┃  (7 earlier lines)
    ┃   ⋈ maitre · kb-click
    ┃  ────────────────────────────────────────────────────────────────────────────────────────────────────
    ┃    Recipes  ▸ Backlog    Path    Board    Log    Sous-Chef    Settings
    ┃   6 items
    ┃   IDEAS                                                                                          4
    ┃  ──────────────────────────────────────────────────────────────────────────────────────────────────
    ┃

       + Thought · 546ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/ui/mouse.go internal/ui/click_test.go && git commit -q -m "feat(ui): the click is drivable from a script
    ┃  …
    ┃
    ┃  588cf5f feat(ui): the click is drivable from a script
    ┃

       Now the docs:

       → Explored: 1 read

    ┃
    ┃  ← Edit docs/05-human-interface.md
    ┃
    ┃    263   | `r` | reload |                                                                                                                                   263   | `r` | reload |
    ┃    264   | `q` | quit |                                                                                                                                     264   | `q` | quit |
    ┃    265                                                                                                                                                      265
    ┃    266   **The mouse is on.** The console asks the terminal for cell-motion reporting, so the                                                               266   **The mouse is on.** The console asks the terminal for cell-motion reporting, so the
    ┃    267 - wheel reaches it. The wheel drives the same scroll the keyboard drives.                                                                            267 + wheel and the click reach it. The wheel drives the same scroll the keyboard drives,
    ┃                                                                                                                                                             268 + and **a click is a move**: it makes the move the key makes, and adds no state of its
    ┃                                                                                                                                                             269 + own.
    ┃    268                                                                                                                                                      270
    ┃    269   | Gesture | Where | What |                                                                                                                         271   | Gesture | Where | What |
    ┃    270   |---|---|---|                                                                                                                                      272   |---|---|---|
    ┃    271   | wheel | over the detail pane | scroll the detail |                                                                                               273   | wheel | over the detail pane | scroll the detail |
    ┃    272   | wheel | anywhere else, **while the list is on screen** | scroll the list, keeping the selection visible |                                        274   | wheel | anywhere else, **while the list is on screen** | scroll the list, keeping the selection visible |
    ┃                                                                                                                                                             275 + | click | on a tab label | switch to that tab — the keyboard's `tab` |
    ┃                                                                                                                                                             276 + | click | on a list row | select it — the keyboard's `↑` `↓` |
    ┃                                                                                                                                                             277 + | click | on the detail pane | open the detail — the keyboard's `enter` |
    ┃    273                                                                                                                                                      278
    ┃    274   The list is the console's one scrollable thing, so a wheel over the header scrolls the                                                             279   The list is the console's one scrollable thing, so a wheel over the header scrolls the
    ┃    275   list rather than nothing. The **full-screen detail** has no list beside it: there the                                                              280   list rather than nothing. The **full-screen detail** has no list beside it: there the
    ┃    276   wheel over the header does nothing at all, and the wheel over the detail scrolls it.                                                               281   wheel over the header does nothing at all, and the wheel over the detail scrolls it.
    ┃    277                                                                                                                                                      282
    ┃    278   The wheel's step is **three lines**: one is too slow to be useful, and a page is too                                                               283   The wheel's step is **three lines**: one is too slow to be useful, and a page is too
    ┃    279   fast to aim with.                                                                                                                                  284   fast to aim with.
    ┃    280                                                                                                                                                      285
    ┃                                                                                                                                                             286 + **A click on a row shows that recipe, and the layout decides how.** In the **split
    ┃                                                                                                                                                             287 + layout** the detail is on screen beside the list, so the click selects and the detail
    ┃                                                                                                                                                             288 + follows the selection. In the **narrow layout** the detail is not on screen at all, so
    ┃                                                                                                                                                             289 + a click that only selected would show the chef nothing: there the click also opens
    ┃                                                                                                                                                             290 + the detail. The **Backlog** has no detail yet, so a click there selects and stops; and
    ┃                                                                                                                                                             291 + a click on a **theme** row selects it — `enter` applies the theme.
    ┃                                                                                                                                                             292 +
    ┃    281   **The mouse is drivable from a script.** `maitre --keys` takes the same steps the                                                                  293   **The mouse is drivable from a script.** `maitre --keys` takes the same steps the
    ┃    282   console takes:                                                                                                                                     294   console takes:
    ┃    283                                                                                                                                                      295
    ┃    284   ```                                                                                                                                                296   ```
    ┃    285   maitre --keys "wheel:down"         # the wheel, over the list                                                                                      297   maitre --keys "wheel:down"         # the wheel, over the list
    ┃    286   maitre --keys "wheel:down@(12,5)"   # the wheel, at column 12, row 5                                                                               298   maitre --keys "wheel:down@(12,5)"   # the wheel, at column 12, row 5
    ┃                                                                                                                                                             299 + maitre --keys "click@(12,2)"        # the click, at column 12, row 2
    ┃    287   ```                                                                                                                                                300   ```
    ┃    288                                                                                                                                                      301
    ┃    289   The coordinates are the terminal's own — zero-based, with (0, 0) at the upper-left                                                                 302   The coordinates are the terminal's own — zero-based, with (0, 0) at the upper-left
    ┃    290   corner. The position is **bracketed**, because the steps are separated by commas and                                                               303   corner. The position is **bracketed**, because the steps are separated by commas and
    ┃

       + Thought · 502ms

       → Explored: 1 read

       + Thought · 2.2s

    ┃
    ┃  ← Edit docs/05-human-interface.md
    ┃
    ┃    252                                                                                                                                                      252
    ┃    253   | Key | Action |                                                                                                                                   253   | Key | Action |
    ┃    254   |---|---|                                                                                                                                          254   |---|---|
    ┃    255   | `↑` `↓` / `j` `k` | move in the list — scroll the detail when it is open |                                                                       255   | `↑` `↓` / `j` `k` | move in the list — scroll the detail when it is open |
    ┃    256 - | `tab` | next recipe |                                                                                                                            256 + | `tab` | next tab |
    ┃    257   | `enter` | expand the detail to full screen; on **Settings**, apply the theme under the cursor |                                                  257   | `enter` | expand the detail to full screen; on **Settings**, apply the theme under the cursor |
    ┃    258   | `esc` | back to the split view |                                                                                                                 258   | `esc` | back to the split view |
    ┃    259   | `n` `N` | move the cursor over the recipe's open tasting notes |                                                                                 259   | `n` `N` | move the cursor over the recipe's open tasting notes |
    ┃    260   | `space` | the primary action for this recipe — sign off, review, triage, serve |                                                                 260   | `space` | the primary action for this recipe — sign off, review, triage, serve |
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/05-human-interface.md && git commit -q -m "docs: the click, in the console's interaction table
    ┃  …
    ┃
    ┃  9c716ce docs: the click, in the console's interaction table
    ┃

       + Thought · 1.8s

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃    13   | Session | Feature | What it lands | State | Brief |                                                                                               13   | Session | Feature | What it lands | State | Brief |
    ┃    14   |--:|---|---|---|---|                                                                                                                               14   |--:|---|---|---|---|
    ┃    15   | **T1** | the identity | the window title carries the project; the header's name stops being a footnote; `--project`; the Settings tab shows the   15   | **T1** | the identity | the window title carries the project; the header's name stops being a footnote; `--project`; the Settings tab shows the
    ┃         plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438, f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81, 8a6866b, cf8c05e | [tui-t1.      plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438, f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81, 8a6866b, cf8c05e | [tui-t1.
    ┃         md](../briefs/tui-t1.md) · [fix 1](../briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |                                                     md](../briefs/tui-t1.md) · [fix 1](../briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |
    ┃    16   | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | ✅ closed · 12 commits · a72ea98, 115b650,      16   | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | ✅ closed · 12 commits · a72ea98, 115b650,
    ┃         284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](../briefs/tui-t2a.md) · [fix 1](../briefs/       284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](../briefs/tui-t2a.md) · [fix 1](../briefs/
    ┃         fix-tui-t2a-1.md) · [fix 2](../briefs/fix-tui-t2a-2.md) |                                                                                                fix-tui-t2a-1.md) · [fix 2](../briefs/fix-tui-t2a-2.md) |
    ┃    17 - | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable from a script | 🔄 **briefed** · next | [tui-t2b.md](../     17 + | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable from a script | ✅ closed · 6 commits · 4dfd5a4, 187d706,
    ┃         briefs/tui-t2b.md) |                                                                                                                                     f64e6a2, ac8df77, 588cf5f, 9c716ce | [tui-t2b.md](../briefs/tui-t2b.md) |
    ┃    18 - | **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | ⏳ planned | — |                           18 + | **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | ⏳ planned · next | — |
    ┃    19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |                                                                          19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |
    ┃    20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef dissolves; D-2 | ⏳ planned | — |                                                 20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef dissolves; D-2 | ⏳ planned | — |
    ┃    21   | **T5** | the Ready tab | `core.Adoptable`; the Path becomes Ready; D-1, D-4 | ⏳ planned | — |                                                    21   | **T5** | the Ready tab | `core.Adoptable`; the Path becomes Ready; D-1, D-4 | ⏳ planned | — |
    ┃    22   | **T6** | the Board, as a visual place | the tree drawn left to right; clicking a node | ⏳ planned | — |                                          22   | **T6** | the Board, as a visual place | the tree drawn left to right; clicking a node | ⏳ planned | — |
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/tui/README.md && git commit -q -m "work: T2b closed — the click lands" && git log --oneline -7
    ┃
    ┃  909ec28 work: T2b closed — the click lands
    ┃  9c716ce docs: the click, in the console's interaction table
    ┃  588cf5f feat(ui): the click is drivable from a script
    ┃  ac8df77 feat(ui): a click on the detail pane is enter
    ┃  f64e6a2 feat(ui): a click on a row shows that recipe
    ┃  187d706 feat(ui): a click on a row selects it
    ┃  4dfd5a4 feat(ui): a click on a tab label switches the tab
    ┃

       + Thought · 83ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && diff /tmp/t2b-final/mouse.go internal/ui/mouse.go; echo "--- model"; diff /tmp/t2b-final/model.go internal/ui/model.go; echo "--- test"; diff /tmp/t2b-final/click_test.go internal/ui/click_test.go | head -60
    ┃
    ┃  (58 earlier lines)
    ┃  >     }
    ┃  87,95c49
    ┃  <   if clicked.selRecipes != keyed.selRecipes || clicked.selBacklog != keyed.selBacklog ||
    ┃  <     clicked.selPath != keyed.selPath || clicked.selBoard != keyed.selBoard ||
    ┃  <     clicked.selLog != keyed.selLog || clicked.selTheme != keyed.selTheme {
    ┃  <     t.Errorf("the click and the keyboard left different cursors: %s vs %s",
    ┃  <       selections(clicked), selections(keyed))
    ┃

       + Thought · 416ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check 2>&1 | tail -25 && git status --short
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
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && git diff --stat && echo "--- goldens untouched:" && git log --oneline -8 --name-only | grep -c golden || true
    ┃
    ┃  --- goldens untouched:
    ┃  0
    ┃

       + Thought · 346ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check-race > /tmp/t2b-race.log 2>&1; echo "exit=$?"; tail -12 /tmp/t2b-race.log
    ┃
    ┃   Background
    ┃

       + Thought · 6.5s

    ┃
    ┃  $ which herdr; herdr --help 2>&1 | head -20
    ┃
    ┃  (14 earlier lines)
    ┃         herdr api <subcommand> ...
    ┃         herdr completion <shell>
    ┃         herdr config <subcommand> ...
    ┃         herdr channel <subcommand> ...
    ┃         herdr workspace <subcommand> ...
    ┃         herdr worktree <subcommand> ...
    ┃         herdr tab <subcommand> ...
    ┃

    ┃
    ┃  $ sleep 240; tail -20 /tmp/t2b-race.log
    ┃
    ┃  (13 earlier lines)
    ┃  ==> test-race
    ┃  GO=go scripts/go-test.sh "-race -shuffle=on -count=2 " "./..."
    ┃  ?     github.com/lthiagol/maitre/cmd/maitre  [no test files]
    ┃  ?     github.com/lthiagol/maitre/cmd/maitre-cli  [no test files]
    ┃  ok    github.com/lthiagol/maitre/internal/brand  1.468s
    ┃  ok    github.com/lthiagol/maitre/internal/cli  8.505s
    ┃  ok    github.com/lthiagol/maitre/internal/core  88.020s
    ┃

    ┃
    ┃  $ sleep 200; tail -14 /tmp/t2b-race.log
    ┃
    ┃  (7 earlier lines)
    ┃  GO=go scripts/dep-budget.sh
    ┃    internal/core: 0 external modules (ceiling 0)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃

       ↳ Shell finished · cd /Users/thiago/Code/Pessoal/maitre && make check-race > /tmp/t2b-race.log 2>&1; echo "exit=$?"; tail -12 /tmp/t2b-race.log

       + Thought · 220ms

       make check and make check-race are green. Now the receipt:

    ┃
    ┃  # Wrote work/inbox/runner-tui-t2b.md
    ┃
    ┃    1 Runner: T2b done. Commits: 6 (4dfd5a4, 187d706, f64e6a2, ac8df77, 588cf5f, 9c716ce) — the state row is 909ec28. Per-AC evidence: AC-01..AC-07 stamped. Pre-notify gate: pass (`make check`, `make check-race`; the tree clean).
    ┃    2
    ┃    3 **A click on a tab** switches to that tab, and it is the keyboard's `tab` made through one function: `Model.switchTab` — the same `m.tab = …`, the same `clampSelections()`, the same status reset. The click lands on the **label**; the separator between two labels belongs to no tab and changes nothing (
    ┃      `TestAClickOnTheGapBetweenTabsDoesNothing`). `TestTheClickOnATabIsTheKeyboardsSwitch` drives the click and the key into two consoles and compares every cursor and the status.
    ┃    4
    ┃    5 **A click on a row** selects it. The index comes from the **map** (`m.hits.rowAt`), the index the tab's own selection uses — the one the render drew the row from. A header row holds no item and selects nothing.
    ┃    6
    ┃    7 | Layout | What the click does | Why |
    ┃    8 |---|---|---|
    ┃    9 | **split** (≥160 columns, the detail beside the list) | selects; the detail follows the selection | the detail is already on screen — nothing else to do |
    ┃   10 | **narrow** (the list alone) | selects **and** opens the detail | the detail is not on screen at all; a click that only selected would show the chef nothing — so it also makes the keyboard's `enter` |
    ┃   11 | **Backlog** | selects and stops | the Backlog has no detail yet — that is T3a |
    ┃   12 | **Settings** | selects the theme, and does **not** apply it | `enter` applies the theme, and applying it is a write |
    ┃   13
    ┃   14 **A click on the detail pane** is the console's `enter` on that recipe: the detail takes the screen, `esc` brings it back. In the full-screen detail it does nothing — because `enter` does nothing there. The brief's premise was false and the fix is the reason: the wheel already follows the pointer, so
    ┃      the click had nothing left to mean *for the wheel*.
    ┃   15
    ┃   16 **The click's spelling** is `click@(x,y)`, and `click` with no position (over the list, the console's default target). It is the **same grammar** as the wheel: the same parser reads it (`parseMouseStep`), and the position is bracketed, so `ui.SplitSteps` does not cut the comma inside it — `SplitSteps(
    ┃      "click@(12,5),tab")` is two steps (`TestTheClickSurvivesTheSplitter`). The command line works: `maitre --render --keys "click@(12,2)" --plain` switches to the Backlog. An unknown step is refused with a message a reader can act on (`click:sideways`, `click@(1`, `click@1,5`).
    ┃   17
    ┃   18 **The goldens: unchanged.** A click changes what the console *answers*, not what it draws; `make goldens-check` regenerated them and found no diff, and no golden file is in the slice's commits. The Settings tab's own **KEYS** list was left alone for the same reason — adding the mouse there would have
    ┃      moved the draw.
    ┃   19
    ┃   20 Notes for the next session:
    ┃   21
    ┃   22 - **The map was not extended.** The tab labels, the rows and the detail pane were all already there; the click needed no new region. T2a's `there` flag and the row's item index did the work as they stand.
    ┃   23 - **The console answers the press, not the release.** `handleMouse` now takes the `tea.MouseMsg` and switches on the concrete type (`tea.MouseWheelMsg`, `tea.MouseClickMsg`); a release would otherwise make every move happen twice.
    ┃   24 - **One answer for the tab and for the detail.** `Model.switchTab` and `Model.openDetail` are the keyboard's moves, and the click calls them — a click that set the tab and skipped the rest would be a second answer to one question. The wheel's rule was not touched.
    ┃   25 - `docs/05-human-interface.md`: the mouse table gained the three click rows, and the key table said `tab` is "next recipe" — it is the next tab; that line was false before this slice and is fixed.
    ┃   26 - No domain change, no CLI change, no dependency: `core.DomainVersion` is 7, the tab count is 7, `internal/ui` still holds 28 external modules.
    ┃   27 - Tripwires proved by mutation: dropping the narrow layout's open fails `TestAClickOnARowInTheNarrowLayoutOpensTheDetail`; dropping the row's selection fails the two row tests; emptying `switchTab` fails `TestTheClickOnATabClampsTheCursors`.
    ┃   28
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/inbox/runner-tui-t2b.md && git commit -q -m "inbox: runner reports tui-t2b" && git log --oneline -1 && git status --short
    ┃
    ┃  f61ef4b inbox: runner reports tui-t2b
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1H:p1" "Runner: T2b done. Commits: 6 (4dfd5a4, 187d706, f64e6a2, ac8df77, 588cf5f, 9c716ce). Per-AC evidence: AC-01..AC-07 stamped. Pre-notify gate: pass. A click on a tab: switches to that tab, through the keyboard's one switch — the same clam
    ┃  ping and status reset; the separator between labels belongs to no tab. A click on a row: selects the row's item from the map; in the split layout the detail follows, in the narrow layout the click also opens the detail (enter), because the detail is not on screen there; on the Backlog it selects and…
    ┃
    ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e737d6341ffeTdoEYTgykU25M8"},"agent_status":"done","completion_seq":166,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/
    ┃  thiago/Code/Pessoal/maitre","pane_id":"w1H:p1","revision":15,"screen_detection_skipped":true,"state_change_seq":166,"tab_id":"w1H:t1","terminal_id":"term_65f0cd2aa89d26","terminal_title":"OC | Project work proposal","terminal_title_stripped":"OC | Project work proposal","workspace_id":"w1H"},"type":
    ┃  "agent_prompted"}}
    ┃

       + Thought · 603ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && git log --oneline -8
    ┃
    ┃  f61ef4b inbox: runner reports tui-t2b
    ┃  909ec28 work: T2b closed — the click lands
    ┃  9c716ce docs: the click, in the console's interaction table
    ┃  588cf5f feat(ui): the click is drivable from a script
    ┃  ac8df77 feat(ui): a click on the detail pane is enter
    ┃  f64e6a2 feat(ui): a click on a row shows that recipe
    ┃  187d706 feat(ui): a click on a row selects it
    ┃  4dfd5a4 feat(ui): a click on a tab label switches the tab
    ┃

       T2b is done. The tree is clean, the gate is green, and the receipt is committed and nudged.

       What landed

       6 commits (plus the state row 909ec28 and the inbox f61ef4b): one per acceptance criterion.

       ┌────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
       │ Commit                                                                                                                                 │ What                                                                                                                                                                   │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ 4dfd5a4                                                                                                                                │ a click on a tab label switches the tab                                                                                                                                │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ 187d706                                                                                                                                │ a click on a row selects it                                                                                                                                            │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ f64e6a2                                                                                                                                │ a click on a row shows that recipe                                                                                                                                     │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ ac8df77                                                                                                                                │ a click on the detail pane is enter                                                                                                                                    │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ 588cf5f                                                                                                                                │ the click is drivable from a script                                                                                                                                    │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ 9c716ce                                                                                                                                │ the click, in docs/05                                                                                                                                                  │
       └────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

       What each click does

       ┌─────────────────────────────────────────────────────────┬───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┬───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
       │ The click                                               │ What it does                                                                                                          │ Why                                                                                                                           │
       ├─────────────────────────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ a tab label                                             │ switches to that tab, through Model.switchTab — the same function the keyboard's tab uses                             │ the separator between labels belongs to no tab                                                                                │
       ├─────────────────────────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ a row (split layout)                                    │ selects; the detail follows the selection                                                                             │ the detail is already on screen                                                                                               │
       ├─────────────────────────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ a row (narrow layout)                                   │ selects and opens the detail                                                                                          │ the detail is not on screen; a select-only click would show the chef nothing                                                  │
       ├─────────────────────────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ a row (Backlog / Settings)                              │ selects and stops                                                                                                     │ no detail yet (T3a); a theme row is selected, not applied — enter applies it, and that is a write                             │
       ├─────────────────────────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ the detail pane                                         │ the console's enter on that recipe — the detail takes the screen, esc brings it back                                  │ the wheel already follows the pointer, so the click had nothing left to mean for the wheel                                    │
       └─────────────────────────────────────────────────────────┴───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┴───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

       The spelling

       click@(x,y), and click with no position (over the list). The same grammar as the wheel — the same parser reads it, the position is bracketed, so ui.SplitSteps does not cut the comma inside it. maitre --render --keys "click@(12,2)" --plain switches tabs.

       Notes

       - The map was not extended — the tab, row and detail regions were all already there.
       - The console answers the press, not the release: handleMouse now switches on the concrete tea.MouseMsg type.
       - switchTab and openDetail are the shared moves; the wheel's rule was not touched.
       - docs/05: the key table said tab is "next recipe" — false before this slice, now "next tab".
       - Goldens: unchanged — a click changes what the console answers, not what it draws; the Settings tab's KEYS list was left alone for the same reason.
       - Tripwires proved by mutation: three mutations each turn the suite red.
       - make check and make check-race: green. core.DomainVersion 7, 7 tabs, no dependency added.

       Build · MiniMax-M3.1-Flash · 17m 42s · 163.4 tok/s

    ┃
    ┃
    ┃
    ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
    ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
    ~/Code/Pessoal/maitre:master                                                                                                                                                                                                                                                             131.3K (13%)  ctrl+p commands
```
