package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// The backlog's moves.
//
// The console does not decide what the chef can do with a backlog item: it **asks**
// the domain — `core.BacklogMoveOf` — and it offers the move the domain names. A
// console-side rule ("is it a small task?") is the defect F-114 closed, and it does
// not come back here.
//
// The three moves:
//
//	shape it           an open idea — the kind changes and a scope is added
//	promote to recipe  an open small task — one criterion, signable
//	none               a settled item — no key, and the footer says so
//
// The dismissal sits beside them: `d` dismisses an item, with the chef's reason,
// and `a` writes a new item down.

// shapeChoices are the choices the shape prompt offers.
//
// The kinds come from the domain — `core.BacklogMoveOf` carries the kinds a shaping
// may pick — so the console carries **no list of its own**. A console-side copy of
// a domain fact is the shape F-114 closed: the console decides for itself and
// nothing holds it.
//
// The letter that stands for a kind is the console's own affordance — the domain
// owns the kind, not the key that picks it.
func shapeChoices(move core.BacklogMove) []promptChoice {
	out := make([]promptChoice, 0, len(move.Kinds))
	for _, kind := range move.Kinds {
		out = append(out, promptChoice{key: kindKey(kind), label: kind, value: kind})
	}
	return out
}

// kindKey is the letter that picks a kind in the prompt.
func kindKey(kind string) string {
	if kind == "" {
		return ""
	}
	return string(kind[0])
}

// backlogPrimary does the move the domain says the selected item has.
//
// The domain answers; the console asks. An item with no move gets no key and no
// attempt: the journal records attempts, and nobody attempted anything.
func (m *Model) backlogPrimary() (tea.Model, tea.Cmd) {
	item, ok := m.selectedBacklogItem()
	if !ok {
		m.status = "nothing here"
		return m, nil
	}
	move := core.BacklogMoveOf(item)
	switch move.Primary {
	case core.BacklogMoveShape:
		m.openShape(item.ID, move)
		return m, nil
	case core.BacklogMovePromote:
		return m.promoteSelected()
	}
	m.status = "nothing for you here — the item is settled"
	return m, nil
}

// openShape asks the chef to shape the item into a small task.
//
// Two answers, in the console's own prompt: what kind of work it is, and what is in
// it. The kinds are the domain's, and the scope is asked for because the domain
// needs one — a small task without a scope is work nobody can tell the end of.
func (m *Model) openShape(id string, move core.BacklogMove) {
	m.prompt = &chefPrompt{id: id, kind: promptShape, choices: shapeChoices(move)}
	m.refusal = ""
	m.status = ""
}

// recordShape writes the shaping.
//
// An **update**, not a promotion: the kind changes and a scope is added, and the id
// does not — `BL01` stays `BL01`. The domain's `UpdateBacklogItem` does it; the
// console asks.
func (m *Model) recordShape(p *chefPrompt) (tea.Model, tea.Cmd) {
	item, err := m.store.LoadBacklogItem(p.id)
	if err != nil {
		m.refusal = refusalText(err)
		m.status = ""
		return m, nil
	}
	// The scope is **added**: what the item already had stays, and the chef's
	// answer joins it.
	scope := append(append([]string{}, item.Scope...), strings.TrimSpace(p.text))
	shaped, err := m.ops.UpdateBacklogItem(p.id, "", "", scope, p.choice, nil)
	if err != nil {
		m.refusal = refusalText(err)
		m.status = ""
		return m, nil
	}
	m.refusal = ""
	m.status = fmt.Sprintf("shaped %s — it is a small task now, with a scope", shaped.ID)
	return m, m.reload()
}

// dismissSelected asks the chef why the item is not worth pursuing.
//
// The reason is asked for: the domain appends it to the item's body, and a
// dismissal without one says nothing to the next reader.
//
// The domain says whether the item may be dismissed at all — the console does not
// decide that a settled item is settled.
func (m *Model) dismissSelected() (tea.Model, tea.Cmd) {
	item, ok := m.selectedBacklogItem()
	if !ok {
		m.status = "nothing to dismiss here"
		return m, nil
	}
	if !core.BacklogMoveOf(item).Dismissable {
		m.status = "nothing for you here — the item is settled"
		return m, nil
	}
	m.prompt = &chefPrompt{id: item.ID, kind: promptDismiss, step: askText}
	m.refusal = ""
	m.status = ""
	return m, nil
}

// recordDismiss writes the dismissal and the reason.
//
// The domain appends the reason to the item's body, so the next reader knows why
// the item was dropped — in the chef's own words.
func (m *Model) recordDismiss(p *chefPrompt) (tea.Model, tea.Cmd) {
	reason := strings.TrimSpace(p.text)
	item, err := m.ops.DismissBacklogItem(p.id, reason)
	if err != nil {
		m.refusal = refusalText(err)
		m.status = ""
		return m, nil
	}
	m.refusal = ""
	m.status = fmt.Sprintf("dismissed %s — %s", item.ID, reason)
	return m, m.reload()
}

// canDismissSelected reports whether the domain says the selected item may still be
// dismissed.
//
// The domain answers it, with the rest of the item's moves — the console does not
// decide that a settled item is settled.
func (m *Model) canDismissSelected() bool {
	item, ok := m.selectedBacklogItem()
	if !ok {
		return false
	}
	return core.BacklogMoveOf(item).Dismissable
}

// openCreate asks the chef what the hunch is.
//
// A backlog item is **not** a recipe: the console may write one down (the chef's
// ruling, Q6) and it may never author a recipe (ADR-0003).
func (m *Model) openCreate() {
	m.prompt = &chefPrompt{kind: promptCreate, step: askText}
	m.refusal = ""
	m.status = ""
}

// recordCreate writes the item down.
//
// The item starts as an **idea** — the kind a hunch starts as — and the chef shapes
// it after, with the move the domain names.
//
// Nothing else is written: no recipe, and nothing outside the plan's backlog.
func (m *Model) recordCreate(p *chefPrompt) (tea.Model, tea.Cmd) {
	item, err := m.ops.CreateBacklogItem(strings.TrimSpace(p.text), "", nil, core.KindIdea, nil)
	if err != nil {
		m.refusal = refusalText(err)
		m.status = ""
		return m, nil
	}
	m.refusal = ""
	m.status = fmt.Sprintf("wrote down %s — an idea; shape it when it is ready", item.ID)
	return m, m.reload()
}
