package core

import (
	"strings"
	"testing"
)

// TestPassReviewNeedsAnOKVerdict covers the hole: a reviewer who says "changes needed"
// must not be able to advance the recipe.
func TestPassReviewNeedsAnOKVerdict(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventGroom); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventSpecify); err != nil {
		t.Fatal(err)
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

	// The reviewer says no — and records no finding.
	if _, err := ops.RecordReviewRound(r.ID, "reviewer", "changes-needed", "not yet"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventPassReview); err == nil {
		t.Fatal("a changes-needed verdict must not advance the recipe")
	} else if !strings.Contains(err.Error(), "changes-needed") {
		t.Fatalf("the error should quote the verdict: %v", err)
	}

	// The reviewer says yes.
	if _, err := ops.RecordReviewRound(r.ID, "reviewer", "ok", "clean now"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventPassReview); err != nil {
		t.Fatalf("an ok verdict should advance the recipe: %v", err)
	}
}

// TestReopenFinding covers the review loop: a fix that does not hold goes back into the
// queue, and the objection is recorded so the next attempt can read it.
func TestReopenFinding(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.AddFinding(r.ID, Finding{
		Phase:       "external",
		Severity:    "high",
		Category:    "correctness",
		Description: "the queue drops findings on refresh",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.ResolveFinding(r.ID, "F-01", "read the queue once"); err != nil {
		t.Fatal(err)
	}

	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Review.Findings[0].Status != "fixed" {
		t.Fatalf("status = %q, want fixed", got.Review.Findings[0].Status)
	}

	// A reason is required.
	if _, err := ops.ReopenFinding(r.ID, "F-01", ""); err == nil {
		t.Fatal("a reopen should need a reason")
	}

	if _, err := ops.ReopenFinding(r.ID, "F-01", "the empty case is still missed"); err != nil {
		t.Fatal(err)
	}
	got, err = ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	f := got.Review.Findings[0]
	if f.Status != "open" {
		t.Fatalf("status = %q, want open", f.Status)
	}
	if f.Resolution != "" || f.Resolved != "" {
		t.Fatalf("the outcome should be cleared: %+v", f)
	}
	// The objection is on the record.
	found := false
	for _, c := range got.Review.Comments {
		if c.Finding == "F-01" && strings.Contains(c.Body, "empty case") {
			found = true
		}
	}
	if !found {
		t.Fatalf("the objection should be recorded as a comment: %+v", got.Review.Comments)
	}

	// Reopening an open finding is refused.
	if _, err := ops.ReopenFinding(r.ID, "F-01", "again"); err == nil {
		t.Fatal("reopening an open finding should be refused")
	}
}

// TestTheJournalCanBeReadFromASequence covers how an agent notices what changed.
func TestTheJournalCanBeReadFromASequence(t *testing.T) {
	ops, store := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	journal, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	seen := journal.Seq
	if seen == 0 {
		t.Fatal("the journal is empty")
	}

	if _, err := ops.AddCriterion(r.ID, "another", "go test ./..."); err != nil {
		t.Fatal(err)
	}
	journal, err = store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	got := journal.Since(seen)
	if len(got) != 1 {
		t.Fatalf("got %d events since %d, want 1", len(got), seen)
	}
	if got[0].Seq <= seen {
		t.Fatalf("the event's sequence %d should be after %d", got[0].Seq, seen)
	}
	if len(journal.Since(journal.Seq)) != 0 {
		t.Fatal("nothing should be after the last sequence number")
	}
}
