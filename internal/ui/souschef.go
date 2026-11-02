package ui

import (
	"fmt"
	"strconv"

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

// sousLines renders the Sous-Chef tab.
//
// Every row is a row of the ticket: the order at the top, the queue under it and
// the stamp at the foot. No row holds a selectable item, so the console's cursor
// has nowhere to land on this tab.
func (m *Model) sousLines() []listLine {
	var out []listLine
	add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }
	add(" " + m.stySec.Render("THE ORDER"))
	add("  " + m.styDim.Render("the automation's order, its queue and the crew's stamp"))
	return out
}
