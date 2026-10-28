Reviewer: slice E cycle 2 done. Verdict: changes-needed. F-81..F-85 closed: yes. Filed 1 new blocker (F-86). Sweep count: 77 lines in 19 files — it reconciles (78 → 77 is F-83's one prose line removed, nothing added), but the command cannot see *maître d'*. Mascot: no. LOW notes: the console's entry now sits under "### The strings" — it describes the console, not a string, and the section says what it carries, but the heading alone can be misread; `docs/03` §6's `#the-version-guard` anchor has no heading (pre-existing since `813956f`) — a broken pointer, one line to fix.

## The five closures, each by a command

| Finding | The command | What it returned |
|---|---|---|
| **F-81** | `grep -n -i -E "line cook\|taster\|expo\|proof\|special\|prep note\|maître d'" GLOSSARY.md` | 7 hits, 7 retired names, one each; **0** outside an `_Avoid_:` line (`GLOSSARY.md:161` carries *maître d'*) |
| **F-82** | the crew section | "### The crew" holds **five** entries (Chef · Sous-chef · Runner · Reviewer · Orchestrator); the console moved to "### The strings" (`:151`), headed `maitre` the way the names table spells it (`:32`), and the section says the console is a **surface** and names what it carries — the window title, the tab bar, the CLI's help |
| **F-83** | the same grep, filtered | the prose no longer names *special* or *prep note*; it points at ADR 0009 § 3, and § 3 does say the thing: "The backlog had **two file types** — a *special* (a hunch) and a *prep note*" |
| **F-84** | the three lists, against `internal/core/format.go:20` | `docs/03:451`, `docs/03:552` and `internal/core/lifecycle.go:200` all read "the id prefixes, the directories, the root files and the field names on disk" — and that is what `planFormatFacts` hashes |
| **F-85** | `go test ./internal/core -run 'TestDomainVersion\|TestDomainShape'` | `core.DomainVersion` is still **6**; a plan stamped **7** is refused, a plan stamped **5** is accepted and restamped; the prose now says the version fences the writer and the fingerprint is the tripwire, in `docs/03` §6/§6.1, the `DomainVersion` comment and `docs/07-vocabulary.md` — and it does not overclaim the other way |

## The sweep's count

The published command, run at the cycle-1 base (`72783a7`) and at HEAD, on the same tree:

| | Lines |
|---|---|
| base | 78 |
| HEAD | 77 |
| delta | exactly one line removed — `GLOSSARY.md:59`, the prose F-83 took the retired words out of |

Nothing was added: the console's `_Avoid_` line already existed, and the name F-81 put on it is one the pattern cannot match. **Every one of the 77 lines names something that still exists** — the legacy reader, the retired-role aliases, the old-id and migration tests, the dated version notes, the tripwires' own mutations, the glossary's `_Avoid_` lines, `prototypes/README.md`.

Extending the pattern with *maître d'* returns **78** — the extra line is exactly the glossary's own `_Avoid_` line. That is F-86: the command has to learn the difference between the glossary's warning and the drift.

## F-86 (medium) — the sweep cannot see the console's retired name

AC-04's list of old names (`work/briefs/slice-e.md:81`) and the command that proves it (commit `c1c1e7c`) carry `Kitchen Brigade`, `scale`, `line-cook`, `taster`, `expo`, `Chopboard`, `special`, `prep note` and `Proof` — **not** *maître d'*. One of the seven renamed terms is outside the instrument's reach, and the count the sweep publishes hides it.

The tree is clean: the only other hits are the dated records (ADR 0009 § 2 and § 4, ADR 0004, ADR 0005). But the proof does not cover the term, so a *maître d'* that came back as a current name would not be caught. The remedy is the command's, not the file's: the pattern must carry the name, and the command must then exclude the glossary's `_Avoid_` lines.

## The mascot and the tree

The mascot **did not land** — `logo/` is untracked and was never added (`git log --all --diff-filter=A -- logo/` is empty); no placement point carries a drawing; AC-01 moves to the drawing session. This slice committed **nothing** of the other session's work: `git status --short` shows only `?? logo/`. The goldens did not move (`internal/ui/testdata/*.golden` untouched) and `make regen-goldens` leaves the tree clean.

## Scope

The five fix commits touched `GLOSSARY.md`, `docs/03-domain-model.md`, `docs/07-vocabulary.md`, `README.md`, `work/README.md` and `internal/core/lifecycle.go` — and in `lifecycle.go` the diff is **comments only** (verified: no non-comment line in the diff). The gate table, the plan format and the transition table did not move, `core.DomainVersion` is still **6**, the command count is still **42**, and `kitchen/activity.json` is untouched. `make lint`, `gofmt -l .`, `make test` and `make regen-goldens` are all clean.
