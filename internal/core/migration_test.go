package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRecordingAMigrationAppendsAControlEvent covers the operation.
//
// The `control` kind had no writer at all until now. This gives it its first: an
// event about the plan itself, carrying no recipe, recording the id form the plan
// came from and the one it goes to.
func TestRecordingAMigrationAppendsAControlEvent(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ops := NewOperations(s, AuditContext{Actor: "runner", Session: "s-migrate", Role: "runner"})

	if err := ops.RecordMigration("R01", "RCP01"); err != nil {
		t.Fatal(err)
	}

	j, err := s.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) != 1 {
		t.Fatalf("the journal holds %d events, want 1", len(j.Events))
	}
	e := j.Events[0]
	if e.Kind != KindControl {
		t.Errorf("the event kind is %q, want %q", e.Kind, KindControl)
	}
	if e.Recipe != "" {
		t.Errorf("the event names the recipe %q — a migration is about the plan", e.Recipe)
	}
	if e.DomainVersion != DomainVersion {
		t.Errorf("the event carries version %d, want %d", e.DomainVersion, DomainVersion)
	}
	var data map[string]any
	if err := json.Unmarshal(e.Data, &data); err != nil {
		t.Fatalf("the event's data is not a document: %v\n%s", err, e.Data)
	}
	if data["from"] != "R01" || data["to"] != "RCP01" {
		t.Errorf("the event records %v, want the old and the new id form", data)
	}
}

// TestTheMigrationStampsThePlan covers the stamp.
//
// The plan's version is the version on its newest journal entry. A plan whose
// entries say nothing about the version reads as version zero — and nothing fences
// a stale tool off from it. The migration is what moves the stamp.
func TestTheMigrationStampsThePlan(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ops := NewOperations(s, AuditContext{Actor: "runner", Session: "s-migrate", Role: "runner"})

	j, err := s.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if got := PlanDomainVersion(j); got != 0 {
		t.Fatalf("a plan with no events reads as version %d, want 0", got)
	}

	if err := ops.RecordMigration("R01", "RCP01"); err != nil {
		t.Fatal(err)
	}

	j, err = s.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if got := PlanDomainVersion(j); got != DomainVersion {
		t.Fatalf("the plan reads as version %d, want %d", got, DomainVersion)
	}
}

// TestTheMigrationDoesNotRewriteTheHistory covers the append-only rule.
//
// The journal is the plan's history and it is never rewritten. A plan whose ids
// were renamed cannot have its old events edited to the new ids: the migration is
// *added*, after them, and the history stays what it was.
func TestTheMigrationDoesNotRewriteTheHistory(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ops := NewOperations(s, AuditContext{Actor: "runner", Session: "s-migrate", Role: "runner"})

	// A plan with a history: a recipe, created and groomed.
	if _, err := ops.CreateRecipe(RecipeSpec{
		Title:      "A recipe",
		Outcome:    "o",
		Problem:    "p",
		InScope:    []string{"a"},
		OutOfScope: []string{"b"},
	}); err != nil {
		t.Fatal(err)
	}
	before, err := s.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(before.Events) == 0 {
		t.Fatal("the journal is empty")
	}
	history, err := json.Marshal(before.Events)
	if err != nil {
		t.Fatal(err)
	}

	if err := ops.RecordMigration("R01", "RCP01"); err != nil {
		t.Fatal(err)
	}

	after, err := s.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Events) != len(before.Events)+1 {
		t.Fatalf("the journal holds %d events, want %d — the history is not rewritten, the event is added",
			len(after.Events), len(before.Events)+1)
	}
	kept, err := json.Marshal(after.Events[:len(before.Events)])
	if err != nil {
		t.Fatal(err)
	}
	if string(kept) != string(history) {
		t.Errorf("the history changed:\nbefore %s\nafter  %s", history, kept)
	}
}

// TestAMigrationNeedsBothIdForms covers the guard.
//
// A migration with no id form records nothing about anything, and a migration that
// changes nothing is not a migration. Both are refused rather than written.
func TestAMigrationNeedsBothIdForms(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ops := NewOperations(s, AuditContext{Actor: "runner", Session: "s-migrate", Role: "runner"})

	for _, m := range [][2]string{
		{"", "RCP01"},
		{"R01", ""},
		{"R01", "R01"},
	} {
		if err := ops.RecordMigration(m[0], m[1]); err == nil {
			t.Errorf("the migration %v should be refused", m)
		}
	}

	j, err := s.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) != 0 {
		t.Fatalf("a refused migration wrote %d events", len(j.Events))
	}
}

// TestTheRepoPlanCarriesItsMigration covers the repo's own plan.
//
// The plan in this repository was migrated by hand: its recipe file was renamed and
// the id inside it changed. The journal was not rewritten — it is append-only — so
// the migration was *added* to it, and the plan now carries a version.
//
// The backlog's collapse was migrated the same way: the plan's shape went from two
// directories to one, and the journal says so.
func TestTheRepoPlanCarriesItsMigration(t *testing.T) {
	root := filepath.Join("..", "..", PlanDirName)
	if _, err := os.Stat(root); err != nil {
		t.Skipf("the repo has no plan at %s", root)
	}
	s := Open(root)
	j, err := s.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) == 0 {
		t.Fatal("the repo's own journal is empty")
	}

	last := j.Events[len(j.Events)-1]
	if last.Kind != KindControl {
		t.Errorf("the newest event is %q, want %q — the migration is the last thing that happened",
			last.Kind, KindControl)
	}
	if last.DomainVersion != DomainVersion {
		t.Errorf("the newest event carries version %d, want %d", last.DomainVersion, DomainVersion)
	}

	// The migrations the plan went through, in the order they happened.
	var migrations []string
	for _, e := range j.Events {
		if e.Kind != KindControl {
			continue
		}
		var data map[string]any
		if len(e.Data) == 0 {
			continue
		}
		if err := json.Unmarshal(e.Data, &data); err != nil {
			t.Fatalf("a control event's data is not a document: %v\n%s", err, e.Data)
		}
		if data["kind"] != "migration" {
			continue
		}
		migrations = append(migrations, fmt.Sprintf("%v -> %v", data["from"], data["to"]))
	}
	for _, want := range []string{"R01 -> RCP01", "specials+prep-notes -> backlog"} {
		found := false
		for _, m := range migrations {
			if m == want {
				found = true
			}
		}
		if !found {
			t.Errorf("the journal should carry the migration %q: %v", want, migrations)
		}
	}

	// The history is untouched: the events before the first migration carry no
	// version — they were written before the version guard existed.
	first := len(j.Events)
	for i, e := range j.Events {
		if e.Kind == KindControl {
			first = i
			break
		}
	}
	for _, e := range j.Events[:first] {
		if e.DomainVersion != 0 {
			t.Errorf("a historical event carries version %d — the history is not rewritten", e.DomainVersion)
		}
	}
	if !strings.Contains(string(mustJSON(t, j.Events)), "R01") {
		t.Error("the history should still carry the old R01 — it is history, not a file to rewrite")
	}
}

// mustJSON renders a value as JSON, or fails.
func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
