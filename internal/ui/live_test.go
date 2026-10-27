package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// liveModel builds a console that watches the disk, over a plan with a recipe that
// carries open tasting notes.
//
// The poll interval is a millisecond, so the test does not wait a second for the
// timer.
func liveModel(t *testing.T) (*Model, *core.Store) {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})
	chef := core.NewOperations(store, core.AuditContext{Actor: "chef", Session: "s-2", Role: "chef"})
	reviewer := core.NewOperations(store, core.AuditContext{Actor: "reviewer", Session: "s-3", Role: "reviewer"})
	addExecutedRecipe(t, ops, chef, reviewer, "The chef watches the plan", 2)

	m := NewLive(store)
	m.pollInterval = time.Millisecond
	updated, _ := m.Update(m.Init()())
	m = updated.(*Model)
	m.width = 100
	m.height = 40
	m.ready = true
	return m, store
}

// pollTick runs one poll and applies what it produced.
//
// It is the poll the interactive program runs once a second, without the timer: the
// test drives the message by hand. The message it returns says what the poll did —
// a reload, or the next poll.
func pollTick(t *testing.T, m *Model) (*Model, tea.Msg) {
	t.Helper()
	updated, cmd := m.Update(pollMsg{})
	m = updated.(*Model)
	if cmd == nil {
		return m, nil
	}
	msg := cmd()
	if reload, ok := msg.(reloadMsg); ok {
		updated, _ = m.Update(reload)
		m = updated.(*Model)
	}
	return m, msg
}

// TestAChangeOnDiskAppears covers the live console.
//
// The chef does not press `r`: the plan changed on disk, and the console sees it.
func TestAChangeOnDiskAppears(t *testing.T) {
	m, store := liveModel(t)

	// An agent renames the recipe on disk.
	ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-9", Role: "runner"})
	title := "Renamed by an agent"
	if _, err := ops.UpdateRecipe("RCP01", core.Fragment{Title: &title}); err != nil {
		t.Fatal(err)
	}

	m, msg := pollTick(t, m)
	if _, ok := msg.(reloadMsg); !ok {
		t.Fatalf("the poll should have reloaded the plan, got %T", msg)
	}
	out := StripANSI(m.View().Content)
	if !strings.Contains(out, title) {
		t.Errorf("the console should show the new title:\n%s", out)
	}
}

// TestThePollDoesNotReloadWhenNothingChanged covers the cost.
//
// A poll a second is cheap; a reload a second is not. The plan is read only when
// something has actually changed.
func TestThePollDoesNotReloadWhenNothingChanged(t *testing.T) {
	m, _ := liveModel(t)

	_, first := pollTick(t, m)
	if _, ok := first.(reloadMsg); ok {
		t.Fatalf("the first poll reloaded a plan that had not changed")
	}
	_, second := pollTick(t, m)
	if _, ok := second.(reloadMsg); ok {
		t.Fatalf("the second poll reloaded a plan that had not changed")
	}
}

// TestTheChefsPlaceSurvivesAReload covers the reload.
//
// A reload replaces the plan's data. It must not replace the chef's place: the
// recipe, the scroll, the note cursor, the theme, the order, the served checkbox and
// the tab are the chef's, not the plan's.
func TestTheChefsPlaceSurvivesAReload(t *testing.T) {
	m, store := liveModel(t)

	// The chef's place.
	m.tab = TabRecipes
	m.order = OrderLane
	m.showServed = false
	m.selRecipes = 0
	m.selNote = 1
	m.detailScroll = 3
	m.theme = newTheme(catppuccinMocha())

	// Something changes on disk.
	ops := core.NewOperations(store, core.AuditContext{Actor: "runner", Session: "s-9", Role: "runner"})
	if _, err := ops.DocumentRecipe("RCP01", "the plan moved on"); err != nil {
		t.Fatal(err)
	}

	m, msg := pollTick(t, m)
	if _, ok := msg.(reloadMsg); !ok {
		t.Fatalf("the poll should have reloaded the plan, got %T", msg)
	}

	if m.tab != TabRecipes {
		t.Errorf("the tab = %v, want %v", m.tab, TabRecipes)
	}
	if m.order != OrderLane {
		t.Errorf("the order = %v, want %v", m.order, OrderLane)
	}
	if m.showServed {
		t.Error("the served checkbox should have survived")
	}
	if m.selRecipes != 0 {
		t.Errorf("the selected recipe = %d, want 0", m.selRecipes)
	}
	if m.selNote != 1 {
		t.Errorf("the note cursor = %d, want 1", m.selNote)
	}
	if m.detailScroll != 3 {
		t.Errorf("the detail scroll = %d, want 3", m.detailScroll)
	}
	if m.theme.name != "Catppuccin Mocha" {
		t.Errorf("the theme = %q, want Catppuccin Mocha", m.theme.name)
	}
}

// TestAFileComingAndGoingDoesNotBreakThePoll covers a plan that is being written to.
//
// A file that disappears between two stats, or appears between two listings, must
// not break the console.
func TestAFileComingAndGoingDoesNotBreakThePoll(t *testing.T) {
	m, store := liveModel(t)

	// The recipe's file, by the name the plan gave it.
	entries, err := os.ReadDir(filepath.Join(store.Root, core.DirRecipes))
	if err != nil {
		t.Fatal(err)
	}
	var path string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			path = filepath.Join(store.Root, core.DirRecipes, e.Name())
		}
	}
	if path == "" {
		t.Fatal("the plan holds no recipe")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	// The file goes away between two polls.
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	m, msg := pollTick(t, m)
	if _, ok := msg.(reloadMsg); !ok {
		t.Fatalf("the poll should have reloaded the plan, got %T", msg)
	}
	if len(m.recipes) != 0 {
		t.Fatalf("the console holds %d recipes, want 0 — the file is gone", len(m.recipes))
	}
	if m.err != nil {
		t.Fatalf("the poll broke: %v", m.err)
	}
	if StripANSI(m.View().Content) == "" {
		t.Error("the console drew nothing")
	}

	// And it comes back.
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	m, msg = pollTick(t, m)
	if _, ok := msg.(reloadMsg); !ok {
		t.Fatalf("the poll should have reloaded the plan, got %T", msg)
	}
	if len(m.recipes) != 1 {
		t.Fatalf("the console holds %d recipes, want 1", len(m.recipes))
	}
}

// TestTheSnapshotPathDoesNotPoll covers the non-interactive programs.
//
// `--render` and `--keys` are not interactive programs. A poll loop in one of them
// would keep it alive, and the script would never end.
func TestTheSnapshotPathDoesNotPoll(t *testing.T) {
	root := filepath.Join(t.TempDir(), core.PlanDirName)
	store := core.Open(root)
	if err := core.Seed(store); err != nil {
		t.Fatalf("seed: %v", err)
	}
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	// The reload the snapshot does schedules no poll.
	updated, cmd := m.Update(m.reload()())
	if cmd != nil {
		t.Fatal("the snapshot path scheduled a poll")
	}
	m = updated.(*Model)
	if cmd := m.poll(); cmd != nil {
		t.Fatal("the snapshot model scheduled a poll")
	}
}

// TestTheInteractiveModelPolls covers the interactive program.
//
// The poll message produces exactly one follow-up poll: not two, not none.
func TestTheInteractiveModelPolls(t *testing.T) {
	m, _ := liveModel(t)

	updated, cmd := m.Update(pollMsg{})
	m = updated.(*Model)
	if cmd == nil {
		t.Fatal("the poll produced no follow-up")
	}
	msg := cmd()
	if _, ok := msg.(pollMsg); !ok {
		t.Fatalf("the poll should produce another poll, got %T", msg)
	}
}
