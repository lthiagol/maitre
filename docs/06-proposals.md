# 06 — Proposals

Three questions the domain has carried since the beginning. Each is a **decision the chef
makes**, not a decision this document makes: every section states the options, what each
costs, and one recommendation with its reason.

| Proposal | The recommendation in one line |
|---|---|
| [1. Courses](#1-courses) | **Do not build a course yet** — the recipe plus `depends_on` already carries it, and a course file with no writer is the v1 `needs_regrooming` shape |
| [2. `design_decisions`](#2-design_decisions) | **The plan already has this** — a design decision is a *comment* on a recipe, in its thread, where the reader already looks |
| [3. Brownfield](#3-brownfield--delta) | **A brownfield recipe is a greenfield recipe whose first step is a reading** — and its criteria are *regression* criteria, proved by the suite that already exists |

The evidence for all three is the audit's ledger ([audit.md](audit.md)) and the v1 lessons
in [../AGENTS.md](../AGENTS.md). A claim with neither is a preference, and is labelled as
one.

---

## 1. Courses

> *"Does `epic` need its own file type, or is it a tag on a recipe plus a dependency
> graph? v1 had work packages (within a recipe) and never had this rung."* — § 9, question 2

### The question

**What is a course, in this domain's terms** — a file, a tag, or a relation? What does it
hold that a recipe does not?

### The evidence

| What | Where |
|---|---|
| The glossary already names the rung: *"One intent across several recipes — a themed set of dishes served together."* | [GLOSSARY.md](../GLOSSARY.md) |
| `courses/` is **created** by `Init` and named in the plan format — and **no code writes it**. A directory no writer ever touches is the shape v1 died of: *"The field, the transition event, and the renderer all exist — no CLI path ever sets it."* | `internal/core/store.go:63`, [research/00-legacy-analysis.md §5.3](../research/00-legacy-analysis.md) |
| The recipe already has the two things a course would hold: an `Intent` (one outcome) and `DependsOn` (the order). | `internal/core/model.go:26,32` |
| The v1 lesson: a *second* source of truth, reconciled only by accident. | [../AGENTS.md](../AGENTS.md) — *Index drift* |

### The options

| Option | What it costs | What it forecloses |
|---|---|---|
| **A. A course file** — `courses/C01-<slug>.json`, holding a title, an intent, and the recipes in it | A new artifact type: its own file, its own read path, its own place in the ladder's promotion, a place in the console, a place in the plan format. And the **same facts in two places** — the course's list of recipes and each recipe's own `depends_on` would have to agree, or the plan says two things about who waits on whom. | The plan format's promise of **one source of truth per fact**. The moment a course exists, "what waits on what" has two answers. |
| **B. A tag on a recipe** — `course: "signing"` | The smallest change. But a tag answers *which* course, not *what waits on what* — and the ordering the course implies is exactly the information a tag cannot carry. | The ordering. A tag says "these belong together"; it does not say which is served first. |
| **C. The relation** — the recipes, and `depends_on` between them, plus the chef's grouping as a *view* | Nothing in the model. A grouping is a query over the plan, not a thing in it. | A named, shareable unit. "The signing course" cannot be written down, only computed. |

### The recommendation

**Option C — and not to build a course yet.** A course is *the relation*: a set of
recipes that share an intent and an order. The plan already carries the order
(`depends_on`) and the recipe already carries the intent. A course file would add a third
place where the same fact lives, and the v1 lessons say exactly what that costs
(*index drift*, *derived state stored twice* — five implementations of one projection, four
of them wrong).

**What would have to be true to build it.** A chef has to ask "show me the signing
course" and get the same answer twice, or two agents have to name the same set of recipes
and the plan has to check they meant the same one. A *view* over the relation answers the
first; the second needs a named unit. Until then, the honest cost of not building it is
that grouping is a query, not a word.

**What it would change**

| | |
|---|---|
| The plan format | **Nothing** — `courses/` is already created and empty; it would simply stop being named. |
| The lifecycle | **Nothing** — a course is not promoted to and not signed off; the recipes inside it are. |
| The console | A *lane* or a *group* over the Board's tree — a view, not a tab. |
| The CLI | Nothing to write; a filter that reads like `--course` would be a query over `depends_on` and the recipes' titles. |

**What it forecloses.** A course cannot be shared between plans, cannot be signed off as a
unit, and cannot have its own acceptance criteria. If a course ever needs its own evidence
("this set is served together and all three are green"), it becomes an artifact — and that
is the moment to build it.

---

## 2. `design_decisions`

> *"How much of the charter (goals / non-goals / principles) belongs in `project.json`
> versus a human-editable document? CPD principle 2 says living context documents carry
> what no single artifact can hold — `AGENTS.md` and friends. There may be a third place:
> the human interface's own settings."* — § 9, question 7

### The question

**What is a design decision in this domain, and where does one live?** What is it *for* —
what does reading one change about what you do next?

### The evidence

| What | Where |
|---|---|
| The plan format already names `decisions.json` — *"project-level decisions (ADRs)"* — and **no code reads or writes it**. It is a file with no writer today. | `internal/core/store.go:24`, [docs/03 §6](03-domain-model.md) |
| The recipe holds `Audit.Decisions` — a record of `{id, summary, context, at}` — and the journal has a `decision` kind. Neither has a writer that records a *choice*. | `internal/core/model.go:270,291`, `internal/core/activity.go:20` |
| The plan's own rule: **one journal, one event log.** A decision about the plan is an event. | [docs/03 §6.1](03-domain-model.md) |
| v1's `needs_regrooming` — a field, an event and a renderer, **no CLI path** — is the shape to avoid. | [research/00-legacy-analysis.md §5.3](../research/00-legacy-analysis.md) |

### The three things called "a decision" here

| | What it is | Where it lives | Does anything write it? |
|---|---|---|---|
| A **deviation** | The method was done differently than described, and why. | `core.Audit.Deviations` on the recipe, and the journal's `deviation` kind | **Yes** — `RecordDeviation` (`operations.go:480`) |
| A **decision** | A choice the *plan* made about itself. | `core.Audit.Decisions` on the recipe (`model.go:270`), and the journal's `decision` kind | **No.** `Audit.Decisions` has a normaliser and a reader and **no writer**; the journal's `decision` kind is written by `DocumentRecipe` (`operations.go:500`) — which records the *audit trail*, not a choice |
| A **design decision** | A choice about *how to build* — the shape, the trade-off, the alternative rejected. | **Nowhere.** This is the gap. | — |

**The middle row is the finding, not the third.** `Audit.Decisions` is a field with a
shape and no writer — the v1 `needs_regrooming` shape exactly (*"The field, the transition event, and the raul renderer all exist — no
CLI path ever sets it"*, §5.3 ✓), in this codebase, today. So the honest question is not only *"where does a design
decision live?"* but also *"what is `Audit.Decisions` for, and who writes it?"*

### The options

| Option | What it costs | What it forecloses |
|---|---|---|
| **A. A new plan file** — `kitchen/design-decisions.json` | A file, a shape, a reader, a writer. And the risk the audit has already priced: a file with no writer is `needs_regrooming`. | The "one journal, one event log" rule, if the file's entries are not also events. |
| **B. A new field on the recipe** — `DesignDecisions []Decision` beside `Audit.Deviations` | The smallest change, and it is **exactly the shape of v1's dead field**: a place the domain claims to keep decisions and no command writes. | Nothing — if the writer exists. The cost is entirely in the writer. |
| **C. Nothing new** — a design decision is a **comment in the recipe's thread**, authored by the sous chef, read by whoever works on it | Nothing in the model. The thread already exists, already renders, and already carries a finding's discussion. | A design decision that outlives its recipe, or one that spans recipes. |

### The recommendation

**Option C — the plan already has this, and it does not need another file.** A design
decision is a choice about *how to build*, and the place a builder looks for it is the
thread of the recipe it is about: `maitre-cli recipe comment RCP01 --body "…"`. The journal
records it, the detail pane shows it, and the audit trail is complete. The alternatives
each add a place where a decision can be written down and not found — which is the
"field with no writer" shape v1 died of ✓.

**And the field that already exists should not be left as it is.** `Audit.Decisions` has
no writer ✓. Either it gets one — the sous chef records the choice, in the recipe, as
the third kind of record the thread already shows ✓ — or it goes, the way the plan
closed the other writer-less fields in the field-fates round (`Execution.StartedAt`,
`Execution.Branch`, `Evidence.Branch`, `Step.Evidence`, `Finding.FixedIn` — all deleted,
and `Evidence.CompletedAt` and `Step.Deviation` given writers ✓). A plan that keeps a
named field for decisions and never writes one is the v1 lesson, reproduced in the room
where the lesson was learned.

**What would have to be true to build it.** A design decision that must outlive the
recipe it was made for — a choice about the *project*, taken before the recipe existed
and binding several recipes — and that an agent must find before it starts. That is a
project-level fact, and the project already has a place for project-level facts: the
**charter** in `project.json`, and the `decisions.json` the plan format already names.

**If that happens, the shape is the one the plan format already promises**: decisions in
`decisions.json`, each also written as a journal `decision` event, so the journal stays
the only event log. Do not put them in a file no journal event mirrors.

**What it would change**

| | |
|---|---|
| The plan format | **Nothing today.** If a project-level decision arrives, `decisions.json` stops being a name in a format and becomes a file with a writer — and every entry is mirrored as a journal event. |
| The lifecycle | **Nothing** — a decision is not a state and is not promoted. |
| The console | **Nothing** — a design decision in the thread already shows under its recipe. A project-level decision would show on the Settings tab, beside the charter. |
| The CLI | `maitre-cli recipe comment` — **it already exists**. |

**What it forecloses.** A design decision cannot be found without knowing which recipe it
belongs to, and it cannot be listed across the plan. For a plan of one recipe that costs
nothing; for a plan of fifty it is the first thing to change.

---

## 3. Brownfield / delta

> *A brownfield job starts somewhere else — there is code, and the question is what to
> change.* — the question, as the user put it

### The question

**What is a delta against a brownfield codebase, in this domain's vocabulary?** And: what
does the recipe's `intent.outcome` mean when the code already exists?

### The evidence

| What | Where |
|---|---|
| The verifier's rule: an evidence is a command **that matched something**. Evidence that matched nothing is rejected — "a green gate that proves nothing" is the project's first v1 lesson. | `internal/core/verifier.go:170`, [../AGENTS.md](../AGENTS.md) |
| `ac.falsifiable` accepts only a command the tool can count matches from, or a `manual:` procedure. A brownfield criterion is verifiable in a way a greenfield one is not: **the suite already exists.** | `internal/core/lifecycle.go:397` |
| The domain's own shape: an intent, a problem, a scope, criteria, and a method. A brownfield job has all of these — the *problem* is "this behaviour is wrong or missing" and the *outcome* is "this behaviour is right". | [docs/03 §2](03-domain-model.md) |
| A finding already carries an anchor to code: `{file, line, side}`. | `internal/core/model.go` |

### The asymmetry, stated plainly

A greenfield criterion is proved by a **new** test: the test is written, and the test's
passing is the evidence. A brownfield criterion is proved by an **existing** suite that
must not change its verdict — which is the same evidence, run in the opposite direction.
The domain has one word for both — an **evidence** — and the plan's machinery does not
care which direction it runs. That is the whole proposal: **the plan needs no brownfield concept;
the recipe's method needs one extra step.**

### The options

| Option | What it costs | What it forecloses |
|---|---|---|
| **A. A `brownfield` flag on the recipe** | A field, and a rule the domain must enforce: what does the flag *change*? Nothing, unless something reads it. | Nothing — but it is a flag that changes no behaviour, which is the "safe unknown value" shape. |
| **B. A delta artifact** — a recipe whose criteria are regression criteria, anchored to code | A vocabulary, and a second kind of criterion the gates must know about. | The uniformity that makes the gate table one table. |
| **C. A reading as the first step** — the brownfield recipe is an ordinary recipe whose **method starts by reading**, and whose criteria are regression criteria proved by the existing suite | One extra step in the method. Nothing in the model, the lifecycle or the plan format. | A brownfield job is expressed in the same words as a greenfield one — there is no `kind: brownfield` to check. |

### The recommendation

**Option C.** A brownfield recipe is a greenfield recipe whose first step is a **reading**:
"read the loader, list what it does" — a step like any other, with a `done_when` like any
other (*"the loader's behaviour is written down"*). Its criteria are **regression
criteria**: *"the queue keeps its order"*, *"a recipe with no criteria is refused"*, each
proved by the command that already exists — `go test ./internal/core -run Board` — which
the verifier counts exactly as it counts a new test. Nothing about the recipe changes
shape; what changes is the *direction* the evidence runs.

**The one thing that must be said, because it is the trap:** a brownfield criterion is
proved by a test that **already passes**. So the evidence is not the criterion's own
existence — it is the criterion's **failing on the broken code**. The method therefore
carries the order: write the test, watch it fail, change the code, watch it pass. A
brownfield recipe whose criteria were green before the work started has proved nothing —
which is v1's lesson (*evidence clobbering*) and the audit's (*zero-match verification*)
in one line.

**What would have to be true to do it differently.** If a brownfield job's evidence has to
be *about* the pre-existing code — "this behaviour is preserved" — then a passing test is
not enough, because the test may not have existed before the change. That needs a
**baseline**: the suite run before and after, compared. That is a real feature, it is not
what the domain does today, and it should be built when a brownfield job's evidence actually
needs it — not before.

**What it would change**

| | |
|---|---|
| The plan format | **Nothing** — a brownfield recipe is a recipe. |
| The lifecycle | **Nothing** — the steps are steps; the criteria are criteria. |
| The console | **Nothing** — the detail pane already shows the method and the criteria. |
| The CLI | **Nothing** — `maitre-cli step add … --done-when "…"` already carries the reading. |

**What it forecloses.** The plan cannot tell a brownfield recipe from a greenfield one by
reading it. A reader has to infer the direction from the criteria's wording. If that
matters, the honest place is **a note in the recipe's own thread** — `maitre-cli recipe comment` — a record that already has
a writer ✓ — and not a new field. (The recipe's `Questions` field looks like the
natural home and is not: it has **no writer** and no CLI path ✓ — the shape the
field-fates round closed elsewhere ✓ — so it is not a home for anything until something
writes it.)

---

## What this document does not decide

These are proposals. Each table's recommendation is the reviewer's reading of the evidence,
and the chef may take the other option, or a fourth one, without the domain breaking —
which is the test of whether a decision belongs here at all.
