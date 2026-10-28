You are the **reviewer** — the same lane that reviewed slice G. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 3 — the last fix of slice G

The runner fixed F-97 and the glyph's scope. **Cycle 3 is blockers only**: file **HIGH
or MEDIUM** blockers. No nits, no cosmetics, no style.

### 1. F-97 — the exemption is the declaration, not the file

**The rule:** the drawing may appear **only in the canonical declaration** — the one
`var` that holds it — and nowhere else in the tree, in any file type. The source files
(`images/logo/ascii/*.txt`) are the drawing's **source**, which is a legitimate home;
the *declaration* is the code's one copy.

**By mutation, as always:**

- a copy pasted into `internal/brand/mark.go` **beside** the declaration must now
  **fail** the wire — that was the one that slipped through;
- a copy in a golden, in a doc, in the plan, in `cmd/`, in `images/`, and in a
  **non-`.go`** file must all fail;
- and the **legitimate** copies must still pass: the `.txt` source, the declaration
  itself, and — for the glyph — the rendered output and the art.

### 2. The glyph's scope is a category, not an extension

> code (go and py), rendered output, art, prose — a script is code, so the glyph in
> `scripts/console-shot.py` now fails.

**Check it:** the glyph pasted into a **`.py`** file must fail; the glyph in the
**goldens** (the header line) must pass; the glyph in the **wallpaper SVG** must pass;
the glyph in **prose** must pass (the brand sheet quotes it). And the four excluded
directories — `.git`, `bin`, `skills`, `work` — must still be excluded, each for the
reason the comment gives: they are **records** or build output, and `work` holds the
transcripts that quote what the console shows, art and all.

### 3. The slip — the runner reported it, and the tree must be clean

The runner wrote a probe over `scripts/console-shot.py` by mistake while proving the
wire, and restored it from git at once.

**Check it yourself:** `git status` clean, `git diff HEAD` empty, and
`scripts/console-shot.py` **byte-identical** to HEAD (`git hash-object` against
`git rev-parse HEAD:…`). A file restored from git is byte-identical only if the
restore worked — and the runner reported the slip instead of letting the next reader
find it.

### 4. How it looks — one last look

- `maitre --render --plain` — the header, **one line**, with the compact glyph.
- `maitre --render --all --plain` — the Sous-Chef tab: the mark above the section,
  **"Not built yet."** still there, the vital signs and the last event whole.
- At **80×24**: **no mark**, and the content whole.
- `maitre-cli` with no arguments — the full form above the title line.

### 5. Scope, and the ledger

- `core.DomainVersion` is **7** and did not move. The gate table, the plan format, the
  command count (**42**) and the names did not move.
- The **card's copy** stayed in the card.
- **The ledger** — `docs/audit.md` is yours. Move **F-97** to closed, with the file
  and the evidence that holds it, and **add the row for slice G** — the mark, its
  three forms, the rule that it yields to the content in both dimensions, and the
  wires that hold it. The ledger is not swept.

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
cat > work/inbox/reviewer-slice-g-cycle-3.md <<EOF
Reviewer: slice G cycle 3 done. Verdict: <ok|changes-needed>. F-97 closed: <yes|no>. Filed <N> new blockers (<F-98.. or none>). Mutation checks: <what you mutated and reverted>. The slip: <the tree's state>. How it looks: <what you saw>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-g cycle 3"
herdr agent prompt "w1Z:p1" "Reviewer: slice G cycle 3 done. Verdict: <ok|changes-needed>. F-97 closed: <yes|no>. Filed <N> new blockers (<F-98.. or none>). Mutation checks: <what you mutated and reverted>. The slip: <the tree's state>. How it looks: <what you saw>. LOW notes: <any>."
```
