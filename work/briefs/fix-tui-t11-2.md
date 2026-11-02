You are a **fresh runner** on a **fix cycle**. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/briefs/tui-t11.md](tui-t11.md) — the contract the slice was built to — and then
this brief. This brief wins where they differ. The orchestrator's pane id is `w1H:p2`.

This is **T11, fix cycle 2**. **F-145 is closed** — the rewrite is textual, the schema
copy is gone, and the reviewer's probe now fails in both packages. **One blocker
remains, and it is about the witness itself.**

The finding is in `docs/audit.md` as **F-146** — read the row, it carries the proof.
The reviewer's receipt is [work/inbox/reviewer-tui-t11-cycle2.md](../inbox/reviewer-tui-t11-cycle2.md).

## F-146 — the guarantee is true by construction and held by no witness

**The defect.** `TestTheTemplateCarriesTheDomainFields` compares the template
against a **hand-written list of names**. So the two go stale **together** and
agree — the witness cannot see the *mechanism* regress.

**The proof.** The rewrite parses and re-serialises again **with** the domain's
field → **`internal/ui` GREEN**. The fixture silently drops the field again, and
the witness — which holds a copy of the same list — is happy.

**And the same hole is live today.** The rewrite blanks **every** `"name"` key **at
any depth**, so a **nested** name is blanked too — proved, both packages green —
while the code's own comment claims *"every other byte is left as the domain wrote
it"*.

**The target behaviour.** The fixture's fidelity to the model is checked against
**the model's own output**, never against a list of names written by hand.

**The remedy the reviewer names, and it costs one `Init()` and no seed:**

> `core.Seed` calls `store.Init` (`internal/core/seed.go:33`), so `core.Open(dir).Init()`
> writes the **domain's OWN** project file — parse both, drop `name` and `created`,
> and compare. It holds the field set and the values, nested included, and **needs
> no list of names**. Proved green today and red under both cases the list cannot
> see.

| | |
|---|---|
| **the witness** | compares the template against the **domain's own** project file. No hand-written list of names anywhere in the witness |
| **the rewrite** | blanks **only the top-level name** — not every `"name"` key at any depth. A nested name survives, and the code's comment is **true** |
| **what it catches** | the parsed rewrite that drops the domain's field, **and** a nested name being blanked. Both green today under the old witness; both must be red |

**One fact, one answer.** The domain's output is the truth. A list of names is a
second copy of the schema — the very thing F-145 was about.

## The three low notes, folded in

| | |
|---|---|
| **a seeder call in a comment counts against the wire** | the wire counts mentions, and a comment is a mention. Say how it should count, and make it so |
| **a stat error that is not *not-exist* prints the wrong reason** | the cleanup reports *the* reason. If the error is something else, the message must say that |
| **`TestTheTemplatesAreRemoved` holds `removeTemplates`, not `TestMain`'s own call** | the witness must cover the removal that **TestMain** actually performs |

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the witness compares the template against the **domain's own** project file — written by `core.Open(dir).Init()` — with `name` and `created` dropped, and **no hand-written list of names**. The reviewer's two cases are **red**: the parsed rewrite that drops the domain's field, and a nested name that gets blanked |
| **AC-02** | the rewrite blanks **only the top-level name**. A nested name survives, and the code's comment is **true** |
| **AC-03** | the wire counts **calls**, not mentions — a seeder call in a comment does not count against it — and the sources' count still has to agree with the counter. The cleanup reports the **actual** reason. The removal witness covers **TestMain's own call** |
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
| **The render path** | **not** this cycle's. Its brief is written (`work/briefs/tui-t12.md`) and goes to the chef first |
| **The poll's frequency** | not this cycle's. It is deliberate — F-140's witness requires the poll to ask while the tab is visible |
| **The SIGKILL leftover** | not a finding: the directory belongs to the OS's tmp, each run makes its own, and a stale one is never read. **Do not prune stale directories at TestMain's start** — two lanes would race each other's pruning |
| **`docs/audit.md`** | the reviewer's. The finding is closed by the cycle-3 reviewer, not by you |
| **`work/`** | the orchestrator's, except the T11 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t11-fix2.md <<'EOF'
Runner: T11 fix 2 done. Commits: <n>. The witness: <how it compares the template against the domain's own project file, and the two cases that now go red>. The rewrite: <that it blanks only the top-level name, and that the comment is true>. The wire: <how it counts calls, the agreement with the counter, the cleanup's reason, and TestMain's own removal>. The goldens: <unchanged>. The cost: <plain and race>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t11-fix2.md && git commit -q -m "inbox: runner reports tui-t11 fix 2"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
