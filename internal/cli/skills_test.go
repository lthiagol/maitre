package cli

import (
	"encoding/json"
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
	skills, _ := skillsFrom(embeddedSkills)
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
	skills, _ := skillsFrom(tree)
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

// TestTheEmbeddedSkillsAgreeWithTheirDirectories holds the rule on the real tree.
//
// Every embedded skill's frontmatter must carry the name of the directory it is in.
// A skill that answers to two names is a skill nobody can address.
func TestTheEmbeddedSkillsAgreeWithTheirDirectories(t *testing.T) {
	skills, refused := skillsFrom(embeddedSkills)
	if len(refused) > 0 {
		for _, ref := range refused {
			t.Errorf("the embedded skill %s is not served: %s", ref.Dir, ref.Reason)
		}
	}
	for _, s := range skills {
		raw, err := fs.ReadFile(embeddedSkills, path.Join("skills", s.Name, "SKILL.md"))
		if err != nil {
			t.Fatal(err)
		}
		if got := frontmatterOf(t, raw)["name"]; got != s.Name {
			t.Errorf("the skill in %s says its name is %q", s.Name, got)
		}
	}
}

// TestASkillThatDisagreesWithItsDirectoryIsRefused holds the rule itself.
//
// The directory's name is the skill's name. A skill whose frontmatter says another
// one is a mistake, and the CLI will not serve it.
func TestASkillThatDisagreesWithItsDirectoryIsRefused(t *testing.T) {
	tree := fstest.MapFS{
		"skills/wrong-name/SKILL.md": &fstest.MapFile{
			Data: []byte("---\nname: other-name\ndescription: It disagrees.\n---\nThe body.\n"),
		},
	}
	skills, refused := skillsFrom(tree)
	if len(skills) != 0 {
		t.Fatalf("the skill should not be served: %v", skillNames(skills))
	}
	if len(refused) != 1 {
		t.Fatalf("got %d refusals, want 1", len(refused))
	}
	if refused[0].Dir != "wrong-name" {
		t.Errorf("the refusal should name the directory: %v", refused[0])
	}
	if !strings.Contains(refused[0].Reason, "other-name") {
		t.Errorf("the refusal should say what the frontmatter says: %v", refused[0])
	}
}

// TestTheCLIDoesNotServeASkillThatDisagrees holds the rule at the surface.
//
// The list does not carry it, and the load refuses it — by the name it claims, and
// by the name of its directory.
func TestTheCLIDoesNotServeASkillThatDisagrees(t *testing.T) {
	tree := fstest.MapFS{
		"skills/wrong-name/SKILL.md": &fstest.MapFile{
			Data: []byte("---\nname: other-name\ndescription: It disagrees.\n---\nThe body.\n"),
		},
		"skills/right-name/SKILL.md": &fstest.MapFile{
			Data: []byte("---\nname: right-name\ndescription: It agrees.\n---\nThe body.\n"),
		},
	}
	old := skillsFS
	skillsFS = tree
	t.Cleanup(func() { skillsFS = old })

	res := cli(t, t.TempDir(), "skill", "list")
	if !res.ok(t) {
		t.Fatalf("skill list failed:\n%s", res.Stderr)
	}
	for _, item := range res.doc(t)["skills"].([]any) {
		if name := item.(map[string]any)["name"]; name == "other-name" {
			t.Errorf("the list should not carry the refused skill: %v", item)
		}
	}
	if !strings.Contains(res.Stderr, "not served") {
		t.Errorf("the list should say the skill is not served:\n%s", res.Stderr)
	}

	// And the load refuses it, by the name it claims.
	res = cli(t, t.TempDir(), "skill", "load", "other-name")
	if res.Code != ExitError {
		t.Fatalf("exit = %d, want %d", res.Code, ExitError)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(res.Stderr), &doc); err != nil {
		t.Fatalf("the refusal should be a document: %v\n%s", err, res.Stderr)
	}
	if doc["error"] != "the skill is not served" {
		t.Errorf("the refusal should say the skill is not served: %v", doc)
	}
	if !strings.Contains(res.Stderr, "other-name") {
		t.Errorf("the refusal should say what the frontmatter says:\n%s", res.Stderr)
	}

	// And the skill that agrees is served.
	res = cli(t, t.TempDir(), "skill", "load", "right-name")
	if res.Code != ExitOK {
		t.Fatalf("the skill that agrees should be served:\n%s", res.Stderr)
	}
}

// TestTheDescriptionRoutes holds the skill's description.
//
// An agent reads the description before it decides, so a description that says only
// "the orchestrator skill" is a description that cannot route. It says what the
// skill is for, when to load it, and what it is not for.
func TestTheDescriptionRoutes(t *testing.T) {
	raw, err := fs.ReadFile(embeddedSkills, path.Join("skills", "orchestrator-role", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	description := frontmatterOf(t, raw)["description"]
	if len(strings.TrimSpace(description)) < 80 {
		t.Errorf("the description should route, not title:\n%s", description)
	}
	for _, want := range []string{
		// What it is for.
		"What:",
		"the orchestrator's role",
		// When to load it.
		"When:",
		"load it when you are about to act as the orchestrator",
		// What it is not for.
		"Near misses:",
		"authoring a recipe",
		"answering the open questions",
	} {
		if !strings.Contains(description, want) {
			t.Errorf("the description should carry %q:\n%s", want, description)
		}
	}
}

// TestTheNamelessSkillIsRefusedForItsOwnReason holds the refusal.
//
// A skill whose frontmatter has no name is refused **for that reason** — not by the
// directory rule catching it by accident, and not with an empty name in the middle
// of a sentence.
func TestTheNamelessSkillIsRefusedForItsOwnReason(t *testing.T) {
	tree := fstest.MapFS{
		"skills/probe-noname/SKILL.md": &fstest.MapFile{
			Data: []byte("---\ndescription: It has no name at all.\n---\nThe body.\n"),
		},
	}
	skills, refused := skillsFrom(tree)
	if len(skills) != 0 {
		t.Fatalf("the skill should not be served: %v", skillNames(skills))
	}
	if len(refused) != 1 {
		t.Fatalf("got %d refusals, want 1", len(refused))
	}
	if !strings.Contains(refused[0].Reason, "no name") {
		t.Errorf("the refusal should say what is missing: %q", refused[0].Reason)
	}
	if strings.Contains(refused[0].Reason, "directory") {
		t.Errorf("the refusal should be the nameless one, not the directory rule: %q", refused[0].Reason)
	}
}

// The CLI says the same thing the loader does.
func TestTheCLISaysWhyTheNamelessSkillIsNotServed(t *testing.T) {
	tree := fstest.MapFS{
		"skills/probe-noname/SKILL.md": &fstest.MapFile{
			Data: []byte("---\ndescription: It has no name at all.\n---\nThe body.\n"),
		},
	}
	old := skillsFS
	skillsFS = tree
	t.Cleanup(func() { skillsFS = old })

	res := cli(t, t.TempDir(), "skill", "list")
	if !res.ok(t) {
		t.Fatalf("skill list failed:\n%s", res.Stderr)
	}
	if !strings.Contains(res.Stderr, "no name") {
		t.Errorf("the warning should say what is missing:\n%s", res.Stderr)
	}
}

// TestTheSkillWithNoDescriptionIsRefused holds the description's absence.
//
// The list is what an agent reads to route, and an entry it cannot describe is an
// entry it cannot route on — so a skill with no description is refused, with the
// reason naming what is missing.
func TestTheSkillWithNoDescriptionIsRefused(t *testing.T) {
	tree := fstest.MapFS{
		"skills/probe-nodesc/SKILL.md": &fstest.MapFile{
			Data: []byte("---\nname: probe-nodesc\n---\nThe body.\n"),
		},
	}
	skills, refused := skillsFrom(tree)
	if len(skills) != 0 {
		t.Fatalf("the skill should not be served: %v", skillNames(skills))
	}
	if len(refused) != 1 {
		t.Fatalf("got %d refusals, want 1", len(refused))
	}
	if !strings.Contains(refused[0].Reason, "no description") {
		t.Errorf("the refusal should say what is missing: %q", refused[0].Reason)
	}
}

// TestTheSkillWithNoDescriptionIsNotInTheList holds it at the surface.
//
// The list carries the skills an agent can route on — and no others.
func TestTheSkillWithNoDescriptionIsNotInTheList(t *testing.T) {
	tree := fstest.MapFS{
		"skills/probe-nodesc/SKILL.md": &fstest.MapFile{
			Data: []byte("---\nname: probe-nodesc\n---\nThe body.\n"),
		},
	}
	old := skillsFS
	skillsFS = tree
	t.Cleanup(func() { skillsFS = old })

	res := cli(t, t.TempDir(), "skill", "list")
	if !res.ok(t) {
		t.Fatalf("skill list failed:\n%s", res.Stderr)
	}
	for _, item := range res.doc(t)["skills"].([]any) {
		if name := item.(map[string]any)["name"]; name == "probe-nodesc" {
			t.Errorf("the list should not carry the refused skill: %v", item)
		}
	}
	if !strings.Contains(res.Stderr, "no description") {
		t.Errorf("the warning should say what is missing:\n%s", res.Stderr)
	}
}

// TestTheSkillListTakesBothSkillsFromTheSkills holds the list's names.
//
// The names and the descriptions come from the skills themselves — proved with **two**
// skills in the tree, so a list written in the CLI cannot pass by naming the one
// skill the tree happens to hold.
func TestTheSkillListTakesBothSkillsFromTheSkills(t *testing.T) {
	tree := fstest.MapFS{
		"skills/first/SKILL.md": &fstest.MapFile{
			Data: []byte("---\nname: first\ndescription: The first probe.\n---\nThe body.\n"),
		},
		"skills/second/SKILL.md": &fstest.MapFile{
			Data: []byte("---\nname: second\ndescription: The second probe.\n---\nThe body.\n"),
		},
	}
	old := skillsFS
	skillsFS = tree
	t.Cleanup(func() { skillsFS = old })

	res := cli(t, t.TempDir(), "skill", "list")
	if !res.ok(t) {
		t.Fatalf("skill list failed:\n%s", res.Stderr)
	}
	items := res.doc(t)["skills"].([]any)
	if len(items) != 2 {
		t.Fatalf("got %d skills, want 2: %v", len(items), items)
	}
	want := map[string]string{
		"first":  "The first probe.",
		"second": "The second probe.",
	}
	for _, item := range items {
		entry := item.(map[string]any)
		name, _ := entry["name"].(string)
		description, _ := entry["description"].(string)
		expected, ok := want[name]
		if !ok {
			t.Errorf("the list should carry the skills of the tree, and it carries %q", name)
			continue
		}
		if description != expected {
			t.Errorf("the skill %s lists %q, and the skill says %q", name, description, expected)
		}
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

// embeddedSkillsList is every embedded skill, and the test fails when there is
// none.
func embeddedSkillsList(t *testing.T) []skill {
	t.Helper()
	skills, _ := skillsFrom(embeddedSkills)
	if len(skills) == 0 {
		t.Fatal("no skill is embedded")
	}
	return skills
}

// embeddedSkill is one embedded skill, and the test fails when it is missing.
func embeddedSkill(t *testing.T, name string) skill {
	t.Helper()
	for _, s := range embeddedSkillsList(t) {
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

// TestTheSkillSaysTheOrchestratorKeepsTheLabelsCurrent holds the role's own
// knowledge.
//
// The label is how the console finds the crew and reads its progress, and the
// orchestrator is the one that writes it: it dispatches the crew, so it is the
// one that knows the round and the interaction count. The skill is the role's
// knowledge, so the skill says so.
func TestTheSkillSaysTheOrchestratorKeepsTheLabelsCurrent(t *testing.T) {
	body := string(embeddedSkill(t, "orchestrator-role").Content)

	for _, want := range []string{
		"The roster's labels",
		"maitre: orchestrator",
		"round <N>",
		"<M> interactions",
		"maitre: runner",
		"herdr tab rename",
		"You write the progress",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("the skill should carry %q", want)
		}
	}
}
