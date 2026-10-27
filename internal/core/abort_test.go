package core

import (
	"strings"
	"testing"
)

// TestAbortClearsTheApproval covers the trap: an aborted recipe must be editable again.
//
// The approval freezes the intent, problem and scope. If an abort left the approval in
// place, the recipe would be back in draft but permanently un-editable.
func TestAbortClearsTheApproval(t *testing.T) {
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
	if _, err := ops.ClaimRecipe(r.ID, "line-cook"); err != nil {
		t.Fatal(err)
	}

	// While approved, the frozen sections cannot be edited.
	outcome := "something else"
	if _, err := ops.UpdateRecipe(r.ID, Fragment{Outcome: &outcome}); err == nil {
		t.Fatal("the intent should be frozen while the recipe is signed off")
	}

	// Abort.
	if _, err := ops.AdvanceWith(r.ID, EventAbortToDraft, TransitionContext{Reason: "the premise was wrong"}); err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Lifecycle != LifecycleDraft {
		t.Fatalf("lifecycle = %q, want draft", got.Lifecycle)
	}
	if got.Approval != nil {
		t.Fatalf("the approval should be cleared, got %+v", got.Approval)
	}
	if got.Execution.ClaimedBy != "" {
		t.Fatalf("the claim should be cleared, got %+v", got.Execution)
	}

	// And the recipe can be edited again.
	if _, err := ops.UpdateRecipe(r.ID, Fragment{Outcome: &outcome}); err != nil {
		t.Fatalf("the recipe should be editable again after an abort: %v", err)
	}

	// The abort is on the record.
	if len(got.Audit.Deviations) == 0 || !strings.Contains(got.Audit.Deviations[0].Note, "the premise was wrong") {
		t.Fatalf("the abort should be recorded: %+v", got.Audit.Deviations)
	}
}
