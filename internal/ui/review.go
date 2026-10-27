package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// promptStep is how far into the prompt the chef is.
type promptStep int

const (
	askChoice promptStep = iota
	askText
)

// promptKind is what the prompt is asking about.
type promptKind int

const (
	promptReview promptKind = iota
	promptTriage
)

// promptChoice is one thing the chef may choose.
type promptChoice struct {
	// key is the character that picks it.
	key string
	// label is what the footer calls it.
	label string
	// value is what it means to the operation.
	value string
}

// reviewChoices are the verdicts the chef may give on a recipe.
var reviewChoices = []promptChoice{
	{key: "o", label: "ok", value: "ok"},
	{key: "c", label: "changes needed", value: "changes-needed"},
}

// chefPrompt is the console asking the chef to decide.
//
// One prompt, two questions: a verdict on a recipe (ok, or changes needed) and an
// action on a tasting note (dismiss, or comment). Both have the same shape — choose,
// then say why — so they share one prompt rather than two frameworks.
//
// The console never decides — it asks. The chef chooses, says why, and the round is
// recorded with their name on it. A verdict the console invented would be a verdict
// nobody gave.
type chefPrompt struct {
	id      string
	kind    promptKind
	choices []promptChoice
	// choice is the value the chef picked.
	choice string
	// text is what the chef typed.
	text string
	step promptStep
	// finding is the tasting note being triaged.
	finding string
}

// title is what the prompt's footer calls the thing being decided.
func (p *chefPrompt) title() string {
	if p.kind == promptTriage {
		return "triage " + p.id + " " + p.finding
	}
	return "review " + p.id
}

// textTitle is what the footer calls the step where the chef types why.
func (p *chefPrompt) textTitle() string {
	if p.kind == promptTriage {
		return "notes for " + p.id + " " + p.finding
	}
	return "notes for " + p.id
}

// cancelled is what the status says when the chef backs out.
func (p *chefPrompt) cancelled() string {
	if p.kind == promptTriage {
		return "triage cancelled"
	}
	return "review cancelled"
}

// openReview asks the chef for a verdict on the selected recipe.
func (m *Model) openReview(id string) {
	m.prompt = &chefPrompt{id: id, kind: promptReview, choices: reviewChoices}
	m.refusal = ""
	m.status = ""
}

// handlePromptKey drives the open prompt.
//
// While the prompt is open every key belongs to it: the console must not move the
// selection out from under a half-typed answer.
func (m *Model) handlePromptKey(k tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	p := m.prompt
	switch k.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.prompt = nil
		m.status = p.cancelled()
		return m, nil
	case "enter":
		if p.step == askChoice {
			return m, nil
		}
		return m.recordPrompt()
	case "backspace":
		if p.step == askText {
			r := []rune(p.text)
			if len(r) > 0 {
				p.text = string(r[:len(r)-1])
			}
		}
		return m, nil
	}
	if p.step == askChoice {
		for _, ch := range p.choices {
			if promptKey(k) == ch.key {
				p.choice = ch.value
				p.step = askText
				break
			}
		}
		return m, nil
	}
	// The text: anything printable goes in.
	if text := typedText(k); text != "" {
		p.text += text
	}
	return m, nil
}

// promptKey is the character a key press stands for.
//
// Text carries the printable characters on a real terminal; the scripted driver
// builds a press from a code alone, so fall back to the code.
func promptKey(k tea.KeyPressMsg) string {
	if k.Text != "" {
		return k.Text
	}
	return k.String()
}

// typedText is what a key press contributes to the text.
//
// Text carries the printable characters on a real terminal; the scripted driver
// builds a press from a code alone, so fall back to the code — but never for a
// control character.
func typedText(k tea.KeyPressMsg) string {
	if k.Text != "" {
		return k.Text
	}
	if k.Mod != 0 {
		return ""
	}
	if k.Code < 32 || k.Code == 127 {
		return ""
	}
	return string(k.Code)
}

// recordPrompt writes what the chef decided.
//
// One prompt, two writes: what it writes depends on what it was asking about.
func (m *Model) recordPrompt() (tea.Model, tea.Cmd) {
	p := m.prompt
	m.prompt = nil
	if p.kind == promptTriage {
		return m.recordTriage(p)
	}
	return m.recordReview(p)
}

// recordReview writes the round the chef just gave.
//
// The reviewer is the chef, and the round goes into the artifact with their name on
// it — the same name the journal records.
func (m *Model) recordReview(p *chefPrompt) (tea.Model, tea.Cmd) {
	r, err := m.ops.RecordReviewRound(p.id, "chef", p.choice, p.text)
	if err != nil {
		m.refusal = refusalText(err)
		m.status = ""
		return m, nil
	}
	m.refusal = ""
	if p.choice != "ok" {
		m.status = fmt.Sprintf("asked for changes on %s — the agents have it", r.ID)
		return m, m.reload()
	}
	// A clean verdict moves the recipe on, if it was waiting to be reviewed.
	if r.Lifecycle == core.LifecycleExecuted {
		r, err = m.ops.Advance(r.ID, core.EventPassReview)
		if err != nil {
			m.refusal = refusalText(err)
			m.status = ""
			return m, m.reload()
		}
		m.status = fmt.Sprintf("reviewed %s — it is clean, serve it when you are ready", r.ID)
		return m, m.reload()
	}
	m.status = fmt.Sprintf("reviewed %s — it is clean, serve it when you are ready", r.ID)
	return m, m.reload()
}

// promptFooter is the footer while the prompt is open.
//
// The prompt takes the footer, because the footer is where the console says what the
// next key does.
func (m *Model) promptFooter() string {
	p := m.prompt
	if p.step == askChoice {
		parts := make([]string, 0, len(p.choices))
		for _, ch := range p.choices {
			parts = append(parts, ch.key+") "+ch.label)
		}
		return " " + m.styAcc.Render("▸ "+p.title()) +
			m.styDim.Render(" · "+strings.Join(parts, " · ")+" · esc cancel")
	}
	// The typed text keeps the tail: the tail is what was just typed.
	room := m.width - 46
	if room < 8 {
		room = 8
	}
	return " " + m.styAcc.Render("▸ "+p.textTitle()) +
		m.styDim.Render(" · ⏎ record · esc cancel · ") +
		m.styFaint.Render(truncate(p.text, room))
}
