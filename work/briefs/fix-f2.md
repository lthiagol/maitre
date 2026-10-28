You are the **runner** — the same lane that wrote F-90 (the healthy one). The
orchestrator's pane id is `w1Z:p1`; use it verbatim in the notify command below.

## The fix cycle for slice F, round 3 — one blocker

The reviewer's cycle 2 filed **F-91**. The orchestrator ruled: **accepted**. One
finding, one commit — and the LOW notes that ride with it.

### F-91 — the tool can write a record it then refuses to read

> `ops.AdvanceWith(id, core.EventApprove, core.TransitionContext{Approver: "line-cook",
> Actor: "chef"})` **succeeds**, the approval carries `"by": "line-cook"`, and the plan
> is then refused by the reader **and** the writer (*not supported*): the tool bricks
> the plan it just wrote.

**The cause: the retired-name rule is enforced per operation, not where a name enters
a record.** The shipped binaries are safe — the CLI refuses those events with their
own verb, and the console approves through `ApproveRecipe` — so the reach is the
**domain's own API**. But an API that can write a record the domain then refuses is a
domain that can brick a plan.

**The fix — one place, where the name enters the record.**

The retired-name rule belongs at the point where an **actor's name is recorded**: the
claim, the review, the approval, the sign-off, the comment, the transition's context,
and the **audit role**. Not one check per operation — one guard, at the entry.

**The rule to hold:** *a retired name cannot enter a live record, by any path.* So
the domain's API cannot write a plan it will not read, and the writer's fence
(`Store.checkWritable`) is the second line, not the first.

**The audit role is part of this.** `--role line-cook` currently lands in the journal
as an event's role. A **new** event carrying a retired name is a live record — refuse
it, at the same entry point. (The journal's **history** keeps its old names: it is
history.)

**The spelling.** `line_cook` — the same word with an underscore — is a different
string and slips past. The retired name must be refused **however it is spelled**:
case, surrounding space, and the separator (`-`, `_`, space). The tool already
compares names that way for the deny-list; the retired-name check must use the *same*
comparison, so there is one answer to "is this the same name?".

**`RequestChanges` does not check its name** — that is the same hole, and the one
place fixes it.

**Prove it:** every one of these must be refused — the claim, the review, the
approval, the sign-off, the comment, the actor, the audit role — and
`ops.AdvanceWith(..., EventApprove, ...{Approver: "line-cook"})` must fail **without
writing anything**.

### The LOW notes that ride with it

| Note | Ruling |
|---|---|
| `docs/03` § 6.2 has a *"the shape moves the shape"* typo, and its facts-list rule does not match the version bump — the real reason version 6 covered **two digest encodings** | **fix it.** The prose must say why the version moved: it moved because the facts list grew *and* because the digest had two encodings — the encoding was the thing the version could not see. |
| the audit role is not checked | **folded into F-91** — same entry point. |
| `RequestChanges` does not check its name | **folded into F-91** — same guard. |
| `line_cook` is a different name under the tool's spelling rules | **folded into F-91** — one comparison for every name. |
| the retired id form is caught by **literals** only | **fix it: the wire reads the canonical list.** The retired names live in *one* place in the code, and the wire must not restate them — a wire that carries its own copy of the list is a second answer to "what is retired". Read the list the code holds. |
| a file whose name and record disagree is read as a recipe | **leave it, and say why in the report.** The record is the truth and the name is derived from it; the *retired* form in a name is refused (that is the rule), and a cross-kind name is a different question. Log it as a note for the ledger's next round. |
| slice E's sweep command now reports the code's own detection and deny-list as drift | **already done by the orchestrator** — `work/README.md`'s command now carries the **literals** only, and its rules say the families need a human's reading. Do not touch `work/`. |

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps **28**
  external modules. The domain has **zero UI imports**.
- `core.DomainVersion` is **7**. Do not move it.
- The **journal is append-only** — its history keeps its old names; a **new** event
  may not carry one.
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

## Two commits — the guard, then the notes

```
Per-AC evidence:
- F-91: a retired name cannot enter a live record by any path; go test ./internal/core -run 'Retired|Name' -v exit 0
```

**One write per file, and commit as soon as the file is coherent.**

## Finish with the receipt — update the one you wrote when you started

```bash
cat > work/inbox/runner-slice-f-fix-1b.md <<EOF
Runner: slice F fix 1b, round 2 done. Commits: <n>. Fixed: F-91 and the notes. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-f fix 1b round 2"
herdr agent prompt "w1Z:p1" "Runner: slice F fix 1b, round 2 done. Commits: <n>. Fixed: F-91 and the notes. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.**
