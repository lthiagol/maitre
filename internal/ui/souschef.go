package ui

import (
	"fmt"
	"strconv"
	"strings"

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
