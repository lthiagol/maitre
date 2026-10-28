package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeRecipeFile puts a recipe file on disk, exactly as it is written.
func writeRecipeFile(t *testing.T, s *Store, name, body string) {
	t.Helper()
	dir := filepath.Join(s.Root, DirRecipes)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestTheCurrentFormatIsRead covers the reader's ordinary work.
//
// A plan in the current format reads. That is the whole point of the check: it
// refuses the retired shapes and nothing else.
func TestTheCurrentFormatIsRead(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatalf("a plan in the current format must read: %v", err)
	}
	if got.ID != r.ID {
		t.Fatalf("read %q, want %q", got.ID, r.ID)
	}
}

// TestTheRetiredBacklogDirectoriesAreRefused covers the directories.
//
// The backlog was two directories; it is one now. A plan that still carries the old
// ones is not supported, and the reader says so — it does not pretend the backlog is
// empty.
func TestTheRetiredBacklogDirectoriesAreRefused(t *testing.T) {
	for _, dir := range []string{"specials", "prep-notes"} {
		t.Run(dir, func(t *testing.T) {
			s := newStore(t)
			if err := s.Init(); err != nil {
				t.Fatal(err)
			}
			old := filepath.Join(s.Root, dir)
			if err := os.MkdirAll(old, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(old, "S01-a-hunch.json"),
				[]byte(`{"id":"S01","title":"A hunch","status":"open"}`), 0o644); err != nil {
				t.Fatal(err)
			}

			_, err := s.ListBacklogItems()
			if err == nil {
				t.Fatal("a plan that carries a retired directory must not be read")
			}
			var unsupported *UnsupportedPlanError
			if !errors.As(err, &unsupported) {
				t.Fatalf("the refusal should be an unsupported plan: %v", err)
			}
			if !strings.Contains(err.Error(), dir+"/") {
				t.Errorf("the refusal should name what the plan carries: %v", err)
			}
			if !strings.Contains(err.Error(), "not supported") {
				t.Errorf("the refusal should say the plan is not supported: %v", err)
			}
		})
	}
}

// TestTheRetiredFieldNamesAreRefused covers the records.
//
// The recipe's completion stamp was `evidence`; it is `completion` now. A
// criterion's record was `proofs`; it is `evidence` now. A record under the old
// name is a record the current reader does not know, and the reader says so.
func TestTheRetiredFieldNamesAreRefused(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string
	}{
		{
			name: "the old completion stamp",
			body: `{"id":"RCP01","title":"t","slug":"t","lifecycle":"closed",
				"evidence":{"completed_at":"2026-10-26T02:21:56Z"}}`,
			want: "`evidence` stamp",
		},
		{
			name: "the old criterion record",
			body: `{"id":"RCP01","title":"t","slug":"t","lifecycle":"executed",
				"acceptance":[{"id":"AC-01","description":"d","verification":"go test ./...",
				"proofs":[{"command":"go test ./...","exit_code":0,"matched":1,"match":"proven"}]}]}`,
			want: "`proofs`",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := newStore(t)
			writeRecipeFile(t, s, "RCP01-t.json", c.body)

			_, err := s.LoadRecipe("RCP01")
			if err == nil {
				t.Fatal("a record under a retired field name must not be read")
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("the refusal should name the retired field: %v", err)
			}

			// And the listing refuses it too: a plan is not half-read.
			if _, err := s.ListRecipes(); err == nil {
				t.Fatal("the listing must not read a plan that carries a retired field")
			}
		})
	}
}

// TestTheRetiredIDFormIsRefused covers the ids.
//
// A recipe's id was `R01`; it is `RCP01` now. A live record under the old form is
// a record of the retired format, and the reader refuses it.
func TestTheRetiredIDFormIsRefused(t *testing.T) {
	s := newStore(t)
	writeRecipeFile(t, s, "R01-t.json", `{"id":"R01","title":"t","slug":"t","lifecycle":"draft"}`)

	_, err := s.LoadRecipe("R01")
	if err == nil {
		t.Fatal("a record under a retired id form must not be read")
	}
	if !strings.Contains(err.Error(), "RCP01") {
		t.Errorf("the refusal should name the current form: %v", err)
	}
}

// TestTheRetiredBacklogIDFormIsRefused covers the backlog's ids.
//
// The backlog's ids were `S01` and `P01`; they are `BL01` now.
func TestTheRetiredBacklogIDFormIsRefused(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(s.Root, DirBacklog)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "S01-a-hunch.json"),
		[]byte(`{"id":"S01","title":"A hunch","kind":"idea","status":"open"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := s.ListBacklogItems()
	if err == nil {
		t.Fatal("a backlog item under a retired id form must not be read")
	}
	if !strings.Contains(err.Error(), "BL01") {
		t.Errorf("the refusal should name the current form: %v", err)
	}
}

// TestTheJournalKeepsTheOldNames covers the exception.
//
// The journal is history: its events keep the names and the ids they were written
// with, forever. The reader resolves what the journal names, and it does not refuse
// a plan because its history is old.
func TestTheJournalKeepsTheOldNames(t *testing.T) {
	ops, store := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendJournalEvent(JournalEntry{
		Seq: 1, At: Now(), Kind: KindControl, Actor: "line-cook", Session: "s-old", Role: "line-cook",
		Recipe: r.ID,
	}); err != nil {
		t.Fatal(err)
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatalf("a journal with an old name must read: %v", err)
	}
	if len(j.Events) == 0 {
		t.Fatal("the journal is empty")
	}
	if j.Events[len(j.Events)-1].Role != "line-cook" {
		t.Errorf("the old role should be kept as it was written: %q", j.Events[len(j.Events)-1].Role)
	}
}

// TestTheRetiredNamesAreOnlyWhereTheyBelong is the wire for "one format".
//
// The retired names may live in exactly two places: the detection that refuses them
// (`retired.go`) and the deny-list that refuses them as names (`model.go`). A third
// place is a second answer to a question the plan has already answered — and the
// tolerance would come back one reader at a time, which is how it came the first
// time.
//
// The scan reads the **code**, not the prose: the history is written down in the
// comments, and the comments are not a reader.
func TestTheRetiredNamesAreOnlyWhereTheyBelong(t *testing.T) {
	allowed := map[string]bool{
		"retired.go": true, // the detection that refuses them
		"model.go":   true, // the deny-list that refuses them as names
	}
	tokens := []string{
		"proofs",
		"specials",
		"prep-notes",
		"ListLegacyItems",
		"LegacyItem",
		"legacyFingerprint",
		"restampLegacy",
	}
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		if allowed[name] {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		code := stripComments(string(raw))
		for _, token := range tokens {
			if strings.Contains(code, token) {
				t.Errorf("%s reads the retired %q — the format is one format", name, token)
			}
		}
	}
}

// stripComments takes the prose out of a Go file, so a scan of the code does not
// trip on the history the comments tell.
func stripComments(src string) string {
	var out []string
	inBlock := false
	for _, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case inBlock:
			if i := strings.Index(trimmed, "*/"); i >= 0 {
				inBlock = false
				line = trimmed[i+2:]
			} else {
				continue
			}
		case strings.HasPrefix(trimmed, "/*"):
			if i := strings.Index(trimmed, "*/"); i >= 0 {
				line = trimmed[i+2:]
			} else {
				inBlock = true
				continue
			}
		case strings.HasPrefix(trimmed, "//"):
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}
