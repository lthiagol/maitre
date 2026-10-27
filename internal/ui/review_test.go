package ui

import (
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// typed expands a word into the key presses that type it.
//
// The driver takes one key per name, and a name it does not know is an error — a
// typo in a script should not be silently typed into the plan.
func typed(word string) []string {
	out := make([]string, 0, len(word))
	for _, r := range word {
		out = append(out, string(r))
	}
	return out
}

// reviewKeys drives a whole review: space, a verdict, a note, enter.
func reviewKeys(verdict, notes string) []string {
	keys := []string{"space", verdict}
	keys = append(keys, typed(notes)...)
	return append(keys, "enter")
}

// TestTheConsoleOffersTheReview checks the footer names the chef's move.
func TestTheConsoleOffersTheReview(t *testing.T) {
	m := selfReviewedModel(t)
	out := StripANSI(m.View().Content)
	if !strings.Contains(out, "space review") {
		t.Errorf("the footer should offer the review:\n%s", out)
	}
}

// TestTheReviewPromptAsksForAVerdict checks the console asks rather than decides.
//
// A verdict the console invented would be a verdict nobody gave.
func TestTheReviewPromptAsksForAVerdict(t *testing.T) {
	m := selfReviewedModel(t)
	out, err := Drive(m, []string{"space"})
	if err != nil {
		t.Fatal(err)
	}
	plain := StripANSI(out)
	if !strings.Contains(plain, "o) ok") || !strings.Contains(plain, "c) changes needed") {
		t.Errorf("the prompt should offer both verdicts:\n%s", plain)
	}
	if !strings.Contains(plain, "esc cancel") {
		t.Errorf("the prompt should say how to back out:\n%s", plain)
	}
}

// TestTheReviewRoundIsRecorded checks the round lands in the artifact.
func TestTheReviewRoundIsRecorded(t *testing.T) {
	m := selfReviewedModel(t)
	if _, err := Drive(m, reviewKeys("o", "cleanwork")); err != nil {
		t.Fatal(err)
	}

	rec, err := m.ops.GetRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	rounds := rec.Review.Rounds
	if len(rounds) != 2 {
		t.Fatalf("expected two rounds, got %d", len(rounds))
	}
	last := rounds[len(rounds)-1]
	if last.Reviewer != "chef" {
		t.Errorf("the reviewer should be the chef, got %q", last.Reviewer)
	}
	if last.Verdict != "ok" {
		t.Errorf("the verdict should be ok, got %q", last.Verdict)
	}
	if last.Notes != "cleanwork" {
		t.Errorf("the notes should be what was typed, got %q", last.Notes)
	}
}

// TestThePromptCanBeCancelled checks esc leaves no trace.
func TestThePromptCanBeCancelled(t *testing.T) {
	m := selfReviewedModel(t)
	if _, err := Drive(m, append([]string{"space", "o"}, append(typed("half"), "esc")...)); err != nil {
		t.Fatal(err)
	}
	if m.prompt != nil {
		t.Fatal("the prompt should be closed")
	}
	rec, err := m.ops.GetRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Review.Rounds) != 1 {
		t.Fatalf("a cancelled prompt should record nothing, got %d rounds", len(rec.Review.Rounds))
	}
}

// TestThePromptHoldsTheSelection checks a half-typed review is not thrown away by a
// stray keypress.
func TestThePromptHoldsTheSelection(t *testing.T) {
	m := selfReviewedModel(t)
	if _, err := Drive(m, []string{"space", "o"}); err != nil {
		t.Fatal(err)
	}
	before := m.selRecipes
	if _, err := Drive(m, []string{"down", "down"}); err != nil {
		t.Fatal(err)
	}
	if m.selRecipes != before {
		t.Errorf("the selection moved while the prompt was open")
	}
	if m.prompt == nil {
		t.Fatal("the prompt closed on its own")
	}
}

// TestTheConsoleOffersToServe checks the last move is the chef's too.
func TestTheConsoleOffersToServe(t *testing.T) {
	m := selfReviewedModel(t)
	if _, err := Drive(m, reviewKeys("o", "clean")); err != nil {
		t.Fatal(err)
	}
	out := StripANSI(m.View().Content)
	if !strings.Contains(out, "space serve") {
		t.Errorf("the footer should offer to serve a reviewed recipe:\n%s", out)
	}
}

// TestServingClosesTheRecipe checks the serve action goes through the gates.
func TestServingClosesTheRecipe(t *testing.T) {
	m := selfReviewedModel(t)
	if _, err := Drive(m, reviewKeys("o", "clean")); err != nil {
		t.Fatal(err)
	}
	if _, err := Drive(m, []string{"space"}); err != nil {
		t.Fatal(err)
	}

	rec, err := m.ops.GetRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Lifecycle != core.LifecycleClosed {
		t.Fatalf("the recipe should be closed, got %q (refusal: %s)", rec.Lifecycle, m.refusal)
	}
}

// TestTheConsoleDoesNotOfferADoomedServe checks the console does not close a recipe
// the domain would refuse.
//
// With an open tasting note the close gate refuses, so the console does not offer to
// serve at all — it offers the triage, which the domain accepts. Either way the
// recipe stays open.
func TestTheConsoleDoesNotOfferADoomedServe(t *testing.T) {
	m := selfReviewedModel(t)
	// An open finding: the recipe is not clean, so it cannot be served.
	if _, err := m.ops.AddFinding("RCP01", core.Finding{
		Phase:       "external",
		Severity:    "high",
		Category:    "correctness",
		Description: "the thread is cut off at the bottom",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := Drive(m, reviewKeys("o", "clean")); err != nil {
		t.Fatal(err)
	}
	if _, err := Drive(m, []string{"space"}); err != nil {
		t.Fatal(err)
	}

	rec, err := m.ops.GetRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Lifecycle == core.LifecycleClosed {
		t.Fatal("the recipe should not have closed")
	}
	if m.prompt == nil {
		t.Error("the console should offer the triage, not a doomed serve")
	}
}
