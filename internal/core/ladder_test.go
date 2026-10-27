package core

import (
	"testing"
)

func TestSpecialLifecycle(t *testing.T) {
	ops, _ := newOps(t)

	s, err := ops.CreateSpecial("A hunch", "something worth noticing", []string{"model"})
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != "S01" || s.Status != StatusOpen {
		t.Fatalf("special = %+v", s)
	}

	// Editing works while it is open.
	if _, err := ops.UpdateSpecial(s.ID, "A better hunch", "sharper", nil); err != nil {
		t.Fatal(err)
	}
	got, err := ops.store.LoadSpecial("S01")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "A better hunch" || got.Body != "sharper" {
		t.Fatalf("special = %+v", got)
	}

	// Promotion is one-way, and the trail survives.
	note, err := ops.PromoteSpecial(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if note.ID != "N01" || note.Title != "A better hunch" || note.Intent != "sharper" {
		t.Fatalf("prep note = %+v", note)
	}
	after, err := ops.store.LoadSpecial("S01")
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != StatusPromoted || after.PromotedTo != "prep-note N01" {
		t.Fatalf("special was not marked promoted: %+v", after)
	}
	// A promoted special cannot be promoted again.
	if _, err := ops.PromoteSpecial(s.ID); err == nil {
		t.Fatal("expected promoting a promoted special to fail")
	}
}

// TestASpecialsEventsAreJournalledMeaningfully covers the journal's kinds.
//
// The kinds are a closed enum and they are meant to mean something: an edit was
// journalled as a `comment` and a dismissal as a `promotion`, so a reader of the log
// saw a promotion for something the chef threw away.
func TestASpecialsEventsAreJournalledMeaningfully(t *testing.T) {
	ops, store := newOps(t)

	s, err := ops.CreateSpecial("A hunch", "something worth noticing", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.UpdateSpecial(s.ID, "A better hunch", "sharper", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.DismissSpecial(s.ID, "not worth pursuing"); err != nil {
		t.Fatal(err)
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) != 3 {
		t.Fatalf("the journal has %d events, want 3", len(j.Events))
	}
	if j.Events[0].Kind != KindCreated {
		t.Fatalf("kind = %q, want %q", j.Events[0].Kind, KindCreated)
	}
	if j.Events[1].Kind != KindDecision {
		t.Fatalf("the edit's kind = %q, want %q", j.Events[1].Kind, KindDecision)
	}
	if j.Events[2].Kind != KindOverlay {
		t.Fatalf("the dismissal's kind = %q, want %q", j.Events[2].Kind, KindOverlay)
	}
}

func TestPrepNotePromotesToRecipe(t *testing.T) {
	ops, _ := newOps(t)

	note, err := ops.CreatePrepNote("Wire the queue", "Findings reach the chef", []string{"the queue lane"})
	if err != nil {
		t.Fatal(err)
	}
	if note.ID != "N01" {
		t.Fatalf("prep note id = %q, want N01", note.ID)
	}

	recipe, err := ops.PromotePrepNote(note.ID)
	if err != nil {
		t.Fatal(err)
	}
	if recipe.Lifecycle != LifecycleDraft {
		t.Fatalf("lifecycle = %q, want draft", recipe.Lifecycle)
	}
	if recipe.Title != "Wire the queue" || recipe.Intent.Outcome != "Findings reach the chef" {
		t.Fatalf("recipe = %+v", recipe)
	}
	if len(recipe.Scope.InScope) != 1 || recipe.Scope.InScope[0] != "the queue lane" {
		t.Fatalf("scope did not carry over: %+v", recipe.Scope)
	}

	after, err := ops.store.LoadPrepNote(note.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != StatusPromoted || after.PromotedTo != "recipe "+recipe.ID {
		t.Fatalf("prep note was not marked promoted: %+v", after)
	}
}

func TestDismissSpecial(t *testing.T) {
	ops, _ := newOps(t)
	s, err := ops.CreateSpecial("A hunch", "body", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.DismissSpecial(s.ID, "not worth it"); err != nil {
		t.Fatal(err)
	}
	got, err := ops.store.LoadSpecial(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusDismissed {
		t.Fatalf("status = %q, want dismissed", got.Status)
	}
	if _, err := ops.DismissSpecial(s.ID, "again"); err == nil {
		t.Fatal("expected dismissing a dismissed special to fail")
	}
}

func TestSpecialsAndPrepNotesAreListed(t *testing.T) {
	ops, store := newOps(t)
	for i := 0; i < 3; i++ {
		if _, err := ops.CreateSpecial("hunch", "body", nil); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := ops.CreatePrepNote("note", "intent", nil); err != nil {
		t.Fatal(err)
	}

	specials, err := ops.ListSpecials()
	if err != nil {
		t.Fatal(err)
	}
	if len(specials) != 3 {
		t.Fatalf("got %d specials, want 3", len(specials))
	}
	notes, err := ops.ListPrepNotes()
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 1 {
		t.Fatalf("got %d prep notes, want 1", len(notes))
	}

	// The journal records the promotions and creations.
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	kinds := map[string]int{}
	for _, e := range j.Events {
		kinds[e.Kind]++
	}
	if kinds[KindCreated] < 4 {
		t.Fatalf("expected the creations to be journalled: %v", kinds)
	}
}
