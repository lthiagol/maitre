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
