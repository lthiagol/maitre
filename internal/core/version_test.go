package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// recordedDomainShape is the shape of the domain's rules when DomainVersion was
// last bumped.
//
// It changes when a gate is added, renamed or dropped, when a transition changes,
// when a lifecycle is added, when an event is added, when the chef gains a move and
// when the rules on a tasting note's outcome change. Any of those means the rules
// changed, and a plan written under the old rules must not be written under the new
// ones by a tool that does not know the difference.
const recordedDomainShape = "0b55ac7113d0d99c"

// domainShape is a fingerprint of the domain's rules: the gates, the transitions, the
// lifecycles, the events, the chef's moves and the rules on a tasting note's
// outcome.
//
// It is the shape, not the behaviour: what the rules are, not what they do. A rule
// that is not in here is a rule the version stamp cannot see — so a rule that changes
// has to be in here, or the plan cannot say which rules it was written under.
func domainShape() string {
	var b strings.Builder
	for _, g := range GateNames() {
		b.WriteString("gate " + g + "\n")
	}
	for _, ev := range AllEvents {
		def := transitions[ev]
		b.WriteString("event " + string(ev))
		for _, from := range def.from {
			b.WriteString(" from " + string(from))
		}
		b.WriteString(" to " + string(def.to))
		for _, g := range def.gates {
			b.WriteString(" gate " + g)
		}
		b.WriteString("\n")
	}
	for _, l := range AllLifecycles {
		b.WriteString("lifecycle " + string(l) + "\n")
	}
	for _, mv := range AllChefMoves {
		b.WriteString("move " + string(mv) + "\n")
	}
	for _, o := range NoteOutcomes {
		b.WriteString("note " + string(o.Move))
		for _, from := range o.From {
			b.WriteString(" from " + from)
		}
		b.WriteString(" to " + o.To + "\n")
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])[:16]
}

// TestDomainShapeIsVersioned keeps the rules and the version together.
//
// The rules changed, so say so with a version: a plan written under the old rules
// must not be written under the new ones by a tool that does not know the difference.
func TestDomainShapeIsVersioned(t *testing.T) {
	got := domainShape()
	if got != recordedDomainShape {
		t.Fatalf("the domain's rules changed: the shape is now %s, the recorded shape is %s. The rules changed, bump DomainVersion and update this fingerprint.", got, recordedDomainShape)
	}
}

// TestDomainVersionRefusesANewerPlan covers the refusal.
//
// A tool older than the plan would apply weaker rules than the plan was written
// under, so it must not write.
func TestDomainVersionRefusesANewerPlan(t *testing.T) {
	ops, store := newOps(t)
	if _, err := ops.CreateRecipe(readySpec()); err != nil {
		t.Fatal(err)
	}
	newer := DomainVersion + 1
	stampPlanVersion(t, store, newer)

	_, err := ops.CreateSpecial("A hunch", "worth noticing", nil)
	if err == nil {
		t.Fatal("a tool older than the plan must not write to it")
	}
	if !strings.Contains(err.Error(), strconv.Itoa(newer)) {
		t.Fatalf("the refusal should name the plan's version: %v", err)
	}
	if !strings.Contains(err.Error(), strconv.Itoa(DomainVersion)) {
		t.Fatalf("the refusal should name the tool's version: %v", err)
	}
	if !strings.Contains(err.Error(), "upgrade") {
		t.Fatalf("the refusal should say what to do: %v", err)
	}

	// The plan is untouched.
	specials, err := store.ListSpecials()
	if err != nil {
		t.Fatal(err)
	}
	if len(specials) != 0 {
		t.Fatalf("the refused write landed anyway: %+v", specials)
	}
}

// TestDomainVersionStampsAnOlderPlan covers the upgrade.
//
// A tool newer than the plan proceeds, and its next write stamps the newer version.
func TestDomainVersionStampsAnOlderPlan(t *testing.T) {
	ops, store := newOps(t)
	if _, err := ops.CreateRecipe(readySpec()); err != nil {
		t.Fatal(err)
	}
	stampPlanVersion(t, store, DomainVersion-1)

	if _, err := ops.CreateSpecial("A hunch", "worth noticing", nil); err != nil {
		t.Fatalf("a newer tool should write to an older plan: %v", err)
	}
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) == 0 {
		t.Fatal("the journal is empty")
	}
	if got := j.Events[len(j.Events)-1].DomainVersion; got != DomainVersion {
		t.Fatalf("the plan is stamped with version %d, want %d", got, DomainVersion)
	}
}

// TestDomainVersionAllowsReadingANewerPlan covers the other half of the refusal.
//
// Reading is allowed: the chef can still look at a plan written under newer rules.
func TestDomainVersionAllowsReadingANewerPlan(t *testing.T) {
	ops, store := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	stampPlanVersion(t, store, DomainVersion+1)

	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatalf("a plan written under newer rules must still be readable: %v", err)
	}
	if got.ID != r.ID {
		t.Fatalf("read %q, want %q", got.ID, r.ID)
	}
	if _, err := store.LoadJournal(); err != nil {
		t.Fatalf("the journal must still be readable: %v", err)
	}
	if _, err := ops.ListRecipes(); err != nil {
		t.Fatalf("the recipe list must still be readable: %v", err)
	}
}

// stampPlanVersion writes a domain version into the plan's journal, as if a tool of
// that version had written it.
func stampPlanVersion(t *testing.T, store *Store, version int) {
	t.Helper()
	j, err := store.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) == 0 {
		t.Fatal("the journal is empty")
	}
	j.Events[len(j.Events)-1].DomainVersion = version
	raw, err := json.Marshal(j)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(store.Root, FileActivity), raw, 0o644); err != nil {
		t.Fatal(err)
	}
}
