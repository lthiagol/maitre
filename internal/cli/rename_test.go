package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// planFiles lists the artifact files in a plan directory.
func planFiles(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		out = append(out, e.Name())
	}
	return out
}

// TestARenameDoesNotMultiplyThePlan covers the plan's index.
//
// The file name follows the title, so a rename used to write a second file and leave
// the first. The plan then held two artifacts with one id, and every reader picked
// one arbitrarily: `show` read the old title, `list` counted the recipe twice, and
// the board drew it twice.
func TestARenameDoesNotMultiplyThePlan(t *testing.T) {
	dir := newPlan(t)

	create := cli(t, dir, "recipe", "create",
		"--title", "Dep",
		"--outcome", "o",
		"--problem", "p",
		"--out-of-scope", "a",
		"--out-of-scope", "b")
	if !create.ok(t) {
		t.Fatalf("create failed:\n%s", create.Stderr)
	}
	update := cli(t, dir, "recipe", "update", "RCP01", "--title", "Renamed")
	if !update.ok(t) {
		t.Fatalf("update failed:\n%s", update.Stderr)
	}

	// One file: the directory is the plan's index.
	files := planFiles(t, filepath.Join(dir, "kitchen", "recipes"))
	if len(files) != 1 {
		t.Fatalf("the recipes directory should hold one file, got %v", files)
	}

	// `show` reads the new title.
	show := cli(t, dir, "recipe", "show", "RCP01")
	if !show.ok(t) {
		t.Fatalf("show failed:\n%s", show.Stderr)
	}
	recipe, _ := show.doc(t)["recipe"].(map[string]any)
	if recipe["title"] != "Renamed" {
		t.Fatalf("the title = %v, want Renamed", recipe["title"])
	}

	// `list` counts the recipe once.
	list := cli(t, dir, "recipe", "list")
	if !list.ok(t) {
		t.Fatalf("list failed:\n%s", list.Stderr)
	}
	if n := list.doc(t)["count"]; n != float64(1) {
		t.Fatalf("the list holds %v recipes, want 1", n)
	}
}

// TestARenameOfASpecialDoesNotMultiplyThePlan is the same, one rung down.
func TestARenameOfASpecialDoesNotMultiplyThePlan(t *testing.T) {
	dir := newPlan(t)

	create := cli(t, dir, "special", "create", "--title", "A hunch", "--body", "b")
	if !create.ok(t) {
		t.Fatalf("create failed:\n%s", create.Stderr)
	}
	update := cli(t, dir, "special", "update", "S01", "--title", "A better hunch")
	if !update.ok(t) {
		t.Fatalf("update failed:\n%s", update.Stderr)
	}

	files := planFiles(t, filepath.Join(dir, "kitchen", "specials"))
	if len(files) != 1 {
		t.Fatalf("the specials directory should hold one file, got %v", files)
	}
	list := cli(t, dir, "special", "list")
	if !list.ok(t) {
		t.Fatalf("list failed:\n%s", list.Stderr)
	}
	if n := list.doc(t)["count"]; n != float64(1) {
		t.Fatalf("the list holds %v specials, want 1", n)
	}
}

// TestARenameOfAPrepNoteDoesNotMultiplyThePlan is the same, one rung down again.
func TestARenameOfAPrepNoteDoesNotMultiplyThePlan(t *testing.T) {
	dir := newPlan(t)

	create := cli(t, dir, "prepnote", "create", "--title", "A plan", "--intent", "i")
	if !create.ok(t) {
		t.Fatalf("create failed:\n%s", create.Stderr)
	}
	update := cli(t, dir, "prepnote", "update", "N01", "--title", "A better plan")
	if !update.ok(t) {
		t.Fatalf("update failed:\n%s", update.Stderr)
	}

	files := planFiles(t, filepath.Join(dir, "kitchen", "prep-notes"))
	if len(files) != 1 {
		t.Fatalf("the prep-notes directory should hold one file, got %v", files)
	}
	list := cli(t, dir, "prepnote", "list")
	if !list.ok(t) {
		t.Fatalf("list failed:\n%s", list.Stderr)
	}
	if n := list.doc(t)["count"]; n != float64(1) {
		t.Fatalf("the list holds %v prep notes, want 1", n)
	}
}
