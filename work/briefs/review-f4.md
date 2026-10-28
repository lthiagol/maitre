You are the **reviewer** — the same lane that reviewed slice F. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 4 — the last round of slice F

The runner fixed F-92, F-93, and then took the gap you and the orchestrator left back.
**Cycle 4 is blockers only**: file **HIGH or MEDIUM** blockers.

**Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's), or your own LOW notes
from cycle 3 (the underscore mapping — the roles the code refuses are unchanged; the
cross-kind file name — logged for the ledger's next round).

### 1. The name check is now by value — and the field list is gone

The orchestrator ruled that the **model does not have to declare its name fields**:
the code needs to know whether **any** field carries a retired name, and that is a
question about the *value*, not the field.

> A **name field**'s value **is** the name — `claimed_by: "line-cook"`, `by: "Line
> Cook"`, `by: "line_cook"`. A **prose field**'s value *mentions* one — *"the
> line-cook wrote this"* — and is therefore **longer than the name**.

**Verify both sides, by command:**

- a field whose **whole value** is a retired name is **refused** — and the error
  **names the field** (`the record's execution.claimed_by carries a retired name`);
- prose that **mentions** a retired name is **not** refused, and there is a test for
  it;
- the comparison uses the tool's **own normalisation** — `Line Cook`, `LINE-COOK`,
  `line_cook` and ` line-cook ` are all the same name;
- **maps and byte slices are not walked**, and the comment says why: the journal's
  **migration events** carry retired *field names* in their `data` (history keeps its
  old names), and a `json.RawMessage` holds a document, not a name.

**And the two views are gone, which was the point.** The writer and the reader run
the **same function** — one answer, two call sites — so the drift F-93 was about
cannot happen, because there is nothing to drift. **The F-93 wire that compared the
two views is gone too.** Check that **nothing was lost** with it: what does the
one-format wire still hold, and is the drift it was built to catch now impossible for
a different reason?

### 2. The one consequence to judge — can a title brick a plan?

The runner flagged it, and it is the right thing to flag:

> *the rule is by value, so a title (or the slug derived from it) of exactly "Line
> Cook" is caught — a recipe with that title would make the plan unsupported.*

**The orchestrator's ruling: the rule stands.** The trade is deliberate — a *loud*
refusal beats a *silent* hole, and a silent hole is the worse failure (this repo has
the lesson: *a silent enum fallback disabled a whole feature, with green tests*).

**But the reviewer's question is the sharp one, and it is yours:**

> **Can a plan be bricked by a title?**

If the writer's check runs **before** the write, a recipe titled *"Line Cook"* is
**refused at creation** — the chef renames it and retries, and the plan was never
unusable. If the record can land first and be refused afterwards, then the tool
**wrote a record it will not read**, which is F-91's shape a third time.

**Prove which of the two it is.** Write a recipe whose title is exactly a retired
name and follow the record: does it land, and does the plan still work afterwards?
**And the same for the slug** — the file's name is derived from the title, so a
recipe whose title slugs to a retired name puts a retired string in the file's name.

### 3. F-92 and F-93 — closed?

**F-92** — `docs/03` § 1 now draws the line the code draws: a name typed by a human
is a **lookup** and is tolerated; a file's name on disk is the **format** and is
refused; the **journal** is history and keeps its old ids. *"A reader that cannot
resolve them cannot read the past"* is now said of the journal. **Check § 1, § 6 and
§ 6.2 agree**, and that the rules-list clause in § 6 no longer reads as making
slice E's change a version bump.

**F-93** — the record's name fields are one list with two views derived from it, and
the comment that lied is gone. (The list is now gone entirely, per the ruling above —
judge the *outcome*, not the intermediate design.)

### 4. Scope, the tree and the ledger

- `core.DomainVersion` is **7** and did not move again. The gate table, the plan
  format, the command count (**42**) and the role names must not have moved.
- **The mascot** — another session's work, in `logo/`. Check the runner committed
  nothing of it.
- **The ledger** — `docs/audit.md` is yours. Move **F-92** and **F-93** to closed,
  with the file and the evidence that holds each, and **add the row for slice F
  itself**: the ruling that this is not brownfield, the tolerance machinery that came
  out, the one-format wire, the one place where a name enters a record, and the check
  that is by value. The ledger is not swept: keep the findings' wording as they were
  filed.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.**

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
cat > work/inbox/reviewer-slice-f-cycle-4.md <<EOF
Reviewer: slice F cycle 4 done. Verdict: <ok|changes-needed>. F-92, F-93 closed: <yes/no>. Can a title brick a plan: <yes/no, and how you proved it>. Filed <N> new blockers (<F-94.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-f cycle 4"
herdr agent prompt "w1Z:p1" "Reviewer: slice F cycle 4 done. Verdict: <ok|changes-needed>. F-92, F-93 closed: <yes/no>. Can a title brick a plan: <yes/no, and how you proved it>. Filed <N> new blockers (<F-94.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>."
```
