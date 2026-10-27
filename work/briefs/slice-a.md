You are the **line cook** for the vocabulary rename of this project. The orchestrator is in pane `w1R:p1`. The decision of record is `docs/07-vocabulary.md` — read it first, then `AGENTS.md` (house style), then `docs/03-domain-model.md` (the spec).

## Your slice — A of 5

The project is `maitre`. The console binary stays `maitre`. The CLI binary becomes `maitre-cli`. Recipe ids become `RCP01`. The domain version becomes 3.

**Scope discipline — do only these four things.** Do NOT touch the roles (runner/reviewer/orchestrator), the tabs (Board/Sous-Chef), the backlog (one item with a kind), or Proof→Evidence. Those are slices B, C and D, and they come after this one. A slice that touches the next slice's files makes the review meaningless.

### AC-01 — the module path
- `go.mod` says `module github.com/lthiagol/maitre`; every import path in the tree follows.
- Verification: `go build ./... && go test ./...` exit 0.

### AC-02 — the binaries
- `cmd/maitre` still builds the console binary named `maitre`.
- `cmd/scale` becomes `cmd/maitre-cli` and builds a binary named `maitre-cli`.
- The Makefile builds both into `./bin`; `scripts/console-shot.py` points at the right binary.
- Verification: `make build && ls bin` shows `maitre` and `maitre-cli`; `make test` exit 0.

### AC-03 — the product name in the strings
- The console's window title and header say `maitre`, not "Kitchen Brigade".
- The CLI's help says `maitre-cli — the maitre agent CLI`.
- The docs' first lines and headers that name the *tool* say maitre.
- Verification: `./bin/maitre --render --plain` shows the maitre title; `./bin/maitre-cli --help` shows the new help; the golden files are regenerated with `make regen-goldens`; `make test` exit 0.

### AC-04 — recipe ids are `RCP01`
- The writer emits `RCP01`, `RCP02`, … and files named `RCP01-<slug>.json`.
- The reader resolves `RCP01`, `rcp01`, `RCP1`, `R01`, `r1` and `1` to the same recipe, and reads a file named `R01-<slug>.json` as well as `RCP01-<slug>.json`.
- `normaliseID` must strip **every** leading letter, not just the first: today's normaliser removes one letter, so `RCP01` would normalise to `CP01` and stop matching `R01`. That is the bug this AC exists to prevent — prove it with a test.
- The demo plan (`core.Seed`) and the repo's own `kitchen/` use the new ids.
- The repo's own plan is migrated: the recipe file is renamed, and its `id` and `depends_on` are updated. **The journal `kitchen/activity.json` is append-only history — do NOT rewrite the ids in it.** A test proves the journal's old `R01` still resolves to the renamed recipe.
- Verification: `go test ./internal/core -run 'ID|Id|Normal' -v` — the new tests pass; `go test ./...` exit 0.

### AC-05 — the domain version is 3
- `core.DomainVersion` is 3, and the comment says why the rules changed (the plan format changed: the ids).
- The guard's direction is unchanged: a tool older than the plan refuses to write; a tool newer than the plan proceeds and stamps the newer version on its next write. The repo's own plan (stamped 2) stays writable.
- `TestDomainShapeIsVersioned` passes.
- Verification: `go test ./internal/core -run 'Domain|Version' -v` exit 0.

### AC-06 — the tripwires still hold
- The doc-vs-code tripwires pass unchanged: the derived-views table, the gate table, the command count (still **46** — the backlog collapse comes in slice C), the plan-format fields.
- Verification: `go test ./internal/core ./internal/cli -run 'Doc' -v` exit 0.

## Invariants you must not break
- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its current external module count (28).
- One fact, one answer: no second source of truth for any fact.
- The domain layer has **zero UI imports**.
- The kitchen theme and the word "recipes" stay. You are renaming, not rewriting the product.
- Do not use `git commit --no-verify`. Do not push. Do not touch other panes. Do not spawn agents.

## Workflow
1. Read `docs/07-vocabulary.md`, `AGENTS.md`, `docs/03-domain-model.md`.
2. **One commit per AC**, with the per-AC evidence stamped in the commit body:
```
Per-AC evidence:
- AC-01: go build ./... && go test ./... exit 0 (N/N pass)
```
3. Before notifying, run the discipline gate: `make lint` clean, `gofmt -l .` empty, `make test` green, `git status` clean.
4. Notify the orchestrator with a **shell command** (not a printed message):
```bash
herdr agent prompt w1R:p1 "Line cook: slice A done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. Notes: <anything the reviewer should know>."
```
