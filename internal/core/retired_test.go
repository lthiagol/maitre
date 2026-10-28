package core

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// writeRecipeFile puts a recipe file on disk, exactly as it is written.
func writeRecipeFile(t *testing.T, s *Store, name, body string) {
	t.Helper()
	dir := filepath.Join(s.Root, DirRecipes)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestTheCurrentFormatIsRead covers the reader's ordinary work.
//
// A plan in the current format reads. That is the whole point of the check: it
// refuses the retired shapes and nothing else.
func TestTheCurrentFormatIsRead(t *testing.T) {
	ops, _ := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatalf("a plan in the current format must read: %v", err)
	}
	if got.ID != r.ID {
		t.Fatalf("read %q, want %q", got.ID, r.ID)
	}
}

// TestTheRetiredBacklogDirectoriesAreRefused covers the directories.
//
// The backlog was two directories; it is one now. A plan that still carries the old
// ones is not supported, and the reader says so — it does not pretend the backlog is
// empty.
func TestTheRetiredBacklogDirectoriesAreRefused(t *testing.T) {
	for _, dir := range []string{"specials", "prep-notes"} {
		t.Run(dir, func(t *testing.T) {
			s := newStore(t)
			if err := s.Init(); err != nil {
				t.Fatal(err)
			}
			old := filepath.Join(s.Root, dir)
			if err := os.MkdirAll(old, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(old, "S01-a-hunch.json"),
				[]byte(`{"id":"S01","title":"A hunch","status":"open"}`), 0o644); err != nil {
				t.Fatal(err)
			}

			_, err := s.ListBacklogItems()
			if err == nil {
				t.Fatal("a plan that carries a retired directory must not be read")
			}
			var unsupported *UnsupportedPlanError
			if !errors.As(err, &unsupported) {
				t.Fatalf("the refusal should be an unsupported plan: %v", err)
			}
			if !strings.Contains(err.Error(), dir+"/") {
				t.Errorf("the refusal should name what the plan carries: %v", err)
			}
			if !strings.Contains(err.Error(), "not supported") {
				t.Errorf("the refusal should say the plan is not supported: %v", err)
			}
		})
	}
}

// TestTheRetiredFieldNamesAreRefused covers the records.
//
// The recipe's completion stamp was `evidence`; it is `completion` now. A
// criterion's record was `proofs`; it is `evidence` now. A record under the old
// name is a record the current reader does not know, and the reader says so.
func TestTheRetiredFieldNamesAreRefused(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string
	}{
		{
			name: "the old completion stamp",
			body: `{"id":"RCP01","title":"t","slug":"t","lifecycle":"closed",
				"evidence":{"completed_at":"2026-10-26T02:21:56Z"}}`,
			want: "`evidence` stamp",
		},
		{
			name: "the old criterion record",
			body: `{"id":"RCP01","title":"t","slug":"t","lifecycle":"executed",
				"acceptance":[{"id":"AC-01","description":"d","verification":"go test ./...",
				"proofs":[{"command":"go test ./...","exit_code":0,"matched":1,"match":"proven"}]}]}`,
			want: "`proofs`",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := newStore(t)
			writeRecipeFile(t, s, "RCP01-t.json", c.body)

			_, err := s.LoadRecipe("RCP01")
			if err == nil {
				t.Fatal("a record under a retired field name must not be read")
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("the refusal should name the retired field: %v", err)
			}

			// And the listing refuses it too: a plan is not half-read.
			if _, err := s.ListRecipes(); err == nil {
				t.Fatal("the listing must not read a plan that carries a retired field")
			}
		})
	}
}

// TestTheRetiredIDFormIsRefused covers the ids.
//
// A recipe's id was `R01`; it is `RCP01` now. A live record under the old form is
// a record of the retired format, and the reader refuses it.
func TestTheRetiredIDFormIsRefused(t *testing.T) {
	s := newStore(t)
	writeRecipeFile(t, s, "R01-t.json", `{"id":"R01","title":"t","slug":"t","lifecycle":"draft"}`)

	_, err := s.LoadRecipe("R01")
	if err == nil {
		t.Fatal("a record under a retired id form must not be read")
	}
	if !strings.Contains(err.Error(), "RCP01") {
		t.Errorf("the refusal should name the current form: %v", err)
	}
}

// TestTheRetiredBacklogIDFormIsRefused covers the backlog's ids.
//
// The backlog's ids were `S01` (a special) and `N01` (a prep note); they are `BL01`
// now.
func TestTheRetiredBacklogIDFormIsRefused(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(s.Root, DirBacklog)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "S01-a-hunch.json"),
		[]byte(`{"id":"S01","title":"A hunch","kind":"idea","status":"open"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := s.ListBacklogItems()
	if err == nil {
		t.Fatal("a backlog item under a retired id form must not be read")
	}
	if !strings.Contains(err.Error(), "BL01") {
		t.Errorf("the refusal should name the current form: %v", err)
	}
}

// TestTheJournalKeepsTheOldNames covers the exception.
//
// The journal is history: its events keep the names and the ids they were written
// with, forever. The reader resolves what the journal names, and it does not refuse
// a plan because its history is old.
func TestTheJournalKeepsTheOldNames(t *testing.T) {
	ops, store := newOps(t)
	r, err := ops.CreateRecipe(readySpec())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendJournalEvent(JournalEntry{
		Seq: 1, At: Now(), Kind: KindControl, Actor: "line-cook", Session: "s-old", Role: "line-cook",
		Recipe: r.ID,
	}); err != nil {
		t.Fatal(err)
	}

	j, err := store.LoadJournal()
	if err != nil {
		t.Fatalf("a journal with an old name must read: %v", err)
	}
	if len(j.Events) == 0 {
		t.Fatal("the journal is empty")
	}
	if j.Events[len(j.Events)-1].Role != "line-cook" {
		t.Errorf("the old role should be kept as it was written: %q", j.Events[len(j.Events)-1].Role)
	}
}

// TestAPlanWithARetiredArtifactIsNotWrittenTo covers the write path.
//
// **Not supported means no writes.** The reader refused such a plan and the writer
// wrote to it anyway — so `maitre-cli backlog create` kept handing records to a
// plan the tool does not support, and the message at read time was the only fence
// there was. The fence is on the write path now.
func TestAPlanWithARetiredArtifactIsNotWrittenTo(t *testing.T) {
	cases := []struct {
		name    string
		prepare func(t *testing.T, s *Store)
	}{
		{
			name: "a retired directory",
			prepare: func(t *testing.T, s *Store) {
				old := filepath.Join(s.Root, "specials")
				if err := os.MkdirAll(old, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(old, "S01-a-hunch.json"),
					[]byte(`{"id":"S01","title":"A hunch","status":"open"}`), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "a retired record",
			prepare: func(t *testing.T, s *Store) {
				writeRecipeFile(t, s, "RCP01-t.json", `{"id":"RCP01","title":"t","slug":"t","lifecycle":"draft",
					"acceptance":[{"id":"AC-01","description":"d","verification":"go test ./...",
					"proofs":[{"command":"go test ./...","exit_code":0,"matched":1,"match":"proven"}]}]}`)
			},
		},
		{
			name: "a retired file name",
			prepare: func(t *testing.T, s *Store) {
				writeRecipeFile(t, s, "R01-t.json", `{"id":"RCP01","title":"t","slug":"t","lifecycle":"draft"}`)
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := newStore(t)
			if err := s.Init(); err != nil {
				t.Fatal(err)
			}
			ops := NewOperations(s, AuditContext{Actor: "runner", Session: "s-1", Role: "runner"})
			c.prepare(t, s)

			_, err := ops.CreateBacklogItem("A hunch", "a body", nil, "idea", nil)
			if err == nil {
				t.Fatal("a plan with a retired artifact must not be written to")
			}
			var unsupported *UnsupportedPlanError
			if !errors.As(err, &unsupported) {
				t.Fatalf("the refusal should be an unsupported plan: %v", err)
			}

			// And the recipe is not written either.
			if _, err := ops.CreateRecipe(readySpec()); err == nil {
				t.Fatal("a plan with a retired artifact must not be written to")
			}

			// And nothing landed.
			entries, err := os.ReadDir(filepath.Join(s.Root, DirBacklog))
			if err != nil {
				t.Fatal(err)
			}
			for _, e := range entries {
				t.Errorf("the refusal wrote %s anyway", e.Name())
			}
		})
	}
}

// TestTheFileNameMustBeInTheCurrentFormat covers the note.
//
// A file named `R01-<slug>.json` is a retired id form on disk, and the file's name
// is part of the format: the plan is not supported, and the refusal says what the
// name should be.
func TestTheFileNameMustBeInTheCurrentFormat(t *testing.T) {
	s := newStore(t)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	writeRecipeFile(t, s, "R01-t.json", `{"id":"RCP01","title":"t","slug":"t","lifecycle":"draft"}`)

	_, err := s.LoadRecipe("RCP01")
	if err == nil {
		t.Fatal("a file named with a retired id form must not be read")
	}
	var unsupported *UnsupportedPlanError
	if !errors.As(err, &unsupported) {
		t.Fatalf("the refusal should be an unsupported plan: %v", err)
	}
	if !strings.Contains(err.Error(), "RCP01") {
		t.Errorf("the refusal should name the current form: %v", err)
	}
	if !strings.Contains(err.Error(), "R01-t.json") {
		t.Errorf("the refusal should name the file: %v", err)
	}

	// And the listing refuses it too: a plan is not half-read.
	if _, err := s.ListRecipes(); err == nil {
		t.Fatal("the listing must not read a plan whose files are named in a retired form")
	}
}

// ---------------------------------------------------------------------------
// The wire for "one format"
// ---------------------------------------------------------------------------

// retiredToken is one retired name the code must not carry.
type retiredToken struct {
	// text is the text the token names.
	text string
	// word is true when the text is a word: it must stand alone.
	//
	// "R01" must not hide inside "RCP01", and "expo" must not hide inside
	// "exported": the token names a retired thing, not a part of a word.
	word bool
	// fold is true when the case does not matter.
	//
	// A family is folded: `legacy` catches `LegacyThing` and `legacyFingerprint`
	// alike, so a retired name cannot grow back under a spelling the list did not
	// think of.
	fold bool
}

// retiredTokens are every retired name the code could grow back.
//
// The list is complete on purpose. A token that is missing is a name the code may
// carry without the wire noticing — and the wire is the only thing that says the
// format is one format.
//
// The retired **id literals** are in it (`R01`, `S01`, `N01`): the code never needs
// to write them, and a list that names them is the affordance the ruling removed.
//
// The retired **families** are in it too (`legacy`, `special`, `prep-note`,
// `proof`): a name the list spells out is a name that grows back the moment
// somebody spells it differently.
//
// The recipe's retired completion stamp is **not** in the list as the word
// `evidence`: `evidence` is a **current** word — the criterion's record — and the
// bare word cannot tell the retired stamp from the current one. Its token is the
// **access** instead: the retired stamp is read from a document by its literal key,
// and the current format never reads a document key by literal.
var retiredTokens = []retiredToken{
	{`["evidence"]`, false, false},
	{"proof", false, true},
	{"special", false, true},
	{"prep-note", false, true},
	{"R01", true, false},
	{"S01", true, false},
	{"N01", true, false},
	{"legacy", false, true},
	{"line-cook", true, false},
	{"taster", true, false},
	{"expo", true, false},
	{"coordinator", true, false},
}

// retiredPlace is one place a retired token may live.
//
// A place is a **declaration in a file**, not a file: the exemption is as narrow as
// the thing it excuses.
type retiredPlace struct {
	// file is the file, relative to the root the scan starts at.
	file string
	// decl is the name of the declaration.
	decl string
}

// retiredPlaces are the places the retired tokens may live.
//
// Two kinds of place, and no more:
//
//   - the **detection** that refuses them (`retired.go`), and
//   - the **deny-list** that refuses them as names (`model.go`).
//
// A retired name anywhere else is a second answer to a question the plan has already
// answered — and the tolerance comes back one reader at a time, which is how it came
// the first time.
//
// The places are declarations, not files. `model.go` may hold the retired names in
// `RetiredNames` and nowhere else, so a retired name cannot hide in it behind the
// exemption.
var retiredPlaces = []retiredPlace{
	{"internal/core/retired.go", "retiredDirectories"},
	{"internal/core/retired.go", "retiredInRecipe"},
	{"internal/core/retired.go", "retiredInBacklog"},
	{"internal/core/model.go", "RetiredNames"},
}

// scannedPackages are the packages the wire scans.
//
// Every package the repository has: a retired name in the console is as retired as
// one in the domain, and a wire that scans one package is a wire with a hole in it.
var scannedPackages = []string{"internal", "cmd", "scripts"}

// TestTheRetiredNamesAreOnlyWhereTheyBelong is the wire for "one format".
//
// The retired names may live in exactly two places: the detection that refuses them
// (`retired.go`) and the deny-list that refuses them as names (`model.go`). A third
// place is a second answer to a question the plan has already answered — and the
// tolerance would come back one reader at a time, which is how it came the first
// time.
//
// The scan reads the **code**, not the prose: the history is written down in the
// comments, and the comments are not a reader.
//
// The scan reads the **whole tree**: every package, not one. A retired name in the
// console is as retired as one in the domain.
func TestTheRetiredNamesAreOnlyWhereTheyBelong(t *testing.T) {
	if err := checkRetiredNames(filepath.Join("..", "..")); err != nil {
		t.Fatal(err)
	}
}

// TestTheOneFormatWireCanFail proves the wire is a wire.
//
// A tripwire that cannot fail is not a tripwire. The wire scans every package, so
// the proof is a retired token in each of them: the token goes into a scratch copy
// of the tree, and the wire must say no — and say which token it found.
func TestTheOneFormatWireCanFail(t *testing.T) {
	cases := []struct {
		pkg string
		// token is the retired name put into the package.
		token string
		// reported is the name the failure must carry — the family, when the
		// token is a name inside one.
		reported string
	}{
		{pkg: "internal/core", token: "proofs", reported: "proof"},
		{pkg: "internal/cli", token: "specials", reported: "special"},
		{pkg: "internal/ui", token: "R01", reported: "R01"},
		{pkg: "cmd/maitre-cli", token: "N01", reported: "N01"},
		{pkg: "cmd/maitre", token: "line-cook", reported: "line-cook"},
		{pkg: "scripts/stamp-plan", token: "prep-notes", reported: "prep-note"},
		// A name the list never spelled out: the family catches it.
		{pkg: "internal/ui", token: "LegacyThing", reported: "legacy"},
	}
	for _, c := range cases {
		t.Run(c.pkg+" carries "+c.token, func(t *testing.T) {
			scratch := scratchTree(t)
			if err := checkRetiredNames(scratch); err != nil {
				t.Fatalf("the scratch copy should be clean: %v", err)
			}
			file := firstGoFile(t, filepath.Join(scratch, c.pkg))
			raw, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			added := fmt.Sprintf("\n// A retired name, on purpose.\nconst retiredProbe = %q\n", c.token)
			if err := os.WriteFile(file, append(raw, []byte(added)...), 0o644); err != nil {
				t.Fatal(err)
			}
			err = checkRetiredNames(scratch)
			if err == nil {
				t.Fatalf("the wire should fail on a retired token in %s", c.pkg)
			}
			if !strings.Contains(err.Error(), c.reported) {
				t.Errorf("the failure should name the retired %q: %v", c.reported, err)
			}
		})
	}
}

// TestTheGlossaryNamesEachRetiredNameOnce is the glossary's rule, checked.
//
// A retired word is named **once** in the glossary: in the `_Avoid_` line of the
// entry that replaced it. A retired word in the prose is a word the tool still
// uses, and the glossary is what a reader trusts to say which word is which.
//
// The retired recipe-level stamp is not in the list: its name — `evidence` — is a
// **current** word too (the criterion's record), and the glossary names it in the
// `Completion` entry's `_Avoid_` line.
func TestTheGlossaryNamesEachRetiredNameOnce(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "GLOSSARY.md"))
	if err != nil {
		t.Fatal(err)
	}
	section := glossaryLanguage(string(raw))
	if section == "" {
		t.Fatal("the glossary has no language section — the wire would check nothing")
	}
	names := append([]string{}, RetiredNames...)
	names = append(names, "special", "prep-note", "proof")
	for _, name := range names {
		found := 0
		for _, line := range strings.Split(section, "\n") {
			if !glossaryNames(line, name) {
				continue
			}
			found++
			if !strings.HasPrefix(strings.TrimSpace(line), "_Avoid_:") {
				t.Errorf("the glossary names the retired %q outside an `_Avoid_` line: %s",
					name, strings.TrimSpace(line))
			}
		}
		if found != 1 {
			t.Errorf("the glossary names the retired %q %d times — it is named once, in the `_Avoid_` line",
				name, found)
		}
	}
}

// checkRetiredNames scans the packages under root and reports every retired token
// the code carries outside the places that may hold them.
func checkRetiredNames(root string) error {
	var found []string
	for _, pkg := range scannedPackages {
		dir := filepath.Join(root, pkg)
		if _, err := os.Stat(dir); err != nil {
			continue
		}
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			name := d.Name()
			if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			violations, err := retiredTokensIn(rel, path)
			if err != nil {
				return err
			}
			for _, v := range violations {
				found = append(found, fmt.Sprintf("%s carries the retired %q", rel, v))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	if len(found) == 0 {
		return nil
	}
	return fmt.Errorf("the retired names live outside the places that may hold them: %s",
		strings.Join(found, "; "))
}

// retiredTokensIn names the retired tokens a file carries outside the places that
// may hold them.
func retiredTokensIn(rel, path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	code, allowed, err := scanSource(rel, string(raw))
	if err != nil {
		return nil, err
	}
	var found []string
	for _, tok := range retiredTokens {
		for _, at := range findToken(code, tok) {
			if inAllowed(at, allowed) {
				continue
			}
			found = append(found, tok.text)
		}
	}
	return found, nil
}

// scanSource parses a Go file, takes the prose out of it, and finds the places in it
// that may hold a retired token.
//
// The comments are **masked**, not removed: every byte of a comment becomes a space,
// so an offset in the masked source is the offset it had in the file — and the scan
// of the code does not trip on the history the comments tell.
func scanSource(rel, src string) (string, [][2]int, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, rel, src, parser.ParseComments)
	if err != nil {
		return "", nil, fmt.Errorf("parse %s: %w", rel, err)
	}
	var b strings.Builder
	prev := 0
	for _, group := range file.Comments {
		for _, c := range group.List {
			start := fset.Position(c.Pos()).Offset
			end := fset.Position(c.End()).Offset
			b.WriteString(src[prev:start])
			b.WriteString(mask(src[start:end]))
			prev = end
		}
	}
	b.WriteString(src[prev:])
	return b.String(), allowedRanges(rel, file, fset), nil
}

// mask takes the prose out of a piece of text without moving the code.
//
// The newlines stay, so the masked text still has the shape of the file.
func mask(text string) string {
	out := make([]byte, 0, len(text))
	for i := 0; i < len(text); i++ {
		if text[i] == '\n' {
			out = append(out, '\n')
			continue
		}
		out = append(out, ' ')
	}
	return string(out)
}

// allowedRanges are the byte ranges of the declarations that may hold a retired
// token.
func allowedRanges(rel string, file *ast.File, fset *token.FileSet) [][2]int {
	var out [][2]int
	for _, place := range retiredPlaces {
		if place.file != rel {
			continue
		}
		start, end, ok := declRange(declNamed(file, place.decl), place.decl, fset)
		if !ok {
			continue
		}
		out = append(out, [2]int{start, end})
	}
	return out
}

// declNamed is the declaration of a package that goes by that name.
func declNamed(file *ast.File, name string) ast.Decl {
	for _, decl := range file.Decls {
		if declName(decl) == name {
			return decl
		}
	}
	return nil
}

// declName is the name a declaration introduces.
func declName(decl ast.Decl) string {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return d.Name.Name
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, n := range vs.Names {
				return n.Name
			}
		}
	}
	return ""
}

// declRange is the byte range of the declaration a place names.
//
// A variable's range is its **value**, not the whole declaration: a `var` block
// with two names is two places, and the exemption is one of them.
func declRange(decl ast.Decl, name string, fset *token.FileSet) (int, int, bool) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		if d.Name.Name != name {
			return 0, 0, false
		}
		return fset.Position(d.Pos()).Offset, fset.Position(d.End()).Offset, true
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, n := range vs.Names {
				if n.Name != name {
					continue
				}
				return fset.Position(vs.Pos()).Offset, fset.Position(vs.End()).Offset, true
			}
		}
	}
	return 0, 0, false
}

// inAllowed reports whether an offset falls inside one of the ranges.
func inAllowed(at int, ranges [][2]int) bool {
	for _, r := range ranges {
		if at >= r[0] && at < r[1] {
			return true
		}
	}
	return false
}

// findToken finds every place the token is in the code.
func findToken(code string, tok retiredToken) []int {
	if tok.fold {
		// A family: the case does not matter, and the name may be a part of a
		// longer one.
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(tok.text))
		var out []int
		for _, m := range re.FindAllStringIndex(code, -1) {
			out = append(out, m[0])
		}
		return out
	}
	if !tok.word {
		var out []int
		for i := 0; i+len(tok.text) <= len(code); {
			j := strings.Index(code[i:], tok.text)
			if j < 0 {
				break
			}
			out = append(out, i+j)
			i = i + j + len(tok.text)
		}
		return out
	}
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(tok.text) + `\b`)
	var out []int
	for _, m := range re.FindAllStringIndex(code, -1) {
		out = append(out, m[0])
	}
	return out
}

// scratchTree copies the packages the wire scans into a scratch directory.
//
// The copy keeps the tree's shape: the places that may hold a retired token are
// named by their path, so a copy that flattened them would not be the same tree.
func scratchTree(t *testing.T) string {
	t.Helper()
	root := filepath.Join("..", "..")
	dir := t.TempDir()
	for _, pkg := range scannedPackages {
		base := filepath.Join(root, pkg)
		if _, err := os.Stat(base); err != nil {
			continue
		}
		err := filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			name := d.Name()
			if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			target := filepath.Join(dir, rel)
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(target, raw, 0o644)
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// firstGoFile is the first Go file in a package directory, in name order.
func firstGoFile(t *testing.T, dir string) string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		return filepath.Join(dir, name)
	}
	t.Fatalf("no Go file in %s", dir)
	return ""
}

// glossaryLanguage is the glossary's language section — the entries, and the
// `_Avoid_` lines under them.
//
// The CPD mapping lives in another section, and it names CPD's own words: the
// glossary says the two vocabularies are different things.
func glossaryLanguage(doc string) string {
	start := strings.Index(doc, "## Language")
	if start < 0 {
		return ""
	}
	rest := doc[start:]
	if i := strings.Index(rest, "## Two vocabularies"); i >= 0 {
		return rest[:i]
	}
	return rest
}

// glossaryNames reports whether a line names the retired name.
//
// The comparison is the way the code compares names: the case does not matter, and
// a space stands for a dash — so "Line cook" is the name "line-cook".
func glossaryNames(line, name string) bool {
	parts := strings.Split(name, "-")
	for i, p := range parts {
		parts[i] = regexp.QuoteMeta(p)
	}
	pattern := `(?i)\b` + strings.Join(parts, `[- ]`) + `\b`
	return regexp.MustCompile(pattern).MatchString(line)
}
