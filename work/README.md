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
| **B** | the roles (`runner`, `reviewer`, `orchestrator`), the tabs (`Board`, `Sous-Chef`) | 🔄 with the line cook | [briefs/slice-b.md](briefs/slice-b.md) |
| **C** | the backlog — one item, `kind: idea \| feature \| bug` | queued · version → 5 | [briefs/slice-c.md](briefs/slice-c.md) |
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

| | |
|---|---|
| **line cook** | pane `w1R:p6` — writes the code, one commit per acceptance criterion |
| **taster** | pane `w1R:p7` — reviews, files findings in [../docs/audit.md](../docs/audit.md), never edits the implementation |
| **verifier** | pane `w1R:p8` — the taster's own second lane; it dispatched F-49 and F-50 |
| **orchestrator** | pane `w1R:p1` — dispatches, verifies, rules on the findings |

```bash
# dispatch a brief (write it into briefs/ FIRST — never into /tmp)
herdr agent prompt line-cook "$(cat work/briefs/slice-b.md)"

# one bounded read to confirm delivery — the only sanctioned read
sleep 8 && herdr pane read w1R:p6 --source recent-unwrapped --lines 40
```

**A prompt sent while a lane is working queues behind the running turn.** A
re-dispatch into a busy lane looks exactly like a stall. `agent_status`, `revision`
and the cumulative `Loading Nm` spinner are **non-signalling** — they prove nothing.
Only activity lines do.

## How to resume

```bash
cd /Users/thiago/Code/Pessoal/new-project
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
