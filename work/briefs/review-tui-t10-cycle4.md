You are the **reviewer** for **T10 fix cycle 3**. You are a fresh lane. You did not
write the fix. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[fix-tui-t10-3.md](fix-tui-t10-3.md) — the contract for this cycle — and then
[tui-t10.md](tui-t10.md), the contract the slice was built to. The orchestrator's pane
id is `w1H:p2`. Use it verbatim.

The fix commits: `3fb3701` (AC-01, the third copy), `4bf7d93` (AC-02, the
specification), `a8d2bab` (AC-03, the flag), and the row `d0f0f15`. The receipt is
[work/inbox/runner-tui-t10-fix3.md](../inbox/runner-tui-t10-fix3.md).

## This cycle is one finding

File **HIGH or MEDIUM** blockers. Do not re-review T10 from scratch. New findings
start at **F-145**. A LOW note goes in the message, not in the ledger as a finding.

| Finding | The target | The commit |
|---|---|---|
| **F-144** | the format lived in **three** places — the skill's table (wired), the console's parser, and the **startup task's** copy at `internal/ui/souschef.go:288`, the copy the console hands the orchestrator at spawn. Nothing parsed it, and its witness was a **substring check** that held the words and not the shape — so the separator could drift with every word intact and **nothing went red** | `3fb3701` |

Close it in `docs/audit.md` against that commit, and write the slice row for this
cycle — and, if T10's findings are all closed by then, **the row that closes the
slice**.

## What the fix claims

**The third copy is one constant** — `rosterExample` in `internal/ui/souschef.go` — and
it is **wired to the parser** the way the skill's copy is. The wire test extracts
**the startup task's own example** from the rendered task and feeds it to the
console's `parseRoster` and `parseProgress`. A **second** test holds that the
task's copy is **one of the skill's own lines**, not a variant of it.

**Your probe now fails.** The startup task's separator drifts from *" · "* to
*" / "*, every word intact — and **two** tests go red: *"the roster's line names the
job 'orchestrator / round 2 / 6 interactions', which is not one of the crew's"* and
*"the startup task's copy is not the skill's own line — the skill's table is the
specification"*. Reproduce it by hand.

**The specification is said in three places** — the skill's table (in its own
section), the console's constant's doc comment, and the wire's comment — and the
skill's own witness holds the claim. The startup task's doc says its copy is **an
example**.

**The flag names the allowance that ran out**: *"the rounds are spent"* · *"the
interactions are spent"* · *"the rounds and the interactions are spent"*.

## What to verify — each by command, not by reading

| # | Check |
|:--:|---|
| 1 | The startup task's copy is **one constant**, and the rendered task's text is unchanged: `herdr tab rename <TAB_ID> "maitre: orchestrator · round <N> · <M> interactions"` |
| 2 | The wire test feeds **the startup task's own example** to the console's `parseRoster` and `parseProgress`, and a second test holds that it is **one of the skill's own lines** |
| 3 | **Your probe now fails**: the separator drifts from *" · "* to *" / "*, every word intact — and **two** tests go red. Reproduce it by hand |
| 4 | The **specification** is said where a reader will find it — the skill's table, the console's constant, and the wire's comment — and the skill's own witness holds the claim |
| 5 | The flag names the allowance: *"the rounds are spent"* · *"the interactions are spent"* · *"the rounds and the interactions are spent"*. A reader must never have to guess which allowance ran out |
| 6 | `internal/core` did not change. `core.DomainVersion` is **7**. The tab count is **8**. The budget is **28 of 30** |
| 7 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The four notes the runner left for you

| | the claim |
|:--:|---|
| 1 | The startup task's copy is a **constant**, not a retyped literal — the console's side of the format lives in one place: the constant and the parser, joined by the wire to the skill's table |
| 2 | The **equality test** is the "one answer" check: the task's copy must be *one of the skill's own lines*. It is a **second, independent** wire — a drift that still parses but is not the skill's line is caught too |
| 3 | The flag's wording is the runner's; the requirement was the reader's |
| 4 | The seeder and the render path are **untouched**, as ruled. The gate's headroom moved from 523.0s to **529.1s** of the 600s default (**88%**); this cycle's four tests cost about what the ones they replace cost, and the **shared seeded store** is the next slice and the only thing that buys the margin |

## Mutate, run, revert

| | The mutation |
|:--:|---|
| 1 | change the startup task's **separator**, every word intact — **two** tests must go red |
| 2 | make the wire test read a **different** example than the startup task's own |
| 3 | make the equality test compare against a **variant** of the skill's line rather than the line itself |
| 4 | make the flag say *"the budget is spent"* again — the compound wording must go red |

Revert every mutation. `git status --short` at the end is empty except your ledger
commit.

## The questions

Read [work/questions.md](../questions.md). Every row in it is struck — Q11–Q13 and
Q25 closed, Q16–Q30 ruled. Do not reopen them and do not answer them. Say whether
anything was newly raised.

## The scope

| | |
|:--:|---|
| **Implementation code** | do not edit it, except to mutate and revert |
| **`docs/audit.md`** | yours. One commit for the ledger, and the slice row for this cycle |
| **`work/`** | the orchestrator's, except your receipt |
| **The runner** | do not prompt it. Do not spawn agents |
| **The seeder, the render path** | the next slice's work — not yours. Say what the shared seeded store must achieve and what it must not break; the brief for it goes to the chef first |

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
cat > work/inbox/reviewer-tui-t10-cycle4.md <<'EOF'
Reviewer: T10 cycle 4 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-145.. or none>). The third copy: <the constant, and that the reviewer's own probe now turns two tests red>. The specification: <where it is said, and what holds the claim>. The flag: <the three wordings>. The four notes: <each one, and whether it holds>. The next slice: <what the shared seeded store must achieve, and what it must not break>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t10-cycle4.md && git commit -q -m "inbox: reviewer reports tui-t10 cycle 4"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
