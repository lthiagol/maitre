package core

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestARefusedTransitionIsJournalled covers the refusal the journal used to swallow.
//
// The journal is the only event log, and the L5 audit trail is read from it. A refusal
// is the moment an attempt was turned away — it is the evidence the chef wants, and
// it was written to nowhere.
func TestARefusedTransitionIsJournalled(t *testing.T) {
	ops, store := newOps(t)
	bad := spec()
	bad.OutOfScope = nil // too few out-of-scope items, so groom is refused
	r, err := ops.CreateRecipe(bad)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ops.Advance(r.ID, EventGroom); err == nil {
		t.Fatal("expected groom to be refused")
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	last := j.Events[len(j.Events)-1]
	if last.Kind != KindRefusal {
		t.Fatalf("kind = %q, want %q — a refused transition must be journalled", last.Kind, KindRefusal)
	}
	if last.Recipe != r.ID {
		t.Fatalf("recipe = %q, want %q", last.Recipe, r.ID)
	}
	if last.Actor != "chef" || last.Session != "s-1" {
		t.Fatalf("the refusal is not stamped with the audit context: %+v", last)
	}

	data := refusalData(t, last)
	if data["event"] != string(EventGroom) {
		t.Fatalf("the refusal does not name the event: %v", data)
	}
	gates, _ := data["gates"].([]any)
	if len(gates) == 0 {
		t.Fatalf("the refusal does not name the gates: %v", data)
	}
	found := false
	for _, g := range gates {
		if g == "scope.out_of_scope_declared" {
			found = true
		}
	}
	if !found {
		t.Fatalf("the refusal does not name the gate that refused it: %v", gates)
	}
	if data["reason"] == nil || data["reason"] == "" {
		t.Fatalf("the refusal does not say why: %v", data)
	}
}

// TestThePlanAndTheJournalAgree covers the refusal's other half.
//
// The refusal writes an event and nothing else: the recipe did not move, so the plan
// and the journal must tell the same story.
func TestThePlanAndTheJournalAgree(t *testing.T) {
	ops, store := newOps(t)
	bad := spec()
	bad.OutOfScope = nil
	r, err := ops.CreateRecipe(bad)
	if err != nil {
		t.Fatal(err)
	}
	before, err := store.LoadRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ops.Advance(r.ID, EventGroom); err == nil {
		t.Fatal("expected groom to be refused")
	}

	after, err := store.LoadRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Lifecycle != before.Lifecycle {
		t.Fatalf("lifecycle = %q, want %q — a refused transition changes nothing",
			after.Lifecycle, before.Lifecycle)
	}
	if after.Approval != nil {
		t.Fatal("a refused transition must not sign the recipe off")
	}
}

// TestARefusedClaimIsJournalled covers the claim path.
//
// The claim is a transition, so a refused claim is a refused transition — and the
// dependency rule is the gate most often turned away.
func TestARefusedClaimIsJournalled(t *testing.T) {
	ops, store := newOps(t)
	dep := approvedRecipe(t, ops, "The dependency")
	r := approvedRecipe(t, ops, "The dependent")
	if _, err := ops.UpdateRecipe(r.ID, Fragment{DependsOn: &[]string{dep.ID}}); err != nil {
		t.Fatal(err)
	}

	if _, err := ops.ClaimRecipe(r.ID, "runner"); err == nil {
		t.Fatal("a recipe with an open dependency must not be claimed")
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	last := j.Events[len(j.Events)-1]
	if last.Kind != KindRefusal {
		t.Fatalf("kind = %q, want %q", last.Kind, KindRefusal)
	}
	data := refusalData(t, last)
	if data["event"] != string(EventStart) {
		t.Fatalf("the refusal does not name the event: %v", data)
	}
}

// TestABadRequestIsNotJournalled covers the other side of the judgement.
//
// "no such recipe" is a mistake in the request, not a decision about the plan. The
// journal is the plan's event log; filling it with typos would bury the refusals that
// matter.
func TestABadRequestIsNotJournalled(t *testing.T) {
	ops, store := newOps(t)
	before, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ops.Advance("R99", EventGroom); err == nil {
		t.Fatal("expected advance on a missing recipe to fail")
	}

	after, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Events) != len(before.Events) {
		t.Fatalf("the journal has %d events, want %d — a bad request writes no event",
			len(after.Events), len(before.Events))
	}
}

// TestARefusedOverlayIsJournalled covers the refusals the state of the artifact
// produces.
//
// A special that is already promoted cannot be promoted again. The console says
// "refused", so the journal has to say it too: the banner and the log must not use
// the word for two different sets.
func TestARefusedOverlayIsJournalled(t *testing.T) {
	ops, store := newOps(t)
	sp, err := ops.CreateSpecial("A hunch", "worth noticing", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.PromoteSpecial(sp.ID); err != nil {
		t.Fatal(err)
	}
	before, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ops.PromoteSpecial(sp.ID); err == nil {
		t.Fatal("promoting a promoted special should be refused")
	}
	after, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Events) != len(before.Events)+1 {
		t.Fatalf("the journal grew by %d events, want 1", len(after.Events)-len(before.Events))
	}
	last := after.Events[len(after.Events)-1]
	if last.Kind != KindRefusal {
		t.Fatalf("the event's kind = %q, want %q", last.Kind, KindRefusal)
	}
}

// TestARefusedEditIsJournalled covers the refusal outside a transition.
//
// The frozen sections are the plan refusing an edit, not a transition. It is the same
// kind of refusal, so it is the same kind of event.
func TestARefusedEditIsJournalled(t *testing.T) {
	ops, store := newOps(t)
	r := approvedRecipe(t, ops, "A recipe")

	outcome := "something else"
	if _, err := ops.UpdateRecipe(r.ID, Fragment{Outcome: &outcome}); err == nil {
		t.Fatal("expected the intent to be frozen after sign-off")
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	last := j.Events[len(j.Events)-1]
	if last.Kind != KindRefusal {
		t.Fatalf("kind = %q, want %q — a refused edit must be journalled", last.Kind, KindRefusal)
	}
	if last.Recipe != r.ID {
		t.Fatalf("recipe = %q, want %q", last.Recipe, r.ID)
	}
	data := refusalData(t, last)
	if !strings.Contains(data["reason"].(string), "frozen") {
		t.Fatalf("the refusal does not say why: %v", data)
	}
}

// TestTheRefusalKeepsThePlanVersion covers the version stamp.
//
// A refusal writes no artifact, so the rules the plan was written under did not
// change. If the refusal stamped the plan with the tool's version, a newer tool
// refusing a write in an older plan would lock the older tool out of it.
func TestTheRefusalKeepsThePlanVersion(t *testing.T) {
	ops, store := newOps(t)
	bad := spec()
	bad.OutOfScope = nil
	r, err := ops.CreateRecipe(bad)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventGroom); err == nil {
		t.Fatal("expected groom to be refused")
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	last := j.Events[len(j.Events)-1]
	if last.Kind != KindRefusal {
		t.Fatalf("kind = %q, want %q", last.Kind, KindRefusal)
	}
	if last.DomainVersion != DomainVersion {
		t.Fatalf("domain version = %d, want %d", last.DomainVersion, DomainVersion)
	}
}

// refusalData reads a journal entry's data as a map.
func refusalData(t *testing.T, e JournalEntry) map[string]any {
	t.Helper()
	if len(e.Data) == 0 {
		t.Fatalf("the event carries no data: %+v", e)
	}
	var data map[string]any
	if err := json.Unmarshal(e.Data, &data); err != nil {
		t.Fatal(err)
	}
	return data
}

// TestTheRefusalNamesTheGateNotItsProse covers the gate list's shape.
//
// The gates are named, not paraphrased: a reader matching them against the gate table
// must find the same names.
func TestTheRefusalNamesTheGateNotItsProse(t *testing.T) {
	ops, store := newOps(t)
	bad := spec()
	bad.OutOfScope = nil
	r, err := ops.CreateRecipe(bad)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventGroom); err == nil {
		t.Fatal("expected groom to be refused")
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	data := refusalData(t, j.Events[len(j.Events)-1])
	for _, g := range data["gates"].([]any) {
		name, _ := g.(string)
		if strings.Contains(name, "(") {
			t.Fatalf("the gate name carries its prose: %q", name)
		}
		if !containsString(GateNames(), name) {
			t.Fatalf("the gate %q is not in the gate table", name)
		}
	}
}
