You are the **reviewer** for **T9 fix cycle 1**. You are a fresh lane. You did not
write the fix. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[fix-tui-t9-1.md](fix-tui-t9-1.md) — the contract for this cycle — and then
[tui-t9.md](tui-t9.md), the contract the slice was built to. The orchestrator's pane
id is `w1H:p2`. Use it verbatim.

The fix commits: `b40eedf` (the witness for the atomic write), `97f1de9` (the doc),
and the row `8a6bf57`. The receipt is
[work/inbox/runner-tui-t9-fix1.md](../inbox/runner-tui-t9-fix1.md).

## This cycle is two findings

File **HIGH or MEDIUM** blockers. Do not re-review T9 from scratch. New findings
start at **F-139**. A LOW note goes in the message, not in the ledger as a finding.

| Finding | The target | The commit |
|---|---|---|
| **F-137** | the atomicity was claimed by a test that could not fail: the reviewer's mutation — the temporary file and the rename replaced by a straight `os.WriteFile` — left **the whole suite green** | `b40eedf` |
| **F-138** | `docs/05` still said the theme's choice was *"for the session: nothing is written to the plan"*, while its own summary table said the opposite | `97f1de9` |

Close both in `docs/audit.md` against those commits, and write the slice row for
this cycle.

## What the fix claims

**Two seams in `internal/settings`** — `writeTemp` (the bytes into the temporary file)
and `renameTemp` (the rename over the old one) — so a fault can be injected at
**either** stage, and two witnesses hold the property:

| Witness | The fault | What it holds |
|---|---|---|
| `TestAHalfWrittenFileCannotLoseTheChefsChoices` | the bytes go **half way** into the temporary file and the write fails | the old file is byte for byte whole, and the half-written temporary file is gone |
| `TestARenameThatFailsCannotLoseTheChefsChoices` | the bytes are staged and the **rename** fails | the old file is byte for byte whole, and the staged bytes went with the fault |

**They are property witnesses, not mechanism watchers.** The property is *"the chef's
file survives a write that fails"*. A write that replaces the old file in place has
neither stage, so the fault cannot be injected at all.

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | Replace the temporary file and the rename with a straight `os.WriteFile` — the reviewer's mutation. **Both** witnesses go red. `internal/cli` and `internal/ui` stay green |
| 2 | Keep the temporary file, replace the **rename** with writing over the target — the **rename** witness goes red |
| 3 | Keep the rename, replace the **bytes** with a write that fails — the **bytes** witness goes red |
| 4 | The seams are **restored**: each witness restores them with `defer`, so `-shuffle=on` and `-count=2` are safe, and the race gate is green. Prove it: run the settings package with `-count=2` and with `-shuffle=on` |
| 5 | `docs/05`'s Settings **prose** says what the code does, **and** its own summary table agrees with it — the Theme row, the Settings row, and the prose all say the same thing about who owns the file |
| 6 | The tab table's stale label is gone: the Sous-Chef tab is **T10**, not *"(T9 — planned)"*. Check every mention of the Sous-Chef tab across `docs/05`, `work/tui/README.md`, `work/briefs/tui-t9.md` and `plan-automation.md` — they must all name the same slice |
| 7 | `internal/core` did not change. `core.DomainVersion` is **7**. The count did not move: still **46 commands, 40 verbs in 8 groups, plus 6 workspace commands**. `internal/settings` is still stdlib only |
| 8 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The judgement calls the fix made

**Two package-level function variables in production code, as test seams.** They are
the slice's only new production surface. Is a seam in the shipped package acceptable
for this, and is the seam restored cleanly enough that a test which forgets to restore
it cannot leak into another?

**The old `TestTheWriteIsAtomic` stayed.** It proves the weaker things — no file left
behind, and a write into a folder that is not there leaves the old file whole. Two
witnesses now sit beside it. Is three tests for one write the right shape, or is the
old one now redundant?

**The commits were rebuilt.** The first AC-01 attempt had one seam, and the second
mutation left it green — so the three commits were re-made over the same file sets.
Say whether the history tells the reader anything, or whether the row naming the
final ids is enough.

## Mutate, run, revert

Break each rule and report the tests that went red:

| | The mutation |
|--:|---|
| 1 | make the write replace the old file in place — no temporary file, no rename |
| 2 | make the rename write over the target instead of renaming |
| 3 | make a witness forget to restore a seam — does the next test see it? |

Revert every mutation. `git status --short` at the end is empty except your ledger
commit.

## The questions

Read [work/questions.md](../questions.md). Q11–Q13 are **closed**, Q25 is **closed**,
and Q16–Q29 are **ruled** — do not reopen them and do not answer them. Say whether
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
cat > work/inbox/reviewer-tui-t9-cycle2.md <<'EOF'
Reviewer: T9 cycle 2 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-139.. or none>). F-137: <the two witnesses, the mutation results, and that the seams are restored>. F-138: <the prose and the table agreeing, and the stale label gone>. The seams: <whether a test seam in the shipped package is acceptable>. The old test: <whether it is redundant now>. The history: <whether the rebuilt commits tell the reader anything>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t9-cycle2.md && git commit -q -m "inbox: reviewer reports tui-t9 cycle 2"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
