Runner: slice G done. Commits: 5. Per-AC evidence: AC-01..AC-05 stamped. Pre-notify gate: pass. The tab's width rule: the column chooses the form — the full form when the column holds its 28 columns plus the one-space indent (29), the half form when it holds 14 plus the indent (15), and no art when it holds neither. Notes: the reviewer should know these.

| Commit | AC |
|---|---|
| `045411e` feat: the mark, declared once, in internal/brand | AC-01 |
| `ff151b0` feat: the mark, in the header, in the help and in the tab | AC-02 |
| `657d4f0` test: the wire that holds the code to the drawing | AC-03 |
| `d3c6a22` test: the layout the mark must not break | AC-04 |
| `d47be09` test: the whole tree, verified | AC-05 |

**The width rule, and why that threshold.** The rule is the **column**, not the
terminal: `m.listColumn()`. The drawing is drawn one column in, like the rest of
the tab, so the column must hold the drawing **plus that one space** — 29 for the
full form, 15 for the half. A drawing that does not fit wraps, and a wrapped row
breaks the layout, so the column picks the form and the tab is never broken by the
mark.

| Column | Form | Terminal width (below the split) |
|--:|---|---|
| ≥ 29 | the full form | ≥ 31 |
| 15–28 | the half form | 17–30 |
| < 15 | no art | unreachable today |

**The third case is unreachable through the console today**: `listColumn()` floors
at 20 (`rowWidth()`'s floor), and the half form needs 15. The branch is there
because the rule says it, and the test checks the rule where it lives
(`markFor`), not through the model.

**The comparison runs from the file to the code.** The drawing in
`images/logo/ascii/` is the source, so the file is what the code must match: the
wire reads the file and compares it to the declaration. Changed a character in the
`.txt` and the wire fails and names the line — verified by hand during the slice,
and the file was restored. The wire covers **both** drawings, and it can fail: the
drawing is edited in memory and the wire says no.

**Two more wires hold the mark to one place**: the half form is a second drawing,
not a crop — no window cut out of the full form is the half form, and the wire
walks every window; and the mark's own characters (the braille block and `⋈`)
appear in no other file in the tree, so the art cannot be pasted a second time.
The boundary holds too: the mark imports nothing of ours, and the domain does not
import the mark.

**The tree.** `make lint && gofmt -l . && make test && make regen-goldens &&
git status --short` is clean; the goldens' diff is the header's one glyph and
nothing else. The budgets: `internal/core` 0, `internal/cli` 0, `internal/ui` 28,
`internal/brand` 0. The tripwires hold — the derived-views table, the gate table,
the plan format, the command count (42). `core.DomainVersion` is still 7.

**The card's copy was not installed in the terminal.** The slots are the ones in
the brief and nothing else: the header carries the compact glyph, the help carries
the full form above the title, the tab carries the drawing above the section. The
card's title, tagline, stamp and contact line stayed in the card — the card is an
image. If the orchestrator wants the tagline in the console header, that is a new
slot and it is the orchestrator's call.

**Nothing outside the code was touched**: `README.md`, `images/`, `docs/`,
`GLOSSARY.md` and `AGENTS.md` are as they were.
