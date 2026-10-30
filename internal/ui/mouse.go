package ui

import "charm.land/bubbletea/v2"

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
// The wheel over the list scrolls the list, and the wheel over the detail scrolls
// the detail. Both are the scrolls the keyboard drives — the same offset, the same
// clamping — so a second scroll variable never creeps in.
func (m *Model) scrollUnderPointer(x, y, delta int) {
	switch {
	case m.hits.overDetail(x, y):
		m.scrollDetail(delta)
	case m.hits.overList(x, y):
		m.move(delta)
	}
}
