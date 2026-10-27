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

// Drive applies a sequence of key names to the model and returns the resulting view.
//
// It is how the console is tested without a terminal: the same model, the same keys, the
// same rendering.
func Drive(m *Model, keys []string) (string, error) {
	for _, name := range keys {
		key, ok := KeyFor(name)
		if !ok {
			return "", fmt.Errorf("unknown key %q", name)
		}
		updated, cmd := m.Update(key)
		m = updated.(*Model)
		if cmd != nil {
			// Run the command the key produced — a reload, for instance.
			if msg := cmd(); msg != nil {
				updated, _ = m.Update(msg)
				m = updated.(*Model)
			}
		}
	}
	return m.View().Content, nil
}
