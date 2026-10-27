package core

import "testing"

// TestEveryMoveIsInTheClosedSet keeps the version fingerprint honest.
//
// The fingerprint hashes the chef's moves, so a move that is not in the set is a rule
// the version stamp cannot see — and an older console files a recipe under "Needs
// you" with nothing to press. Every move the domain can produce has to be in the
// set.
func TestEveryMoveIsInTheClosedSet(t *testing.T) {
	in := func(m ChefMove) bool {
		for _, known := range AllChefMoves {
			if known == m {
				return true
			}
		}
		return false
	}
	for _, s := range everyShape() {
		if m := ChefMoveOf(s.recipe); !in(m) {
			t.Errorf("%s (%s): the move %q is not in the closed set", s.recipe.ID, s.recipe.Lifecycle, m)
		}
	}
}

func TestLaneDerivation(t *testing.T) {
	cases := []struct {
		name string
		r    Recipe
		want Lane
	}{
		{"draft", Recipe{Lifecycle: LifecycleDraft}, LaneGrooming},
		{"groomed", Recipe{Lifecycle: LifecycleGroomed}, LaneGrooming},
		{"ready waits for sign-off", Recipe{Lifecycle: LifecycleReady}, LaneApproval},
		{"approved is on the line", Recipe{Lifecycle: LifecycleApproved}, LaneExecution},
		{"executing is on the line", Recipe{Lifecycle: LifecycleExecuting}, LaneExecution},
		{"executed is in tasting", Recipe{Lifecycle: LifecycleExecuted}, LaneReview},
		{"remediation is in tasting", Recipe{Lifecycle: LifecycleRemediation}, LaneReview},
		{"reviewed is in tasting", Recipe{Lifecycle: LifecycleReviewed}, LaneReview},
		{"closed is served", Recipe{Lifecycle: LifecycleClosed}, LaneDone},
		{
			"blocked wins over the lifecycle",
			Recipe{Lifecycle: LifecycleExecuting, Overlays: Overlays{Blocked: &Overlay{Reason: "x"}}},
			LaneBlocked,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := LaneOf(c.r); got != c.want {
				t.Fatalf("lane = %q, want %q", got, c.want)
			}
		})
	}
}

func TestReviewStateDerivation(t *testing.T) {
	open := Finding{ID: "F-01", Phase: "external", Status: "open"}
	fixed := Finding{ID: "F-01", Phase: "external", Status: "fixed"}

	if got := ReviewStateOf(Recipe{}); got != ReviewNone {
		t.Fatalf("state = %q, want none", got)
	}
	if got := ReviewStateOf(Recipe{Review: Review{Findings: []Finding{open}}}); got != ReviewOpen {
		t.Fatalf("state = %q, want open-findings", got)
	}
	if got := ReviewStateOf(Recipe{Review: Review{Findings: []Finding{fixed}}}); got != ReviewRemediated {
		t.Fatalf("state = %q, want remediated", got)
	}
	if got := ReviewStateOf(Recipe{Review: Review{Rounds: []Round{{Verdict: "ok"}}}}); got != ReviewSelfSigned {
		// Nobody signed for the work, so nobody's signature on the review makes it
		// independent. A round without a claimant is not a clean tasting.
		t.Fatalf("state = %q, want self-reviewed", got)
	}
	if got := ReviewStateOf(Recipe{
		Execution: Execution{ClaimedBy: "runner"},
		Review:    Review{Rounds: []Round{{Reviewer: "reviewer", Verdict: "ok"}}},
	}); got != ReviewClean {
		t.Fatalf("state = %q, want reviewed-clean", got)
	}
}

func TestReadinessDerivation(t *testing.T) {
	if r := ReadinessOf(Recipe{Lifecycle: LifecycleApproved}); !r.Ready {
		t.Fatalf("approved recipe should be ready: %+v", r)
	}
	if r := ReadinessOf(Recipe{Lifecycle: LifecycleReady}); r.Ready {
		t.Fatal("a recipe awaiting sign-off must not be ready")
	}
	if r := ReadinessOf(Recipe{Lifecycle: LifecycleApproved, Overlays: Overlays{Blocked: &Overlay{Reason: "waiting"}}}); r.Ready {
		t.Fatal("a blocked recipe must not be ready")
	}
}

func TestEvidenceSummarySeparatesRejectedFromUnproven(t *testing.T) {
	r := Recipe{Acceptance: []Criterion{
		{ID: "AC-01", Evidence: []Evidence{{ExitCode: 0, Matched: 3, Match: MatchProven}}},
		{ID: "AC-02", Evidence: []Evidence{{ExitCode: 0, Matched: 0, Match: MatchRejected}}},
		{ID: "AC-03", Evidence: []Evidence{{ExitCode: 0, Match: MatchUnknown}}},
		{ID: "AC-04"},
	}}
	s := EvidenceSummaryOf(r)
	if s.Criteria != 4 {
		t.Fatalf("criteria = %d, want 4", s.Criteria)
	}
	if s.Proven != 1 {
		t.Fatalf("proven = %d, want 1", s.Proven)
	}
	if s.Rejected != 1 {
		t.Fatalf("rejected = %d, want 1", s.Rejected)
	}
	if s.Unproven != 1 {
		t.Fatalf("unproven = %d, want 1", s.Unproven)
	}
	if s.Pending != 1 {
		t.Fatalf("pending = %d, want 1", s.Pending)
	}
}

func TestBoardGroupsAndCounts(t *testing.T) {
	recipes := []Recipe{
		{ID: "R01", Lifecycle: LifecycleDraft},
		{ID: "R02", Lifecycle: LifecycleDraft},
		{ID: "R03", Lifecycle: LifecycleReady},
		{ID: "R04", Lifecycle: LifecycleClosed},
	}
	b := BoardOf("/tmp/kitchen", Project{Name: "demo"}, recipes)

	if b.Total() != 4 {
		t.Fatalf("total = %d, want 4", b.Total())
	}
	if b.Count(LaneGrooming) != 2 {
		t.Fatalf("grooming = %d, want 2", b.Count(LaneGrooming))
	}
	if b.Count(LaneApproval) != 1 {
		t.Fatalf("approval = %d, want 1", b.Count(LaneApproval))
	}
	if b.Count(LaneDone) != 1 {
		t.Fatalf("done = %d, want 1", b.Count(LaneDone))
	}
	if b.Count(LaneExecution) != 0 {
		t.Fatalf("execution = %d, want 0", b.Count(LaneExecution))
	}
	if got := b.Summaries(LaneGrooming); len(got) != 2 || got[0].ID != "R01" {
		t.Fatalf("summaries = %+v", got)
	}
}

func TestBoardSortsNumerically(t *testing.T) {
	b := BoardOf("/tmp/kitchen", Project{}, []Recipe{
		{ID: "R10", Lifecycle: LifecycleDraft},
		{ID: "R2", Lifecycle: LifecycleDraft},
	})
	got := b.Summaries(LaneGrooming)
	if got[0].ID != "R2" || got[1].ID != "R10" {
		t.Fatalf("order = %s, %s", got[0].ID, got[1].ID)
	}
}
