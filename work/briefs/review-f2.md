You are the **reviewer** — the same lane that reviewed slice F. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 2 — the fixes to slice F

The findings you filed were fixed across **two lanes**: the first lane landed F-87,
F-88 and F-89, then its turn ended without reporting; a **fresh** lane took F-90 and
the notes. **Cycle 2 is blockers only**: file **HIGH or MEDIUM** blockers. No nits,
no cosmetics, no style.

**Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's), or your own LOW notes
from cycle 1 (the stale comment, the file name's id form, the write path, the
glossary's prose, the wire's token choice — all folded into this cycle's work).

### 1. Are the four findings actually closed?

**F-87 — a retired name cannot buy a gate.** **Prove it:** claim as `line-cook`,
review as `runner`, and confirm `review.reviewer_independent` refuses. Then check:

- a retired name is refused **wherever an actor names itself** — claim, review,
  approve — not only for the approval;
- the **resolution** is still gone: a retired name is not "the runner", it is a name
  the tool refuses;
- **the repository's own plan moved consistently**: `kitchen/recipes/RCP01-*.json`'s
  claim and its review rounds were the same agent, so they must carry the same name —
  otherwise the self-review the gate refuses would stop being refused. **Check they
  match, and that the recipe's own history still reads the same.**
- The **journal's history keeps its old names** — it is history.

**F-88 — the operations doc says what the code says.** Check it: no retired name is
described as *naming the role it stood for*.

**F-89 — the digest's encoding is a format fact, and the version moved.**

- The payload's shape is in the **facts list**, so a change to it moves the shape, and
  the shape moving **must** move the version. `core.DomainVersion` is **7**.
- **The proof that matters:** a plan stamped **6** must be refused by the **version**
  guard, with the honest message — *"written under version 6, this tool speaks 7"* —
  and **not** as tampering. And a plan stamped **7** must be written.
- **The tamper message means content and only content.** Prove it: change a
  criterion's description after the sign-off and confirm the digest refuses.
- The repo's own plan is stamped 7 and its stored digest re-stamped, with the journal
  recording it as an **append** — the 37 historical events untouched.

**F-90 — the one-format wire.**

- It scans the **whole tree** — `internal/`, `cmd/`, `scripts/` — and the exemptions
  are **per declaration**, not per file.
- **Its token list is complete**, and the token choice is **unambiguous**: the retired
  stamp is caught by its *access* (`["evidence"]`), not by the bare word.
- **Prove it can fail**, on your four probes: `specials` in `internal/ui`, `specials`
  in `model.go`, `R01` in `store.go`, and an invented `LegacyThing`. All four must
  fail the wire.

### 2. The notes that rode with it

- **The file's name is part of the format.** A file named `R01-<slug>.json` is refused,
  and the message says what the name should be. **Check the rule stops where it should:**
  the **lookup's** tolerance stands (`RCP1`, `rcp01`, `1`, `R01` all *name* the recipe
  the plan holds), and a file named `B01-<slug>.json` is the same backlog item as
  `BL01-<slug>.json`. A *name* is not a record; a *file's name* is.
- **Not supported means no writes.** `Store.checkWritable` carries the version fence
  **and** the retired fence. **Prove it:** a write to a plan holding a retired
  artifact must fail and change nothing.
- The glossary's rule has a wire; `activity.go`'s version-zero comment says the rule
  as it is now.

### 3. The two tests that changed their premise — judge both

The runner flagged them itself, which is good. **Judge the premises, not the words:**

| Test | What it became |
|---|---|
| `TestAFileNamedWithTheOldIdIsStillRead` | `…IsRefused` — the file's name is part of the format, and the **lookup's** tolerance stands |
| `TestTheRemovalLeavesAnotherKindsFileAlone` | the stray file's record is now a recipe's, because a record whose id is not a recipe's id is an unsupported plan; the removal is still by the whole id, which is what the test is about |

**For each: does the test still assert the property it was written for?** A renamed
test that quietly asserts something *else* is a lost property wearing a new name.

### 4. The question you left open in cycle 1 — now closed

You were asked to judge *"a lookup is tolerant but a record is not"*: the normaliser
strips any leading letters, so `R01` resolves to `RCP01`. The orchestrator's ruling
was that the **lookup** is ergonomic — the number is the identity, the directory
carries the kind — while the **record** is one format.

**The work now draws that line in a second place: the file's name.** A *name* typed by
a human is a lookup and is tolerated; a *file's name* on disk is the format and is not.
**Judge it as one coherent rule**, and say whether the sweep's own list — which names
`R01`-style ids as old names that must be gone — is now satisfied.

### 5. Scope, the tree and the ledger

- The **gate table**, the **plan format**, the **command count (42)** and the **role
  names** must not have moved. `core.DomainVersion` moved **6 → 7** and that is F-89.
- **The mascot** — another session's work, in `logo/`. Check the runner committed
  nothing of it.
- **The ledger** — `docs/audit.md` is yours. Move **F-87, F-88, F-89, F-90** to
  closed, with the file and the evidence that holds each, and **add a row for slice F
  itself**: the ruling that this is not brownfield, the tolerance machinery that came
  out, and the new wire that holds *"one format"*. The ledger is not swept: keep the
  findings' wording as they were filed.

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
- **Do not touch `logo/`** — it is another session's work.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-f-cycle-2.md <<EOF
Reviewer: slice F cycle 2 done. Verdict: <ok|changes-needed>. F-87..F-90 closed: <yes|no>. Filed <N> new blockers (<F-91.. or none>). Mutation checks: <what you mutated and reverted>. The two renamed tests: <what each asserts now>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-f cycle 2"
herdr agent prompt "w1Z:p1" "Reviewer: slice F cycle 2 done. Verdict: <ok|changes-needed>. F-87..F-90 closed: <yes|no>. Filed <N> new blockers (<F-91.. or none>). Mutation checks: <what you mutated and reverted>. The two renamed tests: <what each asserts now>. LOW notes: <any>."
```
