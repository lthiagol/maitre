Runner: slice F fix 1b, round 3 done. Commits: 4 (2 work, 2 receipts). Fixed: F-92, F-93 and the note.

F-93 — the record's name fields are **one list, two views**:
- `nameFields` is the list: the claimant, the approver, the review round's reviewer,
  the comment's author, the record's author, the deviation's author and the three
  overlays' authors. **Declared once.**
- The **writer's view** (`recordNames`) yields every name the list carries;
  `checkRecordNames` refuses a retired one before the write.
- The **reader's view** (`retiredNames`) is the writer's view, **filtered** — the reader
  knows the same fields, because both read the one list. The comment that lied ("the
  reader reads them when it reads") is gone, and the new one is true.
- **The wire compares them.** The one-format wire now checks that the views do not
  reach into the record themselves: which fields carry a name is the list's business,
  and a view that names a field on its own is a second answer. The wire proves it can
  fail — a view that reaches into the record fails the suite (the reviewer's mutation,
  refused).

F-92 — the spec no longer contradicts itself, and the wrong half is gone:
- § 1 now draws the line the code draws: a **name typed by a human** is a **lookup**
  (tolerated), a **file's name on disk** is the **format** (refused — the plan is not
  supported), the **journal** is **history** (its old ids are kept and still resolve).
  "a reader that cannot resolve them cannot read the past" is true of the journal and
  false of a plan's records, and the spec says which is which. The backlog's paragraph
  says the same about `B01-<slug>.json`: a name, not a retired form.
- The note is folded in: § 6's compressed clause ("unless the fact is one a plan
  carries") would have read as making slice E's change — the root files, facts a plan
  carries — a version bump. § 6 and § 6.2 now state the rule the way `docs/07` does:
  **a fact added to the list that the format already had is not a format change; a fact
  the format *changed* is.** (`docs/07` is the orchestrator's and is untouched.)

The LOW notes:
- **`normalizeAgentName`'s underscore** — left as it is, and here is why: the *roles*
  the code refuses are unchanged (agent, sous-chef, runner, reviewer, orchestrator), and
  the wire checks the roles the **spec** names, and the spec names `sous-chef`. A wire
  that named every spelling would be a wire that tracked the spelling rules — and the
  spelling is the domain's business, not the wire's.
- **The sweep** — the orchestrator's, and the code's own instrument is the one-format
  wire. `work/` untouched.
- **A file whose name and record disagree** — left, with the orchestrator's agreement,
  and logged for the ledger's next round.

What the reviewer should know:
- **One gap remains, and it is smaller than it was.** A **new** name-carrying field in
  the model still has to be added to the list by hand — but now **once**, not twice, and
  the wire checks that both views read the list (a view that reaches into the record
  fails the suite). The code cannot tell a name field from a prose field by reflection
  (`Title`, `Description` and `Reason` are strings too), so nothing checks the model
  against the list itself. If the reviewer wants that closed, it needs the model to
  declare the fields — a tag on each — and the list to be derived from it.
- The wire's rule is a source check: the three views (`recordNames`, `retiredNames`,
  `checkRecordNames`) must not touch the record's fields directly.

Untouched: the gate table, the plan format, `core.DomainVersion` (7), the command
count (42), the journal's history, `docs/audit.md`, `docs/07-vocabulary.md`,
`work/`. The invariants hold: `internal/core` 0 · `internal/cli` 0 · `internal/ui` 28.

Gate: make lint && gofmt -l . && make test && make regen-goldens — green.
