package cli

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The wire between the **CLI's dispatch table** and the **boundary**.
//
// The CLI's verbs are the plan's records and nothing else. The crew is runtime
// state, and the runtime is herdr's — so no verb may start a crew, stop one,
// rename a lane's tab, or read the crew's state. The automation is started from
// the console, and the console writes nothing to the plan when it does.
//
// The wire reads the CLI's **own dispatch table**, the way the count wire does,
// so a verb that reached for the automation is caught where it is added.

// classifiedGroups is the dispatch table, as it is.
//
// Every verb is a **plan record** — a record of the plan moves — except the two
// groups that are read-only surfaces over the tool's own files: `skill` (the
// embedded skills) and `settings` (the tool's settings record).
//
// A verb that is not here is a verb this wire has not seen. That is the point:
// a new verb is added here on purpose, and the wire says what it moves.
var classifiedGroups = map[string][]string{
	"recipe":   {"create", "list", "show", "update", "advance", "approve", "request-changes", "claim", "document", "review", "comment", "block", "unblock", "defer", "resume", "cancel"},
	"skill":    {"help", "list", "load"},
	"ac":       {"add", "verify", "signoff", "defer"},
	"wp":       {"add"},
	"step":     {"add", "start", "done", "skip", "fail", "deviate"},
	"finding":  {"add", "resolve", "dismiss", "reopen"},
	"settings": {"show"},
	"backlog":  {"list", "create", "update", "dismiss", "promote"},
}

// cliWorkspaceCommands are the CLI's own commands, outside any group.
var cliWorkspaceCommands = []string{"init", "demo", "gates", "plan", "path", "log"}

// dispatchGroup finds a group in the dispatch table.
//
// The table is read from the CLI's own source, the way the command-count wire
// reads it: a table held in a test would be a copy, and a copy drifts.
var dispatchGroup = regexp.MustCompile(`(?s)r\.dispatch\(args\[1:\], "(\w+)", map\[string\]func\(\[\]string\) int\{(.*?)\n\t\t\}\)`)

// dispatchVerb finds a verb inside a group's map.
var dispatchVerb = regexp.MustCompile(`"[\w-]+":`)

// workspaceCommand finds a command outside any group.
var workspaceCommand = regexp.MustCompile(`case "(\w+)":\n\t\treturn r\.cmd`)

// TestTheDispatchTableIsClassified holds the CLI's verb set.
//
// Every verb the CLI has is classified here, and the table is read from the CLI's
// own source — so a verb that appears without a line in this file fails the
// check, and the file says what a verb is for.
//
// A verb that reached for the automation would be a line here too, and the next
// check is what holds it.
func TestTheDispatchTableIsClassified(t *testing.T) {
	groups, workspace := readDispatchTable(t)

	for name, verbs := range groups {
		want, known := classifiedGroups[name]
		if !known {
			t.Errorf("the group %q is not classified — say what it moves in classifiedGroups", name)
			continue
		}
		if !sameVerbs(verbs, want) {
			t.Errorf("the group %q has the verbs %v, and it is classified as %v",
				name, verbs, want)
		}
	}
	for name := range classifiedGroups {
		if _, known := groups[name]; !known {
			t.Errorf("the group %q is classified, and the CLI has no such group", name)
		}
	}
	if !sameVerbs(workspace, cliWorkspaceCommands) {
		t.Errorf("the CLI's own commands are %v, and they are classified as %v",
			workspace, cliWorkspaceCommands)
	}
}

// TestNoVerbControlsTheAutomation holds the boundary.
//
// The crew is runtime state: herdr carries it, the console starts it and stops
// it, and the plan carries no record of it. So no verb of the CLI starts a crew,
// stops one, renames a lane's tab, or reads the crew's state.
//
// The check is on the words a verb is made of. A verb that reaches for the
// automation names it — the crew, the runtime, the tab, the round.
func TestNoVerbControlsTheAutomation(t *testing.T) {
	groups, workspace := readDispatchTable(t)

	for _, name := range sortedNames(groups) {
		for _, verb := range groups[name] {
			for _, word := range automationWords {
				if strings.Contains(name, word) {
					t.Errorf("the group %q names the automation (%q) — the CLI holds the plan's records", name, word)
				}
				if strings.Contains(verb, word) {
					t.Errorf("the verb %q %q names the automation (%q) — the CLI holds the plan's records",
						name, verb, word)
				}
			}
		}
	}
	for _, cmd := range workspace {
		for _, word := range automationWords {
			if strings.Contains(cmd, word) {
				t.Errorf("the command %q names the automation (%q) — the CLI holds the plan's records", cmd, word)
			}
		}
	}
}

// TestTheAutomationWireBites keeps the wire honest.
//
// A check that cannot fail is not a check. So the boundary is proved on a table
// that does reach for the automation — the crew's own verbs — and the check must
// catch every one of them.
func TestTheAutomationWireBites(t *testing.T) {
	groups := map[string][]string{
		"crew": {"start", "stop"},
	}
	if got := automationOffenders(groups); len(got) == 0 {
		t.Fatal("the wire should catch a group that starts and stops a crew")
	}

	groups = map[string][]string{
		"recipe": {"advance"},
		"sous":   {"start"},
	}
	if got := automationOffenders(groups); len(got) == 0 {
		t.Fatal("the wire should catch a group that starts the automation")
	}

	groups = map[string][]string{
		"recipe": {"advance"},
	}
	if got := automationOffenders(groups); len(got) != 0 {
		t.Fatalf("the wire should hold a table that holds the plan's records, and it caught %v", got)
	}
}

// automationWords are the words that name the automation.
//
// A verb made of one of them is a verb that reaches for the crew, the runtime or
// the console — and the CLI has no such verb.
//
// Three verbs that sound like the automation's are not here, because the plan
// has verbs of its own with those names: "start" (`step start` starts a guidance
// step), "run" (a recipe is run by its executor) and "kill" (which only the
// runtime offers). A verb named after one of them still has to be classified
// above, and the classification is where the judgement is made.
var automationWords = []string{
	"agent", "automation", "crew", "herdr", "label", "launch", "pane", "rename",
	"restart", "round", "sous", "spawn", "stop", "tab",
}

// automationOffenders are the verbs in the table that name the automation.
func automationOffenders(groups map[string][]string) []string {
	var out []string
	for name, verbs := range groups {
		for _, verb := range verbs {
			for _, word := range automationWords {
				if strings.Contains(name, word) || strings.Contains(verb, word) {
					out = append(out, name+" "+verb)
					break
				}
			}
		}
	}
	sort.Strings(out)
	return out
}

// TestTheDocsSayTheSame holds the boundary in the docs.
//
// The boundary is written in three places — the skill, the wire and the docs —
// and a fact written in three places is a fact that drifts. So the docs are held
// to the same words: `docs/03` says which surface is which, and `AGENTS.md`'s
// layer table names the CLI as the artifact's surface and the console as the
// automation's.
func TestTheDocsSayTheSame(t *testing.T) {
	doc := readOwnFile(t, filepath.Join("..", "..", "docs", "03-domain-model.md"))
	architecture := sectionOf(doc, "### 0.1 The architecture")
	if architecture == "" {
		t.Fatal("docs/03 has no architecture section")
	}
	for _, want := range []string{
		"maitre-cli",
		"the console",
		"herdr",
		"The CLI never controls the automation",
		"never by editing the plan's files",
	} {
		if !strings.Contains(architecture, want) {
			t.Errorf("docs/03's architecture should say %q", want)
		}
	}

	agents := readOwnFile(t, filepath.Join("..", "..", "AGENTS.md"))
	for _, want := range []string{
		"`internal/cli` — the **artifact's** surface",
		"`internal/ui` — the **automation's** surface",
	} {
		if !strings.Contains(agents, want) {
			t.Errorf("AGENTS.md's layer table should say %q", want)
		}
	}
}

// sectionOf returns the part of a doc between a heading and the next one.
func sectionOf(doc, heading string) string {
	at := strings.Index(doc, heading)
	if at < 0 {
		return ""
	}
	rest := doc[at:]
	if end := strings.Index(rest, "\n## "); end >= 0 {
		return rest[:end]
	}
	return rest
}

// readDispatchTable reads the CLI's own dispatch table.
//
// The table is in the source, so the table is read from the source: a copy here
// would be a second source of truth, and the second source of truth is what the
// house's lessons are about.
func readDispatchTable(t *testing.T) (groups map[string][]string, workspace []string) {
	t.Helper()
	source := readOwnSource(t, "cli.go")

	groups = map[string][]string{}
	for _, m := range dispatchGroup.FindAllStringSubmatch(source, -1) {
		var verbs []string
		for _, v := range dispatchVerb.FindAllString(m[2], -1) {
			verbs = append(verbs, strings.Trim(v, `":`))
		}
		groups[m[1]] = verbs
	}
	for _, m := range workspaceCommand.FindAllStringSubmatch(source, -1) {
		workspace = append(workspace, m[1])
	}
	if len(groups) == 0 || len(workspace) == 0 {
		t.Fatal("the dispatch table did not parse")
	}
	return groups, workspace
}

// sameVerbs compares two sets of verbs, whatever their order.
func sameVerbs(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	g := append([]string(nil), got...)
	w := append([]string(nil), want...)
	sort.Strings(g)
	sort.Strings(w)
	for i := range g {
		if g[i] != w[i] {
			return false
		}
	}
	return true
}

// sortedNames are the group's names, in order.
func sortedNames(groups map[string][]string) []string {
	out := make([]string, 0, len(groups))
	for name := range groups {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
