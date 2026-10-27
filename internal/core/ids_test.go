package core

import (
	"testing"
)

// TestIdsAreUniqueAfterARemoval covers the id allocation.
//
// The ids were allocated by counting the list, so removing an item made the next id
// collide with one that already existed — v1's `[AC-01, AC-03, AC-03]`. The plan
// directory is hand-editable, so the allocation scans what exists.
func TestIdsAreUniqueAfterARemoval(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	// Three of each family, so every family has at least three ids.
	for i := 0; i < 3; i++ {
		if _, err := ops.AddCriterion(r.ID, "another", "go test ./..."); err != nil {
			t.Fatal(err)
		}
		if _, err := ops.AddStep(r.ID, StepSpec{Action: "do it", DoneWhen: "done"}); err != nil {
			t.Fatal(err)
		}
		if _, err := ops.AddFinding(r.ID, Finding{
			Phase:       FindingPhaseExternal,
			Severity:    SeverityHigh,
			Category:    "correctness",
			Description: "a problem",
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := ops.AddComment(r.ID, "reviewer", "a note", ""); err != nil {
			t.Fatal(err)
		}
	}

	// The plan directory is hand-editable: remove the middle id of every family.
	raw, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	raw.Acceptance = append(raw.Acceptance[:1], raw.Acceptance[2:]...)
	raw.Guidance = append(raw.Guidance[:1], raw.Guidance[2:]...)
	raw.Review.Findings = append(raw.Review.Findings[:1], raw.Review.Findings[2:]...)
	raw.Review.Comments = append(raw.Review.Comments[:1], raw.Review.Comments[2:]...)
	if err := ops.store.SaveRecipe(raw); err != nil {
		t.Fatal(err)
	}

	// One more of each family. No id may collide with one that exists, and the
	// removed id may not come back.
	if _, err := ops.AddCriterion(r.ID, "another", "go test ./..."); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.AddStep(r.ID, StepSpec{Action: "do it", DoneWhen: "done"}); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.AddFinding(r.ID, Finding{
		Phase:       FindingPhaseExternal,
		Severity:    SeverityHigh,
		Category:    "correctness",
		Description: "a problem",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.AddComment(r.ID, "reviewer", "a note", ""); err != nil {
		t.Fatal(err)
	}

	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, family := range []struct {
		name   string
		ids    []string
		absent string
	}{
		{"criterion", criterionIDs(got), "AC-02"},
		{"step", stepIDs(got), "G2"},
		{"finding", findingIDs(got), "F-02"},
		{"comment", commentIDs(got), "C-02"},
	} {
		seen := map[string]bool{}
		for _, id := range family.ids {
			if seen[id] {
				t.Fatalf("the %s id %q was allocated twice: %v", family.name, id, family.ids)
			}
			seen[id] = true
		}
		if seen[family.absent] {
			t.Fatalf("the removed %s id %q came back: %v", family.name, family.absent, family.ids)
		}
	}
}
