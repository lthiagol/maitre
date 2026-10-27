You are the **line cook**. The taster reviewed your slice A and returned `changes-needed` with **six findings**. The orchestrator has ruled on each. Fix them — **only these**, nothing else.

The findings are in `docs/audit.md` (F-66 to F-71). Read each row before you touch its code.

## F-66 (medium) — the usage strings still name the removed binary

`maitre-cli recipe show` with a bad id answers `usage: scale recipe show <id>`. Twenty-nine usage lines across `internal/cli/recipe.go`, `fragments.go` and `backlog.go`, plus `cmd/maitre/main.go` telling the user to run `scale init`.

- Every usage line says **`maitre-cli`**. The console's own errors say **`maitre-cli init`** too — the console is `maitre`, and it must not send the chef to a command that does not exist.
- **Sweep, do not patch:** `grep` the tree for `scale ` as a command word (not as an English word) and fix every hit. Report the grep and its output in the commit message.

## F-67 (medium) — the plan-format wire fires in one direction only

`internal/core/planformat_test.go` checks that the code's directories are named in `docs/03` § 6 — not that the directories § 6 names are the code's. Adding `backlog/BL01-<slug>.json` to § 6 leaves the test green, so the doc can promise a directory the reader cannot find. The wire's own comment claims both ways.

- The wire must compare **both sets**: the code's directories and files must be in the doc, **and** the doc's directories and files must be in the code.
- Prove it: mutate the doc (add a directory that does not exist) and mutate the code (rename one that does), and confirm the test fails in both directions. Report both mutations and that you reverted them.

## F-68 (low) — the migrated plan carries no version

No event in `kitchen/activity.json` has a `domain_version`, so the plan reads as version 0, and the guard cannot fence a stale tool off from it. The journal is append-only — it cannot be rewritten to add the stamp.

**The ruling:** the stamp is written by the **tool**, as an **appended** event.

- Append a **`control`** event to the repo's own plan's journal, through the store, recording the migration (the old id, the new id). `KindControl` is declared in the domain and has **no writer at all** today — this gives it its first one, which is exactly what it is for: a control-plane event about the plan itself.
- The event carries the current `domain_version`. The **32 historical events are not touched** — the journal is appended to, never rewritten.
- The same path must exist as a **tool operation**, not a hand edit: the migration is done through the store, so any plan can be stamped the same way. If the operation needs a new CLI verb, say so in your report — do not add one without telling me.
- If `control` turns out to mean something else in the domain's intent, **say so and propose the alternative** rather than guessing.

## F-69 (low) — the identity variables carry the removed binary's name

`SCALE_ACTOR`, `SCALE_SESSION`, `SCALE_ROLE` are printed in the CLI's own help — a surface the rename did reach.

- The new names are **`MAITRE_CLI_ACTOR`**, **`MAITRE_CLI_SESSION`**, **`MAITRE_CLI_ROLE`**.
- The old names are read as a **fallback for a deprecation window**, and the fallback **warns on stderr** so the chef learns the new name. When both are set, **the new name wins**.
- The help, `docs/04-operations.md` and `AGENTS.md` name the new variables.
- A test proves: the new name is read; the old name is read with a warning; the new name wins when both are set.

## F-70 (low) — the old names left the decision records and the sketches

- **ADR 0004** says plainly *which parts* 0009 supersedes, and `AGENTS.md` points at **0009** as the naming decision. (ADR 0009 itself is written in slice E — for now, point at `docs/07-vocabulary.md`, and say so in the ADR.)
- **ADR 0001**'s open question about the old module path is updated.
- The six `prototypes/*.txt` sketches that still say "Kitchen Brigade · demo" say **maitre**.
- The dozen code comments that still name the old things say the new ones.
- **`docs/audit.md` is exempt.** Its rows record the names as they were when the finding was filed; rewriting them would rewrite the record. Leave them and say so in your report.
- `research/00-legacy-analysis.md` is about **v1** — leave it.

## F-71 (low) — nothing ties a format change to a version

`docs/03` § 6 states the plan format's version and no wire checks it, and the fingerprint (`domainShape()`) hashes the *rules* only — so the format can change without the number moving, and a stale tool can write old-format records into a new-format plan.

- `domainShape()` gains the **format facts**: the recipe's id prefix, the plan's directory names, and the JSON field names of the plan's records. Keep it a **small, explicit, named list** — a fact that is not in the list is a fact the stamp cannot see.
- `internal/core/planformat_test.go` checks **`docs/03` § 6's version paragraph against `core.DomainVersion`**.
- Because the shape now covers the format, **the domain version moves: `core.DomainVersion` is 4**, `docs/03` § 6 says 4, and the recorded fingerprint is updated. The version moved because the *rules about the format* changed — say that in the comment.
- Prove it: change a format fact (the id prefix) and confirm `TestDomainShapeIsVersioned` fails until the version is bumped.

## Rules for this cycle

- Fix **only** these six findings. No new features, no refactoring, no renaming of things the findings do not name.
- One commit per finding, with the per-AC evidence stamped in the body:
```
Per-AC evidence:
- F-66: grep -rn 'scale ' internal/cli cmd | wc -l -> 0; go test ./... exit 0 (N/N pass)
```
- **Do not use `git commit --no-verify`. Do not push. Do not touch other panes. Do not spawn agents.**
- Before notifying: `make lint` clean, `gofmt -l .` empty, `make test` green, `git status` clean, and the goldens regenerated if the console's output changed.
- Notify the orchestrator with a **shell command**:
```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-a-fix.md <<EOF
Runner: slice A fix cycle 1 done. Commits: <n>. Fixed: F-66..F-71. Pre-notify gate: pass. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-a-fix"
herdr agent prompt "$ORCH_PANE" "Runner: slice A fix cycle 1 done. Commits: <n>. Fixed: F-66..F-71. Pre-notify gate: pass. Notes: <anything the reviewer should know>."
```
