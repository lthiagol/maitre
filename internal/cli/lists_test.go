package cli

import (
	"encoding/json"
	"fmt"
	"testing"
)

// listCommands are the commands whose document carries a list.
func listCommands() [][]string {
	return [][]string{
		{"plan"},
		{"path"},
		{"log"},
		{"log", "--since", "0"},
		{"recipe", "list"},
		{"special", "list"},
		{"prepnote", "list"},
		{"gates"},
	}
}

// TestNoListIsEverNull covers the JSON contract.
//
// The plan is read by agents, and the standard readers break on `null`: `jq '.ready[]'`
// on an empty plan is an error, not an empty result. An empty list is `[]`.
func TestNoListIsEverNull(t *testing.T) {
	// An empty plan: the lists are empty, and they must say so.
	empty := newPlan(t)
	for _, args := range listCommands() {
		checkNoNulls(t, empty, args...)
	}

	// A populated plan: the same rule, with lists that have something in them.
	populated := newPlan(t)
	cli(t, populated, "recipe", "create", "--title", "A recipe",
		"--outcome", "o", "--problem", "p",
		"--out-of-scope", "a", "--out-of-scope", "b")
	cli(t, populated, "ac", "add", "R01", "--description", "d", "--verification", "go test ./...")
	cli(t, populated, "wp", "add", "R01", "--name", "core")
	cli(t, populated, "step", "add", "R01", "--action", "do it", "--done-when", "done")
	cli(t, populated, "special", "create", "--title", "A hunch", "--body", "worth noticing")
	cli(t, populated, "prepnote", "create", "--title", "A note", "--intent", "first")
	for _, args := range listCommands() {
		checkNoNulls(t, populated, args...)
	}
	checkNoNulls(t, populated, "recipe", "show", "R01")
}

// checkNoNulls runs a command and asserts that no field in its document is null.
func checkNoNulls(t *testing.T, dir string, args ...string) {
	t.Helper()
	res := cli(t, dir, args...)
	if !res.ok(t) {
		t.Fatalf("scale %v failed:\n%s", args, res.Stderr)
	}
	var doc any
	if err := json.Unmarshal([]byte(res.Stdout), &doc); err != nil {
		t.Fatalf("scale %v did not print a document: %v\n%s", args, err, res.Stdout)
	}
	var nulls []string
	findNulls(doc, "", &nulls)
	if len(nulls) > 0 {
		t.Fatalf("scale %v printed nulls: %v", args, nulls)
	}
}

// findNulls collects the paths of every null in a document.
func findNulls(v any, path string, out *[]string) {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			findNulls(val, path+"/"+k, out)
		}
	case []any:
		for i, val := range t {
			findNulls(val, fmt.Sprintf("%s/%d", path, i), out)
		}
	case nil:
		*out = append(*out, path)
	}
}
