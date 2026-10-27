You are the **runner** — the same lane that wrote slice B and its first fix. The
orchestrator's pane id is `w1Z:p1`; use it verbatim in the notify command below.

## Fix cycle 2 — F-72 is not closed yet

You fixed the rule sentence. **Thank you for reporting the residual difference instead
of picking a side** — that was the right call, and it is why this round is short.

The orchestrator ruled: **F-72 stays open.** The section still contradicts itself — the
sentence now says *the agents' jobs stay plain*, and the table right below it keeps two
jobs in the **Themed** column. The fix and the residual are the same defect, one layer
down.

### The three sources, and which one wins

| Source | What it says about the cast |
|---|---|
| `docs/07-vocabulary.md` § *The rule* — **the decision of record** | The jobs \| plain \| `chef`, `sous-chef`, `runner`, `reviewer`, `orchestrator` — **all plain** |
| `docs/03-domain-model.md`'s table — the spec | Themed \| `chef` · `maître d'` · `sous chef` \| Kept plain \| `runner` · `reviewer` · `orchestrator` |
| `GLOSSARY.md` — the old mapping (slice E rewrites it) | `coordinator → sous chef` · `the human validator → chef` · `the human-facing console → maître d'` |

**The decision of record wins.** `docs/03`'s table is wrong twice over:

1. **It themes two jobs that are plain.** `chef` and `sous chef` are the agents' jobs,
   and `docs/07` puts every job in the plain column. Move them.
2. **It keeps `maître d'`, which was never a job.** Per the old glossary it was *the
   human-facing console*'s themed name. The console is `maitre` now — and it is already
   in the themed column's third row. So `maître d'` does not move: **it goes**.

### The fix — one table, `docs/03-domain-model.md`

| Column | Becomes |
|---|---|
| **Themed** | `recipe` · `special` · `prep note` · `course` · `menu` — **unchanged** |
| **Kept plain** | `chef` · `sous-chef` · `runner` · `reviewer` · `orchestrator` · lifecycle states · gate names · `status` · `review_state` |

Write the two jobs as the code spells them — `chef`, `sous-chef` — not `sous chef`.
The code's role is `sous-chef`; the prose in the spec may say *sous chef* when it is
spoken, but the table is a list of names.

`special` and `prep note` stay: the backlog is slice C's business, and the table
describes the format as it stands today.

### The wire must still hold

`docs/03`'s prose must keep naming every role the code refuses — `agent`,
`sous-chef`, `runner`, `reviewer`, `orchestrator` — and every retired name. The
prose at § the roles already does; the tripwire reads the **whole document**, not the
table, so moving the table's cells must not remove the last mention of anything.

```bash
go test ./internal/core -run 'Doc|DenyList' -v    # the wires still hold
```

If a name would vanish from the document entirely, **that is a signal you are about to
break a wire** — restore the name in the prose, and say so in your report.

### Not yours — do not touch

| | |
|---|---|
| `docs/07-vocabulary.md` | the decision of record. It wins; it does not change. |
| `GLOSSARY.md` | slice E's rewrite. It still carries the old mapping — leave it, and note it. |
| `docs/audit.md` | the reviewer's to write. |
| `docs/decisions/0005-vocabulary.md` | a dated record; slice E's ADR 0009 supersedes it. |
| `internal/core`, `internal/cli`, `internal/ui` | **no code changes in this cycle.** |

### The gate

```bash
make lint && gofmt -l . && make test && git status --short
```

Green, and the tree clean apart from the orchestrator's own `work/` files.

### One commit

One commit for the fix, with the evidence stamped in the body:

```
Per-AC evidence:
- F-72: docs/03's table and its rule sentence now agree with docs/07; go test ./internal/core -run 'Doc|DenyList' -v exit 0
```

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-b-fix-2.md <<EOF
Runner: slice B fix 2 done. Commits: <n>. F-72 closed in docs/03's table. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-b fix 2"
herdr agent prompt "w1Z:p1" "Runner: slice B fix 2 done. Commits: <n>. F-72 closed in docs/03's table. Notes: <anything the reviewer should know>."
```
