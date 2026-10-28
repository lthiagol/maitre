# maitre

> A ground-up rewrite of [`master-plan`](https://github.com/lthiagol/master-plan) v1
> (v1.0.0-rc6), built on
> [CPD](https://github.com/lthiagol/crafted-prompt-development) as its foundation.
> No code carried over. Human interface first, agent CLI derived from it.

| | |
|---|---|
| Agent CLI | `maitre-cli` |
| Human console | `maitre` |
| Plan directory | `kitchen/` |

## Status

**Phase 3 — build.** The domain model is designed explicitly; everything else is
discovered by using it.

| | Status |
|---|---|
| `internal/core` — model, lifecycle, gates, store, derived views, verifier | ✅ |
| `internal/cli` — the `maitre-cli` surface | ✅ slice 1 |
| `internal/ui` — the `maitre` board and detail pane | ✅ slice 2 |
| `docs/03-domain-model.md` · `04-operations.md` · `05-human-interface.md` | 🔄 drafts |

```bash
make build && make test
```

## The board

```
 maitre · my-project
──────────────────────────
12 recipes · /path/to/kitchen

▸ Grooming      2 │ RCP03  draft     Sharpen the AC gate
  Sign-off      1 │ RCP07  ready     Add the backlog view
  On the line   3 │ RCP11  executing  Wire the review queue
  Tasting       2 │
  On hold       1 │
  Served        3 │

────────────────────────────────────────────────────────
↑↓ move  ·  tab lane  ·  enter detail  ·  r reload  ·  q quit
```

`maitre --render` prints the same view as a static snapshot — no terminal needed, so it
works in CI, in a PR, or in a chat message.

## Decisions

| # | Decision |
|---|---|
| [0001](docs/decisions/0001-stack.md) | Stack: Go + Charm |
| [0002](docs/decisions/0002-human-interface-first.md) | Build order: human interface first |
| [0003](docs/decisions/0003-surface-split.md) | Who writes what: artifact vs. management |
| [0004](docs/decisions/0004-naming.md) | Naming: maitre · `maitre-cli` · `maitre` · the ladder — superseded by 0009 |
| [0005](docs/decisions/0005-vocabulary.md) | Vocabulary: generic methodology, specific product — superseded in part by 0009 |
| [0006](docs/decisions/0006-verification.md) | Verification — who runs it, and what evidence is |
| [0007](docs/decisions/0007-operation-layer.md) | The operation layer |
| [0008](docs/decisions/0008-artifact-is-the-channel.md) | The artifact is the channel — no agent-to-agent messaging |
| [0009](docs/decisions/0009-vocabulary.md) | Vocabulary — the kitchen, the crew, and the machine's records; the decision of record is [docs/07-vocabulary.md](docs/07-vocabulary.md) |

## Reference

| Doc | Contents |
|---|---|
| [AGENTS.md](AGENTS.md) | House style, premises, working rules — **read first** |
| [GLOSSARY.md](GLOSSARY.md) | The domain vocabulary and the CPD↔tool mapping — single source |
| [docs/01-vision.md](docs/01-vision.md) | Thesis, goals, non-goals, open questions |
| [docs/02-principles.md](docs/02-principles.md) | 13 principles, each traced to a lesson |
| [docs/03-domain-model.md](docs/03-domain-model.md) | Anatomy, lifecycle, gates, plan format |
| [docs/04-operations.md](docs/04-operations.md) | The operation set — what the tool can do, and who does it |
| [docs/05-human-interface.md](docs/05-human-interface.md) | The console — four jobs, four layout candidates |
| [prototypes/README.md](prototypes/README.md) | The board layout prototypes — four answers, switchable |
| [skills/README.md](skills/README.md) | Vendored third-party skills — what we took, what we left, how to re-sync |
| [research/00-legacy-analysis.md](research/00-legacy-analysis.md) | v1: assets, defects, root causes |
| [research/01-cpd-foundation.md](research/01-cpd-foundation.md) | CPD: principles, stages, ladder, gaps |

## Layout

```
cmd/maitre-cli/  the agent CLI
cmd/maitre/     the human console
internal/core/  domain model · lifecycle · gates · store · views · verifier   ← stdlib only
internal/cli/   the maitre-cli command surface                                   ← stdlib only
internal/ui/    the maitre console (bubbletea v2)
skills/         vendored third-party skills, pinned to an upstream commit
templates/      our own skills (to come)
prototypes/     throwaway design sketches
docs/           design docs
  decisions/    ADRs
research/       reference material
```
