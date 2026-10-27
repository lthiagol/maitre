package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// noteModel builds a console over a recipe with two open tasting notes.
//
// The recipe is executed and the reviewer has filed two notes. The close gate refuses
// it while a note is open, and the chef has two notes to triage.
func noteModel(t *testing.T) *Model {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})
	chef := core.NewOperations(store, core.AuditContext{Actor: "chef", Session: "s-2", Role: "chef"})
	reviewer := core.NewOperations(store, core.AuditContext{Actor: "reviewer", Session: "s-3", Role: "reviewer"})

	addExecutedRecipe(t, ops, chef, reviewer, "The chef can triage a tasting note", 2)

	m := loadModel(t, store)
	m.height = 30
	return m
}

// addExecutedRecipe files a recipe, executes it, and gives it open tasting notes.
//
// The reviewer files the notes: they are somebody else's complaint, and the chef is
// the one who answers them.
func addExecutedRecipe(t *testing.T, ops, chef, reviewer *core.Operations, title string, notes int) string {
	t.Helper()
	r, err := ops.CreateRecipe(core.RecipeSpec{
		Title:      title,
		Outcome:    "The note is a conversation the chef can join",
		Problem:    "The console showed notes with nothing to press",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []core.CriterionSpec{
			{Description: "the console can dismiss a note", Verification: "go test ./internal/ui -run TestTheChefCanDismissANote"},
		},
		WorkPackages: []string{"ui"},
		Steps:        []core.StepSpec{{Action: "triage the note", DoneWhen: "the note is closed", Covers: []string{"AC-01"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range []core.Event{core.EventGroom, core.EventSpecify} {
		if _, err := ops.Advance(r.ID, ev); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := chef.ApproveRecipe(r.ID, "chef"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.ClaimRecipe(r.ID, "runner"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.SetStepStatus(r.ID, "G1", "done"); err != nil {
		t.Fatal(err)
	}
	if _, err := chef.SignOffCriterion(r.ID, "AC-01", "chef", "checked"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, core.EventFinish); err != nil {
		t.Fatal(err)
	}
	descriptions := []string{
		"the queue drops findings on refresh",
		"the thread is cut off at the bottom",
	}
	for i := 0; i < notes; i++ {
		if _, err := reviewer.AddFinding(r.ID, core.Finding{
			Phase:       core.FindingPhaseExternal,
			Severity:    core.SeverityHigh,
			Category:    "correctness",
			Description: descriptions[i%len(descriptions)],
		}); err != nil {
			t.Fatal(err)
		}
	}
	return r.ID
}

// triageKeys drives a whole triage: space, an action, the reason, enter.
func triageKeys(action, text string) []string {
	keys := []string{"space", action}
	keys = append(keys, typed(text)...)
	return append(keys, "enter")
}

// noteLine is the rendered line that carries a tasting note's id.
func noteLine(t *testing.T, raw, id string) string {
	t.Helper()
	for _, l := range strings.Split(raw, "\n") {
		if strings.Contains(StripANSI(l), id) {
			return l
		}
	}
	t.Fatalf("no line carries %q:\n%s", id, StripANSI(raw))
	return ""
}

// TestTheChefCanSelectANote covers the cursor.
//
// One note is selected, and the selection is visible: the console's own selected-row
// style, so "selected" means one thing on this screen.
func TestTheChefCanSelectANote(t *testing.T) {
	m := noteModel(t)
	notes := m.openNotes()
	if len(notes) != 2 {
		t.Fatalf("the fixture should have two open notes, got %d", len(notes))
	}
	// The first note is under the cursor to begin with.
	first, ok := m.selectedNote()
	if !ok || first.ID != notes[0].ID {
		t.Fatalf("the cursor should start on %q", notes[0].ID)
	}

	if _, err := Drive(m, []string{"n"}); err != nil {
		t.Fatal(err)
	}
	note, ok := m.selectedNote()
	if !ok {
		t.Fatal("no note under the cursor")
	}
	if note.ID != notes[1].ID {
		t.Fatalf("the cursor is on %q, want %q", note.ID, notes[1].ID)
	}
	// And back again.
	if _, err := Drive(m, []string{"N"}); err != nil {
		t.Fatal(err)
	}
	back, ok := m.selectedNote()
	if !ok || back.ID != notes[0].ID {
		t.Fatalf("the cursor should be back on %q", notes[0].ID)
	}
	if _, err := Drive(m, []string{"n"}); err != nil {
		t.Fatal(err)
	}
	note, _ = m.selectedNote()

	// The selection is visible: the note under the cursor wears the console's own
	// selected-row style, and the other note does not.
	rec, ok := m.currentRecipe()
	if !ok {
		t.Fatal("no recipe selected")
	}
	raw := strings.Join(m.recipeDetail(rec, m.width), "\n")
	marked := m.stySelected.Render("X")
	prefix := marked[:strings.Index(marked, "X")]
	if prefix == "" {
		t.Fatal("the selected style renders no escape sequences — the check below would prove nothing")
	}
	selected := noteLine(t, raw, note.ID)
	other := noteLine(t, raw, notes[0].ID)
	if !strings.Contains(selected, prefix) {
		t.Errorf("the note under the cursor should wear the selection style:\n%q", selected)
	}
	if strings.Contains(other, prefix) {
		t.Errorf("only the note under the cursor should wear the selection style:\n%q", other)
	}
}

// TestTheChefCanDismissANote covers the chef's move on a note.
//
// A dismissal closes the note without fixing it, and it carries the reason the chef
// typed.
func TestTheChefCanDismissANote(t *testing.T) {
	m := noteModel(t)
	notes := m.openNotes()
	if _, err := Drive(m, triageKeys("d", "the queue drops it on refresh")); err != nil {
		t.Fatal(err)
	}

	rec, err := m.ops.GetRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range rec.Review.Findings {
		if f.ID != notes[0].ID {
			continue
		}
		if f.Status != "dismissed" {
			t.Fatalf("the note's status = %q, want dismissed", f.Status)
		}
		if f.Resolution != "the queue drops it on refresh" {
			t.Fatalf("the note's reason = %q, want what the chef typed", f.Resolution)
		}
		return
	}
	t.Fatalf("the note %q is gone", notes[0].ID)
}

// TestTheChefCanCommentOnANote covers the other chef's move.
//
// The comment belongs to the note, not to the recipe: the thread is the channel
// between the reviewer and the fixer.
func TestTheChefCanCommentOnANote(t *testing.T) {
	m := noteModel(t)
	notes := m.openNotes()
	if _, err := Drive(m, triageKeys("c", "we know, it is on the list")); err != nil {
		t.Fatal(err)
	}

	rec, err := m.ops.GetRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	onNote, onRecipe := 0, 0
	for _, c := range rec.Review.Comments {
		switch c.Finding {
		case notes[0].ID:
			onNote++
			if c.Author != "chef" {
				t.Errorf("the comment should be the chef's, got %q", c.Author)
			}
			if c.Body != "we know, it is on the list" {
				t.Errorf("the comment should be what was typed, got %q", c.Body)
			}
		case "":
			onRecipe++
		}
	}
	if onNote != 1 {
		t.Fatalf("the note should carry one comment, got %d", onNote)
	}
	if onRecipe != 0 {
		t.Errorf("the comment should not be on the recipe, got %d there", onRecipe)
	}

	// The result is visible: the console reloaded, and the comment is in the note's
	// thread on the screen.
	rec, ok := m.currentRecipe()
	if !ok {
		t.Fatal("no recipe selected")
	}
	out := StripANSI(strings.Join(m.recipeDetail(rec, m.width), "\n"))
	if !strings.Contains(out, "chef: we know, it is on the list") {
		t.Errorf("the comment should show in the note's thread:\n%s", out)
	}
}

// TestTheConsoleDoesNotOfferToResolve covers the surface split.
//
// Resolving a note is the fixer's move — the agent that did the work. The console
// offers the chef's moves and nothing else.
func TestTheConsoleDoesNotOfferToResolve(t *testing.T) {
	m := noteModel(t)
	out, err := Drive(m, []string{"space"})
	if err != nil {
		t.Fatal(err)
	}
	plain := StripANSI(out)
	if strings.Contains(plain, "resolve") {
		t.Errorf("the prompt should not offer the fixer's move:\n%s", plain)
	}
	if !strings.Contains(plain, "d) dismiss") || !strings.Contains(plain, "c) comment") {
		t.Errorf("the prompt should offer the chef's moves:\n%s", plain)
	}
}

// TestTheTriageCanBeCancelled checks esc leaves no trace.
func TestTheTriageCanBeCancelled(t *testing.T) {
	m := noteModel(t)
	before, err := m.ops.GetRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}

	keys := append([]string{"space", "d"}, append(typed("half a reason"), "esc")...)
	if _, err := Drive(m, keys); err != nil {
		t.Fatal(err)
	}
	if m.prompt != nil {
		t.Fatal("the prompt should be closed")
	}

	after, err := m.ops.GetRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	for i, f := range after.Review.Findings {
		if f.ID != before.Review.Findings[i].ID {
			t.Fatalf("the notes changed: %q is now %q", before.Review.Findings[i].ID, f.ID)
		}
		if f.Status != before.Review.Findings[i].Status {
			t.Fatalf("a cancelled triage changed the note's status to %q", f.Status)
		}
	}
	if len(after.Review.Comments) != len(before.Review.Comments) {
		t.Errorf("a cancelled triage left %d comments", len(after.Review.Comments)-len(before.Review.Comments))
	}
}

// TestTheFooterNamesTheNoteAndTheTriage covers the contract with the chef.
//
// The footer is where the console says what the next key does. A key that is not
// named is a key the chef cannot find.
func TestTheFooterNamesTheNoteAndTheTriage(t *testing.T) {
	m := noteModel(t)
	out := StripANSI(m.View().Content)
	if !strings.Contains(out, "n/N note") {
		t.Errorf("the footer should say how to move the note cursor:\n%s", out)
	}
	if !strings.Contains(out, "space triage") {
		t.Errorf("the footer should offer the triage:\n%s", out)
	}
}

// TestTheDismissalIsJournalled covers the journal.
//
// An agent polling for changes has to see the dismissal: it is the chef's decision,
// and the journal is the only event log.
func TestTheDismissalIsJournalled(t *testing.T) {
	m := noteModel(t)
	before, err := m.store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := Drive(m, triageKeys("d", "not a bug after all")); err != nil {
		t.Fatal(err)
	}

	after, err := m.store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Events) != len(before.Events)+1 {
		t.Fatalf("the journal grew by %d events, want 1", len(after.Events)-len(before.Events))
	}
	last := after.Events[len(after.Events)-1]
	if last.Kind != core.KindFinding {
		t.Errorf("the event's kind = %q, want %q", last.Kind, core.KindFinding)
	}
	if last.Actor != "chef" {
		t.Errorf("the event should be the chef's, got %q", last.Actor)
	}
}
