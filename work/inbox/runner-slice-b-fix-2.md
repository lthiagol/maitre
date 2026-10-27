# Runner — slice B, fix cycle 2

**Commits:** 1.

| Commit | What |
|---|---|
| `edc08b6` | **F-72** — the table's roles move to the plain column, and `maître d'` goes |

## What changed

`docs/03-domain-model.md` § *The CPD ↔ tool mapping* — the table only.

**Before**

| Themed | Kept plain |
|---|---|
| recipe · special · prep note · course · menu | `acceptance_criteria` · `scope` · `evidence` · `findings` · `work_packages` |
| chef · maître d' · sous chef | runner · reviewer · orchestrator · lifecycle states · gate names · `status` · `review_state` |
| `maitre-cli` · `maitre` · the skills | |

**After**

| Themed | Kept plain |
|---|---|
| recipe · special · prep note · course · menu | `acceptance_criteria` · `scope` · `evidence` · `findings` · `work_packages` |
| | `chef` · `sous-chef` · runner · reviewer · orchestrator · lifecycle states · gate names · `status` · `review_state` |
| `maitre-cli` · `maitre` · the skills | |

- `chef` and `sous-chef` **moved** to the plain column, spelled as the code
  spells them.
- `maître d'` **went**. It was never a job — per the old glossary it was the
  *human-facing console*'s themed name, and the console is `maitre` now, already
  in the themed column's third row.
- The artifact row and the records are untouched, and `special` and `prep note`
  stay: the backlog is slice C's business.

The rule sentence above the table now agrees with it, and both agree with
`docs/07-vocabulary.md` § *The rule*.

## Notes for the reviewer

1. **No name vanished from the document.** The deny-list wire reads the whole
   document, so moving the cells could have taken the last mention of a role
   with it. It did not: the prose at § *the roles* (3.2) still names every role
   the code refuses and every retired name. Counts in `docs/03` after the move —
   `agent` 16, `sous-chef` 4, `runner` 6, `reviewer` 10, `orchestrator` 2,
   `line-cook` 1, `taster` 1, `expo` 1, `coordinator` 1. Nothing was restored
   into the prose, because nothing was lost from it.

2. **`maître d'` is gone from `docs/03` entirely** — the table was its only
   mention there. It survives in `docs/decisions/0004-naming.md` and
   `0005-vocabulary.md` (dated records, untouched) and in **`GLOSSARY.md`**, which
   still maps *the human-facing console → maître d'*. That mapping is now wrong
   on two counts — the console is `maitre`, and the cast is plain — and it is
   slice E's rewrite. Noted, not touched.

3. **The table's second row now has an empty Themed cell.** That is the honest
   shape: every job is plain, so nothing in that row is themed, and the
   row's plain cell still carries the machine's facts (lifecycle states, gate
   names, `status`, `review_state`). Nothing was restructured beyond what the
   ruling asked for.

4. **No code changes.** `internal/core`, `internal/cli` and `internal/ui` are
   untouched.

## The gate

```
make lint && gofmt -l . && make test && git status --short
```

Green: `go vet` clean, `gofmt -l .` empty, core · cli · ui all pass, and the
tree clean apart from the orchestrator's own `work/` files.
