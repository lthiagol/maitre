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
go test ./internal/ui/          # the console's tests
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
make regen-goldens             # after a deliberate design change
```

A golden catches what unit tests miss: **a column that stops lining up.** When it fails,
read the diff before regenerating — a golden that changed because the design changed is
fine; one that changed because a title grew is a bug.

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
