package ui

import (
	"fmt"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
	"github.com/lthiagol/maitre/internal/settings"
)

// sousModel is a console on the Sous-Chef tab.
func sousModel(t *testing.T) *Model {
	t.Helper()
	m := newModel(t)
	m.tab = TabSousChef
	return m
}

// TestTheOrderShowsTheFourControls covers the order.
//
// The ticket's order has four controls — the harness, the crew, the sources and
// the budget — and each one says what it is: the harness says the command it will
// run, and the crew says who reviews.
func TestTheOrderShowsTheFourControls(t *testing.T) {
	m := sousModel(t)
	out := plain(m)

	for _, want := range []string{
		"THE ORDER",
		"harness",
		"command: opencode",
		"crew",
		"orchestrator, runner, reviewer",
		"sources",
		"ready (",
		"budget",
		"rounds · 12 interactions",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the order is missing %q:\n%s", want, out)
		}
	}
}

// TestTheCrewRowSaysWhoReviews covers the crew's row.
//
// The count is a composition, not a number: a lone agent cannot review itself —
// the close gate refuses a self-reviewed recipe — so the one-agent crew is the one
// where the **chef** holds the reviewer's chair.
func TestTheOrderSaysWhoReviews(t *testing.T) {
	m := sousModel(t)

	// The crew of three: one per role.
	if out := plain(m); !strings.Contains(out, "3 · orchestrator, runner, reviewer") {
		t.Errorf("the crew of three should name its three roles:\n%s", out)
	}

	// The crew of one: the chef reviews.
	m.picks.Count = 1
	out := plain(m)
	if !strings.Contains(out, "1 · one agent, and the chef reviews") {
		t.Errorf("the crew of one should say the chef reviews:\n%s", out)
	}
	if strings.Contains(out, "▸ 3 ·") {
		t.Errorf("the crew of three should not be the picked one:\n%s", out)
	}
}

// TestTheOrderIsReadFromTheRecord covers the read.
//
// The order is the record's, not the console's: a console that starts over a plan
// whose record says `pi` runs pi.
func TestTheOrderIsReadFromTheRecord(t *testing.T) {
	store := settingsStore(t)
	rec := settings.Record{
		View: settings.View{Theme: "Kitchen", Order: "triage", ShowServed: true},
		Crew: settings.Crew{Harness: "pi", Count: 2, Sources: "backlog", Rounds: 5, Interactions: 24},
	}
	if err := settings.Write(settings.Path(store.Root), rec); err != nil {
		t.Fatalf("write the settings: %v", err)
	}

	m := loadModel(t, store)
	m.tab = TabSousChef
	out := plain(m)

	if !strings.Contains(out, "command: pi --model sonnet") {
		t.Errorf("the order should carry the harness the record names:\n%s", out)
	}
	if !strings.Contains(out, "2 · runner, and the orchestrator reviews") {
		t.Errorf("the order should carry the crew the record names:\n%s", out)
	}
	if !strings.Contains(out, "backlog (") {
		t.Errorf("the order should carry the sources the record names:\n%s", out)
	}
	if !strings.Contains(out, "5 rounds · 24 interactions") {
		t.Errorf("the order should carry the budget the record names:\n%s", out)
	}
}

// TestChangingAPickWritesTheRecord covers the write.
//
// The console reads the order from the record and writes it when the chef
// changes a pick: an order the console forgot is an order the next console does
// not have.
func TestTheOrderIsWrittenWhenAPickChanges(t *testing.T) {
	store := settingsStore(t)
	m := loadModel(t, store)
	m.tab = TabSousChef

	// The cursor is on the harness: the next key picks the next one.
	m.Update(tea.KeyPressMsg{Code: tea.KeyRight})

	rec, err := settings.Read(settings.Path(store.Root))
	if err != nil {
		t.Fatalf("read the settings: %v", err)
	}
	if rec.Crew.Harness != "pi" {
		t.Errorf("the record's harness is %q, want pi", rec.Crew.Harness)
	}

	// And the crew, the sources and the budget are written the same way.
	m.moveSous(1)
	m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	m.moveSous(1)
	m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	m.moveSous(1)
	m.Update(tea.KeyPressMsg{Code: tea.KeyRight})

	rec, err = settings.Read(settings.Path(store.Root))
	if err != nil {
		t.Fatalf("read the settings: %v", err)
	}
	if rec.Crew.Count != 2 {
		t.Errorf("the record's crew size is %d, want 2", rec.Crew.Count)
	}
	if rec.Crew.Sources != "recipes" {
		t.Errorf("the record's sources are %q, want recipes", rec.Crew.Sources)
	}
	if rec.Crew.Rounds != 5 || rec.Crew.Interactions != 24 {
		t.Errorf("the record's budget is %d rounds · %d interactions, want 5 · 24",
			rec.Crew.Rounds, rec.Crew.Interactions)
	}
}

// TestTheOrderIsNotARecipe covers the tab's own claim.
//
// The console **never authors a recipe** (ADR-0003): changing the order writes the
// tool's settings file and nothing else — the plan's records and the journal are
// untouched.
func TestTheOrderIsNotARecipe(t *testing.T) {
	store := settingsStore(t)
	m := loadModel(t, store)
	m.tab = TabSousChef
	before, err := store.LoadJournal()
	if err != nil {
		t.Fatalf("load the journal: %v", err)
	}

	m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	after, err := store.LoadJournal()
	if err != nil {
		t.Fatalf("load the journal: %v", err)
	}
	if len(after.Events) != len(before.Events) {
		t.Errorf("the journal holds %d events, want %d — the order is the tool's file, not the plan's",
			len(after.Events), len(before.Events))
	}
}

// TestTheOrderCountsTheSources covers the counts.
//
// Each source says how much it would queue, so the chef sees an empty source
// before the crew goes there.
func TestTheOrderCountsTheSources(t *testing.T) {
	m := sousModel(t)
	counts := m.sousCounts()
	for _, source := range sourceNames {
		if counts[source] < 0 {
			t.Fatalf("the source %q has no count", source)
		}
	}
	// The seeded plan has open work, so the ready source is not empty.
	if counts["ready"] == 0 {
		t.Errorf("the seeded plan should have ready work: %v", counts)
	}
	out := plain(m)
	if !strings.Contains(out, "recipes (") {
		t.Errorf("the order should count the recipes source:\n%s", out)
	}
}

// TestTheQueueShowsTheItemsAndTheNext covers the queue.
//
// The queue is the plan's own work, chosen by the order's source, and its first
// item is the **next** one — the crew takes the queue in order.
func TestTheQueueShowsTheItemsAndTheNext(t *testing.T) {
	m := sousModel(t)
	out := plain(m)

	if !strings.Contains(out, "THE QUEUE") {
		t.Fatalf("the tab has no queue:\n%s", out)
	}
	queue := m.sousQueue(m.picks.Sources)
	if len(queue) == 0 {
		t.Fatal("the seeded plan should have queued work")
	}
	if !strings.Contains(out, queue[0].ID) {
		t.Errorf("the queue is missing its first item %q:\n%s", queue[0].ID, out)
	}
	if !strings.Contains(out, "next") {
		t.Errorf("the queue should say which item is next:\n%s", out)
	}
}

// TestTheQueueFollowsTheSourcesPick covers the source.
//
// The queue is the order's source: the chef picks where the work comes from, and
// the queue changes with it.
func TestTheQueueFollowsTheSourcesPick(t *testing.T) {
	m := sousModel(t)

	m.picks.Sources = "recipes"
	recipes := m.sousQueue("recipes")
	if len(recipes) == 0 {
		t.Fatal("the seeded plan should have open recipes")
	}
	out := plain(m)
	if !strings.Contains(out, recipes[0].ID) {
		t.Errorf("the queue should hold the recipes source:\n%s", out)
	}

	m.picks.Sources = "backlog"
	backlog := m.sousQueue("backlog")
	if len(backlog) == 0 {
		t.Fatal("the seeded plan should have open backlog items")
	}
	out = plain(m)
	if !strings.Contains(out, backlog[0].ID) {
		t.Errorf("the queue should hold the backlog source:\n%s", out)
	}
}

// TestTheStampSaysWhenTheAutomationIsNotRunning covers the stamp without a crew.
//
// A tab that is waiting must not look like a tab that is failing: the stamp says
// the automation is not running, and says what to do about it. The console cannot
// tell "never started" from "gone", and it does not guess.
func TestTheStampSaysWhenTheAutomationIsNotRunning(t *testing.T) {
	m := sousModel(t)
	out := plain(m)

	if !strings.Contains(out, "THE STAMP") {
		t.Fatalf("the tab has no stamp:\n%s", out)
	}
	if !strings.Contains(out, "the automation is not running") {
		t.Errorf("the stamp should say the automation is not running:\n%s", out)
	}
	if !strings.Contains(out, "press space to start it") {
		t.Errorf("the stamp should say what to do about it:\n%s", out)
	}
}

// TestTheStampShowsTheVitalSigns covers the stamp with a crew.
//
// The vital signs are the plan's version, this build's version, and whether they
// agree. The tab that starts the automation has to be able to report that it is
// not speaking the plan's rules.
func TestTheStampShowsTheVitalSigns(t *testing.T) {
	m := sousModel(t)
	m.crew = []crewRow{{State: "working", Job: "orchestrator", On: "routing RCP04"}}

	out := plain(m)
	if !strings.Contains(out, fmt.Sprintf("plan %d", core.DomainVersion)) {
		t.Errorf("the stamp should carry the plan's version:\n%s", out)
	}
	if !strings.Contains(out, fmt.Sprintf("build %d", core.DomainVersion)) {
		t.Errorf("the stamp should carry this build's version:\n%s", out)
	}
	if !strings.Contains(out, "they agree") {
		t.Errorf("the stamp should say the versions agree:\n%s", out)
	}
	if !strings.Contains(out, "orchestrator") {
		t.Errorf("the stamp should carry the crew's rows:\n%s", out)
	}
}

// TestTheStampSaysWhenTheVersionsDisagree covers the disagreement.
//
// A build that does not speak the plan's rules says so — the stamp is where the
// chef sees it.
func TestTheStampSaysWhenTheVersionsDisagree(t *testing.T) {
	m := sousModel(t)
	m.crew = []crewRow{{State: "working", Job: "orchestrator", On: "routing RCP04"}}
	// The plan was last written under another version.
	m.journal.Events[len(m.journal.Events)-1].DomainVersion = core.DomainVersion - 1

	out := plain(m)
	if !strings.Contains(out, "they disagree") {
		t.Errorf("the stamp should say the versions disagree:\n%s", out)
	}
}

// TestTheStampCountsTheCrewsProgress covers the progress.
//
// **Policy above, progress below.** The order says the cap — "3 rounds · 12
// interactions" — and the stamp says what has happened against it. The two are
// different facts and neither is repeated.
func TestTheStampCountsTheCrewsProgress(t *testing.T) {
	m := sousModel(t)
	m.crew = []crewRow{{State: "working", Job: "orchestrator", On: "routing RCP04"}}
	m.journal = core.Journal{Schema: core.JournalSchema, Events: []core.JournalEntry{
		{Seq: 1, Kind: core.KindEvidence, Actor: "runner", Role: "runner", Recipe: "RCP02"},
		{Seq: 2, Kind: core.KindEvidence, Actor: "runner", Role: "runner", Recipe: "RCP02"},
		{Seq: 3, Kind: core.KindEvidence, Actor: "reviewer", Role: "reviewer", Recipe: "RCP03"},
		{Seq: 4, Kind: core.KindLifecycle, Actor: "chef", Role: "chef", Recipe: "RCP03"},
	}}

	rounds, interactions := crewProgress(m.journal)
	if rounds != 2 {
		t.Errorf("the rounds are %d, want 2 — two recipes the crew worked", rounds)
	}
	if interactions != 3 {
		t.Errorf("the interactions are %d, want 3 — the chef's own write is not one", interactions)
	}

	out := plain(m)
	if !strings.Contains(out, "round 2 of 3") {
		t.Errorf("the stamp should say the round against the cap:\n%s", out)
	}
	if !strings.Contains(out, "3 of 12 interactions") {
		t.Errorf("the stamp should say the interactions against the cap:\n%s", out)
	}
}
