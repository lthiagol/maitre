package core

import (
	"errors"
	"testing"
)

// groomable returns a recipe that satisfies every groom gate.
func groomable() *Recipe {
	return &Recipe{
		ID:        "R01",
		Title:     "Ship the board",
		Lifecycle: LifecycleDraft,
		Intent: Intent{
			Outcome: "A human can see the plan at a glance",
			Problem: "v1 needed 47 commands for one workflow",
		},
		Scope: Scope{
			InScope:    []string{"the board"},
			OutOfScope: []string{"web UI", "orchestration"},
		},
		Acceptance: []Criterion{{
			ID:           "AC-01",
			Description:  "the board renders 100 recipes",
			Verification: "go test ./internal/core -run Board",
		}},
	}
}

// groomed returns a recipe that satisfies every groom gate, in the groomed state.
// readySpec is a recipe that passes the groom, specify and approve gates.
func readySpec() RecipeSpec {
	return RecipeSpec{
		Title:      "Sharpen the AC gate",
		Outcome:    "a zero-match evidence never satisfies a gate",
		Problem:    "v1 manufactured confidence",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the board renders", Verification: "go test ./internal/core -run Board"},
		},
		WorkPackages: []string{"core"},
		Steps: []StepSpec{{
			Action: "add the gate", DoneWhen: "a rejected evidence blocks", Covers: []string{"AC-01"},
		}},
	}
}

func groomed() *Recipe {
	r := groomable()
	r.Lifecycle = LifecycleGroomed
	return r
}

// ready returns a recipe that satisfies the specify gates, in the ready state.
func ready() *Recipe {
	r := groomable()
	r.Lifecycle = LifecycleReady
	r.WorkPackages = []WorkPackage{{ID: "WP1", Name: "core"}}
	r.Guidance = []Step{{
		ID:       "G1",
		Action:   "implement the board",
		DoneWhen: "the board renders",
		Covers:   []string{"AC-01"},
		Status:   "pending",
	}}
	return r
}

func TestGroomRequiresIntentProblemScopeAndCriteria(t *testing.T) {
	r := groomable()
	r.Intent.Problem = ""
	r.Scope.OutOfScope = nil
	r.Acceptance = nil

	err := Transition(r, EventGroom, TransitionContext{})
	var ge *GateError
	if !errors.As(err, &ge) {
		t.Fatalf("expected a GateError, got %v", err)
	}
	want := map[string]bool{
		"problem.present":             false,
		"scope.out_of_scope_declared": false,
		"ac.present":                  false,
	}
	for _, g := range ge.Gates {
		name := g
		for i := 0; i < len(g); i++ {
			if g[i] == ' ' {
				name = g[:i]
				break
			}
		}
		if _, ok := want[name]; ok {
			want[name] = true
		}
	}
	for name, seen := range want {
		if !seen {
			t.Errorf("expected gate %q to fail", name)
		}
	}
}

func TestGroomRejectsPlaceholderVerification(t *testing.T) {
	r := groomable()
	r.Acceptance[0].Verification = "TODO"

	err := Transition(r, EventGroom, TransitionContext{})
	var ge *GateError
	if !errors.As(err, &ge) {
		t.Fatalf("expected a GateError, got %v", err)
	}
}

func TestDecompositionPrecedesApproval(t *testing.T) {
	// A recipe cannot be approved from draft: the approve event is only legal from ready,
	// and ready requires decomposition. This is the ordering v1 got backwards.
	r := groomable()
	if err := Transition(r, EventApprove, TransitionContext{Approver: "chef"}); err == nil {
		t.Fatal("expected approving a draft to be rejected")
	}
}

func TestSpecifyRequiresCoverageAndDoneWhen(t *testing.T) {
	r := groomed()
	r.WorkPackages = []WorkPackage{{ID: "WP1", Name: "core"}}
	r.Guidance = []Step{{
		ID:       "G1",
		Action:   "implement the board",
		DoneWhen: "the board renders",
		Covers:   []string{"AC-01"},
		Status:   "pending",
	}}
	r.Guidance[0].Covers = nil
	r.Guidance[0].DoneWhen = ""

	err := Transition(r, EventSpecify, TransitionContext{})
	var ge *GateError
	if !errors.As(err, &ge) {
		t.Fatalf("expected a GateError, got %v", err)
	}
}

func TestApprovalRequiresAHuman(t *testing.T) {
	r := ready()
	if err := Transition(r, EventApprove, TransitionContext{}); err == nil {
		t.Fatal("expected approve without an approver to be rejected")
	}
	if err := Transition(r, EventApprove, TransitionContext{Approver: "chef"}); err != nil {
		t.Fatalf("expected approve to succeed: %v", err)
	}
	if r.Lifecycle != LifecycleApproved {
		t.Fatalf("lifecycle = %q, want approved", r.Lifecycle)
	}
	if r.Approval == nil || r.Approval.By != "chef" {
		t.Fatalf("approval not recorded: %+v", r.Approval)
	}
}

func TestApprovalFreezesTheImmutableSections(t *testing.T) {
	r := ready()
	if err := Transition(r, EventApprove, TransitionContext{Approver: "chef"}); err != nil {
		t.Fatal(err)
	}
	want := r.Approval.Fingerprint

	// Changing intent after approval must be detected.
	r.Intent.Outcome = "something else entirely"
	got, err := r.Fingerprint()
	if err != nil {
		t.Fatal(err)
	}
	if got == want {
		t.Fatal("fingerprint did not change after a frozen section was edited")
	}
}

func TestFinishRequiresProvenCriteria(t *testing.T) {
	r := ready()
	r.Lifecycle = LifecycleExecuting
	r.Guidance[0].Status = "done"

	// No evidence yet: the criterion is pending, so finish must be blocked.
	if err := Transition(r, EventFinish, TransitionContext{}); err == nil {
		t.Fatal("expected finish without an evidence to be rejected")
	}

	// An evidence that matched nothing is not evidence.
	r.Acceptance[0].Evidence = []Evidence{{Command: "go test", ExitCode: 0, Matched: 0, Match: MatchRejected}}
	if err := Transition(r, EventFinish, TransitionContext{}); err == nil {
		t.Fatal("expected finish with a zero-match evidence to be rejected")
	}

	// An evidence the tool could not attribute is not evidence either.
	r.Acceptance[0].Evidence = []Evidence{{Command: "make check", ExitCode: 0, Match: MatchUnknown}}
	if err := Transition(r, EventFinish, TransitionContext{}); err == nil {
		t.Fatal("expected finish with an unproven criterion to be rejected")
	}

	// An evidence that matched something passes.
	r.Acceptance[0].Evidence = []Evidence{{Command: "go test", ExitCode: 0, Matched: 3, Match: MatchProven}}
	if err := Transition(r, EventFinish, TransitionContext{}); err != nil {
		t.Fatalf("expected finish to succeed: %v", err)
	}
	if r.Lifecycle != LifecycleExecuted {
		t.Fatalf("lifecycle = %q, want executed", r.Lifecycle)
	}
}

func TestEvidenceIsAppendOnlyAndStatusIsDerived(t *testing.T) {
	c := Criterion{ID: "AC-01"}
	if c.Status() != CriterionPending {
		t.Fatalf("status = %q, want pending", c.Status())
	}
	c.Evidence = append(c.Evidence, Evidence{ExitCode: 1, Match: MatchUnknown})
	if c.Status() != CriterionFailed {
		t.Fatalf("status = %q, want failed", c.Status())
	}
	c.Evidence = append(c.Evidence, Evidence{ExitCode: 0, Matched: 2, Match: MatchProven})
	if c.Status() != CriterionPassed {
		t.Fatalf("status = %q, want passed", c.Status())
	}
	if len(c.Evidence) != 2 {
		t.Fatalf("evidence = %d, want 2 — the earlier evidence must be kept", len(c.Evidence))
	}
}

func TestIllegalTransitionIsRejected(t *testing.T) {
	r := groomable()
	err := Transition(r, EventClose, TransitionContext{})
	var te *TransitionError
	if !errors.As(err, &te) {
		t.Fatalf("expected a TransitionError, got %v", err)
	}
}

func TestAbortRequiresAReason(t *testing.T) {
	r := ready()
	r.Lifecycle = LifecycleExecuting
	if err := Transition(r, EventAbortToDraft, TransitionContext{}); err == nil {
		t.Fatal("expected abort without a reason to be rejected")
	}
	if err := Transition(r, EventAbortToDraft, TransitionContext{Reason: "premise was wrong"}); err != nil {
		t.Fatalf("expected abort to succeed: %v", err)
	}
	if r.Lifecycle != LifecycleDraft {
		t.Fatalf("lifecycle = %q, want draft", r.Lifecycle)
	}
}

func TestCancelledRecipeIsTerminal(t *testing.T) {
	r := groomable()
	if err := Cancel(r, "not needed", "chef", zeroTime()); err != nil {
		t.Fatal(err)
	}
	if err := Transition(r, EventGroom, TransitionContext{}); err == nil {
		t.Fatal("expected a cancelled recipe to reject transitions")
	}
}

func TestEveryStateIsProducible(t *testing.T) {
	// Each state in the table must be reachable from some predecessor. This catches the
	// v1 failure where two states existed as strings that no transition produced.
	reached := map[Lifecycle]bool{}
	for _, def := range transitions {
		reached[def.to] = true
	}
	for _, l := range AllLifecycles {
		if !reached[l] {
			t.Errorf("state %q is not produced by any transition", l)
		}
	}
}

func TestEveryGateIsDefined(t *testing.T) {
	for ev, def := range transitions {
		for _, name := range def.gates {
			if _, ok := gateByName[name]; !ok {
				t.Errorf("event %q references undefined gate %q", ev, name)
			}
		}
	}
}
