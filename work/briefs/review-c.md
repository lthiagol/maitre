You are the **reviewer** — the independent reviewer lane, fresh for slice C. The
orchestrator's pane id is `w1Z:p1` — it is already substituted in this file; wherever
this brief says `$ORCH_PANE`, it means `w1Z:p1`. Use it verbatim in the notify command
below. The repository is `/Users/thiago/Code/Pessoal/maitre`.

## Your slice — the review of slice C

Slice C made the backlog **one artifact with a kind**: the *special* and the *prep
note* are gone, and `BL01` with `kind: idea | feature | bug` took their place.

```bash
cd /Users/thiago/Code/Pessoal/maitre
git log --oneline -12
git diff <base>..HEAD --stat      # the base is the commit before slice C's first commit
```

**Cycle 1: blockers only.** File **HIGH or MEDIUM** blockers. No nits, no cosmetics,
no style, no LOW notes filed as rows. **Do not re-file F-56, F-59, F-60, F-61, P-8**
(the chef's) or **F-73** — F-73 lands with this slice, so judge it here and close it
if it is fixed.

### Verify each acceptance criterion with a command or a mutation

Not by reading the diff and believing it. Mutate, watch it fail, revert, and prove the
tree is clean afterwards.

**AC-01 — the model and the store.**

- One artifact type replaces `Special` and `PrepNote`; one directory replaces
  `kitchen/specials/` and `kitchen/prep-notes/`.
- The reader resolves `BL01`, `bl01`, `BL1`, `B01`, `1` to the same item, and reads a
  file named `B01-<slug>.json` as well as `BL01-<slug>.json`.
- The id normaliser handles a **two-letter** prefix as well as the recipe's
  three-letter one. **Mutate it** — make it accept only `BL` — and confirm the tests
  fail.
- `Init` creates `kitchen/backlog/` and **does not** create the old ones. The old
  directories are gone from the plan format.

**AC-02 — the CLI.**

- One noun for the list: `maitre-cli backlog list · create · update · dismiss ·
  promote`. The old `special` and `prepnote` verb groups are gone.
- `create` takes `--kind idea|feature|bug`, default `idea`; `update` may change the
  kind — that is the shaping step; `dismiss` works on **every** kind.
- **The command count goes from 46 to 42, and the count moved with its tripwire.**
  The doc that says 46→42 and the test that counts must say the **same** number, in
  the same commit. A tripwire whose number moved but did not move with the number is
  no longer a tripwire — **prove the wire fires** by mutating the count.

**AC-03 — how a backlog item becomes work.**

- `promote` creates a recipe **in `draft`**, carrying the item's intent and its scope,
  and — for a small task — **exactly one** acceptance criterion.
- The promotion is **one-way**: the item's status becomes `promoted` and it records
  what it became (`promoted_to`).

**AC-04 — the console.**

- The Backlog tab groups into **Ideas** and **Small tasks**, **derived** from the kind
  and **never stored** as a group. Check the model: a stored group would be a second
  source of truth, and the kind is the only truth.
- A row shows the item's kind.
- **An old plan does not crash the console.** A plan with `kitchen/specials/` on disk
  must render, and the console must be **honest** about what it found — reading them
  is fine, silently ignoring them is not.

**AC-05 — the seed, the demo, the docs and the repo's own plan.**

- `core.Seed` and the chaos seed create backlog items; `maitre-cli demo` still
  produces a plan the console can render.
- `docs/03` § 6 shows `kitchen/backlog/` and the item's shape, and drops the two old
  directories.
- `GLOSSARY.md` replaces the **Special** and **Prep note** entries with one **Backlog
  item** entry. *(The rest of the glossary is slice E's — do not judge it here.)*
- **The repo's own `kitchen/` plan is migrated**: the item's file moves to
  `kitchen/backlog/` with the new id and kind.
- **F-73 is closed here:** `kitchen/project.json` must no longer say
  `"name": "new-project"` — the decision of record names the project **maitre**, and
  the console's header renders the plan's name.

**AC-06 — the tripwires still hold.**

```bash
go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v
```

The derived-views table, the gate table, the plan-format fields and the command count
(now **42**) all agree with the code.

### Then: what did the slice introduce?

1. **The journal is append-only.** `kitchen/activity.json` — the migration of the
   repo's own plan must **not** rewrite it. Check with `git diff` on that file, not by
   reading the prose. And the **old ids must still resolve**: a test proves the
   journal's pre-migration ids (`S01`, `P01`) still resolve after the plan moved.
2. **The promotion weakened nothing.** The brief said it plainly: *if you need to
   weaken a gate or add a plan-format field to make a promotion work, stop and report
   it.* Read the diff for **any gate that got weaker**, any rule that got a special
   case for a promoted recipe, and any new field in the plan format that only the
   promotion needs. A promoted recipe must be guarded exactly like any other.
3. **The version moved with the format.** `core.DomainVersion` goes **4 → 5**,
   `docs/03` § 6's version paragraph says 5, and the recorded fingerprint is updated.
   **Prove the fence still stands, in its designed direction**: the guard refuses a
   write to a plan written under **newer** rules than the tool speaks —
   `internal/core/lock.go` refuses when `plan > DomainVersion`. So a plan stamped
   **6** must be refused by this version-5 tool, and a plan stamped **4** or **5**
   must be accepted. (The fence is one-directional on purpose: an old tool must not
   write a plan it does not understand, but a new tool may write an old plan.)
   > The orchestrator's first brief said "a plan stamped 4 must be refused" — that was
   > wrong, and the runner's report corrected it. Judge against the code, not against
   > the brief.
4. **The id is kind-neutral.** `BL01` stays `BL01` when an idea is shaped into a
   feature. Check that nothing in the model or the store rewrites the id when the kind
   changes.
5. **Scope.** The slice is the backlog. Confirm it did **not** touch the roles, the
   tabs, the recipe ids, the module path, the binaries or `Proof`→evidence (slice D).

### The invariants — check them, do not assume them

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps **28**
  external modules.
  ```bash
  go list -deps -f '{{if not .Standard}}{{.Module.Path}}{{end}}' ./internal/ui | sort -u | grep -v lthiagol/maitre | wc -l
  ```
- The domain layer has **zero UI imports**.
- One fact, one answer: the kind is the only truth; the group is derived.
- The kitchen theme stays; the artifact is still a **recipe**.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.** They go in the
  ledger in a later round.

### The ledger

`docs/audit.md` is **yours** to write. The runner did not touch it — confirm that with
`git diff`. Add any NEW blocker as a new row with an id continuing from **F-73**. Do
not re-file what is closed; the ledger is not swept. **Close F-73** if the plan's name
is fixed. Commit the ledger.

### Hard rules

- Do **not** edit implementation code. Read-only for everything except
  `docs/audit.md`.
- Revert every mutation you make. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-c.md <<EOF
Reviewer: slice C review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-74.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-c"
herdr agent prompt "w1Z:p1" "Reviewer: slice C review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-74.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>."
```
