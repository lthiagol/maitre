# The runner's lane — slice H, captured 2026-10-28

Dumped from pane `w1Z:pJ` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that wrote slice H: the gate, the tools, the CI, and the scripts that hold the
facts (the budget, the pins, the goldens' claims and their commit).
The review that judged it is in [slice-h-reviewer.md](slice-h-reviewer.md).

```
  ┃  === every .golden file in the repo:
  ┃  ./internal/ui/testdata/console-stress.golden
  ┃  ./internal/ui/testdata/console-recipes.golden
  ┃

     + Thought · 1m 57s



  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "goldenPath\|stressGoldenPath" internal/ui/ | grep -v "^Binary"
  ┃
  ┃  (4 earlier lines)
  ┃  internal/ui/golden_test.go:17:const goldenPath = "testdata/console-recipes.golden"
  ┃  internal/ui/golden_test.go:19:// stressGoldenPath is the same, against the pathological plan.
  ┃  internal/ui/golden_test.go:20:const stressGoldenPath = "testdata/console-stress.golden"
  ┃  internal/ui/golden_test.go:47:    if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
  ┃  internal/ui/golden_test.go:50:    if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
  ┃  internal/ui/golden_test.go:57:  want, err := os.ReadFile(goldenPath)
  ┃  internal/ui/golden_test.go:75:  want, err := os.ReadFile(goldenPath)
  ┃

     + Thought · 792ms

  ┃
  ┃  # Wrote scripts/check-golden-claims.sh
  ┃
  ┃    1 #!/bin/sh
  ┃    2 # The golden files' claims.
  ┃    3 #
  ┃    4 # usage: check-golden-claims.sh mark <dir> | report <dir>
  ┃    5 #
  ┃    6 # A golden file that no test writes is a layout nothing pins — and the check
  ┃    7 # that says "the goldens are current" would pass without touching it.
  ┃    8 #
  ┃    9 # The claim is derived, not declared: `mark` moves every golden's timestamp
  ┃   10 # into the past, the golden run rewrites the ones a test claims, and
  ┃   11 # `report` names the ones nobody wrote. No manifest of test names, and no
  ┃   12 # scan for a path in a test file — the paths live in constants, and a scan
  ┃   13 # would be a guess.
  ┃   14 #
  ┃   15 # A failure is loud and names the file. Two causes, one message: no test
  ┃   16 # claims the file, or the golden run's filter misses the test that does —
  ┃   17 # and the second is a real hole too, because `make goldens` would not
  ┃   18 # regenerate the file either.
  ┃   19 set -u
  ┃   20
  ┃   21 # A timestamp in the past, so "not newer than this" means "not written".
  ┃   22 MARK=200001010000
  ┃   23
  ┃   24 phase=${1:-}
  ┃   25 dir=${2:-}
  ┃   26
  ┃   27 case "$phase" in
  ┃   28 mark)
  ┃   29   find "$dir" -name '*.golden' -type f -exec touch -t "$MARK" {} +
  ┃   30   ;;
  ┃   31 report)
  ┃   32   ref=$(mktemp)
  ┃   33   touch -t "$MARK" "$ref"
  ┃   34   orphans=$(find "$dir" -name '*.golden' -type f ! -newer "$ref")
  ┃   35   rm -f "$ref"
  ┃   36   if [ -n "$orphans" ]; then
  ┃   37     echo "these golden files are claimed by no test:" >&2
  ┃   38     printf '%s\n' "$orphans" >&2
  ┃   39     echo "No test in the golden run writes them: either no test claims them, or the run's filter misses the test that
  ┃      does." >&2
  ┃   40     exit 1
  ┃   41   fi
  ┃   42   echo "  every golden file is claimed by a test"
  ┃   43   ;;
  ┃   44 *)
  ┃   45   echo "usage: check-golden-claims.sh mark <dir> | report <dir>" >&2
  ┃   46   exit 2
  ┃   47   ;;
  ┃   48 esac
  ┃   49
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && chmod +x scripts/check-golden-claims.sh && ls -l internal/ui/testdata/ && echo "=
  ┃  == mark:" && scripts/check-golden-claims.sh mark internal/ui/testdata && ls -l internal/ui/testdata/ && echo "=== re…
  ┃
  ┃  No tool named "bash" is currently available. Please use a tool from the available tool list.
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && chmod +x scripts/check-golden-claims.sh && ls -l internal/ui/testdata/ && echo "=
  ┃  == mark:" && scripts/check-golden-claims.sh mark internal/ui/testdata && ls -l internal/ui/testdata/ && echo "=== re…
  ┃
  ┃  (7 earlier lines)
  ┃  === report (nothing ran, so both are orphans):
  ┃  these golden files are claimed by no test:
  ┃  internal/ui/testdata/console-stress.golden
  ┃  internal/ui/testdata/console-recipes.golden
  ┃  No test in the golden run writes them: either no test claims them, or the run's filter misses the test that does.
  ┃  exit 1
  ┃

     + Thought · 1.3s

  ┃
  ┃  ← Edit Makefile
  ┃
  ┃    105   goldens-check: GOLDEN_CHECK := yes                     105   goldens-check: GOLDEN_CHECK := yes
  ┃    106                                                          106
  ┃    107   goldens goldens-check:                                 107   goldens goldens-check:
  ┃    108     $(call header,$@)                                    108     $(call header,$@)
  ┃                                                                 109 +   @if [ "$(GOLDEN_CHECK)" = yes ]; then scripts/check-
  ┃                                                                       golden-claims.sh mark $(GOLDENS); fi
  ┃    109     GO=$(GO) scripts/go-test.sh "$(GOLDEN_FLAGS)" "$(    110     GO=$(GO) scripts/go-test.sh "$(GOLDEN_FLAGS)" "$(
  ┃          GOLDEN_PKG)"                                                 GOLDEN_PKG)"
  ┃                                                                 111 +   @if [ "$(GOLDEN_CHECK)" = yes ]; then scripts/check-
  ┃                                                                       golden-claims.sh report $(GOLDENS); fi
  ┃    110     @if [ "$(GOLDEN_CHECK)" = yes ]; then \              112     @if [ "$(GOLDEN_CHECK)" = yes ]; then \
  ┃    111       git diff --exit-code HEAD -- $(GOLDENS) || { \     113       git diff --exit-code HEAD -- $(GOLDENS) || { \
  ┃    112         echo "the goldens changed — read the diff before 114         echo "the goldens changed — read the diff before
  ┃          committing it." >&2; exit 1; \                               committing it." >&2; exit 1; \
  ┃    113       }; \                                               115       }; \
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make goldens-check; echo "exit $?"; git status --short
  ┃
  ┃  ==> goldens-check
  ┃  GO=go scripts/go-test.sh "-run Golden -update" "./internal/ui"
  ┃  ok    github.com/lthiagol/maitre/internal/ui  5.220s
  ┃    every golden file is claimed by a test
  ┃  exit 0
  ┃   M Makefile
  ┃  ?? scripts/check-golden-claims.sh
  ┃

     + Thought · 837ms



  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -i '' 's/func TestConsoleMatchesTheGolden(/func TestConsoleMatchesTheSnapshot
  ┃  (/' internal/ui/golden_test.go && grep -n "func TestConsoleMatches" internal/ui/golden_test.go && echo "=== goldens-…
  ┃
  ┃  (9 earlier lines)
  ┃  make: *** [goldens-check] Error 1
  ┃  === exit 2
  ┃  === the tree:
  ┃   M Makefile
  ┃   M internal/ui/golden_test.go
  ┃  ?? scripts/check-golden-claims.sh
  ┃

     + Thought · 874ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git checkout -- internal/ui/golden_test.go && make goldens-check && git status --
  ┃  short && echo "(only my own files changed)"
  ┃
  ┃  ==> goldens-check
  ┃  GO=go scripts/go-test.sh "-run Golden -update" "./internal/ui"
  ┃  ok    github.com/lthiagol/maitre/internal/ui  5.220s
  ┃    every golden file is claimed by a test
  ┃   M Makefile
  ┃  ?? scripts/check-golden-claims.sh
  ┃  (only my own files changed)
  ┃



  ┃
  ┃  ← Edit scripts/check-tool-versions.sh
  ┃
  ┃    53     tool=${2:-}                                           53     tool=${2:-}
  ┃    54     pin=${3:-}                                            54     pin=${3:-}
  ┃    55     want=${pin#v}                                         55     want=${pin#v}
  ┃    56     # Every mention that names a version is checked — not 56     # Every mention that names a version is checked — not
  ┃         only the first.                                              only the first.
  ┃    57 -   versions=$(grep -F "$tool" "$doc" | grep -E 'v?[0-9]  57 +   # Only the pin's shape counts: the tool reports two
  ┃         +\.[0-9]+\.[0-9]+' | grep -oE 'v?[0-9]+\.[0-9]+\.[0-9]       versions (the
  ┃         +' || true)
  ┃                                                                 58 +   # release and the module), the pin is the module one,
  ┃                                                                      and a line that
  ┃                                                                 59 +   # carries the other version must not fail the gate.
  ┃                                                                 60 +   versions=$(grep -F "$tool" "$doc" | grep -E 'v[0-9]+\
  ┃                                                                      .[0-9]+\.[0-9]+' | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+'
  ┃                                                                      || true)
  ┃    58     if [ -z "$versions" ]; then                           61     if [ -z "$versions" ]; then
  ┃    59       echo "$doc names no version for $tool." >&2         62       echo "$doc names no version for $tool." >&2
  ┃    60       echo "Change $doc to name $pin." >&2                63       echo "Change $doc to name $pin." >&2
  ┃    61       exit 1                                              64       exit 1
  ┃

     → Explored: 1 read

     + Thought · 641ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp docs/testing.md /tmp/t.bak && echo "=== (a) a line with the release version's
  ┃  shape:" && python3 -c "…
  ┃
  ┃  (12 earlier lines)
  ┃  scripts/check-tool-versions.sh doc docs/testing.md staticcheck v0.8.2
  ┃  docs/testing.md names staticcheck 0.8.1; the pin is v0.8.2.
  ┃  Change docs/testing.md to name v0.8.2.
  ┃  make: *** [versions-check] Error 1
  ┃  exit 2
  ┃  (restored)
  ┃

     + Thought · 3.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
  ┃  def edit(path, old, new):…
  ┃
  ┃  (6 earlier lines)
  ┃    docs/testing.md names govulncheck 1.9.0, as pinned
  ┃   M AGENTS.md
  ┃   M Makefile
  ┃   M docs/testing.md
  ┃   M scripts/check-tool-versions.sh
  ┃  ?? scripts/check-golden-claims.sh
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check > /tmp/c3.log 2>&1; echo "check exit: $?"; grep "^==>" /tmp/c3.log; ec
  ┃  ho "---"; tail -4 /tmp/c3.log; git status --short
  ┃
  ┃  (15 earlier lines)
  ┃  ==> tidy-check
  ┃   M AGENTS.md
  ┃   M Makefile
  ┃   M docs/testing.md
  ┃   M scripts/check-tool-versions.sh
  ┃  ?? scripts/check-golden-claims.sh
  ┃

     ◈ Instructions updated: core/instructions

     + Thought · 452ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && { time make check-race ; } > /tmp/cr3.log 2>&1; echo "check-race exit: $?"; grep
  ┃  "^==>" /tmp/cr3.log; tail -3 /tmp/cr3.log; echo "=== git status --short:"; git status --short; echo "(clean if empty…
  ┃
  ┃  (15 earlier lines)
  ┃   M AGENTS.md
  ┃   M Makefile
  ┃   M docs/testing.md
  ┃   M scripts/check-tool-versions.sh
  ┃  ?? scripts/check-golden-claims.sh
  ┃  (clean if empty)
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add Makefile scripts/check-golden-claims.sh scripts/check-tool-versions.sh AG
  ┃  ENTS.md docs/testing.md && git commit -q -F - <<'EOF'…
  ┃
  ┃  6ddc41b fix: a golden file nobody claims is a layout nothing pins
  ┃  (clean)
  ┃

  ┃
  ┃  # Wrote work/inbox/runner-slice-h.md
  ┃
  ┃    1 Runner: slice H fix 3 done. Commits: 1 (plus this receipt). The orphan check: **built**, and it is derived, not
  ┃      declared. The doc check compares the pinned version only. Notes: `go-test.sh` passing the package first is accepted as
  ┃      it stands — no change, and the next reader should not "fix" it.
  ┃    2
  ┃    3 Per-AC evidence:
  ┃    4
  ┃    5 - **the orphan check — built.** The mechanism is runtime-derived, not a source scan. `scripts/check-golden-claims.sh
  ┃      mark` moves every `.golden` file's timestamp into the past; the golden run rewrites the ones a test claims; `report`
  ┃      names the ones nobody wrote. There is no manifest of test names and no scan for a path in a test file — the paths live
  ┃      in constants (`goldenPath`, `stressGoldenPath`), so a reference-scan would be a guess, and a guess in a gate is worse
  ┃      than the hole.
  ┃    6   **Proved:** with `TestConsoleMatchesTheGolden` renamed away, `make goldens-check` **fails** and names the file —
  ┃      *"these golden files are claimed by no test: internal/ui/testdata/console-recipes.golden / No test in the golden run
  ┃      writes them: either no test claims them, or the run's filter misses the test that does."* The rename was reverted and
  ┃      the check is green again.
  ┃    7   **When it cries wolf, and why that is not a false alarm:** it fires when a golden file is not written by the golden
  ┃      run — either no test claims it, or the run's filter (`-run Golden`) misses the test that does. The second case is a
  ┃      real hole too: `make goldens` would not regenerate that file either. The failure is loud and names the file.
  ┃    8 - **the doc check compares the pinned version only** — only the pin's shape (`vN.N.N`, the module version) is
  ┃      collected; the release version (`2026.2.2`) is not the pin's shape and is not compared. **Proved:** a line in `docs/
  ┃      testing.md` mentioning `staticcheck 2026.2.2` does **not** fail the check, and a second mention carrying the pinned
  ┃      shape with the wrong value (`v0.8.1`) still **fails** — *"docs/testing.md names staticcheck 0.8.1; the pin is v0.8.2.
  ┃      "*
  ┃    9 - **`go-test.sh` passes the package first** — accepted as it stands. The ordering is deliberate (`-update` belongs to
  ┃      the test binary), and a flag nobody knows fails loudly in every binary, which is fine. **No change.**
  ┃   10
  ┃   11 The gate: `make check` green (1m14s), `make check-race` green (3m12s), `git status --short` clean. Counts unchanged: 0
  ┃      / 0 / 28. `core.DomainVersion` is still 7.
  ┃   12
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/inbox && git commit -q -m "inbox: runner reports slice-h fix 3" && g
  ┃  it log --oneline -2 && git status --short && echo "(clean)"
  ┃
  ┃  1a82163 inbox: runner reports slice-h fix 3
  ┃  6ddc41b fix: a golden file nobody claims is a layout nothing pins
  ┃  (clean)
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1Z:p1" "Runner: slice H fix 3 done. Commits: 1 (plus the rec
  ┃  eipt). The orphan check: BUILT, and derived rather than declared — check-golden-claims.sh mark moves every .golden t…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":
  ┃  "herdr:opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"done","completion_seq":656,"cwd":"/Users/thiago/
  ┃  Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,
  ┃  "screen_detection_skipped":true,"state_change_seq":656,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722",
  ┃  "terminal_title":"OC | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":
  ┃  "w1Z"},"type":"agent_prompted"}}
  ┃

     Fix 3 is done. The receipt is committed at work/inbox/runner-slice-h.md and pane w1Z:p1 has been nudged.

     1. The orphan check — built, and derived

     I did not build a reference-scan: the golden paths live in constants (goldenPath, stressGoldenPath), so scanning the test
     sources for a path would be a guess — and a guess in a gate is worse than the hole.

     Instead the claim is derived at runtime, with no manifest:

     ┌──────────────────────────────────────────────────┬──────────────────────────────────────────────────────────────────────┐
     │ Step                                             │ What                                                                 │
     ├──────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────┤
     │ check-golden-claims.sh mark                      │ moves every .golden file's timestamp into the past                   │
     ├──────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────┤
     │ the golden run                                   │ rewrites the ones a test claims                                      │
     ├──────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────┤
     │ check-golden-claims.sh report                    │ names the ones nobody wrote                                          │
     └──────────────────────────────────────────────────┴──────────────────────────────────────────────────────────────────────┘

     Proved: with TestConsoleMatchesTheGolden renamed away, make goldens-check fails and names the file:

     these golden files are claimed by no test:
     internal/ui/testdata/console-recipes.golden
     No test in the golden run writes them: either no test claims them, or the run's filter misses the test that does.

     Renamed back → green again.

     When it cries wolf: it fires when a golden is not written by the golden run — either nothing claims it, or the filter (-run
     Golden) misses the writer. The second is a real hole too: make goldens would not regenerate that file either. The failure
     is loud and names the file.

     2. The doc check compares the pinned version only

     Only the pin's shape (vN.N.N — the module version) is collected. Proved: a line mentioning staticcheck 2026.2.2 (the
     release version) does not fail; a second mention carrying the pinned shape with the wrong value (v0.8.1) still fails.

     3. Accepted as it stands

     go-test.sh passing the package first — no change, and the receipt says so so the next reader does not "fix" it.

     The gate

     make check 1m14s green · make check-race 3m12s green · git status --short clean. Counts unchanged at 0 / 0 / 28; core.
     DomainVersion still 7.

     Build · MiniMax-M3.1-Flash · 8m 47s · 121.3 tok/s

  ┃
  ┃  You are the **runner** — the same lane that wrote slice H. The orchestrator's pane id is
  ┃  `w1Z:p1`; use it verbatim in the notify command below.
  ┃
  ┃  ## The last round for slice H — two notes
  ┃
  ┃  The review came back **`ok`** for the fourth time. The orchestrator ruled: **these two
  ┃  are fixed here, and this is the last round for slice H.** The other two notes are
  ┃  carried to the ledger.
  ┃
  ┃  ### 1. The gate cannot see an untracked golden
  ┃
  ┃  > `git diff --exit-code HEAD -- internal/ui/testdata` does **not** see **untracked**
  ┃  > files, so a new golden that a test writes but nobody commits **passes the gate**.
  ┃
  ┃  **The rule: a gate must not say green when it checked nothing.** This is F-98's shape
  ┃  again — the check that cannot see the thing — and it is the one the repo has paid for
  ┃  most often.
  ┃
  ┃  **The fix:** the check must see an untracked golden. The remedy is the git one — tell
  ┃  git to **intent-to-add** the golden files before the diff, so a new file appears in it
  ┃  as a change rather than as nothing. (`git add -N` is the command; the runner's call
  ┃  how to wire it.)
  ┃
  ┃  **Prove it:** write a **new** golden file that no test claims (temporarily), and
  ┃  confirm `make goldens-check` **fails** — and that the message names the file. Then
  ┃  write a new golden file that **is** claimed but is not committed, and confirm the
  ┃  check **also** fails, because a file nobody committed is not in the commit. Revert
  ┃  both.
  ┃
  ┃  **And the tree must stay clean** — the check must not leave the intent-to-add staged
  ┃  behind it.
  ┃
  ┃  ### 2. The orphan message does not say what to do
  ┃
  ┃  > the orphan message names the file and the two causes but **not the remedy** — so a
  ┃  > new golden whose test is not named `...Golden` will be reported, a true positive by
  ┃  > design that reads like a false alarm to whoever added it.
  ┃
  ┃  **Every message in this gate tells the reader what to do.** The budget's message names
  ┃  the file to change; the version message says both halves; the tool's message says where
  ┃  the binaries land. This one names a cause and stops.
  ┃
  ┃  **The fix:** the orphan message names the **remedy** — the two ways out, which are
  ┃  the two ways a golden stops being claimed:
  ┃
  ┃  1. **name the test so the filter matches it** (the filter is `-run Golden`, so the
  ┃     test's name must carry `Golden`), or
  ┃  2. **extend the filter** so it matches the test.
  ┃
  ┃  And it must say which of the two causes it is, so the reader knows which remedy
  ┃  applies. A message that leaves the reader guessing is a message that will be ignored
  ┃  the next time.
  ┃
  ┃  ### The two notes that are carried as they stand
  ┃
  ┃  | Note | Ruling |
  ┃  |---|---|
  ┃  | after a **failed** `goldens-check`, an orphan golden keeps a `2000-01-01` mtime on disk | **carry it.** A filesystem
  ┃  artifact after a *failed* run is not a git change, and the next `make goldens` rewrites it. The gate already failed — the
  ┃  chef is looking at the failure, not at a green. |
  ┃  | the orphan check's scope is `internal/ui/testdata` (the `GOLDENS` variable) | **carry it.** The scope is a variable; it
  ┃  grows with the need. Say in the script's comment that the scope is the variable, so the next person knows where to widen
  ┃  it. |
  ┃
  ┃  ## The invariants
  ┃
  ┃  - `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  ┃    budget.
  ┃  - One fact, one answer: **one place decides how a test run is invoked.**
  ┃  - `core.DomainVersion` is **7**. Do not move it.
  ┃  - The gate table, the plan format and the command count (**42**) must not move.
  ┃
  ┃  ## The scope — what not to touch
  ┃
  ┃  | | |
  ┃  |---|---|
  ┃  | **The code** | **no behaviour changes.** This round is the gate's own plumbing. |
  ┃  | **`docs/audit.md`** | the reviewer's. |
  ┃  | **`docs/07-vocabulary.md`, `work/`** | the orchestrator's. |
  ┃  | The names, the glossary, the spec | settled. Do not touch. |
  ┃
  ┃  ## The tree is shared
  ┃
  ┃  Stage **by name** — never `git add -A`, never `git add .`. Leave anything you did not
  ┃  write alone. If a file you did not write has changed, **stop and report it**.
  ┃
  ┃  ## The gate
  ┃
  ┃  ```bash
  ┃  make check && make check-race && git status --short
  ┃  ```
  ┃
  ┃  Both green, and the tree clean.
  ┃
  ┃  ## One commit
  ┃
  ┃  ```
  ┃  Per-AC evidence:
  ┃  - goldens-check sees an untracked golden, and the orphan message names the remedy; make check && make check-race exit 0
  ┃  ```
  ┃
  ┃  **One write per file, and commit as soon as the file is coherent.**
  ┃
  ┃  ## Finish with the receipt — update the one you wrote when you started
  ┃
  ┃  ```bash
  ┃  cat > work/inbox/runner-slice-h.md <<EOF
  ┃  Runner: slice H fix 4 done. Commits: <n>. The gate sees an untracked golden; the orphan message names the remedy. Notes:
  ┃  <anything the reviewer should know>.
  ┃  EOF
  ┃  git add work/inbox && git commit -q -m "inbox: runner reports slice-h fix 4"
  ┃  herdr agent prompt "w1Z:p1" "Runner: slice H fix 4 done. Commits: <n>. The gate sees an untracked golden; the orphan
  ┃  message names the remedy. Notes: <anything the reviewer should know>."
  ┃  ```
  ┃
  ┃  **Commit the receipt before you nudge.**
  ┃

     + Thought · 46.3s

     Two fixes. First the scripts — the commit check and the orphan message:

  ┃
  ┃  # Wrote scripts/check-golden-commit.sh
  ┃
  ┃    1 #!/bin/sh
  ┃    2 # The golden files must be in the commit.
  ┃    3 #
  ┃    4 # usage: check-golden-commit.sh <dir>
  ┃    5 #
  ┃    6 # `git diff` does not see an untracked file, so a new golden that nobody
  ┃    7 # committed would pass the check — a gate that says green when it checked
  ┃    8 # nothing. The files are intent-to-add before the diff, so a new file shows
  ┃    9 # up in it as a change. The index is put back as it was afterwards: the
  ┃   10 # tree stays clean.
  ┃   11 set -u
  ┃   12
  ┃   13 dir=${1:-}
  ┃   14
  ┃   15 untracked=$(git ls-files --others --exclude-standard -- "$dir")
  ┃   16
  ┃   17 for f in $untracked; do
  ┃   18   git add -N -- "$f"
  ┃   19 done
  ┃   20
  ┃   21 if git diff --exit-code HEAD -- "$dir"; then
  ┃   22   status=0
  ┃   23 else
  ┃   24   status=1
  ┃   25 fi
  ┃   26
  ┃   27 # Put the index back the way it was: the intent-to-add entries go, and the
  ┃   28 # files are untracked again.
  ┃   29 for f in $untracked; do
  ┃   30   git reset -q -- "$f"
  ┃   31 done
  ┃   32
  ┃   33 if [ "$status" -ne 0 ]; then
  ┃   34   if [ -n "$untracked" ]; then
  ┃   35     echo "these golden files are not in the commit:" >&2
  ┃   36     printf '%s\n' "$untracked" >&2
  ┃   37   fi
  ┃   38   echo "the goldens changed — read the diff before committing it." >&2
  ┃   39   exit 1
  ┃   40 fi
  ┃   41
  ┃

  ┃
  ┃  # Wrote scripts/check-golden-claims.sh
  ┃
  ┃    1 #!/bin/sh
  ┃    2 # The golden files' claims.
  ┃    3 #
  ┃    4 # usage: check-golden-claims.sh mark <dir> | report <dir> <filter>
  ┃    5 #
  ┃    6 # A golden file that no test writes is a layout nothing pins — and the check
  ┃    7 # that says "the goldens are current" would pass without touching it.
  ┃    8 #
  ┃    9 # The claim is derived, not declared: `mark` moves every golden's timestamp
  ┃   10 # into the past, the golden run rewrites the ones a test claims, and
  ┃   11 # `report` names the ones nobody wrote. No manifest of test names, and no
  ┃   12 # scan for a path in a test file — the paths live in constants, and a scan
  ┃   13 # would be a guess.
  ┃   14 #
  ┃   15 # The scope is the directory it is given: widen it where the goldens live.
  ┃   16 #
  ┃   17 # A failure is loud and names the file, the cause and the remedy. Two
  ┃   18 # causes, two remedies:
  ┃   19 #
  ┃   20 #   - no test source mentions the file — write a test that pins it, or
  ┃   21 #     delete it: a golden nothing pins is a layout nothing pins;
  ┃   22 #   - a test source mentions it, but the golden run's filter misses the
  ┃   23 #     test that writes it — name the test so the filter matches it, or
  ┃   24 #     extend the filter so it matches the test.
  ┃   25 set -u
  ┃   26
  ┃   27 # A timestamp in the past, so "not newer than this" means "not written".
  ┃   28 MARK=200001010000
  ┃   29
  ┃   30 phase=${1:-}
  ┃   31 dir=${2:-}
  ┃   32 filter=${3:-Golden}
  ┃   33
  ┃   34 case "$phase" in
  ┃   35 mark)
  ┃   36   find "$dir" -name '*.golden' -type f -exec touch -t "$MARK" {} +
  ┃   37   ;;
  ┃   38 report)
  ┃   39   ref=$(mktemp)
  ┃   40   touch -t "$MARK" "$ref"
  ┃   41   orphans=$(find "$dir" -name '*.golden' -type f ! -newer "$ref")
  ┃   42   rm -f "$ref"
  ┃   43   if [ -n "$orphans" ]; then
  ┃   44     echo "these golden files are claimed by no test:" >&2
  ┃   45     printf '%s\n' "$orphans" >&2
  ┃   46     for f in $orphans; do
  ┃   47       name=$(basename "$f")
  ┃   48       mentioned=$(find . -name '*_test.go' -not -path './.git/*' \
  ┃   49         -exec grep -lF -- "$name" {} + 2>/dev/null || true)
  ┃   50       if [ -n "$mentioned" ]; then
  ┃   51         echo "$f: a test source mentions it, but the golden run's filter (-run $filter) does not match the test that
  ┃      writes it." >&2
  ┃   52         echo "  Name the test so the filter matches it — its name must carry '$filter' — or extend the filter so it
  ┃      matches the test." >&2
  ┃   53       else
  ┃   54         echo "$f: no test source mentions it, so no test claims it." >&2
  ┃   55         echo "  Write a test that pins it, or delete the file: a golden nothing pins is a layout nothing pins." >&2
  ┃   56       fi
  ┃   57     done
  ┃   58     exit 1
  ┃   59   fi
  ┃   60   echo "  every golden file is claimed by a test"
  ┃   61   ;;
  ┃   62 *)
  ┃   63   echo "usage: check-golden-claims.sh mark <dir> | report <dir> <filter>" >&2
  ┃   64   exit 2
  ┃   65   ;;
  ┃   66 esac
  ┃   67
  ┃

     Now wire them into the Makefile:

  ┃
  ┃  ← Edit Makefile
  ┃
  ┃     98   # every other run: the zero-match refusal applies here  98   # every other run: the zero-match refusal applies here
  ┃          too — a                                                      too — a
  ┃     99   # goldens-check that regenerated nothing must fail,     99   # goldens-check that regenerated nothing must fail,
  ┃          not pass with an empty                                       not pass with an empty
  ┃    100   # diff. One recipe, one spelling: the two targets      100   # diff. One recipe, one spelling: the two targets
  ┃          differ in the diff, not                                      differ in the diff, not
  ┃    101   # in the run.                                          101   # in the run.
  ┃    102 - GOLDEN_FLAGS := -run Golden -update                    102 + GOLDEN_FILTER := Golden
  ┃                                                                 103 + GOLDEN_FLAGS := -run $(GOLDEN_FILTER) -update
  ┃    103   GOLDEN_PKG := ./internal/ui                            104   GOLDEN_PKG := ./internal/ui
  ┃    104   goldens: GOLDEN_CHECK :=                               105   goldens: GOLDEN_CHECK :=
  ┃    105   goldens-check: GOLDEN_CHECK := yes                     106   goldens-check: GOLDEN_CHECK := yes
  ┃    106                                                          107
  ┃    107   goldens goldens-check:                                 108   goldens goldens-check:
  ┃    108     $(call header,$@)                                    109     $(call header,$@)
  ┃    109     @if [ "$(GOLDEN_CHECK)" = yes ]; then scripts/check- 110     @if [ "$(GOLDEN_CHECK)" = yes ]; then scripts/check-
  ┃          golden-claims.sh mark $(GOLDENS); fi                         golden-claims.sh mark $(GOLDENS); fi
  ┃    110     GO=$(GO) scripts/go-test.sh "$(GOLDEN_FLAGS)" "$(    111     GO=$(GO) scripts/go-test.sh "$(GOLDEN_FLAGS)" "$(
  ┃          GOLDEN_PKG)"                                                 GOLDEN_PKG)"
  ┃    111 -   @if [ "$(GOLDEN_CHECK)" = yes ]; then scripts/check- 112 +   @if [ "$(GOLDEN_CHECK)" = yes ]; then scripts/check-
  ┃          golden-claims.sh report $(GOLDENS); fi                       golden-claims.sh report $(GOLDENS) "$(GOLDEN_FILTER)";
  ┃                                                                       fi
  ┃    112 -   @if [ "$(GOLDEN_CHECK)" = yes ]; then \              113 +   @if [ "$(GOLDEN_CHECK)" = yes ]; then scripts/check-
  ┃                                                                       golden-commit.sh $(GOLDENS); fi
  ┃    113 -     git diff --exit-code HEAD -- $(GOLDENS) || { \
  ┃    114 -       echo "the goldens changed — read the diff before
  ┃          committing it." >&2; exit 1; \
  ┃    115 -     }; \
  ┃    116 -   fi
  ┃    117                                                          114
  ┃    118   # The old name, kept because it was in the docs and in 115   # The old name, kept because it was in the docs and in
  ┃          a habit.                                                     a habit.
  ┃    119   regen-goldens: goldens                                 116   regen-goldens: goldens
  ┃    120                                                          117
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && chmod +x scripts/check-golden-commit.sh scripts/check-golden-claims.sh && make go
  ┃  ldens-check; echo "exit $?"; git status --short
  ┃
  ┃  (2 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  7.955s
  ┃    every golden file is claimed by a test
  ┃  exit 0
  ┃   M Makefile
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? scripts/check-golden-commit.sh
  ┃

     + Thought · 749ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && printf 'a golden nobody claims\n' > internal/ui/testdata/zz-orphan.golden && echo
  ┃   "=== (a) a new golden that no test claims:" && make goldens-check; echo "=== exit $?"; echo "=== the tree:"; git st…
  ┃
  ┃  (10 earlier lines)
  ┃  === exit 2
  ┃  === the tree:
  ┃   M Makefile
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? internal/ui/testdata/zz-orphan.golden
  ┃  ?? scripts/check-golden-commit.sh
  ┃

     + Thought · 1.4s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && rm internal/ui/testdata/zz-orphan.golden && git rm --cached -q internal/ui/testda
  ┃  ta/console-recipes.golden && echo "=== the file is on disk but not in the index:" && git status --short && echo && e…
  ┃
  ┃  (12 earlier lines)
  ┃    every golden file is claimed by a test
  ┃  === exit 0
  ┃  === the tree after the failed check:
  ┃   M Makefile
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? scripts/check-golden-commit.sh
  ┃

     + Thought · 2.8s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && echo "=== ls-files --others:"; git ls-files --others --exclude-standard -- intern
  ┃  al/ui/testdata; echo "(exit $?)"; echo "=== after add -N:"; git add -N -- internal/ui/testdata/console-recipes.golde…
  ┃
  ┃  (6 earlier lines)
  ┃  === diff HEAD:
  ┃  diff exit 0
  ┃  === diff (no HEAD):
  ┃   M Makefile
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? scripts/check-golden-commit.sh
  ┃

     + Thought · 16.1s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && echo "=== the index still has it:" && git ls-files internal
  ┃  /ui/testdata && echo "=== now remove it from the index without a staged deletion:" && git update-index --force-remov…
  ┃
  ┃  (15 earlier lines)
  ┃   M Makefile
  ┃  DA internal/ui/testdata/console-recipes.golden
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? scripts/check-golden-commit.sh
  ┃  === diff HEAD:
  ┃  diff exit
  ┃

     + Thought · 23.4s

     The staged-deletion case exposed a flaw in my restore. Let me test the exact index mechanics:

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git reset -q HEAD -- internal/ui/testdata/console-recipes.golden 2>/dev/null; git
  ┃   ls-files internal/ui/testdata; echo "=== a brand new untracked golden:"; printf 'x\n' > internal/ui/testdata/zz-new…
  ┃
  ┃  (9 earlier lines)
  ┃   M Makefile
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? internal/ui/testdata/zz-new.golden
  ┃  ?? scripts/check-golden-commit.sh
  ┃  internal/ui/testdata/console-recipes.golden
  ┃  internal/ui/testdata/console-stress.golden
  ┃

     + Thought · 736ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add -N -- internal/ui/testdata/zz-new.golden && echo "=== diff HEAD (worktree
  ┃   vs commit):"; git diff --exit-code HEAD -- internal/ui/testdata | head -6; echo "exit ${PIPESTATUS[0]}"; echo "=== …
  ┃
  ┃  (7 earlier lines)
  ┃  exit
  ┃  === diff --cached HEAD (index vs commit):
  ┃  exit
  ┃   M Makefile
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? scripts/check-golden-commit.sh
  ┃

     + Thought · 3.4s

  ┃
  ┃  # Wrote scripts/check-golden-commit.sh
  ┃
  ┃    1 #!/bin/sh
  ┃    2 # The golden files must be in the commit.
  ┃    3 #
  ┃    4 # usage: check-golden-commit.sh <dir>
  ┃    5 #
  ┃    6 # `git diff` does not see an untracked file, so a new golden that nobody
  ┃    7 # committed would pass the check — a gate that says green when it checked
  ┃    8 # nothing. The files are intent-to-add before the diff, so a new file shows
  ┃    9 # up in it as a change. The index is put back as it was afterwards: the
  ┃   10 # tree stays clean.
  ┃   11 #
  ┃   12 # Two diffs, because they see two things: the worktree against the commit
  ┃   13 # (a regenerated golden that was not committed) and the index against the
  ┃   14 # commit (a golden the commit would drop).
  ┃   15 set -u
  ┃   16
  ┃   17 dir=${1:-}
  ┃   18
  ┃   19 untracked=$(git ls-files --others --exclude-standard -- "$dir")
  ┃   20
  ┃   21 for f in $untracked; do
  ┃   22   git add -N -- "$f"
  ┃   23 done
  ┃   24
  ┃   25 status=0
  ┃   26 git diff --exit-code HEAD -- "$dir" || status=1
  ┃   27 git diff --exit-code --cached HEAD -- "$dir" || status=1
  ┃   28
  ┃   29 # Put the index back the way it was: the intent-to-add entries go, and the
  ┃   30 # files are untracked again. `git rm --cached` removes the entry without
  ┃   31 # restoring anything from HEAD, so a staged deletion stays staged.
  ┃   32 for f in $untracked; do
  ┃   33   git rm --cached -q --force -- "$f" 2>/dev/null ||
  ┃   34     git update-index --force-remove -- "$f"
  ┃   35 done
  ┃   36
  ┃   37 if [ "$status" -ne 0 ]; then
  ┃   38   if [ -n "$untracked" ]; then
  ┃   39     echo "these golden files are not in the commit:" >&2
  ┃   40     printf '%s\n' "$untracked" >&2
  ┃   41   fi
  ┃   42   echo "the goldens changed — read the diff before committing it." >&2
  ┃   43   exit 1
  ┃   44 fi
  ┃   45
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && echo "=== (b) a claimed golden that nobody committed — the file is staged as drop
  ┃  ped:" && git rm --cached -q internal/ui/testdata/console-recipes.golden && git status --short && echo && make golden…
  ┃
  ┃  (64 earlier lines)
  ┃  === the tree after the check:
  ┃   M Makefile
  ┃  D  internal/ui/testdata/console-recipes.golden
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? internal/ui/testdata/console-recipes.golden
  ┃  ?? scripts/check-golden-commit.sh
  ┃

     + Thought · 4.5s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git reset -q HEAD -- internal/ui/testdata/console-recipes.golden && git status --
  ┃  short && echo "=== restored" && cp internal/ui/golden_test.go /tmp/golden_test.go.bak && python3 - <<'PY'…
  ┃
  ┃  (13 earlier lines)
  ┃  \ No newline at end of file
  ┃  these golden files are not in the commit:
  ┃  internal/ui/testdata/console-new.golden
  ┃  the goldens changed — read the diff before committing it.
  ┃  make: *** [goldens-check] Error 1
  ┃  === exit
  ┃

     + Thought · 1.5s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp /tmp/golden_test.go.bak internal/ui/golden_test.go && rm -f internal/ui/testda
  ┃  ta/console-new.golden && git status --short && echo "=== goldens-check is green again:" && make goldens-check; echo …
  ┃
  ┃  (7 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  7.589s
  ┃    every golden file is claimed by a test
  ┃  exit 0
  ┃   M Makefile
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? scripts/check-golden-commit.sh
  ┃

     + Thought · 1.3s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
  ┃  def edit(path, old, new):…
  ┃
  ┃  (6 earlier lines)
  ┃    docs/testing.md names govulncheck 1.9.0, as pinned
  ┃   M AGENTS.md
  ┃   M Makefile
  ┃   M docs/testing.md
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? scripts/check-golden-commit.sh
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check > /tmp/c4.log 2>&1; echo "check exit: $?"; grep "^==>" /tmp/c4.log; ta
  ┃  il -3 /tmp/c4.log; git status --short
  ┃
  ┃  (13 earlier lines)
  ┃  ==> tidy-check
  ┃   M AGENTS.md
  ┃   M Makefile
  ┃   M docs/testing.md
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? scripts/check-golden-commit.sh
  ┃

     ◈ Instructions updated: core/instructions

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && { time make check-race ; } > /tmp/cr4.log 2>&1; echo "check-race exit: $?"; grep
  ┃  "^==>" /tmp/cr4.log; tail -3 /tmp/cr4.log; echo "=== git status --short:"; git status --short
  ┃
  ┃  (14 earlier lines)
  ┃  === git status --short:
  ┃   M AGENTS.md
  ┃   M Makefile
  ┃   M docs/testing.md
  ┃   M scripts/check-golden-claims.sh
  ┃  ?? scripts/check-golden-commit.sh
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add Makefile scripts/check-golden-claims.sh scripts/check-golden-commit.sh AG
  ┃  ENTS.md docs/testing.md && git commit -q -F - <<'EOF'…
  ┃
  ┃  d3c8243 fix: the gate sees a golden nobody committed, and says what to do
  ┃  (clean)
  ┃

  ┃
  ┃  # Wrote work/inbox/runner-slice-h.md
  ┃
  ┃    1 Runner: slice H fix 4 done. Commits: 1 (plus this receipt). The gate sees an untracked golden; the orphan message
  ┃      names the remedy. Notes: the commit check diffs twice — the worktree against the commit and the index against the
  ┃      commit — because they see two things; the index is put back with `git rm --cached --force`, not `git reset`, so a
  ┃      staged deletion stays staged.
  ┃    2
  ┃    3 Per-AC evidence:
  ┃    4
  ┃    5 - **the gate sees an untracked golden** — `scripts/check-golden-commit.sh` intent-to-adds the untracked goldens before
  ┃      the diff, so a new file shows up in it as a change, and puts the index back afterwards: the tree stays clean.
  ┃    6   **Proved (a):** a new golden that no test claims (`internal/ui/testdata/zz-orphan.golden`) → `make goldens-check`
  ┃      **fails** and names it: *"these golden files are claimed by no test: internal/ui/testdata/zz-orphan.golden / internal/
  ┃      ui/testdata/zz-orphan.golden: no test source mentions it, so no test claims it. / Write a test that pins it, or delete
  ┃      the file: a golden nothing pins is a layout nothing pins."* Reverted.
  ┃    7   **Proved (b):** a **new** golden that a test *claims* but that nobody committed (the golden test temporarily wrote
  ┃      `testdata/console-new.golden`) → the check **fails**: *"these golden files are not in the commit: internal/ui/
  ┃      testdata/console-new.golden / the goldens changed — read the diff before committing it."* Also proved the staged-drop
  ┃      case (the golden staged as removed) → the same failure, naming the file. Both reverted; `git status --short` is clean
  ┃      after each, and the check leaves no intent-to-add behind.
  ┃    8 - **the orphan message names the remedy** — it now says which of the two causes it is, and the way out for each:
  ┃    9   - *no test source mentions it* → "Write a test that pins it, or delete the file: a golden nothing pins is a layout
  ┃      nothing pins."
  ┃   10   - *a test source mentions it, but the golden run's filter (-run Golden) does not match the test that writes it* →
  ┃      "Name the test so the filter matches it — its name must carry 'Golden' — or extend the filter so it matches the test."
  ┃   11 - **the two carried notes** — the `2000-01-01` mtime an orphan keeps after a failed run is carried as it stands, and
  ┃      the orphan check's scope is the `GOLDENS` variable: the script's comment says so, so the next person knows where to
  ┃      widen it.
  ┃   12
  ┃   13 The gate: `make check` green (1m14s), `make check-race` green (4m47s), `git status --short` clean. Counts unchanged: 0
  ┃      / 0 / 28. `core.DomainVersion` is still 7.
  ┃   14
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/inbox && git commit -q -m "inbox: runner reports slice-h fix 4" && g
  ┃  it log --oneline -2 && git status --short && echo "(clean)"
  ┃
  ┃  f2a3ff5 inbox: runner reports slice-h fix 4
  ┃  d3c8243 fix: the gate sees a golden nobody committed, and says what to do
  ┃  (clean)
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1Z:p1" "Runner: slice H fix 4 done. Commits: 1 (plus the rec
  ┃  eipt). The gate sees an untracked golden; the orphan message names the remedy. Notes: the commit check diffs twice —…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":
  ┃  "herdr:opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"done","completion_seq":664,"cwd":"/Users/thiago/
  ┃  Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,
  ┃  "screen_detection_skipped":true,"state_change_seq":664,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722",
  ┃  "terminal_title":"OC | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":
  ┃  "w1Z"},"type":"agent_prompted"}}
  ┃

     Fix 4 is done. The receipt is committed at work/inbox/runner-slice-h.md and pane w1Z:p1 has been nudged.

     1. The gate sees an untracked golden

     scripts/check-golden-commit.sh intent-to-adds the untracked goldens before the diff, so a new file shows up in it as a
     change — and puts the index back afterwards.

     Proved (a) — a new golden no test claims → fails, naming the file and the remedy.
     Proved (b) — a new golden a test claims but nobody committed → fails: "these golden files are not in the commit: internal/
     ui/testdata/console-new.golden". Also proved the staged-drop case. Both reverted; the tree is clean after each and the
     check leaves no intent-to-add behind.

     Two details worth flagging:

     - The check diffs twice — the worktree against the commit, and the index against the commit — because they see two things (
       a regenerated golden that was not committed; a golden the commit would drop).
     - The index is restored with git rm --cached --force, not git reset: a git reset would have restored the index entry from
       HEAD and silently undone a staged deletion. My first attempt had that bug; the proof caught it.

     2. The orphan message names the remedy

     It now says which cause it is, and the way out for each:

     - no test source mentions it → "Write a test that pins it, or delete the file: a golden nothing pins is a layout nothing
       pins."
     - a test source mentions it, but the filter (-run Golden) misses the writer → "Name the test so the filter matches it — its
       name must carry 'Golden' — or extend the filter so it matches the test."

     The two carried notes

     The 2000-01-01 mtime an orphan keeps after a failed run is carried as it stands; the orphan check's scope is the GOLDENS
     variable, and the script's comment now says so.

     The gate

     make check 1m14s green · make check-race 4m47s green · git status --short clean. Counts unchanged at 0 / 0 / 28; core.
     DomainVersion still 7.

     Build · MiniMax-M3.1-Flash · 11m 3s · 115.4 tok/s

  ┃
  ┃
  ┃
  ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
  ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
  ~/Code/Pessoal/maitre:master                                                                     209.9K (21%)  ctrl+p commands

```
