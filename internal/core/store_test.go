package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func zeroTime() time.Time { return time.Time{} }

func newStore(t *testing.T) *Store {
	t.Helper()
	root := filepath.Join(t.TempDir(), PlanDirName)
	s := Open(root)
	if err := s.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	return s
}

func TestInitCreatesTheLayout(t *testing.T) {
	s := newStore(t)
	for _, d := range []string{DirRecipes, DirCourses, DirBacklog} {
		info, err := os.Stat(filepath.Join(s.Root, d))
		if err != nil {
			t.Fatalf("stat %s: %v", d, err)
		}
		if !info.IsDir() {
			t.Fatalf("%s is not a directory", d)
		}
	}
	if _, err := os.Stat(filepath.Join(s.Root, FileProject)); err != nil {
		t.Fatalf("project file missing: %v", err)
	}
}

func TestInitIsIdempotent(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatalf("second init: %v", err)
	}
}

func TestNextRecipeIDScanRatherThanCount(t *testing.T) {
	s := newStore(t)

	// Create RCP01 and RCP03 by hand, leaving a gap. The next id must be RCP04 —
	// not RCP03 (a count) and not RCP02 (a len).
	//
	// The number is what must not be handed out twice, so a plan whose files skip a
	// number must not get that number handed out again.
	for _, id := range []string{"RCP01", "RCP03"} {
		r := &Recipe{ID: id, Title: "t " + id, Slug: "t-" + id, Lifecycle: LifecycleDraft}
		if err := s.SaveRecipe(r); err != nil {
			t.Fatal(err)
		}
	}
	got, err := s.NextRecipeID()
	if err != nil {
		t.Fatal(err)
	}
	if got != "RCP04" {
		t.Fatalf("next id = %q, want RCP04", got)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	s := newStore(t)
	in := &Recipe{
		ID:        "RCP01",
		Title:     "Ship the board",
		Lifecycle: LifecycleDraft,
		Intent:    Intent{Outcome: "o", Problem: "p"},
		Scope:     Scope{OutOfScope: []string{"a", "b"}},
		Acceptance: []Criterion{{
			ID: "AC-01", Description: "d", Verification: "go test",
		}},
	}
	if err := s.SaveRecipe(in); err != nil {
		t.Fatal(err)
	}
	out, err := s.LoadRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	if out.Title != in.Title || out.Intent.Outcome != "o" {
		t.Fatalf("round trip lost data: %+v", out)
	}
	if len(out.Acceptance) != 1 || out.Acceptance[0].Verification != "go test" {
		t.Fatalf("acceptance lost: %+v", out.Acceptance)
	}
}

// TestLoadRecipeTolerantOfIDForm covers the lookup.
//
// The **record** is in the current format — `RCP07` — and the **lookup** is
// tolerant: a name resolves to the record whatever form it is written in, because a
// name is not a record. The journal's old ids resolve the same way, and that is the
// history, not a migration.
func TestLoadRecipeTolerantOfIDForm(t *testing.T) {
	s := newStore(t)
	r := &Recipe{ID: "RCP07", Title: "t", Slug: "t", Lifecycle: LifecycleDraft}
	if err := s.SaveRecipe(r); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"RCP07", "RCP7", "rcp07", "7"} {
		if _, err := s.LoadRecipe(id); err != nil {
			t.Errorf("load %q: %v", id, err)
		}
	}
}

func TestLoadMissingRecipe(t *testing.T) {
	s := newStore(t)
	if _, err := s.LoadRecipe("R99"); err == nil {
		t.Fatal("expected an error for a missing recipe")
	}
}

func TestListRecipesSortsNumerically(t *testing.T) {
	s := newStore(t)
	for _, id := range []string{"RCP10", "RCP2", "RCP1"} {
		if err := s.SaveRecipe(&Recipe{ID: id, Title: "t" + id, Slug: "t" + id}); err != nil {
			t.Fatal(err)
		}
	}
	got, err := s.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"RCP1", "RCP2", "RCP10"}
	if len(got) != len(want) {
		t.Fatalf("got %d recipes, want %d", len(got), len(want))
	}
	for i, id := range want {
		if got[i].ID != id {
			t.Fatalf("order[%d] = %q, want %q", i, got[i].ID, id)
		}
	}
}

func TestWriteIsAtomic(t *testing.T) {
	s := newStore(t)
	r := &Recipe{ID: "RCP01", Title: "t", Slug: "t"}
	if err := s.SaveRecipe(r); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(s.RecipeDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() != "RCP01-t.json" {
			t.Fatalf("unexpected file left behind: %s", e.Name())
		}
	}
}

func TestFindRootWalksUp(t *testing.T) {
	s := newStore(t)
	nested := filepath.Join(s.Root, "recipes")
	got, err := FindRoot(nested)
	if err != nil {
		t.Fatal(err)
	}
	if got != s.Root {
		t.Fatalf("root = %q, want %q", got, s.Root)
	}
}
