package core

import (
	"fmt"
	"strings"
)

// The backlog: one artifact, with a kind.
//
// The backlog had two file types and two CLI verb groups. They are one artifact now:
// a **backlog item** whose `kind` says what it is.
//
// The kind is the only truth. The console's two groups — *Ideas* and *Small tasks* —
// are derived from it and never stored. The id is kind-neutral: `BL01` stays `BL01`
// when an idea is shaped into a small task, so the trail survives the shaping.
//
// See docs/07-vocabulary.md § The backlog.

// The statuses a backlog item carries.
//
// A closed set, for the same reason the journal's kinds are: the item's status says
// whether it is live, and a typo would leave an item in a state nothing can move it
// out of.
const (
	// StatusOpen is a live backlog item.
	StatusOpen = "open"
	// StatusDismissed is an item the chef chose not to pursue.
	StatusDismissed = "dismissed"
	// StatusPromoted is an item that became something else.
	StatusPromoted = "promoted"
)

// The kinds a backlog item carries.
//
// A closed set: the kind decides the group the console files the item under, and the
// shape a promotion gives the recipe.
const (
	// KindIdea is a hunch — worth knowing, not yet work.
	KindIdea = "idea"
	// KindFeature is a small task — the work is known and it is ready to run.
	KindFeature = "feature"
	// KindBug is a small task that is a defect.
	KindBug = "bug"
)

// BacklogKinds is every kind a backlog item may carry.
var BacklogKinds = []string{KindIdea, KindFeature, KindBug}

// BacklogItem is the backlog's artifact.
//
// The body is why the item is worth doing — a hunch, or the intent behind a small
// task. The scope is the work, if it is known.
type BacklogItem struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Body       string   `json:"body,omitempty"`
	Scope      []string `json:"scope,omitempty"`
	Kind       string   `json:"kind"`
	Tags       []string `json:"tags,omitempty"`
	Status     string   `json:"status"`
	Created    string   `json:"created"`
	Updated    string   `json:"updated"`
	PromotedTo string   `json:"promoted_to,omitempty"`
}

// isSmallTask reports whether a kind is a small task: the work is known, and it is
// ready to run.
//
// The rule lives here, once, because two surfaces need it — the promotion, which
// gives a small task a criterion, and the console, which files a small task under
// *Small tasks*. Two answers to one question is how the console and the domain drift
// apart.
func isSmallTask(kind string) bool {
	return kind == KindFeature || kind == KindBug
}

// BacklogGroups splits the backlog into the two groups the console files it under.
//
// The groups are derived from the kind — the kind is the only truth, and a group is
// never stored. An idea needs shaping; a feature or a bug is a small task, ready to
// work.
//
// An item whose kind is not one the domain knows is filed under *Ideas*: it needs
// shaping, and the row shows the kind it actually carries, so nothing is hidden.
//
// view
func BacklogGroups(items []BacklogItem) (ideas, tasks []BacklogItem) {
	for _, it := range items {
		if isSmallTask(it.Kind) {
			tasks = append(tasks, it)
			continue
		}
		ideas = append(ideas, it)
	}
	return ideas, tasks
}

// ---------------------------------------------------------------------------
// Store
// ---------------------------------------------------------------------------

// SaveBacklogItem writes a backlog item atomically.
//
// The file is named `<id>-<slug>.json`, and the slug follows the title — so a rename
// writes a new file and the old one would stay. One artifact, one file: the write
// removes the copy the old name left behind.
func (s *Store) SaveBacklogItem(v *BacklogItem) error {
	if v.ID == "" {
		return fmt.Errorf("backlog item has no id")
	}
	dir := filepathJoin(s.Root, DirBacklog)
	name := fmt.Sprintf("%s-%s.json", v.ID, v.Slug())
	if err := writeJSONAtomic(filepathJoin(dir, name), v); err != nil {
		return err
	}
	return removeOtherCopies(dir, name, v.ID)
}

// Slug is the backlog item's filename slug.
func (v *BacklogItem) Slug() string { return Slugify(v.Title) }

// LoadBacklogItem reads one backlog item.
//
// The id is tolerant: `BL01`, `bl01`, `BL1`, `B01` and `1` all name the same item —
// a lookup resolves a name, and a name is not a record.
//
// The record itself is current-format only: an item without a `kind`, or with a
// retired id form, is refused.
func (s *Store) LoadBacklogItem(id string) (*BacklogItem, error) {
	if err := s.checkRetiredLayout(); err != nil {
		return nil, err
	}
	raw, path, err := s.readArtifact(filepathJoin(s.Root, DirBacklog), id)
	if err != nil {
		return nil, err
	}
	var v BacklogItem
	if err := unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := checkRetiredBacklog(raw, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// ListBacklogItems returns every backlog item, ordered by id.
//
// The list is never nil: an empty plan is an empty list, and the CLI marshals an
// empty list as `[]` and a nil slice as `null`. A client cannot iterate `null`.
func (s *Store) ListBacklogItems() ([]BacklogItem, error) {
	if err := s.checkRetiredLayout(); err != nil {
		return nil, err
	}
	entries, err := readDirSafe(filepathJoin(s.Root, DirBacklog))
	if err != nil {
		return nil, err
	}
	out := make([]BacklogItem, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		raw, err := readFile(filepathJoin(s.Root, DirBacklog, e.Name()))
		if err != nil {
			return nil, err
		}
		var v BacklogItem
		if err := unmarshal(raw, &v); err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		if err := checkRetiredBacklog(raw, &v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	// One artifact, one entry: a copy left behind by a rename must not make the
	// item count twice.
	out = dedupeByID(out,
		func(v BacklogItem) string { return v.ID },
		func(v BacklogItem) string { return v.Updated },
	)
	sortByID(out, func(i int) string { return out[i].ID })
	return out, nil
}

// NextBacklogItemID allocates the next backlog item id by scanning what exists.
//
// The prefix is `BL` — two letters, like the kind-neutral id the vocabulary asks for.
// The number is what must not be handed out twice, and the number is what the scan
// counts, so an item written before the rename does not collide with one written
// after it.
//
// Nothing is promoted *to* a backlog item, so no id is reserved here.
func (s *Store) NextBacklogItemID() (string, error) {
	return nextIDIn(filepathJoin(s.Root, DirBacklog), backlogIDPrefix)
}

// ---------------------------------------------------------------------------
// Operations
// ---------------------------------------------------------------------------

// CreateBacklogItem files a backlog item.
//
// The lock is taken before the id is allocated: two creations must not be handed the
// same id, and the id is a scan of what exists.
func (o *Operations) CreateBacklogItem(title, body string, scope []string, kind string, tags []string) (*BacklogItem, error) {
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("a backlog item needs a title")
	}
	if kind == "" {
		kind = KindIdea
	}
	if !isKnownKind(kind) {
		return nil, fmt.Errorf("unknown kind %q (want one of %s)", kind, strings.Join(BacklogKinds, ", "))
	}
	unlock, err := o.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	id, err := o.store.NextBacklogItemID()
	if err != nil {
		return nil, err
	}
	now := Now()
	v := &BacklogItem{
		ID: id, Title: title, Body: body, Scope: emptyIfNil(scope), Kind: kind, Tags: emptyIfNil(tags),
		Status: StatusOpen, Created: now, Updated: now,
	}
	if err := o.store.SaveBacklogItem(v); err != nil {
		return nil, err
	}
	if err := o.journalLocked(KindCreated, "", map[string]any{"kind": "backlog", "id": id, "title": title, "item_kind": kind}); err != nil {
		return nil, err
	}
	return v, nil
}

// UpdateBacklogItem edits a backlog item.
//
// Changing the kind is the **shaping** step: an idea becomes a small task by being
// made a `feature` or a `bug` and gaining a scope. It is an update, not a promotion
// — the item stays one item, and its id does not change.
//
// A field that is empty is left alone: an update says what changes, not what the
// item is. A **list** that was not given is nil, and nil is left alone too: the two
// list fields obey one rule — a list that is given replaces the list it names, a
// list that was not given leaves its field alone.
func (o *Operations) UpdateBacklogItem(id, title, body string, scope []string, kind string, tags []string) (*BacklogItem, error) {
	if kind != "" && !isKnownKind(kind) {
		return nil, fmt.Errorf("unknown kind %q (want one of %s)", kind, strings.Join(BacklogKinds, ", "))
	}
	unlock, err := o.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	v, err := o.store.LoadBacklogItem(id)
	if err != nil {
		return nil, err
	}
	if v.Status != StatusOpen {
		return nil, o.refuseLocked(id, fmt.Sprintf("backlog item %s is %s", id, v.Status))
	}
	if strings.TrimSpace(title) != "" {
		v.Title = title
	}
	if strings.TrimSpace(body) != "" {
		v.Body = body
	}
	if kind != "" {
		v.Kind = kind
	}
	setList(&v.Scope, scope)
	setList(&v.Tags, tags)
	v.Scope = emptyIfNil(v.Scope)
	v.Tags = emptyIfNil(v.Tags)
	v.Updated = Now()
	if err := o.store.SaveBacklogItem(v); err != nil {
		return nil, err
	}
	if err := o.journalLocked(KindDecision, "", map[string]any{"kind": "backlog", "id": id, "item_kind": v.Kind}); err != nil {
		return nil, err
	}
	return v, nil
}

// setList replaces a list field when the list was given, and leaves it alone when it
// was not.
//
// **One rule for every list field of the operation.** nil means "not given"; an
// empty list means "replace it with nothing". Two guards with two different rules
// is how the tags came to be wiped by an update that never mentioned them.
func setList(dst *[]string, given []string) {
	if given != nil {
		*dst = given
	}
}

// emptyIfNil makes a list a list: never nil.
//
// The plan's own rule is that a list is a list — an empty one is `[]`, never `null`.
// The JSON edge honours it: a nil list would be `null` in a document, and the
// standard readers break on that.
func emptyIfNil(list []string) []string {
	if list == nil {
		return []string{}
	}
	return list
}

// DismissBacklogItem marks a backlog item as not worth pursuing.
//
// Every kind can be dismissed: an idea, a feature and a bug are all things the chef
// can decide against, and the asymmetry of the old ladder — only one kind could be
// dismissed — had no reason behind it.
func (o *Operations) DismissBacklogItem(id, reason string) (*BacklogItem, error) {
	unlock, err := o.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	v, err := o.store.LoadBacklogItem(id)
	if err != nil {
		return nil, err
	}
	if v.Status != StatusOpen {
		return nil, o.refuseLocked(id, fmt.Sprintf("backlog item %s is already %s", id, v.Status))
	}
	v.Status = StatusDismissed
	v.Updated = Now()
	v.Body = strings.TrimSpace(v.Body + "\n\nDismissed: " + reason)
	if err := o.store.SaveBacklogItem(v); err != nil {
		return nil, err
	}
	if err := o.journalLocked(KindOverlay, "", map[string]any{"kind": "backlog", "id": id, "status": v.Status}); err != nil {
		return nil, err
	}
	return v, nil
}

// PromoteBacklogItem turns a backlog item into a recipe.
//
// The recipe starts in **draft**: it carries the item's intent, its scope and — for a
// small task — one acceptance criterion, and nothing else. The recipe is not
// pre-groomed: the out-of-scope items, the falsifiable verification and the steps
// are the recipe's own life, exactly as for any recipe. **No gate is weakened** — a
// criterion with no verification cannot pass the falsifiable gate, and a recipe with
// no out-of-scope items cannot be groomed.
//
// Promotion is one-way: the item records where it went, so the trail survives.
//
// The source is written first. A crash between the two writes then leaves an item
// that says where it went and cannot be promoted again — one item, one recipe. The
// other order would leave the item open, and the second promotion would make two
// recipes out of one hunch.
//
// An *error* is different: if the target cannot be written, the source is rolled
// back, so the failure leaves the plan as it was — no status, no pointer, no event.
func (o *Operations) PromoteBacklogItem(id string) (*Recipe, error) {
	unlock, err := o.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	v, err := o.store.LoadBacklogItem(id)
	if err != nil {
		return nil, err
	}
	if v.Status != StatusOpen {
		return nil, o.refuseLocked(id, fmt.Sprintf("backlog item %s is %s", id, v.Status))
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
	if err := o.store.SaveBacklogItem(v); err != nil {
		return nil, err
	}
	// Then the new artifact.
	recipe, err := o.createRecipeLocked(recipeID, RecipeSpec{
		Title:    v.Title,
		Outcome:  v.Title,
		Problem:  v.Body,
		InScope:  v.Scope,
		Criteria: o.promotedCriteria(v),
	})
	if err != nil {
		// The target could not be written. The source must not claim it was
		// promoted: a pointer to an artifact that does not exist is a lie, and a
		// lie that later resolves to a stranger is worse.
		v.Status, v.PromotedTo, v.Updated = wasStatus, wasPromotedTo, wasUpdated
		if rbErr := o.store.SaveBacklogItem(v); rbErr != nil {
			return nil, fmt.Errorf("promote backlog item %s: %w (the source could not be rolled back: %v)", id, err, rbErr)
		}
		return nil, err
	}
	if err := o.journalLocked(KindPromotion, recipe.ID, map[string]any{
		"kind": "backlog", "id": id, "promoted_to": v.PromotedTo,
	}); err != nil {
		return nil, err
	}
	return recipe, nil
}

// promotedCriteria are the criteria a promoted recipe starts with.
//
// A **small task** becomes a recipe with exactly one criterion: the work is small
// enough to be one thing, and a recipe with a criterion is falsifiable from the
// start. Its verification is the recipe's own life — a criterion with no
// verification cannot pass the falsifiable gate, so nothing is claimed that is not
// proven.
//
// An **idea** becomes a recipe with no criteria: it is a big thing, and the criteria
// are what the recipe is groomed into.
func (o *Operations) promotedCriteria(v *BacklogItem) []CriterionSpec {
	if !isSmallTask(v.Kind) {
		return nil
	}
	return []CriterionSpec{{Description: v.Title}}
}

// ListBacklogItems returns every backlog item.
func (o *Operations) ListBacklogItems() ([]BacklogItem, error) {
	return o.store.ListBacklogItems()
}

// isKnownKind reports whether a kind is one the domain knows.
func isKnownKind(kind string) bool {
	for _, k := range BacklogKinds {
		if k == kind {
			return true
		}
	}
	return false
}
