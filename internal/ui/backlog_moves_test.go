package ui

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
)

// TestTheFooterNamesTheMoveTheDomainGives covers the footer on every state of an
// item.
//
// The console asks the domain what the chef can do with the selected item, and the
// footer reads that answer: `space shape it` on an open idea, `space promote to
// recipe` on an open small task, and **no** `space` on a settled item.
//
// The settled case is the one that holds the guard. A promoted item is still a
// small task — the kind says so — so a console that decided by the kind would offer
// the promotion, and the domain refuses it. The console asks instead.
func TestTheFooterNamesTheMoveTheDomainGives(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog

	// The first idea.
	m.selBacklog = 0
	out := plain(m)
	if !strings.Contains(out, "space shape it") {
		t.Errorf("the footer should offer the idea's move:\n%s", out)
	}
	if strings.Contains(out, "space promote to recipe") {
		t.Errorf("the footer offers a promotion the domain refuses:\n%s", out)
	}

	// The small task.
	m.selBacklog = 4
	out = plain(m)
	if !strings.Contains(out, "space promote to recipe") {
		t.Errorf("the footer does not offer the small task's promotion:\n%s", out)
	}

	// The promoted item: history, and no move at all.
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the promotion was refused: %s", m.refusal)
	}
	m.Update(m.reload()())
	// The cursor leaves the promoted item. The footer's claim is about that
	// item, so the check sits on it.
	m.selBacklog = 4
	out = plain(m)
	if strings.Contains(out, "space promote to recipe") || strings.Contains(out, "space shape it") {
		t.Errorf("a settled item has no move — the footer must not offer one:\n%s", out)
	}
}

// TestTheFooterAsksTheDomainForTheMove holds the console's guard.
//
// A console that decided for itself would pass on the two cases where its own rule
// and the domain's answer **agree**: an idea and an open small task. Those two hold
// nothing. The case where they **differ** is a promoted small task — still a small
// task by kind, no move by the domain's answer — and that is the case this holds.
func TestTheFooterAsksTheDomainForTheMove(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 4
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the promotion was refused: %s", m.refusal)
	}
	m.Update(m.reload()())
	// The cursor leaves the promoted item. The footer's claim is about that
	// item, so the check sits on it.
	m.selBacklog = 4

	// The domain's answer, read directly: no move.
	item, err := m.store.LoadBacklogItem("BL05")
	if err != nil {
		t.Fatal(err)
	}
	if move := core.BacklogMoveOf(*item); move.Primary != core.BacklogMoveNone {
		t.Fatalf("a promoted item's move is %q, want none", move.Primary)
	}
	// And the footer reads it.
	if out := plain(m); strings.Contains(out, "space promote to recipe") {
		t.Errorf("the console decided for itself — the domain says there is no move\n%s", out)
	}
}

// TestShapingAnIdeaChangesItsKindAndScopeAndKeepsItsId covers the shaping.
//
// The kind changes and a scope is added — and the id does not: `BL01` stays
// `BL01`, because the id is kind-neutral and the trail survives the shaping.
func TestShapingAnIdeaChangesItsKindAndScopeAndKeepsItsId(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	before, err := m.ops.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}

	steps := append([]string{"space", "f"}, typed("the evidence lane")...)
	steps = append(steps, "enter")
	if _, err := Drive(m, steps); err != nil {
		t.Fatal(err)
	}
	if m.refusal != "" {
		t.Fatalf("the shaping was refused: %s", m.refusal)
	}

	item, err := m.store.LoadBacklogItem("BL01")
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != "BL01" {
		t.Errorf("the id is %q, want BL01 — the id is kind-neutral", item.ID)
	}
	if item.Kind != core.KindFeature {
		t.Errorf("the kind is %q, want %q", item.Kind, core.KindFeature)
	}
	if len(item.Scope) != 1 || item.Scope[0] != "the evidence lane" {
		t.Errorf("the scope is %v, want the one the chef gave", item.Scope)
	}
	if item.Status != core.StatusOpen {
		t.Errorf("the status is %q, want %q — a shaping is not a promotion", item.Status, core.StatusOpen)
	}

	// No recipe was created: this is an update, not a promotion.
	after, err := m.ops.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Fatalf("the plan holds %d recipes, want %d — a shaping creates no recipe", len(after), len(before))
	}
}

// TestTheShapePromptGivesEveryKindItsOwnKey covers the prompt's keys.
//
// The domain owns the kind; the console owns the letter that picks it. And the
// console's own rule is that the letters it hands out **do not collide**: two kinds
// sharing an initial would leave one of them unreachable in the prompt.
//
// The domain's kinds are a `var`, so a kind whose name starts with the same letter
// is a matter of time — and it would vanish silently.
func TestTheShapePromptGivesEveryKindItsOwnKey(t *testing.T) {
	// Two kinds sharing an initial.
	kinds := []string{"feature", "funnel"}
	choices := shapeChoices(core.BacklogMove{Kinds: kinds})
	if len(choices) != len(kinds) {
		t.Fatalf("the prompt offers %d kinds, want %d", len(choices), len(kinds))
	}
	seen := map[string]bool{}
	for i, ch := range choices {
		if ch.value != kinds[i] {
			t.Errorf("the choice is %q, want %q", ch.value, kinds[i])
		}
		if seen[ch.key] {
			t.Errorf("two kinds share the key %q — one of them is unreachable", ch.key)
		}
		seen[ch.key] = true
	}
}

// TestTheShapePromptsKeysAreDistinct covers the prompt the chef sees.
//
// The kinds come from the domain, and the keys the console hands out are distinct.
func TestTheShapePromptsKeysAreDistinct(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	if _, err := Drive(m, []string{"space"}); err != nil {
		t.Fatal(err)
	}
	if m.prompt == nil {
		t.Fatal("the prompt should be open")
	}
	seen := map[string]bool{}
	for _, ch := range m.prompt.choices {
		if ch.key == "" {
			t.Errorf("%q has no key — it is unreachable", ch.value)
		}
		if seen[ch.key] {
			t.Errorf("two kinds share the key %q — one of them is unreachable", ch.key)
		}
		seen[ch.key] = true
	}
}

// TestTheShapePromptOffersTheDomainsKinds covers the prompt's choices.
//
// The kinds are the **domain's**: `core.BacklogMoveOf` carries the kinds a shaping
// may pick, and the prompt offers exactly those. A console-side copy of a domain
// fact is the shape F-114 closed — and a copy that **disagrees** with the domain is
// what this holds.
func TestTheShapePromptOffersTheDomainsKinds(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	item, ok := m.selectedBacklogItem()
	if !ok {
		t.Fatal("there is no item under the cursor")
	}
	move := core.BacklogMoveOf(item)
	if len(move.Kinds) == 0 {
		t.Fatal("the domain names no kind to shape into")
	}

	if _, err := Drive(m, []string{"space"}); err != nil {
		t.Fatal(err)
	}
	if m.prompt == nil {
		t.Fatal("the prompt should be open")
	}
	if len(m.prompt.choices) != len(move.Kinds) {
		t.Fatalf("the prompt offers %d kinds, want the domain's %d (%v)",
			len(m.prompt.choices), len(move.Kinds), move.Kinds)
	}
	for i, ch := range m.prompt.choices {
		if ch.value != move.Kinds[i] {
			t.Errorf("the prompt offers %q, want the domain's %q", ch.value, move.Kinds[i])
		}
	}
}

// TestTheShapePromptAsksForTheKindAndTheScope covers the prompt.
//
// The console asks for what the domain needs and no more: what kind of work this
// is, and what is in it. A small task without a scope is work nobody can tell the
// end of.
func TestTheShapePromptAsksForTheKindAndTheScope(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	out, err := Drive(m, []string{"space"})
	if err != nil {
		t.Fatal(err)
	}
	plainOut := StripANSI(out)
	for _, want := range []string{"f) feature", "b) bug", "esc cancel"} {
		if !strings.Contains(plainOut, want) {
			t.Errorf("the prompt should offer %q:\n%s", want, plainOut)
		}
	}

	// And the scope is asked for in the next step.
	out, err = Drive(m, []string{"f"})
	if err != nil {
		t.Fatal(err)
	}
	plainOut = StripANSI(out)
	if !strings.Contains(plainOut, "scope for BL01") {
		t.Errorf("the prompt should ask for the scope:\n%s", plainOut)
	}
}

// TestAShapingWithoutAScopeIsNotWritten covers the empty answer.
//
// A small task without a scope is not a shaped item, so nothing is written and the
// prompt stays open: the chef types again rather than losing what they wrote.
func TestAShapingWithoutAScopeIsNotWritten(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	if _, err := Drive(m, []string{"space", "f", "enter"}); err != nil {
		t.Fatal(err)
	}
	if m.prompt == nil {
		t.Fatal("the prompt should stay open")
	}
	item, err := m.store.LoadBacklogItem("BL01")
	if err != nil {
		t.Fatal(err)
	}
	if item.Kind != core.KindIdea {
		t.Errorf("the kind is %q, want %q — nothing was written", item.Kind, core.KindIdea)
	}
}

// TestDismissingAnItemRecordsTheReason covers the dismissal.
//
// The domain appends the reason to the item's body, so a dismissal without one
// says nothing to the next reader. The console asks for it, in the chef's own
// words.
func TestDismissingAnItemRecordsTheReason(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	steps := append([]string{"d"}, typed("not now the board already shows it")...)
	steps = append(steps, "enter")
	if _, err := Drive(m, steps); err != nil {
		t.Fatal(err)
	}
	if m.refusal != "" {
		t.Fatalf("the dismissal was refused: %s", m.refusal)
	}

	item, err := m.store.LoadBacklogItem("BL01")
	if err != nil {
		t.Fatal(err)
	}
	if item.Status != core.StatusDismissed {
		t.Fatalf("the status is %q, want %q", item.Status, core.StatusDismissed)
	}
	if !strings.Contains(item.Body, "not now the board already shows it") {
		t.Errorf("the body does not carry the reason:\n%s", item.Body)
	}
}

// TestADismissalWithoutAReasonIsNotWritten covers the empty answer.
//
// A dismissal without a reason says nothing to the next reader, so nothing is
// written and the prompt stays open.
func TestADismissalWithoutAReasonIsNotWritten(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	if _, err := Drive(m, []string{"d", "enter"}); err != nil {
		t.Fatal(err)
	}
	if m.prompt == nil {
		t.Fatal("the prompt should stay open")
	}
	item, err := m.store.LoadBacklogItem("BL01")
	if err != nil {
		t.Fatal(err)
	}
	if item.Status != core.StatusOpen {
		t.Errorf("the status is %q, want %q — nothing was written", item.Status, core.StatusOpen)
	}
}

// TestTheFooterOffersTheDismissOnlyWhileTheDomainSaysSo holds the dismiss key.
//
// A key the domain would refuse is a key the footer must not name: a settled item
// is history, and it cannot be dismissed again.
func TestTheFooterOffersTheDismissOnlyWhileTheDomainSaysSo(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0

	out := plain(m)
	if !strings.Contains(out, "d dismiss") {
		t.Errorf("the footer should offer the dismissal:\n%s", out)
	}

	// Dismiss it: the item is settled, and the key is gone.
	steps := append([]string{"d"}, typed("not now")...)
	steps = append(steps, "enter")
	if _, err := Drive(m, steps); err != nil {
		t.Fatal(err)
	}
	// The cursor leaves the dismissed item. The footer's claim is about that
	// item, so the check sits on it.
	m.selBacklog = 0
	out = plain(m)
	if strings.Contains(out, "d dismiss") {
		t.Errorf("the footer offers a dismissal the domain refuses:\n%s", out)
	}
}

// TestCreatingAnItemWritesABacklogItemAndNothingElse covers the creation.
//
// The chef's ruling (Q6): the console may create a backlog item — a backlog item is
// not a recipe, so ADR-0003 does not forbid it. The item starts as an **idea**,
// the kind a hunch starts as, and the chef shapes it after.
//
// And nothing else is written: no recipe, and nothing outside the plan's backlog.
func TestCreatingAnItemWritesABacklogItemAndNothingElse(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog

	before := planFiles(t, m.store.Root)
	recipes, err := m.ops.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}

	steps := append([]string{"a"}, typed("A weekly digest of what changed")...)
	steps = append(steps, "enter")
	if _, err := Drive(m, steps); err != nil {
		t.Fatal(err)
	}
	if m.refusal != "" {
		t.Fatalf("the creation was refused: %s", m.refusal)
	}

	items, err := m.ops.ListBacklogItems()
	if err != nil {
		t.Fatal(err)
	}
	var created *core.BacklogItem
	for i := range items {
		if items[i].Title == "A weekly digest of what changed" {
			created = &items[i]
		}
	}
	if created == nil {
		t.Fatalf("the item was not written: %+v", items)
	}
	if created.Kind != core.KindIdea {
		t.Errorf("the kind is %q, want %q — a hunch starts as an idea", created.Kind, core.KindIdea)
	}
	if created.Status != core.StatusOpen {
		t.Errorf("the status is %q, want %q", created.Status, core.StatusOpen)
	}

	// No recipe.
	after, err := m.ops.ListRecipes()
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(recipes) {
		t.Fatalf("the plan holds %d recipes, want %d — the console writes no recipe", len(after), len(recipes))
	}

	// And nothing outside the plan's backlog, but the journal — the plan's own
	// event log records the creation.
	changed := changedFiles(before, planFiles(t, m.store.Root))
	for _, path := range changed {
		if path == core.FileActivity {
			continue
		}
		if strings.HasPrefix(path, core.DirBacklog+"/") {
			continue
		}
		t.Errorf("the creation wrote %q — outside the backlog", path)
	}
}

// TestAnItemWithNoTitleIsNotWritten covers the empty answer.
func TestAnItemWithNoTitleIsNotWritten(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog

	if _, err := Drive(m, []string{"a", "enter"}); err != nil {
		t.Fatal(err)
	}
	if m.prompt == nil {
		t.Fatal("the prompt should stay open")
	}
	items, err := m.ops.ListBacklogItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 6 {
		t.Fatalf("the backlog holds %d items, want 6 — nothing was written", len(items))
	}
}

// TestTheBacklogCursorFollowsTheItem covers the selection.
//
// The cursor remembers the item's id. After a shape it is on that id, on its
// new row. After a create it is on the new id. After a dismiss and a promote
// it does not name the id that left: it sits on the nearest remaining row, it
// stays on the Backlog tab, and it does not follow a promotion into the recipe.
// A reload that does not remove the id leaves that id selected.
func TestTheBacklogCursorFollowsTheItem(t *testing.T) {
	m := newModel(t)
	m.tab = TabBacklog
	m.selBacklog = 0
	m.selRecipes = 2

	steps := append([]string{"space", "f"}, typed("the evidence lane")...)
	steps = append(steps, "enter")
	if _, err := Drive(m, steps); err != nil {
		t.Fatal(err)
	}
	item, ok := m.selectedBacklogItem()
	if !ok || item.ID != "BL01" {
		t.Fatalf("after the shape the cursor is on %q, want BL01", idOrNone(item, ok))
	}
	if m.selBacklog == 0 {
		t.Fatal("the cursor stayed on the old index — that row is another idea")
	}
	if m.tab != TabBacklog {
		t.Fatalf("the tab is %v, want the backlog", m.tab)
	}

	steps = append([]string{"a"}, typed("A weekly digest of what changed")...)
	steps = append(steps, "enter")
	if _, err := Drive(m, steps); err != nil {
		t.Fatal(err)
	}
	item, ok = m.selectedBacklogItem()
	if !ok || item.Title != "A weekly digest of what changed" {
		t.Fatalf("after the create the cursor is on %q, want the new item", idOrNone(item, ok))
	}
	created := item.ID
	if m.tab != TabBacklog {
		t.Fatal("a create switched tabs")
	}
	if m.selRecipes != 2 {
		t.Fatalf("the recipe cursor is %d, want 2", m.selRecipes)
	}

	ops := core.NewOperations(m.store, core.AuditContext{Actor: "runner", Session: "s-9", Role: "runner"})
	if _, err := ops.DocumentRecipe("RCP01", "the plan moved on"); err != nil {
		t.Fatal(err)
	}
	m.Update(m.reload()())
	item, ok = m.selectedBacklogItem()
	if !ok || item.ID != created {
		t.Fatalf("the reload left the cursor on %q, want %s", idOrNone(item, ok), created)
	}
	if m.tab != TabBacklog || m.selRecipes != 2 {
		t.Fatalf("the reload moved the chef: tab %v, recipe %d", m.tab, m.selRecipes)
	}

	steps = append([]string{"d"}, typed("not this one")...)
	steps = append(steps, "enter")
	if _, err := Drive(m, steps); err != nil {
		t.Fatal(err)
	}
	item, ok = m.selectedBacklogItem()
	if !ok || item.ID == created {
		t.Fatalf("after the dismissal the cursor is on %q, want the nearest remaining row", idOrNone(item, ok))
	}
	if m.tab != TabBacklog {
		t.Fatal("a dismissal switched tabs")
	}

	// BL05 is the open small task the seed filed. Promoting it must not leave
	// the cursor on it, and must not chase the new recipe.
	if i := backlogRowOf(m.backlogRows(), "BL05"); i >= 0 {
		m.selBacklog = i
	}
	m.Update(tea.KeyPressMsg{Code: ' '})
	if m.refusal != "" {
		t.Fatalf("the promotion was refused: %s", m.refusal)
	}
	m.Update(m.reload()())
	item, ok = m.selectedBacklogItem()
	if !ok || item.ID == "BL05" {
		t.Fatalf("after the promotion the cursor is on %q, want the nearest remaining row", idOrNone(item, ok))
	}
	if item.ID != "BL06" {
		t.Fatalf("after the promotion the cursor is on %q, want BL06", item.ID)
	}
	if m.tab != TabBacklog {
		t.Fatal("a promotion switched tabs")
	}
	if m.selRecipes != 2 {
		t.Fatalf("the cursor followed the recipe: selRecipes = %d, want 2", m.selRecipes)
	}
}

// idOrNone is the id under the cursor, for a failure.
func idOrNone(item core.BacklogItem, ok bool) string {
	if !ok {
		return "(none)"
	}
	return item.ID
}

// planFiles is every file in the plan, by relative path, with its content's digest.
func planFiles(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(raw)
		out[rel] = hex.EncodeToString(sum[:])
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// changedFiles is the paths whose content differs between two snapshots.
func changedFiles(before, after map[string]string) []string {
	var out []string
	for path, sum := range after {
		if before[path] != sum {
			out = append(out, path)
		}
	}
	for path := range before {
		if _, ok := after[path]; !ok {
			out = append(out, path)
		}
	}
	sort.Strings(out)
	return out
}
