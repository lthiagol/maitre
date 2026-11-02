You are a **fresh runner** on a **fix cycle**. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/briefs/tui-t10.md](tui-t10.md) — the contract the slice was built to — and then
this brief. This brief wins where they differ. The orchestrator's pane id is `w1H:p2`.

This is **T10, fix cycle 3**. The review filed **one blocker**, and it is F-143's
failure **one copy over**.

The finding is in `docs/audit.md` as **F-144** — read the row, it carries the proof.
The reviewer's receipt is
[work/inbox/reviewer-tui-t10-cycle2.md](../inbox/reviewer-tui-t10-cycle2.md) and
[work/inbox/reviewer-tui-t10-cycle3.md](../inbox/reviewer-tui-t10-cycle3.md).

## F-144 — the format lives in three places, and one of them is unwired

**The defect.** The format lives in **three** copies:

| | the copy | wired? |
|---|---|---|
| 1 | the **skill's** table | **yes** — `internal/ui/crew_wire_test.go` reads the skill's own example and feeds it to the parser |
| 2 | the **console's** parser | it is the parser |
| 3 | the **startup task** at `internal/ui/souschef.go:288` — the copy the console **hands the orchestrator at spawn**, in the same breath as the work | **no** |

Nothing parses the startup task's copy, and **its witness is a substring check** —
it holds the *words* and not the *shape*. So the **separator can drift with every
word intact** and nothing goes red. The reviewer's probe:

> the startup task's own copy of the format drifts, separator changed, every
> substring its witness checks still present — and **NOTHING red**, internal/ui ok
> at 196.0s.

Then the orchestrator is handed a `herdr tab rename` it can copy **verbatim**,
writes a label the console **cannot read**, and the stamp falls back to the vital
signs alone **with no error anywhere**. F-143's failure, one copy over.

**The target behaviour.** **All three copies are one answer.** The startup task's
copy is wired to the parser the same way the skill's copy is: the wire test feeds
**the startup task's own example** to the console's own `parseRoster` and
`parseProgress`.

**The remedy is cheap** — the parser is right there, in the same package. The
startup task's example is a string in `souschef.go`; extract it the way the skill's
example is extracted, and feed it to the parser.

**A drift must turn a test red.** The reviewer's probe — the startup task's
separator drifts with every word intact — must now fail.

## And which copy is the specification

The reviewer's low note, and it is worth settling because the wire depends on it:
**the skill's table is the specification; the startup task's copy is an example.**
Nothing says so today. Say so — in the skill, in the startup task, or in the wire's
comment — so that a reader knows which one may be edited and which one must follow.

## The other low note — the flag names a compound budget with one word

**"the budget is spent"** is fired by **two** budgets — the rounds (`round 2 of 3`)
and the interactions (`6 of 12`). So *round 2 of 3 · 12 of 12 interactions · the
budget is spent* says the budget is spent **though the round allowance is not**.

**The target behaviour.** The flag says **which** budget it refers to. The wording
is yours; the requirement is not — a reader must not have to guess which of the two
allowances ran out.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the startup task's copy of the format is **wired to the parser**, the way the skill's copy is: the wire test feeds **the startup task's own example** to the console's own `parseRoster` and `parseProgress`. A drift in the startup task's copy — the separator changing with every word intact — **turns a test red** |
| **AC-02** | the skill's table is named as the **specification** and the startup task's copy as an **example**, and it is said where a reader will find it |
| **AC-03** | the over-budget flag says **which** budget it refers to — the rounds or the interactions. A reader must not have to guess which allowance ran out |
| **AC-04** | `work/tui/README.md`'s T10 row names this fix cycle's commits. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7** |
| **`internal/core`** | does not change |
| **The tab count** | stays **8** |
| **The budget** | unchanged: 28 of 30 |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **The seeder** | **not** this cycle's, and **not** the render path's. The shared seeded store is the **next** slice, and the render path's wait belongs to whoever next touches the **render path**. Do not start either |
| **`docs/audit.md`** | the reviewer's. The finding is closed by the cycle-4 reviewer, not by you |
| **`docs/03`, `docs/04`, `docs/06`** | not this slice |
| **`work/`** | the orchestrator's, except the T10 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t10-fix3.md <<'EOF'
Runner: T10 fix 3 done. Commits: <n>. The third copy: <how the startup task's example is wired to the parser, and the drift that turns a test red>. The specification: <where it says which copy is the specification>. The flag: <the wording, and which budget it names>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t10-fix3.md && git commit -q -m "inbox: runner reports tui-t10 fix 3"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
