package cli

import (
	"strings"
)

// The lower rungs of the ladder. Promotion is one-way: a special becomes a prep note,
// a prep note becomes a recipe, and the artifact records where it went.

// ---------------------------------------------------------------------------
// Specials — rung 0
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
	specials, err := store.ListSpecials()
	if err != nil {
		r.fail("list specials", err.Error())
		return ExitError
	}
	return r.emit(map[string]any{"ok": true, "count": len(specials), "specials": specials})
}

func (r *Runner) cmdSpecialCreate(args []string) int {
	fs := newFlagSet("special create")
	title := fs.String("title", "", "what it is")
	body := fs.String("body", "", "what is worth noticing")
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
	special, err := ops.CreateSpecial(*title, *body, *tags)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "special": special})
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
	special, err := ops.UpdateSpecial(fs.arg(0), *title, *body, *tags)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "special": special})
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
	special, err := ops.DismissSpecial(fs.arg(0), *reason)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "special": special})
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
	note, err := ops.PromoteSpecial(fs.arg(0))
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "prepNote": note})
}

// ---------------------------------------------------------------------------
// Prep notes — rung 1
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
	notes, err := store.ListPrepNotes()
	if err != nil {
		r.fail("list prep notes", err.Error())
		return ExitError
	}
	return r.emit(map[string]any{"ok": true, "count": len(notes), "prepNotes": notes})
}

func (r *Runner) cmdPrepNoteCreate(args []string) int {
	fs := newFlagSet("prepnote create")
	title := fs.String("title", "", "what it is")
	intent := fs.String("intent", "", "what it should achieve")
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
	note, err := ops.CreatePrepNote(*title, *intent, *scope)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "prepNote": note})
}

func (r *Runner) cmdPrepNoteUpdate(args []string) int {
	fs := newFlagSet("prepnote update")
	title := fs.String("title", "", "new title")
	intent := fs.String("intent", "", "new intent")
	scope := fs.StringList("scope")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli prepnote update <id> [--title] [--intent] [--scope]")
		return ExitError
	}
	_, ops, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	note, err := ops.UpdatePrepNote(fs.arg(0), *title, *intent, *scope)
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "prepNote": note})
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
	recipe, err := ops.PromotePrepNote(fs.arg(0))
	if err != nil {
		return r.failOp(err)
	}
	return r.emit(map[string]any{"ok": true, "recipe": recipeSummary(recipe)})
}
