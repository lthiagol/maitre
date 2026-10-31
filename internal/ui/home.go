package ui

import (
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/lthiagol/maitre/internal/brand"
	"github.com/lthiagol/maitre/internal/core"
)

// The Home tab — the console's landing place.
//
// **The invariant: the Home tab shows numbers; it shows no items.** No row on Home
// is selectable, `space` does nothing there, and the tab's item count is zero. A
// test holds it, the way slice 9's invariant held *"a recipe is filed under **Needs
// you** if and only if the console has a chef move for it"*.
//
// The precedent is v1's **Overview** lane, dropped because it *restated the board*.
// Home shows the plan's counts and no items. The board shows the recipes; Home
// shows their counts. The recipes' counts appear on both — the chef's ruling, the
// full set. A second board would be a second story about the same plan.
//
// Every number is the domain's. `core.PlanStats` computes them, and the console
// renders them: a number the console counted itself is a number the CLI would count
// differently.

// homeLines renders the Home tab.
//
// Every line is a number, and **no line holds a selectable item** — the item index
// is -1 on every row, so the cursor has nowhere to land and the tab's item count
// is zero.
func (m *Model) homeLines() []listLine {
	content := m.homeContent()
	width := homeWidth(content)
	art := m.homeArt(content, width)

	out := make([]listLine, 0, len(content))
	for i, line := range content {
		text := line
		if art != nil {
			// The mark sits **beside** the numbers, in the column they leave.
			text = padRight(line, width+2)
			if i >= art.start && i < art.start+len(art.lines) {
				text += art.lines[i-art.start]
			}
		}
		out = append(out, listLine{text: text, itemIdx: -1})
	}
	return out
}

// homeWidth is the width of the numbers' column.
//
// The column is as wide as the widest line the tab draws, so the mark beside it
// starts where the numbers end — and a narrower column would let the mark sit on
// top of a number.
func homeWidth(content []string) int {
	w := 0
	for _, line := range content {
		if n := lipgloss.Width(line); n > w {
			w = n
		}
	}
	return w
}

// homeMark is the mark beside the numbers: the drawing, and the row it starts on.
type homeMark struct {
	lines []string
	start int
}

// homeArt is the mark for the tab, when the tab has room for it.
//
// The decoration yields to the content, in both dimensions:
//
//   - **the width** — the numbers keep their column, and the mark is drawn only in
//     what is left. The mark is drawn in its **full** form or not at all: half a
//     mark beside half a screen is neither the emblem nor a number, and the column
//     that would take it is the column the numbers need.
//   - **the height** — the numbers must be whole on the screen. A tab whose numbers
//     are cut off has no room for a face, however wide it is.
//
// At 80×24 — the classic default — the numbers fill the width and the tab is what
// it would be without the mark.
func (m *Model) homeArt(content []string, width int) *homeMark {
	if len(content) > m.listAreaHeight() {
		return nil
	}
	// The mark is drawn one column in, like the rest of the tab, so it needs its own
	// width and a column of air.
	if m.rowWidth()-width-2 < brand.Width(brand.Full)+1 {
		return nil
	}
	art := brand.Full
	// The mark sits in the middle of the numbers, so it reads as the tab's face and
	// not as the header of one block.
	start := (len(content) - len(art)) / 2
	if start < 0 {
		start = 0
	}
	return &homeMark{lines: art, start: start}
}

// planPathColumns is how many columns the plan's path gets.
//
// A path can be longer than the screen, and a row that wraps takes a row the layout
// did not count — the footer, which is the last row, leaves the screen. The tail is
// the part that names the plan, so the tail is what is kept. The name beside it is
// bounded to what the row has left, and the Settings tab's root row is bounded to
// the width: one rule, both places.
const planPathColumns = 56

// homeContent is the tab's content: the plan's numbers, block by block.
//
// Every number comes from `core.PlanStats` — the domain's answer. The console adds
// no count of its own.
func (m *Model) homeContent() []string {
	stats := core.PlanStatsOf(m.recipes, m.items, m.journal)

	var out []string
	add := func(text string) { out = append(out, text) }
	block := func(title string) {
		if len(out) > 0 {
			add("")
		}
		add(" " + m.stySec.Render(title))
	}
	// stat renders one "label value" pair, and join puts a line's pairs together.
	stat := func(label, value string) string {
		return m.styDim.Render(label+" ") + m.styFaint.Render(value)
	}
	join := func(pairs ...string) string {
		return strings.Join(pairs, m.styDim.Render(" · "))
	}

	block("THE PLAN")
	add(m.planRow())

	block("THE RECIPES")
	// Every triage section, so the parts add up to the total. A section left out
	// is a number the row does not make.
	add("  " + join(
		stat("total", strconv.Itoa(stats.Recipes.Total)),
		stat("need you", strconv.Itoa(stats.Recipes.NeedsYou)),
		stat("on the line", strconv.Itoa(stats.Recipes.OnTheLine)),
		stat("being shaped", strconv.Itoa(stats.Recipes.BeingShaped)),
		stat("on hold", strconv.Itoa(stats.Recipes.OnHold)),
		stat("served", strconv.Itoa(stats.Recipes.Served)),
	))

	block("THE EVIDENCE")
	// Every evidence state, so the parts add up to the criteria. Pending is a
	// state the summary carries.
	add("  " + join(
		stat("criteria", strconv.Itoa(stats.Evidence.Criteria)),
		stat("proven", strconv.Itoa(stats.Evidence.Proven)),
		stat("failed", strconv.Itoa(stats.Evidence.Failed)),
		stat("rejected", strconv.Itoa(stats.Evidence.Rejected)),
		stat("unproven", strconv.Itoa(stats.Evidence.Unproven)),
		stat("pending", strconv.Itoa(stats.Evidence.Pending)),
	))

	block("THE TASTING")
	add("  " + join(
		stat("open notes", strconv.Itoa(stats.Tasting.OpenNotes)),
		stat("open findings", strconv.Itoa(stats.Tasting.OpenFindings)),
	))

	block("THE BACKLOG")
	add("  " + join(
		stat("ideas", strconv.Itoa(stats.Backlog.Ideas)),
		stat("small tasks", strconv.Itoa(stats.Backlog.SmallTasks)),
	))

	block("THE AUTOMATION")
	// The Sous-Chef tab said the automation was not built, and that was its honesty.
	// The tab dissolved into this block; the honesty comes with it, or the vital
	// signs read as a running automation.
	add("  " + m.styWarn.Render("the automation is not built"))
	agree := m.styDim.Render("they agree")
	if !stats.Automation.VersionsAgree {
		agree = m.styWarn.Render("they do not agree")
	}
	add("  " + join(
		stat("plan written under", strconv.Itoa(stats.Automation.PlanVersion)),
		stat("this build speaks", strconv.Itoa(stats.Automation.BuildVersion)),
		agree,
	))
	if !stats.Automation.VersionsAgree {
		add("  " + m.styDim.Render("this build does not write plans written under another version"))
	}
	if stats.Automation.LastEvent == nil {
		add("  " + m.styDim.Render("the plan has no events"))
		return out
	}
	e := *stats.Automation.LastEvent
	at := e.At
	if len(at) >= 19 {
		at = at[11:19]
	}
	recipe := e.Recipe
	if recipe == "" {
		recipe = "·"
	}
	add("  " + m.styDim.Render("last event ") +
		m.styFaint.Render(padRight(e.Kind, 12)+padRight(recipe, 6)+padRight(e.Actor, 10)) +
		m.styDim.Render(at))
	return out
}

// planRow is the plan's name and its path, on one line.
//
// The path is capped, and the name takes what the line has left. A name that does
// not fit is truncated to its tail — the same rule as the path — so the row stays
// one line and the footer stays the last row.
func (m *Model) planRow() string {
	path := truncate(humanPath(m.store.Root), planPathColumns)
	name := m.projectName()

	// "  name <name> · path <path>"
	room := m.width - len("  name ") - len(" · path ")
	if room < 2 {
		room = 2
	}
	pathW := lipgloss.Width(path)
	nameW := room - pathW
	if nameW < 1 {
		nameW = 1
		path = truncate(path, room-nameW)
	}
	name = truncate(name, nameW)

	stat := func(label, value string) string {
		return m.styDim.Render(label+" ") + m.styFaint.Render(value)
	}
	return "  " + stat("name", name) + m.styDim.Render(" · ") + stat("path", path)
}
