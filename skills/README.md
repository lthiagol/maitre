# Vendored skills

Third-party skills we depend on, copied into this repo so they can be adapted, pinned,
and shipped with the tool.

**Upstream.** `github.com/mattpocock/skills` — MIT. See [LICENSE](LICENSE) and
[UPSTREAM.md](UPSTREAM.md) for provenance.

**Pinned to a commit, not `main`.** master-plan v1 vendored these at
`upstream_version: main` and never re-synced; both copies had drifted by the time this
repo started (31 changed lines in `codebase-design`, 73 in `diagnosing-bugs`). A pinned
commit is the fix — see [sync.sh](sync.sh).

## What we took, and why

| Skill | Why we need it |
|---|---|
| [codebase-design](codebase-design/SKILL.md) | The `core/`-vs-surface seam is our central architectural commitment. This is the vocabulary for it: module, interface, depth, seam, adapter, leverage, locality. **Our own `core` package gets judged by the deletion test.** |
| [writing-for-agents](writing-for-agents/SKILL.md) | Most of what we write is for agents. This governs skills and `AGENTS.md` — including the positive-voice rule our house style now follows. |
| [domain-modeling](domain-modeling/SKILL.md) | Keeps [GLOSSARY.md](../GLOSSARY.md) and the ADRs sharp. It is the discipline behind our CPD↔tool vocabulary. |
| [diagnosing-bugs](diagnosing-bugs/SKILL.md) | Build a feedback loop before hypothesising. No TDD conflict: its regression-test phase is regression testing, not test-first authoring. |

## What we did not take

| Skill | Why not |
|---|---|
| `code-review` | **Adapted instead of vendored** — it becomes our `taster` role. Its Spec axis must read a recipe from `kitchen/`, not an issue tracker. See `docs/decisions/`. |
| `wayfinder` | Tracker-coupled and decision-shaped. We took its *ideas* — map-as-index, fog of war, name-not-id, HITL-vs-AFK — not the skill. |
| `prototype` | The method transfers; the artifacts assume a web stack. We already have `maitre --render` for the equivalent. |
| `implement`, `implement-spec` | Call `tdd` (premise 2) and assume multi-agent worktree orchestration (premise 3). |
| `improve-codebase-architecture` | Nothing to scan in a small codebase. Its HTML report is a natural fit for the web-shell option later. |

## Re-syncing

```bash
skills/sync.sh          # re-copy at the pinned commit, report what changed
skills/sync.sh --bump   # move to a newer commit, then re-copy
```

The script never edits a vendored file by hand. If a skill needs changing for Kitchen
Brigade, it becomes **our** skill under `templates/skills/` — the vendored copy stays
pristine so the diff against upstream stays readable.
