package core

import (
	"strings"
	"testing"
)

// TestPlanStatsCountsThePlan covers the numbers the Home tab shows.
//
// Every count is the domain's own, and the counts are the ones the triage, the
// evidence summary, the open findings and the backlog's groups already give. A
// number the console counted itself would be a number the CLI counts differently.
func TestPlanStatsCountsThePlan(t *testing.T) {
	ops, store := newOps(t)

	// A recipe the chef must act on: it is ready, waiting for a sign-off.
	ready, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := ops.Advance(ready.ID, ev); err != nil {
			t.Fatal(err)
		}
	}

	// A recipe with a criterion proven and a criterion failed.
	executed, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := ops.Advance(executed.ID, ev); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := ops.ApproveRecipe(executed.ID, "chef"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.ClaimRecipe(executed.ID, "runner"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.SetStepStatus(executed.ID, "G1", "done"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.SignOffCriterion(executed.ID, "AC-01", "chef", "checked"); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.Advance(executed.ID, EventFinish); err != nil {
		t.Fatal(err)
	}
	// An open tasting note, filed by the external tasting.
	if _, err := ops.AddFinding(executed.ID, Finding{
		Phase:       FindingPhaseExternal,
		Severity:    SeverityHigh,
		Category:    "correctness",
		Description: "the queue drops findings on refresh",
	}); err != nil {
		t.Fatal(err)
	}

	// A backlog with an idea and a small task.
	if _, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := ops.CreateBacklogItem("Wire the queue", "it drops findings", nil, KindFeature, nil); err != nil {
		t.Fatal(err)
	}

	recipes, err := store.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	items, err := store.ListBacklogItems()
	if err != nil {
		t.Fatal(err)
	}
	journal, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	stats := PlanStatsOf(recipes, items, journal)

	if stats.Recipes.Total != 2 {
		t.Errorf("total = %d, want 2", stats.Recipes.Total)
	}
	if stats.Recipes.NeedsYou != 2 {
		t.Errorf("need you = %d, want 2 — the ready recipe waits for the chef, and the one with the open note waits for the triage", stats.Recipes.NeedsYou)
	}
	if stats.Recipes.Served != 0 {
		t.Errorf("served = %d, want 0", stats.Recipes.Served)
	}
	if stats.Evidence.Criteria != 2 {
		t.Errorf("criteria = %d, want 2", stats.Evidence.Criteria)
	}
	if stats.Evidence.Proven != 1 {
		t.Errorf("proven = %d, want 1", stats.Evidence.Proven)
	}
	if stats.Tasting.OpenNotes != 1 {
		t.Errorf("open notes = %d, want 1", stats.Tasting.OpenNotes)
	}
	if stats.Tasting.OpenFindings != 1 {
		t.Errorf("open findings = %d, want 1 — the note came from the external tasting", stats.Tasting.OpenFindings)
	}
	if stats.Backlog.Ideas != 1 {
		t.Errorf("ideas = %d, want 1", stats.Backlog.Ideas)
	}
	if stats.Backlog.SmallTasks != 1 {
		t.Errorf("small tasks = %d, want 1", stats.Backlog.SmallTasks)
	}
	if stats.Automation.PlanVersion != DomainVersion {
		t.Errorf("plan version = %d, want %d", stats.Automation.PlanVersion, DomainVersion)
	}
	if stats.Automation.BuildVersion != DomainVersion {
		t.Errorf("build version = %d, want %d", stats.Automation.BuildVersion, DomainVersion)
	}
	if !stats.Automation.VersionsAgree {
		t.Error("the versions should agree — the plan was written by this build")
	}
	if stats.Automation.LastEvent == nil {
		t.Fatal("the plan has events — the last one should be there")
	}
}

// TestPlanStatsOfAnEmptyPlan covers the plan with nothing in it.
//
// A count of zero is a count, and the tab must show it — not nothing at all.
func TestPlanStatsOfAnEmptyPlan(t *testing.T) {
	stats := PlanStatsOf(nil, nil, Journal{})

	if stats.Recipes.Total != 0 || stats.Evidence.Criteria != 0 {
		t.Errorf("an empty plan should count zero: %+v", stats)
	}
	if stats.Automation.LastEvent != nil {
		t.Error("a plan with no journal has no last event")
	}
	if stats.Automation.VersionsAgree {
		t.Error("a plan that was never written agrees with nothing — the versions do not agree")
	}
}

// TestThePlanStatsViewIsHeldInBothDirections covers the wire.
//
// The view is registered in `core.DerivedViews` **and** written down in `docs/03`
// § 8. The wire checks the table both ways, and it must fail when either side goes:
// a view the table does not name is a view a reader cannot find, and a row the table
// has and the code does not is a view that never happens.
func TestThePlanStatsViewIsHeldInBothDirections(t *testing.T) {
	doc := readDoc(t, "03-domain-model.md")

	t.Run("the row deleted from the table", func(t *testing.T) {
		broken := strings.Replace(doc, "| `plan_stats` | `PlanStatsOf` |", "| `invented` | `PlanStatsOf` |", 1)
		if broken == doc {
			t.Fatal("the table does not hold the row — the corruption did not apply")
		}
		if err := checkDerivedViews(broken); err == nil {
			t.Fatal("the wire should fail when the table does not name the view")
		}
	})

	t.Run("the registration removed", func(t *testing.T) {
		orig := DerivedViews
		t.Cleanup(func() { DerivedViews = orig })

		views := make([]DerivedView, 0, len(orig))
		for _, v := range orig {
			if v.Name == "plan_stats" {
				continue
			}
			views = append(views, v)
		}
		DerivedViews = views

		if err := checkViewCoverage("."); err == nil {
			t.Fatal("the wire should fail when a view is registered nowhere")
		}
	})
}
