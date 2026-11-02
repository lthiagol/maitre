You are the **reviewer** for **T11 fix cycle 1**. You are a fresh lane. You did not
write the fix. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[fix-tui-t11-1.md](fix-tui-t11-1.md) — the contract for this cycle — and then
[tui-t11.md](tui-t11.md), the contract the slice was built to. The orchestrator's pane
id is `w1H:p2`. Use it verbatim.

The fix commits: `2accb50` (the textual rewrite), `3c2be62` (the counters and the
cleanup), `6ccfbdc` (the wire), and the row `5a82933`. The receipt is
[work/inbox/runner-tui-t11-fix1.md](../inbox/runner-tui-t11-fix1.md).

## This cycle is one finding

File **HIGH or MEDIUM** blockers. Do not re-review T11 from scratch. New findings
start at **F-146**. A LOW note goes in the message, not in the ledger as a finding.

| Finding | The target | The commit |
|---|---|---|
| **F-145** | the template's project file was rewritten through **a second copy of the plan's schema**, so the fixture could **drift from the model silently** — the domain's project file gained a field and `internal/ui` stayed green | `2accb50` |

Close it in `docs/audit.md` against that commit, and write the slice row for this
cycle — and, if T11's findings are all closed, **the row that closes the slice**.

## What the fix claims

**The name is rewritten TEXTUALLY.** `unnameTheTemplate` reads the seeder's **own
bytes** and empties the name value **where it stands**; **nothing is parsed and
nothing is re-serialised**, and the copy of the plan's schema is **gone**. So a
field the domain adds is in the template — **the template IS the seeder's output**.

**The probe now fails.** The domain's project file gains a field → **a test goes
red** — *"got [charter created name owner settings], want [charter created name
settings]"* — and `internal/core` goes red too. Both reverted.

**The counters are `atomic.Int64`** — shared state a test calling `t.Parallel()`
cannot race. The witness runs two subtests **in parallel** and draws plans at once;
probed, plain counters say **DATA RACE** and these say **nothing**.

**The cleanup's failure says what happened and where the leftover is** — and no
longer prints the nil stat error.

**The wire matches the CALL, not one spelling** — `\bSeed(Stress)?\s*\(` — so a
qualifier, an alias, a space and a newline are all caught. And it asks **how many**:
the count of seeder calls in the sources must agree with the counter, so a seeder
call smuggled into the template's own file is caught too — *the sources call the
seeder 3 times, but the suite seeded 2*.

## What to verify — each by command, not by reading

| # | Check |
|:--:|---|
| 1 | The name is rewritten **textually**: the seeder's **own bytes**, the name value blanked **where it stands**, nothing parsed and nothing re-serialised. A field the domain adds is in the template |
| 2 | **Your probe now fails**: the domain's project file gains a field, and a test goes red — in `internal/ui` **and** in `internal/core`. Reproduce it by hand |
| 3 | The counters are safe under a test that calls `t.Parallel()`. The witness runs parallel subtests and the race detector is **quiet** — and plain counters would say **DATA RACE** |
| 4 | The cleanup's failure message says **what** happened and **where** the leftover is, and does not print the nil stat error |
| 5 | The wire matches the **call**, not one spelling — a qualifier, an alias, a space and a newline are all caught — and the count of seeder calls in the sources agrees with the counter |
| 6 | The **goldens** are unchanged. `git diff` over `internal/ui/testdata/` is empty |
| 7 | The cost holds: about **10s** plain, **under 40s** under `-race -shuffle=on -count=2` |
| 8 | `internal/core` did not change. `core.DomainVersion` is **7**. The budget is unchanged: 0 / 0 / 0 / 28 |
| 9 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The two notes the runner left for you

| | the claim |
|:--:|---|
| 1 | **The wire is a list of field names in the test, not a copy of the schema in the mechanism.** It changes nothing and fails loudly. Is a wire that holds a list of names the right shape, or does it re-create the drift it was meant to prevent? |
| 2 | **The SIGKILL leftover is not fixed** — a killed test binary leaves the template directory in the user's tmp. Say whether that matters, and to whom it belongs |

## Mutate, run, revert

| | The mutation |
|:--:|---|
| 1 | make the template's project file **parsed and re-serialised** again — the domain's field must go missing |
| 2 | make the counters plain ints again |
| 3 | make the wire match only the literal `core.Seed(` |
| 4 | make the cleanup's message print the nil stat error again |

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
| **The render path** | not this slice's. Its brief is written and goes to the chef first: *no process spawn in a snapshot*, the goldens independent of herdr on the PATH, the live console still asking asynchronously, and the crew read is **one spawn, not two**. Say what it must achieve, and what it must not break |

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
cat > work/inbox/reviewer-tui-t11-cycle2.md <<'EOF'
Reviewer: T11 cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-146.. or none>). The project file: <how the name is rewritten textually, and that the probe now fails>. The counters: <the atomic counters, the parallel witness, and the cleanup's message>. The wire: <that it matches the call and the counts agree>. The two notes: <each one, and whether it holds>. The render path: <what its brief must achieve and must not break>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t11-cycle2.md && git commit -q -m "inbox: reviewer reports tui-t11 cycle 2"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
