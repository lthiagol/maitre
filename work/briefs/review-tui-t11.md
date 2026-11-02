You are the **reviewer** for **T11**. You are a fresh lane. You did not write this
code. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[tui-t11.md](tui-t11.md) — the contract. The orchestrator's pane id is `w1H:p2`. Use
it verbatim.

The runner's commits: `825bb82` (the shared store), `90774a4` (the copy is named
after its directory), `f090206` (the store survives `-count=2`), `2fd2c13` (the
wire — no test seeds its own plan), and the row `e9df136`. The receipt is
[work/inbox/runner-tui-t11.md](../inbox/runner-tui-t11.md).

## This slice

**T11 — the shared seeded store.** It changes **no behaviour**. It makes the suite
pay for the plan it builds **once** instead of two hundred and forty times.

| | |
|---|---|
| **the store** | `internal/ui/seed_test.go` — `TestMain` seeds the templates **once** into a suite-owned directory, and every test draws a **byte-for-byte copy** into `<tmp>/demo/kitchen` |
| **the cost** | `internal/ui` **182.2s → 10.5s** plain, and **535.2s → 38.3s** under `-race -shuffle=on -count=2` — **89% → 6%** of go test's 600s limit. A **14×** margin where the slice asked for 3× |
| **the goldens** | **not one moved** |

## What to verify — each by command, not by reading

| # | Check |
|:--:|---|
| 1 | **One template, seeded once.** `TestMain` seeds it into a suite-owned directory, and every test that needs a plan draws a **copy**. The seed does **not** happen per test |
| 2 | The copy is **byte for byte**. `TestTheCopyIsTheTemplate` holds that two copies made at different times have the same content, and that the template holds exactly the **15** files the seeder wrote |
| 3 | The copy lands in a directory named **`demo`**, so the plan is named `demo` and **every golden file is unchanged**. `git diff` over `internal/ui/testdata/` is empty |
| 4 | The template's `project.json` **names nothing**, so a copy's name follows the directory it lands in — and the **stress plan** is a **second** template, seeded on its own |
| 5 | The store survives **`-count=2`**: the seed count does **not** move with the runs and the copies do; the copies do not see each other's writes; `TestMain` removes the template directory and **fails the suite if it survived** |
| 6 | **No test seeds its own plan.** The wire holds it — find it and check it bites |
| 7 | The cost is real: time the suite yourself, plain and under `-race -shuffle=on -count=2`. The target was **well under 200s** under the race gate |
| 8 | `internal/core` did not change. `core.DomainVersion` is **7**. The budget is unchanged: 0 / 0 / 0 / 28 |
| 9 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The judgement calls

**The template names nothing.** Its `project.json` carries **one field** — the name is
empty — so a copy's name follows the directory it lands in. The runner calls this *"the
one judgement call, two lines to flip"*. Is it right?

| | |
|:--:|---|
| **what it buys** | the goldens pin the plan's content, and the project's name comes from the **parent directory's** name — so a copy in `demo/` is named `demo` |
| **what it costs** | the template is not a plan that names itself, and anything that reads the template's own name gets nothing. What breaks, and what would a reader trip over? |

**The render path's wait was not started**, and the runner **named** it rather than
leaving it: `render.go:41` `snapshot()` → `crew.go:54` `readCrew()` →
`herdr.go:76` `exec.Command(herdr)` — two spawns a render. The brief ruled that it
belongs to whoever next touches the **render path**, and the runner put it in the
receipt instead of quietly leaving it. Say whether that placement is right, and
whether the **poll's** two spawns a second (the price of asking a runtime on every
tick, 5ms a call) belongs with it or somewhere else.

## Mutate, run, revert

| | The mutation |
|:--:|---|
| 1 | make a test seed its **own** plan again — the wire must go red |
| 2 | make the copy **not** byte-for-byte — one file changed inside it |
| 3 | make `TestMain` **not** remove the template directory — the suite must fail |
| 4 | make the template name the plan, so a copy lands with the template's name |
| 5 | make the stress plan share the first template |

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
| **The render path** | not this slice's, and not yours. Say where it belongs, and what it must achieve — the brief for it goes to the chef first |

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
cat > work/inbox/reviewer-tui-t11.md <<'EOF'
Reviewer: T11 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-145.. or none>). The store: <where the template lives, and that the seed happens once>. The copy: <that it is byte for byte, and the witness>. The names: <how the copy is named after its directory, and that every golden is unchanged>. The stress plan: <the second template>. Count two: <what survives -count=2, and the cleanup>. The cost: <the plain and race times, measured by you>. The judgement call: <the template naming nothing, and whether that is right>. The render path: <where its wait belongs, and the poll's spawns>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t11.md && git commit -q -m "inbox: reviewer reports tui-t11"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
