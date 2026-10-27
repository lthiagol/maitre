# The inbox

A lane's report lands here **first**, as a file, before it tries to nudge the
orchestrator's pane.

The last runner's turn ended and nobody heard about it — its work existed only in a
working tree and a pane's scrollback. A file in the repo cannot be lost that way.

| File | Written by |
|---|---|
| `runner-<slice>.md` | the runner, when a slice is done |
| `reviewer-<slice>-<cycle>.md` | the reviewer, when a cycle is done |

The orchestrator checks this directory when a lane goes quiet.
