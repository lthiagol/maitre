Reviewer: T13 done. Verdict: **ok**. Filed **0** blockers (none — no F-147).

**The skill.** `internal/cli/skills/orchestrator-role/SKILL.md` has a new section, **The surfaces**
(between *The crew* and *The roster's labels*), and it names the CLI: **7** mentions of `maitre-cli`
(it was **0**) and 33 of herdr. It names the three surfaces and says which moves what — the plan's
records move through `maitre-cli` (one verb per operation, so the gates and the audit apply); the
automation — the crew, the lanes, the round, the label — is the **console's**, started from there,
with herdr carrying the crew's state; the automation's own bookkeeping — the briefs, the receipts,
the ledger — stays in the repository's `work/` and `docs/`, never the plan. And it states the rule
flat: **"The CLI never controls the automation."**

**The wire.** `internal/cli/surfaces_test.go` (new) reads the CLI's **own dispatch table** out of
`cli.go` — the same way the count wire does — and holds it two ways: `TestTheDispatchTableIsClassified`
(every group and verb is classified, **both** directions: an unclassified group fails, a classified
group the CLI does not have fails, and the six workspace commands are listed) and
`TestNoVerbControlsTheAutomation` (no group and no verb is made of an automation word —
`agent`, `automation`, `crew`, `herdr`, `label`, `launch`, `pane`, `rename`, `restart`, `round`,
`sous`, `spawn`, `stop`, `tab`). `TestTheAutomationWireBites` proves the check can fail, and
`TestTheDocsSayTheSame` holds `docs/03`'s § 0.1 and `AGENTS.md`'s layer rows to the same words.
**The CLI's verb set did not move** — counted by me out of `cli.go`: **8 groups, 40 verbs, 6
workspace commands** (46 commands). `internal/cli` (non-test) names `herdr` **zero** times and forks
**zero** processes, and the console's crew path writes nothing to the plan: the boundary holds in the
code, not only in the words.

**The recipe's state.** The skill says a recipe's state moves through `maitre-cli` — `recipe
advance`, `ac verify`, `finding resolve` — **never by editing the plan's files**, in three places:
the protocol's **"7 — Close the slice"**, a row in *What you must know*, and two rows in *What you
must not do* ("hand-edit the plan's files", "start the automation from the CLI").

**The three excluded words — the judgement.** `start`, `run` and `kill` are out of the ban list, and
the reasons hold: `step start` is a plan record, a recipe is run by its executor, and `kill` had to
be left out because the group **`skill` contains it** — the ban is a substring match, so a spelling
cannot tell the plan's vocabulary from the runtime's. **Does the classification arm catch a new verb
on its own? Yes** — proved: a verb `grade`, added to `cli.go`, turns the classification arm red
**alone**, with the ban arm green. **But the classification records no judgement**, so the seam is
at the right *place* and the wrong *shape*: `classifiedGroups` says which verbs **exist**, not what
**kind** each one is, and a verb inherits its group's kind silently. Proved: a verb `session-start`
— one that starts a crew, named with words the ban list does not hold — fails the classification arm,
**passes** the ban arm, and the **whole wire is green** once it is classified. The remedy is one
column — verb → kind (*plan record* · *the tool's own file*), with the wire holding that **no verb is
classified as automation** — and then the word list goes and `step start` and `crew start` are told
apart by their kind. **Worth doing when the next verb lands (courses), not now**: today's 40 verbs
are all plan records, and the wire holds them.

**The docs.** `docs/03-domain-model.md` has a new **§ 0.1 The architecture — the three surfaces**: the
three-row table, the two rules ("The CLI never controls the automation", "the plan's records move
through the CLI, never by editing the plan's files"). `AGENTS.md`'s layer table reads `internal/cli` —
the **artifact's** surface, and `internal/ui` — the **automation's** surface. `docs/04`, `docs/05`
and `docs/06` contradict nothing.

**Mutations — every one proved by command, every one reverted; `git status --short` is empty.**

| # | the mutation | what went red |
|:--:|---|---|
| 1 | a verb `crew-start` in the `recipe` group | **both** arms — `TestTheDispatchTableIsClassified` and `TestNoVerbControlsTheAutomation` |
| 1b | a verb `session-start` (a crew verb, no banned word) | the classification arm **alone** — the ban arm stayed green |
| 1c | …and the same verb **classified** | **nothing** — the whole wire went green (the hole the seam judgement is about) |
| 2 | a verb `grade` that is not classified | `TestTheDispatchTableIsClassified` — the classification arm **alone** |
| 3 | the CLI's name out of the skill's *The surfaces* | `TestTheSkillNamesTheSurfaces` — *"the skill should carry \"maitre-cli recipe advance\""* |
| 4 | the skill says the state moves **by editing the plan's files** | `TestTheSkillMovesTheRecipeThroughTheCLI` |
| 5 | the rule out of `docs/03` § 0.1 | `TestTheDocsSayTheSame` |
| 6 | the skill's *The surfaces* **table** reduced to three vague rows | **nothing** — the whole suite green (LOW note) |

**The questions.** `work/questions.md` has no open row: every row is struck (Q8–Q13, Q25, Q28, Q29),
and Q16–Q30 are ruled in `work/tui/README.md`'s rulings table. **Nothing answered, nothing missed,
nothing newly raised** — this slice raises no question.

**DomainVersion:** 7. `internal/core` did not change (the slice's diff touches `AGENTS.md`,
`docs/03`, the skill, `internal/cli/skills_test.go`, `internal/cli/surfaces_test.go` and
`work/tui/README.md` — no golden, no core file). The budget is **0 / 0 / 0 / 28**.

**The gate:** `make check` **green**, `make check-race` **green** — both run by me, not taken from the
receipt.

**LOW notes** (filed in `docs/audit.md`, not blockers):

| | |
|---|---|
| The wire's **ban arm is a spelling**, and the classification table records no **kind** | The proof is mutation 1c above; the remedy is one column (verb → kind) and the wire holding that no verb is classified as automation |
| The skill's **table** — "which moves what" — is held by **nothing** | The wire holds the boundary's *words* in the whole file, while `docs/03`'s wire reads § 0.1 alone; deleting a row is caught, but what each row says is not (mutation 6) |
| The dispatch table is read by a **regexp**, and both arms read it the same way | A group whose shape the regexp does not match would be invisible to both arms; `gofmt` holds the shape today, so the gate closes the hole, not the wire |

**The ledger:** `docs/audit.md` — the T13 row in the cycle matrix, and the three notes above.
Commit `9401313`.
