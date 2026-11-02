package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// The runtime that holds the crew.
//
// **herdr is a runtime dependency, not a module.** The `herdr` binary is on the
// PATH, and the console asks it for the crew: who is running, what state each
// one is in, and what it is on. Nothing about it is a Go module, so the
// dependency budget does not move — and a console without herdr still runs, on
// a tab that says so.
//
// Every fact has one owner: **herdr knows the crew's state**, and the console
// asks it. The plan carries no crew state, and the console writes none.

// crewRuntime is the console's seam to the runtime that holds the crew.
//
// The console asks four things: is the runtime here, who is in the crew, and
// how to start and stop it. A test answers them without a runtime at all — the
// console's own tests must not depend on a binary being on the PATH.
type crewRuntime interface {
	// available reports whether the runtime is here.
	available() bool
	// crew asks the runtime for the crew: the agents, and the tabs they are on.
	//
	// **One call answers both.** The crew's state and the crew's jobs come from
	// the runtime, and a read that forked twice was a read that waited twice.
	crew() (herdrCrew, error)
	// start opens the crew's tab and runs the orchestrator in it.
	start(planRoot, project, harness, task string) error
	// stop closes the given tabs.
	stop(tabIDs []string) error
}

// herdrCrew is what the runtime knows, in one answer.
//
// The **agents** are the crew's states; the **tabs** are the crew's jobs —
// the roster's line is a tab's label.
type herdrCrew struct {
	// Agents are the runtime's agents.
	Agents []herdrAgent
	// Tabs are the runtime's tabs.
	Tabs []herdrTab
}

// herdrAgent is one agent herdr knows.
//
// The fields the console reads are the ones the design names: the agent's
// **cwd** (where it works — the anchor), its **state**, and the **tab** it is
// on (the roster's line).
type herdrAgent struct {
	PaneID string `json:"pane_id"`
	TabID  string `json:"tab_id"`
	Cwd    string `json:"cwd"`
	// ForegroundCwd is the directory the agent's process is in. It is the
	// fallback: an agent whose cwd is not reported still works somewhere.
	ForegroundCwd string `json:"foreground_cwd"`
	Status        string `json:"agent_status"`
}

// herdrTab is one tab herdr knows.
//
// The **label** is the roster's line: `maitre: runner · RCP03`. It is where the
// job lives — herdr knows an agent is an mcode, not which job it plays.
type herdrTab struct {
	TabID string `json:"tab_id"`
	Label string `json:"label"`
}

// herdrRuntime asks the herdr binary.
//
// The binary is called through `run`, so a test can answer without one.
type herdrRuntime struct {
	// run runs herdr and returns its output.
	run func(args ...string) ([]byte, error)
}

// newHerdrRuntime is the runtime the console uses.
func newHerdrRuntime() *herdrRuntime {
	return &herdrRuntime{run: runHerdr}
}

// runHerdr runs the herdr binary and returns what it printed.
func runHerdr(args ...string) ([]byte, error) {
	out, err := exec.Command("herdr", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("herdr %s: %w", strings.Join(args, " "), err)
	}
	return out, nil
}

// available reports whether herdr is on the PATH.
//
// A console without herdr is still a console: the tab says "herdr is not here"
// and offers nothing. A tab that cannot start the automation must not look like
// one that can.
func (h *herdrRuntime) available() bool {
	_, err := exec.LookPath("herdr")
	return err == nil
}

// herdrSnapshot is what `herdr api snapshot` answers with.
//
// The snapshot carries the **agents** and the **tabs** in one document, so one
// process answers both questions the crew read asks.
type herdrSnapshot struct {
	Result struct {
		Snapshot struct {
			Agents []herdrAgent `json:"agents"`
			Tabs   []herdrTab   `json:"tabs"`
		} `json:"snapshot"`
	} `json:"result"`
}

// crew asks the runtime for the crew.
//
// **One process.** The agents and the tabs come from the same answer, so the
// read forks once and not twice.
func (h *herdrRuntime) crew() (herdrCrew, error) {
	raw, err := h.run("api", "snapshot")
	if err != nil {
		return herdrCrew{}, err
	}
	var doc herdrSnapshot
	if err := json.Unmarshal(raw, &doc); err != nil {
		return herdrCrew{}, fmt.Errorf("read the crew: %w", err)
	}
	return herdrCrew{
		Agents: doc.Result.Snapshot.Agents,
		Tabs:   doc.Result.Snapshot.Tabs,
	}, nil
}

// herdrTabCreated is what `herdr tab create` answers with.
type herdrTabCreated struct {
	Result struct {
		RootPane struct {
			PaneID string `json:"pane_id"`
		} `json:"root_pane"`
	} `json:"result"`
}

// herdrWorkspaceCreated is what `herdr workspace create` answers with.
//
// The workspace's id is the field the tab create needs, and the answer's shape
// is read from the document: an id under the result, or one under a `workspace`
// object.
type herdrWorkspaceCreated struct {
	Result struct {
		WorkspaceID string `json:"workspace_id"`
		Workspace   struct {
			ID string `json:"id"`
		} `json:"workspace"`
		ID string `json:"id"`
	} `json:"result"`
}

// start opens the crew's tab and runs the orchestrator in it.
//
// **Inside a herdr session** the tab joins the chef's own workspace, beside the
// chef. **Outside one** the console makes a workspace for the plan, and the tab
// goes in it. Either way the cwd is the plan's root: the automation works the
// plan, and a tab on the wrong directory is an automation on nobody's plan.
//
// The label is the console's first roster line: `maitre: orchestrator`. After
// that the orchestrator keeps the labels current.
func (h *herdrRuntime) start(planRoot, project, harness, task string) error {
	workspace, err := h.workspace(planRoot, project)
	if err != nil {
		return err
	}
	raw, err := h.run("tab", "create",
		"--workspace", workspace,
		"--label", orchestratorLabel,
		"--cwd", planRoot,
		"--no-focus")
	if err != nil {
		return err
	}
	var created herdrTabCreated
	if err := json.Unmarshal(raw, &created); err != nil {
		return fmt.Errorf("read the tab: %w", err)
	}
	pane := created.Result.RootPane.PaneID
	if pane == "" {
		return fmt.Errorf("the tab was created without a pane")
	}
	// The harness is the command the crew runs; the task is what it is given.
	args := append([]string{"pane", "run", pane}, strings.Fields(harness)...)
	args = append(args, task)
	if _, err := h.run(args...); err != nil {
		return err
	}
	return nil
}

// workspace is the workspace the crew's tab goes in.
//
// The chef's own workspace when the console is inside a herdr session, and a
// workspace made for the plan when it is not.
func (h *herdrRuntime) workspace(planRoot, project string) (string, error) {
	if os.Getenv("HERDR_ENV") == "1" {
		if id := os.Getenv("HERDR_WORKSPACE_ID"); id != "" {
			return id, nil
		}
	}
	raw, err := h.run("workspace", "create",
		"--cwd", planRoot,
		"--label", project,
		"--no-focus")
	if err != nil {
		return "", err
	}
	var created herdrWorkspaceCreated
	if err := json.Unmarshal(raw, &created); err != nil {
		return "", fmt.Errorf("read the workspace: %w", err)
	}
	for _, id := range []string{created.Result.WorkspaceID, created.Result.Workspace.ID, created.Result.ID} {
		if id != "" {
			return id, nil
		}
	}
	return "", fmt.Errorf("the workspace was created without an id")
}

// stop closes the given tabs.
func (h *herdrRuntime) stop(tabIDs []string) error {
	for _, id := range tabIDs {
		if _, err := h.run("tab", "close", id); err != nil {
			return err
		}
	}
	return nil
}

// orchestratorLabel is the console's first roster line.
//
// The roster's lines are `maitre: <job>[ · <what it is on>]`, and the console
// writes the first one — the orchestrator's — when it spawns the crew. After
// that the orchestrator keeps the lines current.
const orchestratorLabel = "maitre: orchestrator"

// rosterPrefix is what a roster's line starts with.
//
// The prefix is the console's own mark: it says the tab belongs to the plan's
// automation, and it is what keeps the chef's own pane out of the crew.
const rosterPrefix = "maitre:"
