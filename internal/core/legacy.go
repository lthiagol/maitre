package core

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// The plan format before the records were renamed.
//
// The recipe's completion stamp was `evidence`; it is `completion` now. The stamp
// says when the recipe was done, and the per-criterion records say what was
// verified — those live on the criterion.
//
// The reader **tolerates** the old names, because a plan written before the rename
// keeps what it had, and a reader that drops a record is a reader that lies about
// the plan it is showing. The writer is the one that changed: a plan the tool
// writes again carries the new names.

// The plan format before the backlog became one artifact.
//
// The backlog was two file types — a *special* in `kitchen/specials/` and a *prep
// note* in `kitchen/prep-notes/`. The format has no such files any more: the writer
// does not create them and `Init` does not make the directories.
//
// The reader **tolerates** them, because a plan written before the rename keeps
// what it had, and a console that hides those files is a console that lies about
// the plan it is showing. The console shows them, in their own group, and says
// where they came from.

// The rungs a legacy item came from.
//
// They are names for the old files, not kinds: a legacy special has no `kind` —
// the field did not exist then.
const (
	// LegacySpecial is a special: a hunch, an observation, an agent discovery.
	LegacySpecial = "special"
	// LegacyPrepNote is a prep note: intent and scope, before acceptance criteria.
	LegacyPrepNote = "prep note"
)

// LegacyItem is an artifact from the plan format before the backlog became one.
//
// It is not a backlog item: it has no `kind`, and it is not in the plan format
// any more. It is read so the console can show a plan written before the rename
// without hiding half of it.
type LegacyItem struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Body   string   `json:"body,omitempty"`
	Intent string   `json:"intent,omitempty"`
	Scope  []string `json:"scope,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Status string   `json:"status"`
	// PromotedTo is where the item went, when it was promoted.
	//
	// The old format kept the pointer, and the pointer is a promise: it says an
	// id was taken. A plan that is read but never written still has to honour
	// the promise — the id must not be handed to a stranger.
	PromotedTo string `json:"promoted_to,omitempty"`
	// Rung is where the item came from — a special or a prep note. It is read
	// from the directory the file lives in, not from the file.
	Rung string `json:"-"`
}

// Text is the item's prose: a special's body, a prep note's intent.
//
// The two old artifacts kept their prose in different fields, and the console
// shows one thing. The field is the only difference between them.
func (v LegacyItem) Text() string {
	if v.Body != "" {
		return v.Body
	}
	return v.Intent
}

// ListLegacyItems returns the artifacts of the plan format before the backlog
// became one.
//
// The list is never nil, and it is empty for a plan written by this tool: the
// directories are not there, and a missing directory is an empty one.
func (s *Store) ListLegacyItems() ([]LegacyItem, error) {
	out := []LegacyItem{}
	for _, dir := range []struct {
		name string
		rung string
	}{
		{"specials", LegacySpecial},
		{"prep-notes", LegacyPrepNote},
	} {
		items, err := s.readLegacyDir(dir.name, dir.rung)
		if err != nil {
			return nil, err
		}
		out = append(out, items...)
	}
	return out, nil
}

// readLegacyDir reads one legacy directory.
func (s *Store) readLegacyDir(name, rung string) ([]LegacyItem, error) {
	dir := filepath.Join(s.Root, name)
	entries, err := readDirSafe(dir)
	if err != nil {
		return nil, err
	}
	out := make([]LegacyItem, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		raw, err := readFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		var v LegacyItem
		if err := unmarshal(raw, &v); err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		if v.ID == "" {
			// The id is the file's name: the old writers put it there too, and a
			// file without one names nothing.
			v.ID = strings.SplitN(e.Name(), "-", 2)[0]
		}
		v.Rung = rung
		out = append(out, v)
	}
	sortByID(out, func(i int) string { return out[i].ID })
	return out, nil
}

// readLegacyFields fills the fields a recipe wrote under the old names.
//
// The recipe's completion stamp was `evidence`; it is `completion` now. A plan
// written before the rename carries the old name, and a reader that ignores it
// loses the moment the recipe was done.
//
// The new name wins when a file carries both: the reader does not merge two
// answers to one question, and the new one is the one the tool writes.
func readLegacyFields(raw []byte, r *Recipe) {
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(raw, &doc); err != nil {
		return
	}
	if _, ok := doc["completion"]; ok {
		return
	}
	old, ok := doc["evidence"]
	if !ok {
		return
	}
	var c Completion
	if err := json.Unmarshal(old, &c); err != nil {
		return
	}
	r.Completion = c
}
