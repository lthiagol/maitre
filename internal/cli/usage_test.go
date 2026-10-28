package cli

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// usageString finds a usage line in the source of any package.
//
// The usage lines are the answer when a command is called wrongly, and they name
// the command the user should run. A usage line that names a binary that does not
// exist sends the user somewhere that cannot help them.
var usageString = regexp.MustCompile(`"usage: ([a-z][a-z0-9-]*) `)

// repoRoot is the repository, two levels up from this package.
const repoRoot = "../.."

// sourceFile is one file the usage-line wire reads.
type sourceFile struct {
	// Path is where the file is, relative to the repository root.
	Path string
	// Content is what the file says.
	Content string
}

// TestTheUsageLinesNameTheBinary keeps the usage lines honest.
//
// The CLI prints one usage line per wrong call. Every one of them names the command
// to run — so every one of them must name a command that exists. The CLI's binary is
// `maitre-cli`; the old name is gone from the tree, and a usage line that still
// carried it would point the user at a command that cannot be run.
//
// **The wire covers the whole tree.** A usage line is not the CLI's alone: the
// scripts print one too, and a usage line added in another package would send the
// user to a command that does not exist. So the wire reads every package that
// prints a usage line, and it checks each one against the commands the repository
// actually ships — a set read from the tree, so a command added is a command the
// wire knows and a command renamed is one the wire stops accepting.
func TestTheUsageLinesNameTheBinary(t *testing.T) {
	commands := shippedCommands(t)
	files := readSources(t, repoRoot)
	if usageLinesIn(files) == 0 {
		t.Fatal("the tree has no usage lines — the check found nothing to hold")
	}
	for _, err := range checkUsageLines(files, commands) {
		t.Error(err)
	}
}

// TestTheUsageWireCoversEveryPackage keeps the wire's scope honest.
//
// The wire used to scan this package alone, so a usage line added in another
// package was never caught — and the tree has more than one package that prints
// one. Every package that ships a command has to be in the scan.
func TestTheUsageWireCoversEveryPackage(t *testing.T) {
	files := readSources(t, repoRoot)
	seen := map[string]bool{}
	for _, f := range files {
		seen[packageOf(f.Path)] = true
	}
	for _, want := range []string{
		"internal/cli",
		"internal/core",
		"internal/ui",
		"cmd/maitre-cli",
		"cmd/maitre",
		"scripts/stamp-plan",
	} {
		if !seen[want] {
			t.Errorf("the usage-line wire should cover %s", want)
		}
	}
}

// TestTheUsageWireTripsOnAnotherPackage keeps the wire honest about its scope.
//
// A usage line added in another package must fail the wire — that is the case the
// old, package-local wire missed. The probe lives outside `internal/cli`, so the
// CLI's own rule cannot catch it, and it names a command the repository does not
// ship.
func TestTheUsageWireTripsOnAnotherPackage(t *testing.T) {
	commands := shippedCommands(t)
	probe := sourceFile{
		Path:    "internal/ui/usage_probe.go",
		Content: "package ui\n\nconst usageProbe = \"usage: scale init\"\n",
	}
	errs := checkUsageLines([]sourceFile{probe}, commands)
	if len(errs) == 0 {
		t.Fatal("the wire should fail on a usage line in another package that names a command the repository does not ship")
	}
	if !strings.Contains(errs[0].Error(), "scale") {
		t.Errorf("the failure should name the command it caught: %v", errs[0])
	}
}

// checkUsageLines is the wire: every usage line names a command the repository
// ships, and the CLI's own name the CLI.
//
// It returns every failure, not the first, so one run says everything that is
// wrong.
func checkUsageLines(files []sourceFile, commands map[string]bool) []error {
	var out []error
	for _, f := range files {
		for _, m := range usageString.FindAllStringSubmatch(f.Content, -1) {
			name := m[1]
			switch {
			case !commands[name]:
				out = append(out, &usageError{path: f.Path, name: name})
			case inPackage(f.Path, "internal/cli") && name != "maitre-cli":
				out = append(out, &usageError{path: f.Path, name: name})
			}
		}
	}
	return out
}

// usageLinesIn counts the usage lines in the given sources.
func usageLinesIn(files []sourceFile) int {
	n := 0
	for _, f := range files {
		n += len(usageString.FindAllStringSubmatch(f.Content, -1))
	}
	return n
}

// usageError is a usage line the wire caught.
type usageError struct {
	path string
	name string
}

func (e *usageError) Error() string {
	return e.path + ": the usage line names " + e.name + ", which the repository does not ship"
}

// shippedCommands are the commands the repository ships: the binaries under `cmd/`
// and the scripts under `scripts/` that have a `main.go`.
//
// The set is read from the tree, so it cannot drift from what the repository has.
func shippedCommands(t *testing.T) map[string]bool {
	t.Helper()
	out := map[string]bool{}
	for _, dir := range []string{"cmd", "scripts"} {
		entries, err := os.ReadDir(filepath.Join(repoRoot, dir))
		if err != nil {
			t.Fatalf("read %s: %v", dir, err)
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			if _, err := os.Stat(filepath.Join(repoRoot, dir, e.Name(), "main.go")); err != nil {
				continue
			}
			out[e.Name()] = true
		}
	}
	if len(out) == 0 {
		t.Fatal("the repository ships no commands — the check found nothing to hold")
	}
	return out
}

// readSources are the Go sources of the repository, the tests left out: a usage
// line in a test is not something a user ever reads.
//
// The vendored skills, the binaries and the working directories are left out too:
// they are not the repository's code.
func readSources(t *testing.T, root string) []sourceFile {
	t.Helper()
	var out []sourceFile
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "skills", "bin", "work", "logo", "kitchen":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		out = append(out, sourceFile{Path: filepath.ToSlash(rel), Content: string(raw)})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// inPackage reports whether the file is in the given package directory.
func inPackage(path, dir string) bool {
	return path == dir || strings.HasPrefix(path, dir+"/")
}

// packageOf is the directory a file is in, relative to the repository root.
func packageOf(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[:i]
	}
	return ""
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
