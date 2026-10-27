package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestADependencyIsReachableFromTheCLI covers the field with no writer.
//
// `depends_on` was in the model and in the gate table, but nothing could write it — so
// a plan built through `scale` could never express a dependency, and the
// `dependencies.satisfied` gate could never be reached from the agent surface.
func TestADependencyIsReachableFromTheCLI(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "The dependent",
		"--outcome", "o", "--problem", "p",
		"--out-of-scope", "a", "--out-of-scope", "b")
	cli(t, dir, "recipe", "create", "--title", "The dependency",
		"--outcome", "o", "--problem", "p",
		"--out-of-scope", "a", "--out-of-scope", "b")

	res := cli(t, dir, "recipe", "update", "RCP01", "--depends-on", "RCP02")
	if !res.ok(t) {
		t.Fatalf("recipe update failed:\n%s", res.Stderr)
	}

	got := readRecipe(t, dir, "RCP01")
	deps, _ := got["depends_on"].([]any)
	if len(deps) != 1 || deps[0] != "RCP02" {
		t.Fatalf("depends_on = %v, want [RCP02]", got["depends_on"])
	}
}

// TestTheDependencyGateIsReachable covers the gate the flag feeds.
//
// A dependency you cannot write is a gate you cannot trip: the claim would always
// succeed, and the Path tab's "waiting" group would always be empty.
func TestTheDependencyGateIsReachable(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "The dependent",
		"--outcome", "o", "--problem", "p",
		"--out-of-scope", "a", "--out-of-scope", "b")
	cli(t, dir, "recipe", "create", "--title", "The dependency",
		"--outcome", "o", "--problem", "p",
		"--out-of-scope", "a", "--out-of-scope", "b")
	cli(t, dir, "recipe", "update", "RCP01", "--depends-on", "RCP02")
	cli(t, dir, "ac", "add", "RCP01",
		"--description", "the recipe waits for its dependency",
		"--verification", "go test ./internal/cli")
	cli(t, dir, "wp", "add", "RCP01", "--name", "core")
	cli(t, dir, "step", "add", "RCP01",
		"--action", "wait for RCP02", "--done-when", "RCP02 is closed", "--covers", "AC-01")

	// RCP01 waits for RCP02.
	path := cli(t, dir, "path")
	if !path.ok(t) {
		t.Fatalf("path failed:\n%s", path.Stderr)
	}
	if !strings.Contains(path.Stdout, "needs RCP02") {
		t.Fatalf("the path does not show RCP01 waiting for RCP02:\n%s", path.Stdout)
	}

	// RCP01 cannot be claimed while RCP02 is open.
	cli(t, dir, "recipe", "advance", "RCP01", "groom")
	cli(t, dir, "recipe", "advance", "RCP01", "specify")
	cli(t, dir, "--actor", "chef", "recipe", "approve", "RCP01", "--by", "chef")
	res := cli(t, dir, "recipe", "claim", "RCP01", "--agent", "line-cook")
	if res.Code != ExitGate {
		t.Fatalf("exit = %d, want %d — a recipe with an open dependency must not be claimed",
			res.Code, ExitGate)
	}
	if !strings.Contains(res.Stderr, "RCP02") {
		t.Fatalf("the refusal should name the dependency:\n%s", res.Stderr)
	}
}

// TestTheDependencyListIsReplaced covers the list semantics.
//
// The other list flags replace the whole list, and a dependency you cannot take off is
// a dependency you cannot undo. So `--depends-on` names the whole set.
func TestTheDependencyListIsReplaced(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o", "--problem", "p")
	cli(t, dir, "recipe", "create", "--title", "Another", "--outcome", "o", "--problem", "p")

	cli(t, dir, "recipe", "update", "RCP01", "--depends-on", "RCP02")
	cli(t, dir, "recipe", "update", "RCP01", "--depends-on", "RCP02", "--depends-on", "RCP03")
	got := readRecipe(t, dir, "RCP01")
	deps, _ := got["depends_on"].([]any)
	if len(deps) != 2 {
		t.Fatalf("depends_on = %v, want two entries", got["depends_on"])
	}
}

// TestTheDependencyFlagAloneChangesNothing covers the fragment contract.
//
// A nil field means "leave alone". The flag is only set when it was given, so an
// update that says nothing about dependencies keeps them.
func TestTheDependencyFlagAloneChangesNothing(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o", "--problem", "p")
	cli(t, dir, "recipe", "create", "--title", "Another", "--outcome", "o", "--problem", "p")
	cli(t, dir, "recipe", "update", "RCP01", "--depends-on", "RCP02")

	res := cli(t, dir, "recipe", "update", "RCP01", "--title", "Renamed")
	if !res.ok(t) {
		t.Fatalf("recipe update failed:\n%s", res.Stderr)
	}
	got := readRecipe(t, dir, "RCP01")
	deps, _ := got["depends_on"].([]any)
	if len(deps) != 1 || deps[0] != "RCP02" {
		t.Fatalf("depends_on = %v, want [RCP02] — an update without the flag leaves it alone",
			got["depends_on"])
	}
}

// readRecipe reads one recipe's file out of the plan.
func readRecipe(t *testing.T, dir, id string) map[string]any {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(dir, "kitchen", "recipes"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), id+"-") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, "kitchen", "recipes", e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		var got map[string]any
		if err := json.Unmarshal(raw, &got); err != nil {
			t.Fatal(err)
		}
		return got
	}
	t.Fatalf("no recipe file for %s", id)
	return nil
}
