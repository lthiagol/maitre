package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEveryWriterWritesTheSameShape covers the hole D-5 found.
//
// The reader normalises a recipe it reads, so a recipe that was read and written back
// says `[]`. A recipe written **first** time — the promotion's, and `recipe create`'s
// — never went through the reader, and it wrote `null`:
//
//	recipe create      →  "scope": {"in_scope": [], "out_of_scope": []}
//	backlog promote    →  "scope": {"in_scope": null, "out_of_scope": null}
//
// Same format, two answers — and one of them is `null`, which the standard readers
// break on. The rule belongs in the writer, so every writer agrees.
func TestEveryWriterWritesTheSameShape(t *testing.T) {
	ops, store := newOps(t)

	// The recipe that is created.
	if _, err := ops.CreateRecipe(RecipeSpec{Title: "A plain recipe"}); err != nil {
		t.Fatal(err)
	}
	// And the recipe that is promoted.
	item, err := ops.CreateBacklogItem("A small task", "worth doing", nil, KindFeature, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.PromoteBacklogItem(item.ID); err != nil {
		t.Fatal(err)
	}

	for _, want := range []struct {
		title      string
		acceptance int
	}{
		{"A plain recipe", 0},
		{"A small task", 1},
	} {
		name := want.title
		path := recipeFile(t, store, name)
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var doc map[string]any
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatal(err)
		}
		scope, ok := doc["scope"].(map[string]any)
		if !ok {
			t.Fatalf("%s: the scope is %v, want an object", name, doc["scope"])
		}
		for _, field := range []string{"in_scope", "out_of_scope"} {
			list, ok := scope[field].([]any)
			if !ok {
				t.Fatalf("%s: the scope's %s is %v, want a list", name, field, scope[field])
			}
			if len(list) != 0 {
				t.Fatalf("%s: the scope's %s holds %v, want none", name, field, list)
			}
		}
		acceptance, ok := doc["acceptance"].([]any)
		if !ok {
			t.Fatalf("%s: the acceptance is %v, want a list", name, doc["acceptance"])
		}
		if len(acceptance) != want.acceptance {
			t.Fatalf("%s: the acceptance holds %d criteria, want %d", name, len(acceptance), want.acceptance)
		}
	}
}

// recipeFile is the path of the recipe file that carries the title.
func recipeFile(t *testing.T, store *Store, title string) string {
	t.Helper()
	entries, err := os.ReadDir(store.RecipeDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), Slugify(title)) {
			return filepath.Join(store.RecipeDir(), e.Name())
		}
	}
	t.Fatalf("no recipe file holds %q", title)
	return ""
}

// TestTheApprovalDigestStillVerifies covers the digest the writer's change must not
// move.
//
// The approval's payload canonicalises its lists, so `null` and `[]` hash the same —
// and a digest written before the change still verifies after it. The recipe is
// approved, written back, and read from the disk: the digest on the file is the one
// the tool recomputes.
func TestTheApprovalDigestStillVerifies(t *testing.T) {
	ops, store := newOps(t)

	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventGroom); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, EventSpecify); err != nil {
		t.Fatal(err)
	}
	approved, err := ops.ApproveRecipe(r.ID, "chef")
	if err != nil {
		t.Fatal(err)
	}
	if approved.Approval == nil {
		t.Fatal("the recipe was not approved")
	}

	// The file on disk, read the way a reader reads it.
	stored, err := store.LoadRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	want, err := stored.Fingerprint()
	if err != nil {
		t.Fatal(err)
	}
	if want != stored.Approval.Fingerprint {
		t.Fatalf("the digest is %q, want %q — the writer's shape moved the digest",
			stored.Approval.Fingerprint, want)
	}
}
