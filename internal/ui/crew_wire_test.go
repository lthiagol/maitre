package ui

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// The wire between the **format** and the **console**.
//
// The format lived in **three** places: the **skill's table**, the **console's
// parser**, and the **startup task** — the copy the console hands the
// orchestrator at spawn. Each had its own witness, and the task's witness held
// the *words* and not the *shape*: the separator could drift with every word
// intact and nothing went red.
//
// The wire joins them to the parser now. **The skill's table is the
// specification** — it is the role's own knowledge, and it may be edited. The
// console's copy (`rosterExample`) and the parser **follow** it, and the
// startup task's copy is **an example** of it.
//
// A drift in any copy turns a test red.

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

// startupTaskExample is a roster's line the startup task hands the orchestrator.
var startupTaskExample = regexp.MustCompile("herdr tab rename [^\"]*\"([^\"]+)\"")

// crewJobs are the jobs the roster's lines name.
//
// The orchestrator's is the one the console reads the progress from, and the
// other two are the crew's other members.
var crewJobs = []string{"orchestrator", "runner", "reviewer"}

// fillPlaceholders puts the numbers in the skill's example.
//
// The console reads numbers, and a placeholder is not one.
func fillPlaceholders(s string) string {
	s = strings.ReplaceAll(s, "<N>", "2")
	return strings.ReplaceAll(s, "<M>", "6")
}

// skillExamples are the roster's lines the skill teaches.
func skillExamples(skill string) []string {
	return skillExample.FindAllString(skill, -1)
}

// startupTaskExamples are the roster's lines the startup task hands over.
func startupTaskExamples(task string) []string {
	var out []string
	for _, match := range startupTaskExample.FindAllStringSubmatch(task, -1) {
		out = append(out, match[1])
	}
	return out
}

// TestTheSkillExampleIsTheFormatTheConsoleReads is the wire, for the skill.
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
		job, ok := checkRosterExample(t, example)
		if !ok {
			continue
		}
		if job != "orchestrator" {
			continue
		}
		orchestrator++
	}
	if orchestrator == 0 {
		t.Errorf("the skill shows no orchestrator's line: %s", skillPath)
	}
}

// TestTheStartupTaskExampleIsTheFormatTheConsoleReads is the wire, for the
// startup task.
//
// The startup task is the copy the console **hands the orchestrator at spawn**,
// in the same breath as the work. Its own witness held the *words* and not the
// *shape*, so the separator could drift with every word intact.
//
// This wire feeds the startup task's **own example** to the console's parser.
// A drift turns this test red.
func TestTheStartupTaskExampleIsTheFormatTheConsoleReads(t *testing.T) {
	m := sousModel(t)
	task := m.startupTask()
	examples := startupTaskExamples(task)
	if len(examples) == 0 {
		t.Fatalf("the startup task shows no roster's line: %q", task)
	}

	for _, example := range examples {
		job, ok := checkRosterExample(t, example)
		if !ok {
			continue
		}
		if job != "orchestrator" {
			t.Errorf("the startup task's example names the job %q, want orchestrator: %q", job, example)
		}
	}
}

// TestTheStartupTaskCopiesTheSkillIsTheSpecification joins the two copies.
//
// The skill's table is the **specification** and the startup task's copy is an
// **example** of it. So the task's copy must be one of the skill's own lines —
// not a variant of it, and not a third format.
func TestTheStartupTaskCopiesTheSkillIsTheSpecification(t *testing.T) {
	raw, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read the skill: %v", err)
	}
	known := map[string]bool{}
	for _, example := range skillExamples(string(raw)) {
		known[example] = true
	}

	m := sousModel(t)
	for _, example := range startupTaskExamples(m.startupTask()) {
		if !known[example] {
			t.Errorf("the startup task's copy %q is not the skill's own line — the skill's table is the specification", example)
		}
	}
}

// checkRosterExample holds one roster's line against the console's parser.
//
// The line must parse as a roster's line, name one of the crew's jobs, and —
// for the orchestrator — carry a round and an interaction count the console
// reads.
func checkRosterExample(t *testing.T, example string) (job string, ok bool) {
	t.Helper()
	// The console reads numbers, and a placeholder is not one: the example is
	// filled in before it is read.
	filled := fillPlaceholders(example)
	job, tail, parsed := parseRoster(filled)
	if !parsed {
		t.Errorf("the console cannot read the roster's line %q", example)
		return "", false
	}
	if !knownCrewJob(job) {
		t.Errorf("the roster's line names the job %q, which is not one of the crew's: %q", job, example)
		return job, false
	}
	if job != "orchestrator" {
		return job, true
	}
	progress, _ := parseProgress(tail)
	if !progress.Round.Reported {
		t.Errorf("the console reads no round from the roster's line %q", example)
	}
	if !progress.Interactions.Reported {
		t.Errorf("the console reads no interactions from the roster's line %q", example)
	}
	return job, true
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
