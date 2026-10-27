package core

import (
	"strings"
	"testing"
)

// TestAFindingRefusesAnUnknownPhase covers the closed set.
//
// The gates read the phase — `review.findings_exist` counts the external ones — so a
// typo would leave the recipe stuck between two gates that disagree about what
// "external" means.
func TestAFindingRefusesAnUnknownPhase(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}

	_, err = ops.AddFinding(r.ID, Finding{
		Phase:       "externl",
		Severity:    SeverityHigh,
		Category:    "correctness",
		Description: "the phase is misspelled",
	})
	if err == nil {
		t.Fatal("an unknown phase must be refused")
	}
	if !strings.Contains(err.Error(), FindingPhaseSelf) || !strings.Contains(err.Error(), FindingPhaseExternal) {
		t.Fatalf("the refusal should name the allowed phases: %v", err)
	}

	// Nothing was written.
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Review.Findings) != 0 {
		t.Fatalf("the refused finding was recorded: %+v", got.Review.Findings)
	}
}

// TestAFindingRefusesAnUnknownSeverity covers the other closed set.
func TestAFindingRefusesAnUnknownSeverity(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}

	_, err = ops.AddFinding(r.ID, Finding{
		Phase:       FindingPhaseExternal,
		Severity:    "catastrophic",
		Category:    "correctness",
		Description: "the severity is not in the set",
	})
	if err == nil {
		t.Fatal("an unknown severity must be refused")
	}
	for _, want := range FindingSeverityLevels {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("the refusal should name %q: %v", want, err)
		}
	}
}

// TestAFindingAcceptsTheClosedSet is the other side: the check is a check, not a wall.
func TestAFindingAcceptsTheClosedSet(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	for _, phase := range FindingPhases {
		for _, severity := range FindingSeverityLevels {
			if _, err := ops.AddFinding(r.ID, Finding{
				Phase:       phase,
				Severity:    severity,
				Category:    "correctness",
				Description: "a real problem",
			}); err != nil {
				t.Fatalf("%s/%s: %v", phase, severity, err)
			}
		}
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Review.Findings) != len(FindingPhases)*len(FindingSeverityLevels) {
		t.Fatalf("got %d findings, want %d", len(got.Review.Findings), len(FindingPhases)*len(FindingSeverityLevels))
	}
}
