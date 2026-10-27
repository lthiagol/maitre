package core

import (
	"testing"
)

// treePlan is a plan with a dependency chain, a branch and a cycle.
//
//	R01  nothing waits on it, but R02, R03 and R04 do
//	R02  waits on R01
//	R03  waits on R01 and R04 — it forks at R01, the first unmet dependency
//	R04  waits on R01
//	R05  waits on R99 — a dependency that is not in the plan
//	R06  waits on nothing, and nothing waits on it
//	R07  waits on R08, and R08 waits on R07 — a dependency cycle
//	R08  waits on R07 — the other half of the cycle
//	R09  waits on R07 and R01 — on a cycle member, and on the trunk
func treePlan(t *testing.T) []Recipe {
	t.Helper()
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(spec())
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := ops.Advance(r.ID, ev); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := ops.ApproveRecipe(r.ID, "chef"); err != nil {
		t.Fatal(err)
	}

	mk := func(id, title string, deps []string) Recipe {
		got, err := ops.CreateRecipe(RecipeSpec{
			Title:      title,
			Outcome:    "o",
			Problem:    "p",
			InScope:    []string{"the tree"},
			OutOfScope: []string{"web shell", "orchestration"},
			Criteria: []CriterionSpec{
				{Description: "the tree forks at the dependency", Verification: "go test ./internal/core -run TestTheTree"},
			},
			WorkPackages: []string{"core"},
			Steps:        []StepSpec{{Action: "draw the tree", DoneWhen: "the tree is right", Covers: []string{"AC-01"}}},
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(deps) > 0 {
			if _, err := ops.UpdateRecipe(got.ID, Fragment{DependsOn: &deps}); err != nil {
				t.Fatal(err)
			}
		}
		return *got
	}

	// The recipes are created in order, so their ids are R02, R03, R04, R05.
	r02 := mk("", "Second", []string{"R01"})
	r03 := mk("", "Third", []string{"R01", "R04"})
	r04 := mk("", "Fourth", []string{"R01"})
	r05 := mk("", "Fifth", []string{"R99"})
	r06 := mk("", "Sixth", nil)
	r07 := mk("", "Seventh", []string{"R08"})
	r08 := mk("", "Eighth", []string{"R07"})
	r09 := mk("", "Ninth", []string{"R07", "R01"})

	var out []Recipe
	for _, id := range []string{"R01", r02.ID, r03.ID, r04.ID, r05.ID, r06.ID, r07.ID, r08.ID, r09.ID} {
		got, err := ops.GetRecipe(id)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, *got)
	}
	return out
}

// treeIDs are the recipes the tree draws, however deep they are.
func treeIDs(tree Tree) map[string]bool {
	ids := map[string]bool{}
	var walk func(n TreeNode)
	walk = func(n TreeNode) {
		ids[n.Summary.ID] = true
		for _, c := range n.Waiting {
			walk(c)
		}
	}
	for _, r := range tree.Roots {
		walk(r)
	}
	for _, o := range tree.Orphans {
		walk(o)
	}
	for _, s := range tree.Stuck {
		walk(s)
	}
	return ids
}

// TestACycleIsDrawn covers the recipes a cycle keeps out of the tree.
//
// Two recipes that wait on each other have no root to walk from, and neither is an
// orphan — its dependency is in the plan. The tree has to draw them anyway: the Path
// tab lists them, and a tree that drops them tells the chef a story the plan does
// not.
func TestACycleIsDrawn(t *testing.T) {
	tree := TreeOf(treePlan(t))
	ids := treeIDs(tree)

	if !ids["R07"] {
		t.Error("R07 is in a cycle and must still be drawn")
	}
	if !ids["R08"] {
		t.Error("R08 is in a cycle and must still be drawn")
	}
}

// TestARecipeWaitingOnACycleIsDrawn covers the transitive case.
//
// A recipe that waits on a cycle member is just as unreachable: the thing it forks
// under is never drawn.
func TestARecipeWaitingOnACycleIsDrawn(t *testing.T) {
	tree := TreeOf(treePlan(t))
	if !treeIDs(tree)["R09"] {
		t.Error("R09 waits on a cycle member and must still be drawn")
	}
}

// findNode is a recipe in the tree, wherever it hangs.
func findNode(tree Tree, id string) (TreeNode, bool) {
	var found TreeNode
	var ok bool
	var walk func(n TreeNode)
	walk = func(n TreeNode) {
		if n.Summary.ID == id {
			found, ok = n, true
			return
		}
		for _, c := range n.Waiting {
			walk(c)
		}
	}
	for _, r := range tree.Roots {
		walk(r)
	}
	for _, o := range tree.Orphans {
		walk(o)
	}
	for _, s := range tree.Stuck {
		walk(s)
	}
	return found, ok
}

// TestACycleIsMarked covers the warning.
//
// The plan cannot be executed as written, and the row says so. A member of the cycle
// and a recipe that waits on one are told apart: the first cannot be executed at all,
// the second can be executed once the cycle is broken.
func TestACycleIsMarked(t *testing.T) {
	tree := TreeOf(treePlan(t))

	r07, ok := findNode(tree, "R07")
	if !ok {
		t.Fatal("R07 is not in the tree")
	}
	if !r07.Cycle {
		t.Error("R07 waits on itself: it is in a cycle and must be marked")
	}
	if r07.WaitsOnCycle {
		t.Error("R07 is in the cycle itself; it does not merely wait on one")
	}

	r09, ok := findNode(tree, "R09")
	if !ok {
		t.Fatal("R09 is not in the tree")
	}
	if r09.Cycle {
		t.Error("R09 waits on a cycle member; it is not in the cycle itself")
	}
	if !r09.WaitsOnCycle {
		t.Error("R09 waits on a cycle and must be marked")
	}
}

// TestTheTreeForksAtTheDependency covers the shape.
//
// The trunk is the ready list, and every recipe that waits on something forks at its
// first unmet dependency — under the recipe it waits for, nested as deep as the
// chain goes.
func TestTheTreeForksAtTheDependency(t *testing.T) {
	tree := TreeOf(treePlan(t))

	if len(tree.Roots) != 2 {
		t.Fatalf("the tree has %d roots, want 2 (R01 and R06)", len(tree.Roots))
	}
	root := tree.Roots[0]
	if root.Summary.ID != "R01" {
		t.Fatalf("the root is %q, want R01", root.Summary.ID)
	}
	if len(root.Waiting) != 3 {
		t.Fatalf("the root has %d branches, want 3 (R02, R03 and R04 fork at R01)", len(root.Waiting))
	}
	if root.Waiting[0].Summary.ID != "R02" {
		t.Errorf("the first branch is %q, want R02", root.Waiting[0].Summary.ID)
	}
	if root.Waiting[1].Summary.ID != "R03" {
		t.Errorf("the second branch is %q, want R03", root.Waiting[1].Summary.ID)
	}
	if root.Waiting[2].Summary.ID != "R04" {
		t.Errorf("the third branch is %q, want R04", root.Waiting[2].Summary.ID)
	}
}

// TestTheTrunkIsThePathTabsReadyList covers the order.
//
// Two tabs showing the same plan must not tell two stories: the tree's roots are
// core.Path's ready entries, in the same order.
func TestTheTrunkIsThePathTabsReadyList(t *testing.T) {
	recipes := treePlan(t)
	tree := TreeOf(recipes)
	ready, _ := Path(recipes)

	if len(tree.Roots) != len(ready) {
		t.Fatalf("the tree has %d roots, the path has %d ready entries", len(tree.Roots), len(ready))
	}
	for i, e := range ready {
		if tree.Roots[i].Summary.ID != e.Summary.ID {
			t.Errorf("root %d is %q, the path's ready entry %d is %q",
				i, tree.Roots[i].Summary.ID, i, e.Summary.ID)
		}
	}
}

// TestARowSaysEveryDependencyItWaitsOn covers the text.
//
// A recipe with two unmet dependencies forks at the first — but the row must say it
// needs both, because the chef reads the row, not the tree.
func TestARowSaysEveryDependencyItWaitsOn(t *testing.T) {
	tree := TreeOf(treePlan(t))
	root := tree.Roots[0]

	var r03 *TreeNode
	for i := range root.Waiting {
		if root.Waiting[i].Summary.ID == "R03" {
			r03 = &root.Waiting[i]
		}
	}
	if r03 == nil {
		t.Fatal("R03 does not fork at R01")
	}
	if len(r03.Needs) != 2 {
		t.Fatalf("R03 waits on %v, want two dependencies", r03.Needs)
	}
	if r03.Needs[0] != "R01" || r03.Needs[1] != "R04" {
		t.Fatalf("R03 waits on %v, want R01 then R04", r03.Needs)
	}
}

// TestADanglingDependencyIsVisible covers the broken plan.
//
// A dependency that is not in the plan cannot be forked at: the recipe waits on
// something that is not there, and the tree has to say so.
func TestADanglingDependencyIsVisible(t *testing.T) {
	tree := TreeOf(treePlan(t))

	if len(tree.Orphans) != 1 {
		t.Fatalf("the tree has %d orphans, want 1", len(tree.Orphans))
	}
	orphan := tree.Orphans[0]
	if orphan.Summary.ID != "R05" {
		t.Fatalf("the orphan is %q, want R05", orphan.Summary.ID)
	}
	if len(orphan.Needs) == 0 || orphan.Needs[0] != "R99" {
		t.Fatalf("the orphan waits on %v, want R99", orphan.Needs)
	}
	if len(orphan.Missing) != 1 || orphan.Missing[0] != "R99" {
		t.Fatalf("the orphan's missing dependencies are %v, want R99", orphan.Missing)
	}
}

// TestARecipeWithNothingWaitingOnItIsOnTheTrunk covers the lonely recipe.
//
// A recipe that can go now and that nothing waits on is on the trunk, at the top
// level — not under anything.
func TestARecipeWithNothingWaitingOnItIsOnTheTrunk(t *testing.T) {
	recipes := treePlan(t)
	tree := TreeOf(recipes)

	// R06 waits on nothing and nothing waits on it.
	var r06 *TreeNode
	for i := range tree.Roots {
		if tree.Roots[i].Summary.ID == "R06" {
			r06 = &tree.Roots[i]
		}
	}
	if r06 == nil {
		t.Fatal("R06 is not on the trunk")
	}
	if len(r06.Needs) != 0 {
		t.Errorf("R06 waits on %v, want nothing", r06.Needs)
	}
	if len(r06.Waiting) != 0 {
		t.Errorf("R06 has %d branches, want none", len(r06.Waiting))
	}
}

// TestEveryRecipeInThePlanIsInTheTreeOnce covers the completeness.
//
// The tree must not lose a recipe and must not draw it twice. The plan has a
// dependency cycle in it, and the recipes that wait on one: those are the shapes that
// broke it.
func TestEveryRecipeInThePlanIsInTheTreeOnce(t *testing.T) {
	recipes := treePlan(t)
	tree := TreeOf(recipes)

	counts := map[string]int{}
	var walk func(n TreeNode)
	walk = func(n TreeNode) {
		counts[n.Summary.ID]++
		for _, c := range n.Waiting {
			walk(c)
		}
	}
	for _, r := range tree.Roots {
		walk(r)
	}
	for _, o := range tree.Orphans {
		walk(o)
	}
	for _, s := range tree.Stuck {
		walk(s)
	}
	for _, r := range recipes {
		if counts[r.ID] != 1 {
			t.Errorf("%s is in the tree %d times, want 1", r.ID, counts[r.ID])
		}
	}
}
