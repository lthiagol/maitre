Runner: slice H fix 4 done. Commits: 1 (plus this receipt). The gate sees an untracked golden; the orphan message names the remedy. Notes: the commit check diffs twice — the worktree against the commit and the index against the commit — because they see two things; the index is put back with `git rm --cached --force`, not `git reset`, so a staged deletion stays staged.

Per-AC evidence:

- **the gate sees an untracked golden** — `scripts/check-golden-commit.sh` intent-to-adds the untracked goldens before the diff, so a new file shows up in it as a change, and puts the index back afterwards: the tree stays clean.
  **Proved (a):** a new golden that no test claims (`internal/ui/testdata/zz-orphan.golden`) → `make goldens-check` **fails** and names it: *"these golden files are claimed by no test: internal/ui/testdata/zz-orphan.golden / internal/ui/testdata/zz-orphan.golden: no test source mentions it, so no test claims it. / Write a test that pins it, or delete the file: a golden nothing pins is a layout nothing pins."* Reverted.
  **Proved (b):** a **new** golden that a test *claims* but that nobody committed (the golden test temporarily wrote `testdata/console-new.golden`) → the check **fails**: *"these golden files are not in the commit: internal/ui/testdata/console-new.golden / the goldens changed — read the diff before committing it."* Also proved the staged-drop case (the golden staged as removed) → the same failure, naming the file. Both reverted; `git status --short` is clean after each, and the check leaves no intent-to-add behind.
- **the orphan message names the remedy** — it now says which of the two causes it is, and the way out for each:
  - *no test source mentions it* → "Write a test that pins it, or delete the file: a golden nothing pins is a layout nothing pins."
  - *a test source mentions it, but the golden run's filter (-run Golden) does not match the test that writes it* → "Name the test so the filter matches it — its name must carry 'Golden' — or extend the filter so it matches the test."
- **the two carried notes** — the `2000-01-01` mtime an orphan keeps after a failed run is carried as it stands, and the orphan check's scope is the `GOLDENS` variable: the script's comment says so, so the next person knows where to widen it.

The gate: `make check` green (1m14s), `make check-race` green (4m47s), `git status --short` clean. Counts unchanged: 0 / 0 / 28. `core.DomainVersion` is still 7.
