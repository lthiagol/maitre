package core

import (
	"errors"
	"strings"
)

// RecordMigration appends the plan's id migration to its journal.
//
// A plan written before the rename carries `R01`; a plan written after it carries
// `RCP01`. The reader takes both, so the plan still reads — but the journal does not
// say *why*, and a reader that finds `R01` in the journal has to guess whether it is
// an old id, a typo, or a different plan.
//
// The migration says it: the old id form and the new one, in one event.
//
// The journal is the plan's only event log and it is **append-only**: the history is
// never rewritten. So a plan whose ids were renamed cannot have its old events edited
// to the new ids — the migration is *added*, after them, and the history stays what
// it was.
//
// The event carries the current domain version, which is also how the plan's stamp
// moves. The version guard reads the newest entry: a plan whose newest entry says
// nothing about the version reads as version zero, and nothing fences a stale tool
// off from it. Appending the migration stamps the plan with the version this build
// speaks, so the next stale tool is refused.
//
// The event names no recipe: the migration is about the plan, not about one recipe.
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
	return o.journal(KindControl, "", map[string]any{
		"kind": "migration",
		"from": from,
		"to":   to,
	})
}
