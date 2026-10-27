package core

import (
	"strings"
	"testing"
)

// selfReviewedRecipe builds a recipe that the executor reviewed itself.
func selfReviewedRecipe(t *testing.T) (*Operations, *Recipe) {
	t.Helper()
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
	if _, err := ops.RecordReviewRound(r.ID, "runner", "ok", "looks good to me"); err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	return ops, got
}

// TestSelfReviewedIsOneAnswer covers the fact the console and the gate both read.
//
// The console used to call a self-reviewed recipe "tasting clean" and filed it under
// "On the line", while the gate refused to let it close. The console was telling the
// chef a lie.
func TestSelfReviewedIsOneAnswer(t *testing.T) {
	_, r := selfReviewedRecipe(t)

	if !SelfReviewed(*r) {
		t.Fatal("the recipe should read as self-reviewed")
	}
	if got := ReviewStateOf(*r); got != ReviewSelfSigned {
		t.Errorf("review state = %q, want %q", got, ReviewSelfSigned)
	}
	if got := TriageOf(*r); got != TriageNeedsYou {
		t.Errorf("triage = %q, want %q", got, TriageNeedsYou)
	}
	reason := ReadinessOf(*r).Reason
	if !strings.Contains(reason, "independent review") {
		t.Errorf("readiness reason = %q, want it to name the independent review", reason)
	}
}

// TestAnIndependentReviewNeedsNobody checks the other side: once somebody else has
// tasted it, the recipe is off the chef's plate.
func TestAnIndependentReviewNeedsNobody(t *testing.T) {
	ops, r := selfReviewedRecipe(t)

	if _, err := ops.RecordReviewRound(r.ID, "reviewer", "ok", "clean"); err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}

	if SelfReviewed(*got) {
		t.Fatal("the recipe should no longer read as self-reviewed")
	}
	if state := ReviewStateOf(*got); state != ReviewClean {
		t.Errorf("review state = %q, want %q", state, ReviewClean)
	}
	if triage := TriageOf(*got); triage != TriageOnTheLine {
		t.Errorf("triage = %q, want %q", triage, TriageOnTheLine)
	}
}

// TestTheGateAndTheTriageAgree is the drift check.
//
// One fact, one answer: if the gate refuses and the console says the recipe is
// settled, the chef is never told what is waiting.
func TestTheGateAndTheTriageAgree(t *testing.T) {
	ops, r := selfReviewedRecipe(t)

	_, err := ops.Advance(r.ID, EventPassReview)
	if err == nil {
		t.Fatal("the gate should refuse")
	}
	if TriageOf(*r) != TriageNeedsYou {
		t.Fatal("the console should say the recipe needs the chef")
	}
}

// TestThePathSaysWhoIsWaiting checks the path tab names the same blocker.
func TestThePathSaysWhoIsWaiting(t *testing.T) {
	_, r := selfReviewedRecipe(t)
	ready, waiting := Path([]Recipe{*r})
	if len(ready) != 1 || ready[0].Note != "needs you" {
		t.Fatalf("the path should say who is waiting, got %+v", ready)
	}
	if len(waiting) != 0 {
		t.Fatalf("nothing is waiting on a dependency, got %+v", waiting)
	}
}
