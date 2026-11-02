You are the **reviewer** for **T10**. You are a fresh lane. You did not write this
code. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/tui/plan-automation.md](../tui/plan-automation.md) § 3–§ 6, § 9 and § 14–§ 15
— the design — and then [tui-t10.md](tui-t10.md), the contract. The orchestrator's
pane id is `w1H:p2`. Use it verbatim.

The runner's commits: `f3e85f8` (AC-01, the tab), `5fe4d68` (AC-05, the record's
schema), `0159fa8` (AC-02, the order), `867b545` (AC-03, the queue and the stamp),
`59b0d12` (AC-04, the crew's rows), `c03282e` (AC-06, the moves), `ed0a13a` (AC-07,
the failures), and the row `9b8468b`. The receipt is
[work/inbox/runner-tui-t10.md](../inbox/runner-tui-t10.md).

## This slice

**T10 — the Sous-Chef tab.** The last slice of the automation's work: the tab that
starts the automation, and shows it working.

| | |
|---|---|
| **the tab** | **Sous-Chef**, between Board and Log. **8 tabs** |
| **the ticket** | the **order** (harness · crew · sources · budget), the **queue**, and the **stamp** (progress · vital signs) |
| **the crew** | found by the agents' **cwd**, the state from **herdr**, the job from the tab's **label** |
| **the moves** | `space` starts the crew and `space` stops it; the footer says which |
| **the record** | the settings file's own schema moves **1 → 2** with the `crew` section. The plan's format does not move — `core.DomainVersion` stays **7** |

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | **8 tabs**, Sous-Chef between Board and Log. The count moved in `AGENTS.md`, in `docs/05` and in the console's own tab list, and the goldens moved with it |
| 2 | the **order**: four controls, read from the record and written when the chef changes a pick. The harness shows **the command it will run**; the crew states **who reviews** — and "1 · one agent, and **the chef reviews**"; the sources' counts are live from the plan; the budget's presets are 2·8, 3·12, 5·24 |
| 3 | the **queue**: the plan's own work under the source pick, up to four rows with the rest counted, and the first marked `next`. The **stamp**: the progress, and the **vital signs** — `plan 7 · build 7 · they agree`, or `they disagree` in the warning style |
| 4 | the **crew**: found by the agents' cwd, the state read from `agent_status`, the job taken from the tab's **label** (`maitre: runner · RCP03`). A tab whose label does not carry the roster's mark is **not** a member — so the chef's own pane is not the crew's to manage. A **blocked** agent's row reads `asked` and says the orchestrator has the question, and the console gains **no move** for it |
| 5 | the **record**: the schema moves 1 → 2, the `crew` section joins `view`, and a **version 1 file is refused by name** — and the console **does not write it either** (a write would replace the chef's record with the defaults) |
| 6 | the **moves**: `space` starts and stops; the footer names which. Starting opens a tab on the plan's root — the chef's own workspace when inside a herdr session, a workspace made for the plan when not — and runs the order's **harness command** with the **startup task**. The console **remembers nothing**: no pane id, no file. Stopping closes the crew's own tabs, and only those |
| 7 | the **failures**: `herdr` not on the PATH → *"herdr is not here"* and **nothing** offered; the spawn fails → it says so and **the move stays**; the crew is gone → not running and the move comes back; the orchestrator gone but the crew not → it says so; the plan's root moved → no crew at the new root |
| 8 | `internal/core`, `internal/cli` and the skills did not change. `core.DomainVersion` is **7**. No new view. `docs/03`, `docs/04` and `docs/06` did not move. The budget is **28 of 30** — herdr is a runtime dependency, not a module |
| 9 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The eight notes the runner left for you

Each is a claim about the code or a decision it made that the design did not settle.

| | the claim |
|--:|---|
| 1 | **The Home block's stale claim was fixed.** It said *"the automation is not built"*; it is built now, so it says *"the automation is started from the Sous-Chef tab"* — the vital signs alone would read as a running automation. And `scripts/check-slice-numbers.sh` now reads the doc's **state word** as well as the number (`planned` or `built`). **Does the check still hold the number, and does it still refuse to pass having found nothing?** |
| 2 | **The progress's derivation is the runner's own.** The design does not say where the round and the interactions come from: the **interactions** are the journal's events the chef did not write, and the **rounds** are the recipes the crew has worked. Say whether you read it differently |
| 3 | **`next` is always shown** on the queue's first item; the mockup shows it only while the crew is running |
| 4 | **The harness runs in the pane** — the console runs the order's harness command, not a fixed `mcode-herdr` |
| 5 | **A member with no roster label is not shown.** The console cannot name its job, and the chef's own pane must not be mistaken for the crew |
| 6 | **The console asks the runtime when the tab comes on screen, and on every poll while it is visible** |
| 7 | **The filter does not reach the tab** (`filterApplies` is unchanged), and **`←`/`→` pick on this tab instead of switching tabs.** Right is a tab-switching key in this console — check what the console does on the other seven tabs, and say whether picking with `→` on one tab only is right |
| 8 | The row names the seven commits; the row's own commit is not among them |

## The gate's own headroom — a risk, not a finding

**The `internal/ui` suite took 488s of go test's 600s default timeout** under
`-race -shuffle=on -count=2`, and the new tests added roughly 50s of that. The
runner did not raise the timeout, because the gate is the gate.

**Check it yourself**: time the suite. And say which of the two the repo should
take — **share a seeded store across the read-only tests** (a test that seeds a
plan costs about a second) or **raise the timeout**. This is a risk to the gate's
reliability, and a reviewer who sees the number can say whether the next slice
breaks it.

## Mutate, run, revert

Break each rule and report the tests that went red:

| | The mutation |
|--:|---|
| 1 | make the tab's order show a crew count without saying **who reviews** |
| 2 | make the stamp's vital signs always say *they agree* |
| 3 | make the console treat a tab **without** the roster's label as a crew member |
| 4 | make the console **write** a settings file it refused to read |
| 5 | make the console's stop close **every** tab in the workspace, not the crew's |
| 6 | make `herdr` missing offer the move anyway |

Revert every mutation. `git status --short` at the end is empty except your ledger
commit.

## The questions

Read [work/questions.md](../questions.md). Q11–Q13 and Q25 are **closed**, and
Q16–Q29 are **ruled** — do not reopen them and do not answer them. Say whether
anything was newly raised.

## The scope

| | |
|---|---|
| **Implementation code** | do not edit it, except to mutate and revert |
| **`docs/audit.md`** | yours. One commit for the ledger, and the slice row for this cycle |
| **`work/`** | the orchestrator's, except your receipt |
| **The runner** | do not prompt it. Do not spawn agents |
| **`docs/03`, `docs/04`, `docs/06`** | not yours. Report a doc that contradicts the code; do not fix it |

## The verdict

| | |
|---|---|
| No HIGH or MEDIUM findings | `verdict: ok` |
| Any HIGH or MEDIUM finding | `verdict: changes-needed` |

Do not invent reasons to fail.

## Notify the orchestrator

The receipt comes first.

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-tui-t10.md <<'EOF'
Reviewer: T10 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-140.. or none>). The tab: <the name, the place, the count>. The order: <the four controls, and who reviews>. The queue and the stamp: <what each says, and the vital signs>. The crew: <how it is found, the states, the labels, and the blocked row>. The record: <the schema's version, the refusal, and that it is not written>. The moves: <the start, the startup task, the stop, and the footer>. The failures: <each one, and what the tab says>. The eight notes: <each one, and whether it holds>. The gate's headroom: <the suite's time, and which of the two the repo should take>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t10.md && git commit -q -m "inbox: reviewer reports tui-t10"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
