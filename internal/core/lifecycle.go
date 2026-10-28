package core

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Lifecycle is the recipe's state. Nine states, every one reachable and producible.
type Lifecycle string

const (
	LifecycleDraft       Lifecycle = "draft"
	LifecycleGroomed     Lifecycle = "groomed"
	LifecycleReady       Lifecycle = "ready"
	LifecycleApproved    Lifecycle = "approved"
	LifecycleExecuting   Lifecycle = "executing"
	LifecycleExecuted    Lifecycle = "executed"
	LifecycleRemediation Lifecycle = "remediation"
	LifecycleReviewed    Lifecycle = "reviewed"
	LifecycleClosed      Lifecycle = "closed"
)

// AllLifecycles lists every state in forward order.
var AllLifecycles = []Lifecycle{
	LifecycleDraft, LifecycleGroomed, LifecycleReady, LifecycleApproved,
	LifecycleExecuting, LifecycleExecuted, LifecycleRemediation,
	LifecycleReviewed, LifecycleClosed,
}

// Valid reports whether l is a known state.
func (l Lifecycle) Valid() bool {
	for _, c := range AllLifecycles {
		if c == l {
			return true
		}
	}
	return false
}

// Terminal reports whether the state admits no further transition.
func (l Lifecycle) Terminal() bool { return l == LifecycleClosed }

// Event is a transition request. It is the only way a recipe's state changes.
type Event string

const (
	EventGroom            Event = "groom"
	EventSpecify          Event = "specify"
	EventApprove          Event = "approve"
	EventRequestChanges   Event = "request_changes"
	EventStart            Event = "start"
	EventFinish           Event = "finish"
	EventEnterRemediation Event = "enter_remediation"
	EventExitRemediation  Event = "exit_remediation"
	EventPassReview       Event = "pass_review"
	EventClose            Event = "close"
	EventAbortToDraft     Event = "abort_to_draft"
	EventReopen           Event = "reopen"
)

// AllEvents lists every transition event.
var AllEvents = []Event{
	EventGroom, EventSpecify, EventApprove, EventRequestChanges, EventStart,
	EventFinish, EventEnterRemediation, EventExitRemediation, EventPassReview,
	EventClose, EventAbortToDraft, EventReopen,
}

// TransitionContext carries what an event needs beyond the recipe's own state.
type TransitionContext struct {
	// Approver is required by EventApprove and must be a human.
	Approver string
	// Reason is required by EventAbortToDraft.
	Reason string
	// Actor is recorded on the audit entry.
	Actor string
	// Session and Role are recorded for the L5 session-boundary check.
	Session string
	Role    string
	// Closed reports whether a recipe is closed.
	//
	// The gate table is package-level and cannot close over the store, so the
	// operations layer fills this in when it builds the context. The
	// `dependencies.satisfied` gate asks it.
	//
	// A context that carries no way to read the plan makes that gate refuse: a gate
	// that cannot verify must not wave the transition through.
	Closed func(id string) (bool, error)
	// Now defaults to time.Now.
	Now time.Time
}

// GateError reports one or more failed gates.
type GateError struct {
	Event Event
	Gates []string
}

func (e *GateError) Error() string {
	return fmt.Sprintf("%s: gate(s) not satisfied: %s", e.Event, strings.Join(e.Gates, ", "))
}

// RefusalError reports that the plan refused an action.
//
// A refusal is not a failure: the plan's rules turned the action away. It is
// journalled, because an attempt that was turned away is evidence — a bad request
// ("no such recipe") is not, because nothing was attempted.
type RefusalError struct {
	Reason string
}

func (e *RefusalError) Error() string { return e.Reason }

// TransitionError reports an illegal transition.
type TransitionError struct {
	Event  Event
	From   Lifecycle
	Reason string
}

func (e *TransitionError) Error() string {
	return fmt.Sprintf("%s: cannot go from %q: %s", e.Event, e.From, e.Reason)
}

// gate is one named check guarding a transition.
type gate struct {
	name  string
	check func(*Recipe, TransitionContext) error
}

// MinOutOfScope is how many out-of-scope items a groomed recipe must declare.
const MinOutOfScope = 2

// DomainVersion is the version of the domain's *rules*: the gate table, the
// transition table, the lifecycles, the events, the chef's moves and the rules on a
// tasting note's outcome.
//
// A plan records the version it was last written under, in its journal. **A tool
// writes a plan only if the plan's version is the tool's version.** A plan stamped
// with an older version is not supported — there is no upgrade path, because a newer
// tool would apply rules the plan was not written under. A plan stamped with a newer
// version is not supported either: the tool does not know the rules it was written
// under. Reading is always allowed: the chef can look at the plan.
//
// The **version** fences the **writer**. The **shape** is the **tripwire**: it hashes
// the gates, the transitions, the lifecycles, the events, the chef's moves, the note's
// outcome rules and the facts of the plan's format. A rule that is not in that hash
// is a rule the tripwire cannot see — so a rule that changes has to be added to it,
// or the plan cannot say which rules it was written under.
//
// **Two facts, two jobs.** The shape is the code-side fact: it says the format's
// facts changed, so the version must move with them. The version is the plan-side
// fact: it says which rules the plan was written under. The version does not say that
// the format changed — the shape does.
//
// So the two move apart. **Bump this when the rules the tool applies change** — a
// gate added, renamed or dropped, a transition added, a lifecycle added, a chef's
// move added, a note's outcome rule changed, a record renamed: the tool would write
// something the plan's rules do not allow, and the version is what fences it.
// `TestDomainShapeIsVersioned` fails until the shape is updated.
//
// **The bump protects the future, not the corpus.** The guard refuses every plan
// whose version is not the tool's version, so a plan stamped 1 is not written by a
// version-6 tool at all: the corpus is not carried along, and nothing is migrated.
//
// **Version 6: the plan format's records were renamed.** The recipe's completion
// stamp was `evidence`; it is `completion` now. A criterion's record was `proofs`;
// it is `evidence` now — the name the journal's event kind already carried, so the
// code ends up agreeing with its own journal. The stamp says when the recipe was
// done; the records say what was verified, and they live on the criterion.
//
// The rules the gates apply did not change — no gate was weakened and no gate rule
// moved — but the **record** did: a version-5 tool would write `proofs` where the
// version-6 writer writes `evidence`.
//
// The version is bumped because the **writer** changed. The version-6 reader does not
// read `proofs` either: the format is one format, and a record under a retired name
// is a record the current reader does not know.
//
// **Version 5: the plan format changed again — the backlog is one artifact.** The
// backlog's two file types (`kitchen/specials/` and `kitchen/prep-notes/`) are one
// (`kitchen/backlog/`), the item carries a `kind`, and its id is `BL01`. The rules
// the gates apply did not change — no gate was weakened and no gate rule moved — but
// the **record** did: a version-4 tool would write a special or a prep note into
// `kitchen/specials/` and `kitchen/prep-notes/`, and the version-5 writer does not.
// The version is bumped because the **writer** changed.
//
// **Version 4: the stamp covers the plan's format.** The fingerprint now hashes the
// facts of the format — the id prefixes, the directories, the root files and the field
// names on disk — so a format change is caught by the tripwire instead of by hand.
// Before this, the format was versioned by hand: a tool could write old-format
// records into a new-format plan and nothing noticed.
//
// **Version 3: the plan format changed.** A recipe's id is `RCP01` and its file is
// named `RCP01-<slug>.json`, where version 2 wrote `R01` and `R01-<slug>.json`. The
// rules the gates apply did not change — the *record* they apply to did, and a tool
// that reads a version-3 plan with the version-2 id machinery would not find the
// recipes in it.
//
// The version is bumped because the **writer** changed: a version-2 tool would write
// `R01` into a plan the version-3 reader expects to find `RCP01`. See
// docs/07-vocabulary.md.
const DomainVersion = 6

// gates is the single gate table. Every gate lives here and nowhere else.
var gates = []gate{
	{"intent.present", func(r *Recipe, _ TransitionContext) error {
		if strings.TrimSpace(r.Intent.Outcome) == "" {
			return fmt.Errorf("intent.outcome is empty")
		}
		return nil
	}},
	{"problem.present", func(r *Recipe, _ TransitionContext) error {
		if strings.TrimSpace(r.Intent.Problem) == "" {
			return fmt.Errorf("intent.problem is empty")
		}
		return nil
	}},
	{"scope.out_of_scope_declared", func(r *Recipe, _ TransitionContext) error {
		if len(r.Scope.OutOfScope) < MinOutOfScope {
			return fmt.Errorf("needs %d out_of_scope items, has %d", MinOutOfScope, len(r.Scope.OutOfScope))
		}
		return nil
	}},
	{"ac.present", func(r *Recipe, _ TransitionContext) error {
		if len(r.Acceptance) == 0 {
			return fmt.Errorf("no acceptance criteria")
		}
		return nil
	}},
	{"ac.falsifiable", func(r *Recipe, _ TransitionContext) error {
		for _, c := range r.Acceptance {
			if err := checkFalsifiable(c); err != nil {
				return fmt.Errorf("%s: %w", c.ID, err)
			}
		}
		return nil
	}},
	{"questions.resolved", func(r *Recipe, _ TransitionContext) error {
		for _, q := range r.Questions {
			if q.Status != "resolved" {
				return fmt.Errorf("%s is still open", q.ID)
			}
		}
		return nil
	}},
	{"guidance.decomposed", func(r *Recipe, _ TransitionContext) error {
		if len(r.WorkPackages) == 0 {
			return fmt.Errorf("no work packages")
		}
		if len(r.Guidance) == 0 {
			return fmt.Errorf("no guidance steps")
		}
		return nil
	}},
	{"ac.covered", func(r *Recipe, _ TransitionContext) error {
		for _, c := range r.Acceptance {
			if !coversAny(r.Guidance, c.ID) {
				return fmt.Errorf("%s has no covering step", c.ID)
			}
		}
		return nil
	}},
	{"guidance.done_when_present", func(r *Recipe, _ TransitionContext) error {
		for _, s := range r.Guidance {
			if strings.TrimSpace(s.DoneWhen) == "" {
				return fmt.Errorf("%s has no done_when", s.ID)
			}
		}
		return nil
	}},
	{"approval.by_human", func(r *Recipe, ctx TransitionContext) error {
		// The approval record is created by the transition, so the gate checks the
		// transition's context.
		//
		// The chef is the human; everyone else in the crew is an agent. The
		// deny-list is the agent roles, so a human may approve under any name they
		// like — but the actor performing the approval may not be one of the
		// agents. That is the permission boundary ADR-0006 promised: the runner is
		// unable to perform a human action, not merely asked not to.
		if strings.TrimSpace(ctx.Approver) == "" {
			return fmt.Errorf("no approver given")
		}
		if IsAgentRole(ctx.Actor) {
			return fmt.Errorf("%s is an agent role: only the chef signs a recipe off", ctx.Actor)
		}
		return nil
	}},
	{"dependencies.satisfied", func(r *Recipe, ctx TransitionContext) error {
		for _, dep := range r.DependsOn {
			if strings.TrimSpace(dep) == "" {
				continue
			}
			if ctx.Closed == nil {
				// A gate that cannot check must refuse, not wave the transition
				// through. v1's lesson: a green gate that proves nothing.
				return fmt.Errorf("cannot check dependency %s: the transition carries no way to read the plan", dep)
			}
			closed, err := ctx.Closed(dep)
			if err != nil {
				return fmt.Errorf("dependency %s: %w", dep, err)
			}
			if !closed {
				return fmt.Errorf("dependency %s is not closed", dep)
			}
		}
		return nil
	}},
	{"steps.settled", func(r *Recipe, _ TransitionContext) error {
		for _, s := range r.Guidance {
			if s.Status != "done" && s.Status != "skipped" {
				return fmt.Errorf("%s is %s", s.ID, s.Status)
			}
		}
		return nil
	}},
	{"ac.all_proven", func(r *Recipe, _ TransitionContext) error {
		for _, c := range r.Acceptance {
			switch st := c.Status(); st {
			case CriterionPassed:
				// Evidence that matched something, or a human sign-off.
			case CriterionDeferred:
				// A deferred criterion is re-filed as a follow-up. It does not
				// keep the recipe open: CPD's rule is that a recipe is done when
				// its in-scope criteria pass.
			default:
				return fmt.Errorf("%s is %s", c.ID, st)
			}
		}
		return nil
	}},
	{"review.findings_exist", func(r *Recipe, _ TransitionContext) error {
		if n := countOpenFindings(r, "external"); n == 0 {
			return fmt.Errorf("no open external findings")
		}
		return nil
	}},
	{"review.findings_closed", func(r *Recipe, _ TransitionContext) error {
		if n := countOpenFindings(r, "external"); n > 0 {
			return fmt.Errorf("%d open external findings", n)
		}
		return nil
	}},
	{"review.verdict_ok", func(r *Recipe, _ TransitionContext) error {
		if len(r.Review.Rounds) == 0 {
			return fmt.Errorf("no review round recorded")
		}
		latest := r.Review.Rounds[len(r.Review.Rounds)-1]
		if latest.Verdict != "ok" {
			return fmt.Errorf("the last review round by %s says %q", latest.Reviewer, latest.Verdict)
		}
		return nil
	}},
	{"review.reviewer_independent", func(r *Recipe, _ TransitionContext) error {
		// The reviewer must not be the agent that did the work. CPD principle 13:
		// "the author should not be the only reviewer", and the L5 law: the session
		// that executes is never the session that reviews.
		//
		// This cannot stop an agent lying about its name — the session id is a string
		// anyone can type. What it does stop is the accident: the executor signing off
		// its own work by reflex, and the anonymous claim that made any reviewer look
		// independent.
		//
		// The check itself lives in views.go, so this gate and the console's triage
		// cannot drift apart. One fact, one answer.
		if len(r.Review.Rounds) == 0 {
			return fmt.Errorf("no review round recorded")
		}
		if SelfReviewed(*r) {
			if r.Execution.ClaimedBy == "" {
				return fmt.Errorf("the work has no claimant, so the review by %s cannot be independent",
					r.Review.Rounds[len(r.Review.Rounds)-1].Reviewer)
			}
			return fmt.Errorf("the last review round was recorded by %s, who did the work", r.Execution.ClaimedBy)
		}
		return nil
	}},
	{"review.no_open_findings", func(r *Recipe, _ TransitionContext) error {
		if n := OpenFindings(r); n > 0 {
			return fmt.Errorf("%d open findings", n)
		}
		return nil
	}},
	{"review.clean", func(r *Recipe, _ TransitionContext) error {
		if n := OpenFindings(r); n > 0 {
			return fmt.Errorf("%d open findings", n)
		}
		return nil
	}},
	{"audit.documented", func(r *Recipe, _ TransitionContext) error {
		if strings.TrimSpace(r.Audit.Documentation) == "" {
			return fmt.Errorf("audit trail has no documentation")
		}
		return nil
	}},
	{"abort.recorded", func(r *Recipe, ctx TransitionContext) error {
		if strings.TrimSpace(ctx.Reason) == "" {
			return fmt.Errorf("abort needs a reason")
		}
		return nil
	}},
	{"changes.recorded", func(r *Recipe, ctx TransitionContext) error {
		if strings.TrimSpace(ctx.Reason) == "" {
			return fmt.Errorf("requesting changes needs a reason")
		}
		return nil
	}},
	{"immutability.respected", func(r *Recipe, _ TransitionContext) error {
		// The sections frozen at the sign-off must still say what they said when
		// the chef signed. A hand edit to them is a material change, and the
		// artifact must not assert a sign-off over content that is no longer
		// there.
		if r.Approval == nil {
			return nil
		}
		want, err := r.Fingerprint()
		if err != nil {
			return err
		}
		if want != r.Approval.Fingerprint {
			return fmt.Errorf("frozen sections changed since approval")
		}
		return nil
	}},
}

// placeholders are verification strings that look filled in but prove nothing.
var placeholders = map[string]bool{
	"todo": true, "tbd": true, "fixme": true, "xxx": true, "...": true, "-": true, "n/a": true,
}

// checkFalsifiable rejects a verification that cannot be run.
//
// A verification is **resolvable** when the tool can determine a match count from it:
// a `go test` command, whose output the tool counts (ADR-0006 §3), or a `manual:`
// procedure a human signs off. Anything else would be *unproven* — and unproven
// evidence never satisfies a gate, so the criterion could never be finished. Saying
// so at grooming time, when the human reads it, beats saying so at finish time.
func checkFalsifiable(c Criterion) error {
	v := strings.TrimSpace(c.Verification)
	if v == "" {
		return fmt.Errorf("verification is empty")
	}
	if placeholders[strings.ToLower(v)] {
		return fmt.Errorf("verification is a placeholder")
	}
	if strings.HasPrefix(strings.ToLower(v), "manual:") {
		if strings.TrimSpace(v[len("manual:"):]) == "" {
			return fmt.Errorf("manual verification has no procedure")
		}
		return nil
	}
	if !isGoTest(v) {
		return fmt.Errorf("verification %q is not a command the tool can count matches for (want a go test command, or manual: with a procedure)", v)
	}
	return nil
}

func coversAny(steps []Step, criterionID string) bool {
	for _, s := range steps {
		for _, c := range s.Covers {
			if c == criterionID {
				return true
			}
		}
	}
	return false
}

func countOpenFindings(r *Recipe, phase string) int {
	n := 0
	for _, f := range OpenFindingsOf(r) {
		if phase == "" || f.Phase == phase {
			n++
		}
	}
	return n
}

// transitionDef describes one legal transition.
type transitionDef struct {
	from  []Lifecycle
	to    Lifecycle
	gates []string
}

// transitions is the single transition table. Guards live in the gates table and run
// inside the transition, so a transition can never skip its gates (constraint C4).
var transitions = map[Event]transitionDef{
	EventGroom: {
		from: []Lifecycle{LifecycleDraft},
		to:   LifecycleGroomed,
		gates: []string{"intent.present", "problem.present", "scope.out_of_scope_declared",
			"ac.present", "ac.falsifiable", "questions.resolved"},
	},
	EventSpecify: {
		from:  []Lifecycle{LifecycleGroomed},
		to:    LifecycleReady,
		gates: []string{"guidance.decomposed", "ac.covered", "guidance.done_when_present"},
	},
	EventApprove: {
		from:  []Lifecycle{LifecycleReady},
		to:    LifecycleApproved,
		gates: []string{"approval.by_human"},
	},
	EventRequestChanges: {
		from:  []Lifecycle{LifecycleReady},
		to:    LifecycleGroomed,
		gates: []string{"changes.recorded"},
	},
	EventStart: {
		from:  []Lifecycle{LifecycleApproved},
		to:    LifecycleExecuting,
		gates: []string{"dependencies.satisfied", "immutability.respected"},
	},
	EventFinish: {
		from:  []Lifecycle{LifecycleExecuting},
		to:    LifecycleExecuted,
		gates: []string{"steps.settled", "ac.all_proven", "immutability.respected"},
	},
	EventEnterRemediation: {
		from:  []Lifecycle{LifecycleExecuted, LifecycleReviewed},
		to:    LifecycleRemediation,
		gates: []string{"review.findings_exist", "immutability.respected"},
	},
	EventExitRemediation: {
		from:  []Lifecycle{LifecycleRemediation},
		to:    LifecycleExecuted,
		gates: []string{"review.findings_closed", "immutability.respected"},
	},
	EventPassReview: {
		from:  []Lifecycle{LifecycleExecuted},
		to:    LifecycleReviewed,
		gates: []string{"review.verdict_ok", "review.reviewer_independent", "review.no_open_findings", "immutability.respected"},
	},
	EventClose: {
		from: []Lifecycle{LifecycleReviewed},
		to:   LifecycleClosed,
		// The review gates are repeated here on purpose: `close` is the last chance
		// to say "this was never actually reviewed by anyone else". A recipe that
		// reached `reviewed` on the executor's own round must not be able to close.
		gates: []string{"review.verdict_ok", "review.reviewer_independent", "review.clean", "audit.documented", "immutability.respected"},
	},
	EventAbortToDraft: {
		from:  []Lifecycle{LifecycleExecuting, LifecycleReady, LifecycleApproved},
		to:    LifecycleDraft,
		gates: []string{"abort.recorded"},
	},
	EventReopen: {
		from:  []Lifecycle{LifecycleClosed},
		to:    LifecycleExecuting,
		gates: nil,
	},
}

// gateByName indexes the gate table.
var gateByName = func() map[string]gate {
	m := make(map[string]gate, len(gates))
	for _, g := range gates {
		m[g.name] = g
	}
	return m
}()

// GateNames lists every gate, sorted.
func GateNames() []string {
	out := make([]string, 0, len(gates))
	for _, g := range gates {
		out = append(out, g.name)
	}
	sort.Strings(out)
	return out
}

// Transition applies an event to a recipe.
//
// The gates run inside the transition, so there is no way to move a recipe without them
// having passed. On success the recipe's lifecycle, approval and audit trail are updated
// in place.
func Transition(r *Recipe, ev Event, ctx TransitionContext) error {
	def, ok := transitions[ev]
	if !ok {
		return fmt.Errorf("unknown event %q", ev)
	}
	if !containsLifecycle(def.from, r.Lifecycle) {
		return &TransitionError{Event: ev, From: r.Lifecycle,
			Reason: "not a legal source state"}
	}
	if r.Overlays.Cancelled != nil {
		return &TransitionError{Event: ev, From: r.Lifecycle, Reason: "recipe is cancelled"}
	}

	var failed []string
	for _, name := range def.gates {
		g, ok := gateByName[name]
		if !ok {
			return fmt.Errorf("gate %q is not defined", name)
		}
		if err := g.check(r, ctx); err != nil {
			failed = append(failed, fmt.Sprintf("%s (%v)", name, err))
		}
	}
	if len(failed) > 0 {
		return &GateError{Event: ev, Gates: failed}
	}

	now := ctx.Now
	if now.IsZero() {
		now = time.Now()
	}
	stamp := now.UTC().Format(time.RFC3339)

	switch ev {
	case EventApprove:
		fp, err := r.Fingerprint()
		if err != nil {
			return err
		}
		r.Approval = &Approval{By: ctx.Approver, At: stamp, Fingerprint: fp}
	case EventFinish:
		// The recipe-level completion stamp. The per-criterion records are the
		// evidence; this is the moment the recipe was done.
		r.Completion.CompletedAt = stamp
	case EventAbortToDraft:
		// The sign-off and the claim are void: the premise was wrong, so the recipe
		// goes back to shaping. Leaving them would freeze the sections the approval
		// fingerprinted, and the recipe could never be edited again.
		r.Approval = nil
		r.Execution = Execution{}
		r.Audit.Deviations = append(r.Audit.Deviations, Deviation{
			At: stamp, Note: "aborted to draft: " + ctx.Reason, By: ctx.Actor,
		})
	}

	r.Lifecycle = def.to
	r.Updated = stamp
	return nil
}

// Block marks the recipe blocked with a reason.
func Block(r *Recipe, reason, by string, now time.Time) error {
	if r.Terminal() {
		return &TransitionError{Event: "block", From: r.Lifecycle, Reason: "recipe is terminal"}
	}
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("block needs a reason")
	}
	r.Overlays.Blocked = &Overlay{Reason: reason, By: by, At: stampOf(now)}
	return nil
}

// Unblock clears the blocked overlay.
func Unblock(r *Recipe) {
	r.Overlays.Blocked = nil
}

// Defer parks the recipe with a reason.
func Defer(r *Recipe, reason, by string, now time.Time) error {
	if r.Terminal() {
		return &TransitionError{Event: "defer", From: r.Lifecycle, Reason: "recipe is terminal"}
	}
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("defer needs a reason")
	}
	r.Overlays.Deferred = &Overlay{Reason: reason, By: by, At: stampOf(now)}
	return nil
}

// Resume clears the deferred overlay.
func Resume(r *Recipe) {
	r.Overlays.Deferred = nil
}

// Cancel marks the recipe cancelled. Terminal.
func Cancel(r *Recipe, reason, by string, now time.Time) error {
	if r.Overlays.Cancelled != nil {
		// The plan's answer, not a bad request: the recipe is cancelled, and
		// cancelling it again is an attempt the plan turns away.
		return &RefusalError{Reason: "recipe is already cancelled"}
	}
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("cancel needs a reason")
	}
	r.Overlays.Cancelled = &Overlay{Reason: reason, By: by, At: stampOf(now)}
	return nil
}

func containsLifecycle(list []Lifecycle, c Lifecycle) bool {
	for _, x := range list {
		if x == c {
			return true
		}
	}
	return false
}

func stampOf(t time.Time) string {
	if t.IsZero() {
		t = time.Now()
	}
	return t.UTC().Format(time.RFC3339)
}
