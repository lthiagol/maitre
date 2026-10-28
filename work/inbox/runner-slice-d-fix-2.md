Runner: slice D fix 2 done. Commits: 1 (db983d4). Notes: two comments, no behaviour.
Pre-notify gate: pass (make lint clean, gofmt -l . empty, make test green, git status
shows only the other session's `logo/`).

**F-78's comment contradicted itself.** It said the candidates were *"at most eight, and
four for a recipe whose scope is empty and whose criteria are not"*, and then that they
were *"the product of the two axes … at most sixteen"* — two sizes from one comment,
because the first paragraph described the **list axis alone**. The comment now says
one thing: the list axis on its own yields at most eight (four for such a recipe), and
the set — the product of the two axes — is **at most sixteen, and eight** for it. The
completeness argument is untouched; only the numbers moved.

**"An unproven evidence"** in `internal/core/gates_test.go` was the last one in the
tree, and the article was wrong. It now reads: a verification the tool cannot
attribute a count from would produce **evidence in the `unknown` state**, and evidence
in that state never satisfies a gate — the code's own word for the match state
(`MatchUnknown`).

Per-AC evidence:
- F-78's comment: the candidate set's numbers match the code; `go test ./internal/core
  -run 'Fingerprint|Legacy' -v` exit 0.
- The tree has no "an evidence" left outside `docs/audit.md`, which is the reviewer's.
- Staged by name: `internal/core/model.go` and `internal/core/gates_test.go` only.
  `logo/` untouched and uncommitted, and nothing was written to the plan.

Notes: none beyond the above — nothing else moved, and the gate table, the plan
format, the command count (42) and `core.DomainVersion` (6) are as they were.
