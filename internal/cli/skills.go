package cli

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// The skills ship inside the binary.
//
// The orchestrator is spawned by the tool, and it loads the role the tool carries: a
// skill read from disk would drift from the tool that spawns the agent. `embed` is
// stdlib, so `internal/cli` stays stdlib-only and `internal/core` does not change.
//
//go:embed skills
var embeddedSkills embed.FS

// skillsFS is where the skills are read from.
//
// It is a variable so a test can put another tree in front of it — one with a skill
// whose directory and frontmatter disagree, which the real tree never has.
var skillsFS fs.FS = embeddedSkills

// A skill is a directory with a SKILL.md inside, and the directory's name is the
// skill's name.
type skill struct {
	// Name is the skill's name, taken from its frontmatter.
	Name string
	// Description says what the skill is for and when an agent should load it.
	Description string
	// Content is the SKILL.md, byte for byte.
	Content []byte
}

// A refusal is a skill the CLI will not serve.
type refusal struct {
	// Dir is the directory the skill is in.
	Dir string
	// Name is the name the skill claims, when it claims one.
	Name string
	// Reason says why it is not served.
	Reason string
}

// skillsFrom reads the skills in the tree.
//
// The list is computed, not written: a skill's name and description come from the
// skill itself, so a new embedded skill appears in the list with no other edit.
//
// A skill whose directory and frontmatter disagree is **refused**. That is a
// mistake, not a variant — the name is the directory's, and a skill that answers to
// two names is a skill nobody can address.
func skillsFrom(fsys fs.FS) ([]skill, []refusal) {
	entries, err := fs.ReadDir(fsys, "skills")
	if err != nil {
		return nil, nil
	}
	var skills []skill
	var refused []refusal
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		raw, err := fs.ReadFile(fsys, path.Join("skills", e.Name(), "SKILL.md"))
		if err != nil {
			refused = append(refused, refusal{Dir: e.Name(), Reason: "it has no SKILL.md"})
			continue
		}
		meta, err := parseSkill(raw)
		if err != nil {
			refused = append(refused, refusal{Dir: e.Name(), Reason: err.Error()})
			continue
		}
		if meta.Name != e.Name() {
			refused = append(refused, refusal{
				Dir:    e.Name(),
				Name:   meta.Name,
				Reason: "its directory is " + e.Name() + " and its frontmatter says its name is " + meta.Name,
			})
			continue
		}
		skills = append(skills, skill{
			Name:        meta.Name,
			Description: meta.Description,
			Content:     raw,
		})
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })
	return skills, refused
}

// cmdSkillHelp prints what the skills menu is.
//
// Text, like the top-level help: the help is read, not parsed.
func (r *Runner) cmdSkillHelp(args []string) int {
	fs := newFlagSet("skill help")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	fmt.Fprint(r.Stdout, `maitre-cli skill — the skills that ship inside the binary

A skill is the machine's own instruction sheet for a job: the role, and the protocol
it follows. The skills ship inside the binary — they are embedded at build time — so
the skill and the tool's version cannot drift. The orchestrator the tool spawns
loads the role the tool carries. There is no directory, no PATH and no file on disk:
the menu works from anywhere.

Usage:
  maitre-cli skill help            what this menu is
  maitre-cli skill list            every embedded skill: its name and its description
  maitre-cli skill load <name>     the skill's content, on stdout, as embedded

list prints a JSON document, like every other list the CLI prints. load prints the
skill's content, byte for byte — it is meant to be read, not parsed. help prints
this text.
`)
	return ExitOK
}

// cmdSkillList lists the embedded skills.
//
// The list is computed from the skills themselves: a skill's name and description
// come from the skill, so a new embedded skill appears in the list with no other
// edit.
func (r *Runner) cmdSkillList(args []string) int {
	fs := newFlagSet("skill list")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	items := make([]map[string]string, 0)
	skills, refused := skillsFrom(skillsFS)
	for _, ref := range refused {
		r.warn("the skill " + ref.Dir + " is not served: " + ref.Reason)
	}
	for _, s := range skills {
		items = append(items, map[string]string{
			"name":        s.Name,
			"description": s.Description,
		})
	}
	return r.emit(map[string]any{"ok": true, "count": len(items), "skills": items})
}

// cmdSkillLoad prints a skill's content.
//
// It is the one command that does not print a JSON document: the content is meant to
// be read, not parsed. It is printed byte for byte, as embedded.
func (r *Runner) cmdSkillLoad(args []string) int {
	fs := newFlagSet("skill load")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	if fs.nArg() != 1 {
		r.fail("usage: maitre-cli skill load <name>")
		return ExitError
	}
	name := fs.arg(0)
	skills, refused := skillsFrom(skillsFS)
	for _, ref := range refused {
		r.warn("the skill " + ref.Dir + " is not served: " + ref.Reason)
	}
	for _, s := range skills {
		if s.Name != name {
			continue
		}
		if _, err := r.Stdout.Write(s.Content); err != nil {
			r.fail("write the skill", err.Error())
			return ExitError
		}
		r.flushWarnings()
		return ExitOK
	}
	for _, ref := range refused {
		if ref.Name == name {
			r.fail("the skill is not served", ref.Reason)
			return ExitError
		}
	}
	known := make([]string, 0, len(skills))
	for _, s := range skills {
		known = append(known, s.Name)
	}
	r.fail("no such skill", name, "known skills: "+strings.Join(known, ", "))
	return ExitError
}

// skillMeta is what a skill's frontmatter says.
type skillMeta struct {
	// Name is the skill's name.
	Name string
	// Description says what the skill is for and when to load it.
	Description string
}

// parseSkill reads a skill's frontmatter.
//
// The frontmatter is the block between two `---` lines, and it carries the name and
// the description. A file without one is not a skill the CLI can serve.
func parseSkill(raw []byte) (skillMeta, error) {
	lines := strings.Split(string(raw), "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return skillMeta{}, fmt.Errorf("the file has no frontmatter")
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			end = i
			break
		}
	}
	if end < 0 {
		return skillMeta{}, fmt.Errorf("the frontmatter is not closed")
	}
	fields := frontmatterFields(lines[1:end])
	meta := skillMeta{
		Name:        fields["name"],
		Description: fields["description"],
	}
	if meta.Name == "" {
		return skillMeta{}, fmt.Errorf("the frontmatter has no name")
	}
	return meta, nil
}

// frontmatterFields reads the frontmatter's keys.
//
// A value is either on the line — `name: orchestrator-role` — or a block under it:
// `description: |` and the indented lines that follow. A literal block (`|`) keeps
// its line breaks; a folded one (`>`) joins its lines with spaces.
func frontmatterFields(lines []string) map[string]string {
	out := map[string]string{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line != strings.TrimLeft(line, " \t") {
			// An indented line belongs to the block above it.
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch value {
		case "|", "|-", "|+":
			block, next := blockLines(lines, i+1)
			out[key] = strings.Join(block, "\n")
			i = next - 1
		case ">", ">-", ">+":
			block, next := blockLines(lines, i+1)
			out[key] = strings.Join(block, " ")
			i = next - 1
		default:
			out[key] = strings.Trim(value, `"'`)
		}
	}
	return out
}

// blockLines collects the indented lines under a block scalar.
//
// It returns the lines, without their indentation, and the index of the first line
// that is not part of the block.
func blockLines(lines []string, from int) ([]string, int) {
	var out []string
	i := from
	for ; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			out = append(out, "")
			continue
		}
		if line == strings.TrimLeft(line, " \t") {
			break
		}
		out = append(out, strings.TrimSpace(line))
	}
	return out, i
}
