package core

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// planFormatLine is one line of the plan format the doc states.
//
// The doc's § 6 shows the plan root as a tree: `recipes/RCP01-<slug>.json` is a
// directory and a file in it, `project.json` is a file in the root.
var planFormatLine = regexp.MustCompile(`^  ([a-z][a-z0-9.-]*)(/([^ \t]+))?([ \t]|$)`)

// planFormatVersion is the version the doc's § 6 states the plan format is.
var planFormatVersion = regexp.MustCompile(`plan format is version (\d+)`)

// TestTheDocStatesThePlanFormatVersion keeps the doc's version and the code's.
//
// The doc states the plan format's version and the code holds it in
// `core.DomainVersion`. Nothing tied the two, so the doc could say 3 while the tool
// spoke 4 — and the version guard would fence off the wrong tool.
//
// The wire reads the number out of the doc's § 6 and compares it with the code.
func TestTheDocStatesThePlanFormatVersion(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "docs", "03-domain-model.md"))
	if err != nil {
		t.Fatal(err)
	}
	section := sectionOf(string(raw), "## 6. Plan format")
	if section == "" {
		t.Fatal("the doc has no § 6. Plan format")
	}
	m := planFormatVersion.FindStringSubmatch(section)
	if m == nil {
		t.Fatalf("the plan format should state its version:\n%s", section)
	}
	if m[1] != strconv.Itoa(DomainVersion) {
		t.Errorf("the doc says the plan format is version %s, the code says %d", m[1], DomainVersion)
	}
}

// TestTheDocNamesThePlanFormatTheCodeHas keeps the plan format honest.
//
// The plan format is what the code writes on disk: the directories under the plan
// root, the files in it, and the name a recipe's file carries. The doc's § 6 states
// it, and the code is the truth — so the two are checked against each other.
//
// **The check runs both ways.** A directory the code creates and the doc does not
// name is a directory the reader cannot find. A directory the doc names and the code
// does not create is a plan the reader was told about that does not exist — and the
// backlog move is exactly that: a doc promising `backlog/BL01-<slug>.json` for a
// directory nothing writes.
//
// The recipe file's name is the third fact: the writer emits `<id>-<slug>.json` with
// the id the writer allocates, and the doc must show that id form and not an older one.
func TestTheDocNamesThePlanFormatTheCodeHas(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "docs", "03-domain-model.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc := string(raw)
	section := sectionOf(doc, "## 6. Plan format")
	if section == "" {
		t.Fatal("the doc has no § 6. Plan format")
	}

	docDirs, docFiles := docPlanFormat(t, section)

	codeDirs := map[string]bool{
		DirRecipes: true,
		DirCourses: true,
		DirBacklog: true,
	}
	codeFiles := map[string]bool{
		FileProject:   true,
		FileActivity:  true,
		FileDecisions: true,
	}

	// The code's directories must be in the doc.
	for _, dir := range sortedKeys(codeDirs) {
		if !docDirs[dir] {
			t.Errorf("the plan format should name the %s/ directory:\n%s", dir, section)
		}
	}
	// The doc's directories must be in the code.
	for _, dir := range sortedKeys(docDirs) {
		if !codeDirs[dir] {
			t.Errorf("the plan format names a %s/ directory the code does not create:\n%s", dir, section)
		}
	}

	// The code's files must be in the doc.
	for _, file := range sortedKeys(codeFiles) {
		if !docFiles[file] {
			t.Errorf("the plan format should name %s:\n%s", file, section)
		}
	}
	// The doc's files must be in the code.
	for _, file := range sortedKeys(docFiles) {
		if !codeFiles[file] {
			t.Errorf("the plan format names %s, which the code does not write:\n%s", file, section)
		}
	}

	// The recipe file's name: the id the writer allocates, and the slug.
	next, err := (&Store{Root: t.TempDir()}).NextRecipeID()
	if err != nil {
		t.Fatal(err)
	}
	want := DirRecipes + "/" + next + "-<slug>.json"
	if !strings.Contains(section, want) {
		t.Errorf("the plan format should name %q — the writer allocates %s:\n%s", want, next, section)
	}
}

// docPlanFormat reads the directories and the files the doc's plan format names.
//
// A file's name is a *template* when it carries a `<` — `RCP01-<slug>.json` is the
// shape of a recipe's file, not a name the plan holds. A template is checked against
// the code by its directory; a literal name is checked against the code's files.
func docPlanFormat(t *testing.T, section string) (dirs, files map[string]bool) {
	t.Helper()
	dirs = map[string]bool{}
	files = map[string]bool{}
	for _, line := range strings.Split(section, "\n") {
		m := planFormatLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		name, sub, file := m[1], m[2], m[3]
		if sub == "" {
			// A file in the plan root.
			files[name] = true
			continue
		}
		// A file in a directory: the directory is a fact, the file is a template.
		dirs[name] = true
		if !strings.Contains(file, "<") {
			files[file] = true
		}
	}
	if len(dirs) == 0 && len(files) == 0 {
		t.Fatalf("the plan format section names no directories and no files:\n%s", section)
	}
	return dirs, files
}

// sortedKeys is a map's keys in a stable order, so the failures read the same way
// every run.
func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
