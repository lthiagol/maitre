package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// TestTheChopboardShowsTheTree covers the tab.
//
// The chef switches to the Chopboard, sees the tree, moves down it, and opens a
// recipe's detail — the same keys as every other tab.
func TestTheChopboardShowsTheTree(t *testing.T) {
	m := newModel(t)
	switchToTab(t, m, TabChopboard)

	out := plain(m)
	if !strings.Contains(out, "Chopboard") {
		t.Fatalf("the tab bar should name the Chopboard:\n%s", out)
	}
	if !strings.Contains(out, "● R") {
		t.Errorf("the tree should draw its trunk:\n%s", out)
	}
	if !strings.Contains(out, "the tree ·") {
		t.Errorf("the summary should say what the tab shows:\n%s", out)
	}

	// The selection moves.
	before := m.selChop
	if _, err := Drive(m, []string{"down"}); err != nil {
		t.Fatal(err)
	}
	if m.selChop == before {
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
	if rec.ID != m.chopTree().recipe[m.selChop] {
		t.Errorf("the detail shows %q, the tree's selection is %q", rec.ID, m.chopTree().recipe[m.selChop])
	}
}

// TestTheChopboardDrawsTheConnectors covers the shape on the screen.
//
// The trunk, the branches and the connectors: the tree must read as a tree.
func TestTheChopboardDrawsTheConnectors(t *testing.T) {
	m := treeModel(t)
	switchToTab(t, m, TabChopboard)
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

// TestTheChopboardSaysWhenItIsEmpty covers the empty plan.
//
// An empty plan says so. A blank screen is a broken console, and an error is a lie:
// there is nothing wrong, there is simply nothing.
func TestTheChopboardSaysWhenItIsEmpty(t *testing.T) {
	root := filepath.Join(t.TempDir(), core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	m.tab = TabChopboard

	out := plain(m)
	if !strings.Contains(out, "The tree is empty") {
		t.Errorf("an empty plan should say so:\n%s", out)
	}
	if strings.Contains(out, "●") {
		t.Errorf("an empty plan should draw no rows:\n%s", out)
	}
}

// TestTheChopboardDoesNotCallACycleEmpty covers the broken plan.
//
// "The tree is empty" means the plan has no recipes. A plan with a dependency cycle
// has recipes — the Path tab lists them — and the tree must not say the same words
// about both.
func TestTheChopboardDoesNotCallACycleEmpty(t *testing.T) {
	m := cycleModel(t)
	switchToTab(t, m, TabChopboard)
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
				{Description: "the cycle is drawn", Verification: "go test ./internal/ui -run TestTheChopboardDoesNotCallACycleEmpty"},
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
				{Description: "the tree draws the forks", Verification: "go test ./internal/ui -run TestTheChopboard"},
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
