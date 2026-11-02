package ui

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/lthiagol/maitre/internal/core"
	"github.com/lthiagol/maitre/internal/settings"
)

// The Sous-Chef tab — the automation's ticket.
//
// The tab is a **ticket**: an order at the top (what the crew is, and what it may
// spend), the queue under it (what it will work) and a stamp at the foot (where
// the work stands, and whether the plan and this build speak the same rules).
//
// **The name is the tab's name, not a job's.** The rows carry the glossary's job
// names — orchestrator, runner, reviewer — so the tab and the glossary's
// sous-chef stay two different things.
//
// The tab holds **no selectable item**: it is a ticket, not a list. The cursor
// over the order's four controls is the tab's own, and it never touches the
// list's selection machinery.

// The order's vocabulary.
//
// The console owns these names — it is what renders them — and it refuses a name
// it does not know, by name. The house's lesson: a safe unknown value disabled a
// whole feature, with green tests.

// harnessNames are the harnesses the crew may run, in the order the tab shows them.
var harnessNames = []string{"opencode", "pi", "custom"}

// harnessCommands is the command each harness runs.
//
// The order says the command, so `custom` is never a mystery: the chef reads what
// will run before it runs.
var harnessCommands = map[string]string{
	"opencode": "opencode",
	"pi":       "pi --model sonnet",
	"custom":   "claude --dangerously-skip-permissions",
}

// sourceNames are the places the crew's work comes from.
var sourceNames = []string{"ready", "recipes", "backlog"}

// crewSizes are the crew's sizes, from the biggest to the smallest.
var crewSizes = []int{3, 2, 1}

// crewConsequence is what each crew size means.
//
// **The count is a composition, not a number.** Three agents, one per role; two,
// and the orchestrator reviews the runner; one, and **the chef reviews** — a lone
// agent cannot review itself, and the close gate refuses a self-reviewed recipe.
var crewConsequence = map[int]string{
	3: "orchestrator, runner, reviewer",
	2: "runner, and the orchestrator reviews",
	1: "one agent, and the chef reviews",
}

// sousBudget is one budget the chef may set: the cap on rounds and the cap on
// interactions.
type sousBudget struct {
	Rounds       int
	Interactions int
}

// sousBudgets are the budgets the order offers.
var sousBudgets = []sousBudget{
	{Rounds: 2, Interactions: 8},
	{Rounds: 3, Interactions: 12},
	{Rounds: 5, Interactions: 24},
}

// sousPicks is the order the chef has set: the four controls of the ticket.
//
// The picks are the console's own copy of what the record holds. They are read
// from the record when the console starts and written back when the chef changes
// one.
type sousPicks struct {
	// Harness is the harness the crew runs.
	Harness string
	// Count is how many agents the crew has.
	Count int
	// Sources is where the work comes from.
	Sources string
	// Budget is the index of the budget in sousBudgets.
	Budget int
}

// defaultSousPicks is the order a plan with no settings file has.
func defaultSousPicks() sousPicks {
	return sousPicks{
		Harness: settings.DefaultHarness,
		Count:   settings.DefaultCrew,
		Sources: settings.DefaultSources,
		Budget:  defaultSousBudget,
	}
}

// defaultSousBudget is the index of the budget a console nobody told otherwise
// uses.
const defaultSousBudget = 1

// sousPicksFrom reads the record's crew section into the picks.
//
// A value the console does not know is refused, by name, and the console's own
// default stands. The refusals come back as sentences: the chef is told which
// value the record carried, because "unknown" alone leaves the chef guessing.
func sousPicksFrom(crew settings.Crew) (sousPicks, []string) {
	picks := defaultSousPicks()
	var said []string

	if knownHarness(crew.Harness) {
		picks.Harness = crew.Harness
	} else {
		said = append(said, unknownName("harness", crew.Harness))
	}

	if knownCrewSize(crew.Count) {
		picks.Count = crew.Count
	} else {
		said = append(said, unknownName("crew size", strconv.Itoa(crew.Count)))
	}

	if knownSource(crew.Sources) {
		picks.Sources = crew.Sources
	} else {
		said = append(said, unknownName("source", crew.Sources))
	}

	if i, ok := sousBudgetIndex(crew.Rounds, crew.Interactions); ok {
		picks.Budget = i
	} else {
		said = append(said, fmt.Sprintf(
			"the budget %d rounds · %d interactions is not one this build has",
			crew.Rounds, crew.Interactions))
	}
	return picks, said
}

// crewRecord writes the picks into the record's crew section.
//
// The round trip is the record's own: the picks the console holds are the picks
// the record holds, and a console that read a record and wrote it back does not
// change it.
func (p sousPicks) crewRecord() settings.Crew {
	budget := sousBudgets[p.Budget]
	return settings.Crew{
		Harness:      p.Harness,
		Count:        p.Count,
		Sources:      p.Sources,
		Rounds:       budget.Rounds,
		Interactions: budget.Interactions,
	}
}

// knownHarness reports whether the console knows that harness.
func knownHarness(name string) bool {
	for _, h := range harnessNames {
		if h == name {
			return true
		}
	}
	return false
}

// knownCrewSize reports whether the console knows that crew size.
func knownCrewSize(n int) bool {
	for _, s := range crewSizes {
		if s == n {
			return true
		}
	}
	return false
}

// knownSource reports whether the console knows that source.
func knownSource(name string) bool {
	for _, s := range sourceNames {
		if s == name {
			return true
		}
	}
	return false
}

// sousBudgetIndex is the index of the budget with those caps, or -1.
func sousBudgetIndex(rounds, interactions int) (int, bool) {
	for i, b := range sousBudgets {
		if b.Rounds == rounds && b.Interactions == interactions {
			return i, true
		}
	}
	return -1, false
}

// harnessCommand is the command the given harness runs.
func harnessCommand(name string) string {
	if cmd, ok := harnessCommands[name]; ok {
		return cmd
	}
	return harnessCommands[settings.DefaultHarness]
}

// sousMoveLabel names what `space` will do on this tab.
//
// The move is **one key with two meanings**: `space` starts the crew when it is
// not running, and stops it when it is. The footer says which — a key that can
// do two things must say which one it is about to do.
func (m *Model) sousMoveLabel() string {
	if m.crewErr != nil {
		// The console does not know whether the crew is running, so it offers no
		// move: a move offered on a guess is a move that can start a second crew.
		return ""
	}
	if m.crewRunning() {
		return "stop the crew"
	}
	return "start the crew"
}

// sousPrimaryAction is the tab's move: start the crew, or stop it.
//
// Starting is the console's first action that is not a plan write: it spawns a
// process. It writes **nothing** to the plan — the crew does the work, and the
// plan's records are the crew's to write.
func (m *Model) sousPrimaryAction() (tea.Model, tea.Cmd) {
	if m.crewErr != nil {
		m.status = m.crewErr.Error()
		return m, nil
	}
	if m.crewRunning() {
		return m, m.stopCrew()
	}
	return m, m.startCrew()
}

// startCrew opens the crew's tab and runs the orchestrator in it.
//
// The console remembers nothing: it does not keep the pane id, and it does not
// write a file. It finds the crew again by asking the runtime, so a console that
// was closed and reopened still shows the running automation.
func (m *Model) startCrew() tea.Cmd {
	err := m.runtime.start(
		m.store.ProjectRoot(),
		m.projectName(),
		harnessCommand(m.picks.Harness),
		m.startupTask(),
	)
	if err != nil {
		// A failed start is not a running automation: the move stays.
		m.status = "the crew did not start: " + err.Error()
		return m.refreshCrew()
	}
	m.status = "the crew is starting: " + m.sousWork()
	return m.refreshCrew()
}

// stopCrew closes the crew's tabs.
//
// **Only the crew's own tabs**: the chef's panes are not the crew's to close.
func (m *Model) stopCrew() tea.Cmd {
	tabs := m.crewTabIDs()
	if len(tabs) == 0 {
		m.status = "the crew is not running"
		return m.refreshCrew()
	}
	if err := m.runtime.stop(tabs); err != nil {
		m.status = "the crew did not stop: " + err.Error()
		return m.refreshCrew()
	}
	m.status = "the crew is stopped"
	return m.refreshCrew()
}

// rosterExample is the console's copy of the format the roster's lines carry.
//
// **The skill's table is the specification** — the roster's table in
// `internal/cli/skills/orchestrator-role/SKILL.md` — and this constant is the
// console's copy of it. The startup task hands it to the orchestrator at spawn,
// and the wire test holds the two together: a drift in either one turns a test
// red.
//
// The parser is the third copy of the format, and it is the one that reads them.
const rosterExample = "maitre: orchestrator · round <N> · <M> interactions"

// startupTask is what the orchestrator is given when the crew starts.
//
// **The skill carries the role; the startup task carries the work.** The task is
// short because the role is not in it: the orchestrator runs
// `maitre-cli skill load orchestrator-role` and follows it.
//
// The task's copy of the format is **an example**, not the specification: it
// follows the skill's table, and the wire test holds it to it.
func (m *Model) startupTask() string {
	return fmt.Sprintf(
		"You are the orchestrator of the plan at %s. Run `maitre-cli skill load orchestrator-role` and follow it. The recipes to execute are: %s. Keep your tab's label current — the console reads the crew's progress from it: `herdr tab rename <TAB_ID> %q` (your own tab's id is in `herdr pane current`).",
		m.store.ProjectRoot(), m.sousWork(), rosterExample)
}

// sousWork is the work the crew is given: the queue's items, up to the budget's
// cap on rounds.
//
// The budget is the chef's cap on how much the crew may do, and the work handed
// to the orchestrator is the queue under that cap.
func (m *Model) sousWork() string {
	queue := m.sousQueue(m.picks.Sources)
	budget := sousBudgets[m.picks.Budget]
	if len(queue) > budget.Rounds {
		queue = queue[:budget.Rounds]
	}
	if len(queue) == 0 {
		return "none — the queue is empty"
	}
	ids := make([]string, 0, len(queue))
	for _, item := range queue {
		ids = append(ids, item.ID)
	}
	return strings.Join(ids, ", ")
}

// crewRow is one member of the crew, as the runtime knows it.
//
// The **state** is the runtime's own — working, idle, blocked, done. The **job**
// is the tab's label, and the **on** is what that label says the member is
// working on.
type crewRow struct {
	State string
	Job   string
	On    string
	// TabID is the tab the member is on. The console keeps it so it can close
	// the crew's own tabs — and only those — when the chef stops the crew.
	TabID string
	// Progress is what the orchestrator reported: the round the crew is on and
	// the interactions it has spent.
	//
	// It is **absent** when the label carries none — and an absent progress is
	// not an error: it means the crew has not reported one.
	Progress *crewProgress
}

// sousItem is one item of the queue.
type sousItem struct {
	ID    string
	Title string
}

// sousQueue is the work the crew will do, in the order it will do it.
//
// The queue is **derived** from the plan and the order's source — never stored,
// never written. The source says where the work comes from:
//
//   - **ready** — the recipes an agent could claim now, and the open small tasks
//   - **recipes** — the plan's open recipes
//   - **backlog** — the backlog's open items, the ideas first
//
// The domain's own views answer "open" and "claimable": the console counts no
// lifecycle of its own.
func (m *Model) sousQueue(source string) []sousItem {
	switch source {
	case "recipes":
		var out []sousItem
		for _, r := range m.recipes {
			if core.TriageOf(r) == core.TriageServed {
				continue
			}
			out = append(out, sousItem{ID: r.ID, Title: r.Title})
		}
		return out
	case "backlog":
		ideas, tasks := core.BacklogGroups(m.items)
		var out []sousItem
		for _, group := range [][]core.BacklogItem{ideas, tasks} {
			for _, it := range group {
				if it.Status != core.StatusOpen {
					continue
				}
				out = append(out, sousItem{ID: it.ID, Title: it.Title})
			}
		}
		return out
	default:
		var out []sousItem
		for _, r := range core.Adoptable(m.recipes) {
			out = append(out, sousItem{ID: r.ID, Title: r.Title})
		}
		for _, it := range openSmallTasks(m.items) {
			out = append(out, sousItem{ID: it.ID, Title: it.Title})
		}
		return out
	}
}

// sousCounts is how many items each source would queue.
//
// The order says where the work comes from, and the count says how much there is:
// a source with nothing in it is a source that would send the crew nowhere.
func (m *Model) sousCounts() map[string]int {
	out := map[string]int{}
	for _, source := range sourceNames {
		out[source] = len(m.sousQueue(source))
	}
	return out
}

// sousLines renders the Sous-Chef tab.
//
// Every row is a row of the ticket: the order at the top, the queue under it and
// the stamp at the foot. No row holds a selectable item, so the console's cursor
// has nowhere to land on this tab.
func (m *Model) sousLines() []listLine {
	var out []listLine
	add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }
	add(" " + m.stySec.Render("THE ORDER"))
	out = append(out, m.sousOrderLines()...)
	out = append(out, m.sousQueueLines()...)
	out = append(out, m.sousStampLines()...)
	return out
}

// sousQueueLines renders the queue: what the crew will work, and what is next.
//
// The queue is the plan's own work, chosen by the order's source. The first item
// is the **next** one: the crew takes the queue in order.
func (m *Model) sousQueueLines() []listLine {
	var out []listLine
	add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }

	add(" " + m.styRule.Render(strings.Repeat("─", max(10, m.listColumn()))))
	add(" " + m.stySec.Render("THE QUEUE"))

	queue := m.sousQueue(m.picks.Sources)
	for i, item := range queue {
		if i >= sousQueueRows {
			break
		}
		row := "  " + m.styFaint.Render(padRight(item.ID, 6)) +
			padRight(item.Title, 42)
		if i == 0 {
			row += m.styAcc.Render("next")
		}
		add(row)
	}
	if len(queue) > sousQueueRows {
		add("  " + m.styFaint.Render(fmt.Sprintf("and %d more", len(queue)-sousQueueRows)))
	}
	add(" " + m.styRule.Render(strings.Repeat("─", max(10, m.listColumn()))))
	return out
}

// sousQueueRows is how many items of the queue the tab draws.
//
// The queue can be long and the ticket is short: the rest is counted, not drawn.
const sousQueueRows = 4

// sousStampLines renders the stamp: where the work stands.
//
// **The stamp says the progress, and the order says the policy.** The budget is
// the cap — "3 rounds · 12 interactions" — and the stamp is what has happened
// against it: "round 2 of 3 · 6 of 12". Two facts, two words, neither repeated.
//
// The stamp also carries the **vital signs**: the version of the rules the plan
// was written under, the version this build speaks, and whether they agree. The
// tab that starts the automation has to be able to report that it is not speaking
// the plan's rules.
func (m *Model) sousStampLines() []listLine {
	var out []listLine
	add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }

	if m.crewErr != nil {
		// The runtime is not here, or it did not answer. The tab says which,
		// and it offers **no move**: a tab that cannot start the automation
		// must not look like one that can.
		add(" " + m.stySec.Render("THE STAMP") + m.styBad.Render(" · "+m.crewErr.Error()))
		add("  " + m.styDim.Render("the automation cannot be started from here"))
		return out
	}
	if len(m.crew) == 0 {
		// The crew is not running. The console cannot tell "never started"
		// from "gone", and it does not guess: it says the automation is not
		// running, and the move comes back.
		add(" " + m.stySec.Render("THE STAMP") + m.styFaint.Render(" · the automation is not running"))
		add("  " + m.styFaint.Render("press space to start it"))
		return out
	}

	versions := m.styDim.Render(fmt.Sprintf("plan %d · build %d · they agree",
		m.planVersion(), core.DomainVersion))
	if !m.versionsAgree() {
		versions = m.styBad.Render(fmt.Sprintf("plan %d · build %d · they disagree",
			m.planVersion(), core.DomainVersion))
	}
	// The heading and the vital signs are **two facts**, and they read as two:
	// the separator is here whether a progress was reported or not.
	head := " " + m.stySec.Render("THE STAMP") + m.styDim.Render(" · ")
	if progress := m.sousProgress(); progress != "" {
		head += progress + m.styDim.Render(" · ")
	}
	add(head + versions)

	out = append(out, m.sousCrewLines()...)
	if !m.crewHasOrchestrator() {
		// A crew with no orchestrator has nobody to ask for directions: the
		// blocked agent's path is broken, and the tab must not pretend
		// otherwise.
		add("  " + m.styWarn.Render("the orchestrator is gone — the crew has nobody to ask for directions"))
	}
	return out
}

// sousCrewLines renders the crew's rows.
func (m *Model) sousCrewLines() []listLine {
	var out []listLine
	for _, row := range m.crew {
		state := row.State
		note := row.On
		if state == "blocked" {
			// **A blocked agent asks the orchestrator for directions**, and the
			// orchestrator decides. The console gains no move for it: the stamp
			// shows that the reviewer asked and that the orchestrator has it.
			state = "asked"
			note = "asked the orchestrator · the orchestrator has it"
		}
		out = append(out, listLine{
			text:    m.sousCrewRow(state, row.Job, note),
			itemIdx: -1,
		})
	}
	return out
}

// sousCrewRow renders one member of the crew.
func (m *Model) sousCrewRow(state, job, note string) string {
	return " " + m.sousStateStyle(state).Render(padRight(state, 8)) +
		m.styDim.Render(padRight(job, 13)) + m.styDim.Render(note)
}

// sousStateStyle is the style a crew member's state wears.
func (m *Model) sousStateStyle(state string) lipgloss.Style {
	switch state {
	case "working":
		return m.styGood
	case "asked", "blocked":
		return m.styWarn
	case "done":
		return m.styDim
	default:
		return m.styDim
	}
}

// sousProgress is the progress the stamp carries.
//
// The progress is the **orchestrator's** to say and the console's to read: it
// comes from the orchestrator's label, and it is printed against the cap the
// chef set. A number the console cannot compute is not printed — a label that
// carries no progress means the crew has not reported one.
func (m *Model) sousProgress() string {
	progress, ok := m.reportedProgress()
	if !ok {
		return ""
	}
	budget := sousBudgets[m.picks.Budget]
	// The budget is spent when the report has **reached** the cap: a crew that
	// has used its third round of three has spent its budget, and the flag
	// arrives **with** the report — not one round later.
	spent := (progress.Round.Reported && progress.Round.Value >= budget.Rounds) ||
		(progress.Interactions.Reported && progress.Interactions.Value >= budget.Interactions)

	var parts []string
	if progress.Round.Reported {
		if progress.Round.Value > budget.Rounds {
			// A progress past the cap is printed **as it was reported**:
			// "round 7 of 3" says the policy is not a policy, and the stamp
			// does not say that.
			parts = append(parts, fmt.Sprintf("round %d", progress.Round.Value))
		} else {
			parts = append(parts, fmt.Sprintf("round %d of %d", progress.Round.Value, budget.Rounds))
		}
	}
	if progress.Interactions.Reported {
		if progress.Interactions.Value > budget.Interactions {
			parts = append(parts, fmt.Sprintf("%d interactions", progress.Interactions.Value))
		} else {
			parts = append(parts, fmt.Sprintf("%d of %d interactions",
				progress.Interactions.Value, budget.Interactions))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	out := m.styDim.Render(strings.Join(parts, " · "))
	if spent {
		out += m.styWarn.Render(" · the budget is spent")
	}
	return out
}

// sousOrderLines renders the order: the four controls.
//
// The **harness** says the command it will run, so `custom` is never a mystery.
// The **crew** says who reviews, because the count is a composition and not a
// number: one agent cannot review itself, and the close gate refuses a
// self-reviewed recipe — so the one-agent crew is the one where the chef holds
// the reviewer's chair.
func (m *Model) sousOrderLines() []listLine {
	var out []listLine
	add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }

	// The label column, so every control's rows line up.
	const indent = "           "

	// The harness, and the command it will run.
	options := make([]string, 0, len(harnessNames))
	for _, h := range harnessNames {
		options = append(options, m.sousOption(h == m.picks.Harness, h))
	}
	add(m.sousLabel(labelHarness, 0) + strings.Join(options, "  "))
	add(indent + m.styFaint.Render("command: "+harnessCommand(m.picks.Harness)))

	// The crew, and who reviews.
	crew := []string{}
	for _, n := range crewSizes {
		crew = append(crew, m.sousOption(n == m.picks.Count, fmt.Sprintf("%d · %s", n, crewConsequence[n])))
	}
	add(m.sousLabel(labelCrew, 1) + crew[0])
	for _, line := range crew[1:] {
		add(indent + line)
	}

	// The sources, and how much each would queue.
	counts := m.sousCounts()
	sources := make([]string, 0, len(sourceNames))
	for _, s := range sourceNames {
		sources = append(sources, m.sousOption(s == m.picks.Sources,
			fmt.Sprintf("%s (%d)", s, counts[s])))
	}
	add(m.sousLabel(labelSources, 2) + strings.Join(sources, "  "))

	// The budget: the cap the chef sets.
	budget := sousBudgets[m.picks.Budget]
	add(m.sousLabel(labelBudget, 3) + m.styDim.Render(fmt.Sprintf("%d rounds · %d interactions",
		budget.Rounds, budget.Interactions)))
	return out
}

// The four controls of the order, in the order the tab shows them.
const (
	labelHarness = "  harness  "
	labelCrew    = "  crew     "
	labelSources = "  sources  "
	labelBudget  = "  budget   "
)

// sousLabel renders a control's label.
//
// The label under the cursor wears the console's selected style: the chef can see
// where the next key lands.
func (m *Model) sousLabel(label string, index int) string {
	if index == m.selSous {
		return m.stySelected.Render(label)
	}
	return m.styDim.Render(label)
}

// sousOption renders one option of a control.
//
// The picked option wears the marker; the others are faint, so the eye finds the
// choice and not the list.
func (m *Model) sousOption(picked bool, text string) string {
	if picked {
		return m.styAcc.Render("▸ " + text)
	}
	return m.styFaint.Render("   " + text)
}

// sousControls is how many controls the order has.
const sousControls = 4

// moveSous moves the cursor over the order's controls.
func (m *Model) moveSous(delta int) {
	m.selSous += delta
	if m.selSous < 0 {
		m.selSous = 0
	}
	if m.selSous >= sousControls {
		m.selSous = sousControls - 1
	}
}

// choosePick changes the pick under the cursor, and writes it.
//
// **The pick is written as soon as it changes.** The record is the order, and an
// order the console forgot is an order the next console does not have.
func (m *Model) choosePick(delta int) {
	switch m.selSous {
	case 0:
		m.picks.Harness = cycleName(harnessNames, m.picks.Harness, delta)
		m.status = "harness: " + m.picks.Harness
	case 1:
		m.picks.Count = cycleInt(crewSizes, m.picks.Count, delta)
		m.status = fmt.Sprintf("crew: %d · %s", m.picks.Count, crewConsequence[m.picks.Count])
	case 2:
		m.picks.Sources = cycleName(sourceNames, m.picks.Sources, delta)
		m.status = "sources: " + m.picks.Sources
	case 3:
		m.picks.Budget = cycleIndex(len(sousBudgets), m.picks.Budget, delta)
		budget := sousBudgets[m.picks.Budget]
		m.status = fmt.Sprintf("budget: %d rounds · %d interactions", budget.Rounds, budget.Interactions)
	}
	m.saveSettings()
}

// applyOrder writes the order and says what was written.
//
// The picks are written as they change, so `enter` writes the same order again —
// and says so, so the chef knows the record holds it.
func (m *Model) applyOrder() {
	m.saveSettings()
	if m.settingsErr != nil {
		return
	}
	m.status = "the order is saved: " + m.sousOrderLine()
}

// sousOrderLine is the order in one line, for the console's message.
func (m *Model) sousOrderLine() string {
	budget := sousBudgets[m.picks.Budget]
	return fmt.Sprintf("%s · crew %d · %s · %d rounds · %d interactions",
		m.picks.Harness, m.picks.Count, m.picks.Sources, budget.Rounds, budget.Interactions)
}

// cycleName is the name next to the given one, wrapping around.
func cycleName(names []string, current string, delta int) string {
	for i, n := range names {
		if n == current {
			return names[(i+delta+len(names))%len(names)]
		}
	}
	return names[0]
}

// cycleInt is the number next to the given one, wrapping around.
func cycleInt(values []int, current, delta int) int {
	for i, v := range values {
		if v == current {
			return values[(i+delta+len(values))%len(values)]
		}
	}
	return values[0]
}

// cycleIndex is the index next to the given one, wrapping around.
func cycleIndex(length, current, delta int) int {
	if length < 1 {
		return 0
	}
	return (current + delta + length) % length
}
