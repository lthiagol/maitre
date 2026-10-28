You are the **reviewer** — the same lane that reviewed slice G. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 2 — the fixes to slice G

The runner fixed F-94, F-95, F-96 and the unpinned boundary. **Cycle 2 is blockers
only**: file **HIGH or MEDIUM** blockers. No nits, no cosmetics, no style.

**Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's), or your own LOW notes
from cycle 1 (the brand sheet's width — the orchestrator corrected the sheet; the
compact glyph's reading — no action wanted; the wallpaper's 90 copies of the glyph —
folded into F-95).

### 1. Are the findings closed? **By mutation, as before.**

**F-94 — the mark yields to the content, in the height too.**

- The tab's area is `height − 5`, the mark costs six rows, the content fifteen. When
  the content and the mark do not both fit, the tab has **no mark**.
- **Prove the chef's screen is whole:** at **80×24** — the classic default — the
  `LAST EVENT` section and the event row must be there, exactly as the baseline had
  them. And at 80×24 there must be **no mark**.
- **The wire renders at 20, 24, 26, 30 and 40 rows.** Check the **content is complete
  at every one of them**, and that the mark appears only at the heights where it fits.
- The mark must **never** be drawn wider than its column either — the width rule and
  the height rule together.

**F-95 — the wire that says where the mark may be, and that can fail.**

- The **drawings** (the full and half forms) must appear **nowhere else in the tree**,
  in **any file type**, and the exemption is **per declaration**, found by parsing.
- The **glyph** `⋈` may be declared in the brand package, and appear in the rendered
  output and in the card's art — but not be **pasted into another code file**.
- **Re-run your two mutations** — a copy in `internal/brand/zz_probe.go`, and a copy
  in `internal/ui/testdata/zz_probe.txt` — and **find the third**: a copy in a file
  the wire does not read. The runner says the wire reads everything except `.git`,
  `bin`, `skills` and `work`; **judge those four exclusions** — is each one a *record*
  or is one of them a hole? (`work/` holds the transcripts, which capture what the
  console shows — that is a record. `skills/` is vendored third-party. `bin/` is
  build output. `.git` is the repository.)
- **The withdrawn claim.** The runner withdrew *"the mark's characters appear in no
  other file in the tree"* — it was false, because the goldens, the wallpaper, the
  brand sheet and the ledger carry the glyph. Check the report now says something
  **true**: the *drawings* appear nowhere else; the glyph is in the brand package, in
  the rendered output and in the card's art.

**F-96 — the wire that keeps the mark out of the domain parses the imports.**

- **Re-run both mutations**: a single-line `import "…/internal/brand"` (the spelling
  `gofmt` writes) and a grouped blank import `_ ""`. Both must now fail.
- **Check the scope**: the rule is about the code that **ships**, so the wire checks
  the non-test files of `internal/core` — and the comment says so.

### 2. The boundary, pinned

**29 takes the full form; 28 does not. 15 takes the half; 14 and 13 do not.** Check
both boundaries, and **mutate each one** — the full form's threshold and the half
form's — and confirm the suite fails. A boundary that can move without a test failing
is a boundary that will move.

### 3. How it looks — the thing the slice is for

Render it and **look** at it:

- `maitre --render --plain` — the header: `⋈ maitre · <project>`, **one line**.
- `maitre --render --all --plain` — the Sous-Chef tab: the mark above the section,
  and **"Not built yet."** still there, with the vital signs and the last event.
- `maitre-cli` with no arguments — the full form above the title line.

**And at 80×24** — the mark is gone and the content is whole. Judge both: that the
tab looks right when the art is there, and that it does not look *broken* when the art
is not.

### 4. Scope, and the ledger

- `core.DomainVersion` is **7** and did not move. The gate table, the plan format, the
  command count (**42**) and the names did not move.
- The **card's copy** stays in the card: no title, tagline, stamp or contact line
  leaked into the terminal.
- `images/`, `README.md`, `docs/07` and `work/` are the orchestrator's.
- **The ledger** — `docs/audit.md` is yours. Move **F-94, F-95, F-96** to closed,
  with the file and the evidence that holds each, and **add the row for slice G** —
  the mark, its three forms, the rule that it yields to the content, and the wires
  that hold it. The ledger is not swept.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.**

### Hard rules

- Do **not** edit implementation code. Read-only for everything except
  `docs/audit.md`.
- Revert every mutation you make. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-g-cycle-2.md <<EOF
Reviewer: slice G cycle 2 done. Verdict: <ok|changes-needed>. F-94..F-96 closed: <yes|no>. Filed <N> new blockers (<F-97.. or none>). Mutation checks: <what you mutated and reverted>. How it looks at 80x24 and at 40: <what you saw>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-g cycle 2"
herdr agent prompt "w1Z:p1" "Reviewer: slice G cycle 2 done. Verdict: <ok|changes-needed>. F-94..F-96 closed: <yes|no>. Filed <N> new blockers (<F-97.. or none>). Mutation checks: <what you mutated and reverted>. How it looks at 80x24 and at 40: <what you saw>. LOW notes: <any>."
```
