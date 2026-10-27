package core

import (
	"strings"
	"testing"
)

// TestPassReviewNeedsAnIndependentReviewer covers the gate that enforces CPD
// principle 13 — the author should not be the only reviewer.
//
// The executor cannot sign off its own work, even with a verdict of "ok".
func TestPassReviewNeedsAnIndependentReviewer(t *testing.T) {
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

	// The executor reviews its own work and says it is fine.
	if _, err := ops.RecordReviewRound(r.ID, "runner", "ok", "looks good to me"); err != nil {
		t.Fatal(err)
	}
	_, err = ops.Advance(r.ID, EventPassReview)
	if err == nil {
		t.Fatal("the executor must not be able to sign off its own work")
	}
	if !strings.Contains(err.Error(), "who did the work") {
		t.Fatalf("the error should say why: %v", err)
	}

	// Somebody else reviews it.
	if _, err := ops.RecordReviewRound(r.ID, "reviewer", "ok", "clean"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventPassReview); err != nil {
		t.Fatalf("an independent review should advance the recipe: %v", err)
	}
}
