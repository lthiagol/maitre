// Package core is the domain layer of maitre.
//
// It holds the plan's data model, its lifecycle and gates, and the store that reads and
// writes the plan directory. It has no UI imports and no CLI imports — every surface is
// a client of it. See docs/03-domain-model.md.
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Recipe is the durable unit of work: CPD's "crafted prompt", themed as a recipe card.
// See docs/03-domain-model.md § 2.
type Recipe struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Slug         string        `json:"slug"`
	Created      string        `json:"created"`
	Updated      string        `json:"updated"`
	Lifecycle    Lifecycle     `json:"lifecycle"`
	Intent       Intent        `json:"intent"`
	Scope        Scope         `json:"scope"`
	Acceptance   []Criterion   `json:"acceptance"`
	Guidance     []Step        `json:"guidance"`
	WorkPackages []WorkPackage `json:"work_packages"`
	Questions    []Question    `json:"questions"`
	DependsOn    []string      `json:"depends_on,omitempty"`
	Execution    Execution     `json:"execution"`
	Evidence     Evidence      `json:"evidence"`
	Review       Review        `json:"review"`
	Audit        Audit         `json:"audit"`
	Approval     *Approval     `json:"approval,omitempty"`
	Overlays     Overlays      `json:"overlays"`
}

// Intent is what we set out to make, and why. Frozen at approval.
type Intent struct {
	Outcome string `json:"outcome"`
	Problem string `json:"problem"`
}

// Scope is what goes in the pan and what does not. Frozen at approval.
type Scope struct {
	InScope    []string `json:"in_scope"`
	OutOfScope []string `json:"out_of_scope"`
}

// Criterion is one falsifiable acceptance criterion.
type Criterion struct {
	ID           string    `json:"id"`
	Description  string    `json:"description"`
	Verification string    `json:"verification"`
	Proofs       []Proof   `json:"proofs"`
	Deferral     *Deferral `json:"deferral,omitempty"`
}

// Status derives the criterion's state from its latest proof.
//
// The exit code dominates: a command that failed is a failure, whatever the count. A
// command that succeeded but matched nothing is *rejected* — a spec defect, not a test
// failure. A command whose match count the tool could not determine is *unproven*.
// Neither satisfies a gate.
//
// A deferral only shows on a criterion with no proof of its own: a criterion that
// already passed stays passed. A deferral re-files the work; it does not undo a proof.
//
// view
func (c Criterion) Status() CriterionStatus {
	if len(c.Proofs) == 0 {
		if c.Deferral != nil {
			return CriterionDeferred
		}
		return CriterionPending
	}
	last := c.Proofs[len(c.Proofs)-1]
	switch {
	case last.ExitCode != 0:
		return CriterionFailed
	case last.Match == MatchProven, last.Match == MatchSignedOff:
		// A machine proof that matched something, or a human who vouched for a
		// procedure the tool could not check.
		return CriterionPassed
	case last.Match == MatchRejected:
		return CriterionRejected
	default:
		return CriterionUnproven
	}
}

// Deferral records why a criterion was deferred and where it was re-filed.
type Deferral struct {
	Reason   string `json:"reason"`
	FollowUp string `json:"follow_up,omitempty"`
}

// Step is one guidance step: a proposed path, revisable, with deviations recorded.
type Step struct {
	ID          string   `json:"id"`
	WorkPackage string   `json:"work_package,omitempty"`
	Order       int      `json:"order"`
	Action      string   `json:"action"`
	Files       []string `json:"files,omitempty"`
	DependsOn   []string `json:"depends_on,omitempty"`
	Covers      []string `json:"covers,omitempty"`
	DoneWhen    string   `json:"done_when,omitempty"`
	Status      string   `json:"status"`
	// Deviation is set when the step was done differently than described, with why.
	// The console renders it, so an operation has to write it.
	Deviation string `json:"deviation,omitempty"`
}

// WorkPackage groups steps owned by one agent or person.
type WorkPackage struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Goal     string `json:"goal,omitempty"`
	Rollback string `json:"rollback,omitempty"`
}

// Execution tracks who claimed the recipe.
//
// The claim is the start: there is no second moment, so there is no second field.
// A branch is a version-control fact, and the plan is not a repository.
type Execution struct {
	ClaimedBy string `json:"claimed_by,omitempty"`
	ClaimedAt string `json:"claimed_at,omitempty"`
}

// Evidence is the recipe-level completion stamp. Per-criterion proofs live on the
// criterion.
type Evidence struct {
	CompletedAt string `json:"completed_at,omitempty"`
}

// Review holds findings, comments and review rounds.
type Review struct {
	Findings []Finding `json:"findings"`
	Comments []Comment `json:"comments"`
	Rounds   []Round   `json:"rounds"`
}

// Finding is a review finding, from the self-review or the external tasting.
type Finding struct {
	ID          string  `json:"id"`
	Phase       string  `json:"phase"`
	Severity    string  `json:"severity"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Anchor      *Anchor `json:"anchor,omitempty"`
	Status      string  `json:"status"`
	Created     string  `json:"created"`
	Resolved    string  `json:"resolved,omitempty"`
	Resolution  string  `json:"resolution,omitempty"`
}

// Finding phases. A closed set: the gates read the phase, so a typo would leave a
// recipe stuck between two gates that disagree about what "external" means.
const (
	FindingPhaseSelf     = "self"
	FindingPhaseExternal = "external"
)

// FindingSeverities. A closed set, for the same reason.
const (
	SeverityLow    = "low"
	SeverityMedium = "medium"
	SeverityHigh   = "high"
)

// FindingPhases lists every phase a finding may carry.
var FindingPhases = []string{FindingPhaseSelf, FindingPhaseExternal}

// FindingSeverityLevels lists every severity a finding may carry.
var FindingSeverityLevels = []string{SeverityLow, SeverityMedium, SeverityHigh}

// NoteMove is a writer's move on a tasting note.
type NoteMove string

// The moves on a tasting note.
const (
	// NoteMoveResolve is the fixer's move: the agent did the work.
	NoteMoveResolve NoteMove = "resolve"
	// NoteMoveDismiss is the chef's move: the note is not worth fixing.
	NoteMoveDismiss NoteMove = "dismiss"
	// NoteMoveReopen is the reviewer's move: the fix does not hold.
	NoteMoveReopen NoteMove = "reopen"
)

// NoteOutcome is one writer's move on a tasting note.
//
// A note's outcome is written once: the fixer closes an open note as *fixed*, the
// chef closes an open note as *dismissed*, and the reviewer sends a closed note back
// to *open*. A second writer is refused — otherwise the chef's dismissal erases the
// fixer's record of what was done.
//
// The table is the rule and the operations read it, so the code cannot drift from it.
// The version fingerprint hashes it too: a plan says which rules it was written under,
// and a tool that does not know the difference refuses to write.
type NoteOutcome struct {
	// Move is the writer's move.
	Move NoteMove
	// From is the status the note must carry for the move to apply. A move that
	// applies to any closed note lists the closed statuses.
	From []string
	// To is the status the note carries after the move.
	To string
}

// NoteOutcomes are the moves on a tasting note, one per writer.
var NoteOutcomes = []NoteOutcome{
	{Move: NoteMoveResolve, From: []string{"open"}, To: "fixed"},
	{Move: NoteMoveDismiss, From: []string{"open"}, To: "dismissed"},
	{Move: NoteMoveReopen, From: []string{"fixed", "dismissed"}, To: "open"},
}

// NoteOutcomeOf is the rule for a writer's move.
func NoteOutcomeOf(move NoteMove) (NoteOutcome, bool) {
	for _, o := range NoteOutcomes {
		if o.Move == move {
			return o, true
		}
	}
	return NoteOutcome{}, false
}

// AppliesTo reports whether a note in the given status may take the move.
func (o NoteOutcome) AppliesTo(status string) bool {
	for _, from := range o.From {
		if from == status {
			return true
		}
	}
	return false
}

// Anchor points at the code a finding is about.
type Anchor struct {
	File string `json:"file"`
	Line int    `json:"line,omitempty"`
	Side string `json:"side,omitempty"`
}

// Comment is a threaded review comment.
type Comment struct {
	ID      string `json:"id"`
	Author  string `json:"author"`
	Body    string `json:"body"`
	At      string `json:"at"`
	Finding string `json:"finding,omitempty"`
}

// Round is one review round's verdict.
type Round struct {
	N        int    `json:"n"`
	Reviewer string `json:"reviewer"`
	Verdict  string `json:"verdict"`
	At       string `json:"at"`
	Notes    string `json:"notes,omitempty"`
}

// Audit holds the recipe's authored record.
//
// Events are **not** here — they live in the plan's journal, which is the only event log.
// The recipe's audit holds what the recipe itself authored: deviations and decisions.
type Audit struct {
	Deviations    []Deviation `json:"deviations"`
	Decisions     []Decision  `json:"decisions"`
	Documentation string      `json:"documentation,omitempty"`
}

// Question is an open question raised while shaping the recipe.
type Question struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Status   string `json:"status"`
	Answer   string `json:"answer,omitempty"`
}

// Deviation records that a step was done differently than described.
type Deviation struct {
	At   string `json:"at"`
	Step string `json:"step,omitempty"`
	Note string `json:"note"`
	By   string `json:"by,omitempty"`
}

// Decision is a project-level decision recorded against the plan.
type Decision struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
	Context string `json:"context,omitempty"`
	At      string `json:"at"`
}

// Approval is the human sign-off. Written once; the fingerprint freezes the sections
// that may not change afterwards.
type Approval struct {
	By          string `json:"by"`
	At          string `json:"at"`
	Fingerprint string `json:"fingerprint"`
}

// Overlays ride alongside the lifecycle.
type Overlays struct {
	Blocked   *Overlay `json:"blocked,omitempty"`
	Deferred  *Overlay `json:"deferred,omitempty"`
	Cancelled *Overlay `json:"cancelled,omitempty"`
}

// Overlay is a reason-stamped overlay state.
type Overlay struct {
	Reason string `json:"reason"`
	By     string `json:"by,omitempty"`
	At     string `json:"at"`
}

// idsOf collects the id of every item in a list.
//
// The ids are allocated from a scan of what exists, never from a count: v1's
// `len()+1` produced [AC-01, AC-03, AC-03] after a removal, and the plan directory
// is hand-editable.
func idsOf[T any](items []T, id func(T) string) []string {
	out := make([]string, 0, len(items))
	for _, it := range items {
		out = append(out, id(it))
	}
	return out
}

// criterionIDs, stepIDs, commentIDs and findingIDs are the ids a recipe already
// holds, so a new one can be allocated from a scan.
func criterionIDs(r *Recipe) []string {
	return idsOf(r.Acceptance, func(c Criterion) string { return c.ID })
}

func stepIDs(r *Recipe) []string {
	return idsOf(r.Guidance, func(s Step) string { return s.ID })
}

func commentIDs(r *Recipe) []string {
	return idsOf(r.Review.Comments, func(c Comment) string { return c.ID })
}

func findingIDs(r *Recipe) []string {
	return idsOf(r.Review.Findings, func(f Finding) string { return f.ID })
}

// Terminal reports whether the recipe has reached a terminal state.
func (r Recipe) Terminal() bool {
	return r.Lifecycle == LifecycleClosed || r.Overlays.Cancelled != nil
}

// normaliseLists replaces every nil list with an empty one.
//
// A list is a list: an empty list is `[]`, never `null`. The plan is read by agents,
// and the standard readers — `jq`, a JSON parser in any language — break on `null`.
// The store calls this when it reads a recipe, so the contract holds everywhere the
// plan is read, and the file that is written back says the same thing.
func (r *Recipe) normaliseLists() {
	if r.Scope.InScope == nil {
		r.Scope.InScope = []string{}
	}
	if r.Scope.OutOfScope == nil {
		r.Scope.OutOfScope = []string{}
	}
	if r.Acceptance == nil {
		r.Acceptance = []Criterion{}
	}
	if r.Guidance == nil {
		r.Guidance = []Step{}
	}
	if r.WorkPackages == nil {
		r.WorkPackages = []WorkPackage{}
	}
	if r.Questions == nil {
		r.Questions = []Question{}
	}
	if r.DependsOn == nil {
		r.DependsOn = []string{}
	}
	if r.Review.Findings == nil {
		r.Review.Findings = []Finding{}
	}
	if r.Review.Comments == nil {
		r.Review.Comments = []Comment{}
	}
	if r.Review.Rounds == nil {
		r.Review.Rounds = []Round{}
	}
	if r.Audit.Deviations == nil {
		r.Audit.Deviations = []Deviation{}
	}
	if r.Audit.Decisions == nil {
		r.Audit.Decisions = []Decision{}
	}
	for i := range r.Acceptance {
		if r.Acceptance[i].Proofs == nil {
			r.Acceptance[i].Proofs = []Proof{}
		}
	}
	for i := range r.Guidance {
		if r.Guidance[i].Files == nil {
			r.Guidance[i].Files = []string{}
		}
		if r.Guidance[i].DependsOn == nil {
			r.Guidance[i].DependsOn = []string{}
		}
		if r.Guidance[i].Covers == nil {
			r.Guidance[i].Covers = []string{}
		}
	}
}

// Fingerprint computes the digest of the sections frozen at approval: intent, scope and
// acceptance criteria. A change after approval changes the fingerprint, which is how
// the immutability gate detects a material edit (constraint C12).
//
// The digest is a pure function of the content: the copy is normalised before it is
// hashed, so `null` and `[]` — the same list, written by different builds — hash the
// same way. A digest that depends on the load path is not a digest of the recipe.
func (r Recipe) Fingerprint() (string, error) {
	canonical := r
	canonical.normaliseLists()
	return canonical.fingerprint()
}

// fingerprint hashes the recipe as it stands, lists and all.
//
// This is the encoding a build before the canonical one used: the lists hash as the
// file happened to hold them. The store compares a stored digest against this to tell
// "written by an older build" from "tampered with".
func (r Recipe) fingerprint() (string, error) {
	payload := struct {
		Intent     Intent      `json:"intent"`
		Scope      Scope       `json:"scope"`
		Acceptance []Criterion `json:"acceptance"`
	}{
		Intent:     r.Intent,
		Scope:      r.Scope,
		Acceptance: r.Acceptance,
	}
	// Only the frozen fields participate. Proofs are excluded because they grow after
	// approval; only the criterion's declared shape is frozen.
	payload.Acceptance = frozenCriteria(r.Acceptance)

	raw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("fingerprint: %w", err)
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

// restampLegacyApproval re-stamps an approval written by an older build.
//
// The fingerprint used to hash the lists as the file happened to hold them, so a plan
// approved before the normalisation carries a digest the canonical one does not match.
// The content did not change; only the encoding did — and a gate that reads that as
// tampering would refuse every transition of a recipe nobody touched.
//
// The older digest is recognised among the candidate encodings, not by recomputing
// the file's current form: an intermediate build may have rewritten the lists since,
// and the file's form is exactly what changed.
//
// The recipe must not be normalised yet: the candidates are read from the file's own
// shape.
func (r *Recipe) restampLegacyApproval() {
	if r.Approval == nil {
		return
	}
	canonical, err := r.Fingerprint()
	if err != nil || canonical == r.Approval.Fingerprint {
		return
	}
	for _, candidate := range r.legacyFingerprints() {
		if candidate == r.Approval.Fingerprint {
			r.Approval.Fingerprint = canonical
			return
		}
	}
	// The content changed, not the encoding. The gate will say so.
}

// legacyFingerprints are the digests an older build could have written.
//
// An older build hashed the file's own form, and in that form each **empty** list was
// either `null` or `[]`. The frozen payload holds three lists — the scope's two and
// the criteria — and a list that is not empty hashes the same either way, so only the
// empty ones vary. The candidates are the product of those choices: at most eight,
// and four for a recipe whose scope is empty and whose criteria are not.
//
// The set is complete: the payload's only list-shaped parts are those three, and the
// only thing an older build could write differently is whether an empty one was
// `null` or `[]`. The canonical form — every list `[]` — is one of the candidates.
//
// A change of *content* changes the digest under every candidate, so a tampered
// recipe matches none of them.
func (r Recipe) legacyFingerprints() []string {
	stringForms := func(l []string) [][]string {
		if len(l) > 0 {
			return [][]string{l}
		}
		return [][]string{nil, {}}
	}
	criterionForms := func(l []Criterion) [][]Criterion {
		if len(l) > 0 {
			return [][]Criterion{l}
		}
		return [][]Criterion{nil, {}}
	}
	var out []string
	for _, inScope := range stringForms(r.Scope.InScope) {
		for _, outOfScope := range stringForms(r.Scope.OutOfScope) {
			for _, acceptance := range criterionForms(r.Acceptance) {
				candidate := r
				candidate.Scope.InScope = inScope
				candidate.Scope.OutOfScope = outOfScope
				candidate.Acceptance = acceptance
				sum, err := candidate.fingerprint()
				if err != nil {
					continue
				}
				out = append(out, sum)
			}
		}
	}
	return out
}

// frozenCriteria strips the parts of a criterion that legitimately change after
// approval, keeping the declared intent and its verification.
func frozenCriteria(in []Criterion) []Criterion {
	if in == nil {
		return nil
	}
	out := make([]Criterion, 0, len(in))
	for _, c := range in {
		out = append(out, Criterion{
			ID:           c.ID,
			Description:  c.Description,
			Verification: c.Verification,
		})
	}
	return out
}

// AgentRoles are the names that may not sign a recipe off.
//
// The chef is the human; everyone else in the brigade is an agent — the sous chef,
// the line cook, the taster and the expo. `agent` is the CLI's default actor, so it
// is on the list too: a caller who does not say who they are is not the chef.
//
// The list is a **deny** list, not an allow list. A human may approve under any name
// they like; what the list refuses is the identity performing the approval. See
// ADR-0006, "the permission deny-list".
var AgentRoles = []string{"agent", "sous-chef", "line-cook", "taster", "expo"}

// IsAgentRole reports whether the name is one of the agent roles.
//
// The comparison is case-insensitive and treats a space as a dash, so "Line Cook"
// and "line-cook" are the same name.
func IsAgentRole(name string) bool {
	n := strings.ToLower(strings.TrimSpace(name))
	n = strings.ReplaceAll(n, " ", "-")
	for _, role := range AgentRoles {
		if n == role {
			return true
		}
	}
	return false
}

// Slugify turns a title into a filesystem-safe slug.
func Slugify(title string) string {
	var b strings.Builder
	prevDash := false
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		default:
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// Now returns the timestamp format used across the plan.
func Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}
