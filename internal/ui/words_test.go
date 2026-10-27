package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// containsWord reports whether the text holds the word as a whole word.
//
// The check is case-sensitive: the board's sections are capitalised ("Tasting") and
// the records are plain ("awaiting tasting"), and one is not the other.
func containsWord(text, word string) bool {
	from := 0
	for {
		i := strings.Index(text[from:], word)
		if i < 0 {
			return false
		}
		i += from
		before := i == 0 || !isWordByte(text[i-1])
		after := i+len(word) == len(text) || !isWordByte(text[i+len(word)])
		if before && after {
			return true
		}
		from = i + 1
	}
}

// isWordByte reports whether the byte is part of a word.
//
// The lane labels hold a dash ("Sign-off"), and a dash does not break a word: it is
// what the label is made of.
func isWordByte(b byte) bool {
	return b != ' ' && b != '\t' && b != '·' && b != '|' && b != ','
}

// TestTheRowShowsThePlainStateWord covers the row.
//
// A recipe's row shows the plain state word — `closed`, `executing` — and the lane is
// a *section* the board groups by. A section is never a second name for the row's
// state.
func TestTheRowShowsThePlainStateWord(t *testing.T) {
	m := newModel(t)
	for _, r := range m.recipes {
		s := core.SummaryOf(r)
		row := StripANSI(m.recipeRow(s, false, m.listColumn()))
		if !strings.Contains(row, string(s.Lifecycle)) {
			t.Errorf("%s: the row should show the state word %q:\n%s", s.ID, s.Lifecycle, row)
		}
		if lane := core.LaneOf(r).Label(); containsWord(row, lane) {
			t.Errorf("%s: the row should not show the lane word %q beside the state:\n%s", s.ID, lane, row)
		}
	}
}

// TestTheDetailShowsOneWordForTheState covers the detail pane.
//
// The pane held the state word and the lane's word side by side — "closed · Served"
// — and said the same thing twice in two vocabularies. One word per fact: the state
// word, and the lane nowhere.
func TestTheDetailShowsOneWordForTheState(t *testing.T) {
	m := newModel(t)
	for _, r := range m.recipes {
		lines := m.recipeDetail(&r, 100)
		if len(lines) < 2 {
			t.Fatalf("%s: the detail should have a meta line", r.ID)
		}
		meta := StripANSI(lines[1])
		if !strings.Contains(meta, string(r.Lifecycle)) {
			t.Errorf("%s: the meta line should show the state word %q:\n%s", r.ID, r.Lifecycle, meta)
		}
		if lane := core.LaneOf(r).Label(); containsWord(meta, lane) {
			t.Errorf("%s: the meta line should not show the lane word %q:\n%s", r.ID, lane, meta)
		}
	}
}

// TestABlockedRecipeStillSaysItIsOnHold covers what the lane word was carrying.
//
// The lane was the only place the detail pane said a recipe was blocked. Now the
// state word holds the lane's place, so the block needs its own line — with the
// reason, which is the fact the chef acts on.
func TestABlockedRecipeStillSaysItIsOnHold(t *testing.T) {
	root := filepath.Join(t.TempDir(), core.PlanDirName)
	store := core.Open(root)
	if err := core.Seed(store); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})
	if _, err := ops.BlockRecipe("RCP05", "waiting on the client"); err != nil {
		t.Fatal(err)
	}
	m := loadModel(t, store)

	var found bool
	for _, r := range m.recipes {
		if r.ID != "RCP05" {
			continue
		}
		found = true
		out := StripANSI(strings.Join(m.recipeDetail(&r, 100), "\n"))
		if !strings.Contains(out, "On hold: waiting on the client") {
			t.Errorf("the detail should say the recipe is on hold, and why:\n%s", out)
		}
	}
	if !found {
		t.Fatal("the blocked recipe is not in the plan")
	}
}
