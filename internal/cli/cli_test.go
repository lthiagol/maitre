package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"
)

// run executes one command in a directory and returns the parsed output.
type result struct {
	Code   int
	Stdout string
	Stderr string
}

// doc parses the stdout document.
func (r result) doc(t *testing.T) map[string]any {
	t.Helper()
	var d map[string]any
	if err := json.Unmarshal([]byte(r.Stdout), &d); err != nil {
		t.Fatalf("stdout is not a JSON document: %v\n%s", err, r.Stdout)
	}
	return d
}

// ok reports whether the command succeeded.
func (r result) ok(t *testing.T) bool {
	t.Helper()
	if r.Code == 0 {
		return true
	}
	// A failure prints to stderr; the document is there instead.
	var d map[string]any
	if err := json.Unmarshal([]byte(r.Stderr), &d); err == nil {
		return false
	}
	t.Fatalf("command failed with %d and no document:\n%s\n%s", r.Code, r.Stdout, r.Stderr)
	return false
}

// cli runs a command in a working directory.
func cli(t *testing.T, dir string, args ...string) result {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	var out, errb bytes.Buffer
	runner := &Runner{Args: args, Stdout: &out, Stderr: &errb}
	code := runner.Run()
	return result{Code: code, Stdout: out.String(), Stderr: errb.String()}
}

// newPlan makes a directory with an initialised plan.
func newPlan(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	res := cli(t, dir, "init")
	if res.Code != 0 {
		t.Fatalf("init failed:\n%s", res.Stderr)
	}
	return dir
}

// TestTheAgentLoop walks a recipe from creation to closed, the way an agent would.
//
// This is the round trip the whole project exists for: the agent authors, the gates
// check, the chef signs off, the agent executes and proves.
func TestTheAgentLoop(t *testing.T) {
	dir := newPlan(t)

	steps := []struct {
		name string
		args []string
	}{
		{"create", []string{"recipe", "create",
			"--title", "Ship the board",
			"--outcome", "A human can see the plan",
			"--problem", "v1 had 47 commands for one workflow",
			"--out-of-scope", "web shell",
			"--out-of-scope", "orchestration"}},
		{"add criterion", []string{"ac", "add", "RCP01",
			"--description", "the board renders 100 recipes",
			"--verification", "go test ./internal/core -run Board"}},
		{"add work package", []string{"wp", "add", "RCP01", "--name", "core"}},
		{"add step", []string{"step", "add", "RCP01",
			"--action", "render the board", "--done-when", "it renders", "--covers", "AC-01"}},
		{"groom", []string{"recipe", "advance", "RCP01", "groom"}},
		{"specify", []string{"recipe", "advance", "RCP01", "specify"}},
		{"approve", []string{"--actor", "chef", "recipe", "approve", "RCP01", "--by", "chef"}},
		{"claim", []string{"recipe", "claim", "RCP01", "--agent", "line-cook"}},
		{"step done", []string{"step", "done", "RCP01", "G1"}},
		{"document", []string{"recipe", "document", "RCP01", "--text", "One board, five tabs."}},
	}
	for _, s := range steps {
		res := cli(t, dir, s.args...)
		if !res.ok(t) {
			t.Fatalf("%s failed:\n%s", s.name, res.Stderr)
		}
	}

	// The proof is what the gates want, so sign one off.
	res := cli(t, dir, "--actor", "chef", "ac", "signoff", "RCP01", "AC-01", "--by", "chef", "--note", "checked by hand")
	if !res.ok(t) {
		t.Fatalf("signoff failed:\n%s", res.Stderr)
	}

	for _, s := range []struct {
		name string
		args []string
	}{
		{"finish", []string{"recipe", "advance", "RCP01", "finish"}},
		{"review", []string{"recipe", "review", "RCP01", "--reviewer", "taster", "--verdict", "ok"}},
		{"pass review", []string{"recipe", "advance", "RCP01", "pass_review"}},
		{"close", []string{"recipe", "advance", "RCP01", "close"}},
	} {
		res := cli(t, dir, s.args...)
		if !res.ok(t) {
			t.Fatalf("%s failed:\n%s", s.name, res.Stderr)
		}
	}

	// It is closed.
	show := cli(t, dir, "recipe", "show", "RCP01")
	if !show.ok(t) {
		t.Fatalf("show failed:\n%s", show.Stderr)
	}
	recipe := show.doc(t)["recipe"].(map[string]any)
	if recipe["lifecycle"] != "closed" {
		t.Fatalf("lifecycle = %v, want closed", recipe["lifecycle"])
	}
}

// TestTheGatesRefuseWhatIsNotReady checks that the CLI reports a refusal as a gate
// failure, not as a crash.
func TestTheGatesRefuseWhatIsNotReady(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o")

	// No criteria, no scope, no steps — grooming must be refused.
	res := cli(t, dir, "recipe", "advance", "RCP01", "groom")
	if res.Code != ExitGate {
		t.Fatalf("exit = %d, want %d", res.Code, ExitGate)
	}
	var d map[string]any
	if err := json.Unmarshal([]byte(res.Stderr), &d); err != nil {
		t.Fatalf("stderr is not a document: %v\n%s", err, res.Stderr)
	}
	if d["ok"] != false {
		t.Fatalf("the refusal should say ok:false: %v", d)
	}
	gates, ok := d["gates"].([]any)
	if !ok || len(gates) == 0 {
		t.Fatalf("the refusal should name the gates that failed: %v", d)
	}
}

// TestAnUnknownVerbIsAClearError checks the surface is navigable.
func TestAnUnknownVerbIsAClearError(t *testing.T) {
	dir := newPlan(t)
	res := cli(t, dir, "recipe", "frobnicate")
	if res.Code != ExitError {
		t.Fatalf("exit = %d, want %d", res.Code, ExitError)
	}
	if !strings.Contains(res.Stderr, "frobnicate") {
		t.Fatalf("the error should name the verb: %s", res.Stderr)
	}
	if !strings.Contains(res.Stderr, "create") {
		t.Fatalf("the error should list the verbs: %s", res.Stderr)
	}
}

// TestTheJournalIsStampedWithTheSession checks the identity reaches every write.
func TestTheJournalIsStampedWithTheSession(t *testing.T) {
	dir := newPlan(t)
	t.Setenv(envActor, "line-cook")
	t.Setenv(envSession, "s-42")
	t.Setenv(envRole, "runner")

	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o")

	res := cli(t, dir, "log")
	if !res.ok(t) {
		t.Fatalf("log failed:\n%s", res.Stderr)
	}
	events := res.doc(t)["events"].([]any)
	if len(events) == 0 {
		t.Fatal("the journal is empty")
	}
	for _, e := range events {
		ev := e.(map[string]any)
		if ev["actor"] != "line-cook" {
			t.Fatalf("actor = %v, want line-cook", ev["actor"])
		}
		if ev["session"] != "s-42" {
			t.Fatalf("session = %v, want s-42", ev["session"])
		}
	}
}

// TestTheBacklogPromotesUpTheLadder checks the one-way promotion, end to end.
func TestTheBacklogPromotesUpTheLadder(t *testing.T) {
	dir := newPlan(t)

	res := cli(t, dir, "special", "create", "--title", "A hunch", "--body", "worth noticing")
	if !res.ok(t) {
		t.Fatalf("special create failed:\n%s", res.Stderr)
	}
	res = cli(t, dir, "special", "promote", "S01")
	if !res.ok(t) {
		t.Fatalf("special promote failed:\n%s", res.Stderr)
	}
	note := res.doc(t)["prepNote"].(map[string]any)
	if note["id"] != "N01" {
		t.Fatalf("prep note id = %v, want N01", note["id"])
	}

	res = cli(t, dir, "prepnote", "promote", "N01")
	if !res.ok(t) {
		t.Fatalf("prepnote promote failed:\n%s", res.Stderr)
	}
	recipe := res.doc(t)["recipe"].(map[string]any)
	if recipe["id"] != "RCP01" {
		t.Fatalf("recipe id = %v, want RCP01", recipe["id"])
	}
}

// TestThePlanIsOneDocument checks the cheap read an agent starts from.
func TestThePlanIsOneDocument(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o")

	res := cli(t, dir, "plan")
	if !res.ok(t) {
		t.Fatalf("plan failed:\n%s", res.Stderr)
	}
	doc := res.doc(t)
	recipes := doc["recipes"].([]any)
	if len(recipes) != 1 {
		t.Fatalf("got %d recipes, want 1", len(recipes))
	}
	if doc["project"] == nil {
		t.Fatal("the plan should carry the project")
	}
}

// TestNoCommandIsHelpfullyRefused checks the entry point.
func TestNoCommandIsHelpfullyRefused(t *testing.T) {
	res := cli(t, t.TempDir())
	if res.Code != ExitError {
		t.Fatalf("exit = %d, want %d", res.Code, ExitError)
	}
	if !strings.Contains(res.Stdout, "maitre-cli") {
		t.Fatalf("the usage should mention the binary: %s", res.Stdout)
	}
}

// TestTheCLIUsesOnlyTheOperations checks the CLI never reaches around the domain.
//
// Every write goes through an operation, so the console and the CLI cannot disagree
// about what is allowed.
func TestTheCLIUsesOnlyTheOperations(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "A recipe",
		"--outcome", "o", "--problem", "p",
		"--out-of-scope", "a", "--out-of-scope", "b")

	// The frozen sections are closed after sign-off, and the CLI says so.
	cli(t, dir, "ac", "add", "RCP01", "--description", "d", "--verification", "go test ./...")
	cli(t, dir, "wp", "add", "RCP01", "--name", "core")
	cli(t, dir, "step", "add", "RCP01", "--action", "do it", "--done-when", "done", "--covers", "AC-01")
	cli(t, dir, "recipe", "advance", "RCP01", "groom")
	cli(t, dir, "recipe", "advance", "RCP01", "specify")
	cli(t, dir, "--actor", "chef", "recipe", "approve", "RCP01", "--by", "chef")

	res := cli(t, dir, "ac", "add", "RCP01", "--description", "another", "--verification", "go test ./...")
	if res.Code == 0 {
		t.Fatal("adding a criterion after sign-off should be refused")
	}
	if !strings.Contains(res.Stderr, "frozen") {
		t.Fatalf("the refusal should say why: %s", res.Stderr)
	}
}

// TestTheReviewLoopGoesBackwards covers the reviewer's move when a fix does not hold.
func TestTheReviewLoopGoesBackwards(t *testing.T) {
	dir := newPlan(t)

	cli(t, dir, "recipe", "create", "--title", "A recipe",
		"--outcome", "o", "--problem", "p",
		"--out-of-scope", "a", "--out-of-scope", "b")
	cli(t, dir, "ac", "add", "RCP01", "--description", "d", "--verification", "go test ./...")
	cli(t, dir, "wp", "add", "RCP01", "--name", "core")
	cli(t, dir, "step", "add", "RCP01", "--action", "do it", "--done-when", "done", "--covers", "AC-01")
	cli(t, dir, "recipe", "advance", "RCP01", "groom")
	cli(t, dir, "recipe", "advance", "RCP01", "specify")
	cli(t, dir, "--actor", "chef", "recipe", "approve", "RCP01", "--by", "chef")
	cli(t, dir, "recipe", "claim", "RCP01", "--agent", "line-cook")
	cli(t, dir, "step", "done", "RCP01", "G1")
	cli(t, dir, "--actor", "chef", "ac", "signoff", "RCP01", "AC-01", "--by", "chef", "--note", "checked")
	cli(t, dir, "recipe", "advance", "RCP01", "finish")

	// The reviewer files a finding and says no.
	res := cli(t, dir, "finding", "add", "RCP01", "--description", "the queue drops findings",
		"--severity", "high", "--category", "correctness")
	if !res.ok(t) {
		t.Fatalf("finding add failed:\n%s", res.Stderr)
	}
	res = cli(t, dir, "recipe", "review", "RCP01", "--reviewer", "taster", "--verdict", "changes-needed")
	if !res.ok(t) {
		t.Fatalf("review failed:\n%s", res.Stderr)
	}

	// The reviewer says no — the recipe must not advance.
	res = cli(t, dir, "recipe", "advance", "RCP01", "pass_review")
	if res.Code != ExitGate {
		t.Fatalf("exit = %d, want %d — a changes-needed verdict must not advance", res.Code, ExitGate)
	}

	// The fixer resolves it, and the reviewer sends it back.
	cli(t, dir, "finding", "resolve", "RCP01", "F-01", "--note", "read the queue once")
	res = cli(t, dir, "finding", "reopen", "RCP01", "F-01", "--reason", "the empty case is still missed")
	if !res.ok(t) {
		t.Fatalf("reopen failed:\n%s", res.Stderr)
	}

	show := cli(t, dir, "recipe", "show", "RCP01")
	if !show.ok(t) {
		t.Fatalf("show failed:\n%s", show.Stderr)
	}
	recipe := show.doc(t)["recipe"].(map[string]any)
	findings := recipe["review"].(map[string]any)["findings"].([]any)
	f := findings[0].(map[string]any)
	if f["status"] != "open" {
		t.Fatalf("status = %v, want open", f["status"])
	}
	comments := recipe["review"].(map[string]any)["comments"].([]any)
	if len(comments) == 0 {
		t.Fatal("the objection should be recorded as a comment")
	}
}

// TestTheJournalCanBeReadFromASequence covers how an agent notices what changed.
func TestTheJournalCanBeReadFromASequence(t *testing.T) {
	dir := newPlan(t)
	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o")

	all := cli(t, dir, "log")
	if !all.ok(t) {
		t.Fatalf("log failed:\n%s", all.Stderr)
	}
	events := all.doc(t)["events"].([]any)
	if len(events) == 0 {
		t.Fatal("the journal is empty")
	}

	// The sequence number of the newest event.
	var seq float64
	for _, e := range events {
		ev := e.(map[string]any)
		if s, ok := ev["seq"].(float64); ok && s > seq {
			seq = s
		}
	}

	cli(t, dir, "recipe", "create", "--title", "Another", "--outcome", "o")

	since := cli(t, dir, "log", "--since", itoa(int(seq)))
	if !since.ok(t) {
		t.Fatalf("log --since failed:\n%s", since.Stderr)
	}
	got := since.doc(t)["events"].([]any)
	if len(got) != 1 {
		t.Fatalf("got %d events since %d, want 1", len(got), int(seq))
	}
}

// itoa is a small helper for the test.
func itoa(n int) string { return strconv.Itoa(n) }
