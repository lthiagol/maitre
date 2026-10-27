package core

import (
	"fmt"
)

// Seed builds a realistic plan directory.
//
// Everything is created through the Operations layer, so the seeded plan cannot drift
// from the model — and seeding is itself an exercise of the write path.
//
// It is used twice: as the console's test fixture, and behind `maitre-cli demo` so a human
// can see the console without setting up a project.
//
// The actors vary by phase, so the journal reads like a real session: the sous chef
// composes, the chef signs off, the runner executes, the reviewer reviews.
func Seed(store *Store) error {
	return seed(store, false)
}

// SeedStress writes the realistic plan plus the pathological recipes.
//
// The extra recipes are the cases that break a layout: a title with a newline, a title
// longer than the column, thirty criteria, fifty steps, twelve findings, an empty
// recipe, and text that is not ASCII. Use it to find out what the console can survive.
//
//	go run ./cmd/maitre-cli demo --dir /tmp/kb --chaos
func SeedStress(store *Store) error {
	return seed(store, true)
}

func seed(store *Store, stress bool) error {
	if err := store.Init(); err != nil {
		return err
	}

	chef := NewOperations(store, AuditContext{Actor: "chef", Session: "s-chef", Role: "chef"})
	sous := NewOperations(store, AuditContext{Actor: "sous-chef", Session: "s-sous", Role: "sous-chef"})
	runner := NewOperations(store, AuditContext{Actor: "runner", Session: "s-runner", Role: "runner"})
	reviewer := NewOperations(store, AuditContext{Actor: "reviewer", Session: "s-reviewer", Role: "reviewer"})

	// RCP01 — being shaped. One criterion cannot be proven: the filter matches nothing.
	if _, err := sous.CreateRecipe(RecipeSpec{
		Title:      "Sharpen the AC gate",
		Outcome:    "A zero-match verification can never satisfy a gate",
		Problem:    "v1 marked criteria passed on empty filters, which manufactured confidence",
		InScope:    []string{"the gate table", "the verifier"},
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the board renders 100 recipes", Verification: "go test ./internal/core -run Board"},
			{Description: "a zero-match proof is rejected", Verification: "go test ./internal/core -run ZeroMatch"},
		},
		WorkPackages: []string{"core"},
		Steps: []StepSpec{
			{Action: "write the board renderer", DoneWhen: "the board renders", Covers: []string{"AC-01"}},
			{Action: "add the zero-match gate", DoneWhen: "a rejected proof blocks finish", Covers: []string{"AC-02"}},
		},
	}); err != nil {
		return fmt.Errorf("seed RCP01: %w", err)
	}

	// RCP02 — waiting for the chef. One criterion proven, one not.
	r02, err := sous.CreateRecipe(RecipeSpec{
		Title:      "Sign off the menu",
		Outcome:    "The chef can approve a recipe from the console",
		Problem:    "v1's TUI was read-only, so every decision went through the CLI",
		InScope:    []string{"the sign-off action"},
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the board renders 100 recipes", Verification: "go test ./internal/core -run Board"},
			{Description: "the console offers sign off", Verification: "go test ./internal/ui -run SignOff"},
		},
		WorkPackages: []string{"ui"},
		Steps: []StepSpec{
			{Action: "render the board", DoneWhen: "the board renders", Covers: []string{"AC-01"}},
			{Action: "add the sign-off action", DoneWhen: "space signs off", Covers: []string{"AC-02"}},
		},
	})
	if err != nil {
		return fmt.Errorf("seed RCP02: %w", err)
	}
	if _, err := sous.Advance(r02.ID, EventGroom); err != nil {
		return fmt.Errorf("seed RCP02 groom: %w", err)
	}
	if _, err := sous.Advance(r02.ID, EventSpecify); err != nil {
		return fmt.Errorf("seed RCP02 specify: %w", err)
	}
	if _, err := chef.SignOffCriterion(r02.ID, "AC-01", "chef", "the board renders 100 recipes"); err != nil {
		return fmt.Errorf("seed RCP02 sign off: %w", err)
	}

	// RCP03 — in the tasting lane. Was executed, then a finding was filed and a criterion
	// regressed, so it is being remediated.
	r03, err := sous.CreateRecipe(RecipeSpec{
		Title:      "Wire the review queue",
		Outcome:    "Findings reach the chef without leaving the board",
		Problem:    "v1 had no review queue — findings lived in the recipe and the TUI never showed them",
		InScope:    []string{"the queue lane"},
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the board renders 100 recipes", Verification: "go test ./internal/core -run Board"},
			{Description: "the filter narrows the queue", Verification: "go test ./internal/ui -run Filter"},
			{Description: "findings close without a restart", Verification: "go test ./internal/ui -run Close"},
		},
		WorkPackages: []string{"ui"},
		Steps: []StepSpec{
			{Action: "add the queue", DoneWhen: "the queue lists findings", Covers: []string{"AC-01"}},
			{Action: "filter findings by severity", DoneWhen: "the filter narrows the list", Covers: []string{"AC-02"}},
			{Action: "close a finding in place", DoneWhen: "the finding closes without a restart", Covers: []string{"AC-03"}},
		},
	})
	if err != nil {
		return fmt.Errorf("seed RCP03: %w", err)
	}
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := sous.Advance(r03.ID, ev); err != nil {
			return fmt.Errorf("seed RCP03 %s: %w", ev, err)
		}
	}
	if _, err := chef.ApproveRecipe(r03.ID, "chef"); err != nil {
		return fmt.Errorf("seed RCP03 approve: %w", err)
	}
	if _, err := runner.ClaimRecipe(r03.ID, "runner"); err != nil {
		return fmt.Errorf("seed RCP03 claim: %w", err)
	}
	for _, g := range []string{"G1", "G2", "G3"} {
		if _, err := runner.SetStepStatus(r03.ID, g, "done"); err != nil {
			return fmt.Errorf("seed RCP03 step %s: %w", g, err)
		}
	}
	for _, ac := range []string{"AC-01", "AC-02", "AC-03"} {
		if _, err := chef.SignOffCriterion(r03.ID, ac, "chef", "checked against the diff"); err != nil {
			return fmt.Errorf("seed RCP03 sign off %s: %w", ac, err)
		}
	}
	if _, err := runner.Advance(r03.ID, EventFinish); err != nil {
		return fmt.Errorf("seed RCP03 finish: %w", err)
	}
	if _, err := reviewer.AddFinding(r03.ID, Finding{
		Phase:       "external",
		Severity:    "high",
		Category:    "correctness",
		Description: "the queue drops findings on refresh",
	}); err != nil {
		return fmt.Errorf("seed RCP03 finding: %w", err)
	}
	if _, err := runner.Advance(r03.ID, EventEnterRemediation); err != nil {
		return fmt.Errorf("seed RCP03 remediation: %w", err)
	}
	// A criterion regressed after the finding — which is why the board shows 1 of 3.
	if _, err := runner.VerifyCriterion(r03.ID, "AC-03"); err != nil {
		return fmt.Errorf("seed RCP03 verify: %w", err)
	}

	// RCP04 — approved and idle. It has no runner yet, and it waits on RCP02.
	r04, err := sous.CreateRecipe(RecipeSpec{
		Title:      "Add the backlog view",
		Outcome:    "Specials and prep notes are visible and promotable",
		Problem:    "v1 had twelve artifact types and no ladder, so nothing connected to anything",
		InScope:    []string{"the backlog lane"},
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the backlog lists specials", Verification: "go test ./internal/core -run Backlog"},
		},
		WorkPackages: []string{"ui"},
		Steps: []StepSpec{
			{Action: "add the backlog lane", DoneWhen: "the lane renders", Covers: []string{"AC-01"}},
		},
	})
	if err != nil {
		return fmt.Errorf("seed RCP04: %w", err)
	}
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := sous.Advance(r04.ID, ev); err != nil {
			return fmt.Errorf("seed RCP04 %s: %w", ev, err)
		}
	}
	if _, err := chef.ApproveRecipe(r04.ID, "chef"); err != nil {
		return fmt.Errorf("seed RCP04 approve: %w", err)
	}
	if _, err := chef.SignOffCriterion(r04.ID, "AC-01", "chef", "the backlog lists specials"); err != nil {
		return fmt.Errorf("seed RCP04 sign off: %w", err)
	}
	// RCP04 cannot run until RCP02 is signed off and closed.
	if _, err := sous.UpdateRecipe(r04.ID, Fragment{DependsOn: &[]string{"RCP02"}}); err != nil {
		return fmt.Errorf("seed RCP04 depends: %w", err)
	}

	// RCP05 — on the line. One criterion proven, one failing.
	r05, err := sous.CreateRecipe(RecipeSpec{
		Title:      "Sharpen the proof contract",
		Outcome:    "A proof carries a command, an exit code and a match count",
		Problem:    "v1's evidence was prose, so nothing could be checked",
		InScope:    []string{"the verifier"},
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the board renders 100 recipes", Verification: "go test ./internal/core -run Board"},
			{Description: "a proof that matched nothing is rejected", Verification: "go test ./internal/core -run ZeroMatch"},
		},
		WorkPackages: []string{"core"},
		Steps: []StepSpec{
			{Action: "add the match states", DoneWhen: "three states exist", Covers: []string{"AC-01"}},
			{Action: "record the unit", DoneWhen: "the unit is recorded", Covers: []string{"AC-02"}},
		},
	})
	if err != nil {
		return fmt.Errorf("seed RCP05: %w", err)
	}
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := sous.Advance(r05.ID, ev); err != nil {
			return fmt.Errorf("seed RCP05 %s: %w", ev, err)
		}
	}
	if _, err := chef.ApproveRecipe(r05.ID, "chef"); err != nil {
		return fmt.Errorf("seed RCP05 approve: %w", err)
	}
	if _, err := runner.ClaimRecipe(r05.ID, "runner"); err != nil {
		return fmt.Errorf("seed RCP05 claim: %w", err)
	}
	if _, err := runner.SetStepStatus(r05.ID, "G1", "done"); err != nil {
		return fmt.Errorf("seed RCP05 step: %w", err)
	}
	if _, err := runner.SetStepStatus(r05.ID, "G2", "in-progress"); err != nil {
		return fmt.Errorf("seed RCP05 step: %w", err)
	}
	if _, err := chef.SignOffCriterion(r05.ID, "AC-01", "chef", "three states exist"); err != nil {
		return fmt.Errorf("seed RCP05 sign off: %w", err)
	}
	if _, err := runner.VerifyCriterion(r05.ID, "AC-02"); err != nil {
		return fmt.Errorf("seed RCP05 verify: %w", err)
	}

	// RCP06 — served.
	r06, err := sous.CreateRecipe(RecipeSpec{
		Title:      "Draft the glossary",
		Outcome:    "One vocabulary for the tool and the methodology",
		Problem:    "v1 said spec, prompt and milestone for the same thing, and nobody knew which was canonical",
		InScope:    []string{"the glossary"},
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "every term has an avoid list", Verification: "go test ./internal/core -run Glossary"},
		},
		WorkPackages: []string{"docs"},
		Steps: []StepSpec{
			{Action: "write the glossary", DoneWhen: "the file exists", Covers: []string{"AC-01"}},
		},
	})
	if err != nil {
		return fmt.Errorf("seed RCP06: %w", err)
	}
	for _, ev := range []Event{EventGroom, EventSpecify} {
		if _, err := sous.Advance(r06.ID, ev); err != nil {
			return fmt.Errorf("seed RCP06 %s: %w", ev, err)
		}
	}
	if _, err := chef.ApproveRecipe(r06.ID, "chef"); err != nil {
		return fmt.Errorf("seed RCP06 approve: %w", err)
	}
	if _, err := runner.ClaimRecipe(r06.ID, "runner"); err != nil {
		return fmt.Errorf("seed RCP06 claim: %w", err)
	}
	if _, err := runner.SetStepStatus(r06.ID, "G1", "done"); err != nil {
		return fmt.Errorf("seed RCP06 step: %w", err)
	}
	if _, err := chef.SignOffCriterion(r06.ID, "AC-01", "chef", "every term has an avoid list"); err != nil {
		return fmt.Errorf("seed RCP06 sign off: %w", err)
	}
	if _, err := runner.Advance(r06.ID, EventFinish); err != nil {
		return fmt.Errorf("seed RCP06 finish: %w", err)
	}
	if _, err := reviewer.AddFinding(r06.ID, Finding{
		Phase:       "external",
		Severity:    "low",
		Category:    "docs",
		Description: "the empty state is unclear",
		Status:      "fixed",
		Resolved:    Now(),
	}); err != nil {
		return fmt.Errorf("seed RCP06 finding: %w", err)
	}
	if _, err := reviewer.AddFinding(r06.ID, Finding{
		Phase:       "external",
		Severity:    "low",
		Category:    "docs",
		Description: "the glossary needs a worked example",
	}); err != nil {
		return fmt.Errorf("seed RCP06 finding: %w", err)
	}
	if _, err := runner.ResolveFinding(r06.ID, "F-02", "added an example"); err != nil {
		return fmt.Errorf("seed RCP06 resolve: %w", err)
	}
	if _, err := reviewer.RecordReviewRound(r06.ID, "reviewer", "ok", "the glossary is clear"); err != nil {
		return fmt.Errorf("seed RCP06 round: %w", err)
	}
	if _, err := reviewer.Advance(r06.ID, EventPassReview); err != nil {
		return fmt.Errorf("seed RCP06 pass review: %w", err)
	}
	if _, err := chef.DocumentRecipe(r06.ID, "One vocabulary, defined once, mapped to CPD."); err != nil {
		return fmt.Errorf("seed RCP06 document: %w", err)
	}
	if _, err := chef.Advance(r06.ID, EventClose); err != nil {
		return fmt.Errorf("seed RCP06 close: %w", err)
	}

	// RCP07 — being shaped, nothing proven yet.
	if _, err := sous.CreateRecipe(RecipeSpec{
		Title:      "Record the review protocol",
		Outcome:    "The hand-off points are named and audited",
		Problem:    "v1's L5 audit was advisory, so the session boundary was a convention rather than a check",
		InScope:    []string{"the four hand-off points"},
		OutOfScope: []string{"web shell", "orchestration"},
		Criteria: []CriterionSpec{
			{Description: "the protocol names four hand-offs", Verification: "go test ./internal/core -run Protocol"},
		},
		WorkPackages: []string{"docs"},
		Steps: []StepSpec{
			{Action: "write the protocol", DoneWhen: "the four points are documented", Covers: []string{"AC-01"}},
		},
	}); err != nil {
		return fmt.Errorf("seed RCP07: %w", err)
	}
	// RCP07 waits on RCP05.
	r07, err := sous.GetRecipe("RCP07")
	if err != nil {
		return fmt.Errorf("seed RCP07 read: %w", err)
	}
	if _, err := sous.UpdateRecipe(r07.ID, Fragment{DependsOn: &[]string{"RCP05"}}); err != nil {
		return fmt.Errorf("seed RCP07 depends: %w", err)
	}

	if stress {
		if err := seedStress(chef); err != nil {
			return err
		}
	}

	// The lower rungs — the backlog the chef works through.
	for _, sp := range []struct {
		title string
		body  string
		tags  []string
	}{
		{"The board needs a way to show proofs", "A proof the chef cannot read is a claim, not a proof.", []string{"board"}},
		{"A recipe with no ACs is probably a special", "If nothing can falsify it, it is a hunch wearing a costume.", []string{"model"}},
		{"The chef wants a weekly digest", "What changed, what is blocked, what needs a decision.", []string{"reporting"}},
		{"Trying the board on a real project", "Run the console against a plan that is not ours.", []string{"trials"}},
	} {
		if _, err := chef.CreateSpecial(sp.title, sp.body, sp.tags); err != nil {
			return fmt.Errorf("seed special %q: %w", sp.title, err)
		}
	}
	for _, pn := range []struct {
		title  string
		intent string
		scope  []string
	}{
		{"Wire the review queue", "Findings reach the chef without leaving the board.", []string{"the queue lane"}},
		{"Sharpen the proof contract", "A proof carries a command, an exit code and a match count.", []string{"the verifier"}},
	} {
		if _, err := chef.CreatePrepNote(pn.title, pn.intent, pn.scope); err != nil {
			return fmt.Errorf("seed prep note %q: %w", pn.title, err)
		}
	}

	return nil
}
