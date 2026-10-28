Reviewer: slice F cycle 4 done. Verdict: ok. F-92, F-93 closed: yes — both. Can a title
brick a plan: **no** — the writer's check runs *before* the write, so the record never lands.
Filed 0 new blockers (none).

The check is by value, and the field list is gone — proved by command:

- A field whose **whole value** is a retired name is refused, and the error **names the
  field**: `the record carries a retired name: execution.claimed_by ("line-cook")` (the
  writer) and `the record's execution.claimed_by carries a retired name ("line-cook")`
  (the reader) — the same field, from the same function.
- Prose that **mentions** a retired name is **not** refused: "The line-cook wrote this"
  (slug `the-line-cook-wrote-this`) is created and the plan reads; there is a test for it
  (`TestProseThatMentionsARetiredNameIsNotRefused`).
- The comparison uses the tool's **own normalisation**: `line-cook`, `Line Cook`,
  `LINE-COOK`, `line_cook` and `" line-cook "` are all refused.
- **Maps and byte slices are not walked**, and the comment says why ✓ (a test covers it).

**The two views are gone.** The writer (`checkRecordNames`) and the reader
(`retiredNames`) run the same function (`retiredNameFields`), so there is nothing to
drift — F-93's outcome, not its intermediate design. **The F-93 wire that compared the
views went with them, and nothing was lost with it**: the one-format wire still holds the
retired names' two places (per declaration), the canonical lists it reads, and its own
proof that it can fail — my four probes (`specials` in `internal/ui`, `specials` in
`model.go`, `R01` in `store.go`, `LegacyThing` in `store.go`) all still fail it, and a
name added to the code's own `RetiredNames` ("wizard") is caught in `internal/ui` without
touching the wire. The drift the F-93 wire guarded against is impossible for a different
reason: there is no list to drift.

**Can a title brick a plan? No.** A recipe titled exactly "Line Cook" is **refused at
creation** — `the record carries a retired name: title ("Line Cook"), slug ("line-cook")`
— and **nothing is written**: the recipes directory stays empty, the plan still reads and
is still writable, and a backlog item is created afterwards without trouble. The **slug**
is caught with it (the title's slug is "line-cook"), and a *later* title change to "Line
Cook" is refused too — the change does not land and the plan keeps working. Only a hand
edit could put such a record in a plan; there the reader refuses it loudly (*"this plan
carries artifacts of a retired format (the record's title carries a retired name
("Line Cook"))"*) and the writer refuses it too, so the tool never writes a record it will
not read. That is F-91's shape avoided for the third time.

F-92 closed: `docs/03` § 1 now draws the line the code draws — a **name typed by a human**
is a **lookup** and is tolerated; a **file's name on disk** is the **format** and a file
named `R01-<slug>.json` makes the plan **not supported**; the **journal** is **history**
and keeps its old ids. "A reader that cannot resolve them cannot read the past" is now
said of the journal. § 6 and § 6.2 agree with § 1, and § 6's rules-list clause now reads
"a fact added to the list that the format already had is not a format change; a fact the
format *changed* is" — it no longer reads as making slice E's change a version bump.

F-93 closed: the record's name fields are no longer declared at all — the walk is by
value over every string field, so a new name-carrying field is caught without anybody
adding it to a list, and the writer and the reader run the same function.

Scope: `core.DomainVersion` is 7 and did not move; the gate table, the plan format, the
command count (42) and the role names did not move; the journal's history is untouched;
`logo/` is untracked and nothing of it is in the runner's commits; `docs/audit.md` was
untouched by the runner.

Mutation checks (every mutation reverted; the tree is as the runner left it, plus this
ledger): the four wire probes above, plus the `RetiredNames` probe. Every mutation was
reverted and the tree is clean.

LOW notes (not filed): `docs/03` § 3.2 explains the by-value rule with a "name field vs
prose field" framing, so the trade the orchestrator ruled on — a prose field whose *whole*
value is a retired name (a title of exactly "Line Cook") is refused — is not written down;
the rule as stated ("by value", "every string field", "one answer") is right, and the
illustration's clause "prose is longer than the name" is what needs the caveat. The
journal's new-event check is by **field** (`actor`, `role`) while the recipe's is by value
over every string field, so a *session* named after a retired name lands in a new event —
harmless, since the journal is history and nothing reads a session as a name. The walk's
"maps are not walked" justification points at the journal's migration events, whose `data`
is a `json.RawMessage` — bytes, which the byte rule already covers, and the journal is
not walked at all; the rule itself stands on its own. The backlog item is not walked, so
a backlog item titled exactly a retired name is accepted where the recipe's title is
refused — an asymmetry in the trade, with no consequence (nothing is bricked).
