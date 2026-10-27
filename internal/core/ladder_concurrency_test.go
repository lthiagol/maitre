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
// the directory at the same time and be handed the same id — two specials with the
// same number, one of them unreachable by id.
func TestConcurrentCreationsTakeDistinctIDs(t *testing.T) {
	ops, _ := newOps(t)

	const n = 16
	var wg sync.WaitGroup
	start := make(chan struct{})
	specialIDs := make([]string, n)
	noteIDs := make([]string, n)
	errs := make([]error, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			s, err := ops.CreateSpecial(fmt.Sprintf("A hunch %d", i), "worth noticing", nil)
			if err != nil {
				errs[i] = err
				return
			}
			specialIDs[i] = s.ID
		}(i)
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			note, err := ops.CreatePrepNote(fmt.Sprintf("A note %d", i), "first", nil)
			if err != nil {
				errs[i] = err
				return
			}
			noteIDs[i] = note.ID
		}(i)
	}
	close(start)
	wg.Wait()

	seen := map[string]bool{}
	for i := 0; i < n; i++ {
		if errs[i] != nil {
			t.Fatalf("creation %d: %v", i, errs[i])
		}
		for _, id := range []string{specialIDs[i], noteIDs[i]} {
			if id == "" {
				t.Fatalf("creation %d got no id", i)
			}
			if seen[id] {
				t.Fatalf("the id %q was handed out twice", id)
			}
			seen[id] = true
		}
	}
}

// TestAFailedPromotionLeavesTheSourceOpen covers the failed promotion.
//
// A promotion writes the source and the new artifact. If the target write fails, the
// source must not claim it was promoted: the artifact is the record, and a record of
// a promotion that did not happen is a lie.
func TestAFailedPromotionLeavesTheSourceOpen(t *testing.T) {
	ops, store := newOps(t)

	s, err := ops.CreateSpecial("A hunch", "worth noticing", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Make the target write fail: a directory where the prep note's file goes.
	noteID, err := store.NextPrepNoteID()
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(store.Root, DirPrepNotes, noteID+"-"+Slugify(s.Title)+".json")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	before, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ops.PromoteSpecial(s.ID); err == nil {
		t.Fatal("the promotion should fail")
	}

	// The source is as it was.
	after, err := store.LoadSpecial(s.ID)
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

	// The state a crash leaves: a prep note that points at a recipe nobody wrote.
	note, err := ops.CreatePrepNote("A note", "first", nil)
	if err != nil {
		t.Fatal(err)
	}
	note.Status = StatusPromoted
	note.PromotedTo = "recipe R01"
	if err := store.SavePrepNote(note); err != nil {
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

// TestAPromotionStillWorksInBothDirections is the other side: the fix does not break
// the promotion it protects.
func TestAPromotionStillWorksInBothDirections(t *testing.T) {
	ops, store := newOps(t)

	s, err := ops.CreateSpecial("A hunch", "worth noticing", nil)
	if err != nil {
		t.Fatal(err)
	}
	note, err := ops.PromoteSpecial(s.ID)
	if err != nil {
		t.Fatalf("promote special: %v", err)
	}
	if note.ID != "N01" {
		t.Fatalf("prep note id = %q, want N01", note.ID)
	}
	after, err := store.LoadSpecial(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != StatusPromoted || after.PromotedTo != "prep-note N01" {
		t.Fatalf("the special should record where it went: %+v", after)
	}
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if last := j.Events[len(j.Events)-1]; last.Kind != KindPromotion {
		t.Fatalf("the journal kind = %q, want %q", last.Kind, KindPromotion)
	}

	// And the other direction.
	recipe, err := ops.PromotePrepNote(note.ID)
	if err != nil {
		t.Fatalf("promote prep note: %v", err)
	}
	if recipe.ID != "R01" {
		t.Fatalf("recipe id = %q, want R01", recipe.ID)
	}
	afterNote, err := store.LoadPrepNote(note.ID)
	if err != nil {
		t.Fatal(err)
	}
	if afterNote.Status != StatusPromoted || afterNote.PromotedTo != "recipe R01" {
		t.Fatalf("the prep note should record where it went: %+v", afterNote)
	}
	j, err = store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if last := j.Events[len(j.Events)-1]; last.Kind != KindPromotion {
		t.Fatalf("the journal kind = %q, want %q", last.Kind, KindPromotion)
	}
}

// TestTwoConcurrentPromotionsCannotBothSucceed covers the one-way ladder.
//
// The promotion takes the plan lock, so the second one finds the source already
// promoted — one source, one target.
func TestTwoConcurrentPromotionsCannotBothSucceed(t *testing.T) {
	ops, store := newOps(t)

	s, err := ops.CreateSpecial("A hunch", "worth noticing", nil)
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
			_, errs[i] = ops.PromoteSpecial(s.ID)
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
	notes, err := store.ListPrepNotes()
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 1 {
		t.Fatalf("got %d prep notes, want 1 — one special, one prep note", len(notes))
	}
}
