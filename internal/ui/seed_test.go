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
// So the suite builds the plan **once**, and every test takes a **copy**. The copy is byte for byte: the goldens
// pin the seeded plan's content, so nothing inside a copy is touched.
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

	// The templates are seeded under a directory named `demo`, so the plan's own
	// name is `demo` — and every copy of it says the same name.
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
