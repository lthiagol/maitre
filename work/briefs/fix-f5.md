You are the **runner** — the same lane that wrote slice F. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## A tiny round — two pieces of prose that mislead

The review came back **`ok`** for slice F. Two of its LOW notes are **not** filed,
because a doc that misleads is a bug in this repo — and both are cheap.

### 1. `docs/03` § 3.2 — the illustration omits the trade

The rule as stated is **right**. The illustration's clause — *"prose is longer than the
name"* — is what needs the caveat, because it is not always true: **a prose field
whose whole value is a retired name is refused too.** A recipe titled exactly
*"Line Cook"* is refused at creation, with its slug named in the error.

The reviewer proved the whole shape:

- the record never lands — the writer's check runs **before** the write, so the plan
  is never bricked and stays readable and writable;
- a *later* title change to *"Line Cook"* is refused and does not land;
- only a hand edit could put such a record in a plan, and there the reader refuses
  it loudly and the writer refuses it too — **the tool never writes a record it will
  not read.**

**The fix:** the section says the rule and the trade, so a reader knows that the
refusal is about the *value*, not about the field's kind. A title that is exactly a
retired name is refused, and the reason is the value, not that titles are name fields.

### 2. The walk's "maps are not walked" — the justification is wrong

The comment points at **the journal's migration events**, whose `data` is a
`json.RawMessage` — which the *bytes* clause already covers — and the journal is not
walked at all. So the reason it gives does not support the rule.

**The fix:** the comment gives the **right** reason, or none. The rule stands on its
own: a map's values are not a field's value, and a retired name inside a document is
the document's business, not a name in a record. Say it that way — or say only what
is true.

### Not yours — do not touch

| | |
|---|---|
| `docs/audit.md` | the reviewer's. |
| `docs/07-vocabulary.md`, `work/` | the orchestrator's. |
| The code's behaviour | **nothing changes.** This round is prose. If a fix would move behaviour, **stop and report it**. |

## The gate

```bash
make lint && gofmt -l . && make test && git status --short
```

Green. **Two commits — one for each piece of prose.**

## Finish with the receipt — update the one you wrote when you started

```bash
cat > work/inbox/runner-slice-f-fix-1b.md <<EOF
Runner: slice F fix 1b, round 5 done. Commits: <n>. The spec says the trade; the comment says the right reason. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-f fix 1b round 5"
herdr agent prompt "w1Z:p1" "Runner: slice F fix 1b, round 5 done. Commits: <n>. The spec says the trade; the comment says the right reason. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.**
