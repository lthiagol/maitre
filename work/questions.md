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
| ~~Q11~~ | ~~one marker or two for *the domain's answer*?~~ | T3b | orchestrator | ~~**two** — the recipe's arrow *asks*, the item's *states*~~ — closed 2026-11-01 |
| ~~Q12~~ | ~~does the unbracketed position's message name its remedy?~~ | T2b | orchestrator | ~~**yes** — one line; each cause names its own remedy~~ — closed 2026-11-01 |
| ~~Q13~~ | ~~the kinds list is built on every call~~ | T3b | orchestrator | ~~**later** — a per-frame allocation on a path that redraws the whole screen~~ — closed 2026-11-01 |
| **Q25** | Q11, Q12 and Q13 have been open since T2b and T3b — the TUI round is over. Close them, or keep them open? | T3b, T2b | orchestrator | **close them** — Q11 and Q12 are answered by what the console now does; Q13 is a per-frame allocation, and the T7 count has the same shape. The chef rules the two that have a product answer |
| ~~Q28~~ | ~~where do the console's settings live — the view's preferences and the crew's configuration?~~ | T9 | orchestrator | ~~**in the plan's folder, in a new file** — `kitchen/settings.json`, owned by the domain. The format moves: `DomainVersion` 7 → 8~~ — ruled 2026-11-01 |

## The rules

| | |
|---|---|
| **Concise** | a question is a question — the options, the recommendation, and what it blocks. Nothing a reader could reconstruct |
| **Ask** | when a question is ready, **ask the chef with the ask tool**. Do not leave it in this file to be answered by reading |
| **Recommend** | every question carries the recommendation, and the ask leads with it |
| **Close** | the answer moves to the **rulings table** in [tui/README.md](tui/README.md), and the entry is struck through with the date |
| **Check** | the reviewer reads this file at every slice's close, and says whether a question was **answered, missed or newly raised** |
| **Nothing is lost** | a question that is not written here is a question that is lost |
