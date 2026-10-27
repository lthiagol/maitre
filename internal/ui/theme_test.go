package ui

import (
	"image/color"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// The published palettes, slot by slot: bg, ink, dim, faint, accent, good, warn, bad,
// selBG, rule. These are the values from the official palettes — not adjusted by eye.
var publishedPalettes = map[string][]string{
	"Kitchen": {"235", "252", "245", "240", "39", "42", "214", "203", "237", "238"},
	"Catppuccin Latte": {
		"#eff1f5", "#4c4f69", "#6c6f85", "#9ca0b0", "#8839ef",
		"#40a02b", "#df8e1d", "#d20f39", "#ccd0da", "#bcc0cc",
	},
	"Catppuccin Frappé": {
		"#303446", "#c6d0f5", "#a5adce", "#737994", "#ca9ee6",
		"#a6d189", "#e5c890", "#e78284", "#414559", "#51576d",
	},
	"Catppuccin Macchiato": {
		"#24273a", "#cad3f5", "#a5adcb", "#6e738d", "#c6a0f6",
		"#a6da95", "#eed49f", "#ed8796", "#363a4f", "#494d64",
	},
	"Catppuccin Mocha": {
		"#1e1e2e", "#cdd6f4", "#a6adc8", "#6c7086", "#cba6f7",
		"#a6e3a1", "#f9e2af", "#f38ba8", "#313244", "#45475a",
	},
	"Dracula": {
		"#282a36", "#f8f8f2", "#6272a4", "#6272a4", "#bd93f9",
		"#50fa7b", "#ffb86c", "#ff5555", "#44475a", "#6272a4",
	},
	"Alucard": {
		"#fffbeb", "#1f1f1f", "#6c664b", "#6c664b", "#644ac9",
		"#14710a", "#a34d14", "#cb3a2a", "#cfcfde", "#6c664b",
	},
}

// TestThePalettesAreThePublishedOnes checks every theme's ten slots.
//
// The values are the published ones from the official palettes. A theme that is
// adjusted by eye is a theme nobody else has.
func TestThePalettesAreThePublishedOnes(t *testing.T) {
	list := palettes()
	if len(list) != len(publishedPalettes) {
		t.Fatalf("there are %d themes, want %d", len(list), len(publishedPalettes))
	}
	for _, p := range list {
		want, ok := publishedPalettes[p.name]
		if !ok {
			t.Errorf("the theme %q is not one of the published ones", p.name)
			continue
		}
		got := []string{
			p.colBg, p.colInk, p.colDim, p.colFaint, p.colAccent,
			p.colGood, p.colWarn, p.colBad, p.colSelBG, p.colRule,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s colours = %v, want %v", p.name, got, want)
		}
	}
}

// TestKitchenIsTheDefault checks that a console nobody told otherwise looks as it did.
//
// The golden snapshots are the real net for that; this says which palette they are.
func TestKitchenIsTheDefault(t *testing.T) {
	m := newModel(t)
	if m.theme.name != "Kitchen" {
		t.Fatalf("the default theme is %q, want Kitchen", m.theme.name)
	}
}

// TestTheThemeRepaintsTheConsole covers the change the chef makes.
//
// A theme that changed one corner of the screen would be a theme half applied.
func TestTheThemeRepaintsTheConsole(t *testing.T) {
	m := newModel(t)
	before := m.View().Content

	latte := paletteNamed(t, "Catppuccin Latte")
	m.theme = newTheme(latte)

	after := m.View().Content
	if before == after {
		t.Fatal("the screen did not change when the theme did")
	}
	esc := lipgloss.NewStyle().Foreground(lipgloss.Color(latte.colAccent)).Render("x")
	seq := esc[:strings.Index(esc, "x")]
	if !strings.Contains(after, seq) {
		t.Fatalf("the new theme's accent colour is not on the screen:\n%q", seq)
	}
}

// TestALightThemeIsLight checks the console's background.
//
// The console sets its background explicitly, so a light theme's ground has to be the
// one it sets — otherwise the text is tuned for a ground it is not drawn on.
func TestALightThemeIsLight(t *testing.T) {
	m := newModel(t)
	latte := paletteNamed(t, "Catppuccin Latte")
	m.theme = newTheme(latte)

	v := m.View()
	if v.BackgroundColor == nil {
		t.Fatal("the console sets no background")
	}
	if !sameColor(v.BackgroundColor, lipgloss.Color(latte.colBg)) {
		t.Fatalf("the background is %v, want Latte's %q", v.BackgroundColor, latte.colBg)
	}
	if m.theme.colInk != latte.colInk {
		t.Fatalf("the ink is %q, want Latte's %q", m.theme.colInk, latte.colInk)
	}
}

// TestTheSettingsTabPicksTheTheme covers the picker.
//
// The Settings tab is the one tab with no cursor; now it has one, the active theme is
// marked, and the footer says how to apply.
func TestTheSettingsTabPicksTheTheme(t *testing.T) {
	m := newModel(t)
	m.tab = TabSettings
	m.width = 100
	m.height = 40

	out := StripANSI(m.View().Content)
	if !strings.Contains(out, "THEME") {
		t.Fatalf("the Settings tab does not list the themes:\n%s", out)
	}
	if !strings.Contains(out, "Catppuccin Mocha") {
		t.Fatalf("the theme list is short:\n%s", out)
	}
	// The active theme is marked.
	if !strings.Contains(out, "● Kitchen") {
		t.Fatalf("the active theme is not marked:\n%s", out)
	}
	// The footer says how to move and how to apply.
	if !strings.Contains(out, "apply theme") {
		t.Fatalf("the footer does not say how to apply a theme:\n%s", out)
	}

	// The cursor moves, and the marked row follows it.
	m.Update(tea.KeyPressMsg{Code: 'j'})
	if m.selTheme != 1 {
		t.Fatalf("the cursor is at %d, want 1", m.selTheme)
	}
	if m.currentSel() != 1 {
		t.Fatalf("the marked row is at %d, want 1", m.currentSel())
	}
	out = StripANSI(m.View().Content)
	if !strings.Contains(out, "enter apply theme") {
		t.Fatalf("the footer does not name the key:\n%s", out)
	}

	// And enter applies it.
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.theme.name != "Catppuccin Latte" {
		t.Fatalf("the theme is %q, want Catppuccin Latte", m.theme.name)
	}
	out = StripANSI(m.View().Content)
	if !strings.Contains(out, "● Catppuccin Latte") {
		t.Fatalf("the applied theme is not marked:\n%s", out)
	}
}

// TestTheThemeSurvivesAReload covers the reload.
//
// The theme is the chef's setting, not the plan's: reloading the plan must not take it
// away.
func TestTheThemeSurvivesAReload(t *testing.T) {
	m := newModel(t)
	mocha := paletteNamed(t, "Catppuccin Mocha")
	m.theme = newTheme(mocha)

	msg := m.reload()()
	loaded, ok := msg.(reloadMsg)
	if !ok {
		t.Fatal("unexpected reload result")
	}
	if loaded.err != nil {
		t.Fatal(loaded.err)
	}
	// The live path: the message goes through the model's own Update.
	m.Update(loaded)

	if m.theme.name != "Catppuccin Mocha" {
		t.Fatalf("the theme is %q after a reload, want Catppuccin Mocha", m.theme.name)
	}
}

// TestNoPackageLevelStyles covers the guard.
//
// A theme the chef can pick cannot live in package-level variables: they cannot be
// swapped at runtime.
func TestNoPackageLevelStyles(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for i, line := range strings.Split(string(raw), "\n") {
			if !strings.HasPrefix(line, "var ") && !strings.HasPrefix(line, "func ") {
				continue
			}
			rest := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "var "), "func "))
			if strings.HasPrefix(rest, "sty") || strings.HasPrefix(rest, "col") {
				t.Errorf("%s:%d has a package-level style or colour: %s", f, i+1, line)
			}
		}
	}
}

// sameColor reports whether two colours are the same colour.
func sameColor(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}

// paletteNamed is the palette with the given name.
func paletteNamed(t *testing.T, name string) palette {
	t.Helper()
	for _, p := range palettes() {
		if p.name == name {
			return p
		}
	}
	t.Fatalf("there is no theme called %q", name)
	return palette{}
}
