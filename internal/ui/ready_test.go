package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// readyModel is a console over a plan with one adoptable recipe, one that
// waits, a draft, an idea, two open small tasks, and one dismissed task.
func readyModel(t *testing.T) (*Model, *core.Store) {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})
	chef := core.NewOperations(store, core.AuditContext{Actor: "chef", Session: "s-2", Role: "chef"})

	spec := func(title string) core.RecipeSpec {
		return core.RecipeSpec{
			Title:        title,
			Outcome:      "the work is known",
			Problem:      "the tab has to tell them apart",
			OutOfScope:   []string{"web shell", "orchestration"},
			Criteria:     []core.CriterionSpec{{Description: "it is listed", Verification: "go test ./internal/ui -run Ready"}},
			WorkPackages: []string{"ui"},
			Steps:        []core.StepSpec{{Action: "list it", DoneWhen: "the row is there", Covers: []string{"AC-01"}}},
		}
	}
	free, err := ops.CreateRecipe(spec("The free recipe"))
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range []core.Event{core.EventGroom, core.EventSpecify} {
		if _, err := ops.Advance(free.ID, ev); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := chef.ApproveRecipe(free.ID, "chef"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.CreateRecipe(spec("Still being shaped")); err != nil {
		t.Fatal(err)
	}
	waiting, err := ops.CreateRecipe(spec("Waiting on the draft"))
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range []core.Event{core.EventGroom, core.EventSpecify} {
		if _, err := ops.Advance(waiting.ID, ev); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := chef.ApproveRecipe(waiting.ID, "chef"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.UpdateRecipe(waiting.ID, core.Fragment{DependsOn: &[]string{"RCP02"}}); err != nil {
		t.Fatal(err)
	}

	if _, err := chef.CreateBacklogItem("An idea", "a hunch", nil, core.KindIdea, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := chef.CreateBacklogItem("The first task", "known work", []string{"the work"}, core.KindFeature, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := chef.CreateBacklogItem("The second task", "known work", []string{"the work"}, core.KindBug, nil); err != nil {
		t.Fatal(err)
	}
	gone, err := chef.CreateBacklogItem("A dismissed task", "known work", []string{"the work"}, core.KindFeature, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := chef.DismissBacklogItem(gone.ID, "not worth it"); err != nil {
		t.Fatal(err)
	}

	m := loadModel(t, store)
	m.tab = TabReady
	m.height = 40
	return m, store
}

// TestTheBoardStillShowsWhatReadyDoesNot covers the shrunk overlap.
//
// A draft, and an approved recipe with an open dependency, stay on the Board.
// Ready does not list them.
func TestTheBoardStillShowsWhatReadyDoesNot(t *testing.T) {
	m := newModel(t)
	m.tab = TabBoard
	m.height = 40
	board := plain(m)
	m.tab = TabReady
	ready := plain(m)
	for _, id := range []string{"RCP01", "RCP04"} {
		if !strings.Contains(board, id) {
			t.Errorf("the board should still show %s:\n%s", id, board)
		}
		if strings.Contains(ready, id) {
			t.Errorf("Ready should not show %s:\n%s", id, ready)
		}
	}
}

// TestSpaceOnAnAdoptableRecipeClaimsNothing covers the console's limit.
//
// The row opens the recipe. space writes nothing, and the footer does not offer
// a claim.
func TestSpaceOnAnAdoptableRecipeClaimsNothing(t *testing.T) {
	m, store := readyModel(t)
	out := plain(m)
	if !strings.Contains(out, "The free recipe") {
		t.Fatalf("the adoptable recipe should be on the tab:\n%s", out)
	}
	if strings.Contains(out, "Still being shaped") || strings.Contains(out, "Waiting on the draft") {
		t.Fatalf("a draft and an approved recipe with an open dependency stay off the tab:\n%s", out)
	}
	if strings.Contains(out, "claim") {
		t.Fatalf("the footer offers a claim:\n%s", out)
	}
	rec, ok := m.currentRecipe()
	if !ok || rec.Title != "The free recipe" {
		t.Fatalf("the cursor should open the adoptable recipe, got %+v", rec)
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	before := len(j.Events)
	if _, err := Drive(m, []string{"space"}); err != nil {
		t.Fatal(err)
	}
	j, err = store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) != before {
		t.Fatalf("space wrote %d events, want %d", len(j.Events), before)
	}
	got, err := store.LoadRecipe(rec.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Lifecycle != core.LifecycleApproved || got.Execution.ClaimedBy != "" {
		t.Fatalf("the recipe is %q claimed by %q — space claimed it", got.Lifecycle, got.Execution.ClaimedBy)
	}
}

// TestSpaceOnASmallTaskPromotesIt covers the move and the cursor.
//
// The row names the two further steps. space promotes, through the move the
// backlog already uses. The cursor does not stay on the id that left, and the
// tab does not change.
func TestSpaceOnASmallTaskPromotesIt(t *testing.T) {
	m, store := readyModel(t)
	out := plain(m)
	if !strings.Contains(out, "promote → sign off → run") {
		t.Fatalf("the row should name the two further steps:\n%s", out)
	}
	if strings.Contains(out, "An idea") || strings.Contains(out, "A dismissed task") {
		t.Fatalf("an idea and a dismissed task are on the tab:\n%s", out)
	}
	for _, id := range []string{"ADOPTABLE", "SMALL TASKS"} {
		if !strings.Contains(out, id) {
			t.Fatalf("the tab is missing %s:\n%s", id, out)
		}
	}

	// The adoptable recipe is the first row. The first small task is the next.
	if _, err := Drive(m, []string{"down"}); err != nil {
		t.Fatal(err)
	}
	item, ok := m.selectedReadyTask()
	if !ok {
		t.Fatal("the cursor should be on a small task")
	}
	if !strings.Contains(plain(m), "space promote to recipe") {
		t.Fatalf("the footer should ask the domain for the move:\n%s", plain(m))
	}
	left := item.ID
	if _, err := Drive(m, []string{"space"}); err != nil {
		t.Fatal(err)
	}
	if m.tab != TabReady {
		t.Fatalf("the tab is %s — a promotion left the tab", m.tab.Label())
	}
	if e, ok := m.currentReadyEntry(); ok && e.id == left {
		t.Fatalf("the cursor still names %s", left)
	}
	got, err := store.LoadBacklogItem(left)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != core.StatusPromoted {
		t.Fatalf("status = %q, want promoted", got.Status)
	}
	for _, e := range m.readyEntries() {
		if e.id == left {
			t.Fatalf("the promoted item %s is still a row", left)
		}
	}
}
