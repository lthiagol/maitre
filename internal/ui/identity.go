package ui

// projectName is the name the console calls this plan.
//
// There is **one** answer, and everything that shows the name reads it — the header
// and the window title. Two computations would drift, and the drift would show as a
// terminal whose tab says one project and whose header says another.
//
// Without an override, the plan's own name is the answer, and a plan that has no
// name falls back to its root.
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
