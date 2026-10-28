Runner: slice F fix 1b, round 2 done. Commits: 4 (2 work, 2 receipts). Fixed: F-91 and the notes.

F-91 — a retired name cannot enter a live record, by any path:
- **The rule moved to where the name enters the record.** `Store.SaveRecipe` runs
  `checkRecordNames` on every record it writes, so the domain's API cannot write a
  plan it will not read — the writer's fence (`Store.checkWritable`) is the second
  line, not the first.
- **The entries are checked too**, with one guard (`checkName`): the claim, the
  review, the approval, the sign-off, the comment (already there), the **transition's
  context** (`Operations.advance` — the one place every approver passes through, so
  `ApproveRecipe`, `RequestChanges` and `AdvanceWith` are all covered) and the
  **audit role** (at the write's entry, before anything is written).
- **One comparison for every name.** `normalizeAgentName` maps a space *and* an
  underscore to a dash, so `line_cook` is the name `line-cook`. The deny-list and
  the retired-name check answer "is this the same name?" the same way.
- **A new journal event may not carry a retired name**; the journal's *history* keeps
  its old names (it is append-only, and it is history). The tests that need such an
  event write the history by hand now — the tool will not write one.
- The reviewer's probe fails **without writing anything**: the plan still reads, the
  approval is not there, the lifecycle did not move.

The notes:
- **The wire reads the list the code holds.** It carried its own copy of the retired
  names — a second answer to "what is retired". The names now live in one place
  (`RetiredNames`, `retiredDirectories`, `retiredStampField`,
  `retiredProofsField`, `retiredIDForms`, `retiredFamilies`) and the wire reads
  them: the id form is caught by the prefix and a digit (`R01`), and the retired
  stamp is still caught by its **access** (`["evidence"]`), because `evidence` is a
  current word. The reviewer's four probes still fail the wire (checked: `specials`
  in `internal/ui/views.go`, `specials` in `internal/core/model.go`, `"R01"` and
  `LegacyThing` in `internal/core/store.go`, `LegacyThing` in `model.go`).
- **`docs/03` § 6.2 and § 2.1 say why the version moved.** The typo is gone, and the
  facts-list rule says the truth: a change to the facts list leaves the version alone
  *unless the fact is one a plan carries* — the digest's encoding is such a fact,
  because version 6 covered **two** encodings and the version is the only thing that
  says which encoding a plan's digest is in.

Two things the reviewer should know:
- **`docs/07-vocabulary.md` carries the same facts-list sentence** (§ 6.2's rule,
  around line 114). It is the orchestrator's file, so it is left alone.
- **A file whose name and record disagree is still read as a recipe** (a
  `C01-rogue.json` holding `RCP02`). Left as it is, and here is why: the record is
  the truth and the name is derived from it; the *retired* form in a name is refused
  (that is the rule), and a cross-kind name is a different question — the plan is
  untidy, not in a retired format. **Logged as a note for the ledger's next round.**

Two tests changed their premise:
- `TestARetiredNameCannotBuyTheIndependenceGate` writes the bad record by hand — the
  writer refuses it now — and the reader still refuses the plan that carries it.
- `TestAnOldRoleNameStillReads` writes the old event by hand, and
  `TestTheJournalKeepsTheOldNames` with it.

Untouched: the gate table, the plan format, `core.DomainVersion` (7), the command
count (42), the journal's history, `docs/audit.md`, `docs/07-vocabulary.md`,
`work/`. The invariants hold: `internal/core` 0 · `internal/cli` 0 · `internal/ui` 28.

Gate: make lint && gofmt -l . && make test && make regen-goldens — green.
