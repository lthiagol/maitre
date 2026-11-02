package ui

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
