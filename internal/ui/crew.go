package ui

import (
	"errors"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// The crew — the agents the runtime holds, and how the console finds them.
//
// **Every fact has one owner.** The runtime knows the crew's **state**; the tab's
// **label** knows the job; and the console knows neither, so it asks.
//
// The **cwd** is the anchor: a workspace's own listing carries no path, so the
// agents' cwd is what says an agent works *this* plan. The label is the roster's
// line — `maitre: runner · RCP03` — and it is where the job lives: the runtime
// knows an agent is an mcode, not which job it plays.

// errNoRuntime is what the console says when the runtime is not here.
//
// A tab that cannot start the automation must not look like one that can, so the
// tab says this and offers nothing.
var errNoRuntime = errors.New("herdr is not here")

// crewMsg carries the crew, as the runtime knows it.
type crewMsg struct {
	rows []crewRow
	err  error
}

// refreshCrew is the command that asks the runtime for the crew.
//
// The read is a call to another program, so it happens in a command and not in
// the render: a frame must never wait on a process.
func (m *Model) refreshCrew() tea.Cmd {
	return func() tea.Msg {
		rows, err := m.readCrew()
		return crewMsg{rows: rows, err: err}
	}
}

// readCrew asks the runtime for the crew.
//
// The console finds the crew by the **agents' cwd** — the plan's root — and reads
// each one's **state** from the runtime. The **job** comes from the tab's label:
// a tab that is not on the roster is not a member of the crew, and the chef's own
// pane is not the crew.
func (m *Model) readCrew() ([]crewRow, error) {
	if !m.runtime.available() {
		return nil, errNoRuntime
	}
	agents, err := m.runtime.agents()
	if err != nil {
		return nil, err
	}
	tabs, err := m.runtime.tabs()
	if err != nil {
		return nil, err
	}

	labels := make(map[string]string, len(tabs))
	for _, t := range tabs {
		labels[t.TabID] = t.Label
	}

	root := filepath.Clean(m.store.ProjectRoot())
	var rows []crewRow
	for _, a := range agents {
		cwd := a.Cwd
		if cwd == "" {
			cwd = a.ForegroundCwd
		}
		if filepath.Clean(cwd) != root {
			continue
		}
		job, tail, ok := parseRoster(labels[a.TabID])
		if !ok {
			continue
		}
		// The tail carries the member's report: the progress for the
		// orchestrator, and what it is on for everybody.
		progress, on := parseProgress(tail)
		row := crewRow{State: a.Status, Job: job, On: on, TabID: a.TabID}
		if progress.Round > 0 || progress.Interactions > 0 {
			row.Progress = &progress
		}
		rows = append(rows, row)
	}
	sortCrew(rows)
	return rows, nil
}

// crewProgress is what the crew's orchestrator reported about its own work.
//
// The **orchestrator** owns the progress: it is the one dispatching the crew, so
// it is the one that knows the round and the interaction count. It writes them
// into its own tab's label — the roster's line — and the console reads them from
// there.
//
// A progress that is not there is **not an error**: it means the crew has not
// reported one.
type crewProgress struct {
	// Round is the round the crew is on.
	Round int
	// Interactions is how many interactions the crew has spent.
	Interactions int
}

// parseProgress reads the progress out of a roster's tail.
//
// The tail is the orchestrator's report: `round 2 · 6 interactions`. The format
// is exact — a round and a count, each with its word — so the console never
// guesses: a part it does not recognise carries no progress, and the rest of
// the tail is what the member is on.
func parseProgress(tail string) (progress crewProgress, on string) {
	var rest []string
	for _, part := range strings.Split(tail, "·") {
		part = strings.TrimSpace(part)
		switch {
		case strings.HasPrefix(part, "round "):
			if n, ok := atoi(strings.TrimPrefix(part, "round ")); ok {
				progress.Round = n
				continue
			}
		case strings.HasSuffix(part, " interactions"):
			if n, ok := atoi(strings.TrimSuffix(part, " interactions")); ok {
				progress.Interactions = n
				continue
			}
		}
		rest = append(rest, part)
	}
	return progress, strings.Join(rest, " · ")
}

// atoi is the number in the string, and whether it was one.
//
// A count that is not a number is not a count: the console does not read "two"
// as two, and it does not read a number it cannot compute.
func atoi(s string) (int, bool) {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || n < 0 {
		return 0, false
	}
	return n, true
}

// parseRoster reads a roster line.
//
// The line is the console's mark, the job, and what the member is on:
// `maitre: runner · RCP03`. A tab whose label does not carry the mark is not on
// the roster, and the console does not pretend it is.
func parseRoster(label string) (job, on string, ok bool) {
	line := strings.TrimSpace(label)
	if !strings.HasPrefix(line, rosterPrefix) {
		return "", "", false
	}
	line = strings.TrimSpace(strings.TrimPrefix(line, rosterPrefix))
	if line == "" {
		return "", "", false
	}
	job, on, _ = strings.Cut(line, "·")
	job = strings.TrimSpace(job)
	on = strings.TrimSpace(on)
	if job == "" {
		return "", "", false
	}
	return job, on, true
}

// crewJobOrder is the order the crew's rows are drawn in.
//
// The roster reads the way the work runs: the orchestrator, then the runner,
// then the reviewer. A member with a job the glossary does not list comes last.
var crewJobOrder = []string{"orchestrator", "runner", "reviewer"}

// sortCrew puts the crew's rows in the order the work runs.
func sortCrew(rows []crewRow) {
	rank := func(job string) int {
		for i, name := range crewJobOrder {
			if name == job {
				return i
			}
		}
		return len(crewJobOrder)
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if rank(rows[i].Job) != rank(rows[j].Job) {
			return rank(rows[i].Job) < rank(rows[j].Job)
		}
		return rows[i].Job < rows[j].Job
	})
}

// reportedProgress is the progress the crew's orchestrator reported.
//
// It is the **orchestrator's** to say: it is the one dispatching the crew, so it
// is the one that knows the round and the interaction count. It writes them into
// its own tab's label, and the console reads them from there.
//
// A progress that is not there is **not an error**: it means the crew has not
// reported one, and the stamp prints no progress at all.
func (m *Model) reportedProgress() (crewProgress, bool) {
	for _, row := range m.crew {
		if row.Job == "orchestrator" && row.Progress != nil {
			return *row.Progress, true
		}
	}
	return crewProgress{}, false
}

// crewRunning reports whether the crew is running.
//
// The crew is running when the runtime holds a member of it: the plan carries no
// crew state, and a console that was closed and reopened still finds the crew.
func (m *Model) crewRunning() bool {
	return len(m.crew) > 0
}

// crewTabIDs are the tabs the crew's members are on.
//
// The console closes these tabs to stop the crew — and only these: the chef's own
// panes are not the crew's to close.
func (m *Model) crewTabIDs() []string {
	out := make([]string, 0, len(m.crew))
	for _, row := range m.crew {
		if row.TabID != "" {
			out = append(out, row.TabID)
		}
	}
	return out
}

// crewHasOrchestrator reports whether the crew has an orchestrator.
//
// A crew with no orchestrator has nobody to ask for directions: the blocked
// agent's path is broken, and the tab must not pretend otherwise.
func (m *Model) crewHasOrchestrator() bool {
	for _, row := range m.crew {
		if row.Job == "orchestrator" {
			return true
		}
	}
	return false
}

// planVersion is the version of the rules the plan was last written under.
func (m *Model) planVersion() int {
	return core.PlanDomainVersion(m.journal)
}

// versionsAgree reports whether the plan and this build speak the same rules.
func (m *Model) versionsAgree() bool {
	return m.planVersion() == core.DomainVersion
}
