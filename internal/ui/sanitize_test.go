package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/kitchen-brigade/internal/core"
)

// stressStore seeds the pathological plan.
func stressStore(t *testing.T) *core.Store {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := core.SeedStress(store); err != nil {
		t.Fatalf("seed stress: %v", err)
	}
	return store
}

// TestNothingRenderedCarriesControlCharacters is the invariant that matters.
//
// Every title, note and command in the console was written by an agent. If one carries a
// newline, a tab or an escape sequence, the screen breaks. So: nothing rendered may
// contain a control character.
func TestNothingRenderedCarriesControlCharacters(t *testing.T) {
	m, err := Snapshot(stressStore(t), 100)
	if err != nil {
		t.Fatal(err)
	}
	m.height = 40

	tabs := []Tab{TabRecipes, TabBacklog, TabPath, TabLog, TabSettings}
	for _, tab := range tabs {
		m.tab = tab
		checkNoControlCharacters(t, "tab "+tab.Label(), m.View().Content)
	}

	// The detail panes too — for every recipe, not just the selected one.
	for _, r := range m.recipes {
		rec := r
		m.tab = TabRecipes
		m.selRecipes = indexOf(m, rec.ID)
		m.mode = ModeDetail
		checkNoControlCharacters(t, "detail "+rec.ID, m.View().Content)
	}
}

// checkNoControlCharacters fails on any control character inside a rendered line.
//
// Our own styling is full of escape sequences, so the ANSI is stripped first — what is
// left must be printable.
func checkNoControlCharacters(t *testing.T, what, content string) {
	t.Helper()
	for i, line := range strings.Split(stripANSI(content), "\n") {
		for _, r := range line {
			if r < 0x20 || r == 0x7f {
				t.Errorf("%s: line %d has the control character %q: %q", what, i+1, r, line)
				return
			}
		}
	}
}

// indexOf finds a recipe's position in the visible list.
func indexOf(m *Model, id string) int {
	for i, s := range m.visibleRows() {
		if s.ID == id {
			return i
		}
	}
	return 0
}

// TestTheConsoleSurvivesThePathologicalPlan renders every tab of the stress plan and
// checks the screen is still a screen.
func TestTheConsoleSurvivesThePathologicalPlan(t *testing.T) {
	m, err := Snapshot(stressStore(t), 100)
	if err != nil {
		t.Fatal(err)
	}
	m.height = 30

	for _, tab := range []Tab{TabRecipes, TabBacklog, TabPath, TabLog, TabSettings} {
		m.tab = tab
		out := StripANSI(m.View().Content)
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if len(lines) != 30 {
			t.Errorf("the %s tab rendered %d rows, want 30", tab.Label(), len(lines))
		}
	}
}

// TestTheStressPlanMatchesTheGolden pins the layout against the pathological plan.
func TestTheStressPlanMatchesTheGolden(t *testing.T) {
	m, err := Snapshot(stressStore(t), 100)
	if err != nil {
		t.Fatal(err)
	}
	m.height = 40
	got := StripANSI(m.View().Content)

	if *update {
		if err := os.MkdirAll(filepath.Dir(stressGoldenPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(stressGoldenPath, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Log("stress golden updated")
		return
	}

	want, err := os.ReadFile(stressGoldenPath)
	if err != nil {
		t.Fatalf("read stress golden: %v — run `make regen-goldens`", err)
	}
	if got != string(want) {
		t.Errorf("the console does not match the stress golden.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestSanitizeKeepsTheTextAndDropsTheDangerous covers the sanitiser directly.
func TestSanitizeKeepsTheTextAndDropsTheDangerous(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"plain title", "plain title"},
		{"one\ntwo", "one two"},
		{"one\ttwo", "one two"},
		{"one\r\ntwo", "one  two"},
		{"\x1b[31mred\x1b[0m", "red"},
		{"a\u202eb", "ab"},
		{"日本語", "日本語"},
	}
	for _, c := range cases {
		if got := sanitize(c.in); got != c.want {
			t.Errorf("sanitize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
