package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/glamour/v2/styles"
	"charm.land/lipgloss/v2"

	"github.com/lthiagol/maitre/internal/brand"
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
//
// The mouse mode is set on **every** view the console draws. In bubbletea v2 the
// mode belongs to the view, not to the program: a view that does not ask for the
// mouse gets no mouse, and the terminal drops the events on the floor.
func (m *Model) View() tea.View {
	if m.err != nil {
		v := tea.NewView(m.styBad.Render("error: " + m.err.Error()))
		v.AltScreen = true
		v.WindowTitle = m.windowTitle()
		v.MouseMode = tea.MouseModeCellMotion
		// Nothing is drawn but the error: the map is empty, so no mouse event can
		// land on a tab that is not on screen.
		m.hits = hitMap{}
		return v
	}
	if !m.ready {
		v := tea.NewView(m.styDim.Render("opening the kitchen…"))
		v.AltScreen = true
		v.WindowTitle = m.windowTitle()
		v.MouseMode = tea.MouseModeCellMotion
		m.hits = hitMap{}
		return v
	}

	v := tea.NewView(m.renderLayout())
	v.AltScreen = true
	v.BackgroundColor = lipgloss.Color(m.colBg)
	v.WindowTitle = m.windowTitle()
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

// renderLayout draws the whole screen: header, the list, and a pinned footer.
func (m *Model) renderLayout() string {
	lines := m.listLines()
	m.scrollToSelection(lines)

	if m.mode == ModeDetail {
		return m.renderDetailFull()
	}
	// Home has no detail: it shows the plan's numbers, and the whole width is
	// the numbers' column. A pane beside it would hold nothing.
	//
	// The Board uses the width for the tree. At 160 and above the forks sit
	// to the right of the trunk, and a pane beside the tree would take those
	// columns. The detail opens on enter, and on a click.
	if m.tab == TabHome || m.tab == TabBoard {
		return m.renderSingle(lines)
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
	bar, tabs := m.tabBar()
	b.WriteString(bar + "\n")
	b.WriteString(m.summaryLine() + "\n")
	bannerRows := m.writeBanner(&b)

	detail := m.detailScrollRows(detailRows(m.detailLines(m.detailPaneWidth())))
	h := m.listAreaHeight()
	start := headerRows + bannerRows
	for i := 0; i < h; i++ {
		if i < len(detail) {
			b.WriteString(detail[i])
		}
		b.WriteString("\n")
	}
	// The detail owns the whole body: the list is not on screen.
	m.hits = m.hitMapOf(tabs, bannerRows, nil, start, 0,
		newHitBox(0, m.width-1, start, start+h-1))
	// The footer is the last row of the screen, on a row of its own. Do not trim:
	// the trailing newline is the padding that holds the footer down.
	return b.String() + m.footer()
}

// renderSingle draws the list across the whole terminal.
func (m *Model) renderSingle(lines []listLine) string {
	var b strings.Builder
	b.WriteString(m.header() + "\n")
	b.WriteString(m.ruleLine() + "\n")
	bar, tabs := m.tabBar()
	b.WriteString(bar + "\n")
	b.WriteString(m.summaryLine() + "\n")
	bannerRows := m.writeBanner(&b)

	h := m.listAreaHeight()
	start := headerRows + bannerRows
	for i := 0; i < h; i++ {
		idx := m.offset + i
		if idx < len(lines) {
			b.WriteString(lines[idx].text)
		}
		b.WriteString("\n")
	}

	m.hits = m.hitMapOf(tabs, bannerRows, lines, start, h, hitBox{})

	// The footer is the last row of the screen. Do not trim: the trailing newlines
	// are the padding that holds the footer down.
	return b.String() + m.footer()
}

// renderSplit draws the list and the detail side by side, both full height.
func (m *Model) renderSplit(lines []listLine) string {
	listWidth := m.listColumn()
	detailWidth := m.detailPaneWidth()

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
	// The pane scrolls here, like it does full screen: the offset is applied where
	// the detail is drawn.
	right := m.detailScrollRows(detailRows(m.detailLines(detailWidth)))
	if len(right) > h {
		right = right[:h]
	}
	for len(right) < h {
		right = append(right, "")
	}

	var b strings.Builder
	b.WriteString(m.header() + "\n")
	b.WriteString(m.ruleLine() + "\n")
	bar, tabs := m.tabBar()
	b.WriteString(bar + "\n")
	b.WriteString(m.summaryLine() + "\n")
	bannerRows := m.writeBanner(&b)
	b.WriteString(joinColumns(left, right, listWidth))

	// The detail sits beside the list: the list's column, one column of air, and
	// the rest of the terminal.
	start := headerRows + bannerRows
	m.hits = m.hitMapOf(tabs, bannerRows, lines, start, h,
		newHitBox(listWidth+1, m.width-1, start, start+h-1))
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
//
// The mark's compact form sits before the product name — one cell, so the header
// stays one line. The header is the one place the console has no room to draw.
//
// The project's name is a **name**, not a footnote: it wears the console's accent,
// where the product's name wears the ink. The two are different things, and they
// read differently.
func (m *Model) header() string {
	lead := m.styTitle.Render(" "+brand.Compact+" maitre") + m.styFaint.Render(" · ")
	// The header is one line. A name that does not fit is truncated to its tail,
	// or the wrap pushes the tabs — and the footer — down.
	budget := m.width - lipgloss.Width(lead)
	if budget < 1 {
		budget = 1
	}
	return lead + m.styAcc.Render(truncate(m.projectName(), budget))
}

// ruleLine is the rule under the header.
func (m *Model) ruleLine() string {
	return m.styRule.Render(strings.Repeat("─", max(10, m.width)))
}

// tabBar is the row of tabs.
//
// It returns the row and the boxes the labels occupy. The mouse is answered against
// the frame, and the renderer is the only thing that knows where the labels are —
// so the renderer says.
func (m *Model) tabBar() (string, []hitBox) {
	var parts []string
	var boxes []hitBox
	col := 0
	for i, name := range TabNames {
		label := "  " + name
		if Tab(i) == m.tab {
			label = "▸ " + name
			parts = append(parts, m.styTabActive.Render(label))
		} else {
			parts = append(parts, m.styTabIdle.Render(label))
		}
		w := lipgloss.Width(label)
		boxes = append(boxes, newHitBox(col, col+w-1, tabBarRow, tabBarRow))
		// The labels are joined by two spaces.
		col += w + 2
	}
	return strings.Join(parts, "  "), boxes
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

// writeBanner writes the refusal banner, if there is one, and returns how many rows
// it took.
//
// The banner sits between the summary and the list, so the list starts lower when it
// is there.
func (m *Model) writeBanner(b *strings.Builder) int {
	lines := m.bannerLines()
	for _, l := range lines {
		b.WriteString(l + "\n")
	}
	return len(lines)
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
	case TabHome:
		return "the plan at a glance · nothing to press here"
	case TabRecipes:
		text := fmt.Sprintf("%d %s · %d need you",
			len(m.recipes), plural(len(m.recipes), "recipe", "recipes"), m.countNeedsYou())
		if m.order == OrderTriage {
			text += "   " + m.servedCheckbox()
		}
		return text
	case TabBacklog:
		text := fmt.Sprintf("%d %s", len(m.items), plural(len(m.items), "item", "items"))
		return text
	case TabReady:
		recipes := core.Adoptable(m.recipes)
		tasks := openSmallTasks(m.items)
		return fmt.Sprintf("%d adoptable · %d %s",
			len(recipes), len(tasks), plural(len(tasks), "small task", "small tasks"))
	case TabBoard:
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
	case m.tab == TabHome:
		// The Home tab shows numbers and nothing to press: there is no item under
		// the cursor, so the footer promises no move — not `space`, not `enter`,
		// not the cursor.
		keys = append(keys, key{"tab", "switch"})
	case m.tab == TabSettings:
		// The Settings tab has no recipe to open: `enter` applies the theme.
		keys = append(keys, key{"↑↓", "move"}, key{"tab", "switch"}, key{"enter", "apply theme"})
	default:
		keys = append(keys, key{"↑↓", "move"}, key{"tab", "switch"})
		// The footer never promises a key that does nothing. `enter` opens the
		// detail, and there is one only where the pane has something to show —
		// the same answer the key reads.
		if m.hasDetail() {
			keys = append(keys, key{"enter", "detail"})
		}
	}
	if label := m.primaryActionLabel(); label != "" {
		keys = append(keys, key{"space", label})
	}
	// The backlog's own keys. `d` is offered only when the domain says the item may
	// still be dismissed — a key the domain would refuse is a key the footer must
	// not name.
	if m.tab == TabBacklog {
		if m.canDismissSelected() {
			keys = append(keys, key{"d", "dismiss"})
		}
		keys = append(keys, key{"a", "add item"})
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
	case TabRecipes:
		rec, ok := m.currentRecipe()
		if !ok {
			return ""
		}
		return core.ChefMoveOf(*rec).Label()
	case TabReady:
		// The footer asks the domain. An adoptable recipe has no move: the
		// console does not claim.
		item, ok := m.selectedReadyTask()
		if !ok {
			return ""
		}
		return core.BacklogMoveOf(item).Primary.Label()
	case TabBacklog:
		item, ok := m.selectedBacklogItem()
		if !ok {
			return ""
		}
		// The domain's answer, not the console's: the footer names the move the
		// current selection has — shape it, promote it, or nothing — and a
		// settled item has none, so the console must not promise a key there.
		return core.BacklogMoveOf(item).Primary.Label()
	}
	return ""
}

// listLines renders the current tab's list, marking which lines hold a selectable item.
func (m *Model) listLines() []listLine {
	// The wide board records a hit per node as it draws. Any other frame has
	// none, or a click on the Recipes tab would land on a node from the Board.
	m.boardHits = nil
	switch m.tab {
	case TabBacklog:
		return m.backlogLines()
	case TabReady:
		return m.readyLines()
	case TabBoard:
		return m.boardTree().rows
	case TabLog:
		return m.logLines()
	case TabHome:
		return m.homeLines()
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
	marks := evidenceMarks(s)
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
//
// The two groups are **derived** from the item's kind — the kind is the only truth,
// and a group is never stored.
func (m *Model) backlogLines() []listLine {
	var out []listLine
	ideas, tasks := core.BacklogGroups(m.items)
	kindWidth := m.backlogKindWidth()

	section := func(name string, items []core.BacklogItem, offset int) {
		if len(items) == 0 {
			return
		}
		out = append(out, listLine{text: m.sectionHeader(name, len(items), m.listColumn()), itemIdx: -1})
		out = append(out, listLine{text: m.styRule.Render(strings.Repeat("─", m.listColumn())), itemIdx: -1})
		for i, it := range items {
			idx := offset + i
			out = append(out, listLine{
				text:    m.backlogRow(idx, it.ID, it.Title, it.Kind, it.Status, kindWidth, m.selBacklog == idx),
				itemIdx: idx,
			})
		}
	}
	section("IDEAS", ideas, 0)
	section("SMALL TASKS", tasks, len(ideas))
	return out
}

// backlogKindWidth is the width of the kind column.
//
// It is the width of the widest kind the tab actually shows, so the column lines up
// whether the tab holds ideas or small tasks.
func (m *Model) backlogKindWidth() int {
	width := 0
	for _, it := range m.items {
		if w := lipgloss.Width(it.Kind); w > width {
			width = w
		}
	}
	return width
}

// backlogRow renders one backlog item.
//
// The row shows the item's kind: the kind is the only truth about the item's shape,
// and the group it is filed under is derived from it.
func (m *Model) backlogRow(index int, id, title, kind, status string, kindWidth int, selected bool) string {
	width := m.listColumn()
	tail := kindWidth + lipgloss.Width(status)
	row := "  " + padRight(id, 6) + padCol(title, titleColumn(width, tail)) +
		m.styFaint.Render(padRight(kind, kindWidth+1)) + m.styDim.Render(status)
	if selected {
		return m.stySelected.Render(padRight(row, width))
	}
	return row
}

// readyLines renders the Ready tab.
//
// Both blocks are always drawn. An empty block is an empty list, present.
func (m *Model) readyLines() []listLine {
	recipes := core.Adoptable(m.recipes)
	tasks := openSmallTasks(m.items)
	var out []listLine
	item := 0

	header := func(name string, n int) {
		out = append(out, listLine{text: m.sectionHeader(name, n, m.listColumn()), itemIdx: -1})
		out = append(out, listLine{text: m.styRule.Render(strings.Repeat("─", m.listColumn())), itemIdx: -1})
	}
	header("ADOPTABLE", len(recipes))
	for _, r := range recipes {
		out = append(out, listLine{text: m.readyRecipeRow(r, m.selReady == item), itemIdx: item})
		item++
	}
	header("SMALL TASKS", len(tasks))
	for _, it := range tasks {
		out = append(out, listLine{text: m.readyTaskRow(it, m.selReady == item), itemIdx: item})
		item++
	}
	return out
}

// readyRecipeRow renders one adoptable recipe.
func (m *Model) readyRecipeRow(r core.Recipe, selected bool) string {
	width := m.listColumn()
	title := padCol(r.Title, titleColumn(width, 0))
	row := "  " + m.styFaint.Render(padRight(r.ID, 6)) + title
	if selected {
		return m.stySelected.Render(padRight(row, width))
	}
	return row
}

// readyTaskRow renders one open small task.
//
// The words are the two steps after the promotion, so the row is not read as
// one keystroke.
func (m *Model) readyTaskRow(it core.BacklogItem, selected bool) string {
	width := m.listColumn()
	note := "promote → sign off → run"
	title := padCol(it.Title, titleColumn(width, lipgloss.Width(note)))
	row := "  " + m.styFaint.Render(padRight(it.ID, 6)) + title + m.styDim.Render(note)
	if selected {
		return m.stySelected.Render(padRight(row, width))
	}
	return row
}

// boardBuilder builds the Board's rows.
//
// The tree is drawn as a trunk with branches: the trunk is the ready list, and each
// branch holds the recipes that fork at the recipe above it.
type boardBuilder struct {
	m      *Model
	out    []listLine
	recipe []string
}

// trunk draws a trunk recipe and the branch that forks at it.
func (b *boardBuilder) trunk(node core.TreeNode) {
	b.out = append(b.out, listLine{text: b.m.boardRow(node, "● ", b.sel()), itemIdx: len(b.recipe)})
	b.recipe = append(b.recipe, node.Summary.ID)
	b.branch(node, "", true)
}

// sel reports whether the row being drawn is the selected one.
func (b *boardBuilder) sel() bool { return len(b.recipe) == b.m.selBoard }

// orphan draws a recipe whose dependency is not in the plan.
//
// It cannot fork at anything, and the plan is broken: it is marked, and it says which
// dependency is missing.
func (b *boardBuilder) orphan(node core.TreeNode) {
	b.out = append(b.out, listLine{text: b.m.boardOrphanRow(node, b.sel()), itemIdx: len(b.recipe)})
	b.recipe = append(b.recipe, node.Summary.ID)
	b.branch(node, "", true)
}

// stuck draws a recipe a dependency cycle keeps out of the tree.
//
// The plan cannot be executed as written, and the tree says so rather than dropping
// the recipe.
func (b *boardBuilder) stuck(node core.TreeNode) {
	b.out = append(b.out, listLine{text: b.m.boardStuckRow(node, b.sel()), itemIdx: len(b.recipe)})
	b.recipe = append(b.recipe, node.Summary.ID)
	b.branch(node, "", true)
}

// boardOrphanRow renders a recipe whose dependency is not in the plan.
//
// The plan is broken, and the marker says so: the row carries the warning, not the
// whole line.
func (m *Model) boardOrphanRow(node core.TreeNode, selected bool) string {
	return m.fitTreeRow(m.styWarn.Render("⚠ "), node, selected)
}

// boardStuckRow renders a recipe a dependency cycle keeps out of the tree.
//
// The plan cannot be executed as written, and the row says so.
func (m *Model) boardStuckRow(node core.TreeNode, selected bool) string {
	return m.fitTreeRow(m.styWarn.Render("⚠ "), node, selected)
}

// branch draws the group of recipes that fork at a node.
//
// The heading names the branch's subject — the dependency they all wait on — and the
// rows under it are the recipes that fork there.
func (b *boardBuilder) branch(node core.TreeNode, prefix string, last bool) {
	if len(node.Waiting) == 0 {
		return
	}
	connector := "├─ "
	if last {
		connector = "└─ "
	}
	b.out = append(b.out, listLine{
		text:    b.m.fitBoardLine(" " + b.m.styDim.Render(prefix+connector+"needs "+node.Summary.ID)),
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
		b.out = append(b.out, listLine{text: b.m.boardRow(c, spine+conn, b.sel()), itemIdx: len(b.recipe)})
		b.recipe = append(b.recipe, c.Summary.ID)
		// Their own branches fork deeper still.
		b.branch(c, spine+"   ", cLast)
	}
}

// boardNote is what a tree row says about what the recipe waits on.
func boardNote(node core.TreeNode) string {
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

// boardRow renders one recipe in the tree.
//
// The row says what the console already says about a recipe — its lifecycle, its
// evidence count, its open tasting notes — and what it waits on.
func (m *Model) boardRow(node core.TreeNode, prefix string, selected bool) string {
	return m.fitTreeRow(prefix, node, selected)
}

// fitTreeRow lays one recipe out so the row is one screen line.
//
// A row wider than the terminal wraps, and the wrap pushes the footer off
// the last row. Rows that already fit are unchanged. A row that does not
// keeps the title, gives up title padding, and then the tail of the note.
func (m *Model) fitTreeRow(prefix string, node core.TreeNode, selected bool) string {
	marks := evidenceMarks(node.Summary)
	note := boardNote(node)
	titleW := titleColumn(m.listColumn(), lipgloss.Width(treeTail(marks, note)))
	row := m.renderTreeRow(prefix, node, marks, note, titleW, selected)
	limit := m.width
	if limit < 1 || lipgloss.Width(row) <= limit {
		return row
	}
	natural := lipgloss.Width(node.Summary.Title)
	// The title stays whole when the row was only over because the note is
	// long. Title padding is what gives way when the title already fits.
	if titleW < natural {
		titleW = natural
		row = m.renderTreeRow(prefix, node, marks, note, titleW, selected)
	}
	for lipgloss.Width(row) > limit && titleW > natural {
		titleW--
		row = m.renderTreeRow(prefix, node, marks, note, titleW, selected)
	}
	for lipgloss.Width(row) > limit && note != "" {
		next := truncateHead(note, lipgloss.Width(note)-1)
		if next == note || next == "" || next == "…" {
			note = ""
		} else {
			note = next
		}
		row = m.renderTreeRow(prefix, node, marks, note, titleW, selected)
	}
	for lipgloss.Width(row) > limit && titleW > 0 {
		titleW--
		row = m.renderTreeRow(prefix, node, marks, note, titleW, selected)
	}
	if lipgloss.Width(row) > limit {
		return truncateHeadStyled(row, limit)
	}
	return row
}

// renderTreeRow is one recipe row at a chosen title width.
func (m *Model) renderTreeRow(prefix string, node core.TreeNode, marks, note string, titleW int, selected bool) string {
	tail := treeTail(marks, note)
	row := prefix + m.styFaint.Render(padRight(node.Summary.ID, 6)) +
		m.lifecycleStyle(string(node.Summary.Lifecycle)).Render(padRight(string(node.Summary.Lifecycle), 12)) +
		padCol(node.Summary.Title, titleW) + m.styFaint.Render(marks) +
		m.styDim.Render(strings.TrimPrefix(tail, marks))
	if selected {
		return m.stySelected.Render(padRight(row, m.listColumn()))
	}
	return row
}

// treeTail is the marks and the note, in the order the row draws them.
func treeTail(marks, note string) string {
	if note == "" {
		return marks
	}
	return marks + " " + note
}

// fitBoardLine keeps one Board line inside the terminal.
//
// The cut keeps the start — the connector, the id — and marks the end.
func (m *Model) fitBoardLine(s string) string {
	if m.width < 1 || lipgloss.Width(s) <= m.width {
		return s
	}
	return truncateHeadStyled(s, m.width)
}

// truncateHeadStyled shortens a styled line to n columns, keeping the start.
func truncateHeadStyled(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if n == 1 {
		return "…"
	}
	var b strings.Builder
	visible := 0
	runes := []rune(s)
	for i := 0; i < len(runes); {
		if runes[i] == 0x1b {
			j := i + 1
			if j < len(runes) && runes[j] == '[' {
				j++
				for j < len(runes) && (runes[j] < 0x40 || runes[j] > 0x7e) {
					j++
				}
				if j < len(runes) {
					j++
				}
			} else if j < len(runes) {
				j++
			}
			b.WriteString(string(runes[i:j]))
			i = j
			continue
		}
		rw := lipgloss.Width(string(runes[i]))
		if visible+rw > n-1 {
			break
		}
		b.WriteRune(runes[i])
		visible += rw
		i++
	}
	b.WriteString("…\x1b[0m")
	return b.String()
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

// settingsLines renders the Settings tab.
func (m *Model) settingsLines() []listLine {
	var out []listLine
	add := func(text string) { out = append(out, listLine{text: text, itemIdx: -1}) }
	add("")
	add(" " + m.stySec.Render("PLAN"))
	// The root is one row. A path that does not fit is truncated to its tail, the
	// same rule as the plan row: a wrap pushes the footer off the bottom.
	rootLabel := "root        "
	budget := m.width - len("  ") - len(rootLabel)
	if budget < 1 {
		budget = 1
	}
	add("  " + m.styDim.Render(rootLabel) + truncate(humanPath(m.store.Root), budget))
	add("  " + m.styDim.Render("recipes     ") + fmt.Sprint(len(m.recipes)))
	add("  " + m.styDim.Render("backlog     ") + fmt.Sprint(len(m.items)))
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

// evidenceMarks is the trailing marker column of a recipe row.
func evidenceMarks(s core.Summary) string {
	marks := ""
	if s.Evidence.Criteria > 0 {
		marks = fmt.Sprintf(" %d/%d", s.Evidence.Proven, s.Evidence.Criteria)
	}
	if s.Evidence.Rejected > 0 {
		marks += " ✗"
	}
	if s.Evidence.Unproven > 0 {
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
