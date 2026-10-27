package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTheDocListsTheBacklogKinds keeps the closed enum honest.
//
// The kind decides the group the console files an item under and the shape a
// promotion gives the recipe, so a kind nobody knows is a typo — and a typo that
// silently disabled a feature is v1's lesson. The doc states the set, and the code
// holds it: the two are checked against each other, both ways.
func TestTheDocListsTheBacklogKinds(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "docs", "03-domain-model.md"))
	if err != nil {
		t.Fatal(err)
	}
	section := sectionOf(string(raw), "## 6. Plan format")
	if section == "" {
		t.Fatal("the doc has no § 6. Plan format")
	}

	listed := map[string]bool{}
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| `kind`") {
			continue
		}
		// The second cell holds the kinds; the first holds the field's name.
		cells := strings.Split(line, "|")
		if len(cells) < 3 {
			continue
		}
		for _, name := range allBackticked(cells[2]) {
			listed[name] = true
		}
	}
	if len(listed) == 0 {
		t.Fatalf("the plan format does not state the kinds:\n%s", section)
	}

	var missing, extra []string
	for _, k := range BacklogKinds {
		if !listed[k] {
			missing = append(missing, k)
		}
	}
	for name := range listed {
		if !containsString(BacklogKinds, name) {
			extra = append(extra, name)
		}
	}
	if len(missing) > 0 || len(extra) > 0 {
		t.Errorf("the kinds do not agree: the doc is missing %v and names %v the code does not have", missing, extra)
	}
}

// TestTheDocNamesTheBacklogGroups keeps the groups derived.
//
// The console's two groups are read off the kind — `core.BacklogGroups` — and never
// stored. A doc that named a third group would be a group somewhere else.
func TestTheDocNamesTheBacklogGroups(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "docs", "03-domain-model.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc := string(raw)
	for _, want := range []string{"Ideas", "Small tasks"} {
		if !strings.Contains(doc, want) {
			t.Errorf("the doc should name the %q group", want)
		}
	}
}
