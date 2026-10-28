package core

import (
	"errors"
	"strings"
)

// RecordMigration appends the plan's migration to its journal.
//
// A plan's records move once — the ids are renamed, a field is renamed, the digest's
// encoding changes — and the journal does not say *why*, so a reader that finds
// `R01` in the journal has to guess whether it is an old id, a typo, or a different
// plan.
//
// The migration says it: the form the plan came from and the form it went to, in one
// event.
//
// **It records a fact; it does not convert anything.** There is no upgrade path: the
// plan's live records are written in the current format, and a plan that carries a
// retired artifact is not supported. What this operation adds is the *history* —
// the journal is append-only, so the events written before the move keep the names
// they were written with, and the move is *added* after them.
//
// The event carries the current domain version, which is also how the plan's stamp
// moves. The version guard reads the newest entry: a plan whose newest entry says
// nothing about the version reads as version zero, and nothing fences a stale tool
// off from it. Appending the migration stamps the plan with the version this build
// speaks, so the next stale tool is refused.
//
// The event names no recipe: the migration is about the plan, not about one recipe.
//
// **The declaration crosses the version fence.** It is the one write that does, and
// it is deliberate: the plan's owner says which rules the plan lives under. It
// changes no record — it appends one event to the journal — and it may only move
// the plan **forward**, so a stale tool cannot use it to unlock a plan it does not
// understand. It is not reachable from the CLI.
func (o *Operations) RecordMigration(from, to string) error {
	if strings.TrimSpace(from) == "" {
		return errors.New("a migration needs the id form it came from")
	}
	if strings.TrimSpace(to) == "" {
		return errors.New("a migration needs the id form it goes to")
	}
	if from == to {
		return errors.New("a migration must change the id form")
	}
	unlock, err := o.store.lockUnder(o.store.checkDeclarationForward)
	if err != nil {
		return err
	}
	defer unlock()
	return o.journalLocked(KindControl, "", map[string]any{
		"kind": "migration",
		"from": from,
		"to":   to,
	})
}
