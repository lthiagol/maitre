# The reviewer's lane — slice E, captured 2026-10-28

Dumped from pane `w1Z:pB` with `herdr pane read --source recent-unwrapped --lines 5000`.
A fresh lane: it re-derived the facts from the code, filed F-81..F-86, and proved
every fix by command. Its last finding was the best of the rename: **the tree was
clean and the instrument was blind** — the sweep's command could not see one of the
seven retired names.

```
  ┃          the guard still refuses a plan stamped 7 and accepts          the guard still refuses a plan stamped 7 and accepts
  ┃          one stamped 5. | `internal/core/version_test.go` — the        one stamped 5. | `internal/core/version_test.go` — the
  ┃          shape test (`recordedDomainShape = "b84b230161f4ab03"`)       shape test (`recordedDomainShape =
  ┃          , `TestDomainVersionRefusesANewerPlan`,                       "b84b230161f4ab03"`),
  ┃          `TestDomainVersionStampsAnOlderPlan`; `internal/core/         `TestDomainVersionRefusesANewerPlan`,
  ┃          lifecycle.go:216` |                                           `TestDomainVersionStampsAnOlderPlan`; `internal/core/
  ┃                                                                        lifecycle.go:216` |
  ┃
  ┃
  ┃
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && awk 'NR>=80 && NR<=150 {n=gsub(/\|/,"|"); print NR": "n}' docs/audit.md | awk -F':
  ┃   ' '{print $2}' | sort | uniq -c; echo "=== any row with a different count? ==="; awk 'NR>=80 && NR<=150 && /^\|/ {n=…
  ┃
  ┃  11 0
  ┃    60 5
  ┃  === any row with a different count? ===
  ┃

     + Thought · 898ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && awk 'NR>=25 && NR<=50 && /^\|/ {n=gsub(/\|/,"|"); if (n!=4) print NR": "n" -> "sub
  ┃  str($0,1,60)}' docs/audit.md; echo "(no output = the round table is well-formed)"; git diff --stat
  ┃
  ┃  (no output = the round table is well-formed)
  ┃   docs/audit.md | 12 +++++++-----
  ┃   1 file changed, 7 insertions(+), 5 deletions(-)
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/audit.md && git commit -q -m "docs: the ledger closes slice E's five
  ┃  findings, and files one against the sweep" && git log --oneline -2
  ┃
  ┃  be543d0 docs: the ledger closes slice E's five findings, and files one against the sweep
  ┃  b70e859 work: slice E's cycle 2 — the last review
  ┃

     + Thought · 2.0s

  ┃
  ┃  # Wrote work/inbox/reviewer-slice-e-cycle-2.md
  ┃
  ┃    1 Reviewer: slice E cycle 2 done. Verdict: changes-needed. F-81..F-85 closed: yes. Filed 1 new blocker (F-86). Sweep
  ┃      count: 77 lines in 19 files — it reconciles (78 → 77 is F-83's one prose line removed, nothing added), but the command
  ┃      cannot see *maître d'*. Mascot: no. LOW notes: the console's entry now sits under "### The strings" — it describes the
  ┃      console, not a string, and the section says what it carries, but the heading alone can be misread; `docs/03` §6's
  ┃      `#the-version-guard` anchor has no heading (pre-existing since `813956f`) — a broken pointer, one line to fix.
  ┃    2
  ┃    3 ## The five closures, each by a command
  ┃    4
  ┃    5 | Finding | The command | What it returned |
  ┃    6 |---|---|---|
  ┃    7 | **F-81** | `grep -n -i -E "line cook\|taster\|expo\|proof\|special\|prep note\|maître d'" GLOSSARY.md` | 7 hits, 7
  ┃      retired names, one each; **0** outside an `_Avoid_:` line (`GLOSSARY.md:161` carries *maître d'*) |
  ┃    8 | **F-82** | the crew section | "### The crew" holds **five** entries (Chef · Sous-chef · Runner · Reviewer ·
  ┃      Orchestrator); the console moved to "### The strings" (`:151`), headed `maitre` the way the names table spells it (`:
  ┃      32`), and the section says the console is a **surface** and names what it carries — the window title, the tab bar, the
  ┃      CLI's help |
  ┃    9 | **F-83** | the same grep, filtered | the prose no longer names *special* or *prep note*; it points at ADR 0009 § 3,
  ┃      and § 3 does say the thing: "The backlog had **two file types** — a *special* (a hunch) and a *prep note*" |
  ┃   10 | **F-84** | the three lists, against `internal/core/format.go:20` | `docs/03:451`, `docs/03:552` and `internal/core/
  ┃      lifecycle.go:200` all read "the id prefixes, the directories, the root files and the field names on disk" — and that is
  ┃      what `planFormatFacts` hashes |
  ┃   11 | **F-85** | `go test ./internal/core -run 'TestDomainVersion\|TestDomainShape'` | `core.DomainVersion` is still **6**;
  ┃      a plan stamped **7** is refused, a plan stamped **5** is accepted and restamped; the prose now says the version fences
  ┃      the writer and the fingerprint is the tripwire, in `docs/03` §6/§6.1, the `DomainVersion` comment and `docs/07-
  ┃      vocabulary.md` — and it does not overclaim the other way |
  ┃   12
  ┃   13 ## The sweep's count
  ┃   14
  ┃   15 The published command, run at the cycle-1 base (`72783a7`) and at HEAD, on the same tree:
  ┃   16
  ┃   17 | | Lines |
  ┃   18 |---|---|
  ┃   19 | base | 78 |
  ┃   20 | HEAD | 77 |
  ┃   21 | delta | exactly one line removed — `GLOSSARY.md:59`, the prose F-83 took the retired words out of |
  ┃   22
  ┃   23 Nothing was added: the console's `_Avoid_` line already existed, and the name F-81 put on it is one the pattern cannot
  ┃      match. **Every one of the 77 lines names something that still exists** — the legacy reader, the retired-role aliases,
  ┃      the old-id and migration tests, the dated version notes, the tripwires' own mutations, the glossary's `_Avoid_` lines,
  ┃      `prototypes/README.md`.
  ┃   24
  ┃   25 Extending the pattern with *maître d'* returns **78** — the extra line is exactly the glossary's own `_Avoid_` line.
  ┃      That is F-86: the command has to learn the difference between the glossary's warning and the drift.
  ┃   26
  ┃   27 ## F-86 (medium) — the sweep cannot see the console's retired name
  ┃   28
  ┃   29 AC-04's list of old names (`work/briefs/slice-e.md:81`) and the command that proves it (commit `c1c1e7c`) carry
  ┃      `Kitchen Brigade`, `scale`, `line-cook`, `taster`, `expo`, `Chopboard`, `special`, `prep note` and `Proof` — **not**
  ┃      *maître d'*. One of the seven renamed terms is outside the instrument's reach, and the count the sweep publishes hides
  ┃      it.
  ┃   30
  ┃   31 The tree is clean: the only other hits are the dated records (ADR 0009 § 2 and § 4, ADR 0004, ADR 0005). But the proof
  ┃      does not cover the term, so a *maître d'* that came back as a current name would not be caught. The remedy is the
  ┃      command's, not the file's: the pattern must carry the name, and the command must then exclude the glossary's `_Avoid_`
  ┃      lines.
  ┃   32
  ┃   33 ## The mascot and the tree
  ┃   34
  ┃   35 The mascot **did not land** — `logo/` is untracked and was never added (`git log --all --diff-filter=A -- logo/` is
  ┃      empty); no placement point carries a drawing; AC-01 moves to the drawing session. This slice committed **nothing** of
  ┃      the other session's work: `git status --short` shows only `?? logo/`. The goldens did not move (`internal/ui/testdata/
  ┃      *.golden` untouched) and `make regen-goldens` leaves the tree clean.
  ┃   36
  ┃   37 ## Scope
  ┃   38
  ┃   39 The five fix commits touched `GLOSSARY.md`, `docs/03-domain-model.md`, `docs/07-vocabulary.md`, `README.md`, `work/
  ┃      README.md` and `internal/core/lifecycle.go` — and in `lifecycle.go` the diff is **comments only** (verified: no non-
  ┃      comment line in the diff). The gate table, the plan format and the transition table did not move, `core.DomainVersion`
  ┃      is still **6**, the command count is still **42**, and `kitchen/activity.json` is untouched. `make lint`, `gofmt -l .`,
  ┃      `make test` and `make regen-goldens` are all clean.
  ┃   40
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/inbox && git commit -q -m "inbox: reviewer reports slice-e cycle 2" &
  ┃  & git log --oneline -3 && git status --short
  ┃
  ┃  8d4bfe4 inbox: reviewer reports slice-e cycle 2
  ┃  be543d0 docs: the ledger closes slice E's five findings, and files one against the sweep
  ┃  b70e859 work: slice E's cycle 2 — the last review
  ┃  ?? logo/
  ┃

     + Thought · 2.6s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && MSG="$(head -1 work/inbox/reviewer-slice-e-cycle-2.md)
  ┃  F-86 (medium) — the sweep's command and AC-04's list of old names omit *maître d'*, so one of the seven renamed terms is…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":
  ┃  "herdr:opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"done","completion_seq":487,"cwd":"/Users/thiago/
  ┃  Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,
  ┃  "screen_detection_skipped":true,"state_change_seq":487,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722",
  ┃  "terminal_title":"OC | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":"w1Z"}
  ┃  ,"type":"agent_prompted"}}
  ┃

     + Thought · 583ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && rm -rf /tmp/base-e /tmp/sweep2.txt /tmp/sweep3.txt /tmp/sweep_base.txt /tmp/sweep_
  ┃  head.txt; git status --short; echo "=== the two commits ==="; git log --oneline -3; echo "=== the tree's own files un…
  ┃
  ┃  (3 earlier lines)
  ┃  be543d0 docs: the ledger closes slice E's five findings, and files one against the sweep
  ┃  b70e859 work: slice E's cycle 2 — the last review
  ┃  === the tree's own files unchanged since the runner ===
  ┃   docs/audit.md                          | 12 ++++++-----
  ┃   work/inbox/reviewer-slice-e-cycle-2.md | 39 ++++++++++++++++++++++++++++++++++
  ┃   2 files changed, 46 insertions(+), 5 deletions(-)
  ┃

     Cycle 2 done — verdict: changes-needed (one medium).

     The five closures — each by a command

     ┌─────────┬─────────────────────────────────────────────────┬──────────────────────────────────────────────────────────────┐
     │ Finding │ Command                                         │ Result                                                       │
     ├─────────┼─────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────┤
     │ F-81    │ grep -n -i -E "line                             │ 7 hits / 7 names, one each; 0 outside an _Avoid_: line (     │
     │         │ cook|taster|expo|proof|special|prep note|maître │ GLOSSARY.md:161)                                             │
     │         │ d'" GLOSSARY.md                                 │                                                              │
     ├─────────┼─────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────┤
     │ F-82    │ the crew section                                │ 5 entries; console in "### The strings" headed maitre; prose │
     │         │                                                 │ says it is a surface and names what it carries               │
     ├─────────┼─────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────┤
     │ F-83    │ same grep, filtered                             │ prose points at ADR 0009 § 3, which says the thing (the two  │
     │         │                                                 │ file types)                                                  │
     ├─────────┼─────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────┤
     │ F-84    │ three lists vs format.go:20                     │ docs/03:451, docs/03:552, lifecycle.go:200 all read id       │
     │         │                                                 │ prefixes · directories · root files · field names — what the │
     │         │                                                 │ code hashes                                                  │
     ├─────────┼─────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────┤
     │ F-85    │ go test -run                                    │ DomainVersion still 6; 7 refused, 5 accepted; prose says     │
     │         │ 'TestDomainVersion|TestDomainShape'             │ version fences the writer, fingerprint is the tripwire — no  │
     │         │                                                 │ overclaim either way                                         │
     └─────────┴─────────────────────────────────────────────────┴──────────────────────────────────────────────────────────────┘

     The sweep's count

     Base 72783a7 → 78 lines; HEAD → 77. The delta is exactly one line: GLOSSARY.md:59, F-83's prose. Nothing added — the
     console's _Avoid_ line already existed, and the name F-81 put on it is one the pattern cannot match. The count reconciles.
     All 77 lines name something that still exists.

     F-86 (medium, new)

     The sweep's list of old names (work/briefs/slice-e.md:81) and its command (commit c1c1e7c) omit maître d' — one of the seven
     renamed terms is outside the instrument's reach, and the published count hides it. The tree is clean (only ADR records), so
     the remedy is the command's: carry the name, exclude the glossary's _Avoid_ lines. Extending the pattern returns78, the
     extra line being exactly that _Avoid_ line.

     Mascot, scope, ledger

     Mascot did not land (logo/ never committed, no placement point carries art); goldens unmoved; git status shows only ?? logo/
     . The five fixes touched only the named files — lifecycle.go's diff is comments only; gate table, plan format, DomainVersion
     6, 42 commands and kitchen/activity.json unmoved; lint/test/goldens clean. Ledger: F-81..F-85 closed with file + evidence,
     F-86 filed, round table rows for Slice E and Slice E · cycle 2 (last slice of the rename).

     LOW notes (message only): the console's entry under "### The strings" describes the console, not a string — but the heading
     alone can be misread; docs/03 §6's #the-version-guard anchor has no heading (pre-existing since 813956f) — a broken pointer,
     one line to fix.

     Commits: be543d0 (ledger), 8d4bfe4 (receipt). Orchestrator pane w1Z:p1 notified.

     Build · MiniMax-M3.1-Flash · 30m 27s · 126.6 tok/s

  ┃
  ┃  You are the **reviewer** — the same lane that reviewed slice E. The orchestrator's pane
  ┃  id is `w1Z:p1`; use it verbatim in the notify command below.
  ┃
  ┃  ## Cycle 3 — the last fixes to slice E
  ┃
  ┃  The runner fixed F-86 and the two notes you left. **Cycle 3 is blockers only**: file
  ┃  **HIGH or MEDIUM** blockers. No nits, no cosmetics, no style.
  ┃
  ┃  ### F-86 — the instrument was blind, not the tree
  ┃
  ┃  **Check the instrument, because the tree did not change.**
  ┃
  ┃  - The sweep's command now carries **all seven** retired names — *maître d'* among
  ┃    them — and it **excludes the glossary's `_Avoid_` lines**, so a warning about a
  ┃    retired word is not counted as drift.
  ┃  - **Run the command yourself** and check the output is the same the commit body
  ┃    shows. The command is in the commit body and in `work/inbox/runner-slice-e-fix-2.md`.
  ┃  - **Check the count reconciles.** You counted **77 lines in 19 files**. The runner
  ┃    reports **72 in 18**, and explains the movement: the old command counted **five of
  ┃    the glossary's six `_Avoid_` lines** as drift, and the complete command sees all six
  ┃    and excludes all six. **Do the arithmetic**: 77 − 5 = 72, and one file drops out
  ┃    (the glossary's own). Confirm that the line the new command finds where the old
  ┃    could not — *maître d'* — is in a place that is legitimately excluded.
  ┃  - **Check the drift that remains is all legitimate**, file by file, as in cycle 2.
  ┃
  ┃  ### The two notes
  ┃
  ┃  - **The glossary's heading is now "### The surfaces"**, and the section says *"The
  ┃    console is a **surface**, not a job"* — ADR-0003's word. Check the taxonomy reads:
  ┃    `grep -n "^### " GLOSSARY.md` → **The plan**, **The crew**, **The work**, **The
  ┃    surfaces**. The console must not be filed as a job, and must not be filed as a
  ┃    string.
  ┃  - **`docs/03` § 6's pointer now lands** on a new `### 6.2 The version guard`, and
  ┃    the implementation is named in the same paragraph (`core.DomainVersion`, and
  ┃    `Store.checkDomainVersion` in `internal/core/lock.go`). **Follow the link** — it
  ┃    must resolve, and the section it lands on must say what the code says.
  ┃
  ┃  ### The mascot and the tree
  ┃
  ┃  **The mascot did not land** — the art is in `logo/` (untracked), and no placement
  ┃  carries a drawing. **AC-01 moves to the drawing session.** Check that nothing of the
  ┃  other session's work was committed, and that `kitchen/activity.json` did not move.
  ┃
  ┃  ### Scope
  ┃
  ┃  F-86's commit **changed no file** — the fix was the instrument, in the commit body
  ┃  and the brief. Check that is true, and that the two notes' commits touched only the
  ┃  glossary's heading and `docs/03`.
  ┃
  ┃  ### The ledger — and this is the last round
  ┃
  ┃  Move **F-86** from open to closed, with the evidence that now holds it. **Add the
  ┃  row for slice E itself**: the rename's last slice — the glossary rewritten, ADR 0009
  ┃  written, the tree swept, and **the mascot's acceptance criterion carried to the
  ┃  session that owns the drawing**. Commit the ledger.
  ┃
  ┃  The ledger is not swept: keep the findings' wording as they were filed.
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
  ┃  cat > work/inbox/reviewer-slice-e-cycle-3.md <<EOF
  ┃  Reviewer: slice E cycle 3 done. Verdict: <ok|changes-needed>. F-86 closed: <yes|no>. Filed <N> new blockers (<F-87.. or
  ┃  none>). Sweep count: <the number, and whether it reconciles>. Mascot: <did it land? yes/no>. LOW notes: <any>.
  ┃  EOF
  ┃  git add work/inbox && git commit -q -m "inbox: reviewer reports slice-e cycle 3"
  ┃  herdr agent prompt "w1Z:p1" "Reviewer: slice E cycle 3 done. Verdict: <ok|changes-needed>. F-86 closed: <yes|no>. Filed <N>
  ┃  new blockers (<F-87.. or none>). Sweep count: <the number, and whether it reconciles>. Mascot: <did it land? yes/no>. LOW
  ┃  notes: <any>."
  ┃  ```
  ┃

     + Thought · 136ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -8 && echo "=== status ===" && git status --short && echo "=== d
  ┃  iff stat since my last ===" && git diff 8d4bfe4..HEAD --stat
  ┃
  ┃  (12 earlier lines)
  ┃   docs/03-domain-model.md            |   7 +-
  ┃   work/briefs/fix-e2.md              | 106 +++++++++++++++++++++++++++++
  ┃   work/briefs/review-e3.md           |  83 +++++++++++++++++++++++
  ┃   work/briefs/slice-e.md             |  11 ++-
  ┃   work/inbox/runner-slice-e-fix-2.md | 133 +++++++++++++++++++++++++++++++++++++
  ┃   6 files changed, 338 insertions(+), 4 deletions(-)
  ┃

     + Thought · 462ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && for c in cea2380 08fa88c 7e12f53 d0e32bd b3e1077; do echo "=== $c ==="; git show -
  ┃  -stat --format='%an | %s' $c | head -12; done
  ┃
  ┃  (18 earlier lines)
  ┃   1 file changed, 133 insertions(+)
  ┃  === b3e1077 ===
  ┃  Luiz Thiago | work: slice E's last review
  ┃
  ┃   work/briefs/review-e3.md | 83 ++++++++++++++++++++++++++++++++++++++++++++++++
  ┃   1 file changed, 83 insertions(+)
  ┃

     + Thought · 246ms

     → Explored: 1 read

     + Thought · 337ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rnE --exclude-dir=.git --exclude-dir=work --exclude-dir=logo --exclude-dir=b
  ┃  in --exclude-dir=kitchen \…
  ┃
  ┃  (14 earlier lines)
  ┃     1 ./internal/core/plan_ids_test.go
  ┃     1 ./internal/core/migration_test.go
  ┃     1 ./internal/core/gates_test.go
  ┃     1 ./internal/core/docs_test.go
  ┃     1 ./internal/core/approval_test.go
  ┃     1 ./docs/04-operations.md
  ┃

     + Thought · 1.4s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && echo "=== complete pattern, WITHOUT the glossary exclusion ==="; grep -rnE --exclu
  ┃  de-dir=.git --exclude-dir=work --exclude-dir=logo --exclude-dir=bin --exclude-dir=kitchen \…
  ┃
  ┃  (14 earlier lines)
  ┃  < ./GLOSSARY.md:_Avoid_: Line cook, executor, implementer
  ┃  < ./GLOSSARY.md:_Avoid_: Proof, test result, verification
  ┃  < ./GLOSSARY.md:_Avoid_: Special, prep note, note, brief, ticket
  ┃  < ./GLOSSARY.md:_Avoid_: Taster, critic, judge
  ┃
  ┃  Exited with code 1
  ┃

     + Thought · 2.3s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cat -n /tmp/sweep_new.txt
  ┃
  ┃  (65 earlier lines)
  ┃      66  ./docs/03-domain-model.md:585:is `completion` now. A criterion's record was `proofs`; it is `evidence` now — the
  ┃      67  ./docs/03-domain-model.md:590:would write `proofs` where the version-6 writer writes `evidence` — the version-6
  ┃      68  ./docs/03-domain-model.md:597:renamed fields: a plan written before the rename keeps `proofs` and `evidence`, and
  ┃      69  ./docs/04-operations.md:242:The names the rename retired — `line-cook`, `taster`, `expo`, `coordinator` — still name
  ┃      70  ./prototypes/README.md:20:the recipe is called the *proof contract*, and the tab bar has five tabs. They are the
  ┃      71  ./scripts/stamp-plan/main.go:13:// (`specials+prep-notes` → `backlog`), or the names a record's field had
  ┃      72  ./scripts/stamp-plan/main.go:14:// (`proofs` → `evidence`).
  ┃

     + Thought · 1.6s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && echo "=== the taxonomy ==="; grep -n "^### \|^## " GLOSSARY.md; echo; echo "=== th
  ┃  e surfaces section ==="; sed -n '151,162p' GLOSSARY.md; echo; echo "=== docs/03 headings ==="; grep -n "^#\{2,3\} " d…
  ┃
  ┃  (31 earlier lines)
  ┃  634:### Writing — decisions (human)
  ┃  644:### Control
  ┃
  ┃  === anchors in docs/03 ===
  ┃  323:defined once, in [§ 4](#4-gates).
  ┃  456:See [the version guard](#62-the-version-guard).
  ┃

     + Thought · 702ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '544,602p' docs/03-domain-model.md
  ┃
  ┃  (52 earlier lines)
  ┃  — so a plan written before the change is still readable. The same is true of the
  ┃  renamed fields: a plan written before the rename keeps `proofs` and `evidence`, and
  ┃  the reader reads them. See [07-vocabulary.md](07-vocabulary.md).
  ┃
  ┃  ---
  ┃
  ┃  ## 7. The operation set
  ┃

     + Thought · 1.2s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && ls internal/core/lock.go && grep -n "checkDomainVersion" internal/core/*.go | head
  ┃  ; echo "=== the guard's body ==="; sed -n "$(grep -n 'func (s \*Store) checkDomainVersion' internal/core/lock.go | cu…
  ┃
  ┃  (38 earlier lines)
  ┃  75-  plan := PlanDomainVersion(j)
  ┃  76-  if plan > DomainVersion {
  ┃  77-    return fmt.Errorf("this plan was written under domain version %d and this tool speaks version %d: upgrade the tool,
  ┃  or read the plan without writing to it", plan, DomainVersion)
  ┃  78-  }
  ┃  79-  return nil
  ┃  80-}
  ┃

     + Thought · 1.3s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "DomainVersion" internal/core/*.go | grep -v "_test" | head -20
  ┃
  ┃  (9 earlier lines)
  ┃  internal/core/lock.go:43:      if err := s.checkDomainVersion(); err != nil {
  ┃  internal/core/lock.go:65:// checkDomainVersion refuses a write to a plan written under newer rules.
  ┃  internal/core/lock.go:70:func (s *Store) checkDomainVersion() error {
  ┃  internal/core/lock.go:75:  plan := PlanDomainVersion(j)
  ┃  internal/core/lock.go:76:  if plan > DomainVersion {
  ┃  internal/core/lock.go:77:    return fmt.Errorf("this plan was written under domain version %d and this tool speaks version
  ┃  %d: upgrade the tool, or read the plan without writing to it", plan, DomainVersion)
  ┃  internal/core/views.go:330:  "PlanDomainVersion",
  ┃

     + Thought · 411ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '125,150p' internal/core/activity.go
  ┃
  ┃  (19 earlier lines)
  ┃    j.Events = append(j.Events, e)
  ┃
  ┃    if err := writeJSONAtomic(filepath.Join(s.Root, FileActivity), j); err != nil {
  ┃      return Journal{}, err
  ┃    }
  ┃    return j, nil
  ┃  }
  ┃

     + Thought · 1.9s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'PY'
  ┃  import re,unicodedata…
  ┃
  ┃  headings: 31
  ┃  OK   4-gates
  ┃  OK   62-the-version-guard
  ┃

     + Thought · 705ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && echo "=== logo/ ever committed? ==="; git log --all --oneline --diff-filter=A -- l
  ┃  ogo/; echo "(empty = never)"; echo "=== kitchen/activity.json ==="; git diff 8d4bfe4..HEAD --stat -- kitchen/; echo "…
  ┃
  ┃  (46 earlier lines)
  ┃  -`core.DomainVersion`.
  ┃  +`core.DomainVersion`; the guard itself is `Store.checkDomainVersion` in
  ┃  +`internal/core/lock.go`.
  ┃
  ┃   **The version fences the writer; the fingerprint is the tripwire.** The fingerprint
  ┃   hashes the gates, the transitions, the lifecycles, the events, the chef's moves, the
  ┃

     + Thought · 10.5s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make lint 2>&1 | tail -2; gofmt -l .; echo "gofmt clean"; make test 2>&1 | tail -5
  ┃  ; make regen-goldens >/dev/null 2>&1; echo "=== after regen ==="; git status --short; echo "=== DomainVersion / comma…
  ┃
  ┃  (15 earlier lines)
  ┃  GLOSSARY.md
  ┃  docs/03-domain-model.md
  ┃  work/briefs/fix-e2.md
  ┃  work/briefs/review-e3.md
  ┃  work/briefs/slice-e.md
  ┃  work/inbox/runner-slice-e-fix-2.md
  ┃

     + Thought · 1.1s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git show 7e12f53 -- work/briefs/slice-e.md | cat
  ┃
  ┃  (21 earlier lines)
  ┃  +
  ┃  +> The orchestrator's first list of names omitted *maître d'*, and the review found it
  ┃  +> (F-86): the tree was clean and the **instrument** was blind. An instrument that
  ┃  +> cannot see a name is not a check.
  ┃
  ┃   - `grep` the tree and prove it. Report the exact commands and their output in the commit message.
  ┃   - **`docs/audit.md` is the ledger of the night — do not sweep it.** Its rows record the names as they were when the finding
  ┃  was filed; rewriting them would rewrite the record. Leave them, and say so in your report.
  ┃

     + Thought · 2.3s

     → Explored: 1 read

     + Thought · 7.5s

     All checks pass. Now the ledger.

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    84   |---|---|---|---|                                        84   |---|---|---|---|
  ┃    85   | **F-59** | medium | **`core.Audit.Decisions` has a     85   | **F-59** | medium | **`core.Audit.Decisions` has a
  ┃         shape and no writer** — it is where a decision *about         shape and no writer** — it is where a decision *about
  ┃         the plan* would live, and nothing writes it. `kitchen/        the plan* would live, and nothing writes it. `kitchen/
  ┃         decisions.json` is named by the plan format and opened        decisions.json` is named by the plan format and opened
  ┃         by nothing. The remedy is on the table in [06-proposals.      by nothing. The remedy is on the table in [06-
  ┃         md](06-proposals.md) § 2. | `internal/core/model.go:          proposals.md](06-proposals.md) § 2. | `internal/core/
  ┃         270`; `internal/core/store.go:24` |                           model.go:270`; `internal/core/store.go:24` |
  ┃    86   | **F-60** | medium | **`core.Questions` has a writer-   86   | **F-60** | medium | **`core.Questions` has a writer-
  ┃         less shape, so the `questions.resolved` gate can never        less shape, so the `questions.resolved` gate can never
  ┃         fire** — nothing ever sets a question, so the list is         fire** — nothing ever sets a question, so the list is
  ┃         always empty and the gate always passes. A gate nobody        always empty and the gate always passes. A gate nobody
  ┃         can fail is not a guard. | `internal/core/model.go:31`;       can fail is not a guard. | `internal/core/model.go:31`;
  ┃         `internal/core/lifecycle.go:197` |                            `internal/core/lifecycle.go:197` |
  ┃    87   | **F-61** | medium | **The `courses/` directory is      87   | **F-61** | medium | **The `courses/` directory is
  ┃         named by the plan format, created by `Init`, and never        named by the plan format, created by `Init`, and never
  ┃         written** — a directory no writer touches. Per [06-           written** — a directory no writer touches. Per [06-
  ┃         proposals.md](06-proposals.md) § 1 the honest answer is       proposals.md](06-proposals.md) § 1 the honest answer is
  ┃         to stop naming it: a course is the relation, not a file.      to stop naming it: a course is the relation, not a
  ┃         | `internal/core/store.go:18,63`; [docs/03 §6](03-            file. | `internal/core/store.go:18,63`; [docs/03 §6](
  ┃         domain-model.md) |                                            03-domain-model.md) |
  ┃    88 - | **F-86** | medium | **The rename's sweep cannot see
  ┃         the console's retired name — the name is missing from
  ┃         the instrument's own list of old names** — AC-04's list
  ┃         of names to sweep (`work/briefs/slice-e.md`:81) and the
  ┃         command that proves it (commit `c1c1e7c`) carry `Kitchen
  ┃         Brigade`, `scale`, `line-cook`, `taster`, `expo`,
  ┃         `Chopboard`, `special`, `prep note` and `Proof`, but
  ┃         **not** *maître d'* — so one of the seven renamed terms
  ┃         is outside the sweep's reach, and the count the sweep
  ┃         publishes (77) hides it. The name was missing because
  ┃         the glossary did not carry it either, and F-81's fix put
  ┃         it in `GLOSSARY.md:161`: the file now names it once, on
  ┃         purpose, and the command still cannot tell that line
  ┃         from drift. The tree is clean — the only other hits are
  ┃         the dated records (ADR 0009 § 2 and § 4, ADR 0004, ADR
  ┃         0005) — but the proof does not cover the term, so a
  ┃         *maître d'* that came back as a current name would not
  ┃         be caught. The remedy is the command's: the pattern must
  ┃         carry the name, and the command must then exclude the
  ┃         glossary's `_Avoid_` lines, which is what tells the
  ┃         warning apart from the drift. | `work/briefs/slice-e.md:
  ┃         81`; commit `c1c1e7c`; `GLOSSARY.md:161` |
  ┃    89   | F-56 | low | **The poll's fingerprint is blind to a    88   | F-56 | low | **The poll's fingerprint is blind to a
  ┃         write that keeps the mtime and the size.** It is path +       write that keeps the mtime and the size.** It is path +
  ┃         mtime + size per file, so two different contents can          mtime + size per file, so two different contents can
  ┃         share one fingerprint. The tool's own writes always           share one fingerprint. The tool's own writes always
  ┃         change the mtime, so in practice only a hand edit that        change the mtime, so in practice only a hand edit that
  ┃         restores the mtime is missed — and the chef's `r`             restores the mtime is missed — and the chef's `r`
  ┃         reloads it. | `internal/core/store.go` —                      reloads it. | `internal/core/store.go` —
  ┃         `DiskFingerprint`. |                                          `DiskFingerprint`. |
  ┃    90   | P-8 | — | `--depends-on=` writes `depends_on: [""]` —  89   | P-8 | — | `--depends-on=` writes `depends_on: [""]` —
  ┃         inert, the same as `--in-scope=`. | Recorded so nobody        inert, the same as `--in-scope=`. | Recorded so nobody
  ┃         calls it a regression later. Not a defect. |                  calls it a regression later. Not a defect. |
  ┃    91                                                            90
  ┃    92   ## Closed findings                                       91   ## Closed findings
  ┃
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    141   | F-82 | medium | **The glossary files the console      141   | F-82 | medium | **The glossary files the console
  ┃          among the crew, where its own new tables say it does          among the crew, where its own new tables say it does
  ┃          not belong.** *Closed in* `GLOSSARY.md` — "### The            not belong.** *Closed in* `GLOSSARY.md` — "### The
  ┃          crew" holds five entries (Chef · Sous-chef · Runner ·         crew" holds five entries (Chef · Sous-chef · Runner ·
  ┃          Reviewer · Orchestrator); the console moved to "### The       Reviewer · Orchestrator); the console moved to "###
  ┃          strings", headed `maitre` the way the names table             The strings", headed `maitre` the way the names table
  ┃          spells it (`:32`), and the section says the console is        spells it (`:32`), and the section says the console is
  ┃          a **surface**, not a job, and names what it carries —         a **surface**, not a job, and names what it carries —
  ┃          the window title, the tab bar and the CLI's help. The         the window title, the tab bar and the CLI's help. The
  ┃          entry itself still describes the console, not a string.       entry itself still describes the console, not a
  ┃          | `GLOSSARY.md:89`; `GLOSSARY.md:151-161`; `GLOSSARY.         string. | `GLOSSARY.md:89`; `GLOSSARY.md:151-161`;
  ┃          md:32` |                                                      `GLOSSARY.md:32` |
  ┃    142   | F-83 | medium | **Two retired words are named in      142   | F-83 | medium | **Two retired words are named in
  ┃          prose as well as in the `_Avoid_` line.** *Closed in*         prose as well as in the `_Avoid_` line.** *Closed in*
  ┃          `GLOSSARY.md` — the backlog item's prose points at [ADR       `GLOSSARY.md` — the backlog item's prose points at [
  ┃          0009](decisions/0009-vocabulary.md) § 3, which says           ADR 0009](decisions/0009-vocabulary.md) § 3, which
  ┃          what the two file types were, so the reader still             says what the two file types were, so the reader still
  ┃          learns what the backlog item replaced; the `_Avoid_`          learns what the backlog item replaced; the `_Avoid_`
  ┃          line keeps the retired words. | `GLOSSARY.md:59-63`;          line keeps the retired words. | `GLOSSARY.md:59-63`;
  ┃          `docs/decisions/0009-vocabulary.md` § 3 |                     `docs/decisions/0009-vocabulary.md` § 3 |
  ┃    143   | F-84 | medium | **The spec lists the fingerprint's    143   | F-84 | medium | **The spec lists the fingerprint's
  ┃          facts twice and the two lists disagree.** *Closed in*         facts twice and the two lists disagree.** *Closed in*
  ┃          `docs/03-domain-model.md` and `internal/core/lifecycle.       `docs/03-domain-model.md` and `internal/core/
  ┃          go` — all three lists (two in the spec, one in the            lifecycle.go` — all three lists (two in the spec, one
  ┃          `DomainVersion` comment) now say the same four things:        in the `DomainVersion` comment) now say the same four
  ┃          the id prefixes, the directories, the root files and          things: the id prefixes, the directories, the root
  ┃          the field names on disk, and that is what the code            files and the field names on disk, and that is what
  ┃          hashes. | `docs/03-domain-model.md:451`; `docs/03-            the code hashes. | `docs/03-domain-model.md:451`;
  ┃          domain-model.md:552`; `internal/core/lifecycle.go:200`;       `docs/03-domain-model.md:552`; `internal/core/
  ┃          `internal/core/format.go:20` — `planFormatFacts` |            lifecycle.go:200`; `internal/core/format.go:20` —
  ┃                                                                        `planFormatFacts` |
  ┃    144   | F-85 | medium | **The spec still states the coupling  144   | F-85 | medium | **The spec still states the coupling
  ┃          the slice split.** *Closed in* `docs/03-domain-model.         the slice split.** *Closed in* `docs/03-domain-model.
  ┃          md` § 6, § 6.1 and the version-4 history, `internal/          md` § 6, § 6.1 and the version-4 history, `internal/
  ┃          core/lifecycle.go`'s `DomainVersion` comment, and             core/lifecycle.go`'s `DomainVersion` comment, and
  ┃          `docs/07-vocabulary.md` — the **version fences the            `docs/07-vocabulary.md` — the **version fences the
  ┃          writer** and the **fingerprint is the tripwire**: a           writer** and the **fingerprint is the tripwire**: a
  ┃          change to what the tool writes moves both, a change to        change to what the tool writes moves both, a change to
  ┃          the facts list moves the fingerprint alone. The code          the facts list moves the fingerprint alone. The code
  ┃          did not move: `core.DomainVersion` is still **6**, and        did not move: `core.DomainVersion` is still **6**, and
  ┃          the guard still refuses a plan stamped 7 and accepts          the guard still refuses a plan stamped 7 and accepts
  ┃          one stamped 5. | `internal/core/version_test.go` — the        one stamped 5. | `internal/core/version_test.go` — the
  ┃          shape test (`recordedDomainShape = "b84b230161f4ab03"`)       shape test (`recordedDomainShape =
  ┃          , `TestDomainVersionRefusesANewerPlan`,                       "b84b230161f4ab03"`),
  ┃          `TestDomainVersionStampsAnOlderPlan`; `internal/core/         `TestDomainVersionRefusesANewerPlan`,
  ┃          lifecycle.go:216` |                                           `TestDomainVersionStampsAnOlderPlan`; `internal/core/
  ┃                                                                        lifecycle.go:216` |
  ┃                                                                  145 + | F-86 | medium | **The rename's sweep cannot see the
  ┃                                                                        console's retired name — the name is missing from the
  ┃                                                                        instrument's own list of old names.** *Closed in* the
  ┃                                                                        **instrument**, not the tree: AC-04's list of names (
  ┃                                                                        `work/briefs/slice-e.md`) now carries *maître d'* and
  ┃                                                                        says the command must carry every name and must tell a
  ┃                                                                        retired word from a **warning about** one; the command
  ┃                                                                        (commit `cea2380`, and `work/inbox/runner-slice-e-fix-
  ┃                                                                        2.md`) matches all seven retired names and excludes
  ┃                                                                        the glossary's `_Avoid_` lines. No file of the tree
  ┃                                                                        changed — the tree was clean, and the proof was blind.
  ┃                                                                        | the command, run at HEAD — **72 lines in 18 files**,
  ┃                                                                        and the count reconciles: 77 − 5 (the glossary's
  ┃                                                                        `_Avoid_` lines the old command counted as drift) =
  ┃                                                                        72; the complete pattern without the glossary
  ┃                                                                        exclusion returns 78, the extra line being the
  ┃                                                                        console's own `_Avoid_` line, and every other hit of
  ┃                                                                        *maître d'* is in the excluded dated records (ADR
  ┃                                                                        0004, 0005, 0009) or in the ledger |
  ┃    145                                                           146
  ┃    146   ## What this ledger cannot tell you                     147   ## What this ledger cannot tell you
  ┃    147                                                           148
  ┃    148   **Cycles 1–7's own findings are not recoverable.**      149   **Cycles 1–7's own findings are not recoverable.**
  ┃          Their per-finding reports lived in a                          Their per-finding reports lived in a
  ┃
  ┃
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃  Could not find oldString in /Users/thiago/Code/Pessoal/maitre/docs/audit.md. It must match exactly, including whitespace and
  ┃  indentation.
  ┃

     + Thought · 96ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "Slice E · cycle 2" docs/audit.md | cut -c1-400
  ┃
  ┃  46:| Slice E · cycle 2 | The fixes to slice E: the glossary's retired names, the console's place, the spec's two lists and
  ┃  the version/fingerprint split — **the last slice of the rename** | F-81..F-85 closed, each verified by a command and not by
  ┃  the diff; the sweep's count reconciles (78 → 77: F-83 removed one prose line and nothing was added, the console's `_Avoid_`
  ┃  line already existed — but th
  ┃

     → Explored: 1 read

     + Thought · 1.3s

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    43   | Slice D | The vocabulary rename: the per-criterion     43   | Slice D | The vocabulary rename: the per-criterion
  ┃         record is evidence, the recipe-level stamp is completion      record is evidence, the recipe-level stamp is
  ┃         | The two words hold, the reader tolerates the old names      completion | The two words hold, the reader tolerates
  ┃         and the tripwires hold both ways; the approval digest         the old names and the tripwires hold both ways; the
  ┃         did not travel with the rename, so a pre-rename plan is       approval digest did not travel with the rename, so a
  ┃         refused as tampered — one high and two mediums filed |        pre-rename plan is refused as tampered — one high and
  ┃                                                                       two mediums filed |
  ┃    44   | Slice D · cycle 2 | The fixes to slice D: the digest   44   | Slice D · cycle 2 | The fixes to slice D: the digest
  ┃         travels with the rename, the glossary covers the record       travels with the rename, the glossary covers the record
  ┃         it names, the version paragraph stops contradicting the       it names, the version paragraph stops contradicting the
  ┃         reader | F-78..F-80 closed — each verified by a               reader | F-78..F-80 closed — each verified by a
  ┃         mutation, and the fixes moved nothing but the three           mutation, and the fixes moved nothing but the three
  ┃         findings' own code |                                          findings' own code |
  ┃    45   | Slice E | The vocabulary rename: the paper — the       45   | Slice E | The vocabulary rename: the paper — the
  ┃         glossary rewritten, ADR 0009 recorded, the tree swept |       glossary rewritten, ADR 0009 recorded, the tree swept |
  ┃         The ids and the guard hold and the sweep's 78 lines are       The ids and the guard hold and the sweep's 78 lines are
  ┃         all records of something that still exists; the retired-      all records of something that still exists; the
  ┃         word pattern holds for seven renamed terms and misses         retired-word pattern holds for seven renamed terms and
  ┃         the console's, and the spec and the `DomainVersion`           misses the console's, and the spec and the
  ┃         comment still describe the fingerprint the way the slice      `DomainVersion` comment still describe the fingerprint
  ┃         found it — five mediums filed |                               the way the slice found it — five mediums filed |
  ┃    46   | Slice E · cycle 2 | The fixes to slice E: the          46   | Slice E · cycle 2 | The fixes to slice E: the
  ┃         glossary's retired names, the console's place, the            glossary's retired names, the console's place, the
  ┃         spec's two lists and the version/fingerprint split —          spec's two lists and the version/fingerprint split —
  ┃         **the last slice of the rename** | F-81..F-85 closed,         **the last slice of the rename** | F-81..F-85 closed,
  ┃         each verified by a command and not by the diff; the           each verified by a command and not by the diff; the
  ┃         sweep's count reconciles (78 → 77: F-83 removed one           sweep's count reconciles (78 → 77: F-83 removed one
  ┃         prose line and nothing was added, the console's               prose line and nothing was added, the console's
  ┃         `_Avoid_` line already existed — but the published            `_Avoid_` line already existed — but the published
  ┃         pattern cannot see the name F-81 put on it) — one medium      pattern cannot see the name F-81 put on it) — one
  ┃         filed, against the sweep's command, which does not look       medium filed, against the sweep's command, which does
  ┃         for the console's retired name |                              not look for the console's retired name |
  ┃                                                                  47 + | Slice E · cycle 3 | The last fixes to slice E: the
  ┃                                                                       sweep's instrument, and the two notes the review left |
  ┃                                                                       F-86 closed — the command carries all seven retired
  ┃                                                                       names and tells a warning from the drift (72 lines in
  ┃                                                                       18 files; the count reconciles, 77 − 5), the glossary's
  ┃                                                                       heading is **The surfaces** and `docs/03`'s pointer
  ┃                                                                       lands on § 6.2 — **the rename's last slice, closed**:
  ┃                                                                       the glossary rewritten, ADR 0009 written, the tree
  ┃                                                                       swept, and **the mascot's acceptance criterion carried
  ┃                                                                       to the session that owns the drawing** (AC-01) |
  ┃    47   | The rounds inside a cycle | Fix → verify → back to the 48   | The rounds inside a cycle | Fix → verify → back to
  ┃         executor → repeat, until the verifier is satisfied |          the executor → repeat, until the verifier is satisfied
  ┃         Three rounds in cycle 8: the first fix introduced three       | Three rounds in cycle 8: the first fix introduced
  ┃         new mediums, the second closed them |                         three new mediums, the second closed them |
  ┃    48                                                            49
  ┃    49   ### The discipline every round runs                      50   ### The discipline every round runs
  ┃    50                                                            51
  ┃
  ┃

     + Thought · 9.2s

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    44   | Slice D · cycle 2 | The fixes to slice D: the digest   44   | Slice D · cycle 2 | The fixes to slice D: the digest
  ┃         travels with the rename, the glossary covers the record       travels with the rename, the glossary covers the record
  ┃         it names, the version paragraph stops contradicting the       it names, the version paragraph stops contradicting the
  ┃         reader | F-78..F-80 closed — each verified by a               reader | F-78..F-80 closed — each verified by a
  ┃         mutation, and the fixes moved nothing but the three           mutation, and the fixes moved nothing but the three
  ┃         findings' own code |                                          findings' own code |
  ┃    45   | Slice E | The vocabulary rename: the paper — the       45   | Slice E | The vocabulary rename: the paper — the
  ┃         glossary rewritten, ADR 0009 recorded, the tree swept |       glossary rewritten, ADR 0009 recorded, the tree swept |
  ┃         The ids and the guard hold and the sweep's 78 lines are       The ids and the guard hold and the sweep's 78 lines are
  ┃         all records of something that still exists; the retired-      all records of something that still exists; the
  ┃         word pattern holds for seven renamed terms and misses         retired-word pattern holds for seven renamed terms and
  ┃         the console's, and the spec and the `DomainVersion`           misses the console's, and the spec and the
  ┃         comment still describe the fingerprint the way the slice      `DomainVersion` comment still describe the fingerprint
  ┃         found it — five mediums filed |                               the way the slice found it — five mediums filed |
  ┃    46   | Slice E · cycle 2 | The fixes to slice E: the          46   | Slice E · cycle 2 | The fixes to slice E: the
  ┃         glossary's retired names, the console's place, the            glossary's retired names, the console's place, the
  ┃         spec's two lists and the version/fingerprint split —          spec's two lists and the version/fingerprint split —
  ┃         **the last slice of the rename** | F-81..F-85 closed,         **the last slice of the rename** | F-81..F-85 closed,
  ┃         each verified by a command and not by the diff; the           each verified by a command and not by the diff; the
  ┃         sweep's count reconciles (78 → 77: F-83 removed one           sweep's count reconciles (78 → 77: F-83 removed one
  ┃         prose line and nothing was added, the console's               prose line and nothing was added, the console's
  ┃         `_Avoid_` line already existed — but the published            `_Avoid_` line already existed — but the published
  ┃         pattern cannot see the name F-81 put on it) — one medium      pattern cannot see the name F-81 put on it) — one
  ┃         filed, against the sweep's command, which does not look       medium filed, against the sweep's command, which does
  ┃         for the console's retired name |                              not look for the console's retired name |
  ┃    47   | Slice E · cycle 3 | The last fixes to slice E: the     47   | Slice E · cycle 3 | The last fixes to slice E: the
  ┃         sweep's instrument, and the two notes the review left |       sweep's instrument, and the two notes the review left |
  ┃         F-86 closed — the command carries all seven retired           F-86 closed — the command carries all seven retired
  ┃         names and tells a warning from the drift (72 lines in 18      names and tells a warning from the drift (72 lines in
  ┃         files; the count reconciles, 77 − 5), the glossary's          18 files; the count reconciles, 77 − 5), the glossary's
  ┃         heading is **The surfaces** and `docs/03`'s pointer           heading is **The surfaces** and `docs/03`'s pointer
  ┃         lands on § 6.2 — **the rename's last slice, closed**:         lands on § 6.2 — **the rename's last slice, closed**:
  ┃         the glossary rewritten, ADR 0009 written, the tree            the glossary rewritten, ADR 0009 written, the tree
  ┃         swept, and **the mascot's acceptance criterion carried        swept, and **the mascot's acceptance criterion carried
  ┃         to the session that owns the drawing** (AC-01) |              to the session that owns the drawing** (AC-01) |
  ┃                                                                  48 + | Slice E · closed | The rename's last slice — the
  ┃                                                                       glossary rewritten, ADR 0009 written, the tree swept |
  ┃                                                                       F-81..F-86 closed across three rounds, each verified by
  ┃                                                                       a command; the sweep's command carries all seven
  ┃                                                                       retired names and tells a warning from the drift, and
  ┃                                                                       its **72 lines in 18 files** are all records of
  ┃                                                                       something that still exists — **AC-01, the mascot's
  ┃                                                                       drawing, is carried to the session that owns it**, and
  ┃                                                                       nothing of that session's work was committed here |
  ┃    48   | The rounds inside a cycle | Fix → verify → back to the 49   | The rounds inside a cycle | Fix → verify → back to
  ┃         executor → repeat, until the verifier is satisfied |          the executor → repeat, until the verifier is satisfied
  ┃         Three rounds in cycle 8: the first fix introduced three       | Three rounds in cycle 8: the first fix introduced
  ┃         new mediums, the second closed them |                         three new mediums, the second closed them |
  ┃    49                                                            50
  ┃    50   ### The discipline every round runs                      51   ### The discipline every round runs
  ┃    51                                                            52
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff docs/audit.md; echo "=== table shape ==="; awk '/^\|/ {n=gsub(/\|/,"|");
  ┃  print NR": "n}' docs/audit.md | awk -F': ' '{print $2}' | sort | uniq -c
  ┃
  ┃  (27 earlier lines)
  ┃   ## What this ledger cannot tell you
  ┃
  ┃  === table shape ===
  ┃    27 3
  ┃    29 4
  ┃    60 5
  ┃

     + Thought · 694ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && awk 'NR>=25 && NR<=50 && /^\|/ {n=gsub(/\|/,"|"); if (n!=4) print NR": "n}' docs/a
  ┃  udit.md; echo "(no output = the cycle matrix is well-formed)"; git add docs/audit.md && git commit -q -m "docs: the l…
  ┃
  ┃  (no output = the cycle matrix is well-formed)
  ┃  4228daa docs: the ledger closes slice E's last finding, and the slice itself
  ┃  b3e1077 work: slice E's last review
  ┃

     + Thought · 2.4s

  ┃
  ┃  # Wrote work/inbox/reviewer-slice-e-cycle-3.md
  ┃
  ┃    1 Reviewer: slice E cycle 3 done. Verdict: ok. F-86 closed: yes. Filed 0 new blockers (none). Sweep count: 72 lines in 18
  ┃      files — it reconciles (77 − 5 = 72, and the complete pattern without the glossary's exclusion returns 78 in 19, the
  ┃      extra line being the console's own `_Avoid_` line). Mascot: no. LOW notes: the complete command itself is in no durable
  ┃      file — the brief carries the *rule* ("the command must carry every one of them") and the command lives in commit
  ┃      `cea2380` and in this lane's receipt, so the next reader has to dig for it.
  ┃    2
  ┃    3 ## F-86 — the instrument, run
  ┃    4
  ┃    5 The command from commit `cea2380`, run here at HEAD, unchanged:
  ┃    6
  ┃    7     $ grep -rnE --exclude-dir=.git --exclude-dir=work --exclude-dir=logo --exclude-dir=bin --exclude-dir=kitchen \
  ┃    8         --include="*.go" --include="*.md" --include="*.py" -i \
  ┃    9         "kitchen brigade|kitchen-brigade|chopboard|\bexpo\b|\btaster\b|line-cook|line cook|\bspecials?\b|prep[ -]note|\
  ┃      bproofs?\b|maître d'" . \
  ┃   10         | grep -v "^./docs/audit" | grep -v "^./docs/audits/" | grep -v "^./research/" \
  ┃   11         | grep -v "^./docs/decisions/0004-naming.md" | grep -v "^./docs/decisions/0005-vocabulary.md" \
  ┃   12         | grep -v "^./docs/decisions/0009-vocabulary.md" | grep -v "^./docs/07-vocabulary.md" \
  ┃   13         | grep -v "^./prototypes/tab-" | grep -v "^./GLOSSARY.md:[0-9]*:_Avoid_:" \
  ┃   14         | cut -d: -f1 | sort | uniq -c | sort -rn
  ┃   15
  ┃   16     15 ./internal/core/legacy.go      11 ./internal/core/roles_test.go    10 ./internal/core/model.go
  ┃   17      9 ./docs/03-domain-model.md        5 ./internal/core/lifecycle.go     4 ./internal/ui/console_test.go
  ┃   18      3 ./internal/core/legacy_test.go   2 ./scripts/stamp-plan/main.go     2 ./internal/core/views.go
  ┃   19      2 ./internal/core/backlog_test.go 2 ./internal/core/backlog_concurrency_test.go
  ┃   20      1 each: prototypes/README.md, internal/core/{plan_ids,migration,gates,docs,approval}_test.go, docs/04-operations.
  ┃      md
  ┃   21
  ┃   22     $ … | wc -l                       → 72
  ┃   23     $ … | cut -d: -f1 | sort -u | wc -l → 18
  ┃   24
  ┃   25 **The output is the commit body's, line for line.**
  ┃   26
  ┃   27 ## The count, reconciled
  ┃   28
  ┃   29 | Command | Lines | Files |
  ┃   30 |--:|--:|--:|
  ┃   31 | the old one (six names, no glossary exclusion) | 77 | 19 |
  ┃   32 | the complete pattern, **without** the glossary exclusion | **78** | 19 |
  ┃   33 | the complete command as it now stands | **72** | **18** |
  ┃   34
  ┃   35 77 − 5 = 72 ✓ — the old command counted five of the glossary's six `_Avoid_` lines as drift; the complete one sees all
  ┃      six and excludes all six, so the glossary drops out of the file list.
  ┃   36
  ┃   37 The line the new command finds where the old could not is *maître d'* at `GLOSSARY.md:161` — the glossary's own
  ┃      `_Avoid_` line, which the command excludes on purpose. The name's other hits are the dated records (ADR 0004, 0005,
  ┃      0009) and the ledger, all excluded. **The complete command finds no new drift.**
  ┃   38
  ┃   39 **All 72 lines name something that still exists**, file by file: the legacy reader (`LegacySpecial`, `LegacyPrepNote`,
  ┃      the old directories, the old `proofs` field it reads), the retired-role aliases and the deny-list's own tests, the old-
  ┃      id and migration fixtures, the dated version notes, the old digests' field names, the migration events `stamp-plan`
  ┃      records, `docs/04`'s deny-list statement, and `prototypes/README.md`'s note that the captures predate the rename.
  ┃   40
  ┃   41 ## The two notes
  ┃   42
  ┃   43 | Note | The command | What it returned |
  ┃   44 |---|---|---|
  ┃   45 | The heading | `grep -n "^### " GLOSSARY.md` | **The plan** (43) · **The crew** (89) · **The work** (114) · **The
  ┃      surfaces** (151) — the console is a surface, not a job and not a string; the crew still holds five entries |
  ┃   46 | The pointer | the anchors of `docs/03`, computed as GitHub computes them | `#62-the-version-guard` **resolves** (the
  ┃      new `### 6.2 The version guard`); `#4-gates` still resolves too |
  ┃   47
  ┃   48 The section the pointer lands on says what the code says: the guard is `Store.checkDomainVersion` in `internal/core/
  ┃      lock.go` — it refuses when `plan > DomainVersion` (`lock.go:70-79`), and the next write stamps the tool's own version (
  ┃      `internal/core/activity.go:139`), while reading is always allowed. `core.DomainVersion` is still **6**.
  ┃   49
  ┃   50 ## The mascot and the tree
  ┃   51
  ┃   52 The mascot **did not land** — `logo/` is untracked and was never added (`git log --all --diff-filter=A -- logo/` is
  ┃      empty); no placement carries a drawing; **AC-01 moves to the drawing session**. Nothing of the other session's work was
  ┃      committed: `git status --short` shows only `?? logo/`. `kitchen/activity.json` did not move.
  ┃   53
  ┃   54 ## Scope
  ┃   55
  ┃   56 | Commit | Files |
  ┃   57 |---|---|
  ┃   58 | `cea2380` (F-86) | **none** — the fix was the instrument, in the commit body and in AC-04's list (`work/briefs/slice-
  ┃      e.md`, which now carries *maître d'* and the rule that the command must tell a retired word from a warning about one) |
  ┃   59 | `08fa88c` (the two notes) | `GLOSSARY.md` — the heading alone; `docs/03-domain-model.md` — the link, the new `### 6.
  ┃      2` heading, and the implementation's name |
  ┃   60
  ┃   61 Nothing else moved: the gate table, the plan format, the transition table, `core.DomainVersion` (**6**), the command
  ┃      count (**42**) and `kitchen/activity.json` are all as the runner left them. `make lint`, `gofmt -l .`, `make test` and
  ┃      `make regen-goldens` are clean, and the goldens did not move.
  ┃   62
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/inbox && git commit -q -m "inbox: reviewer reports slice-e cycle 3" &
  ┃  & git log --oneline -3 && git status --short
  ┃
  ┃  99330d4 inbox: reviewer reports slice-e cycle 3
  ┃  4228daa docs: the ledger closes slice E's last finding, and the slice itself
  ┃  b3e1077 work: slice E's last review
  ┃  ?? logo/
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && MSG="$(head -1 work/inbox/reviewer-slice-e-cycle-3.md)
  ┃  F-86 closed: the instrument now carries all seven retired names and excludes the glossary's _Avoid_ lines; run here it r…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":
  ┃  "herdr:opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"done","completion_seq":495,"cwd":"/Users/thiago/
  ┃  Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,
  ┃  "screen_detection_skipped":true,"state_change_seq":495,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722",
  ┃  "terminal_title":"OC | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":"w1Z"}
  ┃  ,"type":"agent_prompted"}}
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre && rm -f /tmp/sweep_new.txt /tmp/sweep_old.txt /tmp/sweep_complete.txt; git status --
  ┃  short; echo "=== my two commits ==="; git diff --stat b3e1077..HEAD
  ┃
  ┃  ?? logo/
  ┃  === my two commits ===
  ┃   docs/audit.md                          |  4 ++-
  ┃   work/inbox/reviewer-slice-e-cycle-3.md | 61 ++++++++++++++++++++++++++++++++++
  ┃   2 files changed, 64 insertions(+), 1 deletion(-)
  ┃

     Cycle 3 done — verdict: ok, no blockers filed.

     F-86 — the instrument, run

     ┌─────────────────────────────────────────────────────────────────────┬──────────────────────────┬─────────────────────────┐
     │ Command                                                             │ Lines                    │ Files                   │
     ├─────────────────────────────────────────────────────────────────────┼──────────────────────────┼─────────────────────────┤
     │ old (six names, no glossary exclusion)                              │ 77                       │ 19                      │
     ├─────────────────────────────────────────────────────────────────────┼──────────────────────────┼─────────────────────────┤
     │ complete pattern, without the glossary exclusion                    │ 78                       │ 19                      │
     ├─────────────────────────────────────────────────────────────────────┼──────────────────────────┼─────────────────────────┤
     │ complete command as it stands                                       │ 72                       │ 18                      │
     └─────────────────────────────────────────────────────────────────────┴──────────────────────────┴─────────────────────────┘

     77 − 5 = 72 ✓ — the old command counted five of the glossary's six _Avoid_ lines as drift; the complete one sees all six and
     excludes all six, so the glossary leaves the file list. The line the new command finds where the old could not is maître d'
     at GLOSSARY.md:161 — the glossary's own _Avoid_ line, excluded on purpose; its other hits are the excluded dated records (
     ADR 0004/0005/0009) and the ledger. No new drift. The per-file output matches commit cea2380's body line for line.

     All 72 lines name something that still exists, file by file: the legacy reader, the retired-role aliases and deny-list
     tests, the old-id and migration fixtures, the dated version notes, the old digests' field names, stamp-plan's migration
     events, docs/04's deny-list statement, prototypes/README.md's dated capture note.

     The two notes

     - grep -n "^### " GLOSSARY.md → The plan · The crew (five entries) · The work · The surfaces — the console is a surface, not
       a job and not a string.
     - docs/03's pointer resolves: #62-the-version-guard (new ### 6.2 The version guard); #4-gates still resolves. The section
       says what the code says — Store.checkDomainVersion (internal/core/lock.go:70) refuses when plan > DomainVersion, the next
       write stamps the tool's version (activity.go:139), reading stays allowed; core.DomainVersion is still 6.

     Mascot, scope, ledger

     Mascot did not land (logo/ never committed, no placement carries art) — AC-01 moves to the drawing session. Nothing of the
     other session was committed; kitchen/activity.json unmoved; tree shows only ?? logo/. cea2380 changed no file (the
     instrument lives in the commit body and in AC-04's list, now carrying maître d' and the rule that a retired word must be
     told from a warning about one); 08fa88c touched only the glossary's heading and docs/03. Lint, gofmt, tests and goldens
     clean.

     Ledger: F-86 moved to closed with the command and the arithmetic; rows added for Slice E · cycle 3 and Slice E · closed (
     glossary rewritten, ADR 0009 written, tree swept, mascot's AC carried to the drawing session).

     LOW note (message only): the complete command itself is in no durable file — the brief carries the rule, the command lives
     in cea2380 and in the receipt, so the next reader has to dig for it.

     Commits: 4228daa (ledger), 99330d4 (receipt). Orchestrator pane w1Z:p1 notified.

     Build · MiniMax-M3.1-Flash · 2m 10s · 124.2 tok/s

  ┃
  ┃
  ┃
  ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
  ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
  ~/Code/Pessoal/maitre:master                                                                      102.6K (10%)  ctrl+p commands

```
