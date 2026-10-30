package ui

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
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
// The console asks the domain **whether** the item may be promoted and shows the
// domain's reason — it does not decide by itself that an idea is not promotable.
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

// TestTheFooterAsksTheDomainAboutPromotability holds the footer's guard.
//
// The console **asks** the domain whether an item may be promoted. A console that
// decided for itself would pass a test on the two cases where its own rule and the
// domain's answer **agree**: an idea (neither promotes it) and an open small task
// (both promote it). Those two hold nothing.
//
// The case where the answers **differ** is a **promoted** small task: it is still a
// small task — the console's own rule would say yes — and the domain refuses it,
// because it is no longer open. A console that decided by the kind would offer
// `space promote to recipe` on it; one that asks does not.
func TestTheFooterAsksTheDomainAboutPromotability(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	// The small task: the ideas come first, and an idea is not promoted.
	m.selBacklog = 4
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the promotion was refused: %s", m.refusal)
	}
	m.Update(m.reload()())

	// The item is a promoted small task now. The domain refuses it; the kind does
	// not.
	out := plain(m)
	if strings.Contains(out, "space promote to recipe") {
		t.Errorf("the footer offers a promotion the domain refuses — the console decided for itself\n%s", out)
	}
}

// TestSpaceOnAnIdeaWritesNothingToThePlan covers the journal.
//
// The footer offers no move on an idea, so pressing space is not an attempt: the
// plan's only event log stays byte-for-byte as it was. The journal's own rule is
// that a refusal is *"an attempt that was turned away"* — and nobody attempted
// anything.
func TestSpaceOnAnIdeaWritesNothingToThePlan(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	before := activityChecksum(t, m.store.Root)
	m.Update(tea.KeyPressMsg{Code: ' '})
	after := activityChecksum(t, m.store.Root)

	if before != after {
		t.Errorf("the refusal wrote to the journal: %s → %s", before, after)
	}

	// And the reason still reaches the chef: in the banner, and in the item's
	// own detail.
	if !strings.Contains(m.refusal, "acceptance criteria") {
		t.Errorf("the banner does not carry the reason: %q", m.refusal)
	}
	m.width = 200
	if out := plain(m); !strings.Contains(out, "acceptance criteria") {
		t.Errorf("the item's detail does not carry the reason\n%s", out)
	}
}

// activityChecksum is the checksum of the plan's event log.
func activityChecksum(t *testing.T, root string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, core.FileActivity))
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
