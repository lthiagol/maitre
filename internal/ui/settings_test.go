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

// TestTheRecordsDefaultsAreTheConsolesDefaults holds the tie between the record
// and the console.
//
// The record names the console's default theme and the console's default order: the
// CLI prints the record, and the CLI does not know the console's vocabulary. So the
// two names are written twice, and this holds them together — a rename of either
// must come with the other.
func TestTheRecordsDefaultsAreTheConsolesDefaults(t *testing.T) {
	if settings.DefaultTheme != palettes()[0].name {
		t.Errorf("the record's default theme is %q, want the console's %q",
			settings.DefaultTheme, palettes()[0].name)
	}
	if settings.DefaultOrder != OrderNames[0] {
		t.Errorf("the record's default order is %q, want the console's %q",
			settings.DefaultOrder, OrderNames[0])
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
	if !strings.Contains(m.said(), "99") {
		t.Errorf("the console says %q, want it to name the version it refused", m.said())
	}
}

// TestARefusedRecordIsNotWritten covers the fence from the other side.
//
// A record this build cannot read is **not written**: a write would replace the
// chef's file with the defaults, and a file at a version this build does not know
// is a file this build does not touch. The console says so instead.
func TestARefusedRecordIsNotWritten(t *testing.T) {
	store := settingsStore(t)
	path := settings.Path(store.Root)
	whole := `{"version": 1, "view": {"theme": "Dracula", "order": "lane", "showServed": false}}`
	if err := os.WriteFile(path, []byte(whole), 0o644); err != nil {
		t.Fatal(err)
	}

	m := loadModel(t, store)
	// The chef changes something the console writes.
	m.tab = TabSettings
	m.selTheme = themeIndex("Dracula")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != whole {
		t.Errorf("the refused record was written:\n%s", raw)
	}
	if !strings.Contains(m.said(), "version 1") {
		t.Errorf("the console says %q, want it to name the version it refused", m.said())
	}
}

// TestTheConsoleReadsTheCrewSection covers the read of the automation's order.
//
// The record's crew section is the tab's order: the harness, the crew's size, the
// sources and the budget. The console reads them when it starts.
func TestTheConsoleReadsTheCrewSection(t *testing.T) {
	store := settingsStore(t)
	rec := settings.Record{
		View: settings.View{Theme: "Kitchen", Order: "triage", ShowServed: true},
		Crew: settings.Crew{Harness: "pi", Count: 2, Sources: "recipes", Rounds: 5, Interactions: 24},
	}
	if err := settings.Write(settings.Path(store.Root), rec); err != nil {
		t.Fatalf("write the settings: %v", err)
	}

	m := loadModel(t, store)
	if m.picks.Harness != "pi" {
		t.Errorf("the harness is %q, want pi", m.picks.Harness)
	}
	if m.picks.Count != 2 {
		t.Errorf("the crew's size is %d, want 2", m.picks.Count)
	}
	if m.picks.Sources != "recipes" {
		t.Errorf("the sources are %q, want recipes", m.picks.Sources)
	}
	if b := sousBudgets[m.picks.Budget]; b.Rounds != 5 || b.Interactions != 24 {
		t.Errorf("the budget is %+v, want 5 rounds and 24 interactions", b)
	}
}

// TestAnUnknownCrewValueIsRefusedByName covers the refusal of the order's values.
//
// The console does not fall back silently: it names the value it refused, and it
// keeps its own default for that one control only.
func TestAnUnknownCrewValueIsRefusedByName(t *testing.T) {
	store := settingsStore(t)
	rec := settings.Record{
		View: settings.View{Theme: "Kitchen", Order: "triage", ShowServed: true},
		Crew: settings.Crew{Harness: "codex", Count: 7, Sources: "moon", Rounds: 1, Interactions: 1},
	}
	if err := settings.Write(settings.Path(store.Root), rec); err != nil {
		t.Fatalf("write the settings: %v", err)
	}

	m := loadModel(t, store)
	for _, want := range []string{"codex", "7", "moon"} {
		if !strings.Contains(m.status, want) {
			t.Errorf("the console says %q, want it to name %q", m.status, want)
		}
	}
	if m.picks.Harness != settings.DefaultHarness {
		t.Errorf("the harness is %q, want the default %q", m.picks.Harness, settings.DefaultHarness)
	}
	if m.picks.Count != settings.DefaultCrew {
		t.Errorf("the crew's size is %d, want the default %d", m.picks.Count, settings.DefaultCrew)
	}
}

// TestTheWriteKeepsTheCrewSection covers the write.
//
// The record is one file: a write that kept the view and dropped the crew would
// take the automation's order with it.
func TestTheWriteKeepsTheCrewSection(t *testing.T) {
	store := settingsStore(t)
	rec := settings.Record{
		View: settings.View{Theme: "Kitchen", Order: "triage", ShowServed: true},
		Crew: settings.Crew{Harness: "pi", Count: 2, Sources: "recipes", Rounds: 5, Interactions: 24},
	}
	if err := settings.Write(settings.Path(store.Root), rec); err != nil {
		t.Fatalf("write the settings: %v", err)
	}

	m := loadModel(t, store)
	m.tab = TabSettings
	m.selTheme = themeIndex("Dracula")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	got, err := settings.Read(settings.Path(store.Root))
	if err != nil {
		t.Fatalf("read the settings: %v", err)
	}
	if got.Crew != rec.Crew {
		t.Errorf("the crew section is %+v, want %+v", got.Crew, rec.Crew)
	}
	if got.View.Theme != "Dracula" {
		t.Errorf("the theme is %q, want Dracula", got.View.Theme)
	}
}

// TestAnUnknownThemeIsRefusedByName covers the refusal of a theme.
//
// The console does not fall back silently: it says which theme it refused, and it
// uses its own default.
func TestAnUnknownThemeIsRefusedByName(t *testing.T) {
	store := settingsStore(t)
	writeRecord(t, store, settings.View{Theme: "Neon", Order: "triage", ShowServed: true})

	m := loadModel(t, store)
	if m.theme.name != "Kitchen" {
		t.Errorf("the theme is %q, want the default Kitchen", m.theme.name)
	}
	if !strings.Contains(m.status, "Neon") {
		t.Errorf("the console says %q, want it to name the theme it refused", m.status)
	}
	if !strings.Contains(m.status, "not one this build has") {
		t.Errorf("the console says %q, want it to say why", m.status)
	}
	// And the chef sees it: the refusal is on the screen, not in the console's
	// own head.
	if out := plain(m); !strings.Contains(out, "Neon") {
		t.Errorf("the console does not show the refusal:\n%s", out)
	}
}

// TestAnUnknownOrderIsRefusedByName covers the refusal of an order.
//
// And the theme is still read: one unknown value does not cost the chef the rest.
func TestAnUnknownOrderIsRefusedByName(t *testing.T) {
	store := settingsStore(t)
	writeRecord(t, store, settings.View{Theme: "Dracula", Order: "sideways", ShowServed: true})

	m := loadModel(t, store)
	if m.order != OrderTriage {
		t.Errorf("the order is %v, want the default triage", m.order)
	}
	if m.theme.name != "Dracula" {
		t.Errorf("the theme is %q, want Dracula", m.theme.name)
	}
	if !strings.Contains(m.status, "sideways") {
		t.Errorf("the console says %q, want it to name the order it refused", m.status)
	}
	if !strings.Contains(m.status, "not one this build has") {
		t.Errorf("the console says %q, want it to say why", m.status)
	}
}

// TestTheConsoleNeverWritesAValueItDoesNotKnow covers the write.
//
// The console refuses before it writes, the way the domain's writer refuses before
// it writes — so a value it does not know never reaches the file.
func TestTheConsoleNeverWritesAValueItDoesNotKnow(t *testing.T) {
	store := settingsStore(t)
	m := loadModel(t, store)

	m.theme = newTheme(palette{name: "Neon"})
	m.saveSettings()

	if _, err := os.Stat(settings.Path(store.Root)); err == nil {
		t.Fatal("the record should not have been written at all")
	}
	rec, err := settings.Read(settings.Path(store.Root))
	if err != nil {
		t.Fatal(err)
	}
	if rec.View.Theme == "Neon" {
		t.Fatal("a theme the console does not know should not be written")
	}
	if !strings.Contains(m.status, "Neon") {
		t.Errorf("the console says %q, want it to name the value it refused", m.status)
	}
}
