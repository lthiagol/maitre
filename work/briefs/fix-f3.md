You are the **runner** — the same lane that wrote F-91 (the healthy one). The
orchestrator's pane id is `w1Z:p1`; use it verbatim in the notify command below.

## The fix cycle for slice F, round 4 — two findings

The reviewer's cycle 3 filed **F-92** and **F-93**. The orchestrator ruled on both:
**accepted**. Two findings, two commits — one each.

**Both are the same shape: two answers to one question.** That shape has now produced
three findings in this slice, so read them as one lesson, not two.

### F-93 (medium) — the record's name fields are enumerated twice

> `actorNames` (the writer's) and `retiredActorNames` (the reader's) — with the same
> doc comment, and `actorNames`' own comment says *"the reader reads them when it
> reads"*, **which it does not**. Two answers to one question, no wire compares them,
> and the drift's consequence is exactly the hole F-91 just closed.

The reviewer proved the consequence: a mutation adding a **name-holding field** the
reader's list knows and the writer's does not → **the writer accepts `line-cook`, and
the reader then refuses the plan the writer just wrote.** F-91's hole, one list
apart.

**The fix — one list, two views.**

| | |
|---|---|
| **One list** | the fields of a record that carry a name. Declared **once**. |
| **Two views** | the **writer's** view (the fields to check before a write) and the **reader's** view (the fields to check when a record is read) — **derived from the one list**, not written out again. |

The two views may differ in *what they do* — the writer refuses, the reader refuses a
plan — but they may not differ in **which fields they know**. That is the whole
finding.

**And the comment that lies.** `actorNames`' comment says the reader reads those
fields when it reads. It does not — it reads a *different* list. Fix the comment, and
make it true: the reader's view is **derived from** the writer's list, so the comment
can say *one list, two views* and mean it.

**A wire compares them.** The one-format wire should check that the two views are
derived from the same list — so a field added to the list reaches both, and a field
known by only one of them fails the suite. **Prove the wire can fail**: add a field
to one view only, and confirm the suite catches it.

### F-92 (medium) — the spec contradicts itself, and the wrong half comes first

> `docs/03` § 1 still says *"a file named `R01-<slug>.json` is read as readily as one
> named `RCP01-<slug>.json`"* and *"a reader that cannot resolve them cannot read the
> past"* — the code **refuses** that file (the file's name is part of the format), and
> the same document says so twice (§ 6 and § 6.2).

**The passage a reader meets first is the wrong one.** § 1 is the entry point: a
reader arrives there and is told a retired file name is read as readily as a current
one — which the tool refuses.

**The fix:** § 1 says what the tool does, and it says it the way § 6 and § 6.2 say
it. The line the rule draws is the one the code draws:

| | |
|---|---|
| **A name typed by a human** | a **lookup**, and it is tolerated: `RCP01`, `rcp01`, `RCP1`, `1` and `R01` all name the recipe the plan holds. |
| **A file's name on disk** | the **format**, and it is refused: a file named `R01-<slug>.json` is a retired form, and the plan is not supported. |
| **The journal** | **history**: its old ids are kept and still resolve. |

*("a reader that cannot resolve them cannot read the past" is true of the **journal**
and false of a **plan's records**. Say which is which.)*

**And fold in the note:** `docs/03` § 6's compressed clause — *"unless the fact is one
a plan carries"* — would read as making slice E's change a version bump. § 6.2 and
`docs/07` state the rule correctly; § 6 must not contradict them.

### The LOW notes — the rulings

| Note | Ruling |
|---|---|
| the record's name fields are hand-written (complete today, nothing checks it) | **folded into F-93** — that is the finding. |
| `normalizeAgentName`'s underscore mapping makes `sous_chef` an agent role the deny-list wire does not name | **leave it, and say why in the report.** The *roles* the code refuses are unchanged; the wire checks the roles the **spec** names, and the spec names `sous-chef`. A wire that named every spelling would be a wire that tracked the spelling rules. |
| the rebuilt sweep command reports the code's own detection and deny-list as drift | **already done by the orchestrator** — the sweep now scans the **docs**, because the code's instrument is the one-format wire. Do not touch `work/`. |
| a file whose name and record disagree is read as a recipe | **leave it** — logged for the ledger's next round, with the orchestrator's agreement. |

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps **28**
  external modules. The domain has **zero UI imports**.
- `core.DomainVersion` is **7**. Do not move it.
- The **journal is append-only** — its history keeps its old names.
- The gate table, the plan format and the command count (**42**) must not move.

## The tree is shared

**Another session is drawing the mascot in this repository** and has files in the tree
(`logo/`). Stage **by name** — never `git add -A`, never `git add .`. Leave anything
you did not write alone. If a file you did not write has changed, **stop and report
it**.

## Not yours — do not touch

| | |
|---|---|
| `docs/audit.md` | the reviewer's to write. |
| `docs/07-vocabulary.md`, `work/` | the orchestrator's. |
| The gate table, the plan format, the command count (42) | nothing here. If a fix would touch them, **stop and report it**. |

## The gate

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

Green.

## Two commits — one per finding

```
Per-AC evidence:
- F-93: one list, two views; the wire fails when a field is known by only one of them; go test ./internal/core -run 'Retired|OneFormat' -v exit 0
```

**One write per file, and commit as soon as the file is coherent.**

## Finish with the receipt — update the one you wrote when you started

```bash
cat > work/inbox/runner-slice-f-fix-1b.md <<EOF
Runner: slice F fix 1b, round 3 done. Commits: <n>. Fixed: F-92, F-93 and the note. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-f fix 1b round 3"
herdr agent prompt "w1Z:p1" "Runner: slice F fix 1b, round 3 done. Commits: <n>. Fixed: F-92, F-93 and the note. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.**
