package ui

import (
	"errors"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"
)

// The five failures.
//
// Each one has what the tab says and what it offers. A tab that cannot do its
// job says so — a key that cannot work must not be offered, and a tab that is
// waiting must not look like a tab that is failing.

// TestTheRuntimeThatIsNotHereSaysSoAndOffersNothing covers the first failure.
//
// **herdr is not on the PATH**: the tab says "herdr is not here" and offers
// nothing. A tab that cannot start the automation must not look like one that
// can.
func TestTheRuntimeThatIsNotHereSaysSoAndOffersNothing(t *testing.T) {
	m := sousModel(t)
	m.runtime = &fakeRuntime{here: false}

	// The console asks, and the answer is that the runtime is not here.
	cmd := m.refreshCrew()
	msg := cmd().(crewMsg)
	m.Update(msg)

	if m.crewErr == nil {
		t.Fatal("the console should know the runtime is not here")
	}
	out := plain(m)
	if !strings.Contains(out, "herdr is not here") {
		t.Errorf("the tab should say the runtime is not here:\n%s", out)
	}
	if strings.Contains(out, "space start the crew") || strings.Contains(out, "space stop the crew") {
		t.Errorf("the tab should offer no move at all:\n%s", out)
	}
}

// TestAFailedStartSaysSoAndTheMoveStays covers the second failure.
//
// **The spawn fails**: the tab says so, and the move stays. A failed start is
// not a running automation.
func TestAFailedStartSaysSoAndTheMoveStays(t *testing.T) {
	m := sousModel(t)
	runtime := crewRuntimeFixture(t, m)
	runtime.startErr = errors.New("the tab could not be opened")
	m.runtime = runtime

	m.Update(tea.KeyPressMsg{Code: ' '})

	if !strings.Contains(m.status, "the crew did not start") {
		t.Errorf("the console says %q, want it to say the start failed", m.status)
	}
	out := plain(m)
	if !strings.Contains(out, "space start the crew") {
		t.Errorf("the move should stay after a failed start:\n%s", out)
	}
}

// TestTheCrewGoneSaysTheAutomationIsNotRunning covers the third failure.
//
// **The crew is gone** — every agent exited: the tab says the automation is not
// running, and the move comes back.
func TestTheCrewGoneSaysTheAutomationIsNotRunning(t *testing.T) {
	m := sousModel(t)
	m.runtime = &fakeRuntime{here: true}

	cmd := m.refreshCrew()
	m.Update(cmd().(crewMsg))

	out := plain(m)
	if !strings.Contains(out, "the automation is not running") {
		t.Errorf("the tab should say the automation is not running:\n%s", out)
	}
	if !strings.Contains(out, "space start the crew") {
		t.Errorf("the move should come back when the crew is gone:\n%s", out)
	}
}

// TestTheOrchestratorGoneIsSaid covers the fourth failure.
//
// **The orchestrator is gone but the crew is not**: the tab says so. A crew
// with no orchestrator has nobody to ask for directions, and the tab must not
// pretend otherwise.
func TestTheOrchestratorGoneIsSaid(t *testing.T) {
	m := sousModel(t)
	runtime := crewRuntimeFixture(t, m)
	// The orchestrator's tab is not on the roster any more.
	runtime.knownTabs[0] = herdrTab{TabID: "w1:t1", Label: "maitre-cli"}
	m.runtime = runtime

	cmd := m.refreshCrew()
	m.Update(cmd().(crewMsg))

	if len(m.crew) == 0 {
		t.Fatal("the crew should still have members")
	}
	out := plain(m)
	if !strings.Contains(out, "the orchestrator is gone") {
		t.Errorf("the tab should say the orchestrator is gone:\n%s", out)
	}
}

// TestThePlanThatMovedHasNoAutomation covers the fifth failure.
//
// **The plan's root moved**: the console finds the crew by the agents' cwd, so
// a plan that moved is a plan with no automation — the agents are working on
// the old root, and the tab says the automation is not running.
func TestThePlanThatMovedHasNoAutomation(t *testing.T) {
	m := sousModel(t)
	runtime := crewRuntimeFixture(t, m)
	// Every agent works on the plan's **old** root.
	moved := m.store.ProjectRoot() + "-moved"
	for i := range runtime.knownAgents {
		runtime.knownAgents[i].Cwd = moved
	}
	m.runtime = runtime

	cmd := m.refreshCrew()
	m.Update(cmd().(crewMsg))

	if len(m.crew) != 0 {
		t.Fatalf("a plan that moved should have no crew:\n%+v", m.crew)
	}
	out := plain(m)
	if !strings.Contains(out, "the automation is not running") {
		t.Errorf("the tab should say the automation is not running:\n%s", out)
	}
	if !strings.Contains(out, "space start the crew") {
		t.Errorf("the move should come back for a plan that moved:\n%s", out)
	}
}
