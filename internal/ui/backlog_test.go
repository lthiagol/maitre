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

// TestSpaceOnAnIdeaShapesIt covers the idea's move.
//
// T3a refused the space on an idea and named the remedy: *shape it into a small
// task first*. This slice **is** that remedy — so the space on an idea shapes it,
// and there is nothing to refuse. The domain's refusal of the *promotion* of an
// idea stands: `maitre-cli backlog promote` on an idea is still refused, with the
// reason.
func TestSpaceOnAnIdeaShapesIt(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the space on an idea was refused: %s", m.refusal)
	}
	if m.prompt == nil {
		t.Fatal("the space on an idea should open the shaping prompt")
	}
	if !strings.Contains(m.prompt.title(), "BL01") {
		t.Errorf("the prompt should name the item: %q", m.prompt.title())
	}
}

// TestThePromptWritesNothingUntilTheChefAnswers covers the journal.
//
// The journal's own rule is that a refusal is *"an attempt that was turned away"* —
// and an attempt is a write. Nothing is written until the chef answers the prompt,
// so the plan's only event log stays byte-for-byte as it was.
func TestThePromptWritesNothingUntilTheChefAnswers(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	before := activityChecksum(t, m.store.Root)
	m.Update(tea.KeyPressMsg{Code: ' '})
	after := activityChecksum(t, m.store.Root)

	if before != after {
		t.Errorf("the prompt wrote to the journal: %s → %s", before, after)
	}

	// And the item is still an idea.
	item, err := m.store.LoadBacklogItem("BL01")
	if err != nil {
		t.Fatal(err)
	}
	if item.Kind != core.KindIdea {
		t.Errorf("the kind is %q, want %q", item.Kind, core.KindIdea)
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

// TestTheItemDetailSaysWhyAnIdeaIsNotPromotable covers the item's own detail.
//
// The idea's move is a shaping, and the detail says why the item is not a recipe
// yet: it carries no acceptance criteria, and a recipe without them can never be
// signed off. The reason reaches the chef where the item is, not only in a banner.
func TestTheItemDetailSaysWhyAnIdeaIsNotPromotable(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0
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
