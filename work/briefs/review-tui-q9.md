You are the **reviewer** for **Q9**. You are a fresh lane. You did not write this code. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then [tui-q9.md](tui-q9.md) — the contract. The orchestrator's pane id is `w1H:pJ`. Use it verbatim.

The runner's commits: `6714558`, `cfb94aa`, and the state row `1d761da`. The receipt is [work/inbox/runner-tui-q9.md](../inbox/runner-tui-q9.md).

## This slice is two rulings

File **HIGH or MEDIUM** blockers. Do not re-review T4. New findings start at **F-126**. A LOW note goes in the message, not in the ledger as a finding.

| Ruling | The target | The commit |
|---|---|---|
| **Q9** | a save writes `"scope": []` and `"tags": []` when the lists are empty. A nil list and an empty list are the same record on disk. A file that lacks the keys still loads, and the next save writes them. The evidence is the **raw file**, not the loaded struct | `6714558` |
| **Q10** | the cursor remembers the item's **id**. After a shape, that id, on its new row. After a create, the new id. After a dismiss or a promote, not the id that left, still on the Backlog tab, and the recipe cursor does not follow the promotion | `cfb94aa` |

**`DomainVersion` stays 7.** The fingerprint hashes field names. If the shape moved, that is a blocker.

## The runner's two notes — the orchestrator's ruling

| Note | Ruling |
|---|---|
| Footer and detail checks that read a settled item sit back on that item, because the cursor leaves it | **Check the diff.** A test that was moved to re-select the item is fine. A test that no longer fails when the cursor stays on the settled item is a hole — file it |
| When the id that left was the only row, the cursor stays, because there is no other row | **Not a finding.** An empty list has no row to follow. Confirm the cursor does not still name the id that left |

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | Create an item with no scope and no tags. The **raw file** contains `"scope": []` and `"tags": []` |
| 2 | Save a nil list. The raw file still has `[]`, not `null` and not a missing key |
| 3 | A file without the keys loads, the lists are empty, and the next save writes them |
| 4 | Shape an idea into a small task. The cursor's id is the shaped id, not the item that slid into the old index |
| 5 | Create: the cursor's id is the new id. Dismiss and promote: the cursor's id is not the one that left, the tab is still Backlog, and the recipe cursor did not move to the new recipe |
| 6 | `core.DomainVersion` is 7. `TestDomainShapeIsVersioned` is green. No new CLI verb. `kitchen/` was not rewritten |
| 7 | `make check` and `make check-race` are green. `make goldens` leaves `git diff` empty. The goldens were reported unchanged — confirm that |

**Mutate, run, revert.** One mutation for the writer (an empty list omitted again) and one for the cursor (the index kept across a shape). Report the test that went red. Revert both. `git status --short` at the end is empty except your ledger commit.

## The questions

Read [work/questions.md](../questions.md).

| # | What to say |
|--:|---|
| Q8 | answered, not this slice. The free edit stays unbuilt |
| Q9, Q10 | answered, and this slice is where they landed |
| Q11, Q12, Q13 | still open. Do not answer them |

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
cat > work/inbox/reviewer-tui-q9.md <<EOF
Reviewer: Q9 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-126.. or none>). Q9: <the raw file, the nil, the missing keys>. Q10: <shape, create, dismiss, promote — the id under the cursor>. The settled-item tests: <still bite / the hole>. The only row: <the cursor does not name the id that left>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. The goldens: <unchanged / what moved>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-q9.md && git commit -q -m "inbox: reviewer reports tui-q9"
herdr agent prompt "w1H:pJ" "Reviewer: Q9 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-126.. or none>). Q9: <the raw file, the nil, the missing keys>. Q10: <shape, create, dismiss, promote — the id under the cursor>. The settled-item tests: <still bite / the hole>. The only row: <the cursor does not name the id that left>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. The goldens: <unchanged / what moved>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>."
```

After the nudge, end your turn.
