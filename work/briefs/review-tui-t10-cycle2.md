You are the **reviewer** for **T10 fix cycle 1**. You are a fresh lane. You did not
write the fix. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[fix-tui-t10-1.md](fix-tui-t10-1.md) — the contract for this cycle — and then
[tui-t10.md](tui-t10.md), the contract the slice was built to. The orchestrator's pane
id is `w1H:p2`. Use it verbatim.

The fix commits: `2684307` (AC-01, one switch), `1c76e84` (AC-02, the poll),
`f59f9be` (AC-03, the progress from the crew's label), `5b03f48` (AC-04, the stamp),
and the row `e91063e`. The receipt is
[work/inbox/runner-tui-t10-fix1.md](../inbox/runner-tui-t10-fix1.md).

## This cycle is two findings

File **HIGH or MEDIUM** blockers. Do not re-review T10 from scratch. New findings
start at **F-142**. A LOW note goes in the message, not in the ledger as a finding.

| Finding | The target | The commit |
|---|---|---|
| **F-140** | the mouse's tab switch and the filter line's `tab` / `shift+tab` never asked the runtime for the crew, so a click on the *Sous-Chef* label left the stamp saying **"the automation is not running"** and the footer offering **"space start the crew"** — while the crew was running. A move offered on a guess is a move that can start a second crew | `2684307`, `1c76e84` |
| **F-141** | the stamp's progress was the plan's history, read by a false rule (*"an event the chef did not write is a write by the crew"*) — the journal carries no signal separating the crew from the chef, and the stamp printed **"round 7 of 3 · 39 of 12 interactions"** with no crew running at all | `f59f9be`, `5b03f48` |

Close both in `docs/audit.md` against those commits, and write the slice row for
this cycle.

## What the fix claims

**One switch.** `switchTab` returns the command, and the keyboard's `tab`, the
mouse's click on a tab label and the filter line's `tab` / `shift+tab` all go
through it.

**The poll** asks the runtime on every poll **while the Sous-Chef tab is visible**,
and asks **even when the plan has not changed**. A tab that is not on screen does
not spend a call.

**The progress** comes from the **orchestrator's label**:
`maitre: orchestrator · round <N> · <M> interactions`. The console recognises
exactly two shapes — `round <N>` and `<M> interactions` — and the rest of the tail
is what the member is on. `writtenByChef` and the journal derivation are **gone**.

**The startup task** says: *"Keep your tab's label current — the console reads the
crew's progress from it: `herdr tab rename <TAB_ID> "maitre: orchestrator · round
<N> · <M> interactions"` (your tab's id is in `herdr pane list`)."*

**`orchestrator-role` says so too** — a new *The roster's labels* section.

**The stamp** with a reported progress: `round 2 of 3 · 6 of 12 interactions`. With
**no** report: **no progress at all** — the vital signs only. A progress **past** the
cap is printed as reported (`round 7 · 39 interactions`) with **"the budget is
spent"** beside it.

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | The keyboard's `tab`, a **click** on the *Sous-Chef* label, and the filter line's `tab` / `shift+tab` all reach the runtime. With a crew running, a **click** leaves the stamp showing the crew and the footer offering **"space stop the crew"** — never a second start |
| 2 | The console asks the runtime **even when the plan has not changed**, while the Sous-Chef tab is visible. A tab that is **not** on screen does not spend a call |
| 3 | The progress comes from the **orchestrator's label**, and `writtenByChef` and the journal derivation are **gone** — the plan's history cannot be printed as the crew's work any more. On the demo plan, with **no crew**, the stamp carries **no progress** |
| 4 | A label carrying no progress shows **none**, and is not an error. A part the console does not recognise carries **no** progress, and a round that is not a number is **not** a round |
| 5 | The **startup task** tells the orchestrator to keep its label current, and **`orchestrator-role`** says so too — with a witness in `internal/cli/skills_test.go` |
| 6 | A progress **past** the cap is printed as reported, with **"the budget is spent"** beside it |
| 7 | `docs/05`'s two stale sentences are gone: the false equivalence and *"The Sous-Chef tab dissolved into Home"* |
| 8 | `internal/core` did not change. `core.DomainVersion` is **7**. The tab count is **8**. The budget is **28 of 30** |
| 9 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The six notes the runner left for you

| | the claim |
|--:|---|
| 1 | **The label's format is the runner's** — the console recognises exactly two shapes (`round <N>`, `<M> interactions`); anything else in the tail is what the member is on |
| 2 | **The skill's new section is witnessed by a test**, `internal/cli/skills_test.go` — not by the CLI's surface. The brief allowed `SKILL.md` to change, and a claim with no witness is the house's bug |
| 3 | **The history was rewritten once** — the first AC-02 commit failed `make lint` (a staticcheck nit in the test file). The ids in the row are the rewritten ones. Check that the rewrite changed nothing else |
| 4 | **The gate's headroom is thin**: `internal/ui` is **516.1s** of go test's 600s default under `-race -shuffle=on -count=2` — **86%** — and this cycle added ~35s. The remedy — **sharing one seeded store** (a copy costs 0.007s against the seeder's 0.608s) — is a repo-wide change and outside this cycle's scope. **One more slice of this size breaks the gate.** Say what that means for the next slice |
| 5 | **The stamp's wording is the runner's** — *"the budget is spent"*. The rule is *no number it cannot compute* |
| 6 | **The over-budget case is new behaviour** — nothing in the slice or the design said what the stamp should say when the crew reports more than the cap. The runner prints the progress as reported and says the budget is spent |

## Mutate, run, revert

Break each rule and report the tests that went red:

| | The mutation |
|--:|---|
| 1 | make the mouse's switch skip the runtime ask again |
| 2 | make the poll ask only when the plan's fingerprint moved |
| 3 | make the console read the progress from the journal again |
| 4 | make a label with **no** progress show a number anyway |
| 5 | make the startup task stop telling the orchestrator to keep its label current |

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
| **The seeded store** | a repo-wide change, and the next slice's work — not this cycle's. Say what it must achieve, and what it must not break |

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
cat > work/inbox/reviewer-tui-t10-cycle2.md <<'EOF'
Reviewer: T10 cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-142.. or none>). One switch: <the keyboard, the mouse and the filter, and what a click does with a crew running>. The poll: <what it asks, and when>. The progress: <where it comes from, what a label with none shows, and the startup task and the skill>. The over-budget case: <what the stamp says, and whether that is right>. The six notes: <each one, and whether it holds>. The gate's headroom: <the number, and what the next slice must achieve>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t10-cycle2.md && git commit -q -m "inbox: reviewer reports tui-t10 cycle 2"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
