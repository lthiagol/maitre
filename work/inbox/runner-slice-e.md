Runner: slice E done. Commits: 6. Per-AC evidence: AC-02..AC-05 stamped. Pre-notify gate: pass. Mascot: no. Notes: the mascot did not land, the domain version stayed 6, and the old names that remain all name something that still exists.

## The commits

| Commit | What |
|---|---|
| `21a03d7` | AC-02 — the glossary is rewritten |
| `1280149` | AC-03 — ADR 0009, the index, 0004 and 0005 superseded, ADR-0006 § 2 corrected |
| `c1c1e7c` | AC-04 — the sweep |
| `5e29050` | the two LOW notes from slice A's review |
| `7ec436d` | the glossary's mapping covers every term it names |
| `53b8f5d` | AC-05 — the whole tree, verified (no file changed; the commit carries the evidence) |

## AC-01 — the mascot: **not landed**

The art is in `logo/` — untracked, and not mine. **None** of the mascot's placement
points carries a drawing: the console's header still renders ` maitre · <name>`, the
CLI's help still opens with `maitre-cli — the maitre agent CLI`, the README still
starts with `# maitre`, the Sous-Chef tab has no art, and the goldens moved with the
sweep alone (`make regen-goldens` — regenerated, never hand-edited).

I touched none of those files, and I did not commit `logo/`. AC-01 moves to the
session that owns the drawing.

## The two judgements the reviewer should check

**1. `core.DomainVersion` stayed 6 — on purpose.** The LOW note asked for two kinds
of fact in `planFormatFacts`, and the root file names and the journal envelope's
fields are now there. That changed the fingerprint (`cff4269d589690ce` →
`b84b230161f4ab03`), so the recorded shape moved — **and the version did not**, because
the version fences the **writer** and the writer's output did not change (the ledger's
own F-80 says exactly this). The code now says the distinction out loud: the
`recordedDomainShape` comment and the failure message say that a change in *what the
tool writes* moves `DomainVersion`, and that the shape alone fences nothing. If the
reviewer disagrees, the fix is a bump to 7 plus `docs/03` § 6 — say so and it is a
two-line change.

**2. The glossary names each retired word once, on purpose.** `Line cook`, `Taster`,
`Expo`, `Special` and `Prep note` survive in the `_Avoid_` lines, and the *backlog
item* entry says in prose what it replaces. That is the pattern slice D set with
`Proof` → `evidence`: the retired name goes into the `_Avoid_` line, which is where a
reader is told not to use it. The alternative — dropping them from the `_Avoid_` lines
— would leave the glossary unable to say what not to say.

## The sweep: what remains, and why

The proof command and its output are in `c1c1e7c`. It returns **78 lines in 19
files**, and every one of them is a name that still exists:

| Hits | File | Why it stays |
|--:|---|---|
| 15 | `internal/core/legacy.go` | the legacy reader: `LegacySpecial`, `LegacyPrepNote`, the two old directories |
| 11 | `internal/core/roles_test.go` | the tests that hold the retired names on the deny-list |
| 10 | `internal/core/model.go` | `retiredAgentRoles` and the digest's old field name |
| 9 | `docs/03-domain-model.md` | the aliases the domain accepts, `kitchen/specials/`, `kitchen/prep-notes/`, the version-5 writer's `proofs` |
| 6 | `GLOSSARY.md` | the `_Avoid_` lines and the sentence saying what the backlog item replaces |
| 5 | `internal/core/lifecycle.go` | the version-4 and version-5 writers' directories and field names |
| 4 | `internal/ui/console_test.go` | the fixture of an old plan |
| 3 | `internal/core/legacy_test.go` | the same, in the domain |
| 2 each | `internal/core/views.go`, `backlog_test.go`, `backlog_concurrency_test.go`, `scripts/stamp-plan/main.go` | the alias comparison, the migration fixtures, and the names `stamp-plan` records in a journal |
| 1 each | `internal/core/{plan_ids,migration,gates,docs,approval}_test.go`, `docs/04-operations.md`, `prototypes/README.md` | the migration's own ids, the tripwire's mutation, the alias statement, and the note that the captures predate the rename |

`Kitchen Brigade`, `scale`, `Chopboard`, `line-cook` and `R01`-style ids are **zero**
outside that list.

**Not swept, and why each is a record and not drift**

| Excluded | Why |
|---|---|
| `docs/audit.md` | the ledger of the night — its rows keep the names as they were when the finding was filed |
| `docs/audits/` | the per-cycle reports behind the ledger's rows — the same record |
| `research/00-legacy-analysis.md` | v1's vocabulary is history |
| `docs/decisions/0004-naming.md`, `0005-vocabulary.md` | the ADRs whose subject is the vocabulary keep their tables as the record of what they said |
| `docs/decisions/0009-vocabulary.md`, `docs/07-vocabulary.md` | a rename table names the old names; that is what it is for |
| `work/` | the briefs, the transcripts, the receipts — and slice E's own brief, which lists every name to sweep |
| `prototypes/*.txt` | captured console output from before the rename. The README now says so |
| `kitchen/` | the repo's own plan. The journal is append-only; the historical events keep their old ids and their old field names, and `RCP01`'s records keep the runner that wrote them. `kitchen/activity.json` is **byte-identical** |
| `skills/` | the vendored skills are third-party and pinned. **Ours** — `skills/README.md` and `skills/sync.sh` — were swept: `README.md` still said *Kitchen Brigade* and *our `taster` role* (`sync.sh` already said *maitre*) |

## The other changes, in one line each

- **ADR 0006 § 2** now says what the code has: four match states — `proven`,
  `rejected`, `unknown`, `signed-off` — and the correction is dated in the ADR. The
  sentence about `ac.all_proven` says what the gate reads: the criterion's status,
  not the match state.
- **`internal/core/verifier.go`**'s `MatchState` comment had the same misdescription
  the ADR had — three states, one of them named wrong, the sign-off missing. Fixed.
- **ADR 0003**'s `line-cook` → `runner`: 0003's subject is the mechanism, so its
  incidental vocabulary moves. 0004 and 0005 keep every table they had.
- **The seed's demo recipe** said "Specials and prep notes"; the demo is read by the
  chef, so it speaks the current words.
- **The usage-lines wire** now reads every package in the repository and checks each
  usage line against the commands the repository actually ships (read from the tree).
  The proof — a usage line in `internal/ui` naming `scale` — fails the wire, and the
  output is in `5e29050`.
- **The format's facts** now name the root files and the journal envelope's fields,
  and two tests hold it: the facts name them, and **every** fact in the list moves
  the fingerprint.

## Two files I did not touch, on purpose

- **`README.md`** — the mascot's file. Its *Decisions* table lists 0001–0005 and has
  no row for 0009. The orchestrator may want to add one; I did not, because that
  file is the mascot's.
- **`work/README.md`** — the orchestrator's file. It still says the format version is
  **4** (it is 6) and calls the reviewer *the taster* and the writer *the line cook*.
  It also says slice E is *queued*. Those are the orchestrator's to close.

## The invariants

- `internal/core` imports no UI package. Stdlib only.
- `core.DomainVersion` is **6**; `docs/03` § 6 says 6.
- The journal was not touched.
- The kitchen theme stays; the artifact is still a recipe.
- The dependency budget: `internal/core` **0**, `internal/cli` **0**, `internal/ui`
  **28** external modules.
- The doc tripwires: the derived-views table, the gate table, the plan format and the
  command count (**42**) — each run on its own, each green.
- `make lint` clean, `gofmt -l .` empty, `make test` green, `make regen-goldens`
  green, `git status --short` shows only `?? logo/`.
