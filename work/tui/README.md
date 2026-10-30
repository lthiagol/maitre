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
| **T1** | the identity | the window title carries the project; the header's name stops being a footnote; `--project`; the Settings tab shows the plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438, f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81, 8a6866b, cf8c05e | [tui-t1.md](../briefs/tui-t1.md) · [fix 1](../briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |
| **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | ✅ closed · 12 commits · a72ea98, 115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](../briefs/tui-t2a.md) · [fix 1](../briefs/fix-tui-t2a-1.md) · [fix 2](../briefs/fix-tui-t2a-2.md) |
| **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable from a script | ✅ closed · 9 commits · 4dfd5a4, 187d706, f64e6a2, ac8df77, 588cf5f, 9c716ce, 0144216, 6bb7545, f8969ef | [tui-t2b.md](../briefs/tui-t2b.md) · [fix 1](../briefs/fix-tui-t2b-1.md) |
| **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | 🔄 **briefed** | [tui-t3a.md](../briefs/tui-t3a.md) |
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
did not finish its report. **The count is the number of ids the row names**, and the
commits that wrote the row are not among them — a commit cannot carry its own hash.

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

### The lanes

One tab per slice, in **the orchestrator's own workspace**, split for the runner and
the reviewer. **Both lanes are fresh for every slice.**

```bash
# the tab goes in the ORCHESTRATOR's workspace.
# Without --workspace, herdr lands the tab in the FOCUSED one — the user's,
# and often another project. T1's tab opened in training-prod, not in maitre.
RAW=$(herdr tab create --workspace "$HERDR_WORKSPACE_ID" --label "tui: T<N>" --cwd "$PWD" --no-focus)
ROOT_PANE=$(echo "$RAW" | python3 -c "import json,sys; print(json.load(sys.stdin)['result']['root_pane']['pane_id'])")
SPLIT_RAW=$(herdr pane split "$ROOT_PANE" --direction right --ratio 0.5 --cwd "$PWD")
SPLIT_PANE=$(echo "$SPLIT_RAW" | python3 -c "import json,sys; print(json.load(sys.stdin)['result']['pane']['pane_id'])")

herdr agent start runner   --kind opencode --pane "$ROOT_PANE"  --timeout 60000
herdr agent start reviewer --kind opencode --pane "$SPLIT_PANE" --timeout 60000
```

| | |
|---|---|
| **The tab** | `--workspace "$HERDR_WORKSPACE_ID"` — the tab belongs beside the orchestrator. The default is the **focused** workspace, which is the user's |
| **The split** | the runner on the left, the reviewer on the right. `herdr agent start` has **no `--no-focus`** — it pulls the chef's UI to the tab, once per lane |
| **The reviewer** | briefed **after** the runner's work lands, so it reads the code fresh. A reviewer briefed in advance would form its model of the code before the code existed |
| **The dispatch** | the brief plus the notify block with the orchestrator's **real pane id** — `$HERDR_PANE_ID`, not a placeholder |
| **The delivery check** | one `herdr pane read --lines 40`, ~5 s after each dispatch. Activity lines mean delivered; the raw prompt above the spinner means a stuck draft — report it and ask the chef to press Enter, **never blind-resend** |

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
| Q7 | where does the rule that answers *"what is this plan called?"* live in the spec? | **`docs/03` § 6, beside the plan format** — the rule is about what the reader answers from `project.json`, and § 6 is where the reader looks. It is **not** in § 8: the name is a *stored* fact, not a derived view. A wire holds it | 2026-10-30 |

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
