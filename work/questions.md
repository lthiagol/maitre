# The open questions — what waits on the chef

> Every question this work is carrying, and **who answers it**. A question opens when
> it stops being a detail; it closes when it is answered, and the answer moves to the
> **rulings table** in [tui/README.md](tui/README.md).
>
> **The file is checked at every slice's close.** The reviewer reads it and says
> whether a question was answered, missed, or newly raised. **A question that is not
> written here is a question that is lost.**

## The index

| # | The question | Who answers | Where it came from | What it blocks |
|--:|---|---|---|---|
| **Q8** | the **free edit** of a backlog item's fields — is it wanted at all? | **the chef** | T3b's brief ruled it out of the round: the three moves are what the chef needs, and a five-field form is a later slice, *if it is wanted at all* | a later slice |
| **Q9** | a backlog item's `scope` and `tags` are written **absent**; a recipe's are `[]`. Keep the compact record, or write the lists? | **the chef** | T3a's review: the same hazard in kind as the recipe's `null` — a reader that asks for the list gets nothing. The runner's reason for leaving it (*it would be a format change and a version bump*) does not hold: the fingerprint hashes the field **names**, so writing `[]` would not move the version | the next write of every item on disk |
| **Q10** | after a write, does the selection follow the item it was on? | **the chef** | T3b's review: shaping moves the item between the groups and the cursor keeps its **index**, so the selection lands on another item. Pre-existing, and console-wide — a sign-off moves a recipe between the triage sections the same way | the console's feel, everywhere |
| **Q11** | one marker or two for *"the domain's answer"*? | the orchestrator, unless the chef wants it | T3b's review: the detail's arrow on a promoted item reads like an instruction — the recipe's arrow *asks*, the item's *states* | the backlog's detail |
| **Q12** | does the **unbracketed** mouse position's message name its remedy? | the orchestrator | T2b's review: it names the cause (*"malformed mouse step"*) but not the remedy, while the bracketed case says *"the position must be (x,y)"* | one line |
| **Q13** | the kinds list is built on **every call**, and the footer asks it every frame | the orchestrator | T3b's review: a per-frame allocation on a render path that redraws the whole screen anyway | nothing yet |

## The three that wait on the chef

### Q8 — the free edit of a backlog item's fields

**The question.** The console shapes an item (its kind and a scope), dismisses it, and
creates one. It does **not** let the chef rewrite an item's title, body or tags. **Is
that a gap, or is it the right shape?**

| | |
|---|---|
| **(a) never** | the console manages the backlog; it does not rewrite an item's fields. The kind and the scope are the two the domain's shaping needs, and the rest is the chef's own words in the CLI |
| **(b) in a later round** | a five-field form, one prompt per field. It is a *form*, not a move — and the console's prompts are one text at a time |
| **(c) something narrower** | the title and the body, and not the scope or the tags |

**The orchestrator's recommendation: (a), for now.** The moves cover the chef's needs,
and the refusal T3a landed already names the remedy the chef actually wants — *shape it
into a small task first*. **If the chef wants (b) or (c), it is a slice, and it is not
this round.**

### Q9 — a backlog item's lists: absent, or `[]`?

**The question.** The domain's rule is *a list is a list: an empty one is `[]`, never
`null`* — and the **recipe** obeys it, on the writer as well as the reader. The
**backlog item** does not: `scope` and `tags` carry `omitempty`, so an item with no
scope has **no** `scope` key, and a reader that asks for the list gets nothing.

| | |
|---|---|
| **(a) keep it** | the item is a compact record, and an absent key reads the same as a null one to every reader |
| **(b) write the lists** | one rule for every record in the plan. Every item on disk changes shape the next time it is written. **The version does not move** — the fingerprint hashes the field **names**, and `BacklogItem.scope` is already one of them |

**The orchestrator's recommendation: (b).** The rule is the domain's own, and the
recipe's side of it was already fixed in T3a. Leaving the item's side out makes the
rule a rule of one record.

### Q10 — after a write, does the selection follow the item?

**The question.** The chef shapes an idea into a small task. The item **moves** —
from *Ideas* to *Small tasks* — and the cursor keeps its **index**, so the selection
lands on **another item**. The same is true of a sign-off on the Recipes tab, which
moves the recipe between the triage sections. It is pre-existing, and console-wide.

| | |
|---|---|
| **(a) the selection follows the item** | the chef's cursor stays on the thing they just acted on, wherever it went |
| **(b) the cursor keeps its index** | the console's list is a list, and a write does not move the cursor |

**The orchestrator's recommendation: (a).** A selection that points at another item
after a write is a lie about where the chef is — and it is the console's model, so the
fix belongs in the model's selection handling, not in a tab.

## The rules this file follows

| | |
|---|---|
| **A question opens** | when it stops being a detail and starts needing an answer. It goes here, with its options and the orchestrator's recommendation |
| **A question closes** | when it is answered. The answer moves to the **rulings table** in [tui/README.md](tui/README.md), and the entry here is struck through with the date |
| **Who answers** | the **chef** when the decision is about what the console is *for* — its surface, its vocabulary, the shape of the plan. The **orchestrator** when it is about how the console *works*, and the answer is a judgement the house already made |
| **The check** | the reviewer reads this file at every slice's close, and says in its report whether a question was answered, missed, or newly raised |
| **Nothing is lost** | a question that is not written here is a question that is lost. If it is in a brief, a finding or a review's note, and it needs an answer, it belongs here |
