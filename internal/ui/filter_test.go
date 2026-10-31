package ui

import (
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
