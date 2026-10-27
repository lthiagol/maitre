package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestThePreMigrationBacklogIdStillResolves covers the id the plan had before the
// move.
//
// The backlog's ids were `S01` (a special) and `N01` (a prep note); they are `BL01`
// now. The journal is **append-only**, so a plan written before the move still
// carries the old id in its history — and a reader that cannot resolve it cannot
// read the plan's own past.
func TestThePreMigrationBacklogIdStillResolves(t *testing.T) {
	ops, store := newOps(t)

	item, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != "BL01" {
		t.Fatalf("the id is %q, want BL01", item.ID)
	}

	// The journal keeps the old id — the history is not rewritten.
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	for i := range j.Events {
		var data map[string]any
		if len(j.Events[i].Data) == 0 {
			continue
		}
		if err := json.Unmarshal(j.Events[i].Data, &data); err != nil {
			t.Fatal(err)
		}
		if data["kind"] == "backlog" {
			data["id"] = "S01"
			raw, err := json.Marshal(data)
			if err != nil {
				t.Fatal(err)
			}
			j.Events[i].Data = raw
		}
	}
	raw, err := json.Marshal(j)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(store.Root, FileActivity), raw, 0o644); err != nil {
		t.Fatal(err)
	}

	// The old id resolves to the item it names.
	got, err := store.LoadBacklogItem("S01")
	if err != nil {
		t.Fatalf("the pre-migration id does not resolve: %v", err)
	}
	if got.ID != "BL01" {
		t.Fatalf("S01 resolved to %q, want BL01", got.ID)
	}
}

// TestTheRepoPlansJournalIdsStillResolve covers the repo's own plan.
//
// The plan in this repository was migrated: its recipe file was renamed and the id
// inside it changed, and its backlog's shape moved to one directory. The journal was
// not rewritten — it is append-only — so it still carries the ids the plan had
// before the move. Every one of them must still resolve.
func TestTheRepoPlansJournalIdsStillResolve(t *testing.T) {
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

	named := 0
	for _, e := range j.Events {
		if e.Recipe != "" {
			named++
			rec, err := ops.GetRecipe(e.Recipe)
			if err != nil {
				t.Fatalf("the journal's %q does not resolve: %v", e.Recipe, err)
			}
			if rec.ID == "" {
				t.Fatalf("the journal's %q resolved to nothing", e.Recipe)
			}
		}
		// A backlog id the journal names: it resolves to an item.
		var data map[string]any
		if len(e.Data) == 0 {
			continue
		}
		if err := json.Unmarshal(e.Data, &data); err != nil {
			t.Fatalf("an event's data is not a document: %v\n%s", err, e.Data)
		}
		if data["kind"] != "backlog" {
			continue
		}
		id, _ := data["id"].(string)
		if id == "" {
			t.Fatal("a backlog event names no id")
		}
		named++
		item, err := s.LoadBacklogItem(id)
		if err != nil {
			t.Fatalf("the journal's %q does not resolve: %v", id, err)
		}
		if item.ID == "" {
			t.Fatalf("the journal's %q resolved to nothing", id)
		}
	}
	if named == 0 {
		t.Fatal("the journal names no ids — the check found nothing to hold")
	}
}
