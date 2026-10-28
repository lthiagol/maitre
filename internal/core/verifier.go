package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// MatchState says whether the evidence matched anything, and whether the tool could
// tell.
//
//	proven   — the count is known and greater than zero
//	rejected — the count is known and zero: the filter matched nothing
//	unproven — the count could not be determined; not machine evidence
//
// Only proven satisfies a gate. See docs/decisions/0006-verification.md.
type MatchState string

const (
	MatchProven   MatchState = "proven"
	MatchRejected MatchState = "rejected"
	MatchUnknown  MatchState = "unknown"
	// MatchSignedOff is a human vouching for a procedure the tool could not
	// check. It is not machine evidence and it claims no match — the one
	// exception to the zero-match rule, and it is explicit.
	MatchSignedOff MatchState = "signed-off"
)

// Evidence is a verification result a machine produced. Append-only.
type Evidence struct {
	Command      string     `json:"command"`
	ExitCode     int        `json:"exit_code"`
	TimedOut     bool       `json:"timed_out,omitempty"`
	OutputDigest string     `json:"output_digest"`
	Matched      int        `json:"matched"`
	MatchUnit    string     `json:"match_unit,omitempty"`
	Match        MatchState `json:"match"`
	// SignedOff marks a human vouching for a procedure the tool could not check.
	SignedOff bool   `json:"signed_off,omitempty"`
	At        string `json:"at"`
	By        string `json:"by"`
}

// CriterionStatus is a criterion's state, derived from its latest evidence.
type CriterionStatus string

const (
	CriterionPending  CriterionStatus = "pending"
	CriterionPassed   CriterionStatus = "passed"
	CriterionFailed   CriterionStatus = "failed"
	CriterionRejected CriterionStatus = "rejected"
	CriterionUnproven CriterionStatus = "unproven"
	CriterionDeferred CriterionStatus = "deferred"
)

// DefaultVerifyTimeout is the ceiling for one verification. It must accommodate a
// realistic full test suite — v1's 300s default was unreachable for natural AC shapes.
const DefaultVerifyTimeout = 10 * time.Minute

// Verifier runs a verification command and derives what it can from the output.
//
// It never takes a match count on trust: the count is derived from the output the tool
// captured.
type Verifier struct {
	// Dir is the working directory for the command.
	Dir string
	// Timeout bounds one run. Zero means DefaultVerifyTimeout.
	Timeout time.Duration
}

// VerifyResult is a run's outcome plus the captured output.
type VerifyResult struct {
	Evidence Evidence
	Output   string
}

// Run executes the command and returns the evidence.
//
// The command runs through a shell, so a verification may be a pipeline. The sign-off
// gate is the trust boundary: the human approves the recipe — and its commands — before
// anything executes.
//
// The command runs in its own process group, and a timeout kills the whole group. The
// shell is only the parent: the test suite it starts is the child, and killing the
// shell alone leaves the suite running, holding the project directory.
func (v *Verifier) Run(ctx context.Context, command, by string) (VerifyResult, error) {
	if strings.TrimSpace(command) == "" {
		return VerifyResult{}, errors.New("verification command is empty")
	}
	timeout := v.Timeout
	if timeout <= 0 {
		timeout = DefaultVerifyTimeout
	}
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := shellCommand(runCtx, command)
	cmd.Dir = v.Dir
	startGroup(cmd)
	// The context's own kill would stop at the shell. Kill the group instead.
	cmd.Cancel = func() error {
		killGroup(cmd)
		return nil
	}
	// If a child escapes the group and holds the pipes open, do not wait forever.
	cmd.WaitDelay = 2 * time.Second

	output, runErr := cmd.CombinedOutput()
	sum := sha256.Sum256(output)
	digest := hex.EncodeToString(sum[:])[:16]

	evidence := Evidence{
		Command:      command,
		OutputDigest: digest,
		At:           Now(),
		By:           by,
	}

	// A timeout is its own outcome, not a failure with a confusing message.
	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		evidence.TimedOut = true
		evidence.ExitCode = -1
		evidence.Match = MatchUnknown
		return VerifyResult{Evidence: evidence, Output: string(output)}, nil
	}

	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			evidence.ExitCode = exitErr.ExitCode()
		} else {
			evidence.ExitCode = -1
		}
	} else {
		evidence.ExitCode = 0
	}

	count, unit, known := deriveMatch(command, string(output))
	evidence.Matched = count
	evidence.MatchUnit = unit
	switch {
	case !known:
		evidence.Match = MatchUnknown
	case count == 0:
		evidence.Match = MatchRejected
	default:
		evidence.Match = MatchProven
	}

	return VerifyResult{Evidence: evidence, Output: string(output)}, nil
}

// shellCommand wraps a command line for the platform shell.
func shellCommand(ctx context.Context, command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/c", command)
	}
	return exec.CommandContext(ctx, "sh", "-c", command)
}

// deriveMatch counts what the command matched, if it knows how.
//
// The count is derived from the output the tool captured. When the command shape is not
// one we can interpret, the count is unknown and the evidence is unproven — honest,
// rather than silently passed.
func deriveMatch(command, output string) (count int, unit string, known bool) {
	if !isGoTest(command) {
		return 0, "", false
	}
	unit = "test"

	// go test says so explicitly when a filter matched nothing.
	if strings.Contains(output, "[no tests to run]") {
		return 0, unit, true
	}

	// Verbose output names each test result.
	if n := countMatches(output, "--- PASS:", "--- FAIL:", "--- SKIP:"); n > 0 {
		return n, unit, true
	}

	// Quiet output reports one line per package. That tells us something ran, but not
	// how many tests it was.
	if n := countPackageResults(output); n > 0 {
		return n, "package", true
	}

	return 0, unit, true
}

// isGoTest reports whether the command runs the Go test tool.
func isGoTest(command string) bool {
	fields := strings.Fields(command)
	if len(fields) < 2 {
		return false
	}
	// Skip env assignments and flags that precede the subcommand.
	i := 0
	for i < len(fields) && strings.Contains(fields[i], "=") && !strings.HasPrefix(fields[i], "go") {
		i++
	}
	if i+1 >= len(fields) {
		return false
	}
	return fields[i] == "go" && fields[i+1] == "test"
}

// countMatches counts lines containing any of the given markers.
func countMatches(output string, markers ...string) int {
	n := 0
	for _, line := range strings.Split(output, "\n") {
		for _, m := range markers {
			if strings.Contains(line, m) {
				n++
				break
			}
		}
	}
	return n
}

// countPackageResults counts Go's per-package result lines.
func countPackageResults(output string) int {
	n := 0
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "ok ") || strings.HasPrefix(trimmed, "FAIL ") {
			n++
		}
	}
	return n
}
