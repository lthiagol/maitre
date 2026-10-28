# Testing the console

`maitre` is a terminal app, so it has to be testable three ways: by a person, from a
script, and automatically.

---

## 1. By a person

```bash
maitre-cli demo --dir /tmp/kb     # seed a realistic plan
cd /tmp/kb
maitre                        # the real thing, in your terminal
```

`maitre-cli demo` writes seven recipes, six backlog items — four ideas and two small
tasks — and a journal of ~58 events — enough to exercise every tab and every section.

## 2. From a script

The console is driven by keys, so it can be driven by a script too.

```bash
# a static snapshot of every tab
maitre --render --plain

# one recipe
maitre --render --plain --show RCP03

# drive it first, then render — this is the interesting one
maitre --keys "down,down,space" --plain
```

`--keys` takes a comma-separated list: `tab`, `shift+tab`, `enter`, `esc`, `space`,
`up`, `down`, `left`, `right`, or any single character. The console's own commands run, so
a drive that signs off a recipe really signs it off.

```bash
maitre --keys "down,down,space" --plain | head -20
#   "started RCP04" — and the board below it has already moved RCP04 out of "Needs you"
```

That is how you check a behaviour without a terminal, and it is how you paste a repro
into a chat message.

### The stress plan

```bash
maitre-cli demo --chaos --dir /tmp/kb
cd /tmp/kb && maitre
```

`--chaos` adds the pathological recipes: a title with a newline and a tab, a title longer
than the column, thirty acceptance criteria, fifty guidance steps, twelve open tasting
notes, a recipe with no criteria at all, a criterion with no verification, a 400-character
command, and a title with emoji and a right-to-left run.

It is there to **break the console**. Every case is one that has broken a terminal app
before, or would. The invariant the tests check is that **nothing rendered carries a
control character** — see `internal/ui/sanitize_test.go`.

## 3. Automatically

```bash
make test PKG=./internal/ui     # the console's tests
make test                       # everything
```

There are two goldens: `console-recipes.golden` (the normal plan) and
`console-stress.golden` (the pathological one).

The console is a plain `tea.Model`, so the tests drive it headlessly — no terminal, no
pty. `internal/ui/console_test.go` covers the tabs, the ordering keys, the served
checkbox, the primary action, the detail pane, the backlog promotion and the path
grouping.

### The golden

`internal/ui/testdata/console-recipes.golden` pins the Recipes tab's layout.

```bash
make goldens                   # after a deliberate design change
```

A golden catches what unit tests miss: **a column that stops lining up.** When it fails,
read the diff before regenerating — a golden that changed because the design changed is
fine; one that changed because a title grew is a bug.

A golden file that no test writes is a layout nothing pins, and the check that says
"the goldens are current" would pass without touching it. So `make goldens-check`
derives the claim: it marks the goldens, runs the golden tests, and names the ones
nobody wrote. The claim is derived at runtime — no manifest of test names, and no
scan for a path in a test file, because the paths live in constants.

---

## The gate

**One command: `make check`.** The tools are not run by hand.

| Target | What it does |
|---|---|
| `make lint` | `gofmt -l` must be empty, `go vet`, `staticcheck` |
| `make test` | the tests — `PKG=…` and `RUN=…` narrow the run |
| `make test-race` | `-race -shuffle=on -count=2` |
| `make goldens-check` | regenerate the goldens; they must be unchanged, and every one must be claimed by a test |
| `make budget` | the dependency budget, checked against `AGENTS.md` |
| `make tidy-check` | `go mod tidy` must produce no diff |
| `make check` | the fast gate — all of the above |
| `make check-race` | the slow gate — the same, with `-race` |
| `make versions-check` | the pins against the tool the gate runs and the doc that names them |
| `make ci` | what CI runs — `check-race` and `govulncheck` |

A focused check, for one acceptance criterion:

```bash
make test PKG=./internal/core RUN='-run Role|SelfReview'
```

A filter that matches **no** test fails. `go test` alone calls that a pass; a
green that proves nothing is not evidence.

### The tools

```bash
make tools
```

installs two binaries, pinned:

| Tool | Version | Why |
|---|---|---|
| `staticcheck` | `v0.8.2` | `go vet` is weak; staticcheck finds the `SA*`, `S1*` and `ST*` classes vet misses. One binary, no configuration, so it cannot become a project of its own. |
| `govulncheck` | `v1.9.0` | it finds the vulnerabilities that reach the binaries we ship. Not a gate on the code — a gate on the supply chain. |

They are **dev tools**, not module dependencies: `go install` does not touch
`go.mod`, and the budget of 0 / 0 / ≤ 30 is untouched by them.

The pins live in the `Makefile`; the versions named above are the *checked*
copy — `make versions-check` fails when the doc and the pin disagree, the
same shape as the dependency budget. The gate also checks the version of the
tool it runs, so a stale `staticcheck` on the `PATH` fails with `make tools`
rather than quietly running.

A tool that is not installed **fails** the gate with the way to get it. A check
that is quietly skipped is a check that has stopped seeing.

### The hook, and CI

`make install-hooks` copies `scripts/hooks/pre-push.sh` into `.git/hooks/`. It
runs the fast gate, so a broken commit never lands. The slow gate stays in CI
— it takes minutes, and that belongs in CI rather than in a push.

`.github/workflows/ci.yml` runs on every push and every pull request, on the Go
version `go.mod` declares (`go 1.26.0`) — the code must build on the version
it claims. It installs the tools and runs `make ci`, so the list of checks
lives in the `Makefile` and not in the workflow.

---

## Checking the interactive path itself

`--render` and the tests both go through the model, not the terminal. The console uses
the **alt screen**, so what a pty captures is a stream of diffs against a screen buffer
that does not exist in the terminal's scrollback.

To see what the user actually sees, feed the output to a terminal emulator. `pyte` does
this in Python:

```python
import pyte
screen = pyte.Screen(100, 30)
pyte.Stream(screen).feed(raw)
for line in screen.display:
    print(line.rstrip())
```

Cut the capture at the **alt-screen exit** (`ESC[?1049l`). The exit restores the
previous screen, so replaying the whole stream ends up showing the shell instead of the
console:

```python
raw = raw.split("\x1b[?1049l")[0]
```

Stripping the escape codes instead gives a **misleading** picture: the spaces between
columns were written in an earlier frame and never repeated, so a stripped dump looks
broken when the screen is fine.

Set the window size **before** the child starts, or the console sees a different size
than the one you think:

```python
master, slave = os.openpty()
fcntl.ioctl(slave, termios.TIOCSWINSZ, struct.pack("HHHH", rows, cols, 0, 0))
pid = os.fork()
if pid == 0:
    os.setsid()
    os.dup2(slave, 0); os.dup2(slave, 1); os.dup2(slave, 2)
    os.execve("maitre", ["maitre"], dict(os.environ, TERM="xterm-256color"))
```

To drive it under a pty:

```python
import os, pty, select, time

pid, fd = pty.fork()
if pid == 0:
    os.chdir("/tmp/kb")
    os.execve("maitre", ["maitre"], dict(os.environ, TERM="xterm-256color"))
while select.select([fd], [], [], 0.5)[0]:
    os.read(fd, 65536)
os.write(fd, b"q")
```

The output is a stream of *diffs*, not a screen. To read it as a screen, feed it to a
terminal emulator (`pyte` does this in Python) rather than stripping the escape codes —
stripping gives you a misleading picture, because the spaces between columns were
written in an earlier frame and never repeated.

---

## What the tests do not cover yet

| Gap | Why it matters |
|---|---|
| **No live updates under a terminal** | The poll is built and tested at the model level, but the interactive path — raw mode, alt screen, teardown — is only checked by hand |
| **The detail pane cannot scroll** | A long recipe runs off the bottom |
| **Findings cannot be triaged** | The chef can see them, not dismiss them — `TriageFinding` does not exist as an operation |
| **No resize handling** | The layout is computed from the width at load; a resize is not re-laid out |
| **No pty test in CI** | The interactive path — raw mode, alt screen, teardown — is only checked by hand |

Each of these is a slice, not a mystery. They are listed in
[AGENTS.md](../AGENTS.md#open-gaps).
