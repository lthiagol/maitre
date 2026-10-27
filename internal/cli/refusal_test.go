package cli

import (
	"strings"
	"testing"
)

// TestARefusalIsInTheLog covers the refusal the log used to swallow.
//
// The journal is the only event log, and `scale log` is how an agent reads it. A
// refusal that is not there is an attempt nobody can see was turned away.
func TestARefusalIsInTheLog(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o", "--problem", "p")

	// Groom is refused: no criteria and no out-of-scope.
	res := cli(t, dir, "recipe", "advance", "R01", "groom")
	if res.Code != ExitGate {
		t.Fatalf("exit = %d, want %d", res.Code, ExitGate)
	}

	log := cli(t, dir, "log")
	if !log.ok(t) {
		t.Fatalf("log failed:\n%s", log.Stderr)
	}
	if !strings.Contains(log.Stdout, `"kind": "refusal"`) {
		t.Fatalf("the log does not carry the refusal:\n%s", log.Stdout)
	}
	if !strings.Contains(log.Stdout, `"event": "groom"`) {
		t.Fatalf("the refusal does not name the event:\n%s", log.Stdout)
	}
	if !strings.Contains(log.Stdout, "scope.out_of_scope_declared") {
		t.Fatalf("the refusal does not name the gates:\n%s", log.Stdout)
	}
}

// TestASuccessfulWriteIsNotARefusal covers the other side.
//
// The log must distinguish a refusal from a write: a plan whose every entry is a
// refusal is a plan where nothing happened.
func TestASuccessfulWriteIsNotARefusal(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o", "--problem", "p")

	log := cli(t, dir, "log")
	if !log.ok(t) {
		t.Fatalf("log failed:\n%s", log.Stderr)
	}
	if strings.Contains(log.Stdout, `"kind": "refusal"`) {
		t.Fatalf("a write was journalled as a refusal:\n%s", log.Stdout)
	}
	if !strings.Contains(log.Stdout, `"kind": "created"`) {
		t.Fatalf("the write is not in the log:\n%s", log.Stdout)
	}
}
