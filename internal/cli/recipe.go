package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lthiagol/maitre/internal/core"
)

// errorsAs is errors.As, named short because it is used a lot.
func errorsAs(err error, target any) bool { return errors.As(err, target) }

// recipeSummary is the compact projection a mutation returns.
func recipeSummary(rec *core.Recipe) map[string]any {
	criteria := make([]map[string]any, 0, len(rec.Acceptance))
	for _, c := range rec.Acceptance {
		criteria = append(criteria, map[string]any{
			"id":           c.ID,
			"description":  c.Description,
			"verification": c.Verification,
			"status":       string(c.Status()),
		})
	}
	return map[string]any{
		"id":         rec.ID,
		"title":      rec.Title,
		"lifecycle":  string(rec.Lifecycle),
		"outcome":    rec.Intent.Outcome,
		"acceptance": criteria,
		"guidance":   len(rec.Guidance),
		"updated":    rec.Updated,
	}
}

// ---------------------------------------------------------------------------
// The recipe's own state
// ---------------------------------------------------------------------------

func (r *Runner) cmdRecipeCreate(args []string) int {
	fs := newFlagSet("recipe create")
	title := fs.String("title", "", "title")
	outcome := fs.String("outcome", "", "intent.outcome")
	problem := fs.String("problem", "", "intent.problem")
	inScope := fs.StringList("in-scope")
	outScope := fs.StringList("out-of-scope")
	ac := fs.StringList("ac")
	acVerify := fs.StringList("ac-verify")
	wp := fs.StringList("wp")
	step := fs.StringList("step")
	stepDone := fs.StringList("step-done-when")
	stepCovers := fs.StringList("step-covers")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if strings.TrimSpace(*title) == "" {
		r.fail("missing --title")
		return ExitError
	}
	if len(*ac) != len(*acVerify) {
		r.fail("--ac and --ac-verify must be given the same number of times")
		return ExitError
	}
	if len(*step) != len(*stepDone) {
		r.fail("--step and --step-done-when must be given the same number of times")
		return ExitError
	}

	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	spec := core.RecipeSpec{
		Title:        *title,
		Outcome:      *outcome,
		Problem:      *problem,
		InScope:      *inScope,
		OutOfScope:   *outScope,
		WorkPackages: *wp,
	}
	for i := range *ac {
		spec.Criteria = append(spec.Criteria, core.CriterionSpec{
			Description:  (*ac)[i],
			Verification: (*acVerify)[i],
		})
	}
	for i := range *step {
		spec.Steps = append(spec.Steps, core.StepSpec{
			Action:   (*step)[i],
			DoneWhen: (*stepDone)[i],
		})
		if i < len(*stepCovers) {
			spec.Steps[i].Covers = []string{(*stepCovers)[i]}
		}
	}
	rec, err := ops.CreateRecipe(spec)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdRecipeList(args []string) int {
	fs := newFlagSet("recipe list")
	lane := fs.String("lane", "", "only recipes in this lane")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	// A typo must be a loud error, not an empty list: v1's `--type` answered "there
	// are none" for a lane that does not exist.
	if *lane != "" && !core.KnownLane(*lane) {
		r.fail("unknown lane", *lane, "known lanes: "+strings.Join(core.LaneNames(), ", "))
		return ExitError
	}
	store, _, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	recipes, err := store.ListRecipes()
	if err != nil {
		r.fail("list recipes", err.Error())
		return ExitError
	}
	out := make([]map[string]any, 0, len(recipes))
	for _, rec := range recipes {
		if *lane != "" && string(core.LaneOf(rec)) != *lane {
			continue
		}
		out = append(out, recipeSummary(&rec))
	}
	return r.emit(map[string]any{"ok": true, "count": len(out), "recipes": out})
}

func (r *Runner) cmdRecipeShow(args []string) int {
	fs := newFlagSet("recipe show")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli recipe show <id>")
		return ExitError
	}
	store, _, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := store.LoadRecipe(fs.arg(0))
	if err != nil {
		r.fail("no such recipe", fs.arg(0))
		return ExitError
	}
	return r.emit(map[string]any{"ok": true, "recipe": rec})
}

func (r *Runner) cmdRecipeUpdate(args []string) int {
	fs := newFlagSet("recipe update")
	title := fs.String("title", "", "new title")
	outcome := fs.String("outcome", "", "new intent.outcome")
	problem := fs.String("problem", "", "new intent.problem")
	inScope := fs.StringList("in-scope")
	outScope := fs.StringList("out-of-scope")
	dependsOn := fs.StringList("depends-on")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli recipe update <id> [flags]")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	fragment := core.Fragment{}
	if fs.StringValue("title") != "" {
		fragment.Title = title
	}
	if fs.StringValue("outcome") != "" {
		fragment.Outcome = outcome
	}
	if fs.StringValue("problem") != "" {
		fragment.Problem = problem
	}
	if len(*inScope) > 0 {
		fragment.InScope = inScope
	}
	if len(*outScope) > 0 {
		fragment.OutOfScope = outScope
	}
	if len(*dependsOn) > 0 {
		// The list replaces: `--depends-on RCP02` says the recipe depends on RCP02 and
		// on nothing else. Leaving it alone would need a second flag to clear it,
		// and a dependency you cannot take off is a dependency you cannot undo.
		fragment.DependsOn = dependsOn
	}
	rec, err := ops.UpdateRecipe(fs.arg(0), fragment)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

// cmdRecipeAdvance applies a lifecycle event.
//
// This is the one verb for the transitions that have no extra inputs of their own:
// groom, specify, finish, enter_remediation, exit_remediation, pass_review, close,
// abort_to_draft and reopen. The ones that need an approver, a reason or an agent have
// their own verb — `approve`, `request-changes` and `claim` — and `advance` refuses
// them. One intent, one verb: the alias is the weaker path, because `advance approve`
// signs off with the actor as the approver.
func (r *Runner) cmdRecipeAdvance(args []string) int {
	fs := newFlagSet("recipe advance")
	reason := fs.String("reason", "", "why — required by abort_to_draft")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 2 {
		r.fail("usage: maitre-cli recipe advance <id> <event>")
		return ExitError
	}
	ev := core.Event(fs.arg(1))
	known := false
	for _, e := range core.AllEvents {
		if e == ev {
			known = true
			break
		}
	}
	if !known {
		r.fail("unknown event", fs.arg(1))
		return ExitError
	}
	if own := ownVerb(ev); own != "" {
		r.fail("that event has its own verb", string(ev)+": "+own)
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	ctx := core.TransitionContext{Reason: *reason}
	rec, err := ops.AdvanceWith(fs.arg(0), ev, ctx)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

// ownVerb is the verb that owns an event, when the event has one.
//
// One intent, one verb: `advance` is the verb for the transitions that need nothing
// but the event itself. An event that needs an approver, a reason or an agent has a
// verb of its own, and the alias is the weaker path.
func ownVerb(ev core.Event) string {
	switch ev {
	case core.EventApprove:
		return "maitre-cli recipe approve <id> --by <name>"
	case core.EventRequestChanges:
		return "maitre-cli recipe request-changes <id> --by <name> --notes <text>"
	case core.EventStart:
		return "maitre-cli recipe claim <id> --agent <name>"
	}
	return ""
}

// cmdRecipeApprove signs a recipe off.
func (r *Runner) cmdRecipeApprove(args []string) int {
	fs := newFlagSet("recipe approve")
	by := fs.String("by", "", "who is signing it off")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli recipe approve <id> --by <name>")
		return ExitError
	}
	if strings.TrimSpace(*by) == "" {
		r.fail("a sign-off needs --by")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.ApproveRecipe(fs.arg(0), *by)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdRecipeRequestChanges(args []string) int {
	fs := newFlagSet("recipe request-changes")
	by := fs.String("by", "", "who is asking")
	notes := fs.String("notes", "", "what is wrong with it")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli recipe request-changes <id> --by <name> --notes <text>")
		return ExitError
	}
	if strings.TrimSpace(*notes) == "" {
		r.fail("requesting changes needs --notes")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.RequestChanges(fs.arg(0), *by, *notes)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdRecipeClaim(args []string) int {
	fs := newFlagSet("recipe claim")
	agent := fs.String("agent", "", "who is claiming it")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli recipe claim <id> --agent <name>")
		return ExitError
	}
	if strings.TrimSpace(*agent) == "" {
		r.fail("a claim needs --agent")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.ClaimRecipe(fs.arg(0), *agent)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdRecipeDocument(args []string) int {
	fs := newFlagSet("recipe document")
	text := fs.String("text", "", "the documentation")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli recipe document <id> --text <text>")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.DocumentRecipe(fs.arg(0), *text)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdRecipeReview(args []string) int {
	fs := newFlagSet("recipe review")
	reviewer := fs.String("reviewer", "", "who reviewed it")
	verdict := fs.String("verdict", "", "ok or changes-needed")
	notes := fs.String("notes", "", "what they said")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli recipe review <id> --reviewer <name> --verdict <ok|changes-needed>")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.RecordReviewRound(fs.arg(0), *reviewer, *verdict, *notes)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdRecipeComment(args []string) int {
	fs := newFlagSet("recipe comment")
	author := fs.String("author", "", "who is writing (default: the actor)")
	body := fs.String("body", "", "what they said")
	finding := fs.String("finding", "", "the finding this is about")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli recipe comment <id> --author <name> --body <text>")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	// A comment belongs to somebody: the actor, unless the caller says otherwise.
	// The journal records the actor, so a comment that said nobody would leave the
	// artifact disagreeing with the log.
	id := r.resolveIdentity(fs, identity{Actor: r.actor, Session: r.session, Role: r.role})
	rec, err := ops.AddComment(fs.arg(0), firstNonEmpty(*author, id.Actor), *body, *finding)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdRecipeBlock(args []string) int {
	return r.overlay(args, "block", func(ops *core.Operations, id, reason, by string) (*core.Recipe, error) {
		return ops.BlockRecipe(id, reason)
	})
}

func (r *Runner) cmdRecipeUnblock(args []string) int {
	fs := newFlagSet("recipe unblock")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli recipe unblock <id>")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.UnblockRecipe(fs.arg(0))
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdRecipeDefer(args []string) int {
	return r.overlay(args, "defer", func(ops *core.Operations, id, reason, by string) (*core.Recipe, error) {
		return ops.DeferRecipe(id, reason)
	})
}

func (r *Runner) cmdRecipeResume(args []string) int {
	fs := newFlagSet("recipe resume")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli recipe resume <id>")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.ResumeRecipe(fs.arg(0))
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdRecipeCancel(args []string) int {
	return r.overlay(args, "cancel", func(ops *core.Operations, id, reason, by string) (*core.Recipe, error) {
		return ops.CancelRecipe(id, reason)
	})
}

// overlay handles the verbs that need a reason: block, defer, cancel.
func (r *Runner) overlay(args []string, verb string, fn func(*core.Operations, string, string, string) (*core.Recipe, error)) int {
	fs := newFlagSet("recipe " + verb)
	reason := fs.String("reason", "", "why")
	by := fs.String("by", "", "who")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail(fmt.Sprintf("usage: maitre-cli recipe %s <id> --reason <text>", verb))
		return ExitError
	}
	if strings.TrimSpace(*reason) == "" {
		r.fail(fmt.Sprintf("%s needs --reason", verb))
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := fn(ops, fs.arg(0), *reason, *by)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}
