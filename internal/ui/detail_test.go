package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/kitchen-brigade/internal/core"
)

// TestTheConsoleShowsADeviation covers the promise the console makes.
//
// The detail pane renders "deviated: …" under a step, and no operation ever wrote
// it: a rendered field with no writer is a promise the code does not keep.
func TestTheConsoleShowsADeviation(t *testing.T) {
	root := filepath.Join(t.TempDir(), core.PlanDirName)
	store := core.Open(root)
	if err := core.Seed(store); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "line-cook", Session: "s-1", Role: "runner"})
	if _, err := ops.RecordDeviation("R01", "G1", "read the queue twice"); err != nil {
		t.Fatal(err)
	}
	// The console reads the plan after the deviation was recorded.
	m := loadModel(t, store)
	for {
		s, ok := m.currentSummary()
		if ok && s.ID == "R01" {
			break
		}
		m.Update(tea.KeyPressMsg{Code: 'j'})
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	out := plain(m)
	if !strings.Contains(out, "deviated: read the queue twice") {
		t.Fatalf("the detail should show the deviation:\n%s", out)
	}
}

// TestTheLogSelectionMatchesTheList covers the Log tab.
//
// The selection was clamped to the whole journal while the list shows only the
// newest events that fit, so the selection could point at an event that was not on
// the screen — and the detail pane showed it anyway.
func TestTheLogSelectionMatchesTheList(t *testing.T) {
	m := newModel(t)
	m.tab = TabLog
	m.height = 30
	m.Update(tea.KeyPressMsg{Code: 'G'})

	window := m.logWindow()
	if len(window) >= len(m.journal.Events) {
		t.Fatalf("the journal has %d events and the window %d — the test needs a journal longer than the window",
			len(m.journal.Events), len(window))
	}
	if m.selLog != len(window)-1 {
		t.Fatalf("selection = %d, want %d (the last visible event)", m.selLog, len(window)-1)
	}
	out := strings.Join(m.logDetail(100), "\n")
	if !strings.Contains(out, window[len(window)-1].At) {
		t.Fatalf("the detail should show the last visible event:\n%s", out)
	}
}
