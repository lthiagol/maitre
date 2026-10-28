package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"

	"github.com/lthiagol/maitre/internal/brand"
)

// TestTheHeaderCarriesTheMarkAndStaysOneLine covers the console's header.
//
// The header is one line, and the whole reason the compact form exists is that it
// fits in one line. A header that wrapped would push the tabs down and break the
// layout.
func TestTheHeaderCarriesTheMarkAndStaysOneLine(t *testing.T) {
	m := newModel(t)
	header := m.header()

	if !strings.Contains(header, brand.Compact) {
		t.Errorf("the header should carry the mark:\n%s", header)
	}
	if !strings.Contains(header, "maitre") {
		t.Errorf("the header should carry the product name:\n%s", header)
	}
	if strings.Contains(header, "\n") {
		t.Errorf("the header should be one line:\n%s", header)
	}
	if w := lipgloss.Width(header); w > m.width {
		t.Errorf("the header is %d columns wide and the terminal is %d", w, m.width)
	}
}

// TestTheSousChefDrawsTheMarkThatFits covers the width rule.
//
// A drawing that does not fit wraps, and a wrapped row breaks the layout. So the
// column chooses the form: the full form when it fits, the half form when only it
// fits, and no art when neither does.
func TestTheSousChefDrawsTheMarkThatFits(t *testing.T) {
	full := brand.Width(brand.Full)
	half := brand.Width(brand.Half)

	for _, c := range []struct {
		name   string
		column int
		want   []string
	}{
		{"the column fits the full form", full + 1, brand.Full},
		{"the column fits the half form only", half + 1, brand.Half},
		{"the column fits the full form at the wide end", 200, brand.Full},
		{"the column fits neither", half, nil},
		{"a very narrow column", 4, nil},
	} {
		got := markFor(c.column)
		if len(got) != len(c.want) {
			t.Errorf("%s: the mark has %d lines, want %d", c.name, len(got), len(c.want))
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("%s: line %d is %q, want %q", c.name, i+1, got[i], c.want[i])
			}
		}
	}
}

// TestTheSousChefTabShowsTheMark covers the tab.
//
// The art goes at the top, above the section — the tab's content is a list of
// lines, and the drawing is the first of them.
func TestTheSousChefTabShowsTheMark(t *testing.T) {
	m := newModel(t)
	switchToTab(t, m, TabSousChef)
	out := plain(m)

	if !strings.Contains(out, " "+brand.Full[0]) {
		t.Errorf("the tab should draw the mark at the top:\n%s", out)
	}
	art := strings.Index(out, brand.Full[0])
	section := strings.Index(out, "SOUS-CHEF")
	if art < 0 || section < 0 || art > section {
		t.Errorf("the mark should be above the section:\n%s", out)
	}
}

// TestTheSousChefArtFitsTheColumn covers the layout.
//
// The drawing must not exceed the column it is in — a row that wraps breaks the
// layout.
func TestTheSousChefArtFitsTheColumn(t *testing.T) {
	m := newModel(t)
	switchToTab(t, m, TabSousChef)

	column := m.listColumn()
	for _, line := range markFor(column) {
		if w := lipgloss.Width(line) + 1; w > column {
			t.Errorf("the mark's line is %d columns wide and the column is %d", w, column)
		}
	}
}

// TestTheSousChefKeepsItsContentAtEveryHeight covers the height rule.
//
// The decoration yields to the content, in both dimensions. The tab's content must
// be whole at every height — the last event is the last thing that happened, and a
// face is worthless if it costs the chef that — and the mark is drawn only at the
// heights where the content and the mark both fit.
//
// The measure is the chef's screen: 80 columns, the classic default. The mark
// costs six rows, so a 24-row terminal has no room for it.
func TestTheSousChefKeepsItsContentAtEveryHeight(t *testing.T) {
	for _, c := range []struct {
		height  int
		wantArt bool
	}{
		{20, false},
		{24, false},
		{26, true},
		{30, true},
		{40, true},
	} {
		m := newModel(t)
		m.width = 80
		m.height = c.height
		switchToTab(t, m, TabSousChef)
		out := plain(m)

		if !strings.Contains(out, "LAST EVENT") {
			t.Errorf("height %d: the tab should show its last event:\n%s", c.height, out)
		}
		content := m.sousChefContent()
		last := strings.TrimSpace(StripANSI(content[len(content)-1].text))
		if !strings.Contains(out, last) {
			t.Errorf("height %d: the tab's last line should be on the screen — %q is not:\n%s", c.height, last, out)
		}
		if hasArt := strings.Contains(out, brand.Full[0]); hasArt != c.wantArt {
			t.Errorf("height %d: the mark is drawn: %v, want %v:\n%s", c.height, hasArt, c.wantArt, out)
		}
	}
}

// TestTheSousChefStillSaysItIsNotBuilt covers the face.
//
// The tab is a placeholder, and a placeholder that looks broken is worse than no
// tab at all. A mark does not change what the tab says.
func TestTheSousChefStillSaysItIsNotBuilt(t *testing.T) {
	m := newModel(t)
	switchToTab(t, m, TabSousChef)
	out := plain(m)

	for _, want := range []string{"Not built yet.", "VITAL SIGNS", "plan written under", "this build speaks", "LAST EVENT"} {
		if !strings.Contains(out, want) {
			t.Errorf("the tab should still say %q:\n%s", want, out)
		}
	}
}
