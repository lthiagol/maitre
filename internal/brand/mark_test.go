package brand

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// artDir is where the drawings live, relative to this package.
const artDir = "../../images/logo/ascii"

// repoRoot is the repository, two levels up from this package.
const repoRoot = "../.."

// TestTheDrawingsMatchTheSource holds the code to the art.
//
// **The comparison runs from the file to the code.** The drawing in
// `images/logo/ascii/` is the source — the code holds it, so the file is what the
// code must match. A character changed in the `.txt` fails this wire, and the
// failure names the line that moved, so the fix is to move the code.
//
// The wire covers **both** drawings. The full form and the half form are two
// drawings, and a wire that held only one would let the other drift.
func TestTheDrawingsMatchTheSource(t *testing.T) {
	for _, c := range []struct {
		name string
		file string
		code []string
	}{
		{"the full form", "bow-tie-full.txt", Full},
		{"the half form", "bow-tie-half.txt", Half},
	} {
		lines, err := readDrawing(c.file)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if err := checkDrawing(lines, c.code); err != nil {
			t.Errorf("%s: %v", c.name, err)
		}
	}
}

// TestTheDrawingWireCanFail proves the wire is a wire.
//
// A tripwire that cannot fail is not a tripwire. The drawing is edited in memory —
// one character of one line — and the wire must say no.
func TestTheDrawingWireCanFail(t *testing.T) {
	for _, c := range []struct {
		name string
		file string
		code []string
	}{
		{"the full form", "bow-tie-full.txt", Full},
		{"the half form", "bow-tie-half.txt", Half},
	} {
		lines, err := readDrawing(c.file)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		edited := append([]string(nil), lines...)
		edited[0] = editFirstRune(edited[0])
		if err := checkDrawing(edited, c.code); err == nil {
			t.Errorf("%s: the wire should fail when the drawing moves and the code does not", c.name)
		}
	}
}

// TestTheHalfFormIsItsOwnDrawing covers the claim.
//
// The half form is a **second drawing of the same emblem**, not a crop of the full
// one — so no window cut out of the full form can be the half form. The wire walks
// every window and says so.
func TestTheHalfFormIsItsOwnDrawing(t *testing.T) {
	full, half := Width(Full), Width(Half)
	for off := 0; off+half <= full; off++ {
		crop := make([]string, len(Full))
		for i, l := range Full {
			crop[i] = cropRunes(l, off, half)
		}
		if err := checkDrawing(crop, Half); err == nil {
			t.Fatalf("the half form is the full form cut at column %d — it is a crop, not a drawing", off)
		}
	}
}

// TestTheDrawingsAreTheSizesTheSheetSays holds the widths.
//
// The brand sheet says the full form is 28 columns and the half form is 14, and
// the console's rule is built on those numbers.
func TestTheDrawingsAreTheSizesTheSheetSays(t *testing.T) {
	if got := Width(Full); got != 28 {
		t.Errorf("the full form is %d columns, want 28", got)
	}
	if got := Width(Half); got != 14 {
		t.Errorf("the half form is %d columns, want 14", got)
	}
	if got := Width([]string{Compact}); got != 1 {
		t.Errorf("the compact glyph is %d columns, want 1", got)
	}
}

// TestTheArtIsDeclaredOnlyHere holds the art to one place.
//
// The mark is declared once, in this package, and every consumer reads it from
// here. A drawing pasted into the console or the CLI would be a second answer to
// the same question — and the wire above would never see it, because it only reads
// this package. So this wire looks for the mark's own characters anywhere else in
// the tree.
func TestTheArtIsDeclaredOnlyHere(t *testing.T) {
	err := filepath.WalkDir(repoRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "skills", "bin", "work", "kitchen", "images":
				return filepath.SkipDir
			}
			if path == filepath.Join(repoRoot, "internal", "brand") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if line, ok := markCharacter(string(raw)); ok {
			rel, _ := filepath.Rel(repoRoot, path)
			return fmt.Errorf("%s: line %d carries the mark itself — the art is declared once, in internal/brand", rel, line)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestTheMarkIsNotTheDomain covers the boundary.
//
// The art is not the domain: this package imports nothing of ours, and the domain
// does not import it. The domain is the model, and a bow tie is UI.
func TestTheMarkIsNotTheDomain(t *testing.T) {
	if got := oursIn(t, "."); len(got) > 0 {
		t.Errorf("the mark should import nothing of ours, but it imports %v", got)
	}
	if got := oursIn(t, filepath.Join(repoRoot, "internal", "core")); len(got) > 0 {
		t.Errorf("the domain should not import the mark, but it imports %v", got)
	}
}

// oursIn are the packages of this repository that a directory's code imports.
//
// The tests are left out: the rule is about the code that ships.
func oursIn(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, line := range strings.Split(string(raw), "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, `"github.com/lthiagol/maitre/`) {
				continue
			}
			pkg := strings.Trim(line, `"`)
			out = append(out, pkg)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// readDrawing reads a drawing from the source directory.
func readDrawing(name string) ([]string, error) {
	raw, err := os.ReadFile(filepath.Join(artDir, name))
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimRight(string(raw), "\n"), "\n")
	return lines, nil
}

// checkDrawing says whether the drawing in the file and the drawing in the code are
// the same one.
//
// It returns the first difference it finds, named — a wire that says "they differ"
// sends the reader looking.
func checkDrawing(file, code []string) error {
	if len(file) != len(code) {
		return fmt.Errorf("the drawing has %d lines and the code has %d", len(file), len(code))
	}
	for i := range file {
		if file[i] != code[i] {
			return fmt.Errorf("line %d moved: the drawing says %q and the code says %q", i+1, file[i], code[i])
		}
	}
	return nil
}

// editFirstRune changes the first character of a line.
func editFirstRune(line string) string {
	r := []rune(line)
	if len(r) == 0 {
		return "⣿"
	}
	if r[0] == '⣿' {
		r[0] = '⠀'
	} else {
		r[0] = '⣿'
	}
	return string(r)
}

// cropRunes is a window of a line: the runes from `off` to `off+n`.
//
// A window that runs past the end of the line is the empty string, so a crop of a
// short line is compared as a short line.
func cropRunes(line string, off, n int) string {
	r := []rune(line)
	if off >= len(r) {
		return ""
	}
	end := off + n
	if end > len(r) {
		end = len(r)
	}
	return string(r[off:end])
}

// markCharacter finds the mark's own characters in a source file.
//
// The braille block is the drawing, and `⋈` is the compact form. A file that
// carries one of them has the art pasted into it.
func markCharacter(src string) (int, bool) {
	for i, line := range strings.Split(src, "\n") {
		for _, r := range line {
			if (r >= '⠀' && r <= '⣿') || r == '⋈' {
				return i + 1, true
			}
		}
	}
	return 0, false
}
