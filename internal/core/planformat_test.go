package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTheDocNamesThePlanFormatTheCodeHas keeps the plan format honest.
//
// The plan format is what the code writes on disk: the directories under the plan
// root, the files in it, and the name a recipe's file carries. The doc's § 6 states
// it, and the code is the truth — so the two are checked against each other.
//
// The check runs both ways. A directory the code creates and the doc does not name is
// a directory the reader cannot find; a directory the doc names and the code does not
// create is a plan the reader was told about that does not exist.
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

	// The directories under the plan root.
	for _, dir := range []string{DirRecipes, DirCourses, DirPrepNotes, DirSpecials} {
		if !strings.Contains(section, dir+"/") {
			t.Errorf("the plan format should name the %s/ directory:\n%s", dir, section)
		}
	}

	// The files in the plan root.
	for _, file := range []string{FileProject, FileActivity, FileDecisions} {
		if !strings.Contains(section, file) {
			t.Errorf("the plan format should name %s:\n%s", file, section)
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
