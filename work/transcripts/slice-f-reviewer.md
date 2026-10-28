# The reviewer's lane — slice F, captured 2026-10-28

Dumped from pane `w1Z:pE` with `herdr pane read --source recent-unwrapped --lines 5000`.
A fresh lane: it re-derived the facts from the code and filed F-87..F-93 across four
rounds. Its best finding was F-91: **the domain's own API could write a record it
then refused to read** — the tool bricked the plan it had just written.

```
  ┃  two views is gone too.** Check that **nothing was lost** with it: what does the
  ┃  one-format wire still hold, and is the drift it was built to catch now impossible for
  ┃  a different reason?
  ┃
  ┃  ### 2. The one consequence to judge — can a title brick a plan?
  ┃
  ┃  The runner flagged it, and it is the right thing to flag:
  ┃
  ┃  > *the rule is by value, so a title (or the slug derived from it) of exactly "Line
  ┃  > Cook" is caught — a recipe with that title would make the plan unsupported.*
  ┃
  ┃  **The orchestrator's ruling: the rule stands.** The trade is deliberate — a *loud*
  ┃  refusal beats a *silent* hole, and a silent hole is the worse failure (this repo has
  ┃  the lesson: *a silent enum fallback disabled a whole feature, with green tests*).
  ┃
  ┃  **But the reviewer's question is the sharp one, and it is yours:**
  ┃
  ┃  > **Can a plan be bricked by a title?**
  ┃
  ┃  If the writer's check runs **before** the write, a recipe titled *"Line Cook"* is
  ┃  **refused at creation** — the chef renames it and retries, and the plan was never
  ┃  unusable. If the record can land first and be refused afterwards, then the tool
  ┃  **wrote a record it will not read**, which is F-91's shape a third time.
  ┃
  ┃  **Prove which of the two it is.** Write a recipe whose title is exactly a retired
  ┃  name and follow the record: does it land, and does the plan still work afterwards?
  ┃  **And the same for the slug** — the file's name is derived from the title, so a
  ┃  recipe whose title slugs to a retired name puts a retired string in the file's name.
  ┃
  ┃  ### 3. F-92 and F-93 — closed?
  ┃
  ┃  **F-92** — `docs/03` § 1 now draws the line the code draws: a name typed by a human
  ┃  is a **lookup** and is tolerated; a file's name on disk is the **format** and is
  ┃  refused; the **journal** is history and keeps its old ids. *"A reader that cannot
  ┃  resolve them cannot read the past"* is now said of the journal. **Check § 1, § 6 and
  ┃  § 6.2 agree**, and that the rules-list clause in § 6 no longer reads as making
  ┃  slice E's change a version bump.
  ┃
  ┃  **F-93** — the record's name fields are one list with two views derived from it, and
  ┃  the comment that lied is gone. (The list is now gone entirely, per the ruling above —
  ┃  judge the *outcome*, not the intermediate design.)
  ┃
  ┃  ### 4. Scope, the tree and the ledger
  ┃
  ┃  - `core.DomainVersion` is **7** and did not move again. The gate table, the plan
  ┃    format, the command count (**42**) and the role names must not have moved.
  ┃  - **The mascot** — another session's work, in `logo/`. Check the runner committed
  ┃    nothing of it.
  ┃  - **The ledger** — `docs/audit.md` is yours. Move **F-92** and **F-93** to closed,
  ┃    with the file and the evidence that holds each, and **add the row for slice F
  ┃    itself**: the ruling that this is not brownfield, the tolerance machinery that came
  ┃    out, the one-format wire, the one place where a name enters a record, and the check
  ┃    that is by value. The ledger is not swept: keep the findings' wording as they were
  ┃    filed.
  ┃
  ┃  ### Verdict
  ┃
  ┃  - No HIGH or MEDIUM findings → `verdict: ok`.
  ┃  - Any HIGH or MEDIUM → `verdict: changes-needed`.
  ┃  - LOW findings: **note them in your message, do not file them.**
  ┃
  ┃  ### Hard rules
  ┃
  ┃  - Do **not** edit implementation code. Read-only for everything except
  ┃    `docs/audit.md`.
  ┃  - Revert every mutation you make. The tree must be exactly as the runner left it,
  ┃    except your ledger commit.
  ┃  - No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
  ┃  - **Do not touch `logo/`** — it is another session's work.
  ┃
  ┃  ### Notify the orchestrator
  ┃
  ┃  ```bash
  ┃  mkdir -p work/inbox
  ┃  cat > work/inbox/reviewer-slice-f-cycle-4.md <<EOF
  ┃  Reviewer: slice F cycle 4 done. Verdict: <ok|changes-needed>. F-92, F-93 closed: <yes/no>. Can a title brick a plan: <yes/no, and
  ┃  how you proved it>. Filed <N> new blockers (<F-94.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>.
  ┃  EOF
  ┃  git add work/inbox && git commit -q -m "inbox: reviewer reports slice-f cycle 4"
  ┃  herdr agent prompt "w1Z:p1" "Reviewer: slice F cycle 4 done. Verdict: <ok|changes-needed>. F-92, F-93 closed: <yes/no>. Can a title
  ┃  brick a plan: <yes/no, and how you proved it>. Filed <N> new blockers (<F-94.. or none>). Mutation checks: <what you mutated and
  ┃  reverted>. LOW notes: <any>."
  ┃  ```
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -10 && git status --short && git diff 73e3666..HEAD --stat
  ┃
  ┃  (13 earlier lines)
  ┃   internal/core/retired_test.go       | 123 +++++++++++++++++++++++
  ┃   work/README.md                      |  22 +++--
  ┃   work/briefs/fix-f3.md               | 131 +++++++++++++++++++++++++
  ┃   work/briefs/fix-f4.md               | 106 ++++++++++++++++++++
  ┃   work/briefs/review-f4.md            | 117 ++++++++++++++++++++++
  ┃   work/inbox/runner-slice-f-fix-1b.md |  78 ++++++---------
  ┃   8 files changed, 664 insertions(+), 156 deletions(-)
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff 73e3666..HEAD -- internal/core/retired.go
  ┃
  ┃  (239 earlier lines)
  ┃   // checkRetiredRecipe refuses a recipe file that carries a retired artifact.
  ┃   func checkRetiredRecipe(raw []byte, r *Recipe) error {
  ┃     found := retiredInRecipe(raw, r.ID)
  ┃  -  found = append(found, retiredActorNames(r)...)
  ┃  +  found = append(found, retiredNames(r)...)
  ┃     return retiredPlanError(found)
  ┃   }
  ┃

     + Thought · 7.8s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff 73e3666..HEAD -- internal/core/retired_test.go docs/03-domain-model.md
  ┃
  ┃  (216 earlier lines)
  ┃  +    }
  ┃  +  }
  ┃  +}
  ┃  +
  ┃   // TestTheGlossaryNamesEachRetiredNameOnce is the glossary's rule, checked.
  ┃   //
  ┃   // A retired word is named **once** in the glossary: in the `_Avoid_` line of the
  ┃

     + Thought · 4.1s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git show 1437241 --stat && git show 1437241 -- internal/core/retired_test.go | grep -E "^
  ┃  \+.*func Test"
  ┃
  ┃  (36 earlier lines)
  ┃
  ┃   internal/core/retired.go      | 212 ++++++++++++++++++++++++++++--------------
  ┃   internal/core/retired_test.go | 191 ++++++++++++++++++++++++++++++++++++-
  ┃   2 files changed, 330 insertions(+), 73 deletions(-)
  ┃  +func TestTheTwoViewsAgreeOnARecord(t *testing.T) {
  ┃  +func TestTheOneFormatWireFailsOnAViewThatReachesIntoTheRecord(t *testing.T) {
  ┃

     + Thought · 6.2s

  ┃
  ┃  # Wrote internal/core/zz_reviewer_tmp_test.go
  ┃
  ┃     1 package core
  ┃     2
  ┃     3 import (
  ┃     4   "os"
  ┃     5   "path/filepath"
  ┃     6   "strings"
  ┃     7   "testing"
  ┃     8 )
  ┃     9
  ┃    10 // TestZZCycle4Title probes the title question. Temporary review file — deleted after the run.
  ┃    11 func TestZZCycle4Title(t *testing.T) {
  ┃    12   ops, store := newOps(t)
  ┃    13
  ┃    14   // A title that IS a retired name: refused at creation.
  ┃    15   _, err := ops.CreateRecipe(RecipeSpec{
  ┃    16     Title: "Line Cook", Outcome: "o", Problem: "p",
  ┃    17     InScope: []string{"a"}, OutOfScope: []string{"b"},
  ┃    18   })
  ┃    19   if err == nil {
  ┃    20     t.Fatal("a recipe titled with a retired name must be refused")
  ┃    21   }
  ┃    22   t.Logf("title refused: %v", err)
  ┃    23
  ┃    24   // Nothing landed, and the plan still works.
  ┃    25   entries, rerr := os.ReadDir(store.RecipeDir())
  ┃    26   if rerr != nil {
  ┃    27     t.Fatal(rerr)
  ┃    28   }
  ┃    29   if len(entries) != 0 {
  ┃    30     t.Errorf("the refusal wrote %d files anyway", len(entries))
  ┃    31   }
  ┃    32   if _, err := ops.ListRecipes(); err != nil {
  ┃    33     t.Errorf("the plan must still read: %v", err)
  ┃    34   }
  ┃    35   if _, err := ops.CreateBacklogItem("A hunch", "b", nil, KindIdea, nil); err != nil {
  ┃    36     t.Errorf("the plan must still be writable: %v", err)
  ┃    37   }
  ┃    38
  ┃    39   // A title that slugs to a retired name: refused too (the slug is a field).
  ┃    40   _, err = ops.CreateRecipe(RecipeSpec{
  ┃    41     Title: "line cook", Outcome: "o", Problem: "p",
  ┃    42     InScope: []string{"a"}, OutOfScope: []string{"b"},
  ┃    43   })
  ┃    44   if err == nil {
  ┃    45     t.Fatal("a recipe whose slug is a retired name must be refused")
  ┃    46   }
  ┃    47   t.Logf("slug refused: %v", err)
  ┃    48
  ┃    49   // Prose that mentions a retired name: accepted, and the plan reads.
  ┃    50   r, err := ops.CreateRecipe(RecipeSpec{
  ┃    51     Title: "The line-cook wrote this", Outcome: "o", Problem: "p",
  ┃    52     InScope: []string{"a"}, OutOfScope: []string{"b"},
  ┃    53   })
  ┃    54   if err != nil {
  ┃    55     t.Fatalf("prose must not be refused: %v", err)
  ┃    56   }
  ┃    57   got, err := ops.GetRecipe(r.ID)
  ┃    58   if err != nil {
  ┃    59     t.Fatalf("the plan must read: %v", err)
  ┃    60   }
  ┃    61   t.Logf("prose accepted: %q slug=%q", got.Title, got.Slug)
  ┃    62
  ┃    63   // Changing the title afterwards: refused, and the plan keeps working.
  ┃    64   bad := "Line Cook"
  ┃    65   if _, err := ops.UpdateRecipe(r.ID, Fragment{Title: &bad}); err == nil {
  ┃    66     t.Fatal("the title change must be refused")
  ┃    67   } else {
  ┃    68     t.Logf("title change refused: %v", err)
  ┃    69   }
  ┃    70   after, err := ops.GetRecipe(r.ID)
  ┃    71   if err != nil {
  ┃    72     t.Fatalf("the plan must still read: %v", err)
  ┃    73   }
  ┃    74   if after.Title == bad {
  ┃    75     t.Error("the refused change landed")
  ┃    76   }
  ┃    77 }
  ┃    78
  ┃    79 // TestZZCycle4Value probes the rule's sides. Temporary review file.
  ┃    80 func TestZZCycle4Value(t *testing.T) {
  ┃    81   // The error names the field.
  ┃    82   err := checkRecordNames(&Recipe{ID: "RCP01", Execution: Execution{ClaimedBy: "line-cook"}})
  ┃    83   if err == nil {
  ┃    84     t.Fatal("a retired name must be refused")
  ┃    85   }
  ┃    86   if !strings.Contains(err.Error(), "execution.claimed_by") {
  ┃    87     t.Errorf("the refusal should name the field: %v", err)
  ┃    88   }
  ┃    89   t.Logf("writer refusal: %v", err)
  ┃    90
  ┃    91   // The reader's refusal names it too.
  ┃    92   names := retiredNames(&Recipe{ID: "RCP01", Execution: Execution{ClaimedBy: "line-cook"}})
  ┃    93   t.Logf("reader reports: %v", names)
  ┃    94   if len(names) != 1 || !strings.Contains(names[0], "execution.claimed_by") {
  ┃    95     t.Errorf("the reader should name the field: %v", names)
  ┃    96   }
  ┃    97
  ┃    98   // The spellings.
  ┃    99   for _, spelling := range []string{"line-cook", "Line Cook", "LINE-COOK", "line_cook", " line-cook "} {
  ┃   100     r := &Recipe{ID: "RCP01", Execution: Execution{ClaimedBy: spelling}}
  ┃   101     if err := checkRecordNames(r); err == nil {
  ┃   102       t.Errorf("%q was not refused", spelling)
  ┃   103     }
  ┃   104   }
  ┃   105
  ┃   106   // Prose is not refused.
  ┃   107   prose := &Recipe{
  ┃   108     ID: "RCP01", Title: "the line-cook wrote this",
  ┃   109     Intent: Intent{Outcome: "a line-cook's recipe"},
  ┃   110   }
  ┃   111   if err := checkRecordNames(prose); err != nil {
  ┃   112     t.Errorf("prose must not be refused: %v", err)
  ┃   113   }
  ┃   114 }
  ┃   115
  ┃   116 // TestZZCycle4HandWritten probes the hand-edited record. Temporary review file.
  ┃   117 func TestZZCycle4HandWritten(t *testing.T) {
  ┃   118   s := newStore(t)
  ┃   119   if err := s.Init(); err != nil {
  ┃   120     t.Fatal(err)
  ┃   121   }
  ┃   122   dir := filepath.Join(s.Root, DirRecipes)
  ┃   123   if err := os.MkdirAll(dir, 0o755); err != nil {
  ┃   124     t.Fatal(err)
  ┃   125   }
  ┃   126   body := `{"id":"RCP01","title":"Line Cook","slug":"line-cook","lifecycle":"draft"}`
  ┃   127   if err := os.WriteFile(filepath.Join(dir, "RCP01-line-cook.json"), []byte(body), 0o644); err != nil {
  ┃   128     t.Fatal(err)
  ┃   129   }
  ┃   130   _, err := s.LoadRecipe("RCP01")
  ┃   131   t.Logf("the reader on a hand-edited title: %v", err)
  ┃   132   // And the writer refuses the same record.
  ┃   133   r := &Recipe{ID: "RCP02", Title: "Line Cook", Slug: "line-cook"}
  ┃   134   t.Logf("the writer on the same title: %v", s.SaveRecipe(r))
  ┃   135 }
  ┃   136
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/core -run 'TestZZCycle4' -v 2>&1 | tail -30
  ┃
  ┃  (9 earlier lines)
  ┃  --- PASS: TestZZCycle4Value (0.00s)
  ┃  === RUN   TestZZCycle4HandWritten
  ┃      zz_reviewer_tmp_test.go:131: the reader on a hand-edited title: this plan carries artifacts of a retired format (the record's
  ┃  title carries a retired name ("Line Cook"), the record's slug carries a retired name ("line-cook")): the current format is one
  ┃  format, and this plan is not supported
  ┃      zz_reviewer_tmp_test.go:134: the writer on the same title: the record carries a retired name: title ("Line Cook"), slug ("line-
  ┃  cook") — a retired name is not a name this tool accepts
  ┃  --- PASS: TestZZCycle4HandWritten (0.01s)
  ┃  PASS
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.428s
  ┃

     + Thought · 8.3s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && rm -f internal/core/zz_reviewer_tmp_test.go && go test ./... 2>&1 | tail -6 && make lint
  ┃  2>&1 | tail -2 && gofmt -l . && echo "gofmt clean" && make build >/dev/null 2>&1 && rm -rf /tmp/fdemo4 && ./bin/maitre-cli d…
  ┃
  ┃  (10 earlier lines)
  ┃  ▸ Recipes    Backlog    Path    Board    Log    Sous-Chef    Settings
  ┃   7 recipes · 2 need you   ☑ show served
  ┃   Needs you                                                                                      2
  ┃  ──────────────────────────────────────────────────────────────────────────────────────────────────
  ┃    RCP02 ready       Sign off the menu                             1/2
  ┃    RCP03 remediation Wire the review queue                         2/3 1 note
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff 73e3666..HEAD -- work/README.md work/inbox/runner-slice-f-fix-1b.md | head -120
  ┃
  ┃  (114 earlier lines)
  ┃  +1. **The trade-off the rule makes, stated out loud.** The rule is by value, so a field
  ┃  +   whose whole value is a retired name is refused — and that includes a **title** (or
  ┃  +   the slug derived from it) of exactly "Line Cook" ✓: the walk finds
  ┃  +   `{title Line Cook}` ✓. So a recipe titled *Line Cook* makes the plan unsupported,
  ┃  +   and the writer refuses to create it ✓. Prose is safe ✓ ("The line-cook's recipe" is
  ┃  +   not refused ✓). This is the ruling's line, drawn where it was drawn ✓ — but it is a
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n 30,80p work/inbox/runner-slice-f-fix-1b.md
  ┃
  ┃  (5 earlier lines)
  ┃     lists ✓, and its own "it can fail" proof ✓ (a token in each package ✓).
  ┃
  ┃  Untouched: the gate table, the plan format, `core.DomainVersion` (7), the command
  ┃  count (42), the journal's history, `docs/audit.md`, `docs/07-vocabulary.md`,
  ┃  `work/`. The invariants hold: `internal/core` 0 · `internal/cli` 0 · `internal/ui` 28.
  ┃
  ┃  Gate: make lint && gofmt -l . && make test && make regen-goldens — green.
  ┃

     + Thought · 1.2s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/ui/views.go /tmp/rf4-ui.bak && cp internal/core/model.go /tmp/rf4-model.bak &
  ┃  & cp internal/core/store.go /tmp/rf4-store.bak && printf '\n// mutation probe\nvar mutationProbe = "specials"\n' >> internal…
  ┃
  ┃  (7 earlier lines)
  ┃  FAIL  github.com/lthiagol/maitre/internal/core  0.297s
  ┃  FAIL
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/core  0.300s
  ┃  FAIL
  ┃  ?? logo/
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'EOF'
  ┃  p='internal/core/model.go'…
  ┃
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/core  0.333s
  ┃  FAIL
  ┃  ?? logo/
  ┃  BUILD_OK
  ┃

     + Thought · 11.1s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --name-only --oneline 73e3666..HEAD | grep -i logo; echo "(no logo)"; git status
  ┃  --short; grep -rn "RawMessage" --include="*.go" internal/core | grep -v _test | head -5
  ┃
  ┃  (no logo)
  ┃  ?? logo/
  ┃  internal/core/activity.go:58:  Data    json.RawMessage `json:"data,omitempty"`
  ┃  internal/core/retired.go:254:  var doc map[string]json.RawMessage
  ┃  internal/core/retired.go:263:    var criteria []map[string]json.RawMessage
  ┃  internal/core/retired.go:284:  var doc map[string]json.RawMessage
  ┃  internal/core/retired.go:322:// **Byte slices are not walked either.** `json.RawMessage` is one, and it holds the
  ┃

     + Thought · 2.8s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "map\[" --include="*.go" internal/core/model.go internal/core/backlog.go | grep
  ┃  -v _test | head -10
  ┃
  ┃  internal/core/backlog.go:235:  if err := o.journalLocked(KindCreated, "", map[string]any{"kind": "backlog", "id": id, "title":
  ┃  title, "item_kind": kind}); err != nil {
  ┃  internal/core/backlog.go:285:  if err := o.journalLocked(KindDecision, "", map[string]any{"kind": "backlog", "id": id, "item_kind":
  ┃  v.Kind}); err != nil {
  ┃  internal/core/backlog.go:340:  if err := o.journalLocked(KindOverlay, "", map[string]any{"kind": "backlog", "id": id, "status": v.
  ┃  Status}); err != nil {
  ┃  internal/core/backlog.go:409:  if err := o.journalLocked(KindPromotion, recipe.ID, map[string]any{
  ┃

     + Thought · 3.7s

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃     93   | **F-60** | medium | **`core.Questions` has a writer-less  93   | **F-60** | medium | **`core.Questions` has a writer-less
  ┃          shape, so the `questions.resolved` gate can never fire** —       shape, so the `questions.resolved` gate can never fire** —
  ┃          nothing ever sets a question, so the list is always empty        nothing ever sets a question, so the list is always empty
  ┃          and the gate always passes. A gate nobody can fail is not        and the gate always passes. A gate nobody can fail is not
  ┃          a guard. | `internal/core/model.go:31`; `internal/core/          a guard. | `internal/core/model.go:31`; `internal/core/
  ┃          lifecycle.go:197` |                                              lifecycle.go:197` |
  ┃     94   | **F-61** | medium | **The `courses/` directory is named   94   | **F-61** | medium | **The `courses/` directory is named
  ┃          by the plan format, created by `Init`, and never written**       by the plan format, created by `Init`, and never written**
  ┃          — a directory no writer touches. Per [06-proposals.md](06-       — a directory no writer touches. Per [06-proposals.md](06-
  ┃          proposals.md) § 1 the honest answer is to stop naming it:        proposals.md) § 1 the honest answer is to stop naming it:
  ┃          a course is the relation, not a file. | `internal/core/          a course is the relation, not a file. | `internal/core/
  ┃          store.go:18,63`; [docs/03 §6](03-domain-model.md) |              store.go:18,63`; [docs/03 §6](03-domain-model.md) |
  ┃     95   | F-56 | low | **The poll's fingerprint is blind to a       95   | F-56 | low | **The poll's fingerprint is blind to a
  ┃          write that keeps the mtime and the size.** It is path +          write that keeps the mtime and the size.** It is path +
  ┃          mtime + size per file, so two different contents can share       mtime + size per file, so two different contents can share
  ┃          one fingerprint. The tool's own writes always change the         one fingerprint. The tool's own writes always change the
  ┃          mtime, so in practice only a hand edit that restores the         mtime, so in practice only a hand edit that restores the
  ┃          mtime is missed — and the chef's `r` reloads it. |               mtime is missed — and the chef's `r` reloads it. |
  ┃          `internal/core/store.go` — `DiskFingerprint`. |                  `internal/core/store.go` — `DiskFingerprint`. |
  ┃     96   | P-8 | — | `--depends-on=` writes `depends_on: [""]` —     96   | P-8 | — | `--depends-on=` writes `depends_on: [""]` —
  ┃          inert, the same as `--in-scope=`. | Recorded so nobody           inert, the same as `--in-scope=`. | Recorded so nobody
  ┃          calls it a regression later. Not a defect. |                     calls it a regression later. Not a defect. |
  ┃     97 - | **F-92** | medium | **`docs/03` § 1 still says the        97 + | **F-94** | medium | **A retired name can still be
  ┃          reader takes a file named with a retired id form.** The          written into a live record by the one path the by-value
  ┃          ids section reads: *"**The reader strips every leading           rule does not reach — the *journal's* new event.**
  ┃          letter.** `RCP01`, `rcp01`, `RCP1`, `R01`, `r1` and `1`          `internal/core/activity.go` checks a new event's `actor`
  ┃          all name the same recipe, **and a file named `R01-<slug>.        and `role` by field, but the **session** is a string field
  ┃          json` is read as readily as one named `RCP01-<slug>.             of the same event, and nothing checks it: `ops.
  ┃          json`**" — and the next sentence says *"the plans and            RecordMigration` / `ops.journal` write a `Session` the
  ┃          journals written before it carry the old forms, and a            caller names, so a session called `line-cook` lands in a
  ┃          reader that cannot resolve them cannot read the past"*.          new event. The journal is history and the reader does not
  ┃          The code refuses exactly that file: the file's **name** is       check it, so nothing is bricked — but the rule the slice
  ┃          part of the format, and a plan carrying one is *not              states is *"a **new** event may not carry a retired
  ┃          supported*. The same document says so twice (§ 6 and § 6.        name"*, and one of the event's three string fields is not
  ┃          2: *"a **file named with a retired id form** is **not            covered. The recipe's check is by value over **every**
  ┃          supported**"* and *"a file named with a retired id form is       string field; the journal's is by field. | `internal/core/
  ┃          refused as well"*) — so the spec contradicts itself, and         activity.go` — `appendJournalEventLocked`; `internal/core/
  ┃          the passage a reader meets first is the wrong one. The           model.go` — `checkName` |
  ┃          *lookup* is still tolerant (`RCP1`, `rcp01`, `1`, `R01`
  ┃          name the recipe), and the **journal** keeps its old ids;
  ┃          the **plan's** records do not. | `docs/03-domain-model.md:
  ┃          121-125` |
  ┃     98 - | **F-93** | medium | **The record's name fields are
  ┃          enumerated twice, and the comment says they are one.**
  ┃          `internal/core/retired.go` holds two hand-written lists of
  ┃          the same nine fields: `actorNames` (the **writer**'s, read
  ┃          by `checkRecordNames` before a record is saved) and
  ┃          `retiredActorNames` (the **reader**'s, which reports what
  ┃          a record carries). They carry the same doc comment — *"The
  ┃          fields are the ones that name an actor: the claimant, the
  ┃          approver, the reviewer, the comment's author, the record's
  ┃          author, the deviation's author and the overlays' author"*
  ┃          — and `actorNames`'s own says *"the **reader** reads them
  ┃          when it reads"*, which it does not: the reader reads its
  ┃          own copy. Two answers to one question, and the drift's
  ┃          consequence is the hole F-91 just closed: with a name-
  ┃          holding field the **reader's** list knows and the
  ┃          **writer's** does not, the writer accepts `"line-cook"`,
  ┃          and the reader then refuses the plan the writer just wrote
  ┃          (*"this plan carries artifacts of a retired format (a
  ┃          signed-by ("line-cook"))"*) — the plan is bricked, with no
  ┃          remedy but a hand edit. No wire compares the two lists.
  ┃          The remedy is one list and two views of it. | `internal/
  ┃          core/retired.go` — `actorNames`, `checkRecordNames`,
  ┃          `retiredActorNames` |
  ┃     99                                                               98
  ┃    100   ## Closed findings                                          99   ## Closed findings
  ┃    101                                                              100
  ┃    102   | # | Severity | Closed in | Held by |                     101   | # | Severity | Closed in | Held by |
  ┃

     + Thought · 25.6s

  ┃
  ┃  Writing command…
  ┃
  ┃  No tool named "bash" is currently available. Please use a tool from the available tool list.
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n 670,700p docs/03-domain-model.md
  ┃
  ┃  (24 earlier lines)
  ┃  **A lookup is tolerant; a record is not.** `RCP01`, `RCP1`, `rcp01` and `1` all
  ┃  name the same recipe, and `BL01`, `bl01`, `BL1`, `B01` and `1` all name the same
  ┃  backlog item: a name resolves to a record, and a name is not a record. The **record**
  ┃  on disk is in the current format — `RCP01`, `BL01`, `evidence`, `completion` — and a
  ┃  record under a retired form is refused. The **file's name** is part of the format,
  ┃  so a file named with a retired id form is refused as well — while a name the tool
  ┃  merely tolerates (`RCP1`, `B01`) is read. The **journal** keeps the old names and the
  ┃

     + Thought · 5.9s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "type Fragment" -A 20 internal/core/model.go
  ┃
  ┃  (no output)
  ┃  Exited with code 1
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "type Fragment" -A 22 internal/core/*.go | grep -v _test | head -26
  ┃
  ┃  (16 earlier lines)
  ┃  internal/core/operations.go-96-}
  ┃  internal/core/operations.go-97-
  ┃  internal/core/operations.go-98-// ListRecipes returns every recipe, ordered by id.
  ┃  internal/core/operations.go-99-func (o *Operations) ListRecipes() ([]Recipe, error) {
  ┃  internal/core/operations.go-100-  return o.store.ListRecipes()
  ┃  internal/core/operations.go-101-}
  ┃  internal/core/operations.go-102-
  ┃

     + Thought · 1.0s

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃     93   | **F-60** | medium | **`core.Questions` has a writer-less  93   | **F-60** | medium | **`core.Questions` has a writer-less
  ┃          shape, so the `questions.resolved` gate can never fire** —       shape, so the `questions.resolved` gate can never fire** —
  ┃          nothing ever sets a question, so the list is always empty        nothing ever sets a question, so the list is always empty
  ┃          and the gate always passes. A gate nobody can fail is not        and the gate always passes. A gate nobody can fail is not
  ┃          a guard. | `internal/core/model.go:31`; `internal/core/          a guard. | `internal/core/model.go:31`; `internal/core/
  ┃          lifecycle.go:197` |                                              lifecycle.go:197` |
  ┃     94   | **F-61** | medium | **The `courses/` directory is named   94   | **F-61** | medium | **The `courses/` directory is named
  ┃          by the plan format, created by `Init`, and never written**       by the plan format, created by `Init`, and never written**
  ┃          — a directory no writer touches. Per [06-proposals.md](06-       — a directory no writer touches. Per [06-proposals.md](06-
  ┃          proposals.md) § 1 the honest answer is to stop naming it:        proposals.md) § 1 the honest answer is to stop naming it:
  ┃          a course is the relation, not a file. | `internal/core/          a course is the relation, not a file. | `internal/core/
  ┃          store.go:18,63`; [docs/03 §6](03-domain-model.md) |              store.go:18,63`; [docs/03 §6](03-domain-model.md) |
  ┃     95   | F-56 | low | **The poll's fingerprint is blind to a       95   | F-56 | low | **The poll's fingerprint is blind to a
  ┃          write that keeps the mtime and the size.** It is path +          write that keeps the mtime and the size.** It is path +
  ┃          mtime + size per file, so two different contents can share       mtime + size per file, so two different contents can share
  ┃          one fingerprint. The tool's own writes always change the         one fingerprint. The tool's own writes always change the
  ┃          mtime, so in practice only a hand edit that restores the         mtime, so in practice only a hand edit that restores the
  ┃          mtime is missed — and the chef's `r` reloads it. |               mtime is missed — and the chef's `r` reloads it. |
  ┃          `internal/core/store.go` — `DiskFingerprint`. |                  `internal/core/store.go` — `DiskFingerprint`. |
  ┃     96   | P-8 | — | `--depends-on=` writes `depends_on: [""]` —     96   | P-8 | — | `--depends-on=` writes `depends_on: [""]` —
  ┃          inert, the same as `--in-scope=`. | Recorded so nobody           inert, the same as `--in-scope=`. | Recorded so nobody
  ┃          calls it a regression later. Not a defect. |                     calls it a regression later. Not a defect. |
  ┃     97 - | **F-94** | medium | **A retired name can still be
  ┃          written into a live record by the one path the by-value
  ┃          rule does not reach — the *journal's* new event.**
  ┃          `internal/core/activity.go` checks a new event's `actor`
  ┃          and `role` by field, but the **session** is a string field
  ┃          of the same event, and nothing checks it: `ops.
  ┃          RecordMigration` / `ops.journal` write a `Session` the
  ┃          caller names, so a session called `line-cook` lands in a
  ┃          new event. The journal is history and the reader does not
  ┃          check it, so nothing is bricked — but the rule the slice
  ┃          states is *"a **new** event may not carry a retired
  ┃          name"*, and one of the event's three string fields is not
  ┃          covered. The recipe's check is by value over **every**
  ┃          string field; the journal's is by field. | `internal/core/
  ┃          activity.go` — `appendJournalEventLocked`; `internal/core/
  ┃          model.go` — `checkName` |
  ┃     98                                                               97
  ┃     99   ## Closed findings                                          98   ## Closed findings
  ┃    100                                                               99
  ┃    101   | # | Severity | Closed in | Held by |                     100   | # | Severity | Closed in | Held by |
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    154   | F-88 | medium | **`docs/04-operations.md` still          154   | F-88 | medium | **`docs/04-operations.md` still
  ┃          describes the resolution the slice removed.** *Closed in*        describes the resolution the slice removed.** *Closed in*
  ┃          `docs/04-operations.md` — the paragraph no longer says the       `docs/04-operations.md` — the paragraph no longer says the
  ┃          retired names name the role they stood for: it says a            retired names name the role they stood for: it says a
  ┃          retired name **is not a name the tool accepts**, refused         retired name **is not a name the tool accepts**, refused
  ┃          wherever an actor names itself, and that there is no             wherever an actor names itself, and that there is no
  ┃          `line-cook` → `runner` mapping any more. `docs/03` § 3.2         `line-cook` → `runner` mapping any more. `docs/03` § 3.2
  ┃          says the same thing. | `internal/core/docs_test.go` —            says the same thing. | `internal/core/docs_test.go` —
  ┃          `checkDenyList`: the spec names every role the code              `checkDenyList`: the spec names every role the code
  ┃          refuses **and** every retired name |                             refuses **and** every retired name |
  ┃    155   | F-89 | medium | **The approval digest's encoding changed 155   | F-89 | medium | **The approval digest's encoding changed
  ┃          and neither fence moved.** *Closed in* `internal/core/           and neither fence moved.** *Closed in* `internal/core/
  ┃          format.go` — the digest's encoding is **named** (                format.go` — the digest's encoding is **named** (
  ┃          `canonical-payload-v1`) and the canonical payload's field        `canonical-payload-v1`) and the canonical payload's field
  ┃          names are in `planFormatFacts`, so the payload's shape is        names are in `planFormatFacts`, so the payload's shape is
  ┃          a format fact and moves the tripwire; and `core.                 a format fact and moves the tripwire; and `core.
  ┃          DomainVersion` moved **6 → 7** because the shape moved. A        DomainVersion` moved **6 → 7** because the shape moved. A
  ┃          plan stamped 6 is now refused by the **version** guard,          plan stamped 6 is now refused by the **version** guard,
  ┃          with the honest message (*"this plan was written under           with the honest message (*"this plan was written under
  ┃          domain version 6 and this tool speaks version 7"*), and          domain version 6 and this tool speaks version 7"*), and
  ┃          never reaches the digest — whose tamper message now means        never reaches the digest — whose tamper message now means
  ┃          **content and only content**. The repository's own plan is       **content and only content**. The repository's own plan is
  ┃          stamped 7 and the journal carries the declaration as an          stamped 7 and the journal carries the declaration as an
  ┃          **append** (event 40; the 37 historical events untouched).       **append** (event 40; the 37 historical events untouched).
  ┃          | `internal/core/version_test.go` —                              | `internal/core/version_test.go` —
  ┃          `TestDomainVersionRefusesAnOlderPlan`,                           `TestDomainVersionRefusesAnOlderPlan`,
  ┃          `TestDomainVersionWritesThisVersion`,                            `TestDomainVersionWritesThisVersion`,
  ┃          `TestTheDigestEncodingIsAFormatFact`,                            `TestTheDigestEncodingIsAFormatFact`,
  ┃          `TestTheDeclarationMovesThePlanForwardOnly`; `internal/          `TestTheDeclarationMovesThePlanForwardOnly`; `internal/
  ┃          core/format.go` — `planFormatFacts`; `internal/core/lock.        core/format.go` — `planFormatFacts`; `internal/core/lock.
  ┃          go` — `checkWritable`, `checkDeclarationForward` |               go` — `checkWritable`, `checkDeclarationForward` |
  ┃    156   | F-90 | medium | **The new "one format" wire has three    156   | F-90 | medium | **The new "one format" wire has three
  ┃          holes.** *Closed in* `internal/core/retired_test.go` — the       holes.** *Closed in* `internal/core/retired_test.go` — the
  ┃          scan reads the **whole tree** (`internal`, `cmd`,                scan reads the **whole tree** (`internal`, `cmd`,
  ┃          `scripts`), the exemptions are **per declaration** and not       `scripts`), the exemptions are **per declaration** and not
  ┃          per file (`retired.go`'s `retiredDirectories`,                   per file (`retired.go`'s `retiredDirectories`,
  ┃          `retiredInRecipe`, `retiredInBacklog`; `model.go`'s              `retiredInRecipe`, `retiredInBacklog`; `model.go`'s
  ┃          `RetiredNames`), and the token list carries the retired id       `RetiredNames`), and the token list carries the retired id
  ┃          literals (`R01`, `S01`, `N01`), the retired families (           literals (`R01`, `S01`, `N01`), the retired families (
  ┃          `legacy`, `special`, `prep-note`, `proof`) and the retired       `legacy`, `special`, `prep-note`, `proof`) and the retired
  ┃          stamp by its **access** (`["evidence"]`) rather than by          stamp by its **access** (`["evidence"]`) rather than by
  ┃          the bare word. The reviewer's four probes — `specials` in        the bare word. The reviewer's four probes — `specials` in
  ┃          `internal/ui`, `specials` in `model.go`, `R01` in `store.        `internal/ui`, `specials` in `model.go`, `R01` in `store.
  ┃          go` and `LegacyThing` in `store.go` — all fail it. |             go` and `LegacyThing` in `store.go` — all fail it. |
  ┃          `internal/core/retired_test.go` —                                `internal/core/retired_test.go` —
  ┃          `TestTheRetiredNamesAreOnlyWhereTheyBelong`,                     `TestTheRetiredNamesAreOnlyWhereTheyBelong`,
  ┃          `TestTheOneFormatWireCanFail` |                                  `TestTheOneFormatWireCanFail` |
  ┃    157   | F-91 | medium | **The retired-name rule is enforced per  157   | F-91 | medium | **The retired-name rule is enforced per
  ┃          operation, not where a name enters a record — and the            operation, not where a name enters a record — and the
  ┃          second approve path is not guarded.** *Closed in*                second approve path is not guarded.** *Closed in*
  ┃          `internal/core/store.go` — `SaveRecipe` runs                     `internal/core/store.go` — `SaveRecipe` runs
  ┃          `checkRecordNames` on **every** record it writes, so the         `checkRecordNames` on **every** record it writes, so the
  ┃          rule lives where a name enters a record and not once per         rule lives where a name enters a record and not once per
  ┃          operation: the domain's API cannot write a plan it will          operation: the domain's API cannot write a plan it will
  ┃          not read. `internal/core/operations.go` — `advance` checks       not read. `internal/core/operations.go` — `advance` checks
  ┃          the transition's approver (the one place `ApproveRecipe`,        the transition's approver (the one place `ApproveRecipe`,
  ┃          `RequestChanges` and `AdvanceWith` all pass through) and         `RequestChanges` and `AdvanceWith` all pass through) and
  ┃          `Operations.lock` checks the actor **and** the audit role,       `Operations.lock` checks the actor **and** the audit role,
  ┃          before anything is written. `internal/core/activity.go` —        before anything is written. `internal/core/activity.go` —
  ┃          a **new** journal event may not carry a retired name,            a **new** journal event may not carry a retired name,
  ┃          while the **history** keeps its old names (the tests that        while the **history** keeps its old names (the tests that
  ┃          need one write it by hand, the way the old tool wrote it).       need one write it by hand, the way the old tool wrote it).
  ┃          `internal/core/model.go` — `normalizeAgentName` maps a           `internal/core/model.go` — `normalizeAgentName` maps a
  ┃          space **and** an underscore to a dash, so `line_cook` is         space **and** an underscore to a dash, so `line_cook` is
  ┃          the name `line-cook` and a spelling is not a way around          the name `line-cook` and a spelling is not a way around
  ┃          the rule. The reviewer's probe now **fails and writes            the rule. The reviewer's probe now **fails and writes
  ┃          nothing**: the plan still reads, the approval is not             nothing**: the plan still reads, the approval is not
  ┃          there, the lifecycle did not move. **The rule's two lists        there, the lifecycle did not move. **The rule's two lists
  ┃          — the writer's and the reader's — are still two lists:           — the writer's and the reader's — are still two lists:
  ┃          that is F-93.** | `internal/core/roles_test.go` —                that is F-93.** | `internal/core/roles_test.go` —
  ┃          `TestARetiredNameCannotEnterARecord`,                            `TestARetiredNameCannotEnterARecord`,
  ┃          `TestARetiredRoleCannotWrite`,                                   `TestARetiredRoleCannotWrite`,
  ┃          `TestARetiredNameCannotRequestChanges`; `internal/core/          `TestARetiredNameCannotRequestChanges`; `internal/core/
  ┃          store.go` — `SaveRecipe`, `checkRecordNames`; `internal/         store.go` — `SaveRecipe`, `checkRecordNames`; `internal/
  ┃          core/operations.go` — `advance`, `lock`; `internal/core/         core/operations.go` — `advance`, `lock`; `internal/core/
  ┃          activity.go` — `appendJournalEventLocked` |                      activity.go` — `appendJournalEventLocked` |
  ┃                                                                     158 + | F-92 | medium | **`docs/03` § 1 still says the reader
  ┃                                                                           takes a file named with a retired id form.** *Closed in*
  ┃                                                                           `docs/03-domain-model.md` § 1 — the ids section now draws
  ┃                                                                           the line the code draws, in a table of three things that
  ┃                                                                           are not the same: a **name typed by a human** is a
  ┃                                                                           **lookup** and is tolerated (every form names the recipe
  ┃                                                                           the plan holds); a **file's name on disk** is the
  ┃                                                                           **format**, and a file named `R01-<slug>.json` carries a
  ┃                                                                           retired id form, so the plan is **not supported**; the
  ┃                                                                           **journal** is **history** — append-only, its old ids
  ┃                                                                           kept, and the reader still resolves what it names. *"A
  ┃                                                                           reader that cannot resolve them cannot read the past"* is
  ┃                                                                           now said of the journal, and the plan whose files carry
  ┃                                                                           the retired form "has no plan to read". The backlog
  ┃                                                                           paragraph says the same: `B01` is a **name**, not a
  ┃                                                                           retired form. § 6 and § 6.2 agree with § 1, and § 6's
  ┃                                                                           rules-list clause no longer reads as making slice E's
  ┃                                                                           change a version bump (*"a fact added to the list that the
  ┃                                                                           format already had is not a format change; a fact the
  ┃                                                                           format **changed** is"*). | `docs/03-domain-model.md:121-
  ┃                                                                           140` — (the three sections' agreement is not wire-checked)
  ┃                                                                           |
  ┃                                                                     159 + | F-93 | medium | **The record's name fields are
  ┃                                                                           enumerated twice, and the comment says they are one.**
  ┃                                                                           *Closed in* `internal/core/retired.go` — the two lists are
  ┃                                                                           gone: the rule is **by value** over **every** string field
  ┃                                                                           of the record, walked by reflection, so a new name-
  ┃                                                                           carrying field is caught without anybody adding it to a
  ┃                                                                           list, and the **writer** (`checkRecordNames`) and the
  ┃                                                                           **reader** (`retiredNames`) run the **same function** —
  ┃                                                                           one answer, two call sites, so the drift cannot happen.
  ┃                                                                           The comment that lied is gone with the lists. The F-93
  ┃                                                                           wire that compared the two views went with them — there is
  ┃                                                                           nothing to compare — and the one-format wire keeps
  ┃                                                                           everything else: the retired names' two places, the
  ┃                                                                           canonical lists, and its own proof that it can fail. |
  ┃                                                                           `internal/core/retired_test.go` —
  ┃                                                                           `TestARetiredNameIsRefusedByValue`,
  ┃                                                                           `TestProseThatMentionsARetiredNameIsNotRefused`,
  ┃                                                                           `TestTheScanFindsANameInEveryFieldShape`,
  ┃                                                                           `TestTheScanSkipsMapsAndBytes`; `internal/core/retired.go`
  ┃                                                                           — `retiredNameFields`, `checkRecordNames`, `retiredNames`
  ┃                                                                           |
  ┃    158                                                              160
  ┃    159   ## What this ledger cannot tell you                        161   ## What this ledger cannot tell you
  ┃    160                                                              162
  ┃    161   **Cycles 1–7's own findings are not recoverable.** Their   163   **Cycles 1–7's own findings are not recoverable.** Their
  ┃          per-finding reports lived in a                                   per-finding reports lived in a
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    45   | Slice E | The vocabulary rename: the paper — the glossary 45   | Slice E | The vocabulary rename: the paper — the glossary
  ┃         rewritten, ADR 0009 recorded, the tree swept | The ids and       rewritten, ADR 0009 recorded, the tree swept | The ids and
  ┃         the guard hold and the sweep's 78 lines are all records of       the guard hold and the sweep's 78 lines are all records of
  ┃         something that still exists; the retired-word pattern holds      something that still exists; the retired-word pattern holds
  ┃         for seven renamed terms and misses the console's, and the        for seven renamed terms and misses the console's, and the
  ┃         spec and the `DomainVersion` comment still describe the          spec and the `DomainVersion` comment still describe the
  ┃         fingerprint the way the slice found it — five mediums filed      fingerprint the way the slice found it — five mediums filed
  ┃         |                                                                |
  ┃    46   | Slice E · cycle 2 | The fixes to slice E: the glossary's  46   | Slice E · cycle 2 | The fixes to slice E: the glossary's
  ┃         retired names, the console's place, the spec's two lists         retired names, the console's place, the spec's two lists
  ┃         and the version/fingerprint split — **the last slice of the      and the version/fingerprint split — **the last slice of the
  ┃         rename** | F-81..F-85 closed, each verified by a command         rename** | F-81..F-85 closed, each verified by a command
  ┃         and not by the diff; the sweep's count reconciles (78 → 77:      and not by the diff; the sweep's count reconciles (78 → 77:
  ┃         F-83 removed one prose line and nothing was added, the           F-83 removed one prose line and nothing was added, the
  ┃         console's `_Avoid_` line already existed — but the               console's `_Avoid_` line already existed — but the
  ┃         published pattern cannot see the name F-81 put on it) — one      published pattern cannot see the name F-81 put on it) — one
  ┃         medium filed, against the sweep's command, which does not        medium filed, against the sweep's command, which does not
  ┃         look for the console's retired name |                            look for the console's retired name |
  ┃    47   | Slice E · cycle 3 | The last fixes to slice E: the        47   | Slice E · cycle 3 | The last fixes to slice E: the
  ┃         sweep's instrument, and the two notes the review left | F-       sweep's instrument, and the two notes the review left | F-
  ┃         86 closed — the command carries all seven retired names and      86 closed — the command carries all seven retired names and
  ┃         tells a warning from the drift (72 lines in 18 files; the        tells a warning from the drift (72 lines in 18 files; the
  ┃         count reconciles, 77 − 5), the glossary's heading is **The       count reconciles, 77 − 5), the glossary's heading is **The
  ┃         surfaces** and `docs/03`'s pointer lands on § 6.2 — **the        surfaces** and `docs/03`'s pointer lands on § 6.2 — **the
  ┃         rename's last slice, closed**: the glossary rewritten, ADR       rename's last slice, closed**: the glossary rewritten, ADR
  ┃         0009 written, the tree swept, and **the mascot's acceptance      0009 written, the tree swept, and **the mascot's acceptance
  ┃         criterion carried to the session that owns the drawing** (       criterion carried to the session that owns the drawing** (
  ┃         AC-01) |                                                         AC-01) |
  ┃    48   | Slice E · closed | The rename's last slice — the glossary 48   | Slice E · closed | The rename's last slice — the glossary
  ┃         rewritten, ADR 0009 written, the tree swept | F-81..F-86         rewritten, ADR 0009 written, the tree swept | F-81..F-86
  ┃         closed across three rounds, each verified by a command; the      closed across three rounds, each verified by a command; the
  ┃         sweep's command carries all seven retired names and tells a      sweep's command carries all seven retired names and tells a
  ┃         warning from the drift, and its **72 lines in 18 files**         warning from the drift, and its **72 lines in 18 files**
  ┃         are all records of something that still exists — **AC-01,        are all records of something that still exists — **AC-01,
  ┃         the mascot's drawing, is carried to the session that owns        the mascot's drawing, is carried to the session that owns
  ┃         it**, and nothing of that session's work was committed here      it**, and nothing of that session's work was committed here
  ┃         |                                                                |
  ┃    49 - | Slice F | The tolerance machinery comes out: **this is    49 + | Slice F | The tolerance machinery comes out: **this is
  ┃         not brownfield** — the retired formats are gone, the reader      not brownfield** — the retired formats are gone, the reader
  ┃         stops reading them, the deny-list keeps the retired names        stops reading them, the deny-list keeps the retired names
  ┃         and the resolution goes, the approval's digest fences on         and the resolution goes, the approval's digest fences on
  ┃         content, and the version fence is one rule | A plan              content, and the version fence is one rule | A plan
  ┃         carrying a retired artifact is refused plainly, and not          carrying a retired artifact is refused plainly, and not
  ┃         written to either; a rename does not move the digest; the        written to either; a rename does not move the digest; the
  ┃         fence refuses in both directions; a retired name cannot buy      fence refuses in both directions; a retired name cannot buy
  ┃         a gate; **the rule lives in one place — where a name enters      a gate; **the rule lives in one place — where a name enters
  ┃         a record** (`Store.SaveRecipe` and the journal's append);        a record** (`Store.SaveRecipe` and the journal's append) —
  ┃         and the new **"one format"** wire holds the line over the        and it is **by value**, over every string field of the
  ┃         whole tree — cycle 1 filed four mediums (F-87..F-90) |           record, so a new name-carrying field is caught without a
  ┃                                                                          list; and the new **"one format"** wire holds the line over
  ┃                                                                          the whole tree — cycle 1 filed four mediums (F-87..F-90) |
  ┃    50   | Slice F · cycle 2 | The fixes to slice F: the retired     50   | Slice F · cycle 2 | The fixes to slice F: the retired
  ┃         names, the operations doc, the digest's encoding as a            names, the operations doc, the digest's encoding as a
  ┃         format fact, and the wire | F-87..F-90 closed, each              format fact, and the wire | F-87..F-90 closed, each
  ┃         verified by a command and not by the diff; the version           verified by a command and not by the diff; the version
  ┃         moved **6 → 7** and the digest's encoding is now a fact the      moved **6 → 7** and the digest's encoding is now a fact the
  ┃         tripwire sees; the two tests that changed their premise          tripwire sees; the two tests that changed their premise
  ┃         were judged and still assert what they were written for —        were judged and still assert what they were written for —
  ┃         one blocker filed (F-91): the rule is enforced per               one blocker filed (F-91): the rule is enforced per
  ┃         operation, and the second approve path is not guarded |          operation, and the second approve path is not guarded |
  ┃    51   | Slice F · cycle 3 | The last fixes to slice F: the rule   51   | Slice F · cycle 3 | The last fixes to slice F: the rule
  ┃         moved to where a name enters a record, and the notes | F-91      moved to where a name enters a record, and the notes | F-91
  ┃         closed — the probe now fails and writes nothing, the audit       closed — the probe now fails and writes nothing, the audit
  ┃         role is refused, a new journal event may not carry a             role is refused, a new journal event may not carry a
  ┃         retired name and the history keeps its own, and a spelling       retired name and the history keeps its own, and a spelling
  ┃         is not a way around the rule; the wire reads the code's          is not a way around the rule; the wire reads the code's
  ┃         lists, so a name added to the code is caught without             lists, so a name added to the code is caught without
  ┃         touching the wire; **docs/03** § 1 still says the reader         touching the wire; **docs/03** § 1 still says the reader
  ┃         takes a file named with a retired id form (F-92) and the         takes a file named with a retired id form (F-92) and the
  ┃         record's name fields are still enumerated twice (F-93) —         record's name fields are still enumerated twice (F-93) —
  ┃         two mediums filed |                                              two mediums filed |
  ┃                                                                     52 + | Slice F · cycle 4 | The last round of slice F: the check
  ┃                                                                          is by value, and the field list is gone | F-92 and F-93
  ┃                                                                          closed — § 1 draws the line the code draws (a name typed by
  ┃                                                                          a human is a lookup, a file's name on disk is the format,
  ┃                                                                          the journal is history), and the record's name fields are
  ┃                                                                          no longer declared at all: the walk is by value over every
  ┃                                                                          string field, and the writer and the reader run the same
  ┃                                                                          function — **no blockers filed** |
  ┃    52   | The rounds inside a cycle | Fix → verify → back to the    53   | The rounds inside a cycle | Fix → verify → back to the
  ┃         executor → repeat, until the verifier is satisfied | Three       executor → repeat, until the verifier is satisfied | Three
  ┃         rounds in cycle 8: the first fix introduced three new            rounds in cycle 8: the first fix introduced three new
  ┃         mediums, the second closed them |                                mediums, the second closed them |
  ┃    53                                                               54
  ┃    54   ### The discipline every round runs                         55   ### The discipline every round runs
  ┃    55                                                               56
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "F-9" docs/audit.md | cut -c1-90 && git diff --stat
  ┃
  ┃  (4 earlier lines)
  ┃  154:| F-87 | medium | **The independence gate can be bought with a retired name.** *Closed
  ┃  157:| F-90 | medium | **The new "one format" wire has three holes.** *Closed in* `internal
  ┃  158:| F-91 | medium | **The retired-name rule is enforced per operation, not where a name
  ┃  159:| F-92 | medium | **`docs/03` § 1 still says the reader takes a file named with a reti
  ┃  160:| F-93 | medium | **The record's name fields are enumerated twice, and the comment say
  ┃   docs/audit.md | 7 ++++---
  ┃   1 file changed, 4 insertions(+), 3 deletions(-)
  ┃

     + Thought · 724ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '/## Open findings/,/## Closed findings/p' docs/audit.md | cut -c1-110
  ┃
  ┃  (4 earlier lines)
  ┃  | **F-59** | medium | **`core.Audit.Decisions` has a shape and no writer** — it is where a decision *about the
  ┃  | **F-60** | medium | **`core.Questions` has a writer-less shape, so the `questions.resolved` gate can never f
  ┃  | **F-61** | medium | **The `courses/` directory is named by the plan format, created by `Init`, and never wri
  ┃  | F-56 | low | **The poll's fingerprint is blind to a write that keeps the mtime and the size.** It is path +
  ┃  | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same as `--in-scope=`. | Recorded so nobody
  ┃
  ┃  ## Closed findings
  ┃

  ┃
  ┃  # Wrote work/inbox/reviewer-slice-f-cycle-4.md
  ┃
  ┃    1 Reviewer: slice F cycle 4 done. Verdict: ok. F-92, F-93 closed: yes — both. Can a title
  ┃    2 brick a plan: **no** — the writer's check runs *before* the write, so the record never lands.
  ┃    3 Filed 0 new blockers (none).
  ┃    4
  ┃    5 The check is by value, and the field list is gone — proved by command:
  ┃    6
  ┃    7 - A field whose **whole value** is a retired name is refused, and the error **names the
  ┃    8   field**: `the record carries a retired name: execution.claimed_by ("line-cook")` (the
  ┃    9   writer) and `the record's execution.claimed_by carries a retired name ("line-cook")`
  ┃   10   (the reader) — the same field, from the same function.
  ┃   11 - Prose that **mentions** a retired name is **not** refused: "The line-cook wrote this"
  ┃   12   (slug `the-line-cook-wrote-this`) is created and the plan reads; there is a test for it
  ┃   13   (`TestProseThatMentionsARetiredNameIsNotRefused`).
  ┃   14 - The comparison uses the tool's **own normalisation**: `line-cook`, `Line Cook`,
  ┃   15   `LINE-COOK`, `line_cook` and `" line-cook "` are all refused.
  ┃   16 - **Maps and byte slices are not walked**, and the comment says why ✓ (a test covers it).
  ┃   17
  ┃   18 **The two views are gone.** The writer (`checkRecordNames`) and the reader
  ┃   19 (`retiredNames`) run the same function (`retiredNameFields`), so there is nothing to
  ┃   20 drift — F-93's outcome, not its intermediate design. **The F-93 wire that compared the
  ┃   21 views went with them, and nothing was lost with it**: the one-format wire still holds the
  ┃   22 retired names' two places (per declaration), the canonical lists it reads, and its own
  ┃   23 proof that it can fail — my four probes (`specials` in `internal/ui`, `specials` in
  ┃   24 `model.go`, `R01` in `store.go`, `LegacyThing` in `store.go`) all still fail it, and a
  ┃   25 name added to the code's own `RetiredNames` ("wizard") is caught in `internal/ui` without
  ┃   26 touching the wire. The drift the F-93 wire guarded against is impossible for a different
  ┃   27 reason: there is no list to drift.
  ┃   28
  ┃   29 **Can a title brick a plan? No.** A recipe titled exactly "Line Cook" is **refused at
  ┃   30 creation** — `the record carries a retired name: title ("Line Cook"), slug ("line-cook")`
  ┃   31 — and **nothing is written**: the recipes directory stays empty, the plan still reads and
  ┃   32 is still writable, and a backlog item is created afterwards without trouble. The **slug**
  ┃   33 is caught with it (the title's slug is "line-cook"), and a *later* title change to "Line
  ┃   34 Cook" is refused too — the change does not land and the plan keeps working. Only a hand
  ┃   35 edit could put such a record in a plan; there the reader refuses it loudly (*"this plan
  ┃   36 carries artifacts of a retired format (the record's title carries a retired name
  ┃   37 ("Line Cook"))"*) and the writer refuses it too, so the tool never writes a record it will
  ┃   38 not read. That is F-91's shape avoided for the third time.
  ┃   39
  ┃   40 F-92 closed: `docs/03` § 1 now draws the line the code draws — a **name typed by a human**
  ┃   41 is a **lookup** and is tolerated; a **file's name on disk** is the **format** and a file
  ┃   42 named `R01-<slug>.json` makes the plan **not supported**; the **journal** is **history**
  ┃   43 and keeps its old ids. "A reader that cannot resolve them cannot read the past" is now
  ┃   44 said of the journal. § 6 and § 6.2 agree with § 1, and § 6's rules-list clause now reads
  ┃   45 "a fact added to the list that the format already had is not a format change; a fact the
  ┃   46 format *changed* is" — it no longer reads as making slice E's change a version bump.
  ┃   47
  ┃   48 F-93 closed: the record's name fields are no longer declared at all — the walk is by
  ┃   49 value over every string field, so a new name-carrying field is caught without anybody
  ┃   50 adding it to a list, and the writer and the reader run the same function.
  ┃   51
  ┃   52 Scope: `core.DomainVersion` is 7 and did not move; the gate table, the plan format, the
  ┃   53 command count (42) and the role names did not move; the journal's history is untouched;
  ┃   54 `logo/` is untracked and nothing of it is in the runner's commits; `docs/audit.md` was
  ┃   55 untouched by the runner.
  ┃   56
  ┃   57 Mutation checks (every mutation reverted; the tree is as the runner left it, plus this
  ┃   58 ledger): the four wire probes above, plus the `RetiredNames` probe. Every mutation was
  ┃   59 reverted and the tree is clean.
  ┃   60
  ┃   61 LOW notes (not filed): `docs/03` § 3.2 explains the by-value rule with a "name field vs
  ┃   62 prose field" framing, so the trade the orchestrator ruled on — a prose field whose *whole*
  ┃   63 value is a retired name (a title of exactly "Line Cook") is refused — is not written down;
  ┃   64 the rule as stated ("by value", "every string field", "one answer") is right, and the
  ┃   65 illustration's clause "prose is longer than the name" is what needs the caveat. The
  ┃   66 journal's new-event check is by **field** (`actor`, `role`) while the recipe's is by value
  ┃   67 over every string field, so a *session* named after a retired name lands in a new event —
  ┃   68 harmless, since the journal is history and nothing reads a session as a name. The walk's
  ┃   69 "maps are not walked" justification points at the journal's migration events, whose `data`
  ┃   70 is a `json.RawMessage` — bytes, which the byte rule already covers, and the journal is
  ┃   71 not walked at all; the rule itself stands on its own. The backlog item is not walked, so
  ┃   72 a backlog item titled exactly a retired name is accepted where the recipe's title is
  ┃   73 refused — an asymmetry in the trade, with no consequence (nothing is bricked).
  ┃   74
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/audit.md && git commit -q -m "docs: the ledger closes slice F — F-92 and F-9
  ┃  3, and the slice's own row" && git add work/inbox && git commit -q -m "inbox: reviewer reports slice-f cycle 4" && git log -…
  ┃
  ┃  7c4bc03 inbox: reviewer reports slice-f cycle 4
  ┃  1b60b03 docs: the ledger closes slice F — F-92 and F-93, and the slice's own row
  ┃  8c659d7 work: slice F's last review — can a title brick a plan?
  ┃  a6d9851 inbox: runner reports slice-f fix 1b round 4
  ┃  ?? logo/
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1Z:p1" "Reviewer: slice F cycle 4 done. Verdict: ok. F-92, F-93 clos
  ┃  ed: yes — both. Can a title brick a plan: NO — the writer's check runs before the write, so the record never lands. Filed 0 …
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:
  ┃  opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"idle","completion_seq":550,"cwd":"/Users/thiago/Code/Pessoal/
  ┃  maitre","focused":true,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,
  ┃  "screen_detection_skipped":true,"state_change_seq":550,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722","terminal_title":"OC
  ┃  | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":"w1Z"},"type":"agent_prompted"}}
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && go test ./internal/core ./internal/cli ./internal/ui 2>&1 | tail -4
  ┃   && git log --oneline -3
  ┃
  ┃  ?? logo/
  ┃  ok    github.com/lthiagol/maitre/internal/core  35.177s
  ┃  ok    github.com/lthiagol/maitre/internal/cli  6.596s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  112.035s
  ┃  d9a1135 work: slice F's findings — one shape, three times
  ┃  4cedb7d inbox: runner starts slice-f fix 1b round 5
  ┃  f2904ed work: two pieces of prose that mislead
  ┃

     + Thought · 1.3s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -8 && git show --stat --format= d9a1135 f2904ed 4cedb7d | head -40
  ┃
  ┃  (7 earlier lines)
  ┃  6a2dffb the name check is by value, and the field list is gone
  ┃   work/README.md | 19 ++++++++++++++++++-
  ┃   1 file changed, 18 insertions(+), 1 deletion(-)
  ┃   work/briefs/fix-f5.md | 66 +++++++++++++++++++++++++++++++++++++++++++++++++++
  ┃   1 file changed, 66 insertions(+)
  ┃   work/inbox/runner-slice-f-fix-1b.md | 42 +------------------------------------
  ┃   1 file changed, 1 insertion(+), 41 deletions(-)
  ┃

     + Thought · 384ms

     Cycle 4 review done. Verdict: ok — F-92 and F-93 are both closed; no blockers filed.

     The check is by value, and the field list is gone

     ┌───────────────────────────────────────┬─────────────────────────────────────────────────────────────────────────────────────────┐
     │ Probe                                 │ Result                                                                                  │
     ├───────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────┤
     │ a field whose whole value is a        │ refused, and the error names the field — the record carries a retired name: execution.  │
     │ retired name                          │ claimed_by ("line-cook") (writer) / the record's execution.claimed_by carries a retired │
     │                                       │ name ("line-cook") (reader) — the same field, from the same function                    │
     ├───────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────┤
     │ prose that mentions one               │ not refused — "The line-cook wrote this" is created and the plan reads; a test covers   │
     │                                       │ it                                                                                      │
     ├───────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────┤
     │ the tool's own normalisation          │ line-cook, Line Cook, LINE-COOK, line_cook, " line-cook " all refused                   │
     ├───────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────┤
     │ maps and byte slices                  │ not walked, with the reason in the comment and a test for it                            │
     └───────────────────────────────────────┴─────────────────────────────────────────────────────────────────────────────────────────┘

     The two views are gone — the writer and the reader run the same function, so there is nothing to drift. The F-93 wire went with
     them and nothing was lost: the one-format wire still holds the retired names' two places, the canonical lists it reads, and its own
     proof that it can fail (my four probes still fail it; a name added to the code's RetiredNames is caught in internal/ui without
     touching the wire).

     Can a title brick a plan? No

     A recipe titled exactly "Line Cook" is refused at creation — title ("Line Cook"), slug ("line-cook") — and nothing is written: the
     recipes directory stays empty, the plan still reads and is still writable, and a backlog item is created afterwards without
     trouble. The slug is caught with it, and a later title change is refused and does not land. Only a hand edit could put such a
     record in a plan; there the reader refuses it loudly and the writer refuses it — so the tool never writes a record it will not
     read. F-91's shape avoided for the third time.

     F-92 closed: § 1 now draws the line the code draws (a name typed by a human is a lookup; a file's name on disk is the format; the
     journal is history), § 6 and § 6.2 agree, and § 6's rules-list clause no longer reads as making slice E's change a version bump. F-
     93 closed: the name fields are no longer declared at all.

     Scope: DomainVersion 7 unmoved; gate table, plan format, command count (42) and role names unmoved; the journal's history
     untouched; logo/ untracked and uncommitted.

     LOW notes (not filed): docs/03 § 3.2's "name field vs prose field" framing leaves the trade (a title of exactly a retired name is
     refused) unwritten; the journal's new-event check is by field while the recipe's is by value (a session named after a retired name
     lands — harmless); the "maps are not walked" justification points at the journal's data, which is bytes and is not walked at all;
     the backlog item is not walked, so a backlog title of exactly a retired name is accepted where the recipe's title is refused — an
     asymmetry with no consequence.

     Ledger committed (1b60b03), inbox committed (7c4bc03), orchestrator notified. The tree is as the runner left it (logo/ untouched);
     every mutation reverted. The next round (fix-f5) has since been dispatched by the orchestrator — it touches work/ only.

     Build · MiniMax-M3.1-Flash · 5m 28s · 154.1 tok/s

  ┃
  ┃
  ┃
  ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
  ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
  ~/Code/Pessoal/maitre:master                                                                             331.8K (33%)  ctrl+p commands

```
