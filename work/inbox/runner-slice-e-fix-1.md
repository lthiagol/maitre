Runner: slice E fix 1 done. Commits: 5. Fixed: F-81..F-85. Notes: the order is F-83, F-81, F-82, F-84, F-85 — the glossary's proof command counts retired words outside their `_Avoid_` lines, and F-83's prose line was the one that made the count non-zero, so it came first. The mascot did not land. `docs/03`'s `[the version guard](#the-version-guard)` link has no heading to point at — pre-existing, not one of the five findings, left alone.

## The commits

| Commit | Finding | What |
|---|---|---|
| `28c7742` | **F-83** | the backlog item's prose stops naming the retired words; the `_Avoid_` line keeps them |
| `2930fc3` | **F-81** | the console's retired name is carried in its `_Avoid_` line, and the proof reproduces |
| `7d5df0a` | **F-82** | the console leaves the crew for **### The strings**, and its header is spelled `maitre` |
| `4a47494` | **F-84** | the two lists of the fingerprint's facts agree — and a third list, in the code, now agrees too |
| `91e06a6` | **F-85** | the version fences the writer; the fingerprint is the tripwire |

## F-81, F-82, F-83 — the glossary

**The retired-word pattern now holds for every renamed term.** The proof, and it
reproduces:

    $ grep -n -i -E "line cook|taster|expo|proof|special|prep note|maître d'" GLOSSARY.md
    63:_Avoid_: Special, prep note, note, brief, ticket
    103:_Avoid_: Line cook, executor, implementer
    108:_Avoid_: Taster, critic, judge
    112:_Avoid_: Expo, dispatcher, supervisor
    130:_Avoid_: Proof, test result, verification
    161:_Avoid_: maître d', TUI, dashboard, frontend, UI

    $ grep -n -i -E "line cook|taster|expo|proof|special|prep note|maître d'" GLOSSARY.md | grep -v "_Avoid_:" | wc -l
    0

Six retired names on six lines, plus the console's own on the seventh — every one of
them the line that warns the reader off it. (The recipe-level stamp's old name,
`Evidence`, is the one retired name that is also a current word: the term is the
entry, and the warning is in **completion**'s `_Avoid_` line.)

The AC-02 commit's proof returned 0 for all four names it claimed to check — three
because the file capitalises them and one because the word was absent — so it proved
nothing. That command is replaced by one that works and prints its output.

**The console is not one of the crew.** The structure now reads:

    $ grep -n "^### \|^\*\*.*\*\*:$" GLOSSARY.md
    43:### The plan
    45:**Kitchen**:
    50:**Recipe**:
    56:**Backlog item**:
    65:**Course**:
    69:**Menu**:
    74:**Line**:
    79:**Board**:
    84:**Service**:
    89:### The crew
    91:**Chef**:
    96:**Sous-chef**:
    100:**Runner**:
    105:**Reviewer**:
    110:**Orchestrator**:
    114:### The work
    151:### The strings
    158:**maitre**:

Five entries under **The crew** — the five jobs the scope table names — and the
console under **The strings**, spelled the way the code spells it. The section says
what it is: a surface, not a job, per ADR-0003 and ADR 0009 § 4.

**The prose names no retired word.** The backlog item's entry now points at ADR 0009
§ 3 for the two file types it replaces, and the `_Avoid_` line keeps the names.

## F-84, F-85 — the spec and the code

**Three lists, one answer.** `docs/03` listed the fingerprint's facts twice and the
two disagreed; the `DomainVersion` comment listed them a third time, without the root
files. All three now read "the id prefixes, the directories, the root files and the
field names on disk" — which is what `internal/core/format.go:20` hashes.

**The coupling is gone from the prose.** Five sentences taught it:

| Where | Was |
|---|---|
| `docs/03` § 6 | "A format change moves the number: …" |
| `docs/03` § 6.1 | "**What the stamp covers is exactly what the fingerprint hashes**" |
| `docs/03` § 6.1 (version 4) | "a format change moves the version on its own instead of by hand" |
| `internal/core/lifecycle.go` | "What the stamp covers is exactly what `domainShape()` hashes" |
| `internal/core/lifecycle.go` (version 4) | "so a format change moves the version on its own" |

They now say what the code says: **the version fences the writer; the fingerprint is
the tripwire that says the facts moved.** A change to the *facts list* moves the
fingerprint and leaves the version alone; a change to *what the tool writes* moves
both. The version-4 history keeps its facts and loses the belief.

Nothing moved but the prose: `core.DomainVersion` is still **6**, the gate table, the
plan format and the command count are untouched, and
`TestDomainShapeIsVersioned` passes with the same fingerprint.

## The gate

    $ make lint && gofmt -l . && make test && git status --short
    go vet ./...
    go test ./...
    ok  	github.com/lthiagol/maitre/internal/cli	2.454s
    ok  	github.com/lthiagol/maitre/internal/core	(cached)
    ok  	github.com/lthiagol/maitre/internal/ui	(cached)
    ?? logo/

`gofmt -l .` printed nothing. The only thing in `git status` is `logo/` — the
mascot's art, untracked and not mine.

## Two notes

- **The sweep's count moved from 78 to 77.** The command in the AC-04 commit body
  does not include *maître d'*, and F-83 removed the one prose line that carried two
  retired words. The pattern now holds across the glossary: seven retired names, seven
  `_Avoid_` lines, nothing in prose.
- **`docs/03` § 6 links to `[the version guard](#the-version-guard)`, and there is no
  such heading** — the anchor is dangling. It predates this slice and it is not one
  of the five findings, so I left it. It is a one-line fix if the orchestrator wants
  it: the guard is described at the end of § 6.1.

## The mascot

**Still not landed.** `logo/` holds the art, untracked; the console's header, the
CLI's help, the README and the Sous-Chef tab carry no drawing. I touched none of
them, and the goldens did not move — no console output changed in this cycle.
