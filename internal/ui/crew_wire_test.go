package ui

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// The wire between the **skill** and the **console**.
//
// The label's format lives in two places: the **skill** tells the orchestrator
// what to write, and the **console** parses what it wrote. Each side had its own
// witness, and neither witness looked at the other — so the skill could drift,
// the orchestrator was then told to write a format the console cannot read, and
// the console read no progress, with no error anywhere.
//
// This file joins them: it takes the **skill's own example** and feeds it to the
// console's parser. A drift in the skill turns this test red.

// skillPath is the orchestrator's role, from the console's side.
//
// The role is the CLI's file and the console does not import the CLI — but the
// **format** the role teaches is the format the console reads, so the two must
// agree. The path is the wire.
const skillPath = "../cli/skills/orchestrator-role/SKILL.md"

// skillExample is a roster's line the skill shows.
//
// The pattern starts with a letter, so the prose — "`maitre: …`" — is not
// mistaken for an example.
var skillExample = regexp.MustCompile("maitre: [a-z][^`\"\\n]*")

// crewJobs are the jobs the roster's lines name.
//
// The orchestrator's is the one the console reads the progress from, and the
// other two are the crew's other members.
var crewJobs = []string{"orchestrator", "runner", "reviewer"}

// skillExamples are the roster's lines the skill teaches.
//
// The placeholders are filled in: the console reads numbers, and a placeholder is
// not one.
func skillExamples(skill string) []string {
	var out []string
	for _, match := range skillExample.FindAllString(skill, -1) {
		match = strings.ReplaceAll(match, "<N>", "2")
		match = strings.ReplaceAll(match, "<M>", "6")
		out = append(out, match)
	}
	return out
}

// TestTheSkillExampleIsTheFormatTheConsoleReads is the wire.
//
// The skill's own example, fed to the console's parser. If the skill's format
// drifts, this test goes red — the orchestrator is never told to write a format
// the console cannot read.
func TestTheSkillExampleIsTheFormatTheConsoleReads(t *testing.T) {
	raw, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read the skill: %v", err)
	}
	examples := skillExamples(string(raw))
	if len(examples) == 0 {
		t.Fatalf("the skill shows no roster's line: %s", skillPath)
	}

	orchestrator := 0
	for _, example := range examples {
		job, tail, ok := parseRoster(example)
		if !ok {
			t.Errorf("the console cannot read the skill's example %q", example)
			continue
		}
		if !knownCrewJob(job) {
			t.Errorf("the skill's example names the job %q, which is not one of the crew's: %q", job, example)
			continue
		}
		if job != "orchestrator" {
			continue
		}
		orchestrator++
		progress, _ := parseProgress(tail)
		if !progress.Round.Reported {
			t.Errorf("the console reads no round from the skill's example %q", example)
		}
		if !progress.Interactions.Reported {
			t.Errorf("the console reads no interactions from the skill's example %q", example)
		}
	}
	if orchestrator == 0 {
		t.Errorf("the skill shows no orchestrator's line: %s", skillPath)
	}
}

// knownCrewJob reports whether the job is one of the crew's.
func knownCrewJob(job string) bool {
	for _, name := range crewJobs {
		if name == job {
			return true
		}
	}
	return false
}
