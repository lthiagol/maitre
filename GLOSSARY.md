# maitre

The plan a team cooks from: a set of recipes, each one a durable brief an agent
consumes and executes. The kitchen is the shared artifact; the brigade is the people and
agents working it.

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
It becomes a recipe when it is promoted, and shaping an idea into a small task is an
update, not a promotion: the id does not change.
_Avoid_: Special, prep note, note, brief, ticket

**Course**:
One intent across several recipes — a themed set of dishes served together.
_Avoid_: Epic, theme, group

**Menu**:
Coordinated recipes with ordering and cross-dependencies — the whole plan as one
composition.
_Avoid_: Program, roadmap, plan

### The brigade

**Chef**:
The human. Owns the kitchen: decides what to cook next, signs off, and tastes each
dish before it goes out.
_Avoid_: User, approver, owner, manager

**Maître d'**:
The human-facing console. Greets the chef, takes the order, guides them, and checks on
tables.
_Avoid_: TUI, dashboard, frontend, UI

**Sous chef**:
The coordinator agent. Composes recipe cards, preps the line, and runs the tasting.
_Avoid_: Coordinator, planner, scribe

**Line cook**:
The runner agent. Executes the ticket on the line, following the method or deviating
from it with a note.
_Avoid_: Runner, executor, implementer

**Taster**:
The reviewer agent. Tastes the dish independently — ideally a different model, always a
different session — and files tasting notes.
_Avoid_: Reviewer, critic, judge

**Expo**:
The orchestrator agent. Calls the tickets, routes them, and keeps the line moving.
_Avoid_: Orchestrator, dispatcher, supervisor

### The line

**Sign-off**:
The chef's approval of a recipe. Recorded with an approver and a fingerprint of the
sections that are then frozen.
_Avoid_: Approval, gate, signoff

**Tasting**:
The external review of a finished recipe. Clean tasting advances the recipe; open
tasting notes send it back for remediation.
_Avoid_: Review, code review

**Evidence**:
A verification result a machine produced — the command, its exit code, what it matched,
and when. An evidence that matched nothing is not evidence.
_Avoid_: Proof, test result, verification

**Completion**:
The recipe-level stamp: the moment the recipe was done. One field, written on the
`finish` transition.
_Avoid_: Evidence, closed, done

**Lane**:
Where a recipe sits in the chef's workflow: grooming, sign-off, on the line, tasting, on
hold, served. Derived from the recipe's state, never stored.
_Avoid_: Status, column, phase, bucket

**Method**:
The proposed path for executing a recipe: work packages and steps, revisable, with
deviations recorded rather than hidden.
_Avoid_: Steps, plan, guidance, instructions

**Open question**:
A question raised while shaping a recipe that must be resolved before it can be locked.
_Avoid_: Blocker, issue, todo

---

## Two vocabularies

CPD — the methodology — keeps a generic vocabulary. maitre — the tool — uses
the language above. Every term above names its CPD counterpart — or says that it has
none — so the two never drift.

The mapping lives here and nowhere else. It is the single source for the CLI reference,
the skills and the docs.

| CPD | maitre |
|---|---|
| crafted prompt | **recipe** |
| idea / note / brief | **backlog item** |
| decomposed prompt | **course** |
| prompt family | **menu** |
| coordinator | **sous chef** |
| runner | **line cook** |
| reviewer | **taster** |
| orchestrator | **expo** |
| the human validator | **chef** |
| the human-facing console | **maître d'** |
| evidence | **evidence** |
| — | **completion** |
| steps / guidance | **method** |

**The dash is a term CPD does not have.** CPD's *evidence* is the record a criterion
carries; maitre's **completion** is the recipe-level stamp that says the recipe was
done.

Everything else stays generic in both: acceptance criteria, scope, findings, work
packages, lifecycle states, lanes' underlying states.
