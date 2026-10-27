package core

import (
	"testing"
)

// TestACancelledRecipeHasOnePlace covers the cancelled recipe.
//
// The lane filed it under "On the line" and the triage filed it under "Served" — two
// answers to one question, and a cancelled recipe presented as live work.
func TestACancelledRecipeHasOnePlace(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
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
	if _, err := ops.CancelRecipe(r.ID, "the premise was wrong"); err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}

	lane := LaneOf(*got)
	section := TriageOf(*got)
	if lane.Label() != section {
		t.Fatalf("the lane says %q and the triage says %q — one recipe, one place", lane.Label(), section)
	}
	if lane != LaneDone {
		t.Fatalf("lane = %q, want %q", lane, LaneDone)
	}
	if section != TriageServed {
		t.Fatalf("triage = %q, want %q", section, TriageServed)
	}
}

// TestTheChefIsNotOfferedADoomedServe covers the console's move.
//
// The close gate refuses a recipe with an open tasting note (`review.clean`), so a
// console that offers to serve it is offering an action the domain will refuse. The
// console offers the triage instead: dismissing the note, or saying something in its
// thread, is a move the domain accepts.
func TestTheChefIsNotOfferedADoomedServe(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	reviewedRecipe(t, ops, r.ID)

	// An agent files a tasting note after the review.
	if _, err := ops.AddFinding(r.ID, Finding{
		Phase:       FindingPhaseExternal,
		Severity:    SeverityHigh,
		Category:    "correctness",
		Description: "the queue drops findings on refresh",
	}); err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Lifecycle != LifecycleReviewed {
		t.Fatalf("lifecycle = %q, want reviewed", got.Lifecycle)
	}
	if move := ChefMoveOf(*got); move != MoveTriage {
		t.Fatalf("the chef's move = %q, want %q — the close gate would refuse a serve", move, MoveTriage)
	}
}

// TestTheChefCanServeACleanRecipe is the other side: the move is a check, not a wall.
func TestTheChefCanServeACleanRecipe(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	reviewedRecipe(t, ops, r.ID)

	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if move := ChefMoveOf(*got); move != MoveServe {
		t.Fatalf("the chef's move = %q, want %q", move, MoveServe)
	}
}

// TestTheTriageComesBeforeTheReview covers the ordering.
//
// A recipe can be both self-reviewed and carry an open tasting note. The note wins:
// it is the concrete blocker — the close gate refuses a recipe with an open note —
// and a verdict can wait. Triage the note, and the review is offered next.
func TestTheTriageComesBeforeTheReview(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
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
	if _, err := ops.ClaimRecipe(r.ID, "runner"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.SetStepStatus(r.ID, "G1", "done"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.SignOffCriterion(r.ID, "AC-01", "chef", "checked"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventFinish); err != nil {
		t.Fatal(err)
	}
	// The executor tastes its own dish.
	if _, err := ops.RecordReviewRound(r.ID, "runner", "ok", "looks good to me"); err != nil {
		t.Fatal(err)
	}
	// And the reviewer files a note on top.
	if _, err := ops.AddFinding(r.ID, Finding{
		Phase:       FindingPhaseExternal,
		Severity:    SeverityHigh,
		Category:    "correctness",
		Description: "the queue drops findings on refresh",
	}); err != nil {
		t.Fatal(err)
	}

	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !SelfReviewed(*got) {
		t.Fatal("the recipe should be self-reviewed")
	}
	if move := ChefMoveOf(*got); move != MoveTriage {
		t.Fatalf("the chef's move = %q, want %q — the open note is the blocker", move, MoveTriage)
	}
}

// reviewedRecipe walks a recipe to the reviewed state.
func reviewedRecipe(t *testing.T, ops *Operations, id string) {
	t.Helper()
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := ops.Advance(id, ev); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := ops.ApproveRecipe(id, "chef"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.ClaimRecipe(id, "runner"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.SetStepStatus(id, "G1", "done"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.SignOffCriterion(id, "AC-01", "chef", "checked"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(id, EventFinish); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.RecordReviewRound(id, "reviewer", "ok", "clean"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(id, EventPassReview); err != nil {
		t.Fatal(err)
	}
}
