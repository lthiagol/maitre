package ui

import (
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// TestABacklogItemHasADetail covers the item's own pane.
//
// The item is not a recipe: it has its own fields — the body, the scope, the tags,
// the dates — and the console's one detail pane renders them.
func TestABacklogItemHasADetail(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0
	// The split layout: the detail sits beside the list.
	m.width = 200

	out := plain(m)
	for _, want := range []string{
		"BL01",                                   // the id
		"The board needs a way to show evidence", // the title
		"Why",                                    // the body
		"Evidence the chef cannot read is a claim", // the body's own words
		"idea",    // the kind
		"open",    // the status
		"board",   // the tags
		"created", // the dates
		"updated",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the detail is missing %q\n%s", want, out)
		}
	}
}

// TestEnterOpensTheBacklogItem covers the key.
//
// The detail pane is one pane: the same geometry, the same scroll, the same `esc`.
func TestEnterOpensTheBacklogItem(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.mode != ModeDetail {
		t.Fatalf("enter did not open the detail: mode is %v", m.mode)
	}
	out := plain(m)
	if !strings.Contains(out, "The board needs a way to show evidence") {
		t.Errorf("the open detail does not show the item\n%s", out)
	}

	// And `esc` brings it back.
	m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if m.mode != ModeSplit {
		t.Fatalf("esc did not bring the split view back: mode is %v", m.mode)
	}
}

// TestTheBacklogDetailShowsWhereTheItemWent covers the promoted item.
//
// The trail survives: the item says where it went.
func TestTheBacklogDetailShowsWhereTheItemWent(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	// The small task: the ideas come first, and an idea is not promoted.
	m.selBacklog = 4
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the promotion was refused: %s", m.refusal)
	}
	m.Update(m.reload()())

	m.width = 200
	out := plain(m)
	if !strings.Contains(out, "Promoted to recipe") {
		t.Errorf("the detail does not say where the item went\n%s", out)
	}
}

// TestTheFooterDoesNotOfferAMoveTheDomainRefuses covers the footer.
//
// A footer that promises a promotion the domain refuses is a lie: on an idea there
// is no such move, so the console must not say there is one. The small task still
// has its move.
func TestTheFooterDoesNotOfferAMoveTheDomainRefuses(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog

	// The first idea.
	m.selBacklog = 0
	out := plain(m)
	if strings.Contains(out, "space promote to recipe") {
		t.Errorf("the footer promises a promotion the domain refuses\n%s", out)
	}

	// The small task.
	m.selBacklog = 4
	out = plain(m)
	if !strings.Contains(out, "space promote to recipe") {
		t.Errorf("the footer does not offer the small task's promotion\n%s", out)
	}
}

// TestSpaceOnAnIdeaIsRefusedWithTheReason covers the refusal.
//
// The console asks the domain to promote the item and shows the domain's reason —
// it does not decide by itself that an idea is not promotable.
func TestSpaceOnAnIdeaIsRefusedWithTheReason(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal == "" {
		t.Fatal("the space on an idea was not refused")
	}
	for _, want := range []string{"acceptance criteria", "signed off", "small task"} {
		if !strings.Contains(m.refusal, want) {
			t.Errorf("the refusal should name %q: %s", want, m.refusal)
		}
	}

	// And nothing was written.
	items, err := m.ops.ListBacklogItems()
	if err != nil {
		t.Fatal(err)
	}
	for _, it := range items {
		if it.Status != core.StatusOpen {
			t.Fatalf("%s is %q — a refused promotion writes nothing", it.ID, it.Status)
		}
	}
	recipes, err := m.ops.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if len(recipes) != len(m.recipes) {
		t.Fatalf("the plan holds %d recipes, want %d", len(recipes), len(m.recipes))
	}
}
