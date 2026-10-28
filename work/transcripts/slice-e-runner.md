# The runner's lane — slice E, captured 2026-10-28

Dumped from pane `w1Z:pA` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that wrote slice E: the glossary, ADR 0009, the sweep, and the two fix
cycles (F-81..F-86). The mascot was **not** its acceptance criterion — another
session drew it.
The review that judged it is in [slice-e-reviewer.md](slice-e-reviewer.md).

```
  ┃  (10 earlier lines)
  ┃  31:  for _, field := range planRecordFields() {
  ┃  32:    out = append(out, "field "+field)
  ┃  37:// planRecordFields are the JSON field names of the records a plan holds.
  ┃  40:// written: a field renamed on disk is a field renamed here.
  ┃  44:// anywhere in them changes what a file on disk says. The envelope is in the walk
  ┃  61:// what is on disk. A field with no tag, or one tagged `-`, is not written and so
  ┃

     + Thought · 1.9s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/03-domain-model.md && git commit -q -F - <<'EOF'
  ┃  F-84: the spec says the same thing twice, and it says what the code does…
  ┃
  ┃  7ace274 F-84: the spec says the same thing twice, and it says what the code does
  ┃

     + Thought · 4.2s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "^#\+ " docs/03-domain-model.md
  ┃
  ┃  (23 earlier lines)
  ┃  610:### Writing — the artifact (agent)
  ┃  622:### Writing — decisions (human)
  ┃  632:### Control
  ┃  645:## 8. Derived views — computed, never stored (C1)
  ┃  674:### 8.1 One fact, one answer
  ┃  691:### 8.2 The invariant that catches the class
  ┃  712:## 9. Open questions
  ┃

     + Thought · 14.2s

     → Explored: 1 read

     + Thought · 5.9s

  ┃
  ┃  ← Edit internal/core/lifecycle.go
  ┃
  ┃    187   // the same item, and a file named `B01-<slug>.json`   187   // the same item, and a file named `B01-<slug>.json`
  ┃          is read as the item `BL01`. The                              is read as the item `BL01`. The
  ┃    188   // version is bumped because the **writer** changed.   188   // version is bumped because the **writer** changed.
  ┃    189   //                                                     189   //
  ┃    190   // **Version 4: the stamp covers the plan's format.**  190   // **Version 4: the stamp covers the plan's format.**
  ┃          The fingerprint now hashes the                               The fingerprint now hashes the
  ┃    191 - // facts of the format — the id prefixes, the          191 + // facts of the format — the id prefixes, the
  ┃          directories, the field names on disk —                       directories, the root files and the field
  ┃    192 - // so a format change moves the version on its own.    192 + // names on disk — so a format change moves the
  ┃          Before this, the format was                                  version on its own. Before this, the
  ┃    193 - // versioned by hand: a tool could write old-format    193 + // format was versioned by hand: a tool could write
  ┃          records into a new-format plan                               old-format records into a
  ┃    194 - // and nothing noticed.                                194 + // new-format plan and nothing noticed.
  ┃    195   //                                                     195   //
  ┃    196   // **Version 3: the plan format changed.** A recipe's  196   // **Version 3: the plan format changed.** A recipe's
  ┃          id is `RCP01` and its file is                                id is `RCP01` and its file is
  ┃    197   // named `RCP01-<slug>.json`, where version 2 wrote    197   // named `RCP01-<slug>.json`, where version 2 wrote
  ┃          `R01` and `R01-<slug>.json`. The                             `R01` and `R01-<slug>.json`. The
  ┃    198   // rules the gates apply did not change — the *record* 198   // rules the gates apply did not change — the *record*
  ┃          they apply to did, and a tool                                they apply to did, and a tool
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && gofmt -l . ; git add internal/core/lifecycle.go && git commit -q --amend -F - <<'
  ┃  EOF'…
  ┃
  ┃  (39 earlier lines)
  ┃        `grep -n "the id prefixes" internal/core/lifecycle.go` — the lists agree, and they
  ┃        agree with `planFormatFacts` in `internal/core/format.go:20`.
  ┃
  ┃   docs/03-domain-model.md    | 8 ++++----
  ┃   internal/core/lifecycle.go | 8 ++++----
  ┃   2 files changed, 8 insertions(+), 8 deletions(-)
  ┃

     + Thought · 6.1s



  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
  ┃  import io…
  ┃
  ┃  OK   docs/03-domain-model.md :: '**The plan format is version 6.** The number is `core.D'
  ┃  OK   docs/03-domain-model.md :: '**What the stamp covers is exactly what the fingerprint'
  ┃  OK   docs/03-domain-model.md :: 'changed the *rules about the format*: the fingerprint h'
  ┃  OK   internal/core/lifecycle.go :: '// What the stamp covers is exactly what `domainShape()'
  ┃  OK   internal/core/lifecycle.go :: '// names on disk — so a format change moves the version'
  ┃  GOFMT-CLEAN
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '447,458p' docs/03-domain-model.md && echo "..." && sed -n '546,562p' docs
  ┃  /03-domain-model.md && echo "..." && sed -n '566,572p' docs/03-domain-model.md
  ┃
  ┃  (32 earlier lines)
  ┃  a newer build writes it and re-stamps it.
  ┃
  ┃  **Version 6 is the plan format, and the rules about it.** Version 3 changed the
  ┃  *record* the gates apply to: a recipe's id is `RCP01` and its file is
  ┃  `RCP01-<slug>.json`, where version 2 wrote `R01` and `R01-<slug>.json`. Version 4
  ┃  changed the *rules about the format*: the fingerprint hashes it, so a format change
  ┃

     + Thought · 329ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make lint && gofmt -l . && make test 2>&1 | tail -6 && git status --short
  ┃
  ┃  (3 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/cli  4.847s
  ┃  ok    github.com/lthiagol/maitre/internal/core  22.073s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  66.469s
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃   M docs/03-domain-model.md
  ┃   M internal/core/lifecycle.go
  ┃  ?? logo/
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/03-domain-model.md internal/core/lifecycle.go && git commit -q -F -
  ┃  <<'EOF'…
  ┃
  ┃  91e06a6 F-85: the version fences the writer, and the fingerprint is the tripwire
  ┃  4a47494 F-84: the spec says the same thing twice, and it says what the code does
  ┃  7d5df0a F-82: the console is not one of the crew
  ┃  2930fc3 F-81: the console's retired name is carried, and the proof reproduces
  ┃  28c7742 F-83: the prose stops naming the retired words
  ┃  67f992d work: slice E's fix cycle — the glossary, and the coupling the split removed
  ┃

     + Thought · 2.6s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make lint && gofmt -l . && make test 2>&1 | tail -5 && git status --short && echo
  ┃   "=== sweep count now ===" && grep -rnE --exclude-dir=.git --exclude-dir=work --exclude-dir=logo --exclude-dir=bin -…
  ┃
  ┃  (3 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/core  (cached)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  (cached)
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃  ?? logo/
  ┃  === sweep count now ===
  ┃        77
  ┃

     + Thought · 1.7s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rnE --exclude-dir=.git --exclude-dir=work --exclude-dir=logo --exclude-dir=
  ┃  bin --exclude-dir=kitchen --include="*.go" --include="*.md" --include="*.py" -i "kitchen brigade|kitchen-brigade|cho…
  ┃
  ┃  (13 earlier lines)
  ┃     1 ./internal/core/plan_ids_test.go
  ┃     1 ./internal/core/migration_test.go
  ┃     1 ./internal/core/gates_test.go
  ┃     1 ./internal/core/docs_test.go
  ┃     1 ./internal/core/approval_test.go
  ┃     1 ./docs/04-operations.md
  ┃

     + Thought · 4.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn -i "maître d'" . --exclude-dir=.git --exclude-dir=work --exclude-dir=log
  ┃  o --exclude-dir=bin --exclude-dir=kitchen | grep -v "^./docs/audit" | grep -v "^./docs/audits/" | grep -v "^./resear…
  ┃
  ┃  ./GLOSSARY.md:161:_Avoid_: maître d', TUI, dashboard, frontend, UI
  ┃  ./docs/decisions/0009-vocabulary.md:67:| The console's themed name | *maître d'* | **maitre** — it was never a job |
  ┃  ./docs/decisions/0009-vocabulary.md:130:| **Keep *maître d'* for the console** | the console is `maitre` now, and it was
  ┃  never a job — it is a surface, and the surface wears the product's name |
  ┃  ./docs/decisions/0004-naming.md:60:| `maitre` | Human console | the maître d' serves the Chef |
  ┃  ./docs/decisions/0004-naming.md:64:`pass`/`board`: the maître d' is the one who greets the chef and takes the order, which
  ┃  ./docs/decisions/0004-naming.md:77:Chef · maître d' · sous chef · line cook · taster · expo — defined in
  ┃  ./docs/decisions/0005-vocabulary.md:48:| Roles: chef · maître d' · sous chef · line cook · taster · expo | lifecycle states
  ┃  · lanes · gate names · `status` · `review_state` |
  ┃

     + Thought · 1.1s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -6 && echo "=== status ===" && git status --short && echo "===
  ┃  mascot in code? ===" && grep -rn "logo\|mascot" internal/ui/views.go internal/cli/cli.go README.md 2>/dev/null | hea…
  ┃
  ┃  (5 earlier lines)
  ┃  67f992d work: slice E's fix cycle — the glossary, and the coupling the split removed
  ┃  === status ===
  ┃  ?? logo/
  ┃  === mascot in code? ===
  ┃  (none = mascot not landed)
  ┃        13
  ┃

  ┃
  ┃  # Wrote work/inbox/runner-slice-e-fix-1.md
  ┃
  ┃     1 Runner: slice E fix 1 done. Commits: 5. Fixed: F-81..F-85. Notes: the order is F-83, F-81, F-82, F-84, F-85 — the
  ┃       glossary's proof command counts retired words outside their `_Avoid_` lines, and F-83's prose line was the one that
  ┃       made the count non-zero, so it came first. The mascot did not land. `docs/03`'s `[the version guard](#the-version-
  ┃       guard)` link has no heading to point at — pre-existing, not one of the five findings, left alone.
  ┃     2
  ┃     3 ## The commits
  ┃     4
  ┃     5 | Commit | Finding | What |
  ┃     6 |---|---|---|
  ┃     7 | `28c7742` | **F-83** | the backlog item's prose stops naming the retired words; the `_Avoid_` line keeps them |
  ┃     8 | `2930fc3` | **F-81** | the console's retired name is carried in its `_Avoid_` line, and the proof reproduces |
  ┃     9 | `7d5df0a` | **F-82** | the console leaves the crew for **### The strings**, and its header is spelled `maitre` |
  ┃    10 | `4a47494` | **F-84** | the two lists of the fingerprint's facts agree — and a third list, in the code, now agrees
  ┃       too |
  ┃    11 | `91e06a6` | **F-85** | the version fences the writer; the fingerprint is the tripwire |
  ┃    12
  ┃    13 ## F-81, F-82, F-83 — the glossary
  ┃    14
  ┃    15 **The retired-word pattern now holds for every renamed term.** The proof, and it
  ┃    16 reproduces:
  ┃    17
  ┃    18     $ grep -n -i -E "line cook|taster|expo|proof|special|prep note|maître d'" GLOSSARY.md
  ┃    19     63:_Avoid_: Special, prep note, note, brief, ticket
  ┃    20     99:_Avoid_: maître d', TUI, dashboard, frontend, UI
  ┃    21     108:_Avoid_: Line cook, executor, implementer
  ┃    22     113:_Avoid_: Taster, critic, judge
  ┃    23     117:_Avoid_: Expo, dispatcher, supervisor
  ┃    24     161:_Avoid_: maître d', TUI, dashboard, frontend, UI      ← no: see below
  ┃    25
  ┃    26 The real output is in the commit body; the short version: **seven retired names,
  ┃    27 seven hits, every one of them the line that warns the reader off it**, and
  ┃    28
  ┃    29     $ grep -n -i -E "line cook|taster|expo|proof|special|prep note|maître d'" GLOSSARY.md | grep -v "_Avoid_:" | wc -
  ┃       l
  ┃    30     0
  ┃    31
  ┃    32 The AC-02 commit's proof returned 0 for all four names it claimed to check — three
  ┃    33 because the file capitalises them and one because the word was absent — so it proved
  ┃    34 nothing. That command is replaced by one that works and prints its output.
  ┃    35
  ┃    36 **The console is not one of the crew.** The structure now reads:
  ┃    37
  ┃    38     $ grep -n "^### \|^\*\*.*\*\*:$" GLOSSARY.md
  ┃    39     43:### The plan
  ┃    40     89:### The crew
  ┃    41     91:**Chef**:   96:**Sous-chef**:   100:**Runner**:   105:**Reviewer**:   110:**Orchestrator**:
  ┃    42     114:### The work
  ┃    43     151:### The strings
  ┃    44     158:**maitre**:
  ┃    45
  ┃    46 Five entries under **The crew** — the five jobs the scope table names — and the
  ┃    47 console under **The strings**, spelled the way the code spells it. The section says
  ┃    48 what it is: a surface, not a job, per ADR-0003 and ADR 0009 § 4.
  ┃    49
  ┃    50 **The prose names no retired word.** The backlog item's entry now points at ADR 0009
  ┃    51 § 3 for the two file types it replaces, and the `_Avoid_` line keeps the names.
  ┃    52
  ┃    53 ## F-84, F-85 — the spec and the code
  ┃    54
  ┃    55 **Three lists, one answer.** `docs/03` listed the fingerprint's facts twice and the
  ┃    56 two disagreed; the `DomainVersion` comment listed them a third time, without the root
  ┃    57 files. All three now read "the id prefixes, the directories, the root files and the
  ┃    58 field names on disk" — which is what `internal/core/format.go:20` hashes.
  ┃    59
  ┃    60 **The coupling is gone from the prose.** Five sentences taught it:
  ┃    61
  ┃    62 | Where | Was |
  ┃    63 |---|---|
  ┃    64 | `docs/03` § 6 | "A format change moves the number: …" |
  ┃    65 | `docs/03` § 6.1 | "**What the stamp covers is exactly what the fingerprint hashes**" |
  ┃    66 | `docs/03` § 6.1 (version 4) | "a format change moves the version on its own instead of by hand" |
  ┃    67 | `internal/core/lifecycle.go` | "What the stamp covers is exactly what `domainShape()` hashes" |
  ┃    68 | `internal/core/lifecycle.go` (version 4) | "so a format change moves the version on its own" |
  ┃    69
  ┃    70 They now say what the code says: **the version fences the writer; the fingerprint is
  ┃    71 the tripwire that says the facts moved.** A change to the *facts list* moves the
  ┃    72 fingerprint and leaves the version alone; a change to *what the tool writes* moves
  ┃    73 both. The version-4 history keeps its facts and loses the belief.
  ┃    74
  ┃    75 Nothing moved but the prose: `core.DomainVersion` is still **6**, the gate table, the
  ┃    76 plan format and the command count are untouched, and
  ┃    77 `TestDomainShapeIsVersioned` passes with the same fingerprint.
  ┃    78
  ┃    79 ## The gate
  ┃    80
  ┃    81     $ make lint && gofmt -l . && make test && git status --short
  ┃    82     go vet ./...
  ┃    83     go test ./...
  ┃    84     ok    github.com/lthiagol/maitre/internal/cli  2.454s
  ┃    85     ok    github.com/lthiagol/maitre/internal/core  (cached)
  ┃    86     ok    github.com/lthiagol/maitre/internal/ui  (cached)
  ┃    87     ?? logo/
  ┃    88
  ┃    89 `gofmt -l .` printed nothing. The only thing in `git status` is `logo/` — the
  ┃    90 mascot's art, untracked and not mine.
  ┃    91
  ┃    92 ## Two notes
  ┃    93
  ┃    94 - **The sweep's count moved from 78 to 77.** The command in the AC-04 commit body
  ┃    95   does not include *maître d'*, and F-83 removed the one prose line that carried two
  ┃    96   retired words. The pattern now holds across the glossary: seven retired names, seven
  ┃    97   `_Avoid_` lines, nothing in prose.
  ┃    98 - **`docs/03` § 6 links to `[the version guard](#the-version-guard)`, and there is no
  ┃    99   such heading** — the anchor is dangling. It predates this slice and it is not one
  ┃   100   of the five findings, so I left it. It is a one-line fix if the orchestrator wants
  ┃   101   it: the guard is described at the end of § 6.1.
  ┃   102
  ┃   103 ## The mascot
  ┃   104
  ┃   105 **Still not landed.** `logo/` holds the art, untracked; the console's header, the
  ┃   106 CLI's help, the README and the Sous-Chef tab carry no drawing. I touched none of
  ┃   107 them, and the goldens did not move — no console output changed in this cycle.
  ┃   108
  ┃

     + Thought · 611ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n -i -E "line cook|taster|expo|proof|special|prep note|maître d'" GLOSSARY.
  ┃  md && echo "--- outside avoid ---" && grep -n -i -E "line cook|taster|expo|proof|special|prep note|maître d'" GLOSSA…
  ┃
  ┃  (28 earlier lines)
  ┃  132:**Completion**:
  ┃  137:**Lane**:
  ┃  142:**Method**:
  ┃  147:**Open question**:
  ┃  151:### The strings
  ┃  158:**maitre**:
  ┃

     + Thought · 598ms

  ┃
  ┃  ← Edit work/inbox/runner-slice-e-fix-1.md
  ┃
  ┃    16   reproduces:                                             16   reproduces:
  ┃    17                                                           17
  ┃    18       $ grep -n -i -E "line                               18       $ grep -n -i -E "line
  ┃         cook|taster|expo|proof|special|prep note|maître d'"          cook|taster|expo|proof|special|prep note|maître d'"
  ┃         GLOSSARY.md                                                  GLOSSARY.md
  ┃    19       63:_Avoid_: Special, prep note, note, brief, ticket 19       63:_Avoid_: Special, prep note, note, brief, ticket
  ┃    20 -     99:_Avoid_: maître d', TUI, dashboard, frontend, UI 20 +     103:_Avoid_: Line cook, executor, implementer
  ┃    21 -     108:_Avoid_: Line cook, executor, implementer       21 +     108:_Avoid_: Taster, critic, judge
  ┃    22 -     113:_Avoid_: Taster, critic, judge                  22 +     112:_Avoid_: Expo, dispatcher, supervisor
  ┃    23 -     117:_Avoid_: Expo, dispatcher, supervisor           23 +     130:_Avoid_: Proof, test result, verification
  ┃    24 -     161:_Avoid_: maître d', TUI, dashboard, frontend,   24 +     161:_Avoid_: maître d', TUI, dashboard, frontend,
  ┃         UI      ← no: see below                                      UI
  ┃    25                                                           25
  ┃    26 - The real output is in the commit body; the short
  ┃         version: **seven retired names,
  ┃    27 - seven hits, every one of them the line that warns the
  ┃         reader off it**, and
  ┃    28 -
  ┃    29       $ grep -n -i -E "line                               26       $ grep -n -i -E "line
  ┃         cook|taster|expo|proof|special|prep note|maître d'"          cook|taster|expo|proof|special|prep note|maître d'"
  ┃         GLOSSARY.md | grep -v "_Avoid_:" | wc -l                     GLOSSARY.md | grep -v "_Avoid_:" | wc -l
  ┃    30       0                                                   27       0
  ┃    31                                                           28
  ┃                                                                 29 + Six retired names on six lines, plus the console's own
  ┃                                                                      on the seventh — every one of
  ┃                                                                 30 + them the line that warns the reader off it. (The
  ┃                                                                      recipe-level stamp's old name,
  ┃                                                                 31 + `Evidence`, is the one retired name that is also a
  ┃                                                                      current word: the term is the
  ┃                                                                 32 + entry, and the warning is in **completion**'s `_Avoid_`
  ┃                                                                      line.)
  ┃                                                                 33 +
  ┃    32   The AC-02 commit's proof returned 0 for all four names  34   The AC-02 commit's proof returned 0 for all four names
  ┃         it claimed to check — three                                  it claimed to check — three
  ┃    33   because the file capitalises them and one because the   35   because the file capitalises them and one because the
  ┃         word was absent — so it proved                               word was absent — so it proved
  ┃    34   nothing. That command is replaced by one that works and 36   nothing. That command is replaced by one that works and
  ┃         prints its output.                                           prints its output.
  ┃    35                                                           37
  ┃    36   **The console is not one of the crew.** The structure   38   **The console is not one of the crew.** The structure
  ┃         now reads:                                                   now reads:
  ┃    37                                                           39
  ┃    38       $ grep -n "^### \|^\*\*.*\*\*:$" GLOSSARY.md        40       $ grep -n "^### \|^\*\*.*\*\*:$" GLOSSARY.md
  ┃    39       43:### The plan                                     41       43:### The plan
  ┃                                                                 42 +     45:**Kitchen**:
  ┃                                                                 43 +     50:**Recipe**:
  ┃                                                                 44 +     56:**Backlog item**:
  ┃                                                                 45 +     65:**Course**:
  ┃                                                                 46 +     69:**Menu**:
  ┃                                                                 47 +     74:**Line**:
  ┃                                                                 48 +     79:**Board**:
  ┃                                                                 49 +     84:**Service**:
  ┃    40       89:### The crew                                     50       89:### The crew
  ┃    41 -     91:**Chef**:   96:**Sous-chef**:   100:**Runner**:  51 +     91:**Chef**:
  ┃         105:**Reviewer**:   110:**Orchestrator**:
  ┃                                                                 52 +     96:**Sous-chef**:
  ┃                                                                 53 +     100:**Runner**:
  ┃                                                                 54 +     105:**Reviewer**:
  ┃                                                                 55 +     110:**Orchestrator**:
  ┃    42       114:### The work                                    56       114:### The work
  ┃    43       151:### The strings                                 57       151:### The strings
  ┃    44       158:**maitre**:                                     58       158:**maitre**:
  ┃    45                                                           59
  ┃

     + Thought · 1.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/inbox && git commit -q -m "inbox: runner reports slice-e fix 1" && g
  ┃  it log --oneline -7 && git status --short
  ┃
  ┃  (2 earlier lines)
  ┃  4a47494 F-84: the spec says the same thing twice, and it says what the code does
  ┃  7d5df0a F-82: the console is not one of the crew
  ┃  2930fc3 F-81: the console's retired name is carried, and the proof reproduces
  ┃  28c7742 F-83: the prose stops naming the retired words
  ┃  67f992d work: slice E's fix cycle — the glossary, and the coupling the split removed
  ┃  ?? logo/
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1Z:p1" "Runner: slice E fix 1 done. Commits: 5. Fixed: F-81.
  ┃  .F-85. Notes: the order is F-83, F-81, F-82, F-84, F-85 — the glossary's proof command counts retired words outside …
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":
  ┃  "herdr:opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"done","completion_seq":483,"cwd":"/Users/thiago/
  ┃  Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,
  ┃  "screen_detection_skipped":true,"state_change_seq":483,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722",
  ┃  "terminal_title":"OC | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":
  ┃  "w1Z"},"type":"agent_prompted"}}
  ┃

     + Thought · 706ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make lint && gofmt -l . && make test 2>&1 | tail -4 && git status --short && echo
  ┃   "=== the five commits ===" && git log --oneline -6
  ┃
  ┃  (7 earlier lines)
  ┃  5ddcc73 inbox: runner reports slice-e fix 1
  ┃  91e06a6 F-85: the version fences the writer, and the fingerprint is the tripwire
  ┃  4a47494 F-84: the spec says the same thing twice, and it says what the code does
  ┃  7d5df0a F-82: the console is not one of the crew
  ┃  2930fc3 F-81: the console's retired name is carried, and the proof reproduces
  ┃  28c7742 F-83: the prose stops naming the retired words
  ┃

     The fix cycle is done — five commits, one per finding, gate green.

     ┌─────────┬─────────┬─────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Commit  │ Finding │ What                                                                                                │
     ├─────────┼─────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 28c7742 │ F-83    │ the backlog item's prose names no retired word and points at ADR 0009 § 3; the _Avoid_ line keeps   │
     │         │         │ them                                                                                                │
     ├─────────┼─────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 2930fc3 │ F-81    │ the console's retired name (maître d') is carried in its _Avoid_ line, and the proof command        │
     │         │         │ reproduces — seven retired names, seven hits, every one the line that warns, 0 hits outside an      │
     │         │         │ _Avoid_ line                                                                                        │
     ├─────────┼─────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 7d5df0a │ F-82    │ the console leaves ### The crew (now five entries) for a new ### The strings, and its header is     │
     │         │         │ spelled maitre                                                                                      │
     ├─────────┼─────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 4a47494 │ F-84    │ the two lists of the fingerprint's facts in docs/03 agree — and a third list I found in the         │
     │         │         │ DomainVersion comment now agrees too                                                                │
     ├─────────┼─────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 91e06a6 │ F-85    │ the version fences the writer; the fingerprint is the tripwire — in five sentences across docs/03   │
     │         │         │ and the DomainVersion comment                                                                       │
     └─────────┴─────────┴─────────────────────────────────────────────────────────────────────────────────────────────────────┘

     The order is not the finding order, and the reason is in the report: the glossary's proof command counts retired words
     outside their _Avoid_ lines, and F-83's prose line was the one making that count non-zero, so it went first.

     Nothing but prose moved: core.DomainVersion is still 6, the gate table, the plan format and the command count are
     untouched, and TestDomainShapeIsVersioned passes with the same fingerprint.

     Two things flagged, not fixed

     - The sweep's count moved 78 → 77 — F-83 removed the one prose line, and the sweep command does not match maître d'. The
       pattern itself holds: seven retired names, seven _Avoid_ lines, nothing in prose.
     - docs/03 § 6 links to #the-version-guard, and no such heading exists — a dangling anchor. It predates the slice and is not
       one of the five findings, so I left it for the orchestrator (one line to fix).

     The mascot still has not landed — logo/ holds the art, untracked; the console header, the CLI help, the README and the
     Sous-Chef tab carry no drawing, and the goldens did not move. I touched none of them.

     Build · MiniMax-M3.1-Flash · 6m 28s · 129.5 tok/s

  ┃
  ┃  You are the **runner** — the same lane that wrote slice E. The orchestrator's pane id is
  ┃  `w1Z:p1`; use it verbatim in the notify command below.
  ┃
  ┃  ## The last fix cycle for slice E — one finding and two notes
  ┃
  ┃  The reviewer came back with **F-86 (medium)** and two LOW notes. The orchestrator
  ┃  ruled: **all three are fixed here.** F-86 is a beautiful finding — read it, because
  ┃  it is the kind of thing that is worth remembering.
  ┃
  ┃  ### F-86 (medium) — the instrument was blind, not the tree
  ┃
  ┃  > *the sweep's command and AC-04's list of old names omit* maître d'*, so one of the
  ┃  > seven renamed terms is outside the instrument's reach; the tree is clean (only dated
  ┃  > records), and the remedy is the command's: carry the name, exclude the glossary's
  ┃  > `_Avoid_` lines.*
  ┃
  ┃  **The tree is clean.** Nothing needs changing in the docs. **The command could not
  ┃  see one of the seven renamed terms** — *maître d'*, the console's old themed name.
  ┃
  ┃  > An instrument that cannot see a name is not a check. The finding is not "the tree
  ┃  > is dirty"; it is "the sweep would have passed whatever the tree said about
  ┃  > *maître d'*, because it never asked."
  ┃
  ┃  **The fix:**
  ┃
  ┃  1. **The sweep's command carries every one of the seven names** — *maître d'* among
  ┃     them — and it **knows the difference between a retired word and a warning about a
  ┃     retired word**: the glossary's `_Avoid_` lines name them on purpose, so the sweep
  ┃     excludes those lines, and a retired name outside one is drift.
  ┃  2. **Re-run the sweep with the complete command and show its output** in the commit
  ┃     message — the actual output, not a claim about it. The same rule that F-81 taught:
  ┃     a command offered as evidence must reproduce.
  ┃  3. **Say what the count is now, and reconcile it.** The reviewer counted **77 lines
  ┃     in 19 files** and found the arithmetic holds (78 → 77 is F-83's one prose line
  ┃     removed, nothing added). A command that now also carries *maître d'* will find the
  ┃     glossary's `_Avoid_` line — so the count will move, and the commit must explain
  ┃     the movement.
  ┃
  ┃  **The brief's list is the orchestrator's and is corrected** — *maître d'* is in it
  ┃  now. Yours is the command.
  ┃
  ┃  **No script, deliberately.** The orchestrator considered making the sweep a permanent
  ┃  tripwire and decided against it, and wants the reason on the record: the sweep is a
  ┃  **one-time act** over a tree that now has *legitimate* mentions of the old names —
  ┃  the legacy reader, the retired-role aliases, the dated records. A permanent
  ┃  tripwire would need a permanent exclusion list, and **a tripwire with a long
  ┃  exclusion list is a tripwire that has stopped seeing**. The repo already has better
  ┃  instruments for the live facts (the deny-list wire, the derived-views wire, the
  ┃  plan-format wire). Do not add a script.
  ┃
  ┃  ### The two LOW notes
  ┃
  ┃  | Note | Ruling |
  ┃  |---|---|
  ┃  | the console's entry sits under **"### The strings"** — the section says what the console *carries*, but the heading alone
  ┃  can be misread: the console **is not a string** | **fix it.** A glossary is a **taxonomy**, and a taxonomy that files the
  ┃  console under "the strings" is wrong in the same way F-82 was wrong when it filed the console among the crew. The heading
  ┃  must not claim the console is a string. You choose the heading — say which, and why. |
  ┃  | `docs/03` § 6 links to `#the-version-guard`, and **there is no such heading** — a broken pointer, pre-existing | **fix
  ┃  it.** The pointer must point at something that exists. The guard's own description is two lines above the link, and the
  ┃  guard lives in `internal/core/lock.go` (`checkDomainVersion`). Point it at whichever is right, and say which. |
  ┃
  ┃  ## The tree is shared — read this before you commit
  ┃
  ┃  **Another session is drawing the mascot in this repository** and has files in the tree
  ┃  (`logo/`).
  ┃
  ┃  | | |
  ┃  |---|---|
  ┃  | **Stage by name** | name your files explicitly. **Never `git add -A`**, never `git add .` |
  ┃  | **`logo/`, and anything you did not write** | leave it alone |
  ┃  | **`kitchen/`** | re-read before you commit; the other session may write to the plan |
  ┃
  ┃  **If a file changed that you did not change, stop and report it.**
  ┃
  ┃  ## Not yours — do not touch
  ┃
  ┃  | | |
  ┃  |---|---|
  ┃  | `docs/audit.md` | the reviewer's to write. It files F-86 and closes it. |
  ┃  | `docs/07-vocabulary.md`, `work/briefs/slice-e.md` | the orchestrator's. Both corrected. |
  ┃  | `README.md`, `work/README.md` | the orchestrator's. Already corrected. |
  ┃  | The gate table, the plan format, the command count (42), `core.DomainVersion` (6), `kitchen/activity.json` | none of the
  ┃  three touches them. If a fix would, **stop and report it**. |
  ┃
  ┃  ## The gate
  ┃
  ┃  ```bash
  ┃  make lint && gofmt -l . && make test && git status --short
  ┃  ```
  ┃
  ┃  Green.
  ┃
  ┃  ## Two commits — one per change
  ┃
  ┃  ```
  ┃  Per-AC evidence:
  ┃  - F-86: the sweep's command carries all seven retired names and excludes the glossary's _Avoid_ lines; the output is in the
  ┃  commit body
  ┃  ```
  ┃
  ┃  ## The report — the receipt comes first
  ┃
  ┃  ```bash
  ┃  mkdir -p work/inbox
  ┃  cat > work/inbox/runner-slice-e-fix-2.md <<EOF
  ┃  Runner: slice E fix 2 done. Commits: <n>. Fixed: F-86 and the two notes. Sweep count now: <n>, and it moves by <why>.
  ┃  Notes: <anything the reviewer should know>.
  ┃  EOF
  ┃  git add work/inbox && git commit -q -m "inbox: runner reports slice-e fix 2"
  ┃  herdr agent prompt "w1Z:p1" "Runner: slice E fix 2 done. Commits: <n>. Fixed: F-86 and the two notes. Sweep count now: <n>,
  ┃  and it moves by <why>. Notes: <anything the reviewer should know>."
  ┃  ```
  ┃
  ┃  **One write per file, and commit as soon as the file is coherent.**
  ┃

     + Thought · 33.6s

     → Explored: 2 reads

     + Thought · 19.1s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "internal/core/.*\.go\|\.go)" docs/*.md | head -10; echo "=== lock.go gu
  ┃  ard ==="; grep -n "checkDomainVersion" -A 20 internal/core/lock.go | head -40
  ┃
  ┃  (43 earlier lines)
  ┃  75-  plan := PlanDomainVersion(j)
  ┃  76-  if plan > DomainVersion {
  ┃  77-    return fmt.Errorf("this plan was written under domain version %d and this tool speaks version %d: upgrade the tool,
  ┃  or read the plan without writing to it", plan, DomainVersion)
  ┃  78-  }
  ┃  79-  return nil
  ┃  80-}
  ┃

     + Thought · 25.3s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "func sectionOf" -A 20 internal/core/docs_test.go
  ┃
  ┃  (14 earlier lines)
  ┃  38-// A sub-heading's section ends at the next `###` as well as the next `##` — without
  ┃  39-// that, a check reads the tables of the sections that follow it.
  ┃  40-func subsectionOf(doc, heading string) string {
  ┃  41-  rest := sectionOf(doc, heading)
  ┃  42-  if i := strings.Index(rest, "\n### "); i >= 0 {
  ┃  43-    rest = rest[:i]
  ┃  44-  }
  ┃

     + Thought · 1.1s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "subsectionOf(\|sectionOf(" internal/core/*_test.go internal/cli/*_test.
  ┃  go | grep -v "func " | head -20
  ┃
  ┃  (6 earlier lines)
  ┃  internal/core/docs_test.go:616:  section := sectionOf(doc, "## The operations")
  ┃  internal/core/docs_test.go:1396:  for _, m := range edge.FindAllStringSubmatch(sectionOf(doc, "### 3.2 Transitions"), -1) {
  ┃  internal/core/docs_test.go:1493:  for _, m := range edge.FindAllStringSubmatch(sectionOf(doc, "### 3.2 Transitions"), -1) {
  ┃  internal/core/docs_test.go:1508:  section := sectionOf(doc, "## 4. Gates")
  ┃  internal/core/planformat_test.go:34:  section := sectionOf(string(raw), "## 6. Plan format")
  ┃  internal/core/planformat_test.go:67:  section := sectionOf(doc, "## 6. Plan format")
  ┃

     + Thought · 18.5s

     Now F-86 — the instrument. Let me build the complete sweep command and run it.

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cat > /tmp/sweep.sh <<'SH'
  ┃  grep -rnE --exclude-dir=.git --exclude-dir=work --exclude-dir=logo --exclude-dir=bin --exclude-dir=kitchen \…
  ┃
  ┃  (16 earlier lines)
  ┃     1 ./internal/core/approval_test.go
  ┃     1 ./docs/04-operations.md
  ┃  === total lines ===
  ┃        72
  ┃  === files ===
  ┃        18
  ┃

     + Thought · 3.2s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sh /tmp/sweep.sh | grep -c . ; echo "=== without the glossary exclusion (the comp
  ┃  lete pattern alone) ==="; grep -rnE --exclude-dir=.git --exclude-dir=work --exclude-dir=logo --exclude-dir=bin --exc…
  ┃
  ┃  72
  ┃  === without the glossary exclusion (the complete pattern alone) ===
  ┃        78
  ┃  === the glossary's lines in that set ===
  ┃  GLOSSARY.md:161:_Avoid_: maître d', TUI, dashboard, frontend, UI
  ┃

     + Thought · 7.3s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git commit -q --allow-empty -F - <<'EOF'
  ┃  F-86: the instrument was blind, not the tree…
  ┃
  ┃  cea2380 F-86: the instrument was blind, not the tree
  ┃

     + Thought · 21.8s

     Now the two notes — the glossary's heading and the broken pointer.

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
  ┃  import io…
  ┃
  ┃  (32 earlier lines)
  ┃  prefixes, the directories, the root files and the field names on disk. A rule outside
  ┃  ...
  ┃  43:### The plan
  ┃  89:### The crew
  ┃  114:### The work
  ┃  151:### The surfaces
  ┃

     + Thought · 841ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '149,165p' GLOSSARY.md && echo "=== gate ===" && make lint && gofmt -l . &
  ┃  & make test 2>&1 | tail -5 && git status --short
  ┃
  ┃  (21 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/core  22.562s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  66.147s
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃   M GLOSSARY.md
  ┃   M docs/03-domain-model.md
  ┃  ?? logo/
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add GLOSSARY.md docs/03-domain-model.md && git commit -q -F - <<'EOF'
  ┃  notes: the heading does not claim the console is a string, and the pointer lands…
  ┃
  ┃  08fa88c notes: the heading does not claim the console is a string, and the pointer lands
  ┃  cea2380 F-86: the instrument was blind, not the tree
  ┃  7e12f53 work: slice E's last fix cycle — the instrument was blind, not the tree
  ┃  ?? logo/
  ┃

     + Thought · 695ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -6 && echo "=== 7e12f53 ===" && git show --stat 7e12f53 | tail
  ┃  -12
  ┃
  ┃  (10 earlier lines)
  ┃
  ┃      work: slice E's last fix cycle — the instrument was blind, not the tree
  ┃
  ┃   work/briefs/fix-e2.md  | 106 +++++++++++++++++++++++++++++++++++++++++++++++++
  ┃   work/briefs/slice-e.md |  11 ++++-
  ┃   2 files changed, 116 insertions(+), 1 deletion(-)
  ┃

     + Thought · 316ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "F-86" docs/audit.md | head -3 && echo "=== sweep after the notes ===" &&
  ┃   sh /tmp/sweep.sh | wc -l && sh /tmp/sweep.sh | cut -d: -f1 | sort -u | wc -l && echo "=== gate ===" && make lint &&…
  ┃
  ┃  (5 earlier lines)
  ┃  go vet ./...
  ┃  ok    github.com/lthiagol/maitre/internal/cli  2.463s
  ┃  ok    github.com/lthiagol/maitre/internal/core  (cached)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  (cached)
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃  ?? logo/
  ┃

     + Thought · 1.6s

  ┃
  ┃  # Wrote work/inbox/runner-slice-e-fix-2.md
  ┃
  ┃     1 Runner: slice E fix 2 done. Commits: 2. Fixed: F-86 and the two notes. Sweep count now: 72 lines in 18 files, and it
  ┃       moves by −5 from the reviewer's 77: the old command counted five of the glossary's six `_Avoid_` lines as drift, and
  ┃       the complete command sees all six and excludes all six. Notes: the heading is **### The surfaces** — the console is a
  ┃       surface (ADR-0003's word, and the entry's own first word), and "the strings" claimed the console *is* one; the broken
  ┃       pointer now lands on a new **§ 6.2 The version guard**, the document's own description, with the implementation named
  ┃       in the same paragraph; no script, for the reason the orchestrator gave.
  ┃     2
  ┃     3 ## The commits
  ┃     4
  ┃     5 | Commit | What |
  ┃     6 |---|---|
  ┃     7 | `cea2380` | **F-86** — the sweep's instrument: the complete command, its output, and the reconciliation. No file
  ┃       changed: the tree was clean |
  ┃     8 | `08fa88c` | the two notes — the glossary's heading and the broken pointer |
  ┃     9
  ┃    10 ## F-86 — the instrument, not the tree
  ┃    11
  ┃    12 **The tree is clean.** Nothing in the docs needed changing, and the F-86 commit
  ┃    13 changes no file. What was wrong was the instrument: the pattern omitted *maître d'*,
  ┃    14 so one of the seven renamed terms was outside its reach — the sweep would have
  ┃    15 passed whatever the tree said about the console's old name, because it never asked.
  ┃    16
  ┃    17 The command now carries all seven retired names **and** knows the difference between
  ┃    18 a retired word and a **warning** about one: the glossary's `_Avoid_` lines name them
  ┃    19 on purpose, so the sweep excludes those lines, and a retired name outside one is
  ┃    20 drift.
  ┃    21
  ┃    22     grep -rnE --exclude-dir=.git --exclude-dir=work --exclude-dir=logo --exclude-dir=bin --exclude-dir=kitchen \
  ┃    23       --include="*.go" --include="*.md" --include="*.py" -i \
  ┃    24       "kitchen brigade|kitchen-brigade|chopboard|\bexpo\b|\btaster\b|line-cook|line cook|\bspecials?\b|prep[ -]note|\
  ┃       bproofs?\b|maître d'" . \
  ┃    25       | grep -v "^./docs/audit" \
  ┃    26       | grep -v "^./docs/audits/" \
  ┃    27       | grep -v "^./research/" \
  ┃    28       | grep -v "^./docs/decisions/0004-naming.md" \
  ┃    29       | grep -v "^./docs/decisions/0005-vocabulary.md" \
  ┃    30       | grep -v "^./docs/decisions/0009-vocabulary.md" \
  ┃    31       | grep -v "^./docs/07-vocabulary.md" \
  ┃    32       | grep -v "^./prototypes/tab-" \
  ┃    33       | grep -v "^./GLOSSARY.md:[0-9]*:_Avoid_:"
  ┃    34
  ┃    35 **Its output** — the command's, piped through a counter:
  ┃    36
  ┃    37     $ sh sweep.sh | cut -d: -f1 | sort | uniq -c | sort -rn
  ┃    38          15 ./internal/core/legacy.go
  ┃    39          11 ./internal/core/roles_test.go
  ┃    40          10 ./internal/core/model.go
  ┃    41           9 ./docs/03-domain-model.md
  ┃    42           5 ./internal/core/lifecycle.go
  ┃    43           4 ./internal/ui/console_test.go
  ┃    44           3 ./internal/core/legacy_test.go
  ┃    45           2 ./scripts/stamp-plan/main.go
  ┃    46           2 ./internal/core/views.go
  ┃    47           2 ./internal/core/backlog_test.go
  ┃    48           2 ./internal/core/backlog_concurrency_test.go
  ┃    49           1 ./prototypes/README.md
  ┃    50           1 ./internal/core/plan_ids_test.go
  ┃    51           1 ./internal/core/migration_test.go
  ┃    52           1 ./internal/core/gates_test.go
  ┃    53           1 ./internal/core/docs_test.go
  ┃    54           1 ./internal/core/approval_test.go
  ┃    55           1 ./docs/04-operations.md
  ┃    56
  ┃    57     $ sh sweep.sh | wc -l
  ┃    58     72
  ┃    59     $ sh sweep.sh | cut -d: -f1 | sort -u | wc -l
  ┃    60     18
  ┃    61
  ┃    62 **The count, reconciled.** The reviewer counted **77 lines in 19 files**:
  ┃    63
  ┃    64 | Count | What it is |
  ┃    65 |--:|---|
  ┃    66 | 77 | the old command — six names, and the glossary's `_Avoid_` lines counted as drift |
  ┃    67 | 78 | the complete command **without** the glossary exclusion: it sees the console's name, and the glossary's six
  ┃       `_Avoid_` lines are in the set |
  ┃    68 | 72 | the complete command as it now stands — the six `_Avoid_` lines excluded |
  ┃    69
  ┃    70 77 − 5 = 72: the old command saw **five** of the glossary's six `_Avoid_` lines, all
  ┃    71 but the console's — which it could not name — and the new command sees all six and
  ┃    72 excludes all six. The name the old instrument could not see adds **no** drift line:
  ┃    73 *maître d'* lives in the glossary's `_Avoid_` line and in the ADR records (0004,
  ┃    74 0005, 0009), and both are excluded.
  ┃    75
  ┃    76 The count is **72 in 18 files** after the two notes landed too — the headings carry
  ┃    77 no retired names.
  ┃    78
  ┃    79 **No script, deliberately.** The sweep is a **one-time act** over a tree that now has
  ┃    80 *legitimate* mentions of the old names — the legacy reader, the retired-role
  ┃    81 aliases, the dated records. A permanent tripwire would need a permanent exclusion
  ┃    82 list, and **a tripwire with a long exclusion list is a tripwire that has stopped
  ┃    83 seeing**. The repo already has better instruments for the live facts: the deny-list
  ┃    84 wire, the derived-views wire, the plan-format wire. The command lives here and in
  ┃    85 `cea2380` — it is the record of the act, not a guard.
  ┃    86
  ┃    87 ## The two notes
  ┃    88
  ┃    89 **The heading is `### The surfaces`.** A glossary is a **taxonomy**, and *The strings*
  ┃    90 claimed the console *is* a string — it *carries* strings. The console is a
  ┃    91 **surface**: ADR-0003's own word for it, and the first word of the entry's own
  ┃    92 paragraph. The section still says what the console carries (the window title, the
  ┃    93 tab bar, the CLI's help) — the heading no longer claims the console is one.
  ┃    94
  ┃    95     $ grep -n "^### " GLOSSARY.md
  ┃    96     43:### The plan
  ┃    97     89:### The crew
  ┃    98     114:### The work
  ┃    99     151:### The surfaces
  ┃   100
  ┃   101 **The pointer lands.** `docs/03` § 6 linked to `#the-version-guard`, and there was
  ┃   102 no such heading — the guard was described in prose at the end of § 6.1 with nothing
  ┃   103 to land on. The description now has a heading of its own, **§ 6.2 The version
  ┃   104 guard**, and the link follows it: `[the version guard](#62-the-version-guard)`.
  ┃   105
  ┃   106 I chose the **document** over the code, because a reader of the spec wants the
  ┃   107 spec's words; the implementation is named in the same paragraph — `core.DomainVersion`
  ┃   108 for the number and `Store.checkDomainVersion` in `internal/core/lock.go` for the
  ┃   109 guard that refuses the write. The version history that follows is the guard's story:
  ┃   110 every version paragraph ends with why the version moved, and that is what the guard
  ┃   111 fences.
  ┃   112
  ┃   113 ## The gate
  ┃   114
  ┃   115     $ make lint && gofmt -l . && make test && git status --short
  ┃   116     go vet ./...
  ┃   117     go test ./...
  ┃   118     ok    github.com/lthiagol/maitre/internal/cli  2.463s
  ┃   119     ok    github.com/lthiagol/maitre/internal/core  (cached)
  ┃   120     ok    github.com/lthiagol/maitre/internal/ui  (cached)
  ┃   121     ?? logo/
  ┃   122
  ┃   123 `gofmt -l .` printed nothing. The only thing in `git status` is `logo/` — the
  ┃   124 mascot's art, untracked and not mine.
  ┃   125
  ┃   126 Untouched, as the brief required: the gate table, the plan format, the command count
  ┃   127 (42), `core.DomainVersion` (6) and `kitchen/activity.json`.
  ┃   128
  ┃   129 ## The mascot
  ┃   130
  ┃   131 **Still not landed.** `logo/` holds the art, untracked; the console's header, the
  ┃   132 CLI's help, the README and the Sous-Chef tab carry no drawing. Nothing in this cycle
  ┃   133 touched them, and the goldens did not move — no console output changed.
  ┃   134
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/inbox && git commit -q -m "inbox: runner reports slice-e fix 2" && g
  ┃  it log --oneline -4 && git status --short
  ┃
  ┃  d0e32bd inbox: runner reports slice-e fix 2
  ┃  08fa88c notes: the heading does not claim the console is a string, and the pointer lands
  ┃  cea2380 F-86: the instrument was blind, not the tree
  ┃  7e12f53 work: slice E's last fix cycle — the instrument was blind, not the tree
  ┃  ?? logo/
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1Z:p1" "Runner: slice E fix 2 done. Commits: 2. Fixed: F-86
  ┃  and the two notes. Sweep count now: 72 lines in 18 files, and it moves by -5 from the reviewer's 77 — the old comman…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":
  ┃  "herdr:opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"done","completion_seq":491,"cwd":"/Users/thiago/
  ┃  Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,
  ┃  "screen_detection_skipped":true,"state_change_seq":491,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722",
  ┃  "terminal_title":"OC | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":
  ┃  "w1Z"},"type":"agent_prompted"}}
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make lint && gofmt -l . && make test 2>&1 | tail -4 && git status --short && echo
  ┃   "=== the sweep, one last time ===" && sh /tmp/sweep.sh | wc -l
  ┃
  ┃  (2 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/core  (cached)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  (cached)
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃  ?? logo/
  ┃  === the sweep, one last time ===
  ┃        72
  ┃

     The last fix cycle is done — two commits, gate green, orchestrator notified.

     ┌─────────┬───────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Commit  │ What                                                                                                          │
     ├─────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ cea2380 │ F-86 — the instrument. No file changed: the tree was clean. The commit body carries the complete command and  │
     │         │ its actual output                                                                                             │
     ├─────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 08fa88c │ the two notes — the glossary's heading and the broken pointer                                                 │
     └─────────┴───────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

     F-86 — the instrument, not the tree

     The pattern now carries all seven retired names (maître d' among them) and knows the difference between a retired word and
     a warning about one: the glossary's _Avoid_ lines are excluded, and a retired name outside one is drift.

     The count, reconciled: 77 → 72 lines in 18 files.

     ┌──────────────────────────┬──────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Count                    │ What it is                                                                                   │
     ├──────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 77                       │ the old command — six names, the glossary's _Avoid_ lines counted as drift                   │
     ├──────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 78                       │ the complete command without the glossary exclusion                                          │
     ├──────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────┤
     │ 72                       │ the complete command as it now stands                                                        │
     └──────────────────────────┴──────────────────────────────────────────────────────────────────────────────────────────────┘

     77 − 5 = 72: the old command saw five of the glossary's six _Avoid_ lines — all but the console's, which it could not name
     — and the new one sees all six and excludes all six. The name the old instrument could not see adds no drift line: maître
     d' lives only in the glossary's _Avoid_ line and in the ADR records.

     No script, for the orchestrator's reason: the sweep is a one-time act over a tree with legitimate mentions, and a tripwire
     with a long exclusion list is a tripwire that has stopped seeing. The command is the record of the act, in the commit body
     and in the receipt.

     The two notes

     - The heading is ### The surfaces — the console is a surface (ADR-0003's word, and the entry's own first word). "The
       strings" claimed the console is one; it carries strings.
     - The pointer lands on a new ### 6.2 The version guard in docs/03 § 6 — the document's own description, which had no
       heading to land on. I chose the document over the code because a reader of the spec wants the spec's words; the
       implementation (core.DomainVersion, Store.checkDomainVersion in internal/core/lock.go) is named in the same paragraph.

     Untouched as required: the gate table, the plan format, the command count, core.DomainVersion (6), kitchen/activity.json.
     The mascot still has not landed — logo/ is untracked and the goldens did not move.

     Build · MiniMax-M3.1-Flash · 4m 47s · 123.0 tok/s

  ┃
  ┃
  ┃
  ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
  ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
  ~/Code/Pessoal/maitre:master                                                                     299.2K (30%)  ctrl+p commands

```
