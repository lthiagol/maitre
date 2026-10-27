# Slice B — the state, and how to finish it

**This is a continuation brief.** A fresh agent can pick slice B up from this file
alone, with a clean context. The original brief is [slice-b.md](slice-b.md); this one
supersedes it for the work that is left.

## The situation

Slice B was dispatched to the lane then called `line-cook` (now **`runner`**, pane
`w1R:p6`) on 2026-10-27. **That agent's turn ended without notifying anyone** — no
message reached the orchestrator, and **nothing was committed**. Its last visible act
was reading `docs/03-domain-model.md` (lines 130–145 and 465–480).

The work is **not lost**: it sits in the working tree, uncommitted.

| | |
|---|---|
| **Files changed** | 38 |
| **The diff** | 223 insertions, 172 deletions — captured at [../state/slice-b-wip.patch](../state/slice-b-wip.patch) (1,413 lines) |
| **The tree** | `go build ./...` **OK** · `go test ./...` **green** — core, cli and ui all pass in this state |
| **The last commit** | `baea34e` — the rename's topology note. Slice B has landed **nothing**. |

**The state is coherent and green.** Do not revert it: the rename of the roles is
mechanical work that is done and tested, and throwing it away buys nothing.

## Where each acceptance criterion stands

| AC | What | State |
|--:|---|---|
| **AC-01** | the roles are `runner`, `reviewer`, `orchestrator` | 🟡 **mostly done** — see below |
| **AC-02** | the tabs are `Board` and `Sous-Chef` | 🔴 **not started** |
| **AC-03** | the console shows one word per fact | 🔴 **not started** |
| **AC-04** | the tripwires still hold | 🔴 **not started** |

### AC-01 — what landed

| | |
|---|---|
| `core.AgentRoles` | `agent` · `sous-chef` · `runner` · `reviewer` · `orchestrator` ✅ |
| `core.retiredAgentRoles` | `line-cook`→`runner`, `taster`→`reviewer`, `expo`→`orchestrator`, `coordinator`→`sous-chef` ✅ |
| `core.agentRole()` | resolves a name as written **and** as a role ✅ |
| The self-review comparison | compares the names as written and as roles, so a plan written before the rename still cannot be closed by its own executor ✅ |
| `docs/03` § the roles | documents the retired names and what they stood for ✅ |

### AC-01 — what is left

| | |
|---|---|
| **The CLI's documented examples** | `internal/cli/cli.go:37`, `internal/cli/cli.go:238` and `internal/cli/flags.go:27` still show `--actor line-cook`. They must show a current name. |
| **`docs/05-human-interface.md:292`** | still says *console claiming work as `line-cook`* → `runner`. |
| **The commit** | AC-01 is not committed. The WIP plus these fixes commit together as AC-01. |
| **A test for the retired names** | prove that a review round recorded under `line-cook` still cannot close its own recipe — the brief's sharp edge. Check whether the runner already wrote it; if not, write it. |

### AC-02 — the tabs, not started

`TabChopboard` and `TabExpo` are still the constants. What must change:

| | |
|---|---|
| `internal/ui/model.go` | the constants: `TabChopboard` → `TabBoard`, `TabExpo` → `TabSousChef` |
| `internal/ui/render.go:84` | the tab list |
| `internal/ui/detail.go:12,16,26–29` | the pane's cases, and the words *"the Expo is not a list of recipes"* → the new name |
| `internal/ui/views.go:386` | `"the expo · the orchestrator · not built yet"` → the new tab's name |
| `internal/ui/views.go:876,889` | the tab's own copy — *"The Expo is the orchestrator agent…"* → **Sous-Chef**. The tab keeps its honest content: the automation is **not built yet**, and the tab says so, plus the vital signs it already shows. **Do not build a crew view** — the autopilot does not exist, and a tab that is waiting must not look like a tab that is failing. |
| The footer's key hints, `--render`, `--keys`, the goldens | follow |
| `docs/05-human-interface.md` | the tab names |

The labels are **`Board`** (the plan's shape: what waits on what) and **`Sous-Chef`**
(the crew and the machine's vital signs).

### AC-03 — one word per fact

A recipe's row shows the **plain state word** (`closed`, `executing`); a lane header is
a **section** (`Served`, `On the line`) and never a second name for the row's state.
Check the console for any place where both appear for the same thing, and make it one.
**If it already holds, prove it with a test and say so in your report** — do not invent
a change to look busy.

### AC-04 — the tripwires

The derived-views table, the gate table, the command count (**46** — the backlog
collapse is slice C) and the plan-format fields all agree with the code.

```bash
go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v
```

## The invariants — do not break them

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps **28** external modules.
- One fact, one answer. The domain layer has **zero UI imports**.
- **The journal is append-only.** The 32 historical events keep their old `R01`; the
  migration is event **33**, a `control` event stamped 4. Never rewrite it.
- The kitchen theme stays; the artifact is still a **recipe**.
- `core.DomainVersion` stays **4** in this slice — the roles are vocabulary, not a
  format change. (Slices C and D take it to 5 and 6.)

## The scope — what not to touch

Slice B is the roles and the tabs. **Do not** touch the backlog (slice C: one item
with a `kind`), the ids, the module path, the binaries, `Proof`→evidence (slice D),
the glossary, or the mascot (slice E). A slice that touches the next slice's files
makes the review meaningless.

The ledger ([../../docs/audit.md](../../docs/audit.md)) is **the reviewer's** to write.
Do not edit it.

## The discipline gate — run before you report

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

All green; `git status` clean after the goldens are regenerated. **One commit per
acceptance criterion**, with the evidence stamped in the body:

```
Per-AC evidence:
- AC-01: go test ./internal/core -run 'Role|SelfReview' -v exit 0 (N/N pass)
```

## The notification

Report to the orchestrator with a **shell command**, not a printed message:

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-b.md <<EOF
Runner: slice B done. Commits: <n>. Per-AC evidence: AC-01..AC-04 stamped. Pre-notify gate: pass. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-b"
herdr agent prompt "$ORCH_PANE" "Runner: slice B done. Commits: <n>. Per-AC evidence: AC-01..AC-04 stamped. Pre-notify gate: pass. Notes: <anything the reviewer should know>."
```

## Three ways to continue — the orchestrator's recommendation

| | |
|---|---|
| **A — re-prompt the same runner** | cheapest; it holds the context of what it wrote. But its turn already ended once without reporting, and its context is long. |
| **B — a fresh agent, the WIP left in place, this brief dispatched** | ✅ **recommended.** The context is clean, the mechanical work is kept and tested, and the fresh agent's first act is to read the 1,413-line diff and check it against AC-01 before committing. |
| **C — revert the WIP and start the slice from the original brief** | the cleanest context of all, and it throws away coherent, green work for nothing. |

Under **B**, the fresh agent's first three commands are:

```bash
git diff --stat                    # what is already done
git diff > /tmp/read.patch         # read it
go test ./...                      # the state is green — start from here
```
