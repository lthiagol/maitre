package ui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
)

// TestALongRefusalIsVisible covers the refusal the console used to swallow.
//
// The version guard's refusal is long, and the summary line dropped the right side
// when the gap went negative — so the refusal was invisible at the default width. A
// refused write that shows no banner reads as nothing happened.
func TestALongRefusalIsVisible(t *testing.T) {
	m := newModel(t)
	m.width = 80
	m.height = 40
	m.refusal = "refused — this plan was written under domain version 2 and this tool speaks version 1: upgrade the tool, or read the plan without writing to it"

	out := plain(m)
	if !strings.Contains(out, "⛔") {
		t.Fatalf("the refusal has no marker:\n%s", out)
	}
	// The whole reason is on the screen. It wraps, so the text is read with the
	// line breaks flattened.
	if !strings.Contains(flat(out), flat(m.refusal)) {
		t.Fatalf("the refusal does not show the reason:\n%s", out)
	}
}

// TestAShortRefusalStaysOnTheSummaryLine covers the other half.
//
// A refusal that fits keeps the place it has always had: on the summary line, right
// aligned, in the warning style.
func TestAShortRefusalStaysOnTheSummaryLine(t *testing.T) {
	m := newModel(t)
	m.refusal = "refused — review.clean"

	raw := m.View().Content
	if !strings.Contains(raw, "⛔ refused — review.clean") {
		t.Fatalf("the summary line does not carry the refusal:\n%s", plain(m))
	}
	// The left text is still there, beside it.
	if !strings.Contains(flat(plain(m)), flat(m.summaryText())) {
		t.Fatalf("the summary text is gone:\n%s", plain(m))
	}
	// It is right aligned: the left text comes before it on the same row.
	line := summaryRow(t, m)
	if !strings.Contains(flat(line), flat(m.summaryText())) || !strings.Contains(line, "⛔") {
		t.Fatalf("the refusal is not on the summary row:\n%q", line)
	}
	// And it is in the warning style, not the accent style.
	warn := m.styWarn.Render("⛔ refused — review.clean")
	if !strings.Contains(raw, warn) {
		t.Fatalf("the refusal is not in the warning style")
	}
}

// TestTheRefusalBannerTakesRowsFromTheList covers the layout.
//
// The banner cannot be free: the screen is a fixed number of rows, and a banner that
// overflowed it would push the footer off the screen.
func TestTheRefusalBannerTakesRowsFromTheList(t *testing.T) {
	m := newModel(t)
	m.width = 80
	m.height = 40
	m.refusal = strings.Repeat("a long reason ", 40)

	rows := strings.Split(strings.TrimRight(m.View().Content, "\n"), "\n")
	if len(rows) != m.height {
		t.Fatalf("the console drew %d rows, want %d", len(rows), m.height)
	}
	// The banner is under the summary line, near the top.
	at := -1
	for i, r := range rows {
		if strings.Contains(r, "⛔") {
			at = i
			break
		}
	}
	if at < 0 {
		t.Fatalf("the refusal is nowhere on the screen:\n%s", strings.Join(rows[:5], "\n"))
	}
	if at > headerRows {
		t.Fatalf("the refusal is on row %d, want it under the header", at)
	}
}

// TestAWrappedRefusalStillNamesTheGate covers the type assertion.
//
// The banner prints the gate's *names*, not the gate's prose: the prose can carry an
// agent-written reviewer name, and the summary line does not neutralise it. A
// refusal that arrives wrapped is still a refusal, so it has to be recognised.
func TestAWrappedRefusalStillNamesTheGate(t *testing.T) {
	m := newModel(t)
	wrapped := fmt.Errorf("while closing: %w", &core.GateError{
		Event: core.EventClose,
		Gates: []string{"review.clean", "audit.documented (the audit trail has no documentation)"},
	})
	m.refusal = refusalText(wrapped)
	if !strings.Contains(m.refusal, "review.clean") {
		t.Errorf("the refusal should name the gate: %q", m.refusal)
	}
	if strings.Contains(m.refusal, "the audit trail has no documentation") {
		t.Errorf("the refusal should not carry the gate's prose: %q", m.refusal)
	}
}

// flat is the text with the escape sequences, the line breaks and the runs of spaces
// gone, so a wrapped message can be read as one string.
func flat(s string) string {
	sp := strings.Join(strings.Fields(StripANSI(s)), " ")
	return sp
}

// summaryRow is the row of the rendered console that carries the summary.
func summaryRow(t *testing.T, m *Model) string {
	t.Helper()
	for _, l := range strings.Split(m.View().Content, "\n") {
		if strings.Contains(flat(l), flat(m.summaryText())) {
			return l
		}
	}
	t.Fatalf("no row carries the summary text")
	return ""
}
