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
