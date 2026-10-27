package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeRawRecipe puts a recipe file on disk, exactly as it is written.
func writeRawRecipe(t *testing.T, s *Store, name, body string) {
	t.Helper()
	dir := filepath.Join(s.Root, DirRecipes)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestTheReaderReadsTheOldCompletionStamp covers the rename of the recipe's stamp.
//
// The stamp was `evidence`; it is `completion` now. A plan written before the
// rename keeps what it had, and a reader that drops the stamp loses the moment the
// recipe was done — so the reader reads the old name.
func TestTheReaderReadsTheOldCompletionStamp(t *testing.T) {
	s := newStore(t)
	writeRawRecipe(t, s, "RCP01-old.json", `{
  "id": "RCP01",
  "title": "written before the rename",
  "slug": "written-before-the-rename",
  "lifecycle": "closed",
  "evidence": {"completed_at": "2026-10-26T02:21:56Z"}
}`)

	r, err := s.LoadRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	if r.Completion.CompletedAt != "2026-10-26T02:21:56Z" {
		t.Fatalf("completed_at = %q, want the stamp the old field name carried", r.Completion.CompletedAt)
	}
}

// TestTheNewCompletionStampWins covers the other half.
//
// A file that carries both names is a broken plan, and the reader does not merge two
// answers to one question: the new name is the one the tool writes, so it is the one
// that counts.
func TestTheNewCompletionStampWins(t *testing.T) {
	s := newStore(t)
	writeRawRecipe(t, s, "RCP01-both.json", `{
  "id": "RCP01",
  "title": "carrying both names",
  "slug": "carrying-both-names",
  "lifecycle": "closed",
  "evidence": {"completed_at": "2026-10-26T02:21:56Z"},
  "completion": {"completed_at": "2026-10-27T09:00:00Z"}
}`)

	r, err := s.LoadRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	if r.Completion.CompletedAt != "2026-10-27T09:00:00Z" {
		t.Fatalf("completed_at = %q, want the stamp the new field name carried", r.Completion.CompletedAt)
	}
}

// TestThePlanKeepsTheStampWhenItIsWrittenAgain covers the write.
//
// The reader tolerates the old name; the writer does not. A plan the tool writes
// again carries the new field name, and the stamp it carried is still there.
func TestThePlanKeepsTheStampWhenItIsWrittenAgain(t *testing.T) {
	s := newStore(t)
	writeRawRecipe(t, s, "RCP01-old.json", `{
  "id": "RCP01",
  "title": "written before the rename",
  "slug": "written-before-the-rename",
  "lifecycle": "closed",
  "evidence": {"completed_at": "2026-10-26T02:21:56Z"}
}`)

	r, err := s.LoadRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SaveRecipe(r); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(filepath.Join(s.Root, DirRecipes, "RCP01-written-before-the-rename.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"completion"`) {
		t.Fatalf("the written recipe does not carry the new field name:\n%s", raw)
	}
	if strings.Contains(string(raw), `"evidence"`) {
		t.Fatalf("the written recipe still carries the old field name:\n%s", raw)
	}
}
