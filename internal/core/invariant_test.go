package core

import (
	"testing"
)

// TestNeedsYouMeansThereIsAMoveForYou covers a whole class of bug, not one instance.
//
// The console filed a self-reviewed recipe under "Needs you" while the gate refused
// to close it, and counted "0 need you". A section that says *you* and offers nothing
// to press is a lie — and a section that says *you* for something with no move is the
// same lie wearing a different hat.
//
// The invariant: a recipe is filed under "Needs you" if and only if the console has a
// chef move for it.
func TestNeedsYouMeansThereIsAMoveForYou(t *testing.T) {
	for _, s := range everyShape() {
		r := s.recipe
		needs := TriageOf(r) == TriageNeedsYou
		move := ChefMoveOf(r)

		if needs && move == MoveNone {
			t.Errorf("%s (%s): filed under 'Needs you' but the console has no move", r.ID, r.Lifecycle)
		}
		if !needs && move != MoveNone {
			t.Errorf("%s (%s): the console offers %q but it is not filed under 'Needs you'", r.ID, r.Lifecycle, move)
		}
	}
}

// TestARecipeWithAnOpenNoteHasAMove covers the exception that came out.
//
// The console filed a recipe with an open tasting note under "Needs you" and offered
// nothing to press. The console can triage a note now, so the recipe has a move and
// the exception is gone.
func TestARecipeWithAnOpenNoteHasAMove(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	reviewedRecipe(t, ops, r.ID)
	// The taster files a note after the review.
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
	if TriageOf(*got) != TriageNeedsYou {
		t.Fatalf("the recipe should be filed under 'Needs you', got %q", TriageOf(*got))
	}
	if move := ChefMoveOf(*got); move != MoveTriage {
		t.Fatalf("the chef's move = %q, want %q", move, MoveTriage)
	}
}

// shape is one recipe of a given shape, and what the console can do about it.
type shape struct {
	recipe Recipe
}

// everyShape is one recipe of every shape the console can meet: each lifecycle, each
// review arrangement, each overlay.
func everyShape() []shape {
	base := Recipe{ID: "R01", Title: "a shape", Lifecycle: LifecycleDraft}

	var out []shape
	for _, lc := range []Lifecycle{
		LifecycleDraft, LifecycleGroomed, LifecycleReady, LifecycleApproved,
		LifecycleExecuting, LifecycleExecuted, LifecycleRemediation, LifecycleReviewed,
		LifecycleClosed,
	} {
		r := base
		r.Lifecycle = lc
		out = append(out, shape{recipe: r})

		// Claimed by an agent, with that agent's own review round.
		self := r
		self.Execution = Execution{ClaimedBy: "line-cook"}
		self.Review.Rounds = []Round{{N: 1, Reviewer: "line-cook", Verdict: "ok"}}
		out = append(out, shape{recipe: self})

		// Reviewed by somebody else.
		other := r
		other.Execution = Execution{ClaimedBy: "line-cook"}
		other.Review.Rounds = []Round{{N: 1, Reviewer: "taster", Verdict: "ok"}}
		out = append(out, shape{recipe: other})

		// Changes asked for.
		changes := r
		changes.Execution = Execution{ClaimedBy: "line-cook"}
		changes.Review.Rounds = []Round{{N: 1, Reviewer: "chef", Verdict: "changes-needed"}}
		out = append(out, shape{recipe: changes})

		// An open finding. The console can triage it, so it has a move of its own.
		open := r
		open.Execution = Execution{ClaimedBy: "line-cook"}
		open.Review.Findings = []Finding{{ID: "F-01", Status: "open", Phase: "external"}}
		out = append(out, shape{recipe: open})

		// On hold.
		hold := r
		hold.Overlays.Blocked = &Overlay{Reason: "waiting on the client"}
		out = append(out, shape{recipe: hold})
	}
	return out
}
