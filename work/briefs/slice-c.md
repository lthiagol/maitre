You are the **line cook** for the vocabulary rename of this project. The orchestrator is in pane `w1R:p1`. Read `docs/07-vocabulary.md` (the decision of record), `AGENTS.md` (house style) and `docs/03-domain-model.md` (the spec) before you touch anything.

## Your slice — C of 5 — the backlog becomes one artifact

Today the backlog has **two file types**: a *special* (a hunch) in `kitchen/specials/`, and a *prep note* (intent and scope, before criteria) in `kitchen/prep-notes/`. Two CLI verb groups. They become **one artifact with a kind**.

**Scope discipline — do only this.** Do NOT touch the roles, the tabs, the ids, the module path, the binaries, or Proof→Evidence. Slice A and B are already done; leave their files alone unless the backlog's shape forces a change (say so in your report if it does).

### The artifact

A **backlog item** — id `BL01`, `kind` ∈ `idea | feature | bug` — in `kitchen/backlog/`.

| Field | |
|---|---|
| `id` | `BL01`, `BL02`, … |
| `title` | |
| `body` | why it is worth doing — the old special's body and the old prep note's intent, merged |
| `scope` | the work, if it is known — the old prep note's scope |
| `kind` | `idea` \| `feature` \| `bug` |
| `tags` | |
| `status` | `open` \| `dismissed` \| `promoted` |
| `created`, `updated` | |
| `promoted_to` | the recipe it became |

The kind must be **kind-neutral in the id**: `BL01` stays `BL01` when an idea is shaped into a feature.

### AC-01 — the model and the store

- One artifact type replaces `Special` and `PrepNote`. One directory replaces `kitchen/specials/` and `kitchen/prep-notes/`.
- The reader resolves `BL01`, `bl01`, `BL1`, `B01`, `1` to the same item, and reads a file named `B01-<slug>.json` as well as `BL01-<slug>.json`.
- The id normaliser must handle a **two-letter** prefix as well as the recipe's three-letter one — prove both with a test.
- The old directories are gone from the plan format. `Init` creates `kitchen/backlog/` and does **not** create the old ones.
- Verification: `go test ./internal/core -run 'Backlog|ID|Normal' -v` exit 0.

### AC-02 — the CLI: one noun for the list

`scale special list · create · update · dismiss · promote` and `scale prepnote list · create · update · promote` become:

```
scale backlog list · create · update · dismiss · promote
```

- `create` takes `--kind idea|feature|bug` (default: `idea`).
- `update` may change the kind — **that is the shaping step**: an idea becomes a small task by becoming `feature` or `bug` and gaining a scope. It is an update, not a promotion.
- `dismiss` now works on **every** kind. Today only specials could be dismissed, with no reason behind the asymmetry.
- **The command count goes from 46 to 42.** Update the doc that says so (`docs/04-operations.md`) *and* the tripwire that counts it, in the same commit — a tripwire whose number moved must move with the number, or it stops being a tripwire.
- Verification: `go test ./internal/cli -v` exit 0, and the count in the doc and in the test are the same number.

### AC-03 — how a backlog item becomes work

```
idea ──shaped into──▶ a small task   (kind: idea → feature|bug, plus a scope)
idea ──shaped into──▶ a recipe       (promote)
small task ──promoted──▶ a recipe
```

- `promote` creates a recipe **in `draft`**, carrying the item's intent, its scope and — for a small task — **one** acceptance criterion, so the work is falsifiable from the start.
- The recipe is **not** pre-groomed: the second out-of-scope item, the falsifiable verification and the steps are the recipe's own life, exactly as for any recipe. **No gate is weakened, and no gate rule changes.** If you find yourself needing to weaken a gate or to add a plan-format field to make a promotion work, **stop and report it to the orchestrator instead** — that is a design decision, not an implementation detail.
- The promotion is one-way: the item's status becomes `promoted` and it records what it became.
- Verification: a test proves an idea promotes to a recipe in `draft`, that a small task promotes to a recipe with exactly one acceptance criterion, that the item is `promoted`, and that the promotion is one-way.

### AC-04 — the console

- The Backlog tab groups into **Ideas** and **Small tasks** — **derived** from the kind, never stored as a group.
- A row shows the item's kind.
- The tab keeps working with the old plans: a plan with `kitchen/specials/` on disk must not crash the console. Decide and state what it does with them — reading them is fine if it is honest, and silently ignoring them is not.
- Verification: `maitre --render --all --plain` shows both groups; the goldens are regenerated.

### AC-05 — the seed, the demo and the docs

- `core.Seed` and the chaos seed create backlog items, not specials and prep notes. `scale demo` still produces a plan the console can render.
- `docs/03-domain-model.md` § 6 (the plan format) shows `kitchen/backlog/` and the item's shape, and drops the two old directories.
- `GLOSSARY.md`: replace the **Special** and **Prep note** entries with one **Backlog item** entry (and the CPD mapping row). This is the *minimum* glossary work for this slice; slice E rewrites the rest of the file.
- The repo's own `kitchen/` plan: migrate it. It has a promoted item and a served recipe — the item's file moves to `kitchen/backlog/` with the new id and kind. **The journal `kitchen/activity.json` is append-only history: do not rewrite it.** A test proves the journal's old ids still resolve.

### AC-06 — the tripwires still hold

- The derived-views table, the gate table, the plan-format fields and the command count (now 42) all agree with the code.
- Verification: `go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v` exit 0.

## Invariants you must not break
- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps 28 external modules.
- One fact, one answer — no second source of truth for any fact. The kind is the only truth; the group is derived.
- The domain layer has **zero UI imports**.
- The kitchen theme stays; the artifact is still a recipe.
- Do not use `git commit --no-verify`. Do not push. Do not touch other panes. Do not spawn agents.

## Workflow
1. One commit per AC, with the per-AC evidence stamped in the commit body:
```
Per-AC evidence:
- AC-01: go test ./internal/core -run 'Backlog' -v exit 0 (N/N pass)
```
2. Before notifying, run the discipline gate: `make lint` clean, `gofmt -l .` empty, `make test` green, `git status` clean, `make regen-goldens` if the console's output changed.
3. Notify the orchestrator with a **shell command**:
```bash
herdr agent prompt w1R:p1 "Line cook: slice C done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. Notes: <anything the reviewer should know>."
```

### The version moves with the format

Because the domain's fingerprint now covers the **format facts** (the id prefix, the plan's directories, the plan's record field names), a format change fails `TestDomainShapeIsVersioned` until the version moves. This slice changes the format, so:

- `core.DomainVersion` goes **4 → 5**, `docs/03` § 6's version paragraph says 5, and the recorded fingerprint is updated.
- The comment says why: the plan format changed, so a plan written under version 4 must not be written by a tool that does not know the new format.
