# The reviewer's lane — slice H, captured 2026-10-28

Dumped from pane `w1Z:pK` with `herdr pane read --source recent-unwrapped --lines 5000`.
A fresh lane: it re-derived the facts from the code and filed F-98 across five rounds.
Its best finding was the one the whole gate rests on: **the slow gate could pass having
tested nothing**, because `make test-race` never passed the filter.

```
  ┃                                                                                  71 + | The pins | `make versions-check` — the version of the tool **the gate
  ┃                                                                                       runs**, and the version `docs/testing.md` names, against the pin the
  ┃                                                                                       `Makefile` holds |
  ┃                                                                                  72 + | The goldens' claims | `make goldens-check` — every golden file must
  ┃                                                                                       be **written by a test that claims it**; a file no test writes is a
  ┃                                                                                       layout nothing pins |
  ┃    70                                                                            73
  ┃    71   ### The tripwires                                                        74   ### The tripwires
  ┃    72                                                                            75
  ┃    73   The docs are checked against the code, so a doc cannot drift silently.   76   The docs are checked against the code, so a doc cannot drift silently.
  ┃         Each one is                                                                   Each one is
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃     95   | The one format | The retired names may live in exactly two places —    95   | The one format | The retired names may live in exactly two places —
  ┃          the detection that refuses them (`retired.go`) and the deny-list that         the detection that refuses them (`retired.go`) and the deny-list that
  ┃          refuses them as names (`model.go`) — and the exemptions are **per             refuses them as names (`model.go`) — and the exemptions are **per
  ┃          declaration**, not per file; the scan reads the **whole tree** (              declaration**, not per file; the scan reads the **whole tree** (
  ┃          `internal`, `cmd`, `scripts`) and the **code**, not the prose |               `internal`, `cmd`, `scripts`) and the **code**, not the prose |
  ┃     96   | The glossary's rule | A retired word is named **once** in the          96   | The glossary's rule | A retired word is named **once** in the
  ┃          glossary's language section, in the `_Avoid_` line of the entry that          glossary's language section, in the `_Avoid_` line of the entry that
  ┃          replaced it — never in the prose |                                            replaced it — never in the prose |
  ┃     97   | The dependency budget | The live counts of `internal/core`,            97   | The dependency budget | The live counts of `internal/core`,
  ┃          `internal/cli` and `internal/ui` against the ceilings `scripts/dep-           `internal/cli` and `internal/ui` against the ceilings `scripts/dep-
  ┃          budget.sh` holds — and **`AGENTS.md`'s table against the same                 budget.sh` holds — and **`AGENTS.md`'s table against the same
  ┃          numbers**, ceiling and current count, so the doc cannot drift |               numbers**, ceiling and current count, so the doc cannot drift |
  ┃     98   | The tool pins | The version of the tool **the gate runs** against the  98   | The tool pins | The version of the tool **the gate runs** against
  ┃          pin the `Makefile` holds — and **`docs/testing.md`'s copy against the         the pin the `Makefile` holds — and **`docs/testing.md`'s copy against
  ┃          same pin** — so a stale `staticcheck` fails the gate instead of quietly       the same pin** — so a stale `staticcheck` fails the gate instead of
  ┃          running |                                                                     quietly running |
  ┃                                                                                   99 + | The goldens' claims | Every `.golden` file under `internal/ui/
  ┃                                                                                        testdata` must be **written by a test in the golden run** — derived at
  ┃                                                                                        runtime by marking the files, running the golden tests, and naming the
  ┃                                                                                        ones nobody wrote; no manifest, no scan for a path in a test file |
  ┃     99                                                                           100
  ┃    100   ## Open findings                                                        101   ## Open findings
  ┃    101                                                                           102
  ┃    102   | # | Severity | What | Where |                                         103   | # | Severity | What | Where |
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    55   | Slice G · cycle 3 | The last fix of slice G: the exemption is the      55   | Slice G · cycle 3 | The last fix of slice G: the exemption is the
  ┃         **declaration**, not the file, and the glyph's scope is a **category**,       **declaration**, not the file, and the glyph's scope is a **category**,
  ┃         not an extension | F-97 closed, by mutation and not by the diff: the          not an extension | F-97 closed, by mutation and not by the diff: the
  ┃         reviewer's probe — a copy of the drawing pasted into `internal/brand/         reviewer's probe — a copy of the drawing pasted into `internal/brand/
  ┃         mark.go` **beside** the declaration — now **fails the wire** and takes        mark.go` **beside** the declaration — now **fails the wire** and takes
  ┃         the whole suite with it, because the wire masks the two `var` literals        the whole suite with it, because the wire masks the two `var` literals
  ┃         that hold the drawings instead of exempting the file that holds them. A       that hold the drawings instead of exempting the file that holds them. A
  ┃         copy now fails in every file the wire reads — a golden, a doc, the plan,      copy now fails in every file the wire reads — a golden, a doc, the
  ┃         `cmd/`, `images/`, a non-`.go` file — while the `.txt` sources and the        plan, `cmd/`, `images/`, a non-`.go` file — while the `.txt` sources
  ┃         declaration itself still pass. The **glyph's scope is a category**: the       and the declaration itself still pass. The **glyph's scope is a
  ┃         glyph pasted into a `.py` script now fails, and the goldens, the              category**: the glyph pasted into a `.py` script now fails, and the
  ┃         wallpaper and the brand sheet still pass; a script is code, and the rule      goldens, the wallpaper and the brand sheet still pass; a script is
  ┃         is now the thing rather than the spelling. The four exclusions stay           code, and the rule is now the thing rather than the spelling. The four
  ┃         excluded for the reason the comment gives: they are records or build          exclusions stay excluded for the reason the comment gives: they are
  ┃         output. **How it looks**: at 80×24 the tab is **byte-for-byte what it         records or build output. **How it looks**: at 80×24 the tab is **byte-
  ┃         was before the mark** — `LAST EVENT` and the event row on the screen and      for-byte what it was before the mark** — `LAST EVENT` and the event row
  ┃         **no mark**; at 80×40 the full form sits above the section with `Not          on the screen and **no mark**; at 80×40 the full form sits above the
  ┃         built yet.`, the vital signs and the last event whole; the header is one      section with `Not built yet.`, the vital signs and the last event
  ┃         line with the compact glyph; the CLI's help carries the full form above       whole; the header is one line with the compact glyph; the CLI's help
  ┃         the title line. **No blockers filed** |                                       carries the full form above the title line. **No blockers filed** |
  ┃    56   | Slice H | The quality gate: one interface (`make`), the tools the gate 56   | Slice H | The quality gate: one interface (`make`), the tools the
  ┃         needs, and CI | The gate is one command and it **does not skip a check**      gate needs, and CI | The gate is one command and it **does not skip a
  ┃         — every target exists and runs on its own, a **missing tool fails** it        check** — every target exists and runs on its own, a **missing tool
  ┃         with the way to get it, and a filter that matches **no** test **fails**       fails** it with the way to get it, and a filter that matches **no**
  ┃         — the budget counts the way `AGENTS.md` says to count, holds 0 / 0 / ≤        test **fails** — the budget counts the way `AGENTS.md` says to count,
  ┃         30, and **checks the doc's table against itself**, so the doc cannot          holds 0 / 0 / ≤ 30, and **checks the doc's table against itself**, so
  ┃         drift; the tools are pinned, each carries a **written reason**, CI runs       the doc cannot drift; the tools are pinned, each carries a **written
  ┃         the gate **through `make`**, and the pre-push hook runs the fast one —        reason**, CI runs the gate **through `make`**, and the pre-push hook
  ┃         **one medium filed** (F-98): `make test-race` drops the filter it is          runs the fast one — **one medium filed** (F-98): `make test-race` drops
  ┃         given |                                                                       the filter it is given |
  ┃    57   | Slice H · cycle 2 | The fixes to slice H: the slow gate cannot pass    57   | Slice H · cycle 2 | The fixes to slice H: the slow gate cannot pass
  ┃         having tested nothing, and the seven notes | F-98 closed — **one place        having tested nothing, and the seven notes | F-98 closed — **one place
  ┃         decides how a test run is invoked**: `test` and `test-race` are two           decides how a test run is invoked**: `test` and `test-race` are two
  ┃         targets of **one recipe** and differ in the flags they add and in             targets of **one recipe** and differ in the flags they add and in
  ┃         nothing else, so the filter always reaches the script and the zero-match      nothing else, so the filter always reaches the script and the zero-
  ┃         check always applies — `make test-race PKG=./internal/brand RUN='-run         match check always applies — `make test-race PKG=./internal/brand
  ┃         NoSuchTest'` now **fails** (*"the filter matched no test — a green that       RUN='-run NoSuchTest'` now **fails** (*"the filter matched no test — a
  ┃         proves nothing"*) and `RUN='-run TestTheDrawingsMatchTheSource'` **runs       green that proves nothing"*) and `RUN='-run
  ┃         that test**; the seven notes closed with it: the gate checks the              TestTheDrawingsMatchTheSource'` **runs that test**; the seven notes
  ┃         **version it runs** against the pin (an old `staticcheck` on the `PATH`       closed with it: the gate checks the **version it runs** against the pin
  ┃         fails it with *"run `make tools`"*), `docs/testing.md`'s versions are         (an old `staticcheck` on the `PATH` fails it with *"run `make tools`"*)
  ┃         the **checked copy** of the `Makefile`'s pins, `dep-budget.sh` honours        , `docs/testing.md`'s versions are the **checked copy** of the
  ┃         `GO`, `tidy-check` puts the tree back when tidy fails, `install-hooks`        `Makefile`'s pins, `dep-budget.sh` honours `GO`, `tidy-check` puts the
  ┃         says what it replaced, `.NOTPARALLEL:` keeps `make -j` safe, and this         tree back when tidy fails, `install-hooks` says what it replaced, `.
  ┃         ledger's discipline table names the make interface — **no blockers            NOTPARALLEL:` keeps `make -j` safe, and this ledger's discipline table
  ┃         filed** |                                                                     names the make interface — **no blockers filed** |
  ┃    58   | Slice H · cycle 3 | The last fix of slice H: the goldens are a test    58   | Slice H · cycle 3 | The last fix of slice H: the goldens are a test
  ┃         invocation, and the three one-liners | The goldens go through `scripts/       invocation, and the three one-liners | The goldens go through `scripts/
  ┃         go-test.sh` like every other run, and `goldens` and `goldens-check` are       go-test.sh` like every other run, and `goldens` and `goldens-check` are
  ┃         two targets of **one recipe** — so a `goldens-check` that regenerated         two targets of **one recipe** — so a `goldens-check` that regenerated
  ┃         nothing now **fails**, where it passed with an empty diff: with the           nothing now **fails**, where it passed with an empty diff: with the
  ┃         golden tests renamed away the target says *"the filter matched no test —      golden tests renamed away the target says *"the filter matched no test
  ┃         a green that proves nothing"*. The script passes the **package before         — a green that proves nothing"*. The script passes the **package before
  ┃         the flags**, because the flags that are not `go test`'s own (`-update`)       the flags**, because the flags that are not `go test`'s own (`-update`)
  ┃         belong to the test binary. The three one-liners closed with it: `AGENTS.      belong to the test binary. The three one-liners closed with it:
  ┃         md`'s `make check` row names `versions-check`; the version failure says       `AGENTS.md`'s `make check` row names `versions-check`; the version
  ┃         **both** halves (*"run `make tools` — it lands in `<GOPATH>/bin`"* and        failure says **both** halves (*"run `make tools` — it lands in
  ┃         *"put `<GOPATH>/bin` first on your `PATH`: an older tool earlier there        `<GOPATH>/bin`"* and *"put `<GOPATH>/bin` first on your `PATH`: an
  ┃         still fails"*, because an older tool earlier on the `PATH` still wins         older tool earlier there still fails"*, because an older tool earlier
  ┃         after `make tools`); and the doc check reads **every** mention of the         on the `PATH` still wins after `make tools`); and the doc check reads
  ┃         tool, not the first — a second mention naming another version fails it.       **every** mention of the tool, not the first — a second mention naming
  ┃         The relative-path convention stands and is the house's: the scripts are       another version fails it. The relative-path convention stands and is
  ┃         `make`'s implementation, `make` runs its recipes from the repository          the house's: the scripts are `make`'s implementation, `make` runs its
  ┃         root, and every script in `scripts/` behaves that way — **no blockers         recipes from the repository root, and every script in `scripts/`
  ┃         filed** |                                                                     behaves that way — **no blockers filed** |
  ┃                                                                                  59 + | Slice H · cycle 4 | The last round of slice H: the orphan check, and
  ┃                                                                                       the doc check scoped to the pin | The goldens' claims are **derived,
  ┃                                                                                       not declared**: `check-golden-claims.sh mark` moves every `.golden`
  ┃                                                                                       timestamp into the past, the golden run rewrites the ones a test claims
  ┃                                                                                       (the golden test writes the file **unconditionally** under `-update`),
  ┃                                                                                       and `report` names the ones nobody wrote — so a golden file that no
  ┃                                                                                       test claims is a layout nothing pins, and the check says so, **naming
  ┃                                                                                       the file** and **both** causes (*no test claims it* / *the run's filter
  ┃                                                                                       misses the test that does*). The reviewer's probes: with the snapshot
  ┃                                                                                       test renamed away `make goldens-check` **fails** and names `internal/
  ┃                                                                                       ui/testdata/console-recipes.golden`; with the stress test renamed away
  ┃                                                                                       it names `internal/ui/testdata/console-stress.golden`; and it **does
  ┃                                                                                       not cry wolf** — an up-to-date golden is rewritten, so it is not named.
  ┃                                                                                       The doc check now reads **only the pin's shape**: a line carrying the
  ┃                                                                                       release version (`staticcheck 2026.2.2`) passes, and a mention carrying
  ┃                                                                                       the pin's shape with the wrong value (`v0.8.1`) fails. The cycle
  ┃                                                                                       touched **no Go source** — **no blockers filed** |
  ┃    59   | The rounds inside a cycle | Fix → verify → back to the executor →      60   | The rounds inside a cycle | Fix → verify → back to the executor →
  ┃         repeat, until the verifier is satisfied | Three rounds in cycle 8: the        repeat, until the verifier is satisfied | Three rounds in cycle 8: the
  ┃         first fix introduced three new mediums, the second closed them |              first fix introduced three new mediums, the second closed them |
  ┃    60                                                                            61
  ┃    61   ### The discipline every round runs                                      62   ### The discipline every round runs
  ┃    62                                                                            63
  ┃

  ┃
  ┃  $ sleep 150; tail -4 /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/j4-4.log 2>/dev/null; echo "=== ci ==="; tail -3 /private/var/folders
  ┃  /5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/ci4.log 2>/dev/null
  ┃
  ┃  (2 earlier lines)
  ┃
  ┃  ==> tidy-check
  ┃  === ci ===
  ┃  ?     github.com/lthiagol/maitre/cmd/maitre-cli  [no test files]
  ┃  ok    github.com/lthiagol/maitre/internal/brand  1.502s
  ┃  ok    github.com/lthiagol/maitre/internal/cli  15.586s
  ┃

     + Thought · 1.2s

  ┃
  ┃  $ sleep 240; tail -3 /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/ci4.log 2>/dev/null; git status --short
  ┃
  ┃  vulnerabilities in modules you require, but your code doesn't appear to call
  ┃  these vulnerabilities.
  ┃  Use '-show verbose' for more details.
  ┃   M docs/audit.md
  ┃

     ↳ Shell finished · make -j4 check > /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/j4-4.log 2>&1; echo "J4=$?"; make ci > /private/var/f…

     + Thought · 649ms

  ┃
  ┃  $ grep -A2 "Symbol Results" /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/ci4.log | head -6; T=/private/var/folders/5l/jx0d52nd65v87ltmt
  ┃  rwwtf5h0000gn/T/opencode; cp AGENTS.md $T/a4.bak; cp scripts/dep-budget.sh $T/d4.bak; sed -i '' 's#| `internal/ui` | external modules | 30 | 28 |#| `…
  ┃
  ┃  (11 earlier lines)
  ┃    Change AGENTS.md's dependency table to match the gate.
  ┃
  ┃  make: *** [budget] Error 1
  ┃  EXIT=2
  ┃   M docs/audit.md
  ┃  STATUS_END
  ┃

  ┃
  ┃  $ T=/private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode; printf '\nrequire github.com/lthiagol/does-not-exist v1.0.0\n' >> go.mod; cp go.mod $
  ┃  T/broken4.mod; cp go.sum $T/broken4.sum; GOPROXY=off make tidy-check 2>&1 | tail -2; echo "EXIT=${pipestatus[1]}"; diff $T/broken4.mod go.mod && echo…
  ┃
  ┃  go mod tidy failed — go.mod and go.sum were put back.
  ┃  make: *** [tidy-check] Error 1
  ┃  EXIT=2
  ┃  go.mod RESTORED byte-for-byte
  ┃  go.sum RESTORED byte-for-byte
  ┃   M docs/audit.md
  ┃  STATUS_END
  ┃

     + Thought · 2.2s

  ┃
  ┃  $ for t in fmt-check vet staticcheck lint versions-check budget tidy-check goldens goldens-check tools vuln install-hooks; do printf "%-16s " "$t"; if make
  ┃  $t > /dev/null 2>&1; then echo "OK"; else echo "FAIL($?)"; fi; done; git status --short; echo "STATUS_END"; ls -la internal/ui/testdata/
  ┃
  ┃  (13 earlier lines)
  ┃  STATUS_END
  ┃  total 16
  ┃  drwxr-xr-x@  4 thiago  staff   128 Oct 28 13:18 .
  ┃  drwxr-xr-x@ 32 thiago  staff  1024 Oct 28 15:43 ..
  ┃  -rw-r--r--@  1 thiago  staff  2638 Oct 28 15:51 console-recipes.golden
  ┃  -rw-r--r--@  1 thiago  staff  3310 Oct 28 15:51 console-stress.golden
  ┃

     + Thought · 1.2s

  ┃
  ┃  $ make fmt 2>&1|tail -1; make test 2>&1 | tail -3; echo "TEST=${pipestatus[1]}"; make build 2>&1|tail -2; git status --short; echo "STATUS_END"; echo "=== t
  ┃  he runner's cycle-4 diff (no Go source?) ==="; git diff 858cbf1..6ddc41b --stat
  ┃
  ┃  (11 earlier lines)
  ┃   Makefile                       |   2 +
  ┃   docs/testing.md                |   8 ++-
  ┃   scripts/check-golden-claims.sh |  48 ++++++++++++++++++
  ┃   scripts/check-tool-versions.sh |   5 +-
  ┃   work/briefs/fix-h3.md          | 112 +++++++++++++++++++++++++++++++++++++++++
  ┃   6 files changed, 176 insertions(+), 3 deletions(-)
  ┃

     + Thought · 5.0s

  ┃
  ┃  $ git diff docs/audit.md
  ┃
  ┃  (22 earlier lines)
  ┃  @@ -93,6 +97,7 @@ fail.
  ┃   | The glossary's rule | A retired word is named **once** in the glossary's language section, in the `_Avoid_` line of the entry that replaced it — never in
  ┃  the prose |
  ┃   | The dependency budget | The live counts of `internal/core`, `internal/cli` and `internal/ui` against the ceilings `scripts/dep-budget.sh` holds — and
  ┃  **`AGENTS.md`'s table against the same numbers**, ceiling and current count, so the doc cannot drift |
  ┃   | The tool pins | The version of the tool **the gate runs** against the pin the `Makefile` holds — and **`docs/testing.md`'s copy against the same pin** —
  ┃  so a stale `staticcheck` fails the gate instead of quietly running |
  ┃  +| The goldens' claims | Every `.golden` file under `internal/ui/testdata` must be **written by a test in the golden run** — derived at runtime by marking
  ┃  the files, running the golden tests, and naming the ones nobody wrote; no manifest, no scan for a path in a test file |
  ┃
  ┃   ## Open findings
  ┃

     + Thought · 543ms

  ┃
  ┃  $ git add docs/audit.md && git commit -q -m "docs: slice H cycle 4 in the ledger - the goldens' claims, and the three tripwires" && git log --oneline -3
  ┃
  ┃  2791f12 docs: slice H cycle 4 in the ledger - the goldens' claims, and the three tripwires
  ┃  12d385b work: slice H's last review
  ┃  1a82163 inbox: runner reports slice-h fix 3
  ┃

  ┃
  ┃  $ mkdir -p work/inbox
  ┃  cat > work/inbox/reviewer-slice-h-cycle-4.md <<EOF…
  ┃
  ┃  f72c093 inbox: reviewer reports slice-h cycle 4
  ┃  2791f12 docs: slice H cycle 4 in the ledger - the goldens' claims, and the three tripwires
  ┃  12d385b work: slice H's last review
  ┃  STATUS_END
  ┃

  ┃
  ┃  $ herdr agent prompt "w1Z:p1" "Reviewer: slice H cycle 4 done. Verdict: ok. F-98 closed: yes. Filed 0 new blockers (none). Mutation checks: (1) TestConsoleM
  ┃  atchesTheGolden renamed away - make goldens-check failed - reverted; (2) TestTheStressPlanMatchesTheGolden renamed away - make goldens-check failed -…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":
  ┃  "ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"done","completion_seq":660,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,"foreground_cwd":"/
  ┃  Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,"screen_detection_skipped":true,"state_change_seq":660,"tab_id":"w1Z:t1","terminal_id":
  ┃  "term_65ed7354d7f7722","terminal_title":"OC | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":"w1Z"},"type":
  ┃  "agent_prompted"}}
  ┃

     The orchestrator is notified. Here is the cycle-4 review.

     Verdict: ok — F-98 closed, no new blockers

     1. The orphan check — derived, and it does not cry wolf

     ┌────────────────────────────────────────────┬─────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Probe                                      │ Result                                                                                                      │
     ├────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ TestConsoleMatchesTheGolden renamed away   │ make goldens-check fails (exit 2) and names the file: "these golden files are claimed by no test: internal/ │
     │                                            │ ui/testdata/console-recipes.golden"                                                                         │
     ├────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ TestTheStressPlanMatchesTheGolden renamed  │ names internal/ui/testdata/console-stress.golden the same way                                               │
     │ away                                       │                                                                                                             │
     ├────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ Both causes in one message                 │ "either no test claims them, or the run's filter misses the test that does" — and the second is real: make  │
     │                                            │ goldens would not regenerate it either                                                                      │
     ├────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ No wolf                                    │ the golden test writes the file unconditionally under -update (os.WriteFile before it compares), so an up-  │
     │                                            │ to-date golden is rewritten, its timestamp moves, and today's run says "every golden file is claimed by a   │
     │                                            │ test" — proved by the baseline passing with both goldens byte-identical to HEAD                             │
     ├────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ Tree clean                                 │ git status empty after the runs; the marked timestamps do not survive as a modification (git does not track │
     │                                            │ mtimes)                                                                                                     │
     └────────────────────────────────────────────┴─────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

     No manifest, no scan for a path in a test file — the paths live in constants.

     2. The doc check, scoped to the pin

     - The release shape (staticcheck 2026.2.2) → passes ✓
     - The pinned shape with the wrong value (v0.8.1) → fails: "docs/testing.md names staticcheck 0.8.1; the pin is v0.8.2" ✓
     - The installed tool's version is still checked — the old staticcheck v0.6.1 on the PATH still fails the gate with "Run 'make tools' … Put <GOPATH>/bin
       first on your PATH" ✓

     3. The gate, whole

     Every target runs alone; make check, make -j4 check and make ci are green; a tool's directory off the PATH → the gate fails with the instruction; the budget
     is 0 / 0 / 28 against 0 / 0 / 30, checked both ways; make vuln → 0 reachable; the goldens did not move; tidy-check restores go.mod/go.sum byte-for-byte.

     4. Scope and the ledger

     The cycle touched no Go source (AGENTS.md, the Makefile, docs/testing.md, two scripts) — no behaviour change. DomainVersion is 7, the command count is 42,
     and the runner did not touch docs/audit.md.

     The ledger: F-98 closed, the slice H row plus the cycle 2, 3 and 4 rows in the matrix, the discipline table now carries the three new tripwires (the budget
     · the pins · the goldens' claims), and the tripwire table gained the third row.

     LOW notes (noted, not filed — slice H closes)

     (a) After a failed goldens-check an orphan golden keeps a 2000-01-01 mtime on disk — a filesystem artifact, not a git change. (b) The orphan message names
     the file and the two causes but not the remedy — extend GOLDEN_FLAGS, or give the claiming test a name the filter matches. (c) The check's scope is
     internal/ui/testdata; a golden elsewhere would not be checked. (d) git diff --exit-code HEAD does not see untracked files, so a new golden that a test
     writes but nobody commits passes the gate (pre-existing, from slice H).

     Build · MiniMax-M3.1-Flash · 12m 4s · 102.3 tok/s

     ◈ Instructions updated: core/instructions

  ┃
  ┃  You are the **reviewer** — the same lane that reviewed slice H. The orchestrator's pane
  ┃  id is `w1Z:p1`; use it verbatim in the notify command below.
  ┃
  ┃  ## Cycle 5 — the last cycle of slice H
  ┃
  ┃  The runner fixed the last two notes. **Cycle 5 is blockers only**: file **HIGH or
  ┃  MEDIUM** blockers. No nits, no cosmetics, no style.
  ┃
  ┃  **This is the last cycle for slice H.** Anything still open after it is a note for the
  ┃  ledger — **not** another round.
  ┃
  ┃  ### 1. The gate sees a golden nobody committed
  ┃
  ┃  > `git diff --exit-code` does not see untracked files, so a new golden that a test
  ┃  > writes but nobody commits passed the gate.
  ┃
  ┃  **Verify, by command:**
  ┃
  ┃  - **a new golden no test claims** → `make goldens-check` **fails**, names the file,
  ┃    and the message names the **remedy** for that cause (name the test so the filter
  ┃    matches it, or extend the filter) — and says which cause it is;
  ┃  - **a new golden a test claims but nobody committed** → **fails**, naming the file;
  ┃  - **a golden the commit would drop** (staged for deletion) → **fails** the same way;
  ┃  - and the check must **leave the tree exactly as it was** — no intent-to-add left
  ┃    behind, and a staged deletion still staged (the runner says the index is put back
  ┃    with `git rm --cached --force` rather than `git reset`, so a staged deletion stays
  ┃    staged — **check that**).
  ┃
  ┃  ### 2. The messages say what to do
  ┃
  ┃  The rule the orchestrator has applied to this gate twice now: **every message tells
  ┃  the reader what to do.** Check the orphan message and the commit message both name a
  ┃  **remedy**, and that each names the **right** remedy for the cause it reports.
  ┃
  ┃  **And check they cannot cry wolf:** the reviewer should be able to follow each message
  ┃  and be unblocked, without reading the script.
  ┃
  ┃  ### 3. The gate, whole
  ┃
  ┃  - **Every target** runs alone: `fmt` · `fmt-check` · `vet` · `staticcheck` · `lint` ·
  ┃    `test` · `test-race` · `goldens` · `goldens-check` · `budget` · `tidy-check` ·
  ┃    `versions-check` · `check` · `check-race` · `tools` · `vuln` · `ci` · `install-hooks`.
  ┃  - `make check`, `make -j4 check` and `make ci` are green.
  ┃  - The gate **does not skip a check**: with a tool's directory off the `PATH`, it fails
  ┃    with the instruction.
  ┃  - The budget is **0 / 0 / 28** against **0 / 0 / 30**, checked **both ways**.
  ┃  - The tool versions are checked **both ways** — the doc's mention and the tool the
  ┃    gate runs.
  ┃  - The goldens are checked three ways: **unchanged**, **claimed by a test**, and **in
  ┃    the commit**.
  ┃  - `make vuln` reports **0 reachable**; `tidy-check` puts the tree back byte-for-byte
  ┃    when tidy fails.
  ┃
  ┃  ### 4. Scope, and the ledger
  ┃
  ┃  - **No behaviour changes** in the code. `core.DomainVersion` is **7**; the gate table,
  ┃    the plan format, the command count (**42**) and the names did not move.
  ┃  - **The ledger** — `docs/audit.md` is yours. **F-98 closed**, the **slice H** row, the
  ┃    cycle rows, and the discipline table's tripwires: the dependency budget, the tool
  ┃    pins, the goldens' claims — and now **the goldens' commit check**.
  ┃
  ┃  ### Verdict
  ┃
  ┃  - No HIGH or MEDIUM findings → `verdict: ok`.
  ┃  - Any HIGH or MEDIUM → `verdict: changes-needed`.
  ┃  - LOW findings: **note them in your message, do not file them.** They go in the
  ┃    ledger for a later round — **not** another cycle.
  ┃
  ┃  ### Hard rules
  ┃
  ┃  - Do **not** edit implementation code — except `docs/audit.md`, which is yours.
  ┃  - Revert every mutation you make — **including** the new goldens, the staged
  ┃    deletion and the intent-to-add. The tree must be exactly as the runner left it,
  ┃    except your ledger commit.
  ┃  - No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
  ┃
  ┃  ### Notify the orchestrator
  ┃
  ┃  ```bash
  ┃  mkdir -p work/inbox
  ┃  cat > work/inbox/reviewer-slice-h-cycle-5.md <<EOF
  ┃  Reviewer: slice H cycle 5 done. Verdict: <ok|changes-needed>. F-98 closed: <yes|no>. Filed <N> new blockers (<F-99.. or none>). Mutation checks: <what you
  ┃  mutated and reverted>. The goldens' three checks: <what each did>. The messages: <what they tell the reader to do>. LOW notes: <any>.
  ┃  EOF
  ┃  git add work/inbox && git commit -q -m "inbox: reviewer reports slice-h cycle 5"
  ┃  herdr agent prompt "w1Z:p1" "Reviewer: slice H cycle 5 done. Verdict: <ok|changes-needed>. F-98 closed: <yes|no>. Filed <N> new blockers (<F-99.. or none>).
  ┃  Mutation checks: <what you mutated and reverted>. The goldens' three checks: <what each did>. The messages: <what they tell the reader to do>. LOW notes:
  ┃  <any>."
  ┃  ```
  ┃

     + Thought · 221ms

  ┃
  ┃  $ git log --oneline -6; echo "---STATUS---"; git status --short; echo "---RUNNER DIFF---"; git diff 12d385b..HEAD --stat
  ┃
  ┃  (12 earlier lines)
  ┃   scripts/check-golden-claims.sh         |  32 +++++++---
  ┃   scripts/check-golden-commit.sh         |  44 +++++++++++++
  ┃   work/briefs/fix-h4.md                  | 110 +++++++++++++++++++++++++++++++++
  ┃   work/briefs/review-h5.md               |  86 ++++++++++++++++++++++++++
  ┃   work/inbox/reviewer-slice-h-cycle-4.md |   3 +
  ┃   work/inbox/runner-slice-h.md           |  16 ++---
  ┃   10 files changed, 294 insertions(+), 23 deletions(-)
  ┃

  ┃
  ┃  $ git show --stat --format="" d3c8243; echo "=== the runner's diff ==="; git diff 12d385b..d3c8243 -- Makefile scripts/ AGENTS.md docs/testing.md
  ┃
  ┃  (194 earlier lines)
  ┃  +  if [ -n "$untracked" ]; then
  ┃  +    echo "these golden files are not in the commit:" >&2
  ┃  +    printf '%s\n' "$untracked" >&2
  ┃  +  fi
  ┃  +  echo "the goldens changed — read the diff before committing it." >&2
  ┃  +  exit 1
  ┃  +fi
  ┃

     + Thought · 30.9s

  ┃
  ┃  $ make goldens-check 2>&1 | tail -6; echo "EXIT=${pipestatus[1]}"; git status --short; echo "STATUS_END"
  ┃
  ┃  ==> goldens-check
  ┃  GO=go scripts/go-test.sh "-run Golden -update" "./internal/ui"
  ┃  ok    github.com/lthiagol/maitre/internal/ui  7.760s
  ┃    every golden file is claimed by a test
  ┃  EXIT=0
  ┃  STATUS_END
  ┃

     + Thought · 876ms

  ┃
  ┃  $ printf 'not a golden anybody claims\n' > internal/ui/testdata/console-orphan.golden; echo "=== a new golden no test claims ==="; make goldens-check 2>&1 |
  ┃   tail -10; echo "EXIT=${pipestatus[1]}"; echo "=== the tree ==="; git status --short; echo "STATUS_END"
  ┃
  ┃  (8 earlier lines)
  ┃    Write a test that pins it, or delete the file: a golden nothing pins is a layout nothing pins.
  ┃  make: *** [goldens-check] Error 1
  ┃  EXIT=2
  ┃  === the tree ===
  ┃  ?? internal/ui/testdata/console-orphan.golden
  ┃  STATUS_END
  ┃

     + Thought · 1.6s

  ┃
  ┃  $ rm -f internal/ui/testdata/console-orphan.golden; cp internal/ui/golden_test.go /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/gt5.bak;
  ┃   cat >> internal/ui/golden_test.go <<'EOF'…
  ┃
  ┃  (11 earlier lines)
  ┃  EXIT=2
  ┃  === the tree ===
  ┃   M internal/ui/golden_test.go
  ┃  ?? internal/ui/testdata/console-new.golden
  ┃  STATUS_END
  ┃  CACHED_END
  ┃

     + Thought · 3.6s

  ┃
  ┃  $ rm -f internal/ui/testdata/console-new.golden; cp /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/gt5.bak internal/ui/golden_test.go; gi
  ┃  t status --short; echo "=== staging a deletion (the file stays on disk) ==="; git rm --cached -q -- internal/ui/testdata/console-recipes.golden; git …
  ┃
  ┃  (15 earlier lines)
  ┃  ?? internal/ui/testdata/console-recipes.golden
  ┃  total 16
  ┃  drwxr-xr-x@  4 thiago  staff   128 Oct 28 16:09 .
  ┃  drwxr-xr-x@ 32 thiago  staff  1024 Oct 28 15:43 ..
  ┃  -rw-r--r--@  1 thiago  staff  2638 Oct 28 16:09 console-recipes.golden
  ┃  -rw-r--r--@  1 thiago  staff  3310 Oct 28 16:09 console-stress.golden
  ┃

     + Thought · 3.9s

  ┃
  ┃  $ git reset -q HEAD -- internal/ui/testdata/console-recipes.golden; git status --short; echo "STATUS_END(clean expected)"; echo "=== staging a real deletion
  ┃   (the file goes too) ==="; git rm -q -- internal/ui/testdata/console-recipes.golden; git status --short; ls internal/ui/testdata/; echo "=== make gol…
  ┃
  ┃  (14 earlier lines)
  ┃  ?? internal/ui/testdata/console-recipes.golden
  ┃  total 16
  ┃  drwxr-xr-x@  4 thiago  staff   128 Oct 28 16:09 .
  ┃  drwxr-xr-x@ 32 thiago  staff  1024 Oct 28 15:43 ..
  ┃  -rw-r--r--@  1 thiago  staff  2638 Oct 28 16:09 console-recipes.golden
  ┃  -rw-r--r--@  1 thiago  staff  3310 Oct 28 16:09 console-stress.golden
  ┃

     + Thought · 3.8s

  ┃
  ┃  $ git reset -q HEAD -- internal/ui/testdata/console-recipes.golden; git status --short; echo "STATUS_END(clean expected)"; cp internal/ui/golden_test.go /pr
  ┃  ivate/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/gt5b.bak; cat >> internal/ui/golden_test.go <<'EOF'…
  ┃
  ┃  (9 earlier lines)
  ┃    Name the test so the filter matches it — its name must carry 'Golden' — or extend the filter so it matches the test.
  ┃  make: *** [goldens-check] Error 1
  ┃  EXIT=2
  ┃   M internal/ui/golden_test.go
  ┃  ?? internal/ui/testdata/console-orphan2.golden
  ┃  STATUS_END
  ┃

     + Thought · 907ms

  ┃
  ┃  $ cp /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/gt5b.bak internal/ui/golden_test.go; rm -f internal/ui/testdata/console-orphan2.golde
  ┃  n internal/ui/testdata/console-orphan.golden internal/ui/testdata/console-new.golden; git status --short; echo "STATUS_END"; git diff --stat; echo "D…
  ┃
  ┃  STATUS_END
  ┃  DIFF_END
  ┃  total 16
  ┃  drwxr-xr-x@  4 thiago  staff   128 Oct 28 16:09 .
  ┃  drwxr-xr-x@ 32 thiago  staff  1024 Oct 28 15:43 ..
  ┃  -rw-r--r--@  1 thiago  staff  2638 Oct 28 16:09 console-recipes.golden
  ┃  -rw-r--r--@  1 thiago  staff  3310 Oct 28 16:09 console-stress.golden
  ┃

     + Thought · 1.2s

  ┃
  ┃  $ T=/private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode; cp AGENTS.md $T/a5.bak; cp scripts/dep-budget.sh $T/d5.bak; cp docs/testing.md $T/t5.
  ┃  bak; sed -i '' 's#| `internal/ui` | external modules | 30 | 28 |#| `internal/ui` | external modules | 30 | 24 |#' AGENTS.md; echo "--- doctored doc c…
  ┃
  ┃  (9 earlier lines)
  ┃  EXIT=2
  ┃  --- release shape ---
  ┃    docs/testing.md names staticcheck 0.8.2, as pinned
  ┃    docs/testing.md names govulncheck 1.9.0, as pinned
  ┃  EXIT=0
  ┃  STATUS_END
  ┃

  ┃
  ┃  $ T=/private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode; printf '\nrequire github.com/lthiagol/does-not-exist v1.0.0\n' >> go.mod; cp go.mod $
  ┃  T/b5.mod; cp go.sum $T/b5.sum; GOPROXY=off make tidy-check 2>&1 | tail -2; echo "EXIT=${pipestatus[1]}"; diff $T/b5.mod go.mod && echo "go.mod RESTOR…
  ┃
  ┃  (5 earlier lines)
  ┃  staticcheck is version '0.6.1'; the pin is v0.8.2.
  ┃  Run 'make tools' to install the pinned version — it lands in /Users/thiago/go/bin.
  ┃  Put /Users/thiago/go/bin first on your PATH: an older staticcheck earlier there still fails.
  ┃  make: *** [staticcheck] Error 1
  ┃  OLD_TOOL_EXIT=2
  ┃  STATUS_END
  ┃

     + Thought · 6.2s

  ┃
  ┃  $ make check > /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/check5.log 2>&1; echo "CHECK=$?"; make -j4 check > /private/var/folders/5l/
  ┃  jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/j45.log 2>&1; echo "J4=$?"; make ci > /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/ci5…
  ┃
  ┃   Background
  ┃

  ┃
  ┃  $ D=/private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/ita; rm -rf $D; mkdir -p $D; cd $D; git init -q .; git config user.email t@t; git conf
  ┃  ig user.name t; echo tracked > a.txt; git add a.txt; git commit -qm init; echo new > b.txt; git add -N b.txt; echo "--- status with an intent-to-add …
  ┃
  ┃  (11 earlier lines)
  ┃  8b0087286157496d40a00a42e5e2ce46e0f06a37
  ┃
  ┃   a.txt | 1 +
  ┃   1 file changed, 1 insertion(+)
  ┃  --- the committed content of b.txt ---
  ┃  fatal: path 'b.txt' exists on disk, but not in 'HEAD'
  ┃

     + Thought · 1.9s

  ┃
  ┃  $ D=/private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/ita; cd $D; echo "more" >> a.txt; git add a.txt; git add -N b.txt; git status --short;
  ┃   echo "--- commit with both staged ---"; git commit -qm "second"; echo "COMMIT_EXIT=$?"; git show --stat --format="" HEAD; echo "--- is b.txt in the …
  ┃
  ┃  (2 earlier lines)
  ┃  --- commit with both staged ---
  ┃  COMMIT_EXIT=0
  ┃   a.txt | 1 +
  ┃   1 file changed, 1 insertion(+)
  ┃  --- is b.txt in the commit? ---
  ┃  fatal: path 'b.txt' exists on disk, but not in 'HEAD'
  ┃

     + Thought · 11.5s

  ┃
  ┃  $ sleep 240; tail -3 /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/check5.log; echo "=== j4 ==="; tail -2 /private/var/folders/5l/jx0d52
  ┃  nd65v87ltmtrwwtf5h0000gn/T/opencode/j45.log 2>/dev/null; echo "=== ci ==="; tail -2 /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode…
  ┃
  ┃  (3 earlier lines)
  ┃  === j4 ===
  ┃
  ┃  ==> tidy-check
  ┃  === ci ===
  ┃  ?     github.com/lthiagol/maitre/cmd/maitre-cli  [no test files]
  ┃  ok    github.com/lthiagol/maitre/internal/brand  1.590s
  ┃

     + Thought · 7.0s

  ┃
  ┃  $ sleep 200; tail -3 /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/ci5.log; git status --short
  ┃
  ┃  GO=go scripts/check-tool-versions.sh tool govulncheck v1.9.0
  ┃    govulncheck 1.9.0 matches the pin v1.9.0
  ┃  GO=go scripts/require-tool.sh govulncheck ./...
  ┃

  ┃
  ┃  $ sleep 60; tail -8 /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/ci5.log; git status --short; echo "STATUS_END"
  ┃
  ┃  No vulnerabilities found.
  ┃
  ┃  Your code is affected by 0 vulnerabilities.
  ┃  This scan also found 9 vulnerabilities in packages you import and 9
  ┃  vulnerabilities in modules you require, but your code doesn't appear to call
  ┃  these vulnerabilities.
  ┃  Use '-show verbose' for more details.
  ┃  STATUS_END
  ┃

     ↳ Shell finished · make check > /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/check5.log 2>&1; echo "CHECK=$?"; make -j4 check > /priva…

     + Thought · 513ms

  ┃
  ┃  $ cp internal/ui/golden_test.go /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/gt5c.bak; sed -i '' 's/TestConsoleMatchesTheGolden/TestCon
  ┃  soleMatchesTheSnapshot/g' internal/ui/golden_test.go; echo "=== the writer renamed away: which cause does the message report? ==="; make goldens-chec…
  ┃
  ┃  (5 earlier lines)
  ┃  internal/ui/testdata/console-recipes.golden
  ┃  internal/ui/testdata/console-recipes.golden: a test source mentions it, but the golden run's filter (-run Golden) does not match the test that writes it.
  ┃    Name the test so the filter matches it — its name must carry 'Golden' — or extend the filter so it matches the test.
  ┃  make: *** [goldens-check] Error 1
  ┃  EXIT=2
  ┃  STATUS_END
  ┃

     + Thought · 3.1s

  ┃
  ┃  $ for t in fmt-check vet staticcheck lint versions-check budget tidy-check goldens goldens-check regen-goldens tools vuln install-hooks; do printf "%-16s "
  ┃  "$t"; if make $t > /dev/null 2>&1; then echo "OK"; else echo "FAIL"; fi; done; git status --short; echo "STATUS_END"; git diff HEAD --stat -- interna…
  ┃
  ┃  (9 earlier lines)
  ┃  regen-goldens    OK
  ┃  tools            OK
  ┃  vuln             OK
  ┃  install-hooks    OK
  ┃  STATUS_END
  ┃  GOLDENS_UNCHANGED
  ┃

     + Thought · 5.5s

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    70   | The goldens | Byte-identical — except where a feature says otherwise.  70   | The goldens | Byte-identical — except where a feature says otherwise.
  ┃         The Chopboard added a tab, so the tab bar moved and nothing else did. |       The Chopboard added a tab, so the tab bar moved and nothing else did. |
  ┃    71   | The budget | `make budget` — the live counts against the ceilings the  71   | The budget | `make budget` — the live counts against the ceilings the
  ┃         script holds, **and `AGENTS.md`'s table against them**, so the doc            script holds, **and `AGENTS.md`'s table against them**, so the doc
  ┃         cannot drift from the gate |                                                  cannot drift from the gate |
  ┃    72   | The pins | `make versions-check` — the version of the tool **the gate  72   | The pins | `make versions-check` — the version of the tool **the gate
  ┃         runs**, and the version `docs/testing.md` names, against the pin the          runs**, and the version `docs/testing.md` names, against the pin the
  ┃         `Makefile` holds |                                                            `Makefile` holds |
  ┃    73   | The goldens' claims | `make goldens-check` — every golden file must be 73   | The goldens' claims | `make goldens-check` — every golden file must
  ┃         **written by a test that claims it**; a file no test writes is a layout       be **written by a test that claims it**; a file no test writes is a
  ┃         nothing pins |                                                                layout nothing pins |
  ┃                                                                                  74 + | The goldens' commit | `make goldens-check` — every golden file must
  ┃                                                                                       be **in the commit**: an untracked golden, or one the commit would
  ┃                                                                                       drop, fails the gate |
  ┃    74                                                                            75
  ┃    75   ### The tripwires                                                        76   ### The tripwires
  ┃    76                                                                            77
  ┃    77   The docs are checked against the code, so a doc cannot drift silently.   78   The docs are checked against the code, so a doc cannot drift silently.
  ┃         Each one is                                                                   Each one is
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃     98   | The glossary's rule | A retired word is named **once** in the          98   | The glossary's rule | A retired word is named **once** in the
  ┃          glossary's language section, in the `_Avoid_` line of the entry that          glossary's language section, in the `_Avoid_` line of the entry that
  ┃          replaced it — never in the prose |                                            replaced it — never in the prose |
  ┃     99   | The dependency budget | The live counts of `internal/core`,            99   | The dependency budget | The live counts of `internal/core`,
  ┃          `internal/cli` and `internal/ui` against the ceilings `scripts/dep-           `internal/cli` and `internal/ui` against the ceilings `scripts/dep-
  ┃          budget.sh` holds — and **`AGENTS.md`'s table against the same                 budget.sh` holds — and **`AGENTS.md`'s table against the same
  ┃          numbers**, ceiling and current count, so the doc cannot drift |               numbers**, ceiling and current count, so the doc cannot drift |
  ┃    100   | The tool pins | The version of the tool **the gate runs** against the 100   | The tool pins | The version of the tool **the gate runs** against
  ┃          pin the `Makefile` holds — and **`docs/testing.md`'s copy against the         the pin the `Makefile` holds — and **`docs/testing.md`'s copy against
  ┃          same pin** — so a stale `staticcheck` fails the gate instead of quietly       the same pin** — so a stale `staticcheck` fails the gate instead of
  ┃          running |                                                                     quietly running |
  ┃    101   | The goldens' claims | Every `.golden` file under `internal/ui/        101   | The goldens' claims | Every `.golden` file under `internal/ui/
  ┃          testdata` must be **written by a test in the golden run** — derived at        testdata` must be **written by a test in the golden run** — derived at
  ┃          runtime by marking the files, running the golden tests, and naming the        runtime by marking the files, running the golden tests, and naming the
  ┃          ones nobody wrote; no manifest, no scan for a path in a test file |           ones nobody wrote; no manifest, no scan for a path in a test file |
  ┃                                                                                  102 + | The goldens' commit | Every `.golden` under `internal/ui/testdata`
  ┃                                                                                        must be **in the commit** — the files are intent-to-add before the
  ┃                                                                                        diff, so a new one shows up in it as a change; the index is put back
  ┃                                                                                        as it was, and a staged deletion stays staged |
  ┃    102                                                                           103
  ┃    103   ## Open findings                                                        104   ## Open findings
  ┃    104                                                                           105
  ┃    105   | # | Severity | What | Where |                                         106   | # | Severity | What | Where |
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    56   | Slice H | The quality gate: one interface (`make`), the tools the gate 56   | Slice H | The quality gate: one interface (`make`), the tools the
  ┃         needs, and CI | The gate is one command and it **does not skip a check**      gate needs, and CI | The gate is one command and it **does not skip a
  ┃         — every target exists and runs on its own, a **missing tool fails** it        check** — every target exists and runs on its own, a **missing tool
  ┃         with the way to get it, and a filter that matches **no** test **fails**       fails** it with the way to get it, and a filter that matches **no**
  ┃         — the budget counts the way `AGENTS.md` says to count, holds 0 / 0 / ≤        test **fails** — the budget counts the way `AGENTS.md` says to count,
  ┃         30, and **checks the doc's table against itself**, so the doc cannot          holds 0 / 0 / ≤ 30, and **checks the doc's table against itself**, so
  ┃         drift; the tools are pinned, each carries a **written reason**, CI runs       the doc cannot drift; the tools are pinned, each carries a **written
  ┃         the gate **through `make`**, and the pre-push hook runs the fast one —        reason**, CI runs the gate **through `make`**, and the pre-push hook
  ┃         **one medium filed** (F-98): `make test-race` drops the filter it is          runs the fast one — **one medium filed** (F-98): `make test-race` drops
  ┃         given |                                                                       the filter it is given |
  ┃    57   | Slice H · cycle 2 | The fixes to slice H: the slow gate cannot pass    57   | Slice H · cycle 2 | The fixes to slice H: the slow gate cannot pass
  ┃         having tested nothing, and the seven notes | F-98 closed — **one place        having tested nothing, and the seven notes | F-98 closed — **one place
  ┃         decides how a test run is invoked**: `test` and `test-race` are two           decides how a test run is invoked**: `test` and `test-race` are two
  ┃         targets of **one recipe** and differ in the flags they add and in             targets of **one recipe** and differ in the flags they add and in
  ┃         nothing else, so the filter always reaches the script and the zero-match      nothing else, so the filter always reaches the script and the zero-
  ┃         check always applies — `make test-race PKG=./internal/brand RUN='-run         match check always applies — `make test-race PKG=./internal/brand
  ┃         NoSuchTest'` now **fails** (*"the filter matched no test — a green that       RUN='-run NoSuchTest'` now **fails** (*"the filter matched no test — a
  ┃         proves nothing"*) and `RUN='-run TestTheDrawingsMatchTheSource'` **runs       green that proves nothing"*) and `RUN='-run
  ┃         that test**; the seven notes closed with it: the gate checks the              TestTheDrawingsMatchTheSource'` **runs that test**; the seven notes
  ┃         **version it runs** against the pin (an old `staticcheck` on the `PATH`       closed with it: the gate checks the **version it runs** against the pin
  ┃         fails it with *"run `make tools`"*), `docs/testing.md`'s versions are         (an old `staticcheck` on the `PATH` fails it with *"run `make tools`"*)
  ┃         the **checked copy** of the `Makefile`'s pins, `dep-budget.sh` honours        , `docs/testing.md`'s versions are the **checked copy** of the
  ┃         `GO`, `tidy-check` puts the tree back when tidy fails, `install-hooks`        `Makefile`'s pins, `dep-budget.sh` honours `GO`, `tidy-check` puts the
  ┃         says what it replaced, `.NOTPARALLEL:` keeps `make -j` safe, and this         tree back when tidy fails, `install-hooks` says what it replaced, `.
  ┃         ledger's discipline table names the make interface — **no blockers            NOTPARALLEL:` keeps `make -j` safe, and this ledger's discipline table
  ┃         filed** |                                                                     names the make interface — **no blockers filed** |
  ┃    58   | Slice H · cycle 3 | The last fix of slice H: the goldens are a test    58   | Slice H · cycle 3 | The last fix of slice H: the goldens are a test
  ┃         invocation, and the three one-liners | The goldens go through `scripts/       invocation, and the three one-liners | The goldens go through `scripts/
  ┃         go-test.sh` like every other run, and `goldens` and `goldens-check` are       go-test.sh` like every other run, and `goldens` and `goldens-check` are
  ┃         two targets of **one recipe** — so a `goldens-check` that regenerated         two targets of **one recipe** — so a `goldens-check` that regenerated
  ┃         nothing now **fails**, where it passed with an empty diff: with the           nothing now **fails**, where it passed with an empty diff: with the
  ┃         golden tests renamed away the target says *"the filter matched no test —      golden tests renamed away the target says *"the filter matched no test
  ┃         a green that proves nothing"*. The script passes the **package before         — a green that proves nothing"*. The script passes the **package before
  ┃         the flags**, because the flags that are not `go test`'s own (`-update`)       the flags**, because the flags that are not `go test`'s own (`-update`)
  ┃         belong to the test binary. The three one-liners closed with it: `AGENTS.      belong to the test binary. The three one-liners closed with it:
  ┃         md`'s `make check` row names `versions-check`; the version failure says       `AGENTS.md`'s `make check` row names `versions-check`; the version
  ┃         **both** halves (*"run `make tools` — it lands in `<GOPATH>/bin`"* and        failure says **both** halves (*"run `make tools` — it lands in
  ┃         *"put `<GOPATH>/bin` first on your `PATH`: an older tool earlier there        `<GOPATH>/bin`"* and *"put `<GOPATH>/bin` first on your `PATH`: an
  ┃         still fails"*, because an older tool earlier on the `PATH` still wins         older tool earlier there still fails"*, because an older tool earlier
  ┃         after `make tools`); and the doc check reads **every** mention of the         on the `PATH` still wins after `make tools`); and the doc check reads
  ┃         tool, not the first — a second mention naming another version fails it.       **every** mention of the tool, not the first — a second mention naming
  ┃         The relative-path convention stands and is the house's: the scripts are       another version fails it. The relative-path convention stands and is
  ┃         `make`'s implementation, `make` runs its recipes from the repository          the house's: the scripts are `make`'s implementation, `make` runs its
  ┃         root, and every script in `scripts/` behaves that way — **no blockers         recipes from the repository root, and every script in `scripts/`
  ┃         filed** |                                                                     behaves that way — **no blockers filed** |
  ┃    59   | Slice H · cycle 4 | The last round of slice H: the orphan check, and   59   | Slice H · cycle 4 | The last round of slice H: the orphan check, and
  ┃         the doc check scoped to the pin | The goldens' claims are **derived, not      the doc check scoped to the pin | The goldens' claims are **derived,
  ┃         declared**: `check-golden-claims.sh mark` moves every `.golden`               not declared**: `check-golden-claims.sh mark` moves every `.golden`
  ┃         timestamp into the past, the golden run rewrites the ones a test claims       timestamp into the past, the golden run rewrites the ones a test claims
  ┃         (the golden test writes the file **unconditionally** under `-update`),        (the golden test writes the file **unconditionally** under `-update`),
  ┃         and `report` names the ones nobody wrote — so a golden file that no test      and `report` names the ones nobody wrote — so a golden file that no
  ┃         claims is a layout nothing pins, and the check says so, **naming the          test claims is a layout nothing pins, and the check says so, **naming
  ┃         file** and **both** causes (*no test claims it* / *the run's filter           the file** and **both** causes (*no test claims it* / *the run's filter
  ┃         misses the test that does*). The reviewer's probes: with the snapshot         misses the test that does*). The reviewer's probes: with the snapshot
  ┃         test renamed away `make goldens-check` **fails** and names `internal/ui/      test renamed away `make goldens-check` **fails** and names `internal/
  ┃         testdata/console-recipes.golden`; with the stress test renamed away it        ui/testdata/console-recipes.golden`; with the stress test renamed away
  ┃         names `internal/ui/testdata/console-stress.golden`; and it **does not         it names `internal/ui/testdata/console-stress.golden`; and it **does
  ┃         cry wolf** — an up-to-date golden is rewritten, so it is not named. The       not cry wolf** — an up-to-date golden is rewritten, so it is not named.
  ┃         doc check now reads **only the pin's shape**: a line carrying the             The doc check now reads **only the pin's shape**: a line carrying the
  ┃         release version (`staticcheck 2026.2.2`) passes, and a mention carrying       release version (`staticcheck 2026.2.2`) passes, and a mention carrying
  ┃         the pin's shape with the wrong value (`v0.8.1`) fails. The cycle touched      the pin's shape with the wrong value (`v0.8.1`) fails. The cycle
  ┃         **no Go source** — **no blockers filed** |                                    touched **no Go source** — **no blockers filed** |
  ┃                                                                                  60 + | Slice H · cycle 5 | The last cycle of slice H: the gate sees a golden
  ┃                                                                                       nobody committed, and every message says what to do | The goldens are
  ┃                                                                                       checked **three ways** — *unchanged*, *claimed by a test*, *in the
  ┃                                                                                       commit*. The commit check is the new one: `git diff` does not see an
  ┃                                                                                       untracked file, so the files are **intent-to-add** before the diff and
  ┃                                                                                       a new golden shows up in it as a change; the index is put back as it
  ┃                                                                                       was (`git rm --cached --force`, not `git reset`, so **a staged deletion
  ┃                                                                                       stays staged**). The reviewer's probes: a new golden **no test claims**
  ┃                                                                                       → the check fails, names the file and says which cause (*"no test
  ┃                                                                                       source mentions it, so no test claims it"*) and its remedy (*"write a
  ┃                                                                                       test that pins it, or delete the file"*); a new golden **a test claims
  ┃                                                                                       but nobody committed** → the check fails and names it (*"these golden
  ┃                                                                                       files are not in the commit"*); a golden **staged for deletion** → the
  ┃                                                                                       same failure, and the index afterwards is the index before (*"D "* plus
  ┃                                                                                       the file on disk, untracked). The messages were brought to the rule:
  ┃                                                                                       each cause names **its own** remedy — *"name the test so the filter
  ┃                                                                                       matches it — its name must carry 'Golden' — or extend the filter"*, and
  ┃                                                                                       *"write a test that pins it, or delete the file"* — and the reviewer
  ┃                                                                                       could follow each one and be unblocked without reading the script. The
  ┃                                                                                       doc check is scoped to the **pin's shape**: the release version (
  ┃                                                                                       `staticcheck 2026.2.2`) passes and a mention carrying the pin's shape
  ┃                                                                                       with the wrong value (`v0.8.1`) fails. The cycle touched **no Go
  ┃                                                                                       source** — **no blockers filed** |
  ┃    60   | The rounds inside a cycle | Fix → verify → back to the executor →      61   | The rounds inside a cycle | Fix → verify → back to the executor →
  ┃         repeat, until the verifier is satisfied | Three rounds in cycle 8: the        repeat, until the verifier is satisfied | Three rounds in cycle 8: the
  ┃         first fix introduced three new mediums, the second closed them |              first fix introduced three new mediums, the second closed them |
  ┃    61                                                                            62
  ┃    62   ### The discipline every round runs                                      63   ### The discipline every round runs
  ┃    63                                                                            64
  ┃

     + Thought · 755ms

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    111   | **F-61** | medium | **The `courses/` directory is named by the plan   111   | **F-61** | medium | **The `courses/` directory is named by the plan
  ┃          format, created by `Init`, and never written** — a directory no writer        format, created by `Init`, and never written** — a directory no writer
  ┃          touches. Per [06-proposals.md](06-proposals.md) § 1 the honest answer         touches. Per [06-proposals.md](06-proposals.md) § 1 the honest answer
  ┃          is to stop naming it: a course is the relation, not a file. |                 is to stop naming it: a course is the relation, not a file. |
  ┃          `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.md) |            `internal/core/store.go:18,63`; [docs/03 §6](03-domain-model.md) |
  ┃    112   | F-56 | low | **The poll's fingerprint is blind to a write that keeps  112   | F-56 | low | **The poll's fingerprint is blind to a write that keeps
  ┃          the mtime and the size.** It is path + mtime + size per file, so two          the mtime and the size.** It is path + mtime + size per file, so two
  ┃          different contents can share one fingerprint. The tool's own writes           different contents can share one fingerprint. The tool's own writes
  ┃          always change the mtime, so in practice only a hand edit that restores        always change the mtime, so in practice only a hand edit that restores
  ┃          the mtime is missed — and the chef's `r` reloads it. | `internal/core/        the mtime is missed — and the chef's `r` reloads it. | `internal/core/
  ┃          store.go` — `DiskFingerprint`. |                                              store.go` — `DiskFingerprint`. |
  ┃    113   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the same 113   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the
  ┃          as `--in-scope=`. | Recorded so nobody calls it a regression later. Not       same as `--in-scope=`. | Recorded so nobody calls it a regression
  ┃          a defect. |                                                                   later. Not a defect. |
  ┃    114                                                                           114
  ┃                                                                                  115 + ### Notes left for a later round
  ┃                                                                                  116 +
  ┃                                                                                  117 + Low notes, **not findings**: each one is a thing the gate does not
  ┃                                                                                        catch, and none of them was worth another cycle. They are here so a
  ┃                                                                                        later round does not rediscover them. The notes that cycles 3 to 5
  ┃                                                                                        closed are gone from this list.
  ┃                                                                                  118 +
  ┃                                                                                  119 + | | |
  ┃                                                                                  120 + |---|---|
  ┃                                                                                  121 + | The claim's cause is told by a **name scan** | The orphan report
  ┃                                                                                        tells the two causes apart by searching the test sources for the
  ┃                                                                                        file's **name** — so a mention in a constant or a comment counts as
  ┃                                                                                        *"a test source mentions it"*, even when no test writes it any more.
  ┃                                                                                        The gate still fails and the remedy still unblocks; the **cause line**
  ┃                                                                                        can be wrong. |
  ┃                                                                                  122 + | The commit check touches the **shared index** | It adds the
  ┃                                                                                        untracked goldens intent-to-add and takes them out again. A commit in
  ┃                                                                                        another lane in that window is **interrupted** — proved in a scratch
  ┃                                                                                        repo: `git commit` refuses when an intent-to-add is the only change,
  ┃                                                                                        and skips it otherwise, so a commit cannot be corrupted, only delayed.
  ┃                                                                                        |
  ┃                                                                                  123 + | A failed claim leaves a **timestamp in the past** | An orphan golden
  ┃                                                                                        keeps its 2000-01-01 mtime on disk after a failed `goldens-check`. Git
  ┃                                                                                        does not track it; the next `make goldens` rewrites the file. |
  ┃                                                                                  124 + | The goldens' checks read **one directory** | `internal/ui/testdata`
  ┃                                                                                        — a golden file elsewhere in the tree is neither claimed nor checked
  ┃                                                                                        for the commit. |
  ┃                                                                                  125 + | `go-test.sh` hands an **unknown flag to every binary** | The package
  ┃                                                                                        goes first so that the goldens' `-update` reaches the test binary —
  ┃                                                                                        but a flag `go test` does not know is now passed to every test binary
  ┃                                                                                        in a run, not only to the goldens'. |
  ┃                                                                                  126 + | A golden whose name carries a **space** | The scripts walk the
  ┃                                                                                        goldens by word, so a name with a space in it would be mishandled.
  ┃                                                                                        There is none today. |
  ┃                                                                                  127 +
  ┃    115   ## Closed findings                                                      128   ## Closed findings
  ┃    116                                                                           129
  ┃    117   | # | Severity | Closed in | Held by |                                  130   | # | Severity | Closed in | Held by |
  ┃    118   |---|---|---|---|                                                       131   |---|---|---|---|
  ┃

  ┃
  ┃  $ git diff docs/audit.md | head -60; echo "..."; grep -n "Notes left for a later round" docs/audit.md
  ┃
  ┃  (43 earlier lines)
  ┃  +| A golden whose name carries a **space** | The scripts walk the goldens by word, so a name with a space in it would be mishandled. There is none today. |
  ┃  +
  ┃   ## Closed findings
  ┃
  ┃   | # | Severity | Closed in | Held by |
  ┃  ...
  ┃  115:### Notes left for a later round
  ┃

     + Thought · 967ms

  ┃
  ┃  $ git add docs/audit.md && git commit -q -m "docs: slice H cycle 5 in the ledger - the goldens' commit, and the notes for a later round" && git log --onelin
  ┃  e -3; git status --short; echo "STATUS_END"
  ┃
  ┃  a5cbd58 docs: slice H cycle 5 in the ledger - the goldens' commit, and the notes for a later round
  ┃  8a5a5f4 work: slice H's last review
  ┃  f2a3ff5 inbox: runner reports slice-h fix 4
  ┃  STATUS_END
  ┃

  ┃
  ┃  $ mkdir -p work/inbox
  ┃  cat > work/inbox/reviewer-slice-h-cycle-5.md <<EOF…
  ┃
  ┃  b57e33c inbox: reviewer reports slice-h cycle 5
  ┃  a5cbd58 docs: slice H cycle 5 in the ledger - the goldens' commit, and the notes for a later round
  ┃  8a5a5f4 work: slice H's last review
  ┃  STATUS_END
  ┃

  ┃
  ┃  $ herdr agent prompt "w1Z:p1" "Reviewer: slice H cycle 5 done. Verdict: ok. F-98 closed: yes. Filed 0 new blockers (none). Mutation checks: (1) a new golden
  ┃   no test claims - make goldens-check failed - removed; (2) a new golden a test claims but nobody committed - make goldens-check failed - restored; (3…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":
  ┃  "ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"done","completion_seq":668,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,"foreground_cwd":"/
  ┃  Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,"screen_detection_skipped":true,"state_change_seq":668,"tab_id":"w1Z:t1","terminal_id":
  ┃  "term_65ed7354d7f7722","terminal_title":"OC | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":"w1Z"},"type":
  ┃  "agent_prompted"}}
  ┃

     The orchestrator is notified. Here is the cycle-5 review.

     Verdict: ok — F-98 closed, no new blockers

     1. The gate sees a golden nobody committed

     ┌─────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────┐
     │ Probe                                               │ Result                                                                                             │
     ├─────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ A new golden no test claims                         │ make goldens-check fails (exit 2), names internal/ui/testdata/console-orphan.golden, says the      │
     │                                                     │ cause ("no test source mentions it, so no test claims it") and its remedy ("Write a test that pins │
     │                                                     │ it, or delete the file")                                                                           │
     ├─────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ A new golden a test claims but nobody committed     │ fails (exit 2), names it ("these golden files are not in the commit") — and the tree afterwards is │
     │                                                     │ exactly as it was: the file untracked, git diff --cached empty, no intent-to-add left behind       │
     ├─────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────┤
     │ A golden staged for deletion (file on disk, and     │ fails the same way, and the index afterwards is the index before — the staged deletion stays       │
     │ file gone)                                          │ staged (D  + the file on disk, untracked). The runner's claim about git rm --cached --force rather │
     │                                                     │ than git reset is true                                                                             │
     └─────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────┘

     2. The messages say what to do

     Each cause names its own remedy — the filter cause: "name the test so the filter matches it — its name must carry 'Golden' — or extend the filter"; the
     unclaimed cause: "write a test that pins it, or delete the file"; the commit cause: "these golden files are not in the commit" + "read the diff before
     committing it". I followed each one without reading the script and each unblocked me.

     3. The gate, whole

     Every target runs alone; make check, make -j4 check and make ci are green; a tool's directory off the PATH → the gate fails with the instruction; the budget
     is 0 / 0 / 28 against 0 / 0 / 30, checked both ways; the tool versions are checked both ways; the goldens are checked three ways (unchanged · claimed · in
     the commit); make vuln → 0 reachable; tidy-check restores go.mod/go.sum byte-for-byte.

     4. Scope and the ledger

     The cycle touched no Go source — no behaviour change. DomainVersion is 7, the command count is 42, and the runner did not touch docs/audit.md.

     The ledger: F-98 closed; the slice H row and the cycle 2–5 rows in the matrix; the discipline table now carries the four tripwires the gate holds (the
     budget · the pins · the goldens' claims · the goldens' commit); the tripwire table gained the commit row; and the low notes are recorded under "Notes left
     for a later round" — notes, not findings, with no F-number.

     LOW notes (in the ledger, not filed)

     (a) The claim's cause is told by a name scan, so a mention in a constant counts as "a test source mentions it" even when no test writes it any more — the
     cause line can be wrong, though the remedy still unblocks. (b) The commit check touches the shared index; proved in a scratch repo that a concurrent commit
     is interrupted, never corrupted. (c) A failed claim leaves a 2000-01-01 timestamp on the orphan golden (git does not track it). (d) The goldens' checks read
     one directory. (e) go-test.sh hands an unknown flag to every test binary in a run. (f) A golden whose name carries a space would be mishandled.

     Build · MiniMax-M3.1-Flash · 13m 20s · 107.2 tok/s

  ┃
  ┃
  ┃
  ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
  ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
  ~/Code/Pessoal/maitre:master                                                                                                      225.0K (22%)  ctrl+p commands

```
