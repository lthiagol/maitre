You are the **reviewer** for **T9 fix cycle 2**. You are a fresh lane. You did not
write the fix. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[fix-tui-t9-2.md](fix-tui-t9-2.md) — the contract for this cycle — and then
[tui-t9.md](tui-t9.md), the contract the slice was built to. The orchestrator's pane
id is `w1H:p2`. Use it verbatim.

The fix commits: `2fb29f2` (the check), `b170d0c` (the names), and the row and the
gate table `107a961`. The receipt is
[work/inbox/runner-tui-t9-fix2.md](../inbox/runner-tui-t9-fix2.md).

## This cycle is one finding

File **HIGH or MEDIUM** blockers. Do not re-review T9 from scratch. New findings
start at **F-140**. A LOW note goes in the message, not in the ledger as a finding.

| Finding | The target | The commit |
|---|---|---|
| **F-139** | the Sous-Chef tab's slice number was not one number across the plan's own files — **T10** in two places, **T9** and **T8** in three. A reader taking *"Q26 — T9 — the tab is the ticket"* at its word would believe the tab's design had already landed | `2fb29f2`, `b170d0c`, `107a961` |

**The labels are the orchestrator's and were corrected in `dc3aa21`**, before this
cycle. The fix's job was the **check** that holds the number. Close F-139 in
`docs/audit.md` against those commits, and write the slice row for this cycle.

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | `scripts/check-slice-numbers.sh` reads the **thirteen places** in the plan's three files, and `make slices-check` runs it |
| 2 | The check **runs in the gate**: `make check` and `make check-race` both run it. `AGENTS.md`'s gate table carries the target |
| 3 | On a disagreement the check **fails and names the places and the numbers** — the files and the two numbers, so the reader knows where to look |
| 4 | The check **refuses to pass having found nothing**: a place whose number cannot be read fails it — *"the slice check found no number in … (the tab table). The file changed shape — the check must be taught the new shape, or it proves nothing."* |
| 5 | The names: `TestTheWriteIsAtomic` is gone, and the four tests each hold **one** claim — the tidiness, the refusal, the atomicity in the bytes, the atomicity in the rename. Each has a mutation that turns **only** it red |
| 6 | `internal/core` did not change. `core.DomainVersion` is **7**. The count did not move: still 46 commands, 40 verbs in 8 groups. `internal/settings` is still stdlib only |
| 7 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The two judgement calls the fix raised

**The check's written claim is a hand-written list** — `TAB_QUESTIONS="Q16 Q17 Q18
Q19 Q20 Q21 Q22 Q26 Q27"` — of the questions the plan asks *about the tab*. The
runner's reason: the rules about the tab's state (Q19), the journal (Q20), the
startup task (Q21) and the stop (Q22) never say *tab*, so no reading of the text
can pick them out.

**So the check is not exhaustive.** A new ruling about the tab that is not in the
list would not be checked — and the list is the one thing a reader has to read to
know what the check holds. Is a hand-written list acceptable for a check whose
subject is "the tab's slice is one number"? Or can the check be made exhaustive —
for instance by requiring that every ruling the *lands-on* table sends to the tab's
slice **is** in the list, so a new one fails the check instead of escaping it? Say
which, and why.

**`AGENTS.md`'s gate table is not checked against the Makefile.** The fix brief
claimed it was — *that claim was the orchestrator's, and it was false.* The runner
looked, found no wire, and reported it rather than complying silently. So: the new
target is in `AGENTS.md`'s table and in the Makefile, **by hand**. Is that a finding
for this cycle, or the next round's work? And is there a shape the repo already has
for it — `make versions-check`, the golden checks — that this should have used?

## Mutate, run, revert

Break each rule and report what went red:

| | The mutation |
|--:|---|
| 1 | change the number in **each** of the five places the runner listed, one at a time — the check must fail on every one |
| 2 | remove the number from one place — the check must **not** pass having found nothing |
| 3 | add a **new ruling about the tab** to the *"what each ruling lands on"* table, sending it to the tab's slice but **not** adding it to the list — what does the check do? |
| 4 | make one of the four settings tests hold a claim its name does not — the tidiness, the refusal, or the atomicity |

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
cat > work/inbox/reviewer-tui-t9-cycle3.md <<'EOF'
Reviewer: T9 cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-140.. or none>). The check: <what it reads, how it fails, that it runs in the gate>. Its exhaustiveness: <the hand-written list, and whether a new tab ruling escapes it>. The gate table: <whether the table and the Makefile should be wired, and in what shape>. The names: <the four tests, and the mutation that turns only each red>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t9-cycle3.md && git commit -q -m "inbox: reviewer reports tui-t9 cycle 3"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
