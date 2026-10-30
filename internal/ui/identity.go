package ui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/lthiagol/maitre/internal/core"
)

// Option configures the console when it is built.
type Option func(*Model)

// WithProject names the plan for this session.
//
// The name is **display only**: it is what the chef sees — in the header and in the
// window title — and nothing else. The plan keeps its own name on disk, and the
// journal never hears about this.
func WithProject(name string) Option {
	return func(m *Model) { m.project = name }
}

// projectName is the name the console calls this plan.
//
// There is **one** answer, and everything that shows the name reads it — the header
// and the window title. Two computations would drift, and the drift would show as a
// terminal whose tab says one project and whose header says another.
//
// The `--project` flag wins: it is what the chef asked for. Without it, the plan's
// own name is the answer, and a plan that has no name falls back to its root.
func (m *Model) projectName() string {
	if m.project != "" {
		return m.project
	}
	if p, err := m.ops.GetProject(); err == nil && p.Name != "" {
		return p.Name
	}
	return m.store.Root
}

// windowTitle is the title the operating system puts on the terminal's tab.
//
// It is the row that answers *"which of my four terminals is this?"*, so it carries
// the project's name — not the product's.
func (m *Model) windowTitle() string {
	return "maitre · " + m.projectName()
}

// planRoot is the project's directory — the one that holds the plan directory.
//
// The store is rooted at `<project>/kitchen`, and the path a chef needs is the
// project: it is the directory they can open, and the one that tells two plans
// apart when their names are alike.
func (m *Model) planRoot() string {
	if filepath.Base(m.store.Root) == core.PlanDirName {
		return filepath.Dir(m.store.Root)
	}
	return m.store.Root
}

// humanPath shortens a path for a human: the home directory becomes `~`.
//
// A path with the home directory spelled out is longer than it needs to be, and the
// part that is the same in every one of the chef's paths is the part that says
// nothing.
func humanPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return path
	}
	if path == home {
		return "~"
	}
	if strings.HasPrefix(path, home+string(filepath.Separator)) {
		return "~" + path[len(home):]
	}
	return path
}
