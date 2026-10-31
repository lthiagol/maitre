package ui

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// TestSpaceOnTheBoardWritesNothing covers the Board's key.
//
// The footer names no space, so space writes nothing — including on RCP02,
// which is ready and which the Recipes tab signs off. The journal is the
// same bytes afterwards, and the Recipes tab still signs that recipe off.
func TestSpaceOnTheBoardWritesNothing(t *testing.T) {
	m := newModel(t)
	switchToTab(t, m, TabBoard)
	if _, err := Drive(m, []string{"down"}); err != nil {
		t.Fatal(err)
	}
	rec, ok := m.currentRecipe()
	if !ok || rec.ID != "RCP02" {
		t.Fatalf("down should land on RCP02, got %+v", rec)
	}
	if rec.Lifecycle != core.LifecycleReady {
		t.Fatalf("RCP02 is %q, want ready", rec.Lifecycle)
	}
	if move := core.ChefMoveOf(*rec); move != core.MoveSignOff {
		t.Fatalf("the Recipes tab would %q, want a sign-off", move)
	}
	lines := strings.Split(strings.TrimRight(plain(m), "\n"), "\n")
	if strings.Contains(lines[len(lines)-1], "space") {
		t.Errorf("the Board footer names a space:\n%s", lines[len(lines)-1])
	}

	before := activityChecksum(t, m.store.Root)
	if _, err := Drive(m, []string{"space"}); err != nil {
		t.Fatal(err)
	}
	if after := activityChecksum(t, m.store.Root); after != before {
		t.Fatalf("space on the Board changed activity.json (%s -> %s)", before, after)
	}
	got, err := m.ops.GetRecipe("RCP02")
	if err != nil {
		t.Fatal(err)
	}
	if got.Lifecycle != core.LifecycleReady {
		t.Fatalf("RCP02 is %q after space, want ready", got.Lifecycle)
	}

	switchToTab(t, m, TabRecipes)
	if s, ok := m.currentSummary(); !ok || s.ID != "RCP02" {
		t.Fatalf("Recipes should open on RCP02, got %+v", s)
	}
	if _, err := Drive(m, []string{"space"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(m.status, "signed off RCP02") {
		t.Fatalf("Recipes did not sign off RCP02: %q", m.status)
	}
	got, err = m.ops.GetRecipe("RCP02")
	if err != nil {
		t.Fatal(err)
	}
	if got.Lifecycle != core.LifecycleApproved {
		t.Fatalf("lifecycle = %q, want approved", got.Lifecycle)
	}
}

// TestTheBoardShowsTheTree covers the tab.
//
// The chef switches to the Board, sees the tree, moves down it, and opens a
// recipe's detail — the same keys as every other tab.
func TestTheBoardShowsTheTree(t *testing.T) {
	m := newModel(t)
	switchToTab(t, m, TabBoard)

	out := plain(m)
	if !strings.Contains(out, "Board") {
		t.Fatalf("the tab bar should name the Board:\n%s", out)
	}
	if !strings.Contains(out, "● R") {
		t.Errorf("the tree should draw its trunk:\n%s", out)
	}
	if !strings.Contains(out, "the tree ·") {
		t.Errorf("the summary should say what the tab shows:\n%s", out)
	}

	// The selection moves.
	before := m.selBoard
	if _, err := Drive(m, []string{"down"}); err != nil {
		t.Fatal(err)
	}
	if m.selBoard == before {
		t.Error("the selection should move down the tree")
	}

	// And `enter` opens the recipe's detail.
	if _, err := Drive(m, []string{"enter"}); err != nil {
		t.Fatal(err)
	}
	if m.mode != ModeDetail {
		t.Errorf("enter should open the detail, the mode is %v", m.mode)
	}
	rec, ok := m.currentRecipe()
	if !ok {
		t.Fatal("no recipe is selected")
	}
	if rec.ID != m.boardTree().recipe[m.selBoard] {
		t.Errorf("the detail shows %q, the tree's selection is %q", rec.ID, m.boardTree().recipe[m.selBoard])
	}
}

// TestTheBoardDrawsTheConnectors covers the shape on the screen.
//
// The trunk, the branches and the connectors: the tree must read as a tree.
func TestTheBoardDrawsTheConnectors(t *testing.T) {
	m := treeModel(t)
	switchToTab(t, m, TabBoard)
	out := plain(m)

	if !strings.Contains(out, "● RCP01") {
		t.Errorf("the trunk should carry the trunk marker:\n%s", out)
	}
	if !strings.Contains(out, "needs RCP01") {
		t.Errorf("the branch should name the dependency it forks at:\n%s", out)
	}
	if !strings.Contains(out, "├─") || !strings.Contains(out, "└─") {
		t.Errorf("the branches should carry the connectors:\n%s", out)
	}
	if !strings.Contains(out, "RCP04") {
		t.Errorf("the tree should hold every recipe:\n%s", out)
	}
	if !strings.Contains(out, "not in the plan") {
		t.Errorf("a dangling dependency should say so:\n%s", out)
	}
}

// TestTheBoardSaysWhenItIsEmpty covers the empty plan.
//
// An empty plan says so. A blank screen is a broken console, and an error is a lie:
// there is nothing wrong, there is simply nothing.
func TestTheBoardSaysWhenItIsEmpty(t *testing.T) {
	root := filepath.Join(t.TempDir(), core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	m.tab = TabBoard

	out := plain(m)
	if !strings.Contains(out, "The tree is empty") {
		t.Errorf("an empty plan should say so:\n%s", out)
	}
	if strings.Contains(out, "●") {
		t.Errorf("an empty plan should draw no rows:\n%s", out)
	}
}

// TestTheBoardDoesNotCallACycleEmpty covers the broken plan.
//
// "The tree is empty" means the plan has no recipes. A plan with a dependency cycle
// has recipes — Path lists them — and the tree must not say the same words
// about both.
func TestTheBoardDoesNotCallACycleEmpty(t *testing.T) {
	m := cycleModel(t)
	switchToTab(t, m, TabBoard)
	out := plain(m)

	if strings.Contains(out, "The tree is empty") {
		t.Errorf("the console must not call a plan with a cycle empty:\n%s", out)
	}
	if !strings.Contains(out, "RCP01") || !strings.Contains(out, "RCP02") {
		t.Errorf("the cycle's recipes must be drawn:\n%s", out)
	}
	if !strings.Contains(out, "dependency cycle") {
		t.Errorf("the cycle must be marked:\n%s", out)
	}
}

// cycleModel is a console over a plan whose two recipes wait on each other.
func cycleModel(t *testing.T) *Model {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})

	mk := func(title string, deps []string) string {
		got, err := ops.CreateRecipe(core.RecipeSpec{
			Title:      title,
			Outcome:    "the tree draws what it cannot reach",
			Problem:    "a cycle leaves the tree empty",
			OutOfScope: []string{"web shell", "orchestration"},
			Criteria: []core.CriterionSpec{
				{Description: "the cycle is drawn", Verification: "go test ./internal/ui -run TestTheBoardDoesNotCallACycleEmpty"},
			},
			WorkPackages: []string{"ui"},
			Steps:        []core.StepSpec{{Action: "draw the cycle", DoneWhen: "the cycle is visible", Covers: []string{"AC-01"}}},
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(deps) > 0 {
			if _, err := ops.UpdateRecipe(got.ID, core.Fragment{DependsOn: &deps}); err != nil {
				t.Fatal(err)
			}
		}
		return got.ID
	}

	r01 := mk("The first half", nil)
	r02 := mk("The second half", nil)
	if _, err := ops.UpdateRecipe(r01, core.Fragment{DependsOn: &[]string{r02}}); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.UpdateRecipe(r02, core.Fragment{DependsOn: &[]string{r01}}); err != nil {
		t.Fatal(err)
	}

	m := loadModel(t, store)
	m.height = 40
	return m
}

// treeModel is a console over a plan with a dependency tree.
func treeModel(t *testing.T) *Model {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})

	mk := func(title string, deps []string) string {
		got, err := ops.CreateRecipe(core.RecipeSpec{
			Title:      title,
			Outcome:    "the tree forks at the dependency",
			Problem:    "the flat order does not show what waits on what",
			OutOfScope: []string{"web shell", "orchestration"},
			Criteria: []core.CriterionSpec{
				{Description: "the tree draws the forks", Verification: "go test ./internal/ui -run TestTheBoard"},
			},
			WorkPackages: []string{"ui"},
			Steps:        []core.StepSpec{{Action: "draw the tree", DoneWhen: "the tree is right", Covers: []string{"AC-01"}}},
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(deps) > 0 {
			if _, err := ops.UpdateRecipe(got.ID, core.Fragment{DependsOn: &deps}); err != nil {
				t.Fatal(err)
			}
		}
		return got.ID
	}

	r01 := mk("The trunk", nil)
	mk("Waits on the trunk", []string{r01})
	mk("Waits on two", []string{r01, "R99"})
	mk("Waits on a ghost", []string{"R99"})

	m := loadModel(t, store)
	m.height = 40
	return m
}

// TestTheWideBoardDrawsLeftToRight covers the drawing at 160 columns.
//
// The trunk's column is left of its fork's. The node shows the title, the
// dependency it waits on is still readable, and an orphan and a cycle are
// still on the board. The footer is the last row, and a node that does not
// fit is truncated rather than wrapped.
func TestTheWideBoardDrawsLeftToRight(t *testing.T) {
	m := shapeModel(t)
	m.width = 160
	m.height = 40
	m.tab = TabBoard

	out := plain(m)
	lines := strings.Split(out, "\n")
	if len(lines) != m.height {
		t.Fatalf("the screen has %d rows, want %d — a wrapped node pushes the footer", len(lines), m.height)
	}
	if !strings.Contains(lines[len(lines)-1], "quit") {
		t.Errorf("the footer should be the last row, got %q", lines[len(lines)-1])
	}
	for _, line := range lines {
		if w := lipgloss.Width(line); w > m.width {
			t.Errorf("a row is %d columns, the screen is %d — it wraps:\n%s", w, m.width, line)
		}
	}
	if strings.Contains(out, "├─") || strings.Contains(out, "└─") {
		t.Errorf("at 160 columns the tree is left to right:\n%s", out)
	}

	trunk, fork := boardLine(t, lines, "The trunk"), boardLine(t, lines, "The fork")
	if trunk != fork {
		t.Fatalf("the fork should sit on the trunk's row, trunk row %d fork row %d\n%s", trunk, fork, out)
	}
	if colOf(lines[trunk], "The trunk") >= colOf(lines[fork], "The fork") {
		t.Errorf("the trunk's column should be left of the fork's:\n%s", lines[trunk])
	}
	if !strings.Contains(lines[trunk], "needs RCP01") {
		t.Errorf("the fork should still say what it waits on:\n%s", lines[trunk])
	}
	if !strings.Contains(out, "The orphan") || !strings.Contains(out, "not in the plan") {
		t.Errorf("the orphan should still be drawn:\n%s", out)
	}
	if !strings.Contains(out, "Cycle left") || !strings.Contains(out, "Cycle right") || !strings.Contains(out, "dependency cycle") {
		t.Errorf("the cycle should still be drawn:\n%s", out)
	}
	if strings.Contains(lines[len(lines)-1], "claim") {
		t.Errorf("the footer should not offer a claim:\n%s", lines[len(lines)-1])
	}

	// The keyboard still moves down the tree, and enter still opens the detail.
	if _, err := Drive(m, []string{"down", "enter"}); err != nil {
		t.Fatal(err)
	}
	if m.selBoard != 1 {
		t.Errorf("down should move to the fork, the selection is %d", m.selBoard)
	}
	if m.mode != ModeDetail {
		t.Errorf("enter should open the detail, the mode is %v", m.mode)
	}
	rec, ok := m.currentRecipe()
	if !ok || rec.ID != "RCP02" {
		t.Fatalf("the detail should be the fork, got %+v", rec)
	}

	// A title longer than the column is cut, and the footer stays put.
	long := longTitleModel(t)
	long.width = 160
	long.height = 24
	long.tab = TabBoard
	longOut := plain(long)
	longLines := strings.Split(longOut, "\n")
	if len(longLines) != long.height {
		t.Fatalf("a long title wrapped the screen into %d rows, want %d", len(longLines), long.height)
	}
	if !strings.Contains(longLines[len(longLines)-1], "quit") {
		t.Errorf("the footer should be the last row, got %q", longLines[len(longLines)-1])
	}
	for _, line := range longLines {
		if w := lipgloss.Width(line); w > long.width {
			t.Errorf("a row is %d columns, the screen is %d", w, long.width)
		}
	}
	if !strings.Contains(longOut, "…") {
		t.Errorf("a title that does not fit should be truncated:\n%s", longOut)
	}
}

// TestAClickOnABoardNodeOpensItsDetail covers the click, wide and narrow.
//
// The click selects that recipe and opens its detail — the same detail enter
// opens. The coordinates come from the hit map, and the map's box covers the
// title the render just drew.
func TestAClickOnABoardNodeOpensItsDetail(t *testing.T) {
	for _, width := range []int{100, 160} {
		t.Run(fmt.Sprint(width), func(t *testing.T) {
			m := shapeModel(t)
			m.width = width
			m.height = 40
			m.tab = TabBoard
			m.View()

			trunk, ok := hitForItem(m, 0)
			if !ok {
				t.Fatal("the render recorded no hit for the trunk")
			}
			fork, ok := hitForItem(m, 1)
			if !ok {
				t.Fatal("the render recorded no hit for the fork")
			}
			if width >= splitThreshold {
				if trunk.y0 != fork.y0 || trunk.x1 >= fork.x0 {
					t.Errorf("the fork's hit should sit to the right of the trunk's: trunk (%d,%d)-%d fork (%d,%d)-%d",
						trunk.x0, trunk.y0, trunk.x1, fork.x0, fork.y0, fork.x1)
				}
			}

			lines := strings.Split(plain(m), "\n")
			if fork.y0 < 0 || fork.y0 >= len(lines) {
				t.Fatalf("the fork's hit is on row %d, the screen has %d", fork.y0, len(lines))
			}
			col := colOf(lines[fork.y0], "The fork")
			if col < fork.x0 || col > fork.x1 {
				t.Errorf("the hit %d-%d does not cover the title at column %d on %q", fork.x0, fork.x1, col, lines[fork.y0])
			}

			x, y := fork.centre()
			if _, err := Drive(m, []string{fmt.Sprintf("click@(%d,%d)", x, y)}); err != nil {
				t.Fatal(err)
			}
			if m.selBoard != 1 {
				t.Errorf("the click left the selection on %d, want the fork", m.selBoard)
			}
			if m.mode != ModeDetail {
				t.Errorf("the click should open the detail, the mode is %v", m.mode)
			}
			rec, ok := m.currentRecipe()
			if !ok || rec.Title != "The fork" {
				got := ""
				if ok {
					got = rec.Title
				}
				t.Fatalf("the detail should be the fork, got %q", got)
			}
		})
	}
}

// hitForItem is the hit the render recorded for a board recipe.
func hitForItem(m *Model, item int) (hitRow, bool) {
	for _, r := range m.hits.rows {
		if r.item == item {
			return r, true
		}
	}
	return hitRow{}, false
}

// boardLine is the row that holds the text.
func boardLine(t *testing.T, lines []string, text string) int {
	t.Helper()
	for i, line := range lines {
		if strings.Contains(line, text) {
			return i
		}
	}
	t.Fatalf("the board does not show %q", text)
	return -1
}

// colOf is the column text starts at.
func colOf(line, text string) int {
	i := strings.Index(line, text)
	if i < 0 {
		return -1
	}
	return lipgloss.Width(line[:i])
}

// shapeModel is a console over a trunk, a fork, an orphan and a cycle.
func shapeModel(t *testing.T) *Model {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})
	mk := func(title string) string {
		got, err := ops.CreateRecipe(core.RecipeSpec{
			Title:      title,
			Outcome:    "the wide tree draws the fork to the right",
			Problem:    "the vertical tree is the narrow drawing",
			OutOfScope: []string{"web shell", "orchestration"},
			Criteria: []core.CriterionSpec{
				{Description: "the fork is to the right", Verification: "go test ./internal/ui -run TestTheWideBoard"},
			},
			WorkPackages: []string{"ui"},
			Steps:        []core.StepSpec{{Action: "draw the fork", DoneWhen: "the fork is to the right", Covers: []string{"AC-02"}}},
		})
		if err != nil {
			t.Fatal(err)
		}
		return got.ID
	}
	trunk := mk("The trunk")
	fork := mk("The fork")
	orphan := mk("The orphan")
	left := mk("Cycle left")
	right := mk("Cycle right")
	setDeps(t, ops, fork, []string{trunk})
	setDeps(t, ops, orphan, []string{"GONE"})
	setDeps(t, ops, left, []string{right})
	setDeps(t, ops, right, []string{left})

	m := loadModel(t, store)
	m.height = 40
	return m
}

// longTitleModel is a console over one recipe whose title does not fit a node.
func longTitleModel(t *testing.T) *Model {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})
	if _, err := ops.CreateRecipe(core.RecipeSpec{
		Title:      strings.Repeat("n", 200),
		Outcome:    "a node that does not fit is truncated",
		Problem:    "a wrapped node pushes the footer off the screen",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []core.CriterionSpec{
			{Description: "the node is one line", Verification: "go test ./internal/ui -run TestTheWideBoard"},
		},
		WorkPackages: []string{"ui"},
		Steps:        []core.StepSpec{{Action: "truncate the node", DoneWhen: "the footer is the last row", Covers: []string{"AC-02"}}},
	}); err != nil {
		t.Fatal(err)
	}
	return loadModel(t, store)
}

// setDeps sets a recipe's dependencies.
func setDeps(t *testing.T, ops *core.Operations, id string, deps []string) {
	t.Helper()
	if _, err := ops.UpdateRecipe(id, core.Fragment{DependsOn: &deps}); err != nil {
		t.Fatal(err)
	}
}
