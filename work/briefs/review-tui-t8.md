You are the **reviewer** for **T8**. You are a fresh lane. You did not write this code.
Re-derive the facts from the code, not from the runner's report.

The repository is `/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[tui-t8.md](tui-t8.md) — the contract. The orchestrator's pane id is `w1H:p2`. Use it
verbatim.

The runner's commits: `5e2146c` (the skill, embedded), `8ab3b80` (help), `ac3f6d5`
(list), `2e9f3a6` (load), `3333df7` (the refusal), `504cb86` (the description), and
the row `885a46b`. The receipt is
[work/inbox/runner-tui-t8.md](../inbox/runner-tui-t8.md).

## This slice

**T8 — the skills.** `maitre-cli` gains `skill help` · `skill list` · `skill load
<name>`. The skills are **embedded in the binary**, and the first one is
`orchestrator-role` — the role *and* the protocol, for an agent that knows nothing
about this repo.

| | |
|---|---|
| **the skill** | `internal/cli/skills/orchestrator-role/SKILL.md`, embedded with `//go:embed skills` |
| **the list** | computed from the embedded skills — the name and the description come from the skill's own frontmatter. No second list |
| **the refusal** | a skill whose directory and frontmatter disagree is not served |
| **the domain** | `internal/core` and `internal/ui` do not change. `DomainVersion` is **7** |

## What to verify — each by command, not by reading

| # | Check |
|--:|---|
| 1 | `maitre-cli skill help` prints **text** — what the menu is, what the three verbs do, and that the skills ship inside the binary |
| 2 | `maitre-cli skill list` prints a **JSON** document; the name and the description come **from the skill's frontmatter**. Put a second skill into the tree and it appears in the list **with no other edit** |
| 3 | `maitre-cli skill load orchestrator-role` prints the file **byte for byte** — compare the md5 of the output with the md5 of the file in the tree |
| 4 | `maitre-cli skill load` on a skill whose directory and frontmatter disagree **refuses**, and the list does not carry it |
| 5 | The skill carries **the role** (what an orchestrator is — the glossary's definition) **and the protocol** (write the brief → open the lanes → dispatch → wait → receive the receipt → brief the reviewer → close → the ledger) **and the herdr mechanics** including the **two commands that do not work** — `herdr agent prompt` (`agent_not_ready`) and `herdr agent start` (no mcode kind) — **and what it does not do** |
| 6 | The skill names **no recipe** and no pane id. The work arrives in the startup task |
| 7 | The skill is **self-contained**: an agent that knows nothing about this repo can follow it. The repo's places are named — `work/briefs/`, `work/inbox/`, `work/tui/README.md`, `work/questions.md`, `docs/audit.md` |
| 8 | The description routes with **What, When and Near misses** — *What* the skill is for, *When* to load it, and what it is **not** for |
| 9 | `internal/core` and `internal/ui` did not change. `core.DomainVersion` is **7**. The budget is unchanged: 0 / 0 / ≤ 30 |
| 10 | `make check` and `make check-race` are green — confirmed by you, not taken from the receipt |

## The two claims the slice may have left standing

**The runner moved the command count in `AGENTS.md` and `docs/04-operations.md`. Two
claims that the count sits in did not move with it.** Check both, in both files:

| the claim | where |
|---|---|
| *"Every command prints one JSON document"* | `AGENTS.md` and `docs/04-operations.md` — but `skill help` and `skill load` print **text**. `docs/04` was amended; check whether `AGENTS.md` was |
| *"one verb per operation"* | `AGENTS.md` and `docs/04-operations.md` — but the three `skill` verbs are **not operations**. They read embedded files, not the plan. Is the claim now false, or does it need its exception stated, the way the *two commands print text* note was stated? |

**The house's lesson is the doc drift one: the same claim in 2–4 files, nothing
checking the copies.** Say which files carry each claim and what each one now says.

## The two judgement calls the runner raised

**The commit order is not the criterion order.** AC-04 landed first, because
`//go:embed skills` refuses a directory that holds no embeddable file — the skill
has to be in the tree before the embed is. And between AC-01's commit and AC-03's,
**the help named `list` and `load` before they existed.** Every commit is green and
the last one is the whole menu. Is that acceptable — a slice that is green at every
commit, with a menu that briefly advertised verbs it did not yet have? Or is it a
finding? Say which.

**The skill's frontmatter reader** takes `name:` and `description:` on one line or
as a `|` / `>` block, and a skill whose frontmatter has no name is refused. Check
that a skill with no description is also refused or handled sanely — a skill with a
name and no description would appear in the list with nothing to route on.

## Mutate, run, revert

Break each rule and report the tests that went red:

| | The mutation |
|--:|---|
| 1 | make `skill list` take its names from a written list instead of the skills' own frontmatter |
| 2 | serve a skill whose directory and frontmatter disagree |
| 3 | let `skill load` print a JSON document instead of the content |
| 4 | drop the refusal when a skill has no name |

Revert every mutation. `git status --short` at the end is empty except your ledger
commit.

## The questions

Read [work/questions.md](../questions.md). Q11–Q13 are **closed** and Q16–Q28 are
**ruled** — do not reopen them and do not answer them. Say whether anything was
newly raised.

## The scope

| | |
|---|---|
| **Implementation code** | do not edit it, except to mutate and revert |
| **`docs/audit.md`** | yours. One commit for the ledger, and the slice row for this cycle |
| **`work/`** | the orchestrator's, except your receipt |
| **The runner** | do not prompt it. Do not spawn agents |
| **`docs/03`, `docs/05`** | not yours. Report a doc that contradicts the code; do not fix it |

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
cat > work/inbox/reviewer-tui-t8.md <<'EOF'
Reviewer: T8 done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-131.. or none>). help: <what it prints>. list: <the JSON, and that a new skill appears with no other edit>. load: <the content, byte for byte>. The skill: <the role, the protocol, the herdr mechanics, and what it does not do>. Its description: <what it says, and how it routes>. The refusal: <the skill that disagrees with its directory, and that the CLI will not serve it>. The claims: <the two that the count may have left standing, in each file>. The commit order: <your judgement, and why>. Mutations: <what you mutated, what went red, and that you reverted it>. The questions: <answered, missed, or newly raised>. DomainVersion: <7>. The gate: <make check, make check-race>. LOW notes: <any>.
EOF
git add work/inbox/reviewer-tui-t8.md && git commit -q -m "inbox: reviewer reports tui-t8"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

After the nudge, end your turn.
