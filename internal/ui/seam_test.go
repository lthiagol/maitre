package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// TestEnterOpensTheLogEntry covers the Log tab.
//
// The pane already renders the journal entry — `detailLines` has a `case TabLog:` —
// so the key that opens the pane must open it too.
func TestEnterOpensTheLogEntry(t *testing.T) {
	m := newModel(t)
	m.tab = TabLog

	before := plain(m)
	if !strings.Contains(before, "recent activity") {
		t.Fatalf("the Log tab's summary is missing\n%s", before)
	}

	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.mode != ModeDetail {
		t.Fatalf("enter did not open the detail on the Log tab: mode is %v", m.mode)
	}
	after := plain(m)
	// The event's own detail: the data the journal carries. Nothing else on
	// this screen carries a "Data" section.
	if !strings.Contains(after, "Data") {
		t.Errorf("the open detail does not show the event\n%s", after)
	}
}

// TestTheFooterPromisesOnlyWhatItDoes covers the footer's promise.
//
// A key that does nothing is a promise the console cannot keep. The footer offers
// `enter detail` only where there is a detail to open — the Log's entry, the
// backlog item's, a recipe's — and not where there is none: not on the Sous-Chef,
// and not on a tab with nothing selected.
func TestTheFooterPromisesOnlyWhatItDoes(t *testing.T) {
	m := newModel(t)

	// The Log: the pane shows the entry, so the footer offers the key.
	m.tab = TabLog
	if out := plain(m); !strings.Contains(out, "enter detail") {
		t.Errorf("the footer does not offer the Log's detail\n%s", out)
	}

	// The Sous-Chef: there is nothing to open.
	m.tab = TabSousChef
	if out := plain(m); strings.Contains(out, "enter detail") {
		t.Errorf("the footer promises a detail the Sous-Chef does not have\n%s", out)
	}

	// The Recipes tab with nothing selected: nothing to open either.
	empty := emptyModel(t)
	empty.tab = TabRecipes
	if out := plain(empty); strings.Contains(out, "enter detail") {
		t.Errorf("the footer promises a detail with nothing selected\n%s", out)
	}
}

// TestTheKeyAndThePaneAgreeOnEveryTab holds the seam.
//
// The console has **one** detail pane, and the key that opens it reads one answer:
// the footer offers `enter detail` exactly where the key opens something, and the
// key opens something exactly where the footer offers it. Tab by tab.
func TestTheKeyAndThePaneAgreeOnEveryTab(t *testing.T) {
	m := newModel(t)
	for i := range TabNames {
		tab := Tab(i)
		m.tab = tab
		m.mode = ModeSplit
		m.clampSelections()

		promised := strings.Contains(plain(m), "enter detail")
		if promised != m.hasDetail() {
			t.Errorf("%s: the footer promises the detail (%v) and there is one to open (%v)",
				tab.Label(), promised, m.hasDetail())
			continue
		}
		if !promised {
			continue
		}
		// And the key does what the footer says.
		m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
		if m.mode != ModeDetail {
			t.Errorf("%s: the footer promises the detail and the key does not open it", tab.Label())
		}
	}
}

// emptyModel builds a console over a plan with nothing in it.
func emptyModel(t *testing.T) *Model {
	t.Helper()
	store := core.Open(filepath.Join(t.TempDir(), core.PlanDirName))
	if err := store.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	return loadModel(t, store)
}
