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
	for _, d := range []string{DirRecipes, DirCourses, DirPrepNotes, DirSpecials} {
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

	// Create R01 and R03 by hand, leaving a gap. The next id must be R04 — not R03
	// (a count) and not R02 (a len).
	for _, id := range []string{"R01", "R03"} {
		r := &Recipe{ID: id, Title: "t " + id, Slug: "t-" + id, Lifecycle: LifecycleDraft}
		if err := s.SaveRecipe(r); err != nil {
			t.Fatal(err)
		}
	}
	got, err := s.NextRecipeID()
	if err != nil {
		t.Fatal(err)
	}
	if got != "R04" {
		t.Fatalf("next id = %q, want R04", got)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	s := newStore(t)
	in := &Recipe{
		ID:        "R01",
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
	out, err := s.LoadRecipe("R01")
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

func TestLoadRecipeTolerantOfIDForm(t *testing.T) {
	s := newStore(t)
	r := &Recipe{ID: "R07", Title: "t", Slug: "t", Lifecycle: LifecycleDraft}
	if err := s.SaveRecipe(r); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"R07", "R7", "r07", "7"} {
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
	for _, id := range []string{"R10", "R2", "R1"} {
		if err := s.SaveRecipe(&Recipe{ID: id, Title: "t" + id, Slug: "t" + id}); err != nil {
			t.Fatal(err)
		}
	}
	got, err := s.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"R1", "R2", "R10"}
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
	r := &Recipe{ID: "R01", Title: "t", Slug: "t"}
	if err := s.SaveRecipe(r); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(s.RecipeDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() != "R01-t.json" {
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
