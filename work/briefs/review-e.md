You are the **reviewer** — the independent reviewer lane, fresh for slice E, the last
slice. The orchestrator's pane id is `w1Z:p1` — it is already substituted in this
file; wherever this brief says `$ORCH_PANE`, it means `w1Z:p1`. Use it verbatim in the
notify command below. The repository is `/Users/thiago/Code/Pessoal/maitre`.

## Your slice — the review of slice E

Slice E is the **paper**: the glossary is rewritten, ADR 0009 records the decision, and
the tree is swept of every old name.

```bash
cd /Users/thiago/Code/Pessoal/maitre
git log --oneline -12
git diff <base>..HEAD --stat      # the base is the commit before slice E's first commit
```

**Cycle 1: blockers only.** File **HIGH or MEDIUM** blockers. No nits, no cosmetics,
no style, no LOW notes filed as rows. **Do not re-file F-56, F-59, F-60, F-61, P-8**
(the chef's).

### The two judgements the runner made — check them against the code, not against taste

**(1) `core.DomainVersion` stayed at 6, while the format's *shape* moved**
(`b84b230161f4ab03`). The runner's reason: the version fences the **writer**, and the
writer did not change — the slice extended the *facts list* (the guard now covers the
root files and the journal envelope), which is not a change in what the tool writes.

**The orchestrator agrees, and the contract says so in the test's own message:**

> *"the domain's rules changed: the shape is now %s, the recorded shape is %s. Update
> this fingerprint. **If the change is in what the tool writes, bump DomainVersion as
> well — the version fences the writer, and the shape alone fences nothing.**"*

**Check three things:** that the shape really did move and the recorded fingerprint was
updated; that **nothing the tool writes** changed (no plan field, no directory, no file
name); and that the version guard still fences in its designed direction — a plan
stamped **7** is refused, a plan stamped **6** is accepted.

**(2) The glossary names each retired word exactly once, in its `_Avoid_` line** — the
pattern slice D set with *Proof* → *evidence*. Check that the retired words are named
**once each**, that the `_Avoid_` line is where they belong (it warns the reader off
the word), and that no retired word is *also* named in prose that is not an `_Avoid_`
line.

### The sweep's remaining 78 lines — check the claim, not the count

The runner reports that **78 lines in 19 files** still carry an old name, and that every
one of them names something that still exists. **Verify that file by file.** The
legitimate ones are the legacy reader (it must name the old directories and the old
rungs), the retired-role aliases and their tests, the tests of the old ids and of the
migration events, the dated version notes, and the tripwires' own mutations. **A line
that names nothing that still exists is drift the sweep missed**, and it is a finding.

### The mascot did not land

The runner reports: the art is in `logo/` (untracked), the console header, the CLI's
help, the README and the Sous-Chef tab carry no drawing, and **AC-01 moves to the
drawing session**. **Check the tree was not touched by the mascot's absence** — the
goldens moved with the sweep alone, and no file of `logo/` was committed.



### The mascot is not this slice's — and the tree is shared

**Another session is drawing the mascot in this repository, in the same working tree.**
Its files are in the tree (`logo/`), and it may write to the plan.

| | |
|---|---|
| **The mascot's files are not the runner's** | the console header, the CLI's help header, the README's art, the Sous-Chef tab's drawing, and the goldens it moves |
| **The sweep must not touch them** | and must not delete or replace what lands. If the sweep would edit a file the mascot owns, **that is a finding to report**, not a change to make |
| **The runner's report must say whether the mascot landed** | the orchestrator closes the slice on that |

Judge the runner's discipline here: **did it commit anything of the other session's
work?** A sweep is exactly the kind of commit that swallows a neighbour's files.

### AC-02 — the glossary is rewritten

`GLOSSARY.md` is the **single source** for the mapping. Check it says what the
decision of record says, and that the three docs that name it **agree**:

| | |
|---|---|
| `GLOSSARY.md` | the **scope rule** at the top — *the kitchen names the artifact and the place; plain words name the machine's records and the agents' jobs* — with a short table of which layer speaks which vocabulary |
| `docs/07-vocabulary.md` | the decision of record. **It wins** over the other two. |
| `docs/03-domain-model.md` | the spec, which links to the glossary rather than repeating it |

- **maitre** is the project and the console; `maitre-cli` is the CLI.
- The **backlog item** replaces *special* and *prep note*, with its three kinds.
- The roles are **chef, sous-chef, runner, reviewer, orchestrator**.
- A recipe's record is **evidence**; the recipe-level stamp is **completion**.
- The kitchen words that survive — recipe, kitchen, line, board, service, tasting,
  sign-off, method, course, menu — each with its `_Avoid_` line.
- **The Evidence entry carries the sign-off exception** — a result that matched
  nothing does not count, *and the chef's sign-off is the one exception, and it
  claims no match*. (That was F-79; a glossary that went back to a narrower entry
  would reopen it.)

**The CPD mapping table: one row per CPD term.** Slice C collapsed *idea / note* and
*brief* into a single row and left the split in the `_Avoid_` line. The mapping must
be readable **one term at a time**, and the backlog item's kinds and the CPD terms
must be **reconciled** — with the reconciliation **stated**, not left to be inferred.

### AC-03 — the ADR

- `docs/decisions/0009-vocabulary.md` — the vocabulary decision: the scope rule, the
  rename table, the backlog's shape, and **the alternatives that lost, with the reason
  each lost**.
- `docs/decisions/README.md` (the index) and `0004-naming.md` say plainly that 0004 is
  superseded by 0009 — and **0004 keeps its tables as the record of what it said**.
- `0009` and `docs/07-vocabulary.md` point at each other: the ADR is where the choice
  lives, and `docs/07` is the chef's decision of record.
- **ADR-0006's § 2 is corrected**: it said *"one of three match states"* and listed
  three under a name the code does not use. The code has **four** — `proven` ·
  `rejected` · `unknown` · `signed-off` — and the sign-off is the one that ADR's own
  section is about. **Check the table now says what the code has.**

**The line about ADRs, which the orchestrator settled — judge against it:**

| | |
|---|---|
| **An ADR whose subject is the mechanism** (0002, 0006, 0008) | its *incidental* vocabulary may move — a doc that contradicts the code is a bug |
| **An ADR whose subject is the vocabulary** (0004, 0005) | it keeps its tables as the record of what it said; 0009 supersedes it |

**Check that no ADR's substance was rewritten** — a decision that changed, or a
constraint that was lost, is a blocker.

### AC-04 — the sweep

Every remaining mention of an old name is gone from the docs and the code: `Kitchen
Brigade`, `scale` (the binary), `line-cook`, `taster`, `expo` (the role), `Chopboard`,
`special`, `prep note`, `Proof` (where it means the record — and the record is called
**evidence** now, so a leftover *proof* that means the record is drift too),
`R01`-style ids, `github.com/lthiagol/kitchen-brigade`.

- **The runner must prove it with the exact commands and their output in the commit
  message.** Run them yourself.
- **Three exclusions, and they are records, not drift:** `docs/audit.md` is the
  ledger of the night and is **not swept** — its rows keep the names as they were when
  each finding was filed; `research/00-legacy-analysis.md` is about **v1**, and its
  mentions of v1's vocabulary are history; the vendored `skills/` are third-party and
  **must not be edited** — if they mention an old name, the runner says so in its
  report instead.
- **The mascot's files are a fourth exclusion** (see above).
- **A word that means the *idea* may stay.** `ac.all_proven`, the match states
  (`proven` · `rejected` · `unknown`), and prose that means the idea of a proof are
  **not** drift. Check the sweep did not over-sweep: a word banished where it was
  right is a finding.

### AC-05 — the whole tree, verified

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

- The dependency budget holds: `internal/core` and `internal/cli` stdlib-only;
  `internal/ui` at **28** external modules. Count them the way `AGENTS.md` says to
  count them and report the number.
- The doc tripwires hold: the derived-views table, the gate table, the plan format,
  the command count (**42**).

### The invariants — check them, do not assume them

- The domain layer has **zero UI imports**. One fact, one answer.
- `core.DomainVersion` is **6** — the sweep changes no format, so the version does
  not move.
- The **journal is append-only** — the sweep must not touch `kitchen/activity.json`.
  The historical events keep their old ids and their old field names.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.**

### The ledger

`docs/audit.md` is **yours** to write. The runner did not touch it — confirm that with
`git diff`. Add any NEW blocker as a new row with an id continuing from **F-80**. Do
not re-file what is closed; the ledger is not swept. Commit the ledger.

### Hard rules

- Do **not** edit implementation code. Read-only for everything except
  `docs/audit.md`.
- Revert every mutation you make. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
- **Do not touch `logo/`** — it is another session's work.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-e.md <<EOF
Reviewer: slice E review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-81.. or none>). Sweep checks: <the commands you ran and what they returned>. Mascot: <did it land? yes/no>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-e"
herdr agent prompt "w1Z:p1" "Reviewer: slice E review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-81.. or none>). Sweep checks: <the commands you ran and what they returned>. Mascot: <did it land? yes/no>. LOW notes: <any>."
```
