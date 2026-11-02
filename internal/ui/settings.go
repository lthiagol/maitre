package ui

import (
	"github.com/lthiagol/maitre/internal/settings"
)

// The console's view preferences live in a file of their own: `kitchen/settings.json`.
//
// The file is the **tool's**, not the plan's. It is not a plan record, it carries no
// rule the plan was written under, and the domain does not read it — so the plan's
// format does not move, and the console's choices travel with the plan.
//
// The console is the only writer: it owns the vocabulary — the palette names and the
// order names — and it is the one that knows what a theme and an order are.

// settingsPath is where the record lives: in the plan's folder, beside the plan's
// records.
func (m *Model) settingsPath() string {
	return settings.Path(m.store.Root)
}

// loadSettings reads the record and applies the chef's choices.
//
// A plan with no settings file has the defaults, so a console nobody told otherwise
// looks exactly as it did.
func (m *Model) loadSettings() {
	rec, err := settings.Read(m.settingsPath())
	if err != nil {
		m.status = "the settings could not be read: " + err.Error()
		return
	}
	m.applyView(rec.View)
}

// applyView applies the record's view to the console.
//
// Each value is applied on its own: a theme the console does not know does not stop
// the order from being read.
func (m *Model) applyView(view settings.View) {
	if p, ok := themeByName(view.Theme); ok {
		m.theme = newTheme(p)
		m.selTheme = themeIndex(view.Theme)
	}
	if o, ok := orderByName(view.Order); ok {
		m.order = o
	}
	m.showServed = view.ShowServed
	m.sections = m.orderedSections(m.recipes)
}

// saveSettings writes the chef's view choices to the record.
//
// The write is the console's alone: the CLI reads the record and the console writes
// it, because the console is the one that knows the vocabulary.
func (m *Model) saveSettings() {
	view := settings.View{
		Theme:      m.theme.name,
		Order:      OrderNames[int(m.order)],
		ShowServed: m.showServed,
	}
	if err := settings.Write(m.settingsPath(), settings.Record{View: view}); err != nil {
		m.status = "the settings could not be written: " + err.Error()
	}
}

// themeByName finds the palette with that name.
func themeByName(name string) (palette, bool) {
	for _, p := range palettes() {
		if p.name == name {
			return p, true
		}
	}
	return palette{}, false
}

// themeIndex is the row of the palette with that name, or -1.
func themeIndex(name string) int {
	for i, p := range palettes() {
		if p.name == name {
			return i
		}
	}
	return -1
}

// orderByName finds the order with that name.
func orderByName(name string) (Order, bool) {
	for i, n := range OrderNames {
		if n == name {
			return Order(i), true
		}
	}
	return OrderTriage, false
}
