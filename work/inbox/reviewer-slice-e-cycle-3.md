Reviewer: slice E cycle 3 done. Verdict: ok. F-86 closed: yes. Filed 0 new blockers (none). Sweep count: 72 lines in 18 files — it reconciles (77 − 5 = 72, and the complete pattern without the glossary's exclusion returns 78 in 19, the extra line being the console's own `_Avoid_` line). Mascot: no. LOW notes: the complete command itself is in no durable file — the brief carries the *rule* ("the command must carry every one of them") and the command lives in commit `cea2380` and in this lane's receipt, so the next reader has to dig for it.

## F-86 — the instrument, run

The command from commit `cea2380`, run here at HEAD, unchanged:

    $ grep -rnE --exclude-dir=.git --exclude-dir=work --exclude-dir=logo --exclude-dir=bin --exclude-dir=kitchen \
        --include="*.go" --include="*.md" --include="*.py" -i \
        "kitchen brigade|kitchen-brigade|chopboard|\bexpo\b|\btaster\b|line-cook|line cook|\bspecials?\b|prep[ -]note|\bproofs?\b|maître d'" . \
        | grep -v "^./docs/audit" | grep -v "^./docs/audits/" | grep -v "^./research/" \
        | grep -v "^./docs/decisions/0004-naming.md" | grep -v "^./docs/decisions/0005-vocabulary.md" \
        | grep -v "^./docs/decisions/0009-vocabulary.md" | grep -v "^./docs/07-vocabulary.md" \
        | grep -v "^./prototypes/tab-" | grep -v "^./GLOSSARY.md:[0-9]*:_Avoid_:" \
        | cut -d: -f1 | sort | uniq -c | sort -rn

    15 ./internal/core/legacy.go      11 ./internal/core/roles_test.go    10 ./internal/core/model.go
     9 ./docs/03-domain-model.md        5 ./internal/core/lifecycle.go     4 ./internal/ui/console_test.go
     3 ./internal/core/legacy_test.go   2 ./scripts/stamp-plan/main.go     2 ./internal/core/views.go
     2 ./internal/core/backlog_test.go 2 ./internal/core/backlog_concurrency_test.go
     1 each: prototypes/README.md, internal/core/{plan_ids,migration,gates,docs,approval}_test.go, docs/04-operations.md

    $ … | wc -l                       → 72
    $ … | cut -d: -f1 | sort -u | wc -l → 18

**The output is the commit body's, line for line.**

## The count, reconciled

| Command | Lines | Files |
|--:|--:|--:|
| the old one (six names, no glossary exclusion) | 77 | 19 |
| the complete pattern, **without** the glossary exclusion | **78** | 19 |
| the complete command as it now stands | **72** | **18** |

77 − 5 = 72 ✓ — the old command counted five of the glossary's six `_Avoid_` lines as drift; the complete one sees all six and excludes all six, so the glossary drops out of the file list.

The line the new command finds where the old could not is *maître d'* at `GLOSSARY.md:161` — the glossary's own `_Avoid_` line, which the command excludes on purpose. The name's other hits are the dated records (ADR 0004, 0005, 0009) and the ledger, all excluded. **The complete command finds no new drift.**

**All 72 lines name something that still exists**, file by file: the legacy reader (`LegacySpecial`, `LegacyPrepNote`, the old directories, the old `proofs` field it reads), the retired-role aliases and the deny-list's own tests, the old-id and migration fixtures, the dated version notes, the old digests' field names, the migration events `stamp-plan` records, `docs/04`'s deny-list statement, and `prototypes/README.md`'s note that the captures predate the rename.

## The two notes

| Note | The command | What it returned |
|---|---|---|
| The heading | `grep -n "^### " GLOSSARY.md` | **The plan** (43) · **The crew** (89) · **The work** (114) · **The surfaces** (151) — the console is a surface, not a job and not a string; the crew still holds five entries |
| The pointer | the anchors of `docs/03`, computed as GitHub computes them | `#62-the-version-guard` **resolves** (the new `### 6.2 The version guard`); `#4-gates` still resolves too |

The section the pointer lands on says what the code says: the guard is `Store.checkDomainVersion` in `internal/core/lock.go` — it refuses when `plan > DomainVersion` (`lock.go:70-79`), and the next write stamps the tool's own version (`internal/core/activity.go:139`), while reading is always allowed. `core.DomainVersion` is still **6**.

## The mascot and the tree

The mascot **did not land** — `logo/` is untracked and was never added (`git log --all --diff-filter=A -- logo/` is empty); no placement carries a drawing; **AC-01 moves to the drawing session**. Nothing of the other session's work was committed: `git status --short` shows only `?? logo/`. `kitchen/activity.json` did not move.

## Scope

| Commit | Files |
|---|---|
| `cea2380` (F-86) | **none** — the fix was the instrument, in the commit body and in AC-04's list (`work/briefs/slice-e.md`, which now carries *maître d'* and the rule that the command must tell a retired word from a warning about one) |
| `08fa88c` (the two notes) | `GLOSSARY.md` — the heading alone; `docs/03-domain-model.md` — the link, the new `### 6.2` heading, and the implementation's name |

Nothing else moved: the gate table, the plan format, the transition table, `core.DomainVersion` (**6**), the command count (**42**) and `kitchen/activity.json` are all as the runner left them. `make lint`, `gofmt -l .`, `make test` and `make regen-goldens` are clean, and the goldens did not move.
