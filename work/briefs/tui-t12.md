You are a **fresh runner**. The repository is `/Users/thiago/Code/Pessoal/maitre`.
Read `AGENTS.md` first, then this brief. This brief wins where the design and the
docs differ. The orchestrator's pane id is `w1H:p2`.

This slice is **T12 — the render path stops waiting**. It changes **no behaviour the
chef can see**: the Sous-Chef tab still shows the crew, and the live console still
asks the runtime. It changes **when** the console asks, and **how long** a frame
takes.

## The problem, named exactly

The chain, from the reviewer's own ledger note:

```
render.go:41   snapshot()      →  waits
crew.go:54     readCrew()       →  waits
herdr.go:76    exec.Command(herdr)  →  two spawns
```

**A frame of the console waits on a process.** Two `herdr` spawns, each a fork and
an exec, inside the path that draws the screen. The console's own rule — in
`crew.go` — says **a frame must never wait on a process**. The rule is in the code
and the code breaks it.

## The target behaviour

| | |
|---|---|
| **no process spawn in a snapshot** | the render path **never** forks a process. A frame is a frame: it draws what it knows |
| **the live console still asks** | the **poll** asks the runtime, asynchronously, and the tab shows what it last heard. The chef still sees the crew move |
| **the goldens do not need herdr** | the goldens must render **with herdr off the PATH**. A golden that depends on a binary on the machine is a golden that fails on someone else's |
| **one spawn, not two** | the crew read costs **one** process spawn. The two commands that answer "who is the crew" and "what is each one's job" come from **one** spawn |

**The poll's frequency is not this slice's.** It is deliberate — F-140's witness
requires the poll to ask while the tab is visible — and the poll's frequency stays
with the poll. What the two share is the **cause**: the crew read is two spawns. So
this slice fixes the **read** to be one spawn, and the render to spawn **none**.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | **no process spawn in a snapshot.** A frame renders without forking anything. A test holds it: with herdr **off the PATH**, the tab still renders — and the **goldens are unchanged** |
| **AC-02** | the crew read costs **one** spawn, not two. The two answers — who is the crew, and what is each one's job — come from **one** process |
| **AC-03** | the live console still asks: the poll refreshes the crew while the tab is visible, and the tab shows what it last heard. F-140's witness still holds — the mouse's switch and the filter line's switch still reach the runtime, and a click never leaves the tab saying *the automation is not running* while it is |
| **AC-04** | `work/tui/README.md`'s T12 row names this slice's commits. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **What the tab shows** | **unchanged.** The crew's rows, their states, their jobs, the blocked agent's question — all as they are now |
| **The goldens** | **not one moves.** A golden that changed is a screen that changed |
| **`core.DomainVersion`** | stays **7** |
| **`internal/core`** | does not change |
| **The budget** | unchanged: 0 / 0 / 0 / 28 |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **The poll's frequency** | not this slice's. It is deliberate, and the poll's frequency stays with the poll. Change **what** the poll asks, not **how often** |
| **The seeder** | done. The store is landed and closed |
| **`internal/core`** | does not change |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`docs/audit.md`** | the reviewer's |
| **`work/`** | the orchestrator's, except the T12 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t12.md <<'EOF'
Runner: T12 done. Commits: <n>. The snapshot: <what it does instead of forking, and the witness that a frame spawns nothing>. One spawn: <how the crew read costs one process, and what it answers>. The live console: <what the poll asks, and how the tab shows what it last heard>. The goldens: <unchanged, and that they render with herdr off the PATH>. The poll: <that its frequency is unchanged>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t12.md && git commit -q -m "inbox: runner reports tui-t12"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
