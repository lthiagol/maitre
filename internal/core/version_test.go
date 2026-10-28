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

// recordedDomainShape is the shape of the domain's rules: the fingerprint of the
// gates, the transitions, the lifecycles, the events, the chef's moves, the rules on
// a tasting note's outcome, and the facts of the **plan's format**.
//
// It changes when a gate is added, renamed or dropped, when a transition changes,
// when a lifecycle is added, when an event is added, when the chef gains a move,
// when the rules on a tasting note's outcome change, and when a fact of the
// **plan's format** changes — the id prefixes, the directories, the root file
// names, the field names on disk, and the journal envelope's own fields.
//
// The shape and the **version** are not the same thing, and the difference matters:
//
//   - The **version** fences the **writer**. A tool whose rules are older than the
//     plan's must not write to it, because it would write something weaker than the
//     plan was written under. So a change in *what the tool writes* moves
//     `DomainVersion`.
//   - The **shape** is what the version stamp can see. Learning to see a fact it
//     could not see before changes the shape without changing the writer — the plan
//     on disk is the same — and that is not a version bump.
//
// So a format change moves both, and only the shape moves when the fingerprint
// learns to see a fact it could not see before.
const recordedDomainShape = "b84b230161f4ab03"

// domainShape is a fingerprint of the domain's rules: the gates, the transitions, the
// lifecycles, the events, the chef's moves, the rules on a tasting note's outcome,
// and the facts of the plan's format.
//
// It is the shape, not the behaviour: what the rules are, not what they do. A rule
// that is not in here is a rule the version stamp cannot see — so a rule that changes
// has to be in here, or the plan cannot say which rules it was written under.
//
// The format is in it because the format is a rule: a tool that writes the old
// ids, or a record under a field name the reader does not know, writes a plan the
// reader cannot read — and the reader cannot tell that the plan is broken. The root
// file names and the journal envelope's own fields are in it for the same reason:
// the journal's name is the silent one, because a missing journal reads as an empty
// one.
func domainShape() string {
	return shapeOfFacts(planFormatFacts())
}

// shapeOfFacts is the fingerprint of the domain's rules, given the facts of the
// plan's format.
//
// The format's facts are an argument so the wire can be tested: a fact that is in
// the list must change the fingerprint, and a fact that is not in it cannot. The
// hole F-71 was filed to close is exactly a fact that was not in the list — the
// root file names and the journal envelope's own fields — so a change to one of
// them left the fingerprint untouched and the version had to be moved by hand.
func shapeOfFacts(formatFacts []string) string {
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
	for _, f := range formatFacts {
		b.WriteString("format " + f + "\n")
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
		t.Fatalf("the domain's rules changed: the shape is now %s, the recorded shape is %s. Update this fingerprint. If the change is in what the tool writes, bump DomainVersion as well — the version fences the writer, and the shape alone fences nothing.", got, recordedDomainShape)
	}
}

// TestTheFormatFactsCoverTheRootFilesAndTheJournalEnvelope holds the facts the
// fingerprint was missing.
//
// The journal's file name is the silent one: a missing journal reads as an empty
// one. And the envelope's own fields are the format too — a journal that carries a
// different `schema` is a journal a reader cannot read. A fact that is not in the
// list is a fact the version stamp cannot see, so a change to it would have to move
// the version by hand.
func TestTheFormatFactsCoverTheRootFilesAndTheJournalEnvelope(t *testing.T) {
	facts := strings.Join(planFormatFacts(), "\n")
	for _, want := range []string{
		"file " + FileProject,
		"file " + FileActivity,
		"file " + FileDecisions,
		"field Journal.schema",
		"field Journal.seq",
		"field Journal.events",
	} {
		if !strings.Contains(facts, want) {
			t.Errorf("the format's facts should name %q", want)
		}
	}
}

// TestEveryFormatFactMovesTheFingerprint keeps the facts honest.
//
// A fact in the list that the fingerprint cannot see is a dead fact: it reads like
// the format is covered, and a change to it would leave the version where it is.
// Every fact must move the fingerprint, or it is not a fact.
func TestEveryFormatFactMovesTheFingerprint(t *testing.T) {
	facts := planFormatFacts()
	base := shapeOfFacts(facts)
	for _, fact := range facts {
		without := make([]string, 0, len(facts))
		for _, f := range facts {
			if f != fact {
				without = append(without, f)
			}
		}
		if shapeOfFacts(without) == base {
			t.Errorf("the fingerprint does not see the fact %q", fact)
		}
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

	_, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
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
	items, err := store.ListBacklogItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("the refused write landed anyway: %+v", items)
	}
}

// TestDomainVersionRefusesAnOlderPlan covers the other side.
//
// There is no upgrade path: a newer tool does not write a plan written under older
// rules either, because it would apply rules the plan was not written under. The
// refusal says what the chef can do, and it promises no migration.
func TestDomainVersionRefusesAnOlderPlan(t *testing.T) {
	ops, store := newOps(t)
	if _, err := ops.CreateRecipe(readySpec()); err != nil {
		t.Fatal(err)
	}
	older := DomainVersion - 1
	stampPlanVersion(t, store, older)

	_, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, KindIdea, nil)
	if err == nil {
		t.Fatal("a newer tool must not write to a plan written under older rules")
	}
	if !strings.Contains(err.Error(), strconv.Itoa(older)) {
		t.Fatalf("the refusal should name the plan's version: %v", err)
	}
	if !strings.Contains(err.Error(), strconv.Itoa(DomainVersion)) {
		t.Fatalf("the refusal should name the tool's version: %v", err)
	}
	if strings.Contains(err.Error(), "migrat") {
		t.Fatalf("the refusal should not promise a migration: %v", err)
	}

	// The plan is untouched.
	items, err := store.ListBacklogItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("the refused write landed anyway: %+v", items)
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

// TestTheRepoPlanIsWrittenUnderThisToolsVersion covers the repo's own plan.
//
// The version guard refuses a write when the plan's version is not the tool's
// version. The repo's own plan is stamped with the version this tool speaks, so it
// stays writable — and the test says so, so a version bump that leaves the repo's
// own plan behind fails here.
func TestTheRepoPlanIsWrittenUnderThisToolsVersion(t *testing.T) {
	root := filepath.Join("..", "..", PlanDirName)
	if _, err := os.Stat(root); err != nil {
		t.Skipf("the repo has no plan at %s", root)
	}
	s := Open(root)
	j, err := s.LoadJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Events) == 0 {
		t.Fatal("the repo's own journal is empty")
	}
	plan := PlanDomainVersion(j)
	if plan != DomainVersion {
		t.Fatalf("the repo's own plan is stamped %d and this tool speaks %d: the plan is written under another version", plan, DomainVersion)
	}
	if err := s.checkDomainVersion(); err != nil {
		t.Fatalf("the repo's own plan must stay writable: %v", err)
	}
}
