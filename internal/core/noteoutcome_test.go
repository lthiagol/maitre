package core

import (
	"errors"
	"strings"
	"testing"
)

// noteRecipe is a recipe with one open tasting note.
func noteRecipe(t *testing.T, ops *Operations) string {
	t.Helper()
	r, err := ops.CreateRecipe(spec())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.AddFinding(r.ID, Finding{
		Phase:       FindingPhaseExternal,
		Severity:    SeverityHigh,
		Category:    "correctness",
		Description: "the queue drops findings on refresh",
	}); err != nil {
		t.Fatal(err)
	}
	return r.ID
}

// TestADismissalDoesNotOverwriteTheFixersRecord covers a note's outcome.
//
// Whoever writes last used to win: the chef's dismissal destroyed the fixer's record
// of what was done, and the artifact said the note was dismissed when it had been
// fixed.
func TestADismissalDoesNotOverwriteTheFixersRecord(t *testing.T) {
	ops, _ := newOps(t)
	id := noteRecipe(t, ops)

	if _, err := ops.ResolveFinding(id, "F-01", "fixed it properly"); err != nil {
		t.Fatal(err)
	}
	_, err := ops.DismissFinding(id, "F-01", "not needed")
	if err == nil {
		t.Fatal("a dismissal on a fixed note should be refused")
	}
	var refusal *RefusalError
	if !errors.As(err, &refusal) {
		t.Fatalf("the refusal should be the plan's answer, got %T", err)
	}

	got, err := ops.GetRecipe(id)
	if err != nil {
		t.Fatal(err)
	}
	f := got.Review.Findings[0]
	if f.Status != "fixed" {
		t.Fatalf("the note's status = %q, want fixed — the fixer's record must survive", f.Status)
	}
	if f.Resolution != "fixed it properly" {
		t.Fatalf("the note's resolution = %q, want what the fixer wrote", f.Resolution)
	}
}

// TestAResolutionDoesNotOverwriteTheChefsRecord is the other way round.
func TestAResolutionDoesNotOverwriteTheChefsRecord(t *testing.T) {
	ops, _ := newOps(t)
	id := noteRecipe(t, ops)

	if _, err := ops.DismissFinding(id, "F-01", "not worth fixing"); err != nil {
		t.Fatal(err)
	}
	_, err := ops.ResolveFinding(id, "F-01", "fixed it properly")
	if err == nil {
		t.Fatal("a resolution on a dismissed note should be refused")
	}
	var refusal *RefusalError
	if !errors.As(err, &refusal) {
		t.Fatalf("the refusal should be the plan's answer, got %T", err)
	}

	got, err := ops.GetRecipe(id)
	if err != nil {
		t.Fatal(err)
	}
	f := got.Review.Findings[0]
	if f.Status != "dismissed" {
		t.Fatalf("the note's status = %q, want dismissed", f.Status)
	}
	if f.Resolution != "not worth fixing" {
		t.Fatalf("the note's resolution = %q, want the chef's reason", f.Resolution)
	}
}

// TestAReopenCarriesTheOutcomeItErased covers the reviewer's move.
//
// The reopen clears the note's outcome — the note is open again, so it has none — but
// the record of what was done goes into the thread, which is append-only. A reopen
// that dropped it would erase the fixer's record by the back door, which is the
// outcome the write-once rule says cannot happen.
func TestAReopenCarriesTheOutcomeItErased(t *testing.T) {
	ops, _ := newOps(t)
	id := noteRecipe(t, ops)
	if _, err := ops.ResolveFinding(id, "F-01", "fixed it properly in the loader"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.ReopenFinding(id, "F-01", "it still drops findings on refresh"); err != nil {
		t.Fatal(err)
	}

	got, err := ops.GetRecipe(id)
	if err != nil {
		t.Fatal(err)
	}
	f := got.Review.Findings[0]
	if f.Status != "open" {
		t.Fatalf("the note's status = %q, want open", f.Status)
	}
	if f.Resolution != "" {
		t.Fatalf("the note is open again, so it has no outcome: %q", f.Resolution)
	}

	// The thread carries what was erased.
	var found bool
	for _, c := range got.Review.Comments {
		if c.Finding != "F-01" {
			continue
		}
		found = true
		if !strings.Contains(c.Body, "it still drops findings on refresh") {
			t.Errorf("the comment should carry the objection: %q", c.Body)
		}
		if !strings.Contains(c.Body, "fixed it properly in the loader") {
			t.Errorf("the comment should carry the outcome it erased: %q", c.Body)
		}
	}
	if !found {
		t.Fatal("the reopen wrote no comment")
	}
}

// TestARefusedNoteOutcomeIsJournalled covers the refusal as the plan's answer.
//
// The journal is the only event log: an attempt that was turned away is evidence, and
// the plan and the journal have to agree.
func TestARefusedNoteOutcomeIsJournalled(t *testing.T) {
	ops, store := newOps(t)
	id := noteRecipe(t, ops)
	if _, err := ops.ResolveFinding(id, "F-01", "fixed it properly"); err != nil {
		t.Fatal(err)
	}

	before, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.DismissFinding(id, "F-01", "not needed"); err == nil {
		t.Fatal("a dismissal on a fixed note should be refused")
	}
	after, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Events) != len(before.Events)+1 {
		t.Fatalf("the journal grew by %d events, want 1", len(after.Events)-len(before.Events))
	}
	last := after.Events[len(after.Events)-1]
	if last.Kind != KindRefusal {
		t.Fatalf("the event's kind = %q, want %q", last.Kind, KindRefusal)
	}
	if last.Recipe != id {
		t.Fatalf("the event names %q, want %q", last.Recipe, id)
	}
}
