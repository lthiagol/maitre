You are the **line cook** for the vocabulary rename of this project. The orchestrator is in pane `w1R:p1`. Read `docs/07-vocabulary.md` (the decision of record), `AGENTS.md` (house style) and `GLOSSARY.md` before you touch anything.

## Your slice — E of 5, the last — the mascot, the glossary, the ADR and the sweep

**Scope discipline — do only this.** The names are already in the code (slices A–D). This slice is the paper and the mascot.

### AC-01 — the mascot

A **funny agentic maitre**: a chef with a tall toque and a small moustache, holding a clipboard, drawn in ASCII. Funny, not corporate.

- **Two forms**, because the console's header is one line and a multi-line drawing would break the layout:
  - a **full** form — 4 to 6 lines, about 24 columns wide — for the **CLI's help header**, the **README**, and the **Sous-Chef tab** (the tab that says the automation is not built yet gets a face; a waiting tab should not look like a failing one);
  - a **compact** form — **one line**, at most 24 columns — for the **console's header**, before the product name.
- The drawing lives in **one place** in the code, and the two forms are rendered from it or from one constant beside it. Do not paste the same art twice.
- **Propose your drawing in the commit message**, so it can be judged as art and not only as code. If the first one you draw is not funny, draw another and pick the better one. Three candidates, one shipped.
- The mascot must not break the layout: the console's header stays one line, the Sous-Chef tab's content still says the automation is not built, and the goldens are regenerated.

### AC-02 — the glossary is rewritten

`GLOSSARY.md` currently describes "Kitchen Brigade" and the old names. Rewrite it to the decision of record:

- the **scope rule** is stated at the top: *the kitchen names the artifact and the place; plain words name the machine's records and the agents' jobs* — with a short table of which layer speaks which vocabulary;
- **maitre** is the project and the console; `maitre-cli` is the CLI;
- the **backlog item** replaces *special* and *prep note*, with its three kinds;
- the roles are **chef, sous-chef, runner, reviewer, orchestrator**;
- a recipe's record is **evidence**; the recipe-level stamp is **completion**;
- the kitchen words that survive: recipe, kitchen, line, board, service, tasting, sign-off, method, course, menu — each with its `_Avoid_` line, as today;
- the **CPD mapping table** is updated to match, and it stays the single source for the mapping.

### AC-03 — the ADR

`docs/decisions/0004-naming.md` was marked *superseded in part* by slice A. Write the decision properly:

- `docs/decisions/0009-vocabulary.md` — the vocabulary decision: the scope rule, the rename table, the backlog's shape, and the alternatives that lost, with the reason each lost.
- Update `docs/decisions/README.md` (the ADR index) and `0004-naming.md` so it says plainly that it is superseded by 0009, and **keep its tables as the record of what it said**.
- The house rule is *docs state intent; ADRs record choices* — so this ADR is where the choice lives, and `docs/07-vocabulary.md` is where the chef's decision of record lives. Both point at each other.

### AC-04 — the sweep

Every remaining mention of an old name is gone from the docs and the code: `Kitchen Brigade`, `scale` (the binary), `line-cook`, `taster`, `expo` (the role), `Chopboard`, `special`, `prep note`, `Proof` (where it means the record), `R01`-style ids, `github.com/lthiagol/kitchen-brigade`.

- `grep` the tree and prove it. Report the exact commands and their output in the commit message.
- **`docs/audit.md` is the ledger of the night — do not sweep it.** Its rows record the names as they were when the finding was filed; rewriting them would rewrite the record. Leave them, and say so in your report.
- `research/00-legacy-analysis.md` is about **v1** — its mentions of v1's vocabulary are history, not drift. Leave them.
- The vendored `skills/` are third-party, pinned to an upstream commit. **Do not edit them.** If they mention an old name, say so in your report instead.

### AC-05 — the whole tree, verified

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

- All green, the worktree clean after the goldens are regenerated.
- The dependency budget still holds: `internal/core` and `internal/cli` stdlib-only; `internal/ui` at 28 external modules. Count them the way `AGENTS.md` says to count them and report the number.
- The doc tripwires hold: the derived-views table, the gate table, the plan format, the command count (42).

## Invariants you must not break
- One fact, one answer. The domain layer has **zero UI imports**.
- The kitchen theme stays. The artifact is still a recipe.
- Do not use `git commit --no-verify`. Do not push. Do not touch other panes. Do not spawn agents.

## Workflow
1. One commit per AC, with the per-AC evidence stamped in the commit body.
2. Before notifying: the AC-05 commands, all green, `git status` clean.
3. Notify the orchestrator with a **shell command**:
```bash
herdr agent prompt w1R:p1 "Line cook: slice E done. Commits: <n>. Per-AC evidence: AC-01..AC-05 stamped. Pre-notify gate: pass. Mascot: <paste the shipped drawing>. Notes: <anything the reviewer should know>."
```

### Two LOW notes from slice A's review — fix them here

The taster noted two gaps and did not file them as blockers. They are in scope for this slice:

1. **The format facts list misses two kinds of fact** (`internal/core/format.go`): the **root file names** (`project.json`, `activity.json`, `decisions.json`) and the **journal envelope's own fields** (`schema`, `seq`, `events`). Renaming `activity.json` or changing the envelope's `schema` field leaves the fingerprint unchanged, so such a change would need the version moved **by hand** — which is the hole F-71 was filed to close. The journal's file name is the silent one: a missing journal reads as an empty one. **Two lines in `planFormatFacts`** close it.
2. **The usage-lines wire scans `internal/cli` only.** A usage line added in another package would not be caught. The console's errors are checked separately — extend the wire's scope so it covers every package that prints a usage line, and prove it by adding a usage line in another package and confirming the wire fails.
