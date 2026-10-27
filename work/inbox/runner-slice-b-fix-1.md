# Runner — slice B, fix cycle 1

**Commits:** 1.

| Commit | What |
|---|---|
| `8c394cd` | **F-72** — the spec's rule sentence agrees with its table |

## What changed

`docs/03-domain-model.md` § *The CPD ↔ tool mapping*. One sentence, and the
citation that carried the stale wording.

**Before**

> The rule behind the split: the **artifact and the roles** get themed names;
> everything an agent reasons about mechanically stays descriptive. See
> [ADR-0005](decisions/0005-vocabulary.md).

**After**

> The rule behind the split: the kitchen names the **artifact** and the
> **place**; the **agents' jobs** and everything an agent reasons about
> mechanically stay plain. See
> [docs/07-vocabulary.md § The rule](../07-vocabulary.md#the-rule).

The sentence is now true of the table as it stands: the Themed column holds the
artifact (recipe · special · prep note · course · menu) and the tool's own
names; the Kept-plain column holds the machine's records, the **agents' jobs**
(runner · reviewer · orchestrator) and the machine's facts.

## Notes for the reviewer

1. **The citation moved with the sentence.** The rule now points at the
   **decision of record** — `docs/07-vocabulary.md` § *The rule* — and not at
   ADR-0005. That is deliberate: ADR-0005 is the source of the stale wording
   (it keeps the roles in the themed column), and pointing a live sentence at a
   dated record that contradicts it is how the contradiction came back. The ADR
   itself was **not** touched — it is a dated record, and slice E writes ADR
   0009 to supersede it.

2. **ADR-0005 is now stale.** It still says the roles are themed vocabulary and
   still names `line cook`, `taster` and `expo`. Slice E's ADR-0009 supersedes
   it; this cycle left it alone, as instructed.

3. **One residual difference between the decision of record and `docs/03`'s
   table — reported, not fixed.** `docs/07`'s layer table puts **chef** and
   **sous-chef** in *The jobs | plain*, while `docs/03`'s table keeps
   `chef · maître d' · sous chef` in the **Themed** column. The two docs still
   differ on the cast's placement. The ruling was that the table stands, so the
   table was not touched, and the new sentence says nothing about the cast
   either way — it is silent rather than wrong. If the chef wants the cast
   settled, it belongs with slice E's ADR-0009 and the glossary sweep.

4. **No new tripwire**, as instructed: the deny-list wire already fires when a
   role is renamed in the code and the spec is not renamed with it.

5. **No code changes.** `internal/core`, `internal/cli` and `internal/ui` are
   untouched.

## The gate

```
make lint && gofmt -l . && make test && git status --short
```

Green: `go vet` clean, `gofmt -l .` empty, core · cli · ui all pass, and the
tree clean apart from the orchestrator's own `work/` files.
