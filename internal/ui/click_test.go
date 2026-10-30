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
