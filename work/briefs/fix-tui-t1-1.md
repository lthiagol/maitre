You are the **runner**, again — the same lane that wrote T1. The orchestrator's pane id is
`w1H:p1`; use it verbatim in the notify command below.

Read [work/tui/README.md](../tui/README.md) for the state and the protocol, and
[work/briefs/tui-t1.md](tui-t1.md) for the slice you wrote. This is **T1's fix cycle
1**: three findings the orchestrator ruled on, and one note.

## The review's verdict

**`ok`, 0 blockers.** The four mutations failed the suite and were reverted; the
write-nothing claim was proved by command; the gate is green. **The slice's six
acceptance criteria hold.** This cycle is not a re-do of T1 — it is four small things
the review surfaced, and the orchestrator's ruling on how each one lands.

## F-99 (medium) — the name the console renders is sanitised

**The target:** every string the console renders is stripped of the characters that
corrupt the screen, exactly as the plan's titles, notes and commands are. The project's
name is a rendered string, so it goes through the same door.

**The hole, proved by command.** The console strips control characters out of the
plan's records (`internal/ui/sanitize.go`). The project's name does not go through it,
and it reaches the screen in two places — the header, and (since T1) the **window
title**:

```bash
# the plan's own name, written by the tool itself
mkdir -p "$(mktemp -d)/we
ird" && cd "$(mktemp -d)/we
ird"
maitre-cli init          # writes the project's name from the directory's basename
maitre --render --plain  # the header is TWO rows
```

The name in `project.json` is `we\nird`, and the header prints it across two rows. The
tool wrote that data itself — `Init` takes the directory's basename — so this is the
plan's own content, not a hand-made edge case.

The flag has the same shape, and the review filed it as LOW on the grounds that "the
flag is a human's typing". The flag is one door; the plan's own file is another, and
the second one is the console's stated invariant (slice 6a): *a newline in a title
would corrupt the screen*.

**What lands:** the name the console renders is sanitised the way every other rendered
string is. The file on disk keeps what was written — the sanitising is a rendering
concern.

## F-100 (medium) — the plan's name has one answer

**The target:** one question, one answer, in the domain.

**The hole, proved by command.** "What is this plan called?" is answered in three
places, and the answers disagree:

| The case | The answer today |
|---|---|
| `project.json` names the plan | that name |
| `project.json` is **missing** | the project's directory name — `core.LoadProject` synthesises it |
| `project.json` names it **empty** | `m.store.Root` — the **absolute path** of the plan directory |

So the same plan reads as `maitre` in one state and as `/Users/you/Code/maitre/kitchen`
in another, and the absolute path is what lands in the window title. The review saw the
disagreement and filed it as a LOW note about "two paths naming different things"; it is
the shape the ledger calls the recurring one — *two answers to one question*.

**What lands:** the domain answers the plan's name once, and every reader of it gets
that answer. The console asks the domain; it does not answer the question a second
time, and a **name is never a path** — the header and the window title show a name.

## F-101 (medium) — the row labelled `root` shows the root

**The target:** the Settings tab's `root` row shows the plan's **root** — the plan
directory — formatted for a human.

**The hole.** The row is labelled `root`. The domain's root is the plan directory:
`core.FindRoot` returns it, `store.Root` is it, and the glossary names `kitchen/` *the
plan directory*. The row now shows that directory's **parent**, so `kitchen/` is
nowhere in the console, and a label that says `root` does not name what it shows.

This one is **the brief's fault, and the brief is corrected too.** T1's brief said
"the plan's **root path** — the absolute path of the plan directory" and then gave
`~/Code/Pessoal/maitre` as its example, which is the *project's* directory. The runner
followed the example and reported the deviation; the review followed the example too.
The words win: the root is the plan directory, and the example in the brief was wrong.

**What lands:** the `root` row shows the plan's root, humanised — `~/Code/Pessoal/maitre/kitchen`.

## F-102 (low) — `docs/05` says what the window title carries

**The target:** the doc that sketches the console's screen says what the window title
carries, so a reader is not left to infer it.

Today `docs/05-human-interface.md` says the console "owns the terminal: the alt screen,
its own background, and the window title", and the screen sketch's `row 0` reads
`maitre - <project>`. It does not say that the **window title** — the row the operating
system puts on the terminal's tab — now carries the project's name, and the sketch's
`row 0` is the *header*, which is a different row.

One paragraph in `docs/05-human-interface.md`. The screen sketch stays.

## The scope — what not to touch

| | |
|---|---|
| **The CLI** | `maitre-cli` does not change, and no verb is added. `maitre-cli plan` prints the project; it prints the name the domain now answers, and that is the whole of the change |
| **The domain's writes** | `core.Init` still writes the name it writes today. The change is in what the **reader** answers |
| **`core.DomainVersion`** | stays **7**. The plan format on disk does not change: the same fields, the same values. `TestDomainShapeIsVersioned` is the instrument — if it fails, **stop and report** |
| **The tabs** | seven, the same names. No tab is added, removed or renamed |
| **The goldens** | `make goldens` regenerates; `make goldens-check` holds them. **The goldens are expected to be unchanged** — the header's *words* do not change (the fixture's project is `demo`). If a golden moves, say so and say why |
| **The art** | the mark's drawings live in `internal/brand` and nowhere else |
| **`docs/audit.md`** | the ledger is the **reviewer's**. This cycle files no finding and touches no row |
| **`work/`** | the orchestrator's, except the state row (AC-06) |
| **T2a, T2b** | not started, and this cycle starts nothing |

## The scope — what lands

| | |
|---|---|
| **`internal/core`** | the plan's name has one answer. The change is in the **reader**: what `core.LoadProject` answers when the plan names nothing. The **writer** is untouched |
| **`internal/ui`** | the name the console renders is sanitised; the console asks the domain for the name and shows a name; the `root` row shows the root |
| **`docs/05-human-interface.md`** | one paragraph (F-102) |

**A note on the domain change, since T1's brief said "no domain change".** That was
right for T1 and wrong for this finding: the name's answer is a domain fact, and the
domain is where a fact lives. **The brief's scope is corrected here** — the same way
F-101 corrects it. Keep the change to the reader, and keep it small.

## The invariants

- `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
  budget (**≤ 30** external modules). **No dependency is added by this cycle.**
- **The domain layer has zero UI imports.**
- **`core.DomainVersion` is 7** — this cycle changes no format.
- The tab count is 7 and the tab names do not change.
- **One fact, one answer** — a fact the console shows is computed in the domain.
- The gate is `make check`. A focused check is `make test PKG=… RUN=…`.

## The evidence this cycle must produce

For **each** finding, a command and its output — not a claim:

| Finding | The evidence |
|---|---|
| F-99 | the plan whose directory name carries a newline renders a **one-row** header, and the file on disk still holds what was written |
| F-100 | the plan whose `project.json` names it **empty** and the plan whose `project.json` is **missing** show the **same** answer; and that answer is a name, not a path |
| F-101 | the Settings tab's `root` row shows the plan directory, humanised — `~/Code/Pessoal/maitre/kitchen` — and the themes are still listed |
| F-102 | the paragraph is in `docs/05-human-interface.md`, and it says what the window title carries |

Plus: `make check` green, `make check-race` green, the goldens unchanged, and the tree
clean.

## The report — the receipt comes FIRST

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t1-fix1.md <<EOF
Runner: T1 fix 1 done. Commits: <n>. Fixed: F-99, F-100, F-101, F-102. The name is sanitised: <what you proved it with>. The name has one answer: <what the two cases now show>. The root row shows: <what it shows>. DomainVersion: <7, and the version test's result>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports tui-t1 fix 1"
herdr agent prompt "w1H:p1" "Runner: T1 fix 1 done. Commits: <n>. Fixed: F-99, F-100, F-101, F-102. The name is sanitised: <what you proved it with>. The name has one answer: <what the two cases now show>. The root row shows: <what it shows>. DomainVersion: <7, and the version test's result>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
```

**Write the receipt, commit it, and *then* nudge** — the nudge only reaches an idle
receiver; the receipt is the half that cannot be lost.

**One write per file, and commit as soon as the file is coherent.** Every edit to a
large file re-renders that file's whole diff into the terminal, and a turn that renders
dozens of diffs runs out before the work is done.

## Before you report

```bash
make check && git status --short
```

Green, and the tree clean. Report the **real** numbers: the commits, the gate's
result, and what the two name cases now show.
