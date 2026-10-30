package ui

import (
	"fmt"
	"testing"

	"charm.land/bubbletea/v2"
)

// click sends a click to the model at a point.
//
// The click is a press, the way the terminal sends it: the release is the same
// gesture coming back, and the console answers the press only.
func click(t *testing.T, m *Model, x, y int) {
	t.Helper()
	m.Update(tea.MouseClickMsg{X: x, Y: y, Button: tea.MouseLeft})
}

// selections is every cursor of the console, in one string.
//
// Two consoles in the same place have the same cursors: the click and the
// keyboard must leave the same state behind.
func selections(m *Model) string {
	return fmt.Sprintf("recipes=%d backlog=%d path=%d board=%d log=%d theme=%d",
		m.selRecipes, m.selBacklog, m.selPath, m.selBoard, m.selLog, m.selTheme)
}

// TestAClickOnATabLabelSwitchesTheTab covers the tab bar.
func TestAClickOnATabLabelSwitchesTheTab(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	x, y := m.hits.tabs[int(TabBacklog)].centre()
	click(t, m, x, y)
	if m.tab != TabBacklog {
		t.Errorf("the click on the %q label left the console on %q", TabBacklog.Label(), m.tab.Label())
	}
}

// rowFor is the screen row that holds an item, and whether there is one.
func rowFor(m *Model, item int) (int, bool) {
	for _, r := range m.hits.rows {
		if r.item == item {
			return r.y0, true
		}
	}
	return 0, false
}

// headerRow is the screen row of a header — a row that holds no item.
func headerRow(m *Model) (int, bool) {
	for _, r := range m.hits.rows {
		if r.item < 0 {
			return r.y0, true
		}
	}
	return 0, false
}

// TestAClickOnARowSelectsIt covers the list.
func TestAClickOnARowSelectsIt(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	row, ok := rowFor(m, 1)
	if !ok {
		t.Fatal("the map holds no row with an item")
	}
	click(t, m, 5, row)
	if m.selRecipes != 1 {
		t.Errorf("the click on the row left the cursor on %d, want 1", m.selRecipes)
	}
	// The selection is the tab's own selection: the recipe under it is the recipe
	// the row drew.
	rows := m.visibleRows()
	if got := rows[m.selRecipes].ID; got != rows[1].ID {
		t.Errorf("the cursor is on %s, want %s", got, rows[1].ID)
	}
}

// TestAClickOnAHeaderRowSelectsNothing covers the headers.
//
// A header is not a thing to select: the map says the row holds no item, and the
// click leaves the cursor where it was.
func TestAClickOnAHeaderRowSelectsNothing(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	before := m.selRecipes
	row, ok := headerRow(m)
	if !ok {
		t.Fatal("the map holds no header row")
	}
	click(t, m, 5, row)
	if m.selRecipes != before {
		t.Errorf("the click on a header moved the cursor to %d, want %d", m.selRecipes, before)
	}
}

// TestAClickOnARowInTheNarrowLayoutOpensTheDetail covers the narrow layout.
//
// The list is alone there — the detail is not on screen — so a click that only
// selected would show the chef nothing. The click also opens the detail: the
// keyboard's `enter`.
func TestAClickOnARowInTheNarrowLayoutOpensTheDetail(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	if !m.hits.detail.empty() {
		t.Fatalf("the narrow layout should have no detail on screen: %+v", m.hits.detail)
	}
	row, ok := rowFor(m, 1)
	if !ok {
		t.Fatal("the map holds no row with an item")
	}
	click(t, m, 5, row)
	if m.mode != ModeDetail {
		t.Errorf("the click in the narrow layout left the console in %v, want the detail on screen", m.mode)
	}
	rec, ok := m.currentRecipe()
	if !ok {
		t.Fatal("the detail has no recipe")
	}
	want := m.visibleRows()[1].ID
	if rec.ID != want {
		t.Errorf("the detail shows %s, want %s", rec.ID, want)
	}
}

// TestAClickOnARowInTheSplitLayoutOnlySelects covers the split layout.
//
// The detail is on screen there, beside the list, and it follows the selection —
// so the click selects, and nothing else.
func TestAClickOnARowInTheSplitLayoutOnlySelects(t *testing.T) {
	m := newModel(t)
	m.width = 200
	m.height = 30
	m.View()

	if m.hits.detail.empty() {
		t.Fatal("the split layout should have the detail on screen")
	}
	row, ok := rowFor(m, 1)
	if !ok {
		t.Fatal("the map holds no row with an item")
	}
	click(t, m, 5, row)
	if m.mode != ModeSplit {
		t.Errorf("the click in the split layout opened the detail: mode is %v", m.mode)
	}
	if m.selRecipes != 1 {
		t.Errorf("the click left the cursor on %d, want 1", m.selRecipes)
	}
	rec, ok := m.currentRecipe()
	if !ok {
		t.Fatal("the detail has no recipe")
	}
	if want := m.visibleRows()[1].ID; rec.ID != want {
		t.Errorf("the detail shows %s, want %s", rec.ID, want)
	}
}

// TestAClickOnTheBacklogSelectsAndStops covers the tab with no detail.
//
// The Backlog has no detail yet — that is T3a — so a click there selects and
// stops.
func TestAClickOnTheBacklogSelectsAndStops(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.tab = TabBacklog
	m.View()

	row, ok := rowFor(m, 1)
	if !ok {
		t.Fatal("the map holds no backlog row with an item")
	}
	click(t, m, 5, row)
	if m.selBacklog != 1 {
		t.Errorf("the click left the cursor on %d, want 1", m.selBacklog)
	}
	if m.mode != ModeSplit {
		t.Errorf("the click on the backlog opened a detail: mode is %v", m.mode)
	}
}

// TestAClickOnAThemeRowSelectsItAndDoesNotApplyIt covers the Settings tab.
//
// The theme is applied by `enter`, and applying it is a write. A click selects
// the theme under the pointer and stops there.
func TestAClickOnAThemeRowSelectsItAndDoesNotApplyIt(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.tab = TabSettings
	m.View()

	row, ok := rowFor(m, 1)
	if !ok {
		t.Fatal("the map holds no theme row")
	}
	before := m.theme.name
	click(t, m, 5, row)
	if m.selTheme != 1 {
		t.Errorf("the click left the cursor on theme %d, want 1", m.selTheme)
	}
	if m.theme.name != before {
		t.Errorf("the click applied the theme %q", m.theme.name)
	}
	if m.mode != ModeSplit {
		t.Errorf("the click on a theme opened a detail: mode is %v", m.mode)
	}
}

// TestAClickOnTheDetailPaneOpensIt covers the detail pane.
//
// A click inside the detail pane is the console's `enter` on that recipe: the
// detail takes the screen, and `esc` brings it back.
func TestAClickOnTheDetailPaneOpensIt(t *testing.T) {
	m := newModel(t)
	m.width = 200
	m.height = 30
	m.View()

	x, y := m.hits.detail.centre()
	click(t, m, x, y)
	if m.mode != ModeDetail {
		t.Errorf("the click on the detail left the console in %v, want the detail on screen", m.mode)
	}
}

// TestAClickOnTheFullScreenDetailDoesNothing covers the full-screen detail.
//
// The detail already has the screen, and `enter` does nothing there — so the
// click does nothing there too.
func TestAClickOnTheFullScreenDetailDoesNothing(t *testing.T) {
	m := newModel(t)
	m.width = 200
	m.height = 30
	m.View()
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m.View()

	x, y := m.hits.detail.centre()
	click(t, m, x, y)
	if m.mode != ModeDetail {
		t.Errorf("the click on the open detail left the console in %v", m.mode)
	}
}

// stepAt is the scripted spelling of a click at a point.
func stepAt(verb string, x, y int) string {
	return fmt.Sprintf("%s@(%d,%d)", verb, x, y)
}

// TestTheClickIsDrivableFromAScript covers the scripted click.
//
// The console is driven from a script by `maitre --keys`, and CI has no terminal,
// so the click has to be drivable the way the keys are.
func TestTheClickIsDrivableFromAScript(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 12
	m.View()

	// The Backlog label, at the position the last drawn frame says it is at.
	x, y := m.hits.tabs[int(TabBacklog)].centre()
	if _, err := Drive(m, []string{stepAt("click", x, y)}); err != nil {
		t.Fatal(err)
	}
	if m.tab != TabBacklog {
		t.Errorf("the scripted click left the console on %q, want %q", m.tab.Label(), TabBacklog.Label())
	}
}

// TestTheScriptedClickTakesAPosition covers a click on a row.
func TestTheScriptedClickTakesAPosition(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 12
	m.View()

	row, ok := rowFor(m, 1)
	if !ok {
		t.Fatal("the map holds no row with an item")
	}
	if _, err := Drive(m, []string{stepAt("click", 5, row)}); err != nil {
		t.Fatal(err)
	}
	if m.selRecipes != 1 {
		t.Errorf("the scripted click left the cursor on %d, want 1", m.selRecipes)
	}
}

// TestTheClickSurvivesTheSplitter covers the command line.
//
// The steps are separated by commas, and the position is bracketed — so the comma
// inside it is not a separator.
func TestTheClickSurvivesTheSplitter(t *testing.T) {
	steps := SplitSteps("click@(12,5),tab")
	want := []string{"click@(12,5)", "tab"}
	if len(steps) != len(want) {
		t.Fatalf("the steps are %q, want %q", steps, want)
	}
	for i := range want {
		if steps[i] != want[i] {
			t.Errorf("step %d is %q, want %q", i, steps[i], want[i])
		}
	}
}

// TestAnUnknownClickStepIsRefused covers the driver's refusal.
func TestAnUnknownClickStepIsRefused(t *testing.T) {
	m := newModel(t)
	if _, err := Drive(m, []string{"click:sideways"}); err == nil {
		t.Fatal("expected an unknown click to be refused")
	}
	if _, err := Drive(m, []string{"click@(1"}); err == nil {
		t.Fatal("expected a malformed position to be refused")
	}
	if _, err := Drive(m, []string{"click@1,5"}); err == nil {
		t.Fatal("expected an unbracketed position to be refused")
	}
}

// TestTheClickLeavesAnOpenPromptAlone checks the prompt keeps the selection.
//
// While the prompt is open every input belongs to it — the console must not move
// the selection out from under a half-typed answer.
func TestTheClickLeavesAnOpenPromptAlone(t *testing.T) {
	m := newModel(t)
	m.width = 200
	m.height = 12
	m.View()

	m.openReview("RCP02")
	if m.prompt == nil {
		t.Fatal("the prompt should be open")
	}
	before := m.selRecipes
	beforeTab := m.tab

	row, ok := rowFor(m, 1)
	if !ok {
		t.Fatal("the map holds no row with an item")
	}
	click(t, m, 5, row)
	if m.selRecipes != before {
		t.Errorf("the click moved the selection to %d while the prompt was open", m.selRecipes)
	}

	x, y := m.hits.tabs[int(TabBoard)].centre()
	click(t, m, x, y)
	if m.tab != beforeTab {
		t.Errorf("the click went to %q while the prompt was open", m.tab.Label())
	}
}

// TestTheClickOnATabIsTheKeyboardsSwitch checks the click makes the keyboard's move.
//
// The keyboard's `tab` and the mouse's click on the next tab's label must land in
// the same place: one switch, one answer.
func TestTheClickOnATabIsTheKeyboardsSwitch(t *testing.T) {
	clicked := newModel(t)
	clicked.width = 120
	clicked.height = 30
	clicked.View()
	x, y := clicked.hits.tabs[int(clicked.tab)+1].centre()
	click(t, clicked, x, y)

	keyed := newModel(t)
	keyed.width = 120
	keyed.height = 30
	keyed.Update(tea.KeyPressMsg{Code: tea.KeyTab})

	if clicked.tab != keyed.tab {
		t.Errorf("the click landed on %q, the keyboard on %q", clicked.tab.Label(), keyed.tab.Label())
	}
	if selections(clicked) != selections(keyed) {
		t.Errorf("the click and the keyboard left different cursors: %s vs %s",
			selections(clicked), selections(keyed))
	}
	if clicked.status != keyed.status {
		t.Errorf("the click left the status %q, the keyboard %q", clicked.status, keyed.status)
	}
}

// TestTheClickOnATabClampsTheCursors covers the rest of the switch.
//
// The switch is the keyboard's: it resets the status and keeps every cursor
// inside its list. A click that set the tab and skipped the rest would be a
// second answer to one question.
func TestTheClickOnATabClampsTheCursors(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.status = "signed off RCP02"
	m.selBacklog = 99
	m.View()

	x, y := m.hits.tabs[int(TabBacklog)].centre()
	click(t, m, x, y)
	if m.status != "" {
		t.Errorf("the click on a tab left the status %q, want it cleared", m.status)
	}
	if m.selBacklog >= m.backlogCount() {
		t.Errorf("the click on a tab left the cursor at %d, want it inside the %d rows",
			m.selBacklog, m.backlogCount())
	}
}

// TestAClickOnTheGapBetweenTabsDoesNothing covers the separator.
//
// The tab bar joins its parts with a separator, and the separator belongs to no
// tab: a click on it changes nothing.
func TestAClickOnTheGapBetweenTabsDoesNothing(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	before := m.tab
	gap := m.hits.tabs[0].x1 + 1
	click(t, m, gap, tabBarRow)
	if m.tab != before {
		t.Errorf("the click on the gap between the labels went to %q", m.tab.Label())
	}
}
