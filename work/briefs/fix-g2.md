You are the **runner** — the same lane that wrote slice G's fix. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## The last fix for slice G — one finding and one note

The review came back with **F-97** and a LOW note that belongs with it. The
orchestrator ruled: **both are fixed here.** One commit.

### F-97 — the wire exempts the declaration **file**, not the declaration

> The one that still passes: a copy pasted into `internal/brand/mark.go` **itself**,
> next to the declaration — the wire exempts the declaration FILE, not the
> declaration, so the whole suite stays green.

**The rule is *one declaration*, not *one file*.** A second copy of the drawing in the
file that declares it is a **second declaration**, and it is the same defect the wire
exists to catch — just one file closer to home.

**The fix — the exemption is the declaration:**

| | |
|---|---|
| **The drawing** | may appear **only in the canonical declaration** — the one `var` that holds it, in the one file that holds it. Anywhere else in the tree, in **any** file type, is drift — **including in the declaring file itself**. |
| **Inside the declaring file** | the drawing appears **exactly once**, in the declaration. |

So the wire's question becomes: *is every occurrence of the drawing the canonical
declaration?* — and a second copy beside it is not.

**Prove it:** paste a copy into `internal/brand/mark.go` next to the declaration and
confirm the suite **fails**. Then paste one into every other file the wire reads and
confirm it fails there too (the reviewer did: a golden, a doc, the plan, `cmd/`,
`images/`).

### The note that rides with it — the glyph's scope

> the glyph rule is scoped to `.go` files, so the glyph pasted into
> `scripts/console-shot.py` passes — and the runner's own table says the glyph must not
> be "pasted into another code file", and a `.py` file **is** a code file.

**The rule as stated and the rule as implemented must be the same rule.** The glyph's
scope is **code files** — `.go` *and* `.py` — and not "Go files". The glyph may be:

- **declared** in the brand package;
- **rendered** — the goldens carry it on the header line, and that is output;
- **drawn** — the wallpaper's ninety copies are the card's own art.

and it may **not** be pasted into another **code** file.

**So the wire needs a *category*, not an extension.** Three categories — **code**,
**rendered output**, and **art** — and each form's rule says where it may live:

| | may be pasted into |
|---|---|
| **the drawings** | **nothing** but the canonical declaration |
| **the glyph** | the canonical declaration, and nowhere else in **code** — but the rendered output and the art may carry it |

**And the exclusions stand.** The reviewer agreed: `.git` is the repository's objects,
`bin` is build output, `skills` is vendored third-party, and `work` is the night's
scratch — **the transcripts quote what the console shows, art and all**, so `work`
*must* be excluded. Say that in the wire's comment, so the next reader knows each
exclusion is a *record* and not a hole.

## The invariants

- The domain layer has **zero UI imports** — and the wire proves it, whatever the
  spelling.
- One fact, one answer: the art is declared once, in one package, and every consumer
  reads it from there.
- `core.DomainVersion` is **7**. Do not move it.
- The gate table, the plan format and the command count (**42**) must not move.

## The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. Leave anything you did not
write alone. If a file you did not write has changed, **stop and report it**.

## Not yours — do not touch

| | |
|---|---|
| `docs/audit.md` | the reviewer's to write. It files F-97 and closes it. |
| `docs/07-vocabulary.md`, `work/`, `images/`, `README.md` | the orchestrator's. |
| The card's copy | it stays in the card. |

## The gate

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

Green.

## One commit

```
Per-AC evidence:
- F-97: a copy of the drawing beside the declaration fails the wire; the glyph's scope covers the code files, the output and the art; go test ./internal/brand -v exit 0
```

**One write per file, and commit as soon as the file is coherent.**

## Finish with the receipt — update the one you wrote when you started

```bash
cat > work/inbox/runner-slice-g.md <<EOF
Runner: slice G fix 2 done. Commits: <n>. Fixed: F-97 and the glyph's scope. The drawing is declared once, and the wire holds it. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-g fix 2"
herdr agent prompt "w1Z:p1" "Runner: slice G fix 2 done. Commits: <n>. Fixed: F-97 and the glyph's scope. The drawing is declared once, and the wire holds it. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.**
