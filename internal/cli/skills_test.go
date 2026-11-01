package cli

import (
	"strings"
	"testing"
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
