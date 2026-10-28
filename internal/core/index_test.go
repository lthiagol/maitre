package core

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// filesIn lists the files in a plan directory.
func filesIn(t *testing.T, dir string) []string {
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

// TestARenameLeavesOneFile covers the plan's index.
//
// The file name carries the slug and the slug follows the title, so a rename used to
// write a second file and leave the first: two artifacts with one id. Every reader
// then picked one arbitrarily — `show` read the old title, `list` counted the recipe
// twice, and the board drew it twice.
func TestARenameLeavesOneFile(t *testing.T) {
	ops, store := newOps(t)
	r, err := ops.CreateRecipe(spec())
	if err != nil {
		t.Fatal(err)
	}
	title := "A different title"
	if _, err := ops.UpdateRecipe(r.ID, Fragment{Title: &title}); err != nil {
		t.Fatal(err)
	}

	files := filesIn(t, store.RecipeDir())
	if len(files) != 1 {
		t.Fatalf("the directory should hold one file, got %v", files)
	}
	got, err := store.LoadRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != title {
		t.Fatalf("the recipe's title = %q, want %q", got.Title, title)
	}
	list, err := store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("the plan should hold one recipe, got %d", len(list))
	}
	if list[0].Title != title {
		t.Fatalf("the listed recipe's title = %q, want %q", list[0].Title, title)
	}
}

// TestARenameOfABacklogItemLeavesOneFile is the same, one rung down.
func TestARenameOfABacklogItemLeavesOneFile(t *testing.T) {
	ops, store := newOps(t)
	item, err := ops.CreateBacklogItem("A hunch", "a body", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.UpdateBacklogItem(item.ID, "A better hunch", "a body", nil, "", nil); err != nil {
		t.Fatal(err)
	}

	files := filesIn(t, filepath.Join(store.Root, DirBacklog))
	if len(files) != 1 {
		t.Fatalf("the directory should hold one file, got %v", files)
	}
	got, err := store.LoadBacklogItem(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "A better hunch" {
		t.Fatalf("the item's title = %q, want %q", got.Title, "A better hunch")
	}
	list, err := store.ListBacklogItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("the plan should hold one item, got %d", len(list))
	}
}

// TestTheReaderPicksTheCurrentCopy covers a plan that already holds two copies.
//
// The writer removes the stale copy, so a plan written by the fixed tool never holds
// two. A plan written before the fix does — and then every reader has to choose, so
// the choice must be the plan's: the copy the artifact itself says is newest.
//
// Two cases, because the stamps have second resolution and the common tie is two
// copies written in the same second: an older stamp, and an equal one.
func TestTheReaderPicksTheCurrentCopy(t *testing.T) {
	for _, tc := range []struct {
		name    string
		updated string
	}{
		{"an older stamp", "2020-01-01T00:00:00Z"},
		{"an equal stamp", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ops, store := newOps(t)
			r, err := ops.CreateRecipe(spec())
			if err != nil {
				t.Fatal(err)
			}
			// The stale copy: the same recipe under another name, as a plan written
			// by the old build holds it. Its name is different from the current
			// file's, or the test would be writing the current file.
			stale := filepath.Join(store.RecipeDir(), r.ID+"-an-older-title.json")
			if err := writeCopy(t, store, *r, stale, "An older title", tc.updated); err != nil {
				t.Fatal(err)
			}
			// The stale copy was written first: the file time says so.
			old := time.Now().Add(-time.Hour)
			if err := os.Chtimes(stale, old, old); err != nil {
				t.Fatal(err)
			}

			// The reader picks the current copy, not the one that sorts first.
			got, err := store.LoadRecipe(r.ID)
			if err != nil {
				t.Fatal(err)
			}
			if got.Title != r.Title {
				t.Fatalf("the reader picked %q, want %q", got.Title, r.Title)
			}
			// And the listing counts the recipe once.
			list, err := store.ListRecipes()
			if err != nil {
				t.Fatal(err)
			}
			if len(list) != 1 {
				t.Fatalf("the plan should hold one recipe, got %d", len(list))
			}

			// The next write cleans the plan up: the stale copy is gone.
			title := "A third title"
			if _, err := ops.UpdateRecipe(r.ID, Fragment{Title: &title}); err != nil {
				t.Fatal(err)
			}
			files := filesIn(t, store.RecipeDir())
			if len(files) != 1 {
				t.Fatalf("the directory should hold one file, got %v", files)
			}
		})
	}
}

// writeCopy writes a copy of a recipe under another name, as a plan written by the
// old build holds it.
//
// An empty `updated` keeps the copy's own stamp, which is the case the tie-break has
// to survive: the stamps are equal and only the file's time tells them apart.
func writeCopy(t *testing.T, store *Store, r Recipe, path, title, updated string) error {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(store.RecipeDir(), r.ID+"-"+Slugify(r.Title)+".json"))
	if err != nil {
		return err
	}
	var copy Recipe
	if err := unmarshal(raw, &copy); err != nil {
		return err
	}
	copy.Title = title
	copy.Slug = Slugify(title)
	if updated != "" {
		copy.Updated = updated
	}
	out, err := json.MarshalIndent(copy, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(out, '\n'), 0o644)
}

// TestAFailedRemovalDoesNotFailTheWrite covers the write that lands and the cleanup
// that does not.
//
// The write comes first, so the artifact is there; removing the copy the old name
// left behind can fail — an immutable file, a read-only mount. The operation must
// still succeed: the artifact changed, and the journal is the only event log, so a
// failure would leave the plan and the journal disagreeing. The leftover stays, and
// the next write cleans it up.
func TestAFailedRemovalDoesNotFailTheWrite(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("the immutable flag is how a removal fails on this platform")
	}
	ops, store := newOps(t)
	r, err := ops.CreateRecipe(spec())
	if err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(store.RecipeDir(), r.ID+"-"+Slugify(r.Title)+".json")
	if err := exec.Command("chflags", "uchg", stale).Run(); err != nil {
		t.Skipf("chflags: %v", err)
	}
	t.Cleanup(func() { _ = exec.Command("chflags", "nouchg", stale).Run() })

	before, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	title := "Renamed"
	got, err := ops.UpdateRecipe(r.ID, Fragment{Title: &title})
	if err != nil {
		t.Fatalf("the write landed, so the operation should not fail: %v", err)
	}
	if got.Title != title {
		t.Fatalf("the recipe's title = %q, want %q", got.Title, title)
	}

	// The journal grew: the artifact changed.
	after, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Events) != len(before.Events)+1 {
		t.Fatalf("the journal grew by %d events, want 1", len(after.Events)-len(before.Events))
	}

	// And the reader still picks the current copy, even though the leftover is there.
	read, err := store.LoadRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if read.Title != title {
		t.Fatalf("the reader picked %q, want %q", read.Title, title)
	}
}

// TestTheRemovalLeavesAnotherKindsFileAlone covers the kind of the id.
//
// The ids are numbered per kind, so "C01" and "R01" share a number — and a file of
// another kind sitting in the directory is not a copy of this artifact. Deleting it
// would lose an artifact the tool was never asked to touch, with nothing in the
// journal to say so.
//
// The record inside the file is a recipe's: a plan holds one record per file, and a
// record whose id is not a recipe's id is a plan this tool does not support. The
// *name* is the other kind's, and the name is what the removal reads.
func TestTheRemovalLeavesAnotherKindsFileAlone(t *testing.T) {
	ops, store := newOps(t)
	r, err := ops.CreateRecipe(spec())
	if err != nil {
		t.Fatal(err)
	}
	rogue := filepath.Join(store.RecipeDir(), "C01-rogue.json")
	if err := os.WriteFile(rogue, []byte("{\"id\":\"RCP02\",\"title\":\"a course\"}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	title := "A different title"
	if _, err := ops.UpdateRecipe(r.ID, Fragment{Title: &title}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(rogue); err != nil {
		t.Fatalf("the file of another kind was deleted: %v", err)
	}
}
