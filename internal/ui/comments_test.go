package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// recipeWithComments builds a recipe with a finding that has a thread, and a comment
// that belongs to no finding.
func recipeWithComments(t *testing.T) *Model {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "taster", Session: "s-1", Role: "reviewer"})

	r, err := ops.CreateRecipe(core.RecipeSpec{
		Title:      "Wire the review queue",
		Outcome:    "Findings reach the chef without leaving the board",
		Problem:    "v1 had no review queue",
		OutOfScope: []string{"web shell", "orchestration"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ops.AddFinding(r.ID, core.Finding{
		Phase:       "external",
		Severity:    "high",
		Category:    "correctness",
		Description: "the queue drops findings on refresh",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.AddComment(r.ID, "taster", "reproduced — refresh() re-reads the queue", "F-01"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.AddComment(r.ID, "line-cook", "read the queue once", "F-01"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.AddComment(r.ID, "chef", "the review as a whole looks fine", ""); err != nil {
		t.Fatal(err)
	}
	return loadModel(t, store)
}

// detail is the detail pane for the only recipe.
func detail(t *testing.T, m *Model) string {
	t.Helper()
	rec, ok := m.currentRecipe()
	if !ok {
		t.Fatal("no recipe selected")
	}
	return StripANSI(strings.Join(m.recipeDetail(rec, m.width), "\n"))
}

// TestFindingCommentsAreListed covers AC-01.
func TestFindingCommentsAreListed(t *testing.T) {
	m := recipeWithComments(t)
	out := detail(t, m)

	if !strings.Contains(out, "the queue drops findings on refresh") {
		t.Errorf("the finding should be shown:\n%s", out)
	}
	if !strings.Contains(out, "reproduced — refresh() re-reads the queue") {
		t.Errorf("the finding's comment should be listed under it:\n%s", out)
	}
	if !strings.Contains(out, "read the queue once") {
		t.Errorf("the second comment should be listed too:\n%s", out)
	}
}

// TestFindingCommentShowsItsAuthor covers AC-02.
func TestFindingCommentShowsItsAuthor(t *testing.T) {
	m := recipeWithComments(t)
	out := detail(t, m)

	if !strings.Contains(out, "taster:") {
		t.Errorf("a comment should show its author:\n%s", out)
	}
	if !strings.Contains(out, "line-cook:") {
		t.Errorf("each comment should show its own author:\n%s", out)
	}
}

// TestGeneralCommentsAreVisible covers AC-03.
func TestGeneralCommentsAreVisible(t *testing.T) {
	m := recipeWithComments(t)
	out := detail(t, m)

	if !strings.Contains(out, "the review as a whole looks fine") {
		t.Errorf("a comment that belongs to no finding should still be visible:\n%s", out)
	}
	if !strings.Contains(out, "chef:") {
		t.Errorf("it should show its author:\n%s", out)
	}
}

// TestTheCommentSitsUnderItsFinding checks the thread is grouped, not scattered.
func TestTheCommentSitsUnderItsFinding(t *testing.T) {
	m := recipeWithComments(t)
	out := detail(t, m)

	finding := strings.Index(out, "the queue drops findings on refresh")
	comment := strings.Index(out, "reproduced — refresh() re-reads the queue")
	if finding < 0 || comment < 0 {
		t.Fatalf("both should be present:\n%s", out)
	}
	if comment < finding {
		t.Errorf("the comment should come after the finding it belongs to:\n%s", out)
	}
}
