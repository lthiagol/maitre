You are a **fresh runner** on a **fix cycle**. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/briefs/tui-t10.md](tui-t10.md) — the contract the slice was built to — and then
this brief. This brief wins where they differ. The orchestrator's pane id is `w1H:p2`.

This is **T10, fix cycle 1**. The review filed **two blockers**. One is a move
offered on a guess — the console can start a second crew. The other is a number the
console cannot know.

The findings are in `docs/audit.md` as **F-140** and **F-141** — read those rows,
they carry the proof. The reviewer's receipt is
[work/inbox/reviewer-tui-t10.md](../inbox/reviewer-tui-t10.md).

## F-140 — the tab can say the automation is not running while it is

**The defect.** The **keyboard's** switch asks the runtime when the tab is the
Sous-Chef one. The **mouse's** switch (`internal/ui/mouse.go:57`) and the filter
line's `tab` / `shift+tab` (`internal/ui/filter.go:152,155`) **do not**. So the
console holds an empty crew and a nil error, and a click on the *Sous-Chef* label
leaves the stamp saying **"the automation is not running"** and the footer offering
**"space start the crew"** — while the crew is running. **A move offered on a guess
is a move that can start a second crew.**

**And the poll does not save it.** `checkDisk` reloads only when the plan's
fingerprint moved, so a **static** plan is never re-read. The code asks on every
**reload**; a reload is a change on disk.

**The target behaviour.** **One switch.** The mouse and the filter line go through
the same path the keyboard does, and the tab **never** says the automation is not
running while it is. And while the Sous-Chef tab is **visible**, the console asks
the runtime — **even when the plan has not changed**. The crew is runtime state; a
plan that does not move is exactly when the crew's state is most likely to have
moved.

## F-141 — the stamp's progress is the plan's history, not the crew's work

**The defect.** The rule was *"an event the chef did not write is a write by the
crew"* (`Actor == "chef" || Role == "chef"`), and it is **false**: a default
`maitre-cli` write carries **actor `agent` and no role** — and so does a write by
the orchestrator, the runner or the reviewer, because nothing sets
`MAITRE_CLI_ACTOR` or `MAITRE_CLI_ROLE`. The journal carries **no signal** that
separates the crew from the chef.

**The consequence.** On the plan the demo hands the chef, the stamp reads
**"round 7 of 3 · 39 of 12 interactions"** with **no crew running at all** — and
the order above it says the budget as a **policy** (*3 rounds · 12 interactions*).
*7 of 3* says the policy is not a policy.

**Ruled by the chef: the crew's own label carries the progress.** The orchestrator
writes it, and the console reads it.

**The target behaviour.**

| | |
|---|---|
| **the owner** | the **orchestrator**. It is the one that knows the round and the interaction count — it is the one dispatching the crew |
| **the carrier** | the orchestrator's tab **label** — the same roster's line the design already gave the crew |
| **the console** | reads the progress from the crew's **orchestrator** label. A label that carries none shows **no progress** — never a wrong one |
| **the startup task** | tells the orchestrator to keep its label current with the round and the interaction count |
| **the skill** | `orchestrator-role` says so too. The design says the orchestrator keeps the labels current; the **shipped skill does not say it**, and the skill is the role's own knowledge |
| **the stamp** | prints **no number it cannot compute**. The cap stays the policy; the progress comes from the label or it is not shown |

**The format is yours to choose**, with two requirements: the console can read it
**without guessing**, and a label that carries no progress is **not** an error —
it means the crew has not reported one.

**One fact, one answer.** The progress is the orchestrator's to say and the
console's to read. Nothing derives it from the journal any more — the journal
carries no signal, and a number read from a signal that does not exist is a
number the console invented.

## The two stale sentences in `docs/05`

Both are the reviewer's, and both are **in scope** for this cycle:

| | |
|---|---|
| `docs/05:643` | repeats the false equivalence — *"the **interactions** are the crew's writes — the journal's events the chef did not write"* |
| `docs/05:578` | still says *"The Sous-Chef tab dissolved into Home"* — a T9 sentence that now contradicts the paragraph's own next line and the tab T10 built |

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | **one switch**: the mouse's tab switch and the filter line's `tab` / `shift+tab` go through the same path the keyboard's does. A click on the *Sous-Chef* label never leaves the tab saying *the automation is not running* while the crew is running, and never offers to start a second crew |
| **AC-02** | while the Sous-Chef tab is **visible**, the console asks the runtime **even when the plan has not changed** — a static plan is exactly when the crew's state is most likely to have moved |
| **AC-03** | the **progress comes from the crew's orchestrator label**. The console reads it, a label with no progress shows **none**, and the startup task tells the orchestrator to keep its label current. **`orchestrator-role` says so too** |
| **AC-04** | the stamp prints **no number it cannot compute**, and `docs/05`'s two stale sentences — the false equivalence and *"dissolved into Home"* — are gone |
| **AC-05** | `work/tui/README.md`'s T10 row names this fix cycle's commits. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. The plan's format does not move |
| **`internal/core`** | does not change |
| **`internal/cli`** | does not change — but `internal/cli/skills/orchestrator-role/SKILL.md` **does**, because the role owns the label (AC-03) |
| **`internal/ui`**'s budget | unchanged: 28 of 30. herdr is a runtime dependency, not a module |
| **No new view** | nothing is registered in `core.DerivedViews`, so `docs/03` § 8 does not move |
| **The tab count** | stays **8** |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **`docs/audit.md`** | the reviewer's. The findings are closed by the cycle-2 reviewer, not by you |
| **`docs/03`, `docs/04`, `docs/06`** | not this slice |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`work/tui/plan-automation.md`** | the orchestrator's |
| **`work/`** | the orchestrator's, except the T10 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t10-fix1.md <<'EOF'
Runner: T10 fix 1 done. Commits: <n>. One switch: <what the mouse and the filter now go through, and what a click on the tab does with a crew running>. The poll: <what the console asks while the tab is visible, and how a static plan is handled>. The progress: <the label's format, where the console reads it, what the startup task says, and what the skill says>. The stamp: <what it prints when the crew has reported and when it has not>. The docs: <the two stale sentences, gone>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t10-fix1.md && git commit -q -m "inbox: runner reports tui-t10 fix 1"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
