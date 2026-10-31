# The TUI plan — the console, second time round

> **This is the design.** The *state* — where the work stopped — is
> [README.md](README.md). The *contract* for one session is a brief in
> `work/briefs/`, and the brief wins over this file where they differ.
>
> Scope: the console only. The CLI comes after. The domain work the console needs is
> listed separately, because some of it is not optional.

## The answer, in one table

> The *when the round started* column is a **snapshot**: it says what the code said
> before the round began, so the rows whose slices have landed are the record of the
> problem they solved. The round's **state** is
> [work/tui/README.md](README.md)'s table, and the code is the truth.

| # | Your point | What the code said when the round started | The proposal |
|--:|---|---|---|
| 1 | **Title** | the project name **is** already in the header — `⋈ maitre · kbtest` — but dimmed, and the **window title is the hardcoded string `maitre`** (`views.go:53`) | the window title carries the project; the header's project name stops being a footnote; a `--project` flag for two checkouts of one project |
| 2 | **Home tab** | v1's **Overview** was dropped because it *restated the board* | a **Home** tab with **numbers, never items**; the art moves there; the Sous-Chef tab dissolves into it; the tab count stays 7 |
| 3 | **Mouse** | `internal/ui` handles **no mouse event at all** — no `tea.MouseMsg` case exists | bubbletea v2 sets the mouse mode on the `View`; a **hit map the render builds**; wheel, click a row, click a tab |
| 4 | **Backlog** | the cursor **does** move (verified) — but `enter` **does nothing**: a backlog item has **no detail pane**. And `space` promotes **an idea**, which yields a recipe with **no acceptance criteria** | a detail pane for the item; `space` stops being a trap on an idea; **shape it** · dismiss · create reach the console — the operations already exist. The **free edit** of every field is **not in this round**: it is a five-field surface, and the round's need is the moves |
| 5 | **Path & Board** | **Board** is the glossary's own word for the dependency tree (*"what waits on what, drawn as a dependency tree"*). The Path's `READY` means *not blocked* — it holds drafts and things waiting on you | **Path → Ready**, and Ready means *adoptable now*. **Board stays Board.** Merging them is a no |
| 6 | **Grouping** | the glossary's word is **course** (*avoid: epic, theme, group*), and [docs/06 § 1](../../docs/06-proposals.md#1-courses) already ruled **no course file**. Recipes carry **no tags**; backlog items do | a **filter** now (no format change), and **tags** as a later domain decision, taken with the CLI |

## 1 — The project identity

**The gap is not the header. It is the window title.**

| Where | Today | Proposal |
|---|---|---|
| The header, row 1 | `⋈ maitre · kbtest` — the name is there, in `styFaint` | `⋈ maitre · kbtest`, with the name in an accent style. A name in the faint style reads as a footnote, and a footnote is not an identity |
| The **window title** (the OS tab) | `"maitre"` — hardcoded | `"maitre · kbtest"`. This is the row that answers *"which of my four terminals is this?"* |
| **The plan's path** | nowhere | the Settings tab shows the plan's root — the one place that answers "which plan is this" for two plans whose names are alike |
| The flag | — | `maitre --project <name>` overrides the displayed name for the session. A chef with two checkouts of one project needs it |

**The name's source** is `project.json`'s `name`, which `Init` fills with the
basename of the directory **above** `kitchen/`. Two plans can share a basename, and
then the console cannot tell them apart — which is what the Settings tab's path and
the `--project` flag are for.

**A per-project accent colour** — a hue derived from the project's name, so two
terminals are distinguishable at a glance — is the tempting extra. It is cheap
(`fnv` over the name) and it is *identity*, not preference, so it does not fight the
theme. I would hold it until the window title is in, because the window title is the
part that actually solves the problem.

## 2 — The Home tab

**The precedent that must be answered first.** v1 had an **Overview** lane and v2
dropped it: *"the triage-ordered board **is** the overview; a second tab would
restate it"* ([docs/05 § Tabs](../../docs/05-human-interface.md#tabs)). The chef's
tab is not that tab — but the risk is the same one, and the difference has to be
written down.

> **The invariant: the Home tab shows numbers; it shows no items.**
>
> No row on Home is selectable, `space` does nothing there, and the tab's item count
> is zero. A test holds it, the way slice 9's invariant held *"a recipe is filed
> under **Needs you** if and only if the console has a chef move for it"*.

**What it holds** — every row a number, and nothing that duplicates another tab's job:

| Block | The number | Where it comes from |
|---|---|---|
| The plan | the project's name and the path | `project.json` |
| The recipes | need you · on the line · being shaped · served — the triage's sections, and **the parts add up** | the triage sections |
| The **evidence** | proven · failed · rejected · unproven · pending — the states, and **the parts add up** | `core.EvidenceSummaryOf`, summed |
| The **tasting** | open notes · open findings | `core.OpenFindingsOf` |
| The backlog | ideas · small tasks | `core.BacklogGroups` |
| The **automation** | the plan's domain version · this build's · the last journal event | the Sous-Chef tab's content, moved |

**The art moves here.** It leaves the Sous-Chef tab and goes to Home, where it can
sit beside the numbers instead of above them. Slice G's rule stands unchanged: **the
decoration yields to the content, in both dimensions** — the mark is drawn only when
it fits, in width *and* in height, and at 80×24 the tab is what it would be without it.

**The Sous-Chef tab dissolves.** Its two facts (the version agreement, the last event)
are on Home. The tab comes back when the orchestrator does — it was a placeholder
that said so plainly, and a placeholder that has a home is not a placeholder.

**The name.** Candidates, with the glossary's constraint applied (*avoid: dashboard,
TUI, frontend, UI*):

| Name | For | Against |
|---|---|---|
| **Home** | plain, unambiguous, no glossary conflict | a web-ism |
| **Pulse** | kitchen-flavoured, short | a new word the glossary does not carry |
| Kitchen | it *is* the kitchen | it names the plan directory, not a summary of it |
| Overview | v1's dropped lane; it restated the board | the precedent above |
| Status | plain | collides with the lifecycle state word |

**My recommendation: Home.** The glossary's rule is that the kitchen names the
artifact and the place, and plain words name everything else — a tab label is a plain
word, and *Home* is the plain word for the landing tab. *Pulse* is the good
alternative if you want the flavour; nothing forbids it.

## 3 — The mouse

**The facts.** `internal/ui` has no `tea.MouseMsg` case anywhere — the console
receives mouse events and drops them. In bubbletea v2 the mouse mode is **per view**,
not a program option (that was v1):

```go
v := tea.NewView(m.renderLayout())
v.MouseMode = tea.MouseModeCellMotion   // ← the whole of the enabling
```

`tea.Mouse{X, Y, Button, Mod}` carries the coordinates, with `MouseWheelUp` /
`MouseWheelDown` as buttons.

**The hard part: the console must know where it drew things.** The render is a
string; the mouse sends coordinates. So the render has to leave a **hit map** behind:

| What | The box |
|---|---|
| Each tab label | its column range, in the tab bar's row |
| Each list row | its row range, and the item index it holds |
| The detail pane | its region |
| The filter line (T7) | its region |

**The map is built by the render, not beside it.** A second function that guesses
where the tab bar is would be a second answer to one question — the house's own
standing failure. The renderer emits the boxes as it draws.

**What the mouse does:**

| Gesture | Where | What |
|---|---|---|
| wheel | over the list | scroll the list, keeping the selection visible |
| wheel | over the detail pane | scroll the detail |
| click | on a tab label | switch to that tab |
| click | on a row | select it — and open the detail, if the row is something with a detail |
| click | on the detail pane | focus the detail, so the wheel scrolls it |

**The test story matters here, because CI has no terminal.** The console is driven
from a script by `ui.Drive`; the mouse has to be drivable the same way. `ui.Drive`
grows a mouse step, and a test covers a click on a tab and a wheel over the list.

**One caveat worth stating:** the mouse makes the console *look* clickable. It must
not *promise* a click that does nothing — a clickable-looking thing that ignores the
click is worse than a thing that never looked clickable.

## 4 — The Backlog tab

**What I verified, on a seeded plan:**

| The chef's words | What is true |
|---|---|
| "we can't browse it, with mouse" | true — there is no mouse at all |
| "or keyboard" | **half true**: `↓` moves the cursor and the row is highlighted (ANSI confirmed). What is missing is anywhere to *go* |

**The real gap: `enter` does nothing on a backlog item.** `Model.currentRecipe()`
(`model.go:409`) has no `TabBacklog` case, so the detail pane is empty — the item's
body, its scope, its tags and its status are unreachable from the console. You can
move the cursor over a thing you cannot read.

**And `space` is a trap.** On a backlog row, `space` is *promote to recipe* — one
keystroke, no confirmation. And the promotion is not gated by kind:

```
$ maitre-cli backlog promote BL01        # BL01 is an idea
"acceptance": [],  "lifecycle": "draft"
```

An idea promotes to a recipe with **no acceptance criteria** — and a recipe with no
criteria can never pass `ac.falsifiable`, so it can never be signed off. The console
offers, in one keystroke, a move that leads to a dead end. That is the "not optimal"
feeling, and it is a defect, not a taste.

**What the tab should have:**

| Move | When the round started | Proposal |
|---|---|---|
| **Read** | nothing | a detail pane: title, body, scope, kind, status, tags, dates, `promoted_to` |
| **Promote** | `space`, on every row, including ideas | `space` on a **small task** promotes. On an **idea**, `space` offers *shape it* — the kind changes to a small task (the chef picks which) and a scope is added, which is an update, not a promotion. The footer says which |
| **Dismiss** | not in the console | the console dismisses with a reason — `DismissBacklogItem` exists in core and `maitre-cli backlog dismiss` exists; only the console's surface is missing |
| **Create** | not in the console | the console creates a backlog item — `CreateBacklogItem` exists. A backlog item is **not** a recipe, so this does not break ADR-0003 (the console never authors the recipe) |
| **Shape it** | not in the console | the console changes the kind and adds a scope — the **update** that `docs/03` § 1.1 calls *shaping*. The id does not change. `UpdateBacklogItem` exists |
| **Edit the fields** | not in the console | editing title, body and tags as a form is **not in this round**: it is a five-field surface, and the round's need is the moves below |

Every one of those operations already exists in the domain. **This slice is the
console catching up to its own domain, and it needs no CLI work.**

## 5 — Path and Board

**Board keeps its name, because the glossary already gives it one:**

> **Board**: The plan's shape — what waits on what, drawn as a dependency tree.
> Derived from the plan, never stored. — *Avoid*: Tree, graph, dependency view, kanban

So "Tree" is a retired word for it, and "Board" is the decision of record. The Board
tab answers the chef's second question already: *what depends on what, and in what
order*.

**The Path tab answers the wrong question.** Its `READY` list means *not blocked* —
and it holds drafts (*"later"*), recipes waiting for your sign-off (*"needs you"*),
recipes being executed (*"running"*), and recipes you must taste. That is not *"the
ones really ready, no dependency, a menu for what can be adopted at that moment"*.

**The proposal: the Path tab becomes the Ready tab**, and it answers one question:

> **A recipe is in the Ready list if and only if an agent could claim it right now:**
> signed off, no unmet dependency, not blocked, not cancelled, not already claimed.

| Block | The row | The chef's move |
|---|---|---|
| **ADOPTABLE** | recipes an agent could claim now | none — the agents have them. The console does not claim: that is the runner's move (ADR-0003) |
| **SMALL TASKS** | open backlog items of kind `feature`/`bug` | *promote* — and the row says what it takes after that: promote → sign off → run. Two steps, not one |

**This needs one new derived view** — `core.Adoptable`. And it exposes a standing
inconsistency worth fixing in the same breath:

| Today | Problem |
|---|---|
| `core.ReadinessOf(r)` answers *"can this be claimed?"* — and **ignores `depends_on`** | a signed-off recipe that waits on an unfinished one reads as **Ready** |
| `core.Path` answers *"is this blocked?"* — and knows the dependencies | two answers to one question, and they disagree |

`Adoptable` is the one that considers both. `ReadinessOf` either folds into it or
says plainly that it is about the *state* and the state is not the whole answer.

**Merging Path and Board is a no.** They answer different questions — *"what can I
take now?"* and *"in what order, and what is stuck behind what?"* — and the Board's
trunk is already the Path's ready list, so merging would draw the same list twice in
one tab. After the change the overlap **shrinks**: the Board shows everything not
blocked (the shape), the Ready tab shows only what can be adopted (the moment). One
derivation, two questions, two sets.

**The Board gets visual.** The chef liked the tree from the rust version; today's is a
vertical one with `├─ └─` connectors. The proposal:

| Width | The board |
|---|---|
| ≥ 160 cols | the tree drawn **left to right** — the trunk on the left, forks branching right, the recipes as nodes with their titles |
| < 160 cols | today's vertical tree, byte-for-byte |

**A golden text cannot tell you the tree looks good.** A screenshot can — and
`scripts/console-shot.py` already exists to capture one. The visual slices (T4, T6)
need a **screenshot review**, not only the golden.

## 6 — Grouping: what the plan already ruled, and the gap

**The vocabulary is settled.** The chef's "epic" is the glossary's **course**
(*"one intent across several recipes"*, avoid: epic, theme, group) — and
[docs/06 § 1](../../docs/06-proposals.md#1-courses) already ruled against a course
file: *a course is **the relation**; a grouping is a query over the plan, not a thing
in it.*

**But the chef's need is not a course.** *"Work planned for the TUI and some for the
CLI; work on one or the other at a given moment"* is a **label** axis, and the label
axis does not exist on recipes — `BacklogItem` has `Tags []string`, `Recipe` has
nothing.

| Option | What it costs | Verdict |
|---|---|---|
| **A. `tags` on the recipe** — `"tags": ["tui","cli"]` | a new field on `Recipe` → the plan format moves → `DomainVersion` 7 → 8 → the repo's own plan is stamped forward (`scripts/stamp-plan/`). Per the no-migration ruling, that is cheap. The CLI needs a flag later | **the right fact — later, with the CLI** |
| **B. a filter over what the recipe already says** — id, title, intent, scope, and the backlog's tags | nothing in the model | **the TUI-only answer — now** |
| **C. a course as a view over `depends_on`** | nothing | already ruled out, and it does not answer "TUI or CLI" |

**The plan takes B now and leaves A as a decision.** The filter is a console feature
and costs no format change:

| | |
|---|---|
| `/` | opens a filter line under the summary |
| the match | case-insensitive, over the id, the title, the intent and the scope — and the backlog items' tags |
| `esc` | clears it |
| the summary line | says when a filter is on, and on what — a filter the chef cannot see is a filter that loses items |
| the scripts | drivable, like every other key |

**When tags arrive, the filter covers them** — and the tag becomes a *label*, never a
second place where the plan's shape lives. A tag is not an artifact: it cannot be
signed off, it has no acceptance criteria of its own, and it is not shared between
plans. That is the line between a label and a course.

## The domain work the console needs

Small, and none of it needs a CLI verb — the operations exist.

| # | What | Why the console cannot be honest without it |
|--:|---|---|
| D-1 | `core.Adoptable` — a derived view: signed off, no unmet dependency, not blocked, not cancelled | the Ready tab is the chef's question, and it is a question about the plan, so the answer belongs to the domain. Registering it in `core.DerivedViews` is what keeps `docs/03` § 8 honest — the wire checks both directions |
| D-2 | `core.PlanStats` — the numbers Home shows: evidence, tasting, backlog | the same reason. A number the console computes itself is a number the CLI would compute differently |
| D-3 | **an idea cannot promote to a recipe with no criteria** | either the promotion is refused with the reason, or it gives the recipe a criterion, the way a small task's promotion does |
| D-4 | `ReadinessOf` and `Adoptable` cannot disagree | two answers to one question is the house's standing failure |
| D-5 | the promotion writes the same shape as every other writer (`scope: []`, not `null`) | verified: `recipe create` writes `[]`, `backlog promote` writes `null`. Same format, two answers |
| D-6 | *(later)* `tags` on the recipe — `DomainVersion` 7 → 8 | a format decision, taken with the CLI, not with the console |

## The sessions

**One session, one slice.** The slices below are the sessions; the feature names are
T1–T7, and two of them split because they do not fit in one sitting.

| Session | Feature | What it lands | Needs | Size |
|--:|---|---|---|---|
| **T1** | the identity | the window title carries the project; the header's name stops being a footnote; `--project`; the Settings tab shows the plan's path | — | S |
| **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls the list and the detail | — | M |
| **T2b** | the mouse, part 2 | a click on a tab switches; a click on a row selects and opens the detail; the mouse is drivable from a script | T2a | S |
| **T3a** | the backlog is readable | the detail pane renders a backlog item; `space` on an idea stops promoting it into a dead end; D-3 and D-5 land | — | M |
| **T3b** | the backlog moves | shape it · dismiss · create reach the console; D-4's first half | T3a | M |
| **T4** | the Home tab | the landing tab with numbers and no items; the art moves here; the Sous-Chef dissolves into it; D-2 lands | — | M |
| **T5** | the Ready tab | `core.Adoptable` and the tab that answers the chef's question; the Path becomes Ready; D-1 and D-4 land | — | M |
| **T6** | the Board, as a visual place | the tree drawn left to right at ≥ 160 columns; clicking a node opens its detail | T2b | M |
| **T7** | the filter | `/` filters the list of the current tab | — | S |

**The order is not forced.** T1, T3, T5 and T7 are independent of each other; T2b
needs T2a, T6 needs T2b, and T3b needs T3a. The sequence above is the one that
shows the console's shape earliest.

## The test story

| | |
|---|---|
| **Fixtures** | `core.Seed` and `core.SeedStress` build the plans **through the operations**, so a fixture cannot drift from the model. Every slice's data comes from the seeder — no hand-written JSON |
| **Interaction** | `ui.Drive` drives the console with keys, and grows a mouse step. `maitre --keys "tab,down,space"` stays the script |
| **Goldens** | `make goldens` regenerates, `make goldens-check` holds them; the change is claimed by a test |
| **Screenshots** | `scripts/console-shot.py` for the slices where the *look* is the feature — T4 and T6. A text golden cannot see a tree that looks wrong |
| **The domain** | `make test PKG=./internal/core RUN=…` for the views; the derived-views wire checks `docs/03` § 8 against `core.DerivedViews` in both directions |

## Not in this plan

| | |
|---|---|
| **The CLI** | the chef's order: the console first, the CLI after. Every operation this plan needs already exists in `core` — the CLI only ever gains verbs over them |
| **`tags` on the recipe** | a domain decision, taken with the CLI: `DomainVersion` 7 → 8, and the plan format moves |
| **A course** | ruled out by [docs/06 § 1](../../docs/06-proposals.md#1-courses), and not needed for the filter to be useful |
| **The orchestrator** | the Sous-Chef tab comes back when the orchestrator does |
| **The doc tree** | this file lives in `work/`. When the chef rules, the durable parts move to `docs/05-human-interface.md` (the console), `docs/03-domain-model.md` (the new views) and `docs/06-proposals.md` (the grouping) |

## The questions the plan is waiting on

| # | Question | My recommendation |
|--:|---|---|
| Q1 | the Home tab's name | **Home** (Pulse if you want the flavour) |
| Q2 | the Home tab's numbers — the set above, or a smaller set? | the set above; the evidence and tasting rows are the ones the board does not already carry |
| Q3 | does the Ready tab show the small tasks, or only the recipes? | both — but the small tasks' row says *promote → sign off → run*, so nobody reads it as one keystroke |
| Q4 | the filter — before or after the tags? | the filter now; the tags later, with the CLI |
| Q5 | the per-project accent colour | after the window title — it is the extra, not the fix |
| Q6 | is the Backlog tab allowed to **create** items? | yes — a backlog item is not a recipe, so ADR-0003 does not forbid it |
