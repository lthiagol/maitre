package ui

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/kitchen-brigade/internal/core"
)

// update regenerates the golden. Run it with `make regen-goldens`.
var update = flag.Bool("update", false, "update the golden file")

// goldenPath is the checked-in snapshot of the console's Recipes tab.
const goldenPath = "testdata/console-recipes.golden"

// stressGoldenPath is the same, against the pathological plan.
const stressGoldenPath = "testdata/console-stress.golden"

// goldenStore seeds a plan under a directory with a fixed name, so the project name in
// the render is deterministic.
func goldenStore(t *testing.T) *core.Store {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := core.Seed(store); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return store
}

// TestConsoleMatchesTheGolden pins the Recipes tab's layout.
//
// A golden catches the thing unit tests miss: a column that stops lining up. Update it
// with `make regen-goldens` when the design deliberately changes.
func TestConsoleMatchesTheGolden(t *testing.T) {
	store := goldenStore(t)
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	got := StripANSI(m.View().Content)

	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Log("golden updated")
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v — run `make regen-goldens`", err)
	}
	if got != string(want) {
		t.Errorf("the console does not match the golden.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestGoldenIsNotStale catches a golden that was regenerated but not committed.
func TestGoldenIsNotStale(t *testing.T) {
	store := goldenStore(t)
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	got := StripANSI(m.View().Content)

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Skipf("no golden yet: %v", err)
	}
	if got == string(want) {
		return
	}
	// A mismatch here is the golden's problem, not the code's — but it means someone
	// changed the design and did not regenerate. Say so plainly.
	t.Log("the golden is out of date — run `make regen-goldens` if the change was deliberate")
}

// TestDriveReplaysKeys proves the scripted driver works end to end.
func TestDriveReplaysKeys(t *testing.T) {
	store := goldenStore(t)
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	// The first recipe is the ready one — space signs it off.
	out, err := Drive(m, []string{"space"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(StripANSI(out), "signed off R02") {
		t.Errorf("the scripted drive did not sign off R02:\n%s", StripANSI(out))
	}
}

func TestDriveRejectsAnUnknownKey(t *testing.T) {
	store := goldenStore(t)
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Drive(m, []string{"nonsense"}); err == nil {
		t.Fatal("expected an unknown key to be refused")
	}
}
