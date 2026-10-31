package core

import (
	"strings"
	"testing"
)

// TestDomainVersionIsSeven covers the writer.
//
// This slice changes no rule the tool applies. The version stays 7.
func TestDomainVersionIsSeven(t *testing.T) {
	if DomainVersion != 7 {
		t.Fatalf("DomainVersion = %d, want 7", DomainVersion)
	}
}

// TestAdoptableAsksTheClaim covers the one answer.
//
// A recipe is adoptable iff EventStart would be accepted. The view asks that
// transition on a copy, so the cases below are the claim's cases: an approved
// recipe whose dependencies are closed is in; an approved recipe with an unmet
// dependency is out; a cancelled recipe, one already on the line, and one still
// being shaped are out. A gate the claim grows is a gate this set grows.
func TestAdoptableAsksTheClaim(t *testing.T) {
	plan := []Recipe{
		{ID: "RCP01", Lifecycle: LifecycleClosed},
		{ID: "RCP02", Title: "Free", Lifecycle: LifecycleApproved},
		{ID: "RCP03", Title: "After RCP01", Lifecycle: LifecycleApproved, DependsOn: []string{"RCP01"}},
		{ID: "RCP04", Title: "Waiting", Lifecycle: LifecycleApproved, DependsOn: []string{"RCP09"}},
		{ID: "RCP05", Lifecycle: LifecycleApproved, Overlays: Overlays{Cancelled: &Overlay{Reason: "no"}}},
		{ID: "RCP06", Lifecycle: LifecycleExecuting},
		{ID: "RCP07", Lifecycle: LifecycleDraft},
		{ID: "RCP08", Lifecycle: LifecycleReady},
		{ID: "RCP10", Title: "Tampered", Lifecycle: LifecycleApproved, Approval: &Approval{Fingerprint: "not-the-digest"}},
		// EventStart has no blocked gate, so a blocked recipe the claim accepts is in.
		{ID: "RCP11", Lifecycle: LifecycleApproved, Overlays: Overlays{Blocked: &Overlay{Reason: "waiting"}}},
	}
	before := make([]Lifecycle, len(plan))
	for i, r := range plan {
		before[i] = r.Lifecycle
	}

	closed := closedRecipes(plan)
	want := map[string]bool{}
	for _, r := range plan {
		copy := r
		err := Transition(&copy, EventStart, TransitionContext{
			Closed: func(id string) (bool, error) { return closed[id], nil },
		})
		want[r.ID] = err == nil
	}

	got := map[string]bool{}
	var order []string
	for _, r := range Adoptable(plan) {
		got[r.ID] = true
		order = append(order, r.ID)
	}
	for _, r := range plan {
		if got[r.ID] != want[r.ID] {
			t.Errorf("%s: adoptable %v, the claim %v", r.ID, got[r.ID], want[r.ID])
		}
	}
	for i, r := range plan {
		if r.Lifecycle != before[i] {
			t.Errorf("%s: lifecycle = %q, want %q — the view asks a copy", r.ID, r.Lifecycle, before[i])
		}
	}

	if !got["RCP02"] || !got["RCP03"] {
		t.Fatalf("an approved recipe with its dependencies closed must be in: %v", got)
	}
	if got["RCP04"] {
		t.Fatal("an approved recipe with an unmet dependency must be out")
	}
	if got["RCP10"] {
		t.Fatal("a recipe whose sign-off no longer matches must be out — the claim's immutability gate")
	}
	if !got["RCP11"] {
		t.Fatal("the claim accepts a blocked recipe, so the view does too")
	}

	// Plan order, not a second sort.
	wantOrder := []string{"RCP02", "RCP03", "RCP11"}
	if strings.Join(order, ",") != strings.Join(wantOrder, ",") {
		t.Fatalf("order = %v, want %v", order, wantOrder)
	}

	if Adoptable(nil) == nil {
		t.Fatal("an empty plan is an empty list")
	}
}

// TestTheBoardKeepsWhatAdoptableDrops covers the two questions.
//
// An approved recipe with an unmet dependency, and a draft, stay on Path and on
// the tree. Adoptable drops both.
func TestTheBoardKeepsWhatAdoptableDrops(t *testing.T) {
	plan := []Recipe{
		{ID: "RCP01", Lifecycle: LifecycleDraft},
		{ID: "RCP02", Lifecycle: LifecycleApproved, DependsOn: []string{"RCP01"}},
	}
	if got := Adoptable(plan); len(got) != 0 {
		t.Fatalf("adoptable has %d recipes, want none", len(got))
	}
	ready, waiting := Path(plan)
	if len(ready) != 1 || ready[0].Summary.ID != "RCP01" {
		t.Fatalf("path ready = %+v, want RCP01", ready)
	}
	if len(waiting) != 1 || waiting[0].Summary.ID != "RCP02" || !strings.Contains(waiting[0].Note, "RCP01") {
		t.Fatalf("path waiting = %+v, want RCP02 needing RCP01", waiting)
	}
	tree := TreeOf(plan)
	seen := map[string]bool{}
	var walk func(TreeNode)
	walk = func(n TreeNode) {
		seen[n.Summary.ID] = true
		for _, c := range n.Waiting {
			walk(c)
		}
	}
	for _, n := range tree.Roots {
		walk(n)
	}
	if !seen["RCP01"] || !seen["RCP02"] {
		t.Fatalf("the tree shows %v, want RCP01 and RCP02", seen)
	}
}

// TestReadinessAgreesWithAdoptable covers the one predicate.
//
// ReadinessOf.Ready is true iff the recipe is adoptable, including an approved
// recipe that waits on an open one. The reason names that dependency.
func TestReadinessAgreesWithAdoptable(t *testing.T) {
	plan := []Recipe{
		{ID: "RCP01", Lifecycle: LifecycleClosed},
		{ID: "RCP02", Lifecycle: LifecycleApproved},
		{ID: "RCP03", Lifecycle: LifecycleApproved, DependsOn: []string{"RCP01"}},
		{ID: "RCP04", Lifecycle: LifecycleApproved, DependsOn: []string{"RCP09"}},
		{ID: "RCP05", Lifecycle: LifecycleDraft},
		{ID: "RCP06", Lifecycle: LifecycleExecuting},
		{ID: "RCP07", Lifecycle: LifecycleApproved, Overlays: Overlays{Cancelled: &Overlay{Reason: "no"}}},
	}
	in := map[string]bool{}
	for _, r := range Adoptable(plan) {
		in[r.ID] = true
	}
	for _, r := range plan {
		ready := ReadinessOf(r, plan)
		if ready.Ready != in[r.ID] {
			t.Errorf("%s: ready %v, adoptable %v (%q)", r.ID, ready.Ready, in[r.ID], ready.Reason)
		}
	}
	waiting := ReadinessOf(plan[3], plan)
	if waiting.Ready {
		t.Fatal("an approved recipe with an unmet dependency reads as ready")
	}
	if !strings.Contains(waiting.Reason, "RCP09") {
		t.Fatalf("reason = %q, want it to name RCP09", waiting.Reason)
	}
}

// TestTheAdoptableViewIsHeldInBothDirections covers the wire.
//
// The view is registered in core.DerivedViews and written down in docs/03 § 8.
// The wire fails when either side goes.
func TestTheAdoptableViewIsHeldInBothDirections(t *testing.T) {
	doc := readDoc(t, "03-domain-model.md")

	t.Run("the row deleted from the table", func(t *testing.T) {
		broken := strings.Replace(doc, "| `adoptable` | `Adoptable` |", "| `invented` | `Adoptable` |", 1)
		if broken == doc {
			t.Fatal("the table does not hold the row — the corruption did not apply")
		}
		if err := checkDerivedViews(broken); err == nil {
			t.Fatal("the wire should fail when the table does not name the view")
		}
	})

	t.Run("the registration removed", func(t *testing.T) {
		orig := DerivedViews
		t.Cleanup(func() { DerivedViews = orig })

		views := make([]DerivedView, 0, len(orig))
		for _, v := range orig {
			if v.Name == "adoptable" {
				continue
			}
			views = append(views, v)
		}
		DerivedViews = views

		if err := checkViewCoverage("."); err == nil {
			t.Fatal("the wire should fail when a view is registered nowhere")
		}
	})
}
