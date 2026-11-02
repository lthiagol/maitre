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

	cmd := m.switchTab(TabSousChef)
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

// TestTheClickOnTheTabAsksTheRuntime covers the one switch.
//
// The **mouse's** switch is the keyboard's switch: a click on the *Sous-Chef*
// label asks the runtime for the crew, and so the tab never says the automation
// is not running while it is — and never offers to start a second crew.
func TestTheClickOnTheTabAsksTheRuntime(t *testing.T) {
	m := sousModel(t)
	m.runtime = crewRuntimeFixture(t, m)
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
	msg, ok := cmd().(crewMsg)
	if !ok {
		t.Fatalf("the click answered %T, want a crew message", msg)
	}
	if msg.err != nil {
		t.Fatalf("read the crew: %v", msg.err)
	}
	m.Update(msg)

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
}

// TestTheFilterLineSwitchAsksTheRuntime covers the filter line's switch.
//
// The filter line's `tab` is the keyboard's `tab`: it goes through the same path,
// and the same answer.
func TestTheFilterLineSwitchAsksTheRuntime(t *testing.T) {
	m := sousModel(t)
	m.runtime = crewRuntimeFixture(t, m)
	m.tab = TabLog
	m.filterOpen = true

	// `shift+tab` from the Log tab is the Sous-Chef one.
	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	m = updated.(*Model)
	if cmd == nil {
		t.Fatal("the filter line's switch should ask the runtime for the crew")
	}
	msg, ok := cmd().(crewMsg)
	if !ok {
		t.Fatalf("the switch answered %T, want a crew message", msg)
	}
	m.Update(msg)

	if m.tab != TabSousChef {
		t.Fatalf("the filter line left the console on %q, want Sous-Chef", m.tab.Label())
	}
	out := plain(m)
	if strings.Contains(out, "space start the crew") {
		t.Errorf("the footer must not offer to start a second crew:\n%s", out)
	}
}

// TestThePollAsksTheRuntimeWhileTheTabIsVisible covers the poll.
//
// The crew is runtime state: a plan that does not move is exactly when the
// crew's state is most likely to have moved. So while the tab is visible the
// console asks — **even when the plan has not changed**.
func TestThePollAsksTheRuntimeWhileTheTabIsVisible(t *testing.T) {
	m := sousModel(t)
	m.runtime = crewRuntimeFixture(t, m)
	// The plan reads as unchanged.
	fp, err := m.store.DiskFingerprint()
	if err != nil {
		t.Fatalf("fingerprint: %v", err)
	}
	m.disk = fp

	_, cmd := m.Update(pollMsg{})
	crew, ok := crewFrom(cmd)
	if !ok {
		t.Fatal("the poll should ask the runtime for the crew")
	}
	if crew.err != nil {
		t.Fatalf("read the crew: %v", crew.err)
	}
	if len(crew.rows) != 3 {
		t.Errorf("the crew has %d members, want 3", len(crew.rows))
	}
}

// TestThePollDoesNotAskForATabThatIsNotOnScreen covers the rest.
//
// The console asks about the tab the chef is looking at. A tab that is not on
// screen does not spend a call on the runtime.
func TestThePollDoesNotAskForATabThatIsNotOnScreen(t *testing.T) {
	m := sousModel(t)
	m.runtime = crewRuntimeFixture(t, m)
	m.tab = TabLog
	fp, err := m.store.DiskFingerprint()
	if err != nil {
		t.Fatalf("fingerprint: %v", err)
	}
	m.disk = fp

	_, cmd := m.Update(pollMsg{})
	if _, ok := crewFrom(cmd); ok {
		t.Error("the poll should not ask the runtime for a tab that is not on screen")
	}
}

// crewFrom runs the command a message answered and returns the crew message it
// carries, if there is one.
//
// The console asks the runtime and re-arms the poll in one batch, so the crew
// message may be inside a batch or be the answer on its own.
func crewFrom(cmd tea.Cmd) (crewMsg, bool) {
	if cmd == nil {
		return crewMsg{}, false
	}
	msg := cmd()
	batch, ok := msg.(tea.BatchMsg)
	if !ok {
		crew, ok := msg.(crewMsg)
		return crew, ok
	}
	for _, c := range batch {
		if c == nil {
			continue
		}
		if crew, ok := c().(crewMsg); ok {
			return crew, true
		}
	}
	return crewMsg{}, false
}

// TestTheProgressComesFromTheOrchestratorsLabel covers the carrier.
//
// The **orchestrator** owns the progress: it is the one dispatching the crew, so
// it is the one that knows the round and the interaction count. It writes them
// into its own tab's label, and the console reads them from there.
func TestTheProgressComesFromTheOrchestratorsLabel(t *testing.T) {
	m := sousModel(t)
	runtime := crewRuntimeFixture(t, m)
	runtime.knownTabs[0] = herdrTab{
		TabID: "w1:t1",
		Label: "maitre: orchestrator · round 2 · 6 interactions",
	}
	m.runtime = runtime

	rows, err := m.readCrew()
	if err != nil {
		t.Fatalf("read the crew: %v", err)
	}
	m.crew = rows

	progress, ok := m.reportedProgress()
	if !ok {
		t.Fatal("the orchestrator's label carries a progress, and the console did not read it")
	}
	if !progress.Round.Reported || progress.Round.Value != 2 {
		t.Errorf("the round is %+v, want 2", progress.Round)
	}
	if !progress.Interactions.Reported || progress.Interactions.Value != 6 {
		t.Errorf("the interactions are %+v, want 6", progress.Interactions)
	}

	out := plain(m)
	if !strings.Contains(out, "round 2 of 3") {
		t.Errorf("the stamp should carry the round against the cap:\n%s", out)
	}
	if !strings.Contains(out, "6 of 12 interactions") {
		t.Errorf("the stamp should carry the interactions against the cap:\n%s", out)
	}
}

// TestALabelWithNoProgressShowsNone covers the absence.
//
// A label that carries no progress is **not an error**: the crew has not
// reported one, and the stamp prints no progress at all.
func TestALabelWithNoProgressShowsNone(t *testing.T) {
	m := sousModel(t)
	runtime := crewRuntimeFixture(t, m)
	// The orchestrator has reported nothing.
	runtime.knownTabs[0] = herdrTab{TabID: "w1:t1", Label: "maitre: orchestrator"}
	m.runtime = runtime

	rows, err := m.readCrew()
	if err != nil {
		t.Fatalf("read the crew: %v", err)
	}
	m.crew = rows

	if _, ok := m.reportedProgress(); ok {
		t.Error("a label with no progress should carry no progress")
	}
	stamp := stampLine(plain(m))
	if strings.Contains(stamp, "round") || strings.Contains(stamp, "interactions") {
		t.Errorf("the stamp should print no progress it cannot compute: %q", stamp)
	}
	if !strings.Contains(stamp, "plan 7") {
		t.Errorf("the stamp should still carry the vital signs: %q", stamp)
	}
}

// TestTheStartupTaskTellsTheOrchestratorToKeepItsLabelCurrent covers the writer.
//
// The skill carries the role; the startup task carries the work — and it says how
// the orchestrator reports the progress the console reads.
func TestTheStartupTaskTellsTheOrchestratorToKeepItsLabelCurrent(t *testing.T) {
	m := sousModel(t)
	task := m.startupTask()

	for _, want := range []string{"herdr tab rename", "herdr pane current", "maitre: orchestrator", "round <N>", "<M> interactions"} {
		if !strings.Contains(task, want) {
			t.Errorf("the startup task should carry %q: %q", want, task)
		}
	}
}

// TestTheLabelThatIsNotAProgressIsNotOne covers the parser.
//
// The console never guesses: a part it does not recognise carries no progress,
// and the rest of the tail is what the member is on.
func TestTheLabelThatIsNotAProgressIsNotOne(t *testing.T) {
	progress, on := parseProgress("RCP03")
	if progress.Round.Reported || progress.Interactions.Reported {
		t.Errorf("the progress is %+v, want none", progress)
	}
	if on != "RCP03" {
		t.Errorf("the member is on %q, want RCP03", on)
	}

	// A round that is not a number is not a round.
	progress, _ = parseProgress("round two")
	if progress.Round.Reported {
		t.Errorf("the round is %+v, want none — the console does not read a number it cannot compute", progress.Round)
	}
}

// stampLine is the stamp's line in the console's plain text.
//
// The stamp is one line, and the order above it carries the **policy** — the
// budget's cap, which is a different fact. A test that looked at the whole
// screen would read the policy as the progress.
func stampLine(out string) string {
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "THE STAMP") {
			return line
		}
	}
	return ""
}

// TestTheStampPrintsNoNumberItCannotCompute covers the stamp's rule.
//
// A progress past the cap is printed **as it was reported**, with the budget
// spent beside it. "round 7 of 3" says the policy is not a policy, and the
// stamp does not say that.
func TestTheStampPrintsNoNumberItCannotCompute(t *testing.T) {
	m := sousModel(t)
	m.crew = []crewRow{{
		State: "working",
		Job:   "orchestrator",
		Progress: &crewProgress{
			Round:        crewCount{Value: 7, Reported: true},
			Interactions: crewCount{Value: 39, Reported: true},
		},
	}}

	stamp := stampLine(plain(m))
	if strings.Contains(stamp, "7 of 3") || strings.Contains(stamp, "39 of 12") {
		t.Errorf("the stamp should not print a progress against a cap it has passed: %q", stamp)
	}
	if !strings.Contains(stamp, "round 7") {
		t.Errorf("the stamp should carry the round as it was reported: %q", stamp)
	}
	if !strings.Contains(stamp, "the rounds and the interactions are spent") {
		t.Errorf("the stamp should say both allowances are spent: %q", stamp)
	}
}

// TestTheStampMeasuresAProgressInsideTheCap covers the other side.
//
// Inside the budget the progress is measured against the cap — the policy above,
// the progress below.
func TestTheStampMeasuresAProgressInsideTheCap(t *testing.T) {
	m := sousModel(t)
	m.crew = []crewRow{{
		State: "working",
		Job:   "orchestrator",
		Progress: &crewProgress{
			Round:        crewCount{Value: 2, Reported: true},
			Interactions: crewCount{Value: 6, Reported: true},
		},
	}}

	stamp := stampLine(plain(m))
	if !strings.Contains(stamp, "round 2 of 3") {
		t.Errorf("the stamp should measure the round against the cap: %q", stamp)
	}
	if !strings.Contains(stamp, "6 of 12 interactions") {
		t.Errorf("the stamp should measure the interactions against the cap: %q", stamp)
	}
	if strings.Contains(stamp, "the budget is spent") {
		t.Errorf("the stamp should not say the budget is spent inside it: %q", stamp)
	}
}

// TestTheCrewLabelIsSanitisedOnTheSousChefTab covers the label's way in.
//
// The label is **agent-written text**: the console writes the first one, and the
// orchestrator writes it from then on. A newline in it would put a line of the
// agent's text where the crew's rows are a table, and an escape sequence would
// reach the terminal. So it is sanitised where it enters — in `readCrew` — and
// the tab that shows the crew is the tab that proves it.
func TestTheCrewLabelIsSanitisedOnTheSousChefTab(t *testing.T) {
	m := sousModel(t)
	runtime := crewRuntimeFixture(t, m)
	runtime.knownTabs[0] = herdrTab{
		TabID: "w1:t1",
		Label: "maitre: orchestrator · round 2 · 6 interactions\nINJECTED LINE \x1b[31m",
	}
	m.runtime = runtime

	rows, err := m.readCrew()
	if err != nil {
		t.Fatalf("read the crew: %v", err)
	}
	m.crew = rows
	m.width = 100
	m.height = 40

	content := m.View().Content
	checkNoControlCharacters(t, "the Sous-Chef tab", content)
	if strings.Contains(content, "\x1b[31m") {
		t.Errorf("the label's escape reached the render:\n%s", content)
	}
	// The injected text stays on the crew's row: it does not become a line of
	// its own, where the crew's rows are a table.
	for _, line := range strings.Split(stripANSI(content), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "INJECTED") {
			t.Errorf("the label's line broke the crew's table:\n%s", content)
		}
	}
}

// TestTheStampReadsInAllThreeShapes covers the stamp's three shapes.
//
// A progress reported, **no** progress, and a progress **over** the cap — and
// the heading never runs into the vital signs: they are two facts, and they
// read as two.
func TestTheStampReadsInAllThreeShapes(t *testing.T) {
	m := sousModel(t)

	// One: a progress reported.
	m.crew = []crewRow{{
		State: "working",
		Job:   "orchestrator",
		Progress: &crewProgress{
			Round:        crewCount{Value: 2, Reported: true},
			Interactions: crewCount{Value: 6, Reported: true},
		},
	}}
	stamp := stampLine(plain(m))
	if !strings.Contains(stamp, "THE STAMP · plan 7") &&
		!strings.Contains(stamp, "THE STAMP · round 2 of 3") {
		t.Errorf("the heading runs into the vital signs: %q", stamp)
	}
	if !strings.Contains(stamp, "round 2 of 3 · 6 of 12 interactions") {
		t.Errorf("the stamp should carry the progress: %q", stamp)
	}
	if !strings.Contains(stamp, "plan 7 · build 7 · they agree") {
		t.Errorf("the stamp should carry the vital signs: %q", stamp)
	}

	// Two: no progress at all.
	m.crew = []crewRow{{State: "working", Job: "orchestrator"}}
	stamp = stampLine(plain(m))
	if !strings.Contains(stamp, "THE STAMP · plan 7 · build 7 · they agree") {
		t.Errorf("the heading runs into the vital signs: %q", stamp)
	}
	if strings.Contains(stamp, "round") || strings.Contains(stamp, "interactions") {
		t.Errorf("the stamp should print no progress it cannot compute: %q", stamp)
	}

	// Three: a progress over the cap.
	m.crew = []crewRow{{
		State: "working",
		Job:   "orchestrator",
		Progress: &crewProgress{
			Round:        crewCount{Value: 7, Reported: true},
			Interactions: crewCount{Value: 39, Reported: true},
		},
	}}
	stamp = stampLine(plain(m))
	if strings.Contains(stamp, "7 of 3") || strings.Contains(stamp, "39 of 12") {
		t.Errorf("the stamp should not print a progress against a cap it has passed: %q", stamp)
	}
	if !strings.Contains(stamp, "the rounds and the interactions are spent") {
		t.Errorf("the stamp should say both allowances are spent: %q", stamp)
	}
}

// TestTheBudgetIsSpentWhenTheReportReachesTheCap covers the flag's timing.
//
// The flag arrives **with** the report: a crew that has used its third round of
// three has spent its budget, and the stamp says so then — not one round later.
func TestTheBudgetIsSpentWhenTheReportReachesTheCap(t *testing.T) {
	m := sousModel(t)
	m.crew = []crewRow{{
		State: "working",
		Job:   "orchestrator",
		Progress: &crewProgress{
			Round:        crewCount{Value: 3, Reported: true},
			Interactions: crewCount{Value: 12, Reported: true},
		},
	}}

	stamp := stampLine(plain(m))
	if !strings.Contains(stamp, "round 3 of 3 · 12 of 12 interactions") {
		t.Errorf("the stamp should carry the report: %q", stamp)
	}
	if !strings.Contains(stamp, "the rounds and the interactions are spent") {
		t.Errorf("the flag should arrive with the report, and name both allowances: %q", stamp)
	}
}

// TestTheFlagNamesTheAllowanceThatRanOut covers the flag's words.
//
// The budget is **two** allowances — the rounds and the interactions — and a
// reader must not have to guess which of them ran out.
func TestTheFlagNamesTheAllowanceThatRanOut(t *testing.T) {
	m := sousModel(t)

	// The rounds are spent, and the interactions are not.
	m.crew = []crewRow{{
		State: "working",
		Job:   "orchestrator",
		Progress: &crewProgress{
			Round:        crewCount{Value: 3, Reported: true},
			Interactions: crewCount{Value: 6, Reported: true},
		},
	}}
	stamp := stampLine(plain(m))
	if !strings.Contains(stamp, "the rounds are spent") {
		t.Errorf("the flag should name the rounds: %q", stamp)
	}
	if strings.Contains(stamp, "the interactions are spent") {
		t.Errorf("the flag should not name the interactions: %q", stamp)
	}

	// The interactions are spent, and the rounds are not.
	m.crew = []crewRow{{
		State: "working",
		Job:   "orchestrator",
		Progress: &crewProgress{
			Round:        crewCount{Value: 2, Reported: true},
			Interactions: crewCount{Value: 12, Reported: true},
		},
	}}
	stamp = stampLine(plain(m))
	if !strings.Contains(stamp, "the interactions are spent") {
		t.Errorf("the flag should name the interactions: %q", stamp)
	}
	if strings.Contains(stamp, "the rounds are spent") {
		t.Errorf("the flag should not name the rounds: %q", stamp)
	}
}

// TestAZeroedReportIsAReport covers the zero.
//
// "round 0 · 0 interactions" is the crew **reporting** that it has not started a
// round — not a crew that reported nothing. The stamp must not read one as the
// other.
func TestAZeroedReportIsAReport(t *testing.T) {
	m := sousModel(t)
	runtime := crewRuntimeFixture(t, m)
	runtime.knownTabs[0] = herdrTab{
		TabID: "w1:t1",
		Label: "maitre: orchestrator · round 0 · 0 interactions",
	}
	m.runtime = runtime

	rows, err := m.readCrew()
	if err != nil {
		t.Fatalf("read the crew: %v", err)
	}
	m.crew = rows

	if _, ok := m.reportedProgress(); !ok {
		t.Fatal("a zeroed report is a report, and the console read none")
	}
	stamp := stampLine(plain(m))
	if !strings.Contains(stamp, "round 0 of 3 · 0 of 12 interactions") {
		t.Errorf("the stamp should carry the zeroed report: %q", stamp)
	}

	// And it is not the same as no report at all.
	m.crew = []crewRow{{State: "working", Job: "orchestrator"}}
	quiet := stampLine(plain(m))
	if quiet == stamp {
		t.Errorf("a zeroed report and no report read the same: %q", stamp)
	}
}
