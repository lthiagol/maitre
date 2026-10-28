You are the **runner** for slice H — the quality gate. The orchestrator's pane id is
`w1Z:p1` — it is already substituted in this file; wherever this brief says
`$ORCH_PANE`, it means `w1Z:p1`. Use it verbatim in the notify command below. The
repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` (house style and
the working rules) first.

## The chef's ruling — one interface

> *let's ensure it's well documented and organized, so we use just the make tool
> instead of running direct testing code.*

**Every quality command is a `make` target.** The briefs, the docs and the agents call
`make`; nobody runs `go test`, `go vet` or `gofmt` by hand. **The gate is one
command.**

**The tools are dev tools, not module dependencies.** `internal/core` and
`internal/cli` stay stdlib-only and `internal/ui` stays at its budget — `go install`
does not touch `go.mod`.

## Your slice — H — the gate

### AC-01 — the make interface

The `Makefile` grows the whole interface. **Every target below must exist, and every
one must be callable on its own.**

| Target | What it does |
|---|---|
| `make fmt` | `gofmt -w .` |
| `make fmt-check` | `gofmt -l .` must be **empty** |
| `make vet` | `go vet ./...` |
| `make staticcheck` | `staticcheck ./...` |
| `make lint` | `fmt-check` + `vet` + `staticcheck` |
| `make test` | `go test ./...` — **with `PKG` and `RUN` overridable** (see AC-02) |
| `make test-race` | `go test -race -shuffle=on -count=2 ./...` |
| `make goldens` | regenerate the goldens (today's `regen-goldens`) |
| `make goldens-check` | regenerate, then the goldens must be **unchanged** — `git diff --exit-code` over the goldens' path only |
| `make budget` | the dependency budget (AC-03) |
| `make tidy-check` | `go mod tidy` must produce **no diff** |
| `make check` | the **fast gate**: `lint` + `test` + `goldens-check` + `budget` + `tidy-check` |
| `make check-race` | the **slow gate**: `lint` + `test-race` + `goldens-check` + `budget` + `tidy-check` |
| `make tools` | install the two tools, **pinned** (AC-04) |
| `make vuln` | `govulncheck ./...` |
| `make ci` | what CI runs: `check-race` + `vuln` — and it **prints a header per check**, so the log names the check that failed |

**The old `make lint` grows, and `make regen-goldens` becomes `make goldens`** — keep
the old name as an alias if anything still calls it, and say what you kept.

**The gate must not skip a check.** If `staticcheck` is not installed, `make lint`
**fails** with a message that says how to get it (`make tools`) — a gate that quietly
skips a check is a gate that has stopped seeing. Say in your report how a *missing*
tool is reported.

### AC-02 — a focused check, through make

The lanes' per-AC evidence is a focused run, and it must go through `make` too:

```bash
make test PKG=./internal/core RUN='-run Role|SelfReview'
```

`PKG` defaults to `./...` and `RUN` to empty. **Prove the flags reach `go test`** —
and prove that `make test PKG=./internal/core RUN='-run NoSuchTest'` **fails** (a
filter that matches nothing is a green that proves nothing — this repo has the lesson:
*a silent enum fallback disabled a whole feature, with green tests*).

### AC-03 — the dependency budget is checked, and the doc cannot drift

`AGENTS.md`'s dependency table states the budget, and the closed finding **P-6** says
the count and the counting command are written down. **Now they are checked.**

**The requirement:**

1. The live external-module counts for `internal/core`, `internal/cli` and
   `internal/ui` are computed **the way `AGENTS.md` says to compute them**.
2. The **ceilings** (0 / 0 / ≤ 30) are held against the live counts.
3. **`AGENTS.md`'s table is checked against the same numbers**, so the doc cannot
   drift from the gate — including the *current* count it states.
4. The failure message names **the live count, the expected count, and the file to
   change**. A gate that says "budget failed" is a gate that costs an hour.
5. **A dependency that moves the count fails the gate** — and the gate says exactly
   what to change. Prove it: add a dependency to `internal/ui` (a temporary edit,
   reverted) and confirm the gate fails and says the right thing.

**Where the numbers live is yours to choose — but there must be ONE place.** Either
the gate holds them and checks the doc against them, or the doc holds them and the
gate reads them. **Do not have both.** Say which you chose and why.

### AC-04 — the tools, pinned, and CI

**The tools, pinned:**

| Tool | Version | Why |
|---|---|---|
| `staticcheck` | `honnef.co/go/tools/cmd/staticcheck@v0.8.2` | `go vet` is weak; staticcheck finds the `SA*`, `S1*` and `ST*` classes vet misses. One binary, no config, no plugin surface. |
| `govulncheck` | `golang.org/x/vuln/cmd/govulncheck@v1.9.0` | 28 external modules, some very new; it finds the vulnerabilities that reach the binaries we ship. It reads the module graph, so it runs after the modules are downloaded. |

**`make tools` installs them**, and it says where they went (`go env GOPATH/bin`) and
what to do if that directory is not on the `PATH`. **Pinned, not `@latest`** — the
repo pins its vendored skills for the same reason.

**CI — `.github/workflows/ci.yml`:**

- it runs on **push** and **pull_request**;
- `permissions: contents: read`;
- it sets up Go with the **declared floor** — `go.mod` says `go 1.26.0`, and the
  floor is the honest version to test: the code must build on the version it claims;
- it caches the module downloads;
- it installs the tools (`make tools`);
- it runs the gate **through `make`** — the list of checks lives in the `Makefile`,
  not in the workflow, so there is one list. The log names the check that failed.

**And the pre-push hook.** A committed hook (`scripts/hooks/pre-push.sh`) that runs
the **fast** gate, plus a `make install-hooks` that copies it into `.git/hooks/`. The
repo's rule is *a lane commits per acceptance criterion, as it goes*, and the hook
keeps a broken commit from ever landing. The **slow** gate (`-race`) stays in CI: it
takes minutes, and that belongs in CI rather than in a push.

### AC-05 — the reasons, written down

**The tools' reasons are the paper, and the paper is the point.** Each tool gets a
**written reason**, in the same spirit as *"a dependency needs a written reason"*.

| Where | What |
|---|---|
| `AGENTS.md` | the build-commands table gains the targets; the dependency policy points at `make budget`; a new short section **"The quality gate"** with each tool's reason, one line each; and a **working rule**: *the gate is `make check`; a focused check is `make test PKG=… RUN=…`; do not run the tools directly.* |
| `docs/testing.md` | a section on the gate, the tools and CI — it already says *"No pty test in CI"*, so this is where the CI story belongs. |
| **Your own brief** | this file is the example: the per-AC evidence lines are `make test PKG=… RUN=…`, not raw `go test`. |

**The reasons, in the orchestrator's words — use them, or better ones:**

- **`staticcheck`** — `go vet` is weak, and staticcheck finds the classes vet misses.
  It is one binary with no configuration, so it cannot become a project of its own and
  fight the house style.
- **`govulncheck`** — 28 external modules, and it finds the vulnerabilities that
  reach the binaries we ship. It is not a gate on the code; it is a gate on the
  supply chain.
- **`-race`** — the console is a concurrent program: the poll, the terminal, the
  reload. The race detector is the only check that finds a data race.
- **`-shuffle=on`** — a test that depends on another's order is a test that lies.
- **`-count=2`** — it catches state that leaks between runs: a cached global, a
  shared temp directory, a package-level var that a test forgot to reset.
- **`go mod tidy` produces no diff** — a stale `go.mod` is drift.
- **CI** — the gate runs where the code lands, not where a human remembers.

**And one thing the paper should say plainly:** these are **dev tools**, not module
dependencies. The budget of 0 / 0 / ≤ 30 is untouched by them.

### AC-06 — the whole tree, verified

```bash
make check-race && git status --short
```

Green, and the tree clean. Report the **real** numbers: the counts, the tool
versions, and how long the slow gate takes.

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. The tools do not touch `go.mod`.
- The domain layer has **zero UI imports**.
- One fact, one answer: the numbers are in one place; the gate checks the other copy.
- `core.DomainVersion` is **7**. Do not move it.
- The gate table, the plan format and the command count (**42**) must not move.

## The scope — what not to touch

| | |
|---|---|
| **The names, the paper, the glossary, `docs/07`** | settled. Do not touch. |
| **`docs/audit.md`** | the ledger. **Never swept.** |
| **The code** | **no behaviour changes.** This slice is tooling, configuration and documentation. If a fix would change what the tool does, **stop and report it**. |
| **`work/`** | the orchestrator's. |

## The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. Leave anything you did
not write alone. If a file you did not write has changed, **stop and report it**.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-h.md <<EOF
Runner: slice H done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. The numbers live in: <where, and why>. The slow gate takes: <how long>. Tools: <versions>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-h"
herdr agent prompt "w1Z:p1" "Runner: slice H done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. The numbers live in: <where, and why>. The slow gate takes: <how long>. Tools: <versions>. Notes: <anything the reviewer should know>."
```

**Write the receipt, commit it, and *then* nudge.** The nudge only reaches an idle
orchestrator; the receipt is the half that cannot be lost.

**One write per file, and commit as soon as the file is coherent.**
