package ui

import (
	"strings"
	"testing"

	"charm.land/bubbletea/v2"
)

// The console still asks.
//
// The render path stopped asking, and the **console** did not: the poll asks the
// runtime while the tab is visible, and the tab shows what it last heard. The
// witnesses are process level — a fake herdr on the PATH counts the forks —
// because the claim is about processes, not about calls.

// TestThePollAsksTheRuntimeInOneProcess covers the poll.
//
// The poll asks the runtime while the tab is visible, and it asks **once**: the
// crew read is one process.
//
// The poll's **frequency** is not this slice's — it is a second, and it stays a
// second. What changed is what the poll asks for, not how often.
func TestThePollAsksTheRuntimeInOneProcess(t *testing.T) {
	store := goldenStore(t)
	herdr := herdrOnPath(t, herdrAnswer(store.ProjectRoot()))
	t.Setenv("PATH", herdr.dir)

	m := loadModel(t, store)
	m.tab = TabSousChef

	_, cmd := m.Update(pollMsg{})
	crew, ok := crewFrom(cmd)
	if !ok {
		t.Fatal("the poll should ask the runtime for the crew")
	}
	if crew.err != nil {
		t.Fatalf("read the crew: %v", crew.err)
	}
	if len(crew.rows) != 1 {
		t.Errorf("the crew has %d members, want 1:\n%+v", len(crew.rows), crew.rows)
	}

	calls := herdr.calls(t)
	if len(calls) != 1 {
		t.Errorf("the poll forked %d times, want 1: %v", len(calls), calls)
	}
}

// TestTheTabShowsWhatItLastHeard covers the answer the tab holds.
//
// The poll asks, the answer arrives as a message, and the tab shows it. The
// **frame** does not ask: a frame that asked would ask on every frame, and a
// frame never waits on a process.
func TestTheTabShowsWhatItLastHeard(t *testing.T) {
	store := goldenStore(t)
	herdr := herdrOnPath(t, herdrAnswer(store.ProjectRoot()))
	t.Setenv("PATH", herdr.dir)

	m := loadModel(t, store)
	m.tab = TabSousChef

	cmd := m.refreshCrew()
	m.Update(cmd().(crewMsg))

	out := plain(m)
	if !strings.Contains(out, "orchestrator") {
		t.Errorf("the tab should show the crew it last heard:\n%s", out)
	}

	// And the frames that follow ask nothing.
	for i := 0; i < 5; i++ {
		_ = plain(m)
	}
	calls := herdr.calls(t)
	if len(calls) != 1 {
		t.Errorf("the frames forked %d times, want 1 — the ask: %v", len(calls), calls)
	}
}

// TestTheClickOnTheTabReachesTheRuntime covers F-140's witness.
//
// The **mouse's** switch is the keyboard's switch: a click on the *Sous-Chef*
// label asks the runtime for the crew, and so the tab never says the automation
// is not running while it is — and never offers to start a second crew.
func TestTheClickOnTheTabReachesTheRuntime(t *testing.T) {
	store := goldenStore(t)
	herdr := herdrOnPath(t, herdrAnswer(store.ProjectRoot()))
	t.Setenv("PATH", herdr.dir)

	m := loadModel(t, store)
	m.tab = TabLog
	m.width = 120
	m.height = 30
	// The render leaves the map behind, and the map says where the labels are.
	m.View()

	box := m.hits.tabs[int(TabSousChef)]
	updated, cmd := m.Update(tea.MouseClickMsg{
		X: box.x0, Y: box.y0, Button: tea.MouseLeft,
	})
	m = updated.(*Model)
	if cmd == nil {
		t.Fatal("a click on the tab should ask the runtime for the crew")
	}
	m.Update(cmd().(crewMsg))

	if m.tab != TabSousChef {
		t.Fatalf("the click left the console on %q, want Sous-Chef", m.tab.Label())
	}
	out := plain(m)
	if strings.Contains(out, "the automation is not running") {
		t.Errorf("the tab says the automation is not running while it is:\n%s", out)
	}
	if !strings.Contains(out, "space stop the crew") {
		t.Errorf("the footer should offer to stop the crew:\n%s", out)
	}
	if strings.Contains(out, "space start the crew") {
		t.Errorf("the footer must not offer to start a second crew:\n%s", out)
	}

	calls := herdr.calls(t)
	if len(calls) != 1 {
		t.Errorf("the click forked %d times, want 1: %v", len(calls), calls)
	}
}
