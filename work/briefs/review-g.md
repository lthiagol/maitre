You are the **reviewer** — the independent reviewer lane, fresh for slice G. The
orchestrator's pane id is `w1Z:p1` — it is already substituted in this file; wherever
this brief says `$ORCH_PANE`, it means `w1Z:p1`. Use it verbatim in the notify command
below. The repository is `/Users/thiago/Code/Pessoal/maitre`.

## Your slice — the review of slice G

Slice G is **the mark**: a bow tie, drawn as a braille halftone, in the console's
header, the CLI's help and the Sous-Chef tab.

```bash
cd /Users/thiago/Code/Pessoal/maitre
git log --oneline -12
git diff <base>..HEAD --stat      # the base is the commit before slice G's first commit
```

**Cycle 1: blockers only.** File **HIGH or MEDIUM** blockers. No nits, no cosmetics,
no style, no LOW notes filed as rows. **Do not re-file F-56, F-59, F-60, F-61, P-8**
(the chef's).

### The art, and where it came from

| | |
|---|---|
| `images/logo/ascii/bow-tie-full.txt` | **28 × 5** — the full form, **the source** |
| `images/logo/ascii/bow-tie-half.txt` | **14 × 5** — the mid-size form, **a second drawing of the same emblem, not a crop** |
| `images/logo/9-plan-behind.svg` / `.png` | the **banner**, for the README |
| `images/logo/archived/` | the other cards — the reasoning behind the banner |
| `images/logo/README.md` | the brand sheet: the copy, the slots, the sizes |

**The chef chose the emblem.** The art is committed, and it is the **source** the code
is held to.

### Verify each acceptance criterion with a command or a mutation

**AC-01 — the drawing lives in one place.**

- `internal/brand` holds the art and **nothing else**. The **full** form, the **half**
  form and the **compact** glyph are each declared **once**.
- **The consumers are `internal/cli` and `internal/ui`, and `internal/core` does not
  import the brand.** The house rule is that the domain has no UI in it, and a bow tie
  is UI. Check the import graph, not the comments.
- **No art is pasted anywhere else.** The runner claims the mark's characters appear
  in no other file in the tree — **check it**, including the test files and the
  goldens. A second copy of a drawing in a test fixture is a second copy of the
  drawing.

**AC-02 — where each form goes.**

- The **console's header** carries the **compact** glyph before the product name, and
  the header stays **one line**.
- The **CLI's help header** carries the **full** form.
- The **Sous-Chef tab** carries the **full** form — or the **half** when the column is
  too narrow, and **nothing** when neither fits.
- **Render it and look at it.** `maitre --render --plain` and
  `maitre --render --all --plain`; `maitre-cli` with no arguments. A mark that is
  correct in the code and wrong on the screen is a wrong mark.

**The width rule — the sharp end.** The runner's rule: the column holds 28 + the
indent → full; 14 + the indent → half; neither → no art. **Check the boundary
values** (a column of exactly 29, 28, 15, 14, 13) and that a form is **never**
rendered wider than the column it is in — a 28-column drawing in a 20-column column
would **wrap**, and a wrapped row breaks the layout. The runner says the list column
floors at 20 today, so the no-art case is unreachable *through the console* — **check
that claim**, and that the rule is still tested at the level it is tested at.

**AC-03 — the wire that holds the code to the art.**

- The wire compares **from the file to the code** (the drawing is the source), it
  covers **both** drawings, and it **fails when a character changes**. **Prove it**:
  change a character in the `.txt` and confirm the suite fails. A wire that cannot
  fail is not a wire.
- The **half is its own drawing**, not a crop of the full — and there is a test that
  says so.

**AC-04 — the layout must not break.**

- The console's **header is one line** — the goldens prove it, and the diff is one
  line each.
- The **Sous-Chef tab still says the automation is not built** — *"Not built yet."*
  must be there, and the vital signs and the last event must read as they did.
- **A tab that is waiting must not look like a tab that is failing.** A face does not
  change that: judge whether the art makes the tab look *ready* rather than waiting.

**AC-05 — the whole tree, verified.**

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

- The dependency budget: `internal/core` and `internal/cli` stdlib-only; `internal/ui`
  at **28** external modules; the new package at **0**. **Count them yourself.**
- The doc tripwires hold: the derived-views table, the gate table, the plan format,
  the command count (**42**).

### Then: what did the slice introduce?

1. **One fact, one answer.** The art is declared once and every consumer reads it
   from one place. **Check there is no second copy and no second source** — and that
   the *brand sheet* in `images/logo/README.md` and the code agree about the sizes.
2. **The card's copy stayed in the card.** The title, the tagline, the stamp
   (*EVIDENCE REQUIRED*) and the contact line belong to the **image**. The terminal
   got the **emblem** only. **Check none of the card's copy leaked into the
   terminal.**
3. **Scope.** `README.md` and `images/` are the orchestrator's (committed already);
   `docs/`, `GLOSSARY.md` and `AGENTS.md` are the paper and must be untouched. The
   names are settled — none of them moved. `core.DomainVersion` is **7** and did not
   move: the art is not a plan format fact.
4. **The mascot and the archive.** `images/logo/archived/` is the brand's record and
   the orchestrator committed it. Nothing in it may be referenced by the code.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.**

### The ledger

`docs/audit.md` is **yours** to write. The runner did not touch it — confirm that with
`git diff`. Add any NEW blocker as a new row with an id continuing from **F-93**. Do
not re-file what is closed; the ledger is not swept. **Add the row for slice G** — the
mark, its three forms, and where each one goes. Commit the ledger.

### Hard rules

- Do **not** edit implementation code. Read-only for everything except
  `docs/audit.md`.
- Revert every mutation you make. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-g.md <<EOF
Reviewer: slice G review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-94.. or none>). Mutation checks: <what you mutated and reverted>. How it looks: <what you saw on the screen>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-g"
herdr agent prompt "w1Z:p1" "Reviewer: slice G review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-94.. or none>). Mutation checks: <what you mutated and reverted>. How it looks: <what you saw on the screen>. LOW notes: <any>."
```
