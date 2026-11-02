Runner: T13 done. Commits: 5 — `75c4b96` (AC-01), `0e3b25b` (AC-02), `32fd028` (AC-03),
`4ac2158` (AC-04), `17c8d99` (AC-05). This receipt is a sixth commit and is not one of
them.

**The skill.** `internal/cli/skills/orchestrator-role/SKILL.md` has a new section
**"The surfaces"**, between *The crew* and *The roster's labels*. It names the CLI
(`maitre-cli` — the skill went from **0** mentions to **7**), herdr and the console,
and says which surface moves what: the plan's records move through `maitre-cli`; the
automation — the crew, the lanes, the round, the label — is the console's, and it is
started from **there**; the automation's own bookkeeping (the briefs, the receipts,
the ledger) stays in the repository's `work/` and `docs/`, never the plan. It also
states the rule flat: **"The CLI never controls the automation"** — no verb starts a
crew, stops one, renames a lane's tab, or reads the crew's state.

**The wire.** `internal/cli/surfaces_test.go`, new. It reads the CLI's **own dispatch
table** out of `cli.go` — the same way `internal/cli/doc_test.go`'s count wire does,
with a regexp over the source — and holds it two ways:

| test | what it holds |
|---|---|
| `TestTheDispatchTableIsClassified` | every group and verb of the live table is in `classifiedGroups`, and the table's six workspace commands are in `cliWorkspaceCommands`. A verb that appears without a line fails, and the file says what a verb moves |
| `TestNoVerbControlsTheAutomation` | no group name and no verb name is made of an automation word — `agent`, `automation`, `crew`, `herdr`, `label`, `launch`, `pane`, `rename`, `restart`, `round`, `sous`, `spawn`, `stop`, `tab` |
| `TestTheAutomationWireBites` | the check on a synthetic table — the point is that the check can fail |

**The recipe's state.** The skill says it in the protocol's **"7 — Close the slice"**
step: *"A recipe's state moves through `maitre-cli` — `recipe advance`, `ac verify`,
`finding resolve`, and the rest of the verbs. Never by editing the plan's files."*
It says it twice more: a row in **What you must know** ("the plan's records move
through the CLI") and two rows in **What you must not do** ("hand-edit the plan's
files", "start the automation from the CLI").

**The docs.** `docs/03-domain-model.md` has a new **§ 0.1 The architecture — the
three surfaces**: the same three-row table, the same two rules. `AGENTS.md`'s layer
table now reads `internal/cli` — the **artifact's** surface: `maitre-cli`, the plan's
records, and `internal/ui` — the **automation's** surface: the `maitre` console.

**DomainVersion:** 7 — unchanged. No format moved, and `internal/core` was not
touched.

**The gate:** `make check` green, then `make check-race` green.

**Notes for the reviewer:**

- **The CLI's verb set did not move.** No verb added, none removed; the count wire
  still reads 46 commands / 40 verbs in 8 groups / 6 workspace commands. The console
  is untouched, no golden moved, `internal/core` is untouched.
- **The verb set is now pinned.** `classifiedGroups` in `internal/cli/surfaces_test.go`
  is the whole dispatch table, verb by verb. A new verb — a course's, say — fails the
  wire until someone adds a line saying what it moves. That is the intent: AC-02
  holds the set, it does not move it.
- **Three words are deliberately not in `automationWords`:** `start` (`step start`
  starts a guidance **step**), `run` (a recipe is run by its executor) and `kill`
  (only the runtime offers it — `skill` contains "kill"). A verb named after one of
  them still has to be classified in `classifiedGroups`, which is where that judgement
  is made.
- **One wire beyond the brief:** `TestTheDocsSayTheSame` holds `docs/03`'s § 0.1 and
  `AGENTS.md`'s layer rows to the skill's words — the boundary is written in three
  places, and three places drift. It reads the docs from `internal/cli`, because
  `internal/core` is out of scope.
- The `work/tui/README.md` T13 row names the four commits **before** the row's own
  commit, as the house rule says. I did not touch the "next work" row, the
  "Paused here" paragraph or anything else in `work/` — the brief leaves those to the
  orchestrator.
