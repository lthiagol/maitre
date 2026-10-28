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
| **C** | the backlog — one item, `kind: idea \| feature \| bug` | ✅ **closed** — review `ok` after one fix cycle; F-73..F-77 closed | [briefs/slice-c.md](briefs/slice-c.md) |
| **D** | `Proof` → evidence, the recipe-level stamp → completion | ✅ **closed** — review `ok`; F-78..F-80 closed; two comments tidied | [briefs/slice-d.md](briefs/slice-d.md) |
| **E** | the glossary, ADR 0009, the sweep, two LOW notes | ✅ **closed** — review `ok` after three rounds; F-81..F-86 closed. **AC-01, the mascot, was carried to slice G** | [briefs/slice-e.md](briefs/slice-e.md) |

> **The mascot (slice E's AC-01) is being made in another session**, on 2026-10-27.
> This tree is shared, so slice E's brief says plainly: do not draw it, do not touch
> the console header, the CLI's help header, the README's art or the Sous-Chef tab's
> drawing, and **do not delete or replace** what lands. Two sessions in one tree is a
> collision, and a collision in the mascot's files is a merge nobody wins. The runner
> reports whether the mascot landed; the orchestrator closes the slice.

The plan's format version is **6** (`core.DomainVersion`). It moves with the format:
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
| **The review** | the reviewer judges; the runner writes; the orchestrator rules and dispatches. A fix is reviewed by someone who did not write it. |

## The lanes

The lane names are the **domain roles** they fill — one word, one referent. The lanes
are **created fresh, per slice**, and their pane ids are recorded here.

| | |
|---|---|
| **runner** | writes the code, one commit per acceptance criterion. The previous runner's slice-B turn ended without reporting; its work was green and uncommitted, captured at [state/slice-b-wip.patch](state/slice-b-wip.patch). **The current lane is `w1Z:p8`** — it wrote slices D and the fix cycle, so the fix cycles stay with it. |
| **reviewer** | reviews, files findings in [../docs/audit.md](../docs/audit.md), never edits the implementation. **The current lane is `w1Z:p9`.** |
| **verifier** | the reviewer's own second lane; it dispatched F-49 and F-50 |
| **orchestrator** | dispatches, verifies, rules on the findings, writes the paper. **This one is `w1Z:p1`.** |

### Slice D's findings

| Finding | What | Where it goes |
|---|---|---|
| **F-78** (high) | **the approval digest did not travel with the rename.** The digest hashes the criterion's JSON, and the JSON carries the field *name* — so a plan approved before the rename keeps a digest the tool cannot reproduce, and `immutability.respected` reads the difference as **tampering**: every gated transition of a recipe nobody touched is refused. The repo's own RCP01 stored `8813a266…`, which is *exactly* the old encoding's digest, while the tool recomputed `4984098d…`. | **closed** — the candidate set is the product of the two axes (list forms × field name), the repo's own plan is migrated, and a test holds the standing discipline. |
| **F-79** (medium) | the glossary's **Evidence** entry said *"a verification result a machine produced"* — but a human **sign-off** is stored in the same field and `ac.all_proven` takes it as a pass. The single source of the vocabulary was narrower than the record, and the *"matched nothing is not evidence"* line was **false** for a sign-off. **The wording came from the orchestrator's brief.** | **closed** — the entry names the sign-off exception, and `docs/03` § 5 says the same thing. |
| **F-80** (medium) | `docs/03` § 6 said the reader does not see the records a version-5 tool wrote, two paragraphs below saying it reads them. The version bump is right — the *writer* changed; the stated consequence is wrong. | **closed** — the paragraph says what the reader does. |

### Slice E's findings — the paper, and the best one of the rename

| Finding | What | State |
|---|---|---|
| **F-81** | the console's retired name (*maître d'*) was in **no** line of the glossary, while every other retired word carried it once in an `_Avoid_` line — and the commit's proof command returned **0** for all four words it claimed to check, so it proved nothing. | **closed** — the name is carried, and the proof reproduces: seven names, seven hits, every one an `_Avoid_` line, **0 outside**. |
| **F-82** | the glossary filed the console under **"### The crew"** — six entries, the sixth being the console — while its own new tables said the jobs were five, and ADR-0003 says *"the console is not an agent"*. The heading also spelled it `**Maitre**`, which the vocabulary does not use. | **closed** — the console has its own section, **"### The surfaces"**, spelled `maitre`. |
| **F-83** | *special* and *prep note* were named **twice** — in prose and in the `_Avoid_` line. The prose is where a reader meets a retired word with nothing warning them. | **closed** — the prose names no retired word and points at ADR 0009 § 3. |
| **F-84** | `docs/03` listed the fingerprint's facts **twice** and the two lists **disagreed** — one mentioned the root files, the other did not. The slice's own commit updated one and not the other, and nothing checks prose. | **closed** — both lists agree, and so does the third, in the `DomainVersion` comment. |
| **F-85** | the spec, the decision of record and the `DomainVersion` comment all still said *"a format change moves the version on its own"* — **the coupling this slice split**. The runner's judgement was right; the prose still taught the opposite. | **closed** — the version fences the **writer**, the fingerprint is the **tripwire** that says the facts moved. |
| **F-86** | **the sweep's instrument was blind, not the tree.** The command and the brief's list omitted *maître d'*, so one of the seven renamed terms was outside its reach — the sweep would have passed whatever the tree said about it, because **it never asked**. | **closed** — the command carries all seven names and excludes the glossary's `_Avoid_` lines. **The command is written down below**, because an instrument that lives only in a commit message is one the next reader has to dig for. |

### Slice F's findings — one shape, three times

| Finding | What | State |
|---|---|---|
| **F-87** | a retired name could **buy a gate**: `line-cook` claims, `runner` reviews, and `review.reviewer_independent` passed. The rule had been enforced per *operation* instead of at the *entry*. | **closed** — a retired name cannot enter a live record by any path, however it is spelled. |
| **F-88** | the operations doc still said the retired names *name the role they stood for* — the resolution was gone. | **closed** — the doc says what the code says. |
| **F-89** | the digest's **encoding** changed and neither fence moved: a plan from the previous build passed the version guard and was refused as *tampered* — a lie, since nothing was tampered with. | **closed** — the digest's payload shape is a **format fact**, so the shape moves, so the version moves (**7**), so the version guard gives the honest message. |
| **F-90** | the new "one format" wire scanned one package, exempted a whole file, and omitted the retired id literals — three doors open. | **closed** — it scans the whole tree, its exemptions are per *declaration*, and it reads the code's canonical lists instead of carrying its own copy. |
| **F-91** | the domain's own API could **write a record it then refused to read** — the tool bricked the plan it had just written. | **closed** — the rule moved to where a name enters a record, and the writer's check runs **before** the write. |
| **F-92** | `docs/03` contradicted itself: § 1 said a retired file name is read as readily as a current one, while § 6 and § 6.2 said it is refused. | **closed** — § 1 draws the line the code draws: a **typed name** is a lookup and is tolerated, a **file's name on disk** is the format and is refused, the **journal** is history. |
| **F-93** | the record's name fields were enumerated **twice** — the writer's list and the reader's list — with no wire comparing them, and the drift's consequence was F-91's hole. | **closed** — and then **the field list went away entirely**: the check is **by value**, so a new name-carrying field is caught automatically, and the writer and the reader run the *same* function. |

**The shape of all three of the severe ones: two answers to one question.** A rule
enforced per *operation* instead of per *entry*; a digest's *encoding* the version
could not see; a list of *fields* written out twice. The remedy is always the same:
**one answer, and the one that is asked.**

### Slice G's findings — the mark, and the rule it taught

| Finding | What | State |
|---|---|---|
| **F-94** | the art's six rows pushed the Sous-Chef tab's **last event** off the bottom of any terminal under 26 rows — and the tab **cannot scroll**, so at 80×24 the chef simply lost the last thing that happened to the plan. | **closed** — the **decoration yields to the content, in both dimensions**: the width rule and now the height rule. |
| **F-95** | the wire that says "the art is declared only here" read only `.go` files and **skipped the brand package itself** — so a second copy passed, and the runner's claim that the mark's characters appear nowhere else was **false**. | **closed** — the drawings may live only in the canonical declaration; the glyph's scope is a **category** (code, rendered output, art, prose), not an extension. |
| **F-96** | the wire that keeps the mark out of the domain read raw lines, so it was defeated by the import spelling **`gofmt` itself writes** — a single-line import put the brand inside `internal/core` and the wire passed. | **closed** — the wire **parses the imports**, so every spelling is caught. |
| **F-97** | the scope wire exempted the declaration **file**, so a copy pasted *beside* the declaration in `internal/brand/mark.go` passed. | **closed** — the exemption is the **declaration**. |

**The rule slice G taught, which is the one worth keeping: the decoration yields to
the content.** A mark is worth nothing if it costs the chef a line of what happened —
so the mark is drawn only when it fits, in both dimensions, and at 80×24 the tab is
byte-for-byte what it was before the mark.

### The tree is shared — with the mascot's session

Since 2026-10-27 evening, **another session is drawing the mascot in this same
repository.** It has files in the tree (`logo/`) and it may write to the plan.

| | |
|---|---|
| **Stage by name** | never `git add -A`, never `git add .` — a sweep takes the other session's work into your commit |
| **Do not touch** | `logo/`, and anything you did not write. If a file you did not write has changed, **stop and report it** |
| **The journal** | append-only. A command run against the plan writes an event — and an event nobody meant is noise in the record. *(The orchestrator found this the hard way: one `recipe claim` probe wrote a refusal event into the repo's own journal, and it was restored to its committed state.)* |

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

### Carried notes — small, real, and not this cycle's business

| | |
|---|---|
| **A refusal that does not carry the reason** | `backlog update` on a dismissed item says *"backlog item BL01 is dismissed"* — it does not say **why** it was dismissed, and the item knows. A LOW from slice C's review; the ledger's next round. |
| **`LegacyItem.Text`** | a prep note's *intent* is shown under the item's `body` field. The two old artifacts kept their prose in different fields and the console shows one thing — a translation, like *tasting note* → `Finding`. Noted, not a defect. |
| **A duplicated doc comment** | `internal/ui/console_test.go`'s `loadModel` carries two doc comments after the F-74 fix — the old one and the new one. Cosmetic; the next lane that touches the file tidies it. |
| **The brief's ambiguous phrase** | slice C's brief said *"the second out-of-scope item"*, which reads two ways. The runner read it the way the code and `docs/03` read it; the brief is now unambiguous. |

### Slice C's findings

| Finding | What | Where it goes |
|---|---|---|
| **F-74** (high) | the live console's `reloadMsg` branch never assigns `m.legacy`, so an old plan's backlog reads as empty. `--render` assigns it, and the test helper hand-copied it — **the test proved the helper, not the path the chef uses**. | **closed** — `Model.applyReload` is now the one place that knows the model's shape, and the test drives the real `Update`. |
| **F-75** (medium) | `reservedRecipeIDs` scans the backlog's pointers only, so a legacy prep note's dangling `promoted_to` did not reserve its id — the same plan reserved `RCP10` before, and now handed out `RCP04`. | **closed** — the legacy reader carries the pointer, and the scan reads it in the format on disk (`"recipe RCP10"`). |
| **F-76** (medium) | `backlog update` wiped the tags when the flag was absent: the CLI's `StringList` returned a non-nil empty slice. The scope list, one line away, was left alone — two list fields, one operation, two behaviours. | **closed** — nil now means "not given", one rule (`setList`) for both fields, and `emptyIfNil` holds the JSON edge. |
| **F-77** (medium) | `docs/03`'s small-task row said a promoted recipe is *falsifiable from the start*; the paragraph below says a criterion with no verification cannot pass `ac.falsifiable`. The code follows the paragraph. **The phrase came from the orchestrator's brief** — the brief is corrected too. | **closed** — a promoted recipe is **checkable, not yet falsifiable**; the code did not move. |

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

## How a lane's turn dies, and how to keep the work

It happened twice on 2026-10-27: slice B's runner and slice D's runner both ended a
turn **without reporting**, mid-work, with changes uncommitted.

| | |
|---|---|
| **What it looks like** | the lane goes **idle** — not blocked, not waiting for a keypress — with no error and no report |
| **Why it happens** | every edit to a large file re-renders that file's **whole diff** into the terminal. A dozen small edits is a dozen full-file diffs, and the turn runs out in the rendering, not in the thinking. |
| **The rule** | **one write per file** — a scripted pass that applies a file's changes at once — and **commit as soon as the file is coherent**. A commit is a landing place. |
| **The net** | `work/state/slice-b-wip.patch` and `work/state/slice-d-wip.patch` exist because of this. `git apply <patch>` restores the work in flight. |

### How a lane's turn dies, and how to keep the work

It happened **four times** on 2026-10-27 and 28 — slice B's runner, slice D's runner,
and slice F's runner twice. Every one of them ended a turn **without reporting**,
mid-work.

| | |
|---|---|
| **What it looks like** | the lane goes **idle** — not blocked, not waiting for a keypress — with no error and no report |
| **Why it happens** | every edit to a large file re-renders that file's **whole diff** into the terminal. A dozen small edits is a dozen full-file diffs, and the turn runs out in the rendering, not in the thinking. The second time, the turn produced **no tool call at all** — one long thought, then nothing. |
| **The rule** | **one write per file** — a scripted pass that applies a file's changes at once — and **commit as soon as the file is coherent**. A commit is a landing place. |
| **The receipt comes first** | the runner now writes its receipt **before** the work, and updates it when it finishes. A receipt written *before* the work is how the orchestrator knows where the lane stopped — and it is visible in the tree while the work is still in flight. |
| **When a lane has failed twice, do not continue it** | a **fresh** lane took over slice F's last third, with a tight brief, and finished it. The failed lane was retired by name (`retired-runner-f`) so the role name was free. |
| **The net** | `work/state/slice-b-wip.patch` and `work/state/slice-d-wip.patch` exist because of this. `git apply <patch>` restores the work in flight. |

**The receipt protocol is the other half.** A turn that ends leaves a working tree and
a pane's scrollback; the receipt in `work/inbox/` is what the orchestrator can always
read. A lane that stops mid-work should still write its receipt, saying where it
stopped.

### The sweep — the command, and what it now returns

The rename's last slice swept the tree of every old name. The **command** is the
instrument, and an instrument that lives only in a commit message is an instrument the
next reader has to dig for — so it is written down here.

```bash
grep -rnE --exclude-dir=.git --exclude-dir=work --exclude-dir=logo --exclude-dir=bin --exclude-dir=kitchen \
  --include="*.md" -i \
  "kitchen brigade|kitchen-brigade|chopboard|\bexpo\b|\btaster\b|line-cook|line cook|maître d'|\bR01\b|\bS01\b|\bN01\b" . \
  | grep -v "^./docs/audit" | grep -v "^./docs/audits/" | grep -v "^./research/" \
  | grep -v "^./docs/decisions/000[45]-" | grep -v "^./docs/decisions/0009-" \
  | grep -v "^./docs/07-vocabulary.md" | grep -v "^./prototypes/" \
  | grep -v "^./GLOSSARY.md:[0-9]*:_Avoid_:"
```

**The code is not in this command's reach, and that is deliberate.** The code's
instrument is the **one-format wire** — it reads the code's own canonical lists, so a
name added to `RetiredNames` is caught in any package without touching the wire. A
grep over the code would only report the two places where the retired names
*legitimately* live: the detection that refuses them and the deny-list that refuses
them as names.

**What it returns: the places that *explain* the refusal, and nothing else.** The
docs' two sections that say the retired names are refused, and one dated line of
version history. **No live value carries a retired name** — the code refuses them
where they enter a record, and the docs say so where a reader meets them.

**The rules the command encodes, and one it cannot:**

| | |
|---|---|
| **Carry every name** | all the retired **literals**. A command that cannot see a name is not a check — that was **F-86**. |
| **The families need a human** | `special`, `proof` and `legacy` are **English words** as well as retired names — *"no special case in a rule"* is not drift. A grep cannot tell the two apart, so the families are **read**, not grepped. The code's own wire (F-90) does distinguish them, and it scans the tree. |
| **A warning is not drift** | the glossary's `_Avoid_` lines name retired words **on purpose**, so they are excluded. |
| **The records are not swept** | the ledger keeps the names as they were when each finding was filed; `research/` is about v1; the old ADRs keep their tables as the record of what they said; `docs/07` and ADR 0009 are the rename's own tables. |
| **No script** | deliberately. The sweep is a **one-time act** over a tree that now has *legitimate* old names — in the code that refuses them, and in the prose that explains it. A permanent tripwire would need a permanent exclusion list, and **a tripwire with a long exclusion list is a tripwire that has stopped seeing**. The live facts have their own instruments: the deny-list wire, the one-format wire, the derived-views wire, the plan-format wire. |

## The chef's ruling — 2026-10-28: **no migrations**

> *"This project at this moment has no worry with brownfield — we're still making it,
> even the tests files we create can be updated to our new env, no worry with
> migrating from old to new. The idea of rewriting this project is exactly this: take
> new decisions and architecture without worrying with migrations or previous
> versions, taking into consideration the lessons learned previously."*

**What it decides.** This is not brownfield. The retired formats are not *supported*;
they are **gone**. There is no upgrade path, no compatibility layer, and no second
writer for a retired format. The rename's **tolerance machinery is the thing to
remove** — it was built to keep old plans working, and the chef says no plan needs
keeping.

**What it does not touch.** The **lessons** are the point of the sentence. Every
tripwire, the append-only journal, one fact one answer, the docs-current rule, and
the fingerprint's job all stay. And the **journal is history**: its events keep their
old names and ids forever, and the reader still resolves what the journal names —
that is not a migration path, that is a log.

| | |
|---|---|
| **The old specials** | the question is **moot**: the retired directories are not read at all. A plan carrying them is not supported, and the tool says so. |
| **The retired role names** | the *refusal* stays (a name like `line-cook` must never be able to sign off); the *resolution* goes (it is no longer "the runner"). |
| **The old digests** | go — and the root cause is fixed instead: the approval digest hashes a **canonical payload with fixed field names**, so a rename never moves it. |
| **The version fence** | one rule: a tool writes a plan only if the plan's version **is** the tool's version. No "a new tool may write an old plan". |

### Slice F — the no-migration slice

| Slice | What | State |
|--:|---|---|
| **F** | the tolerance machinery comes out; the format is one format | ✅ **closed** — review `ok` after four rounds; F-87..F-93 closed | [briefs/slice-f.md](briefs/slice-f.md) |
| **G** | the mark — a bow tie, drawn as a braille halftone | ✅ **closed** — review `ok` after three rounds; F-94..F-97 closed | [briefs/slice-g.md](briefs/slice-g.md) |
| **H** | the quality gate — one interface: `make` | 🟡 **in flight** — the tools, the CI, and the paper that says why | [briefs/slice-h.md](briefs/slice-h.md) |

**The rename is done.** 26 findings found and closed across slices B to G; the
orchestrator found none of them. The ledger's open table holds only the chef's rows.

## The rules that keep this durable

1. **A brief is written into `work/briefs/` before it is dispatched.** Not into `/tmp` —
   `/tmp` took half the night's record once, and the briefs are the reasoning that
   produced the code. A brief in the repo is in the history next to the change it
   caused.
2. **The gate is `make check`.** One command, and every check is a target: a focused
   check is `make test PKG=… RUN=…`, and the tools are not run directly. The gate
   runs where the code lands — in CI, on every push — and the pre-push hook runs the
   fast one. The **slow** gate (`-race -shuffle=on -count=2`, three minutes) stays in
   CI; a push does not wait for it.
2. **The lanes' scrollback is captured at every slice boundary**, into
   `work/transcripts/`. The scrollback is bounded and rotates: the runner's pane
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
