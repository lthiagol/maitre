#!/bin/sh
# The golden files must be in the commit.
#
# usage: check-golden-commit.sh <dir>
#
# `git diff` does not see an untracked file, so a new golden that nobody
# committed would pass the check — a gate that says green when it checked
# nothing. The files are intent-to-add before the diff, so a new file shows
# up in it as a change. The index is put back as it was afterwards: the
# tree stays clean.
#
# Two diffs, because they see two things: the worktree against the commit
# (a regenerated golden that was not committed) and the index against the
# commit (a golden the commit would drop).
set -u

dir=${1:-}

untracked=$(git ls-files --others --exclude-standard -- "$dir")

for f in $untracked; do
	git add -N -- "$f"
done

status=0
git diff --exit-code HEAD -- "$dir" || status=1
git diff --exit-code --cached HEAD -- "$dir" || status=1

# Put the index back the way it was: the intent-to-add entries go, and the
# files are untracked again. `git rm --cached` removes the entry without
# restoring anything from HEAD, so a staged deletion stays staged.
for f in $untracked; do
	git rm --cached -q --force -- "$f" 2>/dev/null ||
		git update-index --force-remove -- "$f"
done

if [ "$status" -ne 0 ]; then
	if [ -n "$untracked" ]; then
		echo "these golden files are not in the commit:" >&2
		printf '%s\n' "$untracked" >&2
	fi
	echo "the goldens changed — read the diff before committing it." >&2
	exit 1
fi
