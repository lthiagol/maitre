# maitre

The plan a team cooks from: a set of recipes, each one a durable brief an agent
consumes and executes. The kitchen is the shared artifact; the crew is the people and
agents working it.

## The rule

**The kitchen names the artifact and the place. Plain words name the machine's records
and the agents' jobs.**

| Layer | Vocabulary | Examples |
|---|---|---|
| The artifact | kitchen | recipe · the kitchen · the line · the board · service · tasting · sign-off · method |
| The records | plain | criterion · evidence · completion · lifecycle states · `depends_on` |
| The jobs | plain | chef · sous-chef · runner · reviewer · orchestrator |
| The plan's loose end | plain | backlog item · idea · feature · bug |

Two consequences that used to be accidents and are now rules:

- The console says *tasting note*; the code says `Finding`. That is a **translation**,
  not a mismatch — the console is the surface where the kitchen is free to be charming.
- A recipe's row shows the **plain state word** (`closed`, `executing`) and the lane
  headers are **kitchen sections** (`Served`, `On the line`). The console never shows a
  lane word where a state word belongs.

## The names

| Thing | Name |
|---|---|
| The project | **maitre** |
| The console | `maitre` |
| The agent CLI | `maitre-cli` |
| The plan directory | `kitchen/` |
| A recipe's id | `RCP01` |
| A backlog item's id | `BL01` |

The decision of record is [docs/07-vocabulary.md](docs/07-vocabulary.md); the choice is
recorded in [ADR 0009](docs/decisions/0009-vocabulary.md).

## Language

### The plan

**Kitchen**:
The directory holding the plan — recipes, courses, backlog items, and the event
journal. The one place the plan lives.
_Avoid_: Workspace, plan dir, `.mp`

**Recipe**:
The durable unit of work — a structured, reviewed, agent-friendly brief carrying intent,
scope, acceptance criteria, and the audit trail from its execution. CPD calls this a
*crafted prompt*.
_Avoid_: Milestone, task, ticket, prompt

**Backlog item**:
The plan's loose end — one artifact with a **kind**: an `idea` (a hunch, an
observation, an agent discovery), a `feature` or a `bug` (a small task, ready to work).
It replaces the two file types the plan used to hold — see
[ADR 0009](docs/decisions/0009-vocabulary.md) § 3 — and the **kind** is the only truth
about its shape. It becomes a recipe when it is promoted, and shaping an idea into a
small task is an update, not a promotion: the id does not change.
_Avoid_: Special, prep note, note, brief, ticket

**Course**:
One intent across several recipes — a themed set of dishes served together.
_Avoid_: Epic, theme, group

**Menu**:
Coordinated recipes with ordering and cross-dependencies — the whole plan as one
composition.
_Avoid_: Program, roadmap, plan

**Line**:
Where the work happens — the recipes being groomed, signed off, executed and tasted.
The console's lanes are sections of the line, never names for a recipe's state.
_Avoid_: Stage, phase, queue, lane

**Board**:
The plan's shape — what waits on what, drawn as a dependency tree. Derived from the
plan, never stored.
_Avoid_: Tree, graph, dependency view, kanban

**Service**:
The flow of the kitchen — the ordered stages a plan moves through, from shaping to the
pass. It names the flow, not a stage.
_Avoid_: Stage, phase, workflow, pipeline

### The crew

**Chef**:
The human. Owns the kitchen: decides what to cook next, signs off, and tastes each
dish before it goes out.
_Avoid_: User, approver, owner, manager

**Sous-chef**:
The agent that composes recipe cards, preps the line, and runs the tasting.
_Avoid_: Coordinator, planner, scribe

**Runner**:
The runner agent. Executes the recipe on the line, following the method or deviating
from it with a note.
_Avoid_: Line cook, executor, implementer

**Reviewer**:
The reviewer agent. Reviews the recipe independently — ideally a different model,
always a different session — and files tasting notes.
_Avoid_: Taster, critic, judge

**Orchestrator**:
The orchestrator agent. Calls the recipes, routes them, and keeps the line moving.
_Avoid_: Expo, dispatcher, supervisor

### The work

**Sign-off**:
The chef's approval of a recipe. Recorded with an approver and a fingerprint of the
sections that are then frozen.
_Avoid_: Approval, gate, signoff

**Tasting**:
The external review of a finished recipe. A clean tasting advances the recipe; open
tasting notes send it back for remediation.
_Avoid_: Review, code review

**Evidence**:
One verification result — a machine's, or the chef's: the command, its exit code, what
it matched, and when. A result that matched nothing does not count; the chef's sign-off
is the one exception, and it claims no match.
_Avoid_: Proof, test result, verification

**Completion**:
The recipe-level stamp: the moment the recipe was done. One field, written on the
`finish` transition.
_Avoid_: Evidence, closed, done

**Lane**:
Where a recipe sits in the chef's workflow: grooming, sign-off, on the line, tasting,
on hold, served. Derived from the recipe's state, never stored.
_Avoid_: Status, column, phase, bucket

**Method**:
The proposed path for executing a recipe: work packages and steps, revisable, with
deviations recorded rather than hidden.
_Avoid_: Steps, plan, guidance, instructions

**Open question**:
A question raised while shaping a recipe that must be resolved before it can be locked.
_Avoid_: Blocker, issue, todo

### The surfaces

The console is a **surface**, not a job: [ADR-0003](docs/decisions/0003-surface-split.md)
says it manages the plan and never authors the artifact, and the crew above is the
people and the agents. What the console carries is a **string** — the window title,
the tab bar and the CLI's help hold the name, and nothing more.

**maitre**:
The console — the chef's window on the kitchen. It greets the chef, takes the order,
and checks on the tables.
_Avoid_: maître d', TUI, dashboard, frontend, UI

---

## Two vocabularies

CPD — the methodology — keeps a generic vocabulary. maitre — the tool — uses
the language above. Every term above names its CPD counterpart — or says that it has
none — so the two never drift.

The mapping lives here and nowhere else. It is the single source for the CLI reference,
the skills and the docs.

| CPD | maitre |
|---|---|
| idea | **backlog item** — kind `idea` |
| note | **backlog item** — kind `idea` |
| brief | **backlog item** — kind `feature` or `bug` |
| crafted prompt | **recipe** |
| decomposed prompt | **course** |
| prompt family | **menu** |
| coordinator | **sous-chef** |
| runner | **runner** |
| reviewer | **reviewer** |
| orchestrator | **orchestrator** |
| the human validator | **chef** |
| the human-facing console | **maitre** |
| evidence | **evidence** |
| steps | **method** |
| guidance | **method** |
| — | **kitchen** |
| — | **sign-off** |
| — | **tasting** |
| — | **completion** |
| — | **line** |
| — | **board** |
| — | **service** |
| — | **lane** |
| — | **open question** |

**The dash is a term CPD does not have.** maitre's **completion** is the recipe-level
stamp that says the recipe was done, and the per-criterion record CPD calls *evidence*
is maitre's **evidence** — the same word for the same thing.

**The backlog item is where CPD's ladder and ours are reconciled.** CPD splits rung 0
(*idea / note*) from rung 1 (*brief*); maitre holds **one artifact** for both and lets
the **kind** carry the difference — `idea` is rung 0, `feature` and `bug` are rung 1.
From there the ladder is unchanged: rung 2 is a **recipe**, rung 3 a **course**, rung 4
a **menu**.

Everything else stays generic in both: acceptance criteria, scope, findings, work
packages, lifecycle states, lanes' underlying states.
