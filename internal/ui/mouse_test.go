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
		// The map answers the same question the drawing does.
		if at := m.hits.tabAt(box.x0, tabBarRow); at != i {
			t.Errorf("the point at column %d holds tab %d, want %d", box.x0, at, i)
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
		// The map answers the same question the drawing does.
		at, ok := m.hits.rowAt(0, row.y0)
		if want := row.item >= 0; ok != want {
			t.Errorf("the point on row %d holds item %d (%t), want an item: %t", row.y0, at, ok, want)
		}
		if ok && at != row.item {
			t.Errorf("the point on row %d holds item %d, want %d", row.y0, at, row.item)
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

// TestTheWheelOverTheHeaderScrollsTheList covers the fallback.
//
// The list is the console's one scrollable thing: where the pointer is not over the
// detail, the list scrolls.
func TestTheWheelOverTheHeaderScrollsTheList(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	before := m.selRecipes
	m.Update(tea.MouseWheelMsg{X: 1, Y: 0, Button: tea.MouseWheelDown})
	if m.selRecipes == before {
		t.Errorf("the wheel over the header did not scroll the list")
	}
}

// TestTheWheelOverTheFullScreenHeaderDoesNothing checks the fallback stops where the
// list is not on screen.
//
// The full-screen detail has no list beside it, so there is nothing for the wheel to
// scroll but the detail — and the header is not the detail.
func TestTheWheelOverTheFullScreenHeaderDoesNothing(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	// The detail full screen: the detail owns the body.
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m.View()

	before := m.selRecipes
	beforeScroll := m.detailScroll
	m.Update(tea.MouseWheelMsg{X: 1, Y: 0, Button: tea.MouseWheelDown})
	if m.selRecipes != before {
		t.Errorf("the wheel moved the selection to %d while the detail was open", m.selRecipes)
	}
	if m.detailScroll != beforeScroll {
		t.Errorf("the wheel scrolled the detail to %d from over the header", m.detailScroll)
	}
}

// TestTheMouseIsDrivableFromAScript covers the scripted wheel.
//
// The console is driven from a script by `maitre --keys`, and CI has no terminal, so
// the mouse has to be drivable the way the keys are.
func TestTheMouseIsDrivableFromAScript(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 12
	before := m.selRecipes

	// The model has never been drawn: the driver has to draw the frame before it
	// answers the mouse, or the wheel lands on a map that is not on screen.
	if _, err := Drive(m, []string{"wheel:down"}); err != nil {
		t.Fatal(err)
	}
	if m.selRecipes == before {
		t.Errorf("the scripted wheel did not scroll the list: still on %d", m.selRecipes)
	}
}

// TestTheScriptedMouseTakesAPosition covers the coordinates.
//
// The position says where the pointer is, and the console answers what is under it:
// the list on the left, the detail on the right.
func TestTheScriptedMouseTakesAPosition(t *testing.T) {
	m := newModel(t)
	m.width = 200
	m.height = 12
	m.View()

	// Over the list: the list scrolls.
	if _, err := Drive(m, []string{"wheel:down@(1,6)"}); err != nil {
		t.Fatal(err)
	}
	if m.selRecipes == 0 {
		t.Errorf("the wheel over the list did not scroll it")
	}
	if m.detailScroll != 0 {
		t.Errorf("the wheel over the list scrolled the detail to %d", m.detailScroll)
	}

	// Over the detail: the detail scrolls.
	if _, err := Drive(m, []string{"wheel:down@(150,6)"}); err != nil {
		t.Fatal(err)
	}
	if m.detailScroll != wheelStep {
		t.Errorf("the wheel over the detail scrolled it by %d, want %d", m.detailScroll, wheelStep)
	}
}

// TestAnUnknownMouseStepIsRefused covers the driver's refusal.
func TestAnUnknownMouseStepIsRefused(t *testing.T) {
	m := newModel(t)
	if _, err := Drive(m, []string{"wheel:sideways"}); err == nil {
		t.Fatal("expected an unknown mouse step to be refused")
	}
	if _, err := Drive(m, []string{"wheel:down@(1"}); err == nil {
		t.Fatal("expected a malformed position to be refused")
	}
	if _, err := Drive(m, []string{"wheel:down@1,5"}); err == nil {
		t.Fatal("expected an unbracketed position to be refused")
	}
}

// TestTheWheelLeavesAnOpenPromptAlone checks the prompt keeps the selection.
//
// While the prompt is open every input belongs to it — the console must not move
// the selection out from under a half-typed answer.
func TestTheWheelLeavesAnOpenPromptAlone(t *testing.T) {
	m := newModel(t)
	m.width = 200
	m.height = 12
	m.View()

	// A prompt is open on the selected recipe.
	m.openReview("RCP02")
	if m.prompt == nil {
		t.Fatal("the prompt should be open")
	}
	before := m.selRecipes
	beforeScroll := m.detailScroll

	lx, ly := m.hits.list.centre()
	m.Update(tea.MouseWheelMsg{X: lx, Y: ly, Button: tea.MouseWheelDown})
	if m.selRecipes != before {
		t.Errorf("the wheel moved the selection to %d while the prompt was open", m.selRecipes)
	}

	dx, dy := m.hits.detail.centre()
	m.Update(tea.MouseWheelMsg{X: dx, Y: dy, Button: tea.MouseWheelDown})
	if m.detailScroll != beforeScroll {
		t.Errorf("the wheel scrolled the detail to %d while the prompt was open", m.detailScroll)
	}
}

// TestTheWheelScrollsTheSplitDetail covers the split layout's detail pane.
//
// The pane scrolls in **both** layouts. It did not: the split layout drew the detail
// without looking at the offset, so the wheel moved the offset and the frame came
// out the same.
func TestTheWheelScrollsTheSplitDetail(t *testing.T) {
	m := newModel(t)
	m.width = 200
	m.height = 20
	before := StripANSI(m.View().Content)

	x, y := m.hits.detail.centre()
	m.Update(tea.MouseWheelMsg{X: x, Y: y, Button: tea.MouseWheelDown})
	if m.detailScroll != wheelStep {
		t.Fatalf("the wheel over the split detail scrolled it by %d, want %d", m.detailScroll, wheelStep)
	}

	after := StripANSI(m.View().Content)
	if before == after {
		t.Errorf("the frame did not move:\n%s", after)
	}
}

// TestThePositionFormSurvivesTheCommandLine covers the spelling on the command
// line.
//
// The steps are separated by commas, so a position has to be bracketed: the comma
// inside it must not be a separator.
func TestThePositionFormSurvivesTheCommandLine(t *testing.T) {
	steps := SplitSteps("tab,wheel:down@(12,5),space")
	want := []string{"tab", "wheel:down@(12,5)", "space"}
	if len(steps) != len(want) {
		t.Fatalf("the steps are %q, want %q", steps, want)
	}
	for i := range want {
		if steps[i] != want[i] {
			t.Errorf("step %d is %q, want %q", i, steps[i], want[i])
		}
	}
}

// TestThePositionFormIsDrivable covers the command line's spelling, end to end.
func TestThePositionFormIsDrivable(t *testing.T) {
	m := newModel(t)
	m.width = 200
	m.height = 12
	m.View()

	// The wheel, at column 150, row 6 — over the detail pane.
	if _, err := Drive(m, SplitSteps("wheel:down@(150,6)")); err != nil {
		t.Fatal(err)
	}
	if m.detailScroll != wheelStep {
		t.Errorf("the wheel at the position scrolled the detail by %d, want %d", m.detailScroll, wheelStep)
	}
}

// TestTheOriginIsNotAPhantomDetail covers the absent region.
//
// The list layout has no detail pane, and a region that was never drawn must answer
// for nothing: the zero box answered for one cell at the origin, and the wheel at
// the origin scrolled a pane that is not on screen.
func TestTheOriginIsNotAPhantomDetail(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	if !m.hits.detail.empty() {
		t.Errorf("the list layout should have no detail region: %+v", m.hits.detail)
	}
	if m.hits.overDetail(0, 0) {
		t.Errorf("the origin should not be over a detail that is not there")
	}

	// The wheel at the origin scrolls the list.
	before := m.selRecipes
	m.Update(tea.MouseWheelMsg{X: 0, Y: 0, Button: tea.MouseWheelDown})
	if m.selRecipes == before {
		t.Errorf("the wheel at the origin did not scroll the list")
	}
	if m.detailScroll != 0 {
		t.Errorf("the wheel at the origin scrolled the detail to %d", m.detailScroll)
	}
}

// TestAClearedMapAnswersNothing covers the map of a frame that was never drawn.
//
// The error view and the opening view draw no list and no detail, so the map is
// cleared — and a cleared map answers nothing.
func TestAClearedMapAnswersNothing(t *testing.T) {
	m := newModel(t)
	m.width = 120
	m.height = 30
	m.View()

	m.hits = hitMap{}
	if m.hits.overList(0, 0) {
		t.Errorf("a cleared map should have no list")
	}
	if m.hits.overDetail(0, 0) {
		t.Errorf("a cleared map should have no detail")
	}
	if !m.hits.list.empty() || !m.hits.detail.empty() {
		t.Errorf("a cleared map should hold nothing")
	}

	// And the wheel on it scrolls nothing.
	before := m.selRecipes
	m.Update(tea.MouseWheelMsg{X: 0, Y: 0, Button: tea.MouseWheelDown})
	if m.selRecipes != before {
		t.Errorf("the wheel on a cleared map moved the selection to %d", m.selRecipes)
	}
}
