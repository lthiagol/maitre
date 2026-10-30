package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbletea/v2"
)

// KeyFor maps a key name to a key press.
//
// Names are the ones the console's help lists, plus any single character. It lets the
// console be driven from a script — see `maitre --keys`.
func KeyFor(name string) (tea.KeyPressMsg, bool) {
	switch strings.ToLower(name) {
	case "tab":
		return tea.KeyPressMsg{Code: tea.KeyTab}, true
	case "shift+tab", "backtab":
		return tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift}, true
	case "enter", "return":
		return tea.KeyPressMsg{Code: tea.KeyEnter}, true
	case "esc", "escape":
		return tea.KeyPressMsg{Code: tea.KeyEscape}, true
	case "space":
		return tea.KeyPressMsg{Code: ' '}, true
	case "up":
		return tea.KeyPressMsg{Code: tea.KeyUp}, true
	case "down":
		return tea.KeyPressMsg{Code: tea.KeyDown}, true
	case "left":
		return tea.KeyPressMsg{Code: tea.KeyLeft}, true
	case "right":
		return tea.KeyPressMsg{Code: tea.KeyRight}, true
	}
	runes := []rune(name)
	if len(runes) == 1 {
		return tea.KeyPressMsg{Code: runes[0]}, true
	}
	return tea.KeyPressMsg{}, false
}

// Drive applies a sequence of steps to the model and returns the resulting view.
//
// A step is a key name ("tab", "space", "up") or a mouse step ("wheel:down",
// "wheel:down@12,5"). It is how the console is tested without a terminal — and how a
// script drives it — so the mouse has to be drivable the same way the keys are.
//
// The frame is drawn **before** each step. A mouse event is answered against the
// frame the chef is looking at, so the hit map has to be that frame's: a step that
// clicks against a map from a frame that is not on screen clicks the wrong thing.
func Drive(m *Model, steps []string) (string, error) {
	for _, name := range steps {
		// The frame on screen, and the map that goes with it.
		_ = m.View()

		msg, err := stepMsg(m, name)
		if err != nil {
			return "", err
		}
		updated, cmd := m.Update(msg)
		m = updated.(*Model)
		if cmd != nil {
			// Run the command the step produced — a reload, for instance.
			if msg := cmd(); msg != nil {
				updated, _ = m.Update(msg)
				m = updated.(*Model)
			}
		}
	}
	return m.View().Content, nil
}

// stepMsg turns a step name into the message it stands for.
//
// A mouse step is tried first: a key name is never a mouse step, and the mouse is
// the one that needs the frame.
func stepMsg(m *Model, name string) (tea.Msg, error) {
	if msg, ok := MouseFor(name, m); ok {
		return msg, nil
	}
	if isMouseStep(name) {
		return nil, fmt.Errorf("malformed mouse step %q", name)
	}
	key, ok := KeyFor(name)
	if !ok {
		return nil, fmt.Errorf("unknown step %q", name)
	}
	return key, nil
}
