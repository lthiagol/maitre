You are the **runner** — the same lane that wrote slice B. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## The fix cycle for slice B — one finding

The reviewer filed **F-72 (medium)** and the orchestrator ruled on it: **accepted**.
One finding, one fix. Do not touch anything else.

### F-72 — the spec's rule sentence contradicts its own table

`docs/03-domain-model.md` § *The CPD ↔ tool mapping* says:

> The rule behind the split: the **artifact and the roles** get themed names;
> everything an agent reasons about mechanically stays descriptive.

but the table right below it keeps `runner · reviewer · orchestrator` in the
**Kept plain** column. The **table is the right half** — the decision of record
([docs/07-vocabulary.md](../../docs/07-vocabulary.md)) puts the agents' jobs in the
plain vocabulary, and the code says the same. The **sentence is what is stale**, and it
states the opposite of what the section implements, seven lines apart.

**The fix:** rewrite the rule sentence so it agrees with the table and with
`docs/07-vocabulary.md`. Say it in your own words; it must be true of the table as it
now stands — the *artifact and the place* are themed, the *agents' jobs* are plain.

Two things the fix is **not**:

- **No new tripwire.** A wire that parses prose to ask which *column* a role sits in is
  brittle, and the deny-list wire already fires when a role is renamed in the code and
  the spec is not renamed with it. Do not add one.
- **No ADR edit.** `docs/decisions/0005-vocabulary.md` still carries the old wording,
  but an ADR is a dated record of a choice, and slice E writes ADR 0009 to supersede
  it. Say in your report that 0005 is now stale and points at 0009.

**Read `docs/07-vocabulary.md` § *The rule* before you write the sentence.** If the
decision of record and `docs/03`'s table disagree, the decision of record wins and
say so in your report.

### Not yours — do not touch

| | |
|---|---|
| `docs/audit.md` | the reviewer's to write. It files F-72 and closes it. |
| `kitchen/project.json` | F-73 — the plan's own name — is **assigned to slice C**, not here. Leave it. |
| `docs/04-operations.md`, `GLOSSARY.md`, the old ADRs | slice E's sweep. |
| `internal/core`, `internal/cli`, `internal/ui` | **no code changes in this cycle.** The finding is a sentence in a document. |

### The gate

```bash
make lint && gofmt -l . && make test && git status --short
```

Green, and the tree clean apart from the orchestrator's own `work/` files.

### One commit

One commit for the fix, with the evidence stamped in the body:

```
Per-AC evidence:
- F-72: docs/03's rule sentence now matches its table; go test ./internal/core -run 'Doc' -v exit 0
```

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-b-fix-1.md <<EOF
Runner: slice B fix 1 done. Commits: <n>. F-72 fixed in docs/03. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-b fix 1"
herdr agent prompt "w1Z:p1" "Runner: slice B fix 1 done. Commits: <n>. F-72 fixed in docs/03. Notes: <anything the reviewer should know>."
```
