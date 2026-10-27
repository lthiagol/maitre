# Vocabulary — the decision of record

Decided by the chef on 2026-10-27, after the audit ledger closed. This doc is the
paper the rename is built from; `GLOSSARY.md` is rewritten to match it, and where the
two differ this doc wins.

## The rule

**The kitchen names the artifact and the place. Plain words name the machine's records
and the agents' jobs.**

| Layer | Vocabulary | Examples |
|---|---|---|
| The artifact | kitchen | **recipe**, the kitchen, the line, the board, service, tasting, sign-off, method |
| The records | plain | criterion, **evidence**, lifecycle states, `depends_on`, the journal's kinds |
| The jobs | plain | chef, sous-chef, **runner**, **reviewer**, **orchestrator** |
| The plan's loose end | plain | **backlog item**, **idea** / **feature** / **bug** |

Two consequences that used to be accidents and are now rules:

- The console says *tasting note*; the code says `Finding`. That is a **translation**,
  not a mismatch — the console is the surface where the kitchen is free to be charming.
- A recipe's row shows the **plain state word** (`closed`, `executing`) and the lane
  headers are **kitchen sections** (`Served`, `On the line`). The console never shows a
  lane word where a state word belongs.

## The names

| Thing | Was | Is |
|---|---|---|
| The project | Kitchen Brigade | **maitre** |
| The console binary | `maitre` | `maitre` |
| The agent CLI | `scale` | **`maitre-cli`** |
| The module | `github.com/lthiagol/kitchen-brigade` | `github.com/lthiagol/maitre` |
| The plan directory | `kitchen/` | `kitchen/` — the kitchen is the place; the maitre runs it |
| A recipe's id | `R01` | **`RCP01`** |
| A backlog item's id | `S01` / `P01` | **`BL01`** — kind-neutral, so an idea can be shaped without an id change |
| The runner | `line-cook` | **runner** |
| The reviewer | `taster` | **reviewer** |
| The orchestrator | `expo` | **orchestrator** |
| The tree tab | `Chopboard` | **Board** — the plan's shape: what waits on what |
| The order tab | `Path` | **Path** — the plan's order: what can go now |
| The automation's tab | `Expo` | **Sous-Chef** — the crew and the machine's vital signs |
| A criterion's record | `Proof` | **evidence** |
| The recipe-level stamp | `Evidence` | **completion** |

`scale` went because it *weighs*: the binary that records the plan is not a scale. The
two-word `maitre-cli` says what it is, and `maitre` stays the console the chef knows.

## The backlog: one artifact, three kinds

The backlog had **two file types** — *special* (a hunch) and *prep note* (intent and
scope before criteria) — and two CLI verb groups. They became **one artifact** with a
`kind`.

| | |
|---|---|
| The artifact | a **backlog item**, id `BL01`, in `kitchen/backlog/` |
| Its kinds | `idea` · `feature` · `bug` |
| Its groups in the console | **Ideas** (needs shaping) and **Small tasks** (ready to work) — **derived** from the kind, never stored |
| Its fields | title, body, scope, kind, tags, status, created, updated, promoted_to |
| The CLI | `maitre-cli backlog list · create · update · dismiss · promote` — **46 commands → 42** |
| An older plan's specials | read **read-only**, in a third group under the tab — the tool reads an old plan; it does not rewrite it |

One artifact because v1's lesson is that *N×N surfaces over one small domain* is where
drift comes from: two file types for one list is two nouns for one thing, and the
console's two groups do not need two file types to exist.

### How the loose end becomes work

```
idea ──shaped into──▶ a small task  (kind: idea → feature|bug, plus scope)
idea ──shaped into──▶ a recipe       (straight to the big thing)
small task ──promoted──▶ a recipe with ONE criterion and ONE evidence
```

The last line matters: **a small task is too small to be a big recipe, not too small to
be guarded.** It becomes a recipe with a single criterion, so every gate still applies.
The alternative — letting a small task close without gates — is a hole in the domain's
promise, and the third alternative, a second executable artifact type, is a second
anatomy, a second gate set and a second close path.

`dismiss` comes to the small tasks too. Before, only specials could be dismissed, and no
one could say why.

## The plan format changes (domain version 2 → 3, then 4, then 5)

| Change | From | To |
|---|---|---|
| Recipe ids | `R01` | `RCP01` |
| The backlog | `kitchen/specials/` + `kitchen/prep-notes/`, two types | `kitchen/backlog/`, one type with a `kind` |
| A criterion's record | `proofs` | `evidence` |
| The recipe-level stamp | `evidence` | `completion` |
| Journal roles | `line-cook`, `taster`, `expo` | `runner`, `reviewer`, `orchestrator` |

The version guard is what makes this safe: a plan written under the newer version
refuses a stale tool that speaks the older one. The ids are version 3; **version 4**
came when the fingerprint started covering the format's facts, so a format change
moves the version on its own instead of by hand; **version 5** came with the backlog's
collapse into one artifact.

## What is strings only

| | |
|---|---|
| The tab bar | `Board` and `Sous-Chef` replace `Chopboard` and `Expo` |
| The window title | *maitre* replaces *Kitchen Brigade* |
| The CLI's help | `maitre-cli — the maitre agent CLI` |
| The lane headers | they stay kitchen, and they read as sections, never as a second name for a state |
| The mascot | a funny agentic maitre, drawn once and reused in the console header, the CLI's help and the README |

## The alternatives, and why they lost

| Option | Lost because |
|---|---|
| Kitchen everywhere — `TastingNote`, kitchen states, a plan format in metaphor | the plan format changes, agents read a plan in metaphor, and the metaphor strains exactly where the domain is hardest |
| Drop the costume — spec, plan, task, gate, proof, finding | that is v1's vocabulary, and the kitchen is this project's one distinctive idea |
| Coming back to master-plan's names | v1 has **six** ways to record a problem and **four** to record a note; that is the lexical-mismatch lesson this repo exists to hold |
| Two backlog file types, renamed | two nouns for one list |
| A small task closed without gates | work that closes with no proof |
| A second executable artifact type | a second anatomy, a second gate set, a second close path |
| `Proof` kept, `Evidence` dropped | the chef chose otherwise: the journal's event kind is already `evidence`, and making the field agree with the journal is the better end state |

## The order of the work

| # | Step | Touches |
|--:|---|---|
| 1 | The project is **maitre**: module path, titles, help text, docs | `go.mod`, `cmd/`, docs |
| 2 | The CLI is **`maitre-cli`** | `cmd/`, Makefile, `scripts/console-shot.py`, docs |
| 3 | Recipe ids are **`RCP01`**; the domain version goes to **3**; `kitchen/` is migrated | the normaliser, the seeder, the goldens |
| 4 | The roles are **runner**, **reviewer**, **orchestrator** | the journal, the self-review gate and its tests |
| 5 | The tabs are **Board** and **Sous-Chef** | UI strings |
| 6 | The backlog is **one item with a kind** | the model, the store, the CLI, the console |
| 7 | **Evidence** and **completion** | the verifier, the gates, the views, the goldens |
| 8 | The **mascot** | an asset |
| 9 | The docs, the goldens and the tripwires are regenerated | the whole tree |

Each step goes through the reviewer before the next one starts.

## Not in this bundle

| | |
|---|---|
| **F-59** | `Audit.Decisions` has a shape and no writer — [06-proposals.md](06-proposals.md) § 2 holds the remedy |
| **F-60** | the plan has questions, and the `questions.resolved` gate can never fire |
| **F-61** | `courses/` is named, created and never written |
| **The course question** | held by the chef |
| **The mascot's shape** | drawn after the words settle |

The backlog collapse here is a *different* change from those: it renames a rung and
merges two types, and it does not decide what a course is.
