package ui

import (
	"strings"
	"testing"
)

// TestTheDetailPaneScrolls covers a half-built feature: the scroll offset existed and
// nothing moved it, so the bottom of a recipe — its deviations, its review thread —
// was unreachable.
func TestTheDetailPaneScrolls(t *testing.T) {
	m := newModel(t)
	m.width = 100
	m.height = 14

	// Open the detail of the first recipe.
	out, err := Drive(m, []string{"enter"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(StripANSI(out), "↑↓ scroll") {
		t.Errorf("the footer should offer to scroll in the detail:\n%s", StripANSI(out))
	}

	// The title is the first line of the detail, so it is what leaves first.
	if !strings.Contains(StripANSI(out), "RCP02") {
		t.Fatalf("the detail should open at the title:\n%s", StripANSI(out))
	}

	out, err = Drive(m, []string{"down", "down", "down"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(StripANSI(out), "RCP02  ") {
		t.Errorf("the detail should have scrolled past the title:\n%s", StripANSI(out))
	}
	if !strings.Contains(StripANSI(out), "↑↓ scroll") {
		t.Errorf("the footer should still offer to scroll:\n%s", StripANSI(out))
	}
}

// TestTheScrollIndicatorSaysWhereYouAre checks the pane does not lose its place.
func TestTheScrollIndicatorSaysWhereYouAre(t *testing.T) {
	m := newModel(t)
	m.width = 100
	m.height = 14

	if _, err := Drive(m, []string{"enter"}); err != nil {
		t.Fatal(err)
	}
	out, err := Drive(m, []string{"down", "down"})
	if err != nil {
		t.Fatal(err)
	}
	plain := StripANSI(out)
	if !strings.Contains(plain, "3/") {
		t.Errorf("the footer should say where in the recipe you are:\n%s", plain)
	}
}

// TestScrollingBackUpWorks checks the scroll is not a one-way trip.
func TestScrollingBackUpWorks(t *testing.T) {
	m := newModel(t)
	m.width = 100
	m.height = 14

	if _, err := Drive(m, []string{"enter"}); err != nil {
		t.Fatal(err)
	}
	if _, err := Drive(m, []string{"down", "down", "down", "down"}); err != nil {
		t.Fatal(err)
	}
	out, err := Drive(m, []string{"up", "up", "up", "up", "up"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(StripANSI(out), "1/") {
		t.Errorf("scrolling back up should return to the top:\n%s", StripANSI(out))
	}
}

// TestTheDetailPaneFitsTheTerminal covers a layout bug: the pane counted entries, not
// rows, so a markdown paragraph that wrapped pushed the footer off the screen.
func TestTheDetailPaneFitsTheTerminal(t *testing.T) {
	m := newModel(t)
	m.width = 100
	m.height = 20

	out, err := Drive(m, []string{"enter"})
	if err != nil {
		t.Fatal(err)
	}
	rows := strings.Split(StripANSI(out), "\n")
	if len(rows) > m.height {
		t.Errorf("the pane drew %d rows for a %d-row terminal:\n%s", len(rows), m.height, StripANSI(out))
	}
}

// TestTheSplitPaneFitsTheTerminal checks the same for the side-by-side layout.
func TestTheSplitPaneFitsTheTerminal(t *testing.T) {
	m := newModel(t)
	m.width = 200
	m.height = 20

	out, err := Drive(m, []string{"enter"})
	if err != nil {
		t.Fatal(err)
	}
	rows := strings.Split(StripANSI(out), "\n")
	if len(rows) > m.height {
		t.Errorf("the pane drew %d rows for a %d-row terminal:\n%s", len(rows), m.height, StripANSI(out))
	}
}
