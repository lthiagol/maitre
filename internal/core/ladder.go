package core

import (
	"fmt"
	"strings"
)

// The lower rungs of the ladder: rung 0 (special) and rung 1 (prep note).
//
// Rung 2 (recipe) and above live in model.go. These two are simple artifacts — a title
// and a little prose — so they get their own file rather than the full recipe anatomy.

const (
	// StatusOpen is a live special or prep note.
	StatusOpen = "open"
	// StatusDismissed is a special the chef chose not to pursue.
	StatusDismissed = "dismissed"
	// StatusPromoted is an artifact that became something else.
	StatusPromoted = "promoted"
)

// Special is rung 0 — a hunch, an observation, an agent discovery.
type Special struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Body       string   `json:"body,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Status     string   `json:"status"`
	Created    string   `json:"created"`
	Updated    string   `json:"updated"`
	PromotedTo string   `json:"promoted_to,omitempty"`
}

// PrepNote is rung 1 — lightweight intent and scope, before acceptance criteria exist.
type PrepNote struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Intent     string   `json:"intent"`
	Scope      []string `json:"scope,omitempty"`
	Status     string   `json:"status"`
	Created    string   `json:"created"`
	Updated    string   `json:"updated"`
	PromotedTo string   `json:"promoted_to,omitempty"`
}

// ---------------------------------------------------------------------------
// Store
// ---------------------------------------------------------------------------

// SaveSpecial writes a special atomically.
//
// The file is named `<id>-<slug>.json`, and the slug follows the title — so a rename
// writes a new file and the old one would stay. One artifact, one file: the write
// removes the copy the old name left behind.
func (s *Store) SaveSpecial(v *Special) error {
	if v.ID == "" {
		return fmt.Errorf("special has no id")
	}
	dir := filepathJoin(s.Root, DirSpecials)
	name := fmt.Sprintf("%s-%s.json", v.ID, v.Slug())
	if err := writeJSONAtomic(filepathJoin(dir, name), v); err != nil {
		return err
	}
	return removeOtherCopies(dir, name, v.ID)
}

// Slug is the special's filename slug.
func (v *Special) Slug() string { return Slugify(v.Title) }

// LoadSpecial reads one special.
func (s *Store) LoadSpecial(id string) (*Special, error) {
	raw, path, err := s.readArtifact(filepathJoin(s.Root, DirSpecials), id)
	if err != nil {
		return nil, err
	}
	var v Special
	if err := unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &v, nil
}

// ListSpecials returns every special, ordered by id.
//
// The list is never nil: an empty plan is an empty list, and the CLI marshals an
// empty list as `[]` and a nil slice as `null`. A client cannot iterate `null`.
func (s *Store) ListSpecials() ([]Special, error) {
	entries, err := readDirSafe(filepathJoin(s.Root, DirSpecials))
	if err != nil {
		return nil, err
	}
	out := make([]Special, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		raw, err := readFile(filepathJoin(s.Root, DirSpecials, e.Name()))
		if err != nil {
			return nil, err
		}
		var v Special
		if err := unmarshal(raw, &v); err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		out = append(out, v)
	}
	// One artifact, one entry: a copy left behind by a rename must not make the
	// special count twice.
	out = dedupeByID(out,
		func(v Special) string { return v.ID },
		func(v Special) string { return v.Updated },
	)
	sortByID(out, func(i int) string { return out[i].ID })
	return out, nil
}

// SavePrepNote writes a prep note atomically.
//
// The file is named `<id>-<slug>.json`, and the slug follows the title — so a rename
// writes a new file and the old one would stay. One artifact, one file: the write
// removes the copy the old name left behind.
func (s *Store) SavePrepNote(v *PrepNote) error {
	if v.ID == "" {
		return fmt.Errorf("prep note has no id")
	}
	dir := filepathJoin(s.Root, DirPrepNotes)
	name := fmt.Sprintf("%s-%s.json", v.ID, Slugify(v.Title))
	if err := writeJSONAtomic(filepathJoin(dir, name), v); err != nil {
		return err
	}
	return removeOtherCopies(dir, name, v.ID)
}

// LoadPrepNote reads one prep note.
func (s *Store) LoadPrepNote(id string) (*PrepNote, error) {
	raw, path, err := s.readArtifact(filepathJoin(s.Root, DirPrepNotes), id)
	if err != nil {
		return nil, err
	}
	var v PrepNote
	if err := unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &v, nil
}

// ListPrepNotes returns every prep note, ordered by id.
//
// The list is never nil: an empty plan is an empty list, and the CLI marshals an
// empty list as `[]` and a nil slice as `null`. A client cannot iterate `null`.
func (s *Store) ListPrepNotes() ([]PrepNote, error) {
	entries, err := readDirSafe(filepathJoin(s.Root, DirPrepNotes))
	if err != nil {
		return nil, err
	}
	out := make([]PrepNote, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		raw, err := readFile(filepathJoin(s.Root, DirPrepNotes, e.Name()))
		if err != nil {
			return nil, err
		}
		var v PrepNote
		if err := unmarshal(raw, &v); err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		out = append(out, v)
	}
	// One artifact, one entry: a copy left behind by a rename must not make the
	// prep note count twice.
	out = dedupeByID(out,
		func(v PrepNote) string { return v.ID },
		func(v PrepNote) string { return v.Updated },
	)
	sortByID(out, func(i int) string { return out[i].ID })
	return out, nil
}

// NextSpecialID allocates the next special id by scanning what exists.
//
// The prefix is `S`, as the domain model documents it (docs/03 §1): prefixes are
// distinct per kind, so a reference is unambiguous across the directories.
//
// Nothing is promoted *to* a special, so no id is reserved here.
func (s *Store) NextSpecialID() (string, error) {
	return nextIDIn(filepathJoin(s.Root, DirSpecials), specialIDPrefix)
}

// NextPrepNoteID allocates the next prep note id by scanning what exists.
//
// It also scans what is *referenced*: an id a special points at is taken, even
// though no prep note carries it. A promotion records where it went before the
// target is written, so a crash can leave a pointer to an id that does not exist —
// and that pointer must resolve to nothing, never to a stranger.
func (s *Store) NextPrepNoteID() (string, error) {
	reserved, err := s.reservedPrepNoteIDs()
	if err != nil {
		return "", err
	}
	return nextIDIn(filepathJoin(s.Root, DirPrepNotes), prepNoteIDPrefix, reserved...)
}

// reservedPrepNoteIDs are the prep note ids a special points at.
func (s *Store) reservedPrepNoteIDs() ([]string, error) {
	specials, err := s.ListSpecials()
	if err != nil {
		return nil, err
	}
	promoted := make([]string, 0, len(specials))
	for _, sp := range specials {
		promoted = append(promoted, sp.PromotedTo)
	}
	return referencedIDs(promoted, "N"), nil
}

// referencedIDs collects the ids a list of `promoted_to` values names.
//
// The value reads "<kind> <id>" — "prep-note N01", "recipe RCP01". Only the ids with
// the given prefix are collected: a pointer to a recipe does not reserve a prep note
// id.
func referencedIDs(promotedTo []string, prefix string) []string {
	var out []string
	for _, p := range promotedTo {
		fields := strings.Fields(p)
		if len(fields) != 2 {
			continue
		}
		if strings.HasPrefix(fields[1], prefix) {
			out = append(out, fields[1])
		}
	}
	return out
}

// readArtifact finds the current file in dir for the given id.
//
// A plan can hold two files for one id — a copy left behind by a rename. The current
// copy wins, so the reader never picks one arbitrarily.
func (s *Store) readArtifact(dir, id string) ([]byte, string, error) {
	paths, err := copiesWithID(dir, id)
	if err != nil {
		return nil, "", err
	}
	if len(paths) == 0 {
		return nil, "", fmt.Errorf("%s: %w", id, ErrNotFound)
	}
	path := currentCopy(paths)
	raw, err := readFile(path)
	return raw, path, err
}

// nextIDIn scans a directory for the highest numeric id with the given prefix, and
// the ids that are reserved by a pointer.
//
// A reserved id is one an artifact points at but does not carry: the pointer is
// dangling, and handing the id to a new artifact would make the trail resolve to a
// stranger.
func nextIDIn(dir, prefix string, reserved ...string) (string, error) {
	entries, err := readDirSafe(dir)
	if err != nil {
		return "", err
	}
	max := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		base := strings.SplitN(e.Name(), "-", 2)[0]
		// Every file in the directory is of this kind, and the kind's letters have
		// changed: a recipe is RCP01 now and was R01 before. The number is what
		// must not be handed out twice, so it is what is counted.
		n, ok := parseNumericID(base)
		if !ok {
			continue
		}
		if n > max {
			max = n
		}
	}
	for _, id := range reserved {
		n, ok := parseNumericID(id)
		if ok && n > max {
			max = n
		}
	}
	return fmt.Sprintf("%s%02d", prefix, max+1), nil
}

// ---------------------------------------------------------------------------
// Operations
// ---------------------------------------------------------------------------

// CreateSpecial files a special.
//
// The lock is taken before the id is allocated: two creations must not be handed the
// same id, and the id is a scan of what exists.
func (o *Operations) CreateSpecial(title, body string, tags []string) (*Special, error) {
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("a special needs a title")
	}
	unlock, err := o.store.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	id, err := o.store.NextSpecialID()
	if err != nil {
		return nil, err
	}
	now := Now()
	v := &Special{
		ID: id, Title: title, Body: body, Tags: tags,
		Status: StatusOpen, Created: now, Updated: now,
	}
	if err := o.store.SaveSpecial(v); err != nil {
		return nil, err
	}
	if err := o.journalLocked(KindCreated, "", map[string]any{"kind": "special", "id": id, "title": title}); err != nil {
		return nil, err
	}
	return v, nil
}

// UpdateSpecial edits a special.
func (o *Operations) UpdateSpecial(id, title, body string, tags []string) (*Special, error) {
	unlock, err := o.store.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	v, err := o.store.LoadSpecial(id)
	if err != nil {
		return nil, err
	}
	if v.Status != StatusOpen {
		return nil, o.refuseLocked(id, fmt.Sprintf("special %s is %s", id, v.Status))
	}
	if strings.TrimSpace(title) != "" {
		v.Title = title
	}
	v.Body = body
	if tags != nil {
		v.Tags = tags
	}
	v.Updated = Now()
	if err := o.store.SaveSpecial(v); err != nil {
		return nil, err
	}
	if err := o.journalLocked(KindDecision, "", map[string]any{"kind": "special", "id": id}); err != nil {
		return nil, err
	}
	return v, nil
}

// DismissSpecial marks a special as not worth pursuing.
func (o *Operations) DismissSpecial(id, reason string) (*Special, error) {
	unlock, err := o.store.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	v, err := o.store.LoadSpecial(id)
	if err != nil {
		return nil, err
	}
	if v.Status != StatusOpen {
		return nil, o.refuseLocked(id, fmt.Sprintf("special %s is already %s", id, v.Status))
	}
	v.Status = StatusDismissed
	v.Updated = Now()
	v.Body = strings.TrimSpace(v.Body + "\n\nDismissed: " + reason)
	if err := o.store.SaveSpecial(v); err != nil {
		return nil, err
	}
	if err := o.journalLocked(KindOverlay, "", map[string]any{"kind": "special", "id": id, "status": v.Status}); err != nil {
		return nil, err
	}
	return v, nil
}

// PromoteSpecial turns a special into a prep note.
//
// Promotion is one-way: the special records where it went, so the trail survives.
//
// The source is written first. A crash between the two writes then leaves a special
// that says where it went and cannot be promoted again — one special, one prep note.
// The other order would leave the special open, and the second promotion would make
// two prep notes out of one hunch.
//
// An *error* is different: if the target cannot be written, the source is rolled
// back, so the failure leaves the plan as it was — no status, no pointer, no event.
func (o *Operations) PromoteSpecial(id string) (*PrepNote, error) {
	unlock, err := o.store.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	v, err := o.store.LoadSpecial(id)
	if err != nil {
		return nil, err
	}
	if v.Status != StatusOpen {
		return nil, o.refuseLocked(id, fmt.Sprintf("special %s is %s", id, v.Status))
	}
	noteID, err := o.store.NextPrepNoteID()
	if err != nil {
		return nil, err
	}
	now := Now()
	note := &PrepNote{
		ID:      noteID,
		Title:   v.Title,
		Intent:  v.Body,
		Status:  StatusOpen,
		Created: now,
		Updated: now,
	}

	// The source first: it records the target, so the promotion cannot happen twice.
	wasStatus, wasPromotedTo, wasUpdated := v.Status, v.PromotedTo, v.Updated
	v.Status = StatusPromoted
	v.PromotedTo = "prep-note " + noteID
	v.Updated = now
	if err := o.store.SaveSpecial(v); err != nil {
		return nil, err
	}
	// Then the new artifact.
	if err := o.store.SavePrepNote(note); err != nil {
		// The target could not be written. The source must not claim it was
		// promoted: a pointer to an artifact that does not exist is a lie, and a
		// lie that later resolves to a stranger is worse.
		v.Status, v.PromotedTo, v.Updated = wasStatus, wasPromotedTo, wasUpdated
		if rbErr := o.store.SaveSpecial(v); rbErr != nil {
			return nil, fmt.Errorf("promote special %s: %w (the source could not be rolled back: %v)", id, err, rbErr)
		}
		return nil, err
	}
	if err := o.journalLocked(KindPromotion, "", map[string]any{
		"kind": "special", "id": id, "promoted_to": note.PromotedTo,
	}); err != nil {
		return nil, err
	}
	return note, nil
}

// CreatePrepNote files a prep note.
//
// The lock is taken before the id is allocated: two creations must not be handed the
// same id, and the id is a scan of what exists.
func (o *Operations) CreatePrepNote(title, intent string, scope []string) (*PrepNote, error) {
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("a prep note needs a title")
	}
	unlock, err := o.store.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	id, err := o.store.NextPrepNoteID()
	if err != nil {
		return nil, err
	}
	now := Now()
	v := &PrepNote{
		ID: id, Title: title, Intent: intent, Scope: scope,
		Status: StatusOpen, Created: now, Updated: now,
	}
	if err := o.store.SavePrepNote(v); err != nil {
		return nil, err
	}
	if err := o.journalLocked(KindCreated, "", map[string]any{"kind": "prep-note", "id": id, "title": title}); err != nil {
		return nil, err
	}
	return v, nil
}

// UpdatePrepNote edits a prep note.
//
// It is an operation, not a file write: the write takes the plan lock and leaves a
// journal entry, because the journal is the only event log. A note that is no longer
// open is not edited — promotion is one-way.
func (o *Operations) UpdatePrepNote(id, title, intent string, scope []string) (*PrepNote, error) {
	unlock, err := o.store.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	v, err := o.store.LoadPrepNote(id)
	if err != nil {
		return nil, err
	}
	if v.Status != StatusOpen {
		return nil, o.refuseLocked(id, fmt.Sprintf("prep note %s is %s", id, v.Status))
	}
	if strings.TrimSpace(title) != "" {
		v.Title = title
	}
	if strings.TrimSpace(intent) != "" {
		v.Intent = intent
	}
	if len(scope) > 0 {
		v.Scope = scope
	}
	v.Updated = Now()
	if err := o.store.SavePrepNote(v); err != nil {
		return nil, err
	}
	if err := o.journalLocked(KindDecision, "", map[string]any{"kind": "prep-note", "id": id}); err != nil {
		return nil, err
	}
	return v, nil
}

// PromotePrepNote turns a prep note into a recipe.
//
// The recipe starts in draft: the intent and scope carry over, and acceptance criteria
// are the next thing to write.
//
// The source is written first, with the target it went to, so a promotion cannot
// happen twice. An *error* rolls the source back: the failure leaves the plan as it
// was — no status, no pointer, no event.
func (o *Operations) PromotePrepNote(id string) (*Recipe, error) {
	unlock, err := o.store.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	v, err := o.store.LoadPrepNote(id)
	if err != nil {
		return nil, err
	}
	if v.Status != StatusOpen {
		return nil, o.refuseLocked(id, fmt.Sprintf("prep note %s is %s", id, v.Status))
	}
	recipeID, err := o.store.NextRecipeID()
	if err != nil {
		return nil, err
	}

	// The source first: it records the target, so the promotion cannot happen twice.
	wasStatus, wasPromotedTo, wasUpdated := v.Status, v.PromotedTo, v.Updated
	v.Status = StatusPromoted
	v.PromotedTo = "recipe " + recipeID
	v.Updated = Now()
	if err := o.store.SavePrepNote(v); err != nil {
		return nil, err
	}
	// Then the new artifact.
	recipe, err := o.createRecipeLocked(recipeID, RecipeSpec{
		Title:   v.Title,
		Outcome: v.Intent,
		InScope: v.Scope,
	})
	if err != nil {
		// The target could not be written. The source must not claim it was
		// promoted: a pointer to an artifact that does not exist is a lie, and a
		// lie that later resolves to a stranger is worse.
		v.Status, v.PromotedTo, v.Updated = wasStatus, wasPromotedTo, wasUpdated
		if rbErr := o.store.SavePrepNote(v); rbErr != nil {
			return nil, fmt.Errorf("promote prep note %s: %w (the source could not be rolled back: %v)", id, err, rbErr)
		}
		return nil, err
	}
	if err := o.journalLocked(KindPromotion, recipe.ID, map[string]any{
		"kind": "prep-note", "id": id, "promoted_to": v.PromotedTo,
	}); err != nil {
		return nil, err
	}
	return recipe, nil
}

// ListSpecials returns every special.
func (o *Operations) ListSpecials() ([]Special, error) {
	return o.store.ListSpecials()
}

// ListPrepNotes returns every prep note.
func (o *Operations) ListPrepNotes() ([]PrepNote, error) {
	return o.store.ListPrepNotes()
}
