package cli

import (
	"strings"

	"github.com/lthiagol/maitre/internal/core"
)

// ---------------------------------------------------------------------------
// Acceptance criteria
// ---------------------------------------------------------------------------

func (r *Runner) cmdACAdd(args []string) int {
	fs := newFlagSet("ac add")
	description := fs.String("description", "", "what must be true")
	verification := fs.String("verification", "", "the command that checks it")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli ac add <id> --description <text> --verification <command>")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.AddCriterion(fs.arg(0), *description, *verification)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

// cmdACVerify runs the criterion's verification and records the evidence.
//
// The tool runs the command. An evidence that matched nothing is rejected; one the
// tool could not attribute is unproven. Neither satisfies a gate.
func (r *Runner) cmdACVerify(args []string) int {
	fs := newFlagSet("ac verify")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 2 {
		r.fail("usage: maitre-cli ac verify <id> <AC>")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	criterion, err := ops.VerifyCriterion(fs.arg(0), fs.arg(1))
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{
		"ok":        true,
		"criterion": criterion,
		"status":    string(criterion.Status()),
	})
}

func (r *Runner) cmdACSignOff(args []string) int {
	fs := newFlagSet("ac signoff")
	by := fs.String("by", "", "who is vouching")
	note := fs.String("note", "", "the procedure that was followed")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 2 {
		r.fail("usage: maitre-cli ac signoff <id> <AC> --by <name> --note <text>")
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
	rec, err := ops.SignOffCriterion(fs.arg(0), fs.arg(1), *by, *note)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdACDefer(args []string) int {
	fs := newFlagSet("ac defer")
	reason := fs.String("reason", "", "why it is being deferred")
	followUp := fs.String("follow-up", "", "where it was re-filed")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 2 {
		r.fail("usage: maitre-cli ac defer <id> <AC> --reason <text> [--follow-up <id>]")
		return ExitError
	}
	if strings.TrimSpace(*reason) == "" {
		r.fail("a deferral needs --reason")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.DeferCriterion(fs.arg(0), fs.arg(1), *reason, *followUp)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

// ---------------------------------------------------------------------------
// Work packages
// ---------------------------------------------------------------------------

// cmdWPAdd adds a work package — the grouping the steps belong to.
//
// The specify gate requires at least one, so a recipe cannot be locked until it has
// been decomposed.
func (r *Runner) cmdWPAdd(args []string) int {
	fs := newFlagSet("wp add")
	name := fs.String("name", "", "the package's name")
	goal := fs.String("goal", "", "what it achieves")
	rollback := fs.String("rollback", "", "how to undo it")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli wp add <id> --name <text> [--goal] [--rollback]")
		return ExitError
	}
	if strings.TrimSpace(*name) == "" {
		r.fail("missing --name")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.AddWorkPackage(fs.arg(0), *name, *goal, *rollback)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

// ---------------------------------------------------------------------------
// Guidance steps
// ---------------------------------------------------------------------------

func (r *Runner) cmdStepAdd(args []string) int {
	fs := newFlagSet("step add")
	action := fs.String("action", "", "what to do")
	doneWhen := fs.String("done-when", "", "how you will know it worked")
	wp := fs.String("wp", "", "the work package it belongs to")
	covers := fs.StringList("covers")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli step add <id> --action <text> --done-when <text> [--covers AC-01]")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	step := core.StepSpec{Action: *action, DoneWhen: *doneWhen, Covers: *covers, WorkPackage: *wp}
	rec, err := ops.AddStep(fs.arg(0), step)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdStepDone(args []string) int {
	return r.stepStatus(args, "done", "done")
}

func (r *Runner) cmdStepSkip(args []string) int {
	return r.stepStatus(args, "skip", "skipped")
}

// cmdStepStart marks a step as under way. Every documented status has a verb:
// pending · in-progress · done · skipped · failed.
func (r *Runner) cmdStepStart(args []string) int {
	return r.stepStatus(args, "start", "in-progress")
}

// cmdStepFail marks a step as failed.
func (r *Runner) cmdStepFail(args []string) int {
	return r.stepStatus(args, "fail", "failed")
}

func (r *Runner) stepStatus(args []string, verb, status string) int {
	fs := newFlagSet("step " + verb)
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 2 {
		r.fail("usage: maitre-cli step " + verb + " <id> <G1>")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.SetStepStatus(fs.arg(0), fs.arg(1), status)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdStepDeviate(args []string) int {
	fs := newFlagSet("step deviate")
	note := fs.String("note", "", "how it differed from the method")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 2 {
		r.fail("usage: maitre-cli step deviate <id> <G1> --note <text>")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.RecordDeviation(fs.arg(0), fs.arg(1), *note)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

// ---------------------------------------------------------------------------
// Findings
// ---------------------------------------------------------------------------

func (r *Runner) cmdFindingAdd(args []string) int {
	fs := newFlagSet("finding add")
	description := fs.String("description", "", "what is wrong")
	severity := fs.String("severity", "medium", "low, medium or high")
	category := fs.String("category", "", "correctness, docs, performance …")
	phase := fs.String("phase", "external", "self or external")
	file := fs.String("file", "", "the file the finding is about")
	line := fs.Int("line", 0, "the line")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli finding add <id> --description <text> [--severity] [--category] [--phase]")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	finding := core.Finding{
		Phase:       *phase,
		Severity:    *severity,
		Category:    *category,
		Description: *description,
	}
	if *file != "" {
		finding.Anchor = &core.Anchor{File: *file, Line: *line}
	}
	rec, err := ops.AddFinding(fs.arg(0), finding)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdFindingResolve(args []string) int {
	fs := newFlagSet("finding resolve")
	note := fs.String("note", "", "how it was fixed")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 2 {
		r.fail("usage: maitre-cli finding resolve <id> <F-01> [--note <text>]")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.ResolveFinding(fs.arg(0), fs.arg(1), *note)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

// cmdFindingReopen sends a finding back into the queue.
//
// The reviewer's move when a fix does not hold: the finding goes back to open, and the
// objection is recorded so the agent that tries next can read why.
func (r *Runner) cmdFindingReopen(args []string) int {
	fs := newFlagSet("finding reopen")
	reason := fs.String("reason", "", "why it is not fixed")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 2 {
		r.fail("usage: maitre-cli finding reopen <id> <F-01> --reason <text>")
		return ExitError
	}
	if strings.TrimSpace(*reason) == "" {
		r.fail("a reopen needs --reason")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.ReopenFinding(fs.arg(0), fs.arg(1), *reason)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}

func (r *Runner) cmdFindingDismiss(args []string) int {
	fs := newFlagSet("finding dismiss")
	reason := fs.String("reason", "", "why it is not worth fixing")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 2 {
		r.fail("usage: maitre-cli finding dismiss <id> <F-01> --reason <text>")
		return ExitError
	}
	if strings.TrimSpace(*reason) == "" {
		r.fail("a dismissal needs --reason")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := ops.DismissFinding(fs.arg(0), fs.arg(1), *reason)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(rec)})
}
