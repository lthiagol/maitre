package ui

import (
	"errors"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"
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

// TestSpaceStartsTheCrew covers the start.
//
// `space` starts the crew: the console opens a tab on the plan's root and runs
// the orchestrator in it.
func TestSpaceStartsTheCrew(t *testing.T) {
	m := sousModel(t)
	runtime := crewRuntimeFixture(t, m)
	m.runtime = runtime

	m.Update(tea.KeyPressMsg{Code: ' '})

	if len(runtime.started) != 1 {
		t.Fatalf("the runtime was asked to start %d times, want 1", len(runtime.started))
	}
	if !strings.Contains(runtime.started[0], "opencode") {
		t.Errorf("the crew should run the harness the order names: %q", runtime.started[0])
	}
}

// TestTheStartupTaskCarriesTheRoleAndTheWork covers the task.
//
// **The skill carries the role; the startup task carries the work.** The task
// names the plan, tells the orchestrator to load its role, and says what it is
// to execute.
func TestTheStartupTaskCarriesTheRoleAndTheWork(t *testing.T) {
	m := sousModel(t)
	task := m.startupTask()

	if !strings.Contains(task, "maitre-cli skill load orchestrator-role") {
		t.Errorf("the task should load the role: %q", task)
	}
	if !strings.Contains(task, m.store.ProjectRoot()) {
		t.Errorf("the task should name the plan's root: %q", task)
	}
	work := m.sousWork()
	if work == "none — the queue is empty" {
		t.Fatalf("the seeded plan should have work: %q", work)
	}
	if !strings.Contains(task, work) {
		t.Errorf("the task should carry the work %q: %q", work, task)
	}
}

// TestTheWorkIsUnderTheBudgetsCap covers the budget.
//
// The budget is the chef's cap on how much the crew may do, and the work handed
// to the orchestrator is the queue under that cap.
func TestTheWorkIsUnderTheBudgetsCap(t *testing.T) {
	m := sousModel(t)
	m.picks.Budget = 0 // the smallest budget: 2 rounds
	m.picks.Sources = "recipes"

	queue := m.sousQueue("recipes")
	if len(queue) < 3 {
		t.Skipf("the plan has %d recipes, so the cap is not tested", len(queue))
	}
	work := m.sousWork()
	if got := len(strings.Split(work, ", ")); got != 2 {
		t.Errorf("the work is %q — %d items, want 2", work, got)
	}
}

// TestSpaceStopsTheCrew covers the stop.
//
// `space` stops the crew when it is running: the console closes the crew's own
// tabs, and only those.
func TestSpaceStopsTheCrew(t *testing.T) {
	m := sousModel(t)
	runtime := crewRuntimeFixture(t, m)
	m.runtime = runtime
	rows, err := m.readCrew()
	if err != nil {
		t.Fatalf("read the crew: %v", err)
	}
	m.crew = rows

	m.Update(tea.KeyPressMsg{Code: ' '})

	if len(runtime.stopped) != 3 {
		t.Fatalf("the runtime closed %d tabs, want 3", len(runtime.stopped))
	}
	for _, id := range runtime.stopped {
		if id == "w9:t1" {
			t.Errorf("the console closed a tab that is not the crew's: %v", runtime.stopped)
		}
	}
	if len(runtime.started) != 0 {
		t.Errorf("the console started a crew while one was running: %v", runtime.started)
	}
}

// TestTheFooterNamesWhichTheSpaceDoes covers the footer's promise.
//
// One key, two meanings: the footer says which one it is about to do.
func TestTheFooterNamesWhichTheSpaceDoes(t *testing.T) {
	m := sousModel(t)
	m.runtime = crewRuntimeFixture(t, m)

	if out := plain(m); !strings.Contains(out, "space start the crew") {
		t.Errorf("the footer should offer to start the crew:\n%s", out)
	}

	rows, err := m.readCrew()
	if err != nil {
		t.Fatalf("read the crew: %v", err)
	}
	m.crew = rows
	if out := plain(m); !strings.Contains(out, "space stop the crew") {
		t.Errorf("the footer should offer to stop the crew:\n%s", out)
	}
}

// TestTheCrewIsReadWhenTheTabComesOnScreen covers the read.
//
// The crew is runtime state: the console asks for it when the tab comes on
// screen, and the answer comes back as a message.
func TestTheCrewIsReadWhenTheTabComesOnScreen(t *testing.T) {
	m := sousModel(t)
	m.runtime = crewRuntimeFixture(t, m)
	m.tab = TabLog

	cmd := m.switchTabCmd(TabSousChef)
	if cmd == nil {
		t.Fatal("the tab should ask for the crew")
	}
	msg, ok := cmd().(crewMsg)
	if !ok {
		t.Fatalf("the command answered %T, want a crew message", msg)
	}
	if msg.err != nil {
		t.Fatalf("read the crew: %v", msg.err)
	}
	if len(msg.rows) != 3 {
		t.Errorf("the crew has %d members, want 3", len(msg.rows))
	}

	// And the message puts them on the tab.
	m.Update(msg)
	if len(m.crew) != 3 {
		t.Errorf("the console holds %d members, want 3", len(m.crew))
	}
}
