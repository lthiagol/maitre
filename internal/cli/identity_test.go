package cli

import "testing"

// signedOffRecipe builds a recipe that is approved and waiting to be claimed.
func signedOffRecipe(t *testing.T, dir string) {
	t.Helper()
	for _, s := range []struct {
		name string
		args []string
	}{
		{"create", []string{"recipe", "create", "--title", "A recipe",
			"--outcome", "o", "--problem", "p",
			"--out-of-scope", "a", "--out-of-scope", "b"}},
		{"criterion", []string{"ac", "add", "R01",
			"--description", "d", "--verification", "go test ./..."}},
		{"work package", []string{"wp", "add", "R01", "--name", "core"}},
		{"step", []string{"step", "add", "R01",
			"--action", "do it", "--done-when", "done", "--covers", "AC-01"}},
		{"groom", []string{"recipe", "advance", "R01", "groom"}},
		{"specify", []string{"recipe", "advance", "R01", "specify"}},
		{"approve", []string{"--actor", "chef", "recipe", "approve", "R01", "--by", "chef"}},
	} {
		res := cli(t, dir, s.args...)
		if !res.ok(t) {
			t.Fatalf("%s failed:\n%s", s.name, res.Stderr)
		}
	}
}

// TestActorFlagsAreGlobal covers the identity flags in the documented position.
//
// The three identity flags were declared only by `recipe create`. Every other verb
// rejected them as unknown flags, and the documented form — `scale --actor x --session
// y --role z recipe claim R01` — was refused by the top-level dispatcher as an
// unknown command. The L5 session stamp could then only be set through environment
// variables.
func TestActorFlagsAreGlobal(t *testing.T) {
	dir := newPlan(t)
	t.Setenv("SCALE_ACTOR", "")
	t.Setenv("SCALE_SESSION", "")
	t.Setenv("SCALE_ROLE", "")
	signedOffRecipe(t, dir)

	res := cli(t, dir, "--actor", "line-cook", "--session", "s-exec", "--role", "runner",
		"recipe", "claim", "R01", "--agent", "line-cook")
	if !res.ok(t) {
		t.Fatalf("the documented form was refused:\n%s", res.Stderr)
	}

	// The event carries the stamp.
	log := cli(t, dir, "log", "--limit", "1")
	if !log.ok(t) {
		t.Fatalf("log failed:\n%s", log.Stderr)
	}
	events := log.doc(t)["events"].([]any)
	if len(events) == 0 {
		t.Fatal("the journal is empty")
	}
	e := events[0].(map[string]any)
	if e["actor"] != "line-cook" {
		t.Fatalf("actor = %v, want line-cook", e["actor"])
	}
	if e["session"] != "s-exec" {
		t.Fatalf("session = %v, want s-exec", e["session"])
	}
	if e["role"] != "runner" {
		t.Fatalf("role = %v, want runner", e["role"])
	}
}

// TestActorFlagsWorkAfterTheVerb covers the other position.
//
// The flags are the same flags: a caller may put them before the command or after the
// verb.
func TestActorFlagsWorkAfterTheVerb(t *testing.T) {
	dir := newPlan(t)
	t.Setenv("SCALE_ACTOR", "")
	t.Setenv("SCALE_SESSION", "")
	t.Setenv("SCALE_ROLE", "")
	signedOffRecipe(t, dir)

	res := cli(t, dir, "recipe", "claim", "R01", "--agent", "line-cook",
		"--actor", "line-cook", "--session", "s-exec", "--role", "runner")
	if !res.ok(t) {
		t.Fatalf("the flags were refused:\n%s", res.Stderr)
	}

	log := cli(t, dir, "log", "--limit", "1")
	if !log.ok(t) {
		t.Fatalf("log failed:\n%s", log.Stderr)
	}
	e := log.doc(t)["events"].([]any)[0].(map[string]any)
	if e["actor"] != "line-cook" || e["session"] != "s-exec" || e["role"] != "runner" {
		t.Fatalf("the stamp did not reach the journal: %v", e)
	}
}
