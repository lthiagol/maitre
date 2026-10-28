package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newOps(t *testing.T) (*Operations, *Store) {
	t.Helper()
	root := filepath.Join(t.TempDir(), PlanDirName)
	store := Open(root)
	if err := store.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	ops := NewOperations(store, AuditContext{Actor: "chef", Session: "s-1", Role: "chef"})
	return ops, store
}

// spec is a recipe that passes the groom and specify gates.
func spec() RecipeSpec {
	return RecipeSpec{
		Title:      "Sharpen the AC gate",
		Outcome:    "a zero-match result never satisfies a gate",
		Problem:    "v1 manufactured confidence",
		InScope:    []string{"the gate table"},
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the board renders", Verification: "go test ./internal/core -run Board"},
		},
		WorkPackages: []string{"core"},
		Steps: []StepSpec{{
			Action: "add the gate", DoneWhen: "a rejected result blocks", Covers: []string{"AC-01"},
		}},
	}
}

func TestOperationsCreateAndJournal(t *testing.T) {
	ops, store := newOps(t)

	r, err := ops.CreateRecipe(spec())
	if err != nil {
		t.Fatal(err)
	}
	if r.ID != "RCP01" || r.Lifecycle != LifecycleDraft {
		t.Fatalf("recipe = %+v", r)
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) != 1 {
		t.Fatalf("journal has %d events, want 1", len(j.Events))
	}
	e := j.Events[0]
	if e.Kind != KindCreated || e.Recipe != "RCP01" {
		t.Fatalf("event = %+v", e)
	}
	if e.Actor != "chef" || e.Session != "s-1" || e.Role != "chef" {
		t.Fatalf("event is not stamped with the audit context: %+v", e)
	}
	if e.Seq != 1 {
		t.Fatalf("seq = %d, want 1", e.Seq)
	}
}

func TestOperationsAdvanceRunsTheGates(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(spec())
	if err != nil {
		t.Fatal(err)
	}

	// Groom passes.
	if _, err := ops.Advance(r.ID, EventGroom); err != nil {
		t.Fatalf("groom: %v", err)
	}
	// Specify passes.
	if _, err := ops.Advance(r.ID, EventSpecify); err != nil {
		t.Fatalf("specify: %v", err)
	}
	// Approve needs an approver.
	if _, err := ops.ApproveRecipe(r.ID, ""); err == nil {
		t.Fatal("expected approve without an approver to fail")
	}
	if _, err := ops.ApproveRecipe(r.ID, "chef"); err != nil {
		t.Fatalf("approve: %v", err)
	}
}

func TestOperationsRefuseGatesThatDoNotPass(t *testing.T) {
	ops, _ := newOps(t)
	bad := spec()
	bad.OutOfScope = nil // too few out-of-scope items
	r, err := ops.CreateRecipe(bad)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventGroom); err == nil {
		t.Fatal("expected groom to be refused")
	} else if !strings.Contains(err.Error(), "out_of_scope") {
		t.Fatalf("error should name the failed gate: %v", err)
	}
}

func TestOperationsFreezeSectionsAfterApproval(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(spec())
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

	// The frozen sections are closed.
	outcome := "something else"
	if _, err := ops.UpdateRecipe(r.ID, Fragment{Outcome: &outcome}); err == nil {
		t.Fatal("expected the intent to be frozen after sign-off")
	}
	if _, err := ops.AddCriterion(r.ID, "a new one", "go test ./..."); err == nil {
		t.Fatal("expected criteria to be frozen after sign-off")
	}

	// Guidance is not frozen.
	if _, err := ops.AddStep(r.ID, StepSpec{Action: "another step", DoneWhen: "it works"}); err != nil {
		t.Fatalf("guidance should stay revisable: %v", err)
	}
}

func TestRequestChangesSendsItBackToShaping(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(spec())
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := ops.Advance(r.ID, ev); err != nil {
			t.Fatal(err)
		}
	}
	// A reason is required.
	if _, err := ops.RequestChanges(r.ID, "chef", ""); err == nil {
		t.Fatal("expected request-changes without a reason to fail")
	}
	back, err := ops.RequestChanges(r.ID, "chef", "the scope is too wide")
	if err != nil {
		t.Fatalf("request changes: %v", err)
	}
	if back.Lifecycle != LifecycleGroomed {
		t.Fatalf("lifecycle = %q, want groomed", back.Lifecycle)
	}
}

func TestOperationsStampEveryWrite(t *testing.T) {
	ops, store := newOps(t)
	r, err := ops.CreateRecipe(spec())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventGroom); err != nil {
		t.Fatal(err)
	}
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) != 2 {
		t.Fatalf("journal has %d events, want 2", len(j.Events))
	}
	for _, e := range j.Events {
		if e.Actor != "chef" || e.Session != "s-1" {
			t.Fatalf("event not stamped: %+v", e)
		}
	}
	if j.Events[1].Seq != 2 {
		t.Fatalf("seq = %d, want 2", j.Events[1].Seq)
	}
}

func TestJournalKindIsAClosedSet(t *testing.T) {
	store := Open(filepath.Join(t.TempDir(), PlanDirName))
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendJournalEvent(JournalEntry{Kind: "not-a-kind"}); err == nil {
		t.Fatal("expected an unknown journal kind to be refused")
	}
}

func TestOperationsJournalIsTheOnlyEventLog(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(spec())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventGroom); err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	// The recipe carries its authored record, not a copy of its events.
	if len(got.Audit.Deviations) != 0 || len(got.Audit.Decisions) != 0 {
		t.Fatalf("recipe audit should hold only authored records: %+v", got.Audit)
	}
}

func TestSeedProducesARealisticPlan(t *testing.T) {
	root := filepath.Join(t.TempDir(), PlanDirName)
	store := Open(root)
	if err := Seed(store); err != nil {
		t.Fatalf("seed: %v", err)
	}

	recipes, err := store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if len(recipes) != 7 {
		t.Fatalf("seeded %d recipes, want 7", len(recipes))
	}

	want := map[string]Lifecycle{
		"RCP01": LifecycleDraft,
		"RCP02": LifecycleReady,
		"RCP03": LifecycleRemediation,
		"RCP04": LifecycleApproved,
		"RCP05": LifecycleExecuting,
		"RCP06": LifecycleClosed,
		"RCP07": LifecycleDraft,
	}
	for _, r := range recipes {
		if got := want[r.ID]; got != r.Lifecycle {
			t.Errorf("%s (%s) lifecycle = %q, want %q", r.ID, r.Title, r.Lifecycle, got)
		}
	}

	// The lanes come out as the mockups show.
	lanes := map[string]Lane{}
	for _, r := range recipes {
		lanes[r.ID] = LaneOf(r)
	}
	for id, want := range map[string]Lane{
		"RCP01": LaneGrooming, "RCP02": LaneApproval, "RCP03": LaneReview,
		"RCP04": LaneExecution, "RCP05": LaneExecution, "RCP06": LaneDone, "RCP07": LaneGrooming,
	} {
		if lanes[id] != want {
			t.Errorf("%s lane = %q, want %q", id, lanes[id], want)
		}
	}

	// The journal is populated and sequenced.
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) < 20 {
		t.Fatalf("journal has %d events, want a realistic trail", len(j.Events))
	}
	for i, e := range j.Events {
		if e.Seq != i+1 {
			t.Fatalf("event %d has seq %d", i, e.Seq)
		}
		if e.Actor == "" {
			t.Fatalf("event %d has no actor", i)
		}
	}
}

func TestSeedIsIdempotentEnoughToRerun(t *testing.T) {
	root := filepath.Join(t.TempDir(), PlanDirName)
	store := Open(root)
	if err := Seed(store); err != nil {
		t.Fatal(err)
	}
	// Seeding a second time must not crash or corrupt the plan.
	if err := Seed(store); err != nil {
		t.Fatalf("second seed: %v", err)
	}
	recipes, err := store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if len(recipes) == 0 {
		t.Fatal("no recipes after a second seed")
	}
}

func TestLockIsReleasedAfterUse(t *testing.T) {
	ops, store := newOps(t)
	if _, err := ops.CreateRecipe(spec()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(store.Root, LockName)); !os.IsNotExist(err) {
		t.Fatal("the lock file should be gone after the operation")
	}
}

func TestErrorsAreDescriptive(t *testing.T) {
	ops, _ := newOps(t)
	if _, err := ops.GetRecipe("R99"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
