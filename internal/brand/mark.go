// Package brand holds the mark — the artwork — and nothing else.
//
// Two surfaces need the mark: the console's header and the Sous-Chef tab, and the
// CLI's help. Neither may import the other, and the domain must stay free of UI —
// so the art lives in its own package, and every consumer reads it from here.
//
// The drawings are the source in `images/logo/ascii/`. The wire in `mark_test.go`
// holds this code to them: change a character in the drawing and the suite says so.
package brand

import "unicode/utf8"

// Compact is the mark in one cell — the console's header, which stays one line.
const Compact = "⋈"

// Full is the mark's full form: the five lines of
// `images/logo/ascii/bow-tie-full.txt`, 28 columns at its widest.
var Full = []string{
	"⣼⣿⣿⣶⣶⣦⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣤⣴⣶⣶⣿⣿⣧",
	"⢸⣿⣿⣿⣿⣏⠙⠛⠻⠇⢰⣿⣿⣿⣿⣿⣿⡆⠸⠟⠛⠋⣹⣿⣿⣿⣿⡇",
	"⢸⣿⣿⣿⣿⣿⣿⣿⠶⠀⢸⣿⣿⣿⣿⣿⣿⡇⠀⠶⣿⣿⣿⣿⣿⣿⣿⡇",
	"⢸⣿⣿⣿⣿⣏⣠⣤⣴⡆⠸⣿⣿⣿⣿⣿⣿⠇⢰⣦⣤⣄⣹⣿⣿⣿⣿⡇",
	"⠀⠀⢻⣿⣿⠿⠿⠟⠛⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠛⠻⠿⠿⣿⣿⡟",
}

// Half is the mark's mid-size form: the five lines of
// `images/logo/ascii/bow-tie-half.txt`, 14 columns at its widest.
//
// It is a **second drawing of the same emblem**, not a crop of the full form: the
// halves are drawn at the halves' size, and the wire says so.
var Half = []string{
	"⣼⣿⣶⣤⡀⠀⠀⠀⠀⣀⣴⣶⣿",
	"⢸⣿⣿⠙⠻⢰⣿⣿⣿⠸⠛⣹⣿⣿",
	"⢸⣿⣿⣿⠶⢸⣿⣿⣿⠀⣿⣿⣿⣿",
	"⢸⣿⣿⣠⣴⠸⣿⣿⣿⢰⣤⣹⣿⣿",
	"⠀⢻⣿⠿⠛⠁⠀⠀⠀⠀⠉⠻⠿⣿",
}

// Width is the number of columns a drawing needs — its widest line.
//
// The drawings are braille, and a braille cell is one column wide, so the count of
// runes is the count of columns.
func Width(lines []string) int {
	w := 0
	for _, l := range lines {
		if n := utf8.RuneCountInString(l); n > w {
			w = n
		}
	}
	return w
}
