package core

import (
	"os"
	"path/filepath"
	"strings"
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
// A retired name is a name the tool refuses — not the runner, and not any other
// role. The refusal is not a special case for the approval: a retired name is
// refused wherever an actor names itself.
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

// TestARetiredNameCannotNameItself covers the rule.
//
// A retired name is refused wherever an actor names itself: the claim, the review,
// the approval, the sign-off, the comment — and the actor every write is stamped
// with. A name the tool does not accept cannot buy a gate.
//
// The name is refused **however it is spelled**: the case does not matter, and a
// space or an underscore stands for a dash. "line_cook" is the same name as
// "line-cook", and a spelling is not a way around the rule.
func TestARetiredNameCannotNameItself(t *testing.T) {
	retired := []string{"line-cook", "line_cook", "Line Cook", "taster", "expo", "coordinator"}

	for _, name := range retired {
		t.Run(name, func(t *testing.T) {
			ops, _ := newOps(t)
			r, err := ops.CreateRecipe(readySpec())
			if err != nil {
				t.Fatal(err)
			}

			// The claim.
			if _, err := ops.ClaimRecipe(r.ID, name); err == nil {
				t.Errorf("a retired name must not be able to claim: %q", name)
			}
			// The review.
			if _, err := ops.RecordReviewRound(r.ID, name, "ok", "looks good"); err == nil {
				t.Errorf("a retired name must not be able to review: %q", name)
			}
			// The approval.
			if _, err := ops.ApproveRecipe(r.ID, name); err == nil {
				t.Errorf("a retired name must not be able to approve: %q", name)
			}
			// The sign-off.
			if _, err := ops.SignOffCriterion(r.ID, "AC-01", name, "checked by hand"); err == nil {
				t.Errorf("a retired name must not be able to sign off: %q", name)
			}
			// The comment.
			if _, err := ops.AddComment(r.ID, name, "a note", ""); err == nil {
				t.Errorf("a retired name must not be able to comment: %q", name)
			}
		})
	}
}

// TestARetiredActorCannotWrite covers the actor.
//
// Every write is stamped with the actor, so the rule has to hold there too: a
// caller that names itself with a retired name cannot write to the plan at all.
func TestARetiredActorCannotWrite(t *testing.T) {
	store := Open(t.TempDir() + "/" + PlanDirName)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := NewOperations(store, AuditContext{Actor: "line-cook", Session: "s-old", Role: "runner"})

	_, err := ops.CreateRecipe(readySpec())
	if err == nil {
		t.Fatal("a retired name must not be able to write")
	}
	if !strings.Contains(err.Error(), "line-cook") {
		t.Fatalf("the refusal should name the actor: %v", err)
	}
	if !strings.Contains(err.Error(), "not a name this tool accepts") {
		t.Fatalf("the refusal should say what the rule is: %v", err)
	}
}

// TestARetiredRoleCannotWrite covers the audit role.
//
// The role lands in the journal's event, and a **new** event is a live record. The
// check is at the write's entry — before anything is written — so a retired role
// never leaves a half-written plan behind.
func TestARetiredRoleCannotWrite(t *testing.T) {
	for _, role := range []string{"line-cook", "line_cook", "taster"} {
		t.Run(role, func(t *testing.T) {
			store := Open(t.TempDir() + "/" + PlanDirName)
			if err := store.Init(); err != nil {
				t.Fatal(err)
			}
			ops := NewOperations(store, AuditContext{Actor: "chef", Session: "s-old", Role: role})

			if _, err := ops.CreateRecipe(readySpec()); err == nil {
				t.Fatalf("a retired role must not be able to write: %q", role)
			}
			// And nothing landed.
			entries, err := os.ReadDir(store.RecipeDir())
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 0 {
				t.Errorf("the refusal wrote %d files anyway", len(entries))
			}
		})
	}
}

// TestARetiredNameCannotEnterARecord covers the rule the review found.
//
// **A retired name cannot enter a live record, by any path.** The rule is not one
// check per operation: it is at the entry, where the name is recorded — the claim,
// the review, the approval, the sign-off, the comment, the transition's context
// and the audit role.
//
// The hole was the **second road to the approval**: `AdvanceWith` records the
// approver the transition's context carries, and it did not check it. The approval
// went in with `"by": "line-cook"`, and from then on the plan was refused by the
// reader *and* by the writer: the tool bricked the plan it had just written.
func TestARetiredNameCannotEnterARecord(t *testing.T) {
	for _, name := range []string{"line-cook", "line_cook", "Line Cook", "taster", "expo", "coordinator"} {
		t.Run(name, func(t *testing.T) {
			ops, store := newOps(t)
			r, err := ops.CreateRecipe(readySpec())
			if err != nil {
				t.Fatal(err)
			}
			for _, ev := range []Event{EventGroom, EventSpecify} {
				if _, err := ops.Advance(r.ID, ev); err != nil {
					t.Fatal(err)
				}
			}

			// The second road to the approval: the transition's context.
			if _, err := ops.AdvanceWith(r.ID, EventApprove, TransitionContext{
				Approver: name, Actor: "chef",
			}); err == nil {
				t.Fatalf("a retired name must not be able to approve: %q", name)
			}

			// Nothing was written: the plan is not bricked, and it still reads.
			read, err := store.LoadRecipe(r.ID)
			if err != nil {
				t.Fatalf("the plan must still read: %v", err)
			}
			if read.Approval != nil {
				t.Error("the refusal wrote the approval anyway")
			}
			if read.Lifecycle != LifecycleReady {
				t.Errorf("the refusal moved the recipe to %q", read.Lifecycle)
			}
		})
	}
}

// TestARetiredNameCannotRequestChanges covers the other road.
//
// `RequestChanges` did not check the name it was given: the transition drops it —
// the objection is recorded, not the name — but a name the tool does not accept is
// refused wherever it is offered.
func TestARetiredNameCannotRequestChanges(t *testing.T) {
	ops, store := newOps(t)
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

	if _, err := ops.RequestChanges(r.ID, "line-cook", "the scope is wrong"); err == nil {
		t.Fatal("a retired name must not be able to request changes")
	}
	read, err := store.LoadRecipe(r.ID)
	if err != nil {
		t.Fatalf("the plan must still read: %v", err)
	}
	if read.Lifecycle != LifecycleApproved {
		t.Errorf("the refusal moved the recipe to %q", read.Lifecycle)
	}
}

// TestARetiredNameCannotBuyTheIndependenceGate covers the hole.
//
// The executor claims as `line-cook` and the runner reviews as `runner`: with the
// mapping gone, two names stood for one role, and the gate read the review as
// independent. The claim is refused at the door — and a record that carries the
// retired name is refused by the reader.
func TestARetiredNameCannotBuyTheIndependenceGate(t *testing.T) {
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

	// The claim, under the retired name: refused.
	if _, err := ops.ClaimRecipe(r.ID, "line-cook"); err == nil {
		t.Fatal("a retired name must not be able to claim")
	}

	// And the record that carries it: the reader refuses the plan.
	//
	// The record is written **by hand**: the writer refuses a name the tool does
	// not accept, so the tool cannot write this record — and a reader meeting one
	// must refuse it.
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	writeRecipeFile(t, s, "RCP01-claimed-under-the-old-name.json", `{"id":"RCP01",
		"title":"claimed under the old name","slug":"claimed-under-the-old-name",
		"lifecycle":"executing","execution":{"claimed_by":"line-cook","claimed_at":"2026-10-26T00:00:00Z"}}`)
	if _, err := s.LoadRecipe("RCP01"); err == nil {
		t.Fatal("a record carrying a retired name must not be read")
	}
	if _, err := s.ListRecipes(); err == nil {
		t.Fatal("a plan carrying a record with a retired name must not be listed")
	}

	// And the writer refuses it too: the domain cannot write a plan it will not
	// read.
	claimed := &Recipe{
		ID:        "RCP02",
		Title:     "claimed under the old name",
		Slug:      "claimed-under-the-old-name",
		Lifecycle: LifecycleExecuting,
		Execution: Execution{ClaimedBy: "line-cook", ClaimedAt: Now()},
	}
	if err := s.SaveRecipe(claimed); err == nil {
		t.Fatal("the writer must refuse a record carrying a retired name")
	}
}

// TestAnOldRoleNameStillReads covers the journal.
//
// The journal is append-only: a plan written before the rename keeps its old role
// names, and a reader meeting one must not fail.
//
// The old event is written **by hand**: a **new** event may not carry a retired
// name — it is a live record — so the tool will not write this one. The history is
// the old tool's.
func TestAnOldRoleNameStillReads(t *testing.T) {
	_, store := newOps(t)
	writeJournalHistory(t, store, JournalEntry{
		Seq: 1, At: Now(), Kind: "control", Actor: "line-cook", Session: "s-old", Role: "line-cook",
		DomainVersion: DomainVersion,
	})
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

// TestANewEventMayNotCarryARetiredName covers the other side.
//
// The journal's **history** keeps its old names; a **new** event is a live record,
// and a name the tool does not accept cannot be in one.
func TestANewEventMayNotCarryARetiredName(t *testing.T) {
	_, store := newOps(t)
	for _, name := range []string{"line-cook", "line_cook", "Line Cook"} {
		t.Run(name, func(t *testing.T) {
			_, err := store.AppendJournalEvent(JournalEntry{
				Seq: 1, At: Now(), Kind: "control", Actor: name, Session: "s-new", Role: "runner",
			})
			if err == nil {
				t.Fatalf("a new event must not carry the retired name %q", name)
			}
			if !strings.Contains(err.Error(), "not a name this tool accepts") {
				t.Errorf("the refusal should say what the rule is: %v", err)
			}
			// And nothing was written.
			j, err := store.LoadJournal()
			if err != nil {
				t.Fatal(err)
			}
			if len(j.Events) != 0 {
				t.Errorf("the refusal wrote %d events anyway", len(j.Events))
			}
		})
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
