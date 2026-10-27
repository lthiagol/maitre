package cli

import (
	"strings"

	"github.com/lthiagol/maitre/internal/core"
)

// The backlog. One artifact, with a kind.
//
// The CLI's verbs are thin adapters over the operations: the domain owns the rules.

// ---------------------------------------------------------------------------
// The backlog
// ---------------------------------------------------------------------------

func (r *Runner) cmdSpecialList(args []string) int {
	fs := newFlagSet("special list")
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

func (r *Runner) cmdSpecialCreate(args []string) int {
	fs := newFlagSet("special create")
	title := fs.String("title", "", "what it is")
	body := fs.String("body", "", "why it is worth doing")
	tags := fs.StringList("tag")
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
	item, err := ops.CreateBacklogItem(*title, *body, nil, core.KindIdea, *tags)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "item": item})
}

func (r *Runner) cmdSpecialUpdate(args []string) int {
	fs := newFlagSet("special update")
	title := fs.String("title", "", "new title")
	body := fs.String("body", "", "new body")
	tags := fs.StringList("tag")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli special update <id> [--title] [--body] [--tag]")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	item, err := ops.UpdateBacklogItem(fs.arg(0), *title, *body, nil, "", *tags)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "item": item})
}

func (r *Runner) cmdSpecialDismiss(args []string) int {
	fs := newFlagSet("special dismiss")
	reason := fs.String("reason", "", "why it is not worth pursuing")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli special dismiss <id> --reason <text>")
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

func (r *Runner) cmdSpecialPromote(args []string) int {
	fs := newFlagSet("special promote")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli special promote <id>")
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

// ---------------------------------------------------------------------------
// The prep notes
// ---------------------------------------------------------------------------

func (r *Runner) cmdPrepNoteList(args []string) int {
	fs := newFlagSet("prepnote list")
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

func (r *Runner) cmdPrepNoteCreate(args []string) int {
	fs := newFlagSet("prepnote create")
	title := fs.String("title", "", "what it is")
	intent := fs.String("intent", "", "why it is worth doing")
	scope := fs.StringList("scope")
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
	item, err := ops.CreateBacklogItem(*title, *intent, *scope, core.KindFeature, nil)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "item": item})
}

func (r *Runner) cmdPrepNoteUpdate(args []string) int {
	fs := newFlagSet("prepnote update")
	intent := fs.String("intent", "", "new body")
	scope := fs.StringList("scope")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli prepnote update <id> [--intent] [--scope]")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	item, err := ops.UpdateBacklogItem(fs.arg(0), "", *intent, *scope, "", nil)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "item": item})
}

func (r *Runner) cmdPrepNotePromote(args []string) int {
	fs := newFlagSet("prepnote promote")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli prepnote promote <id>")
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
