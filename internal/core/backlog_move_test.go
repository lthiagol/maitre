package core

import "testing"

// TestTheBacklogMoveIsOneAnswer covers the answer the console asks.
//
// The console's footer names the move the selected item has, and its `space` key
// makes it. The answer lives here, once: a console-side rule ("is it a small
// task?") is a second answer to one question, and two answers drift — the defect
// F-114 closed.
func TestTheBacklogMoveIsOneAnswer(t *testing.T) {
	ops, store := newOps(t)

	idea, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	task, err := ops.CreateBacklogItem("A small task", "worth doing", []string{"the work"}, KindFeature, nil)
	if err != nil {
		t.Fatal(err)
	}

	// An open idea: one move, and it is not a promotion.
	move := BacklogMoveOf(*idea)
	if move.Primary != BacklogMoveShape {
		t.Errorf("an idea's move is %q, want %q", move.Primary, BacklogMoveShape)
	}
	if !move.Dismissable {
		t.Error("an open item may be dismissed")
	}

	// An open small task: the promotion.
	move = BacklogMoveOf(*task)
	if move.Primary != BacklogMovePromote {
		t.Errorf("a small task's move is %q, want %q", move.Primary, BacklogMovePromote)
	}

	// A promoted item is history: no move, and nothing to dismiss.
	if _, err := ops.PromoteBacklogItem(task.ID); err != nil {
		t.Fatal(err)
	}
	settled, err := store.LoadBacklogItem(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	move = BacklogMoveOf(*settled)
	if move.Primary != BacklogMoveNone {
		t.Errorf("a promoted item's move is %q, want none", move.Primary)
	}
	if move.Dismissable {
		t.Error("a promoted item may not be dismissed")
	}
}

// TestTheMoveNeverOffersARefusedPromotion holds the two answers together.
//
// The move and the promotion are one fact with two faces: an open small task has the
// promotion, and a settled item has neither. A move that offered a promotion the
// domain refuses would put the chef in a dead end — the trap T3a closed.
func TestTheMoveNeverOffersARefusedPromotion(t *testing.T) {
	ops, _ := newOps(t)

	idea, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	task, err := ops.CreateBacklogItem("A small task", "worth doing", []string{"the work"}, KindBug, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range []BacklogItem{*idea, *task} {
		move := BacklogMoveOf(item)
		allowed := PromotionOf(item).Allowed
		if move.Primary == BacklogMovePromote && !allowed {
			t.Errorf("%s offers the promotion, and the domain refuses it", item.ID)
		}
		if move.Primary != BacklogMovePromote && allowed {
			t.Errorf("%s may be promoted, but the move is %q", item.ID, move.Primary)
		}
	}
}

// TestTheMoveCarriesTheKindsAShapingMayPick covers the kinds.
//
// The console's shape prompt offers the kinds the **domain** names, so the list is
// read from the rule — `isSmallTask` — and not copied: a kind the rule gains is a
// kind the console's prompt offers the same day.
func TestTheMoveCarriesTheKindsAShapingMayPick(t *testing.T) {
	ops, _ := newOps(t)

	idea, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	move := BacklogMoveOf(*idea)
	if move.Primary != BacklogMoveShape {
		t.Fatalf("the move is %q, want %q", move.Primary, BacklogMoveShape)
	}
	if len(move.Kinds) == 0 {
		t.Fatal("a shaping names no kind to shape into")
	}
	for _, kind := range move.Kinds {
		if !isSmallTask(kind) {
			t.Errorf("%q is offered, and it is not a small task", kind)
		}
	}
	// Every small task the domain knows is offered.
	for _, kind := range BacklogKinds {
		if !isSmallTask(kind) {
			continue
		}
		if !containsString(move.Kinds, kind) {
			t.Errorf("%q is a small task, and the move does not offer it", kind)
		}
	}

	// And a settled item offers nothing: there is no shaping to pick a kind for.
	task, err := ops.CreateBacklogItem("A small task", "worth doing", nil, KindFeature, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.PromoteBacklogItem(task.ID); err != nil {
		t.Fatal(err)
	}
	settled, err := ops.ListBacklogItems()
	if err != nil {
		t.Fatal(err)
	}
	for _, it := range settled {
		if it.ID != task.ID {
			continue
		}
		if got := BacklogMoveOf(it); len(got.Kinds) != 0 {
			t.Errorf("a settled item offers %v, want none", got.Kinds)
		}
	}
}

// TestTheMoveIsALabeledSet covers the closed set.
//
// The footer calls the moves by name, and a move the domain gains is a name the
// footer knows the same day.
func TestTheMoveIsALabeledSet(t *testing.T) {
	names := map[string]bool{}
	for _, mv := range AllBacklogMoves {
		names[string(mv)] = true
	}
	for _, want := range []string{"", "shape", "promote"} {
		if !names[want] {
			t.Errorf("the closed set is missing %q", want)
		}
	}
	if BacklogMoveShape.Label() != "shape it" {
		t.Errorf("the shape move is called %q", BacklogMoveShape.Label())
	}
	if BacklogMovePromote.Label() != "promote to recipe" {
		t.Errorf("the promote move is called %q", BacklogMovePromote.Label())
	}
	if BacklogMoveNone.Label() != "" {
		t.Errorf("no move has no name, got %q", BacklogMoveNone.Label())
	}
}
