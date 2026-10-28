package brand

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
// The brand sheet says the full form is 28 columns at its widest and the half form
// is 14, and the console's rule is built on those numbers.
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

// TestTheArtIsDeclaredOnlyHere holds the mark to one place.
//
// **Two things are in play, and the rule for each is different.**
//
// | Thing | Rule |
// |---|---|
// | the drawings | long and unique: they may live in the source files and in the declaration, and **nowhere else** — in any file type |
// | the compact glyph | one character: it may appear in rendered output and in the card's art, but it must not be **pasted into another code file** |
//
// The glyph's rule is narrower because the glyph is legitimate elsewhere: the
// goldens carry it on the header line, and the wallpaper carries ninety copies of
// it — that is the card's own drawing.
//
// **The exemption is per declaration, not per package.** A second copy pasted into
// another file of the same package is still a second copy, and a per-package skip
// is exactly where the copy hides.
func TestTheArtIsDeclaredOnlyHere(t *testing.T) {
	if err := checkArtScope(repoRoot); err != nil {
		t.Fatal(err)
	}
}

// TestTheArtScopeWireCanFail proves the wire is a wire.
//
// The reviewer's two mutations — a second copy in the brand package, and a copy in
// a file that is not Go at all — and the glyph pasted into another code file. Each
// one must fail the wire.
func TestTheArtScopeWireCanFail(t *testing.T) {
	for _, c := range []struct {
		name string
		path string
		body string
	}{
		{"a second copy in the brand package", filepath.Join("internal", "brand", "zz_probe.go"), Full[0]},
		{"a copy in a file that is not Go", filepath.Join("internal", "ui", "testdata", "zz_probe.txt"), Full[0]},
		{"the glyph in another code file", filepath.Join("internal", "ui", "zz_probe.go"), Compact},
	} {
		root := scratchTree(t)
		write(t, filepath.Join(root, c.path), c.body+"\n")
		if err := checkArtScope(root); err == nil {
			t.Errorf("the wire should fail on %s", c.name)
		}
	}
}

// checkArtScope is the scope wire.
//
// The drawings may live in the source files and in the declaration, and nowhere
// else — in any file type. The compact glyph may live in the declaration, in
// rendered output and in the card's art, but it must not be pasted into another
// code file.
//
// The directories that are not read are the ones that are not the code that ships:
// the git objects, the built binaries (they carry the art as string literals), the
// vendored skills, and the night's scratch (the transcripts capture what the
// console shows, art and all).
func checkArtScope(root string) error {
	decl, err := declarationFile(root)
	if err != nil {
		return err
	}
	sources := map[string]bool{
		filepath.Join(root, "images", "logo", "ascii", "bow-tie-full.txt"): true,
		filepath.Join(root, "images", "logo", "ascii", "bow-tie-half.txt"): true,
	}
	lines := drawingLines()

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "skills", "bin", "work":
				return filepath.SkipDir
			}
			return nil
		}
		if path == decl || sources[path] {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(raw)
		for _, line := range lines {
			if strings.Contains(text, line) {
				rel, _ := filepath.Rel(root, path)
				return fmt.Errorf("%s: carries a line of the drawing — the drawings are declared once, in internal/brand", rel)
			}
		}
		if strings.HasSuffix(path, ".go") && strings.Contains(text, Compact) {
			rel, _ := filepath.Rel(root, path)
			return fmt.Errorf("%s: carries the compact glyph — it is declared in internal/brand", rel)
		}
		return nil
	})
}

// declarationFile is the file that declares the drawings.
//
// The scope's exemption is **per declaration**: the file that actually declares
// `Full` and `Half`. It is found by parsing, so a file that does not parse — a
// pasted copy, say — is never mistaken for the declaration, and a rename of the
// declaration does not silently widen the exemption.
func declarationFile(root string) (string, error) {
	dir := filepath.Join(root, "internal", "brand")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			continue
		}
		if declares(f, "Full") && declares(f, "Half") {
			return path, nil
		}
	}
	return "", fmt.Errorf("no file under %s declares the drawings", dir)
}

// declares reports whether the file declares a package-level variable of that name.
func declares(f *ast.File, name string) bool {
	for _, d := range f.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, n := range vs.Names {
				if n.Name == name {
					return true
				}
			}
		}
	}
	return false
}

// drawingLines are the lines of the drawings — the full form and the half form.
//
// A line is long and unique, so a file that carries one of them carries a copy of
// the drawing.
func drawingLines() []string {
	out := make([]string, 0, len(Full)+len(Half))
	out = append(out, Full...)
	out = append(out, Half...)
	return out
}

// scratchTree is a copy of the parts of the tree the scope wire reads.
func scratchTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, c := range []struct{ from, to string }{
		{filepath.Join(repoRoot, "internal", "brand", "mark.go"), filepath.Join("internal", "brand", "mark.go")},
		{filepath.Join(artDir, "bow-tie-full.txt"), filepath.Join("images", "logo", "ascii", "bow-tie-full.txt")},
		{filepath.Join(artDir, "bow-tie-half.txt"), filepath.Join("images", "logo", "ascii", "bow-tie-half.txt")},
	} {
		raw, err := os.ReadFile(c.from)
		if err != nil {
			t.Fatal(err)
		}
		write(t, filepath.Join(root, c.to), string(raw))
	}
	return root
}

// write puts a file on the disk, making its directory.
func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
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
//
// The change is an ordinary letter, so the test file carries none of the mark's own
// characters.
func editFirstRune(line string) string {
	r := []rune(line)
	if len(r) == 0 {
		return "X"
	}
	if r[0] == 'X' {
		r[0] = 'Y'
	} else {
		r[0] = 'X'
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
