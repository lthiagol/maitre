package ui

import (
	"fmt"
	"strings"

	"github.com/lthiagol/maitre/internal/core"
)

// Snapshot loads the plan into a model without a terminal.
//
// It is the same model the console uses, so a script can drive it — see Drive.
func Snapshot(store *core.Store, width int) (*Model, error) {
	return snapshot(store, width, New(store))
}

// SnapshotScripted loads a plan into a console whose writes are stamped as scripted.
//
// The console is the chef's surface, so a human driving it acts as the chef. A script
// driving it acts on the chef's behalf — and the journal should say which, so a
// reader can tell a human's verdict from a driven one.
func SnapshotScripted(store *core.Store, width int) (*Model, error) {
	return snapshot(store, width, NewScripted(store))
}

func snapshot(store *core.Store, width int, m *Model) (*Model, error) {
	msg := m.reload()()
	loaded, ok := msg.(reloadMsg)
	if !ok {
		return nil, fmt.Errorf("unexpected reload result")
	}
	if loaded.err != nil {
		return nil, loaded.err
	}
	m.recipes = loaded.recipes
	m.specials = loaded.specials
	m.notes = loaded.notes
	m.journal = loaded.journal
	m.sections = m.orderedSections()
	m.clampSelections()
	m.width = width
	m.height = 40
	m.ready = true
	return m, nil
}

// Render produces a static text snapshot of the console.
//
// It is the same view the console draws, without the terminal. It serves three purposes:
// sharing a snapshot in chat or a PR, checking the console in CI, and letting an agent
// read the board the way the chef sees it.
func Render(store *core.Store, width int) (string, error) {
	m, err := snapshotModel(store, width)
	if err != nil {
		return "", err
	}
	return m.renderAll(), nil
}

// RenderDetail produces a static snapshot of one recipe.
func RenderDetail(store *core.Store, id string, width int) (string, error) {
	m, err := snapshotModel(store, width)
	if err != nil {
		return "", err
	}
	for i := range m.recipes {
		if m.recipes[i].ID == id {
			m.selRecipes = i
			m.mode = ModeDetail
			return m.viewDetail(width), nil
		}
	}
	return "", fmt.Errorf("recipe %s: %w", id, core.ErrNotFound)
}

// snapshotModel loads the plan into a model without a terminal.
func snapshotModel(store *core.Store, width int) (*Model, error) {
	return Snapshot(store, width)
}

// renderAll renders every tab, one after another.
func (m *Model) renderAll() string {
	var b strings.Builder
	for _, tab := range []Tab{TabRecipes, TabBacklog, TabPath, TabChopboard, TabLog, TabExpo, TabSettings} {
		m.tab = tab
		b.WriteString(m.header() + "\n")
		b.WriteString(m.styRule.Render(strings.Repeat("─", max(10, m.width))) + "\n")
		b.WriteString(m.tabBar() + "\n")
		b.WriteString(m.summaryLine() + "\n")
		for _, l := range m.listLines() {
			b.WriteString(l.text + "\n")
		}
		b.WriteString(m.footer())
		b.WriteString("\n\n")
	}
	return b.String()
}

// StripANSI removes escape sequences and trailing spaces, for plain-text output.
func StripANSI(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(stripANSI(line), " ")
	}
	return strings.Join(lines, "\n")
}

func stripANSI(s string) string {
	var b strings.Builder
	inEsc := false
	for _, r := range s {
		switch {
		case r == 0x1b:
			inEsc = true
		case inEsc && (r == 'm' || r == 'K' || r == 'H'):
			inEsc = false
		case !inEsc:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// LaneNames lists the lanes, for CLI help.
func LaneNames() []string {
	out := make([]string, 0, len(core.AllLanes))
	for _, l := range core.AllLanes {
		out = append(out, string(l))
	}
	return out
}
