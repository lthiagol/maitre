You are the **runner** for the vocabulary rename of this project. The orchestrator's pane id is `w1Z:p1` — it is already substituted in this file; wherever this brief says `$ORCH_PANE`, it means `w1Z:p1`. Use it verbatim in the notify command below. The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `docs/07-vocabulary.md` (the decision of record), `AGENTS.md` (house style) and `GLOSSARY.md` before you touch anything.

## What slices A–D already did — the names you must not undo

The project is `maitre`, the binaries are `maitre` and `maitre-cli`, the ids are `RCP01`
and `BL01`, the roles are `runner`, `reviewer`, `orchestrator`, `sous-chef`, the tabs are
**Board** and **Sous-Chef**, the backlog is one artifact with a `kind`, and the
per-criterion record is **evidence** while the recipe-level stamp is **completion**.
The plan format is version **6** after slice D.

**The rule:** the kitchen names the artifact and the place; plain words name the
machine's records and the agents' jobs. `docs/07-vocabulary.md` is the decision of
record, and it wins over any other doc.

## Your slice — E of 5, the last — the glossary, the ADR and the sweep

**Scope discipline — do only this.** The names are already in the code (slices A–D). This slice is the paper.

### AC-01 — the mascot — **NOT YOURS. Another session is drawing it.**

The mascot is being made **in parallel, in another session**, on 2026-10-27. Do not
draw it, do not ship it, and do not touch the files it lives in:

| | |
|---|---|
| **The art** | the drawing, its full and compact forms, wherever they are placed |
| **The console's header** | the compact form, before the product name |
| **The CLI's help header** · **the README** · **the Sous-Chef tab** | the full form |
| **The goldens** | the mascot will move them — **regenerate, do not hand-edit** |

**Do not delete or replace anything the other session lands.** If you find yourself
needing a change in one of those places — a layout that the mascot breaks, a header
that no longer fits — **stop and report it to the orchestrator** instead of editing
around it. Two sessions in one tree is a collision, and a collision in the mascot's
files is a merge nobody wins.

**What is still yours in AC-01:** nothing. The acceptance criterion moves to the
session that owns the drawing. Your report must say whether the mascot landed, so
the orchestrator can close the slice.

### AC-02 — the glossary is rewritten

`GLOSSARY.md` currently describes "Kitchen Brigade" and the old names. Rewrite it to the decision of record:

- the **scope rule** is stated at the top: *the kitchen names the artifact and the place; plain words name the machine's records and the agents' jobs* — with a short table of which layer speaks which vocabulary;
- **maitre** is the project and the console; `maitre-cli` is the CLI;
- the **backlog item** replaces *special* and *prep note*, with its three kinds;
- the roles are **chef, sous-chef, runner, reviewer, orchestrator**;
- a recipe's record is **evidence**; the recipe-level stamp is **completion**;
- the kitchen words that survive: recipe, kitchen, line, board, service, tasting, sign-off, method, course, menu — each with its `_Avoid_` line, as today;
- the **CPD mapping table** is updated to match, and it stays the single source for the mapping. **One row per CPD term** — slice C collapsed *idea / note* and *brief* into a single row and left the split in the `_Avoid_` line, and a mapping that cannot be read one term at a time is not a mapping. The backlog item's kinds and the CPD terms are reconciled here, and the reconciliation is stated.

### AC-03 — the ADR

`docs/decisions/0004-naming.md` was marked *superseded in part* by slice A. Write the decision properly:

- `docs/decisions/0009-vocabulary.md` — the vocabulary decision: the scope rule, the rename table, the backlog's shape, and the alternatives that lost, with the reason each lost.
- Update `docs/decisions/README.md` (the ADR index) and `0004-naming.md` so it says plainly that it is superseded by 0009, and **keep its tables as the record of what it said**.
- The house rule is *docs state intent; ADRs record choices* — so this ADR is where the choice lives, and `docs/07-vocabulary.md` is where the chef's decision of record lives. Both point at each other.

**One more ADR, and the line that decides it.** The rename moved the *incidental*
vocabulary in the ADRs whose subject is the **mechanism** (0002, 0006, 0008) and left
the ones whose subject is the **vocabulary** (0004, 0005) as the record of what they
said. That line is settled.

**But `docs/decisions/0006-verification.md` § 2 is wrong about the code**, and it is
this slice's to fix because this slice touches the ADRs:

> ### 2. Evidence has one of three match states

The code has **four** — `proven` · `rejected` · `unknown` · `signed-off` — and the
ADR's table lists three, under a name the code does not use (*unproven* for
*unknown*). The sign-off is the one the ADR's own section is about, and it is missing.

**A dated record of a *choice* is not a licence to misdescribe the mechanism.** The
ADR's table must say what the code has. Fix it here, and say in your report exactly
what you changed.

### AC-04 — the sweep

Every remaining mention of an old name is gone from the docs and the code: `Kitchen Brigade`, `scale` (the binary), `line-cook`, `taster`, `expo` (the role), `Chopboard`, `special`, `prep note`, `Proof` (where it means the record — and the record is called **evidence** now, so a leftover *proof* that means the record is drift too), `R01`-style ids, `github.com/lthiagol/kitchen-brigade`.

- `grep` the tree and prove it. Report the exact commands and their output in the commit message.
- **`docs/audit.md` is the ledger of the night — do not sweep it.** Its rows record the names as they were when the finding was filed; rewriting them would rewrite the record. Leave them, and say so in your report.
- `research/00-legacy-analysis.md` is about **v1** — its mentions of v1's vocabulary are history, not drift. Leave them.
- The vendored `skills/` are third-party, pinned to an upstream commit. **Do not edit them.** If they mention an old name, say so in your report instead.

### AC-05 — the whole tree, verified

**Another session may land the mascot while you work.** Before you sweep, and before
you regenerate the goldens, **re-read the files you are about to touch** — the
mascot's files may have changed under you. The sweep is about **old names**; the
mascot is about **new art**, and the two should not meet in the same diff.

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
mkdir -p work/inbox
cat > work/inbox/runner-slice-e.md <<EOF
Runner: slice E done. Commits: <n>. Per-AC evidence: AC-02..AC-05 stamped. Pre-notify gate: pass. Mascot: <did it land? yes/no>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-e"
herdr agent prompt "w1Z:p1" "Runner: slice E done. Commits: <n>. Per-AC evidence: AC-02..AC-05 stamped. Pre-notify gate: pass. Mascot: <did it land? yes/no>. Notes: <anything the reviewer should know>."
```

### Two LOW notes from slice A's review — fix them here

The taster noted two gaps and did not file them as blockers. They are in scope for this slice:

1. **The format facts list misses two kinds of fact** (`internal/core/format.go`): the **root file names** (`project.json`, `activity.json`, `decisions.json`) and the **journal envelope's own fields** (`schema`, `seq`, `events`). Renaming `activity.json` or changing the envelope's `schema` field leaves the fingerprint unchanged, so such a change would need the version moved **by hand** — which is the hole F-71 was filed to close. The journal's file name is the silent one: a missing journal reads as an empty one. **Two lines in `planFormatFacts`** close it.
2. **The usage-lines wire scans `internal/cli` only.** A usage line added in another package would not be caught. The console's errors are checked separately — extend the wire's scope so it covers every package that prints a usage line, and prove it by adding a usage line in another package and confirming the wire fails.
