You are a **fresh runner** on a **fix cycle**. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/briefs/tui-t11.md](tui-t11.md) — the contract the slice was built to — and then
this brief. This brief wins where they differ. The orchestrator's pane id is `w1H:p2`.

This is **T11, fix cycle 1**. The review filed **one blocker**, and it is the house's
own shape once more: a fixture that can **drift from the model silently**.

The finding is in `docs/audit.md` as **F-145** — read the row, it carries the proof.
The reviewer's receipt is [work/inbox/reviewer-tui-t11.md](../inbox/reviewer-tui-t11.md).

## F-145 — the template's project file is rewritten through a second copy of the plan's schema

**The defect.** The template's `project.json` is rewritten through **a second copy
of the plan's schema** — the fixture's own struct — so the suite's fixture can
**drift from the model silently**. The reviewer's proof, by command: **the domain's
project file gains a field, and `internal/ui` stays green.** The fixture no longer
describes the plan it claims to describe, and nothing says so.

**The target behaviour.** The template's project file carries **the same fields the
domain writes**, always. A field the domain adds is **in** the fixture. The
template may differ from the seeder's output in **one** thing only — **the name** —
and that difference is the point of the store.

**And the remedy is a line: rewrite the name TEXTUALLY.** Not through a struct that
knows today's fields. The name is a value in the file; change that value, and
change nothing else. Then a field the domain adds is in the template because the
template **is** the seeder's output.

**The witness.** The reviewer's probe must now fail: **the domain's project file
gains a field, and a test goes red.** The fixture and the model are one thing, and
a wire holds it.

## The three low notes, folded in

Each is cheap, and each is a thing a later reader would trip over.

| | |
|---|---|
| **the cleanup's failure message prints `<nil>`** | a stat error that says nothing. The message must say **what** happened and **where** the leftover is |
| **the counters are plain ints, and no `ui` test calls `t.Parallel()`** | a test that calls it would **race** the counters. The counters are shared state; make them safe for a test that runs in parallel, so the next one does not have to know |
| **the wire matches the literal `core.Seed(`** | a space, an alias, or a different spelling slips past it. The seed **counter** covers the cost — but the wire is the thing a reader trusts. Say how it should match, and make it match that |

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the template's project file is rewritten **textually** — the name only — so the fixture always carries **the fields the domain writes**. The reviewer's probe now fails: **the domain's project file gains a field, and a test goes red** |
| **AC-02** | the counters are safe under a test that calls `t.Parallel()`, and the cleanup's failure message says **what** happened and **where** the leftover is |
| **AC-03** | the wire matches the seed call the way it should, not by one literal spelling — and the seed counter still covers the cost |
| **AC-04** | `work/tui/README.md`'s T11 row names this fix cycle's commits. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7** |
| **`internal/core`** | does not change |
| **The goldens** | **not one moves.** This slice changes no behaviour |
| **The cost** | the store's numbers hold: about **10s** plain, **under 40s** under the race gate |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **The render path** | **not** this cycle's. Its brief goes to the chef first: *no process spawn in a snapshot*, the goldens independent of herdr on the PATH, the live console still asking asynchronously — and the crew read is **one spawn, not two**. Do not start it |
| **The poll's frequency** | not this cycle's either. It is deliberate — F-140's witness requires the poll to ask while the tab is visible. The poll's frequency stays with the poll |
| **`docs/audit.md`** | the reviewer's. The finding is closed by the cycle-2 reviewer, not by you |
| **`work/`** | the orchestrator's, except the T11 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t11-fix1.md <<'EOF'
Runner: T11 fix 1 done. Commits: <n>. The project file: <how the name is rewritten textually, and the probe that now fails>. The counters: <how they are safe under t.Parallel, and what the cleanup's failure says>. The wire: <how it matches, and the counter that still covers the cost>. The goldens: <unchanged>. The cost: <plain and race, before and after>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t11-fix1.md && git commit -q -m "inbox: runner reports tui-t11 fix 1"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
