package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// longDetailModel builds a console over a recipe whose detail **wraps**.
//
// The pane is narrower in the split layout — the list takes a column of it — and a
// narrower pane wraps the same text into more rows. A tripwire for the bound's
// width needs a detail where the two widths disagree: with a short one they agree,
// and the tripwire cannot fail.
func longDetailModel(t *testing.T) *Model {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})

	// A paragraph long enough to wrap several times at the pane's width, and fewer
	// times at the full screen's.
	outcome := strings.Repeat(
		"The chef reads the whole recipe in the console, so the detail pane has to be "+
			"as tall as the work it shows and never stop one row short of the end. ",
		4)

	if _, err := ops.CreateRecipe(core.RecipeSpec{
		Title:      "The detail pane scrolls to its last row",
		Outcome:    outcome,
		Problem:    "The scroll's bound was counted at the wrong width.",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []core.CriterionSpec{
			{Description: "the scroll reaches the pane's last row", Verification: "go test ./internal/ui -run TestTheScrollsBoundIsThePanesWidth"},
		},
		WorkPackages: []string{"ui"},
		Steps:        []core.StepSpec{{Action: "count the bound at the pane's width", DoneWhen: "the scroll reaches the end", Covers: []string{"AC-01"}}},
	}); err != nil {
		t.Fatal(err)
	}

	m := loadModel(t, store)
	switchToTab(t, m, TabRecipes)
	m.width = 200
	m.height = 20
	return m
}

// TestTheScrollsBoundIsThePanesWidth covers the bound the scroll stops at.
//
// The bound is counted at the width the pane **draws with**, and this test pins that
// width: the pane's rows are counted here, at the pane's own width, and not
// through the console's answer — a tripwire that asks the code under test for the
// answer is not a tripwire.
//
// The fixture makes the two widths disagree, or the comparison would be vacuous: the
// bound could be counted at the full screen's width and the rows would be the same.
func TestTheScrollsBoundIsThePanesWidth(t *testing.T) {
	m := longDetailModel(t)
	m.View()

	// The split pane's width, computed here: the list's column, one column of air,
	// and the rest of the terminal.
	paneWidth := m.width - m.listColumn() - 3
	paneRows := len(detailRows(m.detailLines(paneWidth)))
	fullRows := len(detailRows(m.detailLines(m.width - 4)))
	if paneRows == fullRows {
		t.Fatalf("the fixture must make the two widths differ: %d rows at %d columns and %d at %d",
			paneRows, paneWidth, fullRows, m.width-4)
	}
	t.Logf("the pane has %d rows at %d columns; the full screen has %d at %d", paneRows, paneWidth, fullRows, m.width-4)

	// The bound is the pane's rows.
	if m.detailLineCount() != paneRows {
		t.Errorf("the bound is %d rows; the pane has %d rows at %d columns", m.detailLineCount(), paneRows, paneWidth)
	}

	// And the scroll reaches the pane's last row.
	x, y := m.hits.detail.centre()
	for i := 0; i < 500; i++ {
		m.Update(tea.MouseWheelMsg{X: x, Y: y, Button: tea.MouseWheelDown})
	}
	if want := paneRows - 1; m.detailScroll != want {
		t.Errorf("the scroll stopped at %d, want the pane's last row %d", m.detailScroll, want)
	}
}
