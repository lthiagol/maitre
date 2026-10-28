Runner: slice F fix 1b, round 5 done. Commits: 4 (the lane so far; 3 this round — two pieces of prose and this receipt). The spec says the trade; the comment says the right reason. Notes:

- **docs/03 § 3.2 says the rule and the trade.** The illustration's clause — "prose is
  longer than the name" — is not always true, so it now says "usually", and the trade
  is written out: **the refusal is about the value, not about the field's kind**. A
  prose field whose whole value *is* a retired name is refused too — a recipe titled
  exactly "Line Cook" is refused at creation, and the refusal names the title and the
  slug it would have written.
- I checked the whole shape the reviewer proved, before writing it down ✓:
  - the creation is refused ✓ — "the record carries a retired name: title (\"Line
    Cook\"), slug (\"line-cook\")" ✓ — and **nothing lands** ✓ (the writer checks before
    the write ✓);
  - the plan stays **readable and writable** ✓ — the next recipe is created ✓;
  - a later title change to "Line Cook" is refused ✓ and does not land ✓ (the title on
    disk is unchanged ✓);
  - a hand-written record is refused by the **reader**, loudly ✓ — "the record's title
    carries a retired name (\"Line Cook\")" ✓ — and the writer refuses it too ✓: the
    tool never writes a record it will not read ✓.
- **The walk's "maps are not walked" now says why it stands** ✓: a map's values are
  not a field's value, and a retired name inside a document is the document's business,
  not a name in a record ✓. The old reason — the journal's migration events — did not
  support the rule ✓ (their data is a `json.RawMessage`, which the bytes clause
  covers ✓, and the journal is not walked at all ✓). The spec said the same wrong
  reason ✓, and says the right one now ✓.
- **No behaviour changed** ✓ — this round is prose ✓.

Untouched: the code's behaviour, the gate table, the plan format, `core.DomainVersion`
(7), the command count (42), the journal's history, `docs/audit.md`,
`docs/07-vocabulary.md`, `work/`.

Gate: make lint && gofmt -l . && make test — green.
