package core

import (
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

// Lane is where a recipe sits in the human's workflow. It is derived, never stored, so
// the board can be reorganised without a data migration.
type Lane string

const (
	LaneGrooming  Lane = "grooming"
	LaneApproval  Lane = "approval"
	LaneExecution Lane = "execution"
	LaneReview    Lane = "review"
	LaneBlocked   Lane = "blocked"
	LaneDone      Lane = "done"
)

// AllLanes lists the lanes in the order the board shows them.
var AllLanes = []Lane{
	LaneGrooming, LaneApproval, LaneExecution, LaneReview, LaneBlocked, LaneDone,
}

// LaneNames lists every lane, in order. The closed set a caller may filter by.
func LaneNames() []string {
	out := make([]string, 0, len(AllLanes))
	for _, l := range AllLanes {
		out = append(out, string(l))
	}
	return out
}

// KnownLane reports whether the name is a lane.
func KnownLane(name string) bool {
	for _, l := range AllLanes {
		if string(l) == name {
			return true
		}
	}
	return false
}

// LaneLabel is the board's name for a lane.
func (l Lane) Label() string {
	switch l {
	case LaneGrooming:
		return "Grooming"
	case LaneApproval:
		return "Sign-off"
	case LaneExecution:
		return "On the line"
	case LaneReview:
		return "Tasting"
	case LaneBlocked:
		return "On hold"
	case LaneDone:
		return "Served"
	}
	return string(l)
}

// LaneOf derives the lane a recipe belongs to.
//
// The cancelled case comes first: a cancelled recipe is over, whatever else it is,
// and the lane and the triage must say the same thing about it.
// view
func LaneOf(r Recipe) Lane {
	switch {
	case Cancelled(r):
		return LaneDone
	case r.Overlays.Blocked != nil:
		return LaneBlocked
	case r.Lifecycle == LifecycleClosed:
		return LaneDone
	case r.Lifecycle == LifecycleDraft, r.Lifecycle == LifecycleGroomed:
		return LaneGrooming
	case r.Lifecycle == LifecycleReady:
		return LaneApproval
	case r.Lifecycle == LifecycleApproved, r.Lifecycle == LifecycleExecuting:
		return LaneExecution
	case r.Lifecycle == LifecycleExecuted, r.Lifecycle == LifecycleRemediation, r.Lifecycle == LifecycleReviewed:
		return LaneReview
	}
	return LaneGrooming
}

// Cancelled reports whether the recipe was cancelled.
//
// One fact, one answer: the lane and the triage both read this, and both ask it
// first. A cancelled recipe used to be filed under "On the line" by the lane and
// under "Served" by the triage — two answers to one question, and a cancelled
// recipe presented as live work.
func Cancelled(r Recipe) bool {
	return r.Overlays.Cancelled != nil
}

// OpenFindings returns how many of a recipe's findings are still open, of any phase.
//
// One source: the `review.clean` and `review.no_open_findings` gates and the chef's
// move all read this, so the console cannot offer a move the gates would refuse.
func OpenFindings(r *Recipe) int {
	return countOpenFindings(r, "")
}

// OpenFindingsOf lists the findings still open, in the order the artifact holds them.
//
// The count and the list are one fact: the gate that refuses the close, the chef's
// move and the console's cursor over the notes all read the same set, so the console
// cannot select a note the gates do not see.
func OpenFindingsOf(r *Recipe) []Finding {
	var out []Finding
	for _, f := range r.Review.Findings {
		if f.Status == "open" {
			out = append(out, f)
		}
	}
	return out
}

// ChangesRequested reports whether a recipe's latest review round asked for changes.
func ChangesRequested(r Recipe) bool {
	if len(r.Review.Rounds) == 0 {
		return false
	}
	return r.Review.Rounds[len(r.Review.Rounds)-1].Verdict == "changes-needed"
}

// ChefMove is the move the chef can make on a recipe right now.
//
// Derived, never stored: the console's footer and its primary action both read this,
// so the console cannot offer a move the gates would refuse.
type ChefMove string

const (
	MoveNone    ChefMove = ""
	MoveSignOff ChefMove = "sign off"
	MoveReview  ChefMove = "review"
	MoveServe   ChefMove = "serve"
	MoveTriage  ChefMove = "triage"
)

// AllChefMoves is the closed set of the chef's moves.
//
// The console's footer, its primary action and the version fingerprint all read it: a
// new move is a new rule, and a plan has to say which rules it was written under —
// otherwise an older console files a recipe under "Needs you" with nothing to press.
var AllChefMoves = []ChefMove{MoveNone, MoveSignOff, MoveReview, MoveServe, MoveTriage}

// ChefMoveOf derives the chef's move for a recipe.
//
// Only the chef's own moves are here. Claiming the work is an agent's move, and the
// console is not an agent — it must never record itself as one.
// view
func ChefMoveOf(r Recipe) ChefMove {
	// A recipe on hold, or cancelled, offers nothing. The console must not sign off
	// or serve something the chef has explicitly put aside.
	if r.Overlays.Blocked != nil || r.Overlays.Cancelled != nil {
		return MoveNone
	}
	switch r.Lifecycle {
	case LifecycleReady:
		return MoveSignOff
	case LifecycleExecuted, LifecycleRemediation, LifecycleReviewed:
		// An open tasting note comes before the verdict. The note is the concrete
		// blocker — the close gate refuses a recipe with an open note
		// (`review.clean`) — and the verdict can wait. Triage the note, and the
		// review is offered next.
		if OpenFindings(&r) > 0 {
			return MoveTriage
		}
		if SelfReviewed(r) {
			// The executor tasted its own dish. Somebody else has to — and nobody
			// else is coming on their own, so it is the chef's move.
			return MoveReview
		}
		if r.Lifecycle == LifecycleReviewed {
			if ChangesRequested(r) {
				// The chef asked for changes; the agents have it now.
				return MoveNone
			}
			// Reviewed and clean. The last move is the chef's: serve it.
			return MoveServe
		}
	}
	return MoveNone
}

// Label is what the console's footer calls the move.
func (m ChefMove) Label() string {
	switch m {
	case MoveSignOff:
		return "sign off"
	case MoveReview:
		return "review"
	case MoveServe:
		return "serve"
	case MoveTriage:
		return "triage"
	}
	return ""
}

// ReviewState is the review status of a recipe. Derived.
type ReviewState string

const (
	ReviewNone       ReviewState = ""
	ReviewPending    ReviewState = "pending-review"
	ReviewOpen       ReviewState = "open-findings"
	ReviewRemediated ReviewState = "remediated"
	ReviewSelfSigned ReviewState = "self-reviewed"
	ReviewClean      ReviewState = "reviewed-clean"
)

// AllReviewStates are the states a recipe's review can be in.
var AllReviewStates = []ReviewState{
	ReviewNone, ReviewPending, ReviewOpen, ReviewRemediated, ReviewSelfSigned, ReviewClean,
}

// DerivedView is one view the domain computes and never stores.
//
// The function is the one that computes it, kept as a value. The views take different
// things — a recipe, a slice of recipes, a criterion — so the mapping holds them as
// values and reads their names: the specification's table names the function, and a
// name the code derives cannot drift from the code.
//
// A planned view has no function: the domain does not compute it yet, and the table
// says so.
type DerivedView struct {
	// Name is the view's name, as the specification's table calls it.
	Name string
	// Function is the function that computes the view. It is nil for a planned
	// view.
	Function any
	// Values is the view's closed set, when it has one: the lanes, the chef's
	// moves, the review states.
	Values []string
	// Planned is true for a view the domain does not compute yet.
	Planned bool
}

// DerivedViews are the views the domain computes and never stores.
//
// The specification's table of derived views is checked against this set, both ways:
// a view the code has and the table does not is a view a reader cannot find, and a
// row the table has and the code does not is a view that never happens.
//
// The enumerations are read from the code's own sets, never copied — a value that
// the code gains is a value this table names the same day.
var DerivedViews = []DerivedView{
	{Name: "status", Function: (*Criterion).Status},
	{Name: "review_state", Function: ReviewStateOf, Values: reviewStateValues()},
	{Name: "lane", Function: LaneOf, Values: laneValues()},
	{Name: "readiness", Function: ReadinessOf},
	{Name: "coverage", Planned: true},
	{Name: "evidence_summary", Function: EvidenceSummaryOf},
	{Name: "staleness", Planned: true},
	{Name: "chef_move", Function: ChefMoveOf, Values: chefMoveValues()},
	{Name: "self_reviewed", Function: SelfReviewed},
	{Name: "path", Function: Path},
	{Name: "triage", Function: TriageOf, Values: triageSectionValues()},
	{Name: "tree", Function: TreeOf},
	{Name: "backlog_groups", Function: BacklogGroups, Values: []string{"ideas", "tasks"}},
}

// NonViews are the exported functions of this package that are not derived views.
//
// Every exported function is either a view — a function that carries the `// view`
// marker, which says so where the function is defined — or one of these. A function
// that is neither computes something the specification does not know about, and the
// specification's table is how a reader finds a view — so an unregistered view is a
// hole the wire has to close.
//
// A helper that grows into a view is promoted: marked, registered, and given a row.
//
// A name ending in `.*` is a family: every method of that receiver. The families are
// the **write path** and the **boilerplate**, and they are families because a view is
// not a thing they ever become: the operations and the store write the plan, and the
// labels and the error strings are strings. The read-side methods are named one by
// one — a model's own methods are exactly where the domain put a view
// (`Criterion.Status`), so a family there would admit a view nobody classified.
//
// The list cannot rot. A name here that no longer names a function fails the check —
// an exemption for a function that is gone is a hole in the other direction.
var NonViews = []string{
	// The write path: the operations and the store. A family, not a name each.
	"Operations.*",
	"Store.*",
	// The boilerplate: the labels, the error strings, the renderings.
	"*.Label",
	"*.Error",
	"*.String",
	// The read-side methods, named one by one.
	"Board.Count",
	"Board.Total",
	"Board.Summaries",
	"BacklogItem.Slug",
	"Recipe.Terminal",
	"Recipe.Fingerprint",
	"Journal.Since",
	"Journal.Recent",
	"LegacyItem.Text",
	"Lifecycle.Valid",
	"Lifecycle.Terminal",
	"NoteOutcome.AppliesTo",
	"Verifier.Run",
	// The top-level functions that are not views.
	"Block",
	"BoardOf",
	"Cancel",
	"Cancelled",
	"ChangesRequested",
	"Defer",
	"FindRoot",
	"GateNames",
	"IsAgentRole",
	"KnownLane",
	"LaneNames",
	"NewOperations",
	"NoteOutcomeOf",
	"Now",
	"Open",
	"OpenFindings",
	"OpenFindingsOf",
	"PlanDomainVersion",
	"Resume",
	"Seed",
	"SeedStress",
	"Slugify",
	"SummaryOf",
	"Transition",
	"TriageFlatten",
	"TriageSections",
	"Unblock",
}

// funcName is the name of a function, as the specification's table writes it.
//
// The runtime knows a function's name; the table must not keep its own copy of it, or
// a function that was renamed leaves the table naming something that is not there.
func funcName(f any) string {
	v := reflect.ValueOf(f)
	if !v.IsValid() || v.Kind() != reflect.Func {
		// Not a function: the view names nothing the code has.
		return ""
	}
	name := runtime.FuncForPC(v.Pointer()).Name()
	// A method expression: "…/core.(*Criterion).Status".
	name = strings.ReplaceAll(name, "(*", "")
	name = strings.ReplaceAll(name, ")", "")
	// The package path goes; the qualified name stays.
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	// "core.LaneOf" is "LaneOf"; "core.Criterion.Status" is "Criterion.Status".
	if i := strings.Index(name, "."); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// laneValues are the lanes, as names.
func laneValues() []string {
	out := make([]string, 0, len(AllLanes))
	for _, l := range AllLanes {
		out = append(out, string(l))
	}
	return out
}

// chefMoveValues are the chef's moves, as names.
func chefMoveValues() []string {
	out := make([]string, 0, len(AllChefMoves))
	for _, mv := range AllChefMoves {
		out = append(out, string(mv))
	}
	return out
}

// reviewStateValues are the review states, as names.
func reviewStateValues() []string {
	out := make([]string, 0, len(AllReviewStates))
	for _, s := range AllReviewStates {
		out = append(out, string(s))
	}
	return out
}

// triageSectionValues are the board's sections, as names.
func triageSectionValues() []string {
	return append([]string{}, AllTriageSections...)
}

// Label is the board's name for a review state.
func (s ReviewState) Label() string {
	switch s {
	case ReviewPending:
		return "awaiting tasting"
	case ReviewOpen:
		return "findings open"
	case ReviewRemediated:
		return "remediated"
	case ReviewSelfSigned:
		return "self-reviewed"
	case ReviewClean:
		return "tasting clean"
	}
	return "not tasted"
}

// SelfReviewed reports whether a recipe's latest review round cannot count as an
// independent tasting.
//
// It is true when the executor signed off its own work, and when the work has no
// claimant at all: nobody signed for the work, so there is no executor to be separate
// from, and an anonymous claim cannot buy independence.
//
// The executor and the reviewer are compared as **agents**, not as strings. The
// rename retired a name for each role, and an agent that wrote its old name before
// the rename and its new name after is the same agent: comparing the strings alone
// would let a recipe close on a review its own executor gave.
//
// One source of truth: the gate that blocks the transition and the triage that tells
// the chef what to do both read this. v1's lesson was five implementations of one
// projection, and four of them were wrong.
// view
func SelfReviewed(r Recipe) bool {
	if len(r.Review.Rounds) == 0 {
		return false
	}
	claimer := r.Execution.ClaimedBy
	if claimer == "" {
		return true
	}
	return sameAgent(claimer, r.Review.Rounds[len(r.Review.Rounds)-1].Reviewer)
}

// sameAgent reports whether two actor names can name the same agent.
//
// The names are compared as written and as roles: `line-cook` and `runner` are one
// agent under two names, and `taster` and `reviewer` are one agent under two names.
func sameAgent(a, b string) bool {
	na, nb := normalizeAgentName(a), normalizeAgentName(b)
	if na == nb {
		return true
	}
	return agentRole(na) == agentRole(nb)
}

// ReviewStateOf derives the review state.
// view
func ReviewStateOf(r Recipe) ReviewState {
	open := 0
	for _, f := range r.Review.Findings {
		if f.Status == "open" {
			open++
		}
	}
	if len(r.Review.Rounds) == 0 && len(r.Review.Findings) == 0 {
		return ReviewNone
	}
	if open > 0 {
		return ReviewOpen
	}
	if SelfReviewed(r) {
		// The round exists and it says "ok", but the only person who tasted this
		// dish is the one who cooked it. Calling that clean would be a lie.
		return ReviewSelfSigned
	}
	if len(r.Review.Rounds) > 0 {
		return ReviewClean
	}
	return ReviewRemediated
}

// Readiness says whether a recipe can be picked up, and why not if it cannot.
type Readiness struct {
	Ready  bool
	Reason string
}

// ReadinessOf derives whether a recipe is ready to be claimed.
// view
func ReadinessOf(r Recipe) Readiness {
	if r.Overlays.Cancelled != nil {
		return Readiness{false, "cancelled"}
	}
	if r.Overlays.Blocked != nil {
		return Readiness{false, "blocked: " + r.Overlays.Blocked.Reason}
	}
	switch r.Lifecycle {
	case LifecycleApproved:
		return Readiness{true, ""}
	case LifecycleReady:
		return Readiness{false, "awaiting sign-off"}
	case LifecycleExecuting:
		return Readiness{false, "already on the line"}
	case LifecycleExecuted, LifecycleRemediation, LifecycleReviewed:
		if SelfReviewed(r) {
			// Nobody but the executor has tasted this. The recipe cannot leave the
			// stage until someone else looks at it, and that is the chef's move.
			return Readiness{false, "awaiting an independent review"}
		}
		if r.Lifecycle == LifecycleReviewed && ChangesRequested(r) {
			return Readiness{false, "changes requested — the agents have it"}
		}
		if r.Lifecycle == LifecycleReviewed {
			return Readiness{false, "reviewed — it wants serving"}
		}
		return Readiness{false, "awaiting tasting"}
	case LifecycleClosed:
		return Readiness{false, "served"}
	default:
		return Readiness{false, "still being shaped"}
	}
}

// EvidenceSummary counts the evidence attached to a recipe.
//
// Rejected and unproven are separate: a rejected evidence matched nothing (a spec
// defect), an unproven one could not be attributed (not a machine evidence).
// Different remedies.
type EvidenceSummary struct {
	Criteria int
	Proven   int
	Failed   int
	Rejected int
	Unproven int
	Pending  int
}

// EvidenceSummaryOf summarises the evidence on a recipe's criteria.
// view
func EvidenceSummaryOf(r Recipe) EvidenceSummary {
	s := EvidenceSummary{Criteria: len(r.Acceptance)}
	for _, c := range r.Acceptance {
		switch c.Status() {
		case CriterionPassed:
			s.Proven++
		case CriterionFailed:
			s.Failed++
		case CriterionRejected:
			s.Rejected++
		case CriterionUnproven:
			s.Unproven++
		default:
			s.Pending++
		}
	}
	return s
}

// Board is the whole plan, grouped by lane. It is a projection: it reads the recipe
// files and computes everything the human interface shows.
type Board struct {
	Root    string
	Project Project
	Lanes   map[Lane][]Recipe
}

// Project is the plan-level metadata.
type Project struct {
	Name    string
	Created string
}

// BoardOf groups recipes into lanes.
func BoardOf(root string, project Project, recipes []Recipe) Board {
	b := Board{Root: root, Project: project, Lanes: map[Lane][]Recipe{}}
	for _, l := range AllLanes {
		b.Lanes[l] = []Recipe{}
	}
	for _, r := range recipes {
		lane := LaneOf(r)
		b.Lanes[lane] = append(b.Lanes[lane], r)
	}
	for l := range b.Lanes {
		sort.Slice(b.Lanes[l], func(i, j int) bool {
			return recipeLess(b.Lanes[l][i].ID, b.Lanes[l][j].ID)
		})
	}
	return b
}

// Count is the number of recipes in the lane.
func (b Board) Count(l Lane) int { return len(b.Lanes[l]) }

// Total is the number of recipes across every lane.
func (b Board) Total() int {
	n := 0
	for _, l := range AllLanes {
		n += len(b.Lanes[l])
	}
	return n
}

// Summary is the one-line state of a recipe, for the board list.
type Summary struct {
	ID           string
	Title        string
	Lifecycle    Lifecycle
	Lane         Lane
	Review       ReviewState
	Readiness    Readiness
	Evidence     EvidenceSummary
	OpenFindings int
	Updated      string
}

// SummaryOf projects a recipe into its board row.
func SummaryOf(r Recipe) Summary {
	return Summary{
		ID:           r.ID,
		Title:        r.Title,
		Lifecycle:    r.Lifecycle,
		Lane:         LaneOf(r),
		Review:       ReviewStateOf(r),
		Readiness:    ReadinessOf(r),
		Evidence:     EvidenceSummaryOf(r),
		OpenFindings: OpenFindings(&r),
		Updated:      r.Updated,
	}
}

// Summaries projects a lane into board rows.
func (b Board) Summaries(l Lane) []Summary {
	rows := make([]Summary, 0, len(b.Lanes[l]))
	for _, r := range b.Lanes[l] {
		rows = append(rows, SummaryOf(r))
	}
	return rows
}

// TriageSection is a group of recipes by what needs the chef's attention.
//
// This is the console's organising axis — not the lifecycle lane. A lane says *where a
// recipe is*; a triage section says *what the chef should do next*. The two overlap but
// are not the same thing, and the console shows the latter.
type TriageSection struct {
	Name  string
	Items []Summary
}

// TriageSection names, in display order.
const (
	TriageNeedsYou    = "Needs you"
	TriageOnTheLine   = "On the line"
	TriageBeingShaped = "Being shaped"
	TriageOnHold      = "On hold"
	TriageServed      = "Served"
)

// AllTriageSections are the sections the board is grouped into, in display order.
//
// The triage is the board's other axis — the lane says *where a recipe is*, the triage
// says *what the chef should do next*. It is a closed set, so the console and the
// specification read it from here.
var AllTriageSections = []string{
	TriageNeedsYou, TriageOnTheLine, TriageBeingShaped, TriageOnHold, TriageServed,
}

// TriageSections groups recipes into the sections the console shows.
//
// Derived, never stored: the chef's view of the plan is not part of the plan.
func TriageSections(recipes []Recipe) []TriageSection {
	byName := map[string][]Summary{}
	for _, r := range recipes {
		s := SummaryOf(r)
		byName[TriageOf(r)] = append(byName[TriageOf(r)], s)
	}

	out := make([]TriageSection, 0, len(AllTriageSections))
	for _, name := range AllTriageSections {
		items := byName[name]
		if len(items) == 0 {
			continue
		}
		sort.Slice(items, func(i, j int) bool { return recipeLess(items[i].ID, items[j].ID) })
		out = append(out, TriageSection{Name: name, Items: items})
	}
	return out
}

// TriageOf decides which section a recipe belongs to.
//
// The cancelled case comes first, and it reads the same fact the lane reads: a
// cancelled recipe is over, and the two orderings must say so together.
// view
func TriageOf(r Recipe) string {
	switch {
	case Cancelled(r):
		return TriageServed
	case r.Overlays.Blocked != nil:
		return TriageOnHold
	case r.Lifecycle == LifecycleClosed:
		return TriageServed
	case r.Lifecycle == LifecycleDraft, r.Lifecycle == LifecycleGroomed:
		return TriageBeingShaped
	case r.Lifecycle == LifecycleExecuting:
		return TriageOnTheLine
	case r.Lifecycle == LifecycleReady:
		// Waiting for the chef to sign off.
		return TriageNeedsYou
	case r.Lifecycle == LifecycleApproved:
		// Approved and nobody has started it. The agents have it — the console has
		// no move that would unstick it, so it does not claim the chef's attention.
		return TriageOnTheLine
	case r.Lifecycle == LifecycleExecuted, r.Lifecycle == LifecycleRemediation, r.Lifecycle == LifecycleReviewed:
		// Only needs the chef if there is something to triage, if the only review so
		// far is the executor's own — then nobody else has looked at the work — or
		// if it is reviewed and waiting to be served.
		if n := OpenFindings(&r); n > 0 {
			return TriageNeedsYou
		}
		if SelfReviewed(r) {
			return TriageNeedsYou
		}
		if r.Lifecycle == LifecycleReviewed && !ChangesRequested(r) {
			return TriageNeedsYou
		}
		return TriageOnTheLine
	}
	return TriageNeedsYou
}

// TriageFlatten returns the section summaries in display order, skipping the served
// section when it is hidden.
func TriageFlatten(sections []TriageSection, showServed bool) []Summary {
	var out []Summary
	for _, s := range sections {
		if s.Name == TriageServed && !showServed {
			continue
		}
		out = append(out, s.Items...)
	}
	return out
}

// PathEntry is one recipe in the execution order.
type PathEntry struct {
	Summary Summary
	// WaitingFor names the dependency that is not closed yet.
	WaitingFor string
	// Note is a short hint for the chef.
	Note string
}

// Path groups recipes into the execution order: what can go now, and what is waiting.
//
// Neither list is ever nil: an empty plan is an empty list, and the CLI marshals an
// empty list as `[]` and a nil slice as `null`. A client cannot iterate `null`.
// view
func Path(recipes []Recipe) (ready, waiting []PathEntry) {
	closed := closedRecipes(recipes)
	ready = make([]PathEntry, 0, len(recipes))
	waiting = make([]PathEntry, 0, len(recipes))
	for _, r := range recipes {
		if r.Lifecycle == LifecycleClosed || r.Overlays.Cancelled != nil {
			continue
		}
		s := SummaryOf(r)
		entry := PathEntry{Summary: s}
		if deps := unmetDependencies(r, closed); len(deps) > 0 {
			entry.WaitingFor = deps[0]
		}
		switch {
		case entry.WaitingFor != "":
			entry.Note = "needs " + entry.WaitingFor
			waiting = append(waiting, entry)
		case r.Lifecycle == LifecycleReady:
			entry.Note = "needs you"
			ready = append(ready, entry)
		case r.Lifecycle == LifecycleApproved:
			entry.Note = "ready to run"
			ready = append(ready, entry)
		case r.Lifecycle == LifecycleExecuting:
			entry.Note = "running"
			ready = append(ready, entry)
		case SelfReviewed(r):
			// The only review so far is the executor's own, so the chef is the
			// one who has to get a second pair of eyes onto it.
			entry.Note = "needs you"
			ready = append(ready, entry)
		default:
			entry.Note = "later"
			ready = append(ready, entry)
		}
	}
	sort.Slice(ready, func(i, j int) bool { return recipeLess(ready[i].Summary.ID, ready[j].Summary.ID) })
	sort.Slice(waiting, func(i, j int) bool { return recipeLess(waiting[i].Summary.ID, waiting[j].Summary.ID) })
	return ready, waiting
}

// String renders a recipe's board row.
func (s Summary) String() string {
	return fmt.Sprintf("%s  %-10s %s", s.ID, s.Lifecycle, s.Title)
}

// TreeNode is one recipe in the dependency tree.
//
// The tree is the plan as the chef reads it: what can go now, and what is stuck
// behind what.
type TreeNode struct {
	// Summary is the recipe's one-line state.
	Summary Summary
	// Needs are the dependencies that are not closed yet, in the recipe's own order.
	Needs []string
	// Missing names the unmet dependencies that are not in the plan at all. A
	// dependency that is not there is a broken plan, and the tree is where the
	// chef sees it.
	Missing []string
	// Cycle is true when the recipe waits on itself, directly or through other
	// recipes. The plan cannot be executed as written.
	Cycle bool
	// WaitsOnCycle is true when the recipe waits on a recipe that is in a cycle. A
	// recipe in a cycle itself is not counted here: it is the knot, not something
	// downstream of it.
	WaitsOnCycle bool
	// Waiting are the recipes that fork at this one.
	Waiting []TreeNode
}

// Tree is the plan as a dependency tree.
type Tree struct {
	// Roots are the recipes that can go now: the Path tab's ready list, in the same
	// order. Two views of one plan must not tell two stories.
	Roots []TreeNode
	// Orphans are the recipes that wait on something the plan does not have. They
	// cannot fork at anything, so they are their own thing.
	Orphans []TreeNode
	// Stuck are the recipes a dependency cycle keeps out of the tree: the members of
	// the cycle, and the recipes that wait on them. The walk from the roots cannot
	// reach them, and a tree that drops them tells the chef a story the plan does
	// not — so they are drawn, as what they are: a broken plan.
	Stuck []TreeNode
}

// TreeOf builds the plan's dependency tree.
//
// The trunk is the Path tab's ready list, in the same order — `Path` already computes
// the ready set and the first unmet dependency, so the tree reads that answer rather
// than computing a second one.
//
// A recipe that waits on something forks at its **first** unmet dependency, the same
// rule `Path` uses.
//
// The walk starts at the trunk and stops at a recipe it has already drawn, so a
// dependency cycle cannot make the tree infinite. What the walk cannot reach is
// drawn anyway: a tree that silently drops a recipe tells the chef a story the plan
// does not, and the Path tab would list what the tree had forgotten.
// view
func TreeOf(recipes []Recipe) Tree {
	closed := closedRecipes(recipes)
	ready, _ := Path(recipes)

	byID := map[string]Recipe{}
	for _, r := range recipes {
		byID[r.ID] = r
	}

	// The forks: a recipe forks at its first unmet dependency.
	forks := map[string][]string{}
	for _, r := range recipes {
		if r.Lifecycle == LifecycleClosed || r.Overlays.Cancelled != nil {
			continue
		}
		deps := unmetDependencies(r, closed)
		if len(deps) == 0 {
			continue
		}
		forks[deps[0]] = append(forks[deps[0]], r.ID)
	}

	// The cycles: the recipes that wait on themselves.
	cycles := map[string]bool{}
	for _, r := range recipes {
		if r.Lifecycle == LifecycleClosed || r.Overlays.Cancelled != nil {
			continue
		}
		if dependsOn(r.ID, byID, closed)[r.ID] {
			cycles[r.ID] = true
		}
	}

	seen := map[string]bool{}
	var build func(id string) TreeNode
	build = func(id string) TreeNode {
		r := byID[id]
		deps := unmetDependencies(r, closed)
		node := TreeNode{
			Summary: SummaryOf(r),
			Needs:   deps,
		}
		for _, dep := range deps {
			if _, ok := byID[dep]; !ok {
				node.Missing = append(node.Missing, dep)
			}
		}
		// A cycle member cannot be executed, and neither can a recipe that waits
		// on one.
		reaches := dependsOn(id, byID, closed)
		node.Cycle = reaches[id]
		if !node.Cycle {
			for dep := range reaches {
				if cycles[dep] {
					node.WaitsOnCycle = true
					break
				}
			}
		}
		for _, child := range forks[id] {
			if seen[child] {
				continue
			}
			seen[child] = true
			node.Waiting = append(node.Waiting, build(child))
		}
		return node
	}

	tree := Tree{Roots: make([]TreeNode, 0, len(ready)), Orphans: []TreeNode{}, Stuck: []TreeNode{}}
	for _, e := range ready {
		seen[e.Summary.ID] = true
		tree.Roots = append(tree.Roots, build(e.Summary.ID))
	}
	// The leftovers: what the walk could not reach. A recipe that waits on something
	// the plan does not have is an orphan; a recipe whose dependencies are all there
	// and that nothing reaches is in a cycle.
	for _, r := range recipes {
		if r.Lifecycle == LifecycleClosed || r.Overlays.Cancelled != nil {
			continue
		}
		if seen[r.ID] {
			continue
		}
		seen[r.ID] = true
		node := build(r.ID)
		if len(node.Missing) > 0 {
			tree.Orphans = append(tree.Orphans, node)
			continue
		}
		tree.Stuck = append(tree.Stuck, node)
	}
	return tree
}

// dependsOn is the set of recipes a recipe transitively waits on, through the
// dependencies that are not closed yet.
//
// A recipe that is in the set waits on itself: it is in a dependency cycle, and the
// plan cannot be executed as written.
func dependsOn(id string, byID map[string]Recipe, closed map[string]bool) map[string]bool {
	out := map[string]bool{}
	var walk func(cur string)
	walk = func(cur string) {
		for _, dep := range unmetDependencies(byID[cur], closed) {
			if out[dep] {
				continue
			}
			out[dep] = true
			walk(dep)
		}
	}
	walk(id)
	return out
}

// closedRecipes are the ids of the recipes that are closed.
func closedRecipes(recipes []Recipe) map[string]bool {
	closed := map[string]bool{}
	for _, r := range recipes {
		if r.Lifecycle == LifecycleClosed {
			closed[r.ID] = true
		}
	}
	return closed
}

// unmetDependencies are the dependencies of a recipe that are not closed yet, in the
// recipe's own order.
//
// One source: `Path` and the tree both read this, so the flat order and the tree
// cannot disagree about what waits on what.
func unmetDependencies(r Recipe, closed map[string]bool) []string {
	var out []string
	for _, dep := range r.DependsOn {
		if !closed[dep] {
			out = append(out, dep)
		}
	}
	return out
}
