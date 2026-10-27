package cli

import (
	"strings"

	"github.com/lthiagol/maitre/internal/core"
)

// The backlog. One noun for the list.
//
// The verbs are thin adapters over the operations: the domain owns the rules — what a
// kind is, what a promotion does, and what may not be edited.

// cmdBacklogList lists the backlog.
func (r *Runner) cmdBacklogList(args []string) int {
	fs := newFlagSet("backlog list")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	store, _, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	items, err := store.ListBacklogItems()
	if err != nil {
		r.fail("list the backlog", err.Error())
		return ExitError
	}
	return r.emit(map[string]any{"ok": true, "count": len(items), "items": items})
}

// cmdBacklogCreate files a backlog item.
//
// The kind is a flag, and it defaults to `idea`: an item starts as a hunch, and
// shaping it into a small task is an update.
func (r *Runner) cmdBacklogCreate(args []string) int {
	fs := newFlagSet("backlog create")
	title := fs.String("title", "", "what it is")
	body := fs.String("body", "", "why it is worth doing")
	kind := fs.String("kind", core.KindIdea, "idea, feature or bug")
	// The list flags are declared here and read below: a list flag that is not
	// given leaves its field alone.
	fs.StringList("scope")
	fs.StringList("tag")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if strings.TrimSpace(*title) == "" {
		r.fail("missing --title")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	item, err := ops.CreateBacklogItem(*title, *body, fs.listArg("scope"), *kind, fs.listArg("tag"))
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "item": item})
}

// cmdBacklogUpdate edits a backlog item.
//
// Changing the kind is the shaping step: an idea becomes a small task by being made a
// `feature` or a `bug` and gaining a scope. The id does not change — the item stays
// one item, and the trail survives.
func (r *Runner) cmdBacklogUpdate(args []string) int {
	fs := newFlagSet("backlog update")
	title := fs.String("title", "", "new title")
	body := fs.String("body", "", "new body")
	kind := fs.String("kind", "", "idea, feature or bug")
	// The list flags are declared here and read below: a list flag that is not
	// given leaves its field alone.
	fs.StringList("scope")
	fs.StringList("tag")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli backlog update <id> [--title] [--body] [--scope] [--kind] [--tag]")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	item, err := ops.UpdateBacklogItem(fs.arg(0), *title, *body, fs.listArg("scope"), *kind, fs.listArg("tag"))
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "item": item})
}

// cmdBacklogDismiss marks an item as not worth pursuing.
//
// Every kind can be dismissed: an idea, a feature and a bug are all things the chef
// can decide against.
func (r *Runner) cmdBacklogDismiss(args []string) int {
	fs := newFlagSet("backlog dismiss")
	reason := fs.String("reason", "", "why it is not worth pursuing")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli backlog dismiss <id> --reason <text>")
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
	item, err := ops.DismissBacklogItem(fs.arg(0), *reason)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "item": item})
}

// cmdBacklogPromote turns a backlog item into a recipe.
func (r *Runner) cmdBacklogPromote(args []string) int {
	fs := newFlagSet("backlog promote")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli backlog promote <id>")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	recipe, err := ops.PromoteBacklogItem(fs.arg(0))
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(recipe)})
}
