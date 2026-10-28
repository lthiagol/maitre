package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/brand"
)

// TestTheHelpOpensWithTheMark covers the help header.
//
// The help is text, and the mark is the first of it: the full form, above the title
// line.
func TestTheHelpOpensWithTheMark(t *testing.T) {
	var out, errb bytes.Buffer
	r := &Runner{Args: []string{"help"}, Stdout: &out, Stderr: &errb}
	r.Run()

	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(lines) < len(brand.Full)+2 {
		t.Fatalf("the help is %d lines — the mark and the title should be there:\n%s", len(lines), out.String())
	}
	for i, want := range brand.Full {
		if lines[i] != want {
			t.Errorf("line %d is %q, want %q", i+1, lines[i], want)
		}
	}
	title := lines[len(brand.Full)+1]
	if !strings.HasPrefix(title, "maitre-cli — the maitre agent CLI") {
		t.Errorf("the title should follow the mark, got %q", title)
	}
}
