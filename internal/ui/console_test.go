package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/kitchen-brigade/internal/core"
)

// newModel builds a console over a seeded plan.
func newModel(t *testing.T) *Model {
	t.Helper()
	root := filepath.Join(t.TempDir(), core.PlanDirName)
	store := core.Open(root)
	if err := core.Seed(store); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return loadModel(t, store)
}

// loadModel reads a plan into a console model.
func loadModel(t *testing.T, store *core.Store) *Model {
	t.Helper()
	m := New(store)
	msg := m.reload()()
	loaded, ok := msg.(reloadMsg)
	if !ok {
		t.Fatalf("unexpected reload result")
	}
	if loaded.err != nil {
		t.Fatalf("reload: %v", loaded.err)
	}
	m.recipes = loaded.recipes
	m.specials = loaded.specials
	m.notes = loaded.notes
	m.journal = loaded.journal
	m.sections = m.orderedSections()
	m.clampSelections()
	m.width = 100
	m.height = 40
	m.ready = true
	return m
}

// plain is the view with the escape sequences removed.
func plain(m *Model) string { return StripANSI(m.View().Content) }

func TestConsoleShowsTheTabs(t *testing.T) {
	m := newModel(t)
	out := plain(m)
	for _, name := range TabNames {
		if !strings.Contains(out, name) {
			t.Errorf("the tab bar is missing %q\n%s", name, out)
		}
	}
}

func TestRecipesTabGroupsByTriage(t *testing.T) {
	m := newModel(t)
	out := plain(m)
	for _, section := range []string{core.TriageNeedsYou, core.TriageOnTheLine, core.TriageBeingShaped, core.TriageServed} {
		if !strings.Contains(out, section) {
			t.Errorf("the board is missing the %q section\n%s", section, out)
		}
	}
	if !strings.Contains(out, "Sign off the menu") {
		t.Errorf("the board is missing a recipe\n%s", out)
	}
}

func TestServedCheckboxHidesTheServedRecipes(t *testing.T) {
	m := newModel(t)
	if !strings.Contains(plain(m), "Draft the glossary") {
		t.Fatal("the served recipe should be visible by default")
	}
	m.Update(tea.KeyPressMsg{Code: 'c'})
	out := plain(m)
	if strings.Contains(out, "Draft the glossary") {
		t.Errorf("the served recipe should be hidden\n%s", out)
	}
	if !strings.Contains(out, "show served") {
		t.Errorf("the checkbox should still be shown\n%s", out)
	}
}

func TestOrderKeysSwitchTheOrdering(t *testing.T) {
	m := newModel(t)
	for _, k := range []rune{'1', '2', '3'} {
		m.Update(tea.KeyPressMsg{Code: k})
	}
	if m.order != OrderPlain {
		t.Fatalf("order = %v, want plain", m.order)
	}
	m.Update(tea.KeyPressMsg{Code: '1'})
	if m.order != OrderTriage {
		t.Fatalf("order = %v, want triage", m.order)
	}
}

func TestTabKeySwitchesTabs(t *testing.T) {
	m := newModel(t)
	if m.tab != TabRecipes {
		t.Fatal("the console should open on Recipes")
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.tab != TabBacklog {
		t.Fatalf("tab = %v, want Backlog", m.tab)
	}
	for i := 0; i < len(TabNames)-1; i++ {
		m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	}
	if m.tab != TabRecipes {
		t.Fatalf("tab = %v, want Recipes after a full cycle", m.tab)
	}
}

func TestSpaceSignsOffAReadyRecipe(t *testing.T) {
	m := newModel(t)
	// R02 is the ready one, and triage puts it first.
	if s, ok := m.currentSummary(); !ok || s.ID != "R02" {
		t.Fatalf("the first recipe should be R02, got %+v", s)
	}
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the sign-off was refused: %s", m.refusal)
	}
	if !strings.Contains(m.status, "signed off R02") {
		t.Fatalf("status = %q", m.status)
	}
	// The recipe really moved.
	rec, err := m.ops.GetRecipe("R02")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Lifecycle != core.LifecycleApproved {
		t.Fatalf("lifecycle = %q, want approved", rec.Lifecycle)
	}
}

func TestSpaceDoesNotStartWork(t *testing.T) {
	m := newModel(t)
	// Move to R04, which is approved.
	for {
		s, _ := m.currentSummary()
		if s.ID == "R04" {
			break
		}
		m.Update(tea.KeyPressMsg{Code: 'j'})
	}
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("unexpected refusal: %s", m.refusal)
	}
	rec, err := m.ops.GetRecipe("R04")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Lifecycle == core.LifecycleExecuting {
		t.Fatal("the console must not start work — that is an agent's move")
	}
	if rec.Execution.ClaimedBy != "" {
		t.Fatalf("the console must not claim work, got claimed by %q", rec.Execution.ClaimedBy)
	}
}

func TestNothingToDoIsNotARefusal(t *testing.T) {
	m := newModel(t)
	// R07 is a draft — there is nothing for the chef to do.
	for {
		s, _ := m.currentSummary()
		if s.ID == "R07" {
			break
		}
		m.Update(tea.KeyPressMsg{Code: 'j'})
	}
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("there was nothing to refuse: %s", m.refusal)
	}
	if !strings.Contains(m.status, "nothing for you") {
		t.Fatalf("status = %q", m.status)
	}
}

func TestEnterOpensTheDetailFullScreen(t *testing.T) {
	m := newModel(t)
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.mode != ModeDetail {
		t.Fatalf("mode = %v, want detail", m.mode)
	}
	out := plain(m)
	for _, want := range []string{"Outcome", "Acceptance criteria", "Method", "Sign off the menu"} {
		if !strings.Contains(out, want) {
			t.Errorf("the detail pane is missing %q\n%s", want, out)
		}
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if m.mode != ModeSplit {
		t.Fatalf("mode = %v, want split", m.mode)
	}
}

func TestBacklogTabListsSpecialsAndPrepNotes(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	out := plain(m)
	for _, want := range []string{"SPECIALS", "PREP NOTES", "The board needs a way to show proofs", "Wire the review queue"} {
		if !strings.Contains(out, want) {
			t.Errorf("the backlog is missing %q\n%s", want, out)
		}
	}
}

func TestPromotingASpecialFromTheBacklog(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the promotion was refused: %s", m.refusal)
	}
	notes, err := m.ops.ListPrepNotes()
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 3 {
		t.Fatalf("got %d prep notes, want 3", len(notes))
	}
}

func TestPathTabSeparatesReadyFromWaiting(t *testing.T) {
	m := newModel(t)
	m.tab = TabPath
	out := plain(m)
	if !strings.Contains(out, "READY") || !strings.Contains(out, "WAITING") {
		t.Errorf("the path is missing a group\n%s", out)
	}
	if !strings.Contains(out, "needs R02") {
		t.Errorf("the path should say what is blocking\n%s", out)
	}
}

func TestLogTabShowsTheJournal(t *testing.T) {
	m := newModel(t)
	m.tab = TabLog
	out := plain(m)
	if !strings.Contains(out, "recent activity") {
		t.Errorf("the log is missing its header\n%s", out)
	}
	if !strings.Contains(out, "created") {
		t.Errorf("the log should show a creation\n%s", out)
	}
}

func TestSettingsTabShowsThePlan(t *testing.T) {
	m := newModel(t)
	m.tab = TabSettings
	out := plain(m)
	for _, want := range []string{"PLAN", "recipes", "VIEW", "KEYS"} {
		if !strings.Contains(out, want) {
			t.Errorf("the settings are missing %q\n%s", want, out)
		}
	}
}

func TestQuitEmitsTheQuitCommand(t *testing.T) {
	m := newModel(t)
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'q'})
	if cmd == nil {
		t.Fatal("expected a command on quit")
	}
	if msg := cmd(); msg == nil {
		t.Fatal("expected a quit message")
	}
}
