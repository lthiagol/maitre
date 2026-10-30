You are the **runner** for **T1** — the identity. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` (the house style and the working
rules) first, then [work/tui/plan.md](../tui/plan.md) § 1 for the reasoning behind
this slice, and [work/tui/README.md](../tui/README.md) for the state and the protocol.

This is the **first** session of the console's second time round. Nothing else is
started until it closes.

## The chef's ruling — the identity

> *we have the small bowtie and maitre, but i miss some indicator of the project
> we're in. Since we can open more than one maitre and work on different projects at
> the same time, this identification would be a good thing*

**The gap is not the header — it is the window title.** The project name is already
in the header (`⋈ maitre · kbtest`), but it is rendered in `styFaint`, and the
**window title** — the row the operating system puts on the terminal's tab, which is
what the chef actually sees with four terminals open — is the hardcoded string
`"maitre"` (`internal/ui/views.go:53`, in `Model.View`).

## Your slice — T1 — the identity

### AC-01 — the window title carries the project

`Model.View` sets `v.WindowTitle`. It must carry the project's name:
`maitre · kbtest`.

- The name is the plan's `project.json` `name` — the same source the header uses
  today (`m.ops.GetProject()`, falling back to the plan's root).
  **Corrected in fix cycle 1 (F-100):** the domain answers the plan's name once —
  `core.Store.ProjectName()` — and the console asks it. The answer is a name, never
  a path.
- **One source for the name.** The header and the window title must not compute it
  separately; there is one function that answers *"what is this plan called?"*, and
  both read it.
- A test asserts `m.View().WindowTitle` — the window title is not in the rendered
  content, so `--render` cannot see it. A test is the only instrument here.

### AC-02 — the header's project name stops being a footnote

The header is `⋈ maitre · <name>`. The name is currently `styFaint`; it must be
rendered in a style that reads as a **name** — the console's accent, or a dedicated
one. The compact mark and the product name keep their styles.

*State the target:* the project name is the project's name, and it reads like one.

### AC-03 — `--project` overrides the name for the session

`maitre --project <name>` shows `<name>` instead of the plan's own name — in the
header **and** in the window title. A chef with two checkouts of one project needs
it, and it is the answer to two plans whose directories share a basename.

- The flag is **display only**. It writes nothing: it does not touch `project.json`,
  it does not enter the journal, and it does not change the plan's name on disk.
- The default is the plan's own name, so the console looks exactly as it does today
  when the flag is absent.
- The flag reaches the model — `New`, `NewLive` and `NewScripted` are the three
  constructors (`internal/ui/model.go:156-189`) — and `--render` and `--keys` honour
  it too, because they build the same model.

### AC-04 — the Settings tab shows the plan's path

The Settings tab is where the console's own facts live. It shows the plan's **root
path** — the absolute path of the plan directory — so the chef can tell two plans
apart when their names are alike.

- The path is the plan's root, formatted for a human (`~/Code/Pessoal/maitre`), not
  a raw absolute path with the home directory spelled out.
  **Corrected in fix cycle 1 (F-101):** the root is the **plan directory**, and the
  example was the project's directory — `~/Code/Pessoal/maitre/kitchen`. The row is
  labelled `root`, and it shows the root.
- The Settings tab currently lists the themes. The path is one more row, labelled.
  **Corrected in fix cycle 1 (F-101):** the tab already carried a `root` row; it
  showed the plan directory's parent.

### AC-05 — the goldens

`make goldens` regenerates, `make goldens-check` holds them.

**The goldens are expected to be unchanged by this slice.** The golden is
`StripANSI(m.View().Content)` — a *style* is not text — and this slice changes the
header's style, not its words; the header's words do not change; and the Settings tab
is not in the golden. **If a golden changed, something else changed the render — find
it and say so in the report.** A golden that moved for no stated reason is the slice
saying it did more than it said.

### AC-06 — the state row is updated

`work/tui/README.md`'s table gains the row's state, the commit ids, and the brief
link — and it is **committed with the slice**. The row is the resume record: a
session that closes without it has not finished.

```
| **T1** | the identity | … | ✅ closed · 3 commits · 3f9c1a2, 8b0e441 | [tui-t1.md](../briefs/tui-t1.md) |
```

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget. **No dependency is added by this slice.**
- **`core.DomainVersion` is 7.** This slice changes no format, so it does not move.
- The tab count is 7 and the tab names do not change.
- **The domain layer has zero UI imports.** If the identity needs a fact the domain
  does not have, the domain gains it — but this slice should not need one.
- The gate is `make check`. A focused check is `make test PKG=… RUN=…`.

## The scope — what not to touch

| | |
|---|---|
| **The CLI** | `maitre-cli` does not change. The console's flag is `maitre`'s, not `maitre-cli`'s |
| **The tabs** | no tab is added, removed or renamed. T4 and T5 own the tab bar |
| **The domain** | no domain change. If you believe one is needed, **stop and report it** — do not widen the slice |
| **The art** | the mark's drawings live in `internal/brand` and nowhere else. This slice moves no art |
| **The docs** | `docs/05-human-interface.md` mentions the console's screen; if this slice makes a line of it false, fix the line and say so. No other doc changes |
| **`docs/audit.md`** | the ledger. **Never swept.** This slice files no finding |
| **`work/`** | the orchestrator's, except the state row (AC-06) |

## The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. If a file you did not
write has changed, **stop and report it**.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t1.md <<EOF
Runner: T1 done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. The window title now reads: <what it reads>. The goldens: <unchanged / changed, and why>. Notes: <anything the next session should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t1"
```

**Write the receipt, commit it, and *then* nudge** — if a lane is waiting for it. The
nudge only reaches an idle receiver; the receipt is the half that cannot be lost.

**One write per file, and commit as soon as the file is coherent.** Every edit to a
large file re-renders that file's whole diff into the terminal, and a turn that
renders dozens of diffs runs out before the work is done.

## Before you report

```bash
make check && git status --short
```

Green, and the tree clean. Report the **real** numbers: the commits, the gate's
result, and what the window title now reads.
