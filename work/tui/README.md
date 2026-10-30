# The TUI effort — the state

> **This is the state, not the design.** The design is [plan.md](plan.md). The
> *contract* for one session is a brief in `work/briefs/`, and the brief wins over
> both where they differ.
>
> The sibling [../README.md](../README.md) is the **rename's** record — its history,
> and the protocols that outlived it. This file is the same kind of thing for the
> console's second time round.

## Where it stands

| Session | Feature | What it lands | State | Brief |
|--:|---|---|---|---|
| **T1** | the identity | the window title carries the project; the header's name stops being a footnote; `--project`; the Settings tab shows the plan's path | 🔄 **briefed** | [tui-t1.md](../briefs/tui-t1.md) |
| **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | 🔄 **briefed** | [tui-t2a.md](../briefs/tui-t2a.md) |
| **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable from a script | 🔄 **briefed** | [tui-t2b.md](../briefs/tui-t2b.md) |
| **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | ⏳ planned | — |
| **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |
| **T4** | the Home tab | the landing tab; the art; the Sous-Chef dissolves; D-2 | ⏳ planned | — |
| **T5** | the Ready tab | `core.Adoptable`; the Path becomes Ready; D-1, D-4 | ⏳ planned | — |
| **T6** | the Board, as a visual place | the tree drawn left to right; clicking a node | ⏳ planned | — |
| **T7** | the filter | `/` filters the list | ⏳ planned | — |

**The states**, in the order they come: `planned` → `briefed` → `in progress` →
`implemented` → `reviewed` → `closed`. A session that was interrupted and left a
half-built slice is `in progress`, and its commits are already in the log — read the
last commit's body before redoing anything.

**The briefs that exist** are T1, T2a and T2b. The rest are written **just before**
their session — a brief written today is stale by the time the slice runs.

**The commits** land in the table's own row as they happen: `T1 · 2 commits ·
3f9c1a2, 8b0e441`. A row with commits under it and no `closed` is a session that
did not finish its report.

## The session protocol

**One session, one slice.** The slices are independent enough that they can be done
in any order, and each one leaves the tree green.

| Rule | Why |
|---|---|
| **One slice per session** | two slices in one session means two sets of context, and the second one is the one that breaks |
| **A brief is written before the session** | the brief is the contract; a session that starts without one starts by re-deriving the plan |
| **One commit per acceptance criterion, as it goes** | a commit is a landing place. A turn that ends leaves only what was committed |
| **The state row is updated in the same session** | the row above *is* the resume record. A slice that closes without updating it has not finished |
| **The gate is `make check`** | nobody runs `go test`, `go vet` or `gofmt` by hand |
| **The receipt is written, committed, and *then* the nudge** | the file cannot be lost; the nudge is only the bell |

### The receipt

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t1.md <<EOF
Runner: T1 done. Commits: <n>. Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass. Notes: <anything the next session should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t1"
```

**Then** the nudge, if a lane is waiting for it. The nudge only lands while the
receiver is idle, so the receipt is the half that cannot be lost.

### The scrollback

At every slice boundary, the session's scrollback is captured into
`work/transcripts/`. It is bounded and rotates — one night of lanes held 85 lines of
a much longer conversation when it was first captured.

### The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. If a file you did not
write has changed, **stop and report it**.

## The invariants — these do not move

| | |
|---|---|
| **The domain layer** | `internal/core` has **zero UI imports**, and `internal/core` and `internal/cli` stay **stdlib-only**. `internal/ui` keeps its budget (≤ 30 external modules; the count is in `AGENTS.md`) |
| **`core.DomainVersion` is 7** | do not move it. A format change moves it — and T1–T4 and T7 move **no** format |
| **The tab count is 7** | `Home · Recipes · Backlog · Ready · Board · Log · Settings`. T4 and T5 move names; the count stays |
| **One fact, one answer** | a number the console shows is computed in the domain. A view is registered in `core.DerivedViews` **and** written down in `docs/03` § 8 — the wire checks both directions |
| **The console never authors a recipe** | ADR-0003. A *backlog item* is not a recipe, so the console may create one (T3b) |
| **The console offers no move the gates would refuse** | the footer names the move the *current* selection has |
| **The decoration yields to the content, in both dimensions** | slice G's rule. The mark is drawn only when it fits, in width *and* height |
| **The goldens** | `make goldens` regenerates; `make goldens-check` holds them; they are claimed by a test and they are in the commit |

## The rulings — the chef's answers

| # | Question | Answer | Ruled |
|--:|---|---|---|
| Q1 | the Home tab's name | **Home** | 2026-10-30 |
| Q2 | the Home tab's numbers | **the full set** — plan · recipes · evidence · tasting · backlog · automation | 2026-10-30 |
| Q3 | does the Ready tab show the small tasks? | **both** — the small task's row says *promote → sign off → run*, so nobody reads it as one keystroke | 2026-10-30 |
| Q4 | the filter before the tags? | **the filter now; the tags later**, with the CLI (`DomainVersion` 7 → 8) | 2026-10-30 |
| Q5 | the per-project accent colour | **later** — the window title is the fix, the colour is the extra | 2026-10-30 |
| Q6 | may the console create backlog items? | **yes** — a backlog item is not a recipe, so ADR-0003 does not forbid it | 2026-10-30 |

**What each ruling lands on**, so a brief knows what it inherits:

| Ruling | Lands in |
|--:|---|
| Q1, Q2 | **T4** — the tab is *Home*, it holds the full set of numbers, and the art moves there |
| Q3 | **T5** — the Ready tab carries both blocks, and the small task's row names its two further steps |
| Q4 | **T7** — the filter is the console's answer; the recipe's tags stay a domain decision for the CLI |
| Q5 | **after T1** — not a slice of its own until the window title is in |
| Q6 | **T3b** — the console creates and edits backlog items |

**A ruling lands here, and in the brief it came from** — the plan's *Questions* table
is the source, and this table is the record.

## How to resume

```bash
cd /Users/thiago/Code/Pessoal/maitre
git log --oneline -12                       # where the work stopped
make check                                  # the tree must be green
cat work/tui/README.md                      # ← this file: the state
cat work/tui/plan.md                        # the design, for the reasoning
cat work/briefs/tui-t<N>.md                 # the contract for the next slice
```

**Then:**

1. Find the first session in the table above that is not `closed`.
2. If it is `planned`, it has no brief — **write the brief first**, from the plan's
   section for that slice, into `work/briefs/tui-t<N>.md`, and mark the row `briefed`.
3. If it is `in progress`, read the last commit's body — it says where the session
   stopped — and continue from there. **Do not redo what landed.**
4. If the row is `closed`, the next session is the next row that is not `closed`.

**A cold session needs three things and nothing else:** this file, the plan, and the
brief. Everything else — the code, the docs, the history — is in the repository.

## The briefs

| | |
|---|---|
| **Where** | `work/briefs/tui-t<N>.md` — the briefs live with the other briefs, and they are named so the TUI's briefs are one grep away |
| **When** | a brief is written **just before** the session that runs it, and never for two sessions at once: a brief written today is stale by the time the slice runs |
| **What is in it** | the slice's goal · the acceptance criteria · the invariants · the scope — what **not** to touch · the evidence to produce · the report protocol |
| **Which wins** | the brief is the contract. Where the brief and the plan differ, the brief wins and the plan is corrected |
