package cli

import (
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// journalCount is how many events the plan's journal holds.
func journalCount(t *testing.T, dir string) int {
	t.Helper()
	res := cli(t, dir, "log", "--limit", "1000")
	if !res.ok(t) {
		t.Fatalf("log failed:\n%s", res.Stderr)
	}
	events, ok := res.doc(t)["events"].([]any)
	if !ok {
		t.Fatal("the journal document carries no events")
	}
	return len(events)
}

// TestBacklogUpdateGoesThroughAnOperation covers the write that used to bypass the
// domain layer.
//
// The CLI wrote the artifact file itself: no lock, no operation, and — worst — no
// journal entry. The journal is the only event log, so the edit was invisible to
// `maitre-cli log --since`, to the console's Log tab, and to any agent watching for
// changes.
func TestBacklogUpdateGoesThroughAnOperation(t *testing.T) {
	dir := newPlan(t)

	res := cli(t, dir, "backlog", "create", "--title", "An item", "--body", "first")
	if !res.ok(t) {
		t.Fatalf("backlog create failed:\n%s", res.Stderr)
	}
	before := journalCount(t, dir)

	res = cli(t, dir, "backlog", "update", "BL01", "--body", "second")
	if !res.ok(t) {
		t.Fatalf("backlog update failed:\n%s", res.Stderr)
	}

	// The plan file changed.
	items := cli(t, dir, "backlog", "list")
	if !items.ok(t) {
		t.Fatalf("backlog list failed:\n%s", items.Stderr)
	}
	list := items.doc(t)["items"].([]any)
	if len(list) != 1 {
		t.Fatalf("got %d items, want 1", len(list))
	}
	item := list[0].(map[string]any)
	if item["body"] != "second" {
		t.Fatalf("body = %v, want second", item["body"])
	}

	// And the journal knows about it.
	after := journalCount(t, dir)
	if after != before+1 {
		t.Fatalf("the journal has %d events, want %d — the edit must be journalled", after, before+1)
	}
}

// TestBacklogUpdateRefusesAPromotedItem covers the one-way trail.
//
// A promoted item became a recipe; editing it would write to an artifact that is
// already history.
func TestBacklogUpdateRefusesAPromotedItem(t *testing.T) {
	dir := newPlan(t)

	res := cli(t, dir, "backlog", "create", "--title", "An item", "--body", "first")
	if !res.ok(t) {
		t.Fatalf("backlog create failed:\n%s", res.Stderr)
	}
	res = cli(t, dir, "backlog", "promote", "BL01")
	if !res.ok(t) {
		t.Fatalf("backlog promote failed:\n%s", res.Stderr)
	}

	res = cli(t, dir, "backlog", "update", "BL01", "--body", "second")
	if res.Code == 0 {
		t.Fatal("a promoted item must not be edited")
	}
}

// TestTheBacklogKindIsTheShapingStep covers the kind.
//
// An idea becomes a small task by being made a `feature` or a `bug` and gaining a
// scope. It is an update, not a promotion, and the id does not change.
func TestTheBacklogKindIsTheShapingStep(t *testing.T) {
	dir := newPlan(t)

	res := cli(t, dir, "backlog", "create", "--title", "A hunch", "--body", "worth noticing")
	if !res.ok(t) {
		t.Fatalf("backlog create failed:\n%s", res.Stderr)
	}
	item := res.doc(t)["item"].(map[string]any)
	if item["kind"] != "idea" {
		t.Fatalf("the kind = %v, want idea — an item starts as a hunch", item["kind"])
	}
	if item["id"] != "BL01" {
		t.Fatalf("the id = %v, want BL01", item["id"])
	}

	// The shaping.
	res = cli(t, dir, "backlog", "update", "BL01", "--kind", "feature", "--scope", "the queue lane")
	if !res.ok(t) {
		t.Fatalf("backlog update failed:\n%s", res.Stderr)
	}
	item = res.doc(t)["item"].(map[string]any)
	if item["kind"] != "feature" {
		t.Fatalf("the kind = %v, want feature", item["kind"])
	}
	if item["id"] != "BL01" {
		t.Fatalf("the id = %v, want BL01 — the id is kind-neutral", item["id"])
	}
	scope := item["scope"].([]any)
	if len(scope) != 1 || scope[0] != "the queue lane" {
		t.Fatalf("the scope = %v, want the queue lane", scope)
	}

	// A kind nobody knows is refused, not silently accepted.
	res = cli(t, dir, "backlog", "update", "BL01", "--kind", "hunch")
	if res.Code == 0 {
		t.Fatal("an unknown kind must be refused")
	}
}

// TestAnUnknownKindNamesTheKinds covers the closed set.
//
// A kind nobody knows is a typo, and a typo that silently disabled a feature is
// v1's lesson. The refusal names every kind there is, so the caller can see the
// set without reading the docs.
func TestAnUnknownKindNamesTheKinds(t *testing.T) {
	dir := newPlan(t)

	res := cli(t, dir, "backlog", "create", "--title", "A hunch", "--kind", "hunch")
	if res.Code == 0 {
		t.Fatal("an unknown kind must be refused")
	}
	for _, kind := range core.BacklogKinds {
		if !strings.Contains(res.Stderr, kind) {
			t.Errorf("the refusal should name the kind %q: %s", kind, res.Stderr)
		}
	}
}

// TestTheBacklogDismissesEveryKind covers the dismissal.
//
// Only a special could be dismissed before, and no one could say why. An idea, a
// feature and a bug are all things the chef can decide against.
func TestTheBacklogDismissesEveryKind(t *testing.T) {
	dir := newPlan(t)

	for _, kind := range []string{"idea", "feature", "bug"} {
		res := cli(t, dir, "backlog", "create", "--title", "An item", "--kind", kind)
		if !res.ok(t) {
			t.Fatalf("backlog create --kind %s failed:\n%s", kind, res.Stderr)
		}
	}
	for i, id := range []string{"BL01", "BL02", "BL03"} {
		res := cli(t, dir, "backlog", "dismiss", id, "--reason", "not worth it")
		if !res.ok(t) {
			t.Fatalf("backlog dismiss %s failed:\n%s", id, res.Stderr)
		}
		item := res.doc(t)["item"].(map[string]any)
		if item["status"] != "dismissed" {
			t.Fatalf("the status of %s = %v, want dismissed", id, item["status"])
		}
		if i == 0 {
			continue
		}
	}
}
