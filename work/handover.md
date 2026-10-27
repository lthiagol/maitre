# Handover — you are the new orchestrator

**Paste this whole file as your first prompt.** It is written to be self-contained:
you can run the whole rename from it and the repo.

You are the **orchestrator** of the maitre project's vocabulary rename. You have two
lanes — a **runner** and a **reviewer** — and they report to **you**. This is a
handover: the previous orchestrator's session ended, and the work continues from what
is on disk.

## Read first, in this order

1. `work/README.md` — the state, the rulings, the rules, the report protocol.
2. `work/handover.md` — this file.
3. `docs/07-vocabulary.md` — the decision of record. What the rename *is*.
4. `AGENTS.md` — the house style, the working rules, the dependency budget.
5. `docs/audit.md` — the ledger: what is closed, what is open, and why.

## Where the work stands

| Slice | What | State |
|--:|---|---|
| **A** | the project (`maitre`), the binaries (`maitre`, `maitre-cli`), the ids (`RCP01`), the version | ✅ **closed** — the review said `ok`, zero blockers |
| **B** | the roles (`runner`, `reviewer`, `orchestrator`), the tabs (`Board`, `Sous-Chef`) | 🟡 **in flight** — the previous runner's turn ended without reporting; the work is **green and uncommitted** in the tree, captured at `work/state/slice-b-wip.patch`. **Do not revert it.** |
| **C** | the backlog — one item, `kind: idea \| feature \| bug` | queued — the version goes to 5 |
| **D** | `Proof` → evidence, the recipe-level stamp → completion | queued — the version goes to 6 |
| **E** | the mascot, the glossary rewrite, ADR 0009, the sweep, two LOW notes | queued |

The briefs are written: `work/briefs/slice-{b,c,d,e}.md`, and
`work/briefs/slice-b-resume.md` is the continuation brief for the work in flight —
**start there**. `work/briefs/review-a2.md` is the pattern for a review brief.

## Your first four commands

```bash
cd /Users/thiago/Code/Pessoal/maitre
echo "orchestrator pane: ${HERDR_PANE_ID:-none}"     # your pane id, if you have one
git log --oneline -5 && git status --short           # where the work stopped
go test ./...                                         # the WIP is green — start from here
cat work/briefs/slice-b-resume.md
```

## Your lanes

The old ones are closed. **Fresh lanes, fresh tabs, per slice** — a reviewer that has
spent the night inside the code has formed a model of it, and a model is what makes
you accept.

```bash
ORCH_PANE="${HERDR_PANE_ID:-none}"
RAW=$(herdr tab create --label "rename: slice B" --cwd "$PWD" --no-focus)
PANE=$(echo "$RAW" | python3 -c "import json,sys; print(json.load(sys.stdin)['result']['root_pane']['pane_id'])")
herdr agent start runner --kind opencode --pane "$PANE"

# dispatch, with your pane id substituted into the brief
sed "s|\$ORCH_PANE|$ORCH_PANE|g" work/briefs/slice-b-resume.md > /tmp/dispatch-b.md
herdr agent prompt runner "$(cat /tmp/dispatch-b.md)"
```

Record each lane's pane id in `work/README.md`. If `$HERDR_PANE_ID` is empty you have
no pane: **the file receipt is then the only channel**, and you read `work/inbox/`
after each dispatch.

## The report protocol — a lane reports twice, the file first

The last runner's turn ended and **nobody heard about it**: its work existed only in a
working tree and a pane's scrollback. So every lane's report has two halves:

```bash
# 1. the receipt — in the repo, where it cannot be lost
mkdir -p work/inbox
cat > work/inbox/runner-slice-b.md <<EOF
Runner: slice B done. Commits: <n>. Per-AC evidence: AC-01..AC-04 stamped.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-b"

# 2. the fast path — a nudge, if you are in a pane
herdr agent prompt "$ORCH_PANE" "Runner: slice B done. Commits: <n>. …"
```

**Check `work/inbox/` when a lane goes quiet.** The receipt does not depend on a pane.

**The orchestrator waits; it does not poll.** The nudge (`herdr agent prompt
"$ORCH_PANE"`) only reaches an **idle** orchestrator — a turn blocked on a shell
`sleep` cannot receive it. So: end the turn, go idle, and read `work/inbox/` on the
way back in. Never re-dispatch into a busy lane to "check" on it: a queued prompt
looks exactly like a stall.

## The rules you inherit

| | |
|---|---|
| **The orchestrator does not write the implementation.** | It writes the paper — the briefs, the rulings, the decision of record. The runner writes the code and the reviewer judges it. If you write code, the review's independence becomes theatre. |
| **A lane commits per acceptance criterion, as it goes.** | That is what makes a stalled lane recoverable. The last one did not, and 38 files sat uncommitted. |
| **A brief is written into `work/briefs/` before it is dispatched.** | Never `/tmp`. |
| **The lanes' scrollback is captured at every slice boundary**, into `work/transcripts/`. | It is bounded and rotates — the last runner's pane held 85 lines of a much longer conversation. |
| **The reviewer judges; it never edits the implementation.** | It files findings in `docs/audit.md` and commits the ledger. |
| **The ledger is the record.** | Its rows keep the names as they were when each finding was filed: it is **not swept**. `docs/audits/` holds the per-cycle reports. |
| **A fix cycle stays with the lane that wrote the code.** | Freshness is for judgment, not for labour. |
| **One fact, one answer.** | No second source of truth. The domain has zero UI imports; `internal/core` and `internal/cli` are stdlib-only; `internal/ui` keeps 28 external modules. |
| **The journal is append-only.** | The 32 historical events keep their old `R01`; the migration is event 33, a `control` event stamped 4. |

## The rulings — do not re-litigate them

They are in `work/README.md`, and they are the chef's decisions: the vocabulary scope
rule, the two binaries, `RCP01` and `BL01`, **evidence** over proof (against
advice — and it makes the code agree with the journal's own `evidence` event kind),
the backlog's three kinds, a small task promoting into a one-criterion recipe, the
append-only journal, `MAITRE_CLI_*` with the old names as a warned fallback, and the
version guard's fingerprint now covering the format facts.

## What is not yours to decide

Five ledger rows are the **chef's** and stay open: **F-56** (the poll's fingerprint),
**F-59** (`design_decisions`), **F-60** (should the plan have questions at all), **F-61**
(`courses/` is named, created, never written) and **P-8**. Two more design questions
are held: **courses** and **brownfield** — the proposals are in
`docs/06-proposals.md`, and v1 already implements all three, so read v1's
`crates/mp/src/{brownfield,decisions,delta}.rs` before proposing anything.

## How to run the remaining slices

For each of B, C, D, E:

1. Dispatch the runner's brief to a **fresh** lane. B uses `slice-b-resume.md`.
2. When the runner reports (inbox file, or the pane message), verify it yourself:
   the commits, the evidence in each commit body, `make lint`, `gofmt -l .`,
   `make test`, `git status` clean.
3. Write the reviewer's brief (the pattern is `review-a2.md`), dispatch it to a
   **fresh** reviewer, and **tighten** it after the first cycle: blockers only, no
   nits, no re-filing of what is closed.
4. Rule on the findings yourself, then dispatch a focused fix-up to the same runner.
5. When the review says `ok`, capture both lanes' transcripts into
   `work/transcripts/`, update `work/README.md`'s state table, and commit.

Then: the mascot for the chef to judge, the glossary rewrite, ADR 0009, and the
final sweep — slice E's brief has it all, including two LOW notes from slice A's
review that are still open.

## What is worth knowing about how the work went

The orchestrator found **none** of the 65 findings of the night; the reviewer found
every one, including the ones the orchestrator's own fix plans had dropped. The
reason is the separation above. Keep it, and the next four slices are cheap.
