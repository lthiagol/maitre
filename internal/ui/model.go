package ui

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// Tab is a top-level view.
type Tab int

const (
	// TabHome is where the console lands: the plan's numbers, and nothing to press.
	TabHome Tab = iota
	TabRecipes
	TabBacklog
	TabReady
	TabBoard
	TabLog
	TabSettings
)

// TabNames are the tab labels, in order.
//
// The Home tab is first: it is the landing place, and it holds the plan's numbers.
// Ready sits next to the Board: Ready is what an agent could claim now, and the
// Board is the plan's shape.
var TabNames = []string{"Home", "Recipes", "Backlog", "Ready", "Board", "Log", "Settings"}

// Label is the tab's name.
func (t Tab) Label() string { return TabNames[int(t)] }

// Order is how the recipe list is sorted.
type Order int

const (
	// OrderTriage puts what needs the chef first.
	OrderTriage Order = iota
	// OrderLane groups by the recipe's lifecycle lane.
	OrderLane
	// OrderPlain is by id.
	OrderPlain
)

// OrderNames are the order labels.
var OrderNames = []string{"triage", "lane", "plain"}

// Mode is whether the detail pane is split or full screen.
type Mode int

const (
	ModeSplit Mode = iota
	ModeDetail
)

// reloadMsg carries a reloaded plan.
type reloadMsg struct {
	recipes []core.Recipe
	items   []core.BacklogItem
	journal core.Journal
	// disk is the fingerprint of the plan as it was before it was read.
	disk string
	err  error
}

// pollMsg asks the console to look at the disk.
type pollMsg struct{}

// PollInterval is how often the live console looks for a change on disk.
//
// A second is short enough that a change feels immediate, and long enough that the
// scan is free: a stat per file, twenty files, once a second.
const PollInterval = time.Second

// listLine is one rendered line of the list.
//
// itemIdx is the index of the selectable item the line holds, or -1 for a header. The
// list scrolls by line, but the selection is by item — so a section header can scroll
// away while its items stay selected.
type listLine struct {
	text    string
	itemIdx int
}

// Model is the console's state.
type Model struct {
	store *core.Store
	ops   *core.Operations

	// theme is the console's palette and the styles derived from it. It is on the
	// Model because the chef can pick one in the Settings tab, and a theme that
	// lives in the package cannot be swapped at runtime.
	theme

	tab        Tab
	order      Order
	mode       Mode
	showServed bool

	recipes  []core.Recipe
	sections []core.TriageSection

	items []core.BacklogItem

	journal core.Journal

	selRecipes int
	selBacklog int
	// backlogID is the item the backlog cursor names.
	//
	// selBacklog is that item's row. A reload resolves the id to the row, so a
	// shaping that moves an idea into Small tasks does not leave the cursor on
	// the idea that slid into the old index.
	backlogID string
	// backlogFollow is the id a write asked the cursor to land on.
	// backlogLeave is an id the cursor must not name after the next reload:
	// a dismissal and a promotion take the item out of the chef's hands, and
	// the cursor sits on the nearest remaining row.
	backlogFollow string
	backlogLeave  string
	selReady      int
	// readyID is the row the Ready cursor names.
	// readyLeave is an id a promotion took off the tab: the cursor sits on the
	// nearest remaining row, and it does not stay on the id that left.
	readyID    string
	readyLeave string
	selBoard   int
	selLog     int
	// selTheme is the cursor over the theme list on the Settings tab.
	selTheme int
	// selNote is the cursor over the selected recipe's open tasting notes.
	selNote int

	// offset is the first visible line of the list.
	offset int

	// hits is where the last drawn frame put things.
	//
	// It is a fact about the screen, not about the plan, so it lives here and not in
	// the domain. The mouse event arrives after the frame was drawn, so the map is
	// kept between renders and rebuilt by the next one.
	hits hitMap
	// boardHits is where the wide board put its nodes, in that frame.
	//
	// The render fills it as it places the nodes, and hitMapOf copies it onto
	// the screen. Empty means the frame is the vertical tree, whose rows are
	// the hits.
	boardHits []boardHit

	// detailScroll is how far the detail pane has been scrolled.
	detailScroll int

	width  int
	height int
	ready  bool
	err    error

	// project is the name the console shows for this session, when one was given.
	//
	// It is display only: nothing is written and the plan keeps its own name.
	// Empty means the plan's own name.
	project string

	// refusal is the banner shown after a refused action.
	refusal string
	// status is a transient one-line message.
	status string

	// prompt is the prompt asking the chef to decide, while it is open. It carries a
	// verdict on a recipe or an action on a tasting note.
	prompt *chefPrompt

	// live is whether the console watches the disk. The snapshot and the script
	// paths are not interactive programs, and a poll loop in one of them would keep
	// it alive.
	live bool
	// pollInterval is how often the console looks for a change on disk.
	pollInterval time.Duration
	// disk is the fingerprint of the plan as the console last read it.
	disk string
}

// New builds the console model over a store.
//
// It does not watch the disk: the snapshot and the script paths are not interactive
// programs, and a poll loop in one of them would keep it alive. The interactive
// program uses NewLive.
func New(store *core.Store, opts ...Option) *Model {
	return newConsole(store, "console", opts...)
}

// NewLive builds the console model over a store, and it watches the disk.
//
// The chef sees a change on disk without pressing `r`: the console stats the plan's
// files once a second, and reloads when one of them has changed.
func NewLive(store *core.Store, opts ...Option) *Model {
	m := newConsole(store, "console", opts...)
	m.live = true
	return m
}

// NewScripted builds a console whose writes are stamped as scripted.
//
// The console is the chef's surface: a human pressing space in it is the chef
// deciding. A script pressing space for them is the chef's decision mediated by a
// machine — which the journal should be able to show.
func NewScripted(store *core.Store, opts ...Option) *Model {
	return newConsole(store, "console-script", opts...)
}

func newConsole(store *core.Store, session string, opts ...Option) *Model {
	m := &Model{
		store:      store,
		ops:        core.NewOperations(store, core.AuditContext{Actor: "chef", Session: session, Role: "chef"}),
		showServed: true,
		// Kitchen is the default: a console nobody told otherwise looks exactly
		// as it did before the theme was the chef's to choose.
		theme:        newTheme(kitchen()),
		pollInterval: PollInterval,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Init loads the plan.
func (m *Model) Init() tea.Cmd { return m.reload() }

// poll returns the command that asks the console to look at the disk again.
//
// A console that is not live has no poll: the snapshot and the script paths are not
// interactive programs, and a poll loop in one of them would keep it alive.
func (m *Model) poll() tea.Cmd {
	if !m.live {
		return nil
	}
	return tea.Tick(m.pollInterval, func(time.Time) tea.Msg { return pollMsg{} })
}

// checkDisk looks for a change on disk, and reloads the plan if there is one.
//
// The check is a stat per file: cheap enough to run every second, and the plan is
// read only when something has actually changed.
func (m *Model) checkDisk() (tea.Model, tea.Cmd) {
	fp, err := m.store.DiskFingerprint()
	if err != nil {
		// A plan that cannot be stat'ed is a plan that cannot be read either. The
		// chef's `r` still works, and the next poll tries again.
		return m, m.poll()
	}
	if fp != m.disk {
		return m, m.reload()
	}
	return m, m.poll()
}

// reload reads the plan and returns the command that does it.
//
// The command stats the plan before it reads it, so the console knows what it read:
// a change that lands while the plan is being read is caught by the next poll
// rather than missed.
func (m *Model) reload() tea.Cmd {
	return func() tea.Msg {
		disk, _ := m.store.DiskFingerprint()
		recipes, err := m.store.ListRecipes()
		if err != nil {
			return reloadMsg{err: err}
		}
		items, err := m.store.ListBacklogItems()
		if err != nil {
			return reloadMsg{err: err}
		}
		journal, err := m.store.LoadJournal()
		if err != nil {
			return reloadMsg{err: err}
		}
		// Everything here was written by an agent. Sanitise before it reaches the
		// screen, so a newline in a title cannot corrupt the layout.
		sanitizePlan(recipes, items, journal)
		return reloadMsg{recipes: recipes, items: items, journal: journal, disk: disk}
	}
}

// Update handles a message.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.ready = true
		return m, nil

	case reloadMsg:
		if msg.err != nil {
			m.err = msg.err
			// The plan is not what the console has, so the fingerprint stays: the
			// next poll tries again.
			return m, m.poll()
		}
		m.applyReload(msg)
		m.err = nil
		m.disk = msg.disk
		return m, m.poll()

	case pollMsg:
		return m.checkDisk()

	case tea.KeyPressMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)
	}
	return m, nil
}

// applyReload puts a reloaded plan into the model.
//
// **One place knows the model's shape.** The live console and the static render both
// go through here, so a field the render carries and the live path forgets cannot
// happen — a console that shows the chef a plan it is not looking at is worse than a
// console that crashes.
func (m *Model) applyReload(msg reloadMsg) {
	follow, leave := m.backlogFollow, m.backlogLeave
	m.backlogFollow, m.backlogLeave = "", ""
	readyLeave := m.readyLeave
	m.readyLeave = ""
	// The row on screen is the item the chef was on, unless the write named
	// another — a new id, or one that has left.
	if follow == "" && leave == "" {
		if item, ok := m.selectedBacklogItem(); ok {
			m.backlogID = item.ID
		}
	}
	if readyLeave == "" {
		if e, ok := m.currentReadyEntry(); ok {
			m.readyID = e.id
		}
	}
	m.recipes = msg.recipes
	m.items = msg.items
	m.journal = msg.journal
	m.sections = m.orderedSections()
	m.placeBacklog(follow, leave)
	m.placeReady(readyLeave)
	m.clampSelections()
}

// orderedSections applies the current ordering.
func (m *Model) orderedSections() []core.TriageSection {
	switch m.order {
	case OrderTriage:
		return core.TriageSections(m.recipes)
	case OrderLane:
		return m.sectionsByLane()
	default:
		return m.sectionsPlain()
	}
}

// sectionsByLane groups by the lifecycle lane.
func (m *Model) sectionsByLane() []core.TriageSection {
	byLane := map[core.Lane][]core.Recipe{}
	for _, r := range m.recipes {
		l := core.LaneOf(r)
		byLane[l] = append(byLane[l], r)
	}
	var out []core.TriageSection
	for _, l := range core.AllLanes {
		items := byLane[l]
		if len(items) == 0 {
			continue
		}
		sec := core.TriageSection{Name: l.Label()}
		for i := range items {
			sec.Items = append(sec.Items, core.SummaryOf(items[i], m.recipes))
		}
		out = append(out, sec)
	}
	return out
}

// sectionsPlain is a single section, by id.
func (m *Model) sectionsPlain() []core.TriageSection {
	sec := core.TriageSection{Name: "All recipes"}
	for i := range m.recipes {
		sec.Items = append(sec.Items, core.SummaryOf(m.recipes[i], m.recipes))
	}
	if len(sec.Items) == 0 {
		return nil
	}
	return []core.TriageSection{sec}
}

// visibleRows is the flat list of recipes the selection moves through.
func (m *Model) visibleRows() []core.Summary {
	var out []core.Summary
	for _, s := range m.sections {
		if s.Name == core.TriageServed && !m.showServed {
			continue
		}
		out = append(out, s.Items...)
	}
	return out
}

// currentSummary is the selected recipe's row on the Recipes tab.
func (m *Model) currentSummary() (core.Summary, bool) {
	rows := m.visibleRows()
	if m.selRecipes < 0 || m.selRecipes >= len(rows) {
		return core.Summary{}, false
	}
	return rows[m.selRecipes], true
}

// readyEntry is one selectable row on the Ready tab.
//
// Adoptable recipes come first, then open small tasks. A task row carries the
// backlog item; a recipe row does not.
type readyEntry struct {
	id   string
	task bool
}

// readyEntries is the Ready tab's selectable rows, in the order it draws them.
func (m *Model) readyEntries() []readyEntry {
	out := make([]readyEntry, 0)
	for _, r := range core.Adoptable(m.recipes) {
		out = append(out, readyEntry{id: r.ID})
	}
	for _, it := range openSmallTasks(m.items) {
		out = append(out, readyEntry{id: it.ID, task: true})
	}
	return out
}

// openSmallTasks is the Ready tab's second block.
//
// The items come from BacklogGroups, and only the open ones stay: an idea is
// filed with the ideas, and a promoted or dismissed item has left.
func openSmallTasks(items []core.BacklogItem) []core.BacklogItem {
	_, tasks := core.BacklogGroups(items)
	out := make([]core.BacklogItem, 0)
	for _, it := range tasks {
		if it.Status == core.StatusOpen {
			out = append(out, it)
		}
	}
	return out
}

// currentReadyEntry is the row under the Ready cursor.
func (m *Model) currentReadyEntry() (readyEntry, bool) {
	all := m.readyEntries()
	if m.selReady < 0 || m.selReady >= len(all) {
		return readyEntry{}, false
	}
	return all[m.selReady], true
}

// selectedReadyTask is the small task under the cursor, when the cursor is on one.
func (m *Model) selectedReadyTask() (core.BacklogItem, bool) {
	e, ok := m.currentReadyEntry()
	if !ok || !e.task {
		return core.BacklogItem{}, false
	}
	for _, it := range m.items {
		if it.ID == e.id {
			return it, true
		}
	}
	return core.BacklogItem{}, false
}

// board is the tree the Board tab shows.
//
// The rows are the rendered lines; the recipes are the ids of the selectable ones, in
// the same order, so the selection always names a recipe the list draws.
type board struct {
	rows   []listLine
	recipe []string
}

// boardTree is the tree the Board tab shows, and the rows it draws.
//
// Below 160 columns it is the vertical tree. At 160 and above it is the same
// tree drawn left to right. The membership is TreeOf either way.
func (m *Model) boardTree() board {
	tree := core.TreeOf(m.recipes)
	if len(tree.Roots) == 0 && len(tree.Orphans) == 0 && len(tree.Stuck) == 0 {
		// "The tree is empty" means the plan has no recipes. A plan that has
		// recipes and an empty tree is a different thing, and the console must not
		// say the same words about both.
		text := "The tree is empty — no recipes."
		if len(m.recipes) > 0 {
			text = "The plan has recipes, but none of them is open."
		}
		b := &boardBuilder{m: m}
		b.out = append(b.out, listLine{
			text:    "  " + m.styDim.Render(text),
			itemIdx: -1,
		})
		return board{rows: b.out}
	}
	if m.width >= splitThreshold {
		return m.boardTreeWide(tree)
	}
	return m.boardTreeNarrow(tree)
}

// boardTreeNarrow draws the vertical tree.
//
// The connectors are ├─ and └─. Below 160 columns this is the whole drawing.
func (m *Model) boardTreeNarrow(tree core.Tree) board {
	b := &boardBuilder{m: m}
	for _, root := range tree.Roots {
		b.trunk(root)
	}
	for _, orphan := range tree.Orphans {
		b.orphan(orphan)
	}
	for _, stuck := range tree.Stuck {
		b.stuck(stuck)
	}
	return board{rows: b.out, recipe: b.recipe}
}

// currentRecipe is the recipe the detail pane shows.
func (m *Model) currentRecipe() (*core.Recipe, bool) {
	var id string
	switch m.tab {
	case TabRecipes:
		if s, ok := m.currentSummary(); ok {
			id = s.ID
		}
	case TabReady:
		if e, ok := m.currentReadyEntry(); ok && !e.task {
			id = e.id
		}
	case TabBoard:
		if r, ok := m.currentBoardRecipe(); ok {
			id = r
		}
	}
	if id == "" {
		return nil, false
	}
	for i := range m.recipes {
		if m.recipes[i].ID == id {
			return &m.recipes[i], true
		}
	}
	return nil, false
}

// currentBoardRecipe is the recipe the Board's selection points at.
func (m *Model) currentBoardRecipe() (string, bool) {
	ids := m.boardTree().recipe
	if m.selBoard < 0 || m.selBoard >= len(ids) {
		return "", false
	}
	return ids[m.selBoard], true
}

// currentSel is the selection index for the current tab.
func (m *Model) currentSel() int {
	switch m.tab {
	case TabRecipes:
		return m.selRecipes
	case TabBacklog:
		return m.selBacklog
	case TabReady:
		return m.selReady
	case TabBoard:
		return m.selBoard
	case TabLog:
		return m.selLog
	case TabSettings:
		return m.selTheme
	}
	return -1
}

// clampSelections keeps every cursor inside its list.
func (m *Model) clampSelections() {
	clamp := func(n *int, length int) {
		if length == 0 {
			*n = 0
			return
		}
		if *n < 0 {
			*n = 0
		}
		if *n >= length {
			*n = length - 1
		}
	}
	clamp(&m.selRecipes, len(m.visibleRows()))
	clamp(&m.selBacklog, m.backlogCount())
	clamp(&m.selReady, len(m.readyEntries()))
	clamp(&m.selBoard, len(m.boardTree().recipe))
	clamp(&m.selLog, len(m.logWindow()))
	clamp(&m.selTheme, len(palettes()))
	clamp(&m.selNote, len(m.openNotes()))
	m.rememberBacklogRow()
	m.rememberReadyRow()
}

// placeReady resolves the Ready cursor after a reload.
//
// leave does not name that id: a promotion took the item off the tab, and the
// cursor sits on the nearest remaining row. Otherwise the remembered id is
// resolved to its row.
func (m *Model) placeReady(leave string) {
	all := m.readyEntries()
	if leave != "" {
		m.selReady = nearestReadyRow(all, m.selReady, leave)
		m.readyID = readyIDAt(all, m.selReady)
		return
	}
	if i := readyRowOf(all, m.readyID); m.readyID != "" && i >= 0 {
		m.selReady = i
		return
	}
	if len(all) == 0 {
		m.selReady = 0
		m.readyID = ""
		return
	}
	if m.selReady < 0 {
		m.selReady = 0
	}
	if m.selReady >= len(all) {
		m.selReady = len(all) - 1
	}
	m.readyID = all[m.selReady].id
}

// rememberReadyRow stores the id of the row the index points at.
func (m *Model) rememberReadyRow() {
	m.readyID = readyIDAt(m.readyEntries(), m.selReady)
}

// readyRowOf is the index of id in the tab's order, or -1.
func readyRowOf(all []readyEntry, id string) int {
	for i, e := range all {
		if e.id == id {
			return i
		}
	}
	return -1
}

// readyIDAt is the id at index, or empty when the index names no row.
func readyIDAt(all []readyEntry, index int) string {
	if index < 0 || index >= len(all) {
		return ""
	}
	return all[index].id
}

// nearestReadyRow is the row nearest to index that is not left.
func nearestReadyRow(all []readyEntry, index int, left string) int {
	rest := make([]int, 0, len(all))
	for i, e := range all {
		if e.id != left {
			rest = append(rest, i)
		}
	}
	if len(rest) == 0 {
		return 0
	}
	pos := 0
	for _, i := range rest {
		if i < index {
			pos++
		}
	}
	if pos >= len(rest) {
		pos = len(rest) - 1
	}
	return rest[pos]
}

// handleKey dispatches a key press.
func (m *Model) handleKey(k tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	// An open prompt takes every key. The console must not move the selection out
	// from under a half-typed answer.
	if m.prompt != nil {
		return m.handlePromptKey(k)
	}

	switch k.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "r":
		m.refusal, m.status = "", ""
		return m, m.reload()

	case "esc":
		if m.mode == ModeDetail {
			m.mode = ModeSplit
			m.detailScroll = 0
		}
		return m, nil

	case "tab", "right", "l":
		m.switchTab(Tab((int(m.tab) + 1) % len(TabNames)))
		return m, nil

	case "shift+tab", "left", "h":
		m.switchTab(Tab((int(m.tab) - 1 + len(TabNames)) % len(TabNames)))
		return m, nil

	case "up", "k":
		if m.mode == ModeDetail {
			m.scrollDetail(-1)
			return m, nil
		}
		m.move(-1)
		return m, nil

	case "down", "j":
		if m.mode == ModeDetail {
			m.scrollDetail(1)
			return m, nil
		}
		m.move(1)
		return m, nil

	case "pgup":
		if m.mode == ModeDetail {
			m.scrollDetail(-m.listAreaHeight())
			return m, nil
		}
		m.move(-m.listAreaHeight())
		return m, nil

	case "pgdown":
		if m.mode == ModeDetail {
			m.scrollDetail(m.listAreaHeight())
			return m, nil
		}
		m.move(m.listAreaHeight())
		return m, nil

	case "home", "g":
		if m.mode == ModeDetail {
			m.detailScroll = 0
			return m, nil
		}
		m.moveTo(0)
		return m, nil

	case "end", "G":
		if m.mode == ModeDetail {
			m.scrollDetail(m.detailLineCount())
			return m, nil
		}
		m.moveTo(m.listCount() - 1)
		return m, nil

	case "enter":
		if m.tab == TabSettings {
			m.applyTheme()
			return m, nil
		}
		m.openDetail()
		return m, nil

	case "1", "2", "3":
		m.order = Order(int(k.String()[0]) - '1')
		m.sections = m.orderedSections()
		m.clampSelections()
		m.status = "order: " + OrderNames[int(m.order)]
		return m, nil

	case "c":
		m.showServed = !m.showServed
		m.clampSelections()
		return m, nil

	case "n":
		m.moveNote(1)
		return m, nil

	case "N":
		m.moveNote(-1)
		return m, nil

	case "d":
		// The backlog's own key: dismiss the selected item, with the reason.
		if m.tab == TabBacklog {
			return m.dismissSelected()
		}
		return m, nil

	case "a":
		// The backlog's own key: write a new item down.
		if m.tab == TabBacklog {
			m.openCreate()
			return m, nil
		}
		return m, nil

	case " ", "space":
		return m.primaryAction()
	}
	return m, nil
}

// switchTab goes to another tab.
//
// It is the **one** switch: the keyboard's `tab` and the mouse's click on a tab
// label both come through here. A click that set the tab and skipped the rest
// would be a second answer to one question.
func (m *Model) switchTab(t Tab) {
	m.tab = t
	m.status = ""
	// The tab changes the recipe, and the note cursor belongs to the recipe.
	// A footer that promises a triage the space refuses is a lie.
	m.clampSelections()
}

// openDetail gives the detail the screen.
//
// It is the keyboard's `enter`: the detail takes the screen, and `esc` brings it
// back. A tab that has nothing to open — the Settings, the Home — has no
// detail to open, so the move does nothing there.
func (m *Model) openDetail() {
	if m.mode != ModeSplit {
		return
	}
	if !m.hasDetail() {
		return
	}
	m.mode = ModeDetail
	m.detailScroll = 0
}

// hasDetail reports whether the current selection has a detail to open.
//
// The console has **one** detail pane, and it follows the tab: a recipe on Recipes,
// Ready and Board, a backlog item on Backlog, a journal entry on Log. The pane and
// the key read the same answer, so the key can never open something the pane does
// not show — and a backlog item is not a recipe, so it needs its own answer here.
func (m *Model) hasDetail() bool {
	switch m.tab {
	case TabRecipes, TabReady, TabBoard:
		_, ok := m.currentRecipe()
		return ok
	case TabBacklog:
		_, ok := m.selectedBacklogItem()
		return ok
	case TabLog:
		events := m.logWindow()
		return m.selLog >= 0 && m.selLog < len(events)
	}
	return false
}

// scrollDetail moves the detail pane's window.
//
// A recipe is longer than the pane — an intent, evidence per criterion, a review
// thread — so without this the bottom of a recipe is unreachable.
func (m *Model) scrollDetail(delta int) {
	m.detailScroll += delta
	if m.detailScroll < 0 {
		m.detailScroll = 0
	}
	max := m.detailLineCount() - 1
	if max < 0 {
		max = 0
	}
	if m.detailScroll > max {
		m.detailScroll = max
	}
}

// detailLineCount is how many rows the open recipe's detail has.
//
// It is counted at the width the pane draws with — the narrower one in the split
// layout — because a narrower pane wraps the same text into more rows. The bound
// and the pane have to agree, or the scroll stops early against a row count the
// pane never had.
func (m *Model) detailLineCount() int {
	return len(detailRows(m.detailLines(m.detailPaneWidth())))
}

// applyTheme makes the theme under the cursor the console's theme.
//
// The choice is for the session: nothing is written to the plan, so a theme is the
// chef's own setting, not the plan's.
func (m *Model) applyTheme() {
	list := palettes()
	if m.selTheme < 0 || m.selTheme >= len(list) {
		return
	}
	m.theme = newTheme(list[m.selTheme])
	m.status = "theme: " + list[m.selTheme].name
}

// move moves the cursor in the current tab.
func (m *Model) move(delta int) {
	m.status = ""
	switch m.tab {
	case TabRecipes:
		m.selRecipes += delta
	case TabBacklog:
		m.selBacklog += delta
	case TabReady:
		m.selReady += delta
	case TabBoard:
		m.selBoard += delta
	case TabLog:
		m.selLog += delta
	case TabSettings:
		m.selTheme += delta
	}
	m.clampSelections()
	m.detailScroll = 0
}

// moveTo puts the cursor on a specific item.
func (m *Model) moveTo(index int) {
	switch m.tab {
	case TabRecipes:
		m.selRecipes = index
	case TabBacklog:
		m.selBacklog = index
	case TabReady:
		m.selReady = index
	case TabBoard:
		m.selBoard = index
	case TabLog:
		m.selLog = index
	case TabSettings:
		m.selTheme = index
	}
	m.clampSelections()
	m.detailScroll = 0
}

// listCount is the number of selectable items in the current tab.
func (m *Model) listCount() int {
	switch m.tab {
	case TabRecipes:
		return len(m.visibleRows())
	case TabBacklog:
		return m.backlogCount()
	case TabReady:
		return len(m.readyEntries())
	case TabBoard:
		return len(m.boardTree().recipe)
	case TabLog:
		return len(m.logWindow())
	case TabSettings:
		return len(palettes())
	}
	return 0
}

// primaryAction does the one thing the chef most likely wants on the selected item.
func (m *Model) primaryAction() (tea.Model, tea.Cmd) {
	switch m.tab {
	case TabBacklog:
		return m.backlogPrimary()
	case TabReady:
		// An adoptable recipe has no move here. Claiming it is the agent's.
		if it, ok := m.selectedReadyTask(); ok {
			return m.promoteItem(it)
		}
		return m, nil
	case TabRecipes:
		return m.actOnRecipe()
	case TabBoard:
		// The footer names no space here. The sign-off stays on Recipes.
		return m, nil
	}
	return m, nil
}

// actOnRecipe performs the recipe's primary action for its state.
func (m *Model) actOnRecipe() (tea.Model, tea.Cmd) {
	rec, ok := m.currentRecipe()
	if !ok {
		return m, nil
	}
	// Only the chef's own moves are offered. Claiming the work is an agent's move,
	// and the console is not an agent — it must not record itself as one.
	switch core.ChefMoveOf(*rec) {
	case core.MoveSignOff:
		r, err := m.ops.ApproveRecipe(rec.ID, "chef")
		if err != nil {
			m.refusal = refusalText(err)
			m.status = ""
			return m, nil
		}
		m.refusal = ""
		m.status = fmt.Sprintf("signed off %s — an agent can pick it up now", r.ID)
	case core.MoveReview:
		m.openReview(rec.ID)
		return m, nil
	case core.MoveTriage:
		note, ok := m.selectedNote()
		if !ok {
			m.status = "no open tasting note here"
			return m, nil
		}
		m.openTriage(rec.ID, note.ID)
		return m, nil
	case core.MoveServe:
		r, err := m.ops.Advance(rec.ID, core.EventClose)
		if err != nil {
			m.refusal = refusalText(err)
			m.status = ""
			return m, nil
		}
		m.refusal = ""
		m.status = fmt.Sprintf("served %s", r.ID)
	default:
		m.status = "nothing for you here — the agents have it"
	}
	return m, m.reload()
}

// backlogRows is the Backlog tab's rows, in the order the tab draws them.
//
// Ideas first, then small tasks. The cursor, the row and the reload all read
// this one slice, so the highlighted row is the item the cursor names.
func (m *Model) backlogRows() []core.BacklogItem {
	ideas, tasks := core.BacklogGroups(m.items)
	return append(append([]core.BacklogItem{}, ideas...), tasks...)
}

// backlogCount is the number of rows in the Backlog tab.
func (m *Model) backlogCount() int {
	return len(m.backlogRows())
}

// selectedBacklogItem is the backlog item under the cursor.
//
// The second value says whether there is one: the cursor can sit past the last row,
// and then there is nothing under it.
func (m *Model) selectedBacklogItem() (core.BacklogItem, bool) {
	all := m.backlogRows()
	if m.selBacklog < 0 || m.selBacklog >= len(all) {
		return core.BacklogItem{}, false
	}
	return all[m.selBacklog], true
}

// placeBacklog resolves the backlog cursor after a reload.
//
// follow lands on that id — a create, or a shape that kept its id and changed
// its row. leave does not name that id: the item has left the chef's hands, and
// the cursor sits on the nearest remaining row. Otherwise the remembered id is
// resolved to its row. An id that is gone is not kept; the index names the
// nearest row, and that row's id is the one remembered.
func (m *Model) placeBacklog(follow, leave string) {
	all := m.backlogRows()
	if follow != "" {
		if i := backlogRowOf(all, follow); i >= 0 {
			m.selBacklog = i
			m.backlogID = follow
			return
		}
	}
	if leave != "" {
		m.selBacklog = nearestBacklogRow(all, m.selBacklog, leave)
		m.backlogID = backlogIDAt(all, m.selBacklog)
		return
	}
	if i := backlogRowOf(all, m.backlogID); m.backlogID != "" && i >= 0 {
		m.selBacklog = i
		return
	}
	if len(all) == 0 {
		m.selBacklog = 0
		m.backlogID = ""
		return
	}
	if m.selBacklog < 0 {
		m.selBacklog = 0
	}
	if m.selBacklog >= len(all) {
		m.selBacklog = len(all) - 1
	}
	m.backlogID = all[m.selBacklog].ID
}

// rememberBacklogRow stores the id of the row the index points at.
func (m *Model) rememberBacklogRow() {
	m.backlogID = backlogIDAt(m.backlogRows(), m.selBacklog)
}

// backlogRowOf is the index of id in the tab's order, or -1.
func backlogRowOf(all []core.BacklogItem, id string) int {
	for i, it := range all {
		if it.ID == id {
			return i
		}
	}
	return -1
}

// backlogIDAt is the id at index, or empty when the index names no row.
func backlogIDAt(all []core.BacklogItem, index int) string {
	if index < 0 || index >= len(all) {
		return ""
	}
	return all[index].ID
}

// nearestBacklogRow is the row nearest to index that is not left.
//
// The next row is the one that would sit at this index if left had gone. The
// previous row is that one when this was the last. The only row is itself:
// there is nowhere else to sit.
func nearestBacklogRow(all []core.BacklogItem, index int, left string) int {
	rest := make([]int, 0, len(all))
	for i, it := range all {
		if it.ID != left {
			rest = append(rest, i)
		}
	}
	if len(rest) == 0 {
		return 0
	}
	pos := 0
	for _, i := range rest {
		if i < index {
			pos++
		}
	}
	if pos >= len(rest) {
		pos = len(rest) - 1
	}
	return rest[pos]
}

// promoteSelected promotes the selected backlog item into a recipe.
//
// The console **asks** the domain before it acts: `core.PromotionOf` is a read, and
// its answer is the reason the banner shows when the answer is no. The console does
// not call the operation to find out — the operation **journals** every refusal, and
// the journal is the plan's only event log, append-only. Calling it would record an
// attempt the console never offered: the footer, in the same frame, offers no move
// at all on an item the domain refuses.
//
// The recipe path's precedent is the same: the console asks `core.ChefMoveOf` and
// does not call the operation when there is no move. A move that **is** offered and
// is turned away by the plan is an attempt, and the plan's own refusal says so.
func (m *Model) promoteSelected() (tea.Model, tea.Cmd) {
	item, ok := m.selectedBacklogItem()
	if !ok {
		m.status = "nothing to promote here"
		return m, nil
	}
	return m.promoteItem(item)
}

// promoteItem promotes a small task into a recipe.
//
// The console asks the domain before it acts. The cursor does not stay on the
// id that left, on the backlog or on Ready, and the tab does not change.
func (m *Model) promoteItem(item core.BacklogItem) (tea.Model, tea.Cmd) {
	if p := core.PromotionOf(item); !p.Allowed {
		m.refusal = refusalText(&core.RefusalError{Reason: p.Reason})
		m.status = ""
		return m, nil
	}
	recipe, err := m.ops.PromoteBacklogItem(item.ID)
	if err != nil {
		m.refusal = refusalText(err)
		return m, nil
	}
	m.refusal = ""
	m.status = fmt.Sprintf("promoted %s to recipe %s", item.ID, recipe.ID)
	m.backlogLeave = item.ID
	m.readyLeave = item.ID
	return m, m.reload()
}

// refusalText turns a refused action into a banner.
//
// A refusal is information, not an error: the gate working as designed.
func refusalText(err error) string {
	msg := err.Error()
	var ge *core.GateError
	if errors.As(err, &ge) {
		var parts []string
		for _, g := range ge.Gates {
			if i := strings.Index(g, " ("); i > 0 {
				parts = append(parts, g[:i])
				continue
			}
			parts = append(parts, g)
		}
		return fmt.Sprintf("refused — %s", strings.Join(parts, ", "))
	}
	return "refused — " + msg
}
