package ui

import (
	"errors"
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
