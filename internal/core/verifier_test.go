package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDeriveMatchCountsGoTests(t *testing.T) {
	cases := []struct {
		name    string
		command string
		output  string
		want    int
		unit    string
		known   bool
	}{
		{
			name:    "verbose pass",
			command: "go test ./internal/core -run Board -v",
			output:  "=== RUN   TestBoard\n--- PASS: TestBoard (0.00s)\nPASS\nok  \texample.com/m/pkg\t0.003s\n",
			want:    1, unit: "test", known: true,
		},
		{
			name:    "verbose with several results",
			command: "go test ./... -v",
			output:  "--- PASS: TestA (0.00s)\n--- PASS: TestB (0.00s)\n--- FAIL: TestC (0.00s)\n",
			want:    3, unit: "test", known: true,
		},
		{
			name:    "a filter that matched nothing",
			command: "go test ./internal/core -run NoSuchTest",
			output:  "ok  \texample.com/m/pkg\t0.001s [no tests to run]\n",
			want:    0, unit: "test", known: true,
		},
		{
			name:    "quiet run reports packages",
			command: "go test ./internal/core",
			output:  "ok  \texample.com/m/pkg\t0.003s\n",
			want:    1, unit: "package", known: true,
		},
		{
			name:    "a command we cannot interpret",
			command: "make check",
			output:  "everything is fine\n",
			want:    0, unit: "", known: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, unit, known := deriveMatch(c.command, c.output)
			if got != c.want || unit != c.unit || known != c.known {
				t.Fatalf("deriveMatch = (%d, %q, %v), want (%d, %q, %v)",
					got, unit, known, c.want, c.unit, c.known)
			}
		})
	}
}

func TestIsGoTestIgnoresEnvPrefixes(t *testing.T) {
	if !isGoTest("go test ./...") {
		t.Error("go test should be recognised")
	}
	if !isGoTest("GOFLAGS=-count=1 go test ./...") {
		t.Error("an env-prefixed go test should be recognised")
	}
	if isGoTest("go build ./...") {
		t.Error("go build is not go test")
	}
	if isGoTest("echo hi") {
		t.Error("echo is not go test")
	}
}

func TestVerifierRunsTheCommand(t *testing.T) {
	v := &Verifier{Dir: t.TempDir(), Timeout: 30 * time.Second}

	res, err := v.Run(context.Background(), "echo hello", "line-cook")
	if err != nil {
		t.Fatal(err)
	}
	if res.Proof.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0", res.Proof.ExitCode)
	}
	if res.Proof.Match != MatchUnknown {
		t.Fatalf("match = %q, want unknown — a plain command tells us nothing", res.Proof.Match)
	}
	if res.Proof.OutputDigest == "" {
		t.Fatal("no output digest recorded")
	}
	if !strings.Contains(res.Output, "hello") {
		t.Fatalf("output not captured: %q", res.Output)
	}
}

func TestVerifierRecordsAFailure(t *testing.T) {
	v := &Verifier{Dir: t.TempDir(), Timeout: 30 * time.Second}

	res, err := v.Run(context.Background(), "false", "line-cook")
	if err != nil {
		t.Fatal(err)
	}
	if res.Proof.ExitCode == 0 {
		t.Fatal("exit code should be non-zero")
	}
	if res.Proof.Match != MatchUnknown {
		t.Fatalf("match = %q, want unknown", res.Proof.Match)
	}
}

func TestVerifierTimeoutIsDistinguishable(t *testing.T) {
	v := &Verifier{Dir: t.TempDir(), Timeout: 100 * time.Millisecond}

	res, err := v.Run(context.Background(), "sleep 5", "line-cook")
	if err != nil {
		t.Fatal(err)
	}
	if !res.Proof.TimedOut {
		t.Fatal("the proof should record that it timed out")
	}
	// A timeout must not look like an ordinary failure.
	if res.Proof.Match != MatchUnknown {
		t.Fatalf("match = %q, want unknown", res.Proof.Match)
	}
}

func TestVerifierRejectsAnEmptyCommand(t *testing.T) {
	v := &Verifier{Dir: t.TempDir()}
	if _, err := v.Run(context.Background(), "   ", "line-cook"); err == nil {
		t.Fatal("expected an empty command to be refused")
	}
}

// TestATimedOutVerificationLeavesNoChild covers the process group.
//
// A verification is a shell, and the shell has children: the test suite it starts.
// Killing the shell alone leaves the suite running, holding the project directory.
func TestATimedOutVerificationLeavesNoChild(t *testing.T) {
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "child.pid")
	// The shell starts a child, writes its pid, and waits.
	command := fmt.Sprintf("sleep 10 & echo $! > %s; wait", pidFile)

	v := &Verifier{Dir: dir, Timeout: 300 * time.Millisecond}
	start := time.Now()
	res, err := v.Run(context.Background(), command, "line-cook")
	elapsed := time.Since(start)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Proof.TimedOut {
		t.Fatalf("the proof should record that it timed out: %+v", res.Proof)
	}
	// The run must end when the timeout fires. If it waited for the child, the
	// child would be gone by accident — the shell's children would outlive it.
	if elapsed > 5*time.Second {
		t.Fatalf("the run took %s — it waited for the child instead of killing it", elapsed)
	}

	pid := readPID(t, pidFile)
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if !processAlive(pid) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("the child process %d is still running after the timeout", pid)
}

// readPID waits for the command to write its child's pid and reads it.
func readPID(t *testing.T, path string) int {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		raw, err := os.ReadFile(path)
		if err == nil {
			pid, convErr := strconv.Atoi(strings.TrimSpace(string(raw)))
			if convErr == nil {
				return pid
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("the command never wrote its child's pid to %s", path)
	return 0
}
