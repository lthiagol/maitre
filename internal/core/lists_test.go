package core

import (
	"testing"
)

// TestTheListFunctionsReturnAList covers the core's list contract.
//
// A list is a list: an empty one is `[]`, never `null`. The plan is read by agents,
// and the standard readers — `jq`, a JSON parser in any language — break on `null`.
//
// The guarantee belongs to the core, not to the CLI that happens to print it, so it
// is asserted here.
func TestTheListFunctionsReturnAList(t *testing.T) {
	ops, store := newOps(t)

	// An empty plan.
	recipes, err := store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if recipes == nil {
		t.Fatal("ListRecipes returned nil")
	}
	specials, err := store.ListSpecials()
	if err != nil {
		t.Fatal(err)
	}
	if specials == nil {
		t.Fatal("ListSpecials returned nil")
	}
	notes, err := store.ListPrepNotes()
	if err != nil {
		t.Fatal(err)
	}
	if notes == nil {
		t.Fatal("ListPrepNotes returned nil")
	}
	ready, waiting := Path(nil)
	if ready == nil || waiting == nil {
		t.Fatalf("Path returned %v / %v, want two lists", ready, waiting)
	}
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if j.Since(0) == nil {
		t.Fatal("Since returned nil")
	}
	if j.Recent(10) == nil {
		t.Fatal("Recent returned nil")
	}

	// And a plan with something in it.
	if _, err := ops.CreateRecipe(readySpec()); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.CreateSpecial("A hunch", "worth noticing", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.CreatePrepNote("A note", "first", nil); err != nil {
		t.Fatal(err)
	}
	recipes, err = store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if recipes == nil {
		t.Fatal("ListRecipes returned nil")
	}
	specials, err = store.ListSpecials()
	if err != nil {
		t.Fatal(err)
	}
	if specials == nil {
		t.Fatal("ListSpecials returned nil")
	}
	notes, err = store.ListPrepNotes()
	if err != nil {
		t.Fatal(err)
	}
	if notes == nil {
		t.Fatal("ListPrepNotes returned nil")
	}
	j, err = store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if j.Since(0) == nil {
		t.Fatal("Since returned nil")
	}
	if j.Recent(10) == nil {
		t.Fatal("Recent returned nil")
	}
}
