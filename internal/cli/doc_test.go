package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestTheDocsCommandCountMatchesTheCode keeps the headline honest.
//
// The docs and AGENTS.md both say how many commands the CLI has. The number is a
// count of the dispatch table, so it can be counted — and a doc that says "37 verbs"
// beside a table of forty is a doc nobody can trust.
func TestTheDocsCommandCountMatchesTheCode(t *testing.T) {
	total, verbs, groups, workspace := countCommands(t)

	for _, f := range []struct {
		name string
		path string
	}{
		{"docs/04-operations.md", filepath.Join("..", "..", "docs", "04-operations.md")},
		{"AGENTS.md", filepath.Join("..", "..", "AGENTS.md")},
	} {
		text := readOwnFile(t, f.path)
		if err := checkCommandCount(text, total, verbs, groups, workspace); err != nil {
			t.Errorf("%s: %v — the code has %d verbs in %d groups and %d workspace commands",
				f.name, err, verbs, groups, workspace)
		}
	}
}

// commandCount is an occurrence of a command count in a doc.
//
// Two counts are ours: the total, and the number of workspace commands in it. A count
// about v1 is written in another shape — "47 top-level commands" — so it does not
// read as a claim about this CLI.
var commandCount = regexp.MustCompile(`(\d+) (workspace )?commands`)

// commandGroups is an occurrence of the verbs-and-groups count.
//
// The same claim as the total, told in parts: how many verbs, in how many groups. A
// doc that says "42 commands" and "40 verbs in 7 groups" is a doc whose two halves
// disagree, and the total alone would not catch it.
var commandGroups = regexp.MustCompile(`(\d+) verbs in (\d+) groups`)

// checkCommandCount requires every occurrence of a count to agree with the code.
//
// The count is written in more than one place — AGENTS.md says it twice — and a fact
// written in two places is a bug waiting to happen. So every occurrence has to agree,
// not just one of them.
func checkCommandCount(text string, total, verbs, groups, workspace int) error {
	headline := fmt.Sprintf("%d commands", total)
	matches := commandCount.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return fmt.Errorf("the file does not say %q", headline)
	}
	for _, m := range matches {
		want := fmt.Sprintf("%d commands", total)
		if m[2] != "" {
			want = fmt.Sprintf("%d workspace commands", workspace)
		}
		if m[0] != want {
			return fmt.Errorf("the file says %q in one place and %q in another", m[0], want)
		}
	}
	for _, m := range commandGroups.FindAllStringSubmatch(text, -1) {
		want := fmt.Sprintf("%d verbs in %d groups", verbs, groups)
		if m[0] != want {
			return fmt.Errorf("the file says %q and %q", m[0], want)
		}
	}
	return nil
}

// TestTheCommandCountWireTripsOnAHalfEdit keeps the wire honest.
//
// Changing one of the two occurrences in AGENTS.md must fail the check: a half
// finished edit is exactly what the wire is for.
func TestTheCommandCountWireTripsOnAHalfEdit(t *testing.T) {
	total, verbs, groups, workspace := countCommands(t)
	agents := readOwnFile(t, filepath.Join("..", "..", "AGENTS.md"))

	// The first occurrence only — the half-finished edit the wire is for.
	headline := fmt.Sprintf("%d commands", total)
	edited := strings.Replace(agents, headline, fmt.Sprintf("%d commands", total-1), 1)
	if edited == agents {
		t.Fatalf("the test needs AGENTS.md to say %q", headline)
	}
	if err := checkCommandCount(edited, total, verbs, groups, workspace); err == nil {
		t.Fatal("the wire should fail when one occurrence disagrees")
	}
}

// TestTheCommandCountWireTripsOnAHalfEditOfTheGroups keeps the wire honest about
// the verbs and the groups.
//
// The total and the parts are the same claim. A doc that keeps the total and
// moves the parts is half-edited too — and the wire has to catch it.
func TestTheCommandCountWireTripsOnAHalfEditOfTheGroups(t *testing.T) {
	total, verbs, groups, workspace := countCommands(t)
	agents := readOwnFile(t, filepath.Join("..", "..", "AGENTS.md"))

	headline := fmt.Sprintf("%d verbs in %d groups", verbs, groups)
	edited := strings.Replace(agents, headline, fmt.Sprintf("%d verbs in %d groups", verbs+1, groups+1), 1)
	if edited == agents {
		t.Fatalf("the test needs AGENTS.md to say %q", headline)
	}
	if err := checkCommandCount(edited, total, verbs, groups, workspace); err == nil {
		t.Fatal("the wire should fail when the verbs and the groups disagree")
	}
}

// countCommands counts the commands in the dispatch table.
func countCommands(t *testing.T) (total, verbs, groups, workspace int) {
	t.Helper()
	source := readOwnSource(t, "cli.go")

	// The verbs in the dispatch maps.
	dispatch := regexp.MustCompile(`(?s)r\.dispatch\(args\[1:\], "(\w+)", map\[string\]func\(\[\]string\) int\{(.*?)\n\t\t\}\)`)
	for _, m := range dispatch.FindAllStringSubmatch(source, -1) {
		groups++
		verbs += len(regexp.MustCompile(`"[\w-]+":`).FindAllString(m[2], -1))
	}
	// The workspace commands.
	workspace = len(regexp.MustCompile(`case "\w+":\n\t\treturn r\.cmd`).FindAllString(source, -1))
	if verbs == 0 || workspace == 0 {
		t.Fatal("the dispatch table did not parse")
	}
	return verbs + workspace, verbs, groups, workspace
}

// readOwnSource reads a file from this package.
func readOwnSource(t *testing.T, name string) string {
	t.Helper()
	return readOwnFile(t, name)
}

// readOwnDoc reads one of the project's design docs.
func readOwnDoc(t *testing.T, name string) string {
	t.Helper()
	return readOwnFile(t, filepath.Join("..", "..", "docs", name))
}

// readOwnFile reads a file and fails the test if it cannot.
func readOwnFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}
