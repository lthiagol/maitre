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
| **Q19** | where does the console read the automation's state from? | T8 | chef | **herdr** — the console asks the runtime, and finds the crew by the agents' cwd. The plan carries no crew state, so `DomainVersion` stays 7 |
| **Q20** | does starting the automation leave a trace in the plan's journal? | T8 | chef | **no** — the journal records what happened to the plan; the crew's state is runtime state. A trace would mean a new operation and a format change |
| **Q21** | what does the startup task say? | T8 | chef | **the job and the path** — "you are the orchestrator of the plan at <path>". The orchestrator reads the plan and decides what is needed |

## The rules

| | |
|---|---|
| **Concise** | a question is a question — the options, the recommendation, and what it blocks. Nothing a reader could reconstruct |
| **Ask** | when a question is ready, **ask the chef with the ask tool**. Do not leave it in this file to be answered by reading |
| **Recommend** | every question carries the recommendation, and the ask leads with it |
| **Close** | the answer moves to the **rulings table** in [tui/README.md](tui/README.md), and the entry is struck through with the date |
| **Check** | the reviewer reads this file at every slice's close, and says whether a question was **answered, missed or newly raised** |
| **Nothing is lost** | a question that is not written here is a question that is lost |
