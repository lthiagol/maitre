package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
	"github.com/lthiagol/maitre/internal/settings"
)

// settingsStore seeds a plan to hold the tool's settings record.
func settingsStore(t *testing.T) *core.Store {
	t.Helper()
	root := filepath.Join(t.TempDir(), core.PlanDirName)
	store := core.Open(root)
	if err := core.Seed(store); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return store
}

// writeRecord writes the tool's settings record over a plan.
func writeRecord(t *testing.T, store *core.Store, view settings.View) {
	t.Helper()
	if err := settings.Write(settings.Path(store.Root), settings.Record{View: view}); err != nil {
		t.Fatalf("write the settings: %v", err)
	}
}

// TestTheConsoleReadsTheSettingsRecord covers the read at the start.
//
// The chef's choices are in the file, so the console starts with them.
func TestTheConsoleReadsTheSettingsRecord(t *testing.T) {
	store := settingsStore(t)
	writeRecord(t, store, settings.View{Theme: "Dracula", Order: "lane", ShowServed: false})

	m := loadModel(t, store)
	if m.theme.name != "Dracula" {
		t.Errorf("the theme is %q, want Dracula", m.theme.name)
	}
	if m.order != OrderLane {
		t.Errorf("the order is %v, want lane", m.order)
	}
	if m.showServed {
		t.Error("the served recipes are shown, want them hidden")
	}
}

// TestTheConsoleStartsWithTheDefaults covers a plan with no settings file.
//
// A missing file is not an error: the console looks exactly as it did before the
// settings existed.
func TestTheConsoleStartsWithTheDefaults(t *testing.T) {
	store := settingsStore(t)

	m := loadModel(t, store)
	if m.theme.name != "Kitchen" {
		t.Errorf("the theme is %q, want Kitchen", m.theme.name)
	}
	if m.order != OrderTriage {
		t.Errorf("the order is %v, want triage", m.order)
	}
	if !m.showServed {
		t.Error("the served recipes are hidden, want them shown")
	}
	if m.status != "" {
		t.Errorf("the console says %q, want nothing said", m.status)
	}
}

// TestTheChefsChoicesSurviveARestart covers the whole loop.
//
// The chef picks a theme, an order and the served toggle; the console writes them;
// the next console starts with them.
func TestTheChefsChoicesSurviveARestart(t *testing.T) {
	store := settingsStore(t)
	m := loadModel(t, store)

	// The theme, on the Settings tab, the way the chef picks it.
	m.tab = TabSettings
	m.selTheme = themeIndex("Catppuccin Mocha")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	// The order and the served toggle.
	m.Update(tea.KeyPressMsg{Code: '2'})
	m.Update(tea.KeyPressMsg{Code: 'c'})

	// The record on disk holds the choices.
	rec, err := settings.Read(settings.Path(store.Root))
	if err != nil {
		t.Fatalf("read the settings: %v", err)
	}
	if rec.View.Theme != "Catppuccin Mocha" {
		t.Errorf("the record's theme is %q, want Catppuccin Mocha", rec.View.Theme)
	}
	if rec.View.Order != "lane" {
		t.Errorf("the record's order is %q, want lane", rec.View.Order)
	}
	if rec.View.ShowServed {
		t.Error("the record shows the served recipes, want them hidden")
	}

	// And a new console over the same plan starts with them.
	again := loadModel(t, store)
	if again.theme.name != "Catppuccin Mocha" {
		t.Errorf("the theme is %q, want Catppuccin Mocha", again.theme.name)
	}
	if again.order != OrderLane {
		t.Errorf("the order is %v, want lane", again.order)
	}
	if again.showServed {
		t.Error("the served recipes are shown, want them hidden")
	}
}

// TestTheSettingsAreNotThePlansRecords covers the ownership, from the console.
//
// The record is a file beside the plan's records: the journal does not hear of it.
func TestTheSettingsAreNotThePlansRecords(t *testing.T) {
	store := settingsStore(t)
	before, err := store.LoadJournal()
	if err != nil {
		t.Fatalf("load the journal: %v", err)
	}

	m := loadModel(t, store)
	m.tab = TabSettings
	m.selTheme = themeIndex("Dracula")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	after, err := store.LoadJournal()
	if err != nil {
		t.Fatalf("load the journal: %v", err)
	}
	if len(after.Events) != len(before.Events) {
		t.Errorf("the journal holds %d events, want %d — the settings are not the plan's records",
			len(after.Events), len(before.Events))
	}
}

// TestAnUnknownVersionIsRefused covers the version fence, in the console.
//
// The file is not readable, so the console says so and keeps its own defaults.
func TestAnUnknownVersionIsRefused(t *testing.T) {
	store := settingsStore(t)
	writeRecord(t, store, settings.View{Theme: "Dracula", Order: "triage", ShowServed: true})
	path := settings.Path(store.Root)
	if err := os.WriteFile(path, []byte(`{"version": 99, "view": {"theme": "Dracula"}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	m := loadModel(t, store)
	if m.theme.name != "Kitchen" {
		t.Errorf("the theme is %q, want the default Kitchen", m.theme.name)
	}
	if !strings.Contains(m.status, "99") {
		t.Errorf("the console says %q, want it to name the version it refused", m.status)
	}
}
