package core

import (
	"encoding/json"
	"os"
	"path/filepath"
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

// TestARenamedFieldIsNotATamper covers the rename.
//
// The digest hashes a **canonical payload** — a shape the code owns, with fixed field
// names — so the name of a record's field is not part of the digest. A rename is not
// an edit: the content did not change, only the encoding. A digest that cannot tell
// the two apart reads the rename as tampering and refuses every gated transition of a
// recipe nobody touched.
//
// So the payload is the code's, and the recipe's own field names are not in it. The
// wire is the golden digest below: it pins the payload's shape, so a payload that
// starts carrying a record's name fails the test.
func TestARenamedFieldIsNotATamper(t *testing.T) {
	r := Recipe{
		ID:    "RCP01",
		Title: "A recipe",
		Intent: Intent{
			Outcome: "the board renders",
			Problem: "the board does not render",
		},
		Scope: Scope{
			InScope:    []string{"the board"},
			OutOfScope: []string{"the web shell"},
		},
		Acceptance: []Criterion{{
			ID:           "AC-01",
			Description:  "the board renders",
			Verification: "go test ./internal/core -run Board",
		}},
	}
	got, err := r.Fingerprint()
	if err != nil {
		t.Fatal(err)
	}
	if got != canonicalDigestGolden {
		t.Fatalf("the digest is %q, want %q — the payload's shape moved", got, canonicalDigestGolden)
	}

	// The payload names the frozen sections and nothing else: a criterion's records
	// are not in it, so they cannot move the digest when they grow.
	raw, err := json.Marshal(r.approvalPayload())
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"evidence", "proofs", "deferral", "title", "guidance"} {
		if strings.Contains(string(raw), `"`+name+`"`) {
			t.Errorf("the payload names %q — it must cover the frozen sections only:\n%s", name, raw)
		}
	}

	// The records grow after the approval, and the digest does not move.
	grown := r
	grown.Acceptance = []Criterion{{
		ID:           "AC-01",
		Description:  "the board renders",
		Verification: "go test ./internal/core -run Board",
		Evidence: []Evidence{{
			Command: "go test ./internal/core -run Board", ExitCode: 0, Match: MatchProven,
		}},
	}}
	after, err := grown.Fingerprint()
	if err != nil {
		t.Fatal(err)
	}
	if after != got {
		t.Fatalf("the digest moved with the records: %q vs %q", after, got)
	}
}

// canonicalDigestGolden is the digest of the recipe in
// `TestARenamedFieldIsNotATamper`.
//
// It pins the canonical payload: the field names it uses and the shape it hashes. A
// payload that starts carrying a record's name moves this digest, and the test says
// so — which is how the root cause stays fixed.
const canonicalDigestGolden = "49701528b23aec009a608b4803738247bfa43ef95c41d16a4629d1aae33b6d6f"

// TestTheRepoPlansApprovalDigestMatchesItsContent holds the standing discipline.
//
// The repository's own plan is the one artifact in the tree that a rename can break
// silently: its stored digest was computed under the old field name, and the
// immutability gate reads the difference as tampering. Nothing held the digest to the
// canonical recomputation, which is why that could ship.
//
// So it is held now: the stored digest must match what the tool recomputes.
func TestTheRepoPlansApprovalDigestMatchesItsContent(t *testing.T) {
	root := filepath.Join("..", "..", PlanDirName)
	if _, err := os.Stat(root); err != nil {
		t.Skipf("the repo has no plan at %s", root)
	}
	store := Open(root)
	recipes, err := store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range recipes {
		if r.Approval == nil {
			continue
		}
		// The file's own digest, not the store's: the question is what the plan
		// carries on disk.
		stored := rawRecipe(t, store, r.ID)
		want, err := stored.Fingerprint()
		if err != nil {
			t.Fatal(err)
		}
		if stored.Approval.Fingerprint != want {
			t.Errorf("%s carries the digest %q, want the canonical %q — the approval was written under another encoding",
				r.ID, stored.Approval.Fingerprint, want)
		}
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

	// A hand edit: the same content, the other encoding of the same empty list.
	edited := rawRecipe(t, store, r.ID)
	edited.Scope.InScope = nil
	if err := store.SaveRecipe(edited); err != nil {
		t.Fatal(err)
	}

	if _, err := ops.ClaimRecipe(r.ID, "runner"); err != nil {
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
	line := NewOperations(store, AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})
	r, err := line.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}

	_, err = line.SignOffCriterion(r.ID, "AC-01", "runner", "checked")
	if err == nil {
		t.Fatal("an agent must not be able to sign off")
	}
	if !strings.Contains(err.Error(), "runner") {
		t.Fatalf("the refusal should name the actor: %v", err)
	}
	if !strings.Contains(err.Error(), "chef") {
		t.Fatalf("the refusal should name the rule: %v", err)
	}

	// The chef can.
	chef := NewOperations(store, AuditContext{Actor: "chef", Session: "s-2", Role: "chef"})
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
	if _, err := ops.ClaimRecipe(r.ID, "runner"); err != nil {
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
	evidence := c.Evidence[0]
	if evidence.Match != MatchSignedOff {
		t.Fatalf("match = %q, want %q — a sign-off is not machine evidence", evidence.Match, MatchSignedOff)
	}
	if evidence.Match == MatchProven {
		t.Fatal("a sign-off must never claim machine evidence")
	}
	if evidence.Matched != 0 {
		t.Fatalf("matched = %d, want 0", evidence.Matched)
	}

	// And the recipe finishes: the sign-off is the one exception to the zero-match
	// rule, and it still works.
	if _, err := ops.Advance(r.ID, EventFinish); err != nil {
		t.Fatalf("finish: %v", err)
	}
}
