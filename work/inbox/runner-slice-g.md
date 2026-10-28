Runner: slice G fix 2 done. Commits: 1 (plus this receipt). Fixed: F-97 and the glyph's scope. The drawing is declared once, and the wire holds it. Notes: the reviewer should know these.

**F-97 — the exemption is the declaration, not the file.** The wire's question is
now *is every occurrence of the drawing the canonical declaration?* The declaration
is found by parsing the package, and the two literals that hold the drawings are
masked out of the search — so the full drawing may be in the literal that holds it,
the half drawing in the literal that holds it, and nowhere else: not in another
file, not in another corner of the declaring file.

**The glyph's scope is the category, not the extension.** The rule as stated and the
rule as implemented are the same rule now:

| Form | code | rendered output | art | prose |
|---|---|---|---|---|
| the drawings | the canonical declaration, and nowhere else | — | — | — |
| the compact glyph | the canonical declaration, and nowhere else | may ride along | may ride along | may be quoted |

A script is code, so the glyph pasted into `scripts/console-shot.py` now fails the
wire. Prose may quote the glyph — the brand sheet and the reviewer's ledger name
it, and quoting is not pasting.

**The exclusions are records, not holes**, and the comment says so: `.git` is the
repository's own objects, `bin` is build output (the binaries carry the art as
string literals), `skills` is vendored third-party, and `work` is the night's
scratch — the briefs, the transcripts and the receipts quote what the console
shows, art and all.

**Proven by hand in this tree**, each one failing the wire and then removed: a copy
beside the declaration in `internal/brand/mark.go`; a copy in a golden; a copy in a
doc; a copy in the plan (`kitchen/`); a copy in `cmd/`; a copy in `images/`; and
the glyph in a script. The wire's own test proves the same list in a scratch tree.

**One thing to report.** While proving the wire, a probe was written over
`scripts/console-shot.py` by mistake — the probe's path was the script's own path.
It was restored from git at once and is byte-identical to HEAD; `git status` shows
it untouched. No other file was written outside the ones I own.

**The tree.** `make lint && gofmt -l . && make test && make regen-goldens &&
git status --short` is clean; the goldens did not move. The budgets: `internal/core`
0, `internal/cli` 0, `internal/ui` 28, `internal/brand` 0. The tripwires hold — the
derived-views table, the gate table, the plan format, the command count (42).
`core.DomainVersion` is still 7.

**Not touched**: `docs/audit.md`, `docs/07-vocabulary.md`, `work/` (except this
receipt), `images/`, `README.md`, and the card's copy.
