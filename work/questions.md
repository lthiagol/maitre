# The open questions — what wait on the chef

> Each question: the options, the recommendation, and who answers it. **Ask the chef
> with the ask tool** when a question is ready — a question in a file is a question
> waiting to be read, not decided. The reviewer checks this file at every slice's
> close.

| # | Question | From | Who | Recommendation |
|--:|---|---|---|---|
| ~~Q8~~ | ~~the free edit of a backlog item's fields~~ | T3b | chef | **no, for now** — ruled 2026-10-31 |
| ~~Q9~~ | ~~a backlog item's `scope`/`tags`: absent, or `[]`?~~ | T3a | chef | **`[]`** — ruled 2026-10-31 |
| ~~Q10~~ | ~~after a write, does the selection follow the item?~~ | T3b | chef | **yes** — ruled 2026-10-31 |
| **Q11** | one marker or two for *the domain's answer*? | T3b | orchestrator | **two** — the recipe's arrow *asks*, the item's *states* |
| **Q12** | does the unbracketed position's message name its remedy? | T2b | orchestrator | **yes** — one line; each cause names its own remedy |
| **Q13** | the kinds list is built on every call | T3b | orchestrator | **later** — a per-frame allocation on a path that redraws the whole screen |
| **Q17** | the orchestrator's tab: how many tabs, and where does it sit? | T8 | chef | **8 tabs** — Home · Recipes · Backlog · Ready · Board · **Sous-Chef** · Log · Settings. The crew sits with the plan's tabs, before the records and the settings |
| **Q18** | does the Sous-Chef tab offer the chef a move, or only name it? | T8 | chef | **names it** — a blocked agent says what it needs and where the move lives. A move the console cannot perform is a promise it must not make |
| **Q19** | where does the console read the automation's state from? | T8 | chef | **the plan** — the orchestrator writes the crew's state into the plan and the console reads only the plan, as it does now. The cost is a format change: `DomainVersion` 7 → 8 |

## The rules

| | |
|---|---|
| **Concise** | a question is a question — the options, the recommendation, and what it blocks. Nothing a reader could reconstruct |
| **Ask** | when a question is ready, **ask the chef with the ask tool**. Do not leave it in this file to be answered by reading |
| **Recommend** | every question carries the recommendation, and the ask leads with it |
| **Close** | the answer moves to the **rulings table** in [tui/README.md](tui/README.md), and the entry is struck through with the date |
| **Check** | the reviewer reads this file at every slice's close, and says whether a question was **answered, missed or newly raised** |
| **Nothing is lost** | a question that is not written here is a question that is lost |
