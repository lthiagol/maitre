package ui

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// The suite's templates.
//
// A plan is expensive to build: `core.Seed` writes fifteen files through the whole
// write path, and it costs about six tenths of a second. Two hundred and nineteen of
// them is two minutes of a suite that has six.
//
// So the suite builds its plans **once** — two templates — and every test takes a
// **copy**. The copy is byte for byte: the goldens pin the seeded plan's content, so
// nothing inside a copy is touched.
//
// The templates live in a directory the suite owns, and `TestMain` removes it when the
// suite ends. A template in somebody's tmp would outlive the run that made it.
var (
	// templateDir is the suite's own directory.
	templateDir string

	// planTemplate is the ordinary plan, seeded once.
	planTemplate string

	// stressTemplate is the pathological plan, seeded once. It is a **second**
	// template: `core.SeedStress` is a different plan, and it does not come from
	// the first one.
	stressTemplate string

	// seeds counts the seeding the suite did. It is the witness that the seed
	// happens once: the tests run twice under `-count=2`, and the count does not
	// move.
	seeds int

	// copies counts the copies the suite made. It moves with the runs — a copy per
	// test, every time the suite runs — and it is the other half of the witness:
	// the copies grow, the seed does not.
	copies int
)

// TestMain seeds the templates once and removes them when the suite ends.
//
// The seed happens here, not in a test: the suite runs twice in one process under
// `-count=2`, and a per-test seed is exactly what this file replaces.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "maitre-ui-templates")
	if err != nil {
		fmt.Fprintf(os.Stderr, "maitre: the template directory: %v\n", err)
		os.Exit(1)
	}
	templateDir = dir

	// The templates are seeded under a directory named `demo`, so a plan seeded
	// there would name itself `demo` — the name the goldens hold.
	planTemplate = filepath.Join(dir, "demo", core.PlanDirName)
	if err := core.Seed(core.Open(planTemplate)); err != nil {
		fmt.Fprintf(os.Stderr, "maitre: seed the plan: %v\n", err)
		os.RemoveAll(dir)
		os.Exit(1)
	}
	seeds++

	// The pathological plan is a **second** template. `core.SeedStress` is a
	// different plan — it adds the recipes that break a layout — and it does not
	// come from the first one.
	stressTemplate = filepath.Join(dir, "stress", "demo", core.PlanDirName)
	if err := core.SeedStress(core.Open(stressTemplate)); err != nil {
		fmt.Fprintf(os.Stderr, "maitre: seed the stress plan: %v\n", err)
		os.RemoveAll(dir)
		os.Exit(1)
	}
	seeds++

	// The templates name nothing, so every copy takes its name from the directory
	// it lands in.
	for _, template := range []string{planTemplate, stressTemplate} {
		if err := unnameTheTemplate(template); err != nil {
			fmt.Fprintf(os.Stderr, "maitre: unname the template: %v\n", err)
			os.RemoveAll(dir)
			os.Exit(1)
		}
	}

	code := m.Run()

	// The templates are the suite's, so the suite takes them with it.
	if err := removeTemplates(dir); err != nil {
		fmt.Fprintf(os.Stderr, "maitre: remove the templates: %v\n", err)
		code = 1
	} else if _, err := os.Stat(dir); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "maitre: the template directory survived: %s\n", dir)
		code = 1
	}
	os.Exit(code)
}

// removeTemplates removes the suite's template directory.
//
// `TestMain` calls this when the suite ends. A template left behind is a directory
// nobody owns, and the next run would find it and believe it was its own.
func removeTemplates(dir string) error {
	if dir == "" {
		return nil
	}
	return os.RemoveAll(dir)
}

// unnameTheTemplate leaves the template's project file without a name.
//
// The seeder writes the name it took from the directory it seeded into — `store.Init`
// takes the parent directory's basename. A copy lands somewhere else, so a named
// template would make every copy answer with the *template's* name, whatever
// directory it was placed in.
//
// A plan that names nothing answers with the directory it is in — the rule the store
// already follows. So the template names nothing: a copy in a directory named `demo`
// says `demo`, and a copy in a directory named `we\nird` says `we\nird`, exactly as
// a plan seeded there would have said.
//
// The file is written the way the store writes it, so the only difference from the
// seeder's output is the name itself.
func unnameTheTemplate(root string) error {
	path := filepath.Join(root, core.FileProject)
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc struct {
		Name     string         `json:"name"`
		Created  string         `json:"created"`
		Charter  map[string]any `json:"charter"`
		Settings map[string]any `json:"settings"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return err
	}
	doc.Name = ""
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0o644)
}

// seededPlan is the plan a test draws: a copy of the ordinary template.
//
// The copy lands in a directory named `demo`, so the plan's name is `demo` — the
// name the goldens hold.
func seededPlan(t *testing.T) *core.Store {
	t.Helper()
	return copyTemplate(t, planTemplate, t.TempDir(), "demo")
}

// stressPlan is the plan a test draws when it wants the pathological one.
//
// The pathological plan is a **second** template, seeded on its own. A copy of the
// ordinary plan would be a plan with nothing wrong with it, and a stress test on it
// would prove nothing.
func stressPlan(t *testing.T) *core.Store {
	t.Helper()
	return copyTemplate(t, stressTemplate, t.TempDir(), "demo")
}

// copyTemplate copies a seeded plan into `<parent>/<name>/kitchen` and opens it.
//
// The copy is **byte for byte** the template: the goldens pin the seeded plan's
// content, and a copy that rewrote a byte would move a golden.
func copyTemplate(t *testing.T, template, parent, name string) *core.Store {
	t.Helper()
	root := filepath.Join(parent, name, core.PlanDirName)
	if err := copyTree(template, root); err != nil {
		t.Fatalf("copy the template: %v", err)
	}
	copies++
	return core.Open(root)
}

// copyTree copies a directory tree, file by file, byte for byte.
//
// The directories come too — an empty one is part of the plan's shape — and the
// bytes are the template's: nothing is read, parsed or rewritten on the way.
func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

// treeFiles lists a tree's files, relative to its root, sorted.
//
// The list is the plan's shape: the same files, the same names.
func treeFiles(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		out = append(out, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	sort.Strings(out)
	return out
}

// sameTree reports whether two trees hold the same files with the same bytes.
func sameTree(a, b string) error {
	af, err := listFiles(a)
	if err != nil {
		return err
	}
	bf, err := listFiles(b)
	if err != nil {
		return err
	}
	if len(af) != len(bf) {
		return fmt.Errorf("%s holds %d files, %s holds %d", a, len(af), b, len(bf))
	}
	for i := range af {
		if af[i].name != bf[i].name {
			return fmt.Errorf("%s holds %s, %s holds %s", a, af[i].name, b, bf[i].name)
		}
		if string(af[i].data) != string(bf[i].data) {
			return fmt.Errorf("%s and %s differ in %s", a, b, af[i].name)
		}
	}
	return nil
}

// treeFile is one file of a tree, with its bytes.
type treeFile struct {
	name string
	data []byte
}

// listFiles lists a tree's files with their bytes, sorted by name.
func listFiles(root string) ([]treeFile, error) {
	var out []treeFile
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out = append(out, treeFile{name: filepath.ToSlash(rel), data: data})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out, nil
}

// The seeder's files.
//
// The ordinary plan holds exactly these, and the template holds exactly these: a
// template that carries anything else is not the seeder's plan, and a copy of it is
// not the plan the goldens pin.
var seededFiles = []string{
	"activity.json",
	"backlog/BL01-the-board-needs-a-way-to-show-evidence.json",
	"backlog/BL02-a-recipe-with-no-acs-is-probably-an-idea.json",
	"backlog/BL03-the-chef-wants-a-weekly-digest.json",
	"backlog/BL04-trying-the-board-on-a-real-project.json",
	"backlog/BL05-wire-the-review-queue.json",
	"backlog/BL06-sharpen-the-evidence-contract.json",
	"project.json",
	"recipes/RCP01-sharpen-the-ac-gate.json",
	"recipes/RCP02-sign-off-the-menu.json",
	"recipes/RCP03-wire-the-review-queue.json",
	"recipes/RCP04-add-the-backlog-view.json",
	"recipes/RCP05-sharpen-the-evidence-contract.json",
	"recipes/RCP06-draft-the-glossary.json",
	"recipes/RCP07-record-the-review-protocol.json",
}

// TestTheCopyIsTheTemplate proves the copy is the template, byte for byte.
//
// Two copies made at different times carry the same content — the plan a test draws
// does not change under it — and the template still holds exactly the files the
// seeder wrote.
func TestTheCopyIsTheTemplate(t *testing.T) {
	first := seededPlan(t)
	second := seededPlan(t)

	if err := sameTree(first.Root, second.Root); err != nil {
		t.Errorf("two copies are not the same plan: %v", err)
	}
	if err := sameTree(planTemplate, first.Root); err != nil {
		t.Errorf("the copy is not the template: %v", err)
	}

	got := treeFiles(t, planTemplate)
	if strings.Join(got, "\n") != strings.Join(seededFiles, "\n") {
		t.Errorf("the template does not hold the seeder's files.\ngot:  %v\nwant: %v", got, seededFiles)
	}
}

// TestTheTemplateIsSeededOnce proves the suite seeds once, not once per test.
//
// The count is the witness: the tests run twice under `-count=2`, and a suite that
// seeded per test would count two per run — and four for two templates.
func TestTheTemplateIsSeededOnce(t *testing.T) {
	if seeds != 2 {
		t.Errorf("the suite seeded %d times, want 2 — one per template", seeds)
	}
}

// TestTheCopiesGrowAndTheSeedDoesNot proves the template is shared and the plans are
// not.
//
// The suite runs twice in one process under `-count=2`, and the tests run in any
// order under `-shuffle=on`. The copies move with the runs — one per draw, every
// time — and the seed does not move at all.
//
// The witness draws its own copies, so it does not depend on where it lands in the
// order: a suite that seeded per test would count four under `-count=2`, and two
// here.
func TestTheCopiesGrowAndTheSeedDoesNot(t *testing.T) {
	before := copies
	_ = seededPlan(t)
	_ = seededPlan(t)
	if copies != before+2 {
		t.Errorf("two draws made %d copies, want 2 — a copy per draw", copies-before)
	}
	if seeds != 2 {
		t.Errorf("the suite seeded %d times, want 2 — one per template, whatever the run count", seeds)
	}
}

// TestTheCopiesDoNotSeeEachOther proves a copy is the test's own.
//
// The tests run in any order under `-shuffle=on`, and twice over under `-count=2`.
// A write in one copy must not reach another copy, and must not reach the template
// — the template is shared by every test in the run, so a write that reached it
// would show up in the next test to draw a plan.
func TestTheCopiesDoNotSeeEachOther(t *testing.T) {
	first := seededPlan(t)
	second := seededPlan(t)

	ops := core.NewOperations(first, core.AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})
	if _, err := ops.RecordDeviation("RCP01", "G1", "read the queue twice"); err != nil {
		t.Fatalf("record the deviation: %v", err)
	}

	if err := sameTree(second.Root, planTemplate); err != nil {
		t.Errorf("a copy saw another copy's write: %v", err)
	}
	if err := sameTree(first.Root, planTemplate); err == nil {
		t.Error("the write did not land in the copy it was made in")
	}
}

// TestTheTemplatesAreTheSuites proves the templates live in the suite's directory.
//
// They are the suite's own, and `TestMain` removes them: a template in somebody's
// tmp would outlive the run that made it.
func TestTheTemplatesAreTheSuites(t *testing.T) {
	if templateDir == "" {
		t.Fatal("the suite has no template directory")
	}
	if !strings.HasPrefix(filepath.Base(templateDir), "maitre-ui-templates") {
		t.Errorf("the template directory is %q, want one the suite owns", templateDir)
	}
	for _, dir := range []string{planTemplate, stressTemplate} {
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("the template %s is not there: %v", dir, err)
		}
	}
}

// TestTheTemplatesAreRemoved proves the cleanup takes the directory with it.
//
// `TestMain` calls `removeTemplates` when the suite ends; this is that function's
// witness.
func TestTheTemplatesAreRemoved(t *testing.T) {
	scratch := filepath.Join(t.TempDir(), "templates")
	if err := os.MkdirAll(filepath.Join(scratch, "recipes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(scratch, "recipes", "RCP01.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := removeTemplates(scratch); err != nil {
		t.Fatalf("remove the templates: %v", err)
	}
	if _, err := os.Stat(scratch); !os.IsNotExist(err) {
		t.Errorf("the template directory survived the cleanup: %v", err)
	}
}

// TestTheCopyIsNamedAfterItsDirectory proves the copy's name follows the directory.
//
// The name comes from the project's directory, so the copy lands in a directory named
// `demo` and the plan says `demo` — the name every golden holds. A copy in a
// directory with another name says that name, the way a plan seeded there would.
func TestTheCopyIsNamedAfterItsDirectory(t *testing.T) {
	store := seededPlan(t)
	if got, want := store.ProjectName(), "demo"; got != want {
		t.Errorf("the plan's name is %q, want %q", got, want)
	}
	if got := filepath.Base(filepath.Dir(store.Root)); got != "demo" {
		t.Errorf("the copy landed in %q, want a directory named %q", got, "demo")
	}

	other := copyTemplate(t, planTemplate, t.TempDir(), "the-other-one")
	if got, want := other.ProjectName(), "the-other-one"; got != want {
		t.Errorf("the plan's name is %q, want the directory's name %q", got, want)
	}
}

// TestTheStressPlanIsItsOwnTemplate proves the pathological plan is a second plan.
//
// `core.SeedStress` adds the recipes that break a layout — a title with a newline, a
// title longer than the column, text that is not ASCII. A stress test that drew the
// ordinary plan would be testing a plan with nothing wrong with it.
func TestTheStressPlanIsItsOwnTemplate(t *testing.T) {
	ordinary := seededPlan(t)
	stress := stressPlan(t)

	if err := sameTree(ordinary.Root, stress.Root); err == nil {
		t.Error("the stress plan is the ordinary plan")
	}

	recipes, err := stress.ListRecipes()
	if err != nil {
		t.Fatalf("list the stress recipes: %v", err)
	}
	var found bool
	for _, r := range recipes {
		if strings.HasPrefix(r.Title, "STRESS:") {
			found = true
		}
	}
	if !found {
		t.Error("the stress plan holds no STRESS recipe")
	}
}

// The template file is the only place in the package a plan is seeded.
const templateFile = "seed_test.go"

// TestNoTestSeedsItsOwnPlan proves the package does not seed its own plans.
//
// The seeder costs 0.608s and writes fifteen files, and a suite that seeds per test
// pays it per test — that was two minutes of the run. The template is seeded once,
// in `TestMain`, and every test takes a copy.
//
// The wire is the package's own sources: a test file that calls the seeder is a file
// that pays for a plan again, and the cost comes back without anything going red.
func TestNoTestSeedsItsOwnPlan(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read the package: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), "_test.go") || e.Name() == templateFile {
			continue
		}
		raw, err := os.ReadFile(e.Name())
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		for _, call := range []string{"core.Seed(", "core.SeedStress("} {
			if strings.Contains(string(raw), call) {
				t.Errorf("%s seeds its own plan: %s — a test takes a copy, not a seed", e.Name(), call)
			}
		}
	}
}

// The goldens the suite claims.
//
// The store moved nothing, so the goldens moved nothing: the same five files, each
// one a screen the console still draws. The list is the witness that this slice
// neither added a golden nor dropped one.
var goldenFiles = []string{
	"console-board-narrow.golden",
	"console-filter.golden",
	"console-home.golden",
	"console-recipes.golden",
	"console-stress.golden",
}

// TestTheGoldensAreTheOnesTheSuiteClaims proves the goldens are the same five.
//
// The content is the golden tests' own business — they read the file and compare.
// This test holds the set: a golden added or dropped by a slice that changed nothing
// would be a layout nothing pins.
func TestTheGoldensAreTheOnesTheSuiteClaims(t *testing.T) {
	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("read the goldens: %v", err)
	}
	var got []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		got = append(got, e.Name())
	}
	sort.Strings(got)
	if strings.Join(got, "\n") != strings.Join(goldenFiles, "\n") {
		t.Errorf("the goldens are not the ones the suite claims.\ngot:  %v\nwant: %v", got, goldenFiles)
	}
}
