package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

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

// filterNotHere is what the console says where the filter does not reach.
//
// One message, one place: the tab that does not take the filter and the detail
// that leaves no list to filter say the same thing.
const filterNotHere = "the filter does not apply here"

// openFilter opens the filter line.
//
// The filter is a **list** feature: it filters the rows the tab draws. On a tab it
// does not reach, the console says so — a key that does nothing must say so.
func (m *Model) openFilter() (tea.Model, tea.Cmd) {
	if !m.filterApplies() {
		m.status = filterNotHere
		return m, nil
	}
	if m.mode == ModeDetail {
		// The list is not on screen: there is nothing to filter.
		m.status = filterNotHere + " — the detail is open"
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

// handleText types a whole word.
//
// It is the script's way of naming a filter in one step. The word goes into the
// filter when the filter line is open, and nowhere else: a word with no line to go
// into says nothing, and it must not be read as a run of key presses.
func (m *Model) handleText(t textMsg) (tea.Model, tea.Cmd) {
	if !m.filterOpen {
		return m, nil
	}
	m.setFilter(m.filter + string(t))
	return m, nil
}

// setFilter applies a new filter text and puts the selection where it belongs.
//
// The selection **follows the item**: it stays on the same one when that one still
// matches, and it goes to the first match when it does not. An index would not do
// it — the rows under the index change when the filter does, so the cursor would
// jump to whatever took the old place.
//
// The text is sanitised on the way in, like every other text the console renders:
// the filter is drawn on two rows — the line and the summary — and a paste with a
// newline in it must not push the footer off the bottom of the screen.
func (m *Model) setFilter(text string) {
	selected := m.selectedID()
	m.filter = sanitize(text)
	m.selectID(selected)
	m.clampSelections()
}

// selectedID is the id of the item under the cursor, on the current tab.
//
// It is read **before** the filter changes: after it, the cursor points at a
// different item.
func (m *Model) selectedID() string {
	switch m.tab {
	case TabRecipes:
		if s, ok := m.currentSummary(); ok {
			return s.ID
		}
	case TabBacklog:
		if it, ok := m.selectedBacklogItem(); ok {
			return it.ID
		}
	case TabReady:
		if e, ok := m.currentReadyEntry(); ok {
			return e.id
		}
	}
	return ""
}

// selectID puts the cursor on the item with that id.
//
// An item that is not there — it does not match the filter any more — leaves the
// cursor on the first match: the first row of the list is the first thing that
// matches.
func (m *Model) selectID(id string) {
	switch m.tab {
	case TabRecipes:
		m.selRecipes = 0
		for i, r := range m.shownRows(m.filter) {
			if r.ID == id {
				m.selRecipes = i
				return
			}
		}
	case TabBacklog:
		m.selBacklog = 0
		for i, it := range m.backlogItems(m.filter) {
			if it.ID == id {
				m.selBacklog = i
				return
			}
		}
	case TabReady:
		m.selReady = 0
		for i, e := range m.readyItems(m.filter) {
			if e.id == id {
				m.selReady = i
				return
			}
		}
	}
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
// The filter's words are capped to the room the row has left. The summary is one
// row, and a filter the chef typed can be longer than the terminal — a row that
// wrapped would push the footer off the bottom of the screen.
//
// A filter the chef cannot see is a filter that loses items.
func (m *Model) filterSummary() string {
	if m.filter == "" || !m.filterApplies() {
		return ""
	}
	shown, total := m.filterCounts()
	tail := fmt.Sprintf(" · %d of %d", shown, total)
	// The room is what the tab's own summary, the filter's own words, the count
	// and the status on the right leave. The four is the summary row's own rule:
	// the right side needs a column of air beside it.
	const overhead = len(" · filter: ") + 4
	room := m.width - lipgloss.Width(m.tabSummary()) - lipgloss.Width(tail) -
		lipgloss.Width(m.summaryRight()) - overhead
	if room < 1 {
		room = 1
	}
	return "filter: " + truncate(m.filter, room) + tail
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
