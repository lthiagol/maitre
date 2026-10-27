package ui

import (
	"fmt"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// triageChoices are the actions the chef may take on a tasting note.
//
// The fixer's move is not among them: resolving a note is the job of the agent that
// did the work, and the console must not do an agent's job (ADR-0003 §4).
var triageChoices = []promptChoice{
	{key: "d", label: "dismiss", value: "dismiss"},
	{key: "c", label: "comment", value: "comment"},
}

// openNotes are the tasting notes of the selected recipe that are still open.
//
// The same set the gate sees: `review.clean` refuses a recipe with an open note, so
// the chef's cursor and the gate must not disagree about which notes are open.
func (m *Model) openNotes() []core.Finding {
	rec, ok := m.currentRecipe()
	if !ok {
		return nil
	}
	return core.OpenFindingsOf(rec)
}

// selectedNote is the tasting note under the cursor.
func (m *Model) selectedNote() (core.Finding, bool) {
	notes := m.openNotes()
	if m.selNote < 0 || m.selNote >= len(notes) {
		return core.Finding{}, false
	}
	return notes[m.selNote], true
}

// moveNote moves the cursor over the notes.
//
// It wraps: there are only ever a handful of notes, and the chef should not have to
// know how many there are to get back to the one they were reading.
func (m *Model) moveNote(delta int) {
	notes := m.openNotes()
	if len(notes) == 0 {
		m.selNote = 0
		return
	}
	m.selNote = (m.selNote + delta + len(notes)) % len(notes)
}

// openTriage asks the chef what to do with a tasting note.
//
// The console offers the chef's moves only: dismiss the note, or say something in its
// thread. Resolving it is the fixer's move.
func (m *Model) openTriage(id, findingID string) {
	m.prompt = &chefPrompt{
		id:      id,
		kind:    promptTriage,
		choices: triageChoices,
		finding: findingID,
	}
	m.refusal = ""
	m.status = ""
}

// recordTriage writes what the chef did with the note.
//
// The chef is the author: the dismissal and the comment carry the chef's name, the
// same name the journal records.
func (m *Model) recordTriage(p *chefPrompt) (tea.Model, tea.Cmd) {
	switch p.choice {
	case "dismiss":
		r, err := m.ops.DismissFinding(p.id, p.finding, p.text)
		if err != nil {
			m.refusal = refusalText(err)
			m.status = ""
			return m, m.reload()
		}
		m.refusal = ""
		m.status = fmt.Sprintf("dismissed %s on %s — %s", p.finding, r.ID, p.text)
	case "comment":
		r, err := m.ops.AddComment(p.id, "chef", p.text, p.finding)
		if err != nil {
			m.refusal = refusalText(err)
			m.status = ""
			return m, m.reload()
		}
		m.refusal = ""
		m.status = fmt.Sprintf("commented on %s of %s", p.finding, r.ID)
	}
	return m, m.reload()
}
