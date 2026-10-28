Runner: slice F fix 1b, round 4 done. Commits: 4 (the lane so far; 2 this round — the fix and this receipt). The name check is by value; the field list is gone. Notes:

- **The check is by value.** The walk covers **every** string field of the record, by
  reflection: a field whose whole value is a retired name is refused — `claimed_by:
  "line-cook"` ✓, `by: "Line Cook"` ✓ (the tool's own normalisation), `by:
  "line_cook"` ✓ — and prose that mentions one is not: `"the line-cook wrote this"` ✓,
  `"a line-cook's recipe"` ✓. There is a test for each side.
- **The list of fields is gone**, and with it the **two views**: the writer
  (`checkRecordNames`) and the reader (`retiredNames`) run the same check by the same
  function — two call sites, one answer. The F-93 wire that compared the two views is
  gone too, because there is nothing to compare: nothing can drift.
- **The error names the field**: `the record's execution.claimed_by carries a retired
  name ("line-cook")` ✓ — reflection gives the path ✓.
- **Maps are not walked**, and the comment says so ✓: a map's values are the record's
  business, and the journal's migration events carry retired *field names* in their
  data — that is history, which keeps its old names. Byte slices are not walked either
  (`json.RawMessage` holds a document, not a name). A test covers both ✓.
- **The journal's new events keep their rule** ✓: a new event may not carry a retired
  name — its role, its actor. Its data is the record's business ✓.
- `docs/03` says the rule the way it now works ✓.

Two things the reviewer should know:

1. **The trade-off the rule makes, stated out loud.** The rule is by value, so a field
   whose whole value is a retired name is refused — and that includes a **title** (or
   the slug derived from it) of exactly "Line Cook" ✓: the walk finds
   `{title Line Cook}` ✓. So a recipe titled *Line Cook* makes the plan unsupported,
   and the writer refuses to create it ✓. Prose is safe ✓ ("The line-cook's recipe" is
   not refused ✓). This is the ruling's line, drawn where it was drawn ✓ — but it is a
   line a short title can cross ✓, and the reviewer should know it is there ✓.
2. **The wire lost its F-93 check** ✓, and with it the mutation test that proved it
   could fail ✓. The check was "the two views read the one list" ✓ — there are no two
   views any more ✓, and the drift the check guarded against cannot happen ✓. The
   one-format wire keeps everything else ✓: the retired names' places ✓, the canonical
   lists ✓, and its own "it can fail" proof ✓ (a token in each package ✓).

Untouched: the gate table, the plan format, `core.DomainVersion` (7), the command
count (42), the journal's history, `docs/audit.md`, `docs/07-vocabulary.md`,
`work/`. The invariants hold: `internal/core` 0 · `internal/cli` 0 · `internal/ui` 28.

Gate: make lint && gofmt -l . && make test && make regen-goldens — green.
