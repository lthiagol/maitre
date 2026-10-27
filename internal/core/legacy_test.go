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

// TestTheReaderReadsTheOldEvidence covers the rename of a criterion's record.
//
// The record was `proofs`; it is `evidence` now. The records are append-only and they
// are the recipe's history: a reader that drops them cannot say whether the work was
// verified.
func TestTheReaderReadsTheOldEvidence(t *testing.T) {
	s := newStore(t)
	writeRawRecipe(t, s, "RCP01-old.json", `{
  "id": "RCP01",
  "title": "written before the rename",
  "slug": "written-before-the-rename",
  "lifecycle": "executed",
  "acceptance": [
    {
      "id": "AC-01",
      "description": "the board renders",
      "verification": "go test ./internal/core -run Board",
      "proofs": [
        {"command": "go test ./internal/core -run Board", "exit_code": 0, "matched": 1, "match": "proven"}
      ]
    }
  ]
}`)

	r, err := s.LoadRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Acceptance) != 1 {
		t.Fatalf("criteria = %d, want 1", len(r.Acceptance))
	}
	if len(r.Acceptance[0].Evidence) != 1 {
		t.Fatalf("evidence = %d, want 1 — the old field name carried the record", len(r.Acceptance[0].Evidence))
	}
	if r.Acceptance[0].Evidence[0].Match != MatchProven {
		t.Fatalf("match = %q, want %q", r.Acceptance[0].Evidence[0].Match, MatchProven)
	}
	if st := r.Acceptance[0].Status(); st != CriterionPassed {
		t.Fatalf("status = %q, want %q — the record the old field name carried still counts", st, CriterionPassed)
	}
}

// TestTheNewEvidenceWins covers the other half.
//
// A file that carries both names is a broken plan, and the reader does not merge two
// answers to one question: the new name is the one the tool writes, so it is the one
// that counts.
func TestTheNewEvidenceWins(t *testing.T) {
	s := newStore(t)
	writeRawRecipe(t, s, "RCP01-both.json", `{
  "id": "RCP01",
  "title": "carrying both names",
  "slug": "carrying-both-names",
  "lifecycle": "executed",
  "acceptance": [
    {
      "id": "AC-01",
      "description": "the board renders",
      "verification": "go test ./internal/core -run Board",
      "proofs": [
        {"command": "go test ./internal/core -run Board", "exit_code": 0, "matched": 1, "match": "proven"}
      ],
      "evidence": [
        {"command": "go test ./internal/core -run Board", "exit_code": 1, "matched": 0, "match": "rejected"}
      ]
    }
  ]
}`)

	r, err := s.LoadRecipe("RCP01")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Acceptance[0].Evidence) != 1 {
		t.Fatalf("evidence = %d, want 1", len(r.Acceptance[0].Evidence))
	}
	if r.Acceptance[0].Evidence[0].Match != MatchRejected {
		t.Fatalf("match = %q, want %q — the new field name wins", r.Acceptance[0].Evidence[0].Match, MatchRejected)
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
