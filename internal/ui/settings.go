package ui

import (
	"fmt"
	"strings"

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
// looks exactly as it did. A value this build does not know is refused, by name, and
// the console **says so** — it does not fall back silently.
func (m *Model) loadSettings() {
	rec, err := settings.Read(m.settingsPath())
	if err != nil {
		m.status = "the settings could not be read: " + err.Error()
		return
	}
	if said := m.applyView(rec.View); len(said) > 0 {
		m.status = strings.Join(said, "; ")
	}
}

// applyView applies the record's view to the console.
//
// Each value is applied on its own: a theme the console does not know does not stop
// the order from being read. The values it does not know keep the console's own
// default, and the console says which value it refused.
func (m *Model) applyView(view settings.View) []string {
	var said []string
	if p, ok := themeByName(view.Theme); ok {
		m.theme = newTheme(p)
		m.selTheme = themeIndex(view.Theme)
	} else {
		said = append(said, unknownName("theme", view.Theme))
	}
	if o, ok := orderByName(view.Order); ok {
		m.order = o
	} else {
		said = append(said, unknownName("order", view.Order))
	}
	m.showServed = view.ShowServed
	m.sections = m.orderedSections(m.recipes)
	return said
}

// saveSettings writes the chef's view choices to the record.
//
// The write is the console's alone: the CLI reads the record and the console writes
// it, because the console is the one that knows the vocabulary. A value the console
// does not know is refused **before** the write — the way the domain's writer
// refuses before it writes — so the file never holds a value the console cannot
// render.
func (m *Model) saveSettings() {
	view := settings.View{
		Theme:      m.theme.name,
		Order:      OrderNames[int(m.order)],
		ShowServed: m.showServed,
	}
	if said := refuseUnknown(view); len(said) > 0 {
		m.status = strings.Join(said, "; ")
		return
	}
	if err := settings.Write(m.settingsPath(), settings.Record{View: view}); err != nil {
		m.status = "the settings could not be written: " + err.Error()
	}
}

// refuseUnknown names the values of a view the console does not know.
//
// The console never writes a value it does not know.
func refuseUnknown(view settings.View) []string {
	var out []string
	if _, ok := themeByName(view.Theme); !ok {
		out = append(out, unknownName("theme", view.Theme))
	}
	if _, ok := orderByName(view.Order); !ok {
		out = append(out, unknownName("order", view.Order))
	}
	return out
}

// unknownName is the sentence that names a value the console refused.
//
// The refusal says **which** value: a console that says only "unknown" leaves the
// chef guessing which of the two it meant.
func unknownName(kind, name string) string {
	if name == "" {
		return fmt.Sprintf("the record names no %s", kind)
	}
	return fmt.Sprintf("the %s %s is not one this build has", kind, name)
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
