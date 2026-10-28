package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// AuditContext identifies who is acting.
//
// It lives in the Operations struct rather than in each call, so **every write is stamped
// automatically** and a caller cannot forget to identify itself. CPD's L5 law is audited
// by reading the journal, and that requires every write to carry its session.
type AuditContext struct {
	Actor   string
	Session string
	Role    string
}

// Operations is the domain layer's public verb set.
//
// Both surfaces call it. The CLI is a thin adapter; the console is a thin adapter. Neither
// invents an operation the domain does not support. See docs/decisions/0007-operation-layer.md.
type Operations struct {
	store    *Store
	context  AuditContext
	verifier *Verifier
}

// NewOperations returns the verb set over a store, acting as ctx.
func NewOperations(store *Store, ctx AuditContext) *Operations {
	return &Operations{
		store:   store,
		context: ctx,
		verifier: &Verifier{
			Dir:     store.ProjectRoot(),
			Timeout: DefaultVerifyTimeout,
		},
	}
}

// Store exposes the underlying store for reads the domain does not wrap.
func (o *Operations) Store() *Store { return o.store }

// Context returns the audit context.
func (o *Operations) Context() AuditContext { return o.context }

// RecipeSpec is the input to CreateRecipe.
type RecipeSpec struct {
	Title        string
	Outcome      string
	Problem      string
	InScope      []string
	OutOfScope   []string
	Criteria     []CriterionSpec
	WorkPackages []string
	Steps        []StepSpec
}

// CriterionSpec is one acceptance criterion to author.
type CriterionSpec struct {
	Description  string
	Verification string
}

// StepSpec is one guidance step to author.
type StepSpec struct {
	Action   string
	DoneWhen string
	Covers   []string
	// WorkPackage is the package the step belongs to. Empty means the first one.
	WorkPackage string
}

// Fragment is a partial update. A nil field is "leave alone"; a non-nil field replaces.
// That is what makes the update fragment-first: one section at a time.
type Fragment struct {
	Title      *string
	Outcome    *string
	Problem    *string
	InScope    *[]string
	OutOfScope *[]string
	DependsOn  *[]string
}

// ---------------------------------------------------------------------------
// Reading
// ---------------------------------------------------------------------------

// GetRecipe returns one recipe.
func (o *Operations) GetRecipe(id string) (*Recipe, error) {
	return o.store.LoadRecipe(id)
}

// ListRecipes returns every recipe, ordered by id.
func (o *Operations) ListRecipes() ([]Recipe, error) {
	return o.store.ListRecipes()
}

// GetProject returns the plan-level metadata.
func (o *Operations) GetProject() (Project, error) {
	return o.store.LoadProject()
}

// GetJournal returns the plan's event log.
func (o *Operations) GetJournal() (Journal, error) {
	return o.store.LoadJournal()
}

// ---------------------------------------------------------------------------
// Writing — the artifact
// ---------------------------------------------------------------------------

// CreateRecipe authors a new recipe in the draft state.
func (o *Operations) CreateRecipe(spec RecipeSpec) (*Recipe, error) {
	unlock, err := o.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	// The id is allocated under the lock: two creations must not be handed the same
	// id, and the id is a scan of what exists.
	id, err := o.store.NextRecipeID()
	if err != nil {
		return nil, err
	}
	r, err := o.createRecipeLocked(id, spec)
	if err != nil {
		return nil, err
	}
	if err := o.journalLocked(KindCreated, r.ID, map[string]any{"title": r.Title}); err != nil {
		return nil, err
	}
	return r, nil
}

// createRecipeLocked writes a new recipe with the given id. The caller holds the
// plan lock and allocated the id under it.
func (o *Operations) createRecipeLocked(id string, spec RecipeSpec) (*Recipe, error) {
	now := Now()
	r := &Recipe{
		ID:        id,
		Title:     spec.Title,
		Slug:      Slugify(spec.Title),
		Created:   now,
		Updated:   now,
		Lifecycle: LifecycleDraft,
		Intent:    Intent{Outcome: spec.Outcome, Problem: spec.Problem},
		Scope:     Scope{InScope: spec.InScope, OutOfScope: spec.OutOfScope},
	}
	for _, c := range spec.Criteria {
		r.Acceptance = append(r.Acceptance, Criterion{
			ID:           nextID("AC-", criterionIDs(r), 2),
			Description:  c.Description,
			Verification: c.Verification,
		})
	}
	for _, name := range spec.WorkPackages {
		r.WorkPackages = append(r.WorkPackages, WorkPackage{
			ID: fmt.Sprintf("WP%d", len(r.WorkPackages)+1), Name: name,
		})
	}
	for _, s := range spec.Steps {
		step := Step{
			ID:       nextID("G", stepIDs(r), 0),
			Action:   s.Action,
			DoneWhen: s.DoneWhen,
			Order:    len(r.Guidance) + 1,
			Status:   "pending",
			Covers:   s.Covers,
		}
		if s.WorkPackage != "" {
			step.WorkPackage = s.WorkPackage
		} else if len(r.WorkPackages) > 0 {
			step.WorkPackage = r.WorkPackages[0].ID
		}
		r.Guidance = append(r.Guidance, step)
	}
	if err := o.store.SaveRecipe(r); err != nil {
		return nil, err
	}
	return r, nil
}

// UpdateRecipe applies a fragment.
//
// The frozen sections — intent, problem, scope, acceptance — are closed once the recipe
// is signed off. Changing them after approval is an abort or a follow-up, not an edit.
func (o *Operations) UpdateRecipe(id string, f Fragment) (*Recipe, error) {
	return o.mutate(id, KindCreated, func(r *Recipe) error {
		if r.Approval != nil {
			touchesFrozen := f.Outcome != nil || f.Problem != nil ||
				f.InScope != nil || f.OutOfScope != nil
			if touchesFrozen {
				return &RefusalError{Reason: "intent, problem and scope are frozen once the recipe is signed off — abort to draft, or open a follow-up"}
			}
		}
		if f.Title != nil {
			r.Title = *f.Title
			r.Slug = Slugify(*f.Title)
		}
		if f.Outcome != nil {
			r.Intent.Outcome = *f.Outcome
		}
		if f.Problem != nil {
			r.Intent.Problem = *f.Problem
		}
		if f.InScope != nil {
			r.Scope.InScope = *f.InScope
		}
		if f.OutOfScope != nil {
			r.Scope.OutOfScope = *f.OutOfScope
		}
		if f.DependsOn != nil {
			r.DependsOn = *f.DependsOn
		}
		return nil
	})
}

// AddCriterion appends an acceptance criterion.
//
// Criteria are frozen at sign-off, so this is refused once the recipe is approved.
func (o *Operations) AddCriterion(id, description, verification string) (*Recipe, error) {
	return o.mutate(id, KindCriterion, func(r *Recipe) error {
		if r.Approval != nil {
			return &RefusalError{Reason: "acceptance criteria are frozen once the recipe is signed off"}
		}
		if strings.TrimSpace(description) == "" {
			return errors.New("criterion needs a description")
		}
		r.Acceptance = append(r.Acceptance, Criterion{
			ID:           nextID("AC-", criterionIDs(r), 2),
			Description:  description,
			Verification: verification,
		})
		return nil
	})
}

// AddWorkPackage appends a work package.
func (o *Operations) AddWorkPackage(id, name, goal, rollback string) (*Recipe, error) {
	return o.mutate(id, KindCreated, func(r *Recipe) error {
		r.WorkPackages = append(r.WorkPackages, WorkPackage{
			ID:       fmt.Sprintf("WP%d", len(r.WorkPackages)+1),
			Name:     name,
			Goal:     goal,
			Rollback: rollback,
		})
		return nil
	})
}

// AddStep appends a guidance step. Guidance is revisable at any time.
func (o *Operations) AddStep(id string, s StepSpec) (*Recipe, error) {
	return o.mutate(id, KindStep, func(r *Recipe) error {
		step := Step{
			ID:       nextID("G", stepIDs(r), 0),
			Action:   s.Action,
			DoneWhen: s.DoneWhen,
			Order:    len(r.Guidance) + 1,
			Status:   "pending",
			Covers:   s.Covers,
		}
		if len(r.WorkPackages) > 0 {
			step.WorkPackage = r.WorkPackages[0].ID
		}
		r.Guidance = append(r.Guidance, step)
		return nil
	})
}

// SetStepStatus records progress on a step.
func (o *Operations) SetStepStatus(id, stepID, status string) (*Recipe, error) {
	switch status {
	case "pending", "in-progress", "done", "skipped", "failed":
	default:
		return nil, fmt.Errorf("unknown step status %q", status)
	}
	return o.mutate(id, KindStep, func(r *Recipe) error {
		for i := range r.Guidance {
			if r.Guidance[i].ID == stepID {
				r.Guidance[i].Status = status
				return nil
			}
		}
		return fmt.Errorf("no step %q", stepID)
	})
}

// VerifyCriterion runs the criterion's verification and records the evidence.
//
// The tool runs the command. The evidence is append-only. A result that matched
// nothing is rejected; one the tool could not attribute is unproven. See ADR-0006.
//
// The command runs **outside** the plan lock: a verification can run for minutes, and
// the lock goes stale long before that. Holding it let a second writer break the lock
// of a live process, and this save then clobbered the other writer's change. The lock
// is taken only to append the evidence, against the recipe as it is *now*.
func (o *Operations) VerifyCriterion(id, acID string) (*Criterion, error) {
	verification, err := o.verificationFor(id, acID)
	if err != nil {
		return nil, err
	}
	result, err := o.verifier.Run(context.Background(), verification, o.context.Actor)
	if err != nil {
		return nil, err
	}
	var out Criterion
	_, err = o.mutate(id, KindEvidence, func(r *Recipe) error {
		for i := range r.Acceptance {
			if r.Acceptance[i].ID != acID {
				continue
			}
			r.Acceptance[i].Evidence = append(r.Acceptance[i].Evidence, result.Evidence)
			out = r.Acceptance[i]
			return nil
		}
		return fmt.Errorf("no criterion %q", acID)
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// verificationFor reads a criterion's verification command.
func (o *Operations) verificationFor(id, acID string) (string, error) {
	r, err := o.store.LoadRecipe(id)
	if err != nil {
		return "", err
	}
	for _, c := range r.Acceptance {
		if c.ID == acID {
			return c.Verification, nil
		}
	}
	return "", fmt.Errorf("no criterion %q", acID)
}

// SignOffCriterion records a human sign-off for an unproven criterion.
//
// The sign-off is evidence without a command: the human vouches for a procedure
// the tool could not check. It is marked as a sign-off so the console can say who
// vouched — and it claims no match. The zero-match rule stands for machine evidence;
// a sign-off is the one explicit exception, and it has its own state.
//
// A sign-off is a human act, so it takes the same deny-list the approval takes: the
// chef signs a criterion off, not an agent. Otherwise the exception to the zero-match
// rule would be available to the very party the rule exists to constrain, and a
// criterion could be `passed` with nothing behind it.
//
// The name that signs off is a name the tool accepts: a retired one is refused.
func (o *Operations) SignOffCriterion(id, acID, by, note string) (*Recipe, error) {
	if strings.TrimSpace(by) == "" {
		return nil, errors.New("a sign-off needs an approver")
	}
	if err := checkName(by); err != nil {
		return nil, err
	}
	if IsAgentRole(o.context.Actor) {
		return nil, fmt.Errorf("%s is an agent role: only the chef signs a criterion off", o.context.Actor)
	}
	return o.mutate(id, KindEvidence, func(r *Recipe) error {
		for i := range r.Acceptance {
			c := &r.Acceptance[i]
			if c.ID != acID {
				continue
			}
			c.Evidence = append(c.Evidence, Evidence{
				Command:   "manual: " + note,
				Match:     MatchSignedOff,
				SignedOff: true,
				At:        Now(),
				By:        by,
			})
			return nil
		}
		return fmt.Errorf("no criterion %q", acID)
	})
}

// DeferCriterion defers a criterion, with the reason and where it was re-filed.
//
// A deferred criterion does not keep the recipe open — CPD's rule is that a recipe is
// done when its in-scope criteria pass, and a deferred one is re-filed as a follow-up.
func (o *Operations) DeferCriterion(id, acID, reason, followUp string) (*Recipe, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, errors.New("a deferral needs a reason")
	}
	return o.mutate(id, KindCriterion, func(r *Recipe) error {
		for i := range r.Acceptance {
			if r.Acceptance[i].ID == acID {
				r.Acceptance[i].Deferral = &Deferral{Reason: reason, FollowUp: followUp}
				return nil
			}
		}
		return fmt.Errorf("no criterion %q", acID)
	})
}

// DismissFinding closes a finding without fixing it, with the reason.
//
// The chef's move, as distinct from the agent resolving one.
//
// A note's outcome is written once. A note that is already closed keeps the record of
// what was done to it: a dismissal must not overwrite the fixer's explanation of
// what they fixed.
func (o *Operations) DismissFinding(id, findingID, reason string) (*Recipe, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, errors.New("a dismissal needs a reason")
	}
	rule, _ := NoteOutcomeOf(NoteMoveDismiss)
	return o.mutate(id, KindFinding, func(r *Recipe) error {
		for i := range r.Review.Findings {
			f := &r.Review.Findings[i]
			if f.ID != findingID {
				continue
			}
			if !rule.AppliesTo(f.Status) {
				return &RefusalError{Reason: fmt.Sprintf("note %s is %s — a dismissal applies to an open note", findingID, f.Status)}
			}
			f.Status = rule.To
			f.Resolved = Now()
			f.Resolution = reason
			return nil
		}
		return fmt.Errorf("no finding %q", findingID)
	})
}

// ReopenFinding sends a finding back into the queue, with the objection recorded.
//
// This is the reviewer's move when a fix does not hold. The finding's outcome is
// cleared — the substance of the previous attempt lives in the comments, which are
// append-only — and the objection is appended as a comment, so the agent that tries
// next can read why it was sent back.
//
// The cleared outcome goes into that comment. A comment that carried only the
// objection would drop the fixer's record of what they did, and the note's outcome
// is written once.
func (o *Operations) ReopenFinding(id, findingID, reason string) (*Recipe, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, errors.New("a reopen needs a reason")
	}
	rule, _ := NoteOutcomeOf(NoteMoveReopen)
	return o.mutate(id, KindFinding, func(r *Recipe) error {
		for i := range r.Review.Findings {
			f := &r.Review.Findings[i]
			if f.ID != findingID {
				continue
			}
			if !rule.AppliesTo(f.Status) {
				return &RefusalError{Reason: fmt.Sprintf("note %s is %s — a reopen applies to a closed note", findingID, f.Status)}
			}
			body := reason
			if f.Resolution != "" {
				body = fmt.Sprintf("%s (it was closed as %q)", reason, f.Resolution)
			}
			f.Status = rule.To
			f.Resolved = ""
			f.Resolution = ""
			r.Review.Comments = append(r.Review.Comments, Comment{
				ID:      nextID("C-", commentIDs(r), 2),
				Author:  o.context.Actor,
				Body:    body,
				At:      Now(),
				Finding: findingID,
			})
			return nil
		}
		return fmt.Errorf("no finding %q", findingID)
	})
}

// RecordDeviation notes that a step was done differently than described.
//
// The deviation is written on the step as well as on the audit trail: the console
// renders "deviated: …" under the step, and a line nothing can write is a promise
// the code does not keep.
func (o *Operations) RecordDeviation(id, stepID, note string) (*Recipe, error) {
	if strings.TrimSpace(note) == "" {
		return nil, errors.New("a deviation needs a note")
	}
	return o.mutate(id, KindDeviation, func(r *Recipe) error {
		for i := range r.Guidance {
			if r.Guidance[i].ID == stepID {
				r.Guidance[i].Deviation = note
				r.Audit.Deviations = append(r.Audit.Deviations, Deviation{
					At: Now(), Step: stepID, Note: note, By: o.context.Actor,
				})
				return nil
			}
		}
		return fmt.Errorf("no step %q", stepID)
	})
}

// DocumentRecipe writes the closing documentation — the audit-documentation gate.
func (o *Operations) DocumentRecipe(id, text string) (*Recipe, error) {
	return o.mutate(id, KindDecision, func(r *Recipe) error {
		r.Audit.Documentation = text
		return nil
	})
}

// AddFinding files a tasting note.
//
// The phase and the severity are closed sets, checked here so every surface gets the
// check. An unknown value is an error, never a value that quietly means something
// else: the gates read the phase, so a typo would leave the recipe stuck between two
// gates that disagree about what "external" means.
func (o *Operations) AddFinding(id string, f Finding) (*Recipe, error) {
	if err := validFindingPhase(f.Phase); err != nil {
		return nil, err
	}
	if err := validFindingSeverity(f.Severity); err != nil {
		return nil, err
	}
	return o.mutate(id, KindFinding, func(r *Recipe) error {
		if f.ID == "" {
			f.ID = nextID("F-", findingIDs(r), 2)
		}
		if f.Status == "" {
			f.Status = "open"
		}
		if f.Created == "" {
			f.Created = Now()
		}
		r.Review.Findings = append(r.Review.Findings, f)
		return nil
	})
}

// validFindingPhase refuses a phase outside the closed set.
func validFindingPhase(phase string) error {
	if phase == "" {
		return fmt.Errorf("a finding needs a phase (%s)", strings.Join(FindingPhases, ", "))
	}
	for _, p := range FindingPhases {
		if p == phase {
			return nil
		}
	}
	return fmt.Errorf("unknown finding phase %q (want %s)", phase, strings.Join(FindingPhases, ", "))
}

// validFindingSeverity refuses a severity outside the closed set.
func validFindingSeverity(severity string) error {
	if severity == "" {
		return fmt.Errorf("a finding needs a severity (%s)", strings.Join(FindingSeverityLevels, ", "))
	}
	for _, s := range FindingSeverityLevels {
		if s == severity {
			return nil
		}
	}
	return fmt.Errorf("unknown finding severity %q (want %s)", severity, strings.Join(FindingSeverityLevels, ", "))
}

// ResolveFinding closes a tasting note.
//
// The fixer's move: the agent did the work, and says what it did. A note's outcome is
// written once — a resolution must not overwrite the chef's reason for dismissing it.
func (o *Operations) ResolveFinding(id, findingID, resolution string) (*Recipe, error) {
	rule, _ := NoteOutcomeOf(NoteMoveResolve)
	return o.mutate(id, KindFinding, func(r *Recipe) error {
		for i := range r.Review.Findings {
			f := &r.Review.Findings[i]
			if f.ID != findingID {
				continue
			}
			if !rule.AppliesTo(f.Status) {
				return &RefusalError{Reason: fmt.Sprintf("note %s is %s — a resolution applies to an open note", findingID, f.Status)}
			}
			f.Status = rule.To
			f.Resolved = Now()
			f.Resolution = resolution
			return nil
		}
		return fmt.Errorf("no finding %q", findingID)
	})
}

// RecordReviewRound records a review verdict.
//
// Without this the `review.verdict_ok` gate could never be satisfied — the gate
// existed but no verb could satisfy it.
//
// The reviewer's name is a name the tool accepts: a retired one is refused, so an
// executor cannot hide behind the name it used before the rename.
func (o *Operations) RecordReviewRound(id, reviewer, verdict, notes string) (*Recipe, error) {
	switch verdict {
	case "ok", "changes-needed":
	default:
		return nil, fmt.Errorf("unknown review verdict %q (want ok or changes-needed)", verdict)
	}
	if err := checkName(reviewer); err != nil {
		return nil, err
	}
	if strings.TrimSpace(reviewer) == "" {
		return nil, errors.New("a review round needs a reviewer")
	}
	return o.mutate(id, KindComment, func(r *Recipe) error {
		r.Review.Rounds = append(r.Review.Rounds, Round{
			N:        len(r.Review.Rounds) + 1,
			Reviewer: reviewer,
			Verdict:  verdict,
			At:       Now(),
			Notes:    notes,
		})
		return nil
	})
}

// AddComment adds a review comment.
//
// A comment belongs to a tasting note, or to the recipe itself. A note that is not
// there is a bad request, not a refusal: the console groups the comments by the note
// they belong to, so a comment on a note that does not exist would be in no thread at
// all — in the artifact, in the journal, and on no screen.
//
// The author's name is a name the tool accepts: a retired one is refused.
func (o *Operations) AddComment(id, author, body, findingID string) (*Recipe, error) {
	if err := checkName(author); err != nil {
		return nil, err
	}
	if strings.TrimSpace(body) == "" {
		return nil, errors.New("a comment needs a body")
	}
	return o.mutate(id, KindComment, func(r *Recipe) error {
		if findingID != "" {
			found := false
			for _, f := range r.Review.Findings {
				if f.ID == findingID {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("no finding %q", findingID)
			}
		}
		r.Review.Comments = append(r.Review.Comments, Comment{
			ID:      nextID("C-", commentIDs(r), 2),
			Author:  author,
			Body:    body,
			At:      Now(),
			Finding: findingID,
		})
		return nil
	})
}

// ---------------------------------------------------------------------------
// Writing — decisions
// ---------------------------------------------------------------------------

// Advance applies a lifecycle event and records it.
func (o *Operations) Advance(id string, ev Event) (*Recipe, error) {
	return o.AdvanceWith(id, ev, TransitionContext{})
}

// AdvanceWith applies a lifecycle event with a transition context — a reason, an
// approver — and records it.
func (o *Operations) AdvanceWith(id string, ev Event, ctx TransitionContext) (*Recipe, error) {
	if ctx.Approver == "" {
		ctx.Approver = o.context.Actor
	}
	if ctx.Actor == "" {
		ctx.Actor = o.context.Actor
	}
	if ctx.Session == "" {
		ctx.Session = o.context.Session
	}
	if ctx.Role == "" {
		ctx.Role = o.context.Role
	}
	if ev == EventStart {
		// Starting work is the claim. One path owns the state change, so the
		// dependency gate runs on every road to `executing`, and every recipe on
		// the line has a claimant the review gate can be separate from.
		return o.ClaimRecipe(id, ctx.Actor)
	}
	return o.advance(id, ev, ctx.Approver, ctx.Reason)
}

// ApproveRecipe signs a recipe off. The approver is recorded, and the frozen sections
// are sealed.
//
// The approver's name is a name the tool accepts: a retired one is refused.
func (o *Operations) ApproveRecipe(id, by string) (*Recipe, error) {
	if strings.TrimSpace(by) == "" {
		return nil, errors.New("a sign-off needs an approver")
	}
	if err := checkName(by); err != nil {
		return nil, err
	}
	return o.advance(id, EventApprove, by, "")
}

// RequestChanges sends a signed-off-ready recipe back to shaping, with the objection
// recorded. The gates must pass again on the way up.
func (o *Operations) RequestChanges(id, by, notes string) (*Recipe, error) {
	if strings.TrimSpace(notes) == "" {
		return nil, errors.New("requesting changes needs a reason")
	}
	return o.advance(id, EventRequestChanges, by, notes)
}

// BlockRecipe marks the recipe blocked.
func (o *Operations) BlockRecipe(id, reason string) (*Recipe, error) {
	return o.mutate(id, KindOverlay, func(r *Recipe) error {
		return Block(r, reason, o.context.Actor, time.Time{})
	})
}

// UnblockRecipe clears the blocked overlay.
func (o *Operations) UnblockRecipe(id string) (*Recipe, error) {
	return o.mutate(id, KindOverlay, func(r *Recipe) error {
		Unblock(r)
		return nil
	})
}

// DeferRecipe parks the recipe.
func (o *Operations) DeferRecipe(id, reason string) (*Recipe, error) {
	return o.mutate(id, KindOverlay, func(r *Recipe) error {
		return Defer(r, reason, o.context.Actor, time.Time{})
	})
}

// ResumeRecipe clears the deferred overlay.
func (o *Operations) ResumeRecipe(id string) (*Recipe, error) {
	return o.mutate(id, KindOverlay, func(r *Recipe) error {
		Resume(r)
		return nil
	})
}

// CancelRecipe cancels the recipe. Terminal.
func (o *Operations) CancelRecipe(id, reason string) (*Recipe, error) {
	return o.mutate(id, KindOverlay, func(r *Recipe) error {
		return Cancel(r, reason, o.context.Actor, time.Time{})
	})
}

// ClaimRecipe puts the recipe on the line under an agent.
//
// The claim is a lifecycle transition — `EventStart` — so it runs the same gates as
// every other road to `executing`, the dependency rule included, and the journal shows
// it as the transition it is.
//
// A claim needs an agent. An anonymous claim made the independence gate vacuous: the
// executor could sign off its own work by omitting a flag.
//
// The claimant's name is a name the tool accepts: a retired one is refused, so an
// executor cannot claim under the name it used before the rename and be reviewed
// under the current one.
func (o *Operations) ClaimRecipe(id, agent string) (*Recipe, error) {
	if strings.TrimSpace(agent) == "" {
		return nil, errors.New("a claim needs an agent")
	}
	if err := checkName(agent); err != nil {
		return nil, err
	}
	unlock, err := o.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	r, err := o.store.LoadRecipe(id)
	if err != nil {
		return nil, err
	}
	ctx := TransitionContext{
		Actor:   o.context.Actor,
		Session: o.context.Session,
		Role:    o.context.Role,
		Closed:  o.closedLookup,
	}
	from := r.Lifecycle
	if err := Transition(r, EventStart, ctx); err != nil {
		// The claim was refused — the dependency rule, most often. Journal it.
		o.journalRefusalLocked(r.ID, EventStart, err)
		return nil, err
	}
	r.Execution.ClaimedBy = agent
	r.Execution.ClaimedAt = Now()
	r.Updated = Now()
	if err := o.store.SaveRecipe(r); err != nil {
		return nil, err
	}
	if err := o.journalLocked(KindLifecycle, r.ID, map[string]any{
		"from": string(from), "to": string(r.Lifecycle), "agent": agent,
	}); err != nil {
		return nil, err
	}
	return r, nil
}

// ---------------------------------------------------------------------------
// Internals
// ---------------------------------------------------------------------------

// lock takes the plan lock for a write, and refuses an actor the tool does not
// accept.
//
// Every write goes through here, so the rule holds everywhere: a name the tool does
// not accept cannot be the actor that wrote the plan, and so it cannot appear in a
// live record.
//
// The **role** is checked with the actor: it lands in the journal's event, and a new
// event is a live record. The check is here — before anything is written — so a
// retired role never leaves a half-written plan behind.
func (o *Operations) lock() (func(), error) {
	if err := checkName(o.context.Actor); err != nil {
		return nil, err
	}
	if err := checkName(o.context.Role); err != nil {
		return nil, err
	}
	return o.store.lock()
}

// mutate loads a recipe, applies fn, saves it, and journals the change.
//
// The whole sequence happens under the plan's advisory lock, so two processes cannot
// interleave.
func (o *Operations) mutate(id, kind string, fn func(*Recipe) error) (*Recipe, error) {
	unlock, err := o.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()
	return o.mutateLocked(id, kind, fn)
}

// mutateLocked does mutate's work. The caller holds the plan lock.
func (o *Operations) mutateLocked(id, kind string, fn func(*Recipe) error) (*Recipe, error) {
	r, err := o.store.LoadRecipe(id)
	if err != nil {
		return nil, err
	}
	if err := fn(r); err != nil {
		// The mutation was refused by the plan — a frozen section, say. Journal
		// it: the attempt is evidence.
		o.journalRefusalLocked(r.ID, "", err)
		return nil, err
	}
	r.Updated = Now()
	if err := o.store.SaveRecipe(r); err != nil {
		return nil, err
	}
	if err := o.journalLocked(kind, r.ID, nil); err != nil {
		return nil, err
	}
	return r, nil
}

// closedLookup reports whether a recipe is closed.
//
// The gate table is package-level and cannot close over the store, so the operations
// layer hands the gate this closure. It is the only way a gate learns about the plan.
func (o *Operations) closedLookup(id string) (bool, error) {
	r, err := o.store.LoadRecipe(id)
	if err != nil {
		return false, err
	}
	return r.Lifecycle == LifecycleClosed, nil
}

// advance applies a lifecycle event under the lock.
//
// The **approver** is the name the event records — the approval's `by`, the
// transition's `actor` — so it is checked here: this is the one place every
// transition's approver passes through, whether the caller is `ApproveRecipe`,
// `RequestChanges` or `AdvanceWith`.
func (o *Operations) advance(id string, ev Event, approver, reason string) (*Recipe, error) {
	if err := checkName(approver); err != nil {
		return nil, err
	}
	unlock, err := o.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	r, err := o.store.LoadRecipe(id)
	if err != nil {
		return nil, err
	}
	ctx := TransitionContext{
		Approver: approver,
		Reason:   reason,
		Actor:    o.context.Actor,
		Session:  o.context.Session,
		Role:     o.context.Role,
		Closed:   o.closedLookup,
	}
	if ctx.Actor == "" {
		ctx.Actor = o.context.Actor
	}
	from := r.Lifecycle
	if err := Transition(r, ev, ctx); err != nil {
		// The transition was refused. The refusal is journaled: the journal is the
		// only event log, and an attempt that was turned away is evidence.
		o.journalRefusalLocked(r.ID, ev, err)
		return nil, err
	}
	r.Updated = Now()
	if err := o.store.SaveRecipe(r); err != nil {
		return nil, err
	}
	data := map[string]any{"from": string(from), "to": string(r.Lifecycle)}
	if reason != "" {
		data["reason"] = reason
	}
	if approver != "" {
		data["approver"] = approver
	}
	if err := o.journalLocked(KindLifecycle, r.ID, data); err != nil {
		return nil, err
	}
	return r, nil
}

// journal appends an entry, taking the lock.
func (o *Operations) journal(kind, recipeID string, data any) error {
	unlock, err := o.lock()
	if err != nil {
		return err
	}
	defer unlock()
	return o.journalLocked(kind, recipeID, data)
}

// refuseLocked turns an action away and writes it down.
//
// The caller holds the lock. A refusal the plan produced is evidence: somebody asked,
// and the plan said no. The journal is the only event log, so an attempt that was
// turned away belongs in it — whether the refusal came from a gate, a transition, a
// frozen section or the state of the artifact itself.
func (o *Operations) refuseLocked(id, reason string) error {
	err := &RefusalError{Reason: reason}
	o.journalRefusalLocked(id, "", err)
	return err
}

// journalRefusalLocked records a refused action. The caller holds the lock.
//
// A refusal is information, not an error: the gate working as designed. But it is
// still an attempt — somebody asked, and the plan said no. The journal is the only
// event log, so an attempt that was turned away is written down.
//
// Which refusals count: **the ones the plan produced**. A gate that did not pass, a
// transition that is not legal from here, a frozen section that may not be edited —
// these are the plan's answers, and they are what the chef and the next agent need to
// see. An operation-level failure ("no such recipe", "a claim needs an agent", a
// missing argument) is a mistake in the request, not a decision about the plan, and
// journalling it would fill the log with noise about typos.
//
// The error is not swallowed: the caller still gets it, and the CLI still exits 2.
func (o *Operations) journalRefusalLocked(recipeID string, ev Event, err error) {
	var ge *GateError
	var te *TransitionError
	var re *RefusalError
	switch {
	case errors.As(err, &ge):
	case errors.As(err, &te):
	case errors.As(err, &re):
	default:
		// Not the plan's answer: a bad request, not a refusal.
		return
	}
	data := map[string]any{"reason": err.Error()}
	if ev != "" {
		data["event"] = string(ev)
	}
	var gates []string
	if errors.As(err, &ge) {
		for _, g := range ge.Gates {
			// The gate's name, without the parenthetical reason: the name is the
			// fact, the reason is prose.
			if i := strings.Index(g, " ("); i > 0 {
				g = g[:i]
			}
			gates = append(gates, g)
		}
	}
	if gates != nil {
		data["gates"] = gates
	}
	// A refusal that cannot be written must not mask the refusal itself: the
	// caller gets the error either way.
	_ = o.journalLocked(KindRefusal, recipeID, data)
}

// journalLocked appends an entry. The caller holds the lock.
func (o *Operations) journalLocked(kind, recipeID string, data any) error {
	entry := JournalEntry{
		Kind:    kind,
		Recipe:  recipeID,
		Actor:   o.context.Actor,
		Session: o.context.Session,
		Role:    o.context.Role,
	}
	if data != nil {
		raw, err := json.Marshal(data)
		if err != nil {
			return err
		}
		entry.Data = raw
	}
	_, err := o.store.appendJournalEventLocked(entry)
	return err
}
