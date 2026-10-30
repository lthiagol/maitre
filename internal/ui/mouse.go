package ui

import (
	"strconv"
	"strings"

	"charm.land/bubbletea/v2"
)

// wheelStep is how far the wheel scrolls.
//
// A wheel that scrolls one line is too slow to be useful, and a wheel that scrolls
// a whole page is too fast to aim with. Three is what a terminal's own wheel
// reports per notch, and it is small enough to stop on.
const wheelStep = 3

// mouseVerbs are the gestures the console answers.
var mouseVerbs = []string{"wheel", "click"}

// handleMouse answers a mouse event.
//
// The console answers the wheel and the click. The **press** is the gesture: a
// release is the same gesture coming back, and answering both would make every
// move happen twice.
func (m *Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// An open prompt takes every input. The console must not move the selection out
	// from under a half-typed answer.
	if m.prompt != nil {
		return m, nil
	}
	switch ev := msg.(type) {
	case tea.MouseWheelMsg:
		switch ev.Button {
		case tea.MouseWheelUp:
			m.scrollUnderPointer(ev.X, ev.Y, -wheelStep)
		case tea.MouseWheelDown:
			m.scrollUnderPointer(ev.X, ev.Y, wheelStep)
		}
	case tea.MouseClickMsg:
		if ev.Button == tea.MouseLeft {
			m.clickAt(ev.X, ev.Y)
		}
	}
	return m, nil
}

// clickAt answers a click.
//
// A click is a move, and the console's moves are the keyboard's: the click makes
// the move the key makes, and it adds no state of its own.
//
// What is under the pointer decides, and the map is what says — the tab bar, the
// list, the detail pane. The order matters: the tab bar is above the list, and
// the detail pane is beside it.
func (m *Model) clickAt(x, y int) {
	if tab := m.hits.tabAt(x, y); tab >= 0 {
		m.switchTab(Tab(tab))
		return
	}
	if item, ok := m.hits.rowAt(x, y); ok {
		m.selectItem(item)
		// In the narrow layout the detail is not on screen at all, and a click
		// that only selected would show the chef nothing. There the click also
		// opens the detail — the keyboard's `enter`. In the split layout the
		// detail is already on screen and follows the selection.
		if m.hits.detail.empty() {
			m.openDetail()
		}
		return
	}
}

// selectItem puts the cursor on an item of the current tab.
//
// The index is the one the **map** holds — the index the tab's own selection
// uses, the one the render drew the row from — and the move is the keyboard's:
// the same selection state, the same clamping, the same scroll into view.
func (m *Model) selectItem(index int) {
	m.status = ""
	m.moveTo(index)
}

// scrollUnderPointer scrolls whatever is under the pointer.
//
// The detail pane answers the pointer when the pointer is over it. The **list**
// answers the rest: it is the console's one scrollable thing, so a pointer over the
// header scrolls the list rather than nothing. The full-screen detail has no list on
// screen, so there the wheel over the header does nothing at all.
//
// Both scrolls are the scrolls the keyboard drives — the same offset, the same
// clamping — so a second scroll variable never creeps in.
func (m *Model) scrollUnderPointer(x, y, delta int) {
	if m.hits.overDetail(x, y) {
		m.scrollDetail(delta)
		return
	}
	if !m.hits.list.empty() {
		m.move(delta)
	}
}

// mouseStep is a parsed mouse step.
type mouseStep struct {
	verb   string
	arg    string
	x, y   int
	hasPos bool
}

// MouseFor turns a mouse step name into the mouse event it names.
//
// The spelling is `<verb>[:<arg>][@(<x>,<y>)]`:
//
//	"wheel:up"          the wheel up, over the list
//	"wheel:down@(12,5)"  the wheel down, at column 12, row 5
//	"click"             the click, over the list
//	"click@(12,5)"      the click, at column 12, row 5
//
// The coordinates are the terminal's own: zero-based, with (0,0) at the upper-left
// corner. The position is **bracketed**, because the steps are separated by commas
// and the comma inside a position must not be one.
//
// A step with no position goes over the list — the console's default target — and
// over the detail when the list is not on screen.
func MouseFor(name string, m *Model) (tea.MouseMsg, bool) {
	step, ok := parseMouseStep(name)
	if !ok {
		return nil, false
	}
	x, y := m.hits.defaultPoint()
	if step.hasPos {
		x, y = step.x, step.y
	}
	switch step.verb {
	case "wheel":
		switch step.arg {
		case "up":
			return tea.MouseWheelMsg{X: x, Y: y, Button: tea.MouseWheelUp}, true
		case "down":
			return tea.MouseWheelMsg{X: x, Y: y, Button: tea.MouseWheelDown}, true
		}
	}
	return nil, false
}

// parseMouseStep reads a mouse step name.
//
// The spelling is `<verb>[:<arg>][@(<x>,<y>)]`. A name that is not a mouse step is
// refused, so the driver can fall back to the keys.
func parseMouseStep(name string) (mouseStep, bool) {
	head := name
	var step mouseStep
	if i := strings.Index(head, "@"); i >= 0 {
		x, y, ok := parseCoords(head[i+1:])
		if !ok {
			return mouseStep{}, false
		}
		head = head[:i]
		step.x, step.y, step.hasPos = x, y, true
	}
	verb, arg, _ := strings.Cut(head, ":")
	step.verb = strings.ToLower(verb)
	step.arg = strings.ToLower(arg)
	if step.verb == "" {
		return mouseStep{}, false
	}
	return step, true
}

// parseCoords reads the "(x,y)" of a mouse step.
//
// The position is bracketed, so the comma inside it is not a step separator: the
// steps are separated by commas, and "wheel:down@(12,5)" has to survive that.
func parseCoords(s string) (int, int, bool) {
	s = strings.TrimSpace(s)
	if len(s) < 2 || !strings.HasPrefix(s, "(") || !strings.HasSuffix(s, ")") {
		return 0, 0, false
	}
	parts := strings.Split(s[1:len(s)-1], ",")
	if len(parts) != 2 {
		return 0, 0, false
	}
	x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, false
	}
	y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, false
	}
	return x, y, true
}

// isMouseStep reports whether the name names a mouse gesture at all.
//
// A name that names a gesture and is malformed is a mistake in a script, and the
// driver says so — rather than calling it an unknown key.
func isMouseStep(name string) bool {
	verb, _, _ := strings.Cut(name, ":")
	verb = strings.ToLower(verb)
	for _, v := range mouseVerbs {
		if verb == v {
			return true
		}
	}
	return false
}

// SplitSteps splits a command line of steps into steps.
//
// The steps are separated by commas, and a position is **bracketed** — so the comma
// inside it is not a separator. "tab,wheel:down@(12,5)" is two steps, and the
// position survives.
//
// The grammar lives here, next to the parser that reads it: the command line and
// the driver must not disagree about what a step is.
func SplitSteps(s string) []string {
	var out []string
	depth := 0
	start := 0
	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	return append(out, s[start:])
}
