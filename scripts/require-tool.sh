#!/bin/sh
# Run a dev tool, or say how to get it.
#
# usage: require-tool.sh <tool> [args...]
#
# The tools are dev tools, not module dependencies: `make tools` installs
# them with `go install`, which never touches go.mod. A missing tool is a
# failed gate, never a skipped check — a gate that quietly skips a check is
# a gate that has stopped seeing.
set -u

name=$1
shift

go_bin=${GO:-go}

found=$(command -v "$name" 2>/dev/null || true)
if [ -z "$found" ]; then
	gopath=$($go_bin env GOPATH)
	if [ -x "$gopath/bin/$name" ]; then
		found="$gopath/bin/$name"
	fi
fi

if [ -z "$found" ]; then
	echo "$name is not installed." >&2
	echo "Install it with: make tools" >&2
	echo "The binaries land in $($go_bin env GOPATH)/bin — put that directory on your PATH." >&2
	exit 1
fi

exec "$found" "$@"
