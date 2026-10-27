package ui

import (
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// The console owns the terminal.
func TestTheConsoleUsesTheAltScreen(t *testing.T) {
	m := newModel(t)
	m.height = 30
	v := m.View()
	if !v.AltScreen {
		t.Error("the console should use the alt screen — the shell's scrollback stays clean")
	}
	if v.WindowTitle == "" {
		t.Error("the console should set the window title")
	}
	if v.BackgroundColor == nil {
		t.Error("the console should set a background, so the text colours stay legible")
	}
}

// The screen is filled to the bottom, with the footer pinned there.
func TestTheFooterIsPinnedToTheBottom(t *testing.T) {
	m := newModel(t)
	m.height = 30
	out := StripANSI(m.View().Content)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 30 {
		t.Fatalf("the screen has %d rows, want 30", len(lines))
	}
	if !strings.Contains(lines[29], "q") {
		t.Errorf("the footer should be on the last row, got %q", lines[29])
	}
	if !strings.Contains(lines[29], "move") {
		t.Errorf("the footer should list the keys, got %q", lines[29])
	}
}

// The footer says what `space` will do, which depends on what is selected.
func TestTheFooterIsContextual(t *testing.T) {
	m := newModel(t)
	m.height = 30

	// The first recipe is ready — space signs it off.
	if !strings.Contains(StripANSI(m.View().Content), "space sign off") {
		t.Error("the footer should offer to sign off a ready recipe")
	}

	// Move to RCP04, which is approved. The chef has nothing to do with it — an
	// agent claims the work, not the console.
	for {
		s, _ := m.currentSummary()
		if s.ID == "RCP04" {
			break
		}
		m.Update(tea.KeyPressMsg{Code: 'j'})
	}
	if strings.Contains(StripANSI(m.View().Content), "space start") {
		t.Error("the console must not offer to start work — that is an agent's move")
	}
}

// TestTheConsoleDoesNotImpersonateAnAgent checks the console never records itself as
// the runner. Claiming work is the CLI's job.
func TestTheConsoleDoesNotImpersonateAnAgent(t *testing.T) {
	m := newModel(t)
	m.height = 30

	// RCP04 is approved. Space on it must not claim it.
	for {
		s, _ := m.currentSummary()
		if s.ID == "RCP04" {
			break
		}
		m.Update(tea.KeyPressMsg{Code: 'j'})
	}
	m.Update(tea.KeyPressMsg{Code: ' '})

	for _, r := range m.recipes {
		if r.ID == "RCP04" {
			if r.Lifecycle == core.LifecycleExecuting {
				t.Fatal("the console must not claim a recipe — that is an agent's move")
			}
			if r.Execution.ClaimedBy == "runner" {
				t.Fatal("the console must not record itself as an agent")
			}
		}
	}
}

// The list scrolls to keep the selection visible.
func TestTheListScrollsToTheSelection(t *testing.T) {
	m := newModel(t)
	m.height = 12 // only a few rows of list

	// The selection starts at the top.
	if m.offset != 0 {
		t.Fatalf("offset = %d, want 0", m.offset)
	}

	// Move down past the bottom of the window. The scroll happens on render, so the
	// view has to be drawn for the offset to catch up.
	for i := 0; i < 5; i++ {
		m.Update(tea.KeyPressMsg{Code: 'j'})
	}
	m.View()
	if m.offset == 0 {
		t.Fatal("the list should have scrolled to keep the selection visible")
	}

	// And back up. One step does not move the offset — the selection is still
	// visible — so step until it would leave the window.
	before := m.offset
	for i := 0; i < 6; i++ {
		m.Update(tea.KeyPressMsg{Code: 'k'})
	}
	m.View()
	if m.offset >= before {
		t.Fatalf("offset = %d, want less than %d", m.offset, before)
	}
}

// Home and end jump to the ends of the list.
func TestHomeAndEnd(t *testing.T) {
	m := newModel(t)
	m.height = 12
	m.Update(tea.KeyPressMsg{Code: 'G'})
	if m.currentSel() == 0 {
		t.Fatal("end should move the selection to the last item")
	}
	m.Update(tea.KeyPressMsg{Code: 'g'})
	if m.currentSel() != 0 {
		t.Fatal("home should move the selection to the first item")
	}
}

// A resize re-lays out.
func TestResizeRelaysOut(t *testing.T) {
	m := newModel(t)
	m.height = 30
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	out := StripANSI(m.View().Content)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 50 {
		t.Fatalf("the screen has %d rows after a resize, want 50", len(lines))
	}
}

// The split appears on a wide terminal.
func TestTheSplitAppearsOnAWideTerminal(t *testing.T) {
	m := newModel(t)
	m.height = 40
	m.Update(tea.WindowSizeMsg{Width: 180, Height: 40})
	out := StripANSI(m.View().Content)
	if !strings.Contains(out, "Outcome") {
		t.Error("a wide terminal should show the detail pane beside the list")
	}
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	out = StripANSI(m.View().Content)
	if strings.Contains(out, "Outcome") {
		t.Error("a narrow terminal should not show the detail pane")
	}
}

// The Log tab selects entries, newest first.
func TestTheLogTabSelectsEntries(t *testing.T) {
	m := newModel(t)
	m.tab = TabLog
	m.height = 30
	m.Update(tea.KeyPressMsg{Code: 'j'})
	if m.currentSel() != 1 {
		t.Fatalf("selection = %d, want 1", m.currentSel())
	}
	out := StripANSI(m.View().Content)
	if !strings.Contains(out, "recent activity") {
		t.Errorf("the log is missing its header\n%s", out)
	}
}

// The Path tab selects entries and the detail pane follows.
func TestThePathTabSelectsEntries(t *testing.T) {
	m := newModel(t)
	m.tab = TabPath
	m.height = 30
	m.Update(tea.KeyPressMsg{Code: 'j'})
	if m.currentSel() != 1 {
		t.Fatalf("selection = %d, want 1", m.currentSel())
	}
	rec, ok := m.currentRecipe()
	if !ok {
		t.Fatal("the path tab should show the selected recipe in the detail pane")
	}
	if rec.ID == "" {
		t.Fatal("no recipe selected")
	}
}

// A summary line says what the tab holds.
func TestEachTabHasASummary(t *testing.T) {
	m := newModel(t)
	m.height = 30
	want := map[Tab]string{
		TabRecipes:  "recipes",
		TabBacklog:  "specials",
		TabPath:     "the path",
		TabLog:      "recent activity",
		TabSettings: "the plan",
	}
	for tab, needle := range want {
		m.tab = tab
		out := StripANSI(m.View().Content)
		if !strings.Contains(out, needle) {
			t.Errorf("the %s tab is missing %q in its summary", tab.Label(), needle)
		}
	}
}
