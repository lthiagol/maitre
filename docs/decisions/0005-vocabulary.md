# 0005. Vocabulary — generic methodology, specific product

- **Status:** ✅ accepted
- **Date:** 2026-10-25
- **Deciders:** TBD

## Context

The project has two layers that both talk about the same objects:

| Layer | Audience | Purpose |
|---|---|---|
| **CPD** — the methodology | the agentic-coding community, publicly | principles, stages, roles, vocabulary |
| **The tool** — Kitchen Brigade | people using it | the same ideas, implemented |

CPD is currently generic and descriptive: *crafted prompt*, *acceptance criteria*,
*scope*, *evidence*, *findings*. The tool wants specific names: *recipe*, *chef*,
*scale*, *maitre*.

v1 tried to have both without deciding, and produced the drift that is now the project's
worst reference case: the CLI said `spec`, the model kept `spec_status`, the TUI said
`milestone`, CPD said `prompt`. Four words for one artifact, no canonical choice.

## Decision

### 1. CPD stays generic and descriptive

CPD's vocabulary is the reference language for the methodology. It does not take on the
product's theme. CPD is maintained by us and is not widely adopted, so this is a choice
about clarity, not constraint — but the split is worth keeping.

### 2. The tool uses specific names, where the semantics and the feeling match

The kitchen theme applies to the things a person *encounters*: the artifact, the roles,
the binaries, the ladder rungs, the skills.

### 3. Mechanical terms stay plain

Anything an agent reasons about mechanically stays descriptive. Theming
`acceptance_criteria` → *tasting criteria* would cost every agent a translation and
invite confusion at the exact point where precision matters.

| Themed | Kept plain |
|---|---|
| Artifact: recipe · ladder: special → prep note → recipe → course → menu | `acceptance_criteria` · `scope` · `evidence` · `findings` · `work_packages` |
| Roles: chef · maître d' · sous chef · line cook · taster · expo | lifecycle states · lanes · gate names · `status` · `review_state` |
| Binaries: `scale` · `maitre` · skills: service · sous-chef · line-cook · taster · expo | |

### 4. The mapping lives in exactly one place

**[GLOSSARY.md](../../GLOSSARY.md)** is the single source for every term and its CPD
counterpart. It is generated into the CLI reference and the skills, so no doc can hold a
third spelling. This doc states the rule; the glossary holds the table.

## Consequences

### Positive
- CPD stays reusable by people who have never seen the tool.
- The tool keeps a coherent, memorable surface.
- The v1 failure mode — unacknowledged drift — is prevented structurally: one table, one
  source, generated everywhere.

### Negative
- Two vocabularies to hold. Mitigated by the glossary being the single source and the
  mapping being generated, not hand-written.
- Someone reading CPD and then the tool has to translate. Acceptable: the glossary is the
  translation.

### Follow-up
- The CPD README gains a short note pointing at the tool as one implementation, without
  importing the theme.
- The glossary becomes a generated artifact in `docs/` and in the CLI's `--help`
  reference.

## Open questions

- Do the skills use the kitchen names in their *text*, or only in their filenames? Prose
  that says "chef" may confuse an agent that has never read CPD. Open: keep skill
  *prose* generic, use kitchen names for the *artifact* references only.
