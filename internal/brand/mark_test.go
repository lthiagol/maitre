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

// modulePath is the module of this repository.
const modulePath = "github.com/lthiagol/maitre/"

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
// **The rules are per form and per category**, and a category is what a file *is*,
// not what its name ends with:
//
//	| Form | code | rendered output | art | prose |
//	|---|---|---|---|---|
//	| the drawings | the canonical declaration, and nowhere else | — | — | — |
//	| the compact glyph | the canonical declaration, and nowhere else | may ride along | may ride along | may be quoted |
//
// The drawings are long and unique, so a second copy of one is a second
// declaration, wherever it sits. The glyph is one character, and it legitimately
// rides the rendered output — the goldens carry it on the header line — and the
// card's art — the wallpaper carries ninety copies of it. Prose may quote it: the
// brand sheet and the ledger name it.
//
// **The exemption is the declaration, not the file.** A copy pasted into the file
// that declares the drawing is a second declaration — the same defect the wire
// exists to catch, one file closer to home.
func TestTheArtIsDeclaredOnlyHere(t *testing.T) {
	if err := checkArtScope(repoRoot); err != nil {
		t.Fatal(err)
	}
}

// TestTheArtScopeWireCanFail proves the wire is a wire.
//
// Each way a copy can hide: in the brand package, in the rendered output, in the
// art, in a doc, in a file that is not code at all — and the glyph, in code and in
// a script.
func TestTheArtScopeWireCanFail(t *testing.T) {
	for _, c := range []struct {
		name string
		path string
		body string
	}{
		{"a second copy in the brand package", filepath.Join("internal", "brand", "zz_probe.go"), Full[0]},
		{"a copy in the rendered output", filepath.Join("internal", "ui", "testdata", "zz_probe.golden"), Full[0]},
		{"a copy in the art", filepath.Join("images", "logo", "zz_probe.svg"), Full[0]},
		{"a copy in a doc", filepath.Join("docs", "zz_probe.md"), Full[0]},
		{"a copy in a file that is not code", filepath.Join("internal", "ui", "testdata", "zz_probe.txt"), Full[0]},
		{"the glyph in another code file", filepath.Join("internal", "ui", "zz_probe.go"), Compact},
		{"the glyph in a script", filepath.Join("scripts", "zz_probe.py"), Compact},
	} {
		root := scratchTree(t)
		write(t, filepath.Join(root, c.path), c.body+"\n")
		if err := checkArtScope(root); err == nil {
			t.Errorf("the wire should fail on %s", c.name)
		}
	}
}

// TestTheDrawingBesideTheDeclarationFails covers the last hiding place.
//
// The wire's exemption is the **declaration**, not the file: a second copy in the
// file that declares the drawing is a second declaration, and it is the same
// defect the wire exists to catch — one file closer to home.
func TestTheDrawingBesideTheDeclarationFails(t *testing.T) {
	root := scratchTree(t)
	path := filepath.Join(root, "internal", "brand", "mark.go")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	pasted := string(raw) + "\nvar second = []string{\n\t\"" + Full[1] + "\",\n}\n"
	if err := os.WriteFile(path, []byte(pasted), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := checkArtScope(root); err == nil {
		t.Fatal("the wire should fail when a copy of the drawing sits beside the declaration")
	}
}

// checkArtScope is the scope wire.
//
// The drawings may live in the source files and in the canonical declaration, and
// nowhere else — in any category, in any file type. The compact glyph may live in
// the canonical declaration, in the rendered output, in the art and in prose, but
// it may not be pasted into another **code** file.
//
// The directories that are not read are **records**, not holes:
//
//	.git    the repository's own objects
//	bin    build output — the binaries carry the art as string literals
//	skills vendored third-party
//	work   the night's scratch — the briefs, the transcripts and the receipts quote
//	       what the console shows, art and all
func checkArtScope(root string) error {
	decl, err := findCanonicalDeclaration(root)
	if err != nil {
		return err
	}
	sources := map[string]bool{
		filepath.Join(root, "images", "logo", "ascii", "bow-tie-full.txt"): true,
		filepath.Join(root, "images", "logo", "ascii", "bow-tie-half.txt"): true,
	}
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
		if sources[path] {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if err := drawingDrift(rel, raw, decl); err != nil {
			return err
		}
		return glyphDrift(rel, raw, decl)
	})
}

// category is what a file in the tree is.
//
// The rules are per category. A rule that names an extension names a spelling; a
// rule that names a category names the thing.
type category int

const (
	// code is a file that ships code.
	code category = iota
	// output is rendered output — the goldens, which carry the mark because the
	// console draws it.
	output
	// art is the card's own drawing — the banner, the wallpaper, the sheets.
	art
	// prose is everything else: the docs, the ledger, the notes.
	prose
)

// codeExtensions are the extensions of the files that are code.
//
// The glyph's rule is about **code**, not about Go: a script is code too.
var codeExtensions = map[string]bool{
	".go": true, ".py": true, ".sh": true, ".bash": true,
	".js": true, ".ts": true, ".rb": true, ".pl": true, ".mk": true,
}

// codeNames are the files that are code and have no extension.
var codeNames = map[string]bool{"Makefile": true}

// categoryOf is what the file is.
func categoryOf(rel string) category {
	base := filepath.Base(rel)
	switch {
	case codeExtensions[filepath.Ext(base)] || codeNames[base]:
		return code
	case strings.HasSuffix(base, ".golden"):
		return output
	case isUnder(rel, "images"):
		return art
	default:
		return prose
	}
}

// isUnder reports whether the path is inside a directory of that name.
func isUnder(rel, dir string) bool {
	for _, part := range strings.Split(rel, "/") {
		if part == dir {
			return true
		}
	}
	return false
}

// canonicalDeclaration is the one declaration that holds the drawings.
//
// The wire's exemption is the **declaration**, not the file: the drawing may be in
// the one `var` that holds it — and nowhere else, not in another file and not in
// another corner of the declaring file.
type canonicalDeclaration struct {
	// file is the file that declares them, relative to the tree's root.
	file string
	// full is the byte range of the literal that holds the full drawing.
	full [2]int
	// half is the byte range of the literal that holds the half drawing.
	half [2]int
}

// findCanonicalDeclaration finds the one declaration that holds the drawings.
//
// It is found by parsing, so a file that does not parse — a pasted copy, say — is
// never mistaken for the declaration, and a rename of the declaration does not
// silently widen the exemption.
func findCanonicalDeclaration(root string) (*canonicalDeclaration, error) {
	dir := filepath.Join(root, "internal", "brand")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
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
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil, err
		}
		decl := &canonicalDeclaration{file: filepath.ToSlash(rel)}
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
				for i, n := range vs.Names {
					if i >= len(vs.Values) {
						continue
					}
					span := [2]int{
						fset.Position(vs.Values[i].Pos()).Offset,
						fset.Position(vs.Values[i].End()).Offset,
					}
					switch n.Name {
					case "Full":
						decl.full = span
					case "Half":
						decl.half = span
					}
				}
			}
		}
		if decl.full[1] > decl.full[0] && decl.half[1] > decl.half[0] {
			return decl, nil
		}
	}
	return nil, fmt.Errorf("no file under %s declares the drawings", dir)
}

// drawingDrift says whether the file carries a copy of a drawing.
//
// The exemption is the **declaration**: the full drawing may be in the literal
// that holds it, and the half drawing in the literal that holds it. Anywhere else
// — in another file, or in another corner of the declaring file — is a second
// copy.
func drawingDrift(rel string, raw []byte, decl *canonicalDeclaration) error {
	text := string(raw)
	full, half := text, text
	if rel == decl.file {
		full = maskSpan(text, decl.full)
		half = maskSpan(text, decl.half)
	}
	for _, line := range Full {
		if strings.Contains(full, line) {
			return fmt.Errorf("%s: carries a line of the full drawing — the drawing is declared once, in internal/brand", rel)
		}
	}
	for _, line := range Half {
		if strings.Contains(half, line) {
			return fmt.Errorf("%s: carries a line of the half drawing — the drawing is declared once, in internal/brand", rel)
		}
	}
	return nil
}

// glyphDrift says whether the file pastes the glyph into code.
//
// The glyph is one character: the rendered output carries it on the header line,
// and the card's art carries ninety copies of it. It may not be pasted into
// another **code** file — whatever language the code is written in.
func glyphDrift(rel string, raw []byte, decl *canonicalDeclaration) error {
	if rel == decl.file || categoryOf(rel) != code {
		return nil
	}
	if strings.Contains(string(raw), Compact) {
		return fmt.Errorf("%s: pastes the compact glyph into a code file — it is declared in internal/brand", rel)
	}
	return nil
}

// maskSpan takes the bytes of a range out of the text, leaving spaces behind.
//
// The offsets do not move, so a message about the file still points at the right
// place.
func maskSpan(text string, span [2]int) string {
	if span[0] < 0 || span[1] > len(text) || span[0] >= span[1] {
		return text
	}
	b := []byte(text)
	for i := span[0]; i < span[1]; i++ {
		if b[i] != '\n' {
			b[i] = ' '
		}
	}
	return string(b)
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
//
// **The rule is about the code that ships**, so the wire reads the **non-test**
// files of `internal/core` — a test file may import what it likes in order to
// check it.
func TestTheMarkIsNotTheDomain(t *testing.T) {
	if got := importsIn(t, "."); len(got) > 0 {
		t.Errorf("the mark should import nothing of ours, but it imports %v", got)
	}
	if got := importsIn(t, filepath.Join(repoRoot, "internal", "core")); len(got) > 0 {
		t.Errorf("the domain should not import the mark, but it imports %v", got)
	}
}

// TestTheBoundaryWireCanFail proves the wire is a wire.
//
// **Every spelling of an import must be caught** — single-line, grouped, blank —
// because a wire that catches one spelling is a wire that waits for the other
// spelling to arrive.
func TestTheBoundaryWireCanFail(t *testing.T) {
	for _, c := range []struct {
		name string
		src  string
	}{
		{"a single-line import", "package core\n\nimport \"github.com/lthiagol/maitre/internal/brand\"\n"},
		{"a grouped import", "package core\n\nimport (\n\t\"github.com/lthiagol/maitre/internal/brand\"\n)\n"},
		{"a blank import", "package core\n\nimport (\n\t_ \"github.com/lthiagol/maitre/internal/brand\"\n)\n"},
	} {
		dir := t.TempDir()
		write(t, filepath.Join(dir, "zz_probe.go"), c.src)
		if got := importsIn(t, dir); len(got) == 0 {
			t.Errorf("the wire should catch %s", c.name)
		}
	}
}

// importsIn are the packages of this repository that a directory's shipped code
// imports.
//
// **The imports are parsed, not read.** A line that starts with the module path is
// one spelling of an import; a grouped import and a blank import are other
// spellings, and a wire that reads lines waits for the other spellings to arrive.
//
// The tests are left out: the rule is about the code that ships.
func importsIn(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		for _, imp := range file.Imports {
			pkg := strings.Trim(imp.Path.Value, `"`)
			if strings.HasPrefix(pkg, modulePath) {
				out = append(out, pkg)
			}
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
