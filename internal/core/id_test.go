package core

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTheWriterEmitsRCPIds covers the new ids.
//
// A recipe's id is `RCP01`, and the file it lives in is named after it.
func TestTheWriterEmitsRCPIds(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ops := NewOperations(s, AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})

	got, err := ops.CreateRecipe(RecipeSpec{
		Title:      "A recipe",
		Outcome:    "o",
		Problem:    "p",
		InScope:    []string{"a"},
		OutOfScope: []string{"b"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "RCP01" {
		t.Fatalf("the id is %q, want RCP01", got.ID)
	}

	entries, err := os.ReadDir(s.RecipeDir())
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	if len(names) != 1 {
		t.Fatalf("the recipe directory holds %v, want one file", names)
	}
	if names[0][:len("RCP01")] != "RCP01" {
		t.Fatalf("the file is %q, want it named after RCP01", names[0])
	}
}

// TestEveryFormOfTheIdNamesTheSameRecipe covers the reader.
//
// The writer emits `RCP01`; the plans written before the rename carry `R01`; a
// person types what they see. All of them name the same recipe.
//
// **Every** leading letter is stripped, not just the first: a normaliser that stopped
// at one letter would leave `CP01`, which matches neither the new id nor the old one.
func TestEveryFormOfTheIdNamesTheSameRecipe(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ops := NewOperations(s, AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})

	got, err := ops.CreateRecipe(RecipeSpec{
		Title:      "A recipe",
		Outcome:    "o",
		Problem:    "p",
		InScope:    []string{"a"},
		OutOfScope: []string{"b"},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, id := range []string{"RCP01", "rcp01", "RCP1", "R01", "r1", "1", "RCP001"} {
		rec, err := ops.GetRecipe(id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if rec.ID != got.ID {
			t.Errorf("%s resolved to %q, want %q", id, rec.ID, got.ID)
		}
	}
}

// TestTheNormaliserStripsEveryLeadingLetter covers the bug the AC exists to prevent.
//
// The old normaliser removed one letter. `RCP01` would have become `CP01`, which
// matches neither `RCP01` nor `R01` — so a recipe written after the rename could not
// be read at all.
//
// The backlog's prefix is **two** letters and the recipe's is **three**, and the
// normaliser has to answer to both: a stripper that removed a fixed number of letters
// would leave `L01` or `CP01`.
func TestTheNormaliserStripsEveryLeadingLetter(t *testing.T) {
	cases := map[string]string{
		"RCP01": "1",
		"rcp01": "1",
		"R01":   "1",
		"r1":    "1",
		"1":     "1",
		"10":    "10",
		"RCP10": "10",
		"BL01":  "1",
		"bl01":  "1",
		"BL1":   "1",
		"B01":   "1",
		"b01":   "1",
		"BL10":  "10",
		"N01":   "1",
		"S01":   "1",
		"":      "0",
		"R":     "0",
	}
	for id, want := range cases {
		if got := normaliseID(id); got != want {
			t.Errorf("normaliseID(%q) = %q, want %q", id, got, want)
		}
	}
}

// TestAFileNamedWithTheOldIdIsRefused covers the plan on disk.
//
// The ids changed; the files written before the change did not. A file named
// `R01-<slug>.json` is a **retired id form on disk** — and the file's name is part
// of the format, so the plan is not supported: the reader says so plainly instead
// of reading a plan that says one thing in its name and another in its content.
//
// The **lookup** stays tolerant: a name is not a record, and `R01` still names the
// recipe the plan holds.
func TestAFileNamedWithTheOldIdIsRefused(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ops := NewOperations(s, AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})

	got, err := ops.CreateRecipe(RecipeSpec{
		Title:      "A recipe",
		Outcome:    "o",
		Problem:    "p",
		InScope:    []string{"a"},
		OutOfScope: []string{"b"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Rename the file to the id the plans written before the rename carry, and
	// change the id inside it to match.
	old := filepath.Join(s.RecipeDir(), "R01-a-recipe.json")
	if err := os.Rename(filepath.Join(s.RecipeDir(), "RCP01-a-recipe.json"), old); err != nil {
		t.Fatal(err)
	}
	_, err = s.LoadRecipe(got.ID)
	if err == nil {
		t.Fatal("a file named with a retired id form must not be read")
	}
	var unsupported *UnsupportedPlanError
	if !errors.As(err, &unsupported) {
		t.Fatalf("the refusal should be an unsupported plan: %v", err)
	}
	if !strings.Contains(err.Error(), "RCP01") {
		t.Errorf("the refusal should name the current form: %v", err)
	}

	// And the lookup's tolerance stands: the old form still names the recipe the
	// plan holds — as long as the file is named in the current form.
	if err := os.Rename(old, filepath.Join(s.RecipeDir(), "RCP01-a-recipe.json")); err != nil {
		t.Fatal(err)
	}
	rec, err := ops.GetRecipe("R01")
	if err != nil {
		t.Fatalf("R01: %v", err)
	}
	if rec.ID != "RCP01" {
		t.Errorf("R01 resolved to %q, want RCP01", rec.ID)
	}
}

// TestTheJournalsOldIdStillResolves covers the append-only history.
//
// The journal is the plan's history: it is not rewritten when the ids change. An
// entry that says `R01` still names the recipe that is now `RCP01`, and the reader
// must resolve it there.
func TestTheJournalsOldIdStillResolves(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ops := NewOperations(s, AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})

	got, err := ops.CreateRecipe(RecipeSpec{
		Title:      "A recipe",
		Outcome:    "o",
		Problem:    "p",
		InScope:    []string{"a"},
		OutOfScope: []string{"b"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// The journal keeps the old id — the history is not rewritten.
	j, err := s.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	for i := range j.Events {
		if j.Events[i].Recipe == got.ID {
			j.Events[i].Recipe = "R01"
		}
	}
	raw, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(s.Root, FileActivity), raw, 0o644); err != nil {
		t.Fatal(err)
	}

	// Every entry that names the recipe resolves, whatever the form of the id.
	loaded, err := s.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range loaded.Events {
		if e.Recipe == "" {
			continue
		}
		found = true
		rec, err := ops.GetRecipe(e.Recipe)
		if err != nil {
			t.Fatalf("the journal's %q does not resolve: %v", e.Recipe, err)
		}
		if rec.ID != got.ID {
			t.Errorf("the journal's %q resolved to %q, want %q", e.Recipe, rec.ID, got.ID)
		}
	}
	if !found {
		t.Fatal("the journal holds no entries that name a recipe")
	}
}

// TestTheNextIdIsNotAReusedNumber covers the allocation after the rename.
//
// The prefix changed, so a prefix check would see an empty directory and hand out
// `RCP01` again — over a recipe the plan already has.
func TestTheNextIdIsNotAReusedNumber(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ops := NewOperations(s, AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})

	for i := 0; i < 3; i++ {
		if _, err := ops.CreateRecipe(RecipeSpec{
			Title:      "A recipe",
			Outcome:    "o",
			Problem:    "p",
			InScope:    []string{"a"},
			OutOfScope: []string{"b"},
		}); err != nil {
			t.Fatal(err)
		}
	}
	got, err := s.NextRecipeID()
	if err != nil {
		t.Fatal(err)
	}
	if got != "RCP04" {
		t.Fatalf("the next id is %q, want RCP04", got)
	}
}

// TestTheReposOwnJournalStillResolves covers the migration.
//
// The repo's own plan was renamed: the recipe file is `RCP01-<slug>.json` and the id
// inside it is `RCP01`. The journal is the plan's history and was **not** rewritten —
// it still says `R01`.
//
// A reader that cannot resolve the journal's ids cannot read the plan's own past.
func TestTheReposOwnJournalStillResolves(t *testing.T) {
	root := filepath.Join("..", "..", PlanDirName)
	if _, err := os.Stat(root); err != nil {
		t.Skipf("the repo has no plan at %s", root)
	}
	s := Open(root)
	ops := NewOperations(s, AuditContext{Actor: "runner", Session: "s-test", Role: "runner"})

	j, err := s.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) == 0 {
		t.Fatal("the repo's own journal is empty")
	}

	legacy := false
	for _, e := range j.Events {
		if e.Recipe == "" {
			continue
		}
		rec, err := ops.GetRecipe(e.Recipe)
		if err != nil {
			t.Fatalf("the journal's %q does not resolve: %v", e.Recipe, err)
		}
		if e.Recipe == "R01" {
			legacy = true
			if rec.ID != "RCP01" {
				t.Errorf("the journal's R01 resolved to %q, want RCP01", rec.ID)
			}
		}
	}
	if !legacy {
		t.Error("the journal should still carry the old R01 — it is history, not a file to rewrite")
	}
}
