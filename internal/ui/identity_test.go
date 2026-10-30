package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// identityStore seeds a plan under a directory with a fixed name, so the project's
// name in the render is deterministic.
func identityStore(t *testing.T) *core.Store {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := core.Seed(store); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return store
}

// TestWindowTitleCarriesTheProject proves the terminal's tab says which project it
// is.
//
// The window title is not in the rendered content, so `--render` cannot see it: a
// test is the only instrument here.
func TestWindowTitleCarriesTheProject(t *testing.T) {
	store := identityStore(t)
	m := loadInto(t, New(store))

	if got, want := m.View().WindowTitle, "maitre · demo"; got != want {
		t.Errorf("the window title is %q, want %q", got, want)
	}
}

// TestTheWindowTitleAndTheHeaderAgree proves the name has one source.
//
// Two computations would drift, and the drift would show as a terminal whose tab
// says one project and whose header says another.
func TestTheWindowTitleAndTheHeaderAgree(t *testing.T) {
	store := identityStore(t)
	m := loadInto(t, New(store))

	name := m.projectName()
	if !strings.Contains(m.View().WindowTitle, name) {
		t.Errorf("the window title %q does not carry the name %q", m.View().WindowTitle, name)
	}
	if !strings.Contains(plain(m), "· "+name) {
		t.Errorf("the header does not carry the name %q\n%s", name, plain(m))
	}
}

// TestTheNameFallsBackToThePlanRoot proves a plan with no name still has one.
//
// The plan's own name is the answer; without one, the console falls back to the
// plan's root, so the tab is never a bare "maitre".
func TestTheNameFallsBackToThePlanRoot(t *testing.T) {
	store := identityStore(t)
	// A plan whose project file carries no name.
	if err := os.WriteFile(filepath.Join(store.Root, core.FileProject),
		[]byte(`{"name":"","created":"2026-01-01T00:00:00Z"}`), 0o644); err != nil {
		t.Fatalf("write project: %v", err)
	}
	m := loadInto(t, New(store))

	if got, want := m.projectName(), store.Root; got != want {
		t.Errorf("the name is %q, want the plan's root %q", got, want)
	}
	if got, want := m.View().WindowTitle, "maitre · "+store.Root; got != want {
		t.Errorf("the window title is %q, want %q", got, want)
	}
	if !strings.Contains(plain(m), "· "+store.Root) {
		t.Errorf("the header does not carry the fallback\n%s", plain(m))
	}
}

// TestTheProjectOverrideReachesTheConsole proves `--project` names the plan for the
// session.
//
// The flag is what a chef with two checkouts of one project needs, and it is the
// answer to two plans whose directories share a basename.
func TestTheProjectOverrideReachesTheConsole(t *testing.T) {
	store := identityStore(t)
	m := loadInto(t, New(store, WithProject("the-other-one")))

	if got, want := m.View().WindowTitle, "maitre · the-other-one"; got != want {
		t.Errorf("the window title is %q, want %q", got, want)
	}
	if !strings.Contains(plain(m), "· the-other-one") {
		t.Errorf("the header does not carry the override\n%s", plain(m))
	}
	if strings.Contains(plain(m), "· demo") {
		t.Errorf("the header still carries the plan's own name\n%s", plain(m))
	}
}

// TestEveryConstructorTakesTheOverride proves the flag reaches all three.
//
// `New`, `NewLive` and `NewScripted` are the console's three doors, and a flag that
// reaches two of them reaches the wrong two.
func TestEveryConstructorTakesTheOverride(t *testing.T) {
	store := identityStore(t)

	for name, build := range map[string]func(*core.Store, ...Option) *Model{
		"New":         New,
		"NewLive":     NewLive,
		"NewScripted": NewScripted,
	} {
		m := loadInto(t, build(store, WithProject("named")))
		if got, want := m.View().WindowTitle, "maitre · named"; got != want {
			t.Errorf("%s: the window title is %q, want %q", name, got, want)
		}
	}
}

// TestTheProjectOverrideWritesNothing proves the flag is display only.
//
// It must not touch the plan's metadata, and it must not enter the journal: a name
// the chef typed is not a fact about the plan.
func TestTheProjectOverrideWritesNothing(t *testing.T) {
	store := identityStore(t)
	project := filepath.Join(store.Root, core.FileProject)
	activity := filepath.Join(store.Root, core.FileActivity)
	before := map[string][]byte{}
	for _, path := range []string{project, activity} {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		before[path] = raw
	}

	m := loadInto(t, New(store, WithProject("zz-override-name")))
	_ = m.View()

	for path, raw := range before {
		now, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if string(now) != string(raw) {
			t.Errorf("%s changed:\nbefore: %s\nafter:  %s", path, raw, now)
		}
	}

	journal, err := store.LoadJournal()
	if err != nil {
		t.Fatalf("load journal: %v", err)
	}
	for _, e := range journal.Events {
		if strings.Contains(fmt.Sprint(e), "zz-override-name") {
			t.Errorf("the override entered the journal: %+v", e)
		}
	}
}

// TestThePlanKeepsItsOwnName proves the override does not change the plan.
//
// The console shows what the chef asked for; the plan still says what it says.
func TestThePlanKeepsItsOwnName(t *testing.T) {
	store := identityStore(t)
	m := loadInto(t, New(store, WithProject("zz-override-name")))

	p, err := m.ops.GetProject()
	if err != nil {
		t.Fatalf("get project: %v", err)
	}
	if p.Name != "demo" {
		t.Errorf("the plan's name is %q, want %q", p.Name, "demo")
	}
}

// TestTheHeaderNameIsNotAFootnote proves the project's name reads as a name.
//
// A name in the faint style reads as a footnote, and a footnote is not an identity.
func TestTheHeaderNameIsNotAFootnote(t *testing.T) {
	store := identityStore(t)
	m := loadInto(t, New(store))

	header := m.header()
	name := m.projectName()
	if !strings.Contains(header, m.styAcc.Render(name)) {
		t.Errorf("the header does not render the name in the accent style: %q", header)
	}
	if strings.Contains(header, m.styFaint.Render(name)) {
		t.Errorf("the header renders the name in the faint style: %q", header)
	}
}

// TestTheHeaderKeepsTheProductsName proves the slice moved only the project's name.
//
// The mark, the product's name and the separator are the console's own; a slice that
// renames the project and also renames the product is doing more than it said.
func TestTheHeaderKeepsTheProductsName(t *testing.T) {
	store := identityStore(t)
	m := loadInto(t, New(store))

	header := m.header()
	if !strings.Contains(header, " ⋈ maitre") {
		t.Errorf("the header lost the mark and the product's name: %q", header)
	}
	if !strings.Contains(header, " · ") {
		t.Errorf("the header lost the separator: %q", header)
	}
}
