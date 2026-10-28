package core

import (
	"encoding/json"
	"fmt"
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
// The journal is the one thing that keeps the old names, and it is not an exception
// to the rule: the journal is history. Its events keep their old names and their old
// ids forever, and the reader still resolves what the journal names.

// retiredDirectories are the plan directories the rename retired.
//
// The backlog was two file types in two directories — a special and a prep note — and
// it is one artifact in one directory now.
var retiredDirectories = []string{"specials", "prep-notes"}

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
	if _, ok := doc["evidence"]; ok {
		found = append(found, "a recipe's `evidence` stamp (it is `completion` now)")
	}
	if acceptance, ok := doc["acceptance"]; ok {
		var criteria []map[string]json.RawMessage
		if err := json.Unmarshal(acceptance, &criteria); err == nil {
			for i, c := range criteria {
				if _, ok := c["proofs"]; ok {
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

// retiredActorNames names the retired names a live record carries.
//
// A name the tool does not accept cannot appear in a live record: the writer refuses
// it, and the reader refuses a record that carries one. The **journal** is the
// exception, because the journal is history.
//
// The fields are the ones that name an actor: the claimant, the approver, the
// reviewer, the comment's author, the record's author, the deviation's author and
// the overlays' author.
func retiredActorNames(r *Recipe) []string {
	var found []string
	seen := map[string]bool{}
	note := func(what, name string) {
		if !IsRetiredName(name) {
			return
		}
		entry := fmt.Sprintf("%s (%q)", what, name)
		if seen[entry] {
			return
		}
		seen[entry] = true
		found = append(found, entry)
	}
	note("the claimant", r.Execution.ClaimedBy)
	if r.Approval != nil {
		note("the approver", r.Approval.By)
	}
	for _, round := range r.Review.Rounds {
		note("a review round's reviewer", round.Reviewer)
	}
	for _, c := range r.Review.Comments {
		note("a comment's author", c.Author)
	}
	for _, c := range r.Acceptance {
		for _, e := range c.Evidence {
			note("a record's author", e.By)
		}
	}
	for _, d := range r.Audit.Deviations {
		note("a deviation's author", d.By)
	}
	if r.Overlays.Blocked != nil {
		note("a block's author", r.Overlays.Blocked.By)
	}
	if r.Overlays.Deferred != nil {
		note("a deferral's author", r.Overlays.Deferred.By)
	}
	if r.Overlays.Cancelled != nil {
		note("a cancellation's author", r.Overlays.Cancelled.By)
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
