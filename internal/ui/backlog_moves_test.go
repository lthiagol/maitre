package ui

import (
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// TestTheFooterNamesTheMoveTheDomainGives covers the footer on every state of an
// item.
//
// The console asks the domain what the chef can do with the selected item, and the
// footer reads that answer: `space shape it` on an open idea, `space promote to
// recipe` on an open small task, and **no** `space` on a settled item.
//
// The settled case is the one that holds the guard. A promoted item is still a
// small task — the kind says so — so a console that decided by the kind would offer
// the promotion, and the domain refuses it. The console asks instead.
func TestTheFooterNamesTheMoveTheDomainGives(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog

	// The first idea.
	m.selBacklog = 0
	out := plain(m)
	if !strings.Contains(out, "space shape it") {
		t.Errorf("the footer should offer the idea's move:\n%s", out)
	}
	if strings.Contains(out, "space promote to recipe") {
		t.Errorf("the footer offers a promotion the domain refuses:\n%s", out)
	}

	// The small task.
	m.selBacklog = 4
	out = plain(m)
	if !strings.Contains(out, "space promote to recipe") {
		t.Errorf("the footer does not offer the small task's promotion:\n%s", out)
	}

	// The promoted item: history, and no move at all.
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the promotion was refused: %s", m.refusal)
	}
	m.Update(m.reload()())
	out = plain(m)
	if strings.Contains(out, "space promote to recipe") || strings.Contains(out, "space shape it") {
		t.Errorf("a settled item has no move — the footer must not offer one:\n%s", out)
	}
}

// TestTheFooterAsksTheDomainForTheMove holds the console's guard.
//
// A console that decided for itself would pass on the two cases where its own rule
// and the domain's answer **agree**: an idea and an open small task. Those two hold
// nothing. The case where they **differ** is a promoted small task — still a small
// task by kind, no move by the domain's answer — and that is the case this holds.
func TestTheFooterAsksTheDomainForTheMove(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 4
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the promotion was refused: %s", m.refusal)
	}
	m.Update(m.reload()())

	// The domain's answer, read directly: no move.
	item, err := m.store.LoadBacklogItem("BL05")
	if err != nil {
		t.Fatal(err)
	}
	if move := core.BacklogMoveOf(*item); move.Primary != core.BacklogMoveNone {
		t.Fatalf("a promoted item's move is %q, want none", move.Primary)
	}
	// And the footer reads it.
	if out := plain(m); strings.Contains(out, "space promote to recipe") {
		t.Errorf("the console decided for itself — the domain says there is no move\n%s", out)
	}
}

// TestShapingAnIdeaChangesItsKindAndScopeAndKeepsItsId covers the shaping.
//
// The kind changes and a scope is added — and the id does not: `BL01` stays
// `BL01`, because the id is kind-neutral and the trail survives the shaping.
func TestShapingAnIdeaChangesItsKindAndScopeAndKeepsItsId(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	before, err := m.ops.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}

	steps := append([]string{"space", "f"}, typed("the evidence lane")...)
	steps = append(steps, "enter")
	if _, err := Drive(m, steps); err != nil {
		t.Fatal(err)
	}
	if m.refusal != "" {
		t.Fatalf("the shaping was refused: %s", m.refusal)
	}

	item, err := m.store.LoadBacklogItem("BL01")
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != "BL01" {
		t.Errorf("the id is %q, want BL01 — the id is kind-neutral", item.ID)
	}
	if item.Kind != core.KindFeature {
		t.Errorf("the kind is %q, want %q", item.Kind, core.KindFeature)
	}
	if len(item.Scope) != 1 || item.Scope[0] != "the evidence lane" {
		t.Errorf("the scope is %v, want the one the chef gave", item.Scope)
	}
	if item.Status != core.StatusOpen {
		t.Errorf("the status is %q, want %q — a shaping is not a promotion", item.Status, core.StatusOpen)
	}

	// No recipe was created: this is an update, not a promotion.
	after, err := m.ops.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Fatalf("the plan holds %d recipes, want %d — a shaping creates no recipe", len(after), len(before))
	}
}

// TestTheShapePromptAsksForTheKindAndTheScope covers the prompt.
//
// The console asks for what the domain needs and no more: what kind of work this
// is, and what is in it. A small task without a scope is work nobody can tell the
// end of.
func TestTheShapePromptAsksForTheKindAndTheScope(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	out, err := Drive(m, []string{"space"})
	if err != nil {
		t.Fatal(err)
	}
	plainOut := StripANSI(out)
	for _, want := range []string{"f) feature", "b) bug", "esc cancel"} {
		if !strings.Contains(plainOut, want) {
			t.Errorf("the prompt should offer %q:\n%s", want, plainOut)
		}
	}

	// And the scope is asked for in the next step.
	out, err = Drive(m, []string{"f"})
	if err != nil {
		t.Fatal(err)
	}
	plainOut = StripANSI(out)
	if !strings.Contains(plainOut, "scope for BL01") {
		t.Errorf("the prompt should ask for the scope:\n%s", plainOut)
	}
}

// TestAShapingWithoutAScopeIsNotWritten covers the empty answer.
//
// A small task without a scope is not a shaped item, so nothing is written and the
// prompt stays open: the chef types again rather than losing what they wrote.
func TestAShapingWithoutAScopeIsNotWritten(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	if _, err := Drive(m, []string{"space", "f", "enter"}); err != nil {
		t.Fatal(err)
	}
	if m.prompt == nil {
		t.Fatal("the prompt should stay open")
	}
	item, err := m.store.LoadBacklogItem("BL01")
	if err != nil {
		t.Fatal(err)
	}
	if item.Kind != core.KindIdea {
		t.Errorf("the kind is %q, want %q — nothing was written", item.Kind, core.KindIdea)
	}
}
