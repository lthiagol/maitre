package cli

import (
	"strings"
	"testing"
)

// TestACommentOnAMissingNoteIsRefused covers the comment's thread.
//
// A comment belongs to a tasting note or to the recipe itself. A note that is not
// there would put the comment in no thread at all: in the artifact, in the journal,
// and on no screen.
func TestACommentOnAMissingNoteIsRefused(t *testing.T) {
	dir := newPlan(t)
	create := cli(t, dir, "recipe", "create",
		"--title", "T",
		"--outcome", "o",
		"--problem", "p",
		"--out-of-scope", "a",
		"--out-of-scope", "b")
	if !create.ok(t) {
		t.Fatalf("create failed:\n%s", create.Stderr)
	}

	comment := cli(t, dir, "recipe", "comment", "R01", "--finding", "F-99", "--body", "nobody sees this")
	if comment.Code == 0 {
		t.Fatalf("a comment on a note that is not there should be refused:\n%s", comment.Stdout)
	}

	// And nothing was written.
	show := cli(t, dir, "recipe", "show", "R01")
	if !show.ok(t) {
		t.Fatalf("show failed:\n%s", show.Stderr)
	}
	if strings.Contains(show.Stdout, "nobody sees this") {
		t.Errorf("the comment was written anyway:\n%s", show.Stdout)
	}
}

// TestACommentHasAnAuthor covers the default.
//
// The journal records the actor, so a comment that said nobody would leave the
// artifact disagreeing with the log.
func TestACommentHasAnAuthor(t *testing.T) {
	dir := newPlan(t)
	create := cli(t, dir, "recipe", "create",
		"--title", "T",
		"--outcome", "o",
		"--problem", "p",
		"--out-of-scope", "a",
		"--out-of-scope", "b")
	if !create.ok(t) {
		t.Fatalf("create failed:\n%s", create.Stderr)
	}

	comment := cli(t, dir, "recipe", "comment", "R01", "--body", "a note about the recipe")
	if !comment.ok(t) {
		t.Fatalf("comment failed:\n%s", comment.Stderr)
	}

	show := cli(t, dir, "recipe", "show", "R01")
	if !show.ok(t) {
		t.Fatalf("show failed:\n%s", show.Stderr)
	}
	if strings.Contains(show.Stdout, `"author": ""`) {
		t.Errorf("the comment should have an author:\n%s", show.Stdout)
	}
}
