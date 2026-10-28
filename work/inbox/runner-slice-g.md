Runner: slice G fix 1 done. Commits: 4. Fixed: F-94, F-95, F-96 and the boundary. The tree is free of the drawings; the glyph is in the brand package, in the rendered output and in the card's art. Notes: the reviewer should know these.

| Commit | Fix |
|---|---|
| `a919d7f` fix: the mark yields to the content, in the height too | F-94 |
| `f39af35` fix: the wire that says where the mark may be, and that can fail | F-95 |
| `659e14a` fix: the boundary wire parses the imports, whatever the spelling | F-96 |
| `a6065b3` test: both boundaries of the width rule are pinned | the boundary |

**F-94 — the height rule.** The tab's area is `height − 5` rows; the mark costs
**six** (five lines and the blank under them); the content is 15. So the mark is
drawn only when the content and the mark both fit, and otherwise the tab has **no
mark** — the content is whole. The wire renders the tab at 80 columns and at 20,
24, 26, 30 and 40 rows, and asserts that the content is whole at every one of them
(the `LAST EVENT` section and the event row are on the screen) and that the mark is
drawn only at 26 and above. At 80×24 the tab is byte-for-byte what it was before
the mark.

**F-95 — the scope wire, and the truth about the tree.** Two rules, because there
are two things:

| Thing | Rule | Where it is, today |
|---|---|---|
| the drawings | in the source files and in the declaration, and **nowhere else** — any file type | the two `.txt` sources and `internal/brand/mark.go` |
| the compact glyph | must not be **pasted into another code file** | `internal/brand/mark.go`, the two goldens (rendered output), the wallpaper SVG (ninety copies — the card's own art), and prose: the brand sheet, the ledger, the briefs |

**The earlier claim was false and is withdrawn**: *"the mark's characters appear in
no other file in the tree"* — the goldens carry the glyph on the header line, and
so do the wallpaper, the brand sheet and the ledger. What is true: **the drawings
appear nowhere else in the tree**; the glyph is declared in the brand package and
appears in the rendered output and in the card's art.

**The exemption is per declaration, not per package.** The declaration is found by
parsing the package, so a pasted copy is never mistaken for it and a rename does
not widen the exemption. The wire can fail: the reviewer's two mutations (a copy
in `internal/brand/zz_probe.go`, a copy in `internal/ui/testdata/zz_probe.txt`), the
glyph in another code file, and each of them fails it — proven in a scratch tree
and by hand in this tree.

**The directories the wire does not read**: `.git`, `bin`, `skills` and `work`. The
binaries carry the art as string literals, and the transcripts capture what the
console shows — art and all. Everything else is read: `docs`, `images`, `kitchen`,
every file, any type. If the reviewer wants `work/` scanned too, that is one line.

**F-96 — the boundary wire parses the imports.** Every spelling is caught —
single-line, grouped, blank — because a wire that catches one spelling waits for
the other. Proven by hand in this tree: a blank import and a single-line import of
`internal/brand` in `internal/core` both fail the wire, and both compile. The
scope is in the comment: the rule is about the code that **ships**, so the wire
reads the **non-test** files of `internal/core`.

**The boundary is pinned.** The full form is 28 columns and is drawn one column in,
so it needs **29**: 29 takes the full form, 28 does not. The half form needs
**15**: 15 takes it, 14 and 13 do not. Changing the threshold to
`brand.Width(Full)` now fails the suite — it passed before.

**The tree.** `make lint && gofmt -l . && make test && make regen-goldens &&
git status --short` is clean; the goldens did not move (the header did not move).
The budgets: `internal/core` 0, `internal/cli` 0, `internal/ui` 28,
`internal/brand` 0. The tripwires hold — the derived-views table, the gate table,
the plan format, the command count (42). `core.DomainVersion` is still 7.

**Not touched**: `docs/audit.md`, `docs/07-vocabulary.md`, `work/` (except this
receipt), `images/`, `README.md`, and the card's copy.
