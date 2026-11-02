package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

// The documents that state the rule, and the sentence they share with the binary.
var textCommandDocs = []string{
	"../../AGENTS.md",
	"../../docs/04-operations.md",
}

// The documents that state the rule about the verbs, and the exception they share.
var operationDocs = []string{
	"../../AGENTS.md",
	"../../docs/04-operations.md",
}

// The marker of the sentence: the one that names the commands that print text.
const textCommandMarker = "commands print text instead:"

// TestTheTextCommandsAreTheOnesThatPrintText holds the rule against the surface.
//
// The commands that print text are exactly the ones the list names — and the list
// is what the documents and the binary's own help are checked against.
func TestTheTextCommandsAreTheOnesThatPrintText(t *testing.T) {
	isText := map[string]bool{}
	for _, c := range textCommands {
		isText[c.Name] = true

		// Each of them prints text.
		res := cli(t, t.TempDir(), c.argv()...)
		if res.Code != ExitOK {
			t.Fatalf("maitre-cli %s failed:\n%s", c.Name, res.Stderr)
		}
		if strings.TrimSpace(res.Stdout) == "" {
			t.Errorf("maitre-cli %s should print something", c.Name)
			continue
		}
		var doc any
		if err := json.Unmarshal([]byte(res.Stdout), &doc); err == nil {
			t.Errorf("maitre-cli %s should print text, not a document", c.Name)
		}
	}

	// And no other command does.
	for _, cmd := range cliSurface(t) {
		if isText[cmd] {
			continue
		}
		res := cli(t, t.TempDir(), strings.Fields(cmd)...)
		if strings.TrimSpace(res.Stdout) == "" {
			continue
		}
		var doc any
		if err := json.Unmarshal([]byte(res.Stdout), &doc); err != nil {
			t.Errorf("maitre-cli %s prints text, and it is not one of the text commands:\n%s",
				cmd, res.Stdout)
		}
	}
}

// TestTheDocsAndTheHelpSayTheSameThingAboutTheTextCommands holds the three copies.
//
// One list, three copies: the documents and the binary's own help must name exactly
// the commands that print text — and the same count. Put the wrong count back and
// this goes red.
func TestTheDocsAndTheHelpSayTheSameThingAboutTheTextCommands(t *testing.T) {
	want := textCommandNames()

	for _, path := range textCommandDocs {
		count, got := textCommandsIn(t, readOwnFile(t, path))
		if !strings.EqualFold(count, numberWord(len(want))) {
			t.Errorf("%s says %q commands print text; there are %d", path, count, len(want))
		}
		if !sameCommands(got, want) {
			t.Errorf("%s names %v; the text commands are %v", path, got, want)
		}
	}

	// And the binary's own help — the place an agent reads first.
	res := cli(t, t.TempDir(), "help")
	if res.Code != ExitOK {
		t.Fatalf("help failed:\n%s", res.Stderr)
	}
	count, got := textCommandsIn(t, res.Stdout)
	if !strings.EqualFold(count, numberWord(len(want))) {
		t.Errorf("the help says %q commands print text; there are %d", count, len(want))
	}
	if !sameCommands(got, want) {
		t.Errorf("the help names %v; the text commands are %v", got, want)
	}
}

// textCommandsIn reads the sentence that names the text commands out of a text.
//
// The whitespace is normalised first, so a sentence that wraps in the source is
// still one sentence — and the count is read from the word before the marker.
func textCommandsIn(t *testing.T, text string) (count string, commands []string) {
	t.Helper()
	flat := strings.Join(strings.Fields(text), " ")
	idx := strings.Index(flat, textCommandMarker)
	if idx < 0 {
		t.Fatalf("the text does not say which commands print text — %q is missing", textCommandMarker)
	}
	before := strings.Fields(flat[:idx])
	if len(before) == 0 {
		t.Fatalf("the sentence names no count")
	}
	count = before[len(before)-1]

	rest := flat[idx+len(textCommandMarker):]
	if end := strings.Index(rest, "."); end >= 0 {
		rest = rest[:end]
	}
	rest = strings.ReplaceAll(rest, " and ", ",")
	for _, part := range strings.Split(rest, ",") {
		part = strings.TrimSpace(strings.Trim(strings.TrimSpace(part), "`"))
		if part != "" {
			commands = append(commands, part)
		}
	}
	return count, commands
}

// sameCommands compares two sets of command names.
func sameCommands(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	seen := map[string]bool{}
	for _, g := range got {
		seen[g] = true
	}
	for _, w := range want {
		if !seen[w] {
			return false
		}
	}
	return true
}

// TestEveryVerbOutsideTheReadOnlyGroupsMapsToAnOperation holds the exception.
//
// The `skill` group reads the tool's own embedded files, and the `settings` group
// reads the tool's own settings file: both are read-only surfaces, not operations,
// and the documents say so. Every other verb is an adapter over the plan's
// operations, so it reaches the plan.
//
// The exception is asked of the **code**, not of the documents: the groups whose
// verbs never open the plan are the ones the documents must name.
func TestEveryVerbOutsideTheReadOnlyGroupsMapsToAnOperation(t *testing.T) {
	planless := planlessGroups(t)
	if len(planless) == 0 {
		t.Fatal("every group opens the plan — the check found nothing to hold")
	}

	for _, path := range operationDocs {
		text := flatten(readOwnFile(t, path))
		named := exceptionGroups(text)
		if !sameCommands(named, keys(planless)) {
			t.Errorf("%s names %v as the exception; the groups that reach no plan are %v",
				path, named, keys(planless))
		}
		for _, group := range named {
			if !strings.Contains(text, operationExceptionPhrase(group)) {
				t.Errorf("%s does not say what the %s group is: %q",
					path, group, operationExceptionPhrase(group))
			}
		}
	}
}

// operationExceptionPhrase is the sentence the documents share when they state the
// exception.
//
// The two read-only surfaces are not the same surface: the `skill` group reads the
// tool's own embedded files, and the `settings` group reads the tool's own settings
// file. So each has its own sentence, and the documents must carry both.
func operationExceptionPhrase(group string) string {
	if phrase, ok := operationExceptionPhrases[group]; ok {
		return phrase
	}
	return fmt.Sprintf(
		"except the `%s` group, which is a read-only surface over the tool's own embedded files and is not an operation",
		group)
}

// operationExceptionPhrases are the sentences the documents share, by group.
var operationExceptionPhrases = map[string]string{
	"skill":    "except the `skill` group, which is a read-only surface over the tool's own embedded files and is not an operation",
	"settings": "except the `settings` group, which is a read-only surface over the tool's own settings file and is not an operation",
}

// exceptionGroups are the groups a document names as the exception.
func exceptionGroups(flat string) []string {
	var out []string
	for _, m := range regexp.MustCompile("except the `([\\w-]+)` group").FindAllStringSubmatch(flat, -1) {
		out = append(out, m[1])
	}
	return out
}

// flatten normalises the whitespace, so a sentence that wraps in the source is still
// one sentence.
func flatten(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

// keys are the names of a set.
func keys(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	return out
}

// cliSurface is the command surface: the workspace commands, the verbs of every
// group, and the help.
//
// It is read from the dispatch table, so a command added is a command the surface
// knows and a command renamed is one it stops knowing.
func cliSurface(t *testing.T) []string {
	t.Helper()
	source := readOwnFile(t, "cli.go")
	var out []string

	// The workspace commands.
	for _, m := range regexp.MustCompile(`case "(\w+)":\n\t\treturn r\.cmd`).FindAllStringSubmatch(source, -1) {
		out = append(out, m[1])
	}
	// The verbs of every group.
	for _, m := range regexp.MustCompile(`(?s)r\.dispatch\(args\[1:\], "(\w+)", map\[string\]func\(\[\]string\) int\{(.*?)\n\t\t\}\)`).FindAllStringSubmatch(source, -1) {
		for _, v := range regexp.MustCompile(`"([\w-]+)":`).FindAllStringSubmatch(m[2], -1) {
			out = append(out, m[1]+" "+v[1])
		}
	}
	// The help.
	out = append(out, "help")
	return out
}

// cliVerb is one verb of a group.
type cliVerb struct {
	// Verb is what the user types.
	Verb string
	// Func is the Runner method that serves it.
	Func string
}

// cliGroups are the verbs of every group, by group.
func cliGroups(t *testing.T) map[string][]cliVerb {
	t.Helper()
	source := readOwnFile(t, "cli.go")
	out := map[string][]cliVerb{}
	for _, m := range regexp.MustCompile(`(?s)r\.dispatch\(args\[1:\], "(\w+)", map\[string\]func\(\[\]string\) int\{(.*?)\n\t\t\}\)`).FindAllStringSubmatch(source, -1) {
		for _, v := range regexp.MustCompile(`"([\w-]+)":\s+r\.(cmd\w+)`).FindAllStringSubmatch(m[2], -1) {
			out[m[1]] = append(out[m[1]], cliVerb{Verb: v[1], Func: v[2]})
		}
	}
	return out
}

// planlessGroups are the groups whose verbs never open the plan.
//
// A verb that maps to an operation reaches the plan — the operations live there —
// so this asks the code, not the documents.
func planlessGroups(t *testing.T) map[string]bool {
	t.Helper()
	out := map[string]bool{}
	for group, verbs := range cliGroups(t) {
		for _, v := range verbs {
			if !reachesThePlan(t, v.Func, map[string]bool{}) {
				out[group] = true
				break
			}
		}
	}
	return out
}

// reachesThePlan reports whether a method opens the plan.
//
// The calls are followed, because a verb may reach the plan through a helper.
func reachesThePlan(t *testing.T, method string, seen map[string]bool) bool {
	t.Helper()
	if seen[method] {
		return false
	}
	seen[method] = true
	body, ok := runnerMethods(t)[method]
	if !ok {
		return false
	}
	if strings.Contains(body, "openStore") {
		return true
	}
	for _, m := range regexp.MustCompile(`r\.(\w+)\(`).FindAllStringSubmatch(body, -1) {
		if reachesThePlan(t, m[1], seen) {
			return true
		}
	}
	return false
}

// runnerMethods are the Runner's methods, by name.
func runnerMethods(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, chunk := range strings.Split(packageSources(t), "\nfunc ") {
		head := regexp.MustCompile(`^\(\w+ \*?(\w+)\) (\w+)\(`).FindStringSubmatch(chunk)
		if head == nil {
			continue
		}
		out[head[2]] = chunk
	}
	return out
}

// packageSources are the CLI's own sources, read to see what the verbs reach.
func packageSources(t *testing.T) string {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		out.Write(raw)
		out.WriteString("\n")
	}
	return out.String()
}
