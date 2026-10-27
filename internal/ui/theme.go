package ui

import "charm.land/lipgloss/v2"

// palette is a set of ten colours.
//
// The background is set explicitly rather than inherited, because the text colours are
// tuned for the palette's own ground and a light terminal would make them unreadable.
// Bubbletea resets the terminal background when the console exits.
type palette struct {
	name string

	colBg     string
	colInk    string
	colDim    string
	colFaint  string
	colAccent string
	colGood   string
	colWarn   string
	colBad    string
	colSelBG  string
	colRule   string
}

// theme is a palette and the styles derived from it.
//
// It lives on the Model, not in the package: a theme the chef can pick has to be
// swappable at runtime, and a package-level variable cannot be.
type theme struct {
	palette

	styTitle lipgloss.Style
	styDim   lipgloss.Style
	styFaint lipgloss.Style
	styAcc   lipgloss.Style
	styGood  lipgloss.Style
	styWarn  lipgloss.Style
	styBad   lipgloss.Style

	styTabActive lipgloss.Style
	styTabIdle   lipgloss.Style

	stySec lipgloss.Style

	stySelected lipgloss.Style

	styRule lipgloss.Style

	styKey lipgloss.Style
}

// newTheme derives the styles from a palette's colours.
func newTheme(p palette) theme {
	t := theme{palette: p}
	t.styTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.colInk))
	t.styDim = lipgloss.NewStyle().Foreground(lipgloss.Color(p.colDim))
	t.styFaint = lipgloss.NewStyle().Foreground(lipgloss.Color(p.colFaint))
	t.styAcc = lipgloss.NewStyle().Foreground(lipgloss.Color(p.colAccent))
	t.styGood = lipgloss.NewStyle().Foreground(lipgloss.Color(p.colGood))
	t.styWarn = lipgloss.NewStyle().Foreground(lipgloss.Color(p.colWarn))
	t.styBad = lipgloss.NewStyle().Foreground(lipgloss.Color(p.colBad))

	t.styTabActive = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.colAccent))
	t.styTabIdle = lipgloss.NewStyle().Foreground(lipgloss.Color(p.colDim))

	t.stySec = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.colAccent))

	t.stySelected = lipgloss.NewStyle().Background(lipgloss.Color(p.colSelBG)).Foreground(lipgloss.Color(p.colInk))

	t.styRule = lipgloss.NewStyle().Foreground(lipgloss.Color(p.colRule))

	t.styKey = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.colAccent))
	return t
}

// The palettes.
//
// The values are the published ones from the official palettes; they are not adjusted
// by eye. The console's own 256-colour palette is the first of them and stays the
// default, so a console nobody told otherwise looks exactly as it did before the theme
// became the chef's to choose.

// kitchen is the console's own palette, unchanged.
func kitchen() palette {
	return palette{
		name:      "Kitchen",
		colBg:     "235",
		colInk:    "252",
		colDim:    "245",
		colFaint:  "240",
		colAccent: "39",
		colGood:   "42",
		colWarn:   "214",
		colBad:    "203",
		colSelBG:  "237",
		colRule:   "238",
	}
}

// catppuccinLatte is Catppuccin Latte — light.
func catppuccinLatte() palette {
	return palette{
		name:      "Catppuccin Latte",
		colBg:     "#eff1f5",
		colInk:    "#4c4f69",
		colDim:    "#6c6f85",
		colFaint:  "#9ca0b0",
		colAccent: "#8839ef",
		colGood:   "#40a02b",
		colWarn:   "#df8e1d",
		colBad:    "#d20f39",
		colSelBG:  "#ccd0da",
		colRule:   "#bcc0cc",
	}
}

// catppuccinFrappe is Catppuccin Frappé — dark.
func catppuccinFrappe() palette {
	return palette{
		name:      "Catppuccin Frappé",
		colBg:     "#303446",
		colInk:    "#c6d0f5",
		colDim:    "#a5adce",
		colFaint:  "#737994",
		colAccent: "#ca9ee6",
		colGood:   "#a6d189",
		colWarn:   "#e5c890",
		colBad:    "#e78284",
		colSelBG:  "#414559",
		colRule:   "#51576d",
	}
}

// catppuccinMacchiato is Catppuccin Macchiato — dark.
func catppuccinMacchiato() palette {
	return palette{
		name:      "Catppuccin Macchiato",
		colBg:     "#24273a",
		colInk:    "#cad3f5",
		colDim:    "#a5adcb",
		colFaint:  "#6e738d",
		colAccent: "#c6a0f6",
		colGood:   "#a6da95",
		colWarn:   "#eed49f",
		colBad:    "#ed8796",
		colSelBG:  "#363a4f",
		colRule:   "#494d64",
	}
}

// catppuccinMocha is Catppuccin Mocha — dark.
func catppuccinMocha() palette {
	return palette{
		name:      "Catppuccin Mocha",
		colBg:     "#1e1e2e",
		colInk:    "#cdd6f4",
		colDim:    "#a6adc8",
		colFaint:  "#6c7086",
		colAccent: "#cba6f7",
		colGood:   "#a6e3a1",
		colWarn:   "#f9e2af",
		colBad:    "#f38ba8",
		colSelBG:  "#313244",
		colRule:   "#45475a",
	}
}

// dracula is the Dracula palette — dark.
func dracula() palette {
	return palette{
		name:      "Dracula",
		colBg:     "#282a36",
		colInk:    "#f8f8f2",
		colDim:    "#6272a4",
		colFaint:  "#6272a4",
		colAccent: "#bd93f9",
		colGood:   "#50fa7b",
		colWarn:   "#ffb86c",
		colBad:    "#ff5555",
		colSelBG:  "#44475a",
		colRule:   "#6272a4",
	}
}

// alucard is Alucard, Dracula's light companion — light.
func alucard() palette {
	return palette{
		name:      "Alucard",
		colBg:     "#fffbeb",
		colInk:    "#1f1f1f",
		colDim:    "#6c664b",
		colFaint:  "#6c664b",
		colAccent: "#644ac9",
		colGood:   "#14710a",
		colWarn:   "#a34d14",
		colBad:    "#cb3a2a",
		colSelBG:  "#cfcfde",
		colRule:   "#6c664b",
	}
}

// palettes are the palettes the chef can pick, in the order the Settings tab lists
// them.
//
// Kitchen is first and is the default: a console nobody told otherwise looks exactly
// as it did before the theme became the chef's to choose.
func palettes() []palette {
	return []palette{
		kitchen(),
		catppuccinLatte(),
		catppuccinFrappe(),
		catppuccinMacchiato(),
		catppuccinMocha(),
		dracula(),
		alucard(),
	}
}

// lifecycleStyle colours a lifecycle state.
func (t theme) lifecycleStyle(l string) lipgloss.Style {
	switch l {
	case "draft", "groomed":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(t.colDim))
	case "ready":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(t.colWarn))
	case "approved", "executing":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(t.colAccent))
	case "executed", "remediation", "reviewed":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(t.colGood))
	case "closed":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(t.colFaint))
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(t.colInk))
}

// criterionStyle colours a criterion status.
func (t theme) criterionStyle(s string) lipgloss.Style {
	switch s {
	case "passed":
		return t.styGood
	case "failed":
		return t.styBad
	case "rejected":
		return t.styWarn
	case "deferred":
		return t.styDim
	}
	return t.styDim
}
