package ui

import (
	"os"
	"path/filepath"
	"strings"
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
// The answer is the domain's: the console asks *what is this plan called?* and does
// not answer it a second time. The `--project` flag is what the chef asked for, so
// it wins.
//
// The name is a **rendered string**, and it goes through the same door as every
// other one: a newline in a name would corrupt the screen, and the plan's own file
// is written by the tool — `Init` takes the directory's basename, and a directory
// may carry a newline.
func (m *Model) projectName() string {
	if m.project != "" {
		return sanitize(m.project)
	}
	return sanitize(m.store.ProjectName())
}

// windowTitle is the title the operating system puts on the terminal's tab.
//
// It is the row that answers *"which of my four terminals is this?"*, so it carries
// the project's name — not the product's.
func (m *Model) windowTitle() string {
	return "maitre · " + m.projectName()
}

// humanPath shortens a path for a human: the home directory becomes `~`.
//
// A path with the home directory spelled out is longer than it needs to be, and the
// part that is the same in every one of the chef's paths is the part that says
// nothing.
//
// The path is **sanitised** first: it is drawn on one row, and a directory whose
// name carries a newline would split it.
func humanPath(path string) string {
	path = sanitize(path)
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
