Runner: slice F fix 1b done. Commits: 4 (2 work, 2 receipts). Fixed: F-90 and the notes.

F-90 — the "one format" wire (`internal/core/retired_test.go`):
- The scan reads the **whole tree**: every package the repo has (`internal`, `cmd`,
  `scripts`), not one. A token in each of them trips it, and the test proves it —
  one case per package, plus an invented name in the `legacy` family.
- The exemption is **per place, not per file**: the places are declarations
  (`retired.go`'s `retiredDirectories`, `retiredInRecipe`, `retiredInBacklog`; and
  `model.go`'s `RetiredNames`), so a retired name cannot hide in `model.go` behind
  the exemption. The four probes the reviewer ran — `specials` in `internal/ui`,
  `specials` in `internal/core/model.go`, `"R01"` in `internal/core/store.go`,
  `LegacyThing` in `model.go` — all fail the wire now.
- The token list is complete: the retired id literals (`R01`, `S01`, `N01`), the
  retired families (`legacy`, `special`, `prep-note`, `proof`), the retired crew
  names, and the retired functions.
- The recipe's retired stamp is caught by its **access**, not by its word: the token
  is `["evidence"]` — the literal key the retired stamp is read by — because
  `evidence` is a current word (the criterion's record).

Notes:
- **The file name's id form** — a file named `R01-<slug>.json` is a retired id form
  on disk, and the plan is not supported (the reader and the writer both refuse it,
  and the refusal names the file and the current form). The rule stops at the
  **retired forms**, because a name is not a record: the lookup takes `RCP1`,
  `rcp01` and `1` for the same artifact, and the accepted tests rely on
  `B01-<slug>.json` and `RCP1-<slug>.json` being read. A *canonical* file-name rule
  would have contradicted them.
- **Not supported means no writes** — `Store.checkWritable` is the write path's
  guard now: the version fence **and** the retired fence. The declaration
  (`RecordMigration`) still crosses the version fence, but not the retired one.
- **`GLOSSARY.md`** — the Sous-chef's prose was already fixed by F-87; the
  glossary's rule (a retired word is named **once**, in the `_Avoid_` line) now has
  a wire, over the crew's names and the plan's retired file types.
- **`activity.go`** — the version-zero comment was already fixed by F-89: only a
  plan with **no** journal is writable.

Two tests changed their premise; the reviewer should look at both:
- `TestAFileNamedWithTheOldIdIsStillRead` → `…IsRefused`: the file's name is part
  of the format. The **lookup's** tolerance stands — `R01` still names the recipe.
- `TestTheRemovalLeavesAnotherKindsFileAlone`: the stray file's record is now a
  recipe's, because a record whose id is not a recipe's id is a plan the tool does
  not support — and the writer now agrees with the reader about that. The removal is
  still by the whole id, which is what the test is about.

Untouched: the gate table, the plan format, `core.DomainVersion` (7), the command
count (42), the journal, `docs/audit.md`, `docs/07-vocabulary.md`, `work/`.
The invariants hold: `internal/core` 0 · `internal/cli` 0 · `internal/ui` 28.

Gate: make lint && gofmt -l . && make test && make regen-goldens — green.
