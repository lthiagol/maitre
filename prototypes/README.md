# Prototypes

Throwaway artifacts that answer a design question. They are **not** the implementation —
they are the sketch you react to before the real thing gets built.

## The tabs — the bigger picture

```bash
cat prototypes/tab-recipes.txt    # the board — the main surface
cat prototypes/tab-backlog.txt    # specials and prep notes
cat prototypes/tab-path.txt       # the execution order
cat prototypes/tab-log.txt        # the activity journal
cat prototypes/tab-settings.txt   # the plan, the view, the keys
```

**These are generated from the real console** — `maitre --render --plain` output, saved.
They are not sketches; they are what the tool actually draws.

| Tab | Holds | Job |
|---|---|---|
| **Recipes** | every recipe, grouped by what needs the chef | view · decide |
| **Backlog** | specials and prep notes | manage the backlog |
| **Path** | the execution order, and what blocks what | run |
| **Log** | the activity journal | see the audit trail |
| **Settings** | the plan, the view, the keys | — |

v1's `raul` had **seven** lanes — Overview · Milestones · Path · Backlog · Ideas ·
Autopilot · Settings. Four are gone:

| v1 lane | Why it is gone |
|---|---|
| Overview | the triage-ordered board *is* the overview; a second tab would restate it |
| Milestones | → **Recipes**, our vocabulary |
| Ideas | folded into **Backlog** — a special is our idea |
| Autopilot | v2 has no orchestration; the running recipes are a section of the board |

## The layout variants

The frame — a list on the left, the detail on the right — was settled by comparing four
radically different answers. They are kept as the record of what was compared.

```bash
cat prototypes/variant-inventory.txt   # what exists
cat prototypes/variant-kanban.txt      # where everything is
cat prototypes/variant-focus.txt       # what the chef is looking at
```

The decided one is **triage first, inside a split frame** — see
[`docs/05-human-interface.md`](../docs/05-human-interface.md).

## What the data shows

Seven recipes spread across every lane, so each tab has something to say:

| | | |
|---|---|---|
| R01 Sharpen the AC gate | draft | |
| R02 Sign off the menu | ready | awaiting the chef |
| R03 Wire the review queue | remediation | 1 open tasting note |
| R04 Add the backlog view | approved | waiting on R02 |
| R05 Sharpen the proof contract | executing | on the line |
| R06 Draft the glossary | closed | served |
| R07 Record the review protocol | draft | waiting on R05 |

Plus four specials and two prep notes in the backlog.

## Rules for prototypes

1. **Throwaway from day one, and clearly marked as such.** The directory is named
   `prototypes/` so a reader knows.
2. **Trivial to run.** One file, one `cat`. The product is terminal-first, so the
   prototypes are terminal text — the medium the real thing will use.
3. **No persistence.** State lives in the file.
4. **Skip the polish.** The point is to learn something fast.
5. **Capture the answer, not the code.** When a layout is chosen, the decision goes into
   [`docs/05-human-interface.md`](../docs/05-human-interface.md); the prototype stays as
   the record of what was compared.

Adapted from the upstream
[`prototype`](https://github.com/mattpocock/skills/blob/main/skills/engineering/prototype/SKILL.md)
skill (MIT), whose UI branch produces an HTML page with switchable variations. We tried
that first; the text mockups are a better fit for a terminal product and are reviewable in
a terminal, a diff, or a chat message.
