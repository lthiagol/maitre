Runner: slice E fix 2 done. Commits: 2. Fixed: F-86 and the two notes. Sweep count now: 72 lines in 18 files, and it moves by −5 from the reviewer's 77: the old command counted five of the glossary's six `_Avoid_` lines as drift, and the complete command sees all six and excludes all six. Notes: the heading is **### The surfaces** — the console is a surface (ADR-0003's word, and the entry's own first word), and "the strings" claimed the console *is* one; the broken pointer now lands on a new **§ 6.2 The version guard**, the document's own description, with the implementation named in the same paragraph; no script, for the reason the orchestrator gave.

## The commits

| Commit | What |
|---|---|
| `cea2380` | **F-86** — the sweep's instrument: the complete command, its output, and the reconciliation. No file changed: the tree was clean |
| `08fa88c` | the two notes — the glossary's heading and the broken pointer |

## F-86 — the instrument, not the tree

**The tree is clean.** Nothing in the docs needed changing, and the F-86 commit
changes no file. What was wrong was the instrument: the pattern omitted *maître d'*,
so one of the seven renamed terms was outside its reach — the sweep would have
passed whatever the tree said about the console's old name, because it never asked.

The command now carries all seven retired names **and** knows the difference between
a retired word and a **warning** about one: the glossary's `_Avoid_` lines name them
on purpose, so the sweep excludes those lines, and a retired name outside one is
drift.

    grep -rnE --exclude-dir=.git --exclude-dir=work --exclude-dir=logo --exclude-dir=bin --exclude-dir=kitchen \
      --include="*.go" --include="*.md" --include="*.py" -i \
      "kitchen brigade|kitchen-brigade|chopboard|\bexpo\b|\btaster\b|line-cook|line cook|\bspecials?\b|prep[ -]note|\bproofs?\b|maître d'" . \
      | grep -v "^./docs/audit" \
      | grep -v "^./docs/audits/" \
      | grep -v "^./research/" \
      | grep -v "^./docs/decisions/0004-naming.md" \
      | grep -v "^./docs/decisions/0005-vocabulary.md" \
      | grep -v "^./docs/decisions/0009-vocabulary.md" \
      | grep -v "^./docs/07-vocabulary.md" \
      | grep -v "^./prototypes/tab-" \
      | grep -v "^./GLOSSARY.md:[0-9]*:_Avoid_:"

**Its output** — the command's, piped through a counter:

    $ sh sweep.sh | cut -d: -f1 | sort | uniq -c | sort -rn
         15 ./internal/core/legacy.go
         11 ./internal/core/roles_test.go
         10 ./internal/core/model.go
          9 ./docs/03-domain-model.md
          5 ./internal/core/lifecycle.go
          4 ./internal/ui/console_test.go
          3 ./internal/core/legacy_test.go
          2 ./scripts/stamp-plan/main.go
          2 ./internal/core/views.go
          2 ./internal/core/backlog_test.go
          2 ./internal/core/backlog_concurrency_test.go
          1 ./prototypes/README.md
          1 ./internal/core/plan_ids_test.go
          1 ./internal/core/migration_test.go
          1 ./internal/core/gates_test.go
          1 ./internal/core/docs_test.go
          1 ./internal/core/approval_test.go
          1 ./docs/04-operations.md

    $ sh sweep.sh | wc -l
    72
    $ sh sweep.sh | cut -d: -f1 | sort -u | wc -l
    18

**The count, reconciled.** The reviewer counted **77 lines in 19 files**:

| Count | What it is |
|--:|---|
| 77 | the old command — six names, and the glossary's `_Avoid_` lines counted as drift |
| 78 | the complete command **without** the glossary exclusion: it sees the console's name, and the glossary's six `_Avoid_` lines are in the set |
| 72 | the complete command as it now stands — the six `_Avoid_` lines excluded |

77 − 5 = 72: the old command saw **five** of the glossary's six `_Avoid_` lines, all
but the console's — which it could not name — and the new command sees all six and
excludes all six. The name the old instrument could not see adds **no** drift line:
*maître d'* lives in the glossary's `_Avoid_` line and in the ADR records (0004,
0005, 0009), and both are excluded.

The count is **72 in 18 files** after the two notes landed too — the headings carry
no retired names.

**No script, deliberately.** The sweep is a **one-time act** over a tree that now has
*legitimate* mentions of the old names — the legacy reader, the retired-role
aliases, the dated records. A permanent tripwire would need a permanent exclusion
list, and **a tripwire with a long exclusion list is a tripwire that has stopped
seeing**. The repo already has better instruments for the live facts: the deny-list
wire, the derived-views wire, the plan-format wire. The command lives here and in
`cea2380` — it is the record of the act, not a guard.

## The two notes

**The heading is `### The surfaces`.** A glossary is a **taxonomy**, and *The strings*
claimed the console *is* a string — it *carries* strings. The console is a
**surface**: ADR-0003's own word for it, and the first word of the entry's own
paragraph. The section still says what the console carries (the window title, the
tab bar, the CLI's help) — the heading no longer claims the console is one.

    $ grep -n "^### " GLOSSARY.md
    43:### The plan
    89:### The crew
    114:### The work
    151:### The surfaces

**The pointer lands.** `docs/03` § 6 linked to `#the-version-guard`, and there was
no such heading — the guard was described in prose at the end of § 6.1 with nothing
to land on. The description now has a heading of its own, **§ 6.2 The version
guard**, and the link follows it: `[the version guard](#62-the-version-guard)`.

I chose the **document** over the code, because a reader of the spec wants the
spec's words; the implementation is named in the same paragraph — `core.DomainVersion`
for the number and `Store.checkDomainVersion` in `internal/core/lock.go` for the
guard that refuses the write. The version history that follows is the guard's story:
every version paragraph ends with why the version moved, and that is what the guard
fences.

## The gate

    $ make lint && gofmt -l . && make test && git status --short
    go vet ./...
    go test ./...
    ok  	github.com/lthiagol/maitre/internal/cli	2.463s
    ok  	github.com/lthiagol/maitre/internal/core	(cached)
    ok  	github.com/lthiagol/maitre/internal/ui	(cached)
    ?? logo/

`gofmt -l .` printed nothing. The only thing in `git status` is `logo/` — the
mascot's art, untracked and not mine.

Untouched, as the brief required: the gate table, the plan format, the command count
(42), `core.DomainVersion` (6) and `kitchen/activity.json`.

## The mascot

**Still not landed.** `logo/` holds the art, untracked; the console's header, the
CLI's help, the README and the Sous-Chef tab carry no drawing. Nothing in this cycle
touched them, and the goldens did not move — no console output changed.
