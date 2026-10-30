package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// identityStore seeds a plan under a directory with a fixed name, so the project's
// name in the render is deterministic.
func identityStore(t *testing.T) *core.Store {
	t.Helper()
	root := filepath.Join(t.TempDir(), "demo", core.PlanDirName)
	store := core.Open(root)
	if err := core.Seed(store); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return store
}

// TestWindowTitleCarriesTheProject proves the terminal's tab says which project it
// is.
//
// The window title is not in the rendered content, so `--render` cannot see it: a
// test is the only instrument here.
func TestWindowTitleCarriesTheProject(t *testing.T) {
	store := identityStore(t)
	m := loadInto(t, New(store))

	if got, want := m.View().WindowTitle, "maitre · demo"; got != want {
		t.Errorf("the window title is %q, want %q", got, want)
	}
}

// TestTheWindowTitleAndTheHeaderAgree proves the name has one source.
//
// Two computations would drift, and the drift would show as a terminal whose tab
// says one project and whose header says another.
func TestTheWindowTitleAndTheHeaderAgree(t *testing.T) {
	store := identityStore(t)
	m := loadInto(t, New(store))

	name := m.projectName()
	if !strings.Contains(m.View().WindowTitle, name) {
		t.Errorf("the window title %q does not carry the name %q", m.View().WindowTitle, name)
	}
	if !strings.Contains(plain(m), "· "+name) {
		t.Errorf("the header does not carry the name %q\n%s", name, plain(m))
	}
}

// TestTheNameFallsBackToThePlanRoot proves a plan with no name still has one.
//
// The plan's own name is the answer; without one, the console falls back to the
// plan's root, so the tab is never a bare "maitre".
func TestTheNameFallsBackToThePlanRoot(t *testing.T) {
	store := identityStore(t)
	// A plan whose project file carries no name.
	if err := os.WriteFile(filepath.Join(store.Root, core.FileProject),
		[]byte(`{"name":"","created":"2026-01-01T00:00:00Z"}`), 0o644); err != nil {
		t.Fatalf("write project: %v", err)
	}
	m := loadInto(t, New(store))

	if got, want := m.projectName(), store.Root; got != want {
		t.Errorf("the name is %q, want the plan's root %q", got, want)
	}
	if got, want := m.View().WindowTitle, "maitre · "+store.Root; got != want {
		t.Errorf("the window title is %q, want %q", got, want)
	}
	if !strings.Contains(plain(m), "· "+store.Root) {
		t.Errorf("the header does not carry the fallback\n%s", plain(m))
	}
}
