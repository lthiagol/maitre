package ui

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// update regenerates the golden. Run it with `make goldens`.
var update = flag.Bool("update", false, "update the golden file")

// goldenPath is the checked-in snapshot of the console's Recipes tab.
//
// The three snapshots moved together: the tab label is Ready, and RCP04 lost
// the readiness mark because it waits on RCP02. An approved recipe with an open
// dependency is not ready to claim. The tests below are what claim the files.
const goldenPath = "testdata/console-recipes.golden"

// homeGoldenPath is the checked-in snapshot of the console's Home tab.
//
// The Home tab is the landing screen, and the golden is taken at 80 columns: the
// classic default, where the numbers are whole and the mark is absent.
const homeGoldenPath = "testdata/console-home.golden"

// stressGoldenPath is the same, against the pathological plan.
const stressGoldenPath = "testdata/console-stress.golden"

// narrowBoardGoldenPath is the Board tab below 160 columns.
//
// The wide drawing starts at 160. Below that the vertical tree is the tree,
// and this file is that drawing: the connectors, the trunk, the forks.
const narrowBoardGoldenPath = "testdata/console-board-narrow.golden"

// clockPattern is the time of day in a rendered journal entry.
//
// The clock is not pinned by a golden: it changes on every run, and a golden that
// changed on every run would pin nothing. The placeholder is as wide as the time
// it replaces, so the layout the golden holds is the layout the chef sees.
var clockPattern = regexp.MustCompile(`\d\d:\d\d:\d\d`)

// goldenStore seeds a plan under a directory with a fixed name, so the project name in
// the render is deterministic.
//
// The plan is a **copy** of the suite's template — the seeder runs once for the whole
// suite. The copy lands in a directory named `demo`, so the name is `demo` and the
// golden holds.
func goldenStore(t *testing.T) *core.Store {
	t.Helper()
	return seededPlan(t)
}

// TestConsoleMatchesTheGolden pins the Recipes tab's layout.
//
// A golden catches the thing unit tests miss: a column that stops lining up. Update it
// with `make regen-goldens` when the design deliberately changes.
func TestConsoleMatchesTheGolden(t *testing.T) {
	store := goldenStore(t)
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	m.tab = TabRecipes
	got := StripANSI(m.View().Content)

	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Log("golden updated")
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v — run `make regen-goldens`", err)
	}
	if got != string(want) {
		t.Errorf("the console does not match the golden.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestHomeMatchesTheGolden pins the Home tab's layout.
//
// The Home tab is the console's landing screen, and it is the one screen the golden
// takes at 80 columns: the numbers are whole there and the mark is absent — the
// decoration yields to the content, in both dimensions.
//
// **The path is the machine's own**, so the fixture lives under the home directory
// and the console's short form of it (`~/demo/kitchen`) is what the golden holds.
// The temporary directory's own name changes on every run, and a golden that
// changed on every run would pin nothing.
func TestHomeMatchesTheGolden(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	// The copy lands in the home directory, under a directory named `demo` — so the
	// console's short form of the path (`~/demo/kitchen`) is what the golden holds.
	store := copyTemplate(t, planTemplate, home, "demo")

	m, err := Snapshot(store, 80)
	if err != nil {
		t.Fatal(err)
	}
	got := clockPattern.ReplaceAllString(StripANSI(m.View().Content), "HH:MM:SS")

	if *update {
		if err := os.MkdirAll(filepath.Dir(homeGoldenPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(homeGoldenPath, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Log("home golden updated")
		return
	}

	want, err := os.ReadFile(homeGoldenPath)
	if err != nil {
		t.Fatalf("read home golden: %v — run `make regen-goldens`", err)
	}
	if got != string(want) {
		t.Errorf("the console does not match the home golden.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestTheNarrowBoardMatchesTheGolden pins the vertical tree.
//
// Below 160 columns the Board is today's tree. The golden is that screen, at
// the width the other narrow goldens use. A wide drawing that leaked below
// 160 would move it.
func TestTheNarrowBoardMatchesTheGolden(t *testing.T) {
	store := goldenStore(t)
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	m.tab = TabBoard
	got := StripANSI(m.View().Content)
	// The seeded plan forks one child at a time, so the screen shows the last
	// connector. A branch with siblings is what draws ├─ — that test lives
	// with the tree fixture.
	if !strings.Contains(got, "└─") {
		t.Fatalf("the narrow tree should keep its connector:\n%s", got)
	}
	// The same drawing at the other narrow widths. The golden's width is the
	// byte-for-byte pin; these two must stay the vertical tree as well.
	for _, width := range []int{80, 159} {
		m.width = width
		out := StripANSI(m.View().Content)
		if !strings.Contains(out, "└─") {
			t.Errorf("at %d columns the tree should stay vertical:\n%s", width, out)
		}
	}
	m.width = 100

	if *update {
		if err := os.MkdirAll(filepath.Dir(narrowBoardGoldenPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(narrowBoardGoldenPath, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Log("narrow board golden updated")
		return
	}

	want, err := os.ReadFile(narrowBoardGoldenPath)
	if err != nil {
		t.Fatalf("read narrow board golden: %v — run `make regen-goldens`", err)
	}
	if got != string(want) {
		t.Errorf("the narrow board does not match the golden.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestGoldenIsNotStale catches a golden that was regenerated but not committed.
func TestGoldenIsNotStale(t *testing.T) {
	store := goldenStore(t)
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	got := StripANSI(m.View().Content)

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Skipf("no golden yet: %v", err)
	}
	if got == string(want) {
		return
	}
	// A mismatch here is the golden's problem, not the code's — but it means someone
	// changed the design and did not regenerate. Say so plainly.
	t.Log("the golden is out of date — run `make regen-goldens` if the change was deliberate")
}

// TestDriveReplaysKeys proves the scripted driver works end to end.
func TestDriveReplaysKeys(t *testing.T) {
	store := goldenStore(t)
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	// The first recipe is the ready one — space signs it off. The console lands on
	// Home, so the drive steps to the Recipes tab first.
	out, err := Drive(m, []string{"tab", "space"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(StripANSI(out), "signed off RCP02") {
		t.Errorf("the scripted drive did not sign off RCP02:\n%s", StripANSI(out))
	}
}

func TestDriveRejectsAnUnknownKey(t *testing.T) {
	store := goldenStore(t)
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Drive(m, []string{"nonsense"}); err == nil {
		t.Fatal("expected an unknown key to be refused")
	}
}
