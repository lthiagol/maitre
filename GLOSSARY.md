# maitre

The plan a team cooks from: a set of recipes, each one a durable brief an agent
consumes and executes. The kitchen is the shared artifact; the brigade is the people and
agents working it.

## Language

### The plan

**Kitchen**:
The directory holding the plan — recipes, courses, prep notes, specials, and the event
journal. The one place the plan lives.
_Avoid_: Workspace, plan dir, `.mp`

**Recipe**:
The durable unit of work — a structured, reviewed, agent-friendly brief carrying intent,
scope, acceptance criteria, and the audit trail from its execution. CPD calls this a
*crafted prompt*.
_Avoid_: Milestone, task, ticket, prompt

**Special**:
A hunch, an observation, or an agent discovery. The rung below a prep note: something
worth noticing, not yet a problem worth solving.
_Avoid_: Idea, note, ticket

**Prep note**:
A lightweight intent and scope, before acceptance criteria exist. Promoted to a recipe
when success needs to be verifiable.
_Avoid_: Brief, sketch, draft

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

**Proof**:
A verification result a machine produced — the command, its exit code, what it matched,
and when. A proof that matched nothing is not a proof.
_Avoid_: Evidence, test result, verification

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
the language above. Every term above names its CPD counterpart so the two never drift.

The mapping lives here and nowhere else. It is the single source for the CLI reference,
the skills and the docs.

| CPD | maitre |
|---|---|
| crafted prompt | **recipe** |
| idea / note | **special** |
| brief | **prep note** |
| decomposed prompt | **course** |
| prompt family | **menu** |
| coordinator | **sous chef** |
| runner | **line cook** |
| reviewer | **taster** |
| orchestrator | **expo** |
| the human validator | **chef** |
| the human-facing console | **maître d'** |
| evidence | **proof** |
| steps / guidance | **method** |

Everything else stays generic in both: acceptance criteria, scope, findings, work
packages, lifecycle states, lanes' underlying states.
