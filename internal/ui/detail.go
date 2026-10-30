package ui

import (
	"fmt"
	"strings"

	"github.com/lthiagol/maitre/internal/core"
)

// detailLines renders the detail pane for the current tab's selection.
//
// The pane follows the tab: a recipe on Recipes, Path and Board, a backlog item
// on Backlog, a journal entry on Log. One pane, whatever is selected.
func (m *Model) detailLines(width int) []string {
	switch m.tab {
	case TabRecipes, TabPath, TabBoard:
		rec, ok := m.currentRecipe()
		if !ok {
			return []string{m.styDim.Render("  (nothing selected)")}
		}
		return m.recipeDetail(rec, width)
	case TabBacklog:
		return m.backlogDetail(width)
	case TabLog:
		return m.logDetail(width)
	case TabSousChef:
		// There is nothing to open: the Sous-Chef is not a list of recipes. Saying so is
		// better than an empty pane, which looks like a broken console.
		return []string{m.styDim.Render("  Nothing to open here — the Sous-Chef is not a list of recipes.")}
	}
	return []string{m.styDim.Render("  (nothing to show here)")}
}

// backlogDetail renders the selected backlog item.
//
// The item is not a recipe, so it has its own rendering — but it is the console's
// **one** detail pane: the same geometry, the same scroll, the same `esc`.
//
// Every field it shows is a stored fact the domain reads: the title, the body, the
// scope, the kind, the status, the tags, the dates, and where it went when it was
// promoted.
func (m *Model) backlogDetail(width int) []string {
	item, ok := m.selectedBacklogItem()
	if !ok {
		return []string{m.styDim.Render("  (nothing selected)")}
	}
	out := []string{" " + m.styTitle.Render(item.ID+" · "+item.Title)}
	out = append(out, " "+m.styDim.Render(item.Kind+" · "+item.Status))
	out = append(out, " "+m.styFaint.Render(fmt.Sprintf("created %s · updated %s", item.Created, item.Updated)))
	out = append(out, "")
	if item.Body != "" {
		out = append(out, m.stySec.Render("Why"))
		out = append(out, indent(renderMarkdown(item.Body, width-4), "  "))
		out = append(out, "")
	}
	if len(item.Scope) > 0 {
		out = append(out, m.stySec.Render("Scope"))
		for _, s := range item.Scope {
			out = append(out, "  "+m.styGood.Render("in   ")+s)
		}
		out = append(out, "")
	}
	if len(item.Tags) > 0 {
		out = append(out, m.stySec.Render("Tags"))
		out = append(out, "  "+strings.Join(item.Tags, ", "))
		out = append(out, "")
	}
	if item.PromotedTo != "" {
		out = append(out, m.styAcc.Render("Promoted to "+item.PromotedTo))
	}
	return out
}

// logDetail renders the selected journal entry.
//
// It reads the same slice the list renders, so the detail always shows the event the
// highlighted row holds.
func (m *Model) logDetail(width int) []string {
	events := m.logWindow()
	if m.selLog < 0 || m.selLog >= len(events) {
		return []string{m.styDim.Render("  (nothing selected)")}
	}
	e := events[m.selLog]
	out := []string{" " + m.styTitle.Render(fmt.Sprintf("%s %s", e.Kind, e.Recipe))}
	out = append(out, " "+m.styDim.Render(e.At))
	out = append(out, " "+m.styFaint.Render(fmt.Sprintf("%s · %s", e.Actor, e.Role)))
	out = append(out, "")
	if len(e.Data) > 0 {
		out = append(out, m.stySec.Render("Data"))
		out = append(out, indent(renderMarkdown("```json\n"+string(e.Data)+"\n```", width-4), "  "))
		out = append(out, "")
	}
	return out
}

// detailRows flattens the detail into terminal rows.
//
// A markdown paragraph is one entry but several rows, so counting entries is not
// counting rows — and a pane that overflows its height pushes the footer off the
// screen. Everything that measures the detail measures rows.
func detailRows(lines []string) []string {
	rows := make([]string, 0, len(lines))
	for _, l := range lines {
		rows = append(rows, strings.Split(l, "\n")...)
	}
	return rows
}

// viewDetail renders the recipe detail full screen.
func (m *Model) viewDetail(width int) string {
	rec, ok := m.currentRecipe()
	if !ok {
		return m.styDim.Render("  (nothing selected)")
	}
	rows := m.detailScrollRows(detailRows(m.recipeDetail(rec, width)))
	return strings.Join(rows, "\n")
}

// detailScrollRows is the detail's rows, from the scroll's offset.
//
// The offset is applied **here**, once. Two places that applied it were two answers
// to *where in the recipe is the chef*, and the layouts disagreed: the split layout
// moved the offset and drew the same frame.
func (m *Model) detailScrollRows(rows []string) []string {
	if m.detailScroll <= 0 || len(rows) == 0 {
		return rows
	}
	if m.detailScroll >= len(rows) {
		m.detailScroll = len(rows) - 1
	}
	return rows[m.detailScroll:]
}

// detailPaneWidth is the width the detail pane draws with.
//
// The pane is beside the list in the split layout, and the list takes a column of
// it — so the pane is narrower there. A narrower pane wraps the same text into
// **more** rows, so the pane's width is also the width the scroll's bound is
// computed against: one answer for the detail's rows, and the bound cannot stop
// early against a row count the pane never had.
func (m *Model) detailPaneWidth() int {
	if m.mode == ModeDetail {
		return m.width - 4
	}
	if m.width >= splitThreshold {
		return m.width - m.listColumn() - 3
	}
	// The single layout draws no pane. The width only answers for the bound, and
	// the full screen is the honest answer for it.
	return m.width - 4
}

// recipeDetail renders one recipe.
func (m *Model) recipeDetail(r *core.Recipe, width int) []string {
	var out []string

	out = append(out, " "+m.styTitle.Render(r.ID+" · "+r.Title))
	// One word per fact. The state word is the plain one — `closed`, `executing` —
	// and the lane is a *section* the board groups by, not a second name for where
	// the recipe is. A blocked recipe says so on its own line, below, with the
	// reason; nothing here needs the kitchen word for it.
	meta := []string{string(r.Lifecycle)}
	s := core.EvidenceSummaryOf(*r)
	if s.Criteria > 0 {
		meta = append(meta, fmt.Sprintf("%d/%d proven", s.Proven, s.Criteria))
	}
	if rs := core.ReviewStateOf(*r); rs != core.ReviewNone {
		meta = append(meta, rs.Label())
	}
	out = append(out, " "+m.styDim.Render(strings.Join(meta, " · ")))
	// Why it needs the chef, when it does. "awaiting an independent review" says
	// what to do next; "awaiting tasting" does not.
	if core.TriageOf(*r) == core.TriageNeedsYou {
		if reason := core.ReadinessOf(*r).Reason; reason != "" {
			out = append(out, " "+m.styWarn.Render("→ "+reason))
		}
	}
	out = append(out, "")

	if r.Intent.Outcome != "" {
		out = append(out, m.stySec.Render("Outcome"))
		out = append(out, indent(renderMarkdown(r.Intent.Outcome, width-4), "  "))
		out = append(out, "")
	}
	if r.Intent.Problem != "" {
		out = append(out, m.stySec.Render("Problem"))
		out = append(out, indent(renderMarkdown(r.Intent.Problem, width-4), "  "))
		out = append(out, "")
	}
	if len(r.Scope.InScope) > 0 || len(r.Scope.OutOfScope) > 0 {
		out = append(out, m.stySec.Render("Scope"))
		for _, s := range r.Scope.InScope {
			out = append(out, "  "+m.styGood.Render("in   ")+s)
		}
		for _, s := range r.Scope.OutOfScope {
			out = append(out, "  "+m.styDim.Render("out  ")+s)
		}
		out = append(out, "")
	}

	if len(r.Acceptance) > 0 {
		out = append(out, m.stySec.Render("Acceptance criteria"))
		for _, c := range r.Acceptance {
			out = append(out, "  "+padRight(c.ID, 8)+m.criterionStyle(string(c.Status())).Render(padRight(string(c.Status()), 9))+c.Description)
			if c.Verification != "" {
				out = append(out, "          "+m.styFaint.Render("verify: "+c.Verification))
			}
			for _, e := range c.Evidence {
				mark := ""
				if e.SignedOff {
					mark = "  (signed off)"
				} else if e.TimedOut {
					mark = "  (timed out)"
				}
				out = append(out, "          "+m.styFaint.Render(fmt.Sprintf("%s → exit %d, matched %d%s", e.Command, e.ExitCode, e.Matched, mark)))
			}
		}
		out = append(out, "")
	}

	if len(r.Guidance) > 0 {
		out = append(out, m.stySec.Render("Method"))
		for _, st := range r.Guidance {
			out = append(out, "  "+padRight(st.ID, 5)+m.styDim.Render(padRight(st.Status, 9))+st.Action)
			if st.DoneWhen != "" {
				out = append(out, "        "+m.styFaint.Render("done when: "+st.DoneWhen))
			}
			if len(st.Covers) > 0 {
				out = append(out, "        "+m.styFaint.Render("covers: "+strings.Join(st.Covers, ", ")))
			}
			if st.Deviation != "" {
				out = append(out, "        "+m.styWarn.Render("deviated: "+st.Deviation))
			}
		}
		out = append(out, "")
	}

	// The comments are grouped by the finding they belong to, and the ones that belong
	// to no finding get their own section. The recipe is the channel between the
	// reviewer and the fixer, so the thread has to be readable here.
	byFinding := map[string][]core.Comment{}
	var general []core.Comment
	for _, c := range r.Review.Comments {
		if c.Finding == "" {
			general = append(general, c)
			continue
		}
		byFinding[c.Finding] = append(byFinding[c.Finding], c)
	}

	if len(r.Review.Findings) > 0 {
		out = append(out, m.stySec.Render("Tasting notes"))
		selected := ""
		if note, ok := m.selectedNote(); ok {
			selected = note.ID
		}
		for _, f := range r.Review.Findings {
			out = append(out, m.noteLine(f, f.ID == selected, width))
			out = append(out, "        "+f.Description)
			for _, c := range byFinding[f.ID] {
				out = append(out, "          "+m.styAcc.Render(c.Author)+": "+c.Body)
			}
		}
		out = append(out, "")
	}

	if len(general) > 0 {
		out = append(out, m.stySec.Render("Review notes"))
		for _, c := range general {
			out = append(out, "  "+m.styAcc.Render(c.Author)+": "+c.Body)
		}
		out = append(out, "")
	}

	if r.Overlays.Blocked != nil {
		out = append(out, m.styWarn.Render("On hold: "+r.Overlays.Blocked.Reason))
		out = append(out, "")
	}
	if r.Overlays.Deferred != nil {
		out = append(out, m.styDim.Render("Deferred: "+r.Overlays.Deferred.Reason))
		out = append(out, "")
	}
	if r.Approval != nil {
		out = append(out, m.styAcc.Render("Signed off by "+r.Approval.By+" at "+r.Approval.At))
		out = append(out, "")
	}
	if len(r.Audit.Deviations) > 0 {
		out = append(out, m.stySec.Render("Deviations"))
		for _, d := range r.Audit.Deviations {
			out = append(out, "  "+d.Note)
		}
		out = append(out, "")
	}

	return out
}

// noteLine renders one tasting note's header line.
//
// The note under the cursor wears the console's own selected-row style — the same
// style the recipe rows wear, so "selected" means one thing on this screen.
func (m *Model) noteLine(f core.Finding, selected bool, width int) string {
	if selected {
		plain := "  " + padRight(f.ID, 6) + padRight(f.Severity, 9) + padRight(f.Category, 13) + f.Status
		return m.stySelected.Render(padRight(plain, width))
	}
	cls := m.styDim
	switch f.Status {
	case "open":
		cls = m.styBad
	case "fixed":
		cls = m.styGood
	}
	return "  " + padRight(f.ID, 6) + cls.Render(padRight(f.Severity, 9)) + padRight(f.Category, 13) + m.styDim.Render(f.Status)
}

// indent prefixes every line of s.
func indent(s, prefix string) string {
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = prefix + lines[i]
	}
	return strings.Join(lines, "\n")
}
