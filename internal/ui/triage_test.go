package ui

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/kitchen-brigade/internal/core"
)

// selfReviewedModel builds a console over a recipe the executor reviewed itself.
//
// This is the state the gate refuses to close: the only round is the executor's, so
// nobody independent has tasted the work.
func selfReviewedModel(t *testing.T) *Model {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "line-cook", Session: "s-1", Role: "runner"})
	// The chef is the human: only the chef may sign a recipe off.
	chef := core.NewOperations(store, core.AuditContext{Actor: "chef", Session: "s-2", Role: "coordinator"})

	r, err := ops.CreateRecipe(core.RecipeSpec{
		Title:      "The console shows a finding's conversation",
		Outcome:    "The recipe is the channel between the reviewer and the fixer",
		Problem:    "The console showed findings with no thread",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []core.CriterionSpec{
			{Description: "the console groups a finding's comments under it", Verification: "go test ./internal/ui -run TestFindingCommentsAreListed"},
		},
		WorkPackages: []string{"ui"},
		Steps:        []core.StepSpec{{Action: "render the thread", DoneWhen: "the comments show", Covers: []string{"AC-01"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range []core.Event{core.EventGroom, core.EventSpecify} {
		if _, err := ops.Advance(r.ID, ev); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := chef.ApproveRecipe(r.ID, "chef"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.ClaimRecipe(r.ID, "line-cook"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.SetStepStatus(r.ID, "G1", "done"); err != nil {
		t.Fatal(err)
	}
	if _, err := chef.SignOffCriterion(r.ID, "AC-01", "chef", "checked"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(r.ID, core.EventFinish); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.RecordReviewRound(r.ID, "line-cook", "ok", "looks good to me"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.DocumentRecipe(r.ID, "The console groups a finding's thread under it."); err != nil {
		t.Fatal(err)
	}

	m := loadModel(t, store)
	m.height = 30
	return m
}

// TestTheConsoleSaysWhatIsWaiting covers the triage gap.
//
// The gate refused to close the recipe while the console filed it under "On the
// line" and counted "0 need you". The chef was never told.
func TestTheConsoleSaysWhatIsWaiting(t *testing.T) {
	m := selfReviewedModel(t)
	out := StripANSI(m.View().Content)

	if !strings.Contains(out, "1 need you") {
		t.Errorf("the count should say one recipe needs the chef:\n%s", out)
	}
	if !strings.Contains(out, "Needs you") {
		t.Errorf("the recipe should be filed under 'Needs you':\n%s", out)
	}
	if strings.Contains(out, "On the line") {
		t.Errorf("the recipe should not be filed under 'On the line':\n%s", out)
	}
}

// TestTheDetailSaysWhy covers the reason. "Needs you" without a reason is a queue;
// with one, it is an instruction.
func TestTheDetailSaysWhy(t *testing.T) {
	m := selfReviewedModel(t)
	rec, ok := m.currentRecipe()
	if !ok {
		t.Fatal("no recipe selected")
	}
	out := StripANSI(strings.Join(m.recipeDetail(rec, m.width), "\n"))

	if !strings.Contains(out, "awaiting an independent review") {
		t.Errorf("the detail should say what is waiting:\n%s", out)
	}
	if !strings.Contains(out, "self-reviewed") {
		t.Errorf("the detail should not call a self-review clean:\n%s", out)
	}
	if strings.Contains(out, "tasting clean") {
		t.Errorf("the detail should not call a self-review clean:\n%s", out)
	}
}
