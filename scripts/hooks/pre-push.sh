#!/bin/sh
# The fast gate, before a push lands.
#
# The slow gate (-race) stays in CI: it takes minutes, and that belongs in
# CI rather than in a push. `git push --no-verify` skips this hook.
set -e

cd "$(git rev-parse --show-toplevel)"
exec make check
