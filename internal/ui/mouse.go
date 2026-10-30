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

// handleMouse answers a mouse event.
//
// The console answers the wheel and leaves the rest for the click: a click that
// does nothing yet would look like a thing that works.
func (m *Model) handleMouse(mouse tea.Mouse) (tea.Model, tea.Cmd) {
	// An open prompt takes every input. The console must not move the selection out
	// from under a half-typed answer.
	if m.prompt != nil {
		return m, nil
	}
	switch mouse.Button {
	case tea.MouseWheelUp:
		m.scrollUnderPointer(mouse.X, mouse.Y, -wheelStep)
	case tea.MouseWheelDown:
		m.scrollUnderPointer(mouse.X, mouse.Y, wheelStep)
	}
	return m, nil
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
// The spelling is `<verb>[:<arg>][@<x>,<y>]`:
//
//	"wheel:up"          the wheel up, over the list
//	"wheel:down@12,5"   the wheel down, at column 12, row 5
//
// The coordinates are the terminal's own: zero-based, with (0,0) at the upper-left
// corner. A step with no position goes over the list — the console's default
// target — and over the detail when the list is not on screen.
func MouseFor(name string, m *Model) (tea.MouseMsg, bool) {
	step, ok := parseMouseStep(name)
	if !ok || step.verb != "wheel" {
		return nil, false
	}
	x, y := m.hits.defaultPoint()
	if step.hasPos {
		x, y = step.x, step.y
	}
	switch step.arg {
	case "up":
		return tea.MouseWheelMsg{X: x, Y: y, Button: tea.MouseWheelUp}, true
	case "down":
		return tea.MouseWheelMsg{X: x, Y: y, Button: tea.MouseWheelDown}, true
	}
	return nil, false
}

// parseMouseStep reads a mouse step name.
//
// The spelling is `<verb>[:<arg>][@<x>,<y>]`. A name that is not a mouse step is
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

// parseCoords reads the "x,y" of a mouse step.
func parseCoords(s string) (int, int, bool) {
	parts := strings.Split(s, ",")
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
