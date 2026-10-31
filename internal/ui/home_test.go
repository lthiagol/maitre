package ui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/lthiagol/maitre/internal/brand"
	"github.com/lthiagol/maitre/internal/core"
)

// The Home tab — the console's landing place.
//
// **The invariant: the Home tab shows numbers; it shows no items.** No row on Home
// is selectable, `space` does nothing there, and the tab's item count is zero.
//
// The precedent is v1's **Overview** lane, dropped because it restated the board.
// The board shows the recipes; Home shows their counts.

// TestTheConsoleLandsOnHome covers the tab's place.
//
// The console opens on Home: it is where the chef lands, and the numbers are the
// first thing on the screen.
func TestTheConsoleLandsOnHome(t *testing.T) {
	m := newModel(t)

	if m.tab != TabHome {
		t.Fatalf("the console opens on %q, want Home", m.tab.Label())
	}
	if len(TabNames) != 7 {
		t.Fatalf("the console has %d tabs, want 7: %v", len(TabNames), TabNames)
	}
	want := []string{"Home", "Recipes", "Backlog", "Ready", "Board", "Log", "Settings"}
	for i, name := range want {
		if TabNames[i] != name {
			t.Errorf("tab %d is %q, want %q", i, TabNames[i], name)
		}
	}
	out := plain(m)
	if !strings.Contains(out, "▸ Home") {
		t.Errorf("the tab bar should open on Home:\n%s", out)
	}
}

// TestTheHomeTabShowsTheNumbers covers every row of the tab.
//
// Every number is the domain's: the test asks `core.PlanStats` for the plan's
// numbers and looks for them on the screen. A number the console counted itself
// would be a number the CLI counts differently.
func TestTheHomeTabShowsTheNumbers(t *testing.T) {
	m := newModel(t)
	stats := core.PlanStatsOf(m.recipes, m.items, m.journal)
	out := plain(m)

	for _, want := range []string{
		"name " + m.projectName(),
		"path " + truncate(humanPath(m.store.Root), planPathColumns),
		"total " + strconv.Itoa(stats.Recipes.Total),
		"need you " + strconv.Itoa(stats.Recipes.NeedsYou),
		"on the line " + strconv.Itoa(stats.Recipes.OnTheLine),
		"being shaped " + strconv.Itoa(stats.Recipes.BeingShaped),
		"on hold " + strconv.Itoa(stats.Recipes.OnHold),
		"served " + strconv.Itoa(stats.Recipes.Served),
		"criteria " + strconv.Itoa(stats.Evidence.Criteria),
		"proven " + strconv.Itoa(stats.Evidence.Proven),
		"failed " + strconv.Itoa(stats.Evidence.Failed),
		"rejected " + strconv.Itoa(stats.Evidence.Rejected),
		"unproven " + strconv.Itoa(stats.Evidence.Unproven),
		"pending " + strconv.Itoa(stats.Evidence.Pending),
		"open notes " + strconv.Itoa(stats.Tasting.OpenNotes),
		"open findings " + strconv.Itoa(stats.Tasting.OpenFindings),
		"ideas " + strconv.Itoa(stats.Backlog.Ideas),
		"small tasks " + strconv.Itoa(stats.Backlog.SmallTasks),
		"plan written under " + strconv.Itoa(stats.Automation.PlanVersion),
		"this build speaks " + strconv.Itoa(stats.Automation.BuildVersion),
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the tab should show %q:\n%s", want, out)
		}
	}
}

// TestTheHomeTabShowsTheAutomation covers the vital signs.
//
// The Sous-Chef's content became Home's automation block: the vital signs, and the
// honesty the tab carried. The automation is not built, and the block says so —
// the signs alone would read as a running automation.
func TestTheHomeTabShowsTheAutomation(t *testing.T) {
	m := newModel(t)
	out := plain(m)

	for _, want := range []string{"THE AUTOMATION", "the automation is not built", "plan written under", "this build speaks", "they agree", "last event"} {
		if !strings.Contains(out, want) {
			t.Errorf("the tab should show %q:\n%s", want, out)
		}
	}
}

// TestTheHomeTabSaysWhenTheVersionsDisagree covers the mismatch.
//
// A plan written by a different build of the tool is a fact the chef needs, and the
// version guard turns it into a refusal to write.
func TestTheHomeTabSaysWhenTheVersionsDisagree(t *testing.T) {
	m := newModel(t)
	stampPlan(t, m.store, core.DomainVersion+1)

	m = loadModel(t, m.store)
	out := plain(m)

	if !strings.Contains(out, "they do not agree") {
		t.Errorf("the tab should say the versions disagree:\n%s", out)
	}
	if !strings.Contains(out, "does not write plans written under another version") {
		t.Errorf("the tab should say what a mismatch means:\n%s", out)
	}
}

// TestTheHomeTabSaysWhenThereAreNoEvents covers the empty plan.
//
// A plan with no journal has no last event. It does not have a blank one.
func TestTheHomeTabSaysWhenThereAreNoEvents(t *testing.T) {
	root := filepath.Join(t.TempDir(), core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	out := plain(m)

	if !strings.Contains(out, "the plan has no events") {
		t.Errorf("an empty plan should say it has no events:\n%s", out)
	}
}

// TestTheHomeTabShowsNoItems holds the invariant.
//
// **The Home tab shows numbers; it shows no items.** No row on Home is selectable,
// `space` does nothing there, and the tab's item count is zero.
//
// This is the test that stops Home from becoming a second board: make a row on
// Home selectable and this goes red.
func TestTheHomeTabShowsNoItems(t *testing.T) {
	m := newModel(t)

	for i, line := range m.listLines() {
		if line.itemIdx >= 0 {
			t.Errorf("row %d holds the item %d — Home shows numbers, not items", i, line.itemIdx)
		}
	}
	if got := m.listCount(); got != 0 {
		t.Errorf("the tab's item count is %d, want 0", got)
	}
	if got := m.currentSel(); got != -1 {
		t.Errorf("the cursor sits on %d — there is nothing on Home to sit on", got)
	}
}

// TestSpaceDoesNothingOnHome covers the promise.
//
// `space` is the chef's move, and the Home tab has none. A move that is offered
// and refused is an attempt, and the plan's journal would record it.
func TestSpaceDoesNothingOnHome(t *testing.T) {
	m := newModel(t)
	before := journalBytes(t, m)

	m.Update(tea.KeyPressMsg{Code: ' '})

	if after := journalBytes(t, m); after != before {
		t.Errorf("space on Home wrote to the plan:\n%s", after)
	}
	if m.status != "" {
		t.Errorf("space on Home said %q — it promises nothing", m.status)
	}
	if m.refusal != "" {
		t.Errorf("space on Home refused something: %q", m.refusal)
	}
}

// TestTheHomeFooterPromisesNoMove covers the footer.
//
// The chef's footers are promises: a footer that promises a move the console
// refuses is a lie. Home has no item under the cursor, so it promises no move —
// not `space`, not `enter`, not the cursor.
func TestTheHomeFooterPromisesNoMove(t *testing.T) {
	m := newModel(t)
	out := plain(m)

	for _, want := range []string{"space", "enter detail", "move"} {
		if strings.Contains(out, want) {
			t.Errorf("the footer should not promise %q:\n%s", want, out)
		}
	}
	for _, want := range []string{"tab", "reload", "quit"} {
		if !strings.Contains(out, want) {
			t.Errorf("the footer should say %q:\n%s", want, out)
		}
	}
}

// TestTheHomeKeepsItsNumbersAt80x24 covers the content's claim on the screen.
//
// The decoration yields to the content, in both dimensions. At 80×24 — the classic
// default — the numbers are whole and the mark is absent: the tab is what it would
// be without the art.
func TestTheHomeKeepsItsNumbersAt80x24(t *testing.T) {
	m := newModel(t)
	m.width, m.height = 80, 24
	out := plain(m)

	stats := core.PlanStatsOf(m.recipes, m.items, m.journal)
	for _, want := range []string{
		"total " + strconv.Itoa(stats.Recipes.Total),
		"need you " + strconv.Itoa(stats.Recipes.NeedsYou),
		"being shaped " + strconv.Itoa(stats.Recipes.BeingShaped),
		"on hold " + strconv.Itoa(stats.Recipes.OnHold),
		"criteria " + strconv.Itoa(stats.Evidence.Criteria),
		"proven " + strconv.Itoa(stats.Evidence.Proven),
		"pending " + strconv.Itoa(stats.Evidence.Pending),
		"open notes " + strconv.Itoa(stats.Tasting.OpenNotes),
		"ideas " + strconv.Itoa(stats.Backlog.Ideas),
		"plan written under " + strconv.Itoa(stats.Automation.PlanVersion),
	} {
		if !strings.Contains(out, want) {
			t.Errorf("at 80x24 the numbers should be whole — %q is not on the screen:\n%s", want, out)
		}
	}
	if hasMark(out) {
		t.Errorf("at 80x24 the mark should be absent:\n%s", out)
	}
}

// TestTheHomeRowsAddUpOnOneLine holds the two rows whose parts are the states.
//
// A row that names some of the states and not the rest misleads: the parts beside
// the number do not make it. The sections are a partition, and each row stays one
// line — a wrap would push the footer off the bottom.
func TestTheHomeRowsAddUpOnOneLine(t *testing.T) {
	m := newModel(t)
	m.width, m.height = 80, 24
	stats := core.PlanStatsOf(m.recipes, m.items, m.journal)

	sections := stats.Recipes.NeedsYou + stats.Recipes.OnTheLine +
		stats.Recipes.BeingShaped + stats.Recipes.OnHold + stats.Recipes.Served
	if sections != stats.Recipes.Total {
		t.Errorf("the sections sum to %d and the plan holds %d recipes: %+v", sections, stats.Recipes.Total, stats.Recipes)
	}
	states := stats.Evidence.Proven + stats.Evidence.Failed +
		stats.Evidence.Rejected + stats.Evidence.Unproven + stats.Evidence.Pending
	if states != stats.Evidence.Criteria {
		t.Errorf("the states sum to %d and the plan holds %d criteria: %+v", states, stats.Evidence.Criteria, stats.Evidence)
	}

	for _, line := range strings.Split(StripANSI(m.View().Content), "\n") {
		recipes := strings.Contains(line, "total ") && strings.Contains(line, "being shaped ")
		evidence := strings.Contains(line, "criteria ") && strings.Contains(line, "pending ")
		if !recipes && !evidence {
			continue
		}
		if w := lipgloss.Width(line); w > m.width {
			t.Errorf("the row is %d columns and the terminal is %d: %q", w, m.width, line)
		}
	}
}

// TestTheHomeDrawsTheMarkBesideTheNumbers covers the art.
//
// The mark moved here so it could **share** the space: it sits beside the numbers,
// on their rows, and the numbers are not pushed down by it.
func TestTheHomeDrawsTheMarkBesideTheNumbers(t *testing.T) {
	m := newModel(t)
	m.width, m.height = 120, 40

	content := m.homeContent()
	lines := m.listLines()

	marked := -1
	for i, line := range lines {
		if strings.Contains(StripANSI(line.text), brand.Full[0]) {
			marked = i
			break
		}
	}
	if marked < 0 {
		t.Fatalf("at 120x40 the mark should be drawn beside the numbers:\n%s", plain(m))
	}
	// **Beside, not above**: the mark shares the content's rows, and the content
	// starts where it always starts.
	if marked >= len(content) {
		t.Errorf("the mark is on row %d and the content has %d rows — it is below the numbers, not beside them", marked, len(content))
	}
	if !strings.Contains(StripANSI(lines[0].text), "THE PLAN") {
		t.Errorf("the numbers should start at the top — the mark pushed them down:\n%s", plain(m))
	}
}

// TestTheHomeArtNeedsBothDimensions covers the two rules.
//
// The width: the numbers keep their column, and the mark is drawn only in what is
// left. The height: the numbers must be whole, and a tab whose numbers are cut off
// has no room for a face.
func TestTheHomeArtNeedsBothDimensions(t *testing.T) {
	t.Run("the width", func(t *testing.T) {
		m := newModel(t)
		m.width, m.height = 80, 40
		if hasMark(plain(m)) {
			t.Errorf("at 80 columns the mark should be absent — the numbers need the width:\n%s", plain(m))
		}
		m.width, m.height = 120, 40
		if !hasMark(plain(m)) {
			t.Errorf("at 120 columns the mark should be drawn:\n%s", plain(m))
		}
	})

	t.Run("the height", func(t *testing.T) {
		m := newModel(t)
		m.width, m.height = 120, 12
		if hasMark(plain(m)) {
			t.Errorf("at 12 rows the mark should be absent — the numbers need the height:\n%s", plain(m))
		}
	})
}

// TestALongPlanNameStaysOnOneRow covers a name the row cannot hold.
//
// The review measured it: a name of 100 characters at 80 columns wrapped the plan
// row, and the wrap pushed the footer off the bottom. The name is `n` repeated 100
// times. The row stays one line, every line fits the width, and the footer is the
// last row.
func TestALongPlanNameStaysOnOneRow(t *testing.T) {
	const width = 80
	name := strings.Repeat("n", 100)

	root := filepath.Join(t.TempDir(), core.PlanDirName)
	store := core.Open(root)
	if err := core.Seed(store); err != nil {
		t.Fatal(err)
	}
	m := loadInto(t, New(store, WithProject(name)))
	m.width, m.height = width, 24

	out := StripANSI(m.View().Content)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 24 {
		t.Fatalf("the screen has %d rows, want 24 — a wrapped row took one the layout did not count", len(lines))
	}
	if !strings.Contains(lines[23], "quit") {
		t.Errorf("the footer should be the last row, got %q", lines[23])
	}

	plan := ""
	for i, line := range lines {
		if w := lipgloss.Width(line); w > width {
			t.Errorf("row %d is %d columns wide and the terminal is %d: %q", i, w, width, line)
		}
		if strings.Contains(line, "name ") && strings.Contains(line, "path ") {
			plan = line
		}
	}
	if plan == "" {
		t.Fatal("the plan row is not on the screen")
	}
	if strings.Contains(plan, name) {
		t.Errorf("the 100-character name is whole on the row: %q", plan)
	}
	if !strings.Contains(plan, "…") {
		t.Errorf("the name should be truncated to its tail: %q", plan)
	}
}

// hasMark reports whether the screen carries the mark.
func hasMark(out string) bool {
	return strings.Contains(out, brand.Full[0]) || strings.Contains(out, brand.Half[0])
}

// stampPlan writes a domain version into the plan's journal, as if a tool of that
// version had written it.
func stampPlan(t *testing.T, store *core.Store, version int) {
	t.Helper()
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) == 0 {
		t.Fatal("the journal is empty")
	}
	j.Events[len(j.Events)-1].DomainVersion = version
	raw, err := json.Marshal(j)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(store.Root, core.FileActivity), raw, 0o644); err != nil {
		t.Fatal(err)
	}
}

// journalBytes is the plan's journal, byte for byte.
//
// The journal is the plan's only event log, so "nothing happened" is "the journal
// is the same file it was".
func journalBytes(t *testing.T, m *Model) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(m.store.Root, core.FileActivity))
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		t.Fatal(err)
	}
	return string(raw)
}
