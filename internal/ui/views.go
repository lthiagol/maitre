package ui

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/glamour/v2/styles"
	"charm.land/lipgloss/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// The console owns the terminal: the alt screen, so the shell's scrollback stays clean
// and the terminal is restored exactly on exit.
const (
	headerRows = 4 // title, rule, tabs, summary
	footerRows = 1
)

// refusalMarker is what marks a refusal.
//
// The refusal is the one thing the chef must never miss, so it carries a marker of
// its own — a word could be missed in a wall of text.
const refusalMarker = "⛔ "

// splitThreshold is the width at which the detail pane appears beside the list.
//
// Below it, the list gets the whole terminal and the detail opens full screen. A split
// at 100 columns leaves neither side enough room: titles truncate on one side and the
// markdown wraps badly on the other.
const splitThreshold = 160

// View renders the console.
func (m *Model) View() tea.View {
	if m.err != nil {
		v := tea.NewView(m.styBad.Render("error: " + m.err.Error()))
		v.AltScreen = true
		return v
	}
	if !m.ready {
		v := tea.NewView(m.styDim.Render("opening the kitchen…"))
		v.AltScreen = true
		return v
	}

	v := tea.NewView(m.renderLayout())
	v.AltScreen = true
	v.BackgroundColor = lipgloss.Color(m.colBg)
	v.WindowTitle = "maitre"
	return v
}

// renderLayout draws the whole screen: header, the list, and a pinned footer.
func (m *Model) renderLayout() string {
	lines := m.listLines()
	m.scrollToSelection(lines)

	if m.mode == ModeDetail {
		return m.renderDetailFull()
	}
	if m.width >= splitThreshold && m.mode == ModeSplit {
		return m.renderSplit(lines)
	}
	return m.renderSingle(lines)
}

// renderDetailFull draws one recipe across the whole terminal.
func (m *Model) renderDetailFull() string {
	var b strings.Builder
	b.WriteString(m.header() + "\n")
	b.WriteString(m.ruleLine() + "\n")
	b.WriteString(m.tabBar() + "\n")
	b.WriteString(m.summaryLine() + "\n")
	m.writeBanner(&b)

	detail := detailRows(m.detailLines(m.width - 4))
	if m.detailScroll > 0 {
		if m.detailScroll >= len(detail) {
			m.detailScroll = len(detail) - 1
		}
		detail = detail[m.detailScroll:]
	}
	h := m.listAreaHeight()
	for i := 0; i < h; i++ {
		if i < len(detail) {
			b.WriteString(detail[i])
		}
		b.WriteString("\n")
	}
	// The footer is the last row of the screen, on a row of its own. Do not trim:
	// the trailing newline is the padding that holds the footer down.
	return b.String() + m.footer()
}

// renderSingle draws the list across the whole terminal.
func (m *Model) renderSingle(lines []listLine) string {
	var b strings.Builder
	b.WriteString(m.header() + "\n")
	b.WriteString(m.ruleLine() + "\n")
	b.WriteString(m.tabBar() + "\n")
	b.WriteString(m.summaryLine() + "\n")
	m.writeBanner(&b)

	h := m.listAreaHeight()
	for i := 0; i < h; i++ {
		idx := m.offset + i
		if idx < len(lines) {
			b.WriteString(lines[idx].text)
		}
		b.WriteString("\n")
	}

	// The footer is the last row of the screen. Do not trim: the trailing newlines
	// are the padding that holds the footer down.
	return b.String() + m.footer()
}

// renderSplit draws the list and the detail side by side, both full height.
func (m *Model) renderSplit(lines []listLine) string {
	listWidth := m.listColumn()
	detailWidth := m.width - listWidth - 3

	var left []string
	h := m.listAreaHeight()
	for i := 0; i < h; i++ {
		idx := m.offset + i
		if idx < len(lines) {
			left = append(left, lines[idx].text)
		} else {
			left = append(left, "")
		}
	}
	right := detailRows(m.detailLines(detailWidth))
	if len(right) > h {
		right = right[:h]
	}
	for len(right) < h {
		right = append(right, "")
	}

	var b strings.Builder
	b.WriteString(m.header() + "\n")
	b.WriteString(m.ruleLine() + "\n")
	b.WriteString(m.tabBar() + "\n")
	b.WriteString(m.summaryLine() + "\n")
	m.writeBanner(&b)
	b.WriteString(joinColumns(left, right, listWidth))
	return b.String() + m.footer()
}

// listAreaHeight is how many rows the list gets.
//
// A refusal banner takes rows from the list: the banner is the one thing the chef
// must not miss, and the list can give up a row.
func (m *Model) listAreaHeight() int {
	h := m.height - headerRows - len(m.bannerLines()) - footerRows
	if h < 1 {
		h = 1
	}
	return h
}

// scrollToSelection keeps the selected item visible.
//
// The list scrolls by line but the selection is by item, so a section header can scroll
// away while its items stay selected.
func (m *Model) scrollToSelection(lines []listLine) {
	sel := m.currentSel()
	if sel < 0 {
		return
	}
	line := -1
	for i, l := range lines {
		if l.itemIdx == sel {
			line = i
			break
		}
	}
	if line < 0 {
		return
	}
	h := m.listAreaHeight()
	if line < m.offset {
		m.offset = line
	} else if line >= m.offset+h {
		m.offset = line - h + 1
	}
	maxOffset := len(lines) - h
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.offset < 0 {
		m.offset = 0
	}
	if m.offset > maxOffset {
		m.offset = maxOffset
	}
}

// header is the title line.
func (m *Model) header() string {
	name := m.store.Root
	if p, err := m.ops.GetProject(); err == nil && p.Name != "" {
		name = p.Name
	}
	return m.styTitle.Render(" maitre") + m.styFaint.Render(" · "+name)
}

// ruleLine is the rule under the header.
func (m *Model) ruleLine() string {
	return m.styRule.Render(strings.Repeat("─", max(10, m.width)))
}

// tabBar is the row of tabs.
func (m *Model) tabBar() string {
	var parts []string
	for i, name := range TabNames {
		if Tab(i) == m.tab {
			parts = append(parts, m.styTabActive.Render("▸ "+name))
			continue
		}
		parts = append(parts, m.styTabIdle.Render("  "+name))
	}
	return strings.Join(parts, "  ")
}

// summaryLine is the per-tab summary, with the status or refusal on the right.
//
// A refusal that does not fit beside the summary text moves to rows of its own — see
// bannerLines. A refusal the chef cannot see is a refusal that did not happen.
func (m *Model) summaryLine() string {
	left := m.summaryText()
	right := m.summaryRight()
	if right == "" {
		return " " + left
	}
	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 3
	if gap < 1 {
		return " " + truncate(left, max(1, m.width-4))
	}
	style := m.styAcc
	if m.refusal != "" {
		style = m.styWarn
	}
	return " " + left + strings.Repeat(" ", gap) + style.Render(right)
}

// summaryRight is what goes on the right of the summary line.
//
// A refusal outranks the status: the chef asked for a write and the gate said no, and
// that is the one thing on the row that matters.
func (m *Model) summaryRight() string {
	if m.refusal != "" {
		return refusalMarker + m.refusal
	}
	return m.status
}

// summaryFits reports whether the right of the summary line fits beside the left.
func (m *Model) summaryFits() bool {
	right := m.summaryRight()
	if right == "" {
		return true
	}
	return m.width-lipgloss.Width(m.summaryText())-lipgloss.Width(right)-3 >= 1
}

// bannerLines are the rows a refusal takes under the summary line.
//
// A refusal that does not fit beside the summary text takes rows of its own, so the
// chef always reads it. The reason wraps; it is never truncated away.
func (m *Model) bannerLines() []string {
	if m.refusal == "" || m.summaryFits() {
		return nil
	}
	width := max(20, m.width-2)
	rows := wrapText(refusalMarker+m.refusal, width)
	room := m.bannerRoom()
	if len(rows) > room {
		// The screen is too short for the whole reason: keep what fits and say that
		// there is more.
		rows = rows[:room]
		rows[len(rows)-1] = rows[len(rows)-1] + "…"
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, " "+m.styWarn.Render(r))
	}
	return out
}

// writeBanner writes the refusal banner, if there is one.
func (m *Model) writeBanner(b *strings.Builder) {
	for _, l := range m.bannerLines() {
		b.WriteString(l + "\n")
	}
}

// bannerRoom is how many rows the banner may take.
//
// The screen minus the header and the footer, keeping one row for the list.
func (m *Model) bannerRoom() int {
	room := m.height - headerRows - footerRows - 1
	if room < 1 {
		room = 1
	}
	return room
}

// wrapText breaks s into lines of at most width columns.
//
// A word longer than the width is left long: breaking it would make a word unreadable,
// and a word that does not fit is a word the terminal wraps itself.
func wrapText(s string, width int) []string {
	if width < 1 {
		width = 1
	}
	var out []string
	cur := ""
	for _, word := range strings.Fields(s) {
		switch {
		case cur == "":
			cur = word
		case lipgloss.Width(cur)+1+lipgloss.Width(word) <= width:
			cur += " " + word
		default:
			out = append(out, cur)
			cur = word
		}
	}
	if cur != "" {
		out = append(out, cur)
	}
	if len(out) == 0 {
		out = append(out, "")
	}
	return out
}

// summaryText is the per-tab summary line.
func (m *Model) summaryText() string {
	switch m.tab {
	case TabRecipes:
		text := fmt.Sprintf("%d %s · %d need you",
			len(m.recipes), plural(len(m.recipes), "recipe", "recipes"), m.countNeedsYou())
		if m.order == OrderTriage {
			text += "   " + m.servedCheckbox()
		}
		return text
	case TabBacklog:
		return fmt.Sprintf("%d %s · %d %s",
			len(m.specials), plural(len(m.specials), "special", "specials"),
			len(m.notes), plural(len(m.notes), "prep note", "prep notes"))
	case TabPath:
		ready, waiting := core.Path(m.recipes)
		return fmt.Sprintf("the path · %d ready · %d waiting", len(ready), len(waiting))
	case TabChopboard:
		tree := core.TreeOf(m.recipes)
		waiting := 0
		var count func(n core.TreeNode)
		count = func(n core.TreeNode) {
			waiting += len(n.Waiting)
			for _, c := range n.Waiting {
				count(c)
			}
		}
		for _, r := range tree.Roots {
			count(r)
		}
		for _, o := range tree.Orphans {
			waiting++
			count(o)
		}
		for _, s := range tree.Stuck {
			waiting++
			count(s)
		}
		return fmt.Sprintf("the tree · %d recipes · %d can go now · %d waiting",
			len(tree.Roots)+waiting, len(tree.Roots), waiting)
	case TabLog:
		return fmt.Sprintf("recent activity · %d events", len(m.journal.Events))
	case TabExpo:
		return "the expo · the orchestrator · not built yet"
	default:
		return "the plan"
	}
}

// footer is the keybind bar, pinned to the bottom of the screen.
//
// It is contextual: `space` does something different for each kind of selection, so the
// bar says which. In the detail pane the list is not on screen, so the keys say what
// they do there.
func (m *Model) footer() string {
	// An open prompt takes the footer: the footer is where the console says what the
	// next key does.
	if m.prompt != nil {
		return m.promptFooter()
	}

	type key struct{ k, what string }
	var keys []key
	switch {
	case m.mode == ModeDetail:
		keys = append(keys, key{"↑↓", "scroll"}, key{"esc", "back"})
	case m.tab == TabSettings:
		// The Settings tab has no recipe to open: `enter` applies the theme.
		keys = append(keys, key{"↑↓", "move"}, key{"tab", "switch"}, key{"enter", "apply theme"})
	default:
		keys = append(keys, key{"↑↓", "move"}, key{"tab", "switch"}, key{"enter", "detail"})
	}
	if label := m.primaryActionLabel(); label != "" {
		keys = append(keys, key{"space", label})
	}
	// The note cursor is only worth a key when there is a note to point at.
	if len(m.openNotes()) > 0 {
		keys = append(keys, key{"n/N", "note"})
	}
	if m.tab == TabRecipes && m.mode != ModeDetail {
		keys = append(keys, key{"1 2 3", "order"}, key{"c", "served"})
	}
	keys = append(keys, key{"r", "reload"}, key{"q", "quit"})

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, m.styKey.Render(k.k)+m.styDim.Render(" "+k.what))
	}
	bar := " " + strings.Join(parts, m.styDim.Render(" · "))
	if m.mode == ModeDetail {
		bar += m.scrollIndicator()
	}
	return bar
}

// scrollIndicator says where in the recipe you are.
//
// A pane that scrolls without saying so is a pane that loses its place.
func (m *Model) scrollIndicator() string {
	total := m.detailLineCount()
	if total <= m.listAreaHeight() {
		return ""
	}
	at := m.detailScroll + 1
	if at > total {
		at = total
	}
	return m.styDim.Render(fmt.Sprintf("   ·   %d/%d", at, total))
}

// primaryActionLabel names what `space` will do for the current selection.
func (m *Model) primaryActionLabel() string {
	switch m.tab {
	case TabRecipes, TabPath:
		rec, ok := m.currentRecipe()
		if !ok {
			return ""
		}
		return core.ChefMoveOf(*rec).Label()
	case TabBacklog:
		if m.selBacklog < len(m.specials) {
			return "promote to prep note"
		}
		if m.selBacklog-len(m.specials) < len(m.notes) {
			return "promote to recipe"
		}
	}
	return ""
}

// listLines renders the current tab's list, marking which lines hold a selectable item.
func (m *Model) listLines() []listLine {
	switch m.tab {
	case TabBacklog:
		return m.backlogLines()
	case TabPath:
		return m.pathLines()
	case TabChopboard:
		return m.chopTree().rows
	case TabLog:
		return m.logLines()
	case TabExpo:
		return m.expoLines()
	case TabSettings:
		return m.settingsLines()
	default:
		return m.recipeLines()
	}
}

// recipeLines renders the Recipes tab, grouped by triage section.
func (m *Model) recipeLines() []listLine {
	var out []listLine
	item := 0
	for _, sec := range m.sections {
		if sec.Name == core.TriageServed && !m.showServed {
			continue
		}
		out = append(out, listLine{text: m.sectionHeader(sec.Name, len(sec.Items), m.listColumn()), itemIdx: -1})
		out = append(out, listLine{text: m.styRule.Render(strings.Repeat("─", m.listColumn())), itemIdx: -1})
		for _, s := range sec.Items {
			out = append(out, listLine{
				text:    m.recipeRow(s, item == m.selRecipes, m.listColumn()),
				itemIdx: item,
			})
			item++
		}
	}
	return out
}

// countNeedsYou counts the recipes waiting on the chef.
func (m *Model) countNeedsYou() int {
	for _, sec := range core.TriageSections(m.recipes) {
		if sec.Name == core.TriageNeedsYou {
			return len(sec.Items)
		}
	}
	return 0
}

// servedCheckbox is the show/hide control for served recipes.
func (m *Model) servedCheckbox() string {
	box := "☐"
	if m.showServed {
		box = "☑"
	}
	return m.styDim.Render(box + " show served")
}

// idColumn is the width of the id column in a recipe row.
//
// A recipe's id is `RCP01` — five characters — and the column leaves a space after it,
// so the id and the lifecycle never run together.
const idColumn = 6

// recipeRow renders one recipe in the list, sized to fit its column.
//
// A row wider than its column overflows into the pane beside it, so the title is
// truncated to whatever is left over.
func (m *Model) recipeRow(s core.Summary, selected bool, width int) string {
	marks := proofMarks(s)
	title := padCol(s.Title, titleColumn(width, lipgloss.Width(marks)))

	if selected {
		plain := "  " + padRight(s.ID, idColumn) + padRight(string(s.Lifecycle), 12) + title + marks
		return m.stySelected.Render(padRight(plain, width))
	}
	return "  " + m.styFaint.Render(padRight(s.ID, idColumn)) +
		m.lifecycleStyle(string(s.Lifecycle)).Render(padRight(string(s.Lifecycle), 12)) +
		title + m.styFaint.Render(marks)
}

// titleColumn is the width the title gets in a row of the given column width.
//
// The title is capped so the marks after it do not drift to the far edge of a wide
// terminal — a gap that big reads as two separate columns.
func titleColumn(width, tail int) int {
	w := width - 2 - idColumn - 12 - tail - 1
	if w > 44 {
		w = 44
	}
	if w < 8 {
		w = 8
	}
	return w
}

// rowWidth is the width of a full recipe row.
func (m *Model) rowWidth() int {
	w := m.width - 2
	if w < 20 {
		w = 20
	}
	return w
}

// listColumn is the width the list gets.
//
// Below the split threshold the list takes the whole terminal — a truncated title is
// worse than no preview.
func (m *Model) listColumn() int {
	if m.width < splitThreshold {
		return m.rowWidth()
	}
	w := m.width * 45 / 100
	if w < 34 {
		w = 34
	}
	return w
}

// sectionHeader is a section name with its count, right-aligned in the column.
func (m *Model) sectionHeader(name string, count, width int) string {
	n := fmt.Sprint(count)
	gap := width - lipgloss.Width(name) - lipgloss.Width(n) - 2
	if gap < 1 {
		gap = 1
	}
	return " " + m.stySec.Render(name) + strings.Repeat(" ", gap) + m.styDim.Render(n)
}

// backlogLines renders the Backlog tab.
func (m *Model) backlogLines() []listLine {
	var out []listLine
	if len(m.specials) > 0 {
		out = append(out, listLine{text: " " + m.stySec.Render("SPECIALS") + strings.Repeat(" ", max(1, m.listColumn()-10)) + m.styDim.Render(fmt.Sprint(len(m.specials))), itemIdx: -1})
		out = append(out, listLine{text: m.styRule.Render(strings.Repeat("─", m.listColumn())), itemIdx: -1})
		for i, sp := range m.specials {
			out = append(out, listLine{text: m.backlogRow(i, sp.ID, sp.Title, sp.Status, m.selBacklog == i), itemIdx: i})
		}
	}
	if len(m.notes) > 0 {
		out = append(out, listLine{text: " " + m.stySec.Render("PREP NOTES") + strings.Repeat(" ", max(1, m.listColumn()-11)) + m.styDim.Render(fmt.Sprint(len(m.notes))), itemIdx: -1})
		out = append(out, listLine{text: m.styRule.Render(strings.Repeat("─", m.listColumn())), itemIdx: -1})
		for i, n := range m.notes {
			idx := len(m.specials) + i
			out = append(out, listLine{text: m.backlogRow(idx, n.ID, n.Title, n.Status, m.selBacklog == idx), itemIdx: idx})
		}
	}
	return out
}

// backlogRow renders one backlog item.
func (m *Model) backlogRow(index int, id, title, status string, selected bool) string {
	width := m.listColumn()
	row := "  " + padRight(id, 6) + padCol(title, titleColumn(width, lipgloss.Width(status))) + m.styFaint.Render(status)
	if selected {
		return m.stySelected.Render(padRight(row, width))
	}
	return row
}

// pathLines renders the Path tab.
func (m *Model) pathLines() []listLine {
	ready, waiting := core.Path(m.recipes)
	var out []listLine
	item := 0

	section := func(name string, entries []core.PathEntry) {
		if len(entries) == 0 {
			return
		}
		out = append(out, listLine{text: " " + m.stySec.Render(name) + strings.Repeat(" ", max(1, m.listColumn()-8)) + m.styDim.Render(fmt.Sprint(len(entries))), itemIdx: -1})
		out = append(out, listLine{text: m.styRule.Render(strings.Repeat("─", m.listColumn())), itemIdx: -1})
		for _, e := range entries {
			out = append(out, listLine{text: m.pathRow(e), itemIdx: item})
			item++
		}
	}
	section("READY", ready)
	section("WAITING", waiting)
	return out
}

// pathRow renders one entry in the execution order.
func (m *Model) pathRow(e core.PathEntry) string {
	width := m.listColumn()
	note := e.Note
	if e.WaitingFor != "" {
		note = "needs " + e.WaitingFor
	}
	title := padCol(e.Summary.Title, titleColumn(width, lipgloss.Width(note)))
	return "  " + m.styFaint.Render(padRight(e.Summary.ID, 6)) + title + m.styDim.Render(note)
}

// chopBuilder builds the Chopboard's rows.
//
// The tree is drawn as a trunk with branches: the trunk is the ready list, and each
// branch holds the recipes that fork at the recipe above it.
type chopBuilder struct {
	m      *Model
	out    []listLine
	recipe []string
}

// trunk draws a trunk recipe and the branch that forks at it.
func (b *chopBuilder) trunk(node core.TreeNode) {
	b.out = append(b.out, listLine{text: b.m.chopRow(node, "● ", b.sel()), itemIdx: len(b.recipe)})
	b.recipe = append(b.recipe, node.Summary.ID)
	b.branch(node, "", true)
}

// sel reports whether the row being drawn is the selected one.
func (b *chopBuilder) sel() bool { return len(b.recipe) == b.m.selChop }

// orphan draws a recipe whose dependency is not in the plan.
//
// It cannot fork at anything, and the plan is broken: it is marked, and it says which
// dependency is missing.
func (b *chopBuilder) orphan(node core.TreeNode) {
	b.out = append(b.out, listLine{text: b.m.chopOrphanRow(node, b.sel()), itemIdx: len(b.recipe)})
	b.recipe = append(b.recipe, node.Summary.ID)
	b.branch(node, "", true)
}

// stuck draws a recipe a dependency cycle keeps out of the tree.
//
// The plan cannot be executed as written, and the tree says so rather than dropping
// the recipe.
func (b *chopBuilder) stuck(node core.TreeNode) {
	b.out = append(b.out, listLine{text: b.m.chopStuckRow(node, b.sel()), itemIdx: len(b.recipe)})
	b.recipe = append(b.recipe, node.Summary.ID)
	b.branch(node, "", true)
}

// chopOrphanRow renders a recipe whose dependency is not in the plan.
//
// The plan is broken, and the marker says so: the row carries the warning, not the
// whole line.
func (m *Model) chopOrphanRow(node core.TreeNode, selected bool) string {
	width := m.listColumn()
	note := chopNote(node)
	marks := proofMarks(node.Summary)
	tail := marks
	if note != "" {
		tail += " " + note
	}
	title := padCol(node.Summary.Title, titleColumn(width, lipgloss.Width(tail)))
	row := m.styWarn.Render("⚠ ") + m.styFaint.Render(padRight(node.Summary.ID, 6)) +
		m.lifecycleStyle(string(node.Summary.Lifecycle)).Render(padRight(string(node.Summary.Lifecycle), 12)) +
		title + m.styFaint.Render(marks) + m.styDim.Render(strings.TrimPrefix(tail, marks))
	if selected {
		return m.stySelected.Render(padRight(row, width))
	}
	return row
}

// chopStuckRow renders a recipe a dependency cycle keeps out of the tree.
//
// The plan cannot be executed as written, and the row says so.
func (m *Model) chopStuckRow(node core.TreeNode, selected bool) string {
	width := m.listColumn()
	note := chopNote(node)
	marks := proofMarks(node.Summary)
	tail := marks
	if note != "" {
		tail += " " + note
	}
	title := padCol(node.Summary.Title, titleColumn(width, lipgloss.Width(tail)))
	row := m.styWarn.Render("⚠ ") + m.styFaint.Render(padRight(node.Summary.ID, 6)) +
		m.lifecycleStyle(string(node.Summary.Lifecycle)).Render(padRight(string(node.Summary.Lifecycle), 12)) +
		title + m.styFaint.Render(marks) + m.styDim.Render(strings.TrimPrefix(tail, marks))
	if selected {
		return m.stySelected.Render(padRight(row, width))
	}
	return row
}

// branch draws the group of recipes that fork at a node.
//
// The heading names the branch's subject — the dependency they all wait on — and the
// rows under it are the recipes that fork there.
func (b *chopBuilder) branch(node core.TreeNode, prefix string, last bool) {
	if len(node.Waiting) == 0 {
		return
	}
	connector := "├─ "
	if last {
		connector = "└─ "
	}
	b.out = append(b.out, listLine{
		text:    " " + b.m.styDim.Render(prefix+connector+"needs "+node.Summary.ID),
		itemIdx: -1,
	})

	// The recipes that fork there, and their own branches, one level deeper.
	spine := prefix + "│  "
	if last {
		spine = prefix + "   "
	}
	for i, c := range node.Waiting {
		cLast := i == len(node.Waiting)-1
		conn := "├─ "
		if cLast {
			conn = "└─ "
		}
		b.out = append(b.out, listLine{text: b.m.chopRow(c, spine+conn, b.sel()), itemIdx: len(b.recipe)})
		b.recipe = append(b.recipe, c.Summary.ID)
		// Their own branches fork deeper still.
		b.branch(c, spine+"   ", cLast)
	}
}

// chopNote is what a tree row says about what the recipe waits on.
func chopNote(node core.TreeNode) string {
	var parts []string
	if len(node.Needs) > 0 {
		parts = append(parts, "needs "+strings.Join(node.Needs, ", "))
	}
	if len(node.Missing) > 0 {
		parts = append(parts, strings.Join(node.Missing, ", ")+" not in the plan")
	}
	if c := cycleNote(node); c != "" {
		parts = append(parts, c)
	}
	return strings.Join(parts, " — ")
}

// chopRow renders one recipe in the tree.
//
// The row says what the console already says about a recipe — its lifecycle, its
// proof count, its open tasting notes — and what it waits on.
func (m *Model) chopRow(node core.TreeNode, prefix string, selected bool) string {
	width := m.listColumn()
	note := chopNote(node)
	marks := proofMarks(node.Summary)
	tail := marks
	if note != "" {
		tail += " " + note
	}
	title := padCol(node.Summary.Title, titleColumn(width, lipgloss.Width(tail)))
	row := prefix + m.styFaint.Render(padRight(node.Summary.ID, 6)) +
		m.lifecycleStyle(string(node.Summary.Lifecycle)).Render(padRight(string(node.Summary.Lifecycle), 12)) +
		title + m.styFaint.Render(marks) + m.styDim.Render(strings.TrimPrefix(tail, marks))
	if selected {
		return m.stySelected.Render(padRight(row, width))
	}
	return row
}

// cycleNote is what a recipe in a dependency cycle says about itself.
//
// A member of the cycle cannot be executed at all. A recipe that merely waits on a
// member can be executed once the cycle is broken — the two are told apart, because
// the plan tells them apart.
func cycleNote(node core.TreeNode) string {
	switch {
	case node.Cycle:
		return "in a dependency cycle — the plan cannot be executed as written"
	case node.WaitsOnCycle:
		return "waits on a dependency cycle — the plan cannot be executed as written"
	}
	return ""
}

// logWindow is the slice of the journal the Log tab shows: the newest events that
// fit, newest first.
//
// The list, the selection and the detail all read this one slice, so the highlighted
// row is always the event the detail shows. The selection used to be clamped to the
// whole journal, so it could point at an event the list never renders — and the
// detail then showed an event that was not on the screen.
func (m *Model) logWindow() []core.JournalEntry {
	return m.journal.Recent(m.listAreaHeight() + 4)
}

// logLines renders the Log tab.
func (m *Model) logLines() []listLine {
	events := m.logWindow()
	out := make([]listLine, 0, len(events))
	for i, e := range events {
		out = append(out, listLine{text: "  " + m.logRow(e), itemIdx: i})
	}
	return out
}

// logRow renders one journal entry.
func (m *Model) logRow(e core.JournalEntry) string {
	at := e.At
	if len(at) >= 19 {
		at = at[11:19]
	}
	recipe := e.Recipe
	if recipe == "" {
		recipe = "·"
	}
	return m.styFaint.Render(padRight(at, 10)) + padRight(e.Kind, 11) +
		m.styFaint.Render(padRight(recipe, 5)) + m.styDim.Render(e.Actor)
}

// expoLines renders the Expo tab.
//
// The Expo is the orchestrator agent: it calls the tickets, routes them, and keeps the
// line moving. It is not built yet, and this tab says so — a tab that is *waiting*
// must not look like a tab that is failing.
//
// What it shows for now is the automation's vital signs: the rules the plan was last
// written under, the rules this build speaks, and the last thing that happened. The
// recipe counts, the triage and the tree are the other tabs' business.
func (m *Model) expoLines() []listLine {
	var out []listLine
	add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }

	add("")
	add(" " + m.stySec.Render("EXPO"))
	add("  " + m.styDim.Render("The orchestrator agent. It calls the tickets, routes"))
	add("  " + m.styDim.Render("them, and keeps the line moving."))
	add("")
	add("  " + m.styWarn.Render("Not built yet.") + m.styDim.Render(" This tab is where the chef"))
	add("  " + m.styDim.Render("will watch the automation work."))

	add("")
	add(" " + m.stySec.Render("VITAL SIGNS"))
	planVersion := core.PlanDomainVersion(m.journal)
	add("  " + m.styDim.Render("plan written under  ") + m.styFaint.Render(strconv.Itoa(planVersion)))
	add("  " + m.styDim.Render("this build speaks  ") + m.styFaint.Render(strconv.Itoa(core.DomainVersion)))
	if planVersion == core.DomainVersion {
		add("  " + m.styDim.Render("they agree"))
	} else {
		add("  " + m.styWarn.Render("they do not agree") +
			m.styDim.Render(" — the plan was written by a different build of the tool"))
	}

	add("")
	add(" " + m.stySec.Render("LAST EVENT"))
	if len(m.journal.Events) == 0 {
		add("  " + m.styDim.Render("the plan has no events"))
		return out
	}
	e := m.journal.Events[len(m.journal.Events)-1]
	at := e.At
	if len(at) >= 19 {
		at = at[11:19]
	}
	recipe := e.Recipe
	if recipe == "" {
		recipe = "·"
	}
	add("  " + m.styFaint.Render(padRight(e.Kind, 12)) + m.styFaint.Render(padRight(recipe, 6)) +
		m.styDim.Render(padRight(e.Actor, 10)) + m.styDim.Render(at))
	return out
}

// settingsLines renders the Settings tab.
func (m *Model) settingsLines() []listLine {
	var out []listLine
	add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }
	add("")
	add(" " + m.stySec.Render("PLAN"))
	add("  " + m.styDim.Render("root        ") + m.store.Root)
	add("  " + m.styDim.Render("recipes     ") + fmt.Sprint(len(m.recipes)))
	add("  " + m.styDim.Render("specials    ") + fmt.Sprint(len(m.specials)))
	add("  " + m.styDim.Render("prep notes  ") + fmt.Sprint(len(m.notes)))
	add("  " + m.styDim.Render("events      ") + fmt.Sprint(len(m.journal.Events)))
	add("")
	add(" " + m.stySec.Render("VIEW"))
	add("  " + m.styDim.Render("order       ") + OrderNames[int(m.order)])
	add("  " + m.styDim.Render("show served ") + fmt.Sprint(m.showServed))
	add("")
	add(" " + m.stySec.Render("THEME"))
	for i, p := range palettes() {
		marker := "  "
		if p.name == m.theme.name {
			marker = "● "
		}
		row := marker + padRight(p.name, 22)
		if i == m.selTheme {
			row = m.stySelected.Render(row)
		}
		out = append(out, listLine{text: " " + row, itemIdx: i})
	}
	add("")
	add(" " + m.stySec.Render("KEYS"))
	enterWhat := "open the recipe full screen"
	if m.tab == TabSettings {
		enterWhat = "apply the theme"
	}
	for _, k := range [][2]string{
		{"↑ ↓ / j k", "move"},
		{"pgup / pgdn", "page"},
		{"g / G", "first / last"},
		{"tab / shift+tab", "switch tab"},
		{"enter", enterWhat},
		{"esc", "back"},
		{"space", "the primary action"},
		{"n / N", "move over the open tasting notes"},
		{"1 2 3", "order: triage, lane, plain"},
		{"c", "show or hide served recipes"},
		{"r", "reload"},
		{"q", "quit"},
	} {
		add("  " + m.styKey.Render(padRight(k[0], 18)) + m.styDim.Render(k[1]))
	}
	return out
}

// proofMarks is the trailing marker column of a recipe row.
func proofMarks(s core.Summary) string {
	marks := ""
	if s.Proofs.Criteria > 0 {
		marks = fmt.Sprintf(" %d/%d", s.Proofs.Proven, s.Proofs.Criteria)
	}
	if s.Proofs.Rejected > 0 {
		marks += " ✗"
	}
	if s.Proofs.Unproven > 0 {
		marks += " ?"
	}
	if s.OpenFindings > 0 {
		marks += fmt.Sprintf(" %d note", s.OpenFindings)
	}
	if s.Readiness.Ready {
		marks += " ●"
	}
	return marks
}

// joinColumns puts two lists side by side.
func joinColumns(left, right []string, leftWidth int) string {
	n := len(left)
	if len(right) > n {
		n = len(right)
	}
	var b strings.Builder
	for i := 0; i < n; i++ {
		l := ""
		if i < len(left) {
			l = left[i]
		}
		r := ""
		if i < len(right) {
			r = right[i]
		}
		b.WriteString(padRight(l, leftWidth))
		if r != "" {
			b.WriteString(" " + r)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// truncate shortens s to n columns, ANSI-aware.
func truncate(s string, n int) string {
	if lipgloss.Width(s) <= n {
		return s
	}
	if n <= 1 {
		return "…"
	}
	// Keep the tail: titles share a prefix far more often than a suffix, so the tail
	// is the informative part.
	runes := []rune(stripANSI(s))
	if len(runes) <= n {
		return s
	}
	return "…" + string(runes[len(runes)-(n-1):])
}

// padCol puts s in a column of n, leaving one column of air after it so the next
// column never runs into it.
func padCol(s string, n int) string {
	return padRight(truncate(s, n), n+1)
}

// padRight pads s to n columns, ANSI-aware.
func padRight(s string, n int) string {
	w := lipgloss.Width(s)
	if w >= n {
		return s
	}
	return s + strings.Repeat(" ", n-w)
}

// plural picks the singular or plural form for a count.
func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// max is the larger of two ints.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// renderMarkdown renders prose for the terminal.
func renderMarkdown(md string, width int) string {
	if md == "" {
		return ""
	}
	if width < 20 {
		width = 20
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(styles.DarkStyle),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return md
	}
	out, err := r.Render(md)
	if err != nil {
		return md
	}
	return strings.TrimRight(out, "\n")
}
