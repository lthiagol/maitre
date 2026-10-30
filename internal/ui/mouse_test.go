package ui

import (
	"errors"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"
)

// TestTheMouseIsOn covers the enabling.
//
// In bubbletea v2 the mouse mode belongs to the view, not to the program: a view
// that does not ask for the mouse gets no mouse, and the terminal drops the events
// on the floor. Every view the console draws asks for it — the console's, the
// error's and the opening one's.
func TestTheMouseIsOn(t *testing.T) {
	m := newModel(t)
	if got := m.View().MouseMode; got != tea.MouseModeCellMotion {
		t.Errorf("the console's mouse mode = %v, want cell motion", got)
	}

	errored := newModel(t)
	errored.err = errors.New("the plan is on fire")
	if got := errored.View().MouseMode; got != tea.MouseModeCellMotion {
		t.Errorf("the error view's mouse mode = %v, want cell motion", got)
	}

	opening := newModel(t)
	opening.ready = false
	if got := opening.View().MouseMode; got != tea.MouseModeCellMotion {
		t.Errorf("the opening view's mouse mode = %v, want cell motion", got)
	}
}

// TestTheHitMapKnowsWhereTheTabsAre covers the map the render leaves behind.
//
// The map is built by the render, so the boxes and the drawn row are the same
// answer: the label is where the box says it is.
func TestTheHitMapKnowsWhereTheTabsAre(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	rows := strings.Split(plain(m), "\n")

	if len(m.hits.tabs) != len(TabNames) {
		t.Fatalf("the map holds %d tabs, want %d", len(m.hits.tabs), len(TabNames))
	}
	line := []rune(rows[tabBarRow])
	for i, box := range m.hits.tabs {
		if box.y0 != tabBarRow || box.y1 != tabBarRow {
			t.Errorf("the tab %q is on rows %d..%d, want %d", TabNames[i], box.y0, box.y1, tabBarRow)
			continue
		}
		got := string(line[box.x0 : box.x1+1])
		want := "  " + TabNames[i]
		if Tab(i) == m.tab {
			want = "▸ " + TabNames[i]
		}
		if got != want {
			t.Errorf("the tab %q is drawn as %q at columns %d..%d, want %q", TabNames[i], got, box.x0, box.x1, want)
		}
	}
}

// TestTheHitMapKnowsWhereTheRowsAre covers the list's rows.
//
// A row holds an item, and the map says which: the index the tab's own selection
// uses.
func TestTheHitMapKnowsWhereTheRowsAre(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	lines := m.listLines()
	if len(m.hits.rows) != m.listAreaHeight() {
		t.Fatalf("the map holds %d rows, want %d", len(m.hits.rows), m.listAreaHeight())
	}
	for i, row := range m.hits.rows {
		idx := m.offset + i
		want := -1
		if idx < len(lines) {
			want = lines[idx].itemIdx
		}
		if row.item != want {
			t.Errorf("the row %d holds item %d, want %d", i, row.item, want)
		}
		if row.y0 != m.hits.list.y0+i {
			t.Errorf("the row %d is on row %d, want %d", i, row.y0, m.hits.list.y0+i)
		}
	}
}

// TestTheBannerPushesTheListDown covers the refusal banner.
//
// The banner sits between the summary and the list, so the list starts lower when
// it is there. The map says so, and the map is what the mouse is answered against.
func TestTheBannerPushesTheListDown(t *testing.T) {
	m := newModel(t)
	m.width = 40
	m.height = 30
	m.refusal = "the gate said no"
	m.View()

	if len(m.hits.banner) == 0 {
		t.Fatalf("the refusal should have taken rows of its own\n%s", plain(m))
	}
	first := m.hits.banner[0].y0
	if first != headerRows {
		t.Errorf("the banner starts on row %d, want %d", first, headerRows)
	}
	if want := first + len(m.hits.banner); m.hits.list.y0 != want {
		t.Errorf("the list starts on row %d, want %d", m.hits.list.y0, want)
	}
}

// TestTheMapFollowsTheFrame checks the map is rebuilt on every render.
//
// It is one map, and the render is what fills it: a tab that shows other rows shows
// other rows in the map.
func TestTheMapFollowsTheFrame(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()
	recipes := m.hits.rows

	m.tab = TabSettings
	m.View()
	settings := m.hits.rows

	if len(settings) != len(recipes) {
		t.Fatalf("the settings tab shows %d rows, want %d", len(settings), len(recipes))
	}
	lines := m.listLines()
	for i, row := range settings {
		want := -1
		if idx := m.offset + i; idx < len(lines) {
			want = lines[idx].itemIdx
		}
		if row.item != want {
			t.Errorf("the settings row %d holds item %d, want %d", i, row.item, want)
		}
	}
}

// TestTheWheelScrollsTheList covers the wheel over the list.
func TestTheWheelScrollsTheList(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	before, ok := m.currentSummary()
	if !ok {
		t.Fatal("there is no selected recipe")
	}
	x, y := m.hits.list.centre()
	m.Update(tea.MouseWheelMsg{X: x, Y: y, Button: tea.MouseWheelDown})
	after, ok := m.currentSummary()
	if !ok {
		t.Fatal("there is no selected recipe")
	}
	if after.ID == before.ID {
		t.Errorf("the wheel over the list did not scroll it: still on %s", after.ID)
	}
}

// TestTheWheelIsTheKeyboardsScroll checks the wheel drives the same scroll.
//
// Two scrolls would be two answers to one question, so the wheel and the keyboard
// must land in the same place.
func TestTheWheelIsTheKeyboardsScroll(t *testing.T) {
	wheeled := newModel(t)
	wheeled.width = 120
	wheeled.height = 12
	wheeled.View()
	x, y := wheeled.hits.list.centre()
	for i := 0; i < 2; i++ {
		wheeled.Update(tea.MouseWheelMsg{X: x, Y: y, Button: tea.MouseWheelDown})
	}
	wheeled.View()

	keyed := newModel(t)
	keyed.width = 120
	keyed.height = 12
	for i := 0; i < 2*wheelStep; i++ {
		keyed.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	}
	keyed.View()

	if wheeled.selRecipes != keyed.selRecipes {
		t.Errorf("the wheel landed on %d, the keyboard on %d", wheeled.selRecipes, keyed.selRecipes)
	}
	if wheeled.offset != keyed.offset {
		t.Errorf("the wheel scrolled to %d, the keyboard to %d", wheeled.offset, keyed.offset)
	}
	if wheeled.offset == 0 {
		t.Errorf("the list did not scroll at all: the offset is %d", wheeled.offset)
	}
}

// TestTheWheelScrollsTheDetail covers the wheel over the detail pane.
func TestTheWheelScrollsTheDetail(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	// Open the detail full screen: the detail owns the body.
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m.View()
	if m.detailScroll != 0 {
		t.Fatalf("the detail starts at the top, got %d", m.detailScroll)
	}

	x, y := m.hits.detail.centre()
	m.Update(tea.MouseWheelMsg{X: x, Y: y, Button: tea.MouseWheelDown})
	if m.detailScroll != wheelStep {
		t.Errorf("the wheel over the detail scrolled it by %d, want %d", m.detailScroll, wheelStep)
	}
	m.Update(tea.MouseWheelMsg{X: x, Y: y, Button: tea.MouseWheelUp})
	if m.detailScroll != 0 {
		t.Errorf("the wheel back up should return to the top, got %d", m.detailScroll)
	}
}

// TestTheWheelOverTheHeaderDoesNothing checks the wheel answers what is under it.
func TestTheWheelOverTheHeaderDoesNothing(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	before := m.selRecipes
	m.Update(tea.MouseWheelMsg{X: 1, Y: 0, Button: tea.MouseWheelDown})
	if m.selRecipes != before {
		t.Errorf("the wheel over the header moved the selection to %d", m.selRecipes)
	}
}
