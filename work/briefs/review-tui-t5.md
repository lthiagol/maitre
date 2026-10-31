You are the **reviewer** for **T5**. You are a fresh lane. You did not write this code. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then [tui-t5.md](tui-t5.md) — the contract. The orchestrator's pane id is `w1H:pJ`. Use it verbatim.

The runner's commits: `e265382`, `93153c2`, `9ae7789`, and the state row `e6719c4`. The receipt is [work/inbox/runner-tui-t5.md](../inbox/runner-tui-t5.md).

## This slice is the Ready tab

File **HIGH or MEDIUM** blockers. New findings start at **F-126**. A LOW note goes in the message, not in the ledger as a finding.

| What | The target |
|---|---|
| **`core.Adoptable`** | true iff `EventStart` would be accepted, asked on a **copy**. A view: the `// view` marker, a registration, a row in `docs/03` § 8. The wire fails both ways |
| **D-4** | `ReadinessOf.Ready` is true iff the recipe is adoptable. An approved recipe with an unmet dependency is out, and the reason names the dependency |
| **The tab** | **Ready**. Seven names: `Home · Recipes · Backlog · Ready · Board · Log · Settings`. `Path` is not a tab label |
| **The blocks** | ADOPTABLE from `core.Adoptable`. SMALL TASKS from the open `feature` and `bug` items in `core.BacklogGroups`. An idea is not on the tab |
| **The row** | a small task's row says **promote → sign off → run** |
| **The moves** | `space` on an adoptable recipe writes nothing. The footer does not offer *claim*. `space` on a small task promotes. The cursor does not stay on the id that left, and it does not leave the tab |
| **The Board** | still answers `core.Path`. `maitre-cli path` still prints that answer. A recipe the Board shows and Ready does not is the proof |

**`DomainVersion` stays 7.** The view `path` and the CLI command `path` keep their names.

## The orchestrator's ruling — a blocked recipe

The brief's table said a blocked recipe is out. The brief also said **ask the claim, do not restate it**. The runner found `EventStart` has no blocked gate, so a blocked approved recipe is adoptable.

| | |
|---|---|
| **This cycle** | Adoptable matches the claim. A blocked approved recipe is **in** |
| **Not a finding** | a gate that refuses a blocked recipe would move the shape, and this cycle does not move `DomainVersion` |
| **What you check** | the view really asks `EventStart` on a copy. A copy of the gates that forgot the blocked case, or a special case that excludes it, is the drift the brief forbade |

The note-cursor test now drives the Board, because Ready no longer lists an executed recipe. Confirm the test still holds the cursor, and was not emptied to stay green.

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | An approved recipe with its dependencies closed is adoptable. An unmet dependency, a cancellation, an executing recipe, and a draft are not. A sign-off whose digest does not match is not |
| 2 | A blocked approved recipe **is** adoptable, and that is `EventStart`'s answer, not a special case |
| 3 | `ReadinessOf.Ready` agrees, and the unmet dependency's reason names it |
| 4 | The tab bar is the seven names, with Ready where Path was. The two blocks are the two views. The small task's row carries the two further steps |
| 5 | `space` on an adoptable recipe leaves the journal byte-for-byte the same. `space` on a small task promotes. The footer does not offer a claim |
| 6 | RCP04, or whichever approved recipe waits on an open one, is on the Board and on `core.Path`, and not on Ready. `maitre-cli path` still says what it waits on |
| 7 | `make check` and `make check-race` are green. `make goldens` leaves `git diff` empty. Say what moved: the tab label, and the readiness mark on the recipe that waits |

**Mutate, run, revert.** One mutation that makes an unmet dependency adoptable, and one that keeps `ReadinessOf` blind to dependencies. Report the tests that went red. Revert both. `git status --short` at the end is empty except your ledger commit.

## The questions

Read [work/questions.md](../questions.md).

| # | What to say |
|--:|---|
| Q3 | answered, and this slice is where it landed |
| Q8, Q9, Q10 | answered, not this slice |
| Q11, Q12, Q13 | still open. Do not answer them |
| The blocked recipe | not a new question. The ruling above closed it for this cycle |

## The scope

| | |
|---|---|
| **Implementation code** | do not edit it, except to mutate and revert |
| **`docs/audit.md`** | yours. One commit for the ledger. Write the slice row |
| **`work/`** | the orchestrator's, except your receipt |
| **The runner** | do not prompt it. Do not spawn agents |

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
cat > work/inbox/reviewer-tui-t5.md <<EOF
Reviewer: T5 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-126.. or none>). Adoptable: <in, out, and that it asks EventStart on a copy>. The blocked recipe: <in, because the claim says so>. Readiness: <that it agrees, and the reason names the dependency>. The tab: <the seven names, the two blocks, the row's words>. The moves: <space claims nothing, space promotes, the cursor stays on the tab>. The Board: <a recipe it shows that Ready does not>. The note cursor: <the test still holds it>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t5.md && git commit -q -m "inbox: reviewer reports tui-t5"
herdr agent prompt "w1H:pJ" "Reviewer: T5 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-126.. or none>). Adoptable: <in, out, and that it asks EventStart on a copy>. The blocked recipe: <in, because the claim says so>. Readiness: <that it agrees, and the reason names the dependency>. The tab: <the seven names, the two blocks, the row's words>. The moves: <space claims nothing, space promotes, the cursor stays on the tab>. The Board: <a recipe it shows that Ready does not>. The note cursor: <the test still holds it>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>."
```

After the nudge, end your turn.
