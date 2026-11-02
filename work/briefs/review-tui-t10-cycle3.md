You are the **reviewer** for **T10 fix cycle 2**. You are a fresh lane. You did not
write the fix. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[fix-tui-t10-2.md](fix-tui-t10-2.md) — the contract for this cycle — and then
[tui-t10.md](tui-t10.md), the contract the slice was built to. The orchestrator's pane
id is `w1H:p2`. Use it verbatim.

The fix commits: `7edc649` (AC-01, the label), `9003efb` (AC-02, the wire),
`7d58dc8` (AC-03, the stamp), `7597466` (AC-04, the pane's own tab id), and the row
`99412f6`. The receipt is
[work/inbox/runner-tui-t10-fix2.md](../inbox/runner-tui-t10-fix2.md).

## This cycle is two findings

File **HIGH or MEDIUM** blockers. Do not re-review T10 from scratch. New findings
start at **F-144**. A LOW note goes in the message, not in the ledger as a finding.

| Finding | The target | The commit |
|---|---|---|
| **F-142** | the crew's label is **agent-written text** and reached the screen **unsanitised** — a newline put a line of the agent's text where the crew's rows are a table, and an ANSI escape was in the raw render. `sanitizePlan` covers the plan and not the runtime's label | `7edc649` |
| **F-143** | the label's format is **two copies** — the skill's example and the console's parser — and nothing checked that they agreed. The skill could drift, its own witness be updated, and **both packages stay green** | `9003efb` |

Close both in `docs/audit.md` against those commits, and write the slice row for
this cycle.

## What the fix claims

**The label** is sanitised **where it enters** — in `readCrew`, in
`internal/ui/crew.go` — and the witness is on **this** tab:
`TestTheCrewLabelIsSanitisedOnTheSousChefTab` injects a label with a newline and
an escape, renders the tab, and holds that no control character reaches the
render, that the escape is not in the raw render, and that the injected text does
not become a line of its own.

**The wire** is the **skill's own example**: `internal/ui/crew_wire_test.go` reads
`../cli/skills/orchestrator-role/SKILL.md` by path, extracts every roster's line
the skill shows, fills the placeholders, and requires each one to parse.

**The stamp** reads in three shapes — a progress reported, **no** progress, and a
progress **over** the cap — and the flag fires when the report **reaches** the cap.
A **zeroed** report is a report, not an absence: the console carries a **reported**
flag beside each number.

**The pane's own tab id** comes from `herdr pane current`, not `herdr pane list`.

## What to verify — each by command, not by reading

| # | Check |
|:--:|---|
| 1 | The label is sanitised **where it enters**, in `readCrew`. A **newline** or an **ANSI escape** in a label cannot reach the render, and the injected text does not become a line of its own. Taking the `sanitize` out turns the witness red |
| 2 | The wire: the skill's **own example**, read from the skill's file, is fed to the console's parser. Every roster's line the skill shows must parse — a roster's line, a job of the crew's, and the orchestrator's round and interaction count |
| 3 | **Your probe now fails.** The skill drifts to *"round #<N> · <M> interactions"*, the skill's own witness is updated — `internal/cli` stays green and **`internal/ui` goes red**: *"the console reads no round from the skill's example"*. Reproduce it by hand |
| 4 | The stamp in all three shapes: a progress reported, **no** progress (the separator is there), and a progress **over** the cap (the numbers as reported, never *"7 of 3"*) |
| 5 | The over-budget flag fires when the report **reaches** the cap — *round 3 of 3 · 12 of 12 interactions · the budget is spent* — and not one round later |
| 6 | A **zeroed** report reads as a zero, and an **absence** stays an absence. The console carries a **reported** flag beside each number |
| 7 | The startup task and `orchestrator-role` point at **`herdr pane current`**, and the lanes' tab ids are still in `herdr tab list` |
| 8 | `internal/core` did not change. `core.DomainVersion` is **7**. The tab count is **8**. The budget is **28 of 30** |
| 9 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The six notes the runner left for you

| | the claim |
|:--:|---|
| 1 | **The wire lives in `internal/ui` and reads the skill's file by path** — the console does not import the CLI; the path *is* the wire, and a skill that moves breaks the test loudly rather than silently. Is a path a wire, or a coupling that will rot? |
| 2 | **The zeroed report's shape is the runner's**: `round 0 of 3 · 0 of 12 interactions`. A round is 1-based, so a *reported* zero round means *"no round yet"*; the console prints what was reported and does not read it as an absence. Is that right? |
| 3 | **The over-budget wording is unchanged** — *"the budget is spent"* — and it now fires **at** the cap |
| 4 | **Two notes were not fixed and are claimed for the next slice**: the render path still waits on a process, and the poll still spends two process spawns a second while the tab is visible |
| 5 | **The gate's headroom did not move**: `internal/ui` is **523.0s** of go test's 600s default under `-race -shuffle=on -count=2` — **87%**. The next slice of this size lands at ~90%, and the one after it breaks the gate |
| 6 | The five commits are one per criterion and each is green |

## The next slice — the seeded store

**The clock is the reason the next slice is not a console slice.** The reviewer
wrote its acceptance criteria out, and the brief for it goes to the chef before
the runner sees it. Say whether you would hold the same line:

| | |
|:--:|---|
| **it must achieve** | the seeder runs **once** and every test gets a **copy**; `internal/ui` falls from 522s to **well under 200s** under the race gate — a margin of at least **3×** |
| **isolation** | each test its own copy, or the writers see each other under `-shuffle=on -count=2` |
| **the goldens** | the project's name comes from the **parent directory's** name, so the copy must land in a directory named `demo` or `console-home.golden` moves |
| **the stress plan** | a **second** template, not the same one |
| **the cleanup** | `TestMain`, not a user's tmp |
| **timestamps** | no test asserts a fresh one — you checked, and none does |

## Mutate, run, revert

Break each rule and report the tests that went red:

| | The mutation |
|:--:|---|
| 1 | take the `sanitize` out of `readCrew` |
| 2 | make the wire read a **different** example than the skill's own |
| 3 | make the flag wait for the round **after** the cap |
| 4 | make a zeroed report read as an absence |

Revert every mutation. `git status --short` at the end is empty except your ledger
commit.

## The questions

Read [work/questions.md](../questions.md). Q11–Q13 and Q25 are **closed**, and
Q16–Q30 are **ruled** — do not reopen them and do not answer them. Say whether
anything was newly raised.

## The scope

| | |
|:--:|---|
| **Implementation code** | do not edit it, except to mutate and revert |
| **`docs/audit.md`** | yours. One commit for the ledger, and the slice row for this cycle |
| **`work/`** | the orchestrator's, except your receipt |
| **The runner** | do not prompt it. Do not spawn agents |
| **The seeded store** | the next slice's work. Say what it must achieve, and what it must not break — do not start it |

## The verdict

| | |
|:--:|---|
| No HIGH or MEDIUM findings | `verdict: ok` |
| Any HIGH or MEDIUM finding | `verdict: changes-needed` |

Do not invent reasons to fail.

## Notify the orchestrator

The receipt comes first.

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-tui-t10-cycle3.md <<'EOF'
Reviewer: T10 cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-144.. or none>). The label: <where it is sanitised, and the witness on the tab>. One format: <the wire, and that the reviewer's own probe now fails>. The stamp: <the three shapes, the flag at the cap, and the zeroed report>. The pane's own tab id: <what it now says>. The notes: <each one, and whether it holds>. The next slice: <what the seeded store must achieve, and what it must not break>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t10-cycle3.md && git commit -q -m "inbox: reviewer reports tui-t10 cycle 3"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
