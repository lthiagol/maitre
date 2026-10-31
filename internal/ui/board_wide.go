package ui

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// wideGap sits between a recipe and the fork to its right.
//
// The narrow tree's connectors are ├─ and └─. This one runs left to right,
// so the trunk's column stays left of the fork's.
const wideGap = " ── "

// wideKind is how a node is marked.
//
// A trunk recipe leads the row. A fork sits to its right. An orphan and a
// recipe a cycle cannot reach are still drawn, and the mark says so — the
// same mark the vertical tree uses for them.
type wideKind int

const (
	wideTrunk wideKind = iota
	wideFork
	wideWarn
)

// wideCell is one node, placed.
type wideCell struct {
	row, depth, item int
	node             core.TreeNode
	kind             wideKind
}

// wideBuilder lays the tree out left to right.
//
// A node shares its row with its first fork, and further forks stack under
// that one, still in the column to the right. The walk appends recipes in
// the same order as the vertical tree, so the keyboard still moves down it.
type wideBuilder struct {
	m        *Model
	canvas   int
	row      int
	maxDepth int
	recipe   []string
	cells    []wideCell
}

// boardTreeWide draws the tree left to right.
//
// The canvas is the terminal's row. The Board takes the width, so a fork has
// room to sit to the right of the recipe it waits on. A node that does not
// fit is truncated, and the line stays one row, so the footer stays the last
// row of the screen.
func (m *Model) boardTreeWide(tree core.Tree) board {
	w := &wideBuilder{m: m, canvas: m.rowWidth()}
	if w.canvas < 1 {
		w.canvas = 1
	}
	for _, root := range tree.Roots {
		w.walk(root, 0, wideTrunk)
	}
	for _, orphan := range tree.Orphans {
		w.walk(orphan, 0, wideWarn)
	}
	for _, stuck := range tree.Stuck {
		w.walk(stuck, 0, wideWarn)
	}
	return board{rows: w.lines(), recipe: w.recipe}
}

// walk places a node and the forks that wait on it.
func (w *wideBuilder) walk(node core.TreeNode, depth int, kind wideKind) {
	if depth > w.maxDepth {
		w.maxDepth = depth
	}
	item := len(w.recipe)
	w.recipe = append(w.recipe, node.Summary.ID)
	w.cells = append(w.cells, wideCell{
		row: w.row, depth: depth, item: item, node: node, kind: kind,
	})
	if len(node.Waiting) == 0 {
		w.row++
		return
	}
	for _, child := range node.Waiting {
		w.walk(child, depth+1, wideFork)
	}
}

// lines renders the placed nodes.
//
// Columns share the canvas. The gap is what is left between them, and a node
// wider than its column is cut so the row cannot wrap.
func (w *wideBuilder) lines() []listLine {
	cols := w.maxDepth + 1
	gap, gapW, nodeW := wideColumns(w.canvas, cols)
	byRow := map[int]map[int]wideCell{}
	for _, c := range w.cells {
		if byRow[c.row] == nil {
			byRow[c.row] = map[int]wideCell{}
		}
		byRow[c.row][c.depth] = c
	}
	out := make([]listLine, w.row)
	for r := 0; r < w.row; r++ {
		var b strings.Builder
		itemIdx := -1
		for d := 0; d < cols; d++ {
			if d > 0 {
				if _, ok := byRow[r][d]; ok && gapW > 0 {
					b.WriteString(w.m.styDim.Render(gap))
				} else {
					b.WriteString(strings.Repeat(" ", gapW))
				}
			}
			c, ok := byRow[r][d]
			if !ok {
				b.WriteString(strings.Repeat(" ", nodeW))
				continue
			}
			if c.item == w.m.selBoard {
				itemIdx = c.item
			}
			b.WriteString(w.m.wideNode(c, nodeW))
		}
		out[r] = listLine{text: b.String(), itemIdx: itemIdx}
	}
	return out
}

// wideColumns splits the canvas into node columns and the gaps between them.
//
// The line has to fit. A tree deeper than the screen is squeezed until each
// node is one column and the gaps are gone; the node text is then truncated.
func wideColumns(canvas, cols int) (gap string, gapW, nodeW int) {
	if cols < 1 {
		cols = 1
	}
	if canvas < 1 {
		canvas = 1
	}
	gap = wideGap
	gapW = lipgloss.Width(gap)
	gaps := cols - 1
	if gaps > 0 && gapW*gaps >= canvas {
		gap = ""
		gapW = 0
	}
	remain := canvas - gapW*gaps
	if cols > remain {
		gap = ""
		gapW = 0
		remain = canvas
	}
	nodeW = remain / cols
	if nodeW < 1 {
		nodeW = 1
	}
	for nodeW > 1 && nodeW*cols+gapW*(cols-1) > canvas {
		nodeW--
	}
	return gap, gapW, nodeW
}

// wideNode renders one recipe in its column.
//
// The title leads. What the recipe waits on follows it, when the column has
// room. A column that is too narrow keeps the start of the title and cuts
// the rest: the node stays one line.
func (m *Model) wideNode(c wideCell, width int) string {
	prefix := "● "
	if c.kind == wideFork {
		prefix = ""
	}
	if c.kind == wideWarn {
		prefix = "⚠ "
	}
	plain := fitBoardNode(prefix, c.node.Summary.Title, boardNote(c.node), width)
	body := padRight(plain, width)
	if c.item == m.selBoard {
		return m.stySelected.Render(body)
	}
	return body
}

// fitBoardNode fits a node's text into width columns.
//
// The title stays. The note — the dependency it waits on — follows when it
// fits, and is cut before the title is. The cut keeps the start, so a long
// title is still the title.
func fitBoardNode(prefix, title, note string, width int) string {
	if width <= 0 {
		return ""
	}
	lead := prefix
	if lipgloss.Width(lead) > width {
		return truncateHead(lead, width)
	}
	room := width - lipgloss.Width(lead)
	sep := ""
	if note != "" {
		sep = " · "
	}
	if lipgloss.Width(title)+lipgloss.Width(sep)+lipgloss.Width(note) <= room {
		return lead + title + sep + note
	}
	if lipgloss.Width(title) <= room {
		spare := room - lipgloss.Width(title)
		if note != "" && spare > lipgloss.Width(sep) {
			return lead + title + sep + truncateHead(note, spare-lipgloss.Width(sep))
		}
		return lead + title
	}
	return lead + truncateHead(title, room)
}

// truncateHead shortens s to n columns, keeping the start.
//
// The console's other truncate keeps the tail, because titles share a prefix.
// A board node is the title, so the start is the part the chef reads.
func truncateHead(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= n {
		return s
	}
	if n == 1 {
		return "…"
	}
	r := []rune(s)
	for len(r) > 0 && lipgloss.Width(string(r)) > n-1 {
		r = r[:len(r)-1]
	}
	return string(r) + "…"
}
