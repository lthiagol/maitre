package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// TestTheBoardDrawsARenamedRecipeOnce covers the plan's index, seen from the console.
//
// The file name follows the title, so a rename used to write a second file and leave
// the first: the board read both files and drew the recipe twice, for a plan with one
// recipe in it.
func TestTheBoardDrawsARenamedRecipeOnce(t *testing.T) {
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "line-cook", Session: "s-1", Role: "runner"})

	r, err := ops.CreateRecipe(core.RecipeSpec{
		Title:      "Dep",
		Outcome:    "the plan holds one artifact per id",
		Problem:    "a rename left two files behind",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []core.CriterionSpec{
			{Description: "the board draws the recipe once", Verification: "go test ./internal/ui -run TestTheBoardDrawsARenamedRecipeOnce"},
		},
		WorkPackages: []string{"ui"},
		Steps:        []core.StepSpec{{Action: "draw the board", DoneWhen: "one recipe, one row", Covers: []string{"AC-01"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	title := "Renamed"
	if _, err := ops.UpdateRecipe(r.ID, core.Fragment{Title: &title}); err != nil {
		t.Fatal(err)
	}

	m := loadModel(t, store)
	m.height = 30
	if len(m.recipes) != 1 {
		t.Fatalf("the console read %d recipes, want 1", len(m.recipes))
	}
	out := plain(m)
	if n := strings.Count(out, title); n != 1 {
		t.Fatalf("the board draws the recipe %d times, want 1:\n%s", n, out)
	}
}
