package core

import (
	"errors"
	"strings"
	"testing"
)

// TestAnIdeaIsNotPromoted covers the refusal.
//
// An idea is a hunch: it carries no acceptance criteria, and a recipe without them
// can never pass `ac.falsifiable`, so it can never be signed off. Promoting one
// would file the work as **known** — the domain's own definition of a small task —
// and would leave the chef in a dead end one keystroke away.
//
// The refusal names the remedy: shape it into a small task first.
func TestAnIdeaIsNotPromoted(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ops.PromoteBacklogItem(item.ID)
	if err == nil {
		t.Fatal("an idea must not be promoted")
	}
	var re *RefusalError
	if !errors.As(err, &re) {
		t.Fatalf("the refusal should be a refusal: %v", err)
	}
	for _, want := range []string{"acceptance criteria", "signed off", "small task"} {
		if !strings.Contains(re.Reason, want) {
			t.Errorf("the refusal should name %q: %q", want, re.Reason)
		}
	}

	// Nothing was written: no recipe, and the item is still open.
	recipes, err := store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if len(recipes) != 0 {
		t.Fatalf("the plan holds %d recipes, want none", len(recipes))
	}
	after, err := store.LoadBacklogItem(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != StatusOpen {
		t.Fatalf("the item is %q, want %q — a refusal writes nothing", after.Status, StatusOpen)
	}
}

// TestThePromotionAnswerIsOnePlace covers the answer both surfaces read.
//
// The domain refuses the promotion, and the console's footer reads the same answer —
// so the console cannot offer a move the domain refuses. Two implementations of the
// rule would be how the console and the domain drift apart.
func TestThePromotionAnswerIsOnePlace(t *testing.T) {
	ops, _ := newOps(t)

	idea, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	task, err := ops.CreateBacklogItem("A small task", "worth doing", []string{"the work"}, KindFeature, nil)
	if err != nil {
		t.Fatal(err)
	}

	if p := PromotionOf(*idea); p.Allowed {
		t.Error("an idea may not be promoted")
	}
	if p := PromotionOf(*task); !p.Allowed {
		t.Errorf("a small task may be promoted: %s", p.Reason)
	}

	// And the refusal the operation gives is the answer's reason.
	_, err = ops.PromoteBacklogItem(idea.ID)
	if err == nil {
		t.Fatal("an idea must not be promoted")
	}
	var re *RefusalError
	if !errors.As(err, &re) {
		t.Fatalf("the refusal should be a refusal: %v", err)
	}
	if re.Reason != PromotionOf(*idea).Reason {
		t.Errorf("the refusal says %q, want the answer's reason %q",
			re.Reason, PromotionOf(*idea).Reason)
	}
}

// TestAPromotedItemMayNotBePromotedAgain covers the closed set.
//
// The item is history once it is promoted: it says where it went, and it cannot go
// there twice.
func TestAPromotedItemMayNotBePromotedAgain(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A small task", "worth doing", nil, KindFeature, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.PromoteBacklogItem(item.ID); err != nil {
		t.Fatal(err)
	}
	after, err := store.LoadBacklogItem(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if p := PromotionOf(*after); p.Allowed {
		t.Error("a promoted item may not be promoted again")
	}
	if _, err := ops.PromoteBacklogItem(item.ID); err == nil {
		t.Fatal("a promoted item must not be promoted again")
	}
}

// TestASmallTaskPromotesToARecipeWithOneCriterion covers the small task's road.
//
// A small task is too small to be a big recipe, not too small to be guarded: it
// becomes a recipe with **one** criterion, so the work is falsifiable from the start.
func TestASmallTaskPromotesToARecipeWithOneCriterion(t *testing.T) {
	for _, kind := range []string{KindFeature, KindBug} {
		ops, _ := newOps(t)
		item, err := ops.CreateBacklogItem("Wire the queue", "Findings reach the chef", []string{"the queue lane"}, kind, nil)
		if err != nil {
			t.Fatal(err)
		}
		recipe, err := ops.PromoteBacklogItem(item.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(recipe.Acceptance) != 1 {
			t.Fatalf("a %s promoted to %d criteria, want exactly one", kind, len(recipe.Acceptance))
		}
		c := recipe.Acceptance[0]
		if c.Description != item.Title {
			t.Fatalf("the criterion is %q, want %q", c.Description, item.Title)
		}
		if recipe.Scope.InScope[0] != "the queue lane" {
			t.Fatalf("the scope did not carry over: %+v", recipe.Scope)
		}
	}
}

// TestAPromotionIsOneWay covers the trail.
//
// The item's status becomes `promoted` and it records what it became. Promotion
// cannot happen twice, and a promoted item is history: it cannot be edited either.
func TestAPromotionIsOneWay(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A small task", "worth doing", nil, KindFeature, nil)
	if err != nil {
		t.Fatal(err)
	}
	recipe, err := ops.PromoteBacklogItem(item.ID)
	if err != nil {
		t.Fatal(err)
	}

	after, err := store.LoadBacklogItem(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != StatusPromoted {
		t.Fatalf("the status is %q, want %q", after.Status, StatusPromoted)
	}
	if after.PromotedTo != "recipe "+recipe.ID {
		t.Fatalf("promoted_to is %q, want %q", after.PromotedTo, "recipe "+recipe.ID)
	}

	// A second promotion fails: one item, one recipe.
	if _, err := ops.PromoteBacklogItem(item.ID); err == nil {
		t.Fatal("a promoted item must not be promoted again")
	}
	// And it cannot be edited either.
	if _, err := ops.UpdateBacklogItem(item.ID, "A better hunch", "sharper", nil, KindFeature, nil); err == nil {
		t.Fatal("a promoted item must not be edited")
	}
	recipes, err := store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if len(recipes) != 1 {
		t.Fatalf("the plan holds %d recipes, want 1", len(recipes))
	}
}

// TestThePromotionWeakensNoGate covers the promise.
//
// The brief said it plainly: if a promotion needs a weaker gate or a new plan-format
// field, that is a design decision, not an implementation detail. The promoted
// recipe is **not** pre-groomed, and the gates say so: it has no out-of-scope items,
// its criterion has no verification, and it has no steps — so the gates that guard a
// recipe refuse it, exactly as they refuse any other half-written recipe.
func TestThePromotionWeakensNoGate(t *testing.T) {
	ops, _ := newOps(t)

	item, err := ops.CreateBacklogItem("Wire the queue", "Findings reach the chef", []string{"the queue lane"}, KindFeature, nil)
	if err != nil {
		t.Fatal(err)
	}
	recipe, err := ops.PromoteBacklogItem(item.ID)
	if err != nil {
		t.Fatal(err)
	}

	// The recipe carries no out-of-scope items and no steps: the groomer writes
	// them, and the gates hold.
	if len(recipe.Scope.OutOfScope) != 0 {
		t.Fatalf("the recipe carries %d out-of-scope items, want none", len(recipe.Scope.OutOfScope))
	}
	if len(recipe.Guidance) != 0 {
		t.Fatalf("the recipe carries %d steps, want none", len(recipe.Guidance))
	}
	if len(recipe.WorkPackages) != 0 {
		t.Fatalf("the recipe carries %d work packages, want none", len(recipe.WorkPackages))
	}

	// The criterion names the work; its verification is the recipe's own life.
	if recipe.Acceptance[0].Verification != "" {
		t.Fatalf("the criterion carries a verification %q — the recipe is not pre-groomed",
			recipe.Acceptance[0].Verification)
	}

	// So the groom gate refuses it — the same refusal any other half-written
	// recipe gets.
	_, err = ops.Advance(recipe.ID, EventGroom)
	if err == nil {
		t.Fatal("a promoted recipe must not be groomed before it is written")
	}
	var ge *GateError
	if !asGateError(err, &ge) {
		t.Fatalf("the refusal should be a gate refusal: %v", err)
	}
	joined := strings.Join(ge.Gates, " ")
	if !strings.Contains(joined, "scope.out_of_scope_declared") {
		t.Fatalf("the gates should name the missing out-of-scope items: %v", ge.Gates)
	}
	if !strings.Contains(joined, "ac.falsifiable") {
		t.Fatalf("the gates should name the criterion with no verification: %v", ge.Gates)
	}
}

// asGateError is errors.As, spelled out so the test does not import errors for
// one call.
func asGateError(err error, target **GateError) bool {
	ge, ok := err.(*GateError)
	if ok {
		*target = ge
	}
	return ok
}
