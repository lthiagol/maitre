You are the **reviewer** for **T9**. You are a fresh lane. You did not write this
code. Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[tui-t9.md](tui-t9.md) — the contract. The orchestrator's pane id is `w1H:p2`. Use it
verbatim.

The runner's commits: `38d1d48` (the schema, the read, the atomic write), `2aba9f4`
(the file, and that the plan's format does not move), `089ef6b` (the verb, and the
count that moves with it), `1ed7f3f` (the console reads and writes), `2d65fa0` (the
refusal), and the row `17cdefe`. The receipt is
[work/inbox/runner-tui-t9.md](../inbox/runner-tui-t9.md).

## This slice

**T9 — the settings record.** The console's preferences stop dying with the
console. A new stdlib-only package `internal/settings` holds the schema and the
file's mechanics; the record is `kitchen/settings.json`; the CLI gains
`maitre-cli settings show`; the console's Settings tab reads the record on load and
writes it on change.

**The file is the tool's, not the plan's.** `internal/core` does not read it,
`core.DomainVersion` stays **7**, `planFormatFacts()` does not move, and the
repo's own plan is not stamped.

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | `internal/settings` holds the mechanics: a **missing file** yields the defaults and is not an error; the write is **atomic** — the bytes go to a temporary file in the same folder and are renamed over the old one; a write into a folder that is not there **fails and leaves the old file whole** |
| 2 | An **unknown version is refused by name** — *"the settings file has version 99, which is not one this build has (version 1)"*. The tool does not guess, and it does not fall back silently |
| 3 | **Unknown keys are ignored** within the version this build knows, so a file with keys it does not have still reads |
| 4 | `internal/core` did not change and **does not read the file**. A test reads the core's sources and fails if one mentions `settings.json` or `internal/settings`. Separately: a plan with the record written into its folder reads back whole, and the plan's readers ignore it |
| 5 | `core.DomainVersion` is **7**. `planFormatFacts()` did not move. The repo's own `kitchen/` is **not** stamped and holds **no** settings file |
| 6 | `maitre-cli settings show` prints the record as one JSON document, and is **read-only** — the console is the only writer. The count moved with it: **46 commands: 40 verbs in 8 groups, plus 6 workspace commands**, held by `TestTheDocsCommandCountMatchesTheCode` in both documents |
| 7 | The console's Settings tab **reads** the record when it starts and **writes** it when the chef changes the theme, the order or show-served — so the chef's choices survive a restart. The journal's event count does not move when the theme changes: the record is not the plan's record |
| 8 | A theme or an order the console does not know is **refused by name** and the console **says so** — *"the theme Neon is not one this build has"*. It never writes a value it does not know, and each value is applied on its own, so an unknown theme does not cost the chef the order |
| 9 | `internal/settings` has a row in **both** `scripts/dep-budget.sh` and `AGENTS.md`'s dependency table — 0 external modules — and `AGENTS.md`'s Layout section names it |
| 10 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The six notes the runner left for you

Each one is a claim about the code. Check them by command.

| | the claim |
|--:|---|
| 1 | **The exception wire grew a second group.** `settings show` never opens the plan's records, so `planlessGroups` finds the `settings` group beside `skill`, and the documents must name **both** — with a phrase per group. The test was renamed to `TestEveryVerbOutsideTheReadOnlyGroupsMapsToAnOperation`. Is the two-group exception stated in both documents, and is the wire's answer still exactly the groups that reach no plan? |
| 2 | **Forward compatibility, precisely.** The version fence wins: a file at version 2 is refused **before** its keys are read, so the unknown-key rule covers keys added *within* version 1. `plan-automation.md` § 15 says T10 moves the schema to version 2, so a T9 build refuses a T10 file by version. Is that what the code does, and is the forward-compatibility claim in the docs true as stated? |
| 3 | **The record's defaults name the console's defaults** — `Kitchen` and `triage` are written in two places, because the CLI prints the record and the CLI does not know the console's vocabulary. A test holds the tie (`TestTheRecordsDefaultsAreTheConsolesDefaults`). Is one answer, or two that happen to agree? |
| 4 | **The console writes into the plan's folder** — running the console in this repo, even `--keys "c"`, writes `kitchen/settings.json`. The orchestrator has since added it to `.gitignore`. Does anything else in the repo assume `kitchen/` holds only the plan's records? |
| 5 | **The console's write is a synchronous disk write in the key path** — `applyTheme`, the order keys and `c`. The plan's fingerprint hashes every `.json` in the folder, so a write is seen by the poll and the console reloads the plan once. Is that acceptable, or does a settings change reload more than it should? |
| 6 | **The state row had a second T9** — the orchestrator's mistake, and it is fixed: the Sous-Chef tab is **T10**. Confirm the table now reads one T9 and one T10 |

## Mutate, run, revert

Break each rule and report the tests that went red:

| | The mutation |
|--:|---|
| 1 | make the read treat a missing file as an error |
| 2 | make the write non-atomic — write straight over the file |
| 3 | make an unknown version fall back to the defaults instead of refusing |
| 4 | make the console write a theme it does not know |
| 5 | make the console **not** read the record on load |
| 6 | add a second read-only group to the dispatch table — does the exception's witness name it? |

Revert every mutation. `git status --short` at the end is empty except your ledger
commit.

## The questions

Read [work/questions.md](../questions.md). Q11–Q13 are **closed** and Q16–Q29 are
**ruled** — do not reopen them and do not answer them. Say whether anything was
newly raised.

## The scope

| | |
|---|---|
| **Implementation code** | do not edit it, except to mutate and revert |
| **`docs/audit.md`** | yours. One commit for the ledger, and the slice row for this cycle |
| **`work/`** | the orchestrator's, except your receipt |
| **The runner** | do not prompt it. Do not spawn agents |
| **`docs/03`, `docs/05`, `docs/06`** | not yours. Report a doc that contradicts the code; do not fix it |

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
cat > work/inbox/reviewer-tui-t9.md <<'EOF'
Reviewer: T9 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-137.. or none>). The package: <internal/settings, what it owns>. The file: <kitchen/settings.json, and that core does not read it>. The mechanics: <the read, the atomic write, the unknown-version refusal, the unknown keys>. The verb: <settings show, and the count that moved with it>. The console: <the Settings tab reading and writing, and the chef's choices surviving a restart>. The refusal: <a theme or an order the console does not know, and that it says so>. The six notes: <each one, and whether it holds>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t9.md && git commit -q -m "inbox: reviewer reports tui-t9"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
