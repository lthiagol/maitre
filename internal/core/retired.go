package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// The plan format before the rename.
//
// The retired artifacts are **gone**: the reader does not read them, and there is no
// second writer for them. There is one format.
//
// A plan that carries a retired artifact is **not supported**, and the tool says so
// plainly rather than pretending the plan is empty. A reader that hid the retired
// artifacts would show the chef an empty backlog and call it the truth.
//
// **Not supported means no writes.** The refusal is on the write path too: a plan
// the tool does not support is not a plan the tool writes to, and a message at read
// time is not a fence.
//
// The file's name is part of the format: a file named with a retired id form is a
// retired artifact on disk, whatever the record inside it says.
//
// The journal is the one thing that keeps the old names, and it is not an exception
// to the rule: the journal is history. Its events keep their old names and their old
// ids forever, and the reader still resolves what the journal names.

// retiredDirectories are the plan directories the rename retired.
//
// The backlog was two file types in two directories — a special and a prep note — and
// it is one artifact in one directory now.
var retiredDirectories = []string{"specials", "prep-notes"}

// retiredStampField is the recipe's retired completion stamp.
//
// The stamp was `evidence`; it is `completion` now. The word is a **current** word
// too — the criterion's record is `evidence` — so the wire cannot look for the
// word. It looks for the **access**: the literal key the retired stamp is read by,
// which the current format never does.
const retiredStampField = "evidence"

// retiredProofsField is the criterion's retired record name.
//
// The record was `proofs`; it is `evidence` now.
const retiredProofsField = "proofs"

// retiredFamilies are the stems of the retired names whose own names are gone.
//
// The retired reader's functions were named `legacy*` — `ListLegacyItems`,
// `LegacyItem`, `legacyFingerprint`, `restampLegacy` — and the code that named
// them is gone. The stem is what is left of them, and it is what a new one would
// be called: a name the list spells out is a name that comes back the moment
// somebody spells it differently.
var retiredFamilies = []string{"legacy"}

// retiredIDPrefix is one retired id prefix and the one that replaced it.
type retiredIDPrefix struct {
	// retired is the prefix the rename retired.
	retired string
	// current is the prefix in use now.
	current string
}

// retiredIDForms are the id prefixes the rename retired.
//
// A recipe was `R01`; it is `RCP01` now. A backlog item was `S01` (a special) or
// `N01` (a prep note); it is `BL01` now.
var retiredIDForms = []retiredIDPrefix{
	{"R", recipeIDPrefix},
	{"S", backlogIDPrefix},
	{"N", backlogIDPrefix},
}

// retiredIDForm reports whether an id is in a retired form, and which prefix
// replaced it.
//
// The lookup stays tolerant — a name is not a record — but a *record* or a *file*
// in a retired id form is an artifact of the retired format.
func retiredIDForm(id string) (string, bool) {
	for _, form := range retiredIDForms {
		if idHasPrefix(id, form.retired) {
			return form.current, true
		}
	}
	return "", false
}

// UnsupportedPlanError reports that a plan carries artifacts the current format
// retired.
//
// The plan is not half-read and not broken by accident: it is a plan written under
// rules this tool does not speak. The error names what the plan carries, so the
// chef can see what the tool found.
type UnsupportedPlanError struct {
	// Artifacts are the retired things the plan carries, named as they are.
	Artifacts []string
}

func (e *UnsupportedPlanError) Error() string {
	return fmt.Sprintf("this plan carries artifacts of a retired format (%s): the current format is one format, and this plan is not supported",
		strings.Join(e.Artifacts, ", "))
}

// checkRetiredDirectories refuses a plan that still carries a directory the rename
// retired.
//
// An empty directory carries nothing, so it is not an artifact: a plan that has one
// is not broken, it is just untidy. A directory with something in it is a plan whose
// backlog is in the retired shape, and that plan is not supported.
func (s *Store) checkRetiredDirectories() error {
	var found []string
	for _, dir := range retiredDirectories {
		entries, err := readDirSafe(filepath.Join(s.Root, dir))
		if err != nil {
			return err
		}
		if len(entries) > 0 {
			found = append(found, dir+"/")
		}
	}
	if len(found) == 0 {
		return nil
	}
	return &UnsupportedPlanError{Artifacts: found}
}

// retiredFileNames names the retired artifacts the plan's file names carry.
//
// **The file's name is part of the format.** A recipe is `RCP01-<slug>.json` and a
// backlog item is `BL01-<slug>.json`: the name says what the file is before the
// file is read.
//
// A file named with a retired id form — `R01-<slug>.json` — is a retired artifact on
// disk, even when the record inside it is in the current format. The lookup
// normalises the name into the current record, so the plan would answer `RCP01` and
// hold a file called `R01`: a plan that says one thing in its name and another in
// its content is not one format, and the tool says so rather than reading half of
// it.
//
// The rule stops at the retired forms, because a **name** is not a record: the
// lookup takes `RCP1`, `rcp01` and `1` for the same artifact, and a file named
// `B01-<slug>.json` is the same backlog item as one named `BL01-<slug>.json`. What
// the tool does not accept is a **retired** name — and a retired id form in a file's
// name is one.
func retiredFileNames(dir, prefix string) []string {
	entries, err := readDirSafe(dir)
	if err != nil {
		return nil
	}
	var found []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		base := strings.SplitN(e.Name(), "-", 2)[0]
		if _, retired := retiredIDForm(base); !retired {
			continue
		}
		found = append(found, fmt.Sprintf("the file %q carries a retired id form (it is `%s01` now)",
			e.Name(), prefix))
	}
	return found
}

// checkRetiredLayout refuses a plan whose directories or file names carry a retired
// form.
//
// The layout is checked before any record is read: the plan's shape is the format,
// and a plan whose files are named in a retired form is not a plan this tool reads.
func (s *Store) checkRetiredLayout() error {
	if err := s.checkRetiredDirectories(); err != nil {
		return err
	}
	found := append(
		retiredFileNames(s.RecipeDir(), recipeIDPrefix),
		retiredFileNames(s.BacklogDir(), backlogIDPrefix)...)
	return retiredPlanError(found)
}

// retiredInPlan names the retired artifacts the plan's records carry.
//
// Every record is read and checked. A record that does not parse is skipped: a
// broken file is the reader's business, and this check is about the retired format.
func (s *Store) retiredInPlan() []string {
	var found []string
	dirs := []struct {
		dir    string
		recipe bool
	}{
		{s.RecipeDir(), true},
		{s.BacklogDir(), false},
	}
	for _, d := range dirs {
		entries, err := readDirSafe(d.dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
				continue
			}
			raw, err := os.ReadFile(filepath.Join(d.dir, e.Name()))
			if err != nil {
				continue
			}
			if !d.recipe {
				var v BacklogItem
				if err := json.Unmarshal(raw, &v); err != nil {
					continue
				}
				found = append(found, retiredInBacklog(raw, v.ID)...)
				continue
			}
			var r Recipe
			if err := json.Unmarshal(raw, &r); err != nil {
				continue
			}
			found = append(found, retiredInRecipe(raw, r.ID)...)
			found = append(found, retiredActorNames(&r)...)
		}
	}
	return found
}

// checkRetiredPlan refuses a plan that carries a retired artifact.
//
// **Not supported means no writes.** The reader refuses such a plan, and so does the
// writer: a plan the tool does not support is not a plan the tool writes to. A
// message at read time is not a fence — the plan would keep taking writes, and the
// writer would keep writing a format the tool does not speak.
func (s *Store) checkRetiredPlan() error {
	if err := s.checkRetiredLayout(); err != nil {
		return err
	}
	return retiredPlanError(s.retiredInPlan())
}

// retiredInRecipe names the retired artifacts a recipe file carries.
//
// Three things can be retired in a recipe: the stamp the rename moved (`evidence` at
// the recipe's own level — the criterion's records are `evidence` too, and they are
// not retired), the criterion's record under its old name (`proofs`), and the id
// form (`R01` before `RCP01`).
//
// The recipe's own `evidence` is told from a criterion's by where it sits: the
// stamp is a key of the recipe's own document, and the records live under
// `acceptance`.
func retiredInRecipe(raw []byte, id string) []string {
	var found []string
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(raw, &doc); err != nil {
		// A file that does not parse is the reader's business, not this one's.
		return nil
	}
	if _, ok := doc[retiredStampField]; ok {
		found = append(found, "a recipe's `evidence` stamp (it is `completion` now)")
	}
	if acceptance, ok := doc["acceptance"]; ok {
		var criteria []map[string]json.RawMessage
		if err := json.Unmarshal(acceptance, &criteria); err == nil {
			for i, c := range criteria {
				if _, ok := c[retiredProofsField]; ok {
					found = append(found, fmt.Sprintf("criterion %d's `proofs` (it is `evidence` now)", i+1))
				}
			}
		}
	}
	if !canonicalRecipeID(id) {
		found = append(found, fmt.Sprintf("the recipe id %q (it is `RCP01` now)", id))
	}
	return found
}

// retiredInBacklog names the retired artifacts a backlog item file carries.
//
// The old backlog had two file types, and neither had a `kind`: the kind is what
// makes the item one artifact. An item without one is an item of the retired shape.
func retiredInBacklog(raw []byte, id string) []string {
	var found []string
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil
	}
	if _, ok := doc["kind"]; !ok {
		found = append(found, fmt.Sprintf("the backlog item %q has no `kind` (it is `idea`, `feature` or `bug` now)", id))
	}
	if !canonicalBacklogID(id) {
		found = append(found, fmt.Sprintf("the backlog item id %q (it is `BL01` now)", id))
	}
	return found
}

// nameField is one field of a record that carries a name.
//
// The list of them is **one list**: the writer's view and the reader's view are
// both derived from it, so a field cannot be known by one and forgotten by the
// other.
type nameField struct {
	// what names the field, for the reader's message: "the claimant".
	what string
	// read returns the names the field carries.
	read func(*Recipe) []string
}

// nameFields are the fields of a record that carry a name.
//
// The claimant, the approver, the reviewer, the comment's author, the record's
// author, the deviation's author and the overlays' author.
//
// **One list, two views.** The writer's view (`recordNames`) checks every name the
// list yields before a write; the reader's view (`retiredNames`) reports the retired
// ones when a record is read. Both read this list and nothing else, so they cannot
// disagree about which fields carry a name — and a field known by only one of them
// was F-91's hole, one list apart: the writer accepted the name, the reader refused
// the plan, and the tool bricked the plan it had just written.
var nameFields = []nameField{
	{
		what: "the claimant",
		read: func(r *Recipe) []string { return []string{r.Execution.ClaimedBy} },
	},
	{
		what: "the approver",
		read: func(r *Recipe) []string {
			if r.Approval == nil {
				return nil
			}
			return []string{r.Approval.By}
		},
	},
	{
		what: "a review round's reviewer",
		read: func(r *Recipe) []string {
			var out []string
			for _, round := range r.Review.Rounds {
				out = append(out, round.Reviewer)
			}
			return out
		},
	},
	{
		what: "a comment's author",
		read: func(r *Recipe) []string {
			var out []string
			for _, c := range r.Review.Comments {
				out = append(out, c.Author)
			}
			return out
		},
	},
	{
		what: "a record's author",
		read: func(r *Recipe) []string {
			var out []string
			for _, c := range r.Acceptance {
				for _, e := range c.Evidence {
					out = append(out, e.By)
				}
			}
			return out
		},
	},
	{
		what: "a deviation's author",
		read: func(r *Recipe) []string {
			var out []string
			for _, d := range r.Audit.Deviations {
				out = append(out, d.By)
			}
			return out
		},
	},
	{
		what: "a block's author",
		read: func(r *Recipe) []string {
			if r.Overlays.Blocked == nil {
				return nil
			}
			return []string{r.Overlays.Blocked.By}
		},
	},
	{
		what: "a deferral's author",
		read: func(r *Recipe) []string {
			if r.Overlays.Deferred == nil {
				return nil
			}
			return []string{r.Overlays.Deferred.By}
		},
	},
	{
		what: "a cancellation's author",
		read: func(r *Recipe) []string {
			if r.Overlays.Cancelled == nil {
				return nil
			}
			return []string{r.Overlays.Cancelled.By}
		},
	},
}

// namedName is one name a record carries, and the field it came from.
type namedName struct {
	// what names the field the name came from.
	what string
	// name is the name itself.
	name string
}

// recordNames is the writer's view: every name the record carries, with the field
// it came from.
//
// It reads the one list and nothing else — which fields carry a name is the list's
// business, and the wire checks that this view keeps it that way.
func recordNames(r *Recipe) []namedName {
	var out []namedName
	for _, f := range nameFields {
		for _, name := range f.read(r) {
			out = append(out, namedName{what: f.what, name: name})
		}
	}
	return out
}

// retiredNames is the reader's view: the retired names among the record's.
//
// It is the writer's view, filtered — the reader knows the same fields, because
// both read the one list.
func retiredNames(r *Recipe) []namedName {
	var out []namedName
	for _, n := range recordNames(r) {
		if IsRetiredName(n.name) {
			out = append(out, n)
		}
	}
	return out
}

// checkRecordNames refuses a record that carries a name the tool does not accept.
//
// **A retired name cannot enter a live record, by any path.** The writer runs this
// on every record it writes, so the rule does not depend on an operation
// remembering to check: the domain's API cannot write a plan it will not read.
//
// The **journal** is the exception, and only for its **history**: the journal is
// append-only, and its old events keep their old names forever. A **new** event is
// a live record, and it is refused too.
func checkRecordNames(r *Recipe) error {
	for _, n := range recordNames(r) {
		if err := checkName(n.name); err != nil {
			return err
		}
	}
	return nil
}

// retiredActorNames names the retired names a live record carries.
//
// A name the tool does not accept cannot appear in a live record: the writer refuses
// it, and the reader refuses a record that carries one. The **journal** is the
// exception, because the journal is history.
//
// The fields are the one list's (`nameFields`): the claimant, the approver, the
// reviewer, the comment's author, the record's author, the deviation's author and
// the overlays' author.
func retiredActorNames(r *Recipe) []string {
	var found []string
	seen := map[string]bool{}
	for _, n := range retiredNames(r) {
		entry := fmt.Sprintf("%s (%q)", n.what, n.name)
		if seen[entry] {
			continue
		}
		seen[entry] = true
		found = append(found, entry)
	}
	return found
}

// canonicalRecipeID reports whether an id is in the current form: `RCP` and a number.
func canonicalRecipeID(id string) bool { return idHasPrefix(id, recipeIDPrefix) }

// canonicalBacklogID reports whether an id is in the current form: `BL` and a number.
func canonicalBacklogID(id string) bool { return idHasPrefix(id, backlogIDPrefix) }

// idHasPrefix reports whether an id is a kind's prefix followed by a number.
func idHasPrefix(id, prefix string) bool {
	if !strings.HasPrefix(id, prefix) {
		return false
	}
	rest := id[len(prefix):]
	if rest == "" {
		return false
	}
	for _, r := range rest {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// retiredPlanError builds the error for a plan that carries retired artifacts, or nil
// when it carries none.
func retiredPlanError(artifacts []string) error {
	if len(artifacts) == 0 {
		return nil
	}
	return &UnsupportedPlanError{Artifacts: artifacts}
}

// checkRetiredRecipe refuses a recipe file that carries a retired artifact.
func checkRetiredRecipe(raw []byte, r *Recipe) error {
	found := retiredInRecipe(raw, r.ID)
	found = append(found, retiredActorNames(r)...)
	return retiredPlanError(found)
}

// checkRetiredBacklog refuses a backlog item file that carries a retired artifact.
func checkRetiredBacklog(raw []byte, v *BacklogItem) error {
	return retiredPlanError(retiredInBacklog(raw, v.ID))
}
