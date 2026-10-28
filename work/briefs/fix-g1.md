You are the **runner** — the same lane that wrote slice G. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## The fix cycle for slice G — three findings

The reviewer filed **F-94**, **F-95** and **F-96**. The orchestrator ruled on all
three: **accepted**. Three findings, three commits — one each.

**The reviewer's mutations are the reading of this brief.** Two of the three findings
are wires that **cannot fail**, and the third is a layout that costs the chef a line
of content. Read the mutation list in the review's report before you start.

### F-94 (medium) — the art pushed the last event off the screen

> The art's six rows push the Sous-Chef tab's last event off the bottom of any
> terminal under 26 rows, and the tab does not scroll. At 80×24 — the classic default
> — the baseline showed `LAST EVENT` and the event row; the slice shows neither.

**The rule, which is the one this slice needed all along: the decoration yields to
the content, in both dimensions.**

- The **width** rule already says it: the column chooses the form, and no art when
  neither fits.
- The **height** must say it too: the tab knows its area (`listAreaHeight` is
  `height − 5`), so **the art is drawn only when the tab's content *and* the art fit.**
  When they do not, **no art** — and the content is whole.

**The measure that matters is the chef's screen.** At **80×24**, the `LAST EVENT`
section and the event row must be there, exactly as the baseline had them. A face is
worthless if it costs the chef the last thing that happened.

**A wire, at several heights.** Render the tab at a spread of heights — 20, 24, 26,
30, 40 — and assert that the **content is complete at every one of them** (the last
event is present), and that the art appears only at the heights where it fits. The
wire is the *content*; the art is the decoration.

### F-95 (medium) — the wire that says "the art is declared only here" cannot fail

> The wire reads only `.go` files, and it **skips `internal/brand`** — the one place
> a second declaration would land. Two mutations survive the whole suite: a second
> copy in `internal/brand/zz_probe.go`, and a second copy in
> `internal/ui/testdata/zz_probe.txt`.

**Two different things are in play, and the rule for each is different.**

| | |
|---|---|
| **The drawings** — the full and half forms | long, unique, multi-line. They must appear **nowhere else in the tree**, in **any file type**, and the exemption must be **per declaration**, not per package. A per-package skip is where the second copy hides. |
| **The compact glyph** `⋈` | a single character, and it **legitimately** appears elsewhere: the **goldens** carry it on the header line (that is rendered output), and the **wallpaper SVG** carries 90 copies (that is the card's own drawing). So the glyph's rule is narrower: it may be *declared* in the brand package, and it may appear in *rendered output* and in *art* — but it must not be **pasted into another code file**. |

**And the claim in the receipt is false as written** — *"the mark's characters appear in
no other file in the tree"* — because the goldens carry the glyph on the header line.
**Say what is true in your report**: the *drawings* appear nowhere else; the glyph
appears in the brand package, in the rendered output, and in the card's art.

**Prove the wire can fail** on each of the reviewer's two mutations, and on a copy
pasted into a **non-`.go`** file.

### F-96 (medium) — the wire that keeps the mark out of the domain is defeated by the import spelling `gofmt` writes

> `oursIn` reads the raw source line and matches a line that *starts* with the module
> path. A file whose only import is `import "github.com/lthiagol/maitre/internal/brand"`,
> or a grouped `import ( _ "…/internal/brand" )`, compiles, is gofmt-clean, and puts
> the brand inside `internal/core` — and `TestTheMarkIsNotTheDomain` passes.

**The fix is the one the reviewer names: parse the imports.** The repo already has
`go/parser` in `internal/core/retired_test.go` for exactly this. **Any** spelling —
single-line, grouped, blank (`_ ""`) — must be caught, because a wire that only
catches one spelling of an import is a wire that waits for the other spelling.

**And the scope must be stated:** the rule is about the code that **ships**, so the
wire checks the **non-test** files of `internal/core` — and says so in its comment,
because that sentence is what the reviewer found missing.

### The rulings that ride with them

| Note | Ruling |
|---|---|
| the **full form's threshold is unpinned** — changing `brand.Width(Full)+1` to `brand.Width(Full)` passes every test, so a 28-column column would take the 29-column drawing | **fix it: pin both boundaries.** A column of exactly **29** takes the full form; **28** does not (the indent makes it 29 columns); **15** takes the half; **14** does not; **13** does not. The half's boundary is already pinned — the full's is not. |
| the brand sheet says *"28 columns wide"*, but the drawing's **first row is 26** | **the orchestrator's file** — the sheet is being corrected to say *at its widest*. Nothing for you. |
| `⋈` is a math operator, not a bow tie | **no action.** A reader who has not seen the card may read it as a tie-knot symbol — which is what it is meant to be. |
| the wallpaper SVG carries 90 copies of the glyph | **folded into F-95** — the glyph's scope is narrower than the drawings', and the report must say so. |

## The invariants

- The domain layer has **zero UI imports** — and now the wire proves it, whatever the
  spelling.
- `internal/core` must not import `internal/brand`, and `internal/brand` must not
  import `internal/core`.
- One fact, one answer: the art is declared once, in one package, and every consumer
  reads it from there.
- `core.DomainVersion` is **7** — the art is not a plan format fact. Do not move it.
- The gate table, the plan format and the command count (**42**) must not move.

## The tree is shared

**Another session may be working in this repository.** Stage **by name** — never
`git add -A`, never `git add .`. Leave anything you did not write alone. If a file
you did not write has changed, **stop and report it**.

## Not yours — do not touch

| | |
|---|---|
| `docs/audit.md` | the reviewer's to write. It files F-94..F-96 and closes them. |
| `docs/07-vocabulary.md`, `work/`, `images/` | the orchestrator's. The brand sheet's line is being corrected there. |
| The **card's copy** | the title, the tagline, the stamp, the contact line stay in the card. The terminal has the emblem. |
| The names | settled. None of them moved. |

## The gate

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

Green.

## Three commits — one per finding

```
Per-AC evidence:
- F-94: the art never costs content a line; the tab is whole at 80x24; go test ./internal/ui -run 'SousChef|Mark' -v exit 0
```

**One write per file, and commit as soon as the file is coherent.**

## Finish with the receipt — update the one you wrote when you started

```bash
cat > work/inbox/runner-slice-g.md <<EOF
Runner: slice G fix 1 done. Commits: <n>. Fixed: F-94, F-95, F-96 and the boundary. The tree is free of the drawings; the glyph is in the brand package, in the rendered output and in the card's art. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-g fix 1"
herdr agent prompt "w1Z:p1" "Runner: slice G fix 1 done. Commits: <n>. Fixed: F-94, F-95, F-96 and the boundary. The tree is free of the drawings; the glyph is in the brand package, in the rendered output and in the card's art. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.**
