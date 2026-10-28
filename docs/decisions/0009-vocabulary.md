# 0009. Vocabulary — the kitchen, the crew, and the machine's records

- **Status:** ✅ accepted
- **Date:** 2026-10-27
- **Deciders:** the chef

The **decision of record** is [docs/07-vocabulary.md](../07-vocabulary.md) — what the
chef decided, in the chef's words. This record is the **choice** and its reasoning:
why these names and not the others. The two point at each other, and where they differ
the decision of record wins.

## Context

v1's root cause was lexical: the CLI said `spec`, the data model kept `spec_status`,
the TUI said `milestone`, and CPD said `prompt`. Four words, no canonical choice — and
the doc drift that followed, with nothing checking the copies.

This project had the same disease in a different costume. The project was *Kitchen
Brigade*, the agent CLI was `scale`, the roles were themed (*line cook*, *taster*,
*expo*), the machine's records were named in metaphor where they were plainest
(`proofs`, `specials`, `prep-notes`), and the console's own tab was *Chopboard*. Every
name surfaces in a command, a file, a screen or a skill, so every name is a place a
reader can be misled.

Two forces pulled against each other, and the whole decision is the line between them.

| Force | What it wants |
|---|---|
| **The kitchen is this project's one distinctive idea** | the metaphor stays, and it goes deep |
| **A plan is read by agents, and the metaphor strains where the domain is hardest** | the machine's records and the agents' jobs stay plain |

The failure mode to avoid is v1's in reverse: not *no* vocabulary, but **two** — the
kitchen for the artifact and something else for everything an agent reasons about
mechanically — with no rule saying which is which.

## Decision

### 1. The rule

**The kitchen names the artifact and the place. Plain words name the machine's records
and the agents' jobs.**

| Layer | Vocabulary | Examples |
|---|---|---|
| The artifact | kitchen | recipe · the kitchen · the line · the board · service · tasting · sign-off · method |
| The records | plain | criterion · evidence · completion · lifecycle states · `depends_on` |
| The jobs | plain | chef · sous-chef · runner · reviewer · orchestrator |
| The plan's loose end | plain | backlog item · idea · feature · bug |

Two consequences, both rules now:

- The console says *tasting note*; the code says `Finding`. That is a **translation**,
  not a mismatch — the console is the surface where the kitchen is free to be charming.
- A recipe's row shows the **plain state word** (`closed`, `executing`) and the lane
  headers are **kitchen sections** (`Served`, `On the line`). The console never shows a
  lane word where a state word belongs.

### 2. The names

| Thing | Was | Is |
|---|---|---|
| The project | Kitchen Brigade | **maitre** |
| The console binary | `maitre` | `maitre` |
| The agent CLI | `scale` | **`maitre-cli`** |
| The module | `github.com/lthiagol/kitchen-brigade` | `github.com/lthiagol/maitre` |
| The plan directory | `kitchen/` | `kitchen/` — the kitchen is the place; the maitre runs it |
| The console's themed name | *maître d'* | **maitre** — it was never a job |
| A recipe's id | `R01` | **`RCP01`** |
| A backlog item's id | `S01` / `P01` | **`BL01`** — kind-neutral |
| The runner | `line-cook` | **runner** |
| The reviewer | `taster` | **reviewer** |
| The orchestrator | `expo` | **orchestrator** |
| The tree tab | `Chopboard` | **Board** |
| The order tab | `Path` | **Path** |
| The automation's tab | `Expo` | **Sous-Chef** |
| A criterion's record | `Proof` | **evidence** |
| The recipe-level stamp | `Evidence` | **completion** |

`scale` went because it *weighs*: the binary that records the plan is not a scale. The
two-word `maitre-cli` says what it is, and `maitre` stays the console the chef knows.

**The retired names still work on an old plan.** A plan written before the rename
carries `line-cook` or `taster`, and the rename must not loosen what the domain
refuses: a name the rename retired still names the role it stood for, so an old plan is
refused exactly like a new one. The aliases are the tool's, not the vocabulary's — the
records a new plan holds carry only the new names.

### 3. The backlog: one artifact, three kinds

The backlog had **two file types** — a *special* (a hunch) and a *prep note* (intent and
scope before criteria) — and two CLI verb groups. They became **one artifact** with a
`kind`.

| | |
|---|---|
| The artifact | a **backlog item**, id `BL01`, in `kitchen/backlog/` |
| Its kinds | `idea` · `feature` · `bug` |
| Its groups in the console | **Ideas** (needs shaping) and **Small tasks** (ready to work) — **derived** from the kind, never stored |
| The CLI | `maitre-cli backlog list · create · update · dismiss · promote` — **46 commands → 42** |
| An older plan's specials | read **read-only**, in a third group under the tab — the tool reads an old plan; it does not rewrite it |

```
idea ──shaped into──▶ a small task  (kind: idea → feature|bug, plus scope)
idea ──shaped into──▶ a recipe       (straight to the big thing)
small task ──promoted──▶ a recipe with ONE criterion and ONE evidence
```

The last line is the load-bearing one: **a small task is too small to be a big recipe,
not too small to be guarded.** It becomes a recipe with a single criterion, so every
gate still applies.

### 4. What is strings only

| | |
|---|---|
| The tab bar | `Board` and `Sous-Chef` |
| The window title | *maitre* |
| The CLI's help | `maitre-cli — the maitre agent CLI` |
| The lane headers | they stay kitchen, and they read as sections, never as a second name for a state |
| The mascot | a funny agentic maitre, drawn once and reused in the console header, the CLI's help and the README |

## The alternatives, and why they lost

| Option | Lost because |
|---|---|
| **Kitchen everywhere** — `TastingNote`, kitchen states, a plan format in metaphor | the plan format changes, agents read a plan in metaphor, and the metaphor strains exactly where the domain is hardest |
| **Drop the costume** — spec, plan, task, gate, proof, finding | that is v1's vocabulary, and the kitchen is this project's one distinctive idea |
| **Coming back to master-plan's names** | v1 has **six** ways to record a problem and **four** to record a note; that is the lexical-mismatch lesson this repo exists to hold |
| **Keep the roles themed** — line cook, taster, expo | a job is a plain thing: an agent that runs code is a runner, whoever the kitchen calls it. The themed names were the metaphor applied where the metaphor did not help |
| **Keep *maître d'* for the console** | the console is `maitre` now, and it was never a job — it is a surface, and the surface wears the product's name |
| **Two backlog file types, renamed** | two nouns for one list |
| **A small task closed without gates** | work that closes with no proof |
| **A second executable artifact type** | a second anatomy, a second gate set, a second close path |
| **`Proof` kept, `Evidence` dropped** | the chef chose otherwise: the journal's event kind is already `evidence`, and making the field agree with the journal is the better end state |

## Consequences

### Positive

- **One vocabulary, with a rule that says which layer speaks it.** A reader can tell
  why a recipe is a recipe and a criterion is a criterion.
- The kitchen survives whole — it is the project's one distinctive idea, and it is now
  where it earns its keep: the artifact and the place.
- The machine's records and the agents' jobs are plain, so a plan read by an agent
  reads like the tool's own vocabulary.
- The rename is **fenced, not hoped for**: the plan format's version moves with the
  format, and a tool whose rules are older refuses to write.

### Negative

- **The kitchen is thinner at the top.** *Line cook*, *taster* and *expo* were
  memorable, and the plain names are not. The console keeps the flavour in its lanes,
  its sign-off and its tasting notes.
- **Two vocabularies are still two vocabularies.** The rule is what keeps them from
  drifting, and a rule needs a reader.
- The ids changed (`R01` → `RCP01`), and a plan written before the rename needed a
  migration — the version guard is what made that safe.

### Follow-up

- The mascot: a funny agentic maitre, drawn once and reused. Its shape was left open
  until the words settled.
- **Courses are named and not built** — ladder rung 3. The rename did not decide what
  a course is.

## Supersedes

[ADR 0004](0004-naming.md) — the names it chose are replaced by § 2 above. What 0004
decided *beside* the names stands: the console's name, the plan directory, the kitchen
theme, and the artifact ladder.
