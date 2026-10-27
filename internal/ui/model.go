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
	TabRecipes Tab = iota
	TabBacklog
	TabPath
	TabChopboard
	TabLog
	TabExpo
	TabSettings
)

// TabNames are the tab labels, in order.
//
// The Chopboard sits next to the Path: the two are the same plan seen two ways — the
// flat order, and the tree. The Expo sits with the Log: the two are the automation's
// running, seen two ways — the history, and the place that will watch it happen.
var TabNames = []string{"Recipes", "Backlog", "Path", "Chopboard", "Log", "Expo", "Settings"}

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
	recipes  []core.Recipe
	specials []core.Special
	notes    []core.PrepNote
	journal  core.Journal
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

	specials []core.Special
	notes    []core.PrepNote

	journal core.Journal

	selRecipes int
	selBacklog int
	selPath    int
	selChop    int
	selLog     int
	// selTheme is the cursor over the theme list on the Settings tab.
	selTheme int
	// selNote is the cursor over the selected recipe's open tasting notes.
	selNote int

	// offset is the first visible line of the list.
	offset int

	// detailScroll is how far the detail pane has been scrolled.
	detailScroll int

	width  int
	height int
	ready  bool
	err    error

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
func New(store *core.Store) *Model {
	return newConsole(store, "console")
}

// NewLive builds the console model over a store, and it watches the disk.
//
// The chef sees a change on disk without pressing `r`: the console stats the plan's
// files once a second, and reloads when one of them has changed.
func NewLive(store *core.Store) *Model {
	m := newConsole(store, "console")
	m.live = true
	return m
}

// NewScripted builds a console whose writes are stamped as scripted.
//
// The console is the chef's surface: a human pressing space in it is the chef
// deciding. A script pressing space for them is the chef's decision mediated by a
// machine — which the journal should be able to show.
func NewScripted(store *core.Store) *Model {
	return newConsole(store, "console-script")
}

func newConsole(store *core.Store, session string) *Model {
	return &Model{
		store:      store,
		ops:        core.NewOperations(store, core.AuditContext{Actor: "chef", Session: session, Role: "chef"}),
		showServed: true,
		// Kitchen is the default: a console nobody told otherwise looks exactly
		// as it did before the theme was the chef's to choose.
		theme:        newTheme(kitchen()),
		pollInterval: PollInterval,
	}
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
		specials, err := m.store.ListSpecials()
		if err != nil {
			return reloadMsg{err: err}
		}
		notes, err := m.store.ListPrepNotes()
		if err != nil {
			return reloadMsg{err: err}
		}
		journal, err := m.store.LoadJournal()
		if err != nil {
			return reloadMsg{err: err}
		}
		// Everything here was written by an agent. Sanitise before it reaches the
		// screen, so a newline in a title cannot corrupt the layout.
		sanitizePlan(recipes, specials, notes, journal)
		return reloadMsg{recipes: recipes, specials: specials, notes: notes, journal: journal, disk: disk}
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
		m.recipes = msg.recipes
		m.specials = msg.specials
		m.notes = msg.notes
		m.journal = msg.journal
		m.err = nil
		m.sections = m.orderedSections()
		m.clampSelections()
		m.disk = msg.disk
		return m, m.poll()

	case pollMsg:
		return m.checkDisk()

	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}
	return m, nil
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
			sec.Items = append(sec.Items, core.SummaryOf(items[i]))
		}
		out = append(out, sec)
	}
	return out
}

// sectionsPlain is a single section, by id.
func (m *Model) sectionsPlain() []core.TriageSection {
	sec := core.TriageSection{Name: "All recipes"}
	for i := range m.recipes {
		sec.Items = append(sec.Items, core.SummaryOf(m.recipes[i]))
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

// currentPathEntry is the selected entry on the Path tab.
func (m *Model) currentPathEntry() (core.PathEntry, bool) {
	ready, waiting := core.Path(m.recipes)
	all := append(append([]core.PathEntry{}, ready...), waiting...)
	if m.selPath < 0 || m.selPath >= len(all) {
		return core.PathEntry{}, false
	}
	return all[m.selPath], true
}

// chopboard is the tree the Chopboard tab shows.
//
// The rows are the rendered lines; the recipes are the ids of the selectable ones, in
// the same order, so the selection always names a recipe the list draws.
type chopboard struct {
	rows   []listLine
	recipe []string
}

// chopTree is the tree the Chopboard tab shows, and the rows it draws.
func (m *Model) chopTree() chopboard {
	tree := core.TreeOf(m.recipes)
	b := &chopBuilder{m: m}
	if len(tree.Roots) == 0 && len(tree.Orphans) == 0 && len(tree.Stuck) == 0 {
		// "The tree is empty" means the plan has no recipes. A plan that has
		// recipes and an empty tree is a different thing, and the console must not
		// say the same words about both.
		text := "The tree is empty — no recipes."
		if len(m.recipes) > 0 {
			text = "The plan has recipes, but none of them is open."
		}
		b.out = append(b.out, listLine{
			text:    "  " + m.styDim.Render(text),
			itemIdx: -1,
		})
		return chopboard{rows: b.out}
	}
	for _, root := range tree.Roots {
		b.trunk(root)
	}
	for _, orphan := range tree.Orphans {
		b.orphan(orphan)
	}
	for _, stuck := range tree.Stuck {
		b.stuck(stuck)
	}
	return chopboard{rows: b.out, recipe: b.recipe}
}

// currentRecipe is the recipe the detail pane shows.
func (m *Model) currentRecipe() (*core.Recipe, bool) {
	var id string
	switch m.tab {
	case TabRecipes:
		if s, ok := m.currentSummary(); ok {
			id = s.ID
		}
	case TabPath:
		if e, ok := m.currentPathEntry(); ok {
			id = e.Summary.ID
		}
	case TabChopboard:
		if r, ok := m.currentChopRecipe(); ok {
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

// currentChopRecipe is the recipe the Chopboard's selection points at.
func (m *Model) currentChopRecipe() (string, bool) {
	ids := m.chopTree().recipe
	if m.selChop < 0 || m.selChop >= len(ids) {
		return "", false
	}
	return ids[m.selChop], true
}

// currentSel is the selection index for the current tab.
func (m *Model) currentSel() int {
	switch m.tab {
	case TabRecipes:
		return m.selRecipes
	case TabBacklog:
		return m.selBacklog
	case TabPath:
		return m.selPath
	case TabChopboard:
		return m.selChop
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
	clamp(&m.selBacklog, len(m.specials)+len(m.notes))
	clamp(&m.selPath, len(m.pathEntries()))
	clamp(&m.selChop, len(m.chopTree().recipe))
	clamp(&m.selLog, len(m.logWindow()))
	clamp(&m.selTheme, len(palettes()))
	clamp(&m.selNote, len(m.openNotes()))
}

// pathEntries is the Path tab's selectable items.
func (m *Model) pathEntries() []core.PathEntry {
	ready, waiting := core.Path(m.recipes)
	return append(append([]core.PathEntry{}, ready...), waiting...)
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
		m.tab = Tab((int(m.tab) + 1) % len(TabNames))
		m.status = ""
		// The tab changes the recipe, and the note cursor belongs to the recipe.
		// A footer that promises a triage the space refuses is a lie.
		m.clampSelections()
		return m, nil

	case "shift+tab", "left", "h":
		m.tab = Tab((int(m.tab) - 1 + len(TabNames)) % len(TabNames))
		m.status = ""
		m.clampSelections()
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
		if m.mode == ModeSplit {
			if _, ok := m.currentRecipe(); ok {
				m.mode = ModeDetail
				m.detailScroll = 0
			}
		}
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

	case " ", "space":
		return m.primaryAction()
	}
	return m, nil
}

// scrollDetail moves the detail pane's window.
//
// A recipe is longer than the pane — an intent, a proof per criterion, a review
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
func (m *Model) detailLineCount() int {
	return len(detailRows(m.detailLines(m.width - 4)))
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
	case TabPath:
		m.selPath += delta
	case TabChopboard:
		m.selChop += delta
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
	case TabPath:
		m.selPath = index
	case TabChopboard:
		m.selChop = index
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
		return len(m.specials) + len(m.notes)
	case TabPath:
		return len(m.pathEntries())
	case TabChopboard:
		return len(m.chopTree().recipe)
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
		return m.promoteSelected()
	case TabRecipes, TabPath, TabChopboard:
		return m.actOnRecipe()
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

// promoteSelected promotes the selected special or prep note up the ladder.
func (m *Model) promoteSelected() (tea.Model, tea.Cmd) {
	if m.selBacklog < len(m.specials) {
		sp := m.specials[m.selBacklog]
		note, err := m.ops.PromoteSpecial(sp.ID)
		if err != nil {
			m.refusal = refusalText(err)
			return m, nil
		}
		m.refusal = ""
		m.status = fmt.Sprintf("promoted %s to prep note %s", sp.ID, note.ID)
		return m, m.reload()
	}
	idx := m.selBacklog - len(m.specials)
	if idx < len(m.notes) {
		note := m.notes[idx]
		recipe, err := m.ops.PromotePrepNote(note.ID)
		if err != nil {
			m.refusal = refusalText(err)
			return m, nil
		}
		m.refusal = ""
		m.status = fmt.Sprintf("promoted %s to recipe %s", note.ID, recipe.ID)
		return m, m.reload()
	}
	return m, nil
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
