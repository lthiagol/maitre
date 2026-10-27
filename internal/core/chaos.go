package core

import (
	"fmt"
	"strings"
)

// seedStress adds the pathological recipes on top of the realistic plan.
//
// These are not decoration. Each one is a case that has broken a console before, or
// would: a title that breaks the row, a recipe too long to read, an empty one, text that
// is not ASCII. Everything here is written through the Operations layer, like the rest
// of the plan, so it cannot drift from the model.
//
// The titles are prefixed so they are findable while poking at the console.
func seedStress(ops *Operations) error {
	chef := ops

	// A title with a newline and a tab. Before the sanitiser, this one corrupts the
	// screen.
	if _, err := chef.CreateRecipe(RecipeSpec{
		Title:      "STRESS: a title with\na newline and\ta tab",
		Outcome:    "The console survives text that was not written for a terminal",
		Problem:    "Every title in the console came from an agent",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the row does not break", Verification: "go test ./internal/ui -run Sanitize"},
		},
		WorkPackages: []string{"stress"},
		Steps: []StepSpec{
			{Action: "sanitise the title", DoneWhen: "the row stays one row", Covers: []string{"AC-01"}},
		},
	}); err != nil {
		return fmt.Errorf("seed stress newline title: %w", err)
	}

	// A title longer than the column.
	long := "STRESS: a title that is much much much longer than any column the console has"
	if _, err := chef.CreateRecipe(RecipeSpec{
		Title:      long,
		Outcome:    "A long title is truncated, not wrapped",
		Problem:    "Titles have no length limit",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the title is truncated", Verification: "go test ./internal/ui -run Truncate"},
		},
		WorkPackages: []string{"stress"},
		Steps: []StepSpec{
			{Action: "truncate the title", DoneWhen: "the column holds", Covers: []string{"AC-01"}},
		},
	}); err != nil {
		return fmt.Errorf("seed stress long title: %w", err)
	}

	// A title that is not ASCII: emoji, combining marks, and a right-to-left run.
	if _, err := chef.CreateRecipe(RecipeSpec{
		Title:      "STRESS: 日本語のタイトル 🎉 and a RTL run ‮abc",
		Outcome:    "Widths are measured, not guessed",
		Problem:    "A wide character is two cells wide",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the columns still line up", Verification: "go test ./internal/ui -run Width"},
		},
		WorkPackages: []string{"stress"},
		Steps: []StepSpec{
			{Action: "measure the width", DoneWhen: "the columns line up", Covers: []string{"AC-01"}},
		},
	}); err != nil {
		return fmt.Errorf("seed stress unicode title: %w", err)
	}

	// A title of one character.
	if _, err := chef.CreateRecipe(RecipeSpec{
		Title:      "x",
		Outcome:    "A one-character title still works",
		Problem:    "Nothing stops an agent writing one word",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the row still lines up", Verification: "go test ./internal/ui -run Short"},
		},
		WorkPackages: []string{"stress"},
		Steps: []StepSpec{
			{Action: "handle a short title", DoneWhen: "the row lines up", Covers: []string{"AC-01"}},
		},
	}); err != nil {
		return fmt.Errorf("seed stress short title: %w", err)
	}

	// Thirty acceptance criteria — more than fit on a screen.
	many, err := chef.CreateRecipe(RecipeSpec{
		Title:        "STRESS: thirty acceptance criteria",
		Outcome:      "A recipe can have more criteria than the screen has rows",
		Problem:      "Nothing caps the number of criteria",
		OutOfScope:   []string{"web shell", "orchestration"},
		WorkPackages: []string{"stress"},
	})
	if err != nil {
		return fmt.Errorf("seed stress many criteria: %w", err)
	}
	for i := 0; i < 30; i++ {
		if _, err := chef.AddCriterion(many.ID,
			fmt.Sprintf("criterion number %02d is checked", i+1),
			fmt.Sprintf("go test ./internal/core -run Case%02d", i+1)); err != nil {
			return fmt.Errorf("seed stress criterion %d: %w", i, err)
		}
	}

	// A recipe with no acceptance criteria at all.
	if _, err := chef.CreateRecipe(RecipeSpec{
		Title:        "STRESS: no acceptance criteria",
		Outcome:      "An empty recipe has an empty state, not a broken one",
		Problem:      "A recipe can be created before it has criteria",
		OutOfScope:   []string{"web shell", "orchestration"},
		WorkPackages: []string{"stress"},
		Steps: []StepSpec{
			{Action: "show the empty state", DoneWhen: "nothing looks broken", Covers: []string{}},
		},
	}); err != nil {
		return fmt.Errorf("seed stress no criteria: %w", err)
	}

	// A criterion with no verification, and one with a very long command.
	//
	// This one cannot be groomed: the ac.falsifiable gate refuses a criterion with no
	// verification. That is correct — the recipe simply stays in draft, which is
	// exactly where such a recipe belongs.
	if _, err := chef.CreateRecipe(RecipeSpec{
		Title:      "STRESS: odd verifications",
		Outcome:    "A criterion with no verification still reads",
		Problem:    "A criterion can be written with nothing to run",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "this one has no verification", Verification: ""},
			{Description: "this one has a very long command", Verification: "go test " + strings.Repeat("./internal/core ", 30)},
		},
		WorkPackages: []string{"stress"},
		Steps: []StepSpec{
			{Action: "render the odd verifications", DoneWhen: "both read", Covers: []string{"AC-01", "AC-02"}},
		},
	}); err != nil {
		return fmt.Errorf("seed stress odd verifications: %w", err)
	}

	// Fifty guidance steps.
	longSteps, err := chef.CreateRecipe(RecipeSpec{
		Title:        "STRESS: fifty guidance steps",
		Outcome:      "A method can be longer than the screen",
		Problem:      "Nothing caps the number of steps",
		OutOfScope:   []string{"web shell", "orchestration"},
		WorkPackages: []string{"stress"},
	})
	if err != nil {
		return fmt.Errorf("seed stress many steps: %w", err)
	}
	for i := 0; i < 50; i++ {
		if _, err := chef.AddStep(longSteps.ID, StepSpec{
			Action:   fmt.Sprintf("step number %02d of the method", i+1),
			DoneWhen: "the step is done",
		}); err != nil {
			return fmt.Errorf("seed stress step %d: %w", i, err)
		}
	}

	// Twelve open tasting notes.
	found, err := chef.CreateRecipe(RecipeSpec{
		Title:      "STRESS: twelve open tasting notes",
		Outcome:    "A recipe can be swamped with findings",
		Problem:    "Nothing caps the number of findings",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the notes are all listed", Verification: "go test ./internal/ui -run Findings"},
		},
		WorkPackages: []string{"stress"},
		Steps: []StepSpec{
			{Action: "list the findings", DoneWhen: "all twelve show", Covers: []string{"AC-01"}},
		},
	})
	if err != nil {
		return fmt.Errorf("seed stress findings: %w", err)
	}
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := chef.Advance(found.ID, ev); err != nil {
			return fmt.Errorf("seed stress findings %s: %w", ev, err)
		}
	}
	if _, err := chef.ApproveRecipe(found.ID, "chef"); err != nil {
		return fmt.Errorf("seed stress findings approve: %w", err)
	}
	if _, err := chef.ClaimRecipe(found.ID, "runner"); err != nil {
		return fmt.Errorf("seed stress findings claim: %w", err)
	}
	if _, err := chef.SetStepStatus(found.ID, "G1", "done"); err != nil {
		return fmt.Errorf("seed stress findings step: %w", err)
	}
	if _, err := chef.SignOffCriterion(found.ID, "AC-01", "chef", "checked against the diff"); err != nil {
		return fmt.Errorf("seed stress findings sign off: %w", err)
	}
	if _, err := chef.Advance(found.ID, EventFinish); err != nil {
		return fmt.Errorf("seed stress findings finish: %w", err)
	}
	for i := 0; i < 12; i++ {
		if _, err := chef.AddFinding(found.ID, Finding{
			Phase:       "external",
			Severity:    []string{"low", "medium", "high"}[i%3],
			Category:    []string{"correctness", "docs", "performance"}[i%3],
			Description: fmt.Sprintf("finding number %02d, which has something to say about the code", i+1),
		}); err != nil {
			return fmt.Errorf("seed stress finding %d: %w", i, err)
		}
	}
	if _, err := chef.Advance(found.ID, EventEnterRemediation); err != nil {
		return fmt.Errorf("seed stress findings remediation: %w", err)
	}

	// A long deviation note.
	dev, err := chef.CreateRecipe(RecipeSpec{
		Title:      "STRESS: a long deviation note",
		Outcome:    "A long note is wrapped, not truncated",
		Problem:    "An agent can write a lot",
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the note reads", Verification: "go test ./internal/ui -run Deviation"},
		},
		WorkPackages: []string{"stress"},
		Steps: []StepSpec{
			{Action: "record the deviation", DoneWhen: "it is recorded", Covers: []string{"AC-01"}},
		},
	})
	if err != nil {
		return fmt.Errorf("seed stress deviation: %w", err)
	}
	if _, err := chef.RecordDeviation(dev.ID, "G1", strings.Repeat(
		"The method turned out to be different from what was written, and here is the long account of why. ", 6)); err != nil {
		return fmt.Errorf("seed stress deviation note: %w", err)
	}

	// The backlog's rows break the same ways: a title with a newline, a title
	// longer than the column, and a body that is not ASCII.
	for _, it := range []struct {
		title string
		body  string
		kind  string
	}{
		{"STRESS: a backlog title with\na newline and\ta tab", "the row does not break", KindIdea},
		{"STRESS: …much longer than any column the console has, and it keeps going", "the title is cut", KindFeature},
		{"STRESS: 日本語のタイトル 🎉", "not ascii", KindBug},
	} {
		if _, err := chef.CreateBacklogItem(it.title, it.body, nil, it.kind, nil); err != nil {
			return fmt.Errorf("seed stress backlog item %q: %w", it.title, err)
		}
	}

	return nil
}
