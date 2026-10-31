package core

// PlanStats is the plan's numbers, as the Home tab shows them.
//
// It is a **view**: computed from the plan on every read, never stored. A number the
// console counts itself is a number the CLI would count differently, and two answers
// to one question is the house's standing failure.
type PlanStats struct {
	// Recipes counts the recipes by what needs the chef.
	Recipes RecipeStats
	// Evidence counts the evidence on the plan's criteria.
	Evidence EvidenceSummary
	// Tasting counts what the tasting still has open.
	Tasting TastingStats
	// Backlog counts the backlog by the group the item's kind puts it in.
	Backlog BacklogStats
	// Automation is the automation's vital signs.
	Automation AutomationStats
}

// RecipeStats counts the recipes by triage section.
//
// The sections are the triage's own — the same answer the Recipes tab's groups come
// from — so the tab's numbers and the board's rows cannot disagree about what needs
// the chef.
//
// **The sections are a partition**: every recipe is filed under exactly one, and
// the counts add up to the total. A row that names some of them and not the rest
// misleads — the parts beside it do not make the number.
type RecipeStats struct {
	// Total is every recipe in the plan.
	Total int
	// NeedsYou is the recipes filed under `Needs you`.
	NeedsYou int
	// OnTheLine is the recipes filed under `On the line`.
	OnTheLine int
	// BeingShaped is the recipes filed under `Being shaped`.
	BeingShaped int
	// OnHold is the recipes filed under `On hold`.
	OnHold int
	// Served is the recipes filed under `Served`.
	Served int
}

// TastingStats counts what the tasting still has open.
//
// The two numbers are not the same, and the difference is the point:
//
//   - **open notes** is every open finding, of any phase — the notes the chef
//     triages, and the ones the close gate refuses a recipe over.
//   - **open findings** is the open notes the **external** tasting filed: somebody
//     else's objection, not the cook's own.
type TastingStats struct {
	// OpenNotes is every open finding, of any phase.
	OpenNotes int
	// OpenFindings is the open findings the external tasting filed.
	OpenFindings int
}

// BacklogStats counts the backlog by its two groups.
type BacklogStats struct {
	// Ideas is the items of kind `idea`.
	Ideas int
	// SmallTasks is the items of kind `feature` or `bug`.
	SmallTasks int
}

// AutomationStats is the automation's vital signs.
//
// The chef can already see whether the automation is running on the same rules as the
// plan: the version the plan was last written under, the version this build speaks,
// and whether they agree.
type AutomationStats struct {
	// PlanVersion is the version of the rules the plan was last written under.
	PlanVersion int
	// BuildVersion is the version of the rules this build speaks.
	BuildVersion int
	// VersionsAgree says whether the two are the same.
	VersionsAgree bool
	// LastEvent is the newest event in the journal, and nil when the plan has none.
	LastEvent *JournalEntry
}

// PlanStatsOf computes the plan's numbers.
//
// One answer for the whole tab: the recipes, the evidence, the tasting, the backlog
// and the automation. Every count is read from the domain's own view — the triage's
// sections, the evidence summary, the open findings, the backlog's groups — so the
// numbers the console shows and the numbers the CLI would show are one answer.
//
// The plan's name and its path are **not** here: they are stored facts, read from
// `core.Project` and the plan's root, and a view holds what is computed.
// view
func PlanStatsOf(recipes []Recipe, items []BacklogItem, journal Journal) PlanStats {
	stats := PlanStats{
		Recipes: RecipeStats{Total: len(recipes)},
		Automation: AutomationStats{
			PlanVersion:  PlanDomainVersion(journal),
			BuildVersion: DomainVersion,
		},
	}
	stats.Automation.VersionsAgree = stats.Automation.PlanVersion == stats.Automation.BuildVersion

	for _, sec := range TriageSections(recipes) {
		switch sec.Name {
		case TriageNeedsYou:
			stats.Recipes.NeedsYou = len(sec.Items)
		case TriageOnTheLine:
			stats.Recipes.OnTheLine = len(sec.Items)
		case TriageBeingShaped:
			stats.Recipes.BeingShaped = len(sec.Items)
		case TriageOnHold:
			stats.Recipes.OnHold = len(sec.Items)
		case TriageServed:
			stats.Recipes.Served = len(sec.Items)
		}
	}

	for _, r := range recipes {
		e := EvidenceSummaryOf(r)
		stats.Evidence.Criteria += e.Criteria
		stats.Evidence.Proven += e.Proven
		stats.Evidence.Failed += e.Failed
		stats.Evidence.Rejected += e.Rejected
		stats.Evidence.Unproven += e.Unproven
		stats.Evidence.Pending += e.Pending

		for _, f := range OpenFindingsOf(&r) {
			stats.Tasting.OpenNotes++
			if f.Phase == FindingPhaseExternal {
				stats.Tasting.OpenFindings++
			}
		}
	}

	ideas, tasks := BacklogGroups(items)
	stats.Backlog.Ideas = len(ideas)
	stats.Backlog.SmallTasks = len(tasks)

	if len(journal.Events) > 0 {
		last := journal.Events[len(journal.Events)-1]
		stats.Automation.LastEvent = &last
	}
	return stats
}
