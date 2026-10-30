You are the **reviewer** for **T1** — the identity. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md` (the house style and the working
rules) first, then [work/tui/plan.md](../tui/plan.md) § 1 for the reasoning, and
[work/tui/README.md](../tui/README.md) for the state and the protocol.

The orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command below.

**You are a fresh lane.** You did not write this code and you did not read it while it
was written. Re-derive the facts from the code — not from the runner's report, and
not from the brief's claims. A brief says what was *asked*; the code says what was
*done*; the diff between them is your job.

## What T1 asked for

The slice's contract is [tui-t1.md](tui-t1.md). Its six acceptance criteria:

| AC | What |
|--:|---|
| **AC-01** | the window title carries the project; **one source** for the name — the header and the window title must not compute it separately |
| **AC-02** | the header's project name stops being a footnote (`styFaint`) and reads as a name |
| **AC-03** | `--project <name>` overrides the name for the session — in the header **and** the window title — and writes **nothing** |
| **AC-04** | the Settings tab shows the plan's root path, formatted for a human |
| **AC-05** | the goldens are expected to be **unchanged** |
| **AC-06** | `work/tui/README.md`'s state row is updated and **committed with the slice** |

## What to check — and how

**Each check below is a command, not a reading.** A check you cannot run, you do not
report as passed.

### 1. The name has one source

The header and the window title must read the same answer.

- Find every place the console computes the project's name. **There must be one.**
  `grep -rn "GetProject" internal/ui/` — count the call sites, and count how many
  places a name is assembled.
- The **header and the window title must show the same name** under the same
  conditions. Prove it: with the flag absent, with the flag present, and with a plan
  whose `project.json` has no `name` (the fallback path — `m.store.Root`).
- A test asserts `m.View().WindowTitle`. **The window title is not in the rendered
  content**, so `--render` cannot see it — a test is the only instrument here. Find
  the test and confirm it would fail if the title went back to `"maitre"`. **Mutate
  the code to prove the test bites**, and revert the mutation.

### 2. `--project` writes nothing

This is the AC that most easily lies.

- The flag is **display only**. It must not touch `project.json`, must not enter the
  journal, and must not change the plan's name on disk.
- **Prove it by command**: run the console with `--project` against a copy of a plan,
  and show that `project.json` and `activity.json` are **byte-for-byte unchanged**
  afterwards (`git diff` on the copy, or `cmp`).
- The flag must reach **all three** constructors — `New`, `NewLive`, `NewScripted`
  (`internal/ui/model.go`) — and `--render` and `--keys` must honour it, because
  they build the same model. Check the `--render` path and the `--keys` path, not
  only the interactive one.
- The default is the plan's own name, so with the flag absent the console looks
  **exactly as it did before**. Prove it: the goldens are the proof (AC-05).

### 3. The header's name reads as a name

- The project's name must no longer be rendered in `styFaint`. Check the code.
- **The rest of the header must not have moved**: the compact mark (`⋈`), the
  product name (`maitre`), and the `·` separator. A slice that renames the project's
  name and also changes the product's name is doing more than it said.
- The header is **one line** — it is the one place the console has no room to draw.
  Check the header at a narrow width.

### 4. The Settings tab shows the plan's path

- The path is the plan's **root**, formatted for a human — `~/Code/Pessoal/maitre`,
  **not** `/Users/thiago/Code/Pessoal/maitre`.
- It is **one more labelled row**; the themes that are already listed are still
  listed.
- Check it at the Settings tab, and check the path is the **root** and not the
  `kitchen/` directory inside it.

### 5. The goldens

- `make goldens-check` is green, and `git diff` after `make goldens` is **empty**.
- The brief says the goldens are expected to be **unchanged** — the golden is
  `StripANSI(m.View().Content)`, and a *style* is not text. **If a golden changed,
  the slice did more than it said.** Say so, and say what changed.
- The goldens are **claimed by a test** and **in the commit** — `make goldens-check`
  checks both ways.

### 6. The scope

| | |
|---|---|
| **`maitre-cli` did not change** | the console's flag is `maitre`'s, not `maitre-cli`'s. `git diff --stat` over `internal/cli/` and `cmd/maitre-cli/` must be **empty** |
| **No tab was added, removed or renamed** | seven tabs, the same names. `grep` the tab list |
| **No domain change** | `internal/core/` is untouched except — **if** the runner needed a fact the domain did not have, the brief says **stop and report**. Check whether it did |
| **`core.DomainVersion` is still 7** | and the plan format did not move |
| **No dependency was added** | `make budget` still reads **0 / 0 / 28** against **0 / 0 / 30** |
| **The domain layer has zero UI imports** | `internal/core` does not import `internal/ui` |
| **No art moved** | the mark's drawings live in `internal/brand` and nowhere else |
| **`docs/audit.md` was not touched by the runner** | the ledger is yours |
| **`docs/05-human-interface.md`** | if the slice made a line of it false, the line was fixed — and the fix is **in the diff**, named in the report |

### 7. The gate

- `make check` is green — the whole gate, not a part of it.
- `make check-race` is green if the time is there; the slice touches the console, and
  the console is a concurrent program.
- The tree is **clean**: `git status --short` is empty, apart from your own ledger
  commit.

### 8. The state row

- AC-06: `work/tui/README.md`'s T1 row carries the state, the **commit ids** and the
  brief link — and it is **in the slice's commits**, not in a later one.
- The row's commit ids are **real** — they exist in the log and they are the
  slice's own commits.

## The mutation checks — this is what makes a review worth something

A green suite proves the code passes the tests it was written with. It does not prove
the tests *can* fail. **Mutate, run, revert** — and report what you mutated and what
the suite said:

| # | Mutation | What must fail |
|--:|---|---|
| 1 | put the window title back to the hardcoded `"maitre"` | the window-title test |
| 2 | make `--project` write the name into `project.json` | the "writes nothing" check |
| 3 | move the header's name back to `styFaint` | the header's style check |
| 4 | make the name's source differ between the header and the title | the one-source check |

**Revert every mutation.** The tree must be exactly as the runner left it, except
your ledger commit. **Confirm the revert** — `git status --short` and `git diff` at the
end. A reviewer that leaves a mutation behind has damaged the slice it reviewed.

## The ledger — `docs/audit.md` is yours

- File **HIGH or MEDIUM** findings. Number them from **F-99** — the rename's last was
  **F-98**.
- **LOW findings: note them in your message, do not file them.** They go in the
  ledger for a later round — **not** another cycle.
- Add the slice's row to the ledger's slice table, and the cycle's row under it.
- **Never sweep the ledger.** Its rows keep the names as they were when each finding
  was filed.

## The verdict

| | |
|---|---|
| No HIGH or MEDIUM findings | `verdict: ok` |
| Any HIGH or MEDIUM finding | `verdict: changes-needed` |

**If the implementation is right but the spec is ambiguous, give it the benefit of
the doubt** — file a LOW note, not a HIGH finding. And do not invent reasons to
fail: a slice that does what it said is `ok`.

## Hard rules

- Do **not** edit implementation code — **except** `docs/audit.md`, which is yours.
- Revert every mutation you make.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
- **One commit** for the ledger. One commit for the receipt.

## Notify the orchestrator

**The receipt comes FIRST** — the file cannot be lost, the nudge is only the bell.

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-tui-t1.md <<EOF
Reviewer: T1 review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-99.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The window title reads: <what it reads>. The name's source: <where it lives>. --project writes nothing: <yes/no, and how you proved it>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t1"
herdr agent prompt "w1H:p1" "Reviewer: T1 review done. Verdict: <ok|changes-needed>. Filed <N> blockers (<F-99.. or none>). Mutation checks: <what you mutated, what the suite said, and that you reverted it>. The window title reads: <what it reads>. --project writes nothing: <yes/no, and how you proved it>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race if you ran it>. LOW notes: <any>."
```

If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.
