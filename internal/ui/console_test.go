package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
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
// loadModel builds a console over a plan, through the live path.
//
// The message comes from the reload and goes through the model's own Update — the
// same path the running console takes. A helper that copies the fields by hand
// tests the copy, and the copy is where the bug was.
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
	m.Update(loaded)
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
	// RCP02 is the ready one, and triage puts it first.
	if s, ok := m.currentSummary(); !ok || s.ID != "RCP02" {
		t.Fatalf("the first recipe should be RCP02, got %+v", s)
	}
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the sign-off was refused: %s", m.refusal)
	}
	if !strings.Contains(m.status, "signed off RCP02") {
		t.Fatalf("status = %q", m.status)
	}
	// The recipe really moved.
	rec, err := m.ops.GetRecipe("RCP02")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Lifecycle != core.LifecycleApproved {
		t.Fatalf("lifecycle = %q, want approved", rec.Lifecycle)
	}
}

func TestSpaceDoesNotStartWork(t *testing.T) {
	m := newModel(t)
	// Move to RCP04, which is approved.
	for {
		s, _ := m.currentSummary()
		if s.ID == "RCP04" {
			break
		}
		m.Update(tea.KeyPressMsg{Code: 'j'})
	}
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("unexpected refusal: %s", m.refusal)
	}
	rec, err := m.ops.GetRecipe("RCP04")
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
	// RCP07 is a draft — there is nothing for the chef to do.
	for {
		s, _ := m.currentSummary()
		if s.ID == "RCP07" {
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

func TestBacklogTabListsTheItems(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	out := plain(m)
	for _, want := range []string{"IDEAS", "SMALL TASKS", "The board needs a way to show evidence", "Wire the review queue", "idea", "feature"} {
		if !strings.Contains(out, want) {
			t.Errorf("the backlog is missing %q\n%s", want, out)
		}
	}
}

func TestPromotingAnItemFromTheBacklog(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the promotion was refused: %s", m.refusal)
	}
	items, err := m.ops.ListBacklogItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 6 {
		t.Fatalf("got %d items, want 6", len(items))
	}
}

// TestTheBacklogGroupsAreDerivedFromTheKind covers the two groups.
//
// The group is never stored: the kind is the only truth, and the tab's groups are
// read off it.
func TestTheBacklogGroupsAreDerivedFromTheKind(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	out := plain(m)
	ideas := strings.Index(out, "IDEAS")
	tasks := strings.Index(out, "SMALL TASKS")
	if ideas < 0 || tasks < 0 {
		t.Fatalf("the backlog has no groups\n%s", out)
	}
	if ideas > tasks {
		t.Fatalf("the ideas come before the small tasks\n%s", out)
	}
	// An idea is filed under the ideas, a small task under the small tasks.
	ideaRow := strings.Index(out, "BL01")
	taskRow := strings.Index(out, "BL05")
	if ideaRow < ideas || ideaRow > tasks {
		t.Fatalf("the idea is not under the ideas\n%s", out)
	}
	if taskRow < tasks {
		t.Fatalf("the small task is not under the small tasks\n%s", out)
	}
}

// TestAPlanWithRetiredArtifactsIsNotSupported covers the reader's answer.
//
// The plan format is one format. A plan that carries an artifact of a retired one —
// the old backlog directories — is not supported, and the console says so plainly
// instead of showing the chef an empty backlog and calling it the truth.
func TestAPlanWithRetiredArtifactsIsNotSupported(t *testing.T) {
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	writeRetired(t, root, "specials", "S01-a-hunch.json",
		`{"id":"S01","title":"A hunch","body":"worth noticing","status":"open"}`)

	m := New(store)
	msg := m.reload()()
	loaded, ok := msg.(reloadMsg)
	if !ok {
		t.Fatal("unexpected reload result")
	}
	if loaded.err == nil {
		t.Fatal("the console should not read a plan that carries retired artifacts")
	}
	m.Update(loaded)
	m.width, m.height, m.ready = 100, 40, true
	if m.err == nil {
		t.Fatal("the console should hold the error")
	}
	if !strings.Contains(m.err.Error(), "not supported") {
		t.Fatalf("the console should say the plan is not supported: %v", m.err)
	}
	out := plain(m)
	if !strings.Contains(out, "not supported") {
		t.Fatalf("the console should say the plan is not supported\n%s", out)
	}
}

// writeRetired writes one file of the plan format before the rename.
func writeRetired(t *testing.T, root, dir, name, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPathTabSeparatesReadyFromWaiting(t *testing.T) {
	m := newModel(t)
	m.tab = TabPath
	out := plain(m)
	if !strings.Contains(out, "READY") || !strings.Contains(out, "WAITING") {
		t.Errorf("the path is missing a group\n%s", out)
	}
	if !strings.Contains(out, "needs RCP02") {
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
