Runner: T3b done. Commits: 4 (1297cdd, 35688a6, 4799b3b, 7de9a16) + the state row (907bd09). Per-AC evidence: AC-01..AC-07 stamped. Pre-notify gate: pass.

The moves: the domain answers them in ONE place - core.BacklogMoveOf (internal/core/backlog.go), a **view**: the // view marker, a registration in core.DerivedViews (name "backlog_move"), and a row in docs/03 § 8. It answers two things about the selected item: Primary (shape | promote | none) and Dismissable. It reads core.PromotionOf, so the console can never offer a promotion the domain refuses. The console asks it in internal/ui/backlog.go (backlogPrimary), in the footer (primaryActionLabel) and for the dismiss key (canDismissSelected) - the console has no "is it a small task?" rule anywhere.

Shape it: space on an open idea opens the console's prompt, which asks two things - feature or bug (the chef picks; the console does not) and the scope. The scope is *added* to the item's scope and the kind becomes a small task. The id does NOT change: BL01 stays BL01. It is an UPDATE, not a promotion: no recipe is created and the item stays open. The domain's UpdateBacklogItem does it; the console asks. An empty scope is not written and the prompt stays open.

Dismiss: `d` on the Backlog tab asks for the reason in the console's own text prompt; the domain's DismissBacklogItem appends it to the item's body ("Dismissed: later"), so the next reader knows why. An empty reason is not written. The footer names `d dismiss` only while the domain says the item may still be dismissed.

Create: `a` on the Backlog tab asks for one thing - the title, which is what the domain needs and no more. The item starts as an **idea** (the kind a hunch starts as) and the chef shapes it after. A backlog item is not a recipe, so ADR-0003 is not broken (Q6). The console writes a backlog item and nothing else: no recipe, and nothing outside kitchen/backlog/ (the journal records the creation - the plan's own event log).

The footer: on an open idea - "space shape it"; on an open small task - "space promote to recipe"; on a promoted or dismissed item - no space at all, and no `d dismiss` either. The `a add item` key is on the Backlog tab always.

The domain's answer: a **view**, not a declared helper - core.BacklogMoveOf, registered as "backlog_move" with the closed set shape|promote|none, and with a row in docs/03 § 8. The table's name comes from the runtime (funcName). Why a view: it computes something the plan never stores, and it is the same shape as core.ChefMoveOf and core.PromotionOf, which the console already asks.

DomainVersion: 7. No format, no field, no directory. The tab count is 7 and the names did not change. No dependency added: internal/core 0, internal/cli 0, internal/ui 28 (ceiling 30).

The goldens: unchanged. `make goldens` leaves git diff empty - the goldens hold the Recipes tab, and this slice changed the Backlog tab and the prompt framework.

The gate: `make check` green, `make check-race` green, the tree clean.

Evidence, per criterion:
- AC-01 - the footer on each state (the console asked): idea -> "space shape it"; small task -> "space promote to recipe"; promoted -> no space. MUTATED: the console made to decide by kind (if kind == feature|bug -> "promote to recipe") and the suite went RED - TestTheFooterNamesTheMoveTheDomainGives, TestTheFooterAsksTheDomainForTheMove and TestTheFooterAsksTheDomainAboutPromotability all failed. Reverted.
- AC-02 - BL01 shaped from the console: kind feature, scope ["the evidence lane"], id BL01, status open, and the recipe count unchanged.
- AC-03 - BL03 dismissed from the console: status dismissed, body carries "Dismissed: later".
- AC-04 - BL07 written down from the console: kind idea, status open, in kitchen/backlog/; 7 recipes before and after; the only changed files are kitchen/backlog/ and kitchen/activity.json.
- AC-05 - space on an idea shapes it (no refusal banner, nothing written until the chef answers); `maitre-cli backlog promote BL01` is STILL refused, with the reason, and the journal records the refusal as an agent's attempt (seq 59, kind "refusal").
- AC-06 - `make goldens` leaves git diff empty.
- AC-07 - work/tui/README.md: T3b is closed with its 4 commits; the next row that is not closed is T4, the Home tab.

Notes for the next session:
- The plan (work/tui/plan.md, the orchestrator's file) still says T3b lands "dismiss . create . edit". The brief excluded the free edit of every field, so the plan's T3b row and its § 4 "Edit" row need correcting - I did not touch work/ beyond the state row.
- docs/05-human-interface.md: the footer's line about an idea ("the footer offers none") and the `space` row were corrected; the `d` and `a` keys were added to the interaction table. The "four jobs" table still lists "edit" among the backlog's moves - that is the job, not the built surface, and the free edit is not built.
- The console's keys on the Backlog tab: space (the item's move), d (dismiss, with a reason), a (add an item).
- The console's scripted driver types ONE character per step, so `maitre --keys` can only type single letters: a multi-letter word is refused as an unknown step. That is the driver's old rule, unchanged here.
