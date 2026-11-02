You are a **fresh runner** on a **fix cycle**. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/briefs/tui-t10.md](tui-t10.md) — the contract the slice was built to — and then
this brief. This brief wins where they differ. The orchestrator's pane id is `w1H:p2`.

This is **T10, fix cycle 2**. The review filed **two blockers**, and they are the
same shape as the ones this repo has been teaching for seven cycles: a claim that
nothing holds.

The findings are in `docs/audit.md` as **F-142** and **F-143** — read those rows,
they carry the proof. The reviewer's receipt is
[work/inbox/reviewer-tui-t10-cycle2.md](../inbox/reviewer-tui-t10-cycle2.md).

## F-142 — the crew's label reaches the screen unsanitised

**The defect.** The label is **runtime state**, written by an agent, and it reaches
the screen **unsanitised**: a **newline** in a label puts a line of the agent's text
where the crew's rows are a table, and an **ANSI escape** is in the raw render.
`sanitizePlan` covers the plan and **not** the runtime's label, and
`TestNothingRenderedCarriesControlCharacters` proves the invariant on **five** tabs,
not the Sous-Chef one.

**The target behaviour.** The label is **agent-written text**, and it is sanitised
**where it enters** — in `readCrew`, before it is anywhere near the screen. No
newline, no ANSI escape, no control character from a label can reach the render.
And the witness is on **this** tab: the tab that shows the crew is the tab that
proves the crew's rows carry no control characters.

## F-143 — the label's format is two copies, and nothing checks they agree

**The defect.** The format lives in **two** places: the **skill's example** and the
**console's parser**. Nothing checks that they agree.

**The proof.** The reviewer's probe: the skill's format drifts, **the skill's own
witness is updated**, and **both packages stay green** — `internal/cli` green (3.0s),
`internal/ui` green (183.8s), nothing red in either. The orchestrator is then told
to write a format the console cannot read, and **the stamp falls back to the vital
signs alone with no error anywhere**.

**The target behaviour.** **One format, one answer.** The two copies cannot drift.
The remedy is a **wire**: a test that takes **the skill's own example** and feeds it
to the console's parser, or **one format constant** that the skill's witness
compares against. Which of the two is yours to choose; what is not yours to choose
is to leave the two copies unconnected.

**A drift must turn a test red.** The reviewer's probe — the skill's format drifts
and the skill's witness is updated — must now fail somewhere.

## The low notes, folded in

The reviewer recorded these; three of them are cheap and belong to this cycle.

| | |
|---|---|
| **the stamp's heading runs into the vital signs** | with no progress reported the stamp reads *"THE STAMPplan 7 · build 7 · they agree"*. The heading and the vital signs are **two facts** and they must read as two |
| **the budget-is-spent flag arrives one step late** | the over-budget report and the flag must arrive **together**. *Round 3 of 3* and *12 of 12* say nothing — the flag is **over**, not **at** |
| **a zeroed report reads as no report** | the orchestrator reporting `round 0 · 0 interactions` is a report, not an absence. Say so, or say nothing — but do not read one as the other |
| **the startup task and the skill point at `herdr pane list`** | `herdr pane current` is the precise answer for the pane's **own** tab id |
| **the render path waits on a process, and the poll spends two spawns a second** | a note, not this cycle's work. Say if it belongs to the next one |

## What the orchestrator has already done

The chef's ruling on the progress — **the crew's own label carries it** — now has its
row in `work/tui/README.md`'s rulings table (**Q30**) and in
`plan-automation.md`'s question table. **Do not redo those rows.**

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the crew's label is **sanitised where it enters** — in `readCrew`, before it is near the screen. A **newline** or an **ANSI escape** in a label cannot reach the render, and a witness on the **Sous-Chef** tab proves it |
| **AC-02** | **one format, one answer**: the skill's example and the console's parser cannot drift. A wire joins them — the skill's own example fed to the parser, or one format constant the skill's witness compares against. **A drift turns a test red** |
| **AC-03** | the stamp reads well in all three shapes: a progress reported, **no** progress, and a progress **over** the cap. The heading does not run into the vital signs; the over-budget flag arrives **with** the report; a **zeroed** report reads as a report, not as an absence |
| **AC-04** | the startup task and `orchestrator-role` point at **`herdr pane current`** for the pane's own tab id |
| **AC-05** | `work/tui/README.md`'s T10 row names this fix cycle's commits. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7** |
| **`internal/core`** | does not change |
| **The tab count** | stays **8** |
| **The budget** | unchanged: 28 of 30. herdr is a runtime dependency, not a module |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **The seeder** | **not** this cycle's. Sharing one seeded store is the **next** slice, and it must land before any slice that adds tests. Do not start it here |
| **`docs/audit.md`** | the reviewer's. The findings are closed by the cycle-3 reviewer, not by you |
| **`docs/03`, `docs/04`, `docs/06`** | not this slice |
| **`work/tui/plan-automation.md`**, **`work/questions.md`** | the orchestrator's, and the ruling is already recorded. Leave them |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`work/`** | the orchestrator's, except the T10 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t10-fix2.md <<'EOF'
Runner: T10 fix 2 done. Commits: <n>. The label: <where it is sanitised, and the witness on the tab>. One format: <the wire between the skill's example and the console's parser, and the drift that turns a test red>. The stamp: <the three shapes, and what each says>. The pane's own tab id: <what the startup task and the skill now say>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t10-fix2.md && git commit -q -m "inbox: runner reports tui-t10 fix 2"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
