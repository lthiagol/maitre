You are the **runner** — the same lane that wrote slice D. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## A tiny fix — two comments that no longer match the code

The slice D review came back **`ok`**, and its cycle 2 left two LOW notes. The
orchestrator ruled: **both are fixed here.** They are comments, not behaviour, so
there is nothing to prove beyond reading them against the code — but a comment that
states a number the code no longer has is a doc that contradicts the code, and this
repo treats that as a bug.

### 1. `internal/core/model.go` — `legacyFingerprints`' comment contradicts itself

The comment now has two paragraphs that disagree:

- The first says the candidates are *"at most eight, and four for a recipe whose scope
  is empty and whose criteria are not."*
- The second says the candidates are *"the product of the two axes: the list forms and
  the field name, at most sixteen."*

The second is right and the first is stale — the first paragraph still describes the
**list axis alone**, and the numbers it gives are no longer the set's numbers. A
reader gets two different sizes from the same comment.

**The fix:** make the comment say one thing. The first paragraph may keep describing
the list axis on its own — *that* axis is at most eight — and the second says the set
is the product of the two axes. The **product's** numbers are the ones that must be
right: *at most sixteen*, and *eight* for a recipe whose scope is empty and whose
criteria are not.

**The completeness argument stays as it is** — the review confirmed it is right. Only
the numbers move.

### 2. `internal/core/gates_test.go` — "an unproven evidence"

The comment reads *"would produce an unproven evidence, and unproven evidence never
satisfies a gate"*. The word **"an"** is wrong, and it is the last one in the tree —
the slice D sweep removed the rest.

**The fix:** make it read like English. The match state's name in the code is
`unknown` (the code's `MatchUnknown`), so use the code's word if the sentence is
about the state.

## The tree is shared — read this before you commit

**Another session is drawing the mascot in this repository** and has files in the tree
(`logo/`).

| | |
|---|---|
| **Stage by name** | name your files explicitly. **Never `git add -A`**, never `git add .` |
| **`logo/`, and anything you did not write** | leave it alone |
| **`kitchen/`** | re-read before you commit; the other session may write to the plan |

**If a file changed that you did not change, stop and report it.**

## The gate

```bash
make lint && gofmt -l . && make test && git status --short
```

Green.

## One commit

```
Per-AC evidence:
- F-78's comment: the candidate set's numbers match the code; go test ./internal/core -run 'Fingerprint|Legacy' -v exit 0
```

## The report — the receipt comes first

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-d-fix-2.md <<EOF
Runner: slice D fix 2 done. Commits: <n>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-d fix 2"
herdr agent prompt "w1Z:p1" "Runner: slice D fix 2 done. Commits: <n>. Notes: <anything the reviewer should know>."
```

**One write per file, and commit as soon as the file is coherent.**
