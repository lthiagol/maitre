package core

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// rawRecipe reads a recipe file without the store's normalisation.
//
// The store normalises and re-stamps on load, so the intermediate build's file can
// only be reproduced by reading the bytes directly.
func rawRecipe(t *testing.T, store *Store, id string) *Recipe {
	t.Helper()
	path, err := store.findRecipePath(id)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var r Recipe
	if err := json.Unmarshal(raw, &r); err != nil {
		t.Fatal(err)
	}
	return &r
}

// TestANormalisedLegacyApprovalIsRestamped covers the older build's digest.
//
// An intermediate build loaded the plan and re-saved it: the normalisation rewrote
// the lists to `[]` and left the approval's older digest in place. From then on the
// file's own form *is* the canonical one, so the store must recognise the older
// digest among the candidate encodings — not by recomputing the file's current form.
func TestANormalisedLegacyApprovalIsRestamped(t *testing.T) {
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

	// The older encoding: the file's lists were null when the recipe was signed off.
	raw := rawRecipe(t, store, r.ID)
	raw.Scope.InScope = nil
	legacy, err := raw.fingerprint()
	if err != nil {
		t.Fatal(err)
	}
	raw.Approval.Fingerprint = legacy
	if err := store.SaveRecipe(raw); err != nil {
		t.Fatal(err)
	}

	// The intermediate build: it loaded the file and re-saved it, so the lists are
	// normalised and the older digest is still in place.
	intermediate := rawRecipe(t, store, r.ID)
	intermediate.Scope.InScope = []string{}
	if err := store.SaveRecipe(intermediate); err != nil {
		t.Fatal(err)
	}
	onDisk := rawRecipe(t, store, r.ID)
	if onDisk.Approval.Fingerprint != legacy {
		t.Fatal("the test needs the older digest still in place")
	}
	if len(onDisk.Scope.InScope) != 0 {
		t.Fatal("the test needs the lists normalised")
	}

	// The store re-stamps it.
	loaded, err := store.LoadRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := loaded.Fingerprint()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Approval.Fingerprint == legacy {
		t.Fatal("the approval should be re-stamped canonically")
	}
	if loaded.Approval.Fingerprint != canonical {
		t.Fatalf("the approval carries %q, want the canonical %q", loaded.Approval.Fingerprint, canonical)
	}

	// And the gate does not fire on a recipe nobody touched.
	if _, err := ops.ClaimRecipe(r.ID, "line-cook"); err != nil {
		t.Fatalf("the claim should not be refused: %v", err)
	}
}

// TestAFormChangeIsNotATamper covers the same thing from the other side.
//
// `null` versus `[]` is the same empty list. A recipe edited only in that way must
// not be refused.
func TestAFormChangeIsNotATamper(t *testing.T) {
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

	// The older digest, over the file's own null lists.
	raw := rawRecipe(t, store, r.ID)
	raw.Scope.InScope = nil
	legacy, err := raw.fingerprint()
	if err != nil {
		t.Fatal(err)
	}
	raw.Approval.Fingerprint = legacy
	if err := store.SaveRecipe(raw); err != nil {
		t.Fatal(err)
	}

	// A hand edit: the same content, the other encoding of the same empty list.
	edited := rawRecipe(t, store, r.ID)
	edited.Scope.InScope = []string{}
	if err := store.SaveRecipe(edited); err != nil {
		t.Fatal(err)
	}

	if _, err := ops.ClaimRecipe(r.ID, "line-cook"); err != nil {
		t.Fatalf("the claim should not be refused: %v", err)
	}
}

// TestAnAgentMayNotSignOff covers the sign-off's deny-list.
//
// A sign-off is the one exception to the zero-match rule, so it is a human act and
// takes the same deny-list the approval takes. No lie is needed for the hole: the
// actor is recorded truthfully.
func TestAnAgentMayNotSignOff(t *testing.T) {
	store := Open(t.TempDir() + "/" + PlanDirName)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	line := NewOperations(store, AuditContext{Actor: "line-cook", Session: "s-1", Role: "runner"})
	r, err := line.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}

	_, err = line.SignOffCriterion(r.ID, "AC-01", "line-cook", "checked")
	if err == nil {
		t.Fatal("an agent must not be able to sign off")
	}
	if !strings.Contains(err.Error(), "line-cook") {
		t.Fatalf("the refusal should name the actor: %v", err)
	}
	if !strings.Contains(err.Error(), "chef") {
		t.Fatalf("the refusal should name the rule: %v", err)
	}

	// The chef can.
	chef := NewOperations(store, AuditContext{Actor: "chef", Session: "s-2", Role: "coordinator"})
	if _, err := chef.SignOffCriterion(r.ID, "AC-01", "chef", "checked"); err != nil {
		t.Fatalf("the chef should be able to sign off: %v", err)
	}
}

// TestASignedOffCriterionIsPassedButNotProven covers the sign-off's own state.
//
// The chef may sign off under any name, and the sign-off claims no match.
func TestASignedOffCriterionIsPassedButNotProven(t *testing.T) {
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
	if _, err := ops.ClaimRecipe(r.ID, "line-cook"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.SetStepStatus(r.ID, "G1", "done"); err != nil {
		t.Fatal(err)
	}

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
	proof := c.Proofs[0]
	if proof.Match != MatchSignedOff {
		t.Fatalf("match = %q, want %q — a sign-off is not a machine proof", proof.Match, MatchSignedOff)
	}
	if proof.Match == MatchProven {
		t.Fatal("a sign-off must never claim a machine proof")
	}
	if proof.Matched != 0 {
		t.Fatalf("matched = %d, want 0", proof.Matched)
	}

	// And the recipe finishes: the sign-off is the one exception to the zero-match
	// rule, and it still works.
	if _, err := ops.Advance(r.ID, EventFinish); err != nil {
		t.Fatalf("finish: %v", err)
	}
}
