package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// The filter — the console's answer to "which of these am I looking for?".
//
// It is a **console** feature, not a domain one: the plan's format does not move,
// and `core.DomainVersion` stays where it is.
//
// **One rule, in one place.** The match is a single function, and every tab reads
// it. Two implementations of "does this match" is this house's standing failure,
// and the count in the summary must come from the same computation as the list it
// counts.

// matchFilter reports whether any of the fields holds the filter text.
//
// The rule is case-insensitive and a substring. An empty filter matches everything,
// so the list is whole when no filter is on.
func matchFilter(text string, fields ...string) bool {
	if text == "" {
		return true
	}
	needle := strings.ToLower(text)
	for _, f := range fields {
		if strings.Contains(strings.ToLower(f), needle) {
			return true
		}
	}
	return false
}

// recipeFields are the fields a recipe is matched on: its id, its title, its
// intent and its scope.
//
// The intent is the outcome and the problem; the scope is what goes in the pan and
// what does not. A recipe the chef cannot find by what it *says* is a recipe the
// filter does not reach.
func recipeFields(r core.Recipe) []string {
	return []string{
		r.ID,
		r.Title,
		r.Intent.Outcome,
		r.Intent.Problem,
		strings.Join(r.Scope.InScope, " "),
		strings.Join(r.Scope.OutOfScope, " "),
	}
}

// itemFields are the fields a backlog item is matched on: its id, its title, its
// body, its scope and its tags.
//
// The tags are the only labels a backlog item carries today. A recipe has none —
// that is a domain decision taken with the CLI — and the filter does not invent one.
func itemFields(it core.BacklogItem) []string {
	return []string{
		it.ID,
		it.Title,
		it.Body,
		strings.Join(it.Scope, " "),
		strings.Join(it.Tags, " "),
	}
}

// filterApplies reports whether the filter reaches the current tab.
//
// The filter covers the tabs that list **plan items**. The Board is a tree — a
// hidden node misstates what waits on what — the Log holds journal events, and Home
// and Settings hold numbers. A tab the filter does not reach says so, rather than
// pretending.
func (m *Model) filterApplies() bool {
	switch m.tab {
	case TabRecipes, TabBacklog, TabReady:
		return true
	}
	return false
}

// openFilter opens the filter line.
//
// The filter is a **list** feature: it filters the rows the tab draws. On a tab it
// does not reach, the console says so — a key that does nothing must say so.
func (m *Model) openFilter() (tea.Model, tea.Cmd) {
	if !m.filterApplies() {
		m.status = "the filter does not apply here"
		return m, nil
	}
	if m.mode == ModeDetail {
		// The list is not on screen: there is nothing to filter.
		m.status = "the filter does not apply here — the detail is open"
		return m, nil
	}
	m.filterOpen = true
	m.refusal = ""
	m.status = ""
	return m, nil
}

// closeFilter closes the line.
//
// `enter` keeps the filter — the chef is done typing and wants the list filtered.
// `esc` clears it: the list is whole again.
func (m *Model) closeFilter(clear bool) {
	m.filterOpen = false
	if clear {
		m.setFilter("")
		return
	}
	m.clampSelections()
}

// handleFilterKey drives the open filter line.
//
// While the line is open the keys belong to it: the console must not move the
// selection out from under a half-typed filter.
func (m *Model) handleFilterKey(k tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.closeFilter(true)
		return m, nil
	case "enter":
		m.closeFilter(false)
		return m, nil
	case "backspace":
		r := []rune(m.filter)
		if len(r) > 0 {
			m.setFilter(string(r[:len(r)-1]))
		}
		return m, nil
	case "up":
		m.move(-1)
		return m, nil
	case "down":
		m.move(1)
		return m, nil
	case "tab":
		m.switchTab(Tab((int(m.tab) + 1) % len(TabNames)))
		return m, nil
	case "shift+tab":
		m.switchTab(Tab((int(m.tab) - 1 + len(TabNames)) % len(TabNames)))
		return m, nil
	}
	// The text: anything printable goes in, and the list filters as it is typed.
	if text := typedText(k); text != "" {
		m.setFilter(m.filter + text)
	}
	return m, nil
}

// setFilter applies a new filter text.
//
// The list filters **live**: there is no second keypress to apply it.
func (m *Model) setFilter(text string) {
	m.filter = text
	m.clampSelections()
}

// filterRows is how many rows the filter line takes.
//
// **One row.** The list loses one row while the line is open, and the footer stays
// the last row of the screen. A line that pushed the footer off broke the screen.
func (m *Model) filterRows() int {
	if m.filterOpen {
		return 1
	}
	return 0
}

// filterSummary is what the summary says about the filter: what it is on, and how
// many items it leaves.
//
// A filter the chef cannot see is a filter that loses items.
func (m *Model) filterSummary() string {
	if m.filter == "" || !m.filterApplies() {
		return ""
	}
	shown, total := m.filterCounts()
	return fmt.Sprintf("filter: %s · %d of %d", m.filter, shown, total)
}

// filterCounts is how many items the tab holds, and how many the filter leaves.
//
// Both come from the same computation the list comes from, so the count in the
// summary and the rows on the screen cannot disagree.
func (m *Model) filterCounts() (shown, total int) {
	switch m.tab {
	case TabRecipes:
		return len(m.shownRows(m.filter)), len(m.shownRows(""))
	case TabBacklog:
		return len(m.backlogItems(m.filter)), len(m.backlogItems(""))
	case TabReady:
		return len(m.readyItems(m.filter)), len(m.readyItems(""))
	}
	return 0, 0
}

// filterLine is the line the chef types the filter in.
//
// It is one row, whatever the filter says: a line that wrapped would push the
// footer off the bottom of the screen.
func (m *Model) filterLine() string {
	room := m.width - 10
	if room < 1 {
		room = 1
	}
	return " " + m.styAcc.Render("filter:") + " " +
		m.styFaint.Render(truncate(m.filter, room)) + m.styAcc.Render("▏")
}

// writeFilterLine writes the filter line, if it is open, and returns how many rows
// it took.
func (m *Model) writeFilterLine(b *strings.Builder) int {
	if !m.filterOpen {
		return 0
	}
	b.WriteString(m.filterLine() + "\n")
	return 1
}
