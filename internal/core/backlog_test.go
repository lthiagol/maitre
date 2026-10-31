package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTheBacklogIsOneArtifact covers the shape of the plan.
//
// The backlog was two file types in two directories — a special and a prep note. It
// is one artifact now, in one directory, with a kind.
func TestTheBacklogIsOneArtifact(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, []string{"model"})
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != "BL01" {
		t.Fatalf("the id is %q, want BL01", item.ID)
	}
	if item.Kind != KindIdea {
		t.Fatalf("the kind is %q, want %q", item.Kind, KindIdea)
	}
	if item.Status != StatusOpen {
		t.Fatalf("the status is %q, want %q", item.Status, StatusOpen)
	}

	// It lives in the one directory.
	files := filesIn(t, store.BacklogDir())
	if len(files) != 1 {
		t.Fatalf("the backlog directory holds %v, want one file", files)
	}
	if files[0][:len("BL01")] != "BL01" {
		t.Fatalf("the file is %q, want it named after BL01", files[0])
	}
}

// TestInitDoesNotCreateTheOldBacklogDirectories covers the plan format.
//
// The old directories are gone: a plan written by this tool does not have them, and a
// tool that reads a plan does not look in them.
func TestInitDoesNotCreateTheOldBacklogDirectories(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{DirBacklog, DirRecipes, DirCourses} {
		if _, err := os.Stat(filepath.Join(s.Root, dir)); err != nil {
			t.Errorf("the plan should have a %s/ directory: %v", dir, err)
		}
	}
	for _, dir := range []string{"specials", "prep-notes"} {
		if _, err := os.Stat(filepath.Join(s.Root, dir)); err == nil {
			t.Errorf("the plan should not have a %s/ directory any more", dir)
		}
	}
}

// TestTheBacklogIdIsKindNeutral covers the id.
//
// An idea that is shaped into a small task is the same item: it keeps its id, so the
// trail survives the shaping. A kind in the id would break the trail every time an
// idea changed its mind.
func TestTheBacklogIdIsKindNeutral(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	shaped, err := ops.UpdateBacklogItem(item.ID, "", "", []string{"the work"}, KindFeature, nil)
	if err != nil {
		t.Fatal(err)
	}
	if shaped.ID != item.ID {
		t.Fatalf("the id changed: %q became %q", item.ID, shaped.ID)
	}
	if shaped.Kind != KindFeature {
		t.Fatalf("the kind is %q, want %q", shaped.Kind, KindFeature)
	}
	if len(shaped.Scope) != 1 {
		t.Fatalf("the scope did not carry: %+v", shaped.Scope)
	}

	// And the item on disk still answers to its id.
	got, err := store.LoadBacklogItem(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != KindFeature {
		t.Fatalf("the kind on disk is %q, want %q", got.Kind, KindFeature)
	}
}

// TestEveryFormOfTheBacklogIdNamesTheSameItem covers the reader.
//
// The writer emits `BL01`; a plan written before the rename carries `B01`; a person
// types what they see. All of them name the same item.
func TestEveryFormOfTheBacklogIdNamesTheSameItem(t *testing.T) {
	ops, store := newOps(t)

	got, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"BL01", "bl01", "BL1", "bl1", "B01", "b01", "1"} {
		item, err := store.LoadBacklogItem(id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if item.ID != got.ID {
			t.Errorf("%s resolved to %q, want %q", id, item.ID, got.ID)
		}
	}
}

// TestAFileNamedWithTheShortBacklogIdIsStillRead covers the plan on disk.
//
// The prefix is two letters and the plans written before the rename carried one. A
// file named `B01-<slug>.json` is the same artifact as one named `BL01-<slug>.json`,
// and the reader takes both.
func TestAFileNamedWithTheShortBacklogIdIsStillRead(t *testing.T) {
	ops, store := newOps(t)

	got, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(store.BacklogDir(), "B01-a-hunch.json")
	if err := os.Rename(filepath.Join(store.BacklogDir(), "BL01-a-hunch.json"), old); err != nil {
		t.Fatal(err)
	}
	item, err := store.LoadBacklogItem(got.ID)
	if err != nil {
		t.Fatalf("the item under its short file name: %v", err)
	}
	if item.ID != "BL01" {
		t.Errorf("the id is %q, want BL01", item.ID)
	}
	// And it answers to the short form.
	item, err = store.LoadBacklogItem("B01")
	if err != nil {
		t.Fatalf("B01: %v", err)
	}
	if item.ID != "BL01" {
		t.Errorf("B01 resolved to %q, want BL01", item.ID)
	}
}

// TestTheNextBacklogIdIsNotAReusedNumber covers the allocation.
//
// The ids are allocated by scanning what exists, never by counting — v1's `len()+1`
// handed out the same id twice after a removal.
func TestTheNextBacklogIdIsNotAReusedNumber(t *testing.T) {
	ops, store := newOps(t)

	for i := 0; i < 3; i++ {
		if _, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil); err != nil {
			t.Fatal(err)
		}
	}
	got, err := store.NextBacklogItemID()
	if err != nil {
		t.Fatal(err)
	}
	if got != "BL04" {
		t.Fatalf("the next id is %q, want BL04", got)
	}
}

// TestTheBacklogGroupsAreDerivedFromTheKind covers the two groups.
//
// The group is never stored: the kind is the only truth, and the console's groups
// are read off it.
func TestTheBacklogGroupsAreDerivedFromTheKind(t *testing.T) {
	ops, _ := newOps(t)

	if _, err := ops.CreateBacklogItem("An idea", "worth noticing", nil, KindIdea, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.CreateBacklogItem("A feature", "the work", []string{"a"}, KindFeature, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.CreateBacklogItem("A bug", "the work", []string{"b"}, KindBug, nil); err != nil {
		t.Fatal(err)
	}
	items, err := ops.ListBacklogItems()
	if err != nil {
		t.Fatal(err)
	}
	ideas, tasks := BacklogGroups(items)
	if len(ideas) != 1 || ideas[0].Title != "An idea" {
		t.Fatalf("the ideas are %+v, want one", ideas)
	}
	if len(tasks) != 2 {
		t.Fatalf("the small tasks are %+v, want two", tasks)
	}
}

// TestAnUnknownKindIsRefused covers the closed set.
//
// The kind decides the group and the shape a promotion gives the recipe, so a kind
// nobody knows is a typo — and a typo that silently disabled a feature is v1's
// lesson.
func TestAnUnknownKindIsRefused(t *testing.T) {
	ops, _ := newOps(t)

	if _, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, "hunch", nil); err == nil {
		t.Fatal("an unknown kind must be refused")
	}
	item, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.UpdateBacklogItem(item.ID, "", "", nil, "hunch", nil); err == nil {
		t.Fatal("an unknown kind must be refused on update too")
	}
}

// TestAnEmptyBacklogListIsWritten covers the file.
//
// The evidence is the raw file, not the loaded struct. A created item with no
// scope and no tags carries both keys, each `[]`. An update that clears the lists
// does too. A nil list saved directly is the same record. A file that lacks the
// keys still loads, the lists read as empty, and the next save writes the keys.
func TestAnEmptyBacklogListIsWritten(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertEmptyLists(t, backlogRaw(t, store, item.ID))

	if _, err := ops.UpdateBacklogItem(item.ID, "", "", []string{"the work"}, "", []string{"one"}); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.UpdateBacklogItem(item.ID, "", "", []string{}, "", []string{}); err != nil {
		t.Fatal(err)
	}
	assertEmptyLists(t, backlogRaw(t, store, item.ID))

	// A caller that passes nil. The writer is the place that cannot forget.
	bare := &BacklogItem{
		ID: "BL02", Title: "Bare", Kind: KindIdea, Status: StatusOpen,
		Created: item.Created, Updated: item.Updated,
	}
	if err := store.SaveBacklogItem(bare); err != nil {
		t.Fatal(err)
	}
	assertEmptyLists(t, backlogRaw(t, store, bare.ID))

	// A file that lacks the keys still loads. The next save writes them.
	missing := filepath.Join(store.BacklogDir(), "BL03-a-gap.json")
	body := []byte(`{"id":"BL03","title":"A gap","kind":"idea","status":"open","created":"2026-10-31T00:00:00Z","updated":"2026-10-31T00:00:00Z"}`)
	if err := os.WriteFile(missing, body, 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := store.LoadBacklogItem("BL03")
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Scope) != 0 || len(loaded.Tags) != 0 {
		t.Fatalf("the lists read as %v and %v, want empty", loaded.Scope, loaded.Tags)
	}
	if err := store.SaveBacklogItem(loaded); err != nil {
		t.Fatal(err)
	}
	assertEmptyLists(t, backlogRaw(t, store, "BL03"))
}

// assertEmptyLists checks the raw document carries both lists, each empty.
func assertEmptyLists(t *testing.T, raw []byte) {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"scope", "tags"} {
		list, ok := doc[key].([]any)
		if !ok {
			t.Fatalf("%s is %v, want []", key, doc[key])
		}
		if len(list) != 0 {
			t.Fatalf("%s holds %v, want none", key, list)
		}
	}
}

// backlogRaw is the item's file, as it was written.
func backlogRaw(t *testing.T, store *Store, id string) []byte {
	t.Helper()
	entries, err := os.ReadDir(store.BacklogDir())
	if err != nil {
		t.Fatal(err)
	}
	prefix := id + "-"
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), prefix) {
			raw, err := os.ReadFile(filepath.Join(store.BacklogDir(), e.Name()))
			if err != nil {
				t.Fatal(err)
			}
			return raw
		}
	}
	t.Fatalf("no file for %s", id)
	return nil
}

// TestTheBacklogUpdateLeavesTheListsAlone covers the list rule.
//
// nil means "not given" and leaves the field alone; a list means "replace the list
// it names". **Both** list fields obey the rule — the scope and the tags. The tags
// used to be wiped by an update that never mentioned them, because the two fields
// had two different guards.
func TestTheBacklogUpdateLeavesTheListsAlone(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("An item", "why", []string{"the work"}, KindIdea, []string{"one", "two"})
	if err != nil {
		t.Fatal(err)
	}

	// An update that names neither list leaves both alone.
	got, err := ops.UpdateBacklogItem(item.ID, "A better item", "", nil, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Tags) != 2 || len(got.Scope) != 1 {
		t.Fatalf("the update dropped a list: %+v", got)
	}

	// A list that is given replaces the list it names — and only that one.
	got, err = ops.UpdateBacklogItem(item.ID, "", "", nil, "", []string{"three"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Tags) != 1 || got.Tags[0] != "three" {
		t.Fatalf("the tags are %v, want the one the update named", got.Tags)
	}
	if len(got.Scope) != 1 {
		t.Fatalf("the tags wiped the scope: %+v", got)
	}

	// And the list on disk is the same.
	onDisk, err := store.LoadBacklogItem(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(onDisk.Tags) != 1 || len(onDisk.Scope) != 1 {
		t.Fatalf("the item on disk is %+v", onDisk)
	}

	// An **empty** list that was given clears the field: it is not "not given".
	got, err = ops.UpdateBacklogItem(item.ID, "", "", []string{}, "", []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Tags) != 0 || len(got.Scope) != 0 {
		t.Fatalf("the empty lists did not clear the fields: %+v", got)
	}
}

// TestTheBacklogItemLifecycle covers the ordinary path: create, edit, dismiss.
func TestTheBacklogItemLifecycle(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A hunch", "something worth noticing", nil, KindIdea, []string{"model"})
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != "BL01" || item.Status != StatusOpen {
		t.Fatalf("the item is %+v", item)
	}

	// Editing works while it is open.
	if _, err := ops.UpdateBacklogItem(item.ID, "A better hunch", "sharper", nil, "", nil); err != nil {
		t.Fatal(err)
	}
	got, err := store.LoadBacklogItem("BL01")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "A better hunch" || got.Body != "sharper" {
		t.Fatalf("the item is %+v", got)
	}
	// An empty field is left alone: an update says what changes.
	if _, err := ops.UpdateBacklogItem(item.ID, "", "", nil, "", nil); err != nil {
		t.Fatal(err)
	}
	got, err = store.LoadBacklogItem("BL01")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "A better hunch" || got.Body != "sharper" {
		t.Fatalf("an empty update changed the item: %+v", got)
	}

	// A dismissal works on every kind.
	if _, err := ops.DismissBacklogItem(item.ID, "not worth pursuing"); err != nil {
		t.Fatal(err)
	}
	got, err = store.LoadBacklogItem("BL01")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusDismissed {
		t.Fatalf("the status is %q, want %q", got.Status, StatusDismissed)
	}
	if _, err := ops.DismissBacklogItem(item.ID, "again"); err == nil {
		t.Fatal("expected dismissing a dismissed item to fail")
	}
}

// TestTheBacklogItemsEventsAreJournalledMeaningfully covers the journal's kinds.
//
// The kinds are a closed enum and they are meant to mean something: an edit was
// journalled as a `comment` and a dismissal as a `promotion`, so a reader of the log
// saw a promotion for something the chef threw away.
func TestTheBacklogItemsEventsAreJournalledMeaningfully(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A hunch", "something worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.UpdateBacklogItem(item.ID, "A better hunch", "sharper", nil, "", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.DismissBacklogItem(item.ID, "not worth pursuing"); err != nil {
		t.Fatal(err)
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) != 3 {
		t.Fatalf("the journal has %d events, want 3", len(j.Events))
	}
	if j.Events[0].Kind != KindCreated {
		t.Fatalf("kind = %q, want %q", j.Events[0].Kind, KindCreated)
	}
	if j.Events[1].Kind != KindDecision {
		t.Fatalf("the edit's kind = %q, want %q", j.Events[1].Kind, KindDecision)
	}
	if j.Events[2].Kind != KindOverlay {
		t.Fatalf("the dismissal's kind = %q, want %q", j.Events[2].Kind, KindOverlay)
	}
}

// TestTheBacklogItemsAreListed covers the listing.
//
// The list is never nil, and one artifact is one entry even when a rename left a
// copy behind.
func TestTheBacklogItemsAreListed(t *testing.T) {
	ops, store := newOps(t)
	for i := 0; i < 3; i++ {
		if _, err := ops.CreateBacklogItem("A hunch", "body", nil, KindIdea, nil); err != nil {
			t.Fatal(err)
		}
	}
	items, err := ops.ListBacklogItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Fatalf("got %d items, want 3", len(items))
	}
	if items == nil {
		t.Fatal("the list is nil")
	}

	// The journal records the creations.
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	kinds := map[string]int{}
	for _, e := range j.Events {
		kinds[e.Kind]++
	}
	if kinds[KindCreated] < 3 {
		t.Fatalf("expected the creations to be journalled: %v", kinds)
	}
}
