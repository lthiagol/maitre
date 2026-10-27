package cli

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// usageString finds a usage line in the CLI's own source.
//
// The usage lines are the CLI's answer when a command is called wrongly, and they
// name the command the user should run. A usage line that names a binary that does
// not exist sends the user somewhere that cannot help them.
var usageString = regexp.MustCompile(`"usage: ([a-z][a-z0-9-]*) `)

// TestTheUsageLinesNameTheBinary keeps the usage lines honest.
//
// The CLI prints one usage line per wrong call. Every one of them names the command
// to run — so every one of them must name a binary that exists. The binary is
// `maitre-cli`; the old name is gone from the tree, and a usage line that still
// carried it would point the user at a command that cannot be run.
func TestTheUsageLinesNameTheBinary(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var checked int
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range usageString.FindAllStringSubmatch(string(raw), -1) {
			checked++
			if m[1] != "maitre-cli" {
				t.Errorf("%s: the usage line names %q, want maitre-cli", f, m[1])
			}
		}
	}
	if checked == 0 {
		t.Fatal("the CLI has no usage lines — the check found nothing to hold")
	}
}

// TestTheConsoleErrorsNameTheBinary covers the console's own errors.
//
// The console is a separate binary, so its errors are a second place a removed
// command can hide. The chef must never be sent to a command that does not exist.
func TestTheConsoleErrorsNameTheBinary(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "cmd", "maitre", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	if !strings.Contains(src, "`maitre-cli init`") {
		t.Errorf("the console's error should tell the user to run `maitre-cli init`:\n%s", src)
	}
	for _, gone := range []string{"`scale init`", "`scale`"} {
		if strings.Contains(src, gone) {
			t.Errorf("the console's error still names %s", gone)
		}
	}
}
