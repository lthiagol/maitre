package ui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// sousChefModel is a console over a plan with a journal.
func sousChefModel(t *testing.T) *Model {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	ops := core.NewOperations(store, core.AuditContext{Actor: "chef", Session: "s-1", Role: "chef"})
	if _, err := ops.CreateBacklogItem("A hunch", "worth noticing", nil, core.KindIdea, nil); err != nil {
		t.Fatal(err)
	}

	m := loadModel(t, store)
	m.height = 40
	return m
}

// TestTheSousChefTabIsThere covers the tab.
//
// It is a placeholder, and a placeholder that looks broken is worse than no tab at
// all — so it has to render, and it has to say what it is.
func TestTheSousChefTabIsThere(t *testing.T) {
	m := sousChefModel(t)
	switchToTab(t, m, TabSousChef)
	out := plain(m)

	if !strings.Contains(out, "Sous-Chef") {
		t.Errorf("the tab bar should name the Sous-Chef:\n%s", out)
	}
	// The bar's order: the Sous-Chef sits with the Log, the other view of the
	// automation's running.
	if !strings.Contains(out, "Log  ▸ Sous-Chef") && !strings.Contains(out, "Log    Sous-Chef") {
		t.Errorf("the Sous-Chef should sit next to the Log:\n%s", out)
	}
	if !strings.Contains(out, "The orchestrator agent") {
		t.Errorf("the tab should say what the Sous-Chef is:\n%s", out)
	}
	if !strings.Contains(out, "Not built yet.") {
		t.Errorf("the tab should say it is not built yet:\n%s", out)
	}
}

// TestTheSousChefShowsTheVitalSigns covers the versions.
//
// The chef can already see whether the automation is running on the same rules as the
// plan — that is why this tab exists before the autopilot does.
func TestTheSousChefShowsTheVitalSigns(t *testing.T) {
	m := sousChefModel(t)
	switchToTab(t, m, TabSousChef)
	out := plain(m)

	if !strings.Contains(out, "plan written under") {
		t.Errorf("the tab should show the plan's version:\n%s", out)
	}
	if !strings.Contains(out, "this build speaks") {
		t.Errorf("the tab should show this build's version:\n%s", out)
	}
	if !strings.Contains(out, "they agree") {
		t.Errorf("the tab should say whether they agree:\n%s", out)
	}
}

// TestTheSousChefSaysWhenTheVersionsDisagree covers the mismatch.
//
// A plan written by a different build of the tool is a fact the chef needs, and the
// version guard turns it into a refusal to write.
func TestTheSousChefSaysWhenTheVersionsDisagree(t *testing.T) {
	m := sousChefModel(t)
	store := m.store
	stampPlan(t, store, core.DomainVersion+1)

	m = loadModel(t, store)
	switchToTab(t, m, TabSousChef)
	out := plain(m)

	if !strings.Contains(out, "they do not agree") {
		t.Errorf("the tab should say the versions disagree:\n%s", out)
	}
	if !strings.Contains(out, "different build of the tool") {
		t.Errorf("the tab should say what a mismatch means:\n%s", out)
	}
}

// TestTheSousChefShowsTheLastEvent covers the history.
//
// The newest event, with its kind, the recipe, the actor and the time.
func TestTheSousChefShowsTheLastEvent(t *testing.T) {
	m := sousChefModel(t)
	switchToTab(t, m, TabSousChef)
	out := plain(m)

	if !strings.Contains(out, "LAST EVENT") {
		t.Errorf("the tab should show the last event:\n%s", out)
	}
	if !strings.Contains(out, "created") {
		t.Errorf("the last event's kind should be shown:\n%s", out)
	}
	if !strings.Contains(out, "chef") {
		t.Errorf("the last event's actor should be shown:\n%s", out)
	}
}

// TestTheSousChefSaysWhenThereAreNoEvents covers the empty plan.
//
// A plan with no journal has no last event. It does not have a blank one.
func TestTheSousChefSaysWhenThereAreNoEvents(t *testing.T) {
	root := filepath.Join(t.TempDir(), core.PlanDirName)
	store := core.Open(root)
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	m.tab = TabSousChef
	out := plain(m)

	if !strings.Contains(out, "the plan has no events") {
		t.Errorf("an empty plan should say it has no events:\n%s", out)
	}
}

// stampPlan writes a domain version into the plan's journal, as if a tool of that
// version had written it.
func stampPlan(t *testing.T, store *core.Store, version int) {
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
	if err := os.WriteFile(filepath.Join(store.Root, core.FileActivity), raw, 0o644); err != nil {
		t.Fatal(err)
	}
}
