package core

import (
	"encoding/json"
	"strings"
	"testing"
)

// approvedRecipe builds a recipe that is signed off and waiting to be claimed.
func approvedRecipe(t *testing.T, ops *Operations, title string) *Recipe {
	t.Helper()
	s := readySpec()
	s.Title = title
	r, err := ops.CreateRecipe(s)
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
	return r
}

// closeRecipe walks a recipe all the way to closed.
func closeRecipe(t *testing.T, ops *Operations, id string) {
	t.Helper()
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
	if _, err := ops.DocumentRecipe(id, "done and written up"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(id, EventClose); err != nil {
		t.Fatal(err)
	}
}

// TestClaimRefusesAnOpenDependency covers the dependency rule on the claim path.
//
// The claim used to assign `executing` outright, so a recipe whose dependency was
// still open could go on the line while the same recipe could not be started with
// `advance start`. Now the claim is a transition, and the `dependencies.satisfied`
// gate runs on it.
func TestClaimRefusesAnOpenDependency(t *testing.T) {
	ops, _ := newOps(t)
	dep := approvedRecipe(t, ops, "The dependency")
	r := approvedRecipe(t, ops, "The dependent")
	if _, err := ops.UpdateRecipe(r.ID, Fragment{DependsOn: &[]string{dep.ID}}); err != nil {
		t.Fatal(err)
	}

	_, err := ops.ClaimRecipe(r.ID, "runner")
	if err == nil {
		t.Fatal("a recipe whose dependency is open must not be claimed")
	}
	if !strings.Contains(err.Error(), dep.ID) {
		t.Fatalf("the refusal should name the dependency: %v", err)
	}

	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Lifecycle != LifecycleApproved {
		t.Fatalf("lifecycle = %q, want approved — a refused claim changes nothing", got.Lifecycle)
	}
	if got.Execution.ClaimedBy != "" {
		t.Fatalf("a refused claim must not record a claimant: %q", got.Execution.ClaimedBy)
	}
}

// TestClaimRunsWhenTheDependencyIsClosed is the other side: the gate is a check, not
// a wall.
//
// The claim must succeed when the dependency is closed — and it must succeed *through
// the transition*, so the journal records the state change as a transition rather
// than a bare assignment.
func TestClaimRunsWhenTheDependencyIsClosed(t *testing.T) {
	ops, store := newOps(t)
	dep := approvedRecipe(t, ops, "The dependency")
	closeRecipe(t, ops, dep.ID)

	r := approvedRecipe(t, ops, "The dependent")
	if _, err := ops.UpdateRecipe(r.ID, Fragment{DependsOn: &[]string{dep.ID}}); err != nil {
		t.Fatal(err)
	}

	claimed, err := ops.ClaimRecipe(r.ID, "runner")
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if claimed.Lifecycle != LifecycleExecuting {
		t.Fatalf("lifecycle = %q, want executing", claimed.Lifecycle)
	}
	if claimed.Execution.ClaimedBy != "runner" {
		t.Fatalf("claimant = %q, want runner", claimed.Execution.ClaimedBy)
	}
	if claimed.Execution.ClaimedAt == "" {
		t.Fatal("the claim must record when it happened")
	}

	// The state change went through the transition table.
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	last := j.Events[len(j.Events)-1]
	if last.Kind != KindLifecycle {
		t.Fatalf("kind = %q, want %q — the claim must be a transition", last.Kind, KindLifecycle)
	}
}

// TestClaimIsJournalledAsALifecycleEvent covers the journal's shape.
//
// The journal is the only event log. A claim journalled as `control` with no payload
// was invisible to an agent reading the log: the recipe went on the line and the log
// said nothing about it.
func TestClaimIsJournalledAsALifecycleEvent(t *testing.T) {
	ops, store := newOps(t)
	r := approvedRecipe(t, ops, "A recipe")
	if _, err := ops.ClaimRecipe(r.ID, "runner"); err != nil {
		t.Fatal(err)
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) == 0 {
		t.Fatal("the journal is empty")
	}
	last := j.Events[len(j.Events)-1]
	if last.Kind != KindLifecycle {
		t.Fatalf("kind = %q, want %q", last.Kind, KindLifecycle)
	}
	if last.Recipe != r.ID {
		t.Fatalf("recipe = %q, want %q", last.Recipe, r.ID)
	}
	var data map[string]any
	if err := json.Unmarshal(last.Data, &data); err != nil {
		t.Fatalf("the event carries no data: %v", err)
	}
	if data["from"] != string(LifecycleApproved) {
		t.Fatalf("from = %v, want approved", data["from"])
	}
	if data["to"] != string(LifecycleExecuting) {
		t.Fatalf("to = %v, want executing", data["to"])
	}
}

// TestClaimNeedsAnAgent covers the anonymous claim.
//
// An empty claimant made the `review.reviewer_independent` gate vacuous: the
// executor could record its own "ok" round and pass review. The gate the project
// built to stop that was defeated by omitting one flag.
func TestClaimNeedsAnAgent(t *testing.T) {
	ops, _ := newOps(t)
	r := approvedRecipe(t, ops, "A recipe")

	if _, err := ops.ClaimRecipe(r.ID, ""); err == nil {
		t.Fatal("a claim with no agent must be refused")
	}
	if _, err := ops.ClaimRecipe(r.ID, "   "); err == nil {
		t.Fatal("a claim with a blank agent must be refused")
	}

	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Lifecycle != LifecycleApproved {
		t.Fatalf("lifecycle = %q, want approved — a refused claim changes nothing", got.Lifecycle)
	}
}

// TestTheDependenciesGateRefusesWhenItCannotCheck covers the gate that cannot verify.
//
// The gate table is package-level and cannot close over the store, so the gate asks
// the transition context. A context that carries no way to read the plan cannot
// answer the question — and a gate that cannot verify must refuse, not wave the
// transition through. v1's lesson: a green gate that proves nothing.
func TestDependenciesGateRefusesWhenItCannotCheck(t *testing.T) {
	r := &Recipe{
		ID:        "R01",
		Title:     "A recipe",
		Lifecycle: LifecycleApproved,
		DependsOn: []string{"R02"},
	}
	err := Transition(r, EventStart, TransitionContext{})
	if err == nil {
		t.Fatal("the gate should refuse when it cannot check the dependency")
	}
	if !strings.Contains(err.Error(), "R02") {
		t.Fatalf("the refusal should name the dependency: %v", err)
	}
	if r.Lifecycle != LifecycleApproved {
		t.Fatalf("lifecycle = %q, want approved — a refused transition changes nothing", r.Lifecycle)
	}
}

// TestTheClaimIsTheOnlyRoadToTheLine covers the one-path rule.
//
// `advance start` and `claim` are two verbs for one move, so both go through the
// claim: the dependency gate runs on both, and both record a claimant.
func TestClaimIsTheOnlyRoadToTheLine(t *testing.T) {
	ops, _ := newOps(t)
	dep := approvedRecipe(t, ops, "The dependency")
	r := approvedRecipe(t, ops, "The dependent")
	if _, err := ops.UpdateRecipe(r.ID, Fragment{DependsOn: &[]string{dep.ID}}); err != nil {
		t.Fatal(err)
	}

	// The other road to `executing` is guarded the same way.
	if _, err := ops.Advance(r.ID, EventStart); err == nil {
		t.Fatal("advance start must run the dependency gate too")
	} else if !strings.Contains(err.Error(), dep.ID) {
		t.Fatalf("the refusal should name the dependency: %v", err)
	}
}
