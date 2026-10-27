# 0008. The artifact is the channel

- **Status:** ✅ accepted
- **Date:** 2026-10-25
- **Deciders:** TBD

## Context

v1 had **three places to look** for one fact, and agents messaged each other on top of
that:

| | |
|---|---|
| the plan | `milestones/*.json` — findings, review records, the audit trail |
| the reviews registry | `reviews.json` — a plan-level copy of the review state |
| the autopilot session | `<plan>/autopilot/<id>/session.json` — its own event log, its own role state |

An agent that finished a review had to *tell* the agent that would fix the findings. That
is a messaging protocol with no shared vocabulary, and v1's dogfood log is full of the
consequences: findings described in chat instead of on the artifact, verdict and evidence
disagreeing, "the reviewer says ok but the milestone says otherwise".

Meanwhile v1's **schema already had the answer and the CLI never shipped it**. `Finding`
carried `thread: Vec<…>`, `summary`, `rationale`, `confidence`, `tags` — the fields for a
finding carrying its own conversation. But `mp reviews finding add` only took
description, severity, category, phase and an anchor. The verbs were never written.

## Decision

### 1. No agent-to-agent messaging

Agents do not send each other what they found or what they fixed. The **recipe is the
blackboard**: the reviewer writes onto it, the fixer reads from it and writes back.

The task is the recipe. There is no task entity, no inbox, no message queue. A recipe
under review has open findings; the findings are the work queue.

### 2. The finding's description is the review; the comments are the conversation

- `finding.description` — what is wrong and why it matters. Several sentences is fine.
- `review.comments[]` — one flat list, each with an author and an optional finding
  reference. The supporting evidence, the fixer's account of what they did, the
  reviewer's confirmation.
- `review.rounds[]` — the verdict: reviewer, `ok` or `changes-needed`, and notes.

The comments stay **flat** rather than nested under the finding (as v1's `thread` had
them). The console groups them by finding when it renders; grouping is a presentation
concern, and a flat list also supports a comment about the review as a whole.

### 3. The review loop goes backwards

`maitre-cli finding reopen` sends a finding back into the queue when a fix does not hold. The
outcome is cleared — the substance of the previous attempt lives in the comments, which
are append-only — and the objection is appended as a comment, so the next attempt reads
why it was sent back.

Without this the loop only went forward, and "the review happened" was not the same as
"the review is a loop".

### 4. The journal is the notification

An agent is not told that something changed; it looks. `maitre-cli log --since <seq>` returns
what is newer than the last sequence number it saw.

`activity.json` is the **only** event log. A recipe's own audit holds what the recipe
authored — deviations, decisions, documentation — never a copy of its events.

### 5. Detection is a watch concern, not a messaging one

Because the artifact is the channel, noticing that it changed is a polling or watching
problem. The console's mtime scan is therefore not a convenience — **it is the
coordination mechanism**. The same holds for an agent, which polls `maitre-cli log --since`.

### 6. The reviewer cannot advance a recipe that says no

The `review.verdict_ok` gate requires the **latest** round's verdict to be `ok`, not
merely that a round exists. (v1-equivalent hole: a reviewer could record
"changes needed" and the milestone advanced anyway, provided no findings were open.)

## Consequences

### Positive
- One place to look for a fact: the recipe, plus the journal for history.
- No coordination protocol to design, version, or debug.
- The review record and the work queue are the same thing.
- v1's `thread`-shaped schema finally gets used.

### Negative
- **A recipe can grow long.** Every comment, every finding, every evidence is on it. There is
  no summarisation, and a recipe with a hundred findings is a hundred findings of JSON.
  The console has to render it; an agent has to read it.
- **Polling is the cost of coordination.** Every agent that wants to know something
  changed has to look. A watch would be better; that is orchestration's job.
- **Comments are append-only, so a retracted statement stays.** The status fields say what
  is true *now*; the comments say what happened. Reading both is the reader's job.
- The recipe is the unit of concurrency. Two agents editing one recipe at once need the
  plan's advisory lock, which exists but is not exercised by parallel writers yet.

### Follow-up
- The console must render a finding's comments under it, or the thread is unreadable.
- The mtime scan moves up the list: it is coordination, not polish.
- Orchestration remains v2.2+, but it must not introduce a second event log.

## Open questions

- **How does an agent find the work that is waiting for it?** Today: read the plan, look
  for open findings. With orchestration, the queue is derived — but it must still be
  derivable without an orchestrator running.
- **Does a recipe get too big?** If it does, the answer is probably not splitting the
  artifact (that breaks the blackboard) but summarising — a "state of the review" section
  the agent can read first and drill into after.
