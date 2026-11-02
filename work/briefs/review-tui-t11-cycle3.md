You are the **reviewer** for **T11 fix cycle 2**. You are a fresh lane. You did not
write the fix. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[fix-tui-t11-2.md](fix-tui-t11-2.md) — the contract for this cycle — and then
[tui-t11.md](tui-t11.md), the contract the slice was built to. The orchestrator's pane
id is `w1H:p2`. Use it verbatim.

The fix commits: `d0ab68d` (the witness asks the domain), `55368d4` (the rewrite blanks
only the top-level name), `ecebfe4` (the wire counts calls), and the row `d22cc8e`. The
receipt is [work/inbox/runner-tui-t11-fix2.md](../inbox/runner-tui-t11-fix2.md).

## This cycle is one finding

File **HIGH or MEDIUM** blockers. Do not re-review T11 from scratch. New findings
start at **F-147**. A LOW note goes in the message, not in the ledger as a finding.

| Finding | The target | The commit |
|---|---|---|
| **F-146** | **F-145's guarantee was true by construction and held by no witness.** `TestTheTemplateCarriesTheDomainFields` compared the template against a **hand-written list of names**, so the two went stale together and agreed — the rewrite could parse and re-serialise again, dropping the domain's field, and `internal/ui` stayed **green**. And the same hole was live: the rewrite blanked **every** `"name"` key **at any depth**, so a nested name was blanked too, while the code's comment claimed *"every other byte is left as the domain wrote it"* | `d0ab68d`, `55368d4`, `ecebfe4` |

Close it in `docs/audit.md` against those commits, and write the slice row for this
cycle — and, if T11's findings are all closed, **the row that closes the slice**.

## What the fix claims

**The witness asks the DOMAIN, not a list of names.** `core.Seed` calls
`store.Init`, so `core.Open(dir).Init()` writes the domain's **own** project file
**without a seed**. The witness parses the template file and that one, drops `name`
and `created` from both, and compares — the **field set and the values**, nested
included. **The hand-written list is gone.**

**Both of your cases are red**, proved by command and reverted: the rewrite parsing
and re-serialising again *and* the domain file gaining a field; and the domain file
carrying a **nested** name that gets blanked. Both say *"the template project file
is not the domain's own"*.

**The rewrite blanks ONLY the top-level name** — the value is located in the
domain's own JSON through the **decoder's offsets** and replaced by an empty string;
**every other byte is the domain's**. A nested name survives, and the comment is
**true**: the template file is byte-identical to the seeder's apart from the name
value.

**The wire counts CALLS, not mentions** — the sources are parsed with `go/parser` and
`go/ast`, so a comment that says `core.Seed(` is not a call, an alias and a newline
are the same call, and `seededPlan(` is not a call. The count still has to agree with
the counter.

**The cleanup reports the ACTUAL reason** — *is still there* · *could not check* ·
*could not remove*. And the removal witness **drives the function `TestMain` calls**,
while the call itself is held by reading the source: take the call out of
`TestMain` and the witness goes red.

## What to verify — each by command, not by reading

| # | Check |
|:--:|---|
| 1 | The witness asks the **domain**: `core.Open(dir).Init()` writes the domain's own project file, the two are parsed, `name` and `created` are dropped from both, and the **field set and the values** are compared — **nested included**. **No list of names anywhere** in the witness |
| 2 | **Your two cases are red**: the rewrite parsing and re-serialising again *and* the domain file gaining a field; and a **nested** name that gets blanked. Reproduce both by hand |
| 3 | The rewrite blanks **only the top-level name**, located through the decoder's offsets; **every other byte is the domain's**; the template file is byte-identical to the seeder's apart from the name value; a nested name survives; and the code's comment is **true** |
| 4 | The wire counts **calls**, not mentions: a comment saying `core.Seed(` is not a call, an alias and a newline are one call, `seededPlan(` is not a call — and the count still agrees with the counter |
| 5 | The cleanup reports the **actual** reason; the removal witness drives the function `TestMain` calls, and the **call itself** is held — take it out of `TestMain` and the witness goes red |
| 6 | The **goldens** are unchanged. `git diff` over `internal/ui/testdata/` is empty |
| 7 | The cost holds: about **10s** plain, **under 40s** under `-race -shuffle=on -count=2` |
| 8 | `internal/core` did not change. `core.DomainVersion` is **7**. The budget is unchanged: 0 / 0 / 0 / 28 |
| 9 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The note the runner left for you

**A domain that ADDS a field is red in `internal/core` and green in `internal/ui`** —
and the runner argues that is *right there*, because the fixture keeps the field.
Confirm it: the fixture's fidelity to the model is now held by `internal/ui`, and
`internal/core`'s own shape wire is what catches the domain's move.

## Mutate, run, revert

| | The mutation |
|:--:|---|
| 1 | make the rewrite parse and re-serialise again |
| 2 | blank a **nested** name |
| 3 | make the wire count mentions again — a seeder call in a comment |
| 4 | take the removal call out of `TestMain` |

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
| **The render path** | not this slice's. Its brief (`work/briefs/tui-t12.md`) goes to the chef first: *no process spawn in a snapshot*, the crew read is **one spawn**, the goldens render with **herdr off the PATH**, the live console still asks asynchronously — and it must not break what the Sous-Chef tab shows, a single golden, F-140's witness, or the poll's frequency |

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
cat > work/inbox/reviewer-tui-t11-cycle3.md <<'EOF'
Reviewer: T11 cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-147.. or none>). The witness: <that it asks the domain, and the two cases that go red>. The rewrite: <that it blanks only the top-level name, and that the comment is true>. The wire: <that it counts calls, and that the counts agree>. The cleanup and the removal: <the actual reason, and TestMain's own call>. The note: <a domain that adds a field — red in internal/core, green in internal/ui — and whether that is right>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t11-cycle3.md && git commit -q -m "inbox: reviewer reports tui-t11 cycle 3"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
