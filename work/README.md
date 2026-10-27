# The vocabulary rename — the work in progress

The decision is [../docs/07-vocabulary.md](../docs/07-vocabulary.md). This file is
**the state of the work**, so a session that starts cold can pick it up without
re-reading the night's transcripts.

Nothing here is permanent. When the rename is done, `work/` goes away — except this
file, which becomes the ADR's history.

## Where it stands

| Slice | What | State | Brief |
|--:|---|---|---|
| **A** | the project (`maitre`), the binaries (`maitre`, `maitre-cli`), the ids (`RCP01`), the version | ✅ **closed** — review `ok`, 0 blockers | [briefs/slice-a.md](briefs/slice-a.md) |
| **B** | the roles (`runner`, `reviewer`, `orchestrator`), the tabs (`Board`, `Sous-Chef`) | ✅ **closed** — review `ok` after two fix cycles; F-72 closed | [briefs/slice-b-resume.md](briefs/slice-b-resume.md) |
| **C** | the backlog — one item, `kind: idea \| feature \| bug` | 🟡 **in review** — the runner landed it (7 commits, AC-01..AC-06); the review is out | [briefs/slice-c.md](briefs/slice-c.md) |
| **D** | `Proof` → evidence, the recipe-level stamp → completion | queued · version → 6 | [briefs/slice-d.md](briefs/slice-d.md) |
| **E** | the mascot, the glossary, ADR 0009, the sweep, two LOW notes | queued | [briefs/slice-e.md](briefs/slice-e.md) |

The plan's format version is **4** (`core.DomainVersion`). It moves with the format:
each slice that changes the format bumps it, and `docs/03` § 6 says the same number.

## The rulings — do not re-litigate these

| | |
|---|---|
| **The vocabulary** | the kitchen names the artifact and the place; plain words name the machine's records and the agents' jobs. The console's *tasting note* is the code's `Finding` — a translation, not a mismatch. |
| **The binaries** | `maitre` is the console, `maitre-cli` is the CLI, the project is `maitre`. The plan directory stays `kitchen/` — the kitchen is the place. |
| **The ids** | `RCP01` for recipes, `BL01` for backlog items. The backlog's id is kind-neutral, so an idea can be shaped without an id change. |
| **Evidence** | the per-criterion record is **evidence**; the recipe-level stamp is **completion**. The chef's call, against advice — and it makes the code agree with the journal's own `evidence` event kind. |
| **The backlog** | one artifact, three kinds. The console's two groups are **derived** from the kind, never stored. |
| **A small task** | it is too small to be a *big* recipe, not too small to be guarded: it promotes into a recipe with **one criterion**, so every gate still applies. |
| **The journal** | append-only, always. The 32 historical events keep their old `R01`; the migration is event **33**, a `control` event stamped 4. `KindControl` had no writer until then. |
| **The env vars** | `MAITRE_CLI_ACTOR` / `MAITRE_CLI_SESSION` / `MAITRE_CLI_ROLE`. The old `SCALE_*` names still work **with a warning on stderr**; the new name wins when both are set. |
| **The version guard** | the fingerprint now covers the **format facts**, so a format change fails `TestDomainShapeIsVersioned` until the number moves. The format is fenced by the number, not by hand. |
| **The review** | the taster judges; the line cook writes; the orchestrator rules and dispatches. A fix is reviewed by someone who did not write it. |

## The lanes

The lane names are the **domain roles** they fill — one word, one referent. The lanes
are **created fresh, per slice**, and their pane ids are recorded here.

| | |
|---|---|
| **runner** | writes the code, one commit per acceptance criterion. The previous runner's slice-B turn ended without reporting; its work was green and uncommitted, captured at [state/slice-b-wip.patch](state/slice-b-wip.patch). **The current lane is `w1Z:p2`** — it wrote slice B, so the fix cycle stays with it. |
| **reviewer** | reviews, files findings in [../docs/audit.md](../docs/audit.md), never edits the implementation. **The current lane is `w1Z:p4`.** |
| **verifier** | the reviewer's own second lane; it dispatched F-49 and F-50 |
| **orchestrator** | dispatches, verifies, rules on the findings, writes the paper. **This one is `w1Z:p1`.** |

### Slice C's rulings — and two questions for the chef

| | |
|---|---|
| **An older plan's specials** | the console reads them **read-only**, in a third group under the Backlog tab. The tool reads an old plan; it does not rewrite it. Recorded in the decision of record. |
| **No migration verb** | nothing moves a legacy special into the backlog — a verb would move the command count this slice pins at **42**. **Open question for the chef:** is read-only the end state for an old plan's specials, or does the tool owe a migration? Not a blocker for the rename. |
| **The version fence is one-directional** | `internal/core/lock.go` refuses when `plan > DomainVersion`: an **old tool** must not write a plan it does not understand; a **new tool** may write an old plan. The orchestrator's first review brief had this backwards, and the runner's report corrected it — the brief now says judge against the code. |
| **F-73** | fixed — `kitchen/project.json` says `maitre`. The journal took the migration as event **34**, a `control` event stamped 5 — an append, never a rewrite. |

### Two questions for the chef, both from slice C's review

| | |
|---|---|
| **Clearing a list** | `docs/04-operations.md` says *"a list flag replaces the list it names"*, but the parser refuses a list flag with no value (`--tag needs a value`). So a list can be **replaced** but never **emptied** — and `--tag` with no values, by the doc's own rule, should clear the tags. Changing the parser's rule would touch every list flag (`--in-scope`, `--ac`, …), so the fix cycle leaves it alone. **The chef's call.** |
| **The glossary's CPD rows** | slice C's minimal glossary work collapsed two CPD terms (*idea / note* and *brief*) into one row, leaving the split in the `_Avoid_` line. `GLOSSARY.md` is the **single source** for the mapping, and one row per CPD term is its rule. **Slice E rewrites the file** — the rule goes into slice E's brief. |

### Slice C's findings

| Finding | What | Where it goes |
|---|---|---|
| **F-74** (high) | the live console's `reloadMsg` branch never assigns `m.legacy`, so an old plan's backlog reads as empty. `--render` assigns it, and the test helper hand-copied it — **the test proved the helper, not the path the chef uses**. | fix cycle 1 → the runner → cycle 2 |
| **F-75** (medium) | `reservedRecipeIDs` scans the backlog's pointers only, so a legacy prep note's dangling `promoted_to` does not reserve its id — the same plan reserved `RCP10` before, and now hands out `RCP04`. | fix cycle 1 |
| **F-76** (medium) | `backlog update` wiped the tags when the flag was absent: the CLI's `StringList` returns a non-nil empty slice. The scope list, one line away, was left alone — two list fields, one operation, two behaviours. | fix cycle 1 |
| **F-77** (medium) | `docs/03`'s small-task row said a promoted recipe is *falsifiable from the start*; the paragraph below says a criterion with no verification cannot pass `ac.falsifiable`. The code follows the paragraph. **The phrase came from the orchestrator's brief** — the brief is corrected too. | fix cycle 1 |

### Slice B's findings — and where each one lands

| Finding | What | Where it goes |
|---|---|---|
| **F-72** (medium) | `docs/03`'s rule sentence said *the artifact **and the roles** get themed names* while its own table kept the roles plain. The table is right; the sentence was stale. | **closed** — the sentence, the table and `docs/07` now say the same thing. `maître d'` went: it was the console's old themed name, never a job. |
| **F-73** (low) | `kitchen/project.json` says `"name": "new-project"` — the plan's own data, and the console's header renders it. | **slice C** — the slice that migrates the repo's own plan. Not slice B's. |
| LOW notes | the roles in `docs/04`'s prose and the CLI's documented example have no wire; `GLOSSARY.md` still carries the old names; `maître d'` survives in `GLOSSARY.md` and two old ADRs | **slice E** — the sweep owns "every remaining mention of an old name is gone" |

```bash
# a fresh lane, on a fresh tab
RAW=$(herdr tab create --label "rename: slice B" --cwd "$PWD" --no-focus)
PANE=$(echo "$RAW" | python3 -c "import json,sys; print(json.load(sys.stdin)['result']['root_pane']['pane_id'])")
herdr agent start runner --kind opencode --pane "$PANE"
```

### The report protocol — a lane reports **twice**

The last runner's turn ended and **nobody heard about it**: its work was only in a
working tree and a pane's scrollback. So a lane's report now has two halves, and the
**file comes first**:

```bash
# 1. the receipt — in the repo, where it cannot be lost
mkdir -p work/inbox
cat > work/inbox/runner-slice-b.md <<EOF
Runner: slice B done. Commits: <n>. Per-AC evidence: AC-01..AC-04 stamped.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-b"

# 2. the fast path — a nudge, if the orchestrator is in a pane
herdr agent prompt "$ORCH_PANE" "Runner: slice B done. Commits: <n>. …"
```

`$ORCH_PANE` is the orchestrator's pane id, substituted before the brief is
dispatched. **The orchestrator checks `work/inbox/` when a lane goes quiet** — that is
the receipt, and it does not depend on a pane existing at all.

### The orchestrator **waits** — it does not poll

**A nudge only lands while the orchestrator is idle.** A turn blocked on a shell
`sleep` cannot receive it, so the poll and the nudge cancel each other out.

| | |
|---|---|
| **Do** | end the turn and go idle — the lane's `herdr agent prompt "$ORCH_PANE"` arrives as the next turn |
| **Do** | read `work/inbox/` on the next turn, whatever the nudge said — the receipt is the record, the nudge is the bell |
| **Do not** | `sleep` in a shell to wait for a lane — that is polling, and it blocks the only channel that carries the report |
| **Do not** | re-dispatch into a busy lane — a queued prompt looks exactly like a stall |

**A prompt sent while a lane is working queues behind the running turn.** A
re-dispatch into a busy lane looks exactly like a stall. `agent_status`, `revision`
and the cumulative `Loading Nm` spinner are **non-signalling** — they prove nothing.
Only activity lines do.

### A fresh lane per slice

From slice C on, each slice gets a **fresh agent on a fresh tab**, closed when the
slice is reviewed:

```bash
RAW=$(herdr tab create --label "rename: slice C" --cwd "$PWD" --no-focus)
PANE=$(echo "$RAW" | python3 -c "import json,sys; print(json.load(sys.stdin)['result']['root_pane']['pane_id'])")
herdr agent start runner --kind opencode --pane "$PANE"
```

Freshness buys the **reviewer's** independence: a reviewer that has spent the night
inside the code has formed a model of it, and the first blocker of the night came from
checking the reviewer's own example. A fresh reviewer re-derives the facts from the
code instead of from memory of the round before.

The cost is that a fresh agent knows nothing, so the brief must carry what it would
otherwise remember: the invariants, the rulings below, and the ledger's closed rows.

**Two carve-outs.** A *fix cycle* stays with the lane that wrote the code — it holds
the context of what it wrote, and handing a fix to a stranger is re-derivation, not
freshness. And within a slice, the same lane continues: the slice is the unit of work.


## How to resume

```bash
cd /Users/thiago/Code/Pessoal/maitre
git log --oneline -12          # where the work stopped
make lint && make test         # the tree must be green
sed -n '/## Open findings/,/## Closed findings/p' docs/audit.md   # what is open
```

Then: read this file, read the brief for the slice that is not closed, and dispatch it.
If a slice was interrupted mid-flight, its commits are already in the log — read the
last commit's body for the evidence it stamped, and do not re-run what landed.

**To abandon the rename:** `git reset --hard 3236502` returns to the decision of
record, with the night's work intact behind it.

## The rules that keep this durable

1. **A brief is written into `work/briefs/` before it is dispatched.** Not into `/tmp` —
   `/tmp` took half the night's record once, and the briefs are the reasoning that
   produced the code. A brief in the repo is in the history next to the change it
   caused.
2. **The lanes' scrollback is captured at every slice boundary**, into
   `work/transcripts/`. The scrollback is bounded and rotates: the line cook's pane
   held 85 lines of a much longer conversation when it was first captured. The ledger
   holds the findings; the transcript holds how they were found.
3. **The ledger is the record of what was found** — [../docs/audit.md](../docs/audit.md),
   with the per-cycle reports in [../docs/audits/](../docs/audits/). Its rows keep the
   names as they were when each finding was filed: the ledger is not swept.
4. **The findings of the night are not re-litigated by the rename.** F-56, F-59, F-60,
   F-61 and P-8 are the chef's decisions and stay open.
5. **A fresh lane per slice, on a fresh tab** — and the transcript captured before the
   tab is closed. The fix cycle stays with the lane that wrote the code.
6. **The orchestrator does not write the implementation.** It writes the paper: the
   briefs, the rulings, the decision of record. The runner writes the code and the
   reviewer judges it — and if the orchestrator starts writing code, the review's
   independence becomes theatre.
