You are the **runner** — the same lane that wrote F-92 and F-93. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## One more round — the gap you handed back

You wrote:

> *a new name-carrying field in the model still has to be added to the list by hand —
> once, not twice — and nothing checks the model against the list (the code cannot tell
> a name field from a prose field by reflection); closing that needs the model to
> declare the fields.*

**The orchestrator's answer: the model does not have to declare anything. The list of
*fields* goes away, and the check is by value.**

The reason the list exists is that the code must know **which** fields carry a name.
It does not need to know. It needs to know whether **any** field carries a retired
name. And that is a question about the *value*, not about the field:

| | |
|---|---|
| **A name field** | its value **is** the name: `claimed_by: "line-cook"`, `by: "Line Cook"`, `role: "line_cook"`. The field's whole value is a retired name. |
| **A prose field** | its value *mentions* one: *"the line-cook wrote this"*, *"a line-cook's recipe"*. The value is **longer than the name**. |

**So: refuse a string field whose whole value is a retired name.** The whole value,
after the tool's own normalisation (case, space, underscore → dash). Prose that
mentions a retired name is longer than the name, so it is not refused — and that is
the test, written down:

- `claimed_by: "line-cook"` → **refused**.
- `body: "the line-cook wrote this"` → **not refused**, and there is a test that says
  so.
- `by: "Line Cook"` → **refused** (the normalisation).

**What that buys:**

1. **The list of fields is gone.** There is no list to keep, and a new name-carrying
   field is caught **automatically**, because the scan covers **every** string field
   of the record.
2. **The two views are gone too, and that is the better half.** The writer and the
   reader do the *same* check by the *same* function — so there is **one** function,
   not one list with two views over it. Two call sites, one answer: the drift F-93
   was about cannot happen, because there is nothing to drift.
3. **The error names the field.** Reflection gives the path, so the message says
   *which* field carried the name: `the record's claimed_by carries a retired name`.

**One thing to be careful about:** the scan must not reach into the journal event's
`data` map — the **migration events** carry retired *field names* in their data
(`"from": "proofs"`), and that is **history**, which keeps its old names. A map field
is not a string field, so a scan over the record's string fields skips it; make sure
that is true of your scan and say so in the comment.

**And the rule for the journal's new events stays as it is:** a *new* event may not
carry a retired **name** — its role, its actor. Its data is the record's business.

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

## One commit

```
Per-AC evidence:
- the check is by value: a field whose whole value is a retired name is refused, and prose that mentions one is not; go test ./internal/core -run 'Retired|Name' -v exit 0
```

**One write per file, and commit as soon as the file is coherent.**

## Finish with the receipt — update the one you wrote when you started

```bash
cat > work/inbox/runner-slice-f-fix-1b.md <<EOF
Runner: slice F fix 1b, round 4 done. Commits: <n>. The name check is by value; the field list is gone. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-f fix 1b round 4"
herdr agent prompt "w1Z:p1" "Runner: slice F fix 1b, round 4 done. Commits: <n>. The name check is by value; the field list is gone. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.**
