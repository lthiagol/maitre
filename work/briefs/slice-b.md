You are the **runner** for the vocabulary rename of this project. The orchestrator's pane id is `$ORCH_PANE` — the orchestrator sets it before dispatching. Read `docs/07-vocabulary.md` (the decision of record), `AGENTS.md` (house style) and `docs/03-domain-model.md` (the spec) before you touch anything.

## Your slice — B of 5

The roles get their plain names, and the tabs get their honest names.

**Scope discipline — do only these two things.** Do NOT touch the module path, the binaries, the recipe ids, the domain version (slice A, already done), the backlog (slice C), or Proof→Evidence (slice D).

### AC-01 — the roles are runner, reviewer, orchestrator

Today `core.AgentRoles` is `agent, sous-chef, line-cook, taster, expo`. It becomes `agent, sous-chef, runner, reviewer, orchestrator`.

- Every place a role is written or compared follows: the journal's actor roles, the console's footer and prompts, the CLI's `--role`, the seed, the docs.
- **The self-review invariant must survive the rename, and it is the sharp edge of this AC.** A recipe the executor reviewed itself cannot close. That check compares the executor with the reviewer — and it must keep working when the executor's role is recorded in an *old* journal event (`line-cook`) and the reviewer's is a new one (`reviewer`). Read the check before you touch it; if it compares role *strings* through a lookup that only knows the current set, a recipe executed before the rename could be closed by its own executor. Prove it with a test: a review round recorded under the old role names still cannot close its own recipe.
- The journal is append-only history: **do not rewrite the old role names in `kitchen/activity.json`.** A reader that meets an old role name must not fail.

### AC-02 — the tabs are Board and Sous-Chef

- The dependency-tree tab (today `Chopboard`) is **`Board`**.
- The automation's tab (today `Expo`) is **`Sous-Chef`**.
- The tab constants, the labels, the footer's key hints, the `--render`/`--keys` behaviour and the goldens follow.
- The Sous-Chef tab keeps today's honest content: the automation is **not built yet**, and the tab says so, plus the vital signs it already shows (the plan's domain version, this build's, whether they agree, the last journal event). Do not build a crew view — the autopilot does not exist, and a tab that is waiting must not look like a tab that is failing.

### AC-03 — the console shows one word per fact

- A recipe's row shows the **plain state word** (`closed`, `executing`). A lane header is a **section** (`Served`, `On the line`) and never a second name for the row's state.
- Check the console for any place where a lane word and a state word appear for the same thing, and make it one. If it already holds, prove it with a test and say so in your report — do not invent a change to look busy.

### AC-04 — the tripwires still hold

- The derived-views table, the gate table, the command count (still 46 in this slice), and the plan-format fields all stay as they are.
- Verification: `go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v` exit 0.

## Invariants you must not break
- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps 28 external modules.
- One fact, one answer. The domain layer has **zero UI imports**.
- The kitchen theme stays; the *artifact* is still a recipe.
- Do not use `git commit --no-verify`. Do not push. Do not touch other panes. Do not spawn agents.

## Workflow
1. One commit per AC, with the per-AC evidence stamped in the commit body:
```
Per-AC evidence:
- AC-01: go test ./internal/core -run 'Role|SelfReview' -v exit 0 (N/N pass)
```
2. Before notifying, run the discipline gate: `make lint` clean, `gofmt -l .` empty, `make test` green, `git status` clean, and `make regen-goldens` if the console's output changed.
3. Notify the orchestrator with a **shell command**:
```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-b.md <<EOF
Runner: slice B done. Commits: <n>. Per-AC evidence: AC-01..AC-04 stamped. Pre-notify gate: pass. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-b"
herdr agent prompt "$ORCH_PANE" "Runner: slice B done. Commits: <n>. Per-AC evidence: AC-01..AC-04 stamped. Pre-notify gate: pass. Notes: <anything the reviewer should know>."
```
