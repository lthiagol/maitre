You are a **fresh runner** on a **fix cycle**. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/briefs/tui-t9.md](tui-t9.md) — the contract the slice was built to — and then
this brief. This brief wins where they differ. The orchestrator's pane id is `w1H:p2`.

This is **T9, fix cycle 2**. The review filed **one blocker**, and it is a claim that
nothing holds — the house's own shape, for the sixth time.

The finding is in `docs/audit.md` as **F-139** — read the row, it carries the proof.
The reviewer's receipt is [work/inbox/reviewer-tui-t9-cycle2.md](../inbox/reviewer-tui-t9-cycle2.md).

## F-139 — the slice number is not one number across the plan's own files

**The defect.** The Sous-Chef tab's slice was **T10** in two places and **T9** (and
**T8**) in three. A reader who took *"Q26 — T9 — the tab is the ticket"* at its word
would believe the tab's design had already landed. **T10 is the next slice, and its
writer reads these files.**

**The labels are the orchestrator's, and they are corrected in this commit** — every
ruling about the **tab** now lands on **T10**, the **skills'** rulings stay on **T8**,
the **settings'** stay on **T9**, and `plan-automation.md` § 11 is now *"The three
slices"*. **Do not redo those labels.** What is missing is the **check**.

## The check — the witness

**A check holds the number.** The four files that name the Sous-Chef tab's slice are:

| | |
|---|---|
| `work/tui/README.md` | the **state table**, and its *"what each ruling lands on"* table |
| `docs/05-human-interface.md` | the tab table |
| `work/tui/plan-automation.md` | § 10's Q24 ruling, and § 11's table |

A check reads them and **fails when they disagree**. It must **run in the gate** —
`make check` — the way `make goldens-check` runs, so it cannot be skipped.

**The mechanism is yours to choose**: a target in the `Makefile` beside the other
checks, a script under `scripts/` beside `check-golden-claims.sh`, or a test. What
matters:

| | |
|---|---|
| **it runs in the gate** | `make check` runs it, and it fails the gate when the files disagree |
| **it says what disagrees** | the failure names the files and the two numbers — a check that says "mismatch" without saying where is a nuisance, not a witness |
| **it is a real check** | point it at a file where the number differs and it must **fail**. Prove it |
| **the gate table** | `AGENTS.md`'s gate table gains the target, if you add one — the gate's list is checked against the Makefile |

## The second rule — a test's name is a claim

The reviewer's note, and it is worth the cycle: `TestTheWriteIsAtomic` holds **the
tidiness of a successful write** and **the refusal of a write into a folder that is
not there** — the mutation *"a write creates the folder that is not there"* turns
**only that test** red, while the two new witnesses stay green. It does **not** hold
atomicity: the two new witnesses hold that.

**A test's name is a claim, and it must be one the test holds.** Rename it so the
name says what it holds, or split it into two tests that each hold one thing. The
two new witnesses keep the atomicity claim.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | a check reads the four files and **fails when they disagree** about the Sous-Chef tab's slice. It **runs in the gate** — `make check` — and its failure names the files and the two numbers. Prove it: point it at a file whose number differs and it must fail |
| **AC-02** | no test's name claims a rule it does not hold. `TestTheWriteIsAtomic` is renamed (or split) so it says what it holds — the tidiness of a successful write, and the refusal of a write into a missing folder — while the two new witnesses keep the atomicity claim. A test that fails when the atomicity is broken, and one that fails when the tidiness or the refusal is broken |
| **AC-03** | `work/tui/README.md`'s T9 row names this fix cycle's commits. `AGENTS.md`'s gate table carries the new target, if you added one. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. The plan's format does not move |
| **`internal/core`** | does not change |
| **The record's shape** | does not change. The `crew` section comes with **T10** |
| **`internal/settings`** | stays **stdlib only**. 0 external modules |
| **The console** | does not change. The theme, the order and show-served are read and written exactly as they are |
| **The count** | does not move. Still 46 commands, 40 verbs in 8 groups |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **The slice labels** | the orchestrator's, and **already corrected** in this commit. `work/tui/README.md` and `work/tui/plan-automation.md` — do not edit them |
| **`docs/05`** | you may read it. You may **not** edit it in this cycle |
| **`docs/audit.md`** | the reviewer's. The finding is closed by the cycle-3 reviewer, not by you |
| **`internal/core`** | does not change |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`work/`** | the orchestrator's, except the T9 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t9-fix2.md <<'EOF'
Runner: T9 fix 2 done. Commits: <n>. The check: <what it reads, how it fails, that it runs in the gate, and the proof that it fails when the files disagree>. The name: <what the old test is called now, and what each test holds>. The gate: <the target, and AGENTS.md's table>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t9-fix2.md && git commit -q -m "inbox: runner reports tui-t9 fix 2"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
