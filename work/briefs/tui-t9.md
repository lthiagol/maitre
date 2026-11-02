You are a **fresh runner**. The repository is `/Users/thiago/Code/Pessoal/maitre`.
Read `AGENTS.md`, then [work/tui/plan-automation.md](../tui/plan-automation.md) § 15
— the design — and then this brief. This brief wins where they differ. The
orchestrator's pane id is `w1H:p2`.

This slice is **T9 — the settings record**. The console's preferences stop dying
with the console: they are written to a file, in the plan's folder, and read back
when the console starts.

This is the **first of two slices**. **T10** is the Sous-Chef tab, and the
`crew` section of this record comes with it. **This slice carries the `view`
section — the one with a consumer today.** Do not build the `crew` section here.

## The file

```
kitchen/settings.json
```

**The tool owns it, not the domain.** The settings are not a plan *record* — the
glossary's records are the machine's records *about the work* — and they carry no
rules the plan was written under. So:

| | |
|---|---|
| **`internal/core`** | does **not** change. It does not read the file, does not write it, does not know it |
| **`core.DomainVersion`** | stays **7**. `planFormatFacts()` does not move |
| **the repo's own plan** | is **not** stamped. No migration, no `scripts/stamp-plan/` |
| **the plan's readers** | ignore the file. A file in the plan's directory that is not a plan record does not change what a plan is |

The file carries **its own version**, the way the plan's records do and the way the
skill does:

```json
{
  "version": 1,
  "view": { "theme": "Kitchen", "order": "triage", "showServed": true }
}
```

## The package

A new package, **`internal/settings`**, holds the record's schema and the file's
mechanics — the read, the write, and the version check. **stdlib only**, 0 external
modules.

| | |
|---|---|
| **the read** | a missing file yields the **defaults** — it is not an error |
| **the write** | **atomic**: written to a temporary file and renamed, so a half-written file cannot lose the chef's choices |
| **an unknown version** | **refused, by name**. The tool does not guess — the house's *silent enum fallback* lesson: a safe unknown value disabled a whole feature with green tests |
| **unknown keys** | ignored, so an older build reads a newer file |
| **forward compatible** | a file written by a newer tool, with keys this build does not know, still reads |

**The budget moves with the package.** `scripts/dep-budget.sh` holds the numbers and
`AGENTS.md`'s table is checked against them, so **both** need the new row:

```
internal/settings 0 0
```

and in `AGENTS.md`:

```
| `internal/settings` | stdlib only | 0 | 0 |
```

`AGENTS.md`'s **Layout** section gains the package too.

## The CLI

**`maitre-cli settings show`** — the record, as a JSON document, the way every list
the CLI prints is printed.

**One verb, and it is read-only.** The **console** writes the record — it is the
only writer, and it is the one that knows the vocabulary. The CLI reads. (An agent
that wants to know how the plan is set up reads it; the chef changes it in the
console.)

**The command count moves** — one verb, one group. `AGENTS.md` and
`docs/04-operations.md` state the count, and the count is **held** by
`TestTheDocsCommandCountMatchesTheCode`. Move them with the code.

## The console

The **Settings tab's VIEW section** — the theme, the order and show-served — now
**reads the record when the console starts** and **writes it when the chef changes
one of them**. The chef's choices survive a restart.

**One writer, one validator.** The vocabulary — the palette names and the order
names — belongs to the console, because the console is what renders them. So:

| | |
|---|---|
| **the console validates** | a theme or an order it does not know is **refused, by name**, and the console **says so** — on the status line, not in silence |
| **the console does not fall back silently** | an unknown theme means the console says *"the theme X is not one this build has"* and uses the default **with that message** |
| **the file is never written with a value the console does not know** | the console refuses before the write, the way the domain's writer refuses before it writes |

**The plan's own records are untouched by this slice.** The settings are a file
beside the plan's records, not one of them.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | `internal/settings` holds the schema and the mechanics: a missing file yields the defaults, the write is **atomic**, an **unknown version is refused by name**, and unknown keys are ignored so an older build reads a newer file |
| **AC-02** | the file is `kitchen/settings.json`, holding `{"version":1,"view":{...}}`. **`internal/core` does not change** — it does not read the file — and `core.DomainVersion` stays **7**. `planFormatFacts()` does not move. A test holds it |
| **AC-03** | `maitre-cli settings show` prints the record as a JSON document. The command count moves, `AGENTS.md` and `docs/04-operations.md` move with it, and the count wire holds them |
| **AC-04** | the console's Settings tab **reads** the record when it starts and **writes** it when the chef changes the theme, the order or show-served. The chef's choices survive a restart — a test holds it |
| **AC-05** | a theme or an order the console does not know is **refused, by name**, and the console **says so**. It does not fall back silently, and it never writes a value it does not know |
| **AC-06** | `internal/settings` has a row in `scripts/dep-budget.sh` and in `AGENTS.md`'s dependency table — 0 external modules — and `AGENTS.md`'s Layout section names it. `work/tui/README.md`'s T9 row names this slice's commits. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. The plan's format does not move |
| **`internal/core`** | does not change at all |
| **`internal/cli`** | stays **stdlib-only** |
| **`internal/settings`** | **stdlib only**. 0 external modules, ceiling 0 |
| **No new view** | nothing is registered in `core.DerivedViews`, so `docs/03` § 8 does not move |
| **The tab count** | stays **7**. The Sous-Chef tab is **T10** |
| **`kitchen/`** | the repo's own plan is **not** rewritten. Do not stamp it, do not migrate it |
| **Q11–Q13** | closed. Q16–Q29 ruled. Do not reopen them |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **The `crew` section** | it comes with **T10**, the Sous-Chef tab. This slice carries `view` only |
| **The Sous-Chef tab** | **T10**. The console gains no tab in this slice |
| **`internal/core`** | does not change. If you think it must, **report it** — do not change it |
| **The plan's format** | `planFormatFacts()`, `DomainVersion`, `recordedDomainShape` — none of them moves |
| **`kitchen/`** | do not rewrite the repo's plan. The settings file is new, and it lands in the working tree only if a test writes one — the repo's own plan does not get one |
| **`docs/03`, `docs/05`, `docs/06`** | not this slice, except `docs/04-operations.md`'s CLI surface, which the command count moves |
| **`docs/audit.md`** | the reviewer's |
| **`work/tui/plan-automation.md`** | the orchestrator's |
| **`work/`** | the orchestrator's, except the T9 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t9.md <<'EOF'
Runner: T9 done. Commits: <n>. The package: <internal/settings, what it owns>. The file: <kitchen/settings.json, and that core does not read it>. The mechanics: <the read, the atomic write, the unknown-version refusal, the unknown keys>. The verb: <settings show, and the count that moved with it>. The console: <the Settings tab reading and writing, and the chef's choices surviving a restart>. The refusal: <a theme or an order the console does not know, and that it says so>. The budget: <the row in both places, and the layout>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t9.md && git commit -q -m "inbox: runner reports tui-t9"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
