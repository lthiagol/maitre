package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// signedOffRecipe builds a recipe that is approved, claimed and signed off, with its
// step done — ready to finish.
func signedOffRecipe(t *testing.T, ops *Operations, id string) {
	t.Helper()
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := ops.Advance(id, ev); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := ops.ApproveRecipe(id, "chef"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.ClaimRecipe(id, "runner"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.SetStepStatus(id, "G1", "done"); err != nil {
		t.Fatal(err)
	}
}

// TestADeferredCriterionDoesNotWedgeTheRecipe covers the deferral.
//
// A deferred criterion is re-filed as a follow-up: it does not keep the recipe open.
// The gate used to refuse it forever, and the recipe had no path onward.
func TestADeferredCriterionDoesNotWedgeTheRecipe(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	signedOffRecipe(t, ops, r.ID)

	if _, err := ops.DeferCriterion(r.ID, "AC-01", "not now", "follow-up R02"); err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if st := got.Acceptance[0].Status(); st != CriterionDeferred {
		t.Fatalf("status = %q, want %q", st, CriterionDeferred)
	}

	// The recipe can still be finished.
	if _, err := ops.Advance(r.ID, EventFinish); err != nil {
		t.Fatalf("finish: %v", err)
	}
}

// TestADeferralDoesNotUndoAPass covers the other half.
//
// The deferral re-files the work; it must not silently undo the evidence.
func TestADeferralDoesNotUndoAPass(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	signedOffRecipe(t, ops, r.ID)

	if _, err := ops.SignOffCriterion(r.ID, "AC-01", "chef", "checked"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.DeferCriterion(r.ID, "AC-01", "not now", "follow-up R02"); err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if st := got.Acceptance[0].Status(); st != CriterionPassed {
		t.Fatalf("status = %q, want %q — a deferral must not undo a pass", st, CriterionPassed)
	}
}

// TestOnlyTheChefMayApprove covers the permission deny-list.
//
// The chef is the human; everyone else in the crew is an agent. The list is a deny
// list, so a human may approve under any name — but an agent cannot sign off at all.
func TestOnlyTheChefMayApprove(t *testing.T) {
	store := Open(filepath.Join(t.TempDir(), PlanDirName))
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	line := NewOperations(store, AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})
	r, err := line.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := line.Advance(r.ID, ev); err != nil {
			t.Fatal(err)
		}
	}

	// An agent cannot sign it off.
	_, err = line.ApproveRecipe(r.ID, "chef")
	if err == nil {
		t.Fatal("an agent must not be able to approve")
	}
	if !strings.Contains(err.Error(), "runner") {
		t.Fatalf("the refusal should name the actor: %v", err)
	}
	if !strings.Contains(err.Error(), "chef") {
		t.Fatalf("the refusal should name the rule: %v", err)
	}

	// The chef can.
	chef := NewOperations(store, AuditContext{Actor: "chef", Session: "s-2", Role: "chef"})
	approved, err := chef.ApproveRecipe(r.ID, "chef")
	if err != nil {
		t.Fatalf("the chef should be able to approve: %v", err)
	}
	if approved.Approval == nil || approved.Approval.By != "chef" {
		t.Fatalf("the approval should record the approver: %+v", approved.Approval)
	}
}

// TestTheDefaultActorIsAnAgentRole covers the CLI's default identity.
func TestTheDefaultActorIsAnAgentRole(t *testing.T) {
	if !IsAgentRole("agent") {
		t.Error("the default actor is an agent role")
	}
	for _, role := range []string{"runner", "Line Cook", "sous-chef", "reviewer", "expo"} {
		if !IsAgentRole(role) {
			t.Errorf("%q should be an agent role", role)
		}
	}
	for _, human := range []string{"chef", "maitre", "alice"} {
		if IsAgentRole(human) {
			t.Errorf("%q should not be an agent role", human)
		}
	}
}

// TestASignOffIsNotAMachineEvidence covers the sign-off's own state.
//
// The zero-match rule is the project's answer to v1's worst bug. A sign-off is a human
// judgement, not machine evidence, and it claims no match.
func TestASignOffIsNotAMachineEvidence(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	signedOffRecipe(t, ops, r.ID)

	if _, err := ops.SignOffCriterion(r.ID, "AC-01", "chef", "checked by hand"); err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	c := got.Acceptance[0]
	if st := c.Status(); st != CriterionPassed {
		t.Fatalf("status = %q, want %q", st, CriterionPassed)
	}
	if len(c.Evidence) != 1 {
		t.Fatalf("got %d records of evidence, want 1", len(c.Evidence))
	}
	evidence := c.Evidence[0]
	if evidence.Match != MatchSignedOff {
		t.Fatalf("match = %q, want %q — a sign-off claims no match", evidence.Match, MatchSignedOff)
	}
	if evidence.Matched != 0 {
		t.Fatalf("matched = %d, want 0", evidence.Matched)
	}
	if !evidence.SignedOff {
		t.Fatal("the evidence should be marked as a sign-off")
	}

	// And the gate accepts it.
	if _, err := ops.Advance(r.ID, EventFinish); err != nil {
		t.Fatalf("finish: %v", err)
	}
}

// TestAFalsifiableVerificationIsOneTheToolCanCount covers the narrowed gate.
//
// A verification the tool cannot attribute a count from would produce evidence in the
// `unknown` state, and evidence in that state never satisfies a gate — so the
// criterion could never be finished. The gate says so at grooming time.
func TestAFalsifiableVerificationIsOneTheToolCanCount(t *testing.T) {
	// The shape the tool can count from.
	if !isGoTest("go test ./internal/core -run Board") {
		t.Fatal("a go test command is countable")
	}
	good := Criterion{ID: "AC-01", Verification: "go test ./internal/core -run Board"}
	if err := checkFalsifiable(good); err != nil {
		t.Fatalf("a go test verification should be accepted: %v", err)
	}

	// The shape the tool cannot count from.
	if _, unit, known := deriveMatch("make test", "everything is fine\n"); known {
		t.Fatalf("the tool should not pretend to count %s matches", unit)
	}
	bad := Criterion{ID: "AC-01", Verification: "make test"}
	err := checkFalsifiable(bad)
	if err == nil {
		t.Fatal("a verification the tool cannot count from should be refused")
	}
	if !strings.Contains(err.Error(), "make test") {
		t.Fatalf("the refusal should name the verification: %v", err)
	}

	// A manual procedure is a human's job, and it is accepted.
	manual := Criterion{ID: "AC-01", Verification: "manual: read the board and check the lanes"}
	if err := checkFalsifiable(manual); err != nil {
		t.Fatalf("a manual verification should be accepted: %v", err)
	}
}

// TestTheFingerprintIsAFunctionOfTheContent covers the canonical digest.
//
// The same recipe must hash the same way whatever the load path and whatever build
// wrote the file — `null` and `[]` are the same list.
func TestTheFingerprintIsAFunctionOfTheContent(t *testing.T) {
	base := Recipe{
		ID:        "RCP01",
		Lifecycle: LifecycleApproved,
		Scope:     Scope{InScope: []string{"the gate table"}, OutOfScope: []string{"web shell"}},
	}
	nullForm := base
	nullForm.Scope.InScope = nil
	emptyForm := base
	emptyForm.Scope.InScope = []string{}

	a, err := nullForm.Fingerprint()
	if err != nil {
		t.Fatal(err)
	}
	b, err := emptyForm.Fingerprint()
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Fatalf("the same content hashed two ways: %s and %s", a, b)
	}

	// And through the store: a file that says null and a file that says [].
	store := Open(filepath.Join(t.TempDir(), PlanDirName))
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	for _, inScope := range [][]string{nil, {}} {
		r := base
		r.Scope.InScope = inScope
		if err := store.SaveRecipe(&r); err != nil {
			t.Fatal(err)
		}
		loaded, err := store.LoadRecipe("RCP01")
		if err != nil {
			t.Fatal(err)
		}
		got, err := loaded.Fingerprint()
		if err != nil {
			t.Fatal(err)
		}
		if got != a {
			t.Fatalf("the file's form changed the digest: %s vs %s", got, a)
		}
	}
}

// TestTamperedFrozenSectionsAreCaught covers the immutability gate.
//
// The artifact must not assert a sign-off over content that is no longer there.
func TestTamperedFrozenSectionsAreCaught(t *testing.T) {
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

	// A hand edit to the frozen sections — the operations refuse one, but the plan
	// directory is hand-editable.
	raw, err := store.LoadRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	raw.Intent.Outcome = "TOTALLY DIFFERENT NOW"
	if err := store.SaveRecipe(raw); err != nil {
		t.Fatal(err)
	}

	// The next transition is refused.
	_, err = ops.ClaimRecipe(r.ID, "runner")
	if err == nil {
		t.Fatal("a recipe whose frozen sections changed must not be claimed")
	}
	if !strings.Contains(err.Error(), "frozen sections changed") {
		t.Fatalf("the refusal should say why: %v", err)
	}
}

// TestTheDocDescribesWhatTheCodeDoes keeps the spec and the code together.
//
// The doc is the specification: a gate that disagrees with it is a bug, and so is a doc
// that describes a gate the code does not have.
func TestTheDocDescribesWhatTheCodeDoes(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "docs", "03-domain-model.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc := string(raw)
	for _, want := range []string{
		"signed-off",             // the sign-off is its own state
		"or a human sign-off",    // the exception to the zero-match rule
		"deferred",               // a deferred criterion does not wedge the recipe
		"only the chef approves", // the permission deny-list
		"immutability.respected", // the gate that is now wired in
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("docs/03 should mention %q", want)
		}
	}
}
