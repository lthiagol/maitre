package ui

import (
	"errors"
	"strings"
	"testing"
)

// fakeRuntime is a runtime that answers from a fixture.
//
// The console's own tests must not depend on a binary being on the PATH: the
// runtime is a seam, and this is what sits behind it.
type fakeRuntime struct {
	// here is whether the runtime is on the PATH.
	here bool
	// knownAgents and knownTabs are what the runtime knows.
	knownAgents []herdrAgent
	knownTabs   []herdrTab
	// startErr is what the spawn fails with, when it is to fail.
	startErr error
	// started records the starts the console asked for.
	started []string
	// stopped records the tabs the console closed.
	stopped []string
}

// available reports whether the runtime is here.
func (f *fakeRuntime) available() bool { return f.here }

// agents lists the runtime's agents.
func (f *fakeRuntime) agents() ([]herdrAgent, error) { return f.knownAgents, nil }

// tabs lists the runtime's tabs.
func (f *fakeRuntime) tabs() ([]herdrTab, error) { return f.knownTabs, nil }

// start records the start.
func (f *fakeRuntime) start(planRoot, project, harness, task string) error {
	if f.startErr != nil {
		return f.startErr
	}
	f.started = append(f.started, harness+" "+task)
	return nil
}

// stop records the closed tabs.
func (f *fakeRuntime) stop(tabIDs []string) error {
	f.stopped = append(f.stopped, tabIDs...)
	return nil
}

// crewRuntimeFixture builds a runtime holding a crew on the given plan.
//
// The crew is three members on the plan's root, and one agent elsewhere: the
// console finds the crew by the **agents' cwd**, so the agent on another plan is
// not a member.
func crewRuntimeFixture(t *testing.T, m *Model) *fakeRuntime {
	t.Helper()
	root := m.store.ProjectRoot()
	return &fakeRuntime{
		here: true,
		knownAgents: []herdrAgent{
			{PaneID: "w1:p1", TabID: "w1:t1", Cwd: root, Status: "working"},
			{PaneID: "w1:p2", TabID: "w1:t2", Cwd: root, Status: "idle"},
			{PaneID: "w1:p3", TabID: "w1:t3", Cwd: root, Status: "blocked"},
			{PaneID: "w9:p1", TabID: "w9:t1", Cwd: "/somewhere/else", Status: "working"},
		},
		knownTabs: []herdrTab{
			{TabID: "w1:t1", Label: "maitre: orchestrator"},
			{TabID: "w1:t2", Label: "maitre: runner · RCP03"},
			{TabID: "w1:t3", Label: "maitre: reviewer"},
			{TabID: "w9:t1", Label: "maitre: runner · RCP09"},
		},
	}
}

// TestTheCrewIsFoundByTheAgentsCwd covers the anchor.
//
// The console finds the crew by the **agents' cwd**: a workspace's own listing
// carries no path, so the agents' cwd is what says an agent works *this* plan.
func TestTheCrewIsFoundByTheAgentsCwd(t *testing.T) {
	m := sousModel(t)
	runtime := crewRuntimeFixture(t, m)
	m.runtime = runtime

	rows, err := m.readCrew()
	if err != nil {
		t.Fatalf("read the crew: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("the crew has %d members, want 3:\n%+v", len(rows), rows)
	}
	for _, row := range rows {
		if row.Job == "runner" && row.On != "RCP03" {
			t.Errorf("the runner is on %q, want RCP03", row.On)
		}
	}
}

// TestTheCrewRowShowsTheJobAndTheState covers the row.
//
// The **state** is the runtime's own, and the **job** comes from the tab's
// label: the runtime knows an agent is an mcode, not which job it plays.
func TestTheCrewRowShowsTheJobAndTheState(t *testing.T) {
	m := sousModel(t)
	m.runtime = crewRuntimeFixture(t, m)

	rows, err := m.readCrew()
	if err != nil {
		t.Fatalf("read the crew: %v", err)
	}
	m.crew = rows
	out := plain(m)

	for _, want := range []string{"orchestrator", "runner", "reviewer", "working", "idle", "RCP03"} {
		if !strings.Contains(out, want) {
			t.Errorf("the stamp is missing %q:\n%s", want, out)
		}
	}
}

// TestATabThatIsNotOnTheRosterIsNotTheCrew covers the roster.
//
// The console finds the agents by their cwd and reads the job from the tab's
// label. A tab whose label does not carry the roster's mark is not a member of
// the crew — the chef's own pane is not the crew's to manage.
func TestATabThatIsNotOnTheRosterIsNotTheCrew(t *testing.T) {
	m := sousModel(t)
	runtime := crewRuntimeFixture(t, m)
	// The chef's own pane: on the plan's root, and not on the roster.
	runtime.knownAgents = append(runtime.knownAgents, herdrAgent{
		PaneID: "w1:p9", TabID: "w1:t9", Cwd: m.store.ProjectRoot(), Status: "working",
	})
	runtime.knownTabs = append(runtime.knownTabs, herdrTab{TabID: "w1:t9", Label: "maitre-cli"})
	m.runtime = runtime

	rows, err := m.readCrew()
	if err != nil {
		t.Fatalf("read the crew: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("the crew has %d members, want 3:\n%+v", len(rows), rows)
	}
	for _, row := range rows {
		if row.Job == "maitre-cli" {
			t.Errorf("a tab that is not on the roster should not be a member:\n%+v", rows)
		}
	}
}

// TestTheBlockedAgentAsksTheOrchestrator covers the blocked row.
//
// A blocked agent raises a message to the orchestrator, asking for directions.
// The console gains **no move** for it: the stamp shows that the reviewer asked
// and that the orchestrator has it.
func TestTheBlockedAgentAsksTheOrchestrator(t *testing.T) {
	m := sousModel(t)
	m.runtime = crewRuntimeFixture(t, m)

	rows, err := m.readCrew()
	if err != nil {
		t.Fatalf("read the crew: %v", err)
	}
	m.crew = rows
	out := plain(m)

	if !strings.Contains(out, "asked the orchestrator") {
		t.Errorf("the stamp should say the blocked agent asked:\n%s", out)
	}
	if !strings.Contains(out, "the orchestrator has it") {
		t.Errorf("the stamp should say the orchestrator has the question:\n%s", out)
	}
}

// TestTheRuntimeThatIsNotThereIsSaid covers the runtime's absence.
//
// The tab that cannot start the automation must not look like one that can: the
// console says the runtime is not here.
func TestTheRuntimeThatIsNotThereIsSaid(t *testing.T) {
	m := sousModel(t)
	m.runtime = &fakeRuntime{here: false}

	_, err := m.readCrew()
	if err == nil {
		t.Fatal("a runtime that is not there should be an error")
	}
	if !errors.Is(err, errNoRuntime) {
		t.Errorf("the error is %v, want the runtime's absence", err)
	}
}
