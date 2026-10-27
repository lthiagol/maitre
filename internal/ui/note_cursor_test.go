package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// twoNoteModel builds a console over two recipes, each with open tasting notes.
//
// R01 has three notes, R02 has one. The Path tab lists the recipes by id and the
// console keeps a selection per tab, so the two tabs can show different recipes —
// which is how the note cursor could end up pointing past the end of the recipe the
// console was showing.
func twoNoteModel(t *testing.T) *Model {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "line-cook", Session: "s-1", Role: "runner"})
	chef := core.NewOperations(store, core.AuditContext{Actor: "chef", Session: "s-2", Role: "coordinator"})
	taster := core.NewOperations(store, core.AuditContext{Actor: "taster", Session: "s-3", Role: "reviewer"})

	addExecutedRecipe(t, ops, chef, taster, "The first recipe", 3)
	addExecutedRecipe(t, ops, chef, taster, "The second recipe", 1)

	m := loadModel(t, store)
	m.height = 30
	return m
}

// switchToTab drives the console to the named tab.
func switchToTab(t *testing.T, m *Model, tab Tab) {
	t.Helper()
	for i := 0; i < len(TabNames); i++ {
		if m.tab == tab {
			return
		}
		if _, err := Drive(m, []string{"tab"}); err != nil {
			t.Fatal(err)
		}
	}
	t.Fatalf("the console did not reach the %s tab", tab.Label())
}

// TestTheNoteCursorFollowsTheRecipe covers the tab switch.
//
// The cursor was one number for the whole console while the recipe is per tab, so a
// switch could leave it past the end of the new recipe's notes: the footer promised
// a triage and the space answered "no open tasting note here".
//
// The invariant is the one that matters: the footer and the space must agree.
func TestTheNoteCursorFollowsTheRecipe(t *testing.T) {
	m := twoNoteModel(t)

	// The Path tab's selection is pointed at R02, which has one note.
	switchToTab(t, m, TabPath)
	if _, err := Drive(m, []string{"down"}); err != nil {
		t.Fatal(err)
	}
	// Back to the Recipes tab, where R01 — three notes — is selected.
	switchToTab(t, m, TabRecipes)
	if _, err := Drive(m, []string{"n", "n"}); err != nil {
		t.Fatal(err)
	}
	if m.selNote != 2 {
		t.Fatalf("the cursor should be on the third note, got %d", m.selNote)
	}

	// And to the Path tab again: R02, one note.
	switchToTab(t, m, TabPath)

	// The cursor is within the recipe the console is showing.
	if notes := m.openNotes(); m.selNote >= len(notes) {
		t.Fatalf("the cursor is on note %d of %d — it did not follow the recipe", m.selNote, len(notes))
	}

	// And the footer and the space agree.
	out := plain(m)
	if !strings.Contains(out, "space triage") {
		return
	}
	if _, err := Drive(m, []string{"space"}); err != nil {
		t.Fatal(err)
	}
	if m.prompt == nil {
		t.Fatalf("the footer promised a triage and the space refused it:\n%s", out)
	}
}

// TestTheSettingsListNamesTheNoteKey covers the console's own key reference.
//
// The key works and the footer names it, but the one place a chef looks up what the
// keys are has to name it too.
func TestTheSettingsListNamesTheNoteKey(t *testing.T) {
	m := noteModel(t)
	m.height = 60
	// The Settings tab.
	switchToTab(t, m, TabSettings)
	out := plain(m)
	if !strings.Contains(out, "n / N") {
		t.Errorf("the key list should name the note cursor:\n%s", out)
	}
	if !strings.Contains(out, "tasting notes") {
		t.Errorf("the key list should say what the note cursor does:\n%s", out)
	}
}

// TestTheTriageCannotOverwriteANoteTheAgentClosed covers the race.
//
// The chef opens the triage, the agent closes the note over the CLI, and the chef
// types the reason. The dismissal must be refused: the fixer's record of what was
// done survives.
func TestTheTriageCannotOverwriteANoteTheAgentClosed(t *testing.T) {
	m := noteModel(t)
	// The chef opens the triage and starts typing.
	if _, err := Drive(m, []string{"space", "d"}); err != nil {
		t.Fatal(err)
	}
	// The agent resolves the note while the prompt is open.
	agent := core.NewOperations(m.store, core.AuditContext{Actor: "line-cook", Session: "s-9", Role: "runner"})
	if _, err := agent.ResolveFinding("R01", "F-01", "fixed it properly"); err != nil {
		t.Fatal(err)
	}
	// The chef types the reason and hits enter.
	if _, err := Drive(m, append(typed("not needed"), "enter")); err != nil {
		t.Fatal(err)
	}

	rec, err := m.ops.GetRecipe("R01")
	if err != nil {
		t.Fatal(err)
	}
	f := rec.Review.Findings[0]
	if f.Status != "fixed" {
		t.Fatalf("the note's status = %q, want fixed — the fixer's record must survive", f.Status)
	}
	if f.Resolution != "fixed it properly" {
		t.Fatalf("the note's resolution = %q, want what the fixer wrote", f.Resolution)
	}
	if m.refusal == "" {
		t.Error("the console should say the plan refused the dismissal")
	}
}
