package cli

import "testing"

// journalCount is how many events the plan's journal holds.
func journalCount(t *testing.T, dir string) int {
	t.Helper()
	res := cli(t, dir, "log", "--limit", "1000")
	if !res.ok(t) {
		t.Fatalf("log failed:\n%s", res.Stderr)
	}
	events, ok := res.doc(t)["events"].([]any)
	if !ok {
		t.Fatal("the journal document carries no events")
	}
	return len(events)
}

// TestPrepNoteUpdateGoesThroughAnOperation covers the write that used to bypass the
// domain layer.
//
// The CLI wrote the prep note file itself: no lock, no operation, and — worst — no
// journal entry. The journal is the only event log, so the edit was invisible to
// `scale log --since`, to the console's Log tab, and to any agent watching for
// changes.
func TestPrepNoteUpdateGoesThroughAnOperation(t *testing.T) {
	dir := newPlan(t)

	res := cli(t, dir, "prepnote", "create", "--title", "A note", "--intent", "first")
	if !res.ok(t) {
		t.Fatalf("prepnote create failed:\n%s", res.Stderr)
	}
	before := journalCount(t, dir)

	res = cli(t, dir, "prepnote", "update", "N01", "--intent", "second")
	if !res.ok(t) {
		t.Fatalf("prepnote update failed:\n%s", res.Stderr)
	}

	// The plan file changed.
	notes := cli(t, dir, "prepnote", "list")
	if !notes.ok(t) {
		t.Fatalf("prepnote list failed:\n%s", notes.Stderr)
	}
	list := notes.doc(t)["prepNotes"].([]any)
	if len(list) != 1 {
		t.Fatalf("got %d prep notes, want 1", len(list))
	}
	note := list[0].(map[string]any)
	if note["intent"] != "second" {
		t.Fatalf("intent = %v, want second", note["intent"])
	}

	// And the journal knows about it.
	after := journalCount(t, dir)
	if after != before+1 {
		t.Fatalf("the journal has %d events, want %d — the edit must be journalled", after, before+1)
	}
}

// TestPrepNoteUpdateRefusesAPromotedNote covers the one-way ladder.
//
// A promoted note became a recipe; editing it would write to an artifact that is
// already history.
func TestPrepNoteUpdateRefusesAPromotedNote(t *testing.T) {
	dir := newPlan(t)

	res := cli(t, dir, "prepnote", "create", "--title", "A note", "--intent", "first")
	if !res.ok(t) {
		t.Fatalf("prepnote create failed:\n%s", res.Stderr)
	}
	res = cli(t, dir, "prepnote", "promote", "N01")
	if !res.ok(t) {
		t.Fatalf("prepnote promote failed:\n%s", res.Stderr)
	}

	res = cli(t, dir, "prepnote", "update", "N01", "--intent", "second")
	if res.Code == 0 {
		t.Fatal("a promoted note must not be edited")
	}
}
