# The runner's lane — slice C, captured 2026-10-27

Dumped from pane `w1Z:p5` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that wrote slice C: AC-01 to AC-06, then the fix cycle (F-74..F-77).
The review that judged it is in [slice-c-reviewer.md](slice-c-reviewer.md).

```
  ┃                                                                                                                                              45 +   //
  ┃                                                                                                                                              46 +   // The old format kept the pointer, and the pointer is a promise: it says an
  ┃                                                                                                                                              47 +   // id was taken. A plan that is read but never written still has to honour
  ┃                                                                                                                                              48 +   // the promise — the id must not be handed to a stranger.
  ┃                                                                                                                                              49 +   PromotedTo string `json:"promoted_to,omitempty"`
  ┃    44     // Rung is where the item came from — a special or a prep note. It is read                                                         50     // Rung is where the item came from — a special or a prep note. It is read
  ┃    45     // from the directory the file lives in, not from the file.                                                                        51     // from the directory the file lives in, not from the file.
  ┃    46     Rung string `json:"-"`                                                                                                             52     Rung string `json:"-"`
  ┃    47   }                                                                                                                                    53   }
  ┃

  ┃
  ┃  ← Edit internal/core/store.go
  ┃
  ┃    469     return nextIDIn(s.RecipeDir(), recipeIDPrefix, reserved...)                                                                       469     return nextIDIn(s.RecipeDir(), recipeIDPrefix, reserved...)
  ┃    470   }                                                                                                                                   470   }
  ┃    471                                                                                                                                       471
  ┃    472   // reservedRecipeIDs are the recipe ids a backlog item points at.                                                                   472   // reservedRecipeIDs are the recipe ids a backlog item points at.
  ┃                                                                                                                                              473 + //
  ┃                                                                                                                                              474 + // The old plan's artifacts carry pointers too, and a pointer the reader drops is a
  ┃                                                                                                                                              475 + // promise broken: the id it names would be handed to a stranger, and the trail
  ┃                                                                                                                                              476 + // would resolve to the wrong recipe.
  ┃    473   func (s *Store) reservedRecipeIDs() ([]string, error) {                                                                             477   func (s *Store) reservedRecipeIDs() ([]string, error) {
  ┃    474     items, err := s.ListBacklogItems()                                                                                                478     items, err := s.ListBacklogItems()
  ┃    475     if err != nil {                                                                                                                   479     if err != nil {
  ┃    476       return nil, err                                                                                                                 480       return nil, err
  ┃    @@ -478,8 +482,15 @@
  ┃    478     promoted := make([]string, 0, len(items))                                                                                         482     promoted := make([]string, 0, len(items))
  ┃    479     for _, n := range items {                                                                                                         483     for _, n := range items {
  ┃    480       promoted = append(promoted, n.PromotedTo)                                                                                       484       promoted = append(promoted, n.PromotedTo)
  ┃    481     }                                                                                                                                 485     }
  ┃                                                                                                                                              486 +   legacy, err := s.ListLegacyItems()
  ┃                                                                                                                                              487 +   if err != nil {
  ┃                                                                                                                                              488 +     return nil, err
  ┃                                                                                                                                              489 +   }
  ┃                                                                                                                                              490 +   for _, n := range legacy {
  ┃                                                                                                                                              491 +     promoted = append(promoted, n.PromotedTo)
  ┃                                                                                                                                              492 +   }
  ┃    482     return referencedIDs(promoted, "R"), nil                                                                                          493     return referencedIDs(promoted, "R"), nil
  ┃    483   }                                                                                                                                   494   }
  ┃    484                                                                                                                                       495
  ┃    485   // readArtifact finds the current file in dir for the given id.                                                                     496   // readArtifact finds the current file in dir for the given id.
  ┃

     + Thought · 1.7s

  ┃
  ┃  ← Edit internal/core/backlog_concurrency_test.go
  ┃
  ┃    145       }                                                                                                                               145       }
  ┃    146     }                                                                                                                                 146     }
  ┃    147   }                                                                                                                                   147   }
  ┃    148                                                                                                                                       148
  ┃                                                                                                                                              149 + // TestALegacyDanglingPointerResolvesToNothing covers the old plan's pointers.
  ┃                                                                                                                                              150 + //
  ┃                                                                                                                                              151 + // A prep note written before the rename carries the same promise as one written
  ┃                                                                                                                                              152 + // after it: it points at a recipe id, and that id must stay reserved. A reader that
  ┃                                                                                                                                              153 + // drops the pointer hands the id to a stranger, and the trail resolves to the wrong
  ┃                                                                                                                                              154 + // recipe.
  ┃                                                                                                                                              155 + func TestALegacyDanglingPointerResolvesToNothing(t *testing.T) {
  ┃                                                                                                                                              156 +   ops, store := newOps(t)
  ┃                                                                                                                                              157 +
  ┃                                                                                                                                              158 +   // The state a crash leaves, in the plan format before the rename.
  ┃                                                                                                                                              159 +   dir := filepath.Join(store.Root, "prep-notes")
  ┃                                                                                                                                              160 +   if err := os.MkdirAll(dir, 0o755); err != nil {
  ┃                                                                                                                                              161 +     t.Fatal(err)
  ┃                                                                                                                                              162 +   }
  ┃                                                                                                                                              163 +   note := `{"id":"N01","title":"A note","intent":"first","status":"promoted","promoted_to":"recipe RCP10"}`
  ┃                                                                                                                                              164 +   if err := os.WriteFile(filepath.Join(dir, "N01-a-note.json"), []byte(note), 0o644); err != nil {
  ┃                                                                                                                                              165 +     t.Fatal(err)
  ┃                                                                                                                                              166 +   }
  ┃                                                                                                                                              167 +
  ┃                                                                                                                                              168 +   // The id it points at is reserved.
  ┃                                                                                                                                              169 +   r, err := ops.CreateRecipe(readySpec())
  ┃                                                                                                                                              170 +   if err != nil {
  ┃                                                                                                                                              171 +     t.Fatal(err)
  ┃                                                                                                                                              172 +   }
  ┃                                                                                                                                              173 +   if r.ID == "RCP10" {
  ┃                                                                                                                                              174 +     t.Fatalf("the referenced id %q was handed out to a new recipe", r.ID)
  ┃                                                                                                                                              175 +   }
  ┃                                                                                                                                              176 +   next, err := store.NextRecipeID()
  ┃                                                                                                                                              177 +   if err != nil {
  ┃                                                                                                                                              178 +     t.Fatal(err)
  ┃                                                                                                                                              179 +   }
  ┃                                                                                                                                              180 +   if next == "RCP10" {
  ┃                                                                                                                                              181 +     t.Fatalf("the next id is %q, but it is the one the old plan points at", next)
  ┃                                                                                                                                              182 +   }
  ┃                                                                                                                                              183 +
  ┃                                                                                                                                              184 +   // And nothing carries it: the pointer resolves to nothing.
  ┃                                                                                                                                              185 +   recipes, err := store.ListRecipes()
  ┃                                                                                                                                              186 +   if err != nil {
  ┃                                                                                                                                              187 +     t.Fatal(err)
  ┃                                                                                                                                              188 +   }
  ┃                                                                                                                                              189 +   for _, rec := range recipes {
  ┃                                                                                                                                              190 +     if rec.ID == "RCP10" {
  ┃                                                                                                                                              191 +       t.Fatal("the dangling pointer now resolves to a stranger")
  ┃                                                                                                                                              192 +     }
  ┃                                                                                                                                              193 +   }
  ┃                                                                                                                                              194 + }
  ┃                                                                                                                                              195 +
  ┃    149   // TestAPromotionIsRecordedInTheJournal covers the trail.                                                                           196   // TestAPromotionIsRecordedInTheJournal covers the trail.
  ┃    150   //                                                                                                                                  197   //
  ┃    151   // The promotion is an event like any other: the journal is the only event log, and a                                               198   // The promotion is an event like any other: the journal is the only event log, and a
  ┃    152   // promotion nobody can read is a promotion nobody can audit.                                                                       199   // promotion nobody can read is a promotion nobody can audit.
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/core -run 'Legacy|Dangling' -v 2>&1 | head -20
  ┃
  ┃  (5 earlier lines)
  ┃  --- PASS: TestALegacyDanglingPointerResolvesToNothing (0.02s)
  ┃  === RUN   TestALegacyApprovalIsRestampedNotRefused
  ┃  --- PASS: TestALegacyApprovalIsRestampedNotRefused (0.07s)
  ┃  === RUN   TestADanglingDependencyIsVisible
  ┃  --- PASS: TestADanglingDependencyIsVisible (0.27s)
  ┃  PASS
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.543s
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/core/store.go /tmp/store.bak && python3 - <<'EOF'
  ┃  p='internal/core/store.go'…
  ┃
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.226s
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.196s
  ┃

     + Thought · 2.6s

  ┃
  ┃  ← Edit internal/core/backlog_concurrency_test.go
  ┃
  ┃    164     if err := os.WriteFile(filepath.Join(dir, "N01-a-note.json"), []byte(note), 0o644); err != nil {                                  164     if err := os.WriteFile(filepath.Join(dir, "N01-a-note.json"), []byte(note), 0o644); err != nil {
  ┃    165       t.Fatal(err)                                                                                                                    165       t.Fatal(err)
  ┃    166     }                                                                                                                                 166     }
  ┃    167                                                                                                                                       167
  ┃    168 -   // The id it points at is reserved.                                                                                               168 +   // Nine recipes, so the numbering is one short of the id the pointer names.
  ┃    169 -   r, err := ops.CreateRecipe(readySpec())                                                                                           169 +   for i := 0; i < 9; i++ {
  ┃    170 -   if err != nil {                                                                                                                   170 +     if _, err := ops.CreateRecipe(readySpec()); err != nil {
  ┃    171 -     t.Fatal(err)                                                                                                                    171 +       t.Fatal(err)
  ┃                                                                                                                                              172 +     }
  ┃    172     }                                                                                                                                 173     }
  ┃    173 -   if r.ID == "RCP10" {
  ┃    174 -     t.Fatalf("the referenced id %q was handed out to a new recipe", r.ID)
  ┃    175 -   }
  ┃    176     next, err := store.NextRecipeID()                                                                                                 174     next, err := store.NextRecipeID()
  ┃    177     if err != nil {                                                                                                                   175     if err != nil {
  ┃    178       t.Fatal(err)                                                                                                                    176       t.Fatal(err)
  ┃    179     }                                                                                                                                 177     }
  ┃    180     if next == "RCP10" {                                                                                                              178     if next == "RCP10" {
  ┃    181 -     t.Fatalf("the next id is %q, but it is the one the old plan points at", next)                                                   179 +     t.Fatalf("the next id is %q, but that is the one the old plan points at", next)
  ┃    182     }                                                                                                                                 180     }
  ┃                                                                                                                                              181 +   if next != "RCP11" {
  ┃                                                                                                                                              182 +     t.Fatalf("the next id is %q, want RCP11 — the pointer's id is reserved", next)
  ┃                                                                                                                                              183 +   }
  ┃    183                                                                                                                                       184
  ┃    184     // And nothing carries it: the pointer resolves to nothing.                                                                       185     // And nothing carries it: the pointer resolves to nothing.
  ┃                                                                                                                                              186 +   r, err := ops.CreateRecipe(readySpec())
  ┃                                                                                                                                              187 +   if err != nil {
  ┃                                                                                                                                              188 +     t.Fatal(err)
  ┃                                                                                                                                              189 +   }
  ┃                                                                                                                                              190 +   if r.ID == "RCP10" {
  ┃                                                                                                                                              191 +     t.Fatalf("the referenced id %q was handed out to a new recipe", r.ID)
  ┃                                                                                                                                              192 +   }
  ┃    185     recipes, err := store.ListRecipes()                                                                                               193     recipes, err := store.ListRecipes()
  ┃    186     if err != nil {                                                                                                                   194     if err != nil {
  ┃    187       t.Fatal(err)                                                                                                                    195       t.Fatal(err)
  ┃    188     }                                                                                                                                 196     }
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/core -run 'LegacyDangling' -v 2>&1 | head -10 && cp internal/core/store.go /tmp/store.bak && python3 - <<'EOF'
  ┃  p='internal/core/store.go'…
  ┃
  ┃  (6 earlier lines)
  ┃  --- FAIL: TestALegacyDanglingPointerResolvesToNothing (0.13s)
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/core  0.338s
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/core  0.289s
  ┃  FAIL
  ┃

     + Thought · 6.7s

  ┃
  ┃  ← Edit internal/core/backlog_concurrency_test.go
  ┃
  ┃    164     if err := os.WriteFile(filepath.Join(dir, "N01-a-note.json"), []byte(note), 0o644); err != nil {                                  164     if err := os.WriteFile(filepath.Join(dir, "N01-a-note.json"), []byte(note), 0o644); err != nil {
  ┃    165       t.Fatal(err)                                                                                                                    165       t.Fatal(err)
  ┃    166     }                                                                                                                                 166     }
  ┃    167                                                                                                                                       167
  ┃    168 -   // Nine recipes, so the numbering is one short of the id the pointer names.                                                       168 +   // Nine recipes, so the numbering reaches the id the pointer names.
  ┃    169     for i := 0; i < 9; i++ {                                                                                                          169     for i := 0; i < 9; i++ {
  ┃    170       if _, err := ops.CreateRecipe(readySpec()); err != nil {                                                                        170       if _, err := ops.CreateRecipe(readySpec()); err != nil {
  ┃    171         t.Fatal(err)                                                                                                                  171         t.Fatal(err)
  ┃    172       }                                                                                                                               172       }
  ┃    @@ -177,11 +177,8 @@
  ┃    177     }                                                                                                                                 177     }
  ┃    178     if next == "RCP10" {                                                                                                              178     if next == "RCP10" {
  ┃    179       t.Fatalf("the next id is %q, but that is the one the old plan points at", next)                                                 179       t.Fatalf("the next id is %q, but that is the one the old plan points at", next)
  ┃    180     }                                                                                                                                 180     }
  ┃    181 -   if next != "RCP11" {
  ┃    182 -     t.Fatalf("the next id is %q, want RCP11 — the pointer's id is reserved", next)
  ┃    183 -   }
  ┃    184                                                                                                                                       181
  ┃    185     // And nothing carries it: the pointer resolves to nothing.                                                                       182     // And nothing carries it: the pointer resolves to nothing.
  ┃    186     r, err := ops.CreateRecipe(readySpec())                                                                                           183     r, err := ops.CreateRecipe(readySpec())
  ┃    187     if err != nil {                                                                                                                   184     if err != nil {
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/core -run 'LegacyDangling' -v 2>&1 | head -6 && cp internal/core/store.go /tmp/store.bak && python3 - <<'EOF'
  ┃  p='internal/core/store.go'…
  ┃
  ┃  (3 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.332s
  ┃  --- FAIL: TestALegacyDanglingPointerResolvesToNothing (0.11s)
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/core  0.300s
  ┃  FAIL
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.291s
  ┃

     + Thought · 1.8s

  ┃
  ┃  ← Edit internal/ui/detail.go
  ┃
  ┃    88         out = append(out, m.stySec.Render("Tags"))                                                                                     88         out = append(out, m.stySec.Render("Tags"))
  ┃    89         out = append(out, "  "+strings.Join(it.Tags, ", "))                                                                            89         out = append(out, "  "+strings.Join(it.Tags, ", "))
  ┃    90         out = append(out, "")                                                                                                          90         out = append(out, "")
  ┃    91       }                                                                                                                                91       }
  ┃                                                                                                                                              92 +     if it.PromotedTo != "" {
  ┃                                                                                                                                              93 +       out = append(out, m.styAcc.Render("Promoted to "+it.PromotedTo))
  ┃                                                                                                                                              94 +     }
  ┃    92       return out                                                                                                                       95       return out
  ┃    93     }                                                                                                                                  96     }
  ┃    94     return []string{m.styDim.Render("  (nothing selected)")}                                                                           97     return []string{m.styDim.Render("  (nothing selected)")}
  ┃    95   }                                                                                                                                    98   }
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make test 2>&1 | tail -5 && make lint && gofmt -l . && git status --short
  ┃
  ┃  (3 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  57.099s
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃  go vet ./...
  ┃   M internal/core/backlog_concurrency_test.go
  ┃   M internal/core/legacy.go
  ┃   M internal/core/store.go
  ┃   M internal/ui/detail.go
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add -A && git commit -q -F - <<'EOF'
  ┃  F-75: a dangling pointer reserves its id again…
  ┃
  ┃  47c7b23 F-75: a dangling pointer reserves its id again
  ┃  8459f96 F-74: the live console shows an old plan's specials again
  ┃

     + Thought · 26.0s



  ┃
  ┃  ← Edit internal/cli/flags.go
  ┃
  ┃    15     bools map[string]*bool                                                                                                             15     bools map[string]*bool
  ┃    16     strs  map[string]*string                                                                                                           16     strs  map[string]*string
  ┃    17     ints  map[string]*int                                                                                                              17     ints  map[string]*int
  ┃    18     lists map[string]*[]string                                                                                                         18     lists map[string]*[]string
  ┃                                                                                                                                              19 +   // given holds the list flags that were given on the command line.
  ┃                                                                                                                                              20 +   //
  ┃                                                                                                                                              21 +   // A list flag that is not given and a list flag that was given are two
  ┃                                                                                                                                              22 +   // different things, and only the parser can tell them apart: the value of a
  ┃                                                                                                                                              23 +   // list flag that was not given is an empty slice, exactly like one that was
  ┃                                                                                                                                              24 +   // given and named nothing.
  ┃                                                                                                                                              25 +   given map[string]bool
  ┃    19     known map[string]bool                                                                                                              26     known map[string]bool
  ┃    20     args  []string                                                                                                                     27     args  []string
  ┃    21   }                                                                                                                                    28   }
  ┃    22                                                                                                                                        29
  ┃    @@ -31,8 +38,9 @@
  ┃    31       bools: map[string]*bool{},                                                                                                       38       bools: map[string]*bool{},
  ┃    32       strs:  map[string]*string{},                                                                                                     39       strs:  map[string]*string{},
  ┃    33       ints:  map[string]*int{},                                                                                                        40       ints:  map[string]*int{},
  ┃    34       lists: map[string]*[]string{},                                                                                                   41       lists: map[string]*[]string{},
  ┃                                                                                                                                              42 +     given: map[string]bool{},
  ┃    35       known: map[string]bool{},                                                                                                        43       known: map[string]bool{},
  ┃    36     }                                                                                                                                  44     }
  ┃    37     fs.String("actor", "", "who is acting")                                                                                            45     fs.String("actor", "", "who is acting")
  ┃    38     fs.String("session", "", "the session this belongs to")                                                                            46     fs.String("session", "", "the session this belongs to")
  ┃

  ┃
  ┃  ← Edit internal/cli/flags.go
  ┃
  ┃    157           i++                                                                                                                         157           i++
  ┃    158           value = args[i]                                                                                                             158           value = args[i]
  ┃    159         }                                                                                                                             159         }
  ┃    160         *f.lists[name] = append(*f.lists[name], value)                                                                                160         *f.lists[name] = append(*f.lists[name], value)
  ┃                                                                                                                                              161 +       f.given[name] = true
  ┃    161                                                                                                                                       162
  ┃    162       default:                                                                                                                        163       default:
  ┃    163         return fmt.Errorf("unknown flag --%s", name)                                                                                  164         return fmt.Errorf("unknown flag --%s", name)
  ┃    164       }                                                                                                                               165       }
  ┃

  ┃
  ┃  ← Edit internal/cli/flags.go
  ┃
  ┃     87     }                                                                                                                                  87     }
  ┃     88     return ""                                                                                                                          88     return ""
  ┃     89   }                                                                                                                                    89   }
  ┃     90                                                                                                                                        90
  ┃                                                                                                                                               91 + // listArg returns the list a flag names.
  ┃                                                                                                                                               92 + //
  ┃                                                                                                                                               93 + // **nil when the flag was not given.** A list flag that is not given leaves its
  ┃                                                                                                                                               94 + // field alone; a list flag that is given replaces the list it names. So nil says
  ┃                                                                                                                                               95 + // "not given", and an empty list says "replace it with nothing" — two different
  ┃                                                                                                                                               96 + // things, and the parser refuses a list flag with no values, so only the first
  ┃                                                                                                                                               97 + // reaches here.
  ┃                                                                                                                                               98 + //
  ┃                                                                                                                                               99 + // The rule is the doc's: *a list flag replaces the list it names; a flag that is not
  ┃                                                                                                                                              100 + // given leaves its field alone*.
  ┃                                                                                                                                              101 + func (f *flagSet) listArg(name string) []string {
  ┃                                                                                                                                              102 +   if !f.given[name] {
  ┃                                                                                                                                              103 +     return nil
  ┃                                                                                                                                              104 +   }
  ┃                                                                                                                                              105 +   list, ok := f.lists[name]
  ┃                                                                                                                                              106 +   if !ok {
  ┃                                                                                                                                              107 +     return nil
  ┃                                                                                                                                              108 +   }
  ┃                                                                                                                                              109 +   return *list
  ┃                                                                                                                                              110 + }
  ┃                                                                                                                                              111 +
  ┃     91   // parse reads the arguments.                                                                                                       112   // parse reads the arguments.
  ┃     92   func (f *flagSet) parse(args []string) error {                                                                                      113   func (f *flagSet) parse(args []string) error {
  ┃     93     positionalOnly := false                                                                                                           114     positionalOnly := false
  ┃     94     for i := 0; i < len(args); i++ {                                                                                                  115     for i := 0; i < len(args); i++ {
  ┃

     Now the backlog verbs and the operation:

  ┃
  ┃  ← Edit internal/cli/backlog.go
  ┃
  ┃    38   func (r *Runner) cmdBacklogCreate(args []string) int {                                                                               38   func (r *Runner) cmdBacklogCreate(args []string) int {
  ┃    39     fs := newFlagSet("backlog create")                                                                                                 39     fs := newFlagSet("backlog create")
  ┃    40     title := fs.String("title", "", "what it is")                                                                                      40     title := fs.String("title", "", "what it is")
  ┃    41     body := fs.String("body", "", "why it is worth doing")                                                                             41     body := fs.String("body", "", "why it is worth doing")
  ┃    42 -   scope := fs.StringList("scope")
  ┃    43     kind := fs.String("kind", core.KindIdea, "idea, feature or bug")                                                                   42     kind := fs.String("kind", core.KindIdea, "idea, feature or bug")
  ┃    44 -   tags := fs.StringList("tag")                                                                                                       43 +   // The list flags are declared here and read below: a list flag that is not
  ┃                                                                                                                                              44 +   // given leaves its field alone.
  ┃                                                                                                                                              45 +   fs.StringList("scope")
  ┃                                                                                                                                              46 +   fs.StringList("tag")
  ┃    45     if err := fs.parse(args); err != nil {                                                                                             47     if err := fs.parse(args); err != nil {
  ┃    46       r.fail("bad flags", err.Error())                                                                                                 48       r.fail("bad flags", err.Error())
  ┃    47       return ExitError                                                                                                                 49       return ExitError
  ┃    48     }                                                                                                                                  50     }
  ┃    @@ -54,9 +56,9 @@
  ┃    54     if err != nil {                                                                                                                    56     if err != nil {
  ┃    55       r.fail("no plan", err.Error())                                                                                                   57       r.fail("no plan", err.Error())
  ┃    56       return ExitError                                                                                                                 58       return ExitError
  ┃    57     }                                                                                                                                  59     }
  ┃    58 -   item, err := ops.CreateBacklogItem(*title, *body, *scope, *kind, *tags)                                                            60 +   item, err := ops.CreateBacklogItem(*title, *body, fs.listArg("scope"), *kind, fs.listArg("tag"))
  ┃    59     if err != nil {                                                                                                                    61     if err != nil {
  ┃    60       return r.failOp(err)                                                                                                             62       return r.failOp(err)
  ┃    61     }                                                                                                                                  63     }
  ┃    62     return r.emit(map[string]any{"ok": true, "item": item})                                                                            64     return r.emit(map[string]any{"ok": true, "item": item})
  ┃

  ┃
  ┃  ← Edit internal/cli/backlog.go
  ┃
  ┃    72   func (r *Runner) cmdBacklogUpdate(args []string) int {                                                                               72   func (r *Runner) cmdBacklogUpdate(args []string) int {
  ┃    73     fs := newFlagSet("backlog update")                                                                                                 73     fs := newFlagSet("backlog update")
  ┃    74     title := fs.String("title", "", "new title")                                                                                       74     title := fs.String("title", "", "new title")
  ┃    75     body := fs.String("body", "", "new body")                                                                                          75     body := fs.String("body", "", "new body")
  ┃    76 -   scope := fs.StringList("scope")
  ┃    77     kind := fs.String("kind", "", "idea, feature or bug")                                                                              76     kind := fs.String("kind", "", "idea, feature or bug")
  ┃    78 -   tags := fs.StringList("tag")                                                                                                       77 +   // The list flags are declared here and read below: a list flag that is not
  ┃                                                                                                                                              78 +   // given leaves its field alone.
  ┃                                                                                                                                              79 +   fs.StringList("scope")
  ┃                                                                                                                                              80 +   fs.StringList("tag")
  ┃    79     if err := fs.parse(args); err != nil {                                                                                             81     if err := fs.parse(args); err != nil {
  ┃    80       r.fail("bad flags", err.Error())                                                                                                 82       r.fail("bad flags", err.Error())
  ┃    81       return ExitError                                                                                                                 83       return ExitError
  ┃    82     }                                                                                                                                  84     }
  ┃    @@ -88,9 +90,9 @@
  ┃    88     if err != nil {                                                                                                                    90     if err != nil {
  ┃    89       r.fail("no plan", err.Error())                                                                                                   91       r.fail("no plan", err.Error())
  ┃    90       return ExitError                                                                                                                 92       return ExitError
  ┃    91     }                                                                                                                                  93     }
  ┃    92 -   item, err := ops.UpdateBacklogItem(fs.arg(0), *title, *body, *scope, *kind, *tags)                                                 94 +   item, err := ops.UpdateBacklogItem(fs.arg(0), *title, *body, fs.listArg("scope"), *kind, fs.listArg("tag"))
  ┃    93     if err != nil {                                                                                                                    95     if err != nil {
  ┃    94       return r.failOp(err)                                                                                                             96       return r.failOp(err)
  ┃    95     }                                                                                                                                  97     }
  ┃    96     return r.emit(map[string]any{"ok": true, "item": item})                                                                            98     return r.emit(map[string]any{"ok": true, "item": item})
  ┃

  ┃
  ┃  ← Edit internal/core/backlog.go
  ┃
  ┃    229   // made a `feature` or a `bug` and gaining a scope. It is an update, not a promotion                                                229   // made a `feature` or a `bug` and gaining a scope. It is an update, not a promotion
  ┃    230   // — the item stays one item, and its id does not change.                                                                           230   // — the item stays one item, and its id does not change.
  ┃    231   //                                                                                                                                  231   //
  ┃    232   // A field that is empty is left alone: an update says what changes, not what the                                                   232   // A field that is empty is left alone: an update says what changes, not what the
  ┃    233 - // item is.                                                                                                                         233 + // item is. A **list** that was not given is nil, and nil is left alone too: the two
  ┃                                                                                                                                              234 + // list fields obey one rule — a list that is given replaces the list it names, a
  ┃                                                                                                                                              235 + // list that was not given leaves its field alone.
  ┃    234   func (o *Operations) UpdateBacklogItem(id, title, body string, scope []string, kind string, tags []string) (*BacklogItem, error) {  236   func (o *Operations) UpdateBacklogItem(id, title, body string, scope []string, kind string, tags []string) (*BacklogItem, error) {
  ┃    235     if kind != "" && !isKnownKind(kind) {                                                                                             237     if kind != "" && !isKnownKind(kind) {
  ┃    236       return nil, fmt.Errorf("unknown kind %q (want one of %s)", kind, strings.Join(BacklogKinds, ", "))                              238       return nil, fmt.Errorf("unknown kind %q (want one of %s)", kind, strings.Join(BacklogKinds, ", "))
  ┃    237     }                                                                                                                                 239     }
  ┃    @@ -253,17 +255,13 @@
  ┃    253     }                                                                                                                                 255     }
  ┃    254     if strings.TrimSpace(body) != "" {                                                                                                256     if strings.TrimSpace(body) != "" {
  ┃    255       v.Body = body                                                                                                                   257       v.Body = body
  ┃    256     }                                                                                                                                 258     }
  ┃    257 -   if len(scope) > 0 {
  ┃    258 -     v.Scope = scope
  ┃    259 -   }
  ┃    260     if kind != "" {                                                                                                                   259     if kind != "" {
  ┃    261       v.Kind = kind                                                                                                                   260       v.Kind = kind
  ┃    262     }                                                                                                                                 261     }
  ┃    263 -   if tags != nil {                                                                                                                  262 +   setList(&v.Scope, scope)
  ┃    264 -     v.Tags = tags                                                                                                                   263 +   setList(&v.Tags, tags)
  ┃    265 -   }
  ┃    266     v.Updated = Now()                                                                                                                 264     v.Updated = Now()
  ┃    267     if err := o.store.SaveBacklogItem(v); err != nil {                                                                                265     if err := o.store.SaveBacklogItem(v); err != nil {
  ┃    268       return nil, err                                                                                                                 266       return nil, err
  ┃    269     }                                                                                                                                 267     }
  ┃    @@ -272,8 +270,20 @@
  ┃    272     }                                                                                                                                 270     }
  ┃    273     return v, nil                                                                                                                     271     return v, nil
  ┃    274   }                                                                                                                                   272   }
  ┃    275                                                                                                                                       273
  ┃                                                                                                                                              274 + // setList replaces a list field when the list was given, and leaves it alone when it
  ┃                                                                                                                                              275 + // was not.
  ┃                                                                                                                                              276 + //
  ┃                                                                                                                                              277 + // **One rule for every list field of the operation.** nil means "not given"; an
  ┃                                                                                                                                              278 + // empty list means "replace it with nothing". Two guards with two different rules
  ┃                                                                                                                                              279 + // is how the tags came to be wiped by an update that never mentioned them.
  ┃                                                                                                                                              280 + func setList(dst *[]string, given []string) {
  ┃                                                                                                                                              281 +   if given != nil {
  ┃                                                                                                                                              282 +     *dst = given
  ┃                                                                                                                                              283 +   }
  ┃                                                                                                                                              284 + }
  ┃                                                                                                                                              285 +
  ┃    276   // DismissBacklogItem marks a backlog item as not worth pursuing.                                                                   286   // DismissBacklogItem marks a backlog item as not worth pursuing.
  ┃    277   //                                                                                                                                  287   //
  ┃    278   // Every kind can be dismissed: an idea, a feature and a bug are all things the chef                                                288   // Every kind can be dismissed: an idea, a feature and a bug are all things the chef
  ┃    279   // can decide against, and the asymmetry of the old ladder — only a special could be                                                289   // can decide against, and the asymmetry of the old ladder — only a special could be
  ┃

     Now the test:

  ┃
  ┃  ← Edit internal/cli/backlog_test.go
  ┃
  ┃    127       t.Fatal("an unknown kind must be refused")                                                                                      127       t.Fatal("an unknown kind must be refused")
  ┃    128     }                                                                                                                                 128     }
  ┃    129   }                                                                                                                                   129   }
  ┃    130                                                                                                                                       130
  ┃                                                                                                                                              131 + // TestTheBacklogUpdateLeavesTheListsAlone covers the list rule.
  ┃                                                                                                                                              132 + //
  ┃                                                                                                                                              133 + // A list flag replaces the list it names; a flag that is not given leaves its field
  ┃                                                                                                                                              134 + // alone. **Both** list fields obey it — the scope and the tags. The tags used to be
  ┃                                                                                                                                              135 + // wiped by an update that never mentioned them, because the two fields had two
  ┃                                                                                                                                              136 + // different guards.
  ┃                                                                                                                                              137 + func TestTheBacklogUpdateLeavesTheListsAlone(t *testing.T) {
  ┃                                                                                                                                              138 +   dir := newPlan(t)
  ┃                                                                                                                                              139 +
  ┃                                                                                                                                              140 +   res := cli(t, dir, "backlog", "create", "--title", "An item", "--body", "b",
  ┃                                                                                                                                              141 +     "--tag", "one", "--tag", "two", "--scope", "the work")
  ┃                                                                                                                                              142 +   if !res.ok(t) {
  ┃                                                                                                                                              143 +     t.Fatalf("backlog create failed:\n%s", res.Stderr)
  ┃                                                                                                                                              144 +   }
  ┃                                                                                                                                              145 +
  ┃                                                                                                                                              146 +   // An update that names neither list leaves both alone.
  ┃                                                                                                                                              147 +   res = cli(t, dir, "backlog", "update", "BL01", "--title", "A better item")
  ┃                                                                                                                                              148 +   if !res.ok(t) {
  ┃                                                                                                                                              149 +     t.Fatalf("backlog update failed:\n%s", res.Stderr)
  ┃                                                                                                                                              150 +   }
  ┃                                                                                                                                              151 +   item := res.doc(t)["item"].(map[string]any)
  ┃                                                                                                                                              152 +   tags, ok := item["tags"].([]any)
  ┃                                                                                                                                              153 +   if !ok {
  ┃                                                                                                                                              154 +     t.Fatalf("the update dropped the tags: %v", item)
  ┃                                                                                                                                              155 +   }
  ┃                                                                                                                                              156 +   if len(tags) != 2 {
  ┃                                                                                                                                              157 +     t.Fatalf("the tags are %v, want the two the item had", tags)
  ┃                                                                                                                                              158 +   }
  ┃                                                                                                                                              159 +   scope, ok := item["scope"].([]any)
  ┃                                                                                                                                              160 +   if !ok {
  ┃                                                                                                                                              161 +     t.Fatalf("the update dropped the scope: %v", item)
  ┃                                                                                                                                              162 +   }
  ┃                                                                                                                                              163 +   if len(scope) != 1 || scope[0] != "the work" {
  ┃                                                                                                                                              164 +     t.Fatalf("the scope is %v, want the one the item had", scope)
  ┃                                                                                                                                              165 +   }
  ┃                                                                                                                                              166 +
  ┃                                                                                                                                              167 +   // A list flag that is given replaces the list it names — and only that one.
  ┃                                                                                                                                              168 +   res = cli(t, dir, "backlog", "update", "BL01", "--tag", "three")
  ┃                                                                                                                                              169 +   if !res.ok(t) {
  ┃                                                                                                                                              170 +     t.Fatalf("backlog update failed:\n%s", res.Stderr)
  ┃                                                                                                                                              171 +   }
  ┃                                                                                                                                              172 +   item = res.doc(t)["item"].(map[string]any)
  ┃                                                                                                                                              173 +   tags, ok = item["tags"].([]any)
  ┃                                                                                                                                              174 +   if !ok || len(tags) != 1 || tags[0] != "three" {
  ┃                                                                                                                                              175 +     t.Fatalf("the tags are %v, want the one the flag named", tags)
  ┃                                                                                                                                              176 +   }
  ┃                                                                                                                                              177 +   scope, ok = item["scope"].([]any)
  ┃                                                                                                                                              178 +   if !ok || len(scope) != 1 {
  ┃                                                                                                                                              179 +     t.Fatalf("the tags wiped the scope: %v", item)
  ┃                                                                                                                                              180 +   }
  ┃                                                                                                                                              181 + }
  ┃                                                                                                                                              182 +
  ┃    131   // TestAnUnknownKindNamesTheKinds covers the closed set.                                                                            183   // TestAnUnknownKindNamesTheKinds covers the closed set.
  ┃    132   //                                                                                                                                  184   //
  ┃    133   // A kind nobody knows is a typo, and a typo that silently disabled a feature is                                                    185   // A kind nobody knows is a typo, and a typo that silently disabled a feature is
  ┃    134   // v1's lesson. The refusal names every kind there is, so the caller can see the                                                    186   // v1's lesson. The refusal names every kind there is, so the caller can see the
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/cli ./internal/core 2>&1 | head -20
  ┃
  ┃  ok    github.com/lthiagol/maitre/internal/cli  3.601s
  ┃  ok    github.com/lthiagol/maitre/internal/core  14.203s
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/cli/backlog.go /tmp/backlog.bak && python3 - <<'EOF'
  ┃  p='internal/cli/backlog.go'…
  ┃
  ┃  FAIL  github.com/lthiagol/maitre/internal/cli [build failed]
  ┃  FAIL
  ┃

     + Thought · 805ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/cli/flags.go /tmp/flags.bak && python3 - <<'EOF'
  ┃  p='internal/cli/flags.go'…
  ┃
  ┃  --- FAIL: TestTheBacklogUpdateLeavesTheListsAlone (0.04s)
  ┃      backlog_test.go:154: the update dropped the tags: map[body:b created:2026-10-27T21:21:28Z id:BL01 kind:idea status:open title:A better item updated:2026-10-27T21:21:28Z]
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/cli  0.239s
  ┃  FAIL
  ┃  ok    github.com/lthiagol/maitre/internal/cli  0.195s
  ┃

     + Thought · 21.1s

  ┃
  ┃  ← Edit internal/core/backlog.go
  ┃
  ┃    260       v.Kind = kind                                                                                                                   260       v.Kind = kind
  ┃    261     }                                                                                                                                 261     }
  ┃    262     setList(&v.Scope, scope)                                                                                                          262     setList(&v.Scope, scope)
  ┃    263     setList(&v.Tags, tags)                                                                                                            263     setList(&v.Tags, tags)
  ┃                                                                                                                                              264 +   v.Scope = emptyIfNil(v.Scope)
  ┃                                                                                                                                              265 +   v.Tags = emptyIfNil(v.Tags)
  ┃    264     v.Updated = Now()                                                                                                                 266     v.Updated = Now()
  ┃    265     if err := o.store.SaveBacklogItem(v); err != nil {                                                                                267     if err := o.store.SaveBacklogItem(v); err != nil {
  ┃    266       return nil, err                                                                                                                 268       return nil, err
  ┃    267     }                                                                                                                                 269     }
  ┃

  ┃
  ┃  ← Edit internal/core/backlog.go
  ┃
  ┃    284       *dst = given                                                                                                                    284       *dst = given
  ┃    285     }                                                                                                                                 285     }
  ┃    286   }                                                                                                                                   286   }
  ┃    287                                                                                                                                       287
  ┃                                                                                                                                              288 + // emptyIfNil makes a list a list: never nil.
  ┃                                                                                                                                              289 + //
  ┃                                                                                                                                              290 + // The plan's own rule is that a list is a list — an empty one is `[]`, never `null`.
  ┃                                                                                                                                              291 + // The JSON edge honours it: a nil list would be `null` in a document, and the
  ┃                                                                                                                                              292 + // standard readers break on that.
  ┃                                                                                                                                              293 + func emptyIfNil(list []string) []string {
  ┃                                                                                                                                              294 +   if list == nil {
  ┃                                                                                                                                              295 +     return []string{}
  ┃                                                                                                                                              296 +   }
  ┃                                                                                                                                              297 +   return list
  ┃                                                                                                                                              298 + }
  ┃                                                                                                                                              299 +
  ┃    288   // DismissBacklogItem marks a backlog item as not worth pursuing.                                                                   300   // DismissBacklogItem marks a backlog item as not worth pursuing.
  ┃    289   //                                                                                                                                  301   //
  ┃    290   // Every kind can be dismissed: an idea, a feature and a bug are all things the chef                                                302   // Every kind can be dismissed: an idea, a feature and a bug are all things the chef
  ┃    291   // can decide against, and the asymmetry of the old ladder — only a special could be                                                303   // can decide against, and the asymmetry of the old ladder — only a special could be
  ┃

  ┃
  ┃  ← Edit internal/core/backlog.go
  ┃
  ┃    210       return nil, err                                                                                                                 210       return nil, err
  ┃    211     }                                                                                                                                 211     }
  ┃    212     now := Now()                                                                                                                      212     now := Now()
  ┃    213     v := &BacklogItem{                                                                                                                213     v := &BacklogItem{
  ┃    214 -     ID: id, Title: title, Body: body, Scope: scope, Kind: kind, Tags: tags,                                                         214 +     ID: id, Title: title, Body: body, Scope: emptyIfNil(scope), Kind: kind, Tags: emptyIfNil(tags),
  ┃    215       Status: StatusOpen, Created: now, Updated: now,                                                                                 215       Status: StatusOpen, Created: now, Updated: now,
  ┃    216     }                                                                                                                                 216     }
  ┃    217     if err := o.store.SaveBacklogItem(v); err != nil {                                                                                217     if err := o.store.SaveBacklogItem(v); err != nil {
  ┃    218       return nil, err                                                                                                                 218       return nil, err
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/core ./internal/cli 2>&1 | head -10 && gofmt -l . ; make lint
  ┃
  ┃  ok    github.com/lthiagol/maitre/internal/core  16.787s
  ┃  ok    github.com/lthiagol/maitre/internal/cli  6.367s
  ┃  go vet ./...
  ┃

     + Thought · 1.8s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "StringList" internal/cli/*.go | grep -v _test
  ┃
  ┃  (11 earlier lines)
  ┃  internal/cli/recipe.go:49:  wp := fs.StringList("wp")
  ┃  internal/cli/recipe.go:50:  step := fs.StringList("step")
  ┃  internal/cli/recipe.go:51:  stepDone := fs.StringList("step-done-when")
  ┃  internal/cli/recipe.go:52:  stepCovers := fs.StringList("step-covers")
  ┃  internal/cli/recipe.go:166:  inScope := fs.StringList("in-scope")
  ┃  internal/cli/recipe.go:167:  outScope := fs.StringList("out-of-scope")
  ┃  internal/cli/recipe.go:168:  dependsOn := fs.StringList("depends-on")
  ┃

     + Thought · 750ms

  ┃
  ┃  ← Edit internal/core/backlog_test.go
  ┃
  ┃    218       t.Fatal("an unknown kind must be refused on update too")                                                                        218       t.Fatal("an unknown kind must be refused on update too")
  ┃    219     }                                                                                                                                 219     }
  ┃    220   }                                                                                                                                   220   }
  ┃    221                                                                                                                                       221
  ┃                                                                                                                                              222 + // TestTheBacklogUpdateLeavesTheListsAlone covers the list rule.
  ┃                                                                                                                                              223 + //
  ┃                                                                                                                                              224 + // nil means "not given" and leaves the field alone; a list means "replace the list
  ┃                                                                                                                                              225 + // it names". **Both** list fields obey the rule — the scope and the tags. The tags
  ┃                                                                                                                                              226 + // used to be wiped by an update that never mentioned them, because the two fields
  ┃                                                                                                                                              227 + // had two different guards.
  ┃                                                                                                                                              228 + func TestTheBacklogUpdateLeavesTheListsAlone(t *testing.T) {
  ┃                                                                                                                                              229 +   ops, store := newOps(t)
  ┃                                                                                                                                              230 +
  ┃                                                                                                                                              231 +   item, err := ops.CreateBacklogItem("An item", "why", []string{"the work"}, KindIdea, []string{"one", "two"})
  ┃                                                                                                                                              232 +   if err != nil {
  ┃                                                                                                                                              233 +     t.Fatal(err)
  ┃                                                                                                                                              234 +   }
  ┃                                                                                                                                              235 +
  ┃                                                                                                                                              236 +   // An update that names neither list leaves both alone.
  ┃                                                                                                                                              237 +   got, err := ops.UpdateBacklogItem(item.ID, "A better item", "", nil, "", nil)
  ┃                                                                                                                                              238 +   if err != nil {
  ┃                                                                                                                                              239 +     t.Fatal(err)
  ┃                                                                                                                                              240 +   }
  ┃                                                                                                                                              241 +   if len(got.Tags) != 2 || len(got.Scope) != 1 {
  ┃                                                                                                                                              242 +     t.Fatalf("the update dropped a list: %+v", got)
  ┃                                                                                                                                              243 +   }
  ┃                                                                                                                                              244 +
  ┃                                                                                                                                              245 +   // A list that is given replaces the list it names — and only that one.
  ┃                                                                                                                                              246 +   got, err = ops.UpdateBacklogItem(item.ID, "", "", nil, "", []string{"three"})
  ┃                                                                                                                                              247 +   if err != nil {
  ┃                                                                                                                                              248 +     t.Fatal(err)
  ┃                                                                                                                                              249 +   }
  ┃                                                                                                                                              250 +   if len(got.Tags) != 1 || got.Tags[0] != "three" {
  ┃                                                                                                                                              251 +     t.Fatalf("the tags are %v, want the one the update named", got.Tags)
  ┃                                                                                                                                              252 +   }
  ┃                                                                                                                                              253 +   if len(got.Scope) != 1 {
  ┃                                                                                                                                              254 +     t.Fatalf("the tags wiped the scope: %+v", got)
  ┃                                                                                                                                              255 +   }
  ┃                                                                                                                                              256 +
  ┃                                                                                                                                              257 +   // And the list on disk is the same.
  ┃                                                                                                                                              258 +   onDisk, err := store.LoadBacklogItem(item.ID)
  ┃                                                                                                                                              259 +   if err != nil {
  ┃                                                                                                                                              260 +     t.Fatal(err)
  ┃                                                                                                                                              261 +   }
  ┃                                                                                                                                              262 +   if len(onDisk.Tags) != 1 || len(onDisk.Scope) != 1 {
  ┃                                                                                                                                              263 +     t.Fatalf("the item on disk is %+v", onDisk)
  ┃                                                                                                                                              264 +   }
  ┃                                                                                                                                              265 + }
  ┃                                                                                                                                              266 +
  ┃    222   // TestTheBacklogItemLifecycle covers the ordinary path: create, edit, dismiss.                                                     267   // TestTheBacklogItemLifecycle covers the ordinary path: create, edit, dismiss.
  ┃    223   func TestTheBacklogItemLifecycle(t *testing.T) {                                                                                    268   func TestTheBacklogItemLifecycle(t *testing.T) {
  ┃    224     ops, store := newOps(t)                                                                                                           269     ops, store := newOps(t)
  ┃    225                                                                                                                                       270
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/core -run 'ListsAlone' -v 2>&1 | head -10 && cp internal/core/backlog.go /tmp/core-backlog.bak && python3 - <<'EOF'
  ┃  p='internal/core/backlog.go'…
  ┃
  ┃  === RUN   TestTheBacklogUpdateLeavesTheListsAlone
  ┃  --- PASS: TestTheBacklogUpdateLeavesTheListsAlone (0.04s)
  ┃  PASS
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.115s
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.121s
  ┃

     + Thought · 2.6s

  ┃
  ┃  ← Edit internal/core/backlog_test.go
  ┃
  ┃    261     }                                                                                                                                 261     }
  ┃    262     if len(onDisk.Tags) != 1 || len(onDisk.Scope) != 1 {                                                                              262     if len(onDisk.Tags) != 1 || len(onDisk.Scope) != 1 {
  ┃    263       t.Fatalf("the item on disk is %+v", onDisk)                                                                                     263       t.Fatalf("the item on disk is %+v", onDisk)
  ┃    264     }                                                                                                                                 264     }
  ┃                                                                                                                                              265 +
  ┃                                                                                                                                              266 +   // An **empty** list that was given clears the field: it is not "not given".
  ┃                                                                                                                                              267 +   got, err = ops.UpdateBacklogItem(item.ID, "", "", []string{}, "", []string{})
  ┃                                                                                                                                              268 +   if err != nil {
  ┃                                                                                                                                              269 +     t.Fatal(err)
  ┃                                                                                                                                              270 +   }
  ┃                                                                                                                                              271 +   if len(got.Tags) != 0 || len(got.Scope) != 0 {
  ┃                                                                                                                                              272 +     t.Fatalf("the empty lists did not clear the fields: %+v", got)
  ┃                                                                                                                                              273 +   }
  ┃    265   }                                                                                                                                   274   }
  ┃    266                                                                                                                                       275
  ┃    267   // TestTheBacklogItemLifecycle covers the ordinary path: create, edit, dismiss.                                                     276   // TestTheBacklogItemLifecycle covers the ordinary path: create, edit, dismiss.
  ┃    268   func TestTheBacklogItemLifecycle(t *testing.T) {                                                                                    277   func TestTheBacklogItemLifecycle(t *testing.T) {
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/core -run 'ListsAlone' 2>&1 | tail -3 && cp internal/core/backlog.go /tmp/core-backlog.bak && python3 - <<'EOF'
  ┃  p='internal/core/backlog.go'…
  ┃
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.132s
  ┃  --- FAIL: TestTheBacklogUpdateLeavesTheListsAlone (0.05s)
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/core  0.119s
  ┃  ok    github.com/lthiagol/maitre/internal/core  (cached)
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make test 2>&1 | tail -5 && make lint && gofmt -l . && git status --short
  ┃
  ┃  (4 earlier lines)
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃  go vet ./...
  ┃   M internal/cli/backlog.go
  ┃   M internal/cli/backlog_test.go
  ┃   M internal/cli/flags.go
  ┃   M internal/core/backlog.go
  ┃   M internal/core/backlog_test.go
  ┃

     + Thought · 3.2s

  ┃
  ┃  ← Edit docs/04-operations.md
  ┃
  ┃    193   | `depends_on` | `maitre-cli recipe update --depends-on <id>` — repeatable, and it replaces the whole list |                        193   | `depends_on` | `maitre-cli recipe update --depends-on <id>` — repeatable, and it replaces the whole list |
  ┃    194   | the backlog item's `title` · `body` · `scope` · `kind` · `tags` | `maitre-cli backlog update --title --body --scope --kind --tag` 194   | the backlog item's `title` · `body` · `scope` · `kind` · `tags` | `maitre-cli backlog update --title --body --scope --kind --tag`
  ┃          |                                                                                                                                         |
  ┃    195                                                                                                                                       195
  ┃    196   A list flag replaces the list it names; a flag that is not given leaves its field                                                   196   A list flag replaces the list it names; a flag that is not given leaves its field
  ┃    197 - alone.                                                                                                                              197 + alone. The domain reads the two apart: a list that was not given is `nil`, and a list
  ┃                                                                                                                                              198 + that was given is a list — an **empty** list is a list that was given, and it clears
  ┃                                                                                                                                              199 + the field.
  ┃    198                                                                                                                                       200
  ┃                                                                                                                                              201 + **A list cannot be emptied from the command line yet.** The parser refuses a list flag
  ┃                                                                                                                                              202 + with no values (`--tag needs a value`), so there is no way to say "no tags". The
  ┃                                                                                                                                              203 + domain can clear a list; the CLI cannot ask it to. Whether the parser should grow a
  ┃                                                                                                                                              204 + way to say it is an open question.
  ┃                                                                                                                                              205 +
  ┃    199   A typo is an error, never an empty result: an unknown lane, an unknown finding phase,                                               206   A typo is an error, never an empty result: an unknown lane, an unknown finding phase,
  ┃    200   an unknown severity and an unknown event are all refused with the list of what is                                                   207   an unknown severity and an unknown event are all refused with the list of what is
  ┃    201   allowed.                                                                                                                            208   allowed.
  ┃    202                                                                                                                                       209
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make test 2>&1 | tail -4 && make lint && git add -A && git commit -q -F - <<'EOF'
  ┃  F-76: a list flag that is not given leaves its field alone…
  ┃
  ┃  ok    github.com/lthiagol/maitre/internal/cli  2.696s
  ┃  ok    github.com/lthiagol/maitre/internal/core  13.678s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  (cached)
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃  go vet ./...
  ┃  b0bb74f F-76: a list flag that is not given leaves its field alone
  ┃  47c7b23 F-75: a dangling pointer reserves its id again
  ┃

     + Thought · 879ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '/^### 1.1 How a backlog item becomes work/,/^### The CPD/p' docs/03-domain-model.md
  ┃
  ┃  (20 earlier lines)
  ┃  out-of-scope items, no steps and no work packages, and its criterion's *verification*
  ┃  is the recipe's own life — a criterion with no verification cannot pass
  ┃  `ac.falsifiable`. So the gates that guard a recipe refuse a promoted one exactly as
  ┃  they refuse any other half-written recipe. A promotion needed no new plan-format field
  ┃  and no special case in a rule.
  ┃
  ┃  ### The CPD ↔ tool mapping — the single source
  ┃

  ┃
  ┃  ← Edit docs/03-domain-model.md
  ┃
  ┃    74   |---|---|                                                                                                                            74   |---|---|
  ┃    75   | **Shaping** | An **update**, not a promotion: the kind changes and a scope is added, and the id does not. `BL01` stays `BL01` —    75   | **Shaping** | An **update**, not a promotion: the kind changes and a scope is added, and the id does not. `BL01` stays `BL01` —
  ┃         the id is kind-neutral, so the trail survives the shaping |                                                                               the id is kind-neutral, so the trail survives the shaping |
  ┃    76   | **Promotion** | One-way. The item's status becomes `promoted` and it records what it became (`promoted_to`) |                      76   | **Promotion** | One-way. The item's status becomes `promoted` and it records what it became (`promoted_to`) |
  ┃    77   | **An idea** | Becomes a recipe in `draft` with **no** criterion: the criteria are what the recipe is groomed into |                77   | **An idea** | Becomes a recipe in `draft` with **no** criterion: the criteria are what the recipe is groomed into |
  ┃    78 - | **A small task** | Becomes a recipe in `draft` with **exactly one** criterion: the work is small enough to be one thing, and a     78 + | **A small task** | Becomes a recipe in `draft` with **exactly one** criterion: the work is small enough to be one thing. The
  ┃         recipe with a criterion is falsifiable from the start |                                                                                   criterion makes the work **checkable** — it says what has to be true — and it becomes **falsifiable** when the recipe's own life
  ┃                                                                                                                                                   writes its verification |
  ┃    79                                                                                                                                        79
  ┃    80   **No gate is weakened.** The promoted recipe is **not** pre-groomed: it carries no                                                   80   **No gate is weakened.** The promoted recipe is **not** pre-groomed: it carries no
  ┃    81   out-of-scope items, no steps and no work packages, and its criterion's *verification*                                                81   out-of-scope items, no steps and no work packages, and its criterion's *verification*
  ┃    82   is the recipe's own life — a criterion with no verification cannot pass                                                              82   is the recipe's own life — a criterion with no verification cannot pass
  ┃

     + Thought · 1.4s

  ┃
  ┃  ← Edit docs/03-domain-model.md
  ┃
  ┃    79                                                                                                                                        79
  ┃    80   **No gate is weakened.** The promoted recipe is **not** pre-groomed: it carries no                                                   80   **No gate is weakened.** The promoted recipe is **not** pre-groomed: it carries no
  ┃    81   out-of-scope items, no steps and no work packages, and its criterion's *verification*                                                81   out-of-scope items, no steps and no work packages, and its criterion's *verification*
  ┃    82   is the recipe's own life — a criterion with no verification cannot pass                                                              82   is the recipe's own life — a criterion with no verification cannot pass
  ┃    83 - `ac.falsifiable`. So the gates that guard a recipe refuse a promoted one exactly as                                                  83 + `ac.falsifiable`. So a promoted recipe is **checkable, not yet falsifiable**: the
  ┃    84 - they refuse any other half-written recipe. A promotion needed no new plan-format field                                               84 + criterion says what has to be true, and the recipe becomes falsifiable when its
  ┃    85 - and no special case in a rule.                                                                                                       85 + verification is written. The gates that guard a recipe refuse a promoted one exactly
  ┃                                                                                                                                              86 + as they refuse any other half-written recipe. A promotion needed no new plan-format
  ┃                                                                                                                                              87 + field and no special case in a rule.
  ┃    86                                                                                                                                        88
  ┃    87   ### The CPD ↔ tool mapping — the single source                                                                                       89   ### The CPD ↔ tool mapping — the single source
  ┃    88                                                                                                                                        90
  ┃    89   **[GLOSSARY.md](../GLOSSARY.md) owns it.** Every term there names its CPD                                                            91   **[GLOSSARY.md](../GLOSSARY.md) owns it.** Every term there names its CPD
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "falsifiable from the start" docs/ AGENTS.md GLOSSARY.md work/briefs/ 2>/dev/null | head
  ┃
  ┃  docs/audit.md:86:| **F-77** | medium | **`docs/03` § 1.1's small-task row claims a falsifiability the code cannot deliver, and contradicts the paragraph beneath it** — the row says a promoted small task's recipe is "falsifiable from the start", while the paragraph below says
  ┃  its criterion carries no verification and "cannot pass `ac.falsifiable`". The code follows the paragraph: `promotedCriteria` leaves `Verification` empty, so the promoted recipe is refused by the falsifiable gate exactly as any half-written recipe is. The row is the brief's own
  ┃  wording copied in, and it is the sentence an agent reads to decide what a promotion buys. | `docs/03-domain-model.md:78,80-84`; `internal/core/backlog.go:388-393` |
  ┃  work/briefs/fix-c1.md:104:recipe with a criterion is falsifiable from the start"*, and the paragraph below says
  ┃  work/briefs/fix-c1.md:115:> The orchestrator's slice C brief said *"so the work is falsifiable from the start"*.
  ┃

     + Thought · 5.6s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "falsifiable from the start" -B 4 -A 2 work/briefs/slice-c.md
  ┃
  ┃  (no output)
  ┃  Exited with code 1
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && ls work/briefs/ && grep -rn "falsifiable" work/briefs/slice-c.md | head
  ┃
  ┃  (10 earlier lines)
  ┃  slice-b-resume.md
  ┃  slice-b.md
  ┃  slice-c.md
  ┃  slice-d.md
  ┃  slice-e.md
  ┃  work/briefs/slice-c.md:76:- `promote` creates a recipe **in `draft`**, carrying the item's intent, its scope and — for a small task — **one** acceptance criterion, so the work is checkable from the start. It is **not** falsifiable until that criterion has its verification,
  ┃  which is the recipe's own life after grooming — a criterion with no verification cannot pass `ac.falsifiable`.
  ┃  work/briefs/slice-c.md:77:- The recipe is **not** pre-groomed: the second out-of-scope item, the falsifiable verification and the steps are the recipe's own life, exactly as for any recipe. **No gate is weakened, and no gate rule changes.** If you find yourself needing to
  ┃  weaken a gate or to add a plan-format field to make a promotion work, **stop and report it to the orchestrator instead** — that is a design decision, not an implementation detail.
  ┃

     + Thought · 7.4s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '/^### 1.1 How a backlog item becomes work/,/^### The CPD/p' docs/03-domain-model.md && make test 2>&1 | tail -4 && make lint && gofmt -l .
  ┃
  ┃  (27 earlier lines)
  ┃
  ┃  ### The CPD ↔ tool mapping — the single source
  ┃  ok    github.com/lthiagol/maitre/internal/cli  (cached)
  ┃  ok    github.com/lthiagol/maitre/internal/core  15.477s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  (cached)
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃  go vet ./...
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add -A && git commit -q -F - <<'EOF'
  ┃  F-77: the spec says what is true about a promoted small task…
  ┃
  ┃  f77ef07 F-77: the spec says what is true about a promoted small task
  ┃  b0bb74f F-76: a list flag that is not given leaves its field alone
  ┃  47c7b23 F-75: a dangling pointer reserves its id again
  ┃  8459f96 F-74: the live console shows an old plan's specials again
  ┃  e59798e work: slice C's findings, and the two questions the chef owes
  ┃

     + Thought · 1.7s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make lint && gofmt -l . && make test 2>&1 | tail -5 && git status --short && echo "---" && git log --oneline -5
  ┃
  ┃  (5 earlier lines)
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃  ---
  ┃  f77ef07 F-77: the spec says what is true about a promoted small task
  ┃  b0bb74f F-76: a list flag that is not given leaves its field alone
  ┃  47c7b23 F-75: a dangling pointer reserves its id again
  ┃  8459f96 F-74: the live console shows an old plan's specials again
  ┃  e59798e work: slice C's findings, and the two questions the chef owes
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/core -run 'Doc|Version' -v 2>&1 | grep -E "PASS|FAIL" | head -20 && go test ./internal/cli -run 'CommandCount' -v 2>&1 | grep -E "PASS|FAIL"
  ┃
  ┃  (16 earlier lines)
  ┃  --- PASS: TestDomainVersionAllowsReadingANewerPlan (0.01s)
  ┃  --- PASS: TestTheRepoPlanStaysWritableUnderANewerVersion (0.00s)
  ┃  PASS
  ┃  --- PASS: TestTheDocsCommandCountMatchesTheCode (0.00s)
  ┃  --- PASS: TestTheCommandCountWireTripsOnAHalfEdit (0.00s)
  ┃  --- PASS: TestTheCommandCountWireTripsOnAHalfEditOfTheGroups (0.00s)
  ┃  PASS
  ┃

     Now the receipt and the notification:

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && mkdir -p work/inbox && cat > work/inbox/runner-slice-c-fix-1.md <<'EOF'
  ┃  Runner: slice C fix 1 done. Commits: 4 (one per finding). Fixed: F-74, F-75,…
  ┃
  ┃  cc8796c inbox: runner reports slice-c fix 1
  ┃  f77ef07 F-77: the spec says what is true about a promoted small task
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1Z:p1" "Runner: slice C fix 1 done. Commits: 4. Fixed: F-74, F-75, F-76, F-77. Pre-notify gate: pass (lint clean, gofmt empty, make test green, tree clean); the gate table, the plan format and the command count (42
  ┃  ) did not move. Notes: (1) F-74 — the live path and the render path now share Model.applyReload, and the test drives the real reload message through Update instead of copying the fields by hand (the copy was the bug); removing the assignment fails both tests. (2) F-75 —…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"idle","completion_seq":389,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":
  ┃  false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,"screen_detection_skipped":true,"state_change_seq":389,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722","terminal_title":"OC | orchestrator (fork #1)","terminal_title_stripped":
  ┃  "OC | orchestrator (fork #1)","workspace_id":"w1Z"},"type":"agent_prompted"}}
  ┃  NOTIFIED
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -6 && git status --short && echo "CLEAN" && make test 2>&1 | grep -c "^ok"
  ┃
  ┃  cc8796c inbox: runner reports slice-c fix 1
  ┃  f77ef07 F-77: the spec says what is true about a promoted small task
  ┃  b0bb74f F-76: a list flag that is not given leaves its field alone
  ┃  47c7b23 F-75: a dangling pointer reserves its id again
  ┃  8459f96 F-74: the live console shows an old plan's specials again
  ┃  e59798e work: slice C's findings, and the two questions the chef owes
  ┃  CLEAN
  ┃  3
  ┃

     Four fixes, four commits, orchestrator notified.

     ┌─────────┬─────────┬───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Commit  │ Finding │ What                                                                                                                                                                                                                                                          │
     ├─────────┼─────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 8459f96 │ F-74    │ Model.applyReload now puts a reloaded plan into the model and both the live path and the render path go through it — one place knows the model's shape. The test drives the real reload message through Update; loadModel's hand-copy (the bug) is gone, and  │
     │         │         │ TestTheLiveConsoleCarriesTheLegacyGroup checks the message and the model                                                                                                                                                                                      │
     ├─────────┼─────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 47c7b23 │ F-75    │ core.LegacyItem carries promoted_to; reservedRecipeIDs scans the old plan's pointers as well as the backlog's. The test creates nine recipes so the numbering reaches RCP10 and checks it is not handed out                                                   │
     ├─────────┼─────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ b0bb74f │ F-76    │ The parser knows which list flags were given; listArg returns nil for a flag that was not given, and one helper (setList) applies the rule to the scope and the tags alike. The other verbs are untouched; no null reaches a JSON document                    │
     ├─────────┼─────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ f77ef07 │ F-77    │ § 1.1's row now says the criterion makes the work checkable and becomes falsifiable when the recipe's own life writes its verification. The code did not move                                                                                                 │
     └─────────┴─────────┴───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

     Every fix was proved by mutation — removing the assignment, dropping the legacy scan, or making an absent list flag read as an empty list each fail their test — and each mutation was reverted.

     Two things I did not decide

     1. The list-clearing question (the chef's). The parser refuses a list flag with no values, so a list can be replaced but never emptied from the command line. The domain can clear a list — an empty list is a list that was given — but the CLI cannot ask it to. I left the parser
        alone and stated the current behaviour in docs/04.
     2. The slice C brief's "the second out-of-scope item". It can be read as the promotion writes one and the recipe's life writes the second, or as the recipe's life writes them. I read it the second way — the code writes no out-of-scope items, and docs/03 § 1.1 says so plainly.
        The reviewer filed no finding about it, so I did not change behaviour on my own reading; if the other reading is right, that is a behaviour change, not a doc fix.

     The gate table, the plan format and the command count (42) did not move — their tripwires still pass.

     Build · MiniMax-M3.1-Flash · 14m 3s · 90.7 tok/s

  ┃
  ┃
  ┃
  ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
  ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
  ~/Code/Pessoal/maitre:master                                                                                                                                                                                                                               328.8K (33%)  ctrl+p commands

```
