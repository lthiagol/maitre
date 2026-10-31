package ui

import (
	"strconv"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// The filter — the console's answer to "which of these am I looking for?".

// screenLines is the rendered frame, one row per line.
func screenLines(s string) []string {
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}

// footerRow is the row the footer is on: the last one.
//
// The footer is pinned to the bottom of the screen, and a line that pushed it off
// broke the screen.
func footerRow(s string) int {
	return len(screenLines(s)) - 1
}

// typeFilter opens the filter line and types the text, one character at a time.
func typeFilter(t *testing.T, m *Model, text string) {
	t.Helper()
	if _, err := Drive(m, []string{"/"}); err != nil {
		t.Fatal(err)
	}
	for _, r := range text {
		if _, err := Drive(m, []string{string(r)}); err != nil {
			t.Fatal(err)
		}
	}
}

// TestTheFilterLineOpens covers AC-01.
//
// `/` opens a line under the summary, and the line is one row.
func TestTheFilterLineOpens(t *testing.T) {
	m := recipesModel(t)
	if _, err := Drive(m, []string{"/"}); err != nil {
		t.Fatal(err)
	}
	out := plain(m)
	if !m.filterOpen {
		t.Fatalf("`/` should open the filter line:\n%s", out)
	}
	lines := screenLines(out)
	if len(lines) <= headerRows {
		t.Fatalf("the frame is too short to hold the line:\n%s", out)
	}
	if got := lines[headerRows]; !strings.Contains(got, "filter:") {
		t.Errorf("row %d should hold the filter line, got %q", headerRows, got)
	}
}

// TestTheListFiltersAsTheChefTypes covers AC-01.
//
// The filter is **live**: the list narrows as the chef types, with no second
// keypress to apply it.
func TestTheListFiltersAsTheChefTypes(t *testing.T) {
	m := recipesModel(t)
	before := plain(m)
	if !strings.Contains(before, "Sharpen the AC gate") {
		t.Fatalf("the seeded plan should hold the recipe:\n%s", before)
	}

	typeFilter(t, m, "glossary")
	out := plain(m)

	if !strings.Contains(out, "Draft the glossary") {
		t.Errorf("the matching recipe should be on screen:\n%s", out)
	}
	if strings.Contains(out, "Sharpen the AC gate") {
		t.Errorf("the recipes that do not match should be gone:\n%s", out)
	}
}

// TestEnterKeepsTheFilterAndEscClearsIt covers AC-01.
//
// `enter` keeps the filter and closes the line; `esc` clears the filter and closes
// it, and the list is whole again.
func TestEnterKeepsTheFilterAndEscClearsIt(t *testing.T) {
	m := recipesModel(t)
	typeFilter(t, m, "glossary")

	if _, err := Drive(m, []string{"enter"}); err != nil {
		t.Fatal(err)
	}
	if m.filterOpen {
		t.Errorf("`enter` should close the line")
	}
	if m.filter != "glossary" {
		t.Errorf("`enter` should keep the filter, got %q", m.filter)
	}
	if strings.Contains(plain(m), "Sharpen the AC gate") {
		t.Errorf("the filter should still be on:\n%s", plain(m))
	}

	if _, err := Drive(m, []string{"/"}); err != nil {
		t.Fatal(err)
	}
	if _, err := Drive(m, []string{"esc"}); err != nil {
		t.Fatal(err)
	}
	if m.filterOpen {
		t.Errorf("`esc` should close the line")
	}
	if m.filter != "" {
		t.Errorf("`esc` should clear the filter, got %q", m.filter)
	}
	out := plain(m)
	if !strings.Contains(out, "Sharpen the AC gate") {
		t.Errorf("the list should be whole again:\n%s", out)
	}
}

// TestTheFilterLineTakesOneRow covers AC-01.
//
// The line takes **one row**: the list gives up a row, and the footer stays the
// last row of the screen. A line that pushed the footer off broke the screen.
func TestTheFilterLineTakesOneRow(t *testing.T) {
	m := recipesModel(t)
	m.width = 100
	m.height = 24

	closed := plain(m)
	closedArea := m.listAreaHeight()
	if got := footerRow(closed); got != len(screenLines(closed))-1 {
		t.Fatalf("the footer is not the last row:\n%s", closed)
	}

	if _, err := Drive(m, []string{"/"}); err != nil {
		t.Fatal(err)
	}
	open := plain(m)

	if got := m.listAreaHeight(); got != closedArea-1 {
		t.Errorf("the list takes %d rows with the line open, want %d", got, closedArea-1)
	}
	if got := footerRow(open); got != len(screenLines(open))-1 {
		t.Errorf("the footer is not the last row:\n%s", open)
	}
	// The frame is the same height: the line took the row the list gave up.
	if len(screenLines(open)) != len(screenLines(closed)) {
		t.Errorf("the frame is %d rows with the line open and %d without:\n%s",
			len(screenLines(open)), len(screenLines(closed)), open)
	}
}

// TestTheMatchRuleHoldsForBothKinds covers AC-02.
//
// The rule is **one** rule: case-insensitive, a substring, over the fields the
// filter names. A recipe and a backlog item are matched the same way — two
// implementations of "does this match" is the house's standing failure.
func TestTheMatchRuleHoldsForBothKinds(t *testing.T) {
	recipe := core.Recipe{
		ID:    "RCP01",
		Title: "Sharpen the AC gate",
		Intent: core.Intent{
			Outcome: "A zero-match verification can never satisfy a gate",
			Problem: "v1 marked criteria passed on empty filters",
		},
		Scope: core.Scope{
			InScope:    []string{"the gate table", "the verifier"},
			OutOfScope: []string{"web shell"},
		},
	}
	item := core.BacklogItem{
		ID:    "BL01",
		Title: "The board needs a way to show evidence",
		Body:  "Evidence the chef cannot read is a claim, not evidence.",
		Scope: []string{"the board lane"},
		Tags:  []string{"board", "reporting"},
	}

	recipeCases := []struct {
		name   string
		needle string
		want   bool
	}{
		{"the recipe's id", "rcp01", true},
		{"the recipe's title", "SHARPEN", true},
		{"the recipe's outcome", "ZERO-MATCH", true},
		{"the recipe's problem", "empty FILTERS", true},
		{"the recipe's in-scope", "GATE TABLE", true},
		{"the recipe's out-of-scope", "web shell", true},
		{"a word the recipe does not have", "orchestration", false},
	}
	itemCases := []struct {
		name   string
		needle string
		want   bool
	}{
		{"the item's id", "bl01", true},
		{"the item's title", "BOARD", true},
		{"the item's body", "CLAIM, NOT EVIDENCE", true},
		{"the item's scope", "board lane", true},
		{"the item's tags", "REPORTING", true},
		{"a word the item does not have", "orchestration", false},
	}

	// The **same** function answers for both kinds. That is the point: two
	// implementations of "does this match" is the house's standing failure.
	for _, c := range recipeCases {
		if got := matchFilter(c.needle, recipeFields(recipe)...); got != c.want {
			t.Errorf("the recipe and %q: got %t, want %t", c.name, got, c.want)
		}
	}
	for _, c := range itemCases {
		if got := matchFilter(c.needle, itemFields(item)...); got != c.want {
			t.Errorf("the backlog item and %q: got %t, want %t", c.name, got, c.want)
		}
	}

	// An empty filter matches everything: the list is whole when no filter is on.
	if !matchFilter("", recipeFields(recipe)...) {
		t.Errorf("an empty filter should match every recipe")
	}
	if !matchFilter("", itemFields(item)...) {
		t.Errorf("an empty filter should match every backlog item")
	}
}

// TestTheFilterFindsARecipeByItsIntent covers AC-02.
//
// The filter reaches what the recipe **says**, not only what it is called.
func TestTheFilterFindsARecipeByItsIntent(t *testing.T) {
	m := recipesModel(t)
	typeFilter(t, m, "manufactured")
	out := plain(m)

	if !strings.Contains(out, "Sharpen the AC gate") {
		t.Errorf("the recipe whose intent matches should be on screen:\n%s", out)
	}
	if strings.Contains(out, "Draft the glossary") {
		t.Errorf("the recipes that do not match should be gone:\n%s", out)
	}
}

// TestTheFilterFindsABacklogItemByItsTag covers AC-02.
//
// The backlog's tags are the labels it carries, and the filter reaches them.
func TestTheFilterFindsABacklogItemByItsTag(t *testing.T) {
	m := newModel(t)
	switchToTab(t, m, TabBacklog)
	typeFilter(t, m, "trials")
	out := plain(m)

	if !strings.Contains(out, "Trying the board on a real project") {
		t.Errorf("the item whose tag matches should be on screen:\n%s", out)
	}
	if strings.Contains(out, "The chef wants a weekly digest") {
		t.Errorf("the items that do not match should be gone:\n%s", out)
	}
}

// sectionCount is the count a section header shows, or -1 when the header is not
// on screen.
func sectionCount(out, name string) int {
	for _, line := range screenLines(out) {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, name) {
			continue
		}
		fields := strings.Fields(trimmed)
		n, err := strconv.Atoi(fields[len(fields)-1])
		if err != nil {
			continue
		}
		return n
	}
	return -1
}

// TestTheSummarySaysWhatTheFilterIsOn covers AC-03.
//
// The summary says the filter is on and on what: `filter: <text> · N of M` — N
// matches of M items. A filter the chef cannot see is a filter that loses items.
func TestTheSummarySaysWhatTheFilterIsOn(t *testing.T) {
	m := recipesModel(t)
	typeFilter(t, m, "menu")
	out := plain(m)

	if !strings.Contains(out, "filter: menu · 1 of 7") {
		t.Errorf("the summary should say what the filter is on:\n%s", out)
	}
}

// TestTheSectionHeadersCountWhatIsShown covers AC-03.
//
// The header counts the items **shown**, not the section's whole count, and a
// section with no match is not drawn at all: an empty section header is noise that
// says the plan has something it does not have.
func TestTheSectionHeadersCountWhatIsShown(t *testing.T) {
	m := recipesModel(t)
	typeFilter(t, m, "menu")
	out := plain(m)

	// "Needs you" holds two recipes; only one of them matches.
	if got := sectionCount(out, "Needs you"); got != 1 {
		t.Errorf("the Needs you header says %d, want 1:\n%s", got, out)
	}
	// The sections with no match are not drawn.
	for _, name := range []string{"On the line", "Being shaped", "Served"} {
		if got := sectionCount(out, name); got != -1 {
			t.Errorf("the %s header should not be drawn, got the count %d:\n%s", name, got, out)
		}
	}
}

// TestTheFilterCountsTheBacklogToo covers AC-03.
//
// The count is the tab's own: the backlog counts its items.
func TestTheFilterCountsTheBacklogToo(t *testing.T) {
	m := newModel(t)
	switchToTab(t, m, TabBacklog)
	typeFilter(t, m, "board")
	out := plain(m)

	// Three of the six items mention the board.
	if !strings.Contains(out, "filter: board · 3 of 6") {
		t.Errorf("the summary should count the backlog's items:\n%s", out)
	}
}

// selectedRecipe is the recipe under the Recipes cursor.
func selectedRecipe(t *testing.T, m *Model) string {
	t.Helper()
	s, ok := m.currentSummary()
	if !ok {
		t.Fatalf("there is no recipe under the cursor")
	}
	return s.ID
}

// TestTheSelectionStaysOnTheSameItem covers AC-04.
//
// The selection **follows the item**: it stays on the same recipe when that recipe
// still matches, even when the filter moves it to another row.
func TestTheSelectionStaysOnTheSameItem(t *testing.T) {
	m := recipesModel(t)
	// RCP01 is the fifth row of the triage list.
	if _, err := Drive(m, []string{"down", "down", "down", "down"}); err != nil {
		t.Fatal(err)
	}
	if got := selectedRecipe(t, m); got != "RCP01" {
		t.Fatalf("the cursor should be on RCP01, got %s", got)
	}

	// "RCP0" drops RCP07, so every row below the cursor moves up — but the
	// cursor stays on the recipe it was on.
	typeFilter(t, m, "RCP0")
	if got := selectedRecipe(t, m); got != "RCP01" {
		t.Errorf("the cursor should stay on RCP01, got %s", got)
	}
}

// TestTheSelectionGoesToTheFirstMatch covers AC-04.
//
// When the selected recipe does not match any more, the cursor goes to the first
// match — not to whatever took the old row.
func TestTheSelectionGoesToTheFirstMatch(t *testing.T) {
	m := recipesModel(t)
	// RCP03 is the second row of the triage list.
	if _, err := Drive(m, []string{"down"}); err != nil {
		t.Fatal(err)
	}
	if got := selectedRecipe(t, m); got != "RCP03" {
		t.Fatalf("the cursor should be on RCP03, got %s", got)
	}

	// "Sharpen" matches RCP05 and RCP01, and the triage's order puts RCP05
	// first — "On the line" comes before "Being shaped". RCP03 matches neither,
	// so the cursor goes to that first match, and not to the row that took the
	// old place, which an index left alone would have given.
	typeFilter(t, m, "Sharpen")
	if got := selectedRecipe(t, m); got != "RCP05" {
		t.Errorf("the cursor should go to the first match, RCP05, got %s", got)
	}
}

// TestTheListScrollsToTheSelection covers AC-04.
//
// The list scrolls to the selection, so the item the cursor names is on screen.
func TestTheFilterScrollsToTheSelection(t *testing.T) {
	m := recipesModel(t)
	m.width = 100
	m.height = 14

	// The last recipe of the triage list.
	if _, err := Drive(m, []string{"G"}); err != nil {
		t.Fatal(err)
	}
	before := selectedRecipe(t, m)

	// A filter that leaves the selection alone but takes rows off the list.
	typeFilter(t, m, "RCP0")
	if got := selectedRecipe(t, m); got != before {
		t.Fatalf("the cursor should stay on %s, got %s", before, got)
	}

	// The row the cursor names is on screen.
	lines := m.listLines()
	line := -1
	for i, l := range lines {
		if l.itemIdx == m.selRecipes {
			line = i
			break
		}
	}
	if line < 0 {
		t.Fatalf("the selected row is not in the list")
	}
	if line < m.offset || line >= m.offset+m.listAreaHeight() {
		t.Errorf("the selected row is on line %d and the window is %d..%d — the list did not scroll to it",
			line, m.offset, m.offset+m.listAreaHeight()-1)
	}
}

// TestTheFilterSurvivesATabSwitch covers AC-04.
//
// The filter is the console's, not a tab's: it survives the switch and applies to
// the new tab's list.
func TestTheFilterSurvivesATabSwitch(t *testing.T) {
	m := recipesModel(t)
	typeFilter(t, m, "the")

	if _, err := Drive(m, []string{"tab"}); err != nil {
		t.Fatal(err)
	}
	if m.tab != TabBacklog {
		t.Fatalf("the console should be on the Backlog tab, got %s", m.tab.Label())
	}
	if m.filter != "the" {
		t.Errorf("the filter should survive the switch, got %q", m.filter)
	}
	out := plain(m)
	if !strings.Contains(out, "filter: the") {
		t.Errorf("the new tab's summary should say the filter is on:\n%s", out)
	}
	if len(m.backlogItems(m.filter)) == len(m.backlogItems("")) {
		t.Errorf("the filter should apply to the new tab's list")
	}
}

// withoutSummary is the frame without its summary row.
//
// The summary row is where the console puts its one-line message, so it is the one
// row that is allowed to change when the console says something.
func withoutSummary(s string) string {
	lines := screenLines(s)
	if len(lines) <= headerRows {
		return s
	}
	kept := append([]string{}, lines[:3]...)
	kept = append(kept, lines[4:]...)
	return strings.Join(kept, "\n")
}

// TestTheFilterDoesNotApplyToTheOtherTabs covers AC-05.
//
// On Board, Log, Home and Settings, `/` says the filter does not apply here — and
// nothing else changes. The Board is a tree: a hidden node misstates what waits
// on what. The Log holds events, and Home and Settings hold numbers.
//
// **A key that does nothing must say so.** The console does not promise a filter
// it will not honour.
func TestTheFilterDoesNotApplyToTheOtherTabs(t *testing.T) {
	for _, tab := range []Tab{TabHome, TabBoard, TabLog, TabSettings} {
		m := newModel(t)
		switchToTab(t, m, tab)
		before := plain(m)

		if _, err := Drive(m, []string{"/"}); err != nil {
			t.Fatal(err)
		}
		out := plain(m)

		if !strings.Contains(out, "the filter does not apply here") {
			t.Errorf("the %s tab should say the filter does not apply here:\n%s", tab.Label(), out)
		}
		if m.filterOpen {
			t.Errorf("the %s tab should not open a filter line", tab.Label())
		}
		if m.filter != "" {
			t.Errorf("the %s tab should not set a filter, got %q", tab.Label(), m.filter)
		}
		// Nothing else moved: only the summary row changed.
		if withoutSummary(out) != withoutSummary(before) {
			t.Errorf("the %s tab changed more than its summary:\n--- before ---\n%s\n--- after ---\n%s",
				tab.Label(), before, out)
		}
	}
}

// TestTheFilterDoesNotApplyToTheOpenDetail covers AC-05.
//
// The filter is a **list** thing: with the detail open there is no list on screen,
// so there is nothing to filter.
func TestTheFilterDoesNotApplyToTheOpenDetail(t *testing.T) {
	m := recipesModel(t)
	if _, err := Drive(m, []string{"enter"}); err != nil {
		t.Fatal(err)
	}
	if m.mode != ModeDetail {
		t.Fatalf("the detail should be open")
	}
	if _, err := Drive(m, []string{"/"}); err != nil {
		t.Fatal(err)
	}
	if m.filterOpen {
		t.Errorf("the open detail should not take a filter line")
	}
	if !strings.Contains(plain(m), "the filter does not apply here") {
		t.Errorf("the console should say the filter does not apply here:\n%s", plain(m))
	}
}

// TestTheFilterIsDrivenFromAScript covers AC-06.
//
// The filter is drivable, like every other key: `/` opens the line, the
// single-character steps type into it one character at a time, and `text:<string>`
// types a whole word in one step.
func TestTheFilterIsDrivenFromAScript(t *testing.T) {
	t.Run("the single characters", func(t *testing.T) {
		m := recipesModel(t)
		if _, err := Drive(m, []string{"/", "g", "l", "o", "s", "s", "a", "r", "y"}); err != nil {
			t.Fatal(err)
		}
		if m.filter != "glossary" {
			t.Errorf("the filter should read %q, got %q", "glossary", m.filter)
		}
		out := plain(m)
		if !strings.Contains(out, "Draft the glossary") {
			t.Errorf("the matching recipe should be on screen:\n%s", out)
		}
	})

	t.Run("the whole word", func(t *testing.T) {
		m := recipesModel(t)
		if _, err := Drive(m, []string{"/", "text:glossary"}); err != nil {
			t.Fatal(err)
		}
		if m.filter != "glossary" {
			t.Errorf("the filter should read %q, got %q", "glossary", m.filter)
		}
		out := plain(m)
		if !strings.Contains(out, "Draft the glossary") {
			t.Errorf("the matching recipe should be on screen:\n%s", out)
		}
		if strings.Contains(out, "Sharpen the AC gate") {
			t.Errorf("the recipes that do not match should be gone:\n%s", out)
		}
	})

	t.Run("enter keeps it and esc clears it", func(t *testing.T) {
		m := recipesModel(t)
		if _, err := Drive(m, []string{"/", "text:glossary", "enter"}); err != nil {
			t.Fatal(err)
		}
		if m.filterOpen {
			t.Errorf("`enter` should close the line")
		}
		if m.filter != "glossary" {
			t.Errorf("`enter` should keep the filter, got %q", m.filter)
		}

		if _, err := Drive(m, []string{"/", "text:glossary", "esc"}); err != nil {
			t.Fatal(err)
		}
		if m.filterOpen {
			t.Errorf("`esc` should close the line")
		}
		if m.filter != "" {
			t.Errorf("`esc` should clear the filter, got %q", m.filter)
		}
	})

	t.Run("a word with no line to go into says nothing", func(t *testing.T) {
		m := recipesModel(t)
		if _, err := Drive(m, []string{"text:glossary"}); err != nil {
			t.Fatal(err)
		}
		if m.filter != "" {
			t.Errorf("a text step with no filter line open should do nothing, got %q", m.filter)
		}
	})
}

// TestTheTextStepIsNotAKeyRun covers AC-06.
//
// A `text:` step is one word, not a run of key presses: outside the filter line it
// must not press the console's keys — "q" quits, and a script that meant to name a
// filter must not quit the console.
func TestTheTextStepIsNotAKeyRun(t *testing.T) {
	m := recipesModel(t)
	if _, err := Drive(m, []string{"text:quit"}); err != nil {
		t.Fatal(err)
	}
	if m.filter != "" {
		t.Errorf("the filter should be untouched, got %q", m.filter)
	}
	if m.tab != TabRecipes {
		t.Errorf("the tab should be untouched, got %s", m.tab.Label())
	}
}
