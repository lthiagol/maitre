# The runner's lane — slice G, captured 2026-10-28

Dumped from pane `w1Z:pG` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that wrote slice G: the mark in one package, the three forms, the wires that
hold the code to the drawing, and the rule that the decoration yields to the content
in both dimensions.
The review that judged it is in [slice-g-reviewer.md](slice-g-reviewer.md).

```
  ┃      mark_test.go:145: open ../../../../images/logo/ascii/bow-tie-full.txt: no such file or directory
  ┃  --- FAIL: TestTheArtScopeWireCanFail (0.00s)
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/brand  0.326s
  ┃  FAIL
  ┃

     + Thought · 675ms

  ┃
  ┃  ← Edit internal/brand/mark_test.go
  ┃
  ┃    272     t.Helper()                                               272     t.Helper()
  ┃    273     root := t.TempDir()                                      273     root := t.TempDir()
  ┃    274     for _, c := range []struct{ from, to string }{           274     for _, c := range []struct{ from, to string }{
  ┃    275       {filepath.Join(repoRoot, "internal", "brand", "mark.   275       {filepath.Join(repoRoot, "internal", "brand", "mark.
  ┃          go"), filepath.Join("internal", "brand", "mark.go")},            go"), filepath.Join("internal", "brand", "mark.go")},
  ┃    276 -     {filepath.Join(repoRoot, artDir, "bow-tie-full.txt"),  276 +     {filepath.Join(artDir, "bow-tie-full.txt"), filepath.
  ┃          filepath.Join("images", "logo", "ascii", "bow-tie-full.          Join("images", "logo", "ascii", "bow-tie-full.txt")},
  ┃          txt")},
  ┃    277 -     {filepath.Join(repoRoot, artDir, "bow-tie-half.txt"),  277 +     {filepath.Join(artDir, "bow-tie-half.txt"), filepath.
  ┃          filepath.Join("images", "logo", "ascii", "bow-tie-half.          Join("images", "logo", "ascii", "bow-tie-half.txt")},
  ┃          txt")},
  ┃    278     } {                                                      278     } {
  ┃    279       raw, err := os.ReadFile(c.from)                        279       raw, err := os.ReadFile(c.from)
  ┃    280       if err != nil {                                        280       if err != nil {
  ┃    281         t.Fatal(err)                                         281         t.Fatal(err)
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && gofmt -l . ; go test ./internal/brand -v 2>&1 | grep -E '^(=== RUN|--- |ok|FAIL|    )' |
  ┃   head -30
  ┃
  ┃  (7 earlier lines)
  ┃  --- PASS: TestTheDrawingsAreTheSizesTheSheetSays (0.00s)
  ┃  === RUN   TestTheArtIsDeclaredOnlyHere
  ┃  --- PASS: TestTheArtIsDeclaredOnlyHere (0.02s)
  ┃  === RUN   TestTheArtScopeWireCanFail
  ┃  --- PASS: TestTheArtScopeWireCanFail (0.01s)
  ┃  ok    github.com/lthiagol/maitre/internal/brand  0.214s
  ┃

     + Thought · 1.3s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
  ┃  art = open('images/logo/ascii/bow-tie-full.txt', encoding='utf-8').read().rstrip('\n').split('\n')[1]…
  ┃
  ┃  --- FAIL: TestTheArtIsDeclaredOnlyHere (0.00s)
  ┃      mark_test.go:126: internal/brand/zz_probe.go: carries a line of the drawing — the drawings are declared once, in internal/
  ┃  brand
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/brand  0.175s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  68.696s
  ┃  FAIL
  ┃   M internal/brand/mark_test.go
  ┃

     + Thought · 365ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 -c "
  ┃  art = open('images/logo/ascii/bow-tie-full.txt', encoding='utf-8').read().rstrip('\n').split('\n')[1]…
  ┃
  ┃  (6 earlier lines)
  ┃  --- FAIL: TestTheArtIsDeclaredOnlyHere (0.01s)
  ┃      mark_test.go:126: internal/ui/zz_probe.go: carries the compact glyph — it is declared in internal/brand
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/brand  0.185s
  ┃  FAIL
  ┃   M internal/brand/mark_test.go
  ┃

     + Thought · 172ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/brand/mark_test.go && git commit -q -F - <<'EOF'
  ┃  fix: the wire that says where the mark may be, and that can fail…
  ┃
  ┃  42d1e8e fix: the wire that says where the mark may be, and that can fail
  ┃

     + Thought · 5.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git show 657d4f0:internal/brand/mark_test.go | sed -n '/^\/\/ TestTheMarkIsNotTheDomain/
  ┃  ,/^\/\/ readDrawing/p' | head -80
  ┃
  ┃  (41 earlier lines)
  ┃      t.Fatal(err)
  ┃    }
  ┃    return out
  ┃  }
  ┃
  ┃  // readDrawing reads a drawing from the source directory.
  ┃

  ┃
  ┃  ← Edit internal/brand/mark_test.go
  ┃
  ┃    295       t.Fatal(err)                                           295       t.Fatal(err)
  ┃    296     }                                                        296     }
  ┃    297   }                                                          297   }
  ┃    298                                                              298
  ┃                                                                     299 + // TestTheMarkIsNotTheDomain covers the boundary.
  ┃                                                                     300 + //
  ┃                                                                     301 + // The art is not the domain: this package imports
  ┃                                                                           nothing of ours, and the domain
  ┃                                                                     302 + // does not import it. The domain is the model, and a bow
  ┃                                                                           tie is UI.
  ┃                                                                     303 + func TestTheMarkIsNotTheDomain(t *testing.T) {
  ┃                                                                     304 +   if got := oursIn(t, "."); len(got) > 0 {
  ┃                                                                     305 +     t.Errorf("the mark should import nothing of ours, but
  ┃                                                                           it imports %v", got)
  ┃                                                                     306 +   }
  ┃                                                                     307 +   if got := oursIn(t, filepath.Join(repoRoot, "internal",
  ┃                                                                           "core")); len(got) > 0 {
  ┃                                                                     308 +     t.Errorf("the domain should not import the mark, but
  ┃                                                                           it imports %v", got)
  ┃                                                                     309 +   }
  ┃                                                                     310 + }
  ┃                                                                     311 +
  ┃                                                                     312 + // oursIn are the packages of this repository that a
  ┃                                                                           directory's code imports.
  ┃                                                                     313 + //
  ┃                                                                     314 + // The tests are left out: the rule is about the code
  ┃                                                                           that ships.
  ┃                                                                     315 + func oursIn(t *testing.T, dir string) []string {
  ┃                                                                     316 +   t.Helper()
  ┃                                                                     317 +   var out []string
  ┃                                                                     318 +   err := filepath.WalkDir(dir, func(path string, d fs.
  ┃                                                                           DirEntry, err error) error {
  ┃                                                                     319 +     if err != nil {
  ┃                                                                     320 +       return err
  ┃                                                                     321 +     }
  ┃                                                                     322 +     if d.IsDir() || !strings.HasSuffix(path, ".go") ||
  ┃                                                                           strings.HasSuffix(path, "_test.go") {
  ┃                                                                     323 +       return nil
  ┃                                                                     324 +     }
  ┃                                                                     325 +     raw, err := os.ReadFile(path)
  ┃                                                                     326 +     if err != nil {
  ┃                                                                     327 +       return err
  ┃                                                                     328 +     }
  ┃                                                                     329 +     for _, line := range strings.Split(string(raw), "\n")
  ┃                                                                           {
  ┃                                                                     330 +       line = strings.TrimSpace(line)
  ┃                                                                     331 +       if !strings.HasPrefix(line, `"github.com/lthiagol/
  ┃                                                                           maitre/`) {
  ┃                                                                     332 +         continue
  ┃                                                                     333 +       }
  ┃                                                                     334 +       pkg := strings.Trim(line, `"`)
  ┃                                                                     335 +       out = append(out, pkg)
  ┃                                                                     336 +     }
  ┃                                                                     337 +     return nil
  ┃                                                                     338 +   })
  ┃                                                                     339 +   if err != nil {
  ┃                                                                     340 +     t.Fatal(err)
  ┃                                                                     341 +   }
  ┃                                                                     342 +   return out
  ┃                                                                     343 + }
  ┃                                                                     344 +
  ┃    299   // readDrawing reads a drawing from the source directory.  345   // readDrawing reads a drawing from the source directory.
  ┃    300   func readDrawing(name string) ([]string, error) {          346   func readDrawing(name string) ([]string, error) {
  ┃    301     raw, err := os.ReadFile(filepath.Join(artDir, name))     347     raw, err := os.ReadFile(filepath.Join(artDir, name))
  ┃    302     if err != nil {                                          348     if err != nil {
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && gofmt -l . ; go test ./internal/brand 2>&1 | tail -3 && git add internal/brand/mark_test
  ┃  .go && git commit -q --amend --no-edit && git log --oneline -3 && git show --stat HEAD
  ┃
  ┃  (28 earlier lines)
  ┃      Per-AC evidence:
  ┃      - F-95: go test ./internal/brand -v exit 0; a copy in internal/brand, a copy
  ┃        in a .txt file and the glyph in another .go file each fail the wire
  ┃
  ┃   internal/brand/mark_test.go | 214 ++++++++++++++++++++++++++++++++++++--------
  ┃   1 file changed, 178 insertions(+), 36 deletions(-)
  ┃

     + Thought · 4.7s

     → Explored: 1 search, 1 read

     + Thought · 3.3s

  ┃
  ┃  ← Edit internal/brand/mark_test.go
  ┃
  ┃    299   // TestTheMarkIsNotTheDomain covers the boundary.          299   // TestTheMarkIsNotTheDomain covers the boundary.
  ┃    300   //                                                         300   //
  ┃    301   // The art is not the domain: this package imports nothing 301   // The art is not the domain: this package imports
  ┃          of ours, and the domain                                          nothing of ours, and the domain
  ┃    302   // does not import it. The domain is the model, and a bow  302   // does not import it. The domain is the model, and a bow
  ┃          tie is UI.                                                       tie is UI.
  ┃                                                                     303 + //
  ┃                                                                     304 + // **The rule is about the code that ships**, so the wire
  ┃                                                                           reads the **non-test**
  ┃                                                                     305 + // files of `internal/core` — a test file may import what
  ┃                                                                           it likes in order to
  ┃                                                                     306 + // check it.
  ┃    303   func TestTheMarkIsNotTheDomain(t *testing.T) {             307   func TestTheMarkIsNotTheDomain(t *testing.T) {
  ┃    304 -   if got := oursIn(t, "."); len(got) > 0 {                 308 +   if got := importsIn(t, "."); len(got) > 0 {
  ┃    305       t.Errorf("the mark should import nothing of ours, but  309       t.Errorf("the mark should import nothing of ours, but
  ┃          it imports %v", got)                                             it imports %v", got)
  ┃    306     }                                                        310     }
  ┃    307 -   if got := oursIn(t, filepath.Join(repoRoot, "internal",  311 +   if got := importsIn(t, filepath.Join(repoRoot,
  ┃          "core")); len(got) > 0 {                                         "internal", "core")); len(got) > 0 {
  ┃    308       t.Errorf("the domain should not import the mark, but   312       t.Errorf("the domain should not import the mark, but
  ┃          it imports %v", got)                                             it imports %v", got)
  ┃    309     }                                                        313     }
  ┃    310   }                                                          314   }
  ┃    311                                                              315
  ┃    312 - // oursIn are the packages of this repository that a       316 + // TestTheBoundaryWireCanFail proves the wire is a wire.
  ┃          directory's code imports.
  ┃    313   //                                                         317   //
  ┃                                                                     318 + // **Every spelling of an import must be caught** —
  ┃                                                                           single-line, grouped, blank —
  ┃                                                                     319 + // because a wire that catches one spelling is a wire
  ┃                                                                           that waits for the other
  ┃                                                                     320 + // spelling to arrive.
  ┃                                                                     321 + func TestTheBoundaryWireCanFail(t *testing.T) {
  ┃                                                                     322 +   for _, c := range []struct {
  ┃                                                                     323 +     name string
  ┃                                                                     324 +     src  string
  ┃                                                                     325 +   }{
  ┃                                                                     326 +     {"a single-line import", "package core\n\nimport \
  ┃                                                                           "github.com/lthiagol/maitre/internal/brand\"\n"},
  ┃                                                                     327 +     {"a grouped import", "package core\n\nimport (\n\t\
  ┃                                                                           "github.com/lthiagol/maitre/internal/brand\"\n)\n"},
  ┃                                                                     328 +     {"a blank import", "package core\n\nimport (\n\t_ \
  ┃                                                                           "github.com/lthiagol/maitre/internal/brand\"\n)\n"},
  ┃                                                                     329 +   } {
  ┃                                                                     330 +     dir := t.TempDir()
  ┃                                                                     331 +     write(t, filepath.Join(dir, "zz_probe.go"), c.src)
  ┃                                                                     332 +     if got := importsIn(t, dir); len(got) == 0 {
  ┃                                                                     333 +       t.Errorf("the wire should catch %s", c.name)
  ┃                                                                     334 +     }
  ┃                                                                     335 +   }
  ┃                                                                     336 + }
  ┃                                                                     337 +
  ┃                                                                     338 + // importsIn are the packages of this repository that a
  ┃                                                                           directory's shipped code
  ┃                                                                     339 + // imports.
  ┃                                                                     340 + //
  ┃                                                                     341 + // **The imports are parsed, not read.** A line that
  ┃                                                                           starts with the module path is
  ┃                                                                     342 + // one spelling of an import; a grouped import and a
  ┃                                                                           blank import are other
  ┃                                                                     343 + // spellings, and a wire that reads lines waits for the
  ┃                                                                           other spellings to arrive.
  ┃                                                                     344 + //
  ┃    314   // The tests are left out: the rule is about the code that 345   // The tests are left out: the rule is about the code
  ┃          ships.                                                           that ships.
  ┃    315 - func oursIn(t *testing.T, dir string) []string {           346 + func importsIn(t *testing.T, dir string) []string {
  ┃    316     t.Helper()                                               347     t.Helper()
  ┃    317     var out []string                                         348     var out []string
  ┃    318     err := filepath.WalkDir(dir, func(path string, d fs.     349     err := filepath.WalkDir(dir, func(path string, d fs.
  ┃          DirEntry, err error) error {                                     DirEntry, err error) error {
  ┃    319       if err != nil {                                        350       if err != nil {
  ┃    @@ -321,19 +352,18 @@
  ┃    321       }                                                      352       }
  ┃    322       if d.IsDir() || !strings.HasSuffix(path, ".go") ||     353       if d.IsDir() || !strings.HasSuffix(path, ".go") ||
  ┃          strings.HasSuffix(path, "_test.go") {                            strings.HasSuffix(path, "_test.go") {
  ┃    323         return nil                                           354         return nil
  ┃    324       }                                                      355       }
  ┃    325 -     raw, err := os.ReadFile(path)                          356 +     fset := token.NewFileSet()
  ┃                                                                     357 +     file, err := parser.ParseFile(fset, path, nil,
  ┃                                                                           parser.ImportsOnly)
  ┃    326       if err != nil {                                        358       if err != nil {
  ┃    327 -       return err                                           359 +       return fmt.Errorf("parse %s: %w", path, err)
  ┃    328       }                                                      360       }
  ┃    329 -     for _, line := range strings.Split(string(raw), "\n")  361 +     for _, imp := range file.Imports {
  ┃          {
  ┃    330 -       line = strings.TrimSpace(line)                       362 +       pkg := strings.Trim(imp.Path.Value, `"`)
  ┃    331 -       if !strings.HasPrefix(line, `"github.com/lthiagol/   363 +       if strings.HasPrefix(pkg, modulePath) {
  ┃          maitre/`) {
  ┃    332 -         continue                                           364 +         out = append(out, pkg)
  ┃    333         }                                                    365         }
  ┃    334 -       pkg := strings.Trim(line, `"`)
  ┃    335 -       out = append(out, pkg)
  ┃    336       }                                                      366       }
  ┃    337       return nil                                             367       return nil
  ┃    338     })                                                       368     })
  ┃    339     if err != nil {                                          369     if err != nil {
  ┃

  ┃
  ┃  ← Edit internal/brand/mark_test.go
  ┃
  ┃    17                                                               17
  ┃    18   // repoRoot is the repository, two levels up from this      18   // repoRoot is the repository, two levels up from this
  ┃         package.                                                         package.
  ┃    19   const repoRoot = "../.."                                    19   const repoRoot = "../.."
  ┃    20                                                               20
  ┃                                                                     21 + // modulePath is the module of this repository.
  ┃                                                                     22 + const modulePath = "github.com/lthiagol/maitre/"
  ┃                                                                     23 +
  ┃    21   // TestTheDrawingsMatchTheSource holds the code to the art. 24   // TestTheDrawingsMatchTheSource holds the code to the
  ┃                                                                          art.
  ┃    22   //                                                          25   //
  ┃    23   // **The comparison runs from the file to the code.** The   26   // **The comparison runs from the file to the code.** The
  ┃         drawing in                                                       drawing in
  ┃    24   // `images/logo/ascii/` is the source — the code holds it,  27   // `images/logo/ascii/` is the source — the code holds it,
  ┃         so the file is what the                                          so the file is what the
  ┃
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && gofmt -l . ; go test ./internal/brand -v 2>&1 | grep -E '^(=== RUN|--- |ok|FAIL|    )' |
  ┃   head -30
  ┃
  ┃  (11 earlier lines)
  ┃  --- PASS: TestTheArtScopeWireCanFail (0.01s)
  ┃  === RUN   TestTheMarkIsNotTheDomain
  ┃  --- PASS: TestTheMarkIsNotTheDomain (0.00s)
  ┃  === RUN   TestTheBoundaryWireCanFail
  ┃  --- PASS: TestTheBoundaryWireCanFail (0.00s)
  ┃  ok    github.com/lthiagol/maitre/internal/brand  0.260s
  ┃

     + Thought · 7.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cat > internal/core/zz_probe.go <<'EOF'
  ┃  package core…
  ┃
  ┃  BUILD_OK
  ┃  --- FAIL: TestTheMarkIsNotTheDomain (0.00s)
  ┃      mark_test.go:315: the domain should not import the mark, but it imports [github.com/lthiagol/maitre/internal/brand]
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/brand  0.193s
  ┃  FAIL
  ┃   M internal/brand/mark_test.go
  ┃

     + Thought · 272ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cat > internal/core/zz_probe.go <<'EOF'
  ┃  package core…
  ┃
  ┃  BUILD_OK
  ┃  --- FAIL: TestTheMarkIsNotTheDomain (0.00s)
  ┃      mark_test.go:315: the domain should not import the mark, but it imports [github.com/lthiagol/maitre/internal/brand]
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/brand  0.205s
  ┃   M internal/brand/mark_test.go
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/brand/mark_test.go && git commit -q -F - <<'EOF'
  ┃  fix: the boundary wire parses the imports, whatever the spelling…
  ┃
  ┃  659e14a fix: the boundary wire parses the imports, whatever the spelling
  ┃  f39af35 fix: the wire that says where the mark may be, and that can fail
  ┃  a919d7f fix: the mark yields to the content, in the height too
  ┃

     + Thought · 4.8s

     Now the boundary pinning:

  ┃
  ┃  ← Edit internal/ui/mark_test.go
  ┃
  ┃    36   //                                                          36   //
  ┃    37   // A drawing that does not fit wraps, and a wrapped row     37   // A drawing that does not fit wraps, and a wrapped row
  ┃         breaks the layout. So the                                        breaks the layout. So the
  ┃    38   // column chooses the form: the full form when it fits, the 38   // column chooses the form: the full form when it fits,
  ┃         half form when only it                                           the half form when only it
  ┃    39   // fits, and no art when neither does.                      39   // fits, and no art when neither does.
  ┃                                                                     40 + //
  ┃                                                                     41 + // **Both boundaries are pinned.** The full form is 28
  ┃                                                                          columns wide and is drawn one
  ┃                                                                     42 + // column in, so it needs 29: a column of exactly 29 takes
  ┃                                                                          the full form and a
  ┃                                                                     43 + // column of 28 does not. The half form is 14 and needs
  ┃                                                                          15: 15 takes the half form
  ┃                                                                     44 + // and 14 and 13 do not.
  ┃    40   func TestTheSousChefDrawsTheMarkThatFits(t *testing.T) {    45   func TestTheSousChefDrawsTheMarkThatFits(t *testing.T) {
  ┃    41 -   full := brand.Width(brand.Full)
  ┃    42 -   half := brand.Width(brand.Half)
  ┃    43 -
  ┃    44     for _, c := range []struct {                              46     for _, c := range []struct {
  ┃    45       name   string                                           47       name   string
  ┃    46       column int                                              48       column int
  ┃    47       want   []string                                         49       want   []string
  ┃    48     }{                                                        50     }{
  ┃    49 -     {"the column fits the full form", full + 1, brand.Full} 51 +     {"a column of exactly the full form's width", 29,
  ┃         ,                                                                brand.Full},
  ┃    50 -     {"the column fits the half form only", half + 1, brand. 52 +     {"a column one short of the full form", 28, brand.
  ┃         Half},                                                           Half},
  ┃    51 -     {"the column fits the full form at the wide end", 200,  53 +     {"a column of exactly the half form's width", 15,
  ┃         brand.Full},                                                     brand.Half},
  ┃    52 -     {"the column fits neither", half, nil},                 54 +     {"a column one short of the half form", 14, nil},
  ┃    53 -     {"a very narrow column", 4, nil},                       55 +     {"a column of thirteen", 13, nil},
  ┃                                                                     56 +     {"a wide column", 200, brand.Full},
  ┃    54     } {                                                       57     } {
  ┃    55       got := markFor(c.column)                                58       got := markFor(c.column)
  ┃    56       if len(got) != len(c.want) {                            59       if len(got) != len(c.want) {
  ┃    57         t.Errorf("%s: the mark has %d lines, want %d", c.     60         t.Errorf("%s: the mark has %d lines, want %d", c.
  ┃         name, len(got), len(c.want))                                     name, len(got), len(c.want))
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && gofmt -l . ; go test ./internal/ui -run 'SousChef|Mark' -v 2>&1 | grep -E '^(--- |ok|FAI
  ┃  L)'
  ┃
  ┃  (6 earlier lines)
  ┃  --- PASS: TestTheSousChefTabIsThere (0.02s)
  ┃  --- PASS: TestTheSousChefShowsTheVitalSigns (0.02s)
  ┃  --- PASS: TestTheSousChefSaysWhenTheVersionsDisagree (0.02s)
  ┃  --- PASS: TestTheSousChefShowsTheLastEvent (0.02s)
  ┃  --- PASS: TestTheSousChefSaysWhenThereAreNoEvents (0.01s)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  8.339s
  ┃

     + Thought · 478ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -i.bak 's/case column >= brand.Width(brand.Full)+1:/case column >= brand.Width(brand
  ┃  .Full):/' internal/ui/views.go && grep -n "case column" internal/ui/views.go && go test ./internal/ui -run 'SousChef|Mark' …
  ┃
  ┃  (5 earlier lines)
  ┃      mark_test.go:65: a column one short of the full form: line 3 is "⢸⣿⣿⣿⣿⣿⣿⣿⠶⠀⢸⣿⣿⣿⣿⣿⣿⡇⠀⠶⣿⣿⣿⣿⣿⣿⣿⡇", want "⢸⣿⣿⣿⠶⢸⣿⣿⣿⠀⣿⣿⣿⣿"
  ┃      mark_test.go:65: a column one short of the full form: line 4 is "⢸⣿⣿⣿⣿⣏⣠⣤⣴⡆⠸⣿⣿⣿⣿⣿⣿⠇⢰⣦⣤⣄⣹⣿⣿⣿⣿⡇", want "⢸⣿⣿⣠⣴⠸⣿⣿⣿⢰⣤⣹⣿⣿"
  ┃      mark_test.go:65: a column one short of the full form: line 5 is "⠀⠀⢻⣿⣿⠿⠿⠟⠛⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠛⠻⠿⠿⣿⣿⡟", want "⠀⢻⣿⠿⠛⠁⠀⠀⠀⠀⠉⠻⠿⣿"
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/ui  8.464s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  8.208s
  ┃

     + Thought · 995ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/ui/mark_test.go && git commit -q -F - <<'EOF'
  ┃  test: both boundaries of the width rule are pinned…
  ┃
  ┃  a6065b3 test: both boundaries of the width rule are pinned
  ┃  659e14a fix: the boundary wire parses the imports, whatever the spelling
  ┃  f39af35 fix: the wire that says where the mark may be, and that can fail
  ┃  a919d7f fix: the mark yields to the content, in the height too
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make lint && echo "gofmt: [$(gofmt -l .)]" && make test 2>&1 | tail -8 && make regen-gol
  ┃  dens && echo "status: [$(git status --short)]" && for p in internal/core internal/cli internal/ui internal/brand; do echo "…
  ┃
  ┃  (11 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  4.853s
  ┃  status: []
  ┃  internal/core: 0 external modules
  ┃  internal/cli: 0 external modules
  ┃  internal/ui: 28 external modules
  ┃  internal/brand: 0 external modules
  ┃

     + Thought · 4.7s

  ┃
  ┃  # Wrote work/inbox/runner-slice-g.md
  ┃
  ┃    1 Runner: slice G fix 1 done. Commits: 4. Fixed: F-94, F-95, F-96 and the boundary. The tree is free of the drawings; the glyph
  ┃      is in the brand package, in the rendered output and in the card's art. Notes: the reviewer should know these.
  ┃    2
  ┃    3 | Commit | Fix |
  ┃    4 |---|---|
  ┃    5 | `a919d7f` fix: the mark yields to the content, in the height too | F-94 |
  ┃    6 | `f39af35` fix: the wire that says where the mark may be, and that can fail | F-95 |
  ┃    7 | `659e14a` fix: the boundary wire parses the imports, whatever the spelling | F-96 |
  ┃    8 | `a6065b3` test: both boundaries of the width rule are pinned | the boundary |
  ┃    9
  ┃   10 **F-94 — the height rule.** The tab's area is `height − 5` rows; the mark costs
  ┃   11 **six** (five lines and the blank under them); the content is 15. So the mark is
  ┃   12 drawn only when the content and the mark both fit, and otherwise the tab has **no
  ┃   13 mark** — the content is whole. The wire renders the tab at 80 columns and at 20,
  ┃   14 24, 26, 30 and 40 rows, and asserts that the content is whole at every one of them
  ┃   15 (the `LAST EVENT` section and the event row are on the screen) and that the mark is
  ┃   16 drawn only at 26 and above. At 80×24 the tab is byte-for-byte what it was before
  ┃   17 the mark.
  ┃   18
  ┃   19 **F-95 — the scope wire, and the truth about the tree.** Two rules, because there
  ┃   20 are two things:
  ┃   21
  ┃   22 | Thing | Rule | Where it is, today |
  ┃   23 |---|---|---|
  ┃   24 | the drawings | in the source files and in the declaration, and **nowhere else** — any file type | the two `.txt` sources
  ┃      and `internal/brand/mark.go` |
  ┃   25 | the compact glyph | must not be **pasted into another code file** | `internal/brand/mark.go`, the two goldens (rendered
  ┃      output), the wallpaper SVG (ninety copies — the card's own art), and prose: the brand sheet, the ledger, the briefs |
  ┃   26
  ┃   27 **The earlier claim was false and is withdrawn**: *"the mark's characters appear in
  ┃   28 no other file in the tree"* — the goldens carry the glyph on the header line, and
  ┃   29 so do the wallpaper, the brand sheet and the ledger. What is true: **the drawings
  ┃   30 appear nowhere else in the tree**; the glyph is declared in the brand package and
  ┃   31 appears in the rendered output and in the card's art.
  ┃   32
  ┃   33 **The exemption is per declaration, not per package.** The declaration is found by
  ┃   34 parsing the package, so a pasted copy is never mistaken for it and a rename does
  ┃   35 not widen the exemption. The wire can fail: the reviewer's two mutations (a copy
  ┃   36 in `internal/brand/zz_probe.go`, a copy in `internal/ui/testdata/zz_probe.txt`), the
  ┃   37 glyph in another code file, and each of them fails it — proven in a scratch tree
  ┃   38 and by hand in this tree.
  ┃   39
  ┃   40 **The directories the wire does not read**: `.git`, `bin`, `skills` and `work`. The
  ┃   41 binaries carry the art as string literals, and the transcripts capture what the
  ┃   42 console shows — art and all. Everything else is read: `docs`, `images`, `kitchen`,
  ┃   43 every file, any type. If the reviewer wants `work/` scanned too, that is one line.
  ┃   44
  ┃   45 **F-96 — the boundary wire parses the imports.** Every spelling is caught —
  ┃   46 single-line, grouped, blank — because a wire that catches one spelling waits for
  ┃   47 the other. Proven by hand in this tree: a blank import and a single-line import of
  ┃   48 `internal/brand` in `internal/core` both fail the wire, and both compile. The
  ┃   49 scope is in the comment: the rule is about the code that **ships**, so the wire
  ┃   50 reads the **non-test** files of `internal/core`.
  ┃   51
  ┃   52 **The boundary is pinned.** The full form is 28 columns and is drawn one column in,
  ┃   53 so it needs **29**: 29 takes the full form, 28 does not. The half form needs
  ┃   54 **15**: 15 takes it, 14 and 13 do not. Changing the threshold to
  ┃   55 `brand.Width(Full)` now fails the suite — it passed before.
  ┃   56
  ┃   57 **The tree.** `make lint && gofmt -l . && make test && make regen-goldens &&
  ┃   58 git status --short` is clean; the goldens did not move (the header did not move).
  ┃   59 The budgets: `internal/core` 0, `internal/cli` 0, `internal/ui` 28,
  ┃   60 `internal/brand` 0. The tripwires hold — the derived-views table, the gate table,
  ┃   61 the plan format, the command count (42). `core.DomainVersion` is still 7.
  ┃   62
  ┃   63 **Not touched**: `docs/audit.md`, `docs/07-vocabulary.md`, `work/` (except this
  ┃   64 receipt), `images/`, `README.md`, and the card's copy.
  ┃   65
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/inbox && git commit -q -m "inbox: runner reports slice-g fix 1" && git log
  ┃  --oneline -6 && git status --short
  ┃
  ┃  de08e67 inbox: runner reports slice-g fix 1
  ┃  a6065b3 test: both boundaries of the width rule are pinned
  ┃  659e14a fix: the boundary wire parses the imports, whatever the spelling
  ┃  f39af35 fix: the wire that says where the mark may be, and that can fail
  ┃  a919d7f fix: the mark yields to the content, in the height too
  ┃  f00acdf work: slice G's fix cycle, and the brand sheet's width
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1Z:p1" "Runner: slice G fix 1 done. Commits: 4 (plus the receipt).
  ┃  Fixed: F-94, F-95, F-96 and the boundary. The tree is free of the drawings; the glyph is in the brand package, in the rende…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:
  ┃  opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"done","completion_seq":584,"cwd":"/Users/thiago/Code/Pessoal/
  ┃  maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,
  ┃  "screen_detection_skipped":true,"state_change_seq":584,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722","terminal_title":"OC
  ┃  | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":"w1Z"},"type":"agent_prompted"}}
  ┃

     The fix cycle is done, committed, receipted and nudged.

     ┌──────────────────────────────────────────────────────────────────────────────────────────────┬─────────────────────────────────┐
     │ Commit                                                                                       │ Fix                             │
     ├──────────────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────┤
     │ a919d7f fix: the mark yields to the content, in the height too                               │ F-94                            │
     ├──────────────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────┤
     │ f39af35 fix: the wire that says where the mark may be, and that can fail                     │ F-95                            │
     ├──────────────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────┤
     │ 659e14a fix: the boundary wire parses the imports, whatever the spelling                     │ F-96                            │
     ├──────────────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────┤
     │ a6065b3 test: both boundaries of the width rule are pinned                                   │ the boundary                    │
     ├──────────────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────┤
     │ de08e67 inbox: runner reports slice-g fix 1                                                  │ receipt                         │
     └──────────────────────────────────────────────────────────────────────────────────────────────┴─────────────────────────────────┘

     F-94 — the height rule: the tab's area is height − 5, the mark costs 6 rows, the content is 15, so the mark is drawn only when
     both fit. At 80×24 the tab is what it was before the mark — LAST EVENT and the event row are both there. The wire renders at 20/
     24/26/30/40 rows and asserts the content is whole at every one, the mark only from 26 up.

     F-95 — two rules now: the drawings may live in the source files and the declaration and nowhere else, any file type; the glyph
     must not be pasted into another code file (it legitimately rides the goldens, the wallpaper and the prose). The exemption is per
     declaration, found by parsing. The reviewer's two mutations, the glyph probe and the non-Go copy each fail it — proven in a
     scratch tree and by hand in the real tree.

     F-96 — the imports are parsed, not read: single-line, grouped and blank spellings are all caught, and a real import of internal/
     brand in internal/core fails the wire (verified by hand, both spellings compile).

     The boundary — 29 takes the full form, 28 does not; 15 takes the half, 14 and 13 do not. Changing the threshold to brand.Width(
     Full) now fails the suite; it passed before.

     One correction I made to my own earlier claim: "the mark's characters appear in no other file in the tree" was false — the
     goldens, the wallpaper, the brand sheet and the ledger all carry the glyph. What is true is that the drawings appear nowhere else,
     and the glyph is not pasted into another code file. That is stated in the receipt and in the nudge.

     The gate is green: make lint, gofmt empty, every package ok, the goldens unmoved, the tree clean. Budgets unchanged (core 0, cli
     0, ui 28, brand 0); DomainVersion still 7.

     Build · MiniMax-M3.1-Flash · 8m 54s · 129.6 tok/s

  ┃
  ┃  You are the **runner** — the same lane that wrote slice G's fix. The orchestrator's pane
  ┃  id is `w1Z:p1`; use it verbatim in the notify command below.
  ┃
  ┃  ## The last fix for slice G — one finding and one note
  ┃
  ┃  The review came back with **F-97** and a LOW note that belongs with it. The
  ┃  orchestrator ruled: **both are fixed here.** One commit.
  ┃
  ┃  ### F-97 — the wire exempts the declaration **file**, not the declaration
  ┃
  ┃  > The one that still passes: a copy pasted into `internal/brand/mark.go` **itself**,
  ┃  > next to the declaration — the wire exempts the declaration FILE, not the
  ┃  > declaration, so the whole suite stays green.
  ┃
  ┃  **The rule is *one declaration*, not *one file*.** A second copy of the drawing in the
  ┃  file that declares it is a **second declaration**, and it is the same defect the wire
  ┃  exists to catch — just one file closer to home.
  ┃
  ┃  **The fix — the exemption is the declaration:**
  ┃
  ┃  | | |
  ┃  |---|---|
  ┃  | **The drawing** | may appear **only in the canonical declaration** — the one `var` that holds it, in the one file that holds it.
  ┃  Anywhere else in the tree, in **any** file type, is drift — **including in the declaring file itself**. |
  ┃  | **Inside the declaring file** | the drawing appears **exactly once**, in the declaration. |
  ┃
  ┃  So the wire's question becomes: *is every occurrence of the drawing the canonical
  ┃  declaration?* — and a second copy beside it is not.
  ┃
  ┃  **Prove it:** paste a copy into `internal/brand/mark.go` next to the declaration and
  ┃  confirm the suite **fails**. Then paste one into every other file the wire reads and
  ┃  confirm it fails there too (the reviewer did: a golden, a doc, the plan, `cmd/`,
  ┃  `images/`).
  ┃
  ┃  ### The note that rides with it — the glyph's scope
  ┃
  ┃  > the glyph rule is scoped to `.go` files, so the glyph pasted into
  ┃  > `scripts/console-shot.py` passes — and the runner's own table says the glyph must not
  ┃  > be "pasted into another code file", and a `.py` file **is** a code file.
  ┃
  ┃  **The rule as stated and the rule as implemented must be the same rule.** The glyph's
  ┃  scope is **code files** — `.go` *and* `.py` — and not "Go files". The glyph may be:
  ┃
  ┃  - **declared** in the brand package;
  ┃  - **rendered** — the goldens carry it on the header line, and that is output;
  ┃  - **drawn** — the wallpaper's ninety copies are the card's own art.
  ┃
  ┃  and it may **not** be pasted into another **code** file.
  ┃
  ┃  **So the wire needs a *category*, not an extension.** Three categories — **code**,
  ┃  **rendered output**, and **art** — and each form's rule says where it may live:
  ┃
  ┃  | | may be pasted into |
  ┃  |---|---|
  ┃  | **the drawings** | **nothing** but the canonical declaration |
  ┃  | **the glyph** | the canonical declaration, and nowhere else in **code** — but the rendered output and the art may carry it |
  ┃
  ┃  **And the exclusions stand.** The reviewer agreed: `.git` is the repository's objects,
  ┃  `bin` is build output, `skills` is vendored third-party, and `work` is the night's
  ┃  scratch — **the transcripts quote what the console shows, art and all**, so `work`
  ┃  *must* be excluded. Say that in the wire's comment, so the next reader knows each
  ┃  exclusion is a *record* and not a hole.
  ┃
  ┃  ## The invariants
  ┃
  ┃  - The domain layer has **zero UI imports** — and the wire proves it, whatever the
  ┃    spelling.
  ┃  - One fact, one answer: the art is declared once, in one package, and every consumer
  ┃    reads it from there.
  ┃  - `core.DomainVersion` is **7**. Do not move it.
  ┃  - The gate table, the plan format and the command count (**42**) must not move.
  ┃
  ┃  ## The tree is shared
  ┃
  ┃  Stage **by name** — never `git add -A`, never `git add .`. Leave anything you did not
  ┃  write alone. If a file you did not write has changed, **stop and report it**.
  ┃
  ┃  ## Not yours — do not touch
  ┃
  ┃  | | |
  ┃  |---|---|
  ┃  | `docs/audit.md` | the reviewer's to write. It files F-97 and closes it. |
  ┃  | `docs/07-vocabulary.md`, `work/`, `images/`, `README.md` | the orchestrator's. |
  ┃  | The card's copy | it stays in the card. |
  ┃
  ┃  ## The gate
  ┃
  ┃  ```bash
  ┃  make lint && gofmt -l . && make test && make regen-goldens && git status --short
  ┃  ```
  ┃
  ┃  Green.
  ┃
  ┃  ## One commit
  ┃
  ┃  ```
  ┃  Per-AC evidence:
  ┃  - F-97: a copy of the drawing beside the declaration fails the wire; the glyph's scope covers the code files, the output and the
  ┃  art; go test ./internal/brand -v exit 0
  ┃  ```
  ┃
  ┃  **One write per file, and commit as soon as the file is coherent.**
  ┃
  ┃  ## Finish with the receipt — update the one you wrote when you started
  ┃
  ┃  ```bash
  ┃  cat > work/inbox/runner-slice-g.md <<EOF
  ┃  Runner: slice G fix 2 done. Commits: <n>. Fixed: F-97 and the glyph's scope. The drawing is declared once, and the wire holds it.
  ┃  Notes: <anything the reviewer should know>.
  ┃  EOF
  ┃  git add work/inbox && git commit -q -m "inbox: runner reports slice-g fix 2"
  ┃  herdr agent prompt "w1Z:p1" "Runner: slice G fix 2 done. Commits: <n>. Fixed: F-97 and the glyph's scope. The drawing is declared
  ┃  once, and the wire holds it. Notes: <anything the reviewer should know>."
  ┃  ```
  ┃
  ┃  **Commit the receipt before you nudge.**
  ┃

     + Thought · 687ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && ls -t work/inbox/ | head -5 && echo "=====" && cat work/inbox/reviewer-slice-g-fix-1.md
  ┃  2>/dev/null || cat "$(ls -t work/inbox/reviewer* | head -1)"
  ┃
  ┃  (1 earlier line)
  ┃  s. The CLI's help carries the full form above the title line, and maitre --render --all shows the same tab. LOW notes: (a) the
  ┃  four excluded directories are all records - .git is the repository's own objects, bin is build output (the binaries carry the art
  ┃  as string literals), skills is vendored third-party, work is the night's scratch (the briefs, transcripts and receipts quote what
  ┃  the console shows, art and all) - so excluding them is right, and the runner's own note that scanning work/ would be "one line" is
  ┃  a fair offer; (b) the glyph rule is scoped to .go files, so the glyph pasted into scripts/console-shot.py passes - the glyph is
  ┃  one character and the rule is about code, so this is a note rather than a blocker, but the runner's table says the glyph must not
  ┃  be "pasted into another code file" and a .py file is a code file; (c) the brand sheet's "28 columns wide at its widest - the first
  ┃  row is 26" is now correct, and the drawing's first row is indeed 26 - the wire's Width() measures the widest line, so the
  ┃  console's budget of 29 is right.
  ┃

     + Thought · 1m 2s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
  ┃  import os, collections…
  ┃
  ┃  (20 earlier lines)
  ┃  kitchen/backlog:
  ┃
  ┃  kitchen/courses:
  ┃
  ┃  kitchen/recipes:
  ┃  RCP01-the-console-shows-a-finding-s-conversation.json
  ┃

     + Thought · 21.9s

     → Explored: 1 read

     + Thought · 18.3s

     → Explored: 1 read

     + Thought · 3.3s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cat > /tmp/newscope.go.txt <<'BLOCK'
  ┃  // TestTheArtIsDeclaredOnlyHere holds the mark to one place.…
  ┃
  ┃  # github.com/lthiagol/maitre/internal/brand
  ┃  # [github.com/lthiagol/maitre/internal/brand]
  ┃  vet: internal/brand/mark_test.go:310:6: canonicalDeclaration redeclared in this block
  ┃
  ┃  Exited with code 1
  ┃

     + Thought · 225ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
  ┃  p='internal/brand/mark_test.go'…
  ┃
  ┃  (5 earlier lines)
  ┃  --- PASS: TestTheArtIsDeclaredOnlyHere (0.02s)
  ┃  --- PASS: TestTheArtScopeWireCanFail (0.02s)
  ┃  --- PASS: TestTheDrawingBesideTheDeclarationFails (0.00s)
  ┃  --- PASS: TestTheMarkIsNotTheDomain (0.00s)
  ┃  --- PASS: TestTheBoundaryWireCanFail (0.00s)
  ┃  ok    github.com/lthiagol/maitre/internal/brand  0.380s
  ┃

     + Thought · 1.1s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/brand/mark.go /tmp/mark.go.bak && python3 -c "
  ┃  art = open('images/logo/ascii/bow-tie-full.txt', encoding='utf-8').read().rstrip('\n').split('\n')[1]…
  ┃
  ┃  BUILD_OK
  ┃  --- FAIL: TestTheArtIsDeclaredOnlyHere (0.01s)
  ┃      mark_test.go:132: internal/brand/mark.go: carries a line of the full drawing — the drawing is declared once, in internal/brand
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/brand  0.263s
  ┃  FAIL
  ┃   M internal/brand/mark_test.go
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
  ┃  import subprocess, os…
  ┃
  ┃  (8 earlier lines)
  ┃  images/logo/zz_probe.svg: wire fails = True
  ┃      ['    mark_test.go:132: images/logo/zz_probe.svg: carries a line of the full drawing — the drawing is declared once, in
  ┃  internal/brand']
  ┃  scripts/console-shot.py: wire fails = True
  ┃      ['    mark_test.go:132: scripts/console-shot.py: carries a line of the full drawing — the drawing is declared once, in
  ┃  internal/brand']
  ┃   M internal/brand/mark_test.go
  ┃   D scripts/console-shot.py
  ┃

     + Thought · 1.5s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git checkout scripts/console-shot.py && git status --short && ls -l scripts/console-shot
  ┃  .py && head -3 scripts/console-shot.py
  ┃
  ┃  Updated 1 path from the index
  ┃   M internal/brand/mark_test.go
  ┃  -rwxr-xr-x@ 1 thiago  staff  3392 Oct 28 13:12 scripts/console-shot.py
  ┃  #!/usr/bin/env python3
  ┃  """Capture what the console actually shows, as a terminal would draw it.
  ┃

     + Thought · 892ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 -c "
  ┃  open('scripts/zz_probe.py','w',encoding='utf-8').write('probe = \"\u22c8\"\n')…
  ┃
  ┃  --- FAIL: TestTheArtIsDeclaredOnlyHere (0.02s)
  ┃      mark_test.go:132: scripts/zz_probe.py: pastes the compact glyph into a code file — it is declared in internal/brand
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/brand  0.235s
  ┃   M internal/brand/mark_test.go
  ┃

     + Thought · 1.8s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make lint && echo "gofmt: [$(gofmt -l .)]" && make test 2>&1 | tail -8 && make regen-gol
  ┃  dens && echo "status: [$(git status --short)]"
  ┃
  ┃  (7 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/core  21.756s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  78.405s
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃  go test ./internal/ui -run Golden -update
  ┃  ok    github.com/lthiagol/maitre/internal/ui  5.336s
  ┃  status: [ M internal/brand/mark_test.go]
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/brand/mark_test.go && git commit -q -F - <<'EOF'
  ┃  fix: the wire's exemption is the declaration, not the file…
  ┃
  ┃  b5dda82 fix: the wire's exemption is the declaration, not the file
  ┃  e80e760 work: slice G's last fix — the declaration, not the file
  ┃  7c7aca0 docs: slice G cycle 2 in the ledger - three closed, one filed
  ┃

     + Thought · 1.1s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -10 && echo "=== files in the last 3 ===" && git show --stat --oneline
  ┃   e80e760 | head -20 && git show --stat --oneline 7c7aca0 | head -20
  ┃
  ┃  (11 earlier lines)
  ┃  e80e760 work: slice G's last fix — the declaration, not the file
  ┃   work/briefs/fix-g2.md | 112 ++++++++++++++++++++++++++++++++++++++++++++++++++
  ┃   1 file changed, 112 insertions(+)
  ┃  7c7aca0 docs: slice G cycle 2 in the ledger - three closed, one filed
  ┃   docs/audit.md | 10 +++++-----
  ┃   1 file changed, 5 insertions(+), 5 deletions(-)
  ┃

     + Thought · 11.4s

  ┃
  ┃  # Wrote work/inbox/runner-slice-g.md
  ┃
  ┃    1 Runner: slice G fix 2 done. Commits: 1 (plus this receipt). Fixed: F-97 and the glyph's scope. The drawing is declared once,
  ┃      and the wire holds it. Notes: the reviewer should know these.
  ┃    2
  ┃    3 **F-97 — the exemption is the declaration, not the file.** The wire's question is
  ┃    4 now *is every occurrence of the drawing the canonical declaration?* The declaration
  ┃    5 is found by parsing the package, and the two literals that hold the drawings are
  ┃    6 masked out of the search — so the full drawing may be in the literal that holds it,
  ┃    7 the half drawing in the literal that holds it, and nowhere else: not in another
  ┃    8 file, not in another corner of the declaring file.
  ┃    9
  ┃   10 **The glyph's scope is the category, not the extension.** The rule as stated and the
  ┃   11 rule as implemented are the same rule now:
  ┃   12
  ┃   13 | Form | code | rendered output | art | prose |
  ┃   14 |---|---|---|---|---|
  ┃   15 | the drawings | the canonical declaration, and nowhere else | — | — | — |
  ┃   16 | the compact glyph | the canonical declaration, and nowhere else | may ride along | may ride along | may be quoted |
  ┃   17
  ┃   18 A script is code, so the glyph pasted into `scripts/console-shot.py` now fails the
  ┃   19 wire. Prose may quote the glyph — the brand sheet and the reviewer's ledger name
  ┃   20 it, and quoting is not pasting.
  ┃   21
  ┃   22 **The exclusions are records, not holes**, and the comment says so: `.git` is the
  ┃   23 repository's own objects, `bin` is build output (the binaries carry the art as
  ┃   24 string literals), `skills` is vendored third-party, and `work` is the night's
  ┃   25 scratch — the briefs, the transcripts and the receipts quote what the console
  ┃   26 shows, art and all.
  ┃   27
  ┃   28 **Proven by hand in this tree**, each one failing the wire and then removed: a copy
  ┃   29 beside the declaration in `internal/brand/mark.go`; a copy in a golden; a copy in a
  ┃   30 doc; a copy in the plan (`kitchen/`); a copy in `cmd/`; a copy in `images/`; and
  ┃   31 the glyph in a script. The wire's own test proves the same list in a scratch tree.
  ┃   32
  ┃   33 **One thing to report.** While proving the wire, a probe was written over
  ┃   34 `scripts/console-shot.py` by mistake — the probe's path was the script's own path.
  ┃   35 It was restored from git at once and is byte-identical to HEAD; `git status` shows
  ┃   36 it untouched. No other file was written outside the ones I own.
  ┃   37
  ┃   38 **The tree.** `make lint && gofmt -l . && make test && make regen-goldens &&
  ┃   39 git status --short` is clean; the goldens did not move. The budgets: `internal/core`
  ┃   40 0, `internal/cli` 0, `internal/ui` 28, `internal/brand` 0. The tripwires hold — the
  ┃   41 derived-views table, the gate table, the plan format, the command count (42).
  ┃   42 `core.DomainVersion` is still 7.
  ┃   43
  ┃   44 **Not touched**: `docs/audit.md`, `docs/07-vocabulary.md`, `work/` (except this
  ┃   45 receipt), `images/`, `README.md`, and the card's copy.
  ┃   46
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/inbox && git commit -q -m "inbox: runner reports slice-g fix 2" && git log
  ┃  --oneline -3 && git status --short && git diff HEAD --stat
  ┃
  ┃  dee56aa inbox: runner reports slice-g fix 2
  ┃  b5dda82 fix: the wire's exemption is the declaration, not the file
  ┃  e80e760 work: slice G's last fix — the declaration, not the file
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1Z:p1" "Runner: slice G fix 2 done. Commits: 1 (plus the receipt).
  ┃  Fixed: F-97 and the glyph's scope. The drawing is declared once, and the wire holds it. Notes: the wire's exemption is now …
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:
  ┃  opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"done","completion_seq":592,"cwd":"/Users/thiago/Code/Pessoal/
  ┃  maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,
  ┃  "screen_detection_skipped":true,"state_change_seq":592,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722","terminal_title":"OC
  ┃  | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":"w1Z"},"type":"agent_prompted"}}
  ┃

     The last fix is done, committed, receipted and nudged.

     ┌─────────────────────────────────────────────────────────────────────────────────────┬──────────────────────────────────────────┐
     │ Commit                                                                              │ Fix                                      │
     ├─────────────────────────────────────────────────────────────────────────────────────┼──────────────────────────────────────────┤
     │ b5dda82 fix: the wire's exemption is the declaration, not the file                  │ F-97 + the glyph's scope                 │
     ├─────────────────────────────────────────────────────────────────────────────────────┼──────────────────────────────────────────┤
     │ dee56aa inbox: runner reports slice-g fix 2                                         │ receipt                                  │
     └─────────────────────────────────────────────────────────────────────────────────────┴──────────────────────────────────────────┘

     F-97 — the wire's question is now is every occurrence of the drawing the canonical declaration? The declaration is found by
     parsing, and the two literals that hold the drawings are masked out of the search: the full drawing may live in the literal that
     holds it, the half drawing in the literal that holds it, and nowhere else — not in another file, not in another corner of the
     declaring file.

     The glyph's scope is a category, not an extension:

     ┌──────────────────────┬────────────────────────────────────────────────┬───────────────────┬──────────────────┬─────────────────┐
     │ Form                 │ code                                           │ rendered output   │ art              │ prose           │
     ├──────────────────────┼────────────────────────────────────────────────┼───────────────────┼──────────────────┼─────────────────┤
     │ the drawings         │ the canonical declaration, and nowhere else    │ —                 │ —                │ —               │
     ├──────────────────────┼────────────────────────────────────────────────┼───────────────────┼──────────────────┼─────────────────┤
     │ the compact glyph    │ the canonical declaration, and nowhere else    │ may ride along    │ may ride along   │ may be quoted   │
     └──────────────────────┴────────────────────────────────────────────────┴───────────────────┴──────────────────┴─────────────────┘

     A script is code, so the glyph in scripts/console-shot.py now fails. The exclusions (.git, bin, skills, work) are recorded in the
     comment as records, not holes.

     Proven by hand in the real tree — a copy beside the declaration in internal/brand/mark.go, in a golden, in a doc, in the plan, in
     cmd/, in images/, and the glyph in a script: each one failed the wire, and each was removed.

     One thing I must report: while proving the wire, a probe was written over scripts/console-shot.py — the probe's path was the
     script's own path. It was restored from git immediately and is byte-identical to HEAD (git status shows it untouched); no other
     file outside my own was written. That is in the receipt and in the nudge.

     The gate is green: make lint, gofmt empty, every package ok, the goldens unmoved, the tree clean; budgets unchanged (core 0, cli
     0, ui 28, brand 0); DomainVersion still 7.

     Build · MiniMax-M3.1-Flash · 4m 50s · 111.7 tok/s

  ┃
  ┃
  ┃
  ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
  ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
  ~/Code/Pessoal/maitre:master                                                                            173.3K (17%)  ctrl+p commands

```
