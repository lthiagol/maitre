You are a **fresh runner** on a **fix cycle**. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/briefs/tui-t9.md](tui-t9.md) — the contract the slice was built to — and then
this brief. This brief wins where they differ. The orchestrator's pane id is `w1H:p2`.

This is **T9, fix cycle 1**. The slice landed and the review filed **two blockers**.
One is a rule held by nothing. One is a document that contradicts the code.

The findings are in `docs/audit.md` as **F-137** and **F-138** — read those rows,
they carry the proof. The reviewer's receipt is
[work/inbox/reviewer-tui-t9.md](../inbox/reviewer-tui-t9.md).

## F-137 — the atomicity is claimed by a test that cannot fail

**The defect.** `internal/settings`'s `Write` goes to a temporary file in the same
folder and renames it over the old one. `TestTheWriteIsAtomic` proves only that *no
file is left behind* and that *a write into a folder that is not there leaves the old
file whole*.

**The proof.** The reviewer's mutation — the temporary file and the rename replaced by
a straight `os.WriteFile` — leaves **the whole suite green**. `internal/settings`,
`internal/cli` and `internal/ui` all pass. So the rule the brief's AC-01 names, that
`AGENTS.md`'s Layout names, and that the test's own name and comment claim, is **held
by nothing**.

**The target behaviour.** **A half-written file cannot lose the chef's choices.** When
the write fails partway, the file that was there before is **whole** — not truncated,
not empty, not a partial record.

**And it must be held.** A test must go **red** when the write stops being a rename
over the old file. The mechanism is yours to choose — a fault injected between the
temporary write and the rename, with the old file asserted whole; or another witness
that observes the property rather than the mechanism. What matters is the *mutation
result*: replacing the temporary file and the rename with a plain overwrite must turn
**that** test red.

**The house's lesson, the fifth time it has come up:** a green gate that proves
nothing is the v1 lesson about zero-match verification. The implementation is right
today. The claim is unheld.

## F-138 — `docs/05` still says the theme dies with the console

**The defect.** The Settings section reads *"The choice is for the session: nothing is
written to the plan, so a theme is a setting of the console, not of the work"* — while
the **same file's own summary table** says *"it is now persisted in the plan's
`settings.json`"*. The prose is stale in **two** ways: the choice no longer dies with
the console, and the plan's folder **does** gain a file.

**The target behaviour.** The prose says what the code does: the theme, the order and
show-served are persisted in `kitchen/settings.json` — the **tool's** file, in the
plan's folder, **not a plan record**. The plan's format does not move,
`core.DomainVersion` stays **7**, and the plan's readers ignore the file. The console
is the only writer, so a value the console does not know is refused by name and the
console says so.

**`docs/05` was not this slice's scope** — that is why it drifted. The report is the
deliverable, and now the fix is.

## What the orchestrator has already corrected

Three claims were the orchestrator's, and they are fixed in this commit — **do not
redo them**:

| | was | now |
|---|---|---|
| `work/tui/plan-automation.md` § 15 | *"unknown keys are ignored, so an older build reads a newer file"* | the **version fence wins**: a file at a version this build does not know is refused by name, before its keys are read. Unknown keys are ignored only **within** the version this build knows |
| `work/tui/plan-automation.md` § 15 | the CLI verb was `settings show` · `settings set` | **`settings show` only**, and the reason: the console is the only writer |
| `work/questions.md` | Q25 open | **struck** — its subject, Q11–Q13, was closed |

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | the atomic write has a **witness**: a test that goes **red** when the temporary file and the rename are replaced by a plain overwrite. The property it holds is the one the brief names — **a half-written file cannot lose the chef's choices** — and the mutation result is what proves it |
| **AC-02** | `docs/05`'s Settings prose says what the code does: the theme, the order and show-served are persisted in `kitchen/settings.json`, the tool's file in the plan's folder, not a plan record; the plan's format does not move and `core.DomainVersion` stays **7**. The file's own summary table and its prose must agree |
| **AC-03** | `work/tui/README.md`'s T9 row names this fix cycle's commits. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. The plan's format does not move |
| **`internal/core`** | does not change |
| **The record's shape** | does not change. The `crew` section comes with **T10** |
| **`internal/settings`** | stays **stdlib only**. 0 external modules |
| **The console's behaviour** | does not change. The theme, the order and show-served are read and written exactly as they are |
| **The count** | does not move. Still 46 commands, 40 verbs in 8 groups |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **`docs/audit.md`** | the reviewer's. The findings are closed by the cycle-2 reviewer, not by you |
| **`docs/03`, `docs/04`, `docs/06`** | not this slice |
| **`work/tui/plan-automation.md`**, **`work/questions.md`** | the orchestrator's — and already corrected. Leave them |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`work/`** | the orchestrator's, except the T9 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t9-fix1.md <<'EOF'
Runner: T9 fix 1 done. Commits: <n>. The atomicity: <how the write proves it, and the mutation that turns the test red>. The doc: <what docs/05 now says, and that its prose and its table agree>. The orchestrator's corrections: <what you found already fixed>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t9-fix1.md && git commit -q -m "inbox: runner reports tui-t9 fix 1"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
