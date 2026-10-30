package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/brand"
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

// TestTheNameHasOneAnswer proves the two nameless cases answer the same.
//
// A plan whose project file names nothing and a plan whose project file is missing
// are the same case, and they answer the same — and the answer is a **name**, never
// a path.
func TestTheNameHasOneAnswer(t *testing.T) {
	store := identityStore(t)
	project := filepath.Join(store.Root, core.FileProject)

	// A plan whose project file carries no name.
	if err := os.WriteFile(project,
		[]byte(`{"name":"","created":"2026-01-01T00:00:00Z"}`), 0o644); err != nil {
		t.Fatalf("write project: %v", err)
	}
	empty := loadInto(t, New(store))
	emptyName := empty.projectName()

	// The same plan without a project file at all.
	if err := os.Remove(project); err != nil {
		t.Fatalf("remove project: %v", err)
	}
	missing := loadInto(t, New(store))
	missingName := missing.projectName()

	if emptyName != missingName {
		t.Errorf("the two nameless cases answer differently: %q and %q", emptyName, missingName)
	}
	if emptyName != "demo" {
		t.Errorf("the name is %q, want the project's directory name %q", emptyName, "demo")
	}
	if strings.ContainsRune(emptyName, filepath.Separator) {
		t.Errorf("the name is a path: %q", emptyName)
	}
	if got, want := empty.View().WindowTitle, "maitre · "+emptyName; got != want {
		t.Errorf("the window title is %q, want %q", got, want)
	}
	if got, want := missing.View().WindowTitle, "maitre · "+emptyName; got != want {
		t.Errorf("the window title is %q, want %q", got, want)
	}
}

// TestTheNameIsSanitised proves the name goes through the same door as every other
// rendered string.
//
// The plan's own file is written by the tool — `Init` takes the directory's
// basename — and a directory may carry a newline. A name that broke the header
// would break the screen.
func TestTheNameIsSanitised(t *testing.T) {
	store := identityStore(t)
	raw := `{"name":"we\nird","created":"2026-01-01T00:00:00Z"}`
	if err := os.WriteFile(filepath.Join(store.Root, core.FileProject), []byte(raw), 0o644); err != nil {
		t.Fatalf("write project: %v", err)
	}
	m := loadInto(t, New(store))

	name := m.projectName()
	if strings.ContainsAny(name, "\n\r\t") {
		t.Errorf("the name still carries a control character: %q", name)
	}
	first := strings.SplitN(plain(m), "\n", 2)[0]
	if !strings.Contains(first, "we ird") {
		t.Errorf("the header does not carry the sanitised name: %q", first)
	}
	// The header is one row: a name with a newline would take two, and the screen
	// would be corrupt.
	clean := loadInto(t, New(identityStore(t)))
	if got, want := len(strings.Split(plain(m), "\n")), len(strings.Split(plain(clean), "\n")); got != want {
		t.Errorf("the header is %d rows, want %d\n%s", got, want, plain(m))
	}
	// The file on disk keeps what was written: the sanitising is a rendering
	// concern, not a claim about the data.
	onDisk, err := os.ReadFile(filepath.Join(store.Root, core.FileProject))
	if err != nil {
		t.Fatalf("read project: %v", err)
	}
	if string(onDisk) != raw {
		t.Errorf("the file on disk changed:\n%s", onDisk)
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

// TestTheSettingsTabShowsThePlanRoot proves the Settings tab says which plan this
// is.
//
// The row is labelled `root`, and the root is the plan directory — the one
// `core.FindRoot` returns and `store.Root` holds.
func TestTheSettingsTabShowsThePlanRoot(t *testing.T) {
	store := identityStore(t)
	m := loadInto(t, New(store))
	m.tab = TabSettings

	var text strings.Builder
	for _, line := range m.settingsLines() {
		text.WriteString(line.text + "\n")
	}
	rows := text.String()

	if want := humanPath(store.Root); !strings.Contains(rows, want) {
		t.Errorf("the Settings tab does not show the plan's root %q\n%s", want, rows)
	}
	if unwanted := humanPath(filepath.Dir(store.Root)); strings.Contains(rows, unwanted+"\n") {
		t.Errorf("the Settings tab shows the project's directory %q instead of the root\n%s", unwanted, rows)
	}
	for _, p := range palettes() {
		if !strings.Contains(rows, p.name) {
			t.Errorf("the Settings tab lost the theme %q\n%s", p.name, rows)
		}
	}
}

// TestTheRootRowIsOneRow proves a path with a newline in it does not split the row.
//
// The path is not the plan's content — but the console draws it on one line, and a
// directory whose name carries a newline would take two.
func TestTheRootRowIsOneRow(t *testing.T) {
	// A plan whose directory name carries a newline.
	root := filepath.Join(t.TempDir(), "we\nird", core.PlanDirName)
	store := core.Open(root)
	if err := core.Seed(store); err != nil {
		t.Fatalf("seed: %v", err)
	}
	m := loadInto(t, New(store))
	m.tab = TabSettings

	var text strings.Builder
	for _, line := range m.settingsLines() {
		text.WriteString(line.text + "\n")
	}
	rows := text.String()

	for _, line := range strings.Split(strings.TrimRight(rows, "\n"), "\n") {
		if strings.Contains(line, "\n") || strings.Contains(line, "\r") {
			t.Errorf("the root row is split: %q", line)
		}
	}
	if !strings.Contains(rows, "we ird") {
		t.Errorf("the root row does not carry the sanitised path\n%s", rows)
	}
	if strings.Contains(rows, "we\nird") {
		t.Errorf("the root row carries a newline\n%s", rows)
	}
}

// TestHumanPathShortensTheHomeDirectory proves the path is formatted for a human.
//
// The home directory is the same in every one of the chef's paths, and the part
// that is the same is the part that says nothing.
func TestHumanPathShortensTheHomeDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if got, want := humanPath(filepath.Join(home, "Code", "maitre")), "~/Code/maitre"; got != want {
		t.Errorf("humanPath is %q, want %q", got, want)
	}
	if got, want := humanPath(home), "~"; got != want {
		t.Errorf("humanPath is %q, want %q", got, want)
	}
	if got, want := humanPath("/somewhere/else"), "/somewhere/else"; got != want {
		t.Errorf("humanPath is %q, want %q", got, want)
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
	if !strings.Contains(header, " "+brand.Compact+" maitre") {
		t.Errorf("the header lost the mark and the product's name: %q", header)
	}
	if !strings.Contains(header, " · ") {
		t.Errorf("the header lost the separator: %q", header)
	}
}
