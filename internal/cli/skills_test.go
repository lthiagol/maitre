package cli

import (
	"io/fs"
	"path"
	"strings"
	"testing"
	"testing/fstest"
)

// TestTheOrchestratorRoleSkillIsEmbedded holds the first skill in place.
//
// The skill is the role the spawned orchestrator loads, and it ships inside the
// binary — so it is here that it has to be, in the embedded tree and not on disk.
func TestTheOrchestratorRoleSkillIsEmbedded(t *testing.T) {
	skills := skillsFrom(embeddedSkills)
	var found *skill
	for i := range skills {
		if skills[i].Name == "orchestrator-role" {
			found = &skills[i]
		}
	}
	if found == nil {
		t.Fatalf("the orchestrator-role skill is not embedded; the embedded skills are %v", skillNames(skills))
	}
	if found.Description == "" {
		t.Error("the skill should carry a description")
	}
}

// TestTheOrchestratorRoleSkillCarriesTheRoleAndTheProtocol holds what the skill is
// for.
//
// The chef's words: the skill describes *what an orchestrator does and what it
// should know to act as an orchestrator*. So it carries the role, the rules the
// orchestrator works by, the protocol, the herdr mechanics — including the two
// commands that do not work — and what the orchestrator must not do.
func TestTheOrchestratorRoleSkillCarriesTheRoleAndTheProtocol(t *testing.T) {
	body := string(embeddedSkill(t, "orchestrator-role").Content)

	for _, want := range []string{
		// The role.
		"calls the recipes, routes them, and keeps the line moving",
		// What it must know.
		"the brief is the contract",
		"one commit per acceptance criterion",
		"stage by name",
		"`make check` is the gate",
		"the receipt is the record",
		"the open questions belong to the chef",
		// The protocol.
		"Write the brief",
		"Open the lanes",
		"Dispatch the runner",
		"Wait",
		"Receive the receipt",
		"Brief the reviewer",
		"Close the slice",
		"The ledger",
		// The herdr mechanics, and the two commands that do not work.
		"herdr pane send-text",
		"herdr agent wait",
		"herdr agent prompt",
		"agent_not_ready",
		"herdr agent start",
		"no mcode kind",
		// What it does not do.
		"author a recipe",
		"answer the open questions",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("the skill should carry %q", want)
		}
	}
}

// TestTheSkillCarriesNoWork holds the line between the role and the work.
//
// The skill is the role; the startup task is the work. A skill that named the
// recipes to run would carry one plan's work into every plan that loaded it.
func TestTheSkillCarriesNoWork(t *testing.T) {
	body := string(embeddedSkill(t, "orchestrator-role").Content)

	if strings.Contains(body, "RCP") {
		t.Error("the skill should not name a recipe — the startup task carries the work")
	}
	for _, want := range []string{
		"The skill is the role; the startup task is the work",
		"which recipes to run",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("the skill should say %q", want)
		}
	}
}

// TestTheSkillHelpSaysWhatTheMenuIs holds the menu's help.
//
// The help is text — it is read, not parsed. And it says what the menu is, what the
// three verbs do, and that the skills ship inside the binary.
func TestTheSkillHelpSaysWhatTheMenuIs(t *testing.T) {
	res := cli(t, t.TempDir(), "skill", "help")
	if res.Code != ExitOK {
		t.Fatalf("skill help failed:\n%s", res.Stderr)
	}
	if strings.HasPrefix(strings.TrimSpace(res.Stdout), "{") {
		t.Fatalf("the help should be text, not a document:\n%s", res.Stdout)
	}
	for _, want := range []string{
		"maitre-cli skill",
		"maitre-cli skill help",
		"maitre-cli skill list",
		"maitre-cli skill load",
		"inside the binary",
		"embedded",
	} {
		if !strings.Contains(res.Stdout, want) {
			t.Errorf("the help should say %q:\n%s", want, res.Stdout)
		}
	}
}

// TestTheSkillListIsTheEmbeddedTree holds the list against the tree.
//
// The list is computed, not written: every skill in the embedded tree is in the list,
// and nothing else is.
func TestTheSkillListIsTheEmbeddedTree(t *testing.T) {
	entries, err := fs.ReadDir(embeddedSkills, "skills")
	if err != nil {
		t.Fatal(err)
	}
	want := 0
	for _, e := range entries {
		if e.IsDir() {
			want++
		}
	}

	res := cli(t, t.TempDir(), "skill", "list")
	if !res.ok(t) {
		t.Fatalf("skill list failed:\n%s", res.Stderr)
	}
	doc := res.doc(t)
	items, ok := doc["skills"].([]any)
	if !ok {
		t.Fatalf("the list should carry the skills: %v", doc)
	}
	if len(items) != want {
		t.Fatalf("got %d skills, want %d — the tree has %d", len(items), want, want)
	}
	if doc["count"] != float64(want) {
		t.Errorf("count = %v, want %d", doc["count"], want)
	}
}

// TestTheSkillListTakesTheNameAndDescriptionFromTheSkill holds where the list's
// facts come from.
//
// A skill's name and description come from the skill itself — there is no second
// list to keep in step.
func TestTheSkillListTakesTheNameAndDescriptionFromTheSkill(t *testing.T) {
	res := cli(t, t.TempDir(), "skill", "list")
	if !res.ok(t) {
		t.Fatalf("skill list failed:\n%s", res.Stderr)
	}
	items := res.doc(t)["skills"].([]any)
	if len(items) == 0 {
		t.Fatal("the list is empty")
	}
	for _, item := range items {
		entry := item.(map[string]any)
		name, _ := entry["name"].(string)
		if name == "" {
			t.Fatalf("every skill needs a name: %v", entry)
		}
		raw, err := fs.ReadFile(embeddedSkills, path.Join("skills", name, "SKILL.md"))
		if err != nil {
			t.Fatalf("the skill %s is not in the embedded tree: %v", name, err)
		}
		fields := frontmatterOf(t, raw)
		if fields["name"] != name {
			t.Errorf("the skill %s says its name is %q", name, fields["name"])
		}
		if got := strings.TrimSpace(entry["description"].(string)); got != strings.TrimSpace(fields["description"]) {
			t.Errorf("the skill %s lists the description %q, and the skill says %q",
				name, got, strings.TrimSpace(fields["description"]))
		}
	}
}

// TestANewSkillAppearsInTheListWithNoOtherEdit holds the rule that there is no
// second list.
//
// A skill is a directory with a SKILL.md inside. Adding one adds it to the list —
// nothing else has to change.
func TestANewSkillAppearsInTheListWithNoOtherEdit(t *testing.T) {
	tree := fstest.MapFS{
		"skills/first/SKILL.md": &fstest.MapFile{
			Data: []byte("---\nname: first\ndescription: The first skill.\n---\nThe body.\n"),
		},
		"skills/second/SKILL.md": &fstest.MapFile{
			Data: []byte("---\nname: second\ndescription: The second skill.\n---\nThe body.\n"),
		},
	}
	skills := skillsFrom(tree)
	got := skillNames(skills)
	if len(got) != 2 || got[0] != "first" || got[1] != "second" {
		t.Fatalf("the list should hold both skills, in order: %v", got)
	}
}

// TestTheSkillLoadPrintsTheContentByteForByte holds what the load prints.
//
// The content is meant to be read, not parsed — so it is printed as embedded, byte
// for byte, and it is not a JSON document.
func TestTheSkillLoadPrintsTheContentByteForByte(t *testing.T) {
	res := cli(t, t.TempDir(), "skill", "load", "orchestrator-role")
	if res.Code != ExitOK {
		t.Fatalf("skill load failed:\n%s", res.Stderr)
	}
	raw, err := fs.ReadFile(embeddedSkills, path.Join("skills", "orchestrator-role", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Stdout != string(raw) {
		t.Errorf("the load should print the skill, byte for byte:\n%s", res.Stdout)
	}
	if strings.HasPrefix(strings.TrimSpace(res.Stdout), "{") {
		t.Errorf("the load should not print a document:\n%s", res.Stdout)
	}
}

// TestTheSkillLoadRefusesWhatIsNotThere holds the typo.
func TestTheSkillLoadRefusesWhatIsNotThere(t *testing.T) {
	res := cli(t, t.TempDir(), "skill", "load", "no-such-skill")
	if res.Code != ExitError {
		t.Fatalf("exit = %d, want %d", res.Code, ExitError)
	}
	if !strings.Contains(res.Stderr, "no-such-skill") {
		t.Fatalf("the refusal should name the skill: %s", res.Stderr)
	}
	if !strings.Contains(res.Stderr, "orchestrator-role") {
		t.Fatalf("the refusal should list what there is: %s", res.Stderr)
	}
}

// frontmatterOf reads a skill's frontmatter, without the production's parser.
//
// A test that reuses the code it is testing proves nothing.
func frontmatterOf(t *testing.T, raw []byte) map[string]string {
	t.Helper()
	out := map[string]string{}
	lines := strings.Split(string(raw), "\n")
	inFrontmatter := false
	key := ""
	for _, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if inFrontmatter {
				break
			}
			inFrontmatter = true
			continue
		}
		if !inFrontmatter {
			continue
		}
		if key != "" {
			if strings.TrimSpace(line) == "" || strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
				out[key] += "\n" + strings.TrimSpace(line)
				continue
			}
			key = ""
		}
		name, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if value == "|" || value == ">" {
			key = name
			out[name] = ""
			continue
		}
		out[name] = value
	}
	return out
}

// embeddedSkill is one embedded skill, and the test fails when it is missing.
func embeddedSkill(t *testing.T, name string) skill {
	t.Helper()
	for _, s := range skillsFrom(embeddedSkills) {
		if s.Name == name {
			return s
		}
	}
	t.Fatalf("the %s skill is not embedded", name)
	return skill{}
}

// skillNames are the names of the given skills.
func skillNames(skills []skill) []string {
	out := make([]string, 0, len(skills))
	for _, s := range skills {
		out = append(out, s.Name)
	}
	return out
}
