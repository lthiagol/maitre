# The runner's lane — T3b, captured 2026-10-30

Dumped from pane `w1H:pC` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that wrote T3b: the domain answers the item's move, and the console shapes,
promotes, dismisses and creates - and two fix cycles after the review.
The review that judged it is in [tui-t3b-reviewer.md](tui-t3b-reviewer.md).

    ┃
    ┃  ## F-120 (medium) — `docs/05`'s job
    ┃  list speaks in two voices
    ┃
    ┃  **The target:** the job list says
    ┃  **one** thing.
    ┃
    ┃  **The hole.** Cycle 1's review noted,
    ┃  as a LOW note, that `docs/05`'s job
    ┃  list
    ┃  promised the console's free edit — and
    ┃  the orchestrator ruled it into a fix
    ┃  cycle.
    ┃  **That ruling was wrong**, and the
    ┃  reviewer has caught it: the fix made
    ┃  the backlog's
    ┃  row speak in **build state** (*"the
    ┃  free edit of an item's fields is **not
    ┃  built**"*),
    ┃  while the other three rows speak in
    ┃  **remit** — what the console's job
    ┃  *is*:
    ┃
    ┃  | Row | What it names, and whether the
    ┃  console does it |
    ┃  |---|---|
    ┃  | **View** | *"findings anchored to
    ┃  code"* — the domain's `Finding.Anchor`
    ┃  is **never read** by `internal/ui` |
    ┃  | **Decide** | *"accept or reject a
    ┃  scope"* and *"choose what is next"* —
    ┃  the console does neither |
    ┃  | **Run** | *"start"* and *"pause"* —
    ┃  the console does neither; the agent
    ┃  runs the work |
    ┃  | **Manage the backlog** | now says
    ┃  what is **built** — while the three
    ┃  above say what is **remit** |
    ┃
    ┃  So the table promises three things the
    ┃  console does not do, and says one thing
    ┃  it does
    ┃  — in two voices.
    ┃
    ┃  **My ruling, and it is the
    ┃  orchestrator's — the table is the
    ┃  REMIT.** The console's
    ┃  job list says **what the console is
    ┃  for**, and the console's **surface** is
    ┃  described
    ┃  by the keys table below it. That is how
    ┃  ADR-0003 puts it — *the human interface
    ┃  is a
    ┃  **management surface*** — and the
    ┃  orchestrator's ruling on F-117 (that
    ┃  the ADR's own job
    ┃  table **stays**, in the vocabulary of
    ┃  its day) says the same thing about the
    ┃  same
    ┃  table: it is a statement of the
    ┃  **division of labour**, not of the
    ┃  build.
    ┃
    ┃  **What lands:**
    ┃
    ┃  1. **One line above the table** says
    ┃  what the table is — the console's remit
    ┃  — and
    ┃     that **what is built is the keys
    ┃  table below**.
    ┃  2. **The backlog's row returns to the
    ┃  remit's wording.** The clause *"the
    ┃  free edit of
    ┃     an item's fields is not built"*
    ┃  comes **out** of `docs/05` — it is a
    ┃  **build state**,
    ┃     and this table is not where build
    ┃  states live. **It stays in the plan**,
    ┃  which says
    ┃     the free edit is *not in this
    ┃  round*.
    ┃
    ┃  **Say in your report** what you wrote,
    ┃  and that the three other rows are now
    ┃  in the
    ┃  same voice.
    ┃
    ┃  **One thing to be careful about.** The
    ┃  table is the console's **remit**, so it
    ┃  may name
    ┃  things the console does not do **yet**
    ┃  — that is the point of a remit. But it
    ┃  must not
    ┃  name things **no console will ever
    ┃  do**, and it must not name **another
    ┃  agent's** work
    ┃  as the console's own: *"start, pause"*
    ┃  is the **agent's** job, and the chef's
    ┃  console
    ┃  drives it. If a row reads that way,
    ┃  **say so in your report** — the
    ┃  orchestrator
    ┃  ruled the table's voice, and a row that
    ┃  crosses the line is a question the
    ┃  orchestrator
    ┃  answers, not one you settle.
    ┃
    ┃  ## F-121 (low, ruled in by the
    ┃  orchestrator) — the prompt's keys must
    ┃  not collide
    ┃
    ┃  **The target:** every kind the prompt
    ┃  offers is reachable by its own key.
    ┃
    ┃  **The hole.** `kindKey` takes the
    ┃  kind's **first byte** — so two kinds
    ┃  sharing an
    ┃  initial share the prompt's key, and one
    ┃  of them becomes **unreachable**. Today
    ┃  the
    ┃  kinds are `feature` and `bug`, so `f`
    ┃  and `b` — but the domain's kinds are a
    ┃  `var`,
    ┃  and a future kind whose name starts
    ┃  with the same letter would silently
    ┃  disappear
    ┃  from the prompt.
    ┃
    ┃  The **letter** that stands for a kind
    ┃  is the console's own affordance — the
    ┃  runner
    ┃  said so, and it is right: the domain
    ┃  owns the kind, the console owns the
    ┃  key. So the
    ┃  console's own rule must be that the
    ┃  letters it hands out do not collide.
    ┃
    ┃  **What lands:** the prompt's choices
    ┃  carry **distinct** keys, and the test
    ┃  says so.
    ┃
    ┃  **The evidence:** the mutation must
    ┃  fail the suite — two kinds sharing an
    ┃  initial must
    ┃  not produce two choices with the same
    ┃  key.
    ┃
    ┃  ## The notes the review left, and the
    ┃  orchestrator's ruling
    ┃
    ┃  **One stays for a later round, and here
    ┃  is why:** `BacklogMoveOf` builds the
    ┃  kinds
    ┃  list on **every call**, and the footer
    ┃  asks it on every frame. That is a per-
    ┃  frame
    ┃  allocation — a performance question
    ┃  about the console's render path, and
    ┃  the console's
    ┃  render path has its own shape (it
    ┃  redraws the whole screen). It goes in
    ┃  the ledger for
    ┃  a later round; it is not this slice's.
    ┃
    ┃  ## The scope — what not to touch
    ┃
    ┃  | | |
    ┃  |---|---|
    ┃  | **The moves** | **closed.** Shaping,
    ┃  dismissing, creating, and the domain's
    ┃  answer for the move |
    ┃  | **`internal/core`** | untouched by
    ┃  this cycle, unless F-121's distinct
    ┃  letters need it — and then **say so**
    ┃  rather than reaching for it |
    ┃  | **`internal/cli`** | unchanged |
    ┃  | **The tabs** | seven, the same names
    ┃  |
    ┃  | **`core.DomainVersion`** | stays
    ┃  **7** |
    ┃  | **The goldens** | expected unchanged
    ┃  |
    ┃  | **`work/tui/plan.md`** | the
    ┃  **orchestrator's** file — the
    ┃  orchestrator corrected it. **Do not
    ┃  touch it** |
    ┃  | **`docs/audit.md`** | the ledger is
    ┃  the **reviewer's**. This cycle files no
    ┃  finding and touches no row |
    ┃  | **`work/`** | the orchestrator's,
    ┃  except the state row |
    ┃
    ┃  ## The invariants
    ┃
    ┃  - `internal/core` and `internal/cli`
    ┃  stay **stdlib-only**; `internal/ui`
    ┃  keeps its
    ┃    budget. **No dependency is added.**
    ┃  - **The domain layer has zero UI
    ┃  imports.**
    ┃  - **The tab count is 7** and the tab
    ┃  names do not change.
    ┃  - **One fact, one answer** — a fact the
    ┃  console shows is computed in the
    ┃  domain.
    ┃  - The gate is `make check`. A focused
    ┃  check is `make test PKG=… RUN=…`.
    ┃
    ┃  ## The evidence this cycle must produce
    ┃
    ┃  | Finding | The evidence |
    ┃  |---|---|
    ┃  | **F-120** | `docs/05`'s job list says
    ┃  once what the table is, and **every row
    ┃  is in that voice** — the backlog's row
    ┃  no longer carries a build state. Say
    ┃  what you wrote and how each row now
    ┃  reads |
    ┃  | **F-121** | the prompt's choices
    ┃  carry **distinct** keys, and the
    ┃  mutation — two kinds sharing an initial
    ┃  — turns the suite **red** |
    ┃
    ┃  Plus: `make check` green, `make check-
    ┃  race` green, the goldens unchanged, and
    ┃  the tree
    ┃  clean.
    ┃
    ┃  ## The report — the receipt comes FIRST
    ┃
    ┃  ```bash
    ┃  mkdir -p work/inbox
    ┃  cat > work/inbox/runner-tui-t3b-fix2.md
    ┃  <<EOF
    ┃  Runner: T3b fix 2 done. Commits: <n>.
    ┃  Fixed: F-120, F-121. The job list: <the
    ┃  line you added, and how each row now
    ┃  reads>. The keys: <how the prompt's
    ┃  choices keep distinct keys, and the
    ┃  mutation that now fails it>.
    ┃  DomainVersion: <7>. The goldens:
    ┃  <unchanged / changed, and why>. The
    ┃  gate: <make check, make check-race>.
    ┃  Notes: <anything the reviewer should
    ┃  know>.
    ┃  EOF
    ┃  git add work/inbox && git commit -q -m
    ┃  "inbox: runner reports tui-t3b fix 2"
    ┃  herdr agent prompt "w1H:p1" "Runner:
    ┃  T3b fix 2 done. Commits: <n>. Fixed: F-
    ┃  120, F-121. The job list: <the line you
    ┃  added, and how each row now reads>. The
    ┃  keys: <how the prompt's choices keep
    ┃  distinct keys, and the mutation that
    ┃  now fails it>. DomainVersion: <7>. The
    ┃  goldens: <unchanged / changed, and
    ┃  why>. The gate: <make check, make
    ┃  check-race>. Notes: <anything the
    ┃  reviewer should know>."
    ┃  ```
    ┃
    ┃  **Write the receipt, commit it, and
    ┃  *then* nudge** — the nudge only reaches
    ┃  an idle
    ┃  receiver; the receipt is the half that
    ┃  cannot be lost.
    ┃
    ┃  **One write per file, and commit as
    ┃  soon as the file is coherent.**
    ┃
    ┃  **And if you find a file you did not
    ┃  write has changed: stop and report
    ┃  it**, and do
    ┃  not touch it.
    ┃
    ┃  ## Before you report
    ┃
    ┃  ```bash
    ┃  make check && git status --short
    ┃  ```
    ┃
    ┃  Green, and the tree clean. Report the
    ┃  **real** numbers: the commits, the
    ┃  gate's
    ┃  result, and how the job list reads now.
    ┃
    ┃
    ┃  ---
    ┃
    ┃  ## The notify command — verbatim, use
    ┃  it as it is
    ┃
    ┃  The orchestrator's pane id is **w1H:
    ┃  p1**. Your receipt is the record; the
    ┃  nudge is the
    ┃  bell. Send **both**, in this order, and
    ┃  the nudge **last**:
    ┃
    ┃  ```bash
    ┃  mkdir -p work/inbox
    ┃  cat > work/inbox/runner-tui-t3b-fix2.md
    ┃  <<EOF
    ┃  Runner: T3b fix 2 done. Commits: <n>.
    ┃  Fixed: F-120, F-121. The job list: <the
    ┃  line you added, and how each row now
    ┃  reads>. The keys: <how the prompt's
    ┃  choices keep distinct keys, and the
    ┃  mutation that now fails it>.
    ┃  DomainVersion: <7>. The goldens:
    ┃  <unchanged / changed, and why>. The
    ┃  gate: <make check, make check-race>.
    ┃  Notes: <anything the reviewer should
    ┃  know>.
    ┃  EOF
    ┃  git add work/inbox && git commit -q -m
    ┃  "inbox: runner reports tui-t3b fix 2"
    ┃  herdr agent prompt "w1H:p1" "Runner:
    ┃  T3b fix 2 done. Commits: <n>. Fixed: F-
    ┃  120, F-121. The job list: <the line you
    ┃  added, and how each row now reads>. The
    ┃  keys: <how the prompt's choices keep
    ┃  distinct keys, and the mutation that
    ┃  now fails it>. DomainVersion: <7>. The
    ┃  goldens: <unchanged / changed, and
    ┃  why>. The gate: <make check, make
    ┃  check-race>. Notes: <anything the
    ┃  reviewer should know>."
    ┃  ```
    ┃
    ┃  **If the nudge fails** — it only
    ┃  reaches an idle receiver — the receipt
    ┃  is the half
    ┃  that cannot be lost. Write the receipt
    ┃  anyway, and say so in it.
    ┃
    ┃  **Nothing else waits on you.** After
    ┃  the nudge, end your turn. Do not poll
    ┃  for a
    ┃  reply, do not re-dispatch yourself, do
    ┃  not start another slice. The
    ┃  orchestrator
    ┃  decides what happens next.
    ┃
    ┃  ## Two things this brief cannot know
    ┃
    ┃  - **The orchestrator's earlier ruling
    ┃  was wrong, and the reviewer caught it.
    ┃  ** The
    ┃    orchestrator ruled the console's free
    ┃  edit out of `docs/05`'s job list; the
    ┃  fix made
    ┃    that table speak in two voices - one
    ┃  row in build state, three in remit. The
    ┃  ruling
    ┃    now is that the table is the
    ┃  **remit**, with one line saying so, and
    ┃  that what is
    ┃    built is the keys table below it.
    ┃  - **The ADR's own job table stays**,
    ┃  and it is the same table. The
    ┃  orchestrator's ruling
    ┃    on F-117 says why: a decision of
    ┃  record about the division of labour, in
    ┃  the
    ┃    vocabulary of its day, and the
    ┃  console *does* change backlog items -
    ┃  the kind, the
    ┃    scope, the body - so *"edit"* there
    ┃  is not a promise of a five-field form.
    ┃

       + Thought · 15.6s

       → Explored: 1 read

    ┃
    ┃  ← Edit docs/05-human-interface.md
    ┃
    ┃     6   ## The four jobs
    ┃     7
    ┃     8   From [ADR-0003](decisions/
    ┃         0003-surface-split.md), the
    ┃         console owns:
    ┃     9
    ┃    10 + **This is the console's
    ┃         remit** — what the console is
    ┃         *for*, whether or not it is
    ┃    11 + built yet. What is **built**
    ┃         is the [keys](#interaction)
    ┃         below.
    ┃    12 +
    ┃    13   | Job | What it is |
    ┃    14   |---|---|
    ┃    15   | **View** | Read the plan
    ┃         well — the recipe as prose,
    ┃         the acceptance criteria with
    ┃         their evidence, the audit
    ┃         trail, findings anchored to
    ┃         code |
    ┃    16   | **Decide** | Sign off,
    ┃         triage tasting notes, accept
    ┃         or reject a scope, choose
    ┃         what is next |
    ┃    14 - | **Manage the backlog** |
    ┃         Create, shape, promote,
    ┃         dismiss backlog items — the
    ┃         **free edit** of an item's
    ┃         fields is **not built** |
    ┃    17 + | **Manage the backlog** |
    ┃         Create, edit, shape, promote,
    ┃         dismiss backlog items |
    ┃    18   | **Run** | Start, pause,
    ┃         watch, steer execution |
    ┃    19
    ┃    20   It does **not** author the
    ┃         recipe. An agent does that,
    ┃         in a conversation with the
    ┃         chef,
    ┃    21   through `maitre-cli`.
    ┃

       Now F-121 — the prompt's keys must not
       collide:

    ┃
    ┃  ← Edit internal/ui/backlog.go
    ┃
    ┃    32   // a domain fact is the shape
    ┃         F-114 closed: the console
    ┃         decides for itself and
    ┃    33   // nothing holds it.
    ┃    34   //
    ┃    35   // The letter that stands for
    ┃         a kind is the console's own
    ┃         affordance — the domain
    ┃    36 - // owns the kind, not the key
    ┃         that picks it.
    ┃    36 + // owns the kind, not the key
    ┃         that picks it. And the
    ┃         console's own rule is that
    ┃         the
    ┃    37 + // keys it hands out **do not
    ┃         collide**: two kinds sharing
    ┃         an initial would leave one
    ┃    38 + // of them unreachable.
    ┃    39   func shapeChoices(move core.
    ┃         BacklogMove) []promptChoice {
    ┃    40     out := make([]promptChoice,
    ┃         0, len(move.Kinds))
    ┃    41 +   used := map[string]bool{}
    ┃    42     for _, kind := range move.
    ┃         Kinds {
    ┃    40 -     out = append(out,
    ┃         promptChoice{key: kindKey(
    ┃         kind), label: kind, value:
    ┃         kind})
    ┃    43 +     key := kindKey(kind,
    ┃         used)
    ┃    44 +     used[key] = true
    ┃    45 +     out = append(out,
    ┃         promptChoice{key: key, label:
    ┃         kind, value: kind})
    ┃    46     }
    ┃    47     return out
    ┃    48   }
    ┃    49
    ┃    45 - // kindKey is the letter that
    ┃         picks a kind in the prompt.
    ┃    46 - func kindKey(kind string)
    ┃         string {
    ┃    50 + // kindKey is the key that
    ┃         picks a kind in the prompt.
    ┃    51 + //
    ┃    52 + // The kind's initial — and
    ┃         then the shortest prefix of
    ┃         the kind that no other kind
    ┃    53 + // has taken. Two kinds
    ┃         sharing an initial must not
    ┃         leave one of them
    ┃         unreachable:
    ┃    54 + // the domain's kinds are a
    ┃         `var`, and a kind whose name
    ┃         starts with the same letter
    ┃    55 + // would silently vanish from
    ┃         the prompt.
    ┃    56 + func kindKey(kind string,
    ┃         used map[string]bool) string
    ┃         {
    ┃    57     if kind == "" {
    ┃    58       return ""
    ┃    59     }
    ┃    50 -   return string(kind[0])
    ┃    60 +   runes := []rune(kind)
    ┃    61 +   for n := 1; n <= len(runes)
    ┃         ; n++ {
    ┃    62 +     key := string(runes[:n])
    ┃    63 +     if !used[key] {
    ┃    64 +       return key
    ┃    65 +     }
    ┃    66 +   }
    ┃    67 +   // Every prefix is taken,
    ┃         which means two kinds are the
    ┃         same word — the domain
    ┃    68 +   // does not allow it. Take
    ┃         the first free letter, so no
    ┃         choice is unreachable.
    ┃    69 +   for r := 'a'; r <= 'z'; r++
    ┃         {
    ┃    70 +     if !used[string(r)] {
    ┃    71 +       return string(r)
    ┃    72 +     }
    ┃    73 +   }
    ┃    74 +   return ""
    ┃    75   }
    ┃    76
    ┃    77   // backlogPrimary does the
    ┃         move the domain says the
    ┃         selected item has.
    ┃    78   //
    ┃

       Now the tests:

    ┃
    ┃  ← Edit …/inte…/ui/backlog_moves_test.go
    ┃
    ┃    136       t.Fatalf("the plan holds
    ┃          %d recipes, want %d — a
    ┃          shaping creates no recipe",
    ┃          len(after), len(before))
    ┃    137     }
    ┃    138   }
    ┃    139
    ┃    140 + //
    ┃          TestTheShapePromptGivesEvery
    ┃          KindItsOwnKey covers the
    ┃          prompt's keys.
    ┃    141 + //
    ┃    142 + // The domain owns the kind;
    ┃          the console owns the letter
    ┃          that picks it. And the
    ┃    143 + // console's own rule is
    ┃          that the letters it hands
    ┃          out **do not collide**: two
    ┃          kinds
    ┃    144 + // sharing an initial would
    ┃          leave one of them
    ┃          unreachable in the prompt.
    ┃    145 + //
    ┃    146 + // The domain's kinds are a
    ┃          `var`, so a kind whose name
    ┃          starts with the same letter
    ┃    147 + // is a matter of time — and
    ┃          it would vanish silently.
    ┃    148 + func
    ┃          TestTheShapePromptGivesEvery
    ┃          KindItsOwnKey(t *testing.T)
    ┃          {
    ┃    149 +   // Two kinds sharing an
    ┃          initial.
    ┃    150 +   kinds := []string{
    ┃          "feature", "funnel"}
    ┃    151 +   choices := shapeChoices(
    ┃          core.BacklogMove{Kinds:
    ┃          kinds})
    ┃    152 +   if len(choices) != len(
    ┃          kinds) {
    ┃    153 +     t.Fatalf("the prompt
    ┃          offers %d kinds, want %d",
    ┃          len(choices), len(kinds))
    ┃    154 +   }
    ┃    155 +   seen := map[string]bool{}
    ┃    156 +   for i, ch := range choices
    ┃          {
    ┃    157 +     if ch.value != kinds[i]
    ┃          {
    ┃    158 +       t.Errorf("the choice
    ┃          is %q, want %q", ch.value,
    ┃          kinds[i])
    ┃    159 +     }
    ┃    160 +     if seen[ch.key] {
    ┃    161 +       t.Errorf("two kinds
    ┃          share the key %q — one of
    ┃          them is unreachable", ch.
    ┃          key)
    ┃    162 +     }
    ┃    163 +     seen[ch.key] = true
    ┃    164 +   }
    ┃    165 + }
    ┃    166 +
    ┃    167 + //
    ┃          TestTheShapePromptsKeysAreDi
    ┃          stinct covers the prompt the
    ┃          chef sees.
    ┃    168 + //
    ┃    169 + // The kinds come from the
    ┃          domain, and the keys the
    ┃          console hands out are
    ┃          distinct.
    ┃    170 + func
    ┃          TestTheShapePromptsKeysAreDi
    ┃          stinct(t *testing.T) {
    ┃    171 +   m := newModel(t)
    ┃    172 +   m.tab = TabBacklog
    ┃    173 +   m.selBacklog = 0
    ┃    174 +
    ┃    175 +   if _, err := Drive(m, []
    ┃          string{"space"}); err != nil
    ┃          {
    ┃    176 +     t.Fatal(err)
    ┃    177 +   }
    ┃    178 +   if m.prompt == nil {
    ┃    179 +     t.Fatal("the prompt
    ┃          should be open")
    ┃    180 +   }
    ┃    181 +   seen := map[string]bool{}
    ┃    182 +   for _, ch := range m.
    ┃          prompt.choices {
    ┃    183 +     if ch.key == "" {
    ┃    184 +       t.Errorf("%q has no
    ┃          key — it is unreachable",
    ┃          ch.value)
    ┃    185 +     }
    ┃    186 +     if seen[ch.key] {
    ┃    187 +       t.Errorf("two kinds
    ┃          share the key %q — one of
    ┃          them is unreachable", ch.
    ┃          key)
    ┃    188 +     }
    ┃    189 +     seen[ch.key] = true
    ┃    190 +   }
    ┃    191 + }
    ┃    192 +
    ┃    193   //
    ┃          TestTheShapePromptOffersTheD
    ┃          omainsKinds covers the
    ┃          prompt's choices.
    ┃    194   //
    ┃    195   // The kinds are the
    ┃          **domain's**: `core.
    ┃          BacklogMoveOf` carries the
    ┃          kinds a shaping
    ┃    196   // may pick, and the prompt
    ┃          offers exactly those. A
    ┃          console-side copy of a
    ┃          domain
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre
    ┃  && go build ./... && go vet ./..…
    ┃
    ┃  ==> test
    ┃  GO=go scripts/go-test.sh " -run
    ┃  Shape|Backlog|Move" "./internal/ui"
    ┃  ok    github.com/lthiagol/maitre/
    ┃  internal/ui  12.322s
    ┃

       Now the mutation that must fail the
       suite:

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre
    ┃  && cp internal/ui/backlog.go /tm…
    ┃
    ┃  (2 earlier lines)
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/
    ┃  internal/ui  9.987s
    ┃  FAIL
    ┃   M docs/05-human-interface.md
    ┃   M internal/ui/backlog.go
    ┃   M internal/ui/backlog_moves_test.go
    ┃

       + Thought · 1.1s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre
    ┃  && make fmt-check && make vet &&…
    ┃
    ┃  ==> fmt-check
    ┃
    ┃  ==> vet
    ┃  go vet ./...
    ┃  c0f1624 ui: the prompt's keys do not
    ┃  collide
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre
    ┃  && git add docs/05-human-interfa…
    ┃
    ┃  89eaecc docs: the console's job list
    ┃  speaks in one voice
    ┃

       + Thought · 261ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre
    ┃  && ls work/briefs/ | grep t3b &&…
    ┃
    ┃  (4 earlier lines)
    ┃  *The console never authors a recipe** |
    ┃  ADR-0003. A *backlog item* is not a
    ┃  recipe, so the console may create one (
    ┃  T3b) |
    ┃  156:| Q6 | **T3b** — the console
    ┃  creates backlog items. The **free
    ┃  edit** of every field is **not built**:
    ┃  a later round, if it is wanted at all |
    ┃

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃    15   | **T1** | the identity | the
    ┃         window title carries the
    ┃         project; the header's name
    ┃         stops being a footnote; `--
    ┃         project`; the Settings tab
    ┃         shows the plan's path | ✅
    ┃         closed · 11 commits ·
    ┃         94a53ff, 8fc7438, f0b2aa5,
    ┃         82289c6, 0bbcb35, 5a79638,
    ┃         67b455d, b4a3052, 55bdc81,
    ┃         8a6866b, cf8c05e | [tui-t1.
    ┃         md](../briefs/tui-t1.md) · [
    ┃         fix 1](../briefs/fix-tui-t1-
    ┃         1.md) · [fix 2](../briefs/
    ┃         fix-tui-t1-2.md) |
    ┃    16   | **T2a** | the mouse, part 1
    ┃         | the mouse is on; the render
    ┃         builds the hit map; the wheel
    ┃         scrolls | ✅ closed · 12
    ┃         commits · a72ea98, 115b650,
    ┃         284fb65, 0672f9f, 9f971cb,
    ┃         7035ec7, e4b7061, 619b342,
    ┃         47be828, f1cc1ca, a332158,
    ┃         d8d7b95 | [tui-t2a.md](../
    ┃         briefs/tui-t2a.md) · [fix 1](
    ┃         ../briefs/fix-tui-t2a-1.md) ·
    ┃         [fix 2](../briefs/fix-tui-
    ┃         t2a-2.md) |
    ┃    17   | **T2b** | the mouse, part 2
    ┃         | the clicks: a tab, a row, a
    ┃         detail; the mouse drivable
    ┃         from a script | ✅ closed · 9
    ┃         commits · 4dfd5a4, 187d706,
    ┃         f64e6a2, ac8df77, 588cf5f,
    ┃         9c716ce, 0144216, 6bb7545,
    ┃         f8969ef | [tui-t2b.md](../
    ┃         briefs/tui-t2b.md) · [fix 1](
    ┃         ../briefs/fix-tui-t2b-1.md) |
    ┃    18   | **T3a** | the backlog is
    ┃         readable | the detail pane;
    ┃         `space` stops being a trap on
    ┃         an idea; D-3, D-5 | ✅ closed
    ┃         · 9 commits · 4dc8079,
    ┃         1b550c6, 1804ccc, b873821,
    ┃         5edae71, aa4ba5f, da36146,
    ┃         99a7576, 63964ce | [tui-t3a.
    ┃         md](../briefs/tui-t3a.md) · [
    ┃         fix 1](../briefs/fix-tui-t3a-
    ┃         1.md) · [fix 2](../briefs/
    ┃         fix-tui-t3a-2.md) |
    ┃    19 - | **T3b** | the backlog moves
    ┃         | the domain answers the
    ┃         item's move; the console
    ┃         shapes · promotes · dismisses
    ┃         · creates | 🔄 **fix cycle
    ┃         2** — 6 commits · 1297cdd,
    ┃         35688a6, 4799b3b, 7de9a16,
    ┃         4e14f19, 58069e7 | [tui-t3b.
    ┃         md](../briefs/tui-t3b.md) · [
    ┃         fix 1](../briefs/fix-tui-t3b-
    ┃         1.md) · [fix 2](../briefs/
    ┃         fix-tui-t3b-2.md) |
    ┃    19 + | **T3b** | the backlog moves
    ┃         | the domain answers the
    ┃         item's move; the console
    ┃         shapes · promotes · dismisses
    ┃         · creates | ✅ closed · 8
    ┃         commits · 1297cdd, 35688a6,
    ┃         4799b3b, 7de9a16, 4e14f19,
    ┃         58069e7, c0f1624, 89eaecc | [
    ┃         tui-t3b.md](../briefs/tui-
    ┃         t3b.md) · [fix 1](../briefs/
    ┃         fix-tui-t3b-1.md) · [fix 2](.
    ┃         ./briefs/fix-tui-t3b-2.md) |
    ┃    20   | **T4** | the Home tab | the
    ┃         landing tab; the art; the
    ┃         Sous-Chef dissolves; D-2 | ⏳
    ┃         planned | — |
    ┃    21   | **T5** | the Ready tab |
    ┃         `core.Adoptable`; the Path
    ┃         becomes Ready; D-1, D-4 | ⏳
    ┃         planned | — |
    ┃    22   | **T6** | the Board, as a
    ┃         visual place | the tree drawn
    ┃         left to right; clicking a
    ┃         node | ⏳ planned | — |
    ┃    23   | **T7** | the filter | `/`
    ┃         filters the list | ⏳ planned
    ┃         | — |
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre
    ┃  && git add work/tui/README.md &&…
    ┃
    ┃  1932a0a work: T3b's row names the fix
    ┃  cycle 2's two commits
    ┃  89eaecc docs: the console's job list
    ┃  speaks in one voice
    ┃  c0f1624 ui: the prompt's keys do not
    ┃  collide
    ┃  ee35d2c work: T3b's fix cycle 2 - the
    ┃  job list's voice, and the prompt's keys
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre
    ┃  && make check 2>&1 | tail -22
    ┃
    ┃  (16 earlier lines)
    ┃    internal/core: 0 external modules (
    ┃  ceiling 0)
    ┃    internal/cli: 0 external modules (
    ┃  ceiling 0)
    ┃    internal/ui: 28 external modules (
    ┃  ceiling 30)
    ┃  /Library/Developer/CommandLineTools/
    ┃  usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre
    ┃  && make check-race 2>&1 | tail -…
    ┃
    ┃  (14 earlier lines)
    ┃    internal/ui: 28 external modules (
    ┃  ceiling 30)
    ┃  /Library/Developer/CommandLineTools/
    ┃  usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃  --- tree clean ---
    ┃  --- goldens unchanged ---
    ┃

       + Thought · 374ms

       → Explored: 1 read

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre
    ┃  && make build >/dev/null 2>&1 &&…
    ┃
    ┃  (6 earlier lines)
    ┃    ],
    ┃    "kind": "feature",
    ┃    "tags": [
    ┃      "board"
    ┃    ],
    ┃    "status": "open",
    ┃

       Everything still works end-to-end. The
       receipt first:

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre
    ┃  && mkdir -p work/inbox && cat > …
    ┃
    ┃  b15cb94 inbox: runner reports tui-t3b
    ┃  fix 2
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre
    ┃  && herdr agent prompt "w1H:p1" "…
    ┃
    ┃  (1 earlier line)
    ┃  evision":15,"screen_detection_skipped":
    ┃  true,"state_change_seq":354,"tab_id":
    ┃  "w1H:t1","terminal_id":
    ┃  "term_65f0cd2aa89d26","terminal_title":
    ┃  "OC | Project work proposal",
    ┃  "terminal_title_stripped":"OC | Project
    ┃  work proposal","workspace_id":"w1H"},
    ┃  "type":"agent_prompted"}}
    ┃

       T3b's fix cycle 2 is done. 2 commits +
       the state row; the receipt is committed
       and the nudge landed.

       F-120 — the job list speaks in one
       voice

       One line above the table says what the
       table is:

       │ This is the console's remit — what
       │ the console is for, whether or not it
       │ is built yet. What is built is the
       │ keys below.

       and the backlog's row returns to the
       remit's wording — Create, edit, shape,
       promote, dismiss backlog items. The
       clause "the free edit of an item's
       fields is not built" is out of docs/05
       (a build state; the plan says it is not
       in this round).

       ┌─────────┬───────────────────────────┐
       │ Row     │ How it now reads          │
       ├─────────┼───────────────────────────┤
       │ View    │ read the plan well …      │
       │         │ findings anchored to code │
       │         │ — the console does not    │
       │         │ read the anchor yet, but  │
       │         │ reading a finding is the  │
       │         │ chef's reading, not       │
       │         │ another agent's work: in  │
       │         │ voice, build gap          │
       ├─────────┼───────────────────────────┤
       │ Decide  │ sign off, triage tasting  │
       │         │ notes, accept or reject a │
       │         │ scope, choose what is     │
       │         │ next — the scope and the  │
       │         │ choice are the chef's: in │
       │         │ voice, build gap          │
       ├─────────┼───────────────────────────┤
       │ Manage  │ the remit wording         │
       │ the     │                           │
       │ backlog │                           │
       ├─────────┼───────────────────────────┤
       │ Run     │ start, pause, watch,      │
       │         │ steer execution — the one │
       │         │ row that crosses the line │
       │         │ by the brief's own test:  │
       │         │ start and pause are the   │
       │         │ agent's work, and the     │
       │         │ chef's console drives it. │
       │         │ I did not settle it — the │
       │         │ row is untouched and the  │
       │         │ question is the           │
       │         │ orchestrator's            │
       └─────────┴───────────────────────────┘

       F-121 — the prompt's keys do not
       collide

       The console's own rule: the kind's
       initial, then the shortest prefix of
       the kind that no other kind has taken (
       feature and funnel → f and fu); the
       degenerate case falls back to the first
       free letter, so no choice is ever
       unreachable. The domain owns the kind,
       the console owns the key — and the
       collision rule belongs where the key is
       made. internal/core was not touched.

       The mutation: kindKey back to "the
       first byte and nothing else" → the
       suite went red —
       TestTheShapePromptGivesEveryKindItsOwnK
       ey ("two kinds share the key "f" — one
       of them is unreachable"). Reverted. The
       domain's own kinds don't collide today
       (f, b), so the guard is the unit test
       over two kinds sharing an initial.

       make check green · make check-race
       green · make goldens leaves git diff
       empty · tree clean · DomainVersion 7 ·
       no dependency added. The console still
       works end to end: the prompt offers f)
       feature · b) bug, and shaping BL01
       gives it kind feature, scope evidence,
       id BL01.

       Build · MiniMax-M3.1-Flash · 10m 23s ·
       148.4 tok/s

    ┃
    ┃
    ┃
    ┃  Build · MiniMax-M3.1-Fl… · medium
    ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
    ~/Code/Pessoal/maitre:master  203.8K (20%)
```
