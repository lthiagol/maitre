# 0004. Naming

- **Status:** ✅ accepted — superseded in part on 2026-10-27; see
  [docs/07-vocabulary.md](../07-vocabulary.md), the decision of record
- **Date:** 2026-10-25
- **Deciders:** TBD

## Context

Every name in the product surfaces in a command, a file, a screen, or a skill. v1's
naming was its own root cause: the CLI said `spec`, the data model kept `spec_status`, the
TUI said `milestone`, and CPD said `prompt`. Four words, no canonical choice — and the
doc drift that followed.

So: one vocabulary, decided explicitly, mapped to CPD once.

**The names are decided here. Their definitions live in
[GLOSSARY.md](../../GLOSSARY.md)** — one source, one place to change a term.

## Decision

### The project

| | |
|---|---|
| Project | **Kitchen Brigade** |
| Tagline | *a brigade of agents for coding work* |
| Module | `github.com/lthiagol/kitchen-brigade` |
| Plan directory | `kitchen/` |

*Kitchen Brigade* over *Agentic Kitchen Brigade*: shorter, and a brigade is already a team
of workers — the "agentic" part is carried by the metaphor, not by a buzzword.

### The binaries — two, one formula

| Binary | Audience | Verb |
|---|---|---|
| `scale` | Agent CLI | *weigh* a recipe before calling it done — gates, coverage, readiness |
| `maitre` | Human console | the maître d' serves the Chef |

`scale` over `whisk`/`mise`: a scale is the first tool reached for in any kitchen, and
"weigh the recipe before you call it" is exactly what the gate system does. `maitre` over
`pass`/`board`: the maître d' is the one who greets the chef and takes the order, which
is what a management console does.

### The artifact ladder

`special → prep note → recipe → course → menu` — defined in
[GLOSSARY.md](../../GLOSSARY.md#the-plan). The reasoning for the top rung: a recipe card
has exactly our anatomy — the dish, why we make it, **what goes in / what doesn't**, the
method, how to know it's right, tasting notes, and the pass. `scope.in_scope` → *what
goes in*; `scope.out_of_scope` → *what doesn't*.

### The cast

Chef · maître d' · sous chef · line cook · taster · expo — defined in
[GLOSSARY.md](../../GLOSSARY.md#the-brigade).

The human-as-chef mapping is the load-bearing one: CPD's coordinator is the *sous chef* —
it composes, it does not own the kitchen. The chef decides and signs off.

### The skills

| v1 | Tool |
|---|---|
| `mp-flow` | **service** (the flow of service = the stage timeline) |
| `mp-coordinator` | **sous-chef** |
| `mp-runner` | **line-cook** |
| `mp-reviewer` | **taster** |
| `mp-orchestrator` | **expo** |

### The plan directory

```
kitchen/
  recipes/        RCP01-<slug>.json   ← the line
  courses/        C01-<slug>.json
  specials/       S01-<slug>.json     ← stock
  prep-notes/     N01-<slug>.json
  project.json
  activity.json
  decisions.json
```

One root, subdirectories per kind. The "stock vs. the line" distinction — specials and
prep notes are stock, recipes and courses are the line — is a **view in the UI**, not a
directory boundary. Two roots would cost a workspace spanning both and buy nothing,
since per-type subdirectories already allow selective commits (`git add kitchen/recipes`).

## Consequences

### Positive
- One vocabulary, decided explicitly, defined in one place.
- Every name reads on a board and in a terminal.
- The metaphor extends all the way down — a station = work package, a ticket = recipe,
  the pass = the hand-off.

### Negative
- **"Recipe" is not a general term.** CPD is a public reference for the agentic-coding
  community; a reader arriving from CPD will not guess "recipe". Mitigate: the glossary
  names the CPD counterpart for every term.
- **"Scale" and "maitre" are both common words** — discoverability suffers. Mitigate: the
  brew formula and the install output name them explicitly.
- The kitchen theme can over-extend. See [ADR-0005](0005-vocabulary.md) for where the
  theme stops.

## Resolved

- **Names:** project, binaries, plan directory, artifact ladder, roles, skills — all above.
- **Vocabulary rule:** [ADR-0005](0005-vocabulary.md).

## Open questions

- Should the skills keep a `service` meta-skill name, or is one name too many? The stage
  timeline has to live somewhere.
