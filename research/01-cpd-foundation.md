# CPD — the foundation

> Reference. Facts about
> [Crafted-Prompt Development](https://github.com/lthiagol/crafted-prompt-development), no
> prescriptions.
> Source: `github.com/lthiagol/crafted-prompt-development` @ `stable` (4 commits,
> 2026-07-08) — **a single `README.md`, 480 lines**. The public README is a *stripped
> derivative* of a longer private manifesto; it omits the crafted-prompt anatomy, the
> artifact ladder, the lifecycle escape hatches and single-agent mode.

## 1. The thesis

> "Spec-driven development front-loads formal specification at a point where you know the
> least about the problem. Agentic coding is strongest at back-loaded iteration.
>
> The distinction that matters: **CPD front-loads *intent* (what problem, what scope, what
> 'done' means) but refuses to front-load *implementation*.** Intent is knowable early and
> cheap to state; implementation design is discovered by building. That asymmetry is the
> whole argument."

**Already the stance we want.** CPD is anti-spec-driven and TDD-agnostic — it never
mentions TDD. It requires executable verification and distrusts green suites (principle
11). "Avoid TDD and spec-driven development" and "base the project on CPD" are compatible,
not in tension. **CPD is the good parts.**

## 2. The 14 principles

| # | Principle | Operative claim |
|--:|---|---|
| 1 | Discover before specifying | Default for *new* problems. Skippable only when scope is clear — *"Skipping discovery is an explicit, declared choice; defaulting to a prompt for unfamiliar work is the error CPD is built to prevent."* |
| 2 | Living context documents | `AGENTS.md`, `DESIGN.md`, `ARCHITECTURE.md` carry what no single prompt can hold. Keep them updated. |
| 3 | Minimal ceremony by default | *"Heavy ceremony is a deliberate escalation, not the baseline."* |
| 4 | Size-aware routing | *"Forcing large ceremony onto small work is more damaging than missing ceremony on large work."* |
| 5 | Promotion over premature ceremony | Light → heavy is one-way and validated. You never need to know the full scope upfront. |
| 6 | The prompt is the primary artifact — and survives its execution | Intent/problem/scope/ACs **immutable** during execution. Guidance, evidence and the audit trail move. *"The prompt carries the audit trail between code reviews."* |
| 7 | Compose, groom, review — then execute | *"Would I trust a capable agent with this prompt as it stands?"* **Prompt review (stage 4) ≠ code review (stages 6/8/10).** |
| 8 | Prompts define intent, not implementation | Not an upfront blueprint. Orientation docs are welcome; the blueprint is not. |
| 9 | Verification is auditable, not aspirational | *"The tests pass" is a claim. `cargo test -p foo` → exit 0 is evidence. **Bypass is recorded as explicit debt.**" |
| 10 | Scope-down visibility | *"A prompt marked 'done' that silently ships 60% of the intended scope is a broken contract."* Partial scope is legitimate; hidden scope is not. |
| 11 | Green tests don't imply correct behavior | A passing suite is a **precondition** for review, not the conclusion. *"Review must independently try to break the change."* |
| 12 | Agent discoveries feed forward | Discoveries that contradict the prompt's assumptions are captured, not silently acted on. |
| 13 | Code review by agents is non-negotiable | Round 1 self-review (structurally biased) + Round 2 external agent, *"ideally running a different model"*. **"Code review by agents is mandatory. Code review by humans is an extra."** |
| 14 | Human validation and incremental refinement | Small adjustments are **follow-up prompts**, not silent edits. *"The human stays in the loop at the seam between executions — not as a bystander, and not as a micromanager, but as the validator."* |

## 3. The 13-stage lifecycle

| # | Stage | Owner | What happens |
|--:|---|---|---|
| 0 | Discover | coord | *"Hold a chat session: explore the problem, challenge assumptions… **No prompt yet.**"* |
| 1 | Draft | coord | Write the prompt: intent, problem, scope, ACs, design decisions. |
| 2 | Groom | coord | Stress-test scope, verify AC falsifiability. Decompose into work packages and guidance steps. |
| 3 | Specify | coord | Lock it: `prompt_status: ready`. Open questions resolved. |
| 4 | Approve | coord | Pass the readiness gates. Ready for a runner to claim. |
| 5 | Claim & execute | runner | Implement, follow guidance steps *(or deviate with a note)*, pass ACs with evidence. |
| 6 | Self-review | runner | Round-1: lesson-pattern pre-screen, file self-findings. |
| 7 | Complete | runner | Mark complete; lifecycle → `self-reviewed`. |
| 8 | External review | coord | Round-2: examine diff, verify tests, check alignment with intent. File findings. |
| 9 | Remediate | runner | Address findings, fix code, re-run verification, mark findings resolved. |
| 10 | Re-review | coord | Verify remediations. Clean → advance. New issues → back to 9. **"Only the coordinator decides when the loop exits."** |
| 11 | Document | coord | Capture lessons, deviations, evidence. |
| 12 | Hand off | coord | Mark closed. **"The human validates the outcome"** and decides: accept, revise, or refine via a follow-up prompt. |

```
S0 Discover → S1 Draft → S2 Groom → S3 Specify → S4 Approve
  → S5 Claim&execute → S6 Self-review → S7 Complete
  → S8 External review → S9 Remediate → S10 Re-review → S11 Document → S12 Hand off
S10 -.->|findings remain| S9     S12 -.->|next prompt| S0
```

**Runner-owned:** 5, 6, 7, 9. **Coordinator-owned:** 0–4, 8, 10, 11, 12.

### Escape hatches (omitted from the public README)

- **Abort to Draft (stage 5 → 1).** If execution proves the *premise* wrong — neither a
  deviation nor a discovery — *"the runner halts, records what was found, and the prompt
  returns to Draft. Rare, and always declared — never silent."*
- **Partial scope is complete, not open.** *"Done with declared gaps" is a legitimate
  terminal state; "done with hidden gaps" is not.* Deferred ACs are re-filed as follow-up
  prompts.

### Single-agent mode

*"Single-agent mode is the **minimum viable** CPD."* Two accommodations: the session
boundary (L5) degrades to a **turn** boundary, and Round-2 review is **explicitly
weaker** (no model diversity).

## 4. The crafted prompt — seven sections

| Section | Authored | Mutable during execution? |
|---|---|---|
| Intent | Stage 1 | **No** |
| Problem | Stage 1 | **No** |
| Scope (incl. out-of-scope) | Stage 1 | **No** |
| Acceptance criteria | Stage 1 | **No** — new ACs → follow-up prompt |
| Guidance (suggested steps) | Stage 1, decomposed at Groom | **Yes** — with a deviation note |
| Evidence | Stages 5–7 | Grows |
| Audit trail | Stages 5–11 | Grows |

> *"Intent, problem, scope, and ACs are fixed once approved. If execution proves the
> premise wrong, that is an abort — not an edit. Guidance, evidence, and the audit trail
> are the only parts that move during execution."*

### The artifact ladder — size-aware routing, mechanized

| Rung | Artifact | When it fits |
|--:|---|---|
| 0 | Idea / note | A hunch, an observation, an agent discovery |
| 1 | Brief | Problem understood, scope small |
| 2 | Crafted prompt | Intent + ACs + guidance all add value |
| 3 | Decomposed prompt | Multiple independent units, one intent |
| 4 | Prompt family | Many prompts with cross-dependencies |

Promotion is **one-way and validated**; **demotion is allowed and unceremonious**.

## 5. The roles, and how separation is enforced

> *"CPD makes the plan-review boundary **a structural guarantee, not a matter of
> discipline**."*

- **Coordinator** — does not write code. Owns prompt quality and review quality. *"Ideally
  powered by a **different model** from the runner."*
- **Runner** — executes and fixes. *"Does not plan new work or evaluate its own output as
  the final word."* *"Produces evidence — not prose, but command output and exit codes."*

### The L5 law (session boundary)

> *"The runner session that produces code (stages 5–7) is never the same session that fixes
> findings (stage 9). The coordinator's review session (stages 8, 10) is never the same
> session that planned the prompt (stages 0–4). **The boundary is structural, not a matter
> of discipline.**"*

Four hand-off contracts, each carrying **Data / Session-boundary / Evidence**:
(a) 4→5 coordinator→runner · (b) 7→8 runner→coordinator · (c) 8→9 coordinator→runner ·
(d) 9→10 runner→coordinator.

### The honest caveat (from v1's milestone 142)

> *"The session id remains a string the agent can type. **Real L5 enforcement is a separate
> design discussion; the manifesto itself calls single-agent mode 'the minimum viable
> CPD'.** … **This milestone is the audit trail, not the enforcement.**"*

**We inherited this caveat.** CPD's separation is real but not tamper-proof. The tool
should make the honest case explicit rather than overclaim.

## 6. What CPD leaves to the implementer

| Concern | CPD | The tool |
|---|---|---|
| Stage sequence, role binding | ✔ | manifest + lint |
| Two-round review | ✔ | findings registry |
| L5 session boundary | ✔ | ⚠️ **advisory audit only** |
| Prompt anatomy + immutability | ✔ | schema + gates |
| Evidence discipline | ✔ | verifier + lint |
| Intent/implementation split | ✔ | ⚠️ **v1 contradicted it** (steps were gates) |
| Artifact ladder | ✔ | rungs 0–3 (rung 4 unimplemented) |
| **Session identity** | ✖ | the tool's job |
| **Human-in-the-loop surface** | ✔ *"the validator at the seam"* | **the tool's job — and now the main part** |
| **Automation / orchestration** | ✖ | the tool's job |
| **Prompt-injection defense** | ✖ | the tool's job |

### The five resolved open questions — the extension surface

1. Rename `spec_status` → `prompt_status`? → **Defer** (breaking; re-document first)
2. A `guidance` concept distinct from `steps`? → **Re-frame `steps`** as guidance
3. Where do deviation notes live? → **On the prompt's audit trail.** No new array.
4. Does post-execution review need its own kind? → **No** — `finding.phase` distinguishes
5. Rename `mp brief` → `mp compose`? → **Defer** with the terminology shift

### Where CPD is silent — gaps the tool must fill

1. **L5 enforcement.** The audit trail exists; enforcement does not.
2. **Rung 4** (prompt family) — coordinated prompts with ordering.
3. **Human surface.** CPD says *"the human validates the outcome"* — nothing about *how*.
   This is now our main deliverable.
4. **Discovery (stage 0)** is "free-form chat" — no mechanism. A human interface could
   offer real structure here.
5. **No prompt-injection defense** in the methodology.
6. **No LICENSE file** in the public repo (legally ambiguous for reuse).

## 7. Extension candidates

CPD will be updated alongside the tool. Candidates, from what v1 learned:

| Candidate | Rationale |
|---|---|
| **Discovery made structural** | Stage 0 is the most valuable stage and the least supported. A human interface could offer real scaffolding. |
| **Evidence attestation** | CPD requires "commands and exit codes" but doesn't require the tool to *check* the command matched anything. Our thesis: *a proof that matched nothing is not a proof*. |
| **Honest L5** | CPD claims structural separation; v1's audit could not enforce it. Either enforce it, or state the limitation in the methodology. |
| **A human-facing review surface** | CPD principle 14 puts the human at the seam but gives no instrument for it. |
| **Terminology** | v1's CLI said "spec" everywhere while CPD said "prompt", and the data model kept `spec_status`. The vocabulary must match the methodology — see [ADR-0004](../docs/decisions/0004-naming.md). |
| **Ladder rungs in the UI** | Rungs 0–4 should be visible and navigable, not just 0–3. |

**Decide with care.** CPD is published with a stated purpose: *"shared as a reference for
the agentic-coding community. Use it, adapt it, evolve it."* Extensions should be
**additive and general**, not product-specific.

## 8. The most important line

> *"Code review by agents is mandatory. Code review by humans is an extra."*

Combined with principle 14 (*the human validates the outcome*) and premise 3 (*the human
interface is the main part*), the shape is clear:

**The human is not the reviewer — the human is the *validator and the manager*.** The
interface is where the human decides *what* to build next, *whether* it was built
correctly, and *what to ask for next* — while the agents do the reviewing.

v1's TUI was a read-mostly view and every mutation had to go through the CLI. That is the
mistake to not repeat.
