package ui

// The hit map — where the console drew things.
//
// The mouse sends coordinates and the render is a string, so the render leaves a map
// behind. The map is built **by the render**: the renderer is the only thing that
// knows where it put the tabs, the rows and the detail, and a second function that
// re-derived them would be a second answer to one question — this console's
// standing failure.

// tabBarRow is the row the tabs are on: the header's third row.
//
// The header is four rows — the title, the rule, the tabs and the summary — and the
// tabs are the third of them.
const tabBarRow = headerRows - 2

// hitBox is a rectangle on the screen.
//
// The bounds are inclusive: a one-cell box has x0 == x1.
//
// The zero box is **not there**. A region that was never drawn must answer for
// nothing: a zero box that answered for one cell at the origin made the top-left
// corner of a screen with no detail pane look like it had one.
type hitBox struct {
	x0, x1 int
	y0, y1 int
	// there says the region is on screen. The zero box is not.
	there bool
}

// newHitBox is a region that is there.
func newHitBox(x0, x1, y0, y1 int) hitBox {
	return hitBox{x0: x0, x1: x1, y0: y0, y1: y1, there: true}
}

// contains reports whether the point is inside the box.
func (b hitBox) contains(x, y int) bool {
	if !b.there {
		return false
	}
	return x >= b.x0 && x <= b.x1 && y >= b.y0 && y <= b.y1
}

// empty reports whether the box holds nothing.
func (b hitBox) empty() bool {
	if !b.there {
		return true
	}
	return b.x1 < b.x0 || b.y1 < b.y0
}

// centre is the middle of the box.
func (b hitBox) centre() (int, int) {
	return (b.x0 + b.x1) / 2, (b.y0 + b.y1) / 2
}

// boardHit is one node the wide board placed.
//
// The render records it as it writes the node. The columns are the canvas's —
// x 0 is the left of the screen — and line is the row in the board's lines.
// A second walk that guessed the columns would be a second answer, and the
// click would land on the guess.
type boardHit struct {
	item   int
	x0, x1 int
	line   int
}

// hitRow is a list row: where it is, and which item it holds.
//
// The item is the index the tab's own selection uses. A header holds none, and says
// so with -1.
type hitRow struct {
	hitBox
	item int
}

// hitMap is where the last drawn frame put things.
//
// It is rebuilt on every render and kept on the model: a mouse event arrives
// **after** the frame was drawn, and an event is answered against the frame the
// chef is looking at.
type hitMap struct {
	// tabs are the tab labels, in tab order.
	tabs []hitBox
	// rows are the visible list rows, top to bottom.
	rows []hitRow
	// list is the list's region.
	list hitBox
	// detail is the detail pane's region.
	detail hitBox
	// banner are the refusal banner's rows.
	banner []hitBox
}

// tabAt returns the tab under the point, or -1.
func (h hitMap) tabAt(x, y int) int {
	for i, b := range h.tabs {
		if b.contains(x, y) {
			return i
		}
	}
	return -1
}

// rowAt returns the item under the point, and whether there is one.
//
// A header holds no item, and the second value says so.
func (h hitMap) rowAt(x, y int) (int, bool) {
	for _, r := range h.rows {
		if r.contains(x, y) {
			if r.item < 0 {
				return -1, false
			}
			return r.item, true
		}
	}
	return -1, false
}

// overList reports whether the point is over the list.
func (h hitMap) overList(x, y int) bool { return h.list.contains(x, y) }

// overDetail reports whether the point is over the detail pane.
func (h hitMap) overDetail(x, y int) bool { return h.detail.contains(x, y) }

// defaultPoint is where a mouse step goes when it names no position.
//
// The list is the console's default target — it is where the chef looks — and the
// detail is the fallback when the list is not on screen.
func (h hitMap) defaultPoint() (int, int) {
	if !h.list.empty() {
		return h.list.centre()
	}
	if !h.detail.empty() {
		return h.detail.centre()
	}
	return 0, 0
}

// hitMapOf records the boxes of a drawn frame.
//
// Every fact here comes from the render: the tab bar's boxes, the banner's rows, the
// list's rows and the detail's region. Nothing is re-derived — the map is what the
// renderer drew, not a second guess at it.
func (m *Model) hitMapOf(tabs []hitBox, bannerRows int, lines []listLine, start, height int, detail hitBox) hitMap {
	hits := hitMap{tabs: tabs, detail: detail}
	for i := 0; i < bannerRows; i++ {
		// The banner sits under the summary, and the filter line sits under the
		// summary too — so the banner starts lower when the line is open.
		row := headerRows + m.filterRows() + i
		hits.banner = append(hits.banner, newHitBox(0, m.width-1, row, row))
	}
	width := m.listColumn()
	// The wide board takes the screen, so its nodes sit past the list's usual
	// column. The boxes are the ones the render recorded while it placed them.
	if len(m.boardHits) > 0 && m.width > width {
		width = m.width
	}
	if height > 0 {
		hits.list = newHitBox(0, width-1, start, start+height-1)
	}
	if len(m.boardHits) > 0 {
		for _, n := range m.boardHits {
			if height <= 0 || n.line < m.offset || n.line >= m.offset+height {
				continue
			}
			y := start + n.line - m.offset
			hits.rows = append(hits.rows, hitRow{
				hitBox: newHitBox(n.x0, n.x1, y, y),
				item:   n.item,
			})
		}
		return hits
	}
	for i := 0; i < height; i++ {
		idx := m.offset + i
		item := -1
		if idx < len(lines) {
			item = lines[idx].itemIdx
		}
		hits.rows = append(hits.rows, hitRow{
			hitBox: newHitBox(0, width-1, start+i, start+i),
			item:   item,
		})
	}
	return hits
}
