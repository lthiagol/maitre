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
| **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | ✅ closed · 9 commits · 4dc8079, 1b550c6, 1804ccc, b873821, 5edae71, aa4ba5f, da36146, 99a7576, 63964ce | [tui-t3a.md](../briefs/tui-t3a.md) · [fix 1](../briefs/fix-tui-t3a-1.md) · [fix 2](../briefs/fix-tui-t3a-2.md) |
| **T3b** | the backlog moves | the domain answers the item's move; the console shapes · promotes · dismisses · creates | ✅ closed · 8 commits · 1297cdd, 35688a6, 4799b3b, 7de9a16, 4e14f19, 58069e7, c0f1624, 89eaecc | [tui-t3b.md](../briefs/tui-t3b.md) · [fix 1](../briefs/fix-tui-t3b-1.md) · [fix 2](../briefs/fix-tui-t2b-2.md) |
| **T4** | the Home tab | the landing tab with numbers and no items; the art moves there; the Sous-Chef dissolves into it; D-2 lands | ✅ closed · cycle 2 ok · 8 commits · 219b331, 79c7fb5, 36a1a35, dcaac47, 446fd27, a29935c, 7488e02, a75d867 | [tui-t4.md](../briefs/tui-t4.md) · [fix 1](../briefs/fix-tui-t4-1.md) · [cycle 2](../briefs/review-tui-t4-cycle2.md) |
| **Q9** | the backlog's lists, and the selection | empty `scope` and `tags` are `[]`; after a write the selection follows the item | ✅ closed · 2 commits · 6714558, cfb94aa | [tui-q9.md](../briefs/tui-q9.md) |
| **T5** | the Ready tab | `core.Adoptable`; the Path becomes Ready; D-1, D-4 | ✅ closed · 3 commits · e265382, 93153c2, 9ae7789 | [tui-t5.md](../briefs/tui-t5.md) · [review](../briefs/review-tui-t5.md) |
| **T6** | the Board, as a visual place | the tree drawn left to right; clicking a node | ✅ closed · cycle 2 ok · 5 commits · 5713472, 92413be, cdfdc16, a5df3ce, df6d387 | [tui-t6.md](../briefs/tui-t6.md) · [fix 1](../briefs/fix-tui-t6-1.md) · [cycle 2](../briefs/review-tui-t6-cycle2.md) |
| **T7** | the filter | `/` filters the list | ✅ closed · cycle 2 ok, 0 blockers · 13 commits · 7a76e66, 0f264ac, 0dce0c1, c542916, 5ab1660, 1983b1a, af9deb1, 65fdbba, 675b4cc, def86f3, 00c234e, 8c505d7, ad453df | [tui-t7.md](../briefs/tui-t7.md) · [fix 1](../briefs/fix-tui-t7-1.md) · [cycle 2](../briefs/review-tui-t7-cycle2.md) |
| **T8** | the skills | `maitre-cli skill help` · `list` · `load <name>`; the `orchestrator-role` skill, embedded in the binary | ✅ closed · cycle 2 ok · 11 commits · 8ab3b80, ac3f6d5, 2e9f3a6, 5e2146c, 3333df7, 504cb86, 706560c, dc4f4e4, 6d48c1f, 1801999, 27483df | [tui-t8.md](../briefs/tui-t8.md) |
| **T9** | the settings record | `kitchen/settings.json` — the tool's file, the console's VIEW preferences; `maitre-cli settings show` | ✅ closed · cycle 3 ok · 10 commits · 38d1d48, 2aba9f4, 089ef6b, 1ed7f3f, 2d65fa0, 17cdefe, b40eedf, 97f1de9, 2fb29f2, b170d0c | [tui-t9.md](../briefs/tui-t9.md) · [fix 1](../briefs/fix-tui-t9-1.md) · [fix 2](../briefs/fix-tui-t9-2.md) |
| **T10** | the Sous-Chef tab | the ticket: the order, the queue, and the stamp; `space` starts and stops the crew | ✅ implemented · fix 1 · 11 commits · f3e85f8, 5fe4d68, 0159fa8, 867b545, 59b0d12, c03282e, ed0a13a, 2684307, 1c76e84, f59f9be, 5b03f48 | [tui-t10.md](../briefs/tui-t10.md) · [fix 1](../briefs/fix-tui-t10-1.md) |

**The states**, in the order they come: `planned` → `briefed` → `in progress` →
`implemented` → `reviewed` → `closed`. A session that was interrupted and left a
half-built slice is `in progress`, and its commits are already in the log — read the
last commit's body before redoing anything.

**The briefs that exist** are T1 through T6, and Q9. The rest are written **just
before** their session — a brief written today is stale by the time the slice runs.

**Paused here.** T7 is **closed** — cycle 2 ok, 0 blockers. The console's second round
is over: T1 through T7 are all closed.

| | |
|---|---|
| **T6** | closed. Cycle 2 ok, 0 blockers. Receipt: `work/inbox/reviewer-tui-t6-cycle2.md` (`7890f27`) |
| **T7** | closed. The review filed two blockers and one low; the fix closed all three; cycle 2 said **ok**, 0 blockers (`d8a6a08`, ledger `61a910b`) |
| **The fix** | 5 commits (`ad7be4d`). **F-128** the filter drops the *rows*, not the plan they are read from — the `●` survives. **F-129** a message that does not fit takes a row of its own, the way a refusal does. **F-130** `N` is the rows the list draws. The blind *stay-on-the-item* test filters a strict subset now |
| **The round** | T1–T7 are closed. `plan.md` is the design this round followed; the state is this table |
| **T8** | closed. Cycle 2 ok, 0 blockers. The menu, the embedded skill, the refusal — and the six findings of cycle 1 closed, each with a witness of its own. Receipt: `work/inbox/reviewer-tui-t8-cycle2.md` (`048c710`) |
| **T9** | closed. Cycle 3 ok, 0 blockers. The record, the verb, and the console's VIEW preferences surviving a restart — through three cycles of findings, each one a claim nothing held. Receipt: `work/inbox/reviewer-tui-t9-cycle3.md` (`a0b78b8`) |
| **The next work** | **T10** — the Sous-Chef tab, the last slice of the automation's work. Designed in [plan-automation.md](plan-automation.md) § 14–§ 15 |
| **The rulings** | the filter covers **Recipes, Backlog and Ready** (Q14); the filter key lives in **Settings' key list**, not the footer (Q15). Both ruled 2026-10-31 |

Two notes from T6's first review, neither a finding: the note-cursor test's space-triage half does not run on the Board, and the wide board takes the full width so a fork has a column. `docs/07`'s order-tab row says **Ready**. The view and `maitre-cli path` stay `path`.

T4's cycle 2 left two notes, neither a finding: a recipes row with a two-digit count in every section is 84 columns and wraps at 80, and `docs/05`'s Home sentence is held by no test. The free edit of a backlog item stays unbuilt (Q8). Q11, Q12 and Q13 are still open, in [../questions.md](../questions.md), and they are the orchestrator's.

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
| **The receipt is written, committed, and *then* the nudge** | the file cannot be lost; the nudge is only the bell. The nudge is `herdr pane send-text <pane> "<receipt>"` + `herdr pane send-keys <pane> enter` — `herdr agent prompt` refuses a lane that reports through the `mcode-herdr` bridge |
| **The reviewer reads [work/questions.md](../questions.md) at every slice's close** | and says whether a question was **answered, missed, or newly raised**. A question that is not written there is a question that is lost |

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

# the lanes are mcode-herdr agents: mcode, with the bridge that reports the
# pane's state to herdr. `herdr agent start` has no mcode kind — its kinds are
# compiled in — so the lane is launched in the pane instead.
herdr pane run "$ROOT_PANE"  mcode-herdr
herdr pane run "$SPLIT_PANE" mcode-herdr
```

| | |
|---|---|
| **The tab** | `--workspace "$HERDR_WORKSPACE_ID"` — the tab belongs beside the orchestrator. The default is the **focused** workspace, which is the user's |
| **The split** | the runner on the left, the reviewer on the right. Neither `herdr pane run` nor the launch steals the chef's focus |
| **The lane** | `mcode-herdr` — mcode plus the bridge that reports `working` / `idle` / `blocked` for the pane. A lane that waits goes **idle**, so herdr shows the wait instead of a spinner that never ends |
| **The reviewer** | briefed **after** the runner's work lands, so it reads the code fresh. A reviewer briefed in advance would form its model of the code before the code existed |
| **The dispatch** | the brief plus the notify block with the orchestrator's **real pane id** — `$HERDR_PANE_ID`, not a placeholder. `herdr pane send-text <pane> "<text>"`, then `herdr pane send-keys <pane> enter` — `herdr agent prompt` answers `agent_not_ready` for a bridged lane |
| **The delivery check** | one `herdr pane read --lines 40`, ~5 s after each dispatch. Activity lines mean delivered; the raw prompt above the spinner means a stuck draft — report it and ask the chef to press Enter, **never blind-resend** |
| **The wait** | `herdr agent wait <pane> --until idle --timeout 600000`. The lane going **idle** is the signal that its turn is over — the orchestrator waits, it does not poll |

**T7 is in flight.** Its lanes are `mcode-herdr` agents, so herdr sees each lane's state
— a waiting lane goes **idle**, and the receipt is what carries the report back.

| | |
|---|---|
| the **tab** | `w1H:tK` — "tui: T7", in the orchestrator's workspace `w1H`. The lanes are done and idle; close the tab when the next slice's tab opens |
| the **runner** | `w1H:pZ` — `mcode-herdr`, the left pane. Idle |
| the **reviewer** | `w1H:p0` — `mcode-herdr`, the right pane. Idle |

A tab lives until it is closed, so these survive a session that ends. If a lane is
gone, the receipt in `work/inbox/` is the record of what it did, and the slice's
commits are in the log.

### The tree is shared

Stage **by name** — never `git add -A`, never `git add .`. If a file you did not
write has changed, **stop and report it**.

**And the change is often another lane's mutation, in flight.** The reviewer mutates
the code to prove its tests bite, and the mutation is in the working tree while it
runs. So a lane that finds a stranger's edit:

| | |
|---|---|
| **does not touch the file** | the lane that owns it will revert it, and a second revert can lose the experiment |
| **reports it** | the receipt says what changed, what it breaks, and what the change looks like |
| **says whose it probably is** | a mutation that disables a rule the slice just landed is the reviewer's experiment |
| **stands by** | the orchestrator answers, and the answer is usually *wait* |

A lane that reverts a stranger's edit breaks that lane's experiment and hides the
evidence. A lane that waits for the tree to be clean is doing the right thing.

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
| Q8 | the free edit of a backlog item's fields | **no, for now** — shape, promote, dismiss and create cover the moves | 2026-10-31 |
| Q9 | a backlog item's empty `scope` and `tags` | **`[]`** — one rule for every record. `DomainVersion` stays 7 | 2026-10-31 |
| Q10 | after a write, does the selection follow the item? | **yes** — the selection follows the item. The fix belongs in the model | 2026-10-31 |
| Q14 | which tabs the filter covers | **Recipes, Backlog and Ready** — the plan-item lists. The Board is a tree (a hidden node misstates what waits on what); the Log holds events; Home and Settings hold numbers | 2026-10-31 |
| Q11 | one marker or two for *the domain's answer*? | **two** — the recipe's arrow *asks*, the item's *states*. Closed 2026-11-01 |
| Q12 | does the unbracketed position's message name its remedy? | **yes** — one line; each cause names its own remedy. Closed 2026-11-01 |
| Q13 | the kinds list is built on every call | **later** — a per-frame allocation on a path that redraws the whole screen. Closed 2026-11-01; T7's count has the same shape |
| Q26 | the Sous-Chef tab's design | **the ticket** — an order (harness · crew · sources · budget), the queue under it, and a stamp at the foot | 2026-11-01 |
| Q27 | can the chef stop the automation from the tab? | **yes** — `space` stops the crew, and the footer names which key does what. It revises Q22 | 2026-11-01 |
| Q28 | where do the console's settings live? | **in the plan's folder, in a new file** — `kitchen/settings.json` | 2026-11-01 |
| Q29 | who owns that file — the domain, or the tool? | **the tool.** The settings are not a plan *record* and carry no rules, so the plan's format does not move: `DomainVersion` stays **7**, no stamp | 2026-11-01 |
| Q15 | where the console says the filter exists | **the Settings tab's key list** — the footer names the move the *current selection* has, and `/` opens a mode, not a move. The footer's goldens do not move | 2026-10-31 |
| Q16 | the orchestrator's tab's name | **Sous-Chef** — the tab's old name. The rows carry the glossary's job names (`orchestrator`, `runner`, `reviewer`), so the tab's name and the glossary's sous-chef stay two different things | 2026-10-31 |
| Q17 | the orchestrator's tab: the count and the place | **8 tabs — Sous-Chef between Board and Log**. Home · Recipes · Backlog · Ready · Board · **Sous-Chef** · Log · Settings | 2026-10-31 |
| Q18 | does the Sous-Chef tab offer the chef a move? | **it starts the automation** — `space` spawns the orchestrator. A blocked agent raises a message to the **orchestrator**, asking for directions; the orchestrator decides | 2026-10-31 |
| Q19 | where does the console read the automation's state from? | **herdr** — the console asks the runtime, and finds the crew by the agents' cwd. The plan carries no crew state, so `DomainVersion` stays 7 | 2026-10-31 |
| Q20 | does starting the automation leave a trace in the plan's journal? | **no** — the journal records what happened to the plan; the crew's state is runtime state | 2026-10-31 |
| Q21 | what does the startup task say? | **it loads the role** — `maitre-cli skill load orchestrator-role`, plus the work. The skill carries the role's knowledge; the startup task carries the work | 2026-10-31 |
| Q22 | can the chef stop the automation from the tab? | **no** — one move, start. Stopping is the chef's terminal business | 2026-10-31 |
| Q23 | where do the skills live? | **embedded in the binary** — so each version of the tool carries the right version of the skill | 2026-10-31 |
| Q24 | one slice, two, or three? | **three** — T8 the skills, T9 the settings record, T10 the Sous-Chef tab | 2026-10-31 |

**What each ruling lands on**, so a brief knows what it inherits:

| Ruling | Lands in |
|--:|---|
| Q1, Q2 | **T4** — the tab is *Home*, it holds the full set of numbers, and the art moves there |
| Q3 | **T5** — the Ready tab carries both blocks, and the small task's row names its two further steps |
| Q4 | **T7** — the filter is the console's answer; the recipe's tags stay a domain decision for the CLI |
| Q5 | **after T1** — not a slice of its own until the window title is in |
| Q6 | **T3b** — the console creates backlog items |
| Q8 | **not this round** — the free edit stays unbuilt |
| Q9 | **the backlog writer** — an empty `scope` or `tags` is `[]`. `DomainVersion` stays 7 |
| Q10 | **the model** — after a write, the selection follows the item |
| Q14 | **T7** — the filter covers the plan-item lists, and the other tabs say so |
| Q11, Q12, Q13 | **closed** — answered by what the console now does |
| Q26 | **T10** — the tab is the ticket |
| Q27 | **T10** — the tab offers a stop, and says which key does what |
| Q28, Q29 | **T9** — the settings are the tool's file, `kitchen/settings.json`; the plan's format does not move. **T10** adds the `crew` section |
| Q15 | **T7** — nothing to change. The footer stays the selection's move; the key is in Settings |
| Q16 | **T10** — the tab's name |
| Q17 | **T10** — 8 tabs, Sous-Chef between Board and Log |
| Q18 | **T10** — the tab's move is to start the automation; a blocked agent asks the orchestrator |
| Q19 | **T10** — the state comes from herdr; the plan carries no crew state |
| Q20 | **T10** — the journal records the plan, not the runtime |
| Q21 | **T10** — the startup task loads the role |
| Q22 | **T10** — one move: start. Superseded by Q27 |
| Q23 | **T8** — the skills are embedded in the binary |
| Q24 | **T8**, **T9** and **T10** — three slices, the skills first |

**A ruling lands here, and in the brief it came from** — the plan's *Questions* table
is the source, and this table is the record.

**A question that is still open is not here** — it is in
[work/questions.md](../questions.md), with its options and the orchestrator's
recommendation. It moves here when it is answered.

## How to resume

```bash
cd /Users/thiago/Code/Pessoal/maitre
git log --oneline -12                       # where the work stopped
make check                                  # the tree must be green
cat work/tui/README.md                      # ← this file: the state
cat work/tui/plan.md                        # the design, for the reasoning
cat work/briefs/tui-t<N>.md                 # the contract for the next slice
cat work/questions.md                       # ← the questions that wait on the chef
cat work/transcripts/tui-t<N>-*.md           # the lanes' scrollback for the slice
sed -n '/## Open findings/,/## Closed findings/p' docs/audit.md   # the ledger
```

**Then:**

1. Find the first session in the table above that is not `closed`. That is **T7**.
2. If it is `planned`, it has no brief — **write the brief first**, from the plan's
   section for that slice, into `work/briefs/tui-t<N>.md`, and mark the row `briefed`.
3. If it is `in progress`, read the last commit's body — it says where the session
   stopped — and continue from there. **Do not redo what landed.**
4. If the row is `closed`, the next session is the next row that is not `closed`.
5. If a slice was **in flight** when the session ended, its lanes may still be alive —
   their tab and pane ids are in *The lanes* above. Read `work/inbox/` first. The
   receipt may already be there.
6. **The tree may be dirty, and that is by design.** A lane commits per acceptance
   criterion, so its work in flight is uncommitted — `git status` shows another
   lane's hands, not rubbish. **Check `git status` and the receipt before you touch
   anything**, and never `git checkout` a file you did not write.

**A cold session needs four things and nothing else:** this file, the plan, the
brief, and **[work/questions.md](../questions.md)** — the questions that wait on the
chef, with their options and the orchestrator's recommendation. Everything else — the
code, the docs, the history — is in the repository.

## The briefs

| | |
|---|---|
| **Where** | `work/briefs/tui-t<N>.md` — the briefs live with the other briefs, and they are named so the TUI's briefs are one grep away |
| **When** | a brief is written **just before** the session that runs it, and never for two sessions at once: a brief written today is stale by the time the slice runs |
| **What is in it** | the slice's goal · the acceptance criteria · the invariants · the scope — what **not** to touch · the evidence to produce · the report protocol |
| **Which wins** | the brief is the contract. Where the brief and the plan differ, the brief wins and the plan is corrected |
