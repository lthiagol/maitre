You are the **runner** for slice G — the last piece of the rename. The orchestrator's
pane id is `w1Z:p1` — it is already substituted in this file; wherever this brief says
`$ORCH_PANE`, it means `w1Z:p1`. Use it verbatim in the notify command below. The
repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` (house style and
the working rules) and `docs/07-vocabulary.md` (the decision of record) first.

## The mark is decided — the chef chose it

A **bow tie**, drawn as a halftone: two wings with folded-ribbon highlights and a
solid square knot. The terminal draws it with **braille** characters (2 dots wide, 4
tall); the card draws the same dots as circles. One drawing, two media.

**The art is in the repository, and it is the source:**

| File | What |
|---|---|
| `images/logo/ascii/bow-tie-full.txt` | **28 × 5** — the terminal's full form |
| `images/logo/ascii/bow-tie-half.txt` | **14 × 5** — the terminal's mid-size form |
| `images/logo/9-plan-behind.svg` / `.png` | the **banner**, for the README |
| `images/logo/archived/` | the other cards, kept as the reasoning behind the banner |

`images/logo/README.md` is the brand sheet: the copy, the slots, and the note that
*the art survives being halved — at 14 columns it is still unmistakably a bow tie*.

**The orchestrator has already done the README's half** (the banner and the first
line). **Yours is the code.**

## Your slice — G — the mark in the code

### AC-01 — the drawing lives in **one place** in the code

**A new, tiny package: `internal/brand`.** It holds the art, and **nothing else**.

**Why a package and not `internal/core`:** `internal/cli` and `internal/ui` both
import only `internal/core` today, and neither may import the other. The art is
needed by both. `internal/core` is the **domain** — the house rule is that the
domain has no UI in it, and a bow tie is UI. So the art goes in its own package, and
the domain stays clean.

```
internal/brand/mark.go    the full form, the half form, the compact glyph
```

| | |
|---|---|
| **The full form** | the 5 lines of `bow-tie-full.txt`, 28 columns |
| **The half form** | the 5 lines of `bow-tie-half.txt`, 14 columns — **a second drawing of the same emblem at a smaller size**, not a crop of the first, and it is declared once |
| **The compact glyph** | `⋈` — one cell, the console's header |

**Do not paste the art twice.** Every form is declared **once**, in this package, and
every consumer reads it from here. A wire holds the code to the art — see AC-03.

### AC-02 — where each form goes

| Slot | Form | Rule |
|---|---|---|
| **The console's header** | the **compact** glyph, before the product name | the header stays **one line** — that is the whole reason the compact exists |
| **The CLI's help header** | the **full** form, above the title line | the help is text; the art goes above `maitre-cli — the maitre agent CLI` |
| **The Sous-Chef tab** | the **full** form, or the **half** when the column is too narrow | the tab's content is a list of lines; the art goes at the top, above the `SOUS-CHEF` section |
| **The README** | the banner | **done by the orchestrator** — do not touch `README.md` |

**The tab chooses its form by the width it has.** The list column is
`m.listColumn()` — 45% of the width at a wide terminal, the row width otherwise. A
28-column drawing in a 20-column column would **wrap**, and a wrapped row breaks the
layout. So:

- the column fits 28 → the **full** form;
- the column fits 14 but not 28 → the **half** form;
- the column fits neither → **no art**, and the tab is as it is now.

**Say which rule you implemented, and why that threshold.**

### AC-03 — the art's source, and a wire that holds the code to it

`images/logo/ascii/*.txt` is the **source**. The code holds it once, and a **wire**
asserts the two agree — the same shape as the goldens: the expected file, the test
compares.

**The wire must cover both drawings**, and it must be able to fail: change a character
in the `.txt` and the suite must catch that the code has not moved. Say which way the
comparison runs and why.

### AC-04 — the layout must not break

- The console's **header stays one line** — the goldens prove it.
- The art never exceeds the column it is in (AC-02's rule).
- The **Sous-Chef tab still says the automation is not built.** The tab that is
  waiting must not look like a tab that is failing, and a face does not change that.
- The tab's **vital signs** and its last event still read as they do.
- **`make regen-goldens` writes the goldens**, and the diff is only the art and the
  header.

### AC-05 — the whole tree, verified

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

- The dependency budget holds: `internal/core` and `internal/cli` stdlib-only;
  `internal/ui` at **28** external modules; the new package at **0**. Count them the
  way `AGENTS.md` says to count them and report the number.
- The doc tripwires hold: the derived-views table, the gate table, the plan format,
  the command count (**42**).

## The invariants

- The domain layer has **zero UI imports** — and the new package is **not** the
  domain: `internal/core` must not import it, and it must not import `internal/core`.
- One fact, one answer: the art is declared once, in one package, and every
  consumer reads it from there.
- `core.DomainVersion` is **7** — the art is not a plan format fact. Do not move it.
- The gate table, the plan format and the command count must not move.

## The scope — what not to touch

| | |
|---|---|
| **`README.md`** | **the orchestrator's** — the banner and the first line are done. Do not touch it. |
| **`images/`** | the art is the **source**, and the orchestrator commits it. Do not edit the `.txt` files or the cards. |
| **`docs/audit.md`** | the ledger of the night. **Never swept.** |
| **`docs/`, `GLOSSARY.md`, `AGENTS.md`** | the paper. **Do not touch.** The mascot is not in them. |
| **The mascot's placement** | the slots in AC-02, and nothing else. The card's copy — the title, the tagline, the stamp, the contact line — belongs to the **card**, which is an image. **Do not install the card's copy in the terminal.** If you think a slot wants it, say so in your report and let the orchestrator rule. |
| **The names** | settled by the rename. Do not undo one. |

## The discipline gate

**One commit per acceptance criterion**, with the evidence stamped in the body:

```
Per-AC evidence:
- AC-01: the art is declared once, in internal/brand; go test ./internal/brand -v exit 0
```

**One write per file, and commit as soon as the file is coherent** — a commit is a
landing place, and the diff render is a turn's budget.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-g.md <<EOF
Runner: slice G done. Commits: <n>. Per-AC evidence: AC-01..AC-05 stamped. Pre-notify gate: pass. The tab's width rule: <what you implemented>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-g"
herdr agent prompt "w1Z:p1" "Runner: slice G done. Commits: <n>. Per-AC evidence: AC-01..AC-05 stamped. Pre-notify gate: pass. The tab's width rule: <what you implemented>. Notes: <anything the reviewer should know>."
```

**Write the receipt, commit it, and *then* nudge.** The nudge only reaches an idle
orchestrator; the receipt is the half that cannot be lost.
