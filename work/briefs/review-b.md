You are the **reviewer** — the independent reviewer lane. The orchestrator's pane id is
`w1Z:p1` — it is already substituted in this file; wherever this brief says
`$ORCH_PANE`, it means `w1Z:p1`. Use it verbatim in the notify command below. Your own
pane is not the target.

## Your slice — the review of slice B

Slice B renamed the **roles** (`runner`, `reviewer`, `orchestrator`, `sous-chef`,
`chef`) and the **tabs** (`Board`, `Sous-Chef`). The runner's commits are in the log;
its first commit carries the **WIP that was already in the working tree** when the
slice was re-dispatched — 38 files that were green but uncommitted.

```bash
cd /Users/thiago/Code/Pessoal/maitre
git log --oneline -12
git diff <base>..HEAD --stat      # the base is the commit before slice B's first commit
```

**Cycle 1: blockers only.** File **HIGH or MEDIUM** blockers. No nits, no cosmetics,
no style, no LOW notes filed as rows. **Do not re-file F-56, F-59, F-60, F-61, P-8** —
they are the chef's and stay open.

### Verify each acceptance criterion with a command or a mutation

Not by reading the diff and believing it. Mutate, watch it fail, revert, and prove the
tree is clean afterwards.

**AC-01 — the roles.**

- `core.AgentRoles` holds the current roles; `core.retiredAgentRoles` maps
  `line-cook`→`runner`, `taster`→`reviewer`, `expo`→`orchestrator`,
  `coordinator`→`sous-chef`.
- `core.agentRole()` resolves a name **as written** and **as a role**. Prove it:
  a plan written before the rename still cannot be closed by its own executor.
- **The sharp edge**: a review round recorded under `line-cook` still cannot close its
  own recipe. Find the test that holds it. If you mutate the retired-name resolution
  (make `agentRole` return `""`, or drop one mapping), the test must fail — otherwise
  the test proves nothing.
- The CLI's **documented examples** show a current name, not `line-cook`.
  `internal/cli/cli.go` and `internal/cli/flags.go` were named in the brief.
- `docs/05-human-interface.md` no longer says *line-cook*.

**AC-02 — the tabs.**

- The constants are `TabBoard` and `TabSousChef`; the tab list, the detail pane's
  cases, the pane's own copy and the footer's key hints all follow.
- **The retired tab names are gone**: `Chopboard` and `Expo` must not survive in the
  code or the docs — the root `AGENTS.md` and `docs/05-human-interface.md` now say
  **Board** and **Sous-Chef**. Prove it:
  ```bash
  grep -rn "Chopboard\|Expo" --include="*.go" --include="*.md" . \
    | grep -v '^\./docs/audit.md' \
    | grep -v '^\./docs/07-vocabulary.md' \
    | grep -v '^\./research/' \
    | grep -v '^\./work/'
  ```
  The three exclusions are **records, not drift**: the ledger keeps the names as they
  were when each finding was filed, `docs/07-vocabulary.md` is the decision of record
  and its **rename table** names the old word to say what it became, `research/` is
  about v1, and `work/` holds the briefs and the receipts. **But check the rename
  table itself**: every mention of `Chopboard` or `Expo` there must sit in a table
  that says what it became — not in prose that still calls the tab by the old name.
- **The honesty holds**: the tab still says the automation is **not built yet**, and it
  still shows the vital signs it showed before. It must **not** have grown a crew view
  — the autopilot does not exist, and a tab that is waiting must not look like a tab
  that is failing.
- **The layout still holds**: the console's header is one line, and the goldens
  were regenerated (`make regen-goldens` must leave `git status` clean). Render the
  console and read it — `maitre --render --plain` — including the detail pane of a
  recipe, and the Sous-Chef tab.

**AC-03 — one word per fact.**

A recipe's row shows the **plain state word** (`closed`, `executing`); a lane header
is a **section** (`Served`, `On the line`) and never a second name for the row's
state. Find any place where both appear for the same thing. If the invariant already
holds, **say so with the evidence that holds it** — do not invent a change.

**AC-04 — the tripwires.**

```bash
go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v
```

The derived-views table, the gate table, the **command count (46)** and the
plan-format fields all agree with the code.

### Then: what did the slice introduce?

1. **The role resolution is the sharp edge of the whole rename.** Read it as if you
   were writing the exploit: could a recipe be closed by the actor that executed it,
   under any spelling of its role? Could an **unknown** role now *pass* where it used
   to fail — a silent enum fallback is one of this project's standing lessons
   (see `AGENTS.md`, *v1 lessons that keep biting*).
2. **The WIP diff came from an agent whose turn ended without reporting.** Read the
   whole diff, not the parts the brief names. It touched `internal/core/model.go`,
   `views.go`, `seed.go`, `chaos.go` and 33 test files. Is every change a role
   rename, or did something else ride along?
3. **One fact, one answer.** Is any role's name now stored in two places that could
   disagree — the constant, the map, the CLI's flag help, the docs?
4. **Scope.** The slice is the roles and the tabs. Confirm it did **not** touch the
   backlog (slice C), the ids, the module path, the binaries, `Proof`→evidence
   (slice D), the glossary or the mascot (slice E).

### The invariants — check them, do not assume them

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps **28**
  external modules.
  ```bash
  go list -deps -f '{{if not .Standard}}{{.Module.Path}}{{end}}' ./internal/ui | sort -u | grep -v lthiagol/maitre | wc -l
  ```
- The domain layer has **zero UI imports**.
- **The journal is append-only**: `kitchen/activity.json` — the **32 historical
  events keep their old `R01`**; the migration is event **33**, a `control` event
  stamped 4. Check with `git diff` on that file, not by reading the prose.
- `core.DomainVersion` stays **4** — the roles are vocabulary, not a format change.

### One question for you — do not fix it

`kitchen/project.json` still carries `"name": "new-project"`. The decision of record
([docs/07-vocabulary.md](../../docs/07-vocabulary.md)) says the project is **`maitre`**.
Judge it: is that drift, and is it in this slice's scope or slice E's sweep? File it
as a finding if it is drift; do **not** change the file.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.** They go in the
  ledger in a later round.

### The ledger

`docs/audit.md` is **yours** to write. The runner did not touch it — confirm that with
`git diff`. Add any NEW blocker as a new row with an id continuing from **F-71**. Do
not re-file what is closed; the ledger is not swept. Commit the ledger.

### Hard rules

- Do **not** edit implementation code. Read-only for everything except
  `docs/audit.md`.
- Revert every mutation you make. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-b.md <<EOF
Reviewer: slice B review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-72.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-b"
herdr agent prompt "w1Z:p1" "Reviewer: slice B review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-72.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>."
```
