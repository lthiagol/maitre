package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestAnUnknownLaneIsAnError covers the typo.
//
// v1's `--type` answered "there are none" for a value that did not exist, and a
// client could not tell a typo from an empty plan. P7 was written to prevent that.
func TestAnUnknownLaneIsAnError(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o")

	res := cli(t, dir, "recipe", "list", "--lane", "groomng")
	if res.Code == 0 {
		t.Fatalf("an unknown lane should be an error, not an empty list:\n%s", res.Stdout)
	}
	for _, lane := range []string{"grooming", "approval", "execution", "review", "blocked", "done"} {
		if !strings.Contains(res.Stderr, lane) {
			t.Fatalf("the error should name the lane %q: %s", lane, res.Stderr)
		}
	}

	// A known lane still works.
	res = cli(t, dir, "recipe", "list", "--lane", "grooming")
	if !res.ok(t) {
		t.Fatalf("a known lane should work:\n%s", res.Stderr)
	}
}

// TestAdvanceRefusesTheEventsThatHaveTheirOwnVerb covers one intent, one verb.
//
// `advance approve` signed off with the actor as the approver — the alias is the
// weaker of the two paths to a sign-off.
func TestAdvanceRefusesTheEventsThatHaveTheirOwnVerb(t *testing.T) {
	dir := newPlan(t)
	signedOffRecipe(t, dir)

	for _, ev := range []string{"approve", "request_changes", "start"} {
		res := cli(t, dir, "recipe", "advance", "RCP01", ev)
		if res.Code == 0 {
			t.Fatalf("advance %s should be refused", ev)
		}
		if !strings.Contains(res.Stderr, "own verb") {
			t.Fatalf("the refusal should say why: %s", res.Stderr)
		}
	}

	// The verb that owns the event still works — the recipe is signed off by
	// `signedOffRecipe`, and the claim is the only road to the line.
	res := cli(t, dir, "recipe", "claim", "RCP01", "--agent", "line-cook")
	if !res.ok(t) {
		t.Fatalf("claim failed:\n%s", res.Stderr)
	}
}

// TestEveryKindAcceptsItsShortID covers the id tolerance.
//
// "R1" and "r1" worked for a recipe; the other kinds answered only to their exact
// form, so a client had to know each one's padding.
func TestEveryKindAcceptsItsShortID(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o")
	cli(t, dir, "special", "create", "--title", "A hunch", "--body", "worth noticing")
	cli(t, dir, "prepnote", "create", "--title", "A note", "--intent", "first")

	for _, id := range []string{"R1", "r1", "RCP01"} {
		res := cli(t, dir, "recipe", "show", id)
		if !res.ok(t) {
			t.Fatalf("recipe %s: %s", id, res.Stderr)
		}
	}
	for _, id := range []string{"S1", "s1", "S01"} {
		res := cli(t, dir, "special", "update", id, "--body", "sharper")
		if !res.ok(t) {
			t.Fatalf("special %s: %s", id, res.Stderr)
		}
	}
	for _, id := range []string{"N1", "n1", "N01"} {
		res := cli(t, dir, "prepnote", "update", id, "--intent", "sharper")
		if !res.ok(t) {
			t.Fatalf("prep note %s: %s", id, res.Stderr)
		}
	}
}

// TestEveryDocumentedStepStatusIsReachable covers the step verbs.
//
// Two of the five documented statuses had no verb, so a value in the model had no
// path to the CLI.
func TestEveryDocumentedStepStatusIsReachable(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o")
	res := cli(t, dir, "step", "add", "RCP01", "--action", "do it", "--done-when", "done")
	if !res.ok(t) {
		t.Fatalf("step add failed:\n%s", res.Stderr)
	}

	// A new step starts pending.
	if got := stepStatus(t, dir); got != "pending" {
		t.Fatalf("status = %q, want pending", got)
	}
	for _, s := range []struct{ verb, want string }{
		{"start", "in-progress"},
		{"done", "done"},
		{"skip", "skipped"},
		{"fail", "failed"},
	} {
		res := cli(t, dir, "step", s.verb, "RCP01", "G1")
		if !res.ok(t) {
			t.Fatalf("step %s failed:\n%s", s.verb, res.Stderr)
		}
		if got := stepStatus(t, dir); got != s.want {
			t.Fatalf("status = %q, want %q", got, s.want)
		}
	}
}

// stepStatus is the status of the recipe's first step.
func stepStatus(t *testing.T, dir string) string {
	t.Helper()
	res := cli(t, dir, "recipe", "show", "RCP01")
	if !res.ok(t) {
		t.Fatalf("show failed:\n%s", res.Stderr)
	}
	recipe := res.doc(t)["recipe"].(map[string]any)
	steps := recipe["guidance"].([]any)
	if len(steps) == 0 {
		t.Fatal("the recipe has no steps")
	}
	return steps[0].(map[string]any)["status"].(string)
}

// TestAFailurePrintsItsDocumentOnStderr covers where a failure goes.
//
// An agent that parses stdout as JSON gets an empty string on every failure. The
// document is on stderr, and the help says so.
func TestAFailurePrintsItsDocumentOnStderr(t *testing.T) {
	dir := newPlan(t)

	res := cli(t, dir, "recipe", "show", "R99")
	if res.Code == 0 {
		t.Fatal("the command should fail")
	}
	if res.Stdout != "" {
		t.Fatalf("stdout should hold nothing on a failure:\n%s", res.Stdout)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(res.Stderr), &doc); err != nil {
		t.Fatalf("the failure should print a document: %v\n%s", err, res.Stderr)
	}
	if doc["ok"] != false {
		t.Fatalf("the document should say ok:false: %v", doc)
	}

	// The help says the truth.
	help := cli(t, dir, "help")
	if !help.ok(t) {
		t.Fatalf("help failed:\n%s", help.Stderr)
	}
	if !strings.Contains(help.Stdout, "stderr") {
		t.Fatalf("the help should say where a failure goes:\n%s", help.Stdout)
	}
}
