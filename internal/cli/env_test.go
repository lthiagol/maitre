package cli

import (
	"strings"
	"testing"
)

// TestTheCurrentEnvNamesCarryTheIdentity covers the new names.
//
// The names carry the binary's name: `MAITRE_CLI_ACTOR`, `MAITRE_CLI_SESSION`,
// `MAITRE_CLI_ROLE`. A caller that sets them gets what it set.
func TestTheCurrentEnvNamesCarryTheIdentity(t *testing.T) {
	dir := newPlan(t)
	clearIdentityEnv(t)
	t.Setenv(envActor, "runner")
	t.Setenv(envSession, "s-42")
	t.Setenv(envRole, "runner")

	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o")

	res := cli(t, dir, "log")
	if !res.ok(t) {
		t.Fatalf("log failed:\n%s", res.Stderr)
	}
	e := res.doc(t)["events"].([]any)[0].(map[string]any)
	if e["actor"] != "runner" {
		t.Fatalf("actor = %v, want runner", e["actor"])
	}
	if e["session"] != "s-42" {
		t.Fatalf("session = %v, want s-42", e["session"])
	}
	if e["role"] != "runner" {
		t.Fatalf("role = %v, want runner", e["role"])
	}
}

// TestTheDeprecatedEnvNamesStillWork covers the deprecation window.
//
// The old names carry the name of the binary the rename removed. A caller that has
// not migrated keeps working — and it is told, on stderr, that the name is going
// away.
func TestTheDeprecatedEnvNamesStillWork(t *testing.T) {
	dir := newPlan(t)
	clearIdentityEnv(t)
	t.Setenv(deprecatedEnvActor, "runner")
	t.Setenv(deprecatedEnvSession, "s-42")
	t.Setenv(deprecatedEnvRole, "runner")

	cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o")

	res := cli(t, dir, "log")
	if !res.ok(t) {
		t.Fatalf("log failed:\n%s", res.Stderr)
	}
	e := res.doc(t)["events"].([]any)[0].(map[string]any)
	if e["actor"] != "runner" {
		t.Fatalf("actor = %v, want runner", e["actor"])
	}
	if e["session"] != "s-42" {
		t.Fatalf("session = %v, want s-42", e["session"])
	}
	if e["role"] != "runner" {
		t.Fatalf("role = %v, want runner", e["role"])
	}
}

// TestTheDeprecatedEnvNamesWarn covers the notice.
//
// The chef learns the new name from the tool itself, not from a release note. The
// warning names both names: the one that is going away and the one to use.
func TestTheDeprecatedEnvNamesWarn(t *testing.T) {
	dir := newPlan(t)
	clearIdentityEnv(t)
	t.Setenv(deprecatedEnvActor, "runner")

	res := cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o")
	if !res.ok(t) {
		t.Fatalf("the command failed:\n%s", res.Stderr)
	}
	if !strings.Contains(res.Stderr, deprecatedEnvActor) {
		t.Errorf("the warning should name the deprecated variable:\n%s", res.Stderr)
	}
	if !strings.Contains(res.Stderr, envActor) {
		t.Errorf("the warning should name the current variable:\n%s", res.Stderr)
	}
}

// TestTheCurrentEnvNamesWinOverTheDeprecatedOnes covers the precedence.
//
// A caller that has migrated must not be told it is using a deprecated name: it sets
// both during the migration, and the new one is the one it means.
func TestTheCurrentEnvNamesWinOverTheDeprecatedOnes(t *testing.T) {
	dir := newPlan(t)
	clearIdentityEnv(t)
	t.Setenv(envActor, "runner")
	t.Setenv(deprecatedEnvActor, "chef")

	res := cli(t, dir, "recipe", "create", "--title", "A recipe", "--outcome", "o")
	if !res.ok(t) {
		t.Fatalf("the command failed:\n%s", res.Stderr)
	}
	if strings.Contains(res.Stderr, deprecatedEnvActor) {
		t.Errorf("a caller that set the current variable should not be warned:\n%s", res.Stderr)
	}

	log := cli(t, dir, "log")
	if !log.ok(t) {
		t.Fatalf("log failed:\n%s", log.Stderr)
	}
	e := log.doc(t)["events"].([]any)[0].(map[string]any)
	if e["actor"] != "runner" {
		t.Fatalf("actor = %v, want runner — the current name wins", e["actor"])
	}
}

// TestTheHelpNamesTheCurrentEnvVariables covers the help.
//
// The help is where a caller learns the names. It must name the current ones — the
// deprecated ones are a fallback, not the way in.
func TestTheHelpNamesTheCurrentEnvVariables(t *testing.T) {
	res := cli(t, t.TempDir(), "help")
	if !res.ok(t) {
		t.Fatalf("help failed:\n%s", res.Stderr)
	}
	for _, name := range []string{envActor, envSession, envRole} {
		if !strings.Contains(res.Stdout, name) {
			t.Errorf("the help should name $%s:\n%s", name, res.Stdout)
		}
	}
}

// clearIdentityEnv empties every identity variable, so a test does not inherit the
// caller's environment.
func clearIdentityEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{
		envActor, envSession, envRole,
		deprecatedEnvActor, deprecatedEnvSession, deprecatedEnvRole,
	} {
		t.Setenv(name, "")
	}
}
