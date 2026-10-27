package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// TestConcurrentCreationsTakeDistinctIDs covers the id allocation.
//
// The id was allocated before the plan lock was taken, so two creations could scan
// the directory at the same time and be handed the same id — two items with the
// same number, one of them unreachable by id.
func TestConcurrentCreationsTakeDistinctIDs(t *testing.T) {
	ops, _ := newOps(t)

	const n = 16
	var wg sync.WaitGroup
	start := make(chan struct{})
	ids := make([]string, n)
	errs := make([]error, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			item, err := ops.CreateBacklogItem(fmt.Sprintf("A hunch %d", i), "worth noticing", nil, KindIdea, nil)
			if err != nil {
				errs[i] = err
				return
			}
			ids[i] = item.ID
		}(i)
	}
	close(start)
	wg.Wait()

	seen := map[string]bool{}
	for i := 0; i < n; i++ {
		if errs[i] != nil {
			t.Fatalf("creation %d: %v", i, errs[i])
		}
		if ids[i] == "" {
			t.Fatalf("creation %d got no id", i)
		}
		if seen[ids[i]] {
			t.Fatalf("the id %q was handed out twice", ids[i])
		}
		seen[ids[i]] = true
	}
}

// TestAFailedPromotionLeavesTheSourceOpen covers the failed promotion.
//
// A promotion writes the source and the new artifact. If the target write fails, the
// source must not claim it was promoted: the artifact is the record, and a record of
// a promotion that did not happen is a lie.
func TestAFailedPromotionLeavesTheSourceOpen(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Make the target write fail: a directory where the recipe's file goes.
	recipeID, err := store.NextRecipeID()
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(store.RecipeDir(), recipeID+"-"+Slugify(item.Title)+".json")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	before, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ops.PromoteBacklogItem(item.ID); err == nil {
		t.Fatal("the promotion should fail")
	}

	// The source is as it was.
	after, err := store.LoadBacklogItem(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != StatusOpen {
		t.Fatalf("status = %q, want %q — a failed promotion leaves no status", after.Status, StatusOpen)
	}
	if after.PromotedTo != "" {
		t.Fatalf("promoted_to = %q, want empty — a failed promotion leaves no pointer", after.PromotedTo)
	}

	// And the journal knows nothing about a promotion that did not happen.
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) != len(before.Events) {
		t.Fatalf("the journal has %d events, want %d — a failed promotion writes no event",
			len(j.Events), len(before.Events))
	}
}

// TestADanglingPointerResolvesToNothing covers the reserved id.
//
// A crash between the two writes leaves a pointer to an id that does not exist. That
// id must stay reserved: the next artifact of that kind may not take it, or the
// pointer would resolve to a stranger — a trail naming the wrong dish.
func TestADanglingPointerResolvesToNothing(t *testing.T) {
	ops, store := newOps(t)

	// The state a crash leaves: an item that points at a recipe nobody wrote.
	item, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	item.Status = StatusPromoted
	item.PromotedTo = "recipe R01"
	if err := store.SaveBacklogItem(item); err != nil {
		t.Fatal(err)
	}

	// The id it points at is reserved.
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	if r.ID == "R01" {
		t.Fatalf("the referenced id %q was handed out to a new recipe", r.ID)
	}

	// And nothing carries it: the pointer resolves to nothing.
	recipes, err := store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	for _, rec := range recipes {
		if rec.ID == "R01" {
			t.Fatal("the dangling pointer now resolves to a stranger")
		}
	}
}

// TestALegacyDanglingPointerResolvesToNothing covers the old plan's pointers.
//
// A prep note written before the rename carries the same promise as one written
// after it: it points at a recipe id, and that id must stay reserved. A reader that
// drops the pointer hands the id to a stranger, and the trail resolves to the wrong
// recipe.
func TestALegacyDanglingPointerResolvesToNothing(t *testing.T) {
	ops, store := newOps(t)

	// The state a crash leaves, in the plan format before the rename.
	dir := filepath.Join(store.Root, "prep-notes")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	note := `{"id":"N01","title":"A note","intent":"first","status":"promoted","promoted_to":"recipe RCP10"}`
	if err := os.WriteFile(filepath.Join(dir, "N01-a-note.json"), []byte(note), 0o644); err != nil {
		t.Fatal(err)
	}

	// Nine recipes, so the numbering reaches the id the pointer names.
	for i := 0; i < 9; i++ {
		if _, err := ops.CreateRecipe(readySpec()); err != nil {
			t.Fatal(err)
		}
	}
	next, err := store.NextRecipeID()
	if err != nil {
		t.Fatal(err)
	}
	if next == "RCP10" {
		t.Fatalf("the next id is %q, but that is the one the old plan points at", next)
	}

	// And nothing carries it: the pointer resolves to nothing.
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	if r.ID == "RCP10" {
		t.Fatalf("the referenced id %q was handed out to a new recipe", r.ID)
	}
	recipes, err := store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	for _, rec := range recipes {
		if rec.ID == "RCP10" {
			t.Fatal("the dangling pointer now resolves to a stranger")
		}
	}
}

// TestAPromotionIsRecordedInTheJournal covers the trail.
//
// The promotion is an event like any other: the journal is the only event log, and a
// promotion nobody can read is a promotion nobody can audit.
func TestAPromotionIsRecordedInTheJournal(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	recipe, err := ops.PromoteBacklogItem(item.ID)
	if err != nil {
		t.Fatalf("promote: %v", err)
	}
	if recipe.ID != "RCP01" {
		t.Fatalf("recipe id = %q, want RCP01", recipe.ID)
	}
	after, err := store.LoadBacklogItem(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != StatusPromoted || after.PromotedTo != "recipe RCP01" {
		t.Fatalf("the item should record where it went: %+v", after)
	}
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if last := j.Events[len(j.Events)-1]; last.Kind != KindPromotion {
		t.Fatalf("the journal kind = %q, want %q", last.Kind, KindPromotion)
	}
}

// TestTwoConcurrentPromotionsCannotBothSucceed covers the one-way trail.
//
// The promotion takes the plan lock, so the second one finds the source already
// promoted — one source, one target.
func TestTwoConcurrentPromotionsCannotBothSucceed(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	errs := make([]error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, errs[i] = ops.PromoteBacklogItem(item.ID)
		}(i)
	}
	close(start)
	wg.Wait()

	wins := 0
	for _, err := range errs {
		if err == nil {
			wins++
		}
	}
	if wins != 1 {
		t.Fatalf("wins = %d, want 1: %v", wins, errs)
	}
	recipes, err := store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if len(recipes) != 1 {
		t.Fatalf("got %d recipes, want 1 — one item, one recipe", len(recipes))
	}
}
