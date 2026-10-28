package core

import (
	"os"
	"path/filepath"
	"testing"
)

// recipeReviewedBy builds a recipe that was executed by `claimer` and reviewed by
// `reviewer`.
func recipeReviewedBy(t *testing.T, claimer, reviewer string) (*Operations, *Recipe) {
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
	if _, err := ops.ClaimRecipe(r.ID, claimer); err != nil {
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
	if _, err := ops.RecordReviewRound(r.ID, reviewer, "ok", "looks good to me"); err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	return ops, got
}

// TestTheRenamedRolesAreTheAgentRoles covers the deny-list after the rename.
//
// The crew's jobs are plain words now: the sous chef, the runner, the reviewer and
// the orchestrator. The list is a deny list, and the names the rename retired are
// **on it** — they are names the list refuses, and nothing more. A retired name is
// not the runner any more; the mapping that said so was a migration affordance.
func TestTheRenamedRolesAreTheAgentRoles(t *testing.T) {
	for _, role := range []string{"agent", "sous-chef", "runner", "reviewer", "orchestrator"} {
		if !IsAgentRole(role) {
			t.Errorf("%q should be an agent role", role)
		}
	}
	// The names the rename retired: the list refuses them as names.
	for _, retired := range []string{"line-cook", "Line Cook", "taster", "expo", "coordinator"} {
		if !IsAgentRole(retired) {
			t.Errorf("%q should be an agent role", retired)
		}
	}
	for _, human := range []string{"chef", "maitre", "alice"} {
		if IsAgentRole(human) {
			t.Errorf("%q should not be an agent role", human)
		}
	}
}

// TestTheRunnerAndTheReviewerAreTwoAgents covers the other side.
//
// `runner` is the runner and `reviewer` is the reviewer: two jobs, two agents. The
// rename must not make every review a self-review.
func TestTheRunnerAndTheReviewerAreTwoAgents(t *testing.T) {
	ops, r := recipeReviewedBy(t, "runner", "reviewer")
	if SelfReviewed(*r) {
		t.Fatal("the runner's work reviewed by the reviewer is not a self-review")
	}
	if _, err := ops.Advance(r.ID, EventPassReview); err != nil {
		t.Fatalf("an independent review should pass: %v", err)
	}
}

// TestARetiredNameIsNotTheRunner covers the mapping that is gone.
//
// A retired name is a name the deny-list refuses — it is not resolved to the role it
// stood for. So the domain cannot tell that `line-cook` and `runner` are one agent:
// it compares the names it is given.
func TestARetiredNameIsNotTheRunner(t *testing.T) {
	// Both names are refused: the retired names are on the deny-list.
	if !IsAgentRole("runner") || !IsAgentRole("line-cook") {
		t.Fatal("both the runner and the retired name are refused as agent roles")
	}
	// But they are not resolved to each other: a retired name is a name.
	if sameAgent("line-cook", "runner") {
		t.Fatal("a retired name must not resolve to the role it stood for")
	}
	if !sameAgent("Line Cook", "line-cook") {
		t.Fatal("the same name written twice is one name")
	}
}

// TestAnOldRoleNameStillReads covers the journal.
//
// The journal is append-only: a plan written before the rename keeps its old role
// names, and a reader meeting one must not fail.
func TestAnOldRoleNameStillReads(t *testing.T) {
	_, store := newOps(t)
	if _, err := store.AppendJournalEvent(JournalEntry{
		Seq: 1, At: Now(), Kind: "control", Actor: "line-cook", Session: "s-old", Role: "line-cook",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendJournalEvent(JournalEntry{
		Seq: 2, At: Now(), Kind: "control", Actor: "runner", Session: "s-new", Role: "runner",
	}); err != nil {
		t.Fatal(err)
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatalf("a journal with an old role name should read: %v", err)
	}
	if len(j.Events) != 2 {
		t.Fatalf("events = %d, want 2", len(j.Events))
	}
	if j.Events[0].Role != "line-cook" {
		t.Errorf("the old role name should be kept as it was written: %q", j.Events[0].Role)
	}
	if j.Events[1].Role != "runner" {
		t.Errorf("role = %q, want runner", j.Events[1].Role)
	}
}

// TestThePlanJournalsTheNewRoles covers what the tool writes from now on.
//
// The journal records the role of every write, and the roles it records are the plain
// words the rename settled on.
func TestThePlanJournalsTheNewRoles(t *testing.T) {
	ops, store := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range j.Events {
		if e.Recipe != r.ID {
			continue
		}
		found = true
		if e.Role != "chef" {
			t.Errorf("event %d role = %q, want chef", e.Seq, e.Role)
		}
	}
	if !found {
		t.Error("the recipe's creation should be in the journal")
	}
}

// TestTheOldPlanStillReads proves the repository's own plan survives the rename.
//
// The plan was written under the old names and its journal keeps them: the tool reads
// it, and a reader meeting a retired role name does not fail.
func TestTheOldPlanStillReads(t *testing.T) {
	root := filepath.Join("..", "..", "kitchen")
	if _, err := os.Stat(root); err != nil {
		t.Skipf("no plan at %s", root)
	}
	store := Open(root)
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatalf("the old plan's journal should read: %v", err)
	}
	retired := 0
	for _, e := range j.Events {
		if e.Role == "coordinator" || e.Role == "line-cook" || e.Role == "taster" || e.Role == "expo" {
			retired++
		}
	}
	if retired == 0 {
		t.Skip("the old plan's journal carries no retired role name")
	}
}
